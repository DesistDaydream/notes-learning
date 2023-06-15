---
title: Thanos 部署
---

# 概述

> 参考：
> 
> - [官方文档,快速教程](https://thanos.io/tip/thanos/quick-tutorial.md/#sidecar)

# 通过 docker 启动 Thanos

## Sidecar

## Query

# Store

# 在 Kubernetes 集群中部署 Thanos

> 参考：
> 
> - [GitHub,examples(Manifests)](https://github.com/thanos-io/kube-thanos/tree/main/examples)
> - [通过 kubectl 插件 kube-thanos 部署 Store 与 Query 的 Manifests](https://github.com/thanos-io/kube-thanos/tree/main/manifests)

## Sidecar 或 Receiver

由于 Sidecar 的工作性质，所以，Sidecar 组件最好作为 Prometheus 的 sidecar 容器，部署在同一个 Pod 中。

## Query

## Store

## Compact

# 通过 prometheus-operator 部署 Thanos

> 参考：
> 
> - [官方文档,prometheus operator-thanos](https://prometheus-operator.dev/docs/operator/thanos/)
> - [GitHub 文档,prometheus-operator-文档-thanos](https://github.com/prometheus-operator/prometheus-operator/blob/master/Documentation/thanos.md)
