---
title: 基于 DNS 的 Kubernetes 服务发现的规范
---

# 概述

参考：[官方文档](https://github.com/kubernetes/dns/blob/master/docs/specification.md)

任何基于 DNS 的以实现 Kubernetes 服务发现的实现工具，必须提供本规范描述的 Resource Record

DNS Records 类型包含以下几种：A/AAAA、SRV、PTR、CHAME。不同的 k8s 资源，其可用的 Record Type 也不相同。

Resource Records 规范

在下面的描述中，有几个占位符，表示如下含义：

- **<service>** # service 对象的名字

- **<ns>** # namesapce 的名字

- **<zone>** # 集群域名(默认为 cluster.local)。

  - 集群域名可以通过 kubelet 的配置文件 clusterDomain 字段定义。

- **<ttl>** # 一条 Record 的标准 DNS 存活时间

## ClusterIP 类型 Service 的 RR 格式

参考：[官方文档](https://github.com/kubernetes/dns/blob/master/docs/specification.md#23---records-for-a-service-with-clusterip)

现在假定一个 Service 对象名为 `<service>`，在名为 `<ns>` 名称空间，该 Service 对象的 CLUSTER-IP 为 `<cluster-ip>`。则 DNS 实现程序必须具有下列几种类型的 Records(记录)。

### A/AAAA Record

如果 <service> 对象具有 <cluster-ip> 且为 IPv4 地址。

**A 记录 格式：**`**<service>.<ns>.svc.<zone>. <ttl> IN A <cluster-ip>**`

- 请求样例：

  - kubernetes.default.svc.cluster.local. IN A

- 响应样例：

  - kubernetes.default.svc.cluster.local. 4 IN A 10.3.0.1

如果 <service> 对象具有 <cluster-ip> 且为 IPv6 地址。

**AAAA 记录 格式：**`**<service>.<ns>.svc.<zone>. <ttl> IN AAAA <cluster-ip>**`

- 请求样例：

  - kubernetes.default.svc.cluster.local. IN AAAA

- 响应样例：

  - kubernetes.default.svc.cluster.local. 4 IN AAAA 2001:db8::1

### SRV Record

### PTR Record

### 总结

Service 域名与 Service 的 CLUSTER-IP 具有关联关系

Headless 类型的 Service 的 RR 格式

参考：[官方文档](https://github.com/kubernetes/dns/blob/master/docs/specification.md#24---records-for-a-headless-service)
现在假定一个 Headless 类型的 Service 对象名为 `<service>`，在名为 `<ns>` 名称空间，没有 CLUSTER-IP，Service 关联的每个 Pod 的主机名为 `<hostname>`，每个 Pod 的 IP 地址为 `<endpoint-ip>`。则 DNS 实现程序必须具有下列几种类型的 Records(记录)。

### A/AAAA Record

对于 Headless 类型的 Service，Service 对应的每个 endpoint 的 `<endpoint-ip>`(也就是 Pod 的 IP)，都必须有一个 A 记录。如果 endpoint 不存在，则 DNS-Rcode 应该是 NXDOMAIN。并且，<endpoint-ip> 会对应两种 A 记录

- 通过 Service 名称解析到 endpoint-ip 的 Record

- 通过 Pod 主机名解析到 endpoint-ip 的 Record。注意，该记录只有在 Statefulset 生成的 Pod 才有效果。

**针对 Service 名称的 A 记录格式：**`**<service>.<ns>.svc.<zone>. <ttl> IN A <endpoint-ip>**`

- 请求样例：

  - headless.default.svc.cluster.local. IN A

- 响应样例：

  - headless.default.svc.cluster.local. 4 IN A 10.3.0.1

  - headless.default.svc.cluster.local. 4 IN A 10.3.0.2

  - headless.default.svc.cluster.local. 4 IN A 10.3.0.3

  - ..... 有多少个 endpoint 就有多少个响应

注意：Headless 与 ClusterIP 类型的 Service 最大的不同在于下面的 Resource Record。说白了就是多了个子域名

除了上面的 A 记录，还必须要有这样一种 A 记录。Service 对应的 endpoint 的 `<endpoint-ip>`(也就是 Pod 的 IP)，还必须有另一个 A 记录，该记录是 `<endpoint-ip>` 与 包含 `<hostname>`(也就是 Pod 的主机名) 的域名的对应关系。如果一个给定的 hostname 具有多个 IPv4 地址，则每个 IP 都要有一个 A 记录。

**针对 Pod 主机名的 A 记录格式：**`**<hostname>.<service>.<ns>.svc.<zone>. <ttl> IN A <endpoint-ip>**`

- 请求样例：

  - my-pet.headless.default.svc.cluster.local. IN A

- 响应样例：

  - my-pet.headless.default.svc.cluster.local. 4 IN A 10.3.0.100

**AAAA 记录格式**

略，把 A 改为 AAAA，IPv4 改为 IPv6，与 ClusterIP 类型的 Service 的 A 和 AAAA 记录的对应关系一样。

### SRV Record

### PTR Record

### 总结

Service 域名 与 Service 关联的后端 Pod 的 IP 具有解析关系。

带主机名的域名与 Service 关联的后端 Pod 的 IP 具有解析关系

## ExternalName 类型的 Service 的 RR 格式

1. CNAME Record：<service>.<ns>.svc.<zone>. <ttl> IN CNAME <ExtName>.

   1. 请求样例：foo.default.svc.cluster.local. IN A

   2. 响应样例：

      1. foo.default.svc.cluster.local. 10 IN CNAME www.example.com.

      2. www.example.com. 28715 IN A 192.0.2.53

示例：

    # 在default名称空间中有一个名为 myapp 的 servcie ，ip 为 10.108.255.155
    [root@master-1 ~]# kubectl get svc
    NAME         TYPE        CLUSTER-IP       EXTERNAL-IP   PORT(S)        AGE
    kubernetes   ClusterIP   10.96.0.1        <none>        443/TCP        239d
    myapp        NodePort    10.108.255.155   <none>        80:52202/TCP   6d10h
    # 这个 default 名称空间下的 myapp 的 service，在coredns中记录的域名为 myapp.default.svc.cluster.local ，对应的ip为 10.108.255.155
    [root@master-1 ~]# kubectl exec -it -n test myapp-5ccfc89896-zxljh -- /bin/bash
    [root@myapp-5ccfc89896-zxljh /]# ping myapp.default
    PING myapp.default.svc.cluster.local (10.108.255.155) 56(84) bytes of data.
    64 bytes from myapp.default.svc.cluster.local (10.108.255.155): icmp_seq=1 ttl=64 time=0.030 ms

有两种类型的控制器，各自有两种类型的 Service。

    root@lichenhao:~# kubectl get svc -n logging
    NAME                           TYPE        CLUSTER-IP      EXTERNAL-IP   PORT(S)    AGE
    log-bj-test-grafana            ClusterIP   10.103.6.159    <none>        80/TCP     2d4h
    log-bj-test-grafana-headless   ClusterIP   None            <none>        80/TCP     4h37m
    log-bj-test-loki               ClusterIP   10.102.159.19   <none>        3100/TCP   2d4h
    log-bj-test-loki-headless      ClusterIP   None            <none>        3100/TCP   2d4h
    root@lichenhao:~# kubectl get pod -n logging -o wide
    NAME                                   READY   STATUS    RESTARTS   AGE    IP             NODE               NOMINATED NODE   READINESS GATES
    log-bj-test-grafana-7764f5b4d7-28ngk   2/2     Running   0          2d4h   10.244.2.129   master-3.bj-test   <none>           <none>
    log-bj-test-loki-0                     1/1     Running   0          2d4h   10.244.1.211   master-2.bj-test   <none>           <none>

首先实验 ClusterIP 类型的 Service

    ~ # nslookup log-bj-test-grafana.logging.svc.cluster.local
    ......
    Name: log-bj-test-grafana.logging.svc.cluster.local
    Address: 10.103.6.159
    ~ # nslookup log-bj-test-loki.logging.svc.cluster.local
    ......
    Name: log-bj-test-loki.logging.svc.cluster.local
    Address: 10.102.159.19

实验 Headless 类型的 Service

    # 解析 Service 名
    ~ # nslookup log-bj-test-grafana-headless.logging.svc.cluster.local
    ......
    Name: log-bj-test-grafana-headless.logging.svc.cluster.local
    Address: 10.244.2.129
    ~ # nslookup log-bj-test-loki-headless.logging.svc.cluster.local
    ......
    Name: log-bj-test-loki-headless.logging.svc.cluster.local
    Address: 10.244.1.211
    # 解析主机名
    ~ # nslookup log-bj-test-grafana-7764f5b4d7-28ngk.log-bj-test-grafana-headless.logging.svc.cluster.local
    ......
    ** server can't find log-bj-test-grafana-7764f5b4d7-28ngk.log-bj-test-grafana-headless.logging.svc.cluster.local: NXDOMAIN
    ** server can't find log-bj-test-grafana-7764f5b4d7-28ngk.log-bj-test-grafana-headless.logging.svc.cluster.local: NXDOMAIN
    ~ # nslookup log-bj-test-loki-0.log-bj-test-loki-headless.logging.svc.cluster.local
    ......
    Name: log-bj-test-loki-0.log-bj-test-loki-headless.logging.svc.cluster.local
    Address: 10.244.1.211

SRV 记录

    ~ # nslookup _http-metrics._tcp.log-bj-test-loki-headless.logging.svc.cluster.local
    ......
    Name: _http-metrics._tcp.log-bj-test-loki-headless.logging.svc.cluster.local
    Address: 10.244.1.211
