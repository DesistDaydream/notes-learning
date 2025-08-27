---
title: PromQL
linkTitle: PromQL
weight: 1
---

# 概述

> 参考：
>
> - [官方文档，Prometheus - 查询 - 基础](https://prometheus.io/docs/prometheus/latest/querying/basics/)

**Prometheus Query Language(Prometheus 查询语言，简称 PromQL)** 是一种提供了查询功能的编程语言，用来实时选择和汇总 [时间序列数据](/docs/5.数据存储/数据库/时间序列数据/时间序列数据.md)。通过 PromQL 可以对监控数据进行筛选、过滤、组合等等操作。使用 PromQL 编写的语句也可以称为 **Expression(表达式)**，表达式的结果可以通过其他方式显示为图形。

## PromQL 体验

在 graph 页面，可以在红框位置输入表达式

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/lfubxg/1616069073146-722bd1dd-2647-4bf9-84ed-4c2cbf694785.png)

点击红框内的对话框 并输入关键字，系统会自动弹出可用的 metrics name

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/lfubxg/1616069073181-bbf09f64-c80f-4e9f-9644-2916ce031358.jpeg)

表达式直接使用 MetricsName，则展示此时此刻的以 node_cpu_seconds_total 为指标名的所有 TimeSeries(时间序列) 数据

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/lfubxg/1616069073159-7b2a2e59-5866-478c-a87e-a9c6e42a65d5.png)

如果需要筛选则可以输入如下图实例的表达式：node_cpu_seconds_total{job=~"external.\*"}

筛选出来 job 名开头是 external 的 cpu 情况。允许使用正则表达式，=~表示的就是用过正则来匹配后面的值

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/lfubxg/1616069073186-857243c7-7acb-4a2a-9787-ecd8dabbf732.jpeg)

Prometheus 通过 MetricsName(指标名称) 及其对应的一组 LabelSet(标签集) 定义唯一的一条时间序列。指标名称反映了监控样本的基本标识，而 label 则在这个基本特征上为采集到的数据提供了多种特征维度。用户可以基于这些特征维度过滤，聚合，统计从而产生新的计算后的一条时间序列。

PromQL 是 Prometheus 内置的数据查询语言，其提供对时间序列数据丰富的查询，聚合以及逻辑运算能力的支持。并且被广泛应用在 Prometheus 的日常应用当中，包括对数据查询、可视化、告警处理当中。可以这么说，PromQL 是 Prometheus 所有应用场景的基础，理解和掌握 PromQL 是 Prometheus 入门的第一课。

# PromQL 基本语法

PromQL 没有绝对通用的语法，在不同场景查询条件下，语法也不同。但是，语法必须要有语句，这种语句就称为 **Expression(表达式)**，Expression 可以是简单的字符串，也可以是一个指标名称，甚至是一串基于指标的复杂语法。

在 PromQL 中，任何 **Expression(表达式)** 或者 **subExpression(子表达式)** 都可以归为四种类型：

- **Instant Vector Selectors(即时向量选择器)** # 包含每个时间序列的单个样本的一组时间序列，共享相同的时间戳。
  - **即时向量** 在有的地方也被翻译为 **瞬时向量**，都是同一个意思。
- **Range Vector Selectors(范围向量选择器)** # 包含每个时间序列随时间变化的数据点的一组时间序列。
  - **范围向量**
- **String(字符串)** # 一个简单的字符串值(目前未被使用)
- **Scalar(标量)** # 一个简单的数字浮点值

在这四种表达式中，我们还可以通过 **Operators(操作符)** 和 **Functions(函数)** 加工获取到的时间序列数据。

这四种类型，又可以进行统一分类

- Instant Vector Selectors 和 Range Vector Selectors 统称为 **TimeSeries Selectors(时间序列选择器)**
  - 这种表达式会根据 Metrics 来获取指定的时间序列。
- String 和 Scalar 统称为 [**Literal(字面量)**](/docs/2.编程/计算机科学/Literal.md)
  - 给定不同类型的 Literal，就返回对应类型的的值，Prom 里只支持 string 和 scalar 这两种类型

# Expression(表达式)

## Instant Vector Selectors(即时向量选择器)

即时向量选择器就是获取通过 Metric 来获取当前最新时间戳上的值，这就是即时的含义(指的是当前最新的值)

### 根据 MetricName 选择时间序列数据

当我们直接使用 Metrics 查询时，可以查询该 Metrics 下的所有时间序列数据。如：

**promhttp_metric_handler_requests_total** # 该表达式会返回指标名称为 promhttp_metric_handler_requests_total 的所有时间序列数据：

