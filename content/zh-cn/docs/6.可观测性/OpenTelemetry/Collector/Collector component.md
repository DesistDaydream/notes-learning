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

### 配置文件中的可用字段

**key**(STRING) # 要修改的属性的 Key

**value**(STRING) # 要修改的属性的 Value

**action**(STRING) # 要如何修改属性

- 可用的值有：
    - **insert** # 插入属性。若属性已存在则无法插入
    - **update** # 更新属性。若属性不存在则无法更新
    - **upsert** # 更新属性。若属性不存在则插入
    - **delete** # 删除属性。
    - **hash** # 对现有属性进行哈希处理（SHA1）。
    - **extract** # 使用正则表达式，从输入键中提取值到规则中指定的目标键。如果目标键已存在，则会被覆盖。注意：它的行为类似于 Span Processor 的 `to_attributes` 以现有属性为源进行设置。
    - **convert** # 将属性转换为指定类型。

## Resource

https://github.com/open-telemetry/opentelemetry-collector-contrib/tree/main/processor/resourceprocessor

修改 Resource attributes(资源属性)

### 配置文件中的可用字段

与 Attributes 的一样

# Exporter

## debug

https://github.com/open-telemetry/opentelemetry-collector/tree/main/exporter/debugexporter

将遥测数据输出到控制台，以便进行调试。

假如我对文件 `/tmp/01/0001/20251119/9999/a.txt` 写入了一行内容: `DesistDaydream`。可以看到如下信息（这是最详细的 detailed 级别的信息）

```bash
2025-11-19T11:30:45.517+0800    info    Logs    {"resource": {}, "otelcol.component.id": "debug", "otelcol.component.kind": "exporter", "otelcol.signal": "logs", "resource logs": 1, "log records": 1}
2025-11-19T11:30:45.517+0800    info    ResourceLog #0
Resource SchemaURL:
Resource attributes:
     -> host_ip: Str(10.10.4.90)
ScopeLogs #0
ScopeLogs SchemaURL:
InstrumentationScope
LogRecord #0
ObservedTimestamp: 2025-11-19 03:30:45.517526302 +0000 UTC
Timestamp: 1970-01-01 00:00:00 +0000 UTC
SeverityText:
SeverityNumber: Unspecified(0)
Body: Str(DesistDaydream)
Attributes:
     -> log.file.name: Str(a.txt)
     -> log.file.path: Str(/tmp/01/0001/20251119/9999/a.txt)
Trace ID:
Span ID:
Flags: 0
        {"resource": {}, "otelcol.component.id": "debug", "otelcol.component.kind": "exporter", "otelcol.signal": "logs"}
```

Body 是从文件中读到的一行内容

Resource attributes 和 Attributes 来源于程序本身的逻辑以及如下配置文件内容：

```yaml
receivers:
  filelog/operators_demo:
    include:
      - /tmp/01/*/*/*/*.txt
    include_file_path: true

processors:
  resource/host:
    attributes:
      - key: host_ip
        value: "10.10.4.90"
        action: upsert

exporters:
  debug:
    verbosity: detailed

service:
  pipelines:
    logs/demo:
      receivers:
        - filelog/operators_demo
      processors:
        - resource/host
      exporters:
        - debug
```

### 配置文件中的可用字段

**verbosity**(STRING) # 输出内容的详细程度。可用的值: basic, normal,  detailed。`默认值: basic`

# Connector

# Extension

## FileStorage

https://github.com/open-telemetry/opentelemetry-collector-contrib/tree/main/extension/storage/filestorage

**File Storage(文件存储)** 扩展可以将状态持久化到本地文件系统。该扩展程序需要对某个目录拥有读写权限。可以使用默认目录，但该目录必须已存在，扩展程序才能正常运行。

### 配置文件中的可用字段

**directory**

- 默认目录在 Windows 系统上为 `%ProgramData%\Otelcol\FileStorage` ，在其他系统上为 `/var/lib/otelcol/file_storage` 。

# 日志相关的 Receiver 工作逻辑

