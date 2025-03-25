---
title: Map AND Struct
linkTitle: Map AND Struct
weight: 8
---

# 概述

> 参考：
>
> - [官方文档，参考 - 规范 - Map 类型](https://go.dev/ref/spec#Map_types)
> - [官方文档，参考 - 规范 - Struct 类型](https://go.dev/ref/spec#Struct_types)

# Map(映射)

map 是 **key-value pairs(键值对)** 的无序集合。这种结构也称 **关联数组**(associative array)、**字典**(dictionary)、**散列表/哈希表**(hash table)。这是一种快速寻找值的理想结构：给定 Key，对应的 Value 可以迅速定位。

## map 的声明

```go
var MapID map[KeyType]ValueType
```

## map 的实例化

```go
MapID = make(map[KeyType]ValueType)
```

### 赋值

```go
MapID[KEY] = VAL
```

### 实例化的同时进行赋值

```go
MapID := map[KeyType]ValType{
    KEY_1:VAL_1,
    KEY_2:VAL_2,
    ...,
    KEY_n:VAL_n
}
```

这相当于：

```go
MapID := make(map[KeyType]ValueType)
MapID["KEY_1"] = "VALUE_1"
MapID["KEY_2"] = "VALUE_2"
```

## map 的引用

下面是引用 map 中指定 KEY 的 VALUE 的方法：

```go
MapID[KEY]
```

### 引用 map 的长度

map 的长度指的是键值对的个数，有几个键值对，长度就是几

```go
len(MapID)
```

引用 map 中某个 Key 的 Value

```go
MapID["KEY"]
```

## Key/Value 的删除

```go
MapID["KEY"]
```

删除 MapID 这个 map 的 KEY 以及对应的 VAL

# Struct(结构体)

**Struct(结构体)** 是一种**复合[Data type](/docs/2.编程/高级编程语言/Go/Go%20规范与标准库/Data%20type.md)(数据类型)**，Struct 可以看作是一个自定义的 Data type，由一系列的 **Fields(字段，有的地方也称为属性)** 组成，每个字段通常都包含 **名称** 和 **类型**。Struct 把数据聚集在一起，然后访问这些数据的时候，好像这些数据是一个独立实体的一部分。Struct 也是值类型，可以通过 `* new()` 函数创建。

组成结构体的属性分两部分：

- **FieldName(字段名称)** # 字段的名称
- **BaseType(基础类型)** # 基础类型可以是另一个结构体,表示该结构体包含另一个结构体

每个字段都有其对应的基础数据类型，在一个结构体中，FIELD 名字必须是唯一的。代码示例：struct.go

## Struct 的定义

```go
type StructID struct {
    Name1 BaseType1 ["TAG"]
    Name2 BaseType2 [`TAG`]
    ...
}

// 也可以使用简单的方法定义一个结构体:
// 这个结构体中，两个属性的类型都是 int
type T struct {a, b int}
```

## Struct 的声明

```go
var StructVarID StructID
```

## Struct 的实例化

```go
StructVarID := new(StructID)
```

## Struct 的引用

**结构体中属性的引用**
结构体名，中间跟一个点，再接该结构体内的字段名。即可引用该结构体中的某个字段

```go
StrcutID.FIELD1
```

在 Go 语言中，这个 `.` 点符号叫做 **Selector(选择器)**。无论定义的变量是一个结构体类型还是一个结构体类型指针，都是用同样的 **selector-notation(选择器符)** 来引用结构体的字段。

## Tag

除了 FIeld 和 BaseType 之外，还可以给该属性添加 **Tag(标签)**，TAG 使用 `双引号` 或者 `重音号` 来表示。这些 Tag 能被用来做文档或者重要的标签（比如在使用 JSON、YAML 等解析时，这些解析库会读取 Tag 中的内容，将结构体中的每个属性对应到 json 或 yaml 的字段上）。

Tag 里面的内容在正常编程中没有作用。一般在 **Reflect(反射)**、某些第三方库(比如 gin 的数据绑定功能)、等等地方可以起到关键的作用。Strcut 中的 Tag 在通常的编程中无法使用，**只有 reflect 包中的方法可以获取这些 Tag 的信息**。

## Struct 常见用法

在日常使用中，我们常常自己定义一个 函数，用于初始化一个接口体

```go
func NewStructID() *StructID {
 return &StructID(
  FIELED1: XXXX
  FIELED2: XXXX
  ......
 )
}
```

然后直接调用该方法即可，使用方式其实与 new() 函数一样

```go
func main() {
 StructVar := NewStructID
}
```

## 最佳实践

### 将一个结构体中的值拷贝到另一个相似的结构体中

> 参考：
>
> - [GitHub 项目，jinzhu/copier](https://github.com/jinzhu/copier)
> - [知乎，Golang 如何优雅得转换两个相似的结构体](https://www.zhihu.com/question/449267385)
> - [公众号，使用copier提高你的工作效率](https://mp.weixin.qq.com/s/yCI7pKSw0wRlT80k8Bd2vQ)

使用 jinzhu/copier 项目即可轻松得将一个结构体中的值复制到另一个结构体中
