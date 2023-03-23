---
title: Loki 衍生品
---

# 概述

> 参考：
> - [GitHub 项目,production](https://github.com/grafana/loki/tree/main/production)
> - [公众号,Loki 生产环境集群方案](https://mp.weixin.qq.com/s/qnt7JUzHLUU6Qs_tv5V0Hw)

很多新入坑 Loki 的小伙伴当看到 distributor、ingester、querier 以及各种依赖的三方存储时，往往都比较懵逼，不知道从哪儿入手。此外再加上官方的文档里面对于集群部署的粗浅描述，更是让新手们大呼部署太难。其实，除了官方的 helm 外，藏在 Loki 仓库的 production 目录里面有一篇生产环境的集群部署模式。
![](https://notes-learning.oss-cn-beijing.aliyuncs.com/pwhqog/1621302640599-456c53ed-2d2f-445d-af56-43a28becb54f.webp)
原文里面，社区采用的是 docker-compose 的方式来快速拉起一套 Loki 集群。虽然我们正式在生产环境中实施时，不会傻到用 docker-compose 部署在一个 node 上（显然这里我们强行不考虑 docker-swarm）。不过里面关于 Loki 的架构和配置文件却值得我们学习。
那么，与纯分布式的 Loki 集群相比，这套方案有什么特别的呢？首先我们先来看看下面这张图：
![](https://notes-learning.oss-cn-beijing.aliyuncs.com/pwhqog/1621302640583-228525ed-be05-4c2c-90ec-329fc40ed955.webp)
可以看到，最明显的有三大不同点：

1. loki 核心服务 distributor、ingester、querier 没有分离，而是启动在一个实例当中；
2. 抛弃了 consul 和 etcd 外部的 kv 存储，而是直接用 memberlist 在内存中维护集群状态；
3. 使用 boltdb-shipper 替代其他日志索引方案

这样看起来，Loki 集群的整体架构就比较清晰，且更少的依赖外部系统。简单总结了下，除了用于存储 chunks 和 index 而绕不开的 S3 存储外，还需要一个缓存服务用于加速日志查询和写入。

> Loki2.0 版本之后，对于使用 boltdb 存储索引部分做了较大的重构，采用新的 boltdb-shipper 模式，可以让 Loki 的索引存储在 S3 上，而彻底摆脱 Cassandra 或者谷歌的 BigTable。此后服务的横向扩展将变得更加容易。关于 bolt-shipper 的更多细节，可以参考：<https://grafana.com/docs/loki/latest/operations/storage/boltdb-shipper/>

说得这么玄乎，那我们来看看这套方案的配置有哪些不一样呢？

## 原生部分

#### memberlist

    memberlist:
      join_members: ["loki-1", "loki-2", "loki-3"]
      dead_node_reclaim_time: 30s
      gossip_to_dead_nodes_time: 15s
      left_ingesters_timeout: 30s
      bind_addr: ['0.0.0.0']
      bind_port: 7946

Loki 的 memberlist 使用的是 gossip 协议来让集群内的所有节点达到最终一致性的。此部分的配置几乎都是协议频率和超时的控制，保持默认的就好

#### ingester

    ingester:
      lifecycler:
        join_after: 60s
        observe_period: 5s
        ring:
          replication_factor: 2
          kvstore:
            store: memberlist
        final_sleep: 0s

ingester 的状态通过 gossip 协议同步到集群的所有 member 当中，同时让 ingester 的复制因子为 2。即一个日志流同时写入到两个 ingster 服务当中以保证数据的冗余。

## 扩展部分

社区的集群模式配置原生部分仍然显得不太够意思，除了 memberlist 的配置稍显诚意外，其它部分仍然不够我们对生产环境的要求。这里小白简单改造了一下，分享给大家。

#### storage

将 index 和 chunks 的存储统一让 S3 对象存储纳管，让 Loki 彻底摆脱三方依赖。

    schema_config:
      configs:
      - from: 2021-04-25
        store: boltdb-shipper
        object_store: aws
        schema: v11
        index:
          prefix: index_
          period: 24h

    storage_config:
      boltdb_shipper:
        shared_store: aws
        active_index_directory: /loki/index
        cache_location: /loki/boltdb-cache
      aws:
        s3: s3://<S3_ACCESS_KEY>:<S3_SECRET_KEY>@<S3_URL>/<S3__BUCKET>
        s3forcepathstyle: true
        insecure: true

这里值得说明的就是用于存储日志流索引的是 bolt_shipper，它是可以通过共享存储方式写到 s3 当中的。那么`active_index_directory`就是 S3 上的 Bucket 路径，`cache_location`则为 Loki 本地 bolt 索引的缓存数据。

> 事实上 ingester 上传到 s3 的 index 路径为`<S3__BUCKET>/index/`

#### redis

原生的方案里并不提供缓存，这里我们引入 redis 来做查询和写入的缓存。对于很多小伙伴纠结的是一个 redis 共用还是多个 redis 单独使用，这个看你集群规模，不大的情况下，一个 redis 实例足以满足需求。

    query_range:
      results_cache:
        cache:
          redis:
            endpoint: redis:6379
            expiration: 1h
      cache_results: true
    index_queries_cache_config:
      redis:
        endpoint: redis:6379
        expiration: 1h

    chunk_store_config:
      chunk_cache_config:
        redis:
          endpoint: redis:6379
          expiration: 1h
      write_dedupe_cache_config:
        redis:
          endpoint: redis:6379
          expiration: 1h

#### ruler

既然 Loki 以及做了集群化部署，当然 ruler 这个服务也得跟在切分。难以接受的是，社区这部分的配置竟然是缺失的。所以我们得自己补充完整。我们知道日志的 ruler 可以写在 S3 对象存储上，同时每个 ruler 实例也是通过一致性哈希环来分配自己的 rules。所以这部分配置，我们可以如下参考：

    ruler:
      storage:
        type: s3
        s3:
          s3: s3://<S3_ACCESS_KEY>:<S3_SECRET_KEY>@<S3_URL>/<S3_RULES_BUCKET>
          s3forcepathstyle: true
          insecure: true
          http_config:
            insecure_skip_verify: true
        enable_api: true
        enable_alertmanager_v2: true
        alertmanager_url: "http://<alertmanager>"
        ring:
          kvstore:
          store: memberlist

#### 支持 kubernetes

最后，最最最重要的是要让官方的 Loki 集群方案支持在 Kubernetes 中部署，否则一切都是瞎扯。由于篇幅的限制，我将 manifest 提交到 github 上，大家直接 clone 到本地部署。
GitHub 地址：<https://github.com/CloudXiaobai/loki-cluster-deploy/tree/master/production/loki-system>
![image.gif](https://notes-learning.oss-cn-beijing.aliyuncs.com/pwhqog/1621302640377-bc3f0b0d-6c1a-4d28-a2a3-9aa99aa34e4e.gif)
这个 manifest 只依赖一个 S3 对象存储，所以你在部署到生产环境时，请务必预先准备好对象存储的 AccessKey 和 SecretKey。将他们配置到 installation.sh 当中后，直接执行脚本就可以开始安装了。

> 文件中的 ServiceMonitor 是为 Loki 做的 Prometheus Operator 的 Metrics 服务发现，你可以自己选择是否部署

## 总结

本文介绍了官方提供的一种 Loki 生产环境下的集群部署方案，并在此基础上加入了一些诸如缓存、S3 对象存储的扩展配置，并将官方的 docker-compose 部署方式适配到 Kubernetes 当中。官方提供的方案有效的精简了 Loki 分布式部署下复杂的结构，值得我们学习。

**_![image.gif](https://notes-learning.oss-cn-beijing.aliyuncs.com/pwhqog/1621302640324-6122ec11-5c55-4805-b893-9b816cdaa09f.gif)_**
