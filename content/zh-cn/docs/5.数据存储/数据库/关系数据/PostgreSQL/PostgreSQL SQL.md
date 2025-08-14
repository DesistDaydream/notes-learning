---
title: PostgreSQL SQL
linkTitle: PostgreSQL SQL
weight: 20
---

# 概述

> 参考：
>
> - [官方文档，SQL 语言](https://www.postgresql.org/docs/current/sql.html)

TODO: 整理特定于 PostgreSQL 的 SQL

```plsql
# 创建新表
CREATE TABLE user_tbl(name VARCHAR(20), signup_date DATE);

# 插入数据
INSERT INTO user_tbl(name, signup_date) VALUES('张三', '2013-12-22');

# 从表中查询数据
SELECT * FROM user_tbl;

# 更新数据
UPDATE user_tbl set name = '李四' WHERE name = '张三';

# 删除记录
DELETE FROM user_tbl WHERE name = '李四' ;

# 添加栏位
ALTER TABLE user_tbl ADD email VARCHAR(40);

# 更新结构
ALTER TABLE user_tbl ALTER COLUMN signup_date SET NOT NULL;

# 更名栏位
ALTER TABLE user_tbl RENAME COLUMN signup_date TO signup;

# 删除栏位
ALTER TABLE user_tbl DROP COLUMN email;

# 表格更名
ALTER TABLE user_tbl RENAME TO backup_tbl;

# 删除表格
DROP TABLE IF EXISTS backup_tbl;
```

# 函数 与 运算符

> 参考：
>
> - [官方文档，9.函数与运算符](https://www.postgresql.org/docs/current/functions.html)

## 字符串处理

> 参考：
>
> - [官方文档，9.4 字符串 函数与运算符](https://www.postgresql.org/docs/current/functions-string.html)
> - [官方文档，9.5 二进制字符串 函数与运算符](https://www.postgresql.org/docs/current/functions-binarystring.html)

### string_to_array

`string_to_array(STRING text, Delimiter text [, NullString text ]) → text[]`

拆分文本类型的 string，以 Delimiter 作为分隔符，并将结果字段形成文本数组。如果 Delimiter 为 NULL，则字符串中的每个字符将成为数组中的元素。如果分隔符是空字符串，则该字符串将被视为单个字段。如果提供了 NullString 并且不为 NULL，则与该字符串匹配的字段将替换为 NULL。

e.g.

```sql
SELECT string_to_array('xx~~yy~~zz', '~~', 'yy');
-- 结果为: {xx,NULL,zz}
```

---

最佳实践

> `{xx,NULL,zz}` 结果可以利用类型转换 `unnest({xx,NULL,zz}::bigint[])` 转为被过滤语句识别的各种类型

```sql
-- 字符串转数组
SELECT string_to_array('15|8356107', '|')
-- 将 text 数组转为 bigint 数组
SELECT unnest('{15,8356107}'::bigint[]);
-- 可以合起来
SELECT unnest(string_to_array('15|8356107', '|')::bigint[])

-- 这个结果可以放在 WHERE 过滤语句中
SELECT *
FROM my_schema.my_table
WHERE id IN (SELECT unnest(string_to_array('15|8356107', '|')::bigint[]))
-- 其中 15|8356107 可以在 Grafana 中使用变量引入
```

### 字符串格式化

https://www.postgresql.org/docs/current/functions-formatting.html

PostgreSQL 格式化函数提供了一组强大的工具，用于将各种数据类型（日期/时间、整数、浮点、数字）转换为格式化字符串以及从格式化字符串转换为特定数据类型。这些函数都遵循通用的调用约定：第一个参数是要格式化的值，第二个参数是定义输出或输入格式的模板。

```sql
SELECT to_char(timezone('UTC', '2024-12-13T16:23:58.115Z'), 'YYYY-MM-DD HH24:MI:SS')
-- 结果为: 2024-12-13 16:23:58
```

### 编码/解码字符串

`encode ( bytes bytea, format text ) → text` 与 `decode ( string text, format text ) → bytea`

encode 将 bytes 以 format 格式编码为 text；decode 将 string 以 format 解码为 bytea

format 支持的格式有：

- base64
- escape
- hex

```sql
SELECT decode('RGVzaXN0RGF5ZHJlYW0K', 'base64')
-- 结果未: DesistDaydream
```

## 时间处理

> 参考：
>
> - [官方文档，9.9 日期/时间 函数与运算符](https://www.postgresql.org/docs/current/functions-datetime.html)


### 基本的简单函数

PostgreSQL 提供了许多返回与当前日期和时间相关的值的函数。这些 SQL 标准函数的返回时依据当前事务的开始时间：

> 有的函数甚至不需要加参数和 `()` 符号即可直接使用

- CURRENT_DATE
- CURRENT_TIME
- CURRENT_TIMESTAMP
- CURRENT_TIME(precision)
- CURRENT_TIMESTAMP(precision)
- LOCALTIME
- LOCALTIMESTAMP
- LOCALTIME(precision)
- LOCALTIMESTAMP(precision)

示例如下

```sql
SELECT CURRENT_TIME;
_Result:_ `14:39:53.662522-05`
SELECT CURRENT_DATE;
_Result:_ `2019-12-23`
SELECT CURRENT_TIMESTAMP;
_Result:_ `2019-12-23 14:39:53.662522-05`
SELECT CURRENT_TIMESTAMP(2);
_Result:_ `2019-12-23 14:39:53.66-05`
SELECT LOCALTIMESTAMP;
_Result:_ `2019-12-23 14:39:53.662522`
```

### date_bin

https://www.postgresql.org/docs/current/functions-datetime.html#FUNCTIONS-DATETIME-BIN

`date_bin(stride, source, origin)`

函数 date_bin 将输入时间戳“装箱”到与指定原点对齐的指定间隔（stride 步长）中。

人话：修改 source 列的时间，修改逻辑示例如下：

```sql
SELECT date_bin('15 minutes', TIMESTAMP '2020-02-11 15:44:17', TIMESTAMP '2001-01-01');
```

结果为 `2020-02-11 15:30:00`

其中首先要计算从 2001-01-01 开始，每隔 15 分钟的时间点。计算到 2020-02-11 时，会有这么几个时间：

- 2020-02-11 15:15:00
- 2020-02-11 15:30:00
- 2020-02-11 15:45:00

date_bin 会将 source 向下取整，选择上面三个时间中，往旧时间数离自己最近的时间，也就是 2020-02-11 15:30:00

若不想要整点，则可以修改 origin 参数，效果如下：

```sql
SELECT date_bin('15 minutes', TIMESTAMP '2020-02-11 15:44:17', TIMESTAMP '2001-01-01 00:02:30');
```

结果为 `2020-02-11 15:32:30`

> [!Tip]
> data_bin 函数常用来**实现**类似**时间序列数据**的效果，将很多行数据，按照时间块进行分组，以便更好得查看事件的发生频率或其他统计信息。e.g. 有一个表记录了用户访问日志（一次访问一行），可以使用 data_bin 将数据按小时或天进行汇总，查看每小时或每天的访问量
>
> 基于这种统计数据，还可以实现**数据可视化**（e.g. 使用 Grafana）。因为数据被聚合分组了，每隔时间点都有 N 个数据，正好对应 二维 坐标中横轴事件，纵轴数值。

data_bin 函数通常不会单独使用，而是与 COUNT() 函数联动。比如：

```sql
SELECT
  date_bin('1 hour', timestamp_column, 'epoch') AS binned_time,
  COUNT(*)
FROM my_table
GROUP BY binned_time
ORDER BY binned_time;
```

> Tips: 可以省略 ORDER，PostgreSQL 默认会使用 AES 对时间排序。

### 时区

> 参考：
>
> - [官方文档，9.9 日期/时间 函数与运算符 - 9.9.4 AT TIME ZONE 与 AT LOCAL](https://www.postgresql.org/docs/current/functions-datetime.html#FUNCTIONS-DATETIME-ZONECONVERT)
> - 时区文字定义 https://www.postgresql.org/docs/7.2/timezones.html TODO: 这是老版本的，新版本的在哪？

**AT TIME ZONE**

EXAMPLE

```sql
WHERE create_time AT TIME ZONE 'Asia/Shanghai' BETWEEN '2025-08-13 14:37:17' AND '2025-08-14 02:37:17'
```

上面例子中，最后查询的实际上是 2025-08-13 22:37:17 与 2025-08-14 10:37:17 之间的数据

# 最佳实践

## 基础信息

列出所有 Schema

```sql
SELECT *
FROM pg_namespace;
```

或

```sql
SELECT schema_name
FROM information_schema.schemata;
```

列出名为 cheat 的 Schema 下的所有表

```sql
SELECT table_name
FROM information_schema.tables
WHERE table_schema = 'my_schema'
  AND table_type = 'BASE TABLE';
```

列出表的列信息

```sql
SELECT *
FROM information_schema.columns
WHERE table_name = 'my_table';
```

> Notes: 不用写 Schema，因为 information_schema.columns 表中有 table_schema 列来表示表所属的 Schema

## 分页

下面以一个在 Grafana 中的应用为例，利用 LIMIT 和 OFFSET 实现分页查询效果。LIMIT 当作页容量，OFFSET 当作当前页（通过偏移 N 个页容量实现）

根据 `${pageSize}`(页容量) 生成一个总页数的列表，用以通过表单选择 `${currentPage}`(当前页)

```sql
SELECT
  generate_series(
    1,
    ceil(
      (
        SELECT COUNT(*)
        FROM my_schema.my_table
        WHERE $__timeFilter(create_time)
      ) :: numeric / ${pageSize} :: numeric
    )
  )
```

```sql
SELECT
  COUNT(*) OVER() as "总数",
  *
FROM my_schema.my_table
WHERE $__timeFilter(create_time)
LIMIT ${pageSize}
OFFSET (${currentPage} - 1) * ${pageSize}
```
