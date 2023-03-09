---
title: Bash 变量
---

# 概述

> 参考：
> - [Manual(手册),bash(1)-形参](https://www.man7.org/linux/man-pages/man1/bash.1.html#PARAMETERS)-Shell 变量

环境变量是对当前环境起作用的变量，在日常操作中，我们最常用的就是 bash 这个 shell。

# Bash 自带的变量

> 参考：
> - [Manual(手册),bash(1)，形参-](https://www.man7.org/linux/man-pages/man1/bash.1.html#PARAMETERS)Shell 变量

**EDITOR=STRING** # 当 Bash 需要调用编辑器时，使用该变量指定的编辑器。
**IFS=STRING** # (Internal Field Separator)输入字段分隔符。`默认值：IFS 包含空格、制表符和回车`。
Bash 会根据 IFS 中定义的字符来进行字符串拆分。效果如下：

```bash
~]# map=(a,b c)
~]# echo ${map[0]}
a,b
~]# IFS=, && echo ${map[0]}
a b
```

**PATH=\<STRING>** # 命令的搜索路径。以 `:` 分隔的目录列表，bash 执行命令时将会从 $PATH 中查找用户输入的命令，以便执行这些命令，如果在 $PATH 中无法找到，则无法执行。`默认值：取决于操作系统`，通常都是 `/usr/local/bin:/usr/local/sbin:/usr/bin:/usr/sbin:/bin:/sbin`

- 在 bash 的源码(`[./config-top.h](https://git.savannah.gnu.org/cgit/bash.git/tree/config-top.h)`)中，我们可以看到 PATH 变量的默认值由 DEFAULT_PATH_VALUE 定义：

```c
/* The default value of the PATH variable. */
#ifndef DEFAULT_PATH_VALUE
#define DEFAULT_PATH_VALUE \
  "/usr/local/bin:/usr/local/sbin:/usr/bin:/usr/sbin:/bin:/sbin:."
#endif
```

- 我们也可以在系统中通过 `env -i bash -c 'echo "$PATH"'` 命令获取 bash 的 $PATH 变量的默认值

**TMOUT=INT** # Bash 在指定 INT 秒后未操作，则自动登出。`默认值：空`

- 如果设置为大于零的值，则 TMOUT 被视为读取内置的默认超时。当输入来自终端时，如果输入在指定的 X 秒后仍未到达，则 select 命令终止。在交互式 shell 中，该值被解释为发出主提示后等待一行输入的秒数。如果完整的输入行未到达，则 Bash 在等待该秒数后终止。

## 历史记录相关变量

**HISTTIMEFORMAT=STRING** # 历史记录的格式
**HISTSIZE=INT** # 历史记录可以保留的最大命令数
**HISTFILESIZE=INT** # 历史记录可以保留的最大行数
**HISTCONTROL=STRING** #

## 提示符相关变量

### PS1——默认提示符

如下所示，可以通过修改 Linux 下的默认提示符，使其更加实用。在下面的例子中，默认的 PS1 的值是“\s-\v$”,显示出了 shell 的名称的版本。我们通过修改，可以使其显示用户名、主机名和当前工作目录。

```bash
-bash-3.2$ export PS1="\u@\h \w> "
ramesh@dev-db ~> cd /etc/mail
ramesh@dev-db /etc/mail>
```

\[注: 提示符修改为 "username@hostname current-dir>的形式]

本例中 PS1 使用的一些代码如下：
**\u** # 用户名
**\h** # 主机名 建议在\h 值之后有一个空格。从个人角度来讲，使用这个空格可以增加一定的可读性。
**\w** # 当前目录的完整路径。请注意当你在主目录下的时候，如上面所示只会显示～

EXAMPLE

- **export PS1="\[\[\e\[34;1m]\u@\[\e\[0m]\[\e\[32;1m]\H\[\e\[0m] \[\e\[31;1m]\w\[\e\[0m]]\\$ "** # 好看的提示符样式
  - ![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/ic2gz0/1628044442278-702aab5e-f50e-40f9-924a-e00528d1dbff.png)

### PS2——再谈提示符

一个非常长的命令可以通过在末尾加“\”使其分行显示。多行命令的默认提示符是“>”。 我们可以通过修改 PS2 ，将提示符修改为"ABC" 。

    [root@lichenhao ~]# ls \
    > ^C
    [root@lichenhao ~]# PS2="ABC"
    [root@lichenhao ~]# ls \
    ABC^C

当用“\”使长命令分行显示，我发现非常易读。当然我也见过有的人不喜欢分行显示命令

### PS3 # Shell 脚本中使用 select 时的提示符

你可以像下面示范的那样，用环境变量 PS3 定制 shell 脚本的 select 提示：

不使用 PS3 的脚本输出:

cat ps3.sh

执行脚本如下==>

\[注: 缺省的提示符是 #?]

使用 PS3 的脚本输出：

加了一句提示脚本,更加的友好了

### PS4 # PS4-“set -x"用来修改跟踪输出的前缀

如果你像下面那样在调试模式下的脚本中，PS4 环境变量可以定制提示信息：

没有设置 PS4 时的 shell 脚本输出:

```bash
[root@localhost functions]# cat ps4.sh
set -x
echo "PS4 demo script"
ls -l /root/|wc -l
```

\[注: 当使用 sex -x 跟踪输出时的提示符为 ++]

设置 PS4 后的脚本输出:

PS4 在 ps.sh 中定义了下面两个变量

o $0 显示当前的脚本名

o $LINENO 显示的当前的行号

在 ps4.sh 脚本最顶层加一行这个==========>

输出的效果如下===>

如下图所示效果==>

\[注: 使用 PS4 后使 "{script-name}.{line-number}+" 成为 set –x 的命令提示符]
