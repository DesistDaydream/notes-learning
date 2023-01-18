---
title: OpenTelemetry
weight: 2
---

# 概述

> 参考：
> - [GitHub 组织，OpenTelemetry](https://github.com/open-telemetry)
> - [官网](https://opentelemetry.io/)
> - [官方文档](https://opentelemetry.io/docs/)
> - [公众号-OpenTelemetry，OpenTelemetry 核心原理篇 ：怎么理解分布式链路追踪技术？](https://mp.weixin.qq.com/s/bcziZg8RhCrMGYgFeN76cw)
> - [公众号-OpenTelemetry，在生产环境如何选择靠谱的 APM 系统](https://mp.weixin.qq.com/s/3dD0hIuqpXdepLVC6V7aoA)

**OpenTelemetry(开放式遥测技术，简称 OTel)** 于 2019 年 5 月由 [OpenTracing](https://opentracing.io/) 与 OpenCensus 合并而成([Google Open Source](https://opensource.googleblog.com/2019/05/opentelemetry-merger-of-opencensus-and.html))，是一组 API、SDK、工具、更是一种遥测标准，旨在创建和管理 **Telemetry Data(遥测数据)。**通过 OpenTelemetry 标准创建的程序，可以采集 OpenTelemetry 标准的遥测数据，并发送到我们指定的后端中。OpenTelemetry 支持各种流行的开源后端项目，比如 Prometheus、Jaeger 等。

> 遥测数据包括：Traces(链路追踪数据)、Metrics(指标数据)、logs(日志数据)

注意：OpenTelemetry 不是像 Prometheus、Jaeger 那样的可观察性后端。相反，OpenTelemetry 支持将数据导出到各种开源和商业的后端产品中，它提供了一个可插拔的架构，因此可以轻松添加其他技术协议和格式。

OTEL 之于可观测性，类似 OCI 之于容器。

## OpenTelemetry 组件

目前，OpenTelemetry 由以下几个主要组件组成：

- **规范** # 与编程语言无关的规范，规定了遥测数据格式等
- **工具** # 用于采集、转换、导出遥测数据的工具
- **SDK** # 用于为各种编程语言提供编写符合 OpenTelemetry 规范的工具
- **自动 instrumentation 和 贡献包** # 没搞懂这是什么？

## OpenTelemetry 历史

# OpenTelemetry 实现

[GitHub 项目，grafana/agent](https://github.com/grafana/agent) #&#x20;
[GitHub 项目，flashcatcloud/categraf](https://github.com/flashcatcloud/categraf) # 通过配置文件，采集所有数据，然后 Push 给 Prom(Prom 需要使用 `--web.enable-remote-write-receiver` 为自身开启远程写功能)，暂时没有等待 pull 的功能(截止 2022.6.1 v0.1.0 版本)
<https://www.guance.com/> 观测云。。。。这个产品。。怎么说呢。。上来就让人各种注册才能体验的感觉很不好。。而且在云原生社区可观测性 SIG 群里，这家人的表达方式和处理事情的态度给人的感觉也不好~~工作内部矛盾放在群里说。。还揭露个人隐私。。。。o(╯□╰)o~

## Grafana Agent

> 参考：
> - [GitHub 项目，grafana/agent](https://github.com/grafana/agent)
> - [官方文档](https://grafana.com/docs/agent/latest/)

Grafana Agent 收集遥测数据并将其转发到 Grafana Stack、Grafana Cloud 或 Grafana Enterprise 的开源部署，然后可以在其中分析您的数据。您可以在 Kubernetes 和 Docker 上安装 Grafana Agent，或者作为 Linux、macOS 和 Windows 机器的系统进程。

Grafana Agent 是开源的，其源代码可在 GitHub 上的<https://github.com/grafana/agent>上获得。

Grafana Agent 适用于希望收集和转发遥测数据以进行分析和待命警报的工程师、操作员或管理员。那些运行 Grafana Agent 的人必须安装和配置 Grafana Agent 才能正确收集遥测数据并监控正在运行的代理的健康状况。
