---
title: DPI
linkTitle: DPI
date: 2024-02-20T17:01
weight: 1
---

# 概述

> 参考：
> 
> - [Wiki，Deep_packet_inspection](https://en.wikipedia.org/wiki/Deep_packet_inspection)

**Deep packet inspection(深度数据包检测，简称 DPI)** 是一种数据处理的技术或方法，它详细检查计算机网络上传输的数据，并可能根据情况采取警报、阻止、重新路由或记录等行动。

可以通过多种手段获取数据包以进行 DPI。比如 **[Port mirroring](https://en.wikipedia.org/wiki/Port_mirroring)(端口镜像)**、etc. ；也可以通过纯硬件设备，比如 **[Network tap](https://en.wikipedia.org/wiki/Network_tap)(网络分流器)**、etc. ，通过这些硬件设备复制数据流并将其发送到分析器工具以检查数据包。

> Notes: 随着技术的发展，DPI 包含的概念也在增加，有时候从广义上来说 DPI 甚至可以看做是由多个组件组成的系统，其中包括的组件甚至可以直接进行流量过滤，而不单单是复制数据包。

DPI 技术拥有悠久且技术先进的历史，始于 20 世纪 90 年代，之后该技术才进入当今常见的主流部署。该技术的根源可以追溯到 30 多年前，当时许多先驱者通过通用标准和早期创新等方式贡献了他们的发明供行业参与者使用，比如：

- [RMON](https://en.wikipedia.org/wiki/RMON "RMON")
- [Sniffer](https://en.wikipedia.org/wiki/Sniffer_(protocol_analyzer) "Sniffer (protocol analyzer)")
- [Wireshark](https://en.wikipedia.org/wiki/Wireshark "Wireshark")
- etc.

**实现 DPI 的可以是硬件或软件**。基本 DPI 功能包括数据包标头和协议字段的分析。

基于 DPI 可以实现很多功能

- 僵木蠕系统
- 流控系统
- etc.

# Network tap

**Network tap(网络分流器)** 是监视本地网络上的事件的系统。通常是一个硬件设备，一般至少具有三个端口：A 端口、B 端口、monitor(监听) 端口。分流器可以让数据在 A 与 B 之间的传输实时无阻碍得通过，同时还将相同的数据复制到 monitor 端口，从而使得第三方分析数据。

> tap 本身有窃听的意思，Network tap 本质上可以算是 网络窃听器、网络监听器


# 待整理

基础数据管理

活跃资源管理

访问日志管理

信息安全管理

流量采集管理

数据安全管理

网络安全管理

