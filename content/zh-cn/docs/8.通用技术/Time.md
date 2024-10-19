---
title: Time
linkTitle: Time
date: 2024-03-27T17:32
weight: 20
---

# 概述

> 参考：
>
> - [Wiki, ISO 8601](https://en.wikipedia.org/wiki/ISO_8601)
> - [Wiki-cn，各地日期和时间表示法](https://zh.wikipedia.org/zh-hans/%E5%90%84%E5%9C%B0%E6%97%A5%E6%9C%9F%E5%92%8C%E6%97%B6%E9%97%B4%E8%A1%A8%E7%A4%BA%E6%B3%95)
> - https://baike.baidu.com/item/ISO%208601/3910715?fr=aladdin
> - [中国计量科学研究院，SI 基本单位](https://www.nim.ac.cn/520/node/4.html)

1983 年，国际计量大会讨论决定，把 1 米的定义修改为光在 1/299792458 秒内走过的距离

1967 年，国际计量大会定义：1 秒是铯 133 原子基态的两个超精细能量间跃迁对应辐射的 9192631770 个周期的持续时间。

# Timestamps

> 参考：
>
> - [RFC 3339，Date and Time on the Internet: Timestamps](https://tools.ietf.org/html/rfc3339)
> - [Wiki, Timestamp](https://en.wikipedia.org/wiki/Timestamp)

**Timestamps(时间戳)** 是识别特定事件发生时间的字符序列或编码信息，通常给出 date(日期) 和 time(时间)，有时精确到一秒的一小部分。然而，时间戳不必基于某种绝对的时间概念。它们可以具有任何纪元，可以相对于任何任意时间，例如系统的开机时间，或相对于过去的某个任意时间。

**日期**

日期格式为 YYYY-MM-DD，其中 YYYY 为年，MM 为月 (01–12)，DD 为月份日期 (01–31)。例如，2022 年 1 月 1 日显示为 2022-01-01。如果仅显示月和日，则格式为 MM-DD。例如，6 月 11 日显示为 06-11。

**时间**

时间格式为 hh:mm:ss，其中 hh 为小时 (00–24)，mm 为分钟 (00–60)，ss 为秒 (00–60)。如果仅显示小时和分钟，则格式为 hh:mm，例如 23:59。

**日期和时间的组合表示法**

合并表示时，要在时间前面加一大写字母 T，如要表示北京时间 2004 年 5 月 3 日下午 5 点 30 分 8 秒，可以写成 `2004-05-03T17:30:08+08:00` 或 `20040503T173008+08`。

# Unix time

> 参考：
>
> - [Wiki, Unix_time](https://en.wikipedia.org/wiki/Unix_time)
> - [时间戳在线转换](https://www.bejson.com/convert/unix/)

**Unix time** 也叫做 **POSIX 时间** 或 **Epoch Time(纪元时间)**，是计算中广泛使用的日期和时间表示形式。它通过自 1970 年 1 月 1 日 Unix 纪元 00:00:00 UTC 以来经过的非闰秒数来测量时间。在现代计算中，值有时以更高的粒度存储，例如微秒或纳秒。

<https://baike.baidu.com/item/unix> 时间戳

<https://baike.baidu.com/item/2038> 年问题

时间戳转换工具，e.g. https://it-tools.tech/date-converter 、etc. 。可以在 Unix-time 与 Timestamps 之间相互转换。

![600](https://notes-learning.oss-cn-beijing.aliyuncs.com/time/time_converter.png)


