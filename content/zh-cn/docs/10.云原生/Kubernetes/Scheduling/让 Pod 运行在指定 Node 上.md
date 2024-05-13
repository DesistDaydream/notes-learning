---
title: 让 Pod 运行在指定 Node 上
---

# 概述

> 参考：
>
> - [官方文档,概念-调度、抢占、驱逐-让 Pod 运行在指定节点上](https://kubernetes.io/docs/concepts/scheduling-eviction/assign-pod-node/)
> - [官方文档,概念-调度、抢占、驱逐-污点与容忍度](https://kubernetes.io/docs/concepts/scheduling-eviction/taint-and-toleration/)
> - [Kubernetes 内置的 标签、注释、污点](https://kubernetes.io/docs/reference/labels-annotations-taints/)
> - [CSDN,k8s 之 pod 亲和与反亲和的 topologyKey](https://blog.csdn.net/asdfsadfasdfsa/article/details/106027367)

通常情况下 Scheduler(调度器) 将自动进行合理的分配(例如，将 Pods 分散到所有节点上，以防止单独节点上资源使用率远高于其他节点)。但是在某些情况下我们需要更多控制 Pods 落在某个指定的节点上，例如确保一个 Pod 部署在装有 SSD 的机器上，或者将两个不同服务中的 Pods 共同定位到同一可用区域。

所以我们需要 constrain(约束) Pod 运行在指定的节点。可以实现该效果的方式有以下几种：

- **nodeName(节点名称)** #
- **nodeSelector(节点选择器)** # 根据节点的标签，选择 pod 要运行在哪个节点上
  - 这种行为定义 Pod 必须在特定节点上运行。
- **Affinity(亲和) 与 Anti-Affinity(反亲和)** # 根据亲和原则，让 pod 更趋向于与哪些 XXX 运行在同一个节点
  - 这种行为定义 Pod 更倾向于在特定节点上运行。
- **Taint(污点) 与 Toleration(容忍度)** # 根据节点上的污点，以及 pod 是否可以容忍该污点来决定 pod 是否可以运行在哪些节点上

其中 nodeSelector 和 Affinity 与 Anti-Affinity 是通过 [Label Selectors(标签选择器)](/docs/10.云原生/Kubernetes/API%20Resource%20与%20Object/Object%20管理/Label%20and%20Selector(标签和选择器)/Label%20and%20Selector(标签和选择器).md) 来实现的。而 Taint 与 Toleration 是另一套类似于标签选择器的机制。

# nodeName(节点名称)

`nodeName` 是节点选择约束的最简单方法，但是由于其自身限制，通常不使用它。 `nodeName` 是 PodSpec 的一个字段。 如果它不为空，调度器将忽略 Pod，并且给定节点上运行的 kubelet 进程尝试执行该 Pod。 因此，如果 `nodeName` 在 PodSpec 中指定了，则它优先于上面的节点选择方法。
使用 `nodeName` 来选择节点的一些限制：

- 如果指定的节点不存在，
- 如果指定的节点没有资源来容纳 Pod，Pod 将会调度失败并且其原因将显示为， 比如 OutOfmemory 或 OutOfcpu。
- 云环境中的节点名称并非总是可预测或稳定的。

下面的是使用 `nodeName` 字段的 Pod 配置文件的例子：

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: nginx
spec:
  containers:
    - name: nginx
      image: nginx
  nodeName: kube-01
```

上面的 pod 将运行在 kube-01 节点上。

# NodeSelector(节点选择器)

nodeSelector(节点选择器) 是最简单的约束行为。在 Pod 的 manifest 中使用 `nodeSelector` 字段进行配置。`nodeSelector` 字段以键值对的方式描述该 Pod 应该运行在哪些具有指定标签的节点上。为了使 Pod 可以运行在指定的节点上，那么该节点的所有标签中，必须包含 Pod 中 .spec.nodeSelector 字段下定义的每个键值对.

## 使用示例

首先为一个节点设置标签，标签设置方法详见[本章节](Label%20and%20Selector(标签和选择器).md and Selector(标签和选择器).md)。假如现在为节点设置了 disktype: ssd 标签

```bash
# 设置标签
root@desistdaydream:~# kubectl label nodes master-3.bj-test disktype=ssd
node/master-3.bj-test labeled
# 查看该节点标签
root@desistdaydream:~# kubectl get nodes master-3.bj-test --show-labels
NAME               STATUS   ROLES         AGE    VERSION   LABELS
master-3.bj-test   Ready    master,work   183d   v1.19.2   ....,disktype=ssd...这里省略了很多其他标签
```

在 Pod 的 Manifest 中添加 `nodeSelector` 字段，示例如下：

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: nginx
  labels:
    env: test
spec:
  containers:
    - name: nginx
      image: nginx
      imagePullPolicy: IfNotPresent
  nodeSelector:
    disktype: ssd
```

该示例 Pod 启动后，会调度到具有 `disktype=ssd` 标签的节点上。

# Affinity(亲和性)

**Affinity(亲和性)** 可以提供比 NodeSelector 更丰富的调度规则。亲和性主要就是只两个事物之间的 **Affinity(亲和)** 和 **Anti-Affinity(反亲和)。**当两个事物比较亲和，则更愿意在一起，比如 Pod A 和 Node A 亲和，则 Pod 会运行在 Node A 上；而当两个事物反亲和时，则不会在一起，比如 Pod A 和 Pod B 反亲和，则 Pod A 和 Pod B 不会在同一个 Node 上运行。

亲和性调度规则分为 Pod 和 Node 两类：

- **Node Affinity(节点亲和性)** # 与 NodeSelector 类似，但是提供了更灵活的规则
- **Pod Affinity(Pod 亲和性)**# 可以根据某个 Pod 与某个 Pod 的亲和关系来决定该 Pod 要运行在具有(或不具有)某个 Pod 的 Node 上。
  - Pod 亲和可以分为两部分：
    - **Inter-pod Affinity(Pod 之间亲和)**
    - **Inter-pod Anti-affinity(Pod 之间反亲和)**

每种调度规则有 hard 和 soft 两种强度。通过如下两个字段

- **requiredDuringSchedulingIgnoredDuringExecution # 硬。**表示必须遵守调度规则，如果节点不存在或者其他原因导致调度失败，则 Pod 不会创建
- **preferredDuringSchedulingIgnoredDuringExecution # 软。**表示倾向于遵守调度规则，如果节点不存在或其他原因导致调度失败，则 Pod 会在其他不满足调度规则的节点上创建。

Pod Manifest 的 `spec.affinity` 字段可以配置 亲和/反亲和 的调度规则。效果如下：

```yaml
apiVersion: v1
kind: Pod
metadata: ......
spec:
  affinity:
    # 节点亲和
    nodeAffinity:
      preferredDuringSchedulingIgnoredDuringExecution: ......具体的调度规则
      requiredDuringSchedulingIgnoredDuringExecution: ......具体的调度规则
    # Pod 亲和
    podAffinity:
      preferredDuringSchedulingIgnoredDuringExecution: ......具体的调度规则
      requiredDuringSchedulingIgnoredDuringExecution: ......具体的调度规则
    # Pod 反亲和
    PodAntiAffinity:
      preferredDuringSchedulingIgnoredDuringExecution: ......具体的调度规则
      requiredDuringSchedulingIgnoredDuringExecution: ......具体的调度规则
  containers:
    - name: XXX
```

注意：

- 亲和性只在 Pod 调度期间有效。比如修改或删除了 pod 所在节点的标签，Pod 不会被删除。 也就是说，当 Pod 已经运行时，亲和性是不起作用的。
- Affinity 的 Manifests 各字段详见 [Pod Manifest 详解](/docs/10.云原生/2.3.Kubernetes%20 容器编排系统/1.API、Resource(资源)、Object(对象)/API%20 参考/工作负载资源/Pod%20Manifest%20 详解.md 参考/工作负载资源/Pod Manifest 详解.md)

## Node Affinity

Node Affinity 概念上类似于 NodeSelector，属于 Pod 与 Node 之前的亲和性关系。可以根据节点上的标签来约束 Pod 可以调度到哪些节点上。

Affinity 的 Manifests 各字段详见 [Pod Manifest 详解](/docs/10.云原生/2.3.Kubernetes%20 容器编排系统/1.API、Resource(资源)、Object(对象)/API%20 参考/工作负载资源/Pod%20Manifest%20 详解.md 参考/工作负载资源/Pod Manifest 详解.md)

### 应用示例

让 Pod 调度到具有标签 `kubernetes.io/hostname` 的值为 `node-01.test.tjiptv.net` 或者 `node-02.test.tjiptv.net` 的 Node 上

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: with-node-affinity
spec:
  affinity:
    nodeAffinity:
      requiredDuringSchedulingIgnoredDuringExecution:
        nodeSelectorTerms:
          - matchExpressions:
              - key: kubernetes.io/e2e-az-name
                operator: In
                values:
                  - e2e-az1
                  - e2e-az2
      preferredDuringSchedulingIgnoredDuringExecution:
        - weight: 1
          preference:
            matchExpressions:
              - key: another-node-label-key
                operator: In
                values:
                  - another-node-label-value
  containers:
    - name: with-node-affinity
      image: k8s.gcr.io/pause:2.0
```

上面这个示例中的亲和性调度规则分为两部分：

- requiredDuringSchedulingIgnoredDuringExecution 字段下的规则为：Pod 只能被调度到具有 `kubernetes.io/e2e-az-name` 标签，且标签值为 `e2e-az1` 或 `e2e-az2` 的节点上。
  - 除了满足 requiredDuringSchedulingIgnoredDuringExecution 字段的规则外，也就是说，匹配到的这些节点中，还会根据下面的规则，进一步约束 Pod 被调度的节点
- preferredDuringSchedulingIgnoredDuringExecution 字段下的规则为：具有 `another-node-label-key` 标签，且标签值为 `another-node-label-value` 的节点应该优先调度。除非这些节点不可用或由于其他原因不可调度，那么才会选择其他节点。

从上面的例子可以看到，想要匹配到想要的 Node，还是基于 [Selector(选择器)](Label%20and%20Selector(标签和选择器).md and Selector(标签和选择器).md) 来实现的，其中 `operator` 字段中，可以通过设定 NotIn 或者 DoesNotExist 来实现**节点反亲和**的效果。

选择器规则：

- 如果你同时指定了 `nodeSelector` 和 `nodeAffinity`，\_两者\_必须都要满足， 才能将 Pod 调度到候选节点上。
- 如果你指定了多个与 `nodeAffinity` 类型关联的 `nodeSelectorTerms`，则 **如果其中一个** `nodeSelectorTerms` 满足的话，pod 将可以调度到节点上。
- 如果你指定了多个与 `nodeSelectorTerms` 关联的 `matchExpressions`，则 **只有当所有** `matchExpressions` 满足的话，Pod 才会可以调度到节点上。

`weight` 字段的值的范围是 1-100。对于每个符合所有调度要求（资源请求、RequiredDuringScheduling 亲和性表达式等） 的节点，调度器将遍历该字段的元素来计算总和，并且如果节点匹配对应的 MatchExpressions，则添加“权重”到总和。 然后将这个评分与该节点的其他优先级函数的评分进行组合。 总分最高的节点是最优选的

## Pod Affinity

Pod Affinity 表示 Pod 间亲和性与反亲和性，使我们可以基于已经在节点上运行的 Pod 的标签来约束待调度的 Pod 可以调度到的节点，而不是基于节点上的标签。

这种规则可以这么描述：**“如果 X 节点上已经运行了一个或多个满足规则 Y 的 Pod， 则这个 Pod 应该(在反亲和情况下不应该)运行在 X 节点”**

- **Y** 表示一个具有可选的关联命令空间列表的 [**LabelSelector(标签选择器)**](Label%20and%20Selector(标签和选择器).md and Selector(标签和选择器).md)； 与节点不同，因为 Pod 是命名空间限定的(因此 Pod 上的标签也是命名空间限定的)， 因此作用于 Pod 标签的标签选择算符必须指定选择算符应用在哪个命名空间。
- **X** 我们使用 `topologyKey` 字段来表示 X，`topologyKey` 用来指定节点标签的键，用来表示 **Topology Domain(拓扑域)**。 从概念上讲，k8s 的节点、机架、云供应商可用区、云供应商地理区域等，都有区域的概念，都可以称为拓扑域。说白了，`topologyKey` 字段还是用来筛选 Node 的。

> 说明：
>
> - Pod 间亲和性与反亲和性需要大量的处理，这可能会显著减慢大规模集群中的调度。 我们不建议在超过数百个节点的集群中使用它们。
> - Pod 反亲和性需要对节点进行一致的标记，即集群中的每个节点必须具有适当的标签能够匹配 `topologyKey`。如果某些或所有节点缺少指定的 `topologyKey` 标签，可能会导致意外行为。

Pod 间亲和性通过 PodSpec 中 affinity 字段下的 podAffinity 字段进行指定。 而 Pod 间反亲和性通过 PodSpec 中 affinity 字段下的 podAntiAffinity 字段进行指定。

### 应用示例

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: with-pod-affinity
spec:
  affinity:
    podAffinity:
      requiredDuringSchedulingIgnoredDuringExecution:
        - labelSelector:
            matchExpressions:
              - key: security
                operator: In
                values:
                  - S1
          topologyKey: topology.kubernetes.io/zone
    podAntiAffinity:
      preferredDuringSchedulingIgnoredDuringExecution:
        - weight: 100
          podAffinityTerm:
            labelSelector:
              matchExpressions:
                - key: security
                  operator: In
                  values:
                    - S2
            topologyKey: topology.kubernetes.io/zone
  containers:
    - name: with-pod-affinity
      image: k8s.gcr.io/pause:2.0
```

在此示例中，在这个 Pod 的亲和性配置定义了两条规则：

- Pod 亲和性规则
  - 仅当节点和至少一个已运行且有键为“security”且值为“S1”的标签 的 Pod 处于同一区域时，才可以将该 Pod 调度到节点上。(更确切的说，如果节点 N 具有 `topology.kubernetes.io/zone` 标签， 则 Pod 有资格在节点 N 上运行，以便集群中至少有一个具有 `topology.kubernetes.io/zone` 标签的节点正在运行具有键“security”和值 “S1”的标签的 pod。)
- Pod 反亲和性规则
  - 如果在具有 `topology.kubernetes.io/zone` 标签的节点上已经运行了具有键为 `security` 和值 `S2` 标签的 Pod，则同样具有相同标签的 Pod 不能被调度到该节点。
    - 也就是说，具有 security=S2 标签的 Pod 在具有 topology.kubernetes.io/zone 标签的节点上，只能运行一个。

Pod 亲和性与反亲和性的合法 `operator` 字段的值有 In，NotIn，Exists，DoesNotExist。

既然 `topologyKey` 是拓扑域，那 Pod 之间怎样才是属于同一个拓扑域？

- 如果使用  kubernetes.io/hostname 标签，则表示拓扑域为 Node 范围，那么  kubernetes.io/hostname  对应的值不一样就是不同的拓扑域。比如 Pod1 在  kubernetes.io/hostname=node1  的 Node 上，Pod2 在  kubernetes.io/hostname=node2  的 Node 上，Pod3 在  kubernetes.io/hostname=node1  的 Node 上，则 Pod1 和 Pod3 在同一个拓扑域，而 Pod2 则与 Pod1 和 Pod3 不在同一个拓扑域。
- 如果使用  failure-domain.kubernetes.io/zone ，则表示拓扑域为一个区域。同样，Node 的标签  failure-domain.kubernetes.io/zone  对应的值不一样也不是同一个拓扑域，比如 Pod1 在  failure-domain.kubernetes.io/zone=beijing  的 Node 上，Pod2 在  failure-domain.kubernetes.io/zone=hangzhou  的 Node 上，则 Pod1 和 Pod2 不属于同一个拓扑域。
- 当然，topologyKey 也可以使用自定义标签。比如可以给一组 Node 打上标签  custom_topology，那么拓扑域就是针对这个标签了，则该标签相同的 Node 上的 Pod 属于同一个拓扑域。

原则上，topologyKey 可以是任何合法的标签键。 然而，出于性能和安全原因，topologyKey 受到一些限制：

- 在 requiredDuringSchedulingIgnoredDuringExecution 和 preferredDuringSchedulingIgnoredDuringExecution 中，topologyKey 不允许为空。
- 对于 `podAntiAffinity.requiredDuringSchedulingIgnoredDuringExecution`， 准入控制器 `LimitPodHardAntiAffinityTopology` 被引入以确保 `topologyKey` 只能是 `kubernetes.io/hostname`。如果你希望 topologyKey 也可用于其他定制拓扑逻辑，你可以更改准入控制器或者禁用之。
- 除上述情况外，topologyKey 可以是任何合法的标签键。

除了 labelSelector 和 topologyKey，你也可以指定表示命名空间的 namespaces 队列，labelSelector 也应该匹配它 （这个与 labelSelector 和 topologyKey 的定义位于相同的级别）。 如果忽略或者为空，则默认为 Pod 亲和性/反亲和性的定义所在的命名空间。

所有与 requiredDuringSchedulingIgnoredDuringExecution 亲和性与反亲和性 关联的 matchExpressions 必须满足，才能将 pod 调度到节点上。

### 更实际的用例

Pod 间亲和性与反亲和性在与更高级别的集合(例如 ReplicaSets、StatefulSets、 Deployments 等)一起使用时，它们更加有用。 可以轻松配置一组应位于相同定义拓扑（例如，节点）中的工作负载。

在三节点集群中，一个 web 应用程序具有内存缓存(例如 redis)。 我们希望 web 服务器尽可能与缓存放置在同一位置，并且多个 Redis 不在同一个节点上。

#### 永远放置在不同节点上

下面的例子使用 PodAntiAffinity 规则和 topologyKey: "kubernetes.io/hostname" 来部署 redis 集群以便在同一主机上没有两个实例。 它有三个副本和选择器标签 app=store。 Deployment 配置了 PodAntiAffinity，用来确保调度器不会将副本调度到单个节点上。

> [ZooKeeper 教程](https://kubernetes.io/zh/docs/tutorials/stateful-application/zookeeper/#tolerating-node-failure) 中使用了相同的反亲和性配置方式， 来达到高可用性的 StatefulSet 的样例。

这一段反亲和性可以描述为：具有 `kubernetes.io/hostname` 标签的节点上只能运行一个具有 `app=store` 标签的 Pod。这就保证了，3 个 Redis 将会分散到不同的节点。

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: redis-cache
spec:
  selector:
    matchLabels:
      app: store
  replicas: 3
  template:
    metadata:
      labels:
        app: store
    spec:
      affinity:
        podAntiAffinity:
          requiredDuringSchedulingIgnoredDuringExecution:
            - labelSelector:
                matchExpressions:
                  - key: app
                    operator: In
                    values:
                      - store
              topologyKey: "kubernetes.io/hostname"
      containers:
        - name: redis-server
          image: redis:3.2-alpine
```

#### 始终放置在相同节点上

下面 webserver Deployment 的 YAML 代码段中配置了 podAntiAffinity 和 podAffinity。 这将通知调度器将它的所有副本与具有 app=store 选择器标签的 Pod 放置在一起。 这还确保每个 web 服务器副本不会调度到单个节点上。

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: web-server
spec:
  selector:
    matchLabels:
      app: web-store
  replicas: 3
  template:
    metadata:
      labels:
        app: web-store
    spec:
      affinity:
        podAntiAffinity:
          requiredDuringSchedulingIgnoredDuringExecution:
            - labelSelector:
                matchExpressions:
                  - key: app
                    operator: In
                    values:
                      - web-store
              topologyKey: "kubernetes.io/hostname"
        podAffinity:
          requiredDuringSchedulingIgnoredDuringExecution:
            - labelSelector:
                matchExpressions:
                  - key: app
                    operator: In
                    values:
                      - store
              topologyKey: "kubernetes.io/hostname"
      containers:
        - name: web-app
          image: nginx:1.16-alpine
```

如果我们创建了上面的两个 Deployment，我们的三节点集群将如下表所示。

| node-1        | node-2        | node-3        |
| ------------- | ------------- | ------------- |
| _webserver-1_ | _webserver-2_ | _webserver-3_ |
| _cache-1_     | _cache-2_     | _cache-3_     |

可以看到，Web 服务器的 3 个副本，逐一与缓存共存，并且每个节点只有一个 Web 服务器和一个缓存。

```bash
kubectl get pods -o wide
```

输出类似于如下内容：

```bash
NAME                           READY     STATUS    RESTARTS   AGE       IP           NODE
redis-cache-1450370735-6dzlj   1/1       Running   0          8m        10.192.4.2   kube-node-3
redis-cache-1450370735-j2j96   1/1       Running   0          8m        10.192.2.2   kube-node-1
redis-cache-1450370735-z73mh   1/1       Running   0          8m        10.192.3.1   kube-node-2
web-server-1287567482-5d4dz    1/1       Running   0          7m        10.192.2.3   kube-node-1
web-server-1287567482-6f7v5    1/1       Running   0          7m        10.192.4.3   kube-node-3
web-server-1287567482-s330j    1/1       Running   0          7m        10.192.3.2   kube-node-2
```

# Taint(污点) 与 Toleration(容忍度)

该方式可以由 Node 来决定是否拒绝运行 Pod。把 Node 当作第一人称“我”来描述，给 Node 打上一个 Taint，我会拒绝不容忍我的 Taint 的 Pod 运行在我身上；还有一种情况是如果 Pod 已经运行在我身上了，那么这时我多了一个 Pod 不容忍的 Taint，那么如果我身上的 Pod 不容忍，我就要驱逐这些 Pod。

通常情况下，Taint 定义在 Node 上，Toleration 定义在 Pod 上。

## Taint

**taints([]Object)** # 定义 Node 污点。凡是具有污点的节点，都对不容忍污点的 Pod 具有某些 effect(效果)。

- **effect <STRING> # 必须的。**定义当 Pod 不能容忍 Node 的污点的时候，Node 对 Pod 的排斥效果（效果即使要采取的行为是什么)。
  - NoSchedule # 仅影响调度过程(仅影响调度过程，对于已经调度到该 Node 上的 Pod 则没效果，即调度完成后再给 Node 上加的污点就算 Pod 不容忍也没影响)；
  - NoExecut # 即影响调度过程，也影响现存的 Pod(即不容忍 taint 的 Pod 不但不会调度到 Node 上，如果不容忍该 Taint 的 Pod 已经在该 Node 上，会被 Node 驱逐)；
    - 在什么情况下会触发驱逐 Pod 的效果呢？e.g.在 Pod 已经调度到 Node 上之后，再给 Node 添加 Taint，这时候该 Pod 不容忍新增加的污点，那么 Node 就会驱逐该 Pod
  - PreferNoschedule # 最好不调度，实在不行了还可以调度。
- **key: <STRING> # 必须的。**
- **value: <STRING>** #
- **timeAdded: <STRING>**#

可以直接使用 kubectl 命令为节点添加一个污点，比如：

```bash
kubectl taint nodes NODE KEY[=VALUE]:EFFECT
```

### 应用示例

- 为 node-t.tj-test 这个节点添加污点，污点的 key 为 node-role.kubernetes.io/master，没有 value，污点效果为 NodSchedule
  - kubectl taint node node-4.tj-test node-role.kubernetes.io/master:NoSchedule
- 删除所有节点上 node-role.kubernetes.io/master 这个污点。也就是让 master 可以当作 node 使用
  - kubectl taint nodes --all node-role.kubernetes.io/master-
- 该命令可以获取所有节点上的污点
  - kubectl get nodes --template='{{ range .items }}{{ .metadata.name }}{{"\t"}}{{ .spec.taints }}{{"\n"}}{{end}}'

## Toleration

**tolerations([]Object)** # 为 Pod 添加容忍规则

- **effect: <STRING>**# 表明匹配的 taint 的 effect 字段，包括三个效果 NoSchedule, PreferNoSchedule and NoExecute，如果不指定该字段则匹配所有污点效果
- **key: <STRING>** # 指明要容忍的 taint 中的 key
- **operator: <STRING>** # 定义容忍要满足的条件
  - Exists # 只要 Key 一样，则容忍
  - Equal # 需要 Key 和 value 都一样才容忍
- **tolerationSeconds: <INTEGER>** # 定义容忍时间
- **value: <STRING>** # 指明要容忍的 taint 中的 value

### 应用示例

仅需下面两行，就表示容忍所有污点

```yaml
tolerations:
  - operator: Exists
```

## Taint based Evictions(基于污点的驱逐)

前文提到过污点的 effect 值 `NoExecute`会影响已经在节点上运行的 Pod

- 如果 Pod 不能忍受 effect 值为 `NoExecute` 的污点，那么 Pod 将马上被驱逐
- 如果 Pod 能够忍受 effect 值为 `NoExecute` 的污点，但是在容忍度定义中没有指定 `tolerationSeconds`，则 Pod 还会一直在这个节点上运行。
- 如果 Pod 能够忍受 effect 值为 `NoExecute` 的污点，而且指定了 `tolerationSeconds`， 则 Pod 还能在这个节点上继续运行这个指定的时间长度。

当某种条件为真时，节点控制器会自动给节点添加一个污点。当前内置的污点包括：

- `node.kubernetes.io/not-ready`：节点未准备好。这相当于节点状态 `Ready` 的值为 "`False`"。
- `node.kubernetes.io/unreachable`：节点控制器访问不到节点. 这相当于节点状态 `Ready` 的值为 "`Unknown`"。
- `node.kubernetes.io/out-of-disk`：节点磁盘耗尽。
- `node.kubernetes.io/memory-pressure`：节点存在内存压力。
- `node.kubernetes.io/disk-pressure`：节点存在磁盘压力。
- `node.kubernetes.io/network-unavailable`：节点网络不可用。
- `node.kubernetes.io/unschedulable`: 节点不可调度。
- `node.cloudprovider.kubernetes.io/uninitialized`：如果 kubelet 启动时指定了一个 "外部" 云平台驱动， 它将给当前节点添加一个污点将其标志为不可用。在 cloud-controller-manager 的一个控制器初始化这个节点后，kubelet 将删除这个污点。

在节点被驱逐时，节点控制器或者 kubelet 会添加带有 `NoExecute` 效应的相关污点。 如果异常状态恢复正常，kubelet 或节点控制器能够移除相关的污点。

> **说明：** 为了保证由于节点问题引起的 Pod 驱逐 [速率限制](https://kubernetes.io/zh/docs/concepts/architecture/nodes/)行为正常， 系统实际上会以限定速率的方式添加污点。在像主控节点与工作节点间通信中断等场景下， 这样做可以避免 Pod 被大量驱逐。

使用这个功能特性，结合 `tolerationSeconds`，Pod 就可以指定当节点出现一个 或全部上述问题时还将在这个节点上运行多长的时间。
比如，一个使用了很多本地状态的应用程序在网络断开时，仍然希望停留在当前节点上运行一段较长的时间， 愿意等待网络恢复以避免被驱逐。在这种情况下，Pod 的容忍度可能是下面这样的：

    tolerations:
    - key: "node.kubernetes.io/unreachable"
      operator: "Exists"
      effect: "NoExecute"
      tolerationSeconds: 6000

> **说明：**
> Kubernetes 会自动给 Pod 添加一个 key 为 `node.kubernetes.io/not-ready` 的容忍度 并配置 `tolerationSeconds=300`，除非用户提供的 Pod 配置中已经已存在了 key 为 `node.kubernetes.io/not-ready` 的容忍度。
> 同样，Kubernetes 会给 Pod 添加一个 key 为 `node.kubernetes.io/unreachable` 的容忍度 并配置 `tolerationSeconds=300`，除非用户提供的 Pod 配置中已经已存在了 key 为 `node.kubernetes.io/unreachable` 的容忍度。

这种自动添加的容忍度意味着在其中一种问题被检测到时 Pod 默认能够继续停留在当前节点运行 5 分钟。
[DaemonSet](https://kubernetes.io/zh/docs/concepts/workloads/controllers/daemonset/) 中的 Pod 被创建时， 针对以下污点自动添加的 `NoExecute` 的容忍度将不会指定 `tolerationSeconds`：

- `node.kubernetes.io/unreachable`
- `node.kubernetes.io/not-ready`

这保证了出现上述问题时 DaemonSet 中的 Pod 永远不会被驱逐。

## 基于节点状态添加污点

Node 生命周期控制器会自动创建与 Node 条件相对应的带有 `NoSchedule` 效应的污点。 同样，调度器不检查节点条件，而是检查节点污点。这确保了节点条件不会影响调度到节点上的内容。 用户可以通过添加适当的 Pod 容忍度来选择忽略某些 Node 的问题(表示为 Node 的调度条件)。
DaemonSet 控制器自动为所有守护进程添加如下 `NoSchedule` 容忍度以防 DaemonSet 崩溃：

- `node.kubernetes.io/memory-pressure`
- `node.kubernetes.io/disk-pressure`
- `node.kubernetes.io/out-of-disk` (_只适合关键 Pod_)
- `node.kubernetes.io/unschedulable` (1.10 或更高版本)
- `node.kubernetes.io/network-unavailable` (_只适合主机网络配置_)

添加上述容忍度确保了向后兼容，您也可以选择自由向 DaemonSet 添加容忍度。
