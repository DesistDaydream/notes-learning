---
title: "perf 性能分析工具"
weight: 1
---

# 概述

> 参考：
>
> - [GitHub 项目，torvalds/linux - /tools/perf](https://github.com/torvalds/linux/tree/master/tools/perf)
> - [Kernel Wiki，perf](https://perf.wiki.kernel.org/index.php/Main_Page)
> - [PMU 是什么？](https://qastack.cn/unix/326621/what-are-kernel-pmu-event-s-in-perf-events-list)

**Linux Performance Events(Linux 性能事件，简称 LPE)** 是用来分析 Linux 性能的工具，通常称为 **perf**。perf 随 Kernel 2.6+ 一同发布。通过它，应用程序可以利用 PMU，tracepoint 和内核中的特殊计数器来进行性能统计。它不但可以分析指定应用程序的性能问题 (per thread)，也可以用来分析内核的性能问题，当然也可以同时分析应用代码和内核，从而全面理解应用程序中的性能瓶颈。

perf 主要是通过 **Tracing(追踪)** 的方式来实现性能数据的采集。

# perf 安装

**Ubuntu**

```bash
sudo apt install -y linux-tools-generic linux-tools-$(uname -r)
```

> 注意：linux-tools-generic 会安装 linux-tools-common 包，perf 二进制文件在这个包中。

**CentOS**

```bash
yum install -y perf
```

# perf 工具使用详解

**perf \[OPTIONS] COMMAND \[ARGS]**

perf 主要由多个子命令来提供常用功能

## record - 追踪指定的进程，并记录它的 profile 到 perf.data 文件中

> 参考：
>
> - [Manual(手册)，perf-record(1)](https://man7.org/linux/man-pages/man1/perf-record.1.html)

record 子命令将会跟踪指定命令或进程，并采集运行期间的 profile，然后默认将这些数据写入到 perf.data 文件中。**profile** 这个词在这个语境中，可以理解为 **性能分析**，详见 [Flame Graphs(火焰图)](/docs/6.可观测性/性能优化/Flame%20Graphs(火焰图).md)

**perf record \[OPTIONS] \[COMMAND]**

**COMMAND** # 可以指定一个命令，以便采集指定命令运行时的性能数据。或者省略 COMMAND，则采集当前系统下的所有进程。

**OPTIONS**

- **-a, --all-cpus** # 从所有 CPU 采集数据
- **-g** # 启动调用关系分析
- **-p, --pid \<PID>** # 指定要采集数据的进程的 PID

EXAMPLE

- **perf record -a -g -p 5958 -- sleep 30**

## report - 读取 perf.data 文件中记录的内容并展示

## script - 读取 perf.data 文件中记录的内容，并展示追踪效果的数据

## top - 系统分析工具

可以实时显示占用 CPU 时钟最多得函数或进程。以 Symbol 为中心，显示指定 Symbol 的相关信息

```bash
Samples: 2K of event 'cpu-clock:pppH', 4000 Hz, Event count (approx.): 317933941 lost: 0/0 drop: 0/0
Overhead  Shared Object            Symbol
16.96%  perf                     [.] __symbols__insert
7.12%   perf                     [.] rb_next
3.92%   [kernel]                 [k] kallsyms_expand_symbol.constprop.1
```

第一行：

- Samples # 样本数
- event # 事件类型
- event count # 事件总数

第二行

- Overhead # Symbol 的性能事件在所有样本中的百分比
- Shared # Symbol 所在的动态共享对象(Dynamic Shared Object)。如 内核、进程名、动态链接库名、内核模块等等
- Object # 动态共享对象的类型。比如 \[.] 表示用户空间的可执行程序、或动态链接库，而 \[k] 则表示内核空间
- Symbol # 符号名。即函数名。当函数名未知时，用十六进制的地址来表示。

**perf top \[OPTIONS]**

OPTIONS

- **-g** # 开启调用关系分析
- **-p \<PID>** # 分析指定进程的事件，PID 可以是使用 逗号 分隔的多个 PID

EXAMPLE

# 生成火焰图

获取火焰图生成工具

```bash
git clone https://github.com/brendangregg/FlameGraph
cd FlameGraph
```

将 perf record 抓到的记录转换成可读的采样记录

`perf script -i /root/tmp/perf.data`

合并调用栈信息

`./stackcollapse-perf.pl`

生成火焰图

`./flamegraph.pl`

上述命令合并一下：

```bash
perf script -i /root/tmp/perf.data | ./stackcollapse-perf.pl --all | ./flamegraph.pl > /root/tmp/flame.svg
```
