---
title: Iproute 工具包
---

# 概述

> 参考：
> - 官方文档：<https://wiki.linuxfoundation.org/networking/iproute2>
> - [Wiki,Iproute2](https://en.wikipedia.org/wiki/Iproute2)

Iprtoue2 是一组应用程序的集合，用于管理 Linux 网络栈。可以控制、监控 LInux 内核中网络栈的各个方面，包括路由、网络接口、隧道、流量控制、与网络相关的设备驱动程序。

iproute2 基于 Linux 的 [Netlink 接口](https://www.yuque.com/go/doc/44482585) 与 LInux 内核通信，以实现网络栈管理功能。

## 该工具包包含如下工具

1. arpd

2. bridge # 显示或操纵 Linux 网桥 地址和设备

3. cbq

4. ctstat

5. devlink

6. genl

7. ifcfg

8. ifstat

9. [ip](https://www.yuque.com/go/doc/33221919) # 显示或操纵 routing, devices, policy routing and tunnels
10. lnstat

11. nstat

12. rdma

13. routef

14. routel

15. rtacct

16. rtmon

17. rtpr

18. rtstat

19. [ss](https://www.yuque.com/go/doc/33221911) # 转存 Socket 信息
20. tipc #&#x20;

21. tc # 实现 [TC 模块](https://www.yuque.com/go/doc/34380573)进行流量控制的程序

# 配置

**/etc/iproute2/\* **#
