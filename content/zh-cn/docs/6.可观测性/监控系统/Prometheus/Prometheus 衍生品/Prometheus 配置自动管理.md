---
title: Prometheus 配置自动管理
---

# 概述

[jimmidyson/configmap-reload](https://github.com/jimmidyson/configmap-reload)
[prometheus-operator 的 prometheus-config-reloader](https://github.com/prometheus-operator/prometheus-operator/tree/main/cmd/prometheus-config-reloader)

# Prometheus Configmanager

> 参考：
> - [GitHub 项目，facebookincubator/prometheus-configmanager](https://github.com/facebookincubator/prometheus-configmanager)

Prometheus-Configmanager 将规则配置文件抽象成 **tenant(租户)**，每个 tenant 都会有对应唯一的文件名(格式为：**TenantID_rules.yml**)。每当修改这些规则配置文件时，都会对 http://PrometheusIP:PORT/-/reload 发送 POST 请求以便让 Prometheus Server 重新加载配置文件。

> 一般情况，一个文件只有一个规则组，组名与 tenant 名称保持一致

## 基本示例

```json
curl -X POST "http://localhost:9100/v1/lichenhao/alert" -H  "accept: application/json" -H  "Content-Type: application/json" -d '
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
/ # cat lichenhao_rules.yml
groups:
    - name: lichenhao
      rules:
        - alert: test
          expr: string
          for: 1m
          labels:
            additionalProp1: string
            additionalProp2: string
            tenant: lichenhao
          annotations:
            additionalProp1: string
            additionalProp2: string
```

## 部署

### 构建

    docker build -t prometheus-configurer -f  Dockerfile.prometheus .
    docker tag prometheus-configurer:latest lchdzh/prometheus-configurer:latest

### 运行

    docker run -d --name prometheus-configmanager \
    --network=host \
    -v /opt/monitoring/config/prometheus/rules/:/etc/prometheus/rules/ \
    lchdzh/prometheus-configurer:latest \
    -prometheusURL='172.19.42.248:9090' \
    -rules-dir=/etc/prometheus/rules

## API 详解

### GET /v1/{tenant_id}/alert # 获取 {tennat_id}\_rules.yml 文件中所有告警规则

    curl -X GET "http://localhost:9100/v1/lichenhao/alert" -H  "accept: application/json" -H  "Content-Type: application/json"

### POST /v1/{tenant_id}/alert # 创建 {tennat_id}\_rules.yml 文件中的告警规则

若文件不存在则自动创建，若规则不存在则创建，只可一条一条规则创建

### GET /v1/{tenant_id}/alert/{alert_name} # 获取 {tennat_id}\_rules.yml 文件中指定的告警规则

### DELETE /v1/{tenant_id}/alert/{alert_name} # 删除 {tennat_id}\_rules.yml 文件中指定的告警规则

    curl -X DELETE "http://localhost:9100/v1/lichenhao/alert/test" -H  "accept: application/json" -H  "Content-Type: application/json"

### PUT /v1/{tenant_id}/alert/{alert_name} # 更新 {tennat_id}\_rules.yml 文件中已经存在的指定的告警规则

    curl -X PUT "http://localhost:9100/v1/lichenhao/alert/test" -H  "accept: application/json" -H  "Content-Type: application/json" -d '
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

### POST /v1/{tenant_id}/alert/bulk # 批量更新/创建 {tennat_id}\_rules.yml 文件的警报规则

若文件不存在，则

    curl -X POST "http://localhost:9100/v1/lichenhao/alert/bulk" -H  "accept: application/json" -H  "Content-Type: application/json" -d '
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

### GET /v1/tenancy # 查看所有 tennat 的状态

# OpenAPI 3.0

```yaml
openapi: 3.0.1
info:
title: prometheus-configmanager
description: ''
version: 1.0.0
tags:
- name: prometheus-configmanager
paths:
  /v1/tenancy:
    get:
summary: 检索 tenancy 配置
description: ''
      tags:
- prometheus-configmanager
parameters: []
      responses:
        '200':
description: Tenancy configuration
          content:
            application/json:
              schema:
$ref: '#/components/schemas/tenancy_config'
examples: {}
  '/v1/{tenant_id}/alert/{alert_name}':
    delete:
summary: '删除{tenant_id}的告警'
description: ''
      tags:
- prometheus-configmanager
      parameters:
- name: tenant_id
in: path
description: Tenant ID
required: true
example: test
          schema:
type: string
- name: alert_name
in: path
description: 要被删除的告警名称
required: true
example: test
          schema:
type: string
      responses:
        '204':
description: 删除成功
          content:
            '*/*':
              schema:
$ref: '#/components/schemas/UnexpectedError'
examples: {}
        '500':
description: 服务器错误
          content:
            application/json:
              schema:
type: object
                properties:
                  message:
type: string
                required:
- message
examples: {}
    put:
summary: '更新{tenant_id}中的{alert_name}告警规则'
description: ''
      tags:
- prometheus-configmanager
      parameters:
- name: tenant_id
in: path
description: Tenant ID
required: true
example: test
          schema:
type: string
- name: alert_name
in: path
description: Name of alert to be updated
required: true
example: test
          schema:
type: string
      requestBody:
        content:
          application/json:
            schema:
type: object
              properties:
                alert:
type: string
                expr:
type: string
                labels:
type: object
properties: {}
                for:
type: string
                annotations:
type: object
properties: {}
              required:
- alert
- expr
- labels
- for
- annotations
            example:
alert: test
expr: string_update
              labels:
additionalProp1: string
for: 1m
              annotations:
additionalProp1: string
      responses:
        '204':
description: Updated
          content:
            '*/*':
              schema:
type: object
properties: {}
examples: {}
        '400':
description: 规则不存在
          content:
            application/json:
              schema:
type: object
                properties:
                  message:
type: string
                required:
- message
examples: {}
    get:
summary: '获取{tenant_id}中的{alert_name}告警规则'
description: ''
      tags:
- prometheus-configmanager
      parameters:
- name: tenant_id
in: path
description: Tenant ID
required: true
example: test
          schema:
type: string
- name: alert_name
in: path
description: 仅仅是为了与另一个接口进行区分，填任何值都行，不写也行
required: true
example: test
          schema:
type: string
- name: alert_name
in: query
description: Alert Name
required: true
example: test
          schema:
type: string
      responses:
        '200':
description: Alert configuration
          content:
            application/json:
              schema:
type: array
                items:
type: object
                  properties:
                    alert:
type: string
                    expr:
type: string
                    for:
type: string
                    labels:
type: object
properties: {}
                    annotations:
type: object
properties: {}
examples: {}
        '500':
description: '500'
          content:
            application/json:
              schema:
$ref: '#/components/schemas/UnexpectedError'
examples: {}
  '/v1/{tenant_id}/alert':
    post:
summary: '创建{tennat_id}中的告警规则'
description: 若tennat_id不存在，则自动创建，若规则已存在返回400。
      tags:
- prometheus-configmanager
      parameters:
- name: tenant_id
in: path
description: Tenant ID
required: true
example: test
          schema:
type: string
      requestBody:
        content:
          application/json:
            schema:
type: object
              properties:
                alert:
type: string
                expr:
type: string
                labels:
type: object
properties: {}
                for:
type: string
description: '时间，可以是1m,1s等等'
                annotations:
type: object
properties: {}
              required:
- alert
- expr
- labels
- for
- annotations
            example:
alert: test
expr: string
              labels:
additionalProp1: string
additionalProp2: string
for: 1m
              annotations:
additionalProp1: string
additionalProp2: string
      responses:
        '200':
description: 成功
          content:
            '*/*':
              schema:
type: object
                properties:
                  message:
type: string
                required:
- message
examples: {}
        '201':
description: Created
          content:
            application/json:
              schema:
type: object
properties: {}
examples: {}
        '400':
description: 规则已存在
          content:
            application/json:
              schema:
type: object
                properties:
                  message:
type: string
                required:
- message
examples: {}
    get:
summary: '获取{tennat_id}的告警规则'
description: ''
      tags:
- prometheus-configmanager
      parameters:
- name: tenant_id
in: path
description: Tenant ID
required: true
example: test
          schema:
type: string
      responses:
        '200':
description: List of alert configurations
          content:
            application/json:
              schema:
type: array
                items:
type: object
                  properties:
                    alert:
type: string
                    expr:
type: string
                    for:
type: string
                    labels:
type: object
                      properties:
                        tenant:
type: string
                      required:
- tenant
                    annotations:
type: object
properties: {}
                  required:
- alert
- expr
- for
- labels
- annotations
examples: {}
  '/v1/{tenant_id}/alert/bulk':
    post:
summary: 批量更新或创建告警规则
description: ''
      tags:
- prometheus-configmanager
      parameters:
- name: tenant_id
in: path
description: Tenant ID
required: true
example: test
          schema:
type: string
      requestBody:
        content:
          application/json:
            schema:
type: array
              items:
type: object
                properties:
                  alert:
type: string
                  expr:
type: string
                  for:
type: string
                  labels:
type: object
                    properties:
                      tenant:
type: string
                  annotations:
type: object
properties: {}
                required:
- alert
- expr
- for
- labels
- annotations
            example:
- alert: test
expr: string
for: 1m
                labels:
additionalProp1: string
additionalProp2: string
tenant: lichenhao
                annotations:
additionalProp1: string
additionalProp2: string
- alert: test1
expr: string
for: 1m
                labels:
additionalProp1: string
tenant: lichenhao
                annotations:
additionalProp1: string
      responses:
        '200':
description: Success
          content:
            application/json:
              schema:
type: object
                properties:
                  Errors:
type: object
properties: {}
                  Statuses:
type: object
properties: {}
                required:
- Errors
- Statuses
examples: {}
        '500':
description: '500'
          content:
            application/json:
              schema:
$ref: '#/components/schemas/UnexpectedError'
examples: {}
components:
  schemas:
    UnexpectedError:
$ref: '#/components/schemas/error'
    alert_labels:
type: object
      additionalProperties:
type: string
    tenancy_config:
type: object
      properties:
        restrictor_label:
type: string
        restrict_queries:
type: boolean
    alert_bulk_upload_response:
type: object
      required:
- errors
- statuses
      properties:
        errors:
type: object
          additionalProperties:
type: string
        statuses:
type: object
          additionalProperties:
type: string
    error:
type: object
      required:
- message
      properties:
        message:
example: Error string
type: string
    alert_config:
type: object
      required:
- alert
- expr
- labels
- for
- annotations
      properties:
        alert:
type: string
        expr:
type: string
        labels:
$ref: '#/components/schemas/alert_labels'
        for:
type: string
        annotations:
$ref: '#/components/schemas/alert_labels'
    alert_config_list:
type: array
      items:
$ref: '#/components/schemas/alert_config'
```
