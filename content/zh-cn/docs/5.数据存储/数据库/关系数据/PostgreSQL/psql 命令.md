---
title: psql 命令
linkTitle: psql 命令
weight: 20
---


# 概述

> 参考：
>
> - [官方文档](https://www.postgresql.org/docs/current/app-psql.html)

psql 是 PostgreSQL 的交互式终端，可以看作是 PostgreSQL 的 [REPL](/docs/2.编程/Programming%20environment/REPL.md)。

psql 中可以执行有多种类型的命令

- [Meta-Commands](#meta-commands)
- SQL

# Syntax(语法)

**psql \[OPTION] \[DBNAME \[USERNAME]]**

**OPTIONS**

- 连接数据库相关选项
  - **-U, --username USERNAME** # 使用指定的用户连接数据库。`默认值: 当前 Shell 环境的用户`。
  - **-h, --host HOSTNAME** # 指定 PostgreSQL 服务端所在的 HOSTNAME，可以是 IP 或 Domain。如果该以 `\` 开头，则将其用作 Unix 域套接字的目录。
  - **-p, --port PORT** # 指定 PostgreSQL 服务端监听的 TCP 端口或本地 Unix 域套接字文件扩展名。默认为 PGPORT 环境变量的值，`默认值: 5432`。
- SQL 执行相关
  - **-t, --tuples-only** # 关闭列名和结果行计数页脚等的打印。这相当于 `\t` 或 `\pset tuples_only`

# Meta-Commands

```sql
\h：查看SQL命令的解释，比如 \h select。
\?：查看psql命令列表。
\l：列出所有数据库。
\c [database_name]：连接其他数据库。
\d：列出当前数据库的所有表格。
\d [table_name]：列出某一张表格的结构。
\du：列出所有 roles(角色)。由于 “用户” 和 “组” 的概念已统一为 “角色”，因此该命令现在相当于 \dg。
\e：打开文本编辑器。
\conninfo：列出当前数据库和连接的信息。
\password USER # 为指定的 USER 设置密码
```

# 最佳实践

使用 pgadmin 用户连接本地 127.0.0.1 且监听在 5432 端口上的 PgSQL Server 中的 postgres 数据库

- `psql postgres -U pgadmin -h 127.0.0.1 -p 5432`

热加载配置文件

- `psql -U postgres -h 127.0.0.1 -c "SELECT pg_reload_conf();"`
- `/usr/local/pgsql/bin/psql -U pgadmin -d olp_euintf -h 127.0.0.1 -c "SELECT pg_reload_conf();"`
