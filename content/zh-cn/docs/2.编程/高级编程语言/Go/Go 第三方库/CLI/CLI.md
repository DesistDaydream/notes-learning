---
title: "CLI"
linkTitle: "CLI"
weight: 20
---

# 概述

> 参考：

https://github.com/spf13/cobra

https://github.com/urfave/cli

github.com/alecthomas/kingpin

- Prometheus 使用这个 CLI

https://github.com/alecthomas/kong

- kingping 的作者改用 kong

在 Linux 中使用命令时会用到 `SubCommand` 和 `Options或Flag`，这些子命令和选项，以及命令的帮助信息都是通过`命令行参数处理`这个功能里的各种函数来实现的。该功能也也叫`从命令行读取参数`。并且这些`子命令以及选项或标志`统称为`命令行参数`。备注：Linux 中的每一个命令其实就是一个已经编译好的程序。

**flag 包**

使用`flag`包中的相关函数来实现解析命令行的 flag 或 option，详情见：https://golang.org/pkg/flag/#hdr-Usage

下面是其中几种 flag 包的格式说明

格式 1：`flag.TYPE(FlagName, DefaultValue, HelpInfo)`。FlagName 为参数名，DefaultValue 为参数的默认值，HelpInfo 为该参数的帮助信息。返回默认值的指针

格式 2：`flag.TYPE(Pointername ,FlagName, DefaultValue, HelpInfo)`。与上面的格式相同，只不过没有返回值，并且会把`DefaultValue`赋值给`Pointer`指针指向的变量，该变量需要提前定义。

```
test := flag.String("test","testValue","请指定test参数")
flag.Parse() //注意必须要有该行才能让test变量获取用户输入的参数，否则一直是默认值
fmt.Println(test)
```

使用方式：`go run test.go -test 123`结果为`123`；若不指定`-test 123`这个参数，则结果为`testValue`。如果使用`go tun test.go -help`则可获得帮助信息`请指定test参数`

**Args**

使用 Go 自带的 Args 切片变量获取命令参数。Args 切片的第一个位置为文件名的绝对路径，第二个位置是使用程序时输入的参数，以空格作为分隔符，分隔每个参数。每个参数都会保存到切片中。e.g.`go run runCommand.go ip`,这时候`ip`的值就会传递到 `Args[1]` 的位置上。
