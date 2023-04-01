---
title: Go 环境安装与使用
weight: 1
---

# 概述

> 参考：
> - [官方文档，下载并安装 Go](https://golang.org/doc/install)
> - [官方文档，安装多个版本的 Go](https://golang.org/doc/manage-install)
> - [Go 包，标准库-cmd-go](https://pkg.go.dev/cmd/go)
> - [GitHub 项目-goproxy-goproxy.io 与 goproxy.cn 说明](https://github.com/goproxy/goproxy.cn/issues/61)

我们可以通过一个名为 go 的二进制文件实现绝大部分日常的 编码、编译 等工作，只要安装好 Go 的环境即可。

# 安装 Go

## Linux 安装

从[官网](https://golang.org/dl/)下载 linux 版的 `.tar.gz` 包

```bash
export GoVersion=1.20.2
wget https://go.dev/dl/go${GoVersion}.linux-amd64.tar.gz
sudo tar -C /usr/local -xvzf go${GoVersion}.linux-amd64.tar.gz
```

配置环境变量，以便让 shell 可以执行 go 命令并立刻生效

```bash
sudo tee /etc/profile.d/go.sh > /dev/null <<-"EOF"
# export GOPATH=/opt/gopath
# export PATH=$PATH:\$GOPATH/bin:/usr/local/go/bin
export PATH=$PATH:/usr/local/go/bin
export GO111MODULE=on
export GOPROXY=https://goproxy.cn,https://goproxy.io,direct
export CGO_ENABLED=0
EOF

source /etc/profile.d/go.sh
```

> CGO_ENABLED 开启后 Go 代码最终编译的可执行文件都是要有外部依赖的。不过我们依然可以通过 disable CGO_ENABLED 来编译出纯静态的 Go 程序，常用于交叉编译
> CGO_ENABLED 关闭即可编译出纯静态的 Go 程序，可以用于 alpine 镜像中。

## Windows 安装

从[官网](https://golang.org/dl/)下载 Windows 版的 msi 安装包

双击安装 Golang

配置环境变量，执行命令

- go env -w GOPATH=D:\Tools\GoPath
- go env -w GOPROXY=https://goproxy.cn,https://goproxy.io,direct
- go env -w GO111MODULE=on

## 安装多个版本的 Go

获取其余版本的 golang

- go install golang.org/dl/goX.XX.X@latest
- goX.XX.X download

新下载的 golang 版本可以像这样使用，在 go 后面加上版本号

- goX.XX.X version

通过 goX.XX.X env 命令查看该 go 版本的变量，可以看到，默认的 GOROOT 是在 ~/sdk/goX.XX.X 目录下的

所以如果想要切换默认的 go 版本(比如某个程序调用 go 命令时)，只需要将环境变量 GOROOT 指向该目录即可，之后使用 go 命令即可为当前指定的版本

# 初始化项目

Go 的项目通常由 [Go Module](https://www.typescriptlang.org/) 管理，项目目录中通常必须包含如下几个文件：

- go.mod
- go.sum

go 相关工具通过 `go.mod` 与 `go.sum` 两个文件管理项目及其依赖

使用 `go mod init <NAME>` 命令在当前目录会创建一个 go.mod 文件。有任何新的 import，都可以通过 `go mod tidy` 生成依赖文件再生成 `go.sum` 文件。

# 编译 Go

若想在 Windows 中编译依赖 gcc 的项目，则需要安装 gcc 编译器（i.e.[MinGW-w64](https://sourceforge.net/projects/mingw-w64/)）在 file 标签中，下载 [x86_64-win32-seh](https://sourceforge.net/projects/mingw-w64/files/Toolchains%20targetting%20Win64/Personal%20Builds/mingw-builds/8.1.0/threads-win32/seh/x86_64-8.1.0-release-win32-seh-rt_v6-rev0.7z) 这个版本并安装即可。这是一个 tar 包，解压完成后，需要在 Windows 的 ${PATH} 环境变量中，添加解压出来的 bin 目录，通常都在 `PATH\TO\x86_64-8.1.0-release-win32-seh-rt_v6-rev0\mingw64\bin` 这里

## 交叉编译

Linux 下编译出 Windows 的程序。需要安装 Windows 版的 gcc 工具（i.e.[MinGW-w64](https://www.mingw-w64.org/)）(有的环境还需要安装 gcc-multilib 包)

```bash
sudo apt-get install gcc-mingw-w64
GOOS=windows GOARCH=amd64 CGO_ENABLED=1 CC=x86_64-w64-mingw32-gcc go build
```

## 在容器中编译

### golang 镜像

```bash
docker run -it -v /${YourPackageSrc}:/go/work \
  -w /go/work \
  golang:1.17 go build
```

多次使用

```bash
docker run -it --network host --name golang \
  -v /root/projects:/root/projects \
  -v /root/go:/go \
  golang:1.17 /bin/bash
```

一次性构建

```bash
docker run -it -v /root/projects/${Project}:/go/work \
  -w /go/work \
  golang:1.17 go build cmd/XX.go
```

### go-mingw 镜像

镜像，用于使用基于官方 Go Docker 镜像的 MinGW-w64 工具链为 Windows 构建 Go 二进制文件。

```bash
docker run --rm -it -v /${YourPackageSrc}:/go/work \
  -w /go/work \
  -e GOPROXY=https://goproxy.cn,https://goproxy.io,direct
  x1unix/go-mingw go build
```

# Go 关联文件与配置

Go 程序的很多关联文件都可以通过 Go 环境变量进行配置，所以绝大部分关联文件都以变量的形式记录。
**$GOPATH** # GOPATH 环境变量列出了寻找 Go 代码的位置。**同时也是存储 Go 模块的目录，即 go mod 相关命令保存数据的目录**。

- **./pkg/mod/** # 存储下载的源代码。
- **./bin/** # 存储编译的命令。下载的依赖包中若包含二进制文件，也会保存在这个目录中

## Go 环境变量

Go 通过环境变量来配置其运行行为，通过 go env 命令可以看到当前使用的环境变量：
**GO111MODULE="on|off"** # 设置是否使用 go mod，该环境变量将于 1.17 版本删除，并从此开始仅支持 go mod
GOARCH="amd64"
GOBIN=""
GOCACHE="/root/.cache/go-build"
GOENV="/root/.config/go/env"
GOEXE=""
GOFLAGS=""
GOHOSTARCH="amd64"
GOHOSTOS="linux"
GOINSECURE=""
GOMODCACHE="/root/go/pkg/mod"
GONOPROXY=""
GONOSUMDB=""
GOOS="linux"
**GOPATH=<STRING>** # 设置 gopath 所在路径。默认值：`~/go`
GOPRIVATE=""
**GOPROXY=<STRING> **# 设置 go get、go install 命令时，所使用的代理服务器。可以加快获取第三方库的速度。
**GOROOT=<STRING>** # Go 的安装路径。默认值：Go 的安装路径，Linux 中通常为 /usr/local/go
GOSUMDB="sum.golang.org"
GOTMPDIR=""
GOTOOLDIR="/usr/local/go/pkg/tool/linux_amd64"
GOVCS=""
GOVERSION="go1.16.4"
GCCGO="gccgo"
AR="ar"
CC="gcc"
CXX="g++"
**CGO_ENABLED="0"** # CGO_ENABLED 开启后 Go 代码最终编译的可执行文件都是要有外部依赖的。不过我们依然可以通过 disable CGO_ENABLED 来编译出纯静态的 Go 程序，常用于交叉编译。CGO_ENABLED 关闭即可编译出纯静态的 Go 程序，可以用于 alpine 镜像中。
GOMOD="/dev/null"
CGO_CFLAGS="-g -O2"
CGO_CPPFLAGS=""
CGO_CXXFLAGS="-g -O2"
CGO_FFLAGS="-g -O2"
CGO_LDFLAGS="-g -O2"
PKG_CONFIG="pkg-config"
GOGCCFLAGS="-fPIC -m64 -fmessage-length=0 -fdebug-prefix-map=/tmp/go-build1775394647=/tmp/go-build -gno-record-gcc-switches"

## goproxy 说明

我把老哥的 Issue 转移到这里来了哈，这个项目才是 [goproxy.cn](https://goproxy.cn/) 的源代码。[Goproxy](https://github.com/goproxy/goproxy) 是这个项目所基于的底层 Go module proxy 实现，它针对的不只是国内的开发者，所以既不建议用中文也不适合讨论 [goproxy.cn](https://goproxy.cn/)。这里讨论老哥你的问题才更为合适。
我先说一下 [goproxy.cn](https://goproxy.cn/) 和 [goproxy.io](https://goproxy.io/) 的背景。io 是由坤哥（[@oiooj](https://github.com/oiooj)）开发出来的，要早于 cn 出现几个月，目前由他跑在他所任职的公司腾讯云的香港区服务器上。cn 是由我发起的，现在完全属于七牛云，也备案在他们名下，自然所有的 CDN 资源和服务器资源都是由他们提供的，我目前跟 cn 的关系是属于它的维护者，并不是拥有者。
再说一下为什么会有两个这么相似的域名且功能类似的项目存在。我注册 cn 是去年二月底，当时是直接查询的 cn 后缀，因为我想的是这种项目肯定也就咱们中国是刚需要单独再来一套，所以并没有查询别的后缀，因为我认为 cn 实在是太合适了。并且由于当时我还在忙我的本科毕设和其他的一下事情，二月份我并没有开始开发 cn，等开发完了跟七牛云的 CEO 许叔（[@xushiwei](https://github.com/xushiwei)）谈交给他们运营时候，我才第一次听同学说到了 io。但当时我点进 io 的 GitHub 仓库一看发现并没有任何地方提到中国并且全是英文就下意识以为它是个国外项目，就没做过多研究，并且的 io 服务器当时也在美国我访问速度有些慢，就也没在意了。
最后再说一下为什么两个项目没有合并了一起发展。这个我和坤哥是有讨论过的，因为有人找上了我跟坤哥。并且坤哥也同意了最后我的提议合并了两个项目，将 cn 留作国内的公共代理，将 io 的代理类流量重定向到 cn 并后续将 io 用作一款搭代理的开源软件来面向全球提供给大家替换掉 JFrog 的一款商业产品，目前只有这一种解决方案才能保持两个功能独立且都能继续存活下去。坤哥之所以能同意我的这个提议一方面是坤哥所任职的公司加班过于严重，还有一方面是他自身没有精力维护了。然后之所以选择留 cn 做公共代理一方面是这两个域名里面只有 cn 能做备案能挂上 CDN 服务，因为这种类型的项目 CDN 服务是刚需，还有一方面是 io 这个域名后缀做全球化项目比 cn 更为合适，所以留 io 作国内代理把 cn 用作面向全球的搭代理的软件就显得很别扭了。
为什么两个项目现在没有合并呢？这个就不是因为我们两个作者了，因为我们两个作者已经达成了一致可以合并。这里面还有第三方地插足（为保其名誉我暂时不提具体是谁），其认为公共代理必须交由其所掌控的“社区”来运营，cn 已经过户给了一家商业公司无法再过户给其所掌控的“社区”，所以 cn 在其眼里就变成了一个其所描述的邪恶公司所拥有的商业产品，于是其要求我和七牛商谈放弃 cn 并全力投入为其做别的开发，于是被我拒绝。最后坤哥在中间处境比较尴尬，所以合并的事儿也就暂时搁置了。
最后，简单来总结一下就是，[goproxy.cn](https://goproxy.cn/) 和 [goproxy.io](https://goproxy.io/) 目前并无直接关系。或许之后没有了第三方地插足它们两个会合并变得有关系，但在那之前它们唯一能联系在一起的是它们都可以用作为 Go module proxy。至于哪个快、哪个稳、哪个香、用哪个，这个需要老哥你自己做判断了，我跟坤哥的关系并不差，所以我现在不会妄加评论。更何况现在 `GOPROXY` 不是支持逗号列表嘛。❤️

## GOPATH

注意：以下对 GOPATH 的理解是在 golang1.13 版本之前

GOPATH 就是 go 项目的工作目录，是开发人员写代码的目录。

GOPATH 里面一般包含 bin、pkg、src 这 3 个文件夹。

项目文件夹一般是放在 src 目录中

一般情况下，如果自己在开发多个项目，那么最好一个项目对应一个 GOPATH 路径。这时候只需要切换 GOPATH 环境变量的值，就可以编译运行对应的项目了。(比如我有两个项目目录/root/cobra 和/root/bee，这俩目录可以分别作为 GOPATH 变量的值，想开发哪个，就把 GOPATH 变量的值改为对应的目录路径)

这样做的目的主要是为了让每个项目所依赖的库等东西，可以分开而不会冲突

现在有 go module 之后，就可以不用把项目放在 GOPATH 路径下了。具体 go module 的作用详见 1.4.Go module 的介绍及使用.1 新功能 module 的介绍及使用
