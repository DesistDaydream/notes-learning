---
title: Other Exporter
---

# 概述

官方推荐的一些第三方 exporter：<https://prometheus.io/docs/instrumenting/exporters/>

# Kubernetes Event Exporter

> 参考：
> - [GitHub 项目，opsgenie/kubernetes-event-exporter](https://github.com/opsgenie/kubernetes-event-exporter)

将 Kubernetes 中 Event 资源的内容导出到多个目的地

# Nginx Exporter

参考：[GitHub 项目](https://github.com/nginxinc/nginx-prometheus-exporter)

端口默认 9113，Grafana 面板 ID：[GitHub](https://raw.githubusercontent.com/nginxinc/nginx-prometheus-exporter/master/grafana/dashboard.json)

    docker run -d --name nginx-exporter\
      --network host \
      nginx/nginx-prometheus-exporter:0.8.0 \
      -nginx.scrape-uri http://localhost:9123/metrics

# MySQL Exporter

参考：[GitHub 项目](https://github.com/prometheus/mysqld_exporter)
端口默认 9104，Grafana 面板 ID：7362

    docker run -d --name mysql-exporter \
      --network host  \
      -e DATA_SOURCE_NAME="root:oc@2020@(localhost:3306)/" \
      prom/mysqld-exporter

# Redis Exporter

参考：[GitHub 项目](https://github.com/oliver006/redis_exporter)
端口默认 9121，Grafana 面板 ID：11835

    docker run -d --name redis-exporter \
      --network host \
      oliver006/redis_exporter \
      --redis.password='DesistDaydream' \
      --redis.addr='redis://127.0.0.1:6379'

# DELL OMSA 的 Exporter

项目地址：<https://github.com/galexrt/dellhw_exporter>

运行方式：docker run -d --name dellhw_exporter --privileged -p 9137:9137 galexrt/dellhw_exporter

如果 container 无法启动，一般都是因为 OMSA 的资源指标采集服务未安装导致。

该 Exporter 会读取 OMSA 的 SNMP 信息转换成 Metrics 并在 9137 端口输出

# IPMI Exporter

项目地址：<https://github.com/prometheus-community/ipmi_exporter>
依赖于 FreeIPMI 套件中的工具来实现

# HAProxy 的 Exporter

项目地址：<https://github.com/prometheus/haproxy_exporter>
