---
title: PostgreSQL SQL
linkTitle: PostgreSQL SQL
date: 2024-11-28T09:07
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

## Schema 查询

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
WHERE table_schema = 'cheat'
  AND table_type = 'BASE TABLE';
```

# 函数 与 运算符

> 参考：
>
> - [官方文档，9.函数与运算符](https://www.postgresql.org/docs/current/functions.html)

## 时间函数

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