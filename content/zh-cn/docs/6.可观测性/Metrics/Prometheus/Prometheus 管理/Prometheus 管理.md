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

官方表示为了避免夏令时问题，将 UTS 时区写入代码中，任何外部的配置都无法生效（e.g. 配置 /etc/timezone 无效）。

更多讨论在 [issue 500](https://github.com/prometheus/prometheus/issues/500)

其实这个限制是不影响使用的：

- 如果做可视化，Grafana 是可以做时区转换的。
- 如果是调接口，拿到了数据中的时间戳，想怎么处理都可以。
- 如果因为 Prometheus 自带的 UI 不是本地时间，看着不舒服，[2.16 版本](https://github.com/prometheus/prometheus/commit/d996ba20ec9c7f1808823a047ed9d5ce96be3d8f)的新版 Web UI 已经引入了 Local Timezone 的选项
- 如果仍然想改 Prometheus 代码来适应自己的时区，可以参考[这篇文章](https://zhangguanzhang.github.io/2019/09/05/prometheus-change-timezone/)。

# 重大变化

## V2.39

> 参考：
>
> - <https://mp.weixin.qq.com/s/RMtjCiWgTFnKhnTBQc-WLA>

大量的资源优化。改进了 relabeling 中的内存重用，优化了 WAL 重放处理，从 TSDB head series 中删除了不必要的内存使用， 以及关闭了 head compaction 的事务隔离等。
