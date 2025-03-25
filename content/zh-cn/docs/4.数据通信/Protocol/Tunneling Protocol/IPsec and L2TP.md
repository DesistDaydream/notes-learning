---
title: IPsec and L2TP
linkTitle: IPsec and L2TP
weight: 2
---

# 概述

> 参考：
>
> - [Wiki, IPsec](https://en.wikipedia.org/wiki/IPsec)
> - [Wiki, L2TP](https://en.wikipedia.org/wiki/Layer_2_Tunneling_Protocol)

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

# IKEv2 客户端

> 参考：
>
> - https://github.com/hwdsl2/setup-ipsec-vpn/blob/master/docs/ikev2-howto-zh.md

现代操作系统支持 IKEv2 协议标准。因特网密钥交换（英语：Internet Key Exchange，简称 IKE 或 IKEv2）是一种网络协议，归属于 IPsec 协议族之下，用以创建安全关联 (Security Association, SA)。与 IKE 版本 1 相比较，IKEv2 的 [功能改进](https://en.wikipedia.org/wiki/Internet_Key_Exchange#Improvements_with_IKEv2) 包括比如通过 MOBIKE 实现 Standard Mobility 支持，以及更高的可靠性。

Libreswan 支持通过使用 RSA 签名算法的 X.509 Machine Certificates 来对 IKEv2 客户端进行身份验证。该方法无需 IPsec PSK, 用户名或密码。它可以用于 Windows, macOS, iOS, Android, Chrome OS, Linux 和 RouterOS。

默认情况下，运行 VPN 安装脚本时会自动配置 IKEv2。如果你想了解有关配置 IKEv2 的更多信息，请参见 [使用辅助脚本配置 IKEv2](https://github.com/hwdsl2/setup-ipsec-vpn/blob/master/docs/ikev2-howto-zh.md#%E4%BD%BF%E7%94%A8%E8%BE%85%E5%8A%A9%E8%84%9A%E6%9C%AC%E9%85%8D%E7%BD%AE-ikev2)。Docker 用户请看 [配置并使用 IKEv2 VPN](https://github.com/hwdsl2/docker-ipsec-vpn-server/blob/master/README-zh.md#%E9%85%8D%E7%BD%AE%E5%B9%B6%E4%BD%BF%E7%94%A8-ikev2-vpn)。
