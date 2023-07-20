---
title: "K3S 备份与恢复"
linkTitle: "K3S 备份与恢复"
date: "2023-07-20T16:29"
weight: 20
---

# 概述

> 参考：
> 
> - [公众号-边缘计算k3s社区，容器化应用的救命稻草：深入探索 K3s 备份和恢复](https://mp.weixin.qq.com/s/qpfKuvLvQ8E_pJ2WLpJm7g)


# 使用嵌入式 etcd 数据存储进行备份和恢复

## 创建快照

默认情况下，K3s 会在 00:00 和 12:00 自动创建快照，并保留 5 个快照。当然，你也可以禁用自动快照或者通过 `k3s etcd-snapshot save` 来手动创建快照。

快照目录默认为 `${data-dir}/server/db/snapshots`。`data-dir` 的默认值为 `/var/lib/rancher/k3s`，你可以通过设置 `--data-dir` 标志来更改。

## 从快照恢复集群

当 K3s 从备份中恢复时，旧的数据目录将被移动到 `${data-dir}/server/db/etcd-old/`。然后 K3s 将尝试通过创建一个新的数据目录来恢复快照，最后使用具有一个 etcd 成员的新 K3s 集群启动 etcd。

在此示例中有 3 个 K3s Server 节点，分别是 `S1`、`S2`和 `S3`，快照位于 `S1` 上:

1. 在 S1 上，使用 `--cluster-reset` 选项运行 K3s，同时指定 `--cluster-reset-restore-path`：

```
ls /var/lib/rancher/k3s/server/db/snapshots/
on-demand-ip-172-31-3-36-1688025329

systemctl stop k3s

k3s server \
  --cluster-reset \ 
  --cluster-reset-restore-path=/var/lib/rancher/k3s/server/db/snapshots/on-demand-ip-172-31-3-36-1688025329
```

2. 在 `S2` 和 `S3` 上，停止 K3s。然后删除数据目录 `/var/lib/rancher/k3s/server/db/`：

```
systemctl stop k3s
rm -rf /var/lib/rancher/k3s/server/db/
```

3. 在 `S1` 上，再次启动 K3s：

```
systemctl start k3s
```

4. `S1` 启动成功后，在 `S2` 和 `S3` 上，再次启动 K3s 以加入恢复后的集群：

```
systemctl start k3s
```

`S2` 和 `S3` 虽然使用空的数据目录来启动 K3s 服务，但启动时会自动到 `S1` 去同步数据，从而完成 `S2` 和 `S3` 的恢复。

另外，`k3s etcd-snapshot` 子命令支持 S3 兼容的 API，这样我们可以将快照自动或手动的上传到 S3 中存储，并且用于 K3s 的数据恢复。

# K3s 集群的灾难恢复

TODO