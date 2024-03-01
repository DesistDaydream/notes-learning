---
title: Timestamps(时间戳)
---

# 概述

> 参考：
> 
> - [RFC 3339](https://tools.ietf.org/html/rfc3339)
> - [时间戳在线转换](https://www.bejson.com/convert/unix/)

**Timestamps(时间戳)**

<https://baike.baidu.com/item/unix>时间戳

<https://baike.baidu.com/item/2038>年问题

时间戳转换工具，可以把浮点型的毫秒数字串转换成真实时间
![](https://notes-learning.oss-cn-beijing.aliyuncs.com/vym7ql/1616165162261-f89b406f-0967-44d2-b496-baa6ebe57434.png)
时间戳有两种格式

以秒为单位的当前 UNIX 时间戳 # 1562757107, 1562757108, 1562757109

$timestamp

ISO 格式(zero UTC)当前时间戳 # 2019-10-21T06:05:50.000Z

$isoTimestamp

Epoch Time

**Unix 时间(Unix Time) **也叫做 **POSIX 时间 **或 **Epoch Time(纪元时间)**，是用来记录时间的流逝，定义为从 UTC 时间 1970 年 1 月 1 日 00:00 开始流逝的秒数，不考虑闰秒。
