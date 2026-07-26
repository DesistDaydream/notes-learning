---
title: Pub/Sub(发布/订阅)
weight: 4
---

# 概述

> 参考：
>
> - [官方文档，Redis 发布/订阅](https://redis.io/topics/pubsub)
>     - https://redis.io/docs/latest/develop/interact/pubsub/
>     - https://valkey.io/topics/pubsub/
> - [Wiki, 发布/订阅 模式](https://en.wikipedia.org/wiki/Publish%E2%80%93subscribe_pattern)

Redis 可以通过 SUBSCRIBE、UNSUBSCRIBE、PUBLISH 命令及其衍生命令，实现 **Publish/Subscribe(发布/订阅)** 模式。

Redis 的 发布/订阅 模式中，**Messages(消息)** 的**发送者**称为 **Publishers(发布者)**、消息的**接收者**称为 **Subscribers(订阅者)**。而发送者和接收者之间传递消息的**途径**称为 **Channels(频道)**。

- 订阅者可以订阅自己感兴趣的 Channels，并随时等待接收发布到这些 Channels 中的消息，并不需要知道有哪些发布者。
- 发布者可以向任何 Channels 中发布消息，而不需要知道有哪些订阅者

这种将 发布者 与 订阅者 解耦的模式，可以实现更大的可扩展性和更动态的网络拓扑结构。

只要执行了 `SUBSCRIBE` 命令，并指定 Channel 名称，Redis 就会创建一个 Channel，并且执行该命令的客户端就称为 Subscriber。比如，现在执行如下命令：

```bash
# 订阅 Channels
127.0.0.1:6379> SUBSCRIBE test1 test2
Reading messages... (press Ctrl-C to quit)
1) "subscribe"
2) "test1"
3) (integer) 1
1) "subscribe"
2) "test2"
3) (integer) 2
```

此时，创建了两个 Channels，test1 和 test2。

其他客户端发送到这些 Channels 的消息，将被推送到订阅了这俩通道的 Subscriber，也就是订阅者的客户端。效果如下：

```bash
# 发布消息
127.0.0.1:6379> PUBLISH test1 message1
(integer) 1
```

回到刚才订阅 Channels 的客户端中，可以看到：

```bash
# 订阅通道
127.0.0.1:6379> SUBSCRIBE test1 test2
Reading messages... (press Ctrl-C to quit)
1) "subscribe"
2) "test1"
3) (integer) 1
1) "subscribe"
2) "test2"
3) (integer) 2
1) "message"
# 下面是接收到的消息
1) "test1" # 接收到消息的频道
2) "message1" # 该频道接收到的消息内容
```

注意：

- 订阅了一个或多个频道的客户端尽管可以订阅和取消订阅其他频道，但不应发出命令。对订阅和取消订阅操作的回复以消息的形式发送，以便客户端可以读取连贯的消息流，其中第一个元素表示消息的类型。在已订阅客户端的上下文中允许使用的命令是 SUBSCRIBE，PSUBSCRIBE，UNSUBSCRIBE，PUNSUBSCRIBE， PING 和 QUIT。
- 但是，`redis-cli`一旦进入订阅模式，该命令将不接受任何命令，只能通过 `Ctrl-C` 退出该模式。

# Keyspace notifications

https://valkey.io/topics/notifications/

**[Keyspace](/docs/5.数据存储/数据库/键值数据/键值数据.md#Keyspace) notifications(键空间通知)** 是 Pub/Sub 的内置 Channels。这些 Channels 可以让我们实时监控对键和值的更改。

默认情况下，Keyspace 通知相关 Channels 是禁用的（因为该功能会占用一些 CPU 资源）。我们可以通过 valkey.conf 文件中的 notify-keyspace-envents 配置，或 `CONFIG SET` 命令启用键空间通知功能。

以 config 命令为例：

```bash
config set notify-keyspace-events ${Parameters}
```

可用的 Parameter 有：

- **K** # Keyspace 事件，消息会发布到名为 `__keyspace@<DB>__:<KEY>` 的 Channel 中。其中 DB 是数据库号；KEY 是键名
    - 关注点是：某些键发生了什么事件
- **E** # Keyevent 事件，消息会发布到名为 `__keyevent@<DB>__:<EVENT>` 的 Channel 中。其中 DB 是数据库号；EVENT 是事件名
    - 关注点是：某些事件涉及了哪些键
- **g** # Generic commands (non-type specific) like DEL, EXPIRE, RENAME, ...
- **$** # String commands
- **l** # List commands
- **s** # Set commands
- **h** # Hash commands
- **z** # Sorted set commands
- **t** # Stream commands
- **d** # Module key type events
- **x** # Expired events (events generated every time a key expires)
- **e** # Evicted events (events generated when a key is evicted for maxmemory)
- **m** # Key miss events (events generated when a key that doesn't exist is accessed)
- **n** # New key events (Note: not included in the 'A' class)
- **A** # 等价于 `g$lshztxed` 这堆参数。若是 `AKE` 参数，则表示除了 m 和 n 参数之外的所有参数。

> [!Attention]
> <font color="#ff0000">参数中至少应包含 `K` 或 `E`</font>，否则无论参数的其余部分如何，都不会发布消息。
>
> get 相关命令不会产生消息。

这些 Channels 会接收的事件包括：

- 所有影响给定 Keys 的命令
- etc.

## Example

配置 `config set notify-keyspace-events KE$` 后，我们订阅 keyspace 与 keyevent 的相关通道，执行 `set demokey desistdaydream` 命令可以看到这些通道发布过来的消息：

```bash
~]# docker run -it --rm --network host --name valkey-cli valkey/valkey:9.1-alpine valkey-cli
127.0.0.1:6379> config set notify-keyspace-events KE$
OK
127.0.0.1:6379> set demokey desistdaydream
OK
```

订阅 `__keyspace@0__:demokey*` 通道后（订阅数据库 0 中，所有 demokey 开头的键的消息），收到的信息示例：

```bash
~]# docker run -it --rm --network host --name valkey-cli-key valkey/valkey:9.1-alpine valkey-cli
127.0.0.1:6379> psubscribe '__keyspace@0__:demokey*'
1) "psubscribe"
2) "__keyspace@0__:demokey*"
3) (integer) 1
# 下面是通道收到的消息
1) "pmessage"
2) "__keyspace@0__:demokey*"
3) "__keyspace@0__:demokey"
4) "set"
```

订阅 `__keyevent@0__:set*` 通道后（订阅数据库 0 中，所有 set 开头的命令的消息），收到的信息示例：

```bash
~]# docker run -it --rm --network host --name valkey-cli-event valkey/valkey:9.1-alpine valkey-cli
127.0.0.1:6379> psubscribe '__keyevent@0__:set*'
1) "psubscribe"
2) "__keyevent@0__:set*"
3) (integer) 1
# 下面是通道收到的消息
1) "pmessage"
2) "__keyevent@0__:set*"
3) "__keyevent@0__:set"
4) "demokey"
```

可以看到：

- keyspace 显示逻辑是：消息来源是 demokey 这个键，执行了 set 命令；
- keyevent 显示逻辑是：消息来源是 set 这个命令，对 demokey 键执行的。

# 最佳实践
