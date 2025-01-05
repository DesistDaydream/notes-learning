---
title: "LokiQL 常见查询语句"
linkTitle: "LokiQL 常见查询语句"
date: "2023-08-28T11:41"
weight: 20
---

# 概述

> 参考：
> 
> - [官方文档，查询 - 查询示例](https://grafana.com/docs/loki/latest/query/query_examples/)


**多重过滤**

过滤应该首先使用标签匹配器，然后是行过滤器，最后使用标签过滤器：

```logql
{cluster="ops-tools1", namespace="loki-dev", job="loki-dev/query-frontend"} |= "metrics.go" !="out of order" | logfmt | duration > 30s or status_code!="200"
```

**多解析器**

比如要提取以下格式日志行的方法和路径：

```log
level=debug ts=2020-10-02T10:10:42.092268913Z caller=logging.go:66 traceID=a9d4d8a928d8db1 msg="POST /api/prom/api/v1/query_range (200) 1.5s"
```

你可以像下面这样使用多个解析器：

```logql
{job="loki-ops/query-frontend"} | logfmt | line_format "{{.msg}}"
| regexp "(?P<method>\\w+) (?P<path>[\\w|/]+) \\((?P<status>\\d+?)\\) (?P<duration>.*)"
```

Notes: 其实这种 regexp 语法在新版已经可以用 pattern 代替了

```
pattern `<method> <path> (<status>) <duration>` 
```

首先通过 `logfmt` 解析器提取日志中的数据，然后使用 `| line_format` 重新将日志格式化为 `POST /api/prom/api/v1/query_range (200) 1.5s`，然后紧接着就是用 `regexp` 解析器通过正则表达式来匹配提前标签了。

## 格式化

下面的查询显示了如何重新格式化日志行，使其更容易阅读。

```logql
{cluster="ops-tools1", name="querier", namespace="loki-dev"}
  |= "metrics.go" != "loki-canary"
  | logfmt
  | query != ""
  | label_format query="{{ Replace .query \"\\n\" \"\" -1 }}"
  | line_format "{{ .ts}}\t{{.duration}}\ttraceID = {{.traceID}}\t{{ printf \"%-100.100s\" .query }} "
```

其中的 `label_format` 用于格式化查询，而 `line_format` 则用于减少信息量并创建一个表格化的输出。比如对于下面的日志行数据：

```log
level=info ts=2020-10-23T20:32:18.094668233Z caller=metrics.go:81 org_id=29 traceID=1980d41501b57b68 latency=fast query="{cluster=\"ops-tools1\", job=\"loki-ops/query-frontend\"} |= \"query_range\"" query_type=filter range_type=range length=15m0s step=7s duration=650.22401ms status=200 throughput_mb=1.529717 total_bytes_mb=0.994659
level=info ts=2020-10-23T20:32:18.068866235Z caller=metrics.go:81 org_id=29 traceID=1980d41501b57b68 latency=fast query="{cluster=\"ops-tools1\", job=\"loki-ops/query-frontend\"} |= \"query_range\"" query_type=filter range_type=range length=15m0s step=7s duration=624.008132ms status=200 throughput_mb=0.693449 total_bytes_mb=0.432718
```

经过上面的查询过后可以得到如下所示的结果：

```log
2020-10-23T20:32:18.094668233Z	650.22401ms	  traceID = 1980d41501b57b68	{cluster="ops-tools1", job="loki-ops/query-frontend"} |= "query_range"
2020-10-23T20:32:18.068866235Z	624.008132ms	traceID = 1980d41501b57b68	{cluster="ops-tools1", job="loki-ops/query-frontend"} |= "query_range"
```

# 通用

## 无数据时返回 0 值

https://community.grafana.com/t/how-can-i-turn-no-data-to-zero-in-loki/40694/14?u=lchddn0905

修改 Grafana 面板配置 Standard options - No value 设置为 0

## 对字符串的长度进行比较

利用 label_format 通过模板字符串将某个字符串类型的 Label 的长度做作为一个新 Label，添加到日志行上；然后利用行过滤表达式，根据字符串长度筛选日志

```logql
{job="my_job"} |= "mouldId" |= "msgId"
| pattern `[<_>] [<_>] [<level>] - <_> - <vendor>msgId:<msg_id>=[<_>(commandId=<command_id>, mouldId=<mould_id>, dbId=<db_id>)]`
| label_format mouldlen = "{{ .mould_id | len }}"
| mouldlen > 20
```
