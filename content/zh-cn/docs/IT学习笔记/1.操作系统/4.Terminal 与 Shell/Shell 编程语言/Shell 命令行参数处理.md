---
title: Shell 命令行参数处理
---

# 概述

> 参考：
> - <https://www.cnblogs.com/klb561/p/9211222.html>
> - <https://blog.csdn.net/qq_22203741/article/details/77671379>

# shell 代码命令行选项与修传入参数处理

在编写 shell 程序时经常需要处理命令行参数，本文描述在 bash 下的命令行处理方式。

选项与参数：

如下命令行：

./test.sh -f config.conf -v --prefix=/home

1. -f 为选项，它需要一个参数，即 config.conf
2. -v 也是一个选项，但它不需要参数。
3. \--prefix 我们称之为一个长选项，即选项本身多于一个字符，它也需要一个参数，用等号连接，当然等号不是必须的，/home 可以直接写在--prefix 后面，即--prefix/home,更多的限制后面具体会讲到。

在 bash 中，可以用以下三种方式来处理命令行参数，每种方式都有自己的应用场景。

1. 通过位置变量手工处理参数 # 参考变量与环境变量 文章中的位置变量
2. getopts # 简单工具
3. getopt # 复杂工具

依次讨论这三种处理方式。

# 通过位置变量手工处理参数

在手工处理方式中，首先要知道几个变量，还是以上面的命令行为例：

代码如下:

- $0 ： ./test.sh,即命令本身，相当于 c/c++中的 argv\[0]

- $1 ： -f,第一个参数.

- $2 ： config.conf

- $3, $4 ... ：类推。

- $#  参数的个数，不包括命令本身，上例中$#为 4.

- $@ ：参数本身的列表，也不包括命令本身，如上例为 -f config.conf -v --prefix=/home

- $* ：和$@相同，但"$*" 和 "$@"(加引号)并不同，"$*"将所有的参数解释成一个字符串，而"$@"是一个参数数组。

例子：

```bash
#!/bin/bash
for arg in "$*"
do
   echo $arg
done
for arg in "$@"
do
 echo $arg
done
```

执行./test.sh -f config.conf -n 10 会打印：

    -f config.conf -n 10    #这是"$*"的输出
    -f   #以下为$@的输出
    config.conf
    -n
    10

所以，手工处理的方式即对这些变量的处理。因为手工处理高度依赖于你在命令行上所传参数的位置，所以一般都只用来处理较简单的参数。

例如：

./test.sh 10

而很少使用./test -n 10 这种带选项的方式。 典型用法为：

    #!/bin/bash
    if [ x$1 != x ]
    then
        #...有参数
    else
    then
        #...没有参数
    fi

为什么要使用 x$1 != x 这种方式来比较呢？想像一下这种方式比较：

if \[ -n $1 ] #$1 不为空

但如果用户不传参数的时候，$1 为空，这时 就会变成 \[ -n ] ,所以需要加一个辅助字符串来进行比较。

手工处理方式能满足大多数的简单需求，配合 shift 使用也能构造出强大的功能，但在要处理复杂选项的时候建议用下面的两种方法。

# getopts/getopt 工具

处理命令行参数是一个相似而又复杂的事情，为此 go 语言提供 cobra 库来实现，c++ 的 boost 提供了 options 库。

在 shell 中，处理此事的是 getopts 和 getopt 工具。

getopts 和 getopt 功能相似但又不完全相同，其中 getopt 是独立的可执行文件，而 getopts 是由 bash 内置的。

先来看看参数传递的典型用法:

    * ./test.sh -a -b -c  ： 短选项，各选项不需参数

    * ./test.sh -abc   ： 短选项，和上一种方法的效果一样，只是将所有的选项写在一起。

    * ./test.sh -a args -b -c ：短选项，其中-a需要参数，而-b -c不需参数。

    * ./test.sh --a-long=args --b-long ：长选项

先来看 getopts,它不支持长选项。

使用 getopts 非常简单，代码如下:

```bash
#test.sh
#!/bin/bash
while getopts "a:bc" arg; do #选项后面的冒号表示该选项需要参数
case $arg in
         a)
        echo "a's arg:$optarg" #参数存在$optarg中
         b)
        echo "b"
         c)
        echo "c"
         ?)  #当有不认识的选项的时候arg为?
    echo "unkonw argument"
exit 1
esac
done
```

现在就可以使用：
./test.sh -a arg -b -c&#x20;
或
./test.sh -a arg -bc
来加载了。

应该说绝大多数脚本使用该函数就可以了，如果需要支持长选项以及可选参数，那么就需要使用 getopt.

# getopt

## Syntax(语法)

** getopt optstring parameters**
** getopt \[OPTIONS] \[--] optstring parameters**
** getopt \[OPTIONS] -o|--options optstring \[options] \[--] parameters**

OPTIONS

