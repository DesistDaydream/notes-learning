---
title: IP
linkTitle: IP
weight: 1
---

# 概述

> 参考：
>
> - [RFC 791， INTERNET PROTOCOL PROTOCOL SPECIFICATION](https://datatracker.ietf.org/doc/html/rfc791)
> - [Wiki, Internet Protocol](https://en.wikipedia.org/wiki/Internet_Protocol)
> - [Wiki, IPv4](https://en.wikipedia.org/wiki/IPv4)
> - [Wiki, Mask(掩码)](<https://en.wikipedia.org/wiki/Mask_(computing)>)
> - [Wiki, Classful Network(分类网络)](https://en.wikipedia.org/wiki/Classful_network)
> - [IANA,IPv4 地址空间分配情况](https://www.iana.org/assignments/ipv4-address-space/ipv4-address-space.xhtml)
>   - [APNIC](https://www.apnic.net/)(管理亚太地区的 IP 地址注册机构)
>     - [APNIC,帮助-FTP 数据库](https://ftp.apnic.net/stats/apnic/)(亚太地区所有分配的 IP 地址信息)
> - [IANA,IPv4 特殊用途地址注册表](https://www.iana.org/assignments/iana-ipv4-special-registry/iana-ipv4-special-registry.xhtml)

**Internet Protocol(互联网协议，简称 IP)** 是[互联网协议套件](https://en.wikipedia.org/wiki/Internet_protocol_suite)(其中包含 TCP/IP)中的主要通信协议，用于跨网络边界中继数据报。它的路由功能可实现互联网络，并实质上建立了 Internet。

> **Internet protocol suite(互联网协议套件)** 是互联网和类似计算机网络中使用的概念模型和通信协议集。由于该套件中的基本协议是 **TCP(传输控制协议)** 和 **IP(互联网协议)**，因此通常被称为 **TCP/IP**。在其开发过程中，其版本被称为国防部（DoD）模型，因为联网方法的开发是由美国国防部通过 DARPA 资助的。它的实现是一个协议栈。

IP 基于数据包的 Header 中的 IP 地址，将数据包从源主机发送到目标主机。基于此目的，IP 还定义了数据包的封装结构、以及一种寻址方法。寻址方法用来使用源和目标的信息标记数据报。

从历史上看，IP 是在 1974 年由 Vint Cerf 和 Bob Kahn 引入的原始 **Transmission Control Program(传输控制程序)** 中的[无连接](https://en.wikipedia.org/wiki/Connectionless_communication)数据报服务。该服务由一项面向连接的服务补充，成为 [**Transmission Control Protocol(传输控制协议，简称 TCP)**](/docs/4.数据通信/Protocol/TCP_IP/TCP/TCP.md) 的基础。因此 IP 套件通常称为 TCP/IP。IP 的第一个版本是 IPv4，继任者是 IPv6

# IPv4 地址

IPv4 地址最多使用 32 bit 表示，即最多 32 个 1，这 32 bit 以 `点` 分割为 4 组，每组 8 bit，在使用时，使用十进制表示。比如：`192.168.0.1`。

![800](https://notes-learning.oss-cn-beijing.aliyuncs.com/ip/ip_address_2-to-10.png)

## IPv4 地址结构

IPv4 地址的这 32 bit 可以分为两部分

- 网络位 # n bit
- 主机位 # 32 - n bit

这两个部分通过 **Subnet Mast(子网掩码)** 来区分，子网掩码由一连串的 1 和 0 组成，遵从以下规则：

- 1 对应网络位
- 0 对应主机位
- 1 和 0 不能交叉出现

将子网掩码和 IP 地址作“与”操作后，IP 地址的主机部分将被丢弃，剩余的是网络地址和子网地址。

例如：一个 IP 地址为 10.2.45.1，子网掩码为 255.255.252.0，“与” 运算得到：10.2.44.0，则网络设备认为该 IP 地址的网络号与子网号为 10.2.44.0，属于 10.2.44.0/22 网络，其中/22 表示子网掩码长度为 22 位，即从前向后连续的 22 个 1。

```
00001010.00000010.00101101.00000001
与运算
11111111.11111111.11111100.00000000

结果为
00001010.00000010.00101100.00000001 即 10.2.44.0
```

## IPv4 地址分类

- **单播地址**

| 类       | 开头的 bit | 网络位 bit 数 | 主机位 bit 数 | 子网数量            | 每个子网的地址数         | 总地址数                | 起始地址      | 结束地址            | 默认子网掩码        | [CIDR](https://en.wikipedia.org/wiki/CIDR_notation) |
| ------- | ------- | --------- | --------- | --------------- | ---------------- | ------------------- | --------- | --------------- | ------------- | --------------------------------------------------- |
| Class A | 0       | 8         | 24        | 128 (27)        | 16,777,216 (224) | 2,147,483,648 (231) | 0.0.0.0   | 127.255.255.255 | 255.0.0.0     | /8                                                  |
| Class B | 10      | 16        | 16        | 16,384 (214)    | 65,536 (216)     | 1,073,741,824 (230) | 128.0.0.0 | 191.255.255.255 | 255.255.0.0   | /16                                                 |
| Class C | 110     | 24        | 8         | 2,097,152 (221) | 256 (28)         | 536,870,912 (229)   | 192.0.0.0 | 223.255.255.255 | 255.255.255.0 | /24                                                 |

- **组播地址**

  - **D 类 224-239 224.0.0.0 ~ 239.255.255.255**

- **保留地址**

  - **E 类 240 - 254 240.0.0.0 ~ 255.255.255.255**

- **特殊地址**

  - **网络地址** # 网络位不变，主机位全为 0 的 IP 地址代表网络本身
  - **Broadcast Address(广播地址)** # 网络位不变，主机位全为 1 的 IP 地址代表本网络的广播。是专门用于同时向网络中所有工作站进行发送的一个**地址**。在使用 TCP/IP 协议的网络中，主机标识]段 host ID 为全 1 的 IP 地址为广播地址，广播的分组传送给 host ID 段所涉及的所有计算机。例如，对于 10.1.1.0 （255.0.0.0 ）网段，其直播广播地址为 10.255.255.255 （255 即为 2 进制的 11111111 ），当发出一个目的地址为 10.255.255.255 的分组（封包）时，它将被分发给该网段上的所有计算机。
  - **Link Local(链路本地地址)** # 169.254.0.0 ~ 169.254.255.255。用于[链路本地地址](https://en.wikipedia.org/wiki/Link-local_address)两台主机之间的单个链路上时，否则指定 IP 地址，如将有通常被从检索到的 [DHCP](/docs/4.数据通信/Protocol/DHCP.md) 服务器。

- **Private Network(私人网络地址)**

| 名称           | CIDR           | 地址范围                          | 地址数量     | 描述                                      |
| ------------ | -------------- | ----------------------------- | -------- | --------------------------------------- |
| 24-bit block | 10.0.0.0/8     | 10.0.0.0 – 10.255.255.255     | 16777216 | 一个完整的 A 类地址 Single Class A.             |
| 20-bit block | 172.16.0.0/12  | 172.16.0.0 – 172.31.255.255   | 1048576  | Contiguous range of 16 Class B blocks.  |
| 16-bit block | 192.168.0.0/16 | 192.168.0.0 – 192.168.255.255 | 65536    | Contiguous range of 256 Class C blocks. |

# IPv4 Datagram 结构

详见 [IP Header](/docs/4.数据通信/Protocol/TCP_IP/IP/IP%20Header.md)

# IPv4 Fragment

IP Fragment(分片) 主要通过首部中的 Identification、Flags、Fragment Offset 这三个字段对每一个分片进行唯一标识
