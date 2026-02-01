---
title: PXE
linkTitle: PXE
weight: 54
---

# 概述

> 参考：
>
> - [Wiki, Preboot Execution Environment](https://en.wikipedia.org/wiki/Preboot_Execution_Environment)

**Preboot eXecution Environment(预启动执行环境，简称 PXE)** 提供了一种使用网络接口（Network Interface）启动计算机的机制。这种机制让计算机的启动可以不依赖本地数据存储设备（如硬盘）或本地已安装的操作系统。

在服务器开机时，可以使用 PXE 进行启动，该设备的 PXE 作为一个客户端，首先请求 DHCP，在获取到网络参数后，再在广播域里请求 PXE 类型的服务，来引导安装操作系统。一般来说都是通过 TFTP 来进行远程系统文件传输，然后再自动通过传输过来的文件自动进行系统安装
