---
title: API 参考
---

# 概述

> 参考：
> - [官方文档，参考-API 概述-API](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.23)(这里是通过单一页面显示 API 资源各字段详解)
>   - 链接里是 1.23 的，想查看其他版本 API，改变 URL 中的版本即可。
> - [官方文档，参考-KubernetesAPI](https://kubernetes.io/docs/reference/kubernetes-api/)(这里是通过多级页面显示 API 资源各字段详解)
>   - 这些连接的内容，其实是 `kubectl explain` 命令的内容显示在浏览器中了。
> - [OpeaAPI 格式文档](https://github.com/kubernetes/kubernetes/blob/master/api/openapi-spec/swagger.json)

在 [Kubernetes API](https://www.yuque.com/go/doc/33168662) 章节，已经可以看到单一页面的详解中对 API 的分类，在本笔记后面的部分对各资源 Manifest 详解中，其实已经描述了 API 中各个字段的含义。所以本篇文章不会详解每个 API，而是记录一下如何通过 Kubernetes 官网来查找 API 详解，以及如何使用官方文档查看 API 详解。

如果笔记中记录得不够详细，`kubectl explain` 命令也看着不方便，那么通过这篇文章中介绍的官方文档中的 API 详解来查看，将会更加直观。

**Kubernetes API 参考中将会描述每种资源的 Manifests 中每个字段(即.YAML 中的节点)的含义。**

下面是文档中占位符说明：
**\[]TYPE** # 表示该字段由数组组成，数组元素类型为 TYPE，比如 \[]STRING 格式应该就是下面这样

```yaml
args:
  - deletecr
  - --ns
  - --name
```

**map[STRING]STRING** # 表示多个键/值对。键 和 值 的数据类型都是 STRING。

```yaml
labels:
  key1: value1
  key2: value2
```

**OBJECT** # 表示复合结构的 map。

```yaml
resources:
  limits:
    cpu: "2"
    memory: 2Gi
  requests:
    cpu: 500m
    memory: 400Mi
```

**[]OBJECT** # 表示该字段由数组组成，并且数组中的元素都是一个 OBJECT，比如格式应该像下面这样

```yaml
containers:
  - args:
      - AAA
      - BBB
    name: XXX
    image: XXX
  - name: YYY
    image: YYY
```

**在每种资源的 Manifests 中，会有一些共用的部分称为**[**通用定义**](/docs/IT学习笔记/10.云原生/2.3.Kubernetes%20 容器编排系统/1.API、Resource(资源)、Object(对象)/API%20 参考/Common%20Definitions(通用定义).md 参考/Common Definitions(通用定义).md)**（也可以说是功能定义），比如常见的 **[**LabelSelector**](/docs/IT学习笔记/10.云原生/2.3.Kubernetes%20 容器编排系统/1.API、Resource(资源)、Object(对象)/API%20 参考/Common%20Definitions(通用定义)/LabelSelector%20 详解.md 参考/Common Definitions(通用定义)/LabelSelector 详解.md)**，这属于资源的 Manifests 的一部分。很多组件在解析 Manifests 中的通用定义时，都会遵循相同的规则。除了通用定义以外的，都属于 K8S 的资源定义，比如定义 **[**Pod**](/docs/IT学习笔记/10.云原生/2.3.Kubernetes%20 容器编排系统/1.API、Resource(资源)、Object(对象)/API%20 参考/工作负载资源/Pod%20Manifest%20 详解.md 参考/工作负载资源/Pod Manifest 详解.md)** 的 API 参考、定义 **[**Service**](/docs/IT学习笔记/10.云原生/2.3.Kubernetes%20 容器编排系统/1.API、Resource(资源)、Object(对象)/API%20 参考/服务资源/Service%20Manifests%20 详解.md 参考/服务资源/Service Manifests 详解.md)** 的 API 参考等等。**

这是单一页面的样子。左侧是根据对资源的分类而形成的目录，右侧是完整的页面
![](https://notes-learning.oss-cn-beijing.aliyuncs.com/dkxdpv/1616120193938-a171af16-575d-4de6-951a-99cdca271a50.png)
这是多级页面的样子，该 API 详解是内含在官方文档中的，并且对 API 进行了细致的分类
![](https://notes-learning.oss-cn-beijing.aliyuncs.com/dkxdpv/1616120193972-2c538ed5-7a6f-4aca-bf11-732240aa84d6.png)

## Kubernetes API 删除和弃用流程

> 参考：
> - [官方文档，参考-API 概述-Kubernetes 弃用策略](https://kubernetes.io/docs/reference/using-api/deprecation-policy/)

Kubernetes 项目有一个记录良好的特性弃用策略\[1]。该策略规定，只有当同一 API 的更新的、稳定的版本可用时，才可以弃用稳定的 API，并且 API 对于每个稳定性级别都有一个最短的生存期。给弃用的 API，是在未来的 Kubernetes 版本中被标记为删除的 API；它将继续运行，直到给删除（从弃用至少一年），但使用将导致显示警告。删除的 API 在当前版本中不再可用，此时你必须迁移到使用替换的 API。

- GA（Generally available，普遍可用）或稳定的 API 版本可能会被标记为弃用，但不得在 Kubernetes 的主要版本中删除。
- 测试版或预发布 API 版本弃用后，必须支持 3 个版本。
- Alpha 或实验 API 版本可能会在任何版本中被删除，恕不另行通知。

无论某个 API 是因为某个功能从测试版升级到稳定版而被删除，还是因为该 API 没有成功，所有的删除都遵循这个弃用策略。每当删除一个 API 时，迁移选项都会在文档中提供说明。

# API 分类

- [Workloads Resources](/docs/IT学习笔记/10.云原生/2.3.Kubernetes%20 容器编排系统/1.API、Resource(资源)、Object(对象)/API%20 参考/工作负载资源.md 参考/工作负载资源.md)(工作负载资源)
- [Services Resources](https://kubernetes.io/docs/reference/kubernetes-api/service-resources/)(服务资源)
- [Config and Storage Resources](https://kubernetes.io/docs/reference/kubernetes-api/config-and-storage-resources/)(配置与存储资源)
- [Authentication Resources](https://kubernetes.io/docs/reference/kubernetes-api/authentication-resources/)(认证资源)
- [Authorization Resources](https://kubernetes.io/docs/reference/kubernetes-api/authorization-resources/)(授权资源)
- [Policies Resources](https://kubernetes.io/docs/reference/kubernetes-api/policies-resources/)(策略资源)
- [Extend Resources](https://kubernetes.io/docs/reference/kubernetes-api/extend-resources/)(扩展资源)
- [Cluster Resources](https://kubernetes.io/docs/reference/kubernetes-api/cluster-resources/)(集群资源)
- [Common Definitions](https://kubernetes.io/docs/reference/kubernetes-api/common-definitions/)(通用定义) # 在多种资源 API 中，可以使用的 API。比如 节点选择器、meta 字段 等等
- [Common Parameters](https://kubernetes.io/docs/reference/kubernetes-api/common-parameters/common-parameters/)

## Config and Storage Resources

##### [ConfigMap](https://kubernetes.io/docs/reference/kubernetes-api/config-and-storage-resources/config-map-v1/)

ConfigMap holds configuration data for pods to consume.

##### [Secret](https://kubernetes.io/docs/reference/kubernetes-api/config-and-storage-resources/secret-v1/)

Secret holds secret data of a certain type.

##### [Volume](https://kubernetes.io/docs/reference/kubernetes-api/config-and-storage-resources/volume/)

Volume represents a named volume in a pod that may be accessed by any container in the pod.

##### [PersistentVolumeClaim](https://kubernetes.io/docs/reference/kubernetes-api/config-and-storage-resources/persistent-volume-claim-v1/)

PersistentVolumeClaim is a user's request for and claim to a persistent volume.

##### [PersistentVolume](https://kubernetes.io/docs/reference/kubernetes-api/config-and-storage-resources/persistent-volume-v1/)

PersistentVolume (PV) is a storage resource provisioned by an administrator.

##### [StorageClass](https://kubernetes.io/docs/reference/kubernetes-api/config-and-storage-resources/storage-class-v1/)

StorageClass describes the parameters for a class of storage for which PersistentVolumes can be dynamically provisioned.

##### [VolumeAttachment](https://kubernetes.io/docs/reference/kubernetes-api/config-and-storage-resources/volume-attachment-v1/)

VolumeAttachment captures the intent to attach or detach the specified volume to/from the specified node.

##### [CSIDriver](https://kubernetes.io/docs/reference/kubernetes-api/config-and-storage-resources/csi-driver-v1/)

CSIDriver captures information about a Container Storage Interface (CSI) volume driver deployed on the cluster.

##### [CSINode](https://kubernetes.io/docs/reference/kubernetes-api/config-and-storage-resources/csi-node-v1/)

CSINode holds information about all CSI drivers installed on a node.

##### [CSIStorageCapacity v1beta1](https://kubernetes.io/docs/reference/kubernetes-api/config-and-storage-resources/csi-storage-capacity-v1beta1/)

CSIStorageCapacity stores the result of one CSI GetCapacity call.

## Authentication Resources

##### [ServiceAccount](https://kubernetes.io/docs/reference/kubernetes-api/authentication-resources/service-account-v1/)

ServiceAccount binds together: _ a name, understood by users, and perhaps by peripheral systems, for an identity _ a principal that can be authenticated and authorized \* a set of secrets.

##### [TokenRequest](https://kubernetes.io/docs/reference/kubernetes-api/authentication-resources/token-request-v1/)

TokenRequest requests a token for a given service account.

##### [TokenReview](https://kubernetes.io/docs/reference/kubernetes-api/authentication-resources/token-review-v1/)

TokenReview attempts to authenticate a token to a known user.

##### [CertificateSigningRequest](https://kubernetes.io/docs/reference/kubernetes-api/authentication-resources/certificate-signing-request-v1/)

CertificateSigningRequest objects provide a mechanism to obtain x509 certificates by submitting a certificate signing request, and having it asynchronously approved and issued.

## Authorization Resources

##### [LocalSubjectAccessReview](https://kubernetes.io/docs/reference/kubernetes-api/authorization-resources/local-subject-access-review-v1/)

LocalSubjectAccessReview checks whether or not a user or group can perform an action in a given namespace.

##### [SelfSubjectAccessReview](https://kubernetes.io/docs/reference/kubernetes-api/authorization-resources/self-subject-access-review-v1/)

SelfSubjectAccessReview checks whether or the current user can perform an action.

##### [SelfSubjectRulesReview](https://kubernetes.io/docs/reference/kubernetes-api/authorization-resources/self-subject-rules-review-v1/)

SelfSubjectRulesReview enumerates the set of actions the current user can perform within a namespace.

##### [SubjectAccessReview](https://kubernetes.io/docs/reference/kubernetes-api/authorization-resources/subject-access-review-v1/)

SubjectAccessReview checks whether or not a user or group can perform an action.

##### [ClusterRole](https://kubernetes.io/docs/reference/kubernetes-api/authorization-resources/cluster-role-v1/)

ClusterRole is a cluster level, logical grouping of PolicyRules that can be referenced as a unit by a RoleBinding or ClusterRoleBinding.

##### [ClusterRoleBinding](https://kubernetes.io/docs/reference/kubernetes-api/authorization-resources/cluster-role-binding-v1/)

ClusterRoleBinding references a ClusterRole, but not contain it.

##### [Role](https://kubernetes.io/docs/reference/kubernetes-api/authorization-resources/role-v1/)

Role is a namespaced, logical grouping of PolicyRules that can be referenced as a unit by a RoleBinding.

##### [RoleBinding](https://kubernetes.io/docs/reference/kubernetes-api/authorization-resources/role-binding-v1/)

RoleBinding references a role, but does not contain it.

## Policies Resources

##### [LimitRange](https://kubernetes.io/docs/reference/kubernetes-api/policy-resources/limit-range-v1/)

LimitRange sets resource usage limits for each kind of resource in a Namespace.

##### [ResourceQuota](https://kubernetes.io/docs/reference/kubernetes-api/policy-resources/resource-quota-v1/)

ResourceQuota sets aggregate quota restrictions enforced per namespace.

##### [NetworkPolicy](https://kubernetes.io/docs/reference/kubernetes-api/policy-resources/network-policy-v1/)

NetworkPolicy describes what network traffic is allowed for a set of Pods.

##### [PodDisruptionBudget](https://kubernetes.io/docs/reference/kubernetes-api/policy-resources/pod-disruption-budget-v1/)

PodDisruptionBudget is an object to define the max disruption that can be caused to a collection of pods.

##### [PodSecurityPolicy v1beta1](https://kubernetes.io/docs/reference/kubernetes-api/policy-resources/pod-security-policy-v1beta1/)

PodSecurityPolicy governs the ability to make requests that affect the Security Context that will be applied to a pod and container.

## Extend Resources

##### [CustomResourceDefinition](https://kubernetes.io/docs/reference/kubernetes-api/extend-resources/custom-resource-definition-v1/)

CustomResourceDefinition represents a resource that should be exposed on the API server.

##### [MutatingWebhookConfiguration](https://kubernetes.io/docs/reference/kubernetes-api/extend-resources/mutating-webhook-configuration-v1/)

MutatingWebhookConfiguration describes the configuration of and admission webhook that accept or reject and may change the object.

##### [ValidatingWebhookConfiguration](https://kubernetes.io/docs/reference/kubernetes-api/extend-resources/validating-webhook-configuration-v1/)

ValidatingWebhookConfiguration describes the configuration of and admission webhook that accept or reject and object without changing it.

## Cluster Resources

##### [Node](https://kubernetes.io/docs/reference/kubernetes-api/cluster-resources/node-v1/)

Node is a worker node in Kubernetes.

##### [Namespace](https://kubernetes.io/docs/reference/kubernetes-api/cluster-resources/namespace-v1/)

Namespace provides a scope for Names.

##### [Event](https://kubernetes.io/docs/reference/kubernetes-api/cluster-resources/event-v1/)

Event is a report of an event somewhere in the cluster.

##### [APIService](https://kubernetes.io/docs/reference/kubernetes-api/cluster-resources/api-service-v1/)

APIService represents a server for a particular GroupVersion.

##### [Lease](https://kubernetes.io/docs/reference/kubernetes-api/cluster-resources/lease-v1/)

Lease defines a lease concept.

##### [RuntimeClass](https://kubernetes.io/docs/reference/kubernetes-api/cluster-resources/runtime-class-v1/)

RuntimeClass defines a class of container runtime supported in the cluster.

##### [FlowSchema v1beta1](https://kubernetes.io/docs/reference/kubernetes-api/cluster-resources/flow-schema-v1beta1/)

FlowSchema defines the schema of a group of flows.

##### [PriorityLevelConfiguration v1beta1](https://kubernetes.io/docs/reference/kubernetes-api/cluster-resources/priority-level-configuration-v1beta1/)

PriorityLevelConfiguration represents the configuration of a priority level.

##### [Binding](https://kubernetes.io/docs/reference/kubernetes-api/cluster-resources/binding-v1/)

Binding ties one object to another; for example, a pod is bound to a node by a scheduler.

##### [ComponentStatus](https://kubernetes.io/docs/reference/kubernetes-api/cluster-resources/component-status-v1/)

ComponentStatus (and ComponentStatusList) holds the cluster validation info.
