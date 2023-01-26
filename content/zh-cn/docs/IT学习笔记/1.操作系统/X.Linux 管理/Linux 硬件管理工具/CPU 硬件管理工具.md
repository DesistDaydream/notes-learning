---
title: CPU 硬件管理工具
---

# 概述

> 参考：

# Linux 环境查看 CPU 信息

/proc/stat # 参见：[VFS 虚拟文件系统](https://www.yuque.com/desistdaydream/learning/rsm2ly#e87jG)

通过 cat /proc/cpuinfo 命令，可以查看 CPU 相关的信息：

```bash
~]# cat /proc/cpuinfo
processor : 0
vendor_id : GenuineIntel
cpu family : 6
model : 44
model name : Intel(R) Xeon(R) CPU           E5620  @ 2.40GHz
stepping : 2
cpu MHz : 1596.000
cache size : 12288 KB
physical id : 0
siblings : 8
core id : 0
cpu cores : 4
apicid : 0
initial apicid : 0
fpu : yes
fpu_exception : yes
cpuid level : 11
wp : yes
flags : fpu vme de pse tsc msr pae mce cx8 apic sep mtrr pge mca cmov pat pse36 clflush dts acpi mmx fxsr sse sse2 ss ht tm pbe syscall nx pdpe1gb rdtscp lm constant_tsc arch_perfmon pebs bts rep_good xtopology nonstop_tsc aperfmperf pni pclmulqdq dtes64 monitor ds_cpl vmx smx est tm2 ssse3 cx16 xtpr pdcm pcid dca sse4_1 sse4_2 popcnt aes lahf_lm arat epb dts tpr_shadow vnmi flexpriority ept vpid
bogomips : 4800.15
clflush size : 64
cache_alignment : 64
address sizes : 40 bits physical, 48 bits virtual
power management:
......
```

## CPU 核心数相关信息

在查看到的相关信息中，通常有些信息比较让人迷惑，这里列出一些解释：

- **physical id** # 指的是物理封装的处理器的 id。
- **cpu cores** # 位于相同物理封装的处理器中的内核数量。
- **core id **# 每个内核的 id (不一定是按顺序排列的数字) 。
- **siblings** # 位于相同物理封装的处理器中的逻辑处理器的数量。
- **processor** # 逻辑处理器的 id。

我们通常可以用下面这些命令获得这些参数的信息：

```bash
~]# cat /proc/cpuinfo | grep "physical id" | sort|uniq
physical id     : 0
physical id     : 1
~]# cat /proc/cpuinfo | grep "cpu cores" | sort|uniq
cpu cores     : 4
~]# cat /proc/cpuinfo | grep "core id" | sort|uniq
core id          : 0
core id          : 1
core id          : 10
core id          : 9
~]# cat /proc/cpuinfo | grep "siblings" | sort|uniq
siblings     : 8
~]# cat /proc/cpuinfo | grep "processor" | sort -n -k 2 -t: | uniq
processor	: 0
processor	: 1
processor	: 2
processor	: 3
processor	: 4
processor	: 5
processor	: 6
processor	: 7
processor	: 8
processor	: 9
processor	: 10
processor	: 11
processor	: 12
processor	: 13
processor	: 14
processor	: 15
```

通过上面的结果，可以看出这台机器：

1. physical id # 有 2 个物理处理器(i.e.装在主板上的 CPU)（有 2 个）
2. cpu cores # 每个物理处理器有 4 个内核（为 4）
3. siblings # 每个物理处理器有 8 个逻辑处理器（为 8）
   1. 可见台机器的处理器开启了**超线程技术**，每个内核（core）被划分为了 2 个逻辑处理器（processor）
4. processor # 总共有 16 个逻辑处理器（有 16 个）

**超线程技术**：超线程技术就是利用特殊的硬件指令，把两个逻辑处理器模拟成两个物理芯片，让单个处理器都能使用线程级并行计算，进而兼容多线程操作系统和软件，减少了 CPU 的闲置时间，提高的 CPU 的运行效率。
