---
title: Containerization implementation
linkTitle: Containerization implementation
date: 2024-02-27T09:36
weight: 1
---

# 概述

> 参考：
>
> -

# OCI Runtime 规范的实现

> 参考：
>
> - [公众号-k8s 技术圈，Containerd 深度剖析-runtime 篇](https://mp.weixin.qq.com/s/NPxLLhRkpNdTgVcKQSLcFA)

当人们想到容器运行时，可能会想到一连串的相关概念；runc、runv、lxc、lmctfy、Docker（containerd）、rkt、cri-o。每一个都是基于不同的场景而实现的，均实现了不同的功能。如 containerd 和 cri-o，实际均可使用 runc 来运行容器，但其实现了如镜像管理、容器 API 等功能，可以将这些看作是比 runc 具备的更高级的功能。

可以发现，容器运行时是相当复杂的。每个运行时都涵盖了从低级到高级的不同部分，如下图所示

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/ctvy4o/1653965809357-01c7d7f1-81d0-49f1-beaa-15bd63e7acd6.png)

根据功能范围划分，将其分为 **Low level Container Runtime(低级容器运行时)** 和 **High level Container Runtime(高级容器运行时)**

- 低级容器运行时 # 只关注容器的本身运行
- 高级容器运行时 # 支持更多高级功能的运行时，如镜像管理及一些 gRPC/Web APIs，通常被称为

需要注意的是，低级运行时和高级运行时有本质区别，各自解决的问题也不同。

## 低级运行时

低级运行时的功能有限，通常执行运行容器的低级任务。大多数开发者日常工作中不会使用到。其一般指按照 OCI 规范、能够接收可运行 roofs 文件系统和配置文件并运行隔离进程的实现。这种运行时只负责将进程运行在相对隔离的资源空间里，不提供存储实现和网络实现。但是其他实现可以在系统中预设好相关资源，低级容器运行时可通过 config.json 声明加载对应资源。低级运行时的特点是底层、轻量，限制也很一目了然：

- 只认识 rootfs 和 config.json，没有其他镜像能力
- 不提供网络实现
- 不提供持久实现
- 无法跨平台等

### RunC

> 参考：
>
> - [GitHub 项目，opencontainers/runc](https://github.com/opencontainers/runc)

runc 是一个 CLI 工具，用于根据 OCI 规范生成和运行容器。

### youki

> 参考：
>
> - [GitHub 项目，containers/youki](https://github.com/containers/youki)

使用 Rust 语言写的，类似于 Runc 的容器运行时，

### Sysbox

> 参考：
>
> - [GitHub 项目，nestybox/sysbox](https://github.com/nestybox/sysbox)

Sysbox 是一个新型的 OCI 容器运行时，对标 runc。相比于 runc，Sysbox 在以下两个方面做了增强：

- 增强容器隔离性：Sysbox 为所有容器开启 user namespace（即容器中的 root 用户映射为主机中的普通用户），在容器中隐藏宿主机的信息，锁定容器的初始挂载，等等。
- 容器不仅可以运行普通进程，还可以运行 systemd、Docker、K8s、K3s 等系统级软件，一定程度上可以替换虚拟机。

最初 Sysbox 只支持 Docker，但最新版本 v0.4.0 已支持直接作为 Kubernetes 的 CRI 运行时。

### Kata Container

> 参考：
>
> - [GitHub 项目，kata-containers/kata-containers](https://github.com/kata-containers/kata-containers)

Kata Containers 是一个开源项目和社区，致力于构建轻量级虚拟机 (vm) 的标准实现，该虚拟机感觉和性能类似于容器，但提供 vm 的工作负载隔离和安全性优势。

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/ctvy4o/1616122531941-6b13921a-78c5-45a1-9a38-6695b517bca8.png)

## 高级运行时

高级运行时负责容器镜像的传输和管理，解压镜像，并传递给低级运行时来运行容器。通常情况下，高级运行时提供一个守护程序和一个 API，远程应用程序可以使用它来运行容器并监控它们，它们位于低层运行时或其他高级运行时之上。

高层运行时也会提供一些看似很低级的功能。例如，管理网络命名空间，并允许容器加入另一个容器的网络命名空间。

这里有一个类似逻辑分层图，可以帮助理解这些组件是如何结合在一起工作的。

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/ctvy4o/1653966607306-f97afdfd-66fd-4d4d-ab04-364b1b60f27e.png)

### Docker

[Docker](/docs/10.云原生/Containerization%20implementation/Docker/Docker.md)

### Containerd

[Containerd](/docs/10.云原生/Containerization%20implementation/Containerd/Containerd.md)
