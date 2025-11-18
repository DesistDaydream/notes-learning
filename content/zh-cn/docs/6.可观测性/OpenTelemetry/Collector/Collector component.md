---
title: Collector component
linkTitle: Collector component
weight: 20
---

# 概述

> 参考：
>
> - [GitHub 项目，open-telemetry/opentelemetry-collector-contrib](https://github.com/open-telemetry/opentelemetry-collector-contrib)

**Receivers**

- https://github.com/open-telemetry/opentelemetry-collector/blob/main/receiver
- https://github.com/open-telemetry/opentelemetry-collector-contrib/tree/main/receiver

**Processors**

- https://github.com/open-telemetry/opentelemetry-collector/tree/main/processor
- https://github.com/open-telemetry/opentelemetry-collector-contrib/tree/main/processor

**Exporters**

- https://github.com/open-telemetry/opentelemetry-collector/tree/main/exporter
- https://github.com/open-telemetry/opentelemetry-collector-contrib/tree/main/exporter

**Connector**

- https://github.com/open-telemetry/opentelemetry-collector/blob/main/connector
- https://github.com/open-telemetry/opentelemetry-collector-contrib/tree/main/connector

**Extension**

- https://github.com/open-telemetry/opentelemetry-collector/blob/main/extension
- https://github.com/open-telemetry/opentelemetry-collector-contrib/tree/main/extension

# Receiver

## Filelog

https://github.com/open-telemetry/opentelemetry-collector-contrib/tree/main/receiver/filelogreceiver

从文件读取日志

# Processor

## Attributes

https://github.com/open-telemetry/opentelemetry-collector-contrib/tree/main/processor/attributesprocessor

修改 span、log、metric 属性

> Tips: 这个 Processor 命名感觉有点问题，Attributes Processor 并不是处理所有的，只是处理特定 Attributes 的，比如 Resource 的话，就需要用 Resource Attributes 处理。

## Resource

https://github.com/open-telemetry/opentelemetry-collector-contrib/tree/main/processor/resourceprocessor

修改 Resource attributes(资源属性)

# Exporter

## debug

https://github.com/open-telemetry/opentelemetry-collector/tree/main/exporter/debugexporter

# Connector

# Extension

# 日志相关的 Receiver 工作逻辑

[OTel Blog - 2024-05-22, 隆重推出适用于 OpenTelemetry Collector 的全新容器日志解析器](https://opentelemetry.io/blog/2024/otel-collector-container-log-parser/)

## Stanza

> 参考：
>
> - [GitHub 项目，open-telemetry/opentelemetry-collector-contrib - pkg/stanza](https://github.com/open-telemetry/opentelemetry-collector-contrib/tree/main/pkg/stanza)
> - [GitHub 项目，open-telemetry/opentelemetry-collector-contrib - pkg/stanza/docs/operators/README.md](https://github.com/open-telemetry/opentelemetry-collector-contrib/blob/main/pkg/stanza/docs/operators/README.md)

Stanza 是内嵌在各种日志相关的 Receiver 中的 [DataPipeline](docs/6.可观测性/DataPipeline/DataPipeline.md)。可以让诸如 Filelog、etc. 的日志相关的 Receiver 在进入 Processor 之前，进行 解析、过滤、处理、etc. 。

> [!Tip] Stanza 最初由 observIQ 开发为一款独立的日志代理程序（[GitHub 项目，observIQ/stanza](https://github.com/observIQ/stanza)），旨在替代 Fluentd、Fluent Bit 和 Logstash。作为独立代理，它具备读取、处理和导出日志的功能。
>
> 该项目于 2021 年捐赠给 OpenTelemetry 项目，此后被改造为 OpenTelemetry Collector 中多个日志相关的 Receiver 的主要底层代码。

> [!Note] stanza 模块与 Collector 的 Processor 的区别
>
> https://github.com/open-telemetry/opentelemetry-collector-contrib/tree/main/pkg/stanza#faq
>
> Q: Why don't we make every parser and transform operator into a distinct OpenTelemetry processor?
>
> 问：为什么我们不把每个 Opertaor 都做成一个独立的 OpenTelemetry 处理器呢？
>
> A: The nature of a receiver is that it reads data from somewhere, converts it into the OpenTelemetry data model, and emits it. Unlike most other receivers, those which read logs from traditional log media generally require extensive flexibility in the way that they convert an external format into the OpenTelemetry data model. Parser and transformer operators are designed specifically to provide sufficient flexibility to handle this conversion. Therefore, they are embedded directly into receivers in order to encapsulate the conversion logic and to prevent this concern from leaking further down the pipeline.
>
> 答：Receiver 的本质在于从某个地方读取数据，将其转换为 OpenTelemetry 数据模型，然后发送出去。与其它大多数 Receiver 不同，那些从传统日志介质读取日志的 Receiver 通常需要在将外部格式转换为 OpenTelemetry 数据模型方面具有极大的灵活性。Operators 正是为了提供足够的灵活性来处理这种转换而设计的。因此，它们被直接嵌入到 Receiver 中，以封装转换逻辑，并防止这种问题进一步扩散到后续流程中。

Stanza 的核心逻辑有两部分

- **Entry** # 抽象建模。Stanza 对 OpenTelemetry 的每个单独的日志都被建模为 Entry。
- **Operators** # 日志处理的最基本单元。每个 Operator 都负责一项单一职责，例如从文件中读取行，或从字段中解析 JSON 数据。然后，将多个操作符串联起来，形成管道，以实现所需的结果。

Entry 是日志数据在 pipeline 中流转时的基本表示形式。所有 Operator 都会创建、修改、使用 Entry。

> e.g. 一行日志 JSON 解析后的结果称为 Entry；从解析结果提取出来某些字段的值的结果也称为 Entry。解析 JSON 的行为称为 Operator；从解析结果提取字段的值的行为也称为 Operator

多个 Operator 组成 **Operator Sequences(运算序列)**（可以理解为 Pipeline），定义了日志从 Receiver 发出之前应该如何处理它们。

### Operators

[可用的 Operators 操作](https://github.com/open-telemetry/opentelemetry-collector-contrib/blob/main/pkg/stanza/docs/operators/README.md#what-operators-are-available)：

- Inputs:
    - [file_input](https://github.com/open-telemetry/opentelemetry-collector-contrib/blob/main/pkg/stanza/docs/operators/file_input.md)
    - [journald_input](https://github.com/open-telemetry/opentelemetry-collector-contrib/blob/main/pkg/stanza/docs/operators/journald_input.md)
    - [stdin](https://github.com/open-telemetry/opentelemetry-collector-contrib/blob/main/pkg/stanza/docs/operators/stdin.md)
    - [syslog_input](https://github.com/open-telemetry/opentelemetry-collector-contrib/blob/main/pkg/stanza/docs/operators/syslog_input.md)
    - [tcp_input](https://github.com/open-telemetry/opentelemetry-collector-contrib/blob/main/pkg/stanza/docs/operators/tcp_input.md)
    - [udp_input](https://github.com/open-telemetry/opentelemetry-collector-contrib/blob/main/pkg/stanza/docs/operators/udp_input.md)
    - [windows_eventlog_input](https://github.com/open-telemetry/opentelemetry-collector-contrib/blob/main/pkg/stanza/docs/operators/windows_eventlog_input.md)
- Parsers:
    - [csv_parser](https://github.com/open-telemetry/opentelemetry-collector-contrib/blob/main/pkg/stanza/docs/operators/csv_parser.md)
    - [json_parser](https://github.com/open-telemetry/opentelemetry-collector-contrib/blob/main/pkg/stanza/docs/operators/json_parser.md)
    - [json_array_parser](https://github.com/open-telemetry/opentelemetry-collector-contrib/blob/main/pkg/stanza/docs/operators/json_array_parser.md)
    - [regex_parser](https://github.com/open-telemetry/opentelemetry-collector-contrib/blob/main/pkg/stanza/docs/operators/regex_parser.md)
    - [scope_name_parser](https://github.com/open-telemetry/opentelemetry-collector-contrib/blob/main/pkg/stanza/docs/operators/scope_name_parser.md)
    - [syslog_parser](https://github.com/open-telemetry/opentelemetry-collector-contrib/blob/main/pkg/stanza/docs/operators/syslog_parser.md)
    - [severity_parser](https://github.com/open-telemetry/opentelemetry-collector-contrib/blob/main/pkg/stanza/docs/operators/severity_parser.md)
    - [time_parser](https://github.com/open-telemetry/opentelemetry-collector-contrib/blob/main/pkg/stanza/docs/operators/time_parser.md)
    - [trace_parser](https://github.com/open-telemetry/opentelemetry-collector-contrib/blob/main/pkg/stanza/docs/operators/trace_parser.md)
    - [uri_parser](https://github.com/open-telemetry/opentelemetry-collector-contrib/blob/main/pkg/stanza/docs/operators/uri_parser.md)
    - [key_value_parser](https://github.com/open-telemetry/opentelemetry-collector-contrib/blob/main/pkg/stanza/docs/operators/key_value_parser.md)
    - [container](https://github.com/open-telemetry/opentelemetry-collector-contrib/blob/main/pkg/stanza/docs/operators/container.md)
- Outputs:
    - [file_output](https://github.com/open-telemetry/opentelemetry-collector-contrib/blob/main/pkg/stanza/docs/operators/file_output.md)
    - [stdout](https://github.com/open-telemetry/opentelemetry-collector-contrib/blob/main/pkg/stanza/docs/operators/stdout.md)
- General purpose:
    - [add](https://github.com/open-telemetry/opentelemetry-collector-contrib/blob/main/pkg/stanza/docs/operators/add.md)
    - [copy](https://github.com/open-telemetry/opentelemetry-collector-contrib/blob/main/pkg/stanza/docs/operators/copy.md)
    - [filter](https://github.com/open-telemetry/opentelemetry-collector-contrib/blob/main/pkg/stanza/docs/operators/filter.md)
    - [flatten](https://github.com/open-telemetry/opentelemetry-collector-contrib/blob/main/pkg/stanza/docs/operators/flatten.md)
    - [move](https://github.com/open-telemetry/opentelemetry-collector-contrib/blob/main/pkg/stanza/docs/operators/move.md)
    - [noop](https://github.com/open-telemetry/opentelemetry-collector-contrib/blob/main/pkg/stanza/docs/operators/noop.md)
    - [recombine](https://github.com/open-telemetry/opentelemetry-collector-contrib/blob/main/pkg/stanza/docs/operators/recombine.md)
    - [regex_replace](https://github.com/open-telemetry/opentelemetry-collector-contrib/blob/main/pkg/stanza/docs/operators/regex_replace.md)
    - [remove](https://github.com/open-telemetry/opentelemetry-collector-contrib/blob/main/pkg/stanza/docs/operators/remove.md)
    - [retain](https://github.com/open-telemetry/opentelemetry-collector-contrib/blob/main/pkg/stanza/docs/operators/retain.md)
    - [router](https://github.com/open-telemetry/opentelemetry-collector-contrib/blob/main/pkg/stanza/docs/operators/router.md)
    - [sanitize_utf8](https://github.com/open-telemetry/opentelemetry-collector-contrib/blob/main/pkg/stanza/docs/operators/sanitize_utf8.md)
    - [unquote](https://github.com/open-telemetry/opentelemetry-collector-contrib/blob/main/pkg/stanza/docs/operators/unquote.md)
    - [assign_keys](https://github.com/open-telemetry/opentelemetry-collector-contrib/blob/main/pkg/stanza/docs/operators/assign_keys.md)

