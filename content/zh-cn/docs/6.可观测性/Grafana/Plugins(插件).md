---
title: Plugins(插件)
linkTitle: Plugins(插件)
date: 2024-09-30T15:31
weight: 20
---

# 概述

> 参考：
>
> -


# Clickhouse

> 参考：
>
> - [GitHub 项目，grafana/clickhouse-datasource](https://github.com/grafana/clickhouse-datasource)
> - https://grafana.com/grafana/plugins/grafana-clickhouse-datasource/

[ClickHouse](docs/5.数据存储/数据库/关系数据/ClickHouse/ClickHouse.md) 插件

## Macros

https://github.com/grafana/clickhouse-datasource?tab=readme-ov-file#macros

Grafana 的 ClickHouse 插件会改变一些 Grafana 中某些变量、函数的用法；还有一些新增加的 Macros 功能可以使用

`$__timeFilter(COLUMN_NAME)` # 

| Macro                       | Description                                                                                                      | 渲染结果示例                                                              |
| --------------------------- | ---------------------------------------------------------------------------------------------------------------- | ------------------------------------------------------------------- |
| `$__dateFilter(ColumnName)` | Replaced by a conditional that filters the data (using the provided column) based on the date range of the panel | `date >= toDate('2022-10-21') AND date <= toDate('2022-10-23')`     |
| `$__timeFilter(ColumnName)` | 利用 toDateTime() 函数转换 Grafana 面板时间范围选择器选择时间范围的值                                                                   | `time >= toDateTime(1415792726) AND time <= toDateTime(1447328726)` |
| etc.                        |                                                                                                                  |                                                                     |
