---
title: HTTP(新监控标准)
---

# 概述

由于 SNMP 的种种不便，现在更多的是基于 [HTTP](/docs/4.数据通信/Protocol/HTTP/HTTP.md) 协议来实现监控指标的采集。

同样，也是需要一个 Client 采集指标，需要一个 Server 端接收指标后存储指标。

像 SNMP 协议一样，光有协议还不行，基于 HTTP 协议的监控也需要一个数据模型的标准，就像 MIB 和 OID 类似。而现阶段，从 [Prometheus](/docs/6.可观测性/Metrics/Prometheus/Prometheus.md) 的 [Data Model(数据模型)](/docs/6.可观测性/Metrics/Prometheus/Storage(存储)/Data%20Model(数据模型).md) 演化过来的 OpenMetrics 标准，就是这么一种东西。

# OpenMetrics

> 参考：
>
> - [GitHub 项目，OpenObservability/OpenMetrics](https://github.com/OpenObservability/OpenMetrics)
> - [官网](https://openmetrics.io/)
> - [OpenMetrics 规范](https://github.com/OpenObservability/OpenMetrics/blob/main/specification/OpenMetrics.md)

**OpenMetrics** 是新时代的监控指标的标准，由 CNCF 主导，OpenMetrics 定义了大规模传输云原生指标的事实标准。

- **OpenMetricsSpec** 用来定义监控指标的标准

> [!Attention]
> [公众号 - InfoQ，OpenMetrics 归档并合并到 Prometheus](https://mp.weixin.qq.com/s/Wvh8AskHtOe2WoFPyAfVjA)
>
> - 英文帖子: https://horovits.medium.com/openmetrics-is-archived-merged-into-prometheus-d555598d2d04
>
> [GitHub 项目，cncf/toc Issue 1364](https://github.com/cncf/toc/issues/1364) 已在 2024 年 8 月份将 OpenMetrics 项目归档合并到 Prometheus 中。

# Data Model(数据模型)

https://github.com/OpenObservability/OpenMetrics/blob/main/specification/OpenMetrics.md#data-model

平时我们口语交流，一般将随时间变化的数据称为 **Metrics(指标)**。这是监控数据的另一种叫法，与 OID 类似，可以代表一种监控数据、同时也是一种名词，比如我采集某个程序的监控数据，可以说采集这个程序的 Metrics。所以 Metrics 是一个抽象的叫法。

详见 Prometheus [Data Model(数据模型)](/docs/6.可观测性/Metrics/Prometheus/Storage(存储)/Data%20Model(数据模型).md)
