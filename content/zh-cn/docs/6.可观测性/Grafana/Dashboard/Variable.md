---
title: Variable
linkTitle: Variable
weight: 20
---

# 概述

> 参考：
>
> - [官方文档，仪表盘 - 变量](https://grafana.com/docs/grafana/latest/dashboards/variables/)
>     - https://grafana.com/docs/grafana/latest/visualizations/dashboards/variables/

我们可以使用 `node_memory_MemAvailable_bytes{instance="192.168.254.253:9100"}` 这种 PromQL 展示某台设备的可用内存。但是我们想要展示其他设备的时候，都要手动修改，或者额外再创建一个 Panel，这都会导致配置 Dashboard 极其繁琐，我们需要一个动态的方式来展示我们想要展示的内容

Variable 是 value(值) 的 placeholder(占位符)。当我们更改值时，使用该变量的位置也会变为新值。上面例子中的 `192.168.254.253:9100` 就可以使用变量代替，通过改变变量的值，即可简单方便得直接展示一台或多台设备的指标。像下图，设置了一个 server 变量。我们可以将表达式中的部分替换为变量引用 `node_memory_MemAvailable_bytes{instance=~"${server}"}` 。

![800](https://grafana.com/media/docs/grafana/dashboards/screenshot-selected-variables-v12.png)

要查看和设置变量，导航至 Dashboard 中的 **Settings > Variables**，在这里除了可以自己定义变量，还能看到变量的 依赖关系（i.e. 哪些变量引用了哪些变量） 和 使用情况（i.e. 哪些 Panel 应用了哪些变量）

[Dashboard](/docs/6.可观测性/Grafana/Dashboard/Dashboard.md) 可以设置并使用的变量有如下几类（设置变量的时候同时也可以为变量赋值）：

- **Query** # 从数据源中查询并获取变量。
- **Custom** # 手动定义变量，多个变量以 `,` 分隔。
- **Text box** # 定义一个文本框变量，可以在其中输入任意字符串作为变量的值
- **Constant** # 定义一个隐藏的常量变量，对于要共享的仪表板中的指标前缀很有用。TODO: 没看懂有啥用
- **Data source** # 定义一个数据源，可以通过切换该变量，快速切换整个 Dashboard 所有查询所使用的数据源

> [!important] Data source 变量在共享 Dashboard 时尤其有用，不同 Grafana 的数据源的 ID 不同，换了一个 Grafana 并导入相同 Dashboard 后，Panel 的查询可能会由于找不到数据源而丢失。但是在 Panel 中使用 Data source 类型的变量作为数据源，即可完美避免该问题

- **Interval** # TODO: 好像有某些地方可以自动识别这个变量
- **Filter** # TODO: 感觉不到用处。官方说仅限 Prometheus, Loki, InfluxDB, Elasticsearch
- **Switch** # TODO: 有啥用？
- **Global variables** # 内置的全局变量。不用定义可以直接使用。详见下文 [全局内置变量](#全局内置变量)
- **Chained variables** #

# 创建变量

点击 **New avriable** 即可进入创建变量页面

![500](https://notes-learning.oss-cn-beijing.aliyuncs.com/grafana/dashboard/variable-new-1.png)


- **Variable type** # 选择变量类型。根据选择的不同类型，下面会现实不同的可设置项。
- **General** # 所有类型变量都会有的设置项。设置变量的基本信息：
    - Name # 变量的名称。这是<font color="#ff0000">重点</font>，调用变量时，使用的就是这里定义的名称。变量名称要唯一。
    - Label # 在 Dashboard 上现实的名字
    - Description # 变量的描述
    - HIde # 是否在 Dashboard 隐藏变量。Nothing 是不隐藏；Variable 是隐藏变量全部；Label 是仅隐藏变量的 Label
- **其他** # 根据选择的变量类型有不同的可配置项。这里使用的 Query 类型变量，所以展示的是 Query options。

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

---

在 PostgreSQL 中的最佳实践

```sql
-- 利用 timezone 函数改变时区；利用 to_char 改变时间格式。适用于字符串类型的时间
WHERE create_time BETWEEN to_char(timezone('UTC', '${__from:date}'), 'YYYY-MM-DD HH24:MI:SS') AND to_char(timezone('UTC','${__to:date}'), 'YYYY-MM-DD HH24:MI:SS')
```

### `$timeFilter` 或 `$__timeFilter`

https://grafana.com/docs/grafana/latest/dashboards/variables/add-template-variables/#timefilter-or-__timefilter

timeFilter 用以显示时间范围选择器选定的时间间隔。是一种类似把 from 与 to 两个变量合在一起的变量。

该变量主要应用在 SQL 的 WHERE 条件用，用以过滤特定时间范围内的数据。

TODO: 为什么在文本模式下渲染不出来？

### `$__range`

此变量表示时间选择器选择的时间范围之间的秒数。该值的本质是 `$__to - $__from` 两个变量相减。

### `$__interval`

https://grafana.com/docs/grafana/latest/dashboards/variables/add-template-variables/#__interval

可以使用 `$__interval` 变量作为参数按时间（对于 InfluxDB、MySQL、Postgres、MSSQL）、日期直方图间隔（对于 Elasticsearch）进行分组，或作为汇总函数参数（对于 Graphite）。

步长，格式是 30s、1h、5d、etc.

# 变量的引用语法

> 参考：
>
> - [官方文档，数据可视化 - 仪表盘 - 变量 -变量语法](https://grafana.com/docs/grafana/latest/visualizations/dashboards/variables/variable-syntax)

变量插值的格式取决于数据源，但在某些情况下，我们可能需要更改默认格式。

例如，MySQL 数据源默认将变量的多个值用逗号分隔并用引号括起来： `'server01','server02'` 。在某些情况下，您可能希望使用不带引号的逗号分隔字符串： `server01,server02` 。我们可以使用下面记录的引用变量时的语法来实现这一点。

TODO

# 最佳实践

## Variable 多值行为与 Repeat 组合的化学反应

记录于 2025-05-29

在使用某些 [数据库](/docs/5.数据存储/数据库/数据库.md)（e.g. ClickHouse） 的数据源中，如果是**多值**的变量，在引用变量时会在值的两侧使用 `'`(单引号)。若变量是表名，`FROM network_security.${table_name}` 这种语句的渲染结果会是 `FROM network_security.'mytable'`，此时将会报错无法执行查询。

但是若使用了 Repeat 功能，单引号将被摘掉，就可以正常查询了。
