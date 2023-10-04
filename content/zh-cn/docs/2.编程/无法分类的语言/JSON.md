---
title: JSON
---

# 概述

> 参考：
> 
> - 官方文档：<https://www.json.org/json-zh.html>
> - [Wiki，JSON](https://en.wikipedia.org/wiki/JSON)
> - [RFC 8259](https://tools.ietf.org/html/rfc8259)

**JavaScript Object Notation(JS 对象表示法，简称 JSON)** 是一种轻量级的数据交换格式。易于人阅读和编写。同时也易于机器解析和生成。 它基于 JavaScript Programming Language, Standard ECMA-262 3rd Edition - December 1999 的一个子集。 JSON 采用完全独立于语言的文本格式，但是也使用了类似于 C 语言家族的习惯（包括 C, C++, C#, Java, JavaScript, Perl, Python 等）。 这些特性使 JSON 成为理想的数据交换语言。

**JavaScript Object Notation(简称 JSON)** 是一种简单的数据交换格式。从句法上讲，它类似于 JavaScript 的对象和列表。它最常用于Web后端与浏览器中运行的 JavaScript 程序之间的通信，但它也用于许多其他地方。它的主页 json.org 提供了一个清晰，简洁的标准定义。

JSON 建构于两种结构：

- “名称/值”对的集合（A collection of name/value pairs） # 不同的语言中，它被理解为对象（object），映射（mapping），纪录（record），结构（struct），字典（dictionary），哈希表（hash table），有键列表（keyed list），或者关联数组 （associative array）。
- 值的有序列表（An ordered list of values） # 在大部分语言中，它被理解为数组（array）。

JSON 具有以下这些形式：

1. 映射是一个无序的“‘名称/值’对”集合。一个对象以 `{`左大括号 开始， `}`右大括号 结束。每个“名称”后跟一个 :冒号 ；“‘名称/值’ 对”之间使用 ,逗号 分隔。
2. 数组是值（value）的有序集合。一个数组以`[`左中括号 开始，`]`右中括号 结束。值之间使用 ,逗号 分隔。

映射格式样例

```json
 {
  "name": "lichenhao",
  "age": 30
}
```

数组格式样例

```json
 [
  "lichenhao",
  "zhangna"
]
```
混合格式样例：

```json
{
  "family": [
    {
      "name": "lichenhao",
      "age": 30
    },
    {
      "name": "zhangna",
      "age": 30
    }
  ]
}
```
