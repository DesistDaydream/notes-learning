---
title: ClickHouse MGMT
linkTitle: ClickHouse MGMT
date: 2025-01-05T21:03:00
weight: 20
---

# 概述

> 参考：
>
> - [官方文档，管理和部署](https://clickhouse.com/docs/guides/manage-and-deploy-index)
>   - 链接有点乱，官网页面关于这部分大改过。截至本文更新可以从页面的 Server Admin 点进来。

# ClikcHouse 部署后常见操作
## 添加只读用户

```sql
-- TODO
```

### 通过修改文件添加只读用户

添加 developer 用户，密码为 developer（通过 `echo -n "readONLY@123" | sha256sum | tr -d '-'` 命令获取加密后的字符串填到 password_sha256_hex 元素中）

默认数据库为 network_security，允许访问的数据库为 network_security

```xml
    <profiles>
        <developer>
            <readonly>2</readonly>
            <allow_ddl>0</allow_ddl>
        </developer>
    </profiles>
    <users>
        <developer>
                <password_sha256_hex>88fa0d759f845b47c044c2cd44e29082cf6fea665c30c146374ec7c8f3d699e3</password_sha256_hex>
                <default_database>network_security</default_database>
                <allow_databases>
                    <database>network_security</database>
                </allow_databases>
        </developer>
    </users>
```

# 备份与恢复

> 参考：
>
> - [官方文档，管理与部署 - 备份与恢复](https://clickhouse.com/docs/en/operations/backup)
> - [GitHub 项目，AlexAkulov/clickhouse-backup](https://github.com/AlexAkulov/clickhouse-backup)

