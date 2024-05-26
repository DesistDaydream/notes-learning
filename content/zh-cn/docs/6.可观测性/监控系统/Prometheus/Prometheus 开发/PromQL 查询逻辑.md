---
title: PromQL 查询逻辑
---

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
