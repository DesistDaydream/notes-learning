---
title: TLB
linkTitle: TLB
date: 2024-05-22T15:27
weight: 20
tags:
  - CPU
---

# 概述

> 参考：
>
> - [Wiki，Translation lookaside buffer](https://en.wikipedia.org/wiki/Translation_lookaside_buffer)

**translation lookaside buffer(转换后被缓冲区，简称 TLB)** 用于减少访问用户内存位置所需的时间，是 MMU 的一部分

TLB 是一种 [Cache](/docs/8.通用技术/Cache.md) 功能，CPU 在寻址时，若是从 TLB 没有查到虚拟内存地址与物理内存地址的对应关系，称为 **未命中**，此时会去查找常规的页表。

