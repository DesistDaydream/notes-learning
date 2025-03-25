---
title: Admission Controllers 准入控制器
linkTitle: Admission Controllers 准入控制器
weight: 4
---

# 概述

> 参考：
>
> - [官方文档，参考-API 访问控制-使用准入控制器](https://kubernetes.io/docs/reference/access-authn-authz/admission-controllers/)
> - [官方文档，参考-API 访问控制-动态准入控制](https://kubernetes.io/docs/reference/access-authn-authz/extensible-admission-controllers/)
> - [理清 Kubernetes 中的准入控制（Admission Controller）](https://mp.weixin.qq.com/s/nwKO2dmfvXf6dFw-y-vU7A)
> - [公众号 - 运维开发故事，开发一个禁止删除 namespace 的控制器](https://mp.weixin.qq.com/s/GdxSWFEyM1PYP30f3a-FCQ)

准入控制器是**一段代码**，它会在**请求通过 认证和授权 之后**、**对象被持久化之前**，拦截到达 API Server 的请求。

由于准入控制器是拦截 API Server 中最后的持久化逻辑，所以现阶段 准入控制器在 kube-apiserver 自身中实现，一共由于两大类准入控制器

- **静态准入控制器** # kube-apiserver 默认内置的准入控制器，可以从 [准入控制器列表](#Yd0ra) 查看。
  - 比如 istio 为每个 Pod 注入 Sidecar 的功能，就是通过 Mutating 准入控制器实现的。
- **动态准入控制器** # 以 Webhook 的形式运行，请求到达 kube-apiserver 后将会根据 `ValidatingWebhookConfiguration` 资源的定义，将请求转发给自己编写的控制器来处理后再返回给 kube-apiserver。
  - 比如我们编写了一个程序：如果请求是删除 namespace 资源的话，则进制删除。那么将这个程序部署到 k8s 时，再创建一个 ValidatingWebhookConfiguration 对面，以告诉 API Server 将请求转发给咱编写的程序。此时咱的程序处理请求后，会告诉 API Server 是否可以继续执行这个请求。

准入控制器通常用以执行 **Validating(验证)**、**Mutating(变更)** 操作

- 验证操作即验证该请求是否可以执行
- 变更操作就是类似于 Istion，将创建的 Pod 中加入其他字段或减少某些字段。

## 目前版本中，默认启用的准入控制器

CertificateApproval
CertificateSigning
CertificateSubjectRestriction
DefaultIngressClass
DefaultStorageClass
DefaultTolerationSeconds
LimitRanger
MutatingAdmissionWebhook
NamespaceLifecycle
PersistentVolumeClaimResize
PodSecurity
Priority
ResourceQuota
RuntimeClass
ServiceAccount
StorageObjectInUseProtection
TaintNodesByCondition
ValidatingAdmissionWebhook

# 准入控制器列表

## MutatingAdmissionWebhook

当执行变更操作时，通过 Webhook 调用动态准入控制器

## ValidatingAdmissionWebhook

当执行验证操作时，通过 Webhook 调用动态准入控制器
