---
title: compaction failed
---

# 概述

compaction failed 是一个 Prometheus 在压缩数据时产生的错误，导致该问题的因素多种多样，最常见的就是使用 NFS 作为 Prometehus 时序数据库的后端存储。

[官方文档](https://prometheus.io/docs/prometheus/latest/storage/)中曾明确表明不支持 NFS 文件系统
![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/mhabk3/1633918055761-a5d5266e-c5ce-455c-92c1-0b219b2a2c60.png)
该问题的表现形式通常为文件丢失，比如，某个 Block 中的 meta.json 文件丢失

```bash
msg="compaction failed" err="plan compaction: open /prometheus/01FHHPS3NR7M2E8MAV37S61ME6/meta.json: no such file or directory"

msg="Failed to read meta.json for a block during reloadBlocks. Skipping" dir=/prometheus/01FHHPS3NR7M2E8MAV37S61ME6 err="open /prometheus/01FHHPS3NR7M2E8MAV37S61ME6/meta.json: no such file or directory"
```

经过日志筛选，该问题起源于一次 Deleting obsolete block 操作之后的 compact blocks，也就是删除过期块后压缩块。 失败操作源于：

```bash
msg="compaction failed" err="delete compacted block after failed db reloadBlocks:01FHHPS3NR7M2E8MAV37S61ME6: unlinkat /prometheus/01FHHPS3NR7M2E8MAV37S61ME6/chunks: directory not empty"
```

![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/mhabk3/1633919258306-20285a8b-186e-4177-a043-5f06f54f7f2a.png)

这些报错日志信息，可以在 [./prometheus/tsdb/db.go](https://github.com/prometheus/prometheus/blob/release-2.28/tsdb/db.go) 代码中找到

## 解决方式

可以直接删除 01FHHPS3NR7M2E8MAV37S61ME6 块，也就是直接删除这个目录，并重启 Prometheus
