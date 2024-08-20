---
title: Template
linkTitle: Template
date: 2024-08-20T09:56
weight: 3
---

# 概述

> 参考：
>
> - [官方文档，配置 - 模板示例](https://prometheus.io/docs/prometheus/latest/configuration/template_examples/)
> - [官方文档，配置 - 模板参考](https://prometheus.io/docs/prometheus/latest/configuration/template_reference/)

Prometheus 可以在部分配置文件中（[Rules](/docs/6.可观测性/Metrics/Prometheus/Configuration/Rules.md)、etc.）使用 **Template(模板)** 的能力，Prometheus 的模板基于 [Go 语言的 Template](/docs/2.编程/高级编程语言/Go/Go%20规范与标准库/Template.md) 能力。

# Template Function

Prometheus 模板增加了一些函数以便更轻松得处理 PromQL 的查询结果（e.g. 将 Bytes 的数直接转为人类可读的带单位的结果），所有函数列表详见[官方文档](https://prometheus.io/docs/prometheus/latest/configuration/template_reference/#functions)

