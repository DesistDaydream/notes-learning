---
title: "Grafana API"
linkTitle: "Grafana API"
weight: 20
---

# 概述

> 参考：
>
> - [官网文档，开发资源 - HTTP API](https://grafana.com/docs/grafana/latest/developer-resources/api-reference/http-api/)

# 管理 API

https://grafana.com/docs/grafana/latest/developer-resources/api-reference/http-api/admin

热加载 [Provisioning](docs/6.可观测性/Grafana/Grafana%20Configuration/Provisioning.md) 配置

POST /api/admin/provisioning/dashboards/reload

POST /api/admin/provisioning/datasources/reload

POST /api/admin/provisioning/plugins/reload

POST /api/admin/provisioning/access-control/reload

POST /api/admin/provisioning/alerting/reload