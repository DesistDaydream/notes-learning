---
title: Variable
weight: 3
---

# 概述

> 参考：
>
> - [Go 官方文档，参考 - 规范 - 变量](https://go.dev/ref/spec#Variables)
> - [Go 官方文档，参考 - 规范 - 声明和范围](https://go.dev/ref/spec#Declarations_and_scope)

**Variable(变量)**

- How to Name a Variable(如何命名一个变量)
- Scope(范围)
- Constants(常量)
- Defining Multiple Variables 定义多个变量空白标识符`_`用于抛弃值，e.g.值 5 在：`_, b = 5, 7`中被抛弃。`_`空白标识符是一个只写变量，不能获得它的值。这样做是因为 Go 语言中必须使用所有被声明的变量，但有时候并不需要使用从一个函数得到的所有返回值(e.g.上例中的 5 与 7 是通过某个函数获得的值且该函数一定会获得 2 个值，但是其中一个却用不上)。

# 声明变量

格式：`var VarID TYPE = EXP`

- VarID # 变量的标识符
- TYPE # 详见 [Data type](/docs/2.编程/高级编程语言/Go/Go%20规范与标准库/Data%20type.md)
- EXP # 初始化时使用的表达式。i.e.给该变量一个值。

其中 `TYPE` 或者 `= EXP` 这两个部分可以省略其中之一，如果省略 TYPE，那么将根据初始化 EXP 来自动推导变量的类型；如果初始化的 EXP

# 引用变量

# Variables Scope(变量范围)

Variables Scope(变量范围) 就是变量的作用域，定义在哪里的变量，可以在哪里被引用，不可以在哪里被引用，都是变量范围所决定的。
