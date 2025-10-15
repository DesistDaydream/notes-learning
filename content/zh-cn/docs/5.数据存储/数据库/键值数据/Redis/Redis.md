---
title: Redis
linkTitle: Redis
weight: 1
---


# 概述

> 参考：
>
> - [GitHub 项目，redis/redis](https://github.com/redis/redis)
> - [GitHub 项目，valkey-io/valkey](https://github.com/valkey-io/valkey) # [Linux Foundation](/docs/Standard/Foundation/Linux%20Foundation.md) 基于 7.2.4 版本分叉的项目，保持原有 License
> - [官网](https://redis.io/)

> [!Warning]
> [3 月 30 日 Redis 发布博客改变 License](https://redis.com/blog/redis-adopts-dual-source-available-licensing/)，后续 Linux 基金会基于 7.2.4 版本分叉，保持原有 License。开源版本改名称为 Redis OSS（open source）

Redis 是一个 网络化的、内存中的、具有持久化的 [键值数据](docs/5.数据存储/数据库/键值数据/键值数据.md)存储。(是否持久化根据配置决定)

Redis 是一个内存数据库, 所有数据默认都存在于内存当中,可以配置“定时以追加或者快照”的方式储存到硬盘中. 由于 redis 是一个内存数据库, 所以读取写入的速度是非常快的, 所以经常被用来做数据, 页面等的缓存。

## Redis 的组件

- redis-server # 服务端
- redis-cli # [Redis CLI](/docs/5.数据存储/数据库/键值数据/Redis/Redis%20CLI/Redis%20CLI.md)，一个命令行客户端
- redis-benchmark # 压测工具
- redis-check-dump && redis-check-aof # 检测工具

## Redis 的数据类型

[Redis 数据类型](/docs/5.数据存储/数据库/键值数据/Redis/Redis%20数据类型.md)

# Redis 部署

https://redis.io/docs/latest/operate/oss_and_stack/install/

## Docker

https://redis.io/docs/latest/operate/oss_and_stack/install/install-stack/docker/

```bash
docker run -d --name redis \
  --network=host \
  redis:5.0.10-alpine
```

## 二进制

Redis 官方没有提供二进制包，只能通过源码构建。

> 但是我将 redis-server 文件从 RedHat 系系统中拷贝到 Debian 系系统中后，直接执行是可以使用的

在 [GitHub 项目，redis/redis -utils/systemd-redis_server.service](utils/systemd-redis_server.service) 处可以找到 Redis 给 [Systemd](docs/1.操作系统/Systemd/Systemd.md) 使用的 Unit 文件。

```bash
[Unit]
Description=Redis data structure server
Documentation=https://redis.io/documentation
#Before=your_application.service another_example_application.service
#AssertPathExists=/var/lib/redis
Wants=network-online.target
After=network-online.target

[Service]
ExecStart=/usr/local/bin/redis-server --supervised systemd --daemonize no
## Alternatively, have redis-server load a configuration file:
#ExecStart=/usr/local/bin/redis-server /path/to/your/redis.conf
LimitNOFILE=10032
NoNewPrivileges=yes
#OOMScoreAdjust=-900
#PrivateTmp=yes
Type=notify
TimeoutStartSec=infinity
TimeoutStopSec=infinity
UMask=0077
#User=redis
#Group=redis
#WorkingDirectory=/var/lib/redis

[Install]
WantedBy=multi-user.target
```

> [!Attention]
> 其中要注意的是，在配置文件或命令行参数中，daemonize 配置需要设置为 no。因为 Redis 以后台模式启动的话，Systemd 会认为 Redis 没有启动成功，进而频繁重启。

# Redis 关联文件与配置

**/etc/redis.conf** # Redis 主程序的配置文件

**/var/lib/redis/** # 默认的数据存储目录

# Redis 数据持久化的方式

Redis 的数据是保存在内存中的，如果设备宕机，则数据丢失，所以 Redis 提供两种可以将数据从内存中写入硬盘的方式，默认为 RDB，AOF 默认不开启

## Redis Data Base

**Redis Data Base(简称 RDB)** 相关配置在配置文件的 SNAPSHOT 配置环境中。

在默认情况下， Redis 将数据库快照保存在名字为 dump.rdb 的二进制文件中。可以对 Redis 进行设置， 让 Redis 在“ N 秒内数据集至少有 M 个 key 改动”这一条件被满足时， 自动保存一次数据集。也可以通过调用 SAVE 或者 BGSAVE ， 手动让 Redis 进行数据集保存操作。

这种持久化的方式被称为快照（snapshotting），会将数据保存在一个指定的文件中。

RDB 工作原理：

- 触发 RDB 后，redis 会调用 fork(),产生一个与主程序同名的子进程。
- 该子进程会将现有内存中的数据写入到一个临时的 RDB 文件中。(文件名一般为：temp-XXX.rdb)
  - Redis 默认会使用 LZF 算法对数据进行压缩。该算法会消耗大量 CPU，可以在配置中关闭压缩功能，但是数据量写入到磁盘后，会占用大量磁盘空间。
- 当子进程完成对临时 RDB 文件的写入时，Redis 用这个临时的 RDB 文件替换原来的 RDB 文件，并删除旧的 RDB 文件。

这种工作方式使得 Redis 可以从写时复制（copy-on-write）机制中获益。

RDB 优点：

- RDB 文件是一个很简洁的单文件，它保存了某个时间点的 Redis 数据，很适合用于做备份。你可以设定一个时间点对 RDB 文件进行归档，这样就能在需要的时候很轻易的把数据恢复到不同的版本。
- 基于上面所描述的特性，RDB 很适合用于灾备。单文件很方便就能传输到远程的服务器上。
- RDB 的性能很好，需要进行持久化时，主进程会 fork 一个子进程出来，然后把持久化的工作交给子进程，自己不会有相关的 I/O 操作。
- 比起 AOF，在数据量比较大的情况下，RDB 的启动速度更快。

RDB 缺点：

- RDB 容易造成数据的丢失。假设每 5 分钟保存一次快照，如果 Redis 因为某些原因不能正常工作，那么从上次产生快照到 Redis 出现问题这段时间的数据就会丢失了。
- RDB 使用 fork()产生子进程进行数据的持久化，如果数据比较大的话可能就会花费点时间，造成 Redis 停止服务几毫秒。如果数据量很大且 CPU 性能不是很好的时候，停止服务的时间甚至会到 1 秒。

## Append Only File

**Append Only File(简称 AOF)** 配置文件中配置环境 APPEND ONLY MOD E 可以进项相关配置。

默认 redis 使用的是 rdb 方式持久化，这种方式在许多应用中已经足够用了。但是 redis 如果中途宕机，会导致可能有几分钟的数据丢失。

Append Only File 是另一种持久化方式，可以提供更好的持久化特性。Redis 会把每次写入的数据在接收后都写入 appendonly.aof 文件，每次启动时 Redis 都会先把这个文件的数据读入内存里，先忽略 RDB 文件。

AOF 工作方式：

每当 Redis 执行一个改变数据集的命令时（比如 SET）， 这个命令就会被追加到 AOF 文件的末尾。这样的话， 当 Redis 重新启时， 程序就可以通过重新执行 AOF 文件中的命令来达到重建数据集的目的。AOF 重写和 RDB 创建快照一样，都巧妙地利用了写时复制机制:

- Redis 执行 fork() ，现在同时拥有父进程和子进程。
- 子进程开始将新 AOF 文件的内容写入到临时文件。
- 对于所有新执行的写入命令，父进程一边将它们累积到一个内存缓存中，一边将这些改动追加到现有 AOF 文件的末尾,这样样即使在重写的中途发生停机，现有的 AOF 文件也还是安全的。
- 当子进程完成重写工作时，它给父进程发送一个信号，父进程在接收到信号之后，将内存缓存中的所有数据追加到新 AOF 文件的末尾。
- 搞定！现在 Redis 原子地用新文件替换旧文件，之后所有命令都会直接追加到新 AOF 文件的末尾。

AOF 优点：

- 使用 AOF 会让你的 Redis 更加耐久: 你可以使用不同的 fsync 策略：无 fsync,每秒 fsync,每次写的时候 fsync.使用默认的每秒 fsync 策略,Redis 的性能依然很好(fsync 是由后台线程进行处理的,主线程会尽力处理客户端请求),一旦出现故障，你最多丢失 1 秒的数据.
- AOF 文件是一个只进行追加的日志文件,所以不需要写入 seek,即使由于某些原因(磁盘空间已满，写的过程中宕机等等)未执行完整的写入命令,你也也可使用 redis-check-aof 工具修复这些问题.
- Redis 可以在 AOF 文件体积变得过大时，自动地在后台对 AOF 进行重写： 重写后的新 AOF 文件包含了恢复当前数据集所需的最小命令集合。 整个重写操作是绝对安全的，因为 Redis 在创建新 AOF 文件的过程中，会继续将命令追加到现有的 AOF 文件里面，即使重过程中发生停机，现有的 AOF 文件也不会丢失。 而一旦新 AOF 文件创建完毕，Redis 就会从旧 AOF 文件切换到新 AOF 文件，并开始对新 AOF 文件进行追加操作。
- AOF 文件有序地保存了对数据库执行的所有写入操作， 这些写入操作以 Redis 协议的格式保存， 因此 AOF 文件的内容非常容易被人读懂， 对文件进行分析（parse）也很轻松。 导出（export） AOF 文件也非常简单： 举个例子， 如果你不小心执行了 FLUSHALL 命令， 但只要 AOF 文件未被重写， 那么只要停止服务器， 移除 AOF 文件末尾的 FLUSHALL 命令， 并重启 Redis ， 就可以将数据集恢复到 FLUSHALL 执行之前的状态。

AOF 缺点：

- 对于相同的数据集来说，AOF 文件的体积通常要大于 RDB 文件的体积。
- 根据所使用的 fsync 策略，AOF 的速度可能会慢于 RDB 。 在一般情况下， 每秒 fsync 的性能依然非常高， 而关闭 fsync 可以让 AOF 的速度和 RDB 一样快， 即使在高负荷之下也是如此。 不过在处理巨大的写入载入时，RDB 可以提供更有保证的最大延迟时间（latency）。

下图是 redis 启动后读取本地存储文件的过程

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/gghosd/1616134974086-4020ec57-f508-4a12-b30e-aba72d1730e4.jpeg)
