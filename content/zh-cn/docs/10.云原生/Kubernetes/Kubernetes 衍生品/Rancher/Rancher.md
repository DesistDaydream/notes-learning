---
title: Rancher
linkTitle: Rancher
date: 2023-11-23T11:48
weight: 20
---

# 概述

> 参考：
>
> - 官方文档：<https://rancher.com/>

Rancher 是为使用容器的公司打造的容器管理平台。Rancher 简化了使用 Kubernetes 的流程，开发者可以随处运行 Kubernetes（Run Kubernetes Everywhere），满足 IT 需求规范，赋能 DevOps 团队。

Rancher 在现阶段可以看作是一个解决方案，是一套产品的统称，这套产品包括如下几个：

- [K3S](/docs/10.云原生/Kubernetes/Kubernetes%20衍生品/K3S/K3S.md) # 用于运行高可用 Rancher 的底层平台。是一个轻量的 kubernetes，一个 k3s 二进制文件即可包含所有 kubernetes 的主要组件。
- Rancher Server # Rancher 管理程序，常部署于 k3s 之上，用来管理其下游 k8s 集群。
- RKE # Rancher 创建的 kubernetes 集群。是一个可以通过名为 rke 的二进制文件以及一个 yaml 文件，即可启动 kubernetes 集群的引擎。RKE 与 kubernetes 的关系，类似于 docker 与 containerd 的关系。

## Rancher Server 介绍

Rancher Server 由认证代理(Authentication Proxy)、Rancher API Server、集群控制器(Cluster Controller)、数据存储(比如 etcd、mysql 等)和集群代理(Cluster Agent) 组成。除了 Cluster Agent 以外，其他组件都部署在 Rancher Server 中。(这些组件都集中在一起，一般可以通过 docker 直接启动一个 Rancher Server。)

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/kxmsmg/1616114814016-9de5267d-0813-4790-826c-7c4448e40861.png)

Rancher Server 可以管理多种 k8s 集群

1. 通过 Rancher Server 来创建一个 RKE 集群
2. 托管的 kubernetes 集群。e.g.Amazon EKS、Azure AKS、Google GKE 等等
3. 导入已有的 kubernetes 集群。

## Rancher 与下游集群交互的方式

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/kxmsmg/1616114813966-db373999-6c8f-4541-a09f-5f20eaa656ce.png)

通过 Rancher 管理的 kubernetes 集群(不管是导入的还是通过 Rancher 创建的)，都会在集群中部署两种 agent，来与 Rancher 进行交互。

- cattle-cluster-agent # 上图中的 Cluster Agent。用于本集群与 Rancher Server 的 Cluster Controller(集群控制器)的通信
- 连接 Rancher 与本集群的 API Server
- 管理集群内的工作负载，比如 Rancher Server 下发一个部署 pod 的任务，集群代理就会与本集群 API 交互来处理任务
- 根据每个集群的设置，配置 Role 和 RoleBindings
- 实现集群和 Rancher Server 之间的消息传输，包括事件，指标，健康状况和节点信息等。
- cattle-node-agent # 上图中的 Node Agent。用于处理本节点的任务，比如升级 kubernetes 版本以及创建或者还原 etcd 快照等等。
- Note：如果 Cluster Agent 不可用，下游集群中的其中一个 Node Agent 会创建一个通信管道，由节点 Agent 连接到集群控制器，实现下游集群和 Rancher 之间的通信。
- 一般使用 DaemonSet 的方式部署到集群中，以保证每个节点都有一个代理可以执行 Rancher Server 下发的任务。

# Rancher 部署

详见 [Rancher 部署与清理](/docs/10.云原生/Kubernetes/Kubernetes%20衍生品/Rancher/Rancher%20部署与清理.md)

# Rancher 关联文件与配置

Rancher 套件中的各组件配置详见各自组件配置详解

详见 [Rancher 配置](/docs/10.云原生/Kubernetes/Kubernetes%20衍生品/Rancher/Rancher%20配置.md)

## K3S 配置

## Rancher Server 配置

## Rancher 创建的集群配置

Rancher 创建的集群是为 RKE 集群，配置详见：RKE 配置详解
