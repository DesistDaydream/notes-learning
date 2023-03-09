---
title: "source，sh，bash 等执行命令的工具"
linkTitle: "source，sh，bash 等执行命令的工具"
weight: 20
---

# 概述

> 参考：
> -



# source程序：在当前shell环境中从指定文件中读取并执行命令

注意：该命令可以使当前环境的配置文件在此环境中立刻生效而不用重启机器

作用：


1. 这个命令其实只是简单地读取脚本里面的语句依次在当前shell里面执行，没有建立新的子shell。那么脚本里面所有新建、改变变量的语句都会保存在当前shell里面

2. source命令(从 C Shell 而来)是bash shell的内置命令。

3. source(或点)命令通常用于重新执行刚修改的初始化文档。


应用实例：一般用于写在shell脚本中，通过source执行外部文件中是变量赋值，这样不建立子shell的话，外部文件中的变量就可以在该脚本启动的shell中引用

source [选项] FILENAME		# 使环境变量立刻生效

在当前bash环境下读取并执行FileName中的命令。该filename文件可以“无执行权限”

1. EXAMPLE：

   1. source openrc admin admin # 运行openrc这个配置文件，把第一个参数admin和第二个参数admin送给openrc文件中的$1和$2

   2. source .bash_profile #



# sh和bash程序：

实际上，执行sh或者bash就是相当于打开新的子shell，并在新的shell中执行后续的命令。只不过bash与sh是不同的shell，内置的功能有一些细微的区别。

作用:


1. 开启一个新的shell，读取并执行 File 中的命令。该 file 可以“无执行权限”

2. 注：两者在执行文件时的不同，是分别用自己的 shell 来跑文件。

   1. sh使用“-n”选项进行shell脚本的语法检查，使用“-x”选项实现shell脚本逐条语句的跟踪

   2. 可以巧妙地利用shell的内置变量增强“-x”选项的输出信息等。



**bash [OPTIONS] [File]**	#

OPTIONS


1. -n # 对指定文件进行语法检查

2. -x # 打印出执行脚本的过程



# ./ 的命令用法：

作用:

1. 打开一个子shell来读取并执行FileName中命令。

2. 注：运行一个shell脚本时会启动另一个命令解释器.

3. 每个shell脚本有效地运行在父shell(parent shell)的一个子进程里. 这个父shell是指在一个控制终端或在一个xterm窗口中给你命令指示符的进程.shell脚本也可以启动他自已的子进程. 这些子shell(即子进程)使脚本并行地，有效率地地同时运行脚本内的多个子任务.


语法格式：./FileName






