---
title: Kubernetes 证书管理
linkTitle: Kubernetes 证书管理
weight: 1
---

# 概述

> 参考：
>
> -

# PKI 证书和要求

> 参考：
>
> - [官方文档,入门-最佳实践-PKI 证书和要求](https://kubernetes.io/docs/setup/best-practices/certificates/)

Kubernetes 需要 PKI 证书才能进行基于 TLS 的身份验证。如果你是使用 [kubeadm](https://kubernetes.io/zh/docs/reference/setup-tools/kubeadm/) 安装的 Kubernetes， 则会自动生成集群所需的证书。你还可以生成自己的证书。 例如，不将私钥存储在 API 服务器上，可以让私钥更加安全。此页面说明了集群必需的证书。

## 集群是如何使用证书的

Kubernetes 需要 PKI 才能执行以下操作：

- Kubelet 的客户端证书，用于 API 服务器身份验证
- API 服务器端点的证书
- 集群管理员的客户端证书，用于 API 服务器身份认证
- API 服务器的客户端证书，用于和 Kubelet 的会话
- API 服务器的客户端证书，用于和 etcd 的会话
- 控制器管理器的客户端证书/kubeconfig，用于和 API 服务器的会话
- 调度器的客户端证书/kubeconfig，用于和 API 服务器的会话
- [前端代理](https://kubernetes.io/zh/docs/tasks/extend-kubernetes/configure-aggregation-layer/) 的客户端及服务端证书

> **说明：** 只有当你运行 kube-proxy 并要支持 [扩展 API 服务器](https://kubernetes.io/zh/docs/tasks/extend-kubernetes/setup-extension-api-server/) 时，才需要 `front-proxy` 证书

etcd 还实现了双向 TLS 来对客户端和对其他对等节点进行身份验证。

## 证书存放的位置

如果你是通过 kubeadm 安装的 Kubernetes，所有证书都存放在 `/etc/kubernetes/pki` 目录下。本文所有相关的路径都是基于该路径的相对路径。

## 手动配置证书

如果你不想通过 kubeadm 生成这些必需的证书，你可以通过下面两种方式之一来手动创建他们。

### 单根 CA

你可以创建一个单根 CA，由管理员控制器它。该根 CA 可以创建多个中间 CA，并将所有进一步的创建委托给 Kubernetes。
需要这些 CA：

| 路径                   | 默认 CN                   | 描述                                                                                                |
| ---------------------- | ------------------------- | --------------------------------------------------------------------------------------------------- |
| ca.crt,key             | kubernetes-ca             | Kubernetes 通用 CA                                                                                  |
| etcd/ca.crt,key        | etcd-ca                   | 与 etcd 相关的所有功能                                                                              |
| front-proxy-ca.crt,key | kubernetes-front-proxy-ca | 用于 [前端代理](https://kubernetes.io/zh/docs/tasks/extend-kubernetes/configure-aggregation-layer/) |

上面的 CA 之外，还需要获取用于服务账户管理的密钥对，也就是 `sa.key` 和 `sa.pub`。

### 所有的证书

如果你不想将 CA 的私钥拷贝至你的集群中，你也可以自己生成全部的证书。
需要这些证书：

| 默认 CN                       | 父级 CA                   | O (位于 Subject 中) | 类型           | 主机 (SAN)                                          |
| ----------------------------- | ------------------------- | ------------------- | -------------- | --------------------------------------------------- |
| kube-etcd                     | etcd-ca                   |                     | server, client | `localhost`, `127.0.0.1`                            |
| kube-etcd-peer                | etcd-ca                   |                     | server, client | `<hostname>`, `<Host_IP>`, `localhost`, `127.0.0.1` |
| kube-etcd-healthcheck-client  | etcd-ca                   |                     | client         |                                                     |
| kube-apiserver-etcd-client    | etcd-ca                   | system:masters      | client         |                                                     |
| kube-apiserver                | kubernetes-ca             |                     | server         | `<hostname>`, `<Host_IP>`, `<advertise_IP>`, `[1]`  |
| kube-apiserver-kubelet-client | kubernetes-ca             | system:masters      | client         |                                                     |
| front-proxy-client            | kubernetes-front-proxy-ca |                     | client         |                                                     |

\[1]: 用来连接到集群的不同 IP 或 DNS 名 （就像 [kubeadm](https://kubernetes.io/zh/docs/reference/setup-tools/kubeadm/) 为负载均衡所使用的固定 IP 或 DNS 名，`kubernetes`、`kubernetes.default`、`kubernetes.default.svc`、 `kubernetes.default.svc.cluster`、`kubernetes.default.svc.cluster.local`）。
其中，`kind` 对应一种或多种类型的 [x509 密钥用途](https://godoc.org/k8s.io/api/certificates/v1beta1#KeyUsage)：

| kind   | 密钥用途                       |
| ------ | ------------------------------ |
| server | 数字签名、密钥加密、服务端认证 |
| client | 数字签名、密钥加密、客户端认证 |

> **说明：**
> 上面列出的 Hosts/SAN 是推荐的配置方式；如果需要特殊安装，则可以在所有服务器证书上添加其他 SAN。
> **说明：**
> 对于 kubeadm 用户：
>
> - 不使用私钥，将证书复制到集群 CA 的方案，在 kubeadm 文档中将这种方案称为外部 CA。
> - 如果将以上列表与 kubeadm 生成的 PKI 进行比较，你会注意到，如果使用外部 etcd，则不会生成 `kube-etcd`、`kube-etcd-peer` 和 `kube-etcd-healthcheck-client` 证书。

### 证书路径

证书应放置在建议的路径中（以便 [kubeadm](https://kubernetes.io/zh/docs/reference/setup-tools/kubeadm/)使用）。无论使用什么位置，都应使用给定的参数指定路径。

| 默认 CN                       | 建议的密钥路径               | 建议的证书路径               | 命令                    | 密钥参数                   | 证书参数                                                      |
| ----------------------------- | ---------------------------- | ---------------------------- | ----------------------- | -------------------------- | ------------------------------------------------------------- |
| etcd-ca                       | etcd/ca.key                  | etcd/ca.crt                  | kube-apiserver          |                            | --etcd-cafile                                                 |
| kube-apiserver-etcd-client    | apiserver-etcd-client.key    | apiserver-etcd-client.crt    | kube-apiserver          | --etcd-keyfile             | --etcd-certfile                                               |
| kubernetes-ca                 | ca.key                       | ca.crt                       | kube-apiserver          |                            | --client-ca-file                                              |
| kubernetes-ca                 | ca.key                       | ca.crt                       | kube-controller-manager | --cluster-signing-key-file | --client-ca-file, --root-ca-file, --cluster-signing-cert-file |
| kube-apiserver                | apiserver.key                | apiserver.crt                | kube-apiserver          | --tls-private-key-file     | --tls-cert-file                                               |
| kube-apiserver-kubelet-client | apiserver-kubelet-client.key | apiserver-kubelet-client.crt | kube-apiserver          | --kubelet-client-key       | --kubelet-client-certificate                                  |
| front-proxy-ca                | front-proxy-ca.key           | front-proxy-ca.crt           | kube-apiserver          |                            | --requestheader-client-ca-file                                |
| front-proxy-ca                | front-proxy-ca.key           | front-proxy-ca.crt           | kube-controller-manager |                            | --requestheader-client-ca-file                                |
| front-proxy-client            | front-proxy-client.key       | front-proxy-client.crt       | kube-apiserver          | --proxy-client-key-file    | --proxy-client-cert-file                                      |
| etcd-ca                       | etcd/ca.key                  | etcd/ca.crt                  | etcd                    |                            | --trusted-ca-file, --peer-trusted-ca-file                     |
| kube-etcd                     | etcd/server.key              | etcd/server.crt              | etcd                    | --key-file                 | --cert-file                                                   |
| kube-etcd-peer                | etcd/peer.key                | etcd/peer.crt                | etcd                    | --peer-key-file            | --peer-cert-file                                              |
| etcd-ca                       |                              | etcd/ca.crt                  | etcdctl                 |                            | --cacert                                                      |
| kube-etcd-healthcheck-client  | etcd/healthcheck-client.key  | etcd/healthcheck-client.crt  | etcdctl                 | --key                      | --cert                                                        |

注意事项同样适用于服务帐户密钥对：

| 私钥路径 | 公钥路径 | 命令                    | 参数                               |
| -------- | -------- | ----------------------- | ---------------------------------- |
| sa.key   |          | kube-controller-manager | --service-account-private-key-file |
|          | sa.pub   | kube-apiserver          | --service-account-key-file         |

## 为用户帐户配置证书

你必须手动配置以下管理员帐户和服务帐户：

| 文件名                  | 凭据名称                   | 默认 CN                               | O (位于 Subject 中) |
| ----------------------- | -------------------------- | ------------------------------------- | ------------------- |
| admin.conf              | default-admin              | kubernetes-admin                      | system:masters      |
| kubelet.conf            | default-auth               | system:node:`<nodeName>` （参阅注释） | system:nodes        |
| controller-manager.conf | default-controller-manager | system:kube-controller-manager        |                     |
| scheduler.conf          | default-scheduler          | system:kube-scheduler                 |                     |

> **说明：** `kubelet.conf` 中 `<nodeName>` 的值 **必须** 与 kubelet 向 apiserver 注册时提供的节点名称的值完全匹配。 有关更多详细信息，请阅读[节点授权](https://kubernetes.io/zh/docs/reference/access-authn-authz/node/)。

1. 对于每个配置，请都使用给定的 CN 和 O 生成 x509 证书/密钥偶对。
2. 为每个配置运行下面的 `kubectl` 命令：

```bash
KUBECONFIG=<filename> kubectl config set-cluster default-cluster --server=https://<host ip>:6443 --certificate-authority <path-to-kubernetes-ca> --embed-certs
KUBECONFIG=<filename> kubectl config set-credentials <credential-name> --client-key <path-to-key>.pem --client-certificate <path-to-cert>.pem --embed-certs
KUBECONFIG=<filename> kubectl config set-context default-system --cluster default-cluster --user <credential-name>
KUBECONFIG=<filename> kubectl config use-context default-system
```

这些文件用途如下：

| 文件名                  | 命令                    | 说明                                                       |
| ----------------------- | ----------------------- | ---------------------------------------------------------- |
| admin.conf              | kubectl                 | 配置集群的管理员                                           |
| kubelet.conf            | kubelet                 | 集群中的每个节点都需要一份                                 |
| controller-manager.conf | kube-controller-manager | 必需添加到 `manifests/kube-controller-manager.yaml` 清单中 |
| scheduler.conf          | kube-scheduler          | 必需添加到 `manifests/kube-scheduler.yaml` 清单中          |

# Certificate(证书) # 使用证书对集群中的客户端与服务端进行认证

> 参考：
>
> - 官方文档：<https://kubernetes.io/docs/setup/best-practices/certificates/>

etcd 与 etcd 之间，etcd 与 apiserver，apiserver 与 kubelet、scheduler、controller-manager、kube-proxy 等之间的认证，还有 calico 与 apiserver 等等各种组件与组件之间基本都需要认证，认证可以通过多种方式进行，比如证书、token、key/val 对，账号密码等等等

Cluster 中各组件互相通信所用到的 Certificate

- ETCD 的证书，这是集群中的其中一套证书：api-server 作为客户端与服务端 etcd 通信，etcd 集群之间互相对等通信
  - **ca.crt**(证书 CN：etcd-ca) # 给 apiserver 发客户端证书，给 etcd 发服务端证书以及对等证书
  - **peer.crt**(证书 CN：HostName) # etcd 集群各节点属于对等节点，使用 peer 类型证书(一般分为 server 证书和 client 证书，但是 etcd 集群之间不存在服务端和客户端的区别)
  - **apiserver-etcd-client.crt**(证书 CN：kube-apiserver-etcd-client) # 与 server.crt 证书对应。apiserver 作为 etcd 的客户端所用的证书
  - **server.crt**(证书 CN：HostName) # 与 apiserver-etcd-client.crt 证书对应。etcd 作为 apiserver 的服务端所用的证书
- 集群组件间的证书：kube-apiserver 作为服务端与 kubectl，controller-manager，scheduler，kubelet，kube-proxy 通信
  - **ca.crt**(证书 CN：kubernetes) # 给 apiserver 发服务端证书，给其余组件发客户端证书
  - **apiserver.crt**(证书 CN：kube-apiserver)
  - **admin.conf** # 一个在与集群通信时具有最高权限的 user 的认证配置
  - **controller-manager.conf** # KubeConfig 文件，controller-manager 与 apiserver 通信时的认证配置信息
  - **scheduler.conf** # KubeConfig 文件，scheduler 与 apiserver 通信时的认证配置信息
  - **kubelet.conf** # KubeConfig 文件，kubelet 与 apiserver 通信时的认证配置信息
    - kube-apiserver 作为客户端与 kubelet-api 通信，每个节点启动的时候 kubelet-api 的证书会自动从 kubernets 的 ca 证书那里获取自己的 ca 证书
  - **apiserver-kubelet-client.crt**(证书 CN：kube-apiserver-kubelet-client) #
  - **kubelet.crt**(证书 CN：master0@1544020244) #
- 前端代理证书：给用户自定义的 apiserver 使用的证书，kube-aggregator 作为服务端与 extension-apiserver 通信
  - **ca.crt(front-proxy)** # 给自定义的 apiserver 发证书
- 其他证书
  - **sa.key 与 sa.pub** # 用于为集群中所有 ServiceAccount 资源签署 jwt token
