---
title: Server 配置
linkTitle: Server 配置
date: 2023-10-31T22:24
weight: 2
---

# 概述

> 参考：
>
> - [官方文档，配置-配置](https://prometheus.io/docs/prometheus/latest/configuration/configuration/)

Prometheus Server 可通过两种方式来改变运行时行为

- 命令行标志
- 配置文件

## 配置文件热更新

Prometheus Server 可以在运行时重新加载其配置文件(也就俗称的热更新)。如果是新的配置不正确，则 Prometheus Server 则不会应用新的配置，并在日志中显示错误点。

有两种方式可以实现 Prometheus Server 的热更新功能

- 向 Prometheus Server 进程发送 `SIGHUP` 信号。
- 向 Prometheus Server 暴露的 `/-/reload` 端点发送 `HTTP 的 POST` 请求

注意：想要实现热更新功能，需要在 Prometheus Server 中指定 `--web.enable-lifecycle` 标志，这也将重新加载所有的 Rules 配置文件。

# Prometheus Server 命令行标志详解

可以通过 prometheus -h 命令查看所有的可以用标志

prometheus 程序在启动时，可以使用一些标志来对程序进行一些基本设定，比如数据存储路径、存储时间等等

- **--config.file=/PATH/TO/FILE** # prometheus 主配置文件，默认为当前路径的 prometheus.yml
- **--enable-feature=...** # 启动指定的功能特性，多个功能以逗号分割。可以开启的功能详见：[官方文档，已关闭的功能](https://prometheus.io/docs/prometheus/latest/disabled_features/)
- **--web.listen-address="0.0.0.0:9090"** # Prometheus 监听地址。`默认值：0.0.0.0:9090`。该端口用于 Web UI、API 和 Telemetry(遥测)
- **--web.config.file=/PATH/TO/FILE** # \[实验标志]用于开启 TLS 或 身份验证 配置文件路径。
- --web.read-timeout=5m # Maximum duration before timing out read of the request, and closing idle connections.
- **--web.max-connections=INT** # 可以同时连接到 Prometheus Server 的最大数量。`默认值:512`
- **--web.external-url=URL**# 可以从外部访问 Prometheus 的 URL。
  - 例如，如果 Prometheus 是通过反向代理提供的，用于生成返回 Prometheus 本身的相对和绝对链接。如果 URL 具有路径部分，它将被用作所有 HTTP 的前缀 Prometheus 服务的端点。 如果省略，则会自动派生相关的 URL 组件。
    - 注意：该标志在反向代理时似乎问题，详见：<https://github.com/prometheus/prometheus/issues/1583>
  - 例如，Prometheus 产生的的告警，推送到 AlertManager 时，会有一个 `generatorURL` 字段，该字段中所使用的 URL 中的 Endpoint，就是 web.external-url，这个 URL 可以让获取该告警的人，点击 URL 即可跳转到 Prometheus 的 Web 页面并使用对应的 PromQL 查询。
- **--web.route-prefix=PATH** # Web 端内部路由的前缀。 默认为 --web.external-url 标志指定的路径。i.e.后端代码的路由入口路径。一般默认为 / 。
- --web.user-assets= # Path to stat storage.tsdb.max-block-durationic asset directory, available at /user.
- **--web.enable-lifecycle** # 开启配置热更新，开启后，可使用 curl -X POST <http://PrometheusServerIP:9090/-/reload> 命令来重载配置以便让更改后的配置生效，而不用重启 prometheus 进程
- **--web.enable-admin-api** # 开启管理操作 API 端点。通过 admin API，可以删除时序数据。
- --web.console.templates="consoles" # Path to the console template directory, available at /consoles.
- --web.console.libraries="console_libraries" # Path to the console library directory.
- --web.page-title="Prometheus Time Series Collection and Processing Server" # Document title of Prometheus instance.
- --web.cors.origin=".\*" # Regex for CORS origin. It is fully anchored. Example: 'https?://(domain1|domain2).com'
- --web.enable-remote-write-receiver # 开启 Prometheus [Storage(存储)](docs/6.可观测性/监控系统/Prometheus/Storage(存储)/Storage(存储).md) 中的 Remote Storage(远程存储) 功能。
- **--storage.tsdb.path="/PATH/DIR"**# prometheus 存储 metircs 数据的目录(使用绝对路径)
- **--storage.tsdb.retention.time=TIME** # 数据的存储时间，如果既未设置此标志也未设置 storage.tsdb.retention.size 标志，`默认值：15d`。支持的单位：y，w，d，h，m，s，ms。
- --storage.tsdb.retention.size=STORAGE.TSDB.RETENTION.SIZE # [EXPERIMENTAL] Maximum number of bytes that can be stored for blocks. Units supported: KB, MB, GB, TB, PB. This flag is experimental and can be changed in future releases.
- --storage.tsdb.no-lockfile # 不在数据目录创建锁文件。暂时不理解什么意思，待研究
- --storage.tsdb.allow-overlapping-blocks # \[EXPERIMENTAL] Allow overlapping blocks, which in turn enables vertical compaction and vertical query merge.
- --storage.tsdb.wal-compression # Compress the tsdb WAL.
- --storage.remote.flush-deadline= # How long to wait flushing sample on shutdown or config reload.
- --storage.remote.read-sample-limit=5e7 # Maximum overall number of samples to return via the remote read interface, in a single query. 0 means no limit. This limit is ignored for streamed response types.
- --storage.remote.read-concurrent-limit=10 # Maximum number of concurrent remote read calls. 0 means no limit.
- --storage.remote.read-max-bytes-in-frame=1048576 # Maximum number of bytes in a single frame for streaming remote read response types before marshalling. Note that client might have limit on frame size as well. 1MB as recommended by protobuf
- by default.
- --rules.alert.for-outage-tolerance=1h # Max time to tolerate prometheus outage for restoring "for" state of alert.
- --rules.alert.for-grace-period=10m # Minimum duration between alert and restored "for" state. This is maintained only for alerts with configured "for" time greater than grace period.
- **--rules.alert.resend-delay=DURATION**# 向 Alertmanager 重新发送警报前的最少等待时间。`默认值：1m`。
  - 当告警处于 FIRING 状态时，每间隔 1m，就会再次发送一次。注意：重发送之前，还需要一个评估规则的等待期，评估完成后，再等待该值的时间，才会重新发送告警。
- --alertmanager.notification-queue-capacity=10000 # The capacity of the queue for pending Alertmanager notifications.
- --alertmanager.timeout=10s # Timeout for sending alerts to Alertmanager.
- **--query.lookback-delta=DURATION** # 评估 PromQL 表达式时最大的回溯时间。`默认值：5m`
  - 比如，当采集目标的间隔时间为 10m 时，由于该设置，最大只能查询当前时间的前 5m 的数据，这是，即时向量表达式返回的结果将会为空。
- **--query.timeout=DURATION** # 一次查询的超时时间。`默认值：2m`
- --query.max-concurrency=20 # Maximum number of queries executed concurrently.
- --query.max-samples=50000000 # Maximum number of samples a single query can load into memory. Note that queries will fail if they try to load more samples than this into memory, so this also limits the number of samples a query can return.
- **--log.level=STRING** # 设定 Prometheus Server 运行时输出的日志的级别。`默认值：info`。 可用的值有：debug, info, warn, error
- **--log.format=logfmt** # 设定 Prometheus Server 运行时输出的日志的格式。`默认值：logfmt`。可用的值有：logfmt, json

# prometheus.yaml 配置文件详解

下文用到的字段值的占位符说明

- \<BOOLEAN> # 可以采用 true 或 false 值的布尔值
- \<DURATION> # 持续时间。可以使用正则表达式
  - (((\[0-9]+)y)?((\[0-9]+)w)?((\[0-9]+)d)?((\[0-9]+)h)?(((\[0-9]+)m)?(((\[0-9]+)s)?(((\[0-9]+)ms)?|0)，例如：1d、1h30m、5m、10s。
- \<FILENAME> # 当前工作目录中的有效路径
- \<HOST> # 由主机名或 IP 后跟可选端口号组成的有效字符串。
- \<INT> # 一个整数值
- \<LABELNAME> # 与正则表达式\[a-zA-Z \_] \[a-zA-Z0-9 \_] \*匹配的字符串
- \<LABELVALUE> # 一串 unicode 字符
- \<PATH> # 有效的 URL 路径
- \<SCHEME> # 一个字符串，可以使用值 http 或 https
- \<SECRET> # 作为机密的常规字符串，例如密码
- \<STRING> # 常规字符串
- \<TMPL_STRING> # 使用前已模板扩展的字符串

每个字段下

## 顶层字段

- **global**([global](#global)) # 全局配置，所有内容作用于所有配置环境中,若其余配置环境中不再指定同样的配置，则global中的配置作为默认配置
- **rule_files**([rule_files](#rule_files)) #
- **scrape_configs**(\[][scrape_configs](#scrape_configs(占比最大的字段))) # 抓取 Target 的 metrics 时的配置
- **alerting**([alerting](#alerting)) # 与 Alertmanager 相关的配置
  - alert_relabel_configs([relabel_configs](#relabel_configs))
  - alertmanagers
- **remote_write**(\[][remote_write](#remote_write)) # 与远程写入相关功能的配置
- **remote_read**(\[][remote_read](#remote_read)) # 与远程读取相关功能的配置

## global

全局配置，所有内容作用于所有配置环境中,若其余配置环境中不再指定同样的配置，则 global 中的配置作为默认配置

**scrape_interval(DURATION)** # 抓取 targets 的指标频率，`默认值：1m`。

**scrape_timeout(DURATION)** # 对 targets 发起抓取请求的超时时间。`默认值：10s`。

**evaluation_interval(DURATION)** # 评估规则的周期。`默认值：1m`。

该字段主要用于向规则配置文件传递全局的配置。这个值会被规则配置文件中的 `.groups.interval` 覆盖，详见 interval 字段详解

**external_labels(map\[STRING]STRING)** # 与外部系统(federation, remote storage, Alertmanager)通信时添加到任何时间序列或警报的标签。

- **KEY: VAL** # 比如该键值可以是 run: httpd，标签名是 run，run 的值是 httpd，KEY 与 VAL 使用字母，数字，\_，-，.这几个字符且以字母或数字开头；val 可以为空。
- ......

## rule_files

[规则文件配置](/docs/6.可观测性/监控系统/Prometheus/Rules%20配置.md)列表，从所有匹配到的文件中读取配置内容。可以使用正则表达式匹配多个符合的文件。Prometheus 支持两种规则

- recording rules(记录规则)
- alerting rules(告警规则)

## scrape_configs(占比最大的字段)

https://prometheus.io/docs/prometheus/latest/configuration/configuration/#scrape_config

在 [Prometheus](/docs/6.可观测性/监控系统/Prometheus/Prometheus.md) 一文中，粗略介绍了基本的 scrape_configs 配置段的内容，下面是最基本的配置：

```yaml
scrape_configs:
  - job_name: "prometheus"
    scrape_interval: 5s
    static_configs:
      - targets: ["localhost:9090"]
```

scrape_configs 是 Prometheus 采集指标的最重要也是最基本的配置信息，scrape_configs 字段是一个数组，所以可以配置多个 Scrape 配置，不同的 Scrape 配置，所以该段配置至少需要包含以下几个方面：

- **名字** # 每个 scrape 工作都应该具有一个名字。称为 job_name，名字主要起到标识符的作用。
  - 该示例定义了一个抓取配置的 job，名字叫 prometheus
- **目标** # 要抓取的 metrics 的目标。目标可以通过 **静态** 或者 **动态(i.e.各种服务发现)** 这两种方式指定
  - 该示例通过静态配置定义这个 job 中要抓取的目标主机，目标主机由 IP:PORT 组成
- **间隔** # 该 scrape 工作每次抓取 metrics 的时间间隔。就是每隔 X 秒抓一次
  - 该示例每次抓取 metrics 的时间间隔为 5 秒(i.e.每 5 秒获取一次 metrics)
- **其他** # 除了名字、目标、间隔 以外，还可以配置一些额外的抓取配置，比如发起 HTTP 请求时需要携带的 Header 与 Body、抓取策略 等等

### 基本配置

**job_name(STRING)** # 指定抓取 Metrics 的 Job 名字

**scrape_interval(DURATION)** # 指定这个 job 中抓取 targets 的频率。默认使用 global 配置环境中同名参数的值

**scrape_timeout(DURATION)** # 指定这个 job 中抓取 targets 的超时时长。默认使用 global 配置环境中同名参数的值

**metrics_path: PATH** # 从 targets 获取 metrics 时 http 请求的路径。默认为/metrics

**honor_labels(BOOLEAN)** # 控制 Prometheus 如何处理标间之间的冲突。`默认值：false`

- 获取 targets 的 metrics 时(e.g.snmp_exporter|Federate|pushgateway 等)，其中的标签有可能会与本身的标签存在冲突
  - 该参数的值为 true 时，则以抓取数据中的标签为准
  - 值为 false 时，就会重新命名表桥为 exported 形式，然后添加配置文件中的标签。

**honor_timestamps(BOOLEAN)** # 控制 Prometheus 是否尊重抓去到的数据中的时间戳 `默认值：true`

- 比如从 federate、pushgateway 等地方获取指标时，指标中都是带着时间戳的，
  - 若设置为 false，则会忽略这些采集到的时间戳，在入库时加上采集时的时间戳。
  - 若设置为 true，则是在入库时使用抓到到的指标中的时间戳。

**sample_limit(INT)** # 每次抓取 metrics 的数量限制。`默认值：0`。0 表示不限制

### HTTP 配置

Prometheus 抓取目标就是发起 HTTP 请求。

除了 scheme、params 字段以外的其他字段是 Prometheus 共享库中的通用 HTTP 客户端配置，即下面的 `HTTPClientConfig` 结构体中的内容。

代码：[common/config/http_config.go](https://github.com/prometheus/common/blob/v0.30.0/config/http_config.go#L159)

```go
// HTTPClientConfig configures an HTTP client.
type HTTPClientConfig struct {
 // The HTTP basic authentication credentials for the targets.
 BasicAuth *BasicAuth `yaml:"basic_auth,omitempty" json:"basic_auth,omitempty"`
 // The HTTP authorization credentials for the targets.
 Authorization *Authorization `yaml:"authorization,omitempty" json:"authorization,omitempty"`
 // The OAuth2 client credentials used to fetch a token for the targets.
 OAuth2 *OAuth2 `yaml:"oauth2,omitempty" json:"oauth2,omitempty"`
 // The bearer token for the targets. Deprecated in favour of
 // Authorization.Credentials.
 BearerToken Secret `yaml:"bearer_token,omitempty" json:"bearer_token,omitempty"`
 // The bearer token file for the targets. Deprecated in favour of
 // Authorization.CredentialsFile.
 BearerTokenFile string `yaml:"bearer_token_file,omitempty" json:"bearer_token_file,omitempty"`
 // HTTP proxy server to use to connect to the targets.
 ProxyURL URL `yaml:"proxy_url,omitempty" json:"proxy_url,omitempty"`
 // TLSConfig to use to connect to the targets.
 TLSConfig TLSConfig `yaml:"tls_config,omitempty" json:"tls_config,omitempty"`
 // FollowRedirects specifies whether the client should follow HTTP 3xx redirects.
 // The omitempty flag is not set, because it would be hidden from the
 // marshalled configuration when set to false.
 FollowRedirects bool `yaml:"follow_redirects" json:"follow_redirects"`
}
```

**scheme(STRING)** # 指定用于抓取 Metrics 时使用的协议。`默认值：http`

**params: <>** # 发起 http 请求时，URL 里的参数(以键值对的方式表示)。
常用于 snmp_exporter，比如 <http://10.10.100.12:9116/snmp?module=if_mib&target=10.10.100.254>，问号后面就是参数的 key 与 value)

- STRING: STRING

**basic_auth(Object)**# 配置 HTTP 的基础认证信息。

- **username(STRING)** #
- **password(SECRET)** #
- **password_file(STRING)** #

**authorization(Object)** #

- **type(STRING)** # 发起抓取请求时的身份验证类型。`默认值：Bearer`
- **credentials(SECRET)** # 用于身份验证的信息。与 credentials_file 字段互斥。如果是 type 字段是 Bearer，那么这里的值就用 Token 即可。
- **credentials_file(FileName)** # 从文件中读取用于身份验证的信息。与 credentials 字段互斥

**oauth2(Object)** # 配置 OAuth 2.0 的认证配置。与 basic_auth 和 authorization 两个字段互斥

**proxy_url(STRING)** # 指定代理的 URL

**tls_config**([tls_config](#tls_config)) # 指定抓取 metrics 请求时的 TLS 设定

### Scrape 目标配置

Prometheus 将会根据这里的字段配置，以发现需要 Scrape 指标的目标，有两种方式来发现目标：静态 与 动态。

**static_configs**([static_configs](#static_configs)) # 静态配置。直接指定需要抓去 Metrics 的 Targets。

- 具体配置详见下文[静态目标发现](#J021o)

**XX_sd_configs**(\[]OBJECT) # 动态配置。动态需要抓去 Metrics 的 Targets。XXX_sd_configs 中的 sd 全称为 Service Discovery(服务发现)

- 具体配置详见下文[动态目标发现](#IWvg5)
- 不同的服务发现，有不同的配置方式。比如 `kubernetes_sd_configs`、`file_sd_configs` 等等。
- 注意：当 Prometheus 自动发现这些待抓取目标时，会附带一些原始标签，这些标签以 `__meta_XX` 开头，不同的服务发现配置发现标签不同，具体说明详见[《Label 与 Relabel》文章中的 Discovered Labels 章节](/docs/6.可观测性/监控系统/Prometheus/Target(目标)%20与%20Relabeling(重新标记).md) 的说明

`XX_sd_configs` 与 `static_configs` 的区别：静态配置与动态配置就好比主机获取 IP 时是 DHCP 还是 STATIC。动态配置可以动态获取要抓取的 Targets、静态就是指定哪个 Target 就去哪个 Target 抓取 Metrics

### Relabel 配置

**relabel_configs**([relabel_configs](#relabel_configs)) # 在发现目标后，重新配置 targets 的标签。

- 具体配置详见下文 [重设标签](#重设标签)

**metric_relabel_configs**([relabel_configs](#relabel_configs)) # 在抓取到指标后，重新配置 metrics 的标签

- 与 relabel_configs 字段配置内容相同

## alerting

**alert_relabel_configs**([relabel_configs](#relabel_configs))

适用于推送告警时的 Relabel 功能，配置与 [relabel_configs](#PGKul) 相同

**alertmanager**(\[]OBJECT)

https://prometheus.io/docs/prometheus/latest/configuration/configuration/#alertmanager_config

> 该字段配置方式与 scrape_config 字段的配置非常相似，只不过不是配置抓取目标，而是配置推送告警的目标

alertmanager 字段指定了 Prometheus Server 发送警报的目标 Alertmanager，还提供了参数来配置如何与这些 Alertmanager 通信。此外，relabel_configs 允许从已发现的实体中选择 Alertmanagers，并对使用的 API 路径进行高级修改，该路径通过 **alerts_path** 标签暴露。

**alerts.timeout(DURATION)** # 推送警报时，每个目标 Alertmanager 超时，单位：秒。`默认值: 10`。

**timeout(DURATION)** # 推送告警时的超时时间。

**api_version(STRING)** # 推送告警时，应该使用哪个版本的 Alertmanager 路径。`默认值：v2`。

**path_prefix(PATH)** # 推送告警时的，目标路径前缀。`默认值：/`。

- 注意：就算指定了其他路径，也会默认在末尾添加 `/api/v2/alerts`

#### HTTP 配置

**scheme(SCHEME)** # 推送告警时，所使用的协议。`默认值：HTTP`

下面的部分是 HTTP 的认证，是用来配置将告警推送到目标时所需要的认证信息。比如目标是 HTTPS 时，就需要这些配置。发起的 POST 推送请求时，Prometheus 使用 username 和 passwrod 字段的值为这个 HTTP 请求设置 Authorization 请求头。说白了就是发起 HTTP 请求时带着用户名和密码。

**basic_auth(Object)**

- **username(STRING)** #
- **password(SECRET)** # password 和 password_files 字段是互斥的
- **password_file(STRING)** #

**authorization(Object)** #

- **type(STRING)** # 推送告警时的身份验证类型。`默认值：Bearer`
- **credentials(secret)** # 用于身份验证的信息。与 credentials_file 字段互斥。如果是 type 字段是 Bearer，那么这里的值就用 Token 即可。
- **credentials_file(filename)** # 从文件中读取用于身份验证的信息。与 credentials 字段互斥

**oauth2(Object)** # 配置 OAuth 2.0 的认证配置。与 basic_auth 和 authorization 两个字段互斥

**tls_config(Object)** # 指定推送告警时的 TLS 设定

#### Alerts 推送目标的配置

Prometheus 根据这部分配置来推送需要

**static_configs**([static_configs](#static_configs)) # 静态配置。指定推送告警时的目标。

- 具体配置详见下文 [静态目标发现](#静态目标发现)

**XXX_sd_configs**([]OBJECT) # 动态配置。动态发现可供推送告警的 alertmanager- XXXX # 不同的服务发现，有不同的配置方式。与 scrape_configs 字段中的 XXX_sd_configs 配置类似。

- 具体配置详见下文 [动态目标发现](#IWvg5)

#### Relabel 配置

**relabel_configs**([relabel_configs](#relabel_configs)) # 在发现目标后，重新配置 targets 的标签

详见下文 [重设标签](#重设标签)

## remote_write

与远程写相关的配置，详见 [Prometheus 存储章节](/docs/6.可观测性/监控系统/Prometheus/Storage(存储)/Storage(存储).md)

**url(STRING)** # 指定要发送时间序列数据到远程存储的端点的 URL

## remote_read

与远程读相关的配置，详见 [Prometheus 存储章节](/docs/6.可观测性/监控系统/Prometheus/Storage(存储)/Storage(存储).md)

**url(STRING)** # 指定发起查询请求的远程数据库的端点的 URL

# 配置文件中的通用配置字段

## 静态目标发现

这些通用字段会被配置文件中的某些字段共同使用

### static_configs

https://prometheus.io/docs/prometheus/latest/configuration/configuration/#static_config

静态配置。指定用户抓取 metrics 的 targets。静态配置与动态配置就好比主机获取 IP 时是 DHCP 还是 STATIC。动态配置可以动态获取要抓取的 targets、静态就是指定哪个 target 就去哪个 target 抓取 metrics

**targets([]STRING)** # 指定要抓取 metrics 的 targets 的 IP:PORT

- **HOST**

**labels(map\[STRING]STRING)** # 指定该 targets 的标签，可以随意添加任意多个

- **KEY: VAL** # 比如该键值可以是 run: httpd，标签名是 run，run 的值是 httpd，key 与 val 使用字母，数字，\_，-，.这几个字符且以字母或数字开头；val 可以为空。
- ......

## 动态目标发现

### file_sd_configs

https://prometheus.io/docs/prometheus/latest/configuration/configuration/#file_sd_config

基于文件的服务发现提供了一种配置静态目标的更通用的方法，并用作插入自定义服务发现机制的接口。

在 Prometheus 支持的众多服务发现的实现方式中，基于文件的服务发现是最通用的方式。这种方式不需要依赖于任何的平台或者第三方服务。对于 Prometheus 而言也不可能支持所有的平台或者环境。通过基于文件的服务发现方式下，Prometheus 会定时从指定文件中读取最新的 Target 信息，因此，你可以通过任意的方式将监控 Target 的信息写入即可。

用户可以通过 JSON 或者 YAML 格式的文件，定义所有的监控目标。同时还可以通过为这些实例添加一些额外的标签信息，这样从这些实例中采集到的样本信息将包含这些标签信息，从而可以为后续按照环境进行监控数据的聚合。

**files(map\[STRING]STRING)** # Prometheus 将要读取的文件路径，将会从该文件发现待采集的 Target。支持正则表达式

**refresh_interval(DURATION)** # 重新读取文件的间隔时间。`默认值：5m`

通过这种方式，Prometheus 会自动的周期性读取文件中的内容。当文件中定义的内容发生变化时，不需要对 Prometheus Server 进行任何的重启操作。

#### 配置样例

假设现在有一个名为 file_sd.yaml 文件，中分别定义了 2 个采集任务，以及每个任务对应的 Target 列表，内容如下：

```yaml
- targets:
    - "172.19.42.200"
  labels:
    network: "switch"
- targets:
    - "172.19.42.243"
  labels:
    server: "host-1"
```

创建 Prometheus 配置文件/etc/prometheus/prometheus-file-sd.yml，并添加以下内容：

```yaml
global:
  scrape_interval: 15s
  scrape_timeout: 10s
  evaluation_interval: 15s
scrape_configs:
  - job_name: "file_ds"
    file_sd_configs:
      - refresh_interval: 5m # Prometheus 默认每 5m 重新读取一次文件内容，当需要修改时，可以通过refresh_interval进行设置
        files:
          - "file_sd.yaml"
```

这里定义了一个基于 file_sd_configs 的监控采集任务，其中模式的任务名称为 file_ds。在 JSON 文件中可以使用 job 标签覆盖默认的 Job 名称。

在 Prometheus UI 的 Targets 下就可以看到当前从 targets.json 文件中动态获取到的 Target 实例信息以及监控任务的采集状态，同时在 Labels 列下会包含用户添加的自定义标签

这种通用的方式可以衍生了很多不同的玩法，比如与自动化配置管理工具(Ansible)结合、与 Cron Job 结合等等。 对于一些 Prometheus 还不支持的云环境，比如国内的阿里云、腾讯云等也可以使用这种方式通过一些自定义程序与平台进行交互自动生成监控 Target 文件，从而实现对这些云环境中基础设施的自动化监控支持。

### kubernetes_sd_configs

https://prometheus.io/docs/prometheus/latest/configuration/configuration/#kubernetes_sd_config

kubernetes_sd_configs 字段的配置可以让 Prometheus 从 Kubernetes 的 API Server 中自动发现需要抓取目标，并始终与集群状态保持同步。可以抓取的目标有 node、service、pod、endpoints、ingress。

> 注意：如果 Prometheus Server 部署在 Kubernetes 集群外部，通过 k8s 的 API Server 自动发现的 pod ip 是集群内部 IP，一般情况下不互联的。因为 pod 的 ip 一般都是集群内部 IP。所以如果在发现目标后想要采集，需要在 Prometheus Server 所在服务器添加到 Kubernetes 的 Pod IP 的路由条目。

Note：使用该配置进行服务发现，请求都会经过 API Server，集群规模越大，API Server 压力也会跟随增高。

#### API Server 配置

**api_server(STRING)** # 指定 k8s 集群中 API Server 的地址。

- 如果该字段为空，则默认 Prometheus 在 k8s 集群内部运行，将自动发现 apiserver，并使用 Pod 中 /var/run/secrets/kubernetes.io/serviceaccount/ 目录下的的 CA 证书 和 Token。

**basic_auth(Object)**# 如果 apiserver 使用基本认证启动，则使用 basic_auth 字段。`authorization` 字段互斥。password 和 password_file 是互斥的。

- **username(STRING)** #
- **password(SECRET)** #
- **password_file(STRING)** #

**authorization(Object)** # 如果 apiserver 使用证书启动，则使用 authorization 字段。与 `basic_auth` 字段互斥。

- **type(STRING)** # 发起抓取请求时的身份验证类型。`默认值：Bearer`
- **credentials(SECRET)** # 用于身份验证的信息。与 credentials_file 字段互斥。如果是 type 字段是 Bearer，那么这里的值就用 Token 即可。该字段就是老版本的 bearer_token 字段
- **credentials_file(filename)** # 从文件中读取用于身份验证的信息。与 credentials 字段互斥.该字段就是老版本的 bearer_token_file 字段

**oauth2(Object)** # 配置 OAuth 2.0 的认证配置。与 basic_auth 和 authorization 两个字段互斥

**tls_config(Object)** # 指定抓取 metrics 请求时的 TLS 设定

**proxy_url(STRING)** # Optional proxy URL

#### 目标发现的规则配置

**role(STRING)** # 根据 STRING 动态发现地 Target。可用的 STRING 为 endpoints, service, pod, node,ingress。

- 比如，Prometheus 可以自动发现 ep、svc 等等对象作为 scrape 地 target

**namespaces(Object)** # 指定动态发现哪个 namesapce 下的 Target ，如果省略，则 Target 将从所有 namespaces 中动态发现

- **names([]STRING)**

**selectors([]Object)** # 可以根据 selectors 中指定地 label 或者 field 来过滤动态发现的 Target 。如果省略，则不进行任何过滤。

- **role(ROLE)** #
- **label(STRING)** # STRING 使用 key=value 的格式。
- **field(STRING)** #

#### 配置样例

> 参考：
>
> - [官方推荐的样例](https://github.com/prometheus/prometheus/blob/main/documentation/examples/prometheus-kubernetes.yml)

下面的例子是这样的：动态发现 kube-system 名称空间下的所有 pod 作为 target，并且进行过滤，只选择其中标签为 k8s-app=kube-dns 的 pod 作为 target

```yaml
scrape_configs:
  - job_name: "kubernetes-node"
    honor_timestamps: true
    metrics_path: /metrics
    # 注意 scheme 字段，自动发现机制只是会发现 IP:PORT，并不会添加协议，有的 pod 是只提供 https 的。
    # 比如，如果是发现 kubelet、kube-apiserver 等 pod ，则这里应该改为 https
    scheme: http
    bearer_token_file: /etc/prometheus/config_out/serviceaccount/token
    tls_config:
      insecure_skip_verify: true
      ca_file: /etc/prometheus/config_out/serviceaccount/ca.crt
    kubernetes_sd_configs:
      - api_server: "https://172.19.42.234:6443"
        # 这里写了两遍认证信息，这是因为这里的认证则是针对 apiserver 的认证。而 scrape_configs 字段下的认证是针对已经发现的目标进行认证。
        # 需要先认证 apiserver 以发现待抓取的目标，然后再使用 scrape_configs 字段下的认证来采集目标的指标。
        bearer_token_file: /etc/prometheus/config_out/serviceaccount/token
        tls_config:
          insecure_skip_verify: true
          ca_file: /etc/prometheus/config_out/serviceaccount/ca.crt
        role: pod
        namespaces:
          names:
            - kube-system
        selectors:
          - role: pod
            label: k8s-app=kube-dns
```

上面的配置将会自动发现 k8s 集群中的所有 coredns pod

> 注意：这里可以发现，我们是可以访问集群内部的 10.244.0.243，这是因为我加了静态路由配置(`ip route add 10.244.0.0/16 dev ens3 via 172.19.42.231`)，否则，集群外部的 Prometheus 是无法抓取访问不到的目标的。

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/hzhbid/1616049623195-79b06041-01c2-4b81-bcb5-d4efd06de281.png)

可以看到，coredns 的两个端口都发现了，由于我们不需要 53 端口，所以还需要进一步过滤，就是把 53 端口过滤调。可以使用 Relabeling 功能，在配置后面添加如下内容：

```yaml
relabel_configs:
  - source_labels: [__meta_kubernetes_pod_container_port_number]
    regex: 53
    action: drop
```

此时，我们删除了 `__meta_kubernetes_pod_container_port_number` 这个标签的值为 53 的所有指标。这样我们就可以看到，只剩下 9153 端口的指标了

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/hzhbid/1616049623219-a5447656-6c61-40f1-acfe-df6218904b3a.png)

## 重设标签

### relabel_configs

> 参考：
>
> - [官方文档，配置-配置-relabel_config](https://prometheus.io/docs/prometheus/latest/configuration/configuration/#relabel_config)
> - [Label 与 Relabeling](/docs/6.可观测性/监控系统/Prometheus/Target(目标)%20与%20Relabeling(重新标记).md)

relabel 重设标签功能，用于将抓取到的样本中的原始 label 进行重新标记以生成新的 label。

**source_labels(LabelName, ... )** # 从现有的标签中选择将要获取值的标签作为 source_labels。source_labels 可以有多个。

**separator(STRING)** # 指定 source_labels 中所有值之间的分隔符。`默认值: ;`。

**target_label(STRING)** # 通过 regex 字段匹配到的值写入的指定的 target_label 中

**regex(REGEX)** # 从 source_label 获取的值进行正则匹配，匹配到的值写入到 target_label 中。`默认正则表达式为(.*)`。i.e.匹配所有值

**modulus(UINT64)** # 去 source_labels 值的哈希值的模数

**replacement(STRING)** # 替换。指定要写入 target_label 的值，STRING 中可以引用 regex 字段的值，使用正则表达式方式引用。`默认值：$1`。与 action 字段的 replace 值配合使用。

**action(Relabel_Action)** # 对匹配到的标签要执行的动作。`默认值: replace`。

## 其他

### tls_config

https://prometheus.io/docs/prometheus/latest/configuration/configuration/#tls_config

tls_config 字段用来配置 TLS 连接信息。下面描述客户端就是 Prometheus Server，服务端就是要抓取 Metrics 的目标。

**ca_file(FileName)** # CA 证书，用于验证服务端证书

**cert_file(FileName)** # 证书文件，用于客户端对服务器的证书认证。

**key_file(FileName)** # 密钥文件，用于客户端对服务器的证书认证。

**server_name(STRING>** # ServerName 扩展名，用于指示服务器的名称。ServerName extension to indicate the name of the server. ServerName 概念参考：<https://tools.ietf.org/html/rfc4366#section-3.1)

**insecure_skip_verify(BOOLEAN)** # 禁用服务端对证书的验证。类似于 curl 的 -k 选项
