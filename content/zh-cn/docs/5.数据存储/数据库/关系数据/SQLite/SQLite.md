---
title: SQLite
---

# 概述

> 参考：
> 
> - [GitHub 项目，sqlite/sqlite](https://github.com/sqlite/sqlite)
> - [官网](https://www.sqlite.org/index.html)
> - [Wiki，SQLite](https://en.wikipedia.org/wiki/SQLite)

SQLite 是一种 C 语言库，实现小型，快速，自有，高可靠性，全功能，SQL 数据库引擎。 SQLite 是世界上最常用的数据库引擎。 SQLite 内置于所有移动电话和大多数计算机中，并捆绑在人们每天使用的无数其他应用程序内。更多信息...

SQLite 文件格式是稳定的，跨平台和向后兼容的，开发人员承诺通过 2050 年保持这种方式。SQLite 数据库文件通常用作传输系统之间的丰富内容作为数据的长期存档格式。主动使用中有超过 1 万亿（1E12）SQLite 数据库

用白话说，SQLite 通常嵌入在其他程序中，并且 SQLite 存储的数据信息一般只在一个文件中，使用起来非常方便。

通过 sqlite 内置的 sqlite_master 表查看用户创建的表的信息

```sql
# sqlite_master 是一个隐藏的表，维护了 sqlite 数据库自身的数据信息
sqlite> select * from sqlite_master where type="table";
type        name        tbl_name    rootpage    sql
----------  ----------  ----------  ----------  ------------------------------------------
table       memos       memos       2           CREATE TABLE memos(text, priority INTEGER)
```

# sqlite 命令行工具

> 参考：
> 
> - [Ubuntu Manual，sqlite3(1)](https://manpages.ubuntu.com/manpages/jammy/en/man1/sqlite3.1.html)

现阶段通常使用 sqlite3 工具，这个工具可以对接 SQLite 版本 3 的接口。sqlite3 有“交互模式”或“命令行模式”两种模式。

## Syntax(语法)

**sqlite3 \[OPTIONS] \[DatabaseFile] \[SQL]**

SQLite 与 MySQL 不太一样，一个文件是一个数据库。

## 命令行模式

我们可以通过命令行模式查询数据，并且在直接设定输出格式

```sql
$ sqlite3 -header -column mydata.db 'select * from memos;'
text                         priority
---------------------------  ----------
deliver project description  10
lunch with Christine         100
```

**OPTIONS**

sqlite3 的很多选项与交互模式中的点命令具有相同的效果

- **-column** # 以列模式输出查询结果
- **-header** # 输出查询结果时，显示标头

## 交互模式

使用 `sqlite3 DataBase` 命令指定数据库文件即可连接到数据库，并启动交互模式。如果数据库文件不存在，则会自动创建后连接。下面是一个最基本的使用示例，执行了 创建表、插入数据、查询数据 这几个操作：

```sql
$ sqlite3 mydata.db
SQLite version 3.31.1 2020-01-27 19:55:54
Enter ".help" for usage hints.
sqlite> create table memos(text, priority INTEGER);
sqlite> insert into memos values('deliver project description', 10);
sqlite> insert into memos values('lunch with Christine', 100);
sqlite> select * from memos;
deliver project description|10
lunch with Christine|100
```

### 点命令

sqlite3 的交互模式有一组 **Meta-commands(元命令)** 用于 控制输出格式、检查当前连接的数据库文件、对已连接的数据库文件执行管理操作(e.g.重建索引、等等)。所有的元命令始终以 `.` 符号开头，所以以可以称为 **Dot-commands(点命令)**。比如，我可以通过如下方式让 select 查询输出得更加人类可读：

```sql
sqlite> .header on
sqlite> .mode column
sqlite> select * from memos;
text                         priority
---------------------------  ----------
deliver project description  10
lunch with Christine         100
```

所有可用的元命令可以通过 `.help` 这个元命令查看

- **.databases** # 列出已连接的数据库的名称和文件的绝对路径
- **.header \<on|off>** # 输出查询结果时，是否显示标头
- **.separator \<STRING>** # 设置分隔符
- **.show** # 显示所有配置的当前值
- **.tables \[TABLE]** # 列出所有表或指定的表
- **.width \<NUM1 NUM2 ...>** # 为 column 模式设置列宽
