---
title: Leader Election(领导人选举)
---

# 概述

> 参考：
>
> - [官方博客, Kubernetes 的简单领导人选举](https://kubernetes.io/blog/2016/01/simple-leader-election-with-kubernetes/)
> - [zhengyinyong](https://zhengyinyong.com/post/kubernetes-pod-leader-election/)
> - 用法：
>   - [公众号-云原生实验室，巧用 Kubernetes 中的 Leader 选举机制来实现自己的 HA 应用](https://mp.weixin.qq.com/s/LnF9ZIoi-sCUV9rxxEvDvQ)

## 为什么需要 Pod 之间的 Leader Election

一般来说，由 Deployment 创建的 1 个或多个 Pod 都是对等关系，彼此之间提供一样的服务。但是在某些场合，多个 Pod 之间需要有一个 Leader 的角色，即：

- **Pod 之间有且只有一个 Leader**；
- **Leader 在一定周期不可用时，其他 Pod 会再选出一个 Leader**；
- **由处于 Leader 身份的 Pod 来完成某些特殊的业务逻辑（通常是写操作）**；

比如，当**多个 Pod 之间只需要一个写者**时，如果不采用 Leader Election，那么就必须在 Pod 启动之初人为地配置一个 Leader。如果配置的 Leader 在后续的服务中失效且没有对应机制来生成新的 Leader，那么对应 Pod 服务就可能处于不可用状态，违背高可用原则。

典型地，Kubernetes 的核心组件 kube-controller-manager 和 scheduler 就需要一个需要 Leader 的场景。当 kube-controller-manager 的启动参数设置 `--leader-elect=true` 时，对应节点的 kube-controller-manager 在启动时会执行选主操作。当选出一个 Leader 之后，由 Leader 来启动所有的控制器。如果 Leader Pod 不可用，将会自动选出新的 Leader Pod，从而保障控制器仍处于运行状态。

## 一个简单的 Leader Election 的例子

**备注**：该例子取自项目[文档](https://github.com/kubernetes/contrib/tree/master/election)

### 启动一个 leader-elector 的 Pod

- **创建一个**

```bash
$ kubectl run leader-elector \
  --image=k8s.gcr.io/leader-elector:0.5 \
  --replicas=3 \
  -- \
  --election=example \
  --http=0.0.0.0:4040
```

- 副本数为 3，即将生成 3 个 Pod，如果运行成功，可观察到：

```bash
$ kubectl get po
 NAME                              READY   STATUS    RESTARTS   AGE
 leader-elector-68dcb58d55-7dhdz   1/1     Running   0          2m36s
 leader-elector-68dcb58d55-g5zp8   1/1     Running   0          2m36s
 leader-elector-68dcb58d55-q45pd   1/1     Running   0          2m36s
```

- **查看哪个 Pod 成为 Leader**；

可以逐个查看 Pod 的日志：

```bash
kubectl logs -f ${pod_name}
```

如果是 Leader 的话，将会有如下的日志：

```bash
$ kubectl logs leader-elector-68dcb58d55-g5zp8
 leader-elector-9577494c7-l64lp is the leader
 I0122 03:24:31.779331       8 leaderelection.go:296] lock is held by leader-elector-9577494c7-l64lp and has not yet expired
 I0122 03:24:36.101800       8 leaderelection.go:296] lock is held by leader-elector-9577494c7-l64lp and has not yet expired
 I0122 03:24:41.426387       8 leaderelection.go:296] lock is held by leader-elector-9577494c7-l64lp and has not yet expired
 I0122 03:24:45.947321       8 leaderelection.go:215] sucessfully acquired lease default/example
leader-elector-68dcb58d55-g5zp8 is the leader
```

更通用的方式是查看资源锁的身份标识信息：

```bash
kubectl get ep example -o yaml
```

通过查看 annotations 中的 `control-plane.alpha.kubernetes.io/leader` 字段来获得 Leader 的信息；

- **使用** `leader-elector` 实现了一个简单的 HTTP 接口（`:4040`）来查看当前 Leader：

```bash
curl http://localhost:8001/api/v1/namespaces/default/pods/leader-elector-5d77ccc44d-gwsgg:4040/proxy/
{"name":"leader-elector-5d77ccc44d-7tmgm"}
```

# 如何使用 leader-elector

如果自己的项目中需要用到 Leader Election 的逻辑，可以有两种方式：

## 将调用 `leaderelection` 库的逻辑内嵌到自己项目中

Leader Election 库在 <https://github.com/kubernetes/client-go/tree/master/tools/leaderelection>

## 使用 Sidecar 的方式

# Leader Election 的实现

Leader Election 的过程本质上就是一个**竞争分布式锁**的过程。在 Kubernetes 中，这个分布式锁是通过下面几个 **Resource(资源)**实现的：

- **Endpoints** #
- **ConfigMaps** #
- **Leases** # 详见 [集群资源-Lease](/docs/10.云原生/2.3.Kubernetes%20 容器编排系统/1.API、Resource(资源)、Object(对象)/API%20 参考/集群资源.md 参考/集群资源.md)

**谁先创建了某种资源，谁就获得锁**。通常情况下，kube-scheduler 和 kube-controller-manager 使用 leases 资源来实现领导者选举。

按照我们以往的惯例，带着问题去看源码。有这么几个问题：

- **Leader Election 如何竞选** ？
- **Leader 不可用之后如何竞选新的 Leader** ？

不同于 Raft 算法的一致性算法的 Leader 竞选，Pod 之间的 Leader Election 是**无状态**的，也就是说**现在的 Leader 无需同步上一个 Leader 的数据信息，这就把竞选的过程变得非常简单：先到先得**。说白了，就是谁先创建了这个资源，谁就是领导者了。

这部分代码在 `kubernetes/staging/src/k8s.io/client-go/tools/leaderelection` 中，取 1.9.2 版本来分析。

### 资源锁的实现

Kubernetes 实现了两种资源锁（resourcelock）：Endpoint 和 ConfigMap。**如果是基于 Endpoint 的资源锁，获取到锁的 Pod 将会在对应 Namespace 下创建对应的 Endpoint 对象，并在其 Annotations 上记录 Pod 的信息**。

比如 kube-controller-manager：

```bash
$ kubectl get ep -n kube-system | grep kube-controller-manager
kube-controller-manager   <none>   41d

$ kubectl describe ep kube-controller-manager -n kube-system
Name:         kube-controller-manager
Namespace:    kube-system
Labels:       <none>
Annotations:  control-plane.alpha.kubernetes.io/leader:
                {"holderIdentity":"szdc-k8sm-0-5","leaseDurationSeconds":15,"acquireTime":"2018-12-11T0...
Subsets:
Events:  <none>
```

发现在 kube-system 中创建了同名的 Endpoint（kube-controller-manager），并在 Annotations 中以设置了 key 为 `control-plane.alpha.kubernetes.io/leader`，value 为对应 Leader 信息的 JSON 数据。同理，如果采用 ConfigMap 作为资源锁也是类似的实现模式。

resourcelock 是以 interface 的形式对外暴露，在创建过程（`New()`）通过相应的参数来控制具体实例化的过程：

```go
// leaderelection/resourcelock/interface.go
type Interface interface {
 // Get returns the LeaderElectionRecord
 Get() (*LeaderElectionRecord, error)
 // Create attempts to create a LeaderElectionRecord
 Create(ler LeaderElectionRecord) error
 // Update will update and existing LeaderElectionRecord
 Update(ler LeaderElectionRecord) error
 // RecordEvent is used to record events
 RecordEvent(string)
 // Identity will return the locks Identity
 Identity() string
 // Describe is used to convert details on current resource lock
 // into a string
 Describe() string
}
```

其中 `Get()`、`Create()` 和 `Update()` 本质上就是对 `LeaderElectionRecord` 的读写操作。`LeaderElectionRecord` 定义如下：

```go
type LeaderElectionRecord struct {
 // 标示当前资源锁的所有权的信息
 HolderIdentity string `json:"holderIdentity"`
 // 资源锁租约时间是多长
 LeaseDurationSeconds int `json:"leaseDurationSeconds"`
 // 锁获得的时间
 AcquireTime metav1.Time `json:"acquireTime"`
 // 续租的时间
 RenewTime metav1.Time `json:"renewTime"`
 // Leader 进行切换的次数
 LeaderTransitions int `json:"leaderTransitions"`
}
```

理论上，`LeaderElectionRecord` 是保存在资源锁的 Annotations 中，可以是任意的字符串，此处是将 JSON 序列化为字符串来进行存储。
在 `leaderelection/resourcelock/configmaplock.go` 和 `leaderelection/resourcelock/endpointslock.go` 分别是基于 Endpoint 和 ConfigMap 对上面接口的实现。拿 `endpointslock.go` 来看，对这几个接口的实现实际上就是对 Endpoint 资源中 Annotations 的增删查改罢了，比较简单，就不详细展开。

### 竞争锁的过程

完整的 Leader Election 过程在 `leaderelection/leaderelection.go` 中。
整个过程可以简单描述为：

1. **每个 Pod 在启动的时候都会创建**；
2. **在循环中，Pod 会定期（**；
3. **在循环周期中，Leader 会不断 Update 资源锁的对应时间信息，从节点则会不断检查资源锁是否过期，如果过期则尝试更新资源，标记资源所有权。这样一来，一旦 Leader 不可用，则对应的资源锁将得不到更新，过期之后其他从节点会再次创建新的资源锁成为 Leader**；

其中，`LeaderElector.Run()` 的源码为：

```go
func (le *LeaderElector) Run() {
    ...
    // 尝试创建锁
    le.acquire()
    // Leader 更新资源锁的租约
    le.renew()
    ...
}
```

`acquire()` 会周期性地创建锁或探查锁有没有过期：

```go
func (le *LeaderElector) acquire() {
    ...
    wait.JitterUntil(func() {
        // 尝试创建或者续约资源锁
        succeeded := le.tryAcquireOrRenew()
        // leader 可能发生了改变，执行相应的 OnNewLeader() 回调函数
        le.maybeReportTransition()
        // 不成功说明创建资源失败，当前 Leader 是其他 Pod
        if !succeeded {
            ...
            return
        }
        ...
    }, le.config.RetryPeriod, JitterFactor, true, stop)
}
```

执行的周期为 `RetryPeriod`。

我们重点关注 `tryAcquireOrRenew()` 的逻辑：

```go
func (le *leaderElector) tryAcquireOrRenew() bool {
    now := metav1.Now()
    leaderElectionRecord := rl.LeaderElectionRecord{
        HolderIdentity:       le.config.Lock.Identity(),
        LeaseDurationSeconds: int(le.config.LeaseDuration) / time.Second),
        // 将租约改成 now
        RenewTime:            now,
        AcquireTime:          now,
    }

    // 获取当前的资源锁
    oldLeaderElectionRecord, err := le.config.Lock.Get()
    if err != nil {
        ...
        // 执行到这里说明找不到资源锁，执行资源锁的创建动作
        // 由于资源锁对应的底层 Kubernetes 资源 Endpoint 或 ConfigMap 是不可重复创建的，所以此处创建是安全的
        if err = le.config.Lock.Create(leaderElectionRecord); err != nil {
            ...
        }
        ...
    }

    // 如果当前已经有 Leader，进行 Update 操作
    // 如果当前是 Leader：Update 操作就是续租动作，即将对应字段的时间改成当前时间
    // 如果是非 Leader 节点且可运行 Update 操作，则是一个抢夺资源锁的过程，谁先更新成功谁就抢到资源
    ...
    // 如果还没有过期且当前不是 Leader，直接返回
    // 只有 Leader 才进行续租操作且此时其他节点无须抢夺资源锁
    if le.observedTime.Add(le.config.LeaseDuration).After(now.Time) &&
       oldLeaderElectionRecord.HolderIdentity != le.config.Lock.Identity() {
           ...
           return false
    }
    ...
    // 更新资源
    // 对于 Leader 来说，这是一个续租的过程
    // 对于非 Leader 节点（仅在上一个资源锁已经过期），这是一个更新锁所有权的过程
    if err = le.config.Lock.Update(leaderElectionRecord); err != nil {
        ...
    }
}
```

由上可以看出，`tryAcquireOrRenew()` 就是一个不断尝试 Update 操作的过程。

如果执行逻辑从 `le.acquire()` 跳出，往下执行 `le.renew()`，这说明当前 Pod 已经成功抢到资源锁成为 Leader，必须定期续租：

```go
func (le *LeaderElector) renew() {
    stop := make(chan struct{})
    // period 为 0 说明会一直执行
    wait.Until(func() {
        // 每间隔 RetryPeriod 就执行 tryAcquireOrRenew()
        // 如果 tryAcquireOrRenew() 返回 false 跳出 Poll()
        // tryAcquireOrRenew() 返回 false 说明续租失败
        err := wait.Poll(le.config.RetryPeriod, le.config.RenewDeadline, func() (bool, error) {
            return le.tryAcquireOrRenew(), nil
        })

        // 续租失败，说明已经不是 Leader
        ...
    }, 0, stop)
}
```

### 如何使用 `leaderelection` 库

让我们来关注一下 [election](https://github.com/kubernetes/contrib/tree/master/election) 的实现。

主要的逻辑位于 `election/lib/election.go`：

```go
func RunElection(e *leaderelection.LeaderElector) {
    wait.Forever(e.Run, 0)
}
```

主体逻辑很简单，就是不断执行 `Run()`。而 `Run()` 的实现就是上文中 `leaderelection` 的 `Run()` 。

上层应用只需要创建（`NewElection()`）创建 `LeaderElector` 对象，然后在一个 loop 中调用 `Run()` 即可。

综上所述，Kubernetes 中 Pod 的选举过程本质上还是为了服务的高可用。

# 公众号 - k8s技术圈，Kubernetes 源码剖析之 Leader 选举

> - [公众号 - k8s技术圈，Kubernetes 源码剖析之 Leader 选举](https://mp.weixin.qq.com/s/EPKShekTZWe04H1X2E21LQ)

领导者选举要解决什么问题呢？首先，一个分布式集群中运行了多个组件，每个组件负责自身重要的功能。其中有一个组件因为某些原因而退出，此时整个集群的运作都受到了影响。**领导者选举就是要保证每个组件的高可用性**，例如，在 Kubernetes 集群中，允许同时运行多个 kube-scheduler 节点，其中正常工作的只有一个 kube-scheduler 节点（即领导者节点），其他 kube-scheduler 节点为候选（Candidate）节点并处于阻塞状态。在领导者节点因某些原因而退出后，其他候选节点则通过领导者选举机制竞选，有一个候选节点成为领导者节点并接替之前领导者节点的工作。领导者选举机制如下图所示。

## 领导者选举机制

领导者选举机制是**分布式锁**机制的一种，实现分布式锁有多种方式，例如可通过 ZooKeeper、Redis、Etcd 等存储服务。Kubernetes 系统依赖于 Etcd 做存储服务，系统中其他组件也是通过 Etcd 实现分布式锁的。kube-scheduler 组件在 Etcd 上实现分布式锁的原理如下。

- 分布式锁依赖于 Etcd 上的一个 key，key 的操作都是**原子操作**，将 key 作为分布式锁，它有两种状态——存在和不存在。
- key（分布式锁）不存在时：多节点中的一个节点成功创建该 key（获得锁）并写入自身节点的信息，获得锁的节点被称为领导者节点。领导者节点会定时更新（续约）该 key 的信息。
- key（分布式锁）存在时：其他节点处于阻塞状态并定时获取锁，这些节点被称为候选节点。候选节点定时获取锁的过程如下：定时获取 key 的数据，验证数据中领导者租约是否到期，如果未到期则不能抢占它，如果已到期则更新 key 并写入自身节点的信息，更新成功则成为领导者节点。

## 资源锁

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

## 领导者选举过程

领导者选举过程如下图所示：

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

