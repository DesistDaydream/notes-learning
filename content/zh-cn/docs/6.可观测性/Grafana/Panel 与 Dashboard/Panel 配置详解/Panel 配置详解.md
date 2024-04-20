---
title: "Panel 配置详解"
linkTitle: "Panel 配置详解"
weight: 1
---

# 概述

> 参考：
>
> - [官方文档,面板-面板编辑器](https://grafana.com/docs/grafana/latest/panels/panel-editor/)

当我们开始创建一个新的 Panel 时，可以看到下图所示的界面，这个界面分为三大部分，分别用三种颜色的框体括起来

![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/ldaq0w/1636261571312-2e5d4a25-2009-495d-919a-973f3d2cd178.png)

- [**Header(标题)**](https://grafana.com/docs/grafana/latest/panels/panel-editor/#header)**,绿色部分** # 左边是面板的名称，右侧有 4 个按钮，分别是 设置整个 Dashboard、放弃、保存、应用
- [**Visualization preview(可视化的预览)**](https://grafana.com/docs/grafana/latest/panels/panel-editor/#visualization-preview)**,蓝色部分** # 在 数据处理 与 面板样式处理 两部分设置的内容将会反应在这个预览部分
- [**Data section(数据处理)**](https://grafana.com/docs/grafana/latest/panels/panel-editor/#data-section-bottom-pane)**,红色部分** #[](https://grafana.com/docs/grafana/latest/panels/panel-editor/#data-section-bottom-pane)通过数据查询语句来获取数据，以便在面板展示
- [**Panel dispaly options(面板显示选项)**](https://grafana.com/docs/grafana/latest/panels/panel-editor/#panel-display-options-side-pane)**,黄色部分** # 用来配置面板的信息。包括 面板类型、面板名称、展示效果 等等

# Header(标题)

标题部分列出了面板所在的仪表板的名称和一些仪表板命令。您还可以单击**返回**箭头以返回仪表板。

[![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/ldaq0w/1636274879674-172ae5e2-cbdb-42d2-a0de-404c9f32cce7.png)](https://grafana.com/static/img/docs/panel-editor/edit-panel-header-7-0.png)

标题的右侧是以下选项：

- **仪表板设置（齿轮）图标** - 单击以访问仪表板设置。
- **Discard(放弃)** - 放弃自上次保存仪表板以来对面板所做的所有更改。
- **Save(保存)** - 保存仪表板，包括您在面板编辑器中所做的所有更改。
- **Apply(应用)** - 应用您所做的更改，然后关闭面板编辑器，将您返回到仪表板。您必须保存仪表板以保留应用的更改。

# Visualization preview(可视化的预览)

![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/ldaq0w/1636266413475-671d1896-dc09-44f7-8018-772384c8ebac.png)
在可视化的预览部分，可以图像的形式查看从数据源获取到的数据。包含如下几个部分

- Axes 横、纵 坐标轴
  - 横轴是时间
  - 纵轴是值。即该时间点上，数据源中的值。
- Panel title 面板标题
- Legend 图例(即 图片的文字说明)
- 时间范围控件

在左上角是 Grafana 模板

- 在没有配置 Grafana 变量时，这一块是空白的。详见 [模板与变量](/docs/6.可观测性/Grafana/Panel(面板)%20 与%20Dashboard(仪表盘)/Panel(面板)%20 配置详解/Templates%20and%20Variables(模板与变量).md 配置详解/Templates and Variables(模板与变量).md)

在右上角可以控制部分显示效果

- **Table view(表格视图)** # 将可视化预览区域转换为表格，以便查看数据。常用来故障排除
- **Fill(填充)** # 可视化预览将填充预览部分中的可用空间。如果您更改侧窗格的宽度或底部窗格的高度，则可视化将适应以填充任何可用空间。
- **Fit(适合)** # 可视化预览将填充其中的可用空间，但保留面板的纵横比。
- **Exact(确切)** # 可视化预览的大小将与仪表板上的大小完全相同。如果没有足够的可用空间，则可视化将按比例缩小以保留宽高比。
- **Time range controls(时间范围控件)** # 有关更多信息，请参阅[时间范围控件](https://grafana.com/docs/grafana/latest/dashboards/time-range-controls/)。

# Data section(数据处理)

该部分包含一些 tab(标签)，可以在其中 输入查询，转换数据 以及 创建警报规则(如果适用)。

- **Query tab(查询标签)** # 选择数据源并通过查询语句获取数据。参考：[Queries](https://grafana.com/docs/grafana/latest/panels/queries/).
  - ![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/ldaq0w/1636266421210-10e7bbbd-f661-463b-bb0f-1d53b53ffa47.png)
- **Transform tab(转换标签)** # 将 Query 中获取到的数据进行转换。参考：[Transformations](https://grafana.com/docs/grafana/latest/panels/transformations/).
  - ![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/ldaq0w/1636274705492-be88d84d-0e38-40d3-8e80-fe71204480ad.png)
- **Alert tab(告警标签)**# 配置告警规则。参考：[Create alerts](https://grafana.com/docs/grafana/latest/alerting/create-alerts/)
  - ![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/ldaq0w/1636274711891-c3fbf5e9-144d-40d6-a434-170d49a7b3f3.png)

## Query(查询)

在查询标签中，可以配置 Grafana 与数据源并获取可视化数据的方式。

Query 标签的页面由一下几个元素组成

- Data source selector(数据源选择器)
- Query options(查询选项)
- Query inspector button(查询检查器按钮)
- Query editor list(查询编辑器)
- Expressions(表达式)

可以输入数据源的查询语句，查询结果将会实时展现在面板上，比如，从 Prometheus 中获取了节点内存使用率的数据：

![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/ldaq0w/1636274684211-5c3cddae-a4f0-4fc1-b10b-0c43417df8dc.png)

Query 详解见 [Query(查询)](/docs/6.可观测性/Grafana/Panel%20与%20Dashboard/Panel%20配置详解/Query(查询).md)

## Transform(转换)

**Transform(转换)** 用于将查询结果在面板展示之前进行处理。Transform 可以重命名字段、将单独的时间序列连接在一起、在查询中进行数学运算等等。常用于 Table 类型的面板

> 官方文档称为 Tranformations

Transformations process the result set of a query before it’s passed on for visualization. They allow you to rename fields, join separate time series together, do math across queries, and more. For users, with numerous dashboards or with a large volume of queries, the ability to reuse the query result from one panel in another panel can be a huge performance gain.

> **注意**：转换是 Grafana 7.0 测试版的特性。官方文档的内容会在开发者们研究该特性时经常更新。

Transformations sometimes result in data that cannot be graphed. When that happens, Grafana displays a suggestion on the visualization that you can click to switch to table visualization. This often helps you better understand what the transformation is doing to your data

Transform 的用法详见：[Transformations(转换)](/docs/6.可观测性/Grafana/Panel%20与%20Dashboard/Panel%20配置详解/Transformations(转换).md)

## Alert(告警)

# Panel display options(面板显示选项)

该部分包含一些 tab(标签)，可以在这部分配置几乎都有数据可视化的方面。但是并不是所有选项都可用于每个面板类型。这里只介绍所有面板的通用配置，不同的面板，这部分的可配置的项目不同
可以用来配置面板的样式、面板的字段、以及如何为每个字段进行单独配置

- **Panel(面板标签)** # 配置面板
- **Field(字段标签)** # 配置所有字段
- **Overrides(替换标签)** # 根据匹配规则替换指定字段的配置。这是对 Field 的补充，可以实现个性化为每个字段配置不同的样式

## Visualization(可视化) # 指定面板的类型

[Visualizations(可视化)](https://grafana.com/docs/grafana/latest/panels/visualizations/) 用来指定当前面板的类型。Grafana 提供了多种多样的 Visualizations 来适应不同的环境。说白了，一个 Visualizations 就是一个 Panel Type(面板类型)。并且，可以通过 [plugins](https://grafana.com/docs/grafana/latest/plugins/) 来添加更多类型的面板。

![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/ldaq0w/1636262032681-56ebd4dc-2991-4b02-a4e9-c126dcf099e0.png)

在添加面板时，右侧的最上面点击一下就能看到当前可以使用的所有面板。Visualization 标签内可以看到 Grafana 默认自带的一些面板

![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/ldaq0w/1636262106569-8229fbdc-eca7-4e02-85cf-992099d794d1.png)

## Panel # 面板配置

> 参考：
>
> - [官方文档](https://grafana.com/docs/grafana/latest/panels/field-options/)

Grafna 中使用的 [数据模型][/docs/6.可观测性/Grafana/Grafana%20数据模型.md](数据模型)是面向列的表结构，该结构将时间序列和表查询结果统一在一起。此结构中的每一列称为一个 **Field(字段)**。一个字段可以代表`一条时间序列(Prometheus源)`或`表格的列(数据库源)。`

> 每条序列的时间，也算作一个 Field
> 在 Table 类型的面板中，每条时间序列的标签也算作一个 Field

这里只介绍所有面板的通用配置，不同的面板，这部分的可配置的项目不同

### Panel options(面板选项) # 面板的基本信息。名字、描述、透明度。参考[此处](https://grafana.com/docs/grafana/latest/panels/add-a-panel/)

![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/ldaq0w/1636268147075-e492ebe5-f2a3-4cfe-9b8b-ec466199f0ad.png)

### [Standard options](https://grafana.com/docs/grafana/latest/panels/standard-options/)(标准选项)

![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/ldaq0w/1636269601040-63c30a7e-cf11-4cc1-9e86-08466afe414e.png)
设置 单位、显示名、小数点 等等

### Thresholds(阈值)

Thresholds(阈值) 可以用于 Bar、Gauge、Graph、Stat、Table 这几种类型的面板

当在下方的 Alert 标签内开始配置告警后，该标签变为不可用状态
![](https://notes-learning.oss-cn-beijing.aliyuncs.com/ldaq0w/1616067971510-785cd511-5ead-4465-8f72-88c0898a1922.png)
这是一个 Graph 面板配置阈值的样子
![](https://notes-learning.oss-cn-beijing.aliyuncs.com/ldaq0w/1616067971526-3a0fa5b8-96bc-458c-9395-3d43fb0ca76c.png)

### Value mappings(值映射)

### Data links(链接)

## Overrides # 替换，替换 Field

> 参考：[官方文档](https://grafana.com/docs/grafana/latest/panels/field-options/)

这里只介绍所有面板的通用配置，不同的面板，这部分的可配置的项目不同
根据匹配规则，替换面板上某些字段。常用于 Table 类型的面板。Overrides 的概念与 [Graph 类型面板里的 Series overrides](/docs/6.可观测性/Grafana/Panel(面板)%20 与%20Dashboard(仪表盘)/Time%20series%20 类型面板/(弃用)Graph%20 类型面板详解.md series 类型面板/(弃用)Graph 类型面板详解.md) 概念类似
