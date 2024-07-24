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

# debug

可以生成 [pprof](/docs/2.编程/高级编程语言/Go/Go%20工具/pprof/pprof.md) 的 Profile 文件