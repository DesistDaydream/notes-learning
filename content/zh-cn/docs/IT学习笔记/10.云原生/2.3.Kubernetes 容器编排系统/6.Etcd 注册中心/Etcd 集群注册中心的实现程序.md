---
title: Etcd 集群注册中心的实现程序
---

# 概述

> 参考：
> - etcd 详细用法详见 ETCD 简介

作为 kubernetes 集群的存储系统使用(也可以算是集群的注册中心)，保存了集群的所有配置信息，需要高可用，如果需要在生产环境下使用，则需要在单独部署

Note：etcd 只接收 apiserver 的请求

每个 etcd 一般使用两个端口进行工作，一个端口面向客户端提供服务(port/2379)，另一个端口集群内部通信(port/2380)

想要正常运行 ETCD，需要注意以下几点：

1. 配置 etcd 的证书，如果不使用证书，则不法分子有可能直接去修改 etcd 数据

2. ca.crt(证书 CN：etcd-ca) #给 apiserver 发客户端证书，给 etcd 发服务端证书以及对等证书

3. peer.crt(证书 CN：HostName) #etcd 集群各节点属于对等节点，使用 peer 类型证书(一般分为 server 证书和 client 证书，但是 etcd 集群之间不存在服务端和客户端的区别)

4. apiserver-etcd-client.crt(证书 CN：kube-apiserver-etcd-client) #与 server.crt 证书对应。apiserver 作为 etcd 的客户端所用的证书

5. server.crt(证书 CN：HostName) #与 apiserver-etcd-client.crt 证书对应。etcd 作为 apiserver 的服务端所用的证书

6. 修改 etcd 的配置文件

# Etcd Metrics

详见：k8s 主要组件 metrics 获取指南

# etcdctl 命令行工具使用说明

详见：etcdctl 命令行工具

k8s 集群中使用 etcdctl 的技巧：可以创建一个别名，以便后续使用 etcdctl 命令的时候，可以使用别名来自动加载证书，比如下面

测试环境

1. echo 'alias etcdctl="kubectl exec -it -n kube-system etcd-master-1.tj-test -- etcdctl --key=/etc/kubernetes/pki/etcd/peer.key --cert=/etc/kubernetes/pki/etcd/peer.crt --cacert=/etc/kubernetes/pki/etcd/ca.crt --endpoints=172.38.40.212:2379,172.38.40.213:2379,172.38.40.214:2379"' >> /etc/bashrc

测试环境外部 etcd

1. echo 'alias etcdctl="export ETCDCTL_API=3; etcdctl --key=/etc/ssl/etcd/ssl/admin-master-1.tj-test-key.pem --cert=/etc/ssl/etcd/ssl/admin-master-1.tj-test.pem --cacert=/etc/ssl/etcd/ssl/ca.pem --endpoints=172.38.40.212:2379,172.38.40.213:2379,172.38.40.214:2379"' >> /etc/bashrc

无锡环境

1. echo 'alias etcdctl="kubectl exec -it -n kube-system etcd-master-1.wx -- etcdctl --key=/etc/kubernetes/pki/etcd/peer.key --cert=/etc/kubernetes/pki/etcd/peer.crt --cacert=/etc/kubernetes/pki/etcd/ca.crt --endpoints=172.40.0.3:2379,172.40.0.4:2379,172.40.0.5:2379"' >> /root/.bashrc

etcdctl 工具也可以直接从 etcd 容器的目录 /usr/local/bin 下，直接拷贝到宿主机上使用。这样使用起来也更方便

1. echo 'alias etcdctl="export ETCDCTL_API=3; etcdctl --key=/etc/kubernetes/pki/etcd/peer.key --cert=/etc/kubernetes/pki/etcd/peer.crt --cacert=/etc/kubernetes/pki/etcd/ca.crt --endpoints=172.38.40.212:2379,172.38.40.213:2379,172.38.40.214:2379"' >> /root/.bashrc

# Etcd 参数详解

参考：Etcd 命令行 flag 详解
下面是

    apiVersion: v1
    kind: Pod
    metadata:
      creationTimestamp: null
      labels:
        component: etcd
        tier: control-plane
      name: etcd
      namespace: kube-system
    spec:
      containers:
      - command:
        - etcd
        - --advertise-client-urls=https://10.10.100.101:2379
        - --cert-file=/etc/kubernetes/pki/etcd/server.crt
        - --client-cert-auth=true
        - --data-dir=/var/lib/etcd
        - --initial-advertise-peer-urls=https://10.10.100.101:2380
        - --initial-cluster=master-1.test.tjiptv.net=https://10.10.100.101:2380
        - --key-file=/etc/kubernetes/pki/etcd/server.key
        - --listen-client-urls=https://127.0.0.1:2379,https://10.10.100.101:2379
        - --listen-metrics-urls=http://127.0.0.1:2381
        - --listen-peer-urls=https://10.10.100.101:2380
        - --name=master-1.test.tjiptv.net
        - --peer-cert-file=/etc/kubernetes/pki/etcd/peer.crt
        - --peer-client-cert-auth=true
        - --peer-key-file=/etc/kubernetes/pki/etcd/peer.key
        - --peer-trusted-ca-file=/etc/kubernetes/pki/etcd/ca.crt
        - --snapshot-count=10000
        - --trusted-ca-file=/etc/kubernetes/pki/etcd/ca.crt
        image: k8s.gcr.io/etcd:3.3.15-0
        imagePullPolicy: IfNotPresent
        livenessProbe:
          failureThreshold: 8
          httpGet:
            host: 127.0.0.1
            path: /health
            port: 2381
            scheme: HTTP
          initialDelaySeconds: 15
          timeoutSeconds: 15
        name: etcd
        resources: {}
        volumeMounts:
        - mountPath: /var/lib/etcd
          name: etcd-data
        - mountPath: /etc/kubernetes/pki/etcd
          name: etcd-certs
      hostNetwork: true
      priorityClassName: system-cluster-critical
      volumes:
      - hostPath:
          path: /etc/kubernetes/pki/etcd
          type: DirectoryOrCreate
        name: etcd-certs
      - hostPath:
          path: /var/lib/etcd
          type: DirectoryOrCreate
        name: etcd-data
