---
title: Vector
linkTitle: Vector
weight: 1
---

# 概述

> 参考：
>
> - [GitHub 项目，vectordotdev/vector](https://github.com/vectordotdev/vector)
> - [官网](https://vector.dev/)
> - https://www.cnblogs.com/ulricqin/p/17762086.html

Vector 是一种高性能的用于可观测性的 [DataPipeline](/docs/6.可观测性/DataPipeline/DataPipeline.md)(数据管道)，让用户能够控制其可观测性数据。收集、转换 所有日志、指标和跟踪，并将其路由到任意 Vendor 以及明天可能需要的其他 Vendor。

> Notes: Vendor 指使用这些数据的地方，e.g. 数据库、Web 前端、etc. 这些地方都可以对外提供数据，所以称为 Vendor(供应商)，就像数据供应商似的。

Datadog 在 2021 年左右收购了 Vector。Vector 通常用作 ELK 生态中 logstash 的替代品。

Vector 开箱即用，默认支持 [ClickHouse](/docs/5.数据存储/数据库/关系数据/ClickHouse/ClickHouse.md)、etc.

# Vector 架构

![](https://raw.githubusercontent.com/vectordotdev/vector/refs/heads/master/website/static/img/data-model-event.svg)

Vector 将数据通道抽象为 3 部分组件：

- **Sources** # 将可观测性数据源中的数据收集或接收到 Vector 中
- **Transforms** # 在可观测性数据通过拓扑时操纵或更改该数据。（拓扑可以理解为一种网状结构，由处理数据的多个节点组成）
- **Sinks** # 将可观测性数据从 Vector 向前发送到外部服务或目的地

# 部署角色

Vector 可以部署为两个角色，既可以作为数据采集的 Agent，也可以作为数据聚合、路由的 Aggregator，架构示例如下：

![https://vector.dev/docs/setup/deployment/](https://vector.dev/img/deployment.png)

## Agent

## Aggregator

# Vector 部署

> 参考：
>
> - [官方文档，Setup - 安装](https://vector.dev/docs/setup/installation/)

## 包管理器

## 容器

## 二进制文件

# Vector 关联文件与配置

**/var/lib/vector/** # 持久保存 Vector 状态的目录。e.g. 磁盘缓冲、文件检查点、etc. 。可以通过配置文件的 .data_dir 字段指定。

**/etc/vector/vector.yaml** # 已经弃用的默认配置文件。由于 Vector 灵活的设计，可以加载多个配置文件，一般情况都是手动使用 --config-dir 指定配置文件目录。这样更利于配置管理。

详见 [Vector Configuration](/docs/6.可观测性/DataPipeline/Vector/Vector%20Configuration.md)

# API

> 参考：
>
> - [官方文档，参考 - API](https://vector.dev/docs/reference/api/)

Vector 提供 [GraphQL](/docs/2.编程/API/GraphQL.md) API

# 基础用例

[公众号，实战 Vector：开源日志和指标采集工具](https://mp.weixin.qq.com/s/o6bqzJt1M_DNn027Nc91fQ)

模拟 Prometheus 的 Exporter。TODO: 具体都能采集到什么？

```toml
# sample.toml
[sources.prom]
type = "prometheus_scrape"
endpoints = [ "http://localhost:9100/metrics" ]
```