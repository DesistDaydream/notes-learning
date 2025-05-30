---
title: QUIC
linkTitle: QUIC
weight: 20
---


# 概述

> 参考：
>
> - [RFC 9000 - QUIC: A UDP-Based Multiplexed and Secure Transport](https://datatracker.ietf.org/doc/html/rfc9000)
> - [Wiki, QUIC](https://en.wikipedia.org/wiki/QUIC)

尽管其名称最初被提议为“Quick UDP Internet Connections”的缩写，但在 IETF 中使用的 QUIC 一词并不是缩写词；它只是协议的名称。QUIC 提高了当前使用传输控制协议 (TCP) 的面向连接的 Web 应用程序的性能。它通过使用 [UDP](/docs/4.数据通信/Protocol/UDP/UDP.md) 在两个端点之间建立大量多路复用连接来实现这一点，并且旨在为许多应用程序在传输层淘汰 TCP，从而为该协议赢得了偶尔的绰号“TCP/2”

# 如何用 UDP 实现可靠传输？

> 参考：
>
> - 原文：[支付宝二面：如何用 UDP 实现可靠传输？](https://mp.weixin.qq.com/s/Fwxfzzu6QdqyiYUI38afIQ)

> **自 2015 年以来，QUIC 协议开始在 IETF 进行标准化并被国内外各大厂商相继落地。鉴于 QUIC 具备“0RTT 建联”、“支持连接迁移”等诸多优势，并将成为下一代互联网协议：HTTP3.0 的底层传输协议，蚂蚁集团支付宝客户端团队与接入网关团队于 2018 年下半年开始在移动支付、海外加速等场景落地 QUIC。**

本文是综述篇，介绍**QUIC** 在蚂蚁的整体落地情况。

之所以是综述，是因为 QUIC 协议过于复杂，如果对标已有的协议，QUIC 近似等于 HTTP + TLS +TCP，无法详细的毕其功于一役，因此我们通过综述的方式将落地的重点呈现给读者，主要介绍如下几个部分：

- QUIC背景：简单全面的介绍下 QUIC 相关的背景知识
- 方案选型设计：详细介绍蚂蚁的落地方案如何另辟蹊径、优雅的支撑 QUIC 的诸多特性，包括连接迁移等
- 落地场景：介绍 QUIC 在蚂蚁的两个落地场景，包括：支付宝客户端链路以及海外加速链路
- 几项关键技术：介绍落地 QUIC 过程中核心需要解决的问题，以及我们使用的方案，包括：“支持连接迁移”、“提升 0RTT 比例"， "支持 UDP 无损升级”以及“客户端智能选路” 等
- 几项关键的技术专利

本文也是 QUIC 协议介绍的第一篇，后续我们会把更多的落地细节、体验优化手段、性能优化手段、安全与高可用、QUIC 新技术等呈现给大家。

【注】蚂蚁 QUIC 开发团队包括：支付宝客户端团队的梅男、苍茫、述言，以及接入网关的伯琴、子荃、毅丝。

**QUIC 背景介绍**

鉴于读者的背景可能不同，在开始本文之前，我们先简单介绍下 QUIC 相关的背景知识，如果您对这个协议的更多设计细节感兴趣，可以参见相关 Draft：https://datatracker.ietf.org/wg/quic/documents/

**一、QUIC 是什么？**
---------------

简单来说，QUIC (Quick UDP Internet Connections) 是一种基于 UDP 封装的安全  可靠传输协议，他的目标是取代 TCP 并自包含 TLS 成为标准的安全传输协议。下图是 QUIC 在协议栈中的位置，基于 QUIC 承载的 HTTP 协议进一步被标准化为 HTTP3.0。

![](https://mmbiz.qpic.cn/mmbiz_png/nibOZpaQKw0ibaxJGWmOu0IGdj6Gvuclibq2HuO8A1X400lJC27JaVr7cu65mibXqXI7OicJgqvlBlW3f6BHgDLiaJIQ/640?wx_fmt=png)

**二、为什么是 QUIC ？**
-----------------

在 QUIC 出现之前，TCP 承载了 90% 多的互联网流量，似乎也没什么问题，那又为何会出现革命者 QUIC 呢？这主要是因为发展了几十年的 TCP 面临 “协议僵化问题”，表现在几方面：

1. 网络设备支持 TCP 时的僵化，表现在：对于一些防火墙或者 NAT 等设备，如果 TCP 引入了新的特性，比如增加了某些 TCP OPTION 等，可能会被认为是攻击而丢包，导致新特性在老的网络设备上无法工作。
2. 网络操作系统升级困难导致的 TCP 僵化，一些 TCP 的特性无法快速的被演进。
3. 除此之外，当应用层协议优化到 TLS1.3、 HTTP2.0 后， 传输层的优化也提上了议程，QUIC 在 TCP 基础上，取其精华去其糟粕具有如下的硬核优势：

![](https://mmbiz.qpic.cn/mmbiz_png/nibOZpaQKw0ic3rhzShkGTuUia9kN5cHaexPeic7AJGCUvZ2N7sIBb4BWKbRJibRibcdpxhnuaMibSOMCRflaCtM4o40Q/640?wx_fmt=png)

**三、QUIC 生态圈发展简史**
------------------

下图是 QUIC 从创建到现在为止的一些比较重要的时间节点，2021 年，QUIC V1 即将成为 RFC，结束百花齐放的态势。

![](https://mmbiz.qpic.cn/mmbiz_png/nibOZpaQKw0ibaxJGWmOu0IGdj6GvuclibqDRIjeXnIaWQicheKQwUniccFe0N9cuRBCBEGSN9hk4xI8cibwu4CfKCibg/640?wx_fmt=png)

介绍完 QUIC 相关背景，之后我们来介绍蚂蚁的整个落地的内容，这里为了便于阐述，我们用蚂蚁 QUIC 的 一、二、三、四 来进行概括总结，即 “一套落地框架”、“两个落地场景”、“三篇创新专利保护”、“四项关键技术”。

**一套落地框架**

蚂蚁的接入网关是基于多进程的 NGINX 开发的 (内部称为 Spanner，协议卸载的扳手)，而 UDP 在多进程编程模型上存在诸多挑战，典型的像无损升级等。为了设计一套完备的框架，我们在落地前充分考虑了服务端在云上部署上的方便性、扩展性、以及性能问题，设计了如下的落地框架以支撑不同的落地场景：

![](https://mmbiz.qpic.cn/mmbiz_png/nibOZpaQKw0ic3rhzShkGTuUia9kN5cHaex7qQlr5kKXqsgOP2AZdtmtRDDQ58Ho1CiaG8KuTDZn7BYfLd2aKMibuVA/640?wx_fmt=png)

在这套框架中，包括如下两个组件：

- QUIC LB 组件：基于 NGINX 4层 UDP Stream 模块开发，用来基于 QUIC DCID 中携带的服务端信息进行路由，以支持连接迁移。
- NGINX QUIC 服务器：开发了 NGINX\_QUIC\_MODULE，每个 Worker 监听两种类型的端口：
  - BASE PORT ，每个 Worker 使用的相同的端口号，以 Reuseport 的形式监听，并暴露给 QUIC LB，用以接收客户端过来的第一个 RTT 中的数据包，这类包的特点是 DCID 由客户端生成，没有路由信息。
  - Working PORT，每个 Worker 使用的不同的端口号，是真正的工作端口，用以接收第一个 RTT 之后的 QUIC 包，这类包的特定是 DCID 由服务端的进程生成携带有服务端的信息。

当前框架支持的能力包括如下：

1. 在不用修改内核的情况下，完全在用户态支持 QUIC 的连接迁移，以及连接迁移时 CID 的 Update
2. 在不用修改内核的情况下，完全在用户态支持 QUIC 的无损升级以及其他运维问题
3. 支持真正意义上的 0RTT ，并可提升 0RTT 的比例

为何能支持上述能力，我们后面会展开叙述

**两个落地场景**

我们由近及远的两个落地场景如下：

**场景一、支付宝移动端落地**

如下为我们落地架构的示意图，支付宝手机客户端通过 QUIC 携带 HTTP 请求，通过 QUIC LB 等四层网关将请求转发到 Spanner （蚂蚁内部基于 NGINX 开发的 7 层网关），在 Spanner 上我们将 QUIC 请求 Proxy 成 TCP 请求，发送给业务网关（RS）。

![](https://mmbiz.qpic.cn/mmbiz_png/nibOZpaQKw0ibaxJGWmOu0IGdj6GvuclibqZNa2k7XicdX02n90puGXSG3WfJ3j0K30UHpiaibEWlMmIO4aw03Yt5W6Q/640?wx_fmt=png)

具体的方案选型如下：

- 支持的 QUIC 版本是 gQUIC Q46。
- NGINX QUIC MODULE 支持 QUIC 的接入和 PROXY 成 TCP 的能力。
- 支持包括移动支付、基金、蚂蚁森林在内的所有的 RPC 请求。
- 当前选择 QUIC 链路的方式有两种 ：

1. Backup 模式，即在 TCP 链路无法使用的情况下，降级到 QUIC 链路。
2. Smart 模式，即 TCP和 QUIC 竞速，在 TCP 表现力弱于 QUIC 的情况下，下次请求主动使用 QUIC 链路。

在此场景下，通过使用 QUIC 可以获得的红利包括：

1. 在客户端连接发生迁移的时候，可以不断链继续服务
2. 客户端在首次发起连接时，可以节省 TCP 三次握手的时间
3. 对于弱网情况，QUIC 的传输控制可以带来传输性能提升

**场景二、**海外加速落地

蚂蚁集团从 2018 年开始自研了海外的动态加速平台 AGNA（Ant Global Network Accelerator）以替换第三方厂商的加速服务。AGNA 通过在海外部署接入点：Local Proxy(LP) 以及在国内部署接出点：Remote Proxy （RP）的方式，将用户的海外请求通过 LP 和 RP 的加速链路回源国内。如下图所示，我们将 QUIC 部署在 LP 和 RP 之间的链路。

![](https://mmbiz.qpic.cn/mmbiz_png/nibOZpaQKw0ibaxJGWmOu0IGdj6Gvuclibq4hYMIaIYvOZZQS7Jmfw1QoKibUCYBHSpplRicsS0bBneQ9d7dqnicEDtg/640?wx_fmt=png)

在海外接入点上(LP)，每一个 TCP 连接都被 Proxy 成 QUIC 上的一个 Stream 进行承载，在国内接出点上(RP), 每一个 QUIC Stream 又被 Proxy 成一个 TCP 连接，LP 和 RP 之间使用 QUIC 长连接。

在此场景下，通过使用 QUIC 可以获得的红利包括：

1. 通过 QUIC 长连接的上的 Stream 承载 TCP 请求，避免每次的跨海建联。
2. 对于跨海的网络，QUIC 的传输控制可以带来传输性能提升。

**三篇关键专利**

到目前为止，我们把落地过程中一些创新的技术点通过申请专利进行了保护，并积极在 IETF 进行标准化分享我们的研究成果，包括：

### **专利一**

我们将落地场景 2 中，通过 QUIC Stream 进行四层代理的手段来进行海外回源的加速方法进行专利保护，提出：“一种基于 QUIC 协议代理的的链路加速方法”，目前此专利已经获得美国专利授权，专利号：CN110213241A。

### **专利二**

将我们落地框架中的 QUIC LB 组件作为专利进行保护，提出：“一种无状态、一致性、分布式的 QUIC 负载均衡设备”，目前此专利还在受理中。由于通过 QUIC LB 可以很好的支持 QUIC 协议的连接迁移问题，所以目前 IETF QUIC WG 上有关于 QUIC LB 相关的草案，我们目前已经参与到 Draft 的讨论和制定中，后序相关的方案也会继续推广到云上产品。

### **专利三**

将我们解决的 UDP 的无损升级方法进行专利保护，提出 “一种 QUIC 服务器无损升级方案”，目前此专利还在受理中。由于 UDP 无损升级问题是一个业界难题，当前有些手段需要在用户态进行跳转，性能损失较大，我们的方案可以在我们的落地框架中解决当前问题，关于这个方案的细节我们再后面的关键技术中进行介绍。

**四项关键技术**

在整个的落地中，我们设计的方案围绕解决几个核心问题展开，形成了四项关键技术，分别如下

#### **技术点1.优雅的支持连接迁移能力**

先说 **连接迁移面临的问题 ，**上文有提到，QUIC 有一项比较重要的功能是支持连接迁移。这里的连接迁移是指：如果客户端在长连接保持的情况下切换网络，比如从 4G 切换到 Wifi , 或者因为 NAT Rebinding 导致五元组发生变化，QUIC 依然可以在新的五元组上继续进行连接状态。QUIC 之所以能支持连接迁移，一个原因是 QUIC 底层是基于无连接的 UDP，另一个重要原因是因为 QUIC 使用唯一的 CID 来标识一个连接，而不是五元组。

如下图所示，是 QUIC 支持连接的一个示意图，当客户端出口地址从 A 切换成 B 的时候，因为 CID 保持不变，所以在 QUIC 服务器上，依然可以查询到对应的 Session 状态。

![](https://mmbiz.qpic.cn/mmbiz_png/nibOZpaQKw0ibaxJGWmOu0IGdj6Gvuclibq1IMXmRmiakLWriap5ee9KApe3Rlwh8ERU6CByllwkQmeibF7zhLIpdTqA/640?wx_fmt=png)

然而，理论很丰满，落地却很艰难，在端到端的落地过程中，因为引入了负载均衡设备，会导致在连接迁移时，所有依赖五元组 Hash 做转发或者关联 Session 的机制失效。以 LVS 为例，连接迁移后， LVS 依靠五元组寻址会导致寻址的服务器存在不一致。即便 LVS 寻址正确，当报文到达服务器时，内核根据五元组关联进程，依然会寻址出错。同时，IETF Draft 要求，连接迁移时 CID 需要更新掉，这就为仅依靠 CID 来转发的计划同样变的不可行。

再说 **我们的解决方法，**为了解决此问题，我们设计了开篇介绍的落地框架，这里我们将方案做一些简化和抽象，整体思路如下图所示：

1. 在四层负载均衡上，我们设计了 QUIC LoadBalancer 的机制：
1. 我们在 QUIC 的 CID 中扩展了一些字段（ServerInfo）用来关联 QUIC Server 的 IP 和 Working Port 信息。
2. 在发生连接迁移的时候，QUIC LoadBalancer 可以依赖 CID 中的 ServerInfo 进行路由，避免依赖五元组关联 Session 导致的问题。
3. 在 CID 需要 Update 的时候，NewCID 中的 ServerInfo 保留不变，这样就避免在 CID 发生 Update 时，仅依赖 CID Hash 挑选后端导致的寻址不一致问题。
2. 在 QUIC 服务器多进程工作模式上，我们突破了 NGINX 固有的多 Worker 监听在相同端口上的桎梏，设计了多端口监听的机制，每个 Worker 在工作端口上进行隔离，并将端口的信息携带在对 First Initial Packet 的回包的 CID 中，这样代理的好处是：
1. 无论是否连接迁移，QUIC LB 都可以根据 ServerInfo，将报文转发到正确的进程。
2. 而业界普遍的方案是修改内核，将 Reuse port 机制改为 Reuse CID 机制，即内核根据 CID 挑选进程。即便后面可以通过 ebpf 等手段支持，但我们认为这种修改内核的机制对底层过于依赖，不利于方案的大规模部署和运维，尤其在公有云上。
3. 使用独立端口，也有利于多进程模式下，UDP 无损升级问题的解决，这个我们在技术点 3 中介绍。

![](https://mmbiz.qpic.cn/mmbiz_png/nibOZpaQKw0ibaxJGWmOu0IGdj6Gvuclibq4OhaNIfuKQA210xzzgzcWvDQibvXNfyicIElQR3qAZfibOichRiaSll4ibyA/640?wx_fmt=png)

#### **技术点2.提升 0RTT 握手比例**

这里先 **介绍 QUIC 0RTT 原理**。前文我们介绍过， QUIC 支持传输层握手和安全加密层握手都在一个 0RTT 内完成。TLS1.3 本身就支持加密层握手的 0RTT，所以不足为奇。而 QUIC 如何实现传输层握手支持 0RTT 呢？我们先看下传输层握手的目的，即：服务端校验客户端是真正想握手的客户端，地址不存在欺骗，从而避免伪造源地址攻击。在 TCP 中，服务端依赖三次握手的最后一个 ACK 来校验客户端是真正的客户端，即只有真正的客户端才会收到 Sever 的 syn\_ack 并回复。

![](https://mmbiz.qpic.cn/mmbiz_png/nibOZpaQKw0ibaxJGWmOu0IGdj6Gvuclibq01dot1oWD2STkCTWoxx99WTa39ohvY7d4uGlBfm9POiaC98iah960Z5w/640?wx_fmt=png)

QUIC 同样需要对握手的源地址做校验，否则便会存在 UDP 本身的 DDOS 问题，那 QUIC 是如何实现的？依赖 STK(Source Address Token) 机制。这里我们先声明下，跟 TLS 类似，QUIC 的 0RTT 握手，是建立在已经同一个服务器建立过连接的基础上，所以如果是纯的第一次连接，仍然需要一个 RTT 来获取这个 STK。如下图所示，我们介绍下这个原理：

1. 类似于 Session Ticket 原理，Server 会将客户端的地址和当前的 Timestamp 通过自己的 KEY 加密生成 STK。
2. Client 下次握手的时候，将 STK 携带过来，由于 STK 无法篡改，所以 Server 通过自己的 KEY 解密，如果解出来的地址和客户端此次握手的地址一致，且时间在有效期内，则表示客户端可信，便可以建立连接。
3. 由于客户端第一次握手的时候，没有这个 STK，所以服务度会回复 REJ 这次握手的信息，并携带 STK。

![](https://mmbiz.qpic.cn/mmbiz_png/nibOZpaQKw0ibaxJGWmOu0IGdj6Gvuclibq67W9iagJd9YUbGnFdvhTm53u0c3diaY2kKkRyCW3BO8rT8TyR3ctHELg/640?wx_fmt=png)

理论上说，只要客户端缓存了这个 STK，下次握手的时候带过来，服务端便可以直接校验通过，即实现传输层的 0RTT。但是真实的场景却存在**如下两个问题**：

1. 因为 STK 是服务端加密的，所以如果下次这个客户端路由到别的服务器上了，则这个服务器也需要可以识别出来。
2. STK 中 encode 的是上一次客户端的地址，如果下一次客户端携带的地址发生了变化，则同样会导致校验失败。此现象在移动端发生的概率非常大，尤其是 IPV6 场景下，客户端的出口地址会经常发生变化。

再介绍下**我们的解决方法。** 第一个问题比较好解，我们只要保证集群内的机器生成 STK 的秘钥一致即可。**第二个问题，我们的解题思路是**：

1. 我们在 STK 中扩展了一个 Client ID, 这个 Clinet ID 是客户端通过无线保镖黑盒生成并全局唯一不变的，类似于一个设备的 SIMID，客户端通过加密的 Trasnport Parameter 传递给服务端，服务端在 STK 中包含这个 ID。
2. 如果因为 Client IP 发生变化导致校验 STK 校验失败，便会去校验 Client ID，因为 ID 对于一个 Client 是永远不变的，所以可以校验成功，当然前提是，这个客户端是真实的。为了防止 Client ID 的泄露等，我们会选择性对 Client ID 校验能力做限流保护。

![](https://mmbiz.qpic.cn/mmbiz_png/nibOZpaQKw0ibaxJGWmOu0IGdj6Gvuclibqx2BXQFWnG86ZQ6VSrNGEb07U43lwW6zDOmKqkAzc6EZu3uw2WEXbFw/640?wx_fmt=png)

#### **技术点3. 支持 QUIC 无损升级**

我们知道 **UDP 无损升级是业界难题。** 无损升级是指在 reload 或者更新二进制时，老的进程可以处理完存量连接上的数据后优雅退出。以 NGINX 为例，这里先介绍下 TCP 是如何处理无损升级的，主要是如下的两个步骤：

1. 老进程先关闭 listening socket，待存量连接请求都结束后，再关闭连接套接字
2. 新进程从老进程继承 listening socket , 开始 accept 新的请求

而 UDP 无法做到无损升级是因为 UDP 只有一个 listening socket 没有类似 TCP 的连接套接字，所有的收发数据包都在这个 socket 上，导致下面的热升级步骤会存在问题：

1. 在热升级的时候，old process fork 出 new process 后，new process 会继承 listening socket 并开始 recv msg。
2. 而 old process 此时如果关闭 listenging socket， 则在途的数据包便无法接收，达不到优雅退出的目的。
3. 而如果继续监听，则新老进程都会同时收取新连接上的报文，导致老进程无法退出。

![](https://mmbiz.qpic.cn/mmbiz_png/nibOZpaQKw0ibaxJGWmOu0IGdj6GvuclibqyD0ibhw7FU4kqOwMYicooDKJaK3kN8zTTLyJa5TE5b3WZRyYbpX59oUw/640?wx_fmt=png)

这里介绍下**相关的解决方法。** 针对此问题，业界有一些方法，比如：在数据包中携带进程号，当数据包收发错乱后，在新老进程之间做一次转发。考虑到接入层上的性能等原因，我们不希望数据再做一次跳转。结合我们的落地架构，我们设计了如下的 **基于多端口轮转的无损升级方案**，简单来说，我们让新老进程监听在不同的端口组并携带在 CID 中，这样 QUIC LB 就可以根据端口转发到新老进程。为了便于运维，我们采用端口轮转的方式，新老进程会在 reload N 次之后，重新开始之前选中的端口。如下图所示：

1. 无损升级期间，老进程的 Baseport 端口关闭，这样不会再接受 first intial packet, 类似于关闭了 tcp 的 listening socket。
2. 老进程的工作端口，继续工作，用来接收当前进程上残余的流量。
3. 新进程的 Baseport 开始工作，用来接收 first initial packet, 开启新的连接，类似于开启了 tcp 的 listening socket。
4. 新进程的 working port = (I + 1) mod N,  N 是指同时支持新老进程的状态的次数，例如 N = 4， 表示可以同时 reload 四次，四种 Old, New1, New2, New3 四种状态同时并存，I 是上一个进程工作的端口号，这里 + 1 是因为只有一个 worker, 如果 worker 数有 M 个，则加 M。
5. 建好的连接便被 Load Balancer 转移到新进程的监听端口的 Working Port 上。

![](https://mmbiz.qpic.cn/mmbiz_png/nibOZpaQKw0ibaxJGWmOu0IGdj6GvuclibqBn5qV0fOfRrM42KPEqX4ZcDQh0gTtD0I17QvkAApCKEgdAHbHfVGiaA/640?wx_fmt=png)

**技术点4.客户端智能选路**

尽管落地 QUIC 的愿望是好的，但是新事物的发展并不是一帆风顺的。由于 QUIC 是基于 UDP 的，而 UDP 相比于 TCP 在运营商的支持上并非友好，表现在：

1. 在带宽紧张的时候，UDP 会经常被限流。
2. 一些防火墙对于 UDP 包会直接 Drop。
3. NAT 网关针对 UDP 的 Session 存活时间也较短。

同时，根据观察发现，不同的手机厂商对于 UDP 的支持能力也不同，所以在落地过程中，如果盲目的将所有流量完全切为 QUIC 可能会导致一些难以预料的结果。为此，我们在客户端上，设计了开篇介绍的 TCP 和 QUIC 相互 Backup 的链路，如下图所示，我们实时探测 TCP 链路和 QUIC 链路的 RTT、丢包率、请求完成时间、错误率等指标情况，并根据一定的量化方法对两种链路进行打分，根据评分高低，决定选择走哪种链路，从而避免寻址只走一条链路导致的问题。

![](https://mmbiz.qpic.cn/mmbiz_png/nibOZpaQKw0ibaxJGWmOu0IGdj6GvuclibqLLpBUvS8DfFvFZScC3Uz388OlDffAVZQ0d2vnt3OHd2MFpRic3NMiazQ/640?wx_fmt=png)

**做个总结**

本文主要综述性的介绍了 QUIC 在蚂蚁的落地方案、场景以及一些关键技术。**关键技术上**，主要介绍了我们如何通过创造性的提出 QUIC LB 组件、以及多端口监听的机制来优雅的支持 QUIC 的连接迁移机制、QUIC 服务端的无损升级等，依赖这套方案我们的接入网关不需要像业界一样依赖底层内核的改动，这极大的方便了我们方案的部署，尤其在公有云场景下。除了连接迁移以外，我们还提出了 0RTT 建联提升方案、客户端智能选路方案，以最大化 QUIC 在移动端上的收益。截止当前，QUIC 已经在支付宝移动端以及全球加速链路**两个场景**上平稳运行，并带来了较好的业务收益。

**未来规划**

两年来，我们主要以社区的 gQuic 为基础，充分发挥 QUIC 的协议优势，并结合蚂蚁的业务特征以最大化移动端收益为目标，创造性的提出了一些解决方案，并积极向社区和 IETF 进行推广。在未来，随着蚂蚁在更多业务上的开展和探索以及HTTP3.0/QUIC 即将成为标准，我们会主要围绕以下几个方向继续深挖 QUIC 的价值：

1. 我们将利用 QUIC 在应用层实现的优势，设计一套统一的具备自适应业务类型和网络类型的 QUIC 传输控制框架，对不同类型的业务和网络类型，做传输上的调优，以优化业务的网络传输体验。
2. 将 gQUIC 切换成 IETF QUIC，推进标准的 HTTP3.0 在蚂蚁的进一步落地。
3. 将蚂蚁的 QUIC LB 技术点向 IETF QUIC LB 进行推进，并最终演变为标准的 QUIC LB。
4. 探索并落地 MPQUIC（多路径 QUIC） 技术，最大化在移动端的收益。
5. 继续 QUIC 的性能优化工作，使用 UDP GSO， eBPF，io\_uring 等内核技术。
6. 探索 QUIC 在内网承载东西向流量的机会。
