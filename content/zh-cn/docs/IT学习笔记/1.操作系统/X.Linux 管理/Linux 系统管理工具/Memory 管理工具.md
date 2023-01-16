---
title: Memory 管理工具
---

# 概述

# 概述

> ## 参考：

# 查看 Memory 的使用情况

我们可以通过多种方式查看 Memory 信息。

## /proc/memory 文件

> 参考：
> - [RedHat 官方给的解释](https://access.redhat.com/solutions/406773)

该文件报告有关系统上内存使用情况的统计信息。 free 命令使用该文件来报告系统上的可用内存和已使用内存（物理内存和交换内存）以及内核使用的共享内存和缓冲区的数量。该文件是以 `:` 符号分割的 **Key/Value pair(键/值对)** 格式。可用参数及其详解如下：

### MemTotal

总可用 Memory。即.物理 RAM 减去一些保留的 bits 和内核二进制代码所用的量

### MemFree

空闲的 Memory。LowFree 与 HighFree 两个参数的值的和

### MemAvailable

可用的 Memory。估算值，估计有多少内存可用于启动新的应用程序

### Buffers 与 Cached

详见：《[Memory 的缓存机制](✏IT 学习笔记/📄1.操作系统/2.Kernel(内核)/5.Memory%20 管理/Memory%20 的缓存机制.md 管理/Memory 的缓存机制.md)》

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

## free 命令

```bash
~]# free -h
              total        used        free      shared  buff/cache   available
Mem:          3.8Gi       846Mi       506Mi       1.0Mi       2.5Gi       2.9Gi
Swap:            0B          0B          0B
```

Mem：设备上的真实内存

- total # 总计。该设备的总内存大小
- used # 已使用的。linux 对内存的使用量
- free # 空闲的。还剩多少内存可用
- shared # 共享内存
- buff # 缓冲区(buffer)。保存一些将要写入到硬盘中的数据。
- cache # 缓存。从硬盘中读出的数据存放到内存中，以便再次读取相同数据时速度更快。
- availabel # 可用的。free+buff/cache 合起来就是可用的。

## free 命令 与 /proc/meminfo 文件中信息的对应关系。

| free 命令输出     | `/proc/meminfo`文件的字段                              |
| ----------------- | ------------------------------------------------------ |
| `Mem: total`      | `MemTotal`                                             |
| `Mem: used`       | `MemTotal - MemFree - Buffers - Cached - SReclaimable` |
| `Mem: free`       | `MemFree`                                              |
| `Mem: shared`     | `Shmem`                                                |
| `Mem: buff/cache` | `Buffers + Cached + Slab`                              |
| `Mem:available`   | `MemAvailable`                                         |
| `Swap: total`     | `SwapTotal`                                            |
| `Swap: used`      | `SwapTotal - SwapFree`                                 |
| `Swap: free`      | `SwapFree`                                             |

# 一个可以消耗 Linux 内存的 Shell 脚本

```bash
#!/bin/bash
mkdir /tmp/memory
mount -t tmpfs -o size=1024M tmpfs /tmp/memory
dd if=/dev/zero of=/tmp/memory/block
sleep 3600
rm /tmp/memory/block
umount /tmp/memory
rmdir /tmp/memory
```
