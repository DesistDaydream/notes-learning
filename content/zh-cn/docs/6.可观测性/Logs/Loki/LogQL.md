---
title: LogQL
---

# 概述

> 参考：
>
> - [官方文档，LogQL](https://grafana.com/docs/loki/latest/logql/)
> - [公众号，Loki 查询语言 LogQL 使用](https://mp.weixin.qq.com/s/0dXT0fIreZk6_4ZL4S8lHg)

**Log Query Language(日志查询语言，简称 LogQL)** 受 PromQL 启发，可以看作是分布式的 grep 命令，用来从汇总的日志源筛选日志。LogQL 通过 labels(标签) 和 operators(运算符) 进行过滤。

**LogQL 查询有两种类型：**

- **Log Queries(日志查询)** # 根据查询语句返回日志条目，每行是一条日志。
  - 基本的日志查询由两部分组成
    - log stream selector # 日志流选择器
    - filter expression # 过滤表达式
- **Metric Queries(指标查询)** # 用于扩展日志查询并根据 Log Queries 中的日志计数计算值。通过这种查询语句，可以计算将日志数据量化成指标信息，并且，Promtail 可以通过这种查询语句将指标信息，填充到自己暴露的 Metrics 端点中。

注意：由于 Loki 的设计，所有 LogQL 查询都必须包含一个日志流选择器

日志流选择器确定将搜索多少日志流（日志内容的唯一来源，例如文件）。然后，更细粒度的日志流选择器将搜索到的流的数量减少到可管理的数量。这意味着传递给日志流选择器的标签将影响查询执行的相对性能。然后使用过滤器表达式对来自匹配日志流的聚合日志进行分布式 grep。

# Log Queries(日志查询)

基本的日志查询由两部分组成：

- **Log Stream Selector(日志流选择器)**#
- **Log Pipeline(日志管道)**#

注意：由于 Loki 的设计原则，所有的 LogQL 查询必须包含 Log Stream Selector(日志流选择器)

![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/xgx9x4/1621837564259-416660f0-81ef-4d14-9086-dbda268daf9f.png)

**日志流选择器**决定了有多少日志流将被搜索到，一个更细粒度的日志流选择器将搜索到流的数量减少到一个可管理的数量。所以传递给日志流选择器的标签将影响查询执行的性能。

而日志流选择器后面的**日志管道**是可选的，日志管道是一组阶段表达式，它们被串联在一起应用于所过滤的日志流，每个表达式都可以过滤、解析和改变日志行内容以及各自的标签。

下面的例子显示了一个完整的日志查询的操作：

`{container="query-frontend",namespace="loki-dev"} |= "metrics.go" | logfmt | duration > 10s and throughput_mb < 500`

该查询语句由以下几个部分组成：

- 一个日志流选择器 `{container="query-frontend",namespace="loki-dev"}`，用于过滤 `loki-dev` 命名空间下面的 `query-frontend` 容器的日志
- 然后后面跟着一个日志管道 `|= "metrics.go" | logfmt | duration > 10s and throughput_mb < 500`，这管道表示将筛选出包含 `metrics.go` 这个词的日志，然后解析每一行日志提取更多的表达并进行过滤

> 为了避免转义特色字符，你可以在引用字符串的时候使用单引号，而不是双引号，比如 `\w+1` 与 "\w+" 是相同的。

## Log Stream Selector(日志流选择器)

Log Stream Selector 用于确定查询结果中应该包括哪些日志流。Log Stream Selector 由一个或多个 **Key/Value Paire(键/值对)** 组成，其中每个键是一个日志标签，每个值是该标签的值。所谓的日志流就是一行一行的日志，组合在一起，形成的一种类似数据流的感觉，从上到下哗哗流水那种感觉~日志流说白了就是日志的集合。stream(流) 的概念如果在 Prometheus 中描述，那就是 series(序列) 的概念。

Log Stream Selector 中的键值对应包装在一对花括号中，比如：

```json
{job="kube-system/etcd",container="etcd"}
```

在上面这个例子中，所有具有 job 标签，值为 kube-system/etcd 和 container 标签，值为 etcd 的日志流将被包含在查军结果中。

这种语法与 Prometheus 标签选择器 的语法一样。参考 PromQL,prometheus 查询语言 文章中 即时向量 章节中的匹配说明

上面 LogQL 的执行效果如下：

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/xgx9x4/1616129551117-ca616a54-c0f1-48e3-868b-9f43f2138f1a.jpeg)

注意：如果想要获取所有日志，使用这种方式：

`{job=~"..*"}`

## Log Pipeline(日志管道)

Log Pipeline(日志管道) 可以通过 `|` 符号附加到 Log Stream Selector(日志流选择器) 语句后面，以便对日志流进一步处理和过滤。

Log Pipeline 通常由一个或多个 **Expression(表达式)** 组成，多个表达式以 `|` 符号分割。这用法其实跟 Linux 中使用命令并通过管道传递结果给下一个命令的用法一模一样。

可用的 Log Pipeline 有如下几种

- [**Line Filter Expression**](https://grafana.com/docs/loki/latest/logql/#line-filter-expression)**(行过滤表达式)** # 最基本的过滤方式，通过关键字匹配每一行的日志内容
- [**Parser Expression**](https://grafana.com/docs/loki/latest/logql/#parser-expression)**(解析表达式)**# 以指定方式解析日志内容，并将解析结果提取为标签。
- [**Label Filter Expression**](https://grafana.com/docs/loki/latest/logql/#label-filter-expression)**(标签过滤表达式)** #
- [**Line Format Expression**](https://grafana.com/docs/loki/latest/logql/#line-format-expression)**(行格式化表达式)** #
- [**Labels Format Expression**](https://grafana.com/docs/loki/latest/logql/#labels-format-expression)**(标签格式化表达式)** #
- [**Unwrap Expression**](https://grafana.com/docs/loki/latest/logql/#unwrapped-range-aggregations)#
  - 这是一个特殊的表达式，只能在指标查询中使用。

其中一些表达式可以改变日志内容和相应的标签，然后可用于进一步 过滤和处理表达式 或 指标查询。

### Line Filter Expression(行过滤表达式)

通过 日志流选择器 获取到想要的日志后，可以使用 Line Filter Expression(行过滤表达式) 对这些日志进行过滤。过滤表达式 可以只是文本或正则表达式，比如

```bash
# 过滤出日志内容中，包含 timeout 字符串的日志行。
{job="kube-system/etcd",container="etcd"} |= "timeout"
# 匹配 {job="nginx-promtail"} 日志流中所有日志行中，不包含 天津市 字符串的行
{job="nginx-promtail"} != "天津市"
```

注意：过滤表达式不能单独使用，必须基于 日志流选择器 得出的结果，再进行过滤。示例 LogQL 执行效果如下

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/xgx9x4/1616129550979-719a6401-7dd8-4196-9f8a-78f7b14e80a8.jpeg)

在上面的示例中， |= 这个符号作为 filter operators(过滤运算符)，来定义过滤行为。可用的 filter operators 有以下几种：

1. **|=** # 匹配包含指定字符串的日志行
2. **!=** # 匹配不包含指定字符串的日志行
3. **|~** # 匹配包含正则表达式的日志行
4. **!~**# 匹配不包含正则表达式的日志行

### Parser Expression(解析表达式)

Parser Expression 可以将日志内容解析，并提取标签。然后这些被提取出来的标签，可以使用 Label Filter Expression

解析器表达式可以解析和提取日志内容中的标签，这些提取的标签可以用于标签过滤表达式进行过滤，或者用于指标聚合。

提取的标签键将由解析器进行自动格式化，以遵循 Prometheus 指标名称的约定（它们只能包含 ASCII 字母和数字，以及下划线和冒号，不能以数字开头）。

例如下面的日志经过管道 `| json` 将产生以下 Map 数据：

    { "a.b": { "c": "d" }, "e": "f" }

->

    {a_b_c="d", e="f"}

在出现错误的情况下，例如，如果该行不是预期的格式，该日志行不会被过滤，而是会被添加一个新的 `__error__` 标签。

需要注意的是如果一个提取的标签键名已经存在于原始日志流中，那么提取的标签键将以 `_extracted` 作为后缀，以区分两个标签，你可以使用一个标签格式化表达式来强行覆盖原始标签，但是如果一个提取的键出现了两次，那么只有最新的标签值会被保留。

目前支持 `json`、`logfmt`、`regexp` 和 `unpack` 这几种解析器。

我们应该尽可能使用 `json` 和 `logfmt` 等预定义的解析器，这会更加容易，而当日志行结构异常时，可以使用 `regexp`，可以在同一日志管道中使用多个解析器，这在你解析复杂日志时很有用。

#### JSON

json 解析器有两种模式运行。

- 1. 没有参数。如果日志行是一个有效的 json 文档，在你的管道中添加 `| json` 将提取所有 json 属性作为标签，嵌套的属性会使用 `_` 分隔符被平铺到标签键中。

> 注意：数组会被忽略。

- 例如，使用 json 解析器从以下文件内容中提取标签。

```json
{
  "protocol": "HTTP/2.0",
  "servers": ["129.0.1.1", "10.2.1.3"],
  "request": {
    "time": "6.032",
    "method": "GET",
    "host": "foo.grafana.net",
    "size": "55",
    "headers": {
      "Accept": "*/*",
      "User-Agent": "curl/7.68.0"
    }
  },
  "response": {
    "status": 401,
    "size": "228",
    "latency_seconds": "6.031"
  }
}
```

- 可以得到如下所示的标签列表：

```text
"protocol" => "HTTP/2.0"
"request_time" => "6.032"
"request_method" => "GET"
"request_host" => "foo.grafana.net"
"request_size" => "55"
"response_status" => "401"
"response_size" => "228"
"response_size" => "228"
```

- 2. 带有参数的。在你的管道中使用 `|json label="expression", another="expression"` 将只提取指定的 json 字段为标签，你可以用这种方式指定一个或多个表达式，与 `label_format` 相同，所有表达式必须加引号。

> 当前仅支持字段访问（`my.field`, `my["field"]`）和数组访问（`list[0]`），以及任何级别嵌套中的这些组合（`my.list[0]["field"]`）。
> 例如，`|json first_server="servers[0]", ua="request.headers[\"User-Agent\"]` 将从以下日志文件中提取标签：

```json
{
  "protocol": "HTTP/2.0",
  "servers": ["129.0.1.1", "10.2.1.3"],
  "request": {
    "time": "6.032",
    "method": "GET",
    "host": "foo.grafana.net",
    "size": "55",
    "headers": {
      "Accept": "*/*",
      "User-Agent": "curl/7.68.0"
    }
  },
  "response": {
    "status": 401,
    "size": "228",
    "latency_seconds": "6.031"
  }
}
```

- 提取的标签列表为：

```text
"first_server" => "129.0.1.1"
"ua" => "curl/7.68.0"
```

- 如果表达式返回一个数组或对象，它将以 json 格式分配给标签。例如，`|json server_list="services", headers="request.headers` 将提取到如下标签：

```text
"server_list" => `["129.0.1.1","10.2.1.3"]`
"headers" => `{"Accept": "*/*", "User-Agent": "curl/7.68.0"}`
```

#### logfmt

`logfmt` 解析器可以通过使用 `|logfmt` 来添加，它将从 logfmt 格式的日志行中提前所有的键和值。

例如，下面的日志行数据：

```
at=info method=GET path=/ host=grafana.net fwd="124.133.124.161" service=8ms status=200
```

将提取得到如下所示的标签：

```text
"at" => "info"
"method" => "GET"
"path" => "/"
"host" => "grafana.net"
"fwd" => "124.133.124.161"
"service" => "8ms"
"status" => "200"
```

#### regexp

与 `logfmt` 和 `json`（它们隐式提取所有值且不需要参数）不同，`regexp` 解析器采用单个参数 `| regexp "<re>"` 的格式，其参数是使用 Golang RE2 语法的正则表达式。

正则表达式必须包含至少一个命名的子匹配（例如`(?P<name>re)`），每个子匹配项都会提取一个不同的标签。

例如，解析器 `| regexp "(?P<method>\\w+) (?P<path>[\\w|/]+) \\((?P<status>\\d+?)\\) (?P<duration>.*)"` 将从以下行中提取标签：

    POST /api/prom/api/v1/query_range (200) 1.5s

提取的标签为：

    "method" => "POST"
    "path" => "/api/prom/api/v1/query_range"
    "status" => "200"
    "duration" => "1.5s"

#### unpack

`unpack` 解析器将解析 json 日志行，并通过打包阶段解开所有嵌入的标签，一个特殊的属性 `_entry` 也将被用来替换原来的日志行。

例如，使用 `| unpack` 解析器，可以得到如下所示的标签：

```json
{
  "container": "myapp",
  "pod": "pod-3223f",
  "_entry": "original log message"
}
```

允许提取 `container` 和 `pod` 标签以及原始日志信息作为新的日志行。

> 如果原始嵌入的日志行是特定的格式，你可以将 unpack 与 json 解析器（或其他解析器）相结合使用。

### [Label Filter Expression](https://grafana.com/docs/loki/latest/logql/#label-filter-expression)(标签过滤表达式)

标签过滤表达式允许使用其原始和提取的标签来过滤日志行，它可以包含多个谓词。

一个谓词包含一个标签标识符、操作符和用于比较标签的值。

例如 `cluster="namespace"` 其中的 `cluster` 是标签标识符，操作符是 `=`，值是`"namespace"`。

LogQL 支持从查询输入中自动推断出的多种值类型：

- `String（字符串）`用双引号或反引号引起来，例如`"200"`或`us-central1`。
- `Duration（时间）`是一串十进制数字，每个数字都有可选的数和单位后缀，如 `"300ms"`、`"1.5h"` 或 `"2h45m"`，有效的时间单位是 `"ns"`、`"us"`（或 `"µs"`）、`"ms"`、`"s"`、`"m"`、`"h"`。
- `Number（数字）`是浮点数（64 位），如 250、89.923。
- `Bytes（字节）`是一串十进制数字，每个数字都有可选的数和单位后缀，如 `"42MB"`、`"1.5Kib"` 或 `"20b"`，有效的字节单位是 `"b"`、`"kib"`、`"kb"`、`"mib"`、`"mb"`、`"gib"`、`"gb"`、`"tib"`、`"tb"`、`"pib"`、`"bb"`、`"eb"`。

字符串类型的工作方式与 Prometheus 标签匹配器在日志流选择器中使用的方式完全一样，这意味着你可以使用同样的操作符（`=`、`!=`、`=~`、`!~`）。

使用 Duration、Number 和 Bytes 将在比较前转换标签值，并支持以下比较器。

- `==` 或 `=` 相等比较
- `!=` 不等于比较
- `>` 和 `>=` 用于大于或大于等于比较
- `<` 和 `<=` 用于小于或小于等于比较

例如 `logfmt | duration > 1m and bytes_consumed > 20MB` 过滤表达式。

如果标签值的转换失败，日志行就不会被过滤，而会添加一个 `__error__` 标签，要过滤这些错误，请看管道错误部分。

你可以使用 `and`和 `or` 来连接多个谓词，它们分别表示**且**和**或**的二进制操作，`and` 可以用逗号、空格或其他管道来表示，标签过滤器可以放在日志管道的任何地方。

以下所有的表达式都是等价的:

    | duration >= 20ms or size == 20kb and method!~"2.."
    | duration >= 20ms or size == 20kb | method!~"2.."
    | duration >= 20ms or size == 20kb,method!~"2.."
    | duration >= 20ms or size == 20kb method!~"2.."

默认情况下，多个谓词的优先级是从右到左，你可以用圆括号包装谓词，强制使用从左到右的不同优先级。

例如，以下内容是等价的：

    | duration >= 20ms or method="GET" and size <= 20KB
    | ((duration >= 20ms or method="GET") and size <= 20KB)

它将首先评估 `duration>=20ms or method="GET"`，要首先评估 `method="GET" and size<=20KB`，请确保使用适当的括号，如下所示。

    | duration >= 20ms or (method="GET" and size <= 20KB)

### [Line Format Expression](https://grafana.com/docs/loki/latest/logql/#line-format-expression)(行格式化表达式)

日志行格式化表达式可以通过使用 Golang 的 `text/template` 模板格式重写日志行的内容，它需要一个字符串参数 `| line_format "{{.label_name}}"` 作为模板格式，所有的标签都是注入模板的变量，可以用 `{{.label_name}}` 的符号来使用。

例如，下面的表达式：

    {container="frontend"} | logfmt | line_format "{{.query}} {{.duration}}"

将提取并重写日志行，只包含 `query` 和请求的 `duration`。你可以为模板使用双引号字符串或反引号 `{{.label_name}}` 来避免转义特殊字符。

此外 `line_format` 也支持数学函数，例如：

如果我们有以下标签 `ip=1.1.1.1`, `status=200` 和 `duration=3000(ms)`, 我们可以用 `duration` 除以 1000 得到以秒为单位的值：

    {container="frontend"} | logfmt | line_format "{{.ip}} {{.status}} {{div .duration 1000}}"

上面的查询将得到的日志行内容为`1.1.1.1 200 3`。

### [Labels Format Expression](https://grafana.com/docs/loki/latest/logql/#labels-format-expression)(标签格式化表达式)

`| label_format`表达式可以重命名、修改或添加标签，它以逗号分隔的操作列表作为参数，可以同时进行多个操作。

当两边都是标签标识符时，例如 `dst=src`，该操作将把 `src` 标签重命名为 `dst`。

左边也可以是一个模板字符串，例如 `dst="{{.status}} {{.query}}"`，在这种情况下，`dst` 标签值会被 Golang 模板执行结果所取代，这与 `| line_format` 表达式是同一个模板引擎，这意味着标签可以作为变量使用，也可以使用同样的函数列表。

在上面两种情况下，如果目标标签不存在，那么就会创建一个新的标签。

重命名形式 `dst=src` 会在将 `src` 标签重新映射到 `dst` 标签后将其删除，然而，模板形式将保留引用的标签，例如 `dst="{{.src}}"` 的结果是 `dst` 和 `src` 都有相同的值。

> 一个标签名称在每个表达式中只能出现一次，这意味着 `| label_format foo=bar,foo="new"` 是不允许的，但你可以使用两个表达式来达到预期效果，比如 `| label_format foo=bar | label_format foo="new"`。

## 查询示例

**多重过滤**

过滤应该首先使用标签匹配器，然后是行过滤器，最后使用标签过滤器：

    {cluster="ops-tools1", namespace="loki-dev", job="loki-dev/query-frontend"} |= "metrics.go" !="out of order" | logfmt | duration > 30s or status_code!="200"

**多解析器**

比如要提取以下格式日志行的方法和路径：

    level=debug ts=2020-10-02T10:10:42.092268913Z caller=logging.go:66 traceID=a9d4d8a928d8db1 msg="POST /api/prom/api/v1/query_range (200) 1.5s"

你可以像下面这样使用多个解析器：

    {job="cortex-ops/query-frontend"} | logfmt | line_format "{{.msg}}" | regexp "(?P<method>\\w+) (?P<path>[\\w|/]+) \\((?P<status>\\d+?)\\) (?P<duration>.*)"`

首先通过 `logfmt` 解析器提取日志中的数据，然后使用 `| line_format` 重新将日志格式化为 `POST /api/prom/api/v1/query_range (200) 1.5s`，然后紧接着就是用 `regexp` 解析器通过正则表达式来匹配提前标签了。

**格式化**

下面的查询显示了如何重新格式化日志行，使其更容易阅读。

    {cluster="ops-tools1", name="querier", namespace="loki-dev"}
      |= "metrics.go"
      |!= "loki-canary"
      | logfmt
      | query != ""
      | label_format query="{{ Replace .query \"\\n\" \"\" -1 }}"
      | line_format "{{ .ts}}\t{{.duration}}\ttraceID = {{.traceID}}\t{{ printf \"%-100.100s\" .query }} "

其中的 `label_format` 用于格式化查询，而 `line_format` 则用于减少信息量并创建一个表格化的输出。比如对于下面的日志行数据：

    level=info ts=2020-10-23T20:32:18.094668233Z caller=metrics.go:81 org_id=29 traceID=1980d41501b57b68 latency=fast query="{cluster=\"ops-tools1\", job=\"cortex-ops/query-frontend\"} |= \"query_range\"" query_type=filter range_type=range length=15m0s step=7s duration=650.22401ms status=200 throughput_mb=1.529717 total_bytes_mb=0.994659
    level=info ts=2020-10-23T20:32:18.068866235Z caller=metrics.go:81 org_id=29 traceID=1980d41501b57b68 latency=fast query="{cluster=\"ops-tools1\", job=\"cortex-ops/query-frontend\"} |= \"query_range\"" query_type=filter range_type=range length=15m0s step=7s duration=624.008132ms status=200 throughput_mb=0.693449 total_bytes_mb=0.432718

经过上面的查询过后可以得到如下所示的结果：

    2020-10-23T20:32:18.094668233Z 650.22401ms     traceID = 1980d41501b57b68 {cluster="ops-tools1", job="cortex-ops/query-frontend"} |= "query_range"
    2020-10-23T20:32:18.068866235Z 624.008132ms traceID = 1980d41501b57b68 {cluster="ops-tool

# Metric Queries(指标查询)

LogQL 还可以通过 **Functions(函数)** 来对每个日志流进行计算，就是将日志流中的日志计数，并根据技术进行一些查询，这种查询方式与 PromQL 的指标查询类似。

指标查询可用于计算诸如错误消息率或最近 3 个小时内日志数量最多的 N 个日志源之类的信息。

## 范围向量聚合

LogQL 与 Prometheus 具有相同的范围向量概念，不同之处在于所选的样本范围包括每个日志

常用函数主要是如下 4 个：

- `rate`: 计算每秒的日志条目
- `count_over_time`: 对指定范围内的每个日志流的条目进行计数
- `bytes_rate`: 计算日志流每秒的字节数
- `bytes_over_time`: 对指定范围内的每个日志流的使用的字节数

比如计算 nginx 的 qps：

    rate({filename="/var/log/nginx/access.log"}[5m]))

计算 kernel 过去 5 分钟发生 oom 的次数：

    count_over_time({filename="/var/log/message"} |~ "oom_kill_process" [5m]))

## 聚合运算符

与 PromQL 一样，LogQL 支持内置聚合运算符的一个子集，可用于聚合单个向量的元素，从而产生一个具有更少元素但具有集合值的新向量：

- sum: Calculate sum over labels

- min: Select minimum over labels

- max: Select maximum over labels

- avg: Calculate the average over labels

- stddev: Calculate the population standard deviation over labels

- stdvar: Calculate the population standard variance over labels

- count: Count number of elements in the vector

- bottomk: Select smallest k elements by sample value

- topk: Select largest k elements by sample value

- sum：求和

- min：最小值

- max：最大值

- avg：平均值

- stddev：标准差

- stdvar：标准方差

- count：计数

- bottomk：最小的 k 个元素

- topk：最大的 k 个元素

聚合函数我们可以用如下表达式描述：

    <aggr-op>([parameter,] <vector expression) [without|by (label list)]

对于需要对标签进行分组时，我们可以用 `without` 或者 `by` 来区分。比如计算 nginx 的 qps，并按照 pod 来分组：

    sum(rate({filename="/var/log/nginx/access.log"}[5m])) by (pod)

只有在使用 `bottomk` 和 `topk` 函数时，我们可以对函数输入相关的参数。比如计算 nginx 的 qps 最大的前 5 个，并按照 pod 来分组：

    topk(5,sum(rate({filename="/var/log/nginx/access.log"}[5m])) by (pod)))

## Binary Operators(二元运算符)

### 数学计算

Loki 存的是日志，都是文本，怎么计算呢？显然 LogQL 中的数学运算是面向区间向量操作的，LogQL 中的支持的二进制运算符如下：

- `+`：加法
- `-`：减法
- `*`：乘法
- `/`：除法
- `%`：求模
- `^`：求幂

比如我们要找到某个业务日志里面的错误率，就可以按照如下方式计算：

    sum(rate({app="foo", level="error"}[1m])) / sum(rate({app="foo"}[1m]))

### 逻辑运算

集合运算仅在区间向量范围内有效，当前支持

- `and`：并且
- `or`：或者
- `unless`：排除

比如：

    rate({app=~"foo|bar"}[1m]) and rate({app="bar"}[1m])

### 比较运算

LogQL 支持的比较运算符和 PromQL 一样，包括：

- `==`：等于
- `!=`：不等于
- `>`：大于
- `>=`: 大于或等于
- `<`：小于
- `<=`: 小于或等于

通常我们使用区间向量计算后会做一个阈值的比较，这对应告警是非常有用的，比如统计 5 分钟内 error 级别日志条目大于 10 的情况：

    count_over_time({app="foo", level="error"}[5m]) > 10

我们也可以通过布尔计算来表达，比如统计 5 分钟内 error 级别日志条目大于 10 为真，反正则为假：

    count_over_time({app="foo", level="error"}[5m]) > bool 10

## 注释

LogQL 查询可以使用 `#` 字符进行注释，例如：

    {app="foo"} # anything that comes after will not be interpreted in your query

对于多行 LogQL 查询，可以使用 `#` 排除整个或部分行：

    {app="foo"}
        | json
        # this line will be ignored
        | bar="baz" # this checks if bar = "baz"

## Pipeline Errors 管道错误

There are multiple reasons which cause pipeline processing errors, such as:有多种原因导致流水线处理错误，例如：

- A numeric label filter may fail to turn a label value into a number 数字标签过滤器可能无法将标签值转换为数字
- A metric conversion for a label may fail.标签的度量转换可能会失败。
- A log line is not a valid json document.日志行不是有效的 JSON 文档。
- etc…等等…

When those failures happen, Loki won’t filter out those log lines. Instead they are passed into the next stage of the pipeline with a new system label named `__error__`. The only way to filter out errors is by using a label filter expressions. The `__error__` label can’t be renamed via the language.当这些故障发生时，Loki 不会过滤掉这些日志线。相反，它们通过名为\_\_Error\_\_的新系统标签传递到管道的下一个阶段。过滤误差的唯一方法是使用标签过滤表达式。 \_\_Err\_\_标签无法通过语言重命名。
For example to remove json errors:例如要删除 JSON 错误：

    {cluster="ops-tools1",container="ingress-nginx"}
        | json
        | __error__ != "JSONParserErr"

Logql
Alternatively you can remove all error using a catch all matcher such as `__error__ = ""` or even show only errors using `__error__ != ""`.或者，您可以使用捕获器删除所有匹配器，例如\_\_Error\_\_ =“”甚至仅显示使用\_\_Error\_\_！=“”。
The filter should be placed after the stage that generated this error. This means if you need to remove errors from an unwrap expression it needs to be placed after the unwrap.在生成此错误的阶段后应放置过滤器。这意味着如果您需要从未包装中删除从未包装表达式中删除错误，则需要放置在未包装之后。

    quantile_over_time(
     0.99,
     {container="ingress-nginx",service="hosted-grafana"}
     | json
     | unwrap response_latency_seconds
     | __error__=""[1m]
     ) by (cluster)

Logql

> Metric queries cannot contains errors, in case errors are found during execution, Loki will return an error and appropriate status code.度量标准查询不能包含错误，以便在执行期间找到错误，Loki 将返回错误和适当的状态代码。
