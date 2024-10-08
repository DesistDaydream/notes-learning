---
title: Vector
linkTitle: Vector
date: 2024-10-08T11:31
weight: 20
---

# 概述

> 参考：
>
> - [GitHub 项目，vectordotdev/vector](https://github.com/vectordotdev/vector)
> - https://www.cnblogs.com/ulricqin/p/17762086.html

Vector 是一种高性能的用于可观测性的 [DataPipeline](/docs/6.可观测性/DataPipeline/DataPipeline.md)(数据管道)，让用户能够控制其可观测性数据。收集、转换 所有日志、指标和跟踪，并将其路由到任意 Vendor 以及明天可能需要的其他 Vendor。

> Notes: Vendor 指使用这些数据的地方，e.g. 数据库、Web 前端、etc. 这些地方都可以对外提供数据，所以称为 Vendor(供应商)，就像数据供应商似的。

Datadog 在 2021 年左右收购了 Vector。Vector 通常用作 ELK 生态中 logstash 的替代品。

Vector 开箱即用，默认支持 [ClickHouse](/docs/5.数据存储/数据库/关系数据/ClickHouse.md)、etc.

# Vector 架构

Vector 可以部署为两个角色，既可以作为数据采集的 Agent，也可以作为数据聚合、路由的 Aggregator，架构示例如下：

![https://www.cnblogs.com/ulricqin/p/17762086.html](https://download.flashcat.cloud/ulric/20230927153626.png)

## Agent

## Aggregator

