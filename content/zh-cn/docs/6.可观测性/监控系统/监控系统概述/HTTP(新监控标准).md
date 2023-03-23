---
title: HTTP(新监控标准)
---

# 概述

由于 SNMP 的种种不便，现在更多的是基于 HTTP 协议来实现监控指标的采集。

同样，也是需要一个 Client 采集指标，需要一个 Server 端接收指标后存储指标。

就像 SNMP 协议一样，光有协议还不行，基于 HTTP 协议的监控也需要一个数据模型的标准，就像 MIB 和 OID 类似。而现阶段，从 Prometheus 演化过来的 OpenMetrics 标准，就是这么一种东西。

# OpenMetrics

> 参考：
>
> - [GitHub 项目，OpenObservability/OpenMetrics](https://github.com/OpenObservability/OpenMetrics)
> - [官网](https://openmetrics.io/)
> - [OpenMetrics 规范](https://github.com/OpenObservability/OpenMetrics/blob/main/specification/OpenMetrics.md)

**OpenMetrics** 是新时代的监控指标的标准，由 CNCF 主导，OpenMetrics 定义了大规模传输云原生指标的事实标准。

- **OpenMetricsSpec** 用来定义监控指标的标准

# Data Model(数据模型)

平时我们口语交流，一般将随时间变化的数据称为 **Metrics(指标)**。这是监控数据的另一种叫法，与 OID 类似，可以代表一种监控数据、同时也是一种名词，比如我采集某个程序的监控数据，可以说采集这个程序的 Metrics。所以 Metrics 是一个抽象的叫法。

This section MUST be read together with the ABNF section. In case of disagreements between the two, the ABNF's restrictions MUST take precedence. This reduces repetition as the text wire format MUST be supported.

## Data Types

### Values

Metric values in OpenMetrics MUST be either floating points or integers. Note that ingestors of the format MAY only support float64. The non-real values NaN, +Inf and -Inf MUST be supported. NaN MUST NOT be considered a missing value, but it MAY be used to signal a division by zero.

#### Booleans

Boolean values MUST follow 1==true, 0==false.

### Timestamps

Timestamps MUST be Unix Epoch in seconds. Negative timestamps MAY be used.

### Strings

Strings MUST only consist of valid UTF-8 characters and MAY be zero length. NULL (ASCII 0x0) MUST be supported.

### Label

Labels are key-value pairs consisting of strings.
Label names beginning with underscores are RESERVED and MUST NOT be used unless specified by this standard. Label names MUST follow the restrictions in the ABNF section.
Empty label values SHOULD be treated as if the label was not present.

### LabelSet

A LabelSet MUST consist of Labels and MAY be empty. Label names MUST be unique within a LabelSet.

### MetricPoint

Each MetricPoint consists of a set of values, depending on the MetricFamily type.

### Exemplars

Exemplars are references to data outside of the MetricSet. A common use case are IDs of program traces.
Exemplars MUST consist of a LabelSet and a value, and MAY have a timestamp. They MAY each be different from the MetricPoints' LabelSet and timestamp.
The combined length of the label names and values of an Exemplar's LabelSet MUST NOT exceed 128 UTF-8 characters. Other characters in the text rendering of an exemplar such as ",= are not included in this limit for implementation simplicity and for consistency between the text and proto formats.
Ingestors MAY discard exemplars.

### Metric

Metrics are defined by a unique LabelSet within a MetricFamily. Metrics MUST contain a list of one or more MetricPoints. Metrics with the same name for a given MetricFamily SHOULD have the same set of label names in their LabelSet.
MetricPoints SHOULD NOT have explicit timestamps.
If more than one MetricPoint is exposed for a Metric, then its MetricPoints MUST have monotonically increasing timestamps.

### MetricFamily

A MetricFamily MAY have zero or more Metrics. A MetricFamily MUST have a name, HELP, TYPE, and UNIT metadata. Every Metric within a MetricFamily MUST have a unique LabelSet.

#### Name

MetricFamily names are a string and MUST be unique within a MetricSet. Names SHOULD be in snake_case. Metric names MUST follow the restrictions in the ABNF section.
Colons in MetricFamily names are RESERVED to signal that the MetricFamily is the result of a calculation or aggregation of a general purpose monitoring system.
MetricFamily names beginning with underscores are RESERVED and MUST NOT be used unless specified by this standard.

##### Suffixes

The name of a MetricFamily MUST NOT result in a potential clash for sample metric names as per the ABNF with another MetricFamily in the Text Format within a MetricSet. An example would be a gauge called "foo_created" as a counter called "foo" could create a "foo_created" in the text format.
Exposers SHOULD avoid names that could be confused with the suffixes that text format sample metric names use.

- Suffixes for the respective types are:
- Counter: '\_total', '\_created'
- Summary: '\_count', '\_sum', '\_created', '' (empty)
- Histogram: '\_count', '\_sum', '\_bucket', '\_created'
- GaugeHistogram: '\_gcount', '\_gsum', '\_bucket'
- Info: '\_info'
- Gauge: '' (empty)
- StateSet: '' (empty)
- Unknown: '' (empty)

#### Type

Type specifies the MetricFamily type. Valid values are "unknown", "gauge", "counter", "stateset", "info", "histogram", "gaugehistogram", and "summary".

#### Unit

Unit specifies MetricFamily units. If non-empty, it MUST be a suffix of the MetricFamily name separated by an underscore. Be aware that further generation rules might make it an infix in the text format.

#### Help

Help is a string and SHOULD be non-empty. It is used to give a brief description of the MetricFamily for human consumption and SHOULD be short enough to be used as a tooltip.

#### MetricSet

A MetricSet is the top level object exposed by OpenMetrics. It MUST consist of MetricFamilies and MAY be empty.
Each MetricFamily name MUST be unique. The same label name and value SHOULD NOT appear on every Metric within a MetricSet.
There is no specific ordering of MetricFamilies required within a MetricSet. An exposer MAY make an exposition easier to read for humans, for example sort alphabetically if the performance tradeoff makes sense.
If present, an Info MetricFamily called "target" per the "Supporting target metadata in both push-based and pull-based systems" section below SHOULD be first.

## Metric Types

### Gauge

Gauges are current measurements, such as bytes of memory currently used or the number of items in a queue. For gauges the absolute value is what is of interest to a user.
A MetricPoint in a Metric with the type gauge MUST have a single value.
Gauges MAY increase, decrease, or stay constant over time. Even if they only ever go in one direction, they might still be gauges and not counters. The size of a log file would usually only increase, a resource might decrease, and the limit of a queue size may be constant.
A gauge MAY be used to encode an enum where the enum has many states and changes over time, it is the most efficient but least user friendly.

### Counter

Counters measure discrete events. Common examples are the number of HTTP requests received, CPU seconds spent, or bytes sent. For counters how quickly they are increasing over time is what is of interest to a user.
A MetricPoint in a Metric with the type Counter MUST have one value called Total. A Total is a non-NaN and MUST be monotonically non-decreasing over time, starting from 0.
A MetricPoint in a Metric with the type Counter SHOULD have a Timestamp value called Created. This can help ingestors discern between new metrics and long-running ones it did not see before.
A MetricPoint in a Metric's Counter's Total MAY reset to 0. If present, the corresponding Created time MUST also be set to the timestamp of the reset.
A MetricPoint in a Metric's Counter's Total MAY have an exemplar.

### StateSet

StateSets represent a series of related boolean values, also called a bitset. If ENUMs need to be encoded this MAY be done via StateSet.
A point of a StateSet metric MAY contain multiple states and MUST contain one boolean per State. States have a name which are Strings.
A StateSet Metric's LabelSet MUST NOT have a label name which is the same as the name of its MetricFamily.
If encoded as a StateSet, ENUMs MUST have exactly one Boolean which is true within a MetricPoint.
This is suitable where the enum value changes over time, and the number of States isn't much more than a handful.
EDITOR’S NOTE: This might be better as Consideration
MetricFamilies of type StateSets MUST have an empty Unit string.

### Info

Info metrics are used to expose textual information which SHOULD NOT change during process lifetime. Common examples are an application's version, revision control commit, and the version of a compiler.
A MetricPoint of an Info Metric contains a LabelSet. An Info MetricPoint's LabelSet MUST NOT have a label name which is the same as the name of a label of the LabelSet of its Metric.
Info MAY be used to encode ENUMs whose values do not change over time, such as the type of a network interface.
MetricFamilies of type Info MUST have an empty Unit string.

### Histogram

Histograms measure distributions of discrete events. Common examples are the latency of HTTP requests, function runtimes, or I/O request sizes.
A Histogram MetricPoint MUST contain at least one bucket, and SHOULD contain Sum, and Created values. Every bucket MUST have a threshold and a value.
Histogram MetricPoints MUST have at least a bucket with an +Inf threshold. Buckets MUST be cumulative. As an example for a metric representing request latency in seconds its values for buckets with thresholds 1, 2, 3, and +Inf MUST follow value_1 <= value_2 <= value_3 <= value\_+Inf. If ten requests took 1 second each, the values of the 1, 2, 3, and +Inf buckets MUST equal 10.
The +Inf bucket counts all requests. If present, the Sum value MUST equal the Sum of all the measured event values. Bucket thresholds within a MetricPoint MUST be unique.
Semantically, Sum, and buckets values are counters so MUST NOT be NaN or negative. Negative threshold buckets MAY be used, but then the Histogram MetricPoint MUST NOT contain a sum value as it would no longer be a counter semantically. Bucket thresholds MUST NOT equal NaN. Count and bucket values MUST be integers.
A Histogram MetricPoint SHOULD have a Timestamp value called Created. This can help ingestors discern between new metrics and long-running ones it did not see before.
A Histogram's Metric's LabelSet MUST NOT have a "le" label name.
Bucket values MAY have exemplars. Buckets are cumulative to allow monitoring systems to drop any non-+Inf bucket for performance/anti-denial-of-service reasons in a way that loses granularity but is still a valid Histogram.
EDITOR’S NOTE: The second sentence is a consideration, it can be moved if needed
Each bucket covers the values less and or equal to it, and the value of the exemplar MUST be within this range. Exemplars SHOULD be put into the bucket with the highest value. A bucket MUST NOT have more than one exemplar.

### GaugeHistogram

GaugeHistograms measure current distributions. Common examples are how long items have been waiting in a queue, or size of the requests in a queue.
A GaugeHistogram MetricPoint MUST have at least one bucket with an +Inf threshold, and SHOULD contain a Gsum value. Every bucket MUST have a threshold and a value.
The buckets for a GaugeHistogram follow all the same rules as for a Histogram.
The bucket and Gsum of a GaugeHistogram are conceptually gauges, however bucket values MUST NOT be negative or NaN. If negative threshold buckets are present, then sum MAY be negative. Gsum MUST NOT be NaN. Bucket values MUST be integers.
A GaugeHistogram's Metric's LabelSet MUST NOT have a "le" label name.
Bucket values can have exemplars.
Each bucket covers the values less and or equal to it, and the value of the exemplar MUST be within this range. Exemplars SHOULD be put into the bucket with the highest value. A bucket MUST NOT have more than one exemplar.

### Summary

Summaries also measure distributions of discrete events and MAY be used when Histograms are too expensive and/or an average event size is sufficient.
They MAY also be used for backwards compatibility, because some existing instrumentation libraries expose precomputed quantiles and do not support Histograms. Precomputed quantiles SHOULD NOT be used, because quantiles are not aggregatable and the user often can not deduce what timeframe they cover.
A Summary MetricPoint MAY consist of a Count, Sum, Created, and a set of quantiles.
Semantically, Count and Sum values are counters so MUST NOT be NaN or negative. Count MUST be an integer.
A MetricPoint in a Metric with the type Summary which contains Count or Sum values SHOULD have a Timestamp value called Created. This can help ingestors discern between new metrics and long-running ones it did not see before. Created MUST NOT relate to the collection period of quantile values.
Quantiles are a map from a quantile to a value. An example is a quantile 0.95 with value 0.2 in a metric called myapp_http_request_duration_seconds which means that the 95th percentile latency is 200ms over an unknown timeframe. If there are no events in the relevant timeframe, the value for a quantile MUST be NaN. A Quantile's Metric's LabelSet MUST NOT have "quantile" label name. Quantiles MUST be between 0 and 1 inclusive. Quantile values MUST NOT be negative. Quantile values SHOULD represent the recent values. Commonly this would be over the last 5-10 minutes.

### Unknown

Unknown SHOULD NOT be used. Unknown MAY be used when it is impossible to determine the types of individual metrics from 3rd party systems.
A point in a metric with the unknown type MUST have a single value.

## 基本示例

Metrics 数据格式如下图所示
\# HELP http_requests_total The total number of HTTP requests.
\# TYPE http_requests_total counter
http_requests_total{method="post",code="200"} 1027 1395066363000
http_requests_total{method="post",code="400"} 3 1395066363000
\# Escaping in label values:msdos_file_access_time_seconds{path="C:\DIR\FILE.TXT",error="Cannot find file:\n"FILE.TXT""} 1.458255915e9
\# Minimalistic line:metric_without_timestamp_and_labels 12.47
\# A weird metric from before the epoch:something_weird{problem="division by zero"} +Inf -3982045

默认有三行数据来表示

1. \#HELP MetricsName Metrics 的描述
2. \#TYPE MetricsName Metrics 的数据类型
3. MetricsName 与 Metrics 的值
4. 如果有多个 Metrics 的项目，则会有多行

主要由三个部分组成：样本的一般注释信息（HELP），样本的类型注释信息（TYPE）和样本。Prometheus 会对 Exporter 响应的内容逐行解析：
如果当前行以# HELP 开始，Prometheus 将会按照以下规则对内容进行解析，得到当前的指标名称以及相应的说明信息：
\# HELP

如果当前行以# TYPE 开始，Prometheus 会按照以下规则对内容进行解析，得到当前的指标名称以及指标类型:
\# TYPE

TYPE 注释行必须出现在指标的第一个样本之前。如果没有明确的指标类型需要返回为 untyped。 除了# 开头的所有行都会被视为是监控样本数据。 每一行样本需要满足以下格式规范:
metric_name \[
 "{" label_name "=" `"` label_value `"` { "," label_name "=" `"` label_value `"` } \[ "," ] "}"
] value \[ timestamp ]

其中 metric_name 和 label_name 必须遵循 PromQL 的格式规范要求。value 是一个 float 格式的数据，timestamp 的类型为 int64（从 1970-01-01 00:00:00 以来的毫秒数），timestamp 为可选默认为当前时间。具有相同 metric_name 的样本必须按照一个组的形式排列，并且每一行必须是唯一的指标名称和标签键值对组合。
