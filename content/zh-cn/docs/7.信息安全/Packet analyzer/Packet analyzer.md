---
title: Packet analyzer
linkTitle: Packet analyzer
weight: 1
---

# 概述

> 参考：
>
> - [Wiki, Packet analyzer(包分析器)](https://en.wikipedia.org/wiki/Packet_analyzer)

**Packet analyzer(包分析器)** 是一种计算器程序或计算机硬件，可以拦截和记录通过计算机网络的流量，有的地方也称之为 **Packet sniffer(包嗅探器)**。数据包捕获是拦截和记录流量的过程。随着数据流跨网络流流，分析器捕获每个数据包，如果需要，可以解码分组的原始数据，显示分组中的各种字段的值，并根据适当的 [RFC](/docs/Standard/Internet/IETF.md) 或其他规范分析其内容。

## Packet Analyzer 的实现

各种实现的对比: https://en.wikipedia.org/wiki/Comparison_of_packet_analyzers

- [TCPDump](/docs/7.信息安全/Packet%20analyzer/TCPDump/TCPDump.md)
- [WireShark](/docs/7.信息安全/Packet%20analyzer/WireShark/WireShark.md)
- ......等等

# 抓包工具

Reqable

- https://github.com/reqable/reqable-app # 非开源，只是有个仓库
- 官网 https://reqable.com/
- 图标是 小黄鸟，有 移动端  和 PC 端。宣传自己是 Fiddler + Charles + Postman

[Fiddler](/docs/7.信息安全/Packet%20analyzer/Fiddler.md)

[Charles](/docs/7.信息安全/Packet%20analyzer/Charles.md)

mitmproxy

- [GitHub 项目，mitmproxy/mitmproxy](github.com/mitmproxy/mitmproxy)
- Python 编写，为渗透测试人员和软件开发人员提供的交互式、支持 TLS 的拦截 HTTP 代理。

HTTP Debugger

- https://www.httpdebugger.com/
- 可以抓进程的包，而不是通过代理的方式抓包

openQPA

- https://github.com/l7dpi/openQPA, https://gitee.com/l7dpi/openQPA
- http://www.l7dpi.com/
- 基于进程抓包

[SunnyNetTools](https://github.com/qtgolang/SunnyNetTools)

- Sunny网络中间件-抓包工具
