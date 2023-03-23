---
title: RabbitMQ 集群
---

# 概述

> 参考：
> - 官方文档：<https://www.rabbitmq.com/clustering.html>

RabbitMQ 这款消息队列中间件产品本身是基于 Erlang 语言编写，Erlang 语言天生具备分布式特性（通过同步 Erlang 集群各节点的 magic cookie 来实现）。

因此，RabbitMQ 天然支持 Clustering。这使得 RabbitMQ 本身不需要像 ActiveMQ、Kafka 那样通过 ZooKeeper 分别来实现 HA 方案和保存集群的元数据。集群是保证可靠性的一种方式，同时可以通过水平扩展以达到增加消息吞吐量能力的目的。下面先来看下 RabbitMQ 集群的整体方案：

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/iqase5/1616130673298-8615562e-d09c-44b6-be6d-1db0294d30e6.jpeg)

上面图中采用三个节点组成了一个 RabbitMQ 的集群，Exchange A 的元数据信息在所有节点上是一致的，而 Queue 的完整数据则只会存在于它所创建的那个节点上。其他节点只知道这个 queue 的 metadata 信息和一个指向 queue 的 owner node 的指针。

## RabbitMQ 集群中各节点关系

RabbitMQ 集群各节点同步的数据

RabbitMQ 集群的各节点会始终同步四种类型的内部元数据（类似索引）

1. 队列元数据：队列名称和它的属性
2. 交换器元数据：交换器名称、类型和属性
3. 绑定元数据：一张简单的表格展示了如何将消息路由到队列
4. vhost 元数据：为 vhost 内的队列、交换器和绑定提供命名空间和安全属性

因此，当用户访问其中任何一个 RabbitMQ 节点时，通过 rabbitmqctl 查询到的 queue/user/exchange/vhost 等信息都是相同的。

默认情况下，队列中的数据(消息) 是不在各节点互相同步的。如果想要各节点数据保持一致，查看 《RabbitMQ 基于集群的高可用性》章节

Nodes Equal Peers(节点对等)

在某些分布式系统中，节点是具有领导者和追随者概念的。对于 RabbitMQ，通常情况并不是这样的。RabbitMQ 集群中所有节点都是 equal peers(对等的)。

RabbitMQ 集群发送/订阅消息的基本原理

RabbitMQ 集群的工作原理图如下：

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/iqase5/1616130673314-38d3cddf-87ca-4f97-88aa-523c76136332.jpeg)

场景 1：客户端直接连接队列所在节点

如果有一个生产者或者消费者通过 amqp-client 的客户端连接至节点 1 进行消息的发布或者订阅，那么此时的集群中的消息收发只与节点 1 相关，这个没有任何问题；如果客户端相连的是节点 2 或者节点 3（队列 1 数据不在该节点上），那么情况又会是怎么样呢？

场景 2：客户端连接的是非队列数据所在节点

如果消息生产者所连接的是节点 2 或者节点 3，此时队列 1 的完整数据不在该两个节点上，那么在发送消息过程中这两个节点主要起了一个路由转发作用，根据这两个节点上的元数据（也就是上文提到的：指向 queue 的 owner node 的指针）转发至节点 1 上，最终发送的消息还是会存储至节点 1 的队列 1 上。

同样，如果消息消费者所连接的节点 2 或者节点 3，那这两个节点也会作为路由节点起到转发作用，将会从节点 1 的队列 1 中拉取消息进行消费。

# RabbitMQ 基于集群的高可用性

默认情况，RabbitMQ 集群模式下，队列的数据(也就是消息)只会留在该队列所在节点上，并不在所有节点互相同步数据。比如我在 node1 创建了 queue1，则 queue1 中的消息，则只存在于 node1 上，这是为了提高 RabbitMQ 的性能，可以将压力分担在集群中每个节点上。

那么如果想要让集群中的节点复制队列中的数据，以实现高可用效果，RabbitMQ 提供了两种方式实现：

1. Classic Mirrored Queues(传统镜像队列) # 不会随节点故障而转移
2. Quorum Queues(仲裁队列) # 节点故障后，队列自动转移到其他节点

持久化：

1. Durable # 节点故障后，队列处于 down 状态。即 开启持久化
2. Transient # 节点故障后，队列自动消失。即 关闭持久化

Note：Quorum 类型队列的持久化必须开启。

## Classic Mirrored Queues # 不会随节点故障而转移

官方文档：<https://www.rabbitmq.com/ha.html>

Quorum Queues(仲裁队列) 是另一种更现代的队列类型，通过复制提供高可用性，并关注数据安全。在很多情况下，仲裁队列是比传统镜像队列更好的选择

## Quorum Queues # 节点故障后，队列自动转移到其他节点。(该类型从 3.8.0 版本起可用)

官方文档：<https://www.rabbitmq.com/quorum-queues.html>

仲裁队列是 RabbitMQ 实现高可用的队列类型，基于 RAFT 共识算法 ，从 RabbitMQ 3.8.0 版本开始作为默认推荐的高可用方法。

Quorum 队列中的数据在个节点保持一致，当其中一个节点故障时，队列会自动转移到其他节点上。

## Classic Mirrored 与 Quorum 的区别

永昌，在下列情况下，应使用仲裁队列

1. 数据可用性需求较高，比如销售系统中的新订单或选举系统中的投票，这些消息丢失会造成功能产生重大影响

通常，在下列情况下，不应使用仲裁队列

1. 数据可用性需求不高的场景。比如股票行情和即使通讯系统
