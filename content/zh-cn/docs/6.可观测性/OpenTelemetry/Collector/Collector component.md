---
title: Collector component
linkTitle: Collector component
weight: 20
---

# 概述

> 参考：
>
> - [GitHub 项目，open-telemetry/opentelemetry-collector-contrib](https://github.com/open-telemetry/opentelemetry-collector-contrib)

**Receivers** # https://github.com/open-telemetry/opentelemetry-collector/blob/main/receiver

**Processors** # https://github.com/open-telemetry/opentelemetry-collector/tree/main/processor

**Exporters** # https://github.com/open-telemetry/opentelemetry-collector/tree/main/exporter

**Connector** # https://github.com/open-telemetry/opentelemetry-collector/blob/main/connector

**Extension** # https://github.com/open-telemetry/opentelemetry-collector/blob/main/extension

# Receiver

## Filelog

https://github.com/open-telemetry/opentelemetry-collector-contrib/tree/main/receiver/filelogreceiver

从文件读取日志

# Processor

## Attributes

https://github.com/open-telemetry/opentelemetry-collector-contrib/tree/main/processor/attributesprocessor

## Resource

https://github.com/open-telemetry/opentelemetry-collector-contrib/tree/main/processor/resourceprocessor

改变 Attributes(属性)

# Exporter

# Connector

# Extension

# 日志相关的功能模块

## stanza

> 参考：
> - https://github.com/open-telemetry/opentelemetry-collector-contrib/tree/main/pkg/stanza
> - Operator https://github.com/open-telemetry/opentelemetry-collector-contrib/blob/main/pkg/stanza/docs/operators/README.md

> Tips: stanza 最初由 observIQ 开发为一款独立的日志代理程序。作为独立代理，它具备读取、解析和导出日志的功能。该项目于 2021 年捐赠给 OpenTelemetry 项目，此后被改造为 OpenTelemetry Collector 中多个日志接收器的主要代码库。

该模块可以实现类似 Processor 的能力，但是这些能力是内嵌到各种日志相关的 Receiver 中的。

该模块的核心逻辑由 entry.Entry 和 Operators 组成

- **entry.Entry** # stanza 对 OpenTelemetry 的日志数据模型进行了独立表示，其中每个单独的日志记录都被建模为 entry.Entry
- **Operators** # 是日志处理的最基本单元。每个 Operator 都负责一项单一职责，例如从文件中读取行，或从字段中解析 JSON 数据。然后，将多个操作符串联起来，形成管道，以实现所需的结果。

> [!Note] stanza 模块与 Collector 的 Processor 的区别
>
> https://github.com/open-telemetry/opentelemetry-collector-contrib/tree/main/pkg/stanza#faq
>
> Q: Why don't we make every parser and transform operator into a distinct OpenTelemetry processor?
> 问：为什么我们不把每个 Opertaor 都做成一个独立的 OpenTelemetry 处理器呢？
>
> A: The nature of a receiver is that it reads data from somewhere, converts it into the OpenTelemetry data model, and emits it. Unlike most other receivers, those which read logs from traditional log media generally require extensive flexibility in the way that they convert an external format into the OpenTelemetry data model. Parser and transformer operators are designed specifically to provide sufficient flexibility to handle this conversion. Therefore, they are embedded directly into receivers in order to encapsulate the conversion logic and to prevent this concern from leaking further down the pipeline.
>
> 答：Receiver 的本质在于从某个地方读取数据，将其转换为 OpenTelemetry 数据模型，然后发送出去。与其它大多数 Receiver 不同，那些从传统日志介质读取日志的 Receiver 通常需要在将外部格式转换为 OpenTelemetry 数据模型方面具有极大的灵活性。Operators 正是为了提供足够的灵活性来处理这种转换而设计的。因此，它们被直接嵌入到 Receiver 中，以封装转换逻辑，并防止这种问题进一步扩散到后续流程中。

多个 Operator 组成 **Operator Sequences(运算序列)**（可以理解为一个 Pipeline，这个 Pipeline 是日志类型数据的多个处理阶段），定义了从 Receiver 发出日志之前应该如何解析和过滤日志。

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

