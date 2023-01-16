---
title: Kubernetes 衍生品
---

# ClusterAPI 概述

参考：[官方文档](https://cluster-api.sigs.k8s.io/)

## 为什么要建立集群 API？

Kubernetes 是一个复杂的系统，它依赖于正确配置的几个组件才能具有正常运行的集群。社区意识到这是用户的潜在绊脚石，因此专注于简化引导过程。如今，已经创建了[100 多个 Kubernetes 发行版和安装程序](https://www.cncf.io/certification/software-conformance/)，每个[发行版和安装程序](https://www.cncf.io/certification/software-conformance/)都为集群和受支持的基础架构提供程序提供了不同的默认配置。SIG 集群生命周期发现需要一种工具来解决一系列常见的重叠安装问题，因此开始使用 kubeadm。
[Kubeadm](https://kubernetes.io/docs/reference/setup-tools/kubeadm/kubeadm/)被设计为引导最佳实践 Kubernetes 集群的重点工具。kubeadm 项目背后的核心宗旨是创建其他安装程序可以利用的工具，并最终减轻单个安装程序需要维护的配置量。自开始以来，kubeadm 已成为其他多个应用程序（包括 Kubespray，Minikube，kind 等）的基础自举工具。
但是，尽管 kubeadm 和其他引导程序提供程序降低了安装复杂性，但它们并未解决如何长期管理日常群集或 Kubernetes 环境的问题。在设置生产环境时，您仍然面临几个问题，包括

- 如何在多个基础架构提供商和位置之间一致地配置计算机，负载平衡器，VPC 等？
- 如何实现集群生命周期管理的自动化，包括升级和集群删除等操作？
- 如何扩展这些过程以管理任意数量的集群？

SIG 集群生命周期开始了 ClusterAPI 项目，以此作为通过构建声明性的 Kubernetes 风格的 API 来解决这些差距的方法，该 API 使集群的创建，配置和管理自动化。使用此模型，还可以扩展集群 API，以支持所需的任何基础结构提供程序（AWS，Azure，vSphere 等）或引导程序提供程序（默认为 kubeadm）。请参阅越来越多的[可用提供程序](https://cluster-api.sigs.k8s.io/reference/providers.html)列表。

# CNCF 的 Software conformance(软件一致性)

参考：[官方文档](https://www.cncf.io/certification/software-conformance/)
Certified Kubernetes(经过认证的 Kubernetes)
对于使用 Kubernetes 的组织，一致性可以实现从一个 Kubernetes 安装到下一个 Kubernetes 安装的互操作性。它使他们可以灵活地在供应商之间进行选择。

CNCF 运行 Kubernetes 认证合格计划。大多数全球领先的企业软件供应商和云计算提供商都拥有 [经过认证的 Kubernetes](https://www.cncf.io/certification/software-conformance/#logos) 产品。

**有超过 90 种经过认证的 Kubernetes 产品。**邀请所有供应商提交一致性测试结果，以供 CNCF 审核和认证。如果您的公司提供基于 Kubernetes 的软件，我们建议您立即获得认证。

# Kubernetes 管理工具

[Rancher](https://www.yuque.com/go/doc/33161032)
[Kuboard](https://github.com/eip-work/kuboard-press)

# Kubernetes 发行版

K3S
