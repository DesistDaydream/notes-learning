---
title: Push Gateway
---

# 概述

> 参考：
>
> - [GitHub 项目](https://github.com/prometheus/pushgateway)
> - [官方文档,最佳实践-合时使用 Pushgateway](https://prometheus.io/docs/practices/pushing/#when-to-use-the-pushgateway)

Pushgateway 是 Prometheus 生态中一个重要工具，使用它的原因主要是：

- Prometheus 采用 pull 模式，可能由于不在一个子网或者防火墙原因，导致 Prometheus 无法直接拉取各个 target 数据。
- 在监控业务数据的时候，需要将不同数据汇总, 由 Prometheus 统一收集。

由于以上原因，不得不使用 pushgateway，但在使用之前，有必要了解一下它的一些弊端：

- 将多个节点数据汇总到 pushgateway, 如果 pushgateway 挂了，受影响比多个 target 大。
- Prometheus 拉取状态 up 只针对 pushgateway, 无法做到对每个节点有效。
- Pushgateway 可以持久化推送给它的所有监控数据。

因此，即使你的监控已经下线，prometheus 还会拉取到旧的监控数据，需要手动清理 pushgateway 不要的数据。

Note：
pushgateway 无法主动获取获取目标 metrics。目标需要通过脚本、daemon 程序、手动(e.g.通过 curl 获取 metrics 再发送给 pushgateway)等等方式，主动推送自己的 metrics 到 pushgateway 上。

# PushGateway 部署

```bash
docker run -d -p 9091:9091 prom/pushgateway
```

在 Prometheus Server 的配置文件中加入配置以便让 Prometheus Server 获取 pushgateway 中的 metrics

```yaml
- job_name: 'push_node'
    static_configs:
    - targets: ['10.10.100.110:9091']
      labels:
         env: 'pushgateway'
    honor_labels: true
```

# PushGateway 的使用方式

METRICS | curl --data-binary @- http://IP:PORT/metrics/job/JobName/TableName1/TableValue1/..../TableNameN/TableValueN

通过该命令来将 METRICS 中的内容推送到 PushGateway 中，其中 IP:PORT 就是 PushGateway 所在设备的 IP 及其监听的端口

METRICS # 想要推送给 PushGateway 的 metrics 信息。可以通过 curl 来获取指定对象的 metrics，也可以从文件中读取 metcis 格式的内容，等等。

IP:PORT # PushGateway 程序所在设备的 IP，及 PushGateway 监听的端口

JobName # 指定本次推送 metrics 的 job 名称。

TableXXX # JobName 后面的内容可以作为标签，附加在每个指标上

每当使用该命令给 PushGateway 推送信息后，Prometheus 就会从 PushGateway 中获取该数据并存储在本地。

# 推送数据示例

curl <http://10.10.100.205:9137/metrics> | curl --data-binary @- <http://10.10.100.110:9091/metrics/job/get_vs_state/vs_info/dev_phone_vs>

获取 10.10.100.205 设备上的样本，并推送给 PushGateway 。本次推送的 job 名为 get_vs_state。为每个样本中 metrics 添加标签 vs_info="dev_phone_vs"。下图为本次推送示例图。

示例中，job=get_vs_state 在 Prometheus 中，标签变为 exported_job=get_vs_state。vs_info="dev_phone_vs"则原封不动附加到每个指标的标签后面。

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/wvhciw/1616069386870-f7a4bef3-7a2a-4a3f-9b22-76a0e1010f52.jpeg)

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/wvhciw/1616069386943-043d33a2-dc2a-416b-af2e-75310c7d13d2.jpeg)
