---
title: YAML
linkTitle: YAML
weight: 20
---

# 概述

> 参考：
>
> - [官方文档，规范 v1.2.2](https://yaml.org/spec/1.2.2/)
> - [Wiki, YAML](https://en.wikipedia.org/wiki/YAML)

编程免不了要写配置文件，怎么写配置也是一门学问。

YAML 是专门用来写配置文件的语言，非常简洁和强大，远比 JSON 格式方便

**YAML Ain't Markup Language(简称 YAML)** 是一种数据序列化语言。设计目标就是方便人类读写，并且可以在日常工作中与现代编程语言很好的配合。它实质上是一种通用的数据串行化格式。

## YAML 与 JSON 的关系

JSON 和 YAML 都旨在成为人类可读的数据交换格式。但是，JSON 和 YAML 具有不同的优先级。 JSON 的首要设计目标是简单性和通用性。因此，JSON 的生成和解析非常简单，但代价是人类可读性降低。它还使用最低公分母信息模型，以确保所有现代编程环境都可以轻松处理任何 JSON 数据。

相反，YAML 的首要设计目标是人类可读性并支持序列化任意本机数据结构。因此，YAML 允许可读性极强的文件，但生成和解析起来更复杂。此外，YAML 的业务范围超越了最低公分母数据类型，因此在不同的编程环境之间进行转换时，需要进行更复杂的处理。

因此，YAML 可以看作是 JSON 的自然超集，可以提高人类可读性和更完整的信息模型。实际上也是这种情况；每个 JSON 文件也是一个有效的 YAML 文件，JSON 与 YAML 格式可以轻松得互相转换

并且，YAML 格式也可以转换为别的格式

# YAML 基本语法规则

- 大小写敏感
- 使用缩进表示层级关系
  - 缩进时不允许使用 Tab 键，只允许使用空格。
  - 缩进的空格数目不重要，只要相同层级的元素左侧对齐即可
- `#` 表示注释，从这个字符一直到行尾，都会被解析器忽略

# Data Structures(数据结构)

YAML 由多个 **Node(节点)** 组成，每个 Node 都可以是三种 **Native Data Structures(原生数据结构)** 其中之一：

- **Scalars(标量)** # 单个的、不可再分的值，又称为 Strings(字符串)、Numbers(数字)
- **Mappings(映射)** # 键值对的集合，又称为 Object(对象)、哈希(Hashes)、字典(Dictionarys)。转为 json 后使用 `{}` 符号包围。
  - **很多使用 YAML 作为配置的产品的官方文档中，都用 Object(对象) 来表示 Mappings**
- **Sequences(序列)** # 一组按次序排列的值，又称为 arrays(数组)、lists(列表)。转为 json 后使用 `[]` 符号包围。对于序列，人们更多得使用数组来称呼这种原语

上述三种类型的 Node 又可以互相组合形成复合结构：

- 由 mappings 组合成的一种复杂结构，官方文档描述为 **Mapping of Mappings**(意味在一个 Mapping 中有多个 Maapings)。
- 由 Scalars 和 Object 组成的一种复杂结构，很多产品的文档中会看到 `[]Object` 符号，这表示，这个字段下的子字段，是序列和映射的组合体，不是单一的数据结构。`[]Object` 分开看，就是 `[]` 与 Object，而 Object 是用 `{}` 表示。
- 等等

这些概念，与各种编程语言中的数组、映射概念相同。

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/actt7h/1670853760782-b2b178b4-bfcf-42bb-b7df-e6827a16233a.svg)

## Node

> 参考：
>
> - [官方规范，Nodes](https://yaml.org/spec/1.2.2/#3211-nodes)

**其实本人更喜欢将 Node 称为 Field(字段)。。。囧。。。**

### Scalar(标量)

YAML 中标量通常不能独自存在，一般都会在 Mapping 或 Sequence 中作为其中的一部分。

### Mapping(映射) 或 Object(对象)

**Mapping(映射) 或 Object(对象)** 表示一个或多个无序的**键/值对**。这里面的 **值** 即可以是**任意类型的节点**

- 使用 `:`(冒号和空格) 分隔每个键值对
- 映射使用键值的方式表示一个数据(比如“名字: DesistDaydream”，“名字” 是一个数据名，“DesistDaydream” 是该数据的值)
- 映射的一组键值对使用冒号结构表示，并且冒号后面需要跟一个空格，否则该行代码就属于 Scalar 类型的节点

这是一个最基本的 Mapping，Key 的值是 Scalar 类型的节点

```yaml
key1: value1
key2: value2
```

Key 的值可以是一个 Mapping 类型的节点

```yaml
key1:
  child-key1: value
  child-key2: value
```

Key 的值可以是一个 Sequence 类型的节点

```yaml
key1:
  - element1
  - element2
```

### Sequence(序列)

**Sequences(序列)** 类型的节点是具有零个或多个元素的数组，数组中的**每个元素可以是任意类型的节点**

- 使用 `-`(波折号和空格) 来表示数组中的每个元素。
- 相同缩进的 `-` 表示同一个数组中的元素。
- 如果一个对象要定义多个，那么该对象的子对象的第一个前面要加-，比如：如果要定义多个 container 字段，这个 container 包括 name，image，ports 等子字段，那么就需要在 container 对象下一行添加一个-，这样，每个-至下一个-中间的内容表示一个 container 的各种规格，第二个-至第三个-中间表示第二个 container 的各种规则

这是一个没有元素的 Sequences 类型的节点

```yaml
-
```

这是一个简单的 Sequence，有两个元素，这俩元素是 Scalar 类型的节点

```yaml
- element
- element
```

Sequence 甚至可以包含 Sequence，也就是说元素是一个 Sequence 类型的节点

```yaml
-
  - element
  - element
- element
```

Sequence 中的元素还可以是 Mapping

### 复合 Node

```yaml
containers:
  - args:
      - AAA
      - BBB
    name: XXX
    image: XXX
  - name: YYY
    image: YYY
```

# Scalar

纯量是最基本的、不可再分的值。以下数据类型都属于 YAML 的纯量。

- 字符串
- 布尔值
- 整数
- 浮点数
- Null
- 时间
- 日期

```yaml
boolean:
    - TRUE  #true,True都可以
    - FALSE  #false，False都可以
float:
    - 3.14
    - 6.8523015e+5  #可以使用科学计数法
int:
    - 123
    - 0b1010_0111_0100_1010_1110    #二进制表示
null:
    nodeName: 'node'
    parent: ~  #使用~表示null
string:
    - 哈哈
    - 'Hello world'  #可以使用双引号或者单引号包裹特殊字符
    - newline
      newline2    #字符串可以拆成多行，每一行会被转化成一个空格
date:
    - 2018-02-17    #日期必须使用ISO 8601格式，即yyyy-MM-dd
datetime:
    -  2018-02-17T15:02:31+08:00    #时间使用ISO 8601格式，时间和日期之间使用T连接，最后使用+代表时区
```

## Block literal(文字块)

> 参考：
>
> - https://yaml.org/spec/1.2.2/#chapter-8-block-style-productions

这是 YAML 的亮点，特别是与 XML 相比，它的 CDATA 显得相当简陋，Block literal 可以将大块文本细致地插入文件中，常用来在传递配置文件时使用(比如通过 k8s 的 configmap 来传递某个程序的配置文件)

文字块的处理总体分为两大类

- `|` # 保留多行文本
- `>` # 将多行文本合并成一行

使用 `|` 符号在文本中保留行，如：

```yaml
text: |
  This is a really long text
  that spans multiple lines (but preserves new lines).
  It does not need to be escaped with special brackets,
  CDATA tags, or anything like that
```

YAML 编译器将会从第一行的第一个文本字符开始编译（并丢掉所有的缩进空格），但是会在你的文本中保留新行。

对文字块的末尾换行符的处理，一共有 3 中写法

- `|-` # Strip(除去)，去掉末尾的所有换行符
- `|` # Clip(裁剪)，保留末尾一个换行符
- `|+` # Keep(保留)，保留末尾的所有换行符

使用 `>` 符号将多行合并，如：

```yaml
text: >
  这部分内容
  以及这行
  将会合并成一行
```

# YAML 示例

这是一个关于发票信息的配置信息

```yaml
invoice: 34843
date: 2001-01-23
bill-to: &id001
  given: Chris
  family: Dumars
  address:
    lines: |
      458 Walkman Dr.
      Suite #292
    city: Royal Oak
    state: MI
    postal: 48046
ship-to: *id001
product:
  - sku: BL394D
    quantity: 4
    description: Basketball
    price: 450.00
  - sku: BL4438H
    quantity: 1
    description: Super Hoop
    price: 2392.00
tax: 251.42
total: 4443.52
comments: Late afternoon is best.
  Backup contact is Nancy
  Billsmer @ 338-4338.
```

这是一个日志信息的基本示例：

```yaml
---
Time: 2001-11-23 15:01:42 -5
User: ed
Warning: This is an error message
  for the log file
---
Time: 2001-11-23 15:02:31 -5
User: ed
Warning: A slightly different error
  message.
---
Date: 2001-11-23 15:03:17 -5
User: ed
Fatal: Unknown variable "bar"
Stack:
  - file: TopClass.py
    line: 23
    code: |
      x = MoreObject("345\n")
  - file: MoreClass.py
    line: 58
    code: |-
      foo = bar
```

# YAML 数据格式最佳实践

YAML 有两种格式

- Document 格式 yaml 数据。也称为 格式化之后的数据、人类可读类数据 等等
  - 一般作为书面人类可读的格式使用，通过缩进、- 符号来规范格式。一般有多行
- 非 Document 格式 yaml 数据。也称为 格式化之前的数据、人类不可读数据 等等。这种格式与 JSON 格式基本一致
  - 一般用于在代码中传值使用，使用 {} 和 \[] 之类的符号来规范格式。一般只有一行

## 映射

映射的一组键值对，使用冒号结构表示

```yaml
key: value
```

转为 JSON 如下

```json
{ "key": "value" }
```

### 对象(映射的复合结构)

```yaml
object:
  name: Steve
  foo: bar
```

转为 JSON 如下

```json
{"object":{"name":"Steve","foo":"bar"}}{object: { name: 'Steve', foo: 'bar' } }
```

### 注意

映射下不可以有子字段，比如下面这种结构是错误的写法

```yaml
key: value
  keyerror: valueerror
```

## 数组

一组连词线开头的行，构成一个数组。

```yaml
- A
- B
- C
```

转为 JSON 如下

```json
["A", "B", "C"]
```

数据结构的子成员是一个数组，则可以在该项下面缩进一个空格。

```yaml
- - A
  - B
  - C
```

转为 JSON 如下。

```json
[["A", "B", "C"]]
```

数组也可以采用行内表示法。

```yaml
A:
  - B
  - C
```

转为 JSON 如下。

```json
{ "A": ["B", "C"] }
```

### 映射与数组得复合结构

映射和数组可以结合使用，形成复合结构。

```yaml
languages:
  - Ruby:
      description: unknown
      type: XXX
  - Perl
  - Python:
      description: Scripts Programming
      type: scripts
websites:
  YAML: yaml.org
  Ruby: ruby-lang.org
  Python: python.org
  Perl: use.perl.org
```

转为 JSON 如下。

```json
{
  "languages": [
    {
      "Ruby": {
        "description": "unknown",
        "type": "XXX"
      }
    },
    "Perl",
    {
      "Python": {
        "description": "Scripts Programming",
        "type": "scripts"
      }
    }
  ],
  "websites": {
    "YAML": "yaml.org",
    "Ruby": "ruby-lang.org",
    "Python": "python.org",
    "Perl": "use.perl.org"
  }
}
```

## 引用

& 锚点和 * 别名，可以用来引用:

```
defaults: &defaults
  adapter:  postgres
  host:     localhost

development:
  database: myapp_development
  <<: *defaults

test:
  database: myapp_test
  <<: *defaults
```

相当于:

```
defaults:
  adapter:  postgres
  host:     localhost

development:
  database: myapp_development
  adapter:  postgres
  host:     localhost

test:
  database: myapp_test
  adapter:  postgres
  host:     localhost
```

**&** 用来建立锚点（defaults），**<<** 表示合并到当前数据，***** 用来引用锚点。

下面是另一个例子:

```yaml
- &showell Steve
- Clark
- Brian
- Oren
- *showell
```

转为 JavaScript 代码如下:

```js
[ 'Steve', 'Clark', 'Brian', 'Oren', 'Steve' ]
```
