---
title: (弃用)Graph 类型面板详解
---

# 概述

> 参考：

这是一个初始的 Graph 面板，有两条查询语句，更改了序列的名称。

    sum(node_memory_MemTotal_bytes)
    (sum(node_memory_MemTotal_bytes{} - node_memory_MemAvailable_bytes{}) / sum(node_memory_MemTotal_bytes{}))*100

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/qzbd5e/1616067957167-730a2679-0ad0-488a-9c4c-8f3ba5ace79d.png)

Graph 是一个二维的，具有 x/y Axes(轴) 的面板。x 轴(横轴) 以时间分布、y 轴(纵轴) 以样本值分布

下面的文章将只介绍 Graph 面板的独有配置，有很多共有配置详见[此处](https://www.yuque.com/go/doc/33145831)

# Panel # 面板配置

## Display # 显示。设置值的显示样式(柱状、线条、圆点三种)

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/qzbd5e/1616067957211-044eecd5-5b98-425a-8de8-3799545d50f6.png)

在 Display 标签里，有 5 个开关，Bars、Lines、Staircase、Points 用于指定值的显示样式。Alert thresholds 用于指定是否显示告警阈值。

> 其他的配置选项，都是在开启某个样式后，才会显示对应样式专用的选项。
> Min step 设置时间长一点，Bars 与 Points 样式才可以看出来效果。否则都挤到一坨去了~

**Bars** # 柱状图样式。当 X 轴的模式变为 Series、Historgram 时，自动开启
**Lines** # 线条样式。默认样式

- **Staircase** # 开启线条样式时，才可以使用这种样式。
- **Line width** # 线条宽度。
- **Area fill** # 线条到底部这一区域的填充度。默认 1。
- **Fill gradient** # 填充渐变。默认 0。

**Points** # 圆点样式。
**Point radius** # 每个圆点的半径
**Alert thresholds** # 在面板上显示报警阈值和区域

除了设置面板中值的显示样式，还可以设置一些其他的设置
**Stacking and null value(叠加与空值)**
用于在面板上叠加所有 series 的值

**Hover tooltip(悬停提示) #**开启后，鼠标悬停在面板上，会出现一些关于 series 的信息
Mode # 模式。

- All series # 鼠标悬停到面板时，显示所有 series 的信息
- Single # 鼠标选定到面板时，只显示鼠标所在的 series 的信息。

Sort order # 排序。有三种排序方式：None(不排序)、Increasing(由上到下逐渐增大)、Decreasing(由上到下逐渐减小)

## Series overrides # 序列替换。用于个性化每个序列的配置

顾名思义，就是用来替换序列样式的。当一个面板上，配置了多个查询语句，这时就会产生多条 Series。而右侧的配置，是统一的，所有 Series 的配置内容都一样，这样不利于数据展示。所以通过 Series overrides 可以个性化得配置每一条 Series，让不同的 Series 展示出不同的效果(比如多条 Series 可以具有不同的单位、不同的线条宽度不同、不同的显示方式)

点击 `+ Add series override` 即可为指定的 series 进行配置

在 `Alias or regex` 选择要配置的序列。这里也可以使用正则表达式进行多个 series 的匹配。然后点击 `+` 符号，即可为选定的 series 进行单独的配置。

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/qzbd5e/1616067957244-03bc347c-faa5-4145-8a6b-fe3138242f0b.png)

### 应用实例

上面的例子，一个语句是内存用量，一个语句是内存使用率，单位是不一样(一个 KiB、一个百分比)。这时候，就需要使用 Series overrides，为每个 Serie 单独配置。不但单位可以分别配置，还可以将 Serie 移动到右侧的 Y 轴。还可以为不同的 series 配置不用的显示方式(比如有的用圆点、有点用线条、有的用柱状图，都可以在同一个面板显示出来)

比如我现在为 总平均使用率 序列进行单独配置，面板就会变成下面这种效果：
![](https://notes-learning.oss-cn-beijing.aliyuncs.com/qzbd5e/1616067957204-1cae00cd-e82c-4042-9d6b-11daed4f9b8f.png)

## Axis # 设定坐标轴的显示内容

在这里可以更改 x 轴 和 y 轴 的显示内容。常用于配置 metrics 值的 unit(单位)、Decimals(小数位数)。

**Left Y/Right Y** # 更改 Y-axes(Y 轴) 的信息

- Show # 是否显示这个轴
- Unit # 配置 Y 轴 的单位
- Decimals # 配置 Y 轴显示的小数位数。
- Label # 配置 Y 轴 的标签(标签会显示在 Y 轴 的旁边)

**Y-Axes** # Y 轴 配置。一些对齐方式
**X-Axis** # X 轴 配置。可以更改 X 轴的 Mode(模式)

- Mode # 模式。用于改变 X 轴的 显示模式。
  - Time # 时间模式。默认模式。X 轴 表示 时间，数据按时间分组（例如，按小时或分钟）。
  - Series # 序列模式。X 轴 表示 series，数据按照序列分组。Y 轴 仍然代表该序列的值
    - 注意：当 X 轴 切换到 Series 模式时，Display 配置中的将自动使用柱状图的方式
  - Histogram # 直方图模式。X 轴 表示 序列的值，Y 轴 表示 该值的计数。

### X 轴的 Series 模式 示例

当 X 轴 变为 series 模式 时，由于没有时间这种维度，所以一般都使用 当前值。
![](https://notes-learning.oss-cn-beijing.aliyuncs.com/qzbd5e/1616067957169-8ff35969-aa1d-4259-8144-1a88bb33a486.png)

## Legend # 用于配置面板内的 [Legend](/docs/6.可观测性/Grafana/Panel(面板)%20 与%20Dashboard(仪表盘)/Panel(面板)%20 配置详解.md 配置详解.md)

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/qzbd5e/1616067957204-998bd4db-f185-4f42-9138-9edb69081d84.png)

## Time regions # 时间区域。Graph 类型面板不常用

# Field # 字段配置

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/qzbd5e/1616067957229-733ae428-0442-42e4-b76e-d36f89196a3d.png)

# Overrides # 字段替换配置，Graph 类型面板不常用

详见：[Overrides](/docs/6.可观测性/Grafana/Panel(面板)%20 与%20Dashboard(仪表盘)/Panel(面板)%20 配置详解.md 配置详解.md)
