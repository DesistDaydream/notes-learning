---
title: "PostgreSQL MGMT"
linkTitle: "PostgreSQL MGMT"
weight: 20
---

# 概述

> 参考：
>
> - 

# PostgreSQL 部署后常见操作

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

```bash
sudo adduser dbuser
```

然后，切换到 postgres 用户。

```bash
sudo su - postgres
```

下一步，使用 psql 命令登录 PostgreSQL 控制台。

```bash
psql
```

这时相当于系统用户 postgres 以同名数据库用户的身份，登录数据库，这是不用输入密码的。如果一切正常，系统提示符会变为"postgres=#"，表示这时已经进入了数据库控制台。以下的命令都在控制台内完成。

第一件事是使用\password 命令，为 postgres 用户设置一个密码。

```bash
\password postgres
```

第二件事是创建数据库用户 dbuser（刚才创建的是 Linux 系统用户），并设置密码。

```sql
CREATE USER dbuser WITH PASSWORD 'password';
```

第三件事是创建用户数据库，这里为 exampledb，并指定所有者为 dbuser。

```sql
CREATE DATABASE exampledb OWNER dbuser;
```

第四件事是将 exampledb 数据库的所有权限都赋予 dbuser，否则 dbuser 只能登录控制台，没有任何数据库操作权限。

```sql
GRANT ALL PRIVILEGES ON DATABASE exampledb to dbuser;
```

最后，使用\q 命令退出控制台（也可以直接按 ctrl+D）。

```bash
\q
```

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

## 添加只读用户

```sql
-- 1. 创建用户
CREATE USER developer WITH PASSWORD 'developer';
-- 2. 授予连接数据库的权限
GRANT CONNECT ON DATABASE network_security TO developer;
-- 3. 授予 public Schema 的使用权限
GRANT USAGE ON SCHEMA public TO developer;
-- 4. 授予所有表的 public Schema 的只读权限（SELECT）
GRANT SELECT ON ALL TABLES IN SCHEMA public TO developer;
-- 5. 授予所有序列的 public Schema 的只读权限（可以查看序列信息，但不能修改）
-- TODO: 这一步并不是必须的，暂时没找到必须要加这个的原因
-- GRANT SELECT ON ALL SEQUENCES IN SCHEMA public TO developer;
```
