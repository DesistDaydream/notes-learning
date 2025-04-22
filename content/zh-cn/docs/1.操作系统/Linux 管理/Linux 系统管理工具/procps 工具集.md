---
title: procps 工具集
---

# 概述

> 参考：
>
> - https://sourceforge.net/projects/procps-ng/
> - [GitLab 项目，procps-ng/procps](https://gitlab.com/procps-ng/procps)
> - [GitHub 项目，uutils/procps](https://github.com/uutils/procps) # 有人用 Rust 重写了 procps 项目

procps 是一组命令行和全屏实用程序，它们主要从 [Proc 文件系统](/docs/1.操作系统/Kernel/Filesystem/特殊文件系统/proc.md) 中获取信心，该文件系统为内核数据结构提供了一个简单的接口。procps 程序通常集中在描述系统上运行的进程的结构上。包括以下程序(每个标题都是一个程序)

Note：该工具集就算是最小化安装的 linux 发行版系统也是默认包含的~

# free - 显示系统中可用和已用的内存量

详见 [Memory 管理工具](/docs/1.操作系统/Linux%20管理/Linux%20系统管理工具/Memory%20管理工具.md)

# kill - 向指定PID的进程发送信号

可用的信号详见 [Signal(信号)](/docs/1.操作系统/Kernel/Process/Inter%20Process%20Communication/Signal(信号).md)

kill 命令将指定的信号发送到指定的进程或进程组。 如果未指定信号，则发送 TERM 信号。 TERM 信号将杀死不捕获该信号的进程。对于其他过程，由于无法捕获该信号，可能需要使用 KILL（9）信号。

大多数现代 Shell 具有内置的 kill 函数，其用法与此处描述的命令非常相似。'-a' 和 '-p' 选项以及通过命令名称指定进程的可能性是 本地扩展。

如果 sig 为 0，则不发送信号，但仍执行错误检查。

“信号 0” 有点像精神上的 “ping”。在 shell 程序脚本中使用 kill -0 PID 是判断 PID 是否有效的好方法。信号 0 仅用于检查进程是否存在。

## Syntax(语法)

[sysstat 工具集](/docs/1.操作系统/Linux%20管理/Linux%20系统管理工具/sysstat%20工具集.md)

**kill \[-s signal|-p] \[-q sigval] \[-a] \[--] pid...**

# pgrep，pkill，pidwait - 根据名字或其他属性列出进程、发送信号、暂停进程

> 参考：
>
> - [Manual(手册)，pgrep(1)](https://man7.org/linux/man-pages/man1/pgrep.1.html)

pgrep 查看当前正在运行的进程，并列出所有符合匹配模式的进程 ID。比如：`pgrep -u root sshd` 命令将会列出由 root 用户运行的进程命令中包含 `sshd` 字符串的进程 ID。效果如下：

```bash
~]# pgrep -u root sshd -a
1521 sshd: /usr/sbin/sshd -D [listener] 0 of 10-100 startups
16257 sshd: root@pts/1
26155 sshd: root@pts/2
26266 sshd: root@notty
```

## Syntax(语法)

**pgrep \[OPTIONS] PATTERN**

**pkill \[OPTIONS] PATTERN**

**pidwait \[OPTIONS] PATTERN**

PATTERN(模式) 代指正则表达式的匹配模式。比如 pgrep 根据 PATTERN 中的内容匹配进程

> Notes: 默认情况下，仅匹配启动该进程的程序名称，而不匹配完整的命令(用人话说就是不去匹配命令参数、命令路径)。假如一个进程的命令是这样的 `/bin/prometheus --config.file=/etc/prometheus/config_out/prometheus.yaml --enable-feature=remote-write-receiver,exemplar-storage --query.lookback-delta=3h` 此时若是使用 `pgrep bin` 或者 `pgrep 3h` 相搜索路径或者参数都是匹配不到的，默认只能匹配 prometheus 那部分。若想匹配全部进程的命令，可以使用 -f 选项。

**OPTIONS：**

`()` 表示该选项所适用的工具，若没有括号，则说明选项适用于所有三个工具

- **-SIGNAL, --signal SIGNAL** # (pkill)指定要发送的信号。可以使用数字或信号名称。
- **-f, --full** # 对进程的完整命令行进行匹配。默认情况这三个程序通常只会对进程名称进行匹配。
  - 比如 `pgrep -f containerd` 将会匹配到 `3313 /usr/bin/dockerd --containerd=/run/containerd/containerd.sock` 这样的进程。
- **-l, --list-full** # (pgrep)显示出完整的命令行以及进程 ID
- **-t, --terminal \<TERM,...>** # 仅匹配使用指定终端的进程。终端名称不用使用绝对路径。
- **-x, --exact** # 精确匹配。PATTERN 必须与 进程名称 或 进程命令行 完全对应上才会被匹配到。

## EXAMPLE

- 列出名字中包含 docker 的进程号
  - **pgrep docker**
- 列出 containerd 进程的进程号
  - **pgrep -x containerd**
- 踢掉 TTY 为 pts/1 的用户
  - **pkill -kill -t pts/1**
  - 注意：想要获取一个用户所使用的终端，可以通过 [procps 包中的 w 工具](#w%20-%20报告已经登录的用户和这些用户正在执行的命令)即可

# pmap - 报告进程的内存映射

# ps - 报告进程的信息

> 参考：
>
> - [Manual(手册)，ps(1)](https://man7.org/linux/man-pages/man1/ps.1.html)

ps 是 **process status(进程状态)** 的简称

Note：该命令显示出来的带 `[]` 的进程为内核线程，一般不用关注。出现这种情况一般是因为 ps 命令无法获取进程的命令参数，所以会将命令名称放入括号中。毕竟用户态的 ps 命令怎么可能会获得内核内部程序的参数呢~~~

在 Manual 中的 [STANDARD FORMAT SPECIFIERS(标准格式说明符)](https://man7.org/linux/man-pages/man1/ps.1.html#STANDARD_FORMAT_SPECIFIERS) 部分，列出了每个列（i.e. Specifiers(说明符)）表示的含义。这里对常用的 `ps -ef f` 命令输出的内容的 Specifiers 的含义进行说明：

- UID # User ID(用户标识符)。进程所属用户 ID
- PID # Process ID(进程标识符)。
- PPID # Parent Process ID(父进程标识符)。创建该进程的进程的 ID 号。
- C # CPU utilization(CPU 利用率)。
- STIME # Start Time(启动时间)。进程的启动时间，格式为 HH:MM:SS
- TTY # Controlling Terminal(控制终端)。与该进程关联的终端设备
- STAT # 进程的当前状态
  - D # 不可中断的休眠。通常是 IO。
  - R # 运行。正在运行或者在运行队列中等待。
  - S # 休眠。在等待某个事件，信号。
  - T # 停止。进程接收到信息 SIGSTOP，SIGSTP，SIGTIN，SIGTOU 信号。
  - X # 死掉的进程，不应该出现。
  - Z # 僵死进程。
    - 通常还会跟随如下字母表示更详细的状态。
      - < 高优先级
      - N 低优先级
      - L 有 pages 在内存中 locked。用于实时或者自定义 IO。
      - s 进程领导者，其有子进程。
      - l 多线程
        - `-` 位于前台进程组。
- TIME # CPU Time(CPU 时间)。进程在 CPU 上消耗的总时间
- CMD # 启动该进程的命令。包括可执行文件的路径、命令参数、etc.

不常用的说明符

- VSZ # Virtual Memory Size(虚拟内存大小)，包括进程可以访问的所有内存，包括进入交换分区的内容，以及共享库占用的内存。有的地方也称为 total_vm、VIRT
- RRS # Resident Set Size(实际内存用量)，不包括进入交换分区的内存。RSS 包括共享库占用的内存（只要共享库在内存中）。RSS 包括所有分配的栈内存和堆内存。
- LWP # 线程 ID
- NLWP # 线程数量

可以使用 **-o FORMAT** 选项来自定义输出的格式(就是字段)。 FORMAT 是单个参数，格式为空格分隔或逗号分隔的列表，它提供了一种指定单个输出列的方法。 可以在 man 手册的 [STANDARD FORMAT SPECIFIERS(标准格式说明符)](https://man7.org/linux/man-pages/man1/ps.1.html#STANDARD_FORMAT_SPECIFIERS) 部分中找到所有可用的关键字。

```bash
# 标题可以根据需要重命名
~]# ps -o pid,ruser=RealUser -o comm=Command
    PID RealUser Command
   4652 root     bash
   4774 root     ps
# 可以不输出标题行。
~]# ps -o pid= -o comm=
   4652 bash
   4787 ps

# 列宽将根据宽标题增加； 这可以用来加宽WCHAN等列
~]# ps -o pid,wchan=WIDE-WCHAN-COLUMN -o comm
    PID WIDE-WCHAN-COLUMN COMMAND
   4652 -                 bash
   4789 -                 ps
# 可以也提供显式宽度控制
~]# ps opid,wchan:42,cmd
    PID WCHAN                                      CMD
    881 core_sys_select                            /sbin/agetty -o -p -- \u --noclear tty1 linux
   4652 -                                          -bash
   4790 -                                          ps opid,wchan:42,cmd
# 行为因人格而异； 输出可能是名为“ X，comm = Y”的一列或名为“ X”和“ Y”的两列。如有疑问，请使用多个-o选项。
~]# ps -o pid=X,comm=Y
      X Y
   4652 bash
   4791 ps
# 仅输出 启动总时长、PID、进程命令 这三列
~]# ps -p 38095 -o etime,pid,cmd
    ELAPSED     PID CMD
10-03:22:51   38095 /bin/prometheus --web.console.templates=/etc/prometheus/consoles --web.console.libraries=/etc/prometheus/console_libraries --config.file=/etc/prometheus/config_out/prometheus.yml
```

可用的 SPECIFIERS 有很多，下面仅列出常用的几个

- **lstart** # 进程启动的时间。输出格式为 `DDD mmm HH:MM:SS YYY`, 可以由 -D 选项更改（TODO:  用 -D 会报错: error: unsupported SysV option）
- **etime** # 进程启动的总时长，格式为`d-h:m:s`
- **etimes** # 进程启动的总时长，以秒为单位

## Syntax(语法)

**ps \[OPTIONS]**

默认操作显示该 shell 环境下的所有进程

### OPTIONS

#### PROCESS SELECTION(进程选择)

https://man7.org/linux/man-pages/man1/ps.1.html#SIMPLE_PROCESS_SELECTION

一共有两种选择进程的方式，且两种方式互相冲突，比如使用 -e 选项后， -p 选项则毫无意义，依然会输出所有进程

- 整体选择
  - **-e, -A** # 选择所有进程，包括不在本 shell 环境下的进程进行展示
- 按列表选择
  - **-p, --pid \<PIDList>** # 选择 PIDList 中列出来的进程。多个 PID 以逗号分隔
  - **--ppid \<PIDList>** # 选择 PIDList 中列出来的进程的子进程。多个 PID 以逗号分割
- 通用选择
  - **-N, --deselect** # 取消选择。也可以理解为 反向选择。即，选择“通过 整体选择 与 按列表选择 中选择到的”进程以外的所有进程

#### OUTPUT FORMAT CONTROL(输出格式控制)

https://man7.org/linux/man-pages/man1/ps.1.html#OUTPUT_FORMAT_CONTROL

- **-f** # 更多显示信息
- **-l** # 显示进程的详细信息
- **-o \<FORMAT>** # 以自定义的格式 FORMAT 输出信息。FORMAT 是以逗号或空格分隔的参数列表，详见前文
- **-ww** # 更宽的输出，让输出的内容不受屏幕限制，可以换行显示

#### OUTPUT MODIFIERS(输出模式)

https://man.cx/ps#heading8

- **f,--forest** # 以树状结构显示输出结果。与显示线程的选项冲突

#### THREAD DISPLAY(线程显示)

https://man.cx/ps#heading9

显示线程的选项与 -f, --forest 选项不可同时使用。

- **-T** # 显示线程，会多出 SPID 列，这列为 线程 号

## EXAMPLE

简单使用

```bash
~]# ps -elf
F S UID         PID   PPID  C PRI  NI ADDR SZ WCHAN  STIME TTY          TIME CMD
4 S root          1      0  0  80   0 - 32013 ep_pol 15:16 ?        00:00:01 /usr/lib/systemd/systemd --switched-root --system --deserialize 2
1 S root          2      0  0  80   0 -     0 kthrea 15:16 ?        00:00:00 [kthreadd]
```

```bash
~]# ps aux
USER        PID %CPU %MEM    VSZ   RSS TTY      STAT START   TIME COMMAND
root          1  0.0  0.4 128052  6596 ?        Ss   15:16   0:01 /usr/lib/systemd/systemd --switched-root --system --deserialize 22
root          2  0.0  0.0      0     0 ?        S    15:16   0:00 [kthreadd]
```

显示 ps 的完整内容，不受 COMMAND 命令有字符限制影响。说白了就是，让过长的内容可以换行显示，而不是截断

```bash
~]# ps -efww
```

以树状形式显示，且带中括号的内核进程将会放在最上面，与下面的系统进程分开，显示较为直观，效果如下

```bash
~]# ps -ef f
root         1     0  0 Dec24 ?        Ss     0:04 /usr/lib/systemd/systemd --switched-root --system --deserialize 22
.......
root      2827     1  0 Dec24 ?        Ss     0:00 /usr/sbin/sshd -D
root      6400  2827  0 10:51 ?        Ss     0:00  \_ sshd: root@pts/0
root      6402  6400  0 10:51 pts/0    Ss     0:00      \_ -bash
root      6720  6402  0 10:52 pts/0    R+     0:00          \_ ps -ef f
```

不显示内核进程，以树状格式显示

- ps -N -p 2 --ppid 2 -f f

```bash
~]# ps --deselect -p 2 --ppid 2 -f f
UID        PID  PPID  C STIME TTY      STAT   TIME CMD
root         1     0  0 Oct19 ?        Ss     0:16 /usr/lib/systemd/systemd --switched-root --system --deserialize 22
root       478     1  0 Oct19 ?        Ss     0:03 /usr/lib/systemd/systemd-journald
root       496     1  0 Oct19 ?        Ss     0:00 /usr/sbin/lvmetad -f
root       502     1  0 Oct19 ?        Ss     0:00 /usr/lib/systemd/systemd-udevd
root       630     1  0 Oct19 ?        S<sl   0:00 /sbin/auditd
polkitd    653     1  0 Oct19 ?        Ssl    0:01 /usr/lib/polkit-1/polkitd --no-debug
root       654     1  0 Oct19 ?        Ss     0:00 /usr/bin/qemu-ga --method=virtio-serial --path=/dev/virtio-ports/org.qemu.guest_agent.0 --blacklist=guest-file-open,guest-file-close,guest-file-read,guest-file-write,guest-file-seek,guest-file-flush,guest-exec,guest-
root       655     1  0 Oct19 ?        Ss     0:04 /usr/lib/systemd/systemd-logind
root       657     1  0 Oct19 ?        Ss     0:16 /usr/sbin/irqbalance --foreground
dbus       658     1  0 Oct19 ?        Ssl    0:06 /usr/bin/dbus-daemon --system --address=systemd: --nofork --nopidfile --systemd-activation
root       670     1  0 Oct19 ?        Ss     0:02 /usr/sbin/crond -n
chrony     679     1  0 Oct19 ?        S      0:00 /usr/sbin/chronyd
root       689     1  0 Oct19 ?        Ssl    1:14 /usr/sbin/NetworkManager --no-daemon
root       990     1  0 Oct19 ?        Ssl    0:47 /usr/bin/python2 -Es /usr/sbin/tuned -l -P
root       992     1  0 Oct19 ?        Ssl    0:24 /usr/sbin/rsyslogd -n
root     14438     1  0 Oct21 tty1     Ss+    0:00 /sbin/agetty --noclear tty1 linux
root     14445     1  0 Oct21 ttyS0    Ss+    0:00 /sbin/agetty --keep-baud 115200,38400,9600 ttyS0 vt220
root     15151     1  0 Oct21 ?        Ss     0:00 sshd: /usr/sbin/sshd [listener] 0 of 10-100 startups
root     17321 15151  0 09:18 ?        Ss     0:00  \_ sshd: root@pts/0
root     17325 17321  0 09:18 pts/0    Ss     0:00      \_ -bash
root     17365 17325  0 09:23 pts/0    R+     0:00          \_ ps --deselect -p 2 --ppid 2 -f f
```

这是(⊙o⊙)啥？~~

```bash
~]# ps -eo rss,pid,user,command | sort -rn | head -10 | awk '{ hr\[1024**2]="GB"; hr\[1024]="MB";for (x=1024**3; x>=1024; x/=1024) { if ($1>=x) { printf ("%-6.2f %s ", $1/x, hr\[x]); break }} } { printf ("%-6s %-10s ", $2, $3) }{ for ( x=4 ; x<=NF ; x++ ) { printf ("%s ",$x) } print ("\n") }'
15.94  MB 627    root       /usr/bin/python3 /usr/bin/networkd-dispatcher --run-startup-triggers
15.18  MB 683    root       /usr/bin/python3 /usr/share/unattended-upgrades/unattended-upgrade-shutdown --wait-for-signal
```

# pwdx - 报告进程的当前目录

# skill - Obsolete version of pgrep/pkill

# slabtop - 实时显示内核slab缓存信息

# snice - Renice a process

# sysctl - 在运行时读取或写入内核参数

详见 [sysctl](/docs/1.操作系统/Linux%20管理/Linux%20内核管理工具/sysctl.md)

# tload - Graphical representation of system load average

# top - 运行中的进程的实时动态视图

```bash
top - 14:06:23 up 70 days, 16:44,  2 users,  load average: 1.25, 1.32, 1.35
Tasks: 206 total,   1 running, 205 sleeping,   0 stopped,   0 zombie
Cpu(s):  5.9%us,  3.4%sy,  0.0%ni, 90.4%id,  0.0%wa,  0.0%hi,  0.2%si,  0.0%st
Mem:  32949016k total, 14411180k used, 18537836k free,   169884k buffers
Swap: 32764556k total,        0k used, 32764556k free,  3612636k cached
  PID USER      PR  NI  VIRT  RES  SHR S %CPU %MEM    TIME+  COMMAND
28894 root      22   0 1501m 405m  10m S 52.2  1.3   2534:16 java
18249 root      18   0 3201m 1.9g  11m S 35.9  6.0 569:39.41 java
```

第一行解析：任务队列信息，同 uptime 命令的执行结果，具体参数说明情况如下：

- **14:06:23** # 当前系统时间
- **up 70 days, 16:44** # 系统已经运行了 70 天 16 小时 44 分钟
- **2 users** # 当前有 2 个用户登录系统
- **load average: 1.25, 1.32, 1.35** # load average 后面的三个数分别是 1 分钟、5 分钟、15 分钟的负载情况。
  - load average 数据是每隔 5 秒钟检查一次活跃的进程数，然后按特定算法计算出的数值。如果这个数除以逻辑 CPU 的数量，结果高于 5 的时候就表明系统在超负荷运转了。
  - 关于 load average 的说明详见：理解 load average—Linux 系统负荷

第二行解析：Tasks # 任务（进程），具体信息说明如下：

- **total** # 系统全部的进程数。现在共有 206 个进程
- **running** # 运行状态的进程数
- **sleeping** # 睡眠状态的进程数
- **stoped** # 已经停止的状态的进程数
- **zombie** # 僵尸状态的进程数。

第三行解析：cpu 状态信息，显示了基于上次刷新时间间隔内，CPU 使用率的百分比。如果 top 命令每 1 秒刷新一次，则下面的信息为 1 秒时间内，CPU 被占用时间的百分比（也就意味着 us 占用 0.059 秒，sy 占用 0.034 秒，空闲了 0.904 秒，st 占用了 0.002 秒）。
具体属性说明如下：(且所有参数的值加起来应为 100%)

- **us** # user cpu time，未改变过 nice 值的用户空间进程的运行时间
- **sy** # system cpu time，内核空间进程的运行时间
- **ni** # user nice cpu time，已改变过 nice 值的用户空间进程的运行时间
- **id** # idle cpu time，CPU 空闲时间。
- **wa** # io wait cpu time，等待磁盘写入完成的时间。该值较高时，说明 IO 等待比较严重，这可能磁盘大量作随机访问造成的，也可能是磁盘性能出现了瓶颈。
- **hi** # hardware irq，硬[中断](/docs/1.操作系统/Kernel/CPU/Interrupts(中断)/Interrupts(中断).md)（Hardware Interrupts）花费的时间
- **si** # software irq，软[中断](/docs/1.操作系统/Kernel/CPU/Interrupts(中断)/Interrupts(中断).md)（Software Interrupts）花费的时间
- **st** # steal time，使用 hypervisor 管理的虚拟机偷取的时间

第四行解析：[Memory](/docs/1.操作系统/Kernel/Memory/Memory.md) 状态，具体信息如下：

- **total** — 物理内存总量（32GB）
- **used** — 使用中的内存总量（14GB）
- **free** — 空闲内存总量（18GB）
- **buffers** — 缓存的内存量 （169M）

第五行解析：swap 交换分区信息，具体信息说明如下：

- **total** — 交换区总量（32GB）
- **used** — 使用的交换区总量（0K）
- **free** — 空闲交换区总量（32GB）
- **cached** — 缓冲的交换区总量（3.6GB）

第六行解析：以下各进程（任务）的状态监控，项目列信息说明如下：

- **PID** # 进程 id
- **USER** # 进程所有者
- **PR** # 进程优先级
- **NI** # nice 值。负值表示高优先级，正值表示低优先级
- **VIRT** # 进程的虚拟内存总量，单位 KiB。即便还没有真正分配物理内存，也会计算在内。
- **RES** # 进程常驻内存的大小，单位 KiB。是进程实际使用的物理内存的大小，但不包括 Swap 和 共享内存
- **SHR** # 共享内存大小，单位 KiB。与其他进程共同使用的共享内存、加载的动态链接库、程序的代码段等等。
- **S** # 进程状态。D=不可中断的睡眠状态 R=运行 S=睡眠 T=跟踪/停止 Z=僵尸进程
- **%CPU** # 上次更新到现在的 CPU 时间占用百分比。注意：这个 CPU 的使用百分比为所有逻辑 CPU 的使用率总和，所有对于多核 CPU 的设备来说，该值很有可能会超过 100%
- **%MEM** # 进程使用的物理内存百分比
- **TIME+** # 进程使用的 CPU 时间总计，单位 1/100 秒
- **COMMAND** # 进程名称（命令名/命令行）

## Syntax(语法)

**top \[OPTIONS]**

除了标注选项，当 top 运行时，可以通过快捷键进行一些操作

- **P** # 按照 CPU 使用率排序
- **M** # 按照内存使用率排序

OPTIONS

- **-d \<NUM>** # 设定整个进程视图更新的秒数，默认为 5 秒更新一次
- **-H**# 显示线程。 如果没有此命令行选项，则显示每个进程中所有线程的总和。 稍后，可以使用“ H”交互式命令来更改它。
- **-p \<PID>** # 指定 PID 进程进行观察

EXAMPLE

- 1 # 展开第三行的 CPU，显示每个逻辑 CPU 的状态信息
- b # 高亮显示处于 Running 状态的进程
- f # 管理所要展示的字段(i.e.第六行的内容)以及按照指定的字段排序。
  - 按 ↑↓ 选择要操作的字段
  - 按空格表示显示或不显示当前字段
  - 按 → 选中当前行，然后按 ↑↓ 将选中的行移动，以便变更该字段所在位置
  - 按 ← 取消选中当前行
  - 按 s 选择光标所在的行作为 排序 标准
  - 按 q 退出当前编辑界面。i.e.再次显示 top 面板
- R # 按照 f 命令里指定的字段进行排序，倒序或者顺序
- 查看 744 进程及其线程的动态试图
  - top -p 744 -H

# uptime -显示系统运行了多长时间

# vmstat - 报告虚拟内存状态，还有 io、系统、cpu 等信息

```bash
~]# vmstat -w
procs -----------------------memory---------------------- ---swap-- -----io---- -system-- --------cpu--------
 r  b         swpd         free         buff        cache   si   so    bi    bo   in   cs  us  sy  id  wa  st
 0  0            0      7282784         2148       463980    0    0     8     1   33   39   0   0 100   0   0
```

**顶部字段说明：**

procs

- r # 可运行的进程数(正在运行或等待运行)。即就绪队列的长度
- b # 等待 I/O 完成时阻塞(blocked)的进程数。即不可中断睡眠状态的进程数。

memory

- 详见 内存管理。(单位是 KiB)

swap

- si # 每秒从 swap 到内存的大小。(单位是 KiB)
- so # 每秒从内存到 swap 的大小。(单位是 KiB)

io

- bi # 每秒从块设备接收的块数。即磁盘读取速度。
- bo # 每秒发送到块设备的块数。即磁盘写入速度。
  - 注意：bi 与 bo 的单位为`块/秒`。因为 Linux 中块的大小是 1 KiB，所以这个单位也就等价于 KiB/s

system

- in # 每秒中断的次数。包括时钟的中断
- cs # 每秒上下文切换次数

cpu

- 详见：CPU 使用率。(单位是百分比)

## Syntax(语法)

**vmstat \[OTIONS] \[DELAY \[COUNT]]**

DELAY # 延迟时间(单位是秒)。指定 DELAY 后，程序每隔 DELAY 时间运行一次。如果未指定 DELAY，则值输出一行结果，其中包含自 vmstat 启动到结束的平均值。

- **COUNT** # 每隔 DELAY 时间，程序运行的次数。若不指定则一直运行。

OPTIONS

- **-w** # 格式化输出。如果不加 -w ,则输出非常紧凑，不利于人们观察，且每行最大 80 字符。

EXAMPLE

- 每隔一秒输出一行信息，一共输出 5 次
  - vmstat 1 5

# w - 报告已经登录的用户和这些用户正在执行的命令

```bash
~]# w
 09:22:37 up 22:46,  1 user,  load average: 0.00, 0.00, 0.00
USER     TTY      FROM             LOGIN@   IDLE   JCPU   PCPU WHAT
root     pts/0    172.19.42.203    09:06    0.00s  0.21s  0.01s w
```

- USER：用哪个用户登录的
- TTY：为该用户开启的终端名
- FROM：该登录用户从哪个 IP 登录的
- LOGIN@：用户登录的时间
- IDLE：用户登录
- WHAT：该用户当前正在执行的命令

获取到的用户信息中，TTY 的信息可以被 pkill 工具使用，以踢掉用户，让其下线

# watch - 定期执行程序，在全部屏幕上显示输出结果

该工具就是持续执行同一个命令，并实时显示
