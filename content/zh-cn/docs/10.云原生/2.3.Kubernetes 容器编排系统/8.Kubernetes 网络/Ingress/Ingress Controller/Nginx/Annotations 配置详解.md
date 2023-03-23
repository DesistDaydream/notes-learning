---
title: Annotations 配置详解
---

# 概述

> 参考：
>
> - [官方文档，用户指南-Annotations](https://kubernetes.github.io/ingress-nginx/user-guide/nginx-configuration/annotations/)

与 ConfigMap 实现配置 Nginx Ingress Controller 运行时行为类似，只不过，Annotations 的方式，是通过在 Ingress 对象的 `.metadata.annotations` 字段下的内容实现的。

同样，`.metadata.annotations` 字段下的内容也是由无数的 **Key/Value Pairs(键/值对)** 组成。很多 **Key**都会对应一个 Nginx 的 [**Directives(指令)**]([**Directives(指令)**](Nginx%20 配置详解.md 配置详解.md)ntication](https://kubernetes.github.io/ingress-nginx/user-guide/nginx-configuration/annotations/#authentication) # 认证相关配置

可以为 Nginx 所代理的后端配置一些简单的认证，比如 用户名/密码

**nginx.ingress.kubernetes.io/auth-realm: <STRING>** #

**nginx.ingress.kubernetes.io/auth-secret: <STRING>** #

**nginx.ingress.kubernetes.io/auth-type: <STRING>** #

## [Custom Timeouts](https://kubernetes.github.io/ingress-nginx/user-guide/nginx-configuration/annotations/#custom-timeouts) # 自定义超时时间

配置与 upstream 中定义的服务器的连接超时时间。
**nginx.ingress.kubernetes.io/proxy-connect-timeout: <>**# 对应 Nginx 的 proxy_connect_timeout 指令
**nginx.ingress.kubernetes.io/proxy-send-timeout: <>**#
**nginx.ingress.kubernetes.io/proxy-read-timeout: <>** #
**nginx.ingress.kubernetes.io/proxy-next-upstream: <>** #
**nginx.ingress.kubernetes.io/proxy-next-upstream-timeout: <>** #
**nginx.ingress.kubernetes.io/proxy-next-upstream-tries: <>** #
**nginx.ingress.kubernetes.io/proxy-request-buffering: <>** #

## [Canary](https://kubernetes.github.io/ingress-nginx/user-guide/nginx-configuration/annotations/#canary) # 金丝雀/灰度发布相关配置

通过 canary 相关的注释，我们可以实现金丝雀/灰度发布。即.相同的 host，根据不同的规则，将请求转发给不同的后端。
**nginx.ingress.kubernetes.io/canary: "<BOOLEAN>"** # 是否启用 Canary 功能
其他功能详见 《[实现应用灰度发布](/docs/10.云原生/2.3.Kubernetes%20容器编排系统/8.Kubernetes%20网络/Ingress/Ingress%20Controller/Nginx/实现应用灰度发布.md)》
