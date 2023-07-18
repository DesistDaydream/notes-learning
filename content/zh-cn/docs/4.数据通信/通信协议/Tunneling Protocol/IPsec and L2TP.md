---
title: "IPsec and L2TP"
linkTitle: "IPsec and L2TP"
date: "2023-07-18T12:48"
weight: 20
---

# 概述

> 参考：
> 
> - [WIki，IPsec](https://en.wikipedia.org/wiki/IPsec)
> - [Wiki，L2TP](https://en.wikipedia.org/wiki/Layer_2_Tunneling_Protocol)

**Internet Protocol Security(因特网协议安全，简称 IPsec)** 是一种安全网络协议套件，可对数据包进行身份验证和加密，以通过 Internet 协议网络在两台计算机之间提供安全的加密通信。它用于虚拟专用网络 (VPN)。

**Layer 2 Tunneling Protocol(2层隧道协议，简称 L2TP)** 是一种隧道协议，用于支持虚拟专用网络 (VPN) 或作为 ISP 提供服务的一部分。它仅对其自己的控制消息使用加密（“隐藏”）（使用可选的预共享秘密），并且其本身不提供任何内容加密或机密性。相反，它为第 2 层提供隧道（可以加密），并且隧道本身可以通过第 3 层加密协议（例如 IPsec）进行传递。

[知乎，既然IPsec有隧道模式，为什么还有L2TP+IPsec这样的组合？](https://www.zhihu.com/question/279686051)

您好，这是必应。根据我从网络上搜索到的信息，我给您简单介绍一下 ipsec 和 l2tp 的区别：

- ipsec 是一种安全协议，它可以在网络层对数据进行加密、认证和完整性保护，防止数据被窃听、篡改或伪造。
- l2tp 是一种隧道协议，它可以在数据链路层对数据进行封装和传输，实现不同网络之间的透明连接。
- ipsec 和 l2tp 可以单独使用，也可以结合使用。当单独使用时，ipsec可以提供数据的安全性，但不能跨越不同的网络类型；l2tp可以提供数据的透明性，但不能保证数据的安全性。
- 当结合使用时，ipsec 和 l2tp 可以互补，实现数据的安全和透明传输。一般有两种组合方式：ipsec over l2tp 和 l2tp over ipsec。前者是先用ipsec加密数据，再用l2tp封装数据；后者是先用l2tp封装数据，再用 ipsec 加密数据。

# Libreswan

https://github.com/libreswan/libreswan

https://libreswan.org/

https://github.com/hwdsl2/setup-ipsec-vpn/blob/master/README-zh.md

- https://github.com/hwdsl2/docker-ipsec-vpn-server/blob/master/README-zh.md