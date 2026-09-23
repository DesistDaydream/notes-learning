---
title: "Tailcat"
created: "2026-09-20T11:24"
weight: 100
---

# 概述

> 参考：
>
> - [GitHub 项目，tailscale/tailcat](https://github.com/tailscale/tailcat)

Tailcat 类似 [Netcat](/docs/4.数据通信/Utility/Netcat.md)，但是通过 [Tailscale](/docs/4.数据通信/Utility/Tailscale/Tailscale.md) 的数据平面运行，而无需 Tailscale 的控制平面（i.e. 人话，直接通过 DERP 建立简单的通路）。

# 关联文件与配置

**~/.config/tailcat/keys/default.private.json** #