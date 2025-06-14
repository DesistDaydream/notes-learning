---
title: RBAC 相关 Manifests 详解
---

# 概述

> 参考：
>
> - [官方文档，参考 - Kubernetes API - 认证资源](https://kubernetes.io/docs/reference/kubernetes-api/authorization-resources/)

# Role

- **apiVersion**: rbac.authorization.k8s.io/v1
- **kind**: Role
- **metadata** ([ObjectMeta](https://kubernetes.io/docs/reference/kubernetes-api/common-definitions/object-meta/#ObjectMeta))
  Standard object's metadata.
- **rules**(\[]PolicyRule)Rules holds all the PolicyRules for this Role_PolicyRule holds information that describes a policy rule, but does not contain information about who the rule applies to or which namespace the rule applies to.\_
  - **rules.apiGroups** (\[]string)
    APIGroups is the name of the APIGroup that contains the resources. If multiple API groups are specified, any action requested against one of the enumerated resources in any API group will be allowed.
  - **rules.resources** (\[]string)
    Resources is a list of resources this rule applies to. ResourceAll represents all resources.
  - **rules.verbs** (\[]string), required
    Verbs is a list of Verbs that apply to ALL the ResourceKinds and AttributeRestrictions contained in this rule. VerbAll represents all kinds.
  - **rules.resourceNames** (\[]string)
    ResourceNames is an optional white list of names that the rule applies to. An empty set means that everything is allowed.
  - **rules.nonResourceURLs** (\[]string)
    NonResourceURLs is a set of partial urls that a user should have access to. \*s are allowed, but only as the full, final step in the path Since non-resource URLs are not namespaced, this field is only applicable for ClusterRoles referenced from a ClusterRoleBinding. Rules can either apply to API resources (such as "pods" or "secrets") or non-resource URL paths (such as "/api"), but not both

# RoleBinding

- **apiVersion**: rbac.authorization.k8s.io/v1
- **kind**: RoleBinding
- **metadata** ([ObjectMeta](https://kubernetes.io/docs/reference/kubernetes-api/common-definitions/object-meta/#ObjectMeta))
  Standard object's metadata.
- **roleRef**(RoleRef), requiredRoleRef can reference a Role in the current namespace or a ClusterRole in the global namespace. If the RoleRef cannot be resolved, the Authorizer must return an error._RoleRef contains information that points to the role being used_
  - **roleRef.apiGroup** (string), required
    APIGroup is the group for the resource being referenced
  - **roleRef.kind** (string), required
    Kind is the type of resource being referenced
  - **roleRef.name** (string), required
    Name is the name of resource being referenced
- **subjects**(\[]Subject)Subjects holds references to the objects the role applies to._Subject contains a reference to the object or user identities a role binding applies to. This can either hold a direct API object reference, or a value for non-objects such as user and group names._
  - **subjects.kind** (string), required
    Kind of object being referenced. Values defined by this API group are "User", "Group", and "ServiceAccount". If the Authorizer does not recognized the kind value, the Authorizer should report an error.
  - **subjects.name** (string), required
    Name of the object being referenced.
  - **subjects.apiGroup** (string)
    APIGroup holds the API group of the referenced subject. Defaults to "" for ServiceAccount subjects. Defaults to "rbac.authorization.k8s.io" for User and Group subjects.
  - **subjects.namespace** (string)
    Namespace of the referenced object. If the object kind is non-namespace, such as "User" or "Group", and this value is not empty the Authorizer should report an error.
