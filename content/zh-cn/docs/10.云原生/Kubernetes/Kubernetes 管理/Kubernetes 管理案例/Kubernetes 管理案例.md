---
title: Kubernetes 管理案例
linkTitle: Kubernetes 管理案例
weight: 1
---

# 概述

> 参考：

# 资源删除场景

## 处于 Terminating 状态的对象处理

使用 `kubectl edit` 命令来编辑该对象的配置，删除其中 [finalizers](docs/10.云原生/Kubernetes/Controller/Garbage%20Collection(垃圾收集)/Finalizers.md) 字段及其附属字段，即可.

也可以使用 patch 命令来删除 finalizers 字段

```bash
kubectl patch -n NS Resource ResourceName -p '{"metadata":{"finalizers":null}}' -n log
```

或

```bash
kubectl patch -n test configmap mymap \
    --type json \
    --patch='[ { "op": "remove", "path": "/metadata/finalizers" } ]'
```

## 资源无法删除

首先使用命令找到该 ns 还有哪些对象，最后的 NAMESPACE 改为自己想要查找的 ns 名

```bash
export NAMESPACE="test"

kubectl api-resources \
  --verbs=list --namespaced -o name | xargs -n 1 \
  kubectl get --show-kind --ignore-not-found -n NAMESPACE
```

找到对象后，删除，如果删不掉，使用处理 Terminationg 状态对象的方法进行处理
