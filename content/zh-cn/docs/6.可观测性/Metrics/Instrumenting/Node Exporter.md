---
title: Node Exporter
---

# 概述

> 参考：
>
> - [GitHub 项目，prometheus/node_exporter](https://github.com/prometheus/node_exporter)

Node Exporter 用于收集服务器的 metrics，比如 内存、cpu、磁盘、I/O、电源、etc. 。Node Exporter 将采集各种指标的代码逻辑抽象称为 Node 的 **Collector(采集器)**。每类指标都对应一个 Collector，比如 cpu 采集器、meminfo 采集器、etc. 这些名称通常都能直观得看到想要采集的指标是什么

node_exporter 默认监听在 9100 端口上。

Prometheus Server 抓取 metrics 的位置 http://IP:9100/metrics # 获取 node_exporter 所在主机的 metrics 信息

# Node Exporter 部署

## 二进制文件安装 node_exporter

为监控服务器 CPU、内存、磁盘、I/O 等信息，首先需要安装 node_exporter。node_exporter 的作用是服务器数据收集。

- 下载 node_exporter，过程基本与使用 prometheus 程序一样。[下载页面在此](https://github.com/prometheus/node_exporter/releases)

```bash
export VERSION="1.6.0"
wget https://github.com/prometheus/node_exporter/releases/download/v${VERSION}/node_exporter-${VERSION}.linux-amd64.tar.gz
# 解压
mkdir -p /usr/local/prometheus/node_exporter
tar -zxvf node_exporter-${VERSION}.linux-amd64.tar.gz -C /usr/local/prometheus/node_exporter --strip-components=1
```

- 创建 node_exporter 的 systemd 服务，在[这里](https://github.com/prometheus/node_exporter/tree/master/examples/systemd)可以看到官方提供的 systemd 样例

```bash
tee /usr/lib/systemd/system/node-exporter.service > /dev/null << EOF
[Unit]
Description=node_exporter
After=network.target
[Service]
Type=simple
User=root
ExecStart=/usr/local/prometheus/node_exporter/node_exporter \
  --collector.filesystem.ignored-mount-points='^/(dev|proc|sys|var/lib/docker/.+|var/lib/kubelet/pods/.+)($|/)'
Restart=on-failure
[Install]
WantedBy=multi-user.target
EOF
```

- 启动 node_exporter 服务，该服务会默认监听在 9100 端口上，等待 prometheus 主程序来抓取监控数据

```bash
systemctl enable node-exporter --now
```

注意事项：

- 报错 `acpi: no handler for region [powr]`，需要添加参数--no-collector.hwmon，原因应该是与 dell 的硬件信息采集程序冲突

## 容器安装 node_exporter

<https://github.com/prometheus/node_exporter#using-docker>

```bash
export VERSION="1.6.0"

docker run -d --name node-exporter --restart=always \
  --net="host" \
  --pid="host" \
  -v "/:/host:ro,rslave" \
  prom/node-exporter:v${VERSION} \
  --web.listen-address=":9100" \
  --path.rootfs=/host \
  --no-collector.hwmon \
  --no-collector.wifi \
  --collector.filesystem.ignored-mount-points='^/(dev|proc|sys|var/lib/docker/.+|var/lib/kubelet/pods/.+)($|/)'
```

# node_exporter 可采集的数据种类

https://github.com/prometheus/node_exporter?tab=readme-ov-file#collectors

[这里](https://github.com/prometheus/node_exporter#enabled-by-default)有 node_exporter 默认启用的采集器，name 就是采集器名称

[这里](https://github.com/prometheus/node_exporter#disabled-by-default)有 node_exporter 默认禁用的采集器，name 就是采集器名称

如果想要让 node_exporter 采集或者不采集某些数据，可以在启动 node_exporter 程序时，向该程序传递参数。参数中的 NAME 为上面两个连接中，表格中的 name 列

- `--collector.<NAME>` # 标志来启用采集目标。
- `--no-collector.<NAME>` # 标志来禁用采集目标。

## 只让部分已启用的采集器采集指标

[这部分](https://github.com/prometheus/node_exporter/blob/v1.8.1/node_exporter.go#L78) `filters := r.URL.Query()["collect[]"]` 代码是用来可以让服务端在向 node-exporter 发起的 HTTP 请求中，在 [URL](/docs/4.数据通信/通信协议/HTTP/URL%20与%20URI.md) 的 QUERY 部分加入一些内容，以决定采集哪些 Metrics，而不必强制通过本身的 CLI 参数决定。参考 README 的 [Filtering enabled collectors](https://github.com/prometheus/node_exporter#filtering-enabled-collectors)。

> [!Note]
> URL Query 中填写的内容是指采集的，只要使用了 URL Query，那么 node-exporter 则只采集 Query 中指定的指标，其余的全都不采集

若是用 curl 发起请求，就是这样的: `curl 'localhost:9100/metrics?collect[]=cpu&collect[]=meminfo'`。这个表示只让 cpu 和 meminfo 这两个采集器工作采集 cpu 和 内存 的指标。

# Textfile Collector 文本文件采集器

> 参考：
>
> - 官方文档：<https://github.com/prometheus/node_exporter#textfile-collector>
> - 脚本样例：<https://github.com/prometheus-community/node-exporter-textfile-collector-scripts>
> - [公众号,k8s 技术圈-使用 Node Exporter 自定义监控指标](https://mp.weixin.qq.com/s/X73XRrhU_lYMvkJvF1z2uw)

文本采集器逻辑：

1. 在启动 node_exporter 的时候，使用 --collector.textfile.directory=/PATH/TO/DIR 参数。
2. 我们可以将自己收集的数据按照 prometheus 文本格式类型的 metrics 存储到到指定目下的 \*.prom 文件中
   1. 文本格式的 metrics 详见 2.1.Prometehus 时间序列数据介绍.note 文章末尾
3. 每当 prometheus server 到 node_exporter 拉取数据时，node_exporter 会自动解析所有指定目录下的 \*.prom 文件并提供给 prometheus

Note：

1. 采集到的数据没有时间戳
2. 必须使用后缀为 .prom 作为文件的名称
3. 文件中的格式必须符合 prometheus 文本格式的 metrics，否则无法识别会报错。文本格式样例详见：2.1.Prometehus 时间序列数据介绍.note
4. 若用 docker 启动的 node_exporter，则还需要使用 -v 参数，将 textfile 文件所在目录挂载进容器中。

# 为 Node Exporter 添加认证

注意：该功能为实验性的，笔记时间：2021 年 8 月 4 日

![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/ck9bpv/1628068010926-2ae85ce6-13be-4dd3-8ed1-74538c5cf3da.png)

与 Prometheus 添加认证的方式一样，详见：

[Authentication(认证)](/docs/6.可观测性/Metrics/Prometheus/HTTPS%20和%20Authentication.md)

node-exporter 程序使用 `--web.config` 命令行标志来指定 web-config 文件，读取其中内容并开启 TLS 或 认证功能。

# 源码解析

node_exporter.go 中的 `handler.innerHandler()` 方法用于创建 Node 采集器，i.e. 决定要启用哪些 Collector

```go
func (h *handler) innerHandler(filters ...string) (http.Handler, error) {
    // NewNodeCollector()` 方法决定启用哪些 Collector 的主要逻辑，该方法实例化了一个实现了 prometheus.Collector{} 接口的 NodeCollector{} 结构体
 nc, err := collector.NewNodeCollector(h.logger, filters...)

    // ......输出一些信息

 r := prometheus.NewRegistry()
 // 实例化后的 NodeCollector{} 使用 `prometheus.NewRegistry().Register()` 进行注册
 err := r.Register(nc)

   // ......最后就是标准的利用 promhttp.HandlerFor 或 promhttp.InstrumentMetricHandler 返回 http.Handler。具体用哪个以及其中的具体逻辑，与开启哪些 Node 的采集器没有强关联。
}
```

## 日志时区

详见 [Prometheus MGMT - Prometheus UTS 时区问题](/docs/6.可观测性/Metrics/Prometheus/Prometheus%20MGMT/Prometheus%20MGMT.md#Prometheus%20UTS%20时区问题)

# Grafana 面板

- 1860
  - https://github.com/rfmoz/grafana-dashboards
- [8919](https://grafana.com/grafana/dashboards/8919-1-node-exporter-for-prometheus-dashboard-cn-0413-consulmanager/)
  - 国人出的，22 年4月12日之后不维护了
  - [16098](https://grafana.com/grafana/dashboards/16098-1-node-exporter-for-prometheus-dashboard-cn-0417-job/) 新的，代替 8919

# 指标解析

## node_disk_io_time_seconds_total

`node_disk_io_time_seconds_total` 每秒钟增长的值最多就是 1。使用 `irate(node_disk_io_time_seconds_total[5m])` 可以得出磁盘 I/O 的使用率，结果是 0 - 1 之间的数。

原因如下：

在 [Node Exporter 源码](https://github.com/prometheus/node_exporter/blob/v1.6.1/collector/diskstats_linux.go#L320) 可以看到 `node_disk_io_time_seconds_total` 指标是 `diskstatsCollector.descs` 的第 10 号元素。在 `ch <- c.descs[i].mustNewConstMetric(val, dev)` 这里实现了指标写入逻辑

```go
func (c *diskstatsCollector) Update(ch chan<- prometheus.Metric) error {
  // prometheus/procfs 项目中的方法，返回的 IOStats 储存的是 /sys/block/<BLOCK>/stat 文件中的值
	diskStats, err := c.fs.ProcDiskstats()
	for _, stats := range diskStats {
		for i, val := range []float64{
			// ......前 9 个
			float64(stats.IOsTotalTicks) * secondsPerTick, // 这就是 node_disk_io_time_seconds_total 指标
			// ......后 7 个
		}
		ch <- c.descs[i].mustNewConstMetric(val, dev)
	}
}
```

写入的 `stats.IOsTotalTicks` 最终来源是 prometheus/procfs 项目中的 [IOStats.IOsTotalTicks 属性](https://github.com/prometheus/procfs/blob/v0.15.1/blockdevice/stats.go#L61)，IOStats 结构体的信息来源遵循以下几个内核文档的说明

- https://www.kernel.org/doc/Documentation/iostats.txt,
- https://www.kernel.org/doc/Documentation/block/stat.txt
- https://www.kernel.org/doc/Documentation/ABI/testing/procfs-diskstats

也就是说，IOStats 储存的是 `/sys/block/<BLOCK>/stat` 文件中的值，IOsTotalTics 文件中第 10 个字段的值，io_ticks 的值）。

所以，`node_disk_io_time_seconds_total` 本质就是 io_ticks，作为块设备的 I/O 时间如何理解详见 [Block 设备的 I/O 时间](/docs/1.操作系统/Kernel/Hardware/Block.md#I/O%20时间)

最后，node-exporter 代码中，把 io_ticks 除以 1000 得到了 秒 级别的 I/O 时间

```go
secondsPerTick = 1.0 / 1000.0
float64(stats.IOsTotalTicks) * secondsPerTick
```

所以，`node_disk_io_time_seconds_total` 每秒钟增长的值最多就是 1，那么 rate 之后的结果就可以作为使用率。
