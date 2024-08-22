---
title: Scheduler(调度器)
---

# 概述

> 参考：
>
> - [官方文档, 概念 - 调度与驱逐](https://kubernetes.io/docs/concepts/scheduling-eviction/)

**Scheduler(调度器)** 负责决定 Pod 与 Node 的匹配关系，并将 Pod 调度到匹配到的 Node 上，以便 Kubelet 可以运行这些 Pod。Scheduler 在调度时会充分考虑 Cluster 的拓扑结构，当前各个节点的负载，以及应用对高可用、性能、数据亲和性的需求。

Scheduler 通过 Kubernetes 的 watch 机制来发现集群中新创建且尚未被调度到 Node 的 Pod。Scheduler 会将发现的每一个未调度的 Pod 调度到一个合适的 Node 上去运行。

# 调度器的实现

[kube-scheduler](/docs/10.云原生/Kubernetes/Scheduling/kube-scheduler%20实现调度器的程序.md) 是实现 Kubernetes 调度器功能的程序。

# Scheduler 调度策略

> 参考：
>
> - [官方文档, 参考 - 调度 - 调度策略](https://kubernetes.io/docs/reference/scheduling/policies/)

Scheduler 调度的时候，通过以下步骤来完成调度

1. Predicate(预选) # 排除那些不满足 Pod 运行环境的 Node
2. Priorities(优选) # 通过算法，剩余可运行 Pod 的 Node 进行计算后排序，选择结果最高的 Node
3. Select(选定) # 若优选后有多个 Node 得分相同，则随机挑选，将选择的结果告诉 APIServer 用哪个 Node 部署 Pod

调度倾向性：亲合性 Affinity，反亲合性 AntiAffinity，污点 Taints，容忍度 Tolerations

预选策略部分说明：预选策略中需要检查所有项，其中任意一项不通过，则不满足部署条件

1. CheckNodeCondition
2. GeneralPrediactes
   1. HostName
   2. PodFitsHostPorts
   3. MatchNodeSelector
3. NoDiskConflict
4. PodToleratesNodeTaints：检查 Pod 上的 spec.toleratons 可容忍的污点是否完全包含 Node 上的污点，Node 上后加的污点就算不在可容忍范围也可接受
5. PodToleratesNodeNoExecuteTaints：检查 Pod 上的 spec.toleratons 可容忍的污点是否完全包含 Node 上的污点，Node 上后加的污点如果不在可容忍范围，则会被 Node 驱除
6. CheckNodeLabelPresence：检查 Node 标签存在性
7. CheckServiceAffinity：检查 Service 亲合性
8. CheckVolumeBinding：检查 Node 已绑定和未绑定的 PVC
9. NoVolumeZoneConflict：
10. CheckNodeMemoryPressure：检查节点的内存资源压力是否处于过大状态
11. CheckNoePIDPressure：检查节点的 PID 压力是否处于过大状态
12. CheckNodeDiskPressure：检查节点的磁盘压力是否处于过大状态
13. MatchInterPodAffinity：匹配 Pod 内部的亲合性(哪些 Pod 更倾向运行在同一节点等)

优选函数部分说明：

1. LeastRequested：最少被请求的，得分越低，优先级越高
   1. （cpu((capacity-sum(requested))*10/capacity)+（memory((capacity-sum(requested))*10/capacity))/2 # CPU 与内存分别计算（总容量-已被请求的）总容量\*10，并求和，再除以 2,即为该 Node 资源得分
2. MostRequested：得分越多，优先级越高，
3. BalancedResourceAllocation：均衡资源分配
   1. CPU 和内存资源被占用率相近的胜出
4. NodePreferAvoidPods：
   1. 节点注解信息“scheduler.alpha.kubernetes.io/preferAvoidPods”
5. TaintToleration：将 Pod 对象的 spec.tolerations 列表项与 Node 的 taints 列表项进行匹配度检查，匹配的条目越多得分越低
6. SelectorSpreading：与当前 Pod 对象同属相同选择器的 Pod 越多，得分越低
7. InterPodAffinity：
