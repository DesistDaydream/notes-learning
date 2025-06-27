---
title: PodWorker 模块
---

# 概述

PodWorkers 子模块主要的作用就是处理针对每一个的 Pod 的更新事件，比如 Pod 的创建，删除，更新。而 podWorkers 采取的基本思路是：为每一个 Pod 都单独创建一个 goroutine 和更新事件的 channel，goroutine 会阻塞式的等待 channel 中的事件，并且对获取的事件进行处理。而 podWorkers 对象自身则主要负责对更新事件进行下发。

# 准备运行 Pod

## podWorkers.UpdatePod() - 更新事件的 channel

updatePod 将配置更改或终止状态传递到 POD。 POD 可以是可变的，终止或终止，并且如果在 APIServer 上删除，则将转换为终止，它被发现具有终端阶段（成功或失败），或者如果它被 kubelet 驱逐。

为每一个 Pod 都单独创建一个 goroutine 和更新事件的 channel，goroutine 会阻塞式的等待 channel 中的事件，并且对获取的事件进行处理。而 podWorkers 对象自身则主要负责对更新事件进行下发。

源码：`pkg/kubelet/pod_workers.go`

```go
func (p *podWorkers) UpdatePod(options UpdatePodOptions) {
    // 处理当 Pod 是孤儿(无配置)并且我们仅通过仅运行生命周期的终止部分来获得运行时状态
	pod := options.Pod
	var isRuntimePod bool
	if options.RunningPod != nil {
		if options.Pod == nil {
			pod = options.RunningPod.ToAPIPod()
			if options.UpdateType != kubetypes.SyncPodKill {
				klog.InfoS("Pod update is ignored, runtime pods can only be killed", "pod", klog.KObj(pod), "podUID", pod.UID)
				return
			}
			options.Pod = pod
			isRuntimePod = true
		} else {
			options.RunningPod = nil
			klog.InfoS("Pod update included RunningPod which is only valid when Pod is not specified", "pod", klog.KObj(options.Pod), "podUID", options.Pod.UID)
		}
	}
	uid := pod.UID

	p.podLock.Lock()
	defer p.podLock.Unlock()

	// decide what to do with this pod - we are either setting it up, tearing it down, or ignoring it
	now := time.Now()
	status, ok := p.podSyncStatuses[uid]
	if !ok {
		klog.V(4).InfoS("Pod is being synced for the first time", "pod", klog.KObj(pod), "podUID", pod.UID)
		status = &podSyncStatus{
			syncedAt: now,
			fullname: kubecontainer.GetPodFullName(pod),
		}
		// if this pod is being synced for the first time, we need to make sure it is an active pod
		if !isRuntimePod && (pod.Status.Phase == v1.PodFailed || pod.Status.Phase == v1.PodSucceeded) {
			// check to see if the pod is not running and the pod is terminal.
			// If this succeeds then record in the podWorker that it is terminated.
			if statusCache, err := p.podCache.Get(pod.UID); err == nil {
				if isPodStatusCacheTerminal(statusCache) {
					status = &podSyncStatus{
						terminatedAt:       now,
						terminatingAt:      now,
						syncedAt:           now,
						startedTerminating: true,
						finished:           true,
						fullname:           kubecontainer.GetPodFullName(pod),
					}
				}
			}
		}
		p.podSyncStatuses[uid] = status
	}

	// if an update is received that implies the pod should be running, but we are already terminating a pod by
	// that UID, assume that two pods with the same UID were created in close temporal proximity (usually static
	// pod but it's possible for an apiserver to extremely rarely do something similar) - flag the sync status
	// to indicate that after the pod terminates it should be reset to "not running" to allow a subsequent add/update
	// to start the pod worker again
	if status.IsTerminationRequested() {
		if options.UpdateType == kubetypes.SyncPodCreate {
			status.restartRequested = true
			klog.V(4).InfoS("Pod is terminating but has been requested to restart with same UID, will be reconciled later", "pod", klog.KObj(pod), "podUID", pod.UID)
			return
		}
	}

	// once a pod is terminated by UID, it cannot reenter the pod worker (until the UID is purged by housekeeping)
	if status.IsFinished() {
		klog.V(4).InfoS("Pod is finished processing, no further updates", "pod", klog.KObj(pod), "podUID", pod.UID)
		return
	}

	// check for a transition to terminating
	var becameTerminating bool
	if !status.IsTerminationRequested() {
		switch {
		case isRuntimePod:
			klog.V(4).InfoS("Pod is orphaned and must be torn down", "pod", klog.KObj(pod), "podUID", pod.UID)
			status.deleted = true
			status.terminatingAt = now
			becameTerminating = true
		case pod.DeletionTimestamp != nil:
			klog.V(4).InfoS("Pod is marked for graceful deletion, begin teardown", "pod", klog.KObj(pod), "podUID", pod.UID)
			status.deleted = true
			status.terminatingAt = now
			becameTerminating = true
		case pod.Status.Phase == v1.PodFailed, pod.Status.Phase == v1.PodSucceeded:
			klog.V(4).InfoS("Pod is in a terminal phase (success/failed), begin teardown", "pod", klog.KObj(pod), "podUID", pod.UID)
			status.terminatingAt = now
			becameTerminating = true
		case options.UpdateType == kubetypes.SyncPodKill:
			if options.KillPodOptions != nil && options.KillPodOptions.Evict {
				klog.V(4).InfoS("Pod is being evicted by the kubelet, begin teardown", "pod", klog.KObj(pod), "podUID", pod.UID)
				status.evicted = true
			} else {
				klog.V(4).InfoS("Pod is being removed by the kubelet, begin teardown", "pod", klog.KObj(pod), "podUID", pod.UID)
			}
			status.terminatingAt = now
			becameTerminating = true
		}
	}

	// once a pod is terminating, all updates are kills and the grace period can only decrease
	var workType PodWorkType
	var wasGracePeriodShortened bool
	switch {
	case status.IsTerminated():
		// A terminated pod may still be waiting for cleanup - if we receive a runtime pod kill request
		// due to housekeeping seeing an older cached version of the runtime pod simply ignore it until
		// after the pod worker completes.
		if isRuntimePod {
			klog.V(3).InfoS("Pod is waiting for termination, ignoring runtime-only kill until after pod worker is fully terminated", "pod", klog.KObj(pod), "podUID", pod.UID)
			return
		}

		workType = TerminatedPodWork

		if options.KillPodOptions != nil {
			if ch := options.KillPodOptions.CompletedCh; ch != nil {
				close(ch)
			}
		}
		options.KillPodOptions = nil

	case status.IsTerminationRequested():
		workType = TerminatingPodWork
		if options.KillPodOptions == nil {
			options.KillPodOptions = &KillPodOptions{}
		}

		if ch := options.KillPodOptions.CompletedCh; ch != nil {
			status.notifyPostTerminating = append(status.notifyPostTerminating, ch)
		}
		if fn := options.KillPodOptions.PodStatusFunc; fn != nil {
			status.statusPostTerminating = append(status.statusPostTerminating, fn)
		}

		gracePeriod, gracePeriodShortened := calculateEffectiveGracePeriod(status, pod, options.KillPodOptions)

		wasGracePeriodShortened = gracePeriodShortened
		status.gracePeriod = gracePeriod
		// always set the grace period for syncTerminatingPod so we don't have to recalculate,
		// will never be zero.
		options.KillPodOptions.PodTerminationGracePeriodSecondsOverride = &gracePeriod

	default:
		workType = SyncPodWork

		// KillPodOptions is not valid for sync actions outside of the terminating phase
		if options.KillPodOptions != nil {
			if ch := options.KillPodOptions.CompletedCh; ch != nil {
				close(ch)
			}
			options.KillPodOptions = nil
		}
	}

	// the desired work we want to be performing
	work := podWork{
		WorkType: workType,
		Options:  options,
	}

	// 如果 pod worker 协程不存在则启动它
	podUpdates, exists := p.podUpdates[uid]
	if !exists {
		// 创建 channel
        // 我们需要在这里有一个缓冲区，因为将更新放入通道的 checkForUpdates() 方法是从使用通道的同一个 goroutine 调用的。但是，可以保证在这种情况下通道是空的，因此大小为 1 的缓冲区就足够了。
		podUpdates = make(chan podWork, 1)
		p.podUpdates[uid] = podUpdates

		// 确保静态 pod 按照 UpdatePod 接收它们的顺序启动
		if kubetypes.IsStaticPod(pod) {
			p.waitingToStartStaticPodsByFullname[status.fullname] =
				append(p.waitingToStartStaticPodsByFullname[status.fullname], uid)
		}

		// 允许测试 pod 更新通道中的延迟
		var outCh <-chan podWork
		if p.workerChannelFn != nil {
			outCh = p.workerChannelFn(uid, podUpdates)
		} else {
			outCh = podUpdates
		}

        // 启动 goroutine
        // 创建一个新的 Pod Worker 意味着这是一个新的 POD，或者 kubelet 刚刚重新启动。
        // 在任何一种情况下，Kubelet 都愿意相信第一个 POD Worker 同步的 POD 的状态。请参阅 Syncpod 中的相应评论。
		go func() {
			defer runtime.HandleCrash()
			p.managePodLoop(outCh)
		}()
	}

	// 如果没有运行，则向 pod worker 请求
	if !status.IsWorking() {
		status.working = true
		podUpdates <- work
		return
	}

	// 捕获请求的更新与pod worker观察到更新之间的最大延迟
	if undelivered, ok := p.lastUndeliveredWorkUpdate[pod.UID]; ok {
		// track the max latency between when a config change is requested and when it is realized
		// NOTE: this undercounts the latency when multiple requests are queued, but captures max latency
		if !undelivered.Options.StartTime.IsZero() && undelivered.Options.StartTime.Before(work.Options.StartTime) {
			work.Options.StartTime = undelivered.Options.StartTime
		}
	}

	// 始终同步最新数据
	p.lastUndeliveredWorkUpdate[pod.UID] = work

	if (becameTerminating || wasGracePeriodShortened) && status.cancelFn != nil {
		klog.V(3).InfoS("Cancelling current pod sync", "pod", klog.KObj(pod), "podUID", pod.UID, "updateType", work.WorkType)
		status.cancelFn()
		return
	}
}
```

