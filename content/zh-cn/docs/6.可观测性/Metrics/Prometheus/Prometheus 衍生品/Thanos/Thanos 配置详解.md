---
title: Thanos 配置详解
---

# 概述

> 参考：

# Sidecar 配置

## [命令行标志](https://thanos.io/tip/components/sidecar.md/#flags)

**--grpc-addresss=\<STRING>** # 暴露的 StoreAPI 端点。默认值：`0.0.0.0:10901`
**--http-address=\<STRING>** # 监听的 HTTP 端点。/metrics 端点暴露指标。默认值：`0.0.0.0:10902`
**--objstore.config-file=\<FILE>**# 对象存储的配置信息。Sidecar 根据该配置，将 Prometheus 自身保存的数据转存到配置的对象存储中。
**--prometheus.url=\<STRING>**# 与 Prometheus Server 交互的地址。默认值：`http://localhost:9090`
**--tsdb.path=\<PATH>** # Prometheus 存储时间序列数据的路径。默认值：`./data`
**--shipper.upload-compacted** # 开启已压缩数据转存功能。该标志对迁移数据很有用。

## 配置文件

# Receiver 配置

## [命令行标志](https://thanos.io/tip/components/receive.md/#flags)

**--grpc-address=\<STRING>** # 暴露的 StoreAPI 端点。`默认值：0.0.0.0:10901`
**--http-address=\<STRING>** # 监听的 HTTP 端点。/metrics 端点暴露指标。`默认值：0.0.0.0:10902`
**--label=\<KEY="VALUE">** # 为所有序列创建的标签，多个标签指定多个 --label 标志。与 prometheus 配置 external_labels 字段效果一样
**--objstore.config-file=\<FILE>**# 对象存储的配置信息。Receiver 根据该配置，将时序数据转存到对象存储中。
**--receive.hashrings-file=/etc/thanos/receiver-hashring.json** #
**--receive.local-endpoint=127.0.0.1:10901** #
**--receive.replication-factor=1** #
**--remote-write.address=\<STRING>** # 处理 Prometheus 的 Remote Write 请求的地址和端口。`默认值：0.0.0.0:19291`
**--tsdb.path=\<PATH>** # Receiver 的 TSDB 存储数据的路径。`默认值：./data`
**--tsdb.retention=\<DURATION>** # Receiver 的 TSDB 中存储数据的时长。`默认值：15d`
**--tsdb.wal-compression** # 开启压缩 TSDB 的 WAL 功能。

## 配置文件

# Query 配置

## [命令行标志](https://thanos.io/tip/components/query.md/#flags)

**--store=\<URL>** # 发起 PromQL 查询请求时的目标，若有多个查询目标，可以配置多个该选项。通常用来指定某些组件暴露的 StoreAPI。
**--query.lookback-delta=\<DURATION>** # 评估 PromQL 表达式时最大的回溯时间。`默认值：5m`

- 比如，当采集目标的间隔时间为 10m 时，由于该设置，最大只能查询当前时间的前 5m 的数据，这是，即时向量表达式返回的结果将会为空。

## 配置文件

# Store 配置

## [命令行标志](https://thanos.io/tip/components/store.md/#flags)

**--grpc-address=\<STRING>** # 暴露的 StoreAPI 端点。`默认值：0.0.0.0:10901`
**--http-address=\<STRING>** # 监听的 HTTP 端点。/metrics 端点暴露指标。`默认值：0.0.0.0:10902`
**--objstore.config-file=\<FILE>**# 对象存储的配置信息。Store 收到 PromQL 查询请求后，将会根据该配置文件的内容，去相应的对象存储中查询数据。

## 配置文件

# Compactor 配置

## [命令行标志](https://thanos.io/tip/components/compact.md/#flags)

**--log.format=\<STRING>** # 输出的日志格式。`默认值：logfmt`

- 可用的格式有 logfmt、json

**--log.level=\<STRING>** # 输出日志的级别。`默认值：info`
**--retention.resolution-raw=\<DURATION>** # 对象存储中，原始数据的保留时长。
**--retention.resolution-5m=\<DURATION>** # 对象存储中，降采样到样本间隔 5 分钟的数据的保留时长。
**--retention.resolution-1h=\<DURATION>** # 对象存储中，降采样到样本间隔 1 小时的数据的保留时长。

- 通常建议它们的存放时间递增配置（一般只有比较新的数据才会放大看，久远的数据通常只会使用大时间范围查询来看个大致，所以建议将精细程度低的数据存放更长时间）。

**--wait** # 让 Compactor 持续运行，而不是压缩完成后结束进程
**--wait-interval=\<DURATION>** # 每次压缩任务的等待时间和桶刷新时间。仅在指定了 --wait 标志后才起作用。

## 配置文件

# Ruler 配置

## 命令行标志

## 配置文件
