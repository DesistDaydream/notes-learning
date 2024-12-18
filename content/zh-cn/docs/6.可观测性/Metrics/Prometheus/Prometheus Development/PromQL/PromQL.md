---
title: PromQL
linkTitle: PromQL
weight: 1
---

# 概述

> 参考：
>
> - [GitHub 项目，prometheus/prometheus - promql](https://github.com/prometheus/prometheus/tree/main/promql)

# 查询结果数据结构

代码中的查询结果类型，与 Querying API 中描述的一致

从 web/api/v1/api.go 文件可以看到所有可用的 HTTP API，从 `r.Get("/query", wrap(api.query))` 可以跳转到即时向量表达式的查询逻辑。

```go
func (api *API) query(r *http.Request) (result apiFuncResult) {
	......
	return apiFuncResult{&queryData{
		ResultType: res.Value.Type(),
		Result:     res.Value,
		Stats:      qs,
	}, nil, res.Warnings, qry.Close}
}
```

代码中的查询结果类型，与 Querying API 中描述的一致，主要是 4 个字段：status、data 中的 resultType 与 result

```go
type queryData struct {
	ResultType parser.ValueType  `json:"resultType"`
	Result     parser.Value      `json:"result"`
	Stats      *stats.QueryStats `json:"stats,omitempty"`
}
```

parser.Value 是一个接口

```go
type Value interface {
	Type() ValueType
	String() string
}
```

一共有 4 个结构体实现了该接口

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/dg9h9g/1616068620395-100edb14-c157-4c00-9d02-60937abc9df5.png)

这四个结构体也就是 4 个结果类型 matrix、vector、scalar、string

## 直接看最常用的 matrix 和 vector

Matrix 是 Series 类型的数组

```go
type Series struct {
	Metric labels.Labels `json:"metric"`
	Points []Point       `json:"values"`
}
type Point struct {
	T int64
	V float64
}
```

Series 结构体的属性结构与通过 API 查询获取的结果保持一致

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/dg9h9g/1616068620417-18503054-4f3f-4421-9eea-25d62ca25884.png)

有两个字段分别对应 Metrics 与 Points 属性，Metrics 就是表示指标的唯一标识符，即一组标签的合集，Points 就是指标的值数组，Point 里有两个，样本值和时间戳。

# Functions

> 参考：
>
> - [GitHub 项目，prometheus/prometheus - promql/functions.go](https://github.com/prometheus/prometheus/blob/main/promql/functions.go)

## 变化量

下面三个函数，在代码中的实现逻辑都是由同一个方法实现的

`./prometheus/promql/functions.go`

```go
// extrapolatedRate rate、increase、delta 函数功能的具体实现
// 该函数计算速率，如果第一个或最后一个样本接近所选范围的边界，则 extrapolates，并以每秒或整体返回结果
// 如果 isCounter 参数为 true，则允许计数器重置。
// 如果 isRate 参数为 true，则返回结果以每秒作为单位
func extrapolatedRate(vals []parser.Value, args parser.Expressions, enh *EvalNodeHelper, isCounter bool, isRate bool) Vector {
	......
}

func instantValue(vals []parser.Value, out Vector, isRate bool) (Vector, annotations.Annotations) {
	......
}

// === delta(Matrix parser.ValueTypeMatrix) Vector ===
func funcDelta(vals []parser.Value, args parser.Expressions, enh *EvalNodeHelper) Vector {
	return extrapolatedRate(vals, args, enh, false, false)
}

// === rate(node parser.ValueTypeMatrix) Vector ===
func funcRate(vals []parser.Value, args parser.Expressions, enh *EvalNodeHelper) Vector {
	return extrapolatedRate(vals, args, enh, true, true)
}

// === increase(node parser.ValueTypeMatrix) Vector ===
func funcIncrease(vals []parser.Value, args parser.Expressions, enh *EvalNodeHelper) Vector {
	return extrapolatedRate(vals, args, enh, true, false)
}

// === irate(node parser.ValueTypeMatrix) (Vector, Annotations) ===
func funcIrate(vals []parser.Value, args parser.Expressions, enh *EvalNodeHelper) (Vector, annotations.Annotations) {
	return instantValue(vals, enh.Out, true)
}

// === idelta(node model.ValMatrix) (Vector, Annotations) ===
func funcIdelta(vals []parser.Value, args parser.Expressions, enh *EvalNodeHelper) (Vector, annotations.Annotations) {
	return instantValue(vals, enh.Out, false)
}
```

区别

- delta() 与 rate() 的区别在于 isCounter 和 isRate 参数
- delta() 与 increase() 的区别在于 isCounter 参数
- rate() 与 increase() 的区别在于 isRate 参数

```go
extrapolationThreshold := averageDurationBetweenSamples * 1.1 // 外推阈值
sampledInterval := float64(lastT-firstT) / 1000
extrapolateToInterval := sampledInterval
factor := extrapolateToInterval / sampledInterval
```

