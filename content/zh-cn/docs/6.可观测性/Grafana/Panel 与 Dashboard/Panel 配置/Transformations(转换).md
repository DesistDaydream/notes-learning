---
title: Transformations(转换)
---

# 概述

> 参考：
>
> - [官方文档，面板 - 转换](https://grafana.com/docs/grafana/latest/panels/transformations/)
> - [官方文档，面板 - 转换 - 转换类型和选项](https://grafana.com/docs/grafana/latest/panels/transformations/types-options/)(这里就是详解了每一种转换类型)

**Transformations(转换)** 经常用在 Table 面板中，我先以一个 Table 作为基础例子来说明 Transformations 的功能

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/gxof32/1616067878482-d907b3a2-1f60-4894-bc9c-75cc6ba232eb.png)

## 转换顺序

如果同时存在多个转换, 则从上往下一次执行它们，这就有点像 Linux 中的管道，每一个转换会产生一个新的结果，这个结果将会传递给下一个转换继续进行处理。直到所有转换执行完成，在面板展示最终数据。

转换顺序的特性可以在这个例子中可以得到充分体现

# Add field from calculation

Use the row values to calculate a new field

# Filter by name(根据字段名称进行过滤)

使用 regex pattern(正则表达式模式) 删除部分查询结果，模式可以是包含或者排除。过滤的对象是 Field(字段)。

下图经过过滤后，我们仅显示 instance、pod、Value 这几个字段

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/gxof32/1616067878535-d9af4670-b2e2-4973-9b34-63ea43dafacb.png)

# Filter data by query

Filter data by query. This is useful if you are sharing the results from a different panel that has many queries and you want to only visualize a subset of that in this panel.

# Group by

Group the data by a field values then process calculations for each group

# Join by field(连接字段) Outer join

Joins many time series/tables by a field. This can be used to outer join multiple time series on the _time_ field to show many time series in one table.

# Labels to fields

Groups series by time and return labels or tags as fields. Useful for showing time series with labels in a table where each label key becomes a seperate column

# Merge(合并)

Merge many series/tables and return a single table where mergeable values will be combined into the same row. Useful for showing multiple series, tables or a combination of both visualized in a table.

合并多个系列/表并返回一个表，其中可合并的值将被合并到同一行中。用于在一个表中显示多个序列、表或两者的组合。

# Organize fields(组织字段)

重命名，重新排序或隐藏查询返回的字段。

下图经过组织后，字段的名称得以改变

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/gxof32/1616067878497-a4cf0af9-8b5b-4b90-b2a4-d2510cac5503.png)

# Reduce(裁剪)

裁剪所有行或数据点，变为单个值。裁剪之后，可以通过 max、min、mean、last 之类的函数，显示对应的值。

下图就是裁剪之后的样式，原先的多个 Field，被整合到一个 Field，作为整个 Fild 的元素，每个元素都可以有一个或多个值(在对话框中可以选择要显示的值类型)。

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/gxof32/1616067878475-55c70061-ae9d-4057-ad1e-62c82249897b.png)

我们使用 DeBug 模式可以看到这些数据在转换前后的样子：

从下图可以看出来，转换前有 10 个 Fields(包括时间)。转换后，原 Fields 的 name 会作为新 Fields 的元素进行填充。

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/gxof32/1616067878506-6be06d38-2181-4238-bba5-d17711a6525b.png)

# [Rename by regex](https://grafana.com/docs/grafana/latest/panels/transformations/types-options/#rename-by-regex)(通过正则表达式重命名)

使用正则表达式和替换模式重命名查询结果。

在下面的示例中，可以看到 instance 标签的值都带着 `:9100`，实际情况下，我们并不需要显示端口，那么就可以通过 Rename by regex 将其去掉

![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/gxof32/1636270563979-9e684764-6cd5-435c-93c9-be126a8593f1.png)

应用这个转换类型后，可以看到只剩下除了 9100 以外的字符串了。

![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/gxof32/1636270580325-7250a179-2386-4cf9-8918-d372612715d4.png)

# Series to rows

Merge many series and return a single series with time, metric and value as columns. Useful for showing multiple time series visualized in a table.

合并多个系列，并以时间，度量和值作为列返回单个系列。 用于显示表格中可视化的多个时间序列。

将来自多个时间序列数据查询的结果合并为一个结果。 使用表格面板可视化效果时，这很有帮助。
