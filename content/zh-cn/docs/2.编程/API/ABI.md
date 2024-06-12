---
title: ABI
linkTitle: ABI
date: 2024-05-23T22:34
weight: 20
---

# 概述

> 参考：
>
> - [Wiki，Application binary interface](https://en.wikipedia.org/wiki/Application_binary_interface)
> - https://www.jianshu.com/p/bd77c842f281

**Application binary interface(应用程序二进制接口，简称 ABI)** 是两个二进制程序模块之间的接口。通常，这俩模块一个是库或操作系统设施，另一个是用户运行的程序。

ABI 定义如何以机器代码访问数据结构或计算例程，这是一种低级、依赖于硬件的格式。相比之下，应用程序编程接口 ([API](/docs/5.数据存储/存储/存储的基础设施架构/分布式存储/Ceph/API.md)) 在源代码中定义这种访问，这是一种相对高级、独立于硬件、通常是人类可读的格式。 ABI 的一个常见方面是调用约定，它确定如何将数据作为计算例程的输入提供或从计算例程读取数据。 x86 调用约定就是这样的示例。

保持一个稳定的 ABI 要比保持稳定的 API 要难得多。比如，在内核中 `int register_netdevice(struct net_device *dev)` 这个内核函数原型基本上是不会变的，所以保持这个 API 稳定是很简单的，但它的 ABI 就未必了，就算是这个函数定义本身没变，即 API 没变，而 struct net_device 的定义变了，里面多了或者少了某一个字段，它的 ABI 就变了，你之前编译好的二进制模块就很可能会出错了，必须重新编译才行。
