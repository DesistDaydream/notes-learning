---
title: "Label 与 Relabeling"
linkTitle: "Label 与 Relabeling"
date: "2023-08-02T14:28"
weight: 20
---

# 概述

> 参考：

Promtail 的 Label 与 Relabeling 功能与 Prometheus 中的 [Relabeling(重新标记)](/docs/6.可观测性/监控系统/Prometheus/Target(目标)%20与%20Relabeling(重新标记).md) 概念一样。

Promtail 具有一个嵌入式 Web 服务器，可以通过配置文件的 `server` 字段配置监听的端口，默认监听 80 端口

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/mdqko5/1616129665346-dc2414b8-d71a-4d16-864a-019c0706ec01.png)

这个 Web 页面与 Prometheus 的页面基本一样，只不过更简单，只有 Service Dicovery 和 Targets 两个页面。

> 也确实只要两个页面就够了，在 Loki 套件中，Promtail 就是用来 发现目标、重新标记、采集日志、推送日志 的。

在服务发现页面里，也有 Discovered Labels 和 Target Labels 这两个

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/mdqko5/1616129665338-c95fc783-cbd1-4e6c-958f-6a8a78448899.png)

只不过 Journal 这个目标发现程序命名发现了很多标签，但是却显示不出来，这个比较奇怪
