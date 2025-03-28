---
title: Bash
linkTitle: Bash
weight: 1
---

# 概述

> 参考：
>
> - [GNU 官网](https://www.gnu.org/software/bash/)
> - [Wiki, Bash](https://en.wikipedia.org/wiki/Bash_(Unix_shell))
> - [网道，Bash 脚本教程](https://wangdoc.com/bash/index.html)

**Bourne Again Shell(简称 Bash)** 是 Brian Fox 为 GNU 项目编写的 Unix Shell 和编程语言，作为 **Bourne shell(简称 sh)** 的免费软件替代品，于 1989 年首次发布，已被用作绝大多数 Linux 发行版的默认登录 Shell。Bash 是 Linus Torvalds 在移植 GCC 到 Linux 时最先移植的程序之一。

Bash 是一种 Shell，学习 Bash，其实也算是学习一种脚本式的编程语言，Bash 本身就是一种类似编译器似的存在。

# Bash 关联文件与配置

## 全局配置文件，对所有用户生效的配置

**/etc/environment** # 系统的环境变量，所有登录方式都会加载的文件。

**/etc/profile** # 任何用户使用 shell 时都会加载的配置。linux 不推荐直接修改该文件。加载该配置时自动加载 `/etc/profile.d/*.sh` 的所有文件

**/etc/bashrc** # 常用于设置登录功能和命令别名。linux 不推荐直接修改该文件。加载该配置时自动加载 `/etc/profile.d/*.sh` 的所有文件

**/etc/profile.d/\*.sh** # 类似于 include 的效果。通常用来创建自定义配置。

在 **/etc/profile** 和 **/etc/bashrc** 中都会有如下代码块

```bash
for i in /etc/profile.d/*.sh /etc/profile.d/sh.local ; do
    if [ -r "$i" ]; then
        if [ "${-#*i}" != "$-" ]; then
            . "$i"
        else
            . "$i" >/dev/null
        fi
    fi
done
```

这段脚本的大致意思就是：遍历 /etc/profile.d 目录下所有以 .sh 结尾的文件和 sh.local 文件。判断它们是否可读（`[ -r "$i"]`），如果可读，判断当前 Shell启动方式是不是交互式（`$-` 中包含 i）的，如果是交互式的，在当前 Shell 进程中执行该脚本（`. "$i"`，`source "$i"` 的简写， Shell 的模块化方式），否则，也在当前 Shell 进程中执行该脚本，只不过将输出重定向到了 /dev/null 中。

> `${-#*i}` 这个表达式的意思是：从左向右，在 - 变量中找到第一个 i ，并截取 i 之后的子串。

## 用户配置文件，对部分用户生效的配置

> 这些配置文件一般都定义在用户的家目录当中，所以当某一用户使用 shell 时，就会在其家目录中加载这些配置文件。

- **~/.bash_profile** # 仅对当前用户有效。该配置文件会对 ~/.bashrc 进行判断，如果存在，则加载 ~/.bashrc。
- **~/.bash_login** # 仅对当前用户有效。该配置文件会对 ~/.bashrc 进行判断，如果存在，则加载 ~/.bashrc。
- **~/.profile** # 仅对当前用户有效。该配置文件会对 ~/.bashrc 进行判断，如果存在，则加载 ~/.bashrc。
- 其他
  - **~/.bashrc** # 仅对当前用户有效。该配置文件会对 /etc/bashrc 进行判断，如果存在，则加载 /etc/bashrc。

Note：

- 仅在登录的第一次读取一次 bash 这些配置文件，如果在里面加了内容想让其实现，需要重新登录或者用 source FILE 命令来让更改的内容生效）
- 用户登录时逐一加载 ~/.bash_profile、~/.bash_login、~/.profile。当任何一个文件存在时，都不再加载其余文件。
- 为什么配置文件会有这么多分类呢？详见 《shell 的四种登录与交互模式》章节。不同模式，加载的配置文件是不同的
- 这三个文件通常只会存在一个，并且在其中的代码中，包含了判断是否存在 ~/.bashrc 文件并执行的逻辑。

# 登录与交互模式

这里面的概念推荐有一定 Linux 基础了之，尤其是得真正明白 Shell 到底是什么之后再来看。

首先，有几种对登录类型的描述：

- 交互式：一个个地输入命令并及时查看它们的输出结果，整个过程都在跟 Shell 不停地互动。
- 非交互式：运行一个 Shell 脚本 文件，让所有命令批量化、一次性地执行。
- 登录式：需要输入用户名和密码才能使用。
- 非登录式：直接可以使用。

这几种类型的不通组合，决定了运行 Shell 的模式

当我们运行一个 Shell 之后，Shell 会选择下面 4 中模式之一，作为运行本次 Shell 的模式，**不同的模式，加载的配置文件是不同的**。

- **login + interactive # 登录交互。**
  - 首先读取并执行 /etc/profile。
  - 然后逐一加载 ~/.bash_profile、~/.bash_login、~/.profile。当任何一个文件存在时，都不再加载其余文件。
- **login + non-interactive # 登录不交互。**
  - 与 登录交互 模式相同。
- **non-login + interactive # 不登陆交互**
  - 直接加载 ~/.bashrc 文件
- **non-login + non-interactive # 不登陆不交互**
  - 与 不登录交互 模式相同

## 如何判断是否为交互式 Shell? 有两种方式

查看特殊变量 `-` ，如果输出的值包含 `i`，则是交互式，否则是非交互式

```bash
echo $-
# 比如下面的情况
root@desistdaydream:~# ssh 172.19.42.248
root@desistdaydream:~# echo $-
himBHs

# 当使用 ssh 登录时，使用 -T 参数不分配终端，则 $- 没有 i
root@desistdaydream:~# ssh -T 172.19.42.248
echo $-
hBs

```

查看变量 PS1 是否为空，如果不为空，则是交互式，否则为非交互式

```bash
echo $PS1
```

Note：这里需要对“交互式”这三个字进行一下说明。与平时理解的所谓交互式不太一样，这里面虽然人类还是可以与设备交互，但是依然称为“非交互式”。

如何判断是否为登录式 Shell ?

执行命令 shopt login_shell，如果 login_shell 的值为 on 表示登录式，为 off 表示非登录式。

```bash
~]# shopt login_shell
login_shell     on
```

## 典型登录模式总结

- 登陆机器后的第一个 shell：登录 + 交互
- 通过 ssh 登陆到远程主机：登录 + 交互
- 新启动一个 shell 进程，如运行 bash ：不登陆 + 交互
- 执行脚本，如 bash script.sh ：不登陆 + 不交互
- 运行头部有如 #!/usr/bin/env bash 的可执行文件，如 ./executable ：不登陆 + 不交互
- 远程执行脚本，如 ssh user@remote script.sh ：不登陆 + 不交互
- 远程执行脚本，同时请求控制台，如 ssh user@remote -t 'echo $PWD' ：不登陆 + 交互
- 在图形化界面中打开 terminal：不登陆 + 交互

## 登录系统后可自动执行的配置文件

/etc/rc.local

- 注意：centos7 的 rc.local 没有可执行权限，需要添加权限(chmod +x /etc/rc.d/rc.local)，否则无法使用
- 官方推荐使用 [Systemd](/docs/1.操作系统/Systemd/Systemd.md) 来管理启动脚本。而不是这种方式

Systemd # [Systemd](/docs/1.操作系统/Systemd/Systemd.md)

# 在 Bash 中执行命令

## `.` 语法

作用:

- 打开一个子 shell 来读取并执行 FileName 中命令。
- 注：运行一个 shell 脚本时会启动另一个命令解释器.
- 每个 shell 脚本有效地运行在父 shell(parent shell)的一个子进程里. 这个父 shell 是指在一个控制终端或在一个 xterm 窗口中给你命令指示符的进程.shell 脚本也可以启动他自已的子进程. 这些子 shell(即子进程)使脚本并行地，有效率地地同时运行脚本内的多个子任务.

**Syntax(语法)**

**. FileName**

## source 程序

source 程序可以在当前 Shell 环境中从指定文件中读取并执行命令

注意：该命令可以使当前环境的配置文件在此环境中立刻生效而不用重启机器

作用：

- 这个命令其实只是简单地读取脚本里面的语句依次在当前 shell 里面执行，没有建立新的子 shell。那么脚本里面所有新建、改变变量的语句都会保存在当前 shell 里面
- source 命令(从 C Shell 而来)是 bash shell 的内置命令。
- source(或点)命令通常用于重新执行刚修改的初始化文档。

应用实例：一般用于写在 shell 脚本中，通过 source 执行外部文件中是变量赋值，这样不建立子 shell 的话，外部文件中的变量就可以在该脚本启动的 shell 中引用

**source \[选项] FILENAME** # 使环境变量立刻生效

在当前 bash 环境下读取并执行 FileName 中的命令。该 filename 文件可以“无执行权限”

EXAMPLE：

- source openrc admin admin # 运行 openrc 这个配置文件，把第一个参数 admin 和第二个参数 admin 送给 openrc 文件中的$1 和$2
- source .bash_profile #

## sh 和 bash 程序

实际上，执行 sh 或者 bash 就是相当于打开新的子 shell，并在新的 shell 中执行后续的命令。只不过 bash 与 sh 是不同的 shell，内置的功能有一些细微的区别。

作用:

- 开启一个新的 shell，读取并执行 File 中的命令。该 file 可以“无执行权限”
- 注：两者在执行文件时的不同，是分别用自己的 shell 来跑文件。
  - sh 使用“-n”选项进行 shell 脚本的语法检查，使用“-x”选项实现 shell 脚本逐条语句的跟踪
  - 可以巧妙地利用 shell 的内置变量增强“-x”选项的输出信息等。

**bash \[OPTIONS] \[File]** #

OPTIONS

- -n # 对指定文件进行语法检查
- -x # 打印出执行脚本的过程

# 命令行补全

> 参考：
>
> - [Wiki, CommandLineCompletion](https://en.wikipedia.org/wiki/Command-line_completion)

**Command Line Completion(命令行补全)** 也称为 **tab completion**，是命令行解释器的常见功能，在命令行中的程序，可以自动填充部分需要手动输入的命令。

由 bash-completion 程序实现

# Bash 关联文件与配置

**/etc/bash_completion.d/** #

**/usr/share/bash-completion/completions/** # 各种程序补全功能所需文件的保存目录。
