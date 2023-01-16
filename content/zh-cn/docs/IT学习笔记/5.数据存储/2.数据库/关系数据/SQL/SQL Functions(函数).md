---
title: SQL Functions(函数)
---

# 概述

> 参考：
> - [MySQL 官方文档，函数与运算符](https://dev.mysql.com/doc/refman/5.7/en/functions.html)

# 内置函数

> 参考：
> - [MySQL 官方文档，函数与运算符-内置函数和运算符参考](https://dev.mysql.com/doc/refman/8.0/en/built-in-function-reference.html)

| 函数名                                                                                                      | 功能                                             | 启用版本 | 弃用版本 |
| ----------------------------------------------------------------------------------------------------------- | ------------------------------------------------ | -------- | -------- |
| [REPLACE()](https://dev.mysql.com/doc/refman/8.0/en/string-functions.html#function_replace)                 | 替换掉指定字符串                                 |          |          |
| [REGEXP_REPLACE()](https://dev.mysql.com/doc/refman/8.0/en/regexp.html#function_regexp-replace)             | 替换掉使用正则表达式匹配到的字符串               | 8.0.4    |          |
| [SUBSTRING_INDEX()](https://dev.mysql.com/doc/refman/8.0/en/string-functions.html#function_substring-index) | 从指定出现次数的分隔符之前的字符串中返回子字符串 |          |          |
|                                                                                                             |                                                  |          |          |
|                                                                                                             |                                                  |          |          |

# 聚合函数

> 参考：
> - [官方文档，函数和运算符-聚合函数](https://dev.mysql.com/doc/refman/5.7/en/aggregate-functions-and-modifiers.html)

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
> - [官方文档，函数与运算符-聚合函数-MySQL 对 GROUP BY 的处理](https://dev.mysql.com/doc/refman/5.7/en/group-by-handling.html)
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
