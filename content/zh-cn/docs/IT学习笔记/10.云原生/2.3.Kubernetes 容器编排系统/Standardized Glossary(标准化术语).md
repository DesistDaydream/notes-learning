---
title: Standardized Glossary(标准化术语)
---

# 概述

> 参考：
> - [官方文档](https://kubernetes.io/docs/reference/glossary/)

## [Declarative Application Management](https://github.com/kubernetes/community/blob/master/contributors/design-proposals/architecture/declarative-application-management.md)

**Declarative Application Management(声明式应用管理)** 是一种部署和管理应用程序的方式。

## [kubeconfig](https://kubernetes.io/docs/concepts/configuration/organize-cluster-access-kubeconfig/)

kubeconfig 是用于保存集群访问信息的文件，这是引用配置文件的通用方法，并不表示一定会有一个名为 kubeconfig 的文件。
kubeconfig 文件用来组织有关集群、用户、明哼空间的信息和身份验证机制。kubectl 命令行工具使用 kubeconfig 文件来与 Kubernetes 集群进行交互。Kuberntes 集群的某些主要组件，也会使用 kubeconfig 文件进行交互，比如使用 kubeadm 工具部署的 kubernetes 集群，在每个节点的 /etc/kubernetes 目录下，就会有以 .conf 文件结尾的 kubeconfig 文件，以供 kubelet、scheduler、controller-manager 等组件使用。

## [Manifest](https://kubernetes.io/docs/reference/glossary/?fundamental=true#term-manifest)

JSON 或 YAML 格式的 Kubernetes API 对象的规范。
manifest 指定了应用该 manifest 时，Kubernetes 将维护的对象的所需状态。每个配置文件可以包含多个清单
