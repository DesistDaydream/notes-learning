---
title: Alertmanager 配置
---

# 概述

> 参考：
>
> - [官方文档，告警 - 配置](https://prometheus.io/docs/alerting/latest/configuration/)

# Alertmanager 配置文件

下文用到的占位符说明：

- BOOLEAN # 可以采用 true 或 false 值的布尔值
- DURATION # 持续时间。可以使用正则表达式
  - `((([0-9]+)y)?(([0-9]+)w)?(([0-9]+)d)?(([0-9]+)h)?((([0-9]+)m)?((([0-9]+)s)?((([0-9]+)ms)?|0)`，例如：1d、1h30m、5m、10s。
- FILENAME # 当前工作目录中的有效路径
- HOST # 由主机名或 IP 后跟可选端口号组成的有效字符串。
- INT # 一个整数值
- LABELNAME # 与正则表达式 `[a-zA-Z _] [a-zA-Z0-9 _] *` 匹配的字符串
- LABELVALUE # 一串 unicode 字符
- PATH # 有效的 URL 路径
- SCHEME # 一个字符串，可以使用值 http 或 https
- SECRET # 作为机密的常规字符串，例如密码
- STRING # 常规字符串
- TMPL_STRING # 使用前已模板扩展的字符串

## 顶层字段

- **global**([global](#global)) # 全局配置，所有内容作用于所有配置环境中,若其余配置环境中不再指定同样的配置，则 global 中的配置作为默认配置
- **templates**(\[][templates](#templates)) # 指定告警模板文件的路径。若不指定则使用默认模板。可以使用通配符，e.g. 'templates/*.tmpl'
- **route**([route](#route)) # 路由树的节点(Alertmanager 的主要配置)
- **receivers**(\[][receivers](#receivers)) # 告警的接收者列表(Alertmanager 的主要配置)
- **inhibit_rules**([inhibit_rules](#inhibit_rules)) # 抑制规则列表

## global

全局配置，其内的内容作用于所有配置环境中,若其余配置环境中不再指定同样的配置，则 global 中的配置作为默认配置。在这里可以定义告警发送者的信息，比如通过邮件发送告警，那么可以定义全局的 SMTP 配置。

**resolve_timeout**(DURATION) # 如果接收到的告警不包括 EndsAt 字段，那么经过 resolve_timeout 时间后，如果没有重复收到告警，则认为该告警已解决。默认5m。

**http_config**(http_config) # The default HTTP client configuration

### SMTP 相关配置

**smtp_from**: <tmpl_string> # The default SMTP From header field.

**smtp_smarthost**(STRING) # The default SMTP smarthost used for sending emails, including port number.Port number usually is 25, or 587 for SMTP over TLS (sometimes referred to as STARTTLS).

- Example: smtp.example.org:587
 
**smtp_hello**(STRING) # The default hostname to identify to the SMTP server. `默认值: localhost`

**smtp_auth_username**(STRING) # SMTP Auth using CRAM-MD5, LOGIN and PLAIN. If empty, Alertmanager doesn't authenticate to the SMTP server.

**smtp_auth_password**(STRING) # SMTP Auth using LOGIN and PLAIN.

**smtp_auth_identity**(STRING) # SMTP Auth using PLAIN.

**smtp_auth_secret**(STRING) # SMTP Auth using CRAM-MD5.

**smtp_require_tls**(BOOLEAN) # The default SMTP TLS requirement. Note that Go does not support unencrypted connections to remote SMTP endpoints. `默认值: true`

### 用于 Slack 通知的 API URL

**slack_api_url**(secret)

**victorops_api_key**(secret)

**victorops_api_url**(string) | default = "https://alert.victorops.com/integrations/generic/20131114/alert/"

**pagerduty_url**(string) | default = "https://events.pagerduty.com/v2/enqueue"

**opsgenie_api_key**(secret)

**opsgenie_api_url**(string) | default = "https://api.opsgenie.com/"

**hipchat_api_url**(string) | default = "https://api.hipchat.com/"

**hipchat_auth_token**(secret)

**wechat_api_url**(string) | default = "https://qyapi.weixin.qq.com/cgi-bin/"

**wechat_api_secret**(secret)

**wechat_api_corp_id**(string)

## templates

用于定义接收人收到的告警的样式。如 HTML 模板、邮件模板等等。

## route

该字段中，可以配置多个相同接收者的子路由。在一个路由树中，将每个被路由的目标称为 **Node(节点)。**

**group_by**(\[]STRING) # 告警分组策略，凡是具有 \[]STRING 指定的标签名的告警都分为同一个组。

- 可以使用 `group_by: ['...']` 配置禁用聚合功能。

**group_interval**(DURATION) # 发送告警的间隔时间。`默认值：5m`。

**group_wait**(DURATION) # 发送告警前，需要等待分组行为的时间。`默认值：30s`

- 新收到的告警会进行分组聚合，并以组的形式发送，为了方式大量告警频繁触发告警发送，所以有一个等待期，等到多个告警聚合在一个组时，统一发送

**matchers**(\[][matcher](#matcher)) # 匹配规则，凡是符合该规则的告警，将会进入当前节点。说白了，只有匹配上了，才会将告警发出去。

> [!Tip]
> - 如果多个 Label 是“或”的关系，那就只能配置多个相同接收者的路由，每个路由的 matchers 不同。
> - `matchers` 字段代替了在 0.22.0 版本开始被弃用的 `match` 与 `match_re` 字段

**recevier**(STRING) # 当前路由匹配到的告警的接收者。如果 recevier 是整个路由树的根，则就是默认接收者

**routes**(\[][route](#route)) # 子路由配置。`routes` 字段的中的每个元素都是 `route` 字段 。

- 也就是说，`routes` 下每个元素的字段，与 `route` 下的字段相同，这是一个嵌套循环~~

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
      matchers:
        - network_device = "interface-state"
    - receiver: "dev-phone-group"
      group_wait: "10s"
      matchers:
        - instance =~ "dev-phone.*"
    - receiver: "dev-phone-group"
      group_wait: "10s"
      matchers:
        - job =~ "snmp-metrics"
```

## receivers

**receivers(接收者)** 是一个抽象的概念，可以是一个邮箱，也可以是微信、Slack、Webhook 等等。receivers 与 route 配置，根据路由规则将告警发送给指定的接收人。

```yaml
# 指定接收者的名称，用于在 route 配置环境中根据路由规则指定具体的接收者。
name(STRING)
# 不同的接收者有不同的配置环境。
XXXXX_configs:
- <详见下文对应配置环境>, ...
```

现阶段 alertmanager 支持的 XXXX_configs 有 email_configs、pagerduty_configs、pushover_configs、slack_configs、opsgenie_configs、webhook_configs、victorops_configs、wechat_configs。

### email_configs - 邮件 接收者

配置示例

```yaml
receivers:
- name: "default"
  email_configs:
  - to: "desistdaydream@wisetv.com.cn"
    send_resolved: true
```

### webhook_configs - webhook 接收者

webhook 类型的接收者是一种通用的接收者，不像其他类型的接收者，只能发送给特定的服务。而 webhook 只需要指定接收消息的 IP:PORT 即可。Alertmanager 会将指定的消息体已 POST 方法发送给对方，不管对方是什么，只要能处理 Alertamanger 发过去的 JSON 结构的数据即可。

```yaml
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
```

## inhibit_rules

抑制规则配置

**target_matchers**(\[][matcher](#matcher)) # 要静音的目标警报必须满足的匹配器列表。


# 配置文件中的通用配置字段

## matcher

https://prometheus.io/docs/alerting/latest/configuration/#matcher

matcher 部分的配置实现了 [Label matchers](docs/6.可观测性/Metrics/Prometheus/Label%20matchers.md)(标签匹配器) 的功能，根据这部分字段定义的 Label，正确匹配到目标（目标通常指一条告警）。

matcher(匹配器) 的逻辑类似 [PromQL](docs/6.可观测性/Metrics/Prometheus/PromQL/PromQL.md) 的各种 Selectors(选择器)，语法也是类似的，由 3 部分组成

- KEY
- `=`, `!=`, `=~` 三者之一
- VALUE

比如

```yaml
alert_target="test"
```

表示要匹配 KEY 为 alert_tagert；VALUE 为 test 的所有告警。

```yaml
# routes 字段下的 matcher 用法，将 KEY 为 severity，VALUE 为 warning 的告警交给 webhook-warn 这个接收者处理
  routes:
    - matchers:
        - severity = warning
      receiver: "webhook-warn"
# inhibit_rules 字段下 matcher 用法，将 KEY 为 alert_target，VALUE 为 test 的告警静音
inhibit_rules:
  - target_matchers:
      - alert_target="test"
```

TODO: matchers 下多个元素是和还是或的关系？

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
