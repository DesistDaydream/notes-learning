---
title: Pod Monitor
---

# 概述

> 参考：
> - 官方文档：<https://github.com/prometheus-operator/prometheus-operator/blob/master/Documentation/design.md#podmonitor>

Pod Monitor 与 Service Monitor 一样，都是用来生成 Prometheus 配置文件中 scrape 配置段中的内容。

不同点在于 PM 直接与 pod 关联，根据标签选择来定义要监控的 pod，而不再需要通过 service 来暴露 pod 中的端口。

# PodMonitor yaml 详解

    apiVersion: monitoring.coreos.com/v1
    kind: PodMonitor
    metadata:
      name: rabbitmq
    spec:
      # 指定从 pod 中哪个端口采集指标，需要在 pod 的 .spec.containers.ports 字段中指定 containerPort 和 name。
      podMetricsEndpoints:
      - interval: 15s
        port: prometheus # 需要与 .spec.containers.ports.name 相同，则会将端口加入 scrape 配置中
      # 指定要匹配的 pod 的 label，具有相同 label 的将会加入监控配置。
      selector:
        matchLabels:
          app.kubernetes.io/component: rabbitmq
      # 指定要从哪个 namespace 中关联 pod。any: true 为匹配所有 ns 下的 pod
      namespaceSelector:
        any: true

# Pod Monitor 样例

    apiVersion: monitoring.coreos.com/v1
    kind: PodMonitor
    metadata:
      name: rabbitmq
    spec:
      podMetricsEndpoints:
      - interval: 15s
        port: prometheus
      selector:
        matchLabels:
          app.kubernetes.io/component: rabbitmq
      namespaceSelector:
        any: true