## podWorkers.managePodLoop() - 调用 podWorkers.syncPodFn() 方法同步 Pod

managePodLoop 调用 `podWorkers.syncPodFn()` 方法去同步 pod。在完成这次 sync 动作之后，会调用 wrapUp 函数，这个函数将会做几件事情:

- 将这个 pod 信息插入 kubelet 的 workQueue 队列中，等待下一次周期性的对这个 pod 的状态进行 sync
- 将在这次 sync 期间堆积的没有能够来得及处理的最近一次 update 操作加入 goroutine 的事件 channel 中，立即处理。

源码：`pkg/kubelet/pod_workers.go`

```go
func (p *podWorkers) managePodLoop(podUpdates <-chan podWork) {
	for update := range podUpdates {
		err := func() error {
			// 采取适当的行动（UpdatePod阻止了非法阶段）
			switch {
			case update.WorkType == TerminatedPodWork:
			case update.WorkType == TerminatingPodWork:
			default:
                // 这里的 podWorkers.syncPodFn() 实际上是 kubelet.SyncPod() 方法
				err = p.syncPodFn(ctx, update.Options.UpdateType, pod, update.Options.MirrorPod, status)
			}
		}()
    }
}
```

### 引用说明

这里的 `p.syncPodFn()` 引用的是 `kubelet.SyncPod()` 方法。来源如下：

