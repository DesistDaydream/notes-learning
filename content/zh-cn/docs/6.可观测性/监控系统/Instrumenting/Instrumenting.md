---
title: Instrumenting
linkTitle: Instrumenting
date: 2024-03-21T14:24
weight: 1
---

# 概述

> 参考：
> 
> - [官方文档，Instrumenting-Exporter](https://prometheus.io/docs/instrumenting/exporters/)
> - [官方文档，最佳实践-Instrumentation](https://prometheus.io/docs/practices/instrumentation/)

Prometheus 可以通过 3 种方式从目标上 Scrape(抓取) 指标：

1. **Instrumentation(检测仪)** # 内部代码。将 Prometheus 的 Client Libraries(客户端库) 添加到程序代码中，以此暴露一个 endpoint，Prometheus Server 可以通过该 Endpoiint 抓取到指标。
   1. 可以理解为内嵌的 Exporter，比如 Prometheus Server 的 9090 端口的 `/metrics` 就属于此类。
   2. 说白了，就是被监控目标自己就可以吐出符合 Prometheus 格式的指标数据
2. **Exporters** # 外部程序。
3. **Pushgateway** # 针对需要推送指标的应用

# Instrumentation

# Exporter

在[这里](https://github.com/prometheus/prometheus/wiki/Default-port-allocations)可以看到经过 Prometheus 官方注册的各类 Exporter 所默认使用的端口号。

# Push Gateway
