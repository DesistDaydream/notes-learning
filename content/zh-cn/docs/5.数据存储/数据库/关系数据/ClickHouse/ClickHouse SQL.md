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

# SQL 关键字

可以在 [官方文档，SQL 参考 - 语句](https://clickhouse.com/docs/en/sql-reference/statements) 看到 ClickHouse 支持的所有 SQL 基础关键字。诸如常见的 SELECT、INSERT、etc. 还有一些独属于 ClickHouse 的关键字，e.g. KILL、OPTIMIZE、etc.

# Function

> 参考：
>
> - [官方文档，SQL 参考 - 函数](https://clickhouse.com/docs/en/sql-reference/functions)

- Regular Functions(常规函数)
- Aggregate Functions(聚合函数)
- Table Functions(表函数)
- Window Functions(窗口函数)

高阶函数与 lambda 函数，形式为 `params -> expr`。箭头的 左侧是形式参数；右侧是表达式。整体构成一个函数，就像这样：

```text
oneNamedFunc(x -> x * 2)
```

可以把这种函数理解成一种通用的样子：

```text
func oneNamedFunc(x any) {
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

大多数的日期与时间相关的函数都接受时区参数（可选的），e.g. Asia/Shanghai。默认是本地时区。下面是一些最基本的日期与时间函数示例：

```sql
SELECT
    toDateTime('1732947066') AS "from_unix-timestamp",
    toDateTime('2016-06-15 23:00:00') AS time,
    toDate(time) AS date_local,
    toDate(time, 'Asia/Yekaterinburg') AS date_yekat,
    toString(time, 'US/Samoa') AS time_samoa
```

结果为：

```sql
┌─from_unix-timestamp─┬────────────────time─┬─date_local─┬─date_yekat─┬─time_samoa──────────┐
│ 2024-11-30 14:11:06 │ 2016-06-15 23:00:00 │ 2016-06-15 │ 2016-06-15 │ 2016-06-15 04:00:00 │
└─────────────────────┴─────────────────────┴────────────┴────────────┴─────────────────────┘
```

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

## Aggregate Functions


## Table Functions

https://clickhouse.com/docs/en/sql-reference/table-functions

Table 函数可以用来构造一个新的表格式的数据，比如 `select toDate('2010-01-01') + number as d FROM numbers(2);` 可以生成 2 行数据，格式是像这样的表格

| d                   |
| ------------------- |
| 2010-01-01 00:00:00 |
| 2010-01-02 00:00:00 |

## Window Functions


# 最佳实践

一些 SQL 的别名，在 CLI 中可用的快捷指令

- `\l` - SHOW DATABASES
- `\d` - SHOW TABLES
- `\c <DATABASE>` - USE DATABASE
- `.` - repeat the last query

## 显示一些基础信息

**列出所有数据库**

```sql
show databases;
```

**列出 my_database 库中的所有表**（若不指定 from my_database 则列出当前数据库中的所有表）

```sql
show tables from my_database;
```

**显示所有 Tables 的元信息**（e.g. 创建语句、引擎、UUID、etc.）

```sql
SELECT *
FROM system.tables
```

**显示指定表的列信息**（e.g. 列名、类型、默认值、etc.）

```sql
DESCRIBE my_database.my_table
```

或

```sql
-- system 库中的列信息还有 数据压缩情况、etc. 更多信息
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

## 基础增删改查

**清空 my_database.my_table 表中的数据**（仅清空数据保留表结构）

```sql
TRUNCATE TABLE IF EXISTS my_database.my_table;
```

> TRUNCATE 比 DELETE FROM 更高效，因为它不会一条条删除记录，而是直接释放存储空间。

**删除 my_database.my_table 表**

```sql
DROP TABLE IF EXISTS my_database.my_table;
```

**删除 my_database 数据库**

```sql
DROP DATABASE IF EXISTS my_database;
```

> [!Attention] ！！！该操作不可逆！！！`DROP DATABASE` 会删除数据库及其中的所有表

