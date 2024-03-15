---
title: find whereis which 查询工具
---

# which 查看可执行文件的位置

```bash
root@lichenhao:~# which ls
/usr/bin/ls
```

# whereis 查看文件的位置， 如 whereis ls

```bash
root@lichenhao:~# whereis ls
ls: /usr/bin/ls /usr/share/man/man1/ls.1.gz
```

# find # 在目录中搜索文件

## Syntax(语法)

**find \[OPTIONS] \[PATH...] \[EXPRESSION]**

- PATH # 路径名，不写上默认表示在当前路径下搜索
- OPTIONS 与 EXPRESSION # 详见下文，PATH 可以指定一个或多个想要搜索的路径。

如果没有给出 PATH，则使用当前目录；如果未给出任何表达式，则使用表达式 -print

### OPTIONS

-H，-L 和-P 选项控制 find 处理符号链接的行为。这些参数之后的命令行参数将被视为要检查的文件或目录的名称，直到以 - 开头的第一个参数或参数 ( 或 ! 。符号链接的概念详见：ln 链接

- **-P** # 不跟踪符号链接。该选项为 find 的默认行为。当 find 检查或打印文件信息时，该文件是符号链接，则所使用的信息应取自符号链接本身的属性。
- **-L** # 跟踪符号链接。与 -type l 同时使用时仅搜索失效的符号链接。当 find 检查或打印有关文件的信息时，所使用的信息应从链接指向的文件的属性中获取，而不是从链接本身获取（除非它是断开的符号链接，或者 find 无法检查与之相关的文件）链接点）。使用此选项意味着-noleaf。如果以后使用-P 选项，则-noleaf 仍然有效。如果-L 有效，并且 find 在搜索过程中发现到子目录的符号链接，则将搜索该符号链接所指向的子目录。
  - 当-L 选项生效时，-type 谓词将始终与符号链接指向的文件类型匹配，而不是与链接本身匹配（除非符号链接断开）。使用-L 会使-lname 和-ilname 谓词始终返回 false。
- **-H** # 除了处理命令行参数时，不要跟随符号链接。当 find 检查或打印有关文件的信息时，所使用的信息应取自符号链接本身的属性。唯一的例外情况是在命令行上指定的文件是符号链接并且可以解析该链接时，这种情况下，使用的信息取自链接指向的任何位置（即跟随该链接） ）。如果无法检查符号链接指向的文件，则有关链接本身的信息将用作备用。如果-H 有效，并且在命令行上指定的路径之一是指向目录的符号链接，则将检查该目录的内容（尽管-maxdepth 0 当然可以防止此情况）。
- **-D debugoptions** # Print diagnostic information; this can be helpful to diagnose problems with why find is not doing what you want. The list of debug options should be comma separated. Compatibility of the debug options is not guaranteed between releases of findutils. For a complete list of valid debug options, see the output of find -D help. Valid debug options include
- **-Olevel** # Enables query optimisation. The find program reorders tests to speed up execution while preserving the overall effect; that is, predicates with side effects are not reordered relative to each other. The optimisations performed at each optimisation level are as follows.

注意：

1. 上述五个 OPTIONS 必须出现在第一个 PATH 之前。
2. 下面 EXPRESSIONS(表达式)中的 OPTIONS，用来控制 find 的行为，与上述这些 OPTIONS 不同。并且表达式中的 OPTIONS 需要在第一个 PATH 之后指定。
3. find 默认将在当前目录下查找子目录与文件。并且将查找到的子目录和文件全部进行显示

### EXPRESSIONS(表达式)

EXPRESSIONS 由以下几部分组成

- OPTIONS 选项 # （这会影响整体操作，而不是特定文件的处理，并且始终返回 true）
- TESTS 测试 # （返回正确或错误的值）
- ACTIONS 动作 # （具有副作用并返回 true 或 false 值）

上述三部分表达式由运算符分隔。如果表达式除-prune 之外不包含其他任何动作，则对表达式为 true 的所有文件执行-print。

OPTIONS：

- **-type** # 根据文件类型查找 （find PATH -type 类型参数）
  - 参数类型包括：f 普通文件，l 软连接文件，d 目录文件等
  - 例: find -type l # 查找当前目录下的软连接文件
- **-name** # 根据文件名查找 （\* 任意多个字符）（? 1 个字符）（\[] 指范围值，外侧加引号）
  - 例：find /var –name "\*.log" # 查找/var 目录下，所有文件名最后字符是.log 的文件
- **-iname** # 不区分大小写根据文件名查找
- **-inum** # 根据 inode 号查找
- **-size** # 根据大小查找
  - 例:find /etc -size -10k -size +20k # 查找/etc 目录下小于 10k 大于 20k 的文件
  - 注意+-号，如果没有，是精确这么大，通常都会带上+或-号，表示一个范围
- **-user** # 根据所有者查找 （user 可改成 group 则查找属组为 lisi 的文件）
  - 例:find /home -user lisi # 查找/home 目录下，属主为 lisi 的文件
- **-perm** # 根据文件权限查找
  - 例:find /boot -perm 644 # 查找/boot 文件加下，权限为 644 的文件
- **-TIME {+|-}NUM** # 根据时间查找
  - - 表示该时间之前，-表示该时间之内
  - TIME 分为两部分
    - 第一部分，# 表示 ctime，atime，mtime；time 天，min 分钟
      - **c** # change # 表示属性被修改过：所有者、所属组、权限
      - **a** # access # 被访问过(被查看过)
      - **m** # modify # 表示内容被修改过
    - 第二部分
      - **time** # 天
      - **min** # 分钟
  - 逻辑连接符 # -a （and 逻辑与），-o （or 逻辑或)例：
- **-maxdepth \<INT>** # 设置最大目录层级
  - 例：find /home -maxdepth 1 -name \*.log # 查找/home 目录下 1 层的以.log 结尾的文件（即只查找当前目录，不查找当前目录下的子目录中的内容）
  - -ls 假设 find 指令的回传值为 Ture，就将文件或目录名称列出到标准输出
  - 例：find /home -type d -ls # 查找/home 目录下的文件夹，并列出这些文件夹的详细信息

ACTIONS 动作

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

## EXAMPLE

- 删除 /var/log/pods 目录下失效的符号链接
  - find -L /var/log/pods -type l -delete
- 找到/home 目录下的 samlee 用户的所有文件并删除
  - find /home -user samlee -exec rm –r {} ;
  - 注：rm -r 连带目录一起删除。报错原因：-exec 不适合大量传输，速率慢导致。
- 查看当前目录下两天之前修改过得并且 10 分钟之内查看过得文件
  - find -mtime +2 -amin -10
- 查看当前目录下 7 天之内并且 2 天之前的文件
  - find -mtime -7 -a -mtime +2
- 在 etc 目录下查找大于 1k 并且小于 10k 的文件
  - find /etc -size +1k -a -size -10k
- -inum 根据文件 i 节点查询
  - find ./ -inum 2310630 -exec rm {} ; # 有一些文件的硬链接数量很多，有相同的 i 节点，查找其中一个文件的 i 节点号，一次性删除。
- find -maxdepth 1 -mtime +1 -exec ls -l {} ; # 查找出来的如果是文件夹，那么就对这个文件夹执行该命命令，如下图所示，查找出来./test 等 3 个文件夹，那么就这三个文件夹执行 ls -l 的命令，即
  - ls -l ./test
  - ls -l ./lichenhao
  - ls -l ./lost+found
  - 注：查找出来的文字，一字不差的全部放在后面的{}中，等待 COMMAND 执行，所以没法列出目录的详细信息
  - 注意与-ls 参数的区别

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/idlpqz/1616166302372-29738c6b-f92a-43b5-9243-9aa5483629ca.jpeg)

