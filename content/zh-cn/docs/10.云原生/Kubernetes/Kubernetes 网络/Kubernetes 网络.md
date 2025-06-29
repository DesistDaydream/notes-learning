---
title: Kubernetes 网络
linkTitle: Kubernetes 网络
weight: 1
---

# 概述

> 参考：
>
> - [官方文档，概念-集群管理-集群网络](https://kubernetes.io/docs/concepts/cluster-administration/networking/)

Kubernetes 的整体网络分为以下三类

- Node IP(各节点网络) #
- Cluster IP(Service 网络) # 虚拟的，在 Netfilter 结构上，就是主机上 iptables 规则中的地址
- Pod IP(Pod 网络) #

网络是 Kubernetes 的核心部分，Kubernetes 中有下面几个点需要互相通信

- 同一个 Pod 内的多个容器间通信，通过各容器的 lo 通信
- Pod 之间的通信，Pod IP<-->Pod IP
  - overlay 叠加网络转发二层报文，通过隧道方式转发三层报文
- Pod 与 Service 之间的通信，Pod IP<-->Cluster IP。详见 [Service(服务)](/docs/10.云原生/Kubernetes/Kubernetes%20网络/Service(服务).md)。
- Service 与集群外部客户端的通信。详见 [Service(服务)](/docs/10.云原生/Kubernetes/Kubernetes%20网络/Service(服务).md)。

Kubernetes 的宗旨就是在应用之间共享机器。 通常来说，共享机器需要两个应用之间不能使用相同的端口，但是在多个应用开发者之间 去大规模地协调端口是件很困难的事情，尤其是还要让用户暴露在他们控制范围之外的集群级别的问题上。

动态分配端口也会给系统带来很多复杂度 - 每个应用都需要设置一个端口的参数， 而 API 服务器还需要知道如何将动态端口数值插入到配置模块中，服务也需要知道如何找到对方等等。 与其去解决这些问题，Kubernetes 选择了其他不同的方法。

## Kubernetes 网络模型

每一个 Pod 都有它自己的 IP 地址，这就意味着你不需要显式地在每个 Pod 之间创建链接， 你几乎不需要处理容器端口到主机端口之间的映射。 这将创建一个干净的、向后兼容的模型，在这个模型里，从端口分配、命名、服务发现、 负载均衡、应用配置和迁移的角度来看，Pod 可以被视作虚拟机或者物理主机。

Kubernetes 对所有网络设施的实施，都需要满足以下的基本要求（除非有设置一些特定的网络分段策略）：

- 节点上的 Pod 可以不通过 NAT 和其他任何节点上的 Pod 通信
- 节点上的代理（比如：系统守护进程、kubelet） 可以和节点上的所有 Pod 通信

备注：仅针对那些支持 Pods 在主机网络中运行的平台(比如：Linux) ：

- 那些运行在节点的主机网络里的 Pod 可以不通过 NAT 和所有节点上的 Pod 通信

这个模型不仅不复杂，而且还和 Kubernetes 的实现廉价的从虚拟机向容器迁移的初衷相兼容， 如果你的工作开始是在虚拟机中运行的，你的虚拟机有一个 IP ， 这样就可以和其他的虚拟机进行通信，这是基本相同的模型。

Kubernetes 的 IP 地址存在于 Pod 范围内 - 容器分享它们的网络命名空间 - 包括它们的 IP 地址。 这就意味着 Pod 内的容器都可以通过 localhost 到达各个端口。 这也意味着 Pod 内的容器都需要相互协调端口的使用，但是这和虚拟机中的进程似乎没有什么不同， 这也被称为“一个 Pod 一个 IP” 模型。

如何实现这一点是正在使用的容器运行时的特定信息。

也可以在 node 本身通过端口去请求你的 Pod （称之为主机端口）， 但这是一个很特殊的操作。转发方式如何实现也是容器运行时的细节。 Pod 自己并不知道这些主机端口是否存在。

# Network Plugin(网络插件) — 实现 Kubernetes 网络模型的方式

官方文档：<https://kubernetes.io/docs/concepts/extend-kubernetes/compute-storage-net/network-plugins/>

Kubernetes 中的网络插件有几种类型：

- CNI 插件： 遵守 appc/CNI 规约，为互操作性设计。详见：[CNI](/docs/10.云原生/Kubernetes/Kubernetes%20网络/CNI/CNI.md)
- Kubenet 插件：使用 bridge 和 host-local CNI 插件实现了基本的 cbr0
