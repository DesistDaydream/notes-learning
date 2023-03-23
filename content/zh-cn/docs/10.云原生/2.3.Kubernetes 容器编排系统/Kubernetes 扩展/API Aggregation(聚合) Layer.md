---
title: API Aggregation(聚合) Layer
---

# 概述

> 参考；
>
> - 官方文档参考：
>   - <https://kubernetes.io/docs/concepts/extend-kubernetes/api-extension/apiserver-aggregation/>

api aggregation 称为 api 聚合。用于扩展 kubernetes 的 API 。如下所示。其中 v1beta1.metrics.k8s.io 是通过 prometheus-adapter 添加的新 API

    [root@master-1 ~]# kubectl get apiservices.apiregistration.k8s.io
    NAME                                   SERVICE                         AVAILABLE   AGE
    v1.                                    Local                           True        163d
    .........
    v1beta1.metrics.k8s.io                 monitoring/prometheus-adapter   True        120m
    v1beta1.networking.k8s.io              Local                           True        163d
    ......

聚合出来的 API 会关联到一个指定的 service 上，所有对该 API 发起的请求，都会交由该 service 并转发到其后端的 pod 进行处理。

下面是一个扩展 API 的样例，其中指定了该 API 所关联的 service

    apiVersion: apiregistration.k8s.io/v1
    kind: APIService
    metadata:
      name: v1beta1.metrics.k8s.io
    spec:
      group: metrics.k8s.io
      groupPriorityMinimum: 100
      insecureSkipTLSVerify: true
      service:
        name: prometheus-adapter
        namespace: monitoring
      version: v1beta1
      versionPriority: 100

API Aggregation 的核心功能是动态注册、发现汇总、安全代理。

# 配置 API Aggregation

官方文档：<https://kubernetes.io/docs/tasks/extend-kubernetes/configure-aggregation-layer/>

# 应用实例

可以参考 heapster 和 metrics-server 之间的过渡关系：kubectl top 命令解析.note 这里面谈到了 api 聚合 的用例。
