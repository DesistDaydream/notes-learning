---
title: ServiceDiscovery
linkTitle: ServiceDiscovery
weight: 20
---

# 概述

> 参考：
>
> -

**Service Discovery(服务发现)**

# Nacos

> 参考：
> 
> - [官网](https://nacos.io/zh-cn/index.html)
> - [官方文档](https://nacos.io/zh-cn/docs/what-is-nacos.html)
> - [部署文档](https://github.com/nacos-group/nacos-k8s)

发现、配置和管理微服务

# Consul

> 参考：
>
> - [GitHub 项目，hashicorp/consul](https://github.com/hashicorp/consul)


注册

```bash
curl --location --request PUT 'http://10.10.4.90:8500/v1/agent/service/register' \
--header 'Content-Type: application/json' \
--data-raw '{
    "id": "10.10.11.16",
    "name": "node-exporter",
    "address": "10.10.11.16",
    "port": 9100,
    "tags": [
        "linux"
    ],
    "Meta": {
        "custom_house_name": "天津机房",
        "business_type": "迎检"
    }
}'
```

利用服务的 Meta 在 Prometheus 进行 consul_sd 时添加 Lable