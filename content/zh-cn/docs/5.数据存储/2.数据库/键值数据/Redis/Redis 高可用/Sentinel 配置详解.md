---
title: Sentinel 配置详解
---

参考：[官方文档 1](https://redis.io/topics/sentinel#other-sentinel-options)、

Sentinel 的配置与 Redis 配置用法相同，当使用 --sentinel 参数启动 redis 时，则 redis 程序加载配置文件时，将只会特定的配置信息

# Sentinel 配置环境

## **port \<INT>** # Sentinel 监听的端口，Sentinel 之间使用该端口通讯

## sentinel monitor \<master-group-name> \<IP> \<PORT> \<QUORUM> # 指定 Sentinel 要监听的 MASTER

`sentinel montior mymaster 192.168.50.101 6379 1` 这个配置意味着，Sentinel 监控的目标 master 节点的 IP 为 192.168.50.101、端口为 6379，最后一个数字表示投票需要的"最少法定人数"。

> 最少法定人数的理解：比如有 10 个 sentinal 哨兵都在监控某一个 master 节点，如果需要至少 6 个哨兵发现 master 挂掉后，才认为 master 真正 down 掉，那么这里就配置为 6，最小配置 1 台 master，1 台 slave，在二个机器上都启动 sentinal 的情况下，哨兵数只有 2 个，如果一台机器物理挂掉，只剩一个 sentinal 能发现该问题，所以这里配置成 1。

至于 mymaster 只是一个名字，可以随便起，但要保证下面使用同一个名字

## sentinel down-after-milliseconds \<TARGET> \<DURATION> # 监控目标的 SDOWN 等待时长。单位：毫秒

持续 DURATION 时间 TARGET 没响应，就认为 SDOWN。

## sentinel parallel-syncs \<TARGET> \<INT> # 与监控目标并行同步数据的节点数

如果 master 重新选出来后，其它 replica 节点能同时并行从新 master 同步缓存的节点数有多少个。该值越大，所有 replica 节点完成同步切换的整体速度越快，但如果此时正好有人在访问这些 replica，可能造成读取失败，影响面会更广。最保定的设置为 1，只同一时间，只能有一台干这件事，这样其它 replica 还能继续服务，但是所有 replica 全部完成缓存更新同步的进程将变慢。

## sentinel failover-timeout mymaster \<DURATION> # 故障恢复超时时长。单位：毫秒

在指定时间 DURATION 后，故障恢复如果没有成功，则再次进行 Failover 操作

# 配置示例

    #
    sentinel monitor mymaster 172.19.42.231 6379 2
    sentinel down-after-milliseconds mymaster 60000
    sentinel failover-timeout mymaster 180000
    sentinel parallel-syncs mymaster 1
