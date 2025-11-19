---
title: Storage(存储)
linkTitle: Storage(存储)
weight: 1
---

# 概述

> 参考：
>
> - [官方文档，存储](https://grafana.com/docs/loki/latest/storage/)
> - [官方文档，运维 - 存储](https://grafana.com/docs/loki/latest/operations/storage/)
> - [博客 - 2023-12-20，Grafana Loki 简明指南：关于标签您需要了解的一切](https://grafana.com/blog/2023/12/20/the-concise-guide-to-grafana-loki-everything-you-need-to-know-about-labels/)

与其他日志记录系统不同，Loki 是基于仅索引日志的元数据的想法而构建的。从  Loki 的 [Data Model(数据模型)](/docs/6.可观测性/Logs/Loki/Storage(存储)/Data%20Model(数据模型).md) 可知，日志是根据标签进行定位的。 日志数据本身会被压缩成 Chunks，并存储在本地的文件系统中；并且 Loki 还提供了一个 Index 数据，用来根据索引定位日志数据。小索引和高度压缩的 Chunks 简化了操作，并显着降低了 Loki 的成本。

所以 Loki 需要存储两种不同类型的数据，当 Loki 收到 Log Stream 时，会存储两类数据：

- **Chunk(块)** # **日志流本身的信息**。每一个 Chunk 都是将一段时间的日志流压缩后形成的一个文件。
  - 一个 Chunks 就是一个对象，如果是使用本地文件系统存储 Chunks，则可以抽象得将一个 Chunks 文件当做一个对象。在一个 Chunks 文件中一般包含里一段时间的日志流数据。
- **Index(索引)** # **日志流索引的信息**。每一个 Index 都是 键/值 格式的数据库文件，文件中的内容用来关联 日志流的标签 与 Chunks。
  - Index 中的 Key 就是日志流的标签，Value 就是 Chunks 文件所在的绝对路径。

> [!Tip] History
> Loki 在 2.0 版本之前，这两类数据是分开存放的，只有 Chunk 数据可以存在对象存储中。
>
> 直到 2.0 发布，Loki 开发了基于 BoltDB 的 BoltDB-Shipper 数据库用来存储 Index，并且已经可以将 Index 数据也存到对象存储中。也就是 **Single stroe(单存储) 架构**，此时 Chunk 和 Index 都可以同时存在本地文件系统或者同时存在对象存储中。
>
> 从 2.8 版本开始，Loki 使用 TSDB-Shipper 存储代替了 BoltDB-Shipper 来存储 Index。这种方式极大得优化了索引的体积。
>
> 逐渐的，可能也就不再从外部主动区分 Chunk 和 Index，而是将这两个抽象的概念统一成 Loki 存储即可。然后 Loki 存储的内部实现具体用不用 Index，或者 Index 的实现方式是什么样的就不用特殊关注了
>
> ![500](https://notes-learning.oss-cn-beijing.aliyuncs.com/observability/loki/storage/20250326134111466.png)

Loki 可以将 Index 及其 Chunk 数据保存在多种 Backend(后端)：

- **Local Storage(本地存储)** # 通常都是指本地文件系统。这是最简单的 Backend。尽管它也容易因未复制而导致数据丢失，这对于单一二进制使用 Loki 或在项目上进行本地开发的用户来说是很常见的。它在概念上类似于 Prometheus
- **Remote Storage(远程存储)** # 远程存储通常都是指 [对象存储](/docs/5.数据存储/存储/存储的能力/对象存储/对象存储.md)。

Loki 在不同的 Log Stream(日志流) 中接收日志，其中每个 Stream 的 tenantID(租户 ID) 和 一组标签 是该 Stream 的唯一标识。如果 Loki 以单租户模式运行，则所有块都放在名为 **`fake/`** 的目录中，这是用于单个租户模式的合成租户名称。

## Local Storage(本地存储)

与 [Prometheus 的存储概念](/docs/6.可观测性/Metrics/Prometheus/Storage/Storage.md)类似，Loki 也是将日志流数据抽象为一个一个的 Block(块)。由于 Loki 需要存储 Index 与 Chunks 两种数据，所以，数据在磁盘上的布局，与 Prometheus 也就不太一样了。

```bash
~]# ll
total 0
drwxr-xr-x 4 10001 10001  77 Oct 23 17:50 chunks
drwxr-xr-x 4 10001 10001  51 Mar 26 09:02 compactor
drwxr-xr-x 2 10001 10001  10 Mar  7 13:20 rules
drwxr-xr-x 7 10001 10001 109 Mar  7 14:02 tsdb-shipper-active
drwxr-xr-x 5 10001 10001  79 Mar 26 09:02 tsdb-shipper-cache
drwxr-xr-x 4 10001 10001 112 Mar 26 09:07 wal
```

- **chunks** # 持久化保存 Chunks 与 Index 数据。
- **compactor** # Compactor 组件的工作路径。由 compactor.working_directory 配置设置
- **rules**
- **tsdb-shipper-active** # TODO: 好像是 Index 数据临时中专的地方，最后会移动到 chunks/index/ 目录下
- **tsdb-shipper-cache** # 查询时的缓存。若没有查询则该目录不会出现内容
- **wal** # 顾名思义，保存 WAL 数据的。

### 本地 Index

目录组织结构如下：

```bash
~]# tree tsdb-shipper-active/multitenant/
tsdb-shipper-active/multitenant/
└── index_20173
    └── 1742967017-localhost.localdomain-1729669818509348421.tsdb

# 下面是曾经 boltdb 时代记录的
~]# tree boltdb-shipper-active/index_18766/
boltdb-shipper-active/index_18766/
├── 1621413900
├── 1621413900.snapshot
├── 1621414800
└── 1621414800.snapshot
0 directories, 4 files

~]# tree boltdb-shipper-active/index_18766/
boltdb-shipper-active/index_18766/
├── 1621414800
└── 1621414800.snapshot
0 directories, 2 files

~]# tree boltdb-shipper-active/index_18766/
boltdb-shipper-active/index_18766/
├── 1621414800
└── 1621414800.snapshot
├── 1621415700
└── 1621415700.snapshot
0 directories, 4 files
```

Loki 将 Index 数据根据 `schema_config.configs.index.period` 配置，按照时间时间期间进行分组，默认为 168h，也就是 7 天。

如果将 `schema_config.configs.index.period` 设置为 24h，那么对于 BoltDB-Shipper，一个表就是一个目录，表由许多较小的 BoltDB 文件组成(也就是说目录中有很多 BoltDB 文件)，每个文件仅存储 15 分钟的索引值。每天创建的表名称由 **prefix + period-number-since-epoch** 组成。

- **prefix** # 是 `schema_config.configs.index.period` 配置的值。
- **period-number-since-epoch**# 是从 [Epoch 时间](https://en.wikipedia.org/wiki/Unix_time)以来，到开始存储数据时刻的天数。

假如 `schema_config.configs.index.period` 的值为 `loki_index_`，当前时间是 2021 年 5 月 19 日。此时 Loki 收到了日志数据，则 Ingester 会创建一个 名为 **loki_index_18766** 的目录，并且当收到日志流时，会在这个目录中创建以时间戳为名字的文件，这个文件中就是 Index 数据。

> 自 Epoch 时间(1970 年 1 月 1 日) 至今为止，已经过了 18765 天，而我此时处于第 18766 天。那么这个目录名字中的数字就是 18766

上述 Index 数据将会根据 `storage_config.boltdb_shipper` 中的 `active_index_directory` 配置的存储路径存储在本地文件系统的目录内，如果还配置了 `shared_store` 的值为非 filesystem，那么还会将这些数据上传到指定的对象存储中。

在这里，我们为什么不会看到每隔 15 分钟就产生一个文件呢，那是因为 Compactor 这个组件的作用。Compactor 会将每个表中的 Index 数据进行合并，并删除其中的重复数据。你想啊~每隔 15 分钟产生一个文件，那一天 24 小时就有 96 个文件，积累起来是非常多的，在使用查询时也不够方便。

所以，每隔 15 分钟，Compactor 会将旧的 Index 数据压缩到新的 Index 文件中，就像上面的目录结构中展现的，在 1621413900 和 1621414800 存在两个 Index，Compactor 会将 1621413900 中的数据合并到 1621414800 中，并去重，然后留下了唯一一个文件，当下个 15 分钟，又会生成名为 1621415700 的文件，Compactor 又会重复之前的动作，将数据合并到 1621415700 中，最终，在一天结束时，只会留下唯一一个文件。然后，今天的文件，又会被合并到明天的 Index 文件中，最终的最终，只会留下唯一一个 Index 文件。

> Tips: Loki 使用 BoltDb 时，与 Etcd 一样，Index 到最后也会只有一个文件，这个文件就跟 etcd 的 db 文件一样，这个文件中是被压缩的很多很多的键/值对信息

每天的结束，会将这个文件压缩为 .gz 的格式，并发送到 Chunks 数据存储目录中的 `index/` 目录下。

### 本地 Chunks

Chunks 数据的目录结构就很简单了

```bash
~]# tree chunks/
chunks/
├── fake
│   ├── 123e004b6a42face
│   │   ├── MTk1M2MyNDE5YjY6MTk1M2MyNDIyMjg6YThiNmRmZjU=
│   │   ├── ......
│   │   └── MTk1YWMzZmVhYTQ6MTk1YWMzZmVlNTc6NDc0ZGIwMTk=
│   ├── ......
│   └── feb6db45141b0877
├── index
│   ├── index_20170
│   │   └── fake
│   │       └── 1742781768295831640-compactor-1742680944642-1742781341986-644ffb8.tsdb.gz
│   ├── ......
│   └── index_20173
└── loki_cluster_seed.json
```

日志数据集中在 `fake/` 目录下，每个 `MTk1M2MyNDE5YjY6MTk1M2MyNDIyMjg6YThiNmRmZjU=` 这类的字符文件，就是一个 Chunk。并且，Index 的数据压缩后，也会在 Chunks 目录中生成对应的文件。

## Remote Storage(远程存储)

TODO: 需要更新该部分内容

远程存储非常简单，从本地存储的模式可以看出来，Index 和 Chunks 本质上就是一个个的文件，并且互相具有关联关系，所以，Ingester 组件还可以将这些数据，发送到远程存储中。

> Loki 的远程存储与 Prometheus 还有一点区别，Loki 自己的 Ingester 本身就实现了，而 Prometheus 的远程存储则需要其他程序对接。

### Index

数据以 WAL 方式存在本地后，逐步上传到远程存储中。

```bash
├── data
│   ├── loki
│   │   ├── boltdb-shipper-active
│   │   │   └── uploader
│   │   │       └── name
│   │   ├── boltdb-shipper-cache
│   │   ├── compactor
│   │   └── wal
│   │       └── 00000000

```

### Chunks

# Schema(模式) 概念

> 参考：
>
> - [官方文档，运维 - 存储 - 存储模式](https://grafana.com/docs/loki/latest/operations/storage/schema/)

Loki 旨在向后兼容，当 Loki 内部存储发生变化时，通过 **Schema(模式)** 功能，可以让 Loki 的数据迁移更加平滑。在 Schema 概念中，通过一种 **Period(期间)** 的概念，来区分多个 Schema 的配置。本质上，一个 Schema 是一个数组，数组中的每个元素都是一个 Period，表示在这个 Period(期间) 内所使用的存储模式是 XX。

![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/gzp72g/1660102594593-ad7383b9-e99d-414b-b4e4-8547ef28758f.png)

同时，**Schema 中的配置，也可以定义 Loki 储存数据所用的存储类型，Loki 想要正常运行，必须要指定具体的 Schema**。

假如现在配置文件中有如下配置：

```yaml
schema_config:
  configs:
    - from: 2019-07-01
      store: boltdb-shipper
      object_store: filesystem
      schema: v10
      index:
        prefix: index_
        period: 168h
    - from: 2020-07-01
      store: boltdb-shipper
      object_store: s3
      schema: v11
      index:
        prefix: index_
        period: 168h
```

这个配置表示 Loki 现在可以通过两种 Schema 存储数据，在第一个 Period(期间)(2019 年 7 月 1 日 ~ 2020 年 7 月 1 号)，使用 v10 模式；在第二个 Period(期间)(2020 年 7 月 1 日 ~ 至今)，则使用 v11 模式。这种运行方式，显著简化了升级，让用户不再担心升级时如果存储模式发生变化而产生的影响。就算不是升级 Loki，当用户想要切换所使用的存储时，同样可以通过这种方式来平滑过度。

基于此，Loki 可以根据不同的 Schema 配置，将不同时间段的数据存放在不同类型的存储中。

- 比如现在有这么一种场景，公司本来将数据存放在本地，但是随着业务规模扩大，数据存储方式想改变了，这时候如何丝滑的迁移数据呢？
  - 通过 schema_config 字段的配置，2019-10-16 到 2020-11-16，使用本地存放数据；然后我买了公有云的存储产品，从 2020-11-16 到 2021-11-16，使用公有云存放数据。Loki 根据这种配置直到，当 2020-11-16 日开始，就会将数据存放在公有云数据库中了，而老的数据还可以不受影响。然后根据下文提到的 [Retention](#Log%20Retention(日志保留)) 功能，逐步将老数据删除，这种升级过度是非常丝滑的。
- 甚至当 Loki 的存储能力随着开发迭代有新的方式时，也可以通过配置让 Loki 升级更加无缝，升级后从某个时间段开始，使用新的存储方式。

## 在 Schema 中配置储存数据所使用的存储类型

并且，在 Schema 配置中，还可以通过 `schema_config.configs.store` 与 `schema_config.configs.object_store` 两个字段指定 Loki 储存 Index 和 Chunk 数据所使用的存储类型，比如上面这个例子中，对于 2020 年 7 月 1 日之前的数据，使用 boltdb-shipper 存储 Index 数据、filesystem 存储 Chunk 数据；2020 年 7 月 1 日之后的存储方式与之前一样。然后再通过 `storage_config` 字段来配置 boltdb-shipper 和 filesystem 这两类存储的具体使用方式，比如指定认证信息、指定存储路径等等。

同时，还可以通过`schema_config.configs.store` 与 `schema_config.configs.object_store` 两个字段来配置存储数据时 **Table(表)** 的行为，比如 存储时间、文件名前缀 等等。

比如上面这个示例中，配置了储存 Index 类型数据的一些简单行为：`prefix` 字段指定了当前期间表的前缀为 index\_，也就是存储 Index 数据的文件名开头都会是 index\_；`period` 字段则指定了在这个模式期间，表的周期为 168h，也就是每隔 168 小时，创建一张表。

## 依赖于 Schema 的配置

由于 Schema 中配置了 Loki 储存 Index 与 Chunks 两种数据所使用的存储类型。所以，很多配置，都会基于这个信息，来进行个性化配置。比如在 Schema 中设定了 S3 类型的存储，那么就需要设定连接 S3 存储时所需的认证信息。

> 注意：这一段内容会基于对 Loki 的存储概念已经有了详细了解后，才能看懂。

# Single Store

## TSDB

https://lokidex.com/posts/tsdb/

从 Loki v2.8 开始，TSDB 是推荐的 Loki 索引。它很大程度上受到 Prometheus 的 TSDB 子项目的启发。Loki 维护者 [Owen 的博客文章](https://lokidex.com/posts/tsdb/)中，有更多介绍。简而言之，这个新索引更高效、更快且更具可扩展性。它也驻留在对象存储中，就像它前面的 BoltDB-Shipper 索引一样。

TSDB 索引存储的工作原理可以分为 uploadsManager 和 downloadsManager 两部分：

- uploadsManager # 负责上传 `active_index_directory`(默认 tsdb-shipper-active/ 目录) 内的索引分片到配置的共享存储中，同时负责定期清理工作
- downloadsManager # 负责从共享存储下载索引到本地缓存目录 `cache_location`(默认 `tsdb-shipper-cache/` 目录)，同时负责定期同步和清理工作

# Log Retention(日志保留)

https://grafana.com/docs/loki/latest/operations/storage/retention/

> [!tip] 在 Loki 早期（版本 <= 2.7），Log Retention 能力是通过 Table manager 实现。从 2.8 版本开始，由于使用了 TSDB 存储模式，Retention 的行为可以由 Compactor 组件实现。

---

# OTel 支持

在 Loki 3.0 之前，是不支持将 OTel 日志直接导入到 Loki 的，所以创建了

