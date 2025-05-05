---
title: Prometheus MGMT
linkTitle: Prometheus MGMT
weight: 1
---

# 概述

> 参考：
>
> - [官方文档，介绍 - FAQ](https://prometheus.io/docs/introduction/faq/)

[公众号，Prometheus 性能调优：大模型和智能驾驶可观测的高基数问题](https://mp.weixin.qq.com/s/Lw03n5-Vvy5OSMaQTPlW-A)

## Prometheus UTS 时区问题

https://prometheus.io/docs/introduction/faq/#can-i-change-the-timezone-why-is-everything-in-utc

官方表示，为了避免夏令时问题，将 UTS 时区写入代码中，任何外部的配置都无法生效（e.g. 配置 /etc/timezone 无效）。

更多讨论在 [issue 500](https://github.com/prometheus/prometheus/issues/500)

其实这个限制是不影响使用的：

- 如果做可视化，Grafana 是可以做时区转换的。
- 如果是调接口，拿到了数据中的时间戳，想怎么处理都可以。
- 如果因为 Prometheus 自带的 UI 不是本地时间，看着不舒服，[2.16 版本](https://github.com/prometheus/prometheus/commit/d996ba20ec9c7f1808823a047ed9d5ce96be3d8f)的新版 Web UI 已经引入了 Local Timezone 的选项
- 如果仍然想改 Prometheus 代码来适应自己的时区，可以参考[这篇文章](https://zhangguanzhang.github.io/2019/09/05/prometheus-change-timezone/)。

对于 Prometheus 生态的程序，e.g. [Node Exporter](/docs/6.可观测性/Metrics/Instrumenting/Node%20Exporter.md)、etc. 也会有 UTS 时区问题，程序的日志时间就是 UTS 时区的，并且无法在程序实例化阶段通过代码修改，因为这些程序引用的是 [promlog](https://pkg.go.dev/github.com/prometheus/common/promlog) 库。promlog 在 [log.go](https://github.com/prometheus/common/blob/v0.60.0/promlog/log.go#L33) 中定义了日志的时区。

```go
	// This timestamp format differs from RFC3339Nano by using .000 instead
	// of .999999999 which changes the timestamp from 9 variable to 3 fixed
	// decimals (.130 instead of .130987456).
	timestampFormat = log.TimestampFormat(
		func() time.Time { return time.Now().UTC() },
		"2006-01-02T15:04:05.000Z07:00",
	)
```

若想使用其他时区，现阶段（2024-10-14）的解决方案是在编译时直接修改 promlog 库中上面的代码，将 `.UTC()` 去掉

# 重大变化

## v2.39

> 参考：
>
> - <https://mp.weixin.qq.com/s/RMtjCiWgTFnKhnTBQc-WLA>

大量的资源优化。改进了 relabeling 中的内存重用，优化了 WAL 重放处理，从 TSDB head series 中删除了不必要的内存使用， 以及关闭了 head compaction 的事务隔离等。

## v3.0

> 参考：
>
> - [公众号 - CNCF，Prometheus 3.0](https://mp.weixin.qq.com/s/cNvyJx5E3hrHYlfcjSWm0Q)
> - [官方博客，Prometheus 3.0 migration guide](https://prometheus.io/docs/prometheus/3.0/migration/)

由于 TSDB 格式变化，2.0 - 3.0 之间的版本变化需要先升级到 v2.55 才可以升级到 v3