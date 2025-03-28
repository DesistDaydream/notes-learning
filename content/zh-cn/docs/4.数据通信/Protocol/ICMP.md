---
title: ICMP
linkTitle: ICMP
weight: 3
---

# 概述

> 参考：
>
> - [公众号,24 张图搞定 ICMP](https://mp.weixin.qq.com/s/AKiUyMbsGhOZi7cDAhSGkg)

### ICMP

**IP** 是尽力传输的网络协议，提供的数据传输服务是**不可靠**的、无连接的，不能保证数据包能成功到达目的地。那么问题来了：如何确定数据包成功到达目的地？
![](https://notes-learning.oss-cn-beijing.aliyuncs.com/boov5o/1622087956279-d468d429-a172-4d1d-8296-bd855385f80b.png)不可靠传输
这需要一个网络层协议，提供错误检测功能和报告机制功能，于是出现了 **ICMP**（互联网控制消息协议）。ICMP 的主要功能是，**确认 IP 包是否成功送达目的地址**，**通知发送过程中 IP 包被丢弃的原因**。有了这些功能，就可以检查网络是否正常、网络配置是否正确、设备是否异常等信息，方便进行**网络问题诊断**。
![](https://notes-learning.oss-cn-beijing.aliyuncs.com/boov5o/1622087956391-eae4cc3c-7eb0-47aa-875a-154dec57c1cf.png)ICMP 网络诊断功能
**举个栗子**：如果在传输过程中，发生了某个错误，设备便会向源设备返回一条 ICMP 消息，告诉它发生的错误类型。
![](https://notes-learning.oss-cn-beijing.aliyuncs.com/boov5o/1622087956489-289bd365-0f3c-4891-b798-016dcab13bb4.png)ICMP 举例
ICMP 消息是通过 IP 进行传输，但它的目的并不是让 IP 成为一种可靠的协议，而是对传输中发生的问题进行反馈。ICMP 消息的传输同样得不到可靠性保证，也有可能在传输过程中丢失。因此 ICMP 不是传输层的补充，应该把它当做**网络层协议**。

#### ICMP 消息封装

ICMP 消息使用 IP 来封装，**封装格式**如下图。
![](https://notes-learning.oss-cn-beijing.aliyuncs.com/boov5o/1622087956425-eb063720-58a0-4b78-bfa8-39c2e96d569d.png)ICMP 封装格式
其中 **type**（类型）字段表示 ICMP 消息的类型，**code**（代码）字段表示 ICMP 消息的具体含义。例如：type 值为 3 表示目的不可达消息（ Destination Unreachable Message ），若 code 值为 0 表示目的网络不可达（ Network Unreachable ）。常见的 ICMP 消息类型如下图。
![](https://notes-learning.oss-cn-beijing.aliyuncs.com/boov5o/1622087956257-cc3d7e06-5c0f-4f20-9026-48917ba962ee.png)ICMP 消息类型
从功能上，ICMP 的消息可分为两类：一类是通知出错原因的**错误消息**，另一类是用于诊断的**查询消息**。
![](https://notes-learning.oss-cn-beijing.aliyuncs.com/boov5o/1622087956583-0a4b1076-322c-4b5f-bd43-8b5364ce42a4.png)错误消息和查询消息

#### 常见的 ICMP 消息类型

- **回送请求消息**（ Echo Request ）：是由源设备（主机或路由器等）向一个指定的目的设备发出的请求。这种消息用来测试目的地是否可达。
- **回送响应消息**（ Echo Reply ）：对 Echo Request 的响应。目的设备发送 Echo Reply 来响应收到的 Echo Request 。最常用的 ping 命令就是使用 Echo Request 和 Echo Reply 来实现的。

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/boov5o/1622087956612-4e223795-e42d-4c78-8be6-48000ccc7632.png)回送消息

- **目的不可达**（ Destination Unreachable ）：路由器无法将 IP 包发送给目的地址时，会给源设备返回一个 Destination Unreachable 消息，并在消息中显示不可达的具体原因。
  ![](https://notes-learning.oss-cn-beijing.aliyuncs.com/boov5o/1622087956380-16033d96-d332-4f13-9a3a-3d24a50b113d.png)目的不可达实际情况下，经常会遇到的错误代码是 1 ，表示主机不可达，它是指路由表中没有目的设备的信息，或目的设备没有连接到网络。

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/boov5o/1622087956460-afdf7e80-77ee-4e32-8e0d-749e9bc5c483.png)目的不可达类型

- **参数问题**（ Parameter Problem ）：路由器发现 IP 包头出现错误或非法值后，向源设备发送一个 Parameter Problem 消息。这个消息包含有问题的 IP 头，或错误字段的提示信息。

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/boov5o/1622087956530-5dcf521e-a938-4961-b7b1-00d735a998b6.png)参数问题消息

- **重定向**（ Redirect ）：如果路由器发现一条更优的路径发送数据，那么它就会返回一个 Redirect 消息给主机。这个消息包含了最合适的路由信息和源数据。
  实际情况下，这种 Redirect 消息会引发路由问题，所以不进行这种设置。比如：路由器的路由表不准确时，ICMP 有可能就无法正常工作。

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/boov5o/1622087956439-255f3c69-1653-431d-81a4-9592520ed469.png)重定向

- **超时**（ Time Exceeded ）：IP 包中有一个字段是 TTL（生存周期），它的值每经过一次路由器就减 1 ，直到减到 0 时 IP 包会被丢弃。这时，路由器会发送一个 Time Exceeded 消息给源设备，并通知 IP 包已被丢弃。
  设置 TTL 的主要目的，是当路由发生环路时，避免 IP 包无休止的在网络上转发。还可以用 TTL 控制 IP 包的可达范围，比如设置一个较小的 TTL 值。

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/boov5o/1622087956563-69d4450f-834a-4f6e-9ae0-f54f6ad81866.png)超时

- **时间戳请求/时间戳响应**（ Timestamp Request / Timestamp Reply ）：时间戳可以记录 ICMP 消息一次往返所需的时间。源设备发送一个带有发送时间的 Timestamp Request 消息，目的设备收到后，发送一个带有原设备发送时间、目的设备接收时间以及目的设备发送时间的 Timestamp Reply 消息。源设备收到 Timestamp Reply 时，并同时记录到达时间。这些时间戳可以估计网络上的传输时间。

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/boov5o/1622087956537-0334836c-2e44-4f22-83c5-b81cf4d3031a.png)时间戳

### ICMP 的应用

ICMP 被广泛应用于网络测试，最常用的 **ping** 和 **tracert** 网络测试工具，都是使用 ICMP 协议实现的。

#### ping

ping 是 ICMP 最著名的一个应用，通过 ping 可以**测试网络的可达性**，即网络上的报文能否成功到达目的地。使用 ping 命令时，源设备向目的设备发送 _Echo request_ 消息，目的地址是目的设备的 IP 地址。目的设备收到 _Echo request_ 消息后，向源设备回应一个 _Echo reply_ 消息，可知目的设备是可达的。也可以通过 ping 命令来判断目标主机是否启用。
![](https://notes-learning.oss-cn-beijing.aliyuncs.com/boov5o/1622087956406-67c2d46c-ad49-4904-8ead-569eaeb0ef80.png)ping
如果中间某个路由器没有到达目的网络的路由，便会向源设备回应一个 _Destination Unreachable_ 消息，告知目的设备不可达。
![](https://notes-learning.oss-cn-beijing.aliyuncs.com/boov5o/1622087956546-e4b22038-cd77-4ba6-a18a-8f58b0961db2.png)ping 目的不可达
如果源主机在一定时间内无法收到回应报文，就认为目的设备不可达，并显示超时。
![](https://notes-learning.oss-cn-beijing.aliyuncs.com/boov5o/1622087956386-8c862599-839b-4eba-a995-b9267f202d65.png)超时
需要注意的是 ping 过程是双向的消息通信，只有双向都成功传输时，才能说明通信是正常的。另外主机也可能因为防火墙拦截，导致 ping 不通。

#### tracert

ping 工具只能测试目的设备的连通性，但是看不到数据包的传输路径。所以在网络不通的情况下，无法知道网络问题发生在哪个位置。tracert 工具可以查看数据包的**整条传输路径**，包括途中经过的**中间设备**。
![](https://notes-learning.oss-cn-beijing.aliyuncs.com/boov5o/1622087956431-46f8e00a-d012-428d-8bcc-247f97e0b7a6.png)tracert
IP 头部的 **TTL** 字段是为避免数据包循环转发而设计的。每经过一个路由器，数据包头中的 TTL 值减 1 。如果 TTL 值为 0 则丢弃报文，并向源设备回应一个 Time Exceeded 消息，告知错误类型。tracert 就是基于 TTL 字段和 ICMP 协议实现的。在 Windows 中命令是 **tracert** ，在 Unix 、MacOS 中命令是 **traceroute** 。
使用 tracert 命令时，源设备的 tracert **逐跳发送数据包**，并等待每一个响应报文。发送第一个数据包时，TTL 值设为 1 。第一个路由器收到数据包后 TTL 值减 1 ，随即丢弃数据包，并返回一个 _Time Exceeded_ 消息。源设备的 tracert 收到响应报文后，取出源 IP 地址，即路径上的第一个路由器地址。然后 tracert 发送一个 TTL 值为 2 的数据包。第一个路由器将 TTL 值减 1 ，并转发数据包。第二个路由器再将 TTL 值减 1 ，丢弃数据包并返回一个 _Time Exceeded_ 消息。tracert 收到响应报文后，取出源 IP 地址，即路径上的第二个路由器地址。类似步骤，tracert 逐跳获得每一个路由器的地址，并探测到目的设备的可达性。
![](https://notes-learning.oss-cn-beijing.aliyuncs.com/boov5o/1622087956576-de3f09ef-891d-4830-bef8-c84ef221aac4.png)tracert
tracert 过程也是双向的消息通信，只有双向都成功传输时，才能正确探测路径。另外主机安装了防火墙，也可能造成路径探测失败。

### 网络实战

#### ping

在 Windows 电脑上使用 **ping** 命令，并查看返回信息。
![](https://notes-learning.oss-cn-beijing.aliyuncs.com/boov5o/1622087956590-14d0d35c-c869-48f1-be86-6ffd13a9591e.png)ping 命令
同步**抓包**进行验证。
![](https://notes-learning.oss-cn-beijing.aliyuncs.com/boov5o/1622087956373-656b7067-8cd6-4d87-9fcf-f18863eae6ec.png)ping 抓包
还可以直接使用 ping 命令，查看 ping 命令的**使用方法**。
![](https://notes-learning.oss-cn-beijing.aliyuncs.com/boov5o/1622087956428-2fb4954b-c839-4523-b3a2-ff3565ba6188.png)ping 命令用法

#### tracert

在 Windows 电脑上使用 **tracert** 命令，并查看返回信息。
![](https://notes-learning.oss-cn-beijing.aliyuncs.com/boov5o/1622087956498-ede48849-269f-4fe6-883e-1ed11215e544.png)tracert 命令
同步**抓包**进行验证。
![](https://notes-learning.oss-cn-beijing.aliyuncs.com/boov5o/1622087956533-997964bf-42cb-4306-92bc-11b97ccd05e9.png)tracert 抓包
也可以直接使用 tracert 命令，查看 tracert 命令的**使用方法**。
![](https://notes-learning.oss-cn-beijing.aliyuncs.com/boov5o/1622087956443-84a5971e-0af1-4287-b740-d0160533e903.png)tracert 用法

---

**饮水思源：**
TCP/IP 详解 卷 1：协议 - Kevin R.Fall
网络基础 - 田果
图解 TCP/IP - 竹下隆史
路由交换技术 - 杭州华三通信技术有限公司

# ICMPv6

ICMPv6 是 IPv6 的基础协议之一，定义在 RFC 2463 中。用来传递报文转发中产生的信息或者错误。ICMPv6 定义的报文被广泛地应用在其他协议中，包括 NDP(邻居发现协议)、PathMTU 发现机制 等等。