- 源码：`pkg/kubelet/kubelet.go`

```go
func NewMainKubelet(......) (*Kubelet, error) {
	klet := &Kubelet{
		......
	}
	klet.podWorkers = newPodWorkers(klet.syncPod,......)
}
```

- 源码：`pkg/kubelet/pod_workers.go`

```go
func newPodWorkers(syncPodFn syncPodFnType,......) PodWorkers {
	return &podWorkers{
		syncPodFn: syncPodFn,
        ......
	}
}
```

## kubelet.syncPod() - 完成创建容器前的准备工作

在这个方法中，主要完成以下几件事情：

- 如果是删除 pod，立即执行并返回
- 同步 podStatus 到 kubelet.statusManager
- 检查 pod 是否能运行在本节点，主要是权限检查（是否能使用主机网络模式，是否可以以 privileged 权限运行等）。如果没有权限，就删除本地旧的 pod 并返回错误信息
- 创建 containerManagar 对象，并且创建 pod level cgroup，更新 Qos level cgroup
- 如果是 static Pod，就创建或者更新对应的 mirrorPod
- 创建 pod 的数据目录，存放 volume 和 plugin 信息,如果定义了 pv，等待所有的 volume mount 完成（volumeManager 会在后台做这些事情）,如果有 image secrets，去 apiserver 获取对应的 secrets 数据
- 然后调用 kubelet.volumeManager 组件，等待它将 pod 所需要的所有外挂的 volume 都准备好。
- 调用 Runtime() 接口中的 SyncPod 方法，去实现真正的容器创建逻辑

