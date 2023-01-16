---
title: Deployment Manifest 详解
---

# 概述

> 参考：
> - [API 文档，单页](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.21/#deployment-v1-apps)
> - [官方文档,参考-Kubernetes API-工作负载资源-Deployment](https://kubernetes.io/docs/reference/kubernetes-api/workload-resources/deployment-v1/)

## YAML 中的顶层节点

- apiVersion: apps/v1
- kind: Deployment
- [metadata: <Object>](#d12d7d74)
- [spec: <Object>](#5xtmr)
- [status: <Object>](#Sbe0x)

# metadata: <Object>

Deployment 对象的元数据，该字段内容详见通用定义的 [ObjectMeta](✏IT 学习笔记/☁️10.云原生/2.3.Kubernetes%20 容器编排系统/1.API、Resource(资源)、Object(对象)/API%20 参考/Common%20Definitions(通用定义)/ObjectMeta.md Definitions(通用定义)/ObjectMeta.md)

# spec: <Object>

spec 用来描述一个 Deployment 应该具有的属性。也就是用来定义 Deployment 的行为规范。一共分为如下几类

- 描述 Deployment 类型的控制器的行为
- 描述 Deployment 控制器所关联的 Pod 的属性。

## 控制器行为

**minReadySeconds: <INT>** # 新创建的 Pod 在启动后，经过 minReadySeconds 秒后一直没有崩溃，之后，将该 Pod 视为可用。`默认值：0`。
默认值 0 表示 Pod 准备就绪后即被视为可用。&#x20;
**progressDeadlineSeconds: <INT>** # 本 Deployment 对象被视为失败之前的等待时间，单位 秒。`默认值：600`
**replicas: <INT>** # 该控制器运行的 Pod 数量，`默认值：1`。
**revisionHistoryLimit: <INT> **# 可以保留的允许回滚的旧 ReplicaSet 对象的数量。`默认值：10`。控制器的历史可以通过 `kubectl rollout` 命令控制
**selector: <Object> # 必须的。**Pod 的标签选择器，根据标签匹配要控制的 Pod。必须与 `template.metadata.labels` 的内容匹配。

- 该字段内容详见通用定义的[ LabelSelector](✏IT 学习笔记/☁️10.云原生/2.3.Kubernetes%20 容器编排系统/1.API、Resource(资源)、Object(对象)/API%20 参考/Common%20Definitions(通用定义)/LabelSelector%20 详解.md Definitions(通用定义)/LabelSelector 详解.md)。

**strategy: <Ojbect>** # 定义用一个新的 pod 代替现有 pod 的部署策略(更新 pod 的策略)

- **rollingUpdate: <Object>** # 当更新策略为 rollingUpdate 时，需要配置滚动更新的参数
  - **maxSurge: <STRING> **# 设定在更新时最大可用的 Pod 数，就是先添加几个新的 Pod 再删除老的
  - **maxUnavailable: <STRING> **# 设定在更新时最大不可用的 Pod 数
- **type: <STRING>Recreate|RollingUpdate** # 指定更新策略的类型，Recreate(重新创建) 与 RollingUpdate(滚动更新)。`默认值：RollingUpdate`
  - Recreate 是删除一个创建一个

## Pod 属性

**template: <Ojbect> # 必须的**。定义 Pod 的模板,使用 Pod 类型的 metadata 和 spec 字段。

- **metadata:** #与 pod 资源定义的内容基本一致
  - ...
- **spec:** #与 pod 资源定义的内容基本一致
  - ...

# status: <Object>

# Manifests 样例

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: myapp
  labels:
    name: myapp
spec:
  replicas: 1
  selector:
    matchLabels:
      name: myapp
  template:
    metadata:
      name: myapp
      labels:
        name: myapp
    spec:
      containers:
        - name: myapp
          image: lchdzh/network-test
```
