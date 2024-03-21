---
title: Tunneling Protocol
weight: 1
---

# 概述

> 参考：
> 
> - [Wiki，Tunneling Protocol](https://en.wikipedia.org/wiki/Tunneling_protocol)
> - [Wiki，Overlay Network](https://en.wikipedia.org/wiki/Overlay_network)

**Tunneling Protocol(隧道协议)** 是一种通信协议，允许数据从一个网络移动到另一个网络。该协议通过 [通信协议](/docs/4.数据通信/通信协议/通信协议.md) 中 [Encapsulation(封装)](/docs/4.数据通信/通信协议/通信协议.md#Encapsulation(封装)) 的过程跨公共网络发送专用网络通信。因为隧道涉及将流量数据重新打包为不同的形式，可能以加密为标准，它可以隐藏通过隧道运行的流量的性质。隧道协议通过使用数据包的 Payload(数据部分) 来承载实际提供服务的数据包。隧道使用分层协议模型，例如 OSI 或 TCP/IP 协议套件中的那些，但在使用有效载荷承载网络通常不提供的服务时通常会违反分层。通常，在分层模型中，传送协议在与有效载荷协议相同或更高的级别上运行。

隧道技术是一种通过使用互联网络的基础设施在网络之间传递数据的方式。使用隧道传递的数据（或负载）可以是不同协议的数据帧或包。隧道协议将其它协议的数据帧或包重新封装然后通过隧道发送。新的帧头提供路由信息，以便通过互联网传递被封装的负载数据。

隧道的常见用途

- 隧道协议可以允许外部协议在不支持该特定协议的网络上运行，例如在 IPv4 上运行 IPv6。
- 另一个重要用途是提供仅使用底层网络服务提供的不切实际或不安全的服务，例如向其物理网络地址不属于公司网络的远程用户提供公司网络地址。
- 用户还可以使用隧道“潜入”防火墙，使用防火墙通常会阻止的协议，但“包装”在防火墙不会阻止的协议中，例如 HTTP。如果防火墙策略没有明确排除这种“包装”，则此技巧可以绕过预期的防火墙策略（或任何一组互锁的防火墙策略）。
- 另一种基于 HTTP 的隧道方法使用 HTTP CONNECT 方法/命令。客户端向 HTTP 代理发出 HTTP CONNECT 命令。然后，代理与特定的 server:port 建立 TCP 连接，并在该 server:port 和客户端连接之间中继数据。 \[1]因为这会产生安全漏洞，所以支持 CONNECT 的 HTTP 代理通常会限制对 CONNECT 方法的访问。代理仅允许连接到特定端口，例如 HTTPS 的 443。

## Overlay

**Overlay(叠加网络)** 实际上是一种隧道封装技术，是对隧道技术的扩展。传统隧道技术仅限于隧道两端通信，而 Overlay 网络则可以实现 N 个端点之间的互相通信。

## VPN

> 参考：
> 
> - [Wiki，VPN](https://en.wikipedia.org/wiki/Virtual_private_network)

**Virtual Private Network(虚拟专用网络，简称 VPN)** 是通过**隧道协议**建立的虚拟点对点连接。可以从逻辑上，让人们将通过 VPN 将两个或多个互不连接的网络打通，组成一个更大型的局域网。

VPN 程序

- [tinc](https://github.com/gsliepen/tinc) 是一个虚拟专用网络 (VPN) 守护程序

# Tunnel

应用场景：

- 一个公司在天津与北京分别有一个办公地点，天津的内网为 10.0.0.0/24，北京的内网为 10.0.1.0/24。那么如何让两个内网互通呢?可以使用 tunnel 技术，在两地公网出口建立隧道连接。天津访问北京的时候，目的内网地址是封装在公网 IP 里面的，这样就可以让私网地址的数据在公网传输。比如大企业都有自己的隧道网络，当使用个人电脑，安装上某些隧道软件后，那么这台电脑就可以访问公司内部网络了。

Tunnel 技术的实现方式：

- 基于数据包:
  - IP in IP，比 GRE 更小的负载头，并且适合只有负载一个 IP 流的情况。
  - [GRE](/docs/4.数据通信/通信协议/Tunneling%20Protocol/GRE.md)，支持多种网络层协议和多路技术
  - [PPTP](/docs/4.数据通信/通信协议/Tunneling%20Protocol/PPTP.md)（点对点隧道协议）
  - SSTP 安全的 PPTP
  - IPsec/L2TP（数据链接层隧道协议）
  - [WireGuard](/docs/4.数据通信/通信协议/Tunneling%20Protocol/WireGuard/WireGuard.md)
  - 依赖其他协议实现的隧道功能
    - SSL 
    - SSH
    - SOCKS
    - 等

## 实现各种隧道协议的系统

[libreswan](https://github.com/libreswan/libreswan) # IPsec 服务器

[xl2tpd](https://github.com/xelerance/xl2tpd) # L2TP 提供者

[WireGuard](docs/4.数据通信/通信协议/Tunneling%20Protocol/WireGuard/WireGuard.md)

OpenVPN # 基于 SSL 的 VPN 系统，广泛使用 OpenSSL 加密库和 TLS 协议。

# Overlay

overlay 技术，一般都需要一个第三方程序来实现(这个程序可以是一个守护进程、内核等等)。这个程序的作用就是让本身并不互通的内部网络可以互通(比如不同宿主机上的容器互通 2.flannel.note 或者不同宿主机上的虚拟机互通或者)。而让这些网络互通的信息(比如路由、arp 等)就是由这个第三方程序来保存并维护的。如下图所示

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/qw0o0m/1616160946658-7d2f69f8-d44e-4bd6-981c-5752d093bdff.jpeg)

现在想要 10.244.0.1 可以 ping 通 10.244.1.1。那么 node1 就需要知道 10.244.1.1 在哪台宿主机上。而这种信息，就是靠可以实现 overlay 功能的程序或者内核(甚至几条路由表的规则)来保存并维护的。凡是连接到 overlay 程序的容器，其本身的 IP 以及所在宿主机的 IP，都会被保存下来，以便通信时可以使用。而这些数据，宿主机并不需要知道，在把要发送的数据包交给宿主机的网络栈时，数据包中的目的地址，如果是 10.244.1.0/24 网段的，那么 192.168.0.1 宿主机本身并不知道，所以，overlay 程序的其中一个功能，就是在数据包外面进行封装，把本身的目的地址掩盖起来，并把自己维护的信息中的可以被宿主机识别的 IP 地址填上，这样，就会形成 IP 套 IP 的效果，也就是上面的 tunnel 的效果。这时候，当宿主机收到数据包的时候，就会清楚的知道，要发送到哪里，而 node2 再接收到数据包并解封装后，会看到 overlay 程序封装的信息，就会把数据包交给本机的 overlay 程序，进行后续处理.

Overlay 技术的实现就是 VXLAN，关于 VXLAN 的介绍，可以参考 2.flannel.note 中的 VXLAN 模型，其中有关于 VXLAN 工作流程的详细讲解

# 分类

> #网络 #隧道协议
