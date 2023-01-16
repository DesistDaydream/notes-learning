---
title: Ceph 命令行工具
---

# 概述

> 参考：
> - [官方文档,Ceph 存储集群-手册页](https://docs.ceph.com/en/latest/rados/man/)
> - [官方文档,Ceph 对象网关-radosgw-admin 手册页](https://docs.ceph.com/en/latest/man/8/radosgw-admin/#)

# ceph # Ceph 管理工具

一个 Python 实现的脚本工具，用于手动部署和似乎 Ceph 集群。通过很多的子命令，允许部署 ceph-mon、ceph-osd、PG、ceph-mds 等，并可以对集群整体进行维护和管理。

## orch

Orchestrator(编排器，简称 orch)

### Syntax(语法)

COMMAND

- **host **# 对集群中的节点进行管理
  - **add <HOSTNAME> \[ADDR] \[LABELs...] \[--maintenance]** # 向集群中添加一个节点
  - **label add <HOSTNAME> <LABEL>** # 为节点添加一个标签
- **ls** # 列出 Orch 已知的服务
- **rm <ServiceName> **# 移除一个服务

EXAMPLE

# radosgw-admin # RADOS 网关的用户管理工具

radosgw-admin 是一个 RADOS 网关用户的管理工具。可以增删改查用户。该工具通过非常多的子命令进行管理，并且每个子命令可用的选项也各不相同，Ceph 官方对这个工具的提示做的非常不好，子命令需要带的选项并不提示，只能自己尝试~~~

## user

### Syntax(语法)

**radosgw-admin user COMMAND \[OPTIONS]**

COMMAND

- **user create --display-name=<STRING> --uid=<STRING>** # 创建一个新用户
- **user info \[--uid=<STRING> | --access-key=<STRING>]** # 显示一个用户的信息，包括其子用户和密钥。通过 uid 或 ak 指定要显示的用户。
- **user list** # 列出所有用户
- **user modify --uid=<STRING>** # 修改指定的用户

OPTIONS

- **--admin** # 为指定用户设定 admin 标志
- **--display-name=<STRING>** # 为指定用户设定对外显示的名称
- **--email=<STRING>** # 为用户设定邮箱
- **--uid=<STRING>** # 指定用户的 ID。在执行绝大部分与用户相关的命令时，都需要指定该选项，以确定操作的用户。比如 查看用户信息、查看属于指定用户的桶的信息 等等等等
- **--system** # 为指定用户设定 system 标志

### EXAPMEL

创建一个名为 lichenhao 的用户，并添加 system 标志

- **radosgw-admin user create --uid=lichenhao --display-name=lichenhao --system**

<br />
## bucket
### Syntax(语法)
**radosgw-admin bucket COMMAND [OPTIONS]**

COMMAND

- **bucket stats \[OPTIONS] **# 显示桶的统计信息。可以通过选项指定用户下的桶或指定的桶。

OPTIONS

- **--bucket=<STRING>** # 指定桶的名称。可以被 quota 子命令使用。
- **--uid=<STRING>** # 指定用户的 ID。查看桶信息时，将会显示该用户下所有桶的信息。
