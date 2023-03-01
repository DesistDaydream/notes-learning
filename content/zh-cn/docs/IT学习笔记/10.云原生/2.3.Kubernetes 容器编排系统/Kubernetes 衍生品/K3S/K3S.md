---
title: K3S
---

# 概述

> 参考：
> - [GitHub 项目，k3s-io/k3s](https://github.com/k3s-io/k3s)
> - [官方文档](https://rancher.com/docs/k3s/latest/en/)
>   - [中文官方文档](https://docs.rancher.cn/k3s/)
> - [公众号-云原生实验室，K3S 工具进阶完全指南](https://mp.weixin.qq.com/s/ARhxWGypG0wepMqwTLH0mQ)

K3S 是一个轻量的 Kubernetes，具有基本的 kubernetes 功能，将 kubernetes 的主要组件都集成在一个二进制文件中(apiserver、kubelet 等)，这个二进制文件只有不到 100m。内嵌 Containerd，可以通过 Containerd 来启动 coredns 等 kubernetes 的 addone。直接使用 k3s 的二进制文件，即可启动一个 kubernetes 的节点。

Note：K3S 的 kubelet 不支持 systemd 作为 cgroup-driver，原因详见 <https://github.com/rancher/k3s/issues/797>，说是 systemd 的类型无法放进二进制文件里。

k3s 二进制文件包含 kubelet、api-server、kube-controller-manager、kube-scheduler，然后会通过 containerd 拉起 coredns 与 flannel。

# k3s 关联文件与配置

**/etc/rancher/k3s/\*** #

- **./k3s.yaml** #
- **./registries.yaml** #

**/var/lib/rancher/k3s/\*** # k3s 运行时数据存储保存路径

- **./agent/\*** # 作为 k8s 的 node 节点所需要的信息保存路径
  - 包括证书、containerd 数据目录、cni，containerd 的配置文件 等等都在此处
- **./data/\*** #
- **./server/\*** # 作为 k8s 的 master 节点所需要的信息保存路径
  - 包括证书、kube-system 名称空间中的 manifests 文件、etcd 数据 等等都在此处
  - **./db/\*** # 内嵌 Etcd 数据保存路径
  - **./manifests/\*** # k3s 集群启动后，kube-system 名称空间中 pod 的 manifests 文件
  - **./tls/\*** # Kubernetes 主要组件运行所需证书保存路径

**/run/k3s** #
**/run/flannel** #
**/var/lib/kubelet** #
