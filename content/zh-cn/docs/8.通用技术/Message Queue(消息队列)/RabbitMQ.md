---
title: RabbitMQ
---

# 概述

> 参考：
> - 官方网址：<https://www.rabbitmq.com/>
> - 消息队列模拟器：<http://tryrabbitmq.com/>
> - <https://www.yuque.com/noobwo/mq/hpiop0>
> - <https://m.6-km.com/next/quorum-queues.html>
> - <https://zhuanlan.zhihu.com/p/63700605>

Rabbit Message Queue(Rabbit 消息队列，简称：RabbitMQ)。是一个在 AMQP 基础上实现的，可复用的消息队列服务。

Advanced Message Queuing Protocol(高级消息队列协议，简称 AMQP) ，是一个提供统消息服务的应用层(7 层)协议。其设计目标是对于消息的排序、路由（包括点对点和订阅-发布）、保持可靠性、保证安全性\[1]。AMQP 规范了消息传递方和接收方的行为，以使消息在不同的提供商之间实现互操作性，就像 SMTP，HTTP，FTP 等协议可以创建交互系统一样。与先前的中间件标准（如 Java 消息服务）不同的是，JMS 在特定的 API 接口层面和实现行为上进行了统一，而高级消息队列协议则关注于各种消息如何以字节流的形式进行传递。因此，使用了符合协议实现的任意应用程序之间可以保持对消息的创建、传递。

## 工作机制概述

在了解消息通讯之前首先要了解 3 个概念：生产者、消费者和代理。

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/vi79ye/1616130695655-d434cc0d-b15a-48c2-9fad-7fa971d2ffa4.jpeg)

Publisher(生产者)：消息的创建者，负责创建和推送数据到消息服务器；

Consumer(消费者)：消息的接收方，用于处理数据和确认消息；

Broker(代理)：就是 RabbitMQ 本身，用于扮演“快递”的角色，本身不生产消息，只是扮演“快递”的角色。

消息发送原理

首先你必须连接到 Rabbit 才能发布和消费消息，那怎么连接和发送消息的呢？

你的应用程序和 Rabbit Server 之间会创建一个 TCP 连接，一旦 TCP 打开，并通过了认证，认证就是你试图连接 Rabbit 之前发送的 Rabbit 服务器连接信息和用户名和密码，有点像程序连接数据库，使用 Java 有两种连接认证的方式，后面代码会详细介绍，一旦认证通过你的应用程序和 Rabbit 就创建了一条 AMQP 信道（Channel）。

信道是创建在“真实”TCP 上的虚拟连接，AMQP 命令都是通过信道发送出去的，每个信道都会有一个唯一的 ID，不论是发布消息，订阅队列或者介绍消息都是通过信道完成的。

为什么不通过 TCP 直接发送命令？

对于操作系统来说创建和销毁 TCP 会话是非常昂贵的开销，假设高峰期每秒有成千上万条连接，每个连接都要创建一条 TCP 会话，这就造成了 TCP 连接的巨大浪费，而且操作系统每秒能创建的 TCP 也是有限的，因此很快就会遇到系统瓶颈。

如果我们每个请求都使用一条 TCP 连接，既满足了性能的需要，又能确保每个连接的私密性，这就是引入信道概念的原因。

其实在生活中，这种模型用得非常多，就比如我们都会接触的网购快递，可以说是一个典型的消息队列的 case 了：

商家不断的把商品扔给快递公司（注意不是直接将商品给买家），而快递公司则将商品根据地质分发对应的买家

对上面这个过程进行拆解，可以映射扮演的角色

- 商品：Message，传递的消息，由商家投递给快递公司时，需要进行打包（一般 Producer 生产消息也会将实体数据进行封装）
- 商家：Produer 生产者
- 快递公司： Queue，消息的载体
- 买家：Consumer 消费者

那么快递公司时怎么知道要把商品给对应的买家呢？根据包裹上的地址+电话

