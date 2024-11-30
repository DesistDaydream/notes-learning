---
title: Variable
linkTitle: Variable
date: 2024-10-25T15:34
weight: 20
---

# 概述

> 参考：
>
> - [官方文档，仪表盘 - 变量](https://grafana.com/docs/grafana/latest/dashboards/variables/)

可以人为添加如下几类变量

- **Query** # 变量值从数据源查询中获取
- **Custom** # 手动定义一个变量的值
- **Text box** # 定义一个文本框变量，用户可以在其中输入任意字符串作为变量的值
- **Constant** # 定义一个隐藏常量变量，对于要共享的仪表板中的指标前缀很有用。TODO: 没看懂有啥用

# 全局内置变量

https://grafana.com/docs/grafana/latest/dashboards/variables/add-template-variables/#global-variables

## 时间相关
### `$__from` 与 `$__to`

https://grafana.com/docs/grafana/latest/dashboards/variables/add-template-variables/#__from-and-__to

Grafana 有两个内置时间范围变量：`$__from` 和 `$__to`。这两个变量的来源是 Grafana 的时间选择器，下图中的 From 与 To 选择的时间就是这两个变量的值。假如当前时间是 2024 年 11 月 24 日 0 点 0 分 0 秒，选择了 Last 6 hours 这个时间范围，则 `${__from}` 的值为 1732356000000（i.e. 2024-11-23 18:00:00）；`${__to}` 的值为 1732377600000（i.e. 2024-11-24 00:00:00）。

![time-picker_1](https://notes-learning.oss-cn-beijing.aliyuncs.com/grafana/time-picker_1.png)

可以通过如下语法控制显示出来的时间格式：

| Syntax                   | Example result           | Description                                                          |
| ------------------------ | ------------------------ | -------------------------------------------------------------------- |
| `${__from}`              | 1594671549254            | 默认格式。毫秒级 Unix 时间戳                                                    |
| `${__from:date}`         | 2020-07-13T20:19:09.254Z | No args, defaults to ISO 8601/RFC 3339                               |
| `${__from:date:iso}`     | 2020-07-13T20:19:09.254Z | ISO 8601/RFC 3339                                                    |
| `${__from:date:seconds}` | 1594671549               | Unix seconds epoch                                                   |
| `${__from:date:YYYY-MM}` | 2020-07                  | 使用 [date format](https://momentjs.com/docs/#/displaying/) 标准，自定义时间格式 |

### `$timeFilter` 或 `$__timeFilter`

https://grafana.com/docs/grafana/latest/dashboards/variables/add-template-variables/#timefilter-or-__timefilter

timeFilter 用以显示时间范围选择器选定的时间间隔。是一种类似把 from 与 to 两个变量合在一起的变量。

该变量主要应用在 SQL 的 WHERE 条件用，用以过滤特定时间范围内的数据。

<span style="background:rgba(255, 183, 139, 0.55)">TODO: 为什么在文本模式下渲染不出来？</span>

### `$__interval`

https://grafana.com/docs/grafana/latest/dashboards/variables/add-template-variables/#__interval

可以使用 `$__interval` 变量作为参数按时间（对于 InfluxDB、MySQL、Postgres、MSSQL）、日期直方图间隔（对于 Elasticsearch）进行分组，或作为汇总函数参数（对于 Graphite）。

步长，格式是 30s、1h、5d、etc. 

# 手动创建变量

> 参考：
>
> - [官方文档，模板与变量](https://grafana.com/docs/grafana/latest/variables/)
> - [Prometheus 天降奇兵文章](https://yunlzheng.gitbook.io/prometheus-book/part-ii-prometheus-jin-jie/grafana/templating)

在前面的小节中介绍了 Grafana 中 4 中常用的可视化面板的使用，通过在面板中使用 PromQL 表达式，Grafana 能够方便的将 Prometheus 返回的数据进行可视化展示。例如，在展示主机 CPU 使用率时，我们使用了如下表达式：

```promql
1 - (avg(irate(node_cpu{mode='idle'}[5m])) without (cpu))
```

该表达式会返回当前 Promthues 中存储的所有时间序列，每一台主机都会有一条单独的曲线用于体现其 CPU 使用率的变化情况：

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/kct3zl/1616067851128-b841ca2a-6a47-4d36-a92c-1b1b28f4000d.png)

而当用户只想关注其中某些主机时，基于当前我们已经学习到的知识只有两种方式，要么每次手动修改 Panel 中的 PromQL 表达式，要么直接为这些主机创建单独的 Panel。但是无论如何，这些硬编码方式都会直接导致 Dashboard 配置的频繁修改。在这一小节中我们将学习使用 Dashboard 变量的方式解决以上问题。

## Dashboard 变量

在 Grafana 中用户可以为 Dashboard 定义一组变量（Variables），变量一般包含一个到多个可选值。如下所示，Grafana 通过将变量渲染为一个下拉框选项，从而使用户可以动态的改变变量的值：

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/kct3zl/1616067851360-808effbf-a69f-4123-86d3-1a90d373fa18.png)

例如，这里定义了一个名为 node 的变量，用户可以通过在 PromQL 表达式或者 Panel 的标题中通过以下形式使用该变量：

```promql
1 - (avg(irate(node_cpu{mode='idle', instance=~"$node"}[5m])) without (cpu))
```

变量的值可以支持单选或者多选，当对接 Prometheus 时，Grafana 会自动将$node 的值格式化为如 “host1|host2|host3” 的形式。配合使用 PromQL 的标签正则匹配 “**=~**”，通过动态改变 PromQL 从而实现基于标签快速对时间序列进行过滤。

## 变量定义

通过 Dashboard 页面的 Settings 选项，可以进入 Dashboard 的配置页面并且选择 Variables 子菜单:

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/kct3zl/1616067851069-bf8ea850-c067-42af-b078-7fa1bbd154e9.png)
为 Dashboard 添加变量

用户需要指定变量的名称，后续用户就可以通过$variable_name 的形式引用该变量。Grafana 目前支持 6 种不同的变量类型，而能和 Prometheus 一起工作的主要包含以下 5 种类型：

| 类型         | 工作方式                                                                                 |
| ---------- | ------------------------------------------------------------------------------------ |
| Query      | 允许用户通过 Datasource 查询表达式的返回值动态生成变量的可选值                                                |
| Interval   | 该变量代表时间跨度，通过 Interval 类型的变量，可以动态改变 PromQL 区间向量表达式中的时间范围。如 rate(node_cpu\[$interval]) |
| Datasource | 允许用户动态切换当前 Dashboard 的数据源，特别适用于同一个 Dashboard 展示多个数据源数据的情况                            |
| Custom     | 用户直接通过手动的方式，定义变量的可选值                                                                 |
| Constant   | 常量，在导入 Dashboard 时，会要求用户设置该常量的值                                                      |

Label 属性用于指定界面中变量的显示名称，Hide 属性则用于指定在渲染界面时是否隐藏该变量的下拉框。

## 使用变量过滤时间序列

当 Prometheus 同时采集了多个主机节点的监控样本数据时，用户希望能够手动选择并查看其中特定主机的监控数据。这时我们需要使用 Query 类型的变量。

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/kct3zl/1616067851059-59709ed9-d567-4859-89ae-6344f4bf1f92.png)

新建 Query 类型的变量

如上所示，这里我们为 Dashboard 创建了一个名为 node 的变量，并且指定其类型为 Query。Query 类型的变量，允许用户指定数据源以及查询表达式，并通过正则匹配（Regex）的方式对查询结果进行处理，从而动态生成变量的可选值。在这里指定了数据源为 Prometheus，通过使用 node_load1 我们得到了两条时间序列：

```promql
node_load1{instance="foo:9100",job="node"}
node_load1{instance="localhost:9100",job="node"}
```

通过指定正则匹配表达式为`/.*instance="([^"]*).*/`从而匹配出标签 instance 的值作为 node 变量的所有可选项，即：

```yaml
foo:9100
localhost:9100
```

**Selection Options**选项中可以指定该变量的下拉框是否支持多选，以及是否包含全选（All）选项。

保存变量后，用户可以在 Panel 的 General 或者 Metrics 中通过$node 的方式使用该变量，如下所示：

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/kct3zl/1616067851042-b846b8f9-37c3-48b8-9e4b-127319beb2df.png)

