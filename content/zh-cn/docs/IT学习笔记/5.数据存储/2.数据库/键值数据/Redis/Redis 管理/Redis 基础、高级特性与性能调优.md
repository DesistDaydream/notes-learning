---
title: Redis 基础、高级特性与性能调优
---

# 概述

> 参考：
> - [公众号,Redis 基础、高级特性与性能调优](https://mp.weixin.qq.com/s/yUtZrmaE9EbhjoGzMCkFAw)
> - [公众号,万字总结，Redis 性能问题排查解决手册](https://mp.weixin.qq.com/s/7ZlEkS63yl_DahOzq6LxQA)

本文将从 Redis 的基本特性入手，通过讲述 Redis 的数据结构和主要命令对 Redis 的基本能力进行直观介绍。之后概览 Redis 提供的高级能力，并在部署、维护、性能调优等多个方面进行更深入的介绍和指导。
本文适合使用 Redis 的普通开发人员，以及对 Redis 进行选型、架构设计和性能调优的架构设计人员。

## **目录**

- 概述
- Redis 的数据结构和相关常用命令
- 数据持久化
- 内存管理与数据淘汰机制
- Pipelining
- 事务与 Scripting
- Redis 性能调优
- 主从复制与集群分片
- Redis Java 客户端的选择

##

## \*\*

## **概述**

Redis 是一个开源的，基于内存的结构化数据存储媒介，可以作为数据库、缓存服务或消息服务使用。
Redis 支持多种数据结构，包括字符串、哈希表、链表、集合、有序集合、位图、Hyperloglogs 等。
Redis 具备 LRU 淘汰、事务实现、以及不同级别的硬盘持久化等能力，并且支持副本集和通过 Redis Sentinel 实现的高可用方案，同时还支持通过 Redis Cluster 实现的数据自动分片能力。Redis 的主要功能都基于单线程模型实现，也就是说 Redis 使用一个线程来服务所有的客户端请求，同时 Redis 采用了非阻塞式 IO，并精细地优化各种命令的算法时间复杂度，这些信息意味着：

- Redis 是线程安全的（因为只有一个线程），其所有操作都是原子的，不会因并发产生数据异常
- Redis 的速度非常快（因为使用非阻塞式 IO，且大部分命令的算法时间复杂度都是 O(1))
- 使用高耗时的 Redis 命令是很危险的，会占用唯一的一个线程的大量处理时间，导致所有的请求都被拖慢。（例如时间复杂度为 O(N)的 KEYS 命令，严格禁止在生产环境中使用）

##

## **Redis 的数据结构和相关常用命令**

本节中将介绍 Redis 支持的主要数据结构，以及相关的常用 Redis 命令。本节只对 Redis 命令进行扼要的介绍，且只列出了较常用的命令。如果想要了解完整的 Redis 命令集，或了解某个命令的详细使用方法，请参考官方文档：<https://redis.io/commands>

###

### **Key**

Redis 采用 Key-Value 型的基本数据结构，任何二进制序列都可以作为 Redis 的 Key 使用（例如普通的字符串或一张 JPEG 图片）
关于 Key 的一些注意事项：

- 不要使用过长的 Key。例如使用一个 1024 字节的 key 就不是一个好主意，不仅会消耗更多的内存，还会导致查找的效率降低
- Key 短到缺失了可读性也是不好的，例如"u1000flw"比起"user:1000:followers"来说，节省了寥寥的存储空间，却引发了可读性和可维护性上的麻烦
- 最好使用统一的规范来设计 Key，比如"object-type:id:attr"，以这一规范设计出的 Key 可能是"user:1000"或"comment:1234:reply-to"
- Redis 允许的最大 Key 长度是 512MB（对 Value 的长度限制也是 512MB）

###

### **String**

String 是 Redis 的基础数据类型，Redis 没有 Int、Float、Boolean 等数据类型的概念，所有的基本类型在 Redis 中都以 String 体现。与 String 相关的常用命令：

- **SET**：为一个 key 设置 value，可以配合 EX/PX 参数指定 key 的有效期，通过 NX/XX 参数针对 key 是否存在的情况进行区别操作，时间复杂度 O(1)
- **GET**：获取某个 key 对应的 value，时间复杂度 O(1)
- **GETSET**：为一个 key 设置 value，并返回该 key 的原 value，时间复杂度 O(1)
- **MSET**：为多个 key 设置 value，时间复杂度 O(N)
- **MSETNX**：同 MSET，如果指定的 key 中有任意一个已存在，则不进行任何操作，时间复杂度 O(N)
- **MGET**：获取多个 key 对应的 value，时间复杂度 O(N)

上文提到过，Redis 的基本数据类型只有 String，但 Redis 可以把 String 作为整型或浮点型数字来使用，主要体现在 INCR、DECR 类的命令上：

- **INCR**：将 key 对应的 value 值自增 1，并返回自增后的值。只对可以转换为整型的 String 数据起作用。时间复杂度 O(1)
- **INCRBY**：将 key 对应的 value 值自增指定的整型数值，并返回自增后的值。只对可以转换为整型的 String 数据起作用。时间复杂度 O(1)
- **DECR/DECRBY**：同 INCR/INCRBY，自增改为自减。

INCR/DECR 系列命令要求操作的 value 类型为 String，并可以转换为 64 位带符号的整型数字，否则会返回错误。
也就是说，进行 INCR/DECR 系列命令的 value，必须在\[-2^63 ~ 2^63 - 1]范围内。前文提到过，Redis 采用单线程模型，天然是线程安全的，这使得 INCR/DECR 命令可以非常便利的实现高并发场景下的精确控制。

####

#### **例 1：库存控制**

在高并发场景下实现库存余量的精准校验，确保不出现超卖的情况。设置库存总量：

-

<!---->

    SET inv:remain "100"

库存扣减+余量校验：

-

<!---->

    DECR inv:remain

当 DECR 命令返回值大于等于 0 时，说明库存余量校验通过，如果返回小于 0 的值，则说明库存已耗尽。
假设同时有 300 个并发请求进行库存扣减，Redis 能够确保这 300 个请求分别得到 99 到-200 的返回值，每个请求得到的返回值都是唯一的，绝对不会找出现两个请求得到一样的返回值的情况。

####

#### **例 2：自增序列生成**

实现类似于 RDBMS 的 Sequence 功能，生成一系列唯一的序列号设置序列起始值：

-

<!---->

    SET sequence "10000"

获取一个序列值：

-

<!---->

    INCR sequence

直接将返回值作为序列使用即可。
获取一批（如 100 个）序列值：

-

<!---->

    INCRBY sequence 100

假设返回值为 N，那么\[N - 99 ~ N]的数值都是可用的序列值。
当多个客户端同时向 Redis 申请自增序列时，Redis 能够确保每个客户端得到的序列值或序列范围都是全局唯一的，绝对不会出现不同客户端得到了重复的序列值的情况。

### **List**

Redis 的 List 是链表型的数据结构，可以使用 LPUSH/RPUSH/LPOP/RPOP 等命令在 List 的两端执行插入元素和弹出元素的操作。虽然 List 也支持在特定 index 上插入和读取元素的功能，但其时间复杂度较高（O(N)），应小心使用。与 List 相关的常用命令：

- **LPUSH**：向指定 List 的左侧（即头部）插入 1 个或多个元素，返回插入后的 List 长度。时间复杂度 O(N)，N 为插入元素的数量
- **RPUSH**：同 LPUSH，向指定 List 的右侧（即尾部）插入 1 或多个元素
- **LPOP**：从指定 List 的左侧（即头部）移除一个元素并返回，时间复杂度 O(1)
- **RPOP**：同 LPOP，从指定 List 的右侧（即尾部）移除 1 个元素并返回
- **LPUSHX/RPUSHX**：与 LPUSH/RPUSH 类似，区别在于，LPUSHX/RPUSHX 操作的 key 如果不存在，则不会进行任何操作
- **LLEN**：返回指定 List 的长度，时间复杂度 O(1)
- **LRANGE**：返回指定 List 中指定范围的元素（双端包含，即 LRANGE key 0 10 会返回 11 个元素），时间复杂度 O(N)。应尽可能控制一次获取的元素数量，一次获取过大范围的 List 元素会导致延迟，同时对长度不可预知的 List，避免使用 LRANGE key 0 -1 这样的完整遍历操作。

应谨慎使用的 List 相关命令：

- **LINDEX**：返回指定 List 指定 index 上的元素，如果 index 越界，返回 nil。index 数值是回环的，即-1 代表 List 最后一个位置，-2 代表 List 倒数第二个位置。时间复杂度 O(N)
- **LSET**：将指定 List 指定 index 上的元素设置为 value，如果 index 越界则返回错误，时间复杂度 O(N)，如果操作的是头/尾部的元素，则时间复杂度为 O(1)
- **LINSERT**：向指定 List 中指定元素之前/之后插入一个新元素，并返回操作后的 List 长度。如果指定的元素不存在，返回-1。如果指定 key 不存在，不会进行任何操作，时间复杂度 O(N)

由于 Redis 的 List 是链表结构的，上述的三个命令的算法效率较低，需要对 List 进行遍历，命令的耗时无法预估，在 List 长度大的情况下耗时会明显增加，应谨慎使用。换句话说，Redis 的 List 实际是设计来用于实现队列，而不是用于实现类似 ArrayList 这样的列表的。如果你不是想要实现一个双端出入的队列，那么请尽量不要使用 Redis 的 List 数据结构。为了更好支持队列的特性，Redis 还提供了一系列阻塞式的操作命令，如 BLPOP/BRPOP 等，能够实现类似于 BlockingQueue 的能力，即在 List 为空时，阻塞该连接，直到 List 中有对象可以出队时再返回。针对阻塞类的命令，此处不做详细探讨，请参考官方文档（<https://redis.io/topics/data-types-intro>） 中"Blocking operations on lists"一节。

### **Hash**

Hash 即哈希表，Redis 的 Hash 和传统的哈希表一样，是一种 field-value 型的数据结构，可以理解成将 HashMap 搬入 Redis。
Hash 非常适合用于表现对象类型的数据，用 Hash 中的 field 对应对象的 field 即可。
Hash 的优点包括：

- 可以实现二元查找，如"查找 ID 为 1000 的用户的年龄"
- 比起将整个对象序列化后作为 String 存储的方法，Hash 能够有效地减少网络传输的消耗
- 当使用 Hash 维护一个集合时，提供了比 List 效率高得多的随机访问命令

与 Hash 相关的常用命令：

- **HSET**：将 key 对应的 Hash 中的 field 设置为 value。如果该 Hash 不存在，会自动创建一个。时间复杂度 O(1)
- **HGET**：返回指定 Hash 中 field 字段的值，时间复杂度 O(1)
- **HMSET/HMGET**：同 HSET 和 HGET，可以批量操作同一个 key 下的多个 field，时间复杂度：O(N)，N 为一次操作的 field 数量
- **HSETNX**：同 HSET，但如 field 已经存在，HSETNX 不会进行任何操作，时间复杂度 O(1)
- **HEXISTS**：判断指定 Hash 中 field 是否存在，存在返回 1，不存在返回 0，时间复杂度 O(1)
- **HDEL**：删除指定 Hash 中的 field（1 个或多个），时间复杂度：O(N)，N 为操作的 field 数量
- **HINCRBY**：同 INCRBY 命令，对指定 Hash 中的一个 field 进行 INCRBY，时间复杂度 O(1)

应谨慎使用的 Hash 相关命令：

- **HGETALL**：返回指定 Hash 中所有的 field-value 对。返回结果为数组，数组中 field 和 value 交替出现。时间复杂度 O(N)
- **HKEYS/HVALS**：返回指定 Hash 中所有的 field/value，时间复杂度 O(N)

上述三个命令都会对 Hash 进行完整遍历，Hash 中的 field 数量与命令的耗时线性相关，对于尺寸不可预知的 Hash，应严格避免使用上面三个命令，而改为使用 HSCAN 命令进行游标式的遍历，具体请见 <https://redis.io/commands/scan>

###

### **Set**

Redis Set 是无序的，不可重复的 String 集合。与 Set 相关的常用命令：

- **SADD**：向指定 Set 中添加 1 个或多个 member，如果指定 Set 不存在，会自动创建一个。时间复杂度 O(N)，N 为添加的 member 个数
- **SREM**：从指定 Set 中移除 1 个或多个 member，时间复杂度 O(N)，N 为移除的 member 个数
- **SRANDMEMBER**：从指定 Set 中随机返回 1 个或多个 member，时间复杂度 O(N)，N 为返回的 member 个数
- **SPOP**：从指定 Set 中随机移除并返回 count 个 member，时间复杂度 O(N)，N 为移除的 member 个数
- **SCARD**：返回指定 Set 中的 member 个数，时间复杂度 O(1)
- **SISMEMBER**：判断指定的 value 是否存在于指定 Set 中，时间复杂度 O(1)
- **SMOVE**：将指定 member 从一个 Set 移至另一个 Set

慎用的 Set 相关命令：

- **SMEMBERS**：返回指定 Hash 中所有的 member，时间复杂度 O(N)
- **SUNION/SUNIONSTORE**：计算多个 Set 的并集并返回/存储至另一个 Set 中，时间复杂度 O(N)，N 为参与计算的所有集合的总 member 数
- **SINTER/SINTERSTORE**：计算多个 Set 的交集并返回/存储至另一个 Set 中，时间复杂度 O(N)，N 为参与计算的所有集合的总 member 数
- **SDIFF/SDIFFSTORE**：计算 1 个 Set 与 1 或多个 Set 的差集并返回/存储至另一个 Set 中，时间复杂度 O(N)，N 为参与计算的所有集合的总 member 数

上述几个命令涉及的计算量大，应谨慎使用，特别是在参与计算的 Set 尺寸不可知的情况下，应严格避免使用。可以考虑通过 SSCAN 命令遍历获取相关 Set 的全部 member（具体请见 <https://redis.io/commands/scan> ），如果需要做并集/交集/差集计算，可以在客户端进行，或在不服务实时查询请求的 Slave 上进行。

###

### **Sorted Set**

Redis Sorted Set 是有序的、不可重复的 String 集合。Sorted Set 中的每个元素都需要指派一个分数(score)，Sorted Set 会根据 score 对元素进行升序排序。如果多个 member 拥有相同的 score，则以字典序进行升序排序。Sorted Set 非常适合用于实现排名。Sorted Set 的主要命令：

- **ZADD**：向指定 Sorted Set 中添加 1 个或多个 member，时间复杂度 O(Mlog(N))，M 为添加的 member 数量，N 为 Sorted Set 中的 member 数量
- **ZREM**：从指定 Sorted Set 中删除 1 个或多个 member，时间复杂度 O(Mlog(N))，M 为删除的 member 数量，N 为 Sorted Set 中的 member 数量
- **ZCOUNT**：返回指定 Sorted Set 中指定 score 范围内的 member 数量，时间复杂度：O(log(N))
- **ZCARD**：返回指定 Sorted Set 中的 member 数量，时间复杂度 O(1)
- **ZSCORE**：返回指定 Sorted Set 中指定 member 的 score，时间复杂度 O(1)
- **ZRANK/ZREVRANK**：返回指定 member 在 Sorted Set 中的排名，ZRANK 返回按升序排序的排名，ZREVRANK 则返回按降序排序的排名。时间复杂度 O(log(N))
- **ZINCRBY**：同 INCRBY，对指定 Sorted Set 中的指定 member 的 score 进行自增，时间复杂度 O(log(N))

慎用的 Sorted Set 相关命令：

- **ZRANGE/ZREVRANGE**：返回指定 Sorted Set 中指定排名范围内的所有 member，ZRANGE 为按 score 升序排序，ZREVRANGE 为按 score 降序排序，时间复杂度 O(log(N)+M)，M 为本次返回的 member 数
- **ZRANGEBYSCORE/ZREVRANGEBYSCORE**：返回指定 Sorted Set 中指定 score 范围内的所有 member，返回结果以升序/降序排序，min 和 max 可以指定为-inf 和+inf，代表返回所有的 member。时间复杂度 O(log(N)+M)
- **ZREMRANGEBYRANK/ZREMRANGEBYSCORE**：移除 Sorted Set 中指定排名范围/指定 score 范围内的所有 member。时间复杂度 O(log(N)+M)

上述几个命令，应尽量避免传递\[0 -1]或\[-inf +inf]这样的参数，来对 Sorted Set 做一次性的完整遍历，特别是在 Sorted Set 的尺寸不可预知的情况下。可以通过 ZSCAN 命令来进行游标式的遍历（具体请见 <https://redis.io/commands/scan> ），或通过 LIMIT 参数来限制返回 member 的数量（适用于 ZRANGEBYSCORE 和 ZREVRANGEBYSCORE 命令），以实现游标式的遍历。

###

### **Bitmap 和 HyperLogLog**

Redis 的这两种数据结构相较之前的并不常用，在本文中只做简要介绍，如想要详细了解这两种数据结构与其相关的命令，请参考官方文档<https://redis.io/topics/data-types-intro> 中的相关章节 Bitmap 在 Redis 中不是一种实际的数据类型，而是一种将 String 作为 Bitmap 使用的方法。可以理解为将 String 转换为 bit 数组。使用 Bitmap 来存储 true/false 类型的简单数据极为节省空间。HyperLogLogs 是一种主要用于数量统计的数据结构，它和 Set 类似，维护一个不可重复的 String 集合，但是 HyperLogLogs 并不维护具体的 member 内容，只维护 member 的个数。也就是说，HyperLogLogs 只能用于计算一个集合中不重复的元素数量，所以它比 Set 要节省很多内存空间。

###

### **其他常用命令**

- **EXISTS**：判断指定的 key 是否存在，返回 1 代表存在，0 代表不存在，时间复杂度 O(1)
- **DEL**：删除指定的 key 及其对应的 value，时间复杂度 O(N)，N 为删除的 key 数量
- **EXPIRE/PEXPIRE**：为一个 key 设置有效期，单位为秒或毫秒，时间复杂度 O(1)
- **TTL/PTTL**：返回一个 key 剩余的有效时间，单位为秒或毫秒，时间复杂度 O(1)
- **RENAME/RENAMENX**：将 key 重命名为 newkey。使用 RENAME 时，如果 newkey 已经存在，其值会被覆盖；使用 RENAMENX 时，如果 newkey 已经存在，则不会进行任何操作，时间复杂度 O(1)
- **TYPE**：返回指定 key 的类型，string, list, set, zset, hash。时间复杂度 O(1)
- **CONFIG GET**：获得 Redis 某配置项的当前值，可以使用\*通配符，时间复杂度 O(1)
- **CONFIG SET**：为 Redis 某个配置项设置新值，时间复杂度 O(1)
- **CONFIG REWRITE**：让 Redis 重新加载 redis.conf 中的配置

##

## **数据持久化**

Redis 提供了将数据定期自动持久化至硬盘的能力，包括 RDB 和 AOF 两种方案，两种方案分别有其长处和短板，可以配合起来同时运行，确保数据的稳定性。

### 必须使用数据持久化吗？

Redis 的数据持久化机制是可以关闭的。如果你只把 Redis 作为缓存服务使用，Redis 中存储的所有数据都不是该数据的主体而仅仅是同步过来的备份，那么可以关闭 Redis 的数据持久化机制。
但通常来说，仍然建议至少开启 RDB 方式的数据持久化，因为：

- RDB 方式的持久化几乎不损耗 Redis 本身的性能，在进行 RDB 持久化时，Redis 主进程唯一需要做的事情就是 fork 出一个子进程，所有持久化工作都由子进程完成
- Redis 无论因为什么原因 crash 掉之后，重启时能够自动恢复到上一次 RDB 快照中记录的数据。这省去了手工从其他数据源（如 DB）同步数据的过程，而且要比其他任何的数据恢复方式都要快
- 现在硬盘那么大，真的不缺那一点地方

###

### **RDB**

采用 RDB 持久方式，Redis 会定期保存数据快照至一个 rbd 文件中，并在启动时自动加载 rdb 文件，恢复之前保存的数据。可以在配置文件中配置 Redis 进行快照保存的时机：

-

<!---->

    save [seconds] [changes]

意为在\[seconds]秒内如果发生了\[changes]次数据修改，则进行一次 RDB 快照保存，例如

-

<!---->

    save 60 100

会让 Redis 每 60 秒检查一次数据变更情况，如果发生了 100 次或以上的数据变更，则进行 RDB 快照保存。
可以配置多条 save 指令，让 Redis 执行多级的快照保存策略。
Redis 默认开启 RDB 快照，默认的 RDB 策略如下：

-

-

-

<!---->

    save 900 1save 300 10save 60 10000

\`\`也可以通过**BGSAVE**命令手工触发 RDB 快照保存。**RDB 的优点：**

- 对性能影响最小。如前文所述，Redis 在保存 RDB 快照时会 fork 出子进程进行，几乎不影响 Redis 处理客户端请求的效率。
- 每次快照会生成一个完整的数据快照文件，所以可以辅以其他手段保存多个时间点的快照（例如把每天 0 点的快照备份至其他存储媒介中），作为非常可靠的灾难恢复手段。
- 使用 RDB 文件进行数据恢复比使用 AOF 要快很多。

**RDB 的缺点：**

- 快照是定期生成的，所以在 Redis crash 时或多或少会丢失一部分数据。
- 如果数据集非常大且 CPU 不够强（比如单核 CPU），Redis 在 fork 子进程时可能会消耗相对较长的时间（长至 1 秒），影响这期间的客户端请求。

###

### **AOF**

采用 AOF 持久方式时，Redis 会把每一个写请求都记录在一个日志文件里。在 Redis 重启时，会把 AOF 文件中记录的所有写操作顺序执行一遍，确保数据恢复到最新。AOF 默认是关闭的，如要开启，进行如下配置：

-

<!---->

    appendonly yes

AOF 提供了三种 fsync 配置，always/everysec/no，通过配置项\[appendfsync]指定：

- appendfsync no：不进行 fsync，将 flush 文件的时机交给 OS 决定，速度最快
- appendfsync always：每写入一条日志就进行一次 fsync 操作，数据安全性最高，但速度最慢
- appendfsync everysec：折中的做法，交由后台线程每秒 fsync 一次

随着 AOF 不断地记录写操作日志，必定会出现一些无用的日志，例如某个时间点执行了命令**SET key1 "abc"**，在之后某个时间点又执行了**SET key1 "bcd"**，那么第一条命令很显然是没有用的。大量的无用日志会让 AOF 文件过大，也会让数据恢复的时间过长。
所以 Redis 提供了 AOF rewrite 功能，可以重写 AOF 文件，只保留能够把数据恢复到最新状态的最小写操作集。
AOF rewrite 可以通过**BGREWRITEAOF**命令触发，也可以配置 Redis 定期自动进行：

-

-

<!---->

    auto-aof-rewrite-percentage 100auto-aof-rewrite-min-size 64mb

`<br />`上面两行配置的含义是，Redis 在每次 AOF rewrite 时，会记录完成 rewrite 后的 AOF 日志大小，当 AOF 日志大小在该基础上增长了 100%后，自动进行 AOF rewrite。同时如果增长的大小没有达到 64mb，则不会进行 rewrite。**AOF 的优点：**

- 最安全，在启用 appendfsync always 时，任何已写入的数据都不会丢失，使用在启用 appendfsync everysec 也至多只会丢失 1 秒的数据。
- AOF 文件在发生断电等问题时也不会损坏，即使出现了某条日志只写入了一半的情况，也可以使用 redis-check-aof 工具轻松修复。
- AOF 文件易读，可修改，在进行了某些错误的数据清除操作后，只要 AOF 文件没有 rewrite，就可以把 AOF 文件备份出来，把错误的命令删除，然后恢复数据。

**AOF 的缺点：**

- AOF 文件通常比 RDB 文件更大
- 性能消耗比 RDB 高
- 数据恢复速度比 RDB 慢

##

## **内存管理与数据淘汰机制**

### **最大内存设置**

默认情况下，在 32 位 OS 中，Redis 最大使用 3GB 的内存，在 64 位 OS 中则没有限制。在使用 Redis 时，应该对数据占用的最大空间有一个基本准确的预估，并为 Redis 设定最大使用的内存。否则在 64 位 OS 中 Redis 会无限制地占用内存（当物理内存被占满后会使用 swap 空间），容易引发各种各样的问题。通过如下配置控制 Redis 使用的最大内存：

-

<!---->

    maxmemory 100mb

在内存占用达到了 maxmemory 后，再向 Redis 写入数据时，Redis 会：

- 根据配置的数据淘汰策略尝试淘汰数据，释放空间
- 如果没有数据可以淘汰，或者没有配置数据淘汰策略，那么 Redis 会对所有写请求返回错误，但读请求仍然可以正常执行

在为 Redis 设置 maxmemory 时，需要注意：

- 如果采用了 Redis 的主从同步，主节点向从节点同步数据时，会占用掉一部分内存空间，如果 maxmemory 过于接近主机的可用内存，导致数据同步时内存不足。所以设置的 maxmemory 不要过于接近主机可用的内存，留出一部分预留用作主从同步。

###

### **数据淘汰机制**

Redis 提供了 5 种数据淘汰策略：

- volatile-lru：使用 LRU 算法进行数据淘汰（淘汰上次使用时间最早的，且使用次数最少的 key），只淘汰设定了有效期的 key
- allkeys-lru：使用 LRU 算法进行数据淘汰，所有的 key 都可以被淘汰
- volatile-random：随机淘汰数据，只淘汰设定了有效期的 key
- allkeys-random：随机淘汰数据，所有的 key 都可以被淘汰
- volatile-ttl：淘汰剩余有效期最短的 key

最好为 Redis 指定一种有效的数据淘汰策略以配合 maxmemory 设置，避免在内存使用满后发生写入失败的情况。一般来说，推荐使用的策略是 volatile-lru，并辨识 Redis 中保存的数据的重要性。对于那些重要的，绝对不能丢弃的数据（如配置类数据等），应不设置有效期，这样 Redis 就永远不会淘汰这些数据。对于那些相对不是那么重要的，并且能够热加载的数据（比如缓存最近登录的用户信息，当在 Redis 中找不到时，程序会去 DB 中读取），可以设置上有效期，这样在内存不够时 Redis 就会淘汰这部分数据。配置方法：

-

<!---->

    maxmemory-policy volatile-lru   #默认是noeviction，即不进行数据淘汰

**Pipelining**

### **Pipelining**

Redis 提供许多批量操作的命令，如 MSET/MGET/HMSET/HMGET 等等，这些命令存在的意义是减少维护网络连接和传输数据所消耗的资源和时间。
例如连续使用 5 次 SET 命令设置 5 个不同的 key，比起使用一次 MSET 命令设置 5 个不同的 key，效果是一样的，但前者会消耗更多的 RTT(Round Trip Time)时长，永远应优先使用后者。然而，如果客户端要连续执行的多次操作无法通过 Redis 命令组合在一起，例如：

-

-

-

<!---->

    SET a "abc"INCR bHSET c name "hi"

`<br />`此时便可以使用 Redis 提供的 pipelining 功能来实现在一次交互中执行多条命令。
使用 pipelining 时，只需要从客户端一次向 Redis 发送多条命令（以\r\n）分隔，Redis 就会依次执行这些命令，并且把每个命令的返回按顺序组装在一起一次返回，比如：

-

-

-

-

<!---->

    $ (printf "PING\r\nPING\r\nPING\r\n"; sleep 1) | nc localhost 6379+PONG+PONG+PONG

`<br />`大部分的 Redis 客户端都对 Pipelining 提供支持，所以开发者通常并不需要自己手工拼装命令列表。

##### **Pipelining 的局限性**

Pipelining 只能用于执行**连续且无相关性**的命令，当某个命令的生成需要依赖于前一个命令的返回时，就无法使用 Pipelining 了。通过 Scripting 功能，可以规避这一局限性

## **事务与 Scripting**

Pipelining 能够让 Redis 在一次交互中处理多条命令，然而在一些场景下，我们可能需要在此基础上确保这一组命令是连续执行的。比如获取当前累计的 PV 数并将其清 0

-

-

-

-

<!---->

    > GET vCount12384> SET vCount 0OK

`<br />`如果在 GET 和 SET 命令之间插进来一个 INCR vCount，就会使客户端拿到的 vCount 不准确。Redis 的事务可以确保复数命令执行时的原子性。也就是说 Redis 能够保证：一个事务中的一组命令是绝对连续执行的，在这些命令执行完成之前，绝对不会有来自于其他连接的其他命令插进去执行。通过 MULTI 和 EXEC 命令来把这两个命令加入一个事务中：

-

-

-

-

-

-

-

-

-

<!---->

    > MULTIOK> GET vCountQUEUED> SET vCount 0QUEUED> EXEC1) 123842) OK

`<br />`Redis 在接收到 MULTI 命令后便会开启一个事务，这之后的所有读写命令都会保存在队列中但并不执行，直到接收到 EXEC 命令后，Redis 会把队列中的所有命令连续顺序执行，并以数组形式返回每个命令的返回结果。可以使用 DISCARD 命令放弃当前的事务，将保存的命令队列清空。需要注意的是，**Redis 事务不支持回滚**：
如果一个事务中的命令出现了语法错误，大部分客户端驱动会返回错误，2.6.5 版本以上的 Redis 也会在执行 EXEC 时检查队列中的命令是否存在语法错误，如果存在，则会自动放弃事务并返回错误。
但如果一个事务中的命令有非语法类的错误（比如对 String 执行 HSET 操作），无论客户端驱动还是 Redis 都无法在真正执行这条命令之前发现，所以事务中的所有命令仍然会被依次执行。在这种情况下，会出现一个事务中部分命令成功部分命令失败的情况，然而与 RDBMS 不同，Redis 不提供事务回滚的功能，所以只能通过其他方法进行数据的回滚。

###

### **通过事务实现 CAS**

Redis 提供了 WATCH 命令与事务搭配使用，实现 CAS 乐观锁的机制。假设要实现将某个商品的状态改为已售：

-

-

<!---->

    if(exec(HGET stock:1001 state) == "in stock")    exec(HSET stock:1001 state "sold");

`<br />`这一伪代码执行时，无法确保并发安全性，有可能多个客户端都获取到了"in stock"的状态，导致一个库存被售卖多次。使用 WATCH 命令和事务可以解决这一问题：

-

-

-

-

-

-

<!---->

    exec(WATCH stock:1001);if(exec(HGET stock:1001 state) == "in stock") {    exec(MULTI);    exec(HSET stock:1001 state "sold");    exec(EXEC);}

`<br />`WATCH 的机制是：在事务 EXEC 命令执行时，Redis 会检查被 WATCH 的 key，只有被 WATCH 的 key 从 WATCH 起始时至今没有发生过变更，EXEC 才会被执行。如果 WATCH 的 key 在 WATCH 命令到 EXEC 命令之间发生过变化，则 EXEC 命令会返回失败。

### **Scripting**

通过 EVAL 与 EVALSHA 命令，可以让 Redis 执行 LUA 脚本。这就类似于 RDBMS 的存储过程一样，可以把客户端与 Redis 之间密集的读/写交互放在服务端进行，避免过多的数据交互，提升性能。Scripting 功能是作为事务功能的替代者诞生的，事务提供的所有能力 Scripting 都可以做到。Redis 官方推荐使用 LUA Script 来代替事务，前者的效率和便利性都超过了事务。关于 Scripting 的具体使用，本文不做详细介绍，请参考官方文档 <https://redis.io/commands/eval>

## **Redis 性能调优**

尽管 Redis 是一个非常快速的内存数据存储媒介，也并不代表 Redis 不会产生性能问题。
前文中提到过，Redis 采用单线程模型，所有的命令都是由一个线程串行执行的，所以当某个命令执行耗时较长时，会拖慢其后的所有命令，这使得 Redis 对每个任务的执行效率更加敏感。针对 Redis 的性能优化，主要从下面几个层面入手：

- 最初的也是最重要的，确保没有让 Redis 执行耗时长的命令
- 使用 pipelining 将连续执行的命令组合执行
- 操作系统的 Transparent huge pages 功能必须关闭：
-

<!---->

    echo never > /sys/kernel/mm/transparent_hugepage/enabled

- 如果在虚拟机中运行 Redis，可能天然就有虚拟机环境带来的固有延迟。可以通过./redis-cli --intrinsic-latency 100 命令查看固有延迟。同时如果对 Redis 的性能有较高要求的话，应尽可能在物理机上直接部署 Redis。
- 检查数据持久化策略
- 考虑引入读写分离机制

### **长耗时命令**

Redis 绝大多数读写命令的时间复杂度都在 O(1)到 O(N)之间，在文本和官方文档中均对每个命令的时间复杂度有说明。通常来说，O(1)的命令是安全的，O(N)命令在使用时需要注意，如果 N 的数量级不可预知，则应避免使用。例如对一个 field 数未知的 Hash 数据执行 HGETALL/HKEYS/HVALS 命令，通常来说这些命令执行的很快，但如果这个 Hash 中的 field 数量极多，耗时就会成倍增长。
又如使用 SUNION 对两个 Set 执行 Union 操作，或使用 SORT 对 List/Set 执行排序操作等时，都应该严加注意。避免在使用这些 O(N)命令时发生问题主要有几个办法：

- 不要把 List 当做列表使用，仅当做队列来使用
- 通过机制严格控制 Hash、Set、Sorted Set 的大小
- 可能的话，将排序、并集、交集等操作放在客户端执行
- 绝对禁止使用 KEYS 命令
- 避免一次性遍历集合类型的所有成员，而应使用 SCAN 类的命令进行分批的，游标式的遍历

Redis 提供了 SCAN 命令，可以对 Redis 中存储的所有 key 进行游标式的遍历，避免使用 KEYS 命令带来的性能问题。同时还有 SSCAN/HSCAN/ZSCAN 等命令，分别用于对 Set/Hash/Sorted Set 中的元素进行游标式遍历。SCAN 类命令的使用请参考官方文档：<https://redis.io/commands/scanRedis>提供了 Slow Log 功能，可以自动记录耗时较长的命令。相关的配置参数有两个：

-

-

<!---->

    slowlog-log-slower-than xxxms  #执行时间慢于xxx毫秒的命令计入Slow Logslowlog-max-len xxx  #Slow Log的长度，即最大纪录多少条Slow Log

`<br />`使用**SLOWLOG GET \[number]**命令，可以输出最近进入 Slow Log 的 number 条命令。
使用**SLOWLOG RESET**命令，可以重置 Slow Log

###

### **网络引发的延迟**

- 尽可能使用长连接或连接池，避免频繁创建销毁连接
- 客户端进行的批量数据操作，应使用 Pipeline 特性在一次交互中完成。具体请参照本文的 Pipelining 章节

###

### **数据持久化引发的延迟**

Redis 的数据持久化工作本身就会带来延迟，需要根据数据的安全级别和性能要求制定合理的持久化策略：

- AOF + fsync always 的设置虽然能够绝对确保数据安全，但每个操作都会触发一次 fsync，会对 Redis 的性能有比较明显的影响
- AOF + fsync every second 是比较好的折中方案，每秒 fsync 一次
- AOF + fsync never 会提供 AOF 持久化方案下的最优性能
- 使用 RDB 持久化通常会提供比使用 AOF 更高的性能，但需要注意 RDB 的策略配置
- 每一次 RDB 快照和 AOF Rewrite 都需要 Redis 主进程进行 fork 操作。fork 操作本身可能会产生较高的耗时，与 CPU 和 Redis 占用的内存大小有关。根据具体的情况合理配置 RDB 快照和 AOF Rewrite 时机，避免过于频繁的 fork 带来的延迟
  Redis 在 fork 子进程时需要将内存分页表拷贝至子进程，以占用了 24GB 内存的 Redis 实例为例，共需要拷贝 24GB / 4kB \* 8 = 48MB 的数据。在使用单 Xeon 2.27Ghz 的物理机上，这一 fork 操作耗时 216ms。可以通过> **INFO**命令返回的 latest_fork_usec 字段查看上一次 fork 操作的耗时（微秒）

###

### **Swap 引发的延迟**

当 Linux 将 Redis 所用的内存分页移至 swap 空间时，将会阻塞 Redis 进程，导致 Redis 出现不正常的延迟。Swap 通常在物理内存不足或一些进程在进行大量 I/O 操作时发生，应尽可能避免上述两种情况的出现。/proc/<pid>/smaps 文件中会保存进程的 swap 记录，通过查看这个文件，能够判断 Redis 的延迟是否由 Swap 产生。如果这个文件中记录了较大的 Swap size，则说明延迟很有可能是 Swap 造成的。

###

### **数据淘汰引发的延迟**

当同一秒内有大量 key 过期时，也会引发 Redis 的延迟。在使用时应尽量将 key 的失效时间错开。

###

### **引入读写分离机制**

Redis 的主从复制能力可以实现一主多从的多节点架构，在这一架构下，主节点接收所有写请求，并将数据同步给多个从节点。
在这一基础上，我们可以让从节点提供对实时性要求不高的读请求服务，以减小主节点的压力。
尤其是针对一些使用了长耗时命令的统计类任务，完全可以指定在一个或多个从节点上执行，避免这些长耗时命令影响其他请求的响应。关于读写分离的具体说明，请参见后续章节

##

## **主从复制与集群分片**

### **主从复制**

Redis 支持一主多从的主从复制架构。一个 Master 实例负责处理所有的写请求，Master 将写操作同步至所有 Slave。
借助 Redis 的主从复制，可以实现读写分离和高可用：

- 实时性要求不是特别高的读请求，可以在 Slave 上完成，提升效率。特别是一些周期性执行的统计任务，这些任务可能需要执行一些长耗时的 Redis 命令，可以专门规划出 1 个或几个 Slave 用于服务这些统计任务
- 借助 Redis Sentinel 可以实现高可用，当 Master crash 后，Redis Sentinel 能够自动将一个 Slave 晋升为 Master，继续提供服务

启用主从复制非常简单，只需要配置多个 Redis 实例，在作为 Slave 的 Redis 实例中配置：

-

<!---->

    slaveof 192.168.1.1 6379  #指定Master的IP和端口

当 Slave 启动后，会从 Master 进行一次冷启动数据同步，由 Master 触发 BGSAVE 生成 RDB 文件推送给 Slave 进行导入，导入完成后 Master 再将增量数据通过 Redis Protocol 同步给 Slave。之后主从之间的数据便一直以 Redis Protocol 进行同步

#### **使用 Sentinel 做自动 failover**

Redis 的主从复制功能本身只是做数据同步，并不提供监控和自动 failover 能力，要通过主从复制功能来实现 Redis 的高可用，还需要引入一个组件：Redis SentinelRedis Sentinel 是 Redis 官方开发的监控组件，可以监控 Redis 实例的状态，通过 Master 节点自动发现 Slave 节点，并在监测到 Master 节点失效时选举出一个新的 Master，并向所有 Redis 实例推送新的主从配置。Redis Sentinel 需要至少部署 3 个实例才能形成选举关系。关键配置：

-

-

-

-

<!---->

    sentinel monitor mymaster 127.0.0.1 6379 2  #Master实例的IP、端口，以及选举需要的赞成票数sentinel down-after-milliseconds mymaster 60000  #多长时间没有响应视为Master失效sentinel failover-timeout mymaster 180000  #两次failover尝试间的间隔时长sentinel parallel-syncs mymaster 1  #如果有多个Slave，可以通过此配置指定同时从新Master进行数据同步的Slave数，避免所有Slave同时进行数据同步导致查询服务也不可用

另外需要注意的是，Redis Sentinel 实现的自动 failover 不是在同一个 IP 和端口上完成的，也就是说自动 failover 产生的新 Master 提供服务的 IP 和端口与之前的 Master 是不一样的，所以要实现 HA，还要求客户端必须支持 Sentinel，能够与 Sentinel 交互获得新 Master 的信息才行。

### **集群分片**

为何要做集群分片：

- Redis 中存储的数据量大，一台主机的物理内存已经无法容纳
- Redis 的写请求并发量大，一个 Redis 实例以无法承载

当上述两个问题出现时，就必须要对 Redis 进行分片了。
Redis 的分片方案有很多种，例如很多 Redis 的客户端都自行实现了分片功能，也有向 Twemproxy 这样的以代理方式实现的 Redis 分片方案。然而首选的方案还应该是 Redis 官方在 3.0 版本中推出的 Redis Cluster 分片方案。本文不会对 Redis Cluster 的具体安装和部署细节进行介绍，重点介绍 Redis Cluster 带来的好处与弊端。

#### **Redis Cluster 的能力**

- 能够自动将数据分散在多个节点上
- 当访问的 key 不在当前分片上时，能够自动将请求转发至正确的分片
- 当集群中部分节点失效时仍能提供服务

其中第三点是基于主从复制来实现的，Redis Cluster 的每个数据分片都采用了主从复制的结构，原理和前文所述的主从复制完全一致，唯一的区别是省去了 Redis Sentinel 这一额外的组件，由 Redis Cluster 负责进行一个分片内部的节点监控和自动 failover。

#### **Redis Cluster 分片原理**

Redis Cluster 中共有 16384 个 hash slot，Redis 会计算每个 key 的 CRC16，将结果与 16384 取模，来决定该 key 存储在哪一个 hash slot 中，同时需要指定 Redis Cluster 中每个数据分片负责的 Slot 数。Slot 的分配在任何时间点都可以进行重新分配。客户端在对 key 进行读写操作时，可以连接 Cluster 中的任意一个分片，如果操作的 key 不在此分片负责的 Slot 范围内，Redis Cluster 会自动将请求重定向到正确的分片上。

#### **hash tags**

在基础的分片原则上，Redis 还支持 hash tags 功能，以 hash tags 要求的格式明明的 key，将会确保进入同一个 Slot 中。例如：{uiv}user:1000 和{uiv}user:1001 拥有同样的 hash tag {uiv}，会保存在同一个 Slot 中。使用 Redis Cluster 时，pipelining、事务和 LUA Script 功能涉及的 key 必须在同一个数据分片上，否则将会返回错误。如要在 Redis Cluster 中使用上述功能，就必须通过 hash tags 来确保一个 pipeline 或一个事务中操作的所有 key 都位于同一个 Slot 中。
有一些客户端（如 Redisson）实现了集群化的 pipelining 操作，可以自动将一个 pipeline 里的命令按 key 所在的分片进行分组，分别发到不同的分片上执行。但是 Redis 不支持跨分片的事务，事务和 LUA Script 还是必须遵循所有 key 在一个分片上的规则要求。

###

### **主从复制 vs 集群分片**

在设计软件架构时，要如何在主从复制和集群分片两种部署方案中取舍呢？从各个方面看，Redis Cluster 都是优于主从复制的方案

- Redis Cluster 能够解决单节点上数据量过大的问题
- Redis Cluster 能够解决单节点访问压力过大的问题
- Redis Cluster 包含了主从复制的能力

那是不是代表 Redis Cluster 永远是优于主从复制的选择呢？并不是。软件架构永远不是越复杂越好，复杂的架构在带来显著好处的同时，一定也会带来相应的弊端。采用 Redis Cluster 的弊端包括：

- 维护难度增加。在使用 Redis Cluster 时，需要维护的 Redis 实例数倍增，需要监控的主机数量也相应增加，数据备份/持久化的复杂度也会增加。同时在进行分片的增减操作时，还需要进行 reshard 操作，远比主从模式下增加一个 Slave 的复杂度要高。
- 客户端资源消耗增加。当客户端使用连接池时，需要为每一个数据分片维护一个连接池，客户端同时需要保持的连接数成倍增多，加大了客户端本身和操作系统资源的消耗。
- 性能优化难度增加。你可能需要在多个分片上查看 Slow Log 和 Swap 日志才能定位性能问题。
- 事务和 LUA Script 的使用成本增加。在 Redis Cluster 中使用事务和 LUA Script 特性有严格的限制条件，事务和 Script 中操作的 key 必须位于同一个分片上，这就使得在开发时必须对相应场景下涉及的 key 进行额外的规划和规范要求。如果应用的场景中大量涉及事务和 Script 的使用，如何在保证这两个功能的正常运作前提下把数据平均分到多个数据分片中就会成为难点。

所以说，在主从复制和集群分片两个方案中做出选择时，应该从应用软件的功能特性、数据和访问量级、未来发展规划等方面综合考虑，只在**确实有必要**引入数据分片时再使用 Redis Cluster。
下面是一些建议：

1. 需要在 Redis 中存储的数据有多大？未来 2 年内可能发展为多大？这些数据是否都需要长期保存？是否可以使用 LRU 算法进行非热点数据的淘汰？综合考虑前面几个因素，评估出 Redis 需要使用的物理内存。
2. 用于部署 Redis 的主机物理内存有多大？有多少可以分配给 Redis 使用？对比(1)中的内存需求评估，是否足够用？
3. Redis 面临的并发写压力会有多大？在不使用 pipelining 时，Redis 的写性能可以超过 10 万次/秒（更多的 benchmark 可以参考 <https://redis.io/topics/benchmarks> ）
4. 在使用 Redis 时，是否会使用到 pipelining 和事务功能？使用的场景多不多？

综合上面几点考虑，如果单台主机的可用物理内存完全足以支撑对 Redis 的容量需求，且 Redis 面临的并发写压力距离 Benchmark 值还尚有距离，建议采用主从复制的架构，可以省去很多不必要的麻烦。同时，如果应用中大量使用 pipelining 和事务，也建议尽可能选择主从复制架构，可以减少设计和开发时的复杂度。

## **Redis Java 客户端的选择**

Redis 的 Java 客户端很多，官方推荐的有三种：Jedis、Redisson 和 lettuce。在这里对 Jedis 和 Redisson 进行对比介绍**Jedis：**

- 轻量，简洁，便于集成和改造
- 支持连接池
- 支持 pipelining、事务、LUA Scripting、Redis Sentinel、Redis Cluster
- 不支持读写分离，需要自己实现
- 文档差（真的很差，几乎没有……）

**Redisson：**

- 基于 Netty 实现，采用非阻塞 IO，性能高
- 支持异步请求
- 支持连接池
- 支持 pipelining、LUA Scripting、Redis Sentinel、Redis Cluster
- 不支持事务，官方建议以 LUA Scripting 代替事务
- 支持在 Redis Cluster 架构下使用 pipelining
- 支持读写分离，支持读负载均衡，在主从复制和 Redis Cluster 架构下都可以使用
- 内建 Tomcat Session Manager，为 Tomcat 6/7/8 提供了会话共享功能
- 可以与 Spring Session 集成，实现基于 Redis 的会话共享
- 文档较丰富，有中文文档

对于 Jedis 和 Redisson 的选择，同样应遵循前述的原理，尽管 Jedis 比起 Redisson 有各种各样的不足，但也应该在需要使用 Redisson 的高级特性时再选用 Redisson，避免造成不必要的程序复杂度提升。
**Jedis：
**github：<https://github.com/xetorthio/jedis>
文档：<https://github.com/xetorthio/jedis/wiki>
**Redisson：
**github：<https://github.com/redisson/redisson>
文档：<https://github.com/redisson/redisson/wiki>
