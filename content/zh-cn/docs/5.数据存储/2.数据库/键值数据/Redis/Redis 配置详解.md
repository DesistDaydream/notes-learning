---
title: Redis 配置详解
---

# 概述

> 参考：
> - [官方文档](https://redis.io/topics/config)

Redis 可以在不使用配置文件的情况下使用内置的默认配置启动。但是一般情况，都会使用一个 Redis 的配置文件(文件名通常是 redis.conf)来启动 Redis。Redis 启动后，会将 redis.conf 文件的内容加载到内存中，通过 Redis 客户端的 **config get \*** 命令，即可获取当前已经加载到内存中的配置。

```bash
127.0.0.1:6379> config get *
  1) "dbfilename"
  2) "dump.rdb"
  3) "requirepass"
  4) ""
  5) "masterauth"
  6) ""
....... 后略
```

配置文件的写法非常简单。redis.conf 由 **Directives(指令)** 组成，每条指令一行。而 Directives 分为两部分

- **Keyword(关键字)** # 该指令的含义
- **Arguments(参数)** # Redis 执行该指令时的行为

格式如下：

    # 关键字 参数(多个参数以空格分隔)
    Keyword Argument1 Argument2 ... ArugmentN

## 通过命令函参数传递配置

通过命令行传递参数的格式与 redis.conf 文件中配置格式完全相同，只不过关键字前面有个 `--` 前缀。比如：

    redis-server --port 6380 --replicaof 127.0.0.1 6379

生成的内存中配置如下：

    127.0.0.1:6380> config get "replicaof"
    1) "replicaof"
    2) "127.0.0.1 6379"
    127.0.0.1:6380> config get "port"
    1) "port"
    2) "6380"

### 其他基本示例

```bash
Usage: ./redis-server [/path/to/redis.conf] [options]
       ./redis-server - (read config from stdin)
       ./redis-server -v or --version
       ./redis-server -h or --help
       ./redis-server --test-memory <megabytes>

Examples:
       ./redis-server (run the server with default conf)
       ./redis-server /etc/redis/6379.conf
       ./redis-server --port 7777
       ./redis-server --port 7777 --replicaof 127.0.0.1 8888
       ./redis-server /etc/myredis.conf --loglevel verbose
```

## 启动 Redis 后修改配置

Redis 支持在线热更新配置，可以通过 config set ARGUMENT 命令来更改当前配置，若想要配置永久生效，则使用 config rewrite 命令将内存中的配置重写到启动 Redis 时指定的 redis.conf 文件中。但是，并非配置中的所有指令都支持这种方式。

## 简单的配置文件示例

    dir "/data"
    port 6379
    maxmemory 1G
    maxmemory-policy volatile-lru
    min-replicas-max-lag 5
    min-replicas-to-write 1
    rdbchecksum yes
    rdbcompression yes
    repl-diskless-sync yes
    save 900 1
    requirepass redis
    masterauth "redis"

    replica-announce-port 6379
    replica-announce-ip "10.105.180.122"

# redis.conf 文件详解

## Includes 配置环境

- **include /PATH/TO/FILE** # Redis 启动时，除了加载 redis.conf 文件外，还会加载 include 指令指定的文件。

## Network 配置环境

- **bind 127.0.0.1** # 监听的地址
- **port 6379** # redis 监听的端口，默认 6379
- **tcp-backlog 511** # tcp 的等待队列
- **timeout 0** # 客户端连接超时时长。`默认值：0`，不会超时

## General 配置环境

- **daemonize yes|no** # 指定是否在后台运行
- **databases 16** # 可使用的 databases，默认 16 个
- **logfile /PATH/TO/FILE** # 指定 redis 记录日志文件位置

## Snapshotting 配置环境 RDB 功能

- **save TIME NUM** # 在 TIME 秒内有 NUM 个键改变，就做一次快照
- **stop-writes-on-bgsave-error yes|no** # 当 RDB 持久化出现错误后，是否依然进行继续进行工作，yes：不能进行工作，no：可以继续进行工作，可以通过 info 中的 rdb_last_bgsave_status 了解 RDB 持久化是否有错误
- **rdbcompression yes|no** # 是否压缩 rdb 文件，rdb 文件压缩使用 LZF 压缩算法。压缩需要一些 cpu 的消耗；不压缩需要更多的磁盘空间
- **rdbchecksum yes|no** # 是否校验 rdb 文件。从 rdb 格式的第五个版本开始，在 rdb 文件的末尾会带上 CRC64 的校验和。这跟有利于文件的容错性，但是在保存 rdb 文件的时候，会有大概 10%的性能损耗，所以如果你追求高性能，可以关闭该配置。
- **dbfilename FileName** # 指定 snapshot 文件的文件名。默认为 dump.rbd
- **dir PATH** # 指定 snapshot 文件的保存路径(注意：PATH 是目录，不是文件，具体文件名通过 dbfilename 关键字配置)。默认为 /var/lib/redis/ 目录

## Replication 配置环境

- **replicof 192.168.1.2 6379** # 启动主从模式，并设定自己为从服务器，主服务器 IP 为 192.168.1.2，主服务器端口为 6379
- **slave-read-only no** # 作为从服务器是否只读，默认不只读

## Security 配置环境

- **requirepass PASSWORD** # 配置认证密码为 PASSWORD

## Limits 配置环境

- **maxmemory BYTES** # 指定 redis 可使用的最大内存量，单位是 bytes。如果达到限额，则需要配合 maxmemory-policy 配置指定的策略删除 key。note：slave 的输出缓冲区是不计算在 maxmemory 内的。所以为了防止主机内存使用完，建议设置的 maxmemory 需要更小一些。
- **maxmemory-policy POLICY** # 指定 redis 超过内存限额之后的策略，包括以下几种
  - volatile-lru：利用 LRU 算法移除设置过过期时间的 key。
  - volatile-random：随机移除设置过过期时间的 key。
  - volatile-ttl：移除即将过期的 key，根据最近过期时间来删除（辅以 TTL）
  - allkeys-lru：利用 LRU 算法移除任何 key。
  - allkeys-random：随机移除任何 key。
  - noeviction：不移除任何 key，只是返回一个写错误。
  - Note:上面的这些驱逐策略，如果 redis 没有合适的 key 驱逐，对于写命令，还是会返回错误。redis 将不再接收写请求，只接收 get 请求。写命令包括：set setnx setex append incr decr rpush lpush rpushx lpushx linsert lset rpoplpush sadd sinter sinterstore sunion sunionstore sdiff sdiffstore zadd zincrby zunionstore zinterstore hset hsetnx hmset hincrby incrby decrby getset mset msetnx exec sort。

## Append Only Mode 配置环境 AOF 功能配置

## LUA Scripting 配置环境

## Redis Cluster 配置环境

## SLOW LOG 配置环境

## LATENCY MONITOR 配置环境

## EVENT NOTIFICATION 配置环境

## ADVANCED CONFIG 配置环境
