---
title: PostgreSQL
linkTitle: PostgreSQL
date: 2023-11-11T17:30:00
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

## 部署
### Redhat 包部署

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

### Debian 包部署

https://www.postgresql.org/download/linux/ubuntu/

**一、安装**

首先，安装 PostgreSQL 客户端。

`sudo apt-get install postgresql-client`

然后，安装 PostgreSQL 服务器。

`sudo apt-get install postgresql`

正常情况下，安装完成后，PostgreSQL 服务器会自动在本机的 5432 端口开启。

如果还想安装图形管理界面，可以运行下面命令，但是本文不涉及这方面内容。

`sudo apt-get install pgadmin3`

## 为 postgres 用户添加密码

```bash
~]# su - postgres
~]$ psql
Password for user postgres:
psql (13.3)
Type "help" for help.

postgres=# \password postgres
Enter new password:
Enter it again:
postgres=#
```

## 修改配置

https://gist.github.com/AtulKsol/4470d377b448e56468baef85af7fd614

默认情况下 psql 使用对等身份验证通过 UNIX 套接字进行连接，这要求当前 UNIX 用户具有与 psql 相同的用户名。因此，您必须创建 UNIX 用户 postgres，然后以 postgres 身份登录或使用 sudo -u postgres psql 数据库名称 来访问数据库（并且 psql 不应要求输入密码）。若使用 `1` 这种命令连接 PostgreSQL，将会又如下报错：

`psql: FATAL: Peer authentication failed for user “postgres” (or any user)`

但如果打算通过 Unix 套接字而不是对等方法强制进行密码身份验证，修改 pg_hba.conf 配置文件中的如下内容：

```bash
# TYPE DATABASE USER ADDRESS METHOD
local  all      all          peer
```

改为

```bash
# TYPE DATABASE USER ADDRESS METHOD
local  all      all          md5
```

其中 METHOD 可以有三个值：

- `peer` 意味着它将信任 UNIX 用户的身份（真实性）。所以不要求密码。
- `md5` 意味着它始终会要求输入密码，并在使用 MD5 哈希后进行验证。
- `trust`意味着它永远不会要求输入密码，并且始终信任任何连接。

### 修改监听以及允许通过 TCP 连接

```bash
# 修改 postgresql.conf 配置文件，添加监听地址，改为*`
listen_addresses = '*'

# 修改 pg_hba.conf，添加远程主机地址，放在第一行：允许任意用户从任意机器上以密码方式访问数据库，把下行添加为第一条规则：
host    all             all             0.0.0.0/0              md5
```

## 连接数据库

```bash
~]# psql -d postgres -U postgres
Password for user postgres:
psql (13.3)
Type "help" for help.

