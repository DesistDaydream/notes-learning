---
title: Port mirroring
linkTitle: Port mirroring
weight: 2
---

# 概述

> 参考：
>
> - [Wiki, Port mirroring](https://en.wikipedia.org/wiki/Port_mirroring)

**Port mirroring(端口镜像)** 是交换机上的一个功能，可以将一个端口上经过的所有流量复制一份到另一个端口上。这种方式通常用于需要监控网络流量的环境中，原始流量不受影响，被复制的流量发往后端设备以进一步分析，比如 IDS(入侵监测系统)、etc. 。Cisco 公司生成的交换机上的 Port mirroring 功能称为 **Switched Port Analyzer(简称 SPAN)** 或 **Remote Switched Port Analyzer(简称 RSPAN)**，有时候称为 Span Port。

网络工程师或管理员使用端口镜像来分析和调试数据或诊断网络上的错误。它可以帮助管理员密切关注网络性能并在出现问题时向他们发出警报。它可用于镜像单个或多个接口上的入站或出站流量（或两者）。

随着网络分析需求的增加，很多时候也会把 Port mirroring 形容成流量镜像的功能，只有具备了 Port mirroring 的系统，才能真正实现 [Network analysis](/docs/7.信息安全/Network%20analysis/Network%20analysis.md)。

# SPAN

**Switched Port Analyzer(简称 SPAN)**
