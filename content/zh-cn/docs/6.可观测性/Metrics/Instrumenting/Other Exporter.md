---
title: Other Exporter
---

# 概述

官方推荐的一些第三方 exporter: https://prometheus.io/docs/instrumenting/exporters/

process-exporter: 采集进程指标

- https://github.com/ncabatoff/process-exporter

# K8S 集群中常用的 exporter

> - [Vermouth 博客，高可用 Prometheus 问题集锦](http://www.xuyasong.com/?p=1921)


可以在[这里](https://prometheus.io/docs/instrumenting/exporters/)看到官方、非官方的 exporter。如果还是没满足你的需求，你还可以自己编写 exporter，简单方便、自由开放，这是优点。

但是过于开放就会带来选型、试错成本。之前只需要在 zabbix agent 里面几行配置就能完成的事，现在你会需要很多 exporter 搭配才能完成。还要对所有 exporter 维护、监控。尤其是升级 exporter 版本时，很痛苦。非官方 exporter 还会有不少 bug。这是使用上的不足，当然也是 Prometheus 的设计原则。

K8S 生态的组件都会提供 / metric 接口以提供自监控，这里列下我们正在使用的：

- [cadvisor](http://www.xuyasong.com/?p=1483): 集成在 Kubelet 中。
- kubelet: 10255 为非认证端口，10250 为认证端口。
- apiserver: 6443 端口，关心请求数、延迟等。
- scheduler: 10251 端口。
- controller-manager: 10252 端口。
- etcd: 如 etcd 写入读取延迟、存储容量等。
- docker: 需要开启 experimental 实验特性，配置 metrics-addr，如容器创建耗时等指标。
- kube-proxy: 默认 127 暴露，10249 端口。外部采集时可以修改为 0.0.0.0 监听，会暴露：写入 iptables 规则的耗时等指标。
- [kube-state-metrics](http://www.xuyasong.com/?p=1525): K8S 官方项目，采集 pod、deployment 等资源的元信息。
- [node-exporter](http://www.xuyasong.com/?p=1539): Prometheus 官方项目，采集机器指标如 CPU、内存、磁盘。
- blackbox_exporter: Prometheus 官方项目，网络探测，dns、ping、http 监控
- process-exporter: 采集进程指标
- nvidia exporter: 我们有 gpu 任务，需要 gpu 数据监控
- node-problem-detector: 即 npd，准确的说不是 exporter，但也会监测机器状态，上报节点异常打 taint
- 应用层 exporter: mysql、nginx、mq 等，看业务需求。

还有各种场景下的[自定义 exporter](http://www.xuyasong.com/?p=1942)，如日志提取后面会再做介绍。

## Kubernetes Event Exporter

> 参考：
> 
> - [GitHub 项目，opsgenie/kubernetes-event-exporter](https://github.com/opsgenie/kubernetes-event-exporter)

将 Kubernetes 中 Event 资源的内容导出到多个目的地

# Process Exporter

> 参考：
>
> - [GitHub 项目，ncabatoff/process-exporter](https://github.com/ncabatoff/process-exporter)

process_names 下的数组定义进程组名称及该进程组的匹配条件，一共 3 个匹配方式

- **comm** # 与 /proc/${pid}/stat 中第二个字段进行匹配
- **exe** # 
- **cmdline** # 与进程的所有参数进行匹配

```yaml
process_names:
  # comm is the second field of /proc/<pid>/stat minus parens.
  # It is the base executable name, truncated at 15 chars.
  # It cannot be modified by the program, unlike exe.
  - comm:
    - bash

  # exe is argv[0]. If no slashes, only basename of argv[0] need match.
  # If exe contains slashes, argv[0] must match exactly.
  - exe:
    - postgres
    - /usr/local/bin/prometheus

  # cmdline is a list of regexps applied to argv.
  # Each must match, and any captures are added to the .Matches map.
  - name: "{{.ExeFull}}:{{.Matches.Cfgfile}}"
    exe:
    - /usr/local/bin/process-exporter
    cmdline:
    - -config.path\s+(?P<Cfgfile>\S+)
```

若进程组不指定 name 字段，默认为 `{{.ExeBase}}`(可执行文件的 basename)

# Nginx Exporter

参考：[GitHub 项目](https://github.com/nginxinc/nginx-prometheus-exporter)

端口默认 9113，Grafana 面板 ID：[GitHub](https://raw.githubusercontent.com/nginxinc/nginx-prometheus-exporter/master/grafana/dashboard.json)

```bash
docker run -d --name nginx-exporter\
  --network host \
  nginx/nginx-prometheus-exporter:0.8.0 \
  -nginx.scrape-uri http://localhost:9123/metrics
```

# MySQL Exporter

> 参考：
> - [GitHub 项目](https://github.com/prometheus/mysqld_exporter)

端口默认 9104

Grafana 面板：

- 7362
- 官方提供了一个面板在[这里](https://github.com/prometheus/mysqld_exporter/blob/main/mysqld-mixin/dashboards/mysql-overview.json)

```bash
docker run -d --name mysql-exporter \
  --network host  \
  -e DATA_SOURCE_NAME="root:oc@2020@(localhost:3306)/" \
  prom/mysqld-exporter
```

# Redis Exporter

参考：[GitHub 项目](https://github.com/oliver006/redis_exporter)

端口默认 9121，Grafana 面板 ID：11835

```bash
docker run -d --name redis-exporter \
  --network host \
  oliver006/redis_exporter \
  --redis.password='DesistDaydream' \
  --redis.addr='redis://127.0.0.1:6379'
```

# DELL OMSA 的 Exporter

https://github.com/galexrt/dellhw_exporter

运行方式: docker run -d --name dellhw_exporter --privileged -p 9137:9137 galexrt/dellhw_exporter

如果无法启动，一般都是因为 OMSA 的资源指标采集服务未安装导致。

该 Exporter 会读取 OMSA 的 SNMP 信息转换成 Metrics 并在 9137 端口输出

# IPMI Exporter

https://github.com/prometheus-community/ipmi_exporter

依赖于 FreeIPMI 套件中的工具来实现

# HAProxy 的 Exporter

项目地址：<https://github.com/prometheus/haproxy_exporter>
