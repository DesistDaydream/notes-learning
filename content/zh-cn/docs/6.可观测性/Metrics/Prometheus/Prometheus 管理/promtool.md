---
title: promtool
linkTitle: promtool
date: 2024-07-24T10:48
weight: 20
---

# 概述

> 参考：
>
> - [官方文档，命令行工具 - promtool](https://prometheus.io/docs/prometheus/latest/command-line/promtool)
> - https://blog.51cto.com/u_13236892/5968043

promtool 是 Prometheus 一个命令行工具，用以管理、检查 Promethus，包括 规则配置、etc. 。

> [!Warning]
> 截至 2024-07-24 还不够完善，有些功能无法添加认证能力，对以添加认证的 Prom 无能为力。

# Syntax(语法)

| Command | Description                  |
| ------- | ---------------------------- |
| check   | 检查资源的有效性。比如配置文件是否正确、etc.     |
| query   | 运行 PromQL 获取查询结果             |
| debug   | 获取 Debug 信息                  |
| push    | Push to a Prometheus server. |
| test    | 单元测试                         |
| tsdb    | Run tsdb commands.           |
| promql  | PromQL 格式化与编辑器               |

> [!Notes]
> 截至 2024-08-01，--http.config.file 选项的格式在 https://github.com/prometheus/common/blob/v0.55.0/config/http_config.go#L299 ，与 prometheus 的 --web.config.file 配置格式并不一致

# query

- **instant** # 即时向量查询
- **range** # 范围向量查询

## range

范围向量查询。默认返回 302 个样本，当前时间为结束时间，当前时间的前 5 分钟是开始时间，每秒 1 个样本。

> Notes: 若采集周期超过 1 秒，那不在采集时间点的数据，用上一个采集时间点的数据补充，保持一致。

OPTIONS

- **--start**(RFC3339 or Unix-time) # 范围查询的开始时间
- **--end**(RFC3339 or Unix-time) # 范围查询的结束时间
- **--step**(DURATION) # 查询步长（持续时间）。i.e. 每隔 step 时间取一个样本数据。

## EXAMPLE

范围向量查询，查询从 2024-08-01T08:01:01Z 开始到现在的所有数据，每隔 3 分钟取一个样本

```bash
promtool query \
  --http.config.file=http.conf \
  range http://localhost:9090 \
  'hdf_jmr_24_hour_security_log_files{security_data_code="3002"}' \
  --start=2024-08-01T08:01:01Z \
  --step=3m
```

# debug

可以生成 [pprof](/docs/2.编程/高级编程语言/Go/Go%20工具/pprof/pprof.md) 的 Profile 文件

