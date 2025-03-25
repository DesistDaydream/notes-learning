---
title: Configuration
linkTitle: Configuration
weight: 20
---

# 概述

> 参考：
>
> - [官方文档，Collector - 配置](https://opentelemetry.io/docs/collector/configuration)

顶层字段

- **receivers**(map\[STRING][receivers](#receivers)) # 配置 Receivers 组件
- **processors**(map\[STRING][processors](#processors)) # 配置 Processors 组件
- **exporters**(map\[STRING][exporters](#exporters)) # 配置 Exporters 组件
- **extensions**(map\[STRING][extensions](#extensions)) # 配置 扩展
- **connectors** # TODO
- **service**([service](#service)) # 配置在处理各类可观测数据时，使用哪些扩展、使用哪些组件

在 [otelcol/config.go](https://github.com/open-telemetry/opentelemetry-collector/blob/v0.112.0/otelcol/config.go#L21) 可以看到顶层字段的 struct

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

Notes: 这些 STRING 用来表示**组件 ID**。在 [component/identifiable.go](https://github.com/open-telemetry/opentelemetry-collector/blob/v0.112.0/component/identifiable.go#L19) 和 文档中 可以看到，定义组件 ID 遵循 `type[/name]` 格式，e.g. otlp 或 otlp/2。只要 ID 是唯一的，就可以多次定义给定类型的组件。例如：

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

type 可用的字符串可以参考 Collector [Component](/docs/6.可观测性/OpenTelemetry/Component.md)，像下图，在 opentelemetry-collector、opentelemetry-contrib 项目中，在对应的组件目录下，每个目录都是一个可用的 type，目录名的前缀就是 type 的名称，比如 otlpreceiver 是 OTLP Receiver，prometheusreveiver 是 [Prometheus](/docs/6.可观测性/Metrics/Prometheus/Prometheus.md) Receiver，以此类推。

![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/otel/config_type_desc.png)

# receivers

# processors

# exporters

# connectors

# service
