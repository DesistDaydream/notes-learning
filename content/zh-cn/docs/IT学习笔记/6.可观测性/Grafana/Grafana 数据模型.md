---
title: Grafana数据模型
---

# 概述

> 参考：
> - [官方文档，开发者-构建插件-数据框架](https://grafana.com/docs/grafana/latest/developers/plugins/data-frames/)

# Data Frames(数据框架)

Grafana 支持各种不同的数据源，每个数据源都有自己的数据模型。为了实现这一点，Grafana 将来自每个数据源的查询结果合并为一个统一的数据结构，称为 **Data Frames(数据框架)**。
数据框架结构是从 R 编程语言和 Pandas 等数据分析工具中借用的概念。

> 数据帧在 Grafana 7.0+中可用，并且用更通用的数据结构代替了时间序列和表结构，该数据结构可以支持更大范围的数据类型。

本文档概述了数据框架结构以及如何在 Grafana 中处理数据。

## 数据框

数据框是面向列的表结构，这意味着它按列而不是按行存储数据。要了解这意味着什么，让我们看一下 Grafana 使用的 TypeScript 定义：

```go
interface DataFrame {
    name?:  string;
    // reference to query that create the frame
    refId?: string;
    fields: []Field;
}
```

本质上，数据框是 **Fields(字段)**\_ \_的集合，其中每个字段对应于一列。每个字段又由值的集合以及元信息（例如这些值的数据类型）组成。

```go
interface Field {
    name:    string;
// Prometheus like Labels / Tags
    labels?: Record<string, string>;
// For example string, number, time (or more specific primitives in the backend)
type:   FieldType;
// Array of values all of the same type
    values: Vector<T>;
// Optional display data for the field (e.g. unit, name over-ride, etc)
    config: FieldConfig;
}
```

让我们来看一个例子。下表说明了一个具有两个字段（\_时间\_和\_温度）\_的数据框。

|                     |      |
| ------------------- | ---- |
| 时间                | 温度 |
| 2020-01-02 03:04:00 | 45.0 |
| 2020-01-02 03:05:00 | 47.0 |
| 2020-01-02 03:06:00 | 48.0 |

每个字段具有三个值，并且字段中的每个值必须共享相同的类型。在这种情况下，时间字段中的所有值都是时间戳，而温度字段中的所有值都是数字。

数据帧的一个限制是，帧中的所有字段必须具有相同的长度才能成为有效的数据帧。

### 现场配置

数据帧中的每个字段都包含有关字段中值的可选信息，例如单位，缩放比例等。

通过将字段配置添加到数据框，Grafana 可以自动配置可视化。例如，您可以将 Grafana 配置为自动设置数据源提供的单位。

## 转变

除类型信息外，字段配置还支持在 Grafana 中进行\_数据转换\_。

数据转换是任何接受一个数据帧作为输入并返回另一个数据帧作为输出的函数。通过在插件中使用数据框，您可以免费获得一系列转换。

## 数据帧作为时间序列

具有至少一个时间字段的数据帧被视为\_时间序列\_。

有关时间序列的更多信息，请参阅我们的时间序列简介。

### 宽幅

当时间序列集合共享相同的\_时间索引\_时（每个时间序列中的时间字段都相同），它们可以以\_宽\_格式存储在一起。通过重用时间字段，我们可以减少发送到浏览器的数据量。

在此示例中，`cpu`每个主机的使用情况共享时间索引，因此我们可以将它们存储在同一数据帧中。

```shell
Name: Wide
Dimensions: 3 fields by 2 rows
+---------------------+-----------------+-----------------+
| Name: time          | Name: cpu       | Name: cpu       |
| Labels:             | Labels: host=a  | Labels: host=b  |
| Type: []time.Time   | Type: []float64 | Type: []float64 |
+---------------------+-----------------+-----------------+
| 2020-01-02 03:04:00 | 3               | 4               |
| 2020-01-02 03:05:00 | 6               | 7               |
+---------------------+-----------------+-----------------+
```

但是，如果两个时间序列不共享相同的时间值，则它们将表示为两个不同的数据帧。

    Name: cpu
    Dimensions: 2 fields by 2 rows
    +---------------------+-----------------+
    | Name: time          | Name: cpu       |
    | Labels:             | Labels: host=a  |
    | Type: []time.Time   | Type: []float64 |
    +---------------------+-----------------+
    | 2020-01-02 03:04:00 | 3               |
    | 2020-01-02 03:05:00 | 6               |
    +---------------------+-----------------+
    Name: cpu
    Dimensions: 2 fields by 2 rows
    +---------------------+-----------------+
    | Name: time          | Name: cpu       |
    | Labels:             | Labels: host=b  |
    | Type: []time.Time   | Type: []float64 |
    +---------------------+-----------------+
    | 2020-01-02 03:04:01 | 4               |
    | 2020-01-02 03:05:01 | 7               |
    +---------------------+-----------------+

当通过同一过程收集多个时间序列时，通常可以使用宽格式。在这种情况下，每次测量均以相同的间隔进行，因此将共享相同的时间值。

### 长格式

某些数据源以\_长\_格式（也称为\_窄\_格式）返回数据。这是 SQL 数据库返回的常见格式。

在长格式中，字符串值表示为单独的字段，而不是标签。结果，长格式的数据形式可能具有重复的时间值。

Grafana 可以检测长格式的数据帧并将其转换为宽格式。

> **注意：**当前仅在后端中支持长格式：。

例如，以下长格式的数据帧：

    Name: Long
    Dimensions: 4 fields by 4 rows
    +---------------------+-----------------+-----------------+----------------+
    | Name: time          | Name: aMetric   | Name: bMetric   | Name: host     |
    | Labels:             | Labels:         | Labels:         | Labels:        |
    | Type: []time.Time   | Type: []float64 | Type: []float64 | Type: []string |
    +---------------------+-----------------+-----------------+----------------+
    | 2020-01-02 03:04:00 | 2               | 10              | foo            |
    | 2020-01-02 03:04:00 | 5               | 15              | bar            |
    | 2020-01-02 03:05:00 | 3               | 11              | foo            |
    | 2020-01-02 03:05:00 | 6               | 16              | bar            |
    +---------------------+-----------------+-----------------+----------------+

可以转换为宽格式的数据帧：

    Name: Wide
    Dimensions: 5 fields by 2 rows
    +---------------------+------------------+------------------+------------------+------------------+
    | Name: time          | Name: aMetric    | Name: bMetric    | Name: aMetric    | Name: bMetric    |
    | Labels:             | Labels: host=foo | Labels: host=foo | Labels: host=bar | Labels: host=bar |
    | Type: []time.Time   | Type: []float64  | Type: []float64  | Type: []float64  | Type: []float64  |
    +---------------------+------------------+------------------+------------------+------------------+
    | 2020-01-02 03:04:00 | 2                | 10               | 5                | 15               |
    | 2020-01-02 03:05:00 | 3                | 11               | 6                | 16               |
    +---------------------+------------------+------------------+------------------+------------------+

## 技术参考

本节包含技术参考和数据帧实现的链接。

### 阿帕奇箭

数据框架结构受 Apache Arrow Project 启发并使用。Javascript 数据框架使用箭头表作为基础结构，后端 Go 代码在箭头表中序列化其框架以进行传输。

### Java 脚本

JavaScript 实现数据帧是在`/src/dataframe`文件夹和`/src/types/dataframe.ts`该的`@grafana/data`包。

### Go

有关数据帧的 Go 实现的文档，请参阅 github.com/grafana/grafana-plugin-sdk-go/data 软件包。
