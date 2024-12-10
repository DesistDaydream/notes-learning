---
title: Metric Queries
linkTitle: Metric Queries
weight: 20
---

# 概述

> 参考：
>
> - date: "2024-12-10T13:54"
> - [官方文档，查询 - 指标查询](https://grafana.com/docs/loki/latest/query/metric_queries/)

LogQL 还可以通过 **Functions(函数)** 来对每个日志流进行计算以实现 **Metric Queries(指标查询)** 。就是将日志流中的日志计数，并根据技术进行一些查询，这种查询方式与 [PromQL](docs/6.可观测性/Metrics/Prometheus/PromQL/PromQL.md) 的指标查询类似。

指标查询可用于计算诸如 错误消息率、最近 3 个小时内日志数量最多的 N 个日志源、etc. 的信息。

## 范围向量聚合

LogQL 与 Prometheus 具有相同的范围向量概念，不同之处在于所选的样本范围包括每个日志

常用函数主要是如下 4 个：

- `rate`: 计算每秒的日志条目
- `count_over_time`: 对指定范围内的每个日志流的条目进行计数
- `bytes_rate`: 计算日志流每秒的字节数
- `bytes_over_time`: 对指定范围内的每个日志流的使用的字节数

比如计算 nginx 的 qps：

```logql
rate({filename="/var/log/nginx/access.log"}[5m]))
```

计算 kernel 过去 5 分钟发生 oom 的次数：

```logql
count_over_time({filename="/var/log/message"} |~ "oom_kill_process" [5m]))
```

## 聚合运算符

与 PromQL 一样，LogQL 支持内置聚合运算符的一个子集，可用于聚合单个向量的元素，从而产生一个具有更少元素但具有集合值的新向量：

- sum: Calculate sum over labels
- min: Select minimum over labels
- max: Select maximum over labels
- avg: Calculate the average over labels
- stddev: Calculate the population standard deviation over labels
- stdvar: Calculate the population standard variance over labels
- count: Count number of elements in the vector
- bottomk: Select smallest k elements by sample value
- topk: Select largest k elements by sample value
- sum：求和
- min：最小值
- max：最大值
- avg：平均值
- stddev：标准差
- stdvar：标准方差
- count：计数
- bottomk：最小的 k 个元素
- topk：最大的 k 个元素

聚合函数我们可以用如下表达式描述：

```logql
<aggr-op>([parameter,] <vector expression) [without|by (label list)]
```

对于需要对标签进行分组时，我们可以用 `without` 或者 `by` 来区分。比如计算 nginx 的 qps，并按照 pod 来分组：

```logql
sum(rate({filename="/var/log/nginx/access.log"}[5m])) by (pod)
```

只有在使用 `bottomk` 和 `topk` 函数时，我们可以对函数输入相关的参数。比如计算 nginx 的 qps 最大的前 5 个，并按照 pod 来分组：

```logql
topk(5,sum(rate({filename="/var/log/nginx/access.log"}[5m])) by (pod)))
```

## Binary Operators(二元运算符)

### 数学计算

Loki 存的是日志，都是文本，怎么计算呢？显然 LogQL 中的数学运算是面向区间向量操作的，LogQL 中的支持的二进制运算符如下：

- `+`：加法
- `-`：减法
- `*`：乘法
- `/`：除法
- `%`：求模
- `^`：求幂

比如我们要找到某个业务日志里面的错误率，就可以按照如下方式计算：

```logql
sum(rate({app="foo", level="error"}[1m])) / sum(rate({app="foo"}[1m]))
```

### 逻辑运算

集合运算仅在区间向量范围内有效，当前支持

- `and`：并且
- `or`：或者
- `unless`：排除

比如：

```logql
rate({app=~"foo|bar"}[1m]) and rate({app="bar"}[1m])
```

### 比较运算

LogQL 支持的比较运算符和 PromQL 一样，包括：

- `==`：等于
- `!=`：不等于
- `>`：大于
- `>=`: 大于或等于
- `<`：小于
- `<=`: 小于或等于

通常我们使用区间向量计算后会做一个阈值的比较，这对应告警是非常有用的，比如统计 5 分钟内 error 级别日志条目大于 10 的情况：

```logql
count_over_time({app="foo", level="error"}[5m]) > 10
```

我们也可以通过布尔计算来表达，比如统计 5 分钟内 error 级别日志条目大于 10 为真，反正则为假：

```logql
count_over_time({app="foo", level="error"}[5m]) > bool 10
```

## 注释

LogQL 查询可以使用 `#` 字符进行注释，例如：

```logql
{app="foo"} # anything that comes after will not be interpreted in your query
```

对于多行 LogQL 查询，可以使用 `#` 排除整个或部分行：

```logql
{app="foo"}
    | json
    # this line will be ignored
    | bar="baz" # this checks if bar = "baz"
```

## Pipeline Errors 管道错误

There are multiple reasons which cause pipeline processing errors, such as:有多种原因导致流水线处理错误，例如：

- A numeric label filter may fail to turn a label value into a number 数字标签过滤器可能无法将标签值转换为数字
- A metric conversion for a label may fail.标签的度量转换可能会失败。
- A log line is not a valid json document.日志行不是有效的 JSON 文档。
- etc…等等…

When those failures happen, Loki won’t filter out those log lines. Instead they are passed into the next stage of the pipeline with a new system label named `__error__`. The only way to filter out errors is by using a label filter expressions. The `__error__` label can’t be renamed via the language.当这些故障发生时，Loki 不会过滤掉这些日志线。相反，它们通过名为\_\_Error\_\_的新系统标签传递到管道的下一个阶段。过滤误差的唯一方法是使用标签过滤表达式。 \_\_Err\_\_标签无法通过语言重命名。
For example to remove json errors:例如要删除 JSON 错误：

```logql
{cluster="ops-tools1",container="ingress-nginx"}
    | json
    | __error__ != "JSONParserErr"
```
Logql
Alternatively you can remove all error using a catch all matcher such as `__error__ = ""` or even show only errors using `__error__ != ""`.或者，您可以使用捕获器删除所有匹配器，例如\_\_Error\_\_ =“”甚至仅显示使用\_\_Error\_\_！=“”。
The filter should be placed after the stage that generated this error. This means if you need to remove errors from an unwrap expression it needs to be placed after the unwrap.在生成此错误的阶段后应放置过滤器。这意味着如果您需要从未包装中删除从未包装表达式中删除错误，则需要放置在未包装之后。

```logql
quantile_over_time(
 0.99,
 {container="ingress-nginx",service="hosted-grafana"}
 | json
 | unwrap response_latency_seconds
 | __error__=""[1m]
 ) by (cluster)
```

Logql

> Metric queries cannot contains errors, in case errors are found during execution, Loki will return an error and appropriate status code.度量标准查询不能包含错误，以便在执行期间找到错误，Loki 将返回错误和适当的状态代码。
