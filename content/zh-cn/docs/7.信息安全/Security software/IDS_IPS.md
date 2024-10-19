---
title: IDS/IPS
linkTitle: IDS/IPS
date: 2024-06-19T09:00
weight: 20
---

# 概述

> 参考：
>
> - [Wiki, Intrusion detection system](https://en.wikipedia.org/wiki/Intrusion_detection_system)

**Intrusion detection/prevention system(入侵监测系统 与 入侵防御系统，简称 IDS/IPS)** 是一种设备或软件应用程序，用于监视网络或系统是否存在恶意活动或策略违规行为。

这套系统通常包含如下几部分

- **规则** # 一种人类可读的过滤规则
- **规则库** # 特定于某种识别场景的一组规则，识别某些特定的恶意流量
- **引擎** # 将规则翻译成流量过滤的语句以定位某个或某些流量

用 [Snort](/docs/7.信息安全/Security%20software/Snort.md) 举例，Snort 本身可以表示 一种规则格式、一个识别规则的引擎、一个由 N 个规则组合而成的规则库；这些东西组合在一起，可以称之为一套系统。那么当流量来了之后（实时的流量 或者 .pcap 文件），Snort 可以读取流量，根据 Snort 规则库中的 Snort 规则，利用 Snort 识别引擎，对流量进行匹配后，识别出哪些流量是被规则命中的。

> 所谓引擎，应该是一种把自己的定义的人类可读的规则翻译成流量过滤语句的技术。

同理，其他的比如 Suricate、Yara、自研 的都是类似的道理，我们可以把这些系统中的规则通过某种方式进行转换，比如把人类可读的 Snort 规则转换成 Suricata 的规则后，由 Suricata 引擎再对流量识别。
