---
title: RPA
linkTitle: RPA
date: 2024-01-04T14:29
weight: 20
---

# 概述

> 参考：
>
> - [Wiki，Robotic_process_automation](https://en.wikipedia.org/wiki/Robotic_process_automation)

**Robotic process automation(机器人的流程自动化，简称 RPA)** 是一种基于 软件机器人 或 [人工智能](/docs/12.AI/12.AI.md) 的流程自动化形式，这是一种实现自动化的设计思想、概念，并不是特指某个具体的实物。

在传统的工作流程自动化工具中，软件开发人员会生成一系列操作来自动执行任务，并使用内部 [API](/docs/2.编程/API/API.md) 或专用脚本语言连接到后端系统。相比之下，RPA 系统通过观察用户在应用程序的图形用户界面 (GUI) 中执行该任务来开发操作列表，然后通过直接在 GUI 中重复这些任务来执行自动化。这可以降低在可能不具有用于此目的的 API 的产品中使用自动化的障碍。

RPA 工具与 [GUI 测试](https://en.wikipedia.org/wiki/Graphical_user_interface_testing)工具在技术上有很强的相似性。这些工具还可以自动执行与 GUI 的交互，并且通常通过重复用户执行的一组演示操作来实现。 RPA 工具与此类系统的不同之处在于，它们允许在多个应用程序内和之间处理数据，例如，接收包含发票的电子邮件，提取数据，然后将其输入簿记系统。

# 实现工具

UiPath # 国外很有名

UiBot # 按键精灵那波人搞的

影刀 # 国产 UiPath

Power Automate # 微软出的

Tag UI # https://github.com/aisingapore/TagUI 不维护了

PyAutoGUI # Python 第三方库，基于 opencv、 等工具实现的自动化工具，可以识别图像并调用鼠标和键盘操作这些识别到的图像。

# 其他

[RPA 中国官网](http://www.rpa-cn.com/) 居然没有 HTTPS？
