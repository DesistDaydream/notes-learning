---
title: Audit
weight: 2
---

# 概述

> 参考：
>
> - [GitHub 组织，linux-audit](https://github.com/linux-audit)
>   - [GitHub 项目，linux-audit/audit-kernel](https://github.com/linux-audit/audit-kernel)
>   - [GitHub 项目，linux-audit/audit-userspace](https://github.com/linux-audit/audit-userspace)
> - [红帽产品文档，RedHat7 - 安全指南 - 系统审计](https://access.redhat.com/documentation/en-us/red_hat_enterprise_linux/7/html/security_guide/chap-system_auditing)
> - [红帽产品文档，RedHat9 - 安全强化 - 系统审计](https://access.redhat.com/documentation/en-us/red_hat_enterprise_linux/9/html/security_hardening/auditing-the-system_security-hardening)
> - [公众号 - kernsec，Linux Audit 子系统解读](https://mp.weixin.qq.com/s/G6kE52o7OZaGYPqnuUwggQ)
> - [linux audit审计（8）--开启audit对系统性能的影响](https://www.cnblogs.com/xingmuxin/p/8875783.html)

Audit 是实现的 Linux 系统审计的软件包，其中包含两个部分：

- 用户空间应用程序和实用程序
- 内核端系统调用处理。

当我们安装 Audit 后，会在系统中看到两个进程，一个是内核态的，一个是用户态的。

```bash
~]# ps -p $(pgrep audit)
  PID TTY      STAT   TIME COMMAND
  100 ?        S      0:02 [kauditd]
15016 ?        S<sl   0:14 /sbin/auditd
```

并且，不能随意停止 auditd 服务器，内核进程无法通过用户空间的操作终止，使用 `systemctl stop auditd.service` 将会报错：

```text
Failed to stop auditd.service: Operation refused, unit auditd.service may be requested by dependency only (it is configured to refuse manual start/stop).
See system logs and 'systemctl status auditd.service' for details.
```

内核组件从用户空间应用程序接收系统调用，并通过以下过滤器之一对其进行过滤：user, task, fstype, exit

系统调用通过 exclude 过滤器后，它将通过上述其中一个过滤器发送，这些过滤器根据 Audit 规则配置将其发送到 Audit 守护进程，以进行进一步处理。

用户空间审计守护进程从内核收集信息，并在日志文件中创建条目。其他 Audit 用户空间实用程序与 Audit 守护进程、内核审计组件或 Audit 日志文件交互：

- audisp - Audit 分配程序守护进程与 Audit 守护进程交互，并将事件发送到其他应用以进行进一步处理。此守护进程的目的是提供一种插件机制，让实时分析程序能够与审计事件交互。
- auditctl - Audit 控制实用程序与内核审计组件交互，以管理规则并控制事件生成进程的许多设置和参数。
- 剩余的 Audit 实用程序将 Audit 日志文件的内容作为输入，并根据用户的要求生成输出。例如，aureport 实用程序生成所有记录事件的报告。

**注意：由于 Audit 会运行一个内核态的进程，并且监听系统调用，所以像 nfs 这种涉及内核的文件系统工具，可能会跟 Audit 产生冲突。**

# 安装 Audit

```bash
~]# yum install audit audit-libs
```

# Audit 关联文件与配置

**/etc/audit/** #

- **./auditd.conf** # auditd 进程运行时配置
- **./audit.rules** # 由 `rules.d/` 目录下的规则文件生成
- **./rules.d/** # Audit 规则文件

**/usr/share/doc/audit-${VERSION}/rules/** # Audit 软件包安装后根据各种认证标准提供一组预配置的规则文件

**/var/log/audit/** # auditd 记录日志的默认位置。

# 规则语法

## 定义文件系统规则

**-w FILE -p PERMISSIONS -k KEY**

- **FILE** # 是被审计的文件或目录
- **PERMISSIONS** # 记录的权限
  - r - 对文件或目录的读取访问权限.
  - w - 对文件或目录的写入访问权限.
  - x - 执行对文件或目录的访问权限。
  - a - 更改文件或目录的属性.
- **KEY** # 是一个可选字符串，可帮助您识别生成了特定日志条目的规则或一组规则。

## 定义系统调用规则

**-a action,filter -S system_call -F field=value -k key_name**

- **action 和 filter** # 指定记录特定事件的时间。操作可以是 always 或 never。filter 指定将哪个内核规则匹配过滤器应用到事件。rule-matching 过滤器可以是以下之一： task、exit、user 和 exclude。有关这些过滤器的更多信息，请参阅 [第  7.1  节 “Audit 系统架构”](https://access.redhat.com/documentation/zh-cn/red_hat_enterprise_linux/7/html/security_guide/chap-system_auditing#sec-audit_system_architecture) 的开头。
- **system_call** # 指定系统调用的名称。可以在 /usr/include/asm/unistd_64.h 文件中找到所有系统调用的列表。可将多个系统调用分组成一个规则，各自在其自己的 -S 选项后指定。
- **Field=value** # 指定进一步修改规则以根据指定的体系结构、组 ID、进程 ID 和其他选项匹配的额外选项。有关所有可用字段类型及其值的完整列表，请查看 auditctl(8) man page。
- **key_name** # 是一个可选字符串，可帮助您识别生成了特定日志条目的规则或一组规则。

# auditctl 命令行工具

auditctl 命令允许您控制 Audit 系统的基本功能，并定义决定记录哪些审计事件的规则。

# aureport 命令行工具

aureport 实用程序允许您针对 Audit 日志文件中记录的事件生成摘要和列ar 报告

# ausearch 命令行工具

ausearch 实用程序允许您搜索 Audit 日志文件特定事件。默认情况下，ausearch 搜索 `/var/log/audit/audit.log` 文件。
["Linux libc 库"](/docs/1.操作系统/Linux 源码解析/Linux libc 库.md)
