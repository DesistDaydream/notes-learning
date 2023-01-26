---
title: CGroup FS
---

# 概述

> 参考：

# /sys/fs/cgroup/\*

### CGroupV1

CGroupV1 根目录下的每个目录的名称都是一个子系统的名称，每个子系统都有其自己独立的资源控制配置文件。

```bash
~]# ls -l /sys/fs/cgroup/
total 0
dr-xr-xr-x 5 root root  0 Jan 26 21:46 blkio
lrwxrwxrwx 1 root root 11 Jan 26 21:46 cpu -> cpu,cpuacct
lrwxrwxrwx 1 root root 11 Jan 26 21:46 cpuacct -> cpu,cpuacct
dr-xr-xr-x 5 root root  0 Jan 26 21:46 cpu,cpuacct
dr-xr-xr-x 3 root root  0 Jan 26 21:46 cpuset
dr-xr-xr-x 5 root root  0 Jan 26 21:46 devices
dr-xr-xr-x 4 root root  0 Jan 26 21:46 freezer
dr-xr-xr-x 3 root root  0 Jan 26 21:46 hugetlb
dr-xr-xr-x 5 root root  0 Jan 26 21:46 memory
lrwxrwxrwx 1 root root 16 Jan 26 21:46 net_cls -> net_cls,net_prio
dr-xr-xr-x 3 root root  0 Jan 26 21:46 net_cls,net_prio
lrwxrwxrwx 1 root root 16 Jan 26 21:46 net_prio -> net_cls,net_prio
dr-xr-xr-x 3 root root  0 Jan 26 21:46 perf_event
dr-xr-xr-x 5 root root  0 Jan 26 21:46 pids
dr-xr-xr-x 2 root root  0 Jan 26 21:46 rdma
dr-xr-xr-x 5 root root  0 Jan 26 21:46 systemd
dr-xr-xr-x 5 root root  0 Jan 26 21:46 unified
```

#### ./cpu # CPU 子系统

- ./cpu.cfs_quota_us 与 ./cpu.cfs_period_us # 用来限制进程每运行 cfs_period_us 一段时间，只能被分配到的总量为 cfs_quota_us 的 CPU 时间
  - cfs_quota_us 默认值为-1，不做任何限制，如果修改为 20000(20ms)则表示 CPU 只能使用到 20%的
  - cfs_period_us 默认值为 100000(100ms)
- ./cpu.shares #
- ./cpu.stat #
  - nr_periods #
  - nr_throttled #
  - throttled_time #

### CGroupV2

```bash
~]# ls -l /sys/fs/cgroup/
total 0
-r--r--r--  1 root root 0 Feb 18 10:52 cgroup.controllers
-rw-r--r--  1 root root 0 Feb 18 10:54 cgroup.max.depth
-rw-r--r--  1 root root 0 Feb 18 10:54 cgroup.max.descendants
-rw-r--r--  1 root root 0 Feb 18 10:52 cgroup.procs
-r--r--r--  1 root root 0 Feb 18 10:54 cgroup.stat
-rw-r--r--  1 root root 0 Feb 18 10:52 cgroup.subtree_control
-rw-r--r--  1 root root 0 Feb 18 10:54 cgroup.threads
-rw-r--r--  1 root root 0 Feb 18 10:54 cpu.pressure
-r--r--r--  1 root root 0 Feb 18 10:52 cpuset.cpus.effective
-r--r--r--  1 root root 0 Feb 18 10:52 cpuset.mems.effective
drwxr-xr-x  2 root root 0 Feb 18 10:52 init.scope
-rw-r--r--  1 root root 0 Feb 18 10:54 io.cost.model
-rw-r--r--  1 root root 0 Feb 18 10:54 io.cost.qos
-rw-r--r--  1 root root 0 Feb 18 10:54 io.pressure
-rw-r--r--  1 root root 0 Feb 18 10:54 memory.pressure
drwxr-xr-x 44 root root 0 Feb 18 10:53 system.slice
drwxr-xr-x  3 root root 0 Feb 18 10:53 user.slice
```
