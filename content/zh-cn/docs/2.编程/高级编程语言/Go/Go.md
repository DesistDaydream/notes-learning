---
title: Go
linkTitle: Go
weight: 1
---

# 概述

> 参考：
>
> - [GitHub 项目，golang/go](https://github.com/golang/go)
> - [官网](https://golang.org/)
> - [Google 开放源代码](https://cs.opensource.google/go)
> - [GitHub 项目，avelino/awesome-go](https://github.com/avelino/awesome-go)(收录了优秀的 Go 框架、库、软件)
>   - [中文版，GitHub 项目，yinggaozhen/awesome-go-cn](https://github.com/yinggaozhen/awesome-go-cn)
>   - 另一个 go awesome: <https://github.com/shockerli/go-awesome>

Go 是一种开源编程语言，可以轻松构建 **simple(简单)**、**reliable(可靠)** 和 **efficient(高效)** 的软件。

# 学习资料

[菜鸟教程，Go](https://www.runoob.com/go/go-tutorial.html)（快速上手尝试，简单直接）

[Go 语言之旅](https://go.dev/tour)(官方在线教程)

[Go 官方 FAQ](https://go.dev/doc/faq)

[Go by Example](https://gobyexample.com/?tdsourcetag=s_pctim_aiomsg)

- [中文 Go by Example](https://gobyexample-cn.github.io/)

[GitHub 组织，golang-china](https://github.com/golang-china)(Go 语言中国)

[公众号-HelloGitHub，适合 Go 新手学习的开源项目](https://mp.weixin.qq.com/s/pAkjxK6N4shTEtHXQbxChg)

[地鼠文档](https://www.topgoer.cn/)系列文章

- [Go 编程模式](https://coolshell.cn/articles/series/go%e7%bc%96%e7%a8%8b%e6%a8%a1%e5%bc%8f)

电子书

- [GitHub 项目，unknowon/the-way-to-go](https://github.com/unknwon/the-way-to-go_ZH_CN)(Go 入门指南)
- [GitHub 项目，gopl-zh/gopl-zh.github.com](https://github.com/gopl-zh/gopl-zh.github.com)(Go 语言圣经)
  - [The Go Programming Language](https://www.k8stech.net/gopl/chapter0/)

视频

- [B 站-幼麟实验室-Golang 合辑](https://www.bilibili.com/video/BV1hv411x7we?spm_id_from=333.999.0.0&vd_source=708696360de7266de8f3911eef0f7448)

https://github.com/avelino/awesome-go

其他

- [B 站，80%哲学的践行者 ——“够用就行”的Go语言](https://www.bilibili.com/video/BV11FY9zZEwS)

# Hello World

代码：`hello_world.go`

```go
package main

import "fmt"

func main() {
 fmt.Println("Hello World")
}

```

运行

```bash
# go run hello_world.go
Hello World
```

# Go 范儿

> 参考：
>
> - [GitHub 项目 Wiki, golang/go-Wiki-Go 代码审查](https://github.com/golang/go/wiki/CodeReviewComments)
> - [官方文档，有效的 Go-名称](https://go.dev/doc/effective_go#names)
>   - [MakeOptim 博客，Effective Go 中文](https://makeoptim.com/golang/effective-go)(官方文档的中文翻译)
> - [Go 博客，Package names](https://go.dev/blog/package-names)
> - [博客园，不一样的 go 语言-gopher](https://www.cnblogs.com/laud/p/gopher.html)

gopher 原意地鼠，在 go 语言的世界里解释为地道的 go 程序员。在其他语言的世界里也有 PHPer，Pythonic 的说法，反而 Java 是个例外。虽然也有 Javaer 之类的说法，但似乎并不被认可。而地道或者说道地，说的是 gopher 写的代码无不透露出 go 的独特气息，比如项目结构、命名方式、代码格式、编码风格、构建方式等等。用 gopher 的话说，用 go 编写代码就像是在画一幅中国山水画，成品美不胜收，心旷神怡。

## 环境变量

gopher 第一条：把东西放对地方。

go 程序的运行，需要依赖于两个基础的环境变量，GOROOT 与 GOPATH。环境变量几乎在各类编程语言中都存在，比如 java 的 JAVA_HOME，其实也就是编译器及相关工具或标准库所在目录。

但 go 除了 GOROOT 之外，还增加了 GOPATH，它指的是 go 程序依赖的第三方库或自有库所在目录，以指示编译器从这些地方找到依赖。GOPATH 支持多个目录，通常一个目录就是一个项目，并且 GOPATH 目录按约定由 src, pkg, bin 三个目录组成。gopher 们的做法是定义 Global GOPATH, Project GOPATH，而更大的项目还会定义 Module GOPATH。当使用 go get 下载依赖时，会选择 GOPATH 环境变量中的第一个目录存放依赖包。

| 变量   | 含义              | 说明                                                 |
| ------ | ----------------- | ---------------------------------------------------- |
| GOROOT | go 运行环境根目录 | 通常指 go sdk 安装目录，包含编译器、官方工具及标准库 |
| GOPATH | 工作环境目录列表  | 通常指第三方库                                       |

## 项目结构

> 参考：
>
> - [GitHub 项目，golang-standards/project-layout](https://github.com/golang-standards/project-layout)
>   - [MakeOptim 博客，golang 编程规范-项目目录结构](https://makeoptim.com/golang/standards/project-layout)
> - [知乎，该如何组织 Go 项目结构？](https://zhuanlan.zhihu.com/p/346573562)
>   - [Package Oriented Design](https://www.ardanlabs.com/blog/2017/02/package-oriented-design.html)

gopher 第二条：按东西放在约定的地方。

不论采用何种编程语言，良好的项目组织结构都至关重要，因为这将直接影响项目内部依赖的复杂程度以及项目对外提供 API 等服务的灵活性等。最好在项目初期便制定好项目结构约定，甚至可以为其开发脚手架之类的工具来生成项目模板，让开发者尽量按照统一的规范参与项目。

一个常见的 Go 应用项目布局，通常有如下结构：

```latex
- my-go-project
 - cmd
 - pkg
 - internal
 - go.mod && go.sum
 - Makefile
```

其中 cmd 与 pkg 目录是最常见的。一个项目如果具有多个功能，比如 [kubernetes](https://github.com/kubernetes/kubernetes) 项目，具有多个组件，所有组件的入口都在 cmd 目录中，并以组件名命名其下的目录名。而每个组件所调用的各种功能，通常都是放在 pkg 目录下，一个功能一个目录，通常来说，pkg 目录是一个项目中代码量最多的地方。

### cmd

cmd 包是项目的主干，是编译构建的入口，`main()` 所在文件通常放置在此处。一个典型的 cmd 包的目录结构如下所示：

```latex
- cmd
   - app1
     - main.go
   - app2
     - main.go
```

从上述例子可以看出，cmd 下可以允许挂载多个需要编译的应用，只需要在不同的包下编写 main 文件即可。需要注意的是，cmd 中的代码应该尽量「保持简洁」，`main()` 函数中可能仅仅是参数初始化、配置加载、服务启动的操作。

### pkg

pkg 中存放的是可供项目内部/外部所使用的公共性代码，例如：用来连接第三方服务的 client 代码等。也有部分项目将该包命名为 lib，例如：[consul 项目](https://link.zhihu.com/?target=https%3A//github.com/hashicorp/consul) ，所表示的含义其实相同。

### internal

internal 包主要用处在于提供一个项目级别的代码保护方式，存放在其中的代码仅供项目内部使用。具体使用的规则是：.../a/b/c/internal/d/e/f 仅仅可以被.../a/b/c 下的目录导入，.../a/b/g 则不允许。internal 是 Go 1.4 版本中引入的特性，更多信息可以参考[这里](https://link.zhihu.com/?target=https%3A//golang.org/doc/go1.4%23internalpackages)。

在 internal 内部可以继续通过命名对目录的共享范围做区分，例如 internal/myapp 表示该目录下的代码是供 myapp 应用使用的；internal/pkg 表示该目录下的代码是可以供项目内多个应用使用的。

### go.mod && go.sum

go.mod 与 go.sum 是采用 [Go Module](/docs/2.编程/高级编程语言/Go/Go%20环境安装与使用/Go%20Module.md) 进行依赖管理所生成的配置文件。Go Modules 是 Go 1.11 版本中引入的版本管理功能，目前已经是 go 依赖管理的主流方式，所以此处不再讨论 vendor，dep 等依赖管理方式所生成的目录。

### Makefile

Makefile 文件通常存放项目的编译部署脚本。Go 的编译命令虽然简单，但总是手写命令还是效率低下，因此使用 Makefile 写编译部署脚本是工程实践中常见的方式。

### 服务端应用程序目录

#### api

项目对外提供和依赖的 API 文件。比如：OpenAPI/Swagger specs, JSON schema 文件, protocol 定义文件等。
比如，[Kubernetes](https://github.com/kubernetes/kubernetes/tree/master/api) 项目的 api 目录结构如下：

```bash
api
├── api-rules
    └── xxx.plist
├── openapi-spec
    └── swagger.json
```

因此，在 go 中用的比较多的 gRPC proto 文件，也比较适合放在 api 目录下。

```bash
api
└── protobuf-spec
    └── test
        ├── test.pb.go
        └── test.proto
```

### Web 应用程序目录

#### web

Web 应用程序特定的组件，比如 静态资源、服务器端模板和单页应用

### 通用应用程序目录

#### build

打包和持续集成所需的文件。

- build/ci：存放持续集成的配置和脚本，如果持续集成平台对配置文件有路径要求，则可将其 link 到指定位置。
- build/package：存放 AMI、Docker、系统包（deb、rpm、pkg）的配置和脚本等。

例子：

- <https://github.com/cockroachdb/cockroach/tree/master/build>

#### configs

配置文件模板或默认配置。

#### deployments

IaaS，PaaS，系统和容器编排部署配置和模板（docker-compose，kubernetes/helm，mesos，terraform，bosh）。请注意，在某些存储库中（尤其是使用 kubernetes 部署的应用程序），该目录的名字是 /deploy。

#### init

系统初始化（systemd、upstart、sysv）和进程管理（runit、supervisord）配置。

#### scripts

用于执行各种构建，安装，分析等操作的脚本。

这些脚本使根级别的 Makefile 变得更小更简单，例如：<https://github.com/hashicorp/terraform/blob/master/Makefile>。

#### test

外部测试应用程序和测试数据。随时根据需要构建 /test 目录。对于较大的项目，有一个数据子目录更好一些。例如，如果需要 Go 忽略目录中的内容，则可以使用 /test/data 或 /test/testdata 这样的目录名字。请注意，Go 还将忽略以“.”或“\_”开头的目录或文件，因此可以更具灵活性的来命名测试数据目录。

### 其他目录

#### assets

项目中使用的其他资源（图像、logo 等）。

#### docs

设计和用户文档（除了 godoc 生成的文档）。

#### examples

应用程序或公共库的示例程序。

#### githooks

Git 钩子。

#### third_party

外部辅助工具，fork 的代码和其他第三方工具（例如：Swagger UI）。

#### tools

此项目的支持工具。请注意，这些工具可以从 /pkg 和 /internal 目录导入代码。
例子：

- <https://github.com/istio/istio/tree/master/tools>
- <https://github.com/openshift/origin/tree/master/tools>
- <https://github.com/dapr/dapr/tree/master/tools>

#### website

如果不使用 Github pages，则在这里放置项目的网站数据。
例子：

- <https://github.com/hashicorp/vault/tree/master/website>
- <https://github.com/perkeep/perkeep/tree/master/website>

### 不应该包含的目录

项目中不应该包含 src 目录

在 Java 项目中，会常见 src 目录，但在 Go 项目中，并不推荐这么做。在 Go 1.11 之前，Go 项目是放在 `$GOPATH/src` 下，如果项目中再包含 src 目录，那么代码结构就会类似： `$GOPATH/src/my-project/src/app.go`的结构，容易造成混淆。在 Go 引入 modules 之后，项目可以不用写在 `$GOPATH` 下，但是依然不推荐项目中采用`src` 来命名目录。

## 命名规范

gopher 第三条：把名字起得 go 一点。

go 语言的命名与其他语言最大的不同在于首字母的大小写。

- 大写代表公开（导出，可以在其他包内访问）
- 小写代表私有（不导出，只能在包内访问）。

除此之外，与其他语言并无二致，比如不能以数字开头。而由于关键字、保留字的减少，因而减少了一些命名上的忌讳。更为突出的是，go 语言有一些建议性的命名规范，这也是 gophers 的圣经，理应严格遵守。

| 约定       | 范围                                                               | 说明                                                                                             | 示例                                    |
| ---------- | ------------------------------------------------------------------ | ------------------------------------------------------------------------------------------------ | --------------------------------------- |
| 驼峰命名法 | 全局                                                               | 统一使用驼峰命名法                                                                               | var isLive = false                      |
| 大小写一致 | 缩写短语，惯用词                                                   | 如 HTML，CSS, HTTP 等                                                                            | htmlEscape，HTMLEscape                  |
| 简短命名法 | 局部变量                                                           | 方法内、循环等使用的局部变量可以使用简短命名                                                     | 比如 for 循环中的 i，buf 代表 buffer 等 |
| 参数命名法 | 函数参数、返回值、方法接收者                                       | 如果参数类型能说明含义，则参数名可以简短，否则应该采用有文档说明能力的命名                       | 比如 d Duration，t Time                 |
| 通用命名法 | 作用域越大或者使用的地方离声明的地方太远，则应采用清晰有意义的命名 | -                                                                                                |                                         |
| 导出命名法 | 导出变量、函数、结构等                                             | 包名与导出名意义不要重复，同时包的命名要与导出的内容相关，不要使用宽泛的名字，如 common，util    | bytes.Buffer 比 bytes.ByteBuffer 要好   |
| 文件命名   | go 文件，单元测试文件                                              | go 文件名尽量以一个单词来命名，多个单词使用下线划分隔，单元测试文件以对应 go 文件名加\_test 结尾 | proto_test                              |
| 包命名     | 包                                                                 | 包的一级名称应是顶级域名，二级名称则应是项目名称，项目名称单词间以-分隔                          | github.com/mysql                        |

## 代码格式

gopher 第四条：按统一的格式来。

在多人协作团队中，统一的代码格式化模板是第一要义。在 Java 语言中，检验新人经验的一大法宝就是他有没有主动索要代码模板。而在 go 语言中，则没有这个必要了。因为 go 已经有默认的代码格式化工具了，而且代码格式化在 go 语言中是强制规范。所以这使得所有 go 程序员写出来的代码格式都是一样的。

go 默认的代码格式化工具是 gofmt。另外还有一个增强工具 goimport，在 gofmt 的基础上增加了自动删除和引入依赖包。而行长则以不超过 80 个字符为佳，超过请主动以多行展示。

## 编码风格

gopher 第五条：请学会约定

| 项         | 约定                                                                               | 说明                                                                                                                                                       |
| ---------- | ---------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------------- |
| import     | 按标准库、内部包、第三方包的顺序导入包                                             | 只引一个包时使用单行模式，否则使用多行模式                                                                                                                 |
| 变量声明   | 如果连续声明多个变量，应放在一起                                                   | 参见例子                                                                                                                                                   |
| 错误处理   | 不要忽略每一个 error，即使只是打一行日志                                           | go 的 error 处理方式与 C 同出一辙，通过返回值来标明错误或异常，引来的争议也很多，甚至官方已经开始酝酿在 go2 解决这个问题                                   |
| 长语句打印 | 使用格式化方式打印                                                                 | -                                                                                                                                                          |
| 注释规范   | 变量、方法、结构等的注释直接加上声明前，并且不要加空行。废弃方法加 Deprecated:即可 | 其中的第一行注释会被 godoc 识别为简短介绍，第二行开始则被认为是注释详情。注释对 godoc 的生成至关重要，因此关于注释会有一些技巧，我将在后面用专门的章节探讨 |

多变量声明

```go
var (
    name string
    age int
)
```

注释规范

```go
// Add 两数相加
// 两个整数相加，并返回和。
func Add(n1, n2 int)int{
    return n1 + n2
}
```

## 依赖管理

gopher 第六条：使用依赖管理工具管理自有依赖与第三方依赖

一个语言的生态是否完善甚至是否强大，除了 github 上面的开源项目数量之外，还有一大特征就是是否有优秀的依赖管理工具。依赖管理工具在业界已经是无处不在，yum、maven、gradle、pip、npm、cargo 这些工具的大名如雷贯耳。那么 go 有什么呢？

早期 go 的依赖是混乱的，因为没有一个工具能得到普遍认可，而官方又迟迟不出来解决问题。历数存在的工具包括 godep、glide、govender 等等。甚至早期还需要使用 GOPATH 来管理依赖，即项目的所有依赖都通过 go get 下载到指定的 GOPATH 中去。当然这种方案还可以撑大多数时间，但随着时间的流逝，随着开发人员的变动，这种管理依赖的弊端就慢慢显现出来。其实这些老路早期的 java 也走过，曾几何时，每个 java 项目里面都会有一个叫 lib 或 libs 的目录，这里放的就是当前项目依赖的包。当 GO 采用 GOPATH 来管理依赖时，开发人员只能被倒逼着用 java 的方式在源码库中自行管理依赖。这样相当于给依赖包做了隔离，同时又具备了版本管理（因为放在源码库）。

后来在 go1.5 的时候，官方引入了 vender 的概念，其实这也没改变多少，只是官方让大家存放依赖包的目录名称不要乱起了，统一叫 vender 吧。这个方案我觉得比依赖 GOPATH 还糟糕，因为 vendor 目录脱离了版本管理，导致更换依赖包版本很困难，在当前项目对依赖包的版本更新可能会影响其他项目的使用（如果新版本的依赖包有较大变动的话），同时如何将依赖包放到 vendor 下呢？等等。当然官方做出的这些变动可能是想像 maven 那样，推动社区来完成这件事，因而直接推动了上文提到的基于 vendor 的依赖管理工具的诞生。直至后来官方默认的社区做出来 dep，这下安静了，尽管刚开始时也不怎么好用，但有总比没有好。

go1.11 在 vgo 的基础上，官方推出了 go module。在发布前，官方与社区的大神们还为此开吵，认为官方太不厚道且独断专行。完全忽视 dep 社区的存在，无视 dep 在 go 语言中的地位与贡献。喜欢八卦的朋友们，可搜索《关于 Go Module 的争吵》一览大神是怎么吵架的，也可从中学习他们的思想。

相对于 java 的依赖管理工具 maven 或 gradle 来说，gradle 是 maven 的升级版，同时带来了 DSL 与元编程的特性，这无疑使得 gradle 异常地强大。但 gradle.io 在国内的可达情况也不尽如人意，好就好在其与 maven 仓库标准的兼容，F使得从 maven 转到 gradle 几乎没有额外的成本及阻力。

扯了这么多，依赖管理对于一门语言是必不可少的。c 有 cmake，java 有 maven、gradle，rust 有 cargo，那么 go 的 dep 或者 module 就用起来吧，看完大神吵架之后，喜欢哪个就选哪个。是不可能产生一个能满足所有人要求的依赖管理工具的，就连号称最牛逼的 cargo 也不例外。在一般的项目中，能用到的依赖管理功能也就那常用的几个而已，对大多数项目来说，适用好用就行。

## 构建方式

gopher 第七条：按需构建

构建的目标是让代码成为可运行程序。构建的过程应该是低成本并且让人愉悦的，显然 C 在这一方面让人抓狂，而 go 确实做得不错。并且能在任何平台下编译出另外一个平台的可执行程序。不管你的 go 程序是 CLI、GUI、WEB 或者其他形式的网络通讯程序，在 go 的世界里都只需要一个命令构建成可执行程序（依赖也一并被打包），即可在目标系统上运行。在这一点上，java 是望尘莫及了。
下面是用来构建 go 程序常用的参数，其他参数可通过 go help environment 命令查看。

| 参数        | 值                              | 说明                                                     |
| ----------- | ------------------------------- | -------------------------------------------------------- |
| CGO_ENABLED | 0 or 1                          | 是否支持 cgo 命令，如果 go 代码中有 c 代码，需要设置为 1 |
| GOOS        | darwin, freebsd, linux, windows | 可执行程序运行的目标操作系统                             |
| GOARCH      | 386, amd64, arm                 | 可执行程序运行的目标操作系统架构                         |

```bash
# Linux下编译Mac 64位可执行程序
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build main.go
# Linux下编译windows 64位可执行程序
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build main.go
# 如果想减少二进制程序的大小，可以加上-ldflags "-s -w"，但同时会丢掉调试信息，即不能用gdb调试了。
# 如果想更进一步减少程序大小，可以使用加壳工具，比如upx
```
