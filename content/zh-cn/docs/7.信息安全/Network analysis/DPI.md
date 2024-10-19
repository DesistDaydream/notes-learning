---
title: DPI
linkTitle: DPI
date: 2024-02-20T17:01
weight: 3
---

# 概述

> 参考：
>
> - [Wiki, Deep packet inspection](https://en.wikipedia.org/wiki/Deep_packet_inspection)

**Deep packet inspection(深度数据包检测，简称 DPI)** 是一种用于 **Network analysis(网络分析)** 的 行为、技术、方法，它详细检查计算机网络上传输的数据包，并可能根据情况采取 发送警报、阻止数据包传输、重新路由、记录(或者说镜像)流量 等行动。DPI 的这些功能通常用于 确定应用程序行为基线、分析网络使用情况、排除网络性能故障、确保数据格式正确、检查恶意代码、窃听和互联网审查、以及其他目的。

可以通过 软件、硬件、软硬结合 多种方式实现 DPI。

支持 DPI 的设备或程序能够查看 OSI 模型的第 2 层和第 3 层以外的情况。在某些情况下，可以调用 DPI 来查看 OSI 模型的第 2-7 层。这包括标头和数据协议结构以及消息的有效负载。当设备根据 OSI 模型第 3 层之外的信息查找或采取其他操作时，将调用 DPI 功能。 DPI 可以识别的数据包特征又很多：

- 协议
- HTTP 请求类型
- 数据包中的数据特征
- etc.

IP数据包有多个报头；网络设备只需要使用其中第一个（IP 标头）即可正常操作，但使用第二个标头（例如 TCP 或 UDP）通常被认为是  **Stateful packet inspection(浅层数据包检查)**（通常称为状态数据包检查）

> 用白话说：DPI 不是一种特定的技术或协议，而是一种处理流量的方式。

DPI 始于 20 世纪 90 年代。早期的 DPI 实现有：

- [RMON](https://en.wikipedia.org/wiki/RMON "RMON")
- [Sniffer](https://en.wikipedia.org/wiki/Sniffer_(protocol_analyzer) "Sniffer (protocol analyzer)")
- [Wireshark](https://en.wikipedia.org/wiki/Wireshark "Wireshark")
- etc.

在许多情况下，End points 可以利用加密或混淆技术来逃避 DPI 的识别。

实现 DPI 的设备通常都会 **串联** 进现有的网络链路中，以便实现阻止或丢弃数据包的能力。为了保证高可用，串联到链路中的 DPI 设备需要与 [Bypass swtich](/docs/7.信息安全/Network%20analysis/Bypass%20swtich.md) 共同使用，以保证 DPI 设备异常时流量不会中断。

![bypass.drawio.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/information_security/bypass_and_dpi_1.png)

跟复杂网络环境中，DPI 设备也许还会嵌入流量复制的能力，架构可以参考 [Network analysis](/docs/7.信息安全/Network%20analysis/Network%20analysis.md) 中的最佳实践部分。

# Deep/Dynamic Flow Inspection

https://www.telecomtrainer.com/dfi-deep-flow-inspection/

**Deep/Dynamic Flow Inspection(深度/动态流量监测，简称 DFI)** 没有 Wiki, 也没有什么官方说明。通常是做 DFI 企业内部交流时使用的术语

由于 DPI 虽然识别精度高，但需要对特定协议的应用层特征被动跟踪，所以对加密或混淆的流量识别比较困难。此时需要使用 DFI。

DFI 依赖一个“模型”（类似 AI 中的模型的概念），通过流量特征进行识别，比如 P2P 流量。
