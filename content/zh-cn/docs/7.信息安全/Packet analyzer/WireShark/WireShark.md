---
title: WireShark
linkTitle: WireShark
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

## Wireshark 视图

https://www.wireshark.org/docs/wsug_html_chunked/ChUseViewMenuSection.html

我们使用 `tcpdump -i any port 10443 -nn -w demo-http.pcap` 命令，把抓取的数据包保存到 demo-http.pcap 文件，再用 Wireshark 打开它，可以看到三个主要的窗口

- **Packet List(包列表)**
- **Packet Details(包详情)**
- **Packet Bytes(包字节流)**

**Packet List 窗口**：

![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/wireshark/view_packet_list.png)

在 Wireshark 的页面里，可以更加直观的分析数据包，不仅展示各个网络包的头部信息，还会用不同的颜色来区分不同的协议。

**Packet Details 窗口**：

![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/wireshark/view_pakcet_details.png)

在 包列表 窗口中选择某一个网络包后，在 包详情 窗口中，可以更清楚的看到这个网络包在协议栈各层的详细信息，主要是各种协议的 Header 信息和 Payload 信息。e.g. 以编号为 4 的 HTTP GET 的网络包为例子：

- 可以在数据链路层，看到 MAC 包头信息，如源 MAC 地址和目标 MAC 地址等字段；
- 可以在 IP 层，看到 IP 包头信息，如源 IP 地址和目标 IP 地址、TTL、IP 包长度、协议等 IP 协议各个字段的数值和含义；
- 可以在 TCP 层，看到 TCP 包头信息，比如 Flags、Window、etc. 的数值和含义；

**Packet Bytes 窗口**

![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/wireshark/view_pakcet_bytes.png)

与 Pakcet Details 窗口不同，Packet Bytes 窗口展示的是该包的最原始数据内容，i.e. 数据的 Bytes 信息，并没有那些 WireShark 附加的描述性信息。整体分两大部分：

- 左侧是默认以十六进制表示的字节流信息（右键点击可以更改展示的类型），每个 Bytes 占 8 bit；
- 右侧是字节流中每个字节对应 [ASCII 表](/docs/8.通用技术/编码与解码/字符的编码与解码/ASCII%20表.md) 的字符。

e.g. 图中粉框中的 `47 45 54` 对应的就 ASCII 表的 `GET` 这几个字符。

右侧有些杂乱无章的字符是因为一个包中的 Header 的每个部分的信息并不总是占用 8bit

- e.g. [IP Header](/docs/4.数据通信/Protocol/TCP_IP/IP/IP%20Header.md) 中的 version 长度占用 4bit，0100（e.g. 十进制的 4）；HL 长度占用 4bit，0101（e.g. 十进制的 4）
- 连在一起是 01000101 对应的十六进制为 45，对应的 ASCII 字符为 E。
- 所以这里面的 E 并不是指报文里包含了 E 字符，只是一种以每个 Bytes(8bit) 进行分割的展示方式

> 在我们使用 [TCPDump](/docs/7.信息安全/Packet%20analyzer/TCPDump/TCPDump.md) 抓 HTTP 包的时候，可以看到这种命令: `tcpdump -nn 'tcp[20:2]=0x4745'`，其中 `tcp[20:2]` 表示从 TCP 头的第 20 个字节开始，读 2 个字节的数据。图中蓝色背景选中 20 个字节是 TCP Header 部分，后面俩字节是 47 45，i.e. GE 两个字符。所以这个抓包命令就是在抓 GET 请求的数据包。
>
> 这就是 WireShark 的 Packet Bytes 窗口 中信息的一个典型应用，分析流量中的 Bytes

# WireShark 安装

WireShark 依赖 [pcap](/docs/7.信息安全/Packet%20analyzer/pcap.md)，若使用 WireShark 便携包，那么需要手动安装 [Npcap](https://npcap.com/)。

# Protocol Streams(协议流)

https://www.wireshark.org/docs/wsug_html_chunked/ChAdvFollowStreamSection.html#ChAdvFollowStream

WireShark 将具有某些相同特征的多个数据包分为一组，称为 **Protocol Stream(协议流)**，每条 stream 都有一个 Stream ID。若是 TCP 的包，则是 TCP Protocol Streams；若是 HTTP 的包，则是 HTTP Protocol Stream；以此类推。

WireShark 通过 **Follow Protocol Streams(追踪协议流)** 的行为实现 统计、etc. 很多功能

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
