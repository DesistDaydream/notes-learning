---
title: 0_操作系统
weight: 1
---

# 概述

> 参考：
> - [Wiki.Operating System](https://en.wikipedia.org/wiki/Operating_system)
> - [高教书苑,操作系统原理(2020 年版)-全国计算机等级考试四级教程](https://ebook.hep.com.cn/ebooks/index.html#/read?id=693031822086377472)
> - [高教书苑,操作系统原理与实训教程(第 3 版)](https://ebook.hep.com.cn/ebooks/index.html#/read?id=685438574224478208)
> - [公众号-码农的荒岛求生，没有操作系统，程序可以运行起来吗？](https://mp.weixin.qq.com/s/sEv8_o2FABGVtULOGUv3ZQ)

**Operating System(操作系统，简称 OS)** 是一个[系统软件](https://en.wikipedia.org/wiki/System_software)，由很多程序模块组成，用以组织计算机的工作流程、有效得控制和管理计算机系统的各类(软件/硬件)资源，并向用户提供各种服务功能(操作接口)，使用户能够灵活、方便、有效地使用计算机。

白话说，**OS 是计算机系统的**`**管家**`**，是硬件最亲密的**`**伙伴**`**，是人机之间的**`**桥梁**`**，是其他应用软件的**`**基石**`**。**

计算机起初是为了代替人类计算产生的，一台设备只能执行一个程序。如果一台计算机上需要同时运行三个程序，那么会有如下问题产生：

1. 三道程序在内存中如何存放？
2. 什么时候让某个程序占用 CPU？
3. 怎样有序地输出各个程序的运算结果？

对这些问题的解决都必须求助于操作系统。也就是说操作系统必须对内存进行管理，对 CPU 进行管理，对外设机型管理，对存放在磁盘上的文件更是要精心组织和管理。不仅如此，操作系统对这些资源进行管理的基础上，还要给用户提供良好的接口，以便用户可能在某种程度上使用或者操纵这些资源。因此，从操作系统设计者的角度考虑，一个操作系统必须包含以下几个部分

1. 操作系统接口
2. CPU 管理
3. 内存管理
4. 设备管理
5. 文件管理

操作系统就是一个大型的软件而已，与运行在操作系统之上的各种程序基本一样。

1. 操作系统=内核+系统程序
2. 系统程序=编译环境+API(应用程序接口)+AUI(用户接口)
3. 编译环境=编译程序+连接程序+装载程序
4. API=系统调用+语言库函数(C、C++、Java 等)
5. AUI=shell+系统服务例程+应用程序(浏览器、字处理、编辑器等)

一个名为 test 程序执行的简单过程：

1. 用户通过交互界面(shell)告诉操作系统执行 test 程序
2. 操作系统通过文件名在磁盘找到该程序
3. 检查可执行代码首部，找出代码和数据存放的地址
4. 文件系统找到第一个磁盘块
5. 操作系统建立程序的执行环境
6. 操作系统把程序从磁盘装入内存，并跳到程序开始出执行
7. 操作系统检查字符串的位置是否正确
8. 操作系统找到字符串被送往的设备
9. 操作系统将字符串送往输出设备窗口系统确定这是一个合法的操作，然后将字符串转换成像素
10. 窗口系统将像素写入存储映像区
11. 视频硬件将像素表示转换成一组模拟信号控制显示器
12. 显示器发射电子束，最后在屏幕上看到程序执行的结果。

Linux 系统结构详解
![](https://notes-learning.oss-cn-beijing.aliyuncs.com/nm71yz/1616168594662-a90f0c59-7c9b-49ee-ba91-d7065227bfcd.png)

Linux 系统一般有 4 个主要部分：

1. 内核
2. shell、进程
3. 文件系统
4. 应用程序。

内核、shell 和文件系统一起形成了基本的操作系统结构，它们使得用户可以运行程序、管理文件并使用系统。部分层次结构

# 操作系统分类

对 OS 进行分类的方法有很多，从不同的角度可以得到不同的划分。

## 按照计算机硬件和结构与规模分类

- 大型机操作系统
- 中型机操作系统
- 小型机操作系统
- 微型机操作系统
- **(Networking)网络** 操作系统
- **(Embedded)嵌入式** 操作系统
- 等等

## 按照系统所能同时响应的用户与任务分类

- **(Single-user/Multi-user)单用户 或 多用户** 操作系统
  - 单用户操作系统不具有区分用户的功能，但可以允许多个程序串联运行。[\[8\]](https://en.wikipedia.org/wiki/Operating_system#cite_note-8)阿[多用户](https://en.wikipedia.org/wiki/Multi-user)操作系统延伸的与设施识别过程和资源，例如磁盘空间，属于多个用户的多任务处理的基本概念，并且该系统允许多个用户在同一时间与系统进行交互。分时操作系统调度任务以有效使用系统，并且还可以包括用于将处理器时间，大容量存储，打印和其他资源的成本分配给多个用户的计费软件。
- **(Single-tasking/Multi-tasking)单任务 或 多任务** 操作系统
  - 单任务系统一次只能运行一个程序，而[多任务](https://en.wikipedia.org/wiki/Computer_multitasking)操作系统允许[并发](https://en.wikipedia.org/wiki/Concurrent_computing)运行多个程序。这是通过 **Time-shring(**[**分时**](https://en.wikipedia.org/wiki/Time-sharing)**) **实现的，其中可用的处理器时间在多个进程之间分配。这些进程每个都在一个 **time slices(**[**时间片**](https://en.wikipedia.org/wiki/Time_slice)**)** 中被操作系统的任务计划子系统反复中断。多任务可以以抢占式和合作式为特征。在[抢占式](<https://en.wikipedia.org/wiki/Preemption_(computing)>)多任务处理中，操作系统会减少[CPU](https://en.wikipedia.org/wiki/Central_processing_unit)时间，并为每个程序分配一个插槽。[类似于 Unix 的](https://en.wikipedia.org/wiki/Unix-like) 操作系统，例如 [Solaris](<https://en.wikipedia.org/wiki/Solaris_(operating_system)>) 和[Linux](https://en.wikipedia.org/wiki/Linux) 以及非 Unix 操作系统(例如[AmigaOS)均](https://en.wikipedia.org/wiki/AmigaOS)支持抢占式多任务处理。通过依靠每个过程以定义的方式向其他过程提供时间来实现协作式多任务处理。

通过 用户 和 任务 的组合，可以组合出三种：单用户单任务、单用户多任务、多用户多任务。

## 按照系统处理任务的方式分类(更加广泛的分类)

- **多道批处理**操作系统 #
- **(Time-shring)分时 **操作系统 #
  - UNIX
- **(Real-time)实时 **操作系统 #
- **(Distributed)分布式 **操作系统 #
- 等等

## 其他分类

- 个人操作系统
- 等等

# 发展史

[CP/M](https://en.wikipedia.org/wiki/CP/M) 第一个微型机的操作系统

# 查看操作系统信息

**/etc/os-release**
**/etc/issue**
