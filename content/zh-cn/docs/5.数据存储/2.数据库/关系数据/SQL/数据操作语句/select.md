---
title: select
---

# 概述

[select](https://dev.mysql.com/doc/refman/8.0/en/select.html) 用来从一个或多个表中检索匹配到的行，并将结果存储在一个结果表中。

```sql
select
	[Modifiers]
	<SelectExpr>
from
	<TableReferences>
[where <WhereCondition>]
[limit N]
[offset M];
```

- **select** # 查询数据关键字
  - **Modifiers** # [修饰符](#qp9T0)。在 select 关键字之后，可是使用许多影响语句操作的修饰符。
  - **SelectExpr** # 要检索的列，可以使用 `*` 以检索所有列。
- **from** # 关键字，指定要检索行的一个或多个表，多个表使用 , 分隔。
  - **TableReferences** # 表引用。除了基本的直接使用表名，还可以通过 join 和 union 子句使用表达式引用表，甚至可以引用通过 select 返回的表。
- **where** # 从 from 中指定的表中所检索的行需要满足 where 定义的匹配条件。
  - **WhereCondition** # 多个条件以 `and` 和 `or` 关键字连接。
- **limit** # 为返回的记录数
- **offset** # 可以指定 select 语句开始查询的数据偏移量。默认偏移量为 0

**Notes：**

- select 输出的内容本质是一个新表，所以我们时常可以看到下面这种语句

```sql
select t2.Name,
from (
    select t.Name
    from (
        select *
        from t_diversion diversion
            left join t_data_lake lake on diversion.data_lake_code = lake.code
            left join t_private_net net on diversion.private_net_code = net.code
        ) t
) t2
```

最里层的 select 返回了一个表，名命名为 t，然后第二层 select 的 from 来自于最里层 select 生成的表，以此类推。

## 代码处理顺序

虽然 SELECT 位于语句最前面，它在逻辑处理中，基本上是最后一个被执行的部分。下面列出查询子句在逻辑上处理顺序：
1\. FROM
2\. WHERE
3\. GROUP BY
4\. HAVING
5\. SELECT
6\. ORDER BY

## select 后修饰符

在 SELECT 关键字之后，您可以使用许多影响语句操作的修饰符。

> HIGH_PRIORITY、STRAIGHT_JOIN 和以 SQL\_ 开头的修饰符是 MySQL 对标准 SQL 的扩展。

ALL 和 DISTINCT 修饰符指定是否应返回重复行。

- **ALL** # **默认值**。指定应返回所有匹配行，包括重复行。
- **DISTINCT** # 指定从结果集中删除重复行。
- **DISTINCTROW** # 是 DISTINCT 的同义词。

注意：同时指定两个修饰符是错误的。

# 应用示例

从 inventory 表中返回所有行的所有列

```sql
MariaDB [caredaily]> select * from inventory;
+----+--------------+------+-----------+
| id | type         | size | inventory |
+----+--------------+------+-----------+
|  1 | 果C精品      | NB   |         1 |
|  2 | IP小飞侠     | NB   |         5 |
|  3 | 果C精品      | NB   |        50 |
|  4 | 新款丝薄     | XXL  |       100 |
+----+--------------+------+-----------+
```

从 inventory 表中返回所有行的 type、size 列

```sql
MariaDB [caredaily]> select type,size from inventory;
+--------------+------+
| type         | size |
+--------------+------+
| 果C精品      | NB   |
| IP小飞侠     | NB   |
| 果C精品      | NB   |
| 新款丝薄     | XXL  |
+--------------+------+
```

- 在 inventory 表中，对 amount 列的数据进行求和操作，条件是 product 列中所有的“新款丝薄”
  - select sum(amount) from inventory where product="新款丝薄";

## 分组编号

若 MySQL 为 8.0+ 版本，可以使用 `row_number()` 函数。

## 多列模糊查询

```sql
SELECT * FROM card_descs WHERE concat(sc_name,effect,evo_cover_effect) LIKE '%奥米加%';
```

从 sc_name,effect,evo_cover_effect 这三列中搜索包含“奥米加”的行
等效于：

```sql
SELECT *
FROM card_descs
WHERE sc_name LIKE '%奥米加%'
    OR effect LIKE '%奥米加%'
    OR evo_cover_effect LIKE '%奥米加%';
```

## 获取某所有值并去重

distinct 简单来说就是用来去重的，而 group by 的设计目的则是用来聚合统计的，两者在能够实现的功能上有些相同之处，但应该仔细区分，因为用错场景的话，效率相差可以倍计。
单纯的去重操作使用 distinct，速度是快于 group by 的。

```sql
SELECT
	distinct level
FROM
	card_descs
```
