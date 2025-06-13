---
title: tailscale CLI
linkTitle: tailscale CLI
weight: 20
---

# 概述

> 参考：
>
> - [官方文档，参考 - CLI](https://tailscale.com/kb/1080/cli)

# Syntax(语法)

**tailscale COMMAND**

- **netcheck** # 打印本地网络状况分析。主要是显示当前可用的 DERP 以及本机到各个 DERP 的连接延迟 等信息。
- **ping \<HOST>** # 通过 Tailscale ping 指定主机，看看本机是如何路由到目标的（是直通还是经过了 DERP）
- **set** # 改变指定的配置。

# up

# set

**--accept-routes** # 是否接受其他节点公开的路由信息。`默认值: false`

- Tips: 对应 /var/lib/tailscale/tailscaled.state 文件中 `.profile-XXX` 字段中的 `.RouteAll` 字段

**--advertise-routes** # 向整个 Tailscale 网络公开本机的路由。也就是说告诉其他节点访问哪些 IP 要经过本机。`默认值: 空`，值是以 `,` 分割的 CIDR 格式的子网

# EXAMPLE
