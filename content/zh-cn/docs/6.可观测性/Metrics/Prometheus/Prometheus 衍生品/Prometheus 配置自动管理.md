---
title: Prometheus 配置自动管理
created: 2026-07-01T17:57
weight: 100
---

# 概述

> 参考：
>
> - 

[jimmidyson/configmap-reload](https://github.com/jimmidyson/configmap-reload)

[prometheus-operator 的 prometheus-config-reloader](https://github.com/prometheus-operator/prometheus-operator/tree/main/cmd/prometheus-config-reloader)

# Prometheus Configmanager

> 参考：
>
> - [GitHub 项目，facebookincubator/prometheus-configmanager](https://github.com/facebookincubator/prometheus-configmanager)

Prometheus-Configmanager 将规则配置文件抽象成 **tenant(租户)**，每个 tenant 都会有对应唯一的文件名(格式为：**TenantID_rules.yml**)。每当修改这些规则配置文件时，都会对 http://PrometheusIP:PORT/-/reload 发送 POST 请求以便让 Prometheus Server 重新加载配置文件。

> 一般情况，一个文件只有一个规则组，组名与 tenant 名称保持一致

## 基本示例

```json
curl -X POST "http://localhost:9100/v1/desistdaydream/alert" -H  "accept: application/json" -H  "Content-Type: application/json" -d '
{
  "alert": "test",
  "expr": "string",
  "labels": {
    "additionalProp1": "string",
    "additionalProp2": "string"
  },
  "for": "1m",
  "annotations": {
    "additionalProp1": "string",
    "additionalProp2": "string"
  }
}'
```

上面的 curl 请求会生成如下文件

```yaml
~]# cat desistdaydream_rules.yml
groups:
    - name: desistdaydream
      rules:
        - alert: test
          expr: string
          for: 1m
          labels:
            additionalProp1: string
            additionalProp2: string
            tenant: desistdaydream
          annotations:
            additionalProp1: string
            additionalProp2: string
```

## 部署

### 构建

```bash
docker build -t prometheus-configurer -f  Dockerfile.prometheus .
docker tag prometheus-configurer:latest lchdzh/prometheus-configurer:latest
```

### 运行

```bash
docker run -d --name prometheus-configmanager \
  --network=host \
  -v /opt/monitoring/config/prometheus/rules/:/etc/prometheus/rules/ \
  lchdzh/prometheus-configurer:latest \
  -prometheusURL='172.19.42.248:9090' \
  -rules-dir=/etc/prometheus/rules
```

## API 详解

### GET /v1/{tenant_id}/alert # 获取 {tennat_id}\_rules.yml 文件中所有告警规则

```bash
curl -X GET "http://localhost:9100/v1/desistdaydream/alert" -H  "accept: application/json" -H  "Content-Type: application/json"
```

### POST /v1/{tenant_id}/alert # 创建 {tennat_id}\_rules.yml 文件中的告警规则

若文件不存在则自动创建，若规则不存在则创建，只可一条一条规则创建

### GET /v1/{tenant_id}/alert/{alert_name} # 获取 {tennat_id}\_rules.yml 文件中指定的告警规则

### DELETE /v1/{tenant_id}/alert/{alert_name} # 删除 {tennat_id}\_rules.yml 文件中指定的告警规则

```bash
curl -X DELETE "http://localhost:9100/v1/desistdaydream/alert/test" -H  "accept: application/json" -H  "Content-Type: application/json"
```

### PUT /v1/{tenant_id}/alert/{alert_name} # 更新 {tennat_id}\_rules.yml 文件中已经存在的指定的告警规则

```bash
curl -X PUT "http://localhost:9100/v1/desistdaydream/alert/test" -H  "accept: application/json" -H  "Content-Type: application/json" -d '
{
  "alert": "test",
  "expr": "string_update",
  "labels": {
    "additionalProp1": "string"
  },
  "for": "1m",
  "annotations": {
    "additionalProp1": "string"
  }
}'
```

### POST /v1/{tenant_id}/alert/bulk # 批量更新/创建 {tennat_id}\_rules.yml 文件的警报规则

若文件不存在，则

```bash
~]# curl -X POST "http://localhost:9100/v1/desistdaydream/alert/bulk" -H  "accept: application/json" -H  "Content-Type: application/json" -d '
[{
  "alert": "test1",
  "expr": "string",
  "labels": {
    "additionalProp1": "string"
  },
  "for": "1m",
  "annotations": {
    "additionalProp1": "string"
  }
},
{
  "alert": "test2",
  "expr": "string",
  "labels": {
    "additionalProp1": "string"
  },
  "for": "1m",
  "annotations": {
    "additionalProp1": "string"
  }
}
]'
```

### GET /v1/tenancy # 查看所有 tennat 的状态

# OpenAPI 3.0

