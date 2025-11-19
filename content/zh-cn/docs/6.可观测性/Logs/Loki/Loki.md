---
title: Loki
linkTitle: Loki
weight: 1
---

# 概述

> 参考：
>
> - [GitHub 项目，grafana/loki](https://github.com/grafana/loki)
> - [官方文档，基础 - 概述](https://grafana.com/docs/loki/latest/fundamentals/overview/)

Loki 是受 Prometheus 启发的水平可扩展，高度可用的多租户日志聚合系统。它的设计具有很高的成本效益，并且易于操作。它不索引日志的内容，而是为每个日志流设置一组标签。
与其他日志聚合系统相比，Loki 有以下特点：

- 不对日志进行全文本索引。通过存储压缩的，非结构化的日志以及仅索引元数据，Loki 更加易于操作且运行成本更低。
- 使用与 Prometheus 相同的标签对日志流进行索引和分组，从而使您能够使用与 Prometheus 相同的标签在指标和日志之间无缝切换。
- 特别适合存储 Kubernetes Pod 日志。诸如 Pod 标签之类的元数据会自动被抓取并建立索引。
- 在 Grafana 中具有本机支持（需要 Grafana v6.0）。

基于 Loki 的日志包含 3 个程序：

- Loki 是主服务器，负责存储日志和处理查询。
- ClientAgent 客户端代理，负责收集日志并将其发送给 Loki。
    - [Promtail](/docs/6.可观测性/DataPipeline/Promtail/Promtail.md) 是其中一种 Agent，是 loki 原配。
- Grafana 用于查询和显示日志。

Loki 像 Prometheus 一样，但是是用于处理日志的：我们更喜欢基于多维标签的索引方法，并且希望使用没有依赖关系的单一二进制，易于操作的系统。Loki 与 Prometheus 的不同之处在于，它侧重于日志而不是指标，并通过推送而不是拉取交付日志。

> Tips: Loki 与 ClientAgent 加一起才相当于 Prometheus，因为 Promtail 是发现目标，采集日志的程序。然后主动 Push 给 Loki，由 Loki 存储日志数据。
> 而 Promtheus，可以自己发现目标，采集指标，存储指标。

## Loki Observability(可观察性)

> 参考：
>
> - [官方文档，运维 - 可观测性](https://grafana.com/docs/loki/latest/operations/observability/)

Loki 在 `/metrics` 端点上公开了指标，该端点暴露了 [OpenMetrics](/docs/6.可观测性/Metrics/监控系统概述/OpenMetrics.md) 格式的指标。

Loki 存储库具有一个[混合包](https://github.com/grafana/loki/tree/main/production/loki-mixin)，其中包括一组仪表板，记录规则和警报。总之，mixin 为您提供了一个全面的软件包，用于监视生产中的 Loki。

有关 mixin 的更多信息，请参阅 [monitoring-mixins 项目](https://github.com/monitoring-mixins/docs) 的文档 。

## Multi Tenancy(多租户)

Loki 支持多租户，以使租户之间的数据完全分离。当 Loki 在多租户模式下运行时，所有数据（包括内存和长期存储中的数据）都由租户 ID 分区，该租户 ID 是从请求中的 `X-Scope-OrgID` HTTP 头中提取的。当 Loki 不在多租户模式下时，将忽略 Header 头，并将租户 ID 设置为 `fake`，这将显示在索引和存储的块中。

# 学习资料

[公众号 - gopher云原生，浅析 Grafana Loki 日志聚合系统](https://mp.weixin.qq.com/s/kGY_zNhlXjnEAgqRUccMtg)

- https://mp.weixin.qq.com/s/6RC7RP2l5nqKNidsVVYoNA

# Loki 架构概述

> 参考：
>
> - [官方文档，基础 - 架构](https://grafana.com/docs/loki/latest/fundamentals/architecture/)

Loki 由多个组件组成，每个组件都可以实现特定的功能：

- **写入日志数据**
  - **Distributor(分配器)** # 对应 distributor 组件。负责处理客户端写入的日志，它是日志数据写入路径中的**第一站**，一旦 Distributor 收到日志数据，会将其拆分为多个批次，然后并行发送给一个或多个 Ingester
  - **Ingester(摄取器)** # 对应 ingester 组件。负责将日志数据写入 本地文件系统 或 指定的存储后端(DynamoDB、S3、Cassandra 等)
- **读取日志数据**，处理 LogQL 请求
  - **Querier(查询器)** # 对应 querier 组件。接收客户端发送的 LogQL 请求并从定的存储中查询日志数据并返回给客户端
  - **Query Frontend(查询前端)** # 对应 query-frontend 组件。为 Querier 组件提供负载均衡功能。
- **其他**
  - **Table Manager 表管理器)** # 对应 table-manager 组件。负责所有数据中，Table 的维护工作。根据配置文件中 schema_config.configs 字段中的相关配置，在指定时间开始之前创建周期表，并在根据 table_manager 字段中的相关配置，将数据时间范围超过保留期的数据删除。
  - **Compactor(压缩器)** # 2.6 版本时，Compactor 组件被设置为默认的用来实现数据保留功能的组件，暂时只支持 boltdb-shipper。准备要代替 table-manager 组件。
  - **Ruler(规则管理器)** # 对应 ruler 组件。从存储中读取数据，根据规则发送给告警处理程序。

loki 二进制文件的设计方式与 thanos 非常类似，都是在单一二进制文件中，可以运行指定的一个或多个组件。

Loki 内部将组件称为 **Modules(模块)**。如果想要运行指定的模块，有两种方式：

- 命令行标志 # loki 二进制文件的 `-target` 命令行标志
- 配置文件 # 配置文件中的 `target` 字段。

target 可用的值有：

- **all** # 表示 Loki 将以 Monolithic 架构运行，这也是默认的运行方式。Monolithic 模式非常适合测试或运行一个小规模的 Loki；而 Microservices 架构则提供了 Loki 的水平扩展性。
- **read** # 运行 Ingestor 和 Distributor 组件
- **write** # 运行 Querier、Query frontend、Ruler 组件
- **ingester** # 只运行 Ingester 组件
- **distributor** # 只运行 Distributor 组件
- **query-frontend** # 只运行 Query Frontend 组件
- **query-scheduler** # 只运行
- **querier** # 只运行 Querier 组件
- **index-gateway** # 只运行
- **ruler** # 只运行 Ruler 组件
- **compactor** # 只运行 Compactor 组件

## 最基本的运行条件

这些组件中，可以和存储直接交互的有 Ingester、Querier、Ruler。**最重要的组件是 Distributor、Ingester、Querier**这三个，这是 Loki 基本运行的最低要求。

Distributor 接收客户端(比如 Promtail) 推送的日志，处理后交给 Ingester 转存到本地或对象存储中，Querier 接收 LogQL 查询请求。

## 架构分类

> 参考：
>
> - [官方文档，基础 - 架构 - 部署模式](https://grafana.com/docs/loki/latest/fundamentals/architecture/deployment-modes/)

作为一个应用程序，Loki 由许多组件微服务构建而成，旨在作为一个可水平扩展的分布式系统运行。Loki 的独特设计将整个分布式系统的代码编译成单个二进制或 Docker 映像。该单个二进制文件的行为由-target 命令行标志控制，并定义了三种操作模式之一。

Loki 旨在根据需求变化轻松地在不同架构下重新部署集群，无需更改配置或进行最少的配置更改。

- Monolithic 架构对于快速开始试验 Loki 以及每天高达约 100GB 的小读/写量非常有用。
- Loki 的简单可扩展部署可以扩展到每天数 TB 甚至更多的日志。
- 对于非常大的 Loki 集群或需要对扩展和集群操作进行更多控制的集群，建议使用微服务模式。

### Monolithic(统一) 架构

这种架构需要通过 loki 二进制文件只启动 1 个进程，使该进程用 `-target=all` 以便在一个进程中运行所有组件。

![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/gcp6zx/1660115619129-e6fa2017-8e05-46d7-ab56-207ee3cfc90b.png)

这是最经典的模式，早期 Loki 通常都是以这种模式被大家部署。这种模式是 loki 以单个二进制文件运行 Loki 的所有组件，如上图，instance 可以看作一个单独的二进制文件。
Monolithic 模式非常适合于本地开发、小规模等场景，Monolithic 模式可以通过多个进程进行扩展，但有以下限制：

- 当运行带有多个副本的单体模式时，当前无法使用本地索引和本地存储，因为每个副本必须能够访问相同的存储后端，但是本地存储对于并发访问并不安全。主要是因为 BoltDB 仅允许一个进程在同一时间锁定数据库。如果使用远程存储不受影响。
- 各个组件无法独立缩放，因此读取组件的数量不能超过写入组件的数量。

这个进程产生一个 gRPC 监听(默认 9095 端口)和一个 HTTP 监听(默认 3100 端口)。各个组件内部在同一个进程的共享内存中进行数据交互。

### Simple scalable(简单可扩展) 架构

这种架构需要通过 loki 二进制文件至少启动 2 个进程，保证两个进程分别具有 读 和 写 的功能

- 其中一个进程使用 `-target=write` 运行具有写功能的组件，包括 Ingestor 和 Distributor
- 另一个进程使用 `-target=read` 运行具有读功能的组件，包括 Querier、Query frontend、Ruler

![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/gcp6zx/1660120707493-6efe2870-f1a3-446f-b760-6f520236c358.png)

这种架构将 Loki 的读/写分离。这个图里少了一点，通常来说，5 个 Loki 实例前面还有一个负载均衡设备，用来接收客户端的 读/写请求，以便将请求转发给对应的 Loki 实例。

### Microservices(微服务)  架构

这种架构需要通过 loki 二进制文件至少启动 4 个进程，整套架构由多个单一功能的进程组成

- `-target=distributor` # 运行分配器
- `-target=ingester`# 运行摄取器
- `-target=querier`# 运行查询器
- `-target=query-frontend`# 运行查询前端

![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/gcp6zx/1660115629058-db37e36d-3ed5-4bd9-86bc-3fbb05df38d0.png)

这种微服务架构与 Thanos 类似，可以通过一个 Loki 的二进制文件，使用子命令来启动不同的功能。

- 每个组件都产生一个 gRPC 监听(默认 9095 端口)和一个 HTTP 监听(默认 3100 端口)。
  - 通常情况下，gRPC 端口用于组件间通信；HTTP 端口用于暴露一些管理 API(比如 指标、运行状态、就绪性)
- 各个组件可以暴露的 HTTP API 详见 [Loki API](/docs/6.可观测性/Logs/Loki/Loki%20API.md) 笔记。通过 API，我们可以更清晰得了解到，每个组件可以实现的具体功能
- 各个组件通过 memberlist 统一到一个哈希环上，以互相发现。当我们部署在 K8S 中时，将会配置 `memberlist.join_members` 字段，并且需要创建对应的 service 资源，service 的 endpoint 将会关联到所有 Distributor、Ingester、Querier 组件。

### Gateway(网关)

在我们使用 Simple scalable 和 Microservices 这两种架构时，通常会使用一个 `loki-gateway` ，这是一个 Nginx，配置很简单：

```nginx
http {
  server {
    listen             8080;

    location = / {
      return 200 'OK';
      auth_basic off;
    }

    ......略

    location = /loki/api/v1/push {
      set $loki_api_v1_push_backend http://loki-loki-distributed-distributor.logging.svc.cluster.local;
      proxy_pass       $loki_api_v1_push_backend:3100$request_uri;
    }

    location ~ /loki/api/.* {
      set $loki_api_backend http://loki-loki-distributed-query-frontend.logging.svc.cluster.local;
      proxy_pass       $loki_api_backend:3100$request_uri;
    }
  }
}
```

可以看到，这个 `loki-gateway` 用来为 Loki 进行读/写分离的。loki-gateway 会根据客户端发起请求的 URL 判断这个请求应该由哪个组件进行处理。

Nginx 的配置依据两种架构的不同而有细微区别，但是总归是需要一个 Gateway 的。不管是 Promtail 推送数据，还是客户端查询数据，都可以先经过 loki-gateway

## 数据写入路径

![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/gcp6zx/1660123438275-02b6febb-0f26-431b-9450-9b5f6f125305.png)

整体的日志写入路径如下所示：

- `Distributor` 收到一个 HTTP 请求，以存储流的数据。
- 每个流都使用哈希环进行哈希操作。
- `Distributor` 将每个流发送到合适的 `Ingester` 和他们的副本（基于配置的复制因子）。
- 每个 `Ingester` 将为日志流数据创建一个块或附加到一个现有的块上。每个租户和每个标签集的块是唯一的。
- `Distributor` 通过 HTTP 连接响应一个成功代码。

## 数据读取路径

![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/gcp6zx/1660118903925-bca6ba6b-f991-4a28-a407-9c6febb38a36.png)

日志读取路径的流程如下所示：

- `Querier` 收到一个对数据的 HTTP 请求。
- `Querier` 将查询传递给所有 `Ingesters` 以获取内存数据。
- `Ingesters` 收到读取请求，并返回与查询相匹配的数据（如果有的话）。
- 如果没有 `Ingesters` 返回数据，查询器会从后端存储(比如 S3)加载数据，并对其运行查询。
- 查询器对所有收到的数据进行迭代和重复计算，通过 HTTP 连接返回最后一组数据。

# Loki 主要组件概述

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/gcp6zx/1621238613211-cedcd7da-602a-4c15-9b27-dbcb797317d8.png)

## Distributor(分配器)

Distributor 服务负责处理客户端写入的日志，它本质上是日志数据写入路径中的**第一站**，一旦 Distributor 收到日志数据，会将其拆分为多个批次，然后并行发送给多个 Ingester。

Distributor 通过 gRPC 与 Ingester 通信，它们都是无状态的，可以根据需要扩大或缩小规模。

## Ingester(摄取器)

Ingester 服务负责将日志数据写入长期存储后端（DynamoDB、S3、Cassandra 等）。此外 Ingester 会验证摄取的日志行是按照时间戳递增的顺序接收的（即每条日志的时间戳都比前面的日志晚一些），当 Ingester 收到不符合这个顺序的日志时，该日志行会被拒绝并返回一个错误。

注意：虽然 Ingester 支持 BoltDB 写入本地文件系统，但是这仅适用于[单进程模式](/docs/6.可观测性/日志系统/Loki/Loki%20 部署.md 部署.md)，因为 Querier 也需要访问相同的存储，而 BoltDB 仅允许一个进程在同一时间锁定数据库。

## Querier(查询器)

Querier(查询器) 使用 LogQL 处理查询，从 Ingesters 和长期存储中获取日志。

## Query Frontend(查询前端)

Query Frontend(查询前端) 是一个可选的组件。当 Loki 以微服务架构运行时，且存在多个 Querier(查询器)，则查询前端可以平均得调度 LogQL 请求到查询器上，说白了就是实现负载均衡的效果。并且查询前端还可以并行处理请求、并缓存这些数据。

# Loki 关联文件与配置

**/etc/loki/local-config.yaml** # loki 程序运行时默认配置文件

**${StorageConfig}/index** # loki 的 BoltDB 中存储索引数据保存路径，无默认值，根据配置文件中 `.strorage_confg.boltdb.directory` 字段指定。

**${StorageConfig}/chunks** # loki 的 chunks(块) 存储数据保存路径，无默认值，根据配置文件中 `.strorage_confg.filesystem.directory` 字段指定。

# Loki 与其他日志系统相比

官方文档：<https://grafana.com/docs/loki/latest/fundamentals/overview/comparisons/>

Loki / Promtail / Grafana vs EFK

EFK（Elasticsearch，Fluentd，Kibana）堆栈用于从各种来源提取，可视化和查询日志。

Elasticsearch 中的数据作为非结构化 JSON 对象存储在磁盘上。每个对象的键和每个键的内容都被索引。然后可以使用 JSON 对象或定义为 Lucene 的查询语言来查询数据以定义查询（称为查询 DSL）。

相比之下，Loki 在单二进制模式下可以将数据存储在磁盘上，但是在水平可伸缩模式下，数据存储在诸如 S3，GCS 或 Cassandra 之类的云存储系统中。日志以纯文本格式存储，并带有一组标签名称和值，其中仅对标签对进行索引。这种折衷使得它比全索引更便宜，并且允许开发人员从其应用程序积极地进行日志记录。使用 LogQL 查询 Loki 中的日志。但是，由于这种设计上的折衷，基于内容（即日志行中的文本）进行过滤的 LogQL 查询需要加载搜索窗口中与查询中定义的标签匹配的所有块。

Fluentd 通常用于收集日志并将其转发到 Elasticsearch。Fluentd 被称为数据收集器，它可以从许多来源提取日志，对其进行处理，然后将其转发到一个或多个目标。

相比之下，Promtail 的用例专门针对 Loki 量身定制。它的主要操作模式是发现存储在磁盘上的日志文件，并将与一组标签关联的日志文件转发给 Loki。Promtail 可以为与 Promtail 在同一节点上运行的 Kubernetes Pod 进行服务发现，充当容器边车或 Docker 日志记录驱动程序，从指定的文件夹中读取日志并尾随系统日志。

Loki 用一组标签对表示日志的方式类似于 Prometheus 表示度量的方式。当与 Prometheus 一起部署在环境中时，由于使用相同的服务发现机制，Promtail 的日志通常具有与应用程序指标相同的标签。具有相同级别的日志和指标使用户可以在指标和日志之间无缝地进行上下文切换，从而有助于根本原因分析。

Kibana 用于可视化和搜索 Elasticsearch 数据，并且在对该数据进行分析时非常强大。Kibana 提供了许多可视化工具来进行数据分析，例如位置图，用于异常检测的机器学习以及用于发现数据关系的图形。可以将警报配置为在发生意外情况时通知用户。

相比之下，Grafana 专门针对来自 Prometheus 和 Loki 等来源的时间序列数据量身定制。可以设置仪表板以可视化指标（即将提供日志支持），并且可以使用浏览视图对数据进行临时查询。与 Kibana 一样，Grafana 支持根据您的指标进行警报。

- kibana 启动速度比 grafana 慢了 10 倍
- es 启动时，内存使用达到 1.5G，后续存储同样内容的情况下，es 内存使用率 1G 多，loki 内存使用率 200 多 M
- promtail 使用 yaml 作为 配置文件格式，与 prom 配置逻辑一致。fluentd 配置文件格式类似 html
- grafana 页面可以直接通过标签用鼠标点击过滤。kibana 则需要输入内容。
