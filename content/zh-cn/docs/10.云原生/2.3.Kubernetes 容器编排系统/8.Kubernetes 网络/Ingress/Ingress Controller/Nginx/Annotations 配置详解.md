---
title: Annotations 配置详解
---

# 概述

> 参考：
>
> - [官方文档，用户指南-Annotations](https://kubernetes.github.io/ingress-nginx/user-guide/nginx-configuration/annotations/)

与 ConfigMap 实现配置 Nginx Ingress Controller 运行时行为类似，只不过，Annotations 的方式，是通过设置 Ingress 资源的 `.metadata.annotations` 字段下的内容实现的。

`.metadata.annotations` 字段下的内容也是由无数的 **Key/Value Pairs(键/值对)** 组成。绝大部分 **Key** 都会对应一个 Nginx 的 [**Directives(指令)**](/docs/Web/Nginx/Nginx%20配置详解/Nginx%20配置详解.md#Directives(指令))

Nginx controoler 程序默认读取 Ingress 对象中 `metadata.annotations` 字段下前缀为 `nginx.ingress.kubernetes.io` 的字段，作为运行程序时的配置信息。

> 注意：
> 
> - 所有 Key 都是以 `nginx.ingress.kubernetes.io` 作为前缀，比如配置认证相关，那么 Key 就是 `nginx.ingress.kubernetes.io/auth-realm`
> - 可以为 nginx-ingress-controller 程序添加 `--annotations-prefix` [命令行标志](docs/10.云原生/2.3.Kubernetes%20容器编排系统/8.Kubernetes%20网络/Ingress/Ingress%20Controller/Nginx/命令行标志.md)以改变前缀

# Key 详解

## Authentication - 认证相关配置

https://kubernetes.github.io/ingress-nginx/user-guide/nginx-configuration/annotations/#authentication

可以为 Nginx 所代理的后端配置一些简单的认证，比如 用户名/密码

**nginx.ingress.kubernetes.io/auth-realm**(STRING) #

**nginx.ingress.kubernetes.io/auth-secret**(STRING) #

**nginx.ingress.kubernetes.io/auth-type**(STRING) #

## Backend Protocol - 后端协议配置

https://kubernetes.github.io/ingress-nginx/user-guide/nginx-configuration/annotations/#backend-protocol

使用后端协议注释可以指示 NGINX 应如何与后端服务通信。 

**nginx.ingress.kubernetes.io/backend-protocol**(STRING) # 指示 NGINX 应如何与后端服务通信。`默认值: HTTP`。可用的值：HTTP、HTTPS、GRPC、GRPCS 和 FCGI

- 对应 proxy_pass 指令的值中的协议，比如设置为 HTTPS 时，则生成类似 `proxy_pass https://XXXX` 这种配置。

## Custom Timeouts - 自定义超时时间

https://kubernetes.github.io/ingress-nginx/user-guide/nginx-configuration/annotations/#custom-timeouts

配置与 upstream 中定义的服务器的连接超时时间。

**nginx.ingress.kubernetes.io/proxy-connect-timeout** # 对应 Nginx 的 proxy_connect_timeout 指令

**nginx.ingress.kubernetes.io/proxy-send-timeout**

**nginx.ingress.kubernetes.io/proxy-read-timeout**

**nginx.ingress.kubernetes.io/proxy-next-upstream**

**nginx.ingress.kubernetes.io/proxy-next-upstream-timeout**

**nginx.ingress.kubernetes.io/proxy-next-upstream-tries**

**nginx.ingress.kubernetes.io/proxy-request-buffering**

## Canary - 金丝雀/灰度发布相关配置

https://kubernetes.github.io/ingress-nginx/user-guide/nginx-configuration/annotations/#canary

通过 canary 相关的注释，我们可以实现金丝雀/灰度发布。即.相同的 host，根据不同的规则，将请求转发给不同的后端。

**nginx.ingress.kubernetes.io/canary**(BOOLEAN) # 是否启用 Canary 功能

其他功能详见 《[实现应用灰度发布](/docs/10.云原生/2.3.Kubernetes%20容器编排系统/8.Kubernetes%20网络/Ingress/Ingress%20Controller/Nginx/实现应用灰度发布.md)》
