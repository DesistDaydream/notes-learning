---
title: server 组
linkTitle: server 组
date: 2024-06-14T10:20
weight: 20
---


# 概述

> 参考：
>
> - https://redis.io/docs/latest/commands/?group=server

BGREWRITEAOF - # Asynchronously rewrite the append-only file

BGSAVE -

summary: Asynchronously save the dataset to disk

since: 1.0.0

CLIENT GETNAME -

summary: Get the current connection name

since: 2.6.9

CLIENT ID -

summary: Returns the client ID for the current connection

since: 5.0.0

CLIENT KILL \[ip:port] \[ID client-id] \[TYPE normal|master|slave|pubsub] \[ADDR ip:port] \[SKIPME yes/no]

summary: Kill the connection of a client

since: 2.4.0

CLIENT LIST -

summary: Get the list of client connections

since: 2.4.0

CLIENT PAUSE timeout

summary: Stop processing commands from clients for some time

since: 2.9.50

CLIENT REPLY ON|OFF|SKIP

summary: Instruct the server whether to reply to commands

since: 3.2

CLIENT SETNAME connection-name

summary: Set the current connection name

since: 2.6.9

CLIENT UNBLOCK client-id \[TIMEOUT|ERROR]

summary: Unblock a client blocked in a blocking command from a different connection

since: 5.0.0

COMMAND -

summary: Get array of Redis command details

since: 2.8.13

COMMAND COUNT -

summary: Get total number of Redis commands

since: 2.8.13

COMMAND GETKEYS -

summary: Extract keys given a full Redis command

since: 2.8.13

COMMAND INFO command-name \[command-name ...]

summary: Get array of specific Redis command details

since: 2.8.13

# config
## config get PARAMETER - 获取 Redis 启动后，内存中的配置参数

EXAMPLE

- **`config get *`** # 获取所有配置参数的值

- **confige get maxmemory** # 获取 maxmemory 配置参数的值

## config resetstat - 重置 INFO 命令的状态信息

## config rewrite - 用内存中的配置重新配置文件

该命令会重写 Redis 启动时指定的 redis.conf 文件。有时候使用 config set 命令更改了内存中的配置，为了防止重启 Redis 时，配置丢失，所以需要执行 config rewrite 命令，以便让 Redis 下次启动时，可以加载正确的配置文件。

注意：执行该命令的客户端需要对 redis.conf 文件具有写权限才可以正常执行命令，否则报错 `(error) ERR Rewriting config file: Invalid argument`

## config set PARAMETER VALUE - 设置 redis 中指定配置参数的值

CONFIG SET 命令用于在服务器运行期间重写某些配置，而不用重启 Redis。你可以使用此命令更改不重要的参数或从一个参数切换到另一个持久性选项。

可以通过 CONFIG GET \*获得 CONFIG SET 命令支持的配置参数列表，该命令是用于获取有关正在运行的 Redis 实例的配置信息的对称命令。

所有使用 CONFIG SET 设置的配置参数将会立即被 Redis 加载，并从下一个执行的命令开始生效。

所有支持的参数与 redis.conf 文件中使用的等效配置参数具有相同含义，但有以下重要区别：

- 在指定字节或其他数量的选项中，不能使用在 redis.conf 中使用的简写形式（如 10k，2gb 等），所有内容都应该指定为格式良好的 64 位整数，以配置指令的基本单位表示。但从 Redis3.0 以及更高版本开始，可以将 CONFIG SET 与内存单元一起用于 maxmemory、客户端输出缓冲以及复制积压大小（repl-backlog-size）指定内存单位。

- save 参数是一个以空格分隔的整数字符串。每对整数代表一个秒/修改阈值。

例如在 redis.conf 中看起来像这样：

```bash
save 900 1save 300 10
```

