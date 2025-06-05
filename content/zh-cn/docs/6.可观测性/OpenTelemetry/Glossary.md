---
title: Glossary
linkTitle: Glossary
weight: 2
---

# 概述

> 参考：
>
> - [官方文档，概念 - 术语表](https://opentelemetry.io/docs/concepts/glossary)


## Attribute

Attribute(属性) 是 OpenTelemetry 中用于表示 Metadata(元数据) 的概念。

> [!Note] Metadata(元数据) 在 OpenTelemetry 中就是 key/value pair(键值对) 的抽象描述。**本质上，Attribute 就是一堆 key/value 的集合**
>
> e.g. [Prometheus](/docs/6.可观测性/Metrics/Prometheus/Prometheus.md) 系列的 Metrics 的 **Labels**、某种描述场景下的 **Fields**、**Dimensions(维度)**、etc. 都可以称为 Attribute

Attribute 作为遥测数据的键值对信息，可以跨 Signals、Resources、etc. 使用。

[Resources](https://opentelemetry.io/docs/specs/otel/resource/sdk/), [Instrumentation Scopes](https://opentelemetry.io/docs/specs/otel/glossary/#instrumentation-scope), [Metric points](https://opentelemetry.io/docs/specs/otel/metrics/data-model/#metric-points), [Spans](https://opentelemetry.io/docs/specs/otel/trace/api/#set-attributes), Span [Events](https://opentelemetry.io/docs/specs/otel/trace/api/#add-events), Span [Links](https://opentelemetry.io/docs/specs/otel/trace/api/#link) and [Log Records](https://opentelemetry.io/docs/specs/otel/logs/data-model/) 可能包含一组属性。每个此类集合中的键都是唯一的，即，不能存在多个具有相同键的键值对。可以通过多种方式来强制执行唯一性，以使其最适合特定实现的限制。

详见 [Attribute 规范](https://opentelemetry.io/docs/specs/otel/common/#attributes)

## Resource

Resource 是 Attribute 的一种，是产生遥测数据的实体的 Attribute。e.g. [Kubernetes](/docs/10.云原生/Kubernetes/Kubernetes.md) 中容器运行的进程产生的遥测数据，会包含 进程名称、Pod 名称、命名空间、etc. 这些都属于 **Resource Attributes(资源属性)**

## Signal

Metrics、Logs、Traces、Baggage、etc.  都是 Signals 之一
