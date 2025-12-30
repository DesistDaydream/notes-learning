---
title: Glossary
linkTitle: Glossary
weight: 101
---

# 概述

> 参考：
>
> -

## DataPath(数据路径)

网络数据在内核中进行网络传输时，所经过的所有点组合起来，称为数据路径。

## Ring Buffer

由网卡驱动程序创建的一种 Ring Buffer 的[数据结构](/docs/2.编程/计算机科学/Data%20type/Data%20structure.md)，保存在内存中，驱动程序会将这些内存地址告诉 [NIC](/docs/4.数据通信/Networking%20device/NIC.md) 硬件，以便 NIC 可以通过 DMA(直接内存访问) 将接收到的数据包直接写入这些地址，而无需 CPU 介入。

同时，驱动程序会消费 Ring Buffer 中的网络数据包，并将这些数据包封装成内核通用的 skb 接口，交给内核网络栈。

## Socket Buffer(简称 sk_buff 或 skb)

在内核代码中是一个名为 [**sk_buff**](https://www.kernel.org/doc/html/latest/networking/kapi.html#c.sk_buff) 的结构体。内核显然需要一个数据结构来储存报文的信息。这就是 skb 的作用。

sk_buff 结构自身并不存储报文内容，它通过多个指针指向真正的报文内存空间:

![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/efrsi8/1617849698535-471768e0-dcf8-4471-8dd2-605a1bc4e020.png)

sk_buff 是一个贯穿整个协议栈层次的结构，在各层间传递时，内核只需要调整 sk_buff 中的指针位置就行。

![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/efrsi8/1617849692989-54095177-b85c-449e-8c66-3b026e4925da.png)

## DEVICE(设备)

在内核代码中，是一个名为 [**net_device**](https://www.kernel.org/doc/html/latest/networking/kapi.html#c.net_device) 的结构体。一个巨大的数据结构，描述一个网络设备的所有 属性、数据 等信息。


