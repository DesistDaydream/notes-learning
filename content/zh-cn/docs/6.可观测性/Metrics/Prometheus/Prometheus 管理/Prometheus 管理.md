---
title: Prometheus 管理
linkTitle: Prometheus 管理
date: 2024-06-28T17:16
weight: 20
---

# 概述

> 参考：
>
> - [官方文档，介绍 - FAQ](https://prometheus.io/docs/introduction/faq/)

## Prometheus UTS 时区问题

https://prometheus.io/docs/introduction/faq/#can-i-change-the-timezone-why-is-everything-in-utc

官方表示为了避免夏令时问题，将 UTS 时区写入代码中，任何外部的配置都无法生效。

更多讨论在 [issue 500](https://github.com/prometheus/prometheus/issues/500)

# 重大变化

## V2.39

> 参考：
>
> - <https://mp.weixin.qq.com/s/RMtjCiWgTFnKhnTBQc-WLA>

大量的资源优化。改进了 relabeling 中的内存重用，优化了 WAL 重放处理，从 TSDB head series 中删除了不必要的内存使用， 以及关闭了 head compaction 的事务隔离等。
