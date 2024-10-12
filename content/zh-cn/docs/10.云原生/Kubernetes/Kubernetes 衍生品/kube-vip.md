---
title: kube-vip
linkTitle: kube-vip
date: 2022-10-12T15:55:00
weight: 20
---

# 概述

> 参考：
>
> - [GitHub 项目，kube-vip/kube-vip](https://github.com/kube-vip/kube-vip)
> - [官网](https://kube-vip.io/)
> - [官方文档，使用静态 Pod 部署 kube-vip](https://kube-vip.io/hybrid/static/)

kube-vip 是一个用于 Kubernetes 控制平面的 VIP 和 负载均衡器。他可以实现 Keepalived + HAProxy 的功能。

Kube-VIP 最初是为 Kubernetes 控制平面提供 HA 解决方案而创建的，随着时间的推移，它已经发展为将相同的功能合并到 Kubernetes 的 load-banlancers 类型的 Service 资源。

# 配置

Kube-VIP 通过命令行标志变更运行时行为

# 部署

## 作为静态 Pod 运行

生成 Manifests

```bash
export VIP=172.38.180.213
export INTERFACE=eth0
export KVVERSION=$(curl -sL https://api.github.com/repos/kube-vip/kube-vip/releases | jq -r ".[0].name")

docker run --rm --net host docker.io/plndr/kube-vip:${KVVERSION} \
/kube-vip manifest pod \
--interface ${INTERFACE} \
--vip ${VIP} \
--controlplane \
--services \
--arp \
--leaderElection | tee  /etc/kubernetes/manifests/kube-vip.yaml
```

等待 kubelet 将 pod 启动后，就会自动生成 VIP

## 作为 DaemonSet 运行

```bash
kubectl apply -f https://kube-vip.io/manifests/rbac.yaml
```