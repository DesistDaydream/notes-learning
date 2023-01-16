---
title: IPv6
---

# 概述

> 参考：
> - [Wiki,IPv6](https://en.wikipedia.org/wiki/IPv6)
> - [RFC,8200](https://datatracker.ietf.org/doc/html/rfc8200)
> - [IANA,IPv6 地址空间分配情况](https://www.iana.org/assignments/ipv6-address-space/ipv6-address-space.xhtml)
> - [GitHub 项目，IPv6-CN/ipv6cn.github.io](https://github.com/IPv6-CN/ipv6-cn.github.io)(这个资料。。。怎么说呢。。。)
> - <https://www.bilibili.com/video/BV1J7411c7ae>
> - <https://www.bilibili.com/video/BV1aF411v7hU>
> - <https://www.zhihu.com/question/361275417>
> - <https://cloud.tencent.com/developer/article/1468099>

**Internet Protocol version 6(互联网协议版本 6，简称 IPv6)** 是 Internet Protocol(互联网协议，简称 IP) 的最新版本(截止 2022 年 1 月 28 日)，该协议为网络上的计算机提供识别和定位系统并通过 Internet 路由流量。IPv6 由 Internet Engineering Task Force(互联网工程任务组，简称 IETF) 开发，旨在解决 IPv4 地址耗尽的问题。1998 年 12 月，IPv6 成为 IETF 的标准草案，2017 年 7 月 14 日将其批准为互联网标准。
[
](https://www.bilibili.com/video/BV1J7411c7ae)

## 困扰 1. IPv4 和 IPv6 只有地址格式不同吗

除了地址格式不同，IPv4 与 IPv6 协议栈也不同，他们在逻辑上是**完全不同的 2 个世界**
以下实践中经常会遇到的 4 个不同之处：
▷ **基本通讯过程**：ND 替代 ARP、多播替代广播、fe80 地址成为标配、ICMP 成为通讯核心
▷ **IP 配置方式**：客户端以无状态自动配置 IP 成为主流，弱化 DHCP
▷ **DNS**[域名解析](https://cloud.tencent.com/product/cns?from=10680)：AAAA 记录替代 IPv4 的 A 记录、对应用存在优先级问题（优先解析 AAAA 还是 A）
▷ **应用层适应性**：socket 编程中 AF_INET 仅支持 IPv4，AF_INET6 仅支持 IPv6

---

## 困扰 2. IPv4 到 IPv6 对应用程序是透明无感知的吗

**错，是有感知的**，上层应用程序需要进行改造。
举个例子，当访问 fzxiaomange.com 时候，是要优先解析 IPv6 地址（AAAA）还是 IPv4 地址（A），因为总得选一条线路来发送请求。现在许多框架会优先选择 IPv6。
注意：如果解析出 AAAA 记录，即使本机没有可路由 IPv6 地址，也有可能**依然尝试通过 IPv6 进行请求，导致请求失败**。
还有一个典型的例子，是程序会在**应用层里交互底层 IP 地址**，比如 FTP 主动模式会在应用层里交互 IPv4 地址，而如果实际可用的是 IPv6 地址，就可能导致后续连接的异常。
**无法做到透明无感知，是导致产生 IPv4 到 IPv6 的部分过渡方案的原因之一。**

---

## 困扰 3. 提供 WEB 服务，需要每台服务器都配置 IPv6 地址吗

现在有一种言论，说“IPv6 地址无限多，每台服务器都可以配一个 IP 地址，不用做 NAT”。
**这很容易误导人，具体如何使用 IPv6，还得根据场景而定**。比如笔者的个人博客 fzxiaomange.com，由 nginx->php->mysql 组成，分别位于 3 台服务器上，那只需要在 nginx 上配置 IPv6 地址，并在 DNS 上添加一条 AAAA 记录指向 L7 的 IPv6 地址即可。完全没必要在 php、mysql 服务器上配置 IPv6 地址，而且一旦配置了，就直接暴露内网了。
每个设备都配置 IPv6，主要适用于**偏客户端以及地址需求量大的场景**，诸如物联网设备、手机 4G、家庭宽带等。
另外，**IPv6 有 NAT**，适用于办公 PC、机房服务器等需要访问 IPv6 网络，而不想被别人主动访问的场景。

---

## 困扰 4. IPv4 和 IPv6 要配在同一张网卡上吗

都可以，首先需要先了解 2 个词**“Single-Stack(单栈)”**和**“Dual-Stack(双栈)”**：
**以节点为角度**（通用的解释）：

- 单栈：表示一个 IPv6 节点，也就是一台服务器，或一部手机，仅有 IPv6 地址，或仅有 IPv4 地址，前者叫做**“IPv6 单栈”**或**“IPv6-Only”**，后者叫**“IPv4 单栈”**或**“IPv4-Only”**。
- 双栈：表示一个 IPv6 节点，同时拥有 IPv6 地址和 IPv4 地址

**以网卡为角度**：
▷ 单栈：表示一张网卡仅有 IPv6 地址，或仅有 IPv4 地址，示意图如下
![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/mzn9uc/1643343522816-f045574f-93d4-4ed9-9356-4045d7acaa38.png)
▷ 双栈：表示一张网卡同时拥有 IPv6 地址和 IPv4 地址
![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/mzn9uc/1643343532878-4cb46010-1684-4e57-88f5-29e2472d56d2.png)
IPv4 和 IPv6 在逻辑上是两个完全不相交的世界。如果终端处位于同一个物理层，比如同一个 VLAN，那么网卡就只能同时配置 IPv6 地址和 IPv4 地址；反之，就必须一张网卡配置 IPv6，另一张网卡配置 IPv4。所以，**关键看网络架构是如何设计，各有利弊**。比如放同一张网卡上，就可以做到带宽共享，而放不同网卡，可以做到带宽分别限制与计费。

---

## 困扰 5. 我的网卡有 fe80 开头的地址，可以用来上公网吗

当网卡启动的时候，会自动生成**“链路本地地址”（Link-Local Address）**，这是一个**fe80::/10**的单播地址。“链路本地地址”用于**IP 自动配置**、**邻居发现**等。
**注意事项**：
▷ 核心：每张网卡都会存在“链路本地地址”，这是 IPv6 协议通讯的核心，不应当删掉
▷ 范围：仅在同一个二层范围内进行传播，不会被路由器做转发
▷ 地址：“链路本地地址”的算法不统一，有的操作系统会根据 mac 地址计算而来（EUI-64），而有的则是随机或其他某种算法计算而来
▷ 服务：“链路本地地址”虽然可以在二层内互通，但主要用于核心通讯以及某些网络高级协议。不适用于上层应用业务之间的通讯。因此不能用来上公网，也不能用于对外提供服务

---

## 困扰 8. IPv6 使用多播替代广播，需要做哪些改造

IPv6 使用多播替代了广播，多播的特点是不会像广播那样完全泛洪，而是数据包只发送给加入了多播组的机器。
但是，这有个前提，就是**交换机要能识别并维护多播组的信息**，主流交换机都具备此功能，然而并不都是默认开启的。对于二层交换机来说，需要开启**MLDv2 Snooping**。
顾名思义，就是交换机会识别**“MLDv2 成员报告”**报文从哪个端口发来的，并记录下来，之后当交换机收到多播包后，会先查找其多播地址是否能在缓存里匹配上
▷ 匹配成功：仅会将数据包从相应的端口发出
▷ 匹配失败：就会泛洪，此时和广播毫无差异

---

## 困扰 9. IPv6 真的安全吗

**理想很美好**，IPv6 从设计之初，就进行了大量的安全方面的设计，“完整的”IPv6 在安全方面有至少以下 3 个优势：
▷ 原生支持的端到端加密
▷ 安全的邻居发现（**Secure Neighbor Discovery，简称 SEND**）
▷ 更大的地址空间
**现实很残酷**，只有第 3 点发挥了作用，更大的地址空间，减少了被非法扫描到的概率。而第 1、2 点并没真正普及起来，因为**协议本身就很复杂、学习难度很大、实现起来也很不容易**。因此 IETF 为了加速 IPv6 的普及，**对安全性不再强制要求**。这也导致了 IPv6 实际上并没有预期中的那么安全，在 IPv4 里存在的地址欺骗、虚假网关等情况，在 IPv6 里依然存在。

---

## 困扰 10. 如何学习 IPv6

网上能找到非常多的 IPv6 教程，其中有很多教程都是通篇讲 IPv6 地址、IP 包格式、ICMP 包格式，这很容易让初学者打退堂鼓。笔者虽然不是专业的网工，但愿能抛砖引玉，推荐的学习步骤如下：
1️⃣ IPv6 的历史、设计理念
2️⃣ IPv6 的地址格式、分类、前缀计算，以及与 IPv4 的对比
3️⃣ IP 地址、网关路由的配置与查看
4️⃣ 服务端实践，尝试给自己的网站增加 IPv6
5️⃣ 客户端实践，让自己的 PC 访问 IPv6 互联网
6️⃣ 应用层实践，自己写一对 C/S 程序，能同时支持 IPv4 与 IPv6
7️⃣ IPv6 通讯原理，抓包分析每个包，熟悉 ND、DHCPv6 等
8️⃣ IPv4 与 IPv6 的互访、过渡
9️⃣ IPv6 安全
🔟 具体领域的 IPv6，例如移动 IPv6\[

]\(https://www.bilibili.com/video/BV1J7411c7ae)

# IPv6 地址

IPv6 地址最多使用 128 bit 表示，即最多 128 个 1，这 128 bit 以 `冒号` 分割为 8 组，每组 16 bit，在使用时，使用十六进制表示。比如：`2031:0000:130F:0000:0000:09C0:876A:130B`。
![](https://notes-learning.oss-cn-beijing.aliyuncs.com/mzn9uc/1631498641648-aa666528-2155-4dfb-8969-d44a5eee7e35.jpeg)
为了书写方便，IPv6 的十六进制表示方法，可以进行压缩，具体压缩规则为：

- 每组的前导 `0` 可以省略，所以上面的例子可以写为：`2031:0:130F:0:0:9C0:876A:130B`
- 地址中包含的连续两个或多个均为 0 的组，可以用 `::`(双冒号) 代替，所以上面的例子可以进一步简写为：`2031:0:130F::9C0:876A:130B`
  - 注意：一个 IPv6 地址中，只能使用一次双冒号，否则当计算机将压缩后的地址恢复成 128 bit 时，无法确定每段中 0 的个数。

## IPv6 地址结构

IPv6 地址的这 128 bit 可以分为两部分：

- **Network prefix(网络前缀)** # n bit，相当于 IPv4 地址中的网络 ID。
  - 网络前缀由 IANA 一层层分配。IANA 组织当前划定的单播地址是 `2000::/3`。也就是说从 `2000::` 到 `3FFF:FFFF:FFFF:FFFF:FFFF:FFFF:FFFF:FFFF`
- **Interface Identify(接口标识)** # 128-n bit，相当于 IPv4 地址中的主机 ID
- **Prefix length(前缀长度)** # IPv6 没有子网掩码，通过在地址后面添加 `/NUM` 来区分一个地址的 网络前缀 和 接口表示，比如 `/64` 表示前 64 位是网络前缀，后 64 位为接口标识

接口标识可以通过三种方式生成

- **手动配置** #
- **通过软件自动生成** # 在有的地方称为有状态配置
- **通过 IEEE 的 EUI-64 规范生成** # 这是最常用的方式。在有的地方，也称为无状态配置

IEEE EUI-64 规范是将网络设备的 MAC 地址转换为 IPv6 接口标识的过程，如下图所示，MAC 地址的前 24 位（用 c 表示的部分）为公司标识，后 24 位（用 m 表示的部分）为扩展标识符。高 7 位是 0 表示了 MAC 地址本地唯一。转换的第一步将 FFFE 插入 MAC 地址的公司标识和扩展标识符之间， 第二步将高 7 位的 0 改为 1 表示此接口标识全球唯一。
![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/mzn9uc/1631458464439-6142be5c-2a5e-4c08-8aa2-87abcc2ba854.png)
例如：MAC 地址：`00-0E-0C-82-C4-D4`；转换后：`020E:0CFF:FE:82:C4D4`。&#x20;
这种由 MAC 地址产生 IPv6 地址接口标识的方法可以减少配置的工作量，尤其是当采用无状态地址自动配置时，只需要获取一个 IPv6 前缀就可以与接口标识形成 IPv6 地址。但是使用这种方式最大的缺点是任何人都可以通过二层 MAC 地址推算出三层 IPv6 地址。

## IPv6 地址分类

![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/mzn9uc/1631432277341-adbd47b9-16b3-4918-91d3-b473806fe36e.png)
IPv6 没有广播地址(以更丰富的组播地址代替)，分类如下

- **Unicast Address(单播地址)** # 全局单播地址可以分配给任何一个想要接入互联网的设备。
  - **共有地址**
    - **Global Unicast Address(全局单播地址) **# 前缀为 `2000::/3`
  - **私有地址** #
    - **Link-local Address(链路本地地址)** # 前缀为 `FE80::/10`。当两个支持 IPv6 特性的设备直连时，直连的接口会自动给自己分配一个链路本地地址，用来在没有管理员的配置下设备之间就能相互通信，并完成邻居发现等工作。这类地址的前 10 bit 是 `FE80`，后 54 bit 全是 0，最后 64 bit 是 EUI-64 地址，所以，上面例子中的链路本地地址是 `FE80::0000:09C0:876A:130B`。
      - IPv4 中的 169.254.0.0/8 就是链路本地地址。
    - **Unique Local Address(唯一本地地址)** # 前缀为 `FC00::/7`。与链路本地地址一样，区别在于链路本地地址只能用于共享链路上的设备，而站点本地地址可以用于站点内部，获得站点本地地址的设备是不能将数据包路由到站点之外的，也就是说，站点本地地址将限制数据包的传递。
      - **Site-local Address(站点本地地址)** # 在 RFC 3879 中阐述了站点本地地址会产生的问题，于 2004 年 9 月被启用，使用 Unique Local Address 代替。
    - **loopbak(环回地址)** # 前缀为 `::1/128`
- **Anycat Address(任意播地址)** #
- **Multicast Address(组播地址)** # 前缀为 `FF00::/8`，作用和 IPv4 的组播相同
  - 被请求节点组播地址 # 前缀为 `FF02::1FF00:0000/104`

> 备注：
>
> - 为什么这里叫做“全局单播地址”，而“唯一本地地址”却不叫做“唯一本地单播地址”，好吧，其实都是简称，在 RFC 里是这么定义的：“Global Unicast Addresses”、“Link-Local IPv6 Unicast Addresses”。其实“全局单播地址”是可以叫做“全局地址”的，只是这样显得有点别扭。

这是一个设备上 IPv6 的地址

```bash
以太网适配器 以太网:

   连接特定的 DNS 后缀 . . . . . . . : lan
   IPv6 地址 . . . . . . . . . . . . : 2408:8210:3c36:c1c0::a2a
   IPv6 地址 . . . . . . . . . . . . : 2408:8210:3c36:c1c0:42e:b860:ccb3:81e9
   临时 IPv6 地址. . . . . . . . . . : 2408:8210:3c36:c1c0:f4ba:1e1e:cc62:938a
   本地链接 IPv6 地址. . . . . . . . : fe80::42e:b860:ccb3:81e9%10
   IPv4 地址 . . . . . . . . . . . . : 192.168.254.245
   子网掩码  . . . . . . . . . . . . : 255.255.255.0
   默认网关. . . . . . . . . . . . . : fe80::ded8:7cff:fe11:ebdf%10
                                       192.168.254.1
```

## 全局单播地址

当我们需要给设备配置 IPv6 地址时，与 IPv4 公网地址类似，分为这么几部分：

- 网络前缀
- 子网 ID
- 接口标识

![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/mzn9uc/1631505253266-7dddac0e-c557-4e7d-973a-ea776134ef66.png)
IANA 组织当前划定的单播地址是 `2000::/3`，这个地址段占整个 IPv6 地址空间的 1/8，然后再将这个地址段逐级分下去。`/23` 是注册机构前缀，`/32` 是 ISP(运营商) 前缀，`/48` 是站点前缀，`/64` 是子网前缀。

- 根据 [RFC 4191-2.5.1](https://www.rfc-editor.org/rfc/rfc4291.html#section-2.5.1) 中的描述，对于 IPv6 所有单播地址，如果地址的前 3 bit 不是 000，则接口标识必须为 64 bit；如果地址前 3 bit 是 000，则没有此限制。
  - 这里前 3 bit 为 0，也就意味着，IP 地址从 `2000::` 开始，如果是小于 `2000::` 的话，前 3 bit 肯定是 0。所以，这里的规定也与 `2000::` 开始的地址都是单播地址这个规定对应上了。
  - 并且，通常来说，单播地址都是 `/64` 的。
- 前 64 bit
  - 前 48 位由互联网网络提供商(IANA、RIR、ISP 等)分配
  - 后面 16 是个人根据前 48 位而划分的子网
- 后 64 bit
  - 就是子网下的主机了

在 IPv6 地址空间中的脚注中，可以看到很多 IANA 保留的全局单播地址前缀：

- **2001:0::/23** # 为 IETF 协议分配 \[ RFC2928 ] 保留。
- **2001:0::/32** # 为 TEREDO \[ RFC4380 ] 保留。
- **2001:2::/48** # 保留用于基准测试 \[ RFC5180 ]\[ RFC Errata 1752 ]。
- **2001:3::/32** # 为 AMT \[ RFC7450 ] 保留。
- **2001:4:112::/48** # 为 AS112-v6 \[ RFC7535 ] 保留。
- **2001:10::/28** # 已弃用（以前称为 ORCHID）\[ RFC4843 ]。
- **2001:20::/28** # 为 ORCHIDv2 \[ RFC7343 ] 保留。
- **2001:db8::/32** # 为文档保留，也就是说在书籍、文档等地方使用 IPv6 地址示例时，使用这个前缀的地址。参考：[RFC3849](https://www.iana.org/go/rfc3849)。
- **2002::/16** # 为 6to4 \[ RFC3056 ] 保留。

# IPv6 地址的配置方式

> 参考：
> - [GitHub 组织，radvd-project](https://github.com/radvd-project)
> - [RFC 4862，IPv6 Stateless Address Autoconfiguration](https://datatracker.ietf.org/doc/html/rfc4862)

与 IPv4 一样，在 IPv6 里，任何单播地址都可以自动生成，也可以手工配置固定 IP，具体看应用场景：

- 客户端：如果我想访问 ipv6 互联网，而不对外提供服务，那么使用自动生成即可，无需使用固定的 ip 地址
- 服务端：如果需要对外提供服务，那么 ip 地址就需要固定了，不能使用自动生成

自动配置地址的场景，在 IPv6 里分为 2 种方法：“有状态”与“无状态”

- **Stateful(有状态)** # 地址由 DHCPv6 Server 统一管理，DHCPv6 Client 从中获得一个可用的 IP 地址
- **Stateless(无状态)** # 路由器发出“路由通告”报文（Router Advertisement，简称 RA），报文内包含了 IPv6 地址的前缀信息。当收到 RA 包后，就会根据其中前缀信息，自动生成一个或多个 IP 地址

## Stateless(无状态)

**Stateless Address Autoconfiguration(无状态地址自动配置，简称 SLAAC) **是 IPv6 中用于自动生成 IPv6 地址的策略。

在操作系统启动时，会自动在每个启用 IPv6 的接口上创建“链路本地地址”，当连接到网络时，会执行冲突解决。这个地址默认使用 `fe80::/64` 作为前缀。该行为使用 **Neighbor Discovery Protocol(邻居发现协议，简称 DNP)** 的一个组件，通过 SLAAC 独立进行

> 当自动生成的链路本地地址在网络上冲突时，该 IP 地址的状态会提示 [tentative noprefixroute dadfailed](#VNOuu)

除了“链路本地地址”，路由器会通过 **Router Advertisement(路由通过，简称 RA)** 提供给主机一个网络前缀，以便操作系统自动配置一个 IPv6 “单播地址”。所以，尽管存在 DHCPv6，但是 IPv6 主机通常都是使用 RA 来自动创建全局可路由的单播地址。

与 IPv4 一样，IPv6 支持全球唯一的 IP 地址。 IPv6 的设计旨在通过淘汰网络地址转换，重新强调网络设计的端到端原则，该原则最初是在早期 Internet 建立期间构想的。因此，网络上的每个设备都可以直接从任何其他设备全局寻址。

一个稳定的、唯一的、全球可寻址的 IP 地址将有助于跨网络跟踪设备。因此，对于笔记本电脑和手机等移动设备而言，此类地址是一个特殊的隐私问题。为了解决这些隐私问题，SLAAC 协议包括通常称为“隐私地址”或更准确地说是 **Temporary Addresses(临时地址)** 的内容，编入 RFC 4941，“IPv6 中无状态地址自动配置的隐私扩展”。临时地址是随机且不稳定的。典型的消费设备每天都会生成一个新的临时地址，一周后会忽略发往旧地址的流量。 Windows 自 XP SP1 起默认使用临时地址，macOS 自 (Mac OS X) 10.7 起，Android 自 4.0 起，iOS 自 4.3 版起默认使用。 Linux 发行版对临时地址的使用各不相同。

### Neighbor Discovery Protocol(邻居发现协议)

ARP 协议是 IPv4 用于解析目标 MAC 地址的协议，而在 IPv6 里，解析地址采用的是**邻居发现（Neighbor Discovery Protocol，简称 NDP 或 ND）**

ND 不是一个具体协议，而是用来描述多个相关功能的协议的**抽象集合**，所涵盖的所有协议均是基于 ICMPv6。其中有 2 种报文与解析 MAC 地址有关：

- **邻居请求报文 NS（Neighbor Solicitation）**：请求解析
- **邻居通告报文 NA（Neighbor Advertisement）**：响应解析

这与 ping 是非常类似的：

- ping：发送 icmp 的 echo request 报文，对端响应 icmp 的 echo reply 报文
- 地址解析：发送 icmp 的 ns 报文，对端响应 icmp 的 na 报文

更多详情见《[ARP 与 NDP](✏IT 学习笔记/🌐4.数据通信/通信协议/2.ARP%20 与%20NDP.md 与 NDP.md)》

### Router Advertisement(路由通告)

**Router Advertisement(路由通告，简称 RA)**

在 Linux 中，有一个名为 [radvd](https://github.com/radvd-project) 的程序，可以让服务器发送 RA 报文。

### Temporary addresses(临时地址)

无状态地址自动配置使用全局唯一和静态 MAC 地址来创建接口标识符，提供了跟踪用户设备（跨时间和 IPv6 网络前缀更改）以及用户的机会。

> 因为知道了 MAC，就可以推导出来 IPv6 地址

为了减少用户身份永久绑定到 IPv6 地址部分的可能性，节点可以创建具有基于随时间变化的随机位字符串和相对较短的生命周期（数小时到数天）的接口标识符的临时地址，之后它们替换为新地址。

临时地址可用作发起连接的源地址，而外部主机通过查询域名系统使用公共地址。

默认情况下，为 IPv6 配置的网络接口在 OS X Lion 和更高版本的 Apple 系统以及 Windows Vista、Windows 2008 Server 和更高版本的 Microsoft 系统中使用临时地址。

## 常用 IPv6 地址

运营商前缀

- 电信为 240e::/20
- 移动为 2409:8000::/20
- 联通为 2408:8000::/20

2001:db8:0:1::/64 # 通常使用这个网段的地址，很多网站(比如 Redhat、ISC-DHCP 等)的例子都是这个段，比如：

- 2001:db8:0:1::1
- 2001:db8:0:1::2
- ......
- 2001:db8:42:1::
- ......
- 2001:db8:0:1::fffe

# IPv6 Datagram 结构

IPv6 数据报被封装在链路层的 Frame 中

首部长度固定为 40 Bytes，即 320 bit。
![](https://notes-learning.oss-cn-beijing.aliyuncs.com/mzn9uc/1631514605323-4532e443-754c-44ce-8c4f-70c31b7065e8.jpeg)

- **Version(版本)** # 和 IPv4 一样，由 4 比特构成。IPv6 其版本号为 6，因此在这个字段上的值为“6”。
- **Traffic Class(通信类)** # 相当于 IPv4 的 Type Of Service(TOS) 字段，也由 8 bit 构成。但 TOS 在 IPv4 中几乎没有什么建树，未能成为卓有成效的技术，本来计划在 IPv6 中删掉这个字段，不过出于今后研究的考虑还是保留了该字段。
- **Flow Label(流标签)** # 由 20 bit 构成，准备用于 Quality Of Service(服务质量，简称 QOS) 控制。使用这个字段提供怎样的服务已经成为未来研究的课题。不适用 Qos 时每一位可以全部设置为 0。 在进行服务质量控制的时，将流标号设置为一个随机数，然后利用一种可以设置流的协议 RSVP（Resource Reservation Protocol ）在路由器上进行 Qos 设置。当某个包在发送途中需要 Qos 时，需要附上 RSVP 预想的流标号。路由器接收到这样的 IP 包后现先将流标号作为查找关键字，迅速从服务质量控制信息中查找并做相应处理。此外，只有流标号、源地址以及目标地址三项完全一致时，才被认为是一个流。
- **Payload Length(有效荷载长度)** # 上层 PDU 与 扩展首部的和，单位是字节。 IPv4 的 TL(Total Length) 是指包含首部在内的所有长度。然而 IPv6 中的这个 Playload Length 不包括首部，只表示数据部分的长度。由于 IPv6 的可选项是指连接 IPv6 首部的数据，只有当有可选项时，此处包含可选项数据的所有长度就是 Playload Length。
- **Next Header(下一个首部)** # 相当于 IPv4 中的协议字段。由 8 比特构成。通常表示 IP 的上一层协议是 TCP 或 UDP。不过在有 IPv6 扩展首部的情况下，该字段表示后面第一个扩展首部的协议类。
- **Hop Limit(跳数限制)** # 由 8 bit 构成。与 IPv4 中的 TTL 意思相同。为了强调“可通过路由器个数”这个概念，才将名字改为 Hop Limit。数据每经过一次路由器就减 1，减到 0 则丢弃数据。
- **Source Address(源地址)** # 由 128 bit 构成，表示发送端 IP 地址。
- **Destination Address(目标地址)** # 由 128 bit 构成，表示接收端 IP 地址。
- **Extension Head(扩展首部)** # IPv6 的首部长度固定，无法将可选项将入其中，取而代之的是通过扩展首部对功能进行了有效扩展。 扩展首部通常介于 IPv6 首部与 TCP/UDP 首部中间。在 IPv4 中可选项长度固定为 40 字节，但是在 IPv6 中没有这样的限制。也就是说，IPv6 的扩展首部可以是任意长度。扩展首部当中还可以包含扩展首部协议以及下一个扩展首部字段。
  - IPv6 的分片机制，也是通过 IPv6 的扩展首部来实现的。不再像 IPv4 一样，通过首部中的 标示、标志和偏移字段 进行分片识别。IPv6 网络中的中间路由器不再处理分片，只在产生数据包的源节点处理分片。
  - IPv6 数据报中包含 0 个、1 个或多个扩展首部。

## IPv6 扩展首部

![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/mzn9uc/1631515308493-c6b27a8f-9a24-4c8e-b6c7-f5e8caaa9ab6.png)

# IPv6 问题总结

使用 WireGuard 时，如果中继服务器的两端的 Peer，一个是 IPv4、一个是 IPv6，那么交互时，如果使用 ps -ef 这种命令，显示将会卡住，无法刷出后半部分内容。

## tentative noprefixroute dadfailed 状态

![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/mzn9uc/1643038958219-72693ffb-e644-42ac-b20f-320e6bb6c9eb.png)
可能的原因是 IP 地址冲突，比如虚拟机的 MAC 相同了，生成的 IPv6 就会相同

如果不换 MAC，还可以修改内核参数解决

```bash
net.ipv6.conf.all.accept_dad = 0
net.ipv6.conf.default.accept_dad = 0
net.ipv6.conf.all.use_tempaddr = 0
net.ipv6.conf.default.use_tempaddr = 0
```

# IPv6 地址的分配机制

网络前缀由 IANA 一层层分配。[IPv6 地址空间](https://www.iana.org/assignments/ipv6-address-space/ipv6-address-space.xhtml) 中包含了当前 IPv6 地址的总体分配情况，[IPv6 全局单播地址分配](https://www.iana.org/assignments/ipv6-unicast-address-assignments/ipv6-unicast-address-assignments.xhtml)中，则是单播地址的分配情况。

详见[《IP》章节的 Ip 地址分配机制](.md)部分

## ISP 分配机制

ISP 给我们分配 IPv6 时，会同时分配两个不同网段的 IPv6 地址

- 唯一地址 # 为光猫分配的地址，不会作下行分配(即，光猫内的设备不用使用这个地址进行访问)
- 前缀地址 # 为光猫内部设备分配地址的地址段

按照 IPv6 的分配规则，前缀地址必须在/60 位以内。目前中国的运营商有两种前缀，一个是 56 位，一个是 60 位。电信多是 56，联通移动多是 60 的。

> 电信为 240e::/18
> 移动为 2409:8000::/20
> 联通为 2408:8000::/20

当然，这个长度没有好与坏，不管是 56 还是 60，你获得的剩余地址量都是用不完的，哪怕是给你家里每一粒灰尘都分配上公网 IP。
问题就在于，运营商只会给你分配一次前缀地址，当你的路由器获取到前缀地址后，你其它的设备都会通过这个前缀地址向下分配剩余的地址。打个比方，通过 PPPoE 拨号，你将会获取到 WAN 口 IP 地址，这个机制和 IPv4 是一样的：
`2408:8210:4703:130b:9032:f35a:a8ab:fdf1/64`
然后你还可以获取到一个前缀 IP 地址：
`2408:8210:3c36:c1c0::1/56`
当你的手机连接 WIFI 的时候，路由器就会按照这个前缀地址给你的手机向下分配剩余的地址，如：
`2408:8210:3c36:c1c0:XXXX:XXXX:XXXX:XXXX`
\[

]\(https://cloud.tencent.com/developer/article/1468099)

##
