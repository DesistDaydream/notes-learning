---
title: lsof 列出打开的文件
---

# 概述

> 参考：
> - [Manual(手册)](<https://man.cx/lsof(8)>)

该工具以文件为主体，用于列出打开文件的进程，进程打开的端口(TCP、UDP)等、找回/恢复删除的文件。是十分方便的系统监视工具，因为 lsof 需要访问核心内存和各种文件，所以需要 root 用户执行。

```bash
[root@master0 ~]# lsof | more
COMMAND    PID  TID           USER   FD      TYPE             DEVICE  SIZE/OFF       NODE NAME
systemd      1                root  cwd       DIR              253,0       238         64 /
systemd      1                root  rtd       DIR              253,0       238         64 /
systemd      1                root  txt       REG              253,0   1612152   17149941 /usr/lib/systemd/systemd
......
kthreadd     2                root  cwd       DIR              253,0       238         64 /
kthreadd     2                root  rtd       DIR              253,0       238         64 /
kthreadd     2                root  txt   unknown                                         /proc/2/exe
......
lsof      1893                root  mem       REG              253,0    155784      72860 /usr/lib64/libselinux.so.1
lsof      1893                root  mem       REG              253,0    164240      41015 /usr/lib64/ld-2.17.so
lsof      1893                root    4r     FIFO                0,9       0t0      37707 pipe
lsof      1893                root    7w     FIFO                0,9       0t0      37708 pipe
```

每行显示一个打开的文件，若不指定条件默认将显示所有进程打开的所有文件。lsof 输出各列信息的意义如下：

COMMAND 进程的名称 | PID 进程标识符 | USER 进程所有者 | FD 文件描述符 | TYPE 文件类型 | DEVICE 磁盘的名称 | SIZE 文件的大小 | NODE 索引节点 | NAME 文件的绝对路径

## FD # 表示该文件被打开的 FD 号或其他信息

- cwd：表示 current work dirctory，即：应用程序的当前工作目录，这是该应用程序启动的目录，除非它本身对这个目录进行更改
- txt # 该类型的文件是程序代码，如应用程序二进制文件本身或共享库，如上列表中显示的 /sbin/init 程序
- er # FD 信息错误(参考 NAME 列);
- ltx：shared library text (code and data);
- mxx ：hex memory-mapped type number xx.
- mem # 内存映射文件
- mmap # memory-mapped device;
- pd # 父目录
- rtd # root 目录
- v86 VP/ix mapped file;
- 0：表示标准输出
- 1：表示标准输入
- 2：表示标准错误

> 一般在标准输出、标准错误、标准输入后还跟着文件状态模式：r、w、u 等

- u：表示该文件被打开并处于读取/写入模式
- r：表示该文件被打开并处于只读模式
- w：表示该文件被打开并处于
- 空格 # 表示该文件的状态模式为 unknow，且没有锁定
- - # 表示该文件的状态模式为 unknow，且被锁定

同时在文件状态模式后面，还跟着相关的锁

- N：for a Solaris NFS lock of unknown type;
- r：for read lock on part of the file;
- R：for a read lock on the entire file;
- w：for a write lock on part of the file;（文件的部分写锁）
- W：for a write lock on the entire file;（整个文件的写锁）
- u：for a read and write lock of any length;
- U：for a lock of unknown type;
- x：for an SCO OpenServer Xenix lock on part of the file;
- X：for an SCO OpenServer Xenix lock on the entire file;
- space：if there is no lock.

## TYPE # 文件类型

- DIR：表示目录
- CHR：表示字符类型
- BLK：块设备类型
- UNIX： UNIX 域套接字
- FIFO：先进先出 (FIFO) 队列
- IPv4：网际协议 (IP) 套接字
- .......

# Syntax(语法)

**lsof \[ -?abChKlnNOPRtUvVX ] \[ -A A ] \[ -c c ] \[ +c c ] \[ +|-d d ] \[ +|-D D ] \[ +|-e s ] \[ +|-f \[cfgGn] ] \[ -F \[f] ] \[ -g \[s] ] \[ -i \[i] ] \[ -k k] \[ +|-L \[l] ] \[ +|-m m ] \[ +|-M ] \[ -o \[o] ] \[ -p s ] \[ +|-r \[t\[m<fmt>]] ] \[ -s \[p:s] ] \[ -S \[t] ] \[ -T \[t] ] \[ -u s ] \[ +|-w ] \[ -x \[fl] ] \[ -z \[z] ] \[ -Z \[Z] ] \[ -- ] \[names]**

## OPTIONS

- **-a** # 过滤选项之间进行 AND 运算。比如我使用 -d 和 -p，则结果要两个筛选都满足才可以。默认情况是 或 运算。列出满足任意过滤选项的所有结果。
  - 说白了，这 -a 选项就是个逻辑运算符。
- **-c <STRING>** # 列出以 STRING 字符开头的命令的进程的文件列表。其实就是通过进程名筛选
- **+d <DIR>** # 列出目录下被打开的文件
- **-d <FD>** # 列出占用指定文件描述符的进程。可以使用 2-10 这种方式来列出 2 到 10 号描述符的文件。
- **+D <DIR>** # 递归列出目录下被打开的文件
- **-g <GroupID>** # 列出 GID 号进程详情
- **-i \[<STRING>]** # 列出符合条件的网络连接相关。（4、6、协议、:PORT、 @IP ）
- **-n** # 直接显示 IP 而不是主机名
- **-N <DIR>** # 列出使用 NFS 的文件
- **-p \<PID\[,PID,PID....]>** # 列出指定进程号所打开的文件。多个 PID 以逗号分隔
- **-P** # 直接显示端口号，而不是端口号的名称
- **-u <USERNAME>** # 列出指定用户所打开的文件
- **-w** # 关闭程序运行中产生的警告信息。

## EXAMPLE

- 查看谁正在使用某个文件，也就是说查找某个文件相关的进程
  - **lsof /bin/bash **
- 显示除了 root 用户下的 sshd 进程所用的文件
  - **lsof -u ^root -c sshd**
- 列出目前连接主机 peida.linux 上端口为：20，21，22，25，53，80 相关的所有文件信息，且每隔 3 秒不断的执行 lsof 指令
  - **lsof -i @peida.linux:20,21,22,25,53,80 -r 3**
- 列出所有的网络连接
  - **lsof -i**
- 列出所有 tcp 网络连接信息
  - lsof -i tcp
- 列出正在使用 3306 端口的进程信息
  - lsof -i :3306
- 列出 9267 号进程打开的文件描述符为 132 的文件
  - lsof -p 9267 -d 132 -a
- 列出 9267 号进程打开的所有文件，以及文件描述符为 132 的所有文件
  - lsof -p 9267 -d 13
- 列出谁在使用某个特定的 udp 端口
  - lsof -i udp:55
- 列出某个用户的所有活跃的网络端口
  - lsof -a -u test -i
