---
title: Thanos
---

# 概述

> 参考：
> - [GitHub 项目，thanos-io/thanos](https://github.com/thanos-io/thanos)
> - [官网](https://thanos.io/)
> - [K8S 训练营，Kubernetes 监控-Thanos](https://www.qikqiak.com/k8strain/monitor/thanos/#thanos)
> - [公众号-k8s 技术圈，使用 Thanos 集中管理多 Prometheus 实例数据](https://mp.weixin.qq.com/s/VRnf0BMZfWSM4QobsHSiYg)
> - [公众号-k8s 技术圈，Thanos Ruler 组件的使用](https://mp.weixin.qq.com/s/o3Gr8X5Br3DsOATvXXlZBw)

首先需要明确一点，Thanos 是一组程序的统称。这一组程序可以组成具有无限存储容量的高可用 Metrics 系统。可以将其无缝添加到现有 Prometheus 之上。

单独使用 Prometheus 可能产生的问题

- 长、短期数据未分层，一套 Prometheus 在被查询长周期指标时，Prometheus 所在服务器的内存、CPU 使用率飙升，甚至可能导致监控、告警服务不可用，原因在于两点：
  - 查询长周期数据时，Prometheus 会将大量数据载入内存
  - Prometheus 载入的不是降采样数据
- 查询的时间范围越大，需要的内存就越多。在另一个生产的方案中，采用 VictoriaMetrics 单机版作为远端存储，服务器内存高达 128G。同时，这种方式还存在丢数据的情况。
- Prometheus 联邦的方式，只是结局了将多个 Prometheus 聚合起来的情况，并没有提供抽样的能力，不能加快长周期指标的查询，不适用于当前远端存储的场景。

综上所属

- Thanos Compact 组件能对指标数据进行降采样，以提高大时间范围查询的效率
- Thanos 的 Sidecar 和 Receiver 组件都可以将指标数据转存到对象存储中
- Thanos Querier 组件可以
- 此时，通过 Thanos 将数据分了层
  - 短期数据保存在 Receiver 或 Prometheus 中，用于告警系统的高频查询以及 Grafana 的展示
  - 长期数据保存在对象存储中，以供后续分析使用

# Thanos 架构概述

Thanos 遵循 [KISS ](https://en.wikipedia.org/wiki/KISS_principle)和 Unix 哲学，由一组组件组成，每个组件都可以实现特定的功能：

- Prometheus 本地数据处理
  - **Sidecar** # 暴露 StoreAPI。连接 Prometheus，读取其数据以便进行下一步处理。将已压缩数据上传到对象存储，或读取其数据以进行查询和/或将其上传到云存储中。还可以动态重载 Prometheus 配置文件。
  - **Receiver **# 暴露 StoreAPI。实现了 Prometheus 的 Remote Write API。接收 Prometheus 发送过来的 WAL 数据存储到 Receiver 本地的 TSDB 中，也就是说，Receiver 本身就是一个 TSDB。同时，也可以将这些数据上传到对象存储中。
- PromQL 查询处理。分为两个部分，一个前端一个后端。
  - **Querier** # 实现 Prometheus 的 API、以及一个类似 Prometheus 的 Web 页面。Querier 从 StoreAPI 查询数据并返回结果。
  - **Query Frontend **# 实现 Prometheus 的 API、以及一个类似 Prometheus 的 Web 页面。可以将请求负载均衡到指定的多个 Querier 上，同时可以缓存响应数据、也可以按查询日拆分。有点像 Redis 的效果。
- 对象存储中的数据处理
  - **Store Gateway **# 暴露 StoreAPI。将对象存储中的数据暴露出来，供 Querier 组件通过 PromQL 查询
  - **Compactor **# 压缩数据，对保存在对象存储中的数据进行降采样。
- 其他
  - **Ruler **# 针对 Thanos 中的数据评估记录和警报规则以进行展示/上传。

通过 thanos 程序的子命令可以运行指定的组件

```bash
root@lichenhao:/mnt# docker run --rm quay.io/thanos/thanos:v0.20.0 --help
usage: thanos [<flags>] <command> [<args> ...]

A block storage based long-term storage for Prometheus.

Flags:
  -h, --help               Show context-sensitive help (also try --help-long and
...... 略

Commands:
  help [<command>...]
    Show help.

  sidecar [<flags>]
    Sidecar for Prometheus server.

  store [<flags>]
    Store node giving access to blocks in a bucket provider. Now supported GCS,
    S3, Azure, Swift, Tencent COS and Aliyun OSS.

  query [<flags>]
    Query node exposing PromQL enabled Query API with data retrieved from
    multiple store nodes.

  rule [<flags>]
    Ruler evaluating Prometheus rules against given Query nodes, exposing Store
    API and storing old blocks in bucket.

  compact [<flags>]
    Continuously compacts blocks in an object store bucket.
  ...... 略
  receive [<flags>]
    Accept Prometheus remote write API requests and write to local tsdb.

  query-frontend [<flags>]
    Query frontend command implements a service deployed in front of queriers to
    improve query parallelization and caching.
```

可以看到：

- **Sidecar** # 由 **thanos sidecar** 子命令实现。
- **Receiver **# 由 **thanos receive** 子命令实现。
- **Store Gateway **# 由 **thanos store** 子命令实现。
- **Querier** # 由 **thanos query** 子命令实现。
- **Query Frontend** # 由 **thanos query-frontend** 子命令实现
- **Compactor **# 由 **thanos compact** 子命令实现。
- **Ruler **# 由 **thanos rule** 子命令实现。

thanos 二进制文件的设计方式与 loki 非常类似，都是在单一二进制文件中，可以运行指定的一个或多个组件。

## StoreAPI

**StoreAPI **并不算一个组件，而是 Thanos 的一个核心概念。Sidecar、Store、Receiver 组件可以暴露 StoreAPI 端点。而 Query 可以使用其他组件爆出来的 StoreAPI 端点。

Query 向 StoreAPI 发送 PromQL 查询请求，Sidecar、Store、Receiver 在 StoreAPI 上收到 PromQL 请求后，从各自关联的后端数据库中查询数据，并返回给 Query。

StoreAPI 一般都是 gRPC 调用。

## 架构分类

在 Thanos 中，Sidecar 和 Receiver 组件是最主要的组件，这俩组件，可以让 Thanos 实现两种架构模型。

- Sidecar 架构 # 通过推送方式，将时序数据转存到对象存储中
- Receiver 架构 # 自身就是一个 TSDB，也实现了 Prometheus 的 Remote Write API。接收 Prometheus 发送过来的 WAL 数据存储到 Receiver 本地的 TSDB 中。同时，也可以将这些数据上传到对象存储中。

### Sidecar 架构

![Thanos High Level Arch Diagram.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/gq2u99/1620877646742-0d57551f-e57e-4ecb-a5d6-b2c14ad6f979.png)

### Receiver 架构

![Thanos High Level Arch Diagram Receive v0.19.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/gq2u99/1620877652987-b4b6f910-11f7-4fb1-b74a-ac5cf02c772b.png)

## 总结

Compactor、Sidecar、Receiver、Ruler 是唯一对对象存储具有写权限的组件，且只有 Comactor 才能删除数据。

- 通常情况下，Compactor 不需要高可用，因为 Compactor 没有针对所有对象存储提供的安全锁定机制，所以当多个 Compactor 同时操作同一个 Bucket 内的对象，则会出现争抢问题，导致数据异常。

### Sidecar 与 Receiver 两种架构的对比

Sidecar

- 优势:
  - 查询时从多个源来获取数据，降低单一数据源压力
- 缺点:
  - 近期数据需要 Query 与 Sidecar 之间网络请求完成，会增加额外耗时
  - 需要 Store Gateway 能访问每个 Prometheus 实例

Receiver

- 优势：
  - 数据集中
  - Prometheus 无状态
  - 只需要暴露 Receiver 给 Prometheus 访问
- 缺点:
  - Receiver 承受大量 Prometheus 的 remote write 写入

总的来说

- 当 Prometheus 数量多的情况下，推荐使用 Sidecar 模式，因为所有 Prometheus 如果同时通过 Remote Write 向同一个 Receiver 写入数据，并且 Receiver 还要处理来自 Querier 的查询请求。
- 当需要 Prometheus 完全无状态的时候，推荐使用 Receiver，并且这种模式也利于将整体架构转移到其他监控系统(比如 Grafana 的 Mimir)，此时 Prometheus 仅作为发现目标、采集数据所用，将时序数据的存储能力交给其他更优质的产品。

# Thanos 主要组件概述

**Sidecar 和 Receiver 组件是最基本的 Thanos 功能。也就是说一般情况下，值部署一个 Sidecar 组件或者 Receiver 组件即可处理 Prometheus 本地的时序数据**

不同的是：

- Sidecar 是通过推送的方式，将 Prometheus 本地数据推送到对象存储。
  - 不过 Sidecar 还有一些其他功能，比如动态加载配置文件等等。
- Receiver 自身就是一个 TSDB，也实现了 Prometheus 的 Remote Write API。接收 Prometheus 发送过来的 WAL 数据存储到 Receiver 本地的 TSDB 中。
  - 当然，Receiver 也可以将本地时序数据推送到对象存储。

## [Sidecar(边车)](https://thanos.io/tip/components/sidecar.md/)

Sidecar 可以连接 Prometheus Server，读取其数据以便进行下一步处理。Sidecar 可以提供下列功能：

- **实时配置更新** # 实时监控 Prometheus Server 配置文件，当文件变化时让 Prometheus Server 重载配置。或者通过配置文件模板，生成配置文件。
- **转存数据 **# 将 Prometheus Server 的数据转存到其他地方。比如 对象存储、公有云存储、本地文件系统 等等
- **暴露 StoreAPI **# Querier 组件可以通过 Sidecar 的 StoreAPI 向 Prometheus Server 发起 PromQL 请求，以获查询数据。

注意：

- 该组件并不是通过 Prometheus 的 Remote Write 将数据写入到对象存储中的。而是将 Prometheus 已经压缩的数据块，上传到对象存储中。所以，如果直接从对象存储中查询数据，并不是实时的，远程存储的数据与 Prometheus 中的数据的时间间隔受 Prometheus 压缩数据间隔影响。
  - 所以，Thanos 如果想要通过 Remote Write 来转存数据，则是通过 **Receiver **功能实现的。
- 要使 Sidecar 模式正常运行，必须修改 Prometheus 的 `--storage.tsdb.min-block-duration` 与 `--storage.tsdb.max-block-duration` 这俩命令行标志的值为相同的值以禁用 Prometheus 的本地压缩，通常推荐设置为默认的 `2h`。
  - 这个设置是为了避免当 Sidecar 在上传数据会读取 Prometheus 的本地数据时而产生的问题。比如正在读取的数据正在被压缩，那么该数据上传到对象存储中将会出现问题，甚至都无法正常读取与上传

## [Receiver(接收器)](https://thanos.io/tip/components/receive.md/)

Receiver 是一个实现了 Prometheus 的 Remote Write API 的 TSDB，具有如下功能。

- **Remote Write API** # Receiver 接收 Prometheus Remote Write 传输过来的 WAL 数据，并存储在本地 TSDB 中。
- **转存数据 **# Receiver 还可以将这些数据转存到对象存储中。
- **暴露 StoreAPI **# Queriver 组件可以直接向 Receiver 的 StoreAPI 发起 PromQL 请求。
  - 注意：由于 Receiver 自身就是 TSDB，所以这个 PromQL 可以直接从 Receiver 中查询数据而不用像 Sidecar 似的，还得把查询请求转发给 Prometheus。

# Thanos 其他组件概述

## [Querier(查询器)](https://thanos.io/tip/components/query.md/)

**Query 组件分为两部分**

- **Querier(查询器) **# **实现了 Prometheus API**，可以通过 Querier 发起 PromQL 查询请求，以获取数据；甚至可以从 Prometheus Server 的时序数据库中删除数据。每个从 Querier 发起的 PromQL 查询请求都会发送到可以暴露 StoreAPI 的组件上，并获取查询结果。
- **Query Fronted(查询前端) **#** 实现 Prometheus API**，可以将请求负载均衡到指定的多个 Querier 上，同时可以缓存响应数据、也可以按查询日拆分。有点像 Redis 之于 Mysql 的效果

![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/gq2u99/1620707578416-24269f9e-5260-465a-ab5b-083980a3a016.png)

## [Store Gateway(存储网关)](https://thanos.io/tip/components/store.md/)

当 Sidcar 或 Receiver 将 Prometheus Server 中的数据转存到对象存储中后，下一步我们就会希望如何用起来这些数据，比如，我们可不可以直接使用 PromQL 查询对象存储中的数据呢？当然可以！这就是 Store Gateway 的工作。

**Store Gateway 也会暴露 StoreAPI 端点**，Querier 组件通过这个 API，使用 PromQL 从对象存储中查询数据。

注意：当 Query 组件通过 Store Gateway 查询数据时，由于对象存储中数据并不是实时的，所以，并不适合查询即时数据。

## [Compactor(压缩器)](https://thanos.io/tip/components/compact.md/)

Prometheus Server 会定期压缩旧数据以提高查询效率。Compactor 组件可以实现相同功能，只不过，压缩目标是 Sidecar 或 Receiver 转存到对象存储中的数据。

由于我们有数据长期存储的能力，也就可以实现查询较大时间范围的监控数据，当时间范围很大时，查询的数据量也会很大，这会导致查询速度非常慢。通常在查看较大时间范围的监控数据时，我们并不需要那么详细的数据，只需要看到大致就行。Thanos Compact 这个组件应运而生，它读取对象存储的数据，**对其进行压缩(下采样)再上传到对象存储**，这样在查询大时间范围数据时就可以只读取压缩和降采样后的数据，极大地减少了查询的数据量，从而加速查询。

## [Ruler(规则管理器)](https://thanos.io/tip/components/rule.md/)

由于 Ruler 组件获取评估数据的路径是 `ruler --> query --> sidecar --> prometheus`，需要经整个查询链条，这也提升了发生故障的风险，而且评估原本就可以在 Prometheus 中进行，所以在非必要的情况下更加推荐使用原本的 Prometheus 方式来做报警和记录规则的评估。

# Thanos 关联文件与配置

Thanos 的配置就是指各个组件自己的配置

# Thanos 监控面板

<https://github.com/thanos-io/thanos/blob/main/examples/dashboards/dashboards.md>
