---
title: fio 参数详解
---

# fio Job file 参数详解

官方文档：<https://fio.readthedocs.io/en/latest/fio_doc.html#job-file-parameters>

## 参数类型：

Job file 的参数有多种类型，每种类型的参数的值可用的类型各不相同，比如时间类型的参数的值就需要填写时间相关的值。

str # 字符串类型。值为字符数字字符序列

time # 时间类型。值为带时间后缀的整数。时间的默认单位是秒(seconds)。可以指定其他单位： 天(d)、小时(h)、分钟(m)、毫秒(ms 或 msec)、微秒(us 或 usec)。e.g.使用 10m 表示 10 分钟。

int # 整数类型。整数值，可以包含整数前缀和整数后缀

bool # 布尔类型。通常解析为整数，但是仅定义为 true 和 false（1 和 0）

irange # 范围类型。带后缀的整数范围。允许给出值范围，例如 1024-4096。冒号也可以用作分隔符，例如。 1k：4k。如果该选项允许两组范围，则可以使用 `,`或 `/` 定界符来指定它们：1k-4k / 8k-32k。另请参见 int.

float_list # 浮点列表类型。浮点数列表，以 `:` 字符分隔。

# 描述 JOB 的相关参数

name=STR # 这可以用来覆盖由 fio 为该作业打印的名称。否则，使用作业名称。在命令行上，此参数的特殊目的还用于发信号通知新作业的开始。

# 时间相关参数

官方文档：<https://fio.readthedocs.io/en/latest/fio_doc.html#time-related-parameters>

runtime=TIME # 指定 Job 运行的时间(默认单位：秒)。到时间后，不管指定的 size 大小有没有读写完。

time_based # 如果设置，则 fio 将在 `runtime 的值`这个时间内内运行，即使已完全读取或写入文件。它会在运行时允许的情况下简单地循环相同的工作负载。

> 该参数一般配合 runtime 一起使用，单独使用没有效果。

# 要测试的目标文件/设备相关参数

directory=STR # 测试目录。

filename=STR # 测试文件名称。通常选择需要测试的盘的 data 目录

注意：！！当使用 fio 的 filename 参数指定某个要测试的裸设备（硬盘或分区），切勿在系统分区做测试，会破坏系统分区，而导致系统崩溃。若一定要测试系统分区较为安全的方法是：在根目录下创建一个空目录，在测试命令中使用 directory 参数指定该目录，而不使用 filename 参数。现在假设 /dev/vda3 设备挂载在 / 目录下，那么不要执行 fil --filename=/dev/vda 这种操作

# I/O TYPE 相关参数

官方文档：<https://fio.readthedocs.io/en/latest/fio_doc.html#i-o-type>

direct=BOOL # 测试过程绕过系统的 buffer。使测试结果更真实。`默认值：false`i.e.使用缓存。

rw=STR # 测试时的 I/O 模式。默认为 read。可用的模式有：

1. read # 顺序读
2. write # 顺序写
3. randread # 随机读
4. randwrite # 随机写
5. randrw # 随机混合读和写
6. 等.....其余模式见官方文档

fdatasync=INT # 与 fsync 类似，但使用 fdatasync(2)只同步数据而不同步元数据块。在 Windows, FreeBSD, DragonFlyBSD 或 OSX 中没有 fdatasync(2)，所以这就回到使用 fsync(2)。默认值为 0，这意味着 fio 不会定期发出问题，并等待数据同步完成。

# Block 大小相关参数

官方文档：<https://fio.readthedocs.io/en/latest/fio_doc.html#block-size>

bs=INT # 单次 I/O 的块文件大小

bsrange=iRANGE # 与 bs 参数类似，只不过是指定一个块文件大小的范围

# Buffers and memory 相关参数

zero_buffers # 用 0 初始化系统 buffer

# I/O Size 相关参数

size=INT # 本次的测试文件大小为 INT。fio 程序将持续运行，直到传输了 INT 这些数据。

# I/O engine(引擎) 相关参数

ioengine=STR # 告诉 fio 使用什么样的方式去测试 I/O。有如下方式可用：

- sync # 也就是最通常的 read/write 操作。基本的读(2)或写(2)I/O。lseek(2)用于定位 I/O 位置。请参阅 fsync 和 fdatasync 来同步写 I/O。
- psync # 基本的 pread(2)或 pwrite(2) I/O。除 Windows 外，所有支持的操作系统都是默认值。
    - pvsync / pvsync2 - 对应的 preadv / pwritev，以及 preadv2 / p writev2
- vsync # 使用 readv / writev，主要是会将相邻的 I/O 进行合并
- libaio # Linux 原生的异步 I/O，这也是通常我们这边用的最多的测试盘吞吐和延迟的方法
- 1. 对于 libaio engine 来说，还需要考虑设置 iodepth

测试多了，就会很悲催的发现，libaio 很容易就把盘给打死，但 sync 这些还需要启动几个线程。。。

并且对于 fio --rw=write --ioengine=XXXX --filename=fiotest --direct=1 --size=2G --bs=4k --name="Max throughput" --iodepth=4 命令，sync 引擎测试结果只有 libaio 引擎测试结果的三分之一。

# I/O engine(引擎) 特定参数

官方文档：<https://fio.readthedocs.io/en/latest/fio_doc.html#i-o-engine-specific-parameters>

# I/O depth(深度) 相关参数

官方文档：<https://fio.readthedocs.io/en/latest/fio_doc.html#i-o-depth>

iodepth=INT # 针对文件保留的 I/O 单元数。 `默认值：1`。

> 请注意，将 iodepth 增加到 1 以上不会影响同步 ioengine(使用 verify_async 时的小角度除外)。 甚至异步引擎也可能会施加 OS 限制，从而导致无法实现所需的深度。 在 Linux 上使用 libaio 且未设置 direct = 1 时可能会发生这种情况，因为在该 OS 上缓冲的 I / O 并不异步。 密切注意 fio 输出中的 I / O 深度分布，以验证所达到的深度是否符合预期。

# I/O rate 相关参数

官方文档：<https://fio.readthedocs.io/en/latest/fio_doc.html#i-o-rate>

# I/O latency

官方文档：<https://fio.readthedocs.io/en/latest/fio_doc.html#i-o-latency>

# I/O replay

官方文档：<https://fio.readthedocs.io/en/latest/fio_doc.html#i-o-replay>

# 线程、进程、Job 同步相关参数

官方文档：<https://fio.readthedocs.io/en/latest/fio_doc.html#threads-processes-and-job-synchronization>

thead # 如果指定了此选项，则 fio 将使用 POSIX Threads 的函数 pthread_create(3) 创建线程来创建作业。Fio 默认使用 fork 创建 Job。i.e.使用进程来执行 Job

> 使用 thread 在一定程度上可以节省系统开销

# 认证参数

官方文档：<https://fio.readthedocs.io/en/latest/fio_doc.html#verification>

# Steady state

官方文档：<https://fio.readthedocs.io/en/latest/fio_doc.html#steady-state>

# 测量和报告相关参数

官方文档：<https://fio.readthedocs.io/en/latest/fio_doc.html#measurements-and-reporting>

# 错误处理相关参数

官方文档：<https://fio.readthedocs.io/en/latest/fio_doc.html#error-handling>

rwmixwrite=30 # 在混合读写的模式下，写占 30%

group_reporting # 关于显示结果的，汇总每个进程的信息

lockmem=1G # 只使用 1g 内存进行测试

nrfiles=8 # 每个进程生成文件的数量

numjobs=30 # 本次的测试线程为 30 个
