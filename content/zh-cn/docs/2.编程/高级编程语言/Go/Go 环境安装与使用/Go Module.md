---
title: Go Module
weight: 2
---

# 概述

> 参考：
>
> - [官方文档，参考 - Go Modules 参考](https://go.dev/ref/mod)
> - [公众号，Go Modules 终极入门](https://mp.weixin.qq.com/s/6gJkSyGAFR0v6kow2uVklA)

**Go Module(Go 模块)** 是实现 [Modular Programming(模块化编程)](https://en.wikipedia.org/wiki/Modular_programming) 的工具。是 Go 语言中正式官宣的项目依赖解决方案，Go modules（前身为 vgo）发布于 Go1.11，成长于 Go1.12，丰富于 Go1.13，正式于 Go1.14 已经准备好，并且可以用在生产上（ready for production）了，Go 官方也鼓励所有用户从其他依赖项管理工具迁移到 Go modules。

module 是一个相关 Go 包的集合，它是源代码更替和版本控制的单元。模块由源文件形成的 go.mod 文件的根目录定义，包含 go.mod 文件的目录也被称为模块根。moudles 取代旧的的基于 GOPATH 方法来指定在工程中使用哪些源文件或导入包。模块路径是导入包的路径前缀，go.mod 文件定义模块路径，并且列出了在项目构建过程中使用的特定版本。

使用 Go Module 时，GOPATH 不再用于解析导入。但是，它仍然用于存储下载的源代码（在$GOPATH/pkg/mod 中）和编译的命令（在 GOPATH / bin 中）。

当程序编译时，会读取 go.mod 文件中的路径，来加载其编译所需的各种库

Go moudles 目前集成在 Go 的工具链中，只要安装了 Go，自然而然也就可以使用 Go moudles 了，而 Go modules 的出现也解决了在 Go1.11 前的几个常见争议问题：

- Go 语言长久以来的依赖管理问题。
- “淘汰”现有的 [GOPATH 的使用模式(即.解析导入能力)](https://pkg.go.dev/cmd/go#hdr-GOPATH_environment_variable)。
- 统一社区中的其它的依赖管理工具（提供迁移功能）。

# GOPATH

Go Module 出现后，GOPATH 路径变为纯粹的第三方依赖库的保存路径。目录结构通常如下：

```bash
~]# tree -L 3 $GOPATH
.
├── bin
│   ├── godef
│   ├── gomodifytags
│   ├── go-outline
│   ├── gopkgs
│   ├── goplay
│   ├── gopls
│   ├── gotests
│   └── staticcheck
└── pkg
    ├── mod
    │   ├── cache
    │   ├── github.com
    │   └── gorm.io
    └── sumdb
        └── sum.golang.org

```

bin/ 下是安装的某些第三方依赖库后生成的可执行文件

pkg/ 下是安装的第三方库

## 早期 GOPATH 模式痛点

我们先看看第一个问题，GOPATH 是什么，我们可以输入如下命令查看：

```bash
$ go env
GOPATH="/root/go"
...
```

我们输入 go env 命令行后可以查看到 GOPATH 变量的结果，我们进入到该目录下进行查看，如下：

```bash
go
├── bin
├── pkg
└── src
    ├── github.com
    ├── golang.org
    ├── google.golang.org
    ├── gopkg.in
....
```

GOPATH 目录下一共包含了三个子目录，分别是：

- bin：存储所编译生成的二进制文件。
- pkg：存储预编译的目标文件，以加快程序的后续编译速度。
- src：存储所有项目的源代码。在编写 Go 应用程序，程序包和库时，一般会以$GOPATH/src/github.com/foo/bar 的路径进行存放。

因此在使用 GOPATH 模式下，我们需要将项目代码存放在固定的 `$GOPATH/src` 目录下，并且如果执行 go get 来拉取外部依赖会自动下载并安装到 `$GOPATH` 目录下。

**为什么弃用 GOPATH 模式？**

在 GOPATH 的 `$GOPATH/src/` 下进行 .go 文件或源代码的存储，我们可以称其为 GOPATH 的模式，这个模式，看起来好像没有什么问题，那么为什么我们要弃用呢，参见如下原因：

- GOPATH 模式下没有版本控制的概念，具有致命的缺陷，至少会造成以下问题：
  - 在执行 go get 的时候，你无法传达任何的版本信息的期望，也就是说你也无法知道自己当前更新的是哪一个版本，也无法通过指定来拉取自己所期望的具体版本。
  - 在运行 Go 应用程序的时候，你无法保证其它人与你所期望依赖的第三方库是相同的版本，也就是说在项目依赖库的管理上，你无法保证所有人的依赖版本都一致。
  - 你没办法处理 v1、v2、v3 等等不同版本的引用问题，因为 GOPATH 模式下的导入路径都是一样的，都是 github.com/foo/bar。
- Go 语言官方从 Go1.11 起开始推进 Go modules（前身 vgo），Go1.13 起不再推荐使用 GOPATH 的使用模式，Go modules 也渐趋稳定，因此新项目也没有必要继续使用 GOPATH 模式。

在 GOPATH 模式下的产物

Go1 在 2012 年 03 月 28 日发布，而 Go1.11 是在 2018 年 08 月 25 日才正式发布（数据来源：GitHub Tag），在这个空档的时间内，并没有 Go modules 这一个东西，最早期可能还好说，因为刚发布，用的人不多，所以没有明显暴露，但是后期 Go 语言使用的人越来越多了，那怎么办？

这时候社区中逐渐的涌现出了大量的依赖解决方案，百花齐放，让人难以挑选，其中包括我们所熟知的 vendor 目录的模式，以及曾经一度被认为是“官宣”的 dep 的这类依赖管理工具。

但为什么 dep 没有正在成为官宣呢，其实是因为随着 Russ Cox 与 Go 团队中的其他成员不断深入地讨论，发现 dep 的一些细节似乎越来越不适合 Go，因此官方采取了另起 proposal 的方式来推进，其方案的结果一开始先是释出 vgo（Go modules 的前身，知道即可，不需要深入了解），最终演变为我们现在所见到的 Go modules，也在 Go1.11 正式进入了 Go 的工具链。

因此与其说是 “在 GOPATH 模式下的产物”，不如说是历史为当前提供了重要的教训，因此出现了 Go modules。

# Go Module 的使用和管理

可以这么说，一个自己新建的项目，就是一个模块，一个模块就是一个目录下的所有文件的集合。所以才说一个模块就是一个 Go Package 的合集。

## Go Module 相关环境变量

```bash
$ go env
GO111MODULE="auto" # 使用module功能必须要让该变量变为on
GOPROXY="https://proxy.golang.org,direct"
GONOPROXY=""
GOSUMDB="sum.golang.org"
GONOSUMDB=""
GOPRIVATE=""
...
```

### GO111MODULE

Go 语言提供了 GO111MODULE 这个环境变量来作为 Go modules 的开关，其允许设置以下参数：

- auto：只要项目包含了 go.mod 文件的话启用 Go modules，目前在 Go1.11 至 Go1.14 中仍然是默认值。
- on：启用 Go modules，推荐设置，将会是未来版本中的默认值。
- off：禁用 Go modules，不推荐设置。

GO111MODULE 的小历史

你可能会留意到 GO111MODULE 这个名字比较 “奇特”，实际上在 Go 语言中经常会有这类阶段性的变量， GO111MODULE 这个命名代表着 Go 语言在 1.11 版本添加的，针对 Module 的变量。

像是在 Go1.5 版本的时候，也发布了一个系统环境变量 GO15VENDOREXPERIMENT，作用是用于开启 vendor 目录的支持，当时其默认值也不是开启，仅仅作为 experimental。其随后在 Go1.6 版本时也将默认值改为了开启，并且最后作为了 official，GO15VENDOREXPERIMENT 系统变量就退出了历史舞台。

而未来 GO111MODULE 这一个系统环境变量也会面临这个问题，也会先调整为默认值为 on（曾经在 Go1.13 想想改为 on，并且已经合并了 PR，但最后因为种种原因改回了 auto），然后再把 GO111MODULE 的支持给去掉，我们猜测应该会在 Go2 将 GO111MODULE 给去掉，因为如果直接去掉 GO111MODULE 的支持，会存在兼容性问题。

### GOPROXY

这个环境变量主要是用于设置 Go 模块代理（Go module proxy），其作用是用于使 Go 在后续拉取模块版本时能够脱离传统的 VCS 方式，直接通过镜像站点来快速拉取。

GOPROXY 的默认值是：<https://proxy.golang.org,direct>，这有一个很严重的问题，就是 proxy.golang.org 在国内是无法访问的，因此这会直接卡住你的第一步，所以你必须在开启 Go modules 的时，同时设置国内的 Go 模块代理，执行如下命令：

```bash
go env -w GOPROXY=https://goproxy.cn,direct
```

GOPROXY 的值是一个以英文逗号 “,” 分割的 Go 模块代理列表，允许设置多个模块代理，假设你不想使用，也可以将其设置为 “off” ，这将会禁止 Go 在后续操作中使用任何 Go 模块代理。

direct 是什么

而在刚刚设置的值中，我们可以发现值列表中有 “direct” 标识，它又有什么作用呢？

实际上 “direct” 是一个特殊指示符，用于指示 Go 回源到模块版本的源地址去抓取（比如 GitHub 等），场景如下：当值列表中上一个 Go 模块代理返回 404 或 410 错误时，Go 自动尝试列表中的下一个，遇见 “direct” 时回源，也就是回到源地址去抓取，而遇见 EOF 时终止并抛出类似 “invalid version: unknown revision...” 的错误。

### GOSUMDB

它的值是一个 Go checksum database，用于在拉取模块版本时（无论是从源站拉取还是通过 Go module proxy 拉取）保证拉取到的模块版本数据未经过篡改，若发现不一致，也就是可能存在篡改，将会立即中止。

GOSUMDB 的默认值为：sum.golang.org，在国内也是无法访问的，但是 GOSUMDB 可以被 Go 模块代理所代理（详见：Proxying a Checksum Database）。

因此我们可以通过设置 GOPROXY 来解决，而先前我们所设置的模块代理 goproxy.cn 就能支持代理 sum.golang.org，所以这一个问题在设置 GOPROXY 后，你可以不需要过度关心。

另外若对 GOSUMDB 的值有自定义需求，其支持如下格式：

- 格式 1：+。
- 格式 2：+ 。

也可以将其设置为 “off”，也就是禁止 Go 在后续操作中校验模块版本。

### GOPRIVATE/GONOPROXY/GONOSUMDB

这三个环境变量都是用在当前项目依赖了私有模块，例如像是你公司的私有 git 仓库，又或是 github 中的私有库，都是属于私有模块，都是要进行设置的，否则会拉取失败。

更细致来讲，就是依赖了由 GOPROXY 指定的 Go 模块代理或由 GOSUMDB 指定 Go checksum database 都无法访问到的模块时的场景。

而一般建议直接设置 GOPRIVATE，它的值将作为 GONOPROXY 和 GONOSUMDB 的默认值，所以建议的最佳姿势是直接使用 GOPRIVATE。

并且它们的值都是一个以英文逗号 “,” 分割的模块路径前缀，也就是可以设置多个，例如：

```bash
go env -w GOPRIVATE="git.example.com,github.com/eddycjy/mquote"
```

设置后，前缀为 git.xxx.com 和 github.com/eddycjy/mquote 的模块都会被认为是私有模块。

如果不想每次都重新设置，我们也可以利用通配符，例如：

```bash
go env -w GOPRIVATE="*.example.com"
```

这样子设置的话，所有模块路径为 example.com 的子域名（例如：git.example.com）都将不经过 Go module proxy 和 Go checksum database，需要注意的是不包括 example.com 本身。

具体使用步骤：

- 通过 go 命令行，进入到你当前的工程目录下，在命令行设置临时环境变量 set GO111MODULE=on；
- 执行命令 go mod init NAME 在当前目录下生成一个 go.mod 文件，执行这条命令时，当前目录不能存在 go.mod 文件。如果之前生成过，要先删除；
- 如果你工程中存在一些不能确定版本的包，那么生成的 go.mod 文件可能就不完整，因此继续执行下面的命令；
- 执行 go mod tidy 命令，它会添加缺失的模块以及移除不需要的模块。执行后会生成 go.sum 文件(模块下载条目)。添加参数-v，例如 go mod tidy -v 可以将执行的信息，即删除和添加的包打印到命令行；
- 执行命令 go mod verify 来检查当前模块的依赖是否全部下载下来，是否下载下来被修改过。如果所有的模块都没有被修改过，那么执行这条命令之后，会打印 all modules verified。
- 执行命令 go mod vendor 生成 vendor 文件夹，该文件夹下将会放置你 go.mod 文件描述的依赖包，文件夹下同时还有一个文件 modules.txt，它是你整个工程的所有模块。在执行这条命令之前，如果你工程之前有 vendor 目录，应该先进行删除。同理 go mod vendor -v 会将添加到 vendor 中的模块打印出来；

## go.mod 文件

> 参考：
>
> - [官方文档，参考 - Go Modules 参考 - go.mod 文件](https://go.dev/ref/mod#go-mod-file)

go.mod 文件定义 module 路径以及列出其他需要在 build 时引入的模块的特定的版本。例如下面的例子中，go.mod 声明 example.com/m 路径是 module 的根目录，同时也声明了 module 依赖特定版本的 golang.org/x/text 和 gopkg.in/yaml.v2。

go.mod 文件中有如下几个关键字：

- **module** # 定义 module 路径，该路径不用与当前路径相同，只是 module 所用的一个名称，可以代指当前目录。(比如/root/desistdaydream/cobra/目录下，创建一个 go.mod 文件，可以定义 module 路径为 cobratest，这个 cobratest 模块路径名，就表示/root/desistdaydream/cobra/这个目录)to define the module path;
- **go** # 设置期望的 Go 语言版本
- **require** # to require a particular module at a given version or later;
- **exclude** # to exclude a particular module version from use; and
- **replace** # replace 指令可以将特定版本的模块或模块的所有版本的内容替换为其他位置的内容。替换可以通过另一个模块路径和版本，或者特定平台的文件路径来指定。
    - Notes: 在不方便导入互联网包的时候，可以用来导入本地文件系统中的其他项目。

```go

module github.com/eddycjy/module-repo

go 1.13

require (
    github.com/eddycjy/mquote v0.0.0-20200220041913-e066a990ce6f
)
```

go.mod 文件还可以指定要替换和排除的版本，命令行会自动根据 go.mod 文件来维护需求声明中的版本。如果想获取更多的有关 go.mod 文件的介绍，可以使用命令 go help go.mod。

go.mod 文件用 // 注释，而不用 /\*\*/。文件的每行都有一条指令，由一个动作加上参数组成。例如：

```go
module my/thing
require other/thing  v1.0.2
require new/thing   v2.3.4
exclude old/thing   v1.2.3
replace bad/thing   v1.4.5  => good/thing v1.4.5
```

上面三个动词 require、exclude、replace 分别表示：项目需要的依赖包及版本、排除某些包的特别版本、取代当前项目中的某些依赖包。

相同动作的命令可以放到一个动词+括号组成的结构中，例如：

```go
require (
    new/thing v2.3.4
    old/thing v1.2.3
)
```

### replace

https://go.dev/ref/mod#go-mod-file-replace

replace 指令将 Module 内容（特定或所有版本）替换为其它地方找到的内容。其它地方找到内容可以是：

- 网络上的 Module
- 文件系统中包含 go.mod 文件的目录
- etc.

### 其他命令的支持

旧的版本，构建编译命令 `go build` 中的参数没有 `-mod` 参数，最新的版本现在多了这个，用来对 `go.mod` 文件进行更新或其他使用控制。形式如：`go build -mod [mode]`，其中 mode 有以下几种取值：readonly，release，vendor。当执行 `go build -mod=vendor` 的时候，会在生成可执行文件的同时将项目的依赖包放到主模块的 `vendor` 目录下。

`go get -m [packages]` 会将下载的依赖包放到 `GOPATH/pkg/mod` 目录下，并且将依赖写入到 `go.mod` 文件。`go get -u=patch` 会更新主模块下的所有依赖包。

如果遇到不熟悉的导入包，任何可以查找包含该引入包模块的 `go` 命令，都会自动将该模块的最新版本添加到 `go.mod` 文件中。同时也会添加缺失的模块，以及删除无用的 module。例如：go build, go test 或者 go list 命令。另外，有一个专门的命令 `go mod tidy`，用来查看和添加缺失的 module 需求声明以及移除不必要的。

`go.mod` 文件是可读，也是可编辑的。`go` 命令行会自动更新 `go.mod` 文件来维持一个标准格式以及精确的引入声明。

## go.sum 文件

在第一次拉取模块依赖后，会发现多出了一个 go.sum 文件，其详细罗列了当前项目直接或间接依赖的所有模块版本，并写明了那些模块版本的 SHA-256 哈希值以备 Go 在今后的操作中保证项目所依赖的那些模块版本不会被篡改。

```go
github.com/eddycjy/mquote v0.0.1 h1:4QHXKo7J8a6J/k8UA6CiHhswJQs0sm2foAQQUq8GFHM=
github.com/eddycjy/mquote v0.0.1/go.mod h1:ZtlkDs7Mriynl7wsDQ4cU23okEtVYqHwl7F1eDh4qPg=
github.com/eddycjy/mquote/module/tour v0.0.1 h1:cc+pgV0LnR8Fhou0zNHughT7IbSnLvfUZ+X3fvshrv8=
github.com/eddycjy/mquote/module/tour v0.0.1/go.mod h1:8uL1FOiQJZ4/1hzqQ5mv4Sm7nJcwYu41F3nZmkiWx5I=
...
```

我们可以看到一个模块路径可能有如下两种：

```go
github.com/eddycjy/mquote v0.0.1 h1:4QHXKo7J8a6J/k8UA6CiHhswJQs0sm2foAQQUq8GFHM=
github.com/eddycjy/mquote v0.0.1/go.mod h1:ZtlkDs7Mriynl7wsDQ4cU23okEtVYqHwl7F1eDh4qPg=
```

h1 hash 是 Go modules 将目标模块版本的 zip 文件开包后，针对所有包内文件依次进行 hash，然后再把它们的 hash 结果按照固定格式和算法组成总的 hash 值。

而 h1 hash 和 go.mod hash 两者，要不就是同时存在，要不就是只存在 go.mod hash。那什么情况下会不存在 h1 hash 呢，就是当 Go 认为肯定用不到某个模块版本的时候就会省略它的 h1 hash，就会出现不存在 h1 hash，只存在 go.mod hash 的情况。

# go mod CLI

go mod 提供了一系列操作模块的命令，所有的 go 命令中现在已经内置了对 module 的支持，而不仅仅是 go mod 命令。例如使用 go get 时，会经常自动在后台添加、移除、升级、降级依赖包版本。

## Syntax(语法)

**go mod \[ARGUMENTS]**

COMMAND：

- **download** # 下载模块到本地缓存，具体可以通过命令 go env 查看，其中环境变量 GOCACHE 就是缓存的地址，如果该文件夹的内容太大，可以通过命令 go clean -cache
- **edit** # 从工具或脚本中编辑 go.mod 文件
- **graph** # 打印模块需求图
- **init** # 在当前目录下初始化新的模块
- **tidy** # 添加缺失的模块以及移除无用的模块
- **vendor** # 导出项目所有的依赖到 vendor 目录
- **verify** # 验证依赖项是否达到预期的目的
- **why** # 查看为什么需要包或模块

## go mod download

**go mod download \[-dir] \[-json] \[modules]**

使用此命令来下载指定的模块，模块的格式可以根据主模块依赖的形式或者 path@version 形式指定。如果没有指定参数，此命令会将主模块下的所有依赖下载下来。

go mod download 命令非常有用，主要用来预填充本地缓存或者计算 Go 模块代理的回答。默认情况下，下载错误会输出到标准输出，正常情况下没有任何输出。-json 参数会以 JSON 的格式打印下载的模块对象，对应的 Go 对象结构是这样。

type Module struct { Path string //module path Version string //module version Error string //error loading module Info string //absolute path to cached .info file GoMod string //absolute path to cached .mod file Zip string //absolute path to cached .zip file Dir string //absolute path to cached source root directory Sum string //checksum for path, version (as in go.sum) GoModSum string //checksum for go.mod (as in go.sum)}

## go mod init

**go mod init \[ModuleName]**

一般情况 ModuleName 是以后 import 时所使用的路径

此命令会在当前目录中初始化和创建一个新的 go.mod 文件，当然你也可以手动创建一个 go.mod 文件，然后包含一些 module 声明，这样就比较麻烦。go mod init 命令可以帮助我们自动创建

例如：`go mod init example.com/m`

使用这条命令时，go.mod 文件必须提前不能存在。初始化会根据引入包声明来推测模块的路径或者如果你工程中之前已经存在一些依赖包管理工具，例如 godep，glide 或者 dep。那么 go mod init 同样也会根据依赖包管理配置文件来推断。

## go mod tidy

**go mod tidy \[-v]**

默认情况下，Go 不会移除 go.mod 文件中的无用依赖。所以当你的依赖中有些使用不到了，可以使用 go mod tidy 命令来清除它。

它会添加缺失的模块以及移除不需要的模块。执行后会生成 go.sum 文件(模块下载条目)。添加参数-v，例如 go mod tidy -v 可以将执行的信息，即移除的模块打印到标准输出。

## go mod vendor

go mod vendor \[-v]

此命令会将 build 阶段需要的所有依赖包放到主模块所在的 vendor 目录中，并且测试所有主模块的包。同理 go mod vendor -v 会将添加到 vendor 中的模块打印到标准输出。

## go mod verify

此命令会检查当前模块的依赖是否已经存储在本地下载的源代码缓存中，以及检查自从下载下来是否有修改。如果所有的模块都没有修改，那么会打印 all modules verified，否则会打印变化的内容。

虚拟版本号

go.mod 文件和 go 命令通常使用语义版本作为描述模块版本的标准形式，这样可以比较不同版本的先后顺序。例如模块的版本是 v1.2.3，那么通过重新对版本号进行标签处理，得到该版本的虚拟版本。形式如：v0.0.0-yyyymmddhhmmss-abcdefabcdef。其中时间是提交时的 UTC 时间，最后的后缀是提交的哈希值前缀。时间部分确保两个虚拟版本号可以进行比较，以确定两者顺序。

下面有三种形式的虚拟版本号：

- vX.0.0-yyyymmddhhmmss-abcdefabcdef，这种情况适合用在在目标版本提交之前 ，没有更早的的版本。（这种形式本来是唯一的形式，所以一些老的 go.mod 文件使用这种形式）
- vX.Y.Z-pre.0.yyyymmddhhmmss-abcdefabcdef，这种情况被用在当目标版本提交之前的最新版本提交是 vX.Y.Z-pre。
- vX.Y.(Z+1)-0.yyyymmddhhmmss-abcdefabcdef，同理，这种情况是当目标版本提交之前的最新版本是 vX.Y.Z。

虚拟版本的生成不需要你去手动操作，go 命令会将接收的 commit 哈希值自动转化为虚拟版本号。

# 最佳实践

## 获取私有仓库的包

参考 [GOPRIVATE/GONOPROXY/GONOSUMDB](#GOPRIVATE/GONOPROXY/GONOSUMDB)

TODO: 是否需要改下面的配置文件待确认

修改本地 `.gitconfig` 文件

```ini
# 添加信息
[url "ssh://git@github.com/"]
  insteadOf=https://github.com/
```

修改后再使用 `go mod download` 或者 `go mod tidy` 就可以正常下载文件了

## 使用本地文件系统中的包

```bash
├── my-project
│   ├── go.mod
│   └── main.go
└── my-locallib
    ├── go.mod
    └── mypackage.go
```

假如我有一个名为 my-locallib 的本地项目，想要在 my-project 项目中导入 my-locallib 作为依赖包（不通过互联网，有敏感信息）

1. 在项目的 `go.mod` 文件中使用 `replace` 指令：

```go
module my-project

go 1.24

require (
    need-import-my-locallib v0.0.0
)

// replace 包名 => 目录路径
replace need-import-my-locallib => /path/to/my-locallib
```

> [!Tip] 路径中最后的目录名与导入的包名并不需要完全一样
>
> 另外，replace 的目录路径除了使用绝对路径，还可以使用相对路径（避免多平台（windows, linux, etc.）的项目路径格式不一致），比如：
>
> `replace need-import-my-locallib => ../my-locallib`

> [!Attention] replace 指定的目录必须是一个使用 go mod init 初始化成功的项目根目录。（i.e. 必须包含 go.mod 之类的文件）

2. 然后在代码中导入：

```go
import "need-import-my-locallib"

// 若想使用 my-locallib 库中 pkg/utils/ 目录下的包，则可以这么导入
import "need-import-my-locallib/pkg/utils/"
```



