---
title: Proc File System
---

# 概述

> 参考：
>
> - [Manual(手册)，proc(5)](https://man7.org/linux/man-pages/man5/proc.5.html)

**process information pseudo-filesystem(进程信息伪文件系统，简称 proc)**， 提供了内核数据结构的接口。`一般挂载到 /proc 目录`。一般情况是由操作系统自动挂载的，也可以通过`mount -t proc proc /proc`命令手动挂载。proc 文件系统中的大多数文件都是只读的，但是有些文件是可写的，用于改变内核参数。

proc 文件系统不用于存储。其主要目的是为硬件，内存，运行的进程和其他系统组件提供基于文件的接口。通过查看相应的 /proc 文件，可以检索许多系统组件上的实时信息。/proc 中的某些文件也可以（由用户和应用程序）操纵以配置内核。

# /proc/PID/ - 每个进程自己的独立信息

**每个进程在 /proc 下有一个名为自己进程号的目录，该目录记载了该进程相关的 proc 信息。**

## ./cgroup - 进程所属的控制组信息

详见 [Cgroup](/docs/10.云原生/Containerization/2.CGroup/2.CGroup.md)

## ./cmdline - 该进程的完整命令

- 除非进程是僵尸和内核态进程。在后一种情况下，这个文件中什么也没有：也就是说，对这个文件的读取将返回 0 个字符。命令行参数在这个文件中以一组字符串的形式出现，用空字节('\0')隔开，在最后一个字符串之后还有一个空字节。
- 把这个文件看作是进程希望你看到的命令行。

## ./exe - 具有该 PID 的实际运行的程序的绝对路径。是一个符号链接

## ./fd/ - 其中包含 PID 进程打开的每个文件的一个条目

- 该条目由其文件描述符命名，并且是指向实际文件的符号链接。 因此，0 是标准输入，1 是标准输出，2 是标准错误，依此类推。详解见：[File Descriptor(文件描述符)](/docs/1.操作系统/Kernel/Filesystem/文件管理/File%20Descriptor(文件描述符).md)

## ./fdinfo/ - 其中包含 PID 进程打开的每个文件的一个条目，该条目由其文件描述符命名

该目录中的文件仅由进程所有者读取。 可以读取每个文件的内容以获得有关相应文件描述符的信息。 内容取决于相应文件描述符所引用的文件类型。详解见：[File Descriptor(文件描述符)](/docs/1.操作系统/Kernel/Filesystem/文件管理/File%20Descriptor(文件描述符).md)

## ./maps - 进程的内存映射信息

```bash
> cat /proc/1751/maps
00400000-00401000 r-xp 00000000 fd:01 100897359                          /opt/java/jdk1.8.0_231/bin/java
00600000-00601000 r--p 00000000 fd:01 100897359                          /opt/java/jdk1.8.0_231/bin/java
00601000-00602000 rw-p 00001000 fd:01 100897359                          /opt/java/jdk1.8.0_231/bin/java
01542000-01563000 rw-p 00000000 00:00 0                                  [heap]
6c7c00000-6e0100000 rw-p 00000000 00:00 0
6e0100000-76d400000 ---p 00000000 00:00 0
76d400000-797580000 rw-p 00000000 00:00 0
797580000-7c0000000 ---p 00000000 00:00 0
7c0000000-7c18a0000 rw-p 00000000 00:00 0
7c18a0000-800000000 ---p 00000000 00:00 0
```

- address 字段表示进程中内存映射占据的地址空间，格式为十六进制的 BeginAddress-EndAddress。
- perms 字段表示权限，共四个字符，依次为 rwxs 或 rwxp，其中 r 为 read，w 为 write，x 为 execute，s 为- -shared，p 为 private，对应位置没有权限时用一个短横线代替。
- offset 字段表示内存映射地址在文件中的字节偏移量。
- dev 字段表示 device，格式为 major:minor。
- inode 字段表示对应 device 的 inode，0 表示内存映射区域没有关联的 inode，如未初始化的 BSS 数据段就是这种情况。
- pathname 字段用于内存映射的文件，对于 ELF 格式的文件来说，可以通过命令 readelf -l 查看 ELF 程序头部的 Offset 字段，与 maps 文件的 offset 字段作对比。pathname 可能为空，表示匿名映射，这种情况下难以调试进程，如 gdb、strace 等命令。除了正常的文件路径之外，pathname 还可能是下面的值：
  - \[stack]     初始进程（主线程）的 stack
  - \[stack:\<tid>]     线程 ID 为 tid 的 stack.  对应于/proc/\[pid]/task/\[tid]/路径
  - \[vdso]    Virtual Dynamically linked Shared Object
  - \[heap]     进程的 heap

## ./mountinfo - 进程的 mount namespace 的信息

该文件比 /proc/PID/mounts 的信息更全免并修复了一些问题，适用于云原生时代的大量 mount namesapces。

该文件中每一行都是一条挂载信息，每条挂载信息由如下几个部分组成：

| 挂载 ID | 父 ID | major:minor | root  | mount point | mount options | optional fields | separator | filesystem type | mount source | super options      |
| ------- | ----- | ----------- | ----- | ----------- | ------------- | --------------- | --------- | --------------- | ------------ | ------------------ |
| 36      | 35    | 98:00:00    | /mnt1 | /mnt2       | rw,noatime    | master:1        | -         | ext3            | /dev/root    | rw,errors=continue |

- mount ID # 挂载的唯一 ID
- parent ID # the ID of the parent mount (or of self for the root of this mount namespace's mount tree).
  - If a new mount is stacked on top of a previous existing mount (so that it hides the existing mount) at pathname P, then the parent of the new mount is the previous mount at that location. Thus, when looking at all the mounts stacked at a particular location, the top-most mount is the one that is not the parent of any other mount at the same location. (Note, however, that this top-most mount will be accessible only if the longest path subprefix of P that is a mount point is not itself hidden by a stacked mount.)
  - If the parent mount lies outside the process's root directory (see chroot(2)), the ID shown here won't have a corresponding record in mountinfo whose mount ID (field 1) matches this parent mount ID (because mounts that lie outside the process's root directory are not shown in mountinfo). As a special case of this point, the process's root mount masy have a parent mount (for the initramfs filesystem) that lies outside the process's root directory, and an entry for that mount will not appear in mountinfo.
- major:minor # the value of st_dev for files on this filesystem (see [stat(2)](https://man7.org/linux/man-pages/man2/stat.2.html)).
- root # the pathname of the directory in the filesystem which forms the root of this mount.
- mount point # 挂载点 the pathname of the mount point relative to the process's root directory.
- mount options # 挂载选项 per-mount options (see mount(2)).
- optional fields # zero or more fields of the form "tag\[:value]"; see below.
- separator # the end of the optional fields is marked by a single hyphen.
- filesystem type # 挂载的文件系统类型 the filesystem type in the form "type\[.subtype]".
- mount source # filesystem-specific information or "none".
- super options # 超级快选项。per-superblock options (see mount(2)).

## ./root/ - 每个进程的文件系统的 `/` 目录

/proc/PID/root/ 目录是一个指向进程根目录的软链接，效果如下：

```bash
~]# ls -l /proc/1192/root
lrwxrwxrwx 1 root root 0 Nov 11 10:42 root -> /

```

该目录通过 chroot(2) 系统调用设置。

该目录常被用来查看容器内的文件系统。与容器的 Merged 不同，该目录会包含所有挂载，这些挂载信息，来源于 /proc/PID/mountinfo 文件。

## ./smaps - 每个进程的内存映射的使用信息

仅当内核选项 CONFIG_PROC_PAGE_MONITOR 配置了之后，才存在该文件。`pmap` 命令会读取该文件，并以人类易读的形式显示信息。

进程的每一个映射，都有其对应的信息，文件格式如下：

```bash
00400000-00aa4000 r-xp 00000000 fc:01 1710957                            /bin/node_exporter
Size:               6800 kB
KernelPageSize:        4 kB
MMUPageSize:           4 kB
Rss:                5852 kB
......
00aa4000-01077000 r--p 006a4000 fc:01 1710957                            /bin/node_exporter
Size:               5964 kB
KernelPageSize:        4 kB
MMUPageSize:           4 kB
Rss:                5508 kB
......
7fe089fd7000-7fe08c708000 rw-p 00000000 00:00 0
......
7fe08c708000-7fe09c888000 ---p 00000000 00:00 0
......
7fe09c888000-7fe09c889000 rw-p 00000000 00:00 0
......
......
ffffffffff600000-ffffffffff601000 --xp 00000000 00:00 0                  [vsyscall]
......
```

该进程的映射信息与 /proc/PID/maps 文件中的内容相同。该文件的每一个映射信息下面的都包含内存使用量

- **Size** # 该映射的内存大小
  - 通过 `cat smaps| grep ^Size | awk '{print $2}' | awk '{sum += $1};END {print sum}' && cat status | grep VmSize` 命令，可以看到 smaps 文件中的 Size 与 status 文件中的 VmSize 是相同的
- **Rss** #
- **Pss** #
- ......

## ./smaps_rollup - 汇总 smaps 文件中每个映射的内存信息为一条结果

smaps_rollup 文件中的内容，是将 smaps 文件中每个映射的内存信息汇总之后的结果

```bash
[root@hw-cloud-xngy-jump-server-linux-2 /proc/1185]# cat smaps_rollup
00400000-7ffc50bfa000 ---p 00000000 00:00 0                              [rollup]
Rss:               21952 kB
Pss:               21948 kB
Pss_Anon:          10372 kB
Pss_File:          11576 kB
Pss_Shmem:             0 kB
Shared_Clean:          4 kB
Shared_Dirty:          0 kB
Private_Clean:     11576 kB
Private_Dirty:     10372 kB
Referenced:        21952 kB
Anonymous:         10372 kB
LazyFree:              0 kB
AnonHugePages:         0 kB
ShmemPmdMapped:        0 kB
FilePmdMapped:        0 kB
Shared_Hugetlb:        0 kB
Private_Hugetlb:       0 kB
Swap:                  0 kB
SwapPss:               0 kB
Locked:                0 kB
```

通过 awk 计算 Rss 的大小，可以看到 smaps 文件中的聚合值与 smaps_rollup 的值一样

```bash
[root@hw-cloud-xngy-jump-server-linux-2 /proc/1185]# cat smaps| grep ^Rss | awk '{print $2}' | awk '{sum += $1};END {print sum}'
21304
[root@hw-cloud-xngy-jump-server-linux-2 /proc/1185]# cat smaps_rollup
00400000-7ffc50bfa000 ---p 00000000 00:00 0                              [rollup]
Rss:               21304 kB
```

## ./statm - 进程的内存使用情况信息

该文件中只有一行，每列信息以空格分割，共 7 列

- size # 进程使用的内存。等同于 status 文件中的 VmSize，man 手册里这个描述与实际不符
- resident # 进程的 RSS。等同于 status 文件中的 VmRSS，man 手册里这个描述与实际不符
- shared # 与其他进程共享的内存。等同于 status 文件中的 RssFile + RssShmem，man 手册里这个描述与实际不符

由于内核内部可伸缩性优化，文件中的一些值并不准确，如果需要准确的值，可以查看 smaps 和 smaps_rollup 文件。

## ./status - 该进程的状态信息

包括但不限于 PID、该进程 CPU 与 Memory 的使用情况、等等。在这个文件中，包含了 ./stat 和 ./statm 文件中的许多信息。

- Name # 进程名称
- State # 进程状态
- VmSize # 进程申请的总内存。与 statm 文件中第一个字段的值相同
- VmRSS # 也就是进程当前时刻占用的物理内存。与 statm 文件中第二个字段的值相同

## ./task/TID/ - 进程的线程信息目录

该目录中是进程的每个线程的信息，目录名称(TID)为线程 ID

# /proc/cmdline - 引导系统时，传递给内核的参数

通常通过引导管理器（如 lilo（8）或 grub（8））完成。

# /proc/cpuinfo - CPU 信息

如 cpu 的类型、制造商、 型号和性能等。

# /proc/devices - 当前运行的核心配置的设备驱动的列表

# /proc/dma - 显示当前使用的 dma 通道

# /proc/filesystems - 内核可用的文件系统信息

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

详见：《[Memory 管理工具](/docs/1.操作系统/Linux%20管理/Linux%20系统管理工具/Memory%20管理工具.md)》 章节

# /proc/modules

存放当前加载了哪些核心模块信息。

# /proc/net/ - 网络层的信息

**软连接文件，被连接到 self/net 目录。主要是链接到的进程的 network namespace 的信息**

**./nf_conntrack** # 链接跟踪表，该文件用于记录已跟踪的连接

## 进程间通信所用 Socket 信息

**./tcp** # 所有的 TCP 连接信息。

**./tcp6** # 所有的基于 IPv6 的 TCP 连接信息。

参考：[GitHub Linux 项目文档](https://github.com/torvalds/linux/blob/master/Documentation/networking/proc_net_tcp.rst)

保存 TCP 套接字表的转储。除了调试之外，大部分信息都没有什么用。

- sl # 值是套接字的内核哈希槽位
- local_address # 是本地地址和端口号对
- rem_address # 是远程地址和端口号对(如果连接)
- St # 是套接字的内部状态。根据内核内存使用情况，
- tx_queue 和 rx_queue # 是传出和传入的数据队列。
- tr、tm->when 和 rexmits # 字段保存内核套接字状态的内部信息，仅在调试时有用。
- uid # 字段保存套接字创建者的有效 uid。
- inode # 该 socket 的 inode 号，后面一串 16 进制的字符是该 socket 在内存中的地址。

```bash
root@lichenhao:~# cat /proc/net/tcp
  sl  local_address rem_address   st tx_queue rx_queue tr tm->when retrnsmt   uid  timeout inode
   0: 0100007F:177A 00000000:0000 0A 00000000:00000000 00:00000000 00000000     0        0 12975942 1 ffff923dd621a300 100 0 0 10 0
   1: 3500007F:0035 00000000:0000 0A 00000000:00000000 00:00000000 00000000   101        0 28017 1 ffff923ef9dd08c0 100 0 0 10 0
   2: 00000000:0016 00000000:0000 0A 00000000:00000000 00:00000000 00000000     0        0 33221 1 ffff923eecf6c600 100 0 0 10 0
   3: F82A13AC:0016 CB2A13AC:FD4C 01 00000000:00000000 02:00025EB2 00000000     0        0 12973284 4 ffff923dd621e900 20 4 29 10 20
   4: F82A13AC:0016 CB2A13AC:FD48 01 00000000:00000000 02:000A6D4A 00000000     0        0 12944563 2 ffff923dd621bd40 20 4 31 10 23
```

注意：

这里用 16 进制表示的 IP 有点奇葩比如 `F82A13AC` 转换成 IP 地址是 `248.42.19.172`，真实 IP 地址是 `172.19.42.248`，也就是说反过来了。。。。`F82A13AC` 应该是 `AC132AF8`

**./udp** # 所有 UDP 连接信息

**./udp6** # 所有基于 IPv6 的 UDP 连接信息

**./unix** # 所有 Unix Domain Socket 连接信息

# /proc/softirqs - **软中断统计信息**

# /proc/self

**当某一进程访问此软连接时，该软连接将指向该进程自己的 /proc/PID/ 目录。**

# /proc/sys/

**sys 目录是可写的，可以通过它来访问或修改内核的参数。详见 [Kernel 参数](/docs/1.操作系统/Kernel/Linux%20Kernel/Kernel%20参数/Kernel%20参数.md) 文章**

# /proc/stat

**系统的不同状态，例如，系统启动后页面发生错误的次数。**

该文件包含系统启动以来的很多系统和内核的统计信息，平时大家比较关心的比如包括 CPU 运行情况、中断情况、启动时间、上线文切换次数、运行中的进程等信息都在其中。

一、文件全貌

```bash
# Linux下查看/proc/stat的具体信息如下
[root@WSC-31-2 ~]# cat /proc/stat
cpu  60382 1 80032 198934063 2349 0 109 0 0 0
cpu0 2405 0 2084 4140924 682 0 6 0 0 0
...  # 此处较多冗余信息，简化之
cpu47 200 0 134 4147222 10 0 0 0 0 0
intr 33622492 64 ... 0 0 0 0 # 此处较多冗余信息，简化之
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

# /proc/version # 内核**版本信息**

# 应用实例

1. cat /proc/net/bonding/\* # 在该目录查看网卡 bonding 信息
2. cat /proc/cpuinfo # 查看 CPU 信息

# 分类

# 文件系统