```promql
promhttp_metric_handler_requests_total{code="200", instance="172.38.40.250:9090", job="prometheus"}=(98@1518096812.326)
promhttp_metric_handler_requests_total{code="500", instance="172.38.40.250:9090", job="prometheus"}=(0@1518096812.326)
promhttp_metric_handler_requests_total{code="503", instance="172.38.40.250:9090", job="prometheus"}=(0@1518096812.326)
```

### 过滤选择到的时间序列数据

如果想要过滤查询结果，可以根据 metric 中的 Label 来进行匹配，通过在 `{ }` 符号中使用一组 Label 来进一步过滤这些时间序列数据，支持两种匹配模式：完全匹配和正则匹配。可以使用 4 种用于标签匹配的操作符

- `=` # 完全匹配。通过使用 `Label="Value"` 可以选择那些标签值满足表达式定义的时间序列；
- `!=` # 完全不匹配。反之使用 `Label!="Value"` 则可以根据标签值匹配排除时间序列；
- `=~` # 正则匹配。使用 `Label=~"RegEx"` 表示选择那些标签符合正则表达式定义的时间序列；
- `!~` # 正则不匹配。反之使用 `Label!~"RegEx"` 进行排除；
  - 正则匹配中，多个表达式之间使用 `|` 进行分隔：

例如，如果我们只需要查询所有 promhttp_metric_handler_requests_total 指标中满足标签 code 的值为 200 的时间序列，则可以使用如下表达式：

promhttp_metric_handler_requests_total{code="200"}

反之使用 code!="200" 则可以排除这些时间序列：

promhttp_metric_handler_requests_total{code!="200"}

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/lfubxg/1616069073195-72aba6c2-dabb-47d8-960b-b26cdf73e04a.png)

如果想查询多个环节下的时间序列序列可以使用如下表达式：

promhttp_metric_handler_requests_total{code=~"200|500"}

### 根据 Label 选择时间序列数据

我们可以把 Label 当作表达式来获取所有具有这些 Label 的 Metrics 的时间序列数据

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/lfubxg/1616069073324-a1b41ace-fb9a-4f25-af64-b18536239692.png)

在 Label 中还可以使用原始标签

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/lfubxg/1616069073217-6f22cfb6-c05c-44ef-b079-6d91f6e295a2.png)

## Range Vector Selectors(范围向量选择器)

直接通过表达式 http_requests_total 查询时间序列时，返回值中只会包含该时间序列中的最新的一个样本值，这样的返回结果我们称之为瞬时向量。而相应的这样的表达式称之为**瞬时向量表达式**。

而如果我们想要获取过去一段时间范围内的样本数据时，我们则需要使用**范围向量表达式**。范围向量表达式 和 瞬时向量表达式 基本一样，唯一的区别在于，范围向量表达式 中我们需要定义时间选择的范围，时间范围通过 `[ ]` (这个符号表示：**时间范围选择器**) 进行定义。

例如：

`promhttp_metric_handler_requests_total{}[4m]` # 该表达式将会返回查询到的时间序列中最近 4 分钟的所有样本数据：数据如下

> [!Notes]
>
> - 这种查询方式只可以获取到值，并不能生成图表，因为图表中的一条向量是由很多个点组成的，每个点的位置由横轴是 time、纵轴是 value 互相确认得到的。但是这种查询结果会使每个点由多个 time 与多个 value 组成，一个点怎么可能由多个点合成一个呢？这在二维图标上是没法显示出来的。
> - 由于有多个值而且还没法展示，所以范围向量一般不会单独使用，而是与 irate()等函数一起使用。以便让多个值根据指定的函数规则聚合成唯一的一个值。

```promql
promhttp_metric_handler_requests_total{code="200", instance="172.38.40.250:9090", job="prometheus"}=[
    223@1518096812.326
    224@1518096817.326
    225@1518096822.326
    226@1518096827.326
]
promhttp_metric_handler_requests_total{code="200", instance="172.38.40.250:9090", job="prometheus"}=[
    0@1518096812.326
    0@1518096817.326
    0@1518096822.326
    0@1518096827.326
]
```

除了使用 m 表示分钟以外，PromQL 的时间范围选择器支持其它时间单位：

- ms # 毫秒
- s # 秒
- m # 分钟
- h # 小时
- d # 天
- w # 周
- y # 年

### Subquery(子查询)

在 `[ ]` 符号表示的 时间范围选择器 中，可以使用 **Subquery(子查询)** 功能为时间范围添加 **Resolution(分辨率)**（有时候也成为 **Step(步长)**），比如：

`promhttp_metric_handler_requests_total{}[4m:30s]`

