---
title: "K3S 管理"
linkTitle: "K3S 管理"
weight: 20
---

# 概述

> 参考：
> 
> - [官方文档，高级选项和配置](https://docs.k3s.io/zh/advanced)

# 使用 etcdctl

```bash
sudo etcdctl \
  --cacert=/var/lib/rancher/k3s/server/tls/etcd/server-ca.crt \
  --cert=/var/lib/rancher/k3s/server/tls/etcd/client.crt \
  --key=/var/lib/rancher/k3s/server/tls/etcd/client.key \
  --endpoint=A,B,C \
  member list
```

# 更换 Master 节点

K3S 可以快速更换所有 Master 节点的系统，只需要剔除一个新增一个，逐一替换即可。就算没有 `k3s server --cluster-init` 这种命令的节点也是可以的。

# 其他推荐

## 关于旧版 iptables

## 建议关闭 firewalld