---
title: ClickHouse
linkTitle: ClickHouse
weight: 1
---

# 概述

> 参考：
>
> - [GitHub 项目，ClickHouse/ClickHouse](https://github.com/ClickHouse/ClickHouse)
> - [官网](https://clickhouse.com/)

存算分离，查询性能过剩

https://clickhouse.com/docs/en/guides/sre/network-ports

| 端口号   | 描述                                                                                                                               |
| ----- | -------------------------------------------------------------------------------------------------------------------------------- |
| 2181  | ZooKeeper default service port. **Note: see `9181` for ClickHouse Keeper**                                                       |
| 8123  | HTTP default port                                                                                                                |
| 8443  | HTTP SSL/TLS default port                                                                                                        |
| 9000  | 原生协议端口（也称为 ClickHouse TCP 协议）。由 ClickHouse 生态的应用程序和进程使用（e.g. 各种语言利用 SDK 编写的程序、clickhouse-client 等自带程序、etc.）。也用于分布式查询的内部服务器之间的通信。 |
| 9440  | 与 9000 的功能相同，但是带有 SSL/TLS                                                                                                        |
| 9004  | MySQL emulation port                                                                                                             |
| 9005  | PostgreSQL emulation port (also used for secure communication if SSL is enabled for ClickHouse).                                 |
| 9009  | Inter-server communication port for low-level data access. Used for data exchange, replication, and inter-server communication.  |
| 9010  | SSL/TLS for inter-server communications                                                                                          |
| 9011  | Native protocol PROXYv1 protocol port                                                                                            |
| 9019  | JDBC bridge                                                                                                                      |
| 9100  | gRPC port                                                                                                                        |
| 9181  | Recommended ClickHouse Keeper port                                                                                               |
| 9234  | Recommended ClickHouse Keeper Raft port (also used for secure communication if `<secure>1</secure>` enabled)                     |
| 9363  | 在 /metrics 路径下暴露 Prometheus 格式的 Metric 指标                                                                                        |
| 9281  | Recommended Secure SSL ClickHouse Keeper port                                                                                    |
| 42000 | Graphite default port                                                                                                            |

# 学习资料

[B 站 - 蜂蜜柠檬水HLN，带你快速认识 ClickHouse 数据库 | 了解 OLTP 与 OLAP](https://www.bilibili.com/video/BV1vd4FeSErh)

# Engine

**[Engine](/docs/5.数据存储/数据库/关系数据/ClickHouse/Engine.md)(引擎)** 是 ClickHouse 实现数据处理功能的核心抽象。数据库 以及 表 都由各种各样的 Engine 实现

- **Database Engine(数据库引擎)**
- **Table Engine(表引擎)**

# 关联文件与配置

https://clickhouse.com/docs/en/operations/configuration-files

https://clickhouse.com/docs/en/operations/settings

**/etc/clickhouse-server/**

- **./config.xml** # ClickHouse Server 运行配置。
- **./config.d/** # 配置文件可以拆分到该目录，程序运行时会将该目录下的文件合并到 config.xml 主配置文件
- **./metrika.xml** # 默认的 include_from 文件。该文件中的配置用来替换主配置文件 config.xml 中的配置。
  - e.g. config.xml 中有 `<remote_servers incl="clickhouse_remote_server"/>`，那么 metrika.xml 中的 `<clickhouse_remote_servers>` 部分配置就会作为 config.xml 中的 remote_servers。
- **./users.xml** # e.g. 认证信息、etc. 相关配置
- **./users.d/** # 配置文件可以拆分到该目录，程序运行时会将该目录下的文件合并到 users.xml 主配置文件

# ClickHouse 部署

https://clickhouse.com/docs/en/install

# CLI

https://clickhouse.com/docs/en/operations/utilities

## clickhouse-server

## clickhouse-client

https://clickhouse.com/docs/integrations/sql-clients/cli

### OPTIONS

https://clickhouse.com/docs/integrations/sql-clients/cli#command-line-options

连接选项

- **-h, --host**(STRING) # ClickHouse Server 的 IP
- **--port**(INT) # ClickHouse Server 的 PORT
- **-u, --user**(STRING) # 连接数据库所使用的用户名。`默认值: default`
- **--possword**(STRING) #  连接数据库的用户的密码
- **-d, --database**(STRING) # 连接的数据库名称。`默认值: default`

查询选项

- **-m, --multiline** # 允许多行查询（按 Enter 键时不发送查询），仅当查询以分号结尾时才会发送查询。
- **-q, --query**(STRING) # 指定查询 SQL。可以将 SQL 保存到文件中，使用 `--query="$(cat query.sql)"` 这种方式执行查询。
 
### EXAMPLE

```bash
clickhouse-client -u default --password d1234567 -h 127.0.0.1 --port 9000 -d my_database --query="cat tmp/query.sql"
```

# ClickHouse 生态

> 参考：
>
> - [官方文档，集成](https://clickhouse.com/docs/en/integrations)

Grafana 数据源插件 https://github.com/grafana/clickhouse-datasource 。详见 Grafana [Plugins](/docs/6.可观测性/Grafana/Plugins.md)

- 在 https://github.com/grafana/clickhouse-datasource/tree/main/src/dashboards 有一些内置的利用 ClickHouse 本身的数据创建出来的 Grafana 仪表盘
- [官方文档，可观测性 - Grafana](https://clickhouse.com/docs/en/observability/grafana) 有一些最佳实践和示例

https://github.com/clickvisual/clickvisual 一个基于 clickhouse 构建的轻量级日志分析和数据可视化 Web 平台。

https://github.com/metrico/promcasa 通过 ClickHouse 的 SQL，将查询结果转为 OpenMetrics 格式数据。

## 驱动与接口

https://clickhouse.com/docs/en/interfaces/overview

[可视化接口](https://clickhouse.com/docs/en/interfaces/third-party/gui)

# Cluster

> 参考：
>
> - [官方文档，架构 - 水平扩展](https://clickhouse.com/docs/architecture/horizontal-scaling)
> - [B 站，WordScenesTV - 【clickhouse】clickhouse集群架构、部署和使用](https://www.bilibili.com/video/BV1qz421h7BX)

![](https://clickhouse.com/docs/assets/ideal-img/scaling-out-1.3666d1c.600.png)

**Shard** # 数据的分片

**Replica** # 每个分片的副本

**ClickHouseKeeper** # ClickHouse 集群的协调系统，通知 Shard 的副本关于状态变化，使用 RAFT [共识算法](/docs/3.集群与分布式/分布式算法/共识算法.md)实现。ClickHouseKeeper 必须单数节点，最少 3 个来保证选举。

- ClickHouseKeeper 的逻辑也在 ClickHouse 程序的逻辑中，所以可以有两种运行方式
  - 与 ClickHouse 一起运行，作为其内部逻辑
  - 独立运行

---

![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/clickhouse/20250325154658383.png)

比如

通过如下 SQL 可以查看集群的拓扑结构：

```sql
SELECT cluster, shard_num, replica_num, host_name, port
FROM system.clusters;
```

结果像这样

| cluster    | shard_num | replica_num | host_name | port |
| ---------- | --------- | ----------- | --------- | ---- |
| my_cluster | 1         | 1           | host1     | 9000 |
| my_cluster | 1         | 2           | host2     | 9000 |
| my_cluster | 2         | 1           | host3     | 9000 |
| my_cluster | 2         | 2           | host4     | 9000 |

这种结果的配置来源于下面这种配置：

```xml
<remote_servers>
  <my_cluster>
    <shard>
      <replica>
        <host>host1</host>
        <port>9000</port>
      </replica>
      <replica>
        <host>host3</host>
        <port>9000</port>
      </replica>
    </shard>
    <shard>
      <replica>
        <host>host2</host>
        <port>9000</port>
      </replica>
      <replica>
        <host>host4</host>
        <port>9000</port>
      </replica>
    </shard>
  </my_cluster>
</remote_servers>
```

这个集群共两个分片，将数据分别保存在 host1/host3 和 host2/host4 上，每个分片都有一个自己的备份
