---
title: Etcd 配置详解
---

# 概述

> 参考：
>
> - [官方文档](https://etcd.io/docs/current/op-guide/configuration/)

Etcd 运行时的行为可以通过三种方式进行配置

1. 配置文件
2. 命令行标志
3. 环境变量

而一般情况，配置文件中的关键字 与 命令行标志 和 环境变量 是 一一对应的。比如：

1. 配置文件中关键字：ETCD_DATA_DIR
2. 对应的环境变量中的变量名：ETCD_DATA_DIR
3. 对应的 flag： --data-dir

优先级：配置文件 > 命令行标志 > 环境变量

## Member 成员相关标志

--name # member 的名称。`默认值：default`
--data-dir # etcd 数据存储路径。`默认值：${name}.etcd`。一般大家都修改到 /var/lib/etcd 下。
--wal-dir
--snapshot-count
--heartbeat-interval # 心跳检测的间隔时间，时间单位是 milliseconds(毫秒)。`默认值：100`

- 注意：修改心跳值的同时要修改 election-timeout 标志。因为 选举超时 时间至少需要是 心跳检测间隔的 5 倍，如果达不到 5 倍，则 etcd 启动失败

--election-timeout # 选举超时时间，时间单位是 milliseconds(毫秒)。`默认值：1000`
--listen-peer-urls # 监听的用于节点之间通信的 url，可监听多个，集群内部将通过这些 url 进行数据交互(如选举，数据同步等)
--listen-client-urls # 监听的用于客户端通信的 url，同样可以监听多个。
--max-snapshots
--max-wals
--cors
--quota-backend-bytes # etcd 可储存的数据配额上限。`默认值：0`。

> 默认值 0 表示最低配额。在 3.4 版本时，最低配额是 2G，也就是说 etcd 最多可以保存 2G 的数据。

--backend-batch-limit
--backend-bbolt-freelist-type
--backend-batch-interval
--max-txn-ops
--max-request-bytes
--grpc-keepalive-min-time
--grpc-keepalive-interval
--grpc-keepalive-timeout

## Clustering 集群相关标志

> 注意：
> --initial-advertise-peer-urls，-initial-cluster，-initial-cluster-state 和 --initial-cluster-token 这 4 个标志是比较特殊的存在。只在 etcd 第一次启动并加入集群之前生效。
> 上面这 4 个标志用于引导（静态引导，服务发现引导 or 运行时配置）新成员，并且已经在集群中的成员重新启动时，将忽略这些标志。
> 使用发现服务时，需要设置 --discovery 前缀标志。

**--initial-advertise-peer-urls** # 用于节点间通信的 URL，节点间以该值进行通信。

- 默认值： <http://localhost:2380>

**--initial-cluster** # 用来引导初始集群的配置。一般是集群中所有 --initial-advertise-peer-urls 标志值的合集，每个值以逗号分隔

- 默认值：default=<http://localhost:2380>
- default 是每个节点的 etcd 的 --name 标志的值。--name 标志的默认值就是 default

**--initial-cluster-state** # 初始群集状态(两种状态：new 或 existing)。

- 默认值： new
- new # 对于在初始静态或 DNS 引导过程中存在的所有成员，将其设置为 new。
- existing # 设为 existing 状态的 etcd 将尝试加入 --initial-cluster 标志指定的集群。如果设置了错误的值，则 etcd 将尝试启动但安全失败。

**--initial-cluster-token** # 初始集群引导时所使用的 token。设置该值后集群将生成唯一 id，并为每个节点也生成唯一 id，当使用相同配置文件再启动一个集群时，只要该 token 值不一样，etcd 集群就不会相互影响。

- 默认值：etcd-cluster

**--advertise-client-urls** # 建议使用的客户端通信 url，该值用于 etcd 代理或 etcd 成员与 etcd 节点通信。

- 默认值：<http://localhost:2379>

--discovery #
--discovery-srv #
--discovery-srv-name #
--discovery-fallback #
--discovery-proxy #
--strict-reconfig-check #
--auto-compaction-retention #
--auto-compaction-mode #
--enable-v2 #

## Proxy 代理相关标志

--proxy

--proxy-failure-wait

--proxy-refresh-interval

--proxy-dial-timeout

--proxy-write-timeout

--proxy-read-timeout

## Security 安全相关标志

安全相关的标志用来帮助[建立一个安全的 etcd 集群](https://github.com/etcd-io/etcd/blob/master/Documentation/op-guide/security.md)

--ca-file

--cert-file

--key-file

--client-cert-auth

--client-crl-file

--client-cert-allowed-hostname

--trusted-ca-file

--auto-tls

--peer-ca-file

--peer-cert-file

--peer-key-file

--peer-client-cert-auth

--peer-crl-file

--peer-trusted-ca-file

--peer-auto-tls

--peer-cert-allowed-cn

--peer-cert-allowed-hostname

--cipher-suites

## Logging flags

--logger

--log-outputs

--log-level

--debug

--log-package-levels

## Unsafe flags

--force-new-cluster

## Miscellaneous flags

--version

--config-file

## Profiling flags

--enable-pprof

--metrics

--listen-metrics-urls

## Auth flags

--auth-token

--bcrypt-cost

## Experimental flags

-experimental-corrupt-check-time

--experimental-compaction-batch-limit

--experimental-peer-skip-client-san-verification
