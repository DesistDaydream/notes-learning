---
title: pprof
weight: 1
---

# 概述

> 参考：
>
> - [Go 官方文档，诊断](https://go.dev/doc/diagnostics)
> - [GitHub 项目，google/pprof](https://github.com/google/pprof)
> - [GitHub 项目-文档，google/pprof/doc](https://github.com/google/pprof/tree/master/doc)
> - [Go 包，net/http/pprof](https://pkg.go.dev/net/http/pprof)
> - [Go 博客，分析 Go 程序](https://go.dev/blog/pprof)
> - [思否，Golang 大杀器之性能剖析 PProf](https://segmentfault.com/a/1190000016412013)
> - [公众号-云原生实验室，忙碌的开发人员的 Go Profiling、跟踪和可观察性指南](https://mp.weixin.qq.com/s/SveQPLr7abKXccLpYKkNKA)

pprof 是 **go 程序的性能分析器**，一个可视化和分析 Profiling 数据的工具。pprof 可以从目标获取运行数据并生成 profile.proto 格式的 Profiles 文件，还可以读取 profile.proto 格式的 Profiling 样本集合，并生成报告。

profile.proto 是一个协议缓冲区，描述了一组调用堆栈和符号化信息。详见 <https://developers.google.com/protocol-buffers>

可以通过本地文件或 HTTP 读取 Profiles 文件。同时也可以聚合或比较多个 Profiles 文件。每个 profile.prot 格式的 Profile 样本的集合。

## 使用 pprof

想要使用 pprof 程序非常简单，只需要引入 `net/http/pprof` 包，并启动监听即可

```go
package main

import (
 "log"
 "net/http"
 _ "net/http/pprof"
)

func main() {
 if err := http.ListenAndServe("localhost:18080", nil); err != nil {
  log.Fatal("ListenAndServe: ", err)
 }
}
```

注意：若不使用 `DefaultServeMux`，则我们需要手动注册处理程序，效果如下：

```go
package main

import (
    "log"
    "net/http"
    _ "net/http/pprof"
)

func main() {
    mux := http.NewServeMux()
    mux.Handle("/debug/pprof/", http.HandlerFunc(pprof.Index))
    mux.Handle("/debug/pprof/cmdline", http.HandlerFunc(pprof.Cmdline))
    mux.Handle("/debug/pprof/profile", http.HandlerFunc(pprof.Profile))
    mux.Handle("/debug/pprof/symbol", http.HandlerFunc(pprof.Symbol))
    mux.Handle("/debug/pprof/trace", http.HandlerFunc(pprof.Trace))
    if err := http.ListenAndServe("localhost:18080", nil); err != nil {
        log.Fatal("ListenAndServe: ", err)
    }
}
```

此时我们可以在 pprof 启动的 HTTP 服务端的 `/debug/pprof` 端点查看本程序 profile 信息：

![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/gz5a6m/1652319650333-cea6d5e1-0f15-4981-93f9-86d5833ecf57.png)

pprof 库会暴露多个端点

- **/debug/pprof/allocs** #
  - 用与分析 Mem 申请内存频率过高的问题。比如 Go 频繁垃圾回收等
- **/debug/pprof/block** # 查看导致阻塞同步的堆栈跟踪
- **/debug/pprof/cmdline** #
- **/debug/pprof/goroutine** # 查看当前所有运行的 goroutines 堆栈跟踪
  - 用于分析 Goroutine 泄露问题
- **/debug/pprof/heap** # 查看活动对象的内存分配情况
  - 用于分析 Mem 使用率高的问题
- **/debug/pprof/mutex** # 查看导致互斥锁的竞争持有者的堆栈跟踪
  - 用于分析锁的抢占问题
- **/debug/pprof/profile** # 默认获取程序 30s 对 CPU 的使用情况的 Profile 文件。
  - 用于分析 CPU 使用率高的问题

## 不通过 net/http 标准库使用 pprof

若不直接使用 `net/http` 标准库，则需要手动注册路由

### gin 库

```go
package main

import (
 "github.com/gin-contrib/pprof"
 "github.com/gin-gonic/gin"
)

func main() {
 router := gin.Default()
 pprof.Register(router)
 router.Run(":8080")
}
```

### chi 库

```go
package main

import (
 "net/http"

 "github.com/go-chi/chi"
 "github.com/go-chi/chi/middleware"
)

func main() {
 r := chi.NewRouter()
 r.Mount("/debug", middleware.Profiler())
 http.ListenAndServe(":8080", r)
}
```

## top 信息

- flat：给定函数上运行耗时
- flat%：同上的 CPU 运行耗时总比例
- sum%：给定函数累积使用 CPU 总比例
- cum：当前函数加上它之上的调用运行总耗时
- cum%：同上的 CPU 运行耗时总比例

# pprof 工具

pprof 工具运行时，会在 ${HOME}/pprof/ 目录下生成临时的 `*.pb.gz` 文件

## Syntax(语法)

**pprof \[FORMAT] \[OPTIONS] \[BINARY] \<SOURCE> ...**

从 SOURCE 处获取性能信息数据，并在当前目录下生成 FORMAT 格式的 Profile 文件，文件名默认为 `profileXXX.pb.gz`。若省略 FORMAT，则将会进入交互式 CLI。

在省略 FORMAT 时，提供 `-http` 参数，pprof 会启动 HTTP 服务，可以通过浏览器浏览 Profile 信息。

**FORMAT**

- **-proto** # 以压缩的 protobuf 格式输出 Profile 文件。`默认生成的缓存 Profile 文件就是这种格式`

**SOURCE OPTIONS**

- **-seconds=\<INT>** # 采集 SOURCE 的持续时间，单位：秒

## 交互式 CLI

**top -cum 5** # 按照资源使用率排序并查看前 5 个

# 最佳实践

当我们在程序中注册了 pprof 之后，就可以开始使用 pprof 工具对获取到的性能数据进行分析。

## 生成 Profile 文件

首先，我们需要先使用 pprof 工具从 `/debug/pprof/profile` **端点**获取性能数据。默认情况会采集 30 秒程序运行数据，并缓存 potol 格式的 Profile 文件到 `${HOME}/pprof/` 目录中，同时在当前目录生成指定格式的 Profile 文件

```bash
~]# go tool pprof -proto http://localhost:18080/debug/pprof/profile
Fetching profile over HTTP from http://localhost:18080/debug/pprof/profile
Saved profile in /home/desistdaydream/pprof/pprof.main.samples.cpu.001.pb.gz
Generating report in profile001.pb.gz
```

> 注意：
>
> - 若待采集 profile 文件的目标程序需要使用 HTTPS，则将 http 改为 https+insecure

## 分析 Profile 文件

想要分析 profile.proto 文件，同样需要使用 pprof 程序。我们可以通过两种方式分析数据文件

- 交互式 CLI
- Web

### 交互式 CLI

```bash
~]# go tool pprof profile001.pb.gz
File: main
Type: cpu
Time: May 12, 2022 at 9:07am (CST)
Duration: 10s, Total samples = 10ms (  0.1%)
Entering interactive mode (type "help" for commands, "o" for options)
(pprof)
```

CLI 中最常见的操作就是 top、list、web，web 命令是打开一个 Web 服务端，这里可以查看识图、火焰图等等

`top` 命令可以获取最消耗资源的函数名，然后通过 `list FuncName` 列出其中函数中最消耗资源的几行代码并且会展示出资源消耗的数值。

### Web

使用 `-http` 参数以便让 pprof 工具监听一个端口，通过 Web 来分析数据

```bash
~]# go tool pprof -http=:8080 profile001.pb.gz
Serving web UI on http://localhost:8080
```

> 也可以在交互式 CLI 中使用 `web` 命令打开该界面。若出现 `Could not execute dot; may need to install graphviz.` 报错，安装 graphviz 即可。

在 Web 中可以将 profile 数据解析成火焰图、以图形和线条方式展示调用链、等等等

效果如下：

![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/gz5a6m/1652326690785-770c9f40-8461-4904-9958-aefc3a88ee9f.png)

## 最简单的排查 CPU 使用率问题

在 `pkg/webhook/webhook.go` 中添加 `h.router.Mount("/debug", middleware.Profiler())` 代码以引入 pprof

```go
h.router.Use(middleware.Timeout(2 * h.WebhookTimeout))
h.router.Mount("/debug", middleware.Profiler())
h.router.Get("/receivers", h.handler.ListReceivers)
```

程序启动后开始获取 CPU 信息的 Profile 文件

```bash
~]# go tool pprof -proto http://localhost:19093/debug/pprof/profile
Fetching profile over HTTP from http://localhost:19093/debug/pprof/profile
Saved profile in /home/desistdaydream/pprof/pprof.main.samples.cpu.002.pb.gz
Generating report in profile001.pb.gz
```

同时发起一些请求，触发其中的代码逻辑

以 Web 形式分析 Profile 文件

```bash
~]# go tool pprof -http=":8080" profile001.pb.gz
Serving web UI on http://localhost:8080
```

首先可以看到调用逻辑关系

![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/gz5a6m/1652341730532-08dc5241-a436-4288-bb7e-76e7e5f168ef.png)

这里可以查看每个函数调用所占用的 CPU 时间

![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/gz5a6m/1652341794711-3758bcfe-b3e6-4611-8c7b-0c02338a97fe.png)

这里可以查看火焰图

![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/gz5a6m/1652341761629-901f7396-c816-4d99-bedf-0a51ecda903c.png)

# pprof 实战

原文链接：[腾讯云+社区，golang pprof 实战-CPU,heap,alloc,goroutine,mutex,block](https://cloud.tencent.com/developer/article/1485112)

## 前言

如果要说在 golang 开发过程进行性能调优，pprof 一定是一个大杀器般的工具。但在网上找到的教程都偏向简略，难寻真的能应用于实战的教程。这也无可厚非，毕竟 pprof 是当程序占用资源异常时才需要启用的工具，而我相信大家的编码水平和排场问题的能力是足够高的，一般不会写出性能极度堪忧的程序，且即使发现有一些资源异常占用，也会通过排查代码快速定位，这也导致 pprof 需要上战场的机会少之又少。即使大家有心想学习使用 pprof，却也常常相忘于江湖。

**既然如此，那我就送大家一个性能极度堪忧的“炸弹”程序吧！**

这程序没啥正经用途却极度占用资源，基本覆盖了常见的性能问题。本文就是希望读者能一步一步按照提示，使用 pprof 定位这个程序的的性能瓶颈所在，借此学习 pprof 工具的使用方法。

因此，本文是一场“实验课”而非“理论课”，请读者腾出时间，脚踏实地，一步一步随实验步骤进行操作，这会是一个很有趣的冒险，不会很无聊，希望你能喜欢。

## 实验准备

这里假设你有基本的 golang 开发功底，拥有 golang 开发环境并配置了 $GOPATH，能熟练阅读简单的代码或进行简单的修改，且知道如何编译运行 golang 程序。此外，需要你大致知道 pprof 是干什么的，有一个基本印象即可，你可以花几分钟时间读一下[《Golang 大杀器之性能剖析 PProf》](https://blog.wolfogre.com/posts/go-ppof-practice/)的开头部分，这不会耽误太久。

此外由于你需要运行一个“炸弹”程序，请务必确保你用于做实验的机器有充足的资源，你的机器至少需要：

- 2 核 CPU；
- 2G 内存。

注意，以上只是最低需求，你的机器配置能高于上述要求自然最好。实际运行“炸弹”时，你可以关闭电脑上其他不必要的程序，甚至 IDE 都不用开，我们的实验操作基本上是在命令行里进行的。

此外，务必确保你是在个人机器上运行“炸弹”的，能接受机器死机重启的后果（虽然这发生的概率很低）。请你务必不要在危险的边缘试探，比如在线上[服务器](https://cloud.tencent.com/product/cvm?from=10680)运行这个程序。

可能说得你都有点害怕了，为打消你顾虑，我可以剧透一下“炸弹”的情况，让你安心：

- 程序会占用约 2G 内存；
- 程序占用 CPU 最高约 100%（总量按“核数 \* 100%”来算）；
- 程序不涉及网络或文件读写；
- 程序除了吃资源之外没有其他危险操作。

且程序所占用的各类资源，均不会随着运行时间的增长而增长，换句话说，只要你把“炸弹”启动并正常运行了一分钟，就基本确认安全了，之后即使运行几天也不会有更多的资源占用，除了有点费电之外。

## 获取“炸弹”

炸弹程序的代码我已经放到了 [GitHub](https://links.jianshu.com/go?to=https%3A%2F%2Fblog.wolfogre.com%2Fredirect%2Fv3%2FA_4-v86v-9Btg9a9FuRKCcgSAwM8Cv46xcU7LxImWv3FQQYW3DshxTsGzDw8cyzMPIIcSogxEgMDPAr-OsXFWhYGO25BBhbcOyH9xTwGTQrFOwbMPDwFzDyCHEqIxQ) 上，你只需要在终端里运行 `go get` 便可获取，注意加上 `-d` 参数，避免下载后自动安装：

```bash
go get -d github.com/wolfogre/go-pprof-practice
cd $GOPATH/src/github.com/wolfogre/go-pprof-practice
```

我们可以简单看一下 `main.go` 文件，里面有几个帮助排除性能调问题的关键的的点，我加上了些注释方便你理解，如下：

```go
package main

import (
    _ "net/http/pprof"
)

func main() {
    runtime.GOMAXPROCS(1)
    runtime.SetMutexProfileFraction(1)
    runtime.SetBlockProfileRate(1)
    go func() {
        if err := http.ListenAndServe(":6060", nil); err != nil {
            log.Fatal(err)
        }
        os.Exit(0)
    }()
}
```

除此之外的其他代码你一律不用看，那些都是我为了模拟一个“逻辑复杂”的程序而编造的，其中大多数的问题很容易通过肉眼发现，但我们需要做的是通过 pprof 来定位代码的问题，所以为了保证实验的趣味性请不要提前阅读代码，可以实验完成后再看。

接着我们需要编译一下这个程序并运行，你不用担心依赖问题，这个程序没有任何外部依赖。

```bash
go build
./go-pprof-practice
```

运行后注意查看一下资源是否吃紧，机器是否还能扛得住，坚持一分钟，如果确认没问题，咱们再进行下一步。

控制台里应该会不停的打印日志，都是一些“猫狗虎鼠在不停地吃喝拉撒”的屁话，没有意义，不用细看。

![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/gz5a6m/1652342670773-707536b2-d1e0-4dcb-a837-97c73ce2ec20.png)

## 使用 pprof

保持程序运行，打开浏览器访问 `http://localhost:6060/debug/pprof/`，可以看到如下页面：

![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/gz5a6m/1652342670421-8ffb9b61-ff9a-4b0b-bf6e-57df41394c08.png)

页面上展示了可用的程序运行采样数据，分别有：

| 类型         | 描述                       | 备注                                                              |
| ------------ | -------------------------- | ----------------------------------------------------------------- |
| allocs       | 内存分配情况的采样信息     | 可以用浏览器打开，但可读性不高                                    |
| blocks       | 阻塞操作情况的采样信息     | 可以用浏览器打开，但可读性不高                                    |
| cmdline      | 显示程序启动命令及参数     | 可以用浏览器打开，这里会显示 ./go-pprof-practice                  |
| goroutine    | 当前所有协程的堆栈信息     | 可以用浏览器打开，但可读性不高                                    |
| heap         | 堆上内存使用情况的采样信息 | 可以用浏览器打开，但可读性不高                                    |
| mutex        | 锁争用情况的采样信息       | 可以用浏览器打开，但可读性不高                                    |
| profile      | CPU 占用情况的采样信息     | 浏览器打开会下载文件                                              |
| threadcreate | 系统线程创建情况的采样信息 | 可以用浏览器打开，但可读性不高                                    |
| trace        | 程序运行跟踪信息           | 浏览器打开会下载文件，本文不涉及，可另行参阅《深入浅出 Go trace》 |

因为 cmdline 没有什么实验价值，trace 与本文主题关系不大，threadcreate 涉及的情况偏复杂，所以这三个类型的采样信息这里暂且不提。除此之外，其他所有类型的采样信息本文都会涉及到，且炸弹程序已经为每一种类型的采样信息埋藏了一个对应的性能问题，等待你的发现。

由于直接阅读采样信息缺乏直观性，我们需要借助 `go tool pprof` 命令来排查问题，这个命令是 go 原生自带的，所以不用额外安装。

我们先不用完整地学习如何使用这个命令，毕竟那太枯燥了，我们一边实战一边学习。

以下正式开始。

## 排查 CPU 占用过高

我们首先通过活动监视器（或任务管理器、top 命令，取决于你的操作系统和你的喜好），查看一下炸弹程序的 CPU 占用：

![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/gz5a6m/1652342670264-907c836e-e6a7-4962-8809-1be829ae7078.png)

可以看到 CPU 占用相当高，这显然是有问题的，我们使用 `go tool pprof` 来排场一下：

```bash
go tool pprof http://localhost:6060/debug/pprof/profile
```

等待一会儿后，进入一个交互式终端：

![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/gz5a6m/1652342670211-12943d99-488c-41f9-bcc0-ebb23e9125b9.png)

输入 top 命令，查看 CPU 占用较高的调用：

![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/gz5a6m/1652342670434-f5e5887e-67e4-40cc-8235-30f14e1a2e31.png)

很明显，CPU 占用过高是 `github.com/wolfogre/go-pprof-practice/animal/felidae/tiger.(*Tiger).Eat` 造成的。

> 注：为了保证实验节奏，关于图中 `flat`、`flat%`、`sum%`、`cum`、`cum%` 等参数的含义这里就不展开讲了，你可以先简单理解为数字越大占用情况越严重，之后可以在[《Golang 大杀器之性能剖析 PProf》](https://links.jianshu.com/go?to=https%3A%2F%2Fblog.wolfogre.com%2Fredirect%2Fv3%2FA3jsjsv0r3pCsn4x_qmqFKwSAwM8Cv46xcU7LxImWv3F_wdFRERZQ0pZxVoWBjvFWhYGWsWtTRvFOwaJVMX_BDIwMTjM_wIwOcz_AjE1zP5HBolU_1AlMjAlRTUlQTQlQTclRTYlOUQlODAlRTUlOTklQTglRTQlQjklOEIlRTYlODAlQTclRTglODMlQkQlRTUlODklOTYlRTYlOUUlOTAlMjBQUHMsbi0YMRIDAzwK_jrFxVoWBjtuQQYW3Dsh_cU8Bk0KxTsGzDw8Bcw8ghxKiMU)等其他资料中深入学习。

输入 `list Eat`，查看问题具体在代码的哪一个位置：

![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/gz5a6m/1652342672266-67b39808-d532-47f5-9f51-6a3597bc0620.png)

可以看到，是第 24 行那个一百亿次空循环占用了大量 CPU 时间，至此，问题定位成功！

接下来有一个扩展操作：图形化显示调用栈信息，这很酷，但是需要你事先在机器上安装 `graphviz`，大多数系统上可以轻松安装它：

```bash
brew install graphviz # for macos
apt install graphviz # for ubuntu
yum install graphviz # for centos
```

或者你也可以访问 [graphviz 官网](https://links.jianshu.com/go?to=https%3A%2F%2Fblog.wolfogre.com%2Fredirect%2Fv3%2FA421Yoc_xEV4GG_UO8tV1nMSAwM8Cv46xcU7gjwSbQjbbjsviVpukMUYBkEJFgboxTESAwM8Cv46xcVaFgY7bkEGFtw7If3FPAZNCsU7Bsw8PAXMPIIcSojF)寻找适合自己操作系统的安装方法。

安装完成后，我们继续在上文的交互式终端里输入 `web`，注意，虽然这个命令的名字叫“web”，但它的实际行为是产生一个 .svg 文件，并调用你的系统里设置的默认打开 .svg 的程序打开它。如果你的系统里打开 .svg 的默认程序并不是浏览器（比如可能是你的代码编辑器），这时候你需要设置一下默认使用浏览器打开 .svg 文件，相信这难不倒你。

你应该可以看到：

![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/gz5a6m/1652342672880-6d25d002-b6d6-4517-9462-7ba47d2fe941.png)

图中，`tiger.(*Tiger).Eat` 函数的框特别大，箭头特别粗，pprof 生怕你不知道这个函数的 CPU 占用很高，这张图还包含了很多有趣且有价值的信息，你可以多看一会儿再继续。

至此，这一小节使用 pprof 定位 CPU 占用的实验就结束了，你需要输入 `exit` 退出 pprof 的交互式终端。

为了方便进行后面的实验，我们修复一下这个问题，不用太麻烦，注释掉相关代码即可：

```go
func (t *Tiger) Eat() {
    log.Println(t.Name(), "eat")
}
```

之后修复问题的的方法都是注释掉相关的代码，不再赘述。你可能觉得这很粗暴，但要知道，这个实验的重点是如何使用 pprof 定位问题，我们不需要花太多时间在改代码上。

## 排查内存占用过高

重新编译炸弹程序，再次运行，可以看到 CPU 占用率已经下来了，但是内存的占用率仍然很高：

![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/gz5a6m/1652342673383-82d8e798-d1b6-42fc-ab58-7f70d1c2d0d4.png)

我们再次运行使用 pprof 命令，注意这次使用的 URL 的结尾是 heap：

```bash
go tool pprof http://localhost:6060/debug/pprof/heap
```

再一次使用 `top`、`list` 来定问问题代码：

![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/gz5a6m/1652342675687-d36c06ac-c6b8-4739-9fd0-ae5272b71cb6.png)

可以看到这次出问题的地方在 `github.com/wolfogre/go-pprof-practice/animal/muridae/mouse.(*Mouse).Steal`，函数内容如下：

```go
func (m *Mouse) Steal() {
    log.Println(m.Name(), "steal")
    max := constant.Gi
    for len(m.buffer) * constant.Mi < max {
        m.buffer = append(m.buffer, [constant.Mi]byte{})
    }
}
```

可以看到，这里有个循环会一直向 m.buffer 里追加长度为 1 MiB 的数组，直到总容量到达 1 GiB 为止，且一直不释放这些内存，这就难怪会有这么高的内存占用了。

使用 `web` 来查看图形化展示，可以再次确认问题确实出在这里：

![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/gz5a6m/1652342675266-a848fa57-5cb9-4188-bedd-0c84ed45c6db.png)

现在我们同样是注释掉相关代码来解决这个问题。

再次编译运行，查看内存占用：

![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/gz5a6m/1652342675895-2104ed19-316d-4df4-816d-a3fecb1e99d1.png)

可以看到内存占用已经将到了 35 MB，似乎内存的使用已经恢复正常，一片祥和。

但是，内存相关的性能问题真的已经全部解决了吗？

## 排查频繁内存回收

你应该知道，频繁的 GC 对 golang 程序性能的影响也是非常严重的。虽然现在这个炸弹程序内存使用量并不高，但这会不会是频繁 GC 之后的假象呢？

为了获取程序运行过程中 GC 日志，我们需要先退出炸弹程序，再在重新启动前赋予一个环境变量，同时为了避免其他日志的干扰，使用 grep 筛选出 GC 日志查看：

```bash
GODEBUG=gctrace=1 ./go-pprof-practice | grep gc
```

日志输出如下：

![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/gz5a6m/1652342676751-565d6c37-649d-4a45-abbc-22de6a6764d6.png)

可以看到，GC 差不多每 3 秒就发生一次，且每次 GC 都会从 16MB 清理到几乎 0MB，说明程序在不断的申请内存再释放，这是高性能 golang 程序所不允许的。

如果你希望进一步了解 golang 的 GC 日志可以查看[《如何监控 golang 程序的垃圾回收》](https://links.jianshu.com/go?to=https%3A%2F%2Fblog.wolfogre.com%2Fredirect%2Fv3%2FA9DNc05mRFLA-ZPsjfPhLuZDu-oKbuLF_wQyMDE2xf8CMDfF_wIwMcUtHy8qzDsGiVTMOxzFMRIDAzwK_jrFxVoWBjtuQQYW3Dsh_cU8Bk0KxTsGzDw8Bcw8ghxKiMU),为保证实验节奏，这里不做展开。

所以接下来使用 pprof 排查时，我们在乎的不是什么地方在占用大量内存，而是什么地方在不停地申请内存，这两者是有区别的。

由于内存的申请与释放频度是需要一段时间来统计的，所有我们保证炸弹程序已经运行了几分钟之后，再运行命令：

```bash
go tool pprof http://localhost:6060/debug/pprof/allocs
```

同样使用 top、list、web 大法：

![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/gz5a6m/1652342678255-9fbe8591-fd0c-426d-b10b-384a944a79fa.png)

![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/gz5a6m/1652342677694-f7eb46e4-bbf3-4f83-a2eb-6fb0574cba27.png)

可以看到 `github.com/wolfogre/go-pprof-practice/animal/canidae/dog.(*Dog).Run` 会进行无意义的内存申请，而这个函数又会被频繁调用，这才导致程序不停地进行 GC:

```go
func (d *Dog) Run() {
    log.Println(d.Name(), "run")
    _ = make([]byte, 16 * constant.Mi)
}
```

这里有个小插曲，你可尝试一下将 `16 * constant.Mi` 修改成一个较小的值，重新编译运行，会发现并不会引起频繁 GC，原因是在 golang 里，对象是使用堆内存还是栈内存，由编译器进行逃逸分析并决定，如果对象不会逃逸，便可在使用栈内存，但总有意外，就是对象的尺寸过大时，便不得不使用堆内存。所以这里设置申请 16 MiB 的内存就是为了避免编译器直接在栈上分配，如果那样得话就不会涉及到 GC 了。

我们同样注释掉问题代码，重新编译执行，可以看到这一次，程序的 GC 频度要低很多，以至于短时间内都看不到 GC 日志了：

![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/gz5a6m/1652342677869-a929d664-3bb1-4a81-b44b-9e085c70a69f.png)

## 排查协程泄露

由于 golang 自带内存回收，所以一般不会发生内存泄露。但凡事都有例外，在 golang 中，协程本身是可能泄露的，或者叫协程失控，进而导致内存泄露。

我们在浏览器里可以看到，此时程序的协程数已经多达 106 条：
![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/gz5a6m/1652342920133-d94f6be1-3553-4e5b-9969-d90ee3e8eaee.png)
虽然 106 条并不算多，但对于这样一个小程序来说，似乎还是不正常的。为求安心，我们再次是用 pprof 来排查一下：

```bash
go tool pprof http://localhost:6060/debug/pprof/goroutine
```

同样是 top、list、web 大法：
![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/gz5a6m/1652342931465-9bce9a37-20cb-4a3f-832a-ea349f816a14.png)

![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/gz5a6m/1652342954617-7ec317ec-6c0b-444e-9107-591e4b80fa1f.png)

可能这次问题藏得比较隐晦，但仔细观察还是不难发现，问题在于 `github.com/wolfogre/go-pprof-practice/animal/canidae/wolf.(*Wolf).Drink` 在不停地创建没有实际作用的协程：

```go
func (w *Wolf) Drink() {
    log.Println(w.Name(), "drink")
    for i := 0; i < 10; i++ {
        go func() {
            time.Sleep(30 * time.Second)
        }()
    }
}
```

可以看到，Drink 函数每次回释放 10 个协程出去，每个协程会睡眠 30 秒再退出，而 Drink 函数又会被反复调用，这才导致大量协程泄露，试想一下，如果释放出的协程会永久阻塞，那么泄露的协程数便会持续增加，内存的占用也会持续增加，那迟早是会被操作系统杀死的。

我们注释掉问题代码，重新编译运行可以看到，协程数已经降到 4 条了：

![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/gz5a6m/1652343010037-df499f46-85e0-4835-887e-e6c36439ac23.png)

## 排查锁的争用

到目前为止，我们已经解决这个炸弹程序的所有资源占用问题，但是事情还没有完，我们需要进一步排查那些会导致程序运行慢的性能问题，这些问题可能并不会导致资源占用，但会让程序效率低下，这同样是高性能程序所忌讳的。

我们首先想到的就是程序中是否有不合理的锁的争用，我们倒一倒，回头看看上一张图，虽然协程数已经降到 4 条，但还显示有一个 mutex 存在争用问题。

相信到这里，你已经触类旁通了，无需多言，开整。

```bash
go tool pprof http://localhost:6060/debug/pprof/mutex
```

同样是 top、list、web 大法：

![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/gz5a6m/1652343071279-b6754799-45a6-4769-9f73-7d474636fe36.png)

![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/gz5a6m/1652343075415-a3376489-a540-447c-b144-7fef8fe838b5.png)

可以看出来这问题出在 `github.com/wolfogre/go-pprof-practice/animal/canidae/wolf.(*Wolf).Howl`。但要知道，在代码中使用锁是无可非议的，并不是所有的锁都会被标记有问题，我们看看这个有问题的锁那儿触雷了。

```go
func (w *Wolf) Howl() {
    log.Println(w.Name(), "howl")
    m := &sync.Mutex{}
    m.Lock()
    go func() {
        time.Sleep(time.Second)
        m.Unlock()
    }()
    m.Lock()
}
```

可以看到，这个锁由主协程 Lock，并启动子协程去 Unlock，主协程会阻塞在第二次 Lock 这儿等待子协程完成任务，但由于子协程足足睡眠了一秒，导致主协程等待这个锁释放足足等了一秒钟。虽然这可能是实际的业务需要，逻辑上说得通，并不一定真的是性能瓶颈，但既然它出现在我写的“炸弹”里，就肯定不是什么“业务需要”啦。

所以，我们注释掉这段问题代码，重新编译执行，继续。

## 排查阻塞操作

好了，我们开始排查最后一个问题。

在程序中，除了锁的争用会导致阻塞之外，很多逻辑都会导致阻塞。

![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/gz5a6m/1652343092686-180c61a0-e682-4666-84ba-528f17fd346d.png)

可以看到，这里仍有 2 个阻塞操作，虽然不一定是有问题的，但我们保证程序性能，我们还是要老老实实排查确认一下才对。

```bash
go tool pprof http://localhost:6060/debug/pprof/block
```

top、list、web，你懂得。

![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/gz5a6m/1652343109220-a11fcfbf-b808-4da2-b253-3d4339319bd9.png)

![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/gz5a6m/1652343120730-030aaaa4-47fa-41be-873b-ee97e32c730f.png)

可以看到，阻塞操作位于 `github.com/wolfogre/go-pprof-practice/animal/felidae/cat.(*Cat).Pee`：

```go
func (c *Cat) Pee() {
    log.Println(c.Name(), "pee")
    <-time.After(time.Second)
}
```

你应该可以看懂，不同于睡眠一秒，这里是从一个 channel 里读数据时，发生了阻塞，直到这个 channel 在一秒后才有数据读出，这就导致程序阻塞了一秒而非睡眠了一秒。

这里有个疑点，就是上文中是可以看到有两个阻塞操作的，但这里只排查出了一个，我没有找到其准确原因，但怀疑另一个阻塞操作是程序监听端口提供 porof 查询时，涉及到 IO 操作发生了阻塞，即阻塞在对 HTTP 端口的监听上，但我没有进一步考证。

不管怎样，恭喜你完整地完成了这个实验。

## 思考题

另有一些问题，虽然比较重要，但碍于篇幅（其实是我偷懒），就以思考题的形式留给大家了。

1. 每次进入交互式终端，都会提示“type ‘help’ for commands, ‘o’ for options”，你有尝试过查看有哪些命令和哪些选项吗？有重点了解一下“sample_index”这个选项吗？
2. 关于内存的指标，有申请对象数、使用对象数、申请空间大小、使用空间大小，本文用的是什么指标？如何查看不同的指标？（提示：在内存实验中，试试查看、修改“sample_index”选项。）
3. 你有听说过火焰图吗？要不要在试验中生成一下火焰图？
4. main 函数中的 `runtime.SetMutexProfileFraction` 和 `runtime.SetBlockProfileRate` 是如何影响指标采样的？它们的参数的含义是什么？

## 最后

碍于我的水平有限，实验中还有很多奇怪的细节我只能暂时熟视无睹（比如“排查内存占用过高”一节中，为什么实际申请的是 1.5 GiB 内存），如果这些奇怪的细节你也发现了，并痛斥我假装睁眼瞎，那么我的目的就达到了……

——还记得吗，本文的目的是让你熟悉使用 pprof，消除对它的陌生感，并能借用它进一步得了解 golang。而你通过这次试验，发现了程序的很多行为不同于你以往的认知或假设，并抱着好奇心，开始向比深处更深处迈进，那么，我何尝不觉得这是本文的功德呢？

与君共勉。
