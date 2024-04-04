---
title: Network analysis
linkTitle: Network analysis
date: 2024-04-03T11:02
weight: 1
---

# 概述

> 参考：
>
> -

Network analysis(网络分析) 依赖很多基础的流量处理功能以组成完整的系统

- [Port mirroring](/docs/7.信息安全/Network%20analysis/Port%20mirroring.md) # 流量镜像是网络分析的基础。没有端口镜像将流量镜像出来，那么所有的流量分析程序都要使用原始流量，这是绝对不可接受的。除了具体的流量封堵外，任何分析，都应该使用镜像出来的流量。
- etc.

# 最佳实践

## 基于 DPI 的网络分析

在网络分析系统中，[DPI](/docs/7.信息安全/Network%20analysis/DPI.md) 设备除了数据包的检查与处置外，通常还可能具有 [Port mirroring](/docs/7.信息安全/Network%20analysis/Port%20mirroring.md)(端口镜像)、[Fiber-optic splitter](/docs/4.数据通信/Networking%20device/Fiber-optic%20splitter.md)、etc. 相关的流量复制能力，这些被 DPI 复制的流量将会送到 [Network packet broker](/docs/7.信息安全/Network%20analysis/Network%20packet%20broker.md) 设备中以进行聚合、过滤，然后再转发给后端的业务系统。这是一种常见的网络流量分析系统，拿到了流量就相当于有了数据，至于数据如何用，根据具体业务情况而定。

这些业务系统通常包括

- 僵木蠕监测与处置系统
- 流控系统
- 话单系统
- 流量深度分析系统
- etc.

![network_analysis_system.excalidraw](Excalidraw/network_analysis_system.excalidraw.md)
