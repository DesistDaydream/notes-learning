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

**Deep packet inspection(深度数据包检测，简称 DPI)** 是一种数据处理的 行为、技术、方法，它详细检查计算机网络上传输的数据，并可能根据情况采取警报、阻止、重新路由或记录等行动。

可以通过多种手段获取数据包以进行 DPI。比如 **[Port mirroring](https://en.wikipedia.org/wiki/Port_mirroring)(端口镜像)**、etc. ；也可以通过纯硬件设备，比如 **[Network tap](https://en.wikipedia.org/wiki/Network_tap)(网络分流器)**、etc. ，通过这些硬件设备复制数据流并将其发送到分析器工具以检查数据包。

> Notes: 随着技术的发展，DPI 包含的概念也在增加，有时候从广义上来说 DPI 甚至可以看做是由多个组件组成的系统，其中包括的组件甚至可以直接进行流量过滤，而不单单是复制数据包。

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

## Network tap

> 参考：
> 
> - [Wiki，Network tap](https://en.wikipedia.org/wiki/Network_tap)

**Network tap(网络分流器)** 是监视本地网络上的事件的系统。通常是一个硬件设备，一般至少有 3 个端口：A 端口、B 端口、monitor(监听) 端口。分流器可以让数据在 A 与 B 之间的传输实时无阻碍得通过，同时还将相同的数据复制到 monitor 端口，从而使得第三方分析数据。

> tap 本身有窃听的意思，Network tap 本质上可以算是 网络窃听器、网络监听器

### Network tap 与 Fiber-optic splitter

下面是来自 ChatGPT 3.5 的回答：Network Tap 和 [Fiber-optic splitter](docs/4.数据通信/Networking%20device/Fiber-optic%20splitter.md) 在网络中起着类似但不完全相同的作用

- **功能相似：** Network Tap 设备和光纤分路器都用于在网络中监视和管理流量。它们都允许用户获取网络流量并将其转发到指定的目的地进行分析、存储或其他处理。
- **不同之处：** 光纤分路器主要用于光纤网络中的信号分发，将光信号从一个源头分发到多个目的地。而 Network Tap 设备则是一种专门设计用于复制网络数据流的设备，通常用于复制传统以太网网络中的数据流，使得监视和分析网络流量变得更容易。
- **配合使用：** 在一些情况下，Network Tap 设备和光纤分路器可以结合使用，特别是在光纤网络中。光纤分路器可以用于将光信号分发到多个 Network Tap 设备，从而使得网络流量可以被复制到多个目的地进行监视和分析。

综上所述，Network Tap 和 Fiber-optic splitter 都是网络中重要的监视和管理工具，它们可以在光纤网络中配合使用，以实现更全面和高效的网络流量监视和分析。

我个人理解，Fiber-optic splitter 的底层原理是物理规则对光子的控制；而 Network tap 可以控制的流量不是单只光纤的流量，还包含电口中的流量，本质是通过设备将传入的流量完全复制一份到另一个端口传出。

## 最佳实践

有时候 Network tap 镜像的流量并不需要转发到其他的流量分析工具中，而是让一个 DPI 设备具有 Network tap、流量分析、etc. 一系列功能，组成一个完整而全面的安全设备。

# Bypass tap

> 参考：
> 
> - [Wiki，Bypass tap](https://en.wikipedia.org/wiki/Bypass_switch)

**Bypass tap(旁路分路器)** 也称为 Bypass switch(旁路交换机)，在中文环境中有时候也称为 Optical swap(光切换设备，Optiswap)、光开关、等。Bypass tap 是一种硬件设备，可以与安全设备并联并串联到网络链路中，为安全设备提供 **fail-safe access(故障时可安全访问)** 的能力。

Bypass tap 通常至少有 4 个端口。

- A 和 B 两个端口串联，且中间不经过任何电路，Bypass tap 本身不运行的情况下也可以保证 A 到 B 之间的链路畅通。当安全设备正常运行时，A 到 B 的连接是断开的；
- C 和 D 是用来监控安全设备的端口，安全设备正常运行时，流量通过 C 和 D 端口，相当于将安全设备串联进网络中。
- 当检测到安全设备出现异常时，将会切断 C 和 D 的端口，将流量转交给 A 和 B 以保证网络链路上的数据不间断。 

![image.png|500](https://notes-learning.oss-cn-beijing.aliyuncs.com/information_security/202402220011706.png)

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

