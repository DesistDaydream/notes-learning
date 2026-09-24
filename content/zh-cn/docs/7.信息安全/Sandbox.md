---
title: "Sandbox"
created: "2026-09-24T13:09"
weight: 100
---

# 概述

> 参考：
>
> - [Wiki, Sandbox(computer security)](https://en.wikipedia.org/wiki/Sandbox_(computer_security))
> - [Wiki, Sandbox(software development)](https://en.wikipedia.org/wiki/Sandbox_(software_development))

**Sandbox(沙箱)** 是一类**受限执行环境**：让不完全可信的代码在受控边界内运行，即使它被攻陷或被滥用，能造成的影响也被限制在边界之内。

关键的是：**沙箱不是某一种具体技术，而是一类抽象**。它由两个动作组成——**隔离**（让里面的东西看不见外面的东西）与**限制**（即使看得见也不许动）。不同的实现只是在这两个动作上各管一段。

也正因为如此，Chroot、Namespaces、seccomp、Landlock、bubblewrap、Containerization、Virtualization、etc. 这些东西**并不是同一条演化链上的前后代**，而是在不同层次上各解决一部分问题的并列机制。

# 实现族谱

## 视图隔离（换根）

- [Chroot](/docs/1.操作系统/Kernel/Process/Chroot.md) # 更改进程及其子进程看到的 `/`。最早、最弱的一种；**root 可逃逸**，所以它是「视图切换」而不是安全边界

## 内核命名空间

- [Namespaces](/docs/10.云原生/Containerization/1.Namespaces/1.Namespaces.md) # Linux 内核提供的隔离原语：mount / pid / net / ipc / uts / user / cgroup / time。容器、bubblewrap、nsjail 都是把它们组合起来用

## 资源限制

- [CGroup](/docs/10.云原生/Containerization/2.CGroup/2.CGroup.md) # CPU / 内存 / IO 等资源的上限与统计。注意它**不隔离可见性**，只管用量

## 访问控制与能力

- [Access Control](/docs/1.操作系统/登录%20Linux%20与%20访问控制/Access%20Control/Access%20Control.md) # 传统 DAC、Linux Capabilities、SELinux / AppArmor 等。这是「限制」这一半最正统的一支

## 系统调用过滤

- **seccomp-bpf** # 用 BPF 程序过滤 syscall。粒度细、开销低，但**不能做基于路径的文件系统控制**——它只看得到 fd 和 syscall 号
- **Landlock** # 内核 LSM，做**路径级**的文件访问限制（以及 ABI v4+ 的网络端口、IPC）。非特权即可用，不需要 user namespace；随内核 ABI 版本演进，旧内核可能只能覆盖访问类别的子集

## 用户态沙箱运行时

- [Bubblewrap](#Bubblewrap) # 用非特权 user namespace 拼出 mount + PID 等命名空间，再 `execve`。Flatpak 的底座。**不提供策略**，参数即策略
- **nsjail/firejail** # 同类工具，各自把 namespace + seccomp + cgroup + profile 打包成更好用的 CLI

## 完整环境隔离

- [Containerization](/docs/10.云原生/Containerization/Containerization.md) # namespace + cgroup + 联合文件系统 + 镜像/OCI 规范。隔离的是**整个能力环境**，不只是某个进程
- [Virtualization](/docs/10.云原生/Virtualization/Virtualization.md) # 换掉整个内核（含 microVM）。隔离强度最高，代价是启动与资源开销

## 服务级组合

- [systemd.exec 类指令](/docs/1.操作系统/Systemd/Unit%20File/systemd.exec%20类指令.md) # 其中的 `SANDBOXING(沙盒)` 章节：`ProtectSystem=`、`ProtectHome=`、`PrivateTmp=`、`DynamicUser=` 等，由 PID 1 帮你把上面这些机制组合好

## 其他语境里的「沙箱」

「沙箱」这个词在别的领域也大量出现，含义都是「受限的执行环境」，但机制完全不同，别混为一谈：

- **eBPF 的 Sandbox Programs** # 见 [eBPF](/docs/1.操作系统/Kernel/BPF/eBPF.md)——指在内核中受限运行的 BPF 程序
- **Web 浏览器沙箱** # 如 Chrome 的 renderer sandbox，防止网页代码逃逸到系统
- **Kubernetes 的 PodSandbox** # CRI 术语，对应 Pod 的 infra(pause) 容器所维持的那组 namespace
- **CNCF Sandbox** # 项目成熟度分级，与安全无关
- **Mock / 测试沙盒环境** # 指与生产隔离的试验环境
- **Agent 沙箱** # 见 [Agent](/docs/12.AI/Agent.md) — 约束 LLM Agent 执行命令时的文件影响，把「不受限的 shell」收进一个可审批的边界

# 一份粗略对照

| 机制 | 层次 | 需要特权/userns | 主要管的维度 |
| --- | --- | --- | --- |
| chroot | 用户态 syscall | 需 `CAP_SYS_CHROOT` | 文件系统视图（可逃逸） |
| Namespaces | 内核 | 视类型而定 | 视图：fs / pid / net / ipc / uts / user |
| CGroup | 内核 | 通常需特权 | 资源用量 |
| Capabilities / SELinux / AppArmor | 内核 LSM | 需策略配置 | 凭据、访问控制 |
| seccomp-bpf | 内核 | 非特权（`no_new_privs`） | syscall |
| Landlock | 内核 LSM | 非特权 | 路径级文件访问 |
| bubblewrap | 用户态 + namespace | 非特权 userns | 文件系统布局 + 可选 pid/net 等 |
| 容器(runc 等) | 用户态 + namespace | 默认 root，rootless 需 userns | 完整环境 |
| microVM | 硬件虚拟化 | 需 KVM 等 | 完整机器 |

# 常见误区

- **「进了 chroot 就安全了」** # 不是。chroot 只换视图，持有特权的进程可以逃出，历史上的容器逃逸 CVE 多与此有关
- **「挂了新的 `/proc` 就隔离了」** # 不是。procfs 只显示**执行 mount 的那个进程所在 PID Namespace 内**的进程；不建私有 PID Namespace 的话，`/proc/<pid>/root`、`/proc/<pid>/fd` 这类魔法链接可以穿透挂载白名单。**必须两者配套**
- **「沙箱就是容器」** # 容器是沙箱的一种**实现形态**，且不含虚拟化；反过来沙箱可以是进程内的一层路径检查（如进程内 fs 围栏），完全不涉及内核边界
- **「隔离强度」不能跨维度比较** # 一个能限制 `open()` 路径的机制，和一个能卸载整个网络栈的机制，没有谁「更强」——它们管的维度不同

# Bubblewrap

> 参考：
>
> - [GitHub 项目，containers/bubblewrap](https://github.com/containers/bubblewrap)

最初的 bubblewrap 代码存在于用户命名空间之前 - 它继承了 xdg-app helper 的代码，而 又远远地源自 [linux-user-chroot](https://git.gnome.org/browse/linux-user-chroot)

有些 AI [Agent](/docs/12.AI/Agent.md) （e.g. Deepseek Harness、etc.）在 Linux 中运行，会调用 bubblewrap 开启沙箱来运行某些命令。
