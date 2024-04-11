---
title: WireShark
linkTitle: WireShark
date: 2023-11-01T21:03
weight: 1
---

# 概述

> 参考
>
> - [官网](https://www.wireshark.org/)
> - [官方文档，用户指南](https://www.wireshark.org/docs/wsug_html_chunked/)
> - <https://help.aliyun.com/document_detail/112990.html>(Wireshark 常见提示)
> - <https://blog.csdn.net/qq_15437629/article/details/116565673>
> - [公众号-马哥 Linux 运维，8 个常用的 Wireshark 使用技巧](https://mp.weixin.qq.com/s/yWuDodOpCClZT36yVBeeaQ)
> - [公众号-小林 coding，一文搞定 Wireshark 网络数据包分析](https://mp.weixin.qq.com/s/hL96imOvuodILIhI70fbTg)（Notes: 一文搞不定）

Wireshark 除了可以抓包外，还提供了可视化分析网络包的图形页面，同时，还内置了一系列的汇总分析工具。

[TCPDump](/docs/7.信息安全/Packet%20analyzer/TCPDump/TCPDump.md) 和 Wireshark，这两大利器把我们“看不见”的数据包，呈现在我们眼前，一目了然。这两个工具就是最常用的网络抓包和分析工具，更是分析网络性能必不可少的利器。

- tcpdump 仅支持命令行格式使用，常用在 Linux 服务器中抓取和分析网络包。
- Wireshark 除了可以抓包外，还提供了可视化分析网络包的图形页面。

tcpdump 虽然功能强大，但是输出的格式并不直观。所以，在工作中 tcpdump 只是用来抓取数据包，不用来分析数据包，而是把 tcpdump 抓取的数据包保存成 pcap 后缀的文件，接着用 Wireshark 工具进行数据包的可视化分析。

## Wireshark 工具如何分析数据包？

我们使用 `tcpdump -i eht1 icmp and host 183.232.231.174 -w ping.pcap` 命令，把抓取的数据包保存到 ping.pcap 文件

接着把 ping.pcap 文件拖到电脑，再用 Wireshark 打开它。打开后，你就可以看到下面这个界面：

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/wireshark/1616160823025-9f523c60-bfc7-4c13-a9e4-87589e631637.jpeg)

在 Wireshark 的页面里，可以更加直观的分析数据包，不仅展示各个网络包的头部信息，还会用不同的颜色来区分不同的协议，由于这次抓包只有 ICMP 协议，所以只有紫色的条目。

接着，在网络包列表中选择某一个网络包后，在其下面的网络包详情中，可以更清楚的看到，这个网络包在协议栈各层的详细信息。比如，以编号 1 的网络包为例子：

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/wireshark/1616160823020-a671d5ea-af21-4aee-911e-1c07ec76848c.jpeg)

- 可以在数据链路层，看到 MAC 包头信息，如源 MAC 地址和目标 MAC 地址等字段；
- 可以在 IP 层，看到 IP 包头信息，如源 IP 地址和目标 IP 地址、TTL、IP 包长度、协议等 IP 协议各个字段的数值和含义；
- 可以在 ICMP 层，看到 ICMP 包头信息，比如 Type、Code 等 ICMP 协议各个字段的数值和含义；

Wireshark 用了分层的方式，展示了各个层的包头信息，把“不可见”的数据包，清清楚楚的展示了给我们，还有理由学不好计算机网络吗？是不是相见恨晚？

从 ping 的例子中，我们可以看到网络分层就像有序的分工，每一层都有自己的责任范围和信息，上层协议完成工作后就交给下一层，最终形成一个完整的网络包。

# WireShark 安装

WireShark 依赖 [pcap](/docs/7.信息安全/Packet%20analyzer/pcap.md)，若使用 WireShark 便携包，那么需要手动安装 [Npcap](https://npcap.com/)。

# 过滤语法

`!tcp.analysis.flags` # 去掉 Bad TCP 的包

# Following Protocol Streams(追踪协议流)

https://www.wireshark.org/docs/wsug_html_chunked/ChAdvFollowStreamSection.html#ChAdvFollowStream

WireShark 将具有某些相同特征的多个数据包分为一组，称为 Protocol Stream(协议流)，每条 stream 都有一个 Stream ID。若是 TCP 的包，则是 TCP Protocol Streams；若是 HTTP 的包，则是 HTTP Protocol Stream；以此类推。

想要将多个数据包分为一组，通常要有一些前提条件，比如 TCP 流要保证这些数据包的 src port 和 dest port 且 src ip 和 dest ip 完全相同、etc. 

现在 WireShark 可以为如下追踪如下协议的流：

- DCCP
- HTTP
- HTTP/2
- QUIC
- SIP Call
- TCP
- TLS
- UDP
- WebSocket

# Statistics(统计)

https://www.wireshark.org/docs/wsug_html_chunked/ChStatistics.html

WireShark 提供了广泛的网络统计信息，可以通过 “Statistics(统计)” 菜单访问。

这些统计信息的范围从有关加载的捕获文件的一般信息（如捕获的数据包的数量）到有关特定协议的统计信息（例如，有关捕获的 HTTP 请求和响应数量的统计信息）。

## Conversations(会话)

https://www.wireshark.org/docs/wsug_html_chunked/ChStatConversations.html

**Conversations(会话)** 是两个指定 Endpoint 之间的流量。例如，IP 会话是两个 IP 地址之间的所有流量、Ethernet 会话是两个 Mac 地址之间的流量、etc.

会话统计窗口中的标签会根据标签名的特征合并数据包，以形成统计信息，合并后的数据包将会被到同一个组中

- Ethernet 标签 # 合并 src mac addr 和 dest mac addr 相同的数据包
- IPv4 标签 # 合并 src ip 和 dest ip 相同的数据包
- UDP 标签/TCP 标签 # 合并 src port 和 dest port 相同；且 src ip 与 dest ip 相同的数据包，每个组的数据包当做同一条 Protocol Stream(协议流)，具有相同的 StreamID。
  - 比如，若总共 2000 个包，其中有一些包的 src ip:port 和 dest ip:port 一样，那么该统计-会话中的 UDP 的协议流将会小于 2000

