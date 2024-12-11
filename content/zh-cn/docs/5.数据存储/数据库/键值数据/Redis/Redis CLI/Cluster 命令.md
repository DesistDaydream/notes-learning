---
title: Cluster 命令
linkTitle: Cluster 命令
date: 2024-06-14T15:28
weight: 20
---

# 概述

> 参考：
>
> - https://redis.io/docs/latest/commands/?group=cluster


一、以下命令是Redis Cluster集群所独有的，执行下面命令需要先登录redis：

\[root@manage redis]# redis-cli -c -p 6382 -h 192.168.10.12     （客户端命令：redis-cli -c -p port -h ip）

192.168.10.12:6382>  登录redis后，在里面可以进行下面命令操作

集群

cluster info ：打印集群的信息

cluster nodes ：列出集群当前已知的所有节点（ node），以及这些节点的相关信息。

节点

cluster meet   ：将 ip 和 port 所指定的节点添加到集群当中，让它成为集群的一份子。

cluster forget  ：从集群中移除 node\_id 指定的节点。

cluster replicate  ：将当前从节点设置为 node\_id 指定的master节点的slave节点。只能针对slave节点操作。

cluster saveconfig ：将节点的配置文件保存到硬盘里面。

槽(slot)

cluster addslots  \[slot ...] ：将一个或多个槽（ slot）指派（ assign）给当前节点。

cluster delslots  \[slot ...] ：移除一个或多个槽对当前节点的指派。

cluster flushslots ：移除指派给当前节点的所有槽，让当前节点变成一个没有指派任何槽的节点。

cluster setslot  node  ：将槽 slot 指派给 node\_id 指定的节点，如果槽已经指派给

另一个节点，那么先让另一个节点删除该槽>，然后再进行指派。

cluster setslot  migrating  ：将本节点的槽 slot 迁移到 node\_id 指定的节点中。

cluster setslot  importing  ：从 node\_id 指定的节点中导入槽 slot 到本节点。

cluster setslot  stable ：取消对槽 slot 的导入（ import）或者迁移（ migrate）。

键

cluster keyslot  ：计算键 key 应该被放置在哪个槽上。

cluster countkeysinslot  ：返回槽 slot 目前包含的键值对数量。

cluster getkeysinslot   ：返回 count 个 slot 槽中的键 。
