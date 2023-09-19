---
title: SQL
---

# 概述

> 参考：
>
> - [Wiki,SQL](https://en.wikipedia.org/wiki/SQL)
> - [菜鸟教程，SQL](https://www.runoob.com/sql/sql-tutorial.html)

**Structured Query Language(结构化查询语言，简称 SQL)** 是一种特定领域的编程语言，用于管理 RDBMS(关系数据库管理系统) 中保存的数据。使用 SQL 编写的语句也可以称为 **Expression(表达式)**。

SQL 在 1986 年成为 [ANSI](docs/x_标准化术语/计算机标准/ANSI.md) 的一项标准，在 1987 年成为国际标准化组织（ISO）标准。

每种关系型数据库所使用的 SQL 基本都一样，但是又有其自身特殊的 SQL。由于 MySQL 的使用率非常高，所以 SQL 文档的各种例子都以 MySQL 为主。

## 关键字

> 参考：
>
> - [MySQL 官方文档，语言结构-关键字和保留字](https://dev.mysql.com/doc/refman/8.0/en/keywords.html)

SELECT
FROM
ORDER BY
GROUP BY
WHERE

# SQL 标准

| Year | Name                                               | Alias                                                                               | Comments                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                       |
| ---- | -------------------------------------------------- | ----------------------------------------------------------------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| 1986 | SQL-86                                             | SQL-87                                                                              | First formalized by ANSI                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                       |
| 1989 | SQL-89                                             | [FIPS](https://en.wikipedia.org/wiki/Federal_Information_Processing_Standard) 127-1 | Minor revision that added integrity constraints adopted as FIPS 127-1                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                          |
| 1992 | [SQL-92](https://en.wikipedia.org/wiki/SQL-92)     | SQL2, FIPS 127-2                                                                    | Major revision (ISO 9075), _Entry Level_ SQL-92 adopted as FIPS 127-2                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                          |
| 1999 | [SQL:1999](https://en.wikipedia.org/wiki/SQL:1999) | SQL3                                                                                | Added regular expression matching, [recursive queries](https://en.wikipedia.org/wiki/Hierarchical_and_recursive_queries_in_SQL) (e.g. [transitive closure](https://en.wikipedia.org/wiki/Transitive_closure)), [triggers](https://en.wikipedia.org/wiki/Database_trigger), support for procedural and control-of-flow statements, nonscalar types (arrays), and some object-oriented features (e.g. [structured types](https://en.wikipedia.org/wiki/Structured_type)), support for embedding SQL in Java ([SQL/OLB](https://en.wikipedia.org/wiki/SQL/OLB)) and vice versa ([SQL/JRT](https://en.wikipedia.org/wiki/SQL/JRT)) |
| 2003 | [SQL:2003](https://en.wikipedia.org/wiki/SQL:2003) |                                                                                     | Introduced [XML](https://en.wikipedia.org/wiki/XML)-related features ([SQL/XML](https://en.wikipedia.org/wiki/SQL/XML)), [window functions](https://en.wikipedia.org/wiki/SQL_window_function), standardized sequences, and columns with autogenerated values (including identity columns)                                                                                                                                                                                                                                                                                                                                     |
| 2006 | [SQL:2006](https://en.wikipedia.org/wiki/SQL:2006) |                                                                                     | ISO/IEC 9075-14:2006 defines ways that SQL can be used with XML. It defines ways of importing and storing XML data in an SQL database, manipulating it within the database, and publishing both XML and conventional SQL data in XML form. In addition, it lets applications integrate queries into their SQL code with [XQuery](https://en.wikipedia.org/wiki/XQuery), the XML Query Language published by the World Wide Web Consortium ([W3C](https://en.wikipedia.org/wiki/W3C)), to concurrently access ordinary SQL-data and XML documents.[\[33\]](https://en.wikipedia.org/wiki/SQL#cite_note-SQLXML2006-36)           |
| 2008 | [SQL:2008](https://en.wikipedia.org/wiki/SQL:2008) |                                                                                     | Legalizes ORDER BY outside cursor definitions. Adds INSTEAD OF triggers, TRUNCATE statement,[\[34\]](https://en.wikipedia.org/wiki/SQL#cite_note-iablog.sybase.com-paulley-37) FETCH clause                                                                                                                                                                                                                                                                                                                                                                                                                                    |
| 2011 | [SQL:2011](https://en.wikipedia.org/wiki/SQL:2011) |                                                                                     | Adds temporal data (PERIOD FOR)[\[35\]](https://en.wikipedia.org/wiki/SQL#cite_note-feature_temporal-38) (more information at [Temporal database#History](https://en.wikipedia.org/wiki/Temporal_database#History)). Enhancements for [window functions](https://en.wikipedia.org/wiki/SQL_window_function) and FETCH clause.[\[36\]](https://en.wikipedia.org/wiki/SQL#cite_note-features_2011-39)                                                                                                                                                                                                                            |
| 2016 | [SQL:2016](https://en.wikipedia.org/wiki/SQL:2016) |                                                                                     | Adds row pattern matching, polymorphic table functions, [JSON](https://en.wikipedia.org/wiki/JSON)                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                             |
| 2019 | SQL:2019                                           |                                                                                     | Adds Part 15, multidimensional arrays (MDarray type and operators)                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                             |

# SQL 语句

> 参考：
>
> - [MySQL 官方文档，SLQ 语句](https://dev.mysql.com/doc/refman/8.0/en/sql-statements.html)

SQL 语言被细分为几个语言元素，包括：

- 子句，它们是语句和查询的组成部分。（在某些情况下，这些是可选的。）
- 表达式，可以生成[标](<https://en.wikipedia.org/wiki/Scalar_(computing)>)量值或由数据的[列](<https://en.wikipedia.org/wiki/Column_(database)>)和[行](<https://en.wikipedia.org/wiki/Row_(database)>)组成的[表](<https://en.wikipedia.org/wiki/Table_(database)>)
- Predicates，指定可以评估为 SQL[三值逻辑 (3VL)](https://en.wikipedia.org/wiki/Ternary_logic)（真/假/未知）或[布尔](https://en.wikipedia.org/wiki/Boolean_logic) [真值](https://en.wikipedia.org/wiki/Truth_value)的条件，用于限制语句和查询的效果，或更改程序流程。
- 查询，根据特定条件检索数据。这是 SQL 的一个重要元素。
- 语句，可能对模式和数据产生持久影响，或者可能控制[事务](https://en.wikipedia.org/wiki/Database_transaction)、程序流、连接、会话或诊断。
  - SQL 语句还包括[分号](https://en.wikipedia.org/wiki/Semicolon)(";") 语句终止符。虽然不是每个平台都需要它，但它被定义为 SQL 语法的标准部分。
- [SQL 语句和查询中通常会忽略无关紧要的空格](<https://en.wikipedia.org/wiki/Whitespace_(computer_science)>)，从而更容易格式化 SQL 代码以提高可读性。

在 MySQL 中按照功能将各种语句进行了分类：

- [Data Definition Statements](https://dev.mysql.com/doc/refman/8.0/en/sql-data-definition-statements.html)(数据定义语句)
- [Data Manipulation Statements](/docs/5.数据存储/2.数据库/关系数据/SQL/数据操作语句/数据操作语句.md)(数据操作语句)
- [Transactional and Locking Statements](https://dev.mysql.com/doc/refman/8.0/en/sql-transactional-statements.html)(事务和锁语句)
- [Replication Statements](https://dev.mysql.com/doc/refman/8.0/en/sql-replication-statements.html)(复制语句)
- [Prepared Statements](https://dev.mysql.com/doc/refman/8.0/en/sql-prepared-statements.html)(预处理语句)
- [Compound Statement Syntax](https://dev.mysql.com/doc/refman/8.0/en/sql-compound-statements.html)(符合语句)
- [Database Administration Statements](https://dev.mysql.com/doc/refman/8.0/en/sql-server-administration-statements.html)(数据库管理语句)
- [Utility Statements](https://dev.mysql.com/doc/refman/8.0/en/sql-utility-statements.html)(使用程序语句)

## 基础示例

显示当前存在哪些数据库

```sql
mysql> show databases;
+--------------------+
| Database           |
+--------------------+
| information_schema |
| mysql              |
| performance_schema |
| sys                |
+--------------------+
4 rows in set (0.00 sec)
```

> 通常 MySQL 部署成功后，都有几个默认的数据库
>
> - information_schema
> - mysql
> - performance_schema
> - sys

创建数据库

```sql
mysql> create database menagerie;
Query OK, 1 row affected (0.00 sec)
```

开始使用数据库

```sql
mysql> use menagerie
Database changed
```

创建表

```sql
mysql> create table pet (name VARCHAR(20), owner VARCHAR(20),species VARCHAR(20), sex CHAR(1), birth DATE, death DATE);
Query OK, 0 rows affected (0.01 sec)
```

显示表、查看表信息

```sql
mysql> show tables;
+---------------------+
| Tables_in_menagerie |
+---------------------+
| pet                 |
+---------------------+
1 row in set (0.00 sec)

mysql> describe pet;
+---------+-------------+------+-----+---------+-------+
| Field   | Type        | Null | Key | Default | Extra |
+---------+-------------+------+-----+---------+-------+
| name    | varchar(20) | YES  |     | NULL    |       |
| owner   | varchar(20) | YES  |     | NULL    |       |
| species | varchar(20) | YES  |     | NULL    |       |
| sex     | char(1)     | YES  |     | NULL    |       |
| birth   | date        | YES  |     | NULL    |       |
| death   | date        | YES  |     | NULL    |       |
+---------+-------------+------+-----+---------+-------+
6 rows in set (0.00 sec)
```

将数据加载到表中

```sql
-- 直接从文件中将数据加载表中
mysql> LOAD DATA LOCAL INFILE '/pet.txt' INTO TABLE pet;

-- 使用 SQL 语句一条一条插入数据到表中
INSERT INTO pet VALUES ('Puffball','Diane','hamster','f','1999-03-30',NULL);
```

查询数据

```sql
mysql> select * from pet;
+----------+--------+---------+------+------------+------------+
| name     | owner  | species | sex  | birth      | death      |
+----------+--------+---------+------+------------+------------+
| name     | owner  | species | s    | 0000-00-00 | 0000-00-00 |
| Fluffy   | Harold | cat     | f    | 1993-02-04 | 0000-00-00 |
| Claws    | Gwen   | cat     | m    | 1994-03-17 | 0000-00-00 |
| Buffy    | Harold | dog     | f    | 1989-05-13 | 0000-00-00 |
| Fang     | Benny  | dog     | m    | 1990-08-27 | 0000-00-00 |
| Bowser   | Diane  | dog     | m    | 1979-08-31 | 1995-07-29 |
| Chirpy   | Gwen   | bird    | f    | 1998-09-11 | 0000-00-00 |
| Whistler | Gwen   | bird    |      | 1997-12-09 | 0000-00-00 |
| Slim     | Benny  | snake   | m    | 1996-04-29 | 0000-00-00 |
+----------+--------+---------+------+------------+------------+
9 rows in set (0.00 sec)
```

# 数据定义语句

## create # 创建数据库或者表

**create database NAME \[ARGS];**
**create table NAME (Column_Name1 Column_Type1 \[ARGS],....,Column_NameN Column_TypeN \[ARGS])\[ARGS];**
其中 Column_Type 可以使用的参数，涉及每种关系型数据库的各自实现方式

### ARGS

适用于创建 database 时使用的 ARGS：

- character set # 指定字符集为 CharacterName。定义了字符以及字符的编码。
- collation # 指定字符序为 CollationName。定义了字符的培训规则。

适用于创建 table 时使用的 ARGS：
在 () 外设置的 ARGS 将对所有列生效，也可以在 () 内对指定的列设置 ARGS

- null | not null # 指定该列是否可以插入 null 值。默认为 yes，可以插入。一般情况使用设置为 not null,原因见下面说明。
- default \[VALUE] # 指定该列在插入数据为空时的默认值。默认插入 NULL。
    - Note:如果当前列不设定 default 的 VALUE ，在插入数据时，如果不指定列的值。则会根据列的 null 或者 not null 来插入值
    - 当 null 为 yes 时，默认插入 null
    - 当 null 为 no 时，默认根据当前列的类型插入值,对于数值类型插入 0，字符串类型插入空字符串，时间戳类型插入当前日期和时间，ENUM 类型插入枚举组的第一条。
    - e.g.当设置列为 not null、default 不指定 VALUE 时。在插入一个空值时，会报错。因为插入空值，会根据 default 的规则插入 null，但是又不能插入 null，所以插入失败
- comment # 指定该列的注释
- key | primary key | unique key | foreign key <(Column_Name)> # 指定 Column_Name 为该表的索引
- primary key # 关键字用于定义列为主键。 您可以使用多列来定义主键，列间以逗号分隔。
- engine # 设置存储引擎
- charset=ENCODE # 设置编码为 ENCODE。一般设定为 utf8
- collate

适用于 int 类型数据的参数

- auto_increment # 指定该字段自动递增。该参数只适用于 int 类型。每个表有且只能有一个自动递增列，且必须将该列定义为 key
- unsigned # 指定该列的数值无符号，i.e.不会为负数

### EXAMPLE

- 创建一个名为 caredaily 的数据库
  - create database caredaily set utf8 collate utf8_general_ci;
- 创建名为 practice 的数据库，并且指定字符集为 utf8，字符序为 utf8_general_ci。
  - create database `practice` character set utf8 collate utf8_general_ci;
- 创建一个名为 test 的表，其中只有名为 id 的列，类型是 int。
  - create table `test` (`id` int);
- 创建一个名为 product 的数据表。第一列名为 id 的 int 类型 ，无符号，自动递增。第二列名为 size 的 varchar(64) 类型。设定 id 列为主键。自动生成默认值，引擎为 innodb。编码为 utf8

```sql
create table `product` (
`id` int unsigned not null auto_increment,
`product` varchar(64) not null,
primary key (`id`)
)engine=innodb default charset=utf8;
```

- 创建一个名为 size 的表。

```sql
create table `size` (
`id` int unsigned not null auto_increment,
`size` varchar(64) not null,
primary key (`id`)
)engine=innodb default charset=utf8;
```

- 创建一个名为 inventory 的表。第四列名为 create_time 的 时间戳 类型

```sql
create table `inventory` (
`id` int unsigned not null auto_increment,
`src` varchar(64) not null default "总部",
`product` varchar(64) not null,
`size` varchar(64) not null,
`amount` int not null,
`create_time` timestamp not null,
primary key (`id`)
)engine=innodb default charset=utf8;
```

## drop # 删除数据库或者表

**drop {database|table} NAME;**

## truncate # 清空数据

EXAMPLE

- truncate table inventory; # 清空 emp 表，这个命令删除上万条记录特别快，因为他不记录日志

## alter # 修改数据库或者数据表

**alter database NAME;** # 修改数据库的信息
**alter table NAME;** # 修改数据表的信息
EXAMPLE

- alter table inventory add column `create_time` timestamp; # 修改 inventory 表，添加名为 create_time 的列
- alter table inventory drop column `create_time` # 删除 inventory 表中名为 create_time 的列

# 数据操作语句

详见《[数据操作语句](/docs/5.数据存储/2.数据库/关系数据/SQL/数据操作语句/数据操作语句.md)》

# 事务和锁定语句

# 复制语句

# Prepared(预处理) 语句

# Compound(复合) 语句

> 参考：
>
> - [MySQL 官方文档，SQL 语句-复合语句](https://dev.mysql.com/doc/refman/8.0/en/sql-compound-statements.html)

## 控制结构语句

### [case](https://dev.mysql.com/doc/refman/8.0/en/case.html)

### if

# 数据库管理语句

## show # 显示信息

显示有关数据库、表、列的信息，或有关服务器状态的信息。

### EXAMPLE

- 查看创建 test 数据库的 sql 语句
  - show create database test;
- 查看创建 test 表的 sql 语句
  - show create table test;
- 显示数据库的状态信息
  - show status;

# 实用程序语句

## use # 使用指定数据库作为后续 SQL 执行的目标库

# 特殊符号

## @

声明变量、调用变量

# 其他

创建一个 utf8mb4 类型的数据库
create database db2 DEFAULT CHARACTER SET utf8mb4;
创建表
CREATE TABLE students (id int UNSIGNED NOT NULL PRIMARY KEY,name VARCHAR（20）NOT NULL,age tinyint UNSIGNED);

为 emp 表添加记录(有 id，name，sex，age 字段)
insert into emp (id,name,sex,age) values(1,'xiaoming','m',30);

修改 emp 表的内容（第几行第几个字段）
update emp set age=18 where id=4;

批量执行 sql 程序
mysql < hellodb_innodb.sql

备注：也可不进入数据库的情况下查看数据库
mysql -e 'show databases'
