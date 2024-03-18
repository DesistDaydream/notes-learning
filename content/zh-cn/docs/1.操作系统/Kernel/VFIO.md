---
title: VFIO
linkTitle: VFIO
date: 2024-03-18T15:29
weight: 20
---

# 概述

> 参考：
> 
> - [Linux 内核文档，驱动 API - VFIO](https://docs.kernel.org/driver-api/vfio.html)

**Virtual Function I/O(简称 VFIO)** 驱动程序是一个与 IOMMU/设备无关的框架，用于在受 IOMMU 保护的安全环境中公开对**用户空间**的**直接设备**访问。换句话说，这允许安全非特权用户空间驱动程序。

人话：可以让用户空间的进程直接控制物理硬件设备，而不用经过内核。比如 [DPDK](docs/4.数据通信/DPDK.md) 可以通过 vfio 模块，让使用 DPDK 的程序掠过内核直接控制网卡，进而避免了流量过大导致的内核软中断过高的问题。
