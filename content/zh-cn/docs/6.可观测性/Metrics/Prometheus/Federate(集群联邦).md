---
title: Federate(集群联邦)
linkTitle: Federate(集群联邦)
date: 2023-12-20T14:48
weight: 6
---

# 概述

> 参考：
>
> - [官方文档，Prometheus-联邦](https://prometheus.io/docs/prometheus/latest/federation/)

通过 Remote Storage 可以分离监控样本采集和数据存储，解决 Prometheus 的持久化问题。这一部分会重点讨论如何利用联邦集群特性对 Promthues 进行扩展，以适应不同监控规模的变化。

Prometheus Federate 还可以充当代理功能，让 Prometheus Server 获取无法直接访问网段的 Metrics

# 使用联邦集群

对于大部分监控规模而言，我们只需要在每一个数据中心(例如：EC2 可用区，Kubernetes 集群)安装一个 Prometheus Server 实例，就可以在各个数据中心处理上千规模的集群。同时将 Prometheus Server 部署到不同的数据中心可以避免网络配置的复杂性。

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/gx0oz1/1616069518476-78bbd4f5-fc64-4a96-bde6-2309bd716812.jpeg)

如上图所示，在每个数据中心部署单独的 Prometheus Server，用于采集当前数据中心监控数据。并由一个中心的 Prometheus Server 负责聚合多个数据中心的监控数据。这一特性在 Promthues 中称为联邦集群。

联邦集群的核心在于每一个 Prometheus Server 都包含额一个用于获取当前实例中监控样本的接口/federate(用于 web 打开 localhost:9090/federate 即可，初始是空白的，需要详细指明要匹配的内容，才可以获取 metrics)。对于中心 Prometheus Server 而言，无论是从其他的 Prometheus 实例还是 Exporter 实例中获取数据实际上并没有任何差异。其实其他的 promeheus 就相当于中心 prometheus 的一个 exporter

```yaml
scrape_configs:
  - job_name: "federate"
    scrape_interval: 15s
    honor_labels: true
    metrics_path: "/federate"
    params:
      "match[]":
        - 'up{job=~"external.*"}'
    static_configs:
      - targets:
          - "192.168.77.11:9090"
          - "192.168.77.12:9090"
```

配置说明：

1. 通过 URL 中的 match\[]参数指定我们可以指定需要获取的时间序列。match\[]参数必须是一个瞬时向量选择器，例如 up 或者{job="api-server"}。配置多个 match\[]参数，用于获取多组时间序列的监控数据。该例中表示 job 名字开头是 external 的 metric 都抓取
   1. 如果只指定指标名，则获取当前 prom 中该指标名的所有样本；如果只指定标签，则获取所有符合该标签的样本。
   2. 事例中 up{job=~"external._"} 表示获取 up 指标中，标签符合正则 external._ 的所有样本
   3. 注意，可以使用 '{job=~"..\*"}' 来匹配所有 job 的 metric，但是官方不建议这么用，防止意外情况发生
2. horbor_labels 配置 true 可以确保当采集到的监控指标冲突时，能够自动忽略冲突的监控数据。如果为 false 时，prometheus 会自动将冲突的标签替换为”exported\_“的形式。
3. targets 的目标选择要抓取的另一个 prometheus 运行时所监听的 IP:PORT

可以在 web 上进行 match 语句的测试，例如下面，如果获取到的数据和自己预期的一样，那么该配置就没问题

1. [http://172.38.40.214:30001/federate?match\[\]=up{job%3D~"external.\*"}](http://172.38.40.214:30001/federate?match%5B%5D=up%7Bjob%253D~%22external.*%22%7D)
   1. 使用 curl 命令：
      1. curl 'http://172.38.40.214:30001/federate?match\[]=up{job%3D~"external.*"}'
   2. url 中的规则为
      1. {job=~"prometheus"}

功能分区

联邦集群的特性可以帮助用户根据不同的监控规模对 Promthues 部署架构进行调整。例如如下所示，可以在各个数据中心中部署多个 Prometheus Server 实例。每一个 Prometheus Server 实例只负责采集当前数据中心中的一部分任务(Job)，例如可以将不同的监控任务分离到不同的 Prometheus 实例当中，再有中心 Prometheus 实例进行聚合。

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/gx0oz1/1616069518457-5ea7fc2c-1edc-4ce8-acc4-b5bd7534e556.jpeg)

功能分区，即通过联邦集群的特性在任务级别对 Prometheus 采集任务进行划分，以支持规模的扩展。
