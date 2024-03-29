---
title: Hash
linkTitle: Hash
date: 2024-03-29T18:40
weight: 20
---

# 概述

> 参考：
>
> - [Wiki 类别，Hashing](https://en.wikipedia.org/wiki/Category:Hashing)
> - [Wiki，Hash function](https://en.wikipedia.org/wiki/Hash_function)

Hashing 是一种实现数据 Retrieval(检索) 的算法，有多种 Hashing 算法，比如

- [Consistent hashing](/docs/5.数据存储/Retrieval/Consistent%20hashing.md)

# Hash table

> 参考：
>
> - [Wiki，Hash table](https://en.wikipedia.org/wiki/Hash_table)

**Hash table(哈希表)** 也称为 hash map(哈希映射) 或 hash set(哈希集)，是一种实现关联数组的数据结构，也称为 dictionary(字典)，它是一种将键映射到值的抽象数据类型。哈希表使用哈希函数来计算索引（也称为哈希码）到桶或槽数组中，从中可以找到所需的值。在查找过程中，对键进行哈希处理，生成的哈希值指示相应值的存储位置。

理想情况下，哈希函数会将每个键分配给一个唯一的存储桶，但大多数哈希表设计都采用不完善的哈希函数，这可能会导致哈希冲突，即哈希函数为多个键生成相同的索引。此类冲突通常以某种方式进行调节
