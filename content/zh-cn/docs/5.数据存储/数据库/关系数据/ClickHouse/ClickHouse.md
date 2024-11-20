---
title: ClickHouse
linkTitle: ClickHouse
date: 2024-09-30T15:29
weight: 1
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

# Engine

> 参考：
>
> - [官方文档，SQL 参考 - 引擎](https://clickhouse.com/docs/en/engines)
> - [流式数据同步：一种PostgreSQL到ClickHouse的高效数据同步方案](https://juejin.cn/post/7375275474006016011)

- Database Engine(数据库引擎)
- Table Engine(表引擎)

## Database Engine


## Table Engine

Table Engine(表引擎) 本质上是用来定义表的类型。ClickHouse 的表甚至可以通过 Engine 从其他数据库中读取数据（e.g. 直接读取 PostgreSQL 中某个表的数据）

比如用下面找个创建 Table 的语法举例：

```sql
CREATE TABLE my_database.my_table (
  `id` UInt64, 
  `command_source` Nullable(Int64), 
  `source_system` String, 
  `version` String, 
  `command_id` String,
)
ENGINE = PostgreSQL('10.53.192.45:5432', 'PG_Database', 'PG_Table', 'PG_Username', 'PG_Password', 'CH_ClusterName')
```

创建完成后，我们在 CH 中查询的 my_database.my_table 表中的数据实际上是直接获取的 PostgreSQL 中的 PG_Database.PG_Table 表的数据。

Table Engine 可以决定：

- How and where data is stored, where to write it to, and where to read it from.数据如何存储、在何处存储、将其写入何处以及从何处读取。
- Which queries are supported, and how.支持哪些查询以及如何支持。
- Concurrent data access.并发数据访问。
- Use of indexes, if present.使用索引（如果存在）。
- Whether multithread request execution is possible.是否可以执行多线程请求。
- Data replication parameters.数据复制参数。

# ClickHouse 生态

Grafana 数据源插件 https://github.com/grafana/clickhouse-datasource

https://github.com/clickvisual/clickvisual 一个基于 clickhouse 构建的轻量级日志分析和数据可视化 Web 平台。

https://github.com/metrico/promcasa 通过 ClickHouse 的 SQL，将查询结果转为 OpenMetrics 格式数据。

