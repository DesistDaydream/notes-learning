---
title: Kafka
linkTitle: Kafka
date: 2024-05-08T20:14
weight: 20
---

# 概述

> 参考：
>
> - [官网](https://kafka.apache.org/)

Kafka 原理详解

原创 马哥企业教练团队 马哥 Linux 运维 3 天前

Kafka 是什么？

Kafka 是 Apache 旗下的一款分布式流媒体平台，Kafka 是一种高吞吐量、持久性、分布式的发布订阅的消息队列系统。它最初由 LinkedIn(领英)公司发布，使用 Scala 语言编写，与 2010 年 12 月份开源，成为 Apache 的顶级子项目。它主要用于处理消费者规模网站中的所有动作流数据。动作指(网页浏览、搜索和其它用户行动所产生的数据)。

消息系统分类

我们知道常见的消息系统有 Kafka、RabbitMQ、ActiveMQ 等等，但是这些消息系统中所使用的消息模式如下两种：

Peer-to-Peer (Queue)

简称 PTP 队列模式，也可以理解为点到点。例如单发邮件，我发送一封邮件给小徐，我发送过之后邮件会保存在服务器的云端，当小徐打开邮件客户端并且成功连接云端服务器后，可以自动接收邮件或者手动接收邮件到本地，当服务器云端的邮件被小徐消费过之后，云端就不再存储(这根据邮件服务器的配置方式而定)。

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/takfzz/1616130276202-0f04f19a-eec4-44db-a3a5-4880460e16bd.jpeg)

名词解释：

Producer=生产者

Queue=队列

Consumer=消费者

Peer-to-Peer 模式工作原理：

1.消息生产者 Producer1 生产消息到 Queue，然后 Consumer1 从 Queue 中取出并且消费消息。2.消息被消费后，Queue 将不再存储消息，其它所有 Consumer 不可能消费到已经被其它 Consumer 消费过的消息。3.Queue 支持存在多个 Producer，但是对一条消息而言，只会有一个 Consumer 可以消费，其它 Consumer 则不能再次消费。4.但 Consumer 不存在时，消息则由 Queue 一直保存，直到有 Consumer 把它消费。

Publish/Subscribe（Topic）

简称发布/订阅模式。例如我微博有 30 万粉丝，我今天更新了一条微博，那么这 30 万粉丝都可以接收到我的微博更新，大家都可以消费我的消息。

注：以下图示中的 Pushlisher 是错误的名词，正确的为 Publisher

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/takfzz/1616130276240-fca326b2-c114-47c9-92f9-a54b5e47dc99.jpeg)

名词解释：

Publisher=发布者

Topic=主题

Subscriber=订阅者

Publish/Subscribe 模式工作原理：

1.消息发布者 Publisher 将消息发布到主题 Topic 中，同时有多个消息消费者 Subscriber 消费该消息。2.和 PTP 方式不同，发布到 Topic 的消息会被所有订阅者消费。3.当发布者发布消息，不管是否有订阅者，都不会报错信息。4.一定要先有消息发布者，后有消息订阅者。

注意：Kafka 所采用的就是发布/订阅模式，被称为一种高吞吐量、持久性、分布式的发布订阅的消息队列系统。

常用消息系统对比

•RabbitMQ Erlang 编写，支持多协议 AMQP，XMPP，SMTP，STOMP。支持负载均衡、数据持久化。同时 支持 Peer-to-Peer 和发布/订阅模式•Redis 基于 Key-Value 对的 NoSQL 数据库，同时支持 MQ 功能，可做轻量级队列服务使用。就入队操作而言， Redis 对短消息(小于 10KB)的性能比 RabbitMQ 好，长消息的性能比 RabbitMQ 差。•ZeroMQ 轻量级，不需要单独的消息服务器或中间件，应用程序本身扮演该角色，Peer-to-Peer。它实质上是 一个库，需要开发人员自己组合多种技术，使用复杂度高•ActiveMQ JMS 实现，Peer-to-Peer，支持持久化、XA 事务•Kafka/Jafka 高性能跨语言的分布式发布/订阅消息系统，数据持久化，全分布式，同时支持在线和离线处理•MetaQ/RocketMQ 纯 Java 实现，发布/订阅消息系统，支持本地事务和 XA 分布式事务

