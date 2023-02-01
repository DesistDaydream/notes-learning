---
title: Node 资源预留与 Pod 驱逐
---

# 概述

> 参考：
> - 原文：[SegmentFault(思否)，k8s 节点资源预留与 pod 驱逐](https://segmentfault.com/a/1190000021402192)

# 节点资源管理

## 为什么

K8S 的节点上的资源会被 pod 和系统进程所使用，如果默认什么都不配置，那么节点上的全部资源都是可以分配给 pod 使用的，系统进程本身没有保障，这样做很危险：

- 集群雪崩：如果节点上调度了大量 pod，且 pod 没有合理的 limit 限制，节点资源将被耗尽，sshd、kubelet 等进程 OOM，节点变成 not ready 状态，pod 重新继续调度到其他节点，新节点也被打挂，引起集群雪崩。
- 系统进程异常：就算 pod 设置了 limit，但如果机器遇到资源不足，系统进程如 docker 没有资源保障，会频繁 OOM，或者进程 hang 住无响应，虽然能运行，但容器会反复出问题

节点资源主要分为两类：

- 可压缩资源：如 CPU，即使 cpu 超配，也可以划分时间片运行，只是运行变慢，进程不会挂。
- 不可压缩资源：Memory/Storage，内存不同于 CPU，系统内存不足时，会触发 OOM 杀死进程，按照 oom score 来确定先 kill 谁，oom_score_adj 值越高，被 kill 的优先级越高。

oom 分数：
![](https://notes-learning.oss-cn-beijing.aliyuncs.com/251afead-3706-4ed7-9dbe-ed79c2f36290/1460000021402196)
所以，OOM 的优先级如下：

    BestEffort Pod > Burstable Pod > 其它进程 > Guaranteed Pod > kubelet/docker 等 > sshd 等进程

因此需要对节点的内存等资源进行配置，以保证节点核心进程运行正常。

## 怎么做

节点资源的配置一般分为 2 种：

1. 资源预留：为系统进程和 k8s 进程预留资源
2. pod 驱逐：节点资源到达一定使用量，开始驱逐 pod

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/251afead-3706-4ed7-9dbe-ed79c2f36290/1460000021402198)

- Node Capacity：Node 的所有硬件资源
- kube-reserved：给 kube 组件预留的资源：kubelet,kube-proxy 以及 docker 等
- system-reserved：给 system 进程预留的资源
- eviction-threshold：kubelet eviction 的阈值设定
- Allocatable：真正 scheduler 调度 Pod 时的参考值（保证 Node 上所有 Pods 的 request resource 不超过 Allocatable）

allocatable 的值即对应 describe node 时看到的 allocatable 容量，pod 调度的上限

    计算公式：节点上可配置值 = 总量 - 预留值 - 驱逐阈值

    Allocatable = Capacity - Reserved(kube+system) - Eviction Threshold

以上配置均在 kubelet 中添加，涉及的参数有：

    --enforce-node-allocatable=pods,kube-reserved,system-reserved
    --kube-reserved-cgroup=/system.slice/kubelet.service
    --system-reserved-cgroup=/system.slice
    --kube-reserved=cpu=200m,memory=250Mi
    --system-reserved=cpu=200m,memory=250Mi
    --eviction-hard=memory.available<5%,nodefs.available<10%,imagefs.available<10%
    --eviction-soft=memory.available<10%,nodefs.available<15%,imagefs.available<15%
    --eviction-soft-grace-period=memory.available=2m,nodefs.available=2m,imagefs.available=2m
    --eviction-max-pod-grace-period=30
    --eviction-minimum-reclaim=memory.available=0Mi,nodefs.available=500Mi,imagefs.available=500Mi

## 配置含义

配置的含义如下：

（1）--enforce-node-allocatable

    含义：指定kubelet为哪些进程做硬限制，可选的值有：

    * pods
    * kube-reserved
    * system-reserve

    这个参数开启并指定pods后kubelet会为所有pod的总cgroup做资源限制(通过cgroup中的kubepods.limit\_in\_bytes)，限制为公式计算出的allocatable的大小。

    假如想为系统进程和k8s进程也做cgroup级别的硬限制，还可以在限制列表中再加system-reserved和kube-reserved，同时还要分别加上--kube-reserved-cgroup和--system-reserved-cgroup以指定分别限制在哪个cgroup里。


    配置：--enforce-node-allocatable\=pods,kube-reserved,system-reserved

（2）设置 k8s 组件的 cgroup

    含义：这个参数用来指定k8s系统组件所使用的cgroup。

    注意，这里指定的cgroup及其子系统需要预先创建好，kubelet并不会为你自动创建好。


    配置：--kube-reserved-cgroup=/system.slice/kubelet.service

（3）设置系统守护进程的 cgroup

    含义：这个参数用来指定系统守护进程所使用的cgroup。

    注意，这里指定的cgroup及其子系统需要预先创建好，kubelet并不会为你自动创建好。

    配置：--system\-reserved-cgroup=/system.slice

（4）配置 k8s 组件预留资源的大小，CPU、Mem

    指定为k8s系统组件（kubelet、kube-proxy、dockerd等）预留的资源量，

    如：--kube-reserved=cpu=1,memory=2Gi,ephemeral-storage=1Gi。

    这里的kube-reserved只为非pod形式启动的kube组件预留资源，假如组件要是以static pod（kubeadm）形式启动的，那并不在这个kube-reserved管理并限制的cgroup中，而是在kubepod这个cgroup中。

    （ephemeral storage需要kubelet开启feature\-gates，预留的是临时存储空间（log，EmptyDir），生产环境建议先不使用）

    ephemeral-storage是kubernetes1.8开始引入的一个资源限制的对象，kubernetes 1.10版本中kubelet默认已经打开的了,到目前1.11还是beta阶段，主要是用于对本地临时存储使用空间大小的限制，如对pod的empty dir、/var/lib/kubelet、日志、容器可读写层的使用大小的限制。

（5）配置 系统守护进程预留资源的大小（预留的值需要根据机器上容器的密度做一个合理的值）

    含义：为系统守护进程(sshd, udev等)预留的资源量，

    如：--system\-reserved=cpu=500m,memory\=1Gi,ephemeral-storage=1Gi。

    注意，除了考虑为系统进程预留的量之外，还应该为kernel和用户登录会话预留一些内存。


    配置：--system\-reserved=cpu=200m,memory\=250Mi

（6）配置 驱逐 pod 的硬阈值

    含义：设置进行pod驱逐的阈值，这个参数只支持内存和磁盘。

    通过

    kubelet 将会对pod进行驱逐。


    配置：--eviction-hard=memory.available<5%,nodefs.available<10%,imagefs.available<10%

（7）配置 驱逐 pod 的软阈值

    --eviction-soft=memory.available<10%,nodefs.available<15%,imagefs.available<15%

（8）定义达到软阈值之后，持续时间超过多久才进行驱逐

    --eviction-soft-grace-period=memory.available=2m,nodefs.available=2m,imagefs.available=2m

（9）驱逐 pod 前最大等待时间 = min(pod.Spec.TerminationGracePeriodSeconds, eviction-max-pod-grace-period)，单位为秒

    --eviction-max-pod-grace-period=30

（10）至少回收的资源量

    --eviction-minimum-reclaim=memory.available=0Mi,nodefs.available=500Mi,imagefs.available=500Mi

以上配置均为百分比，举例：

以 2 核 4GB 内存 40GB 磁盘空间的配置为例，Allocatable 是 1.6 CPU，3.3Gi 内存，25Gi 磁盘。当 pod 的总内存消耗大于 3.3Gi 或者磁盘消耗大于 25Gi 时，会根据相应策略驱逐 pod。

## 硬驱逐与软驱逐

### 硬驱逐

kubelet 利用 metric 的值作为决策依据来触发驱逐行为，下面内容来自于 Kubelet summary API。

一旦超出阈值，就会触发 kubelet 进行资源回收的动作（区别于软驱逐，有宽限期），指标如下：
![](https://notes-learning.oss-cn-beijing.aliyuncs.com/251afead-3706-4ed7-9dbe-ed79c2f36290/1460000021402197)

- nodefs： 机器文件系统
- imagesfs: Kubelet 能够利用 cAdvisor 自动发现这些文件系统，镜像存储空间

例如如果一个 Node 有 10Gi 内存，我们希望在可用内存不足 1Gi 时进行驱逐，就可以选取下面的一种方式来定义驱逐阈值：

- memory.available<10%
- memory.available<1Gi

可以配置百分比或者实际值，但是操作符只能使用小于号，即 <

### 软驱逐

软阈值需要和一个宽限期参数协同工作。当系统资源消耗达到软阈值时，这一状况的持续时间超过了宽限期之前，Kubelet 不会触发任何动作。如果没有定义宽限期，Kubelet 会拒绝启动。

另外还可以定义一个 Pod 结束的宽限期。如果定义了这一宽限期，那么 Kubelet 会使用 pod.Spec.TerminationGracePeriodSeconds 和最大宽限期这两个值之间较小的那个（进行宽限），如果没有指定的话，kubelet 会不留宽限立即杀死 Pod。

软阈值的定义包括以下几个参数：

- eviction-soft：驱逐阈值，例如 memory.available<1.5Gi，如果满足这一条件的持续时间超过宽限期，就会触发对 Pod 的驱逐动作。
- eviction-soft-grace-period：驱逐宽限期，例如 memory.available=1m30s，用于定义达到软阈值之后，持续时间超过多久才进行驱逐。
- eviction-max-pod-grace-period：达到软阈值之后，到驱逐一个 Pod 之前的最大宽限时间（单位是秒）

### 判断周期

Housekeeping interval 参数定义一个时间间隔，Kubelet 每隔这一段就会对驱逐阈值进行评估。

- housekeeping-interval：容器检查的时间间隔。

### 节点表现

如果触发了硬阈值，或者符合软阈值的时间持续了与其对应的宽限期，Kubelet 就会认为当前节点压力太大，下面的节点状态定义描述了这种对应关系。
![](https://notes-learning.oss-cn-beijing.aliyuncs.com/251afead-3706-4ed7-9dbe-ed79c2f36290/1460000021402200)
Kubelet 会持续报告节点状态的更新过程，这一频率由参数 —node-status-update-frequency 指定，缺省情况下取值为 10s。

如果一个节点的状况在软阈值的上下波动，但是又不会超过他的宽限期，将会导致该节点的状态持续的在是否之间徘徊，最终会影响降低调度的决策过程。

要防止这种状况，下面的标志可以用来通知 Kubelet，在脱离 pressure 之前，必须等待。

`eviction-pressure-transition-period` 定义了在脱离 pressure 状态之前要等待的时间

Kubelet 在把 pressure 状态设置为 False 之前，会确认在周期之内，该节点没有达到阈值

---

如果达到了驱逐阈值，并且超出了宽限期，那么 Kubelet 会开始回收超出限量的资源，直到回到阈值以内。

Kubelet 在驱逐用户 Pod 之前，会尝试回收节点级别的资源。如果服务器为容器定义了独立的 imagefs，他的回收过程会有所不同。

**有 Imagefs**
如果 nodefs 文件系统到达了驱逐阈值，kubelet 会按照下面的顺序来清理空间:

- 删除死掉的 Pod / 容器

如果 imagefs 文件系统到达了驱逐阈值，kubelet 会按照下面的顺序来清理空间:

- 删掉所有无用镜像

**没有 Imagefs**

如果 nodefs 文件系统到达了驱逐阈值，kubelet 会按照下面的顺序来清理空间。

1. 删除死掉的 Pod / 容器
2. 删掉所有无用镜像

## pod 驱逐策略

Kubelet 会按照下面的标准对 Pod 的驱逐行为进行评判：

- 根据服务质量：即 BestEffort、Burstable、Guaranteed
- 根据 Pod 调度请求的被耗尽资源的消耗量

接下来，Pod 按照下面的顺序进行驱逐（QOS）：

1. BestEffort：消耗最多紧缺资源的 Pod 最先驱逐。
2. Burstable：请求（request）最多紧缺资源的 Pod 被驱逐，如果没有 Pod 超出他们的请求，会驱逐资源消耗量最大的 Pod。
3. Guaranteed：请求（request）最多紧缺资源的 Pod 被驱逐，如果没有 Pod 超出他们的请求，会驱逐资源消耗量最大的 Pod。

参考 POD 的 QOS：[服务质量等级](https://link.segmentfault.com/?url=https%3A%2F%2Fk8smeetup.github.io%2Fdocs%2Ftasks%2Fconfigure-pod-container%2Fquality-service-pod%2F)

Guaranteed Pod 不会因为其他 Pod 的资源被驱逐。如果系统进程（例如 kubelet、docker、journald 等）消耗了超出 system-reserved 或者 kube-reserved 的资源，而且这一节点上只运行了 Guaranteed Pod，那么为了保证节点的稳定性并降低异常请求对其他 Guaranteed Pod 的影响，必须选择一个 Guaranteed Pod 进行驱逐。

本地磁盘是一个 BestEffort 资源。如有必要，kubelet 会在 DiskPressure 的情况下，kubelet 会按照 QoS 进行评估。如果 Kubelet 判定缺乏 inode 资源，就会通过驱逐最低 QoS 的 Pod 的方式来回收 inodes。如果 kubelet 判定缺乏磁盘空间，就会通过在相同 QoS 的 Pods 中，选择消耗最多磁盘空间的 Pod 进行驱逐。

---

**有 Imagefs**

- 如果 nodefs 触发了驱逐，Kubelet 会用 nodefs 的使用对 Pod 进行排序 – Pod 中所有容器的本地卷和日志。
- 如果 imagefs 触发了驱逐，Kubelet 会根据 Pod 中所有容器的消耗的可写入层进行排序。

---

**没有 Imagefs**

- 如果 nodefs 触发了驱逐，Kubelet 会对各个 Pod 的所有容器的总体磁盘消耗进行排序 —— 本地卷 + 日志 + 写入层。
- 在某些场景下，驱逐 Pod 可能只回收了很少的资源。这就导致了 kubelet 反复触发驱逐阈值。另外回收资源例如磁盘资源，是需要消耗时间的。
- 要缓和这种状况，Kubelet 能够对每种资源定义 minimum-reclaim。kubelet 一旦发现了资源压力，就会试着回收至少 minimum-reclaim 的资源，使得资源消耗量回到期望范围。

例如下面的配置：

    --eviction-hard=memory.available<500Mi,nodefs.available<1Gi,imagefs.available<100Gi

    --eviction-minimum-reclaim="memory.available=0Mi,nodefs.available=500Mi,imagefs.available=2Gi"

- 如果 memory.available 被触发，Kubelet 会启动回收，让 memory.available 至少有 500Mi。
- 如果是 nodefs.available，Kubelet 就要想法子让 nodefs.available 回到至少 1.5Gi。
- 而对于 imagefs.available， kubelet 就要回收到最少 102Gi。

缺省情况下，所有资源的 eviction-minimum-reclaim 为 0。

_在节点资源紧缺的情况下，调度器将不再继续向此节点部署新的 Pod_

## 节点 OOM 时

如果节点在 Kubelet 能够回收内存之前，遭遇到了系统的 OOM (内存不足)，节点就依赖 oom_killer 进行响应了。

kubelet 根据 Pod 的 QoS 为每个容器设置了一个 oom_score_adj 值。
![](https://notes-learning.oss-cn-beijing.aliyuncs.com/251afead-3706-4ed7-9dbe-ed79c2f36290/1460000021402201)
如果 kubelet 无法在系统 OOM 之前回收足够的内存，oom_killer 就会根据根据内存使用比率来计算 oom_score，得出结果和 oom_score_adj 相加，最后得分最高的 Pod 会被首先驱逐。

跟 Pod 驱逐不同，如果一个 Pod 的容器被 OOM 杀掉，他是可能被 kubelet 根据 RestartPolicy 重启的。

## Daemonset 的处理

因为 DaemonSet 中的 Pod 会立即重建到同一个节点，所以 Kubelet 不应驱逐 DaemonSet 中的 Pod。

但是目前 Kubelet 无法分辨一个 Pod 是否由 DaemonSet 创建。如果 Kubelet 能够识别这一点，那么就可以先从驱逐候选列表中过滤掉 DaemonSet 的 Pod。

一般来说，强烈建议 DaemonSet 不要创建 BestEffort Pod，而是使用 Guaranteed Pod，来避免进入驱逐候选列表。

## 已知问题

### Kubelet 无法及时监测到内存压力

Kubelet 目前从 cAdvisor 定时获取内存使用状况统计。如果内存使用在这个时间段内发生了快速增长，Kubelet 就无法观察到 MemoryPressure，可能会触发 OOMKiller。我们正在尝试将这一过程集成到 memcg 通知 API 中，来降低这一延迟，而不是让内核首先发现这一情况。

如果用户不是希望获得终极使用率，而是作为一个过量使用的衡量方式，对付这一个问题的较为可靠的方式就是设置驱逐阈值为 75% 容量。这样就提高了避开 OOM 的能力，提高了驱逐的标准，有助于集群状态的平衡。

### Kubelet 可能驱逐超出需要的更多 Pod

这也是因为状态搜集的时间差导致的。未来会加入功能，让根容器的统计频率和其他容器分别开来（[https://github.com/google/cad...](https://link.segmentfault.com/?url=https%3A%2F%2Fgithub.com%2Fgoogle%2Fcadvisor%2Fissues%2F1247)）。

### Kubelet 如何在 inode 耗尽的时候评价 Pod 的驱逐

目前不可能知道一个容器消耗了多少 inode。如果 Kubelet 觉察到了 inode 耗尽，他会利用 QoS 对 Pod 进行驱逐评估。在 cadvisor 中有一个 issue，来跟踪容器的 inode 消耗，这样我们就能利用 inode 进行评估了。例如如果我们知道一个容器创建了大量的 0 字节文件，就会优先驱逐这一 Pod

## 最佳实践

### 资源预留

1、资源预留需要设置，pod 的 limit 也要设置。
2、cpu 是可压缩资源，内存、磁盘资源是不可压缩资源。内存一定要预留，CPU 可以根据实际情况来调整
3、预留多少合适：根据集群规模设置阶梯，如下 (GKE 建议)：

Allocatable = Capacity - Reserved - Eviction Threshold

对于内存资源：

- 内存少于 1GB，则设置 255 MiB
- 内存大于 4G，设置前 4GB 内存的 25％
- 接下来 4GB 内存的 20％（最多 8GB）
- 接下来 8GB 内存的 10％（最多 16GB）
- 接下来 112GB 内存的 6％（最高 128GB）
- 超过 128GB 的任何内存的 2％
- 在 1.12.0 之前的版本中，内存小于 1GB 的节点不需要保留内存

对于 CPU 资源：

- 第一个核的 6％
- 下一个核的 1％（最多 2 个核）
- 接下来 2 个核的 0.5％（最多 4 个核）
- 4 个核以上的都是总数的 0.25％

对于磁盘资源（不是正式特性，仅供参考）：
![](https://notes-learning.oss-cn-beijing.aliyuncs.com/251afead-3706-4ed7-9dbe-ed79c2f36290/1460000021402199)
效果：查看节点的可分配资源：

    kubectl describe node \[NODE\_NAME\] | grep Allocatable -B 4 -A 3

### 驱逐配置

    \--eviction-hard=memory.available<5%,nodefs.available<10%,imagefs.available<10%

    --eviction-soft\=memory.available<10%,nodefs.available<15%,imagefs.available<15%

    --eviction-soft\-grace-period=memory.available=2m,nodefs.available=2m,imagefs.available=2m

    --eviction-max\-pod-grace-period=30

    --eviction-minimum-reclaim=memory.available=0Mi,nodefs.available=500Mi,imagefs.available=500Mi

原文链接：[http://www.xuyasong.com/?p=1725](https://link.segmentfault.com/?url=http%3A%2F%2Fwww.xuyasong.com%2F%3Fp%3D1725)

## Reference

- [https://kubernetes.io/docs/admin/out-of-resource/](https://link.segmentfault.com/?url=https%3A%2F%2Fkubernetes.io%2Fdocs%2Fadmin%2Fout-of-resource%2F)
- [https://k8smeetup.github.io/docs/tasks/configure-pod-container/quality-service-pod/](https://link.segmentfault.com/?url=https%3A%2F%2Fk8smeetup.github.io%2Fdocs%2Ftasks%2Fconfigure-pod-container%2Fquality-service-pod%2F)
- [https://cloud.google.com/kubernetes-engine/docs/concepts/cluster-architecture#node_allocatable](https://link.segmentfault.com/?url=https%3A%2F%2Fcloud.google.com%2Fkubernetes-engine%2Fdocs%2Fconcepts%2Fcluster-architecture%23node_allocatable)
