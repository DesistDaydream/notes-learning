---
title: Memroy 的 Overcommit 与 OOM
linkTitle: Memroy 的 Overcommit 与 OOM
weight: 20
---

# 概述

> 参考：
>
> - [Linux Kernel 文档，子系统 - 内存管理 - Overcommit 设计](https://www.kernel.org/doc/html/latest/mm/overcommit-accounting.html)
> - [Linux Kernel 文档，子系统 - 内存管理 - OOM](https://www.kernel.org/doc/html/latest/mm/oom.html)
>     - 在 [2022 年 5 月 10 日 Commit 481cc97](https://github.com/torvalds/linux/commit/481cc97349d694e3211e14a886ad2b7ef55b5a2c) 中，创建了一系列内存管理的 .rst 文档文件，其中有 oom 处理。这次提交说明：是为了紧跟 Mel Gorman 的书 "Understanding the Linux Virtual Memory Manager"，而创建了一系列的大纲，以便以后可以将文档转换到更合理的位置

over commit memory 与 out of memory 机制

## over-commit memory 机制

Linux 内核根据应用程序的要求分配内存，通常来说应用程序分配了内存但是并没有实际全部使用，为了提高性能，这部分没用的内存可以留作它用，这部分内存是属于每个进程的，内核直接回收利用的话比较麻烦，所以内核采用一种 **over-commit memory(过度分配内存)** 的办法来间接利用这部分 “空闲” 的内存，提高整体内存的使用效率。一般来说这样做没有问题，但当大多数应用程序都消耗完自己的内存的时候麻烦就来了，因为这些应用程序的内存需求加起来超出了物理内存（包括 swap）的容量，内核（OOM killer）必须杀掉一些进程才能腾出空间保障系统正常运行。用银行的例子来讲可能更容易懂一些，部分人取钱的时候银行不怕，银行有足够的存款应付，当全国人民（或者绝大多数）都取钱而且每个人都想把自己钱取完的时候银行的麻烦就来了，银行实际上是没有这么多钱给大家取的。

## out of memory(OOM) 机制

某时刻应用程序大量请求内存导致系统内存不足造成的，这通常会触发 Linux 内核里的 Out of Memory (OOM) killer，OOM killer 会杀掉某个进程以腾出内存留给系统用，不致于让系统立刻崩溃

内核检测到系统内存不足、挑选并杀掉某个进程的过程可以参考内核源代码 linux/mm/oom_kill.c，当系统内存不足的时候，out_of_memory() 被触发，然后调用 select_bad_process() 选择一个 “bad” 进程杀掉，如何判断和选择一个 “bad” 进程呢，总不能随机选吧？挑选的过程由 oom_badness() 决定，挑选的算法和想法都很简单很朴实：最 bad 的那个进程就是那个最占用内存的进程。

## OOM 触发后的 Message 信息

```bash
Nov 24 19:52:22 dr-2 kernel: dsm_sa_datamgrd invoked oom-killer: gfp_mask=0x201da, order=0, oom_adj=0, oom_score_adj=0
Nov 24 19:52:22 dr-2 kernel: dsm_sa_datamgrd cpuset=/ mems_allowed=0-1
Nov 24 19:52:22 dr-2 kernel: Pid: 4917, comm: dsm_sa_datamgrd Not tainted 2.6.32-279.19.1.el6.x86_64 #1
Nov 24 19:52:22 dr-2 kernel: Call Trace:
Nov 24 19:52:22 dr-2 kernel: [<ffffffff810c29e1>] ? cpuset_print_task_mems_allowed+0x91/0xb0
Nov 24 19:52:22 dr-2 kernel: [<ffffffff81112d40>] ? dump_header+0x90/0x1b0
Nov 24 19:52:22 dr-2 kernel: [<ffffffff810e064e>] ? __delayacct_freepages_end+0x2e/0x30
Nov 24 19:52:22 dr-2 kernel: [<ffffffff8120dfec>] ? security_real_capable_noaudit+0x3c/0x70
Nov 24 19:52:22 dr-2 kernel: [<ffffffff811131c2>] ? oom_kill_process+0x82/0x2a0
Nov 24 19:52:22 dr-2 kernel: [<ffffffff811130be>] ? select_bad_process+0x9e/0x120
Nov 24 19:52:22 dr-2 kernel: [<ffffffff81113600>] ? out_of_memory+0x220/0x3c0
Nov 24 19:52:22 dr-2 kernel: [<ffffffff8112331e>] ? __alloc_pages_nodemask+0x89e/0x940
Nov 24 19:52:22 dr-2 kernel: [<ffffffff811574ea>] ? alloc_pages_current+0xaa/0x110
Nov 24 19:52:22 dr-2 kernel: [<ffffffff811101c7>] ? __page_cache_alloc+0x87/0x90
Nov 24 19:52:22 dr-2 kernel: [<ffffffff81125cfb>] ? __do_page_cache_readahead+0xdb/0x210
Nov 24 19:52:22 dr-2 kernel: [<ffffffff81125e51>] ? ra_submit+0x21/0x30
Nov 24 19:52:22 dr-2 kernel: [<ffffffff811114f3>] ? filemap_fault+0x4c3/0x500
Nov 24 19:52:22 dr-2 kernel: [<ffffffff8113a754>] ? __do_fault+0x54/0x510
Nov 24 19:52:22 dr-2 kernel: [<ffffffff8113ad07>] ? handle_pte_fault+0xf7/0xb50
Nov 24 19:52:22 dr-2 kernel: [<ffffffff8105a5c3>] ? perf_event_task_sched_out+0x33/0x80
Nov 24 19:52:22 dr-2 kernel: [<ffffffff8113b99a>] ? handle_mm_fault+0x23a/0x310
Nov 24 19:52:22 dr-2 kernel: [<ffffffff810432d9>] ? __do_page_fault+0x139/0x480
Nov 24 19:52:22 dr-2 kernel: [<ffffffff81095bdf>] ? hrtimer_try_to_cancel+0x3f/0xd0
Nov 24 19:52:22 dr-2 kernel: [<ffffffff81095c92>] ? hrtimer_cancel+0x22/0x30
Nov 24 19:52:22 dr-2 kernel: [<ffffffff814eb723>] ? do_nanosleep+0x93/0xc0
Nov 24 19:52:22 dr-2 kernel: [<ffffffff81095d64>] ? hrtimer_nanosleep+0xc4/0x180
Nov 24 19:52:22 dr-2 kernel: [<ffffffff81094af0>] ? hrtimer_wakeup+0x0/0x30
Nov 24 19:52:22 dr-2 kernel: [<ffffffff814ef68e>] ? do_page_fault+0x3e/0xa0
Nov 24 19:52:22 dr-2 kernel: [<ffffffff814eca45>] ? page_fault+0x25/0x30
Nov 24 19:52:22 dr-2 kernel: Mem-Info:
Nov 24 19:52:22 dr-2 kernel: Node 0 Normal per-cpu:
Nov 24 19:52:22 dr-2 kernel: CPU    0: hi:  186, btch:  31 usd:   0
Nov 24 19:52:22 dr-2 kernel: CPU    1: hi:  186, btch:  31 usd:   0
Nov 24 19:52:22 dr-2 kernel: CPU    2: hi:  186, btch:  31 usd:   2
Nov 24 19:52:22 dr-2 kernel: CPU    3: hi:  186, btch:  31 usd:   0
Nov 24 19:52:22 dr-2 kernel: CPU    4: hi:  186, btch:  31 usd:   0
Nov 24 19:52:22 dr-2 kernel: CPU    5: hi:  186, btch:  31 usd:   0
Nov 24 19:52:22 dr-2 kernel: CPU    6: hi:  186, btch:  31 usd:   0
Nov 24 19:52:22 dr-2 kernel: CPU    7: hi:  186, btch:  31 usd:   0
Nov 24 19:52:22 dr-2 kernel: CPU    8: hi:  186, btch:  31 usd:   0
Nov 24 19:52:22 dr-2 kernel: CPU    9: hi:  186, btch:  31 usd:   0
Nov 24 19:52:22 dr-2 kernel: CPU   10: hi:  186, btch:  31 usd:   0
Nov 24 19:52:22 dr-2 kernel: CPU   11: hi:  186, btch:  31 usd:   0
Nov 24 19:52:22 dr-2 kernel: CPU   12: hi:  186, btch:  31 usd:   0
Nov 24 19:52:22 dr-2 kernel: CPU   13: hi:  186, btch:  31 usd:   0
Nov 24 19:52:22 dr-2 kernel: CPU   14: hi:  186, btch:  31 usd:   0
Nov 24 19:52:22 dr-2 kernel: CPU   15: hi:  186, btch:  31 usd:   0
Nov 24 19:52:22 dr-2 kernel: Node 1 DMA per-cpu:
Nov 24 19:52:22 dr-2 kernel: CPU    0: hi:    0, btch:   1 usd:   0
Nov 24 19:52:22 dr-2 kernel: CPU    1: hi:    0, btch:   1 usd:   0
Nov 24 19:52:22 dr-2 kernel: CPU    2: hi:    0, btch:   1 usd:   0
Nov 24 19:52:22 dr-2 kernel: CPU    3: hi:    0, btch:   1 usd:   0
Nov 24 19:52:22 dr-2 kernel: CPU    4: hi:    0, btch:   1 usd:   0
Nov 24 19:52:22 dr-2 kernel: CPU    5: hi:    0, btch:   1 usd:   0
Nov 24 19:52:22 dr-2 kernel: CPU    6: hi:    0, btch:   1 usd:   0
Nov 24 19:52:22 dr-2 kernel: CPU    7: hi:    0, btch:   1 usd:   0
Nov 24 19:52:22 dr-2 kernel: CPU    8: hi:    0, btch:   1 usd:   0
Nov 24 19:52:22 dr-2 kernel: CPU    9: hi:    0, btch:   1 usd:   0
Nov 24 19:52:22 dr-2 kernel: CPU   10: hi:    0, btch:   1 usd:   0
Nov 24 19:52:22 dr-2 kernel: CPU   11: hi:    0, btch:   1 usd:   0
Nov 24 19:52:22 dr-2 kernel: CPU   12: hi:    0, btch:   1 usd:   0
Nov 24 19:52:22 dr-2 kernel: CPU   13: hi:    0, btch:   1 usd:   0
Nov 24 19:52:22 dr-2 kernel: CPU   14: hi:    0, btch:   1 usd:   0
Nov 24 19:52:22 dr-2 kernel: CPU   15: hi:    0, btch:   1 usd:   0
Nov 24 19:52:22 dr-2 kernel: Node 1 DMA32 per-cpu:
Nov 24 19:52:22 dr-2 kernel: CPU    0: hi:  186, btch:  31 usd:  69
Nov 24 19:52:22 dr-2 kernel: CPU    1: hi:  186, btch:  31 usd:  31
Nov 24 19:52:22 dr-2 kernel: CPU    2: hi:  186, btch:  31 usd:   0
Nov 24 19:52:22 dr-2 kernel: CPU    3: hi:  186, btch:  31 usd:  46
Nov 24 19:52:22 dr-2 kernel: CPU    4: hi:  186, btch:  31 usd:   0
Nov 24 19:52:22 dr-2 kernel: CPU    5: hi:  186, btch:  31 usd:  99
Nov 24 19:52:22 dr-2 kernel: CPU    6: hi:  186, btch:  31 usd:   0
Nov 24 19:52:22 dr-2 kernel: CPU    7: hi:  186, btch:  31 usd:  55
Nov 24 19:52:22 dr-2 kernel: CPU    8: hi:  186, btch:  31 usd:  42
Nov 24 19:52:22 dr-2 kernel: CPU    9: hi:  186, btch:  31 usd:  20
Nov 24 19:52:22 dr-2 kernel: CPU   10: hi:  186, btch:  31 usd:   0
Nov 24 19:52:22 dr-2 kernel: CPU   11: hi:  186, btch:  31 usd:  30
Nov 24 19:52:22 dr-2 kernel: CPU   12: hi:  186, btch:  31 usd:   0
Nov 24 19:52:22 dr-2 kernel: CPU   13: hi:  186, btch:  31 usd: 168
Nov 24 19:52:22 dr-2 kernel: CPU   14: hi:  186, btch:  31 usd:   0
Nov 24 19:52:22 dr-2 kernel: CPU   15: hi:  186, btch:  31 usd:  30
Nov 24 19:52:22 dr-2 kernel: Node 1 Normal per-cpu:
Nov 24 19:52:22 dr-2 kernel: CPU    0: hi:  186, btch:  31 usd:   0
Nov 24 19:52:22 dr-2 kernel: CPU    1: hi:  186, btch:  31 usd:   0
Nov 24 19:52:22 dr-2 kernel: CPU    2: hi:  186, btch:  31 usd:   0
Nov 24 19:52:22 dr-2 kernel: CPU    3: hi:  186, btch:  31 usd:   0
Nov 24 19:52:22 dr-2 kernel: CPU    4: hi:  186, btch:  31 usd:   0
Nov 24 19:52:22 dr-2 kernel: CPU    5: hi:  186, btch:  31 usd:   2
Nov 24 19:52:22 dr-2 kernel: CPU    6: hi:  186, btch:  31 usd:   0
Nov 24 19:52:22 dr-2 kernel: CPU    7: hi:  186, btch:  31 usd:   2
Nov 24 19:52:22 dr-2 kernel: CPU    8: hi:  186, btch:  31 usd:   0
Nov 24 19:52:22 dr-2 kernel: CPU    9: hi:  186, btch:  31 usd:   0
Nov 24 19:52:22 dr-2 kernel: CPU   10: hi:  186, btch:  31 usd:   0
Nov 24 19:52:22 dr-2 kernel: CPU   11: hi:  186, btch:  31 usd:   0
Nov 24 19:52:22 dr-2 kernel: CPU   12: hi:  186, btch:  31 usd:   0
Nov 24 19:52:22 dr-2 kernel: CPU   13: hi:  186, btch:  31 usd:   2
Nov 24 19:52:22 dr-2 kernel: CPU   14: hi:  186, btch:  31 usd:   0
Nov 24 19:52:22 dr-2 kernel: CPU   15: hi:  186, btch:  31 usd:   0
Nov 24 19:52:22 dr-2 kernel: active_anon:3693 inactive_anon:959 isolated_anon:0
Nov 24 19:52:22 dr-2 kernel: active_file:17 inactive_file:1177 isolated_file:0
Nov 24 19:52:22 dr-2 kernel: unevictable:0 dirty:6 writeback:0 unstable:0
Nov 24 19:52:22 dr-2 kernel: free:19262 slab_reclaimable:94836 slab_unreclaimable:3898229
Nov 24 19:52:22 dr-2 kernel: mapped:230 shmem:0 pagetables:5388 bounce:0
Nov 24 19:52:22 dr-2 kernel: Node 0 Normal free:15468kB min:45120kB low:56400kB high:67680kB active_anon:0kB inactive_anon:136kB active_file:12kB inactive_file:516kB unevictable:0kB isolated(anon):0kB isolated(file):0kB present:8273920kB mlocked:0kB dirty:4kB writeback:0kB mapped:0kB shmem:0kB slab_reclaimable:349292kB slab_unreclaimable:7792080kB kernel_stack:2824kB pagetables:10744kB unstable:0kB bounce:0kB writeback_tmp:0kB pages_scanned:1098 all_unreclaimable? yes
Nov 24 19:52:22 dr-2 kernel: lowmem_reserve[]: 0 0 0 0
Nov 24 19:52:22 dr-2 kernel: Node 1 DMA free:15740kB min:80kB low:100kB high:120kB active_anon:0kB inactive_anon:0kB active_file:0kB inactive_file:0kB unevictable:0kB isolated(anon):0kB isolated(file):0kB present:15352kB mlocked:0kB dirty:0kB writeback:0kB mapped:0kB shmem:0kB slab_reclaimable:0kB slab_unreclaimable:0kB kernel_stack:0kB pagetables:0kB unstable:0kB bounce:0kB writeback_tmp:0kB pages_scanned:0 all_unreclaimable? yes
Nov 24 19:52:22 dr-2 kernel: lowmem_reserve[]: 0 3243 8041 8041
Nov 24 19:52:22 dr-2 kernel: Node 1 DMA32 free:36536kB min:18112kB low:22640kB high:27168kB active_anon:14772kB inactive_anon:3644kB active_file:56kB inactive_file:4008kB unevictable:0kB isolated(anon):0kB isolated(file):0kB present:3321540kB mlocked:0kB dirty:16kB writeback:0kB mapped:916kB shmem:0kB slab_reclaimable:840kB slab_unreclaimable:2982720kB kernel_stack:0kB pagetables:0kB unstable:0kB bounce:0kB writeback_tmp:0kB pages_scanned:27584 all_unreclaimable? yes
Nov 24 19:52:22 dr-2 kernel: lowmem_reserve[]: 0 0 4797 4797
Nov 24 19:52:22 dr-2 kernel: Node 1 Normal free:9304kB min:26788kB low:33484kB high:40180kB active_anon:0kB inactive_anon:56kB active_file:0kB inactive_file:184kB unevictable:0kB isolated(anon):0kB isolated(file):0kB present:4912640kB mlocked:0kB dirty:4kB writeback:0kB mapped:4kB shmem:0kB slab_reclaimable:29212kB slab_unreclaimable:4818116kB kernel_stack:824kB pagetables:10808kB unstable:0kB bounce:0kB writeback_tmp:0kB pages_scanned:380 all_unreclaimable? yes
Nov 24 19:52:22 dr-2 kernel: lowmem_reserve[]: 0 0 0 0
Nov 24 19:52:22 dr-2 kernel: Node 0 Normal: 2528*4kB 407*8kB 4*16kB 0*32kB 0*64kB 0*128kB 0*256kB 0*512kB 0*1024kB 1*2048kB 0*4096kB = 15480kB
Nov 24 19:52:22 dr-2 kernel: Node 1 DMA: 1*4kB 1*8kB 1*16kB 1*32kB 1*64kB 0*128kB 1*256kB 0*512kB 1*1024kB 1*2048kB 3*4096kB = 15740kB
Nov 24 19:52:22 dr-2 kernel: Node 1 DMA32: 1*4kB 537*8kB 320*16kB 307*32kB 135*64kB 45*128kB 0*256kB 0*512kB 0*1024kB 1*2048kB 0*4096kB = 35692kB
Nov 24 19:52:22 dr-2 kernel: Node 1 Normal: 1670*4kB 8*8kB 0*16kB 4*32kB 8*64kB 7*128kB 2*256kB 1*512kB 0*1024kB 0*2048kB 0*4096kB = 9304kB
Nov 24 19:52:22 dr-2 kernel: 2609 total pagecache pages
Nov 24 19:52:22 dr-2 kernel: 1323 pages in swap cache
Nov 24 19:52:22 dr-2 kernel: Swap cache stats: add 1919969, delete 1918646, find 72763/93573
Nov 24 19:52:22 dr-2 kernel: Free swap  = 25878484kB
Nov 24 19:52:22 dr-2 kernel: Total swap = 32767992kB
Nov 24 19:52:22 dr-2 kernel: 4194303 pages RAM
Nov 24 19:52:22 dr-2 kernel: 115262 pages reserved
Nov 24 19:52:22 dr-2 kernel: 1110 pages shared
Nov 24 19:52:22 dr-2 kernel: 3295246 pages non-shared
Nov 24 19:52:22 dr-2 kernel: [ pid ]   uid  tgid total_vm      rss cpu oom_adj oom_score_adj name
Nov 24 19:52:22 dr-2 kernel: [  624]     0   624     2763        0   1     -17         -1000 udevd
Nov 24 19:52:22 dr-2 kernel: [ 1981]     0  1981    23294        2   0     -17         -1000 auditd
Nov 24 19:52:22 dr-2 kernel: [ 1997]     0  1997    82373      196   0       0             0 rsyslogd
Nov 24 19:52:22 dr-2 kernel: [ 2027]     0  2027     2287       94   0       0             0 irqbalance
Nov 24 19:52:22 dr-2 kernel: [ 2054]    81  2054     5377        1   1       0             0 dbus-daemon
Nov 24 19:52:22 dr-2 kernel: [ 2078]     0  2078     1020        0   9       0             0 acpid
Nov 24 19:52:22 dr-2 kernel: [ 2087]    68  2087     6338        1   7       0             0 hald
Nov 24 19:52:22 dr-2 kernel: [ 2088]     0  2088     4527        1   0       0             0 hald-runner
Nov 24 19:52:22 dr-2 kernel: [ 2116]     0  2116     5063        1   2       0             0 hald-addon-inpu
Nov 24 19:52:22 dr-2 kernel: [ 2131]    68  2131     4452        1   1       0             0 hald-addon-acpi
Nov 24 19:52:22 dr-2 kernel: [ 2153]     0  2153    16019        0   1     -17         -1000 sshd
Nov 24 19:52:22 dr-2 kernel: [ 2169]     0  2169    29303        1   7       0             0 crond
Nov 24 19:52:22 dr-2 kernel: [ 2180]     0  2180     5364        0   0       0             0 atd
Nov 24 19:52:22 dr-2 kernel: [ 2473]     0  2473   272563      264   1       0             0 dsm_sa_datamgrd
Nov 24 19:52:22 dr-2 kernel: [ 2586]     0  2586   122414        0   7       0             0 dsm_sa_datamgrd
Nov 24 19:52:22 dr-2 kernel: [ 2601]     0  2601    73220       30   1       0             0 dsm_sa_eventmgr
Nov 24 19:52:22 dr-2 kernel: [ 2660]     0  2660   125789       89  12       0             0 dsm_sa_snmpd
Nov 24 19:52:22 dr-2 kernel: [ 2734]     0  2734   159845        1   1       0             0 dsm_om_shrsvcd
Nov 24 19:52:22 dr-2 kernel: [ 2767]     0  2767     1016        1   8       0             0 mingetty
Nov 24 19:52:22 dr-2 kernel: [ 2769]     0  2769     1016        1   4       0             0 mingetty
Nov 24 19:52:22 dr-2 kernel: [ 2771]     0  2771     1016        1   7       0             0 mingetty
Nov 24 19:52:22 dr-2 kernel: [ 2773]     0  2773     1016        1  12       0             0 mingetty
Nov 24 19:52:22 dr-2 kernel: [ 2779]     0  2779     1016        1  10       0             0 mingetty
Nov 24 19:52:22 dr-2 kernel: [ 5065]     0  5065  1028466        1   3       0             0 console-kit-dae
Nov 24 19:52:22 dr-2 kernel: [11520]     0 11520     1016        1   0       0             0 mingetty
Nov 24 19:52:22 dr-2 kernel: [23919]    38 23919     6485        1   7       0             0 ntpd
Nov 24 19:52:22 dr-2 kernel: [14167]     0 14167    27401       34  15       0             0 keepalived
Nov 24 19:52:22 dr-2 kernel: [14168]     0 14168    27961       88  15       0             0 keepalived
Nov 24 19:52:22 dr-2 kernel: [14169]     0 14169    27927       31  15       0             0 keepalived
Nov 24 19:52:22 dr-2 kernel: [ 2448]     0  2448     2762        0   2     -17         -1000 udevd
Nov 24 19:52:22 dr-2 kernel: [ 2546]     0  2546     2762        0   3     -17         -1000 udevd
Nov 24 19:52:22 dr-2 kernel: [20557]   188 20557  2140125     3183   0       0             0 haproxy
Nov 24 19:52:22 dr-2 kernel: [16066]     0 16066    95625      285   7       0             0 snmpd
Nov 24 19:52:22 dr-2 kernel: Out of memory: Kill process 20557 (haproxy) score 133 or sacrifice child
Nov 24 19:52:22 dr-2 kernel: Killed process 20557, UID 188, (haproxy) total-vm:8560500kB, anon-rss:12468kB, file-rss:264kB
Nov 24 19:52:40 dr-2 kernel: __ratelimit: 8552 callbacks suppressed
```

## message 信息中的名词解释

进程所占资源列表相关信息

- **pid** # 进程 ID 号
- **uid** # 该进程用户的 UserID
- **tgid** #
- **total_vm** # 该进程所占用的虚拟内存页,1page=4k 内存，所以实际占用需要用该值乘以 4
- **rss** # 该进程所占用的实际内存页,1page=4k 内存，所以实际占用需要用该值乘以 4
- **cpu** #
- **oom_adj** # oom 计算出的该进程的。详见本文《oom killer 怎么挑选进程》章节
- **oom_score_adj** # oom 给该进程的分数。详见本文《oom killer 怎么挑选进程》章节
- **name** # 进程名

**invoked** # 本次 oom 起始于哪个进程的内存申请。通常在第一行

**Killed** # 本次 oom 最终杀掉了哪个进程

# OOM Killer

原文: https://www.jianshu.com/p/ba1cdf92a602

是 Linux 内核设计的一种机制，在内存不足的会后，选择一个占用内存较大的进程并 kill 掉这个进程，以满足内存申请的需求（内存不足的时候该怎么办，其实是个两难的事情，oom killer 算是提供了一种方案吧）

在什么时候触发？

前面说了，在内存不足的时候触发，主要牵涉到【linux 的物理内存结构】和【overcommit 机制】

2.1 内存结构 node、zone、page、order

对于物理内存内存，linux 会把它分很多区（zone），zone 上还有 node 的概念，zone 下面是很多 page，这些 page 是根据 buddy 分配算法组织的，看下面两张图：

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/rwsybc/1616167905260-9c71a405-512b-4594-8868-a3026cd6748b.jpeg)

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/rwsybc/1616167905277-eaf8430e-3eb8-4499-9639-9bf2b8c4572d.jpeg)

上面的概念做下简单的介绍，对后面分析 oom killer 日志很有必要：

- Node：每个 CPU 下的本地内存节点就是一个 Node，如果是 UMA 架构下，就只有一个 Node0,在 NUMA 架构下，会有多个 Node
- Zone：每个 Node 会划分很多域 Zone，大概有下面这些：
- ZONE_DMA：定义适合 DMA 的内存域，该区域的长度依赖于处理器类型。比如 ARM 所有地址都可以进行 DMA，所以该值可以很大，或者干脆不定义 DMA 类型的内存域。而在 IA-32 的处理器上，一般定义为 16M。
- ZONE_DMA32：只在 64 位系统上有效，为一些 32 位外设 DMA 时分配内存。如果物理内存大于 4G，该值为 4G，否则与实际的物理内存大小相同。
- ZONE_NORMAL：定义可直接映射到内核空间的普通内存域。在 64 位系统上，如果物理内存小于 4G，该内存域为空。而在 32 位系统上，该值最大为 896M。
- ZONE_HIGHMEM：只在 32 位系统上有效，标记超过 896M 范围的内存。在 64 位系统上，由于地址空间巨大，超过 4G 的内存都分布在 ZONE_NORMA 内存域。
- ZONE_MOVABLE：伪内存域，为了实现减小内存碎片的机制。
- 分配价值链
- 除了只能在某个区域分配的内存（比如 ZONE_DMA），普通的内存分配会有一个“价值”的层次结构，按分配的“廉价度”依次为：ZONE_HIGHMEM > ZONE_NORMAL > ZONE_DMA。
- 即内核在进行内存分配时，优先从高端内存进行分配，其次是普通内存，最后才是 DMA 内存
- Page：zone 下面就是真正的内存页了，每个页基础大小是 4K，他们维护在一个叫 free_area 的数组结构中
- order：数组的 index，也叫 order，实际对应的是 page 的大小，比如 order 为 0，那么就是一堆 1 个空闲页（4K）组成的链表，order 为 1，就是一堆 2 个空闲页（8K）组成的链表，order 为 2，就是一堆 4 个空闲页（16K）组成的链表

  2.2 overcommit

根据 2.1，已经知道物理内存的大概结构以及分配的规则，不过实际上还有虚拟内存的存在，他的 overcommit 机制和 oom killer 有很大关系：

在实际申请内存的时候，比如申请 1G，并不会在物理区域中分配 1G 的真实物理内存，而是分配 1G 的虚拟内存，等到需要的时候才去真正申请物理内存，也就是说申请不等于分配

所以说，可以申请比物理内存实际大的内存，也就是 overcommit，这样会面临一个问题，就是当真的需要这么多内存的时候怎么办—>oom killer!

vm.overcommit_memory 接受三种值：

- 0 – Heuristic overcommit handling. 这是缺省值，它允许 overcommit，但过于明目张胆的 overcommit 会被拒绝，比如 malloc 一次性申请的内存大小就超过了系统总内存
- 1 – Always overcommit. 允许 overcommit，对内存申请来者不拒。
- 2 – Don’t overcommit. 禁止 overcommit。

## oom killer 怎么挑选进程？

linux 会为每个进程算一个分数，最终他会将分数最高的进程 kill

- /proc/oom_score_adj
  - 在计算最终的 badness score 时，会在计算结果是中加上 oom_score_adj，取值范围为-1000 到 1000
  - 如果将该值设置为-1000，则进程永远不会被杀死，因为此时 badness score 永远返回 0。
- /proc/oom_adj
  - 取值是-17 到+15，取值越高，越容易被干掉。如果是-17，则表示不能被 kill
  - 该设置参数的存在是为了和旧版本的内核兼容
- /proc/oom_score
  - 这个值是系统综合进程的内存消耗量、CPU 时间(utime + stime)、存活时间(uptime - start time)和 oom_adj 计算出的，消耗内存越多分越高，存活时间越长分越低
- 子进程内存：Linux 在计算进程的内存消耗的时候，会将子进程所耗内存的一半同时算到父进程中。这样，那些子进程比较多的进程就要小心了。
- 其他参数（不是很关键，不解释了）
  - /proc/sys/vm/oom_dump_tasks
  - /proc/sys/vm/oom_kill_allocating_task
  - /proc/sys/vm/panic_on_oom
- 关闭 OOM killer

  - sysctl -w vm.overcommit_memory=2
  - echo "vm.overcommit_memory=2" >> /etc/sysctl.conf

    3.1 找出最有可能被杀掉的进程

```bash
tee oomscore.sh > /dev/null <<"EOF"
#!/bin/bash
for proc in $(find /proc -maxdepth 1 -regex '/proc/[0-9]+'); do
    printf "%2d %5d %s\n" \
        "$(cat $proc/oom_score)" \
        "$(basename $proc)" \
        "$(cat $proc/cmdline | tr '\0' ' ' | head -c 50)"
done 2>/dev/null | sort -nr | head -n 10
EOF
chmod +x oomscore.sh
./oomscore.sh
18   981 /usr/sbin/mysqld
 4 31359 -bash
 4 31056 -bash
 1 31358 sshd: root@pts/6
 1 31244 sshd: vpsee [priv]
 1 31159 -bash
 1 31158 sudo -i
 1 31055 sshd: root@pts/3
 1 30912 sshd: vpsee [priv]
 1 29547 /usr/sbin/sshd -D
```

3.2 避免的 oom killer 的方案

- 直接修改 /proc/PID/oom_adj 文件，将其置为 -17
- 修改 /proc/sys/vm/lowmem_reserve_ratio
- 直接关闭 oom-killer

参考：

- node & zone
- 理解 LINUX 的 MEMORY OVERCOMMIT
- linux OOM-killer 机制（杀掉进程，释放内存）
- Taming the OOM killer
- linux OOM 机制分析
- 理解和配置 Linux 下的 OOM Killer
- ubuntu 解决 cache 逐渐变大导致 oom-killer 将某些进程杀死的情况

# OOM Killer 日志分析实践

原文: https://www.jianshu.com/p/8dd45fdd8f33

下面是一台 8G 内存上的一次 oom killer 的日志，上面跑的是 RocketMQ 3.2.6，java 堆配置：-server -Xms4g -Xmx4g -Xmn2g -XX:PermSize=128m -XX:MaxPermSize=320m

```bash
Jun  4 17:19:10 iZ23tpcto8eZ kernel: AliYunDun invoked oom-killer: gfp_mask=0x201da, order=0, oom_score_adj=0
Jun  4 17:19:10 iZ23tpcto8eZ kernel: AliYunDun cpuset=/ mems_allowed=0
Jun  4 17:19:10 iZ23tpcto8eZ kernel: active_anon:1813257 inactive_anon:37301 isolated_anon:0 active_file:84 inactive_file:0 isolated_file:0 unevictable:0 dirty:0 writeback:0 unstable:0 free:23900 slab_reclaimable:34218 slab_unreclaimable:5636 mapped:1252 shmem:100531 pagetables:68092 bounce:0 free_cma:0
Jun  4 17:19:10 iZ23tpcto8eZ kernel: Node 0 DMA free:15900kB min:132kB low:164kB high:196kB active_anon:0kB inactive_anon:0kB active_file:0kB inactive_file:0kB unevictable:0kB isolated(anon):0kB isolated(file):0kB present:15992kB managed:15908kB mlocked:0kB dirty:0kB writeback:0kB mapped:0kB shmem:0kB slab_reclaimable:0kB slab_unreclaimable:8kB kernel_stack:0kB pagetables:0kB unstable:0kB bounce:0kB free_cma:0kB writeback_tmp:0kB pages_scanned:0 all_unreclaimable? yes
Jun  4 17:19:10 iZ23tpcto8eZ kernel: lowmem_reserve[]: 0 2801 7792 7792
Jun  4 17:19:10 iZ23tpcto8eZ kernel: Node 0 DMA32 free:43500kB min:24252kB low:30312kB high:36376kB
active_anon:2643608kB(2.5G)  inactive_anon:61560kB active_file:40kB inactive_file:40kB unevictable:0kB isolated(anon):0kB isolated(file):0kB present:3129216kB managed:2869240kB mlocked:0kB dirty:0kB writeback:0kB mapped:748kB shmem:160024kB slab_reclaimable:54996kB slab_unreclaimable:6816kB kernel_stack:704kB pagetables:67440kB unstable:0kB bounce:0kB free_cma:0kB writeback_tmp:0kB pages_scanned:275 all_unreclaimable? yes
Jun  4 17:19:10 iZ23tpcto8eZ kernel: lowmem_reserve[]: 0 0 4990 4990
Jun  4 17:19:10 iZ23tpcto8eZ kernel: Node 0 Normal free:36200kB min:43192kB low:53988kB high:64788kB
active_anon:4609420kB(4.3G) inactive_anon:87644kB active_file:296kB inactive_file:0kB unevictable:0kB isolated(anon):0kB isolated(file):0kB present:5242880kB managed:5110124kB mlocked:0kB dirty:0kB writeback:0kB mapped:4260kB shmem:242100kB slab_reclaimable:81876kB slab_unreclaimable:15720kB kernel_stack:1808kB pagetables:204928kB unstable:0kB bounce:0kB free_cma:0kB writeback_tmp:0kB pages_scanned:511 all_unreclaimable? yes
Jun  4 17:19:10 iZ23tpcto8eZ kernel: lowmem_reserve[]: 0 0 0 0
Jun  4 17:19:10 iZ23tpcto8eZ kernel: Node 0 DMA: 1*4kB (U) 1*8kB (U) 1*16kB (U) 0*32kB 2*64kB (U) 1*128kB (U) 1*256kB (U) 0*512kB 1*1024kB (U) 1*2048kB (R) 3*4096kB (M) = 15900kB
Jun  4 17:19:10 iZ23tpcto8eZ kernel: Node 0 DMA32: 1281*4kB (UEM) 825*8kB (UEM) 1404*16kB (UEM) 290*32kB (EM) 0*64kB 0*128kB 0*256kB 0*512kB 0*1024kB 0*2048kB 0*4096kB = 43468kB
Jun  4 17:19:10 iZ23tpcto8eZ kernel: Node 0 Normal: 1441*4kB (UEM) 3177*8kB (UEM) 315*16kB (UEM) 0*32kB 0*64kB 0*128kB 0*256kB 0*512kB 0*1024kB 0*2048kB 0*4096kB = 36220kB
Jun  4 17:19:10 iZ23tpcto8eZ kernel: Node 0 hugepages_total=0 hugepages_free=0 hugepages_surp=0 hugepages_size=2048kB
Jun  4 17:19:10 iZ23tpcto8eZ kernel: 100592 total pagecache pages
Jun  4 17:19:10 iZ23tpcto8eZ kernel: 0 pages in swap cache
Jun  4 17:19:10 iZ23tpcto8eZ kernel: Swap cache stats: add 0, delete 0, find 0/0
Jun  4 17:19:10 iZ23tpcto8eZ kernel: Free swap  = 0kB
Jun  4 17:19:10 iZ23tpcto8eZ kernel: Total swap = 0kB
Jun  4 17:19:10 iZ23tpcto8eZ kernel: 2097151 pages RAM
Jun  4 17:19:10 iZ23tpcto8eZ kernel: 94167 pages reserved
Jun  4 17:19:10 iZ23tpcto8eZ kernel: 284736 pages shared
Jun  4 17:19:10 iZ23tpcto8eZ kernel: 1976069 pages non-shared
Jun  4 17:19:10 iZ23tpcto8eZ kernel: [ pid ]   uid  tgid total_vm      rss nr_ptes swapents oom_score_adj name
Jun  4 17:19:10 iZ23tpcto8eZ kernel: [  338]     0   338    10748      844      25        0             0 systemd-journal
Jun  4 17:19:10 iZ23tpcto8eZ kernel: [  351]     0   351    26113       61      20        0             0 lvmetad
Jun  4 17:19:10 iZ23tpcto8eZ kernel: [  368]     0   368    10509      149      23        0         -1000 systemd-udevd
Jun  4 17:19:10 iZ23tpcto8eZ kernel: [  521]     0   521   170342      908     178        0             0 rsyslogd
Jun  4 17:19:10 iZ23tpcto8eZ kernel: [  525]     0   525     8671       82      21        0             0 systemd-logind
Jun  4 17:19:10 iZ23tpcto8eZ kernel: [  526]    81   526     7157       96      19        0          -900 dbus-daemon
Jun  4 17:19:10 iZ23tpcto8eZ kernel: [  530]     0   530    31575      162      17        0             0 crond
Jun  4 17:19:10 iZ23tpcto8eZ kernel: [  540]    28   540   160978      131      37        0             0 nscd
Jun  4 17:19:10 iZ23tpcto8eZ kernel: [  548]     0   548    27501       30      10        0             0 agetty
Jun  4 17:19:10 iZ23tpcto8eZ kernel: [  588]     0   588     1621       26       9        0             0 iprinit
Jun  4 17:19:10 iZ23tpcto8eZ kernel: [  590]     0   590     1621       25       9        0             0 iprupdate
Jun  4 17:19:10 iZ23tpcto8eZ kernel: [  601]     0   601     9781       23       8        0             0 iprdump
Jun  4 17:19:10 iZ23tpcto8eZ kernel: [  838]    38   838     7399      169      18        0             0 ntpd
Jun  4 17:19:10 iZ23tpcto8eZ kernel: [  881]     0   881      386       44       4        0             0 aliyun-service
Jun  4 17:19:10 iZ23tpcto8eZ kernel: [ 5973]  1000  5973    41595      165      32        0             0 gmond
Jun  4 17:19:10 iZ23tpcto8eZ kernel: [ 3829]     0  3829    33413      292      67        0             0 sshd
Jun  4 17:19:10 iZ23tpcto8eZ kernel: [ 3831]  1000  3831    33582      476      68        0             0 sshd
Jun  4 17:19:10 iZ23tpcto8eZ kernel: [ 3832]  1000  3832    29407      622      16        0             0 bash
Jun  4 17:19:10 iZ23tpcto8eZ kernel: [14638]     0 14638    20697      210      42        0         -1000 sshd
Jun  4 17:19:10 iZ23tpcto8eZ kernel: [11531]     0 11531    33413      293      66        0             0 sshd
Jun  4 17:19:10 iZ23tpcto8eZ kernel: [11533]  1000 11533    33413      292      64        0             0 sshd
Jun  4 17:19:10 iZ23tpcto8eZ kernel: [11534]  1000 11534    29361      584      15        0             0 bash
Jun  4 17:19:10 iZ23tpcto8eZ kernel: [ 3172]     0  3172     6338      161      17        0             0 AliYunDunUpdate
Jun  4 17:19:10 iZ23tpcto8eZ kernel: [ 3224]     0  3224    32867     2270      61        0             0 AliYunDun
Jun  4 17:19:10 iZ23tpcto8eZ kernel: [ 5417]  1000  5417    28279       51      14        0             0 sh
Jun  4 17:19:10 iZ23tpcto8eZ kernel: [ 5421]  1000  5421    28279       53      13        0             0 sh
Jun  4 17:19:10 iZ23tpcto8eZ kernel: [ 5424]  1000  5424 36913689  1537770   66407        0             0 java
Jun  4 17:19:10 iZ23tpcto8eZ kernel: [17132]     0 17132    21804      215      44        0             0 zabbix_agentd
Jun  4 17:19:10 iZ23tpcto8eZ kernel: [17133]     0 17133    21804      285      43        0             0 zabbix_agentd
Jun  4 17:19:10 iZ23tpcto8eZ kernel: [17134]     0 17134    21866      290      44        0             0 zabbix_agentd
Jun  4 17:19:10 iZ23tpcto8eZ kernel: [17135]     0 17135    21866      290      44        0             0 zabbix_agentd
Jun  4 17:19:10 iZ23tpcto8eZ kernel: [17136]     0 17136    21841      290      44        0             0 zabbix_agentd
Jun  4 17:19:10 iZ23tpcto8eZ kernel: [17137]     0 17137    21804      245      43        0             0 zabbix_agentd
Jun  4 17:19:10 iZ23tpcto8eZ kernel: [13669]  1000 13669    28279       51      14        0             0 sh
Jun  4 17:19:10 iZ23tpcto8eZ kernel: [13673]  1000 13673    28279       50      13        0             0 sh
Jun  4 17:19:10 iZ23tpcto8eZ kernel: [13675]  1000 13675   879675   204324     494        0             0 java
Jun  4 17:19:10 iZ23tpcto8eZ kernel: Out of memory: Kill process 5424 (java) score 800 or sacrifice child
Jun  4 17:19:10 iZ23tpcto8eZ kernel: Killed process 5424 (java) total-vm:147654756kB, anon-rss:6151080kB, file-rss:0kB
```

各列字段解释：

- min 下的内存是保留给内核使用的；当到达 min，会触发内存的 direct reclaim
- low 水位比 min 高一些，当内存可用量小于 low 的时候，会触发 kswapd 回收内存，当 kswapd 慢慢的将内存 回收到 high 水位，就开始继续睡眠

## 1. 谁申请内存以及谁被 kill 了？

这两个问题，可以从头部和尾部的日志分析出来：

```text
Jun  4 17:19:10 iZ23tpcto8eZ kernel: AliYunDun invoked oom-killer: gfp_mask=0x201da, order=0, oom_score_adj=0
Jun  4 17:19:10 iZ23tpcto8eZ kernel: Killed process 5424 (java) total-vm:147654756kB, anon-rss:6151080kB, file-rss:0kB
```

AliYunDun 申请内存，kill 掉了 java 进程 5424，他占用的内存是 6151080K(5.8G)

还有一个小问题可能会有疑问，那就是进程 5424 的 RSS（1537770）明明小于 6151080，实际是因为这里的 RSS 是 4K 位单位的，所以要乘以 4，算出来就对了

物理内存申请我们在上一篇分析了，会到不同的 Node 不同的 zone，那么这次申请的是哪一部分？这个可以从 gfp_mask=0x201da, order=0 分析出来，gfp_mask(get free page)是申请内存的时候，会传的一个标记位，里面包含三个信息：区域修饰符、行为修饰符、类型修饰符：

```text
0X201da = 0x20000 | 0x100| 0x80 | 0x40 | 0x10 | 0x08 | 0x02
也就是下面几个值：
___GFP_HARDWAL | ___GFP_COLD | ___GFP_FS | ___GFP_IO | ___GFP_MOVABLE| ___GFP_HIGHMEM
```

同时设置了 `___GFP_MOVABLE` 和 `___GFP_HIGHMEM` 会扫描 `ZONE_MOVABLE`，其实也就是会在 ` ZONE_NORMAL`，再贴一次神图

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/rwsybc/1616167905288-838e8528-b134-45bb-bbec-81d175c4c272.jpeg)

