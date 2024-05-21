---
title: Thanos 组件详解
weight: 1
---

# 概述

> 参考：
> 
> - [官方文档，组件](https://thanos.io/tip/components/)
> - <https://zhuanlan.zhihu.com/p/137248127>

# Compactor(压实器)

**注意：Compactor 在持久运行状态，会对对象存储发起大量的 GET 请求。最好间隔一段时间，运行一次，压缩一次数据即可，不必持久运行**

![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/ilh4m6/1636990514133-53c917f1-6f57-4bb5-a59d-0776a5fef235.png)

# Receiver(接收器)

> 参考：
> 
> - [官方文档,组件-接收器](https://thanos.io/tip/components/receive.md)

# Querier(查询器)

> 参考：
> 
> - [官方文档,组件-查询器](https://thanos.io/tip/components/query.md)

Querier 组件分为两部分

- **Querier(查询器)** # 实现了 Prometheus API，可以通过 Querier 发起 PromQL 查询请求，以获取数据；甚至可以从 Prometheus Server 的时序数据库中删除数据。每个从 Querier 发起的 PromQL 查询请求都会发送到可以暴露 StoreAPI 的组件上，并获取查询结果。
- **Query Fronted(查询前端)** # 实现了 Prometheus API，可以将请求负载均衡到指定的多个 Querier 上，同时可以缓存响应数据、也可以按查询日拆分。有点像 Redis 的效果

Querier 组件向一个或多个暴露 StoreAPI 的组件发起查询请求，并将结果去重后，返回给查询客户端。
![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/ilh4m6/1627118034684-9a4e56ce-7593-4bca-bc38-7fd72bff563c.png)

## Deduplication(重复数据删除)

The query layer can deduplicate series that were collected from high-availability pairs of data sources such as Prometheus. A fixed single or multiple replica labels must be chosen for the entire cluster and can then be passed to query nodes on startup.查询层可以从高可用性对数据源（如 Prometheus）收集的档案。必须为整个群集选择固定的单个或多个副本标签，然后可以通过启动时传递到查询节点。

Two or more series that are only distinguished by the given replica label, will be merged into a single time series. This also hides gaps in collection of a single data source.仅通过给定的副本标签区分的两个或更多系列，将合并为单个时间序列。这也隐藏了收集单个数据源的间隙。

### 单副本标签示例

假如现在 Query 从不同的 StoreAPI 中获取了三条时序数据

- Prometheus + sidecar "A": cluster=1,env=2,replica=A
- Prometheus + sidecar "B": cluster=1,env=2,replica=B
- Prometheus + sidecar "A" in different cluster: cluster=2,env=2,replica=A

我们像下面这样配置 query

```bash
If we configure Querier like this:

thanos query \
    --http-address        "0.0.0.0:9090" \
    --query.replica-label "replica" \
    --store               "<store-api>:<grpc-port>" \
    --store               "<store-api2>:<grpc-port>" \
```

当我们在运行 query 时指定了 `--query.replica-label` 标志时，我们会将具有相同标签的时序去重只保留一个，此时我们将获得 2 个结果：

- up{job="prometheus",env="2",cluster="1"} 1
- up{job="prometheus",env="2",cluster="2"} 1

如果没有此副本标志（关闭数据删除），我们将获得 3 个结果：

- up{job="prometheus",env="2",cluster="1",replica="A"} 1
- up{job="prometheus",env="2",cluster="1",replica="B"} 1
- up{job="prometheus",env="2",cluster="2",replica="A"} 1

### 多副本标签示例

- Prometheus + sidecar "A": cluster=1,env=2,replica=A,replicaX=A
- Prometheus + sidecar "B": cluster=1,env=2,replica=B,replicaX=B
- Prometheus + sidecar "A" in different cluster: cluster=2,env=2,replica=A,replicaX=A

上面的例子有两个需要去重的标签，那么就可以使用两次 `--query.replica-label` 标志。

```bash
thanos query \
    --http-address        "0.0.0.0:9090" \
    --query.replica-label "replica" \
    --query.replica-label "replicaX" \
    --store               "<store-api>:<grpc-port>" \
    --store               "<store-api2>:<grpc-port>" \
```

# Query Frontend(查询前端)

> 参考：
> 
> - [官方文档，组件-查询前端](https://thanos.io/tip/components/query-frontend.md/)
> - [公众号-k8s 技术圈，使用 Thanos 查询前端优化查询性能](https://mp.weixin.qq.com/s/W9diP0OKt_-ajAXM_wgogg)

# Sidecar(边车)

> 参考：
> 
> - [官方文档，组件-边车](https://thanos.io/tip/components/sidecar.md/)
