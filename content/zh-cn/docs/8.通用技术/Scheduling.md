---
title: Scheduling
linkTitle: Scheduling
weight: 20
---

# 概述

> 参考：
>
> - [Wiki, Scheduling(computing)](https://en.wikipedia.org/wiki/Scheduling_(computing))
> - [自己在科学笔记中对”生产与分配“的记录](https://github.com/DesistDaydream/notes-science/blob/main/%E7%A7%91%E5%AD%A6/%E7%A7%91%E5%AD%A6.md)

在计算机中，**Scheduling(调度)** 是分配资源来执行任务的行为。资源可以是 CPU、网络、etc. ，任务可以是 线程、进程、数据流、etc. 。

> [!Note] 调度是生产与分配中，分配的实践，是对[控制理论](https://zh.wikipedia.org/wiki/%E6%8E%A7%E5%88%B6%E7%90%86%E8%AE%BA)的实践（一个完善的调度系统可以间接理解为一种要无限趋近于[最优控制](https://zh.wikipedia.org/zh-hans/%E6%9C%80%E4%BC%98%E6%8E%A7%E5%88%B6)的系统）
>
> 但是由于现实的复杂性，通过纯粹的数学计算往往不能离最优控制还相去甚远。此时可以向系统中添加额外的组件（e.g. [Cache](/docs/8.通用技术/Cache.md)、负载均衡、虚拟化、etc.）来改善资源的分配。

调度活动由一种称为 **scheduler(调度器)** 的机制执行。调度器通常被设计为使所有计算机资源保持忙碌状态（如在负载均衡中），允许多个用户有效共享系统资源。

调度是计算本身的基础，也是计算机系统执行模型的内在部分；调度的概念使得在单个 CPU 上实现计算机多任务称为可能。

> [!Tip] 在英语体系里，有一个 [Schedule](https://en.wikipedia.org/wiki/Schedule) 单词，翻译成中文有”时间表“的含义，作为一种 [Time](/docs/8.通用技术/Time.md) 管理工具。Schedule 与 Scheduling 虽然意思差距很大，但是其内涵异曲同工。

# 调度系统设计精要

> 参考：
>
> - [调度系统设计精要](https://draveness.me/system-design-scheduler/)
