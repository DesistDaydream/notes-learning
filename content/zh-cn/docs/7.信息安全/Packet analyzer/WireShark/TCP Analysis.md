---
title: TCP Analysis
linkTitle: TCP Analysis
date: 2024-04-11T16:32
weight: 20
---

# 概述

> 参考：
>
> - [官方文档，TCP Analysis](https://www.wireshark.org/docs/wsug_html_chunked/ChAdvTCPAnalysis.html)
> - [公众号-小林 coding，一文搞定 Wireshark 网络数据包分析](https://mp.weixin.qq.com/s/hL96imOvuodILIhI70fbTg)

使用 Wireshark 分析 [TCP](/docs/4.数据通信/Protocol/TCP_IP/TCP/TCP.md)。

## 通过 WireShark 理解三次握手和四次挥手

通过 WireShark 分析网络包，并理解 TCP 三次握手和四次挥手的工作原理。

本次例子，我们将要访问的 <http://192.168.3.200> 服务端。在终端用 tcpdump 命令抓取数据包：

```bash
# 客户端执行 tcpdump 命令抓包
tcpdump -i any tcp and host 192.168.3.200 and port 80 -w http.pcap
```

接着，在终端二执行下面的 curl 命令 `curl http://192.168.3.200`

最后，回到终端一，按下 Ctrl+C 停止 tcpdump，并把得到的 http.pcap 取出到电脑。

使用 Wireshark 打开 http.pcap 后，你就可以在 Wireshark 中，看到如下的界面：

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/wireshark/1616160823027-31c76057-664b-4c75-b496-01d9f48d775f.jpeg)

我们都知道 HTTP 是基于 TCP 协议进行传输的，那么：

- 最开始的 3 个包就是 TCP 三次握手建立连接的包
- 中间是 HTTP 请求和响应的包
- 而最后的 3 个包则是 TCP 断开连接的挥手包

Wireshark 可以用时序图的方式显示数据包交互的过程，从菜单栏中，点击 `统计 (Statistics) -> 流量图 (Flow Graph)`，然后，在弹出的界面中的 `「流量类型」选择 「TCP Flows」`，你可以更清晰的看到，整个过程中 TCP 流的执行过程：

![TCP 流量图](https://notes-learning.oss-cn-beijing.aliyuncs.com/wireshark/1616160823032-86a082df-77eb-4417-89d8-e952091e1f35.jpeg)

你可能会好奇，为什么三次握手连接过程的 Seq 是 0 ？

实际上是因为 Wireshark 工具帮我们做了优化，它默认显示的是序列号 seq 是相对值，而不是真实值。

如果你想看到实际的序列号的值，可以右键菜单， 然后找到「协议首选项」，接着找到「Relative Seq」后，把它给取消，操作如下：

![取消序列号相对值显示](https://notes-learning.oss-cn-beijing.aliyuncs.com/wireshark/1616160823088-c6c0d567-2a1b-4aff-9f26-9e7831bf4819.jpeg)

取消后，Seq 显示的就是真实值了：

![TCP 流量图](https://notes-learning.oss-cn-beijing.aliyuncs.com/wireshark/1616160823067-984c8d93-59af-43d2-b732-cd458926b9db.jpeg)

可见，客户端和服务端的序列号实际上是不同的，序列号是一个随机值。

这其实跟我们书上看到的 TCP 三次握手和四次挥手很类似，作为对比，你通常看到的 TCP 三次握手和四次挥手的流程，基本是这样的：

![TCP 三次握手和四次挥手的流程|500](https://notes-learning.oss-cn-beijing.aliyuncs.com/tcp/tcp-three-way-handshake-and-four-way-wave.png)

为什么抓到的 TCP 挥手是三次，而不是书上说的四次？

因为服务器端收到客户端的 FIN 后，服务器端同时也要关闭连接，这样就可以把 ACK 和 FIN 合并到一起发送，节省了一个包，变成了“三次挥手”。

而通常情况下，服务器端收到客户端的 FIN 后，很可能还没发送完数据，所以就会先回复客户端一个 ACK 包，稍等一会儿，完成所有数据包的发送后，才会发送 FIN 包，这也就是四次挥手了。

如下图，就是四次挥手的过程：

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/wireshark/1616160823036-56c913f2-2e69-4d3c-a6bc-18a549713c6b.jpeg)

## TCP 快速建立连接

客户端在向服务端发起 HTTP GET 请求时，一个完整的交互过程，需要 2.5 个 RTT 的时延。

由于第三次握手是可以携带数据的，这时如果在第三次握手发起 HTTP GET 请求，需要 2 个 RTT 的时延。

但是在下一次（不是同个 TCP 连接的下一次）发起 HTTP GET 请求时，经历的 RTT 也是一样，如下图：

![常规 HTTP 请求|600](https://notes-learning.oss-cn-beijing.aliyuncs.com/wireshark/1616160823176-61c2f1f2-9549-4e63-bcdf-1114b232e864.jpeg)

在 Linux 3.7 内核版本中，提供了 TCP Fast Open 功能，这个功能可以减少 TCP 连接建立的时延。

![常规 HTTP 请求 与 Fast Open HTTP 请求|600](https://notes-learning.oss-cn-beijing.aliyuncs.com/wireshark/1616160823163-e73ac138-f4ef-4436-b04e-1c24bb743a31.jpeg)

- 在第一次建立连接的时候，服务端在第二次握手产生一个 Cookie （已加密）并通过 SYN、ACK 包一起发给客户端，于是客户端就会缓存这个 Cookie，所以第一次发起 HTTP Get 请求的时候，还是需要 2 个 RTT 的时延；
- 在下次请求的时候，客户端在 SYN 包带上 Cookie 发给服务端，就提前可以跳过三次握手的过程，因为 Cookie 中维护了一些信息，服务端可以从 Cookie 获取 TCP 相关的信息，这时发起的 HTTP GET 请求就只需要 1 个 RTT 的时延；

注：客户端在请求并存储了 Fast Open Cookie 之后，可以不断重复 TCP Fast Open 直至服务器认为 Cookie 无效（通常为过期）

可以通过设置 net.ipv4.tcp_fastopn 内核参数，来打开 Fast Open 功能。

net.ipv4.tcp_fastopn 各个值的意义:

- 0 关闭
- 1 作为客户端使用 Fast Open 功能
- 2 作为服务端使用 Fast Open 功能
- 3 无论作为客户端还是服务器，都可以使用 Fast Open 功能

TCP Fast Open 抓包分析。在下图，数据包 7 号，客户端发起了第二次 TCP 连接时，SYN 包会携带 Cooike，并且有长度为 5 的数据。

服务端收到后，校验 Cooike 合法，于是就回了 SYN、ACK 包，并且确认应答收到了客户端的数据包，ACK = 5 + 1 = 6

![TCP Fast Open 抓包分析](https://notes-learning.oss-cn-beijing.aliyuncs.com/wireshark/1616160823189-cedf088a-663f-4fa8-9df7-e3357bde8c42.jpeg)

# TCP 提示信息

> 参考：
>
> - [官方文档，TCP Analysis](https://www.wireshark.org/docs/wsug_html_chunked/ChAdvTCPAnalysis.html)
> - 原文：[程序员宅基地，TCP报文（ tcp dup ack 、TCP Retransmission）](https://www.cxyzjd.com/article/ynchyong/109110028)
>   - [CSDN，TCP报文（ tcp dup ack 、TCP Retransmission）](https://blog.csdn.net/ynchyong/article/details/109110028)

WireShark 的 Info 列会根据包的前后文和内容分析该包的状态，以提示使用者该包可能有什么问题，常见的提示有：

- TCP ACKed unseen segment
- TCP Dup ACK
- TCP Fast Retransmission
- TCP Keep-Alive
- etc.

## TCP Dup ACK(重复确认) 和 TCP Fast Retransmission(快速重传)

`TCP Dup ACK <frame>#<acknowledgment number>` 表示第几次重新请求某一个包，frame 表示第几个包（不是 Seq），acknowledgment number 表示第几次请求。

> TODO: 细节说法待确认

**丢包**或者**乱序**的情况下，会出现该标志。

一般快速重传算法在收到**三次冗余的 Ack**，即三次 `TCP Dup ACK <frame>#<acknowledgment number>` 后，**发送端**进行快速重传。

为什么是三次呢？因为**两次 duplicated ACK** 肯定是**乱序**造成的，**丢包**肯定会造成**三次 duplicated ACK**。

当接收方收到乱序数据包时，会发送重复的 ACK，以使告知发送方要重发该数据包，当发送方收到 3 个重复 ACK 时，就会触发快速重传，立刻重发丢失数据包。

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/wireshark/1616160823137-f4bebd6f-37ff-40e8-bc73-6d5db00516b5.jpeg)

TCP 重复确认和快速重传的一个案例，用 Wireshark 分析，显示如下：

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/wireshark/1616160823162-89692a6f-2b62-4d7f-8f82-97e0b04c838f.jpeg)

- 数据包 1 期望的下一个数据包 Seq 是 1，但是数据包 2 发送的 Seq 却是 10945，说明收到的是乱序数据包，于是回了数据包 3 ，还是同样的 Seq = 1，Ack = 1，这表明是重复的 ACK；
- 数据包 4 和 6 依然是乱序的数据包，于是依然回了重复的 ACK；
- 当对方收到三次重复的 ACK 后，于是就快速重传了 Seq = 1 、Len = 1368 的数据包 8；
- 当收到重传的数据包后，发现 Seq = 1 是期望的数据包，于是就发送了确认报文 ACK；

注意：快速重传和重复 ACK 标记信息是 Wireshark 的功能，非数据包本身的信息。

以上案例在 TCP 三次握手时协商开启了选择性确认 SACK，因此一旦数据包丢失并收到重复 ACK ，即使在丢失数据包之后还成功接收了其他数据包，也只需要重传丢失的数据包。如果不启用 SACK，就必须重传丢失包之后的每个数据包。

如果要支持 SACK，必须双方都要支持。在 Linux 下，可以通过 net.ipv4.tcp_sack 参数打开这个功能（Linux 2.4 后默认打开）。

## TCP Out-Of-Order(报文乱序)

\[TCP Out-Of-Order] 指的是 TCP 发送端传输过程中报文乱序了。

有时候我们使用 tcpdump 抓包的时候 使用 -i any 时，若数据包经过了 bond 或者 bridge 类型的网络设备，也会出现这种情况，这是因为抓的 bond 设备与 bond salve 设备都发了包（或者 brdige 跳转到 brdige 的 salve 网络设备），所以 Wireshare 将不同网络设备的包当做乱序的。

![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/wireshark/1659010050504-b955aabe-b52a-4df1-b50c-04839ad1002f.png)

这是一个具有 Bond 设备的服务器，抓取了一次对 https://www.baidu.com 的请求的数据包，从 TCP 三次握手就可以看出来，每个包都是两份。

## TCP Retransmission(超时重传)

超时重传，如果一个包的丢了，又**没有后续包**可以在接收方触发\[Dup Ack]，或者 **\[Dup Ack]也丢失**的情况下，TCP 会触发超时重传机制。

![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/wireshark/1659010179388-13370a7d-1a0c-4611-a9e4-e2492cde492e.png)

## TCP Previous segment not captured

在 TCP 发送端传输过程中，该 Seq 前的报文缺失了。一般在网络拥塞的情况下，造成 TCP 报文乱序、丢包时，会出现该标志。

需要注意的是，\[TCP Previous segment not captured] 解析文字是 **wireshark 添加的标记**，**并非 TCP 报文内容**。

## TCP Window 控制

TCP 为了防止发送方无脑的发送数据，导致接收方缓冲区被填满，所以就有了滑动窗口的机制，它可利用接收方的接收窗口来控制发送方要发送的数据量，也就是流量控制。

接收窗口是由接收方指定的值，存储在 TCP 头部中，它可以告诉发送方自己的 TCP 缓冲空间区大小，这个缓冲区是给应用程序读取数据的空间：

- 如果应用程序读取了缓冲区的数据，那么缓冲空间区的就会把被读取的数据移除
- 如果应用程序没有读取数据，则数据会一直滞留在缓冲区。

接收窗口的大小，是在 TCP 三次握手中协商好的，后续数据传输时，接收方发送确认应答 ACK 报文时，会携带当前的接收窗口的大小，以此来告知发送方。

假设接收方接收到数据后，应用层能很快的从缓冲区里读取数据，那么窗口大小会一直保持不变，过程如下：

![理想状态下的窗口变化|500](https://notes-learning.oss-cn-beijing.aliyuncs.com/wireshark/1616160823208-5740027e-4b25-4e05-985e-4f5b3aac0181.jpeg)

但是现实中服务器会出现繁忙的情况，当应用程序读取速度慢，那么缓存空间会慢慢被占满，于是为了保证发送方发送的数据不会超过缓冲区大小，则服务器会调整窗口大小的值，接着通过 ACK 报文通知给对方，告知现在的接收窗口大小，从而控制发送方发送的数据大小。

![服务端繁忙状态下的窗口变化|500](https://notes-learning.oss-cn-beijing.aliyuncs.com/wireshark/1616160823169-a3dfea42-74f4-4a5f-8836-ff78dca0557f.jpeg)

### TCP ZeroWindow(零窗口) 与 窗口探测

TCP ZeroWindow 是作为**接收方**发出现的标志，表示**接收缓冲区已经满**了。假设接收方处理数据的速度跟不上接收数据的速度，缓存就会被占满，从而导致接收窗口为 0，当发送方接收到零窗口通知时，就会停止发送数据。**wireshark 会打上\[TCP ZeroWindow]标签**。

如下图，可以看到接收方的窗口大小在不断的收缩至 0：

![窗口大小在收缩](https://notes-learning.oss-cn-beijing.aliyuncs.com/wireshark/1616160823151-9f129444-b09c-4bd6-be72-3dafdac2d9c6.jpeg)

接着，发送方会定时发送窗口大小探测报文，以便及时知道接收方窗口大小的变化。

以下图 Wireshark 分析图作为例子说明：

![零窗口 与 窗口探测](https://notes-learning.oss-cn-beijing.aliyuncs.com/wireshark/1616160823171-33894ce8-f739-488e-965e-8e64b267ec30.jpeg)

- 发送方发送了数据包 1 给接收方，接收方收到后，由于缓冲区被占满，回了个零窗口通知；
- 发送方收到零窗口通知后，就不再发送数据了，直到过了 3.4 秒后，发送了一个 TCP Keep-Alive 报文，也就是窗口大小探测报文；
- 当接收方收到窗口探测报文后，就立马回一个窗口通知，但是窗口大小还是 0；
- 发送方发现窗口还是 0，于是继续等待了 6.8（翻倍） 秒后，又发送了窗口探测报文，接收方依然还是回了窗口为 0 的通知；
- 发送方发现窗口还是 0，于是继续等待了 13.5（翻倍） 秒后，又发送了窗口探测报文，接收方依然还是回了窗口为 0 的通知；

可以发现，这些窗口探测报文以 3.4s、6.5s、13.5s 的间隔出现，说明超时时间会翻倍递增。

这连接暂停了 25s，想象一下你在打王者的时候，25s 的延迟你还能上王者吗？

### 发送窗口的分析

在 Wireshark 看到的 Windows size 也就是 " win = "，这个值表示发送窗口吗？

这不是发送窗口，而是在向对方声明自己的接收窗口。

你可能会好奇，抓包文件里有「Window size scaling factor」，它其实是算出实际窗口大小的乘法因子，「Windos size value」实际上并不是真实的窗口大小，真实窗口大小的计算公式如下：

`「Windos size value」 * 「Window size scaling factor」 = 「Caculated window size 」`

对应的下图案例，也就是 32 \* 2048 = 65536。

![700](https://notes-learning.oss-cn-beijing.aliyuncs.com/wireshark/1616160823197-f831c96d-7d0a-4acb-9a6c-e62ce143f763.jpeg)

实际上是 Caculated window size 的值是 Wireshark 工具帮我们算好的，Window size scaling factor 和 Windos size value 的值是在 TCP 头部中，其中 Window size scaling factor 是在三次握手过程中确定的，如果你抓包的数据没有 TCP 三次握手，那可能就无法算出真实的窗口大小的值，如下图：

![700](https://notes-learning.oss-cn-beijing.aliyuncs.com/wireshark/1616160823193-da2d3691-94dc-4076-8461-8a776fc69885.jpeg)

如何在包里看出发送窗口的大小？

- 很遗憾，没有简单的办法，发送窗口虽然是由接收窗口决定，但是它又可以被网络因素影响，也就是拥塞窗口，实际上发送窗口是值是 min(拥塞窗口，接收窗口)。

发送窗口和 MSS 有什么关系？

- 发送窗口决定了一口气能发多少字节，而 MSS 决定了这些字节要分多少包才能发完。
- 举个例子，如果发送窗口为 16000 字节的情况下，如果 MSS 是 1000 字节，那就需要发送 1600/1000 = 16 个包。

发送方在一个窗口发出 n 个包，是不是需要 n 个 ACK 确认报文？

- 不一定，因为 TCP 有累计确认机制，所以当收到多个数据包时，只需要应答最后一个数据包的 ACK 报文就可以了。

### TCP window update

接收方消耗缓冲数据后，更新 TCP 窗口， 可以看到从 win=0 逐渐变大，这时 **wireshark 会打上\[TCP window update]** 标签

### TCP window Full

作为**发送方的标识**，当前**发送包的大小已经超过了接收端窗口大小**，wireshark 会打上此标识，标识不能在发送。

## TCP 延迟确认与 Nagle 算法

当我们 TCP 报文的承载的数据非常小的时候，例如几个字节，那么整个网络的效率是很低的，因为每个 TCP 报文中都有会 20 个字节的 TCP 头部，也会有 20 个字节的 IP 头部，而数据只有几个字节，所以在整个报文中有效数据占有的比重就会非常低。

这就好像快递员开着大货车送一个小包裹一样浪费。

那么就出现了常见的两种策略，来减少小报文的传输，分别是：

- Nagle 算法
- 延迟确认

Nagle 算法是如何避免大量 TCP 小数据报文的传输？

Nagle 算法做了一些策略来避免过多的小数据报文发送，这可提高传输效率。

Nagle 算法的策略：

- 没有已发送未确认报文时，立刻发送数据。
- 存在未确认报文时，直到「没有已发送未确认报文」或「数据长度达到 MSS 大小」时，再发送数据。

只要没满足上面条件中的一条，发送方一直在囤积数据，直到满足上面的发送条件。

![禁用 Nagle 算法 与 启用 Nagle 算法|600](https://notes-learning.oss-cn-beijing.aliyuncs.com/wireshark/1616160823193-330e95bc-5b83-4a95-92c9-d301e3148437.jpeg)

上图右侧启用了 Nagle 算法，它的发送数据的过程：

- 一开始由于没有已发送未确认的报文，所以就立刻发了 H 字符；
- 接着，在还没收到对 H 字符的确认报文时，发送方就一直在囤积数据，直到收到了确认报文后，此时就没有已发送未确认的报文，于是就把囤积后的 ELL 字符一起发给了接收方；
- 待收到对 ELL 字符的确认报文后，于是把最后一个 O 字符发送出去

可以看出，Nagle 算法一定会有一个小报文，也就是在最开始的时候。

另外，Nagle 算法默认是打开的，如果对于一些需要小数据包交互的场景的程序，比如，telnet 或 ssh 这样的交互性比较强的程序，则需要关闭 Nagle 算法。

可以在 Socket 设置 TCP_NODELAY 选项来关闭这个算法（关闭 Nagle 算法没有全局参数，需要根据每个应用自己的特点来关闭）。

![关闭 Nagle 算法](https://notes-learning.oss-cn-beijing.aliyuncs.com/wireshark/202311011544423.png)

那延迟确认又是什么？

事实上当没有携带数据的 ACK，他的网络效率也是很低的，因为它也有 40 个字节的 IP 头 和 TCP 头，但没有携带数据。

为了解决 ACK 传输效率低问题，所以就衍生出了 TCP 延迟确认。

TCP 延迟确认的策略：

- 当有响应数据要发送时，ACK 会随着响应数据一起立刻发送给对方
- 当没有响应数据要发送时，ACK 将会延迟一段时间，以等待是否有响应数据可以一起发送
- 如果在延迟等待发送 ACK 期间，对方的第二个数据报文又到达了，这时就会立刻发送 ACK

![TCP 延迟确认|600](https://notes-learning.oss-cn-beijing.aliyuncs.com/wireshark/1616160823199-a0f7bdbb-e7e2-4fb3-b14f-6d1c517c96b3.jpeg)

延迟等待的时间是在 Linux 内核中的定义的，如下图：

![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/wireshark/202311011548648.png)

关键就需要 HZ 这个数值大小，HZ 是跟系统的时钟频率有关，每个操作系统都不一样，在我的 Linux 系统中 HZ 大小是 1000，如下：

```bash
cat /boot/config-2.6.32-431.el6.x86_64 | grep '^CONFIG_HZ='
CONFIG_HZ=1000
```

知道了 HZ 的大小，那么就可以算出：

- 最大延迟确认时间是 200 ms （1000/5）
- 最短延迟确认时间是 40 ms （1000/25）

TCP 延迟确认可以在 Socket 设置 TCP_QUICKACK 选项来关闭这个算法。

![关闭 TCP 延迟确认](https://notes-learning.oss-cn-beijing.aliyuncs.com/wireshark/202311011549867.png)

延迟确认 和 Nagle 算法混合使用时，会产生新的问题

当 TCP 延迟确认 和 Nagle 算法混合使用时，会导致时耗增长，如下图：

![TCP 延迟确认 和 Nagle 算法混合使用|500](https://notes-learning.oss-cn-beijing.aliyuncs.com/wireshark/1616160823201-97a64317-4a76-449e-b319-1dca774ee157.jpeg)

发送方使用了 Nagle 算法，接收方使用了 TCP 延迟确认会发生如下的过程：

- 发送方先发出一个小报文，接收方收到后，由于延迟确认机制，自己又没有要发送的数据，只能干等着发送方的下一个报文到达；
- 而发送方由于 Nagle 算法机制，在未收到第一个报文的确认前，是不会发送后续的数据；
- 所以接收方只能等待最大时间 200 ms 后，才回 ACK 报文，发送方收到第一个报文的确认报文后，也才可以发送后续的数据。

很明显，这两个同时使用会造成额外的时延，这就会使得网络"很慢"的感觉。

要解决这个问题，只有两个办法：

- 要么发送方关闭 Nagle 算法
- 要么接收方关闭 TCP 延迟确认

# TCP 三次握手异常情况实战分析

TCP 三次握手的这三个异常情况：

- TCP 第一次握手的 SYN 丢包了，会发生了什么？
- TCP 第二次握手的 SYN、ACK 丢包了，会发生什么？
- TCP 第三次握手的 ACK 包丢了，会发生什么？

再继续问：

- 那会重传几次？
- 超时重传的时间 RTO 会如何变化？
- 在 Linux 下如何设置重传次数？
- ….

接下里我用三个实验案例，带大家一起探究探究这三种异常。

实验场景

本次实验用了两台虚拟机，一台作为服务端，一台作为客户端，它们的关系如下：

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/wireshark/1616160823067-15e439a6-ae99-47c0-a3dc-2de0b9cce21e.jpeg)

实验环境

- 客户端和服务端都是 CentOs 6.5 Linux，Linux 内核版本 2.6.32
- 服务端 192.168.12.36，apache web 服务
- 客户端 192.168.12.37

### 实验一：TCP 第一次握手 SYN 丢包

为了模拟 TCP 第一次握手 SYN 丢包的情况，我是在拔掉服务器的网线后，立刻在客户端执行 curl 命令：

![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/wireshark/202311011500940.png)

在客户端执行 tcpdump 命令，抓取访问服务端的 HTTP 数据包

```
tcpdump -i eht0 tcp and host 192.168.12.36 and port 80 -w tcp_sys_timeout.pcap
```

过了一会， curl 返回了超时连接的错误：

![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/wireshark/202311011501520.png)

从 date 返回的时间，可以发现在超时接近 1 分钟的时间后，curl 返回了错误。

接着，把 tcp_sys_timeout.pcap 文件用 Wireshark 打开分析，显示如下图：

![SYN 超时重传五次](https://notes-learning.oss-cn-beijing.aliyuncs.com/wireshark/1616160823088-ecfeac4f-b707-4cc0-ac5e-7a777e6c4e21.jpeg)

从上图可以发现， 客户端发起了 SYN 包后，一直没有收到服务端的 ACK ，所以一直超时重传了 5 次，并且每次 RTO 超时时间是不同的：

- 第一次是在 1 秒超时重传
- 第二次是在 3 秒超时重传
- 第三次是在 7 秒超时重传
- 第四次是在 15 秒超时重传
- 第五次是在 31 秒超时重传

可以发现，每次超时时间 RTO 是指数（翻倍）上涨的，当超过最大重传次数后，客户端不再发送 SYN 包。

在 Linux 中，第一次握手的 SYN 超时重传次数，是如下内核参数指定的：

```
$ cat /c/sys/net/ipv4/tcp_syn_retries
5
```

`tcp_syn_retries` 默认值为 5，也就是 SYN 最大重传次数是 5 次。

接下来，我们继续做实验，把 tcp_syn_retries 设置为 2 次：

```
echo 2 > /proc/sys/net/ipv4/tcp_syn_retries
```

重传抓包后，用 Wireshark 打开分析，显示如下图：

![SYN 超时重传两次](https://notes-learning.oss-cn-beijing.aliyuncs.com/wireshark/1616160823073-af405abf-7d5a-4062-9987-a7233d7c1938.jpeg)

#### 实验一的实验小结

通过实验一的实验结果，我们可以得知，当客户端发起的 TCP 第一次握手 SYN 包，在超时时间内没收到服务端的 ACK，就会在超时重传 SYN 数据包，每次超时重传的 RTO 是翻倍上涨的，直到 SYN 包的重传次数到达 tcp_syn_retries 值后，客户端不再发送 SYN 包。

![SYN 超时重传|200](https://notes-learning.oss-cn-beijing.aliyuncs.com/wireshark/1616160823075-f999dcb2-2aec-494a-aba9-54d059717f02.jpeg)

### 实验二：TCP 第二次握手 SYN、ACK 丢包

为了模拟客户端收不到服务端第二次握手 SYN、ACK 包，我的做法是在客户端加上防火墙限制，直接粗暴的把来自服务端的数据都丢弃，客户端防火墙的配置如下：

`iptables -I INPUT -s 192.168.12.36 -j DROP`

接着，在客户端执行 curl 命令

![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/wireshark/202311011504547.png)

从 date 返回的时间前后，可以算出大概 1 分钟后，curl 报错退出了。

客户端在这其间抓取的数据包，用 Wireshark 打开分析，显示的时序图如下：

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/wireshark/1616160823111-a8b3d7e4-eabe-4c1f-a22b-1ddd0e125afa.jpeg)

从图中可以发现：

- 客户端发起 SYN 后，由于防火墙屏蔽了服务端的所有数据包，所以 curl 是无法收到服务端的 SYN、ACK 包，当发生超时后，就会重传 SYN 包
- 服务端收到客户的 SYN 包后，就会回 SYN、ACK 包，但是客户端一直没有回 ACK，服务端在超时后，重传了 SYN、ACK 包，接着一会，客户端超时重传的 SYN 包又抵达了服务端，服务端收到后，超时定时器就重新计时，然后回了 SYN、ACK 包，所以相当于服务端的超时定时器只触发了一次，又被重置了。
- 最后，客户端 SYN 超时重传次数达到了 5 次（tcp_syn_retries 默认值 5 次），就不再继续发送 SYN 包了。

所以，我们可以发现，当第二次握手的 SYN、ACK 丢包时，客户端会超时重发 SYN 包，服务端也会超时重传 SYN、ACK 包。

咦？客户端设置了防火墙，屏蔽了服务端的网络包，为什么 tcpdump 还能抓到服务端的网络包？

添加 iptables 限制后， tcpdump 是否能抓到包 ，这要看添加的 iptables 限制条件：

- 如果添加的是 INPUT 规则，则可以抓得到包
- 如果添加的是 OUTPUT 规则，则抓不到包

网络包进入主机后的顺序如下：

- 进来的顺序 Wire -> NIC -> tcpdump -> netfilter/iptables
- 出去的顺序 iptables -> tcpdump -> NIC -> Wire

tcp_syn_retries 是限制 SYN 重传次数，那第二次握手 SYN、ACK 限制最大重传次数是多少？

TCP 第二次握手 SYN、ACK 包的最大重传次数是通过 tcp_synack_retries 内核参数限制的，其默认值如下：

```
$ cat /proc/sys/net/ipv4/tcp_synack_retries
5
```

是的，TCP 第二次握手 SYN、ACK 包的最大重传次数默认值是 5 次。

为了验证 SYN、ACK 包最大重传次数是 5 次，我们继续做下实验，我们先把客户端的 tcp_syn_retries 设置为 1，表示客户端 SYN 最大超时次数是 1 次，目的是为了防止多次重传 SYN，把服务端 SYN、ACK 超时定时器重置。

接着，还是如上面的步骤：

1. 客户端配置防火墙屏蔽服务端的数据包
2. 客户端 tcpdump 抓取 curl 执行时的数据包

把抓取的数据包，用 Wireshark 打开分析，显示的时序图如下：

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/wireshark/1616160823110-6ad95c7d-5b11-4b5d-a160-0f17e55e4a20.jpeg)
从上图，我们可以分析出：

- 客户端的 SYN 只超时重传了 1 次，因为 tcp_syn_retries 值为 1
- 服务端应答了客户端超时重传的 SYN 包后，由于一直收不到客户端的 ACK 包，所以服务端一直在超时重传 SYN、ACK 包，每次的 RTO 也是指数上涨的，一共超时重传了 5 次，因为 tcp_synack_retries 值为 5

接着，我把 tcp_synack_retries 设置为 2，tcp_syn_retries 依然设置为 1:

```
echo 2 > /proc/sys/net/ipv4/tcp_synack_retries
echo 1 > /proc/sys/net/ipv4/tcp_syn_retries
```

依然保持一样的实验步骤进行操作，接着把抓取的数据包，用 Wireshark 打开分析，显示的时序图如下：

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/wireshark/1616160823097-cbfe5e48-374e-4bb9-a651-510ab2cd9826.jpeg)
可见：

- 客户端的 SYN 包只超时重传了 1 次，符合 tcp_syn_retries 设置的值；
- 服务端的 SYN、ACK 超时重传了 2 次，符合 tcp_synack_retries 设置的值

#### 实验二的实验小结

通过实验二的实验结果，我们可以得知，当 TCP 第二次握手 SYN、ACK 包丢了后，客户端 SYN 包会发生超时重传，服务端 SYN、ACK 也会发生超时重传。

客户端 SYN 包超时重传的最大次数，是由 tcp_syn_retries 决定的，默认值是 5 次；服务端 SYN、ACK 包时重传的最大次数，是由 tcp_synack_retries 决定的，默认值是 5 次。

### 实验三：TCP 第三次握手 ACK 丢包

为了模拟 TCP 第三次握手 ACK 包丢，我的实验方法是在服务端配置防火墙，屏蔽客户端 TCP 报文中标志位是 ACK 的包，也就是当服务端收到客户端的 TCP ACK 的报文时就会丢弃，服务端 iptables 配置命令如下

```bash
iptables -I INPUT -s 192.168.12.37 -p tcp --tcp-flag ACK ACK -j DROP
```

接着，在客户端执行如下 tcpdump 命令：

```bash
tcpdump -i eth0 tcp and host 192.168.12.36 and port 80 -w tcp_thir_ack_timeout.pacp
```

然后，客户端向服务端发起 telnet，因为 telnet 命令是会发起 TCP 连接，所以用此命令做测试：

![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/wireshark/202311011507328.png)

此时，由于服务端收不到第三次握手的 ACK 包，所以一直处于 SYN_RECV 状态：

![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/wireshark/202311011508124.png)

而客户端是已完成 TCP 连接建立，处于 ESTABLISHED 状态：

![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/wireshark/202311011508277.png)

过了 1 分钟后，观察发现服务端的 TCP 连接不见了：

![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/wireshark/202311011509533.png)

过了 30 分别，客户端依然还是处于 ESTABLISHED 状态：

![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/wireshark/202311011509022.png)

接着，在刚才客户端建立的 telnet 会话，输入 123456 字符，进行发送：

![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/wireshark/202311011509202.png)

持续「好长」一段时间，客户端的 telnet 才断开连接：

![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/wireshark/202311011510759.png)

以上就是本次实验三的现象，这里存在两个疑点：

- 为什么服务端原本处于 SYN_RECV 状态的连接，过 1 分钟后就消失了？
- 为什么客户端 telnet 输入 123456 字符后，过了好长一段时间，telnet 才断开连接？

不着急，我们把刚抓的数据包，用 Wireshark 打开分析，显示的时序图如下：

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/wireshark/1616160823161-35d5ac76-4d22-4cf8-a22f-2aec49af8c71.jpeg)

上图的流程：

- 客户端发送 SYN 包给服务端，服务端收到后，回了个 SYN、ACK 包给客户端，此时服务端的 TCP 连接处于 SYN_RECV 状态；
- 客户端收到服务端的 SYN、ACK 包后，给服务端回了个 ACK 包，此时客户端的 TCP 连接处于 ESTABLISHED 状态；
- 由于服务端配置了防火墙，屏蔽了客户端的 ACK 包，所以服务端一直处于 SYN_RECV 状态，没有进入 ESTABLISHED 状态，tcpdump 之所以能抓到客户端的 ACK 包，是因为数据包进入系统的顺序是先进入 tcpudmp，后经过 iptables；
- 接着，服务端超时重传了 SYN、ACK 包，重传了 5 次后，也就是超过 tcp_synack_retries 的值（默认值是 5），然后就没有继续重传了，此时服务端的 TCP 连接主动中止了，所以刚才处于 SYN_RECV 状态的 TCP 连接断开了，而客户端依然处于 ESTABLISHED 状态；
- 虽然服务端 TCP 断开了，但过了一段时间，发现客户端依然处于 ESTABLISHED 状态，于是就在客户端的 telnet 会话输入了 123456 字符；
- 此时由于服务端已经断开连接，客户端发送的数据报文，一直在超时重传，每一次重传，RTO 的值是指数增长的，所以持续了好长一段时间，客户端的 telnet 才报错退出了，此时共重传了 15 次。

通过这一波分析，刚才的两个疑点已经解除了：

- 服务端在重传 SYN、ACK 包时，超过了最大重传次数 tcp_synack_retries，于是服务端的 TCP 连接主动断开了。
- 客户端向服务端发送数据包时，由于服务端的 TCP 连接已经退出了，所以数据包一直在超时重传，共重传了 15 次， telnet 就 断开了连接。

TCP 第一次握手的 SYN 包超时重传最大次数是由 tcp_syn_retries 指定，TCP 第二次握手的 SYN、ACK 包超时重传最大次数是由 tcp_synack_retries 指定，那 TCP 建立连接后的数据包最大超时重传次数是由什么参数指定呢？

TCP 建立连接后的数据包传输，最大超时重传次数是由 tcp_retries2 指定，默认值是 15 次，如下：

```bash
$ cat /proc/sys/net/ipv4/tcp_retries2
15
```

如果 15 次重传都做完了，TCP 就会告诉应用层说：“搞不定了，包怎么都传不过去！”

那如果客户端不发送数据，什么时候才会断开处于 ESTABLISHED 状态的连接？

这里就需要提到 TCP 的 保活机制。这个机制的原理是这样的：

定义一个时间段，在这个时间段内，如果没有任何连接相关的活动，TCP 保活机制会开始作用，每隔一个时间间隔，发送一个「探测报文」，该探测报文包含的数据非常少，如果连续几个探测报文都没有得到响应，则认为当前的 TCP 连接已经死亡，系统内核将错误信息通知给上层应用程序。

在 Linux 内核可以有对应的参数可以设置保活时间、保活探测的次数、保活探测的时间间隔，以下都为默认值：

```bash
net.ipv4.tcp_keepalive_time=7200
net.ipv4.tcp_keepalive_intvl=75
net.ipv4.tcp_keepalive_probes=9
```

- tcp_keepalive_time=7200：表示保活时间是 7200 秒（2 小时），也就 2 小时内如果没有任何连接相关的活动，则会启动保活机制
- tcp_keepalive_intvl=75：表示每次检测间隔 75 秒；
- tcp_keepalive_probes=9：表示检测 9 次无响应，认为对方是不可达的，从而中断本次的连接。

也就是说在 Linux 系统中，最少需要经过 2 小时 11 分 15 秒才可以发现一个「死亡」连接。

![700](https://notes-learning.oss-cn-beijing.aliyuncs.com/wireshark/1616160823144-22ffe810-8621-4ffb-bcef-c2721b441f0b.jpeg)

这个时间是有点长的，所以如果我抓包足够久，或许能抓到探测报文。

#### 实验三的实验小结

在建立 TCP 连接时，如果第三次握手的 ACK，服务端无法收到，则服务端就会短暂处于 SYN_RECV 状态，而客户端会处于 ESTABLISHED 状态。

由于服务端一直收不到 TCP 第三次握手的 ACK，则会一直重传 SYN、ACK 包，直到重传次数超过 tcp_synack_retries 值（默认值 5 次）后，服务端就会断开 TCP 连接。

而客户端则会有两种情况：

- 如果客户端没发送数据包，一直处于 ESTABLISHED 状态，然后经过 2 小时 11 分 15 秒才可以发现一个「死亡」连接，于是客户端连接就会断开连接。
- 如果客户端发送了数据包，一直没有收到服务端对该数据包的确认报文，则会一直重传该数据包，直到重传次数超过 tcp_retries2 值（默认值 15 次）后，客户端就会断开 TCP 连接。
