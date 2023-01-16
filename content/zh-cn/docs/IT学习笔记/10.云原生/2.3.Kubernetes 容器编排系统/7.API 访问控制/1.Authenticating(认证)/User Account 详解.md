---
title: User Account 详解
---

# 概述

> 参考：
> - [官方文档,概念-配置-使用 kubeconfig 文件访问集群](https://kubernetes.io/docs/concepts/configuration/organize-cluster-access-kubeconfig/)
> - [官方文档,任务-访问集群中的应用程序-配置多集群访问](https://kubernetes.io/docs/tasks/access-application-cluster/configure-access-multiple-clusters/)

**User Account(用户账户，简称 UA)** 使用 KubeConfig 文件进行认证。KubeConfig 是一个允许各个客户端与集群通信时所用到的认证配置文件，由于与 kubernetes 交互的途径只有通过 API Server 这一条途径，所以就相当于 API Server 的各客户端(kubelet、scheduler、controller-manager、kube-proxy 等)与其进行通信时使用的认证、配置文件。

KubeConfig 是对 UserAccount 的扩展，KubeConfig 会创建 UserAccount 并关联到指定的集群上

使用 KubeConfig 的原因：可以不用进行双向证书交换，节省交互开销。仅用于对安全性不那么高的情况，否则依然使用双向认证，比如 etcd 与 apiserver 的交互

1. 首先，Kubeconfig 可以是任意文件名的文件，Kubeconfig 只是一个概念，并以文本文件的形式展示出来。

2. 在开启了 TLS 的集群中，每当与集群交互的时候少不了的是身份认证，使用证书和 token(令牌)两种认证方式是最简单也最通用的认证方式。

以 kubectl 为例，kubectl 只是个 go 编写的可执行程序，只要为 kubectl 配置合适的 KubeConfig，就可以在集群中的任意节点使用。kubectl 默认会从 ~/.kube 目录下查找文件名为 config 的文件，也可以使用 --kubeconfig 命令行标志时指明具体的 KubeConfig 文件。(注意：下文中的用户指的是 kubernetes 中的用户，与 linux 的用户不同)

- 比如
- 当 kubectl 去 get 或者 delete 等资源的时候，相当于是对集群请求执行该指令，而集群是通过 API Server 来接收这些指令的

- 那么首先要确认的是使用 kubectl 进行操作的这个 User 是谁，这个 User 是否有证书来对我发起这些操作。如果我都不认可这个 User，那么我都不会接受这些指令请求，这就是 KubeConfig 的作用

- 然后 KubeConfig 可以指明一个 User 与一个 cluster 绑定，当绑定之后，即证明该 User 可以通过 kubectl 来对该绑定集群的 API Server 发起请求，一个 User 可以绑定多个集群。一个集群也可绑定多个 User

- 当 API Server 认可这个 User 通过 kubectl 发送的请求后，就需要下一步授权来对该请求中指令进行鉴权，鉴别这个 User 是否有权利执行这个指令

总之，KubeConfig 就是访问集群所需认证信息文件。

# KubeConfig Manifests 详解

> 参考：
> - 官方文档中还没有对这个配置文件的描述
> - 代码：<https://github.com/kubernetes/client-go/blob/master/tools/clientcmd/api/types.go>

使用命令修改：

- 配置方法详见 [kubectl 的 config 子命令](https://www.yuque.com/go/doc/33163778)，当使用 --kubeconfig 指定文件时，如果文件不存在，则会自动创建，并包含基本模板
- 直接对文件中的各个字段进行修改

## apiVersion: v1

## kind: Config

## clusters: <\[]Object> # 定义访问指定集群所用的证书、访问入口、名称，可指定多个集群

- **cluster: \<map\[STRING]STRING>** # 集群列表
  - **certificate-authority-data: <STRING>** # 集群的认证信息。一般为集群 ca 证书的 base64 格式的字符串
  - **server: <STRING>** # 集群的入口，一般为 API Server 的 `https://IP:PORT`
  - **insecure-skip-tls-verify: <BOOLEAN>** # 是否跳过验证服务端证书有效性的行为。
- **name: <STRING>** # 指定该集群的名称

## contexts: <\[]Object> # 指名用户与集群的绑定关系。

比如有一台主机作为客户端(kubelet)，想控制多个 k8s 的集群，为了让一个 kubectl 控制多个集群且多个不同用户账号可以访问多个不同的集群。

- **context: \<map\[STRING]STRING>** # 上下文列表
  - **cluster: <STRING>** # 指明与 user 关联的 cluster
  - **user: <STRING>** # 指明与 cluster 关联的 user
  - **namespace: <STRING>** #可省，指定该配置环境默认操作的 namespace，省略表示默认为 default 名称空间。
- **name: <STRING>** #指明该上下文的名称。默认格式为"用户名@集群名"，表示这个用户被授权到这个集群中。其实可以使用任意字符串。

## current-context: <STRING> # 当前所使用的上下文

kubectl config current-context 命令就是获取该字段的值。

## users: <\[]Object> # 定义用户信息

- **name: <STRING>** # 指定用户名称
- **user: \<map\[STRING]STRING>** # 用户信息
  - **client-certificate-data: REDACTED** #证书，一般使用集群 ca 证书的 base64 格式的字符串。指明这个用户用来与集群 api 通信时所用客户端的证书
  - **client-key-data: REDACTED** #密钥，一般使用集群 ca 证书的密钥的 base64 格式的字符串。指明这个用户用来与集群 api 通信时所用客户端的密钥

## preferences: <\[]Object>
