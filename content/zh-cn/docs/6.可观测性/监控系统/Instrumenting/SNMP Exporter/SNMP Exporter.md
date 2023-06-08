---
title: SNMP Exporter
weight: 1
---

# 概述

> 参考：
>
> - [GitHub 项目，prometheus/snmp_exporter](https://github.com/prometheus/snmp_exporter)

Snmp Exporter 通过 snmp 采集监控数据，并转换成[ OpenMetrics 格式](</docs/6.可观测性/监控系统/监控系统概述/HTTP(新监控标准).md>>)的指标。

在这个项目中，有两个组件，

- **Exporter(导出器)** # 通过 snmp 抓去指标数据并转换成 OpenMetrics 格式
- [**Generator(生成器)**](https://github.com/prometheus/snmp_exporter/blob/master/generator) # 生成 Exporter 的配置文件。

## Exporter(导出器)

snmp_exporter 启动后默认监听在 9116 端口上。**snmp_exporter 会根据 snmp.yml 配置文件中的配置规则抓取 snmp 数据并转换成 Metrics 格式的数据。**

Prometheus Server 抓取 metircs 的 http 请求样例： `http://IP:PORT/snmp?module=if_mib&target=TargetIP` # 获取 TargetIP 上的 snmp 信息，并转换成 metrics 格式，其中 `module=if_mib` 是可省的，若不指定 module，则抓取所有 module。

snmp exporter 源码简单解析

```go
// 这个结构体实现了 prometheus.Collector 接口
type collector struct {
 // ......略
}
// 采集 Metrics 的主要逻辑在这里，这里省略了很多不相关的代码
func (c collector) Collect(ch chan<- prometheus.Metric) {
 // 这里使用了 gosnmp/gosnmp 这个库，通过这个库来执行类似 snmpwalk 这样的命令获取 snmp 数据
 pdus, err := ScrapeTarget(c.ctx, c.target, c.module, c.logger)
 oidToPdu := make(map[string]gosnmp.SnmpPDU, len(pdus))
PduLoop:
 // 为每个 pdu 查找匹配到的 Metrics
 for oid, pdu := range oidToPdu {
  head := metricTree
  oidList := oidToList(oid)
  for i, o := range oidList {
   var ok bool
   head, ok = head.children[o]
   if !ok {
    continue PduLoop
   }
   if head.metric != nil {
    // 在这里获取 snmp 数据并转换为 Metrics 格式的数据
    samples := pduToSamples(oidList[i+1:], &pdu, head.metric, oidToPdu, c.logger)
    for _, sample := range samples {
     ch <- sample
    }
    break
   }
  }
 }
}
```

## Generator(生成器)

用于生成 Exporter 运行时所需的配置文件(snmp.yml)。

为什么 snmp exporter 这个导出器的配置需要生成呢？~

首先得先从 snmp exporter 的运行逻辑说起，snmp exporter 的运行，必须依赖于 snmp.yml 这个配置文件。snmp.yaml 指明了每一个 OID 转换成 Metrics 之后的格式及内容。Prometheus Server 每次对 snmp exporter 发起 http 请求获取 Metrics 时，snmp exporter 都会使用 [gosnmp 这个第三方库中的功能](https://github.com/gosnmp/gosnmp) 向目标执行类似 snmpwalk 的命令，获取 snmp 的数据，并逐一将获取到的 snmp 数据转换为 Metrics 格式的数据。

基于这个运行机制，那么 snmp.yml 文件中，必须就必须包含 Metrics 的名字、OID、Metrics 的类型、Metircs 的帮助信息。而这些信息如何填写到这个文件中呢？总不能手写吧。。。这么多指标。。。查来查去，再手写进去。。。。无法想象。。。所以，此时就需要一个工具，可以根据某些信息来自动生成这个 snmp.yml 文件，而依据内容，当然就是 MIB 啦！所以 **Generator 将会根据 MIB 的内容，生成 snmp.yml 文件**。

- MIB 中的 DESCRIPTION，将会变为 Metrics 的帮助信息
- MIB 中的 Object 名称，将会变为 Metrics 的名称
- MIB 中 Object 的值类型，将会变为 Metrics 的值类型

Generator 使用 MIB 中的哪些信息、转换后是否需要设置标签、是否忽略某些 OID 等等这种行为，是由 generator.yml 文件进行控制的

**总结：Generator(生成器) 通过 `MIB 库文件` 以及 `generator.yml 文件` 这两种东西，来生成 snmp.yml 文件**

> MIB 库文件一般是放在 generator 程序运行时所在目录的 mibs 目录下的，generator.yml 文件一般是放在 generator 程序运行时所在目录下。
> 
> 如果运行 generator 时无法在 MIB 库文件中找到 generator.yml 文件中配置的 OID，则 generator 程序运行将会报错，提示无法找到对应的 Object。此时就需要将必要的 MIB 库文件，拷贝到 mibs/ 目录下即可。

generator.yml 文件详解见[此处](/docs/6.可观测性/监控系统/Instrumenting/SNMP%20Exporter/配置详解.md)

**generator.yml 文件最简单示例**

```yaml
modules:
  # Default IF-MIB interfaces table with ifIndex.
  if_mib:
    # 指定要获取的 OID，在生成 snmp.yaml 时，会根据这里面的定义去 MIB 中查找对应的 Object
    # 这个示例表示要获取 sysUpTime 与 ifXTable 这俩 OID 的数据以及 interfaces 这个 Object 组的数据
    # 用白话说就是对下面这些 Object 执行 walk 命令
    walk: [sysUpTime, interfaces, ifXTable]
    # 指定要使用的 snmp 版本
    version: 2
    # 获取 snmp 数据时的认证信息
    auth:
      community: public
```

# Snmp Exporter 部署

## 通过二进制方式安装

从 [GitHub Release](https://github.com/prometheus/snmp_exporter/releases) 处下载二进制文件，以后台方式运行即可

## 通过 docker 启动

```bash
docker run -d --name snmp_exporter --restart=always \
  --net=host \
  lchdzh/snmp-exporter:0.19.1 \
  --config.file=/etc/snmp_exporter/snmp.yml
```

注意：该项目没有现成的 docker 镜像，需要先手动构建，参考 [fork 到自己仓库的项目](https://github.com/DesistDaydream/snmp_exporter)，修改 Dockerfile 后，执行构建

```bash
git clone https://github.com/prometheus/snmp_exporter.git
cd snmp_exporter
go build .
docker build -t lchdzh/snmp-exporter:0.19.1 .
```
# snmp_exporter 配置

## Exporter 配置

**snmp.yml** # snmp_exporter 程序运行时根据该文件，将 OID 转换为 Metircs。通过 generator 程序自动生成。

> snmp_exporter 默认使用 snmp v2 来获取 snmp 格式的数据，如果想要修改配置 snmp_exporter 的配置文件，比如使用 snmp v3 的方式获取 snmp 格式的数据，则需要使用 generator 来生成配置文件

## Generator 配置

**mibs/\***# 用来存放 MIB 文件。
**generator.yml** # 用来配置生成 snmp.yml 的行为。

# Prometheus 中的 scrape_configs 配置示例

Prometheus 默认从 snmp_exporter 所监控的 9116 端口获取数据，其路径为：`http://localhost:9116/snmp?target=1.2.3.4`（其中 1.2.3.4 是 snmp_exporter 程序要获取 SNMP 数据的目的设备的 ip）

**有两种方式配置 job：**

这两种配置方式，虽然配置方式不一样，但是最终结果是一样的，只不过过程不一样，适应场景不一样

- 第一种配置方式可以极大简化配置文件，让所有要采集的 SNMP 目标，都包含在同一个 job 中，但是思维方式比较抽象。
- 第二种配置方式比较直观，但是每个要采集的 SNMP 目标都需要单独一个 job，并且配置文件内容过长。
  - 并且，第二种配置方式是 prometheus operator 的 serviceMonitor 自动生成的配置文件的格式。

## 第一种：一个 job 配置多个 SNMP 目标

```yaml
scrape_configs:
  - job_name: "snmp"
    static_configs:
      - targets:
          - 172.19.42.200 # 指定 snmp_exporter 要采集的目标设备ip
          - 172.19.42.243
    metrics_path: /snmp
    params:
      module: [if_mib]
    # 主要用来将 instance 标签的值修改为待采集的 snmp 目标，否则所有 instance 的值都是 snmp exporter 程序的 IP 了。
    # 最后将 __address__ 标签替换为 snmp_exporter 的监听地址
    relabel_configs:
      - source_labels: [__address__]
        target_label: __param_target
      - source_labels: [__param_target]
        target_label: instance
      - target_label: __address__
        replacement: 172.19.42.210:9116 # snmp_exporter 监听的IP:PORT
```

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/re86m3/1616069331853-cfb8bf5e-eef0-40b6-9f7f-51dc914e9bab.png)

为什么需要 relabel 呢？如果不写 relabel_configs 字段的话，会出现这种情况，endpoint 就是 target 中的值，效果如下图

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/re86m3/1616069331885-e0d59911-2da2-4ab7-83d2-50af6adbdc42.png)

## 第二种：一个 job 配置一个 SNMP 目标

```yaml
scrape_configs:
  - job_name: "snmp1"
    static_configs:
      - targets:
          - 172.19.42.210:9116 # 指定 snmp_exporter 监听的 IP:PORT
    metrics_path: /snmp
    params:
      module:
        - if_mib
      target:
        - 172.19.42.200 # 指定要采集的 SNMP 信息的目标设备 IP
  - job_name: "snmp2"
    static_configs:
      - targets:
          - 172.19.42.210:9116 # 指定 snmp_exporter 监听的 IP:PORT
    metrics_path: /snmp
    params:
      module:
        - if_mib
      target:
        - 172.19.42.243 # 指定要采集的 SNMP 信息的目标设备 IP
```

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/re86m3/1616069331867-03f9cfc0-f33f-40a5-8701-fdbb6089b8fa.png)

为什么要写多个 job 呢？因为如果将配置写成这样，将会出现下面这种情况

```yaml
scrape_configs:
  - job_name: "snmp1"
    static_configs:
      - targets:
          - 172.19.42.210:9116
    metrics_path: /snmp
    params:
      module:
        - if_mib
      target:
        - 172.19.42.200
        - 172.19.42.243
```

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/re86m3/1616069331849-200a90f1-f81e-40d1-99a9-bb96cc390971.png)

## 总结

上述两种配置实际上，都是使用类似下面的 URL 发送向 snmp_exporter 发送一个 http 的 GET 请求：

http://172.19.42.210:9116/snmp?module=if_mib&target=172.19.42.200
