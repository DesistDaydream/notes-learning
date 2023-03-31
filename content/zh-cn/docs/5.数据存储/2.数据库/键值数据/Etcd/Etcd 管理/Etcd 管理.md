---
title: Etcd 管理
---

# 概述

> 参考：
> - GitHub
>     - https://github.com/etcd-io/website/blob/main/content/en/docs/v3.5/op-guide/maintenance.md
> - [官方文档,-运维指南-维护](https://etcd.io/docs/latest/op-guide/maintenance/)

Etcd 集群需要定期 **Maintenacne(维护)** 才能保持可靠性。根据 etcd 应用程序的需求，通常可以自动执行该维护，而无需停机或性能显着降低。

所有 etcd 维护都管理 etcd 键空间消耗的存储资源。存储空间配额可以防止无法充分控制键空间大小；如果 etcd 成员的空间不足，则配额将触发群集范围的警报，这将使系统进入有限操作维护模式。为了避免空间不足以写入键空间，必须压缩 etcd 键空间历史记录。可以通过对 etcd 成员进行碎片整理来回收存储空间本身。最后，etcd 成员状态的定期快照备份使恢复由于操作错误引起的意外逻辑数据丢失或损坏成为可能。

# Raft Log Retention(Raft 日志保留)

https://etcd.io/docs/v3.5/op-guide/maintenance/#raft-log-retention

# Auto Compaction(自动压缩)

https://etcd.io/docs/v3.5/op-guide/maintenance/#auto-compaction

# Defragmentation(碎片整理)

https://etcd.io/docs/v3.5/op-guide/maintenance/#defragmentation

在压缩 keyspace 之后，Etcd 数据库可能会出现内部 **Fragmentation(碎片)**。任何内部碎片都是后端可以免费使用但仍会占用存储空间的空间。通过在后端数据库中留下空白来在内部压缩旧修订版碎片 etcd。碎片空间可供 etcd 使用，但主机文件系统不可用。换句话说，删除应用程序数据不会回收磁盘空间。

碎片整理过程将此存储空间释放回文件系统。碎片整理是针对每个成员进行的，因此可以避免集群范围内的延迟峰值。

在 kube-prometheus-stack 项目中，会自带碎片所占空间的告警，告警名称为 **etcdDatabaseHighFragmentationRatio**，当出现该告警时，即可执行碎片整理操作。

具体用法详见 [etcdctl](/docs/5.数据存储/2.数据库/键值数据/Etcd/Etcd%20命令行工具/etcdctl.md#defrag)

# Etcd Space Quota(Etcd 空间配额)

> 参考：
> 
> - [官方文档，运维指南-维护-空间配额](https://etcd.io/docs/v3.5/op-guide/maintenance/#space-quota)

etcd 通过 **Space Quota(空间配额)** 确保集群以可靠的方式运行，空间配额指的是 etcd 可以储存的数据量上限。没有空间配额，如果密钥空间过大，etcd 可能会遭受性能不佳的影响，或者它可能只是用尽了存储空间，从而导致了不可预测的集群行为。

默认情况下，etcd 的空间配额适合大多数应用程序的使用情况。不过，可以通过 quota-backend-bytes 命令行参数修改配额的值

注意如果 etcd 中的数据超**过了配额的值**，则**无法再写入新数据**。并且 etcd 会在集群中发出一个 alarm(警报)，该警报会告诉各节点，并且集群将会变为 maintenance mode(维护模式)，处于维护模式的集群仅接受 key 的读取和删除操作。并且如果想让集群恢复正常运行，需要进行如下操作

### (测试用)写入数据触发告警

```bash
# 使用一个循环填满 keyspace(键空间)，空间爆满后，触发了告警
[root@lichenhao ~]# while [ 1 ]; do dd if=/dev/urandom bs=1024 count=1024  | etcdctl put key  || break; done
......
1048576 bytes (1.0 MB, 1.0 MiB) copied, 0.0132167 s, 79.3 MB/s
{"level":"warn","ts":"2020-11-19T22:20:16.018+0800","caller":"clientv3/retry_interceptor.go:62","msg":"retrying of unary invoker failed","target":"endpoint://client-e3fb4b20-987e-4de7-b6d4-bb41ceb2ff59/127.0.0.1:2379","attempt":0,"error":"rpc error: code = ResourceExhausted desc = etcdserver: mvcc: database space exceeded"}
Error: etcdserver: mvcc: database space exceeded

# 查看 etcd 状态，也可以看到告警，告警内容 alarm: NOSPACE
[root@lichenhao ~]# etcdctl endpoint status -wtable
+----------------+------------------+---------+---------+-----------+------------+-----------+------------+--------------------+--------------------------------+
|    ENDPOINT    |        ID        | VERSION | DB SIZE | IS LEADER | IS LEARNER | RAFT TERM | RAFT INDEX | RAFT APPLIED INDEX |             ERRORS             |
+----------------+------------------+---------+---------+-----------+------------+-----------+------------+--------------------+--------------------------------+
| 127.0.0.1:2379 | 656f8f6ebad83496 |  3.4.13 |  2.2 GB |      true |      false |         2 |      28412 |              28412 |   memberID:7309218425989510294 |
|                |                  |         |         |           |            |           |            |                    |                 alarm:NOSPACE  |
+----------------+------------------+---------+---------+-----------+------------+-----------+------------+--------------------+--------------------------------+

# 此时已经无法向 etcd 中写入数据
[root@lichenhao ~]# etcdctl put /hello world
{"level":"warn","ts":"2020-11-19T22:35:16.870+0800","caller":"clientv3/retry_interceptor.go:62","msg":"retrying of unary invoker failed","target":"endpoint://client-ed96eb23-4bc4-41e9-84bf-1a357fde2b6f/127.0.0.1:2379","attempt":0,"error":"rpc error: code = ResourceExhausted desc = etcdserver: mvcc: database space exceeded"}
Error: etcdserver: mvcc: database space exceeded
```

### 释放足够的空间并对数据库进行碎片整理

```bash
# 获取当前 revision
~]# rev=$(etcdctl endpoint status --write-out="json" | egrep -o '"revision":[0-9]*' | egrep -o '[0-9].*')
# 压缩所有旧的 revisions
~]# etcdctl compaction $rev
compacted revision 28406
# 使用碎片整理，处理多余的空间
~]# etcdctl defrag
Finished defragmenting etcd member[127.0.0.1:2379]
```

### 清除警报

```bash
# 解除告警
~]# etcdctl alarm disarm
memberID:7309218425989510294 alarm:NOSPACE
```

etcd_mvcc_db_total_size_in_use_in_bytes 指标指示历史记录压缩后的实际数据库使用情况 etcd_mvcc_db_total_size_in_bytes 指标则显示数据库大小，包括等待进行碎片整理的可用空间。后者仅在前者接近时才增加，这意味着当这两个指标都接近配额时，需要进行历史压缩以避免触发空间配额告警。
