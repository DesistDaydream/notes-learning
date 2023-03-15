---
title: etcdctl 命令行工具
---

# 概述

> 参考：
>
> - 官方文档：<https://github.com/etcd-io/etcd/tree/master/etcdctl>
> - [etcd 可用的库和客户端](https://etcd.io/docs/latest/integrations/)

# etcdctl \[GlobalOptions] COMMAND \[CommandOptions] \[Arguments...]

使用说明：

1. export ETCDCTL_API=3 使用该命令使得 etcdctl 通过 v3 版本来进行操作
2. 如果在 etcd 的配置文件中的 Security 段落，开启了验证证书，则在使用命令时，需要使用--cert、--key、--cacert 选项来指定验证所需证书，否则无法操纵服务器
   1. v2 版本中使用如下方式 etcdctl --key-file=/etc/kubernetes/pki/etcd/peer.key --cert-file=/etc/kubernetes/pki/etcd/peer.crt --ca-file=/etc/kubernetes/pki/etcd/ca.crt --endpoints="https://IP:PORT" COMMAND
   2. v3 版本中使用如下方式 etcdctl --key=/etc/kubernetes/pki/etcd/peer.key --cert=/etc/kubernetes/pki/etcd/peer.crt --cacert=/etc/kubernetes/pki/etcd/ca.crt --endpoints="https://IP:PORT" COMMAND
   3. 在下面的 EXAMPLE 则不再输入认证相关参数，以便查阅方便。但是实际使用中需要使用，否则无法连接 etcd 服务端

## GLOBAL OPTIONS

- **--cacert=/PATH/FILE** # 使用此 CA 包验证启用 TLS 的安全服务器的证书。即 etcd 的 ca，用该 ca 来验证 cert 选项中提供的证书是否正确
- **--cert=/PATH/FILE**# 使用指定的 TLS 证书文件鉴定客户端是否安全。即 etcd 的 peer 证书，peer 证书对于 etcdctl 来说就是与它交互的服务端的证书
- **--key=/PATH/FILE** # 使用指定的 TLS 证书的密钥文件鉴定客户端是否安全。即 etcd 的 peer 证书的私钥
- **--endpoints=\[IP1:PORT1,IP2:PORT2,.....]** # 指定后端服务器的 IP 和 Port
- --command-timeout=5s # timeout for short running command (excluding dial timeout)
- --debug\[=false] # enable client-side debug logging
- --dial-timeout=2s # dial timeout for client connections
- --hex\[=false] # print byte strings as hex encoded strings
- --insecure-skip-tls-verify\[=false # skip server certificate verification
- --insecure-transport\[=true] # disable transport security for client connections
- --user="" # username\[:password] for authentication (prompt if password is not supplied)
- -w, --write-out="simple" # 指定输出内容的格式，格式可有有这么几个 (fields, json, protobuf, simple, table)(一般常用 json)
  - Note：输出的 json 格式只有一行，可以使用 jq 程序来对 json 进行格式化，可以把每个{}分行，以便人类阅读，下图为样例

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/miobxe/1616136392283-c0b50823-df6d-49d3-8d85-2aed1c7de3e0.jpeg)

# 子命令详解

# Base 基本

## get # 获取键或者键的范围。Gets the key or a range of keys

**etcdctl get \[OPTIONS] \[Range-End]**
OPTIONS

- --consistency="l" Linearizable(l) or Serializable(s)
- --from-key\[=false] Get keys that are greater than or equal to the given key using byte compare
- **--keys-only\[=false]** # 仅获取键而不显示该键所对应的值
- --limit=0 # Maximum number of results
- --order="" # Order of results; ASCEND or DESCEND (ASCEND by default)
- **--prefix\[=false]** # 获取 KEY 前缀匹配到的所有的键。前缀就是键名的从开头开始的的字符串。可以指定`''`匹配所有 KEY
- --print-value-only\[=false] Only write values when using the "simple" output format
- --rev=0 Specify the kv revision
- --sort-by="" Sort target; CREATE, KEY, MODIFY, VALUE, or VERSION

EXAMPLE

- etcdctl get --prefix --keys-only '' # 获取所有键，并且只显示键名。
- etcdctl get --prefix --keys-only / # 获取以`/`开头的所有键，并且只显示键名。
- etcdv3 get /registry/events/kube-system/kube-flannel-ds-amd64-47cnw.15966b23d2027e45 -w=json | jq . # 以 json 格式输出指定键的值，并使用 jq 命令对 json 内容进行格式化输出以便人类阅读

put # 写入一个 key/value 到 etcd 存储中。

del # 删除指定的 key 或一个范围的 keys

txn Txn processes all the requests in one transaction

compaction Compacts the event history in etcd

# alarm # etcd 中告警相关命令

## alerm disarm # 解除所有告警

## alarm list # 列出 etcd 中所有的告警

# check # 检查 etcd 的性能

## check datascale # Check the memory usage of holding data for different workloads on a given server endpoint

## check perf # 检查 etcd 的性能

检查 60 秒的 etcd 群集性能。经常运行检查性能可以创建一个较大的键空间历史记录，可以使用--auto-compact 和--auto-defrag 选项（如下所述）对其进行自动压缩和碎片整理。

OPTIONS

1. \--load # 性能检查的工作负载模型。可接受的工作负载：s(small 小)，m(medium 中)，l(large 大)，xl(x 大)

# defrag # 对指定 endpoints 的 etcd 成员的存储空间进行碎片整理

# endpoint # 用于查询 etcd 中各个端点的信息

endpoint health # Checks the healthiness of endpoints specified in `--endpoints` flag

endpoint status # 打印出 --endpoints 标志中指定的 endpoints 状态

endpoint hashkv # Prints the KV history hash for each endpoint in --endpoints

# lease 相关命令

lease grant Creates leases

lease revoke Revokes leases

lease timetolive Get lease information

lease keep-alive Keeps leases alive (renew)

# member # 用于管理 etcd 集群中的成员

member add # 将新成员作为新对等方引入 etcd 集群中。

member remove # 从参与集群共识的成员中删除 etcd 集群的成员。

member update # 为 etcd 集群中现有成员设置对等 URL。

member list # 列出集群中的所有成员

EXAMPLE

1. etcdctl member list # 列出 etcd 集群中的成员
2. etcdctl member list --write-out=json | jq . # 通过 json 可以看到 etcd 集群中，哪个是节点是 leader

# snapshot 快照相关命令。用来让用户对 etcd 的数据进行备份与恢复

**etcdctl snapshot save** # 存储一个 etcd 节点后端快照到指定文件。Stores an etcd node backend snapshot to a given file
EXAMPLE

- etcdctl snapshot save snapshot.db # 备份指定后端节点的 etcd 数据到 snapshot.db 文件

**etcdctl snapshot restore \[OPTIONS]** # 恢复一个 etcd 成员的快照到一个 etcd 的文件夹中。
OPTIONS

- **--data-dir=/PATH/FILE**# 把指定的路径作为 snapshot 文件的恢复目录，会把数据写到指定的目录下。Path to the data directory

EXAMPLE

- etcdctl snapshot restore snapshot.db --data-dir=/var/lib/etcd/ # 从 snapshot.db 文件中恢复数据到/var/lib/etcd/目录下

**etcdctl snapshot status** # 取指定文件的后端快照状态。Gets backend snapshot status of a given file
EXAMPLE

- etcdv3 snapshot status snapshot.db # 获取 snapshot.db 的状态，包括键/值对有多少，占多少空间

# 其他

make-mirror Makes a mirror at the destination etcd cluster

migrate Migrates keys in a v2 store to a mvcc store

lock Acquires a named lock

elect Observes and participates in leader election

auth enable Enables authentication

auth disable Disables authentication

## user # etcd 用户相关的命令

详见：[etcdctl user 命令](https://www.yuque.com/go/doc/33190532)

## role # etcd 的 role 相关的命令

role add # Adds a new role

role delete # Deletes a role

role get # Gets detailed information of a role

role list # Lists all roles

role grant-permission # Grants a key to a role

role revoke-permission # Revokes a key from a role

watch Watches events stream on keys or prefixes
