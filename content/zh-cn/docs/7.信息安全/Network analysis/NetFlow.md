---
title: NetFlow
linkTitle: NetFlow
date: 2024-02-25T18:14
weight: 20
---

# 概述

> 参考：
>
> - [Wiki，NetFlow](https://en.wikipedia.org/wiki/NetFlow)

NetFlow 是 1996 年左右在 Cisco 路由器上引入的一项功能，它能够在 IP 网络流量进入或退出接口时收集该流量。通过分析 NetFlow 提供的数据，网络管理员可以确定流量的来源和目的地、服务类别以及拥塞原因等信息。

NetFlow 通常由三个部分组成

- **Flow exporter(流导出器)** # 将数据包聚合成流，并将流记录导出到一个或多个流收集器。
- **Flow collector(流收集器)** # 负责从流导出器接收的流数据的接收、存储和预处理。
- **Analysis application(分析程序)** # 在入侵检测或流量分析的情况下分析接收到的流数据。

截至 2012 年，大多数防火墙和基于软件的 IP 路由器仍在使用类似于 NetFlow 交换的技术。例如 Linux 使用的 [Netfilter](/docs/1.操作系统/Kernel/Network/Linux%20网络流量控制/Netfilter/Netfilter.md) 框架的 conntrack 功能。

# NetFlow 与 DPI

下面回答来自 ChatGPT 3.5

NetFlow（网络流量流）和 DPI（深度数据包检测）是两种不同的网络监测和分析技术，它们在某些方面有联系，但主要关注的方面有所不同。

**NetFlow:**

- NetFlow是一种网络流量分析技术，通常由Cisco开发和使用。
- 它通过在网络设备（如路由器和交换机）上捕获和记录网络流量的摘要信息来工作。这些信息包括源IP地址、目标IP地址、端口号、协议等。
- NetFlow生成的报告提供了对流量模式、流量源和目标等的高层次摘要，适用于流量分析和网络性能监测。

**DPI:**

- DPI是深度数据包检测的缩写，是一种更为深入地分析网络数据包内容的技术。
- DPI不仅关注流量的源和目标，还深入解析数据包的内容，可以识别应用程序、协议、用户行为等。
- DPI通常用于实时检测网络中的应用层协议，以便对网络流量进行更深入的分析，包括识别特定应用程序或服务，检测潜在的威胁和执行策略。

联系和区别：

- **联系：**
  - NetFlow和DPI都提供了对网络流量的分析能力，但关注的层次和信息的深度不同。
  - 在一些情况下，可以将两者结合使用，通过NetFlow提供的摘要信息定位特定流量，然后使用DPI深入分析该流量的内容。
- **区别：**
  - NetFlow更专注于流量的高层次概要，而DPI更关注深入的数据包内容分析。
  - NetFlow通常用于流量监测、性能优化和容量规划，而DPI主要用于安全分析、应用层识别和用户行为分析。

# 分类

Wiki 中的分类: [Computer network analysis](https://en.wikipedia.org/wiki/Category:Computer_network_analysis)，还分类下还包括 [Network analyers](https://en.wikipedia.org/wiki/Category:Network_analyzers)，这里有 DPI
