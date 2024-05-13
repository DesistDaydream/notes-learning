---
title: MySQL 命令行工具
---

# 概述

> 参考：

# mysql

> 参考：
> - [官方文档，MySQL 程序-客户端程序-mysql](https://dev.mysql.com/doc/refman/8.0/en/mysql.html)

mysql 是一个简单的 SQL Shell。 它支持交互和非交互使用。 交互使用时，查询结果以 ASCII 表格式显示。 非交互使用（例如，用作过滤器）时，结果以制表符分隔的格式显示。 可以使用命令选项更改输出格式。

## Syntax(语法)

**mysql \[OPTIONS] \[DATABASE]**
**DATABASE** # 指定连接 mysql 后要操作的数据库。若不指定，则需要在交互模式下使用 `use` 指令选择数据库，否则对数据库的操作将会报 `No database selected` 错误：

```bash
mysql> show tables;
ERROR 1046 (3D000): No database selected
```

**OPTIONS：**

- **-h, --host \<HostName>** # 指定要连接的 mysql 主机。如果链接本机 mysql，可以省略。
- **-P, --port \<PORT>** # 指定要连接的 mysql 的端口。默认值：`3306`
- **-u, --user \<UserName>** # 指定要登录 mysql 的用户名
- **-p, --password \<PASSWORD>** # 使用密码来登录。如果指定要登录 mysql 的用户密码为空，则该选项可省

## 命令行模式

我们可以通过 `mysql db_name <FILE.sql > output.tab` 命令直接执行写在文件中的 SQL 语句

## 交互模式

### 斜线命令

> 参考：
> - [官方文档，MySQL 程序-客户端程序-myslq 客户端命令](https://dev.mysql.com/doc/refman/8.0/en/mysql-commands.html)

在 mysql 的交互模式中有一组 mysql 程序自带的命令，用以 控制输出格式、检查、获取数据信息 等等，这些命令以 `\` 开头，不过也有与之相对应的字符串命令

- **\u, use \<DBName>** # 选择想要操作的数据库。与 MySQL 的 SQL 中的 use 语句功能一致

### 基础示例

- grant select,insert,update,delete,create,drop ON mysql.\* TO 'desistdaydream'@'localhost' identified by 'desistdaydream'; # 为名为 mysql 的数据库创建名为 desistdaydream 的用户，密码为 desistdaydream，具有 select、insert、update、delete、create、drop 这些命令的执行权限。
- flush privileges; # 刷新权限。由权限账号信息是在 MYSQLD 服务启动的时候就加载到内存中的，所以你在原权限表中的任何直接修改都不会直接生效。用 flush privileges 把中表中的信息更新到内存。
- select user(); # 查看当前登录的用户。
- show databases; # 列出所有已经存在的数据库
- use mysql; # 切换当前要操作的数据库为 mysql
- show tables; # 显示当前数据库中所有的表
- show columns from db; # 显示当前数据库中名为 db 的表的属性。效果如下
  - desc test; # 与该命令效果相同
  - Field # 该表中都有哪些列
  - Type # 该列的数据类型
  - Null # 该列是否可以插入 null
  - Key # 索引类型
  - Default # 该列插入空值时。默认插入什么值。
  - Extra # 该列额外的参数。

```sql
MariaDB [mysql]> SHOW COLUMNS FROM db;
+-----------------------+---------------+------+-----+---------+-------+
| Field                 | Type          | Null | Key | Default | Extra |
+-----------------------+---------------+------+-----+---------+-------+
| Host                  | char(60)      | NO   | PRI |         |       |
| Db                    | char(64)      | NO   | PRI |         |       |
| User                  | char(16)      | NO   | PRI |         |       |
| Select_priv           | enum('N','Y') | NO   |     | N       |       |
.......
```

- select Host,db from db; # 显示 db 表中，Host 和 Db 列及其内容，效果如下

```sql
MariaDB [mysql]> SELECT Host,db from db;
+-----------+---------+
| Host      | db      |
+-----------+---------+
| %         | test    |
| %         | test\_% |
| localhost | mysql   |
+-----------+---------+
```

# mysqladmin

EXAMPLE

- mysqladmin -u root -p password "my_password" # 修改 root 密码，密码为：my_password。如果默认密码为空，则可以不加-p。
