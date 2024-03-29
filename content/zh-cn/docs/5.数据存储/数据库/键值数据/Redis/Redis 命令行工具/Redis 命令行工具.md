---
title: redis 命令行工具
weight: 1
---
# 概述

> 参考：
>
> - [官方文档](https://redis.io/commands)
>   - 中文

# redis-cli # 命令行客户端

help # 按 TAB 键遍历所有组以及对应的命令

help @ # 获取在指定 group 中的命令列表

help # 获取指定 command 的帮助信息

# Group 列表

# generic # 通用的

详见：[generic 组](/docs/5.数据存储/数据库/键值数据/Redis/Redis%20命令行工具/generic%20组.md)

主要是对 key 的操作，比如列出所有 key，删除 key 等等。一般在列出 key 时无法获取对应的 value

# string # 字符串

详见：[string 组](/docs/5.数据存储/数据库/键值数据/Redis/Redis%20命令行工具/string%20组.md)

主要是对于 value 的操作，比如对指定的一个或多个 key 的 value 进行增、删、改、查。

# list

# set

# sorted_set

# hash

# pubsub # 发布与订阅

### PUBSUB subcommand \[argument \[argument ...]] # 检查发布/订阅子系统的状态

**PUBSUB CHANNELS \[pattern]**# 列出当前活动的频道。

活动频道是具有一个或多个订阅者（不包括订阅模式的客户端）的发布/订阅频道。如果未`pattern`指定，则列出所有通道，否则，如果指定了模式，则仅列出与指定的 glob 样式模式匹配的通道。

**PUBSUB NUMSUB \[channel-1 ... channel-N]**# 获取指定频道的订阅数

PUBSUB NUMPAT # 与 NUMSUB 命令类似，只不过是获取通过 PSUBSCRIBE 命令订阅的频道的订阅数

### SUBSCRIBE channel \[channel ...] # 订阅指定频道，以获取发布到这些频道上的消息

> 注意：多个频道以空格分隔

PSUBSCRIBE pattern \[pattern ...] # 与 SUBSCRIBE 命令类似，只不过是通过表达式来匹配多个频道

### UNSUBSCRIBE \[channel \[channel ...]] # 退订指定频道，停止接收发布到这些频道上的消息

PUNSUBSCRIBE \[pattern \[pattern ...]] # 与 UNSUBSCRIBE 命令类似，只不过是通过表达式来匹配多个频道

### PUBLISH channel message # 向一个频道发送一条消息

# transactions

# connection # redis 的客户端与服务端连接相关的命令

包括用密码认证以操作数据库、选择连接哪个数据库等

### AUTH password # 输入密码进行认证后，即可操作数据库

修改配置文件`requirepass PASSWORD`设置密码

ECHO message # Echo the given string

PING \[message] # Ping the server

QUIT - # Close the connection

### SELECT index # 为当前连接更改所选数据库，可以把每个数据库当做 namespace 来看。默认为 0 号数据库

# server # 关于 redis 服务器的相关命令，包括查看配置等

详见 [server 组](/docs/5.数据存储/数据库/键值数据/Redis/Redis%20命令行工具/server%20组.md)

# scripting

# hyperloglog

# cluster # 与集群模式相关的命令

# geo

# stream
