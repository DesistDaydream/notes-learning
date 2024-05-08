---
title: Snort Rule body
linkTitle: Snort Rule body
date: 2024-05-06T13:32
weight: 20
---

# 概述

> 参考：
>
> - [官方文档，Rule Options](https://docs.snort.org/rules/options/index.html)

Snort Rule body 由很多的 Options 组成，<font color="#ff0000">Options 是 Snort 规则的核心和灵魂</font>，这些 Options 决定了被 Snort 规则影响的数据包是否应该通过并发往目的地，或是还是应该应该就此停止。

> 其中最关键的就是 Payload Detection(载荷监测) 选项，Payload 检测是对 Payload 进行匹配的核心规则，只要满足了匹配规则，就可以对匹配到的数据包进行一系列操作。

每个 Options 都有一个 **name(名称)**，关键字后面跟 `:`，冒号后面时 Options 的具体内容（也可以成为 Options 的**条件**），最后每个 Options 都要以 `;` 结尾。有的选项带有 **arguments(参数)**，可以在 `;` 之间用 `,` 分割 Options 和 Option arguments。

比如: `content:" pizza", within 6;`

- `content` 是 Option 的名称，表示这是一个名为 content 的选项。
- `"pizza"` 是 Option 的内容
- `, within 6` 是 Option 的参数

Rule Options 总共分为 4 类:

- **General(常规)**
- **Payload Detection(Payload 检测)**
- **Non-Payload Detection(非 Payload 检测)**
- **Post-Detection(监测后)** # 指定的规则触发后要对数据包采取的操作

# 常规

## sid

Snort ID(简称 sid) 是 Snort 规则的唯一标识符

# Payload 检测

https://docs.snort.org/rules/options/payload/

## content

https://docs.snort.org/rules/options/payload/content

Syntax(语法): `content:[!]"content_modifer"[,content_modifer_argument];`

- content_modifer 是要匹配的数据内容
- content_modifer_argument （可选的）修饰符，用来对要匹配的内容进行额外的评估要求。

content 用于执行 基本字符串 和/或 十六进制 模式匹配。十六进制表示法必须使用 `| |` 包裹起来

比如:

- `content:"|61 62 63|"` 表示 Playload 中要包含 abc 这三个字符

## offset, depth, distance, within

这四个 Options 是对匹配内容的修饰符，用以指定如何对 Payload 进行查找，比如: 从 Payload 何处开始查找、查找 Payload 的前多少个字节、etc.

depth 选项允许规则编写者指定在 Snort 数据包或缓冲区中查找指定模式的深度。例如，将 depth 设置为 5 将告诉 Snort 仅查找 Payload 的前 5 个字节内的模式。

## pcre

https://docs.snort.org/rules/options/payload/pcre

pcre 选项可以根据 Perl 兼容的 [Regular Expression(正则表达式)](/docs/8.通用技术/Regular%20Expression(正则表达式).md) 对数据包进行匹配。

> Notes: 由于从性能角度来看，正则表达式的成本相对较高，因此使用 PCRE 的选项还应该至少有一个 [content](#content) 选项，以利用 Snort 的快速模式引擎。
