---
title: "perf 性能分析工具"
weight: 1
---

# 概述

> 参考：
>
> - [GitHub 项目，torvalds/linux - /tools/perf](https://github.com/torvalds/linux/tree/master/tools/perf)
> - [Kernel Wiki，perf](https://perf.wiki.kernel.org/index.php/Main_Page)
> - perf 事件列表中的内核 pmu 事件是什么
>   - https://unix.stackexchange.com/questions/326621/what-are-kernel-pmu-event-s-in-perf-events-list
>   - https://qastack.cn/unix/326621/what-are-kernel-pmu-event-s-in-perf-events-list
> - [brendangregg 博客，perf](https://www.brendangregg.com/perf.html)

**Linux Performance Events(Linux 性能事件，简称 LPE)** 是用来分析 [Linux Kernel](/docs/1.操作系统/Kernel/Linux%20Kernel/Linux%20Kernel.md) 性能的工具，通常称为 **perf**。perf 随 Kernel 2.6+ 一同发布。通过它，应用程序可以利用 PMU，tracepoint 和内核中的特殊计数器来进行性能统计。它不但可以分析指定应用程序的性能问题 (per thread)，也可以用来分析内核的性能问题，当然也可以同时分析应用代码和内核，从而全面理解应用程序中的性能瓶颈。

perf 主要是通过 **Tracing(追踪)** 的方式来实现性能数据的采集。

perf 和其他调试工具一样，需要 **symbol(符号信息)**。它们用于将内存地址转换为函数和变量名称，以便我们人类可以读取它们。如果没有符号，您将看到代表所分析的内存地址的十六进制数字。

> [!Note]
> perf 命令甚至有时候可以跟 [strace 工具](/docs/1.操作系统/Linux%20管理/Linux%20系统管理工具/strace%20工具.md) 实现类似的效果，比如 `perf stat -e syscalls:*` 统计系统调用的计数，就很像 `strace -c XX`

## Events

使用 `perf list` 命令可以列出可以分析的所有事件。

TODO: perf 可用的 Events 列表是从如何获取到的？

- perf_event_open() 系统调用？
- /sys/kernel/debug/tracing/events/ 目录？是特定于 tracepoint 类型事件的？
- 由于 perf 本身就是与 Linux 内核强耦合的工具，所以获取 Events 应该也是通过某种方式动态获取的？

perf_event_open [系统调用](/docs/1.操作系统/Kernel/System%20Call/System%20Call.md)用以设置性能监控，其中 **perf_event_attr** 参数（源码: [include/uapi/linux/perf_event.h - struct perf_event_attr {}](https://github.com/torvalds/linux/blob/v6.10/include/uapi/linux/perf_event.h#L389)） 是一个结构体，为正在创建的 Event 提供详细的配置信息，[perf_event_open 的 Manual]([perf_event_open](https://man7.org/linux/man-pages/man2/perf_event_open.2.html)) 中列出了所有 Events 的类型：

- PERF_TYPE_HARDWARE
- PERF_TYPE_SOFTWARE
- PERF_TYPE_TRACEPOINT
- PERF_TYPE_HW_CACHE
- PERF_TYPE_RAW
- PERF_TYPE_BREAKPOINT (since Linux 2.6.33)
- dynamic PMU
- kprobe and uprobe (since Linux 4.17)

通过 `perf list` 命令显示出来的分类：

- hw
- sw
- cache
- tracepoint
- pmu
- sdt
- event_glob

[brendangregg 对 Events 的分类](https://www.brendangregg.com/perf.html#Events)：

- **Hardware Events** # CPU性能监控计数器。
- **Software Events** # 这些是基于内核计数器的低级事件。e.g. CPU migrations, minor faults, major faults, etc.
- **Kernel Tracepoint Events** # 这是静态的内核级检测点，它们被硬编码在内核中有趣且合乎逻辑的位置。（TODO: /sys/kernel/debug/tracing/events/ 目录？）
- **User Statically-Defined Tracing (USDT)** # 这些是用于用户级程序和应用程序的静态跟踪点。
- **Dynamic Tracing** # 软件可以动态检测，在任何位置创建事件。对于内核软件，这使用 kprobes 框架。对于用户级软件，uprobes。
- **Timed Profiling**: Snapshots can be collected at an arbitrary frequency, using perf record -F_Hz_. This is commonly used for CPU usage profiling, and works by creating custom timed interrupt events.

### Tracepoint 事件

> 参考：
>
> - [内核官方文档，核心 API - tracepoint](https://docs.kernel.org/core-api/tracepoint.html)
> - https://www.brendangregg.com/perf.html#Tracepoints

Tracepoint 事件可以在 Linux Kernel 代码的 [include/trace/events/](https://github.com/torvalds/linux/blob/v6.10/include/trace/events/) 处找到所有注册的 Tracepoint，这里还可以从代码注释中看到每个 Tracepoint 事件的详细描述

比如，irq 相关 Tracepoint，可以在 include/trace/events/irq.h 文件中找到

```c
TRACE_EVENT(irq_handler_entry,......)
TRACE_EVENT(irq_handler_exit,......)
DEFINE_EVENT(softirq, softirq_entry,......)
DEFINE_EVENT(softirq, softirq_exit,......)
......
```

这与 `perf list irq:*` 命令列出的内容一致

```bash
~]# perf list irq:*

List of pre-defined events (to be used in -e):

  irq:irq_handler_entry                              [Tracepoint event]
  irq:irq_handler_exit                               [Tracepoint event]
  irq:softirq_entry                                  [Tracepoint event]
  irq:softirq_exit                                   [Tracepoint event]
  irq:softirq_raise                                  [Tracepoint event]
```

在内核官方文档里，可以可以从核心 API 里找到 tracepoint 的描述

通过如下 `perf list` 命令，可以查看每种 Tracepoint 类型下可用的 Tracepoint 总数

```bash
perf list | awk -F: '/Tracepoint event/ { lib[$1]++ } END { for (l in lib) { printf " %-16.16s %d\n", l, lib[l] } }' | sort | column
```

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

# Syntax(语法)

> 参考：
>
> - [Manual(手册), perf(1)](https://www.man7.org/linux/man-pages/man1/perf.1.html)
> - [官方手册](https://perf.wiki.kernel.org/index.php/Latest_Manual_Pages)

**perf \[OPTIONS] COMMAND \[ARGS]**

perf 主要由多个子命令来提供常用功能

**COMMAND**

在源码 [tools/perf/perf.c](https://github.com/torvalds/linux/blob/v6.10/tools/perf/perf.c#L51) 处有所有子命令的列表

- **list** # 列出所有可用事件
  - 命令入口: [tools/perf/builtin-list.c - cmd_list() 函数](https://github.com/torvalds/linux/blob/v6.10/tools/perf/builtin-list.c#L505)
- **stat** # 收集性能计数器统计信息

通用 OPTIONS（下面这些 OPTIONS 不是所有命令都支持，但是有些命令的 OPTIONS 逻辑是类似的，这些命令都是 perf 进行性能监控的常用命令，比如 stat、record、top、etc.）

- **-a, --all-cpus** # 从所有 CPU 采集数据
- **-e, --event \<EVENT>** # Event 选择器。让 perf 指定的事件进行操作。可以通过 `perf list` 命令列出所有可用的 Events
- **-g** # 启动调用关系分析
- **-p, --pid \<PID>** # 指定要采集数据的进程的 PID

## record - 追踪指定的进程，并记录它的 profile 到 perf.data 文件中

> 参考：
>
> - [Manual(手册), perf-record(1)](https://man7.org/linux/man-pages/man1/perf-record.1.html)

record 子命令将会跟踪指定命令或进程，并采集运行期间的 profile，然后默认将这些数据写入到 perf.data 文件中。**profile** 这个词在这个语境中，可以理解为 **性能分析**，详见 [Flame Graphs(火焰图)](/docs/6.可观测性/性能优化/Flame%20Graphs(火焰图).md)

**perf record \[OPTIONS] \[COMMAND]**

**COMMAND** # 可以指定一个命令，以便采集指定命令运行时的性能数据。或者省略 COMMAND，则采集当前系统下的所有进程。

**OPTIONS**

EXAMPLE

- **perf record -a -g -p 5958 -- sleep 30**

## report - 读取 perf.data 文件中记录的内容并展示

## script - 读取 perf.data 文件中记录的内容，并展示追踪效果的数据

## top - 系统分析工具

> 参考：
>
> - [Manual(手册), perf-top(1)](https://www.man7.org/linux/man-pages/man1/perf-top.1.html)

可以实时显示占用 CPU 时钟最多的函数或进程。以 Symbol 为中心，显示指定 Symbol 的相关信息

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
- Object # 动态共享对象的类型。比如 `[.]` 表示用户空间的可执行程序、或动态链接库，而 `[k]` 则表示内核空间
- Symbol # 符号名。即函数名。当函数名未知时，用十六进制的地址来表示。

**perf top \[OPTIONS]**

OPTIONS

EXAMPLE

# 最佳实践

> 参考：
>
> - [brendangregg 博客，perf - 例子](https://www.brendangregg.com/perf.html#Examples)

## 生成火焰图

获取火焰图生成工具

```bash
git clone https://github.com/brendangregg/FlameGraph
cd FlameGraph
```

将 `perf record XXX` 命令抓到的 perf.data 文件中的记录转换成可读的采样记录

`perf script -i /root/tmp/perf.data`

合并调用栈信息

`./stackcollapse-perf.pl`

生成火焰图

`./flamegraph.pl`

上述命令合并一下：

```bash
perf script -i /root/tmp/perf.data | ./stackcollapse-perf.pl --all | ./flamegraph.pl > /root/tmp/flame.svg
```

## 追踪中断

`perf top -e irq:irq_handler_entry`

输出效果类似下面这样

```bash
Samples: 51K of event 'irq:irq_handler_entry', 1 Hz, Event count (approx.): 1436 lost: 0/0 drop: 0/0
Overhead  Trace output
  94.01%  irq=16 name=enp0s8
   1.88%  irq=19 name=ehci_hcd:usb1
   1.88%  irq=19 name=enp0s3
   1.11%  irq=21 name=ahci[0000:00:0d.0]
   1.11%  irq=21 name=snd_intel8x0
```

```bash
Samples: 6K of event 'irq:irq_handler_entry', 1 Hz, Event count (approx.): 303 lost: 0/0 drop: 0/0
Overhead  Trace output
  15.84%  irq=53 name=ahci[0000:00:1f.2]
  10.23%  irq=44 name=enp9s0-TxRx-1
   8.58%  irq=43 name=enp9s0-TxRx-0
   7.59%  irq=45 name=enp9s0-TxRx-2
   2.97%  irq=46 name=enp9s0-TxRx-3
   0.99%  irq=48 name=enp12s0-TxRx-0
```


