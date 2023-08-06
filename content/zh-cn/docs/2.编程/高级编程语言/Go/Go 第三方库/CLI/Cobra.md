---
title: Cobra
---

# 概述

> 参考：
> - [GitHub 项目，spf13/cobra](https://github.com/spf13/cobra)
> - [官网](https://cobra.dev/)
> - <https://zhangguanzhang.github.io/2019/06/02/cobra/>

Cobra 是一个 [Go](/docs/2.编程/高级编程语言/Go/Go.md)  语言的库，其提供简单的接口来创建强大现代的 CLI 接口，类似于 git 或者 go 工具。同时，它也是一个应用，用来生成个人应用框架，从而开发以 Cobra 为基础的应用。热门的 docker 和 k8s 源码中都使用了 Cobra

Cobra 结构由三部分组成：

- **Command(命令)** #
- **Args(参数)** #
- **Flag(标志)** #

```go
type Command struct {
    Use   string  // The one-line usage message.
    Short string  // The short description shown in the 'help' output.
    Long  string  // The long message shown in the 'help<this-command>' output.
    Run   func(cmd *Command, args []string)  // Run runs the command.
    ...
}
```

传统 Linux 和 unix 的话命令规范为情况为下面几种

```bash
# 单独命令,例如date
date

# 带选项的命令
ls -l

# 选项有值
last -n 3

# 短选项合起来写,注意合起来写的时候最后一个选项以外的选项都必须是无法带值的，例如last -n 3 -R只能合起来写成下面的
last -Rn 3

# 无值的长选项
rm --force

# 带值的长选项
last --num 3
last --num=3
find -type f

# 值能追加的命令
command --host ip1 --host ip2 #命令内部能完整读取所有host做处理

# 带args的命令
rm file1 file2
cat -n file1 file2

# 多级命令
ip addr show
ip addr delete xxx

# 所有情况的命令
cmd sub_cmd1 subcmd2 --host 10.0.0.2 -nL3 -d ':' --username=admin '^a' '^b'
```

而 cobra 是针对长短选项和多级命令都支持的库，单独或者混合都是支持的，不过大多数还是用来写多级命令的 cli tool 用的。命令的格式为下列

```bash
rootCommand subcommand1 subcommand2 -X value --XXXX value -Y a -Y b --ZZ c --ZZ d args1  args2
```

前三个是不同场景下的说明，最后一个是要执行的函数

## 使用 Cobra 编写的典型项目

Cobra 用于许多 Go 项目，例如 [Kubernetes](https://kubernetes.io/)、 [Hugo](https://gohugo.io/) 和 [GitHub CLI](https://github.com/cli/cli) 等等。[此列表](https://github.com/spf13/cobra/blob/main/projects_using_cobra.md)包含更广泛的使用 Cobra 的项目列表。
<https://github.com/gohugoio/hugo>
<https://github.com/containerd/nerdctl>

# 安装与导入

安装

```bash
go get -u github.com/spf13/cobra@latest
```

导入

```go
import "github.com/spf13/cobra"
```

## Cobra 命令行工具

cobra-cli 是一个命令行程序，用于生成 Cobra 应用程序和命令文件。它将引导您的应用程序脚手架以快速开发基于 Cobra 的应用程序。这是将 Cobra 合并到您的应用程序中的最简单方法。

```bash
go install github.com/spf13/cobra-cli@latest
```

安装后会创建一个可执行文件 cobra-cli 位于 `${GOPATH}/bin` 目录中

```bash
$ go env | grep GOPATH
GOPATH="/home/lichenhao/go"
$ which cobra-cli
/home/lichenhao/go/bin/cobra-cli
```

# Cobra 的基本使用

我们使用 `go mod init github.com/DesistDaydream/go-cobra` 初始化一个项目。
Cobra 的应用程序目录结构通常如下：

```bash
$ tree
.
├── LICENSE
├── cmd
│   └── root.go
├── go.mod
├── go.sum
└── main.go
```

> cobra-cli 默认情况下，Cobra 将添加 Apache 许可证。如果您不想这样，可以将标志添加 `-l none` 到所有生成器命令。但是，它会在每个文件顶部添加 `Copyright © 2022 NAME HERE <EMAIL ADDRESS>` 这样的添加版权声明。如果通过选项 `-a YOUR NAME` 则索赔将包含您的姓名。
> **注意：使用 cobra-cli 生成的目录结构在真正使用时并不灵活，我们通常会将 XXXCmd 变量封装到函数数，以便可以对变量进行更多的处理。灵活性更大。下面的使用示例并不是生产推荐的结构和用法。**

`main.go` 文件非常简单，只有一个目的，初始化 Cobra

```go
package main

import "github.com/DesistDaydream/go-cobra/cmd"

func main() {
	cmd.Execute()
}
```

## 创建命令

`Command{}` 是 Cobra 命令的**核心结构体**，只有有了这个结构体，才能围绕命令执行方法、设置命令行标志等。

### 创建根命令(rootCmd)

根命令通常放在 `cmd/root.go` 文件中

```go
// rootCmd 表示在没有任何子命令调用的情况时的基本命令
var rootCmd = &cobra.Command{
	Use:   "go-cobra",
	Short: "这个应用简要的描述",
	Long: `横跨多行的较长描述，可能包含示例和使用应用程序的用法。 例如：
当我运行程序时，会显示该描述内容
	如果使用缩进，这行在界面展示时有缩进。`,
	Run: func(cmd *cobra.Command, args []string) {
		// 如果这个应用没有任何子命令，直接使用 go-cobra 执行的话，将会执行这里面的代码
	},
}
```

### 创建子命令

使用 `Command.AddCommand()` 方法将一个或多个命令添加到父命令中，下面的示例可以为根命令添加一个 version 子命令。

```go
func init() {
  rootCmd.AddCommand(versionCmd)
}

var versionCmd := &cobra.Command{
    Use:   "version",
    Short: "这个命令的简要描述",
    Long:  `横跨多行的较长描述，可能包含示例和使用命令的用法。`,
    Run: func(cmd *cobra.Command, args []string) {
        fmt.Println("version called")
    },
}
```

## Flag(命令行标志)

标志可以是“持久的”，这意味着该标志将可用于分配给它的命令以及该命令下的每个命令。对于全局标志，将标志分配为根上的持久标志。

```go
rootCmd.PersistentFlags().StringVarP(&rootFlags.CfgFile, "config", "c", "", "指定配置文件")
```

也可以在本地分配一个标志，它只适用于该特定命令

```go
rootCmd.Flags().BoolP("toggle", "t", false, "关于toggle标志的帮助信息")
```

实际命令都有选项，分为持久和本地，持久例如`kubectl`的`-n`可以用在很多二级命令下，本地命令选项则不会被继承到子命令。我们给 remove 添加一个移除指定名字的选项，修改`cmd/remove.go`的 init 函数：
添加 Flags 使用 `Command.Flags()` 或 `cmd.PersistentFlags()` 方法，具体有以下使用规律

- \<type>
- \<type>P
- \<type>Var
- \<type>VarP

带 P 的相对没带 P 的多了个短选项,没带 P 的选项只能用`--long-iotion`这样，而不能使用 `-l` 这种。

- 获取选项的值用`cmd.Flags().GetString("name")`

不带 Var 的获取值使用`Get<type>("FlagName")`，这样似乎非常麻烦，实际中都是用后面俩种 Var 直接传入地址自动注入的，例如

```go
var dates int32
cmd.Flags().Int32VarP(&dates,"date", "d", 1234, "this is var test")
```

- type 有 `Slice`，`Count`，`Duration`,`IP`,`IPMask`,`IPNet` 之类的类型,Slice 类型可以多个传入，直接获取就是一个切片，例如 `--master ip1 --master ip2`
- 类似 `--force` 这样的开关型选项，实际上用 Bool 类型即可，默认值设置为 false，单独给选项不带值就是 true，也可以手动传入 false 或者 true
- MarkDeprecated 告诉用户放弃这个标注位，应该使用新标志位，MarkShorthandDeprecated 是只放弃短的，长标志位依然可用。MarkHidden 隐藏标志位
- `MarkFlagRequired("region")` 表示 region 是必须的选项，不设置下选项都是可选的

## 读取配置文件

\~~类似~~`~~kubectl~~`~~ 的~~`~~~/.kube/config~~`~~ 和 ~~`~~gcloud~~`~~这些 ~~`~~cli~~`~~ 都会读取一些配置信息，也可以从命令行指定信息。细心观察的话可以看到这个是一直存在在命令帮助上的~~

    Global Flags:
          --config string   config file (default is $HOME/.cli.yaml)

\~~spf13 里的 viper 包的几个方法就是干这个的，viper 是 cobra 集成的配置文件读取的库
可以通过环境变量读取~~

    removeCmd.Flags().StringP("name", "n", viper.GetString("ENVNAME"), "The application to be executed")

\~~默认可以在 ~~`~~cmd/root.go~~`~~ 文件里看到默认配置文件是家目录下的.应用名，这里我是~~`~~$HOME/.cli.yaml~~`~~，创建并添加下面内容~~

    name: "Billy"
    greeting: "Howdy"

\~~Command 的 Run 里提取字段~~

```go
Run: func(cmd *cobra.Command, args []string) {
    greeting := "Hello"
    name, _ := cmd.Flags().GetString("name")
    if name == "" {
        name = "World"
    }
    if viper.GetString("name")!=""{
        name = viper.GetString("name")
    }
    if viper.GetString("greeting")!=""{
        greeting = viper.GetString("greeting")
    }
    fmt.Println(greeting + " " + name)
},
```

也可以将配置文件中的值绑定到命令行 Flag 里。在下面的示例中，通过 viper 包获取到的 author 的值将会绑定到命令行 Flag 的 author 中：

```go
var author string

func init() {
    rootCmd.PersistentFlags().StringVar(&author, "author", "YOUR NAME", "Author name for copyright attribution")
    viper.BindPFlag("author", rootCmd.PersistentFlags().Lookup("author"))
}
// 不想使用的话相关可以注释掉 viper 相关的，编译出来的程序能小几M
```

# root.go 文件简单示例

`rootCmd` 的声明通常会被封装在一个函数中，这个封装函数会被 Execute() 执行。

```go
/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"os"

	vipercmd "github.com/DesistDaydream/go-cobra/cmd/viper"
	"github.com/DesistDaydream/go-cobra/config"
	"github.com/spf13/cobra"
)

type RootFlags struct {
	// 这里定义的变量，可以在下面的 init 函数中，通过 rootCmd.PersistentFlags().StringVar(&CfgFile, "config", "", "指定配置文件(默认在$HOME/.cobracli.yaml)") 进行绑定
	// 也可以通过 viper 进行绑定
	CfgFile string
}

var rootFlags RootFlags

// Execute 将所有子命令添加到根命令并设置 Flags。这由 main.main() 调用。它只需要对 rootCmd 发生一次。
func Execute() {
	app := newApp()
	err := app.Execute()
	if err != nil {
		os.Exit(1)
	}

}

func newApp() *cobra.Command {
	// rootCmd 表示在没有任何子命令调用的情况时的基本命令。
	var rootCmd = &cobra.Command{
		Use:   "go-cobra",
		Short: "这个应用简要的描述",
		Long: `横跨多行的较长描述，可能包含示例和使用应用程序的用法。 例如：
当我运行程序时，会显示该描述内容
	如果使用缩进，这行在界面展示时有缩进。`,
		// 如果这个应用没有任何子命令，直接使用 go-cobra 执行的话，将会执行下面 Run 字段指定的函数
		Run: rootRun,
	}

	// 我们可以在这里定义命令行 Flags 和 配置设置。
	// 这里可以做一些初始化的工作，比如初始化数据库连接、初始化日志、读取配置文件等等

	// ######## 添加 命令行Flags ########
	// Cobra 支持 持久性flags (i.e. Global Flags)，如果在这个位置定义，则这些 flags 对应用程序来说是全局的。
	// 第一个参数是变量，用于存储该flag的值；第二个参数为该flag的名字；第三个参数为该flag的默认值,无默认值可以为空；第四个参数是该flag的描述信息
	// 比如我现在使用如下命令: go-cobra --config abc 。那么 cfgFile 的值为abc。
	rootCmd.PersistentFlags().StringVarP(&rootFlags.CfgFile, "config", "c", "", "指定配置文件")
	// Cobra 还支持本地 flags ，仅在直接调用此命令时才有意义。
	rootCmd.Flags().BoolP("toggle", "t", false, "关于toggle标志的帮助信息")

	// ######## 添加 配置 ########
	// ！！！注意！！！：Cobra 只有在上面的 Run 字段定义的函数运行之前才会解析手动指定的命令行 Flags，否则只能获取到代码中设置的 Flags 默认值。
	// 比如运行 go run main.go --config="abc.yaml" 时，rootFlags.CfgFile 并不会被赋值为 abc.yaml，而是默认值。
	// 此时有两种方式解决这个问题：
	// 1. 使用 Prase() 函数，提前解析 Flags：
	// rootCmd.PersistentFlags().Parse(os.Args)
	// 2. 使用 OnInitialize() 函数，该函数会在 Command.Run 字段指定的函数执行前，先执行 initConfig 函数。
	// 查看 Cobra 源码，OnInitialize() 中的 initializers 变量会在 preRun() 函数中被执行。
	cobra.OnInitialize(initConfig)
	// 假如我现在在这里执加了一行 config.NewConfig(rootFlags.CfgFile)，那么这个函数其实是会在 OnInitialize 函数执行之前执行的。
	// config.NewConfig(rootFlags.CfgFile)

	// ######## 添加 子命令 ########
	// 为了更好的管理子命令，我们通常会将子命令放在不同的文件中，然后在这里进行注册
	rootCmd.AddCommand(
		NewVersionCmd(),
		vipercmd.NewViperCmd(),
	)

	return rootCmd
}

func initConfig() {
	// 使用 Viper 简化处理配置文件的过程。Viper 可以从 JSON、TOML、YAML、HCL、环境变量和命令行参数等等地方中读取配置。
	config.NewConfig(rootFlags.CfgFile)
}

func rootRun(cmd *cobra.Command, args []string) {
	fmt.Println("主程序运行后执行的代码块。如果注销 Run，则运行主程序会显示上面Long上的信息")
	fmt.Println("在 Run 字段指定的函数中，我们可以获取到 Flags 的值：", rootFlags.CfgFile)
}
```

# cobra.Command 结构体解析（待整理）

## 别名(Aliases)

现在我们想添加一个别名

```
cli
|----app
|----remove|rm
```

我们修改下初始化值即可

```go
var removeCmd = &cobra.Command{
	Use:   "remove",
    Aliases: []string{"rm"},
```

## 命令帮助添加示例(Example)

我们修改下 remove 的 Run 为下面

```go
Run: func(cmd *cobra.Command, args []string) {
           if len(args) == 0 {
              cmd.Help()
              return
           }
},
```

运行输出里 example 是空的

```go
[root@k8s-m1 cli]# go run main.go app remove
A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.

Usage:
  cli app remove [flags]

Aliases:
  remove, rm

Flags:
  -h, --help          help for remove
  -n, --name string   The application to be executed

Global Flags:
      --config string   config file (default is $HOME/.cli.yaml)
```

添加 example

```go
var removeCmd = &cobra.Command{
	Use:   "remove",
        Aliases: []string{"rm"},
        Example: `
cli remove -n test
cli remove --name test
`,
```

```go
go run main.go app remove
A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.

Usage:
  cli app remove [flags]

Aliases:
  remove, rm

Examples:

cli remove -n test
cli remove --name test


Flags:
  -h, --help          help for remove
  -n, --name string   The application to be executed

Global Flags:
      --config string   config file (default is $HOME/.cli.yaml)
```

## 参数验证器(Args)

该字段接收类型为`type PositionalArgs func(cmd *Command, args []string) error`
内置的为下面几个:

- `NoArgs`: 如果存在任何位置参数，该命令将报告错误。
- `ArbitraryArgs`: 该命令将接受任何 args。
- `OnlyValidArgs`: 如果存在任何不在 ValidArgs 字段中的位置参数，该命令将报告错误 Command。
- `MinimumNArgs(int)`: 如果没有至少 N 个位置参数，该命令将报告错误。
- `MaximumNArgs(int)`: 如果有多于 N 个位置参数，该命令将报告错误。
- `ExactArgs(int)`: 如果没有确切的 N 位置参数，该命令将报告错误。
- `RangeArgs(min, max):` 如果 args 的数量不在预期 args 的最小和最大数量之间，则该命令将报告错误。
- 自己写的话传入符合类型定义的函数即可

```go
  Args: func(cmd *cobra.Command, args []string) error {
  if len(args) < 1 {
    return errors.New("requires at least one arg")
  }
  if myapp.IsValidColor(args[0]) {
    return nil
  }
  return fmt.Errorf("invalid color specified: %s", args[0])
},
```

前面说的没传递选项和任何值希望打印命令帮助也可以用`MinimumNArgs(1)`来触发



# 自定义 help,usage 输出

help

```go
command.SetHelpCommand(cmd *Command)
command.SetHelpFunc(f func(*Command, []string))
command.SetHelpTemplate(s string)
```

usage

```go
command.SetUsageFunc(f func(*Command) error)
command.SetUsageTemplate(s string)
```

# Run 的 hook

Run 功能的执行先后顺序如下：
- PersistentPreRun
- PreRun
- Run
- PostRun
- PersistentPostRun

接收 `func(cmd *Command, args []string)` 类型的函数，Persistent 的能被下面的子命令继承

RunE 功能的执行先后顺序如下：
- PersistentPreRunE
- PreRunE
- RunE
- PostRunE
- PersistentPostRunE

接收 `func(cmd *Command, args []string) error` 的函数

## 预处理相关函数说明

当具有多级子命令时，`PersistentXXX()` 相关函数只会执行一次

比如现在创建了一个 cobra 命令，具有如下几个子命令：
- add
  - command
  - args
- del

如果在 cobra 和 add 中都使用了 PersistentPreRun() 函数的话，只会有第一个执行，并且是子命令的方法优先，参考 Issue：
- https://github.com/spf13/cobra/issues/216
- https://github.com/spf13/cobra/issues/252

可以在最底层的子命令中，通过如下方式执行父命令的 `PersistenXXX()` 函数

```go
func CreateCommand() *cobra.Command {
	subCmd := &cobra.Command{
		PersistentPreRun: subPersistentPreRun,
	}

	subCmd.AddCommand(
		CreateSubSubCommand(),
	)

	return subCmd
}

func subPersistentPreRun(cmd *cobra.Command, args []string) {
	// 执行父命令的预运行逻辑
	parent := cmd.Parent()
	if parent.PersistentPreRun != nil {
		parent.PersistentPreRun(parent, args)
	}

	// 本子命令的预运行逻辑
}
```

### OnInitialize与OnFinalize函数

除了 PersistentXXX() 这种函数以外，我们还可以使用 `OnInitialize()` 函数来执行预运行的逻辑，该函数可以避免在每一个子命令的  `PersistentPreRun()` 中重复调用 `parent.PersistentPreRun`。
- `OnInitialize()` 函数会在**调用**每个命令的 `Execute()` 方法时运行。
- `OnFinalize()` 函数则会在**完成**每个命令的 `Execute()` 方法时运行。

`OnInitialize()` 会将其参数赋值给 initializers 变量(这个变量的类型是一个函数)，该变量会在 Command.preRun() 函数中被执行
> 注意，这个 Command.preRun() 与 Command.PreRun() 不同。前者是 Command 结构体的方法，后者是 Command 的一个属性。
> 在 execute() 中，preRun 在 PreRun 之前执行。

