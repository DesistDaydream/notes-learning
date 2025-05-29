---
title: Filter
linkTitle: Filter
weight: 20
---

# 概述

> 参考：
>

WireShark 有两类过滤器

- DisplayFilters(显示过滤器) # 过滤已捕获的数据包
- 捕获过滤器 # 只捕获符合条件的数据包

在 `视图 - 内部 - 支持的协议` 菜单项中，可以查看当前 WireShark 支持的所有可用于编写过滤表达式的协议关键字和协议中的字段关键（比如 tcp、tcp.port、etc.）。

# 显示过滤器

> 参考：
>
> - [官方 Wiki，显示过滤器](https://wiki.wireshark.org/DisplayFilters)
> - [官方文档，6.3 查看时过滤数据包](https://www.wireshark.org/docs/wsug_html_chunked/ChWorkDisplayFilterSection.html)

在 WireShark GUI 过滤器工具栏中，可以通过 **Filter Expressions(过滤器表达式)** 隐藏不敢兴趣的数据包，仅显示某些特定的包，比如

- 特定协议的包
- 特定字段的包
- 满足某些字段的值进行比较后的包
- 等等等等

通过 `文件 - 导出特定分组` 菜单项将过滤的结果导出到一个新的文件中

## 显示过滤器表达式

> 参考：
>
> - [官方文档，6.4 构建显示过滤器表达式](https://www.wireshark.org/docs/wsug_html_chunked/ChWorkBuildDisplayFilterSection.html)


### 运算符

|    英文    |  符号  | 描述   | 示例                                     | 备注                                                            |
| :------: | :--: | ---- | -------------------------------------- | ------------------------------------------------------------- |
|    eq    | `==` | 完全匹配 | `ip.addr == 192.168.1.1`               |                                                               |
| contains |      | 模糊匹配 | `_ws.col.info contains "Client Hello"` | 仅对 Protocol, field or slice 有效。e.g. ip.add 不可以使用 contains 运算符 |
|   etc.   |      |      |                                        |                                                               |

### EXAMPLE

https://gitlab.com/wireshark/wireshark/-/wikis/DisplayFilters

`!tcp.analysis.flags` # 去掉 Bad TCP 的包

# 捕获过滤器

> 参考：
>
> - [官方 Wiki，捕获过滤器](https://wiki.wireshark.org/CaptureFilters)

点击 “捕获选项” ，打开捕获选项窗口，在下方可以输入捕获过滤器表达式。

![center|1000](https://notes-learning.oss-cn-beijing.aliyuncs.com/wireshark/CaptureFilters_view.png)

## 捕获过滤器表达式
