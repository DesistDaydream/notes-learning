---
title: proc
linkTitle: proc
weight: 20
---

# 概述

> 参考：
>
> - [Manual(手册)，proc(5)](https://man7.org/linux/man-pages/man5/proc.5.html)
> - [GitHub 项目，torvalds/linux - Documentation/filesystems/proc.rst](https://github.com/torvalds/linux/blob/master/Documentation/filesystems/proc.rst)
>   - https://www.kernel.org/doc/html/latest/filesystems/proc.html

**process information pseudo-filesystem(进程信息伪文件系统，简称 proc)**， 提供了内核数据结构的接口。`一般挂载到 /proc 目录`。一般情况是由操作系统自动挂载的，也可以通过`mount -t proc proc /proc`命令手动挂载。proc 文件系统中的大多数文件都是只读的，但是有些文件是可写的，用于改变内核参数。

proc 文件系统不用于存储。其主要目的是为硬件，内存，运行的进程和其他系统组件提供基于文件的接口。通过查看相应的 /proc 文件，可以检索许多系统组件上的实时信息。/proc 中的某些文件也可以（由用户和应用程序）操纵以配置内核。

# /proc/PID/ - 每个进程自己的独立信息

**每个进程在 /proc 下有一个名为自己进程号的目录，该目录记载了该进程相关的 proc 信息。**

这些目录的详细用处详见 [Process info](/docs/1.操作系统/Kernel/Process/Process%20info.md)

# /proc/cmdline - 引导系统时，传递给内核的参数

通常通过引导管理器（如 lilo（8）或 grub（8））完成。

# /proc/cpuinfo - CPU 信息

如 cpu 的类型、制造商、 型号和性能等。

# /proc/devices - 当前运行的核心配置的设备驱动的列表

# /proc/dma - 显示当前使用的 dma 通道

# /proc/filesystems - 内核可以支持的文件系统列表

# /proc/interrupts - 系统中断统计信息

详见：[Interrupts(中断)](/docs/1.操作系统/Kernel/CPU/Interrupts(中断)/Interrupts(中断).md)

这用于记录每个 IO 设备每个 CPU 的中断数。从 Linux 2.6.24 开始，至少对于 i386 和 x86-64 体系结构，这还包括系统内部的中断（即与设备本身不相关的中断），例如 NMI（不可屏蔽中断），LOC（本地）。计时器中断），而对于 SMP 系统，则是 TLB（TLB 刷新中断），RES（重新安排中断），CAL（远程功能调用中断）等。格式非常容易读取，以 ASCII 格式完成。

# /proc/ioports - 当前使用的 i/o 端口

# /proc/kcore - 系统物理内存映像

与物理内存大小完全一样，然而实际上没有 占用这么多内存；它仅仅是在程序访问它时才被创建。(注意：除非你把它拷贝到什么地方，否则/proc 下没有任何东西占用任何磁盘空间。)

# /proc/kmsg - 内核输出的信息

这些信息也会被送到 syslog。dmesg 命令获取该文件中的内容并展示

# /proc/loadavg - 系统 load average 信息

# /proc/meminfo - 系统上内存使用情况的统计信息

详见 [Memory Info](/docs/1.操作系统/Kernel/Memory/Memory%20Info.md) 章节

# /proc/modules

存放当前加载了哪些核心模块信息。

# /proc/net/ - 网络层的信息

**软连接文件，被连接到 self/net 目录。主要是链接到的进程的 network namespace 的信息**

详见 [Linux 网络栈管理](/docs/1.操作系统/Kernel/Network/Linux%20网络栈管理/Linux%20网络栈管理.md#/proc/net) 中的 /proc/net 章节

# /proc/softirqs - **软中断统计信息**

# /proc/self

**当某一进程访问此软连接时，该软连接将指向该进程自己的 /proc/PID/ 目录。**

# /proc/sys/

**sys 目录是可写的，可以通过它来访问或修改内核的参数。详见 [Kernel 参数](/docs/1.操作系统/Kernel/Linux%20Kernel/Kernel%20参数/Kernel%20参数.md) 文章**

# /proc/stat

https://man7.org/linux/man-pages/man5/proc_stat.5.html

**系统的不同状态，例如，系统启动后页面发生错误的次数。**

该文件包含系统启动以来的很多系统和内核的统计信息，平时大家比较关心的比如包括 CPU 运行情况、中断情况、启动时间、上线文切换次数、运行中的进程等信息都在其中。

一、文件全貌

```bash
# Linux下查看/proc/stat的具体信息如下
[root@WSC-31-2 ~]# cat /proc/stat
cpu  60382 1 80032 198934063 2349 0 109 0 0 0
cpu0 2405 0 2084 4140924 682 0 6 0 0 0
...  # 此处较多信息，简化之
cpu47 200 0 134 4147222 10 0 0 0 0 0
intr 33622492 64 ... 0 0 0 0 # 此处较多信息，简化之
ctxt 68533835
btime 1528905555
processes 318904
procs_running 1
procs_blocked 0
softirq 16567860 0 3850777 8555 5448802 116727 0 1 3577293 1290 3564415
```

这里将上述内容划分成几个模块进行分析

二、字段含义分析

```bash
name   user  nice   system      idle      iowait  irrq  softirq  steal guest guest_nice
cpu    60382   1     80032     198934063   2349     0     109      0     0       0
cpu0   2405    0     2084      4140924     682      0     6        0     0       0
...  # 此处较多冗余信息，简化之
cpu47  200     0     134       4147222     10       0     0        0     0       0
```

|            |                 |          |                                                                |
| ---------- | --------------- | -------- | -------------------------------------------------------------- |
| cpu 指标   | 含义            | 时间单位 | 备注                                                           |
| user       | 用户态时间      | jiffies  | 一般/高优先级，仅统计 nice<=0                                  |
| nice       | nice 用户态时间 | jiffies  | 低优先级，仅统计 nice>0                                        |
| system     | 内核态时间      | jiffies  |                                                                |
| idle       | 空闲时间        | jiffies  | 不包含 IO 等待时间                                             |
| iowait     | I/O 等待时间    | jiffies  | 硬盘 IO 等待时间                                               |
| irq        | 硬中断时间      | jiffies  |                                                                |
| softirq    | 软中断时间      | jiffies  |                                                                |
| steal      | 被盗时间        | jiffies  | 虚拟化环境中运行其他操作系统上花费的时间（since Linux 2.6.11） |
| guest      | 来宾时间        | jiffies  | 操作系统运行虚拟 CPU 花费的时间(since Linux 2.6.24)            |
| guest_nice | nice 来宾时间   | jiffies  | 运行一个带 nice 值的 guest 花费的时间(since Linux 2.6.33)      |

说明：

1. 1 jiffies = 0.01s = 10ms。也就是 100 分之一秒。该数值除以 100 就得到以秒为单位的数值
2. 如果把单独一个 CPU 的所有字段的时间加起来，其实就是系统的运行时间。而第一行的 CPU 是指的所有 CPU 时间之和，这个时间会超过系统运行时间。
3. 常用计算等式：CPU 时间 = user + system + nice + idle + iowait + irq + softirq
4. man 手册中 iowait 有单独说明，iowait 时间是不可靠值，具体原因如下：
   1. CPU 不会等待 I/O 执行完成，而 iowait 是等待 I/O 完成的时间。当 CPU 进入 idle 状态，很可能会调度另一个 task 执行，所以 iowait 计算时间偏小；
   2. 多核 CPU 中，iowait 的计算并非某一个核，因此计算每一个 cpu 的 iowait 非常困难
   3. 这个值在某些情况下会减少

# /proc/uptime

**系统启动的时间长度**

该文件包含两个数字:

1. 系统的正常运行时间（秒）
2. 空闲进程所花费的时间（秒）

# /proc/version - 内核**版本信息**

# 应用实例

1. cat /proc/net/bonding/\* # 在该目录查看网卡 bonding 信息
2. cat /proc/cpuinfo # 查看 CPU 信息
