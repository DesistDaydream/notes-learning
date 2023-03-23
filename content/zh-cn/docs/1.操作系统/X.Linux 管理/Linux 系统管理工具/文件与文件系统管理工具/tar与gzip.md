---
title: tar与gzip
---

# tar

> 参考：
>
> - [Manual(手册),tar(1)](https://man7.org/linux/man-pages/man1/tar.1.html)

tar 是一个归档工具

## Syntax(语法)

**tar \[OPTIONS] /PATH/FILE**
OPTIONS:

- **-f**# 指定被处理的文件（在所有选项里一定要放最后一个，否则会报错）
- **-x** # 解包文件
- **-c** # 打包文件
- **-C, --directory=\<DIR>**# 解包到指定目录
- **-v**# 显示执行命令时的详细过程
- **-t** # 查看已经打包的文件中的内容
- **-z** # 通过 gzip 指令处理打包的文件；
- **--strip-components=NUM** # 去除前缀目录，i.e.默认会自动创建与压缩包同名的目录来存放压缩包内的文件，当 NUM 为 1 时，则不再创建该目录，直接将压缩包内的文件全部解压到当前目录

EXAMPLE:

- tar -zcvf xxx.tar.gz /tmp/test/\* # 把/tmp/test/下所有的文件和目录都打包并压缩成 xxx.tar.gz
- tar -zxvf xxx.tar.gz \[-C /XXX] # 解压 xxx.tar.gz 文件（加-C /tmp/test2/ 指定解压缩的路径，解压缩到/tmp/test2/目录下）
- tar -tvf xxx.tar.gz # 不解包查看压缩文件里有哪些文件和目录
- zcat xxx.tar.gz # 读取打包文件的内容等详细信息

# gzip

命令用来压缩文件。gzip 是个使用广泛的压缩程序，文件经它压缩过后，其名称后面会多处“.gz”扩展名。

gzip 是在 Linux 系统中经常使用的一个对文件进行压缩和解压缩的命令，既方便又好用。gzip 不仅可以用来压缩大的、较少使用的文件以节省磁盘空间，还可以和 tar 命令一起构成 Linux 操作系统中比较流行的压缩文件格式。据统计，gzip 命令对文本文件有 60%～ 70%的压缩率。减少文件大小有两个明显的好处，一是可以减少存储空间，二是通过网络传输文件时，可以减少传输的时间。

## Syntax(语法)

**gzip \[OPTIONS] /PATH/FILE**

**OPTINOS**

- **-l** # 列出压缩文件的详细信息
- **-d** # 解开压缩文件
- **-n** # 压缩文件时，不保存原来的文件名称及时间戳记；
- **-N** # 压缩文件时，保存原来的文件名称及时间戳记
- **-r**# 递归处理，将指定目录下的所有文件及子目录一并处理，该处理方式是把文件夹下的每一个文件都压缩成一个新的文件，不是对目录整体进行，注意与 tar 的区别
- **-v** # 显示指令执行过程
- -<压缩效率>：压缩效率是一个介于 1~9 的数值，默认值为 6，指定越大的数值，压缩效率就会越高

**EXAMPLE**
![](https://notes-learning.oss-cn-beijing.aliyuncs.com/xft64c/1616166289989-3b7ed972-966f-410f-a3e8-78d2e6836881.jpeg)

- gzip -r .backup # 递归压缩.backup 目录下的所有文件以及子目录的文件，效果如右图所示，压缩只能对单个文件压缩，注意与 tar 打包命令的区别

# 常见问题

当我们在 Linux 中使用绝对路径归档文件时，会出现如下提示：

```bash
$ tar -zcvf test.tar.gz /mnt/d/Projects/DesistDaydream/go-learning/test_files
tar: Removing leading `/' from member names
```

tar 工具会提示我们将开头的 `/` 移除，这是出于安全考虑，如果我们使用绝对路径归档文件，当我们提取这些文件时，就会覆盖原本的文件，这有可能会产生很严重的影响

# 通过其他工具查看 tar.gz 中的内容

得益于 Linux 社区，有很多命令行工具可以来达成上面的目标。下面就让我们来看看使用它们的一些示例。

## 使用 vim 工具

vim 不只是一个编辑器，使用它我们可以干很多事情。下面的命令展示的是在没有解压的情况下使用 vim 查看一个压缩的归档文件的内容：

```bash
" tar.vim version v29
" Browsing tarfile /root/projects/keepalived-ehualu-2.0.20.tar.gz
" Select a file with cursor and press ENTER

keepalived/
keepalived/config/
keepalived/config/check_ports.sh
keepalived/config/keepalived.conf
keepalived/docker-compose-backup.yaml
keepalived/docker-compose-master.yaml
keepalived/README.md
```

无需解压如何查看一个归档或压缩文件的内容无需解压如何查看一个归档或压缩文件的内容

你甚至还可以浏览归档文件的内容，打开其中的文本文件（假如有的话）。要打开一个文本文件，只需要用方向键将鼠标的游标放置到文件的前面，然后敲 ENTER 键来打开它。

## 使用 tar 命令

为了列出一个 tar 归档文件的内容，可以运行：

$ tar -tf ostechnix.tar
ostechnix/
ostechnix/image.jpg
ostechnix/file.pdf
ostechnix/song.mp3

或者使用-v 选项来查看归档文件的具体属性，例如它的文件所有者、属组、创建日期等等。
$ tar -tvf ostechnix.tar
drwxr-xr-x sk/users 0 2018-07-02 19:30 ostechnix/
-rw-r--r-- sk/users 53632 2018-06-29 15:57 ostechnix/image.jpg
-rw-r--r-- sk/users 156831 2018-06-04 12:37 ostechnix/file.pdf
-rw-r--r-- sk/users 9702219 2018-04-25 20:35 ostechnix/song.mp3

## 使用 rar 命令

要查看一个 rar 文件的内容，只需要执行：
$ rar v ostechnix.rar
RAR 5.60 Copyright (c) 1993-2018 Alexander Roshal 24 Jun 2018
Trial version Type 'rar -?' for help
Archive: ostechnix.rar
Details: RAR 5
Attributes Size Packed Ratio Date Time Checksum Name
----------- --------- -------- ----- ---------- ----- -------- ----
rw-r--r-- 53632 52166 97% 2018-06-29 15:57 70260AC4 ostechnix/image.jpg
-rw-r--r-- 156831 139094 88% 2018-06-04 12:37 C66C545E ostechnix/file.pdf
-rw-r--r-- 9702219 9658527 99% 2018-04-25 20:35 DD875AC4 ostechnix/song.mp3
---------- --------- -------- ----- ---------- ----- -------- ----
9912682 9849787 99% 3

## 使用 unrar 命令

你也可以使用带有 l 选项的 unrar 来做到与上面相同的事情，展示如下：
$ unrar l ostechnix.rar
UNRAR 5.60 freeware Copyright (c) 1993-2018 Alexander Roshal
Archive: ostechnix.rar
Details: RAR 5
Attributes Size Date Time Name
----------- --------- ---------- ----- ----
rw-r--r-- 53632 2018-06-29 15:57 ostechnix/image.jpg
-rw-r--r-- 156831 2018-06-04 12:37 ostechnix/file.pdf
-rw-r--r-- 9702219 2018-04-25 20:35 ostechnix/song.mp3
----------- --------- ---------- ----- ----
9912682 3

## 使用 zip 命令

为了查看一个 zip 文件的内容而无需解压它，可以使用下面的 zip 命令：
$ zip -sf ostechnix.zip
Archive contains:
Life advices.jpg
Total 1 entries (597219 bytes)

## 使用 unzip 命令

你也可以像下面这样使用 -l 选项的 unzip 命令来呈现一个 zip 文件的内容：
$ unzip -l ostechnix.zip
Archive: ostechnix.zip
Length Date Time Name
--------- ---------- ----- ----
597219 2018-04-09 12:48 Life advices.jpg
-------- -------
597219 1 file

## 使用 zipinfo 命令

$ zipinfo ostechnix.zip
Archive: ostechnix.zip
Zip file size: 584859 bytes, number of entries: 1
-rw-r--r-- 6.3 unx 597219 bx defN 18-Apr-09 12:48 Life advices.jpg
1 file, 597219 bytes uncompressed, 584693 bytes compressed: 2.1%

如你所见，上面的命令展示了一个 zip 文件的内容、它的权限、创建日期和压缩百分比等等信息。

## 使用 zcat 命令

要一个压缩的归档文件的内容而不解压它，使用 zcat 命令，我们可以得到：
$ zcat ostechnix.tar.gz
zcat 和 gunzip -c 命令相同。所以你可以使用下面的命令来查看归档或者压缩文件的内容：
$ gunzip -c ostechnix.tar.gz

## 使用 zless 命令

要使用 zless 命令来查看一个归档或者压缩文件的内容，只需：
$ zless ostechnix.tar.gz

这个命令类似于 less 命令，它将一页一页地展示其输出。

## 使用 less 命令

可能你已经知道 less 命令可以打开文件来交互式地阅读它，并且它支持滚动和搜索。

运行下面的命令来使用 less 命令查看一个归档或者压缩文件的内容：

$ less ostechnix.tar.gz
