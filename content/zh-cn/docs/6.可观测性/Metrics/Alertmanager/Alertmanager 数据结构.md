---
title: Alertmanager 数据结构
---

# 概述

> 参考：
>
> - [官方文档，告警 - 客户端](https://prometheus.io/docs/alerting/latest/clients/)(接收告警的数据结构)

## AlertManager 接收告警的数据结构

这个数据结构，其实就是 Prometheus 推送告警的数据结构.。。。详见 [Prometheus Alerting](/docs/6.可观测性/Metrics/Prometheus/Alerting(告警).md) 章节

```json
[
  {
    "labels": {
      "alertname": "<requiredAlertName>",
      "<labelname>": "<labelvalue>",
      ...
    },
    "annotations": {
      "<labelname>": "<labelvalue>",
    },
    "startsAt": "<rfc3339>",
    "endsAt": "<rfc3339>",
    "generatorURL": "<generator_url>"
  },
  ...
]
```

## Alertmanager 通过 Webhook 推送告警的数据结构

> 参考:
>
> - [官方文档，告警 - 配置 - webhook_config](https://prometheus.io/docs/alerting/latest/configuration/#webhook_config)(通过 Webhook 推送告警的数据结构)

下面就是 Alertmanager 在 Webhook 配置中，以 POST 请求发送的 JSON 结构的数据格式：

```json
{
  "version": "4",
  "groupKey": <string>,              // key identifying the group of alerts (e.g. to deduplicate)
  "truncatedAlerts": <int>,          // how many alerts have been truncated due to "max_alerts"
  // 当前发送的告警状态，如果是激活的告警则是 firing，如果是已经解决的则是 resolved
  "status": "<resolved|firing>",
  "receiver": <string>,
  // 根据 AlertManager 配置中 group_by 字段获取。
  // 也就是说通过某个标签分组，那么这个标签的名和值都会被写到这个字段中
  "groupLabels": <object>,
  // 这一组告警中，相同的 Labels 和 Annotations
  "commonLabels": <object>,
  "commonAnnotations": <object>,
  // 反向连接到 Alertmanager，
  "externalURL": <string>,
  "alerts": [
    {
      "status": "<resolved|firing>",
      "labels": {
        "alertname": "<RequiredAlertName>",
        "<labelname>": "<labelvalue>",
        ......
      },
      "annotations": <object>,
      // 告警的触发时间
      "startsAt": "<rfc3339>",
      // 若 status 字段的值为 resolved，则会显示该告警的结束时间
      "endsAt": "<rfc3339>",
      // 反向链接到 Prometheus
      // URL 包括 prometheus 域名和由生成告警的 promQL 语句。通过该连接，可以直接在 prometheus web 中打开使用指定 promql 查询。
      "generatorURL": <string>
    },
    ...
  ]
}
```

### 示例

```json
{
  "receiver": "webhook",
  "status": "firing",
  "alerts": [
    {
      "status": "firing",
      "labels": {
        "instance": "192.168.0.1:9100",
        "job": "node-exporter",
        "severity": "warning",
        "alertname": "测试告警1"
      },
      "annotations": {
        "description": "测试告警1的告警详情",
        "summary": "测试告警1概要"
      },
      "startsAt": "2021-04-24T15:22:27.944457098Z",
      "endsAt": "0001-01-01T00:00:00Z",
      "generatorURL": "http://desistdaydream:9090/graph?g0.expr=vector%281%29\u0026g0.tab=1",
      "fingerprint": "5adc80257c32889a"
    },
    {
      "status": "firing",
      "labels": {
        "instance": "192.168.0.1:9100",
        "severity": "warning",
        "alertname": "测试告警2"
      },
      "annotations": {
        "description": "测试告警2的告警详情",
        "summary": "测试告警2概要"
      },
      "startsAt": "2021-04-24T15:22:27.944457098Z",
      "endsAt": "0001-01-01T00:00:00Z",
      "generatorURL": "http://desistdaydream:9090/graph?g0.expr=vector%281%29\u0026g0.tab=1",
      "fingerprint": "e6532a92e438cdbf"
    }
  ],
  "groupLabels": {
    "alertname": "测试告警组"
  },
  "commonLabels": {
    "instance": "192.168.0.1:9100",
    "severity": "warning"
  },
  "commonAnnotations": {},
  "externalURL": "http://desistdaydream:9093",
  "version": "4",
  "groupKey": "{}:{instance=\"192.168.0.1:9100\"}",
  "truncatedAlerts": 0
}
```

```json
{
  "receiver": "webhook",
  "status": "firing",
  "alerts": [
    {
      "status": "firing",
      "labels": {
        "alertname": "测试告警3",
        "label_2": "value-1",
        "severity": "critical",
        "tenant": "test"
      },
      "annotations": {
        "additionalProp1": "string"
      },
      "startsAt": "2021-04-24T15:22:27.944457098Z",
      "endsAt": "0001-01-01T00:00:00Z",
      "generatorURL": "http://desistdaydream:9090/graph?g0.expr=vector%281%29\u0026g0.tab=1",
      "fingerprint": "496f742ac98e2398"
    }
  ],
  "groupLabels": {},
  "commonLabels": {
    "alertname": "测试告警3",
    "label_2": "value-1",
    "severity": "critical",
    "tenant": "test"
  },
  "commonAnnotations": {
    "additionalProp1": "string"
  },
  "externalURL": "http://desistdaydream:9093",
  "version": "4",
  "groupKey": "{}:{}",
  "truncatedAlerts": 0
}
```

```json
{
	"receiver": "webhook",
	"status": "resolved",
	"alerts": [{
		"status": "resolved",
		"labels": {
			"alert_event": "内存",
			"alert_target": "node",
			"alertname": "内存使用率过高！",
			"instance": "localhost:9100",
			"job": "node-exporter",
			"severity": "minor"
		},
		"annotations": {
			"description": "localhost:9100 内存持续一小时使用率大于 95% (目前可用:90.1%)",
			"summary": "内存使用率过高！"
		},
		"startsAt": "2024-06-28T03:42:37.186Z",
		"endsAt": "2024-06-28T03:43:22.186Z",
		"generatorURL": "http://bj-test-lichenhao-1:9090/graph?g0.expr=node_memory_MemAvailable_bytes+%2F+node_memory_MemTotal_bytes+%2A+100+%3C+0.5\u0026g0.tab=1",
		"fingerprint": "4a2564565982cb6d"
	}],
	"groupLabels": {
		"instance": "localhost:9100"
	},
	"commonLabels": {
		"alert_event": "内存",
		"alert_target": "node",
		"alertname": "内存使用率过高！",
		"instance": "localhost:9100",
		"job": "node-exporter",
		"severity": "minor"
	},
	"commonAnnotations": {
		"description": "localhost:9100 内存持续一小时使用率大于 95% (目前可用:90.1%)",
		"summary": "内存使用率过高！"
	},
	"externalURL": "http://bj-test-lichenhao-1:9093",
	"version": "4",
	"groupKey": "{}:{instance=\"localhost:9100\"}",
	"truncatedAlerts": 0
}
```