---
title: Rules
linkTitle: Rules
weight: 2
---

# 概述

> 参考:
>
> -

Prometheus 规则分为两种：

- **Recording Rule(记录规则)** #
- **Alerting Rule(告警规则)** #
  - ！！！注意编写告警规则的逻辑，由于 Prometheus 会定期评估告警，所以会定期读取数据，尽相避免读取大范围的数据，以免造成性能问题

Prometheus 规则配置文件需要在 [Prometheus Server 配置](/docs/6.可观测性/Metrics/Prometheus/Configuration/Promethesu%20Server.md) 文件中的 rule_files 字段中指定，让 Prometheus 加载指定的文件并读取其配置(这个过程称为 **Evaluation(评估)**)。

一个规则封装了一个向量表达式，该向量表达式在指定的时间间隔内进行评估并采取行动（目前要么记录，要么用于报警）。

可以通过发送 SIGHUP 到 Prometheus 进程在运行时重新加载规则文件。仅当所有规则文件格式正确时，才会应用更改。

> [!Tip] 规则语法检查
> 可以使用 [promtool](/docs/6.可观测性/Metrics/Prometheus/Prometheus%20MGMT/promtool.md) 程序在不启动 Prometheus Server 的情况下检查文件中的语法是否正确。.e.g. `promtool check rules /path/to/example.rules.yml`

# Recording Rule(记录规则)