[OTel Blog - 2024-05-22, 隆重推出适用于 OpenTelemetry Collector 的全新容器日志解析器](https://opentelemetry.io/blog/2024/otel-collector-container-log-parser/)

## Stanza

> 参考：
>
> - [GitHub 项目，open-telemetry/opentelemetry-collector-contrib - pkg/stanza](https://github.com/open-telemetry/opentelemetry-collector-contrib/tree/main/pkg/stanza)
> - [GitHub 项目，open-telemetry/opentelemetry-collector-contrib - pkg/stanza/docs/operators/README.md](https://github.com/open-telemetry/opentelemetry-collector-contrib/blob/main/pkg/stanza/docs/operators/README.md)

Stanza 是内嵌在各种日志相关的 Receiver 中的 [DataPipeline](/docs/6.可观测性/DataPipeline/DataPipeline.md)。可以让诸如 Filelog、etc. 的日志相关的 Receiver，在进入 Processor 之前，进行 解析、过滤、处理、etc. 。

> [!Tip] Stanza 最初由 observIQ 开发为一款独立的日志代理程序（[GitHub 项目，observIQ/stanza](https://github.com/observIQ/stanza)），旨在替代 Fluentd、Fluent Bit 和 Logstash。作为独立代理，它具备读取、处理和导出日志的功能。
>
> 该项目于 2021 年捐赠给 OpenTelemetry 项目，此后被改造为 OpenTelemetry Collector 中多个日志相关的 Receiver 的主要底层代码。

### 架构

- **Entry** # 抽象建模。Stanza 对 OpenTelemetry 的每个单独的日志都被建模为 Entry。
- **Operators** # 日志处理的最基本单元。每个 Operator 都表示一个行为，例如从文件中读取行，或从字段中解析 JSON 数据。然后，将多个操作符串联起来，形成管道，以实现所需的结果。
- **Operator Sequences(运算序列)** # 由多个 Operator 组成的日志处理 Pipeline，定义了日志从 Receiver 发出之前应该如何处理它们。

Entry 是日志数据在 Pipeline 中流转时的基本表示形式。所有 Operator 都会创建、修改、使用 Entry。

> e.g. 一行日志 JSON 解析后的结果称为 Entry；从解析结果提取出来某些字段的值的结果也称为 Entry。解析 JSON 的行为称为 Operator；从解析结果提取字段的值的行为也称为 Operator

> [!Note] Stanza 模块与 Processor 的区别
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

### Operators

可用的 Operators：

https://github.com/open-telemetry/opentelemetry-collector-contrib/blob/main/pkg/stanza/docs/operators/README.md#what-operators-are-available

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
    - [**regex_parser**](https://github.com/open-telemetry/opentelemetry-collector-contrib/blob/main/pkg/stanza/docs/operators/regex_parser.md) # 使用 Go 语言格式的[正则](/docs/2.编程/高级编程语言/Go/Go%20规范与标准库/文本处理/正则.md)表达式，解析 `parse_from` 字段给定的字符串。
        - 解析结果可以 保存到 Attributes 中
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

### 最佳实践

#### 动态得提取文件路径中某个目录名作为 Resource attribute

```yaml
receivers:
  filelog/test:
    include:
      - /tmp/01/*/*/*/*.txt
    # 将文件的路径添加为 log.file.path 属性
    include_file_path: true
    operators:
      # 正则匹配 log.file.path 属性的值（i.e. 文件路径），将 01/ 后面 / 之前的部分作为 type_code 属性的值
      - type: regex_parser
        regex: '^/tmp/01/(?P<type_code>\d+)/'
        parse_from: attributes["log.file.path"]
      # 把 type_code 属性拷贝到 type_code 资源属性中
      - type: copy
        from: attributes.type_code
        to: resource.type_code
```

若路径名为: `/tmp/01/0001/20251119/9999/a.txt`，则 type_code 的值为 0001

使用 debug Exporter 检查结果，效果如下

```bash
2025-11-19T10:37:12.915+0800    info    Logs    {"resource": {}, "otelcol.component.id": "debug", "otelcol.component.kind": "exporter", "otelcol.signal": "logs", "resource logs": 1, "log records": 1}
2025-11-19T10:37:12.916+0800    info    ResourceLog #0
Resource SchemaURL:
Resource attributes:
     -> type_code: Str(0001)
     -> host_ip: Str(10.10.4.90)
ScopeLogs #0
ScopeLogs SchemaURL:
InstrumentationScope
LogRecord #0
ObservedTimestamp: 2025-11-19 02:37:12.694847928 +0000 UTC
Timestamp: 1970-01-01 00:00:00 +0000 UTC
SeverityText:
SeverityNumber: Unspecified(0)
Body: Str(DesistDaydream)
Attributes:
     -> type_code: Str(0001)
     -> log.file.name: Str(a.txt)
     -> log.file.path: Str(/tmp/01/0001/20251119/9999/a.txt)
Trace ID:
Span ID:
Flags: 0
        {"resource": {}, "otelcol.component.id": "debug", "otelcol.component.kind": "exporter", "otelcol.signal": "logs"}
```

可以看到，Attributes 生成了 type_code，且 Resource attributes 中多了 type_code

> Notes: 可以把 copy 改为 move，把 type_code 从 Attributes 移到 Resource attributes。避免某些程序（e.g. Loki）读两遍

