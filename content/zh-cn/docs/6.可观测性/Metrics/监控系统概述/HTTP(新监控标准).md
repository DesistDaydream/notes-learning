---
title: HTTP(新监控标准)
linkTitle: HTTP(新监控标准)
weight: 20
---

# 概述

> 参考：
>
> - date: "2024-12-06T15:26"

由于 SNMP 的种种不便，现在更多的是基于 [HTTP](/docs/4.数据通信/Protocol/HTTP/HTTP.md) 协议来实现监控指标的采集。

同样，也是需要一个 Client 采集指标，需要一个 Server 端接收指标后存储指标。

像 SNMP 协议一样，光有协议还不行，基于 HTTP 协议的监控也需要一个数据模型的标准，就像 MIB 和 OID 类似。而现阶段，从 [Prometheus](/docs/6.可观测性/Metrics/Prometheus/Prometheus.md) 的 [Data Model(数据模型)](/docs/6.可观测性/Metrics/Prometheus/Storage(存储)/Data%20Model(数据模型).md) 演化过来的 OpenMetrics 标准，就是这么一种东西。

# OpenMetrics

详见 [OpenMetrics](docs/6.可观测性/Metrics/监控系统概述/OpenMetrics.md)