---
title: 进程与CPU核心的绑定
linkTitle: 进程与CPU核心的绑定
weight: 20
---

# 概述

> 参考:
>
> - [Linux 有问必答：如何知道进程运行在哪个 CPU 内核上？](https://linux.cn/article-6307-1.html?pr)
>   - https://www.xmodulo.com/cpu-core-process-is-running.html

问题：我有个 Linux 进程运行在多核处理器系统上。怎样才能找出哪个 CPU 内核正在运行该进程？

当你在 多核 NUMA 处理器上运行需要较高性能的 HPC（高性能计算）程序或非常消耗网络资源的程序时，CPU/memory 的亲和力是限度其发挥最大性能的重要因素之一。在同一 NUMA 节点上调度最相关的进程可以减少缓慢的远程内存访问。像英特尔 Sandy Bridge 处理器，该处理器有一个集成的 PCIe 控制器，你可以在同一 NUMA 节点上调度网络 I/O 负载（如网卡）来突破 PCI 到 CPU 亲和力限制。

作为性能优化和故障排除的一部分，你可能想知道特定的进程被调度到哪个 CPU 内核（或 NUMA 节点）上运行。

这里有几种方法可以**找出哪个 CPU 内核被调度来运行给定的 Linux 进程或线程**。

## 方法一

如果一个进程使用 taskset 命令明确的被固定（pinned）到 CPU 的特定内核上，你可以使用 taskset 命令找出被固定的 CPU 内核：

$ taskset -c -p

例如, 如果你对 PID 5357 这个进程有兴趣:

$ taskset -c -p 5357

pid 5357's current affinity list: 5

输出显示这个过程被固定在 CPU 内核 5 上。

但是，如果你没有明确固定进程到任何 CPU 内核，你会得到类似下面的亲和力列表。

pid 5357's current affinity list: 0-11

输出表明该进程可能会被安排在从 0 到 11 中的任何一个 CPU 内核。在这种情况下，taskset 不能识别该进程当前被分配给哪个 CPU 内核，你应该使用如下所述的方法。

## 方法二

ps 命令可以告诉你每个进程/线程目前分配到的 （在“PSR”列）CPU ID。

$ ps -o pid,psr,comm -p

PID PSR COMMAND

5357 10 prog

输出表示进程的 PID 为 5357（名为”prog”）目前在 CPU 内核 10 上运行着。如果该过程没有被固定，PSR 列会根据内核可能调度该进程到不同内核而改变显示。

## 方法三

top 命令也可以显示 CPU 被分配给哪个进程。首先，在 top 命令中使用 P 选项。然后按“f”键，显示中会出现 “Last used CPU” 列。目前使用的 CPU 内核将出现在 “P”（或“PSR”）列下。

```bash
top -p 5357 -H
```

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/hsfm16/1616167411830-a73b375d-c071-4953-a44d-eb0390f44258.jpeg)
相比于 ps 命令，使用 top 命令的好处是，你可以连续监视随着时间的改变， CPU 是如何分配的。

## 方法四

另一种来检查一个进程/线程当前使用的是哪个 CPU 内核的方法是使用 htop 命令。

从命令行启动 htop。按 键，进入”Columns”，在”Available Columns”下会添加 PROCESSOR。

每个进程当前使用的 CPU ID 将出现在“CPU”列中。

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/hsfm16/1616167411836-ec26f183-e19b-4050-a5bb-5c4b033215cd.jpeg)

请注意，所有以前使用的命令 taskset，ps 和 top 分配 CPU 内核的 IDs 为 0，1，2，…，N-1。然而，htop 分配 CPU 内核 IDs 从 1 开始（直到 N）。
