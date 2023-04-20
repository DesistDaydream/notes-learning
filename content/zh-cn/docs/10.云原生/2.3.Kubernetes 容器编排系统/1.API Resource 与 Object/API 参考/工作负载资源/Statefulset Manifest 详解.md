---
title: Statefulset Manifest 详解
---

# 概述

> 参考：
>
> - [API 文档单页](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.21/#statefulset-v1-apps)
> - [官方文档，参考-KubernetesAPI-工作负载资源-StatefulSet](https://kubernetes.io/docs/reference/kubernetes-api/workload-resources/stateful-set-v1/)

# Manifest 中的顶层字段
- **apiVersion**: apps/v1
- **kind**: StatefulSet
- **metadata**([metadata](#metadata))
- **spec**([spec](#spec)) # 指明该 StatefulSet 的规格
- **status**([status](#status))

# metadata

Statefulset 对象的元数据，该字段内容详见通用定义的 [ObjectMeta](/docs/10.云原生/2.3.Kubernetes%20 容器编排系统/1.API、Resource(资源)、Object(对象)/API%20 参考/Common%20Definitions(通用定义).md Definitions(通用定义).md)

# spec

spec 用来描述一个 Statefulset 应该具有的属性。也就是用来定义 Statefulset 的行为规范。一共分为如下几类

- 描述 Statefulset 类型的控制器的行为
- 描述 Statefulset 控制器所关联的 Pod 的属性。

## 控制器行为

**podManagementPolicy: \<STRING>** # Pod 管理策略。`默认值：OrderedReady`
此配置只影响扩、缩 Pod 的行为，更新 Pod 不受此配置控制。可用的值有以下两个：

- OrderedReady # 按照 Pod 的次序依次创建每个 Pod 并等待 Ready 之后才创建后面的 Pod
- Parallel # 并行创建或删除 Pod（不等待前面的 Pod Ready 就开始创建所有的 Pod）

**replicas: \<INT>** # 该控制器运行的 Pod 数量，`默认值：1`。
**selector: \<Object> # 必须的**。Pod 的选择器，根据标签匹配要控制的 Pod。必须与 `template.metadata.labels` 的内容匹配。

- 该字段内容详见通用定义的[ LabelSelector](/docs/10.云原生/2.3.Kubernetes%20 容器编排系统/1.API、Resource(资源)、Object(对象)/API%20 参考/Common%20Definitions(通用定义)/LabelSelector%20 详解.md Definitions(通用定义)/LabelSelector 详解.md)。

**serviceName: \<STRING> # 必须的**。serviceName 是管理此 StatefulSet 的服务的名称。

该服务必须在 StatefulSet 之前存在，并且负责该集合的网络标识。 Pod 会遵循以下格式获取 DNS 或 hostname：pod-specific-string.serviceName.default.svc.cluster.local，其中"pod-specific-string"由 StatefulSet 控制器管理。

说白了，就是该字段指定的 service 名称将会自动生成子域名(而只有 headless 类型的 svc 才具有子域名)，假如现在有如下 pod 和 svc

```bash
~]# kubectl get pod -n cs-monitoring -o wide
NAME              READY   STATUS    RESTARTS   AGE     IP            NODE           NOMINATED NODE   READINESS GATES
cs-prometheus-0   2/2     Running   0          4m15s   10.244.4.54   node-2.bj-cs   <none>           <none>
cs-prometheus-1   2/2     Running   0          4m15s   10.244.5.55   node-3.bj-cs   <none>           <none>
~]# kubectl get svc -n cs-monitoring
NAME                     TYPE        CLUSTER-IP      EXTERNAL-IP   PORT(S)                         AGE
cs-prometheus            NodePort    10.98.204.140   <none>        9090:31001/TCP,9091:31002/TCP   5d
cs-prometheus-headless   ClusterIP   None            <none>        9090/TCP                        35s
```

我们可以发现，`cs-prometheus-0.cs-prometheus-headless.cs-monitoring.svc` 这个域名将会固定解析 `10.244.4.54`。而 `cs-prometheus-1.cs-prometheus-headless.cs-monitoring.svc` 这个域名将会固定解析到 `10.244.5.55`。而两个 svc 的域名，则是正常的解析结果，headless 是轮询解析到 pod 的 IP，而正常的 svc 直接解析到 10.98.204.140。

**volumeClaimTemplates: <\[]Object>** # PVC 模板。用于从指定的 storageclass 中申请 PVC，可直接当做 volume，在 volumeMonut 中使用

注意：通过 volumeClaimTemplates 定义的 PVC 在 statefulset 删除后不会自动删除，详见：[官方 issue](https://github.com/kubernetes/kubernetes/issues/55045)

- metadata:
  - name: STRING # 指定 volumeClaimTemplates 的名称，该名称用于在 volumeMount 时使用
- spec:
  - accessModes: # 指定该 volume 的访问模式
    - ReadWriteOnce # 样例为读写模式
  - resources: # 指定存储资源的申请量，样例为需求 30G
    - requests:
      - storage: 30Gi
  - storageClassName: STRING # 指定要从哪个 storageclass 中申请资源
  - volumeMode: Filesystem # 指定卷模式，样例为 Filesystem

## Pod 属性

**template: \<Ojbect> # 必须的**。定义 Pod 的模板,使用 Pod 类型的 metadata 和 spec 字段。

- **metadata**([PodMetadata](/docs/10.云原生/2.3.Kubernetes%20容器编排系统/1.API%20Resource%20与%20Object/API%20参考/工作负载资源/Pod%20Manifest%20详解.md#metadata)) # 与 pod 资源定义的内容基本一致
- **spec**([PodSpec](/docs/10.云原生/2.3.Kubernetes%20容器编排系统/1.API%20Resource%20与%20Object/API%20参考/工作负载资源/Pod%20Manifest%20详解.md#spec)) # 与 pod 资源定义的内容基本一致