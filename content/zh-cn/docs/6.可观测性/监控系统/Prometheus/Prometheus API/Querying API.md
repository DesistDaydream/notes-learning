---
title: Querying API
---

# 概述

> 参考：
>
> - [官方文档，Prometheus-查询-HTTP API](https://prometheus.io/docs/prometheus/latest/querying/api/)
> - [OpenAPI](https://app.swaggerhub.com/apis/DesistDaydream/PrometheusAPI/v1)

Querying API 可以查询 Prometheus 的 运行时配置、时间序列数据、运行时状态 等等。官方称之为 HTTP API

## 在 HTTP API 中使用 PromQL

Prometheus 当前稳定的 HTTP API 可以通过 /api/v1 访问。

### API 响应格式

Prometheus API 使用了 JSON 格式的响应内容。 当 API 调用成功后将会返回 2xx 的 HTTP 状态码。

反之，当 API 调用失败时可能返回以下几种不同的 HTTP 状态码：

- 404 Bad Request # 当参数错误或者缺失时。
- 422 Unprocessable Entity # 当表达式无法执行时。
- 503 Service Unavailiable # 当请求超时或者被中断时。

这类 API 有一个统一的返回体

```json
{
  // 本次请求的结果
  "status": "success" | "error",
  // 本次请求所返回的具体数据，就是一个PromQL的查询结果
  "data": <data>,


  // 当 status 值为 error 时，会显示这两个字段。data 字段仍然会保留一些附加数据。
  "errorType": "<string>",
  "error": "<string>",

  // 当执行http请求出现warnings时，会显示这个字段。data 字段中仍然会有数据
  "warnings": ["<string>"]
}
```

status 与 data 字段是默认包含的。当 status 值不是 success 时，会出现其他字段。

不同的接口，data 字段下的格式就各不相同了。详见各 API 详解。

## API 中的 POST 请求

在下面详解 API 中，只是介绍了 GET 方法示例，其实还可以使用 `POST` 方法，并设定请求头 `Content-Type:application/x-www-form-urlencoded`，直接在请求体中以 x-www-form-urlencoded 格式填写那些所需参数。

使用 POST 请求有多个好处

- 当 PromQL 非常长，可能会导致 URL 的字符数超过服务端规定的最大 URL 字符数时，此功能很有用。
- 防止 curl 命令中使用编码

### 当 PromQL 非常长，可能会导致 URL 的字符数超过服务端规定的最大 URL 字符数时，此功能很有用。

比如，即时查询中的例子，还可以这么写：

```json
curl -X POST 'http://localhost:9090/api/v1/query' \
--header 'Content-Type: application/x-www-form-urlencoded' \
--data-urlencode 'query=up' \
--data-urlencode 'time=2021-05-18T14:28:51.781Z'

{
    "status": "success",
    "data": {
        "resultType": "vector",
        "result": [
            {
                "metric": {
                    "__name__": "up",
                    "instance": "localhost:9090",
                    "job": "prometheus"
                },
                "value": [
                    1621348131.781,
                    "1"
                ]
            },
            {
                "metric": {
                    "__name__": "up",
                    "instance": "localhost:9100",
                    "job": "node-exporter"
                },
                "value": [
                    1621348131.781,
                    "1"
                ]
            }
        ]
    }
}
```

### 防止 curl 命令中使用编码

此外，该功能对于使用 curl 调试也非常有用。比如该示例的符号，都进行了编码，而通过 POST 请求，则不必编码，即可使用

```bash
curl 'http://localhost:9090/api/v1/series?match[]=up&match[]=process_start_time_seconds%7Bjob%3D%22prometheus%22%7D'

# 上述语句可以替换为
curl -X POST 'http://localhost:9090/api/v1/series' \
--header 'Content-Type: application/x-www-form-urlencoded' \
--data-urlencode 'match[]=up' \
--data-urlencode 'match[]=process_start_time_seconds{job="prometheus"}'
```

# Expression Queries(表达式查询)

通过这类 API，可以使用 PromQL 来获取时间序列数据。这类 API 分两类：**Instant Queries(即时查询)** 与 **Range Queries(范围查询)**

## Instant Queries(即时查询)

```text
GET /api/v1/query
POST /api/v1/query
```

URL 请求参数：

- **query=\<STRING>** # PromQL 表达式。
- **time=\<TIMESTAMP>** # (可选参数)用于指定用于计算 PromQL 的时间戳。默认情况下使用当前系统时间。
- **timeout=\<DURATION>** # 超时设置。可选参数，默认情况下使用-query,timeout 的全局设置。

例如使用以下表达式查询表达式 up 在时间点 2015-07-01T20:10:51.781Z 的计算结果：

```json
$ curl 'http://localhost:9090/api/v1/query?query=up&time=2021-02-28T6:23:51.781Z' | jq .
{
  "status": "success",
  "data": {
    "resultType": "vector",
    "result": [
      {
        "metric": {
          "__name__": "up",
          "instance": "172.19.42.248:9100",
          "job": "node_exporter"
        },
        "value": [ 1614493431.781, "1"
]
      },
      {
        "metric": {
          "__name__": "up",
          "instance": "localhost:9090",
          "job": "prometheus"
        },
        "value": [
           1614493431.781,
           "1"
				]
      }
    ]
  }
}
```

## Range Queries(范围查询)

```text
GET /api/v1/query_range
POST /api/v1/query_range
```

URL 请求参数：

- **query=\<STRING>** # PromQL 表达式。
- **start=\<rfc3339 | unix_timestamp>** # 起始时间
- **end=\<rfc3339 | unix_timestamp>** # 结束时间
- **step=\<duration | float>** # 步长。起始时间与结束时间之间获取的所有数据的间隔时间。假如 step=10 则每隔 10 秒获取一次样本值。
  - 就好像人走路，一步迈多长，这里就是返回样本值时，每隔多久返回一次。
- **timeout=\<duration>** # 超时设置。可选参数，默认情况下使用-query,timeout 的全局设置。Evaluation timeout. Optional. Defaults to and is capped by the value of the -query.timeout flag.

例如使用以下表达式查询表达式 up 在 30 秒范围内以 15 秒为间隔计算 PromQL 表达式的结果。

```json
$ curl 'http://localhost:9090/api/v1/query_range?query=up&start=2015-07-01T20:10:30.781Z&end=2015-07-01T20:11:00.781Z&step=15s'
{
	"status": "success",
	"data": {
		"resultType": "matrix",
		"result": [{
					"metric": {
						"__name__": "up",
						"job": "prometheus",
						"instance": "localhost:9090"
					},
					"values": [
						["1435781430.781", "1"],
						["1435781445.781", "1"],
						["1435781460.781", "1"]
					]
				},
				{
					"metric": {
						"__name__": "up",
						"job": "node",
						"instance": "localhost:9091"
					},
					"values": [
						["1435781430.781", "0"],
						["1435781445.781", "0"],
						["1435781460.781", "1"]
					]
				}
```

## 表达式查询结果的响应数据的格式

当 API 调用成功后，Prometheus 会返回 JSON 格式的响应内容，格式如文章开头所示。并且在 `data` 字段中返回查询结果，`data` 字段包含 `resultType` 与` result` 两个字段：

```json
{
  "status": "success",
  "data": {
    "resultType": "matrix" | "vector" | "scalar" | "string",
    "result": <VALUE>
  }
}
```

`**resultType**`** 字段**表示该时间序列的值类型：

- **matrix(矩阵)** # 范围向量表达式查询结果
- **vector(向量)** # 即时向量表达式查询结果
- **scalar(标量)** # 标量表达式查询结果
- **string(字符串)** # 字符串表达式查询结果

**`result` 字段** 就是查询结果，该字段的值将会随 resultType 字段值的变化而不尽相同。这与 PromQL 章节描述的概念保持一致

### matrix

当返回数据类型 resultType 为 matrix 时，result 响应格式如下：

```json
[
  {
    "metric": { "<label_name>": "<label_value>", ... },
    "values": [ [ <unix_time>, "<sample_value>" ], ... ]
  },
  ...
]
```

其中 metrics 表示当前时间序列的特征维度，values 包含当前事件序列的一组样本。

### vector

当返回数据类型 resultType 为 vector 时，result 响应格式如下：

```json
[
  {
    "metric": { "<label_name>": "<label_value>", ... },
    "value": [ <unix_time>, "<sample_value>" ]
  },
  ...
]
```

其中 metrics 表示当前时间序列的特征维度，value 只包含一个唯一的样本。

### scalar

当返回数据类型 resultType 字段值为 scalar 时，result 字段格式如下：

```text
[ <unix_time>, "<scalar_value>" ]
```

由于标量不存在时间序列一说，因此 result 表示为当前系统时间一个标量的值。

### string

当返回数据类型 resultType 字段值为 string 时，result 字段格式如下：

```text
[ <unix_time>, "<string_value>" ]
```

字符串类型的响应内容格式和标量相同。

# Querying Metadata(查询元数据)

这类 API 可以获取相关数据的源数据列表，比如所有数据中，全部的标签列表，符合匹配规则的时间序列列表、甚至可以直接获取指定标签的标签值

## 获取 Series，仅获取序列，不包含该序列的值

```text
GET /api/v1/series
POST /api/v1/series
```

URL 请求参数：

- `match[]=<series_selector>` # 时间序列选择器，也就是 PromQL，该 API 会列出所有匹配到的序列，但是不包含序列的值。至少要有一个 match\[] 参数
- `start=<rfc3339 | unix_timestamp>`: 起始时间。可选的
- `end=<rfc3339 | unix_timestamp>`: 结束时间。可选的

`data` 字段中由一个对象列表组成，该对象包含匹配到的每个序列的标签

以下示例返回与选择器 up 或 process_start_time_seconds {job =“ prometheus”} 匹配的所有序列：

```json
$ curl 'http://localhost:9090/api/v1/series?match[]=up&match[]=process_start_time_seconds%7Bjob%3D%22prometheus%22%7D'
{
   "status" : "success",
   "data" : [
      {
         "__name__" : "up",
         "job" : "prometheus",
         "instance" : "localhost:9090"
      },
      {
         "__name__" : "up",
         "job" : "node",
         "instance" : "localhost:9091"
      },
      {
         "__name__" : "process_start_time_seconds",
         "job" : "prometheus",
         "instance" : "localhost:9090"
      }
   ]
}
```

## 获取 Labels 列表

根据匹配规则，获取匹配到的所有序列中包含的所有标签。若是不加参数，则返回所有的 Labels

```text
GET /api/v1/labels
POST /api/v1/labels
```

URL 请求参数：

- `start=<rfc3339 | unix_timestamp>`: 起始时间。可选的
- `end=<rfc3339 | unix_timestamp>`: 结束时间。可选的
- `match[]=<series_selector>` # 时间序列选择器，也就是 PromQL，该 API 从匹配到的序列中获取这些序列包含的所有标签。可选的

`data` 字段字符串类型的标签名列表

```json
curl --location --request GET 'http://test-prometheus.desistdaydream.ltd/api/v1/labels?match[]={job=%22prometheus%22}'
{
    "status": "success",
    "data": [
        "__name__",
        "alertmanager",
        "branch",
        "call",
        "code",
        "config",
        "dialer_name",
        "endpoint",
        "event",
        "goversion",
        "handler",
        "instance",
        "interval",
        "job",
        "le",
        "listener_name",
        "name",
        "quantile",
        "reason",
        "revision",
        "role",
        "rule_group",
        "scrape_job",
        "slice",
        "version"
    ]
}
```

## 获取 Labels 的值

根据匹配规则获取指定标签名称的标签值的列表

```text
GET /api/v1/label/<label_name>/values
```

URL 请求参数:

- `start=<rfc3339 | unix_timestamp>`: 起始时间。可选的
- `end=<rfc3339 | unix_timestamp>`: 结束时间。可选的
- `match[]=<series_selector>` # 时间序列选择器，也就是 PromQL，该 API 从匹配到的序列中获取指定指标的所有值。可选的

This example queries for all label values for the `job` label:

```json
$ curl http://localhost:9090/api/v1/label/job/values
{
   "status" : "success",
   "data" : [
      "node",
      "prometheus"
   ]
}
```

# Targets

```
GET /api/v1/targets
```

# Rules

```
GET /api/v1/rules
```

# Alerts

```
GET /api/v1/alerts
```

# Querying Target Metadata

```
GET /api/v1/targets/metadata
```

# Querying Metric Metadata

```
GET /api/v1/metadata
```

# Alertmanagers

```
GET /api/v1/alertmanagers
```

# Status

```
GET /api/v1/status/config
```
