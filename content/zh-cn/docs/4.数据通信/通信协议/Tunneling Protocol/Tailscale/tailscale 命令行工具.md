---
title: tailscale 命令行工具
linkTitle: tailscale 命令行工具
date: 2024-03-21T14:12
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
