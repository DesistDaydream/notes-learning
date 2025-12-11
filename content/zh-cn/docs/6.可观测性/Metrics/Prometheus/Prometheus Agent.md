---
title: Prometheus Agent
linkTitle: Prometheus Agent
weight: 7
---

# 概述

> 参考：
>
> - [Prometheus Blog, Introducing Prometheus Agent Mode, an Efficient and Cloud-Native Way for Metric Forwarding](https://prometheus.io/blog/2021/11/16/agent/)
> - [官方文档，Prometheus Agent 模式](https://prometheus.io/docs/prometheus/latest/prometheus_agent)

Prometheus Agent 是 Prometheus 二进制文件中内置的一种操作模式，它具有相同的抓取 API、语义、配置和发现机制；此代理模式禁用了 Prometheus 的一些常用功能（TSDB、警报和规则评估），并针对抓取和远程写入远程位置优化了二进制文件。

![](https://prometheus.io/assets/blog/2021-11-16/agent.png)
