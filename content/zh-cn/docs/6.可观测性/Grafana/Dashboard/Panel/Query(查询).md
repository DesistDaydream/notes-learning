---
title: Query(查询)
---

# 概述

> 参考：
>
> - [官方文档，面板 - 查询](https://grafana.com/docs/grafana/latest/panels/queries)
>   - 新 URL https://grafana.com/docs/grafana/latest/panels-visualizations/query-transform-data/

Query 标签由如下几个元素组成

- Data source selector(数据源选择器)
- Query options(查询选项)
- Query inspector button()
- Query editor list
- Expressions

# Data source(数据源) # 数据源选择器

![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/grafana/dashboard/1636275083307-1fa893ed-814b-434d-9e51-e2c6499c6f45.png)

这部分是一个下拉列表，通过这里选择要使用的数据源，Query editor 中的查询语句，必须要是此数据源可以可以识别的。

# Query options(查询选项)

https://grafana.com/docs/grafana/latest/panels-visualizations/query-transform-data/#query-options

## Max data points(最大数据点)

## Min interval(最小间隔)

## Interval(间隔)

Interval 是一个，通过时间来聚合或分组的一些数据点时，使用的 **time span(时间跨度)**。该选项也可以实现查询编辑器中的 Min step 和 Resolution 类似的效果。但是，更多的是用在 Prometheus 范围向量查询语句中，比如 `rate(http_requests_total[$__interval])`。interval 选项可以为两个 Grafana 的内置变量提供值，`$__interval` 与 `$__interval_ms`。

也就是说，Interval 常用来计算 总和、平均值、速率 等一段时间范围的变化量。

除了在这里可以定义 Interval，还可以通过 [Variable](/docs/6.可观测性/Grafana/Dashboard/Variable.md) 定义。

## Relative time(相对时间)

为 Panel 配置独立时间范围选择器。利用 Relative time 的配置，可以为 Panel 配置独立于 Dashboard 时间选择器（i.e. 这可以让配了 Relative time 的 Panel 与 Dashboard 右上角的**时间范围选择器**中选择的时间不同）

下面的示例中，将 Relative time 设为 7d，此时可以看到 Dashboard 选择的是 30min，但是这个 Panel 展示的是 7d 的数据。

![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/grafana/dashboard/202408021344061.png)

## Time shift(时移)

# Query inspector(查询检查器)

用于调试查询编辑器中的查询语句，里面会显示 HTTP 的请求和响应的原始数据

# Query editor(查询编辑器)

查询编辑器可以编写查询语句，以便从数据源中获取数据，不同的数据源，其查询编辑器也各不相同。

### Legend(图例) - 改变 series 的名称

Legend 可以通过一种模式来改变 series 的名称。这个模式类似于 Go 模板的语法，使用 `{{ }}` 符号，引用 series 名称中，标签的值。

例如，上图我在 Legend 框内填 `设备：{{instance}}` 那么，将会出现这种效果。

![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/grafana/dashboard/1636266519763-257b682f-43aa-42aa-ae64-4dc0dd9df523.png)

> 这里的`{{instance}}` 符号会获取 instance 这个标签的值

### Min step - 最小步长

可以控制 Prometheus 查询数据时的最小步长(Step)，从而减少或增加从 Prometheus 返回的数据量。

- 步长指起始时间与结束时间之间获取的所有数据的间隔时间。假如 step=10 则每隔 10 秒获取一次样本值。
  - 就好像人走路，一步迈多长，这里就是返回样本值时，每隔多久返回一次。
- 比如，我想要查询 14 点到 15 点之间的数据，假如 Min step 为 1m，则一共返回 60 个样本。假如 Min step 为 30m，则一共返回 2 个样本，效果如下

![image.png|550](https://notes-learning.oss-cn-beijing.aliyuncs.com/grafana/dashboard/1636273461647-2dbebd5e-b5ca-47d8-a23e-07d4e0d72ebe.png) ![image.png|550](https://notes-learning.oss-cn-beijing.aliyuncs.com/grafana/dashboard/1636273507618-65072c3a-f3bd-46c9-aa79-979c6bf388b8.png)

### Resolution - 分辨率

则可以控制 Grafana 自身渲染的数据量。例如，如果**Resolution**的值为**1/10**，Grafana 会将 Prometeus 返回的 10 个样本数据合并成一个点。因此 **Resolution**越小可视化的精确性越高，反之，可视化的精度越低。

### Format - 格式化获取到的样本数据

- Time series # 时间序列格式。默认格式
- Table # 表格式。用于 Table 面板
- Heatmap # 热力图格式。用于 Heatmap 面板

### Instant - 瞬时。控制是否获取指标的瞬时值

开启后，只会显示最近一次的序列的值。常用于 Stat、Gauge 这种面板，以及 Graph 面板下 Series 模式的 X 轴。

因为开启瞬时值，只会显示当前值，是没有时间的概念的。
