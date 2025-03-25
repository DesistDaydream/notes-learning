---
title: Kubernetes
linkTitle: Kubernetes
weight: 1
---

# 概述

> 参考：
>
> - [官方文档，概念 - 概述](https://kubernetes.io/docs/concepts/overview/)
> - [play with kubernetes](https://labs.play-with-k8s.com/)

Kubernetes 是一套编排系统，编排目标是实现了 [Containerization](/docs/10.云原生/Containerization/Containerization.md)(容器化) 的容器。

Borg 是谷歌内部的容器管理系统，Kuberntes 根据 Borg 的思路使用 [Go](/docs/2.编程/高级编程语言/Go/Go.md) 语言重新开发，2015 年 7 月份发布

特性：

1. 自我修复：一个 pod 崩了，可以在 1 秒启动，pod 比较轻量，kill 掉崩的容器再启动一个，所以一般情况一个 deployment 会启动多个 pod
2. 自动实现水平扩展：一个 pod 不够，再起一个
3. 自动服务发现和自动负载均衡：当在 k8s 上运行很多程序的时候，通过服务发现，找到所依赖的服务，且多个相同 pod 可以实现自动负载均衡
4. 自动发布与回滚
5. 支持密钥和配置管理：云原声应用，基于环境变量进行配置，需要一个外部组件，当镜像启动为容器的时候，可以自动去外部组件加载相关配置，这个配置中心就是 etcd
6. 存储编排
7. 任务的批量处理执行

Google 成立 CNCF，让各大公司共同管理，并把 Kubernetes 贡献给 [CNCF](/docs/10.云原生/云原生/CNCF.md)，所以 Kubernetes 不会闭源。

# Kubernetes 架构

> 参考：
>
> [官方文档，概念 - 集群架构](https://kubernetes.io/docs/concepts/architecture/)

Kubernetes 集群由代表 Control Palne(控制平面) 和 一组 Nodes(节点) 的机器组成。

![k8s-arch.excalidraw|1000](Excalidraw/k8s-arch.excalidraw.md)

## Control Plane Components(控制平面组件)

### API Server

实现程序：kube-apiserver

[API Server](/docs/10.云原生/Kubernetes/API%20Resource%20与%20Object/API%20Server/API%20Server.md)

### Etcd

[Etcd](/docs/10.云原生/Kubernetes/Etcd/Etcd.md)

### Scheduler

实现程序：kube-scheduler

[Scheduling](/docs/10.云原生/Kubernetes/Scheduling/Scheduling.md)

### Controller Manager

实现程序：kube-controller-manager

[Controller](/docs/10.云原生/Kubernetes/Controller/Controller.md)

## Node Components(节点组件)

### Kubelet

[Kubelet](/docs/10.云原生/Kubernetes/Kubelet/Kubelet.md)

### kube-proxy

[kube-proxy](/docs/10.云原生/Kubernetes/Kubernetes%20网络/kube-proxy/kube-proxy.md)

### Container runtime

[Runtime](/docs/10.云原生/Kubernetes/Kubelet/Runtime.md)

## Addons(附加组件)

### DNS

DNS，core

### WebUI

Dashboard 提供 web 界面的

### Container Resource Monitoring(容器资源监控)

- heapster：是 Kubernetes 原生的集群监控方案。Heapster 以 Pod 的形式运行，它会自动发现集群节点、从节点上的 Kubelet 获取监控数据。Kubelet 则是从节点上的 cAdvisor 收集数据。
  - Heapster 将数据按照 Pod 进行分组，将它们存储到预先配置的 backend 并进行可视化展示。Heapster 当前支持的 backend 有 InfluxDB（通过 Grafana 展示），Google Cloud Monitoring 等。
- ingress

### Cluster-level Logging(集群级日志)

# Kuberntes API 接口

官方文档: https://kubernetes.io/docs/concepts/overview/kubernetes-api/

Kubernetes API 使您可以查询和操纵 Kubernetes 中对象的状态。 Kubernetes 控制平面的核心是 API 服务器和它公开的 HTTP API。用户，集群的不同部分以及外部组件都通过 API 服务器相互通信。

# Kubernetes Objects(对象)

官方文档: https://kubernetes.io/docs/concepts/overview/working-with-objects/

Kubernetes [Object](/docs/10.云原生/Kubernetes/API%20Resource%20与%20Object/Object.md) 是 Kubernetes 系统中的持久实体。 Kubernetes 使用这些实体来表示您的集群状态。了解 Kubernetes 对象模型以及如何使用这些对象。

## kubernetes 所有用 kubectl creat 出来的都可以理解为是一种对象

- workload：Pod，ReplicaSet，Deployment，StatefuSet()，DaemonSet，Job
- 服务发现及均衡：Service，Ingress
- 配置与存储：Volume
  - ConfiMap，secret
  - DownwardAPI
- 集群级对象：Namesapces,Node,Role,ClusterRole,RoleBinding,ClusterRoleBinding
- 元数据型对象：PodTemplate，LimitRange

每个对象所引用的路径格式为：/api/GROUP/VERSION/namespaces/NAMESPACES/TYPE/NAME

可以使用命令 kubectl api-resources 命令查看所有可以创建为对象的资源

# 基本概念

**Cluster：所有运行 kubernetes 的设备的合计**

Cluster 是计算、存储和网络资源的集合，Kubernetes 利用这些资源运行各种基于容器的应用。

**Master ：控制 kubernetes 的 cluster**

Master 是 Cluster 的大脑，它的主要职责是调度，即决定将应用放在哪里运行。Master 运行 Linux 操作系统，可以是物理机或者虚拟机。为了实现高可用，可以运行多个 Master。

**Node ：运行 kuberntes 的 node**

Node 的职责是运行容器应用。Node 由 Master 管理，Node 负责监控并汇报容器的状态，并根据 Master 的要求管理容器的生命周期。Node 运行在 Linux 操作系统，可以是物理机或者是虚拟机。

**Pod：Kubernetes 的最小工作单元**

Pod 是 Kubernetes 的最小工作单元。每个 Pod 包含一个或多个容器。Pod 中的容器会作为一个整体被 Master 调度到一个 Node 上运行。
Kubernetes 引入 Pod 主要基于下面两个目的：

- 可管理性。
  - 有些容器天生就是需要紧密联系，一起工作。Pod 提供了比容器更高层次的抽象，将它们封装到一个部署单元中。Kubernetes 以 Pod 为最小单位进行调度、扩展、共享资源、管理生命周期。
- 通信和资源共享。
  - Pod 中的所有容器使用同一个网络 namespace，即相同的 IP 地址和 Port 空间。它们可以直接用 localhost 通信。同样的，这些容器可以共享存储，当 Kubernetes 挂载 volume 到 Pod，本质上是将 volume 挂载到 Pod 中的每一个容器。user,mnt,pnt。

Pods 有两种使用方式：

- 运行单一容器。
  - one-container-per-Pod 是 Kubernetes 最常见的模型，这种情况下，只是将单个容器简单封装成 Pod。即便是只有一个容器，Kubernetes 管理的也是 Pod 而不是直接管理容器。
- 运行多个容器。
  - 这些容器联系必须非常紧密，而且需要直接共享资源的应该放到一个 Pod 中(注意：当使用多容器的时候，其中一个容器要加上 command 的参数，否则其中一个起不来来)
    - 比如：File Puller 会定期从外部的 Content Manager 中拉取最新的文件，将其存放在共享的 volume 中。Web Server 从 volume 读取文件，响应 Consumer 的请求。这两个容器是紧密协作的，它们一起为 Consumer 提供最新的数据；同时它们也通过 volume 共享数据。所以放到一个 Pod 是合适的。

**Scheduler（kube-scheduler）：调度 POD**

- Scheduler 负责决定将 Pod 放在哪个 Node 上运行。Scheduler 在调度时会充分考虑 Cluster 的拓扑结构，当前各个节点的负载，以及应用对高可用、性能、数据亲和性的需求。

**Controller：执行运行 POD 的任务**

控制器，Kubernetes 一般情况人们不会直接创建 Pod，而是通过创建 Controller 来管理 Pod 的。Controller 中定义了 Pod 的部署特性，比如有几个副本，在什么样的 Node 上运行等。为了满足不同的业务场景，Kubernetes 提供了多种 Controller，包括 Deployment、ReplicaSet、DaemonSet、StatefuleSet、Job 等，我们逐一讨论。一般创建 POD，都是直接创建 Deployment 的 kind，然后定义该 Deployment 下有几个 pod 的副本，一般情况至少有俩，保证 pod 的高可用。注意：deployment 下创建的多个 pod 的功能和内容是一模一样的，多个 pod 被分配到多个节点，以便实现负载均衡和高可用，pod 比较轻量，就算挂了一个，还可以自动销毁后再自动启动一个，所以，不要把一个 deployment 下的多个 pod 分开理解，他们是一个整体

**label selector：标签选择器，简称 selector**

可以给 kubernetes 中所有 node，resource 等等打上标签，然后让某个资源使用 selector 来选择具有相同标签的 Node 或 resource 成为同一组来协调工作或者进行各种限定

比如具有相同标签的 Pod 和 Node，该 Pod 会使用 selector 选择在该 Node 上运行，该 Pod 对该 Node 具有倾向性；或者把具有相同标签的 Service 和 Pod 关联起来，使 Service 使用 selector 知道可以选择哪些 Pod 来进行调度

**Service：服务发现，执行访问 POD 的任务**

- Deployment 可以部署多个副本，每个 Pod 都有自己的 IP，外界如何访问这些副本呢？通过 Pod 的 IP 吗？
- 要知道 Pod 很可能会被频繁地销毁和重启，它们的 IP 会发生变化，用 IP 来访问不太现实。答案是 Service。Service 作为访问 Pod 的接入层来使用
- Kubernetes Service 定义了外界访问一组特定 Pod 的方式。Service 有自己的 IP 和端口，Service 为 Pod 提供了负载均衡。
- 可以把 service 想象成负载均衡功能的前端，该 Service 下的 pod 是负载均衡功能的后端,通过类似 nat 的方式，访问 service 的 IP:PORT，然后转发数据到后端的 pod，注意：在转发到后端 Pod 之前，Service 会先把请求转发到 Endpoints 后再转发到 Pod

**kube-proxy：转发 Service 的流量到 POD**

- service 在逻辑上代表了后端的多个 Pod，外界通过 service 访问 Pod。service 接收到的请求是如何转发到 Pod 的呢？这就是 kube-proxy 要完成的工作。接管系统的 iptables，所有到达 Service 的请求，都会根据 proxy 所定义的 iptables 的规则，进行 nat 转发
- 每个 Node 都会运行 kube-proxy 服务，它负责将访问 service 的 TCP/UPD 数据流转发到后端的容器。如果有多个副本，kube-proxy 会实现负载均衡。
- 每个 Service 的变动(创建，改动，摧毁)都会通知 proxy，在 proxy 所在的本节点创建响应的 iptables 规则，如果 Service 后端的 Pod 摧毁后重新建立了，那么就是靠 proxy 来把 pod 信息提供给 Service。

**Kubernetes 的网络**

kubernetes 的整体网络分为以下三类

- Node IP，各节点网络
- Cluster IP，Service 网络，虚拟的，是主机上 iptables 规则中的地址
- Pod IP，Pod 网络
  - 同一个 Pod 内的多个容器间通信，通过各容器的 lo 通信
  - 各 Pod 之间的通信
    - overlay 叠加网络转发二层报文，通过隧道方式转发三层报文
  - Pod 与 Service 之间的通信，

通过 CNI(Container Network Interface 容器网络接口)来使用第三方 plugin 实现网络的解决方案

- flannel，叠加网络，不支持网络策略
- calico，三层隧道网络，可基于 BGP 协议，即支持网络配置也支持网络策略
- canel，

**Namespace：隔离资源**

官方文档: https://kubernetes.io/docs/concepts/overview/working-with-objects/namespaces/

该 Namespace 与平时所接触的 Namespace 不一样，这是 kubernetes 专用的另一种。如果有多个用户或项目组使用同一个 Kubernetes Cluster，如何将他们创建的 Controller、Pod 等资源分开呢？

答案就是 Namespace。

Namespace 可以将一个物理的 Cluster 逻辑上划分成多个虚拟 Cluster，每个 Cluster 就是一个 Namespace。不同 Namespace 里的资源是完全隔离的。

Kubernetes 默认创建了两个 Namespace。

- default -- 创建资源时如果不指定，将被放到这个 Namespace 中。
- kube-system -- Kubernetes 自己创建的系统资源将放到这个 Namespace 中

**API Server（kube-apiserver）**

API Server 提供 HTTP/HTTPS RESTful API，即 Kubernetes API。API Server 是 Kubernetes Cluster 的前端接口，各种客户端工具（CLI 或 UI）以及 Kubernetes 其他组件可以通过它管理 Cluster 的各种资源。kubectl 就是 API Server 的客户端程序，实现 k8s 各种资源的增删改查

**ETCD**

作为 kubernetes 集群的存储系统使用，保存了集群的所有配置信息，需要高可用，如果需要在生产环境下使用，则需要在单独部署

**Volume 卷**

Volume 的工作流程：可以把 volume 想象成一个中间人，数据流走向：Container—Volum—StorageResource

Volume 的应用场景

在 container 中的磁盘文件是短暂的，这对于 fornon-trivial 类型的 APP 来说会有一些问题。第一，当 container 崩溃时，kubelet 会重启它，但是文件都将丢失并且 container 以最干净的状态启动；第二，当在 Pod 中运行多个 container 的时候，这些 container 需要共享文件以实现功能。Volume 就是为了解决上面两种情况出现的。

volume 定义了一个逻辑卷，该逻辑卷有多种类型，不同的类型可以把不同的存储资源当成 volume 使用(比如内存，文件，分区，网络存储等等)。当我们给 Pod 指定一个 volume 类型后，还需要给该类型的 volume 指定一个可以存放数据的地方；这样，在 container 使用 volume 的时候，可以把自己的数据存放在 volume 所指定的存储资源的地方

**认证**

etcd 内部，etcd 与 apiservice，apiservice-客户端，apiservice 与 kubectl，apiservice 与 kube-proxy

客户端与服务端的概念

谁向谁发请求，前者就是客户端，所在在这里，客户端与服务端没有绝对，一个服务既可以是客户端也可以是服务端

## Pod 被创建的简单流程

![800](https://notes-learning.oss-cn-beijing.aliyuncs.com/te78l0/1616120984034-51654ec9-735a-4eb1-b033-c4dd648cd2d7.png)

1. kubectl 发送部署 Deployment 资源的请求到 API Server。
2. API Server 通知 Controller Manager 创建一个 Deployment 资源。
3. API Server 通知 Scheduler 执行调度任务，将两个副本 Pod 分发到 k8s-node1 和 k8s-node2。
4. API Server 通知 k8s-node1 和 k8s-node2 上的 kubelet 在各自的节点上创建并运行 Pod。

补充几点：

- 更详细的流程详见 [Pod 是如何出现的](/docs/10.云原生/Kubernetes/Kubernetes%20机制与特性/Pod%20是如何出现的.md)
- API Server 的通知是利用 [Watch and Informer](/docs/10.云原生/Kubernetes/Kubernetes%20机制与特性/Watch%20and%20Informer.md) 机制实现的。各个组件主动与 API Server 建立长连接，接收 API Server 的通知
- 应用的配置和当前状态信息保存在 Etcd 中，每一步操作完的结果都会经由 API Server 将信息更新到 Etcd 中。最后执行 `kubectl get pod` 时 API Server 会从 etcd 中读取这些数据。
- flannel 会为每个 Pod 都分配 IP。因为没有创建 service，目前 kube-proxy 还没参与进来。

# 待整理文章

## Node 节点

## 控制平面到 Node 的通信

本文列举控制面节点（确切说是 API 服务器）和 Kubernetes 集群之间的通信路径。 目的是为了让用户能够自定义他们的安装，以实现对网络配置的加固，使得集群能够在不可信的网络上 （或者在一个云服务商完全公开的 IP 上）运行。

节点到控制面

Kubernetes 采用的是中心辐射型（Hub-and-Spoke）API 模式。 所有从集群（或所运行的 Pods）发出的 API 调用都终止于 apiserver（其它控制面组件都没有被设计为可暴露远程服务）。 apiserver 被配置为在一个安全的 HTTPS 端口（443）上监听远程连接请求， 并启用一种或多种形式的客户端身份认证机制。 一种或多种客户端鉴权机制应该被启用， 特别是在允许使用匿名请求 或服务账号令牌的时候。

应该使用集群的公共根证书开通节点，这样它们就能够基于有效的客户端凭据安全地连接 apiserver。 例如：在一个默认的 GCE 部署中，客户端凭据以客户端证书的形式提供给 kubelet。 请查看 kubelet TLS 启动引导 以了解如何自动提供 kubelet 客户端证书。

想要连接到 apiserver 的 Pod 可以使用服务账号安全地进行连接。 当 Pod 被实例化时，Kubernetes 自动把公共根证书和一个有效的持有者令牌注入到 Pod 里。 kubernetes 服务（位于所有名字空间中）配置了一个虚拟 IP 地址，用于（通过 kube-proxy）转发 请求到 apiserver 的 HTTPS 末端。

控制面组件也通过安全端口与集群的 apiserver 通信。

这样，从集群节点和节点上运行的 Pod 到控制面的连接的缺省操作模式即是安全的，能够在不可信的网络或公网上运行。

控制面到节点

从控制面（apiserver）到节点有两种主要的通信路径。 第一种是从 apiserver 到集群中每个节点上运行的 kubelet 进程。 第二种是从 apiserver 通过它的代理功能连接到任何节点、Pod 或者服务。

API 服务器到 kubelet

从 apiserver 到 kubelet 的连接用于：

- 获取 Pod 日志
- 挂接（通过 kubectl）到运行中的 Pod
- 提供 kubelet 的端口转发功能。

这些连接终止于 kubelet 的 HTTPS 末端。 默认情况下，apiserver 不检查 kubelet 的服务证书。这使得此类连接容易受到中间人攻击， 在非受信网络或公开网络上运行也是 不安全的。

为了对这个连接进行认证，使用 --kubelet-certificate-authority 标志给 apiserver 提供一个根证书包，用于 kubelet 的服务证书。

如果无法实现这点，又要求避免在非受信网络或公共网络上进行连接，可在 apiserver 和 kubelet 之间使用 SSH 隧道。

最后，应该启用 Kubelet 用户认证和/或鉴权 来保护 kubelet API。

apiserver 到节点、Pod 和服务

从 apiserver 到节点、Pod 或服务的连接默认为纯 HTTP 方式，因此既没有认证，也没有加密。 这些连接可通过给 API URL 中的节点、Pod 或服务名称添加前缀 https: 来运行在安全的 HTTPS 连接上。 不过这些连接既不会验证 HTTPS 末端提供的证书，也不会提供客户端证书。 因此，虽然连接是加密的，仍无法提供任何完整性保证。 这些连接 目前还不能安全地 在非受信网络或公共网络上运行。

SSH 隧道

Kubernetes 支持使用 SSH 隧道来保护从控制面到节点的通信路径。在这种配置下，apiserver 建立一个到集群中各节点的 SSH 隧道（连接到在 22 端口监听的 SSH 服务） 并通过这个隧道传输所有到 kubelet、节点、Pod 或服务的请求。 这一隧道保证通信不会被暴露到集群节点所运行的网络之外。

SSH 隧道目前已被废弃。除非你了解个中细节，否则不应使用。 Konnectivity 服务是对此通信通道的替代品。

Konnectivity 服务

FEATURE STATE: Kubernetes v1.18 \[beta]

作为 SSH 隧道的替代方案，Konnectivity 服务提供 TCP 层的代理，以便支持从控制面到集群的通信。 Konnectivity 服务包含两个部分：Konnectivity 服务器和 Konnectivity 代理，分别运行在 控制面网络和节点网络中。Konnectivity 代理建立并维持到 Konnectivity 服务器的网络连接。 启用 Konnectivity 服务之后，所有控制面到节点的通信都通过这些连接传输。

请浏览 Konnectivity 服务任务 在你的集群中配置 Konnectivity 服务。

## Controller 控制器

官方文档：<https://kubernetes.io/docs/concepts/architecture/controller/>

**控制器模式是 Kubernetes 的重要设计原则之一**

在机器人技术和自动化领域，控制回路（Control Loop）是一个非终止回路，用于调节系统状态。

这是一个控制环的例子：房间里的温度自动调节器。

当你设置了温度，告诉了温度自动调节器你的期望状态（Desired State）。 房间的实际温度是当前状态（Current State）。 通过对设备的开关控制，温度自动调节器让其当前状态接近期望状态。

控制器通过 apiserver 监控集群的公共状态，并致力于将当前状态转变为期望的状态。

**控制器模式**

一个控制器至少追踪一种类型的 Kubernetes 资源。这些 对象 有一个代表期望状态的 spec 字段。 该资源的控制器负责确保其当前状态接近期望状态。

控制器可能会自行执行操作；在 Kubernetes 中更常见的是一个控制器会发送信息给 API 服务器，这会有副作用。 具体可参看后文的例子。

通过 API 服务器来控制

Job 控制器是一个 Kubernetes 内置控制器的例子。 内置控制器通过和集群 API 服务器交互来管理状态。

Job 是一种 Kubernetes 资源，它运行一个或者多个 Pod， 来执行一个任务然后停止。 （一旦被调度了，对 kubelet 来说 Pod 对象就会变成了期望状态的一部分）。

在集群中，当 Job 控制器拿到新任务时，它会保证一组 Node 节点上的 kubelet 可以运行正确数量的 Pod 来完成工作。 Job 控制器不会自己运行任何的 Pod 或者容器。Job 控制器是通知 API 服务器来创建或者移除 Pod。 控制面中的其它组件 根据新的消息作出反应（调度并运行新 Pod）并且最终完成工作。

创建新 Job 后，所期望的状态就是完成这个 Job。Job 控制器会让 Job 的当前状态不断接近期望状态：创建为 Job 要完成工作所需要的 Pod，使 Job 的状态接近完成。

控制器也会更新配置对象。例如：一旦 Job 的工作完成了，Job 控制器会更新 Job 对象的状态为 Finished。

（这有点像温度自动调节器关闭了一个灯，以此来告诉你房间的温度现在到你设定的值了）。

**直接控制**

相比 Job 控制器，有些控制器需要对集群外的一些东西进行修改。

例如，如果你使用一个控制环来保证集群中有足够的节点，那么控制就需要当前集群外的一些服务在需要时创建新节点。

和外部状态交互的控制器从 API 服务器获取到它想要的状态，然后直接和外部系统进行通信并使当前状态更接近期望状态。

（实际上有一个控制器可以水平地扩展集群中的节点。请参阅 集群自动扩缩容）。

**期望状态与当前状态**

Kubernetes 采用了系统的云原生视图，并且可以处理持续的变化。

在任务执行时，集群随时都可能被修改，并且控制回路会自动修复故障。这意味着很可能集群永远不会达到稳定状态。

只要集群中控制器的在运行并且进行有效的修改，整体状态的稳定与否是无关紧要的。

**设计**

作为设计原则之一，Kubernetes 使用了很多控制器，每个控制器管理集群状态的一个特定方面。 最常见的一个特定的控制器使用一种类型的资源作为它的期望状态， 控制器管理控制另外一种类型的资源向它的期望状态演化。

使用简单的控制器而不是一组相互连接的单体控制回路是很有用的。 控制器会失败，所以 Kubernetes 的设计正是考虑到了这一点。

说明：

可以有多个控制器来创建或者更新相同类型的对象。 在后台，Kubernetes 控制器确保它们只关心与其控制资源相关联的资源。

例如，你可以创建 Deployment 和 Job；它们都可以创建 Pod。 Job 控制器不会删除 Deployment 所创建的 Pod，因为有信息 （标签）让控制器可以区分这些 Pod。

**运行控制器的方式**

Kubernetes 内置一组控制器，运行在 kube-controller-manager 内。 这些内置的控制器提供了重要的核心功能。

Deployment 控制器和 Job 控制器是 Kubernetes 内置控制器的典型例子。 Kubernetes 允许你运行一个稳定的控制平面，这样即使某些内置控制器失败了， 控制平面的其他部分会接替它们的工作。

你会遇到某些控制器运行在控制面之外，用以扩展 Kubernetes。 或者，如果你愿意，你也可以自己编写新控制器。 你可以以一组 Pod 来运行你的控制器，或者运行在 Kubernetes 之外。 最合适的方案取决于控制器所要执行的功能是什么
