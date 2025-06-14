---
title: Etcd 性能测试
---

# 概述

> 参考：
>
> - [官方文档，运维指南-性能](https://etcd.io/docs/current/op-guide/performance/)

安装 etcd 压测工具 benchmark：

```bash
$ go get go.etcd.io/etcd/tools/benchmark
# GOPATH should be set
$ ls $GOPATH/bin
benchmark
```

# 官方推荐的 etcd 性能数据

其中官方使用的设备信息为：

- Google Cloud Compute Engine
- 3 machines of 8 vCPUs + 16GB Memory + 50GB SSD
- 1 machine(client) of 16 vCPUs + 30GB Memory + 50GB SSD
- Ubuntu 17.04
- etcd 3.2.0, go 1.8.3

### etcd 写性能

| Key 数量  | 每个 Key 的大小 | 每个值的大小    | 连接数量 | 客户端数量 | 目标 etcd 节点数 | 写性能的平均 QPS | 每个请求的平均延迟 | 服务器 RRS 的平均值 |
| ------- | ---------- | --------- | ---- | ----- | ----------- | ---------- | --------- | ------------ |
| 10,000  | 8 bytes    | 256 bytes | 1    | 1     | 只有一个 leader | 583        | 1.6ms     | 48 MB        |
| 100,000 | 8 bytes    | 256 bytes | 100  | 1000  | 只有一个 leader | 44,341     | 22ms      | 124MB        |
| 100,000 | 8 bytes    | 256 bytes | 100  | 1000  | 所有 etcd 节点  | 50,104     | 20ms      | 126MB        |

### etcd 读性能

该测试的 --endpoint 参数指定所有节点。

Linearizable(线性化) 读请求经过集群仲裁达成共识以获取最新数据，Serializable(串行化)读取请求比线性化读取要廉价一些，因为他们是通过任意一台 etcd 成员来相应请求，而不是具有法定人数的成员，这种请求获取到的数据有可能是过期的。由于 etcd 是强一致性的，其默认读取测试就是线性化读取。

| 请求数     | 每个 Key 的大小 | 每个值的大小    | 连接数量 | 客户端数量 | Consistency(一致性) | 读性能的平均 QPS | 每个请求的平均延迟 |
| ------- | ---------- | --------- | ---- | ----- | ---------------- | ---------- | --------- |
| 10,000  | 8 bytes    | 256 bytes | 1    | 1     | Linearizable     | 1,353      | 0.7ms     |
| 10,000  | 8 bytes    | 256 bytes | 1    | 1     | Serializable     | 2,909      | 0.3ms     |
| 100,000 | 8 bytes    | 256 bytes | 100  | 1000  | Linearizable     | 141,578    | 5.5ms     |
| 100,000 | 8 bytes    | 256 bytes | 100  | 1000  | Serializable     | 185,758    | 2.2ms     |

官方鼓励在新环境中首次设置 etcd 集群时运行基准测试，以确保该集群获得足够的性能；群集延迟和吞吐量可能会对较小的环境差异敏感。

# etcdctl check perf 官方测试工具

直接使用 etcdctl 命令行工具的子命令，即可进行简单的测试

```bash
~# etcdctl check perf
60 / 60 Booooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooo! 100.00% 1m0s
PASS: Throughput is 150 writes/s
PASS: Slowest request took 0.101178s
PASS: Stddev is 0.002695s
PASS
```

# benchmarks 官方测试工具

对于写入测试，按照官方文档的测试方法指定不同数量的客户端和连接数以及 key 的大小，对于读取操作，分别测试了线性化读取以及串行化读取，由于 etcd 是强一致性的，其默认读取测试就是线性化读取。

Note:

1. 下面的测试环境为 etcd 3.4.3 版本，k8s 集群已部署完成。测试命令就是根据上文官方的测试结果所使用的命令，一共测试 7 次
2. 先后创建 benchmark 别名，在别名中加指定证书和 endpoint。
   1. echo 'alias benchmark="benchmark --key=/etc/kubernetes/pki/etcd/peer.key --cert=/etc/kubernetes/pki/etcd/peer.crt --cacert=/etc/kubernetes/pki/etcd/ca.crt --endpoints=172.40.0.3:2379,172.40.0.4:2379,172.40.0.5:2379"' >> /etc/bashrc
3. 测试结果中，关于延迟数据的单位是 secs，是秒的意思，1 秒=1000 毫秒(ms)

### 写入测试

```bash
# write to leader
benchmark --endpoints=${HOST_1} --target-leader --conns=1 --clients=1 \
    put --key-size=8 --sequential-keys --total=10000 --val-size=256
benchmark --endpoints=${HOST_1} --target-leader  --conns=100 --clients=1000 \
    put --key-size=8 --sequential-keys --total=100000 --val-size=256

# write to all members
benchmark --endpoints=${HOST_1},${HOST_2},${HOST_3} --conns=100 --clients=1000 \
    put --key-size=8 --sequential-keys --total=100000 --val-size=256
```

### 读取测试

```bash
# Single connection read requests
benchmark --endpoints=${HOST_1},${HOST_2},${HOST_3} --conns=1 --clients=1 \
    range YOUR_KEY --consistency=l --total=10000
benchmark --endpoints=${HOST_1},${HOST_2},${HOST_3} --conns=1 --clients=1 \
    range YOUR_KEY --consistency=s --total=10000

# Many concurrent read requests
benchmark --endpoints=${HOST_1},${HOST_2},${HOST_3} --conns=100 --clients=1000 \
    range YOUR_KEY --consistency=l --total=100000
benchmark --endpoints=${HOST_1},${HOST_2},${HOST_3} --conns=100 --clients=1000 \
    range YOUR_KEY --consistency=s --total=100000
```

```bash
$ benchmark --conns=1 --clients=1  \
range foo --consistency=l --total=10000

$ benchmark --conns=1 --clients=1  \
range foo --consistency=s --total=10000

$ benchmark --conns=100 --clients=1000  \
range foo --consistency=l --total=100000

$ benchmark --conns=100 --clients=1000  \
range foo --consistency=s --total=100000
```

# 第三方 fio 工具测试

参考：<https://www.chenshaowen.com/blog/the-use-of-etcd-and-etcdctl.html>

[https://www.ibm.com/cloud/blog/using-fio-to-tell-whether-your-storage-is-fast-enough-for-etcd?mhsrc=ibmsearch_a\&mhq=fio%20etcd](https://www.ibm.com/cloud/blog/using-fio-to-tell-whether-your-storage-is-fast-enough-for-etcd?mhsrc=ibmsearch_a&mhq=fio%2520etcd)

<https://github.com/etcd-io/etcd/issues/10577>

etcd 的一些 Prometheus 指标。其中之一是 wal_fsync_duration_seconds。 etcd docs 建议该度量标准的第 99 个百分位数应小于 10 毫秒。如果您正在考虑 etcd 在 Linux 机器上运行 集群，并且需要评估存储（例如 SSD）是否足够快，那么一种选择是使用流行的 I/O 测试工具 fio。为此，您可以运行以下命令(其中 /var/lib/etcd 是要测试的存储设备的安装点下的目录,根据 etcd 存储路径进行修改)

fio --rw=write --ioengine=sync --fdatasync=1 --directory=/var/lib/etcd --size=22m --bs=2300 --name=mytest

然后，您要做的就是查看输出并检查 fdatasync 持续时间的第 99 个百分位数 是否小于 10ms。如果真是这样，则您的存储空间足够快。这是示例输出：

```bash
fsync/fdatasync/sync_file_range:
  sync (usec): min=534, max=15766, avg=1273.08, stdev=1084.70
  sync percentiles (usec):
   | 1.00th=[ 553], 5.00th=[ 578], 10.00th=[ 594], 20.00th=[ 627],
   | 30.00th=[ 709], 40.00th=[ 750], 50.00th=[ 783], 60.00th=[ 1549],
   | 70.00th=[ 1729], 80.00th=[ 1991], 90.00th=[ 2180], 95.00th=[ 2278],
   | 99.00th=[ 2376], 99.50th=[ 9634], 99.90th=[15795], 99.95th=[15795],
   | 99.99th=[15795]
```

一些注意事项：

- 在上面的示例中，我们针对特定情况调整了--size 和 --bs 参数的值 。为了从中获得有意义的见解 fio，您应该使用最适合您的案例的值。要学习如何派生它们，请阅读我们 如何发现如何配置 fio 的信息。
- 在测试过程中，fio 生成的 I / O 负载 是唯一的 I / O 活动。在实际情况下，除与 wal_fsync_duration_seconds 相关的写入外，可能还会对存储进行其他写入。这样的额外负载会使 wal_fsync_duration_seconds 更大。因此，如果您观察到的第 99 个百分位数 fio 仅略低于 10 毫秒，则可能是存储速度不够快。
- 您需要的 Fio 版本至少应为 3.5， 因为较旧的版本不会报告 fdatasync 持续时间百分位数。
- 上面的输出只是来自的整个输出的一小部分摘录 fio。

```bash
fio --randrepeat=1 \
  --ioengine=libaio \
  --direct=1 \
   --gtod_reduce=1 \
   --name=etcd-disk-io-test \
   --filename=etcd_read_write.io \
   --bs=4k --iodepth=64 --size=4G \
   --readwrite=randrw --rwmixread=75
```

# etcd 官方推荐的硬件配置

官方文档：<https://github.com/etcd-io/etcd/blob/master/Documentation/op-guide/hardware.md>

etcd通常在开发或测试的时候用很少的资源就可以了，比如说使用普通的笔记本或者是廉价的云主机就可以，但是在生产环境上，还是需要按推荐的硬件配置进行部署，虽然这不是必须的，但是这样做可以增加集群的健壮性。一如既往，在上生产环境之前，需要先进行负载模拟测试。

CPUs

很少有etcd部署需要大量的CPU资源。典型的etcd部署节点，需要2-4个CPU就可以顺利运行。负载很高的etcd集群，比如上千客户端或者每秒超过上万请求，倾向于CPU绑定，可以直接从内存获取请求。即便这样重的负载，通常需要8-16个CPU就可以了。

内存

etcd占用的内存相对比较小，但是etcd性能仍然取决于是否拥有足够的内存。一个etcd服务器会积极的缓存key-value数据和大部分的跟踪watcher。通常8GB内存就够了，对于负载高的集群，比如有上千watcher和超过百万的keys，相应的就需要16GB-64GB内存。

磁盘

高速磁盘是保证 etcd 部署性能和稳定性的关键因素。

缓慢的磁盘会增加 etcd 请求的延迟，并有可能损害群集的稳定性。由于 etcd 的共识协议依赖于将元数据持久存储在日志中，因此大多数 etcd 集群成员必须将每个请求写入磁盘。另外，etcd 还将以增量方式将其状态检查到磁盘，以便可以截断该日志。如果这些写入花费的时间太长，则心跳可能会超时并触发选举，从而破坏了群集的稳定性。通常，要判断磁盘是否足够快用于 etcd，可以使用诸如 fio 之类的基准测试工具。阅读此处的示例

etcd对磁盘写入延时非常敏感。通常稳定达到50 IOPS（比如：一个7200转的磁盘）是必须的，对于负载很高的集群，推荐能稳定达到500 IOPS（比如：一个典型的本地SSD盘或者高性能的虚拟块设备盘）。注意，大多数云服务提供商发布的是瞬时并发IOPS，并不是稳定的IOPS，瞬时并发IOPS可能十倍于稳定连续的IOPS（说明：因为瞬时并发IOPS可能会写缓存，或者测试时无其他用户竞争磁盘资源，所以会很高，当测试时间很长后，就会测试出设备的真实IOPS能力，这个在国内云厂商基本没有这个问题）。测试稳定连续IOPS，我们建议使用磁盘基准测试工具，比如 diskbench 或者 fio。

etcd对磁盘带宽没什么要求，但是更大的磁盘带宽可以在失败节点加入集群时，更快的完成恢复操作。通常10MB/s带宽的磁盘15s可以恢复100MB的数据，对于大型集群，100MB/s或更高带宽的磁盘可以在15s内恢复1GB数据。

如果有可能，etcd后端存储就用SSD。一个SSD磁盘和机械盘相比，通常会提供更低的写入延时和更少的数据跳变（variance），因此可以提高etcd集群的稳定性和可靠性。如果使用机械盘，尽可能使用最快的（15000转）。使用RAID 0也是一种有效提高磁盘性能的方法，不管是机械盘还是SSD都可以。etcd集群至少有3个节点，磁盘使用RAID做镜像或者做奇偶校验都是不必要的，因为etcd自身的一致性复制已经保证了数据的高可用。

网络

多节点部署的etcd集群会受益于快速和可靠的网络。为了满足etcd集群的一致性和分区容忍，一个不可靠网络出现网络分区会导致部分节点无效。低延时可以保证etcd成员之间快速通信，高带宽可以减少etcd故障节点的恢复时间。1Gb网络就可以满足常见的etcd部署场景，对于大型etcd集群，使用10Gb网络可以减少平均故障恢复时间。

如果有可能，尽量将所有etcd成员节点部署在同一个数据中心，这样可以避免网络延时开销，降低发生网络分区的可能性。如果需要另外的数据中心级故障域，尽量选择和当前数据中心离得比较近的。也可以阅读性能调优文档，了解跨数据中心部署的更多信息。

## 示例硬件配置

这有一些在AWS和GCE环境上的硬件配置例子。如上所述，但是还是有必要再强调一下，无论如何，管理员在将etcd集群投入生产环境使用之前，都应该做一下模拟负载测试。

请注意：这些配置假设这些服务器只用来跑etcd服务。如果在这些服务器上还跑其他服务，可能会导致其他服务和etcd抢资源，存在资源竞争问题，导致etcd集群不稳定。

**小型集群**

一个小型集群服务少于100个客户端，访问请求小于每秒200，并且存储不超过100MB的数据。

示例应用负载：一个 50 节点的 kubernetes 集群

| 提供商 | 类型                          | vCPUs | 内存 (GB) | 最大并发 IOPS | 磁盘带宽 (MB/s) |
| --- | --------------------------- | ----- | ------- | --------- | ----------- |
| AWS | m4.large                    | 2     | 8       | 3600      | 56.25       |
| GCE | n1-standard-2 + 50GB PD SSD | 2     | 7.5     | 1500      | 25          |

**中型集群**

一个中型集群服务少于500个客户端，访问请求小于每秒1000，并且存储不超过500MB的数据。

示例应用负载：一个 250 节点的 kubernetes 集群

| 提供商 | 类型                           | vCPUs | 内存 (GB) | 最大并发 IOPS | 磁盘带宽 (MB/s) |
| --- | ---------------------------- | ----- | ------- | --------- | ----------- |
| AWS | m4.xlarge                    | 4     | 16      | 6000      | 93.75       |
| GCE | n1-standard-4 + 150GB PD SSD | 4     | 15      | 4500      | 75          |

**大型集群**

    一个大型集群服务少于1500个客户端，访问请求小于每秒10000，并且存储不超过1GB的数据。

示例应用负载：一个 1000 节点的 kubernetes 集群

| 提供商 | 类型                           | vCPUs | 内存 (GB) | 最大并发 IOPS | 磁盘带宽 (MB/s) |
| --- | ---------------------------- | ----- | ------- | --------- | ----------- |
| AWS | m4.2xlarge                   | 8     | 32      | 8000      | 125         |
| GCE | n1-standard-8 + 250GB PD SSD | 8     | 30      | 7500      | 125         |

**超大型集群**

一个超大型集群服务超过1500个客户端，访问请求超过每秒10000，并且存储超过1GB的数据。

示例应用负载：一个 3000 节点的 kubernetes 集群

| 提供商 | 类型                            | vCPUs | 内存 (GB) | 最大并发 IOPS | 磁盘带宽 (MB/s) |
| --- | ----------------------------- | ----- | ------- | --------- | ----------- |
| AWS | m4.4xlarge                    | 16    | 64      | 16,000    | 250         |
| GCE | n1-standard-16 + 500GB PD SSD | 16    | 60      | 15,000    | 250         |
