---
title: MySQL SQL
linkTitle: MySQL SQL
weight: 20
---

# 概述

> 参考：
>
> - [MySQL 官方文档，SLQ 语句](https://dev.mysql.com/doc/refman/8.0/en/sql-statements.html)

SQL 语言被细分为几个语言元素，包括：

- 子句，它们是语句和查询的组成部分。（在某些情况下，这些是可选的。）
- 表达式，可以生成[标](<https://en.wikipedia.org/wiki/Scalar_(computing)>)量值或由数据的[列](<https://en.wikipedia.org/wiki/Column_(database)>)和[行](<https://en.wikipedia.org/wiki/Row_(database)>)组成的[表](<https://en.wikipedia.org/wiki/Table_(database)>)
- Predicates，指定可以评估为 SQL[三值逻辑 (3VL)](https://en.wikipedia.org/wiki/Ternary_logic)（真/假/未知）或[布尔](https://en.wikipedia.org/wiki/Boolean_logic) [真值](https://en.wikipedia.org/wiki/Truth_value)的条件，用于限制语句和查询的效果，或更改程序流程。
- 查询，根据特定条件检索数据。这是 SQL 的一个重要元素。
- 语句，可能对模式和数据产生持久影响，或者可能控制[事务](https://en.wikipedia.org/wiki/Database_transaction)、程序流、连接、会话或诊断。
  - SQL 语句还包括[分号](https://en.wikipedia.org/wiki/Semicolon)(";") 语句终止符。虽然不是每个平台都需要它，但它被定义为 SQL 语法的标准部分。
- [SQL 语句和查询中通常会忽略无关紧要的空格](<https://en.wikipedia.org/wiki/Whitespace_(computer_science)>)，从而更容易格式化 SQL 代码以提高可读性。

在 MySQL 中按照功能将各种语句进行了分类：

- [Data Definition Statements](https://dev.mysql.com/doc/refman/8.0/en/sql-data-definition-statements.html)(数据定义语句)
- [Data Manipulation Statements]()(数据操作语句)
- [Transactional and Locking Statements](https://dev.mysql.com/doc/refman/8.0/en/sql-transactional-statements.html)(事务和锁语句)
- [Replication Statements](https://dev.mysql.com/doc/refman/8.0/en/sql-replication-statements.html)(复制语句)
- [Prepared Statements](https://dev.mysql.com/doc/refman/8.0/en/sql-prepared-statements.html)(预处理语句)
- [Compound Statement Syntax](https://dev.mysql.com/doc/refman/8.0/en/sql-compound-statements.html)(符合语句)
- [Database Administration Statements](https://dev.mysql.com/doc/refman/8.0/en/sql-server-administration-statements.html)(数据库管理语句)
- [Utility Statements](https://dev.mysql.com/doc/refman/8.0/en/sql-utility-statements.html)(使用程序语句)

# 关键字

> 参考：
>
> - [MySQL 官方文档，语言结构 - 关键字和保留字](https://dev.mysql.com/doc/refman/8.0/en/keywords.html)

# 函数与运算符

> 参考：
>
> - [MySQL 官方文档，函数与运算符](https://dev.mysql.com/doc/refman/5.7/en/functions.html)

# 内置函数

> 参考：
>
> - [MySQL 官方文档，函数与运算符 - 内置函数和运算符参考](https://dev.mysql.com/doc/refman/8.0/en/built-in-function-reference.html)

| 函数名                                                                                                      | 功能                                             | 启用版本 | 弃用版本 |
| ----------------------------------------------------------------------------------------------------------- | ------------------------------------------------ | -------- | -------- |
| [REPLACE()](https://dev.mysql.com/doc/refman/8.0/en/string-functions.html#function_replace)                 | 替换掉指定字符串                                 |          |          |
| [REGEXP_REPLACE()](https://dev.mysql.com/doc/refman/8.0/en/regexp.html#function_regexp-replace)             | 替换掉使用正则表达式匹配到的字符串               | 8.0.4    |          |
| [SUBSTRING_INDEX()](https://dev.mysql.com/doc/refman/8.0/en/string-functions.html#function_substring-index) | 从指定出现次数的分隔符之前的字符串中返回子字符串 |          |          |
|                                                                                                             |                                                  |          |          |
|                                                                                                             |                                                  |          |          |

# 聚合函数

> 参考：
>
> - [官方文档，函数和运算符 - 聚合函数](https://dev.mysql.com/doc/refman/5.7/en/aggregate-functions-and-modifiers.html)

Aggregate Function(聚合函数) 通常与 `group by` 关键字一起使用，用于将值进行分组。

| **Name**                                                                                                     | **Description**                                  | **Introduced** |
| ------------------------------------------------------------------------------------------------------------ | ------------------------------------------------ | -------------- |
| [AVG()](https://dev.mysql.com/doc/refman/5.7/en/aggregate-functions.html#function_avg)                       | Return the average value of the argument         |                |
| [BIT_AND()](https://dev.mysql.com/doc/refman/5.7/en/aggregate-functions.html#function_bit-and)               | Return bitwise AND                               |                |
| [BIT_OR()](https://dev.mysql.com/doc/refman/5.7/en/aggregate-functions.html#function_bit-or)                 | Return bitwise OR                                |                |
| [BIT_XOR()](https://dev.mysql.com/doc/refman/5.7/en/aggregate-functions.html#function_bit-xor)               | Return bitwise XOR                               |                |
| [COUNT()](https://dev.mysql.com/doc/refman/5.7/en/aggregate-functions.html#function_count)                   | Return a count of the number of rows returned    |                |
| [COUNT(DISTINCT)](https://dev.mysql.com/doc/refman/5.7/en/aggregate-functions.html#function_count-distinct)  | Return the count of a number of different values |                |
| [GROUP_CONCAT()](https://dev.mysql.com/doc/refman/5.7/en/aggregate-functions.html#function_group-concat)     | Return a concatenated string                     |                |
| [JSON_ARRAYAGG()](https://dev.mysql.com/doc/refman/5.7/en/aggregate-functions.html#function_json-arrayagg)   | Return result set as a single JSON array         | 5.7.22         |
| [JSON_OBJECTAGG()](https://dev.mysql.com/doc/refman/5.7/en/aggregate-functions.html#function_json-objectagg) | Return result set as a single JSON object        | 5.7.22         |
| [MAX()](https://dev.mysql.com/doc/refman/5.7/en/aggregate-functions.html#function_max)                       | Return the maximum value                         |                |
| [MIN()](https://dev.mysql.com/doc/refman/5.7/en/aggregate-functions.html#function_min)                       | Return the minimum value                         |                |
| [STD()](https://dev.mysql.com/doc/refman/5.7/en/aggregate-functions.html#function_std)                       | Return the population standard deviation         |                |
| [STDDEV()](https://dev.mysql.com/doc/refman/5.7/en/aggregate-functions.html#function_stddev)                 | Return the population standard deviation         |                |
| [STDDEV_POP()](https://dev.mysql.com/doc/refman/5.7/en/aggregate-functions.html#function_stddev-pop)         | Return the population standard deviation         |                |
| [STDDEV_SAMP()](https://dev.mysql.com/doc/refman/5.7/en/aggregate-functions.html#function_stddev-samp)       | Return the sample standard deviation             |                |
| [SUM()](https://dev.mysql.com/doc/refman/5.7/en/aggregate-functions.html#function_sum)                       | Return the sum                                   |                |
| [VAR_POP()](https://dev.mysql.com/doc/refman/5.7/en/aggregate-functions.html#function_var-pop)               | Return the population standard variance          |                |
| [VAR_SAMP()](https://dev.mysql.com/doc/refman/5.7/en/aggregate-functions.html#function_var-samp)             | Return the sample variance                       |                |
| [VARIANCE()](https://dev.mysql.com/doc/refman/5.7/en/aggregate-functions.html#function_variance)             | Return the population standard variance          |                |

## GROUP BY

> 参考：
>
> - [官方文档，函数与运算符 - 聚合函数 - MySQL 对 GROUP BY 的处理](https://dev.mysql.com/doc/refman/5.7/en/group-by-handling.html)
> - [mysql 查询先 having 再 group by](https://codeantenna.com/a/UeE9vTxxdO)

`group by` 可以脱离聚合函数单独使用，此时不在 group by 分组的列将会取其中第一行的值。
如果想要对数据进行分组后取到最新的值，还需要使用 `row_number()` 函数(8.0+版本)，或者某些分组编号功能(此时编号为 1 的就是每组最新的值)。

### HAVING

having 与 group by 之间的用法，暂时找不到一个可以描述清楚的文章，但是事实是，having 可以单独使用，并在外层再套一个 select 后执行 group by 以实现分组后取最新值的效果。

```sql
select
    *
from
    (
        SELECT

        FROM
            t_diversion diversion
        WHERE
            diversion.is_deleted = 0
        HAVING
            1
        ORDER BY
            diversion.updated_at desc
    ) t
GROUP BY
    t.datalakeName,
    t.privateNetName
```

此时可以取到 diversion 表中，按照 datalakeName 和 privateNetName 分组后，每组中最新的一条数据。


# 运算符

> 参考：
>
> - [官方文档，函数与运算符 - 运算符](https://dev.mysql.com/doc/refman/8.0/en/non-typed-operators.html)

| **运算符**                                                                                                        | **描述**                                                                                                                                                                                     | **Introduced**                                                                | **弃用版本** |
| ----------------------------------------------------------------------------------------------------------------- | -------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ----------------------------------------------------------------------------- | ------------ | --- | --- |
| [&](https://dev.mysql.com/doc/refman/8.0/en/bit-functions.html#operator_bitwise-and)                              | Bitwise AND                                                                                                                                                                                  |                                                                               |              |
| [>](https://dev.mysql.com/doc/refman/8.0/en/comparison-operators.html#operator_greater-than)                      | Greater than operator                                                                                                                                                                        |                                                                               |              |
| [>>](https://dev.mysql.com/doc/refman/8.0/en/bit-functions.html#operator_right-shift)                             | Right shift                                                                                                                                                                                  |                                                                               |              |
| [>=](https://dev.mysql.com/doc/refman/8.0/en/comparison-operators.html#operator_greater-than-or-equal)            | Greater than or equal operator                                                                                                                                                               |                                                                               |              |
| [<](https://dev.mysql.com/doc/refman/8.0/en/comparison-operators.html#operator_less-than)                         | Less than operator                                                                                                                                                                           |                                                                               |              |
| [<>,!=](https://dev.mysql.com/doc/refman/8.0/en/comparison-operators.html#operator_not-equal)                     | Not equal operator                                                                                                                                                                           |                                                                               |              |
| [<<](https://dev.mysql.com/doc/refman/8.0/en/bit-functions.html#operator_left-shift)                              | Left shift                                                                                                                                                                                   |                                                                               |              |
| [<=](https://dev.mysql.com/doc/refman/8.0/en/comparison-operators.html#operator_less-than-or-equal)               | Less than or equal operator                                                                                                                                                                  |                                                                               |              |
| [<=>](https://dev.mysql.com/doc/refman/8.0/en/comparison-operators.html#operator_equal-to)                        | NULL-safe equal to operator                                                                                                                                                                  |                                                                               |              |
| [%,MOD](https://dev.mysql.com/doc/refman/8.0/en/arithmetic-functions.html#operator_mod)                           | Modulo operator                                                                                                                                                                              |                                                                               |              |
| [\*](https://dev.mysql.com/doc/refman/8.0/en/arithmetic-functions.html#operator_times)                            | Multiplication operator                                                                                                                                                                      |                                                                               |              |
| [+](https://dev.mysql.com/doc/refman/8.0/en/arithmetic-functions.html#operator_plus)                              | Addition operator                                                                                                                                                                            |                                                                               |              |
| [-](https://dev.mysql.com/doc/refman/8.0/en/arithmetic-functions.html#operator_minus)                             | Minus operator                                                                                                                                                                               |                                                                               |              |
| [-](https://dev.mysql.com/doc/refman/8.0/en/arithmetic-functions.html#operator_unary-minus)                       | Change the sign of the argument                                                                                                                                                              |                                                                               |              |
| [->](https://dev.mysql.com/doc/refman/8.0/en/json-search-functions.html#operator_json-column-path)                | Return value from JSON column after evaluating path; equivalent to JSON_EXTRACT().                                                                                                           |                                                                               |              |
| [->>](https://dev.mysql.com/doc/refman/8.0/en/json-search-functions.html#operator_json-inline-path)               | Return value from JSON column after evaluating path and unquoting the result; equivalent to JSON_UNQUOTE(JSON_EXTRACT()).                                                                    |                                                                               |              |
| [/](https://dev.mysql.com/doc/refman/8.0/en/arithmetic-functions.html#operator_divide)                            | Division operator                                                                                                                                                                            |                                                                               |              |
| [:=](https://dev.mysql.com/doc/refman/8.0/en/assignment-operators.html#operator_assign-value)                     | 变量赋值                                                                                                                                                                                     |                                                                               |              |
| [=](https://dev.mysql.com/doc/refman/8.0/en/assignment-operators.html#operator_assign-equal)                      | 变量赋值（作为[ SET](https://dev.mysql.com/doc/refman/8.0/en/set-variable.html) 语句的一部分，或作为 [UPDATE](https://dev.mysql.com/doc/refman/8.0/en/update.html) 语句中 SET 子句的一部分） |                                                                               |              |
| [=](https://dev.mysql.com/doc/refman/8.0/en/comparison-operators.html#operator_equal)                             | Equal operator                                                                                                                                                                               |                                                                               |              |
| [^](https://dev.mysql.com/doc/refman/8.0/en/bit-functions.html#operator_bitwise-xor)                              | Bitwise XOR                                                                                                                                                                                  |                                                                               |              |
| [AND,&&](https://dev.mysql.com/doc/refman/8.0/en/logical-operators.html#operator_and)                             | Logical AND                                                                                                                                                                                  |                                                                               |              |
| [BETWEEN ... AND ...](https://dev.mysql.com/doc/refman/8.0/en/comparison-operators.html#operator_between)         | Whether a value is within a range of values                                                                                                                                                  |                                                                               |              |
| [BINARY](https://dev.mysql.com/doc/refman/8.0/en/cast-functions.html#operator_binary)                             | Cast a string to a binary string                                                                                                                                                             |                                                                               | 8.0.27       |
| [CASE](https://dev.mysql.com/doc/refman/8.0/en/flow-control-functions.html#operator_case)                         | Case operator                                                                                                                                                                                |                                                                               |              |
| [DIV](https://dev.mysql.com/doc/refman/8.0/en/arithmetic-functions.html#operator_div)                             | Integer division                                                                                                                                                                             |                                                                               |              |
| [IN()](https://dev.mysql.com/doc/refman/8.0/en/comparison-operators.html#operator_in)                             | Whether a value is within a set of values                                                                                                                                                    |                                                                               |              |
| [IS](https://dev.mysql.com/doc/refman/8.0/en/comparison-operators.html#operator_is)                               | Test a value against a boolean                                                                                                                                                               |                                                                               |              |
| [IS NOT](https://dev.mysql.com/doc/refman/8.0/en/comparison-operators.html#operator_is-not)                       | Test a value against a boolean                                                                                                                                                               |                                                                               |              |
| [IS NOT NULL](https://dev.mysql.com/doc/refman/8.0/en/comparison-operators.html#operator_is-not-null)             | NOT NULL value test                                                                                                                                                                          |                                                                               |              |
| [IS NULL](https://dev.mysql.com/doc/refman/8.0/en/comparison-operators.html#operator_is-null)                     | NULL value test                                                                                                                                                                              |                                                                               |              |
| [LIKE](https://dev.mysql.com/doc/refman/8.0/en/string-comparison-functions.html#operator_like)                    | Simple pattern matching                                                                                                                                                                      |                                                                               |              |
| [MEMBER OF()](https://dev.mysql.com/doc/refman/8.0/en/json-search-functions.html#operator_member-of)              | Returns true (1) if first operand matches any element of JSON array passed as second operand, otherwise returns false (0)                                                                    | 8.0.17                                                                        |              |
| [NOT,!](https://dev.mysql.com/doc/refman/8.0/en/logical-operators.html#operator_not)                              | Negates value                                                                                                                                                                                |                                                                               |              |
| [NOT BETWEEN ... AND ...](https://dev.mysql.com/doc/refman/8.0/en/comparison-operators.html#operator_not-between) | Whether a value is not within a range of values                                                                                                                                              |                                                                               |              |
| [NOT IN()](https://dev.mysql.com/doc/refman/8.0/en/comparison-operators.html#operator_not-in)                     | Whether a value is not within a set of values                                                                                                                                                |                                                                               |              |
| [NOT LIKE](https://dev.mysql.com/doc/refman/8.0/en/string-comparison-functions.html#operator_not-like)            | Negation of simple pattern matching                                                                                                                                                          |                                                                               |              |
| [NOT REGEXP](https://dev.mysql.com/doc/refman/8.0/en/regexp.html#operator_not-regexp)                             | Negation of REGEXP                                                                                                                                                                           |                                                                               |              |
| [OR,                                                                                                              |                                                                                                                                                                                              | ](https://dev.mysql.com/doc/refman/8.0/en/logical-operators.html#operator_or) | Logical OR   |     |     |
| [REGEXP](https://dev.mysql.com/doc/refman/8.0/en/regexp.html#operator_regexp)                                     | Whether string matches regular expression                                                                                                                                                    |                                                                               |              |
| [RLIKE](https://dev.mysql.com/doc/refman/8.0/en/regexp.html#operator_regexp)                                      | Whether string matches regular expression                                                                                                                                                    |                                                                               |              |
| [SOUNDS LIKE](https://dev.mysql.com/doc/refman/8.0/en/string-functions.html#operator_sounds-like)                 | Compare sounds                                                                                                                                                                               |                                                                               |              |
| [XOR](https://dev.mysql.com/doc/refman/8.0/en/logical-operators.html#operator_xor)                                | Logical XOR                                                                                                                                                                                  |                                                                               |              |
| [                                                                                                                 | ](https://dev.mysql.com/doc/refman/8.0/en/bit-functions.html#operator_bitwise-or)                                                                                                            | Bitwise OR                                                                    |              |     |
| [~](https://dev.mysql.com/doc/refman/8.0/en/bit-functions.html#operator_bitwise-invert)                           | Bitwise inversion                                                                                                                                                                            |                                                                               |              |

# 数据库管理语句

## show - 显示信息

显示有关数据库、表、列的信息，或有关服务器状态的信息。

### EXAMPLE

- 查看创建 test 数据库的 sql 语句
  - show create database test;
- 查看创建 test 表的 sql 语句
  - show create table test;
- 显示数据库的状态信息
  - show status;

# 实用程序语句

## use - 使用指定数据库作为后续 SQL 执行的目标库

# 特殊符号

## @

声明变量、调用变量
