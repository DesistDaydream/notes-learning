---
title: MySQL
---

# 概述

> 参考：
> - [官网](https://www.mysql.com/)

MySQL 的社区版本 MariaDB ，使用安装 MySQL 的 时候，会自动安装 MariaDB 。同时安装 mariadb-server ，即可开始使用了

# MySQL 部署

## docker 启动 MySQL

```bash
mkdir -p /opt/mysql/config
mkdir -p /opt/mysql/data

docker run -d --name mysql --restart always \
--network host \
-v /opt/mysql/data:/var/lib/mysql
-v /opt/mysql/config:/etc/mysql/conf.d
-e MYSQL_ROOT_PASSWORD=mysql \
mysql:8
```

对于 5.7+ 版本，推荐设置 SQL Mode，去掉默认的 ONLY_FULL_GROUP_BY。

```bash
tee /opt/mysql/config/mysql.cnf > /dev/null <<EOF
[mysqld]
sql_mode=STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_AUTO_CREATE_USER,NO_ENGINE_SUBSTITUTION
EOF
```

如果不去掉这个模式，当我们使用 `group by` 时，`select` 中选择的列如果不在 group by 中，将会报错：
`ERROR 1055 (42000): Expression #2 of SELECT list is not in GROUP BY clause  and contains nonaggregated column 'kalacloud.user_id' which is not functionally  dependent on columns in GROUP BY clause; this is incompatible  with sql_mode=only_full_group_by`
如果想解决该错误，除了修改 SQL 模式外，还可以使用 `ANY_VALUE()` 函数处理每一个 select 选中的列但是没有参与 group by 分组的字段。详见：[官方文档，函数和运算符-聚合函数-MySQL 对 GROUP BY 的处理](https://dev.mysql.com/doc/refman/5.7/en/group-by-handling.html)

## Kubernetes 中部署 MySQL

> 参考：
> - [阳明公众号](https://mp.weixin.qq.com/s/C0EYTBJ7sLw823-qE5TjTA)

# MySQL 关联文件与配置

**/etc/my.cnf **# MariaDB 基础配置文件
**/var/lib/myql/\*** # 数据存储路径

# MySQL 数据类型

MySQL 中定义数据字段的类型对你数据库的优化是非常重要的。
MySQL 支持多种类型，大致可以分为三类：数值、日期/时间和字符串(字符)类型。

## 数值类型

MySQL 支持所有标准 SQL 数值数据类型。

这些类型包括严格数值数据类型(INTEGER、SMALLINT、DECIMAL 和 NUMERIC)，以及近似数值数据类型(FLOAT、REAL 和 DOUBLE PRECISION)。

关键字 INT 是 INTEGER 的同义词，关键字 DEC 是 DECIMAL 的同义词。

BIT 数据类型保存位字段值，并且支持 MyISAM、MEMORY、InnoDB 和 BDB 表。

作为 SQL 标准的扩展，MySQL 也支持整数类型 TINYINT、MEDIUMINT 和 BIGINT。下面的表显示了需要的每个整数类型的存储和范围。

| 类型           | 大小                                          | 范围（有符号）                                                                                                                      | 范围（无符号）                                                    | 用途       |
| -------------- | --------------------------------------------- | ----------------------------------------------------------------------------------------------------------------------------------- | ----------------------------------------------------------------- | ---------- |
| TINYINT        | 1 byte                                        | (-128，127)                                                                                                                         | (0，255)                                                          | 小整数值   |
| SMALLINT       | 2 bytes                                       | (-32 768，32 767)                                                                                                                   | (0，65 535)                                                       | 大整数值   |
| MEDIUMINT      | 3 bytes                                       | (-8 388 608，8 388 607)                                                                                                             | (0，16 777 215)                                                   | 大整数值   |
| INT 或 INTEGER | 4 bytes                                       | (-2 147 483 648，2 147 483 647)                                                                                                     | (0，4 294 967 295)                                                | 大整数值   |
| BIGINT         | 8 bytes                                       | (-9,223,372,036,854,775,808，9 223 372 036 854 775 807)                                                                             | (0，18 446 744 073 709 551 615)                                   | 极大整数值 |
| FLOAT          | 4 bytes                                       | (-3.402 823 466 E+38，-1.175 494 351 E-38)，0，(1.175 494 351 E-38，3.402 823 466 351 E+38)                                         | 0，(1.175 494 351 E-38，3.402 823 466 E+38)                       | 单精度     |
| 浮点数值       |
| DOUBLE         | 8 bytes                                       | (-1.797 693 134 862 315 7 E+308，-2.225 073 858 507 201 4 E-308)，0，(2.225 073 858 507 201 4 E-308，1.797 693 134 862 315 7 E+308) | 0，(2.225 073 858 507 201 4 E-308，1.797 693 134 862 315 7 E+308) | 双精度     |
| 浮点数值       |
| DECIMAL        | 对 DECIMAL(M,D) ，如果 M>D，为 M+2 否则为 D+2 | 依赖于 M 和 D 的值                                                                                                                  | 依赖于 M 和 D 的值                                                | 小数值     |

## 日期和时间类型

表示时间值的日期和时间类型为 DATETIME、DATE、TIMESTAMP、TIME 和 YEAR。

每个时间类型有一个有效值范围和一个"零"值，当指定不合法的 MySQL 不能表示的值时使用"零"值。

TIMESTAMP 类型有专有的自动更新特性，将在后面描述。

| 类型                                                                                                   | 大小( bytes)    | 范围                                    | 格式                | 用途             |
| ------------------------------------------------------------------------------------------------------ | --------------- | --------------------------------------- | ------------------- | ---------------- |
| DATE                                                                                                   | 3               | 1000-01-01/9999-12-31                   | YYYY-MM-DD          | 日期值           |
| TIME                                                                                                   | 3               | '-838:59:59'/'838:59:59'                | HH:MM:SS            | 时间值或持续时间 |
| YEAR                                                                                                   | 1               | 1901/2155                               | YYYY                | 年份值           |
| DATETIME                                                                                               | 8               | 1000-01-01 00:00:00/9999-12-31 23:59:59 | YYYY-MM-DD HH:MM:SS | 混合日期和时间值 |
| TIMESTAMP                                                                                              | 4               | 1970-01-01 00:00:00/2038                |
| 结束时间是第 2147483647 秒，北京时间 2038-1-19 11:14:07，格林尼治时间 2038 年 1 月 19 日 凌晨 03:14:07 | YYYYMMDD HHMMSS | 混合日期和时间值，时间戳                |

## 字符串类型

字符串类型指 CHAR、VARCHAR、BINARY、VARBINARY、BLOB、TEXT、ENUM 和 SET。该节描述了这些类型如何工作以及如何在查询中使用这些类型。

| 类型       | 大小                  | 用途                            |
| ---------- | --------------------- | ------------------------------- |
| CHAR       | 0-255 bytes           | 定长字符串                      |
| VARCHAR    | 0-65535 bytes         | 变长字符串                      |
| TINYBLOB   | 0-255 bytes           | 不超过 255 个字符的二进制字符串 |
| TINYTEXT   | 0-255 bytes           | 短文本字符串                    |
| BLOB       | 0-65 535 bytes        | 二进制形式的长文本数据          |
| TEXT       | 0-65 535 bytes        | 长文本数据                      |
| MEDIUMBLOB | 0-16 777 215 bytes    | 二进制形式的中等长度文本数据    |
| MEDIUMTEXT | 0-16 777 215 bytes    | 中等长度文本数据                |
| LONGBLOB   | 0-4 294 967 295 bytes | 二进制形式的极大文本数据        |
| LONGTEXT   | 0-4 294 967 295 bytes | 极大文本数据                    |

注意：char(n) 和 varchar(n) 中括号中 n 代表字符的个数，并不代表字节个数，比如 CHAR(30) 就可以存储 30 个字符。

CHAR 和 VARCHAR 类型类似，但它们保存和检索的方式不同。它们的最大长度和是否尾部空格被保留等方面也不同。在存储或检索过程中不进行大小写转换。

BINARY 和 VARBINARY 类似于 CHAR 和 VARCHAR，不同的是它们包含二进制字符串而不要非二进制字符串。也就是说，它们包含字节字符串而不是字符字符串。这说明它们没有字符集，并且排序和比较基于列值字节的数值值。

BLOB 是一个二进制大对象，可以容纳可变数量的数据。有 4 种 BLOB 类型：TINYBLOB、BLOB、MEDIUMBLOB 和 LONGBLOB。它们区别在于可容纳存储范围不同。

有 4 种 TEXT 类型：TINYTEXT、TEXT、MEDIUMTEXT 和 LONGTEXT。对应的这 4 种 BLOB 类型，可存储的最大长度不同，可根据实际情况选择。

# mysql 命令行工具

> 参考：
> - [官方文档，MySQL 程序-客户端程序-mysql](https://dev.mysql.com/doc/refman/8.0/en/mysql.html)

mysql 是一个简单的 SQL Shell。 它支持交互和非交互使用。 交互使用时，查询结果以 ASCII 表格式显示。 非交互使用（例如，用作过滤器）时，结果以制表符分隔的格式显示。 可以使用命令选项更改输出格式。

## Syntax(语法)

**mysql \[OPTIONS] \[DATABASE]**
**DATABASE **# 指定连接 mysql 后要操作的数据库。若不指定，则需要在交互模式下使用 `use` 指令选择数据库，否则对数据库的操作将会报 `No database selected` 错误：

```bash
mysql> show tables;
ERROR 1046 (3D000): No database selected
```

**OPTIONS：**

- **-h, --host <HostName> **# 指定要连接的 mysql 主机。如果链接本机 mysql，可以省略。
- **-P, --port <PORT> **# 指定要连接的 mysql 的端口。默认值：`3306`
- **-u, --user <UserName> **# 指定要登录 mysql 的用户名
- **-p, --password <PASSWORD>** # 使用密码来登录。如果指定要登录 mysql 的用户密码为空，则该选项可省

## 命令行模式

略

## 交互模式

### 斜线命令

在 mysql 的交互模式中有一组 mysql 程序自带的命令，用以 控制输出格式、检查、获取数据信息 等等，这些命令以 `\` 开头，不过也有与之相对应的字符串命令

- **\u, use <DBName> **# 选择想要操作的数据库

### 基础示例

- grant select,insert,update,delete,create,drop ON mysql.\* TO 'lichenhao'@'localhost' identified by 'lichenhao'; # 为名为 mysql 的数据库创建名为 lichenhao 的用户，密码为 lichenhao，具有 select、insert、update、delete、create、drop 这些命令的执行权限。
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


    MariaDB [mysql]> SHOW COLUMNS FROM db;
    +-----------------------+---------------+------+-----+---------+-------+
    | Field                 | Type          | Null | Key | Default | Extra |
    +-----------------------+---------------+------+-----+---------+-------+
    | Host                  | char(60)      | NO   | PRI |         |       |
    | Db                    | char(64)      | NO   | PRI |         |       |
    | User                  | char(16)      | NO   | PRI |         |       |
    | Select_priv           | enum('N','Y') | NO   |     | N       |       |
    .......

- select Host,db from db; #显示 db 表中，Host 和 Db 列及其内容，效果如下


    MariaDB [mysql]> SELECT Host,db from db;
    +-----------+---------+
    | Host      | db      |
    +-----------+---------+
    | %         | test    |
    | %         | test\_% |
    | localhost | mysql   |
    +-----------+---------+

# mysqladmin

EXAMPLE

- mysqladmin -u root -p password "my_password" #修改 root 密码，密码为：my_password。如果默认密码为空，则可以不加-p。
