---
title: kube-vip
---

# 概述

> 参考：
> - [官网](https://kube-vip.io/)
> - [官方文档,使用静态 Pod 部署 kube-vip](https://kube-vip.io/hybrid/static/)

kube-vip 是一个用于 Kubernetes 控制平面的 VIP 和 负载均衡器。他可以实现 Keepalived + HAProxy 的功能。

Kube-VIP 最初是为 Kubernetes 控制平面提供 HA 解决方案而创建的，随着时间的推移，它已经发展为将相同的功能合并到 Kubernetes 的 load-banlancers 类型的 Service 资源。

# 配置

Kube-VIP 通过命令行标志变更运行时行为

# 部署

## 使用静态 Pod 部署 kube-vip

生成 Manifests

```bash
docker run --rm --net host docker.io/plndr/kube-vip:v0.3.5 \
/kube-vip manifest pod \
--interface ${INTERFACE} \
--vip ${VIP} \
--controlplane \
--services \
--arp \
--leaderElection | tee  /etc/kubernetes/manifests/kube-vip.yaml
```

等待 kubelet 将 pod 启动后，就会自动生成 VIP
