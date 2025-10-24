---
title: Redis CLI
linkTitle: Redis CLI
weight: 1
---

# 概述

> 参考：
>
> - [官方文档](https://redis.io/commands)
>   - 中文

# redis-cli - 命令行客户端

help # 按 TAB 键遍历所有组以及对应的命令

help @ # 获取在指定 group 中的命令列表

help # 获取指定 command 的帮助信息

# Group 列表

- [generic](#generic) # 通用的
- 不同数据类型的数据处理命令
  - [string](#string)
  - etc.
- [pubsub](#pubsub) # 发布与订阅
- [connection](#connection) # redis 的客户端与服务端连接相关的命令
- [server](#server) # 关于 redis 服务器的相关命令，包括查看配置等
- etc.

# generic

详见：[generic 组](/docs/5.数据存储/数据库/键值数据/Redis/Redis%20CLI/generic%20组.md)

主要是对 key 的操作，比如列出所有 key，删除 key 等等。一般在列出 key 时无法获取对应的 value

# 不同类型的数据处理命令组

可以通过 [generic 组](docs/5.数据存储/数据库/键值数据/Redis/Redis%20CLI/generic%20组.md) 的 TYPE 指令确定 Key 的类型

## string

详见：[string 组](/docs/5.数据存储/数据库/键值数据/Redis/Redis%20CLI/string%20组.md)

主要是对于 value 的操作，比如对指定的一个或多个 key 的 value 进行增、删、改、查。

## list

## set

## sorted_set

## hash

hash 组下的所有命令用于管理 hash 数据类型的数据

**HGETALL key** # 获取 key 的所有字段名和值

**HKEYS key** # 获取 key 的所有字段的名

**HVALS keys** # 获取 key 的所有字段的值

**HGET key field** # 获取 hash 类型 key 的 field 的值。

## stream

# pubsub

https://redis.io/docs/latest/commands/?group=pubsub

### PUBSUB - 检查发布/订阅子系统的状态

**PUBSUB subcommand \[argument \[argument ...]]**

**PUBSUB CHANNELS \[pattern]** # 列出当前活动的频道。

活动频道是具有一个或多个订阅者（不包括订阅模式的客户端）的发布/订阅频道。如果未`pattern`指定，则列出所有通道，否则，如果指定了模式，则仅列出与指定的 glob 样式模式匹配的通道。

**PUBSUB NUMSUB \[channel-1 ... channel-N]** # 获取指定频道的订阅数

**PUBSUB NUMPAT** # 与 NUMSUB 命令类似，只不过是获取通过 PSUBSCRIBE 命令订阅的频道的订阅数

### SUBSCRIBE - 订阅指定频道，以获取发布到这些频道上的消息

**SUBSCRIBE channel \[channel ...]**

> 注意：多个频道以空格分隔

PSUBSCRIBE pattern \[pattern ...] # 与 SUBSCRIBE 命令类似，只不过是通过表达式来匹配多个频道

### UNSUBSCRIBE - 退订指定频道，停止接收发布到这些频道上的消息

**UNSUBSCRIBE \[channel \[channel ...]]**

PUNSUBSCRIBE \[pattern \[pattern ...]] # 与 UNSUBSCRIBE 命令类似，只不过是通过表达式来匹配多个频道

### PUBLISH channel message - 向一个频道发送一条消息

# transactions

# connection

包括用密码认证以操作数据库、选择连接哪个数据库等

### AUTH password - 输入密码进行认证后，即可操作数据库

修改配置文件`requirepass PASSWORD`设置密码

ECHO message # Echo the given string

PING \[message] # Ping the server

QUIT - # Close the connection

### SELECT index - 为当前连接更改所选数据库，可以把每个数据库当做 namespace 来看。默认为 0 号数据库

# server

详见 [server 组](/docs/5.数据存储/数据库/键值数据/Redis/Redis%20CLI/server%20组.md)

# scripting

# hyperloglog

# cluster - 与集群模式相关的命令

详见 [Cluster 命令](/docs/5.数据存储/数据库/键值数据/Redis/Redis%20CLI/Cluster%20命令.md)

# geo

# 最佳实践

`select 0` 进入 0 号数据库，默认从 0 - 15。

`info` 命令返回的结果中，Keyspace 部分可以看到当前在哪几个 Database 中存了多少个 Keys。

```bash
# Keyspace
db0:keys=9,expires=0,avg_ttl=0
db15:keys=14,expires=5,avg_ttl=682335
```

`keys *` 列出当前数据库中所有的键
