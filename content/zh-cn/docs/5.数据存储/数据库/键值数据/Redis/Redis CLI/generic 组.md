---
title: generic 组
linkTitle: generic 组
date: 2024-06-14T10:21
weight: 20
---

# 概述

> 参考：
>
> - https://redis.io/docs/latest/commands/?group=generic


# DEL - 删除一个 key

since: 1.0.0

**DEL key \[key ...]**

EXAMPLE

DUMP key

summary: Return a serialized version of the value stored at the specified key.

since: 2.6.0

# EXISTS - 判断指定的 key 是否存在

since: 1.0.0

**EXISTS key \[key ...]**

EXPIRE key seconds

summary: Set a key's time to live in seconds

since: 1.0.0

EXPIREAT key timestamp

summary: Set the expiration for a key as a UNIX timestamp

since: 1.2.0

# KEYS - 查找与指定 pattern 匹配到的所有 keys。

since: 1.0.0

**KEYS pattern**

**EXAMPLE**

- `keys *` # 获取所有的键

MIGRATE host port key| destination-db timeout \[COPY] \[REPLACE] \[KEYS key]

summary: Atomically transfer a key from a Redis instance to another one.

since: 2.6.0

MOVE key db

summary: Move a key to another database

since: 1.0.0

OBJECT subcommand \[arguments \[arguments ...]]

summary: Inspect the internals of Redis objects

since: 2.2.3

PERSIST key

summary: Remove the expiration from a key

since: 2.2.0

PEXPIRE key milliseconds

summary: Set a key's time to live in milliseconds

since: 2.6.0

PEXPIREAT key milliseconds-timestamp

summary: Set the expiration for a key as a UNIX timestamp specified in milliseconds

since: 2.6.0

PTTL key

summary: Get the time to live for a key in milliseconds

since: 2.6.0

RANDOMKEY -

summary: Return a random key from the keyspace

since: 1.0.0

RENAME key newkey

summary: Rename a key

since: 1.0.0

RENAMENX key newkey

summary: Rename a key, only if the new key does not exist

since: 1.0.0

RESTORE key ttl serialized-value \[REPLACE]

summary: Create a key using the provided serialized value, previously obtained using DUMP.

since: 2.6.0

# SCAN - 增量迭代 keys 空间(类似 KEYS 指令可以遍历所有 key)

since: 2.8.0

## Syntax(语法)

**SCAN cursor \[MATCH pattern] \[COUNT count] \[TYPE type]**


SORT key \[BY pattern] \[LIMIT offset count] \[GET pattern \[GET pattern ...]] \[ASC|DESC] \[ALPHA] \[STORE destination]

summary: Sort the elements in a list, set or sorted set

since: 1.0.0

TTL key

summary: Get the time to live for a key

since: 1.0.0

# TYPE - 确定已存储的 key 的类型

since: 1.0.0

TYPE 命令可以确定指定 key 的 [Redis 数据类型](/docs/5.数据存储/数据库/键值数据/Redis/Redis%20数据类型.md)

## Syntax(语法)

**TYPE key**


WAIT numslaves timeout

summary: Wait for the synchronous replication of all the write commands sent in the context of the current connection

since: 3.0.0

XRANGE key arg arg arg ...options...

summary: Help not available

since: not known

XGROUP arg arg ...options...

summary: Help not available

since: not known

LOLWUT arg ...options...

summary: Help not available

since: not known

XREADGROUP key arg arg arg arg arg arg ...options...

summary: Help not available

since: not known

PFDEBUG arg arg arg ...options...

summary: Help not available

since: not known

GEORADIUS_RO key arg arg arg arg arg ...options...

summary: Help not available

since: not known

HOST: arg ...options...

summary: Help not available

since: not known

XCLAIM key arg arg arg arg arg ...options...

summary: Help not available

since: not known

XPENDING key arg arg ...options...

summary: Help not available

since: not known

XACK key arg arg arg ...options...

summary: Help not available

since: not known

XDEL key arg arg ...options...

summary: Help not available

since: not known

XTRIM key arg ...options...

summary: Help not available

since: not known

REPLCONF arg ...options...

summary: Help not available

since: not known

XREVRANGE key arg arg arg ...options...

summary: Help not available

since: not known

BZPOPMIN key arg arg ...options...

summary: Help not available

since: not known

LATENCY arg arg ...options...

summary: Help not available

since: not known

TOUCH key arg ...options...

summary: Help not available

since: not known

XLEN key arg

summary: Help not available

since: not known

RESTORE-ASKING key arg arg arg ...options...

summary: Help not available

since: not known

GEORADIUSBYMEMBER_RO key arg arg arg arg ...options...

summary: Help not available

since: not known

ZPOPMIN key arg ...options...

summary: Help not available

since: not known

POST arg ...options...

summary: Help not available

since: not known

XREAD key arg arg arg ...options...

summary: Help not available

since: not known

PSYNC arg arg arg

summary: Help not available

since: not known

XINFO arg arg ...options...

summary: Help not available

since: not known

BZPOPMAX key arg arg ...options...

summary: Help not available

since: not known

ASKING arg

summary: Help not available

since: not known

SWAPDB arg arg arg

summary: Help not available

since: not known

SUBSTR key arg arg arg

summary: Help not available

since: not known

XSETID key arg arg

summary: Help not available

since: not known

UNLINK key arg ...options...

summary: Help not available

since: not known

MODULE arg arg ...options...

summary: Help not available

since: not known

PFSELFTEST arg

summary: Help not available

since: not known

REPLICAOF arg arg arg

summary: Help not available

since: not known

MEMORY arg arg ...options...

summary: Help not available

since: not known

XADD key arg arg arg arg ...options...

summary: Help not available

since: not known

ZPOPMAX key arg ...options...

summary: Help not available

since: not known
