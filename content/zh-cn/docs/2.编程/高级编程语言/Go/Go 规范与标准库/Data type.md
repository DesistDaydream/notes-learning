---
title: Data type
weight: 2
---

# 概述

> 参考：
>
> - [官方文档，参考-规范-类型](https://go.dev/ref/spec#Types)

**Data Type(数据类型)** 用来对一组相关值进行分类，描述可对其执行的操作并定义它们的存储方式。 通常也会称为 **Literal(字面量)**

Go 语言将数据类型分为四类：基础类型、复合类型、引用类型和接口类型。虽然数据类型有很多，但是这些数据类型都是对程序中一个变量或状态的间接引用。这意味着对任一引用类型数据的修改都会影响所有该引用的拷贝。所谓的引用，是对值的引用。注意引用与指针的区别，详见 pointer.go

## Underlying Type(基本类型)

数据中最基本的类型，是构成其余数据类型以及对象的最小单位，当定义其他数据类型时，同样需要定义基础数据类型。基础数据类型也是 Go 语言的**内置数据类型**

- **Numeric(数字类型)**
  - Integer Type(整数类型)
  - Floating Point Numbers(浮点数型)
  - complex 复数共两种复数，complex64 和 complex128，分别对应 float32 和 float64 两种浮点数精度。内置的 complex 函数用于构建复数，内建的 real 和 imag 函数分别返回复数的实部和虚部
- **Strings(字符类型)**
- **Booleans(布尔类型)**
  - 注意：很多时候，Booleans 类型的值可以用数字表示
    - 1 表示 true(真)
    - 0 表示 false(假)

## Composite Type(复合类型)

是通过组合基础数据类型，来表达更复杂的数据结构

- Arrays(数组) # 多个相同基础类型的数据组合在一起
- Slices(切片)
- Maps(字典)
- Functions(函数，这里面主要指的是函数的参数的数据类型)
- Structs(结构体) #
- Interfaces(接口)
- Channels(通道)
- Pointers(指针)

## 自定义数据类型(类型定义)

变量或表达式的 Type 定义了对应存储值的属性特征，例如数值在内存的存储大小（或者是元素的 bit 个数），它们在内部是如何表达的，是否支持一些操作符，以及它们自己关联的方法集等。

在任何程序中都会存在一些变量有着相同的内部结构，但是却表示完全不同的概念。例如，一个 int 类型的变量可以用来表示一个循环的迭代索引、或者一个时间戳、或者一个文件描述符、或者一个月份；一个 float64 类型的变量可以用来表示每秒移动几米的速度、或者是不同温度单位下的温度；一个字符串可以用来表示一个密码或者一个颜色的名称。这些基于基本数据类型所生成的新数据类型都叫数据类型。再比如数组、切片、字典等，虽然在有的时候他们的基础数据类型可以使一样的，但是他们本身所表示的数据类型是不同的含义。

一个类型声明语句创建了一个新的类型名称，和现有类型具有相同的底层结构。新命名的类型提供了一个方法，用来分隔不同概念的类型，这样即使它们底层类型相同也是不兼容的。

格式：`type TypeID BaseType`

如果自定义数据类型的 BaseType 是由零个或多个任意类型的值(每个值对应一个类型)聚合成的实体，则需要使用 [Struct](/docs/2.编程/高级编程语言/Go/Go%20 标准库/Maps(映射)%20 与%20Struct(结构体).md 与 Struct(结构体).md)

**Type Definition(类型定义)** 是 Go 实现面向对象编程的最基本要素

# 类型检测

详见 [Reflection](/docs/2.编程/高级编程语言/Go/Go%20规范与标准库/Reflection.md) 特性

# 类型转换

> 参考：
>
> - [官方 Tour，基础-13](https://go.dev/tour/basics/13)
> - [标准库，strconv](https://pkg.go.dev/strconv)

语法：`Type(Expression)`

- 将 Expression 获取到的值转换为 Type 类型。

> 注意：这种语法属于强制类型转换，若出现错误将会 panic。如果想要处理错误，可以使用 [strconv](https://pkg.go.dev/strconv) 库等

注意：若 Type 为指针类型，最好写成这样：`(Type)(Expression)`，比如：

```go
 type stringValue string
 p := new(string)
 a := (*stringValue)(p)
 fmt.Println(a)
```

## \[]byte 与 String 互相转换

```go
	// string 转 []byte
	str := "hello"
	bytes := []byte(str)

	// []byte 转 string
	str2 := string(bytes)
```

## strconv

strconv 全称 **string conversion(字符串转换)**，可以实现字符串类型与其它[Underlying Type(基本类型)](#Underlying%20Type(基本类型)) 之间的转换。