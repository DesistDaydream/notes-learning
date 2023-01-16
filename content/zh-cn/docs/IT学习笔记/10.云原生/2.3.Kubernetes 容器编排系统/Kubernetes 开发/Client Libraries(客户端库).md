---
title: Client Libraries(客户端库)
---

# 概述

> 参考：
> - [官方文档](https://kubernetes.io/docs/reference/using-api/client-libraries/)

**Client Libraries(客户端库)** 是各种编程语言的**第三方库的统称**。这些库可以用来让各种编程语言通过代码的方式访问 Kubernetes API。在使用这些库编写代码时，并不需要自己实现对 Kubernetes API 的调用和 处理 Request/Response，这些处理逻辑都在 Client Libraries 中包括了。客户端库还会处理诸如身份验证之类的行为。

如果代码在 Kubernetes 集群中运行，代码中的 Client Libraires 可以发现并使用 Kubernetes 的 ServiceAccount 进行身份验证。

如果代码在 Kubernetes 集群外运行，代码中的 Client Libraires 能够理解 [kubeconfig 文件](https://kubernetes.io/zh/docs/tasks/access-application-cluster/configure-access-multiple-clusters/) 格式来读取凭据和 API 服务器地址。

Kubernetes 现阶段官方支持 Go、Python、Java、 dotnet、Javascript 和 Haskell 语言的客户端库。还有一些其他客户端库由对应作者而非 Kubernetes 团队提供并维护。 参考客户端库了解如何使用其他语言 来访问 API 以及如何执行身份认证。
