---
title: sysstat 工具集
---

# 概述

> 参考：
>
> - [GitHub 项目，sysstat/sysstat](https://github.com/sysstat/sysstat)
> - [官网](http://sebastien.godard.pagesperso-orange.fr/)
> - [Manual(手册)，pidstat(1)](https://man7.org/linux/man-pages/man1/pidstat.1.html)
> - [Manual(手册)，sar(1)](https://man7.org/linux/man-pages/man1/sar.1.html)

sysstat 包包含很多类 UNIX 的应用程序，用以监控系统性能和使用活动。

# cifsiostat

# iostat - 报告设备和分区的 I/O 统计数据

iostat 命令用于通过观察设备活动的时间及其平均传输速率来监视系统 I/O 设备的负载。 iostat 命令生成可用于更改系统配置的报告，以更好地平衡物理磁盘之间的输入/输出负载。

iostat 命令生成两种类型的报告，即 CPU 利用率报告和设备利用率报告。

CPU 利用率报告

Device 利用率报告

```bash
~]# iostat -xd
Linux 4.18.0-193.19.1.el8_2.x86_64 (ansible.tj-test)  11/24/2020  _x86_64_ (4 CPU)
Device            r/s     w/s     rkB/s     wkB/s   rrqm/s   wrqm/s  %rrqm  %wrqm r_await w_await aqu-sz rareq-sz wareq-sz  svctm  %util
sda              0.04    0.17      1.56     11.25     0.00     0.26   0.19  60.08   10.57   54.93   0.01    40.14    65.48   0.97   0.02
scd0             0.00    0.00      0.00      0.00     0.00     0.00   0.00   0.00    6.41    0.00   0.00    38.52     0.00   1.00   0.00
dm-0             0.04    0.43      1.51     11.25     0.00     0.00   0.00   0.00    9.71  108.32   0.05    40.01    26.33   0.44   0.02
dm-1             0.00    0.00      0.00      0.00     0.00     0.00   0.00   0.00    1.69    0.00   0.00    21.57     0.00   0.39   0.00
```

| 指标       | 含义               | 提示                                                                                                    |
| -------- | ---------------- | ----------------------------------------------------------------------------------------------------- |
| r/s      | 每秒发送给磁盘的读请求数     | 合并后的请求数                                                                                               |
| w/s      | 每秒发送给磁盘的写请求数     | 合并后的请求数                                                                                               |
| rkB/s    | 每秒从磁盘读取的数据量      | 单位为 kB                                                                                                |
| wkB/s    | 每秒向磁盘写入的数据量      | 单位为 kB                                                                                                |
| rrqm/s   | 每秒合并的读请求数        | %rrqm表示合并读请求的百分比                                                                                      |
| wrqm/s   | 每秒合并的写请求数        | %wrqm表示合并写请求的百分比                                                                                      |
| r_await  | 读请求处理完成等待时间      | 包括队列中的等待时间和设备实际处理的时间，单位为毫秒                                                                            |
| w_await  | 写请求处理完成等待时间      | 包括队列中的等待时间和设备实际处理的时间，单位为毫秒                                                                            |
| aqu-sz   | 平均请求队列长度         | 旧版中为avgqu-sz                                                                                          |
| rareq-sz | 平均读请求大小          | 单位为 kB                                                                                                |
| wareq-sz | 平均写请求大小          | 单位为 kB                                                                                                |
| svctm    | 处理 I/O 请求所需的平均时间 | 单位为 毫秒。注意这是推断的数据，(不包括等待时间)并不保证完全准确                                                                    |
| %util    | 磁盘处理 I/O 的时间百分比  | 即使用率（详见 [Block](/docs/1.操作系统/Kernel/Hardware/Block.md) 的 “磁盘 I/O 时间” 章节），由于可能存在并行 I/O，100% 并不一定表明磁盘饱和 |

这些指标中你要注意：

- %util ，就是磁盘 I/O 使用率
- r/s + w/s ，就是 IOPS
- rkB/s + wkB/s ，就是吞吐量
- r_await + w_await ，就是响应时间

## Syntax(语法)

**iostat \[OPTIONS] \[INTERVAL \[COUNT]]**

- INTERVAL # 间隔时间，单位是秒，指定 INTERVAL 时，mpstat 根据该时间每隔 INTERVAL 秒输出一次信息，并在最后输出平均值。
  - COUNT # 每隔 INTERVAL 时间，输出信息的数量。若不指定 CONUNT，则 mpstat 会根据间隔时间持续输出统计信息。

OPTIONS

- **-c** # 只显示 CPU 利用率报告
- **-d** # 只显示磁盘利用率报告
- **-x** # 显示扩展信息。输出更多的统计信息

EXAMPLE

# mpstat - 显示处理器相关的统计信息

**mpstat \[OPTIONS] \[INTERVAL \[COUNT]]**

- INTERVAL # 间隔时间，单位是秒。指定 INTERVAL 时，mpstat 根据该时间每隔 INTERVAL 秒输出一次信息，并在最后输出平均值。
  - COUNT # 每隔 INTERVAL 时间，输出信息的数量。若不指定 CONUNT，则 mpstat 会根据间隔时间持续输出统计信息。

OPTIONS

- **-A** # 此选项等效于 -u -I ALL -P ALL
- **-P** **\<CPU\[,CPU2....]>** # 指定要监控的 CPU，ALL 为监控所有 CPU

EXAMPLE

- 显示所有 CPU 的统计信息。每隔 5 秒输出 1 次信息，总共输出 2 次。
  - **mpstat -P ALL 5 2**

# nfsiostat-sysstat

# pidstat - 显示 Linux 进程的统计信息

> 参考：
>
> - [Manual(手册)，pidstat(1)](https://man7.org/linux/man-pages/man1/pidstat.1.html)

pidstat 是一个以 Task(任务) 为主体，显示 Task 相关系统使用情况的工具。

> Task(任务) 是 进程、线程 之类的统称。

根据命令的不同选项，显示进程的不同信息。

## -d - 报告 I/O 统计信息

```bash
~]$ pidstat -d 1 --human
Linux 5.4.0-88-generic (hw-cloud-xngy-jump-server-linux-2)  10/03/2021  _x86_64_ (2 CPU)

10:27:48 PM   UID       PID   kB_rd/s   kB_wr/s kB_ccwr/s iodelay  Command

10:27:49 PM   UID       PID   kB_rd/s   kB_wr/s kB_ccwr/s iodelay  Command
10:27:50 PM  1000     16829      0.0B      4.0k      0.0B       0  bash

10:27:50 PM   UID       PID   kB_rd/s   kB_wr/s kB_ccwr/s iodelay  Command
^C

Average:      UID       PID   kB_rd/s   kB_wr/s kB_ccwr/s iodelay  Command
Average:     1000     16829      0.0B      1.3k      0.0B       0  bash
```

## -u - 默认选项。报告进程的 CPU 利用率

```bash
~]# pidstat
Linux 4.18.0-193.19.1.el8_2.x86_64 (ansible.tj-test)  10/27/2020  _x86_64_ (4 CPU)
09:51:15 PM   UID       PID    %usr %system  %guest   %wait    %CPU   CPU  Command
09:51:15 PM     0         1    0.00    0.01    0.00    0.00    0.02     1  systemd
09:51:15 PM     0         2    0.00    0.00    0.00    0.00    0.00     0  kthreadd
```

## -w - 报告进程的上下文切换情况

```bash
~]# pidstat -w
Linux 4.18.0-193.19.1.el8_2.x86_64 (ansible.tj-test)  10/27/2020  _x86_64_ (4 CPU)
10:00:14 PM   UID       PID   cswch/s nvcswch/s  Command
10:00:14 PM     0         1      0.11      0.03  systemd
10:00:14 PM     0         2      0.02      0.00  kthreadd
```

- **cswch/s** # 每秒自愿上下文切换(voluntary context switches)的次数
- **nvcswch** # 每秒非自愿上下文切换(non voluntary context switches)的次数

## -r - 报告进程的内存使用情况统计信息

```bash
~]# pidstat -r
Linux 4.18.0-193.28.1.el8_2.x86_64 (desistdaydream.bj-net)  11/18/2020  _x86_64_ (2 CPU)
09:58:16 PM   UID       PID  minflt/s  majflt/s     VSZ     RSS   %MEM  Command
09:58:16 PM     0         1      5.57      0.05  176812   10844   0.28  systemd
09:58:16 PM     0       664      0.43      0.00   91980    8980   0.23  systemd-journal
```

- **minflt/s** # 每秒任务执行的次要故障总数，即不需要从磁盘加载内存页面的次要故障总数。
- **majflt/s** # 每秒任务执行的主要故障总数，即需要从磁盘加载内存页面的主要故障总数。
- **VSZ** # Virtual Size(虚拟大小)。整个任务的虚拟内存使用量，单位: KiB。
- **RSS** # Resident Set Size(常驻集大小)。任务使用的未交换的物理内存，单位: KiB。

## -d - 报告进程的磁盘 I/O 统计信息

```bash
~]# pidstat -d
Linux 4.18.0-193.28.1.el8_2.x86_64 (desistdaydream.bj-net)  11/18/2020  _x86_64_ (2 CPU)
09:58:37 PM   UID       PID   kB_rd/s   kB_wr/s kB_ccwr/s iodelay  Command
09:58:37 PM     0         1     56.46      0.18      0.01      33  systemd
09:58:37 PM     0         7      0.18      0.00      0.00       0  kworker/u4:0-events_unbound
```

## Syntax(语法)

**pidstat \[OPTIONS] \[INTERVAL \[COUNT]]**

- **INTERVAL** # 间隔时间，单位是秒，指定 INTERVAL 时，mpstat 根据该时间每隔 INTERVAL 秒输出一次信息，并在最后输出平均值。
  - **COUNT** # 每隔 INTERVAL 时间，输出信息的数量。若不指定 CONUNT，则 mpstat 会根据间隔时间持续输出统计信息。

OPTIONS

- **-C STRING** # 仅显示 Task 名称包含 STRING 的任务。该字符串可以是正则表达式。就是根据进程名过滤。
- **-p \<PID\[,PID2....]>** # 只显示指定的一个或多个进程的信息。默认关键字为 ALL，显示所有信息。可以用 SELF 关键字来只显示自身的信息
- **-t** # 显示进程所关联的线程的统计信息。

EXAMPLE

- **pidstat -u 5 1**

# sadf - 以多种格式显示 sar 工具收集到的数据

# sar - 系统活动报告

**system activity report(系统活动报告，简称 SAR)**。sar 是以系统为主体，报告系统相关信息的工具，包括 CPU 使用率、RAM 使用情况、磁盘 IO、网络活动状态等等。

根据命令的不同选项，显示不同信息。

## -b - 报告 I/O 和传输速率的统计信息

```bash
~]$ sar --human -b 1
Linux 5.4.0-88-generic (hw-cloud-xngy-jump-server-linux-2)  10/03/2021  _x86_64_ (2 CPU)

10:07:46 PM       tps      rtps      wtps      dtps   bread/s   bwrtn/s   bdscd/s
10:07:47 PM      0.00      0.00      0.00      0.00      0.00      0.00      0.00
10:07:48 PM      0.00      0.00      0.00      0.00      0.00      0.00      0.00
10:07:49 PM      0.00      0.00      0.00      0.00      0.00      0.00      0.00
10:07:50 PM      2.00      0.00      2.00      0.00      0.00     16.00      0.00


^CAverage:         0.50      0.00      0.50      0.00      0.00      3.99      0.00
```

- **tpc** # Transfers Per Second(每秒传输总数)，即每秒发送到物理设备的传输总数。也就是物理设备的 I/O 请求数。比如上面的例子，就是在 10 点 07 分 49 秒时，发起的两次 I/O 请求，这里面并不关心 I/O 的传输总量，或者具体是读还是写，仅记录次数。
- **rtps** # Read TPS，每秒向物理设备发送的读请求总数
- **wtps** # Write TPS，每秒向物理设备发送的写请求总数
- **dtps** # Discard TPS，每秒向物理设备发送的丢弃请求总数
- **bread/s** # Black Read，每秒从设备中读取的数据总量(以块为单位)
- **bwrtn/s** # Black Written，每秒写入设备的数据总量(以块为单位)
- **bdscd/s** # Black Discarded，每秒为设备丢弃的丢弃的数据总量(以块为单位)

## -n KEYWORD - 报告网络统计信息

KEYWORD 可用的值有 DEV、EDEV、NFS、NFSD、SOCK、IP、EIP、ICMP、EICMP、TCP、ETCP、UDP、SOCK6、IP6、EIP6、ICMP6、EICMP6、UDP6。当然，也可以用使用 ALL 来报告所有的网络统计信息

### DEV - 报告指定网络设备的统计信息

```bash
~]# sar -n DEV 1
Linux 4.18.0-193.28.1.el8_2.x86_64 (desistdaydream.bj-net)  11/18/2020  _x86_64_ (2 CPU)
10:18:38 PM     IFACE   rxpck/s   txpck/s    rxkB/s    txkB/s   rxcmp/s   txcmp/s  rxmcst/s   %ifutil
10:18:39 PM   docker0      0.00      0.00      0.00      0.00      0.00      0.00      0.00      0.00
10:18:39 PM      ens3      4.00      1.00      0.23      0.10      0.00      0.00      0.00      0.00
10:18:39 PM        lo      0.00      0.00      0.00      0.00      0.00      0.00      0.00      0.00
```

- **rxpck/s** 和 **txpck/s** 分别表示每秒接收、发送的网络帧数，也就是 PPS
- **rxkB/s** 和 **txkB/s** 分别表示每秒接收、发送的千字节数，也就是 BPS
- **rxcmp/s** 和 **txcmp/s**  分别表示每秒接收、发送的压缩数据包
- **rxmcst/s** # 表示每秒接收的多播数据包

### TCP - 报告 TCPv4 网络流量的统计信息

```bash
~]$ sar -n TCP 1
Linux 5.4.0-88-generic (hw-cloud-xngy-jump-server-linux-2)  10/03/2021  _x86_64_ (2 CPU)

09:46:03 PM  active/s passive/s    iseg/s    oseg/s
09:46:04 PM      0.00      0.00      2.00      2.00
09:46:05 PM      0.00      0.00      2.00      3.00
09:46:06 PM      0.00      0.00      2.00      2.00
09:46:07 PM      0.00      0.00      2.00      2.00
^C
Average:         0.00      0.00      2.00      2.25
```

- **active/s** # TCP 连接每秒从 CLOSED 状态直接转换到 SYN-SENT 状态的次数。
  - 每秒本地发起 TCP 连接数，例如通过 connect() 系统调用。(待确认描述)
- **passive/s** # TCP 连接每秒从 LISTEN 状态直接转换到 SYN-RCVD 状态的次数。
  - 每秒远程发起的 TCP 连接数，例如通过 accept() 系统调用。(待确认描述)
- **iseg/s** # 每秒接受的 TCP segments 总数，包括错误的。
  - 也就是每秒收到了多少个 TCP 包
- **oseg/s** # 每秒发送的 TCP segments 总数。不包括重传的。
  - 也就是每秒发送了多少个 TCP 包

## Syntax(语法)

**sar \[OPTIONS] \[INTERVAL \[COUNT]]**

- **INTERVAL** # 间隔时间，单位是秒。指定 INTERVAL 时，mpstat 根据该时间每隔 INTERVAL 秒输出一次信息，并在最后输出平均值。
  - **COUNT** # 每隔 INTERVAL 时间，输出信息的数量。若不指定 CONUNT，则 mpstat 会根据间隔时间持续输出统计信息。

OPTIONS

- **-p \<CPU\[,CPU2....]>** # 只显示指定的一个或多个 CPU 的信息，以需要表示，多个 CPU 以逗号分割。

EXAMPLE

- 输出 1，3，5，7 这 4 个 CPU 中，idle 小于 10 的 CPU
  - **sar -P 1,3,5,7 1 | tail -n+3 | awk '$NF<10 {print $0}'**

# tapestat
