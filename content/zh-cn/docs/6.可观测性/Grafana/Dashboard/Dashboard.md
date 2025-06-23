---
title: Dashboard
linkTitle: Dashboard
weight: 1
---

# 概述

> 参考：
>
> - [官方文档，仪表盘](https://grafana.com/docs/grafana/latest/dashboards/)
> - 借助 [Grafana Play](https://play.grafana.org/d/000000041/)，您可以探索并了解 Grafana 的各种 Panel 工作原理，从实际示例中学习以加速您的开发
>   - https://play.grafana.org/dashboards/f/PGJ1Fr4Zz/demo3a-grafana-features 这是各种 Panel 的 Demo 集合。

**Panel(面板)** 与 **Dashboard(仪表盘)**

**Panel(面板)** 是 Grafana 用于展示的基本 **Visualization(可视化)** 模块。多个 Panel(面板) 组成了一个 **Dashboard(仪表盘)**。每个面板都有各种各样的样式和格式设置选项。 面板可以在仪表板上拖放和重新排列。 它们也可以调整大小

对于 Grafana 来说，页面处理的数据实际上是一个一个的 **Field(字段)**。从 Grafana 数据模型 章节，可以发现，Grafana 从数据源拿到的数据都是统一的格式，Grafana 在面板处理数据时，其实就是对一堆 Field(字段) 来操作。从各个数据源获取到的数据，统一被放在 Field 中了。

# 时间选择器

![time-picker_1](https://notes-learning.oss-cn-beijing.aliyuncs.com/grafana/time-picker_1.png)


# Panels(面板) 的类型

> 其实就是说有多少种 Visualizations

### Alert list

警报列表面板允许您显示仪表板警报。您可以配置列表以显示当前状态或最近的状态更改。您可以在“警报”概述中了解有关警报的更多信息。

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/pkl4xq/1616067984679-ae1c7be6-6e4e-4446-a674-ecad72d5ea97.png)

### Bar gauge

条形表通过将每个字段减小为单个值来简化数据。您选择 Grafana 如何计算减少量

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/pkl4xq/1616067984672-8ab42ac2-73be-4c36-8ccf-6f87db0bcf48.png)

### Dashboard list

仪表板列表面板允许您显示指向其他仪表板的动态链接。可以将列表配置为使用加星标的仪表板，最近查看的仪表板，搜索查询和仪表板标签。

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/pkl4xq/1616067984683-7e6736e3-adc8-4a04-bd39-f5b6725de555.png)

### Gauge

仪表是一个单值面板，可以为每个系列，列或行重复一个仪表。

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/pkl4xq/1616067984690-0e8784ba-2105-4902-8bf9-e279372d17d8.png)

### Time series - 最常用的面板

该可视化是 Grafana 生态系统中最常用的。它可以渲染为一条线，一条点的路径或一系列条形图。这种类型的图具有足够的通用性，几乎可以显示任何时间序列数据。

详见: [Time series 面板](/docs/6.可观测性/Grafana/Dashboard/Time%20series%20面板/Time%20series%20面板.md)

老版本称为 Graph

### Heatmap

通过“热图”面板可视化，您可以查看一段时间内的直方图。有关直方图的更多信息，请参阅直方图和热图简介。

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/pkl4xq/1616067984704-f817d555-9bb6-4b2c-99c6-b61580ea35dc.png)

### Logs

日志面板可视化显示了来自支持日志的数据源的日志行，例如 Elastic，Influx 和 Loki。通常，您将在图形面板旁边使用此面板来显示相关过程的日志输出。

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/pkl4xq/1616067984673-d77300c9-b33b-45ab-a171-ca6c64f394c7.png)

### News

此面板显示 RSS feed。默认情况下，它显示 Grafana Labs 博客中的文章。

在“显示”部分的“ URL”字段中输入 RSS 的 URL。此面板类型不接受任何其他查询。

### Stat

“状态”面板显示一个较大的状态值，并带有可选的图形迷你图。您可以使用阈值控制背景或值颜色。

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/pkl4xq/1616067984686-ae274958-840c-4cdf-83be-10279eb5bc68.png)

### Table

表格面板非常灵活，支持时间序列以及表格，注释和原始 JSON 数据的多种模式。该面板还提供日期格式，值格式和着色选项

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/pkl4xq/1616067984704-6cac5049-5d7b-4eb9-bd9f-111d6867b0ec.png)

### Text

文本面板使您可以为仪表板制作信息和描述面板。

在“模式”下，选择要使用 markdown 还是 HTML 设置文本样式，然后在下面的框中输入内容。 Grafana 包含标题和段落以帮助您入门，或者您可以从其他编辑器粘贴内容。
