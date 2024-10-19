---
title: CIDR
---

# 概述

> 参考：
>
> - [Wiki, CIDR](https://en.wikipedia.org/wiki/Classless_Inter-Domain_Routing)

**Classless Inter-Domain Routing(无类别域间路由，简称 CIDR)** 是一个用于给用户分配[IP 地址](https://zh.wikipedia.org/wiki/IP%E5%9C%B0%E5%9D%80)以及在[互联网](https://zh.wikipedia.org/wiki/%E4%BA%92%E8%81%94%E7%BD%91)上有效地路由 IP[数据包](https://zh.wikipedia.org/wiki/%E6%95%B0%E6%8D%AE%E5%8C%85)的对 IP 地址进行归类的方法。

> Notes: CIDR 可以看作是网段的一种表示方法。是一种表示 **IP Blocks(网络地址块)** 或 **SubNet(子网)** 的方法

在[域名系统](https://zh.wikipedia.org/wiki/%E5%9F%9F%E5%90%8D%E7%B3%BB%E7%BB%9F)出现之后的第一个十年里，基于[分类网络](https://zh.wikipedia.org/wiki/%E5%88%86%E7%B1%BB%E7%BD%91%E7%BB%9C)进行地址分配和路由 IP 数据包的设计就已明显显得[可扩充性](https://zh.wikipedia.org/wiki/%E5%8F%AF%E6%89%A9%E6%94%BE%E6%80%A7)不足（参见 RFC 1517）。为了解决这个问题，[互联网工程工作小组](https://zh.wikipedia.org/wiki/%E4%BA%92%E8%81%94%E7%BD%91%E5%B7%A5%E7%A8%8B%E5%B7%A5%E4%BD%9C%E5%B0%8F%E7%BB%84)在 1993 年发布了一新系列的标准 RFC 1518 和 RFC 1519 以定义新的分配 IP 地址块和路由[IPv4](https://zh.wikipedia.org/wiki/IPv4)数据包的方法。

一个 IP 地址包含两部分：标识网络的**前缀**和紧接着的在这个网络内的**主机地址**。在之前的[分类网络](https://zh.wikipedia.org/wiki/%E5%88%86%E7%B1%BB%E7%BD%91%E7%BB%9C)中，IP 地址的分配把 IP 地址的 32 位按每 8 位为一段分开。这使得前缀必须为 8，16 或者 24 位。因此，可分配的最小的地址块有 256（24 位前缀，8 位主机地址，28=256）个地址，而这对大多数企业来说太少了。大一点的地址块包含 65536（16 位前缀，16 位主机，216=65536）个地址，而这对大公司来说都太多了。这导致不能充分使用 IP 地址和在路由上的不便，因为大量的需要单独路由的小型网络（C 类网络）因在地域上分得很开而很难进行[聚合路由](https://zh.wikipedia.org/w/index.php?title=%E8%81%9A%E5%90%88%E8%B7%AF%E7%94%B1&action=edit&redlink=1)，于是给路由设备增加了很多负担。

无类别域间路由是基于 **Variable-length Subnet Masking(可变长子网掩码，简称 VLSM)** 来进行任意长度的前缀的分配。在 RFC 950（1985）中有关于可变长子网掩码的说明。CIDR 包括：

- 指定任意长度的前缀的可变长子网掩码技术。遵从 CIDR 规则的地址有一个后缀说明前缀的位数，例如：192.168.0.0/16。这使得对日益缺乏的 IPv4 地址的使用更加有效。
- 将多个连续的前缀聚合成[**超网**](https://zh.wikipedia.org/w/index.php?title=%E8%B6%85%E7%BD%91&action=edit&redlink=1)，以及，在互联网中，只要有可能，就显示为一个聚合的网络，因此在总体上可以减少路由表的表项数目。聚合使得互联网的路由表不用分为多级，并通过 VLSM 逆转“划分子网”的过程。
- 根据机构的实际需要和短期预期需要而不是分类网络中所限定的过大或过小的地址块来管理 IP 地址的分配的过程。

因为在 [IPv6](https://zh.wikipedia.org/wiki/IPv6)中也使用了 IPv4 的用后缀指示前缀长度的 CIDR，所以 IPv4 中的分类在 IPv6 中已不再使用。
