---
title: Promtail 配置详解
---

# 概述

> 参考：
>
> - [Loki 官方文档，客户端-Promtail-配置](https://grafana.com/docs/loki/latest/clients/promtail/configuration/)
> - [GitHub 官方文档](https://github.com/grafana/loki/blob/master/docs/sources/clients/promtail/configuration.md)

# promtail.yaml 配置文件详解

Promtail 在 YAML 文件（通常称为 config.yaml）中进行配置，该文件包含 Promtail 运行时信息，抓取到的日志存储位置，以及抓取日志的行为

顶级字段

- **server**([server](#server)) #  配置 promtail 程序运行时行为。如指定监听的ip、port等信息。
- **clients**([clients](#clients)) # 配置 Promtail 如何连接到 Loki 的多个实例，并向每个实例发送日志。
  - Note：如果其中一台远程Loki服务器无法响应或发生任何可重试的错误，这将影响将日志发送到任何其他已配置的远程Loki服务器。
  - 发送是在单个线程上完成的！ 如果要发送到多个远程Loki实例，通常建议并行运行多个Promtail客户端。
- **positions**([positions](#positions)) # positions 文件用于记录 Promtail 发现的目标。该字段用于定义如何保存 postitions.yaml 文件。
  - Promtail 发现的目标就是指日志文件。
- **scrape_configs** # 配置 Promtail 如何发现日志文件，以及如何从这些日志文件抓取日志。


# 配置如何 tail 目标
target_config: <target_config>

## server

## clients

## positions

positions 文件用于记录 Promtail 发现的目标。该字段用于定义如何保存 postitions.yaml 文件。Promtail 发现的目标就是指日志文件。

**filename(STRING)** # 指定 positions 文件路径。`默认值：/var/log/positions.yaml`

**sync_period(DURATION)**# 更新 positions 文件的时间间隔。`默认值：10s`

**ignore_invalid_yaml(BOOLEAN)** # Whether to ignore & later overwrite positions files that are corrupted。`默认值：false`

## scrape_configs(\[]OBJECT)(占比最大的字段)

> 参考：
>
> - https://grafana.com/docs/loki/latest/clients/promtail/configuration/#scrape_configs
> - [Scraping 功能官方文档](https://grafana.com/docs/loki/latest/clients/promtail/scraping)

Promtail 根据 scrape_configs 字段的内容，使用指定的发现方法从一系列目标中抓取日志。

### 基本配置

**job_name(STRING)** # 指定抓取日志的 Job 名字

**pipeline_stages(pipeline_stages)** # 定义从指定的目标抓取日志的行为。`默认值：docker{}`。详见：[Pipeline](docs/6.可观测性/Logs/Loki/Promtail/Pipeline/Pipeline.md) 与 [Stages(阶段) 详解](docs/6.可观测性/Logs/Loki/Promtail/Pipeline/Stages(阶段)%20详解.md)

**loki_push_api(loki_push_api_config)** # 定义日志推送的路径 (e.g. from other Promtails or the Docker Logging Driver)

### Scrape 目标配置

Promtail 会根据这里的字段的配置，以发现需要 Scrape 日志的目标，有两种方式来发现目标：**静态** 与 **动态**

**static_configs(\[]Object)** # 静态配置。直接指定需要抓去 Metrics 的 Targets。

- 具体配置详见下文[静态目标发现](#静态目标发现)

**XX_sd_configs**(OBJECT) # 动态配置

- 具体配置详见下文[动态目标发现](#动态目标发现)

**journal**([journal](#journal)) # 动态配置

- 具体配置详见下文[动态目标发现](#动态目标发现)

**syslog**([syslog](#syslog)) # 动态配置

- 具体配置详见下文[动态目标发现](#动态目标发现)

### Relabel 配置

**relabel_configs([]OBJECT)** # 为本 Job 下抓取日志的过程定义 Relabeling 行为。与 Prometheus 的 Relabeling 行为一致

- 具体配置详见下文[重设标签](#重设标签)

# 配置文件中的通用配置字段

## 静态目标发现

### static_configs([]Object)

**targets([]STRING)** # 指定要抓取 metrics 的 targets 的 IP:PORT

- **HOST**

**labels(map\[STRING]STRING)** # 指定该 targets 的标签，可以随意添加任意多个。

这个字段与 Prometheus 的配置有一点区别。Promtail 中必须要添加 `__path__` 这个键，以指定要抓去日志的文件路径。

- **KEY: VAL** # 比如该键值可以是 run: httpd，标签名是 run，run 的值是 httpd，key 与 val 使用字母，数字，\_，-，.这几个字符且以字母或数字开头；val 可以为空。
- ......

#### 配置示例

```yaml
- job_name: system
  pipeline_stages:
  static_configs:
    - targets: # 指定抓取目标，i.e.抓取哪台设备上的文件
        - localhost
      labels: # 指定该日志流的标签
        job: varlogs # 指定一个标签，至少需要一个非 __ 开头的标签，这样才能为日志流定义唯一标识符，否则日志流没有名字。
        __path__: /var/log/host/* # 指定抓取路径，该匹配标识抓取 /var/log/host 目录下的所有文件。注意：不包含子目录下的文件。
```

## 动态目标发现

我们可以从 grafana/loki 项目代码 [clients/pkg/promtail/scrapeconfig/scrapeconfig.go](https://github.com/grafana/loki/blob/v2.6.1/clients/pkg/promtail/scrapeconfig/scrapeconfig.go#L53) 中找到所有可以动态发现目标的配置。

### journal(OBJECT)

在具有 systemd 的 Linux 系统上，Loki 可以通过 journal 程序获取日志。

```yaml
# 从 Journal 获取的日志保留所有原始字段，并将这些信息转变为 JSON 格式。默认值：false
json: <BOOLEAN>

# 当 Promtail 启动时，从 Journal 日志文件中，获取的最老时间的日志。默认值：7h
# 比如值为7h的话,则 Promatail 于 17:00 启动，则会抓取 10:00 到 17:00 之间的日志内容
max_age: <DURATION>

# 为本次通过 Journal 日志文件采集日志的任务添加标签。
labels:
  <LabelName>: <LabelValue>
  ......

# 获取 Journal 日志文件的路径。默认值：/var/log/journal 和 /run/log/journal
path: <STRING>
```

注意：由于 Journal 程序存储日志的路径问题，所以我们我们在容器中运行 Promtail 时，必须挂载相关路径，否则 Promtail 读取不到 Journal 生成的日志。比如可以通过下面的 docker 命令运行

> /run/log/journal 一般不用挂载，大部分系统都不适用这个目录了，虽然 Journal 还是会处理该目录~

```bash
docker run \
  -v /var/log/journal/:/var/log/journal/ \
  -v /run/log/journal/:/run/log/journal/ \
  -v /etc/machine-id:/etc/machine-id \
  grafana/promtail:latest
```

下面是 journal 自动发现日志流后，自动发现的标签。

    __journal__audit_loginuid
    __journal__audit_session
    __journal__boot_id
    __journal__cap_effective
    __journal__cmdline
    __journal__comm
    __journal__exe
    __journal__gid
    __journal__hostname # 主机名
    __journal__machine_id
    __journal__pid
    __journal__selinux_context
    __journal__source_realtime_timestamp
    __journal__stream_id
    __journal__systemd_cgroup
    __journal__systemd_invocation_id
    __journal__systemd_slice
    __journal__systemd_unit
    __journal__transport
    __journal__uid
    __journal_code_file
    __journal_code_func
    __journal_code_line
    __journal_cpu_usage_nsec
    __journal_invocation_id
    __journal_message
    __journal_message_id
    __journal_priority
    __journal_priority_keyword
    __journal_syslog_facility
    __journal_syslog_identifier
    __journal_syslog_pid
    __journal_syslog_timestamp
    __journal_unit # 该标签是 unit 的名称，标签值是所有 .service，比如 ssh.service、dockerd.service 等等

#### json 字段说明

这是开启的样子：
![](https://notes-learning.oss-cn-beijing.aliyuncs.com/sxgd83/1616129621041-ee0d0d3a-b256-4a34-9b14-12bdbbc159a1.png)
这是关闭的样子：
![](https://notes-learning.oss-cn-beijing.aliyuncs.com/sxgd83/1616129621334-5638249b-63aa-446a-b276-5a621df8be5d.png)
可以看见，Parsed Fields 中，多出来很多字段。json 字段开启后，除了正常的 Message，还有很多关于该日志消息的属性。

#### 配置示例

```yaml
- job_name: systemd-journal
  journal:
    labels:
      job: systemd-journal
  relabel_configs:
    - action: labelmap
      regex: __(journal__systemd_unit)
    - action: labelmap
      regex: __(journal__hostname)
    - action: drop
      source_labels: [journal__systemd_unit]
      regex: session-.*scope
```

### kubernetes_sd_configs(\[]Object)

https://grafana.com/docs/loki/latest/clients/promtail/configuration/#kubernetes_sd_config

与 Prometheus 中的 kubernetes 的服务发现机制基本一致。与 Prometheus 配置的不同点在于，Promtail 的 kubernetes 服务发现配置一般都会使用 Relabeling 机制弄出来一个 `__path__` 标签

具体字段内容详见《[Prometheus Server 配置](/docs/6.可观测性/Metrics/Prometheus/Configuration/Promethesu%20Server.md)》文章中 [kubernetes_sd_configs](/docs/6.可观测性/Metrics/Prometheus/Configuration/Promethesu%20Server.md#kubernetes_sd_configs) 章节

#### 配置示例

```yaml
- job_name: kubernetes-pods
  pipeline_stages:
    - docker: {}
  kubernetes_sd_configs:
    - role: pod
  relabel_configs:
    # 为日志流配标签
    - source_labels: [__meta_kubernetes_namespace]
      action: replace
      target_label: namespace
      # 为日志流配置标签
    - source_labels: [__meta_kubernetes_pod_name]
      action: replace
      target_label: pod_name
      # 配置抓取日志的路径
    - source_labels:
        - __meta_kubernetes_pod_annotation_kubernetes_io_config_mirror
        - __meta_kubernetes_pod_container_name
      separator: /
      regex: (.*)
      replacement: /var/log/pods/*$1/*.log
      target_label: __path__
```

这里有一个注意事项，最后的一段，则是比 Prometheus 多出来的部分，因为 Promtail 必须需要一个 **path** 字段来获取采集日志的路径。

### docker_sd_configs(\[]Object)

https://grafana.com/docs/loki/latest/clients/promtail/configuration/#docker_sd_config

**host(STRING)** # Docker 守护进程的地址。通常设置为：`unix:///var/run/docker.sock`

**filters([]Object)** # 过滤器，用于过滤发现的容器。只有满足条件的容器的日志，才会被 Promtail 采集并上报。

> 可用的过滤器取决于上游 Docker 的 API：<https://docs.docker.com/engine/api/v1.41/#operation/ContainerList>，在这个链接中，可以从 Available filters 部分看到，等号左边就是 name 字段，等号右边就是 values 字段。
>
> 这个 name 与 values 的用法就像 `docker ps` 命令中的 `--filter` 标志，这个标志所使用的过滤器，也是符合 Docker API 中的 ContainerList。

- **name(STRING)** #
- **values([]STRING)** #

**refresh_interval(DURATION)** # 刷新间隔。每隔 refresh_interval 时间，从 Docker 的守护进程发现一次可以采集日志的容器。

#### 配置示例

```yaml
- job_name: flog_scrape
  docker_sd_configs:
    - host: unix:///var/run/docker.sock
      refresh_interval: 60s
  relabel_configs:
    - source_labels: ["__meta_docker_container_name"]
      regex: "/(.*)"
      target_label: "container"
```

## 重设标签

### relabel_configs(Object)

详见 [Label 与 Relabeling](/docs/6.可观测性/Logs/Loki/Promtail/Label%20与%20Relabeling.md)

# 配置文件示例

## 采集 Docker 容器日志

```yaml
server:
  http_listen_port: 9080
  grpc_listen_port: 0

positions:
  filename: /tmp/positions.yaml

clients:
  - url: http://gateway:3100/loki/api/v1/push
    tenant_id: tenant1

scrape_configs:
  - job_name: flog_scrape
    docker_sd_configs:
      - host: unix:///var/run/docker.sock
        refresh_interval: 5s
    relabel_configs:
      - source_labels: ["__meta_docker_container_name"]
        regex: "/(.*)"
        target_label: "container"
```
