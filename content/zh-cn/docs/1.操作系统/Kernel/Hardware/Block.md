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

**Block(块) 设备** 包含的信息常见于如下几处

- [Sysfs 中的磁盘信息](#Sysfs%20中的磁盘信息)
- [Proc 中的磁盘信息](#Proc%20中的磁盘信息)

## I/O 时间

**I/O 时间** 指 Block(块设备、磁盘、硬盘、etc.) 在执行 [I/O](/docs/0.计算机/I_O.md) 操作上花费的时间。

在 [`/sys/block/<DEV>/stat`](https://github.com/torvalds/linux/blob/master/Documentation/block/stat.rst) 文件中，io_ticks 字段记录了毫秒级的 磁盘 I/O 时间。如果磁盘在 1 秒内持续执行 I/O 操作，那么在 1 秒钟后，io_ticks 的值是 1000；如果 1 秒内磁盘一直空闲，那么 io_ticks 的值是 0。也就是说，io_ticks 的值每秒最多增加 1000。

> [!Note]
> 在 [Node Exporter](/docs/6.可观测性/Metrics/Instrumenting/Node%20Exporter.md) 源码中，[这里](https://github.com/prometheus/node_exporter/blob/v1.6.1/collector/diskstats_linux.go#L320) 可以看到 `stats.IOsTotalTicks` 对应 `diskstatsCollector.descs[10]`(i.e. node_disk_io_time_seconds_total 指标)。而 `IOsTotalTics` 对应到 [prometheus/procfs 项目，blockdevice/stats.go 中的 IOStats 结构体得 IOsTotalTicks 属性](https://github.com/prometheus/procfs/blob/v0.15.1/blockdevice/stats.go#L61)。这些结构体的信息来源遵循以下几个内核文档的说明
>
> - https://www.kernel.org/doc/Documentation/iostats.txt,
> - https://www.kernel.org/doc/Documentation/block/stat.txt
> - https://www.kernel.org/doc/Documentation/ABI/testing/procfs-diskstats

一般情况下，io_ticks 的值可以当作磁盘的使用率，比如计算某区间时间中，io_ticks 每秒的变化率（基于上面的逻辑，这个变化率一定是 0 到 1 之间的小数）。比如我们统计 1 分钟时间 io_ticks 增加了 60000，那说明这一分钟的时间中，磁盘一直在执行 I/O，i.e. 使用率是 100%

### io_ticks 与 cpu 中 wa 的区别

Claude AI 回答：

top 命令中的 wa（wait）。"wa" 代表 I/O wait 时间，表示 CPU 等待 I/O 操作完成的时间百分比。具体来说：

- 它表示 CPU 空闲（idle）且有未完成的 I/O 请求的时间比例。
- 高 wa 值通常表明系统有 I/O 瓶颈，可能是磁盘 I/O，也可能是网络 I/O。
- wa 是从 CPU 的角度来看待 I/O 等待时间。

磁盘 I/O 的 io_ticks。io_ticks 是 /proc/diskstats 文件中的一个字段，表示磁盘花在处理 I/O 请求上的时间（以毫秒为单位）。具体来说：

- 它直接反映了磁盘设备处理 I/O 请求所花费的时间。
- 这个值包括了所有类型的 I/O 操作（读和写）。
- io_ticks 是从磁盘设备的角度来看待 I/O 时间。

联系：

- 两者都与 I/O 操作的时间有关。
- 高 wa 值通常会对应较高的 io_ticks，因为 CPU 等待的 I/O 操作很可能就是磁盘 I/O。

区别：

- wa 是 CPU 时间的一个比例，而 io_ticks 是绝对时间值。
- wa 包括所有类型的 I/O 等待（包括网络 I/O），而 io_ticks 专指特定磁盘设备的 I/O 时间。
- 多核系统中，wa 可能会被平均到所有 CPU 上，而 io_ticks 是针对特定磁盘的。

实际应用：

```bash
# 查看 CPU 的 wa 百分比
top -bn1 | grep "Cpu(s)" | sed "s/.*, *\([0-9.]*\)%* wa.*/\1/"

# 查看特定磁盘的 io_ticks
cat /proc/diskstats | awk '$3 == "sda" {print $13}'

# 使用 iostat 查看更详细的 I/O 统计信息
iostat -x 1
```

通过这些命令，您可以同时监控 wa 和 io_ticks，从而更全面地了解系统的 I/O 性能状况。

相关性：

- 如果看到高 wa 值，但 io_ticks 并不高，可能说明 I/O 瓶颈不在磁盘，而在网络或其他设备。
- 如果 wa 和 io_ticks 都很高，这通常表明系统确实存在磁盘 I/O 瓶颈。

总的来说，虽然 wa 和 io_ticks 有关联，但它们提供了不同层面的 I/O 性能视角。结合使用这两个指标，可以更准确地诊断系统的 I/O 性能问题。

# Sysfs 中的磁盘信息

**/sys/block/\<BLOCK>/queue/**

- **./rotational** # 块设备旋转的类型，旋转就是 HHD，不旋转就是 SSD，非常形象生动得比喻磁盘使用的情况~哈哈。`0 表示 SSD`，`1 表示 HDD`
  - 注意：如果磁盘已经被做了 Raid，那么这个值将会一直都是 1。这个说法忘记了出处，找到后补充。

**/sys/block/\<BLOCK>/stat** # stat 文件提供了有关块设备 \<BLOCK> 状态的多项统计信息，共 17 个字段。内容详见上文

- https://github.com/torvalds/linux/blob/master/Documentation/block/stat.rst
- https://www.kernel.org/doc/Documentation/iostats.txt
- https://www.kernel.org/doc/Documentation/block/stat.txt
- https://www.kernel.org/doc/Documentation/ABI/testing/procfs-diskstats

| 名称              | 单位           | 描述                                                    |
| --------------- | ------------ | ----------------------------------------------------- |
| read I/Os       | requests     | 已处理的读 I/O 总数                                          |
| read merges     | requests     | number of read I/Os merged with in-queue I/O          |
| read sectors    | sectors      | number of sectors read                                |
| read ticks      | milliseconds | total wait time for read requests                     |
| write I/Os      | requests     | 已处理的写 I/O 总数                                          |
| write merges    | requests     | number of write I/Os merged with in-queue I/O         |
| write sectors   | sectors      | number of sectors written                             |
| write ticks     | milliseconds | total wait time for write requests                    |
| in_flight       | requests     | number of I/Os currently in flight                    |
| io_ticks        | milliseconds | 该块设备处于活动状态（执行 I/O 操作）的总时间。i.e. [磁盘 I/O 时间](#I/O%20时间) |
| time_in_queue   | milliseconds | total wait time for all requests                      |
| discard I/Os    | requests     | number of discard I/Os processed                      |
| discard merges  | requests     | number of discard I/Os merged with in-queue I/O       |
| discard sectors | sectors      | number of sectors discarded                           |
| discard ticks   | milliseconds | total wait time for discard requests                  |
| flush I/Os      | requests     | number of flush I/Os processed                        |
| flush ticks     | milliseconds | total wait time for flush requests                    |

# Proc 中的磁盘信息

**/proc/diskstats**

对应 /sys/block/\<BLOCK>/stat 的内容

