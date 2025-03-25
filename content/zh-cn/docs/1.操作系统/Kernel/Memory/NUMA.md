---
title: NUMA
linkTitle: NUMA
weight: 20
---

# 概述

> 参考：
>
> - [Linux Kernel 文档，管理员指南 - 内存管理 - 概念 - Nodes](https://www.kernel.org/doc/html/latest/admin-guide/mm/concepts.html#nodes)
> - [Linux Kernel 文档，内存管理 - NUMA](https://www.kernel.org/doc/html/latest/mm/numa.html)
> - [Linux Kernel 文档，管理员指南 - 内存管理 - NUMA 内存策略](https://www.kernel.org/doc/html/latest/admin-guide/mm/numa_memory_policy.html)
> - [Wiki, Non-uniform memory access](https://en.wikipedia.org/wiki/Non-uniform_memory_access)

**Non-uniform memory access(非均匀内存访问，简称 NUMA)** 是一种用于多处理结构的计算机的内存设计，其中内存访问时间取决于相对于处理器的内存位置。在 NUMA 下，处理器可以比非本地内存（另一个处理器的本地内存或处理器之间共享的内存）更快地访问自己的本地内存。 NUMA 的优势仅限于特定的工作负载，特别是在数据通常与某些任务或用户密切相关的服务器上。

![《深入浅出PDPK》-2.9 | 900](https://notes-learning.oss-cn-beijing.aliyuncs.com/os/kernel/memory/numa_1.png)

Linux 将系统的硬件资源按照抽象的 **Nodes(节点)** 概念进行划分。Linux 将 Nodes 映射到硬件平台的物理 Cells(单元) 上，抽象出某些架构的一些细节。与物理单元一样，软件 Nodes 可能包含 0 个或多个 CPU、内存和/或 IO 总线。而且，与对更远程单元的访问相比，对“较近” Node 上的存储器的存储器访问通常会经历更快的访问时间和更高的有效带宽。
