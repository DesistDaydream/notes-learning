---
title: Tempo
linkTitle: Tempo
weight: 20
---

# 概述

> 参考：
>
> - 原文链接: https://zhuanlan.zhihu.com/p/272506092

`Tempo`是 Grafana Labs 在`ObservabilityCON 2020`大会上新开源的一个用于做分布式式追踪的后端服务。它和 Cortex、Loki 一样，Tempo 也是一个兼备`高扩展`和`低成本`效应的系统。

之前小白有提到 Grafana Labs 的云原生 Observability 宇宙只剩下 trace 部分，那么今天就拿 Loki 的分布式追踪来体验下这 Observability 的最后一环吧。正式开始前，先看下小白精心准备的 Tempo 体验视频吧。

## 关于 Tempo

Tempo 本质上来说还是一个存储系统，它兼容一些开源的 trace 协议（包含 Jaeger、Zipkin 和 OpenCensus 等），将他们存在廉价的 S3 存储中，并利用 TraceID 与其他监控系统（比如 Loki、Prometheus）进行协同工作。

![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/tempo/202311041716613.png)

可以看到 Tempo 的架构仍然分为`distributor`、`ingester`、`querier`、`tempo-query`、`compactor`这几个架构，熟悉 Loki 和 Cortex 的朋友可能光看名字就知道他们大概是做什么的。不熟悉的同学也没关系，下面简单说下各模块的作用：

- distributor

监听多个端口，分别接受来自 Jaeger、Zipkin 和 OpenCensus 协议的数据，按照 TraceID 进行哈希并映射到哈希环上，并交由 ingester 进行存储处理。当前 distributor 支持的 trace 协议如下：

| Protocol                | Port  |
| ----------------------- | ----- |
| OpenTelemetry           | 55680 |
| Jaeger - Thrift Compact | 6831  |
| Jaeger - Thrift Binary  | 6832  |
| Jaeger - Thrift HTTP    | 14268 |
| Jaeger - GRPC           | 14250 |
| Zipkin                  | 9411  |

- ingester

具体负责 trace 数据的块存储（memcache、GCS、S3）、缓存（Memcache）和索引的处理

- querier

负责从 ingester 和后端存储里面捞取 trace 数据，并提供 api 给查询者

- compactor

负责后端存储块的压缩，减少数据块数量

- tempo-query

tempo 的一个可视化界面，用的`jaeger query`，可以在上面查询 tempo 的 trace 数据。

## Loki 链路跟踪

> 要体验的同学，可以先下载小白在 GitHub 上的 Docker-Compose，推荐配合本篇内容一起实践
>
> <https://github.com/CloudXiaobai/loki-cluster-deploy/tree/master/demo/docker-compose-with-tempo>

### Loki 方面

在做之前我们先看下 Loki 的文档是怎么描述的:

_The tracing_config block configures tracing for Jaeger. Currently limited to disable auto-configuration per environment variables only._

可以看到当前 Loki 对于 Trace 的支持集中在 Jaeger，而且配置是默认开启的，并且只能在环境变量里面读取 jaeger 的信息。docker-compose 下的案例如下：

    querier-frontend:  image: grafana/loki:1.6.1  runtime: runc  scale: 2  environment:    - JAEGER_AGENT_HOST=tempo    \\tempo的地址    - JAEGER_ENDPOINT=http://tempo:14268/api/traces    - JAEGER_SAMPLER_TYPE=const   \\采样率类型    - JAEGER_SAMPLER_PARAM=100    \\采样率100

### API 网关方面

API 网关并不是 Loki 的原生组件，而是在 Loki 分布式部署的情况下，需要有一个统一的入口对 Loki API 进行路由。之前小白用的 Nginx，但是`原生的Nginx并不支持OpenTracing`。小白根据 nginx1.14 版本做了一个带 jaeger 模块的镜像用于 Loki 入口的 trace 生成和日志采集。

    gateway:  image: quay.io/cloudxiaobai/nginx-jaeger:1.14.0  runtime: runc  restart: always  ports:    - 3100:3100  volumes:    - ./nginx.conf:/etc/nginx/nginx.conf    - ./jaeger-config.json:/etc/jaeger-config.json    - 'gateway_trace_log:/var/log/nginx/'

