---
title: 更新 APIServer 证书
---

# 概述

> 参考：
>
> - 原文链接：<https://mp.weixin.qq.com/s/bs0urFxOG71nq9K34H1b6Q>

本文我们将了解如何将一个新的 DNS 名称或者 IP 地址添加到 Kubernetes APIServer 使用的 TLS 证书中。在某些情况下默认的证书包含的名称可能不能满足我们的要求，又或者是 APIServer 地址有所变化，都需要重新更新证书。

我们这里的集群是使用 kubeadm 搭建的单 master 集群，使用的也是 kubeadm 在启动集群时创建的默认证书授权 CA，对于其他环境的集群不保证本文也同样适用。

Kubernetes APIServer 使用数字证书来加密 APIServer 的相关流量以及验证到 APIServer 的连接。所以如果我们想使用命令行客户端（比如 kubectl）连接到 APIServer，并且使用的主机名或者 IP 地址不包括在证书的 subject 的备选名称（SAN）列表中的话，访问的时候可能会出错，会提示对指定的 IP 地址或者主机名访问证书无效。要解决这个问题就需要更新证书，使 SAN 列表中包含所有你将用来访问 APIServer 的 IP 地址或者主机名。

# 步骤

## 生成 kubeadm 配置文件

因为集群是使用 kubeadm 搭建的，所以我们可以直接使用 kubeadm 来更新 APIServer 的证书，来保证在 SAN 列表中包含额外的名称。

首页我们一个 kubeadm 的配置文件，如果一开始安装集群的时候你就是使用的配置文件，那么我们可以直接更新这个配置文件，但是如果你没有使用配置文件，直接使用的 kubeadm init 来安装的集群，那么我们可以从集群中获取 kubeadm 的配置信息来创建一个配置文件，因为 kubeadm 会将其配置写入到 kube-system 命名空间下面一个名为 kubeadm-config 的 ConfigMap 中。可以直接执行如下所示的命令将该配置导出：

    $ kubectl -n kube-system get configmap kubeadm-config -o jsonpath='{.data.ClusterConfiguration}' > kubeadm-config.yaml

上面的命令会导出一个名为 kubeadm.yaml 的配置文件，内容如下所示：

```yaml
apiServer:
  extraArgs:
    authorization-mode: Node,RBAC
  timeoutForControlPlane: 4m0s
apiVersion: kubeadm.k8s.io/v1beta2
certificatesDir: /etc/kubernetes/pki
clusterName: kubernetes
controllerManager: {}
dns:
  type: CoreDNS
etcd:
  local:
    dataDir: /var/lib/etcd
imageRepository: registry.aliyuncs.com/k8sxio
kind: ClusterConfiguration
kubernetesVersion: v1.17.11
networking:
  dnsDomain: cluster.local
  podSubnet: 10.244.0.0/16
  serviceSubnet: 10.96.0.0/12
scheduler: {}
```

## 添加 certSANs

上面的配置中并没有列出额外的 SAN 信息，我们要添加一个新的数据，需要在 apiServer 属性下面添加一个 certsSANs 的列表。如果你在启动集群的使用就使用的了 kubeadm 的配置文件，可能里面就已经包含 certSANs 列表了，如果没有我们就需要添加它，比如我们这里要添加一个新的域名 api.k8s.local 以及 ydzs-master2 和 ydzs-master3 这两个主机名和 10.151.30.70、10.151.30.71 这两个新的 IP 地址，那么我们需要在 apiServer 下面添加如下所示的数据：

```yaml
apiServer:
  certSANs:
    - api.k8s.local
    - ydzs-master2
    - ydzs-master3
    - 10.151.30.11
    - 10.151.30.70
    - 10.151.30.71
```

上面我只列出了 apiServer 下面新增的 certSANs 信息，这些信息是包括在标准的 SAN 列表之外的，所以不用担心这里没有添加 kubernetes、kubernetes.default 等等这些信息，因为这些都是标准的 SAN 列表中的。

## 备份老证书

