---
title: DMA
created: 2026-10-18T23:19
weight: 100
tags:
  - iommu
---

# 概述

> 参考：
>
> - [Wiki, Direct memory access](https://en.wikipedia.org/wiki/Direct_memory_access)

**Direct Memory Access(直接内存访问，简称 DMA)** 是一种让硬件子系统直接访问内存的方式

通过 DMA 可以让应用程序直接通过内存管理硬件。没有 DMA 之前，应用程序需要通过 CPU 来读取硬件信息，但是这会占用大量 CPU 时间，导致系统整体性能无法提高。

使用 DMA 的方式，各种应用程序就可以绕过 CPU 直接访问内存中的硬件数据，e.g. [DPDK](/docs/4.数据通信/DPDK/DPDK.md)、etc.

# IOMMU

> 参考：
>
> - [Wiki, Input-Output Memory Management Unit](https://en.wikipedia.org/wiki/Input-output_memory_management_unit)

**Input-Output Memory Management Unit(输入输出内存管理单元，简称 IOMMU)**

需要在 BIOS 中控制 IOMMU 的开关