- 同样消息队列也需要一个映射规则，实现 Message 和 Consumer 之间的路由

# RabbitMQ 基本概念

1. Exchange
2. Queue
3. Routing key # 一种规则，用于 Exchange 路由消息到匹配到规则的 Queue 上
4. Binding # 交换器 与 队列(或另一个交换器) 的关联行为

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/vi79ye/1616130695576-789c9f61-914e-4f88-951c-b1a0915fdc4a.jpeg)

## Broker # 指安装了 RabbitMQ 的服务器

## Virtual Host(虚拟主机) # 类似 RabbitMQ 虚拟化的感觉

每个 Rabbit 都能创建很多 vhost 我们称之为虚拟主机，每个虚拟主机其实都是 mini 版的 RabbitMQ，拥有自己的队列，交换器和绑定，拥有自己的权限机制。

1. RabbitMQ 默认的 vhost 是 / 开箱即用；
2. 多个 vhost 是隔离的，多个 vhost 无法通讯，并且不用担心命名冲突（队列和交换器和绑定），实现了多层分离；
3. 创建用户的时候必须指定 vhost；

注意，从上层角度看，RabbitMQ 实现主要依靠 Exchange 与 Queue 来实现。

## Queue(队列) # 用于存储消息

官方文档：<https://www.rabbitmq.com/queues.html>

## Exchange(交换器) # 用于接受、分配消息

官方文档：未知

客户端向 RebbitMQ 发送消息，并不会直接发送到 Queue 中，而是从 5674 端口到达 Exchange，然后根据路由规则，决定收到消息应该投递到哪个队列上，这些规则其中一部分称为 **Routing key(路由键)**，另外一部分规则是根据 Exchange 的类型来决定。

Queue 可以通过 Routing key 关联到 Exchange 上(也可以省略 Routing key 直接绑定)，这种关联行为称为 **Binding(绑定)**。

注意：

1. 绑定时所使用的 Routing key 可以使用通配符。
2. 为了与每个消息中的 Routing key 区分，我们一般称绑定时使用的 Routing key 为 **Binding key(绑定键)**。

绑定完成后，当 Exchange 收到消息(**每个消息也可以指定 Routing key**)后，凡是 Routing key 匹配 Binding key 条件的，则该消息会被发送到具有 Binding key 的队列中。如果匹配多分，则消息在 Exchange 中会复制多份。

例如：现在发送一条消息，Routing key 为 test.1，两个队列绑定到交换器时，Binding key 分别为 _.1 和 test._，则该消息到达 Exchange 时，会复制为两份，分别发送到这两各队列中。

**Exchange 的 4 种类型**

1. direct # Routing key 与 Binding key 必须完全匹配，不接受通配符匹配。
2. fanot # 不需要 routing key，收到消息后，自动将消息复制多份并发送给与自己绑定的队列上。
3. headers
4. topic # Routing key 与 Binding key 可以根据通配符匹配。匹配规则如下：
   1. - 匹配一个单词
   2. # 匹配 0 个或多个字符
   3. \*，# 只能写在.号左右，且不能挨着字符
   4. 单词和单词之间需要用.隔开。
5. 特殊的 Exchange
   1. 如果不手动绑定 Exchange 与 Queue，发送消息时不指定 Exchange，但是消息依然可以路由到队列中。
   2. 这是因为每个 vhost 下都有一个 default Exchange。默认的 Exchange 不可以手动绑定，没有名字，不可以删除，direct 类型。每个 Queue 都会与当前 vhost 下的 默认 Exchange 隐式绑定在一起，且 Binding key 为 Queue 名称。
   3. 消息中的 Routing key 必须指定为 Queue 名称，这时消息通过默认的 Exchange 被路由到对应的 Queue 上。这保证了消息准确投递。

## Connection(连接) # 应用程序与 Rabbit 之间建立的连接，程序代码中使用

## Channel(信道) # 消息推送使用的通道

# RabbitMQ 配置
