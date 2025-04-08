---
title: "Memory Info"
linkTitle: "Memory Info"
weight: 20
---

# 概述

> 参考：
>
> - 

# /proc/meminfo

> 参考：
>
> - [Manual(手册)，proc_meminfo(5)](https://man7.org/linux/man-pages/man5/proc_meminfo.5.html)
> - [GitHub 项目，torvalds/linux - Documentation/filesystems/proc.rst#meminfo](https://github.com/torvalds/linux/blob/v5.19/Documentation/filesystems/proc.rst#meminfo)
> - [RedHat 官方给的解释](https://access.redhat.com/solutions/406773)

该文件报告有关系统上 [Memory](/docs/1.操作系统/Kernel/Memory/Memory.md) 使用情况的统计信息。 free 命令使用该文件来报告系统上的可用内存和已使用内存（物理内存和交换内存）以及内核使用的共享内存和缓冲区的数量。该文件是以 `:` 符号分割的 **Key/Value pair(键/值对)** 格式。可用参数及其详解如下：

### MemTotal

总可用 Memory。即.物理 RAM 减去一些保留的 bits 和内核二进制代码所用的量

### MemFree

空闲的 Memory。LowFree 与 HighFree 两个参数的值的和

### MemAvailable

可用的 Memory。估算值，估计有多少内存可用于启动新的应用程序

### Buffers 与 Cached

详见：《[Memory 的缓存](/docs/1.操作系统/Kernel/Memory/Memory%20的缓存.md)》

### Active

最近使用过的 Memory。除非必要，否则通常不会回收。

### Inactive

最近使用比较收的 Memory。这些内存会被优先回收。

### Slab

内核数据结构缓存。dentry、inode_cache 等

### SReclaimable

Slab Reclaimable。Slab 的一部分，可以被 reclaimed(回收)。例如 dentry、inode 的缓存等等。

### SUnreclaim

Slab UnReclaim。Slab 的一部分，不可以被 reclaimed(回收)。即使内存有压力也无法回收

### CommitLimit

提交限制。当前可以分配的内存上限。只有当 [/proc/sys/vm/overcommit_memory](net(网络相关参数).md Kernel/Kernel 参数/net(网络相关参数).md) 的参数值为 2 的时候，该限制才生效。这个上限是指当程序向系统申请内存时，如果申请的内存加上现在已经分配的内存，超过了 commitlimit 的值，则该申请将会失败。

该值通过如下公式：

`CommitLimit = (total_RAM - total_huge_TLB) * overcommit_ratio / 100 + total_swap`

- totaol_RAM # 系统内存总量(就是物理内存)
- total_huge_TLB # 为 huge pages 保留的内存量，一般没有保留，都是 0
- overcommit_ratio # /proc/sys/vm/overcommit_ratio 内核参数的值。
- total_swap # swap 空间的总量

### Committed_AS

> Allocated Size(已经分配的大小，简称 AS)

当前已经分配的内存总量。注意：不是正在使用的，而是已经分配的。

当 overcommit_memory 参数的值为 2 时，该值不能超过 CommitLimit 的值。其余时候该值可以无限大。

### VmallocXXX

参考：<https://zhuanlan.zhihu.com/p/77827102>

- VmallocTotal
- VmallocUsed
- VmallocChunk

注意：**VmallocTotal 会非常大，这是正常的**

### Percpu

### HardwareCorrupted

### AnonHugePages

### CmaXXX

- CmaTotal
- CmaFree

### HugePagesXXX

[Huge Pages](/docs/1.操作系统/Kernel/Memory/Huge%20Pages.md) 相关信息

- HugePages_Total
- HugePages_Free
- HugePages_Rsvd
- HugePages_Surp
- Hugepagesize

### DirectMapXXX

- DirectMap4k
- DirectMap2M
- DirectMap1G