---
title: WireGuard
linkTitle: WireGuard
date: 2023-11-03T22:27
weight: 1
---

# 概述

> 参考：
>
> - [官网](https://www.wireguard.com/)
> - [zx2c4 源码，wireguard-linux](https://git.zx2c4.com/wireguard-linux)
>   - [GitHub 项目，WrieGuard/wireguard-linux](https://github.com/WireGuard/wireguard-linux)
> - [Wiki，WireGuard](https://en.wikipedia.org/wiki/WireGuard)
> - [张馆长博客，个人办公用 wireguard 组网笔记](https://zhangguanzhang.github.io/2020/08/05/wireguard-for-personal/)
> - [米开朗基杨博客，WireGuard 教程：WireGuard 的工作原理](https://fuckcloudnative.io/posts/wireguard-docs-theory/)

WireGuard 是一种可以实现加密 VPN 的通信协议。通常也表示为实现该通信协议的软件。

WireGuard 是由 Jason Donenfeld 等人用 C 语言编写的一个开源 VPN 协议，被视为下一代 VPN 协议，旨在解决许多困扰 IPSec/IKEv2、OpenVPN 或 L2TP 等其他 VPN 协议的问题。它与 Tinc 和 MeshBird 等现代 VPN 产品有一些相似之处，即加密技术先进、配置简单。

> [!Tip]
> 从 2020 年 1 月开始，Wireguard 已经并入了 [Linux 内核的 5.6 版本](https://github.com/torvalds/linux/blob/v5.6/drivers/net/wireguard/version.h)，这意味着大多数 Linux 发行版的用户将拥有一个开箱即用的 WireGuard。

WireGuard 没有传统的 Server 端、Client 端的概念，在 WireGuard 构建的 VPN 环境中，使用 **Peer** 来描述 VPN 中的每一个网络节点，这个 Peer 可以是 服务器、路由器、etc. 。通常来说，一个具有固定公网 IP 的 Peer，非官方得称为 **Bounce Server/Relay Server(弹跳服务器/中继服务器)**。各个在 NAT 后面的 Peer，可以通过 Bounce Server 这个 Peer 直接互通。

## Wireguard 的优缺点

[公众号-云原生实验室，WireGuard 真的很香吗？香个屁！](https://mp.weixin.qq.com/s/OvqpL9aO6oMSL4GgjE6zbw)

- 翻译自: [IPFire 博客，Why Not WireGuard](https://blog.ipfire.org/post/why-not-wireguard)

[米开朗基杨博客，WireGuard 教程：WireGuard 的工作原理](https://fuckcloudnative.io/posts/wireguard-docs-theory/)

WireGuard 优点：

- 配置精简，可直接使用默认值
- 只需最少的密钥管理工作，每个主机只需要 1 个公钥和 1 个私钥。
- 就像普通的以太网接口一样，以 Linux 内核模块的形式运行，资源占用小。
- 能够将部分流量或所有流量通过 VPN 传送到局域网内的任意主机。
- 能够在网络故障恢复之后自动重连，戳到了其他 VPN 的痛处。
- 比目前主流的 VPN 协议，连接速度要更快，延迟更低（见上图）。
- 使用了更先进的加密技术，具有前向加密和抗降级攻击的能力。
- 支持任何类型的二层网络通信，例如 ARP、DHCP 和 ICMP，而不仅仅是 TCP/HTTP。
- 可以运行在主机中为容器之间提供通信，也可以运行在容器中为主机之间提供通信。

WireGuard 不能做的事：

- 类似 gossip 协议实现网络自愈。
- 通过信令服务器突破双重 NAT。
- 通过中央服务器自动分配和撤销密钥。
- 发送原始的二层以太网帧。

当然，你可以使用 WireGuard 作为底层协议来实现自己想要的功能，从而弥补上述这些缺憾。

## WireGuard 的性能

WireGuard 声称其性能比大多数 VPN 协议更好，但这个事情有很多争议，比如某些加密方式支持硬件层面的加速。

WireGuard 直接在内核层面处理路由，直接使用系统内核的加密模块来加密数据，和 Linux 原本内置的密码子系统共存，原有的子系统能通过 API 使用 WireGuard 的 Zinc 密码库。WireGuard 使用 UDP 协议传输数据，在不使用的情况下默认不会传输任何 UDP 数据包，所以比常规 VPN 省电很多，可以像 55 一样一直挂着使用，速度相比其他 VPN 也是压倒性优势。

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/kpbis3/1616160904933-319867cc-3391-4d97-bf43-e8c40786c553.jpeg)

关于性能比较的更多信息可以参考下面几篇文档：

- wireguard.com/performance
- reddit.com/r/linux/comments/9bnowo/wireguard_benchmark_between_two_servers_with_10
- restoreprivacy.com/openvpn-ipsec-wireguard-l2tp-ikev2-protocols

# WireGuard 的工作原理

原文: [米开朗基杨博客，WireGuard 教程：WireGuard 的工作原理](https://fuckcloudnative.io/posts/wireguard-docs-theory/)

## WireGuard 术语

**Peer/Node/Device**

连接到 VPN 并为自己注册一个 VPN 子网地址（如 192.0.2.3）的主机。还可以通过使用逗号分隔的 [CIDR 表示法](/docs/4.数据通信/数据通信/CIDR.md) 指定子网范围，为其自身地址以外的 IP 地址选择路由。

**Bounce/Relay Server(中继服务器)**

一个公网可达的对等节点，可以将流量中继到 NAT 后面的其他对等节点。Bounce/Relay Server 并不是特殊的节点，它和其他对等节点一样，唯一的区别是它有公网 IP，并且开启了内核级别的 IP 转发，可以将 VPN 的流量转发到其他客户端。

**Public Endpoint(公开端点)**

节点的公网 IP 地址:端口，例如 123.124.125.126:1234，或者直接使用域名 some.domain.tld:1234。如果对等节点不在同一子网中，那么节点的公开端点必须使用公网 IP 地址。

**Private key(私钥)**

单个节点的 WireGuard 私钥，生成方法是：wg genkey > example.key。

**Public key(公钥)**

单个节点的 WireGuard 公钥，生成方式为：wg pubkey < example.key > example.key.pub。

## WireGuard 工作原理

中继服务器工作原理

中继服务器（Bounce Server）和普通的对等节点一样，它能够在 NAT 后面的 VPN 客户端之间充当中继服务器，可以将收到的任何 VPN 子网流量转发到正确的对等节点。事实上 WireGuard 并不关心流量是如何转发的，这个由系统内核和 iptables 规则处理。

如果所有的对等节点都是公网可达的，则不需要考虑中继服务器，只有当有对等节点位于 NAT 后面时才需要考虑。

在 WireGuard 里，客户端和服务端基本是平等的，差别只是谁主动连接谁而已。双方都会监听一个 UDP 端口，谁主动连接，谁就是客户端。主动连接的客户端需要指定对端的公网地址和端口，被动连接的服务端不需要指定其他对等节点的地址和端口。如果客户端和服务端都位于 NAT 后面，需要加一个中继服务器，客户端和服务端都指定中继服务器作为对等节点，它们的通信流量会先进入中继服务器，然后再转发到对端。

WireGuard 是支持漫游的，也就是说，双方不管谁的地址变动了，WireGuard 在看到对方从新地址说话的时候，就会记住它的新地址（跟 mosh 一样，不过是双向的）。所以双方要是一直保持在线，并且通信足够频繁的话（比如配置 persistent-keepalive），两边的 IP 都不固定也不影响的。

Wireguard 如何路由流量

利用 WireGuard 可以组建非常复杂的网络拓扑，这里主要介绍几个典型的拓扑：

① 端到端直接连接

这是最简单的拓扑，所有的节点要么在同一个局域网，要么直接通过公网访问，这样 WireGuard 可以直接连接到对端，不需要中继跳转。

② 一端位于 NAT 后面，另一端直接通过公网暴露

这种情况下，最简单的方案是：通过公网暴露的一端作为服务端，另一端指定服务端的公网地址和端口，然后通过 persistent-keepalive 选项维持长连接，让 NAT 记得对应的映射关系。

③ 两端都位于 NAT 后面，通过中继服务器连接

大多数情况下，当通信双方都在 NAT 后面的时候，NAT 会做源端口随机化处理，直接连接可能比较困难。可以加一个中继服务器，通信双方都将中继服务器作为对端，然后维持长连接，流量就会通过中继服务器进行转发。

④ 两端都位于 NAT 后面，通过 UDP NAT 打洞

上面也提到了，当通信双方都在 NAT 后面的时候，直接连接不太现实，因为大多数 NAT 路由器对源端口的随机化相当严格，不可能提前为双方协调一个固定开放的端口。必须使用一个信令服务器（STUN），它会在中间沟通分配给对方哪些随机源端口。通信双方都会和公共信令服务器进行初始连接，然后它记录下随机的源端口，并将其返回给客户端。这其实就是现代 P2P 网络中 WebRTC 的工作原理。有时候，即使有了信令服务器和两端已知的源端口，也无法直接连接，因为 NAT 路由器严格规定只接受来自原始目的地址（信令服务器）的流量，会要求新开一个随机源端口来接受来自其他 IP 的流量（比如其他客户端试图使用原来的通信源端口）。运营商级别的 NAT 就是这么干的，比如蜂窝网络和一些企业网络，它们专门用这种方法来防止打洞连接。更多细节请参考下一部分的 NAT 到 NAT 连接实践的章节。

如果某一端同时连接了多个对端，当它想访问某个 IP 时，如果有具体的路由可用，则优先使用具体的路由，否则就会将流量转发到中继服务器，然后中继服务器再根据系统路由表进行转发。你可以通过测量 ping 的时间来计算每一跳的长度，并通过检查对端的输出（wg show wg0）来找到 WireGuard 对一个给定地址的路由方式。

## WireGuard 报文格式

WireGuard 使用加密的 UDP 报文来封装所有的数据，UDP 不保证数据包一定能送达，也不保证按顺序到达，但隧道内的 TCP 连接可以保证数据有效交付。WireGuard 的报文格式如下图所示：

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/kpbis3/1616160904918-24e86676-6f0a-4804-8669-45d68a22e23d.jpeg)

关于 WireGuard 报文的更多信息可以参考下面几篇文档：

- wireshark.org/docs/dfref/w/wg.html
- Lekensteyn/wireguard-dissector
- nbsoftsolutions.com/blog/viewing-wireguard-traffic-with-tcpdump

## WireGuard 安全模型

WireGuard 使用以下加密技术来保障数据的安全：

- 使用 ChaCha20 进行对称加密，使用 Poly1305 进行数据验证。
- 利用 Curve25519 进行密钥交换。
- 使用 BLAKE2 作为哈希函数。
- 使用 HKDF 进行解密。

WireGuard 的加密技术本质上是 Trevor Perrin 的 Noise 框架的实例化，它简单高效，其他的 VPN 都是通过一系列协商、握手和复杂的状态机来保障安全性。WireGuard 就相当于 VPN 协议中的 qmail，代码量比其他 VPN 协议少了好几个数量级。

关于 WireGuard 加密的更多资料请参考下方链接：

- wireguard.com/papers/wireguard.pdf
- eprint.iacr.org/2018/080.pdf
- courses.csail.mit.edu/6.857/2018/project/He-Xu-Xu-WireGuard.pdf
- wireguard.com/talks/blackhat2018-slides.pdf
- arstechnica.com/gadgets/2018/08/wireguard-vpn-review-fast-connections-amaze-but-windows-support-needs-to-happen

## WireGuard 密钥管理

WireGuard 通过为每个对等节点提供简单的公钥和私钥来实现双向认证，每个对等节点在设置阶段生成密钥，且只在对等节点之间共享密钥。每个节点除了公钥和私钥，不再需要其他证书或预共享密钥。

在更大规模的部署中，可以使用 Ansible 或 Kubernetes Secrets 等单独的服务来处理密钥的生成、分发和销毁。

下面是一些有助于密钥分发和部署的服务：

- pypi.org/project/wireguard-p2p
- trailofbits/algo
- StreisandEffect/streisand
- its0x08/wg-install
- brittson/wireguard_config_maker
- wireguardconfig.com

如果你不想在 wg0.conf 配置文件中直接硬编码，可以从文件或命令中读取密钥，这使得通过第三方服务管理密钥变得更加容易：

```bash
[Interface]
...
PostUp = wg set %i private-key /etc/wireguard/wg0.key <(cat /some/path/%i/privkey)
```

从技术上讲，多个服务端之间可以共享相同的私钥，只要客户端不使用相同的密钥同时连接到两个服务器。但有时客户端会需要同时连接多台服务器，例如，你可以使用 DNS 轮询来均衡两台服务器之间的连接，这两台服务器配置相同。大多数情况下，每个对等节点都应该使用独立的的公钥和私钥，这样每个对等节点都不能读取到对方的流量，保障了安全性。

# WireGuard 关联文件与配置

**/etc/wireguard/** # WireGuard 运行时配置文件的存放路径。

# 命令行工具

https://github.com/WireGuard/wireguard-tools 包含如下两个工具

- wg
- wg-quick

## wg

> 参考：
>
> - [Manual(手册)，wg](https://www.man7.org/linux/man-pages/man8/wg.8.html)

## wg-quick

> 参考：
>
> - [Manual(手册)，wg-quick(8)](https://man7.org/linux/man-pages/man8/wg-quick.8.html)

### Syntax(语法)

**wg-quick [ up | down | save | strip ] [ CONFIG_FILE | INTERFACE ]**

# WireGuard 衍生品

[Tailscale](/docs/4.数据通信/Protocol/Tunneling%20Protocol/Tailscale/Tailscale.md)

- 自研 DERP 协议
- 一种基于 WireGuard 的虚拟组网工具

NetBird

- https://github.com/netbirdio/netbird
- https://mp.weixin.qq.com/s/amPzZb7NZCtSls0p8k-2HQ
- 简要来说 NetBird 是一个配置简易的，基于 WireGuard 的 VPN。它与 Tailscale 很像，但是区别也比较明显。**Tailscale 是在用户态实现了 WireGuard 协议**，无法使用 WireGuard 原生的命令行工具来进行管理。而 **NetBird 直接使用了内核态的 WireGuard**，可以使用命令行工具 wg 来查看和管理。

EasyTier

- https://github.com/EasyTier/EasyTier
- 一个简单、安全、去中心化的内网穿透 VPN 组网方案，使用 Rust 语言和 Tokio 框架实现。

# 待整理内容

原文: [米开朗基杨，WireGuard 教程：WireGuard 的搭建使用与配置详解](https://icloudnative.io/posts/wireguard-docs-practice/)

## 高级特性

### IPv6

前面的例子主要使用 `IPv4`，WireGuard 也支持 `IPv6`。例如：

```ini
[Interface]
AllowedIps = 192.0.2.3/24, 2001:DB8::/64
[Peer]
...
AllowedIPs = 0.0.0.0/0, ::/0
```

### 转发所有流量

如果你想通过 VPN 转发所有的流量，包括 VPN 子网和公网流量，需要在 `[Peer]` 的 `AllowedIPs` 中添加 `0.0.0.0/0, ::/0`。

即便只转发 `IPv4` 流量，也要指定一个 `IPv6` 网段，以避免将 `IPv6` 数据包泄露到 VPN 之外。详情参考：**reddit.com/r/WireGuard/comments/b0m5g2/ipv6_leaks_psa_for_anyone_here_using_wireguard_to**

例如：

```ini
[Interface]
# Name = phone.example-vpn.dev
Address = 192.0.2.3/32
PrivateKey = <private key for phone.example-vpn.dev>
[Peer]
# Name = public-server1.example-vpn.dev
PublicKey = <public key for public-server1.example-vpn.dev>
Endpoint = public-server1.example-vpn.dev:51820
AllowedIPs = 0.0.0.0/0, ::/0
```

一般只有把 VPN 当做武当纵云梯来用时，才会需要转发所有流量，不多说，点到为止。

### NAT-to-NAT 连接

如果两个对等节点（peer）都位于 NAT 后面，想不通过中继服务器直接连接，需要保证至少有一个对等节点（peer）具有稳定的公网出口，使用静态公网 IP 或者通过 `DDNS` 动态更新 `FQDN` 都可以。

`WebRTC` 协议可以动态配置两个 NAT 之间的连接，它可以通过信令服务器来检测每个主机的 `IP:Port` 组合。而 WireGuard 没有这个功能，它没有没有信令服务器来动态搜索其他主机，只能硬编码 `Endpoint+ListenPort`，并通过 `PersistentKeepalive` 来维持连接。

总结一下 NAT-to-NAT 连接的前提条件：

- 至少有一个对等节点（peer）有固定的公网 IP，如果都没有固定的公网 IP，也可以使用 `DDNS` 来维护一个稳定的域名。
- 至少有一个对等节点（peer）指定 UDP `ListenPort`，而且它的 NAT 路由器不能做 UDP 源端口随机化，否则返回的数据包将被发送到之前指定的 `ListenPort`，并被路由器丢弃，不会发送到新分配的随机端口。
- 所有的对等节点（peer）必须在 `[Peer]` 配置中启用其他对等节点（peer）的 `PersistentKeepalive`，这样就可以维持连接的持久性。

对于通信双方来说，只要**服务端**所在的 NAT 路由器没有指定到 NAT 后面的对等节点（peer）的转发规则，就需要进行 UDP 打洞。

UDP 打洞的原理：

- `Peer1` 向 `Peer2` 发送一个 UDP 数据包，不过 `Peer2` 的 NAT 路由器不知道该将这个包发给谁，直接丢弃了，不过没关系，这一步的目的是让 `Peer1` 的 NAT 路由器能够接收 UDP 响应并转发到后面的 `Peer1`。
- `Peer2` 向 `Peer1` 发送一个 UDP 数据包，由于上一步的作用，`Peer1` 的 NAT 路由器已经建立临时转发规则，可以接收 UDP 响应，所以可以接收到该数据包，并转发到 `Peer1`。
- `Peer1` 向 `Peer2` 发送一个 UDP 响应，由于上一步的作用，由于上一步的作用，`Peer2` 的 NAT 路由器已经可以接收 UDP 响应，所以可以接收到该数据包，并转发到 `Peer2`。

**这种发送一个初始的数据包被拒绝，然后利用路由器已建立的转发规则来接收响应的过程被称为 『UDP 打洞』。**

当你发送一个 UDP 数据包出去时，路由器通常会创建一个临时规则来映射源地址/端口和目的地址/端口，反之亦然。从目的地址和端口返回的 UDP 数据包会被转发到原来的源地址和端口，这就是大多数 UDP 应用在 NAT 后面的运作方式（如 BitTorrent、Skype 等）。这个临时规则会在一段时间后失效，所以 NAT 后面的客户端必须通过 `PersistentKeepalive` 定期发送数据包来维持连接的持久性。

当两个对等节点（peer）都位于 NAT 后面时，要想让 UDP 打洞生效，需要两个节点在差不多的时间向对方发送数据包，这就意味着双方需要提前知道对方的公网地址和端口号，可以在 `wg0.conf` 中指定。

#### UDP 打洞的局限性

从 2019 年开始，很多以前用过的老式打洞方法都不再有效了。以前很著名的就是 **pwnat**\[6] 开创的一种新的打洞方法，它能够在不需要代理、第三方服务器、upnp、DMZ、sproofing、dns 转换的情况下实现 NAT 中的 P2P 通信。它的原理也很简单：

通过让客户端假装成为一个互联网上任意的 `ICMP` 跳跃点（ a random hop on the Internet）来解决这个问题，从而让服务端能够获取到客户端的 IP 地址。`traceroute` 命令也是使用这项技术来检测 Internet 上的跳跃点。

具体来说，当服务器启动时，它开始向固定地址 `3.3.3.3` 发送固定的 **ICMP 回应请求包**（ICMP echo request packets）。显然，我们无法从 `3.3.3.3` 收到返回的 **ICMP 回应数据包**（ICMP echo packets）。然而，`3.3.3.3` 并不是我们可以访问的主机，我们也不是想伪装成它来发 ICMP 回应数据包。相反，pwnat 技术的实现原理在于，当我们的客户端想要连接服务端时，客户端（知道服务器 IP 地址）会向服务端送 **ICMP 超时数据包**（ICMP Time Exceeded packet）。这个 ICMP 数据包里面包含了服务端发送到 `3.3.3.3` 的原始固定 **ICMP 回应请求包**。

为什么要这样做呢？好吧，我们假装是互联网上的一个 ICMP 跳越点，礼貌地告诉服务器它原来的 **ICMP 回应请求包**无法传递到 `3.3.3.3`。而你的 NAT 是一个聪明的设备，它会注意到 **ICMP 超时数据包**内的数据包与服务器发出 **ICMP 回应请求包**相匹配。然后它将 **ICMP 超时数据包**转发回 NAT 后面的服务器，包括来自客户端的完整 IP 数据包头，从而让服务端知道客户端 IP 地址是什么！

现在这种类似的 UDP 打洞方法受到了很多的限制，详情可以参考[上篇文章](https://mp.weixin.qq.com/s/o6OyuFBFanTcp3-XnlYjlw)，这里不过多阐述。除了 UDP 打洞之外，我们仍然可以使用硬编码的方式指定两个对等节点（peer）的公网地址和端口号，这个方法对大多数 NAT 网络都有效。

#### 源端口随机化

如果所有的对等节点（peer）都在具有严格的 UDP 源端口随机化的 NAT 后面（比如大多数蜂窝网络），那么无法实现 `NAT-to-NAT` 连接。因为双方都无法协商出一个 `ListenPort`，并保证自己的 NAT 在发出 ping 包后能够接收发往该端口的流量，所以就无法初始化打洞，导致连接失败。因此，一般在 `LTE/3G` 网络中无法进行 p2p 通信。

#### 使用信令服务器

上节提到了，如果所有的对等节点（peer）都在具有严格的 UDP 源端口随机化的 NAT 后面，就无法直接实现 `NAT-to-NAT` 连接，但通过第三方的信令服务器是可以实现的。信令服务器相当于一个中转站，它会告诉通信双方关于对方的 `IP:Port` 信息。这里有几个项目可以参考：

- **takutakahashi/wg-connect**
- **git.zx2c4.com/wireguard-tools/tree/contrib/nat-hole-punching**

#### 动态 IP 地址

WireGuard 只会在启动时解析域名，如果你使用 `DDNS` 来动态更新域名解析，那么每当 IP 发生变化时，就需要重新启动 WireGuard。目前建议的解决方案是使用 `PostUp` 钩子每隔几分钟或几小时重新启动 WireGuard 来强制解析域名。

总的来说，`NAT-to-NAT` 连接极为不稳定，而且还有一堆其他的限制，所以还是建议通过中继服务器来通信。

NAT-to-NAT 配置示例：

Peer1：

```ini
[Interface]
...
ListenPort = 12000
[Peer]
...
Endpoint = peer2.example-vpn.dev:12000
PersistentKeepalive = 25
```

Peer2：

```ini
[Interface]
...
ListenPort = 12000
[Peer]
...
Endpoint = peer1.example-vpn.dev:12000
PersistentKeepalive = 25
```

更多资料：

- samyk/pwnat
- en.wikipedia.org/wiki/UDP_hole_punching
- stackoverflow.com/questions/8892142/udp-hole-punching-algorithm
- stackoverflow.com/questions/12359502/udp-hole-punching-not-going-through-on-3g
- stackoverflow.com/questions/11819349/udp-hole-punching-not-possible-with-mobile-provider
- WireGuard/WireGuard@master/contrib/examples/nat-hole-punching
- staaldraad.github.io/2017/04/17/nat-to-nat-with-wireguard
- golb.hplar.ch/2019/01/expose-server-vpn.html

### 动态分配子网 IP

这里指的是对等节点（peer）的 VPN 子网 IP 的动态分配，类似于 DHCP，不是指 `Endpoint`。

WireGuard 官方已经在开发动态分配子网 IP 的功能，具体的实现可以看这里：**WireGuard/wg-dynamic**

当然，你也可以使用 `PostUp` 在运行时从文件中读取 IP 值来实现一个动态分配 IP 的系统，类似于 Kubernetes 的 CNI 插件。例如：

```ini
[Interface]
...
PostUp = wg set %i allowed-ips /etc/wireguard/wg0.key <(some command)
```

### 奇技淫巧

#### 共享一个 peers.conf 文件

介绍一个秘密功能，可以简化 WireGuard 的配置工作。如果某个 `peer` 的公钥与本地接口的私钥能够配对，那么 WireGuard 会忽略该 `peer`。利用这个特性，我们可以在所有节点上共用同一个 peer 列表，每个节点只需要单独定义一个 `[Interface]` 就行了，即使列表中有本节点，也会被忽略。具体方式如下：

- 每个对等节点（peer）都有一个单独的 `/etc/wireguard/wg0.conf` 文件，只包含 `[Interface]` 部分的配置。
- 每个对等节点（peer）共用同一个 `/etc/wireguard/peers.conf` 文件，其中包含了所有的 peer。
- Wg0.conf 文件中需要配置一个 PostUp 钩子，内容为 `PostUp = wg addconf /etc/wireguard/peers.conf`。

关于 `peers.conf` 的共享方式有很多种，你可以通过 `ansible` 这样的工具来分发，可以使用 `Dropbox` 之类的网盘来同步，当然也可以使用 `ceph` 这种分布式文件系统来将其挂载到不同的节点上。

#### 从文件或命令输出中读取配置

WireGuard 也可以从任意命令的输出或文件中读取内容来修改配置的值，利用这个特性可以方便管理密钥，例如可以在运行时从 `Kubernetes Secrets` 或 `AWS KMS` 等第三方服务读取密钥。

### 容器化

WireGuard 也可以跑在容器中，最简单的方式是使用 `--privileged` 和 `--cap-add=all` 参数，让容器可以加载内核模块。

你可以让 WireGuard 跑在容器中，向宿主机暴露一个网络接口；也可以让 WireGuard 运行在宿主机中，向特定的容器暴露一个接口。

下面给出一个具体的示例，本示例中的 `vpn_test` 容器通过 WireGuard 中继服务器来路由所有流量。本示例中给出的容器配置是 `docker-compose` 的配置文件格式。

中继服务器容器配置：

```yaml
version: '3'
services:
  wireguard:
    image: linuxserver/wireguard
    ports:
      - 51820:51820/udp
    cap_add:
      - NET_ADMIN
      - SYS_MODULE
    volumes:
      - /lib/modules:/lib/modules
      - ./wg0.conf:/config/wg0.conf:ro
```

中继服务器 WireGuard 配置 `wg0.conf`：

```ini
[Interface]
# Name = relay1.wg.example.com
Address = 192.0.2.1/24
ListenPort = 51820
PrivateKey = oJpRt2Oq27vIB5/UVb7BRqCwad2YMReQgH5tlxz8YmI=
DNS = 1.1.1.1,8.8.8.8
PostUp = iptables -A FORWARD -i wg0 -j ACCEPT; iptables -t nat -A POSTROUTING -o eth0 -j MASQUERADE; ip6tables -A FORWARD -i wg0  -j ACCEPT; ip6tables -t nat -A POSTROUTING -o eth0 -j MASQUERADE
PostDown = iptables -D FORWARD -i wg0 -j ACCEPT; iptables -t nat -D POSTROUTING -o eth0 -j MASQUERADE; ip6tables -D FORWARD -i wg0 -j ACCEPT; ip6tables -t nat -D POSTROUTING -o eth0 -j MASQUERADE
[Peer]
# Name = peer1.wg.example.com
PublicKey = I+hXRAJOG/UE2IQvIHsou2zTgkUyPve2pzvHTnd/2Gg=
AllowedIPs = 192.0.2.2/32
```

客户端容器配置：

```yaml
version: '3'
services:
  wireguard:
    image: linuxserver/wireguard
    cap_add:
      - NET_ADMIN
      - SYS_MODULE
    volumes:
      - /lib/modules:/lib/modules
      - ./wg0.conf:/config/wg0.conf:ro

  vpn_test:
    image: curlimages/curl
    entrypoint: curl -s http://whatismyip.akamai.com/
    network_mode: 'service:wireguard'
```

客户端 WireGuard 配置 `wg0.conf`：

```ini
[Interface]
# Name = peer1.wg.example.com
Address = 192.0.2.2/32
PrivateKey = YCW76edD4W7nZrPbWZxPZhcs32CsBLIi1sEhsV/sgk8=
DNS = 1.1.1.1,8.8.8.8
[Peer]
# Name = relay1.wg.example.com
Endpoint = relay1.wg.example.com:51820
PublicKey = zJNKewtL3gcHdG62V3GaBkErFtapJWsAx+2um0c0B1s=
AllowedIPs = 192.0.2.1/24,0.0.0.0/0
PersistentKeepalive = 21
```
