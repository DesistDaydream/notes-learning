---
title: pt-ioprofile
---

## 概述

> 参考：
> - [官方文档](https://www.percona.com/doc/percona-toolkit/LATEST/pt-ioprofile.html#environment)

pt-ioprofile 是 [Percona ](https://www.percona.com/)出的 IO 查看工具。Percona 用来监视进程 IO 并打印文件和 I/O 活动表。

**pt-ioprofile** 本质上就是一个 shell 脚本，只做两件事：

1. 通过 lsof 和 strace 两个工具获取指定进程的一段时间的数据，并保存到文件中
2. 使用 awk 等工具聚合两个文件的内容。

**pt-ioprofile** 使用 strace 和 lsof 工具监视进程的 IO 并打印出一个文件和 I/O 活动表。默认情况下，它监视 mysqld 进程 30 秒。输出如下：

```bash
Tue Dec 27 15:33:57 PST 2011
Tracing process ID 1833
     total       read      write      lseek  ftruncate filename
  0.000150   0.000029   0.000068   0.000038   0.000015 /tmp/ibBE5opS
```

- read：从文件中读出数据。要读取的文件用文件描述符标识，数据读入一个事先定义好的缓冲区。
- write：把缓冲区的数据写入文件中。
- pread：由于 lseek 和 read 调用之间，内核可能会临时挂起进程，所以对同步问题造成了问题，调用 pread 相当于顺序调用了 lseek 和 read，这两个操作相当于一个捆绑的原子操作。
- pwrite：由于 lseek 和 write 调用之间，内核可能会临时挂起进程，所以对同步问题造成了问题，调用 pwrite 相当于顺序调用了 lseek 和 write，这两个操作相当于一个捆绑的原子操作。
- fsync：确保文件所有已修改的内容已经正确同步到硬盘上，该调用会阻塞等待直到设备报告 IO 完成。
- open：打开一个文件，并返回这个文件的描述符。
- close：close 系统调用用于“关闭”一个文件，close 调用终止一个文件描述符以其文件之间的关联。文件描述符被释放，并能够重新使用。
- lseek：对文件描述符指定文件的读写指针进行设置，也就是说，它可以设置文件的下一个读写位置。
- fcntl：针对(文件)描述符提供控制。

**pt-ioprofile** 通过使用 附加`strace`到进程来工作`ptrace()`，这将使其运行非常缓慢，直到`strace`分离。除了冻结服务器之外，还有一些风险，即进程在与服务器`strace`分离后崩溃或性能不佳，或者`strace`没有完全分离并使进程处于睡眠状态。因此，这应该被视为一种侵入性工具，除非您对此感到满意，否则不应在生产服务器上使用。

**WARNING**: **pt-ioprofile** freezes the server and may crash the process, or make it perform badly after detaching, or leave it in a sleeping state! Before using this tool, please:

- Read the tool’s documentation
- Review the tool’s known “BUGS”
- Test the tool on a non-production server
- Backup your production server and verify the backups

> **pt-ioprofile** should be considered an intrusive tool, and should not be used on production servers unless you understand and accept the risks.

## 安装

从 <https://www.percona.com/downloads/percona-toolkit/LATEST/> 选择指定环境下的指定版本下载安装包

## Syntax(语法)

**pt-ioprofile \[OPTIONS] \[FILE]**

**OPTIONS**

- **--aggregate, -a** # 聚合结果的方式，可用的值有 sum 与 avg。`默认值：sum`
- 如果求和，则每个单元格将包含其中的值的总和。如果 avg，则每个单元格将包含其中值的平均值。
- **--cell, -c \<STRING>** # 统计的数据。`默认值：times`
  - count  # I/O 次数
  - sizes  # I/O 大小
  - times  # I/O 时间
- **--group-by, -g \<STRING>** # 对输出结果进行分组 `默认值：filename`
  - all # 所有输出都在一行
  - filename # 每个文件名输出一行
  - pid # 每个进程 ID 输出一行
- **--profile-pid \<INT>** # 指定要分析的进程的 PID，该值会覆盖 --profile-process 选项。
- **--profile-process \<STRING>** # 指定要分析的进程名称。`默认值：mysqld`
- **--run-time \<INT>** # 程序运行时长，单位秒。`默认值：30`
- **--save-samples \<FILE>** # 将 strace 与 lsof 的输出结果保存到指定的 FILE 中。

**EXAMPLE**

- pt-ioprofile -p=9267 --cell=sizes

## 问题

原文: https://tusundong.top/post/centos7-pt-ioprofile.html

最近在排查 io wait 需要使用到 pt-ioprofile，结果发现在 CentOS7.8 下执行没有结果。

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/zgyenr/1622780378714-e4fa5820-6042-4558-9cbd-51fc885b6dc0.png)

最后在大神`@轻松的鱼`指导下，修改源码，编辑 /usr/bin/pt-ioprofile 文件，添加此行

    • /^strace: Process/ { mode = "strace"; }

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/zgyenr/1622780378720-19f72459-84e1-448e-81d2-7fe581bdb917.png)
