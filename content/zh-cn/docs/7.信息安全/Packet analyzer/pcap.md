---
title: pcap
linkTitle: pcap
weight: 20
---

# 概述

> 参考：
>
> - [Wiki, pcap(包捕获)](https://en.wikipedia.org/wiki/Pcap)

在计算机网络管理领域，**Packet Capture(包捕获，简称 pcap)** 是一个用于捕获网络流量的 **API**。很多数据包分析器都依赖于 pcap 来运行。所以，pcap 准确来说，应该称为 **PCAP API**

- WinPcap # [Microsoft OS](/docs/1.操作系统/Operating%20system/Microsoft%20OS/Microsoft%20OS.md) 下最早的 pcap
  - [Npcap](https://nmap.org/npcap/) # Windows 新的 pcap
- [libpcap](#libpcap) # [Unix-like OS](/docs/1.操作系统/Operating%20system/Unix-like%20OS/Unix-like%20OS.md) 下的 pcap

> Notes: 虽然该名称是 packet capture 的缩写，但这并不是 API 的正确名称。类 Unix 系统在 libpcap 库中实现 pcap；对于 Windows，有一个名为 WinPcap 的 libpcap 端口不再受支持或开发，而对于 Windows 7 及更高版本，仍支持一个名为 Npcap 的端口。

很多实现 pcap 能力的语言若想开发 PCAP 能力必须依赖 libpcap（e.g. go 语言需要开启 CGO_ENABLED=1，且保证系统中安装了 pcap）

# pcap 安装

## Npcap

https://npcap.com/#download

无注意事项，直接安装即可。

## libpcap

> 参考:
>
> - [TCPDump 官方文档，pcap](https://www.tcpdump.org/manpages/pcap.3pcap.html)

libpcap 是 [Unix-like OS](/docs/1.操作系统/Operating%20system/Unix-like%20OS/Unix-like%20OS.md) 下的 pcap C 库，若想基于 libpcap 进行开发，通常需要在系统中安装 [C](/docs/2.编程/高级编程语言/C/C.md) 语言的头文件（i.e. libpcap 的开发库）

- Ubuntu 系系统
  - `apt install libpcap-dev`
- RedHat 系系统
  - `yum install libpcap-devel`

libpcap 与 [TCPDump](/docs/7.信息安全/Packet%20analyzer/TCPDump/TCPDump.md) 项目一起进行维护

# 其它抓包工具

**netcap** # 一个基于 [BCC 工具集](/docs/1.操作系统/Linux%20管理/Linux%20系统管理工具/BCC%20工具集.md)，可对含skb作为参数的系统函数，或者基于 [DPDK](/docs/4.数据通信/DPDK/DPDK.md) 的mbuf抓包分析的工具。

- https://github.com/bytedance/netcap
- [公众号 - 深入浅出BPF，字节跳动开源Linux内核网络抓包工具netcap](https://mp.weixin.qq.com/s/ezMEjx3mzHOaxvfMQwkn9g)