这里所有的事情都和具体的容器没有关系，可以看到该方法是创建 pod 实体（即容器）之前需要完成的准备工作。

源码：`pkg/kubelet/kubelet.go`

```go
func (kl *Kubelet) syncPod(ctx context.Context, updateType kubetypes.SyncPodType, pod, mirrorPod *v1.Pod, podStatus *kubecontainer.PodStatus) error {
	// 主要工作流的延迟测量是相对于 kubelet 第一次发现 Pod 的时间
	var firstSeenTime time.Time
	if firstSeenTimeStr, ok := pod.Annotations[kubetypes.ConfigFirstSeenAnnotationKey]; ok {
		firstSeenTime = kubetypes.ConvertToTimestamp(firstSeenTimeStr).Get()
	}

	// 如果创建，记录 Pod Worker 启动延迟
	// TODO: make pod workers record their own latencies
	if updateType == kubetypes.SyncPodCreate {
		if !firstSeenTime.IsZero() {
			// 这是我们第一次同步 pod。如果设置了 firstSeenTime，则记录自 kubelet 第一次看到 pod 以来的延迟。
			metrics.PodWorkerStartDuration.Observe(metrics.SinceInSeconds(firstSeenTime))
		} else {
			klog.V(3).InfoS("First seen time not recorded for pod",
				"podUID", pod.UID,
				"pod", klog.KObj(pod))
		}
	}

	// Generate final API pod status with pod and status manager status
	apiPodStatus := kl.generateAPIPodStatus(pod, podStatus)
	// The pod IP may be changed in generateAPIPodStatus if the pod is using host network. (See #24576)
	// TODO(random-liu): After writing pod spec into container labels, check whether pod is using host network, and
	// set pod IP to hostIP directly in runtime.GetPodStatus
	podStatus.IPs = make([]string, 0, len(apiPodStatus.PodIPs))
	for _, ipInfo := range apiPodStatus.PodIPs {
		podStatus.IPs = append(podStatus.IPs, ipInfo.IP)
	}

	if len(podStatus.IPs) == 0 && len(apiPodStatus.PodIP) > 0 {
		podStatus.IPs = []string{apiPodStatus.PodIP}
	}

    // 检查 Pod 是否可以运行在本节点。如果 Pod 不应该运行，将 Pod 的容器 stop，这与 termination 不同(我们希望 stop Pod，但如果软准入机制允许稍后重启它)
    // 适当设置状态和阶段
	runnable := kl.canRunPod(pod)
	if !runnable.Admit {
	}

	// 如果设置了 firstSeenTime，记录自 kubelet 首次看到 Pod 以来 Pod 运行所需的时间。
	existingStatus, ok := kl.statusManager.GetPodStatus(pod.UID)
	if !ok || existingStatus.Phase == v1.PodPending && apiPodStatus.Phase == v1.PodRunning &&
		!firstSeenTime.IsZero() {
		metrics.PodStartDuration.Observe(metrics.SinceInSeconds(firstSeenTime))
	}

    // 更新 Pod 状态
	kl.statusManager.SetPodStatus(pod, apiPodStatus)

	// 必须停止不可运行的 Pod，并向 PodWorker 返回一个错误类型
	if !runnable.Admit {
	}

	// 加载网络插件，如果网络插件没有准备好，只有在 Pod 使用宿主机的网络时才启动它
	if err := kl.runtimeState.networkErrors(); err != nil && !kubecontainer.IsHostNetworkPod(pod) {
	}

	// 确保 kubelet 知道 Pod 使用的 secrets 和 configmaps 资源
	if !kl.podWorkers.IsPodTerminationRequested(pod.UID) {
	}

	// 为 Pod 创建 Cgroups，并在启用 cgroups-per-qos 标志的情况下对其应用资源参数。
	pcm := kl.containerManager.NewPodContainerManager()

	// 为静态 Pod 创建 Mirror Pod
	if kubetypes.IsStaticPod(pod) {
	}

	// 为 Pod 创建数据目录
	if err := kl.makePodDataDirs(pod); err != nil {
	}

    // 挂载 Volume
	// Volume 管理器不会为 terminating 状态的 Pod 挂载卷
	// TODO: 一旦添加上下文取消，可以删除此检查
	if !kl.podWorkers.IsPodTerminationRequested(pod.UID) {
		// 等待卷 attach/mount
		if err := kl.volumeManager.WaitForAttachAndMount(pod); err != nil {
		}
	}

	// 获取 Pod 的 secret 信息
	pullSecrets := kl.getPullSecretsForPod(pod)

    // 调用 Runtime 接口中的 SyncPod() 方法以开始创建容器
    // 这里的 kl.containerRuntime.SyncPod() 实际上是 kubeGenericRuntimeManager.SyncPod() 方法
	result := kl.containerRuntime.SyncPod(pod, podStatus, pullSecrets, kl.backOff)
	kl.reasonCache.Update(pod.UID, result)
}
```

