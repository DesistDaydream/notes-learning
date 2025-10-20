---
title: PostgreSQL
linkTitle: PostgreSQL
weight: 1
---

# 概述

> 参考：
>
> - [官网](https://www.postgresql.org/)
> - [PostgreSQL 新手入门-阮一峰](http://www.ruanyifeng.com/blog/2013/12/getting_started_with_postgresql.html)

PostgreSQL 是一个功能强大的开源对象关系数据库系统，经过 30 多年的积极开发，在可靠性、特性健壮性和性能方面赢得了很高的声誉。

一个 Database(库) 中包含多个 Schemas(模式)，一个 Schema(模式) 中包含多个 Tables(表)

## Schema

https://www.postgresql.org/docs/current/ddl-schemas.html

PostgreSQL 的数据库中包含 1 个或多个 Schema，所有的 Table 是归属在 Schema 下的。

> 可以讲 Schema 理解为 Namespace（PostgreSQL 也是通过 pg_namespace 元表（元数据表）查看所有 Schema）

默认情况下，创建的 Table 自动放入名为 **public** 的 Schema 下。每个数据库都会包含 public Schema。

要访问非 public Schema 下的 Table，使用 `.` 符号。e.g. `SchemaName.TableName`，如果用最简单的 SQL 举例就是:  `select * from schema_demo.table_one` 列出名为 schema_demo 模式中的 table_one 表下的所有列。

PostgreSQL 内置了如下几个 Schemas

- **public** # 在不指定 Schema 的情况下，新建的 Table 都默认保存在 public Schema 中。
- **pg_catalog** # System catalogs(系统目录)，保存 PostgreSQL 运行常见的
- **information_schema** # 与 Schema 相关的内部信息

# PostgreSQL 部署

## Redhat 包部署

https://www.postgresql.org/download/linux/redhat/

```bash
yum install postgresql-server postgresql -y
```

> 除了 postgresl 客户端和服务端以外，还有两个包可以安装 postgresql-contrib（额外提供的模块）、postgresql-devel（C 语言开发的库和头文件）

初始化数据库，为 /var/lib/pgsql/data/ 目录填充数据，若目录为空则无法启动 postgresql。

```bash
postgresql-setup --initdb
```

启动 PostgreSQL 服务端

```bash
systemctl enable postgresql.service --now
```

## Debian 包部署

https://www.postgresql.org/download/linux/ubuntu/

**一、安装**

首先，安装 PostgreSQL 客户端。

`sudo apt-get install postgresql-client`

然后，安装 PostgreSQL 服务器。

`sudo apt-get install postgresql`

正常情况下，安装完成后，PostgreSQL 服务器会自动在本机的 5432 端口开启。

如果还想安装图形管理界面，可以运行下面命令，但是本文不涉及这方面内容。

`sudo apt-get install pgadmin3`

# PostgreSQL 关联文件与配置

**pg_hba.conf** # 控制如何访问以及哪些可以访问 PgSQL Server

- https://www.postgresql.org/docs/current/auth-pg-hba-conf.html

**postgresql.conf** # 可以改监听地址

# 元数据

## System catalogs

> 参考：
>
> - [官方文档，内部 - 51. 系统目录](https://www.postgresql.org/docs/current/catalogs.html)

**System catalogs(系统目录)** 是关系数据库管理系统存储模式元数据的地方，例如有关表和列的信息以及内部簿记信息。 PostgreSQL 的 <font color="#ff0000">System catalogs 是常规表</font>。您可以删除并重新创建表、添加列、插入和更新值，并以这种方式严重扰乱您的系统。通常，不应手动更改系统目录，通常有 SQL 命令可以做到这一点。 （例如，CREATE DATABASE 会在 pg_database 目录中插入一行，并实际上在磁盘上创建数据库。）对于特别深奥的操作有一些例外，但随着时间的推移，其中许多操作已作为 SQL 命令提供，因此需要对系统目录的直接操作正在不断减少。

这些 System catalogs 常规表默认保存在 `pg_catalog` Schema 中，还可以通过 `\dS`

| Catalog 名称                                                                                                | 用途                                                  |
| --------------------------------------------------------------------------------------------------------- | --------------------------------------------------- |
| [`pg_namespace`](https://www.postgresql.org/docs/current/catalog-pg-namespace.html "51.32. pg_namespace") | 记录 Schems 的基本元信息。包含 oid, nspname, nspowner,nspacl 列 |
| TODO                                                                                                      |                                                     |

比如 `SELECT * FROM pg_namespace;` 可以查看所有 Schemas 的信息。

## Information Schema

> 参考：
>
> - [官方文档，客户端接口](https://www.postgresql.org/docs/current/information-schema.html)

Information Schema 由一组视图组成，这些视图包含有关当前数据库中定义的对象的信息。Information Schema 是在 SQL 标准中定义的，因此可以预期是可移植的并保持稳定，与 System catalogs 不同，Information Schema 特定于 PostgreSQL 并且根据实现问题进行建模。然而，Information Schema 视图不包含有关 PostgreSQL 特定功能的信息；要查询这些信息，您需要查询 System catalogs 或其他 PostgreSQL 特定的视图。

Information Schema 有一个名为 information_schema 的 Schema。该模式自动存在于所有数据库中。该模式的所有者是集群中的初始数据库用户，该用户自然拥有该模式的所有权限，包括删除它的能力（但由此节省的空间微乎其微）。

默认情况下，information_schema 不在模式搜索路径中，因此需要通过限定名称访问其中的所有对象。由于 information_schema 中某些对象的名称是用户应用程序中可能出现的通用名称，因此如果要将信息模式放入路径中，则应小心。

> [!Tip]
> 在 information_schema.schemate 查看 Schema 信息时，有一列名为 catalog_name，可以从 `SELECT datname FROM pg_catalog.pg_database;` 获取到，这 catalog_name 就是类似 Database(数据库) 的概念。相当于在 psql 中执行 `\l`

## 查询性能

下面内容来自 Claude AI 回答：

PostgreSQL 确实有查询统计信息的系统视图。我来帮你创建类似的查询分析面板。

```sql
-- 1. 查询最耗时的 SQL 语句
SELECT
    queries.query,
    calls as executions,
    round(total_exec_time::numeric, 2) as total_time_ms,
    round(mean_exec_time::numeric, 2) as avg_time_ms,
    round(stddev_exec_time::numeric, 2) as stddev_time_ms,
    round((100 * total_exec_time / sum(total_exec_time) over ())::numeric, 2) as percentage_cpu,
    rows as total_rows
FROM pg_stat_statements queries
ORDER BY total_exec_time DESC
LIMIT 10;

-- 2. 查询执行次数最多的 SQL
SELECT
    queries.query,
    calls as executions,
    round(total_exec_time::numeric, 2) as total_time_ms,
    round(mean_exec_time::numeric, 2) as avg_time_ms,
    round((total_exec_time/calls)::numeric, 2) as avg_time_per_call_ms
FROM pg_stat_statements queries
ORDER BY calls DESC
LIMIT 10;

-- 3. 按小时统计查询数量（需要结合 pg_stat_statements 的重置时间）
WITH RECURSIVE hours AS (
    SELECT
        date_trunc('hour', now()) as hour
    UNION ALL
    SELECT
        hour - interval '1 hour'
    FROM hours
    WHERE hour > date_trunc('hour', now() - interval '24 hour')
)
SELECT
    hours.hour,
    count(*) as query_count
FROM hours
LEFT JOIN pg_stat_statements stats
    ON date_trunc('hour', stats.last_call) = hours.hour
GROUP BY hours.hour
ORDER BY hour DESC;

-- 4. 查询平均执行时间分布
SELECT
    CASE
        WHEN mean_exec_time < 1 THEN '<1ms'
        WHEN mean_exec_time < 10 THEN '1-10ms'
        WHEN mean_exec_time < 100 THEN '10-100ms'
        WHEN mean_exec_time < 1000 THEN '100ms-1s'
        ELSE '>1s'
    END as execution_time_bucket,
    count(*) as query_count
FROM pg_stat_statements
GROUP BY 1
ORDER BY
    CASE execution_time_bucket
        WHEN '<1ms' THEN 1
        WHEN '1-10ms' THEN 2
        WHEN '10-100ms' THEN 3
        WHEN '100ms-1s' THEN 4
        ELSE 5
    END;

```

要使用这些查询，你需要先确保启用了 `pg_stat_statements` 扩展。在 PostgreSQL 中执行：

```sql
CREATE EXTENSION IF NOT EXISTS pg_stat_statements;
```

然后在 postgresql.conf 中添加：

```
shared_preload_libraries = 'pg_stat_statements'
pg_stat_statements.track = all
```

在 Grafana 中，你可以创建以下面板：

1. **最耗时查询 Top 10（Table Panel）**
   - 使用第一个查询
   - 显示查询文本、执行次数、总时间、平均时间和行数
   - 可以按总时间降序排序

2. **查询执行次数 Top 10（Table Panel）**
   - 使用第二个查询
   - 显示最频繁执行的查询及其性能指标

3. **查询数量时间趋势（Time Series Panel）**
   - 使用第三个查询
   - 显示每小时查询数量的变化趋势
   - X轴设置为时间，Y轴为查询数量

4. **查询耗时分布（Pie Chart or Bar Chart）**
   - 使用第四个查询
   - 展示不同执行时间范围的查询数量分布

你还可以添加以下监控指标：

```sql
-- 5. 数据库活跃连接数
SELECT
    datname,
    count(*) as active_connections
FROM pg_stat_activity
GROUP BY datname;

-- 6. 缓存命中率
SELECT
    sum(heap_blks_hit) * 100.0 / (sum(heap_blks_hit) + sum(heap_blks_read)) as cache_hit_ratio
FROM pg_statio_user_tables;
```

需要注意的是：

1. pg_stat_statements 的数据会持续累积，需要定期重置：

   ```sql
   SELECT pg_stat_statements_reset();
   ```

2. 有些查询可能包含敏感信息，建议在生产环境中谨慎展示完整的查询文本
3. 可以根据需要调整查询的时间范围和限制条件

# GUI 工具

[数据库管理工具](/docs/5.数据存储/数据管理工具/数据库管理工具.md)
