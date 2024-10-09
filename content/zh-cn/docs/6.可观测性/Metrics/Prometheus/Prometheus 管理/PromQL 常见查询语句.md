---
title: PromQL 常见查询语句
linkTitle: PromQL 常见查询语句
date: 2024-06-08T09:46
weight: 2
---

# 概述

> 参考：
>
> - [GitHub 项目，samber/awesome-prometheus-alerts](https://github.com/samber/awesome-prometheus-alerts)
>   - https://samber.github.io/awesome-prometheus-alerts/
> - [腾讯云+社区，prometheus 告警指标](https://cloud.tencent.com/developer/article/1758801)
> - [公众号，云原生小白-监控容器 OOMKill 的正确指标](https://mp.weixin.qq.com/s/rPxTBYmwG_7HnZRpRXMFuQ)
> - https://panzhongxian.cn/cn/2023/09/grafana-pannel-skills/ Grafana 常用但难配的图表。一些真实场景的查询语句写法以及对应 Grafana 图标如何用

# 问题

如何获取范围向量中的第一个和最后一个值。 https://stackoverflow.com/questions/68895729/how-to-get-the-first-and-last-element-of-a-range-vector-in-promql

- MetricsQS 中有 `first_over_time()` 函数

如何获取范围向量中，指定的值。 https://stackoverflow.com/questions/45213745/prometheus-how-to-calculate-proportion-of-single-value-over-time ，比如 `count_over_time(my_metric[1m] != 0)`获取 1 分钟内所有值中不为 0 的值

- MetricsQL 中有 `count_ne_over_time(my_metric[1h], 0)` 函数

# SLO/SLI

## 根据过去一段时间的统计数据监测异常值

参考 [Statistics](/docs/5.数据存储/Statistics/Statistics.md) 中的 “检测和处理异常值”，使用 **Z-Score** 法，通过下面的公式实现

$$
Z-score = \frac{x - \mu}{\sigma}
$$

- **x** 是当前值
- **μ** 是总体的 **mean(平均值)**
- **σ** 是总体的 **standard deviation(标准差)**。

> Tips: 这里的 population(总体) 的意思对应到 Prometheus 中就是指 **范围向量**

这里使用网卡流量速率举例说明

为了使用 Z-score 方法来检测网卡流量的异常情况，需要完成以下几个步骤：

1. **计算过去 n 小时的平均值和标准差**。
2. **计算当前值与平均值的差异，并标准化**。
  1. Notes: 在统计学中，“标准化” 通常指的是将数据转换为具有特定性质的标准形式，以便进行比较或进一步分析。具体来说，标准化数据意味着将数据调整为均值为 0、标准差为 1 的形式。这通常是通过计算 Z-score 来实现的
3. **根据标准化的值（Z-score）来判断是否异常**。

下面是一个 PromQL 示例，假设想计算过去 1 小时的平均值和标准差，并与当前值进行对比：

```promql
# 计算当前值
current_value = irate(hdf_hdf_network_receive_bytes_total[15m])

# 计算过去 1 小时的平均值
avg_over_time(irate(hdf_hdf_network_receive_bytes_total[15m])[1h:])

# 计算过去 1 小时的标准差
stddev_over_time(irate(hdf_hdf_network_receive_bytes_total[15m])[1h:])

# 计算 Z-score
z_score = (current_value - avg_over_time(irate(hdf_hdf_network_receive_bytes_total[15m])[1h:])) / stddev_over_time(irate(hdf_hdf_network_receive_bytes_total[15m])[1h:])
```

为了判断是否异常，需要设定一个阈值，通常 Z-score 大于 3 或小于 -3 被认为是异常的（下面使用 abs 取绝对值）

```promql
abs(
  (
    irate(hdf_hdf_network_receive_bytes_total[15m])
    -
    avg_over_time(irate(hdf_hdf_network_receive_bytes_total[15m])[1h:])
  )
  /
  stddev_over_time(irate(hdf_hdf_network_receive_bytes_total[15m])[1h:])
)
> 3
```

这个查询将返回当前网卡流量接收字节数是否与过去 1 小时的平均值相比存在显著异常。如果想使用不同的时间窗口，只需调整 `[1h]` 为需要的值，比如 `[2h]` 或 `[30m]`。

### 网卡收/发流量速率变化异常

```promql
abs(
  (
    irate(hdf_hdf_network_receive_bytes_total[15m])
    -
    avg_over_time(irate(hdf_hdf_network_receive_bytes_total[15m])[6h:])
  )
  /
  stddev_over_time(irate(hdf_hdf_network_receive_bytes_total[15m])[6h:])
)
> 3
```

# 通用

### 服务中断了多少时间

这个前提是，服务状态为 0 或者 1。

```bash
avg_over_time(up{}[1h])
```

此时，如果结果为 0.9，则表示服务在 90% 的时间处于运行状态。如果按照 1 小时算，则这 1 小时中，有 6 分钟是停机的。

# 数通设备资源查询语句

查询端口接收的实时带宽。注意:是带宽

```promql
irate(ifHCInOctets{instance="IP.IP.IP.IP",ifAlias="XXXX"}[6m]) * 8
```

查询端口发送的实时带宽。注意:是带宽

```promql
irate(ifHCOutOctets{instance="IP.IP.IP.IP",ifAlias="XXXX"}[6m]) * 8
```

# 物理机资源查询语句

## CPU

### CPU 的使用率

显示 cpu 的每个逻辑 core 的使用率

```text
avg(irate(node_cpu_seconds_total{mode="idle"}[5m])) by(instance,job)
* 100 < 20
```

查询物理机 CPU 的使用率，显示总体使用率

```text
100 - avg (irate(node_cpu_seconds_total{instance="XXXX",mode="idle"}[5m])) by (instance) * 100
```

### 上下文切换越来越多

```text
(rate(node_context_switches_total[5m])) / (count (node_cpu_seconds_total{mode="idle"}) without(cpu, mode)) > 1000
```

- `rate(node_context_switches_total[5m]` # 设备上下文切换在 5 分钟之间的变化量
- `(count (node_cpu_seconds_total{mode="idle"}) without(cpu, mode))` # 获取 instance 的 CPU 总核数
- 两个序列相除，即可获得每个 CPU 核上，5 分钟的上下文切换次数的变化量

## 内存

### 内存使用率

```text
node_memory_MemAvailable_bytes{}
/
node_memory_MemTotal_bytes{}
* 100 < 10
```

### OOM

检测主机是否发生了 oom

```text
increase(node_vmstat_oom_kill[5m]) > 0
```

## 磁盘

### 磁盘使用率

不能直接用 node_filesystem_avail_bytes / node_filesystem_size_bytes，需要通过 node_filesystem_free_bytes 作为中转，把 inode 等系统占用的磁盘空间也算上。否则告警不准。

```promql
(
  node_filesystem_size_bytes{fstype=~"ext.*|xfs|nfs",mountpoint!~".*pod.*"}
  -
  node_filesystem_free_bytes{fstype=~"ext.*|xfs|nfs",mountpoint!~".*pod.*"}
)
/
(
  node_filesystem_avail_bytes{fstype=~"ext.*|xfs|nfs",mountpoint!~".*pod.*"}
  +
  (
    node_filesystem_size_bytes{fstype=~"ext.*|xfs|nfs",mountpoint!~".*pod.*"}
    -
    node_filesystem_free_bytes{fstype=~"ext.*|xfs|nfs",mountpoint!~".*pod.*"}
  )
) * 100
```

### 磁盘将满

根据磁盘 1 小时的变化速率，预测 4 小时内会不会被写满

```promql
predict_linear(node_filesystem_free_bytes{fstype!~"tmpfs"}[1h], 4 * 3600) < 0
```

### IO 使用率

```text
100
-
(avg(irate(node_disk_io_time_seconds_total[5m])) by(instance,job))
* 100 < 20
```

### 读写速率

读取速率

```text
sum by (instance) (irate(node_disk_read_bytes_total[2m]))
/ 1024 / 1024 > 200
```

写入速率

```text
sum by (instance) (irate(node_disk_written_bytes_total[2m]))
/ 1024 / 1024 > 200
```

### 读写延迟

读取延迟

```yaml
rate(node_disk_read_time_seconds_total[1m])
/
rate(node_disk_reads_completed_total[1m])
> 0.1 and rate(node_disk_reads_completed_total[1m])
> 0
```

写入延迟

```text
rate(node_disk_write_time_seconds_total[1m])
/
rate(node_disk_writes_completed_total[1m])
> 0.1 and rate(node_disk_writes_completed_total[1m])
> 0
```

### inode

```promql
(
1 - node_filesystem_files_free{fstype=~"ext4|xfs"}
/
node_filesystem_files{fstype=~"ext4|xfs"}
) * 100
```

## 网络

### 流量过高

接收和发送的带宽大于 2.5GiB/s 时告警

**接收**

```yaml
max(rate(node_network_receive_bytes_total{device!~'tap.*|veth.*|br.*|docker.*|virbr*|lo*'}[2m])) by(instance,job)
/ 1024 / 1024 / 1024 * 8 > 2.5
```

**发送**

```yaml
max(rate(node_network_transmit_bytes_total{device!~'tap.*|veth.*|br.*|docker.*|virbr*|lo*'}[2m])) by(instance,job)
/ 1024 / 1024 / 1024 * 8 > 2.5
```

### 错误包过多

接收

```yaml
increase(node_network_receive_errs_total[5m]) > 0
```

发送

```yaml
increase(node_network_transmit_errs_total[5m]) > 0
```

### TCP_ESTABLISHED 过高

```yaml
node_netstat_Tcp_CurrEstab > 50000
```

这篇文章 <https://mp.weixin.qq.com/s/rPxTBYmwG_7HnZRpRXMFuQ> 详细介绍了容器内触发 OOM 的机制，以及应该使用的监控指标。

# 其他

[prometheus 告警指标 - 云 + 社区 - 腾讯云](https://cloud.tencent.com/developer/article/1758801)

### 主机和硬件监控

#### 内存

节点内存压力大。主要页面故障率高

```yaml
  - alert: HostMemoryUnderMemoryPressure
    expr: rate(node_vmstat_pgmajfault[1m]) > 1000
    for: 5m
    labels:
      severity: warning
    annotations:
      summary: Host memory under memory pressure (instance {{ $labels.instance }})
      description: The node is under heavy memory pressure. High rate of major page faults\n  VALUE = {{ $value }}\n  LABELS: {{ $labels }}
```

#### 主机网络接口流入流量异常

主机网络接口可能接收了太多的数据（\> 100 MB/s）。阀值根据自己机器背板网卡决定

```yaml
  - alert: HostUnusualNetworkThroughputIn
    expr: sum by (instance) (rate(node_network_receive_bytes_total[2m])) / 1024 / 1024 > 100
    for: 5m
    labels:
      severity: warning
    annotations:
      summary: Host unusual network throughput in (instance {{ $labels.instance }})
      description: Host network interfaces are probably receiving too much data (> 100 MB/s)\n  VALUE = {{ $value }}\n  LABELS: {{ $labels }}
```

#### 主机网络接口流出流量异常

主机网络接口可能发送了太多的数据（\> 100 MB/s）。

```yaml
  - alert: HostUnusualNetworkThroughputOut
    expr: sum by (instance) (rate(node_network_transmit_bytes_total[2m])) / 1024 / 1024 > 100
    for: 5m
    labels:
      severity: warning
    annotations:
      summary: Host unusual network throughput out (instance {{ $labels.instance }})
      description: Host network interfaces are probably sending too much data (> 100 MB/s)\n  VALUE = {{ $value }}\n  LABELS: {{ $labels }}
```
 
#### 主机 swap 分区使用

主机 swap 交换分区使用情况 (> 80%)

```yaml
  - alert: HostSwapIsFillingUp
    expr: (1 - (node_memory_SwapFree_bytes / node_memory_SwapTotal_bytes)) * 100 > 80
    for: 5m
    labels:
      severity: warning
    annotations:
      summary: Host swap is filling up (instance {{ $labels.instance }})
      description: Swap is filling up (>80%)\n  VALUE = {{ $value }}\n  LABELS: {{ $labels }}
```

#### 主机 systemctl 管理的服务 down 了

主机上systemctl 管理的服务不正常，failed了，根据自己的实际情况来判断哪些服务

```yaml
  - alert: HostSystemdServiceCrashed
    expr: node_systemd_unit_state{state="failed"} == 1
    for: 5m
    labels:
      severity: warning
    annotations:
      summary: Host SystemD service crashed (instance {{ $labels.instance }})
      description: SystemD service crashed\n  VALUE = {{ $value }}\n  LABELS: {{ $labels }}
```

#### 主机物理元设备(有的虚拟机可能没有此指标)

物理机温度过高

```yaml
  - alert: HostPhysicalComponentTooHot
    expr: node_hwmon_temp_celsius > 75
    for: 5m
    labels:
      severity: warning
    annotations:
      summary: Host physical component too hot (instance {{ $labels.instance }})
      description: Physical hardware component too hot\n  VALUE = {{ $value }}\n  LABELS: {{ $labels }}
```

#### 主机节点超温报警(有的虚拟机可能没有此指标)

触发物理节点温度报警

```yaml
  - alert: HostNodeOvertemperatureAlarm
    expr: node_hwmon_temp_alarm == 1
    for: 5m
    labels:
      severity: critical
    annotations:
      summary: Host node overtemperature alarm (instance {{ $labels.instance }})
      description: Physical node temperature alarm triggered\n  VALUE = {{ $value }}\n  LABELS: {{ $labels }}
```

#### 主机RAID 卡阵列失效(虚拟机可能没有此指标)

RAID阵列{{$labels.device }}由于一个或多个磁盘故障而处于退化状态。备用硬盘的数量不足以自动修复问题。

```yaml
  - alert: HostRaidArrayGotInactive
    expr: node_md_state{state="inactive"} > 0
    for: 5m
    labels:
      severity: critical
    annotations:
      summary: Host RAID array got inactive (instance {{ $labels.instance }})
      description: RAID array {{ $labels.device }} is in degraded state due to one or more disks failures. Number of spare drives is insufficient to fix issue automatically.\n  VALUE = {{ $value }}\n  LABELS: {{ $labels }}
```

#### 主机RAID磁盘故障(虚拟机可能没有此指标)

在{{ \Extra close brace or missing open bracelabels.md_device }}需要注意，可能需要进行磁盘更换

```yaml
  - alert: HostRaidDiskFailure
    expr: node_md_disks{state="failed"} > 0
    for: 5m
    labels:
      severity: warning
    annotations:
      summary: Host RAID disk failure (instance {{ $labels.instance }})
      description: At least one device in RAID array on {{ $labels.instance }} failed. Array {{ $labels.md_device }} needs attention and possibly a disk swap\n  VALUE = {{ $value }}\n  LABELS: {{ $labels }}
```

#### 主机内核版本偏差

不同的内核版本正在运行

```yaml
  - alert: HostKernelVersionDeviations
    expr: count(sum(label_replace(node_uname_info, "kernel", "$1", "release", "([0-9]+.[0-9]+.[0-9]+).*")) by (kernel)) > 1
    for: 5m
    labels:
      severity: warning
    annotations:
      summary: Host kernel version deviations (instance {{ $labels.instance }})
      description: Different kernel versions are running\n  VALUE = {{ $value }}\n  LABELS: {{ $labels }}
```

#### 检测主机 OOM 杀进程

```yaml
  - alert: HostOomKillDetected
    expr: increase(node_vmstat_oom_kill[5m]) > 0
    for: 5m
    labels:
      severity: warning
    annotations:
      summary: Host OOM kill detected (instance {{ $labels.instance }})
      description: OOM kill detected\n  VALUE = {{ $value }}\n  LABELS: {{ $labels }}
```

#### 检测到主机EDAC可纠正的错误

{{ \Extra close brace or missing open brace

```yaml
  - alert: HostEdacCorrectableErrorsDetected
    expr: increase(node_edac_correctable_errors_total[5m]) > 0
    for: 5m
    labels:
      severity: info
    annotations:
      summary: Host EDAC Correctable Errors detected (instance {{ $labels.instance }})
      description: {{ $labels.instance }} has had {{ printf "%.0f" $value }} correctable memory errors reported by EDAC in the last 5 minutes.\n  VALUE = {{ $value }}\n  LABELS: {{ $labels }}
```

#### 检测到主机EDAC不正确的错误

```yaml
  - alert: HostEdacUncorrectableErrorsDetected
    expr: node_edac_uncorrectable_errors_total > 0
    for: 5m
    labels:
      severity: warning
    annotations:
      summary: Host EDAC Uncorrectable Errors detected (instance {{ $labels.instance }})
      description: {{ $labels.instance }} has had {{ printf "%.0f" $value }} uncorrectable memory errors reported by EDAC in the last 5 minutes.\n  VALUE = {{ $value }}\n  LABELS: {{ $labels }}
```

### Docker 容器

#### 一个容器消失

```yaml
  - alert: ContainerKilled
    expr: time() - container_last_seen > 60
    for: 5m
    labels:
      severity: warning
    annotations:
      summary: Container killed (instance {{ $labels.instance }})
      description: A container has disappeared\n  VALUE = {{ $value }}\n  LABELS: {{ $labels }}
```

#### 容器 cpu 的使用量

[容器](https://cloud.tencent.com/product/tke?from_column=20065&from=20065)CPU使用率超过80%。

```yaml
 # cAdvisor有时会消耗大量的CPU，所以这个警报会不断地响起。
  # If you want to exclude it from this alert, just use: container_cpu_usage_seconds_total{name!=""}
  - alert: ContainerCpuUsage
    expr: (sum(rate(container_cpu_usage_seconds_total[3m])) BY (instance, name) * 100) > 80
    for: 5m
    labels:
      severity: warning
    annotations:
      summary: Container CPU usage (instance {{ $labels.instance }})
      description: Container CPU usage is above 80%\n  VALUE = {{ $value }}\n  LABELS: {{ $labels }}
```

#### 容器内存的使用量

容器内存使用率超过 80%。

```yaml
 # See https://medium.com/faun/how-much-is-too-much-the-linux-oomkiller-and-used-memory-d32186f29c9d
  - alert: ContainerMemoryUsage
    expr: (sum(container_memory_working_set_bytes) BY (instance, name) / sum(container_spec_memory_limit_bytes > 0) BY (instance, name) * 100) > 80
    for: 5m
    labels:
      severity: warning
    annotations:
      summary: Container Memory usage (instance {{ $labels.instance }})
      description: Container Memory usage is above 80%\n  VALUE = {{ $value }}\n  LABELS: {{ $labels }}
```

#### 容器磁盘的使用量

容器磁盘使用量超过 80%

```yaml
  - alert: ContainerVolumeUsage
    expr: (1 - (sum(container_fs_inodes_free) BY (instance) / sum(container_fs_inodes_total) BY (instance)) * 100) > 80
    for: 5m
    labels:
      severity: warning
    annotations:
      summary: Container Volume usage (instance {{ $labels.instance }})
      description: Container Volume usage is above 80%\n  VALUE = {{ $value }}\n  LABELS: {{ $labels }}
```

### Redis 相关报警信息

#### redis down

[redis](https://cloud.tencent.com/product/crs?from_column=20065&from=20065) 服务 down 了，报警

```yaml
  - alert: RedisDown
    expr: redis_up == 0
    for: 5m
    labels:
      severity: critical
    annotations:
      summary: Redis down (instance {{ $labels.instance }})
      description: Redis instance is down\n  VALUE = {{ $value }}\n  LABELS: {{ $labels }}
```

#### redis 缺少主节点(集群，或者sentinel 模式才有)

redis 集群中缺少标记的主节点

```yaml
  - alert: RedisMissingMaster
    expr: (count(redis_instance_info{role="master"}) or vector(0)) < 1
    for: 5m
    labels:
      severity: critical
    annotations:
      summary: Redis missing master (instance {{ $labels.instance }})
      description: Redis cluster has no node marked as master.\n  VALUE = {{ $value }}\n  LABELS: {{ $labels }}
```

#### Redis 主节点过多

redis 集群中被标记的主节点过多

```yaml
  - alert: RedisTooManyMasters
    expr: count(redis_instance_info{role="master"}) > 1
    for: 5m
    labels:
      severity: critical
    annotations:
      summary: Redis too many masters (instance {{ $labels.instance }})
      description: Redis cluster has too many nodes marked as master.\n  VALUE = {{ $value }}\n  LABELS: {{ $labels }}
```

#### Redis 复制中断

Redis实例丢失了一个slave

```yaml
  - alert: RedisReplicationBroken
    expr: delta(redis_connected_slaves[1m]) < 0
    for: 5m
    labels:
      severity: critical
    annotations:
      summary: Redis replication broken (instance {{ $labels.instance }})
      description: Redis instance lost a slave\n  VALUE = {{ $value }}\n  LABELS: {{ $labels }}
```

#### Redis 集群 flapping

在Redis副本连接中检测到变化。当复制节点失去与主节点的连接并重新连接（也就是flapping）时，会发生这种情况。

```yaml
  - alert: RedisClusterFlapping
    expr: changes(redis_connected_slaves[5m]) > 2
    for: 5m
    labels:
      severity: critical
    annotations:
      summary: Redis cluster flapping (instance {{ $labels.instance }})
      description: Changes have been detected in Redis replica connection. This can occur when replica nodes lose connection to the master and reconnect (a.k.a flapping).\n  VALUE = {{ $value }}\n  LABELS: {{ $labels }}
```

#### Redis缺少备份

Redis已经有24小时没有备份了。

```yaml
  - alert: RedisMissingBackup
    expr: time() - redis_rdb_last_save_timestamp_seconds > 60 * 60 * 24
    for: 5m
    labels:
      severity: critical
    annotations:
      summary: Redis missing backup (instance {{ $labels.instance }})
      description: Redis has not been backuped for 24 hours\n  VALUE = {{ $value }}\n  LABELS: {{ $labels }}
```

#### Redis内存不足

Redis内存耗尽（>90%）。

```
#需要 redis 实例设置 maxmemory maxmemory-policy 最大使用内存参数
- alert: RedisOutOfMemory
    expr: redis_memory_used_bytes / redis_total_system_memory_bytes * 100 > 90
    for: 5m
    labels:
      severity: warning
    annotations:
      summary: Redis out of memory (instance {{ $labels.instance }})
      description: Redis is running out of memory (> 90%)\n  VALUE = {{ $value }}\n  LABELS: {{ $labels }}
```

#### Redis连接数过多

Redis实例有太多的连接

```yaml
  - alert: RedisTooManyConnections
    expr: redis_connected_clients > 100
    for: 5m
    labels:
      severity: warning
    annotations:
      summary: Redis too many connections (instance {{ $labels.instance }})
      description: Redis instance has too many connections\n  VALUE = {{ $value }}\n  LABELS: {{ $labels }}
```

#### Redis连接数不足

Redis实例应该有更多的连接（> 5）。

```yaml
  - alert: RedisNotEnoughConnections
    expr: redis_connected_clients < 5
    for: 5m
    labels:
      severity: warning
    annotations:
      summary: Redis not enough connections (instance {{ $labels.instance }})
      description: Redis instance should have more connections (> 5)\n  VALUE = {{ $value }}\n  LABELS: {{ $labels }}
```

#### Redis拒绝连接

一些与Redis的连接已被拒绝

```yaml
  - alert: RedisRejectedConnections
    expr: increase(redis_rejected_connections_total[1m]) > 0
    for: 5m
    labels:
      severity: critical
    annotations:
      summary: Redis rejected connections (instance {{ $labels.instance }})
      description: Some connections to Redis has been rejected\n  VALUE = {{ $value }}\n  LABELS: {{ $labels }}
```

### rabbitmq 监控 : \[rabbitmq/rabbitmq-prometheus \]

#### rabbitmq 节点 down

节点数量少于 1 个

```yaml
  - alert: RabbitmqNodeDown
    expr: sum(rabbitmq_build_info) < 3
    for: 5m
    labels:
      severity: critical
    annotations:
      summary: Rabbitmq node down (instance {{ $labels.instance }})
      description: Less than 3 nodes running in RabbitMQ cluster\n  VALUE = {{ $value }}\n  LABELS: {{ $labels }}  - alert: RabbitmqDown
    expr: rabbitmq_up == 0
    for: 5m
    labels:
      severity: critical
    annotations:
      summary: Rabbitmq down (instance {{ $labels.instance }})
      description: RabbitMQ node down\n  VALUE = {{ $value }}\n  LABELS: {{ $labels }}
```

#### Rabbitmq实例的不同版本

在同一集群中运行不同版本的Rabbitmq，可能会导致失败。

```yaml
  - alert: RabbitmqInstancesDifferentVersions
    expr: count(count(rabbitmq_build_info) by (rabbitmq_version)) > 1
    for: 5m
    labels:
      severity: warning
    annotations:
      summary: Rabbitmq instances different versions (instance {{ $labels.instance }})
      description: Running different version of Rabbitmq in the same cluster, can lead to failure.\n  VALUE = {{ $value }}\n  LABELS: {{ $labels }}  - alert: RabbitmqClusterPartition
    expr: rabbitmq_partitions > 0
    for: 5m
    labels:
      severity: critical
    annotations:
      summary: Rabbitmq cluster partition (instance {{ $labels.instance }})
      description: Cluster partition\n  VALUE = {{ $value }}\n  LABELS: {{ $labels }}
```

#### Rabbitmq内存高

一个节点使用了90%以上的内存分配。

```yaml
  - alert: RabbitmqMemoryHigh
    expr: rabbitmq_process_resident_memory_bytes / rabbitmq_resident_memory_limit_bytes * 100 > 90
    for: 5m
    labels:
      severity: warning
    annotations:
      summary: Rabbitmq memory high (instance {{ $labels.instance }})
      description: A node use more than 90% of allocated RAM\n  VALUE = {{ $value }}\n  LABELS: {{ $labels }}  - alert: RabbitmqOutOfMemory
    expr: rabbitmq_node_mem_used / rabbitmq_node_mem_limit * 100 > 90
    for: 5m
    labels:
      severity: warning
    annotations:
      summary: Rabbitmq out of memory (instance {{ $labels.instance }})
      description: Memory available for RabbmitMQ is low (< 10%)\n  VALUE = {{ $value }}\n  LABELS: {{ $labels }}
```

#### Rabbitmq文件描述符的用法

一个节点使用90%以上的文件描述符。

```yaml
  - alert: RabbitmqFileDescriptorsUsage
    expr: rabbitmq_process_open_fds / rabbitmq_process_max_fds * 100 > 90
    for: 5m
    labels:
      severity: warning
    annotations:
      summary: Rabbitmq file descriptors usage (instance {{ $labels.instance }})
      description: A node use more than 90% of file descriptors\n  VALUE = {{ $value }}\n  LABELS: {{ $labels }}  - alert: RabbitmqTooManyConnections
    expr: rabbitmq_connectionsTotal > 1000
    for: 5m
    labels:
      severity: warning
    annotations:
      summary: Rabbitmq too many connections (instance {{ $labels.instance }})
      description: RabbitMQ instance has too many connections (> 1000)\n  VALUE = {{ $value }}\n  LABELS: {{ $labels }}
```

#### Rabbitmq连接数太多

节点的总连接数过高。

```yaml
  - alert: RabbitmqTooMuchConnections
    expr: rabbitmq_connections > 1000
    for: 5m
    labels:
      severity: warning
    annotations:
      summary: Rabbitmq too much connections (instance {{ $labels.instance }})
      description: The total connections of a node is too high\n  VALUE = {{ $value }}\n  LABELS: {{ $labels }}  - alert: RabbitmqTooManyMessagesInQueue
    expr: rabbitmq_queue_messages_ready{queue="my-queue"} > 1000
    for: 5m
    labels:
      severity: warning
    annotations:
      summary: Rabbitmq too many messages in queue (instance {{ $labels.instance }})
      description: Queue is filling up (> 1000 msgs)\n  VALUE = {{ $value }}\n  LABELS: {{ $labels }}
```

#### Rabbitmq无队列消费

一个队列的消费者少于1个

```yaml
  - alert: RabbitmqNoQueueConsumer
    expr: rabbitmq_queue_consumers < 1
    for: 5m
    labels:
      severity: warning
    annotations:
      summary: Rabbitmq no queue consumer (instance {{ $labels.instance }})
      description: A queue has less than 1 consumer\n  VALUE = {{ $value }}\n  LABELS: {{ $labels }}  - alert: RabbitmqSlowQueueConsuming
    expr: time() - rabbitmq_queue_head_message_timestamp{queue="my-queue"} > 60
    for: 5m
    labels:
      severity: warning
    annotations:
      summary: Rabbitmq slow queue consuming (instance {{ $labels.instance }})
      description: Queue messages are consumed slowly (> 60s)\n  VALUE = {{ $value }}\n  LABELS: {{ $labels }}
```

#### Rabbitmq不可路由的消息

一个队列有不可更改的消息

```yaml
  - alert: RabbitmqUnroutableMessages
    expr: increase(rabbitmq_channel_messages_unroutable_returned_total[5m]) > 0 or increase(rabbitmq_channel_messages_unroutable_dropped_total[5m]) > 0
    for: 5m
    labels:
      severity: warning
    annotations:
      summary: Rabbitmq unroutable messages (instance {{ $labels.instance }})
      description: A queue has unroutable messages\n  VALUE = {{ $value }}\n  LABELS: {{ $labels }}  - alert: RabbitmqNoConsumer
    expr: rabbitmq_queue_consumers == 0
    for: 5m
    labels:
      severity: critical
    annotations:
      summary: Rabbitmq no consumer (instance {{ $labels.instance }})
      description: Queue has no consumer\n  VALUE = {{ $value }}\n  LABELS: {{ $labels }}
```
