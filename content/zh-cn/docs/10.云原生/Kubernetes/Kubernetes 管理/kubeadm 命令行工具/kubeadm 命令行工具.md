---
title: kubeadm 命令行工具
---

# 概述

> 参考：
>
> - [官方文档，参考 - kubeadm](https://kubernetes.io/docs/reference/setup-tools/kubeadm/)
> - [kubeadm 库](https://pkg.go.dev/k8s.io/kubernetes/cmd/kubeadm/app/apis/kubeadm)
> - [v1beta2 版本的 kubeadm 包的配置文件字段详解](https://pkg.go.dev/k8s.io/kubernetes/cmd/kubeadm/app/apis/kubeadm/v1beta2)

Kubeadm 是一个工具，它提供了 kubeadm init 以及 kubeadm join 这两个命令作为快速创建 [Kubernetes](/docs/10.云原生/Kubernetes/Kubernetes.md) 集群的最佳实践。

kubeadm 通过执行必要的操作来启动和运行一个最小可用的集群。它被故意设计为只关心启动集群，而不是准备节点环境的工作。同样的，诸如安装各种各样的可有可无的插件，例如 Kubernetes 控制面板、监控解决方案以及特定云提供商的插件，这些都不在它负责的范围。

相反，我们期望由一个基于 kubeadm 从更高层设计的更加合适的工具来做这些事情；并且，理想情况下，使用 kubeadm 作为所有部署的基础将会使得创建一个符合期望的集群变得容易。

## kubeadm 中的资源

实际上，kubeadm 继承了 kubernetes 的哲学，一切介资源，只不过由于 kubeadm 并没有控制器逻辑、也并不需要将这些资源实例化为一个个的对象。这些资源主要是为了让 kubeadm 的概念以及使用方式，更贴近 Kubernetes，所以 **kubeadm 的资源仅仅作为定义配置所用**。在 kubeadm 的 [API 代码](https://github.com/kubernetes/kubernetes/blob/master/cmd/kubeadm/app/apis/kubeadm/v1beta2/types.go)中，也可以看到这些资源的结构体定义。

kubeadm 的运行时行为通常由下面几个 API 资源来控制：

1. **InitConfiguration(初始化配置)** #
2. **ClusterConfiguation(集群配置)** #
3. **KubeletConfiguration(kubelet 程序配置)** #
4. **KubeProxyConfiguration(kube-proxy 程序配置)** #
5. **JoinConfiguration(加入集群配置)** #

其中 InitConfiguration、ClusterConfiguation、JoinConfiguration 资源属于 kubeadm 在控制集群时所用的配置

而 KubeletConfiguration 与 KubeProxyConfiguration 资源，实际上就是 kubelet 和 kube-proxy 程序的配置文件，kubeadm 可以通过其自身的配置文件，在控制集群时，修改 kubelet 与 kube-proxy 程序的配置文件。

可以通过 kubeadm config print init-defaults 命令可以输出这些资源的 Manifests 模板，该命令默认会输出 InitConfiguration 与 ClusterConfiguration 的默认配置，可以通过使用 --component-configs STRING 选项来输出 KubeletConfiguration 和 KubeProxyConfiguration 的默认配置

**而 kubeadm 的这些资源的 Manifests，其实就是 kubeadm 在部署集群时所使用的配置文件。**

# kubeadm 安装

> 参考：
>
> - [官方文档，安装 - 生产环境 - 工具 - kubeadm - 安装 kubeadm](https://kubernetes.io/docs/setup/production-environment/tools/kubeadm/install-kubeadm/)

### 安装 run-time(运行时)

Kubernetes 使用 container runtime 以便在 Pods 运行容器。

默认情况下，Kubernetes 使用 **Container Runtime Interface(容器运行时接口，简称 CRI**) 与我们选择的容器运行时交互。

如果未指定运行时(可以通过 kubelet 程序的 `--cri-socket` 标志指定运行时的 Sokcet 路径)，则 kubeadm 会通过扫描众所周知的 Unix 域套接字列表自动尝试检测已安装的容器运行时。下表列出了容器运行时及其关联的套接字路径：

| Runtime    | Unix Socket 的路径                 |
| ---------- | ------------------------------- |
| docker     | /var/run/docker.sock            |
| containerd | /run/containerd/containerd.sock |
| CRI-O      | /var/run/crio/crio.sock         |

如果同时检测到 Docker 和 containerd，则 Docker 优先。这是必需的，因为 Docker 18.09 附带了容器，即使您仅安装了 Docker，也可以检测到两者。如果检测到其他两个或更多运行时，则 kubeadm 退出并显示错误。

> [!Warning]
> kubelet 通过内置的 dockershim CRI 实现与 Docker 集成。2020 年 12 月 2 日[官方博客发文](https://kubernetes.io/blog/2020/12/02/dont-panic-kubernetes-and-docker/)称，在 v1.20，您将收到 Docker 弃用警告。在将来的 Kubernetes 版本（目前计划在 2021 年下半年为 1.22 版本）中删除 Docker 运行时支持时，它将不再受支持，您将需要切换到其他兼容的容器运行时之一，例如 containerd 或 CRI-O 。

# kubeadm 关联文件与配置

kubeadm 的配置文件主要用来**部署集群所用**，其中包括初始化集群所需的所有信息。

**kubeadm-config.yaml** # kubeadm 所需的配置文件，一般使用这个名字。可以通过 --config 参数指定其他的文件。详见 [kubeadm 配置详解](/docs/10.云原生/Kubernetes/Kubernetes%20管理/kubeadm%20命令行工具/kubeadm%20配置详解.md)
