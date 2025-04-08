---
title: Memory 管理工具
weight: 2
---

# 概述

> 参考：

# 查看 Memory 的使用情况

我们可以通过多种方式查看 Memory 信息。

## free 命令

> 参考：
>
> - https://gitlab.com/procps-ng/procps

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

## free 命令 与 /proc/meminfo 文件中信息的对应关系

| free 命令输出         | `/proc/meminfo`文件的字段                                   |
| ----------------- | ------------------------------------------------------ |
| `Mem: total`      | `MemTotal`                                             |
| `Mem: used`       | `MemTotal - MemFree - Buffers - Cached - SReclaimable` |
| `Mem: free`       | `MemFree`                                              |
| `Mem: shared`     | `Shmem`                                                |
| `Mem: buff/cache` | `Buffers + Cached + SReclaimable`                      |
| `Mem: available`  | `MemAvailable`                                         |
| `Swap: total`     | `SwapTotal`                                            |
| `Swap: used`      | `SwapTotal - SwapFree`                                 |
| `Swap: free`      | `SwapFree`                                             |

buff/cache = Buffers + Cached + SReclaimable

从 procps 项目代码中查找 free 命令的具体逻辑 https://gitlab.com/procps-ng/procps/-/blob/v4.0.5/src/free.c?ref_type=tags#L329

```c
		if (flags & FREE_WIDE) {
			printf(" %11s", scale_size(MEMINFO_GET(mem_info, MEMINFO_MEM_BUFFERS, ul_int),
				    args.exponent, flags & FREE_SI, flags & FREE_HUMANREADABLE));
			printf(" %11s", scale_size(MEMINFO_GET(mem_info, MEMINFO_MEM_CACHED_ALL, ul_int)
				    , args.exponent, flags & FREE_SI, flags & FREE_HUMANREADABLE));
		} else {
			printf(" %11s", scale_size(MEMINFO_GET(mem_info, MEMINFO_MEM_BUFFERS, ul_int) +
				    MEMINFO_GET(mem_info, MEMINFO_MEM_CACHED_ALL, ul_int), args.exponent, flags & FREE_SI, flags & FREE_HUMANREADABLE));
		}
```

寻找 MEMINFO_MEM_CACHED_ALL 来源

[`MEMINFO_MEM_CACHED_ALL, // ul_int derived from MEM_CACHED + MEM_SLAB_RECLAIM`](https://gitlab.com/procps-ng/procps/-/blob/v4.0.5/library/include/meminfo.h#L45)

[`MEM_set(MEM_CACHED_ALL, ul_int, derived_mem_cached)`](https://gitlab.com/procps-ng/procps/-/blob/v4.0.5/library/meminfo.c#L168)

寻找 derived_mem_cached 来源

[`mHr(derived_mem_cached) = mHr(Cached) + mHr(SReclaimable);`](https://gitlab.com/procps-ng/procps/-/blob/v4.0.5/library/meminfo.c#L729)

验证代码

```bash
~]# free -w && cat /proc/meminfo | grep -i -E 'MemTotal|MemFree|MemAvailable|Buffers|Cached|Slab|SReclaimable|SUnreclaim'
               total        used        free      shared     buffers       cache   available
Mem:        30374804     1467916    25966764        2384       20108     2920016    28444040
Swap:              0           0           0
MemTotal:       30374804 kB
MemFree:        25966736 kB
MemAvailable:   28444012 kB
Buffers:           20108 kB
Cached:          1902468 kB
SwapCached:            0 kB
Slab:            1153644 kB
SReclaimable:    1017548 kB
SUnreclaim:       136096 kB
```

buffers(20108) = Buffers(20108) 

cache(2920016) = Cached(1902468) + SReclaimable(1017548)

## OPTIONS

**-h, --human** # 以人类可读的方式输出结果。i.e. 带着 Ki, Mi, Gi, etc. 这种单位

**-w, --wide** # 以宽模式输出结果。buff/cache 将会拆分成两列，分别现实其各自的内存使用情况。

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