在 Metrics 中使用变量

这里需要注意的是，如果允许用户多选在 PromQL 表达式中应该使用标签的正则匹配模式，因为 Grafana 会自动将多个选项格式化为如“foo:9100|localhost:9100”的形式。

使用 Query 类型的变量能够根据允许用户能够根据时间序列的特征维度对数据进行过滤。在定义 Query 类型变量时，除了使用 PromQL 查询时间序列以过滤标签的方式以外，Grafana 还提供了几个有用的函数：

- label_values(LABEL) # 返回 Promthues 所有监控指标中，标签名为 LABEL 的所有可选值
- label_values(METRIC, LABEL) # 返回监控指标 METRIC 中，标签名为 LABEL 的所有可选值
- metrics(RegEx) # 返回所有满足 RegEx(正则表达式) 匹配到的指标名称
- query_result(PromQL) # 返回 PromQL 的查询结果

例如，当需要监控 Prometheus 所有采集任务的状态时，可以使用如下方式，获取当前所有采集任务的名称：

```promql
label_values(up, job)
```

例如，有时候我们想要动态修改变量查询结果。比如某一个节点绑定了多个 ip，一个用于内网访问，一个用于外网访问，此时 prometheus 采集到的指标是内网的 ip，但我们需要的是外网 ip。这里我们想要能在 Grafana 中动态改变标签值，进行 ip 段的替换，而避免从 prometheus 或 exporter 中修改采集指标。

