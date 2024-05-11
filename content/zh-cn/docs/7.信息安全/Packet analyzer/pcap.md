---
title: pcap
linkTitle: pcap
date: 2024-02-22T17:27
weight: 20
---

# 概述

> 参考：
> 
> - [Wiki，pcap(包捕获)](https://en.wikipedia.org/wiki/Pcap)

在计算机网络管理领域，**Packet Capture(包捕获，简称 pcap)** 是一个用于捕获网络流量的 **API**。很多数据包分析器都依赖于 pcap 来运行。所以，pcap 准确来说，应该称为 **PCAP API**

- WinPcap # Windows 系统下最早的 pcap
- [Npcap](https://nmap.org/npcap/) # Windows 新的 pcap
- [libpcap](https://www.tcpdump.org/) # 类 Unix 系统下的 pcap

> Notes: 虽然该名称是 packet capture 的缩写，但这并不是 API 的正确名称。类 Unix 系统在 libpcap 库中实现 pcap；对于 Windows，有一个名为 WinPcap 的 libpcap 端口不再受支持或开发，而对于 Windows 7 及更高版本，仍支持一个名为 Npcap 的端口。

# pcap 安装

## Npcap

https://npcap.com/#download

无注意事项，直接安装即可。
