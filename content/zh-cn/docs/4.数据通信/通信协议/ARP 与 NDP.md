---
title: ARP 与 NDP
linkTitle: ARP 与 NDP
date: 2024-02-23T12:23
weight: 20
---

# 概述

> 参考：
>
> - [Wiki，ARP](https://en.wikipedia.org/wiki/Address_Resolution_Protocol)
> - [Wiki，NDP](https://en.wikipedia.org/wiki/Neighbor_Discovery_Protocol)
> - [RFC 826](https://www.rfc-editor.org/rfc/rfc826.html)
> - [公众号，36 张图详解 ARP](https://mp.weixin.qq.com/s/_5Wgsx4mEoDZgwv9-yHYEA)

**Address Resolution Protoco(地址解析协议，简称 ARP)** 是一种通信协议，该协议可以通过给定的网络层地址(通常是 IPv4 地址)，发现与之相关联的链路层地址(通常你是 MAC 地址)。ARP 于 1982 年在 [RFC 826](https://www.rfc-editor.org/rfc/rfc826.html) 中定义。说白了，就是根据 IP 地址查询对应 MAC 地址的协议。

注意：在 IPv6 网络环境下，APR 的功能已经被 NDP 替代

对应关系：一个 ip 地址对应一个 MAC 地址。多个 ip 地址可以对应一个 MAC 地址(e.g.一个网卡上配置两个 ip)

## ARP 报文

在抓包时，可以抓到如下几种 ARP 包

- ARP, Request who-has 10.10.100.254 tell 10.10.100.101, length 28
  - 在局域网中询问谁有 10.10.100.254，告诉自己，自己就是 10.10.100.101
- ARP, Reply 10.10.100.254 is-at 00:0f:e2:ff:05:92, length 46
  - 当 10.10.100.254 收到 arp 广播包之后，就会进行 reply(响应)，10.10.100.254 的 mac 地址是 00:0f:e2:ff:05:92
- 这时候 10.10.100.101 就会更新自己的 arp 表，记录下来 10.10.100.254 与 00:0f:e2:ff:05:92 的对应关系

ARP 报文分为 **ARP 请求报文**和 **ARP 应答报文**，它们的报文格式相同，但是各个字段的取值不同。

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/nguycm/1623911614198-6f1cde9d-cf37-4711-b6ce-cd9b8b1ccdd5.png)ARP 报文格式

ARP 报文中各个字段的含义如下。

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/nguycm/1623911614222-e8fe87d5-ca7b-494c-99c1-b8b1be472301.png)

# ARP 原理

ARP 是通过 **ARP Request(请求)**和 **ARP Reply(响应)**报文确定 MAC 地址的。
![](https://notes-learning.oss-cn-beijing.aliyuncs.com/nguycm/1623910582328-062e8f26-c6da-4ee5-97e4-857507dbb707.png)

## 同网段交互流程

```bash
root@desistdaydream:~# tcpdump -i any host 172.19.42.202
tcpdump: verbose output suppressed, use -v or -vv for full protocol decode
listening on any, link-type LINUX_SLL (Linux cooked v1), capture size 262144 bytes
14:13:38.036170 ARP, Request who-has 172.19.42.202 tell desistdaydream.bj-test, length 28
14:13:38.036549 ARP, Reply 172.19.42.202 is-at 00:be:d5:ef:24:4e (oui Unknown), length 46
14:13:38.036583 IP desistdaydream.bj-test > 172.19.42.202: ICMP echo request, id 3, seq 1, length 64
14:13:38.036821 IP 172.19.42.202 > desistdaydream.bj-test: ICMP echo reply, id 3, seq 1, length 64
```

如上所示，当前主机没有 172.19.42.202 主机的 arp 关系表，当我们向 172.19.42.202 发送 icmp 请求时，数据包将会经历如下过程：

假如主机 A 向同一网段上的主机 B 发送数据。主机 A 的 IP 地址为 `10.0.0.1` ，主机 B 的 IP 地址为 `10.0.0.2` ，主机 C 的 IP 地址为 `10.0.0.3` 。它们都不知道对方的 MAC 地址。ARP 地址解析过程如下：
![](https://notes-learning.oss-cn-beijing.aliyuncs.com/nguycm/1623910652117-bd47d4a3-1fb2-4f50-9b4f-bc8e5c2a1b2b.png)

1. **主机 A** 首先查看自己的 ARP 表（即 ARP 缓存表），确定是否有主机 B 的 IP 地址对应表项。如果有，则直接使用表项中的 MAC 地址进行封装，封装成帧后发送给主机 B 。

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/nguycm/1623910669257-9b3f2f80-ddd3-4501-807e-96acfe414c20.png)

2. 如果**主机 A** 的 ARP 表没有对应的表项，就发送一个广播帧，源 IP 和源 MAC 地址是主机 A ，目的 IP 地址是主机 B ，目的 MAC 地址是广播 MAC 地址，即 `FFFF-FFFF-FFFF` 。这就是 **ARP 请求报文**。

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/nguycm/1623910669279-c6354c1b-c098-4974-b4cd-22e93fc54efe.png)

3. ARP 请求是广播报文，同一个网段的所有主机都能收到。只有主机 B 发现报文中的目的 IP 地址是自己，于是**主机 B** 发送响应报文给主机 A ，源 MAC 地址和源 IP 地址是主机 B ，目的 MAC 地址和目的 IP 地址是主机 A ，这个报文就叫 **ARP 响应报文**。同时，主机 B 的 ARP 表记录主机 A 的映射关系，即主机 A 的 IP 地址和 MAC 地址的对应关系。

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/nguycm/1623910669203-16790059-e9ef-45f9-9c9b-28f5f835c556.png)

4. **主机 C** 也收到了 ARP 请求报文，但目的 IP 地址不是自己，所以不会进行响应。于是主机 C 添加主机 A 的映射关系到 ARP 表，并丢弃 ARP 请求报文。

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/nguycm/1623910669305-20a3b990-6012-45d5-8462-176106ad6dd6.png)

5. 主机 A 收到 ARP 响应报文后，添加主机 B 的映射关系，同时用主机 B 的 MAC 地址做为目的地址封装成帧，并发送给主机 B 。

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/nguycm/1623910669300-be2dccaa-122c-41d2-92b6-3b035e70335d.png)
如果每发送一个 IP 报文就要进行一次 ARP 请求，来确定 MAC 地址，那将会造成不必要的网络流量，通常的做法是用 ARP 表记录 IP 地址和 MAC 地址的映射关系。主机发送报文时，首先会查看它的 **ARP 表**，目的是为了确定是否是已知的设备 MAC 地址。如果有，就直接使用；如果没有，就发起 ARP 请求获取。不过，缓存是有一定期限的。ARP 表项在**老化时间**（ `aging time` ）内是有效的，如果老化时间内未被使用，表项就会被删除。
![](https://notes-learning.oss-cn-beijing.aliyuncs.com/nguycm/1623910669502-4614b748-fac7-4301-bc85-68ccec45c91a.png)
ARP 表项分为**动态 ARP 表项**和**静态 ARP 表项**：

- **动态 ARP 表项**由 ARP 动态获取，因此在网络通信中，无需事先知道 MAC 地址，只要有 IP 地址即可。如果老化时间内未被使用，表项就会被自动删除。
- **静态 ARP 表项**是手工配置，不会老化。静态 ARP 表项的优先级高于动态 ARP 表项，可以将相应的动态 ARP 表项覆盖。

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/nguycm/1623910669307-374038c8-0ab6-4a06-ad5a-b55240326db1.png)