这时需要使用 grafana 的 query_result 函数

```promql
# 将10.10.15.xxx段的ip地址替换为10.20.15.xxx段 注：替换端口同理
query_result(label_replace(kube_pod_info{pod=~"$pod"}, "node", "10.20.15.$1", "node", "10.10.15.(.*)"))
```


```promql
# 通过正则从返回结果中匹配出所需要的ip地址
regex：/.*node="(.*?)".*/
```

在 grafana 中配置如图：

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/kct3zl/1616067851071-4ddd8696-d110-474a-8014-49824a3f7119.png)

## 使用变量动态创建 Panel 和 Row

当在一个 Panel 中展示多条时间序列数据时，通过使用变量可以轻松实现对时间序列的过滤，提高用户交互性。除此以外，我们还可以使用变量自动生成 Panel 或者 Row。 如下所示，当需要可视化当前系统中所有采集任务的监控任务运行状态时，由于 Prometheus 的采集任务配置可能随时发生变更，通过硬编码的形式实现，会导致 Dashboard 配置的频繁变更：

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/kct3zl/1616067851119-aa1bcd6d-8e23-4cd9-95be-3db2e40ce360.png)

Prometheus 采集任务状态

如下所示，这里为 Dashboard 定义了一遍名为 job 的变量：

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/kct3zl/1616067851059-05a9e06b-9215-4768-9126-542c61e9d70f.png)

使用变量获取当前所有可选任务

通过使用 label_values 函数，获取到当前 Promthues 监控指标 up 中所有可选的 job 标签的值：

```promql
label_values(up, job)
```

如果变量启用了 Multi-value 或者 Include All Option 选项的变量，那么在 Panel 的 General 选项的 Repeat 中可以选择自动迭代的变量，这里使用了 Singlestat 展示所有监控采集任务的状态：

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/kct3zl/1616067851091-b6fcaadd-b2aa-41c1-957e-039bc878dbc6.png)

General 中的 Repeat 选项

Repeat 选项设置完成后，Grafana 会根据当前用户的选择，自动创建一个到多个 Panel 实例。 为了能够使 Singlestat Panel 能够展示正确的数据，如下所示，在 Prometheus 中，我们依然使用了`$job`变量，不过此时的 `$job` 反应的是当前迭代的值：

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/kct3zl/1616067851115-f60eb991-5dbb-4c33-8cd2-58639c7001c6.png)
在 Metric 中使用变量

而如果还希望能够自动生成 Row，只需要在 Row 的设置中，选择需要 Repeat 的变量即可：

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/kct3zl/1616067851130-ac817567-fee3-40c3-b4e1-d5f07f86cf8d.png)