Kafka 介绍

Kafka 三大特点

1.高吞吐量：可以满足每秒百万级别消息的生产和消费。

2.持久性：有一套完善的消息存储机制，确保数据高效安全且持久化。

3.分布式：基于分布式的扩展；Kafka 的数据都会复制到几台服务器上，当某台故障失效时，生产者和消费者转而使用其它的 Kafka。

Kafka 的几个概念

1.Kafka 作为一个集群运行在一个或多个服务器上，这些服务器可以跨多个机房，所以说 kafka 是分布式的发布订阅消息队列系统。

2.Kafka 集群将记录流存储在称为 Topic 的类别中。

3.每条记录由键值；"key value"和一个时间戳组成。

Kafka 的四个核心 API：

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/takfzz/1616130276252-1256f05d-0cf7-402b-867c-3ba5568c86b3.jpeg)

1.Producer API：生产者 API 允许应用程序将一组记录发布到一个或多个 Kafka Topic 中。2.Consumer API：消费者 API 允许应用程序订阅一个或多个 Topic，并处理向他们传输的记录流。3.Streams API：流 API 允许应用程序充当流处理器，从一个或者多个 Topic 中消费输入流，并将输出流生成为一个或多个输出主题，从而将输入流有效地转换为输出流。4.Connector API：连接器 API 允许构建和运行可重用的生产者或消费者，这些生产者或消费者将 Kafka Topic 连接到现有的应用程序或数据系统。例如：连接到关系数据库的连接器可能会捕获对表的每次更改。

在 Kafka 中，客户端和服务器之间的通信采用 TCP 协议完成，该协议经过版本控制，新版本与旧版本保存向后兼容性，我们为 Kafka 提供了一个 Java 客户端，但是客户端可以使用多种语言。

Kafka 架构简介

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/takfzz/1616130276237-feaedd53-58e7-4d88-9693-a6c869df6110.jpeg)

上图中画了一下 Kafka 的架构，较为凌乱，相同颜色的线条表示一种功能，但是希望大家能够仔细看下

•Producer：消息和数据的生产者，主要负责生产 Push 消息到指定 Broker 的 Topic 中。•Broker：Kafka 节点就是被称为 Broker，Broker 主要负责创建 Topic，存储 Producer 所发布的消息，记录消息处理的过程，现是将消息保存到内存中，然后持久化到磁盘。•Topic：同一个 Topic 的消息可以分布在一个或多个 Broker 上，一个 Topic 包含一个或者多个 Partition 分区，数据被存储在多个 Partition 中。•replication-factor：复制因子；这个名词在上图中从未出现，在我们下一章节创建 Topic 时会指定该选项，意思为创建当前的 Topic 是否需要副本，如果在创建 Topic 时将此值设置为 1 的话，代表整个 Topic 在 Kafka 中只有一份，该复制因子数量建议与 Broker 节点数量一致。•Partition：分区；在这里被称为 Topic 物理上的分组，一个 Topic 在 Broker 中被分为 1 个或者多个 Partition，也可以说为每个 Topic 包含一个或多个 Partition，(一般为 kafka 节. 点数 CPU 的总核心数量)分区在创建 Topic 的时候可以指定。分区才是真正存储数据的单元。•Consumer：消息和数据的消费者，主要负责主动到已订阅的 Topic 中拉取消息并消费，为什么 Consumer 不能像 Producer 一样的由 Broker 去 push 数据呢？因为 Broker 不知道 Consumer 能够消费多少，如果 push 消息数据量过多，会造成消息阻塞，而由 Consumer 去主动 pull 数据的话，Consumer 可以根据自己的处理情况去 pull 消息数据，消费完多少消息再次去取。这样就不会造成 Consumer 本身已经拿到的数据成为阻塞状态。•ZooKeeper：ZooKeeper 负责维护整个 Kafka 集群的状态，存储 Kafka 各个节点的信息及状态，实现 Kafka 集群的高可用，协调 Kafka 的工作内容。

