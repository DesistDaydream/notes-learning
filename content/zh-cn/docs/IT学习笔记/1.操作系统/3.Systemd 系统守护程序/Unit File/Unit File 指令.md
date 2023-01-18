---
title: Unit File 指令
---

# 概述

> 参考：
> - [Manual(手册),systemd.directives(7)](https://man7.org/linux/man-pages/man7/systemd.directives.7.html) # Unit File 中全部的指令列表
> -

一个 Unit File 具有多个 Sections，大体可以分为 2 类

- 通用 Sections # 与 Unit 类型无关的部分
  - \[Unit] 与 \[Install]
- 特殊 Sections # 特定于 Unit 类型的部分
  - \[Service]、\[Socket]、\[Timer]、\[Mount]、\[Path] 等等

除了 Unit 和 Install 以外的其余每个 Sections(部分) 都有其特定的 Directives(指令)，同时，也有一些通用的 Directives(指令) 可以用在多个 Sections(部分) 中。

# 通用部分的指令

## [\[Unit\]](https://man7.org/linux/man-pages/man5/systemd.unit.5.html#[UNIT]_SECTION_OPTIONS) 部分的指令

unit 本身的说明，以及与其他相依赖的 daemon 的设置，包括在什么服务之前或之后启动等设置
**Description=<STRING>** # Unit 描述，用 systemctl list-units 和 systemctl status 查看服务时候的描述内容就是这里定义的
**Documentation=<STRING>** # 提供该 Unit 可以进一步文件查询的地址或者位置
**After(Before)=<STRING>** # 在哪些之后(前)启动，说明该 Unit 可以在哪些 daemon 启动后(前)才能够启动，非强制性，只是推荐规范
**Requires=<STRING>** # 需要启动哪些，说明启动该 Unit 前需要启动哪些 Unit，强制性的，如果不启动该项定义的 Unit 则无法启动该 Unit
**Wants=<STRING>** # 想要启动哪些，与 Requires 相反，说明启动该 Unit 后想要启动哪些 Unit，非强制
**Conflicts=<STRING>** # 代表该 Unit 与列表中的 daemon 有冲突，如果该设置里的服务启动了，那么这个 Unit 就不能启动

## [\[Install\]](https://man7.org/linux/man-pages/man5/systemd.unit.5.html#[INSTALL]_SECTION_OPTIONS) 部分的指令

Install 部分包含 Unit 的启动信息。通常是配置文件的最后一个区块，用来定义如何启动，以及是否开机启动等等。
**WantedBy=<STRING>** # 它的值是一个或多个 Target，当前 Unit 激活时（enable）符号链接会放入/etc/systemd/system/目录下面以 Target 名 + .wants 后缀构成的子目录中
**RequiredBy=<STRING>** # 它的值是一个或多个 Target，当前 Unit 激活时，符号链接会放入/etc/systemd/system 目录下面以 Target 名 + .required 后缀构成的子目录中
**Alias=<STRING>** # 当前 Unit 可用于启动的别名
**Also=<STRING>** # 当前 Unit 激活（enable）时，会被同时激活的其他 Unit

# 特殊部分的指令

不同的 Unit 类型就使用对应的部分，在这里面设定 启动程序的命令、环境配置、重启方式 等等。

每个 特殊 Sections(部分) 都有其特定的 Directives(指令)，同时，也有一些通用的 Directives(指令) 可以用在多个 特殊 Sections(部分) 中。

## 通用指令

### [systemd.exec](https://man7.org/linux/man-pages/man5/systemd.exec.5.html) 类指令

systemd.exec 类的指令用于配置进程执行时的环境，比如 环境变量、运行程序的用户和组、运行路径 等等

- 用于 service、socket、mount、swap 部分

[**PATHS(路径)**](https://man7.org/linux/man-pages/man5/systemd.exec.5.html#PATHS)** 相关指令**

- **WorkingDirectory=<STRING>** # 采用相对于由 RootDirectory 指令 或特殊值 `~` 指定的服务根目录的目录路径。

[**USER/GROUP IDENTITY(用户/组标识)**](https://man7.org/linux/man-pages/man5/systemd.exec.5.html#USER/GROUP_IDENTITY)** 相关指令**

- **User=<STRING>** # 指定运行该 Unit 使用的用户。

[**CAPABILITIES(能力)**](https://man7.org/linux/man-pages/man5/systemd.exec.5.html#CAPABILITIES)** 相关指令**

[**SECURITY(安全)**](https://man7.org/linux/man-pages/man5/systemd.exec.5.html#SECURITY)** 相关指令**

[**MANDATORY ACCESS CONTROL(强制访问控制)**](https://man7.org/linux/man-pages/man5/systemd.exec.5.html#MANDATORY_ACCESS_CONTROL)** 相关指令**

[**PROCESS PROPERITES(进程属性)**](https://man7.org/linux/man-pages/man5/systemd.exec.5.html#PROCESS_PROPERTIES)** 相关指令**
为执行的进程设置各种资源的软限制和硬限制。

[**SCHEDULING(调度)**](https://man7.org/linux/man-pages/man5/systemd.exec.5.html#SCHEDULING)** 相关指令**

[**SANDBOXING(沙盒)**](https://man7.org/linux/man-pages/man5/systemd.exec.5.html#SANDBOXING)** 相关指令**

[**SYSTEM CALL FILTERING(系统调用过滤)**](https://man7.org/linux/man-pages/man5/systemd.exec.5.html#SYSTEM_CALL_FILTERING)** 相关指令**

[**ENVIRONMENT(环境变量)**](https://man7.org/linux/man-pages/man5/systemd.exec.5.html#ENVIRONMENT)** 相关指令**
**Environment=<STRING>** # 指定 Unit 启动时所使用的环境变量

[**LOGGING AND STANDARD INPUT/OUTPUT(日志的标准输入/输出)**](https://man7.org/linux/man-pages/man5/systemd.exec.5.html#LOGGING_AND_STANDARD_INPUT/OUTPUT)** 相关指令**

### [systemd.kill](https://man7.org/linux/man-pages/man5/systemd.kill.5.html) 类指令

systemd.kill 类的指令用于配置进程停止时，应该使用方式方法。

- 用于 service、socket、swap、mount、scope 部分

**KillMode=<STRING>** # 这个选项用来处理 Containerd 进程被杀死的方式。默认情况下，systemd 会在进程的 cgroup 中查找并杀死 Containerd 的所有子进程，这肯定不是我们想要的。`KillMode`字段可以设置的值如下。我们需要将 KillMode 的值设置为 `process`，这样可以确保升级或重启 Containerd 时不杀死现有的容器。

- **control-group**（默认值）：当前控制组里面的所有子进程，都会被杀掉
- **process**：只杀主进程
- **mixed**：主进程将收到 SIGTERM 信号，子进程收到 SIGKILL 信号
- **none**：没有进程会被杀掉，只是执行服务的 stop 命令。

### [systemd.resource-control](https://man7.org/linux/man-pages/man5/systemd.resource-control.5.html) 类指令

用于对 Unit 启动的进程进行资源限制相关的指令
**Delegate=<STRING>** # 这个选项允许进程(比如 containerd)以及运行时自己管理自己创建的容器的 `cgroups`。如果不设置这个选项，systemd 就会将进程移到自己的 `cgroups` 中，从而导致该进程无法正确获取容器的资源使用情况。

## [\[Service\]](https://man7.org/linux/man-pages/man5/systemd.service.5.html) 部分的指令

详见 [service Unit](/docs/IT学习笔记/1.操作系统/3.Systemd%20 系统守护程序/Unit%20File/service%20Unit.md Unit.md)

## [\[timer\]](https://man7.org/linux/man-pages/man5/systemd.timer.5.html) 部分指令

详见 [timer Unit](/docs/IT学习笔记/1.操作系统/3.Systemd%20 系统守护程序/Unit%20File/timer%20Unit.md Unit.md)
