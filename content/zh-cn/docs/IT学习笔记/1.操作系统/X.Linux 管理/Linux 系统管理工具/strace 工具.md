---
title: strace 工具
slug: strace
---

# 概述

> 参考：
>
> - [GitHub 项目，strace/strace](https://github.com/strace/strace)
> - [官网](https://strace.io/)
> - [Manual(手册),strace(1)](https://man7.org/linux/man-pages/man1/strace.1.html)

**strace** 是一个用来跟踪 **system calls(系统调用)** 和 **signals(信号)** 的工具。

在最简单的情况下，strace 会运行指定的命令，直到退出为止。 它截获并记录 **由进程调用的系统调用** 和 **由进程接收的信号**。 每个系统调用的名称，其参数及其返回值都会在 标准错误 或使用 -o 选项指定的文件上显示。

strace 是有用的诊断，说明和调试工具。系统管理员，诊断人员和疑难解答人员将发现，对于不容易获取源代码的程序而言，这是无价的，因为它们无需重新编译即可跟踪它们。学生，黑客和过于好奇的人会发现，即使跟踪普通程序，也可以学到很多有关系统及其系统调用的知识。而且程序员会发现，由于系统调用和信号是在用户/内核界面上发生的事件，因此仔细检查此边界对于错误隔离，健全性检查和尝试捕获竞争状况非常有用。

## strace 输出内容介绍

### 追踪到系统调用时输出的信息

下面是一个最基本，最简单的追踪，strace 程序执行时，输出的每一行内容都是一个 syscall(系统调用)。基本格式如下：

`SyscallName(Parameter) = ReturnValue`

假如我追踪 cat /dev/null 命令，则输出中有这么一段：

```bash
# 使用了 openat 这个系统调用，参数为 "/dev/null",O_RDONLY，返回值为 3
openat(AT_FDCWD, "/dev/null", O_RDONLY) = 3
```

如果命令出现错误，通常 ReturenVale 为 -1，并附加 errno 符号和错误信息

```bash
openat(AT_FDCWD, "123", O_RDONLY)       = -1 ENOENT (No such file or directory)
```

### 追踪到信号时输出的信息

如果追踪到信号，则输出内容的基本格式如下：

`SignalName{si_signo=SignalName, si_code=SI_USER, si_pid=PID, ...}`

假如我同时最总两个进程，并像其中一个进程发送了 SIGTERM 信号，则输出中有这么一段：

```bash
[pid  5314] --- SIGTERM {si_signo=SIGTERM, si_code=SI_USER, si_pid=27467, si_uid=0} ---
```

# starce 语法

strace \[-ACdffhikqrtttTvVxxy] \[-I n] \[-b execve] \[-e expr]... \[-a column] \[-o file] \[-s strsize] \[-X format] \[-P path]... \[-p pid]... { -p pid | \[-D] \[-E var\[=val]]... \[-u username] COMMAND \[args] }

strace -c \[-df] \[-I n] \[-b execve] \[-e expr]... \[-O overhead] \[-S sortby] \[-P path]... \[-p pid]... { -p pid | \[-D] \[-E var\[=val]]... \[-u username] COMMAND \[args] }

## OPTIONS

### General 通用选项

**-e EXPR** # 用来指定要追踪的事件以及如何追踪。

EXPR(表达式) 的格式为 `QUALIFIER=[!]VALUE[,VALUE...]`

- **QUALIFIER(限定词)** # 可用的限定词有 trace、abbrev、verbose、raw、signal、read、write、fault、inject、status、quiet、decode-fds、kvm。`默认值：trace`。
- **VALUE** # 是与 qualifier 相关的字符串或数字。

Note：

- 由于 QUALIFIER 的默认值为 trace。所以 -e trace=sendto,read 也可以写成 -e sendto,read。
- QUALIFIER 限定词根据功能不通，在 filtering、tampering、Output format 等选项中，有具体的使用说明。
- 使用 `!` 会否定该组值。比如，-e trace=open 表示仅追踪 open 系统调用；而 -e trace='!open' 表示追踪除了 open 以外的所有系统调用
  - 注意加单引号，否则无法识别，并报错提示：`-bash: !XXXX: event not found`

### Startup 启动选项

- **-p PID** # 追踪指定 PID 的进程的系统调用。

### Tracing 跟踪选项

- **-f** # 跟踪子进程，并显示 PID 号。些子进程是由 fork(2)，vfork(2) 和 clone(2) 系统调用而由当前跟踪的进程创建的。

### Filtering 过滤选项

- **-e trace=SYSCALL_SET** # 指定要追踪的系统调用。
  - 可用的 SYSCALL_SET 有如下这些
    - **SYSCALL1\[,SYSCALL2,...]** # 直接指定系统调用的名称，多个名称以逗号分隔。
    - **/REGEX** # 前面加上 `/` ，后面可以使用正则表达式进行匹配，来匹配系统调用的名称。
    - **%SyscallSet** # 前面加上 `%`，就会追踪一类系统调用的集合。比如：
      - %clock    Trace all system calls that read or modify system clocks.
      - %creds    Trace all system calls that read or modify user and group identifiers or capability sets.
      - %desc     Trace all file descriptor related system calls.
      - **%file** # 追踪所有以文件名为参数的系统调用。可以看作是 -e trace=open,stat,chmod,unlink,..... 的简写。
      - %fstat    Trace fstat and fstatat syscall variants.
      - %fstatfs  Trace fstatfs, fstatfs64, fstatvfs, osf_fstatfs, and osf_fstatfs64 system calls.
      - %ipc      Trace all IPC related system calls.
      - %lstat    Trace lstat syscall variants.
      - %memory   Trace all memory mapping related system calls.
      - %network  Trace all the network related system calls.
      - %process  Trace all system calls which involve process management.
      - %pure     Trace syscalls that always succeed and have no arguments.
      - **%signal** # 追踪所有与信号相关的系统调用。
      - %stat     Trace stat syscall variants.
      - %statfs   Trace statfs, statfs64, statvfs, osf_statfs, and osf_statfs64 system calls.
      - %%stat    Trace syscalls used for requesting file status.
      - %%statfs  Trace syscalls related to file system statistic
    - ........等等
- **-e signal=SIGNAL_SET** # 指定要追踪的信号。
- **-e status=STATUS_SET** # 指定要追踪的系统调用的返回码

### Output format 输出格式选项

- **-a COLUMN** # 设定列的间隔为 COLUMN，默认为 40。i.e. `=` 与前面的间隔
- **-o, --output \<FILE>** # 将追踪结果输出到文件中(默认标准错误)。
  - 与 -ff 参数一起使用时，会把每个线程的追踪写到单独的文件中，以 FileName.PID 格式命名。
- **-q, --quiet=STRING** # 抑制有关附加、分离、个性的消息。当 strace 的输出被重定向到文件中时，会自动添加该选项。
  - 可用的值有：attach,personality,exit,all。这些可用的值只在 --quiet 选项时可用，我们还可以使用 -q、-qq、-qqq 以添加不同的抑制信息，q 越多，抑制的信息就越多。
- **-s, --string-limit \<STRSIZE>** # 设定要输出的最大字符串长度为 STRSIZE。`默认值：32`。Note:文件名不作为字符串，并始终完整打印。
  - 示例如下，在 sendto 和 read 系统调用中，参数只显示了 32 个字符。当指定 -s 选项后，可以输出更多字符。

```bash
~]# strace -p 22863 -e trace=sendto,read
sendto(6, "GET / HTTP/1.0\r\nUser-Agent: Keep"..., 71, 0, NULL, 0) = 71
read(6, "HTTP/1.1 426 Upgrade Required\r\nd"..., 4096) = 129
~]# strace -p 22863 -s 1000 -e trace=sendto,read
strace: Process 22863 attached
sendto(7, "GET / HTTP/1.0\r\nUser-Agent: KeepAliveClient\r\nHost: 10.0.9.213:50080\r\n\r\n", 71, 0, NULL, 0) = 71
read(7, "HTTP/1.1 426 Upgrade Required\r\ndate: Fri, 24 Jul 2020 07:53:01 GMT\r\nserver: istio-envoy\r\nconnection: close\r\ncontent-length: 0\r\n\r\n", 4096) = 129
```

- **-t, -tt, -ttt** # 显示追踪时间(在输出的行开头显示)。2 个 t 显示微秒，3 个 t 显示时间戳
- **-T** # 显示追踪花费的时间(在输出的行末尾显示)
- **-y, -yy** # 打印与文件描述符参数相关联的路径。2 个 y，打印与套接字文件描述符相关的特定协议信息，以及与设备文件描述符相关的块/字符设备号。

```bash
# 这是一个建立 http 连接的系统调用追踪
# 不加 -y，只显示数字 3，表示当前文件描述符的编号为3
[pid  8675] write(3, "POST /api/v1/auth/tokens:login HTTP/1.1\r\nHost: 10.20.5.98:8056\r\nUser-Agent: Go-http-client/1.1\r\nContent-Length: 44\r\nContent-Type: application/json\r\nAccept-Encoding: gzip\r\n\r\n{\"auth\":{\"name\":\"admin\",\"password\":\"admin\"}}", 217) = 217
# 使用 -y 参数，显示编号为3的文件描述符的路径
[pid  8667] write(3<socket:[80219]>, "POST /api/v1/auth/tokens:login HTTP/1.1\r\nHost: 10.20.5.98:8056\r\nUser-Agent: Go-http-client/1.1\r\nContent-Length: 44\r\nContent-Type: application/json\r\nAccept-Encoding: gzip\r\n\r\n{\"auth\":{\"name\":\"admin\",\"password\":\"admin\"}}", 217) = 217
# 使用 -yy 参数，显示3号文件描述符的Socket的具体信息(源和目的地址)
[pid  8623] write(3<TCP:[172.38.40.250:27436->10.20.5.98:8056]>, "POST /api/v1/auth/tokens:login HTTP/1.1\r\nHost: 10.20.5.98:8056\r\nUser-Agent: Go-http-client/1.1\r\nContent-Length: 44\r\nContent-Type: application/json\r\nAccept-Encoding: gzip\r\n\r\n{\"auth\":{\"name\":\"admin\",\"password\":\"admin\"}}", 217) = 217
# 查看一下这个进程的3号文件描述符
[root@ansible fd]# ll /proc/8685/fd/3
lrwx------ 1 root root 64 Jan 24 10:55 /proc/8675/fd/3 -> 'socket:[80219]'
[root@ansible fdinfo]# cat /proc/net/tcp
# 从这里查看socket号为80219的连接信息，16进制转换过去就是 172.38.40.250:27436 与 10.20.5.98:8056
```

### Statistics 统计选项

- **-c** # 统计每一次系统调用的执行时间、次数、错误次数。输出效果如下：
  - -c 参数常用来在排障之前，查看当前进程使用了哪些系统调用，然后在后续排障中单独追踪指定的系统调用

```bash
[root@dr-02 keepalived]# strace -p 22863 -c
strace: Process 22863 attached
% time     seconds  usecs/call     calls    errors syscall
------ ----------- ----------- --------- --------- ----------------
 39.04    0.002437          32        76           select
 16.42    0.001025           6       147           fcntl
 14.48    0.000904          69        13        13 connect
  9.55    0.000596          59        10           close
  7.59    0.000474          36        13           sendto
  4.37    0.000273          11        23           read
  3.57    0.000223           8        26           getsockopt
  3.22    0.000201          15        13           socket
  1.75    0.000109           8        13           setsockopt
------ ----------- ----------- --------- --------- ----------------
100.00    0.006242                   334        13 total
```

- **-o FILE** # 将追踪的输出写入 FILE。 如果提供了-ff 选项，则使用 filename.pid 格式。 如果参数以“ |”开头 或“！”，则将其余参数视为命令，并将所有输出通过管道传递给它。 这对于将调试输出管道传输到程序而不影响已执行程序的重定向非常方便。 后者当前与 -ff 选项不兼容。

### Tampering 干预选项

- **-e inject=SYSCALL_SET** #

### Miscellaneous 选项

# 应用示例

## 设置 strace 命令的运行时间

`strace` 命令本身并不支持指定运行最大时间的选项。但是，你可以使用 `timeout` 命令来限制 `strace` 命令的运行时间。`timeout` 命令可以在指定的时间内运行一个命令，并在超时时终止该命令的执行。

例如，如果你想在 `strace` 命令在 5 秒内运行，你可以使用以下命令：

`timeout 5 strace -p 123456`

这将运行 `<your command>` 命令并使用 `strace` 进行跟踪，但最多只运行 5 秒。如果命令在 5 秒内完成，则 `timeout` 命令将返回该命令的退出状态码。否则，`timeout` 命令将终止该命令并返回一个非零的退出状态码。

## 其他

- 追踪 ls 命令的系统调用情况。
  - **starce ls**
- 统计 df 命令的系统调用信息。
  - **strace -c df**
- 追踪 22863 进程的系统调用，只追踪网络与 read 相关的系统调用。输出更多信息，扩大输出字符串到 1000。
  - **strace -p 22863 -s 1000 -e trace=%network,read**
- 追踪新编译的 main 程序，显示时间、追踪线程、扩大输出字符、追踪 write()、追踪 SIGHUP 信号
  - **strace -t -f -s 1000 -e trace=write -e signal='SIGHUP' ./main --xsky-pass=admin**

**分析进程 I/O 情况**

- 追踪 1234 进程及其子进程，去掉所有字符串，在末尾显示花费的时间，将结果保存到 strace.file 文件中
  - **strace -T -s 0 -f -p 1234 -o strace.file**
- 文件中的内容如下

```bash
~]# cat strace.file
31090 pread64(227, ""..., 16384, 81920) = 16384 <0.000991>
31090 pread64(227, ""..., 16384, 65536) = 16384 <0.001292>
31090 pread64(227, ""..., 16384, 98304) = 16384 <0.000176>
31090 pread64(129, ""..., 16384, 131072) = 16384 <0.002121>
31090 pread64(129, ""..., 16384, 16384) = 16384 <0.000932>
31090 pread64(128, ""..., 16384, 49152) = 16384 <0.001072>
31090 pread64(128, ""..., 16384, 16384) = 16384 <0.000820>
```

- 通过 awk 命令，聚合系统调用中的第三个参数，即可得出来追踪这一段时间，这个程序从磁盘中读取了多少数据

```bash
awk -F, '{print $3}' test.file | awk '{sum += $1} END {print sum}'
```