### 查找重复文件

在工作生活当中，我们很可能会遇到查找重复文件的问题。比如从某游戏提取的游戏文本有重复的，我们希望找出所有重复的文本，让翻译只翻译其中一份，而其他的直接替换。那么这个问题该怎么做呢？当然方法多种多样，而且无论那种方法应该都不会太难，但笔者第一次遇到这个问题的时候第一反应是是用 Linux 的 Shell 脚本，所以文本介绍这种方式。
先上代码：
`find -not -empty -type f -printf "%sn" | sort -rn |uniq -d | xargs -I{} -n1 find -type f -size {}c -print0 | xargs -0 md5sum | sort | uniq -w32 --all-repeated=separate | cut -b 36-`
大家先 cd 到自己想要查找重复文件的文件夹，然后 copy 上面代码就可以了，系统会对当前文件夹及子文件夹内的所有文件进行查重。
下面分析一下上面的命令。
首先看第一句：
`find -not -empty -type f -printf "%sn"`
find 是查找命令；-not -empty 是要寻找非空文件；-type f 是指寻找常规文件；-printf “%sn”比较具有迷惑性，这里的%s 并非 C 语言中的输出字符串，它实际表示的是文件的大小，单位为 bytes（不懂就 man，man 一下 find，就可以看到了），n 是换行符。所以这句话的意思是输出所有非空文件的大小。
通过管道，上面的结果被传到第二句：
`sort -rn`
sort 是排序，-n 是指按大小排序，-r 是指从大到小排序（逆序 reverse）。
第三句：
`uniq -d`
uniq 是把重复的只输出一次，而-d 指只输出重复的部分（如 9 出现了 5 次，那么就输出 1 个 9，而 2 只出现了 1 次，并非重复出现的数字，故不输出）。
第四句：
`xargs -I{} -n1 find -type f -size {}c -print0`
这一部分分两部分看，第一部分是 xargs -I{} -n1，xargs 命令将之前的结果转化为参数，供后面的 find 调用，其中-I{}是指把参数写成{}，而-n1 是指将之前的结果一个一个输入给下一个命令（-n8 就是 8 个 8 个输入给下一句，不写-n 就是把之前的结果一股脑的给下一句）。后半部分是 find -type f -size {}c -print0，find 指令我们前面见过，-size{}是指找出大小为{}bytes 的文件，而-print0 则是为了防止文件名里带空格而写的参数。
第五句：
`xargs -0 md5sum`
xargs 我们之前说过，是将前面的结果转化为输入，那么这个-0 又是什么意思？man 一下 xargs，我们看到-0 表示读取参数的时候以 null 为分隔符读取，这也不难理解，毕竟 null 的二进制表示就是 00。后面的 md5sum 是指计算输入的 md5 值。
第六句：sort 是排序，这个我们前面也见过。
第七句：
`uniq -w32 --all-repeated=separate`
uniq -w32 是指寻找前 32 个字符相同的行，原因在于 md5 值一定是 32 位的，而后面的--all-repeated=separate 是指将重复的部分放在一类，分类输出。
第八句：
`cut -b 36-`
由于我们的结果带着 md5 值，不是很好看，所以我们截取 md5 值后面的部分，cut 是文本处理函数，这里-b 36-是指只要每行 36 个字符之后的部分。
我们将上述每个命令用管道链接起来，存入 result.txt：
`find -not -empty -type f -printf "%sn" | sort -rn |uniq -d | xargs -I{} -n1 find -type f -size {}c -print0 | xargs -0 md5sum | sort | uniq -w32 --all-repeated=separate | cut -b 36- >result.txt`
虽然结果很好看，但是有一个问题，这是在 Linux 下很好看，实际上如果有朋友把输出文件放到 Windows 上，就会发现换行全没了，这是由于 Linux 下的换行是 n，而 windows 要求 nr，为了解决这个问题，我们最后执行一条指令，将 n 转换为 nr：
`cat result.txt | cut -c 36- | tr -s 'n'`
