---
title: Filter
linkTitle: Filter
weight: 20
---


# 概述

> 参考：
>
> - [官方文档，6.3 查看时过滤数据包](https://www.wireshark.org/docs/wsug_html_chunked/ChWorkDisplayFilterSection.html)
> - [官方文档，6.3 构建显示过滤器表达式](https://www.wireshark.org/docs/wsug_html_chunked/ChWorkBuildDisplayFilterSection.html)

在 WireShark GUI 过滤器工具栏中，可以通过 **Filter Expressions(过滤器表达式)** 隐藏不敢兴趣的数据包，仅显示某些特定的包，比如

- 特定协议的包
- 特定字段的包
- 满足某些字段的值进行比较后的包
- 等等等等

通过 `文件 - 导出特定分组` 菜单项将过滤的结果导出到一个新的文件中

在 `视图 - 内部 - 支持的协议` 菜单项中，可以查看当前 WireShark 支持的所有可用于编写过滤表达式的协议关键字和协议中的字段关键（比如 tcp、tcp.port、etc.）。

# Syntax(语法)

# EXAMPLE

https://gitlab.com/wireshark/wireshark/-/wikis/DisplayFilters

`!tcp.analysis.flags` # 去掉 Bad TCP 的包
