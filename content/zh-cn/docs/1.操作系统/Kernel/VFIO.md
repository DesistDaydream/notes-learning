---
title: VFIO
weight: 20
tags:
  - iommu
  - dma
---

# 概述

> 参考：
>
> - [Linux 内核文档，驱动 API - VFIO](https://docs.kernel.org/driver-api/vfio.html)

许多现代系统现在都提供 [DMA](/docs/1.操作系统/DMA.md) 与中断重映射设施，帮助确保 I/O 设备在其被分配的边界内活动。这类设施包括带 AMD-Vi 和 Intel VT-d 的 x86 硬件、带 Partitionable Endpoints（PE）的 POWER 系统，以及像 Freescale PAMU 这样的嵌入式 PowerPC 系统。

**Virtual Function I/O(简称 VFIO)** 驱动程序是一个与 IOMMU 设备无关的框架，用于在受 IOMMU 保护的安全环境中向**用户空间**提供**直接访问设备**的能力。换句话说，这允许安全非特权用户空间驱动程序。

> [!Tip] 人话
> 可以让用户空间的进程直接控制物理硬件设备，而不用经过内核。比如 [DPDK](/docs/4.数据通信/DPDK/DPDK.md) 可以通过 vfio 模块，让使用 DPDK 的程序略过内核直接控制网卡，进而避免了流量过大导致的内核软中断过高的问题。

我们为什么要这样做？虚拟机在配置为追求尽可能高的 I/O 性能时，常会采用**直接访问设备**的方式（i.e. device assignment(设备分配)，把宿主机的硬件设备设备直接分配给虚拟机）。从设备和宿主的视角看，这相当于把该虚拟机变成了一个用户态驱动，换来延迟显著降低、带宽更高，以及可以直接复用裸机设备驱动。

某些应用，尤其是高性能计算领域的应用（e.g. DPDK），同样能受益于用户空间中低开销的直接设备访问。这类应用包括网络适配器（往往不是基于 TCP/IP 的）和计算加速器。在 VFIO 出现之前，这些驱动要么必须走完整的开发流程、变成正式的上游驱动，要么在树外维护，要么改用 UIO 框架，而 UIO 没有 IOMMU 保护的概念、中断支持有限，并且访问 PCI 配置空间这类资源需要 root 权限。

VFIO 驱动框架的意图是把这些做法统一起来：一方面取代 [KVM](/docs/10.云原生/Virtualization%20implementation/KVM_QEMU/KVM_QEMU.md) 中专用于 PCI 的设备分配代码，另一方面提供比 UIO 更安全、功能更丰富的用户态驱动环境。

# 最佳实践

要想使用 VFIO 必须在 BIOS 中开启 #IOMMU 。

开启 IOMMU 后，可以直接添加 PCI 设备（xml 中 hostdev.devices 字段，`mode="subsystem" type="pci"`）。若不开启就添加 PCI 设备，启动虚拟机时会报错：`host doesn't support passthrough of host PCI devices`
