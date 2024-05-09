---
title: TCP Header
linkTitle: TCP Header
date: 2024-05-09T15:55
weight: 20
---

# 概述

> 参考：
>
> - [RFC 9293，3.1.Header Format](https://datatracker.ietf.org/doc/html/rfc9293#name-header-format)

TCP 段被封装在 IP 数据报中

![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/tcp/1628820358483-a9e565df-371d-4e47-b0d0-0f1fb6077945.png)

首部长度：一般为 20 字节，选项最多 40 字节，限制 60 字节。下图中的位，即代表 bit，也就是说，首部一共 160 bit，即 20 Byte。

![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/tcp/tcp-segment.jpg)

对照在 WireShark 中展示的内容看，排除 `[]` 中的内容，WireShark 中展示的一个 SYN TCP 段的内容，每一行就是包头中的一个内容

![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/tvcktp/1628819589583-6eb31754-8352-45b6-b4b9-b7a61d26433e.png)

- **Source Port(源端口号)** #
- **Destination Port(目的端口号)** #
    - 每个 TCP 报文段都包含源和目的的端口号，这两个端口号用于寻找发送端与接收端的应用进程。这两个值加上 IP 首部中的源和目的的 IP 地址，**组成 TCP 四元组，用于确定唯一一个 TCP 连接**。
- **Sequence Number(序号，简称 SeqNum)** # TCP 报文段的唯一标识符，该标识符具有先后顺序。如果不为每一个包编号，则没法确认哪个包先来哪个包后来。
    - **SeqNum 用来解决网络包乱序的问题。**
    - **Initial Sequence Number(初始序号，简称 ISN)** # TCP 交互的两端，有一个初始的 SeqNum，就是 A 发送给 B 或者 B 发送给 A 的第一个 TCP 段，这第一个 TCP 段的 SeqNum 就是 ISN。
    - 注意：TCP 为应用层提供全双工服务，这意味着数据能在两个方向上独立进行传输。因此，一个 TCP 连接的两端都会有自己独立的 SeqNum。所以首次建立连接时客户端和服务端都会生成一个 ISN。ISN 是一个随机生成的数。
    - SeqNum 最大值为 232-1，到达最大值后，回到 0 开始。
- **Acknowledgment Number(确认序号，简称 AckNum)** # 下一次期望收到数据中报文段的 SeqNum。发出去的包应该有确认，要不然怎么知道对方有没有收到呢？如果没有收到就应该重新发送，直到送达。
    - **AckNum 用来解决丢包的问题**。
    - AckNum 可以用来确认上次发送的数据大小。
        - 假如 172.19.42.244 向 172.19.42.248 发送了一个 PSH,ACK 的报文段，其中 SeqNum=x1,AckNum=y1，发送了 90Bytes 的数据，
        - 那么 172.19.42.248 向 172.19.42.244 就会发送一个 ACK 的报文段，其中 SeqNum=x2,AckNum=x1+90
- **Header Length(首部长度)** #
- **Flag — TCP 标志** # 用来定义当前 TCP 报文段的类型。
- **Window size value(窗口大小)** # TCP 流量控制。通信双方各声明一个窗口，标识自己当前能够的处理能力，别发送的太快，撑死我，也别发的太慢，饿死我。
- **Chceksum(校验和)** #
- **Urgent pointer(紧急指针)** #
- **Options(选项)** # 告诉对方本次传输的一些限制。比如 MSS、SACK 等等
- **TCP Payload(数据)** # 这个字段需要在 HTTP 包中才可以看到，TCP 的有效载荷就是当前传输的数据

## SeqNum 与 AckNum 的计算

由于 TCP 的全双工机制，所以 TCP 交互的两端都有独立的 SeqNum，其两端的 SeqNum 互相之间没有绝对的关联关系，只有一端的 SeqNum 与对端的 AckNum 有关系。

现在假设有 A 和 B 两个系统想要建立 TCP 连接并进行数据交互,并且通信环境正常可以正常建立连接，且以数字表示 A 发送给 B 或者 B 发送给 A 的第几号数据包。

那么：

- A 发送给 B 的第 N 号 TCP 段中 SeqNum 的值等于 `A 的 N-1 号 TCP 段中发送的数据大小(Bytes)` 与 `A 的 N-1 号 TCP 段的 SeqNum` 之和。即：
    - `A_N_SeqNum = A_ISN + A_N-1_TCPPlayload + A_N-1_SeqNum`
- B 响应给 A 的第 N 号 TCP 段中 AckNum 的值等于` A 的 N 号 TCP 段中发送的数据大小(Bytes)` 与 ` A 的 N 号 TCP 段中的 SeqNum` 之和。即：
    - `B_N_AckNum = A_N_TCPPlayload + A_N_SeqNum`

反之亦然：

- `A_N_AckNum = B_N_TCPPlayload + B_N_SeqNum`
- `B_N_SeqNum = B_ISN + B_N-1_TCPPlayload + B_N-1_SeqNum`

## TCP Flag

TCP 报文段的标志内容，将会包含所有标志通过设置标志的值来启用或禁用这些标志(1 表示设置(即.启用)，0 表示未设置(即.禁用))。一个 TCP 报文段中，可以同时启用多个 TCP 标志。下图就是一个 TCP 三次握手中，第二次交互的 TCP 标志内容：

![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/tvcktp/1628827678903-275e8388-c654-4656-b321-d5d501ff2803.png)

当前可用的 TCP 标志有如下几个：

- **ACK(Acknowledgement)** # 确认、响应
    - 除了第一个 SYN 报文段意外，其余报文段都要启用 ACK。因为除了建立连接时，发送的第一个报文段，其余所有的报文段都需要响应发送给自己的报文段，用来表示已收到消息。
- **CWR** # 拥塞窗口减少
- **ECE** # 显式拥塞提醒回应
- **FIN** # 终止、关闭连接
    - 用于释放连接，表示此报文段的发送方数据已发送完毕并要求释放连接
- **PSH(Push**) # 推送、数据传输
    - 当应用程序双方进行交互式通信时，若一端希望在键入命令后就能收到对方响应。此时可采用推送操作。发送方将会立即创建一个报文段并发送出去，接收方接收到 PSH=1 的报文段会进快递交付，而不会等到整个缓存都填满后在向上交付。
- **RST** # 复位、连接重置。
    - 表示连接中出现严重错误必须释放连接再重新建立传输连接，也可用来拒绝一个非法报文段或拒绝打开一个连接。
- **SYN(Synchronize)** # 同步、建立连接
    - SYN 是 TCP/IP 建立连接时使用的握手信号，SYN 仅在三次握手建立 TCP 连接时有效。
    - 客户端和服务端建立连接时，客户端首先会发出一个 SYN 报文段用来建立连接，服务端使用 SYN+ACK 应答表示接收到该消息，最后客户端再以 ACK 消息进行响应。
    - SYN 用于请求和建立连接，也可用于设备间的 SEQ 序列号同步，SYN=1 表示是一个连接请求或连接接收的报文段，当 SYN=1 且 ACK=0 表示连接请求报文段，若对方同意则响应 SYN=1 且 ACK=1。
- **URG** # 紧急
    - 表示报文段中有紧急数据应尽快发送，不要按原来的排队顺序来传送。
