---
title: I/O Virtualization
linkTitle: I/O Virtualization
date: 2024-03-18T17:43
weight: 20
---

# 概述

> 参考：
> 
> - [Wiki，I/O virtualization](https://en.wikipedia.org/wiki/I/O_virtualization)

**I/O Virtualization(输入/输出虚拟化)** 是一种广义的虚拟化概念。

[VFIO](docs/1.操作系统/Kernel/VFIO.md) 是 I/O 虚拟化的一种具体实现方式，专注于将设备直通给虚拟机，让虚拟机可以直接控制物理设备，从而获得接近于物理机的性能。

