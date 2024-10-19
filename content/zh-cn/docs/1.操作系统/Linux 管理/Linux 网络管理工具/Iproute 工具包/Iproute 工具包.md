---
title: Iproute 工具包
weight: 1
---

# 概述

> 参考：
>
> - [官方文档](https://wiki.linuxfoundation.org/networking/iproute2)
> - [Wiki, Iproute2](https://en.wikipedia.org/wiki/Iproute2)

Iprtoue2 是一组应用程序的集合，用于管理 Linux 网络栈。可以控制、监控 LInux 内核中网络栈的各个方面，包括路由、网络接口、隧道、流量控制、与网络相关的设备驱动程序。

Iproute2 基于 Linux 的 [Netlink](/docs/2.编程/高级编程语言/Go/Go%20第三方库/网络栈控制/Netlink/Netlink.md) 接口与 LInux 内核通信，以实现网络栈管理功能。Iproute2 的发展与内核网络组件的发展密切相关，原作者 Alexey Kuznetsov 负责 Linux 内核中的 QoS 实现，并且各种文档由 [Linux 基金会 Wiki](https://wiki.linuxfoundation.org/) 维护，且代码也存在于 [Linux 内核代码](https://git.kernel.org/pub/scm/network/iproute2/iproute2.git)中

## 该工具包包含如下工具

1. arpd
2. bridge # 显示或操纵 Linux 网桥 地址和设备
3. cbq
4. ctstat
5. devlink
6. genl
7. ifcfg
8. ifstat
9. [ip](/docs/1.操作系统/Linux%20管理/Linux%20网络管理工具/Iproute%20工具包/ip/ip.md) # 显示或操纵 routing, devices, policy routing and tunnels
10. lnstat
11. nstat
12. rdma
13. routef
14. routel
15. rtacct
16. rtmon
17. rtpr
18. rtstat
19. [ss](/docs/1.操作系统/Linux%20管理/Linux%20网络管理工具/Iproute%20工具包/ss.md) # 转存 Socket 信息
20. tipc #
21. tc # 实现 [TC 模块](/docs/1.操作系统/Kernel/Network/Linux%20网络流量控制/TC%20模块/TC%20模块.md) 进行流量控制的程序

# 关联文件与配置

**/etc/iproute2/** #
