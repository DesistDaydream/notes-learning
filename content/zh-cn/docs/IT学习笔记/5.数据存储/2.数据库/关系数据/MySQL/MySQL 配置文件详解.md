---
title: MySQL 配置文件详解
---

# 概述

> ## 参考：

# my.cnf

\[mysqld]
skip-grant-tables # 登录时跳过权限检查

设置时区
default-time_zone='+8:00'
==========================

\#  开启 binlog
log-bin=mysql-bin
binlog-format=Row
server-id=1
expire_logs_days=7
max_binlog_size=10m

binlog
<https://dev.mysql.com/doc/refman/5.7/en/replication-howto-masterbaseconfig.html>
