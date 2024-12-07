---
title: Data Model(数据模型)
linkTitle: Data Model(数据模型)
date: 2021-10-22T22:05:00
weight: 20
---

# 概述

> 参考：
>
> - 官方文档没有专门讲 Log Stream 的章节，Stream 的概念都是在其他章节提到的
> - [官方文档，入门-标签](https://grafana.com/docs/loki/latest/getting-started/labels/)
> - [官方文档，运维-存储](https://grafana.com/docs/loki/latest/operations/storage/)

# Log Stream(日志流) 概念

Loki 通过一种称为 **Log Stream(日志流)** 的概念组织所有日志数据。**Log Stream(日志流) 之于 Loki 类似于 Time series(时间序列) 之于 Prometheus**

Loki 使用 **Stream(流)** 这个词来描述保存的日志数据，并根据 **Label(标签)** 来定位日志流，Label 是日志流的元数据。Label 的概念和用法与 Prometheus 中的 Label 一致。如果 Loki 与 Prometheus 同时使用，那么他们之间得标签是一致的，通过 Label，很容易得就可以将应用程序指标和日志数据关联起来。

Stream 与 Label 是强关联的，在 Loki 中，Label 是唯一可以定义 Log Stream 的东西。每个标签键和值的组合定义了一条 log stream。如果一个标签值发生了变化，则这会生成一个新的 Log stream。在 Prometheus 中，类似 Log Stream 概念的是 time series(stream 对应 series)。但是不同的是，在 Prometheus 中还有一个维度，是 metrics name(指标名称)。但是在 Loki 中则谁 Path，一个 采集日志的 Path 实际上是会采集很多很多日志的。也正是由于此，所以 Loki 将这种概念称为 Stream，而不是 Series。

用白话说，所谓的 Log Stream 可以是下面事物的一种：

- **File** # 一个文件就是一个 Log Stream。一般情况，客户端(比如 Promtail)从文件中 tail 内容以获取日志信息，所以，一个日志就相当于一个 日志流。
- **STDOUT** # 标准输出。
- ....等等

所以 Log Stream 就是上述事物的一种通用抽象。

## Log Line(日志行) 概念

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

通过使用一个标签就可以查询很多日志流了，通过组合多个不同的标签，可以创建非常灵活的日志查询。Label 标签是 Loki 日志数据的索引，它们用于查找压缩后的日志内容，这些内容被单独存储为块。标签和值的每一个唯一组合都定义了一个流 ，一个流的日志被分批，压缩，并作为块进行存储。

# Cardinality(基数) 概念

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

**Automatic stream sharding(自动流分片)** 功能会在数据发送到 Loki 的速率超过配置的限额时，为日志流添加 `__stream_shard__` 标签，标签的只是从 0 开始的数字。

该功能是 Loki 在处理大型日志流产生的解决方案，避免数据丢失。详见: https://grafana.com/docs/loki/latest/operations/automatic-stream-sharding/#when-to-use-automatic-stream-sharding

