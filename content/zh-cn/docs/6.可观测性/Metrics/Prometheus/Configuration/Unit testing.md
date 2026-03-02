---
title: "Unit testing"
linkTitle: "Unit testing"
weight: 100
---

# 概述

> 参考：
>
> - [官方文档，配置 - 规则单元测试](https://prometheus.io/docs/prometheus/latest/configuration/unit_testing_rules/)

# 配置文件详解

**顶级字段**

- **rule_files**(\[]STRING) # 规则文件列表
- **evaluation_interval**
- **fuzzy_compare**
- **group_eval_order**
- **tests**(\[][tests](#tests)) # 单元测试规则列表

## tests

**input_series** # 时间序列数据。人为定义

**name** # 单元测试的名称

**start_timestamp**

**alert_rule_test** # 告警规则的单元测试逻辑

**promql_expr_test** # PromQL 的单元测试逻辑

**external_labels**

# 最佳实践

[官方最佳实践案例](https://prometheus.io/docs/prometheus/latest/configuration/unit_testing_rules/#example)

## 测试 Promehteus 模板渲染结果

test_rules # 告警规则文件

```yaml
groups:
  - name: "模板测试"
    rules:
      - alert: "测试时间相关模板函数"
        expr: |-
          my_process_boot_time_seconds
          and
          changes(my_process_pid[15m]) > 0
        labels:
          alert_target: "节点"
        annotations:
          summary: "进程已重启"
          description: |-
            当前系统启动时间: {{ ($value | toTime).Format "2006-01-02 15:04:05" }}
      - alert: "测试数字相关模板函数"
        expr: |-
          irate(node_network_receive_bytes_total[15m]) * 8 >= 0
        labels:
          alert_event: "网络"
        annotations:
          summary: "{{ $labels.instance }} 设备的 {{ $labels.device_name }} 网卡接受流量异常"
          description: |-
            带宽1024: {{ humanize1024 $value }}B/s
            带宽1000: {{ humanize $value }}B/s
```

unit_test_alert.yaml # 单元测试文件

```yaml
rule_files:
  - ./test_rules.yaml

tests:
  - interval: 3m # input_series 给定的时间序列中，每个样本之间的时间间隔
    # 时间序列数据。手动创建的时间序列及其样本值。
    input_series:
      - series: 'my_process_boot_time_seconds{instance="192.168.1.100:9100"}'
        # 1767196800 是 2026-01-01 00:00:00。但是 Prometheus 时区默认是 UTC，最后结果是: 2025-12-31 16:00:00
        # +0x5 的意思是一共 5 个样本值，每个样本值之间的差值是 0。相当于 1767196800, 1767196800, 1767196800, 1767196800, 1767196800
        values: "1767196800+0x5"
      - series: 'my_process_pid{instance="192.168.1.100:9100"}'
        values: "1000+100x5" # 1000, 1100, 1200, 1300, 1400
    # 测试规则。使用 input_series 中给定的时间序列数据，评估 rule_files 文件中与 alertname 指定的名称相同的告警规则
    alert_rule_test:
      - eval_time: 15m # 评估时间，如果想要看到多个样本值的变化，至少要是上面 interval 设定值的 2 倍。
        alertname: "测试时间相关模板函数" # 需要进行测试的告警规则，该名称必须要在 rule_files 定义的文件中存在
        # 期望的告警结果。
        # 注意：如果想要检查告警中的模板渲染结果，不要设置 exp_alerts 字段。
        # 如果设定了期望的结果，并且测试成功通过，那么执行结果不会显示渲染结果，只会显示 SUCCESS 这个字符串
        # exp_alerts:
        #   - exp_labels:
        #       alert_target: "node"
        #       alertname: "测试时间相关模板函数"
        #       instance: "192.168.1.100:9100"
        #     exp_annotations:
        #       summary: "进程已重启"
        #       description: "当前系统启动时间: 2025-12-31 16:00:00"
  - interval: 3m
    input_series:
      - series: 'node_network_receive_bytes_total{instance="192.168.1.100:9100", device_name="eth0"}'
        values: "1000000000000+100000000000x10"
    alert_rule_test:
      - eval_time: 15m
        alertname: "测试数字相关模板函数"
```

执行 `promtool test rules unit_test_alert.yaml` 即可进行测试

> [!Tip] 如果想要检查 Prometheus 规则中模板的渲染结果，只需要单元测试文件中的 exp_alerts 字段指定的期望结果 与 真实结果不一样即可

