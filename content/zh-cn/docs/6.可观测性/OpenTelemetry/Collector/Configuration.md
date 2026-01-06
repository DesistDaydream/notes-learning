---
title: Configuration
linkTitle: Configuration
weight: 2
---

# 概述

> 参考：
>
> - [官方文档，Collector - 配置](https://opentelemetry.io/docs/collector/configuration)

顶层字段

- **receivers**(map\[STRING][receivers](#receivers)) # 配置 Receivers 管道组件
- **processors**(map\[STRING][processors](#processors)) # 配置 Processors 管道组件
- **exporters**(map\[STRING][exporters](#exporters)) # 配置 Exporters 管道组件
- **extensions**(map\[STRING][extensions](#extensions)) # 配置 扩展
- **connectors**(map\[STRING][connectors](#connectors)) # TODO: 配置 [Connectors](https://opentelemetry.io/docs/collector/configuration/#connectors) 管道组件
- **service**([service](#service)) # 在处理各类可观测数据时，使用哪些扩展、使用哪些组件。每个 service 可以简单理解为一个 Pipeline(管道)。

在 [otelcol/config.go](https://github.com/open-telemetry/opentelemetry-collector/blob/v0.126.0/otelcol/config.go#L21) 可以看到顶层字段的 struct

```go
// Config defines the configuration for the various elements of collector or agent.
type Config struct {
 // Receivers is a map of ComponentID to Receivers.
 Receivers map[component.ID]component.Config `mapstructure:"receivers"`

 // Exporters is a map of ComponentID to Exporters.
 Exporters map[component.ID]component.Config `mapstructure:"exporters"`

 // Processors is a map of ComponentID to Processors.
 Processors map[component.ID]component.Config `mapstructure:"processors"`

 // Connectors is a map of ComponentID to connectors.
 Connectors map[component.ID]component.Config `mapstructure:"connectors"`

 // Extensions is a map of ComponentID to extensions.
 Extensions map[component.ID]component.Config `mapstructure:"extensions"`

 Service service.Config `mapstructure:"service"`
}
```

> Notes: map 中的 key 用来表示**组件 ID**。

## 组件 ID 命名规范

`TYPE[/NAME]`

定义组件 ID 遵循 `TYPE[/NAME]` 格式，在 [component/identifiable.go](https://github.com/open-telemetry/opentelemetry-collector/blob/v0.140.0/component/identifiable.go#L85) 和 [官方文档中](https://opentelemetry.io/docs/collector/configuration/#basics) 可以看到。

TYPE 可用的字符串可以参考 Collector [Collector component](/docs/6.可观测性/OpenTelemetry/Collector/Collector%20component.md)，像下图，在 opentelemetry-collector、opentelemetry-collector-contrib 项目中，其对应的组件目录下，每个目录都是一个可用的 TYPE，目录名的前缀就是 TYPE，比如 otlpreceiver 是 OTLP Receiver（TYPE 是 otlp），prometheusreveiver 是 [Prometheus](/docs/6.可观测性/Metrics/Prometheus/Prometheus.md) Receiver（TYPE 是 prometheus），以此类推。

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/otel/config_type_desc.png)

e.g. otlp 或 otlp/2。只要 ID 是唯一的，就可以多次定义给定种类的组件。例如：

```yaml
receivers:
  otlp:
    protocols:
      grpc:
        endpoint: 0.0.0.0:4317
      http:
        endpoint: 0.0.0.0:4318
  otlp/2:
    protocols:
      grpc:
        endpoint: 0.0.0.0:55690

processors:
  batch:
  batch/test:

exporters:
  otlp:
    endpoint: otelcol:4317
  otlp/2:
    endpoint: otelcol2:4317

extensions:
  health_check:
  pprof:
  zpages:

service:
  extensions: [health_check, pprof, zpages]
  pipelines:
    traces:
      receivers: [otlp]
      processors: [batch]
      exporters: [otlp]
    traces/2:
      receivers: [otlp/2]
      processors: [batch/test]
      exporters: [otlp/2]
    metrics:
      receivers: [otlp]
      processors: [batch]
      exporters: [otlp]
    logs:
      receivers: [otlp]
      processors: [batch]
      exporters: [otlp]
```

# receivers

配置内容取决于使用的 Receiver 类型

# processors

配置内容取决于使用的 Processor 类型

# exporters

配置内容取决于使用的 Exporter 类型

# connectors

配置内容取决于使用的 Connector 类型

# service

https://github.com/open-telemetry/opentelemetry-collector/blob/ee8646eb459f8cbd3b858a5871910065a5a68c03/service/config.go#L13

**extensions**(\[]STRING)

**pipelines**(map\[STRING][pipelines](#pipelines)) # 定义管道。map 中的 key 是管道 ID，也遵循 `TYPE[/NAME]` 格式。TYEP 可以用的值有: traces, metrics, logs

**telemetry**(Object) # 配置 Collector 本身的与组件无关的配置。e.g. 内部指标暴露端口、日志级别、etc.

- 默认值: https://github.com/open-telemetry/opentelemetry-collector/blob/ee8646eb459f8cbd3b858a5871910065a5a68c03/service/telemetry/otelconftelemetry/factory.go#L43
- 示例：

```yaml
service:
  telemetry:
    metrics:
      readers:
        - pull:
            exporter:
              prometheus:
                host: 0.0.0.0
                port: 8888
                without_units: true
```

## pipelines

**receivers**(\[]CommandID)

**processors**(\[]CommandID)

**exporters**(\[]CommandID)

# 配置示例

将本地文件日志写入到 Loki 的最简单配置

```yaml
receivers:
  filelog:
    include: [/var/log/syslog]

exporters:
  otlphttp/loki:
    endpoint: http://localhost:3100/otlp

service:
  pipelines:
    logs:
      receivers: [filelog]
      processors: []
      exporters: [otlphttp/loki]
```