### 引用说明

这里的 `kl.containerRuntime.SyncPod()` 引用的是 `kubeGenericRuntimeManager.SyncPod()` 方法，来源如下：

源码：`pkg/kubelet/kuberuntime`-`NewMainKubelet()`

```go
func NewMainKubelet() (*Kubelet, error) {
	klet := &Kubelet{......}
	runtime, err := kuberuntime.NewKubeGenericRuntimeManager(......)
	klet.containerRuntime = runtime
}

func (kl *Kubelet) syncPod(ctx context.Context, updateType kubetypes.SyncPodType, pod, mirrorPod *v1.Pod, podStatus *kubecontainer.PodStatus) error {
	result := kl.containerRuntime.SyncPod(pod, podStatus, pullSecrets, kl.backOff)
}
```

# 开始运行 Pod(CRI 在这里)

KubeRuntimeManager(pkg/kubelet/kuberuntime) 子模块的 SyncPod() 方法是真正完成 Pod 内容器实体的创建。

`kubeGenericRuntimeManager.runtimeService` **就是 CRI**，都是由第三方对接的，具体时间逻辑也在第三方，比如 Containerd，Docker 等。该结构体中的很多方法都调用了 runtimeService 接口中的方法，以控制 Pod 和 容器。

## kubeGenericRuntimeManager.SyncPod() - 创建容器

