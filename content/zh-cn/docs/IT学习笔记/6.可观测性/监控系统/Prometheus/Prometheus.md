---
title: Prometheus
---

# 概述

> 参考：
> - [官方文档](https://prometheus.io/docs/introduction/overview/)
> - [yunlzheng 写的电子书](https://yunlzheng.gitbook.io/prometheus-book/)
> - [GitHub 项目，Tencent-Cloud-Native/tkedocs](https://github.com/Tencent-Cloud-Native/tkedocs)(我个人总结完绝大部分文档后发现的这个项目)
> - Prometheus 纪录片
>   - [YouTube](https://www.youtube.com/watch?v=rT4fJNbfe14)
>   - [B 站翻译](https://www.bilibili.com/video/BV1aW4y147GX)

Prometheus 是由 SoundCloud 开发的 开源监控报警系统 和 时间序列数据库(TSDB) 。**Time Series(时间序列)** 概念详见：[Prometheus 数据模型](https://www.yuque.com/go/doc/33147376)。使用 Go 语言开发，是 Google BorgMon 监控系统的开源版本。

> 题外话：Google 的 Borg 诞生了 kuberntes、Google 的 Borgmon 诞生了 Prometheus

2016 年由 Google 发起 Linux 基金会旗下的 Cloud Native Computing Foundation(云原生计算基金会), 将 Prometheus 纳入其下第二大开源项目。Prometheus 目前在开源社区相当活跃。

## Prometheus 架构概述

Prometheus 的基本原理是通过 HTTP 协议周期性抓取被监控组件的状态，任意组件只要提供对应的 HTTP 接口就可以接入监控。不需要任何 SDK 或者其他的集成过程。这样做非常适合做虚拟化环境监控系统，比如 VM、Docker、Kubernetes 等。输出被监控组件信息的 HTTP 接口被叫做 exporter 。

下面这张图说明了 Prometheus 的整体架构，以及生态中的一些组件作用：
![](https://notes-learning.oss-cn-beijing.aliyuncs.com/usvgfl/1616041189523-5ca97287-5886-4ab9-a4f8-6c249117e314.jpeg)
Prometheus 生态圈中包含了多个组件，其中许多组件是可选的，多数 Prometheus 组件是 Go 语言写的，使得这些组件很容易编译和部署：

- **Prometheus Server** # 主要负责数据抓取和存储，提供 PromQL 查询语言的支持。用于收集和存储时间序列数据。
  - 定期从配置好的 Jobs 中**拉取 **Exporters 采集的** Metrics(指标)** 数据；或者**接收**来自 **Pushgateway**(类似 zabbix 的 proxy) 发过来的 Metrics；或者从其他的 Prometheus Server 中拉取 Metrics。
  - Prometheus Server 在本地存储收集到的 Metrics，并通过一定 **RecordingRule(记录规则)** 进行清理和整理数据，并把得到的结果存储到新的时间序列中。还会运行已定义好的 **AlertingRule(告警规则)**，记录新的时间序列或者向 Alertmanager 推送警报。
  - 由于 Metrics 都是通过 HTTP 或者 HTTPS 协议提供的，所以 Prometheus Server 在抓取 Metrics 时，也就是发起一次 HTTP 或者 HTTPS 的 GET 请求
- **Instrumenting** # 为 Prometheus 提供指标的工具或代码
  - **Exporters **# 导出器。Exporter 是 Prometheus 的一类数据采集组件的总称。它负责从设备上搜集数据，并将其转化为 Prometheus 支持的格式(一般情况下 exporter 是安装在需要采集数据的设备上的程序，并监听某个 port。但是如果想要收集 snmp 信息的话，则有专门的 snmp-exporter 安装在某个地方；再收集指定设备的 snmp 信息，然后 prometheus 再找 snmp-exporter 去收集数据)。与传统的数据采集组件不同的是，它并不向中央服务器发送数据，而是等待中央服务器主动前来抓取。Prometheus 提供多种类型的 Exporter 用于采集各种不同服务的运行状态。目前支持的有数据库、硬件、消息中间件、存储系统、HTTP 服务器、JMX 等。
  - **Client Library** # 客户端库(客户端 SDK)，官方提供的客户端类库有 go、java、scala、python、ruby，其他还有很多第三方开发的类库，支持 nodejs、php、erlang 等。为需要监控的服务生成相应的 Metrics 并暴露给 Prometheus server。当 Prometheus server 来 pull 时，直接返回实时状态的 Metrics。
  - **Push Gateway **# 支持 Client 主动推送 Metrics 到 PushGateway，而 PrometheusServer 只是定时去 Gateway 上抓取数据。
- **Alertmanager** # 警告管理器，用来进行报警。从 Prometheus server 端接收到 alerts 后，会进行去除重复数据，分组，并路由到对收的接受方式，发出报警。常见的接收方式有：电子邮件，pagerduty，OpsGenie, webhook 等。
- **prometheus_cli** # 命令行工具。
- **其他辅助性工具**
  - Prometheus 通过 PromQL 和其他 API 可视化地展示收集的数据。Prometheus 支持很多方式的图表可视化，例如 Grafana、自带的 PrometheusDashboard 以及自身提供的模版引擎等等。Prometheus 还提供 HTTP API 的查询方式，自定义所需要的输出。

Prometheus 适用的场景

- Prometheus 在记录纯数字时间序列方面表现非常好。它既适用于面向服务器等硬件指标的监控，也适用于高动态的面向服务架构的监控。对于现在流行的微服务，Prometheus 的多维度数据收集和数据筛选查询语言也是非常的强大。Prometheus 是为服务的可靠性而设计的，当服务出现故障时，它可以使你快速定位和诊断问题。它的搭建过程对硬件和服务没有很强的依赖关系。

Prometheus 不适用的场景

- Prometheus 它的价值在于可靠性，甚至在很恶劣的环境下，你都可以随时访问它和查看系统服务各种指标的统计信息。 如果你对统计数据需要 100%的精确，它并不适用，例如：它不适用于实时计费系统。

### 总结：prometheus 从 Instrumenting 那里抓取监控数据，储存。完了~~~~哈哈哈哈哈

## Instrumenting(检测仪表装置) 的实现方式

Prometheus 可以通过 3 种方式从目标上 Scrape(抓取) 指标：

- **Exporters **# 外部抓取程序
- **Instrumentation** # 可以理解为内嵌的 Exporter，比如 Prometheus Server 的 9090 端口的 `/metrics` 就属于此类。
  - 说白了，就是目标自己就可以吐出符合 Prometheus 格式的指标数据
- **Pushgateway** # 针对需要推送指标的应用

## Label 与 Relabeling

详见 [Label 与 Relabeling 章节](✏IT 学习笔记/👀6.可观测性/监控系统/Prometheus/Target(目标)%20 与%20Relabeling(重新标记).md 与 Relabeling(重新标记).md)

## Instrumenting 的安装与使用

详见 [Instrumenting 章节](/docs/IT学习笔记/6.可观测性/监控系统/Instrumenting/Instrumenting.md)

# Prometheus 部署

> 参考：
> - [官方文档，Prometheus-安装](https://prometheus.io/docs/prometheus/latest/installation/)

官方系统版本可在这里下载：<https://prometheus.io/download/>

Prometheus 官方有多种部署方案，比如：Docker 容器、Ansible、Chef、Puppet、Saltstack 等。Prometheus 用 Golang 实现，因此具有天然可移植性(支持 Linux、Windows、macOS 和 Freebsd)。

## 二进制文件运行 Prometheus Server

- <https://github.com/prometheus/prometheus/releases/> 在该页面下直接下载 prometheus 的进制文件 `prometheus-版本号.linux-amd64.tar.gz` 并解压，其中包含 prometheus 的主程序还有 yaml 格式的配置文件以及运行所需要的依赖库

```bash
export PromVersion="2.25.1"
wget https://github.com/prometheus/prometheus/releases/download/v${PromVersion}/prometheus-${PromVersion}.linux-amd64.tar.gz
```

- 创建/usr/local/prometheus 目录，并将解压的所有文件移动到该目录下

```bash
mkdir /usr/local/prometheus
tar -zxvf prometheus-${PromVersion}.linux-amd64.tar.gz -C /usr/local/prometheus/ --strip-components=1
```

- 创建 Systemd 服务,在 ExecStart 字段上，使用运行时标志来对 prometheus 进行基本运行配置，标志说明详见下文

```bash
cat > /usr/lib/systemd/system/prometheus.service << EOF
[Unit]
Description=prometheus
After=network.target
[Service]
Type=simple
User=root
ExecStart=/usr/local/prometheus/prometheus \
  --web.console.templates=/usr/local/prometheus/consoles \
  --web.console.libraries=/usr/local/prometheus/console_libraries \
  --config.file=/usr/local/prometheus/prometheus.yml \
  --storage.tsdb.path=/var/lib/prometheusData \
  --web.enable-lifecycle
Restart=on-failure
[Install]
WantedBy=multi-user.target
EOF
```

- 启动 Prometheus

```bash
systemctl start prometheus
```

## 容器运行 prometheus

获取配置文件.

> 也可以不获取配置文件，去掉启动时的 -v /etc/monitoring/prometheus:/etc/prometheus/config_out 与 --config.file=/etc/prometheus/config_out/prometheus.yml 这两行即可
> 获取配置文件主要是为了让后续测试可以更方便得修改文件

```bash
mkdir -p /opt/monitoring/prometheus
docker run -d --name prometheus --rm prom/prometheus
docker cp prometheus:/etc/prometheus /opt/monitoring/prometheus
mv /opt/monitoring/prometheus/prometheus /opt/monitoring/prometheus/config
docker stop prometheus
```

运行 Prometheus Server

```bash
docker run -d --name prometheus --restart=always \
  --network host \
  -v /etc/localtime:/etc/localtime \
  -v /opt/monitoring/prometheus/config:/etc/prometheus/config_out \
  prom/prometheus \
  --config.file=/etc/prometheus/config_out/prometheus.yml
```

# Prometheus 关联文件与配置

**/etc/prometheus/prometheus.yml **# Prometheus Server 运行时的配置文件。可通过 --config.file 标志指定其他文件。
**/etc/prometheus/rule.yml** # Prometheus Rule 配置文件。该文件默认不存在，需手动创建。可以在 prometheus.yml 配置中指定其他文件。

## Prometheus 配置示例

### 默认配置文件

```yaml
# 全局配置
global:
  scrape_interval:     15s # 默认抓取间隔, 15秒向目标抓取一次数据。
  evaluation_interval: 15s # 每15秒评估一次规则，默认为1分钟。
  # scrape_timeout is set to the global default (10s).

 # 告警报警配置，设置prometheus主程序对接alertmanager程序的
 alerting:
  alertmanagers:
  - static_configs:
    - targets:
      # - alertmanagerIP:9093

# 指定要使用的规则文件位置，并加载一次，根据全局配置中的 evaluation_interval 来定期评估
# 从所有匹配到的文件中读取配置内容。可以使用正则表达式匹配多个符合的文件。Prometheus支持两种规则
# 其一是记录规则(recording rules)
# 其二是告警规则(alerting rules)
rule_files:
  # - "first.rules"
  # - "second.rules"

# 抓取配置，prometheus抓取exporter上的数据时的配置，一个job就是一个抓取工作，其中可以包括1个或者多个目标
# 目标指的是可以被 prometheus 采集的服务器、服务等。
# 默认配置里，Prometheus Server 会抓取本地9090端口上数据。该端口上的 exporter 就是 PrometheusServer 自己的 exporter
scrape_configs:
# job_name 指定要 scrape(抓取) 的 job(工作) 名称，名称必须是唯一的
# 并且在这个配置内的时间序例，每一条都会自动添加上这个{job_name:"prometheus"}的标签。
- job_name: 'prometheus'
  # 设定该job的抓取时间间隔
  scrape_interval: 5s
  static_configs:
  - targets: ['localhost:9090']
```

### 具有 node_exporter 的配置简单文件

抓取部署了 node_exporter 设备的监控数据的方式及 prometheus.yml 配置文件说明
prometheus 会从 Node Exporter 所在服务器的 http://IP:9100/metrics 这个地址里的内容来获取该设备的监控数据
所以需要给 prometheus 创建一个工作(i.e.job)。一个 job 就是一个抓取监控数据的工作，其中包括要抓取目标的 ip 和 port，还可以设置标签进行分类，还能进行抓取筛选等等，下面提供一个基本的配置
修改 prometheus.yml，加入下面的监控目标，以便让 prometheus 监控这个已经安装了 node_exporter 的设备

```yaml
- job_name: "linux" #新增一个job，名为linux
  static_configs: # 使用静态配置
    - targets: ["10.10.100.101:9100"] #添加一个要抓取数据的目标，指定IP与PORT 。node_exporter所安装的设备
      labels:
        instance: lchTest #给该目标添加一个标签
```

现在，prometheus.yml 配置文件中中一共定义了两个监控：一个是默认自带监控 prometheus 自身服务，另一个是我们新增的 job，这个 job 就是要抓取目标是 10.10.100.101 这台服务器上的监控数据

```yaml
scrape_configs:
  - job_name: "prometheus"
    static_configs:
      - targets: ["localhost:9090"]
  - job_name: "linux" #指定job名称
    static_configs: #设定静态配置
      - targets: ["10.10.100.101:9100"] #指定node_exporter所安装设备的ip:port
        labels:
          instance: lchTest #给该target一个label来分类，常用于在查询语句中的筛选条件
```

访问 Prometheus Web，在 Status->Targets 页面下，我们可以看到我们配置的两个 Target，它们的 State 为 UP
![](https://notes-learning.oss-cn-beijing.aliyuncs.com/usvgfl/1616041189541-1dfdddd7-ee74-4f32-8df6-8821cf415a14.jpeg)

# Prometheus 的基本使用方式

Prometheus 运行后默认会监听在 9090 端口，可以通过访问 9090 端口来打开 Prometheus 的 web 界面
![](https://notes-learning.oss-cn-beijing.aliyuncs.com/usvgfl/1616041189523-861a164c-3f79-42af-bd88-44c4baf2e349.jpeg)
Prometheus 本身也是自带 exporter 的,我们通过请求 http://ip:9090/metrics 可以查看从 exporter 中能具体抓到哪些 metrics。
![](https://notes-learning.oss-cn-beijing.aliyuncs.com/usvgfl/1616041189563-4125f137-160f-48dd-b4f6-dfd6af94aed0.jpeg)

这里以 Prometheus 本身数据为例，简单演示下在 Web 中查询指定表达式及图形化显示查询结果。
![](https://notes-learning.oss-cn-beijing.aliyuncs.com/usvgfl/1616041189526-ee545ef0-965e-499c-b80f-b6cdaf05c974.jpeg)
[

# ](Instrumenting.md)
