---
title: DPI
linkTitle: DPI
date: 2024-02-20T17:01
weight: 2
---

# 概述

> 参考：
> 
> - [Wiki，Deep_packet_inspection](https://en.wikipedia.org/wiki/Deep_packet_inspection)

**Deep packet inspection(深度数据包检测，简称 DPI)** 是一种用于 **Network analysis(网络分析)** 的 行为、技术、方法，它详细检查计算机网络上传输的数据包，并可能根据情况采取警报、阻止、重新路由或记录等行动。可以通过 软件、硬件、软硬结合 多种方式实现 DPI。

> 用白话说：DPI 不是一种特定的技术或协议，而是一种处理流量的方式，凡是可以处理 OSI 模型 4 层往上的程序或设备，都可以称为 DPI 的实现。

可以通过多种手段获取数据包以进行 DPI。比如 **[Port mirroring](https://en.wikipedia.org/wiki/Port_mirroring)(端口镜像)**、etc. ；也可以通过纯硬件设备，比如 [Network tap(网络分流器)](docs/7.信息安全/Network%20analysis/Network%20tap.md)、etc. ，通过这些硬件设备复制数据流并将其发送到分析器工具以检查数据包。

> Notes: 随着技术的发展，DPI 包含的概念也在增加，有时候从广义上来说 DPI 甚至可以看做是由多个组件组成的系统，其中包括的组件甚至可以是 Network tap 以便在复制数据包的同时处理数据包，还不影响原始数据包的传递。

DPI 始于 20 世纪 90 年代。早期的 DPI 实现有：

- [RMON](https://en.wikipedia.org/wiki/RMON "RMON")
- [Sniffer](https://en.wikipedia.org/wiki/Sniffer_(protocol_analyzer) "Sniffer (protocol analyzer)")
- [Wireshark](https://en.wikipedia.org/wiki/Wireshark "Wireshark")
- etc.

**实现 DPI 的可以是硬件或软件**。基本 DPI 功能包括数据包标头和协议字段的分析。

基于 DPI 可以实现的功能

- 僵木蠕系统
- 流控系统
- etc.

# Bypass tap

> 参考：
> 
> - [Wiki，Bypass tap](https://en.wikipedia.org/wiki/Bypass_switch)

**Bypass tap(旁路分路器)** 也称为 Bypass switch(旁路交换机)，在中文环境中有时候也称为 Optical swap(光切换设备，Optiswap)、光开关、等。Bypass tap 是一种硬件设备，可以与安全设备并联并串联到网络链路中，为安全设备提供 **fail-safe access(故障时可安全访问)** 的能力。

Bypass tap 通常至少有 4 个端口。

- A 和 B 两个端口串联，且中间不经过任何电路，Bypass tap 本身不运行的情况下也可以保证 A 到 B 之间的链路畅通。当安全设备正常运行时，A 到 B 的连接是断开的；
- C 和 D 是用来监控安全设备的端口，安全设备正常运行时，流量通过 C 和 D 端口，相当于将安全设备串联进网络中。
- 当检测到安全设备出现异常时，将会切断 C 和 D 的端口，将流量转交给 A 和 B 以保证网络链路上的数据不间断。 

![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/information_security/202402220011706.png)

通常来说，这两种情况可以用两种模式来概括

- 控制模式
- 直通模式

# 最佳实践

https://en.wikipedia.org/wiki/Beam_splitter

通常来说，为了高可用，DPI 等安全设备需要与 Bypass tap 共同使用，以保证 DPI 设备异常时流量不会中断。

在有些不要求高可用的场景中，不需要 DPI 的流量过滤，只需要流量分析的能力时，也可以不用 Bypass tap，而使用市面上的 [Fiber-optic splitter(分光器)](docs/4.数据通信/Networking%20device/Fiber-optic%20splitter.md)，将光的 20% 分到 DPI 或其他流量处理设备中进行后续的流量分析。

# 待整理

基础数据管理

活跃资源管理

访问日志管理

信息安全管理

流量采集管理

数据安全管理

网络安全管理

