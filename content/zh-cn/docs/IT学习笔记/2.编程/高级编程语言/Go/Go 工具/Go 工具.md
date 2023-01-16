---
title: Go 工具
---

# 概述

> 参考：
> - [官方文档，命令文档](https://go.dev/doc/cmd)

go 是用来管理 Go 编程语言源代码的工具

# go

> 参考：
> - [Go 包，标准库-cmd-go](https://pkg.go.dev/cmd/go)

go 是一个工具，用来管理 Go 语言编写的代码。该工具由多个子命令组成。每个子命令可以实现不同类型的功能。

## bug # start a bug report

## [build](https://pkg.go.dev/cmd/go#hdr-Compile_packages_and_dependencies) # 编译 package 及其依赖

默认编译当前路径下的代码包及其依赖，生成一个可执行文件
OPTIONS

- **-o <NAME>** # 指定构建完成后生成的文件名为 NAME
- **-x** # 输出 Go 程序编译、链接、打包的全过程。包括都使用了哪些库、执行了什么命令、等等

EXAMPLE

- 指定构建名称
  - go build -o jhs_cli cmd/jhs_cli/main.go

## clean # remove object files and cached files

EXAMPLE

- go clean -i github.com/spf13/cobra/cobra #

doc show documentation for package or symbol

env print Go environment information

fix update packages to use new APIs

fmt gofmt (reformat) package sources

generate generate Go files by processing source

## get #下载并安装 package 及其依赖

OPTIONS

- -u #更新现有依赖，强制更新它所依赖的其他全部模块，不包括自身
- -t #更新所有直接依赖和间接依赖的模块版本，包括单元测试中用到的。

install compile and install packages and dependencies

list list packages or modules

## mod # go 模块维护与管理命令

详见《[Go Module](✏IT 学习笔记/👨‍💻2.编程/高级编程语言/Go/Go%20 环境安装与使用/Go%20Module.md 环境安装与使用/Go Module.md)》章节

## run # 编译并运行 Go 程序

## test # test packages

详见 《[Go 单元测试](✏IT 学习笔记/👨‍💻2.编程/高级编程语言/Go/Go%20 单元测试.md 单元测试.md)》 章节

## tool # run specified go tool

## vet # report likely mistakes in packages

# 其他工具

很多 Go 语言生态的工具为我们编写代码提供了强大的支持，这些工具通常会作为 IDE 的插件被安装
比如 VSCode 中，当我们安装完 Go 的所有工具后，右键点击代码会出现如下提示：
![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/gy06h4/1654832963071-167da116-2e44-4818-a22c-37dc041ebafc.png)
通过这些工具，我们可以

- 检查代码问题
- 自动创建测试代码
- 自动格式化代码
- 等等

## gopls

> 参考：
> - [GitHub 项目，golang/tools](https://github.com/golang/tools)
> - [VSCode 建议你启用 gopls，它到底是个什么东东](https://www.modb.pro/db/87143)

gopls 是一个用以实现 [LSP](https://en.wikipedia.org/wiki/Language_Server_Protocol) 的官方工具。

## gotests

> 参考：
> - [GitHub 项目，cweill/gotests](https://github.com/cweill/gotests)

gotests 工具可以让我们更容易得编写 Go 单元测试。该工具可以根据目标源文件的 函数 和 方法 自动生成测试用例。测试文件中的任何新依赖项都会自动导入。
gotests 可以作为 IDE 的插件提供更方便的使用，下面是一个 Sublime Text3 插件的示例
![](https://notes-learning.oss-cn-beijing.aliyuncs.com/gy06h4/1654843875110-6dbf3a8d-1512-4022-bb7d-210084311509.gif)
在 [Emacs](https://github.com/damienlevin/GoTests-Emacs), also [Emacs](https://github.com/s-kostyaev/go-gen-test), [Vim](https://github.com/buoto/gotests-vim), [Atom Editor](https://atom.io/packages/gotests), [Visual Studio Code](https://github.com/Microsoft/vscode-go), and [IntelliJ Goland](https://www.jetbrains.com/help/go/run-debug-configuration-for-go-test.html) 等 IDE 中也有这个插件。当然，如果不想在 IDE 中使用，也可以在命令行直接使用 gotests。、、、、、、

### 简单示例

假如有一个文件 unit_tests.go 如下代码：

```go
func UnitTests(needArgs string) bool {
	if needArgs == "unittests" {
		return true
	} else {
		return false
	}
}
```

gotests 将会创建一个 unit_tests_test.go 文件，并为 UnitTests() 函数生成测试用例：

```go
func TestUnitTests(t *testing.T) {
	// 这是是测试时需要传递给 UnitTests() 的参数
	type args struct {
		needArgs string
	}
	// 可以创建多个测试
	tests := []struct {
		// 测试名称
		name string
		// 需要传递给 UnitTests() 的参数
		args args
		// 需要判断 UnitTests() 的返回值
		want bool
	}{
		// TODO: 在这里写具体的测试用例，也就是执行 UnitTests() 时想要传递的参数和想要获取到的返回值
		// 这是一个 struct 类型的数组，注意书写格式。
		{
			name: "这是第一条测试在下面填写测试想要传递的参数以及想要获取到的返回值",
			args: args{"unittests"},
			want: true,
		},
		{
			name: "这里是第二条测试用例中需要用到的信息",
			args: args{"这里的参数会导致返回值为 false,进而会导致本次测试失败"},
			want: true,
		},
	}
	// 执行我们提供的每一条测试用例
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := UnitTests(tt.args.needArgs); got != tt.want {
				// 如果 UnitTests() 的返回值与我们填写的 want(即想要获得的返回值) 不一致，那么将会报错
				t.Errorf("UnitTests() = %v, want %v", got, tt.want)
			}
		})
	}
}
```

执行测试后效果如下：
![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/gy06h4/1654845554292-4eb67713-5aac-400a-894f-55242b3fe799.png)

### Syntax(语法)

**gotests \[OPTIONS] PATH...**

## dlv

> 参考：
> - [GitHub 项目，go-delve/delve](https://github.com/go-delve/delve)

Delve 是 Go 编程语言的调试器。该项目的目标是为 Go 提供一个简单、功能齐全的调试工具。 Delve 应该易于调用和使用。如果您使用的是调试器，那么事情可能不会如您所愿。考虑到这一点，Delve 应该尽可能地远离你。

## impl

> 参考：
> - [GitHub 项目，josharian/impl](https://github.com/josharian/impl)

impl 用于生成实现接口的 [Method stub](<✏IT 学习笔记/👨‍💻2.编程/Programming(编程)/Programming(编程).md>>)

### 简单示例

通过 Go: Generate Interface Stubs 可以快速生成某个接口下的方法
![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/gy06h4/1642038271876-e8806fd8-1531-4c24-b1f4-a7f4d9ae306a.png)
语法：`VAR *STRUCT INTERFACE`
比如，我想让 `File` 结构体实现 `io.Closer` 接口，则输入:`f *File io.Closer`，将会生成如下方法：

```go
func (f *File) Close() error {
	panic("not implemented") // TODO: Implement
}
```

> 也可以通过命令行，使用 `impl 'f *File' io.Closer` 命令生成方法。

若提示 `Cannot stub interface: unrecognized interface: handler.YuqeData`导致无法生成方法，则对接口使用一下 `Find All Implementations`
![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/gy06h4/1642045391841-a1d01b46-deda-4561-b9b6-de714d8ee672.png)

## gomodifytags

> 参考：
> - [GitHub 项目，fatih/gomodifytags](https://github.com/fatih/gomodifytags)

## staticcheck

> 参考：
> - [GitHub 项目，dominikh/go-tools](https://github.com/dominikh/go-tools)
> - [官网](https://staticcheck.io/)

Staticcheck 是一个高级 Go Linter，即用于 Go 的代码检查工具，使用静态分析，可以发现错误和性能问题，提供简化，并强制执行样式规则

## goplay

> 参考：
> - [GitHub 项目，haya14busa/goplay](https://github.com/haya14busa/goplay)

goplay 可以让代码通过 <https://play.golang.org/> 打开（这是一个在线运行 Go 代码的网站）。
