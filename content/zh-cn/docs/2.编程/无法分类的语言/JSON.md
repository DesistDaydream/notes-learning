---
title: JSON
---

# 概述

> 参考：
>
> - [官方文档](https://www.json.org/json-zh.html)
> - [Wiki, JSON](https://en.wikipedia.org/wiki/JSON)
> - [RFC 8259](https://tools.ietf.org/html/rfc8259)

**JavaScript Object Notation(JS 对象表示法，简称 JSON)** 是一种轻量级的数据交换格式。易于人阅读和编写。同时也易于机器解析和生成。 它基于 JavaScript Programming Language, Standard ECMA-262 3rd Edition - December 1999 的一个子集。 JSON 采用完全独立于语言的文本格式，但是也使用了类似于 C 语言家族的习惯（包括 C, C++, C#, Java, JavaScript, Perl, Python 等）。 这些特性使 JSON 成为理想的数据交换语言。

**JavaScript Object Notation(简称 JSON)** 是一种简单的数据交换格式。从句法上讲，它类似于 JavaScript 的对象和列表。它最常用于Web后端与浏览器中运行的 JavaScript 程序之间的通信，但它也用于许多其他地方。它的主页 json.org 提供了一个清晰，简洁的标准定义。

JSON 建构于两种结构：

- “名称/值”对的集合（A collection of name/value pairs） # 不同的语言中，它被理解为对象（object），映射（mapping），纪录（record），结构（struct），字典（dictionary），哈希表（hash table），有键列表（keyed list），或者关联数组 （associative array）。
- 值的有序列表（An ordered list of values） # 在大部分语言中，它被理解为数组（array）。

JSON 具有以下这些形式：

1. 映射是一个无序的“‘名称/值’对”集合。一个对象以 `{` 左大括号 开始， `}`右大括号 结束。每个“名称”后跟一个 `:` (冒号)，“‘名称/值’ 对”之间使用 `,`(逗号)分隔。
2. 数组是值（value）的有序集合。一个数组以`[`左中括号 开始，`]`右中括号 结束。值之间使用 ,逗号 分隔。

映射格式样例

```json
 {
  "name": "desistdaydream",
  "age": 30
}
```

数组格式样例

```json
 [
  "desistdaydream",
  "zhangna"
]
```

混合格式样例：

```json
{
  "family": [
    {
      "name": "desistdaydream",
      "age": 30
    },
    {
      "name": "zhangna",
      "age": 30
    }
  ]
}
```

# JSON Schema

## 概述

> 参考：
>
> - [GitHub 项目，json-schema-org/json-schema-spec](https://github.com/json-schema-org/json-schema-spec)
> - [官网](https://json-schema.org/)
> - [官方文档，参考 - 理解 JSON Schema - 什么是 schema](https://json-schema.org/understanding-json-schema/about)

JSON Schema 是一种声明性语言，用于定义 JSON 数据的结构和约束。JSON Schema 有点类似 [ASN.1](/docs/2.编程/无法分类的语言/ASN.1.md)、XML Schema、etc.

人们可以使用 JSON Schema 语法定义一系列的预期目标，然后利用 [Validator](#Validator) 验证 JSON 数据是否符合预期。e.g. 字段是否为数字、长度是否为 64、某个字段是否是必须的、etc.

JSON Schema 写出来的内容通常长这样：

```json
{
  // 使用 JSON Schema 的哪个规范版本
  "$schema": "http://json-schema.org/draft-07/schema#",
  // 定义了一个对象类型。就是顶层字段是一个 {}
  "type": "object",
  // 该对象中可以包含 name、age、email 三个字段
  "properties": {
    // name 字段必须是至少 2 个字符的字符串
    "name": {
      "type": "string",
      "minLength": 2
    },
    // age 必须是 0-120 的整数
    "age": {
      "type": "integer",
      "minimum": 0,
      "maximum": 120
    },
    // email 必须是有效的邮箱格式的字符串
    "email": {
      "type": "string",
      "format": "email"
    }
  },
  // name 和 age 是必填字段（i.e. 收到的数据如果没有这俩字段则验证不通过）
  "required": ["name", "age"]
}
```

JSON Schema 的验证能力还可以用在如下场景：

- **文档生成** # 为 API 提供清晰的数据格式说明。可以在 [OpenAPI](/docs/2.编程/API/OpenAPI.md) 中集成 JSON Schema 的定义
- **代码生成** # 可以根据 Schema 自动生成相应的数据模型代码
- **IDE 支持** # 提供自动补全和验证功能

若是按照验证功能分类的话，JSON Schema 的验证能力可以分为如下几类：

- 基础类型验证 # 验证字段是否是某个类型
  - string
  - number
  - integer
  - boolean
  - array
  - object
  - null
- 字符串验证 # 验证字段中的字符串是否满足某个格式
  - minLength/maxLength
  - pattern (正则表达式)
  - format (如 email、date、uri 等)
- 数值验证 # 验证数值是否在某个范围区间
  - minimum/maximum
  - multipleOf
  - exclusiveMinimum/exclusiveMaximum
- 数组验证
  - minItems/maxItems
  - uniqueItems
  - items (数组元素的 Schema)
- 对象验证
  - required (必填属性)
  - properties (属性定义)
  - additionalProperties (额外属性控制)
  - dependencies (属性依赖关系)

JSON Schema 的历史可以追溯到 Kris Zyp 于 2007 年 10 月 2 日向 json.com 提交的第一个 JSON Schema 提案。

## Validator

https://json-schema.org/tools

**Validator(验证器)** 是实现 JSON Schema 规范的工具。此类工具可以轻松地将 JSON Schema 集成到任何规模的项目中。

https://github.com/kaptinlin/jsonschema

https://github.com/invopop/jsonschema

## 参考

https://json-schema.org/understanding-json-schema/reference

基本示例

```json
{
  "title": "Match anything",
  "description": "This is a schema that matches anything.",
  "default": "Default value",
  "examples": [
    "Anything",
    4035
  ],
  "deprecated": true,
  "readOnly": true,
  "writeOnly": false
}
```

### 字段关键字

#### 通用关键字

#### 特定于类型的关键字



## JSON Schema 规范

> 参考：
>
> - [官方文档，规范](https://json-schema.org/specification)

- 2022-12
- 2019-09
- draft-07
- draft-06
- draft-05

[JSON Schema Core](https://json-schema.org/draft/2020-12/json-schema-core.html) 定义了基本规范、通用的关键字及用法

[JSON Schema Validation](https://json-schema.org/draft/2020-12/json-schema-validation.html) 包含各种用于验证功能的关键字及用法

