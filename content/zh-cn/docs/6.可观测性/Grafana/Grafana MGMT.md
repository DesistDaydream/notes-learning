---
title: Grafana MGMT
linkTitle: Grafana MGMT
weight: 20
---

# 概述

> 参考：
>
> - [官方文档，管理](https://grafana.com/docs/grafana/latest/administration/)

# Grafana 密码重置

> 参考：
>
> - [官方文档，管理 - CLI - 重置 admin 密码](https://grafana.com/docs/grafana/latest/administration/cli/#reset-admin-password)

## SQLite3 重置

首先需要安装 sqlite3 命令行工具，然后通过 `sqlite3 /PATH/TO/grafana.db` 命令进入 Grafana 数据库

通过 `select login, password, salt from user where login = 'admin';` 语句可以查询到 admin 的当前密码信息

使用下面的 SQL，可以更新 admin 用户的密码为 admin：

```plsql
sqlite> update user set password = '59acf18b94d7eb0694c61e60ce44c110c7a683ac6a8f09580d626f90f4a242000746579358d77dd9e570e83fa24faa88a8a6', salt = 'F3FAxVm33R' where login = 'admin';
```

## API 修改密码

前提是没有忘记密码

```bash
curl -X PUT -H "Content-Type: application/json" -d '{
  "oldPassword": "旧密码",
  "newPassword": "新密码",
  "confirmNew": "新密码"
}' http://账号:旧密码@IP:PORT/api/user/password
```

## grafana-cli 密码重置

```bash
grafana-cli admin reset-admin-password 新密码
```

# 常用 Dashboard 推荐

Kubernetes

- [13105](https://grafana.com/grafana/dashboards/13105-k8s-dashboard-cn-20240513-starsl-cn/)

Node exporter

- [Node Exporter](/docs/6.可观测性/Metrics/Instrumenting/Node%20Exporter.md)

Process exporter

- [Process Exporter](/docs/6.可观测性/Metrics/Instrumenting/Other%20Exporter.md#Process%20Exporter)

Nginx

- [9614](https://grafana.com/grafana/dashboards/9614-nginx-ingress-controller/)
- [12559](https://grafana.com/grafana/dashboards/12559-loki-nginx-service-mesh-json-version/) # Loki NGINX Service Mesh - JSON version

[如何用 Loki 来绘制 Ingress Nginx 监控大屏](https://mp.weixin.qq.com/s/zcY_8c_7eVcRpKh9IHasLg)

# 在代理后面使用 Grafana

https://grafana.com/tutorials/run-grafana-behind-a-proxy/

```nginx
# This is required to proxy Grafana Live WebSocket connections.
map $http_upgrade $connection_upgrade {
  default upgrade;
  '' close;
}

upstream grafana {
  server localhost:3000;
}

server {
  listen 80;
  root /usr/share/nginx/html;
  index index.html index.htm;

  location / {
    proxy_set_header Host $host;
    proxy_pass http://grafana;
  }

  # Proxy Grafana Live WebSocket connections.
  location /api/live/ {
    proxy_http_version 1.1;
    proxy_set_header Upgrade $http_upgrade;
    proxy_set_header Connection $connection_upgrade;
    proxy_set_header Host $host;
    proxy_pass http://grafana;
  }
}
```

> Notes: proxy_pass 字段的值应该替换成真实的 Grafana 地址

# Datasource

## PostgreSQL

当列的 data_type 为 `timestamp without time zone` 时， Grafana 的 PG 数据源会在查询到时间之后，把时间修改为页面设置的时区的值。此时将会产生一个问题，若真实时间是 2004-05-03T08:00:00+08:00，并且也存的是 2004-05-03 08:00:00 这个值，当 Grafana 的 PG 数据源插件在读取到该列的数据类型是 `timestamp without time zone` 时，会把该时间认为时 UTC 时间，若 Grafana 页面设置的是 CST，那么返回结果就会变成 2004-05-03 16:00:00。

> 所以，对于使用 PostgreSQL 存储数据时，如果想要使用 data_type 的信息，那么存 2004-05-03T08:00:00+08:00 这个时间，应该存 2004-05-03 00:00:00 这个字符串

下面有几种方式可以在 Grafana 的时间选择器设置为 “UTC+08:00” 时区时，避免这个情况：

- `::text` 直接忽略 data_type，以字符串形式读取数据
- `AT TIME ZONE` 可以让时间以指定时区的结果返回。相当于 sql 语句 与 数据源插件 一加一减，相互抵消。
- 利用 Grafana 的 Transformations 中的 Format time 功能，设置 data_type 为 `timestamp without time zone` 的列的 timezone 为 UTC

```sql
SELECT
  *,
  create_time::text AS create_time_text,
  create_time AT TIME ZONE 'Asia/Shanghai' AS create_time_tz
FROM my_teble
```

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/grafana/mgmt_postgresql_datasource_data_type_without_time_zone.png)
