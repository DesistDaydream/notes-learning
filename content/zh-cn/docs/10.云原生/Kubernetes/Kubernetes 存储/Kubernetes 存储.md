---
title: Kubernetes 存储
linkTitle: Kubernetes 存储
weight: 1
---

# 概述

> 参考：
>
> - [官方文档，概念-存储](https://kubernetes.io/docs/concepts/storage/)
> - [公众号-CNCF，卷扩展现在是个稳定特性](https://mp.weixin.qq.com/s/hNR5XkMeZbDVInUOX_5MAg)
> - [公众号-CNCF，存储容量跟踪在 Kubernetes1.24 中正式 GA](https://mp.weixin.qq.com/s/EBghRVRQvnPSTf4YdCkp2w)

在 [Container](/docs/10.云原生/Containerization/Containerization.md) 中的文件在磁盘上是临时存储的(这与 Docker 一样，容器删除后，容器内的文件也随着删除)，这给 Container 中运行的需要持久化存储的应用程序带来了很多问题。

- 第一，当 Container 崩溃时，kubelet 会重启它，但是文件都将丢失并且 Container 以最干净的状态启动
- 第二，当在 Pod 中运行多个 Container 的时候，这些 Container 需要共享文件以实现功能。

**Volume(卷)** 就是为了解决上面两种情况出现的。

从本质上讲，Volume(卷) 只是一个包含一些数据目录，Pod 中 Container 可以访问这个目录。至于该目录是如何形成的是由所使用的 Volume 类型决定的。这个 Volume 的类型可以是：host 的内存，host 的文件，host 的目录，nfs、glusterfs、甚至是云厂商所提供的各种类型的存储

**可以说，Kubernetes 存储功能的基础，就是 Volume(卷)。**

Volume 功能详解见 [Volume](/docs/10.云原生/Kubernetes/Kubernetes%20存储/Volume/Volume.md) 章节

## 与 Docker 中的 Volume 的概念比较

Kubernetse 为什么不直接复用 Docker 中的 Volume，而是要自己实现呢?~

Kubernetes Volume 和 Docker Volume 概念相似，但是又有不同的地方，Kubernetes Volume 与 Pod 的生命周期相同，但与容器的生命周期不相关。当容器终止或重启时，Volume 中的数据也不会丢失。当 Pod 被删除时，Volume 才会被清理。并且数据是否丢失取决于 Volume 的具体类型，比如 emptyDir 类型的 Volume 数据会丢失，而持久化类型的数据则不会丢失。另外 Kubernetes 提供了将近 20 种 Volume 类型。

# Volume 的实现-Volume Plugins(卷插件)

**Volume Plugins(卷插件)** 是实现 Kubernetes 存储功能的方式。说白了，**Volume Plugins 就是用来实现 Volume 功能的**。一共可以分为两类：

- **In-Tree(树内)** # 代码逻辑在 K8S 官方仓库中；表示源码是放在 Kubernetes 内部的(常见的 NFS、cephfs 等)，和 Kubernetes 一起发布、管理与迭代，缺点是迭代速度慢、灵活性差；
- **Out-of-Tree(树外)** # 代码逻辑在 K8s 官方仓库之外，实现与 K8s 代码的解耦；表示代码独立于 Kubernetes，它是由存储提供商实现的，目前主要有 **Flexvolume** 或 **CSI** 两种实现机制，可以根据存储类型实现不同的存储插件

## In-Tree(树内)

**In-Tree 类型的卷插件的代码 与 Kubernetes 代码 在一起**。比如 emptyDir、hostPath、ConfigMap、PVC 等等类型的 Volume，凡是在[官方文档的卷类型](https://kubernetes.io/docs/concepts/storage/volumes/)中的都属于 In-Tree 类型的卷插件。所以 卷插件类型 也可以说是 卷类型。

## Out-of-Tree(树外)

**OUt-of-Tree 类型的卷插件代码 与 Kubernetes 代码 不在一起**。这种类型的插件，可以让存储供应商创建自定义的存储插件而无需将他们添加到 Kubernetes 代码仓库。

这类卷插件的实现方式分两位 **CSI** 和 **FlexVolume** 两类。都允许独立于 Kubernetes 代码库开发卷插件，并作为 Pod 部署在 Kubernetes 集群中。

### CSI - Container Storage Interface(容器存储接口)

CSI 为容器编排系统定义标准规范，以将任意存储系统暴露给它们的容器。

详见 [CSI](/docs/10.云原生/Kubernetes/Kubernetes%20存储/CSI/CSI.md)

### FlexVolume

FlexVolume 是一个自 1.2 版本（在 CSI 之前）以来在 Kubernetes 中一直存在的树外插件接口。 它使用基于 exec 的模型来与驱动程序对接。 用户必须在每个节点（在某些情况下是主控节点）上的预定义卷插件路径中安装 FlexVolume 驱动程序可执行文件。
Pod 通过 `flexvolume` 树内插件与 Flexvolume 驱动程序交互。 更多详情请参考 [FlexVolume](https://github.com/kubernetes/community/blob/master/contributors/devel/sig-storage/flexvolume.md) 示例。

## In-Tree 向 CSI 迁移

从 Kubernetes [2021 年 12 月 10 日的博客](https://kubernetes.io/zh-cn/blog/2021/12/10/storage-in-tree-to-csi-migration-status-update/)中报告了迁移工作的进展。从这里可以看到，Kubernetes 官方希望我们无论使用那哪种 In-Tree 模型的卷插件，都尽早迁移至使用 CSI 驱动的模型。

容器存储接口旨在帮助 Kubernetes 取代其现有的树内存储驱动机制 ── 特别是供应商的特定插件。自 v1.13 起，Kubernetes 对[容器存储接口](https://github.com/container-storage-interface/spec/blob/master/spec.md#README)的支持工作已达到[正式发布阶段](https://kubernetes.io/blog/2019/01/15/container-storage-interface-ga/)。引入对 CSI 驱动的支持，将使得 Kubernetes 和存储后端技术之间的集成工作更易建立和维护。使用 CSI 驱动可以实现更好的可维护性（驱动作者可以决定自己的发布周期和支持生命周期）、减少出现漏洞的机会（得益于更少的树内代码，出现错误的风险会降低。另外，集群操作员可以只选择集群需要的存储驱动）。

随着更多的 CSI 驱动诞生并进入生产就绪阶段，Kubernetes 存储特别兴趣组希望所有 Kubernetes 用户都能从 CSI 模型中受益 ── 然而，我们不应破坏与现有存储 API 类型的 API 兼容性。对此，我们给出的解决方案是 CSI 迁移：该功能实现将树内存储 API 翻译成等效的 CSI API，并把操作委托给一个替换的 CSI 驱动来完成。

CSI 迁移工作使存储后端现有的树内存储插件（如 kubernetes.io/gce-pd 或 kubernetes.io/aws-ebs）能够被相应的 [CSI 驱动](https://kubernetes-csi.github.io/docs/introduction.html) 所取代。如果 CSI 迁移功能正确发挥作用，Kubernetes 终端用户应该不会注意到有什么变化。现有的 StorageClass、PersistentVolume 和 PersistentVolumeClaim 对象应继续工作。当 Kubernetes 集群管理员更新集群以启用 CSI 迁移功能时，利用到 PVCs[1](https://kubernetes.io/zh-cn/blog/2021/12/10/storage-in-tree-to-csi-migration-status-update/#fn:1)（由树内存储插件支持）的现有工作负载将继续像以前一样运作 ── 不过在幕后，Kubernetes 将所有存储管理操作（以前面向树内存储驱动的）交给 CSI 驱动控制。

举个例子。假设你是 kubernetes.io/gce-pd 用户，在启用 CSI 迁移功能后，你仍然可以使用 kubernetes.io/gce-pd 来配置新卷、挂载现有的 GCE-PD 卷或删除现有卷。所有现有的 API/接口 仍将正常工作 ── 只不过，底层功能调用都将通向 [GCE PD CSI 驱动](https://github.com/kubernetes-sigs/gcp-compute-persistent-disk-csi-driver)，而不是 Kubernetes 的树内存储功能。

不过，这里面没提到那些 Kubernetes 内置资源类型的卷，比如 ConfigMap、Secret、等等。人们主要需要迁移的是那些类似 glusterfs、cephfs 之类第三方的 In-Tree 类型的卷插件。

迁移进展文章

- <https://mp.weixin.qq.com/s/6nhv2zQIAOAfUJ661YmDsQ>

# Kubernetes 存储模型

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/aplkpr/1616117503767-42e19ed6-fbd6-4b5b-bc38-db7e7a699432.jpeg)

- **Volume Controller** # K8S 的卷控制器。[代码位置](https://github.com/kubernetes/kubernetes/tree/master/pkg/controller/volume)
  - **PV Controller** # 负责 PV/PVC 的绑定、生命周期管理，并根据需求进行数据卷的 **Provision/Delete** 操作。[代码位置](https://github.com/kubernetes/kubernetes/tree/master/pkg/controller/volume/persistentvolume)
  - **AD Controller** # 负责存储设备的 **Attach/Detach** 操作，将设备挂载到目标节点。[代码位置](https://github.com/kubernetes/kubernetes/tree/master/pkg/controller/volume/attachdetach)
- **Kubelet** # Kubelet 是在每个 Node 节点上运行的主要 “节点代理”，功能是 Pod 生命周期管理、容器健康检查、容器监控等；
  - **Volume Manager** # 属于 kubelet 中的组件，管理卷的 **Mount/Unmount** 操作、卷设备的格式化的操作.
    - 注意：Volume Manager 也可以负责数据卷的 **Attach/Detach** 操作，需要配置 kubelet 开启特性
- **Volume Plugins** # 它主要是对上面所有挂载功能的实现。PV Controller、AD Controller、Volume Manager 主要是进行操作的调用，而具体操作则是由 Volume Plugins 实现的。根据源码的位置可将 Volume Plugins 分为 In-Tree 和 Out-of-Tree 两类
  - In-Tree # 与 Kubernetes 代码强耦合的卷插件
  - Out-of-Tree # 与 Kubernetes 代码无关的卷插件。
- **Scheduler** # 实现对 Pod 的调度能力，会根据一些存储相关的的定义去做存储相关的调度
- 其他
  - **External Provioner：** External Provioner 是一种 sidecar 容器，作用是调用 Volume Plugins 中的 CreateVolume 和 DeleteVolume 函数来执行 **Provision/Delet**e 操作。因为 K8s 的 PV 控制器无法直接调用 Volume Plugins 的相关函数，故由 External Provioner 通过 gRPC 来调用
  - **External Attacher：** External Attacher 是一种 sidecar 容器，作用是调用 Volume Plugins 中的 ControllerPublishVolume 和 ControllerUnpublishVolume 函数来执行 **Attach/Detach** 操作。因为 K8s 的 AD 控制器无法直接调用 Volume Plugins 的相关函数，故由 External Attacher 通过 gRPC 来调用。
