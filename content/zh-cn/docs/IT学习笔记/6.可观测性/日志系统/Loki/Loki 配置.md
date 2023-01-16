---
title: Loki 配置
---

# 概述

> 参考：
> - [官方文档，配置](https://grafana.com/docs/loki/latest/configuration/)
> - [官方文档，告警规则和记录规则](https://grafana.com/docs/loki/latest/rules/)

Loki 可以通过两种方式配置 Loki 的运行时行为

- 命令行标志
- 配置文件

配置文件的一部分字段的值，可以通过命令行标志设置。在官方文档中，我们可以查看到配置文件中，所有与命令行标志对应的字段，效果如下：

```yaml
# HTTP server listen host
# CLI flag: -server.http-listen-address
[http_listen_address: <string>]
```

凡是注释中，有 `CLI flag` 的字段，都可以通过命令行标签设置其值。

# Loki 命令行标志详解

**-target <STRING>** # 指定要启用的模块

- 可用的模块有 distributor、ingester、querier、query-frontend、table-manager。
- 可以使用 read、write 来让 loki 运行在只读或只写的模式
- 可以使用 all 表示启用所有模块

# loki.yaml 配置文件详解

文档中包含配置文件关键字与命令行 flag 的对应值，配置文件中的很多配置，都可以通过命令行 flag 来实现。
配置文件的基本结构：

```yaml
# 指定 loki 二进制文件要运行的组件列表。默认值：all，即运行所有组件
target: <STRING>
# 通过 X-Scope-OrgID 标头启用身份验证，如果为 true，则必须存在。 如果为 false，则 OrgID 将始终设置为“ fake”。默认值：true
auth_enabled: <BOOLEAN>
# 用于配置 loki 提供 http 和 gRPC 这两种服务的行为
server: <OBJECT>
# 通用配置。用于配置一些其他配置部分可以共享的配置，比如存储。优先级低，若其他部分指定了相同的配置，则该配置在对应的其他部分的配置将被忽略。
# 从 2.4 版本开始，common 字段将会逐步代替其他描述不清晰的字段，比如 common.storage 将会代替 storage_cofig 字段
common: <OBJECT>

######## 存储架构配置 ########
# 配置储存 Chunk 与 Index 类型数据的模式，以及指定储存这些数据所用的存储类型。
schema_config: <schema_config>
# 为 schema_config 字段指定的存储类型配置详细信息。比如 数据存储位置、连接存储的方式 等等。
# 注意：该字段的配置会根据 schema_config 字段中指定的信息来选择可用的字段。
# 未来将会逐步被 common.storage 字段代替
storage_config: <storage_config>

######## 组件配置 ########
# Distributor(分配器) 组件的配置
distributor: <distributor_config>
# Querier(查询器) 组件的配置.
querier: <querier_config>
# Ingester(摄取器) 组件的配置。还可以配置摄取器如何将自己注册到哈希环上
ingester: <ingester_config>
# 配置 distributor 如何连接到 ingesters
ingester_client: <OBJECT>
# Query Frontend(查询前端) 组件的配置
frontend: <query_frontend_config>
# Ruler(规则器) 组件的配置
ruler: <ruler_config>
# Compactor(压缩器) 组件的配置
compactor: <compactor_config>
# Table Manager(表管理器) 组件的配置，以规定数据保留的行为
table_manager: <OBJECT>

######## 其他配置 ########
# The queryrange_config configures the query splitting and caching in the Loki query-frontend.
query_range: <queryrange_config>
chunk_store_config: <chunk_store_config> # 配置 Loki 如何将数据存放在指定的存储中
# 配置每个租户的限制或全局的限制
limits_config: <limits_config>
# The frontend_worker_config configures the worker - running within the Loki
# querier - picking up and executing queries enqueued by the query-frontend.
frontend_worker: <frontend_worker_config>
# Configuration for "runtime config" module, responsible for reloading runtime configuration file.
runtime_config: <runtime_config>
# Configuration for tracing
tracing: <tracing_config>
```

## target: <STRING> # 指定二进制文件要运行的组件列表。

可用的值有：all、read、write、ingester、distributor、query-frontend、query-scheduler、querier、index-gateway、ruler、compactor。

## auth_enabled: <BOOLEAN>

## server: <Object>

用于配置 loki 提供 http 和 gRPC 这两种服务的行为

```yaml
server:
  http_listen_address: <STRING> # 指定 http 服务监听的端口
```

## [common: <OBJECT>](https://grafana.com/docs/loki/next/configuration/#common)

> 2.4 版本之前并没有这个字段，早期 Loki 的配置文件解读起来非常混乱。但是 2.4 版本之后，可以通过 common 字段统一定义一些之前带有歧义的字段，`common.storage` 可以代替 `storage_config` 用以配置后端存储的信息。

通用配置。**在配置 Loki 组件所使用的 哈希环、存储、等等 时，可以不在每个组件单独配置，而是直接使用这里定义的通用配置。**

### path_prefix: <string>

When defined, the given prefix will be present in front of the endpoint paths.

### replication_factor: <int> | default = 3

How many times incoming data should be replicated to the ingester component.

### ring: <OBJECT>

所有使用哈希环的组件的通用哈希环配置。. If a common ring is given, its values are used to define any undefined ring values. For instance, you can expect the `heartbeat_period` defined in the common section to be used by the distributor's ring, but only if the distributor's ring itself # doesn't have a `heartbeat_period` set.
**kvstore: <OBJECT>** #

- **store: <STRING>** # 用于保存哈希环的存储。`默认值：memberlist`

### storage: <OBJECT>

> 该字段可以代替 `storage_config` 字段。比如 ruler.storage.type 的值为 s3 的话，就会使用这里的 s3 字段的配置；若值为 local，则会使用这里的 filesystem 字段的配置

Loki 不同组件共享使用的存储配置。该字段配置存储信息，用以告诉 Loki 如何使用各种类型的存储。
**s3: <OBJECT> **# S3 类型存储的信息。包括 连接方式、数据要保存的桶 等信息

- 详见下文通用配置字段 [s3: <OBJECT>](#J3m3x)

**azure: \<Azure_Store_Config>** #&#x20;
**gcs: <>** #&#x20;
**swift: <>** #&#x20;
**filesystem: **[**<OBJECT>**](https://grafana.com/docs/loki/next/configuration/#filesystem) # 将本地文件系统作为 Loki 组件存储数据的地方

- **chunks_directory: <STRING> **# 存储 chunks 数据的目录
- **rules_directory: <STRING>** # 存储 Loki Rules 文件的目录

**bos: <OBJECT>** # Baidu Object Storage(百度对象存储) 的信息。
**hedging: **[**<OBJECT>**](https://grafana.com/docs/loki/next/configuration/#hedging) #

### 配置示例

通用的 S3 存储配置

```yaml
common:
  storage:
    s3:
      access_key_id: minioadmin
      bucketnames: chunks
      endpoint: 172.19.42.215:9000
      insecure: true
      s3forcepathstyle: true
      secret_access_key: minioadmin
  ring:
    kvstore:
      store: memberlist
```

## 配置如何存储 chunk 与 index 数据

影响 chunk 与 index 两类数据如何存储的最重要配置只有两个字段：`schema_config` 和 `storage_config`。其他字段都是对存储方式的补充。不过随着版本的更迭，从 2.4 版本开始，`storage_config` 字段会逐渐被 `common.storage` 字段顶替。
简单描述的话：

- schema_config 用来定义使用什么储存数据
- strage_config 用来定义如何连接存储以及存储储存数据的路径。

### schema_config: <Object>

配置存储 chunk 与 index 两类数据的 schema(模式)。该字段用途详见 [Loki 存储](/docs/IT学习笔记/6.可观测性/日志系统/Loki/Storage(存储).md)
schema_config 下只有一个单独的 `configs` 字段，其实用 period_config 更准确~~~`configs` 字段下这是一个数组，每个数组都可以用来定义"某一时间段 loki 存储所使用的 schema"。所以，`configs` 字段用来定义从 哪个时间段开始使用哪种模式将 index 与 chunk 类型的数据存储到哪里去。

#### configs: <\[]Object>

**from: 2018-04-15** # 该模式的起始时间
注意：store 与 object_store 字段的配置将会决定 Loki 使用 storage_config 中的哪个字段作为存储数据的地方
**schema: <STRING>** # 模式的版本，当前推荐为 v11。
**store: <STRING>** # 存放 Index 数据的存储类型。可用的值有：aws, aws-dynamo, gcp, bigtable, bigtable-hashed,cassandra, boltdb-shipper
**object_store: <STRING>** # 存放 Chunks 数据的存储类型。可用的值有：s3、aws、azure、gcp、bigtable、gcs、cassandra、swift、filesystem。`默认值：与 store 字段的值相同`。
**index: <Object>** # 指定储存 Index 数据的行为。

- **prefix: <STRING>** # 表的前缀，也就是 index 文件的前缀。
- **period: <DURATION>** # 表的周期(在当前期间中，每隔 DURATION 的时间创建一张表)。该值必须为 24h 的倍数。`默认值：168h`

**chunks: <Ojbect>** # 指定储存 Chunks 数据的行为。`默认值：复制 index 字段的配置`。其内字段含义与 index 字段下的子字段功能一样。

- **prefix: <STRING>** #
- **period: <DURATION>** #

\~~注意~~~~：~~`~~store~~`~~ 与 ~~`~~object_store~~`~~ 字段的值，将会影响 ~~`~~storage_config~~`~~ 字段下可以使用的字段。比如 store 为 boltdb-shipper，则 storage_config 中只有 boltdb-shipper 字段可以配置，其他无法配置，配置了就会报错。~~Loki 2.4 版本之后，推荐使用 `common.storage` 字段。

### storage_config: <Objcet>

> Loki 2.4 版本之后，推荐使用 `common.storage` 字段。

对 `schema_config` 字段配置的扩充。主要用来定义储存 index 和 chunks 类型数据的存储的行为。比如 连接存储的方式、存储储存数据的位置 等等信息。
有多种存储类型可用，该字段中的配置需要根据 `schema_config.configs.store` 与 `schema_config.configs.object_store` 字段的值来编写。

> 比如，在 schema_config.configs.store 中使用 aws，那么 storage_config 中就可以使用 aws 配置

#### boltdb: <Object> # boltdb 存储类型的配置。

仅当 schema_config.configs.store 为 boltdb 时，才配置该字段

- **directory: <STRING>** # 存放 BoltDB 索引数据的绝对路径

#### boltdb_shipper: <Ojbect> # boltdb_shipper 存储类型的配置。

仅当 schema_config.configs.store 为 boltdb_shipper 时，才配置该字段

- **active_index_directory: <STRING>** #
- **cache_location: <STRING>** #
- **cache_ttl: <DURATION>** # `默认值：24h`
- **shared_store: <STRING>** # 用于保存 BoltDB 文件的存储。
  - 在 2.4 版本之后，若 `common.storage` 定义了 s3，且 `schema_config.object_storage` 定义为 s3，则这个字段的值也为 s3。也就是说，Index 数据也会存到 S3。这个说法待验证。

#### filesystem: <Object> # filesystem 存储类型的配置。

仅当 schema_config.configs.object_store 为 filesystem 时，才配置该字段

- **directory: <STRING>** # 存放 chunks 数据的绝对路径

#### aws: <Object> # S3 配置

仅当 schema_config.configs.object_store 为 aws 时，才配置该字段。该字段配置与通用存储配置中的 `[s3](#zJRSQ)` 字段相同

- **bucketnames: <STRING>** #
- **endpoint: localhost:9000** #
- **access_key_id: minioadmin** #
- **secret_access_key: minioadmin** #
- **insecure: true** #
- **s3forcepathstyle: true** #

## Distributor 组件配置

### distributor: <Object>

Loki 的 distributor(分配器) 组件配置。

## Ingester 组件配置

### ingester_client: <Object>

### ingester: <Object>

Loki 的 Ingester(摄取器) 配置，以及配置采集器如何将自己注册到键值存储
**lifecycler:** #

- **address: 127.0.0.1 **#
- **ring:** #
  - **kvstore:** #
    - **store: <STRING>** # 用于 ring 的后端存储类型。值为 consul, etcd,inmemory, memberlist
  - **replication_factor: 1** #
- **final_sleep: 0s** #

**chunk_idle_period: 5m** #&#x20;
**chunk_retain_period: 30s** #&#x20;
**max_transfer_retries: 0** #&#x20;
**wal: <Object>** # Ingester 的 WAL 配置。

- **enabled: <BOOLEAN>**
- **dir: \</PATH/TO/DIR>** # WAL 存放目录。`默认值: wal`，即默认数据存储目录下的 /wal 目录。

## Querier 组件配置

### querier: <[OBJECT](https://grafana.com/docs/loki/latest/configuration/#querier)>

## Query frontend 组件配置

### frontend: <[OBJECT](https://grafana.com/docs/loki/latest/configuration/#frontend)>

## Ruler 组件配置

### ruler: <Object>

Ruler 组件配置。
**storage: <Ojbect>** # 根据 type 的值，则会优先默认选择[通用存储](#SJMUR)。可用的值有：azure, gcs, s3, swift, local, bos。若没有通用存储，则使用 storage 字段下对应的字段。

- **type: <STRING> **#
- **s3: <OBJECT>** # 配置用于存储规则文件的存储信息
  - 详见下文通用配置字段 [s3: <OBJECT>](#J3m3x)

**rule_path: /loki/tmprules** #&#x20;
**alertmanager_url: http://localhost** #&#x20;
**ring:** #

- **kvstore:** #
  - **store: inmemory** #

### 配置示例

将规则文件保存在本地文件系统

```yaml
ruler:
  alertmanager_url: http://monitor-hw-cloud-k8s-alertmanager.monitoring.svc.cluster.local.:9093
  enable_alertmanager_v2: true
  external_url: https://alertmanager.xx
  ring:
    kvstore:
      store: memberlist
  rule_path: /tmp/loki/scratch
  storage:
    local:
      directory: /etc/loki/rules
    type: local
```

## Table manager 组件配置

### table_manager: <Object>

Table Manager(表管理器) 组件配置，以规定数据保留的行为。该配置环境用途详见《[Loki 存储](/docs/IT学习笔记/6.可观测性/日志系统/Loki/Storage(存储).md)》

> 注意：
>
> - Table Manager 无法管理存放在对象存储(比如 S3)中的数据，如果要使用对象存储来储存 Index 与 Chunks 数据，则应该自行设置 Bucket 的策略，以删除旧数据。

**retention_deletes_enabled: <BOOLEAN>** # 是否开启删除保留数据的行为。`默认值：false`。
**retention_period: <DURATION>** # 指定要保留多长时间的表。

- DURATION 的值必须是 schema_config.configs.index(或 chunks).period 字段值的倍数。`默认值：0s`，即保留所有时间的表，不删除
- 注意，为了避免查询超出保留期限的数据，`chunk_store_config.max_look_back_period` 字段的值必须小于或等于 retention_period 的值

**creation_grace_period: <DURATION>** # 提前 DURATION 时间创建新表。`默认值：10m`

## limits_config: <Object> # 配置各个组件处理数据的最大值

**ingestrion_rate_mb: <FLOAT>** # 每秒可以摄取日志量的大小，单位 MB。`默认值：4`
**enforce_metric_name: <BOOLEAN> **# 强制每个样本都有一个 metric 名称。`默认值：true`

- 通常设为 false

**reject_old_samples: <BOOLEAN> **# 旧样本是否会被拒绝。`默认值：true`
**reject_old_samples_max_age: <DURATION>** # 拒绝前可以接收的最大样本年龄。`默认值：168h`

- 如果拒绝旧样本，那么旧样本不能早于 reject_old_samples_max_age 时间

## 其他

### chunk_store_config: <Object>

配置 Loki 如何将数据存放在指定存储中。该配置环境用途详见《[Loki 存储](/docs/IT学习笔记/6.可观测性/日志系统/Loki/Storage(存储).md)》
**max_look_back_period: <DURATION>** # 限制可以查询多长时间的数据。`默认值：0s`，即不做限制。DURATION 必须小于或等于 table_manager.retention_period 字段的值

# loki.yaml 配置文件中的通用字段

这里面说明的通用字段会被配置文件中的某些字段共同使用。与 common 字段不同，这里指的字段是需要在配置文件中真实书写的；而 common 中定义的配置类似于默认值。

## 后端存储信息

用来定义 如何连接存储、数据在存储中的路径 等等

### [s3: <OBJECT>](https://grafana.com/docs/loki/next/configuration/#s3_storage_config) # S3 存储配置

**endpoint: <STRING>** # 连接 S3 的 endpoint。`默认值：空`
**access_key_id: <STRING>** # 连接 S3 的 AK。`默认值：空`
**secret_access_key: <STRING>** # 连接 S3 的 SK。`默认值：空`
**bucketnames: <STRING>** # 以逗号分割的桶名称列表。`默认值：空`。多个桶可以均匀得分布 chunks
**insecure: <BOOLEAN>** # 是否使用不安全的连接去连接 S3，i.e.是否使用 HTTP 连接 S3。`默认值：false`
**s3forcepathstyle: <BOOLEAN>** #
**http_config: <OBJECT>**

- **insecure_skip_verify: <BOOLEAN>** # 是否跳过证书验证。`默认值：false`

# 配置文件示例

## loki 启动时的最小配置

### 使用本地文件系统

```yaml
auth_enabled: false

server:
  http_listen_port: 3100

common:
  path_prefix: /loki
  storage:
    filesystem:
      chunks_directory: /loki/chunks
      rules_directory: /loki/rules
  replication_factor: 1
  ring:
    instance_addr: 127.0.0.1
    kvstore:
      store: inmemory

schema_config:
  configs:
    - from: 2020-10-24
      store: boltdb-shipper
      object_store: filesystem
      schema: v11
      index:
        prefix: index_
        period: 24h

ruler:
  alertmanager_url: http://localhost:9093
```

### 使用 S3

```yaml
auth_enabled: false

server:
  http_listen_port: 3100

common:
  path_prefix: /loki
  storage:
    s3:
      s3forcepathstyle: true
      bucketnames: loki-lch-test
      endpoint: localhost:9000
      access_key_id: minioadmin
      secret_access_key: minioadmin
      insecure: true
  replication_factor: 1
  ring:
    instance_addr: 127.0.0.1
    kvstore:
      store: inmemory

schema_config:
  configs:
    - from: 2020-10-24
      store: boltdb-shipper
      object_store: s3
      schema: v11
      index:
        prefix: index_
        period: 24h

ruler:
  alertmanager_url: http://localhost:9093
```

## Index 与 Chunk 都使用 S3

这里的 S3 使用 Mini

```yaml
schema_config:
  configs:
    - from: 2020-07-01
      store: boltdb-shipper
      object_store: aws
      schema: v11
      index:
        prefix: index_
        period: 24h
common:
  storage:
    s3:
      access_key_id: minioadmin
      bucketnames: chunks
      endpoint: 172.19.42.215:9000
      insecure: true
      s3forcepathstyle: true
      secret_access_key: minioadmin
```

## 简单完整配置

```yaml
auth_enabled: false
chunk_store_config:
  max_look_back_period: 0s
common:
  ring:
    kvstore:
      store: memberlist
  storage:
    s3:
      access_key_id: minioadmin
      bucketnames: chunks
      endpoint: 172.19.42.215:9000
      insecure: true
      s3forcepathstyle: true
      secret_access_key: minioadmin
compactor:
  shared_store: filesystem
distributor:
  ring:
    kvstore:
      store: memberlist
frontend:
  compress_responses: true
  log_queries_longer_than: 5s
  tail_proxy_url: http://loki-loki-distributed-querier:3100
frontend_worker:
  frontend_address: loki-loki-distributed-query-frontend:9095
ingester:
  chunk_block_size: 262144
  chunk_encoding: snappy
  chunk_idle_period: 1h
  chunk_retain_period: 1m
  chunk_target_size: 1536000
  lifecycler:
    ring:
      kvstore:
        store: memberlist
      replication_factor: 1
  max_chunk_age: 1h
  max_transfer_retries: 0
  wal:
    dir: /var/loki/wal
limits_config:
  enforce_metric_name: false
  max_cache_freshness_per_query: 10m
  reject_old_samples: true
  reject_old_samples_max_age: 168h
  split_queries_by_interval: 15m
memberlist:
  join_members:
    - loki-loki-distributed-memberlist
query_range:
  align_queries_with_step: true
  cache_results: true
  max_retries: 5
  results_cache:
    cache:
      enable_fifocache: true
      fifocache:
        max_size_items: 1024
        validity: 24h
ruler:
  alertmanager_url: https://alertmanager.xx:9093
  enable_alertmanager_v2: true
  external_url: https://alertmanager.xx
  ring:
    kvstore:
      store: memberlist
  rule_path: /tmp/loki/scratch
  storage:
    local:
      directory: /etc/loki/rules
    type: local
schema_config:
  configs:
    - from: "2022-06-21"
      index:
        period: 24h
        prefix: loki_index_
      object_store: s3
      schema: v12
      store: boltdb-shipper
server:
  http_listen_port: 3100
storage_config:
  boltdb_shipper:
    active_index_directory: /var/loki/index
    cache_location: /var/loki/cache
    cache_ttl: 168h
    shared_store: s3
  filesystem:
    directory: /var/loki/chunks
table_manager:
  retention_deletes_enabled: false
  retention_period: 0s
```
