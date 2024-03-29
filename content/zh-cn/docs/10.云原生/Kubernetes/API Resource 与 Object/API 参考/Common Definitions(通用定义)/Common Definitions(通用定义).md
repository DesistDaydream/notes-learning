---
title: "Common Definitions(通用定义)"
weight: 1
---

# 概述

> 参考：
>
> - [官方文档，参考-KubernetesAPI-通用定义](https://kubernetes.io/docs/reference/kubernetes-api/common-definitions/)

Common Definitions(通用定义)

## [DeleteOptions](https://kubernetes.io/docs/reference/kubernetes-api/common-definitions/delete-options/)

DeleteOptions may be provided when deleting an API object.

## [LabelSelector](/docs/10.云原生/Kubernetes/API%20Resource%20与%20Object/API%20参考/Common%20Definitions(通用定义)/LabelSelector.md)

A label selector is a label query over a set of resources.

## [ListMeta](https://kubernetes.io/docs/reference/kubernetes-api/common-definitions/list-meta/)

ListMeta describes metadata that synthetic resources must have, including lists and various status objects.

## [LocalObjectReference](https://kubernetes.io/docs/reference/kubernetes-api/common-definitions/local-object-reference/)

LocalObjectReference contains enough information to let you locate the referenced object inside the same namespace.

## [NodeSelectorRequirement](https://kubernetes.io/docs/reference/kubernetes-api/common-definitions/node-selector-requirement/)

A node selector requirement is a selector that contains values, a key, and an operator that relates the key and values.

## [ObjectFieldSelector](https://kubernetes.io/docs/reference/kubernetes-api/common-definitions/object-field-selector/)

ObjectFieldSelector selects an APIVersioned field of an object.

## [ObjectMeta](/docs/10.云原生/Kubernetes/API%20Resource%20与%20Object/API%20参考/Common%20Definitions(通用定义)/ObjectMeta.md)

ObjectMeta 是所有资源持久化成对象后必须要具有的元数据，其中包括对象的 名称、所在名称空间 等等。对应的 yaml 字段为 `.metadata`

## [ObjectReference](https://kubernetes.io/docs/reference/kubernetes-api/common-definitions/object-reference/)

ObjectReference contains enough information to let you inspect or modify the referred object.

## [Patch](https://kubernetes.io/docs/reference/kubernetes-api/common-definitions/patch/)

Patch is provided to give a concrete name and type to the Kubernetes PATCH request body.

## [Quantity](https://kubernetes.io/docs/reference/kubernetes-api/common-definitions/quantity/)

Quantity is a fixed-point representation of a number.

## [ResourceFieldSelector](https://kubernetes.io/docs/reference/kubernetes-api/common-definitions/resource-field-selector/)

ResourceFieldSelector represents container resources (cpu, memory) and their output format.

## [Status](/docs/10.云原生/Kubernetes/API%20Resource%20与%20Object/API%20参考/Common%20Definitions(通用定义)/Status.md)

Status 字段是不返回其他对象的调用的返回值。

## [TypedLocalObjectReference](https://kubernetes.io/docs/reference/kubernetes-api/common-definitions/typed-local-object-reference/)

TypedLocalObjectReference contains enough information to let you locate the typed referenced object inside the same namespace.
