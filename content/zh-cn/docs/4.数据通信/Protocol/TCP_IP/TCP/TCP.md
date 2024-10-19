---
title: TCP
weight: 1
---

# 概述

> 参考：
>
> - [RFC 675,](https://datatracker.ietf.org/doc/html/rfc675)
> - [RFC 793, TRANSMISSION CONTROL PROTOCOL - DARPA INTERNET PROGRAM PROTOCOL SPECIFICATION](https://datatracker.ietf.org/doc/html/rfc793)
>   - [RFC 9293, Transmission Control Protocol (TCP)](https://datatracker.ietf.org/doc/html/rfc9293)
> - [Wiki, TCP](https://en.wikipedia.org/wiki/Transmission_Control_Protocol)
> - [极客时间,趣谈网络协议](https://time.geekbang.org/column/intro/100007101)
> - <https://www.jianshu.com/p/1118f497a425>
> - <https://www.jianshu.com/p/3c7a0771b67e>
> - [公众号-小林coding，通过动图学习 TCP 的滑动窗口和流量控制的工作方式](https://mp.weixin.qq.com/s/WG1Is0HMAHYMgRvJQpd3KA)
>   - <https://www2.tkn.tu-berlin.de/teaching/rn/animations/gbn_sr/>
>   - <https://www2.tkn.tu-berlin.de/teaching/rn/animations/flow/>

**Transmission Control Protocol(传输控制协议，简称 TCP)** 是[互联网协议套件](https://en.wikipedia.org/wiki/Internet_protocol_suite)的最主要协议之一。它起源于最初的网络实现，补充了 Internet Protocol。因此整个套件通常称为 **TCP/IP**。

IP 地址后面的端口的作用：当从外部访问该 IP 地址的机器时候，是通过该 IP 地址的端口来访问这台机器的某个程序，然后程序向访问者提供该程序所具有的功能（服务）。web 界面默认是 80 端口，那么当你访问一个网页的时候，这个 IP 就会带你访问该机器的占用 80 端口的程序，然后该程序去调用首页脚本本间展示给访问者

例如：你通过远程 SSH 访问一台设备 192.168.0.1 的话，那么需要设置一下这台机器 SSH 服务程序所占用的端口号，比如 22，那么你就是通过 192.168.0.1:22 这个来访问这台机器的 SSH 进程。

上面的描述，就是一个基本的 TCP。

TCP 天然认为网络环境是恶劣的，丢包、乱序、重传，拥塞都是常有的事情，一言不合就可能送达不了，因而要从算法层面来保证可靠性。TCP 是靠谱的协议，但是这不能说明它面临的网络环境好。从 IP 层面来讲，如果网络状况的确那么差，是没有任何可靠性保证的，而作为 IP 的上一层 TCP 也无能为力，唯一能做的就是更加努力，不断重传，通过各种算法保证。也就是说，对于 TCP 来讲，IP 层你丢不丢包，我管不着，但是我在我的层面上，会努力保证可靠性。这有点像如果你在北京，和客户约十点见面，那么你应该清楚堵车是常态，你干预不了，也控制不了，你唯一能做的就是早走。打车不行就改乘地铁，尽力不失约。

## 总结

通过对 TCP 头的解析，我们知道要掌握 TCP 协议，重点应该关注以下几个问题：

- 顺序问题 ，稳重不乱；
- 丢包问题，承诺靠谱；
- 连接维护，有始有终；
- 流量控制，把握分寸；
- 拥塞控制，知进知退。

# TCP Segment 结构

详见 [TCP Header](/docs/4.数据通信/Protocol/TCP_IP/TCP/TCP%20Header.md)

# TCP 状态

一个 TCP 连接在它的声明周期内会有不同的状态。下图说明了 TCP 连接可能会有的状态，以及基于事件的状态转换。事件中有的是应用程序的操作，有的是接收到了网络发过来的请求。

![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/tvcktp/1628824939321-772e1691-86d4-4ca9-b65c-7a3db704ff9b.png)

- **LISTEN** # 服务端。等待来自远程的 TCP 请求
- **SYN-SEN**T # TCP 第一次握手后客户端所处的状态。发送连接请求后，等待来自服务端的确认。
  - TCP 默认 SYN 报文最大 retry 5 次，每次超时了翻倍，1s -> 3s -> 7s -> 15s -> 31s -> 63s。[参考资料](https://blog.csdn.net/u010039418/article/details/78234570)
- **SYN-RECEIVED** # TCP 第二次握手后服务端所处的状态。服务端已经接收到连接请求并发送确认。服务端正在等待最终确认。
- **ESTABLISHED** # TCP 第三次握手后服务端与客户端所处的状态。代表连接已经建立起来了。这是连接数据传输阶段的正常状态。
- **FIN-WAIT-1** # 等待来自远程 TCP 的终止连接请求或终止请求的确认
- **FIN-WAIT-2** # 在此端点发送终止连接请求后，等待来自远程 TCP 的连接终止请求
- **CLOSE_WAIT** # 该端点已经收到来自远程端点的关闭请求，此 TCP 正在等待本地应用程序的连接终止请求
- **CLOSING** # 等待来自远程 TCP 的连接终止请求确认
- **LAST_ACK** # 等待先前发送到远程 TCP 的连接终止请求的确认
- **TIME-WAIT** # 主动断开连接一方的状态。等待足够的时间来确保远程 TCP 接收到其连接终止请求的确认
  - TCP 主动关闭连接的一方在发送最后一个 ACK 后进入  `TIME_AWAIT`  状态，再等待 2 个 MSL 时间后才会关闭(因为如果 server 没收到 client 第四次挥手确认报文，server 会重发第三次挥手 FIN 报文，所以 client 需要停留 2 MSL 的时长来处理可能会重复收到的报文段；同时等待 2 MSL 也可以让由于网络不通畅产生的滞留报文失效，避免新建立的连接收到之前旧连接的报文)，了解更详细的过程请参考 TCP 四次挥手。

# TCP 行为的过程

![image.png|500](https://notes-learning.oss-cn-beijing.aliyuncs.com/tcp/tcp-three-way-handshake-and-four-way-wave.png)

## TCP 三次握手

参考：<https://hit-alibaba.github.io/interview/basic/network/TCP.html>

三次握手的过程的示意图如下：

![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/tvcktp/1628831488474-13bf079c-3419-4f9d-b92d-ca1620d7b6a7.png)

所谓三次握手(Three-way Handshake)，是指建立一个 TCP 连接时，需要客户端和服务器总共发送 3 个包。三次握手的目的是连接服务器指定端口，建立 TCP 连接，并同步连接双方的序列号和确认号，交换 TCP 窗口大小信息。

在 socket 编程中，客户端执行 connect() 时。将触发三次握手：

- **第一次握手(SYN，SeqNum=client_isn)** # 客户端请求同步，发送 SYN 报文段
  - 客户端生成一个随机数 client_isn
  - 设置 TCP 首部字段
    - 将 client_isn 填入到 Sequence Number 字段中。
    - 将 TCP 标志中 SYN 的值设为 1。
  - 发送完毕后，客户端进入 SYN_SEND 状态。
- **第二次握手(SYN,ACK，SeqNum=server_isn，AckNum=client_isn+1)** # 服务端回应并请求同步，发送 SYN + ACK 报文段
  - 服务端收到客户端的 SYN 报文段后，生成一个随机数 server_isn
  - 设置 TCP 首部字段
    - 将 server_isn 填入 Sequence Number 字段中
    - 将 client_isn+1 填入 Acknowledgement Number 字段中。
    - 将 TCP 标志中的 SYN 和 ACK 的值设为 1。
  - 发送完毕后，服务器端进入 SYN_RCVD 状态。
- **第三次握手(ACK，SeqNum=client_isn+1，AckNum=server_isn+1)** # 客户端回应确认，建立连接。发送 ACK 报文段
  - 客户端收到服务端报文后，还需要向服务端回应最后一个 ACK 报文段
  - 设置 TCP 首部字段
    - 将 server_isn+1 填入 Acknowledgement Number 字段中。
    - 将 TCP 标志中的 ACK 的值设为 1。
  - 发送完毕后，客户端进入 ESTABLISHED 状态，当服务器端接收到这个包时，也进入 ESTABLISHED 状态，TCP 握手结束。

注意：

- 第三次握手时，此时客户端已经处于 ESTABLISHED 状态。对于客户端来说，他已经建立起连接了，并且已经知道服务器的接收和发送能力是正常的。所以也就可以携带数据了。
- 但是第一次握手是不可以带数据发送给服务端的，因为还不知道对方的窗口大小，并且也容易让客户端发起大量数据进行攻击。

## TCP 四次挥手

四次挥手的示意图如下：

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/tvcktp/1616161421390-7655ce0e-3e07-436c-93b1-708b3389996d.jpeg)

TCP 的连接的拆除需要发送四个包，因此称为四次挥手(Four-way handshake)，也叫做改进的三次握手。客户端或服务器均可主动发起挥手动作，在 socket 编程中，任何一方执行 close() 操作即可产生挥手操作。

第一次挥手(FIN=1，seq=x)

- 假设客户端想要关闭连接，客户端发送一个 FIN 标志位置为 1 的包，表示自己已经没有数据可以发送了，但是仍然可以接受数据。
- 发送完毕后，客户端进入 FIN_WAIT_1 状态。

第二次挥手(ACK=1，ACKnum=x+1)

- 服务器端确认客户端的 FIN 包，发送一个确认包，表明自己接受到了客户端关闭连接的请求，但还没有准备好关闭连接。
- 发送完毕后，服务器端进入 CLOSE_WAIT 状态，客户端接收到这个确认包之后，进入 FIN_WAIT_2 状态，等待服务器端关闭连接。

第三次挥手(FIN=1，seq=y)

- 服务器端准备好关闭连接时，向客户端发送结束连接请求，FIN 置为 1。
- 发送完毕后，服务器端进入 LAST_ACK 状态，等待来自客户端的最后一个 ACK。

第四次挥手(ACK=1，ACKnum=y+1)

- 客户端接收到来自服务器端的关闭请求，发送一个确认包，并进入 TIME_WAIT 状态，等待可能出现的要求重传的 ACK 包。
- 服务器端接收到这个确认包之后，关闭连接，进入 CLOSED 状态。
- 客户端等待了某个固定时间（两个最大段生命周期，2MSL，2 Maximum Segment Lifetime）之后，没有收到服务器端的 ACK ，认为服务器端已经正常关闭连接，于是自己也关闭连接，进入 CLOSED 状态。

## 数据传输

![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/tvcktp/1628842098522-1dc9dcc9-026c-4e88-a288-cd0e73613e77.png)

![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/tvcktp/1628841642687-befcb683-4d7e-41b2-849d-14b028d2a170.png)

## 通过 WireShark 理解三次握手和四次挥手

[TCP Analysis](/docs/7.信息安全/Packet%20analyzer/WireShark/TCP%20Analysis.md)

# TCP 报文

> 参考：
>
> - 原文：[程序员宅基地，TCP报文（ tcp dup ack 、TCP Retransmission）](https://www.cxyzjd.com/article/ynchyong/109110028)
>   - [CSDN，TCP报文（ tcp dup ack 、TCP Retransmission）](https://blog.csdn.net/ynchyong/article/details/109110028)

## TCP 慢启动

慢启动是 TCP 的一个拥塞控制机制，慢启动算法的基本思想是当 TCP 开始在一个网络中传输数据或发现数据丢失并开始重发时，首先慢慢的对网路实际容量进行试探，避免由于发送了过量的数据而导致阻塞。

慢启动为发送方的 TCP 增加了另一个窗口：拥塞窗口(congestion window)，记为 cwnd。当与另一个网络的主机建立 TCP 连接时，拥塞窗口被初始化为 1 个报文段（即另一端通告的报文段大小）。每收到一个 ACK，拥塞窗口就增加一个报文段（cwnd 以字节为单位，但是慢启动以报文段大小为单位进行增加）。发送方取拥塞窗口与通告窗口中的最小值作为发送上限。拥塞窗口是发送方使用的流量控制，而通告窗口则是接收方使用的流量控制。发送方开始时发送一个报文段，然后等待 ACK。当收到该 ACK 时，拥塞窗口从 1 增加为 2，即可以发送两个报文段。当收到这两个报文段的 A C K 时，拥塞窗口就增加为 4。这是一种指数增加的关系。

## 拥塞避免算法

![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/wireshark/1659010179458-470eeff4-bf7b-4d4b-a5d2-215b501d9f61.png)

网络中拥塞的发生会导致数据分组丢失，需要尽量避免。在实际中，拥塞算法与慢启动通常在一起实现，其基本过程：

1. 对一个给定的连接，初始化 cwnd 为 1 个报文段，ssthresh 为 65535 个字节。
2. TCP 输出例程的输出不能超过 cwnd 和接收方通告窗口的大小。拥塞避免是发送方使用 的流量控制，而通告窗口则是接收方进行的流量控制。前者是发送方感受到的网络拥塞的估 计，而后者则与接收方在该连接上的可用缓存大小有关。
3. 当拥塞发生时（超时或收到重复确认），ssthresh 被设置为当前窗口大小的一半（cwnd 和接收方通告窗口大小的最小值，但最少为 2 个报文段）。此外，如果是超时引起了拥塞，则 cwnd 被设置为 1 个报文段（这就是慢启动）。
4. 当新的数据被对方确认时，就增加 cwnd，但增加的方法依赖于是否正在进行慢启动或拥塞避免。如果 cwnd 小于或等于 ssthresh，则正 在进行慢启动，否则正在进行拥塞避免。慢启动一直持续到回到当拥塞发生时所处位置的半时候才停止（因为记录了在步骤 2 中制造麻烦的窗口大小的一半），然后转为执行拥塞避免。

慢启动算法初始设置 cwnd 为 1 个报文段，此后每收到一个确认就加 1。那样，这会使窗口按指数方式增长：发送 1 个报文段，然后是 2 个，接着是 4 个……。

## TCP 协议中的计时器

TCP 中有四种计时器（Timer），分别为：

1. 重传计时器：Retransmission Timer
2. 坚持计时器：Persistent Timer
3. 保活计时器：Keeplive Timer
4. 时间等待计时器：Timer_Wait Timer

### 重传计时器

大家都知道 TCP 是保证数据可靠传输的。怎么保证呢？带确认的重传机制。在滑动窗口协议中，接受窗口会在连续收到的包序列中的最后一个包向接收端发送一个 ACK，当网络拥堵的时候，发送端的数据包和接收端的 ACK 包都有可能丢失。TCP 为了保证数据可靠传输，就规定在重传的“时间片”到了以后，如果还没有收到对方的 ACK，就重发此包，以避免陷入无限等待中。

当 TCP 发送报文段时，就创建该特定报文的重传计时器。可能发生两种情况： 1.若在计时器截止时间到之前收到了对此特定报文段的确认，则撤销此计时器。 2.若在收到了对此特定报文段的确认之前计时器截止时间到，则重传此报文段，并将计时器复位。

### 持久计时器

先来考虑一下情景：

发送端向接收端发送数据包知道接受窗口填满了，然后接受窗口告诉发送方接受窗口填满了停止发送数据。此时的状态称为“**零窗口**”状态，发送端和接收端窗口大小均为 0.直到接受 TCP 发送确认并宣布一个非零的窗口大小。但这个确认会丢失。我们知道 TCP 中，对确认是不需要发送确认的。若确认丢失了，接受 TCP 并不知道，而是会认为他已经完成了任务，并等待着发送 TCP 接着会发送更多的报文段。但发送 TCP 由于没有收到确认，就等待对方发送确认来通知窗口大小。双方的 TCP 都在永远的等待着对方。

要打开这种死锁，TCP 为每一个链接使用一个持久计时器。当发送 TCP 收到窗口大小为 0 的确认时，就坚持启动计时器。当坚持计时器期限到时，发送 TCP 就发送一个特殊的报文段，叫做**探测报文**。这个报文段只有一个字节的数据。他有一个序号，但他的序号永远不需要确认；甚至在计算机对其他部分的数据的确认时该序号也被忽略。探测报文段提醒接受 TCP：确认已丢失，必须重传。

坚持计时器的值设置为重传时间的数值。但是，若没有收到从接收端来的响应，则需发送另一个探测报文段，并将坚持计时器的值加倍和复位。发送端继续发送探测报文段，将坚持计时器设定的值加倍和复位，直到这个值增大到门限值（通常是 60 秒）为止。在这以后，发送端每个 60 秒就发送一个探测报文，直到窗口重新打开。

### 保活计时器

保活计时器使用在某些实现中，用来防止在两个 TCP 之间的连接出现**长时间**的空闲。假定客户打开了到服务器的连接，传送了一些数据，然后就保持静默了。也许这个客户出故障了。在这种情况下，这个连接将永远的处理打开状态。

要解决这种问题，在大多数的实现中都是使服务器设置保活计时器。每当服务器收到客户的信息，就将计时器复位。

通常设置为**两小时**。若服务器过了两小时还没有收到客户的信息，他就发送探测报文段。若发送了**10 个**探测报文段（每一个相隔 75 秒）还没有响应，就假定客户除了故障，因而就终止了该连接。

这种连接的断开当然不会使用四次握手，而是直接硬性的中断和客户端的 TCP 连接。

### 时间等待计时器

时间等待计时器是在四次握手（挥手）的时候使用的。四次握手的简单过程是这样的：假设客户端准备中断连接，首先向服务器端发送一个 FIN 的请求关闭包（FIN=final），然后由 established 过渡到 FIN-WAIT1 状态。服务器收到 FIN 包以后会发送一个 ACK，然后自己有 established 进入 CLOSE-WAIT.此时通信进入**半双工状态**，即留给服务器一个机会将剩余数据传递给客户端，传递完后服务器发送一个 FIN+ACK 的包，表示我已经发送完数据可以断开连接了，就这便进入 LAST_ACK 阶段。客户端收到以后，发送一个 ACK 表示收到并同意请求，接着由 FIN-WAIT2 进入 TIME-WAIT 阶段。服务器收到 ACK，结束连接。此时（即客户端发送完 ACK 包之后），客户端还要等待 2MSL（MSL=maxinum segment lifetime 最长报文生存时间，2MSL 就是两倍的 MSL）才能真正的关闭连接。

![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/wireshark/1659010180989-f89f5334-9e7b-426f-a24b-20b6ec74b929.png)

附加：

![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/wireshark/1659010180979-c35b1d03-9cfd-4197-bb63-6168a4110ae5.png)
