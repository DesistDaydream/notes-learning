---
title: Quota(配额)
---

# 概述

> 参考：
>
> - [官方文档,概念-策略-LimitRange](https://kubernetes.io/docs/concepts/policy/limit-range/)
> - [官方文档,概念-策略-ResourceQuotas](https://kubernetes.io/docs/concepts/policy/resource-quotas/)

# namespace 中的资源配额

官方文档：<https://kubernetes.io/docs/tasks/administer-cluster/manage-resources/quota-memory-cpu-namespace/>

当多个团队或者用户共用同一个集群的时候难免会有资源竞争的情况发生，这时候就需要对不同团队或用户的资源使用配额做出限制。比如，不同团队使用不同的 namespace，然后给该 namespace 进行资源限制即可

目前有两种 k8s 对象分配管理相关的控制策略

## LimitRange(限制范围)

设定 pod 等对象的默认资源消耗以及可以消耗的资源范围

官方文档：

- 概念：<https://kubernetes.io/docs/concepts/policy/limit-range/>
- 用法：
  - <https://kubernetes.io/docs/tasks/administer-cluster/manage-resources/cpu-constraint-namespace/>
  - <https://kubernetes.io/docs/tasks/administer-cluster/manage-resources/memory-constraint-namespace/>
  - <https://kubernetes.io/docs/tasks/administer-cluster/manage-resources/cpu-default-namespace/>
  - <https://kubernetes.io/docs/tasks/administer-cluster/manage-resources/memory-default-namespace/>
  - .....等等

## ResourceQuota(资源配额)

基于 namespace，限制该 namesapce 下的总体资源的创建和消耗

官方文档：

- 概念：<https://kubernetes.io/docs/concepts/policy/resource-quotas/>
- 用法：
  - <https://kubernetes.io/docs/tasks/administer-cluster/manage-resources/quota-memory-cpu-namespace/>
  - <https://kubernetes.io/docs/tasks/administer-cluster/manage-resources/quota-pod-namespace/>
  - <https://kubernetes.io/docs/tasks/administer-cluster/quota-api-object/> # 为指定的 API 对象设置 resourceQuota

资源配额分为三种类型：

- 计算资源配额
- 存储资源配额
- 对象数量配额

## 总结

- 仅设置 ResourceQuota 时，如果不再 pod 上设置资源的需求和限制，则无法成功创建 pod，需要配合 LimitRange 设置 pod 的默认需求和限制，才可成功创建 pod
- 两种控制策略的作用范围都是对于某一 namespace
  - ResourceQuota 用来限制 namespace 中所有的 Pod 占用的总的资源 request 和 limit
  - LimitRange 是用来设置 namespace 中 Pod 的默认的资源 request 和 limit 值，还有，Pod 的可用资源的 request 和 limit 值的最大与最小值。

# 简单的应用示例

Note：polinux/stress 这是一个非常好用的压测容器，可以对容器指定其所使用的内存和 cpu 等资源的大小。当创建完资源配合等资源限制的对象后，可以通过该容器来测试资源限制是否生效。

## 配置计算资源配额

为 test 名称空间分配了如下配合，最多能建立 2 个 pod，最多 request 的 cpu 数量为 2 个，内存为 10G，最多 limit 的 cpu 数量为 4 个，内存为 20G

```yaml
apiVersion: v1
kind: ResourceQuota
metadata:
  name: compute-resources
  namespace: test
spec:
  hard:
    pods: "2"
    requests.cpu: "2"
    requests.memory: 10Gi
    limits.cpu: "4"
    limits.memory: 20Gi
```

## 配置 API 对象数量限制

```yaml
apiVersion: v1
kind: ResourceQuota
metadata:
  name: object-counts
  namespace: test
spec:
  hard:
    configmaps: "1"
    persistentvolumeclaims: "4"
    replicationcontrollers: "2"
    secrets: "10"
    services: "10"
    services.loadbalancers: "2"
```

## 配置 CPU 和内存 LimitRange

test 名称空间下的 pod 启动后，默认 request 的 cpu 为 0.5，内存为 256M，默认 limit 的 cpu 为 1，内存为 512M

```yaml
apiVersion: v1
kind: LimitRange
metadata:
  name: limit-range
  namespace: test
spec:
  limits:
    - default:
        memory: 512Mi
        cpu: 1
      defaultRequest:
        memory: 256Mi
        cpu: 0.5
      type: Container
```

Note:

- default 即 limit 的值
- defaultRequest 即 request 的值

在 limits 字段下还有其他的可用字段如下：

- max 代表 limit 的最大值
- min 代表 request 的最小值
- maxLimitRequestRatio 代表 limit / request 的最大值。由于节点是根据 pod request 调度资源，可以做到节点超卖，maxLimitRequestRatio 代表 pod 最大超卖比例。
