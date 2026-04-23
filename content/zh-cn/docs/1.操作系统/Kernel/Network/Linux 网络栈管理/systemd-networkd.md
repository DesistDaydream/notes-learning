---
title: systemd-networkd
linkTitle: systemd-networkd
weight: 32
---

# 概述

> 参考：
>
> - [Linux man pages, systemd-networkd(8)](https://man7.org/linux/man-pages/man8/systemd-networkd.8.html)

**systemd-networkd** 是一个用于管理网络的系统服务。检测并配置出现的网络设备，以及创建虚拟网络设备。

systemd-networkd 与 [systemd-resolved.service](/docs/1.操作系统/Kernel/Network/Linux%20DNS%20管理/systemd-resolved.service.md) 一起组成了由 [Systemd](/docs/1.操作系统/Systemd/Systemd.md) 驱动的完整网络访问能力。

# 配置

**/run/systemd/network/** # 读取网络设备配置的目录
