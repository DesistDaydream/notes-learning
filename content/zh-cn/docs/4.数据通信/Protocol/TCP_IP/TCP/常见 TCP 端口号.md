---
title: 常见 TCP 端口号
---

# Socket 网络端口含义

> 参考：
>
> - [Wiki, TCP 与 UDP 端口号列表](https://en.wikipedia.org/wiki/List_of_TCP_and_UDP_port_numbers)
> - [IANA，分配-服务名和端口号注册列表](https://www.iana.org/assignments/service-names-port-numbers/service-names-port-numbers.xhtml)

分类情况：

- WellKnownPorts(公认端口)
  - 从 0 到 1023，它们紧密绑定（binding）于一些服务。通常这些端口的通讯明确表明了某种服务的协议。这些端口由 IANA 分配管理，IANA(The Internet Assigned Numbers Authority，互联网数字分配机构)是负责协调一些使 Internet 正常运作的机构
- RegisteredPorts(注册端口)
  - 从 1024 到 49151。是公司和其他用户向互联网名称与数字地址分配机构（ICANN）登记的端口号，利用因特网的传输控制协议（TCP）和用户数据报协议（UDP）进行通信的应用软件需要使用这些端口。在大多数情况下，这些应用软件和普通程序一样可以被非特权用户打开。
- Dynamicand/orPrivatePorts(动态和/或私有端口)
  - 从 49152 到 65535。这类端口号仅在客户进程运行时才动态选择，因此又叫做短暂端口号。被保留给客户端进程选择暂时使用的。也可以理解为，客户端启动的时候操作系统随机分配一个端口用来和服务器通信，客户端进程关闭下次打开时，又重新分配一个新的端口。

# TCP/UDP 端口列表

计算机之间依照互联网[传输层](https://zh.wikipedia.org/wiki/%E4%BC%A0%E8%BE%93%E5%B1%82)[TCP/IP 协议](https://zh.wikipedia.org/wiki/TCP/IP%E5%8D%8F%E8%AE%AE)的协议通信，不同的协议都对应不同的[端口](https://zh.wikipedia.org/wiki/%E9%80%9A%E8%A8%8A%E5%9F%A0)。并且，利用数据报文的[UDP](https://zh.wikipedia.org/wiki/%E7%94%A8%E6%88%B7%E6%95%B0%E6%8D%AE%E6%8A%A5%E5%8D%8F%E8%AE%AE)也不一定和[TCP](https://zh.wikipedia.org/wiki/%E4%BC%A0%E8%BE%93%E6%8E%A7%E5%88%B6%E5%8D%8F%E8%AE%AE)采用相同的端口号码。以下为两种通信协议的端口列表链接

## 众所周知的端口

0 到 1023 号端口

0 到 1023（0 到 2^10 − 1）范围内的端口号是众所周知的端口或系统端口。 它们由提供广泛使用的网络服务类型的系统进程使用。在类 Unix 操作系统上，进程必须以超级用户权限执行才能使用众所周知的端口之一将网络套接字绑定到 IP 地址。

以下列表仅列出常用端口，详细的列表请参阅 [IANA](https://zh.wikipedia.org/wiki/IANA) 网站。

## 注册端口

https://en.wikipedia.org/wiki/List_of_TCP_and_UDP_port_numbers#Registered_ports

1024 - 49151 号端口

从1024 到 49151 (210到214 215-1) 的端口号范围是注册的端口。IANA在请求实体申请时将它们分配给特定服务。[2] 在大多数系统上，注册端口可以在没有超级用户特权的情况下使用。

## 动态、私有或临时端口

https://en.wikipedia.org/wiki/List_of_TCP_and_UDP_port_numbers#Dynamic,_private_or_ephemeral_ports

参见：[临时端口](https://zh.wikipedia.org/wiki/%E4%B8%B4%E6%97%B6%E7%AB%AF%E5%8F%A3)

49152 - 65535 号端口

根据定义，该段端口属于“动态端口”范围，没有端口可以被正式地注册占用。