syncPod 主要执行以下几个操作：

1. 计算 Sandbox 和 Container 是否发生变化
2. 必要时 kill 调 Pod Sandbox
3. kill 调不应该运行的所有容器
4. 必要时创建 Sandbox 容器
5. 创建临时容器
6. 创建初始化容器
7. 创建业务容器
8. 在创建容器中调用 `kubeGenericRuntimeManager.startContainer()` 启动容器

initContainers 可以有多个，多个 container 严格按照顺序启动，只有当前一个 container 退出了以后，才开始启动下一个 container。

源码：`pkg/kubelet/kuberuntime/kuberuntime_manager.go`- `containerRuntime.SyncPod()`

```go
func (m *kubeGenericRuntimeManager) SyncPod(pod *v1.Pod, podStatus *kubecontainer.PodStatus, pullSecrets []v1.Secret, backOff *flowcontrol.Backoff) (result kubecontainer.PodSyncResult) {
	// Step 1：计算 sandbox 和 container 是否发生变化
	podContainerChanges := m.computePodActions(pod, podStatus)
	klog.V(3).InfoS("computePodActions got for pod", "podActions", podContainerChanges, "pod", klog.KObj(pod))
	if podContainerChanges.CreateSandbox {
		ref, err := ref.GetReference(legacyscheme.Scheme, pod)
		if err != nil {
			klog.ErrorS(err, "Couldn't make a ref to pod", "pod", klog.KObj(pod))
		}
		if podContainerChanges.SandboxID != "" {
			m.recorder.Eventf(ref, v1.EventTypeNormal, events.SandboxChanged, "Pod sandbox changed, it will be killed and re-created.")
		} else {
			klog.V(4).InfoS("SyncPod received new pod, will create a sandbox for it", "pod", klog.KObj(pod))
		}
	}

	// Step 2：kill 掉 sandbox 已经改变的 Pod
	if podContainerChanges.KillPod {
		if podContainerChanges.CreateSandbox {
			klog.V(4).InfoS("Stopping PodSandbox for pod, will start new one", "pod", klog.KObj(pod))
		} else {
			klog.V(4).InfoS("Stopping PodSandbox for pod, because all other containers are dead", "pod", klog.KObj(pod))
		}

		killResult := m.killPodWithSyncResult(pod, kubecontainer.ConvertPodStatusToRunningPod(m.runtimeName, podStatus), nil)
		result.AddPodSyncResult(killResult)
		if killResult.Error() != nil {
			klog.ErrorS(killResult.Error(), "killPodWithSyncResult failed")
			return
		}

		if podContainerChanges.CreateSandbox {
			m.purgeInitContainers(pod, podStatus)
		}
	} else {
		// Step 3：kill 掉非 running 状态的容器
		for containerID, containerInfo := range podContainerChanges.ContainersToKill {
			klog.V(3).InfoS("Killing unwanted container for pod", "containerName", containerInfo.name, "containerID", containerID, "pod", klog.KObj(pod))
			killContainerResult := kubecontainer.NewSyncResult(kubecontainer.KillContainer, containerInfo.name)
			result.AddSyncResult(killContainerResult)
			if err := m.killContainer(pod, containerID, containerInfo.name, containerInfo.message, containerInfo.reason, nil); err != nil {
				killContainerResult.Fail(kubecontainer.ErrKillContainer, err.Error())
				klog.ErrorS(err, "killContainer for pod failed", "containerName", containerInfo.name, "containerID", containerID, "pod", klog.KObj(pod))
				return
			}
		}
	}

	// Step 4：如果必要，为 Pod 创建 sandbox
	podSandboxID := podContainerChanges.SandboxID
	if podContainerChanges.CreateSandbox {
		// ConvertPodSysctlsVariableToDotsSeparator converts sysctl variable
		// in the Pod.Spec.SecurityContext.Sysctls slice into a dot as a separator.
		// runc uses the dot as the separator to verify whether the sysctl variable
		// is correct in a separate namespace, so when using the slash as the sysctl
		// variable separator, runc returns an error: "sysctl is not in a separate kernel namespace"
		// and the podSandBox cannot be successfully created. Therefore, before calling runc,
		// we need to convert the sysctl variable, the dot is used as a separator to separate the kernel namespace.
		// When runc supports slash as sysctl separator, this function can no longer be used.
		sysctl.ConvertPodSysctlsVariableToDotsSeparator(pod.Spec.SecurityContext)

		podSandboxID, msg, err = m.createPodSandbox(pod, podContainerChanges.Attempt)

		resp, err := m.runtimeService.PodSandboxStatus(podSandboxID, false)

        // 如果 pod 网络是 host 模式，容器也相同；其他情况下，容器会使用 None 网络模式，让 kubelet 的网络插件自己进行网络配置
		if !kubecontainer.IsHostNetworkPod(pod) {
			podIPs = m.determinePodSandboxIPs(pod.Namespace, pod.Name, resp.GetStatus())
			klog.V(4).InfoS("Determined the ip for pod after sandbox changed", "IPs", podIPs, "pod", klog.KObj(pod))
		}
	}

    // 为容器获取 Sandbox 配置(如：元数据、集群DNS 、容器的端口映射 等等)
	configPodSandboxResult := kubecontainer.NewSyncResult(kubecontainer.ConfigPodSandbox, podSandboxID)
	result.AddSyncResult(configPodSandboxResult)
	podSandboxConfig, err := m.generatePodSandboxConfig(pod, podContainerChanges.Attempt)

    // 用于启动容器的行为，适用于任何类型的容器，容器类型包括：container(容器)、init_container(初始化容器)、ephemeral_container(临时容器)
    // 上述三种对容器的分类描述，在 日志消息 与 监控指标的标签 中会出现，用来定位容器。
    // 下面代码中启动容器时，都会调用 start，也就是 `func(typeName, metricLabel string, spec *startSpec) error{}` 函数
    // 启动容器的核心是 m.startContainer() 方法
	start := func(typeName, metricLabel string, spec *startSpec) error {
        // Step 最终：调用 m.startContainer() 启动容器
        // 注意（Aramase）Podips填充单堆栈和双堆栈集群。只发送Podips。
		m.startContainer(podSandboxID, podSandboxConfig, spec, pod, podStatus, pullSecrets, podIP, podIPs)
	}

    // Step 5：启动 ephemeral_container(临时容器)，调用上面定义的 start。
    for _, idx := range podContainerChanges.EphemeralContainersToStart {
        start("ephemeral container", metrics.EphemeralContainer, ephemeralContainerStartSpec(&pod.Spec.EphemeralContainers[idx]))
    }

    // Step 6: 启动 init_container(初始化容器)，调用上面定义的 start。
	if container := podContainerChanges.NextInitContainerToStart; container != nil {
		start("init container", metrics.InitContainer, containerStartSpec(container))
	}

    // Step 7：启动 container(容器)。调用上面定义的 start。
	for _, idx := range podContainerChanges.ContainersToStart {
		start("container", metrics.Container, containerStartSpec(&pod.Spec.Containers[idx]))
	}

	return
}
```

