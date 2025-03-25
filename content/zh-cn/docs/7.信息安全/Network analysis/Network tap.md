---
title: Network tap
linkTitle: Network tap
weight: 5
---

# 概述

> 参考：
>
> - [Wiki, Network tap](https://en.wikipedia.org/wiki/Network_tap)

**Network tap(网络分流器)** 是一种通过类似 [Port mirroring](/docs/7.信息安全/Network%20analysis/Port%20mirroring.md) 技术实现的流量监听方式。实现了 Network tap 能力的硬件设备可以用来监视本地网络上的事件。Network tap 设备一般至少有 3 个端口：A 端口、B 端口、monitor(监听) 端口。Network tap 可以让数据在 A 与 B 之间的传输实时无阻碍得通过，同时还将相同的数据复制到 monitor 端口，从而使得第三方分析数据。

> [!tips]
> tap 本身有窃听的意思，Network tap 本质上可以算是 网络窃听器、网络监听器。厂商们为了更好的售卖产品而美化了该词的含义，而且也根据 Tap 设备具体的功能（复制流量分到另一个端口就相当于分流了），把 tap 的含义改为了分流。还扩展了一下 tap 这个词成 terminal access point(终端接入点)

Network tap 在狭义上通常只有 1 分 多 的能力，适用在简单的流量不大的网络环境中，若是在流量巨大且复杂的网络中，想实现 多 分 多的能力（汇聚与分流），则需要使用 [Network packet broker](/docs/7.信息安全/Network%20analysis/Network%20packet%20broker.md) 技术。

# Network tap 与 Network packet broker

Network tap 通常是在小流量场景中串接进链路中；而 Network packet broker 则通常是在大流量、复杂的网络场景中，并接到链路上。

在 多 分 多 的场景中，数据来源是多条链路，如何处理多条链路中的流量而不会产生混乱，则是 Network packet broker 需要解决的重要问题。

# Network tap 与 Fiber-optic splitter

Network Tap 和 [Fiber-optic splitter](/docs/4.数据通信/Networking%20device/Fiber-optic%20splitter.md) 在网络中起着类似但不完全相同的作用。都是用来实现类似流量复制的逻辑。

- 作用的层次
  - 分光器：物理层（i.e. L1 层  ）。基于物理规则对光子产生影响。
  - 分流器：数据链路层及以上（L2 层  - L7 层）。
- 是否有源（i.e. 是否有 5-tuple 需要处理）
  - 分光器：无源
  - 分流器：有源
- 功能
  - 分光器：只能做 1 对多全流量复制
  - 分流器：支持更加智能的功能，比如流量汇聚、负载均衡、报文过滤、报文编辑、流量分发（需要 NPB 的能力）

在一些情况下，Network tap 和 光纤分路器 可以结合使用。光纤分路器可以用于将光信号分发到 Network tap 设备，从而使得网络流量可以被复制到多个目的地进行监视和分析。
