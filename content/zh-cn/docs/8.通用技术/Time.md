---
title: Time
linkTitle: Time
date: 2024-03-27T17:32
weight: 20
---

# 概述

> 参考：
> 
> - [Wiki，ISO 8601](https://en.wikipedia.org/wiki/ISO_8601)
> - https://zh.wikipedia.org/zh-hans/%E5%90%84%E5%9C%B0%E6%97%A5%E6%9C%9F%E5%92%8C%E6%97%B6%E9%97%B4%E8%A1%A8%E7%A4%BA%E6%B3%95

## 时间格式

日期

日期格式为 YYYY-MM-DD，其中 YYYY 为年，MM 为月 (01–12)，DD 为月份日期 (01–31)。例如，2022 年 1 月 1 日显示为 2022-01-01。如果仅显示月和日，则格式为 MM-DD。例如，6 月 11 日显示为 06-11。

时间

时间格式为 hh:mm:ss，其中 hh 为小时 (00–24)，mm 为分钟 (00–60)，ss 为秒 (00–60)。如果仅显示小时和分钟，则格式为 hh:mm，例如 23:59。

# Timestamps

> 参考：
> 
> - [RFC 3339](https://tools.ietf.org/html/rfc3339)
> - [时间戳在线转换](https://www.bejson.com/convert/unix/)

**Timestamps(时间戳)**

<https://baike.baidu.com/item/unix> 时间戳

<https://baike.baidu.com/item/2038> 年问题

时间戳转换工具，可以把浮点型的毫秒数字串转换成真实时间

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/vym7ql/1616165162261-f89b406f-0967-44d2-b496-baa6ebe57434.png)

时间戳有两种格式

以秒为单位的当前 UNIX 时间戳 # 1562757107, 1562757108, 1562757109

$timestamp

ISO 格式(zero UTC)当前时间戳 # 2019-10-21T06:05:50.000Z

$isoTimestamp

Epoch Time

**Unix 时间(Unix Time)** 也叫做 **POSIX 时间** 或 **Epoch Time(纪元时间)**，是用来记录时间的流逝，定义为从 UTC 时间 1970 年 1 月 1 日 00:00 开始流逝的秒数，不考虑闰秒。