更新完 kubeadm 配置文件后我们就可以更新证书了，首先我们移动现有的 APIServer 的证书和密钥，因为 kubeadm 检测到他们已经存在于指定的位置，它就不会创建新的了。

    $ mv /etc/kubernetes/pki/apiserver.{crt,key} ~

## 生成新证书

然后直接使用 kubeadm 命令生成一个新的证书：

    $ kubeadm init phase certs apiserver --config kubeadm-config.yaml
    W0902 10:05:28.006627     832 validation.go:28] Cannot validate kubelet config - no validator is available
    W0902 10:05:28.006754     832 validation.go:28] Cannot validate kube-proxy config - no validator is available
    [certs] Generating "apiserver" certificate and key
    [certs] apiserver serving cert is signed for DNS names [ydzs-master kubernetes kubernetes.default kubernetes.default.svc kubernetes.default.svc.cluster.local api.k8s.local ydzs-master2 ydzs-master3] and IPs [10.96.0.1 123.59.188.12 10.151.30.11 10.151.30.70 10.151.30.71]

通过上面的命令可以查看到 APIServer 签名的 DNS 和 IP 地址信息，一定要和自己的目标签名信息进行对比，如果缺失了数据就需要在上面的 certSANs 中补齐，重新生成证书。

该命令会使用上面指定的 kubeadm 配置文件为 APIServer 生成一个新的证书和密钥，由于指定的配置文件中包含了 certSANs 列表，那么 kubeadm 会在创建新证书的时候自动添加这些 SANs。

## 重启 kube-apiserver

最后一步是重启 APIServer 来接收新的证书，最简单的方法是直接杀死 APIServer 的容器：

    $ docker ps | grep kube-apiserver | grep -v pause
    7fe227a5dd3c        aa63290ccd50                               "kube-apiserver --ad…"   14 hours ago        Up 14 hours                             k8s_kube-apiserver_kube-apiserver-ydzs-master_kube-system_6aa38ee2d66b7d9b6660a88700d00581_0
    $ docker kill 7fe227a5dd3c
    7fe227a5dd3c

容器被杀掉后，kubelet 会自动重启容器，然后容器将接收新的证书，一旦 APIServer 重启后，我们就可以使用新添加的 IP 地址或者主机名来连接它了，比如我们新添加的 api.k8s.local。

## 验证

要验证证书是否更新我们可以直接去编辑 kubeconfig 文件中的 APIServer 地址，将其更换为新添加的 IP 地址或者主机名，然后去使用 kubectl 操作集群，查看是否可以正常工作。

当然我们可以使用 openssl 命令去查看生成的证书信息是否包含我们新添加的 SAN 列表数据：

    $ openssl x509 -in /etc/kubernetes/pki/apiserver.crt -text
    Certificate:
    ......
            Subject: CN=kube-apiserver
    ......
                X509v3 Subject Alternative Name:
                    DNS:ydzs-master, DNS:kubernetes, DNS:kubernetes.default, DNS:kubernetes.default.svc, DNS:kubernetes.default.svc.cluster.local, DNS:api.k8s.local, DNS:ydzs-master2, DNS:ydzs-master3, IP Address:10.96.0.1, IP Address:123.59.188.12, IP Address:10.151.30.11, IP Address:10.151.30.70, IP Address:10.151.30.71
    ......

## 更新集群配置

如果上面的操作都一切顺利，最后一步是将上面的集群配置信息保存到集群的 kubeadm-config 这个 ConfigMap 中去，这一点非常重要，这样以后当我们使用 kubeadm 来操作集群的时候，相关的数据不会丢失，比如升级的时候还是会带上 certSANs 中的数据进行签名的。

    $ kubeadm config upload from-file --config kubeadm.yaml

使用上面的命令保存配置后，我们同样可以用下面的命令来验证是否保存成功了：

    $ kubectl -n kube-system get configmap kubeadm-config -o yaml

更新 APIServer 证书的名称在很多场景下都会使用到，比如在控制平面前面添加一个负载均衡器，或者添加新的 DNS 名称或 IP 地址来使用控制平面的端点，所以掌握更新集群证书的方法也是非常有必要的。
