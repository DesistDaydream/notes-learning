---
title: "LokiQL 常见查询语句"
linkTitle: "LokiQL 常见查询语句"
date: "2023-08-28T11:41"
weight: 20
---

# 概述

> 参考：
> 
> -

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
