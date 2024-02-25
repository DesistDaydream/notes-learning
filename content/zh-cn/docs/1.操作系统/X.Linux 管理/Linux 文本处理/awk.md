---
title: awk
---

# 概述

> 参考：
>
> - [官网](https://www.gnu.org/software/gawk/)
> - [官方文档](https://www.gnu.org/software/gawk/manual/gawkinet/)

awk 以字段为单位进行处理

其实就是把一行的数据分割,然后对每个字段进行处理,包含 cut 等命令)，支持变量，条件判断，循环，数组等编程基本语言

# Syntax(语法)

**awk \[OPTIONS] 'COMMAND' FILE**

**awk \[OPTIONS] 'PATTERN1{ACTION1} PATTERN2{ACTION2}....' FILE**

**OPTIONS**

- **-f FILE** # 指定要使用的 awk 代码文件。
- **-F \[“\[分隔符]”]** # 指定分隔符，默认分隔符为一个或多个的“空格键”或者“tab 键”，也可以具体指定一个或多个
  - e.g.当使用-F “ \[/:]”的时候即是空格、/、:这三个符号出现任意一个都算作一个分隔符

# AWK 语言

awk 其实本质上可以看作编程语言，只不过这个语言只是用来处理文本的而已。在使用命令时，可以使用 `-f` 选项指定要使用的代码文件。

## awk 语言的基本结构

awk 代码由 `PATTERN {ACTION}` 组成，PATTERN 是可省略的。

- PATTERN 用来进行匹配的模式，匹配到的内容将会执行 ACTION 中定义的操作
  - /搜索模式/
  - 判断模式
  - BEGIN 执行 ACTION 前的准备工作，比如给 awk 中的自带变量赋值,在 print 前在屏幕输出点内容
  - END 执行 ACTION 后的收尾工作
- ACTION 用来执行具体的动作
  - print $NUM "输出内容" $NUM...... # 在屏幕输出哪几个字段以及哪些内容，内容可以是各种分隔符

一个最简单的 awk 代码如下：

```bash
{print $0}
```

这里省略的匹配模式，也就是说匹配所有内容。直接使用一个 ACTION，输出文件中的所有内容了。

### 特殊的 BEGIN 与 END 模式

BEGIN 与 END 与正常的模式不同，不用于匹配输入记录。通常用来为本次 awk 的运行提供启动和清理操作。

> BEGIN 模式通常用来进行变量复制。END 模式通常用来清理数据

## Hello World

假设现在有如下文件：

```bash
Hello World Text
```

awk 代码如下：

```bash
BEGIN {
    string = "Hello World"
    printf("%s%s",string,RS)
}

{
    printf("%s%s",$0,RS)
}
```

输出结果：

```bash
# 可以直接执行 BEGIN 模式中的动作
~]# awk -f helloworld.awk
Hello World
^C
# 需要指定待处理文件，才可以执行匹配模式中的动作
~]# awk -f helloworld.awk text
Hello World
Hello World Text
```

## awk 中的变量

引用 awd 中的内置变量时不用 `$` 符号

- **$0** # 全部输入
- **ARGC** #
- **FS** # 字段分隔字符，`默认值：空格`
- **NF** # 每一行拥有的字段总数(当每行字段数不一样，可以 print NF 来打印每行的最后一个字段)
- **NR** # 目前 awk 所处理的是第几行的数据
- **OFS** # 当前输出内容 print 的时候使用的分隔符，默认为空白，print 中以逗号分隔每个字段
- **RS** # 分隔符。`默认值：\n`，即换行符

# 应用示例

去除相同的行

- `awk '!seen[$0]++' test.log`
  - `!seen[$0]++` 的含义是，如果某一行第一次出现， `seen[$0]` 的值为 0，使用 ! 取反的话，值为 1，i.e. 结果为真，所以 awk 会输出该行。然后 `seen[$0]` 会自增 1，对于后续该行的重复，`!seen[$0]` 为假，因此 awk 不会再次输出该行。

输出文本最后一行

  - awk 'END {print}'

从 FILE 文件中，删除每行第一列，输出剩余的
  - awk '{ $1=""; print $0 }' FILE

查找 hcs 的 access 实时日志的带 HIT 字符的行，取出第五段内容然后排序总结，该日志可以实时查看用户的命中情况以及访问的资源
  - `cat accesslog | grep 'HIT' | awk '{print $5}' | sort | uniq -c`

在/etc/passwd 文件中，每行以:为分隔符，打印 username:这几个字符，后面跟以:分割的第一个字段内容，后面跟换行(\n)然后 uid:这几个字符，再跟第三个字段内容,效果如右图所示
  - `awk -F":" '{ print "username:" $1 "\nuid:" $3 }' /etc/passwd`

搜索模式以冒号为分隔符,找到第七个字段以 bash 结尾的所有行,输出每行的第一个和第三个字段(注意：这里面的~在 shell 环境中是用=~表示的)
- `cat /etc/passwd | awk -F : '$7~/bash$/{print $1,$3}'`
  - 还可以写成判断模式，判断第七字段的字符是否等于/bin/bash，
  - `cat /etc/passwd | awk -F : '$7=="/bin/bash"/{print $1,$3}'`

以冒号为分隔符的第三个字段数小于 10 的那些行，输出其中的第一和第三个字段

- `cat /etc/passwd | awk -F : '$3<10{print $1,$3}'`

- `cat /etc/passwd | awk 'BEGIN{FS=":"}{print "UserName\n-----------"}$3<10{print $1,$3}'`

以=号为分隔符，不包含开头带#或者空白行的所有行，显示这些行的第一个字段,判断模式与搜索模式并用
  - `awk -F = '!/^#|^$/{print $1}' /etc/sysctl.conf`
    - `awk -F = '/^[^#]/{print $1}' /etc/sysctl.conf` # 不含 # 的行

查看当前的普通用户个数
  - `cat /etc/passwd | awk -F : '$3>1000 && $7=="/bin/bash"{print $1,$3}'`

其他

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/aokdnm/1616166312890-e4478d40-8a1b-4f08-a348-81ee7c69c9e0.jpeg)

