---
title: "MySQL 部署"
linkTitle: "MySQL 部署"
weight: 20
---

# 概述

# docker 启动 MySQL

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

# Kubernetes 中部署 MySQL

> 参考：
>
> - [阳明公众号](https://mp.weixin.qq.com/s/C0EYTBJ7sLw823-qE5TjTA)

## Helm 部署 MySQL

> 参考：
> - [bitnami/mysql](https://github.com/bitnami/charts/tree/main/bitnami/mysql)