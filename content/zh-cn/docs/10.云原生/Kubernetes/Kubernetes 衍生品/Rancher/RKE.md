---
title: RKE
linkTitle: RKE
weight: 20
---

# 概述

> 参考：
>
> - 官方文档：
>   - https://rancher.com/docs/rke/latest/en/
>   - https://rancher2.docs.rancher.cn/docs/rke/_index/

**Rancher Kubernetes Engine(RKE)**，是经过 CNCF 认证的 Kubernetes 发行版，完全在 Docker 容器内运行。它适用于裸机和虚拟机。RKE 解决了安装复杂性的问题，这是 Kubernetes 社区中的常见问题。借助 RKE，Kubernetes 的安装和操作既简单又自动化，并且完全独立于所运行的操作系统和平台。只要服务器可以运行受支持的 Docker 版本，就可以使用 RKE 部署和运行 Kubernetes。

使用 rke 工具，仅需通过一个 yaml 的配置文件以及 docker 环境，即可启动一个功能完全的 kubernetes 集群。其中所有系统组件(包括 kubelet)都是以容器的方式运行的。通过 Rancher 创建的 kubernetes 集群，就是 RKE 集群。

## RKE 集群与原生 K8S 集群的区别

RKE 与 sealos 实现高可用的方式类似。不同点是 RKE 集群的 node 节点是通过 ngxin 来连接 API Server。

# RKE 集群部署

参考：RKE 部署与清理

- 下载 rke 二进制文件。(在 github 上下载 rke 命令行工具)
- 创建集群配置文件。
  - RKE 默认使用名为 cluster.yml 的集群配置文件来确定集群中应该包含哪些节点以及如何部署 Kubernetes。
  - 下面是一个单节点 cluster.yml 文件示例，

```bash
cat > cluster.yml <<
# 指定要部署集群的节点信息
nodes:
# 指定该节点的IP
- address: 1.2.3.4
  # 指定部署集群时，所使用的用户
  user: ubuntu
  # 指定该集群的角色，controlplane运行k8s主要组件，etcd运行etcd，worker运行用户创建的非k8s主要组件的pod。
  role:
  - controlplane # 对应 k8s master 节点
  - etcd
  - worker # 对应 k8s node 节点
EOF
```

- 在 cluster.yml 执行 `rke up` 命令

## RKE 清理

参考：<https://rancher.com/docs/rke/latest/en/managing-clusters/>。

中文：<https://rancher2.docs.rancher.cn/docs/rke/managing-clusters/_index>

在 cluster.yaml 文件所在目录

# RKE 配置

RKE 默认通过一个名为 cluster.yml 的文件配置集群参数。可以通过 --config 选项来指定其他的 yaml 格式的文件

cluster.yml #
