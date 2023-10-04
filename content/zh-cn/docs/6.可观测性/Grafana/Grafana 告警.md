---
title: Grafana 告警
linkTitle: Grafana 告警
date: 2023-12-10T11:10
weight: 20
---

# 概述

> 参考：
>
> - [Grafana 8.0 告警使用](https://mp.weixin.qq.com/s/1aJOqhGOXaOas2lPUcP2-g)

Grafana 除了支持丰富的数据源和图表功能之外，还支持告警功能，该功能也使得 Grafana 从一个数据可视化工具成为了一个真正的监控利器。Grafana 可以通过 Alerting 模块的配置把监控数据中的异常信息进行告警，告警的规则可以直接基于现有的数据图表进行配置，在告警的时候也会把出现异常的图表进行通知，使得我们的告警通知更加友好。

## 渠道

Grafana Alerting 支持多种告警渠道，比如钉钉、Discord、Email、Kafka、Pushover、Telegram、Webhook 等等，我们这里可以使用钉钉和 Email 进行展示说明。

### Email

邮箱告警通常是最常见的告警接收方式，通过 Grafana 告警需要在 Grafana 的配置文件中配置 stmp 服务。在配置文件 `/etc/grafana/grafana.ini` 文件中添加 `SMTP/Emailing` 配置块并开启 `Alerting`：

```ini
#################################### SMTP / Emailing ##########################
[smtp]
enabled = true
host = smtp.163.com:465  # 我们这里使用 163 的邮箱
user = <xxx@163.com>
password =   # 使用网易邮箱的授权码
skip_verify = true
from_address = <xxx@163.com>

#################################### Alerting ############################
[alerting]
enabled = true
execute_alerts = true
```

需要注意的是这里我们使用的是 163 的邮箱进行发送，在配置 `smtp` 的时候需要在邮箱中开启 `IMAP/SMTP` 和 `POP3/SMTP` 两个服务，并添加一个授权码，上面的 password 密码使用的就是授权码进行认证：

配置完成后重新启动 Grafana：

```bash
systemctl daemon-reload
systemctl restart grafana-server
```

回到 Grafana 页面中点击左侧的 `Notification channels` 开始添加消息通知渠道：

点击 `Add channel` 按钮新建一个通知渠道，这里我们选择渠道类型为 `Email`，添加接收通知的邮件地址，此外还可以对通知进行简单的配置：

点击下方的 `Test` 按钮可以测试是否可以正常发送邮件，如果出现 `Test notification sent` 的提示证明发送成功，正常也可以收到一封如下所示的告警通知邮件：

测试成功后点击 `Save` 按钮，保存这个通知渠道。

### 钉钉

Grafana 还内置支持了钉钉，所以如果我们想把告警消息接入钉钉群也是非常方便的。创建一个自定义群机器人，需要注意的是现在的钉钉群机器人新增了 3 种安全认证方式，这里我们选择关键字的方式即可，设置关键字 `alert` 即可：

创建后会生成一个 Webhook 的地址，复制该地址：

然后回到 Grafana 中新建一个新的通知渠道，选择类型为 `DingDing`，将上面复制的 Webhook 地址拷贝到 `Url` 栏目中：

同样点击 `Test` 按钮可以测试消息：

测试通过后点击 `Save` 保存该通知渠道即可。这样我们就创建了两个通知渠道，也可以根据需要设置一个默认的渠道，如果还有其他的渠道需求，可以自行添加即可。

## 规则

在通知渠道的左侧就是一个 `Alert rules` 告警规则的选项卡，点击该页面下面的 `How to add an alert` 按钮就有提示如果创建一个告警：

提示非常清晰，在任何仪表板图形面板的 `Alert` 选项卡中添加和配置告警即可，可以使用现有查询构建和可视化告警，所以我们需要前往仪表板的图形面板中进行配置，这里我们同样以 CPU 使用率这个面板进行说明。

进入 CPU 使用率的面板编辑页面：

在编辑页面图形下方有 `Query`、`Transform`、`Alert` 三个选项卡，前两个我们都已经使用过了，这里需要使用到的是 `Alert` 这个选项卡：

但是我们切换到 `Alert` 选项卡页面的时候出现了 `Template variables are not supported in alert queries` 这样的提示信息，意思就是在告警查询中是不支持模板变量的，但是我们这里的图形查询中定义了好几个变量，应该怎么处理该问题呢？

首先我们需要在该面板中添加一个没有变量的查询语句，该语句用于报警使用，比如我们对节点总的使用率进行监控报警，添加新的查询语句 `(1 - sum(rate(node_cpu_seconds_total{instance=~"node1:9100", mode="idle"}[5m])) by (instance) / sum(rate(node_cpu_seconds_total{instance=~"node1:9100"}[5m])) by (instance) ) * 100`，去掉对节点参数的使用，因为 Grafana 的报警不支持多维数据，所以这里我们暂时只对 `node1` 节点进行监控，然后需要将该查询设置成 `Disable query`，这样图表中就不会有该指标数据了，因为该指标是用来监控报警的：

现在我们再切换到 `Alert` 选项卡页面就可以正常创建报警规则了：

点击 `Create Alert` 按钮创建报警规则：

在 Rule Name 中，添加一个描述性名称，该名称显示在警报规则列表中，后面的 `Evaluate every` 表示的评估时间，这里我们设置 `1m`，表示每隔 1 分钟 Grafana 会来评估我们的报警规则，`For` 表示的是 `Pending Duration` 的时长，意思就是如果报警规则持续 `1m` 的时间则表示要真正去触发报警了。

然后就是配置报警的条件，在 `WHEN` 后面可以点击选择各种计算方式，我们这里选择 `avg()` 表示平均值，`OF` 后面的查询就是我们真正用于监控报警的语句，点击可以选择用于查询的语句，这里我们需要选择上面新建的语句 `D`，`query(D, 5m, now)` 就表示语句 D 从现在开始的前 5 分钟内平均值大于（IS ABOVE）1 这个阈值。

在下方还可以配置用于报警通知的渠道：

然后配置好过后 `Apply` 该面板并保存 Dashboard，正常隔一会儿就可以收到报警通知了：

但是通知里面的图形并没有渲染出来，这是因为我们没有安装 `grafana-image-renderer` 插件，在 Grafana 安装节点上执行下面命令安装：

```bash
~]# grafana-cli plugins install grafana-image-renderer

✔ Downloaded grafana-image-renderer v3.3.0 zip successfully

Please restart Grafana after installing plugins. Refer to Grafana documentation for instructions if necessary.
~]# systemctl restart grafana-server
```

再一次触发报警的时候可能还是不会正常渲染图形，查看 `Grafana` 的日志可以了解到相关错误信息：

`☸ ➜ journalctl -u grafana-server -f ...... Nov 30 18:19:01 node2 grafana-server[62536]: t=2021-11-30T18:19:01+0800 lvl=eror msg="Render request failed" logger=plugins.backend pluginId=grafana-image-renderer url="http://localhost:3000/d-solo/oq26nAFnz/nodejie-dian-jian-kong?orgId=1&panelId=2&render=1" error="Error: Failed to launch the browser process!\n/var/lib/grafana/plugins/grafana-image-renderer/chrome-linux/chrome: error while loading shared libraries: libatk-1.0.so.0: cannot open shared object file: No such file or directory\n\n\nTROUBLESHOOTING: https://github.com/puppeteer/puppeteer/blob/main/docs/troubleshooting.md\n" ......`

要解决这个问题我们需要安装几个 `puppeteer` 的依赖包：

`☸ ➜ yum install atk at-spi2-atk libxkbcommon-x11-devel libXcomposite gtk3 -y`

依赖安装完成后正常收到的告警消息通知就包含图形数据了：

只是渲染的图形中文是乱码，这主要是 Linux 字体库对中文支持不好的原因，我们只需要给服务器的 Linux 系统安装支持的中文字体库即可，这里我们安装文泉驿字体库：

```bash
☸ ➜ yum search wqy
Loaded plugins: fastestmirror
Loading mirror speeds from cached hostfile
 *base: mirrors.aliyun.com
 * epel: mirrors.bfsu.edu.cn
 *extras: mirrors.aliyun.com
 * updates: mirrors.aliyun.com
\================================================ N/S matched: wqy =================================================
wqy-microhei-fonts.noarch : Compact Chinese fonts derived from Droid
wqy-unibit-fonts.noarch : WenQuanYi Unibit Bitmap Font
wqy-zenhei-fonts.noarch : WenQuanYi Zen Hei CJK Font

Name and summary matches only, use "search all" for everything.

☸ ➜ yum install wqy-microhei-fonts.noarch wqy-unibit-fonts.noarch wqy-zenhei-fonts.noarch -y
```

这个时候渲染的图形就可以正常显示了：

但是钉钉通知中没有将图形显示出来：

这是因为在邮件告警通知中的图片是通过邮件服务器发送出去的，是一个图片附件，所以可以正常看到，但是对于钉钉、webhook 这些告警渠道则是直接显示的图片，所以需要设置 `external storage` 才能进行显示，设置后图片变为一个指向 `external_image_storage` 中图片地址的链接，如果本机无法访问 `external_image_storage` 图片是无法显示的，所以如果设置成 local 的时候邮件中的图片可能也不能显示了，因为这个时候是直接一张图片链接。不过 `external_image_storage` 可设置的 provider 包括 s3、webdav、gcs、azure_blob、local，如果是线上服务的话建议接入对象存储，比如 s3 服务，这里我们使用阿里云 OSS 来配置 s3 进行说明。

在 Grafana 配置文件 `/etc/grafana/grafana.ini` 中配置 `external_image_storage`：

```ini
#################################### External image storage ##########################
[external_image_storage]
provider = s3  # 使用 s3 模式

[external_image_storage.s3]
endpoint = oss-cn-beijing.aliyuncs.com
bucket =
region = oss-cn-beijing
access_key =   # 使用阿里云后台的 ak 和 sk 进行配置
secret_key =

#################################### Server ####################################
[server]
domain = 192.168.31.46  # 设置 Grafana 访问地址为内网 IP
```

另外注意需要将 Grafana 的访问域名设置成内网 IP，否则在局域网其他节点上访问不到，配置完成后重启 Grafana 即可：

`systemctl daemon-reload ☸ ➜ systemctl restart grafana-server`

配置完成后我们重新去触发下报警，正常在邮件和钉钉中收到的图片都可以正常显示了：
