---
title: Redis 管理
linkTitle: Redis 管理
weight: 1
---

# 概述

> 参考：
>
> - [官方文档，Redis 管理](https://redis.io/topics/admin)
>   - https://redis.io/docs/latest/operate/oss_and_stack/management/admin/

在生产中配置和管理 Redis 的建议。

# 概述

> 参考：
>
> - [官方中文文档](http://www.redis.cn/topics/admin.html)

官方优化建议

1. 我们建议使用 linux 部署 Redis。Redis 也在 osx，FreeBSD，OpenBSD 上经过测试，但 Linux 经过所有主要的压力测试，并且最多产品部署。 确保设置 Linux 内核 overcommit memory setting 为 1。向/etc/sysctl.conf 添加 vm.overcommit_memory = 1 然后重启，或者运行命令 sysctl vm.overcommit_memory=1 以便立即生效。
2. 确保禁用 Linux 内核特性 transparent huge pages，它对内存使用和延迟有非常大的负面影响。通过命令 echo never > /sys/kernel/mm/transparent_hugepage/enabled 来完成。
3. 确保你的系统设置了一些 swap（我们建议和内存一样大）。如果 linux 没有 swap 并且你的 redis 实例突然消耗了太多内存，或者 Redis 由于内存溢出会宕掉，或者 Linux 内核 OOM Killer 会杀掉 Redis 进程。
4. 设置一个明确的 maxmemory 参数来限制你的实例，以便确保实例会报告错误而不是当接近系统内存限制时失败
5. 如果你对一个写频繁的应用使用 redis，当向磁盘保存 RDB 文件或者改写 AOF 日志时，redis 可能会用正常使用内存 2 倍的内存。额外使用的内存和保存期间写修改的内存页数量成比例，因此经常和这期间改动的键的数量成比例。确保相应的设置内存的大小。
6. 当在 daemontools 下运行时，使用 daemonize no
7. 即使你禁用持久化，如果你使用复制，redis 会执行 rdb 保存，除非你使用新的无磁盘复制特性，这个特性目前还是实验性的。
8. 如果你使用复制，确保要么你的 master 激活了持久化，要么它不会在当掉后自动重启。slave 是 master 的完整备份，因此如果 master 通过一个空数据集重启，slave 也会被清掉。

Redis 延迟问题疑难解答

<http://www.redis.cn/topics/latency.html>

# 管理工具

https://github.com/qishibo/AnotherRedisDesktopManager # GUI 客户端

# 故障处理

## redis 进程占用 CPU 很高-达到 100

问题说明：

    监控发现，redis进程占用CPU很高-达到100%。并且会有2个redis进程。如下图：

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/sq1d5g/1616134552751-f644a68e-d162-4bad-8a09-a5909044b5b2.jpeg)

    分析了一下，因为redis在持久化保存的时候，会fork出一个进程来。仔细观察进程号PID，会发现，占用CPU很高的那个进程，不是redis的主进程。而是fork出来的那个。这个fork出来的进程，由于任务就是持久化，所以它的工作是：把内存中的数据（此时内存数据，约2.18G），拷贝出来到新的进程中，然后进行压缩，保存到硬盘上（硬盘数据大约是700M）。在压缩的过程中，是要用CPU的。

解决办法：

    个人觉得，如果主进程CPU占用不高，并且没有服务延迟，那不管用这个fork进程CPU跑的有多高。如果十分在意，那可以考虑，更改redis配置，不压缩数据保存。
