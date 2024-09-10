---
title: Label matchers
linkTitle: Label matchers
date: 2024-09-02T12:41
weight: 20
---

# 概述

> 参考：
>
> - [官方文档，告警 - 配置，标签匹配器](https://prometheus.io/docs/alerting/latest/configuration/#label-matchers)

**Label matchers(标签匹配器)** 类似 K8S 中的 [Label and Selector(标签和选择器)](/docs/10.云原生/Kubernetes/API%20Resource%20与%20Object/Object%20管理/Label%20and%20Selector(标签和选择器)/Label%20and%20Selector(标签和选择器).md)，是 Prometheus 中常见用以查找数据的功能，与 [PromQL](/docs/6.可观测性/Metrics/Prometheus/PromQL/PromQL.md) 的 XX选择器也有类似的效果。

在 [Alertmanager 配置](/docs/6.可观测性/Metrics/Alertmanager/Alertmanager%20配置.md) 中可以通过 Label matchers 过滤出想要处理的告警条目；在 [Promethesu Server](/docs/6.可观测性/Metrics/Prometheus/Configuration/Promethesu%20Server.md) 的配置中也有类似（没有明确指出，比如 .remote_read.required_matchers）的逻辑，尽管有时候语法可能并不相同。