postgres=#
```

## 添加新用户和新数据库

初次安装后，默认生成一个名为 postgres 的数据库和一个名为 postgres 的数据库用户。这里需要注意的是，同时还生成了一个名为 postgres 的 Linux 系统用户。

下面，我们使用 postgres 用户，来生成其他用户和新数据库。好几种方法可以达到这个目的，这里介绍两种。

**第一种方法，使用 PostgreSQL 控制台。**

首先，新建一个 Linux 新用户，可以取你想要的名字，这里为 dbuser。

> sudo adduser dbuser

然后，切换到 postgres 用户。

> sudo su - postgres

下一步，使用 psql 命令登录 PostgreSQL 控制台。

> psql

这时相当于系统用户 postgres 以同名数据库用户的身份，登录数据库，这是不用输入密码的。如果一切正常，系统提示符会变为"postgres=#"，表示这时已经进入了数据库控制台。以下的命令都在控制台内完成。

第一件事是使用\password 命令，为 postgres 用户设置一个密码。

> \password postgres

第二件事是创建数据库用户 dbuser（刚才创建的是 Linux 系统用户），并设置密码。

> CREATE USER dbuser WITH PASSWORD 'password';

第三件事是创建用户数据库，这里为 exampledb，并指定所有者为 dbuser。

> CREATE DATABASE exampledb OWNER dbuser;

第四件事是将 exampledb 数据库的所有权限都赋予 dbuser，否则 dbuser 只能登录控制台，没有任何数据库操作权限。

> GRANT ALL PRIVILEGES ON DATABASE exampledb to dbuser;

最后，使用\q 命令退出控制台（也可以直接按 ctrl+D）。

> \q

**第二种方法，使用 shell 命令行。**

添加新用户和新数据库，除了在 PostgreSQL 控制台内，还可以在 shell 命令行下完成。这是因为 PostgreSQL 提供了命令行程序 createuser 和 createdb。还是以新建用户 dbuser 和数据库 exampledb 为例。

首先，创建数据库用户 dbuser，并指定其为超级用户。

```bash
sudo -u postgres createuser --superuser dbuser
```

然后，登录数据库控制台，设置 dbuser 用户的密码，完成后退出控制台。

```bash
sudo -u postgres psql
\password dbuser
\q
```

接着，在 shell 命令行下，创建数据库 exampledb，并指定所有者为 dbuser。

```bash
sudo -u postgres createdb -O dbuser exampledb
```

## 登录数据库

添加新用户和新数据库以后，就要以新用户的名义登录数据库，这时使用的是 psql 命令。

> psql -U dbuser -d exampledb -h 127.0.0.1 -p 5432

上面命令的参数含义如下：-U 指定用户，-d 指定数据库，-h 指定服务器，-p 指定端口。

输入上面命令以后，系统会提示输入 dbuser 用户的密码。输入正确，就可以登录控制台了。

psql 命令存在简写形式。如果当前 Linux 系统用户，同时也是 PostgreSQL 用户，则可以省略用户名（-U 参数的部分）。举例来说，我的 Linux 系统用户名为 ruanyf，且 PostgreSQL 数据库存在同名用户，则我以 ruanyf 身份登录 Linux 系统后，可以直接使用下面的命令登录数据库，且不需要密码。

> psql exampledb

此时，如果 PostgreSQL 内部还存在与当前系统用户同名的数据库，则连数据库名都可以省略。比如，假定存在一个叫做 ruanyf 的数据库，则直接键入 psql 就可以登录该数据库。

> psql

另外，如果要恢复外部数据，可以使用下面的命令。

> psql exampledb exampledb.sql

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

# GUI 工具

[数据库管理工具](/docs/5.数据存储/数据管理工具/数据库管理工具.md)

# SQL

TODO: 整理特定于 PostgreSQL 的 SQL

```plsql
# 创建新表
CREATE TABLE user_tbl(name VARCHAR(20), signup_date DATE);

# 插入数据
INSERT INTO user_tbl(name, signup_date) VALUES('张三', '2013-12-22');

# 从表中查询数据
SELECT * FROM user_tbl;

# 更新数据
UPDATE user_tbl set name = '李四' WHERE name = '张三';

# 删除记录
DELETE FROM user_tbl WHERE name = '李四' ;

# 添加栏位
ALTER TABLE user_tbl ADD email VARCHAR(40);

# 更新结构
ALTER TABLE user_tbl ALTER COLUMN signup_date SET NOT NULL;

# 更名栏位
ALTER TABLE user_tbl RENAME COLUMN signup_date TO signup;

# 删除栏位
ALTER TABLE user_tbl DROP COLUMN email;

# 表格更名
ALTER TABLE user_tbl RENAME TO backup_tbl;

# 删除表格
DROP TABLE IF EXISTS backup_tbl;
```

## Schema 查询

列出所有 Schema

```sql
SELECT *
FROM pg_namespace;
```

或

```sql
SELECT schema_name
FROM information_schema.schemata;
```

列出名为 cheat 的 Schema 下的所有表

```sql
SELECT table_name
FROM information_schema.tables
WHERE table_schema = 'cheat'
  AND table_type = 'BASE TABLE';
```

