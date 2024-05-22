---
title: "Cache"
linkTitle: "Cache"
date: "2023-08-06T16:15"
weight: 20
---

# 概述

> 参考：
>
> - [WIki，Cache(computing)](https://en.wikipedia.org/wiki/Cache_(computing))

**Cache(缓存)** 在 [Computer](/docs/0.计算机/Computer.md)(计算机) 中，是指一种用于存储数据的软件或硬件，以便可以更快的满足未来对该数据的请求；存储在缓存中的数据可能是早期计算的结果或存储在其他地方的数据的副本。作为缓存的组件的最低要求就是速度够快，避免对存储原始数据的组件产生过大压力。

缓存是计算机系统中最基本的性能改进技术之一，一次又一次地用于让“常见的情况更快”[^2]。硬件缓存背后的思想是利用指令和数据引用的 **Locality(局部性原理)**。通常有两种局部性: **temporal locality(时间局部性)** 和 **spatial locality(空间局部性)**。时间局部性是指，最近访问过的指令或数据项可能很快会再次访问。想想循环中的循环变量或指令，它们被多次反复访问。空间局部性是指，当程序访问内存地址 x 时，可能很快会访问邻近 x 的内存。想想遍历某种数组，访问一个接一个的元素。当然，这些性质取决于程序的特点，并不是绝对的定律，而更像是一种经验法则。[^1]

硬件缓存，无论是指令、数据还是地址转换（如 TLB），都利用了局部性，在小而快的芯片内存储器中保存一份内存副本。处理器可以先检查缓存中是否存在就近的副本，而不是必须访问（缓慢的）内存来满足请求。如果存在，处理器就可以很快地访问它（例如在几个 CPU 时钟内），避免花很多时间来访问内存（好多纳秒）。[^1]

你可能会疑惑：既然像 TLB 这样的缓存这么好，为什么不做更大的缓存，装下所有的数据？可惜的是，这里我们遇到了更基本的定律，就像物理定律那样。如果想要快速地缓存，它就必须小，因为光速和其他物理限制会起作用。大的缓存注定慢，因此无法实现目的。所以，我们只能用小而快的缓存。剩下的问题就是如何利用好缓存来提升性能。[^1]

[^1]: 操作系统导论 - 19.2
[^2]: “Computer Architecture: A Quantitative Approach” John Hennessy and David Patterson - Morgan-Kaufmann, 2006

# 写入策略

https://en.wikipedia.org/wiki/Cache_(computing)#Writing_policies

由于缓存中的数据不能长久保持，所以要根据一定的策略将数据写入到后端存储，具体的写入行为，由 **Writing policies(写入策略)** 决定

- **Write-throught(直写)** # 向缓存和后端存储同步写入
- **Write-back(回写)** # 先写入缓存，直到某一条件触发后，再写入到后端存储。

问 Bing 回写是什么的回答：

感谢您提供更多的信息。我找到了一些关于计算机领域的回写（write back）的资料，希望对您有用。😊

回写（write back）是一种缓存（cache）的写入策略，它指的是在数据更新时只写入缓存，而不写入后端存储（如内存或硬盘）。只有当缓存中的数据被替换出去时，才会把修改过的数据写回后端存储。这样可以提高数据写入的速度，但也有可能造成数据不一致的问题。1 2

回写（write back）通常与写分配（write allocate）一起使用，这意味着当缓存未命中时，会先把数据从后端存储加载到缓存中，然后再更新缓存。这样可以减少对后端存储的访问次数，同时也为后续对同一位置的写入或读取提供便利。3

回写（write back）与另一种缓存写入策略透写（write through）相对，透写指的是在数据更新时同时写入缓存和后端存储，这样可以保持数据的一致性，但也会降低数据写入的速度。透写通常与无写分配（no write allocate）一起使用，这意味着当缓存未命中时，只更新后端存储，而不加载数据到缓存中。3 4

如果您想了解更多关于回写（write back）和透写（write through）的区别和优缺点，请参考以下链接：

1: [计算机的Cache和Memory访问时Write-back，Write-through及write allocate的区别 - 国静德远 - 博客园](https://www.cnblogs.com/guojingdeyuan/p/7626983.html)
2: [cache 写回策略_写回cache_xingzhe22222的博客-CSDN博客](https://blog.csdn.net/xingzhe22222/article/details/81988101)
3: [透写和回写缓存（Write Through and Write Back in Cache）](https://zhuanlan.zhihu.com/p/571429282)
4: [Cache写机制：Write-through与Write-back - CSDN博客](https://blog.csdn.net/zhangj95/article/details/81199272)