另外 order 表示了本次申请内存的大小 0，也就是 4KB

也就是说 AliYunDun 尝试从 ZONE_NORMAL 申请 4KB 的内存，但是失败了，导致了 OOM KILLER

## 2. 各个 zone 的情况如何？

接下来，自然就会问，连 4KB 都没有，那到底还有多少？看这部分日志：

```bash
Jun  4 17:19:10 iZ23tpcto8eZ kernel: Node 0 DMA free:15900kB min:132kB low:164kB high:196kB active_anon:0kB inactive_anon:0kB active_file:0kB inactive_file:0kB unevictable:0kB isolated(anon):0kB isolated(file):0kB present:15992kB managed:15908kB mlocked:0kB dirty:0kB writeback:0kB mapped:0kB shmem:0kB slab_reclaimable:0kB slab_unreclaimable:8kB kernel_stack:0kB pagetables:0kB unstable:0kB bounce:0kB free_cma:0kB writeback_tmp:0kB pages_scanned:0 all_unreclaimable? yes
Jun  4 17:19:10 iZ23tpcto8eZ kernel: lowmem_reserve[]: 0 2801 7792 7792

Jun  4 17:19:10 iZ23tpcto8eZ kernel: Node 0 DMA32 free:43500kB min:24252kB low:30312kB high:36376kB
active_anon:2643608kB(2.5G)  inactive_anon:61560kB active_file:40kB inactive_file:40kB unevictable:0kB isolated(anon):0kB isolated(file):0kB present:3129216kB managed:2869240kB mlocked:0kB dirty:0kB writeback:0kB mapped:748kB shmem:160024kB slab_reclaimable:54996kB slab_unreclaimable:6816kB kernel_stack:704kB pagetables:67440kB unstable:0kB bounce:0kB free_cma:0kB writeback_tmp:0kB pages_scanned:275 all_unreclaimable? yes
Jun  4 17:19:10 iZ23tpcto8eZ kernel: lowmem_reserve[]: 0 0 4990 4990

Jun  4 17:19:10 iZ23tpcto8eZ kernel: Node 0 Normal free:36200kB min:43192kB low:53988kB high:64788kB
active_anon:4609420kB(4.3G) inactive_anon:87644kB active_file:296kB inactive_file:0kB unevictable:0kB isolated(anon):0kB isolated(file):0kB present:5242880kB managed:5110124kB mlocked:0kB dirty:0kB writeback:0kB mapped:4260kB shmem:242100kB slab_reclaimable:81876kB slab_unreclaimable:15720kB kernel_stack:1808kB pagetables:204928kB unstable:0kB bounce:0kB free_cma:0kB writeback_tmp:0kB pages_scanned:511 all_unreclaimable? yes
Jun  4 17:19:10 iZ23tpcto8eZ kernel: lowmem_reserve[]: 0 0 0 0
```

