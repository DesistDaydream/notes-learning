---
title: Aggregation Operators(聚合运算符)
---

# 概述

> 参考：
>
> - [官方文档,Prometheus-查询-运算符-聚合运算符](https://prometheus.io/docs/prometheus/latest/querying/operators/#aggregation-operators)

Prometheus 还提供了下列内置的聚合运算符，这些运算符仅作用于瞬时向量。可以将瞬时表达式返回的样本数据进行聚合，形成一个新的时间序列。

- sum (求和)
- min (最小值)
- max (最大值)
- avg (平均值)
- stddev (标准差)
- stdvar (标准差异)
- count (计数)
- count_values (对 value 进行计数)
- bottomk (后 n 条时序)
- topk (前 n 条时序)
- quantile (分布统计)

## Syntax(语法)

**\[AggregationOperators] \[without|by (LabelName,....)] (\[Parameter,...] VectorExpression)**

- **Aggregation Operators** # 聚合运算符
- **without|by (LabelName,...)** # 若不指定该选项，则聚合全部数据的值。如果指定了，则按照指定的 LabelName 进行聚合。通过 without 和 by 可以按照样本的问题对数据进行聚合。该用法的示例图详见文末
  - **by** # 聚合 by 后面指定的 LabelName 样本数据，并将聚合以外的标签的移除
  - **without** # 与 by 相反，聚合 without 后面没有指定的 LabelName 样本数据。并将聚合以外的标签的移除
- **Parameter #** 参数，其中只有 count_values, quantile, topk, bottomk 支持
- **VectorExpression** # 向量表达式。详见 [PromQL](/docs/6.可观测性/Metrics/Prometheus/PromQL/PromQL.md) 章节

# sum 与 min 与 max 与 avg 详解

# count 与 count_values 详解

EXAMPLE

- 计算 up 序列 的值为 1 的序列总数
  - count by(job, namespace, service) (up == 1)

# quantile 详解

> 参考:
>
> - https://cloud.tencent.com/developer/news/319419
> - https://www.zhihu.com/question/20575291

其实，在 Prometheus 中，把 quantile(分位数)改为 percentile(百分位数)更准确

在这么一组数中

18,6,250,4,21,10,1,1,4274,5,102,15,5,3,10,1,5,3,177,5,34,45,1,5,15

quantile(0.5,上述所有数的集合)，如果使用该公式计算 0.5 百分位数，则值为 6。也就是说，比 6 大数的有 50%，比 6 小的数有 50%

quantile(0.99,上述所有数的集合)，如果使用该公式计算 0.99 百分位数，则值为 3308.239999999992。比 3308 大的数有 1%，比 3308 小的有 99%

如果使用响应时间来举例，那么，0.99 百分位的值为 10，意味着：在此时此刻，所有请求的响应时间中，有 99%都是 10 以上秒，剩下的不到 10 秒

# Example

如果 http_requests_total 指标一共有 3 个 label，分别是 application、instance、group。那么下面的两种运算方式得出的结果是一样的

- sum without (instance) (http_requests_total)
- 等价于
- sum by (code,handler,job,method) (http_requests_total)

如果只需要计算整个应用的 HTTP 请求总量，可以直接使用表达式

- sum(http_requests_total)

count_values 用于时间序列中每一个样本值出现的次数。count_values 会为每一个唯一的样本值输出一个时间序列，并且每一个时间序列包含一个额外的标签。

- count_values("count", http_requests_total)

topk 和 bottomk 则用于对样本值进行排序，返回当前样本值前 n 位，或者后 n 位的时间序列。获取 HTTP 请求数前 5 位的时序样本数据，可以使用表达式：

- topk(5, http_requests_total)

quantile 用于计算当前样本数据值的分布情况 quantile(φ, express)其中 0 ≤ φ ≤ 1。例如，当 φ 为 0.5 时，即表示找到当前样本数据中的中位数

- quantile(0.5, http_requests_total)

# 效果示例图

不进行聚合运算的数据

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/fybu67/1616069162982-ceaed866-8e18-4a11-ba5b-a248b397ef1d.jpeg)

不使用 by 或者 without，聚合运算所有数据，得出唯一一个值

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/fybu67/1616069162971-82c3f7c8-4fd3-4e75-be31-b7c41e43b2d5.jpeg)

使用 by，聚合具有相同 namespace 的样本值，得出一个或多个值

![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/fybu67/1660618023658-45664731-5ddb-455f-8a7b-42ad1c3c3dfa.png)
