---
title: Control structure
linkTitle: Control structure
weight: 1
---

# 概述

> 参考:
>
> - [官方文档，参考 - 语言规范 - 语句](https://go.dev/ref/spec#Statements)

Go 语言提供了 4 种条件结构和分支结构用作 [Control structure(控制结构)](/docs/2.编程/计算机科学/Control%20structure.md)。在结构中，可以使用 `break` 和 `continue` 这样的关键字来中途改变结构的状态。还可以使用 `return` 来结束某个函数的执行，或使用 `goto` 和标签来调整程序的执行位置

> Tips: 在 Go 语言种，将这种控制结构描述为 flow of control

# for 循环

> 参考：
>
> - [Go 官方文档，参考 - 语言规范 - For 语句](https://go.dev/ref/spec#For_statements)
> - [知乎，Golang那些坑 - 使用 for 循环的注意事项](https://juejin.cn/post/7153633858309586975)

用于测试某个条件(布尔型或逻辑型)的语句，初始化语句执行完成之后；如果该条件成立，则会执行 if 后由大括号括起来的代码块，然后执行修饰语句，之后再次判断条件语句是否成立，如此循环；直到条件语句不成立时，就忽略该代码块继续执行后续的代码。

- 基本格式：`for 初始化语句;条件语句;修饰语句 {代码块}`

## for range

> [!Warning]
> 注意 for range 的 [内存逃逸](/docs/2.编程/高级编程语言/Go/Go%20常见问题/内存逃逸.md) 问题

# if else 判断

# switch 判断

给定一个变量，当该变量满足某个条件时执行某个代码。

# select

与 switch 类似

# break 与 continue

- `break` 用于退出当前当前代码块
- `continue` 用于忽略当前循环，继续执行后续循环，只用于 for 结构体中 Note：注意！是退出当前代码块，如果循环有多层嵌套，那么只是退出当前循环；如果循环中套用 select 等，则也是退出当前控制结构。

# 标签与 goto

- 标签用于在出现标签关键字的时候，代码回到标签定义行再继续执行下面的代码。某一行以`:`冒号结尾的单词即可定义标签。标签区分大小写
