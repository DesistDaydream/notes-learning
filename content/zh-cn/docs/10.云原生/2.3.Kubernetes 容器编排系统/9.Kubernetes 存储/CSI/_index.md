---
title: CSI
---

# 概述

> 参考：
>
> - [GitHub 项目,规范](https://github.com/container-storage-interface/spec)

CSI 是与 [CNI](/docs/10.云原生/2.3.Kubernetes%20容器编排系统/8.Kubernetes%20网络/CNI/CNI.md) 类似的东西，都是一种规范。

Container Storage Interface，容器存储接口（CSI）为容器编排系统（如 Kubernetes）定义了一个标准接口，以将任意存储系统暴露给其容器工作负载。

## CSI Specification 规范介绍

官方文档：<https://github.com/container-storage-interface/spec/blob/master/spec.md>

# 背景

CSI 出现之前，很多存储类型的 PV，比如 iSCSI、NFS、CephFS 等等([详见 k8s 官方支持列表](https://kubernetes.io/docs/concepts/storage/persistent-volumes/#types-of-persistent-volumes))。这些类型的存储的代码，放在 Kubernetes 代码仓库中，这种称为 **in-tree** 类型的存储(也就是在代码树中)。这里代理的问题是 Kubernetes 代码与 第三方存储厂商的代码**强耦合**。

- 更改 in-tree 类型的存储代码，用户必须更新 K8s 组件，成本较高
- in-tree 存储代码中的 bug 会引发 K8s 组件不稳定
- K8s 社区需要负责维护及测试 in-tree 类型的存储功能
- in-tree 存储插件享有与 K8s 核心组件同等的特权，存在安全隐患
- 三方存储开发者必须遵循 K8s 社区的规则开发 in-tree 类型存储代码

**CSI(容器存储接口)**标准的出现解决了上述问题，将三方存储代码与 K8S 代码解耦，使得三方存储厂商研发人员只需实现 CSI 接口（无需关注容器平台是 K8s 还是 Swarm 等）。
