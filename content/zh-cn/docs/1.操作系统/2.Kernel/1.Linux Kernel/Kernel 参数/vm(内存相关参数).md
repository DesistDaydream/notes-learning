---
title: vm(内存相关参数)
---

# 概述

> 参考：
>
> - [官方文档，Linux 内核用户和管理员指南-/proc/sys 文档-/proc/sys/vm 文档](https://www.kernel.org/doc/html/latest/admin-guide/sysctl/vm.html)

/proc/sys/vm/ 目录下的文件可用于调整 Linux Kernel 中有关 Virtual Memory(虚拟内存) 子系统的操作。

### vm.drop_caches = NUM

写入该文件可以清理内存中的缓存。详见 [Memory 的缓存](/docs/1.操作系统/2.Kernel/5.Memory/Memory%20的缓存.md#缓存的清理) 章节中“缓存清理”部分

### vm.swappiness = 10

这个内核参数可以用来调整系统使用 swap 的时机。`默认值：60`

设为 60 即表示：当内存中空闲空间低于 60%的时候，就会开始使用 swap 空间(也就是说系统使用了 40%的内存之后，就开始使用 swap)

### vm.max_map_count = 262144

一个进程可以拥有的 VMA(虚拟内存区域) 的数量。`默认值：65530`

常用于运行 Elasticsearch 服务的设备上。

### vm.overcommit_memory = 1

配置是否允许[内存 overcommit](Memroy%20 的%20Over%20Commit%20 与%20OOM.md 管理/Memroy 的 Over Commit 与 OOM.md)，有 0、1、2 三种模式。`默认值：0`

- **0** # heuristic overcommit(试探式允许 overcommit)。
  - 当应用进程尝试申请内存时，内核会做一个检测。内核将检查是否有足够的**可用内存**可以分配。如果有足够的可用内存，内存申请允许；否则，内存申请失败。
- **1** # always overcommit,never check(总是允许 overcommit)
  - 对于内存的申请请求，内核不会做任何检测，并直接分配内存。
- **2** # never overcommit,always check(永不允许 overcommit)
  - 说是永不允许 overcommit，其实只是通过其他参数来控制 overcommit(过量使用) 的大小。可以分配的总内存将会受到 /proc/meminfo 中的 CommitLimit 这个参数限制。
  - `CommitLimit = (total_RAM - total_huge_TLB) * overcommit_ratio / 100 + total_swap`
    - totaol_RAM # 系统内存总量(就是物理内存)
    - total_huge_TLB # 为 huge pages 保留的内存量，一般没有保留，都是 0
    - overcommit_ratio # /proc/sys/vm/overcommit_ratio 内核参数的值。
    - total_swap # swap 空间的总量
  - 比如我现在有一个 16G 内存的服务器，swap 空间为 16，overcommit_ratio 参数设为 50，那么 CommitLimit 的计算结果为 24G。
    - 此时，如果 /proc/meminfo 中的 Commited_AS 参数 值为 23G，当一个程序申请超过 1G 内存时，则会失败。
  - 所以从根本上讲，模式 2 下，可以分配的内存总量，受 overcommit_ration 这个内核参数控制。所谓的永远不会 overcommit，则是指 overcommit_ration 参数的值小于 100。
  - 注意：从 Linux 内核 3.14 开始，如果 /proc/sys/vm/overcommit_kbytes 参数的值不为 0，则`CommitLimit = overcommit_kbytes + total_swap`

> 所以所有模式都可能会触发 OOM 机制。只不过模式 0 和 1 在程序申请内存时，行为不同；而模式 2 则受 overcommit_ration 参数的限制。

### vm.overcommit_ration = 50

内存 overcommit(过量使用) 的百分比。`默认值：50`。与 vm.overcommit_memory = 2 配合使用，其他情况该参数无效。

如果该值小于 100 ，那么系统分配给程序的内存，则永远不会超过真实内存，所以也就不存在 overcommit。如果该值超过 100，那么 CommitLimit 的值则会超过真实内存，这时候就相当于 overcommit 了，并且在内存使用超过真实值时会触发 OOM。

```bash
# 在没有swap 和 大页预留的情况下
# vm.overcommit_ration=100 时，CommitLimit 等于 Mem 总量
[root@lichenhao ~]# sysctl -w vm.overcommit_ratio=100
vm.overcommit_ratio = 100
[root@lichenhao ~]# cat /proc/meminfo  | grep Commit
 && free -k
CommitLimit:     3868968 kB
Committed_AS:     806056 kB
              total        used        free      shared  buff/cache   available
Mem:        3868968      282904     3195972        8748      390092     3368548
Swap:             0           0           0
# vm.overcommit_ration=50 时，CommitLimit 等于 Mem 总量的一半
[root@lichenhao ~]# sysctl -w vm.overcommit_ratio=100
vm.overcommit_ratio = 50
[root@lichenhao ~]# cat /proc/meminfo  | grep Commit && free -k
CommitLimit:     1934484 kB
Committed_AS:     800516 kB
              total        used        free      shared  buff/cache   available
Mem:        3868968      282712     3196160        8748      390096     3368740
Swap:             0           0           0
# vm.overcommit_ration=200 时，CommitLimit 等于 Mem 总量的一倍。此时内存可以 overcommit，有可能会触发 OOM
[root@lichenhao ~]# sysctl -w vm.overcommit_ratio=100
vm.overcommit_ratio = 200
[root@lichenhao ~]# cat /proc/meminfo  | grep Commit && free -k
CommitLimit:     7737936 kB
Committed_AS:     800516 kB
              total        used        free      shared  buff/cache   available
Mem:        3868968      282872     3196000        8748      390096     3368580
Swap:             0           0           0
```

### vm.panic_on_oom = 0

触发 oom 之后，内核 panic 机制。`默认值：0`。

- 0 # oom 触发后，内核不会 panic。
- 1 # oom 触发后，内核会出现 panic 情况。但是，如果某个进程通过内存/ cpusets 限制使用节点，并且这些节点成为内存耗尽状态，则一个进程可能会被 oom-killer 杀死。在这种情况下不内核不会 panic。因为其他节点的内存可能是空闲的。这意味着系统总体状态可能还不是致命的。
- 2 # oom 触发后，内核直接 panic。

### vm.oom_kill_allocating_task = 0

触发 oom 后，内核 kill 进程的行为。`默认值：0`。

- 0 # 内核将检查每个进程的分数，分数最高的进程将被 kill 掉
- 1 # 那么内核将 kill 掉当前申请内存的进程
