---
title: Alertmanager 配置
---

# 概述

> 参考：
> - [官方文档，告警-配置](https://prometheus.io/docs/alerting/latest/configuration/)

# Alertmanager 配置文件

下文用到的占位符说明：

- \<DURATION> # 与正则表达式匹配的持续时间  \[0-9]+(ms|\[smhdwy])
- \<LabelName> # 与正则表达式匹配的字符串  \[a-zA-Z\_]\[a-zA-Z0-9\_]\*
- \<LabelValue> # 一串 unicode 字符
- \<FilePath> # 当前工作目录中的有效路径
- \<BOOLEAN> # 可以接受值的布尔值，true 或 false
- \<STRING> # 一个普通的字符串
- \<SECRET> # 是秘密的常规字符串，例如密码
- \<tmpl_string> # a string which is template-expanded before usage
- \<tmpl_secret> # a string which is template-expanded before usage that is a secret

下面是一个配置文件的基本结构

```yaml
# 全局配置，所有内容作用于所有配置环境中,若其余配置环境中不再指定同样的配置，则global中的配置作为默认配置
global:
  # The default SMTP From header field.
  [ smtp_from: <tmpl_string> ]
  # The default SMTP smarthost used for sending emails, including port number.
  # Port number usually is 25, or 587 for SMTP over TLS (sometimes referred to as STARTTLS).
  # Example: smtp.example.org:587
  [ smtp_smarthost: <string> ]
  # The default hostname to identify to the SMTP server.
  [ smtp_hello: <string> | default = "localhost" ]
  # SMTP Auth using CRAM-MD5, LOGIN and PLAIN. If empty, Alertmanager doesn't authenticate to the SMTP server.
  [ smtp_auth_username: <string> ]
  # SMTP Auth using LOGIN and PLAIN.
  [ smtp_auth_password: <secret> ]
  # SMTP Auth using PLAIN.
  [ smtp_auth_identity: <string> ]
  # SMTP Auth using CRAM-MD5.
  [ smtp_auth_secret: <secret> ]
  # The default SMTP TLS requirement.
  # Note that Go does not support unencrypted connections to remote SMTP endpoints.
  [ smtp_require_tls: <bool> | default = true ]

  # The API URL to use for Slack notifications.
  [ slack_api_url: <secret> ]
  [ victorops_api_key: <secret> ]
  [ victorops_api_url: <string> | default = "https://alert.victorops.com/integrations/generic/20131114/alert/" ]
  [ pagerduty_url: <string> | default = "https://events.pagerduty.com/v2/enqueue" ]
  [ opsgenie_api_key: <secret> ]
  [ opsgenie_api_url: <string> | default = "https://api.opsgenie.com/" ]
  [ hipchat_api_url: <string> | default = "https://api.hipchat.com/" ]
  [ hipchat_auth_token: <secret> ]
  [ wechat_api_url: <string> | default = "https://qyapi.weixin.qq.com/cgi-bin/" ]
  [ wechat_api_secret: <secret> ]
  [ wechat_api_corp_id: <string> ]

  # The default HTTP client configuration
  [ http_config: <http_config> ]

  # 如果接收到的告警不包括 EndsAt 字段，那么经过 resolve_timeout 时间后，如果没有重复收到告警，则认为该告警已解决。默认5m。
  resolve_timeout: <DURATION>

# 指定告警模板文件的路径。若不指定则使用默认模板。可以使用通配符，e.g. 'templates/*.tmpl'
templates:
  [ - <filepath> ... ]

# 路由树的节点(Alertmanager 的主要配置)
route:
  详见下文单独章节
  routes:
  - 详见下文单独章节

# 告警的接收者列表(Alertmanager 的主要配置)
receivers:
- <receiver> ...

# 抑制规则列表
inhibit_rules:
  [ - <inhibit_rule> ... ]
```

## global(OBJECT)

全局配置，其内的内容作用于所有配置环境中,若其余配置环境中不再指定同样的配置，则 global 中的配置作为默认配置。在这里可以定义告警发送者的信息，比如通过邮件发送告警，那么可以定义全局的 SMTP 配置。

## templates([]OBJECT)

用于定义接收人收到的告警的样式。如 HTML 模板、邮件模板等等。

### 配置示例

```yaml
templates: [- <filepath> ...]
```

## route(OBJECT)

该字段中，可以配置多个相同接收者的子路由。在一个路由树中，将每个被路由的目标称为 **Node(节点)。**
**group_by([]STRING)** # 告警分组策略，凡是具有 \[]STRING 指定的标签名的告警都分为同一个组。

- 可以使用 `group_by: ['...']` 配置禁用聚合功能。

**group_interval(DURATION)** # 发送告警的间隔时间。`默认值：5m`。
**group_wait(DURATION)** # 发送告警前，需要等待分组行为的时间。`默认值：30s`

- 新收到的告警会进行分组聚合，并以组的形式发送，为了方式大量告警频繁触发告警发送，所以有一个等待期，等到多个告警聚合在一个组时，统一发送

**matchers([]OBJECT)** # 匹配规则，凡是符合该规则的告警，将会进入当前节点。说白了，只有匹配上了，才会将告警发出去。

> 注意：
>
> - 如果多个 Label 是“或”的关系，那就只能配置多个相同接收者的路由，每个路由的 matchers 不同。
> - `matchers` 字段代替了在 0.22.0 版本开始被弃用的 `match` 与 `match_re` 字段

**recevier(STRING)** # 当前路由匹配到的告警的接收者。如果 recevier 是整个路由树的根，则就是默认接收者
**routes: []OBJECT** # 子路由配置。`routes` 字段的中的每个元素其实就是 `route: <OBJECT>`。

- 也就是说，`routes: <[]OBJECT>` 下每个元素的字段，与 `route: <OBJECT>` 下的字段相同，这是一个嵌套循环~~

### 配置示例

其中具有标签 instance=dev-phone.\* 和 job=snmp-metrics 的告警都会路由给名为 dev-phone-group 的接收者

```yaml
route:
  group_by:
    - "namespace"
  group_interval: "5m"
  group_wait: "30s"
  receiver: "default"
  repeat_interval: "6h"
  routes:
    - receiver: "network-group"
      group_wait: "10s"
      match:
        network_device: "interface-state"
    - receiver: "dev-phone-group"
      group_wait: "10s"
      match_re:
        instance: "dev-phone.*"
    - receiver: "dev-phone-group"
      group_wait: "10s"
      match_re:
        job: "snmp-metrics"
```

## receivers([]OBJECT)

**receivers(接收者)** 是一个抽象的概念，可以是一个邮箱，也可以是微信、Slack、Webhook 等等。receivers 与 route 配置，根据路由规则将告警发送给指定的接收人。

    # 指定接收者的名称，用于在 route 配置环境中根据路由规则指定具体的接收者。
    name(STRING)
    # 不同的接收者有不同的配置环境。
    XXXXX_configs:
    - \<详见下文对应配置环境>, ...

现阶段 alertmanager 支持的 XXXX_configs 有 email_configs、pagerduty_configs、pushover_configs、slack_configs、opsgenie_configs、webhook_configs、victorops_configs、wechat_configs。

### email_configs 字段。邮件 接收者

配置示例

    receivers:
    - name: "default"
      email_configs:
      - to: "lichenhao@wisetv.com.cn"
        send_resolved: true

### webhook_configs 字段。webhook 接收者

webhook 类型的接收者是一种通用的接收者，不像其他类型的接收者，只能发送给特定的服务。而 webhook 只需要指定接收消息的 IP:PORT 即可。Alertmanager 会将指定的消息体已 POST 方法发送给对方，不管对方是什么，只要能处理 Alertamanger 发过去的 JSON 结构的数据即可。

    # Whether or not to notify about resolved alerts.
    [ send_resolved: <boolean> | default = true ]

    # The endpoint to send HTTP POST requests to.
    url: <string>

    # The HTTP client's configuration.
    [ http_config: <http_config> | default = global.http_config ]

    # The maximum number of alerts to include in a single webhook message. Alerts
    # above this threshold are truncated. When leaving this at its default value of
    # 0, all alerts are included.
    [ max_alerts: <int> | default = 0 ]

## inhibit_rules: <OBJECT>

抑制规则配置

# 配置文件示例

```yaml
config:
  global:
    resolve_timeout: "5m"
  route:
    group_by:
      - "job"
      - "ssc_pool_type"
    group_wait: "30s"
    group_interval: "5m"
    repeat_interval: "6h"
    receiver: "webhook"
    routes:
      - repeat_interval: "5m"
        match:
          alertname: Watchdog
        receiver: "developer"
      - match:
          job: console-server-exporter
        receiver: "developer"
  receivers:
    - name: "webhook"
      webhook_configs:
        - url: "http://gateway.ssc.svc.cluster.local.:9010/alarmService/api/v1/alerts"
          send_resolved: true
        - url: "http://alertmanager-webhook-dingtalk.monitoring.svc.cluster.local.:8060/dingtalk/webhook1/send"
          send_resolved: true
    - name: "developer"
      webhook_configs:
        - url: "http://alertmanager-webhook-dingtalk.monitoring.svc.cluster.local.:8060/dingtalk/webhook_mention_users/send"
  templates:
    - "/etc/alertmanager/config/*.tmpl"
```
