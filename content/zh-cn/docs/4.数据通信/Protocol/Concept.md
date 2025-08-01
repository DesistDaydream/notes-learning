---
title: 概念
linkTitle: Concept
weight: 2
---

# MTU 与 MSS

> 参考：
>
> - [Wiki, MTU](https://en.wikipedia.org/wiki/Maximum_transmission_unit)
> - [Wiki, MSS](https://en.wikipedia.org/wiki/Maximum_segment_size)

**Maximum Transmission Unit(即最大传输单元，简称 MTU)** 是一个二层的概念；以太网最大的 mtu 就是 1500（它是不包含二层头部的，加上头部应该为 1518 bytes，2bit 的以太网类型+6bit 的 DMAC+6bit 的 SMAC+4bit 的 FCS），每个以太网帧都有最小的大小 64bytes，最大不能超过 1518bytes

注：

1. 小于 64Bytes 的数据帧一般是由于以太网冲突产生的 “碎片”或者线路干扰或者坏的以太网接口产生的，对于大于 1518Bytes 的数据帧我们一般把它叫做 Giant 帧，这种一般是由于线路干扰或者坏的以太网口产生
2. 以太网 EthernetII 最大的数据帧是 1518Bytes，是指包含以太网帧的帧头（DMAC 目的 MAC 地址 48bit=6Bytes+SMAC 源 MAC 地址 48bit=6Bytes+Type 域 2bytes）14Bytes 和帧尾 CRC 校验部分 4Bytes （这个部份有时候大家也把它叫做 FCS）

IP MTU 是一个三层概念，它包含了三层头部及所有载荷，根据下层为上层服务的，上层基于下层才能做进一步的扩展的原则，尽管 IP MTU 的变化范围很大（68-65535），但也不得不照顾以太网 MTU 的限制,说白了就是 ip 对以太网的妥协。

网络层 IP 协议会检查每个从上层协议下来的数据包的大小，并根据本机 MTU 的大小决定是否作“分片”处理

**Maximum Segment Size(最大段长度，简称 MSS)** 是 TCP 里面的一个概念，它是 TCP 数据包每次能够传输的最大数据分段，不包含包头部分，它与 IP MTU 满足如下关系：

IP MTU=MSS+20bytes（IP 包头）+20bytes（TCP 包头）

当然，如果传输的时候还承载有其他协议，还要加些包头在前面。

注：为了达到最佳的传输效能，TCP 协议在建立连接的时候通常要协商双方的 MSS 值，这个值 TCP 协议在实现的时候往往用 MTU 值代替（需要减去 IP 数据包报头的大小 20Bytes 和 TCP 数据段的包头 20Bytes），所以往往 MSS 为 1460。通讯双方会根据双方提供的 MSS 值得最小值确定为这次连接的最大 MSS 值。

简言之，mtu 就是总的最后发出去的报文大小，MSS 就是需要发出去的数据大小，比如 PPPoE，就是在以太网上承载 PPP 协议（点到点连接协议），它包括 6bytes 的 PPPoE 头部和 2bytes 的 PPP 协议 ID 号，此时，由于以太网的 MTU 值为 1500，所以上层 PPP 负载数据不能超过 1492 字节，也就是相当于在 PPPOE 环境下的 MTU 是 1492 字节，MSS 是 1452 字节（1492 字节-20-20）。

重点：

MTU 不包含 帧头（18byte） 指帧头后面的所有负载，与 ip mtu 的区别就是在帧头和 ip 头之间可能会有其他协议头（比如 GRE 头、pppoe 头、MPLS 标签，这些协议头都是在帧头后 ip 头前）

ip MTU 包含 ip 头（20byte） 指 ip 头本身及后面的所有负载，一个普通的以太网数据包 mtu=ip mut，只有封装了其他协议头部时 mtu=ip mut+其他协议头部+负载（tcp 头+tcp-mss）

TCP-MSS 不包含 tcp 头（20byte） 指 tcp 头后面的所有负载

IP MTU=tcp-MSS+20bytes（IP 包头）+20bytes（TCP 包头）

<https://serverfault.com/questions/500448/mysterious-fragmentation-required-rejections-from-gateway-vm>

# MSL、RTT、TTL

原文链接：<https://my.oschina.net/vbird/blog/1525869>

在[《TCP 关闭状态分析》](https://my.oschina.net/vbird/blog/1507479)一文中，段落“TIME_WAIT”中提到：

之后等待 2 个最大的报文存活周期（这是因为：一是保证残留网络报文不会被新连接接收而产生数据错乱，由于自己上一次发送的数据报文可能还残留在网络中，等待 **2MSL** 时间可以保证所有残留的网络报文在自己关闭前都已经超时。二是确保自己最后 ACK 发送到对端，因为 ACK 发送也可能会发送失败，这时对端会重新发送 FIN，如果已经 CLOSED 了那么对端将收到 RST 而不是 ACK 了，这不符合 TCP 可靠关闭的策略。）

## MSL

Maximum Segment Lifetime 缩写，译为“报文最大生存时间”。它指的是任何报文在网络上存在的最长时间，超过这个时间的报文将会被丢弃。标准规范中规定 MSL 为 2 分钟，实际应用中常用的是 30s，1min 和 2min 等。
2MSL 即两倍的 MSL，TCP 的 TIME_WAIT 状态也称为**2MSL 等待状态**，当 TCP 的一端发起主动关闭，在发出最后一个 ACK 包后，即第 3 次挥手完成后发送了第四次挥手的 ACK 包后就进入了 TIME_WAIT 状态，必须在此状态上停留两倍的 MSL 时间，等待 2MSL 时间主要目的是怕最后一个 ACK 包对方没收到，那么对方在超时后将重发第三次挥手的 FIN 包，主动关闭端接到重发的 FIN 包后可以再发一个 ACK 应答包。在 TIME_WAIT 状态时两端的端口不能使用，要等到 2MSL 时间结束才可继续使用。当连接处于 2MSL 等待阶段时任何迟到的报文段都将被丢弃。不过在实际应用中可以通过设置 **SO_REUSEADDR**选项达到不必等待 2MSL 时间结束再使用此端口。

在 [Linux Kernel 5.12 版本的源码码](https://github.com/torvalds/linux/blob/v5.12/include/net/tcp.h#L121) 可以看到默认值是 60

## RTT

> To this end, we define the **round-trip time (RTT)**, which is the time it takes for a small packet to travel from client to server and then back to the client.

**Round-Trip Time(往返时间，简称 RTT)**，简单理解的话，RTT 指的是客户端到服务端往返所花费的时间。其实 RTT 的定义是一个很小的分组(这里的很小也就是说对于发送方来说“开始发送”和“发送完”是同一个时刻。换句话说这个分组是没有长度的，传输时延可以忽略不计)，从客户端发送到接收端再返回客户端的时间。TCP 三次握手中发送的 SYN、ACK 以及四次挥手中使用到的 FIN 的交换就是这样的例子，ping 命令使用的也是这么一个非常小的包来试探网络的情况。如果是长度不可忽略的分组，接收方在接收到最后一比特数据后才会发送 ACK，实际也等价于发送最后一 bit 到接收到 ACK 的时间间隔。
TCP 含有动态估算 RTT 的算法。TCP 还持续估算一个给定连接的 RTT，这是因为 RTT 受网络传输拥塞程序的变化而变化。

## TTL

详见 [IP Header](/docs/4.数据通信/Protocol/TCP_IP/IP/IP%20Header.md) 的 TTL 部分。

# DMZ 与 UPnP

DMZ 主机的作用：设置 DMZ 主机，就是让一台内网电脑完全暴露在外网，也就是任何由外部发起的连接到路由的数据都会被转发到设定的内网主机。这样做的坏处就是路由无法提供自身的防火墙来保护内网的这台机器，存在一定的风险，一般情况下不推荐。

# Interpacket gap

> 参考：
>
> - [Wiki, Interpacket gap](https://en.wikipedia.org/wiki/Interpacket_gap)

**Interpacket gap(包间隙，简称 IPG)** 也称为 **interframe gap(IFG)** 是网络数据包或网络帧之间可能需要的暂停。根据所使用的物理层协议或编码，暂停可能是必要的，以允许接收器时钟恢复，允许接收器准备另一个数据包（例如，从低功耗状态上电）或其他目的。可以将其视为保护间隔的特定情况。

# PDU

> 参考：
>
> - [Wiki, PDU](https://en.wikipedia.org/wiki/Protocol_data_unit)
