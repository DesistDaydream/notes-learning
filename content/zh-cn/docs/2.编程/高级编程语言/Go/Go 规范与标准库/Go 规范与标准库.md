---
title: Go 规范与标准库
linkTitle: Go 规范与标准库
weight: 1
---

# 概述

> 参考：
>
> - [官方文档，参考 - 规范](https://go.dev/ref/spec)
>   - [官方文档，参考 - 规范 的翻译](https://github.com/bekcpear/mypelicanconfandarticles/blob/master/content/Tech/gospec.rst)
> - [go.dev, Tour(Go 语言之旅，通过在线解析器体验 Go 语言的各种特性)](https://go.dev/tour/list)
> - [公众号，11个现代Go特性：用 gopls/modernize 让你的代码焕然一新](https://mp.weixin.qq.com/s/mQehW07uSvfMkSMfcrSEsA)

Go 是一种通用语言，专为系统编程而设计。它是一种强类型且自带垃圾回收功能的语言，并具有显式支持并发编程的能力(称为 goroutine)。Go 程序由 Packages(包) 构建，其属性允许有效得管理依赖关系。

- **Go 语言参考**描述了 Go 语言的具体语法和语义
- **Go 标准库则**是与 Go 语言一起发行的一些可选功能，以便人们可以从一开始就轻松得使用 Go 进行编程。

# Keywords

> 参考：
>
> - [官方文档，参考 - 规范 - 关键字](https://go.dev/ref/spec#Keywords)

Go 语言非常简单，只有 25 个`关键字(Keywords)`可以使用，记住这 25 个关键字，就掌握了最基本的 Go 语言用法。这些关键字是 go 语言保留的，不能用作标识符

`关键字`在编程语言中是指该语言的一个功能，比如下文面的 `var`，就是指声明一个变量，`func` 就是定义一个函数等等。

> Note: if-else 算两个关键字所以在这里一共只写了 24 个。

1. **break** # 控制结构
2. **case** # 控制结构
3. **chan** # 用于 channel 通讯
4. **const** # 语言基础里面的常量申明
5. **continue** # 用在 for 控制结构中，用以忽略本次循环的后续所有逻辑，执行下一次循环
6. **default** # 控制结构
7. **defer** # 用于在函数退出之前执行某语句的功能
8. **fallthrough** # 控制结构
9. **for** # 控制结构
10. **func** # 用于定义函数和方法
11. **go** # 用于并发
12. **goto** 控制结构
13. **if-else** # 控制结构
14. **import** 用于定义该文件引用某个包
15. **interface** # 用于定义接口
16. **map** # 用于声明 map 类型数据
17. **package** # 用于定义该文件所属的包
18. **range** # 用于读取 slice、map、channel 数据
19. **return** # 用于从函数返回。有时候也用来直接跳出当前函数，回到主程序继续执行
20. **select** # 用于选择不同类型的通讯
21. **struct** # 用于定义抽象数据类型
22. **switch** # 控制结构
23. **type** # 用于 Type Declarations(类型声明)，有两种形式：
24. Definitions(定义) 自定义类型
25. Declarations(声明) 一个类型的别名。
    1. 其实所谓的类型的别名，也可以当作一种自定义的类型。
26. **var** # 用于 Declarations(声明) 变量

# Lexical elements(词汇元素)

一些 Go 语言中抽象或具象名词，用于描述某些实体或行为。

## Identifier

**Identifier(标识符)** 是一个抽象的概念，代表已命名的实体（e.g. [Variable](/docs/2.编程/高级编程语言/Go/Go%20规范与标准库/Variable.md)、自定义的 [Data type](/docs/2.编程/高级编程语言/Go/Go%20规范与标准库/Data%20type.md)、etc.）。Identifiers 由一个或多个字母和数字组成，Identifier 的第一个字符必须是字母。

有一些 Indentifiers 是 [predeclared(预先声明的)](https://go.dev/ref/spec#Predeclared_identifiers)（e.g. int, int8, rune, true, false, append, print, new, etc.）这些预声明的 Identifier 是一种类似 Keywords 的存在，可以是 类型、常量、零值、函数.

> [!Note] 随着 Go 语言版本的迭代，会逐渐加入一些新的预声明 Identifier（e.g. 用于快速比较获取最大值/最小值的 max, min 内置函数是在 1.21 版本加入的；1.21 版本后，删除数组中的元素也有了可以直接使用的 slices.Delete() 方法；etc.）

# Notation(表示法)

Go 语言的语法遵从 [EBNF](https://en.wikipedia.org/wiki/Extended_Backus%E2%80%93Naur_form) 表示法

# Block(块)

Block(块) 是由一对 `{}`(大括号) 括起来的一系列声明和语句。Block 可以是空。

除了显式的 Block 外，Go 语言还存在 implicit block(隐式块)

Block 的用法会直接影响 [scoping(作用域)](#Declarations%20and%20scope(声明与范围))

## implicit block(隐式块)

1. **universe block(全域块)** # universe block 代表所有 Go 源代码文本。Notes: 这里的所有指编译时用到的所有 Go 文件。
2. **package block(包块)** # 每个 Package 有一个 package block，是包含该包的所有 Go 源代码文本
3. **file block(文件快)** # 每个文件都有一个 file block(文件块)，包含该文件中的所有 Go 源代码文本。
4. 每个 if, for, switch 语句都被视为处于其自己的隐式代码块中。
5. switch 或 select 语句中的每个子句都作为一个隐式代码块。 代码块可以嵌套，并影响作用域。

# Declarations and scope(声明与范围)

https://go.dev/ref/spec#Declarations_and_scope

这部分是介绍作用域的，用来定义各种 [Identifier](#Identifier) 是否可以被引用、是否可以被使用

## exported and unexported(导出与未导出)

https://go.dev/ref/spec#Exported_identifiers

> Note: **uppercase(大写)** 或 **lowercase(小写)**

[Identifier](#identifier) 可以被 **exported(导出)** 以允许其他[包](#Packages(包))访问 ta。当满足以下条件时，Identifier 将被导出：

- Identifier 名称的第一个字符是 **uppercase** letter(大写字母)
- Identifier 必须在 package block 中声明。在其他地方声明的均不会被导出。
  - Note: package block 是一种 [implicit block(隐式块)](#implicit%20block(隐式块))

# Packages(包)

https://go.dev/ref/spec#Packages

更多介绍见 [Go Module](docs/2.编程/高级编程语言/Go/Go%20环境安装与使用/Go%20Module.md)

## Import(导入)

https://go.dev/ref/spec#Import_declarations

# Go 标准库

> 参考：
>
> - [Go 包，标准库](https://pkg.go.dev/std)
>   - [中文文档](https://studygolang.com/pkgdoc)

**Go Standard Library(Go 标准库)** 是 Go 内置 **Package(包)** 的集合，每个 package 都可以实现一类功能。每个 package 里有他们对应的常量、变量、函数、方法等。每个库就是一类功能，比如 bufio 库，这里面就是关于实现读写功能的各种内容；而 fmt 库则是关于实现格式化输入输出等功能。在[这里](https://pkg.go.dev/std?tab=packages)可以看到 go 语言 原生支持的所有标准库。

与 标准库 相对应的就是 [Go 第三方库](/docs/2.编程/高级编程语言/Go/Go%20第三方库/Go%20第三方库.md) ，第三方库一般属于由个人开发，实现更多丰富功能的库。在 [Go.dev](https://pkg.go.dev/)可以搜索自己想要使用的所有库。

[Go.dev](https://pkg.go.dev/) 是 golang.org 的配套网站。 Golang.org 是开源项目和发行版的所在地，而 go.dev 是 Go 用户的中心，可从整个 Go 生态系统中提供集中和精选的资源。还可以在专门的[标准库](https://pkg.go.dev/std)页面看到所有标准库下的 Go 内置包。

Go.dev 提供：

- 在 index.golang.org 上发布的 Go 软件包和模块的集中信息。
- 基本学习资源
- 关键用例和案例研究

Go.dev 当前处于 MVP 状态。我们为自己的建设感到自豪，并很高兴与社区分享。我们希望您能在使用 go.dev 的过程中找到价值和乐趣。 Go.dev 只有一小部分我们打算构建的功能，我们正在积极寻求反馈。如果您有任何想法，建议或问题，请告诉我们.
