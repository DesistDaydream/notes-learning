---
title: Binary Operators(二元运算符)
---

# 概述

> 参考：
>
> - [官方文档，查询 - 运算符](https://prometheus.io/docs/prometheus/latest/querying/operators/#binary-operators)

PromQL 支持基本的 逻辑 和 算术 运算符。 对于两个即时向量之间的运算，可以修改匹配行为。

使用 PromQL 除了能够方便的按照查询和过滤时间序列以外，PromQL 还支持丰富的运算符，用户可以使用这些运算符对进一步的对事件序列进行二次加工。这些运算符包括：数学运算符，逻辑运算符，布尔运算符等等。

**官方文档中，将时间序列中的标签称为 element(元素)**

# Arithmetic(算术) 二元运算符

PromQL 支持以下算术二元运算符：

- - (加法)
- - (减法)
- - (乘法)
- / (除法)
- % (求余)
- ^ (幂运算)

算术二元运算符可以实现如下三种类型的运算：

- **Between two scalars(标量与标量)**
- **Between an instant vector and a scalar(即时向量与标量)**
- **Between two instant vectors(即时向量与即时向量)**

## Between two scalars(标量与标量)

就是普通的数学运算，类似于 1+1、2\*3 等等，直接获取一个标量结果

## Between an instant vector and a scalar(即时向量与标量)

当瞬时向量与标量之间进行数学运算时，数学运算符会依次作用域瞬时向量中的每一个样本值，从而得到一组新的时间序列。

与 标量之间 的二元运算一样，只不过将即时向量表达式获取到的所有时间序列的值与标量进行运算，效果如下：

![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/zpuhbm/1626438137638-8b462402-4042-4c0d-b030-fc186973d5ab.png)

经过运算后：

![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/zpuhbm/1626438160237-51de2eda-a03a-494a-8b86-16ffdaeb0551.png)

## Between two instant vectors(即时向量与即时向量)

如果是即时向量与即时向量之间进行数学运算时，过程会相对复杂一点。 例如，如果我们想根据 node_disk_bytes_written 和 node_disk_bytes_read 获取主机磁盘 IO 的总量，可以使用如下表达式：

```promql
node_disk_read_bytes_total + node_disk_written_bytes_total
```

那这个表达式是如何工作的呢？依次找到与左边向量表达式的标签完全匹配的右边向量表示，并将两者进行进行运算。同时新的时间序列将不会包含指标名称。 该表达式返回结果的示例如下所示：

![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/zpuhbm/1626440024200-621fc788-29f8-4455-b933-4dbc3fa3a881.png)
如果运算符左右两边的向量表达式没有匹配到，则直接丢弃，效果如下：
![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/zpuhbm/1626440170722-a5bee3be-4df9-4553-9525-3aadd8c193f0.png)

# Comparison(比较) 运算符

目前，Prometheus 支持以下布尔运算符如下：

- \== (相等)
- != (不相等)
- > (大于)
- < (小于)
- > \= (大于等于)
- <= (小于等于)

在 PromQL 通过标签匹配模式，用户可以根据时间序列的特征维度对其进行查询。而 比较运算 则支持用户根据时间序列中样本的值，对时间序列进行过滤。

比如有这么一种场景：

- 通过数学运算符我们可以很方便的计算出，当前所有主机节点的内存使用率：
  - (node_memory_bytes_total - node_memory_free_bytes_total) / node_memory_bytes_total
- 而系统管理员在排查问题的时候可能只想知道当前内存使用率超过 95%的主机呢？通过使用比较运算，就可以方便的获取到该结果：
  - (node_memory_bytes_total - node_memory_free_bytes_total) / node_memory_bytes_total > 0.95

即时向量与标量进行布尔运算时，PromQL 依次比较向量中的所有时间序列样本的值，如果比较结果为 true 则保留，反之丢弃。

即时向量与即时向量直接进行布尔运算时，同样遵循默认的匹配模式：依次找到与左边向量元素匹配（标签完全一致）的右边向量元素进行相应的运算，如果没找到匹配元素，则直接丢弃。

## 使用 bool 修饰符改变比较运算符的行为

布尔运算符的默认行为是对时序数据进行过滤。而在其它的情况下我们可能需要的是真正的布尔结果。例如，只需要知道当前模块的 HTTP 请求量是否>=1000，如果大于等于 1000 则返回 1（true）否则返回 0（false）。这时可以使用 bool 修饰符改变布尔运算的默认行为。 例如：

    http_requests_total > bool 1000

使用 bool 修改符后，布尔运算不会对时间序列进行过滤，而是直接依次瞬时向量中的各个样本数据与标量的比较结果 0 或者 1。从而形成一条新的时间序列。

    http_requests_total{code="200",handler="query",instance="localhost:9090",job="prometheus",method="get"}  1
    http_requests_total{code="200",handler="query_range",instance="localhost:9090",job="prometheus",method="get"}  0

同时需要注意的是，如果是在两个标量之间使用布尔运算，则必须使用 bool 修饰符

2 == bool 2 # 结果为 1

# Logical(逻辑) 运算符

目前，Prometheus 支持以下逻辑运算符(这些运算符只能作用在瞬时向量上)：

- and (并且)
- or (或者)
- unless (排除)

使用即时向量表达式能够获取到一个包含多个时间序列的集合，我们称为瞬时向量。 通过集合运算，可以在两个瞬时向量与瞬时向量之间进行相应的集合操作。

vector1 and vector2 会产生一个由 vector1 的元素组成的新的向量。该向量包含 vector1 中完全匹配 vector2 中的元素组成。

vector1 or vector2 会产生一个新的向量，该向量包含 vector1 中所有的样本数据，以及 vector2 中没有与 vector1 匹配到的样本数据。

vector1 unless vector2 会产生一个新的向量，新向量中的元素由 vector1 中没有与 vector2 匹配的元素组成。

# 运算符优先级

对于复杂类型的表达式，需要了解运算操作的运行优先级

例如，查询主机的 CPU 使用率，可以使用表达式：

```promql
100 * (1 - avg (irate(node_cpu{mode='idle'}[5m])) by(job) )
```

其中 irate 是 PromQL 中的内置函数，用于计算区间向量中时间序列每秒的即时增长率。关于内置函数的部分，会在下一节详细介绍。

在 PromQL 操作符中优先级由高到低依次为：

1. ^
2. \*, /, %
3. +, -
4. \==, !=, <=, <, >=, >
5. and, unless
6. or
