---
title: Redis 数据类型
linkTitle: Redis 数据类型
date: 2024-06-14T13:46
weight: 20
---

# 概述

> 参考：
>
> - [官方文档，数据类型](https://redis.io/topics/data-types/)
>   - https://redis.io/docs/latest/develop/data-types/


string(字符串)，hash(哈希)，list(列表)，set(集合) 及 zset(sorted set：有序集合)。

后面增加了：

Bit arrays (或者说 simply bitmaps)

在 2.8.9 版本添加了 HyperLogLog 结构

一般情况下，每种类型的数据都有与之相关的 [Redis CLI](/docs/5.数据存储/数据库/键值数据/Redis/Redis%20CLI/Redis%20CLI.md) 中的子命令对应进行处理