---
title: Load 高/CPU 使用率 问题及实用脚本
---

# Linux 资源&瓶颈分析概述

> 参考：
> - [公众号,CPU 飙高，系统性能问题如何排查？](https://mp.weixin.qq.com/s/fzLcAkYwKhj-9hgoVkTzaw)
> - [阿里云,ECS 运维指南之 LInux 系统诊断-找到 Linux 虚机 Load 高的元凶](https://developer.aliyun.com/topic/download?id=143)

Load Average 和 CPU 使用率 可被细分为不同的子域指标，指向不同的资源瓶颈。总体来说，指标与资源瓶颈的对应关系基本如下图所示。

注意：Load 与 CPU 使用率 之间没有必然的联系。有可能 Load 很高，而 CPU 使用率很低；也有可能 Load 很低而 CPU 使用率很高。具体原因详见 CPU 管理 与 Process 进程管理 章节中关于 Load 与 CPU 使用率的概念。简单说就是因为 Load Average 在计算时，包含了对 I/O 的统计

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/go9krg/1616164005685-f19dcd0b-9837-42cb-aeaa-a48586cf1cca.jpeg)

## Load 高 & CPU 高

这是我们最常遇到的一类情况，即 load 上涨是 CPU 负载上升导致。根据 CPU 具体资源分配表现，可分为以下几类：

**CPU sys 高**

这种情况 CPU 主要开销在于系统内核，可进一步查看上下文切换情况。

- 如果非自愿上下文切换较多，说明 CPU 抢占较为激烈，大量进程由于时间片已到等原因，被系统强制调度，进而发生的上下文切换。
- 如果自愿上下文切换较多，说明可能存在 I/O、内存等系统资源瓶颈，大量进程无法获取所需资源，导致的上下文切换。

**CPU si 高**

这种情况 CPU 大量消耗在软中断，可进一步查看软中断类型。一般而言，网络 I/O 或者线程调度引起软中断最为常见：

- NET_TX & NET_RX。NET_TX 是发送网络数据包的软中断，NET_RX 是接收网络数据包的软中断，这两种类型的软中断较高时，系统存在网络 I/O 瓶颈可能性较大。
- SCHED。SCHED 为进程调度以及负载均衡引起的中断，这种中断出现较多时，系统存在较多进程切换，一般与非自愿上下文切换高同时出现，可能存在 CPU 瓶颈。

**CPU us 高**

这种情况说明资源主要消耗在应用进程，可能引发的原因有以下几类：

- 死循环或代码中存在 CPU 密集计算。这种情况多核 CPU us 会同时上涨。
- 内存问题，导致大量 FULLGC，阻塞线程。这种情况一般只有一核 CPU us 上涨。
- 资源等待造成线程池满，连带引发 CPU 上涨。这种情况下，线程池满等异常会同时出现。

## Load 高 & CPU 低

这种情况出现的根本原因在于不可中断睡眠态(TASK_UNINTERRUPTIBLE)进程数较多，即 CPU 负载不高，但 I/O 负载较高。可进一步定位是磁盘 I/O 还是网络 I/O 导致。

# 排查策略

利用现有常用的工具，我们常用的排查策略基本如下图所示：
![](https://notes-learning.oss-cn-beijing.aliyuncs.com/go9krg/1616164005644-f386b76a-e4ac-447e-b4fc-185c5393c19f.jpeg)

从问题发现到最终定位，基本可分为四个阶段：

## 资源瓶颈定位

这一阶段通过全局性能检测工具，初步定位资源消耗异常位点。

常用的工具有：

- top、vmstat、tsar(历史)
- 中断：/proc/softirqs、/proc/interrupts
- I/O：iostat、dstat

## 热点进程定位

定位到资源瓶颈后，可进一步分析具体进程资源消耗情况，找到热点进程。

常用工具有：

- 上下文切换：pidstat -w
- CPU：pidstat -u
- I/O：iotop、pidstat -d
- 僵尸进程：ps

## 线程&进程内部资源定位

找到具体进程后，可细化分析进程内部资源开销情况。

常用工具有：

- 上下文切换：pidstat -w -p \[pid]
- CPU：pidstat -u -p \[pid]
- I/O: lsof

## 热点事件&方法分析

获取到热点线程后，我们可用 trace 或者 dump 工具，将线程反向关联，将问题范围定位到具体方法&堆栈。

常用的工具有：

- perf：Linux 自带性能分析工具，功能类似 hotmethod，基于事件采样原理，以性能事件为基础，支持针对处理器相关性能指标与操作系统相关性能指标的性能剖析。
- jstack
  - 结合 ps -Lp 或者 pidstat -p 一起使用，可初步定位热点线程。
  - 结合 zprofile-threaddump 一起使用，可统计线程分布、等锁情况，常用与线程数增加分析。
- strace：跟踪进程执行时的系统调用和所接收的信号。
- tcpdump：抓包分析，常用于网络 I/O 瓶颈定位。

# 实用脚本

## 找出系统中 Load 高时处于运行队列的进程

```bash
#!/bin/bash
LANG=C
PATH=/sbin:/usr/sbin:/bin:/usr/bin
interval=1
length=86400
for i in $(seq 1 $(expr ${length} / ${interval}));do
date
LANG=C ps -eTo stat,pid,tid,ppid,comm --no-header | sed -e 's/^ *//' |
perl -nE 'chomp;say if (m!^\S*[RD]+\S*!)'
date
cat /proc/loadavg
echo -e "\n"
sleep ${interval}
done
```

## 查 CPU 使用率高的线程

```bash
#!/bin/bash
LANG=C
PATH=/sbin:/usr/sbin:/bin:/usr/bin
interval=1
length=86400
for i in $(seq 1 $(expr ${length} / ${interval}));do
date
LANG=C ps -eT -o%cpu,pid,tid,ppid,comm | grep -v CPU | sort -n -r | head -20
date
LANG=C cat /proc/loadavg
{ LANG=C ps -eT -o%cpu,pid,tid,ppid,comm | sed -e 's/^ *//' | tr -s ' ' |
grep -v CPU | sort -n -r | cut -d ' ' -f 1 | xargs -I{} echo -n "{} + " &&
echo '0'; } | bc -l
sleep ${interval}
done
fuser -k $0
```
