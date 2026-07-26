---
title: Tunneling Protocol
weight: 1
---

# 概述

> 参考：
>
> - [Wiki, Tunneling Protocol](https://en.wikipedia.org/wiki/Tunneling_protocol)

**Tunneling Protocol(隧道协议)** 是一种通信协议，允许数据从一个网络移动到另一个网络。该协议通过通信协议中 [Encapsulation(封装)](/docs/4.数据通信/Protocol/Communication%20protocol.md#Encapsulation(封装)) 的过程跨公共网络发送专用网络通信。因为隧道涉及将流量数据重新打包为不同的形式，可能以加密为标准，它可以隐藏通过隧道运行的流量的性质。隧道协议通过使用数据包的 Payload(数据部分) 来承载实际提供服务的数据包。隧道使用分层协议模型，例如 OSI 或 TCP/IP 协议套件中的那些，但在使用有效载荷承载网络通常不提供的服务时通常会违反分层。通常，在分层模型中，传送协议在与有效载荷协议相同或更高的级别上运行。

隧道技术是一种通过使用互联网络的基础设施在网络之间传递数据的方式。使用隧道传递的数据（或负载）可以是不同协议的数据帧或包。隧道协议将其它协议的数据帧或包重新封装然后通过隧道发送。新的帧头提供路由信息，以便通过互联网传递被封装的负载数据。

隧道的常见用途

- 隧道协议可以允许外部协议在不支持该特定协议的网络上运行，例如在 IPv4 上运行 IPv6。
- 另一个重要用途是提供仅使用底层网络服务提供的不切实际或不安全的服务，例如向其物理网络地址不属于公司网络的远程用户提供公司网络地址。
- 用户还可以使用隧道“潜入”防火墙，使用防火墙通常会阻止的协议，但“包装”在防火墙不会阻止的协议中，例如 HTTP。如果防火墙策略没有明确排除这种“包装”，则此技巧可以绕过预期的防火墙策略（或任何一组互锁的防火墙策略）。
- 另一种基于 HTTP 的隧道方法使用 HTTP CONNECT 方法/命令。客户端向 HTTP 代理发出 HTTP CONNECT 命令。然后，代理与特定的 server:port 建立 TCP 连接，并在该 server:port 和客户端连接之间中继数据。 \[1]因为这会产生安全漏洞，所以支持 CONNECT 的 HTTP 代理通常会限制对 CONNECT 方法的访问。代理仅允许连接到特定端口，例如 HTTPS 的 443。

应用场景：

- 一个公司在天津与北京分别有一个办公地点，天津的内网为 10.0.0.0/24，北京的内网为 10.0.1.0/24。那么如何让两个内网互通呢?可以使用 tunnel 技术，在两地公网出口建立隧道连接。天津访问北京的时候，目的内网地址是封装在公网 IP 里面的，这样就可以让私网地址的数据在公网传输。比如大企业都有自己的隧道网络，当使用个人电脑，安装上某些隧道软件后，那么这台电脑就可以访问公司内部网络了。


## 扩展技术

**[Overlay network](/docs/4.数据通信/Protocol/Tunneling%20Protocol/Overlay%20network.md)(叠加网络)** 实际上是一种隧道封装技术，是对隧道技术的扩展。传统隧道技术仅限于隧道两端通信，而 Overlay network 则可以实现 N 个端点之间的互相通信。

**[VPN](/docs/4.数据通信/VPN/VPN.md)(虚拟专用网络)**

# Protocol

Tunnel 技术的实现方式：

- 基于数据包:
    - IP in IP，比 GRE 更小的负载头，并且适合只有负载一个 IP 流的情况。
    - [GRE](/docs/4.数据通信/Protocol/Tunneling%20Protocol/GRE.md)，支持多种网络层协议和多路技术
    - [PPTP](/docs/4.数据通信/Protocol/Tunneling%20Protocol/PPTP.md)（点对点隧道协议）
    - SSTP 安全的 PPTP
    - IPsec/L2TP（数据链接层隧道协议）
    - [WireGuard](/docs/4.数据通信/Protocol/Tunneling%20Protocol/WireGuard/WireGuard.md)
    - 依赖其他协议实现的隧道功能
        - SSL
        - SSH
        - SOCKS
        - 等

**Encrypted Tunnel Protocol(加密的隧道协议)** 用以保障通信安全。很长时间以来，Linux 中加密隧道的标准解决方案是 IPsec

## Obfuscation Tunnel Protocol

> [!Question]
> 加密的隧道协议虽然可以让流量安全得传输，但是却没法规避**流量审查**。当审查机关想要禁止用户访问流量时候，只要识别到加密隧道协议的特征，就能禁掉。怎么办？

所以，我们需要一种手段，来让流量**不那么 “显眼”**，隐藏代理特征，让流量看起来像正常流量。比如让实现隧道协议的服务端，对外看起来就像普通的 HTTPS 网站一样。

**Obfuscation Tunnel Protocol(混淆的隧道协议)** 在 Encrypted Tunnel Protocol 之上对流量进行混淆。很多这类协议都被用来规避审查，以便从中国境内可以顺利访问国际互联网。

> 对于中国用户来说，这种协议主要是用来对抗 **Great Firewall of China(简称 GFW)**，GFW 在中国互联网审查中的作用是阻止对选定外国网站的访问，并减缓跨境互联网流量。

- Shadowsocks
    - https://en.wikipedia.org/wiki/Shadowsocks
- [VMess](/docs/4.数据通信/Protocol/Tunneling%20Protocol/VMess.md)
- VLESS # VMess 简化版
- Trojan/Trojan-Go
- Hysteria
    - https://github.com/apernet/hysteria

# 支持多种隧道协议的客户端

除了每个协议所属项目自己实现的 服务端、客户端之外，在 [Proxy Client](/docs/Web/Proxy/Proxy%20Client.md) 文章中，还可以找到一些支持多种协议的客户端