> 参考：
>
> - [官方文档，配置 - 记录规则](https://prometheus.io/docs/prometheus/latest/configuration/recording_rules/)

在我们使用 Prometheus 的过程中，随着时间的推移，存储在 Prometheus 中的监控指标数据越来越多，查询频率也在不断的增加，当我们用 Grafana 添加更多的 Dashboard 的时候，可能会慢慢的体验到 Grafana 已经无法按时渲染图表，并且偶尔还会出现超时的情况，特别是当我们在长时间汇总大量的指标数据的时候，Prometheus 查询超时的情况可能更多了，这时就需要一种能够类似于后排批处理的机制在后台完成这些复杂运算的计算，对于使用者而言只需要查询这些运算结果即可。

当我们有频繁使用的复杂查询时，如果直接将语句写在 Grafana 的 query 中，Grafana 每次刷新都对 Promethus 提交实时查询，会增加 Prometheus 的性能消耗并且降低了响应速度。 这时候我们就可以用到 Recoding rules 了。

记录规则允许我们预先计算经常使用或计算成本高的表达式，并将其结果保存为一组新的时间序列。因此，查询预先计算的结果通常比每次需要时执行原始表达式快得多。

## 配置示例

这是一个简单的记录规则。使用一个表达式 `sum by (job) (http_inprogress_requests)` 生成了一条新的名为 `job:http_inprogress_requests:sum` 的时间序列

```yaml
groups:
  - name: example
    rules:
      - record: job:http_inprogress_requests:sum
        expr: sum by (job) (http_inprogress_requests)
```

# Alerting Rule(告警规则)

> 参考：
>
> - [官方文档，配置 - 告警规则](https://prometheus.io/docs/prometheus/latest/configuration/alerting_rules/)

**Alerting Rules(告警规则)** 可以让我们基于 PromQL 的表达式，定义告警的触发条件，当满足触发条件时，Prometheus Server 会将触发的告警通知发送到指定的服务。这个服务默认是 Prometheus 官方提供的 [Alertmanager](/docs/6.可观测性/Metrics/Alertmanager/Alertmanager.md)。详见 [Alerting](/docs/6.可观测性/Metrics/Prometheus/Alerting.md)

在 Prometheus 中一条告警规则主要由以下几部分组成：

- 告警名称：用户需要为告警规则命名，当然对于命名而言，需要能够直接表达出该告警的主要内容
- 告警规则：告警规则实际上主要由 PromQL 进行定义，其实际意义是当表达式（PromQL）查询结果持续多长时间（During）后出发告警

在 Prometheus 中，还可以通过 Group(告警组) 对一组相关的告警进行统一定义。当然，这些定义都是通过 YAML 文件来统一管理的。

## 配置示例

```yaml
groups:
  - name: test
    rules:
      - alert: TestAlert
        expr: prometheus_http_requests_total > 20
        for: 1m
        labels:
          alert_test: test
        annotations:
          message: "{{$labels.instance}}: 测试告警(current value is: {{ $value }}"
```

在告警规则文件中，我们可以将一组相关的规则设置定义在一个 group 下。在每一个 group 中我们可以定义多个告警规则。一条告警规则主要由以下几部分组成：

- alert # 告警规则的名称。
- expr # 基于 PromQL 表达式告警触发条件，用于计算是否有时间序列满足该条件。
  - 该样例的意思是某个值减去某几个值的和再除以某个值如果大于 20 就产生告警
- for # evaluation(评估) 等待时间，可选参数。用于表示只有当触发条件持续一段时间后才发送告警。在等待期间新产生告警的状态为 pending。
- labels # 自定义标签，允许用户指定要附加到告警信息上的一组附加标签。
- annotations # 用于指定一组注释，用于描述告警的详细信息，annotations 的内容在告警产生时会一同作为参数发送到 Alertmanager。这里面的 key 与 value 都可以自己定义。这一部分的内容是在讲告警发到接收者的时候，接收者能看到的信息。常用语描述告警信息以便管理员定位问题

## 告警规则配置进阶

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/prometheus/1616069617782-d847748c-9878-4e48-abbb-413799232424.jpeg)

如图所示，在一个告警产生时，会有其自身的 Labels，这些 Labels 信息可以填写进告警规则配置文件中，引用这些 Labels，就相当于把这些 Labels 中的值当做文件内容进行输出。(类似于引用变量)
引用语法：

1. {{ $labels.LabelName }} # 这就是引用该 label 的值，其中 LabelName 就是上图红框中 Labels 的键
2. {{ $value }} # 引用 expr 中 PromQL 表达式中获取到的值

# 配置文件详解

**顶层字段**

**groups**(\[]OBJECT)

- **name**(STRING) # group 的名称，在一个文件中必须是唯一的
- **interval**(DURATION) # 对这个组中的规则进行 evaluated(评估) 的频率。`默认值: PrometheusServer 配置文件中的 global.evaluation_interval 的值`
- **limit**(INT) # 限制警报规则和记录规则可以生成的警报系列的数量。 `默认值: 0`，0 表示没有限制。
- **rules**([]OBJECT) # 定义 Prometheus Rule 详情。rules 字段下使用不同的子字段会对应不同的规则。
  - 包含 record 字段时，则该规则为 [Recording Rule](#Recording%20Rule)
  - 包含 alert 字段时，则该规则为 [Alerting Rule](#Alerting%20Rule)

Recording Rule 与 Alerting Rule 存在于规则组中。组中的规则以规定的时间间隔顺序运行，并具有相同的规则评估时间。Recording Rule 的名称必须是有效的 Metrics 名称。Alerting Rule 名称则比较宽泛，可以随意定义，一般来说，满足有效的标签值即可。

> [!Notes]
> groups 在 recording rule 中并没有实际意义，只是与 alerting rule 同在一个配置文件中，所以两种规则格式要相同罢了，不管在哪个组下定义的记录规则，都可以在任何地方直接使用。
>
> 生成新的时间序列后，Prometheus 会以新的时间序列名称保存数据，该数据与原始 expr 中的表达式所得出的值虽然一样，但是存储的指标是不一样的。

所谓 Evaluated(评估) 规则，就是指 PrometheusServer 会检查规则的状态，如果告警规则的状态是 FIRING，则发送告警。

interval 字段的值 加上 PrometheusServer 的命令行标志 --rules.alert.resend-delay 的值(默认 1m)，才是<font color="#ff0000">**真实的评估周期**</font>。这个说明在官方文档中没有，请参考 [Prometheus 规则处理逻辑中的 - 评估告警规则](/docs/6.可观测性/Metrics/Prometheus/Prometheus%20Development/Rule/Rule.md#评估告警规则)

## Recording Rule

**record**(STRING) # 新的时间序列的名字。必须是有效的 Metrics 名称。

**expr**(STRING) # PromQL 表达式，用于生成新时间序列的。每个评估周期都会在当前时间进行评估，并将结果记录为一组新的时间序列，该时间序列作为 Metrics 的名称由 record 字段定义。

**labels**(map\[STRING]STRING) # 为新的时间序列添加标签集

## Alerting Rule

**alert**(STRING) # 告警名称。

**expr**(STRING) # PromQL 表达式，用于产生告警的。 每个评估周期都会在当前时间进行评估，所有结果时间序列都会变为待处理/触发警报。

**for**(DURATION) # 发送告警的等待时间。`默认值: 0s`，即没有等待期。告警产生后，默认是立刻发送的。配置该字段，可以指定在产生告警后的多长时间再发送告警。

- 在等待期的告警状态为 Pending，超过等待期后，变为 Firing。

**labels**(map\[STRING]STRING) # 为该告警添加或覆盖标签

**annotations**(map\[STRING]STRING) # 为该告警添加注释。

# 最佳实践

> 参考：
>
> - [官方文档，最佳实践 - 记录规则](https://prometheus.io/docs/practices/rules/#recording-rules)

## Recording Rule 命名

Recording 规则名称一般采用 `level:metric:operations` 形式

- **level** # 表示聚合级别，i.e. 根据哪些标签聚合的。e.g. by (instance)，则 level 是 instance
- **metric** # 指标名称，除了在使用 rate() 或 irate() 函数时，需要将指标名称的 `_total` 之外，整体名称应该保持不变
- **operations** # 应用于指标的操作列表

比如：

```yaml
      - record: instance_device:node_disk_io_time_seconds:rate5m
        expr: |-
          rate(node_disk_io_time_seconds_total{job="node-exporter", device=~"(/dev/)?(mmcblk.p.+|nvme.+|rbd.+|sd.+|vd.+|xvd.+|dm-.+|md.+|dasd.+)"}[5m])
```

再比如：

```yaml
groups:
  - name: node-exporter.rules
    rules:
      - record: instance:node_num_cpu:sum
        expr: |
          count without (cpu) (
            count without (mode) (
              node_cpu_seconds_total{job="node-exporter"}
            )
```

如图所示 expr 字段写的表达式与 record 指定的表达式查询结果相同

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/prometheus/1616069617793-9bf3f46e-a10f-4b18-9555-369ba0d8d17f.jpeg)

等同于

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/prometheus/1616069617804-82e7c2f5-8bd7-4932-a5aa-9b702e780e59.jpeg)
