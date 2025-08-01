---
title: 通信协议
linkTitle: Communication protocol
weight: 1
---

# 概述

> 参考：
>
> - [Wiki, Communication Protocol](https://en.wikipedia.org/wiki/Communication_protocol)(通信协议)
> - [Wiki, Encapsulation](<https://en.wikipedia.org/wiki/Encapsulation_(networking)>)(封装)
> - [Wiki, PDU](https://en.wikipedia.org/wiki/Protocol_data_unit)
>   - 注：Wiki 中将网络层 PDU 描述为 Packet 不够准确，详见 RFC 1594 中 13 节 Packet 的名词解释
> - [Wiki, SDU](https://en.wikipedia.org/wiki/Service_data_unit)
> - [RFC 1325](https://datatracker.ietf.org/doc/html/rfc1325)、[RFC 1594](https://datatracker.ietf.org/doc/html/rfc1594)、[RFC 2664](https://datatracker.ietf.org/doc/html/rfc2664)
>   - 这几个 RFC 是一些关于互联网的仅供参考的常见问答，里面包含一些名词解释，2664 是最新版

**Communication Protocol(通信协议)** 是一个规则系统，允许通信系统的两个或多个实体通过[物理量](https://en.wikipedia.org/wiki/Physical_quantity)的任何变化来传输信息。该协议定义了通信的规则、语法、语义和同步以及可能的错误恢复方法。协议可以通过硬件、软件或两者的组合来实现。

# 网络系统的分层架构

因特网是一个极为复杂的系统，这个系统有许多部分：大量的应用程序和协议、各种类型的端系统、分组交换机和各种类型的链路级媒体 等等等等。面对这种巨大的复杂性，我们迫切得需要组织整个网络体系结构。

网络设计者以 **Layer(分层)** 的方式组织协议以及实现这些协议的硬件/软件。每个协议属于这些层次之一，各层的所有协议被统称为 **Protocol Stack(协议栈)**。

![image.png|800](https://notes-learning.oss-cn-beijing.aliyuncs.com/data_comm/protocol/1628914014057-a14c5659-364a-4bfb-ad46-0dbec7375997.png)

# Encapsulation(封装)

**Encapsulation(封装)** 是一种设计模块化通信协议的方法，是将上层数据经过处理，变为下层数据的过程。处理完成后的实体称为 **Protocol Data Unit(协议数据单元，简称 PDU)** 或 **Service Data Unit(服务数据单元，简称 SDU)**。

> PDU 与 SDU 通常来说可以一起理解，楞要说区别，可有有以下几点
>
> - SDU 并不会跨主机，可以说，比如说同一个主机应用层往传输层发送的内容。而 PDU 则是两个不同主机，由 A 主机应用层发送到 B 主机应用层的内容。
> - SDU 是 PDU 的一个子集
> - SDU 是 PDU 中的一个 Payload(有效载荷)

比如，主机 A 要向主机 B 发送一条信息。这条信息就称为 **Data(数据)/Payload(有效载荷)**。这条消息从主机 A 发送出去之前，会被各种协议进行处理，这个处理的过程，就是 Encapsulation(封装)，封装之后的产物就是 PDU，不同网络层的 PDU 叫法不同：

- **Message(报文)** # 应用层的协议封装的 PDU
  - 比如 HTTP Message(HTTP 报文)
- **Segment(段)** # 传输层的协议封装的 PDU
  - 比如 TCP Segment(TCP 段)。为什么称为 Segment(段) 呢，可以想象，从应用层接收到的数据通常会大于 MSS(1460 Bytes)，只要大于 MSS 的数据，都会被分为一段一段的，逐一发送。
- **Datagram(数据报)** # 网络层的协议封装的 PDU
  - 比如 IP Datagram(IP 数据报)。当 IP Datagram 的大小超过 MTU 时，将会被 **Fragment(分片(动词))**，被拆分的每一个部分称为 **Fragment(片(名词))**
    - 每一个 Framment 在很多场景和日常交流中，也被称为 **Packet(包)**，所以很多文档也将网络层协议封装的 PDU 称为 Packet
- **Frame(帧)** # 链路层的协议封装的 PDU
- **bit(比特)** # 物理层的协议封装的 PDU

随着发展，名词 **Packet(包)** 具有更广泛的意义，甚至，可以把 Packet 当做 PDU 来看。这个术语使用松散，虽然某些互联网文献使用它来专门指，通过物理网络发送的数据，但很多其他文献将互联网视为分组交换网络，并将 IP Datagram 描述为 Packet(这一段描述可以在 [RFC1325-术语](https://datatracker.ietf.org/doc/html/rfc1325#page-30)中找到)。

其实，Message、Segment，Datagram，Packet，Frame 是存在于同条记录中的，是基于所在协议层不同而取了不同的名字。我们可以用一个形象的例子对数据包的概念加以说明：我们在邮局邮寄产品时，虽然产品本身带有自己的包装盒，但是在邮寄的时候只用产品原包装盒来包装显然是不行的。必须把内装产品的包装盒放到一个邮局指定的专用纸箱里，这样才能够邮寄。这里，产品包装盒相当于数据包，里面放着的产品相当于可用的数据，而专用纸箱就相当于帧，且一个帧中通常只有一个数据包。

## 封装内的字段结构

一个完整的 PDU 由两个字段组成

- **Header(首部)** # 用来定义数据要传输的目标
  - 不同协议的首部信息各不相同。
- **Payload(有效载荷)** # 需要传输的数据
  - 可以表示为最原本的数据，从下层协议来看，也可以把上层协议封装而成的 SDU 当做 Payload。

通过抓包工具，抓出来的包，在 Wireshark 上查看，就可以看到各种协议封装后，PDU 中的信息。这里以一个 HTTP 的响应包为例

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/data_comm/protocol/1629077712291-6d02d74d-54f4-494b-9b6f-1a850d99005e.png)

可以看到每一层的封装信息，我们现在展开传输层，可以看到如下内容

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/data_comm/protocol/1629077829616-4487b427-08b0-4799-8d16-37e6e3717286.png)

如果展开网络层，则可以看到如下内容：

![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/data_comm/protocol/1629078088587-4b67122f-443c-4e5e-a08b-41f438267ffb.png)

这里就是 IP Datagram 首部的所有信息

在 各个协议的详解中，可以看到各种首部的所有信息，与抓包显示出来的都是可以一一对上的。

# 协议

## 纯二层协议

**Point-to-Point Protocol(点对点协议，简称 PPP)** https://en.wikipedia.org/wiki/Point-to-Point_Protocol

[ARP 与 NDP](/docs/4.数据通信/Protocol/Data%20Link%20Layer/ARP%20与%20NDP.md)

## 纯三层协议

[IP](/docs/4.数据通信/Protocol/TCP_IP/IP/IP.md)

## 纯四层协议

[TCP](/docs/4.数据通信/Protocol/TCP_IP/TCP/TCP.md)

## 上层协议

[HTTP](/docs/4.数据通信/Protocol/HTTP/HTTP.md)

**Remote Authentication Dial-In User Service(远程用户拨号认证，简称 RADUS)** 是一种 7 层协议，RADUS 是一种分布式的、客户端/服务器结构的信息交互协议，能保护网络不受未授权访问的干扰，常应用在既要求较高安全性、又允许远程用户访问的各种网络环境中。RADIUS 协议为标准协议，基本所有主流设备均支持，在实际网络中应用最多。RADIUS 使用UDP（User Datagram Protocol，用户数据报协议）作为传输协议，具有良好的实时性；同时也支持重传机制和备用服务器机制，具有较好的可靠性；实现较为简单，适用于大用户量时服务器端的多线程结构。

- https://info.support.huawei.com/info-finder/encyclopedia/zh/RADIUS.html
