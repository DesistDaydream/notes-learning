---
title: fio 磁盘性能测试工具
linkTitle: fio 磁盘性能测试工具
date: 2024-05-13T09:19
weight: 1
---

# 概述

> 参考:
>
> 项目地址：<https://github.com/axboe/fio>
> 官方文档：<https://fio.readthedocs.io/en/latest/>

注意：！！当使用 fio 的 filename 参数指定某个要测试的裸设备（硬盘或分区），切勿在系统分区做测试，会破坏系统分区，从而导致系统崩溃。若一定要测试系统分区较为安全的方法是：在根目录下创建一个空目录，在测试命令中使用 directory 参数指定该目录，而不使用 filename 参数。现在假设 /dev/vda3 设备挂载在 / 目录下，那么不要执行 fil --filename=/dev/vda 这种操作！！

## 1 性能的基本概念

### 1.1 什么是一个 IO

IO 即 Input 和 Output，可以分为读 IO 和写 IO。读 IO，就是发指令，从磁盘读取某段扇区的内容。指令一般是通知磁盘开始扇区的位置，然后给出需要从这个初始扇区往后读取的连续扇区个数，同时给出的动作是读，还是写。磁盘收到这条指令，就会按照指令的要求，读或者写数据。控制器发出的这种指令+数据，就是一次 IO，读或者写。

### 1.2 顺序 IO 和随机 IO

顺序和随机，可以简单地理解为本次 IO 给出的初始扇区地址，和上一次 IO 的结束扇区地址，是否是按顺序的，如果相差很大，就算一次随机 IO。

### 1.3 IO 大小

一次 IO 需要读或者写的数据块大小。

### 1.4 带宽

每秒读出或写入的数据量。常用单位包括 KB/s、MB/s、GB/s 等。

### 1.5 延时

客户端发出 IO 请求直到收到请求并响应是需要一段时间的，这段时间就是 IO 延时。IO 延 时一般都是毫秒级的。随着 IO 压力的增大，IO 延时也会随之增大。对于存储来说，由于写 是前台操作，而读是后台操作，因此通常写的 IO 延时要低于读。相同 IO 模型下，IO 延时越小，存储性能越好。一般，IO 延时如果超过 30ms 就说明存储已经比较吃力了。

### 1.6 IOPS

每秒收到的 IO 响应数。读写比例、队列深度、随机度和块大小描述了 IO 模型。在这个模型下，每个 IO 的延时最终体现了 IO 性能。IOPS 和 带宽从两个不同的方面反映了存储性能。

### 1.7 队列深度

IO 队列使得 IO 可以并行处理，为何要对磁盘 IO 进行并行处理呢？主要目的是提高磁盘处理 IO 的效率。这一点对于多物理磁盘组成的虚拟磁盘（或 LUN）显得尤为重要。如果 一次提交一个 IO，虽然响应时间较短，但系统的带宽很小。相比较而言，一次提交多个 IO 既缩短了磁头移动距离（通过电梯算法），同时也能够提升 IOPS。

## fio 的工作方式

想要运行 fio 可以通过两种方式

1. 给定一个或多个 job files
2. 直接使用命令行参数

如果 job file 中仅包含一个作业，则最好只在命令行中提供参数。命令行参数与作业参数相同，还有一些额外的参数可控制全局参数。这里面提到的 job file 中的擦书，其实也可以理解为指令、关键字、字段等等。就是一个配置文件中用来描述程序运行行为的东西。

1. 例如，对于 job files 中的参数 iodepth = 2，在命令行选项为--iodepth 2 或--iodepth = 2。
2. 还可以使用命令行来提供多个作业条目。对于 fio 看到的每个--name 选项，它将使用该名称开始一个新作业。 --name 条目之后的命令行条目将应用于该作业，直到没有更多的条目或看到新的--name 条目为止。这类似于作业文件选项，其中每个选项都适用于当前作业，直到看到新的 \[] 作业条目为止。

开始模拟 I/O 工作负载的第一步，是编写一个描述特定配置的 job file。作业文件可以包含任何数量的线程和/或文件-作业文件的典型内容是定义共享参数的全局部分，以及一个或多个描述所涉及作业的作业部分。运行时，fio 会分析该文件并按照说明进行所有设置。如果我们从上到下分解一个 job ，它包含以下基本参数：

