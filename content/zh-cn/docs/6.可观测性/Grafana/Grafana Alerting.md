---
title: Grafana Alerting
linkTitle: Grafana Alerting
weight: 20
---

# 概述

> 参考：
>
> - [官方文档，告警](https://grafana.com/docs/grafana/latest/alerting/)

**Alert rules(警报规则)** # 确定是否触发警报的规则。e.g. 评估规则的周期（i.e. 多久检查一次告警是否应该触发）、etc.

**Notification policies(通知策略)** # 确定警报如何路由到联络点。e.g. 根据什么条件分组、发送前等待多久、相同警报重新发送的时间间隔、etc.

**Contack points(联络点)** # 当警报实例触发时如何通知联系人。e.g. 设置通过企业微信，使用 XXX 模板渲染消息，并发送给联系人。联系人是真实世界的实体，联络点是通知方式的抽象。

# Alert rules

一个完整的 Alert rules 由三部分组成

- **Query(查询语句)**
- **Condition Expressions(条件表达式)** # 可选。简称 Expressions
- **Evaluation behavior(评估行为)** # [评估行为](#Evaluation%20behavior) 与 Prometheus 中的 [Alerting](/docs/6.可观测性/Metrics/Prometheus/Alerting.md) 的评估行为类似（有一点不同是：Prom 的评估是决定是否将警报发送出去，但 Grafana 内部集成了类似 Alertmanager 的逻辑，所以评估是决定是否将警报交给 [Notifications](#notifications) 组件）

Grafana 会评估 Expressions 的处理结果（对 Query 查询结果的处理结果），满足条件的将会交给 Notifications 组件。

## Evaluation behavior

https://grafana.com/docs/grafana/latest/alerting/fundamentals/alert-rule-evaluation/

为 Alert rules 定义 **Evaluation behavior(评估行为)**，评估成功后，会改变警报的状态，并将某些状态的警报交给 Notification policies(通知策略) 以便对警报进行后续处理。

等待多久、间隔多久、查询出错或查询没数据时怎么办、etc. 都属于评估行为。包括如下几个部分

- **Evaluation name** # 评估行为名称
- **Evaluation interval** # 每次评估行为的间隔
- **Period during** # 满足 Alert rules 中定义的条件后等待多长时间触发警报
- **No data and Error handing** # 当 Alert rules 没有数据或者执行错误时的处理方式
  - <font color="#ff0000">Notes: 截至 11.5.1 版本，默认的 No data 处理方式是 No Data，这会导致每隔默认的 4h 就发送一次警报</font>，最好将 No data 的行为改为 Normal。

![600](https://notes-learning.oss-cn-beijing.aliyuncs.com/grafana/alerting/20250213083518945.png)

# Notifications

https://grafana.com/docs/grafana/latest/alerting/fundamentals/notifications/

Notifications 由下面几部分组成

- Notification policies
- Contact points

开始定义您的联系点，以指定如何接收您的警报通知。然后，配置您的警报规则将其警报发送到联系点或使用通知策略树以灵活地将警报路由到联系点。

## Notification policies

## Contact points

# Template

> 参考：
>
> - https://grafana.com/docs/grafana/latest/alerting/fundamentals/templates/

可以在下面这些地方使用模板

- Alert rule annotations
- Alert rule labels
- Notifcation templates

获取某个 Query(查询) 或 Condition Expressions(条件表达式) 的值

```go
// 获取表达式 A 的值
{{ index $values "A" }}
// 或
{{ $values.A.Value }}
```

# 警报通知数据结构

https://grafana.com/docs/grafana/latest/alerting/configure-notifications/template-notifications/reference/

# Grafana 8.0 告警使用

> 参考：
>
> - [公众号 - k8s技术圈，Grafana 8.0 告警使用](https://mp.weixin.qq.com/s/1aJOqhGOXaOas2lPUcP2-g)

Grafana 除了支持丰富的数据源和图表功能之外，还支持告警功能，该功能也使得 Grafana 从一个数据可视化工具成为了一个真正的监控利器。Grafana 可以通过 Alerting 模块的配置把监控数据中的异常信息进行告警，告警的规则可以直接基于现有的数据图表进行配置，在告警的时候也会把出现异常的图表进行通知，使得我们的告警通知更加友好。

## 渠道

Grafana Alerting 支持多种告警渠道，比如钉钉、Discord、Email、Kafka、Pushover、Telegram、Webhook 等等，我们这里可以使用钉钉和 Email 进行展示说明。

### Email

邮箱告警通常是最常见的告警接收方式，通过 Grafana 告警需要在 Grafana 的配置文件中配置 stmp 服务。在配置文件 `/etc/grafana/grafana.ini` 文件中添加 `SMTP/Emailing` 配置块并开启 `Alerting`：

```ini
#################################### SMTP / Emailing ##########################
[smtp]
enabled = true
host = smtp.163.com:465  # 我们这里使用163的邮箱
user = xxx@163.com
password = <email password>  # 使用网易邮箱的授权码
skip_verify = true
from_address = xxx@163.com

#################################### Alerting ############################
[alerting]
enabled = true
execute_alerts = true
```

需要注意的是这里我们使用的是 163 的邮箱进行发送，在配置 `smtp` 的时候需要在邮箱中开启 `IMAP/SMTP` 和 `POP3/SMTP` 两个服务，并添加一个授权码，上面的 password 密码使用的就是授权码进行认证：

![](https://mmbiz.qpic.cn/mmbiz_png/z9BgVMEm7YtaXicIicsE3YsiatkqOAgK4z0RfekRDLvica25up5ia9acEPZ3fMc2h9dXjiac6lzv4eHcNvsgPzWu2anw/640?wx_fmt=png)

配置完成后重新启动 Grafana：

`☸ ➜ systemctl daemon-reload
☸ ➜ systemctl restart grafana-server
`

回到 Grafana 页面中点击左侧的 `Notification channels` 开始添加消息通知渠道：

![](https://mmbiz.qpic.cn/mmbiz_png/z9BgVMEm7YtaXicIicsE3YsiatkqOAgK4z0lwlC33pOjfeFrRMDr9g4UOIg8J9gah7CicoMdrfM3cPmiaPHu9NvvxcQ/640?wx_fmt=png)

点击 `Add channel` 按钮新建一个通知渠道，这里我们选择渠道类型为 `Email`，添加接收通知的邮件地址，此外还可以对通知进行简单的配置：

![](https://mmbiz.qpic.cn/mmbiz_png/z9BgVMEm7YtaXicIicsE3YsiatkqOAgK4z0MLiaztiaF8NMFxesqceKJUDyDJFYa6kUcxNYNlp5B3yn3fVM7k0XibADw/640?wx_fmt=png)

点击下方的 `Test` 按钮可以测试是否可以正常发送邮件，如果出现 `Test notification sent` 的提示证明发送成功，正常也可以收到一封如下所示的告警通知邮件：

![](https://mmbiz.qpic.cn/mmbiz_png/z9BgVMEm7YtaXicIicsE3YsiatkqOAgK4z0MC6icVz8JudR420YvbCfrJDTLRW2KGRmBMXknPJqhicq7JfdQle4jEjA/640?wx_fmt=png)

测试成功后点击 `Save` 按钮，保存这个通知渠道。

### 钉钉

Grafana 还内置支持了钉钉，所以如果我们想把告警消息接入钉钉群也是非常方便的。创建一个自定义群机器人，需要注意的是现在的钉钉群机器人新增了 3 种安全认证方式，这里我们选择关键字的方式即可，设置关键字 `alert` 即可：

![](https://mmbiz.qpic.cn/mmbiz_png/z9BgVMEm7YtaXicIicsE3YsiatkqOAgK4z0jIOUMuUXMMwP4L2KYtpS7wJvvSYibnHKSInsjL8hpdq8FaWQ7XRaBHA/640?wx_fmt=png)

创建后会生成一个 Webhook 的地址，复制该地址：

![](https://mmbiz.qpic.cn/mmbiz_png/z9BgVMEm7YtaXicIicsE3YsiatkqOAgK4z05fyduxX0syj0ibeX9UUS4Ogvk81ByHHhGiaA2EicQfrxQiaaiaG64ialC0hg/640?wx_fmt=png)

然后回到 Grafana 中新建一个新的通知渠道，选择类型为 `DingDing`，将上面复制的 Webhook 地址拷贝到 `Url` 栏目中：

![](https://mmbiz.qpic.cn/mmbiz_png/z9BgVMEm7YtaXicIicsE3YsiatkqOAgK4z0sRPpTrMTLfjqSwwoQBicicAUpk7qpnbicicSTECrxxLVu6P2k7IDibKz2Gw/640?wx_fmt=png)

同样点击 `Test` 按钮可以测试消息：

![](https://mmbiz.qpic.cn/mmbiz_png/z9BgVMEm7YtaXicIicsE3YsiatkqOAgK4z0R4cfiaibVLATf5BCTQ0NDVCkmLEX4QwlRx1f58ZAsGcRiau2KX801NvUA/640?wx_fmt=png)

测试通过后点击 `Save` 保存该通知渠道即可。这样我们就创建了两个通知渠道，也可以根据需要设置一个默认的渠道，如果还有其他的渠道需求，可以自行添加即可。

## 规则

在通知渠道的左侧就是一个 `Alert rules` 告警规则的选项卡，点击该页面下面的 `How to add an alert` 按钮就有提示如果创建一个告警：

![](https://mmbiz.qpic.cn/mmbiz_png/z9BgVMEm7YtaXicIicsE3YsiatkqOAgK4z0W4A1MTD1GV8npT0LATgbW8z1IOibKxn3RmGcibrMg2rWHazagiciaSrw4g/640?wx_fmt=png)

提示非常清晰，在任何仪表板图形面板的 `Alert` 选项卡中添加和配置告警即可，可以使用现有查询构建和可视化告警，所以我们需要前往仪表板的图形面板中进行配置，这里我们同样以 CPU 使用率这个面板进行说明。

进入 CPU 使用率的面板编辑页面：

![](https://mmbiz.qpic.cn/mmbiz_png/z9BgVMEm7YtaXicIicsE3YsiatkqOAgK4z0yAF2LGgmQiaQZeJsxSUyaKqpK9mDTE4pibY25ia7VvO7TBibJwIrpfQbvQ/640?wx_fmt=png)

在编辑页面图形下方有 `Query`、`Transform`、`Alert` 三个选项卡，前两个我们都已经使用过了，这里需要使用到的是 `Alert` 这个选项卡：

![](https://mmbiz.qpic.cn/mmbiz_png/z9BgVMEm7YtaXicIicsE3YsiatkqOAgK4z0CT0RFrju0icNb2icGnSbvhOJXjxP3FO8fffFMZU5ANgeXAURevkRgkbg/640?wx_fmt=png)

但是我们切换到 `Alert` 选项卡页面的时候出现了 `Template variables are not supported in alert queries` 这样的提示信息，意思就是在告警查询中是不支持模板变量的，但是我们这里的图形查询中定义了好几个变量，应该怎么处理该问题呢？

首先我们需要在该面板中添加一个没有变量的查询语句，该语句用于报警使用，比如我们对节点总的使用率进行监控报警，添加新的查询语句 `(1 - sum(rate(node_cpu_seconds_total{instance=~"node1:9100", mode="idle"}[5m])) by (instance) / sum(rate(node_cpu_seconds_total{instance=~"node1:9100"}[5m])) by (instance) ) * 100`，去掉对节点参数的使用，因为 Grafana 的报警不支持多维数据，所以这里我们暂时只对 `node1` 节点进行监控，然后需要将该查询设置成 `Disable query`，这样图表中就不会有该指标数据了，因为该指标是用来监控报警的：

![](https://mmbiz.qpic.cn/mmbiz_png/z9BgVMEm7YtaXicIicsE3YsiatkqOAgK4z0fZm2OffTP5nibCJhaooy8NvX2guwia51CQFTwzpib2UOTYW5Wt46iaibyKg/640?wx_fmt=png)

现在我们再切换到 `Alert` 选项卡页面就可以正常创建报警规则了：

![](https://mmbiz.qpic.cn/mmbiz_png/z9BgVMEm7YtaXicIicsE3YsiatkqOAgK4z0uxicTQODUoDrUe1undfUckkjXqmhZXicXSnhGLlDnhsObD1EWnZPUb0w/640?wx_fmt=png)

点击 `Create Alert` 按钮创建报警规则：

![](https://mmbiz.qpic.cn/mmbiz_png/z9BgVMEm7YtaXicIicsE3YsiatkqOAgK4z0icdiagkticFRjQ6bCfiaRiaDwAIXhMdW1nPGXRhESq0Ety9VFibRcsbqsmpg/640?wx_fmt=png)

在 Rule Name 中，添加一个描述性名称，该名称显示在警报规则列表中，后面的 `Evaluate every` 表示的评估时间，这里我们设置 `1m`，表示每隔 1 分钟 Grafana 会来评估我们的报警规则，`For` 表示的是 `Pending Duration` 的时长，意思就是如果报警规则持续 `1m` 的时间则表示要真正去触发报警了。

然后就是配置报警的条件，在 `WHEN` 后面可以点击选择各种计算方式，我们这里选择 `avg()` 表示平均值，`OF` 后面的查询就是我们真正用于监控报警的语句，点击可以选择用于查询的语句，这里我们需要选择上面新建的语句 `D`，`query(D, 5m, now)` 就表示语句 D 从现在开始的前 5 分钟内平均值大于（IS ABOVE）1 这个阈值。

在下方还可以配置用于报警通知的渠道：

![](https://mmbiz.qpic.cn/mmbiz_png/z9BgVMEm7YtaXicIicsE3YsiatkqOAgK4z07s98S0iaKAW3URlFDBktyfNLUPo6dMHKeibkdWFCPj03KSvf5HqXl9PQ/640?wx_fmt=png)

然后配置好过后 `Apply` 该面板并保存 Dashboard，正常隔一会儿就可以收到报警通知了：

![](https://mmbiz.qpic.cn/mmbiz_png/z9BgVMEm7YtaXicIicsE3YsiatkqOAgK4z0o1aWM2bNypzEM9WRHFOus80fHT6nZlkdcOP7TUhfgjOv4bzoeBsNicg/640?wx_fmt=png)

但是通知里面的图形并没有渲染出来，这是因为我们没有安装 `grafana-image-renderer` 插件，在 Grafana 安装节点上执行下面命令安装：

`☸ ➜ grafana-cli plugins install grafana-image-renderer

✔ Downloaded grafana-image-renderer v3.3.0 zip successfully

Please restart Grafana after installing plugins. Refer to Grafana documentation for instructions if necessary.
☸ ➜ systemctl restart grafana-server

`

再一次触发报警的时候可能还是不会正常渲染图形，查看 `Grafana` 的日志可以了解到相关错误信息：

```
☸ ➜ journalctl -u grafana-server -f
......
Nov 30 18:19:01 node2 grafana-server[62536]: t=2021-11-30T18:19:01+0800 lvl=eror msg="Render request failed" logger=plugins.backend pluginId=grafana-image-renderer url="http://localhost:3000/d-solo/oq26nAFnz/nodejie-dian-jian-kong?orgId=1&panelId=2&render=1" error="Error: Failed to launch the browser process!\n/var/lib/grafana/plugins/grafana-image-renderer/chrome-linux/chrome: error while loading shared libraries: libatk-1.0.so.0: cannot open shared object file: No such file or directory\n\n\nTROUBLESHOOTING: https://github.com/puppeteer/puppeteer/blob/main/docs/troubleshooting.md\n"
......
```

要解决这个问题我们需要安装几个 `puppeteer` 的依赖包：

`☸ ➜ yum install atk at-spi2-atk libxkbcommon-x11-devel libXcomposite gtk3 -y
`

依赖安装完成后正常收到的告警消息通知就包含图形数据了：

![](https://mmbiz.qpic.cn/mmbiz_png/z9BgVMEm7YtaXicIicsE3YsiatkqOAgK4z027ibzcTGzLbvAmZ00xdYEpbC9lDCyxrWzV255NqWytDHvia4BsfIP3PA/640?wx_fmt=png)

只是渲染的图形中文是乱码，这主要是 Linux 字体库对中文支持不好的原因，我们只需要给服务器的 Linux 系统安装支持的中文字体库即可，这里我们安装文泉驿字体库：

```
☸ ➜ yum search wqy
Loaded plugins: fastestmirror
Loading mirror speeds from cached hostfile
 * base: mirrors.aliyun.com
 * epel: mirrors.bfsu.edu.cn
 * extras: mirrors.aliyun.com
 * updates: mirrors.aliyun.com
================================================ N/S matched: wqy =================================================
wqy-microhei-fonts.noarch : Compact Chinese fonts derived from Droid
wqy-unibit-fonts.noarch : WenQuanYi Unibit Bitmap Font
wqy-zenhei-fonts.noarch : WenQuanYi Zen Hei CJK Font

  Name and summary matches only, use "search all" for everything.

☸ ➜ yum install wqy-microhei-fonts.noarch wqy-unibit-fonts.noarch wqy-zenhei-fonts.noarch -y
```

这个时候渲染的图形就可以正常显示了：

![](https://mmbiz.qpic.cn/mmbiz_png/z9BgVMEm7YtaXicIicsE3YsiatkqOAgK4z0h6j6R0tKz2OVDiaCH4AXJxtF9o42xE02SYwPiaicdUkria6OAmbaMILSpw/640?wx_fmt=png)

但是钉钉通知中没有将图形显示出来：

![](https://mmbiz.qpic.cn/mmbiz_png/z9BgVMEm7YtaXicIicsE3YsiatkqOAgK4z0KwN7HWQlhNpiboSnicOnf47kkg3r8nnepKiaa8KUvEd0kIzwS5qZV0uNA/640?wx_fmt=png)

这是因为在邮件告警通知中的图片是通过邮件服务器发送出去的，是一个图片附件，所以可以正常看到，但是对于钉钉、webhook 这些告警渠道则是直接显示的图片，所以需要设置 `external storage` 才能进行显示，设置后图片变为一个指向 `external_image_storage` 中图片地址的链接，如果本机无法访问 `external_image_storage` 图片是无法显示的，所以如果设置成 local 的时候邮件中的图片可能也不能显示了，因为这个时候是直接一张图片链接。不过 `external_image_storage` 可设置的 provider 包括 s3、webdav、gcs、azure\_blob、local，如果是线上服务的话建议接入对象存储，比如 s3 服务，这里我们使用阿里云 OSS 来配置 s3 进行说明。

在 Grafana 配置文件 `/etc/grafana/grafana.ini` 中配置 `external_image_storage`：

```
#################################### External image storage ##########################
[external_image_storage]
provider = s3  # 使用 s3 模式

[external_image_storage.s3]
endpoint = oss-cn-beijing.aliyuncs.com
bucket = <bucket>
region = oss-cn-beijing
access_key = <ak>  # 使用阿里云后台的ak和sk进行配置
secret_key = <sk>

#################################### Server ####################################
[server]
domain = 192.168.31.46  # 设置 Grafana 访问地址为内网 IP
```

另外注意需要将 Grafana 的访问域名设置成内网 IP，否则在局域网其他节点上访问不到，配置完成后重启 Grafana 即可：

```
☸ ➜ systemctl daemon-reload
☸ ➜ systemctl restart grafana-server
```

配置完成后我们重新去触发下报警，正常在邮件和钉钉中收到的图片都可以正常显示了：

![](https://mmbiz.qpic.cn/mmbiz_png/z9BgVMEm7YtaXicIicsE3YsiatkqOAgK4z0PLlkfQxVOor0gvMf79ricofpVk5lEgmhlibezkQj5B5NaLiaI5S6g4klQ/640?wx_fmt=png)
