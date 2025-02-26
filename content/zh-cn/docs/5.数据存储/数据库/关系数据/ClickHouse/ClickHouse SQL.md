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

SELECT 有 EXCEPT 修饰符，可以用于排除某些列，比如下面的 SQL 可以从表中排除 domain 与 url 两列（i.e 返回结果没有这俩列）

```sql
SELECT * EXCEPT (domain,url)
FROM nginx_logs.nginx_access
```

# Function

> 参考：
>
> - [官方文档，SQL 参考 - 函数](https://clickhouse.com/docs/en/sql-reference/functions)

- [Regular Functions](#Regular%20Functions)(常规函数)
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

### JSON

https://clickhouse.com/docs/en/sql-reference/functions/json-functions

解析 JSON 数据，主要包含两类

- [`simpleJSON*` (`visitParam*`)](https://clickhouse.com/docs/sql-reference/functions/json-functions#simplejson-visitparam-functions)  # 这是为了极为快速解析JSON的有限子集而制造的。
- [`JSONExtract*`](https://clickhouse.com/docs/sql-reference/functions/json-functions#jsonextract-functions) # 是为了解析普通json而制造的。可以处理的逻辑更复杂

### 字典

https://clickhouse.com/docs/sql-reference/functions/ext-dict-functions

**Dictionary(字典)** 相关函数用于从 [Dictionary](https://clickhouse.com/docs/dictionary) 类型的数据中，根据 Key 获取各种 Value

比如，假设有一个名为 `products` 的字典，包含 产品ID、名称、价格信息

```sql
SELECT dictGet('products','name', products_id) AS product_name
FROM shopping_record
```

这个 SQL 的意思是：

- 使用 shopping_record 表的 products_id 列的值作为对比用的 key
- 从 products 表中找出来表的 **主键 与 shopping_record 表中 products_id 列的值相等**的行，获取该行 name 列的值，作为 product_name
- <font color="#ff0000">Notes: 这里函数中没有填写要与 products 表中哪个列的值与 products_id 进行匹配。这是因为在创建 Dictionary 类型数据时，会指定哪个列作为 Key，所有调用 dict 相关函数的请求都会与 Key 列进行对比</font>。比如下面的 structure 中，除了 attribute 还有 id，那么 id 就是 key，所有的 attribute 就是该 key 的值。所以上面的例子就是 shopping_record 中的 products_id 列的值与 products 表中的 id 列进行对比，相同的就会找到对应 name 的值返回

```xml
<dictionary>
    <name>products</name>
    <structure>
        <id>
            <name>product_id</name>
        </id>
        <attribute>
            <name>name</name>
            <type>String</type>
        </attribute>
        <attribute>
            <name>price</name>
            <type>Float32</type>
        </attribute>
    </structure>
    <!-- 其他配置 -->
</dictionary>
```

---

从字典表中获取值的集中基础函数：

```sql
dictGet('dict_name', attr_names, id_expr)  
dictGetOrDefault('dict_name', attr_names, id_expr, default_value_expr)  
dictGetOrNull('dict_name', attr_name, id_expr)
```

- **dict_name** # 表名。要从哪个表获取值
- **attr_names** # 属性名。要从表中获取哪列的值
- **id_expr** # 根据本表中哪个字段从字典表中查找对应的值
- **default_value_expr** # 如果没查到对应的值返回的默认值。

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

## 分页

获取分页列表（1 至 n 的页码列表）

```sql
WITH 
    total_count AS (
        SELECT COUNT(*) AS total 
        FROM network_security.${tableName}
        WHERE $__timeFilter(create_time)
          AND command_id IN (${commandID})
    )
SELECT 
    arrayJoin(range(1, 
        CAST(
            CEIL(total_count.total * 1.0 / ${pageSize}) AS UInt32
        ) + 1
    )) AS page_number
FROM total_count;
```

分页查询

```sql
SELECT * FROM (
  SELECT 
    COUNT(*) OVER() as "指令总数",
    *
  FROM network_security.${tableName} 
  WHERE $__timeFilter(create_time)
  AND command_id IN (${commandID})
) sub
LIMIT ${pageSize} 
OFFSET (${currentPage} - 1) * ${pageSize}
```

# 处理 JSON 结构与 Base64 编码

```sql
SELECT
    JSONExtractString(base64Decode(detail_info), 'accuracy') AS accuracy_value,
    COUNT(*) AS accuracy_count
FROM network_security.td_host_control_event
WHERE $__timeFilter(found_time)
GROUP BY accuracy_value
```
