---
title: EvictionManager 模块
---

# 概述

> 参考：
> - [公众号-云原生实验室，Kubernetes 单机侧的驱逐策略总结](https://mp.weixin.qq.com/s/ehECtQiXSHLpCrH5vuBX_w)
>   - 本文转自 Edwardesire 的博客，原文：[**https://edwardesire.com/posts/process-eviction-under-k8s/**](https://edwardesire.com/posts/process-eviction-under-k8s/)

进程驱逐：当机器存在资源压力时，可能是由于有恶意程序在消耗系统资源，或者是 overcommit 导致。系统通过控制机器上的进程存活来减少单个程序对系统的整体影响。驱逐阶段最关键的就是选择合适的进程，通过最小代价来保证系统的稳定。在执行层面上可以分为两种驱逐方式：

- 用户空间驱逐：通过守护进程之类的机制，触发式主动清理进程。
- 内核空间驱逐：内核在无法分配内存时，通过 oom_killer 选择进程终止来释放资源。

本文从 k8s 出发，总结不同层次下的驱逐流程和进程选择策略。

## Kubelet 驱逐策略

k8s 除了支持 API 发起的主动驱逐，还支持用户空间的 pod 驱逐（将资源大户的进程终止）。对于不可压缩资源：内存、disk（nodefs）、pid，kubelet 会监控相应的指标来触发 pod 驱逐。K8S 依据 pod 的资源消耗和优先级来驱逐 pod 来回收资源：

- 如果 pod 资源使用量超过资源请求值，则优先驱逐
- 依据 pod priority 驱逐
- pod 真实资源使用量越高则越优先驱逐

我们可以得出：

- 当 BestEffort 和 Burstable pod 的资源使用量超过请求值时，会依据 pod  priority 和超过请求多少来判断驱逐顺序。也不会有特例的 pod 能够不被驱逐的风险。当 Guaranteed 和 Busrtable  pod 的使用量低于请求值时，基于 pod priority 确定驱逐顺序。

这一切的逻辑都在 kubelet 的 eviction manager 实现。

### Eviction manager

Manager 的接口定义包含主流程的启动函数以及提供给 kubelet 上报节点状态的：

- Start()：开始驱逐控制循环，获取监控数据，并判断资源是否到阈值，触发 pod 的驱逐，以及当节点出现压力时将本地的节点状态更新。
- IsUnderMemoryPressure()：判断节点是否达到内存限制压力，通过控制循环内更新的节点状态判断。
- IsUnderDiskPressure()：判断节点是否达到磁盘限制压力，通过控制循环内更新的节点状态判断。
- IsUnderPIDPressure()：判断节点是否达到 PID 限制压力，通过控制循环内更新的节点状态判断。

kubelet 在 tryUpdateNodeStatus 上报节点状态循环中，会调用上述方法判断节点的资源压力情况。

kubelet 在初始化 evictionManager 之后会调用 evictionManager.Start() 启动驱逐，之后再同步节点状态时调用上述压力判断方法。除了实现 Manager 的接口外，还实现了在 pod 生命周期负责评估允许 pod 执行的 PodAdmitHandler 接口。evictionManager 主要是依据 pod 的性质判断是否能够在已经有资源压力的机器上创建容器。

### 驱逐控制循环

#### 初始化阶段

kubelet 主程解析配置并初始化 evictionManager，解析单机的资源阈值参数 ParseThresholdConfig()

kubelet 以 signal 维度设定资源阈值，每个 signal 标识一种资源指标，定义资源的阈值和其他驱逐参数。比如 `memory.available` 表示节点可用内存驱逐标识（memory.available = capacity - workingSet）。

kubelet 通过下列参数确定资源 signal 属性，构建相应资源的阈值。

- `--eviction-hard mapStringString`：资源驱逐硬下线，默认为：imagefs.available<15%,memory.available<100Mi,nodefs.available<10%
- `--eviction-soft mapStringString`：资源驱逐的软下线，当触发时，pod 有优雅退出时间。
- `--eviction-soft-grace-period mapStringString`：触发黄线时，pod 驱逐的优雅退出时间。
- `--eviction-minimum-reclaim mapStringString`：资源的最小释放量。默认为 0。

其中，同一个资源的 eviction-soft 和 soft-grace-period 配置必须都存在，`grace period must be specified for the soft eviction threshold`。

通过解析配置项，设置各资源 signal 的阈值之后，kubelet 调用 evictionManager.Start() 驱动 evictionManager 工作。

#### evictionManager 的启动

在启动控制循环之前，evictionManager 会增加对 cgroup 内存子系统监控的预处理。这个预处理通过 **cgroup notifier** 的机制监听 mem cgroup 的使用情况，并且在控制循环中周期性更新 cgroup notifier 的阈值配置。

##### MemoryThresholdNotifier

evictionManager 分别给 `memory.available` 和 `allocatableMemory.available` signal 配置 MemoryThresholdNotifier 的，监控的 cgroup 路径不同。`allocatableMemory.available` 的 cgroupRoot 根目录，即节点上 pods 的根 cgroup。`memory.available` 则监控 `/proc/cgroups/memory` 目录。

MemoryThresholdNotifier 的工作流程是：

- 初始化 MemoryThresholdNotifier
  MemoryThresholdNotifier 需要获取 cgroup 目录的 cgoup 内存子系统路径，并设置 evictionManager.synchronize() 为阈值处理函数 thresholdHandler
- 创建 goroutine 来启动 MemoryThresholdNotifier
  在 MemoryThresholdNotifier.Start() 循环中：监听事件 channel，并调用驱逐函数（调用 synchronize）
- 在 synchronize 阶段调用 UpdateThreshold() 更新 memcg 的阈值，并激活 MemoryThresholdNotifier。
  依据当前的采集指标配置，计算 cgroup 内存使用阈值。
  如果 MemoryThresholdNotifier 已经存在 notifier 实例，则创建新的 cgroupNotifier 替换。cgroupNotifier 通过 epoll 上述 eventfd 描述符的方式，监听内存超过阈值的事件。

这里有两个关键点：

1、在 UpdateThreshold 函数中计算 cgroup 内存使用阈值

如上述，通过监听 memory.usage_in_bytes 文件，获取内存使用情况（不包含 swap），当内存使用阈值。而内存使用阈值 **memcgThreshold** 通过监控数据得来：

`// Set threshold on usage to capacity - eviction_hard + inactive_file,  // since we want to be notified when working_set = capacity - eviction_hard  inactiveFile := resource.NewQuantity(int64(*memoryStats.UsageBytes-*memoryStats.WorkingSetBytes), resource.BinarySI)  capacity := resource.NewQuantity(int64(*memoryStats.AvailableBytes+*memoryStats.WorkingSetBytes), resource.BinarySI)  evictionThresholdQuantity := evictionapi.GetThresholdQuantity(m.threshold.Value, capacity)  memcgThreshold := capacity.DeepCopy()  memcgThreshold.Sub(*evictionThresholdQuantity)  memcgThreshold.Add(*inactiveFile)`

计算内存使用阈值 **memcgThreshold** 的绝对值通过 capacity - eviction_hard（如果红线不为绝对值，则依据 capacity \* 百分比） + inactive_file 计算而来。

其中

- 内存容量 capacity = memoryStats.AvailableBytes + memoryStats.WorkingSetBytes，即 内存可用量 + workload 已使用量（两个值都从监控模块得到）
- 硬下线 eviction_hard 是参数值
- 不活跃的文件内存页 inactive_file = memoryStats.UsageBytes - memoryStats.WorkingSetBytes，即   内存已使用量 - workload 已使用量（包含最近使用的内存、待回收的脏内存和内核占用内存，两个值也从监控模块得到）。

2、在 UpdateThreshold 函数中创建 cgroupNotifier

**cgroup notifier** 的机制是通过 eventfd 监听 cgroup 中内存使用超过阈值的事件。

- memory.usage_in_bytes：监听内存使用文件对象。
- cgroup.event_control：阈值监控控制接口，依据 `<event_fd> <fd of memory.usage_in_bytes> <threshold>` 的格式配置 event_fd，watchfd 和阈值 threshold。

`/sys/fs/cgroup/memory $ cat memory.usage_in_bytes 92459601920 $ ls -lt cgroup.event_control --w--w--w- 1 root root 0 Nov 24 12:05 cgroup.event_control     # an interface for event_fd()`

cgroupNotifier 会依据 cgroup 事件向 channel 压入事件，触发事件消费者（evictionManager）处理。这里 channel 不会传递具体的事件内容，只做任务触发功能。

注册 cgroup 的 threshold，需要有 3 步：

- 使用 eventfd(2) 创建 eventfd
- 创建打开 memory.usage_in_bytes 或者 memory.memsw.usage_in_bytes 文件描述符
- 在 cgroup.event_control 写入 "\<event_fd> " 信息

在 evictionManager.Start() 的最后启动控制循环 synchronize 周期性检查驱逐的阈值条件是否达到，并进行下一步动作。

#### 控制循环 synchronize

在 evictionManager 的控制循环中，维持 10s 调用 synchronize 函数来选择 pod 驱逐。驱逐首要判断的就是驱逐的触发条件，通过监控系统资源的方式来判断资源使用情况是否触及阈值。evictionManager 有两种触发方式：

1、基于 cgroup 触发驱逐（基于事件）：上述已经描述了内存的 CgroupNotifier 机制

2、依据监控数据触发驱逐（周期性检查）

2.1 通过 summaryProvider 获取节点和 pods 的资源使用情况

2.2 在 signalObservations 函数中依据监控数据，获取各资源的使用情况 signalObservations

单个 signalObservation 记录着资源的总量和可用量：

`// signalObservation is the observed resource usage type signalObservation struct {  // The resource capacity  capacity *resource.Quantity  // The available resource  available *resource.Quantity  // Time at which the observation was taken  time metav1.Time }`

2.3 在 thresholdsMet 函数中判断是否需要驱逐来释放资源

当上述观测到的资源可用量低于各 signal 的阈值时，返回需要释放的资源类型。

无论哪种方式，都会执行 synchronize 后段逻辑来判断是否需要驱逐 pod。

3、更新节点状态，将资源压力状态更新，并上报到集群 API

集群内其他组件能够观测到节点状态，从节点外部处理。

4、如果开启了 featuregate LocalStorageCapacityIsolation 本地存储，会首先尝试清理影响本地磁盘

这个是依据 featuregate 来控制是否开启，会检查 pod 下列资源使用是否超过 limit 值。

- emptyDir 的 sizeLimit
- ephemeralStorage 的 limit
- container 的 ephemeralStorage limit

这种驱逐是立即的，没有优雅退出时间。当触发到本地磁盘触发条件时，会忽略其他资源的驱逐行为。

当驱逐流程走到这，会判断是否存在资源紧张的驱逐资源。如果 thresholdsMet 返回的空数组，则表示没有资源触及到驱逐阈值。否则继续执行节点资源的回收。

5、回收节点级别的资源

5.1 reclaimNodeLevelResource：回收节点级别资源

首先尝试回收节点资源：nodefs/imagefs，这部分可以通过删除没使用的容器和镜像，而不侵害执行中的 pod。调用完节点资源回收函数之后，再采集一次指标。如果空闲资源大于阈值，则跳过本次驱逐的后续流程：pod 级别的驱逐。

5.2 rank 阶段：判断触发驱逐条件的资源优先级

每次 synchronize 只会选择一个超过阈值的资源进行回收。当多个资源出现触碰到阈值时，资源驱逐优先级如下：

- 内存资源的驱逐优先级最高
- 没有资源 signal 的优先级最低

  5.3 尝试回收用户 pods 的资源

依据上一个步骤获得的资源 signal 来判断节点上活跃 pod 的驱逐优先级，将 pod 依据驱逐优先级排序：

比如依据内存资源评判 pod 驱逐优先级规则有：

- 依据 pod 是否超出资源请求值：没有资源使用指标的首先被驱逐。超过请求值的首先被驱逐。
- 依据 pod 的 spec.priority：依据 pod 配置的优先级排序，默认为 0。priority 越高的 pod，驱逐序列越靠后。
- 依据内存资源消费：依据 pod 消费内存超过请求值的部分排序。超过的资源绝对值越高的 pod 越优先被驱逐。

kubelet 实现了 multiSorter 的功能：依据上述顺序将活跃的 pod 排序。如果当前规则的结果是等序，才进入下个规则判断 pod 优先级。上述内存资源评判逻辑翻译过来就是，首先找出资源使用量超过请求值的 pod（包含没有指标的 pod），然后依据 pod 的 spec.priority 排序。在同 priority 的 pod 内部再依据超过的资源绝对值越高的 pod 排序。

除了 rankMemoryPressure 的逻辑，还有 rankPIDPressure，rankDiskPressure 的逻辑。

5.4 驱逐

在依据可回收资源的排序后，每次驱逐周期只会执行一次 pod 的删除。如果不是 HardEviction，还会给 MaxPodGracePeriodSeconds 的时间来让 pod 内的容器进程退出。具体的驱逐动作操作在发送事件，删除 pod 并更新 pod 的驱逐状态。

## 系统驱逐策略

上面描述的是用户态中 kubelet 通过驱逐来限制节点资源、pod 资源。在内核内存管理中，通过 OOM killer 来限制单机层面的内存使用。

### OOM killer

OOM killer（Out Of Memory killer）是一种 Linux 内核的一种内存管理机制：在系统可用内存较少的情况下，内核为保证系统还能够继续运行下去，会选择结束进程来释放内存资源。

#### 运行机制

running processes require more memory than is physically available.   内核在调用 alloc_pages() 分配内存时，如果所需要的内存超过物理内存时，通过调用 out_of_memory() 函数来选择进程释放资源。OOM killer 会检查所有运行中的进程，选择结束一个活多个进程来释放系统内存。

out_of_memory() 函数：先做一部分检查，避免通过结束进程的方式来释放内存。如果只能通过结束进程的方式来释放，那么函数会继续选择目标进程来回收。如果这个阶段也无法释放资源，kernel 最终报错异常退出。函数源码地址：[https://elixir.bootlin.com/linux/v5.17.2/source/mm/oom_kill.c#L1052，流程如下：](https://elixir.bootlin.com/linux/v5.17.2/source/mm/oom_kill.c#L1052%EF%BC%8C%E6%B5%81%E7%A8%8B%E5%A6%82%E4%B8%8B%EF%BC%9A)

1. 首先通知 oom_notify_list 链表的订阅者：依据通知链（notification chains）机制，通知注册了 oom_notify_list 的模块释放内存。如果订阅者能够处理 OOM，释放了内存则会退出 OOM killer，不执行后续操作。
2. 如果当前 task 存在 pending 的 SIGKILL，或者已经退出的时，会释放当前进程的资源。包括和 task 共享同一个内存描述符 mm_struct 的进程、线程也会被杀掉。
3. 对于 IO-less 的回收，依据 gfp_mask 判断，如果 1) 分配的是非 FS 操作类型的分配，并且 2）不是 cgroup 的内存 OOM -> 直接退出 oom-killer。
4. 检查内存分配的约束（例如 NUMA），有 CONSTRAINT_NONE, CONSTRAINT_CPUSET，CONSTRAINT_MEMORY_POLICY, CONSTRAINT_MEMCG 类型。
5. 检查 `/proc/sys/vm/panic_on_oom` 的设置，做操作。可能 panic，也可能尝试 oom_killer。如果 panic_on_oom 设置的为 2，则进程直接 panic 强制退出。
6. `/proc/sys/vm/oom_kill_allocating_task` 为 true 的时候，调用 oom_kill_process 直接 kill 掉当前想要分配内存的进程 (此进程能够被 kill 时)。
7. select_bad_process()，选择最合适的进程，调用 oom_kill_process。
8. 如果没有合适的进程，如果非 sysrq 和 memcg，则 panic 强制退出。

上述流程中有几个细节：

##### gfp_mask 约束

`/*   * The OOM killer does not compensate for IO-less reclaim.   * pagefault_out_of_memory lost its gfp context so we have to   * make sure exclude 0 mask - all other users should have at least   * ___GFP_DIRECT_RECLAIM to get here. But mem_cgroup_oom() has to   * invoke the OOM killer even if it is a GFP_NOFS allocation.   */  if (oc->gfp_mask && !(oc->gfp_mask & __GFP_FS) && !is_memcg_oom(oc))   return true;`

gfp_mask 是申请内存（get free  page）时传递的标志位。前四位表示内存域修饰符（\_\_\_GFP_DMA、\_\_\_GFP_HIGHMEM、\_\_\_GFP_DMA32、\_\_\_GFP_MOVABLE），从第 5 位开始为内存分配标志。定义：[**https://elixir.bootlin.com/linux/v5.17.2/source/include/linux/gfp.h#L81**](https://elixir.bootlin.com/linux/v5.17.2/source/include/linux/gfp.h#L81)。默认为空，从 ZONE_NORMAL 开始扫描，ZONE_NORMAL 是默认的内存申请类型。

OOM killer 不对非 IO 的回收进行补偿，所以分配的 gfp_mask 是非 FS 操作类型的分配的 OOM 会直接退出。

##### oom_constraint 约束

检查内存分配是否有限制，有几种不同的限制策略。仅适用于 NUMA 和 memcg 场景。oom_constraint 可以是：CONSTRAINT_NONE,CONSTRAINT_CPUSET,CONSTRAINT_MEMORY_POLICY,CONSTRAINT_MEMCG 类型。对于 UMA 架构而言，oom_constraint 永远都是 CONSTRAINT_NONE，表示系统并没有约束产生的 OOM。而在 NUMA 的架构下，有可能附加其他的约束导致 OOM 的情况出现。

然后调用 `check_panic_on_oom(oc)` 检查是否配置了 /proc/sys/kernel/panic_on_oom，如果有则直接触发 panic。

当走到这一步，oom killer 需要选择终止的进程，有两种选择逻辑选择合适的进程通过：

- 谁触发 OOM 就终止谁：通过 sysctl_oom_kill_allocating_task 控制，是否干掉当前申请内存的进程
- 谁最 “坏” 就制止谁：通过打分判断最 “坏” 的进程

sysctl_oom_kill_allocating_task 来自 `/proc/sys/vm/oom_kill_allocating_task`。当参数为 true 的时候，调用 oom_kill_process 直接 kill 掉当前想要分配内存的进程。

##### select_bad_process：选择最 “坏” 的进程

普通场景下通过 oom_evaluate_task 函数，评估进程分数选择需要终止的进程。如果是 memory cgroup 的情况调用 mem_cgroup_scan_tasks 来选择。先看看 oom_evaluate_task 的逻辑

- mm->flags 为 MMF_OOM_SKIP 的进程则跳过，遍历下一个进程评估
- oom_task_origin 分数最高，该标志表示 task 已经被分配大量内存并标记为 oom 的潜在原因，所以优先杀掉。
- 其他情况的进程通过 oom_badness 函数计算分数

最后分数最高的进程被终止的优先级最高。

oom_badness 函数计算的进程终止优先级**分数**由两部分组成，由下列两个参数提供。

参数：

- oom_score_adj：OOM kill score adjustment，调整值由用户打分。范围在 OOM_SCORE_ADJ_MIN（-1000） 到  OOM_SCORE_ADJ_MAX（1000）。数值越大，进程被终止的优先级越高。用户可以通过该值来保护某个进程。
- totalpages：当前可分配的内存上限值，提供系统打分的依据。

计算公式：

`/*   * The baseline for the badness score is the proportion of RAM that each   * task's rss, pagetable and swap space use.   */ points = get_mm_rss(p->mm) + get_mm_counter(p->mm, MM_SWAPENTS) +   mm_pgtables_bytes(p->mm) / PAGE_SIZE; adj *= totalpages / 1000; points += adj;`

基础分数 process_pages 由 3 部分组成：

- get_mm_rss(p->mm)：rss 部分
- get_mm_counter(p->mm, MM_SWAPENTS)：swap 占用内存
- mm_pgtables_bytes(p->mm) / PAGE_SIZE：页表占用内存

将 3 个部分相加，并结合 oom_score_adj：将归一化后的 adj 和 points 求和，作为当前进程的分数。

所以进程得分 `points=process_pages + oom_score_adj*totalpages/1000`

之前老的内核版本还会有一些复杂的计算逻辑考虑，比如对于特权进程的处理。如果是 root 权限的进程，有 3% 的内存使用特权。points=process_pages_0.97 + oom_score_adj_totalpages/1000。v4.17 移除，使得计算逻辑更加简洁和可预测。

`/*   * Root processes get 3% bonus, just like the __vm_enough_memory()   * implementation used by LSMs.   */  if (has_capability_noaudit(p, CAP_SYS_ADMIN))   points -= (points * 3) / 100;`

mem_cgroup_scan_tasks：memory cgroup cgroup 的处理会需要遍历 cgroup 的层次结构，调用 oom_evaluate_task 计算 task 的分数。回收父进程的内存也会回收子进程的内存。

###### oom_kill_process

接下来进入终止进程的逻辑，oom_kill_process 函数在终止进程之前会先检查，task 是否已经退出，占用的内存会被释放，防止重复处理；获取 memory cgroup 消息，判断是否需要删除 cgroup 下所有的 tasks。然后是 dump 信息，将 OOM 的原因打印出来，保留 OOM 的线索。

之后在 \_\_oom_kill_process 函数内调用 put_task_struct 释放内核栈，释放系统资源。唤醒 oom_reaper 内核线程收割 wake_oom_reaper(victim)。

oom_reaper 会在有清理任务之前一直保持休眠。wake_oom_reaper 会将任务压入 oom_reaper_list 链表，oom_reaper 通过 oom_reaper_list 链表判断需要调用 oom_reap_task_mm 清理地址空间。清理时会遍历 vma，跳过 VM_LOCKED|VM_HUGETLB|VM_PFNMAP 的 VMA 区域。具体的释放操作通过 unmap_page_range 完成：

\` for (vma = mm->mmap ; vma; vma = vma->vm_next) {
  if (!can_madv_lru_vma(vma))
   continue;

/\_
   * Only anonymous pages have a good chance to be dropped
   * without additional steps which we cannot afford as we
   * are OOM already.
   \_
   * We do not even care about fs backed pages because all
   * which are reclaimable have already been reclaimed and
   * we do not want to block exit*mmap by keeping mm ref
   * count elevated without a good reason.
   \_/
  if (vma_is_anonymous(vma) || !(vma->vm_flags & VM_SHARED)) {
   struct mmu_notifier_range range;
   struct mmu_gather tlb;

mmu_notifier_range_init(\&range, MMU_NOTIFY_UNMAP, 0,
      vma, mm, vma->vm_start,
      vma->vm_end);
   tlb_gather_mmu(\&tlb, mm);
   if (mmu_notifier_invalidate_range_start_nonblock(\&range)) {
    tlb_finish_mmu(\&tlb);
    ret = false;
    continue;
   }
   unmap_page_range(\&tlb, vma, range.start, range.end, NULL);
   mmu_notifier_invalidate_range_end(\&range);
   tlb_finish_mmu(\&tlb);
  }
 }

\`

<https://elixir.bootlin.com/linux/v5.17.2/source/mm/oom_kill.c#L528>

##### 控制 oom killer 的行为

上述有提及几个文件参数来控制控制 oom killer 的行为：

1. /proc/sys/vm/panic_on_oom，当出现 oom 时，该值设定允许或者禁止 kernel panic（默认为 0）

- 0: 发生 oom 时，内核会选择调用 oom-killer 来选择进程删除
- 1: 发生 oom 时，内核通常情况会直接 panic，除了特定条件：通过 mempolicy/cpusets 限制使用的进程则会被 oom-killer 删除时，不会 panic
- 2: 发生 oom 时，内核无条件直接 panic

2. /proc/sys/vm/oom_kill_allocating_task，可以取值为 0 或者非 0（默认为 0），0 代表发送 oom 时，进行遍历任务链表，选择一个进程去杀死，而非 0 代表，发送 oom 时，直接 kill 掉引起 oom 的进程，并不会去遍历任务链表。
3. /proc/sys/vm/oom_dump_tasks：可以取值为 0 或者非 0（默认为 1），表示是否在发送 oom killer 时，打印 task 的相关信息。
4. /proc//oom_score_adj：配置进程的评分调整分，通过该在值来保护某个进程不被杀死或者每次都杀某个进程。其取值范围为 - 1000 到 1000 。
5. /proc/sys/vm/overcommit_memory：控制内存超售，oom-killer 功能，默认为 0

- 0：**启发式策略** ，比较严重的 Overcommit 将不能得逞，比如你突然申请了 128TB 的内存。而轻微的 overcommit 将被允许。另外，root 能 Overcommit 的值比普通用户要稍微多。默认
- 1：**永远允许 overcommit** ，这种策略适合那些不能承受内存分配失败的应用，比如某些科学计算应用。
- 2：**永远禁止 overcommit** ，在这个情况下，系统所能分配的内存不会超过 **swap+RAM\*系数** （/proc/sys/vm/overcmmit_ratio，默认 50%，你可以调整），如果这么多资源已经用光，那么后面任何尝试申请内存的行为都会返回错误，这通常意味着此时没法运行任何新程序。

Memory cgroup 子系统的控制：

1. memory.use_hierarchy：指定 cgroup 层次结构。（default 为 0）

- 0：父进程不从子进程回收内存
- 1：会从超出内存限制（memory limit）的子进程中回收

2. memory.oom_control：oom 控制，（默认为 0：每个 cgroup 内存子系统）

- 0：当进程消费更多的内存时会被 oom_killer 杀掉
- 1：关闭 oom_killer，当 task 尝试使用更多的内存时，会卡住直到内存充足。
- 读文件时，描述 oom 状态：oom_kill_disable（是否开启）、under_oom（是否处于 oom 状态）

## 用户空间的 oom killer

最后再简单介绍一个用户空间的 oom killer：[https://github.com/facebookincubator/oomd。oomd](https://github.com/facebookincubator/oomd%E3%80%82oomd) 的目标是在用户空间，解决内存资源使用的问题。

### 运行机制

- 使用 PSI、cgroupv2 来监控系统上的内存使用情况，oomd 在内核的 oom_killer 处理之前，先进行内存资源的释放。
- 监控系统和 cgroup 的内存压力。

并且配置上可以做到如此驱逐策略：

- 当 workload 有内存压力 / 系统有内存压力时，通过内存大小或增长率选择一个 memory hog（资源大户）删除。
- 当系统有内存压力时，通过内存大小或增长率选择一个 memory hog 删除。
- 当系统有 swap 压力时，选择使用 swap 最多的 cgroup 来删除。

可以看到，oomd 充当了 kubelet 的功能，是单机上 oom 管理的 agent。

## 总结

可以看到用户空间和内核空间的驱逐策略的不同。用户空间通过监控系统资源来触发驱逐流程，内核空间通过分配内存时触发驱逐流程。因为用户空间的驱逐需要在内核驱逐之前来

除了进程驱逐手段，还有其他手段来做到资源保障和稳定性，比如资源抑制和回收。通过 cgroup v2 的 Memory Qos 的能力

- 当整机内存出现压力时，保障 container 的内存分配性能，降低其内存分配延迟
- 对过度申请内存的 container 进行抑制和快速回收，降低整机内存的使用压力
- 对整机保留内存进行保护

#### 参考

- Memory Resource Controller: [**https://www.kernel.org/doc/Documentation/cgroup-v1/memory.txt**](https://www.kernel.org/doc/Documentation/cgroup-v1/memory.txt)
- liux 下 oom killer: [**https://www.mo4tech.com/oom-killer-mechanism-for-the-linux-kernel.html**](https://www.mo4tech.com/oom-killer-mechanism-for-the-linux-kernel.html)
- 内存分配掩码（gfp_mask）：[**https://blog.csdn.net/farmwang/article/details/66975128**](https://blog.csdn.net/farmwang/article/details/66975128)
- oom-killer 日志分析：[**https://bhsc881114.github.io/2018/06/24/oom-killer%E7%90%86%E8%A7%A3%E5%92%8C%E6%97%A5%E5%BF%97%E5%88%86%E6%9E%90-%E6%97%A5%E5%BF%97%E5%88%86%E6%9E%90/**](https://bhsc881114.github.io/2018/06/24/oom-killer%E7%90%86%E8%A7%A3%E5%92%8C%E6%97%A5%E5%BF%97%E5%88%86%E6%9E%90-%E6%97%A5%E5%BF%97%E5%88%86%E6%9E%90/)
- memory managemnent：[**https://learning-kernel.readthedocs.io/en/latest/mem-management.html**](https://learning-kernel.readthedocs.io/en/latest/mem-management.html)
- Linux 内存管理 (21)OOM：[**https://www.cnblogs.com/arnoldlu/p/8567559.html**](https://www.cnblogs.com/arnoldlu/p/8567559.html)
- Liux OOM 的参数：[**http://www.wowotech.net/memory_management/oom.html**](http://www.wowotech.net/memory_management/oom.html)

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/722c5b2e-6e47-406c-b75b-d475e0296a7b/640)

<https://mp.weixin.qq.com/s/ehECtQiXSHLpCrH5vuBX_w>