可以看到 Normal 还有 36200KB，DMA32 还有 43500KB，DMA 还有 15900KB，其中 Normal 的 free 确实小于 min，但是 DMA32 和 DMA 的 free 没问题啊？从上篇文章分析来看，分配是有链条的，Normal 不够了，会从 DMA32 以及 DMA 去请求分配，所以为什么分配失败了呢？

### 2.1 lowmem_reserve

虽然说分配内存会按照 Normal、DMA32、DMA 的顺序去分配，但是低端内存相对来说更宝贵些，为了防止低端内存被高端内存用完，linux 设计了保护机制，也就是 lowmen_reserve，从上面的日志看，他们的值是这样的：

- DMA（index=0）: lowmem_reserve\[]:0 2801 7792 7792
- DMA32（index=1）: lowmem_reserve\[]: 0 0 4990 4990
- Normal（index=2）: lowmem_reserve\[]: 0 0 0 0

lowmen_reserve 的值是一个数组，当 Normal(index=2)像 DMA32 申请内存的时候，需要满足条件：DMA32 high+lowmem_reserve\[2] < free，才能申请，来算下：

- Normal：从自己这里申请，free(36200) < min(43192)，所以申请失败了(watermark\[min]以下的内存属于系统的自留内存，用以满足特殊使用，所以不会给用户态的普通申请来用)
- Normal 转到 DMA32 申请:`high(36376KB) + lowmem_reserve[2](4990)\*4=56336KB > DMA32 Free(43500KB)`,不允许申请
- Normal 转到 DMA 申请:`high(196KB) + lowmem_reserve[2](7792)\*4 = 31364KB > DMA Free(15900KB)`,不允许申请,所以....最终失败了

