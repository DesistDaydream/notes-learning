---
title: Garbage Collection(垃圾回收)
---

# kubelet Garbage Collection 介绍

> 参考：
> - 官方文档：<https://kubernetes.io/docs/concepts/cluster-administration/kubelet-garbage-collection/>

垃圾回收是 kubelet 的一个有用功能，它将清理未使用的镜像和容器。 Kubelet 将每分钟对容器执行一次垃圾回收，每五分钟对镜像执行一次垃圾回收。

注意！！不建议使用外部垃圾收集工具，因为这些工具可能会删除原本期望存在的容器进而破坏 kubelet 的行为。

比如：
使用 docker container prune -f 命令，清理了节点上不再使用的容器，这时候，在 /var/log/pods/ContainerNAME/\* 目录下的日志软链接是不会清除的，因为这个软连接由 kubelet 管理，并且只有在日志关联的容器被 kubelet 清理时，才会清理该软链接。所以容器没了，软连接 kubelet 也就不管了~

解决办法：
find -L /var/log/pods -type l -delete 直接直接该命令即可

代码路径：./pkg/kubelet/kuberuntime/kuberuntime_gc

## 镜像回收

Kubernetes 借助于 cadvisor 通过 imageManager 来管理所有镜像的生命周期。

镜像垃圾回收策略只考虑两个因素：HighThresholdPercent 和 LowThresholdPercent。 磁盘使用率超过上限阈值（HighThresholdPercent）将触发垃圾回收。 垃圾回收将删除最近最少使用的镜像，直到磁盘使用率满足下限阈值（LowThresholdPercent）。

## 容器回收

容器垃圾回收策略考虑三个用户定义变量。 MinAge 是容器可以被执行垃圾回收的最小生命周期。 MaxPerPodContainer 是每个 pod 内允许存在的死亡容器的最大数量。 MaxContainers 是全部死亡容器的最大数量。 可以分别独立地通过将 MinAge 设置为 0，以及将 MaxPerPodContainer 和 MaxContainers 设置为小于 0 来禁用这些变量。

kubelet 将处理无法辨识的、已删除的以及超出前面提到的参数所设置范围的容器。 最老的容器通常会先被移除。 MaxPerPodContainer 和 MaxContainer 在某些场景下可能会存在冲突， 例如在保证每个 pod 内死亡容器的最大数量（MaxPerPodContainer）的条件下可能会超过 允许存在的全部死亡容器的最大数量（MaxContainer）。 MaxPerPodContainer 在这种情况下会被进行调整： 最坏的情况是将 MaxPerPodContainer 降级为 1，并驱逐最老的容器。 此外，pod 内已经被删除的容器一旦年龄超过 MinAge 就会被清理。

不被 kubelet 管理的容器不受容器垃圾回收的约束。

# 用户配置

用户可以使用以下 kubelet 参数调整相关阈值来优化镜像垃圾回收：

- image-gc-high-threshold #触发镜像垃圾回收的磁盘使用率百分比。默认值为 85%。
- image-gc-low-threshold #镜像垃圾回收试图释放资源后达到的磁盘使用率百分比。默认值为 80%。

我们还允许用户通过以下 kubelet 参数自定义垃圾收集策略：

- minimum-container-ttl-duration # 完成的容器在被垃圾回收之前的最小年龄，默认是 0 分钟。 这意味着每个完成的容器都会被执行垃圾回收。
- maximum-dead-containers-per-container # 每个容器要保留的旧实例的最大数量。默认值为 1。
- maximum-dead-containers # 要全局保留的旧容器实例的最大数量。 默认值是 -1，意味着没有全局限制。

容器可能会在其效用过期之前被垃圾回收。这些容器可能包含日志和其他对故障诊断有用的数据。 强烈建议为 maximum-dead-containers-per-container 设置一个足够大的值，以便每个预期容器至少保留一个死亡容器。 由于同样的原因，maximum-dead-containers 也建议使用一个足够大的值。

查阅这个 Issue 获取更多细节。

弃用

这篇文档中的一些 kubelet 垃圾收集（Garbage Collection）功能将在未来被 kubelet 驱逐回收（eviction）所替代。

包括:

|                                         |                                       |                                            |
| --------------------------------------- | ------------------------------------- | ------------------------------------------ |
| 现存参数                                | 新参数                                | 解释                                       |
| --image-gc-high-threshold               | --eviction-hard 或 --eviction-soft    | 现存的驱逐回收信号可以触发镜像垃圾回收     |
| --image-gc-low-threshold                | --eviction-minimum-reclaim            | 驱逐回收实现相同行为                       |
| --maximum-dead-containers               |                                       | 一旦旧日志存储在容器上下文之外，就会被弃用 |
| --maximum-dead-containers-per-container |                                       | 一旦旧日志存储在容器上下文之外，就会被弃用 |
| --minimum-container-ttl-duration        |                                       | 一旦旧日志存储在容器上下文之外，就会被弃用 |
| --low-diskspace-threshold-mb            | --eviction-hard or eviction-soft      | 驱逐回收将磁盘阈值泛化到其他资源           |
| --outofdisk-transition-frequency        | --eviction-pressure-transition-period | 驱逐回收将磁盘压力转换到其他资源           |
