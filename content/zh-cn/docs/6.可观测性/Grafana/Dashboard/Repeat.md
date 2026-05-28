---
title: Repeat
linkTitle: Repeat
created: 2026-05-28T23:54
weight: 100
---

# 概述

> 参考：
>
> - https://grafana.com/docs/grafana/latest/panels-visualizations/configure-panel-options/#configure-repeating-panels
> - https://grafana.com/docs/grafana/latest/visualizations/dashboards/build-dashboards/create-dashboard/#configure-repeat-options
> - 官方示例: https://grafana.com/docs/grafana/latest/visualizations/dashboards/build-dashboards/create-dashboard/#configure-repeat-options

利用 [Variable](/docs/6.可观测性/Grafana/Dashboard/Variable.md) 与 [Panel](/docs/6.可观测性/Grafana/Dashboard/Panel/Panel.md)/Row 的 Repeat options，我们可以动态得生成多个 Panel 或 Row

配置 Repeat options 之前，我们展示多个目标，要不在一个 Panel 中，要不就要手动创建另一个额外的 Panel。使用了 Repeat 功能后，就会如图下半部分一样，自动生成了一个 Panel

> Note: 这里的 Repeat 使用了 instance 变量

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/grafana/repeat-demo-1.png)