我们可以看到上图，Broker 和 Consumer 有使用到 ZooKeeper，而 Producer 并没有使用到 ZooKeeper，因为 Kafka 从 0.8 版本开始，Producer 并不需要根据 ZooKeeper 来获取集群状态，而是在配置中指定多个 Broker 节点进行发送消息，同时跟指定的 Broker 建立连接，来从该 Broker 中获取集群的状态信息，这是 Producer 可以知道集群中有多少个 Broker 是否在存活状态，每个 Broker 上的 Topic 有多少个 Partition，Prodocuer 会讲这些元信息存储到 Producuer 的内存中。如果 Producoer 像集群中的一台 Broker 节点发送信息超时等故障，Producer 会主动刷新该内存中的元信息，以获取当前 Broker 集群中的最新状态，转而把信息发送给当前可用的 Broker，当然 Prodocuer 也可以在配置中指定周期性的去刷新 Broker 的元信息以更新到内存中。

注意：只有 Broker 和 ZooKeeper 才是服务，而 Producer 和 Consumer 只是 Kafka 的 SDK 罢了

主题和日志

主题和日志官方被称为是 Topic and log。Topic 是记录发布到的类别或者订阅源的名称，Kafka 的 Topic 总是多用户的；也就是说，一个 Topic 可以有零个、一个或者多个消费者订阅写入它的数据。每个 Topic Kafka 集群都为一个 Partition 分区日志，如下图所示：

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/takfzz/1616130276222-33e53691-f9e4-4eff-8871-5d6a56e51622.jpeg)

每个 Partition 分区都是一个有序的记录序列(不可变),如果有新的日志会按顺序结构化添加到末尾，分区中的记录每个都按顺序的分配一个 ID 号，称之为偏移量，在整个 Partition 中具有唯一性。如上图所示，有 Partition、Partition1、Partition2，其中日志写入的顺序从 Old 到 New，ID 号从 0-12 等。

Kafka 集群发布过的消息记录会被持久化到硬盘中，无论该消息是否被消费，发布记录都会被 Kafka 保留到硬盘当中，我们可以设置保留期限。例如，如果保留策略我们设置为两天，则在发布记录的两天内，该消息可供使用，之后则被 Kafka 丢弃以释放空间，Kafka 的性能在数据大小方面是非常出色的，可以长时间保留数据不成问题。

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/takfzz/1616130276252-41c286e2-c058-45e1-8dbb-bd53b942b190.jpeg)

实际上，以消费者为单位地保留的唯一元数据是消费者在日志中的偏移或位置。这个偏移量由消费者控制的：消费者通常会在读取记录时线性地推进偏移量，但事实上，由于消费者的位置时由消费者控制的，所以它可以按照自己喜欢的任何顺序进行消费记录。例如，消费者可以重置之前的偏移量来处理之前的数据，或者直接从最新的偏移量开始消费。这些功能的组合意味着 Kafka 消费者非常的不值一提，他们可以很随便，即使这样，对集群或者其他消费者没有太大影响。例如：可以使用命令工具来“tail”任何 Topic 的内容，而不会更改任何现有使用者所使用的内容。

日志中分区有几个用途。首先，他们允许日志的大小超出适合单台服务器的大小，每个单独的分区必须适合托管它的服务器，但是一个主题可能有许多分区，因此它可以处理任意数量的数据，其次，他们作为并行的单位-更多的是在一点上。

Distribution(分布)

日志 Partition 分区分布在 Kafka 集群中的服务器上，每台服务器都处理数据并请求共享分区。为了实现容错，每个 Partition 分区被复制到多个可配置的 Kafka 集群中的服务器上。

名词介绍：

leader：领导者

followers：追随者

每个 Partition 分区都有一个 leader(领导者)服务器，是每个 Partition 分区，假如我们的 Partition1 分区分别被复制到了三台服务器上，其中第二台为这个 Partition 分区的领导者，其它两台服务器都会成为这个 Partition 的 followers(追随者)。其中 Partition 分片的 leader(领导者)处理该 Partition 分区的所有读和写请求，而 follower(追随者)被动地复制 leader(领导者)所发生的改变，如果该 Partition 分片的领导者发生了故障等，两个 follower(追随者)中的其中一台服务器将自动成为新的 leader 领导者。每台服务器都充当一些分区的 leader(领导者)和一些分区的 follower(追随者)，因此集群内的负载非常平衡。

注意：上面讲的 leader 和 follower 仅仅是每个 Partition 分区的领导者和追随者，并不是我们之前学习到的整个集群的主节点和备节点，希望大家不要混淆。

