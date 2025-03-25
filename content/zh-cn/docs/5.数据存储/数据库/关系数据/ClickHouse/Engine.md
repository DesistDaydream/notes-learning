---
title: Engine
linkTitle: Engine
weight: 20
date: 2025-03-25T12:47:00
---

# 概述

> 参考：
>
> - [官方文档，SQL 参考 - 引擎](https://clickhouse.com/docs/en/engines)
> - [流式数据同步：一种PostgreSQL到ClickHouse的高效数据同步方案](https://juejin.cn/post/7375275474006016011)

**Engine(引擎)** 是 ClickHouse 实现数据处理功能的核心抽象。数据库 以及 表 都由各种各样的 Engine 实现

- **Database Engine(数据库引擎)**
- **Table Engine(表引擎)**

# Database Engine


# Table Engine

Table Engine(表引擎) 本质上是用来定义表的类型。ClickHouse 的表甚至可以通过 Engine 从其他数据库中读取数据（e.g. 直接读取 PostgreSQL 中某个表的数据）

Table Engine 可以决定：

- How and where data is stored, where to write it to, and where to read it from. 数据如何存储、在何处存储、将其写入何处以及从何处读取。
- Which queries are supported, and how. 支持哪些查询以及如何支持。
- Concurrent data access. 并发数据访问。
- Use of indexes, if present. 使用索引（如果存在）。
- Whether multithread request execution is possible. 是否可以执行多线程请求。
- Data replication parameters. 数据复制参数。

Table Engine 分为如下几大类：

- **MergeTree Family** # MergeTree 相关引擎。<font color="#ff0000">这些引擎是 ClickHouse 数据存储功能的核心，也是最常用的表引擎</font>。提供了大多数的功能：列式存储、定制分区、稀疏主索引、etc.
- **Log Family** # 日志相关引擎。这些引擎都是具有最小功能的轻量级引擎。当需要快速写入许多小表（最多约 100 万行）并在以后整体读取时，最有效。
- **Integrations** # 集成引擎。用于与其他的数据存储与处理系统集成的引擎。e.g. Kafka, MySQL, ODBC, JDBC, HDFC, etc.
- **Special** # 其他特定功能的引擎

e.g. 创建一个 PostgreSQL 引擎的表：

```sql
CREATE TABLE my_database.my_table (
  `id` UInt64,
  `version` String,
  `command_id` String,
)
ENGINE = PostgreSQL('10.53.192.45:5432', 'PG_Database', 'PG_Table', 'PG_Username', 'PG_Password', 'CH_ClusterName')
```

创建完成后，我们在 CH 中查询的 my_database.my_table 表中的数据实际上是直接获取的 PostgreSQL 中的 PG_Database.PG_Table 表的数据。

## MergeTree Family

### MergeTree

## Special

### Distributed

Distributed(分布式) 表常用在集群模式的 ClickHouse 中。**Distribuited 表本身不存储任何数据**，而是提供了一个接口来访问集群中多个分片中的数据。 

当查询 Distributed 表时，Distributed 表会将查询转发给所有主机，等待来自每个分片的查询结果，然后计算并返回整个查询结果。

当写入 Distributed 表时，可以自定定义或让引擎自动分布在各个分片中。

下面是创建 Distributed 表的基本语法。创建一个 my_distributed 表，该表相当于汇总了 my_cluster 集群中所有节点中的 my_database.my_table 的数据：

```sql
CREATE TABLE [IF NOT EXISTS] [db.]my_distributed [ON CLUSTER cluster]
(
    name1 [type1] [DEFAULT|MATERIALIZED|ALIAS expr1],
    name2 [type2] [DEFAULT|MATERIALIZED|ALIAS expr2],
    ...
) ENGINE = Distributed(my_cluster, my_database, my_table[, sharding_key[, policy_name]])
[SETTINGS name=value, ...]
```

