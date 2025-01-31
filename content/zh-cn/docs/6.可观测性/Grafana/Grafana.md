---
title: Grafana
weight: 1
---

# 概述

> 参考：
>
> - [官网](https://grafana.com/)
> - [GitHub 项目，grafana/grafana](https://github.com/grafana/grafana)

Grafana 是开源的可视化和分析软件。它使我们可以查询，可视化，警报和浏览指标，无论它们存储在哪里。它为您提供了将 [时间序列数据](/docs/5.数据存储/数据库/时间序列数据/时间序列数据.md) 转换为精美的图形和可视化效果的工具。

# Grafana 部署

> 参考：
>
> - [官方文档，安装 - 安装 Grafana](https://grafana.com/docs/grafana/latest/setup-grafana/installation/)

## docker 方式运行 grafana

获取配置文件

```shell
mkdir -p /opt/monitoring/server/config/grafana
mkdir -p /opt/monitoring/server/data/grafana
chown -R 472 /opt/monitoring/server/data/grafana
docker run -d --name grafana --rm grafana/grafana
docker cp grafana:/etc/grafana /opt/monitoring/server/config
docker stop grafana
```

运行 Grafana

```shell
docker run -d --name grafana \
  --network host \
  -v /opt/monitoring/server/config/grafana:/etc/grafana \
  -v /opt/monitoring/server/data/grafana:/var/lib/grafana \
  -v /etc/localtime:/etc/localtime \
  grafana/grafana
```

# Grafana 关联文件与配置

**/etc/grafana/** # Grafana 配置文件保存路径

- **./grafana.ini** # Grafana 运行所需配置文件
- **./provisioning/** # Grafana 的 Provisioning 功能默认要读取的路径。该功能详见 [Provisioning 配置](docs/6.可观测性/Grafana/Grafana%20Configuration/Provisioning%20配置.md)。可以通过 grafana.ini 的 .paths.provisioning 字段修改读取路径

**/var/lib/grafana/** # Grafana 数据保存路径

- **./grafana.db** # Grafana 数据文件，包括 用户信息、dashboard、datasource 等等。这是一个 SQLite3 数据库文件。
- **./plugins/** # Grafana 安装的插件保存在该目录下
