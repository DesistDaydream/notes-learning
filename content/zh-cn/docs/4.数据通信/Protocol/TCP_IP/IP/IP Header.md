---
title: IP Header
linkTitle: IP Header
date: 2024-02-29T09:32
weight: 20
---

# 概述

> 参考：
>
> - [RFC 791，3.1.Internet Header Format](https://datatracker.ietf.org/doc/html/rfc791#section-3.1)

IPv4 数据报被封装在链路层的 Frame 中

IPv4 数据报首部共 14 个字段，其中 13 个是必须的，第 14 个是可选的。前 13 个字段长度固定为 20 Bytes，即 160 bit；第 14 个字段长度在 0 ~ 40 Bytes 之间。

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/ip/ip_datagram.png)

对照 [WireShark](/docs/7.信息安全/Packet%20analyzer/WireShark/WireShark.md) 中展示的内容看，排除 `[]` 中的内容，每一行就是首部中的一个字段

![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/ip/ip_header_in_wireshark.png)

- **Version(版本)** # IP 协议的版本号。IPv4 其版本号为 4，因此在这个字段上的值为“6”。
- **Internet Header Length(首部长度，简称 IHL)** # 由于 Options 字段的长度是可变的。所以 IPv4 的首部长度也是可变的。该字段的值在 5 ~ 15 之间(该字段只有 4 bits，1111 即为 15)
  - 首部长度的计算方式如下：`IHL * 32 bits`。
    - 若 IHL 的值为 5，也就是说 Options 字段为 0，那么 IPv4 首部长度就是 5 \* 32 bits = 160 bits = 20 Bytes
  - 就像上面的 IPv4 的 Datagram 结构图一样，每行都是 32 bit，不算 Options 字段和 Payload，那么刚好是 5 行。
- **Differentiated Services Field** # 差异化的服务字段，基本没啥用。。。。o(╯□╰)o
  - **Differentiated Services Code Point** # 最初定义为 Type Of Service(服务类型，简称 TOS)，
  - **Explicit Congestion Notification** # 该字段定义在 RFC3168 中，
- **Total Length** # 定义了整个 IP 数据报的大小，最小为 20 字节(Payload 字段无内容)，最大为 65535 字节。
- **Identification**# 主要用于唯一标识单个 IP 数据报的片段组。
  - 一些实验工作建议将 ID 字段用于其他目的，例如添加数据包跟踪信息以帮助跟踪具有欺骗源地址的数据报，\[31] 但 RFC 6864 现在禁止任何此类使用。
- **Flags** # 用来控制或识别 IP 分片之后的每个片段，这 3 个 bit 分别表示不同的含义，若字段值为 0 表示未设置，值为 1 表示设置，类似 TCP 首部中 Flags 字段的用法。
  - 第一个 # Reserved，保留字段，必须为 0
  - 第二个 # Don't Fragment(DF)
  - 第三个 # More Fragment(MF)
- **Fragment Offset(分片偏移)** # IP 分片之后的偏移量
- **Time To Live(存活时间，简称 TTL)** # 其实用 Hop Limit 的描述更准确，封包每经过一个路由器，就会将 TTL 字段的值减 1，减到 0 是，该包将会被丢弃。
- **Protocol**# 封装 IP 数据报的上层协议，比如 6 表示 TCP、1 表示 ICMP
  - 每种协议根据 [RFC 1700](https://datatracker.ietf.org/doc/html/rfc1700) 都分配了一个固定的编号，该 RFC 1700 最终被 RFC 3232 废弃，并将协议编号的维护工作，转到[IANA 的在线数据库](https://www.iana.org/assignments/protocol-numbers/protocol-numbers.xhtml)中
- **Header Checksum** # 当数据包到达路由器时，路由器会计算标头的校验和，并将其与校验和字段进行比较。如果值不匹配，则路由器会丢弃该数据包。
- **Source Address(源地址)** # 发送端 IP 地址。
- **Destination Address(目标地址)** # 接收端 IP 地址。
- **Options(选项)** # 可变长度，0-40 Bytes。

# TTL

**Time To Live(存活时间，简称 TTL)**。这个生存时间是由源主机设置（IP 头里有一个 TTL 域）初始值但不是存的具体时间，而是存储了一个 IP 数据报可以经过最大路由数，每经过一个处理它的路由器此值就减 1，当此值为 0 则数据报将被丢弃，同时发送 ICMP 报文通知源主机。

**注意**：域名也有 TTL 的概念，只是这里的 TTL 指的是 Time- To-Live，简单的说它表示一条域名解析记录在 DNS 服务器上缓存时间。当各地的 DNS 服务器接受到解析请求时，就会向域名指定的 DNS 服务器发出解析请求从而获得解析记录。在获得这个记录之后，记录会在 DNS 服务器中保存一段时间，这段时间内如果再接到这个域名的解析请求，DNS 服务器将不再向 DNS 服务器发出请求，而是直接返回刚才获得的记录。而这个记录在 DNS 服务器上保留的时间，就是 TTL 值。
