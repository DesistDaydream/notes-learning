---
title: Storage(存储)
---

# 概述

> 参考：
>
> - [官方文档,存储](https://grafana.com/docs/loki/latest/storage/)
> - [官方文档,运维-存储](https://grafana.com/docs/loki/latest/operations/storage/)
> - [官方文档,运维-存储-BoltDB-Shipper](https://grafana.com/docs/loki/latest/operations/storage/boltdb-shipper/)

与其他日志记录系统不同，Loki 是基于仅索引日志的元数据的想法而构建的。从 [Loki 的数据模型](/docs/6.可观测性/日志系统/Loki/Storage(存储)/Data%20Model(数据模型).md Model(数据模型).md)可知，日志是根据标签进行定位的。 日志数据本身会被压缩成 Chunks，并存储在本地的文件系统中；并且 Loki 还提供了一个 Index 数据，用来根据索引定位日志数据。小索引和高度压缩的 Chunks 简化了操作，并显着降低了 Loki 的成本。

所以 Loki 需要存储两种不同类型的数据，当 Loki 收到 Log Stream 时，会存储两类数据：

- **Chunks(块)** # **日志流本身的信息**。每一个 Chunks 都是将一段时间的日志流压缩后形成的一个文件。
  - 一个 Chunks 就是一个对象，如果是使用本地文件系统存储 Chunks，则可以抽象得将一个 Chunks 文件当做一个对象。在一个 Chunks 文件中一般包含里一段时间的日志流数据。
- **Indexes(索引)** # **日志流索引的信息**。每一个 Index 都是 键/值 格式的数据库文件，文件中的内容用来关联 日志流的标签 与 Chunks。
  - Index 中的 Key 就是日志流的标签，Value 就是 Chunks 文件所在的绝对路径。

**Chunks 与 Index 的存储方式**

- **Chunks** # 直接以压缩格式的文件形式存储。
- **Index** # Loki 自身实现了一个本地的 **Key/Value Database(键值数据库)**，这个数据库是基于 BoltDB 开发的，称为 **BoltDB-Shipper**。BoltDB-Shipper 用来存储 Index 数据。

> 小记：Loki 在 2.0 版本之前，这两类数据是分开存放的，只有 Chunks 数据可以存在对象存储中。直到 2.0 发布，Loki 开发了基于 BoltDB 的 BoltDB-Shipper 数据库用来存储 Index，并且已经可以将 Index 数据也存到对象存储中。

同时，Loki 还可以将这些数据，同时存储到 远程存储 中去(比如对象存储)。这些功能都是通过 Ingester 组件实现的。

Loki 在不同的 Log Stream(日志流) 中接收日志，其中每个 Stream 的 tenantID(租户 ID) 和 一组标签 是该 Stream 的唯一标识。如果 Loki 以单租户模式运行，则所有块都放在名为 `**fake**` 的文件夹中，这是用于单个租户模式的合成租户名称。

## Local Storage(本地存储)

### On-disk Layout(磁盘上的布局)

与 [Prometheus 的存储概念](/docs/6.可观测性/监控系统/Prometheus/Storage(存储)/Storage(存储).md)类似，Loki 也是将日志流数据抽象为一个一个的 Block(块)，只不过，在 Loki 这里，称之为 [Table(表)](#Table(表)%20概念)。由于 Loki 需要存储 Index 与 Chunks 两种数据，所以，数据在磁盘上的布局，与 Prometheus 也就不太一样了。

### Index

目录组织结构如下：

```bash
[root@nfs-1 loki]# tree boltdb-shipper-active/index_18766/
boltdb-shipper-active/index_18766/
├── 1621413900
├── 1621413900.snapshot
├── 1621414800
└── 1621414800.snapshot
0 directories, 4 files

[root@nfs-1 loki]# tree boltdb-shipper-active/index_18766/
boltdb-shipper-active/index_18766/
├── 1621414800
└── 1621414800.snapshot
0 directories, 2 files

[root@nfs-1 loki]# tree boltdb-shipper-active/index_18766/
boltdb-shipper-active/index_18766/
├── 1621414800
└── 1621414800.snapshot
├── 1621415700
└── 1621415700.snapshot
0 directories, 4 files
```

Loki 将 Index 数据根据 `schema_config.configs.index.period` 这个配置，按照时间时间期间进行分组，莫认为 168h，也就是 7 天。

如果将 `schema_config.configs.index.period` 设置为 24h，那么对于 BoltDB-Shipper，一个表就是一个目录，表由许多较小的 BoltDB 文件组成(也就是说目录中有很多 BoltDB 文件)，每个文件仅存储 15 分钟的索引值。每天创建的表名称由 **prefix + period-number-since-epoch** 组成。

- **prefix** # 是 `schema_config.configs.index.period` 配置的值。
- **period-number-since-epoch**# 是从 [Epoch 时间](https://en.wikipedia.org/wiki/Unix_time)以来，到开始存储数据时刻的天数。

假如 `schema_config.configs.index.period` 的值为 `loki_index_`，当前时间是 2021 年 5 月 19 日，此时 Loki 收到了日志数据，则 Ingester 会创建一个 名为 **loki_index_18766** 的目录，并且当收到日志流时，会在这个目录中创建以时间戳为名字的文件，这个文件中就是 Index 数据。

> 自 Epoch 时间(1970 年 1 月 1 日) 至今为止，已经过了 18765 天，而我此时处于第 18766 天。那么这个目录名字中的数字就是 18766

上述 Index 数据将会根据 `storage_config.boltdb_shipper` 中的 `active_index_directory` 配置的存储路径存储在本地文件系统的目录内，如果还配置了 `shared_store` 的值为非 filesystem，那么还会将这些数据上传到指定的对象存储中。

在这里，我们为什么不会看到每隔 15 分钟就产生一个文件呢，那是因为 Compactor 这个组件的作用。Compactor 会将每个表中的 Index 数据进行合并，并删除其中的重复数据。你想啊~每隔 15 分钟产生一个文件，那一天 24 小时就有 96 个文件，积累起来是非常多的，在使用查询时也不够方便。

所以，每隔 15 分钟，Compactor 会将旧的 Index 数据压缩到新的 Index 文件中，就像上面的目录结构中展现的，在 1621413900 和 1621414800 存在两个 Index，Compactor 会将 1621413900 中的数据合并到 1621414800 中，并去重，然后留下了唯一一个文件，当下个 15 分钟，又会生成名为 1621415700 的文件，Compactor 又会重复之前的动作，将数据合并到 1621415700 中，最终，在一天结束时，只会留下唯一一个文件。然后，今天的文件，又会被合并到明天的 Index 文件中，最终的最终，只会留下唯一一个 Index 文件。

特别说明：Loki 与 Etcd 都是使用的 BoltDB 实现的，Index 到最后也会只有一个文件，这个文件就跟 etcd 的 db 文件一样，这个文件中是被压缩的很多很多的键/值对信息

每天的结束，会将这个文件压缩为 .gz 的格式，并发送到 Chunks 数据存储目录中的 index 目录下。

### Chunks

Chunks 数据的目录结构就很简单了

```bash
[root@nfs-1 chunks]# tree |more
.
├── index
│   ├── index_18738
│   │   └── compactor-1619056534.gz
│   ├── index_18739
│   │   └── compactor-1619142934.gz
......省略
│   └── index_18766
│       ├── compactor-1621432533.gz
│       └── loki-bj-net-0-1615528533442739252-1621431900.gz
├── ZmFrZS80MDgzM2JkNDYxNWM0MWIzOjE3OGZmOTgxM2NmOjE3OGZmOTgxM2QwOjUyN2VlY2Yy
├── ZmFrZS80MDgzM2JkNDYxNWM0MWIzOjE3OTdhM2JlOTUxOjE3OTdhM2JlOTUyOjI3NDZlMDdj
├── ZmFrZS80MDIwYTRmN2I2MzVmZGIyOjE3OGY3MDYwZTIwOjE3OGY3MDYxZTE3OmFhMzZhNzg=
├── ZmFrZS80MDIwYTRmN2I2MzVmZGIyOjE3OGY5NzU4MzI3OjE3OGY5NzVkZTNiOmIxOTIzNWY4
├── ZmFrZS80MDIwYTRmN2I2MzVmZGIyOjE3OGZjYWZkZTA1OjE3OGZjYWZkZjE1OjIzNWIzMjVi
.......省略
```

那些一长串的字符文件，就是一个 Chunks。并且，Index 的数据压缩后，也会在 Chunks 目录中生成对应的文件。

## Remote Storage(远程存储)

远程存储非常简单，从本地存储的模式可以看出来，Index 和 Chunks 本质上就是一个个的文件，并且互相具有关联关系，所以，Ingester 组件还可以将这些数据，发送到远程存储中。

现阶段：

- Chunks 支持以下远程存储
  - **Cassandra**
  - **S3** # 任何实现 S3 接口的服务都可以用来存储 Chunks 数据，比如开源的 [MinIO](/docs/5.数据存储/1.存储/存储的基础设施架构/Distributed%20Storage(分布式存储)/MinIO.md Storage(分布式存储)/MinIO.md)。
  - ......等等，详见官方文档
- Index 支持以下远程存储
  - **Cassandra**
  - ......等等，详见官方文档

> Loki 的远程存储与 Prometheus 还有一点区别，Loki 自己的 Ingester 本身就实现了，而 Prometheus 的远程存储则需要其他程序对接。

现在用 MinIO 中存储的 Index 与 Chunks 查看一下目录结构
![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/gzp72g/1621436054664-ec14823d-2330-40c2-a6a1-155e6fc2b3b9.png)
fake 目录中就是 Chunks 数据
![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/gzp72g/1621436066582-d381149e-921e-47a1-9155-b716614528f7.png)
index 目录总就是 Index 数据
![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/gzp72g/1621436085579-707d7396-de4d-4761-b753-0bdd940de96b.png)

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

## Chunk 格式

      -------------------------------------------------------------------
      |                               |                                 |
      |        MagicNumber(4b)        |           version(1b)           |
      |                               |                                 |
      -------------------------------------------------------------------
      |         block-1 bytes         |          checksum (4b)          |
      -------------------------------------------------------------------
      |         block-2 bytes         |          checksum (4b)          |
      -------------------------------------------------------------------
      |         block-n bytes         |          checksum (4b)          |
      -------------------------------------------------------------------
      |                        # blocks (uvarint)                        |
      -------------------------------------------------------------------
      | #entries(uvarint) | mint, maxt (varint) | offset, len (uvarint) |
      -------------------------------------------------------------------
      | #entries(uvarint) | mint, maxt (varint) | offset, len (uvarint) |
      -------------------------------------------------------------------
      | #entries(uvarint) | mint, maxt (varint) | offset, len (uvarint) |
      -------------------------------------------------------------------
      | #entries(uvarint) | mint, maxt (varint) | offset, len (uvarint) |
      -------------------------------------------------------------------
      |                      checksum(from # blocks)                     |
      -------------------------------------------------------------------
      |           metasOffset - offset to the point with # blocks        |
      -------------------------------------------------------------------

`mint` 和 `maxt`分别描述了最小和最大的 Unix 纳秒时间戳。

### Block 格式

一个 block 由一系列日志 entries 组成，每个 entry 都是一个单独的日志行。
请注意，一个 block 的字节是用 Gzip 压缩存储的。以下是它们未压缩时的形式。

    -------------------------------------------------------------------
      |    ts (varint)    |     len (uvarint)    |     log-1 bytes      |
      -------------------------------------------------------------------
      |    ts (varint)    |     len (uvarint)    |     log-2 bytes      |
      -------------------------------------------------------------------
      |    ts (varint)    |     len (uvarint)    |     log-3 bytes      |
      -------------------------------------------------------------------
      |    ts (varint)    |     len (uvarint)    |     log-n bytes      |
      -------------------------------------------------------------------

`ts` 是日志的 Unix 纳秒时间戳，而 len 是日志条目的字节长度。

# Schema(模式) 概念

> 参考：
>
> - [官方文档，运维-存储-存储模式](https://grafana.com/docs/loki/latest/operations/storage/schema/)

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

- 比如现在有这么一种场景，公司本来将数据存放在本地，但是随着业务规模扩大，钱赚多了，数据存储方式想改变了，这时候如何丝滑的迁移数据呢？
- 通过 schema_config 字段的配置，2019-10-16 到 2020-11-16，使用本地存放数据；然后我买了公有云的存储产品，从 2020-11-16 到 2021-11-16，使用公有云存放数据。Loki 根据这种配置直到，当 2020-11-16 日开始，就会将数据存放在公有云数据库中了，而老的数据还可以不受影响。然后根据下文提到的 Retention 功能，逐步将老数据删除，这种升级过度是非常丝滑的。

## 在 Schema 中配置储存数据所使用的存储类型

并且，在 Schema 配置中，还可以通过 `schema_config.configs.store` 与 `schema_config.configs.object_store` 两个字段指定 Loki 储存 Index 和 Chunk 数据所使用的存储类型，比如上面这个例子中，对于 2020 年 7 月 1 日之前的数据，使用 boltdb-shipper 存储 Index 数据、filesystem 存储 Chunk 数据；2020 年 7 月 1 日之后的存储方式与之前一样。然后再通过 `storage_config` 字段来配置 boltdb-shipper 和 filesystem 这两类存储的具体使用方式，比如指定认证信息、指定存储路径等等。

同时，还可以通过`schema_config.configs.store` 与 `schema_config.configs.object_store` 两个字段来配置存储数据时 **Table(表)** 的行为，比如 存储时间、文件名前缀 等等。

比如上面这个示例中，配置了储存 Index 类型数据的一些简单行为：`prefix` 字段指定了当前期间表的前缀为 index\_，也就是存储 Index 数据的文件名开头都会是 index\_；`period` 字段则指定了在这个模式期间，表的周期为 168h，也就是每隔 168 小时，创建一张表。

## 依赖于 Schema 的配置

由于 Schema 中配置了 Loki 储存 Index 与 Chunks 两种数据所使用的存储类型。所以，很多配置，都会基于这个信息，来进行个性化配置。比如在 Schema 中设定了 S3 类型的存储，那么就需要设定连接 S3 存储时所需的认证信息。

> 注意：这一段内容会基于对 Loki 的存储概念已经有了详细了解后，才能看懂。

# BoltDB-Shipper 运行细节

Ingester 组件用于将 Index 与 Chunks 数据写入存储；Querier 组件用于从存储中读取 Index 和 Chunks 以处理 LogQL 查询请求。

## 写入数据

在深入了解细节之前，需要知道 Ingester 如何管理存储中的 Index 数据。

## 读取数据

Queriers lazily loads BoltDB files from shared object store to configured `cache_location`. When a querier receives a read request, the query range from the request is resolved to period numbers and all the files for those period numbers are downloaded to `cache_location`, if not already. Once we have downloaded files for a period we keep looking for updates in shared object store and download them every 5 Minutes by default. Frequency for checking updates can be configured with `resync_interval` config.
查询者将 BoltDB 文件从共享对象存储延迟加载到已配置的 cache_location。 当查询器接收到读取请求时，该请求的查询范围将解析为期间号，并将那些期间号的所有文件下载到 cache_location（如果尚未下载）。 下载文件一段时间后，我们会继续在共享库中查找更新，默认情况下每 5 分钟下载一次。 可以使用 resync_interval config 来配置检查更新的频率。

To avoid keeping downloaded index files forever there is a ttl for them which defaults to 24 hours, which means if index files for a period are not used for 24 hours they would be removed from cache location. ttl can be configured using `cache_ttl` config.为了避免永久保存下载的索引文件，有一个 ttl，默认值为 24 小时，这意味着如果一段时间内未使用索引文件 24 小时，它们将从缓存位置中删除。 可以使用 cache_ttl config 来配置 ttl。

**Note:** For better read performance and to avoid using node disk it is recommended to run Queriers as statefulset(when using k8s) with persistent storage for downloading and querying index files
.注意：为了获得更好的读取性能并避免使用节点磁盘，建议将 Queriers 作为 statefulset 运行（使用 k8s 时），并带有持久性存储，以下载和查询索引文件。

# Chunk 存储

Chunk 存储是 Loki 的长期数据存储，旨在支持交互式查询和持续写入，不需要后台维护任务。它由以下部分组成:

- 一个 chunks 索引，这个索引可以通过以下方式支持：Amazon DynamoDB、Google Bigtable、Apache Cassandra。
- 一个用于 chunk 数据本身的键值（KV）存储，可以是：Amazon DynamoDB、Google Bigtable、Apache Cassandra、Amazon S3、Google Cloud Storage。

> 与 Loki 的其他核心组件不同，块存储不是一个单独的服务、任务或进程，而是嵌入到需要访问 Loki 数据的 `ingester` 和 `querier` 服务中的一个库。

块存储依赖于一个统一的接口，用于支持块存储索引的 `NoSQL` 存储（DynamoDB、Bigtable 和 Cassandra）。这个接口假定索引是由以下项构成的键的条目集合。

- 一个哈希 key，对所有的读和写都是必需的。
- 一个范围 key，写入时需要，读取时可以省略，可以通过前缀或范围进行查询。

该接口在支持的数据库中的工作方式有些不同：

- `DynamoDB` 原生支持范围和哈希键，因此，索引条目被直接建模为 DynamoDB 条目，哈希键作为分布键，范围作为 DynamoDB 范围键。
- 对于 `Bigtable` 和 `Cassandra`，索引条目被建模为单个列值。哈希键成为行键，范围键成为列键。

一组模式集合被用来将读取和写入块存储时使用的匹配器和标签集映射到索引上的操作。随着 Loki 的发展，Schemas 模式也被添加进来，主要是为了更好地平衡写操作和提高查询性能。

# Table(表) 概念

> 参考：
>
> - [官方文档，运维-存储-表管理器](https://grafana.com/docs/loki/latest/operations/storage/table-manager/)
> - [开源中国博客，Loki|数据过期启动删除策略设计](https://my.oschina.net/u/1787735/blog/4429161)

Loki 可以将 Index 和 Chunks 数据以 **Table(表)** 的形式储存起来，凡是支持表结构的存储产品，都可以用来存储 Loki 的数据。使用 Table 时，会在一段时间内创建多个 Table，每个 Table 包含特定时间范围内的数据，所以 Table 也称为 **Periodic Table(周期表)。**

**注意：表的概念不适用与存储在 S3 的 Index 与 Chunk 数据**

如图所示，数据在存储中就是这样一种结构：
![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/gzp72g/1621406186450-f1bdfeb2-3225-472f-ae51-feee150458e8.png)
Loki 接收到的 Log Stream，会根据时间被分配到一个 Periodic Table 中

这是一种数据整合的方法，以便于更好的管理数。将数据基于表来进行分组有两个好处：

- **Schema(模式)** # 不同周期的表具有不同的 schema(模式) 和 版本。模式用来指定这个周期表内的数据使用什么类型的存储来存放数据。
- **Retention(保留)** # Retention 是通过删除整个表来实现的。
  - 通过将数据放在不同的周期内，在删除时也可以按照周期来删除，大大提高了数据处理效率

## Table Manager(表管理器)

**Table Manager(表管理器)** 是 Loki 的一个组件，用于创建和删除表。根据配置文件中 schema_config.configs 字段中的相关配置，在指定时间开始之前创建周期表，并在根据 table_manager 字段中的相关配置，将数据时间范围超过保留期的表内的数据删除

Table Manager 并不是可以管理所有存储类型中的日志流数据的，现阶段仅支持以下类型的存储：

- Index 数据
  - [BoltDB Shipper](https://grafana.com/docs/loki/latest/operations/storage/table-manager/boltdb-shipper/)
  - [Cassandra](https://cassandra.apache.org/)
  - [DynamoDB](https://aws.amazon.com/dynamodb)
  - ......等等，详见官方文档
- Chunk 数据
  - [Cassandra](https://cassandra.apache.org/)
  - [DynamoDB](https://aws.amazon.com/dynamodb)
  - Filesystem
  - ......等等，详见官方文档

注意：

- Table Manager 无法管理存放在对象存储(比如 S3)中的数据，如果要使用对象存储来储存 Index 与 Chunks 数据，则应该自行设置桶策略，以删除旧数据。

周期表存储相对于特定时间段的索引或块数据。在 schema_config 配置环境中配置储存在单个表中的数据的时间范围的持续时间及其存储类型。

schema_config 配置环境可以包含一个或多个配置。每个配置都定义了 .configs.from 字段指定的日期和下一个 .configs.from 配置之间这一段时间使用的存储模式，对于最后一个模式配置条目，则视为 .configs.from 指定的时间到当前时间这一段时间。这允许在一段时间内具有多个不重叠的架构配置，以便执行架构版本升级或更改存储设置（包括更改存储类型）。效果如下图：
![](https://notes-learning.oss-cn-beijing.aliyuncs.com/gzp72g/1616129694870-76a7338c-4755-41d8-8702-e6c9b37a5bc8.jpeg)
写入路径将命中日志条目时间戳记所在的表（通常是最后一个表，除了靠近表末尾和下一个表的开头的短时间段外），而读取路径将命中包含查询数据的表时间范围。

表管理器会在配置中设定的开始时间之前创建新表，以确保一旦到达当前表的结束时间，新表就准备就绪。

## Loki Storage Retention(数据保留)

> 官方文档：<https://grafana.com/docs/loki/latest/operations/storage/retention/>

Table Manager 可以实现 Loki 的 Retention(数据保留) 行为。如果想要启动 Retention，在 table_manager 配置环境中，需要开启删除数据保留的行为，以及确定数据可以保留的 period(期限，也就是时间范围)，详见 《Loki 配置详解》 table_manager 配置环境说明。

这个 Retention Period(保留期限) 是这么理解的：比如设置保留期限为 7d，表周期为 7d，那么当第三张表激活时，则第一张表被删除。
![](https://notes-learning.oss-cn-beijing.aliyuncs.com/gzp72g/1616129694876-b393cced-7824-4974-9c32-207b52852e71.jpeg)
这种设计允许进行快速删除操作，但要以保留期限由表格周期控制为代价。所以 Retention 必须为 Period 的倍数，比如当表周期为 7 天时，那么要保留的就必须是 7 天、14 天、21 天.....这种。否则无法正确删除一整张表。像 14 天，21 天这种，其实就意味着保留 2 张表、3 张表。

这是一个保留 28 天 日志数据的配置示例：主要配置点在最后 6 行

```yaml
schema_config:
  configs:
    - from: 2018-04-15
      store: boltdb
      object_store: filesystem
      schema: v11
      index:
        prefix: index_
        period: 24h

storage_config:
  boltdb:
    directory: /loki/index
  filesystem:
    directory: /loki/chunks

chunk_store_config:
  max_look_back_period: 672h

table_manager:
  retention_deletes_enabled: true
  retention_period: 672h
```

这个示例效果如下，当第六张新表创建之后，最早的表就会删除，至少保留 28 天，4 张表。![](https://notes-learning.oss-cn-beijing.aliyuncs.com/gzp72g/1616129694882-bf6a427c-5f6a-4ba5-93a0-a28c0ecd729e.jpeg)
注意：官方示例写的 30 天，是错误的，详见：<https://github.com/grafana/loki/pull/2772>
