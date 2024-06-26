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

Prometheus 可以从如下几类 Intrumenting 中 Scrape(抓取) 指标：

- **Instrumentation(检测仪/仪表化)** # 内部仪表。本质上是 Prometheus 的 **Client Libraries(客户端库)** 添加到程序代码中，以此暴露一个 endpoint，Prometheus Server 可以通过该 Endpoiint 抓取到指标。
   - 可以理解为内嵌的 Exporter，比如 Prometheus Server 的 9090 端口的 `/metrics` 就属于此类。
   - 说白了，就是被监控目标自己就可以吐出符合 Prometheus 格式的指标数据
- **Exporters** # 外部仪表。
  - 概念更为宽泛，除了使用到 Instrumentation 实现的各种程序外，还有一些通过脚本产生的符合 Prometheus [Data Model(数据模型)](/docs/6.可观测性/监控系统/Prometheus/Storage(存储)/Data%20Model(数据模型).md) 的纯文本的程序也可以称为 Exporter。
- **Pushgateway** # 针对需要推送指标的应用

# Instrumentation

> 参考:
>
> - [官方文档，Instrumentation - 客户端库](https://prometheus.io/docs/instrumenting/clientlibs/)

**Instrumentation(仪表化)**，顾名思义，将某个东西变为 Instrumenting（也可以说变为 Exporter）。所以 Instrumentation 是一组 Library，当我们在编写的程序代码中引入了 Instrumentation，并使用其提供的各种方法、接口，那么我们的程序就可以变成像仪表一样的东西，以展示出想要的 观测或监控 数据。

Prometheus 官方维护了某些语言的 Library:

- [Go](https://github.com/prometheus/client_golang)
- [Java or Scala](https://github.com/prometheus/client_java)
- [Python](https://github.com/prometheus/client_python)
- [Ruby](https://github.com/prometheus/client_ruby)
- [Rust](https://github.com/prometheus/client_rust)

除了官方自己维护的，还有一些官方推荐的第三方维护的特定语言的 Library，详见官网页。

## Prometheus 开发工具包

> 参考:
>
> [GitHub 项目，prometheus/exporter-toolkit](https://github.com/prometheus/exporter-toolkit)
> - https://pkg.go.dev/github.com/prometheus/exporter-toolkit/web

Prometheus Exporter Toolkit 为开发 Prometheus Exporter 提供能力的工具包。

借助 exporter-toolkit 可以为自己开发的 Exporter 加入像 [HTTPS 和 Authentication(认证)](/docs/6.可观测性/监控系统/Prometheus/HTTPS%20和%20Authentication(认证).md) 的能力。比如 node-exporter、push gateway、etc. 甚至 Prometheus 自身，都可以通过 --web.config.file 参数实现认证功能，这就是通过 exporter-toolkit 库实现的。

# Exporter

在[这里](https://github.com/prometheus/prometheus/wiki/Default-port-allocations)可以看到经过 Prometheus 官方注册的各类 Exporter 所默认使用的端口号。

# Push Gateway

[Push Gateway](/docs/6.可观测性/监控系统/Instrumenting/Push%20Gateway.md)