### 2.2 min_free_kbytes

这里属于扩展知识了，和分析 oom 问题不大

我们知道了每个区都有 min、low、high，那他们是怎么计算出来的，就是根据 min_free_kbytes 计算出来的，他本身在系统初始化的时候计算，最小 128K，最大 64M

- watermark\[min] = min_free_kbytes 换算为 page 单位即可，假设为 min_free_pages。（因为是每个 zone 各有一套 watermark 参数，实际计算效果是根据各个 zone 大小所占内存总大小的比例，而算出来的 per zone min_free_pages）
- watermark\[low] = watermark\[min] \* 5 / 4
- watermark\[high] = watermark\[min] \* 3 / 2

min 和 low 的区别：

- min 下的内存是保留给内核使用的；当到达 min，会触发内存的 direct reclaim
- low 水位比 min 高一些，当内存可用量小于 low 的时候，会触发 kswapd 回收内存，当 kswapd 慢慢的将内存 回收到 high 水位，就开始继续睡眠

## 3. 最后的问题：java 为什么占用了这么多内存？

内存不足申请失败的细节都分析清楚了，剩下的问题就是为什么 java 会申请这么多内存(5.8G)，明明-Xmx 配置的是 4G，加上 PermSize，也就最多 4.3G。

因为这上面跑的是 RocketMQ，他会有文件映射 mmap，所以在仔细分析 oom 日志之前，怀疑是 pagecache 占用，导致 RSS 为 5.8G，这带来了另一个问题，为什么 pagecache 没有回收？分析了日志以后，发现和 pagecache 基本没关系，看这个日志(换行是我后来加上的)：

```text
Jun  4 17:19:10 iZ23tpcto8eZ kernel: active_anon:1813257 inactive_anon:37301 isolated_anon:0 active_file:84
inactive_file:0 isolated_file:0
unevictable:0 dirty:0 writeback:0
unstable:0 free:23900
slab_reclaimable:34218
slab_unreclaimable:5636 mapped:1252
shmem:100531 pagetables:68092 bounce:0
free_cma:0
```

当时的内存大部分都是活跃的匿名页(active_anon 18132574KB=6.9G)，其他的非活跃匿名页（inactive_anon 145M），活跃文件页（active_file 844=336KB），非活跃文件页（inactive_file 0），也就是说当时基本没有 pagecache，因为 pagecache 会属于文件页

并且，这台机器上的 gc log 没配置好，进程重启以后 gc 文件被覆盖了，另外被 oom killer 也没有 java dump，所以…..真的不知道到底为什么 java 占了 5.8G!!! 悬案还是没有解开 T_T

如果上层申请内存的速度太快，导致空闲内存降至 `watermark[min]` 后，内核就会进行 direct reclaim（直接回收），即直接在应用程序的进程上下文中进行回收，再用回收上来的空闲页满足内存申请，因此实际会阻塞应用程序，带来一定的响应延迟
