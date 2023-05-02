---
title: set 命令
---

# 概述

> 参考：
> 
> - [GNU 文档，Bash 参考手册-Shell 内置命令-修改 Shell 行为-Set 内置命令](https://www.gnu.org/software/bash/manual/html_node/The-Set-Builtin.html)
> - <https://morven.life/posts/how-to-write-robust-shell-script/>

## 前言

Shell 脚本会有各种坑，经常导致 Shell 脚本因为各种原因不能正常执行成功。实际上，编写健壮可靠的 Shell 脚本也是有一定的技巧的。

```bash
# 在执行Shell脚本的时候，通常都会创建一个新的Shell，比如，当我们执行：
bash script.sh
```

Bash 会创建一个新的 Shell 来执行 script.sh，同时也默认给定了这个执行环境的各种参数。set 命令可以用来修改 Shell 环境的运行参数，不带任何参数的 set 命令，会显示所有的环境变量和 Shell 函数。我们重点介绍其中最常用的四个。

set -euxo pipefail

## set -x

默认情况下，Shell 脚本执行后只显示运行结果，不会展示结果是哪一行代码的输出，如果多个命令连续执行，它们的运行结果就会连续输出，导致很难分清一串结果是什么命令产生的。
set -x 用来在运行结果之前，先输出执行的那一行命令，行首以+表示是命令而非命令输出，同时，每个命令的参数也会展开，我们可以清晰地看到每个命令的运行实参，这对于 Shell 的 debug 来说非常友好。

```bash
#!/bin/bash
set -x
v=5
echo $v
echo "hello"
# output:
# + v=5
# + echo 5
# 5
# + echo hello
# hello
```

实际上，set -x 还有另一种写法 set -o xtrace。

## set -u

Shell 脚本不像其他高级语言，如 Python, Ruby 等，Shell 脚本默认不提供安全机制，举个简单的例子，Ruby 脚本尝试去读取一个没有初始化的变量的内容的时候会报错，而 Shell 脚本默认不会有任何提示，只是简单地忽略。

    #!/bin/bash
    echo $v
    echo "hello"
    # output:
    #
    # hello

可以看到，echo $v 输出了一个空行，Bash完全忽略了不存在的$v 继续执行后面的命令 echo "hello"。这其实并不是开发者想要的行为，对于不存在的变量，脚本应该报错且停止执行来防止错误的叠加。set -u 就用来改变这种默认忽略未定义变量行为，脚本在头部加上它，遇到不存在的变量就会报错，并停止执行。

    #!/bin/bash
    set -u
    echo $a
    echo bar
    # output:
    # ./script.sh: line 4: v: unbound variable

set -u 另一种写法是 set -o nounset

## set -e

对于默认的 Shell 脚本运行环境，有运行失败的命令（返回值非 0），Bash 会继续执行后面的命令：

    #!/bin/bash
    unknowncmd
    echo "hello"
    # output:
    # ./script.sh: line 3: unknowncmd: command not found
    # hello

可以看到，Bash 只是显示有错误，接着继续执行 Shell 脚本，这种行为很不利于脚本安全和排错。实际开发中，如果某个命令失败，往往需要脚本停止

set -e 从根本上解决了这个问题，它使得脚本只要发生错误，就终止执行：

    #!/bin/bash
    set -e
    unknowncmd
    echo "hello"
    # output:
    # ./script.sh: line 4: unknowncmd: command not found

可以看到，第 4 行执行失败以后，脚本就终止执行了。

set -e 根据命令的返回值来判断命令是否运行失败。但是，某些命令的非零返回值可能不表示失败，或者开发者希望在命令失败的情况下，脚本继续执行下去：

    #!/bin/bash
    set -e
    $(ls foobar)
    echo "hello"
    # output:
    # ls: cannot access 'foobar': No such file or directory

可以看到，打开 set -e 之后，即使 ls 是一个已存在的命令，但因为 ls 命令的运行参数 foobar 实际上并不存在导致命令的返回非 0 值，这有时候并不是我们看到的。

可以暂时关闭 set -e，该命令执行结束后，再重新打开 set -e：

    #!/bin/bash
    set -e
    set +e
    $(ls foobar)
    set -e
    echo "hello"
    # output:
    # ls: cannot access 'foobar': No such file or directory
    # hello

上面代码中，set +e 表示关闭-e 选项，set -e 表示重新打开-e 选项。

还有一种方法是使用 command || true，使得该命令即使执行失败，脚本也不会终止执行。

set -e 还有另一种写法 set -o errexit。

## set -o pipefail

set -e 有一个例外情况，就是不适用于管道命令。对于管道命令，Bash 会把最后一个子命令的返回值作为整个命令的返回值。也就是说，只要最后一个子命令不失败，管道命令总是会执行成功，因此它后面命令依然会执行，set -e 就失效了。

请看下面这个例子。

    #!/bin/bash
    set -e
    foo | echo "bar"
    echo "hello"
    # output:
    # ./script.sh: line 4: foo: command not found
    # bar
    # hello

可以看到，foo 是一个不存在的命令，但是 foo | echo bar 这个管道命令还是会执行成功，导致后面的 echo hello 会继续执行。

set -o pipefail 用来解决这种情况，只要一个子命令失败，整个管道命令就失败，脚本就会终止执行：

    #!/bin/bash
    set -e
    set -o pipefail
    foo | echo "bar"
    echo "hello"
    # output:
    # ./script.sh: line 5: foo: command not found
    # bar

可以看到，echo hello 命令并没有执行。

合并四个参数

对于上面提到的四个 set 命令参数，一般都放在一起使用。

# 写法一 set -euxo pipefail# 写法二 set -euxset -o pipefail

这两种写法任选其一放在所有 Shell 脚本的头部。

当然，也可以在在执行 Shell 脚本的时候，从 Bash 命令行传入这些参数：

bash -euxo pipefail script.sh

# Shell 脚本防御式编程

编写 Shell 脚本的时候应该考虑不可预期的程序输入，如文件不存在或者目录没有创建成功…其实 Shell 命令有很多选项可以解决这类问题，例如，使用 mkdir 创建目录的时候，如果父目录不存在，mkdir 默认返回错误，但如果加上-p 选项，mkdir 在父目录不存在的情况下先创建父目录；rm 在删除一个不存在的文件会失败，但如果加上-f 选项，即使文件不能存在也能执行成功。

注意字符串中的空格

我们必须时刻注意字符串中的空格字符，如文件名中的空格，命令参数中的空格等等，对于这些空格字符安全的最佳时实践是使用"括住相应的字符串：

# will fail if $filename contains spacesif \[ $filename = "foo" ];# will success even if $filename contains spacesif \[ "$filename" = "foo" ];

Someone will always use spaces in filenames or command line arguments and you should keep this in mind when writing shell scripts. In particular you should use quotes around variables.

if \[ $filename = “foo” ]; will fail if $filename contains a space. This can be fixed by using:

if \[ “$filename” = “foo” ];

类似的情况是，我们在使用$@或者其他包含由空格分割的多个字符串也要注意使用"括住相应的变量，实际上，使用"括住相应的变量没有任何副作用，只会是我们的 Shell 脚本更加健壮：

foo() { for i in $@; do printf "%s\n" "$i"; done }; foo bar "baz quux"barbazquuxfoo() { for i in "$@"; do printf "%s\n" "$i"; done }; foo bar "baz quux"barbaz quux

# 使用 trap 命令

关于 Shell 脚本一个常见的情况是，脚本执行失败导致文件系统处于不一致的状态，比如文件锁、临时文件或者 Shell 脚本的错误只更新了部分文件。为了达到“事务的完整性”我们需要解决这些不一致的问题，要么删除文件锁、临时文件，要么将状态恢复到更新之前的状态。实际上，Shell 脚本确实提供了一种在捕捉到特定的 unix 信号的情况下执行一段命令或者函数的功能：

trap command signal \[signal ...]

其实 Shell 脚本可以捕捉很多类型的信号（完整信号列表可以使用 kill -l 命令获取），但是我们通常只关心在问题发生之后用来恢复现场的三种信号：INT，TERM 和 EXIT

|        |                                                                                                                                                                                  |
| ------ | -------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| Signal | Description                                                                                                                                                                      |
| INT    | Interrupt – this signal is sent when someone kills the script by pressing ctrl-c                                                                                                 |
| TERM   | Terminate – this signal is sent when someone sends the TERM signal using the kill command                                                                                        |
| EXIT   | Exit – this is a pseudo-signal and is triggered when your script exits, either through reaching the end of the script, an exit command or by a command failing when using set -e |

一般情况下，我们在操作对应的共享区之前先创建文件锁：

```bash
if [ ! -e $lockfile ]; then
    touch $lockfile
    critical-section
    rm $lockfile
else
    echo "critical-section is already running"
fi
```

但是当 Shell 脚本操作对应的共享区的时候有人手动 Kill 掉对应的 Shell 脚本进程，这个时候文件锁的存在会导致 Shell 脚本不能再次操作对应的共享区。使用 trap 我们可以捕捉到对应的 Kill 信号并做相应的恢复操作：

```bash
if [ ! -e $lockfile ]; then
    trap "rm -f $lockfile; exit" INT TErm EXIT
    touch $lockfile
    $lockfile
    rm $lockfile
    trap - INT TErm EXIT
else
    echo "critical-section is already running"
fi
```

有了上面这段 trap 命令，即使当 Shell 脚本操作对应的共享区的时候有人手动 Kill 掉对应的 Shell 脚本进程，文件锁也会被清理干净。需要注意的是，我们在捕捉到 Kill 信号之后删除完文件锁之后直接退出而不是继续执行

。

Be Atomic

很多时候我们需要一次更新一批文件，但是有可能在更新了一半之后 Shell 脚本出错或者有人 kill 掉了 Shell 脚本的进程。你可能会想到，就使用刚才学到的 trap 知识，同时对就文件做备份，一旦捕捉到出错的信号，就恢复备份。这看起来没错，但是很多时候只能解决一部分的问题。例如，我们要把一个网站里面的 URL 从www.example.org全部更新为www.example.com，Shell脚本的主要逻辑类似于这样：

for file in $(find /var/www -type f -name "\*.html"); do perl -pi -e 's/www.example.org/www.example.com/' $filedone

正确的做法是尽量使更新操作原子化，实现操作的“事务一致性”：1. 拷贝旧目录；2. 在拷贝的目录中进行更新操作；3. 替换原目录

cp -a /var/www /var/www-tmpfor file in $(find /var/www-tmp -type f -name "\*.html"); do perl -pi -e 's/www.example.org/www.example.com/' $filedonemv /var/www /var/www-oldmv /var/www-tmp /var/www

在类 Unix 文件系统上进行最后的两次 mv 操作是非常快的（因为只需要替换两个目录的 inode，而不用执行实际的拷贝操作），换句话说，容易出错的地方是批量的更新操作，而我们将更新操作全部在拷贝的目录中执行，这样，更新操作即使出错，也不会影响原目录。这里的技巧是，使用双倍的硬盘空间来进行操作，任何是需要长时间打开文件的操作都是在备份目录中进行。事实上，保持一系列操作的原子性对于某些容易出错的 Shell 脚本来说非常重要，同时操作前备份文件也是一个好的编程习惯。
