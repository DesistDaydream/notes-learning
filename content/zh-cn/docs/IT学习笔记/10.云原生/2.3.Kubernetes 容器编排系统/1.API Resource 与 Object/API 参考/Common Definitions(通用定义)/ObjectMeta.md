---
title: ObjectMeta
---

# 概述

> 参考：
> - [官方文档，参考-KubernetesAPI-通用定义-ObjectMeta](https://kubernetes.io/docs/reference/kubernetes-api/common-definitions/object-meta/)

ObjectMeta 是所有持久化资源必须具有的元数据信息，也就是每个对象所具有的元数据信息。

# 基本字段

## name: \<STRING>

Name must be unique within a namespace. Is required when creating resources, although some resources may allow a client to request the generation of an appropriate name automatically. Name is primarily intended for creation idempotence and configuration definition. Cannot be updated. More info: <http://kubernetes.io/docs/user-guide/identifiers#names>

## generateName (string)

GenerateName is an optional prefix, used by the server, to generate a unique name ONLY IF the Name field has not been provided. If this field is used, the name returned to the client will be different than the name passed. This value will also be combined with a unique suffix. The provided value has the same validation rules as the Name field, and may be truncated by the length of the suffix required to make the value unique on the server.If this field is specified and the generated name exists, the server will NOT return a 409 - instead, it will either return 201 Created or 500 with Reason ServerTimeout indicating a unique name could not be found in the time allotted, and the client should retry (optionally after the time indicated in the Retry-After header).Applied only if Name is not specified. More info: <https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#idempotency>

## namespace: \<STRING>

Namespace defines the space within which each name must be unique. An empty namespace is equivalent to the "default" namespace, but "default" is the canonical representation. Not all objects are required to be scoped to a namespace - the value of this field for those objects will be empty.Must be a DNS_LABEL. Cannot be updated. More info: <http://kubernetes.io/docs/user-guide/namespaces>

## labels: \<map\[STRING]STRING)>

Map of string keys and values that can be used to organize and categorize (scope and select) objects. May match selectors of replication controllers and services. More info: <http://kubernetes.io/docs/user-guide/labels>

## annotations: \<map\[STRING]STRING>

- Annotations is an unstructured key value map stored with a resource that may be set by external tools to store and retrieve arbitrary metadata. They are not queryable and should be preserved when modifying objects. More info: <http://kubernetes.io/docs/user-guide/annotations>

# 系统字段

## finalizers: \[]string

Must be empty before the object is deleted from the registry. Each entry is an identifier for the responsible component that will remove the entry from the list. If the deletionTimestamp of the object is non-nil, entries in this list can only be removed. Finalizers may be processed and removed in any order. Order is NOT enforced because it introduces significant risk of stuck finalizers. finalizers is a shared field, any actor with permission can reorder it. If the finalizer list is processed in order, then this can lead to a situation in which the component responsible for the first finalizer in the list is waiting for a signal (field value, external system, or other) produced by a component responsible for a finalizer later in the list, resulting in a deadlock. Without enforced ordering finalizers are free to order amongst themselves and are not vulnerable to ordering changes in the list.

## managedFields: <\[]Object>

managedFields 主要用以声明该对象的各个**字段的管理者**是谁。简单点说，可以理解为，对象中每个字段是通过哪些程序更新的，都会在这里记录。更新字段的 程序、代码 等，即称为**字段的管理者**。这是
详见：[服务端应用-字段管理章节](https://kubernetes.io/docs/reference/using-api/server-side-apply/#field-management)(自己的笔记在：[《Kubernetes 对象管理》章节-服务端应用-字段管理](/docs/IT学习笔记/10.云原生/2.3.Kubernetes%20 容器编排系统/1.API、Resource(资源)、Object(对象)/Kubernetes%20 对象管理.md 对象管理.md))
ManagedFields maps workflow-id and version to the set of fields that are managed by that workflow. This is mostly for internal housekeeping, and users typically shouldn't need to set or understand this field. A workflow can be the user's name, a controller's name, or the name of a specific apply path like "ci-cd". The set of fields is always in the version that the workflow used when modifying the object._ManagedFieldsEntry is a workflow-id, a FieldSet and the group version of the resource that the fieldset applies to._
**apiVersion: \<STRING>** # APIVersion defines the version of this resource that this field set applies to. The format is "group/version" just like the top-level APIVersion field. It is necessary to track the version of a field set because it cannot be automatically converted.
**fieldsType** (string)FieldsType is the discriminator for the different fields format and version. There is currently only one possible value: "FieldsV1"
**fieldsV1** (FieldsV1)FieldsV1 holds the first JSON version format as described in the "FieldsV1" type._FieldsV1 stores a set of fields in a data structure like a Trie, in JSON format.Each key is either a '.' representing the field itself, and will always map to an empty set, or a string representing a sub-field or item. The string will follow one of these four formats: 'f:', where is the name of a field in a struct, or key in a map 'v:', where is the exact json formatted value of a list item 'i:', where is position of a item in a list 'k:', where is a map of a list item's key fields to their unique values If a key maps to an empty Fields value, the field that key represents is part of the set.The exact format is defined in sigs.k8s.io/structured-merge-diff_
**manager: \<STRING> **# Manager 是管理这些字段的工作流的标识符。比如 kubectl、kubectl-replace、kubctl-run、kubelet、Go-http-client 等等

- 这个主要是用来清晰表明该对象最近几次创建、更新，是由哪些程序更新的

**operation: \<STRING>** # 导致创建 managedField 字段下的这个条目的操作类型。此字段的唯一有效值是“Apply”和“Update”。

- 这个主要是用来清晰表明该对象最近几次创建、更新，是 manager 中标明的程序执行的哪些操作

**subresource** (string)Subresource is the name of the subresource used to update that object, or empty string if the object was updated through the main resource. The value of this field is used to distinguish between managers, even if they share the same name. For example, a status update will be distinct from a regular update using the same manager name. Note that the APIVersion field is not related to the Subresource field and it always corresponds to the version of the main resource.
**time** (Time)Time is timestamp of when these fields were set. It should always be empty if Operation is 'Apply'_Time is a wrapper around time.Time which supports correct marshaling to YAML and JSON. Wrappers are provided for many of the factory methods that the time package offers._

## ownerReferences (\[]OwnerReference)

\_Patch strategy: merge on key uid*List of objects depended by this object. If ALL objects in the list have been deleted, this object will be garbage collected. If this object is managed by a controller, then an entry in this list will point to this controller, with the controller field set to true. There cannot be more than one managing controller.\_OwnerReference contains enough information to let you identify an owning object. An owning object must be in the same namespace as the dependent, or be cluster-scoped, so there is no namespace field.*

- **ownerReferences.apiVersion** (string), requiredAPI version of the referent.
- **ownerReferences.kind** (string), requiredKind of the referent. More info: <https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds>
- **ownerReferences.name** (string), requiredName of the referent. More info: <http://kubernetes.io/docs/user-guide/identifiers#names>
- **ownerReferences.uid** (string), requiredUID of the referent. More info: <http://kubernetes.io/docs/user-guide/identifiers#uids>
- **ownerReferences.blockOwnerDeletion** (boolean)If true, AND if the owner has the "foregroundDeletion" finalizer, then the owner cannot be deleted from the key-value store until this reference is removed. Defaults to false. To set this field, a user needs "delete" permission of the owner, otherwise 422 (Unprocessable Entity) will be returned.
- **ownerReferences.controller** (boolean)If true, this reference points to the managing controller.

# 只读字段

## creationTimestamp (Time)

CreationTimestamp is a timestamp representing the server time when this object was created. It is not guaranteed to be set in happens-before order across separate operations. Clients may not set this value. It is represented in RFC3339 form and is in UTC.Populated by the system. Read-only. Null for lists. More info: <https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata>_Time is a wrapper around time.Time which supports correct marshaling to YAML and JSON. Wrappers are provided for many of the factory methods that the time package offers._

## deletionGracePeriodSeconds (int64)

Number of seconds allowed for this object to gracefully terminate before it will be removed from the system. Only set when deletionTimestamp is also set. May only be shortened. Read-only.

## deletionTimestamp (Time)

DeletionTimestamp is RFC 3339 date and time at which this resource will be deleted. This field is set by the server when a graceful deletion is requested by the user, and is not directly settable by a client. The resource is expected to be deleted (no longer visible from resource lists, and not reachable by name) after the time in this field, once the finalizers list is empty. As long as the finalizers list contains items, deletion is blocked. Once the deletionTimestamp is set, this value may not be unset or be set further into the future, although it may be shortened or the resource may be deleted prior to this time. For example, a user may request that a pod is deleted in 30 seconds. The Kubelet will react by sending a graceful termination signal to the containers in the pod. After that 30 seconds, the Kubelet will send a hard termination signal (SIGKILL) to the container and after cleanup, remove the pod from the API. In the presence of network partitions, this object may still exist after this timestamp, until an administrator or automated process can determine the resource is fully terminated. If not set, graceful deletion of the object has not been requested.Populated by the system when a graceful deletion is requested. Read-only. More info: <https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata>_Time is a wrapper around time.Time which supports correct marshaling to YAML and JSON. Wrappers are provided for many of the factory methods that the time package offers._

## generation (int64)

A sequence number representing a specific generation of the desired state. Populated by the system. Read-only.

## resourceVersion (string)

An opaque value that represents the internal version of this object that can be used by clients to determine when objects have changed. May be used for optimistic concurrency, change detection, and the watch operation on a resource or set of resources. Clients must treat these values as opaque and passed unmodified back to the server. They may only be valid for a particular resource or set of resources.Populated by the system. Read-only. Value must be treated as opaque by clients and . More info: <https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#concurrency-control-and-consistency>

## selfLink (string)

SelfLink is a URL representing this object. Populated by the system. Read-only.DEPRECATED Kubernetes will stop propagating this field in 1.20 release and the field is planned to be removed in 1.21 release.

## uid (string)

UID is the unique in time and space value for this object. It is typically generated by the server on successful creation of a resource and is not allowed to change on PUT operations.Populated by the system. Read-only. More info: <http://kubernetes.io/docs/user-guide/identifiers#uids>

# Ignored

- **clusterName** (string)The name of the cluster which the object belongs to. This is used to distinguish resources with same name and namespace in different clusters. This field is not set anywhere right now and apiserver is going to ignore it if set in create or update request.
