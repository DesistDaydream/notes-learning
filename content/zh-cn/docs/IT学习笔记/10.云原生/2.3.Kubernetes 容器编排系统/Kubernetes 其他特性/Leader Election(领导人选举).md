---
title: Leader Election(领导人选举)
---

# 概述

> 参考：
> - [官方博客,Kubernetes 的简单领导人选举](https://kubernetes.io/blog/2016/01/simple-leader-election-with-kubernetes/)
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

**备注**：该例子取自项目[文档](https://github.com/kubernetes/contrib/tree/master/election)。

### 启动一个 leader-elector 的 Pod

- **创建一个 **

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
$ kubectl logs -f ${pod_name}
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
$ kubectl get ep example -o yaml
```

通过查看 annotations 中的 `control-plane.alpha.kubernetes.io/leader` 字段来获得 Leader 的信息；

- **使用 **；
  `leader-elector` 实现了一个简单的 HTTP 接口（`:4040`）来查看当前 Leader：

```bash
curl http://localhost:8001/api/v1/namespaces/default/pods/leader-elector-5d77ccc44d-gwsgg:4040/proxy/
{"name":"leader-elector-5d77ccc44d-7tmgm"}
```

###

# 如何使用 leader-elector

如果自己的项目中需要用到 Leader Election 的逻辑，可以有两种方式：

## 将调用 `leaderelection` 库的逻辑内嵌到自己项目中

Leader Election 库在 <https://github.com/kubernetes/client-go/tree/master/tools/leaderelection>

## 使用 Sidecar 的方式

# Leader Election 的实现

Leader Election 的过程本质上就是一个**竞争分布式锁**的过程。在 Kubernetes 中，这个分布式锁是通过下面几个 **Resource(资源) **实现的：

- **Endpoints** #
- **ConfigMaps** #
- **Leases** # 详见 [集群资源-Lease](✏IT 学习笔记/☁️10.云原生/2.3.Kubernetes%20 容器编排系统/1.API、Resource(资源)、Object(对象)/API%20 参考/集群资源.md 参考/集群资源.md)

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

1. **每个 Pod 在启动的时候都会创建 **；
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
综上所述，Kubernetes 中 Pod 的选举过程本质上还是为了服务的高可用。**希望大家研究得愉快**！
