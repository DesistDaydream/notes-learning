---
title: BPF 流量控制
linkTitle: BPF 流量控制
weight: 1
---

# 概述

> 参考：
> 
> - [Kernel 网络官方文档，Linux Socket Filtering aka Berkeley Packet Filter](https://www.kernel.org/doc/html/latest/networking/filter.html#)

# 学习资料

[\[译\] 利用 eBPF 支撑大规模 K8S Service](https://mp.weixin.qq.com/s/KnNcM2OaBqOgfVDghaHy8g)

[为容器时代设计的高级 eBPF 内核特性（FOSDEM, 2021）](https://mp.weixin.qq.com/s/ZCprEJi9zrHxRSO1XNRsqQ)

# BPF 在网络领域的实现

> 参考：
>
> - [arthurchiao.art 的文章](http://arthurchiao.art/index.html)
>   - [\[译\] 为容器时代设计的高级 eBPF 内核特性（FOSDEM, 2021)](http://arthurchiao.art/blog/advanced-bpf-kernel-features-for-container-age-zh/)

## eBPF 架构

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/gt3mhv/1617851739655-d700e4b4-2880-4a5f-afbf-182baaba18f3.png)

eBPF 能够让你**在内核中**创建新的 DataPath。eBPF 就是内核本身的代码，想象空间无限，并且热加载到内核；换句话说，一旦加载到内核，内核的行为就变了。在 eBPF 之前，改变内核行为这件事情，只能通过修改内核再重新编译，或者开发内 核模块才能实现。

由于上述原因，真正的 eBPF，应该是基于 eBPF 实现的数据路径，由于 eBPF 可以修改内核，所以可以在内核创建新的类似 Netfilter 的 Hook 点，以便跳过复杂的 Netfilter。甚至可以直接在网卡驱动中运行 eBPF 代码，而无需将数据包送到复杂的协议栈进行处理。
