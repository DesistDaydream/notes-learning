---
title: Go 标准库
---

# 概述

> 参考：
> - [GitHub,DesistDaydream/go-learning](https://github.com/DesistDaydream/go-learning)(学习代码)
> - [Go 包，标准库](https://pkg.go.dev/std)
> - [中文文档](https://studygolang.com/pkgdoc)
> - [go.dev,Tour(Go 语言之旅，通过在线解析器体验 Go 语言的各种特性)](https://go.dev/tour/list)

**Go Standard Library(Go 标准库)** 是 Go 内置 **Package(包)** 的集合，每个 package 都可以实现一类功能。每个 package 里有他们对应的常量、变量、函数、方法等。每个库就是一类功能，比如 bufio 库，这里面就是关于实现读写功能的各种内容；而 fmt 库则是关于实现格式化输入输出等功能。在[这里](https://pkg.go.dev/std?tab=packages)可以看到 go 语言 原生支持的所有标准库。

与 标准库 相对应的就是 [第三方库](https://www.teambition.com/project/5f90e312755d8a00446050eb/app/5eba5fba6a92214d420a3219/workspaces/5f90e312c800160016ea22fb/folders/5f9d3b0f4cc5830001c21a7c) ，第三方库一般属于由个人开发，实现更多丰富功能的库。在 [Go.dev ](https://pkg.go.dev/)可以搜索自己想要使用的所有库。

[Go.dev](https://pkg.go.dev/) 是 golang.org 的配套网站。 Golang.org 是开源项目和发行版的所在地，而 go.dev 是 Go 用户的中心，可从整个 Go 生态系统中提供集中和精选的资源。还可以在专门的[标准库](https://pkg.go.dev/std)页面看到所有标准库下的 Go 内置包。

Go.dev 提供：

- 在 index.golang.org 上发布的 Go 软件包和模块的集中信息。
- 基本学习资源
- 关键用例和案例研究

Go.dev 当前处于 MVP 状态。我们为自己的建设感到自豪，并很高兴与社区分享。我们希望您能在使用 go.dev 的过程中找到价值和乐趣。 Go.dev 只有一小部分我们打算构建的功能，我们正在积极寻求反馈。如果您有任何想法，建议或问题，请告诉我们.
