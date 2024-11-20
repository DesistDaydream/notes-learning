---
title: ClickHouse SQL
linkTitle: ClickHouse SQL
date: 2024-11-20T09:43
weight: 20
---

# 概述

> 参考：
>
> - [官方文档，SQL 参考](https://clickhouse.com/docs/en/sql-reference)



# Function

> 参考：
>
> - [官方文档，SQL 参考 - 函数](https://clickhouse.com/docs/en/sql-reference/functions)

- Regular Functions(常规函数)
- Aggregate Functions(聚合函数)
- Table Functions(表函数)
- Window Functions(窗口函数)

高阶函数与 lambda 函数，形式为 `params -> expr`。箭头左侧是一个形式参数，右侧是一个表达式。整体构成一个函数

```
onNameFunc(x -> x * 2)
对应一个下面这种函数
func onNameFunc(x int) {
  x * 2
}
```

## Regular Functions

https://clickhouse.com/docs/en/sql-reference/functions

### 数组

https://clickhouse.com/docs/en/sql-reference/functions/array-functions

`arrayExists([func,] arr1, ...)` # func 是一个高阶函数，可以接受 lambda 函数

- e.g. `arrayExists(x -> x = src_ip, [${example_array}])` # `${example_array}` # 变量是一个数组，作为参数传递给 x，arrayExists 将会逐一一对 x 中的元素，执行 `x = src_ip` 表达式。主要用于判断 变量中的元素是否等于 src_ip

### 日期与时间

https://clickhouse.com/docs/en/sql-reference/functions/date-time-functions

toStartOfXXX

类似 [PostgreSQL](/docs/5.数据存储/数据库/关系数据/PostgreSQL/PostgreSQL.md) 的 date_bin 函数，将日期和时间向下舍入到 XXX，主要是用来依据时间进行数据聚合以实现时间序列数据的效果，并可在 Grafana 中绘制时间序列图表。

有很多现成的函数和 1 个通用函数

- [toStartOfWeek](https://clickhouse.com/docs/en/sql-reference/functions/date-time-functions#tostartofweek)
- [toLastDayOfWeek](https://clickhouse.com/docs/en/sql-reference/functions/date-time-functions#tolastdayofweek)
- [toStartOfDay](https://clickhouse.com/docs/en/sql-reference/functions/date-time-functions#tostartofday)
  - e.g. `SELECT toStartOfDay(toDateTime('2023-04-21 10:20:30'))` 结果为 `2023-04-21 00:00:00`
- [toStartOfHour](https://clickhouse.com/docs/en/sql-reference/functions/date-time-functions#tostartofhour)
  - e.g. `SELECT toStartOfHour(toDateTime('2023-04-21 10:20:30'))` 结果为 `2023-04-21 10:00:00`
- [toStartOfMinute](https://clickhouse.com/docs/en/sql-reference/functions/date-time-functions#tostartofminute)
- [toStartOfSecond](https://clickhouse.com/docs/en/sql-reference/functions/date-time-functions#tostartofsecond)
- [toStartOfMillisecond](https://clickhouse.com/docs/en/sql-reference/functions/date-time-functions#tostartofmillisecond)
- [toStartOfMicrosecond](https://clickhouse.com/docs/en/sql-reference/functions/date-time-functions#tostartofmicrosecond)
- [toStartOfNanosecond](https://clickhouse.com/docs/en/sql-reference/functions/date-time-functions#tostartofnanosecond)
- [toStartOfFiveMinutes](https://clickhouse.com/docs/en/sql-reference/functions/date-time-functions#tostartoffiveminutes)
- [toStartOfTenMinutes](https://clickhouse.com/docs/en/sql-reference/functions/date-time-functions#tostartoftenminutes)
- [toStartOfFifteenMinutes](https://clickhouse.com/docs/en/sql-reference/functions/date-time-functions#tostartoffifteenminutes)
- [toStartOfInterval](https://clickhouse.com/docs/en/sql-reference/functions/date-time-functions#tostartofinterval) # 通用函数，可以在参数中指定 向上/向下 舍入的具体逻辑。
  - e.g. `toStartOfInterval(t, INTERVAL 1 YEAR)` 的返回值与 `toStartOfYear(t)` 相同
  - e.g. `SELECT toStartOfInterval(toDateTime('2023-01-01 14:45:00'), INTERVAL 1 MINUTE, toDateTime('2023-01-01 14:35:30'))` 结果为 `2023-01-01 14:44:30`

# 最佳实践

## 显示一些基础信息

**列出所有数据库**

```sql
SHOW DATABASES;
```

**列出所有 Tables**

```sql
SELECT *
FROM system.tables
```

**显示指 Table 的信息**

```sql
DESCRIBE network_security.td_his_evil_sample_query
```

或

```sql
SELECT *
FROM system.columns
WHERE table = 'my_table' AND database = 'my_database';
```

**列出所有 Views**

```sql
SELECT *
FROM system.tables
WHERE engine = 'View'
```

