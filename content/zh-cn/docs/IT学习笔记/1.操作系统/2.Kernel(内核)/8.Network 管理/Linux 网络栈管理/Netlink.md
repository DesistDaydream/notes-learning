---
title: Netlink
---

# 概述

> 参考：
> - [Manual(手册),netlink](https://man7.org/linux/man-pages/man7/netlink.7.html)
> - [Manual(手册),rtnetlink](https://man7.org/linux/man-pages/man7/rtnetlink.7.html)
> - [Wiki,Netlink](https://en.wikipedia.org/wiki/Netlink)
> - [内核官方文档,Linux 网络文档-通用 Netlink](https://www.kernel.org/doc/html/latest/networking/generic_netlink.html)

**Netlink** 是一个 Linux 内核接口，用于在 内核 与 用户空间进程 之间传输信息。还可以用作两个用户空间进程之间、甚至内核子系统之间的数据通信。说白了，就是一个通过 Socket 实现 IPC 的方式。

[Iproute2](✏IT 学习笔记/📄1.操作系统/X.Linux%20 管理/Linux%20 网络管理工具/Iproute%20 工具包.md 管理/Linux 网络管理工具/Iproute 工具包.md)、keepalived、ethtool 等等 应用程序，很多功能都是基于 Netlink 开发的。

Netlink 由两部分组成：

## Rtnetlink 概述

rtnetlink 是 Linux 路由套接字

RTNETLINK 允许读取和更改内核的路由表。它在内核中使用以在各种子系统之间进行通信，尽管此处未记录此使用，并且与用户空间程序通信。可以通过 NetLink_Route 套接字来控制网络路由，IP 地址，链接参数，邻居设置，排队学科，流量类和数据包分类器。它基于[NetLink](https://man7.org/linux/man-pages/man7/netlink.7.html) 消息;有关更多信息。
