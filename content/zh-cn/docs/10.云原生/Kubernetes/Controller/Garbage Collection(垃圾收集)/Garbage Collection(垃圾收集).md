---
title: Garbage Collection(垃圾收集)
linkTitle: Garbage Collection(垃圾收集)
weight: 1
---

# 概述

> 参考：
>
> - [官方文档，概念-Kubernetes 架构-垃圾收集](https://kubernetes.io/docs/concepts/architecture/garbage-collection/)
> - [宋净超-云原生资料库-Kubernetes 基础教程，集群资源管理-垃圾收集](https://lib.jimmysong.io/kubernetes-handbook/concepts/garbage-collection/)

**Garbage Collection(垃圾收集)** 功能用来删除曾经拥有 owner(拥有者) 但不再拥有 owner 的某些对象。

比如张三拥有 100 块钱，则张三就是这 100 块钱的 owner，当张三死亡后，那么这 100 块钱则不再具有 owner。

注意：垃圾收集是 beta 特性，在 Kubernetes 1.4 及以上版本默认启用。

## Garbage Collector 垃圾收集器

Garbage Collector 是 k8s 垃圾收集功能的具体实现。Garbage Collector 属于 Kubernetes Controller 的一部分。kube-controller-manager 的 `--controllers` 命令行标志的值中包含 garbagecollector，用以控制是否启用该控制器。

垃圾收集器在 kubernetes 的代码中的位置在这里：[pkg/controller/garbagecollector](https://github.com/kubernetes/kubernetes/tree/master/pkg/controller/garbagecollector)。

# Owner 和 Dependent

> 参考：
>
> - [官方文档，概念-概述-使用 Kubernetes 对象-属主与从属](https://kubernetes.io/docs/concepts/overview/working-with-objects/owners-dependents/)

在 Kubernetes 中

- **Owner(拥有者/属主)** # 一些 Kubernetes 对象是其它对象的 Owner。例如，一个 Deployment 是一组 Pod 的 Owner(i.e.deployment 拥有这些 pod)。
- **Dependent(依赖他人者/从属)** # 被 Owner 拥有的对象被称 Dependent(从属) 对象，该 Dependent 属于该 Owner。

每个 Dependent 对象具有一个指向其 owner 对象的 metadata.ownerReferences 字段。

有时，Kubernetes 会自动设置 ownerReference 的值。例如，当创建一个 ReplicaSet 时，Kubernetes 自动设置 ReplicaSet 中每个 Pod 的 ownerReference 字段值。在 1.6 版本，Kubernetes 会自动为一些对象设置 ownerReference 的值，这些对象是由 ReplicationController、ReplicaSet、StatefulSet、DaemonSet 和 Deployment 所创建或管理。

也可以通过手动设置 ownerReference 的值，来指定 Owner 和 Dependent 之间的关系。

下面的配置文件，表示一个具有 3 个 Pod 的 ReplicaSet：

```yaml
apiVersion: extensions/v1beta1
kind: ReplicaSet
metadata:
  name: my-repset
spec:
  replicas: 3
  selector:
    matchLabels:
      pod-is-for: garbage-collection-example
  template:
    metadata:
      labels:
        pod-is-for: garbage-collection-example
    spec:
      containers:
        - name: nginx
          image: nginx
```

如果创建该 ReplicaSet，然后查看 Pod 的 metadata 字段，能够看到 OwnerReferences 字段：

    kubectl create -f https://k8s.io/docs/concepts/abstractions/controllers/my-repset.yaml
    kubectl get pods --output=yaml

输出显示了 Pod 的 Owner 是名为 my-repset 的 ReplicaSet：

    apiVersion: v1
    kind: Pod
    metadata:
      ...
      ownerReferences:
      - apiVersion: extensions/v1beta1
        controller: true
        blockOwnerDeletion: true
        kind: ReplicaSet
        name: my-repset
        uid: d9607e19-f88f-11e6-a518-42010a800195
      ...

# 垃圾收集器删除 Dependent 的方式

当删除对象时，可以指定是否该对象的 Dependent 也自动删除掉。自动删除 Dependent 的行为称为 **Cascading Deletion(级联删除)**。Kubernetes 中有两种级联删除的模式：**Background(后台) 模式**和 **Foreground(前台) 模式**。

如果删除对象时，不自动删除它的 Dependent，这些 Dependent 被称为 orphaned(孤立) 对象。

## Foreground(前台) 级联删除

在 foreground 级联删除 模式下，根对象首先进入 “deletion in progress(删除中)” 状态。在该状态时会有如下的情况：

- 对象仍然可以通过 REST API 可见
- 会设置对象的 deletionTimestamp 字段
- 对象的 metadata.finalizers 字段包含了值 foregroundDeletion

一旦被设置为 deletion in progress 状态，垃圾收集器会删除对象的所有 Dependent。垃圾收集器删除了所有 “Blocking” 的 Dependent（对象的 ownerReference.blockOwnerDeletion=true）之后，它会删除 Owner 对象。

注意，在 “foreground 删除” 模式下，Dependent 只有通过 ownerReference.blockOwnerDeletion 才能阻止删除 Owner 对象。在 Kubernetes 1.7 版本中将增加 admission controller，基于 Owner 对象上的删除权限来控制用户去设置 blockOwnerDeletion 的值为 true，所以未授权的 Dependent 不能够延迟 Owner 对象的删除。

如果一个对象的 ownerReferences 字段被一个 Controller（例如 Deployment 或 ReplicaSet）设置，blockOwnerDeletion 会被自动设置，没必要手动修改这个字段。

## Background(后台) 级联删除

在 background 级联删除 模式下，Kubernetes 会立即删除 Owner 对象，然后垃圾收集器会在后台删除这些 Dependent。

## 设置级联删除策略

通过为 Owner 对象设置 deleteOptions.propagationPolicy 字段，可以控制级联删除策略。可能的取值包括：“orphan”、“Foreground” 或 “Background”。

对很多 Controller 资源，包括 ReplicationController、ReplicaSet、StatefulSet、DaemonSet 和 Deployment，默认的垃圾收集策略是 orphan。因此，除非指定其它的垃圾收集策略，否则所有 Dependent 对象使用的都是 orphan 策略。

注意：本段所指的默认值是指 REST API 的默认值，并非 kubectl 命令的默认值，kubectl 默认为级联删除，后面会讲到。

下面是一个在后台删除 Dependent 对象的例子：

    kubectl proxy --port=8080
    curl -X DELETE localhost:8080/apis/extensions/v1beta1/namespaces/default/replicasets/my-repset \
    -d '{"kind":"DeleteOptions","apiVersion":"v1","propagationPolicy":"Background"}' \
    -H "Content-Type: application/json"

下面是一个在前台删除 Dependent 对象的例子：

    kubectl proxy --port=8080
    curl -X DELETE localhost:8080/apis/extensions/v1beta1/namespaces/default/replicasets/my-repset \
    -d '{"kind":"DeleteOptions","apiVersion":"v1","propagationPolicy":"Foreground"}' \
    -H "Content-Type: application/json"

下面是一个孤儿 Dependent 的例子：

    kubectl proxy --port=8080
    curl -X DELETE localhost:8080/apis/extensions/v1beta1/namespaces/default/replicasets/my-repset \
    -d '{"kind":"DeleteOptions","apiVersion":"v1","propagationPolicy":"Orphan"}' \
    -H "Content-Type: application/json"

kubectl 也支持级联删除。 通过设置 --cascade 为 true，可以使用 kubectl 自动删除 Dependent 对象。设置 --cascade 为 false，会使 Dependent 对象成为孤儿 Dependent 对象。--cascade 的默认值是 true。

下面是一个例子，使一个 ReplicaSet 的 Dependent 对象成为孤立 Dependent：

    kubectl delete replicaset my-repset --cascade=false
