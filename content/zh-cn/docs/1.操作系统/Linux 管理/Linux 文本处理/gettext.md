---
title: "gettext"
linkTitle: "gettext"
date: "2023-06-05T10:58"
weight: 20
---

# 概述

> 参考：
> 
> - [官网](https://www.gnu.org/software/gettext/)

gettext 是 GNU 旗下的一组工具集合，提供了一个框架以帮助其他 GNU 包生成多语言消息。

# 安装 gettext

一般的发行版都默认自带 gettext 工具集，如果没有，使用包管理器安装 gettext 包即可，通常包含如下几个工具：

- envsubst
- gettext
- gettext.sh
- ngettext

Ubuntu

- apt install gettext-base

# envsubst

> 参考：
> 
> - [官方手册，envsubst](https://www.gnu.org/software/gettext/manual/html_node/envsubst-Invocation.html)

envsubst 程序可以用来替换环境变量的值。正常情况下，与 cat 命令类似，所有的标准输入都会复制到标准输出，但是不同的地方在于，如果标准输入中包含变量引用，比如 `$VARIABLE` 或 `${VARIABLE}` 这种形式，则这些引用将会被替换为变量的值：

```bash
~]# envsubst
HOME
HOME
$HOME
/root
${HOME}   
/root
```

若我们将标准输入改为由文件提供，那么我们就可以将文件中的所有变量引用的地方都替换为对应的值。比如：

```bash
sudo tee ~/tmp/test.txt <<-"EOF"
HOME = ${HOME}
PATH = ${PATH}
API_URL = ${API_URL}
EOF
```

执行 `envsubst < test.txt > test2.txt` 命令以替换文件中的环境变量，生成的 test.txt 内容如下：

```bash
~]# cat test2.txt 
HOME = /root
PATH = /usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin:/usr/games:/usr/local/games:/snap/bin:/usr/local/go/bin
API_URL = 
```

也可以使用 `envsubst '${HOME}' < test.txt > test1.txt` 这种方式只替换 `${HOME}` 变量，多个变量以 `,` 分隔。