Geo-Replication(地域复制)

Kafka Mirrormaker 为集群提供地域复制支持，使用 MirrorMaker，可以跨多个机房或云端来复制数据，可以在主动/被动方案中使用它进行备份和恢复；在主动方案中，可以使数据更接近用户，或支持数据位置要求。

Producers(生产者)

生产者将数据发布到他们选择的 Topic，生产者负责选择分配给 Topic 中的哪个分区的记录。这可以通过循环方式来完成，只是为了负载均衡，或者可以根据一些语义分区函数(比如基于记录中的某个键)来完成。

Consumers(消费者)

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/takfzz/1616130276249-40cdfc85-7b14-4783-9add-ea99966847a1.jpeg)

名词介绍

Consumers：消费者

Consumers Group：消费者组

Consumers Group name：消费者组名

Consumers 使用 Consumers Group name 标记自己，并且发布到 Topic 的每个记录被传递到每个订阅 Consumers Group 中的一个 Consumers 实例，Consumers 实例可以在单独的进程中，也可以在不同的机器，如果所有 Consumers 实例具有相同的 Consumers Group，则记录将有效地在 Consumers 上进行负载均衡。

如果所有 Consumers 实例在不同的 Consumers Group 中，则每个记录将广播到所有 Consumers 进程中。

两个 Kafka Cluster，托管了四个 Partition(分区)，从 P0-P3，包含两个 Consumers Group 分别是 Consumer Group A 和 Consumer Group B，Consumners Group A 有两个 Consumers 实例，B 有四个 Consumers 实例。也就是消费者 A 组有两个消费者，B 组有四个消费者。然后，更常见的是，我们发现 Topic 有少量的 Consumers Group，每个消费者对应一个用户组，每个组有许多消费者实例组成，用于可伸缩和容错，这只不过是发布/订阅语义，其中订阅者是一组消费者，而不是单个进程。

在 Kfaka 中实现消费者的方式是通过在消费者实例上划分日志中的 Partition 分区，以便每个实例在任何时间点都是分配的“相同份额”，维护消费者组成功资格的过程由 Kafka 动态协议实现，如果新的消费者实例加入该消费者组，新消费者实例将从该组的其它成员手里接管一些分区；如果消费者实例故障，其分区将分发给其余消费者实例。

Kafka 仅提供分区内记录的总顺序，而不是 Topic 中不同分区之间的记录。对于大多数应用程序而言，按分区排序和按键分许数据的能力已经足够，但是如果你需要记录总顺序，则可以使用只有一个分区的 Topic 来实现，尽管这意味着每个消费者组只有一个消费者进程。

Consumer Group

我们开始处有讲到消息系统分类：P-T-P 模式和发布/订阅模式，也有说到我们的 Kafka 采用的就是发布订阅模式，即一个消息产生者产生消息到 Topic 中，所有的消费者都可以消费到该条消息，采用异步模型；而 P-T-P 则是一个消息生产者生产的消息发不到 Queue 中，只能被一个消息消费者所消费，采用同步模型 其实发布订阅模式也可以实现 P-T-P 的模式，即将多个消费者加入一个消费者组中，例如有三个消费者组，每个组中有 3 个消息消费者实例，也就是共有 9 个消费者实例，如果当 Topic 中有消息要接收的时候，三个消费者组都会去接收消息，但是每个人都接收一条消息，然后消费者组再将这条消息发送给本组内的一个消费者实例，而不是所有消费者实例，到最后也就是只有三个消费者实例得到了这条消息，当然我们也可以将这些消费者只加入一个消费者组，这样就只有一个消费者能够获得到消息了。

Guarantees(担保)

在高级别的 Kafka 中提供了以下保证：

•生产者发送到特定 Topic 分区的消息将按照其发送顺序附加。也就是说，如果一个 Producers 生产者发送了 M1 和 M2，一般根据顺序来讲，肯定是先发送的 M1，随后发送的 M2，如果是这样，假如 M1 的编号为 1，M2 的编号为 2，那么在日志中就会现有 M1，随后有 M2。•消费者实例按照他们存储在日志中的顺序查看记录。•对于具有复制因子 N 的 Topic，Kafka 最多容忍 N-1 个服务器故障，则不会丢失任何提交到日志的记录。
