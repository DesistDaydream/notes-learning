---
title: Packet analyzer
linkTitle: Packet analyzer
date: 2024-01-14T19:03
weight: 1
---

# 概述

> 参考：
>
> - [Wiki，Packet analyzer(包分析器)](https://en.wikipedia.org/wiki/Packet_analyzer)
> - [Wiki，pcap(包捕获)](https://en.wikipedia.org/wiki/Pcap)

**Packet analyzer(包分析器)** 是一种计算器程序或计算机硬件，可以拦截和记录通过计算机网络的流量，有的地方也称之为 **Packet sniffer(包嗅探器)**。数据包捕获是拦截和记录流量的过程。随着数据流跨网络流流，分析器捕获每个数据包，如果需要，可以解码分组的原始数据，显示分组中的各种字段的值，并根据适当的 [RFC](docs/x_标准化/Internet/IETF.md) 或其他规范分析其内容。

## Packet Analyzer 的实现

各种实现的对比：<https://en.wikipedia.org/wiki/Comparison_of_packet_analyzers>

- tcpdump
- Wireshark
- ......等等

# pcap 概述

在计算机网络管理领域，**Packet Capture(包捕获，简称 pcap)** 是一个用于捕获网络流量的 **API**。很多数据包分析器都依赖于 pcap 来运行。所以，pcap 准确来说，应该称为 **PCAP API**

- WinPcap # Windows 系统下最早的 pcap
- [Npcap](https://nmap.org/npcap/) # Windows 新的 pcap
- [libpcap](https://www.tcpdump.org/) # 类 Unix 系统下的 pcap

# 抓包工具

Fiddler

- 官网：<https://www.telerik.com/fiddler>

[Charles](docs/7.信息安全/Packet%20analyzer/Charles.md)

HTTP Debugger

- https://www.httpdebugger.com/