- -a, --alternative Allow long options starting with single -
- -h, --help This small usage guide
- -l, --longoptions <LongOPTS> # 要被识别的长选项
- -n, --name <progname> The name under which errors are reported
- -o, --options <OPTString> # 要被识别的短选项
- -q, --quiet Disable error reporting by getopt(3)
- -Q, --quiet-output No normal output
- -s, --shell <shell> Set shell quoting conventions
- -T, --test Test for getopt(1) version
- -u, --unquoted Do not quote the output
- -V, --version Output version information

getopt 工具基本使用示例

```bash
#!/bin/bash
# 请注意，我们使用 "$@" 来使每个命令行参数扩展为一个单独的单词。 $@ 周围的引号是必不可少的！
# -o表示短选项，两个冒号表示该选项有一个可选参数，可选参数必须紧贴选项
# 如-carg 而不能是-c arg
# -l或--long选项后面接可接受的长选项，每个选项用逗号分开，冒号的意义同短选项。
# -n:出错时的信息
# -- ：举一个例子比较好理解：
#我们要创建一个名字为 "-f"的目录你会怎么办？
# mkdir -f #不成功，因为-f会被mkdir当作选项来解析，这时就可以使用
# mkdir -- -f 这样-f就不会被作为选项。
# 我们需要使用 temp 作为`eval set-`来抵消getopt的返回值。
temp=`getopt -o ab:c:: --long a-long,b-long:,c-long:: \
     -n 'example.bash' -- "$@"`
if [ $? != 0 ] ; then echo "terminating..." >&2 ; exit 1 ; fi
# 注意 $temp 周围的引号：它们是必不可少的！
# set 会重新排列参数的顺序，也就是改变$1,$2...$n的值，这些值在getopt中重新排列过了
eval set -- "$temp"
# 经过getopt的处理，下面处理具体选项。
while true ; do
    case "$1" in
        -a|--a-long) echo "option a" ; shift ;;
        -b|--b-long) echo "option b, argument \`$2'" ; shift 2 ;;
        -c|--c-long)
            # c 有一个可选参数。 由于我们处于引用模式，如果找不到可选参数，则会生成一个空参数。
            case "$2" in
                "") echo "option c, no argument"; shift 2 ;;
                *)  echo "option c, argument \`$2'" ; shift 2 ;;
            esac ;;
        --) shift ; break ;;
        *) echo "internal error!" ; exit 1 ;;
    esac
done
echo "remaining arguments:"
for arg do
   echo '--> '"\`$arg'" ;
done
```

上述代码会循环处理参数，每处理一个参数，就用 shift 命令剔除，直到所有参数全部处理完，则跳出循环

执行效果如下：

    ./parse.bash -a par1 'another arg' --c-long 'wow!*\?' -cmore -b " very long "
    option a
    option c, no argument
    option c, argument `more'
    option b, argument ` very long '
    remaining arguments:
    --> `par1'
    --> `another arg'
    --> `wow!*\?'

比如使用

./test -a -b arg arg1 -c

你可以看到,命令行中多了个 arg1 参数，在经过 getopt 和 set 之后，命令行会变为：

-a -b arg -c -- arg1

$1 指向 -a,$2 指向 -b,$3 指向 arg,$4 指向 -c,$5 指向 --,而多出的 arg1 则被放到了最后。

3，总结

一般小脚本手工处理也就够了，getopts 能处理绝大多数的情况，getopt 较复杂，功能也更强大!/bin/bash

```bash
function Help(){
    echo  ' ================================================================ '
    echo  ' 请使用下列 flags 来运行脚本，若无需改变默认值，使用 ./test.sh -- 即可；'
    echo  ' 如果环境无互联网，则安装的docker版本根据离线包的版本决定，dockerversion参数无用；'
    echo  ' --cgroupdriver: 指定 cgroupdrive 的类型，默认为 systemd；'
    echo  ' --dockerversion: 指定要安装的 docer-ce 版本，默认为 19.03.11；'
    echo  ' 使用示例:'
    echo  ' ./test.sh --cgroupdriver=systemd --dockerversion=19.03.11'
    echo  ' ================================================================'
}

case "$1" in
    -h|--help) Help; exit;;
esac

if [[ $1 == '' || ! $1 =~ '--' ]];then
    Help;
    exit;
fi

CMDFLAGS="$*"
for FLAGS in $CMDFLAGS;
do
    key=$(echo ${FLAGS} | awk -F"=" '{print $1}' )
    value=$(echo ${FLAGS} | awk -F"=" '{print $2}' )
    case "$key" in
      --cgroupdriver) CgroupDriver=$value ;;
      --dockerversion) DockerVersion=$value ;;
    esac
done

# Docker 配置
CgroupDriver=${CgroupDriver:-systemd}
DockerVersion=${DockerVersion:-19.03.11}
```
