---
title: exec,xargs,管道符等组合命令
---

# `|`(管道符)：把|前面的标准输出内容当作|后面的标准输入内容

1. EXAMPLE

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/sl7f9m/1616165959227-38b8c228-113d-4a87-a335-22560fdc396b.jpeg)

- echo “--help” | cat #把--help 当作 cat 的标准输入输出到屏幕上，如图所示，注意与 xargs 应用实例 a 的区别

# exec

exec \[-cl] \[-a name] \[COMMAND \[ARGUMENTS...]] #如果指定了 command，它将替换 shell。 没有创建新进程。 参数成为命令的参数。 如果提供了-l 选项，则 shell 会在传递给 command 的第 0 个参数的开头放置一个破折号。 这是 login（1）的作用。 -c 选项导致命令在空环境中执行。 如果提供了-a，则 shell 将 name 作为第 0 个参数传递给执行的命令。 如果由于某种原因无法执行命令，则会退出非交互式 shell，除非启用了 shell 选项 execfail，在这种情况下它将返回失败。 如果无法执行文件，则交互式 shell 将返回失败。 如果未指定 command，则任何重定向在当前 shell 中生效，返回状态为 0.如果存在重定向错误，则返回状态为 1

通过 exec 来执行的命令会顶替掉当前 shell 的进程，但是进程 PID 保持不变

# xargs # 把从|前面获得的输入内容，当成 xargs 后面的命令的选项或者参数来执行

xargs \[OPTIONS] COMMAND #不指定 COMMAND 则默认输出到屏幕上

OPTIONS：

1. -d 指定获得输入内容的分隔符，默认分隔符为空白或换行

2. -n 每次传递 NUM 个字符给|后的内容

EXAMPLE：

- echo “--help” | xargs cat #把--help 当作 cat 的选项或者参数，如图所示，注意与管道符的应用实例 a 中的区别

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/sl7f9m/1616165959237-c4a1a831-c6a0-4c52-ba6f-9f071cb59ba3.jpeg)

- ls /var/run/docker/netns | xargs -I {} nsenter --net=/var/run/docker/netns/{} ip addr #遍历 netns 目录下的所有文件,通过 xargs 命令把所有文件名传递给后面命令，作为后面命令的参数