## kubeGenericRuntimeManager.startContainer() - 启动容器

最终由 `kubeGenericRuntimeManager.startContainer()` 完成容器的启动，其主要有以下几个步骤：

- 拉取镜像
- 创建容器
- 启动容器
- 运行 post start lifecycle hooks

源码：`pkg/kubelet/kuberuntime/kuberuntime_container.go`

```go
func (m *kubeGenericRuntimeManager) startContainer(podSandboxID string, podSandboxConfig *runtimeapi.PodSandboxConfig, spec *startSpec, pod *v1.Pod, podStatus *kubecontainer.PodStatus, pullSecrets []v1.Secret, podIP string, podIPs []string) (string, error) {
	container := spec.container

	// Step 1：拉取镜像
	imageRef, msg, err := m.imagePuller.EnsureImageExists(pod, container, pullSecrets, podSandboxConfig)
	if err != nil {
		s, _ := grpcstatus.FromError(err)
		m.recordContainerEvent(pod, container, "", v1.EventTypeWarning, events.FailedToCreateContainer, "Error: %v", s.Message())
		return msg, err
	}

	// Step 1：创建容器
	// 对于一个新的容器，RestartCount 变量的值应该为 0
	restartCount := 0
	containerStatus := podStatus.FindContainerStatusByName(container.Name)

	target, err := spec.getTargetID(podStatus)
	m.generateContainerConfig(container, pod, restartCount, podIP, imageRef, podIPs, target)
	m.internalLifecycle.PreCreateContainer(pod, container, containerConfig)
	m.runtimeService.CreateContainer(podSandboxID, containerConfig, podSandboxConfig)
	m.internalLifecycle.PreStartContainer(pod, container, containerID)

	// 3、启动容器
	err = m.runtimeService.StartContainer(containerID)

	// Symlink container logs to the legacy container log location for cluster logging
	// support.
	// TODO(random-liu): Remove this after cluster logging supports CRI container log path.
	containerMeta := containerConfig.GetMetadata()
	sandboxMeta := podSandboxConfig.GetMetadata()
	legacySymlink := legacyLogSymlink(containerID, containerMeta.Name, sandboxMeta.Name,
		sandboxMeta.Namespace)
	containerLog := filepath.Join(podSandboxConfig.LogDirectory, containerConfig.LogPath)
	// only create legacy symlink if containerLog path exists (or the error is not IsNotExist).
	// Because if containerLog path does not exist, only dangling legacySymlink is created.
	// This dangling legacySymlink is later removed by container gc, so it does not make sense
	// to create it in the first place. it happens when journald logging driver is used with docker.
	if _, err := m.osInterface.Stat(containerLog); !os.IsNotExist(err) {
		if err := m.osInterface.Symlink(containerLog, legacySymlink); err != nil {
			klog.ErrorS(err, "Failed to create legacy symbolic link", "path", legacySymlink,
				"containerID", containerID, "containerLogPath", containerLog)
		}
	}

	// 4、执行启动后的 Hook，就是一些启动后的检查，如果检查不通过，容器将会处于异常状态，并根据策略决定是否重启
	if container.Lifecycle != nil && container.Lifecycle.PostStart != nil {
		kubeContainerID := kubecontainer.ContainerID{
			Type: m.runtimeName,
			ID:   containerID,
		}
        // runner.Run 这个方法的主要作用就是在业务容器起来的时候，
        // 首先会执行一个 container hook(PostStart 和 PreStop),做一些预处理工作。
        // 只有 container hook 执行成功才会运行具体的业务服务，否则容器异常。
		msg, handlerErr := m.runner.Run(kubeContainerID, pod, container, container.Lifecycle.PostStart)
		if handlerErr != nil {
			klog.ErrorS(handlerErr, "Failed to execute PostStartHook", "pod", klog.KObj(pod),
				"podUID", pod.UID, "containerName", container.Name, "containerID", kubeContainerID.String())
			m.recordContainerEvent(pod, container, kubeContainerID.ID, v1.EventTypeWarning, events.FailedPostStartHook, msg)
			if err := m.killContainer(pod, kubeContainerID, container.Name, "FailedPostStartHook", reasonFailedPostStartHook, nil); err != nil {
				klog.ErrorS(err, "Failed to kill container", "pod", klog.KObj(pod),
					"podUID", pod.UID, "containerName", container.Name, "containerID", kubeContainerID.String())
			}
			return msg, ErrPostStartHook
		}
	}

	return "", nil
}
```
