---
title: TOML
linkTitle: TOML
date: 2024-04-07T08:48
weight: 20
---

# 概述

> 参考：
>
> - [GitHub 项目，toml-lang/toml](https://github.com/toml-lang/toml)
> - [官方文档](https://toml.io/en/latest)
> - [Wiki, TOML](https://en.wikipedia.org/wiki/TOML)
> - [知乎](https://zhuanlan.zhihu.com/p/50412485)
> - [格式对比](https://www.cnblogs.com/sunsky303/p/9208848.html)

**Tom's Obvious, Minimal Language(简称 TOML)** 是一种配置语言，旨在称为一种最小的配置文件结构，并且易于阅读、具有显而易见的语义。

## TOML 规范

- TOML 大小写敏感
- TOML 必须是有效的 UTF-8 编码的 Unicode 文档
- 空白表示 Tab(0x09) 或 空格(0x20)
- 换行表示 LF(0x0a) 或 CRLF(0x0D 0x0A)

## TOML 特点

TOML 的原子单位也是 **Key/Value pair(键值对)**。多个 Key/Value pair 组成一个 Table(表)。

所以，一个 TOML 格式的配置文件，本质上是 **Table(表)** 的集合。

TOML 放弃了括号或缩进的底层原理，而是以 `.` 符号来表示层级关系(实现类似缩进的效果)

## TOML 基本示例

```toml
# This is a TOML document.
title = "TOML Example"

[owner]
name = "Tom Preston-Werner"
dob = 1979-05-27T07:32:00-08:00 # First class dates

[database]
server = "192.168.1.1"
ports = [ 8000, 8001, 8002 ]
connection_max = 5000
enabled = true

[servers]
  # 可以使用缩进让结构更清晰，但是并不是必须要缩进的
  [servers.alpha]
  ip = "10.0.0.1"
  dc = "eqdc10"

  [servers.beta]
  ip = "10.0.0.2"
  dc = "eqdc10"

[clients]
data = [ ["gamma", "delta"], [1, 2] ]

# 在数组内部，是可以使用换行符的，主要是为了方便人类阅读
hosts = [
  "alpha",
  "omega"
]
```

### 转换成 JSON 后是这样的

```json
{
  "title": "TOML Example",
  "owner": {
    "name": "Tom Preston-Werner",
    "dob": {
      "date": "1979-05-27 07:32:00",
      "timezone_type": 1,
      "timezone": "-08:00"
    }
  },
  "database": {
    "server": "192.168.1.1",
    "ports": [8000, 8001, 8002],
    "connection_max": 5000,
    "enabled": true
  },
  "servers": {
    "alpha": {
      "ip": "10.0.0.1",
      "dc": "eqdc10"
    },
    "beta": {
      "ip": "10.0.0.2",
      "dc": "eqdc10"
    }
  },
  "clients": {
    "data": [
      ["gamma", "delta"],
      [1, 2]
    ],
    "hosts": ["alpha", "omega"]
  }
}
```

## TOML 与 INI、JSON、YAML 的对比

配置文件是一种非常基础的文件格式，但远没有数据文件格式（如 `SQLite`）、文档文件格式（如 `Markdown`）、编程语言（如 `JavaScript`）、甚至二进制文件格式（如 `PNG`）需求那么复杂。

只要严谨但不严苛、支持必要的数据类型和嵌套，又易于人类手工直接阅读和编辑就可以了。

但就是这样一种广泛需要而又简单的应用场景，却反而长期以来一直没有一种足够好的文件格式。

---

INI（`.ini`）文件是一种非常原始的基础形式，但各家有各家的用法，而且它最多只能解决一层嵌套。只适合非常非常简单的配置文件，一旦需要两层嵌套，或需要数组，就力不从心了。

```ini
; 最简单的结构
a = a;
b = b; 这些等号后面的值是字符串（句末分号不是必须的；它后面的都是注释）
; 稍微复杂一点的单层嵌套结构
[c]
x = c.x
y = c.y
[d]
x = d.x
y = d.y
```

---

JSON（`.json`）是一种非常好的数据存放和传输的格式，但阅读和编辑它实在不方便。即便 `JSON5`（`.json5` - `ECMAScript 5.1 JSON`）这种扩展格式允许了你像写 `JavaScript` 对象那样书写裸键名、允许尾逗号，并且可以有注释，写多行字符串依然麻烦。即便它将来加上了多行字符串语法，依然不行，因为它虽然是基于括号嵌套语法的层级关系，在不缩进的情况下，却根本没法阅读。

```json
{
  "a": "a",
  "b": "b",
  "c": {
    "x": "c.x",
    "y": "c.y"
  },
  "d": {
    "x": "d.x",
    "y": "d.y"
  },
  "e": [
    { "x": "e[0].x", "y": "e[0].y" },
    { "x": "e[1].x", "y": "e[1].y" }
  ]
}
```

---

YAML（`.yaml` 或 `.yml`）干脆将 `JSON` 中有了不够、没有不行的括号结构去掉了，只保留缩进。但编辑和阅读它总令人非常慌张，生怕数错了层次（实际上，对于阅读，语法关键字并不是越小越好）。而且在不支持统一缩进、反缩进、自动在换行时缩进的编辑环境下，这非常麻烦——这本来对编程语言来说不是什么事，但配置文件最常用的使用场景却恰恰是这样。
　　另外，`YAML` 的语法实在太多了，而且不是循序渐进的，即便你不需要复杂的功能，为了保证自己的简单功能不出错，也要对那些复杂的语法有所了解并加以避免（比如究竟什么键名可以不加引号，什么字符串可以不加引号；你总不能为了避免歧义全都加上引号，那和 `JSON` 也就差球不多了）。更糟的是，纵使如此复杂，想要配置一段精确的多行字符串（精确控制首尾空行数）时，却显得力不从心。再加上缩进问题，编辑多行文本实在烦不胜烦。如果你还需要转义字符……

```yaml
a1: abc # string
a2: true # boolean
b1: nil # string
b2: null # null
b3: NULL # null
b4: NuLL # string
b5: Null # null
c:
  x: c.x
  y: c.y
d:
  x: d.x
  y: d.y
e:
  - x: e[0].x
    y: e[0].y
  - x: e[1].x
    y: e[1].y
```

---

终于，TOML（`.toml`）横空出世。它彻底放弃了括号或缩进的底层原理，而是采取了显式键名链的方式。

为了方便（同时看起来更清楚——这种读和写的契合非常关键！），你可以指定小节名。妙的是，小节名也是可以链式声明的。

另外，某些数据可能使用内联数组或表更合适以避免臃肿，这也是支持的。

```toml
a = "a"
b = "b"
c.x = "c.x"
c.y = "c.y"
[d]
x = "d.x"
y = "d.y"
[[e]]
x = "e[0].x"
y = "e[0].y"
[[e]]
x = "e[1].x"
y = "e[1].y"
[f.A]
x.y = "f.A.x.y"
[f.B]
x.y = """
f.
B.
x.
y
"""
[f.C]
points = [
{ x=1, y=1, z=0 },
{ x=2, y=4, z=0 },
{ x=3, y=9, z=0 },
]
```

# TOML 原语

## Key/Value pair(键/值对)

TOML 文档的主要结构也是 `Key/Value pair(键/值对)` 格式。key 与 value 以 `=` 符号分割

### Array(数组)

Array 类型的值以 `[]` 表示，每个元素以 `,` 分割

## Table(表)

> 类似于 INI 中的 Sections(部分)

**Table(表)** 是 `键值对` 的集合，也称为 Hash Tables(哈希表) 或 Dictionaries(字典)，以 `[]` 符号表示。从 Table 的 `[]` 符号开始到下一个 `[]` 符号为止，所有键值对都属于该 Table。

Table 的名称则用 `[]` 符号内的字符串表示。Table 的命名规则与 Key 的命名规则相同，同样是可以使用 `.` 符号来表示 Table 与 Table 之间的层级关系。

配置文件的开头没有任何 `[]` 表示的部分，也称为 **root Table(根表)**。根表不用 \[] 符号，也就没有名称，所有**属于根表的 Key/Value pair 都只能写在文件开头**。

### 示例

**Table 1**

```toml
[table-1]
key1 = "some string"
key2 = 123

[table-2]
key1 = "another string"
key2 = 456
```

转为 JSON：

```json
{
  "table-1": {
    "key1": "some string",
    "key2": 123
  },
  "table-2": {
    "key1": "another string",
    "key2": 456
  }
}
```

**Table 2**

```toml
[dog."tater.man"]
type.name = "pug"
```

转为 JSON：

```json
{ "dog": { "tater.man": { "type": { "name": "pug" } } } }
```

### Inline Tables(内联表)

比如：

```text
name = { first = "Tom", last = "Preston-Werner" }
point = { x = 1, y = 2 }
animal = { type.name = "pug" }
```

表示：

```toml
[name]
first = "Tom"
last = "Preston-Werner"

[point]
x = 1
y = 2

[animal]
type.name = "pug"
```

### Array of Tables(表的数组)

Table 的数组使用 `[[]]` 符号表示。
比如，下面的配置：

```toml
[[products]]
name = "Hammer"
sku = 738594937

[[products]]  # empty table within the array

[[products]]
name = "Nail"
sku = 284758393

color = "gray"
```

转换为 JSON 为：

```json
{
  "products": [
    { "name": "Hammer", "sku": 738594937 },
    {},
    { "name": "Nail", "sku": 284758393, "color": "gray" }
  ]
}
```

# 总结

从某种成都上来说，TOML 也可以使用类似系统中的路径格式来表示，不管是 Table 还是 `.` 符号，这些原语组合成一个 Key 并确定唯一一个值，非常像 Kubernetes 在 Etcd 中存储的数据格式。
比如前文的[基本示例](#TOML%20基本示例)中的所有 Key，可以看成下面这个样子：

```text
/title
/owner/name
/owner/dob
/database/server
/database/ports
/database/connection_max
/database/enabled
/servers/alpha/ip
/servers/alpha/dc
/servers/beta/ip
/servers/beta/dc
/clients/data
/clients/hosts
```

