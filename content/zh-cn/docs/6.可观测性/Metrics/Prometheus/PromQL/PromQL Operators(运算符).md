---
title: PromQL Operators(运算符)
---

# 概述

> 参考：
>
> - [官方文档，Prometheus-查询-运算符](https://prometheus.io/docs/prometheus/latest/querying/operators/)

# Binary Operators(二元运算符)

PromQL 支持基本的 逻辑 和 算术 运算符。 对于两个即时向量之间的运算，可以修改匹配行为。

使用 PromQL 除了能够方便的按照查询和过滤时间序列以外，PromQL 还支持丰富的运算符，用户可以使用这些运算符对进一步的对事件序列进行二次加工。这些运算符包括：数学运算符，逻辑运算符，布尔运算符等等。

详见《[Binary Operators(二元运算符)](/docs/6.可观测性/Metrics/Prometheus/PromQL/Binary%20Operators(二元运算符).md)》章节

## Vector Matching(向量匹配)

向量与向量之间进行运算操作时会基于默认的匹配规则：

- 右边向量表达式获取到的时间序列与左边向量表达式获取到的时间序列的标签进行一一匹配，完全匹配到的两个时间序列进行运算，没有匹配到的直接丢弃

说白话就是，两条时间序列要想进行二元运算，他们的标签必须完全相同才可以。

**这时就产生问题了，如果我想让两条标签不同的序列进行二元运行，怎么办呢？这时候就需要使用 Vector Matching(向量匹配) 来扩展二元运算的功能。**通过 Vector Matching，我们可以根据标签的匹配规则，选择出来符合要求的时间序列进行二元运算。

接下来将介绍在 PromQ L 中有两种典型的匹配模式：一对一(one-to-one)、多对一(many-to-one)、一对多(one-to-many)。

向量匹配常用于生成新的时间序列

### 一对一匹配

一对一匹配模式会从操作符两边表达式获取的瞬时向量依次比较并找到唯一匹配(标签完全一致)的样本值，然后进行二元运算。

#### Syntax(语法)

```promql
<VectorExpr1> <BinaryOperators> ignoring(LabelList) <VectorExpr2>
<VectorExpr1> <BinaryOperators> on(LabelList) <VectorExpr2>
```

on 与 ignoring 关键字会将其左右两侧表达式中标签进行匹配，根据其指定的 LabelList 来匹配标签，匹配到的序列将会执行二元运算

- **ignoreing(LabelList)** # 匹配不包含 LabelList 的序列。
- **on(LabelList)** # 匹配包含 LabelList 的序列。

例如当存在样本：

```text
method_code:http_errors:rate5m{method="get", code="500"}  24
method_code:http_errors:rate5m{method="get", code="404"}  30
method_code:http_errors:rate5m{method="put", code="501"}  3
method_code:http_errors:rate5m{method="post", code="500"} 6
method_code:http_errors:rate5m{method="post", code="404"} 21

method:http_requests:rate5m{method="get"}  600
method:http_requests:rate5m{method="del"}  34
method:http_requests:rate5m{method="post"} 120
```

使用 PromQL 表达式：

```promql
method_code:http_errors:rate5m{code="500"} / ignoring(code) method:http_requests:rate5m
或者
method_code:http_errors:rate5m{code="500"} / on(method) method:http_requests:rate5m
```

该表达式会返回在过去 5 分钟内，HTTP 请求状态码为 500 的在所有请求中的比例。如果没有使用 ignoring(code)，操作符两边表达式返回的瞬时向量中将找不到任何一个标签完全相同的匹配项。

因此结果如下：

```promql
{method="get"}  0.04            //  值就是 24 / 600 得到的结果
{method="post"} 0.05            //  值就是 6 / 120 得到的结果
```

同时由于 method 为 put 和 del 的样本找不到匹配项，因此不会出现在结果当中。

### 多对一 && 一对多

多对一 和 一对多 两种匹配模式指的是“一”侧的每一个向量元素可以与"多"侧的多个元素匹配的情况。这里所谓的 `一` 和 `多`，其实就是相对表达式中，`左侧向量表达式` 和 `右侧向量表达式`。在这种情况下，通过 group_left 或者 group_right 这两个修饰符来确定 `左边/右边` 哪边的表达式充当“多”的角色(也就指拥有更高的基数(higher cardinality))。

多对一 和 一对多 两种模式一定是出现在操作符两侧表达式返回的向量标签不一致的情况。因此需要使用 ignoring 和 on 修饰符来排除或者限定匹配的标签列表。

#### Syntax(语法)

```promql
<VectorExpr1> <BinaryOperators> ignoring(LabelList) group_left(LabelList) <VectorExpr2>
<VectorExpr1> <BinaryOperators> ignoring(LabelList) group_right(LabelList) <VectorExpr2>
<VectorExpr1> <BinaryOperators> on(LabelList) group_left(LabelList) <VectorExpr2>
<VectorExpr1> <BinaryOperators> on(LabelList) group_right(LabelList) <VectorExpr2>
```

group_left 与 group_right 修饰符用来指定以左边或右边的向量表达式为主：

- **group_left(LabelList)** # 表示以**左边的 VectorExpr1 为主**，生成的新序列将会包含 **左侧序列** 中的所有标签以及 LabelList 里指定的标签
- **group_right(LabelList)** # 表示以**右边的 VectorExpr2 为主**，生成的新序列将会包含 **右侧序列** 中的所有标签以及 LabelList 里指定的标签

例如,使用表达式：

```promql
method_code:http_errors:rate5m / on(method) group_left method:http_requests:rate5m
```

该表达式中，左向量 method_code:http_errors:rate5m 包含两个标签 method 和 code。而右向量 method:http_requests:rate5m 中只包含一个标签 method，因此匹配时需要使用 ignoring 限定匹配的标签为 method。 在限定匹配标签后，右向量中的元素可能匹配到多个左向量中的元素 因此该表达式的匹配模式为多对一，需要使用 group_left 修饰符指定左向量具有更好的基数。

最终的运算结果如下：

```promql
{method="get", code="500"}  0.04            //  24 / 600
{method="get", code="404"}  0.05            //  30 / 600
{method="post", code="500"} 0.05            //   6 / 120
{method="post", code="404"} 0.175           //  21 / 120
```

提醒：group 修饰符只能在比较和数学运算符中使用。在逻辑运算 and,unless 和 or 才注意操作中默认与右向量中的所有元素进行匹配。

# Aggregation Operators(聚合运算符)

> 参考：
>
> - [官方文档](https://prometheus.io/docs/prometheus/latest/querying/operators/#aggregation-operators)

聚合运算符可以聚合单个即时向量的元素，从而产生一个包含较少元素且具有聚合值得新向量。

聚合运算符既可以用于聚合所有标签维度，也可以通过包含一个 without 或 by 子句来保留不同的维度。这些从句可以用在短语的前面或后面。

详见：[Aggregation Operators(聚合运算符)](/docs/6.可观测性/Metrics/Prometheus/PromQL/Aggregation%20Operators(聚合运算符).md)
