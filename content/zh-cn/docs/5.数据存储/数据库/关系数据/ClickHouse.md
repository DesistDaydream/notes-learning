---
title: ClickHouse
linkTitle: ClickHouse
date: 2024-09-30T15:29
weight: 20
---

# 概述

> 参考：
>
> - [GitHub 项目，ClickHouse/ClickHouse](https://github.com/ClickHouse/ClickHouse)
> - [官网](https://clickhouse.com/)


https://clickhouse.com/docs/en/guides/sre/network-ports

| 端口号   | 描述                                                                                                                                                                                                                                                         |
| ----- | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| 2181  | ZooKeeper default service port. **Note: see `9181` for ClickHouse Keeper**                                                                                                                                                                                 |
| 8123  | HTTP default port                                                                                                                                                                                                                                          |
| 8443  | HTTP SSL/TLS default port                                                                                                                                                                                                                                  |
| 9000  | Native Protocol port (also referred to as ClickHouse TCP protocol). Used by ClickHouse applications and processes like `clickhouse-server`, `clickhouse-client`, and native ClickHouse tools. Used for inter-server communication for distributed queries. |
| 9004  | MySQL emulation port                                                                                                                                                                                                                                       |
| 9005  | PostgreSQL emulation port (also used for secure communication if SSL is enabled for ClickHouse).                                                                                                                                                           |
| 9009  | Inter-server communication port for low-level data access. Used for data exchange, replication, and inter-server communication.                                                                                                                            |
| 9010  | SSL/TLS for inter-server communications                                                                                                                                                                                                                    |
| 9011  | Native protocol PROXYv1 protocol port                                                                                                                                                                                                                      |
| 9019  | JDBC bridge                                                                                                                                                                                                                                                |
| 9100  | gRPC port                                                                                                                                                                                                                                                  |
| 9181  | Recommended ClickHouse Keeper port                                                                                                                                                                                                                         |
| 9234  | Recommended ClickHouse Keeper Raft port (also used for secure communication if `<secure>1</secure>` enabled)                                                                                                                                               |
| 9363  | 在 /metrics 路径下暴露 Prometheus 格式的 Metric 指标                                                                                                                                                                                                                  |
| 9281  | Recommended Secure SSL ClickHouse Keeper port                                                                                                                                                                                                              |
| 9440  | Native protocol SSL/TLS port                                                                                                                                                                                                                               |
| 42000 | Graphite default port                                                                                                                                                                                                                                      |

### CLI

https://clickhouse.com/docs/en/operations/utilities

**clickhouse-server**

**clickhouse-client**

https://clickhouse.com/docs/en/integrations/sql-clients/cli

- clickhouse-client -u default --password 12345678  -m -n --port 9100 -h 127.0.0.1 -d network_security


# 其他

Grafana 数据源插件 https://github.com/grafana/clickhouse-datasource

https://github.com/clickvisual/clickvisual 一个基于 clickhouse 构建的轻量级日志分析和数据可视化 Web 平台。

https://github.com/metrico/promcasa 通过 ClickHouse 的 SQL，将查询结果转为 OpenMetrics 格式数据。

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

`arrayExists(\[func,] arr1, ...)` # func 是一个高阶函数，可以接受 lambda 函数

- e.g. `arrayExists(x -> x = src_ip, [${example_array}])` # `${example_array}` # 变量是一个数组，作为参数传递给 x，arrayExists 将会逐一一对 x 中的元素，执行 `x = src_ip` 表达式。主要用于判断 变量中的元素是否等于 src_ip

### 日期与时间

https://clickhouse.com/docs/en/sql-reference/functions/date-time-functions

toStartOfXXX

类似 [PostgreSQL](docs/5.数据存储/数据库/关系数据/PostgreSQL/PostgreSQL.md) 的 date_bin 函数，将日期和时间向下舍入到 XXX，主要是用来依据时间进行数据聚合以实现时间序列数据的效果，并可在 Grafana 中绘制时间序列图表。

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
