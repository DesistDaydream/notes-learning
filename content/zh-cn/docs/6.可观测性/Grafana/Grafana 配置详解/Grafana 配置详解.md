---
title: Grafana 配置详解
weight: 1
---


# 概述

> 参考：
> 
> - [官方文档，Setup-配置 Grafana](https://grafana.com/docs/grafana/latest/setup-grafana/configure-grafana/)

Grafna 可以通过 [INI](/docs/2.编程/无法分类的语言/INI.md) 格式的配置文件、命令行标志、环境变量来配置运行时行为。

环境变量 与 配置文件 中的配置具有一一对应的关系。环境变量可以覆盖配置文件中的配置(即.环境变量的优先级更高，如果有相通配置，以环境变量的配置为主)。

环境变量格式：`GF_<SectionName>_<KeyName>`

- SectionName 对应配置文件中 `[ ]` 中的内容
- KeyName 对应配置文件中的关键字。
- 配置文件中的 `.` 和 `-` 两个符号，到环境变量中则变为 `_` 符号。环境变量的文本全是大写的

假如现在的配置文件内容如下：

```git
# default section
instance_name = ${HOSTNAME}
[security]
admin_user = admin
[auth.google]
client_secret = 0ldS3cretKey
[plugin.grafana-image-renderer]
rendering_ignore_https_errors = true
```

对应环境变量，则是：

```bash
GF_DEFAULT_INSTANCE_NAME=my-instance
GF_SECURITY_ADMIN_USER=owner
GF_AUTH_GOOGLE_CLIENT_SECRET=newS3cretKey
GF_PLUGIN_GRAFANA_IMAGE_RENDERER_RENDERING_IGNORE_HTTPS_ERRORS=true
```

可以看到，Grafana 的配置具有层次感，配置文件中 `[ ]` 表示一套配置环境，配置环境下方，有具体的配置关键字。

## 备注

Grafana 容器镜像会默认配置一些环境变量，以指定一些基本的配置路径。[此处](https://grafana.com/docs/grafana/latest/administration/configure-docker/)是官方对镜像的解释

# grafana.ini 配置文件详解

> Grafana 的配置文件开头 `;` 表示注释

## [paths]

**data = \<STRING>** # Grafana 数据存储路径。`默认值：/var/lib/grafana/data`
**logs = \<STRING>** # Grafana 日志模式为 file 时，记录日志的路径。`默认值：/var/log/grafana`
**plugins = \<STRING>** # Grafana 插件的安装路径。`默认值：/var/lib/grafana/plugins`
**provisioning = \<STRING>** # Grafana 的 provisioning 功能加载配置文件的路径。`默认值：/etc/grafana/provisioning`

## [server]

**http_port = \<INT>** # Grafana 监听的端口。`默认值：3000`。
**root_url = \<STRING>** # 通过 Web 浏览器访问 Grafana 的完整 URL。`默认值：%(protocol)://%(domain)s:%(http_port)s/`

- %(protocol)、%(domain)、%(http_port) 对应配置文件中 \[server] 部分的 protocol、domain、http_port 字段
- root_url 常用于重定向和发送电子邮件时填写 URL

## [security]

**admin_user = \<STRING>** # Grafana Web UI 的管理员账号的用户命。默认值：admin
**admin_password = \<STRING>** # Grafana Web UI 的管理员账号的密码。默认值：admin

## [auth]

## [auth.anonymous]

匿名访问的配置，配置匿名访问，可以使无需任何登录操作即可访问 Grafana
**enabled = \<BOOL>** # 是否开启匿名访问，开启后，可以匿名访问指定的 Organization(组织) 的仪表盘。`默认值：false`。
**org_name = \<STRING>** # 匿名用户可以访问的组织名称。`默认值：Main Org.`。
**org_role = \<STRING>** # 通过匿名访问的组织应该具有的权限。`默认值：Viewer`。

- 可用的值有 Editor 和 Admin。

## [log]

Grafana 日志配置

**mode = \<console | file | syslog>** # Grafana 记录日志的模式，多种模式以空格分隔。`默认值： console file`

**level = \<debug | info | warn | error | critical>** # 日志级别。`默认值：info`

## [smtp]

https://grafana.com/docs/grafana/latest/administration/configuration/#smtp

Grafana 的 Email 服务端配置。通过 smtp 部分的配置，Grafana 可以通过邮件 发送告警、重置密码 等。

注意，当我们使用邮箱重置密码时，会发现 Grafana 发送的重置连接的域名是 `http://localhost:3000`

![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/shgqef/1639992484686-2a5d4b08-6d1e-412c-a557-a82ade2ea984.png)

这个域名是从过配置文件中 `[server]` 部分的 `root_rul` 字段获取的。

### 配置示例

```yaml
smtp:
  enabled: true
  host: "smtp.263.net:25"
  user: "lich_wb@ehualu.com"
  password: "邮箱密码"
  from_address: "lich_wb@ehualu.com"
  from_name: Grafana
```

```ini
[smtp]
enabled = true
from_address = lich_wb@ehualu.com
from_name = Grafana
host = smtp.263.net:25
password = 邮箱密码
user = lich_wb@ehualu.com
```

# grafana.ini 配置示例

```git
[server]
http_port = 3000
[analytics]
check_for_updates = true
[log]
mode = console
level = info
[paths]
data = /var/lib/grafana/data
logs = /var/log/grafana
plugins = /var/lib/grafana/plugins
provisioning = /etc/grafana/provisioning
```
