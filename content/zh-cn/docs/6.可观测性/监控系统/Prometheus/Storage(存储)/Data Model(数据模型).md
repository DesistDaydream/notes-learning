---
title: Data Model(数据模型)
---

# 概述

> 参考：
>
> - [官方文档，概念-数据模型](https://prometheus.io/docs/concepts/data_model/)
> - [yunlzheng 文档](https://yunlzheng.gitbook.io/prometheus-book/parti-prometheus-ji-chu/promql)

Prometheus 从根本上将所有数据存储为 [Time Series(时间序列)](https://en.wikipedia.org/wiki/Time_series)：属于同一度量标准和同一组标注维的带有时间戳的值流。除了存储的时间序列外，Prometheus 可能会生成临时派生的时间序列作为查询的结果。

## Time-Series Data(时间序列数据) 概念

> 参考：
>
> - [Wiki](https://en.wikipedia.org/wiki/Time_series)
> - [InfluxDB 对时间序列数据的定义](https://www.influxdata.com/what-is-time-series-data/)
> - [这是论文](http://get.influxdata.com/rs/972-GDU-533/images/why%20time%20series.pdf)

**Time Series(时间序列)** 是一组按照时间发生先后顺序进行排列的数据点序列。通常一组时间序列的时间间隔为一恒定值（如 1 秒，5 分钟，12 小时，7 天，1 年），因此时间序列可以作为离散时间数据进行分析处理。时间序列广泛应用于数理统计、信号处理、模式识别、计量经济学、数学金融、天气预报、地震预测、脑电图、控制工程、航空学、通信工程以及绝大多数涉及到时间数据测量的应用科学与工程学。

**[Time Series Data](/docs/5.数据存储/数据库/时间序列数据/时间序列数据.md)(时间序列数据，简称 series)** 是在一段时间内通过重复 Measurement(测量) 而获得的观测值的集合；可以将这些观测值绘制于图形之上，它会有一个数据轴和一个时间轴。

从另一个角度看，时间序列数据是在不同时间上收集到的数据，用于所描述现象随时间变化的情况。这类数据反映了某一事物、现象等随时间的变化状态或程度。

## Prometheus 中时间序列数据的组成

**时间序列(Time Series,简称 series) 有序列、系列的意思**。比如有这么几种描述：一系列的书、这一系列操作、等等。可以通过这种语境来理解 series(比如可以这么描述：这一系列数据)。

与传统意义上定义的时序数据一样，由两部分组成：

- **Metrics(指标)** # 用来描述要采集的数据指标，是时序数据的唯一标识符。例如：检测各个城市的风力、系统内存已使用的字节数 等等。相当于关系型数据库中的表。
- **Sample(样本)** # 针对监测对象的某项指标(由 Metric 和 Tag 定义)按特定时间间隔采集到的每个 Metric 值就是一个 Sample(样本)。类似关系型数据库中的一行。

首先需要明确一个概念：

**Vector(向量)(也称为欧几里得向量、几何向量、矢量)**，指具有大小和方向的 **Magnitude(量)**。它可以形象化地表示为带箭头的线段。箭头所指：代表向量的方向；线段长度：代表向量的大小。与向量对应的量叫做数量（物理学中称[标量](https://baike.baidu.com/item/%E6%A0%87%E9%87%8F/1530843)），数量（或标量）只有大小，没有方向。

Prometheus 会将所有采集到的样本数据以 **TimeSeries(时间序列)** 的方式保存在内存数据库中，并且定时保存到硬盘上。TimeSeriesData 是按照**时间戳**和**值**的序列顺序存放的一条不规则有方向的线段，我们称之为 **Vector(向量)**。每条 TimeSeriesData 通过 **MetricsName(指标名称)** 和一组 **LabelSet(标签集)** 作为唯一标识符。如下所示，可以将 TimeSeries 理解为一个以时间为 x 轴、值为 y 轴的数字矩阵；而这个矩阵中的每一个点都是一个 **Sample(样本)**，相同 MetricName 和 LabelSet 的多个样本之间连成的线段就是时间序列数据。

```bash
  ^
  │   . . . . . . . . . . . . . . . . . . - >   node_cpu{cpu="cpu0",mode="idle"} # node_cpu 是 MetricsName(指标名称)
  │     . . . . . . . . . . . . . . . . . - >   node_cpu{cpu="cpu0",mode="system"} # {cpu="cpu0",mode="system"} 是 LabelSet(标签集)
  值    . . . . . . . . . . . . . . . . . - >   node_load1{} # 这个时间序列只有指标名称，没有标签集
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

## Metric(指标) 结构

指标的样式一：在形式上(输出到某个程序供人阅读)，指标(Metrics)都通过如下格式标识(指标名称(metrics name)和一组标签集(LabelSet))

```text
<Metrics Name>{<Label Name>=<Label Value>, ...}
```

指标的样式二：在时间序列数据库中，指标(Metrics)则是使用下面的格式标识

```json
{__name__=<Metrics Name>, <Label Name>=<Label Value>, ...}
```

1. **Metrics Name(指标的名称)** # 可以反映被监控数据的含义（比如，http_request_total - 表示当前系统接收到的 HTTP 请求总量）。指标名称只能由 ASCII 字符、数字、下划线以及冒号组成并必须符合正则表达式\[a-zA-Z\_:]\[a-zA-Z0-9\_:]\*。
2. **LabelSet(标签集)** # 反映了当前样本的特征维度，通过这些维度 Prometheus 可以对样本数据进行过滤，聚合等。标签的名称只能由 ASCII 字符、数字以及下划线组成并满足正则表达式\[a-zA-Z\_]\[a-zA-Z0-9\_]\*。
3. 其中以\_\_作为前缀的标签(两个\_)，是获取到 metrics 后自动生成的原始标签。标签的值则可以包含任何 Unicode 编码的字符。在 Prometheus 的底层实现中指标名称实际上是以\_\_name\_\_=的形式保存在数据库中的，详见文章最后的图例

因此以下两种方式均表示的同一条 time-series ：

```text
api_http_requests_total{method="POST", handler="/messages"}
等同于：
{__name__="api_http_requests_total", method="POST", handler="/messages"}
```

在 Prometheus 的源码中也可以看到指标(Metric)对应的数据结构，如下所示：

```
type Metric LabelSet
type LabelSet map[LabelName]LabelValue
type LabelName string
type LabelValue string
```

## 白话说

有一条名叫内存使用率的时间序列数据，"内存使用率"就叫做 metric name。在 2019 年 10 月 1 日 00:00 的值为 100M，在 2019 年 10 月 1 日 01:00 的值为 110M。时间就是样本里的时间戳。值就是该样本的值。所有这些具有时间标识的值连在一起组成一条线，就叫时间序列数据，这条线的名字就叫“内存使用率”

可以看到，所谓的 Time Series，是使用一组标签作为唯一标识符的，可以这么说，所有标签都属于时间序列的名字，而不只是 name 字段。

# Metrics(指标) 的类型

> 参考：
> 
> [官网文档，概念-metric 类型](https://prometheus.io/docs/concepts/metric_types/)

在 Prometheus 的存储实现上所有的监控样本都是以 time-series 的形式保存在 Prometheus 的 TSDB(时序数据库) 中，而 TimeSeries 所对应的 Metric(监控指标) 也是通过 LabelSet 进行唯一命名的。

从存储上来讲所有的 Metrics 都是相同的，但是在不同的场景下这些 Metrics 又有一些细微的差异。 例如，在 Node Exporter 返回的样本中指标 node_load1 反应的是当前系统的负载状态，随着时间的变化这个指标返回的样本数据是在不断变化的。而指标 node_cpu 所获取到的样本数据却不同，它是一个持续增大的值，因为其反应的是 CPU 的累积使用时间，从理论上讲只要系统不关机，这个值是会无限变大的。

为了能够帮助用户理解和区分这些不同监控指标之间的差异，Prometheus 定义了 4 中不同的 **Metric Type(指标类型)**：Counter(计数器)、Gauge(计量器)、Histogram(直方图)、Summary(摘要)。

在 Exporter 返回的样本数据中，其注释中也包含了该样本的类型。例如：

> 其中 TYPE node_cpu counter 表明 node_cpu 的指标类型为 counter

```text
# HELP node_cpu Seconds the cpus spent in each mode.
# TYPE node_cpu counter
node_cpu{cpu="cpu0",mode="idle"} 362812.789625
```

## Counter(计数器) # 只增不减的计数器

Counter 类型的指标其工作方式和计数器一样，只增不减（除非系统发生重置）。常见的监控指标，如 http_requests_total，node_cpu 都是 Counter 类型的监控指标。 一般在定义 Counter 类型指标的名称时推荐使用\_total 作为后缀。

Counter 是一个简单但有强大的工具，例如我们可以在应用程序中记录某些事件发生的次数，通过以时序的形式存储这些数据，我们可以轻松的了解该事件产生速率的变化。 PromQL 内置的聚合操作和函数可以让用户对这些数据进行进一步的分析：

例如，通过 rate()函数获取 HTTP 请求量的增长率：

rate(http_requests_total\[5m])

查询当前系统中，访问量前 10 的 HTTP 地址：

topk(10, http_requests_total)

## Gauge(仪表盘) # 可增可减的 Gauge

与 Counter 不同，Gauge 类型的指标侧重于反应系统的当前状态。因此这类指标的样本数据可增可减。常见指标如：node_memory_MemFree(主机当前空闲的内容大小)、node_memory_MemAvailable(可用内存大小)都是 Gauge 类型的监控指标。

通过 Gauge 指标，用户可以直接查看系统的当前状态：

`node_memory_MemFree`

对于 Gauge 类型的监控指标，通过 PromQL 内置函数 delta() 可以获取样本在一段时间返回内的变化情况。例如，计算 CPU 温度在两个小时内的差异：

`delta(cpu_temp_celsius{host="zeus"}[2h])`

还可以使用 deriv() 计算样本的线性回归模型，甚至是直接使用 predict_linear() 对数据的变化趋势进行预测。例如，预测系统磁盘空间在 4 个小时之后的剩余情况：

`predict_linear(node_filesystem_free{job="node"}[1h], 4 * 3600)`

## Histogram(直方图) 与 Summary(摘要)

> 参考：
>
> - [官方文档](https://prometheus.io/docs/practices/histograms/)
> - [云原生实验室](https://fuckcloudnative.io/posts/prometheus-histograms/)

除了 Counter 和 Gauge 类型的监控指标以外，Prometheus 还定义了 Histogram 和 Summary 的指标类型。Histogram 和 Summary 主用用于统计和分析样本的分布情况。

在大多数情况下人们都倾向于使用某些量化指标的平均值，例如 CPU 的平均使用率、页面的平均响应时间。这种方式的问题很明显，以系统 API 调用的平均响应时间为例：如果大多数 API 请求都维持在 100ms 的响应时间范围内，而个别请求的响应时间需要 5s，那么就会导致某些 WEB 页面的响应时间落到中位数的情况，而这种现象被称为长尾问题。

为了区分是平均的慢还是长尾的慢，最简单的方式就是按照请求延迟的范围进行分组。例如，统计延迟在 0~10ms 之间的请求数有多少而 10~20ms 之间的请求数又有多少。通过这种方式可以快速分析系统慢的原因。Histogram 和 Summary 都是为了能够解决这样问题的存在，通过 Histogram 和 Summary 类型的监控指标，我们可以快速了解监控样本的分布情况。

直方图和摘要均是样本观察值，也就是说在一段时间内持续观察某个样本后得出的数据。

这两种指标类型都属于**统计学范畴的指标**。

### Histogram(直方图)

Histogram 在**一段时间范围内**观察某指标(通常是 请求的持续时间 或 响应时间的长短 等)，并对该指标的样本进行采样，将其计入可配置的 **bucket(储存区)** 中。Histogram 还提供所有 observed(被观察指标) 的样本在这一段时间范围内的总和。

**传统意义上的直方图**
假设我们想获取某个应用在不同响应时间的次数，则首先需要获取该应用在一段时间内的响应时间，收集这些样本。假设最后得到的所有样本的响应时间范围是 0s~10s。那么我们将样本的值划分为不同的区间，这个区间就是 **bucket(存储区)**，假设每个 bucket 的宽度是 0.2s，那么第一个 bucket 则表示响应时间小于 0.2s 的所有样本数量；第二个 bucket 表示响应时间大于 0.2s 且小于 0.4s 的样本数量；以此类推。效果如图：

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/gnzmdt/1617542272313-06f1c876-f41b-431b-99ed-e1e14b152761.jpeg)

**Prometheus 中的直方图**

Prometheus 中的直方图与传统意义的直方图有一些差别，准确描述，应该称为 **累计直方图**。主要差别在 bucket 的定义，在 Prometheus 的累计直方图中，还是假设 bucket 的宽度为 0.2s，那么第一个 bucket 表示响应时间小于等于 0.2s 的样本数量，第二个 bucket 表示响应时间小于等于 0.4s 的样本数量，以此类推。也就是说，**每一个 bucket 中的样本都包含了卡面所有 bucket 中的样本**，所以称为 累计直方图。而每个 bucket 范围中的最大值，称为 **upper inclusive bound(上边界)**。效果如图：

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/gnzmdt/1617543132157-1f35ddae-ef51-4ca7-9ece-72886154cc1f.jpeg)

Histogram 类型的指标在同一时间具有多条时间序列(假设指标名称为 <basename)：(这些时间序列大体分为 3 种)

- **`<basename>_bucket{le="<上边界>"}`**# 要观察的样本分布在 bucket 中的数量。解释的更通俗易懂一点，这个值表示 要观察的样本的值 小于等于 上边界的值 的所有样本数量。
  - le 通常用来表示该 bucket 的上限。le 这俩字符按照关系运算符来理解，就是“小于或等于”的意思。。。。。。。
  - 用白话说就是，le 是 bucket 的标识符，比如下面的示例，就可以这么描述：0 到 0.00025 储存区，含有 332 个样本；0 到 0.0005 储存区，含有 336 个样本。
- **`<basename>_sum`** # 所有 被观察样本的值 的总和。
- **`<basename>_count`** # 观察次数。(该值和上面的 \<basename>\_bucket{le="+Inf"} 相同)
  - 本质上是一个 Counter 类型的指标。

在 coredns 的样本数据中，我们能找到类型为 Histogram 的监控指标 `coredns_dns_request_duration_seconds`

在这些时间序列中，被观察的样本是“每个 dns 的解析请求所花费的时间”

```bash
# HELP coredns_dns_request_duration_seconds Histogram of the time (in seconds) each request took.
# TYPE coredns_dns_request_duration_seconds histogram
# 在总共336次解析请求的花费时间中，小于0.00025秒的有332次
coredns_dns_request_duration_seconds_bucket{server="dns://:53",zone=".",le="0.00025"} 332
# 在总共336次解析请求的花费时间中，小于0.0005秒的有336次
coredns_dns_request_duration_seconds_bucket{server="dns://:53",zone=".",le="0.0005"} 336
........
# 在总共336次解析请求的花费时间中，小于8.192秒的有336次
coredns_dns_request_duration_seconds_bucket{server="dns://:53",zone=".",le="8.192"} 336
coredns_dns_request_duration_seconds_bucket{server="dns://:53",zone=".",le="+Inf"} 336
# 所有的336次解析请求，总的花费时间为 0.03502086400000001 秒
coredns_dns_request_duration_seconds_sum{server="dns://:53",zone="."} 0.03502086400000001
# DNS 解析请求一共336次
coredns_dns_request_duration_seconds_count{server="dns://:53",zone="."} 336
```

可以通过 histogram_quantile() 函数 来计算 Histogram 类型样本的 分位数。分位数可能不太好理解，你可以理解为分割数据的点。我举个例子，假设样本的 9 分位数（quantile=0.9）的值为 x，即表示小于 x 的采样值的数量占总体采样值的 90%。Histogram 还可以用来计算应用性能指标值（Apdex score）。

注意：

bucket 可以理解为是对数据指标值域的一个划分，划分的依据应该基于数据值的分布。注意后面的采样点是包含前面的采样点的，假设 xxx_bucket{...,le="0.01"} 的值为 10，而 xxx_bucket{...,le="0.05"} 的值为 30，那么意味着这 30 个采样点中，有 10 个是小于 10 ms 的，其余 20 个采样点的响应时间是介于 10 ms 和 50 ms 之间的。

用白话说：直方图与 Counter 和 Gauge 的本质区别在于，直方图是对一组样本进行统计获得的结果，而 Counter 和 Gauge 仅仅是一个单一的样本。
直方图的应用场景：在 1 小时的 http 请求中，有多少请求的响应时间少于 1 秒，有多少请求的响应时间少于 2 秒，总有有多少请求，所有请求的平均的响应时间是多少。

如果是 Guage 的话，则只能表示每一个请求的具体响应时间，或者总共有多少个请求。

所以才说，直方图就是一种**统计学上的指标**

### Summary(摘要)

与 Histogram 类型类似，用于表示一段时间内的数据采样结果（通常是请求持续时间或响应大小等），但它 bucket 表示分位数（通过客户端计算，然后展示出来），而不是通过区间来计算。

例如，指标 prometheus_tsdb_wal_fsync_duration_seconds 的指标类型为 Summary。 它记录了 Prometheus Server 中 wal_fsync 处理的处理时间，通过访问 Prometheus Serve r 的 /metrics 地址，可以获取到以下监控样本数据：

```text
# HELP prometheus_tsdb_wal_fsync_duration_seconds Duration of WAL fsync.
# TYPE prometheus_tsdb_wal_fsync_duration_seconds summary
prometheus_tsdb_wal_fsync_duration_seconds{quantile="0.5"} 0.012352463
prometheus_tsdb_wal_fsync_duration_seconds{quantile="0.9"} 0.014458005
prometheus_tsdb_wal_fsync_duration_seconds{quantile="0.99"} 0.017316173
prometheus_tsdb_wal_fsync_duration_seconds_sum 2.888716127000002
prometheus_tsdb_wal_fsync_duration_seconds_count 216
```

从上面的样本中可以得知当前 Prometheus Server 进行 wal_fsync 操作的总次数为 216 次，耗时 2.888716127000002s。其中中位数（quantile=0.5）的耗时为 0.012352463，9 分位数（quantile=0.9）的耗时为 0.014458005s。

### Summary 类型 与 Histogram 类型 的异同

1. 两类样本同样会反应当前指标的记录的总数(以\_count 作为后缀)以及其值的总量（以\_sum 作为后缀）
2. 不同在于 Histogram 指标直接反应了在不同区间内样本的个数，区间通过标签 len 进行定义。
3. 对于分位数的计算而言，Histogram 通过 histogram_quantile 函数是在服务器端计算的分位数。 而 Sumamry 的分位数则是直接在客户端计算完成。
4. Summary 在通过 PromQL 进行查询时有更好的性能表现，而 Histogram 则会消耗更多的资源。反之对于提供指标的服务而言，Histogram 消耗的资源更少。在选择这两种方式时用户应该按照自己的实际场景进行选择。

# Prometheus 底层保存的时间序列数据的样例

详见 [Querying API](/docs/6.可观测性/监控系统/Prometheus/Prometheus%20API/Querying%20API.md)

下面红框的地方就是

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/gnzmdt/1616069604244-2f845e27-61ec-4a5b-ab9a-0634bf8907b2.jpeg)

# Prometheus 格式的 Metrics 详解

> 参考：
>
> - [官方文档](https://prometheus.io/docs/instrumenting/exposition_formats/)

通过各种 Exporter 暴露的 HTTP 服务，Prometheus 可以采集到 当前时间 主机所有监控指标的样本数据。数据格式示例如下：

```text
# HELP http_requests_total The total number of HTTP requests.
# TYPE http_requests_total counter
http_requests_total{method="post",code="200"} 1027 1395066363000
http_requests_total{method="post",code="400"}    3 1395066363000
# Escaping in label values:
msdos_file_access_time_seconds{path="C:\\DIR\\FILE.TXT",error="Cannot find file:\n\"FILE.TXT\""} 1.458255915e9
# Minimalistic line:
metric_without_timestamp_and_labels 12.47
# A weird metric from before the epoch:
something_weird{problem="division by zero"} +Inf -3982045
```

Note：上面通过 http 采集到的数据就是文本格式的 Metrics，格式一定是上述的样子，每个时间序列都分为 3 个部分。

1. `# HELP` 时间序列名称 时间序列描述
2. `# TYPE` 时间序列名称 时间序列类型
3. 非#开头的每一行表示当前 Node Exporter 采集到的一个监控样本：node_cpu 和 node_load1 表明了当前指标的名称、大括号中的标签则反映了当前样本的一些特征和维度、浮点数则是该监控样本的具体值。
   1. 如果有多个 Metrics 的项目，则会有多行

主要由三个部分组成：样本的一般注释信息（HELP），样本的类型注释信息（TYPE）和样本。Prometheus 会对 Exporter 响应的内容逐行解析：

如果当前行以 `# HELP` 开始，Prometheus 将会按照以下规则对内容进行解析，得到当前的指标名称以及相应的说明信息：

## # HELP

如果当前行以 `# TYPE` 开始，Prometheus 会按照以下规则对内容进行解析，得到当前的指标名称以及指标类型:

## # TYPE

TYPE 注释行必须出现在指标的第一个样本之前。如果没有明确的指标类型需要返回为 untyped。

## MetricsName 与 Metrics 的值

除了 `#` 开头的所有行都会被视为是监控样本数据。 每一行样本需要满足以下格式规范:

```bash
metric_name [{ label_name = "label_value" { , label_name = "label_value" } [ ,... ] }] value [ timestamp ]
```

其中 metric_name 和 label_name 必须遵循 PromQL 的格式规范要求。value 是一个 float 格式的数据，timestamp 的类型为 int64（从 1970-01-01 00:00:00 以来的毫秒数），timestamp 为可选默认为当前时间。具有相同 metric_name 的样本必须按照一个组的形式排列，并且每一行必须是唯一的指标名称和标签键值对组合。