30s 就是 Resolution，表示在 4m 的时间范围中，每隔 30 秒取一个样本值。

Resolution 通常是可省略的，默认值为 [Promethesu Server](/docs/6.可观测性/Metrics/Prometheus/Configuration/Promethesu%20Server.md) 的 `.global.evaluation_interval`

子查询之所以叫子查询，通常用在多个 [PromQL Functions](/docs/6.可观测性/Metrics/Prometheus/PromQL/PromQL%20Functions.md) 的场景，比如

`rate(avg_over_time(node_network_receive_bytes_total{device="eth0"}[5m])[6h:5m])`是合法的

`rate(avg_over_time(node_network_receive_bytes_total{device="eth0"}[5m])[6h:])` 是合法的

`rate(avg_over_time(node_network_receive_bytes_total{device="eth0"}[5m])[6h])` 是非法的

> [!Attention]
> 其中 `avg_over_time(node_network_receive_bytes_total{device="eth0"}[5m]` 的结果为即时向量的时间序列，若想基于该结果使用 rate 函数，则必须<font color="#ff0000">显式</font>使用 Subquery 指定 Resolution。像上面两个合法的例子示例

### Offset modifier(位移修饰符) - 时间位移操作

在瞬时向量表达式或者范围向量表达式中，都是以当前时间为基准：

`http_request_total{}` # 瞬时向量表达式，选择当前最新的数据

`http_request_total{}[5m]` # 范围向量表达式，选择以当前时间为基准，5 分钟内的数据

而如果我们想查询，5 分钟前的瞬时样本数据，或昨天一天的区间内的样本数据呢? 这个时候我们就可以使用位移操作，位移操作的关键字为 offset。

可以使用 offset 时间位移操作：

`http_request_total{} offset 5m` #

`http_request_total{}[5m] offset 1d` #

## String(字符串) 表达式

直接使用字符串，作为 PromQL 表达式，则会直接返回字符串。

```text
"this is a string"
'these are unescaped: \n \\ \t'
`these are not unescaped: \n ' " \t`
```

比如下图：

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/lfubxg/1616069073187-d54f3195-edf1-4215-af6d-97b1132832f6.png)

如果通过 API，则是这种样子：

```bash
~]# curl "http://172.38.40.244:30001/api/v1/query?query=%22HelloWorld%22"
{"status":"success","data":{"resultType":"string","result":[1610673038.083,"HelloWorld"]}}
```

## Scalar(标量) 表达式

除了使用瞬时向量表达式和区间向量表达式以外，PromQL 还直接支持用户使用标量(Scalar)和字符串(String)。

标量只有一个数字，没有时序。

例如：

10

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/lfubxg/1616069073229-d7297a02-c472-4d5d-a21f-0265ecc5ca63.png)

需要注意的是，当使用表达式 count(http_requests_total)，返回的数据类型，依然是瞬时向量。用户可以通过内置函数 scalar()将单个瞬时向量转换为标量。

## Operators(运算符)

Prometheus 支持许多二进制和聚合运算符。详见《[PromQL Operators](/docs/6.可观测性/Metrics/Prometheus/PromQL/PromQL%20Operators(运算符).md)》章节

## Functions(函数)

Prometheus 支持多种对数据进行操作的函数。详见《[PromQL Functions](/docs/6.可观测性/Metrics/Prometheus/PromQL/PromQL%20Functions.md)》章节

# 表达式样例

所有的 PromQL 表达式都必须至少包含一个指标名称(例如 http_request_total)，或者一个不会匹配到空字符串的标签过滤器(例如{code="200"})。

因此以下两种方式，均为合法的表达式：
合法的表达式：

- http_request_total
- http_request_total{}
- {method="get"}

不合法的表达式：

- `{job=~".*"}`
  - 注意，可以使用 `'{job=~"..*"}'` 来匹配所有 job 的 metric，但是官方不建议这么用，防止意外情况发生

同时，除了使用{label=value}的形式以外，我们还可以使用内置的\_\_name\_\_标签来指定监控指标名称：

`{__name__=~"http_request_total"}` # 合法

`{__name__=~"node_disk_bytes_read|node_disk_bytes_written"}` # 合法

## 即时向量查询结果的示意图

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/lfubxg/1616069073170-822867a1-8a50-4fa6-ae70-e34c802af202.jpeg)

## 范围向量查询结果示意图

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/lfubxg/1616069073260-b6507c67-3c3c-4c52-bf34-773772f11743.jpeg)

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/lfubxg/1616069073175-af75032f-496c-4fe2-bfbf-4bd44d56e836.jpeg)