对于支持 OpenTracing 的 Nginx，我们需要修改 nginx.conf 配置文件如下：

    ...#加载opentracing库load_module modules/ngx_http_opentracing_module.so;http {   #启用opentracing  opentracing on;   #加载jaeger库  opentracing_load_tracer /usr/local/lib/libjaegertracing_plugin.so /etc/jaeger-config.json;   #日志格式，打印traceid  log_format opentracing '"traceID":"$opentracing_context_uber_trace_id"';   server {    listen               3100 default_server;    location = / {      #向upstream转发时带上trace的头信息      opentracing_operation_name $uri;      opentracing_trace_locations off;      opentracing_propagate_context;      proxy_pass      http://querier:3100/ready;    }  }}

> 以上小白只截取了 Nginx 部分配置，完整的要参考 docker-compose 里的 nginx.conf

此外，nginx 还需要一个 jaeger-config.json，用于将 trace 数据转给 agent 处理。

    {  "service_name": "gateway", \\服务名  "diabled": false,  "reporter": {    "logSpans": true,    "localAgentHostPort": "jaeger-agent:6831"  \\jaeger-agent地址  },  "sampler": {    "type": "const",    "param": "100"  \\采样率  }}

> 为了方便演示，小白配置的采样率均为 100%

最后，我们为 API 网关启用一个 Jaeger-agent 用于收集 trace 信息并转给 Tempo，它的配置如下：

    jaeger-agent:  image: jaegertracing/jaeger-agent:1.20  runtime: runc  restart: always  # 转发给tempo  command: ["--reporter.grpc.host-port=tempo:14250"]  ports:    - "5775:5775/udp"    - "6831:6831/udp"    - "6832:6832/udp"    - "5778:5778"

> 为什么 API 网关不直接发给 Tempo 要经过 Jaeger-agent 转发一下，小白认为用 agent 的方式更加灵活一些。

以上，我们就完成了 Loki 分布式追踪的配置部分，接下来我们用`docker-compose up -d`将服务都运行起来。

### Grafana 方面

当 docker 的所有服务运行正常后，我们访问 grafana 并添加两个数据源

- 添加 tempo 数据源

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/243bc194-36a1-4ab1-b90f-2f1f4062a53e/640)

- 添加 Loki 数据源，并解析 API 网关 TraceID

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/243bc194-36a1-4ab1-b90f-2f1f4062a53e/640)
![](https://notes-learning.oss-cn-beijing.aliyuncs.com/243bc194-36a1-4ab1-b90f-2f1f4062a53e/640)

> Loki 提取 TraceID 的正则部分是从 API 网关的日志中匹配

## 体验 Tempo

数据源设置 OK 后，我们进入 Explore 选择 loki 查询 trace.log 就可以得到 API 网关的日志了。

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/243bc194-36a1-4ab1-b90f-2f1f4062a53e/640)
从 Parsed Fields 里面我们就可以看到，Grafana 从 API 网关的日志里面提取了 16 位字符串作为 TraceID 了，而它关联了 Tempo 的数据源，我们点击`Tempo`按钮就可以直接切到 Trace 的信息如下：

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/243bc194-36a1-4ab1-b90f-2f1f4062a53e/640)

展开 Trace 信息，我们可以看到 Loki 的一次查询的链路会经过下面几个部分

    gateway -> query-frontend -> querier -> ingester                                    |-> SeriesStore.GetChunkRefs

并且得出结论，本次查询的耗时主要落在 Ingeter 上，原因是查询的日志还没被 flush 到存储当中，querier 需从 ingester 中取日志的数据。

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/243bc194-36a1-4ab1-b90f-2f1f4062a53e/640)

我们再来看一个 Loki 接收日志的案例：

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/243bc194-36a1-4ab1-b90f-2f1f4062a53e/640)

从 trace 的链路来看，当日志采集端往 Loki Post 日志时，请求的链路会经过如下部分：

    gateway -> distributor -> ingester

同时，我们还看到了这次的提交的日志流经过两个 ingester 实例的处理，且处理时间没有明显差异。

## 总结

关于`Logging`和`Tracing`两部分在 Grafana 上的展示还没有达到 ObservabilityCON 2020 上的流畅度，不过根据会上的消息，更精细话的`trace <--> log`、`metrics <--> trace`和`metrics <--> log`这三部分互相协作部分应该很快会发布。届时 Grafana 将是云上可观测性应用系统里的王者级产品（虽然有额外的各种查询语句学习成本）
