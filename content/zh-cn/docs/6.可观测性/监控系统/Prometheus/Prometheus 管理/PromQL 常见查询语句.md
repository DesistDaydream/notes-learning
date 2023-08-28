---
title: PromQL 常见查询语句
---

# 概述

> 参考：
> 
> - [GitHub 项目，samber/awesome-prometheus-alerts](https://github.com/samber/awesome-prometheus-alerts)
> - [腾讯云+社区，prometheus 告警指标](https://cloud.tencent.com/developer/article/1758801)
> - [公众号，云原生小白-监控容器 OOMKill 的正确指标](https://mp.weixin.qq.com/s/rPxTBYmwG_7HnZRpRXMFuQ)

# 问题

如何获取范围向量中的第一个和最后一个值。<https://stackoverflow.com/questions/68895729/how-to-get-the-first-and-last-element-of-a-range-vector-in-promql>

- MetricsQS 中有 `first_over_time()` 函数

如何获取范围向量中，指定的值。<https://stackoverflow.com/questions/45213745/prometheus-how-to-calculate-proportion-of-single-value-over-time>，比如 `count_over_time(my_metric[1m] != 0) `获取 1 分钟内所有值中不为 0 的值

- MetricsQL 中有 `count_ne_over_time(my_metric[1h], 0) ` 函数

# SLO/SLI

## 在过去一段时间内内，服务中断了多少时间

这个前提是，服务状态为 0 或者 1。

```bash
avg_over_time(up{}[1h])
```

此时，如果结果为 0.9，则表示服务在 90% 的时间处于运行状态。如果按照 1 小时算，则这 1 小时中，有 6 分钟是停机的。

# 数通设备资源查询语句

查询端口接收的实时带宽。注意:是带宽

    irate(ifHCInOctets{instance="IP.IP.IP.IP",ifAlias="XXXX"}[6m]) * 8

查询端口发送的实时带宽。注意:是带宽

    irate(ifHCOutOctets{instance="IP.IP.IP.IP",ifAlias="XXXX"}[6m]) * 8

# 物理机资源查询语句

## CPU

### CPU 的使用率

显示 cpu 的每个逻辑 core 的使用率

    avg(irate(node_cpu_seconds_total{mode="idle"}[5m])) by(instance,job)
    * 100 < 20

查询物理机 CPU 的使用率，显示总体使用率

    100 - avg (irate(node_cpu_seconds_total{instance="XXXX",mode="idle"}[5m])) by (instance) * 100

### 上下文切换越来越多

    (rate(node_context_switches_total[5m])) / (count (node_cpu_seconds_total{mode="idle"}) without(cpu, mode)) > 1000

- `rate(node_context_switches_total[5m]` # 设备上下文切换在 5 分钟之间的变化量
- `(count (node_cpu_seconds_total{mode="idle"}) without(cpu, mode))` # 获取 instance 的 CPU 总核数
- 两个序列相除，即可获得每个 CPU 核上，5 分钟的上下文切换次数的变化量

## 内存

### 内存使用率

    node_memory_MemAvailable_bytes{}
    /
    node_memory_MemTotal_bytes{}
    * 100 < 10

### OOM

检测主机是否发生了 oom

    increase(node_vmstat_oom_kill[5m]) > 0

## 磁盘

### 磁盘使用率

**使用率过高**

    (node_filesystem_avail_bytes{fstype=~"ext4|xfs"}
    /
    node_filesystem_size_bytes {fstype=~"ext4|xfs"}
    * 100)
    < 20

**磁盘将满**

- 根据磁盘 1 小时的变化速率，预测 4 小时内会不会被写满


    predict_linear(node_filesystem_free_bytes{fstype!~"tmpfs"}[1h], 4 * 3600) < 0

### IO 使用率

    100
    -
    (avg(irate(node_disk_io_time_seconds_total[5m])) by(instance,job))
    * 100 < 20

### 读写速率

读取速率

    sum (irate(node_disk_read_bytes_total[2m])) by (instance)
    / 1024 / 1024 > 200

写入速率

    sum (irate(node_disk_written_bytes_total[2m])) by (instance)
    / 1024 / 1024 > 200

### 读写延迟

读取延迟

    rate(node_disk_read_time_seconds_total[1m])
    /
    rate(node_disk_reads_completed_total[1m])
    > 0.1 and rate(node_disk_reads_completed_total[1m])
    > 0

写入延迟

    rate(node_disk_write_time_seconds_total[1m])
    /
    rate(node_disk_writes_completed_total[1m])
    > 0.1 and rate(node_disk_writes_completed_total[1m])
    > 0

## 网络

### 流量过高

接收和发送的带宽大于 2.5GiB/s 时告警
**接收**

    max(rate(node_network_receive_bytes_total{device!~'tap.*|veth.*|br.*|docker.*|virbr*|lo*'}[2m])) by(instance,job)
    / 1024 / 1024 / 1024 * 8 > 2.5

**发送**

    max(rate(node_network_transmit_bytes_total{device!~'tap.*|veth.*|br.*|docker.*|virbr*|lo*'}[2m])) by(instance,job)
    / 1024 / 1024 / 1024 * 8 > 2.5

### 错误包过多

接收

    increase(node_network_receive_errs_total[5m]) > 0

发送

    increase(node_network_transmit_errs_total[5m]) > 0

### TCP_ESTABLISHED 过高

    node_netstat_Tcp_CurrEstab > 50000

# 容器

### 容器消失

当前时间 - 最后以检测到容器的时间，如果大于 60 秒就认为容器消失

    time() - container_last_seen
    > 60

### 容器 cpu 的使用量

容器 CPU 使用率超过 80%。

    (
    sum by(instance, name) (rate(container_cpu_usage_seconds_total[3m])
    )
    * 100
    )
    > 80

### 容器内存的使用量

容器内存使用率超过 80%。

    (
    sum by(instance,name) (container_memory_working_set_bytes)
    /
    sum by(instance,name) (container_spec_memory_limit_bytes > 0)
    * 100
    )
    > 80

这篇文章 <https://mp.weixin.qq.com/s/rPxTBYmwG_7HnZRpRXMFuQ> 详细介绍了容器内触发 OOM 的机制，以及应该使用的监控指标。

### 容器磁盘的使用量

容器磁盘使用量超过 80%

    (1 - (sum(container_fs_inodes_free) by (instance)
    /
    sum(container_fs_inodes_total) by (instance))
    * 100)
    > 80

# 其他

[prometheus 告警指标 - 云 + 社区 - 腾讯云](https://cloud.tencent.com/developer/article/1758801)

## 主机和硬件监控

### 内存

节点内存压力大。主要页面故障率高

      \- alert: HostMemoryUnderMemoryPressure
        expr: rate(node\_vmstat\_pgmajfault\[1m\]) \> 1000
        for: 5m
        labels:
          severity: warning
        annotations:
          summary: Host memory under memory pressure (instance {{ $labels.instance }})
          description: The node is under heavy memory pressure. High rate of major page faults\\n  VALUE \= {{ $value }}\\n  LABELS: {{ $labels }}

### 主机中 inode 文件句柄报警

磁盘可用的 inode 快用完了（<10%）。

      \- alert: HostOutOfInodes
        expr: node\_filesystem\_files\_free{mountpoint \="/rootfs"} / node\_filesystem\_files{mountpoint \="/rootfs"} * 100 < 10
        for: 5m
        labels:
          severity: warning
        annotations:
          summary: Host out of inodes (instance {{ $labels.instance }})
          description: Disk is almost running out of available inodes ( 10% left)\\n  VALUE \= {{ $value }}\\n  LABELS: {{ $labels }}

### 主机 swap 分区使用

主机 swap 交换分区使用情况 (> 80%)

      \- alert: HostSwapIsFillingUp
        expr: (1 \- (node\_memory\_SwapFree\_bytes / node\_memory\_SwapTotal\_bytes)) * 100 \> 80
        for: 5m
        labels:
          severity: warning
        annotations:
          summary: Host swap is filling up (instance {{ $labels.instance }})
          description: Swap is filling up (\>80%)\\n  VALUE \= {{ $value }}\\n  LABELS: {{ $labels }}

### 主机物理元设备 (有的虚拟机可能没有此指标)

物理机温度过高

      \- alert: HostPhysicalComponentTooHot
        expr: node\_hwmon\_temp\_celsius \> 75
        for: 5m
        labels:
          severity: warning
        annotations:
          summary: Host physical component too hot (instance {{ $labels.instance }})
          description: Physical hardware component too hot\\n  VALUE \= {{ $value }}\\n  LABELS: {{ $labels }}

### 主机节点超温报警 (有的虚拟机可能没有此指标)

触发物理节点温度报警

      \- alert: HostNodeOvertemperatureAlarm
        expr: node\_hwmon\_temp\_alarm \== 1
        for: 5m
        labels:
          severity: critical
        annotations:
          summary: Host node overtemperature alarm (instance {{ $labels.instance }})
          description: Physical node temperature alarm triggered\\n  VALUE \= {{ $value }}\\n  LABELS: {{ $labels }}

### 主机 RAID 卡阵列失效 (虚拟机可能没有此指标)

RAID 阵列 {{$labels.device}} 由于一个或多个磁盘故障而处于退化状态。备用硬盘的数量不足以自动修复问题。

      \- alert: HostRaidArrayGotInactive
        expr: node\_md\_state{state\="inactive"} \> 0
        for: 5m
        labels:
          severity: critical
        annotations:
          summary: Host RAID array got inactive (instance {{ $labels.instance }})
          description: RAID array {{ $labels.device }} is in degraded state due to one or more disks failures. Number of spare drives is insufficient to fix issue automatically.\\n  VALUE \= {{ $value }}\\n  LABELS: {{ $labels }}

### 主机 RAID 磁盘故障 (虚拟机可能没有此指标)

在 {{ \Extra close brace or missing open bracelabels.md\_device }} 需要注意，可能需要进行磁盘更换

      \- alert: HostRaidDiskFailure
        expr: node\_md\_disks{state\="failed"} \> 0
        for: 5m
        labels:
          severity: warning
        annotations:
          summary: Host RAID disk failure (instance {{ $labels.instance }})
          description: At least one device in RAID array on {{ $labels.instance }} failed. Array {{ $labels.md\_device }} needs attention and possibly a disk swap\\n  VALUE \= {{ $value }}\\n  LABELS: {{ $labels }}

### 主机内核版本偏差

不同的内核版本正在运行

      \- alert: HostKernelVersionDeviations
        expr: count(sum(label\_replace(node\_uname\_info, "kernel", "$1", "release", "(\[0-9\]+.\[0-9\]+.\[0-9\]+).*")) by (kernel)) \> 1
        for: 5m
        labels:
          severity: warning
        annotations:
          summary: Host kernel version deviations (instance {{ $labels.instance }})
          description: Different kernel versions are running\\n  VALUE \= {{ $value }}\\n  LABELS: {{ $labels }}

### 检测到主机 EDAC 可纠正的错误

{{ \Extra close brace or missing open brace

      \- alert: HostEdacCorrectableErrorsDetected
        expr: increase(node\_edac\_correctable\_errors\_total\[5m\]) \> 0
        for: 5m
        labels:
          severity: info
        annotations:
          summary: Host EDAC Correctable Errors detected (instance {{ $labels.instance }})
          description: {{ $labels.instance }} has had {{ printf "%.0f" $value }} correctable memory errors reported by EDAC in the last 5 minutes.\\n  VALUE \= {{ $value }}\\n  LABELS: {{ $labels }}

### 检测到主机 EDAC 不正确的错误

{{ \Extra close brace or missing open brace

      \- alert: HostEdacUncorrectableErrorsDetected
        expr: node\_edac\_uncorrectable\_errors\_total \> 0
        for: 5m
        labels:
          severity: warning
        annotations:
          summary: Host EDAC Uncorrectable Errors detected (instance {{ $labels.instance }})
          description: {{ $labels.instance }} has had {{ printf "%.0f" $value }} uncorrectable memory errors reported by EDAC in the last 5 minutes.\\n  VALUE \= {{ $value }}\\n  LABELS: {{ $labels }}

## Redis 相关报警信息

### redis down

redis 服务 down 了，报警

      \- alert: RedisDown
        expr: redis\_up \== 0
        for: 5m
        labels:
          severity: critical
        annotations:
          summary: Redis down (instance {{ $labels.instance }})
          description: Redis instance is down\\n  VALUE \= {{ $value }}\\n  LABELS: {{ $labels }}

### redis 缺少主节点 (集群，或者 sentinel 模式才有)

redis 集群中缺少标记的主节点

      \- alert: RedisMissingMaster
        expr: (count(redis\_instance\_info{role\="master"}) or vector(0)) < 1
        for: 5m
        labels:
          severity: critical
        annotations:
          summary: Redis missing master (instance {{ $labels.instance }})
          description: Redis cluster has no node marked as master.\\n  VALUE \= {{ $value }}\\n  LABELS: {{ $labels }}

### Redis 主节点过多

redis 集群中被标记的主节点过多

      \- alert: RedisTooManyMasters
        expr: count(redis\_instance\_info{role\="master"}) \> 1
        for: 5m
        labels:
          severity: critical
        annotations:
          summary: Redis too many masters (instance {{ $labels.instance }})
          description: Redis cluster has too many nodes marked as master.\\n  VALUE \= {{ $value }}\\n  LABELS: {{ $labels }}

### Redis 复制中断

Redis 实例丢失了一个 slave

      \- alert: RedisReplicationBroken
        expr: delta(redis\_connected\_slaves\[1m\]) < 0
        for: 5m
        labels:
          severity: critical
        annotations:
          summary: Redis replication broken (instance {{ $labels.instance }})
          description: Redis instance lost a slave\\n  VALUE \= {{ $value }}\\n  LABELS: {{ $labels }}

### Redis 集群 flapping

在 Redis 副本连接中检测到变化。当复制节点失去与主节点的连接并重新连接（也就是 flapping）时，会发生这种情况。

      \- alert: RedisClusterFlapping
        expr: changes(redis\_connected\_slaves\[5m\]) \> 2
        for: 5m
        labels:
          severity: critical
        annotations:
          summary: Redis cluster flapping (instance {{ $labels.instance }})
          description: Changes have been detected in Redis replica connection. This can occur when replica nodes lose connection to the master and reconnect (a.k.a flapping).\\n  VALUE \= {{ $value }}\\n  LABELS: {{ $labels }}

### Redis 缺少备份

Redis 已经有 24 小时没有备份了。

      \- alert: RedisMissingBackup
        expr: time() \- redis\_rdb\_last\_save\_timestamp\_seconds \> 60 * 60 * 24
        for: 5m
        labels:
          severity: critical
        annotations:
          summary: Redis missing backup (instance {{ $labels.instance }})
          description: Redis has not been backuped for 24 hours\\n  VALUE \= {{ $value }}\\n  LABELS: {{ $labels }}

### Redis 内存不足

Redis 内存耗尽（>90%）。

    #需要 redis 实例设置 maxmemory maxmemory\-policy 最大使用内存参数
    \- alert: RedisOutOfMemory
        expr: redis\_memory\_used\_bytes / redis\_total\_system\_memory\_bytes * 100 \> 90
        for: 5m
        labels:
          severity: warning
        annotations:
          summary: Redis out of memory (instance {{ $labels.instance }})
          description: Redis is running out of memory (\> 90%)\\n  VALUE \= {{ $value }}\\n  LABELS: {{ $labels }}

### Redis 连接数过多

Redis 实例有太多的连接

      \- alert: RedisTooManyConnections
        expr: redis\_connected\_clients \> 100
        for: 5m
        labels:
          severity: warning
        annotations:
          summary: Redis too many connections (instance {{ $labels.instance }})
          description: Redis instance has too many connections\\n  VALUE \= {{ $value }}\\n  LABELS: {{ $labels }}

### Redis 连接数不足

Redis 实例应该有更多的连接（> 5）。

      \- alert: RedisNotEnoughConnections
        expr: redis\_connected\_clients < 5
        for: 5m
        labels:
          severity: warning
        annotations:
          summary: Redis not enough connections (instance {{ $labels.instance }})
          description: Redis instance should have more connections (\> 5)\\n  VALUE \= {{ $value }}\\n  LABELS: {{ $labels }}

### Redis 拒绝连接

一些与 Redis 的连接已被拒绝

      \- alert: RedisRejectedConnections
        expr: increase(redis\_rejected\_connections\_total\[1m\]) \> 0
        for: 5m
        labels:
          severity: critical
        annotations:
          summary: Redis rejected connections (instance {{ $labels.instance }})
          description: Some connections to Redis has been rejected\\n  VALUE \= {{ $value }}\\n  LABELS: {{ $labels }}

## rabbitmq 监控 : \[rabbitmq/rabbitmq-prometheus ]

### rabbitmq 节点 down

节点数量少于 1 个

      \- alert: RabbitmqNodeDown
        expr: sum(rabbitmq\_build\_info) < 3
        for: 5m
        labels:
          severity: critical
        annotations:
          summary: Rabbitmq node down (instance {{ $labels.instance }})
          description: Less than 3 nodes running in RabbitMQ cluster\\n  VALUE \= {{ $value }}\\n  LABELS: {{ $labels }}  \- alert: RabbitmqDown
        expr: rabbitmq\_up \== 0
        for: 5m
        labels:
          severity: critical
        annotations:
          summary: Rabbitmq down (instance {{ $labels.instance }})
          description: RabbitMQ node down\\n  VALUE \= {{ $value }}\\n  LABELS: {{ $labels }}

### Rabbitmq 实例的不同版本

在同一集群中运行不同版本的 Rabbitmq，可能会导致失败。

      \- alert: RabbitmqInstancesDifferentVersions
        expr: count(count(rabbitmq\_build\_info) by (rabbitmq\_version)) \> 1
        for: 5m
        labels:
          severity: warning
        annotations:
          summary: Rabbitmq instances different versions (instance {{ $labels.instance }})
          description: Running different version of Rabbitmq in the same cluster, can lead to failure.\\n  VALUE \= {{ $value }}\\n  LABELS: {{ $labels }}  \- alert: RabbitmqClusterPartition
        expr: rabbitmq\_partitions \> 0
        for: 5m
        labels:
          severity: critical
        annotations:
          summary: Rabbitmq cluster partition (instance {{ $labels.instance }})
          description: Cluster partition\\n  VALUE \= {{ $value }}\\n  LABELS: {{ $labels }}

### Rabbitmq 内存高

一个节点使用了 90% 以上的内存分配。

      \- alert: RabbitmqMemoryHigh
        expr: rabbitmq\_process\_resident\_memory\_bytes / rabbitmq\_resident\_memory\_limit\_bytes * 100 \> 90
        for: 5m
        labels:
          severity: warning
        annotations:
          summary: Rabbitmq memory high (instance {{ $labels.instance }})
          description: A node use more than 90% of allocated RAM\\n  VALUE \= {{ $value }}\\n  LABELS: {{ $labels }}  \- alert: RabbitmqOutOfMemory
        expr: rabbitmq\_node\_mem\_used / rabbitmq\_node\_mem\_limit * 100 \> 90
        for: 5m
        labels:
          severity: warning
        annotations:
          summary: Rabbitmq out of memory (instance {{ $labels.instance }})
          description: Memory available for RabbmitMQ is low ( 10%)\\n  VALUE \= {{ $value }}\\n  LABELS: {{ $labels }}

### Rabbitmq 文件描述符的用法

一个节点使用 90% 以上的文件描述符。

      \- alert: RabbitmqFileDescriptorsUsage
        expr: rabbitmq\_process\_open\_fds / rabbitmq\_process\_max\_fds * 100 \> 90
        for: 5m
        labels:
          severity: warning
        annotations:
          summary: Rabbitmq file descriptors usage (instance {{ $labels.instance }})
          description: A node use more than 90% of file descriptors\\n  VALUE \= {{ $value }}\\n  LABELS: {{ $labels }}  \- alert: RabbitmqTooManyConnections
        expr: rabbitmq\_connectionsTotal \> 1000
        for: 5m
        labels:
          severity: warning
        annotations:
          summary: Rabbitmq too many connections (instance {{ $labels.instance }})
          description: RabbitMQ instance has too many connections (\> 1000)\\n  VALUE \= {{ $value }}\\n  LABELS: {{ $labels }}

### Rabbitmq 连接数太多

节点的总连接数过高。

      \- alert: RabbitmqTooMuchConnections
        expr: rabbitmq\_connections \> 1000
        for: 5m
        labels:
          severity: warning
        annotations:
          summary: Rabbitmq too much connections (instance {{ $labels.instance }})
          description: The total connections of a node is too high\\n  VALUE \= {{ $value }}\\n  LABELS: {{ $labels }}  \- alert: RabbitmqTooManyMessagesInQueue
        expr: rabbitmq\_queue\_messages\_ready{queue\="my-queue"} \> 1000
        for: 5m
        labels:
          severity: warning
        annotations:
          summary: Rabbitmq too many messages in queue (instance {{ $labels.instance }})
          description: Queue is filling up (\> 1000 msgs)\\n  VALUE \= {{ $value }}\\n  LABELS: {{ $labels }}

### Rabbitmq 无队列消费

一个队列的消费者少于 1 个

      \- alert: RabbitmqNoQueueConsumer
        expr: rabbitmq\_queue\_consumers < 1
        for: 5m
        labels:
          severity: warning
        annotations:
          summary: Rabbitmq no queue consumer (instance {{ $labels.instance }})
          description: A queue has less than 1 consumer\\n  VALUE \= {{ $value }}\\n  LABELS: {{ $labels }}  \- alert: RabbitmqSlowQueueConsuming
        expr: time() \- rabbitmq\_queue\_head\_message\_timestamp{queue\="my-queue"} \> 60
        for: 5m
        labels:
          severity: warning
        annotations:
          summary: Rabbitmq slow queue consuming (instance {{ $labels.instance }})
          description: Queue messages are consumed slowly (\> 60s)\\n  VALUE \= {{ $value }}\\n  LABELS: {{ $labels }}

### Rabbitmq 不可路由的消息

一个队列有不可更改的消息

      \- alert: RabbitmqUnroutableMessages
        expr: increase(rabbitmq\_channel\_messages\_unroutable\_returned\_total\[5m\]) \> 0 or increase(rabbitmq\_channel\_messages\_unroutable\_dropped\_total\[5m\]) \> 0
        for: 5m
        labels:
          severity: warning
        annotations:
          summary: Rabbitmq unroutable messages (instance {{ $labels.instance }})
          description: A queue has unroutable messages\\n  VALUE \= {{ $value }}\\n  LABELS: {{ $labels }}  \- alert: RabbitmqNoConsumer
        expr: rabbitmq\_queue\_consumers \== 0
        for: 5m
        labels:
          severity: critical
        annotations:
          summary: Rabbitmq no consumer (instance {{ $labels.instance }})
          description: Queue has no consumer\\n  VALUE \= {{ $value }}\\n  LABELS: {{ $labels }}

##
