---
title: Data Model(数据模型)
linkTitle: Data Model(数据模型)
weight: 2
---

# 概述

> 参考：
>
> - 官方文档没有专门讲 Log stream 的章节，Stream 的概念都是在其他章节提到的
> - [官方文档，入门 - 标签](https://grafana.com/docs/loki/latest/get-started/labels/)
> - [官方文档，运维 - 存储](https://grafana.com/docs/loki/latest/operations/storage/)
> - https://www.qikqiak.com/post/grafana-loki-usage/

基础概念

- Indexed Label(索引标签)
- Log Stream(日志流)
- Log Line(日志行)
- Structured metadata(结构化元数据)

**Index Label(索引标签)** 简称 **Label(标签)** 是 Loki 的重要组成部分。它们使 Loki 能够将日志消息组织和分组到 **Log Stream(日志流)** 中。每个日志流必须至少有一个标签才能在 Loki 中存储和查询。

**Structured metadata(结构化元数据)** 附加到日志行上，以解决高基数的元数据（e.g. 用户名、进程 ID、etc.）作为标签导致的查询缓慢问题。并且可以适配 OTel 的资源属性如何存储的问题（因为 OTel 的资源属性设计并不会考虑这些键/值对格式的元数据是否会产生高基数，仅仅是做了一个规范定义）。

# Log Stream(日志流) 与 Indexed Labels(索引标签)

Loki 通过一种称为 **Log Stream(日志流)** 的概念，组织所有日志数据。

> [!Tip] Log Stream(日志流) 之于 Loki 类似于 Time series(时间序列) 之于 Prometheus
>
> Loki 中 Label 的概念和用法与 Prometheus 中的 Label 类似

Loki 使用 **Stream(流)** 这个词来描述保存的**一组日志数据**，由 **Indexed Label(索引标签)** 来定位每一条日志流（i.e. 每组日志）。

> [!Attention] 每条日志流必须至少拥有一个标签，才能在 Loki 中存储和查询

Stream 与 Label 是强关联的，在 Loki 中，Label 是唯一可以定义 Log Stream 的东西。每个标签键和值的组合定义了一条 log stream。如果一个标签值发生了变化，则这会生成一个新的 Log stream。在 Prometheus 中，类似 Log Stream 概念的是 time series(stream 对应 series)。但是不同的是，

用白话说，所谓的 Log Stream 可以是下面事物的一种：

- **File** # 一个文件可以是一条 Log Stream。一般情况，客户端（e.g Promtail）从文件中 tail 内容以获取日志信息。所以，在没有其他过滤条件的情况下，一个日志文件中的所有日志就相当于一条 日志流。
- **STDOUT** # 标准输出。
- ....etc.

所以 Log Stream 就是上述事物的一种通用抽象。

## Log Line(日志行)

Log Line 就是指日志流中的每一行日志，称为 **Log Line(日志行)**。Log Line 之于 Loki 类似于 Series 之于 Prometheus。

## 标签示例

下面的示例将说明 Loki 中 Label 标签的基本使用和概念。

首先看下下面的示例，这是一个 promtail 的抓取配置示例：

```yaml
scrape_configs:
  - job_name: system
    pipeline_stages:
    static_configs:
      - targets:
          - localhost
        labels:
          job: syslog
          __path__: /var/log/syslog
```

这个配置将获取日志文件数据并添加一个 job=syslog 的标签，我们可以这样来查询：

```json
{job="syslog"}
```

这将在 Loki 中创建一个流。现在我们再新增一些任务配置：

```yaml
scrape_configs:
  - job_name: system
    pipeline_stages:
    static_configs:
      - targets:
          - localhost
        labels:
          job: syslog
          __path__: /var/log/syslog
  - job_name: system
    pipeline_stages:
    static_configs:
      - targets:
          - localhost
        labels:
          job: apache
          __path__: /var/log/apache.log
```

现在我们采集两个日志文件，每个文件有一个标签与一个值，所以 Loki 会存储为两个流。我们可以通过下面几种方式来查询这些流：

```json
{job="apache"} <- 显示 job 标签为 apache 的日志
{job="syslog"} <- 显示 job 标签为 syslog 的日志
{job=~"apache|syslog"} <- 显示 job 标签为 apache 或者 syslog 的日志
```

最后一种方式我们使用的是一个 regex 标签匹配器来获取 job 标签值为 apache 或者 syslog 的日志。接下来我们看看如何使用额外的标签：

要获取这两个任务的日志可以用下面的方式来代替 regex 的方式：

```json
{env="dev"} <- 将返回所有带有 env=dev 标签的日志
```

通过使用一个标签就可以查询很多日志流了，通过组合多个不同的标签，可以创建非常灵活的日志查询。Label(标签) 是 Loki 日志数据的索引，它们用于查找压缩后的日志内容，这些内容被单独存储为块。标签和值的每一个唯一组合都定义了一个流 ，一个流的日志被分批，压缩，并作为块进行存储。

# Structured metadata(结构化元数据)

https://grafana.com/docs/loki/latest/get-started/labels/structured-metadata/

痛点：选择适当的**低基数标签**对于有效操作和查询 Loki 至关重要。某些元数据，特别是与基础设施相关的元数据（e.g. Kubernetes Pod 名称、进程 ID、etc.），可能难以嵌入日志行中，而且基数过高的话，也无法有效地存储为索引标签。

**Structured metadata(结构化元数据)** 是一种将元数据附加到日志行的方法，无需将这些元数据索引或将其包含在日志行内容本身中。

> [!Tip] 用人话说: Structured metadata 也可以看作是一种 Label（只不过 Loki 将 Label 又抽象了一层，称为 Fields，而 Fields 分为 Index Lable 和 Structured metadata。就像下图，在 Grafana 中查询的某条日志，结构化元数据与标签都在 Fields 中）
>
> ![](https://notes-learning.oss-cn-beijing.aliyuncs.com/loki/storage/metadata_and_lable_1.png)
>
> 不过这种结构化元数据的标签<font color="#ff0000">**并不会被索引**</font>
>
> 也不像 Indexed Label 会产生高基数的影响，并且没有 Indexed Label 可以决定日志流如何分组的特点。Structured metadata 只是一种附加到日志流上的标识。
>
> 有的官方文档或者博客里，还会将结构化元数据成为 Extracted Labels

> [!Note]
> 现阶段（截止到 2025-11-19），只有使用 Grafana Alloy 或 OpenTelemetry Collector 以 OpenTelemetry 格式摄取数据需要入库到 Loki 的时候，才会产生结构化元数据
>
> 结构化元数据设计之初，是旨在支持 [OpenTelemetry](/docs/6.可观测性/OpenTelemetry/OpenTelemetry.md) 数据的原生摄取。以及为了解决高基数的问题，而避免将所有键/值对的元数据都存到 Indexed Label 中（索引过多就失去了索引的作用）

在 [Log Queries](/docs/6.可观测性/Logs/Loki/LogQL/Log%20Queries.md) 中，结构化元数据无法使用日志流选择器查询，只能基于日志流选择器，通过标签过滤表达式查询。

## Indexed label 与 Structured metadata

|         | Indexed Label                            | Structured metadata                   |
| ------- | ---------------------------------------- | ------------------------------------- |
| 用途      | 作为日志流的**索引**                             | 作为日志行的**附加数据**                        |
| 储存位置    | Index                                    | Chunk                                 |
| 适合存放的数据 | 低基数字段。e.g. service_name, host, env, etc. | 任意字段。e.g. user_id, trace_id, ip, etc. |
| 查询方式    | 日志流选择器 `{KeyName="Value"}`               | 标签过滤表达式 `KeyName = "Value"`           |

## 结构化元数据转为索引标签

配置 `.limits_config.otlp_config.resource_attributes.attributes_config` 字段，将指定的 OTel 的资源属性转为标签

示例

```yaml
limits_config:
  # 将 OTLP 的 Resource attributes 映射为 Loki 的 Indexed label，以便可以通过标签过滤
  # https://grafana.com/docs/loki/latest/send-data/otel/#changing-the-default-mapping-of-otlp-to-loki-format
  otlp_config:
    resource_attributes:
      attributes_config:
        - action: index_label
          attributes:
            - host_ip
            - province
            - region
            - custom_house_name
            - log_type
            - docking_system
```

# Cardinality(基数)

https://grafana.com/docs/loki/latest/get-started/labels/cardinality/

前面的示例使用的是静态定义的 Label 标签，只有一个值；但是有一些方法可以动态定义标签。比如我们有下面这样的日志数据：

```text
11.11.11.11 - frank [25/Jan/2000:14:00:01 -0500] "GET /1986.js HTTP/1.1" 200 932 "-" "Mozilla/5.0 (Windows; U; Windows NT 5.1; de; rv:1.9.1.7) Gecko/20091221 Firefox/3.5.7 GTB6"
```

我们可以使用下面的方式来解析这条日志数据：

```yaml
- job_name: system
   pipeline_stages:
      - regex:
        expression: "^(?P<ip>\\S+) (?P<identd>\\S+) (?P<user>\\S+) \\[(?P<timestamp>[\\w:/]+\\s[+\\-]\\d{4})\\] \"(?P<action>\\S+)\\s?(?P<path>\\S+)?\\s?(?P<protocol>\\S+)?\" (?P<status_code>\\d{3}|-) (?P<size>\\d+|-)\\s?\"?(?P<referer>[^\"]*)\"?\\s?\"?(?P<useragent>[^\"]*)?\"?$"
    - labels:
        action:
        status_code:
   static_configs:
   - targets:
      - localhost
     labels:
      job: apache
      env: dev
      __path__: /var/log/apache.log
```

这个 regex 匹配日志行的每个组件，并将每个组件的值提取到一个 capture 组里面。在 pipeline 代码内部，这些数据被放置到一个临时的数据结构中，允许在处理该日志行时将其用于其他处理（此时，临时数据将被丢弃）。

从该 regex 中，我们就使用其中的两个 capture 组，根据日志行本身的内容动态地设置两个标签：

```text
action (例如 action="GET", action="POST") status_code (例如 status_code="200", status_code="400")
```

假设我们有下面几行日志数据：

```text
11.11.11.11 - frank [25/Jan/2000:14:00:01 -0500] "GET /1986.js HTTP/1.1" 200 932 "-" "Mozilla/5.0 (Windows; U; Windows NT 5.1; de; rv:1.9.1.7) Gecko/20091221 Firefox/3.5.7 GTB6"
11.11.11.12 - frank [25/Jan/2000:14:00:02 -0500] "POST /1986.js HTTP/1.1" 200 932 "-" "Mozilla/5.0 (Windows; U; Windows NT 5.1; de; rv:1.9.1.7) Gecko/20091221 Firefox/3.5.7 GTB6"
11.11.11.13 - frank [25/Jan/2000:14:00:03 -0500] "GET /1986.js HTTP/1.1" 400 932 "-" "Mozilla/5.0 (Windows; U; Windows NT 5.1; de; rv:1.9.1.7) Gecko/20091221 Firefox/3.5.7 GTB6"
11.11.11.14 - frank [25/Jan/2000:14:00:04 -0500] "POST /1986.js HTTP/1.1" 400 932 "-" "Mozilla/5.0 (Windows; U; Windows NT 5.1; de; rv:1.9.1.7) Gecko/20091221 Firefox/3.5.7 GTB6"
```

则在 Loki 中收集日志后，会创建为如下所示的流：

```text
{job="apache",env="dev",action="GET",status_code="200"} 11.11.11.11 - frank [25/Jan/2000:14:00:01 -0500] "GET /1986.js HTTP/1.1" 200 932 "-" "Mozilla/5.0 (Windows; U; Windows NT 5.1; de; rv:1.9.1.7) Gecko/20091221 Firefox/3.5.7 GTB6"
{job="apache",env="dev",action="POST",status_code="200"} 11.11.11.12 - frank [25/Jan/2000:14:00:02 -0500] "POST /1986.js HTTP/1.1" 200 932 "-" "Mozilla/5.0 (Windows; U; Windows NT 5.1; de; rv:1.9.1.7) Gecko/20091221 Firefox/3.5.7 GTB6"
{job="apache",env="dev",action="GET",status_code="400"} 11.11.11.13 - frank [25/Jan/2000:14:00:03 -0500] "GET /1986.js HTTP/1.1" 400 932 "-" "Mozilla/5.0 (Windows; U; Windows NT 5.1; de; rv:1.9.1.7) Gecko/20091221 Firefox/3.5.7 GTB6"
{job="apache",env="dev",action="POST",status_code="400"} 11.11.11.14 - frank [25/Jan/2000:14:00:04 -0500] "POST /1986.js HTTP/1.1" 400 932 "-" "Mozilla/5.0 (Windows; U; Windows NT 5.1; de; rv:1.9.1.7) Gecko/20091221 Firefox/3.5.7 GTB6"
```

这 4 行日志将成为 4 个独立的流，并开始填充 4 个独立的块。任何与这些 标签/值 组合相匹配的额外日志行将被添加到现有的流中。如果有另一个独特的标签组合进来（比如 status_code=“500”）就会创建另一个新的流。

比如我们为 IP 设置一个 Label 标签，不仅用户的每一个请求都会变成一个唯一的流，每一个来自同一用户的不同 action 或 status_code 的请求都会得到自己的流。

如果有 4 个共同的操作（GET、PUT、POST、DELETE）和 4 个共同的状态码（可能不止 4 个！），这将会是 16 个流和 16 个独立的块。然后现在乘以每个用户，如果我们使用 IP 的标签，你将很快就会有数千或数万个流了。

这个 Cardinality 太高了，这足以让 Loki 挂掉。

当我们谈论 Cardinality 的时候，我们指的是标签和值的组合，以及他们创建的流的数量，高 Cardinality 是指使用具有较大范围的可能值的标签，如 IP，或结合需要其他标签，即使它们有一个小而有限的集合，比如 status_code 和 action。

高 Cardinality 会导致 Loki 建立一个巨大的索引（成本高），并将成千上万的微小块存入对象存储中（速度慢），Loki 目前在这种配置下的性能非常差，运行和使用起来非常不划算的。

# Automatic stream sharding

**Automatic stream sharding(自动流分片)** 功能会在数据发送到 Loki 的速率超过配置的限额时，为日志流添加 `__stream_shard__` 标签，标签的值是从 0 开始的数字。

该功能是 Loki 在处理大型日志流产生的解决方案，避免数据丢失。详见: https://grafana.com/docs/loki/latest/operations/automatic-stream-sharding/#when-to-use-automatic-stream-sharding
