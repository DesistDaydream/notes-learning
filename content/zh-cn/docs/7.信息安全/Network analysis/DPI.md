---
title: DPI
linkTitle: DPI
date: 2024-02-20T17:01
weight: 2
---

# 概述

> 参考：
>
> - [Wiki，Deep packet inspection](https://en.wikipedia.org/wiki/Deep_packet_inspection)

**Deep packet inspection(深度数据包检测，简称 DPI)** 是一种用于 **Network analysis(网络分析)** 的 行为、技术、方法，它详细检查计算机网络上传输的数据包，并可能根据情况采取警报、阻止、重新路由或记录等行动。可以通过 软件、硬件、软硬结合 多种方式实现 DPI。

> 用白话说：DPI 不是一种特定的技术或协议，而是一种处理流量的方式。

可以通过多种手段获取数据包以进行 DPI。比如 [Port mirroring(端口镜像)](https://en.wikipedia.org/wiki/Port_mirroring)、etc. ；也可以通过纯硬件设备，比如 [Network tap(网络分流器)](/docs/7.信息安全/Network%20analysis/Network%20tap.md)、etc. ，通过这些硬件设备复制数据流并将其发送到分析器工具以检查数据包。

> Notes: 随着技术的发展，DPI 包含的概念也在增加，有时候从广义上来说 DPI 甚至可以看做是由多个组件组成的系统，其中包括的组件甚至可以是 Network tap 以便在复制数据包的同时处理数据包，还不影响原始数据包的传递。

DPI 始于 20 世纪 90 年代。早期的 DPI 实现有：

- [RMON](https://en.wikipedia.org/wiki/RMON "RMON")
- [Sniffer](https://en.wikipedia.org/wiki/Sniffer_(protocol_analyzer) "Sniffer (protocol analyzer)")
- [Wireshark](https://en.wikipedia.org/wiki/Wireshark "Wireshark")
- etc.

支持 DPI 的设备或程序能够查看 OSI 模型的第 2 层和第 3 层以外的情况。在某些情况下，可以调用 DPI 来查看 OSI 模型的第 2-7 层。这包括标头和数据协议结构以及消息的有效负载。当设备根据 OSI 模型第 3 层之外的信息查找或采取其他操作时，将调用 DPI 功能。 DPI 可以识别的数据包特征又很多：

- 协议
- HTTP 请求类型
- etc.

在许多情况下，End points 可以利用加密或混淆技术来逃避 DPI 的识别。

**实现 DPI 的可以是硬件或软件**。基本 DPI 功能包括数据包标头和协议字段的分析。

基于 DPI 可以实现的功能

- 僵木蠕系统
- 流控系统
- etc.

# 最佳实践

https://en.wikipedia.org/wiki/Beam_splitter

通常来说，为了保证高可用，DPI 等安全设备需要与 [Bypass tap](/docs/7.信息安全/Network%20analysis/Bypass%20tap.md) 共同使用，以保证 DPI 设备异常时流量不会中断。

在有些不要求高可用的场景中，不需要 DPI 的流量过滤，只需要流量分析的能力时，也可以不用 Bypass tap，而使用市面上的 [Fiber-optic splitter(分光器)](/docs/4.数据通信/Networking%20device/Fiber-optic%20splitter.md)，将光的 20% 分到 DPI 或其他流量处理设备中进行后续的流量分析。

# Deep/Dynamic Flow Inspection

https://www.telecomtrainer.com/dfi-deep-flow-inspection/

**Deep/Dynamic Flow Inspection(深度/动态流量监测，简称 DFI)** 没有 Wiki，也没有什么官方说明。通常是做 DFI 企业内部交流时使用的术语

由于 DPI 虽然识别精度高，但需要对特定协议的应用层特征被动跟踪，所以对加密或混淆的流量识别比较困难。此时需要使用 DFI。

DFI 依赖一个“模型”（类似 AI 中的模型的概念），通过流量特征进行识别，比如 P2P 流量。
