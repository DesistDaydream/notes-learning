---
title: Computing Virtualization
linkTitle: Computing Virtualization
weight: 2
---

# 概述

> 参考：

CPU 虚拟化(vCPU=virtual CPU # 虚拟 CPU)

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/gd8p4x/1616124392775-cccdd43b-0f21-4877-8c78-c4d6c7352728.png)

使用如下命令可以查看该 CPU 是否支持虚拟化

egrep -o '(vmx|svm)' /proc/cpuinfo

如果有输出 vmx 或者 svm，就说明当前的 CPU 支持 KVM。CPU 厂商 Intel 和 AMD 都支持虚拟化了，除非是非常老的 CPU。

在 CUP 虚拟化的图片中，宿主机有两个物理 CPU，上面起了两个虚机 VM1 和 VM2。 VM1 有两个 vCPU，VM2 有 4 个 vCPU。可以看到 VM1 和 VM2 分别有两个和 4 个线程在两个物理 CPU 上调度。

虚机的 vCPU 总数可以超过物理 CPU 数量，这个叫 CPU overcommit（超配）。 KVM 允许 overcommit，这个特性使得虚机能够充分利用宿主机的 CPU 资源，但前提是在同一时刻，不是所有的虚机都满负荷运行。 当然，如果每个虚机都很忙，反而会影响整体性能，所以在使用 overcommit 的时候，需要对虚机的负载情况有所了解，需要测试。
