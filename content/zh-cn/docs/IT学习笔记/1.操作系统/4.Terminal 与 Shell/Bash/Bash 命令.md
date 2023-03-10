---
title: Bash 命令
---

# 概述
> 参考：
> - [Manual(手册)，bash(1)-Shell 内置命令](https://www.man7.org/linux/man-pages/man1/bash.1.html#SHELL_BUILTIN_COMMANDS)

# source，sh，bash 等执行命令的工具
## source 程序

source 程序可以在当前 Shell 环境中从指定文件中读取并执行命令

注意：该命令可以使当前环境的配置文件在此环境中立刻生效而不用重启机器

作用：
1. 这个命令其实只是简单地读取脚本里面的语句依次在当前 shell 里面执行，没有建立新的子 shell。那么脚本里面所有新建、改变变量的语句都会保存在当前 shell 里面
2. source 命令(从 C Shell 而来)是 bash shell 的内置命令。
3. source(或点)命令通常用于重新执行刚修改的初始化文档。

应用实例：一般用于写在 shell 脚本中，通过 source 执行外部文件中是变量赋值，这样不建立子 shell 的话，外部文件中的变量就可以在该脚本启动的 shell 中引用

**source \[选项] FILENAME** # 使环境变量立刻生效

在当前 bash 环境下读取并执行 FileName 中的命令。该 filename 文件可以“无执行权限”

1. EXAMPLE：
   1. source openrc admin admin # 运行 openrc 这个配置文件，把第一个参数 admin 和第二个参数 admin 送给 openrc 文件中的$1 和$2
   2. source .bash_profile #

sh 和 bash 程序：

实际上，执行 sh 或者 bash 就是相当于打开新的子 shell，并在新的 shell 中执行后续的命令。只不过 bash 与 sh 是不同的 shell，内置的功能有一些细微的区别。

作用:

1. 开启一个新的 shell，读取并执行 File 中的命令。该 file 可以“无执行权限”
2. 注：两者在执行文件时的不同，是分别用自己的 shell 来跑文件。
   1. sh 使用“-n”选项进行 shell 脚本的语法检查，使用“-x”选项实现 shell 脚本逐条语句的跟踪
   2. 可以巧妙地利用 shell 的内置变量增强“-x”选项的输出信息等。

**bash \[OPTIONS] \[File]** #

OPTIONS
1. -n # 对指定文件进行语法检查
2. -x # 打印出执行脚本的过程

## ./ 的命令用法

作用:
1. 打开一个子 shell 来读取并执行 FileName 中命令。
2. 注：运行一个 shell 脚本时会启动另一个命令解释器.
3. 每个 shell 脚本有效地运行在父 shell(parent shell)的一个子进程里. 这个父 shell 是指在一个控制终端或在一个 xterm 窗口中给你命令指示符的进程.shell 脚本也可以启动他自已的子进程. 这些子 shell(即子进程)使脚本并行地，有效率地地同时运行脚本内的多个子任务.

语法格式：./FileName

# 命令行补全

> 参考：
> - [Wiki,CommandLineCompletion](https://en.wikipedia.org/wiki/Command-line_completion)

**Command Line Completion(命令行补全)** 也称为 **tab completion**，是命令行解释器的常见功能，在命令行中的程序，可以自动填充部分需要手动输入的命令。

## 关联文件与配置

### bash-completion

**/usr/share/bash-completion/completions/\* **# 各种程序补全功能所需文件的保存目录。