这意味着，如果数据集有 1 个以上变更，则在 900 秒后保存；如果有 10 个以上变更，则在 300 秒后就保存，应使用\`CONFIG SET SAVE “900 1 300 10”来设置。

可以使用 CONFIG SET 命令将持久化从 RDB 快照切换到 AOF 文件（或其他相似的方式）。 有关如何执行此操作的详细信息，请查看 persistencepage。

一般来说，你应该知道将 appendonly 参数设置为 yes 将启动后台进程以保存初始 AOF 文件（从内存数据集中获取），并将所有后续命令追加到 AOF 文件，从而达到了与一个 Redis 服务器从一开始就开启了 AOF 选项相同的效果。

如果你愿意，可以同时开启 AOF 和 RDB 快照，这两个选项不是互斥的。

DBSIZE -

summary: Return the number of keys in the selected database

since: 1.0.0

DEBUG OBJECT key

summary: Get debugging information about a key

since: 1.0.0

DEBUG SEGFAULT -

summary: Make the server crash

since: 1.0.0

FLUSHALL \[ASYNC]

summary: Remove all keys from all databases

since: 1.0.0

FLUSHDB \[ASYNC]

summary: Remove all keys from the current database

since: 1.0.0

# info \[SECTION] - 获取 Redis 服务器的信息和统计信息

http://www.redis.cn/commands/info.html

**info \[SECTION]**

info 命令以一种易于理解和阅读的格式，返回关于 Redis 服务器的各种信息和统计数值。

通过给定可选的参数 SECTION，可以让命令只返回某一部分的信息。如果没有指定任何 SECTION，默认为 default。

- server # Redis 服务器的一般信息
- clients # 客户端的连接部分
- memory #  内存消耗相关信息
- persistence # RDB 和 AOF 相关信息
- stats # 一般统计
- replication # 主/从复制信息
- cpu # 统计 CPU 的消耗
- commandstats # Redis 命令统计
- cluster # Redis 集群信息
- keyspace # 数据库的相关统计

它也可以采取以下值:

- all # 返回所有信息
- default # 值返回默认设置的信息

LASTSAVE -

summary: Get the UNIX time stamp of the last successful save to disk

since: 1.0.0

MEMORY DOCTOR -

summary: Outputs memory problems report

since: 4.0.0

MEMORY HELP -

summary: Show helpful text about the different subcommands

since: 4.0.0

MEMORY MALLOC-STATS -

summary: Show allocator internal stats

since: 4.0.0

MEMORY PURGE -

summary: Ask the allocator to release memory

since: 4.0.0

MEMORY STATS -

summary: Show memory usage details

since: 4.0.0

MEMORY USAGE key \[SAMPLES count]

summary: Estimate the memory usage of a key

since: 4.0.0

## monitor - 实时监听 Redis 服务器收到的所有请求

这些请求中也包括 Channel(频道) 中发布的消息。

效果如下：

```bash
127.0.0.1:6379> monitor
OK
1613227324.358400 [0 172.19.42.232:24454] "PING"
.......
1613227325.569008 [0 172.19.42.233:60388] "PUBLISH" "__sentinel__:hello" "172.19.42.233,26379,dc3b5dd2ab650290b08971ac719e6df182300109,2,mymaster,172.19.42.231,6379,2"
1613227325.783351 [0 172.19.42.232:24454] "PUBLISH" "__sentinel__:hello" "172.19.42.232,26379,f387733988ef9590e7616f04701b9da4a387915b,2,mymaster,172.19.42.231,6379,2"
1613227325.922619 [0 172.19.42.231:33198] "PUBLISH" "__sentinel__:hello" "172.19.42.231,26379,eff7ec8b5312fe23cece432d8b89e3926e8ea92b,2,mymaster,172.19.42.231,6379,2"
```

## replicaof HOST POST - 使该节点变为 replica 或 master 角色

该命令可以动态更改 Replication 模式的配置。

EXAMPLE

- 将节点变为 master

  - **replicaof no one**

- 将该节点作为 172.19.42.231:6379 节点的 Replica 节点

  - **replicaof 172.19.42.231 6379**

## role - 显示当前节点的 Replication 模式的信息

效果如下：

```bash
127.0.0.1:6379> role
1) "master"
 # 当前节点的角色
2) (integer) 60490530
3) 1) 1) "172.19.42.232"
 # 该节点关联的 Replica 节点信息
      2) "6379"
      3) "60490389"
   2) 1) "172.19.42.233"
      2) "6379"
      3) "60490389"
```

```bash
127.0.0.1:6379> role
1) "slave"

 # 当前节点的角色
2) "172.19.42.231"
 # 该节点关联的 Master 节点信息
3) (integer) 6379
4) "connected"
5) (integer) 60503163
```

SAVE -

summary: Synchronously save the dataset to disk

since: 1.0.0

SHUTDOWN \[NOSAVE|SAVE]

summary: Synchronously save the dataset to disk and then shut down the server

since: 1.0.0

SLOWLOG subcommand \[argument]

summary: Manages the Redis slow queries log

since: 2.2.12

SYNC -

summary: Internal command used for replication

since: 1.0.0

TIME -

summary: Return the current server time

since: 2.6.0
