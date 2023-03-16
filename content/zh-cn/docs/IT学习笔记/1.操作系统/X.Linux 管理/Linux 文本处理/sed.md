---
title: sed
---

# 概述

> 参考：
> - 官方文档：<https://www.gnu.org/software/sed/>
> - 官方手册：<https://www.gnu.org/software/sed/manual/sed.html>
> - 官方文档：<https://www.gnu.org/software/sed/manual/sed.html#Execution-Cycle>
> - <https://opus.konghy.cn/sed-refer/>
> - [https://github.com/liquanzhou/ops_doc/blob/master/shell%E5%AE%9E%E4%BE%8B%E6%89%8B%E5%86%8C.sh#L2925](https://github.com/liquanzhou/ops_doc/blob/master/shell%25E5%25AE%259E%25E4%25BE%258B%25E6%2589%258B%25E5%2586%258C.sh#L2925)
> - <https://mp.weixin.qq.com/s/tKvg69WvAFLJSfHsgRe1Yw>

sed 是一种新型的，非交互式的编辑器。它能执行与编辑器 vi 相同的编辑任务。

sed 是一种 **stream editor(流编辑器)**，逐行处理文件(或输入)，并将结果发送到标准输出。处理时，把当前处理的行存储在临时缓冲区中，称为 **pattern space( 模式空间)**，接着用 sed 命令处理缓冲区中的内容，处理完成后，把缓冲区的内容送往屏幕。然后读入下行，执行下一个循环。如果没有使诸如‘D’的特殊命令，那会在两个循环之间清空模式空间，但不会清空保留空间。这样不断重复，直到文件末尾。文件内容并没有改变，除非你使用重定向存储输出。

功能：主要用来自动编辑一个或多个文件,简化对文件的反复操作,编写转换程序等。
![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/ec3zxx/1639446629207-32f3f321-c8c9-429d-880d-7098bc7ceed1.png)

## sed 工作原理

sed 实际上是一个循环结构，该循环用来对输入给 sed 文本的每一行执行下列操作：

1. sed 从输入内容中读取一行，删除任何尾随的换行符，并将该行保存到 模式空间 中。
2. 对模式空间中的内容，执行 COMMAND。执行之前需要根据 ADDR(行定位)，验证模式空间中的行是否符合 ADDR 中给定的条件。
   1. 只有符合条件的，COMMAND 才会执行
3. 上一步完成后，模式空间中的内容会打印到标准输出中。如果删除了末尾的换行符，则会添加回去。(如果使用了 -n 选项，则不会将模式空间中的行输出)
4. 然后回到循环体的开头，继续处理输入内容的下一行

除非有特殊命令（例如“d 使用'），则在两个循环之间删除模式空间。另一方面，保持空间可在周期之间保持其数据（请参见命令'H'，'H'，'X'，'G'，'G'在两个缓冲区之间移动数据）

# Syntax(语法)

**sed \[OPTIONS] SCRIPT FILE(s)**

## OPTIONS

- **-e, --expression=\<SCRIPT>** # 以选项中的指定的 SCRIPT 来处理输入的文本文件，常用来在一行命令中，执行两个 sed SCRIPT
- **-f, --file=\<SCRIPT>** # 以选项中指定的 SCRIPT 文件来处理输入的文本文件,把 sed 相关命令写进文件里，直接引用该文件中的命令进行操作
- **-i** # 直接编辑原文件，sed 操作的内容不输出到屏幕，直接更改文件内容
- **-n, --quiet, -silent** # 禁止模式空间中的内容在标准输出中打印
  - 通常与 p 命令一同使用，用来仅显示 sed 操作的行。
- **-r, --regexp-extended** # 允许在 SCRIPT 中使用扩展的正则表达式。如果在 SCRIPT 中使用正则，且不使用该选项，运行就会报错

## SCRIPT

SCRIPT 是 sed 在处理文本时主要依赖的部分，脚本包含多个部分，至少要具有一个 COMMAND

SCRIPT 语法：**\[ADDR]COMMAND\[OPTIONS]**

- **ADDR **# 行定位，**全称 Addresses**。用于确定 sed 当前操作的文本需要处理哪些行。
  - 如果指定了 ADDR ，则 COMMAND 仅对被定位的行执行操作。
  - ADDR 可以是单个行号、通过 pattern(正则表达式) 来匹配指定的行、通过 X,Y 来匹配一个范围内的行
- **COMMAND** # 用于执行通过 行定位 匹配到的行的操作。是添加内容、还是替换内容、还是删除内容等等
- **OPTIONS** # 选项仅在 COMMAND 有可用的 OPTIONS 时才有用。比如 s 命令具有多个 OPTIONS

Note：行定位与 COMMAND 不分先后，不分左右，不同的 COMMAND，会出现在 SCRIPT 不同的位置

### 定界符

一个复杂的 SCRIPT，需要定界符来区分语法格式中每一部分。在 sed 中可以用任意字符作为定界符。sed 会自动将第一个出现的非命令字符作为定界符。一般情况定界符使用 / 。示例如下：

- sed 's:test:TEXT:g' file
- sed 's|test|TEXT|g' file

上面的示例通过定界符，区分了 COMMAND、正则表达式、要替换的内容、COMMAND 对应的 flag
定界符出现在样式内部时，需要进行转义：

1. sed 's//bin//usr/local/bin/g' file

### [ADDR(行定位)](https://www.gnu.org/software/sed/manual/sed.html#sed-addresses)

sed 工具将会更具 ADDR 来决定在哪一行或哪些行执行命令，比如下面的示例表示：仅在第 144 行将字符串 hello 替换为 world

```bash
sed '144s/hello/world/' input.txt > output.txt
```

ADDR 可以通过多种方式进行行定位：
一、使用数值定位行

- **NUMBER** # 使用行号 NUMBER，来定位指定的行
- **$** # 这个符号与输入的最后一个文件的最后一行匹配，或者在指定-i 或-s 选项时与每个文件的最后一行匹配。
- **first~step** # 从 first 行开始，每隔 step 的行，被定位

EXAMPLE：

- sed -n 10p passwd # 输出 passwd 文件中的第 10 行
- sed -n $p passwd # 输出 passwd 文件中最后一行
- sed -n 1~5p passwd # 输出 passwd 文件中第一行、第六行、第十一行....以此类推

二、使用 RegExp(正则表达式) 定位行
Note：通过正则表达式来匹配指定的行，正则表达式两边需要添加定界符

- **/RegExp/** # 根据 RegExp 匹配到内容，来定位包含这些内容的行
- **/RegExp/I** # 进行正则匹配时，不区分大小写
- **/RegExp/M** # 正则表达式匹配的 M 修饰符是 GNU sed 扩展，它指示 GNU sed 在多行模式下匹配正则表达式

EXAMPLE

- sed -n /root/p passwd # 输出 passwd 文件中所有带 root 字符的行

三、定位一个范围内的行

- **NUM1,NUM2** # 使用行号定位 NUM1 到 NUM2 的所有行
- **NUM,+N** # 定位 NUM1 行及其后两行
- **NUM,~N** # 匹配 NUM1 和 NUM1 之后的行，直到行号是 N 的倍数的下一行为止。以下命令从第 6 行开始打印，直到下一行是 4 的倍数（即第 8 行）：
  - seq 10 | sed -n '6,~4p' # 输出 6 7 8 行
- 注意：任意 NUM 可以 使用 RegExp 代替，i.e.通过正则来匹配开始行或结束行
  - **/RegExp1/,/RegExp2/** # 使用正则，定位 pattern1 匹配到的行，到 pattern2 匹配到的行，这两行中间的所有行

EXAMPLE

- sed -n 1,5p passwd # 输出 passwd 文件中的第一行到第五行
- sed -n /root/,/sshd/p passwd # 输出带有 root 的行，到 带有 sshd 的行，中间的所有行
- sed -n /root/,+2p /etc/passwd # 输出带有 root 的行，以及带有 root 行下面的 2 行

### COMMAND(操作行为)

- **a \<TEXT>** # 在匹配行的下一行插入 TEXT
- **i \<TEXT>** # 在匹配行的上一行插入 TEXT
- **a\ \<TEXT>** # 在当前行下面插入 TEXT。
- **i\ \<TEXT>** # 在当前行上面插入 TEXT。
- **c\ \<TEXT>** # 把选定的行改为新的 TEXT。
- **d** # 删除，删除选择的行。
- **D** # 删除模板块的第一行。
  - 删除命令用于删除匹配的行，而且删除命令还会改变 sed 脚本中命令的执行操作顺序，因为匹配的行一旦被删除，模式空间将变为“空”，自然不会再执行哪个 sed 脚本后续的命令。删除命令会导致读取新的输入行（下一行），而 sed 脚本中的命令则从头开始执行。需要注意的是删除时是删除整行，而不是删除匹配的内容（如要删除匹配的内容，可以使用替换）。
- **s/RegExp/REPLACEMENT/FLAGS** # 在已经定位的每行中，将 RegExp 匹配到的内容，替换成 REPLACEMENT
- **h** # 拷贝模板块的内容到内存中的缓冲区。
- **H** # 追加模板块的内容到内存中的缓冲区。
- **g** # 获得内存缓冲区的内容，并替代当前模板块中的文本。
- **G** # 获得内存缓冲区的内容，并追加到当前模板块文本的后面。
- **l** # 列表不能打印字符的清单。
- **L** # 同 l，不显示非打印字符
- **n** # 读取下一个输入行，用下一个命令处理新的行而不是用第一个命令。
- **N** # 追加下一个输入行到模式空间中，并在二行间嵌入一个新行，改变当前行号码。
- **p **# 打印模式空间中执行了 COMMAND 的内容
- **P** # 打印模板块的第一行。
- **q** # 退出 Sed。
- b lable 分支到脚本中带有标记的地方，如果分支不存在则分支到脚本的末尾。
- r file 从 file 中读行。
- t label if 分支，从最后一行开始，条件一旦满足或者 T，t 命令，将导致分支到带有标号的命令处，或者到脚本的末尾。
- T label 错误分支，从最后一行开始，一旦发生错误或者 T，t 命令，将导致分支到带有标号的命令处，或者到脚本的末尾。
- w file 写并追加模板块到 file 末尾。
- W file 写并追加模板块的第一行到 file 末尾。
- ! # 表示后面的命令对所有没有被定位到的行发生作用。i.e.对行定位操作匹配到的行取反。
- \= 打印当前行号码。
- # 把注释扩展到下一个换行符以前。

# 特殊 COMMAND

## 替换指令(S, SUBSTITUTION)

指令格式：\[address]s/pattern/replacement/flags

- address # 操作地址
- s # 替换指令
- pattern # 匹配需要替换的内容
- replacement # 为替换的内容
- flags # 标记可以是如下内容：
  - n # 1 - 512 之间的数字，表示对模式空间中指定模式的第 n 次出现进行替换。如一行中有 3 个 A，而只想替换第二个 A。
  - g # 对模式空间的所有匹配进行全局更改。没有 g 则只有第一次匹配被替换。如一行中有 3 个 A，则仅替换第一个 A。
  - p # 打印模式空间的内容，即表示打印行。与-n 选项一起使用可以只打印匹配的行。
  - w file # 将模式空间的内容写到文件 file 中。 即表示把行写入一个文件。

replacement 为字符串，用来替换这则表达式匹配的内容。在 replacement 部分，下列字符有特殊含义：

- & 用正则表达式匹配的内容进行替换
- \n 匹配第 n 个子串，该子串之前在 pattern 中用 `\(\)`指定，即正则表达式分组。
- \ 转义（转义替换部分包含：&、\等）

EXAMPLE

- set -n 's/123/234/'p test # 把 test 文件中的 123 替换为 234，并显示所有完成替换的行
- sed 's/|/\n/' FILE # 把 FILE 文件中的每个|符号变成换行符，即让一行的内容变成多行
- sed ':t;N;s/\n//;b t' FILE # 把 FILE 问文件中的换行符清空，即然多行内容变为一行
- sed 's/^\[ \t]\*//' FILE # 删除行首 tab 键
- echo this is a test line | sed 's/\w+/\[&]/g' # 将所有的单词用中括号 \[] 包裹起来：

## 转换指令（Y）

按字符转换（Transform）的语法格式为：\[address]y/yousource-chars/dest-chars/

- address 用于定位需要修改的行
- source-chars 为需要修改的字符
- dest-chars 为准备替换的字符。

EXAMPLE

- 就文件中的 china 转换为大写：
  - sed '/china/y/abcdefghijklmnopqrstuvwxyz/ABCDEFGHIJKNOPQRSTUVWXYZ/' file

# 其他

## 命令组合

用时候我们可能会对一个文件或者输入做连续的 sed 处理，例如：

sed '表达式' | sed '表达式' | sed '表达式'

其实，我们可以将所有的表达式用分号(;)组合起来，即：

sed '表达式; 表达式;表达式'

分号的含义就是将前面的处理完结果传给后边的表达式继续处理。

我们也可以对多个命令用大括号 {} 进行组合，然后作用于同一匹配地址，即：

address{commad1; command2; command3}

也可以放在多行：

address{ commad1 command2 command3}

示例，如果文件中含 test 的行，则将其下一行的 aa 替换为 bb：

sed '/test/{ n; s/aa/bb/; }' file

选定行的范围

选择要处理行的范围，可以用逗号(,)来分割。例如选定所有在模板 test 和 check 所确定的范围内的行：

sed -n '/test/,/check/p' file

打印从第 5 行开始到第一个包含以 test 开始的行之间的所有行：

sed -n '5,/^test/p' file

对于模板 test 和 west 之间的行，每行的末尾用字符串 aaa bbb 替换：

sed '/test/,/west/s/$/aaa bbb/' file

## Sed 常用格式

sed 命令行常用的基本格式大致有一下三种形式：

（1）Sed \[options] 'script' file1 file2 ...

script 结构为 /PATTERN/action，PATTERN 为正则表达式，action 为要执行的动作。例如：sed ‘/\[\[:upper:]]/d’ binary.sh 表示删除所有的大写字母的行。script 的结构还可以是 /PATTERN1/,/PATTERN2/action，这表示从第一次被 PATTERN1 匹配到的行到第一次被 PATTERN2 匹配到的中间的所有行执行 action 动作。

这里需要注意的是，当进行字符串替换时，需要在 PATTERN 前加上 s 动作，并在末尾加上替换的范围，例如： sed 's/abc/123/g' 表示将匹配到的字符串 abc 替换成 123，g 表示替换所有的行。

（2）Sed \[options] –f scriptfile file1 file2 ...

scriptfile 表示脚本文件，即 sed 支持将要执行的操作写在文件里边，然后通过使用 -f 参数来加载文件。

（3）Sed \[options] 'ADDR1,ADDR2command' file1 file2 ...

该格式应用于以行为单位的操作，例如：

sed ’1,2d’ file

就可以将 file 的前两行删除并显示出来，但是它不会改变源文件。

Sed ‘1,2!d’ file

表示删除除第一行和第二行之外的所有行。

注： 在这种格式中的 & 表示引用前面匹配到的所有字符。并且在该种格式中可以引入分组。

示例：

$ sed 's/bc/-&-/' testfile

这里表示在匹配到 bc 字符两端加上字符'-'

$ sed 's/ /-\1-~\2~/' testfile

这里的 \1 和 \2 表示正则表达式的分组 1 和分组二所匹配的内容。

## Sed 高级应用

正常的 Sed 数据处理流程是读取文档的一行至模式空间，然后对该行应用相应的 Sed 指令。当指令完成后输出该行并清空模式空间，一次循环读入文档的下一行数据，直至文档数据结尾。然而在真实环境中的数据可能并不会那么有规律，有时我们会把数据分多行写入文档，如：

姓名：张三邮箱：zhangsan@gmail.com姓名：李四邮箱：lisi@gmail.com

从上面的模板文件中可以看出，实际每两行位一条完整的记录，而此时如果需要用 Sed 对文档进行处理，就需要对 Sed 工作流程进行人工干预。

多行操作 Next

Next（N）指令通过读取新的输入行，并将它追加至模式空间的现有内容之后，来创建多行模式空间。模式空间的最初内容与新的输入含之间用换行符分隔。在模式空间中插入的换行符可以用 \n 匹配。

列举一个范例，范例所用样本文件如下（test.txt）：

Name:HuotyMail:huoty@gmail.comName:KonghyMail:konghy@163.com

我们要做的处理是，当读入的内容与 Name 匹配时，立刻读取下一行，再输入模式空间中的内容。处理脚本如下所示（sed.sh）:

\#n/Name/{NL}

其中，#n 放在脚本文件中表示屏蔽自动输出，L 表示不打印非打印字符（小写 l 标识打印非打印字符），即行尾的 \n。用 sed 执行操作如下：

sed -f sed.sh test.txt

多行操作 Print

Print（p）表示仅输出多行模式空间中的第一部分直到第一个插入的 \n 换行符为止。如模式空间中的内容为 “aaa\nbbb”，则 P 只输出 aaa。

多行删除 Delete（D）

Delete 删除模式空间中直到第一个插入的换行符（\n）前的内容。由于 d 命令的作用是删除模式空间中的内容并读取新的输入行，而如果 sed 在 d 指令后还有多条命令，则余下的指令将不再执行。而返回第一条指令对新度入行进行处理。多行指令 D 则不会读入因的行，而是放回 sed 脚本的顶端，使得剩余指令继续应用于模式空间中的剩余部分内容。

Hold（h,H）, Get（g,G）

Sed 还有一个称为保持空间（hold space）的缓冲区。模式空间的内容可以复制到保持空间，保持空间同样可以复制到模式空间。由一组 Sed 命令用于两者之间移动数据：

Hold（h|H） 将模式空间的内容复制或者追加到保持空间 Get（g|G） 将保持空间的内容复制或者追加到模式空间 Exchange（x） 交换保持空间与模式空间中的内容

举一个使用范例，样本文件如下（test.txt）：

aaabbbcccddd

Sed 教程文件如下（sed.sh）：

/aaa/{hd}/ccc/{G}

执行处理命令：

sed -f sed.sh test.txt

结果如下所示：

bbbcccaaaddd

sed{

# 先读取资料、存入模式空间、对其进行编辑、再输出、再用下一行替换模式空间内容

# 调试工具 sedsed (参数 -d) <http://aurelio.net/sedsed/sedsed-1.0>

-n # 输出由编辑指令控制(取消默认的输出,必须与编辑指令一起配合)

-i # 直接对文件操作

-e # 多重编辑

-r # 正则可不转移特殊字符

b # 跳过匹配的行

p # 打印

d # 删除

s # 替换

g # 配合 s 全部替换

i # 行前插入

a # 行后插入

r # 读

y # 转换

q # 退出

& # 代表查找的串内容

- # 任意多个 前驱字符(前导符)

? # 0 或 1 个 最小匹配 没加-r 参数需转义 ?

$ # 最后一行

.\* # 匹配任意多个字符

(a) # 保存 a 作为标签 1(\1)

模式空间{

    # 模式空间(两行两行处理) 模式匹配的范围，一般而言，模式空间是输入文本中某一行，但是可以通过使用N函数把多于一行读入模式空间

    # 暂存空间里默认存储一个空行

    n   # 读入下一行(覆盖上一行)

    h   # 把模式空间里的行拷贝到暂存空间

    H   # 把模式空间里的行追加到暂存空间

    g   # 用暂存空间的内容替换模式空间的行

    G   # 把暂存空间的内容追加到模式空间的行后

    x   # 将暂存空间的内容于模式空间里的当前行互换

    ！  # 对其前面的要匹配的范围取反

    D   # 删除当前模式空间中直到并包含第一个换行符的所有字符(/.*/匹配模式空间中所有内容，匹配到就执行D,没匹配到就结束D)

    N   # 追加下一个输入行到模式空间后面并在第二者间嵌入一个换行符，改变当前行号码,模式匹配可以延伸跨域这个内嵌换行

    p   # 打印模式空间中的直到并包含第一个换行的所有字符

}

标签函数{

    : lable # 建立命令标记，配合b，t函数使用跳转

    b lable # 分支到脚本中带有标记的地方，如果分支不存在则分支到脚本的末尾。

    t labe  # 判断分支，从最后一行开始，条件一旦满足或者T,t命令，将导致分支到带有标号的命令出，或者到脚本末尾。与b函数不同在于t在执行跳转前会先检查其前一个替换命令是否成功，如成功，则执行跳转。

    sed -e '{:p1;/A/s/A/AA/;/B/s/B/BB/;/[AB]\{10\}/b;b p1;}'     # 文件内容第一行A第二行B:建立标签p1;两个替换函数(A替换成AA,B替换成BB)当A或者B达到10个以后调用b,返回

    echo 'sd  f   f   [a    b      c    cddd    eee]' | sed ':n;s#\(\[[^ ]*\)  *#\1#;tn'  # 标签函数t使用方法,替换[]里的空格

    echo "198723124.03"|sed -r ':a;s/([0-9]+)([0-9]{3})/\1,\2/;ta'  # 每三个字符加一个逗号

}

引用外部变量{

    sed -n ''$a',10p'

    sed -n ""$a",10p"

}

sed 10q # 显示文件中的前 10 行 (模拟"head")

sed -n '$=' # 计算行数(模拟 "wc -l")

sed -n '5,/^no/p' # 打印从第 5 行到以 no 开头行之间的所有行

sed -i "/^$f/d" a 　　 　 # 删除匹配行

sed -i '/aaa/,$d' # 删除匹配行到末尾

sed -i "s/=/:/" c # 直接对文本替换

sed -i "/^pearls/s/$/j/" # 找到 pearls 开头在行尾加 j

sed '/1/,/3/p' file # 打印 1 和 3 之间的行

sed -n '1p' file # 取出指定行

sed '5i\aaa' file # 在第 5 行之前插入行

sed '5a\aaa' file # 在第 5 行之后抽入行

echo a|sed -e '/a/i\b' # 在匹配行前插入一行

echo a|sed -e '/a/a\b' # 在匹配行后插入一行

echo a|sed 's/a/&\nb/g' # 在匹配行后插入一行

seq 10| sed -e{1,3}'s/./a/' # 匹配 1 和 3 行替换

sed -n '/regexp/!p' # 只显示不匹配正则表达式的行

sed '/regexp/d' # 只显示不匹配正则表达式的行

sed '$!N;s/\n//' # 将每两行连接成一行

sed '/baz/s/foo/bar/g' # 只在行中出现字串"baz"的情况下将"foo"替换成"bar"

sed '/baz/!s/foo/bar/g' # 将"foo"替换成"bar"，并且只在行中未出现字串"baz"的情况下替换

echo a|sed -e 's/a/#&/g' # 在 a 前面加#号

sed 's/foo/bar/4' # 只替换每一行中的第四个字串

sed 's/(.\*)foo/\1bar/' # 替换每行最后一个字符串

sed 's/(.*)foo(.*foo)/\1bar\2/' # 替换倒数第二个字符串

sed 's/\[0-9]\[0-9]$/&5' # 在以\[0-9]\[0-9]结尾的行后加 5

sed -n ' /^eth|em\[01]\[^:]/{n;p;}' # 匹配多个关键字

sed -n -r ' /eth|em\[01]\[^:]/{n;p;}' # 匹配多个关键字

echo -e "1\n2"|xargs -i -t sed 's/^/1/' {} # 同时处理多个文件

sed '/west/,/east/s/$/_VACA_/' # 修改 west 和 east 之间的所有行，在结尾处加*VACA*

sed 's/\[^1-9]_(\[0-9]+)._/\1/' # 取出第一组数字，并且忽略掉开头的 0

sed -n '/regexp/{g;1!p;};h' # 查找字符串并将匹配行的上一行显示出来，但并不显示匹配行

sed -n ' /regexp/{n;p;}' # 查找字符串并将匹配行的下一行显示出来，但并不显示匹配行

sed -n 's/(mar)got/\1ianne/p' # 保存(mar)作为标签 1

sed -n 's/(\[0-9]+).\*(t)/\2\1/p' # 保存多个标签

sed -i -e '1,3d' -e 's/1/2/' # 多重编辑(先删除 1-3 行，在将 1 替换成 2)

sed -e 's/@.\*//g' -e '/^$/d' # 删除掉@后面所有字符，和空行

sed -n -e "{s/^ _\[0-9]_//p}" # 打印并删除正则表达式的那部分内容

echo abcd|sed 'y/bd/BE/' # 匹配字符替换

sed '/^#/b;y/y/P/' 2 # 非#号开头的行替换字符

sed '/suan/r readfile' # 找到含 suan 的行，在后面加上读入的文件内容

sed -n '/no/w writefile' # 找到含 no 的行，写入到指定文件中

sed '/regex/G' # 在匹配式样行之后插入一空行

sed '/regex/{x;p;x;G;}' # 在匹配式样行之前和之后各插入一空行

sed 'n;d' # 删除所有偶数行

sed 'G;G' # 在每一行后面增加两空行

sed '/^$/d;G' # 在输出的文本中每一行后面将有且只有一空行

sed 'n;n;n;n;G;' # 在每 5 行后增加一空白行

sed -n '5~5p' # 只打印行号为 5 的倍数

seq 1 30|sed '5~5s/.\*/a/' # 倍数行执行替换

sed -n '3,${p;n;n;n;n;n;n;}' # 从第 3 行开始，每 7 行显示一次

sed -n 'h;n;G;p' # 奇偶调换

seq 1 10|sed '1!G;h;$!d' # 倒叙排列

ls -l|sed -n '/^.rwx.\*/p' # 查找属主权限为 7 的文件

sed = filename | sed 'N;s/\n/\t/' # 为文件中的每一行进行编号(简单的左对齐方式)

sed 's/^\[ \t]\*//' # 将每一行前导的"空白字符"(空格，制表符)删除,使之左对齐

sed 's/^\[ \t]_//;s/\[ \t]_$//' # 将每一行中的前导和拖尾的空白字符删除

sed '/{abc,def}/\[111,222]/s/^/00000/' # 匹配需要转行的字符: } / \[

echo abcd\nabcde |sed 's/\n/@/g' |tr '@' '\n' # 将换行符转换为换行

cat tmp|awk '{print $1}'|sort -n|sed -n '$p' # 取一列最大值

sed -n '{s/^\[^/]_//;s/:._//;p}' /etc/passwd # 取用户家目录(匹配不为/的字符和匹配:到结尾的字符全部删除)

sed = filename | sed 'N;s/^/ /; s/ \*(.{6,})\n/\1 /' # 对文件中的所有行编号(行号在左，文字右端对齐)

/sbin/ifconfig |sed 's/._inet addr:(._) Bca.\*/\1/g' |sed -n '/eth/{n;p}' # 取所有 IP

修改 keepalive 配置剔除后端服务器{

    sed -i '/real_server.*10.0.1.158.*8888/,+8 s/^/#/' keepalived.conf

    sed -i '/real_server.*10.0.1.158.*8888/,+8 s/^#//' keepalived.conf

}

模仿 rev 功能{

    echo 123 |sed '/\n/!G;s/\(.\)\(.*\n\)/&\2\1/;//D;s/.//;'

    /\n/!G;         　　　　　　# 没有\n换行符，要执行G,因为保留空间中为空，所以在模式空间追加一空行

    s/\(.\)\(.*\n\)/&\2\1/;     # 标签替换 &\n23\n1$ (关键在于& ,可以让后面//匹配到空行)

    //D;            　　　　　　# D 命令会引起循环删除模式空间中的第一部分，如果删除后，模式空间中还有剩余行，则返回 D 之前的命令，重新执行，如果 D 后，模式空间中没有任何内容，则将退出。  //D 匹配空行执行D,如果上句s没有匹配到,//也无法匹配到空行, "//D;"命令结束

    s/.//;          　　　　　　# D结束后,删除开头的 \n

}

}

# 应用示例

- 获取双引号之间的字符。假如 sed.txt 中的内容为 `"bitnami "` 。那么下面命令会输出 `bitnami`
  - **sed 's/^"(._)"._/\1/' sed.txt **
- 打印 passwd 文件的内容，等效于 cat passwd 命令。
  - sed -n p passwd
  - 注意： 如果不加 -n，则 passwd 每行内容输出两次，因为 sed 本身的逻辑在从模式空间到标准输出一行，然后 p 命令还会再将该行输出一遍。
- 删除开头带#的行
  - sed '/^#/d' FILE
- 删掉空白行：
  - sed '/^$/d' file
- 搜索 resolv.conf 文件中，开头带有 nameserver 字符串的行，并在行首添加#
  - sed '/^nameserver/s/^/#/' /etc/resolv.conf
- 将 resolv.conf 文件中，具有 nameserver 关键字的行开头的 # 符号去掉。
  - sed 's/#(nameserver.\*)/\1/' /etc/resolv.conf
- 在 hostname 行的前一行添加 ${STRING} 变量中的内容。其中 `STRING="\ \ \ \ \ \ labels:"`
  - sed -i "/hostname/i${STRING}" prometheus.yml #
- 在 hostname 行的行首添加两个空格
  - sed "s/hostname/ &/" prometheus.yml
- 在文件最后一行添加变量中的内容，注意 $a 前加 `\` 符号以便让 sed 认出 $a 表示最后一行
  - sed -i "$a${Masters\[${i}]%%=*} ${Masters\[${i}]##\*=}" /tmp/hosts #
- 在开头是 kind: Deployment 这行的下一行的下一行，添加 namespace: redis 行


    sed -n '/^kind: Deployment/{N;a\  namespace: redis
    p}' all-redis-operator-resources.yaml

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/ec3zxx/1616166346275-89fa46ea-529e-45ab-b8c9-501e0d4f8b43.jpeg)

## 匹配行的下 N 行替换内容

定位到包含字符串 console-agent-exporter 的行下的第 1 行，将 scrape_interval.\* 替换为 scrape_interval: 180s

- sed '/console-agent-exporter/{n;s/scrape_interval.\*/scrape_interval: 180s/;}' prometheus.yml

定位到包含字符串 console-agent-exporter 的行下的第 2 行，将 scrape_timeout:.\* 替换为 scrape_timeout: 180s

- sed '/console-agent-exporter/{n;n;s/scrape_timeout.\*/scrape_timeout: 180s/;}' prometheus.yml

## 其他

- if \[\[ ! `grep 'web.config=' /opt/monitoring/client/docker-compose.yml` ]]; then sed -i '/web.listen-address=:9100/i\ \ \ \ - --web.config=/etc/prometheus/config_out/web-config.yml' /opt/monitoring/client/docker-compose.yml; fi
