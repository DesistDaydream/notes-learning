---
title: CPU
linkTitle: CPU
weight: 1
---

# 概述

> 参考：
>
> - [极客时间，Linux 性能优化实战-03 基础篇：经常说的 CPU 上下文切换是什么意思](https://time.geekbang.org/column/article/69859)
> - [LinuxPerformance 博客，进程切换：自愿与强制](http://linuxperf.com/?p=209)

在 [Linux Kernel](/docs/1.操作系统/Kernel/Linux%20Kernel/Linux%20Kernel.md) 中，CPU 的管理，绝大部分时间都是在进行任务的调度，所以很多时候也称为**调度管理**。

## CPU 多线程、并发、并行 概念

Node：在这里时间片只是一种描述，理解 CPU 的并行与并发概念就好

1、CPU 时间分片、多线程？

如果线程数不多于 CPU 核心数，会把各个线程都分配一个核心，不需分片，而当线程数多于 CPU 核心数时才会分片。

2、并发和并行的区别

- 并发：当有多个线程在操作时,如果系统只有一个 CPU，把 CPU 运行时间划分成若干个时间片,分配给各个线程执行，在一个时间段的线程代码运行时，其它线程处于挂起状态。这种方式我们称之为 **Concurrent(并发)**。并发=间隔发生
- 并行：当系统有一个以上 CPU 时,则线程的操作有可能非并发。当一个 CPU 执行一个线程时，另一个 CPU 可以执行另一个线程，两个线程互不抢占 CPU 资源，可以同时进行，这种方式我们称之为 **Parallel(并行)**。 并行=同时进行

区别：并行是指两个或者多个事件在同一时刻发生；而并发是指两个或多个事件在同一时间间隔内发生。

并行是同时做多件事情。

并发表示同时发生了多件事情，通过时间片切换，哪怕只有单一的核心，也可以实现“同时做多件事情”这个效果。

根据底层是否有多处理器，并发与并行是可以等效的，这并不是两个互斥的概念。

举个我们开发中会遇到的例子，我们说资源请求并发数达到了 1 万。这里的意思是有 1 万个请求同时过来了。但是这里很明显不可能真正的同时去处理这 1 万个请求的吧！

如果这台机器的处理器有 4 个核心，不考虑超线程，那么我们认为同时会有 4 个线程在跑。也就是说，并发访问数是 1 万，而底层真实的并行处理的请求数是 4。如果并发数小一些只有 4 的话，又或者你的机器牛逼有 1 万个核心，那并发在这里和并行一个效果。也就是说，并发可以是虚拟的同时执行，也可以是真的同时执行。而并行的意思是真的同时执行。

结论是：并行是我们物理时空观下的同时执行，而并发则是操作系统用线程这个模型抽象之后站在线程的视角上看到的“同时”执行。

### time slice(时间片) 概念

> 参考：<https://en.wikipedia.org/wiki/Preemption_(computing)#Time_slice>

The period of time for which a process is allowed to run in a preemptive multitasking system is generally called the _time slice_ or _quantum_.

**time slice(时间片)** 是一个程序运行在[抢占式多任务系统](<https://en.wikipedia.org/wiki/Preemption_(computing)>)中的一段时间。也可以称为 quantum(量子)。

## CPU 使用率概念

CPU 不像硬盘、内存，并不具备逻辑上数量、大小、空间之类的概念。只要使用 CPU，就是使用了这个 CPU 的全部，也就无法通过大小之类的概念来衡量一个 CPU，所以我们日常所说的 CPU 的使用率 ，实际上是指的在一段时间范围内，CPU 执行 **Tasks(任务)** 时，<font color="#ff0000">**花费时间的百分比**</font>。比如 60 分钟内，一颗 CPU 执行各种任务花费了 6 分钟，则 CPU 在这一小时时间内的使用率为 10%。

> 上文说的 **Tasks(任务)**，即会指系统中的进程、线程，也代表各种硬件去请求 CPU 执行的各种事情，比如网卡接收到数据，就会告诉 CPU 需要处理(i.e.中断)。

在 Linux 系统中，CPU 的使用率一般可分为 4 大类：

1. User Time(用户进程运行时间)
2. System Time(系统内核运行时间)
3. Idle Time(空闲时间)
4. Steal Time(被抢占时间)

除了 Idle Time 外，CPU 在其余时间都处于工作运行状态

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/srucoz/1616168021555-68fba1de-f5d5-462d-bef6-a78b476521ad.png)

通常而言，我们泛指的整体 CPU 使用率为 User Time 和 Systime 占比之和(例如 tsar 中 CPU util)，即：

$$CPUutil = \frac{(UserTime + SystemTime)}{UserTime + SystemTime + IdleTime + StealTime}$$

为了便于定位问题，大多数性能统计工具都将这 4 类时间片进一步扩展成了 8 类，如下图，是在 top 命令的 man 手册中对 CPU 使用率的分类。

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/srucoz/1616168021546-ebe53556-f50b-49f2-8477-c10cf2b8f2f5.png)

- us：用户进程空间中未改变过优先级的进程占用 CPU 百分比
- sy：内核空间占用 CPU 百分比
- ni：用户进程空间内改变过优先级的进程占用 CPU 百分比
- id：空闲时间百分比
- wa：等待磁盘 I/O 操作的时间百分比
- hi：硬中断时间百分比
- si：软中断时间百分比
- st：虚拟化时被其余 VM 窃取时间百分比

这 8 类分片中，除 wa 和 id 外，其余分片 CPU 都处于工作态。

# 调度算法

> 详见：[CPU 调度算法](/docs/1.操作系统/Kernel/CPU/调度算法.md)

首先明确一个概念：**Task(任务)**，一个进程从处理到结束就算一个任务，处理网卡收到的数据包也算一个任务。一般来说，CPU 就是在处理一个个的 **Task(任务)**，并度过其一生。

在 Linux 内核中，进程和线程都是用 tark_struct 结构体表示的，区别在于线程的 tark_struct 结构体里部分资源是共享了进程已创建的资源，比如内存地址空间、代码段、文件描述符等，所以 Linux 中的线程也被称为轻量级进程，因为线程的 tark_struct 相比进程的 tark_struct 承载的 资源比较少，因此以「轻」得名。

一般来说，没有创建线程的进程，是只有单个执行流，它被称为是主线程。如果想让进程处理更多的事情，可以创建多个线程分别去处理，但不管怎么样，它们对应到内核里都是 tark_struct。

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/srucoz/1616168021545-596ecf70-ac19-4620-8845-bfe72ef7bdce.jpeg)

所以，Linux 内核里的调度器，调度的对象就是 tark_struct，接下来我们就把这个数据结构统称为任务。

在 Linux 系统中，根据任务的优先级以及响应要求，主要分为两种，其中优先级的数值越小，优先级越高：

- 实时任务，对系统的响应时间要求很高，也就是要尽可能快的执行实时任务，优先级在 0~99 范围内的就算实时任务；
- 普通任务，响应时间没有很高的要求，优先级在 100~139 范围内都是普通任务级别；

也就是说，在 LInux 内核中，实时任务总是比普通任务的优先级要高。
