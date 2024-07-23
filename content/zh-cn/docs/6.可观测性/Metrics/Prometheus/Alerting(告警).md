---
title: Alerting(告警)
linkTitle: Alerting(告警)
date: 2023-12-20T14:48
weight: 5
---

# 概述

> 参考：
>
> - [官方文档，告警 - 告警概述](https://prometheus.io/docs/alerting/latest/overview/)
> - [官方文档，告警 - 客户端](https://prometheus.io/docs/alerting/latest/clients/)
> - [OpenAPI](https://github.com/prometheus/alertmanager/blob/main/api/v2/openapi.yaml)

Prometheus 本身不提告警的通知的功能！告警能力在 Prometheus 的架构中被划分成两个独立的部分。如下所示，通过在 Prometheus 中定义 AlertRule（告警规则），Prometheus 会周期性的对告警规则进行 **Evaluate(评估)**，如果满足告警触发条件就会向 Alertmanager 发送告警信息。

**Evaluate(评估)** 就是指，Prometheus Server 会定期执行规则配置文件中的 PromQL，获得结果并与阈值进行匹配，当超过设置的阈值时，会产生告警。这个过程，就称为 **Evaluate(评估)**。在代码中，通过 Eval() 方法来评估规则。
![](https://notes-learning.oss-cn-beijing.aliyuncs.com/sw6o6t/1616069590594-41190e69-d023-4ef4-87ad-fdc1a7cf8b6f.png)

Alertmanager 处理客户端应用程序(如 Prometheus Server)发送的警报。它负责对它们进行重复数据删除，分组和路由，以及正确的接收器集成，例如 email，PagerDuty 或 OpsGenie。它还负责警报的静音和抑制。

即：Prometheus Server 只负责根据 PromQL 语句定义的规则产生告警并发送给 Alertmanager(告警管理器)。

注意：

- Alertmanager 是一个单独的程序，需要独立安装使用
- Alertmanager 既可以描述为一类具有处理告警功能的应用程序。也可以描述为一个 Prometheus 官方推出的名为 Alertmanager 的程序。以后的描述一般都不加区分

## 关联 Alertmanager 与 Prometheus

由于 Alertmanager 与 Prometheus 是两个程序。所以需要修改 Prometheus Server 的配置文件，以便让自己产生的告警可以发送到正确地方，配置效果如下（Prometheus 推出的 Alertmanager 默认监听在 9093 端口上）

```yaml
alerting:
  alertmanagers:
    - static_configs:
        - targets: ["localhost:9093"]
```

告警规则的配置，由于 Prometheus Server 自己产生告警，所以还需要在配置文件中指定具体根据哪个《告警规则的配置文件》来生成告警

```yaml
rule_files:
  - /etc/prometheus/rules.yml
```

下文会介绍配置文件的详细用法

### Prometheus 推出的 Alertmanager 程序简介

详见：[Alertmanager](/docs/6.可观测性/Metrics/Alertmanager/Alertmanager.md)

Prometheus 推出的 Alertmanager 作为一个独立的组件，可以实现告警管理功能，负责接收并处理来自 Prometheus Server(也可以是其它的客户端程序)的告警信息。Alertmanager 可以对这些告警信息进行进一步的处理，比如当接收到大量重复告警时能够消除重复的告警信息，同时对告警信息进行分组并且路由到正确的通知方，Alertmanager 内置了对邮件，Slack 等多种通知方式的支持。同时 AlertManager 还提供了静默和告警抑制机制来对告警通知行为进行优化。

## 查看告警的状态

在 prometheus server 的 web 页面中的 `Alerts` 标签查看到所有其所配置和产生的告警信息，效果如图：

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/sw6o6t/1616069590604-e9eaacdf-e275-4662-b8f1-8d1739a63fc4.jpeg)

其中每行都是一条告警规则，绿色表示没有达到设定的阈值，不会产生告警；红色的表示达到设定的阈值并已经持续了一段时间，所以产生了告警，并推送给 alermanager。在绿条或者红条中间的位置是“路径>组名”这里表示其下的所有告警都是这个组里的。点开一个告警，就能看到其中的配置，包括告警规则的名称、告警触发条件、等待时长等等信息

## 告警发送过程

在 State 中，有两种状态，一个是 PENDING，一个是 FIRING。当告警刚刚出发时，处于 PENDING 状态，此时告警并不会推送到 Alertmanager，当该状态无法恢复且持续配置中定义一段时间后，则变为 FIRING，并向 Alertmanager 推送告警。

当一个告警解决后，会具有一个隐藏的 Pending 状态，持续 15 分钟，在这 15 分钟之内，依然会重复发送告警，只不过发送的每个告警的结束时间，都是同一个，就是解决告警的时间。这个 15 分钟是不可变的，在代码 github.com/prometheus/prometheus/rules/alerting.go 这个里，有一个常量 `resolvedRetention` 就是用来判断何时删除一个未激活告警的条件之一。

# Prometheus 告警规则配置

详见：[Rules 配置](/docs/6.可观测性/Metrics/Prometheus/Rules%20配置.md)

# 告警数据结构

免责声明：Prometheus 会自动负责发送由其配置的 **[警报规则](https://prometheus.io/docs/prometheus/latest/configuration/alerting_rules/)** 生成的 **[警报](https://prometheus.io/docs/prometheus/latest/configuration/alerting_rules/)**。强烈建议在 Prometheus 中根据时间序列数据配置警报规则，而不是实现直接客户端。

**<font color="#ff0000">也就是说，不要自己写一个程序，频繁对 Prometheus 发起 PromQL 查询请求，来生成告警。</font>**

Alertmanager 现阶段有两个 API，v1 和 v2，这两个 API 都是用来监听发送到自身的告警。

Prometheus 产生告警后，会通过 POST 请求将下列 JSON 格式内容向 Alertamanger 推送告警：

```json
[
  {
    "labels": {
      "alertname": "<Prometheus Server 的规则配置文件中的 .groups.rules.alert 字段的值>",
      "<LabelName>": "<LabelValue>",
      ...
    },
    "annotations": {
      "<LabelName>": "<LabelValue>",
   ...
    },
    "startsAt": "<RFC3339>",
    "endsAt": "<RFC3339>",
    "generatorURL": "<GeneratorURL>" // 生成告警的 URL。就是可以向 Prometheus API 发送的包含 PromQL 的完整 URL
  },
  ...
]
```

推送路径根据 Prometheus Server 配置文件中 `alerting.alertmanagers.api_version` 和 `alerting.alertmanagers.path_prefix` 这两个字段决定。

默认推送路径为 /api/v2/alerts。如果 api_version 为 v2，path_prefix 值为 /test，最终的路径就是 /test/api/v2/alerts

## labels 与 annotations 字段

**labels(标签)** 是告警的唯一标识符。具有相同标签的告警，则称为重复数据，重复数据只会保留最新的一个。

**annotations(注释)** 顾名思义，就是用来注释一个告警

labels 包含如下内容：

- alertname 字段
  - 该字段的的值就是 Prometheus Server 的 Rules 配置文件中的 .groups.rules.alert 字段的值
- 告警规则配置文件中定义的标签
- 产生告警的时间序列所具有的标签

## startsAt 与 endsAt 字段

**startsAt** # 告警的开始时间

**endsAt** # 告警的结束时间

- 结束时间可以这么理解：从开始时间到结束时间，如果 Alertmanager 没有再收到相同的告警，则认为告警已经处理

> 注意：对于 Prometheus 官方的 Alertmanager 来说，startsAt 和 endsAt 时间戳都是可选的。如果省略了 startAt，则由 Alertmanager 分配当前时间。 endsAt 只有在已知警报的结束时间时才会设置。否则，它将被设置为从最后一次收到警报的时间开始的一个可配置的超时时间段。

## 示例

### 配置文件

```yaml
groups:
  - name: test
    rules:
      - alert: test
        expr: vector(1)
        labels:
          level: warning
        annotations:
          additionalExample: test
```

### 生成的告警数据

```json
[
  {
    "labels": {
      "alertname": "test",
      "level": "warning"
    },
    "annotations": {
      "additionalExample": "test"
    },
    "startsAt": "2021-02-23T03:56:42.944457098Z",
    "endsAt": "2021-02-23T04:04:27.944457098Z",
    "generatorURL": "http://cs-cs-prometheus.desistdaydream.ltd/graph?g0.expr=vector%281%29\u0026g0.tab=1"
  }
]
```

抓包内容如下：

![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/sw6o6t/1621754366379-909c188e-f854-4c8e-8a9d-e75b6e671d2c.png)