## 跨网段交互流程

1. 在 client 向 server 发送数据时，内核会在该数据包外封装 源 IP、源 MAC、目的 IP，由于目的 MAC 未知，所以不填
2. 数据包经过交换、路由等设备后，到达 server 前的交换机，交换机根据自身的 arp 表，找到目的 IP 与哪个口的 mac 地址相关联，则把数据包交给响应的网口，以便顺利到达 server。
3. server 收到数据包后，内核会把最外层的 IP 与 MAC 都玻璃，并交给相对应的用户空间内的程序，处理完成后，再发送数据响应 client
4. 此时 server 已经有个 client 的 MAC 与 IP，所以在内核封装的时候，会填写源 IP、源 MAC、目的 IP、目的 MAC。
5. 之后数据包到达 client 前面的交换机时，交换机发现目的 IP 与 MAC，则更新或者保持不变 arp 表，并把数据包交给 MAC 为目的 MAC 的端口，以便数据送达 client
6. 这样就完成了两台设备之间的互相通信以及数据报文的完整封装
7. 否则第一次发送数据的时候，如果是不在本网段的地址，则目的 MAC 一般都是未知的，除非已经建立连接之后，才能根据数据包的源 MAC 知道响应的时候目的 MAC 填什么

# 免费 ARP 广播 与 一般 ARP 广播

- **免费 arp 广播** # 在设备首次接通网线之后，会进行 arp 广播，告诉大家自己的 IP 与 MAC 对应关系。当局域网内的设备收到这个免费的 arp 广播时，没有该 arp 记录则会添加，如果该 arp 记录改变了则会更新。
- **一般 arp 广播** # 在对本机未知同网段的设备发送数据时，会进行 arp 广播，询问该设备在哪里。
- “免费 arp 广播”与“一般 arp 广播”的区别：
  1. 普通的 arp 请求如果目的地址不是本机，则本机就直接丢弃了，但是免费的 arp 广播，则会在本机保留或者更新。

如果目的 IP 不在本网段，则不会进行 arp 广播，因为目的地址不在本网段的设备，发送出去直接就到网关了，而网关在接入的时候、每隔 N 时间，都会发送 arp 来询问网关在哪。至于交换机测，也是同样的道理，本网段的设备，在需要发送数据包的时候，如果 arp 表里没有，则会先进行 arp 广播再发送。因为这些都是通过 mac 地址来进行二层发送的。

**免费 ARP** 是一种特殊的 ARP 请求，它并非通过 IP 找到对应的 MAC 地址，而是当主机启动的时候，发送一个免费 ARP 请求，即请求自己的 IP 地址的 MAC 地址。
![](https://notes-learning.oss-cn-beijing.aliyuncs.com/nguycm/1623911576682-1a99a60a-4fa9-4480-8061-cec725fe9073.webp)
与普通 ARP 请求报文的区别在于报文中的目标 IP 地址。普通 ARP 报文中的目标 IP 地址是其它主机的 IP 地址；而免费 ARP 的请求报文中，**目标 IP 地址是自己的 IP 地址**。
![](https://notes-learning.oss-cn-beijing.aliyuncs.com/nguycm/1623911577030-243f8668-60e2-4952-bac9-beef00f9cba4.png)
免费 ARP 的作用：

- 起到一个**宣告**作用。它以广播的形式将数据包发送出去，不需要得到回应，只为了告诉其它主机自己的 IP 地址和 MAC 地址。
- 可用于**检测 IP 地址冲突**。当一台主机发送了免费 ARP 请求报文后，如果收到了 ARP 响应报文，则说明网络内已经存在使用该 IP 地址的主机。
- 可用于**更新**其它主机的 ARP 缓存表。如果该主机更换了网卡，而其它主机的 ARP 缓存表仍然保留着原来的 MAC 地址。这时，通过免费的 ARP 数据包，更新其它主机的 ARP 缓存表。

# NDP
