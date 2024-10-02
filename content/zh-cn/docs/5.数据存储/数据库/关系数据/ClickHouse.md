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