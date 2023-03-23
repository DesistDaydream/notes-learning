---
title: 使用 Thanos 查询前端优化查询性能
---

Thanos 中的 Query 组件可以提供一个统一的查询入口，但是当查询的数据规模较大的时候，对 querier 组件也会有很大的压力，为此 Thanos 也提供了一个 `Query Frontend` 的组件来提升性能。Thanos Query Frontend 是 Thanos Query 的前端，它的目标是将大型查询拆分为多个较小的查询，并缓存查询结果来提升性能。

## 概述

Thanos Query Frontend 组件通过 `thanos query-frontend` 命令实现了一个放在 querier 前面的服务，以改进读取路径。它基于 Cortex Query Frontend 组件，所以你可以找到一些 Cortex 常见的特性，如查询拆分和结果缓存。Thanos Query Frontend 是无状态和水平扩展的，可以使用下列命令来启动 Thanos Query Frontend：

`thanos query-frontend \     --http-address     "0.0.0.0:9090" \     --query-frontend.downstream-url="<thanos-querier>:<querier-http-port>"`

在接收到查询请求后 query frontend 不会立即响应，而是将查询请求放入一个查询队列中，querier 会连接到 query frontend 并消费这个队列，执行从队列中获取的查询请求并响应给 query frontend，query frontend 会对这些响应的结果进行聚合。

## 特性

### 查询队列

query frontend 的队列机制有以下用途。

-

确保可能导致 OOM 的大型查询在发生错误时能够得到重试。

-

防止多个大的查询请求打在单个 querier 上。

-

可以分配租户所对应的 querier，避免单个租户使用 DOS 拒绝服务攻击其他租户。

### 查询拆分

query frontend 会将多天的的查询拆分为多个单天的查询，游下游的 querier 去并行处理这些已拆分的查询。返回的查询结果由 query frontend 进行汇聚。这样可以防止大时间跨度的查询导致 queier 发生 OOM，并且能够更快的执行查询以及更好的查询负载均衡。查询前端根据配置的 `--query-range.split-interval` 标志将长查询拆分为多个短查询，`--query-range.split-interval` 的默认值为 24 小时，启用缓存时，它应该大于 0。

### 重试机制

query frontend 支持在 HTTP 请求失败时重试查询的重试机制，有一个 `--query-range.max-retries-per-request` 标志来限制最大重试次数。

### 查询缓存

query frontend 支持将查询结果进行缓存用以加速后续的查询。当缓存的结果不够完整时，query frontend 会计算出所需要的子查询并分配给下游的 querier 并行执行，子查询的步长会对齐以提升查询结果的可缓存性。当前支持的缓存后端有：memcached，redis 和内存，下面是这些缓存后端的一些配置。

**内存缓存**

`type: IN-MEMORY config:   max_size: ""   max_size_items: 0   validity: 0s`

max_size 表示在内存中缓存的最大尺寸，单位可以是 KB、MB、GB。如果 max_size 和 max_size_items 都没有设置，就不会创建缓存。如果只设置 max_size 或 max_size_items 中的任意一个，则对其他字段没有限制。

**memcached**

`type: MEMCACHED config:   addresses: [your-memcached-addresses]   timeout: 500ms   max_idle_connections: 100   max_item_size: 1MiB   max_async_concurrency: 10   max_async_buffer_size: 10000   max_get_multi_concurrency: 100   max_get_multi_batch_size: 0   dns_provider_update_interval: 10s   expiration: 24h`

expiration 表示指定 memcached 缓存有效时间，如果设置为 0，则使用默认的 24 小时过期时间，上面的配置是默认的 memcached 配置。

**Redis**

默认的 Redis 配置如下：

`type: REDIS config:   addr: ""   username: ""   password: ""   db: 0   dial_timeout: 5s   read_timeout: 3s   write_timeout: 3s   pool_size: 100   min_idle_conns: 10   idle_timeout: 5m0s   max_conn_age: 0s   max_get_multi_concurrency: 100   get_multi_batch_size: 100   max_set_multi_concurrency: 100   set_multi_batch_size: 100   expiration: 24h0m0s`

