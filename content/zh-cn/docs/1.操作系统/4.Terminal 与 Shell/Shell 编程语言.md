---
title: Shell 编程语言
weight: 1
---

# 概述

> 参考：
> 
> - [Wiki，Shell_script](https://en.wikipedia.org/wiki/Shell_script)

**Shell Programming Language(shell 编程语言)** 不是编译型的语言，而是一种脚本语言，通常称为 **Shell script**。

在类 Unix 系统中，我们最常使用的就是 Bash，通常 Shell 编程语言狭义上直接指的是 Bash 编程语言。

而在 Microsoft 系统中，早期是一种 **Batch file(批处理)** 文件，然后发展出 CMD，到近代可以使用 [PowerShell](docs/1.操作系统/4.Terminal%20与%20Shell/WindowsShell/PowerShell/PowerShell.md) 脚本语言。

# Shell 的基本结构与要素

Shell 语言的运行环境依赖于其所使用的 shell，是 bash、sh 还是 zsh 等等。想要执行 shell 语言，并不需要下载一个编译器，直接在指定的 shell 中执行代码即可。

脚本式的语言是编写完代码之后，一条一条执行，所以可以把平时在 Linux 上操作的 Bash Shell 想象成一个大型的文本编辑器，每输入一条命令，就相当于一行代码，直接通过这个 Bash 的 shell 就执行了，而把很多命令组合起来，放在一个文件里，直接运行该文件，与在界面输入很多内容，有异曲同工之妙。

由于 Shell 语言不需要编译器，所以 Shell 代码的第一行，必须指定其内的代码使用什么 Shell 来运行。

## Hello World

```bash
#!/bin/bash # 告诉内该脚本用什么shell运行，必须是脚本第一行
printf 'Hello World!'
```

其实如果是在某个系统下运行代码，第一行也是可以省略的，第一行的意思其实就是代表运行后续命令的环境，而第一行其实也是调用系统 /bin/ 目录下的 bash 二进制文件，来执行后续的代码。

## Shell 语言中的关键字

Shell 语言关键字主要看是哪种 Shell，在 Linux 下通常都是 Bash。

## Shell 语言中的标准库

Shell 语言中的标准库，就是该 Shell 所在系统下的可用命令+该 Shell 自带的一些命令。所以在 Shell 的 代码中，书写的大部分代码都是一个个的 Linux 命令。而那些非系统自带的命令(或者说工具)，可以理解为 Shell 编程语言的第三方库。

