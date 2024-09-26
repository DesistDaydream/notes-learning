---
title: Block
linkTitle: Block
date: 2024-09-26T17:31
weight: 20
---

# 概述

> 参考：
>
> -

**Block(块) 设备**

https://www.kernel.org/doc/Documentation/iostats.txt,

https://www.kernel.org/doc/Documentation/block/stat.txt

https://www.kernel.org/doc/Documentation/ABI/testing/procfs-diskstats

# Sysfs 中的信息

**/sys/block/\<BLOCK>/queue/**

- **./rotational** # 块设备旋转的类型，旋转就是 HHD，不旋转就是 SSD，非常形象生动得比喻磁盘使用的情况~哈哈。`0 表示 SSD`，`1 表示 HDD`
  - 注意：如果磁盘已经被做了 Raid，那么这个值将会一直都是 1。这个说法忘记了出处，找到后补充。

**/sys/block/\<BLOCK>/stat**

stat 文件提供了有关块设备 \<BLOCK> 状态的多项统计信息。共 17 个字段

https://github.com/torvalds/linux/blob/master/Documentation/block/stat.rst

