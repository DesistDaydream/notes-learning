---
title: Bash 变量
---

# 概述

> 参考：
> 
> - [Manual(手册),bash(1)-形参](https://www.man7.org/linux/man-pages/man1/bash.1.html#PARAMETERS)-Shell 变量

Bash 可以从逻辑上分为如下几种变量

- **局部变量** # 局部变量在脚本或命令中定义，仅在当前shell实例中有效，其他shell启动的程序不能访问局部变量。
- **环境变量** # 继承自操作系统的环境变量。所有的程序，包括 Shell 启动的程序，都能访问环境变量，有些程序需要环境变量来保证其正常运行。必要的时候 Shell 脚本也可以定义环境变量。

局部变量与环境变量的区别主要在于是否可以被子进程继承

```bash
~]# TestVar="This is Normal Var"
~]# export TestEnvVar="This is Enviroment Var"

~]# bash # 在这里进入了一个新的 Bash 程序，不通过 export 声明的非环境变量
~]# echo $TestVar

~]# echo $TestEnvVar
This is Enviroment Var
```

这个示例中可以看到，当我们启动一个新 Bash 时，不使用 `export` 声明的变量，将不会被子 Bash 继承。这就是环境变量与普通变量的基本区别。

## 局部变量

没什么特殊的说明

## 环境变量

根据变量的定义位置，环境变量分为多个作用域：

- 系统范围
- 用户范围
- 进程范围

# 声明变量与取消变量

对于 Shell 编程语言来说，一个变量其实不声明也是可以的，默认任意字符串组合成的变量的值都为空，比如

```bash
~]# env | grep randomVar
~]# echo $randomVar

 ~]# echo $?
0
```

可以看到，一个不存在的变量依然可以正常输出，只不过值为空。

## 局部变量

直接使用 `变量名=变量值` 的方式即可声明一个变量并为该变量赋值。

VarName=VALUE

**给变量赋予默认值**

`VarName=${VarName:=VALUE}` # 如果 VarName 不空，则其值不变；否则，VarName 会使用 VALUE 作为其默认值

## 环境变量

**export \[OPTIONS] \[VarName\[=VALUE] ...]** # 设置或显示环境变量(export 的效力仅作用于该次登陆操作)。

用户创建的变量仅可用于当前 Shell，子 Shell 默认读取不到父 Shell 定义的变量。为了把变量传递给子 Shell，需要使用 export 命令。这样输出的变量，对于子 Shell 来说就是环境变量。

OPTIONS

- **-f** # 代表\[NAME]中为函数名称
- **-n** # 删除指定的变量。变量实际上并未删除，只是不会输出到后续指令的执行环境中
- **-p** # 列出所有的 shell 赋予程序的环境变量。

EXAMPLE

- export VarName="Value" #
- export VarName #

**\[set] VarName="Value"**

EXAMLE

- test="test.test" # 设定一个名为 test 的变量的值为 test.test

**unset VarName**

EXAMPLE

- unset TestVar # 取消变量名为 TestVar 的值

注意：如果想要给 `$PATH` 变量增加内容，则需要用命令 `PATH=$PATH:/NAME/NAME`，如果前面不加 `$PATH`，那么这个变量就等于被改写成 `/NAME/NAME`，这点在修改变量时候尤为重要，必须要在定义 PATH 的引入本身已经定义好的 `$PATH`

**declare # 声明 shell 变量**

declare 为 shell 命令，在第一种语法中可用来声明变量并设置变量的属性，在第二种语法中可用来显示 shell 函数。若不加上任何参数，则会显示全部的 shell 变量与函数(与执行 set 指令的效果相同)。

语法格式：

declare \[+/-]\[OPTIONS] VarName

OPTIONS

- **-** # 给变量添加类型属性
- **+** # 取消变量的类型属性
- **-a** # 将变量声明为数组型
- **-i** # 将变量声明为整型
- **-x** # 将变量声明为环境变量
- **-r** # 将变量声明为只读变量
- **-p** # 查看变量的被声明的类型
- **-f** #

# 引用变量

Linux 中，普通变量与环境变量的引用方式相同。

引用变量有两种写法 `$VAR_NAME` 或者 `${VAR_NAME}` ，推荐使用第二种方法，两种方法的差别如下

```bash
~]# name='hello DesistDaydream!'; name="$nameceshi"; echo $name
# 这里输出为空，不是手打的回车
~]# name='hello DesistDaydream!'; name="${name}ceshi"; echo $name
hello DesistDaydream!ceshi
```

我们还可以把命令的执行结果赋值给变量,使用 `$(COMMAND)` 方式

```bash
~]# time=$(date)
~]# echo $time
2020年 08月 28日 星期五 12:44:46 CST
```

变量的拼接

```bash
~]# name="小明"
~]# echo $name
小明
~]# name=$name"和老王"
~]# echo $name
小明和老王
~]# name="${name}在一起"
~]# echo $name
小明和老王在一起
```

删除变量

```bash
~]# name='hello Desist Daydream!'
~]# echo ${name}
hello Desist Daydream!
~]# unset name
~]# echo ${name}
~]#
```

# Bash 自带的变量

> 参考：
> 
> - [Manual(手册),bash(1)，形参-](https://www.man7.org/linux/man-pages/man1/bash.1.html#PARAMETERS)Shell 变量

**EDITOR=STRING** # 当 Bash 需要调用编辑器时，使用该变量指定的编辑器。

**IFS=STRING** # (Internal Field Separator)输入字段分隔符。`默认值：IFS 包含空格、制表符和回车`。

- Bash 会根据 IFS 中定义的字符来进行字符串拆分。效果如下：

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

详见 [Bash 操作历史记录](/docs/1.操作系统/4.Terminal%20与%20Shell/Bash/Bash%20操作历史记录.md)

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

- **`\u`** # 用户名
- **`\h`** # 主机名 建议在\h 值之后有一个空格。从个人角度来讲，使用这个空格可以增加一定的可读性。
- **`\w`** # 当前目录的完整路径。请注意当你在主目录下的时候，如上面所示只会显示～

**EXAMPLE**

好看的提示符样式

`export PS1="[\[\e[34;1m\]\u@\[\e[0m\]\[\e[32;1m\]\H\[\e[0m\] \[\e[31;1m\]\w\[\e[0m\]]\\$ "`

> 这个配置可以直接修改 `~/.bashrc` 文件中的 `$PS1` 变量即可

![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/ic2gz0/1628044442278-702aab5e-f50e-40f9-924a-e00528d1dbff.png)

### PS2——再谈提示符

一个非常长的命令可以通过在末尾加“\”使其分行显示。多行命令的默认提示符是“>”。 我们可以通过修改 PS2 ，将提示符修改为"ABC" 。

```bash
~]# ls \
> ^C
~]# PS2="ABC"
~]# ls \
ABC^C
```

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
~]# cat ps4.sh
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

## 特殊变量

**位置变量**：`$0` 该脚本所在的绝对路径 `$1` 脚本的第一个参数 `$2`.....等等

- `$数字` # 是位置参数的用法。如果运行脚本的时候带参数，那么可以在脚本里通过 $1 获取第一个参数，$2 获取第二个参数......依此类推，一共可以直接获取 9 个参数（称为位置参数）。$0 用于获取脚本名称。相应地，如果 $+数字 用在函数里，那么表示获取函数的传入参数，$0 表示函数名。
  - `$0` # 该脚本所在的绝对路径
  - `$1` # 脚本的第一个参数
  - `$2` # 脚本的第二个参数
  - .....等等
  - 比如：
    - ./test.sh a b c # 运行该脚本时候，a 就是变量 $1(第一个参数)，b 就是变量 $2(第二个参数)，变量 $# 为 3,一共 3 个位置参数
- 其他位置参数：
  - `$#` # 位置变量的个数
  - `$*` # 引用所有的位置参数，引用后就是显示或者执行引用的字符串
  - `$@` # 引用所有的位置参数，引用后就是显示或者执行引用的字符串

**特殊变量**

- $? # 上一条命令执行的结果的返回值，成功为 0，失败为不为 0
- `$$` # 当前 shell 的 PID 号
- $! # Shell 最后运行的后台进程的 PID

# 操作变量的高级技巧

## 使用 `${ }` 符号的功能来处理变量中的字符串

在日常使用变量中，经常会需要对变量中的值进行操作，比如删除、替换、截取等等。通过`${ }`符号，就可以对指定变量中的值进行各种操作

想要处理字符串，通常有如下几种方式：

- 从指定位置截取字符串
   - `${VARIABLE:START:LENGTH}` # 从 VARIABLE 值的 左边 起第 START 个字符开始，向右截取 LENGTH 个字符
   - `${VARIABLE:0-START:LENGTH}` # 从 VARIABLE 值的 右边 起第 START 个字符开始，向右截取 LENGTH 个字符
   - Note：
      - :LENGTH 可省略。省略的话表示截取到变量值的末尾
      - 从左边开始计数时，起始数字是 0（这符合程序员思维）；从右边开始计数时，起始数字是 1（这符合常人思维）。计数方向不同，起始数字也不同。
      - 左数和右数的区别就是其实位置那个位置有没有 0- 这个标识符
      - 不管从哪边开始计数，截取方向都是从左到右。
- 从指定字符串处截取字符串
   - `${VARIABLE#*CHARS}` # 从左边开始到第一个 CHARS 为止的字符全部忽略，只留下右边的所有字符
   - `${VARIABLE##*CHARS}` # 从左边开始到最后一个 CHARS 为止的字符全部忽略，只留下右边的所有字符
   - `${VARIABLE%CHARS*}` # 从右边开始到第一个 CHARS 为止的字符全部忽略，只留下左边的所有字符
   - `${VARIABLE%%CHARS*}` # 从右边开始到最后一个 CHARS 为止的字符全部忽略，只留下左边的所有字符
   - Note:
      - 截取时，CHARS 不被包含在内。CHARS 可以是一个字符，也可以是一个字符串，当作一个整体看待，不要把 CHARS 拆开。当 CHARS 为字符时，则在计数开始时，表示从出现该字符串整体开始算。
      - `-` 符号只是一个通配符，可以省略。上面的语法中的 \* 就是表示 CHARS 左侧或者右侧的所有字符
      - 比如 `${VARIABLE%%CHARS*}` 其实就是删掉从右数最后一个 CHARS 右侧的所有字符
      - 这种截取方式无法指定字符串长度。
- 替换变量中匹配到的字符串
   - `${VARIABLE/OldChars/NewChars}` # 将 VARIABLE 值中匹配到第一个的 OldChars 替换为 NewChars
   - `${VARIABLE//OldChars/NewChars}` # 将 VARIABLE 值中匹配到所有的 OldChars 替换为 NewChars
- 获取变量值的长度
   - `${#VARIABLE}` # 变量名前加 # 符号

**EXAMPLE**

**从指定位置截取字符串**

`export split='www.desistdaydream.com'`

-  `echo ${split:4:14}` # 从第四位字符开始(包括第四位)，一共截取14个字符。由于第一个字符是0号位置，所以第四位字符，按照人类的理解应该是第五个字符。
  -  输出结果：desistdaydream
-  echo ${split:4} # 省略 length，截取到字符串末尾
  -  输出结果：desistdaydream.com
-  `echo ${split: 0-18: 14}` # 从右边数，b是第 13 个字符。
  -  输出结果：desistdaydream
-  `echo ${split: 0-18}` # 省略 length，直接截取到字符串末尾
  -  输出结果：desistdaydream.com

**从指定字符串处截取字符串**

`split="http://www.desistdaydream.com/index.html"`

- `echo ${split#*/}` # 也可以写为 `${split#*p:/}` 或 `${split#http:/}`，效果相同。注意带 * 和不带 * 的区别。
  -  输出结果：/www.desistdaydream.com/index.html
-  `echo ${split##*/}`
  -  输出结果：index.html
-  `echo ${split%/*}`
  -  输出结果：[http://www.desistdaydream.com](http://www.desistdaydream.com)
-  `echo ${split%%/*}`
  -  输出结果：http:


## 使用 eval 命令让变量的值作为另一个变量的变量名

```bash
root@lichenhao:~# varname=name
root@lichenhao:~# name=lichenhao
root@lichenhao:~# echo $$varname
209409varname
root@lichenhao:~# echo '$'$varname
$name
root@lichenhao:~# eval echo '$'$varname
lichenhao
```