expiration 表示指定 Redis 缓存有效时间，如果设置为 0，则使用默认的 24 小时过期时间。

### 慢查询日志

query frontend 还支持通过配置 `--query-frontend.log-queries-longer-than` 标志来记录运行时间超过某个持续时间的查询。

## 安装

安装 query frontend 最重要的就是通过 `-query-frontend.downstream-url` 指定下游的 querier，如下所示：

\`# thanos-query-frontend.yamlapiVersion: apps/v1
kind: Deployment
metadata:
  name: thanos-query-frontend
  namespace: kube-mon
  labels:
    app: thanos-query-frontend
spec:
  selector:
    matchLabels:
      app: thanos-query-frontend
  template:
    metadata:
      labels:
        app: thanos-query-frontend
    spec:
      containers:
        - name: thanos
          image: thanosio/thanos:v0.25.1
          imagePullPolicy: IfNotPresent
          ports:
            - containerPort: 9090
              name: http
          args:
            - query-frontend
            - --log.level=info
            - --log.format=logfmt
            - --query-frontend.compress-responses
            - --http-address=0.0.0.0:9090
            - --query-frontend.downstream-url=http://thanos-querier.kube-mon.svc.cluster.local:9090
            - --query-range.split-interval=12h
            - --query-range.max-retries-per-request=10
            - --query-frontend.log-queries-longer-than=10s
            - --labels.split-interval=12h
            - --labels.max-retries-per-request=10
            - |-
              --query-range.response-cache-config="config":
                max_size: "200MB"
                max_size_items: 0
                validity: 0s
              type: IN-MEMORY
            - |-
              --labels.response-cache-config="config":
                max_size: "200MB"
                max_size_items: 0
                validity: 0s
              type: IN-MEMORY
          env:
            - name: HOST_IP_ADDRESS
              valueFrom:
                fieldRef:
                  fieldPath: status.hostIP
          livenessProbe:
            failureThreshold: 4
            httpGet:
              path: /-/healthy
              port: 9090
              scheme: HTTP
            periodSeconds: 30
          readinessProbe:
            failureThreshold: 20
            httpGet:
              path: /-/ready
              port: 9090
              scheme: HTTP
            periodSeconds: 5
          resources:
            requests:
              memory: 512Mi
              cpu: 500m
            limits:
              memory: 512Mi
              cpu: 500m

---

apiVersion: v1

kind: Service

metadata:

name: thanos-query-frontend

namespace: kube-mon

labels:

app: thanos-query-frontend

spec:

ports:

- name: http

port: 9090

targetPort: 9090

selector:

app: thanos-query-frontend

type: NodePort

\`

这里我们开启了缓存，使用内存缓存，直接通过 `-query-range.response-cache-config` 参数来配置缓存配置，也可以通过 `-query-range.response-cache-config-file` 指定缓存配置文件，两种方式均可，为了验证结果，这里创建了一个 NodePort 类型的 Service，直接创建上面的资源对象即可：

`☸ ➜ kubectl apply -f https://p8s.io/docs/thanos/manifests/thanos-query-frontend.yaml ☸ ➜ kubectl get pods -n kube-mon -l app=thanos-query-frontend NAME                                     READY   STATUS    RESTARTS   AGE thanos-query-frontend-78954bc857-nkxkk   1/1     Running   0          151m ☸ ➜ kubectl get svc -n kube-mon -l app=thanos-query-frontend NAME                    TYPE       CLUSTER-IP      EXTERNAL-IP   PORT(S)          AGE thanos-query-frontend   NodePort   10.97.182.150   <none>        9090:30514/TCP   152m`

然后可以通过任意节点的 30514 端口访问到查询前端：

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/eb1b6218-460e-40ab-8fdc-0d2e104245d9/640)

该前端页面和 Query 组件是一致的，但是在 Query 前面缓存包装了一层。
