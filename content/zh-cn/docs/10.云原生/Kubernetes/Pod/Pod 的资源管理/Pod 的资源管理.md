---
title: Pod 的资源管理
linkTitle: Pod 的资源管理
weight: 1
---

# 概述

> 参考：
>
> - [官方文档，概念 - 配置 - Pod 和容器的资源管理](https://kubernetes.io/docs/concepts/configuration/manage-resources-containers/)

# Pod 中 Container 的资源需求与资源限制

可以在 Pod 的 yaml 中定义该 Pod 中各个 Container 对内存与 CPU 的最低需求量和最大使用量

- requests：资源需求，最低保障，资源最少需要多少
- limits：限制，硬限制，限额，资源最大不能超过多少

当对 Container 进行资源制定后，会出现 QoS(服务质量)的属性，下列 3 个属性从上往下优先级下降；当节点资源不够时，优先级越高，越会保证其正常运行，其余不够提供资源的 Container 则不再运行

- Guarateed：有保证的，Pod 中每个 Container 同时设置 CPU 和内存的 requests 和 limits，且 request 和 limits 的值相同
- Burstable：超频，Pod 中至少有一个 Container 设置了 CPU 或内存资源的 requests 属性
- BestEffort：尽力努力(尽力而为)没有任何一个 Container 设置了 requests 和 limits 属性

![image.jpeg](https://notes-learning.oss-cn-beijing.aliyuncs.com/xu9mbw/1617283210163-cdf33748-8346-40d3-96b1-c3da2cf90df5.jpeg)

关于在 yaml 中如何写资源限制中数值的说明：

- kubernetes 中的一个 CPU 是一个逻辑 CPU，1CPU 的核心数=1000millicores 毫核心(也就是说 500m 相当于 0.5 个 CPU)
- 限制 cpu 可以用整数写，1 就是 1 个 cpu、0.5 就是 0.5 个 cpu。也可以带单位，1000m 就是 1 个 cpu，500m 就是 0.5 个 cpu。
- 限制内存需要使用单位(IEC 标准、公制标准都可以)，即 Mi、Gi 或者 M、G 等。例如 1024Mi、1Gi 等

Note：polinux/stress 这是一个非常好用的压测容器，可以对容器指定其所使用的内存和 cpu 等资源的大小。当创建完资源配合等资源限制的对象后，可以通过该容器来测试资源限制是否生效。

分配内存资源给 容器 和 Pods

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: memory-demo
  namespace: mem-example
spec:
  containers:
    - name: memory-demo-ctr
      image: polinux/stress
      resources:
        limits:
          memory: "200Mi"
        requests:
          memory: "100Mi"
      command: ["stress"]
      args: ["--vm", "1", "--vm-bytes", "150M", "--vm-hang", "1"]
```

分配 CPU 资源给 容器 和 Pods

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: cpu-demo
  namespace: cpu-example
spec:
  containers:
    - name: cpu-demo-ctr
      image: vish/stress
      resources:
        limits:
          cpu: "1"
        requests:
          cpu: "0.5"
      args:
        - -cpus
        - "2"
```
