---
title: LogQL
linkTitle: LogQL
weight: 1
---

# 概述

> 参考：
>
> - date: "2024-12-10T08:47"
> - [官方文档，查询](https://grafana.com/docs/loki/latest/query/)
> - [公众号，Loki 查询语言 LogQL 使用](https://mp.weixin.qq.com/s/0dXT0fIreZk6_4ZL4S8lHg)

**Log Query Language(日志查询语言，简称 LogQL)** 受 PromQL 启发，可以看作是分布式的 grep 命令，用来从汇总的日志源筛选日志。LogQL 通过 Labels(标签) 和 Operators(运算符) 进行过滤。

LogQL 查询有两种类型：

- **Log Queries(日志查询)** # 根据查询语句返回日志条目，每行是一条日志。
- **Metric Queries(指标查询)** # 用于扩展日志查询并根据 Log Queries 中的日志计数计算值。通过这种查询语句，可以计算将日志数据量化成指标信息，并且，Promtail 可以通过这种查询语句将指标信息，填充到自己暴露的 Metrics 端点中。

注意：由于 Loki 的设计，所有 LogQL 查询都必须包含一个 Log Queries 中的 日志流选择器

日志流选择器确定将搜索多少日志流（日志内容的唯一来源，例如文件）。然后，更细粒度的日志流选择器将搜索到的流的数量减少到可管理的数量。这意味着传递给日志流选择器的标签将影响查询执行的相对性能。然后使用过滤器表达式对来自匹配日志流的聚合日志进行分布式 grep。

# Log Queries(日志查询)

详见 [Log Queries](docs/6.可观测性/Logs/Loki/LogQL/Log%20Queries.md)

# Metric Queries(指标查询)

详见 [Metric Queries](docs/6.可观测性/Logs/Loki/LogQL/Metric%20Queries.md)
