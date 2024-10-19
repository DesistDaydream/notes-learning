---
title: Pub/Sub(发布/订阅)
---

# 概述

> 参考：
>
> - [官方文档](https://redis.io/topics/pubsub)
>   - https://redis.io/docs/latest/develop/interact/pubsub/
> - [Wiki, 发布/订阅 模式](https://en.wikipedia.org/wiki/Publish%E2%80%93subscribe_pattern)

Redis 可以通过 SUBSCRIBE、UNSUBSCRIBE、PUBLISH 这类命令实现 Publish/Subscribe(发布/订阅)模式。

Redis 的 发布/订阅模式中，**Messages(消息)** 的**发送者称为 Publishers(发布者)**、消息的**接收者称为 Subscribers(订阅者)**。而发送者和接收者之间传递消息的**途径称为 Channels(频道)**。

- 订阅者可以订阅自己感兴趣的 Channels，并随时等待接收发布到这些 Channels 中的消息，并不需要知道有哪些发布者。

- 发布者可以向任何 Channels 中发布消息，而不需要知道有哪些订阅者

这种将 发布者 与 订阅者 解耦的模式，可以实现更大的可扩展性和更动态的网络拓扑结构。

Redis 中并没有默认已经存在的 Channels(频道)。 只要执行了 `SUBSCRIBE` 命令，并指定 Channel 名称，Redis 就会创建一个频道，并且执行该命令的客户端就称为 Subscriber。比如，现在执行如下命令：

```basic
127.0.0.1:6379> SUBSCRIBE test1 test2
Reading messages... (press Ctrl-C to quit)
1) "subscribe"
2) "test1"
3) (integer) 1
1) "subscribe"
2) "test2"
3) (integer) 2
```

此时，创建了两个 Channels，test1 和 test2，其他客户端发送到这些 Channels 的消息，将被推送到订阅了这俩通道的 Subscriber，也就是订阅者的客户端。效果如下：

```
# 订阅 Channels
127.0.0.1:6379> SUBSCRIBE test1 test2
Reading messages... (press Ctrl-C to quit)
1) "subscribe"
2) "test1"
3) (integer) 1
1) "subscribe"
2) "test2"
3) (integer) 2

# 发布消息
127.0.0.1:6379> PUBLISH test1 message1
(integer) 1

# 接收到了消息
127.0.0.1:6379> SUBSCRIBE test1 test2
Reading messages... (press Ctrl-C to quit)
1) "subscribe"
2) "test1"
3) (integer) 1
1) "subscribe"
2) "test2"
3) (integer) 2
1) "message"
 # 这里就是接收到的消息
2) "test1" # 接收到消息的频道
3) "message1" # 该频道接收到的消息内容
```

注意：

- 订阅了一个或多个频道的客户端尽管可以订阅和取消订阅其他频道，但不应发出命令。对订阅和取消订阅操作的回复以消息的形式发送，以便客户端可以读取连贯的消息流，其中第一个元素表示消息的类型。在已订阅客户端的上下文中允许使用的命令是 SUBSCRIBE，PSUBSCRIBE，UNSUBSCRIBE，PUNSUBSCRIBE， PING 和 QUIT。
- 但是，`redis-cli`一旦进入订阅模式，该命令将不接受任何命令，只能通过退出该模式`Ctrl-C`。
