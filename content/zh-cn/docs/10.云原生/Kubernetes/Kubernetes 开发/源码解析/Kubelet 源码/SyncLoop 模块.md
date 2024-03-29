---
title: "SyncLoop 模块"
---

# 概述

> 参考：

SyncLoop 模块，Kubelet 同步循环

![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/paseql/1645680266069-2cc34f9d-ed55-45dd-8df4-ac0ed0c8c388.png)

kubelet 的工作核心就是在围绕着不同的生产者生产出来的不同的有关 pod 的消息来调用相应的消费者（不同的子模块）完成不同的行为(创建和删除 pod 等)，即图中的控制循环（SyncLoop），通过不同的事件驱动这个控制循环运行。

本文仅分析新建 pod 的流程，当一个 pod 完成调度，与一个 node 绑定起来之后，这个 pod 就会触发 kubelet 在循环控制里注册的 handler，上图中的 HandlePods 部分。此时，通过检查 pod 在 kubelet 内存中的状态，kubelet 就能判断出这是一个新调度过来的 pod，从而触发 Handler 里的 ADD 事件对应的逻辑处理。然后 kubelet 会为这个 pod 生成对应的 podStatus，接着检查 pod 所声明的 volume 是不是准备好了，然后调用下层的容器运行时。如果是 update 事件的话，kubelet 就会根据 pod 对象具体的变更情况，调用下层的容器运行时进行容器的重建。

## kubelet 创建 pod 的流程![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/paseql/1645680598998-2bf06936-f8ac-48fa-b370-796e5c334545.png)

Note：注意 14，,15 步，kubelet 会先将生成配置(volume 挂载、配置主机名等等)，才会去启动 pod，哪怕 pod 启动失败，挂载依然存在。

# 同步循环

## kubelet.syncLoop() # 同步循环入口

syncLoop 中首先定义了一个 syncTicker 和 housekeepingTicker，即使没有需要更新的 pod 配置，kubelet 也会定时去做同步和清理 pod 的工作。然后在 for 循环中一直调用 syncLoopIteration，如果在每次循环过程中出现比较严重的错误，kubelet 会记录到 runtimeState 中，遇到错误就等待 5 秒中继续循环。

源码：`pkg/kubelet/kubelet.go`

```go
func (kl *Kubelet) syncLoop(updates <-chan kubetypes.PodUpdate, handler SyncHandler) {
	klog.InfoS("Starting kubelet main sync loop")

    // 每秒检测一次是否有需要同步的 pod workers。同步周期默认 10 秒
	syncTicker := time.NewTicker(time.Second)
	defer syncTicker.Stop()

    // 每两秒检测一次是否有需要清理的 pod
	housekeepingTicker := time.NewTicker(housekeepingPeriod)
	defer housekeepingTicker.Stop()

    // pod 的生命周期变化
	plegCh := kl.pleg.Watch()
	const (
		base   = 100 * time.Millisecond
		max    = 5 * time.Second
		factor = 2
	)
	duration := base
	// Responsible for checking limits in resolv.conf
	// The limits do not have anything to do with individual pods
	// Since this is called in syncLoop, we don't need to call it anywhere else
	if kl.dnsConfigurer != nil && kl.dnsConfigurer.ResolverConfig != "" {
		kl.dnsConfigurer.CheckLimitsForResolvConf()
	}

	for {
		if err := kl.runtimeState.runtimeErrors(); err != nil {
			klog.ErrorS(err, "Skipping pod synchronization")
			// exponential backoff
			time.Sleep(duration)
			duration = time.Duration(math.Min(float64(max), factor*float64(duration)))
			continue
		}
		// reset backoff if we have a success
		duration = base

		kl.syncLoopMonitor.Store(kl.clock.Now())
        // 第二个参数为 SyncHandler 类型，SyncHandler 是一个 interface，
        // 在该文件开头处定义
		if !kl.syncLoopIteration(updates, handler, syncTicker.C, housekeepingTicker.C, plegCh) {
			break
		}
		kl.syncLoopMonitor.Store(kl.clock.Now())
	}
}
```

## kubelet.syncLoopIteration() # 监听 Pod 变化

syncLoopIteration 这个方法就会对多个管道进行遍历，发现任何一个管道有消息就交给 handler 去处理。它会从以下管道中获取消息：

- configCh：该信息源由 kubeDeps 对象中的 PodConfig 子模块提供，该模块将同时 watch 3 个不同来源的 pod 信息的变化（file，http，apiserver），一旦某个来源的 pod 信息发生了更新（创建/更新/删除），这个 channel 中就会出现被更新的 pod 信息和更新的具体操作。
- syncCh：定时器管道，每隔一秒去同步最新保存的 pod 状态
- houseKeepingCh：housekeeping 事件的管道，做 pod 清理工作
- plegCh：该信息源由 kubelet 对象中的 pleg 子模块提供，该模块主要用于周期性地向 container runtime 查询当前所有容器的状态，如果状态发生变化，则这个 channel 产生事件。
- livenessManager.Updates()：健康检查发现某个 pod 不可用，kubelet 将根据 Pod 的 restartPolicy 自动执行正确的操作

源码：`pkg/kubelet/kubelet.go`-`kubelet.syncLoopIteration()`

## Pod 处理入口

kubelet.syncLoopIteration() 中的 SyncHandler 参数是一个接口，kubelet 结构体已实现该接口。后续 Pod 的增删改查行为，使用的都是通过实现了 SyncHandler 接口的 kubelet 结构体上的方法。

所以 `kubelet.HandlePodXXX()` 就是 `SyncHandler.HandlePodXXX()`，在这些处理 Pod 的方法中，都会调用 `kubelet.dispatchWork()` 方法，把对 Pod 的操作下发给 podWorkers()。

### kubelet.HandlePodAdditions() # 处理新增 Pod

### kubelet.HandlePodUpdates() # 处理更新 Pod

### kubelet.HandlePodRemoves() # 处理删除 Pod

### kubelet.HandlePodReconcile() # 处理调谐 Pod

## kubelet.dispatchWork() # 下发工作

`kubelet.dispatchWorker()` 的主要作用是把某个对 Pod 的操作（创建/更新/删除）下发给 PodWorkers 模块。这里面说的对 Pod 的操作，就是上面 SyncHandler 接口下的几个方法。

源码：`pkg/kubelet/kubelet.go`

```go
func (kl *Kubelet) dispatchWork(pod *v1.Pod, syncType kubetypes.SyncPodType, mirrorPod *v1.Pod, start time.Time) {
	// 在一个异步工作器中运行同步。所有对 Pod 的操作都转到 PodWorker 模块中
    // 由实现了 PodWorkers interface{} 的 podWorkers struct{} 的各种方法进行处理
	kl.podWorkers.UpdatePod(UpdatePodOptions{
		Pod:        pod,
		MirrorPod:  mirrorPod,
		UpdateType: syncType,
		StartTime:  start,
	})
	// Note the number of containers for new pods.
	if syncType == kubetypes.SyncPodCreate {
		metrics.ContainersPerPodCount.Observe(float64(len(pod.Spec.Containers)))
	}
}
```
