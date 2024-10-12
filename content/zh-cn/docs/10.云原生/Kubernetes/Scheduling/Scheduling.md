---
title: Scheduling
linkTitle: Scheduling
date: 2020-04-02T10:36:00
weight: 1
---

# 概述

> 参考：
>
> - [官方文档, 概念 - 调度、抢占与驱逐](https://kubernetes.io/docs/concepts/scheduling-eviction/)

**Scheduling(调度)** 是一个行为，用来让 Pod 匹配到 Node，以便 Node 上的 Kubelet 可以运行这些 Pod。如果没有调度系统，Kubernetes 集群就不知道 Pod 应该运行在哪里。这种调度的概念，与 Linux 中调度任务来使用 CPU 是一个意思。可以看看 [Scheduler](/docs/8.通用技术/Scheduler.md) 相关文章，调度是在 IT 行业中，很多程序都很重要的概念。

与 Scheduling(调度) 伴生的，还有 **Preemption(抢占)** 与 **Eviction(驱逐)** 两个概念。顾名思义：

- **Preemption(抢占)** 是指终止优先级较低的 Pod 的行为，以便优先级较高的 Pod 可以在节点上调度。
  - 抢占行为通常发生在资源不足时，当一个新 Pod 需要调度，但是资源不足，那么就可能需要抢占优先级低的 Pod，这个低优先级的 Pod 将会被驱逐，以便让优先级高的 Pod 运行在节点上。
- **Eviction(驱逐)** 是指终止节点上一个或多个 Pod 的行为。

由 抢占 与 驱逐 两个行为，还引申出了 **Pod Disruption(中断)** 的概念。[Pod Disruption(中断)](https://kubernetes.io/docs/concepts/workloads/pods/disruptions/) 是指节点上的 Pod 自愿或者非资源终止运行的行为。

- 自愿中断是由应用程序所有者或者集群管理故意启动的(比如.维护节点前手动驱逐 Pod)
- 非自愿中断是无意的，可能由不可避免的问题触发(比如.节点资源耗尽或意外删除)
