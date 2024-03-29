---
title: Fiber-optic splitter
linkTitle: Fiber-optic splitter
date: 2024-02-22T13:32
weight: 4
---

# 概述

> 参考：
>
> - [Wiki，Fiber-optic splitter](https://en.wikipedia.org/wiki/Fiber-optic_splitter)
> - [Wiki，Beam_splitter](https://en.wikipedia.org/wiki/Beam_splitter)
> - [百度百科，光纤分路器](https://baike.baidu.com/item/%E5%85%89%E7%BA%A4%E5%88%86%E8%B7%AF%E5%99%A8)
> - [百度百科，分光器](https://baike.baidu.com/item/%E5%88%86%E5%85%89%E5%99%A8)

**Fiber-optic splitter(分光器、光纤分路器)** 是基于石英基板的集成波导光功率分配装置，类似于同轴电缆传输系统。光网络系统使用耦合到分支分配的光信号。光纤分路器是光纤链路中最重要的无源器件之一。它是一种具有多个输入和输出端子的光纤汇接设备，特别适用于无源光网络（EPON、GPON、BPON、FTTX、FTTH等）连接主配线架和终端设备并对光信号进行分支。

> 有的地方也会用 Optical splitter 或 Beam splitter

分光器仅仅是单纯针对物理链路进行按比例分光，不会涉及上层的流量识别和处理，通常按照 2/8 比例分光，80% 的光能量在原始链路传输以保证通信不受影响；20% 的光能量作为复制的流量发送给下一个设备以进一步处理（可以分析流量中恶意数据等）。通常这 20% 的光需要经过 [Optical amplifier](/docs/4.数据通信/Networking%20device/Optical%20amplifier.md)(光放大器) 才能被其他设备正常使用。

# 分类

> #网络硬件
