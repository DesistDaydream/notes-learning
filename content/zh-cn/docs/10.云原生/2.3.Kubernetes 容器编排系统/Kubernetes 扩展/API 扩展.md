---
title: API 扩展
---

# 概述

> 参考：
>
> - [官方文档，概念-扩展 Kubernetes-扩展 Kubernetes API](https://kubernetes.io/docs/concepts/extend-kubernetes/api-extension/)

Resource(资源) 是 Kubernetes API 中的一个 endpoint(端点)， 其中存储的是某个类别的 API 对象 的一个集合。 例如内置的 pods 资源包含一组 Pod 对象。(这里面“点”的意思是这么一种语境。我说的几点记住了吗？知识点。等等)

而扩展 Kubernetes API 实际上就是添加 **Custom Resource(自定义的资源)**。

## Custom Resource 自定义资源

**什么是自定义资源呢？**
如 [API 与 Resource](/docs/10.云原生/2.3.Kubernetes%20容器编排系统/1.API%20Resource%20与%20Object/API%20与%20Resource.md) 中介绍的，Kubernetes 自身的一切都抽象为 Resource(资源)。顾名思义，Custom Resource(自定义资源) 就是非 Kubernetes 核心的资源。如果要类比的话，那么 Custom Resource 与 Kubernetes 的关系，类似于 Linux 中，Module(模块) 与 Kernel(内核) 的关系。其实，再准确的说法应该是下文将要提到的 Operator，Operator 与 Kubernetes 的关系，类似于 Linux 中，Module(模块) 与 Kernel(内核) 的关系。现在很多 Kubernetes 核心功能现在都用自定义资源来实现，这使得 Kubernetes 更加模块化。

自定义资源可以像普通资源(比如.pod)一样被创建和销毁。一旦某个自定义资源被安装，就可以使用 kubectl 来创建和访问其中的对象，就像为 pods 这种内置资源所做的一样。

添加 Custom Resource 方式有以下两种

1. CustomResourceDefinitions(自定义资源)
   1. 使用 CustomResourceDefinition 对象来创建一个或多个 CRD 资源。
   2. 相对简单，创建 CRD 可以不必编程。详见：Custom Resource Definitions(CRD)
2. API Aggregation(API 聚合)
   1. 使用 APIService 对象来创建一个或多个 API 聚合 资源。
   2. 需要编程，但支持对 API 行为进行更多的控制，例如数据如何存储以及在不同 API 版本间如何转换等。详见：API Aggregation(聚合) Layer

Kubernetes 提供的上述两种选项以满足不同用户的需求，这样就既不会牺牲易用性也不会牺牲灵活性。

- 聚合 API 指的是一些下位的 API 服务器，运行在主 API 服务器后面；主 API 服务器以代理的方式工作。这种组织形式称作 API 聚合（API Aggregation，AA） 。 对用户而言，看起来仅仅是 Kubernetes API 被扩展了。
- CRD 允许用户创建新的资源类别同时又不必添加新的 API 服务器。 使用 CRD 时，你并不需要理解 API 聚合。

无论以哪种方式安装自定义资源，新的资源都会被当做自定义资源，以便与内置的 Kubernetes 资源(如 Pods) 和 内置的 API(如 v1.apps) 相区分。

## Custom Controller 自定义控制器

**为什么需要 Custom Controller 呢？**
需要 Custom Controller 的原因，就要从 Kubernetes 的 声明式 API 的特点说起，基于这种特点，所有资源都基于 Kubernetes 的控制器模型来维护其自身的状态。那么当我们创建一个 Custom Resource 时，谁又来维护它的状态的呢？这就是 Custom Controller 的由来。Custom Resource 不能由集群的 controller-manager 来维护，必须要自定义一个 Controller 才可以。

就 Custom Resource 本身而言，它只能用来存取结构化的数据。 当你将 Custom Resource 与 Custom Controller 相结合时，就能够提供真正的声明式 API（Declarative API）。

使用[声明式 API](https://kubernetes.io/zh/docs/concepts/overview/kubernetes-api/)， 你可以声明或者设定你的资源的期望状态，并尝试让 Kubernetes 对象的当前状态 同步到其期望状态。控制器负责将结构化的数据解释为用户所期望状态的记录，并持续地维护该状态。

Custom Controller 与 Custom Resource 不同，并不能通过原生对象来直接创建。Custom Controller 是一个逻辑上的概念。

一般通过代码编写一个程序，以 Pod 方式运行在集群中，专门用来监听 Custom Resource 的状态，这种运作模型的程序，就是 Custom Controller。在一个运行中的集群上部署和更新自定义控制器，这类操作与集群的生命周期无关。 自定义控制器可以用于任何类别的资源，不过它们与自定义资源结合起来时最为有效。

## Operator 模式

将 Custom Resource 与 Custom Controller 结合起来使用就是 [Operator 模式](https://coreos.com/blog/introducing-operators.html)。你可以使用 Custom Controller 来将特定于某应用的领域知识组织起来，以编码的形式构造对 Kubernetes API 的扩展。

详见：Operator 模式介绍

所以，Operator 也可以称为 Custom Controller。

# 应用示例

kube-prometheus 项目就是一个典型的 API 扩展，其中包含 prometheus 这个 operater，还有多种 CRD ，还有名为 v1beta1.metrics.k8s.io 的聚合 API。这一整套扩展组合出来一个自动化的 prometheus 产品。prometheus-operator 会管理各种 CRD 创建出来的资源，以及 v1beta1.metrics.k8s.io 这个聚合 API 所关联的 service 的后端 pod 的状态。
