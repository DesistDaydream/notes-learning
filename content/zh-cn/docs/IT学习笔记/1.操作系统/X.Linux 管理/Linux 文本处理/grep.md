---
title: grep
---

# 概述

> 参考：[man 手册](https://man.cx/grep)

grep 是文本搜索工具，可以使用正则表达式在文件内查找符合条件的字串行

# Syntax(语法)

**grep \[OPTIONS] PATTERNS \[FILE...]**

grep 根据 PATTERNS(模式) 过滤给定的内容。其实就是使用正则表达式，过滤内容。

## OPTIONS：

### Pattern Syntaz(模式语法)

用于定义过滤时所解析的正则表达式语法

- **-E, --extended-regexp** # 将 PATTERNS 解释为扩展的正则表达式（ERE，请参见下文）。
- **-P,--perl-regexp** # 将 PATTERNS 解释为与 Perl 兼容的正则表达式(PCREs)。与-z（--null-data）选项结合使用时，此选项是实验性的，并且 grep -P 可能会警告未实现的功能。

### Matching Control(配置控制)

- **-i **# 忽略大小写
- **-v, --invert-match **# 反向选择，选择没有要匹配的字符串的行

### General Output Control(通用输出控制) 选项

- **-c, --count** # 计算找到的符合行的次数
- **--color=auto** # 把查找到内容自动表上默认颜色，auto 可改成别的颜色英文
- **-l,--files-with-matches** # 在筛选时，只输出文件名。常用来在一堆文件中，筛选指定内容，只查看哪些文件有指定的内容。
- **-o, --only-matching** # 仅显示被匹配到的字符串，不显示整行
- **-s, --no-messages** # 不显示错误信息

### Output Line Prefix Control(控制输出行的前缀)

- **-n** # 顺便输出行号

### Context Line Control(控制输出内容的上下行)

- **-A NUM, --after-context=NUM** # 打印出查找到的行的下面 NUM 行
- **-B NUM, --before-context=NUM** # 打印出查找到的行的上面 NUM 行
- **-C NUM, --context=NUM** # -A 与 -B 选项的结合体，控制输出内容的 上面 和下面 NUM 行

### File and Directory Selection(文件和目录选择)

- **-a, --text** # 像对待文本一样处理二进制文件；这等效于--binary-files = text 选项。
- **-R, --dereference-recursive** # 递归地阅读每个目录下的所有文件并进行 grep 操作;该选项相当于-d recurse

## EXAMPLE

- 查看 accesslog 文件的实时更新，并筛选出不包含两个字符串的行
  - tailf accesslog | grep -vE '(miguvideo|mgtv)'
- grep --color=auto -i R.\*h ./boot.log | grep -Evi "star|net" #不区分大小写搜索 boot.log 文件中包含 Rh 中间含有任意字符的，并且不包含 Star 或 net 的所有行，并以高亮颜色显示搜索的字符串，|表示或的关系，正则表达式扩展内容，选项必须有 E 才能生效。
- grep -i '/bin/bash' /etc/passwd | sort -n -t: -k3 | tail -1|cut -d: -f1 #取出默认 shell 为 bash 且其 ID 号最大的用户
- grep "^#\[\[:space:]]{1,}\[^\[:space:]]{1,}" /etc/X #取出/etc/x 文件中井号开头后最少一个空白字符后最少一个非空白字符的行
- ifconfig | egrep --color=auto -n '\[0-9]{1,3}.\[0-9]{1,3}.\[0-9]{1,3}.\[0-9]{1,3}'
- egrep --color=auto -n '\[0-9]+.\[0-9]+.\[0-9]+.\[0-9]+' #匹配 ifconfig 中的所有 IP 地址，带匹配高亮，带行号（egrep 就是 grep -E）
- egrep --color=auto -n '<\[0-9]{2}>.\[0-9]+.\[0-9]+.\[0-9]+' #可以搜索第一段是两位数的 IP，比如 10.0.0.0 网段
- grep -i --color=auto '\[0-9]+.\[0-9]+.\[0-9]+.\[0-9]+' ./interfaces #不适用 egrep 的方法

应用示例

## 筛选 `{{ }}` 之间的内容

    hi,hello {{A1}}
    {{B0B}}test{{CC_CC}}
    @{{D-DD}}
    {{E#@EEE}}

### 筛选后文本&#xA;

**cat content.txt | grep -oP "(?<={{)(\w|-|#|@)+(?=}})"**

    A1
    B0B
    CC_CC
    D-DD
    E#@EEE

**cat content.txt | grep -oE "{{(\w|-|#|@)+}}"**

**cat content.txt | grep -oP "{{(\w|-|#|@)+}}"**

    {{A1}}
    {{B0B}}
    {{CC_CC}}
    {{D-DD}}
    {{E#@EEE}}
