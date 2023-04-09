---
title: "Go 规范与标准库"
weight: 1
---

# 概述

> 参考：
> - [Go 官方文档，参考-语言规范](https://go.dev/ref/spec)
> - [Go 包，标准库](https://pkg.go.dev/std)
>     - [中文文档](https://studygolang.com/pkgdoc)
> - [go.dev,Tour(Go 语言之旅，通过在线解析器体验 Go 语言的各种特性)](https://go.dev/tour/list)
> - [GitHub,DesistDaydream/go-learning](https://github.com/DesistDaydream/go-learning)(学习代码)

Go 是一种通用语言，专为系统编程而设计。它是一种强类型且自带垃圾回收功能的语言，并具有显式支持并发编程的能力(称为 goroutine)。Go 程序由 Packages(包) 构建，其属性允许有效得管理依赖关系。

- **Go 语言参考**描述了 Go 语言的具体语法和语义
- **Go 标准库则**是与 Go 语言一起发行的一些可选功能，以便人们可以从一开始就轻松得使用 Go 进行编程。

# Go 语言关键字

> 参考：
> - [官方文档，参考-规范-关键字](https://go.dev/ref/spec#Keywords)

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

# Go 语言规范

> 参考：
> 
> - [官方文档，参考-规范](https://go.dev/ref/spec)
>     - [官方文档，参考-规范 的翻译](https://github.com/bekcpear/mypelicanconfandarticles/blob/master/content/Tech/gospec.rst)

## Notation(表示法)

Go 语言的语法遵从 [EBNF](https://en.wikipedia.org/wiki/Extended_Backus%E2%80%93Naur_form) 表示法

# Go 标准库

> 参考：
> 
> - [Go 包，标准库](https://pkg.go.dev/std)

**Go Standard Library(Go 标准库)** 是 Go 内置 **Package(包)** 的集合，每个 package 都可以实现一类功能。每个 package 里有他们对应的常量、变量、函数、方法等。每个库就是一类功能，比如 bufio 库，这里面就是关于实现读写功能的各种内容；而 fmt 库则是关于实现格式化输入输出等功能。在[这里](https://pkg.go.dev/std?tab=packages)可以看到 go 语言 原生支持的所有标准库。

与 标准库 相对应的就是 [第三方库](https://www.teambition.com/project/5f90e312755d8a00446050eb/app/5eba5fba6a92214d420a3219/workspaces/5f90e312c800160016ea22fb/folders/5f9d3b0f4cc5830001c21a7c) ，第三方库一般属于由个人开发，实现更多丰富功能的库。在 [Go.dev ](https://pkg.go.dev/)可以搜索自己想要使用的所有库。

[Go.dev](https://pkg.go.dev/) 是 golang.org 的配套网站。 Golang.org 是开源项目和发行版的所在地，而 go.dev 是 Go 用户的中心，可从整个 Go 生态系统中提供集中和精选的资源。还可以在专门的[标准库](https://pkg.go.dev/std)页面看到所有标准库下的 Go 内置包。

Go.dev 提供：

- 在 index.golang.org 上发布的 Go 软件包和模块的集中信息。
- 基本学习资源
- 关键用例和案例研究

Go.dev 当前处于 MVP 状态。我们为自己的建设感到自豪，并很高兴与社区分享。我们希望您能在使用 go.dev 的过程中找到价值和乐趣。 Go.dev 只有一小部分我们打算构建的功能，我们正在积极寻求反馈。如果您有任何想法，建议或问题，请告诉我们.