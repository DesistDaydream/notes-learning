---
title: Prometheus API
linkTitle: Prometheus API
weight: 1
---

# 概述

> 参考：
>
> - [官方文档，Prometheus - 管理 API](https://prometheus.io/docs/prometheus/latest/management_api/)
> - [官方文档，Prometheus - 查询 - HTTP API](https://prometheus.io/docs/prometheus/latest/querying/api/)

Prometheus 提供多种类型的 API 以满足不同需求。但是唯独没有可以修改配置的 API，Prometheus 的各种配置，只能通过重新加载修改后的配置文件这种方式来修改

Prometheus API 分两大块

- **HTTP API** # HTTP 接口，用于查询数据、状态等。所以也称为 Querying API
- **Management API** # 管理接口，用于简单管理 Prometheus Server，重载配置，健康检查等
