---
title: K3S
weight: 1
---

# 概述

> 参考：
>
> - [GitHub 项目，k3s-io/k3s](https://github.com/k3s-io/k3s)
> - [官方文档](https://rancher.com/docs/k3s/latest/en/)
>   - [中文官方文档](https://docs.rancher.cn/k3s/)
> - [公众号-云原生实验室，K3S 工具进阶完全指南](https://mp.weixin.qq.com/s/ARhxWGypG0wepMqwTLH0mQ)

K3S 是一个轻量的 Kubernetes，具有基本的 kubernetes 功能，将 kubernetes 的主要组件都集成在一个二进制文件中(apiserver、kubelet 等)，这个二进制文件只有不到 100m。内嵌 Containerd，可以通过 Containerd 来启动 coredns 等 kubernetes 的 addone。直接使用 k3s 的二进制文件，即可启动一个 kubernetes 的节点。

Note: K3S 的 kubelet 不支持 systemd 作为 cgroup-driver，原因详见 https://github.com/rancher/k3s/issues/797 ，说是 systemd 的类型无法放进二进制文件里。

k3s 二进制文件包含 kubelet、api-server、kube-controller-manager、kube-scheduler，然后会通过 containerd 拉起 coredns 与 flannel。

## K3S 封装的组件

> [官方文档，安装 - 管理封装的组件](https://docs.k3s.io/zh/installation/packaged-components)

K3S 封装了部分非 K8S 核心组件，比如 `coredns`、`traefik`、`local-storage`、`metrics-server`、`servicelb`。这些组件通常都会以 manifests 文件的方式保存在 `/var/lib/rancher/k3s/server/manifests/` 目录中，当 K3S 启动时，自动拉起这些组件。

> 嵌入式 `servicelb` LoadBalancer controller 没有 manifests 文件，但由于历史原因，它可以像 `AddOn` 一样被禁用

# k3s 关联文件与配置

**/etc/rancher/k3s/** # 各种配置存放目录

- **./k3s.yaml** # kubeconfig 文件
- **./config.yaml** # 运行时配置文件，与环境变量和命令行参数等效。
- **./registries.yaml** # 镜像仓库配置，可以配置加速、私有仓库。配置方式详见[官方文档-安装-私有镜像仓库配置](https://docs.k3s.io/zh/installation/private-registry)

**/run/k3s/** # K3S 所使用的容器 Runtime 的数据保存路径。

- **./containerd/** # 与 [Containerd](/docs/10.云原生/Containerization%20implementation/Containerd/Containerd.md#Containerd%20关联文件与配置 关联文件与配置>) 中的 /run/containerd/ 目录功能一致。

**/run/flannel/** # 与 [Flannel](/docs/10.云原生/Kubernetes/8.Kubernetes%20网络/CNI/Flannel.md#Flannel%20关联文件与配置%20容器编排系统/8.Kubernetes%20网络/CNI/Flannel#Flannel%20关联文件与配置 容器编排系统/8.Kubernetes 网络/CNI/Flannel.md#Flannel 关联文件与配置 容器编排系统/8.Kubernetes 网络/CNI/Flannel#Flannel 关联文件与配置>) 中 /run/flannel/ 目录功能一致。

**/var/lib/rancher/k3s/** # k3s 运行时数据存储保存路径

- **./server/** # 作为 k8s 的 master 节点所需要的信息保存路径
  - 包括证书、kube-system 名称空间中的 manifests 文件、etcd 数据 等等都在此处
  - **./db/** # 内嵌 Etcd 数据保存路径
  - **./manifests/** # 功能与 [Kubelet](/docs/10.云原生/Kubernetes/Kubelet/Kubelet.md) 管理的 `/etc/kubernetes/manifets/` 目录功能一样。k3s 集群启动后，读取该目录下的 manifets 文件以运行 Pod。
  - **./tls/** # Kubernetes 主要组件运行所需证书保存路径
- **./agent/** # 作为 k8s 的 node 节点所需要的信息保存路径。对于 K3S 来说，master 节点也属于 node 节点，所以 master 节点在该目录也会保存数据。
  - 包括证书、containerd 数据目录、cni，containerd 的配置文件 等等都在此处
  - **./etc/** # 各种组件的配置文件保存路径。比如 CNI、Containerd、Flannel 等等，相当于各个组件自己所使用的 etc 目录。
  - **./containerd/** # 与 [Containerd](/docs/10.云原生/Containerization%20implementation/Containerd/Containerd.md#Containerd%20关联文件与配置 关联文件与配置>) 中的 /var/lib/containerd/ 目录功能一致。
- **./data/**

**/var/lib/kubelet/** #

K3S 所有可能使用的目录可以参考 [清理 K3S](/docs/10.云原生/2.3.Kubernetes%20容器编排系统/Kubernetes%20衍生品/K3S/K3S%20部署与清理.md#清理 K3S)

# 进入容器的文件系统

详见 [进入容器文件系统](/docs/10.云原生/Containerization%20implementation/容器管理/容器运行时管理/进入容器文件系统.md)。在 k3s 中，如果是 containerd 的话，则是在 /run/k3s/containerd/ 目录代替 /run/containerd/ 目录

`/run/k3s/containerd/io.containerd.runtime.v2.task/k8s.io/${ContainerID}/rootfs/``

# K3S 与 K8S 的区别

> 参考：
>
> - [公众号-边缘计算k3s社区，K3s vs K8s：轻量级和全功能的对决](https://mp.weixin.qq.com/s/575ZBryg4bv9k01To1QY7w)
