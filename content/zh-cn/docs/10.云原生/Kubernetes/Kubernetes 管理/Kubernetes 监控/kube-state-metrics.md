---
title: kube-state-metrics
---

# 概述

> 参考：
> - [GitHub 项目，kubernetes/kube-stsate-metrics](https://github.com/kubernetes/kube-state-metrics)
> - [GitHub 文档,可暴露的所有指标列表](https://github.com/kubernetes/kube-state-metrics/tree/master/docs)

已经有了 cadvisor、Metric Server，几乎容器运行的所有指标都能拿到，但是下面这种情况却无能为力：

- 我调度了多少个 replicas？现在可用的有几个？
- 多少个 Pod 是 running/stopped/terminated 状态？
- Pod 重启了多少次？
- 我有多少 job 在运行中

而这些则是 kube-state-metrics 提供的内容，它基于 client-go 开发，轮询 Kubernetes API，并将 Kubernetes 的结构化信息转换为 Metrics。

kube-state-metrics 提供的指标，按照阶段分为三种类别：

- 1.实验性质的：k8s api 中 alpha 阶段的或者 spec 的字段。
- 2.稳定版本的：k8s 中不向后兼容的主要版本的更新
- 3.被废弃的：已经不在维护的。

指标类别包括：

- [CertificateSigningRequest Metrics](https://github.com/kubernetes/kube-state-metrics/blob/master/docs/certificatessigningrequest-metrics.md)
- [ConfigMap Metrics](https://github.com/kubernetes/kube-state-metrics/blob/master/docs/configmap-metrics.md)
- [CronJob Metrics](https://github.com/kubernetes/kube-state-metrics/blob/master/docs/cronjob-metrics.md)
- [DaemonSet Metrics](https://github.com/kubernetes/kube-state-metrics/blob/master/docs/daemonset-metrics.md)
- [Deployment Metrics](https://github.com/kubernetes/kube-state-metrics/blob/master/docs/deployment-metrics.md)
- [Endpoint Metrics](https://github.com/kubernetes/kube-state-metrics/blob/master/docs/endpoint-metrics.md)
- [Horizontal Pod Autoscaler Metrics](https://github.com/kubernetes/kube-state-metrics/blob/master/docs/horizontalpodautoscaler-metrics.md)
- [Ingress Metrics](https://github.com/kubernetes/kube-state-metrics/blob/master/docs/ingress-metrics.md)
- [Job Metrics](https://github.com/kubernetes/kube-state-metrics/blob/master/docs/job-metrics.md)
- [Lease Metrics](https://github.com/kubernetes/kube-state-metrics/blob/master/docs/lease-metrics.md)
- [LimitRange Metrics](https://github.com/kubernetes/kube-state-metrics/blob/master/docs/limitrange-metrics.md)
- [MutatingWebhookConfiguration Metrics](https://github.com/kubernetes/kube-state-metrics/blob/master/docs/mutatingwebhookconfiguration-metrics.md)
- [Namespace Metrics](https://github.com/kubernetes/kube-state-metrics/blob/master/docs/namespace-metrics.md)
- [NetworkPolicy Metrics](https://github.com/kubernetes/kube-state-metrics/blob/master/docs/networkpolicy-metrics.md)
- [Node Metrics](https://github.com/kubernetes/kube-state-metrics/blob/master/docs/node-metrics.md)
- [PersistentVolume Metrics](https://github.com/kubernetes/kube-state-metrics/blob/master/docs/persistentvolume-metrics.md)
- [PersistentVolumeClaim Metrics](https://github.com/kubernetes/kube-state-metrics/blob/master/docs/persistentvolumeclaim-metrics.md)
- [Pod Disruption Budget Metrics](https://github.com/kubernetes/kube-state-metrics/blob/master/docs/poddisruptionbudget-metrics.md)
- [Pod Metrics](https://github.com/kubernetes/kube-state-metrics/blob/master/docs/pod-metrics.md)
- [ReplicaSet Metrics](https://github.com/kubernetes/kube-state-metrics/blob/master/docs/replicaset-metrics.md)
- [ReplicationController Metrics](https://github.com/kubernetes/kube-state-metrics/blob/master/docs/replicationcontroller-metrics.md)
- [ResourceQuota Metrics](https://github.com/kubernetes/kube-state-metrics/blob/master/docs/resourcequota-metrics.md)
- [Secret Metrics](https://github.com/kubernetes/kube-state-metrics/blob/master/docs/secret-metrics.md)
- [Service Metrics](https://github.com/kubernetes/kube-state-metrics/blob/master/docs/service-metrics.md)
- [StatefulSet Metrics](https://github.com/kubernetes/kube-state-metrics/blob/master/docs/statefulset-metrics.md)
- [StorageClass Metrics](https://github.com/kubernetes/kube-state-metrics/blob/master/docs/storageclass-metrics.md)
- [ValidatingWebhookConfiguration Metrics](https://github.com/kubernetes/kube-state-metrics/blob/master/docs/validatingwebhookconfiguration-metrics.md)
- [VerticalPodAutoscaler Metrics](https://github.com/kubernetes/kube-state-metrics/blob/master/docs/verticalpodautoscaler-metrics.md)
- [VolumeAttachment Metrics](https://github.com/kubernetes/kube-state-metrics/blob/master/docs/volumeattachment-metrics.md)

可以通过 prometheus 配置 scrape 的 target 为 kube-state-metrics ，将数据持久保存起来。

不过 metrics-server 和 kube-state-metrics 之间还是有很大不同的，二者的主要区别如下：

官方说明的区别：<https://github.com/kubernetes/kube-state-metrics#kube-state-metrics-vs-metrics-server>

- metrics-server 主要关注的是[资源度量 API](https://github.com/kubernetes/community/blob/master/contributors/design-proposals/instrumentation/resource-metrics-api.md) 的实现，比如 CPU、文件描述符、内存、请求延时等指标。
- kube-state-metrics 主要关注的是业务相关的一些元数据，比如 Deployment、Pod、副本状态等。

metric-server 的对比

- metric-server（或 heapster）是从 api-server 中获取 cpu、内存使用率这种监控指标，并把他们发送给存储后端，如 influxdb 或云厂商，他当前的核心作用是：为 HPA 等组件提供决策指标支持。
- kube-state-metrics 关注于获取 k8s 各种资源的最新状态，如 deployment 或者 daemonset，之所以没有把 kube-state-metrics 纳入到 metric-server 的能力中，是因为他们的关注点本质上是不一样的。metric-server 仅仅是获取、格式化现有数据，写入特定的存储，实质上是一个监控系统。而 kube-state-metrics 是将 k8s 的运行状况在内存中做了个快照，并且获取新的指标，但他没有能力导出这些指标
- 换个角度讲，kube-state-metrics 本身是 metric-server 的一种数据来源，虽然现在没有这么做。
- 另外，像 Prometheus 这种监控系统，并不会去用 metric-server 中的数据，他都是自己做指标收集、集成的（Prometheus 包含了 metric-server 的能力），但 Prometheus 可以监控 metric-server 本身组件的监控状态并适时报警，这里的监控就可以通过 kube-state-metrics 来实现，如 metric-serverpod 的运行状态。

## 与 Kubernetes 的兼容性

kube-state-metrics 使用 client-go 与 Kubernetes 集群进行交互。支持的 Kubernetes 集群版本由 client-go 决定。client-go 和 Kubernetes 集群的兼容性矩阵可以在这里找到。所有额外的兼容性只是尽最大努力，或者碰巧仍然/已经被支持。

如果 kube-state-metrics 与 kubernetes 版本不兼容，通常会出现如下问题：

- 获取某些资源的指标错误，因为 kubernetes 资源的 API 版本与 kube-state-metrics 所使用的 client-go 采集的资源的 API 版本不一致，报错信息类似 `failed to list *${APIVersion}.${Resource}: the server could not find the requested resource`：

```bash
pkg/mod/k8s.io/client-go@v0.24.1/tools/cache/reflector.go:167: failed to list *v1.PodDisruptionBudget: the server could not find the requested resource
pkg/mod/k8s.io/client-go@v0.24.1/tools/cache/reflector.go:167: Failed to watch *v1.CronJob: failed to list *v1.CronJob: the server could not find the requested resource
```

### 兼容矩阵

| **kube-state-metrics** | **Kubernetes 1.19** | **Kubernetes 1.23** | **Kubernetes 1.24** |
| ---------------------- | ------------------- | ------------------- | ------------------- |
| **v2.3.0**             | ✓                   | ✓                   | -                   |
| **v2.4.2**             | -/✓                 | ✓                   | -                   |
| **v2.5.0**             | -                   | ✓                   | ✓                   |

- `✓` 完全支持的版本
- `-` Kubernetes 集群具有 client-go 库无法使用的功能(其他 API 对象，不推荐使用的 API 等)。

# [Node Metrics](https://github.com/kubernetes/kube-state-metrics/blob/master/docs/node-metrics.md)(节点指标)

**kube_node_status_allocatable** # 可调度的节点中，各种资源的可分配额度
这里面的资源指的 memory、cpu、pods 等等，也就是说，这个指标反应了集群中每个节点上可以被使用内存有多少、CPU 有多少、可以部署多少个 Pod 等等。

# [Pod Metrics](https://github.com/kubernetes/kube-state-metrics/blob/master/docs/pod-metrics.md)(Pod 指标)

**kube_pod_container_resource_limits** # 为 Pod 中每个容器配置的 `.spec.container.resources.limits` 字段的值
**kube_pod_container_resource_requests** # 为 Pod 中每个容器配置的 `.spec.container.resources.requests` 字段的值
