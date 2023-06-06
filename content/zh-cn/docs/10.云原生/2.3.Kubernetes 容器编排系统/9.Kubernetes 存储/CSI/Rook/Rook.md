---
title: Rook
weight: 1
---

# 概述

> 参考：
> 
> - 官方文档：<https://rook.github.io/docs/rook/master/>

Rook 是一个开源的 cloud-native storage orchestrator(云原生存储协调器), 提供平台和框架；为各种存储解决方案提供平台、框架和支持，以便与云原生环境本地集成。Rook 作为云原生存储的编排系统，可以直接帮助维护人员管理有状态的存储程序。可以说是一个 SAAS，Storage as a service。

Rook 将存储软件转变为自我管理、自我扩展和自我修复的存储服务，它通过自动化部署、引导、配置、置备、扩展、升级、迁移、灾难恢复、监控和资源管理来实现此目的。

Rook 目前支持 Ceph、NFS、Minio Object Store、CockroachDB 等。这些存储的作用其中之一就是为各个业务提供存储空间，以便统一管理。e.g.使用 Rook 创建一个分布式的 ceph 集群，然后作为 k8s 集群的 storageClass，在某 Pod 需要使用卷的时候，可以直接从 ceph 集群中，拿去存储空间来使用。

Rook 使用底层云本机容器管理、调度和编排平台提供的工具来实现它自身的功能。

Rook 使用 Kubernetes 原语使 Ceph 存储系统能够在 Kubernetes 上运行。下图说明了 Ceph Rook 如何与 Kubernetes 集成：

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/kll50k/1616117730705-c83c3c8d-40a2-4f54-a1bc-eff7e8f6f84a.jpeg)

随着 Rook 在 Kubernetes 集群中运行，Kubernetes 应用程序可以挂载由 Rook 管理的块设备和文件系统，或者可以使用 S3 / Swift API 提供对象存储。Rook oprerator 自动配置存储组件并监控群集，以确保存储处于可用和健康状态。

Rook oprerator 是一个简单的容器，具有引导和监视存储集群所需的全部功能。oprerator 将启动并监控 ceph monitor pods 和 OSDs 的守护进程，它提供基本的 RADOS 存储。oprerator 通过初始化运行服务所需的 pod 和其他组件来管理池，对象存储（S3 / Swift）和文件系统的 CRD。

oprerator 将监视存储后台驻留程序以确保群集正常运行。Ceph mons 将在必要时启动或故障转移，并在群集增长或缩小时进行其他调整。oprerator 还将监视 api 服务请求的所需状态更改并应用更改。

Rook oprerator 还创建了 Rook agent。这些 agent 是在每个 Kubernetes 节点上部署的 pod。每个 agent 都配置一个 Flexvolume 插件，该插件与 Kubernetes 的 volume controller 集成在一起。处理节点上所需的所有存储操作，例如附加网络存储设备，安装卷和格式化文件系统。

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/kll50k/1616117730704-e5377336-af63-4167-a704-e1cd1eeaa124.jpeg)

该 rook 容器包括所有必需的 Ceph 守护进程和工具来管理和存储所有数据 - 数据路径没有变化。 rook 并没有试图与 Ceph 保持完全的忠诚度。 许多 Ceph 概念（如 placement groups 和 crush maps）都是隐藏的，因此您无需担心它们。 相反，Rook 为管理员创建了一个简化的用户体验，包括物理资源，池，卷，文件系统和 buckets。 同时，可以在需要时使用 Ceph 工具应用高级配置。

Rook 在 golang 中实现。Ceph 在 C ++中实现，其中数据路径被高度优化。我们相信这种组合可以提供两全其美的效果。
