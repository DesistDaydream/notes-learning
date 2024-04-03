---
title: Network tap
linkTitle: Network tap
date: 2024-02-23T11:52
weight: 20
---

# 概述

> 参考：
>
> - [Wiki，Network tap](https://en.wikipedia.org/wiki/Network_tap)

**Network tap(网络分流器)** 是监视本地网络上的事件的系统。通常是一个硬件设备，一般至少有 3 个端口：A 端口、B 端口、monitor(监听) 端口。分流器可以让数据在 A 与 B 之间的传输实时无阻碍得通过，同时还将相同的数据复制到 monitor 端口，从而使得第三方分析数据。

> tap 本身有窃听的意思，Network tap 本质上可以算是 网络窃听器、网络监听器。为了美化一下，也根据具体的功能（复制流量就相当于分流了），把 tap 的含义改为了分流。

Network tap 通常只有 1 分 2 的能力，若是想实现 1 分 多，多 分 多的能力，可以使用 [Network Packet Broker](docs/7.信息安全/Network%20analysis/Network%20Packet%20Broker.md) 技术。

## 实现 Network tap 的技术

有多种监控网络的方法。根据网络技术、监控目标、可用资源和目标网络的规模，可以使用多种 tap 方法

- [Port mirroring](docs/7.信息安全/Network%20analysis/Port%20mirroring.md)
- etc.

# Network tap 与 Fiber-optic splitter

下面是来自 ChatGPT 3.5 的回答：

Network Tap 和 [Fiber-optic splitter](/docs/4.数据通信/Networking%20device/Fiber-optic%20splitter.md) 在网络中起着类似但不完全相同的作用

- **功能相似：** Network Tap 设备和光纤分路器都用于在网络中监视和管理流量。它们都允许用户获取网络流量并将其转发到指定的目的地进行分析、存储或其他处理。
- **不同之处：** 光纤分路器主要用于光纤网络中的信号分发，将光信号从一个源头分发到多个目的地。而 Network Tap 设备则是一种专门设计用于复制网络数据流的设备，通常用于复制传统以太网网络中的数据流，使得监视和分析网络流量变得更容易。

我个人理解，Fiber-optic splitter 的底层原理是物理规则对光子的控制；而 Network tap 可以控制的流量不是单只光纤的流量，还包含电口中的流量，本质是通过设备将传入的流量完全复制一份到另一个端口传出。

- 作用的层次
  - 分光器：物理层（i.e. L1 层  ）
  - 分流器：L2 - L4 层
- 是否有源（i.e. 是否有 5-tuple 需要处理）
  - 分光器：无源
  - 分流器：有源
- 功能
  - 分光器：只能做 1 对多全流量复制  
  - 分流器：支持更加智能的功能，比如流量汇聚、负载均衡、报文过滤、报文编辑、流量分发（需要 NPB 的能力）

在一些情况下，Network tap 和 光纤分路器 可以结合使用，特别是在光纤网络中。光纤分路器可以用于将光信号分发到多个 Network tap 设备，从而使得网络流量可以被复制到多个目的地进行监视和分析。
