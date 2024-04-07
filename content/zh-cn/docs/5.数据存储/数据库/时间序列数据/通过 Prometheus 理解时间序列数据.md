---
title: 通过 Prometheus 理解时间序列数据
---

## 时间序列(Time Series,简称 series) # 有序列、系列的意思

比如有这么几种描述：一系列的书、这一系列操作、等等。可以通过这种语境来理解 series(比如可以这么描述：这一系列数据)。

首先需要明确一个概念：

**Vector(向量)(也称为欧几里得向量、几何向量、矢量)**，指具有大小和方向的 **Magnitude(量)**。它可以形象化地表示为带箭头的线段。箭头所指：代表向量的方向；线段长度：代表向量的大小。与向量对应的量叫做数量（物理学中称[标量](https://baike.baidu.com/item/%E6%A0%87%E9%87%8F/1530843)），数量（或标量）只有大小，没有方向。

Prometheus 会将所有采集到的样本数据以 **TimeSeries(时间序列)** 的方式保存在内存数据库中，并且定时保存到硬盘上。TimeSeriesData 是按照**时间戳**和**值**的序列顺序存放的一条不规则有方向的线段，我们称之为 **Vector(向量)**。每条 TimeSeriesData 通过指标 **MetricsName(指标名称)** 和一组 **LabelSet(标签集)** 作为唯一标识符。如下所示，可以将 TimeSeries 理解为一个以时间为 x 轴、值为 y 轴的数字矩阵；而这个矩阵中的每一个点都是一个 **Sample(样本)**，相同 MetricName 和 LabelSet 的多个样本之间连成的线段就是时间序列数据。

```bash
  ^
  │   . . . . . . . . . . . . . . . . . . - >   node_cpu{cpu="cpu0",mode="idle"}
  │     . . . . . . . . . . . . . . . . . - >   node_cpu{cpu="cpu0",mode="system"}
  值    . . . . . . . . . . . . . . . . . - >   node_load1{}
  │     . . . . . . . . . . . . . . . . . - >
  v   . . . . . . . . . . . . . . . . . . - >
    <------------------ 时间 ---------------->
```

在 TimeSeries(时间序列) 中的每一个点称为一个 Sample(样本)，**Sample(样本)** 与 **Metric(指标) 构成了时间序列数据**，每条时间序列数据由这两部分组成：

> 下面这个例子可以看到，Prometheus 返回的时间序列数据主要有两个字段，resultType(结果类型) 与 result(结果)。result 就是这条时间序列的数据内容
> 而 result 又分为 metric 和 value。其中 value 就是指的 sample

```json
// 获取 prometheus_http_requests_total 指标。发送GET请求
// http://172.38.40.244:30001/api/v1/query?query=prometheus_http_requests_total
// 获取如下结果
{
  "status": "success",
  "data": {
    "resultType": "vector",
    "result": [
      {
        // 下面就是这条时间序列的内容
        "metric": {
          "__name__": "prometheus_http_requests_total",
          "code": "200",
          "handler": "/api/v1/label/:name/values",
          "instance": "localhost:8080",
          "job": "prometheus"
        },
        "value": [1610437100.544, "1"]
      }
    ]
  }
}
```

### Metric(指标)，就是 metric 字段

- 一个 Metrics 由 MetricName 和 描述当前样本特征的 LabelSets(所有标签的集合) 组成。

### Sample(样本)，就是 value 字段

样本中包括一个时间戳和一个样本值。有时也可以称为 指标值、时间序列值 等等，毕竟在响应体中，value 字段

- **TimesTamp(时间戳)**：一个精确到毫秒的时间戳。时间戳概念

- **SampleValue(样本值)**： 一个 float64(也可以是别的类型) 的浮点型数据表示当前样本的值。

> 样本也可以当作名词来描述这个序列的值的含义(i.e.一个数字代表了什么事物)。
> 怎么好理解怎么来，根据对 prom 的学习的不同阶段会有不同的理解。

```bash
<--------------- metric ----------------------------------------->         <-timestamp -><-value->
"__name__":"http_request_total","method":"get","statuscode":"200"},"value":[1568996888.215,"2"]
http_request_total{status="200", method="GET"}=1434417561287 => 94334
http_request_total{status="404", method="GET"}=1434417560938 => 38473
http_request_total{status="404", method="GET"}=1434417561287 => 38544
http_request_total{status="200", method="POST"}=1434417560938 => 4748
http_request_total{status="200", method="POST"}=1434417561287 => 4785
```

## Metric(指标) 详解

指标的样式一：在形式上(输出到某个程序供人阅读)，指标(Metrics)都通过如下格式标识(指标名称(metrics name)和一组标签集(LabelSet))

    <Metrics Name>{<Label Name>=<Label Value>, ...}

指标的样式二：在时间序列数据库中，指标(Metrics)则是使用下面的格式标识

    {__name__=<Metrics Name>, <Label Name>=<Label Value>, ...}

1. **Metrics Name(指标的名称)** #  可以反映被监控数据的含义（比如，http_request_total - 表示当前系统接收到的 HTTP 请求总量）。指标名称只能由 ASCII 字符、数字、下划线以及冒号组成并必须符合正则表达式\[a-zA-Z\_:]\[a-zA-Z0-9\_:]\*。
2. **LabelSet(标签集)** # 反映了当前样本的特征维度，通过这些维度 Prometheus 可以对样本数据进行过滤，聚合等。标签的名称只能由 ASCII 字符、数字以及下划线组成并满足正则表达式\[a-zA-Z\_]\[a-zA-Z0-9\_]\*。
3. 其中以\_\_作为前缀的标签(两个\_)，是获取到 metrics 后自动生成的原始标签。标签的值则可以包含任何 Unicode 编码的字符。在 Prometheus 的底层实现中指标名称实际上是以\_\_name\_\_=的形式保存在数据库中的，详见文章最后的图例

因此以下两种方式均表示的同一条 time-series ：

    api_http_requests_total{method="POST", handler="/messages"}
    等同于：
    {__name__="api_http_requests_total", method="POST", handler="/messages"}

在 Prometheus 的源码中也可以看到指标(Metric)对应的数据结构，如下所示：

    type Metric LabelSet
    type LabelSet map[LabelName]LabelValue
    type LabelName string
    type LabelValue string

## 白话说

有一条名叫内存使用率的时间序列数据，"内存使用率"就叫做 metric name。在 2019 年 10 月 1 日 00:00 的值为 100M，在 2019 年 10 月 1 日 01:00 的值为 110M。时间就是样本里的时间戳。值就是该样本的值。所有这些具有时间标识的值连在一起组成一条线，就叫时间序列数据，这条线的名字就叫“内存使用率”

可以看到，所谓的 Time Series，是使用一组标签作为唯一标识符的，可以这么说，所有标签都属于时间序列的名字，而不只是 name 字段。
