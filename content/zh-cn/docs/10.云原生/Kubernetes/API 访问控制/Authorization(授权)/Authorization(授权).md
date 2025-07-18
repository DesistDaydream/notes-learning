---
title: Authorization(授权)
linkTitle: Authorization(授权)
weight: 1
---

# 概述

> 参考：
>
> - [官方文档，参考 - API 访问控制 - 授权](https://kubernetes.io/docs/reference/access-authn-authz/authorization/)

在 Kubernetes 中，在 **Authorization(i.e.授予访问权限，简称：授权)** 之前必须进行过 [Authentication(认证)](/docs/10.云原生/Kubernetes/API%20访问控制/Authentication(认证)/Authentication(认证).md)

## 授权流程

### 确定是允许还是拒绝请求

Kubernetes 使用 API 服务器授权 API 请求。它根据所有策略评估所有请求属性来决定允许或拒绝请求。 一个 API 请求的所有部分必须被某些策略允许才能继续。这意味着默认情况下拒绝权限。

（尽管 Kubernetes 使用 API 服务器，但是依赖于特定种类对象的特定字段的访问控制和策略由准入控制器处理。）

配置多个授权模块时，将按顺序检查每个模块。 如果任何授权模块批准或拒绝请求，则立即返回该决定，并且不会与其他授权模块协商。 如果所有模块对请求没有意见，则拒绝该请求。一个拒绝响应返回 HTTP 状态代码 403 。

### 审查您的请求属性

Kubernetes 仅审查以下 API 请求属性：

- user - 身份验证期间提供的 user 字符串。
- group - 经过身份验证的用户所属的组名列表。
- extra - 由身份验证层提供的任意字符串键到字符串值的映射。
- API - 指示请求是否针对 API 资源。
- Request path - 各种非资源端点的路径，如 /api 或 /healthz。
- API request verb - API 动词 get，list，create，update，patch，watch，proxy，redirect，delete 和 deletecollection 用于资源请求。要确定资源 API 端点的请求动词，请参阅确定请求动词。
- HTTP request verb - HTTP 动词 get，post，put 和 delete 用于非资源请求。
- Resource - 正在访问的资源的 ID 或名称（仅限资源请求） - 对于使用 get，update，patch 和 delete 动词的资源请求，您必须提供资源名称。
- Subresource - 正在访问的子资源（仅限资源请求）。
- Namespace - 正在访问的对象的名称空间（仅适用于命名空间资源请求）。
- API group - 正在访问的 API 组（仅限资源请求）。空字符串表示核心 API 组。

### 确定请求动词

要确定资源 API 端点的请求动词，需要检查所使用的 HTTP 动词以及请求是否对单个资源或资源集合起作用：

| HTTP 动词 | request 动词                                   |
| --------- | ---------------------------------------------- |
| POST      | create                                         |
| GET, HEAD | get (单个资源)，list (资源集合)                |
| PUT       | update                                         |
| PATCH     | patch                                          |
| DELETE    | delete (单个资源)，deletecollection (资源集合) |

Kubernetes 有时使用专门的动词检查授权以获得额外的权限。例如：

- [Pod 安全策略](https://kubernetes.io/docs/concepts/policy/pod-security-policy/) 检查 policy API 组中 podsecuritypolicies 资源的 use 动词的授权。
- [RBAC](https://kubernetes.io/docs/reference/access-authn-authz/rbac/#privilege-escalation-prevention-and-bootstrapping)检查 rbac.authorization.k8s.io API 组中 roles 和 clusterroles 资源的 bind 动词的授权。
- [认证](https://kubernetes.io/docs/reference/access-authn-authz/authentication/) layer 检查核心 API 组中 users，groups 和 serviceaccounts 的 impersonate 动词的授权，以及 authentication.k8s.io API 组中的 userextras

# 授权的实现方式

在 Kubernetes 中，可以通过多种方式来实现 Authorization(授权) 功能

## RBAC 授权

> 参考：
>
> - 官方文档：<https://kubernetes.io/docs/reference/access-authn-authz/rbac/>
> - RBAC 概念详见：[RBAC](/docs/7.信息安全/Access%20Control/RBAC.md)

**RBAC** # 基于角色的访问控制（RBAC）是一种基于企业内个人用户的角色来管理对计算机或网络资源的访问的方法。在这种语境中，权限是单个用户执行特定任务的能力，例如查看，创建或修改文件。要了解有关使用 RBAC 模式的更多信息，请参阅 [RBAC 模式](https://kubernetes.io/docs/reference/access-authn-authz/rbac/)。

- 当指定的 RBAC（基于角色的访问控制）使用 rbac.authorization.k8s.io API 组来驱动授权决策时，允许管理员通过 Kubernetes API 动态配置权限策略。
- 要启用 RBAC，请使用 --authorization-mode = RBAC 启动 apiserver 。

详见 ：[RBAC 授权](/docs/10.云原生/Kubernetes/API%20访问控制/Authorization(授权)/RBAC%20授权.md)

## ABAC 授权

> 参考：
>
> - 官方文档：<https://kubernetes.io/docs/reference/access-authn-authz/abac/>

**ABAC** # 基于属性的访问控制（ABAC）定义了一种访问控制范例，通过使用将属性组合在一起的策略，将访问权限授予用户。策略可以使用任何类型的属性（用户属性，资源属性，对象，环境属性等）。要了解有关使用 ABAC 模式的更多信息，请参阅 [ABAC 模式](https://kubernetes.io/docs/reference/access-authn-authz/abac/)。

## Node 授权

> 参考：
>
> - 官方文档：<https://kubernetes.io/docs/reference/access-authn-authz/node/>

**Node** # 一个专用授权程序，根据计划运行的 pod 为 kubelet 授予权限。了解有关使用节点授权模式的更多信息，请参阅[节点授权.](https://kubernetes.io/docs/reference/access-authn-authz/node/)

## Webhook 授权

> 参考：
>
> - 官方文档：<https://kubernetes.io/docs/reference/access-authn-authz/webhook/>

**Webhook**# WebHook 是一个 HTTP 回调：发生某些事情时调用的 HTTP POST；通过 HTTP POST 进行简单的事件通知。实现 WebHook 的 Web 应用程序会在发生某些事情时将消息发布到 URL。要了解有关使用 Webhook 模式的更多信息，请参阅 Webhook 模式。
