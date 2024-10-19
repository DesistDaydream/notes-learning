---
title: TCP/IP
linkTitle: TCP/IP
date: 2023-11-16T21:14
weight: 1
---

# 概述

> 参考：
>
> - [RFC，791](https://datatracker.ietf.org/doc/html/rfc791)
> - [Wiki, Internet protocol suite](https://en.wikipedia.org/wiki/Internet_protocol_suite)

**Internet protocol suite(互联网协议簇)** 通常称为 **TCP/IP**

## 封装

当应用程序用 TCP 传输数据时，数据被送入协议栈中，然后逐个通过每一层直到被当做一串比特流送入网络。其中每一层对收到的数据都要增加一些首部信息(有时还需要增加尾部信息)，过程如图所示：

![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/oc8ill/1628821300980-52f384b2-d2c9-4227-a1c5-6481d6cbf20e.png)

以太网帧的帧头和帧尾下面所标注的数字是典型以太网帧首部的字节查高难度。以太网数据帧的物理特性是其查高难度必须在 **46~1500 字节之间(也就是 MTU 的长度)**

> 注意：所有的 Internet 标准和大多数有关 TCP/IP 的书都使用 octe 这个术语来表示字节。使用这个过分雕琢的术语是有历史原因的，因为 TCP/IP 的很多工作都是在 DEC-10 系统上进行的，但是它并不使用 8bit 的字节。由于现在几乎所有的计算机系统都采用 8bit 的字节，因此我们在本书中使用 Byte(字节) 这个术语。

由于应用数据受 MSS 长度限制，IP 首部 + TCP 首部 + 应用数据受 MTU 长度限制。所以，当一个 IP 报文超过 MTU 时就会进行 **Packet(分片/分组)**。分组既可以是一个 IP 数据报，也可以是 IP 数据报的一个 **Fragment(片段)**。

![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/oc8ill/1628821542180-6dae0209-e7ac-494e-b6f3-715ce143c6d5.png)
