---
title: Process info
linkTitle: Process info
weight: 20
---

# 概述

> 参考：
>
> - 

Linux 中有多种途径可以获取进程信息

- [/proc/PID/](#/proc/PID/) 目录
- etc.

# /proc/PID/

> 参考：
>
> - https://github.com/torvalds/linux/blob/v5.19/Documentation/filesystems/proc.rst#11-process-specific-subdirectories

>[!Tip]
>这下面的文件或目录的用途和信息，都有对应的 Manual，Manual 的名字是 proc_pid_XXX 的形式。e.g. stat 文件的 Manual 是 proc_pid_stat(5)

**每个进程在 /proc/ 目录下有一个名为自己进程号的目录，该目录记载了该进程相关的 proc 信息。**

## ./cgroup - 进程所属的控制组信息

详见 [Cgroup](/docs/10.云原生/Containerization/2.CGroup/2.CGroup.md)

## ./cmdline - 该进程的完整命令

- 除非进程是僵尸和内核态进程。在后一种情况下，这个文件中什么也没有：也就是说，对这个文件的读取将返回 0 个字符。命令行参数在这个文件中以一组字符串的形式出现，用空字节('\0')隔开，在最后一个字符串之后还有一个空字节。
- 把这个文件看作是进程希望你看到的命令行。

## ./exe - 具有该 PID 的实际运行的程序的绝对路径。是一个符号链接

## ./fd/ - 其中包含 PID 进程打开的每个文件的一个条目

- 该条目由其文件描述符命名，并且是指向实际文件的符号链接。 因此，0 是标准输入，1 是标准输出，2 是标准错误，依此类推。详解见：[File Descriptor(文件描述符)](/docs/1.操作系统/Kernel/Filesystem/文件管理/File%20Descriptor(文件描述符).md)

## ./fdinfo/ - 其中包含 PID 进程打开的每个文件的一个条目，该条目由其文件描述符命名

该目录中的文件仅由进程所有者读取。 可以读取每个文件的内容以获得有关相应文件描述符的信息。 内容取决于相应文件描述符所引用的文件类型。详解见：[File Descriptor(文件描述符)](/docs/1.操作系统/Kernel/Filesystem/文件管理/File%20Descriptor(文件描述符).md)

## ./maps - 进程的内存映射信息

```bash
~]# cat /proc/1751/maps
00400000-00401000 r-xp 00000000 fd:01 100897359                          /opt/java/jdk1.8.0_231/bin/java
00600000-00601000 r--p 00000000 fd:01 100897359                          /opt/java/jdk1.8.0_231/bin/java
00601000-00602000 rw-p 00001000 fd:01 100897359                          /opt/java/jdk1.8.0_231/bin/java
01542000-01563000 rw-p 00000000 00:00 0                                  [heap]
6c7c00000-6e0100000 rw-p 00000000 00:00 0
6e0100000-76d400000 ---p 00000000 00:00 0
76d400000-797580000 rw-p 00000000 00:00 0
797580000-7c0000000 ---p 00000000 00:00 0
7c0000000-7c18a0000 rw-p 00000000 00:00 0
7c18a0000-800000000 ---p 00000000 00:00 0
```

- address 字段表示进程中内存映射占据的地址空间，格式为十六进制的 BeginAddress-EndAddress。
- perms 字段表示权限，共四个字符，依次为 rwxs 或 rwxp，其中 r 为 read，w 为 write，x 为 execute，s 为- -shared，p 为 private，对应位置没有权限时用一个短横线代替。
- offset 字段表示内存映射地址在文件中的字节偏移量。
- dev 字段表示 device，格式为 major:minor。
- inode 字段表示对应 device 的 inode，0 表示内存映射区域没有关联的 inode，如未初始化的 BSS 数据段就是这种情况。
- pathname 字段用于内存映射的文件，对于 ELF 格式的文件来说，可以通过命令 readelf -l 查看 ELF 程序头部的 Offset 字段，与 maps 文件的 offset 字段作对比。pathname 可能为空，表示匿名映射，这种情况下难以调试进程，如 gdb、strace 等命令。除了正常的文件路径之外，pathname 还可能是下面的值：
  - \[stack]     初始进程（主线程）的 stack
  - \[stack:\<tid>]     线程 ID 为 tid 的 stack.  对应于/proc/\[pid]/task/\[tid]/路径
  - \[vdso]    Virtual Dynamically linked Shared Object
  - \[heap]     进程的 heap

## ./mountinfo - 进程的 mount namespace 的信息

该文件比 /proc/PID/mounts 的信息更全免并修复了一些问题，适用于云原生时代的大量 mount namesapces。

该文件中每一行都是一条挂载信息，每条挂载信息由如下几个部分组成：

| 挂载 ID | 父 ID | major:minor | root  | mount point | mount options | optional fields | separator | filesystem type | mount source | super options      |
| ------- | ----- | ----------- | ----- | ----------- | ------------- | --------------- | --------- | --------------- | ------------ | ------------------ |
| 36      | 35    | 98:00:00    | /mnt1 | /mnt2       | rw,noatime    | master:1        | -         | ext3            | /dev/root    | rw,errors=continue |

- mount ID # 挂载的唯一 ID
- parent ID # the ID of the parent mount (or of self for the root of this mount namespace's mount tree).
  - If a new mount is stacked on top of a previous existing mount (so that it hides the existing mount) at pathname P, then the parent of the new mount is the previous mount at that location. Thus, when looking at all the mounts stacked at a particular location, the top-most mount is the one that is not the parent of any other mount at the same location. (Note, however, that this top-most mount will be accessible only if the longest path subprefix of P that is a mount point is not itself hidden by a stacked mount.)
  - If the parent mount lies outside the process's root directory (see chroot(2)), the ID shown here won't have a corresponding record in mountinfo whose mount ID (field 1) matches this parent mount ID (because mounts that lie outside the process's root directory are not shown in mountinfo). As a special case of this point, the process's root mount masy have a parent mount (for the initramfs filesystem) that lies outside the process's root directory, and an entry for that mount will not appear in mountinfo.
- major:minor # the value of st_dev for files on this filesystem (see [stat(2)](https://man7.org/linux/man-pages/man2/stat.2.html)).
- root # the pathname of the directory in the filesystem which forms the root of this mount.
- mount point # 挂载点 the pathname of the mount point relative to the process's root directory.
- mount options # 挂载选项 per-mount options (see mount(2)).
- optional fields # zero or more fields of the form "tag\[:value]"; see below.
- separator # the end of the optional fields is marked by a single hyphen.
- filesystem type # 挂载的文件系统类型 the filesystem type in the form "type\[.subtype]".
- mount source # filesystem-specific information or "none".
- super options # 超级快选项。per-superblock options (see mount(2)).

## ./root/ - 每个进程的文件系统的 `/` 目录

/proc/PID/root/ 目录是一个指向进程根目录的软链接，效果如下：

```bash
~]# ls -l /proc/1192/root
lrwxrwxrwx 1 root root 0 Nov 11 10:42 root -> /

```

该目录通过 chroot(2) 系统调用设置。

该目录常被用来查看容器内的文件系统。与容器的 Merged 不同，该目录会包含所有挂载，这些挂载信息，来源于 /proc/PID/mountinfo 文件。

## ./smaps - 每个进程的内存映射的使用信息

> [!Note]
> 仅当内核选项 CONFIG_PROC_PAGE_MONITOR 配置了之后，才存在该文件。`pmap` 命令会读取该文件，并以人类易读的形式显示信息。

进程的每一个映射，都有其对应的信息，文件格式如下：

```bash
00400000-00aa4000 r-xp 00000000 fc:01 1710957                            /bin/node_exporter
Size:               6800 kB
KernelPageSize:        4 kB
MMUPageSize:           4 kB
Rss:                5852 kB
......
00aa4000-01077000 r--p 006a4000 fc:01 1710957                            /bin/node_exporter
Size:               5964 kB
KernelPageSize:        4 kB
MMUPageSize:           4 kB
Rss:                5508 kB
......
7fe089fd7000-7fe08c708000 rw-p 00000000 00:00 0
......
7fe08c708000-7fe09c888000 ---p 00000000 00:00 0
......
7fe09c888000-7fe09c889000 rw-p 00000000 00:00 0
......
......
ffffffffff600000-ffffffffff601000 --xp 00000000 00:00 0                  [vsyscall]
......
```

该进程的映射信息与 /proc/PID/maps 文件中的内容相同。该文件的每一个映射信息下面的都包含内存使用量

- **Size** # 该映射的内存大小
  - 通过 `cat smaps| grep ^Size | awk '{print $2}' | awk '{sum += $1};END {print sum}' && cat status | grep VmSize` 命令，可以看到 smaps 文件中的 Size 与 status 文件中的 VmSize 是相同的
- **Rss** #
- **Pss** #
- ......

## ./smaps_rollup - 汇总 smaps 文件中每个映射的内存信息为一条结果

smaps_rollup 文件中的内容，是将 smaps 文件中每个映射的内存信息汇总之后的结果

```bash
[root@hw-cloud-xngy-jump-server-linux-2 /proc/1185]# cat smaps_rollup
00400000-7ffc50bfa000 ---p 00000000 00:00 0                              [rollup]
Rss:               21952 kB
Pss:               21948 kB
Pss_Anon:          10372 kB
Pss_File:          11576 kB
Pss_Shmem:             0 kB
Shared_Clean:          4 kB
Shared_Dirty:          0 kB
Private_Clean:     11576 kB
Private_Dirty:     10372 kB
Referenced:        21952 kB
Anonymous:         10372 kB
LazyFree:              0 kB
AnonHugePages:         0 kB
ShmemPmdMapped:        0 kB
FilePmdMapped:        0 kB
Shared_Hugetlb:        0 kB
Private_Hugetlb:       0 kB
Swap:                  0 kB
SwapPss:               0 kB
Locked:                0 kB
```

通过 awk 计算 Rss 的大小，可以看到 smaps 文件中的聚合值与 smaps_rollup 的值一样

```bash
[root@hw-cloud-xngy-jump-server-linux-2 /proc/1185]# cat smaps| grep ^Rss | awk '{print $2}' | awk '{sum += $1};END {print sum}'
21304
[root@hw-cloud-xngy-jump-server-linux-2 /proc/1185]# cat smaps_rollup
00400000-7ffc50bfa000 ---p 00000000 00:00 0                              [rollup]
Rss:               21304 kB
```

## ./statm - 进程的内存状态信息

https://man7.org/linux/man-pages/man5/proc_pid_statm.5.html

> Tips: 由于内核内部可伸缩性优化，文件中的一些值并绝对不准确，如果需要准确的值，可以查看 smaps 和 smaps_rollup 文件。

该文件中只有一行，每列信息以空格分隔，共 7 列

- **size** # Virtual memory size(虚拟内存大小，简称 vms)。statm 文件中该字段是进程在虚拟内存拥有的 **页数**。
    - 该值与系统 PageSize(页容量) 相乘的计算结果是 status 文件中 VmSize 的值
- **resident** # Resident Set Size(常驻集大小，简称 rss)。statm 文件中的该字段表示的是进程在实内存中拥有的 **页数**，i.e. rss 页数。
    - 该值与系统 PageSize(页容量) 相乘的计算结果是 status 文件中 VmRSS 的值
- **shared** # 与其他进程共享的内存。
- **trs**
- **lrs**
- **drs**
- **dt**

```bash
~]# PID=10000
~]# echo $(($(cat /proc/$PID/statm | awk '{print $1}') * $(getconf PAGE_SIZE) / 1024)); grep VmSize /proc/$PID/status
2083364
VmSize:  2083364 kB
~]# echo $(($(cat /proc/$PID/statm | awk '{print $2}') * $(getconf PAGE_SIZE) / 1024)); grep VmRSS /proc/$PID/status
126876
VmRSS:    126876 kB
```

通过如下方式，可以启动一个申请了 1 MiB 内存的进程，以便验证各种信息

```bash
python3 -c '
import time
import os
# 申请1 MiB的内存
data = bytearray(1024 * 1024)
print(f"进程ID: {os.getpid()}")
print(f"已分配1 MiB内存")
print(f"按Ctrl+C结束进程")
while True:
    time.sleep(1)
'
```

## ./stat - 进程的状态

https://man7.org/linux/man-pages/man5/proc_pid_stat.5.html

https://stackoverflow.com/questions/39066998/what-are-the-meaning-of-values-at-proc-pid-stat

22. **starttime** # 系统 boot 成功后，进程的时间。该值以时钟周期表示。

> starttime 除以 时钟频率（`getconf -a | grep CLK_TCK`），加上系统启动时间（`cat /proc/stat | grep btime`）。可以计算出进程的启动时间（Unix 时间戳格式）

23. **vsize** # 该值的单位是 Bytes，该值除以 `$(getconf PAGE_SIZE)` 的结果与 /proc/PID/statm 文件中第一字段 size 的值相同
24. **rss** # 与 /proc/PID/statm 文件中第二字段 resident 的值相同

计算进程启动时间的脚本

```bash
#!/bin/bash

# 获取进程的 PID
pid=$1

# 获取进程的 starttime（以时钟周期表示）
starttime=$(cat /proc/$pid/stat | awk '{print $22}')

# 获取系统的时钟频率（CLK_TCK）
clk_tck=$(getconf CLK_TCK)

# 获取系统的启动时间（btime）
btime=$(cat /proc/stat | grep btime | awk '{print $2}')

# 计算进程的启动时间（Unix 时间戳）
starttime_seconds=$(echo "$starttime / $clk_tck" | bc)
process_start_time=$(echo "$btime + $starttime_seconds" | bc)

# 输出进程的启动时间（Unix 时间戳）
echo "Process start time (Unix timestamp): $process_start_time"

# 如果需要转换为可读的时间格式，可以使用 date 命令
echo "Process start time (readable format): $(date -d @$process_start_time)"
```

## ./status - 以人类可读形式呈现的进程状态

https://man7.org/linux/man-pages/man5/proc_pid_status.5.html

包括但不限于 PID、该进程 CPU 与 Memory 的使用情况、等等。在这个文件中，包含了 ./stat 和 ./statm 文件中的许多信息。

- **Name** # 进程名称
- **State** # 进程状况
- **VmSize** # 进程申请的总内存。结果与 statm 文件中第一个字段的值相关
- **VmRSS** # 进程当前时刻占用的物理内存。结果与 statm 文件中第二个字段的值相关
- etc.

## ./task/TID/ - 进程的线程信息目录

该目录中是进程的每个线程的信息，目录名称(TID)为线程 ID
