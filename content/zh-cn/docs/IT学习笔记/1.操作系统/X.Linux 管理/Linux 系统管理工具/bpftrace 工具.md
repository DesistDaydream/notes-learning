---
title: bpftrace 工具
---

# 概述

> 参考：
> - [GitHub 项目，iovisor/bpftrace](https://github.com/iovisor/bpftrace)

bpftrace 是用于 Linux 增强型 eBPF 的高级跟踪语言，可在最近的 Linux 内核 (4.x) 中使用。 bpftrace 使用 LLVM 作为后端，将脚本编译为 BPF 字节码，并利用 BCC 与 Linux BPF 系统交互，以及现有的 Linux 跟踪功能：内核动态跟踪（kprobes）、用户级动态跟踪（uprobes）、和跟踪点。 bpftrace 语言的灵感来自 awk 和 C，以及 DTrace 和 SystemTap 等前身跟踪器。 bpftrace 是由 Alastair Robertson 创建的。、

**简单示例**

```bash
# Files opened by process
bpftrace -e 'tracepoint:syscalls:sys_enter_open { printf("%s %s\n", comm, str(args->filename)); }'

# Syscall count by program
bpftrace -e 'tracepoint:raw_syscalls:sys_enter { @[comm] = count(); }'

# Read bytes by process:
bpftrace -e 'tracepoint:syscalls:sys_exit_read /args->ret/ { @[comm] = sum(args->ret); }'

# Read size distribution by process:
bpftrace -e 'tracepoint:syscalls:sys_exit_read { @[comm] = hist(args->ret); }'

# Show per-second syscall rates:
bpftrace -e 'tracepoint:raw_syscalls:sys_enter { @ = count(); } interval:s:1 { print(@); clear(@); }'

# Trace disk size by process
bpftrace -e 'tracepoint:block:block_rq_issue { printf("%d %s %d\n", pid, comm, args->bytes); }'

# Count page faults by process
bpftrace -e 'software:faults:1 { @[comm] = count(); }'

# Count LLC cache misses by process name and PID (uses PMCs):
bpftrace -e 'hardware:cache-misses:1000000 { @[comm, pid] = count(); }'

# Profile user-level stacks at 99 Hertz, for PID 189:
bpftrace -e 'profile:hz:99 /pid == 189/ { @[ustack] = count(); }'

# Files opened, for processes in the root cgroup-v2
bpftrace -e 'tracepoint:syscalls:sys_enter_openat /cgroup == cgroupid("/sys/fs/cgroup/unified/mycg")/ { printf("%s\n", str(args->filename)); }'

```