1. I/O type # 定义发布给文件的 I / O 模式。我们可能只从该文件中顺序读取，或者我们可能在随机写入。甚至顺序或随机混合读写。我们应该执行缓冲 I / O 还是直接/原始 I / O？
2. Block size # 指定模拟 I/O 数据流时，每次 I/O 的块大小。可以是单个值，也可以描述块大小的范围。
3. I/O size # 指定本次 job 将要读取或写入多少数据
4. I/O engine # 定义如何向文件发出 I/O。我们可以使用内存映射文件，可以使用常规读/写，可以使用拼接，异步 I / O 甚至是 SG（SCSI 通用 sg）
5. I/O depth # 定义在 I/O engine 是异步的时，我们要保持多大的队列深度
6. Target file/device # How many files are we spreading the workload over.
7. Threads, processes and job synchronization # How many threads or processes should we spread this workload over.

# Job File 格式

Job file 参数详见：fio 参数详解

如前所述，fio 接受一个或多个描述该操作的作业文件。作业文件格式是经典的 ini 文件，其中\[]括号中的名称定义了作业名称。您可以随意使用任何所需的 ASCII 名称，但具有特殊含义的 global 除外。作业名称后面是零个或多个参数的序列，每行一个，用于定义作业的行为。如果一行中的第一个字符是“;”或“＃”，则整行都将作为注释被丢弃。

全局部分为该文件中描述的作业设置默认值。作业可以覆盖全局节参数，并且如果需要的话，作业文件甚至可以具有多个全局节。作业仅受位于其上方的全局部分影响。

因此，让我们看一个非常简单的作业文件，该文件定义了两个过程，每个过程都随机读取 128MiB 文件：

```ini
; -- start job file --
[global]
rw=randread
size=128m
[job1]
[job2]
; -- end job file --
```

如您所见，作业文件部分本身为空，因为所有描述的参数都是共享的。由于未提供文件名选项，因此，fio 会为每个作业组成一个合适的文件名。在命令行上，此作业如下所示：

```bash
$ fio --name=global --rw=randread --size=128m --name=job1 --name=job2
```

让我们看一个示例，其中有许多进程随机写入文件：

```ini
; -- start job file --
[random-writers]
ioengine=libaio
iodepth=4
rw=randwrite
bs=32k
direct=0
size=64m
numjobs=4
; -- end job file --
```

这里我们没有全局部分，因为我们只定义了一项工作。我们想在这里使用异步 I / O，每个文件的深度为 4。我们还将缓冲区大小增加到 32KiB，并将 numjobs 定义为 4，以分叉 4 个相同的作业。结果是 4 个进程，每个进程随机写入其自己的 64MiB 文件。您可以在命令行上指定参数，而不使用上面的作业文件。对于这种情况，您可以指定：

```bash
$ fio --name=random-writers --ioengine=libaio --iodepth=4 --rw=randwrite --bs=32k --direct=0 --size=64m --numjobs=4
```

# fio 命令行工具

## Syntax(语法)

**fio \[OPTIONS] \[JOB OPTIONS] \[job file(s)]**

OPTIONS

JOB OPTIONS

- 由于 Jobfile 中的参数与命令行选项基本保持一一对应的关系，所以对于 fio 的命令行参数，参考 Job file 参数即可 fio 参数详解

EXAMPLE

- 注意，下面两条命令直接对整块磁盘进行写操作，会破坏文件，谨慎操作
   - fio -filename=/dev/vdb1 -direct=1 -iodepth 64 -thread -rw=randwrite -ioengine=libaio -bs=4K -numjobs=8 -runtime=60 -group_reporting -name=test1
   - fio -filename=/dev/vdb3 -direct=1 -iodepth 64 -thread -rw=write -ioengine=libaio -bs=512K -numjobs=8 -runtime=60 -group_reporting -name=test2