## /proc/softirqs 文件处理示例

由于 /proc/softirqs 文件中第一行与其他行，前面少了一些内容，所有第一行系列要多空出来一些字符

```bash
BEGIN{
  cpucount = 6
}
# 处理第一行
NR == 1{
  num = 5;
  # 由于第一行少了一列，所以让第一列与第二行的第二列对其
  printf("%30s",$1);
  # 从第二列开始循环，每隔15个字符便输出一列
  for (i=2;i<=cpucount;i++) {
    printf("%15s\t",$i);
  }
  printf(RS);
}
# 处理第二行及以后的行
NR > 1{
  # 通过循环，每隔15个字符便输出一列。输出完成后换行。
  for (i=1;i<=cpucount+1;i++) {
    printf("%15s\t",$i);
  }
  printf(RS);
}
```

## 聚合一个文件中指定字段的数字，求和

- 文件中的内容如下

```bash
[root@dengrui test_dir]# cat strace.file
31090 pread64(227, ""..., 16384, 81920) = 16384 <0.000991>
31090 pread64(227, ""..., 16384, 65536) = 16384 <0.001292>
31090 pread64(227, ""..., 16384, 98304) = 16384 <0.000176>
31090 pread64(129, ""..., 16384, 131072) = 16384 <0.002121>
31090 pread64(129, ""..., 16384, 16384) = 16384 <0.000932>
31090 pread64(128, ""..., 16384, 49152) = 16384 <0.001072>
31090 pread64(128, ""..., 16384, 16384) = 16384 <0.000820>
```

- 通过 awk 命令，聚合系统调用中的第三个参数

```bash
awk -F, '{print $3}' test.file | awk '{sum += $1} END {print sum}'
```
