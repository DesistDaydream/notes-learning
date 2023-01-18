---
title: Storage(存储)
---

# 概述

> 参考：
> - [官方文档,存储](https://prometheus.io/docs/prometheus/latest/storage/)
> - [GitHub,TSDB](https://github.com/prometheus/prometheus/tree/main/tsdb)
> - [GitHub 文档,TSDB format](https://github.com/prometheus/prometheus/blob/main/tsdb/docs/format/README.md)
> - [简书,Prometheus 存储机制](https://www.jianshu.com/p/ef9879dfb9ef)
> - [公众号,Prometheus 存储流向](https://mp.weixin.qq.com/s/J3oK0idEFbvErOwBEBrNSg)
> - 以下所有内容均基于 Prometheus 2.27+ 版本

Prometheus 自身就包含一个 **Time Series Database(时间序列数据库)**，所以 Prometheus 采集完指标数据后，可以保存在本地，由 Prometheus 自身来管理这些数据。当然，Prometheus 也可以通过一种称为 **Remote Write** 的技术，将数据存储到 **Remote Storage Systems(远程存储系统)**。

本地存储限制了 Prometheus 的可扩展性，带来了数据持久化、高科用等一系列的问题。为了解决单节点存储的限制，Prometheus 没有自己实现集群存储，而是提供了远程读写的接口，让用户自己选择合适的时序数据库来实现 Prometheus 的扩展性。

# Local Storage(本地存储)

**注意：** Prometheus 的本地存储不支持不兼容 POSIX 的文件系统，因为可能会发生不可恢复的损坏。不支持 NFS 文件系统（包括 AWS 的 EFS）。NFS 可能符合 POSIX，但大多数实现均不符合。强烈建议使用本地文件系统以提高可靠性。Prometheus 启动时会有如下 warn：
![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/lh6032/1623820971678-3d263b32-2760-4e77-9a22-b2c438bc62d5.png)
并且，经过实践，在数据量足够多时，当 Prometheus 压缩数据时，有不小的概率会丢失某个 Block 中的 meta.json 文件。进而导致压缩失败，并频繁产生告警，详见故障：[compaction failed](/docs/IT学习笔记/6.可观测性/监控系统/Prometheus/Prometheus%20 管理/故障处理/compaction%20failed.md 管理/故障处理/compaction failed.md)

Prometheus 的本地时间序列数据库将数据以自定义的高效格式存储在本地存储上。也就是说，Prometheus 采集到的指标数据，以文件的形式直接保存在操作系统的文件系统中。On-disk Layout 章节将会详细介绍这些数据在本地存储中布局。

## On-disk Layout(磁盘上的布局)

本地存储的目录看起来应该是下面这个样子：

    ./data
    ├── 01F5JX01DJSHFY98CREKE3F2FX
    │   ├── chunks
    │   │   └── 000001
    │   ├── index
    │   ├── meta.json
    │   └── tombstones
    ├── 01F5MQ3BR42QWB0JVKA1T4BBHP
    │   ├── chunks
    │   │   └── 000001
    │   ├── index
    │   ├── meta.json
    │   └── tombstones
    ├── 01F5MXZ46MQYP9QH0G0XVTQF0D
    │   ├── chunks
    │   │   └── 000001
    │   ├── index
    │   ├── meta.json
    │   └── tombstones
    ├── chunks_head
    │   ├── 000009
    │   └── 000010
    ├── queries.active
    └── wal
        ├── 00000008
        ├── 00000009
        ├── 00000010
        └── checkpoint.00000007
            └── 00000000

Prometheus 的存储大致可以分为两类

- Block(块) # 以 01 开头的那些目录。根据 [ULID](https://github.com/ulid/spec) 原理命名。
- Wal(预写日志) # wal 目录部分

注意：虽然持久化后的 Block 数据都是上述结构，但是在持久化之前，时序数据是保存在内存中，并且实现了 WAL 机制。

最新写入的数据保存在内存中的 Block 中，每隔 2 小时都会持久化到磁盘中(也就是生成 01F5JX01DJSHFY98CREKE3F2FX 这种块目录)。为了防止程序崩溃导致数据丢失，实现了 **Write Ahead Log(预写日志，简称 WAL)** 机制，启动时会以写入日志(WAL)的方式来实现重播，从而恢复数据。wal 目录中的这些原始数据尚未被 **Compaction(压缩)**，因为，它们的大小明显要超过 Block 中 chunks 目录中数据的大小。Prometheus 最少保留 3 个 WAL 文件。

通过 Block 的形式保存所有的时序数据，可以明显提高 Prometheus 的查询效率，当查询一段时间范围内的所有样本数据时，只需要简单的从落在该范围内的 Block 中查询数据即可。

### Block(块)

Prometheus 存储在本地的时间序列数据，被抽象为一个一个的 **Block(块)。**每个 Block 都是一个单独的目录，Block 由 4 个部分组成：

- **chunks/\*** # Block(块) 中的所有时序数据所在的子目录。
  - chunks 目录中的时序数据被分组为一个或多个分段文件，默认情况下，每个文件的最大容量为 512MiB。
- **meta.json** # 元数据文件
- **index** # 索引文件。根据指标名称和标签索引到 chunks 目录中的时间序列数据
- **tombstones** # 如果通过 API 删除时序数据，删除记录会保存在单独的逻辑文件 `tombstone` 当中。
  - 也就是说，被删除的数据不会直接立即删除。而是通过 tombstones 文件建立一个删除记录，在通过 PromQL 查找数据时，不会搜索 tombstones 文件中标记的数据。

默认情况下，一个 Block(块) 最少包含 2 个小时的时序数据。可以通过下面这些参数设置每个 Block 所包含数据的时间周期。

- \--storage.tsdb.min-block-duration # 一个存储 Block 的最小时间。默认 2 小时
- \--storage.tsdb.max-block-duration # 一个存储 Block 的最大时间
  - 每隔一段时间，这些 2 小时的 Block 将会通过 Compaction 机制，压缩成时间周期更长的 Block，以节省存储空间。通常这个时间周期是 --storage.tsdb.retention 标志指定的时间的 10%，若是 10% 的结果小于 31 天，则默认最大时间为 31 天。
- \--storage.tsdb.retention # 块的过期时间.
- **举个栗子**:
- 假设有如下设置:
  - \--storage.tsdb.max-block-duration=1h
  - \--storage.tsdb.max-block-duration=15m
  - \--storage.tsdb.retention=2h
- 再假设你在今天的 16:00 搜索了数据,那么你最多可以搜索到今天 13:00(即 16-(2-1))的数据.而最少也可以搜索到 14:45(如果期间数据在产生)往后的数据。

我们将存储层划分为一个一个的 Block(块)，每个块在一段时间内保存所有序列。每个块充当独立数据库。
![1889435-999d351beafab3c6.jpg](https://notes-learning.oss-cn-beijing.aliyuncs.com/lh6032/1620917933638-6655ade5-1636-43c7-8c72-20889f3218ed.jpeg)
这样每次查询，仅检查所请求的时间范围内的块子集，查询执行时间自然会减少。
这种布局也使删除旧数据变得非常容易，一旦块的时间范围完全落后于配置的保留边界，它就可以完全丢弃。
![1889435-af09c18b8bbeb5fc.jpg](https://notes-learning.oss-cn-beijing.aliyuncs.com/lh6032/1620917933635-d92e1ace-518f-4c73-b33e-03688b64b9ec.jpeg)

### Index(索引)

一般 Prometheus 的查询是把 metric+label 做关键字的，而且是很宽泛，完全用户自定义的字符，因此没办法使用常规的 sql 数据库，prometheus 的存储层使用了全文检索中的[倒排索引](https://nlp.stanford.edu/IR-book/html/htmledition/a-first-take-at-building-an-inverted-index-1.html)概念，将每个时间序列视为一个小文档。而 metric 和 label 对应的是文档中的单词。
例如，requests_total{path="/status", method="GET", instance="10.0.0.1:80"}是包含以下单词的文档：

- name="requests_total"
- path="/status"
- method="GET"
- instance="10.0.0.1:80"

### Compaction(压缩)

这些 2 小时的 Block 会在后台压缩成更大的 Block，数据压缩合并成更高级别的 Block 文件后删除低级别的 Block 文件。一个高级别的块通常包含数据保留时间 10%的时间周期的时序数据，若是 10% 小于 31 天，则默认为 31 天。

这个和 leveldb、rocksdb 等 LSM 树的思路一致。这些设计和 Gorilla 的设计高度相似，所以 Prometheus 几乎就是等于一个缓存 TSDB。它本地存储的特点决定了它不能用于 long-term 数据存储，只能用于短期窗口的 timeseries 数据保存和查询，并且不具有高可用性（宕机会导致历史数据无法读取）。

所以，Prometheus 实现了下文的 [Remote Storage 功能](</docs/IT学习笔记/6.可观测性/监控系统/Prometheus/Storage(存储).md>>)，可以通过该功能，将数据通过网络转存到其他存储中。但是，需要仔细评估它们，性能和效率方面会产生很大的变化。

现有存储层的样本压缩功能在 Prometheus 的早期版本中发挥了重要作用。单个原始数据点占用 16 个字节的存储空间。但当普罗米修斯每秒收集数十万个数据点时，可以快速填满硬盘。但，同一系列中的样本往往非常相似，我们可以利用这一类样品（同样 label）进行有效的压缩。批量压缩一系列的许多样本的块，在内存中，将每个数据点压缩到平均 1.37 字节的存储。这种压缩方案运行良好，也保留在新版本 2 存储层的设计中。具体压缩算法可以参考：[Facebook 的“Gorilla”论文中](http://www.vldb.org/pvldb/vol8/p1816-teller.pdf)

### 基准测试

cpu、内存、查询效率都比 1.x 版本得到了大幅度的提升
具体测试结果参考：<https://dzone.com/articles/prometheus-2-times-series-storage-performance-anal>

## 存储配置

对于本地存储，prometheus 提供了一些配置项，主要包括：

- \--storage.tsdb.path: 存储数据的目录，默认为 data/，如果要挂外部存储，可以指定该目录
- \--storage.tsdb.retention.time: 数据过期清理时间，默认保存 15 天
- \--storage.tsdb.retention.size: 实验性质，声明数据块的最大值，不包括 wal 文件，如 512MB

Prometheus 将所有当前使用的块保留在内存中。此外，它将最新使用的块保留在内存中，最大内存可以通过 storage.local.memory-chunks 标志配置。

### 容量规划

容量规划除了上边说的内存，还有磁盘存储规划，这和你的 Prometheus 的架构方案有关。

- 如果是单机 Prometheus，计算本地磁盘使用量。
- 如果是 Remote-Write，和已有的 Tsdb 共用即可。
- 如果是 Thanos 方案，本地磁盘可以忽略（2H)，计算对象存储的大小就行。

在一般情况下，Prometheus 中存储的每一个样本大概占用 1-2 字节大小。如果需要对 Prometheus Server 的本地磁盘空间做容量规划时，可以通过以下公式计算：

    磁盘大小 = 保留时间 * 每秒获取样本数 * 样本大小

**保留时间(retention_time_seconds) **和 **样本大小(bytes_per_sample)** 不变的情况下，如果想减少本地磁盘的容量需求，只能通过减少每秒获取样本数(ingested_samples_per_second)的方式。

因此有两种手段，一是减少时间序列的数量，二是增加采集样本的时间间隔。

考虑到 Prometheus 会对时间序列进行压缩，因此减少时间序列的数量效果更明显。

Prometheus 每 2 小时将已缓冲在内存中的数据压缩到磁盘上的块中。包括 Chunks、Indexes、Tombstones、Metadata，这些占用了一部分存储空间。一般情况下，Prometheus 中存储的每一个样本大概占用 1-2 字节大小（1.7Byte）。可以通过 PromQL 来查看每个样本平均占用多少空间：

    rate(prometheus_tsdb_compaction_chunk_size_bytes_sum[2h])
    /
    rate(prometheus_tsdb_compaction_chunk_samples_sum[2h])

     {instance="0.0.0.0:8890", job="prometheus"}  1.252747585939941

查看当前每秒获取的样本数：

```shell
rate(prometheus_tsdb_head_samples_appended_total[1h])
```

有两种手段，一是减少时间序列的数量，二是增加采集样本的时间间隔。考虑到 Prometheus 会对时间序列进行压缩，因此减少时间序列的数量效果更明显。

举例说明：

- 采集频率 30s，机器数量 1000，Metric 种类 6000，1000_6000_2_60_24 约 200 亿，30G 左右磁盘。
- 只采集需要的指标，如 match\[], 或者统计下最常使用的指标，性能最差的指标。

以上磁盘容量并没有把 wal 文件算进去，wal 文件 (Raw Data) 在 Prometheus 官方文档中说明至少会保存 3 个 Write-Ahead Log Files，每一个最大为 128M(实际运行发现数量会更多)。

因为我们使用了 Thanos 的方案，所以本地磁盘只保留 2H 热数据。Wal 每 2 小时生成一份 Block 文件，Block 文件每 2 小时上传对象存储，本地磁盘基本没有压力。

关于 Prometheus 存储机制，可以看[这篇](http://www.xuyasong.com/?p=1601)。

## 故障恢复

如果怀疑数据库中的损坏引起的问题，则可以通过使用 storage.local.dirtyflag 配置，来启动服务器来强制执行崩溃恢复。
如果没有帮助，或者如果您只想删除现有的数据库，可以通过删除存储目录的内容轻松地启动

# Remote Storage(远程存储)

Prometheus 的本地存储在可伸缩性和持久性方面受到单个节点的限制。Prometheus 并没有尝试解决 Prometheus 本身中的集群存储，而是提供了一组允许与远程存储系统集成的接口。

Prometheus 通过下面几种方式与远程存储系统集成：

- Prometheus 可以以标准格式将其采集到的样本数据写入到指定的远程 URL。
- Prometheus 可以以标准格式从指定的远程 URL 读取(返回)样本数据。
- Prometheus 可以以标准格式从其他 Prometheus 接收样本。![](https://notes-learning.oss-cn-beijing.aliyuncs.com/lh6032/1616069469195-edb3fcc9-e672-43be-b6b9-fcc52d6ed497.jpeg)

说白了，Prometheus 规定了一种标准格式，可以将采集到的指标数据实时发送给 Adapter，然后由 Adapter 处理后，在存储在第三方存储中(比如 InfluxDB、OpenTSDB 等等)。

同时，Prometheus 自身也自带了一个 Adapter，可以在启动程序时，指定 `--web.enable-remote-write-receiver` 标志即可，此时，Prometheus 会在 `/api/v1/write` 端点上暴露 Remote Write API，其他 Prometheus 可以将采集到的指标数据发送到 `http://PrometheusIP:PORT:9090/api/v1/write` 上，这与 Federate(联邦) 功能有点类似，都可以用来汇总数据的。此时，这个开启了 Remote Write API 的 Prometheus 通常被称为 **Receiver(接收器)**，象征着这个 Prometheus 可以接收其他符合 Prometheus 标准格式的指标数据。

其他的集成在 Adapter 要么可以自己实现，要么就继承在第三方存储中，在 [官方文档,集成方式-远程端点和存储](https://prometheus.io/docs/operating/integrations/#remote-endpoints-and-storage) 章节中可以看到现阶段所有可以实现 Remote Write API 的 Adapter 以及 第三方存储。

有关在 Prometheus 中配置远程存储集成的详细信息，请参阅 Prometheus 配置文档的 [远程写入](https://prometheus.io/docs/prometheus/latest/configuration/configuration/#remote_write) 和[ 远程读取](https://prometheus.io/docs/prometheus/latest/configuration/configuration/#remote_read) 部分。

有关请求和响应消息的详细信息，请参阅[远程存储协议缓冲区定义](https://github.com/prometheus/prometheus/blob/master/prompb/remote.proto)。

注意：

- 读写协议都使用基于 HTTP 的快速压缩协议缓冲区编码。该协议尚未被认为是稳定的 API，当可以安全地假定 Prometheus 和远程存储之间的所有跃点都支持 HTTP/2 时，该协议将来可能会更改为在 HTTP/2 上使用 gRPC。
- 在 Remote Read 的实现中，读取路径上，Prometheus 仅从远端获取一组标签选择器和时间范围的原始系列数据。PromQL 对原始数据的所有评估仍然在 Prometheus 本身中进行。这意味着远程读取查询具有一定的可伸缩性限制，因为所有必需的数据都需要先加载到查询的 Prometheus 服务器中，然后再在其中进行处理。但是，暂时认为支持 PromQL 的完全分布式评估是不可行的。

## 远程读

在远程读的流程当中，当用户发起查询请求后，Promthues 将向 remote_read 中配置的 URL 发起查询请求(matchers,ranges)，Adaptor 根据请求条件从第三方存储服务中获取响应的数据。同时将数据转换为 Promthues 的原始样本数据返回给 Prometheus Server。
当获取到样本数据后，Promthues 在本地使用 PromQL 对样本数据进行二次处理。

## 远程写

用户可以在 Promtheus 配置文件中指定 Remote Write(远程写) 的 URL 地址，一旦设置了该配置项，Prometheus 将样本数据通过 HTTP 的形式发送给 Adaptor(适配器)。而用户则可以在适配器中对接外部任意的服务。外部服务可以是真正的存储系统，公有云的存储服务，也可以是消息队列等任意形式。

## 配置

配置非常简单，只需要将对应的地址配置下就行

    remote_write:
      - url: "http://localhost:9201/write"
    remote_read:
      - url: "http://localhost:9201/read"

####

# 压缩示例

```bash
[root@hw-cloud-xngy-jump-server-linux-2 /mnt/sfs_turbo/monitoring-prometheus-prometheus-monitor-hw-cloud-k8s-prometheus-0-pvc-9ca02cc7-33f2-4059-807d-196c78a1e728/prometheus-db]# ll
total 180
drwxrwxrwx 40 root      root  4096 Oct 11 09:00 ./
drwxrwxrwx  3 root      root  4096 Aug 17 09:40 ../
drwxr-xr-x  3 lichenhao 2000  4096 Aug 30 13:26 01FEARAQ8BSD82FA4TDR516476/
drwxr-xr-x  3 lichenhao 2000  4096 Aug 30 13:26 01FEARB5KS3BXX3GPQTMDQ5ZFD/
drwxr-xr-x  3 lichenhao 2000  4096 Sep  3 13:01 01FEN0FGFRRF6VKPM1RT41SXJ4/
drwxr-xr-x  3 lichenhao 2000  4096 Sep 10 07:02 01FF6CNJS2V4QK65SXG31QMQZ6/
drwxr-xr-x  3 lichenhao 2000  4096 Sep 17 01:02 01FFQRVFSQHH62X3K9CPA2Y7MZ/
drwxr-xr-x  3 lichenhao 2000  4096 Sep 23 19:02 01FG95195EKPKK985Z1XJQ4M1P/
drwxr-xr-x  3 lichenhao 2000  4096 Sep 30 13:03 01FGTH81ZQ5RPGTD5VNWXYWTAK/
drwxr-xr-x  3 lichenhao 2000  4096 Oct  7 07:05 01FHBXE6D4J7994S2FWQTBSN1K/
drwxr-xr-x  3 lichenhao 2000  4096 Oct  8 01:00 01FHDV5JG0JT9V0V5F7EJPGB7R/
drwxr-xr-x  2 lichenhao 2000  4096 Oct  9 13:01 01FHFRZ1B9EKAHBYA9VF7MABH1.tmp-for-deletion/
drwxr-xr-x  3 lichenhao 2000  4096 Oct  9 11:00 01FHHFWEFK55FNSAVZQQWJ5XZ6/
drwxr-xr-x  3 lichenhao 2000  4096 Oct  9 13:00 01FHHPR5QJG9C5Q98WHDAEJJM1/
drwxr-xr-x  3 lichenhao 2000  4096 Oct  9 13:00 01FHHPRNDK82JK6X12V6CX9SSK/
drwxr-xr-x  3 lichenhao 2000  4096 Oct  9 13:01 01FHHPS3NR7M2E8MAV37S61ME6/
drwxr-xr-x  3 lichenhao 2000  4096 Oct  9 15:00 01FHHXKY9K68GB1DST73HWM4W4/
drwxr-xr-x  3 lichenhao 2000  4096 Oct  9 17:00 01FHJ4FM7G4WVNH4SR7N4909PQ/
drwxr-xr-x  3 lichenhao 2000  4096 Oct  9 19:00 01FHJBBBGRB14DDKQ6G0MM3N36/
drwxr-xr-x  3 lichenhao 2000  4096 Oct  9 21:00 01FHJJ72QHW1JMZY44466EAFPE/
drwxr-xr-x  3 lichenhao 2000  4096 Oct  9 23:00 01FHJS2SZH369ZGEQF2WJCX97N/
drwxr-xr-x  3 lichenhao 2000  4096 Oct 10 01:00 01FHJZYH7G903FDNBKKGKGGZM8/
drwxr-xr-x  3 lichenhao 2000  4096 Oct 10 03:00 01FHK6T8FKEQE3F61EZAT0PB6C/
drwxr-xr-x  3 lichenhao 2000  4096 Oct 10 05:00 01FHKDNZQJZKKCV9AXFDGQ48M3/
drwxr-xr-x  3 lichenhao 2000  4096 Oct 10 07:00 01FHKMHPZH8DSC18W69H6E1H7X/
drwxr-xr-x  3 lichenhao 2000  4096 Oct 10 09:00 01FHKVDE7F7AXTCMAT437TAW8X/
drwxr-xr-x  3 lichenhao 2000  4096 Oct 10 11:00 01FHM295FJ6GWFC1YSBVN1P2SC/
drwxr-xr-x  3 lichenhao 2000  4096 Oct 10 13:00 01FHM94WQJP3BYQV6X4RTG5G6T/
drwxr-xr-x  3 lichenhao 2000  4096 Oct 10 15:00 01FHMG0KZHBH8J58D874RTHCEP/
drwxr-xr-x  3 lichenhao 2000  4096 Oct 10 17:00 01FHMPWB7FRVXRAGZBDB8JM1FG/
drwxr-xr-x  3 lichenhao 2000  4096 Oct 10 19:00 01FHMXR2FK0VZTDTH2M43KXTQ5/
drwxr-xr-x  3 lichenhao 2000  4096 Oct 10 21:00 01FHN4KSQFXY291W9PDWR15XPC/
drwxr-xr-x  3 lichenhao 2000  4096 Oct 10 23:00 01FHNBFGZGXTR6KZHSQ7HPG19E/
drwxr-xr-x  3 lichenhao 2000  4096 Oct 11 01:00 01FHNJB87GHRBG8BAVBGZG5FHE/
drwxr-xr-x  3 lichenhao 2000  4096 Oct 11 03:00 01FHNS6ZFJ3Z4AP894F4YE5RKP/
drwxr-xr-x  3 lichenhao 2000  4096 Oct 11 05:00 01FHP02PQG6NA8REWSAEKP1STW/
drwxr-xr-x  3 lichenhao 2000  4096 Oct 11 07:00 01FHP6YDZFB88DFPYB7RQBSMA2/
drwxr-xr-x  3 lichenhao 2000  4096 Oct 11 09:00 01FHPDT57J3XW0SSQRDSAAASF1/
drwxr-xr-x  2 lichenhao 2000  4096 Oct 11 09:00 chunks_head/
-rw-r--r--  1 lichenhao 2000     0 Sep 27 21:49 lock
-rw-r--r--  1 lichenhao 2000 20001 Oct 11 10:09 queries.active
drwxr-xr-x  3 lichenhao 2000  4096 Oct 11 09:00 wal/
[root@hw-cloud-xngy-jump-server-linux-2 /mnt/sfs_turbo/monitoring-prometheus-prometheus-monitor-hw-cloud-k8s-prometheus-0-pvc-9ca02cc7-33f2-4059-807d-196c78a1e728/prometheus-db]# mv 01FHHPS3NR7M2E8MAV37S61ME6 /root/backup/

[root@hw-cloud-xngy-jump-server-linux-2 /mnt/sfs_turbo/monitoring-prometheus-prometheus-monitor-hw-cloud-k8s-prometheus-0-pvc-9ca02cc7-33f2-4059-807d-196c78a1e728/prometheus-db]# [root@hw-cloud-xngy-jump-server-linux-2 /mnt/sfs_turbo/monitoring-prometheus-prometheus-monitor-hw-cloud-k8s-prometheus-0-pvc-9ca02cc7-33f2-4059-807d-196c78a1e728/prometheus-db]#
[root@hw-cloud-xngy-jump-server-linux-2 /mnt/sfs_turbo/monitoring-prometheus-prometheus-monitor-hw-cloud-k8s-prometheus-0-pvc-9ca02cc7-33f2-4059-807d-196c78a1e728/prometheus-db]#
[root@hw-cloud-xngy-jump-server-linux-2 /mnt/sfs_turbo/monitoring-prometheus-prometheus-monitor-hw-cloud-k8s-prometheus-0-pvc-9ca02cc7-33f2-4059-807d-196c78a1e728/prometheus-db]# ll
total 176
drwxrwxrwx 39 root      root  4096 Oct 11 10:09 ./
drwxrwxrwx  3 root      root  4096 Aug 17 09:40 ../
drwxr-xr-x  3 lichenhao 2000  4096 Aug 30 13:26 01FEARAQ8BSD82FA4TDR516476/
drwxr-xr-x  3 lichenhao 2000  4096 Aug 30 13:26 01FEARB5KS3BXX3GPQTMDQ5ZFD/
drwxr-xr-x  3 lichenhao 2000  4096 Sep  3 13:01 01FEN0FGFRRF6VKPM1RT41SXJ4/
drwxr-xr-x  3 lichenhao 2000  4096 Sep 10 07:02 01FF6CNJS2V4QK65SXG31QMQZ6/
drwxr-xr-x  3 lichenhao 2000  4096 Sep 17 01:02 01FFQRVFSQHH62X3K9CPA2Y7MZ/
drwxr-xr-x  3 lichenhao 2000  4096 Sep 23 19:02 01FG95195EKPKK985Z1XJQ4M1P/
drwxr-xr-x  3 lichenhao 2000  4096 Sep 30 13:03 01FGTH81ZQ5RPGTD5VNWXYWTAK/
drwxr-xr-x  3 lichenhao 2000  4096 Oct  7 07:05 01FHBXE6D4J7994S2FWQTBSN1K/
drwxr-xr-x  3 lichenhao 2000  4096 Oct  8 01:00 01FHDV5JG0JT9V0V5F7EJPGB7R/
drwxr-xr-x  2 lichenhao 2000  4096 Oct  9 13:01 01FHFRZ1B9EKAHBYA9VF7MABH1.tmp-for-deletion/
drwxr-xr-x  3 lichenhao 2000  4096 Oct  9 11:00 01FHHFWEFK55FNSAVZQQWJ5XZ6/
drwxr-xr-x  3 lichenhao 2000  4096 Oct  9 13:00 01FHHPR5QJG9C5Q98WHDAEJJM1/
drwxr-xr-x  3 lichenhao 2000  4096 Oct  9 13:00 01FHHPRNDK82JK6X12V6CX9SSK/
drwxr-xr-x  3 lichenhao 2000  4096 Oct  9 15:00 01FHHXKY9K68GB1DST73HWM4W4/
drwxr-xr-x  3 lichenhao 2000  4096 Oct  9 17:00 01FHJ4FM7G4WVNH4SR7N4909PQ/
drwxr-xr-x  3 lichenhao 2000  4096 Oct  9 19:00 01FHJBBBGRB14DDKQ6G0MM3N36/
drwxr-xr-x  3 lichenhao 2000  4096 Oct  9 21:00 01FHJJ72QHW1JMZY44466EAFPE/
drwxr-xr-x  3 lichenhao 2000  4096 Oct  9 23:00 01FHJS2SZH369ZGEQF2WJCX97N/
drwxr-xr-x  3 lichenhao 2000  4096 Oct 10 01:00 01FHJZYH7G903FDNBKKGKGGZM8/
drwxr-xr-x  3 lichenhao 2000  4096 Oct 10 03:00 01FHK6T8FKEQE3F61EZAT0PB6C/
drwxr-xr-x  3 lichenhao 2000  4096 Oct 10 05:00 01FHKDNZQJZKKCV9AXFDGQ48M3/
drwxr-xr-x  3 lichenhao 2000  4096 Oct 10 07:00 01FHKMHPZH8DSC18W69H6E1H7X/
drwxr-xr-x  3 lichenhao 2000  4096 Oct 10 09:00 01FHKVDE7F7AXTCMAT437TAW8X/
drwxr-xr-x  3 lichenhao 2000  4096 Oct 10 11:00 01FHM295FJ6GWFC1YSBVN1P2SC/
drwxr-xr-x  3 lichenhao 2000  4096 Oct 10 13:00 01FHM94WQJP3BYQV6X4RTG5G6T/
drwxr-xr-x  3 lichenhao 2000  4096 Oct 10 15:00 01FHMG0KZHBH8J58D874RTHCEP/
drwxr-xr-x  3 lichenhao 2000  4096 Oct 10 17:00 01FHMPWB7FRVXRAGZBDB8JM1FG/
drwxr-xr-x  3 lichenhao 2000  4096 Oct 10 19:00 01FHMXR2FK0VZTDTH2M43KXTQ5/
drwxr-xr-x  3 lichenhao 2000  4096 Oct 10 21:00 01FHN4KSQFXY291W9PDWR15XPC/
drwxr-xr-x  3 lichenhao 2000  4096 Oct 10 23:00 01FHNBFGZGXTR6KZHSQ7HPG19E/
drwxr-xr-x  3 lichenhao 2000  4096 Oct 11 01:00 01FHNJB87GHRBG8BAVBGZG5FHE/
drwxr-xr-x  3 lichenhao 2000  4096 Oct 11 03:00 01FHNS6ZFJ3Z4AP894F4YE5RKP/
drwxr-xr-x  3 lichenhao 2000  4096 Oct 11 05:00 01FHP02PQG6NA8REWSAEKP1STW/
drwxr-xr-x  3 lichenhao 2000  4096 Oct 11 07:00 01FHP6YDZFB88DFPYB7RQBSMA2/
drwxr-xr-x  3 lichenhao 2000  4096 Oct 11 09:00 01FHPDT57J3XW0SSQRDSAAASF1/
drwxr-xr-x  2 lichenhao 2000  4096 Oct 11 09:00 chunks_head/
-rw-r--r--  1 lichenhao 2000     0 Sep 27 21:49 lock
-rw-r--r--  1 lichenhao 2000 20001 Oct 11 10:09 queries.active
drwxr-xr-x  3 lichenhao 2000  4096 Oct 11 10:09 wal/
[root@hw-cloud-xngy-jump-server-linux-2 /mnt/sfs_turbo/monitoring-prometheus-prometheus-monitor-hw-cloud-k8s-prometheus-0-pvc-9ca02cc7-33f2-4059-807d-196c78a1e728/prometheus-db]# ll 01FHHPS3NR7M2E8MAV37S61ME6
ls: cannot access '01FHHPS3NR7M2E8MAV37S61ME6': No such file or directory
[root@hw-cloud-xngy-jump-server-linux-2 /mnt/sfs_turbo/monitoring-prometheus-prometheus-monitor-hw-cloud-k8s-prometheus-0-pvc-9ca02cc7-33f2-4059-807d-196c78a1e728/prometheus-db]# ll
total 136
drwxrwxrwx 29 root      root  4096 Oct 11 10:35 ./
drwxrwxrwx  3 root      root  4096 Aug 17 09:40 ../
drwxr-xr-x  3 lichenhao 2000  4096 Aug 30 13:26 01FEARAQ8BSD82FA4TDR516476/
drwxr-xr-x  3 lichenhao 2000  4096 Aug 30 13:26 01FEARB5KS3BXX3GPQTMDQ5ZFD/
drwxr-xr-x  3 lichenhao 2000  4096 Sep  3 13:01 01FEN0FGFRRF6VKPM1RT41SXJ4/
drwxr-xr-x  3 lichenhao 2000  4096 Sep 10 07:02 01FF6CNJS2V4QK65SXG31QMQZ6/
drwxr-xr-x  3 lichenhao 2000  4096 Sep 17 01:02 01FFQRVFSQHH62X3K9CPA2Y7MZ/
drwxr-xr-x  3 lichenhao 2000  4096 Sep 23 19:02 01FG95195EKPKK985Z1XJQ4M1P/
drwxr-xr-x  3 lichenhao 2000  4096 Sep 30 13:03 01FGTH81ZQ5RPGTD5VNWXYWTAK/
drwxr-xr-x  3 lichenhao 2000  4096 Oct  7 07:05 01FHBXE6D4J7994S2FWQTBSN1K/
drwxr-xr-x  3 lichenhao 2000  4096 Oct  8 01:00 01FHDV5JG0JT9V0V5F7EJPGB7R/
drwxr-xr-x  3 lichenhao 2000  4096 Oct  9 13:00 01FHHPRNDK82JK6X12V6CX9SSK/
drwxr-xr-x  3 lichenhao 2000  4096 Oct 10 17:00 01FHMPWB7FRVXRAGZBDB8JM1FG/
drwxr-xr-x  3 lichenhao 2000  4096 Oct 10 19:00 01FHMXR2FK0VZTDTH2M43KXTQ5/
drwxr-xr-x  3 lichenhao 2000  4096 Oct 10 21:00 01FHN4KSQFXY291W9PDWR15XPC/
drwxr-xr-x  3 lichenhao 2000  4096 Oct 10 23:00 01FHNBFGZGXTR6KZHSQ7HPG19E/
drwxr-xr-x  3 lichenhao 2000  4096 Oct 11 01:00 01FHNJB87GHRBG8BAVBGZG5FHE/
drwxr-xr-x  3 lichenhao 2000  4096 Oct 11 03:00 01FHNS6ZFJ3Z4AP894F4YE5RKP/
drwxr-xr-x  3 lichenhao 2000  4096 Oct 11 05:00 01FHP02PQG6NA8REWSAEKP1STW/
drwxr-xr-x  3 lichenhao 2000  4096 Oct 11 07:00 01FHP6YDZFB88DFPYB7RQBSMA2/
drwxr-xr-x  3 lichenhao 2000  4096 Oct 11 09:00 01FHPDT57J3XW0SSQRDSAAASF1/
drwxr-xr-x  3 lichenhao 2000  4096 Oct 11 10:33 01FHPK51H3650KFJ4QVSD1VRBX/
drwxr-xr-x  3 lichenhao 2000  4096 Oct 11 10:33 01FHPK5B0QGC4WTV8HXBAJC3HH/
drwxr-xr-x  3 lichenhao 2000  4096 Oct 11 10:33 01FHPK5KH2PG7G01DQECF0CN3Q/
drwxr-xr-x  3 lichenhao 2000  4096 Oct 11 10:33 01FHPK5WG49AN9B78W08PH8Q2M/
drwxr-xr-x  3 lichenhao 2000  4096 Oct 11 10:34 01FHPK65DB93Q1C71QSMBGM36C/
drwxr-xr-x  3 lichenhao 2000  4096 Oct 11 10:35 01FHPK7W5TEA4JV0388AFMH0TM/
drwxr-xr-x  2 lichenhao 2000  4096 Oct 11 10:32 chunks_head/
-rw-r--r--  1 lichenhao 2000     0 Oct 11 10:31 lock
-rw-r--r--  1 lichenhao 2000 20001 Oct 11 10:35 queries.active
drwxr-xr-x  3 lichenhao 2000  4096 Oct 11 10:32 wal/
[root@hw-cloud-xngy-jump-server-linux-2 /mnt/sfs_turbo/monitoring-prometheus-prometheus-monitor-hw-cloud-k8s-prometheus-0-pvc-9ca02cc7-33f2-4059-807d-196c78a1e728/prometheus-db]# ll
total 96
drwxrwxrwx 19 root      root  4096 Oct 11 10:36 ./
drwxrwxrwx  3 root      root  4096 Aug 17 09:40 ../
drwxr-xr-x  3 lichenhao 2000  4096 Aug 30 13:26 01FEARAQ8BSD82FA4TDR516476/
drwxr-xr-x  3 lichenhao 2000  4096 Aug 30 13:26 01FEARB5KS3BXX3GPQTMDQ5ZFD/
drwxr-xr-x  3 lichenhao 2000  4096 Sep  3 13:01 01FEN0FGFRRF6VKPM1RT41SXJ4/
drwxr-xr-x  3 lichenhao 2000  4096 Sep 10 07:02 01FF6CNJS2V4QK65SXG31QMQZ6/
drwxr-xr-x  3 lichenhao 2000  4096 Sep 17 01:02 01FFQRVFSQHH62X3K9CPA2Y7MZ/
drwxr-xr-x  3 lichenhao 2000  4096 Sep 23 19:02 01FG95195EKPKK985Z1XJQ4M1P/
drwxr-xr-x  3 lichenhao 2000  4096 Sep 30 13:03 01FGTH81ZQ5RPGTD5VNWXYWTAK/
drwxr-xr-x  3 lichenhao 2000  4096 Oct  7 07:05 01FHBXE6D4J7994S2FWQTBSN1K/
drwxr-xr-x  3 lichenhao 2000  4096 Oct 11 05:00 01FHP02PQG6NA8REWSAEKP1STW/
drwxr-xr-x  3 lichenhao 2000  4096 Oct 11 07:00 01FHP6YDZFB88DFPYB7RQBSMA2/
drwxr-xr-x  3 lichenhao 2000  4096 Oct 11 09:00 01FHPDT57J3XW0SSQRDSAAASF1/
drwxr-xr-x  3 lichenhao 2000  4096 Oct 11 10:35 01FHPK919H3DMC9YQMQX6D2MPY/
drwxr-xr-x  3 lichenhao 2000  4096 Oct 11 10:36 01FHPK9V6C09SHQ0VK189R8VR4/
drwxr-xr-x  3 lichenhao 2000  4096 Oct 11 10:36 01FHPKA8TKRQ080T2KYKCZGBHX/
drwxr-xr-x  3 lichenhao 2000  4096 Oct 11 10:36 01FHPKAQNCTF4TKTPHCTW87EKV/
drwxr-xr-x  2 lichenhao 2000  4096 Oct 11 10:32 chunks_head/
-rw-r--r--  1 lichenhao 2000     0 Oct 11 10:31 lock
-rw-r--r--  1 lichenhao 2000 20001 Oct 11 10:37 queries.active
drwxr-xr-x  3 lichenhao 2000  4096 Oct 11 10:32 wal/
[root@hw-cloud-xngy-jump-server-linux-2 /mnt/sfs_turbo/monitoring-prometheus-prometheus-monitor-hw-cloud-k8s-prometheus-0-pvc-9ca02cc7-33f2-4059-807d-196c78a1e728/prometheus-db]#
```

可以看到，Prometheus 逐步压缩一天的所有 Block，并逐步压缩到单一的 Block 中。10 月 9 日 与 10 日的 Block 逐步压缩，统一到了 10 月 7 日的 Block 中。