- **fio --rw=write --ioengine=sync --fdatasync=1 --directory=/var/lib/etcd --size=22m --bs=2300 --name="fioEtcdTest" --time_based --runtime=2m**
- **fio --rw=write --ioengine=libaio --iodepth=4 --direct=1 --filename=fiotest --size=2G --bs=4k --name="Max throughput" --time_based --runtime=60**

## 结果分析

官方文档：<https://fio.readthedocs.io/en/latest/fio_doc.html#interpreting-the-output>

使用 fio -ioengine=libaio -bs=4k -direct=1 -thread -rw=write -size=2G -filename=test -name="Max throughput" -iodepth=4 -runtime=60 命令得到如下结果

```bash
# 前面几行是命令运行时，实时显示的信息
# 名为 Max throughput 的 Job 基本信息
Max throughput: (g=0): rw=write, bs=(R) 4096B-4096B, (W) 4096B-4096B, (T) 4096B-4096B, ioengine=libaio, iodepth=4
fio-3.7
 # fio 版本号
Starting 1 thread
 # 本次此时启动了 1 个线程
# 1 表示运行的IO线程数；[W(1)] 表示使用的模式；[100.0%] 表示当前命令的执行进度；[r=0KiB/s,w=137MiB/s] 表示瞬时吞吐率；
# [r=0,w=35.0k IOPS] 表示 IOPS 瞬时值；[eta 00m:00s] 表示持续时间
# 该行内容在命令执行期间，可以看到值在实时变化
Jobs: 1 (f=1): [W(1)][100.0%][r=0KiB/s,w=137MiB/s][r=0,w=35.0k IOPS][eta 00m:00s]

# 从本行开始为命令执行完成后每个 Job 的统计信息。
# Job名称：(当前的GID,Job个数)：错误个数：本次Job的PID，本次Job的结束时间
Max throughput: (groupid=0, jobs=1): err= 0: pid=7767: Thu Nov 12 16:09:04 2020
# 本次Job的测试模式(这里是写入)，IOPS平均值，带宽平均值，(带宽最大值)(数据总量/运行总时间)
# BW 是 BandWidth 的缩写。2048MiB 就是指的 -size 参数指定的 大小，是本次测试读/写的数据总量
  write: IOPS=34.3k, BW=134MiB/s (140MB/s)(2048MiB/15289msec)

# latency(延迟)相关信息。(单位是：微妙)。注意后面信息括号中的单位会改变，fio 会根据本次测试结果得出的时间，合理给出一个单位。msec毫秒，usec微秒，nsec纳秒
# slat 是 Submission latency(提交延迟)，就是提交到实际执行 I/O 的时间
    slat (usec): min=3, max=6347, avg= 6.09, stdev=12.68
# clat 是 Completion latency(完成延迟)，就是从提交到完成的时间
    clat (usec): min=57, max=74949, avg=108.13, stdev=230.83
# lat 是 Total latency(总延迟)，就是 fio 从创建这个 I/O 单元到完成的总时间
     lat (usec): min=63, max=74954, avg=114.80, stdev=231.22
# 完成延迟的百分位数(单位是：微妙)，比如99.00th=[  149] 表示这组样本的 99 百分位的延迟的值为 149
    clat percentiles (usec):
     |  1.00th=[   86],  5.00th=[   92], 10.00th=[   94], 20.00th=[   97],
     | 30.00th=[   99], 40.00th=[  101], 50.00th=[  103], 60.00th=[  106],
     | 70.00th=[  109], 80.00th=[  113], 90.00th=[  119], 95.00th=[  125],
     | 99.00th=[  149], 99.50th=[  192], 99.90th=[  330], 99.95th=[  988],
     | 99.99th=[ 5276]
# 基于一组样本的带宽信息(单位是 KiB/s):最小值，最大值，该线程在其组中接收的总聚合带宽的大约百分比，平均值，标准偏差，本次采样的样本总数
   bw (  KiB/s): min=113928, max=149024, per=99.98%, avg=137137.07, stdev=7370.84, samples=30
# 基于一组样本的IOPS信息，与 bw 一样
   iops        : min=28482, max=37256, avg=34284.27, stdev=1842.71, samples=30
# I/O 完成延迟的分布，这里的信息适用于一组报告的所有 I/O
# 500=0.27% 表示 0.27% 的 I/O 在 500 微妙以内完成
  lat (usec)   : 100=34.55%, 250=65.12%, 500=0.27%, 750=0.01%, 1000=0.01%
  lat (msec)   : 2=0.01%, 4=0.02%, 10=0.02%, 20=0.01%, 100=0.01%
# cpu 使用率
  cpu          : usr=10.24%, sys=23.20%, ctx=189677, majf=0, minf=1
# IO 深度在整个工作周期中分布
  IO depths    : 1=0.1%, 2=0.1%, 4=100.0%, 8=0.0%, 16=0.0%, 32=0.0%, >=64=0.0%
# 一个提交调用中要提交多少IO
     submit    : 0=0.0%, 4=100.0%, 8=0.0%, 16=0.0%, 32=0.0%, 64=0.0%, >=64=0.0%
# 一个完成调用中要完成多少IO
     complete  : 0=0.0%, 4=100.0%, 8=0.0%, 16=0.0%, 32=0.0%, 64=0.0%, >=64=0.0%
# 发出的读取/写入/修剪请求的数量，以及其中有多少个请求被缩短或丢弃。
     issued rwts: total=0,524288,0,0 short=0,0,0,0 dropped=0,0,0,0
     latency   : target=0, window=0, percentile=100.00%, depth=4

# 在上面将所有 Job 的统计信息都列出来之后，这下面显示所有 Job 最终的统计信息。这些数据是所有 Job 的平均值。
Run status group 0 (all jobs):
# 本次测试的模式(写测试)：带宽平均值，带宽最小值-带宽最大值，io(i.e.读写数据的总量)，运行时间
# 括号外的值是2的幂的格式，括号中的值是10的幂的等效值
  WRITE: bw=134MiB/s (140MB/s), 134MiB/s-134MiB/s (140MB/s-140MB/s), io=2048MiB (2147MB), run=15289-15289msec


# 当前测试数据所在磁盘的统计信息
# ios 表示所有组执行的 I/O 数，ios就是 I/Os。
# merge 表示 I/O 调度程序执行的合并数。
# ticks 表示我们保持磁盘活跃的 ticks 数。
# in_queue 表示在磁盘队列中花费的总时间
# util 表示磁盘利用率。在命令执行期间 100％表示我们使磁盘一直处于繁忙状态，而50％的磁盘将有一半的时间处于空闲状态
# aggr 前缀的信息官方没有说明,应该是聚合的意思
Disk stats (read/write):
    dm-0: ios=0/519394, merge=0/0, ticks=0/47273, in_queue=47273, util=97.70%, aggrios=0/524290, aggrmerge=0/0, aggrticks=0/48534, aggrin_queue=1310, aggrutil=97.56%
  vda: ios=0/524290, merge=0/0, ticks=0/48534, in_queue=1310, util=97.56%
```

可以看到，在一个非常强悍的 Optane 盘上面，使用 sync engine，每次都 sync 写盘，性能还是很差的，吞吐不到 300 MB，其他的盘可能就更差了。我们主要关注几个指标：

slat/clat/lat：这几个是 latency 指标，slat 就是 Submission latency，也就是提交到实际执行 I/O 的时间，在 sync 测试里面这个是没有的，因为 slat 就是 clat。clat 就是 Completion latency，也就是从提交到完成的时间。lat 就是 Total latency，包括 fio 从创建这个 I/O 单元到完成的总的时间。

另外需要关注的指标就是 BW，和 IOPS，这两这个很直观了，就不解释了。最下面是 ios，也就是总的 I/O 操作次数，merge 就是被 I/O 调度合并的次数，ticks 就是让磁盘保持忙碌的次数，in_queue 就是总的在磁盘队列里面的耗时，而 util 则是磁盘的利用率。

## 其他测试命令

fio --filename=/dev/vdb -direct=1 -bs=1M -rw=randwrite -ioengine=libaio -size=50g -numjobs=32 -iodepth=32  -group_reporting -name=mytest -thread --time_based --runtime=120

<https://mp.weixin.qq.com/s/zpkheD6Izn0RsipSukHA5Q>
