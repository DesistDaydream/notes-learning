---
title: Exemplars
---

# 概述

> 参考；
>
> - <https://prometheus.io/docs/prometheus/latest/feature_flags/#exemplars-storage>
> - <https://grafana.com/docs/grafana/latest/basics/exemplars/>

这是啥？CNCF 可观测性白皮书最后提到了这个

### Trace ID 实际应用

我们讨论了在多个信号之间相互跳转的方法，但是它真的是有用的吗？让我们简单的看两个基本案例:

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/vx3gmg/1656494100400-266391e7-2b19-4845-a69a-2112f8128bbc.jpeg)

- 我们收到了一个关于超出 SLO (service level objectives) 的意外高错误率的告警。告警来源于错误的计数器值，我们看到请求暴增导致 501 errors。我们使用\_exemplar\_ 跳转到事例的 logs 以了解准确的可供人类阅读的错误消息中。错误似乎来自于依赖深层次的内部微服务系统，由于存在与 trace ID 匹配的 request ID，所以可以跳转到 traces。多亏了这一点，我们确切的了解到哪个 service/process 导致了这个问题，并进一步挖掘更多的信息。

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/vx3gmg/1656494100616-81e6c43f-e5c1-46c5-af8c-788ff1bc5f37.jpeg)

- 我们去 debug 慢请求，我们使用 trace 采样手动触发请求并获得 trace ID。多亏了 tracing view，我们可以在请求方式的几个进程中看到，对于基本操作而说，ABC-1 请求的速度非常的慢。由于目标元数据和时间，我们选择了相关的 CPU 使用率 metrics。我们看到 CPU 使用率很高，接近了机器的限制值，表明 CPU 已经饱和。为了了解 CPU 使用率高的原因 (特别是当它是容器中仅存的进程)，我们使用相同的 目标元数据 和 time 选择跳转到 CPU profile。

总结一下，好像是通过一个 ID 可以在 Metrics、Log、Trace 数据之间相互跳转。即一个 ID 关联了一个或多个应用所有的可观测性数据
