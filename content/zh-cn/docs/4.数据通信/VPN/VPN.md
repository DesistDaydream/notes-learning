---
title: VPN
created: 2026-07-27T17:23
weight: 1
---

# 概述

> 参考：
>
> - [Wiki, VPN](https://en.wikipedia.org/wiki/Virtual_private_network)

**Virtual Private Network(虚拟专用网络，简称 VPN)** 是一种组网方式，通过 **[Tunneling Protocol](/docs/4.数据通信/Protocol/Tunneling%20Protocol/Tunneling%20Protocol.md)(隧道协议)** 建立的虚拟点对点连接。可以从逻辑上，让人们将通过 VPN 将两个或多个互不连接的网络打通，组成一个更大型的局域网。

可以实现 VPN 的常见 技术、协议、解决方案：

- OpenVPN # 基于 SSL 的 VPN 系统，广泛使用 OpenSSL 加密库和 TLS 协议。
- [libreswan](https://github.com/libreswan/libreswan) # IPsec 服务器
- [xl2tpd](https://github.com/xelerance/xl2tpd) # L2TP 提供者
- [tinc](https://github.com/gsliepen/tinc) 是一个虚拟专用网络 (VPN) 守护程序
- [WireGuard](/docs/4.数据通信/Protocol/Tunneling%20Protocol/WireGuard/WireGuard.md)
- [Tailscale](/docs/4.数据通信/Utility/Tailscale/Tailscale.md)
- https://github.com/EasyTier/EasyTier # 对标 Tailsacle，可以使用 WireGuard 作为基础的 VPN
- etc.

Tips: 有时候越过 GFW(中国不让国内用户访问国外的东西) 的程序也被称为 VPN。只要建立了一层专用的网络，好像都可以套在 VPN 上。
