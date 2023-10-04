---
title: Leader 选举
---

# 概述

> 参考：
> - [原文链接](https://mp.weixin.qq.com/s/EPKShekTZWe04H1X2E21LQ)

领导者选举要解决什么问题呢？首先，一个分布式集群中运行了多个组件，每个组件负责自身重要的功能。其中有一个组件因为某些原因而退出，此时整个集群的运作都受到了影响。**领导者选举就是要保证每个组件的高可用性**，例如，在 Kubernetes 集群中，允许同时运行多个 kube-scheduler 节点，其中正常工作的只有一个 kube-scheduler 节点（即领导者节点），其他 kube-scheduler 节点为候选（Candidate）节点并处于阻塞状态。在领导者节点因某些原因而退出后，其他候选节点则通过领导者选举机制竞选，有一个候选节点成为领导者节点并接替之前领导者节点的工作。领导者选举机制如下图所示。
![](https://notes-learning.oss-cn-beijing.aliyuncs.com/d05e3808-6527-434a-b7bb-bdf3bc826284/640)

## 领导者选举机制

领导者选举机制是**分布式锁**机制的一种，实现分布式锁有多种方式，例如可通过 ZooKeeper、Redis、Etcd 等存储服务。Kubernetes 系统依赖于 Etcd 做存储服务，系统中其他组件也是通过 Etcd 实现分布式锁的。kube-scheduler 组件在 Etcd 上实现分布式锁的原理如下。

- 分布式锁依赖于 Etcd 上的一个 key，key 的操作都是**原子操作，**将 key 作为分布式锁，它有两种状态——存在和不存在。
- key（分布式锁）不存在时：多节点中的一个节点成功创建该 key（获得锁）并写入自身节点的信息，获得锁的节点被称为领导者节点。领导者节点会定时更新（续约）该 key 的信息。
- key（分布式锁）存在时：其他节点处于阻塞状态并定时获取锁，这些节点被称为候选节点。候选节点定时获取锁的过程如下：定时获取 key 的数据，验证数据中领导者租约是否到期，如果未到期则不能抢占它，如果已到期则更新 key 并写入自身节点的信息，更新成功则成为领导者节点。

# 资源锁

Kubernetes 支持 3 种资源锁，资源锁的意思是基于 Etcd 集群的 key 在依赖于 Kubernetes 的某种资源下创建的分布式锁。3 种资源锁介绍如下：

- **EndpointsResourceLock**：依赖于 Endpoints 资源，默认资源锁为该类型。
- **ConfigMapsResourceLock**：依赖于 Configmaps 资源。
- **LeasesResourceLock**：依赖于 Leases 资源。

可通过 `--leader-elect-resource-lock` 参数指定使用哪种资源锁，如不指定则`EndpointsResourceLock` 为默认资源锁。它的 key（分布式锁）存在于 Etcd 集群的 `/registry/services/endpoints/kube-system/kube-scheduler` 中。该 key 中存储的是竞选为领导者节点的信息，它通过 `LeaderElectionRecord` 结构体进行描述：

```go
# 源码路径：vendor/k8s.io/client-go/tools/leaderelection/resourcelock/interface.go

type LeaderElectionRecord struct {
    HolderIdentity string
    LeaseDurationSeconds int
    AcquireTime metav1.Time
    RenewTime metav1.Time
    LeaderTransitions int
}
```

- **HolderIdentity**：领导者身份标识，通常为 Hostname\_\<hash 值>。
- **LeaseDurationSeconds**：领导者租约的时长。
- **AcquireTime**：领导者获得锁的时间。
- **RenewTime**：领导者续租的时间。
- **LeaderTransitions**：领导者选举切换的次数。

每种资源锁实现了对 key（资源锁）的操作方法，它的接口定义如下：

```go
# 代码路径：vendor/k8s.io/client-go/tools/leaderelection/resourcelock/interface.go

type Interface interface {
    Get() (LeaderElectionRecord, error)
    Create(ler LeaderElectionRecord) error
    Update(ler LeaderElectionRecord) error
    RecordEvent(string)
    Identity() string
    Describe() string
}
```

Get 方法用于获取资源锁的所有信息，Create 方法用于创建资源锁，Update 方法用于更新资源锁信息，RecordEvent 方法通过 EventBroadcaster 事件管理器记录事件，Identity 方法用于获取领导者身份标识，Describe 方法用于获取资源锁的信息。

# 领导者选举过程

领导者选举过程如下图所示：![](https://notes-learning.oss-cn-beijing.aliyuncs.com/d05e3808-6527-434a-b7bb-bdf3bc826284/640)

领导者选举过程

`le.acquire` 函数尝试从 Etcd 中获取资源锁，领导者节点获取到资源锁后会执行 kube-scheduler 的主要逻辑（即 `le.config.Callbacks.OnStartedLeading` 回调函数），并通过 `le.renew` 函数定时（默认值为 2 秒）对资源锁续约。候选节点获取不到资源锁，它不会退出并定时（默认值为 2 秒）尝试获取资源锁，直到成功为止。代码示例如下：

```go
# 代码路径：vendor/k8s.io/client-go/tools/leaderelection/leaderelection.go

func (le *LeaderElector) Run(ctx context.Context) {
    defer func() {
      runtime.HandleCrash()
      le.config.Callbacks.OnStoppedLeading()
    }()
    if !le.acquire(ctx) {
      return
    }
    ...
    go le.config.Callbacks.OnStartedLeading(ctx)
    le.renew(ctx)
}
```

> 1. 资源锁获取过程

```go
func (le *LeaderElector) acquire(ctx context.Context) bool {
  ...
  wait.JitterUntil(func() {
      succeeded = le.tryAcquireOrRenew()
      le.maybeReportTransition()
      if !succeeded {
        return
      }
      ...
      cancel()
    }, le.config.RetryPeriod, JitterFactor, true, ctx.Done())
  return succeeded
}
```

获取资源锁的过程通过 `wait.JitterUntil` 定时器定时执行，它接收一个 func 匿名函数和一个 stopCh Chan，内部会定时调用匿名函数，只有当 stopCh 关闭时，该定时器才会停止并退出。

执行 `le.tryAcquireOrRenew` 函数来获取资源锁。如果其获取资源锁失败，会通过 return 等待下一次定时获取资源锁。如果其获取资源锁成功，则说明当前节点可以成为领导者节点，退出 acquire 函数并返回 true。`le.tryAcquireOrRenew` 代码示例如下。

（1）首先，通过 `le.config.Lock.Get` 函数获取资源锁，当资源锁不存在时，当前节点创建该 key（获取锁）并写入自身节点的信息，创建成功则当前节点成为领导者节点并返回 true。

```go
oldLeaderElectionRecord, err := le.config.Lock.Get()
if err != nil {
    if !errors.IsNotFound(err) {
      return false
    }
    if err = le.config.Lock.Create(leaderElectionRecord); err != nil {
        return false
    }
    le.observedRecord = leaderElectionRecord
    le.observedTime = le.clock.Now()
    return true
}
```

（2）当资源锁存在时，更新本地缓存的租约信息。

```go
if !reflect.DeepEqual(le.observedRecord, *oldLeaderElectionRecord) {
    le.observedRecord = *oldLeaderElectionRecord
    le.observedTime = le.clock.Now()
}
```

（3）候选节点会验证领导者节点的租约是否到期，如果尚未到期，暂时还不能抢占并返回 false。

```go
if len(oldLeaderElectionRecord.HolderIdentity) > 0 &&
le.observedTime.Add(le.config.LeaseDuration).After(now.Time) &&
!le.IsLeader() {
    ...
    return false
}
```

（4）如果是领导者节点，那么 AcquireTime（资源锁获得时间）和 LeaderTransitions（领导者进行切换的次数）字段保持不变。如果是候选节点，则说明领导者节点的租约到期，给 LeaderTransitions 字段加 1 并抢占资源锁。

```go
if le.IsLeader() {
    leaderElectionRecord.AcquireTime = oldLeaderElectionRecord.AcquireTime
    leaderElectionRecord.LeaderTransitions = oldLeaderElectionRecord.LeaderTransitions
} else {
    leaderElectionRecord.LeaderTransitions = oldLeaderElectionRecord.LeaderTransitions + 1
}
```

（5）通过 `le.config.Lock.Update` 函数尝试去更新租约记录，若更新成功，函数返回 true。

```go
if err = le.config.Lock.Update(leaderElectionRecord); err != nil {
    klog.Errorf("Failed to update lock: %v", err)
    return false
}
...
return true
```

> 2. 领导者节点定时更新租约过程

在领导者节点获取资源锁以后，会定时（默认值为 2 秒）循环更新租约信息，以保持长久的领导者身份。若因网络超时而导致租约信息更新失败，则说明被候选节点抢占了领导者身份，当前节点会退出进程。代码示例如下：

```go
# 代码路径：vendor/k8s.io/client-go/tools/leaderelection/leaderelection.go

func (le *LeaderElector) renew(ctx context.Context) {
    ...
    wait.Until(func() {
        ...
        err := wait.PollImmediateUntil(le.config.RetryPeriod, func() (bool, error) {
            done := make(chan bool, 1)
            go func() {
                defer close(done)
                done <- le.tryAcquireOrRenew()
            }()
            ...
        }, timeoutCtx.Done())
        ...
        if err == nil {
            klog.V(5).Infof("successfully renewed lease %v", desc)
            return
        }
        ...
        cancel()
    }, le.config.RetryPeriod, ctx.Done())

    if le.config.ReleaseOnCancel {
        le.release()
    }
}
```

领导者节点续约的过程通过 `wait.PollImmediateUntil` 定时器定时执行，它接收一个 func 匿名函数（条件函数）和一个 stopCh，内部会定时调用条件函数，当条件函数返回 true 或 stopCh 关闭时，该定时器才会停止并退出。

执行 `le.tryAcquireOrRenew` 函数来实现领导者节点的续约，其原理与资源锁获取过程相同。le.tryAcquireOrRenew 函数返回 true 说明续约成功，并进入下一个定时续约；返回 false 则退出并执行 le.release 函数且释放资源锁。

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/d05e3808-6527-434a-b7bb-bdf3bc826284/640)
本文授权转载于：Kubernetes 源码剖析，作者：郑东旭
