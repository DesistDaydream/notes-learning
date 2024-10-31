---
title: Collector
linkTitle: Collector
date: 2024-10-31T11:13
weight: 20
---

# 概述

> 参考：
>
> -

OpenTelemetry Collector 提供了一种与供应商无关的接收、处理和导出遥测数据的实现。它消除了运行、操作和维护多个代理/收集器的需要。这具有改进的可扩展性，并支持开源可观测性数据格式（例如 Jaeger、Prometheus、Fluent Bit 等）发送到一个或多个开源或商业后端。

Collector 会根据配置定时采集数据或被动接收数据以缓存，然后可以主动推送或被动等待拉取。Prometheus 可以配置 static_configs 从 OTel Collector 抓取其缓存的最新数据。