---
title: Computing platform
linkTitle: Computing platform
created: 2026-04-29T12:58
weight: 1
---

# 概述

> 参考：
>
> - [知乎，全球AI计算平台对比：英伟达CUDA、华为CANN和海光ROCm](https://zhuanlan.zhihu.com/p/1946846455777191464)

**AI Computing platform(人工智能计算平台)** 现阶段是我自己造的词，用来表示 CUAD, CANN, etc. 相关生态、框架、程序。

通过 Computing platform，可以让 [Torch](/docs/12.AI/科学计算/Torch.md), 昇思, etc. 框架，使用专用的计算设备（e.g. [GPU](/docs/0.计算机/GPU/GPU.md), [NPU](/docs/0.计算机/NPU.md)），避免了只能用 CPU 进行缓慢计算的尴尬。

当前可用的计算平台：

- CUDA
- [CANN](/docs/12.AI/Computing%20platform/CANN.md)
- ROCm

# CUDA

> 参考：
>
> - [Wiki, CUDA](https://en.wikipedia.org/wiki/CUDA)

**Compute Unified Device Architecture( 统一计算设备架构，简称 CUDA)**