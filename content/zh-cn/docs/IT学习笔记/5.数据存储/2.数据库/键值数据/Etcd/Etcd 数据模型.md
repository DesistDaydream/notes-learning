---
title: Etcd 数据模型
---

# 概述

> 参考：
> - [官方文档](https://etcd.io/docs/latest/learning/data_model/)

# WAL

数据库通常使用 [预写日志记录](https://en.wikipedia.org/wiki/Write-ahead_logging)； etcd 也使用它。有关预写日志记录的详细信息不在本文讨论范围之内，但是出于此目的，我们需要知道的是-每个 etcd 集群成员都在持久性存储上保留一个预写日志（WAL）。 etcd 在将某些操作（例如更新）写入键值存储之前，将其写入 WAL。如果成员崩溃并在快照之间重新启动，则可以通过查看 WAL 的内容在本地恢复自上次快照以来完成的事务。

因此，每当客户将密钥添加到键值存储或更新现有密钥的值时，客户端都会 etcd 在 WAL 上附加一个记录操作的条目，WAL 是持久性存储上的普通文件。在继续进行之前， etcd 必须 100％确信 WAL 条目已被实际保留。要在 Linux 上实现此目的，仅使用 write 系统调用是不够的， 因为实际写入物理存储的时间可能会延迟。例如，Linux 可能会将写入的 WAL 条目在内核内存高速缓存（例如页面高速缓存）中保留一段时间。为了确保将数据写入持久性存储中，您必须在“”之后调用 fdatasync 系统调用，write 这正是该 etcd 操作（如以下 strace 所示） 输出，其中 8 是 WAL 文件的文件描述符）：

    21:23:09.894875 lseek(8, 0, SEEK_CUR)   = 12808 <0.000012>
    21:23:09.894911 write(8, ".\0\0\0\0\0\0\202\10\2\20\361\223\255\266\6\32$\10\0\20\10\30\26\"\34\"\r\n\3fo"..., 2296) = 2296 <0.000130>
    21:23:09.895041 fdatasync(8)            = 0 <0.008314>

不幸的是，写入持久性存储需要时间。如果 fdatasync 花费太长时间，则 etcd 系统性能会降低。 [etcd 文档建议](https://github.com/etcd-io/etcd/blob/master/Documentation/faq.md#what-does-the-etcd-warning-failed-to-send-out-heartbeat-on-time-mean) 为了使存储足够快，fdatasync 写入 WAL 文件时调用的第 99 个百分点 **必须小于 10ms**。还有其他与存储相关的指标，但这是本文的重点。

# bbolt

<https://github.com/DesistDaydream/go-library/tree/master/bbolt> 代码练习
