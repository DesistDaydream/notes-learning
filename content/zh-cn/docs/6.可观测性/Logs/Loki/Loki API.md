---
title: Loki API
linkTitle: Loki API
date: 2024-09-26T08:37
weight: 20
---

# 概述

> 参考：
>
> - [官方文档, HTTP API](https://grafana.com/docs/loki/latest/api/)

每个组件都会暴露一些基本的 API

- [GET /ready](https://grafana.com/docs/loki/latest/api/#identify-ready-loki-instance)
- [GET /metrics](https://grafana.com/docs/loki/latest/api/#return-exposed-prometheus-metrics)
- [GET /config](https://grafana.com/docs/loki/latest/api/#list-current-configuration)
- [GET /services](https://grafana.com/docs/loki/latest/api/#list-running-services)
- [GET /loki/api/v1/status/buildinfo](https://grafana.com/docs/loki/latest/api/#list-build-information)

除了这几个基本的 API 以外，每个组件都会暴露一些专用的 API，若是以 Monolithic 架构启动 Loki，则下面的所有 API 都会在这个进程暴露。

# Querier API 与 Query Frontend API

查询器与查询前端暴露的 API 是我们最常用的 API，用来处理客户端发来的 LogQL。

- [GET /loki/api/v1/query](https://grafana.com/docs/loki/latest/api/#query-loki)
- [GET /loki/api/v1/query_range](https://grafana.com/docs/loki/latest/api/#query-loki-over-a-range-of-time)
  - [Step vs Interval](https://grafana.com/docs/loki/latest/api/#step-vs-interval)
- [GET /loki/api/v1/labels](https://grafana.com/docs/loki/latest/api/#list-labels-within-a-range-of-time)
- [GET /loki/api/v1/label/\<name>/values](https://grafana.com/docs/loki/latest/api/#list-label-values-within-a-range-of-time)
- [GET /loki/api/v1/series](https://grafana.com/docs/loki/latest/api/#list-series)
- [GET /loki/api/v1/index/stats](https://grafana.com/docs/loki/latest/api/#index-stats)
- [GET /loki/api/v1/tail](https://grafana.com/docs/loki/latest/api/#stream-log-messages)

# Distributor API

- [POST /flush](https://grafana.com/docs/loki/latest/api/#post-flush)
- [POST /ingester/flush_shutdown](https://grafana.com/docs/loki/latest/api/#post-ingesterflush_shutdown)

# Ingester API

- [POST /flush](https://grafana.com/docs/loki/latest/api/#flush-in-memory-chunks-to-backing-store)
- [POST /ingester/shutdown](https://grafana.com/docs/loki/latest/api/#flush-in-memory-chunks-and-shut-down)

# Ruler API

以 /loki/ 开头的 API 与 [Prometheus API](https://prometheus.io/docs/prometheus/latest/querying/api/) 兼容，结果格式可以互换使用

- [GET /ruler/ring](https://grafana.com/docs/loki/latest/api/#ruler-ring-status)
- [GET /loki/api/v1/rules](https://grafana.com/docs/loki/latest/api/#list-rule-groups)
- [GET /loki/api/v1/rules/{namespace}](https://grafana.com/docs/loki/latest/api/#get-rule-groups-by-namespace)
- [GET /loki/api/v1/rules/{namespace}/{groupName}](https://grafana.com/docs/loki/latest/api/#get-rule-group)
- [POST /loki/api/v1/rules/{namespace}](https://grafana.com/docs/loki/latest/api/#set-rule-group)
- [DELETE /loki/api/v1/rules/{namespace}/{groupName}](https://grafana.com/docs/loki/latest/api/#delete-rule-group)
- [DELETE /loki/api/v1/rules/{namespace}](https://grafana.com/docs/loki/latest/api/#delete-namespace)
- [GET /api/prom/rules](https://grafana.com/docs/loki/latest/api/#list-rule-groups)
- [GET /api/prom/rules/{namespace}](https://grafana.com/docs/loki/latest/api/#get-rule-groups-by-namespace)
- [GET /api/prom/rules/{namespace}/{groupName}](https://grafana.com/docs/loki/latest/api/#get-rule-group)
- [POST /api/prom/rules/{namespace}](https://grafana.com/docs/loki/latest/api/#set-rule-group)
- [DELETE /api/prom/rules/{namespace}/{groupName}](https://grafana.com/docs/loki/latest/api/#delete-rule-group)
- [DELETE /api/prom/rules/{namespace}](https://grafana.com/docs/loki/latest/api/#delete-namespace)
- [GET /prometheus/api/v1/rules](https://grafana.com/docs/loki/latest/api/#list-rules)
- [GET /prometheus/api/v1/alerts](https://grafana.com/docs/loki/latest/api/#list-alerts)

# Compactor API

- [GET /compactor/ring](https://grafana.com/docs/loki/latest/api/#get-compactorring)
- [POST /loki/api/v1/delete](https://grafana.com/docs/loki/latest/api/#post-lokiapiv1delete)
- [GET /loki/api/v1/delete](https://grafana.com/docs/loki/latest/api/#get-lokiapiv1delete)
- [DELETE /loki/api/v1/delete](https://grafana.com/docs/loki/latest/api/#delete-lokiapiv1delete)

# Series API

- GET /loki/api/v1/series
- POST /loki/api/v1/series
- GET /api/prom/series
- POST /api/prom/series
