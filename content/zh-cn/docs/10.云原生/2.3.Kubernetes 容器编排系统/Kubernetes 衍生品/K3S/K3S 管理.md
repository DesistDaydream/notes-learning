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


sudo etcdctl version \
  --cacert=/var/lib/rancher/k3s/server/tls/etcd/server-ca.crt \
  --cert=/var/lib/rancher/k3s/server/tls/etcd/client.crt \
  --key=/var/lib/rancher/k3s/server/tls/etcd/client.key

# 其他推荐

## 关于旧版 iptables

## 建议关闭 firewalld