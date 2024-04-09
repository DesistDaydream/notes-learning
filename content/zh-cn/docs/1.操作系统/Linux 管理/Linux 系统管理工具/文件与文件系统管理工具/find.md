---
title: find
linkTitle: find
date: 2024-04-09T08:54
weight: 20
---

# 概述

> 参考：
> 
> - [Manual(手册)，find(1)](https://man7.org/linux/man-pages/man1/find.1.html)

find 工具可以在目录中搜索文件

# Syntax(语法)

**find \[OPTIONS] \[PATH...] \[EXPRESSION]**

- PATH # 路径名，不写上默认表示在当前路径下搜索
- OPTIONS 与 EXPRESSION # 详见下文，PATH 可以指定一个或多个想要搜索的路径。

如果没有给出 PATH，则使用当前目录；如果未给出任何表达式，则使用表达式 -print

## OPTIONS

-H，-L 和-P 选项控制 find 处理 [文件管理](docs/1.操作系统/Kernel/Filesystem/文件管理/文件管理.md) 章节中提到的 **符号链接** 的行为。这些参数之后的命令行参数将被视为要检查的文件或目录的名称，直到以 `-` 开头的第一个参数或参数 `(` 或 `!` 。符号链接的概念详见：

- **-P** # 不跟踪符号链接。该选项为 find 的默认行为。当 find 检查或打印文件信息时，该文件是符号链接，则所使用的信息应取自符号链接本身的属性。
- **-L** # 跟踪符号链接。与 -type l 同时使用时仅搜索失效的符号链接。当 find 检查或打印有关文件的信息时，所使用的信息应从链接指向的文件的属性中获取，而不是从链接本身获取（除非它是断开的符号链接，或者 find 无法检查与之相关的文件）链接点）。使用此选项意味着-noleaf。如果以后使用-P 选项，则-noleaf 仍然有效。如果-L 有效，并且 find 在搜索过程中发现到子目录的符号链接，则将搜索该符号链接所指向的子目录。
  - 当-L 选项生效时，-type 谓词将始终与符号链接指向的文件类型匹配，而不是与链接本身匹配（除非符号链接断开）。使用-L 会使-lname 和-ilname 谓词始终返回 false。
- **-H** # 除了处理命令行参数时，不要跟随符号链接。当 find 检查或打印有关文件的信息时，所使用的信息应取自符号链接本身的属性。唯一的例外情况是在命令行上指定的文件是符号链接并且可以解析该链接时，这种情况下，使用的信息取自链接指向的任何位置（即跟随该链接） ）。如果无法检查符号链接指向的文件，则有关链接本身的信息将用作备用。如果-H 有效，并且在命令行上指定的路径之一是指向目录的符号链接，则将检查该目录的内容（尽管-maxdepth 0 当然可以防止此情况）。
- **-D debugoptions** # Print diagnostic information; this can be helpful to diagnose problems with why find is not doing what you want. The list of debug options should be comma separated. Compatibility of the debug options is not guaranteed between releases of findutils. For a complete list of valid debug options, see the output of find -D help. Valid debug options include
- **-Olevel** # Enables query optimisation. The find program reorders tests to speed up execution while preserving the overall effect; that is, predicates with side effects are not reordered relative to each other. The optimisations performed at each optimisation level are as follows.

注意：

- 上述五个 OPTIONS 必须出现在第一个 PATH 之前。
- 下面 EXPRESSIONS(表达式)中的 OPTIONS，用来控制 find 的行为，与上述这些 OPTIONS 不同。并且表达式中的 OPTIONS 需要在第一个 PATH 之后指定。
- find 默认将在当前目录下查找子目录与文件。并且将查找到的子目录和文件全部进行显示

## EXPRESSIONS(表达式)

EXPRESSIONS 由以下几部分组成

- **Options(选项)** # 这会影响整体操作，而不是特定文件的处理，并且始终返回 true
- **Tests(测试)** # Tests 下涉及的表达式会返回 true 或 false，通常基于我们正在考虑的文件的某些属性。例如，-empty 用来测试仅在当前文件为空时才为真。
  - 说人话：<font color="#ff0000">Tests 部分涉及的表达式是核心</font>。这些表达式是各种判断条件以便让 find 程序可以过滤出想要的文件。这个就有点像 Shell 脚本中的 `[[ ]]` 符号
- **Actions(动作)** # 具有副作用并返回 true 或 false 值
- **Operators(运算符)** # 运算符可以将多个不同条件的表达式连接起来以表达更丰富的过滤条件。比如 -a 表示逻辑与；-o 表示逻辑或；etc.

如果表达式除 -prune 之外不包含其他任何动作，则对表达式为 true 的所有文件执行 -print。

### Options

- **-maxdepth \<INT>** # 设置最大目录层级
  - 例：find /home -maxdepth 1 -name \*.log # 查找/home 目录下 1 层的以.log 结尾的文件（即只查找当前目录，不查找当前目录下的子目录中的内容）
  - -ls 假设 find 指令的回传值为 Ture，就将文件或目录名称列出到标准输出
  - 例：find /home -type d -ls # 查找/home 目录下的文件夹，并列出这些文件夹的详细信息

### Tests

**-newerXY** # 

**-type** # 根据文件类型查找 （find PATH -type 类型参数）

- 参数类型包括：f 普通文件，l 软连接文件，d 目录文件等
- 例: find -type l # 查找当前目录下的软连接文件

**-name** # 根据文件名查找 （\* 任意多个字符）（? 1 个字符）（\[] 指范围值，外侧加引号）

- 例：find /var –name "\*.log" # 查找/var 目录下，所有文件名最后字符是.log 的文件

**-iname** # 不区分大小写根据文件名查找

**-inum** # 根据 inode 号查找

**-size** # 根据大小查找

- 例:find /etc -size -10k -size +20k # 查找/etc 目录下小于 10k 大于 20k 的文件
- 注意+-号，如果没有，是精确这么大，通常都会带上+或-号，表示一个范围

**-user** # 根据所有者查找 （user 可改成 group 则查找属组为 lisi 的文件）

- 例:find /home -user lisi # 查找/home 目录下，属主为 lisi 的文件

**-perm** # 根据文件权限查找

- 例:find /boot -perm 644 # 查找/boot 文件加下，权限为 644 的文件

**-TIME {+|-}NUM** # 根据时间查找

- - 表示该时间之前，-表示该时间之内
- TIME 分为两部分
  - 第一部分，# 表示 ctime，atime，mtime；time 天，min 分钟
    - **c** # change # 表示属性被修改过：所有者、所属组、权限
    - **a** # access # 被访问过(被查看过)
    - **m** # modify # 表示内容被修改过
  - 第二部分
    - **time** # 天
    - **min** # 分钟

### Actions

**\[ -exec 或者-ok command ] {} \[];**

- **-delete** # 删除找到的文件。使用-delete 会自动打开“ -depth”选项。
- **-exec \<COMMAND>** # 假设 find 指令的回传值为 True，就执行该指令；-ok 与-exec 类似，只不过是交互型。
   - 格式:find PATH 选项 选项内容 -exec COMMAND {} ;
   - 该格式的意思是使用 find 查找出来的内容放到{}中，再对{}中的内容逐条执行 COMMAND 命令
   - 它的终止是以;为结束标志的，所以这句命令后面的分号是不可缺少的，考虑到各个系统中分号会有不同的意义，所以前面加反斜杠。（反斜杠的意思参考正则表达式）
   - {} 花括号里面存放前面 find 查找出来的文件名。
   - 注意：固定格式，只能这样写。注意中间的空格。
- **-ok** #

注意该语法格式中不要少了最后的分好和最后大括号周围的空格，-ok 为会询问，交互式

### Operators

**-a, -and** # 逻辑运算符 AND(与)

**-o, -or** # 逻辑运算符 OR(或)

**-not, !** # 逻辑运算符 NOT(非)

# EXAMPLE

根据时间查找文件

- 查找今天修改过的文件
  - `find -type f -newermt $(date +"%Y-%m-%d")`
- 查找今天某个时间段（比 20:15 新且比 20:30 旧）内修改过的文件
  - `find -type f -newermt "today 20:15" -and ! -newermt "today 20:30"`
- 查找 2024 年 3 月 5 号之后修改过的文件
  - `find -type f -newermt "2024-03-05 9:00"`
- 查看当前目录下两天之前修改过得并且 10 分钟之内查看过得文件
  - find -mtime +2 -amin -10
- 查看当前目录下 7 天之内并且 2 天之前的文件
  - find -mtime -7 -a -mtime +2

有一些文件的硬链接数量很多，有相同的 i 节点，查找其中一个文件的 i 节点号，一次性删除。

- `find . -inum 2310630 -exec rm {} ;` 

- 删除 /var/log/pods 目录下失效的符号链接
  - find -L /var/log/pods -type l -delete
- 找到/home 目录下的 samlee 用户的所有文件并删除
  - find /home -user samlee -exec rm –r {} ;
  - 注：rm -r 连带目录一起删除。报错原因：-exec 不适合大量传输，速率慢导致。

- 在 etc 目录下查找大于 1k 并且小于 10k 的文件
  - find /etc -size +1k -a -size -10k
- find -maxdepth 1 -mtime +1 -exec ls -l {} ; # 查找出来的如果是文件夹，那么就对这个文件夹执行该命令，如下图所示，查找出来./test 等 3 个文件夹，那么就这三个文件夹执行 ls -l 的命令，即
  - ls -l ./test
  - ls -l ./lichenhao
  - ls -l ./lost+found
  - 注：查找出来的文字，一字不差的全部放在后面的 `{}` 中，等待 COMMAND 执行，所以没法列出目录的详细信息
  - 注意与-ls 参数的区别

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/idlpqz/1616166302372-29738c6b-f92a-43b5-9243-9aa5483629ca.jpeg)

