---
title: Sentinel 命令
linkTitle: Sentinel 命令
weight: 20
---

# 概述

> 参考：
>
> - [官方文档](https://redis.io/topics/sentinel#sentinel-api)
>   - https://redis.io/docs/latest/operate/oss_and_stack/management/sentinel/#sentinel-api
> - [博客园大佬](https://www.cnblogs.com/biglittleant/p/7770960.html)

## 基本命令

### PING - 健康检查

### sentinel masters - 显示被监控的所有 master 以及它们的状态

### sentinel master \<MASTER> - 显示指定 MASTER 的信息和状态

### sentinel slaves \<MASTER> - 显示指定 MASTER 的所有 slave 以及它们的状态

### sentinel get-master-addr-by-name \<MASTER> - 返回指定 MASTER 的 ip 和端口

如果正在进行 failover 或者 failover 已经完成，将会显示被提升为 master 的 slave 的 ip 和端口。

### sentinel reset \<pattern> - 重置名字匹配该正则表达式的所有的 master 的状态信息

清楚其之前的状态信息，以及 slaves 信息。

### sentinel failover \<MASTER> - 强制 sentinel 为指定的 MASTER 执行 failover

执行该操作的 Sentinel 节点并不需要得到其他 Sentinel 的同意。但是 failover 后会将最新的配置发送给其他 sentinel。

## 动态修改 Sentinel 配置

注意：如果你通过 API 修改了一个 sentinel 的配置，sentinel 不会把修改的配置告诉其他 sentinel。你需要自己手动地对多个 sentinel 发送修改配置的命令。

### sentinel flushconfig - 重写运行 Sentinel 时指定的配置文件内容

与 Redis 的 config rewrite 命令效果一样。也就是将内存中的配置，写入到文件中。

### SENTINEL MONITOR \<MASTER> \<IP> \<PORT> \<QUORUM> - 让 sentinel 监听指定 MASTER

### sentinel remove \<MASTER> - 让 sentinel 放弃对指定 MASTER 的监听

### SENTINEL SET \<name> \<option> \<value> - 这个命令很像 Redis 的 CONFIG SET 命令，用来改变指定 master 的配置。支持多个\<option>\<value>

只要是配置文件中存在的配置项，都可以用`SENTINEL SET`命令来设置。

EXAMPLE

- 设置 master 的属性，比如说**quorum**(票数)，而不需要先删除 master，再重新添加 master。

  - sentinel set objects-cache-master quorum 5

## 增加或删除 Sentinel

由于有 sentinel 自动发现机制，所以添加一个 sentinel 到你的集群中非常容易，你所需要做的只是监控到某个 Master 上，然后新添加的 sentinel 就能获得其他 sentinel 的信息以及 master 所有的 slaves。

如果你需要添加多个 sentinel，建议你一个接着一个添加，这样可以预防网络隔离带来的问题。你可以每个 30 秒添加一个 sentinel。最后你可以用`SENTINEL MASTER mastername`来检查一下是否所有的 sentinel 都已经监控到了 master。删除一个 sentinel 显得有点复杂：因为 sentinel 永远不会删除一个已经存在过的 sentinel，即使它已经与组织失去联系很久了。

要想删除一个 sentinel，应该遵循如下步骤：

1. 停止所要删除的 sentinel

2. 发送一个`SENTINEL RESET *`命令给所有其它的 sentinel 实例，如果你想要重置指定 master 上面的 sentinel，只需要把\*号改为特定的名字，注意，需要一个接一个发，每次发送的间隔不低于 30 秒。

3. 检查一下所有的 sentinels 是否都有一致的当前 sentinel 数。使用`SENTINEL MASTER mastername` 来查询。

## 删除旧 master 或者不可达 slave

sentinel 永远会记录好一个 Master 的 slaves，即使 slave 已经与组织失联好久了。这是很有用的，因为 sentinel 集群必须有能力把一个恢复可用的 slave 进行重新配置。

并且，failover 后，失效的 master 将会被标记为新 master 的一个 slave，这样的话，当它变得可用时，就会从新 master 上复制数据。

然后，有时候你想要永久地删除掉一个 slave(有可能它曾经是个 master)，你只需要发送一个`SENTINEL RESET master`命令给所有的 sentinels，它们将会更新列表里能够正确地复制 master 数据的 slave。

## 发布/订阅

客户端可以向一个 sentinel 发送订阅某个频道的事件的命令，当有特定的事件发生时，sentinel 会通知所有订阅的客户端。需要注意的是客户端只能订阅，不能发布。

订阅频道的名字与事件的名字一致。例如，频道名为 sdown 将会发布所有与 SDOWN 相关的消息给订阅者。

如果想要订阅所有消息，只需简单地使用`PSUBSCRIBE *`

以下是所有你可以收到的消息的消息格式，如果你订阅了所有消息的话。第一个单词是频道的名字，其它是数据的格式。

注意：以下的 instance details 的格式是：

`<instance-type> <name> <ip> <port> @ <master-name> <master-ip> <master-port>`

如果这个 redis 实例是一个 master，那么@之后的消息就不会显示。

```bash
+reset-master <instance details> -- 当master被重置时.
  +slave <instance details> -- 当检测到一个slave并添加进slave列表时.
  +failover-state-reconf-slaves <instance details> -- Failover状态变为reconf-slaves状态时
  +failover-detected <instance details> -- 当failover发生时
  +slave-reconf-sent <instance details> -- sentinel发送SLAVEOF命令把它重新配置时
  +slave-reconf-inprog <instance details> -- slave被重新配置为另外一个master的slave，但数据复制还未发生时。
  +slave-reconf-done <instance details> -- slave被重新配置为另外一个master的slave并且数据复制已经与master同步时。
  -dup-sentinel <instance details> -- 删除指定master上的冗余sentinel时 (当一个sentinel重新启动时，可能会发生这个事件).
  +sentinel <instance details> -- 当master增加了一个sentinel时。
  +sdown <instance details> -- 进入SDOWN状态时;
  -sdown <instance details> -- 离开SDOWN状态时。
  +odown <instance details> -- 进入ODOWN状态时。
  -odown <instance details> -- 离开ODOWN状态时。
  +new-epoch <instance details> -- 当前配置版本被更新时。
  +try-failover <instance details> -- 达到failover条件，正等待其他sentinel的选举。
  +elected-leader <instance details> -- 被选举为去执行failover的时候。
  +failover-state-select-slave <instance details> -- 开始要选择一个slave当选新master时。
  no-good-slave <instance details> -- 没有合适的slave来担当新master
  selected-slave <instance details> -- 找到了一个适合的slave来担当新master
  failover-state-send-slaveof-noone <instance details> -- 当把选择为新master的slave的身份进行切换的时候。
  failover-end-for-timeout <instance details> -- failover由于超时而失败时。
  failover-end <instance details> -- failover成功完成时。
  switch-master <master name> <oldip> <oldport> <newip> <newport> -- 当master的地址发生变化时。通常这是客户端最感兴趣的消息了。
  +tilt -- 进入Tilt模式。
  -tilt -- 退出Tilt模式。
```

TILT 模式

redis sentinel 非常依赖系统时间，例如它会使用系统时间来判断一个 PING 回复用了多久的时间。

然而，假如系统时间被修改了，或者是系统十分繁忙，或者是进程堵塞了，sentinel 可能会出现运行不正常的情况。

当系统的稳定性下降时，TILT 模式是 sentinel 可以进入的一种的保护模式。当进入 TILT 模式时，sentinel 会继续监控工作，但是它不会有任何其他动作，它也不会去回应`is-master-down-by-addr`这样的命令了，因为它在 TILT 模式下，检测失效节点的能力已经变得让人不可信任了。

如果系统恢复正常，持续 30 秒钟，sentinel 就会退出 TITL 模式。

## -BUSY 状态

> 注意：该功能还未实现。

当一个脚本的运行时间超过配置的运行时间时，sentinel 会返回一个-BUSY 错误信号。如果这件事发生在触发一个 failover 之前，sentinel 将会发送一个 SCRIPT KILL 命令，如果 script 是只读的话，就能成功执行。
