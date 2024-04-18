---
title: Unit File 指令
---

# 概述

> 参考：
>
> - [Manual(手册)，systemd.unit(5)](https://man7.org/linux/man-pages/man5/systemd.unit.5.html) # Unit 的介绍
> - [Manual(手册)，systemd.directives(7)](https://man7.org/linux/man-pages/man7/systemd.directives.7.html) # Unit File 中全部的指令列表

一个 Unit File 具有多个 Sections(部分)，大体可以分为 2 类

- **通用 Sections** # 与 Unit 类型无关的部分
  - \[Unit] 与 \[Install]
- **特殊 Sections** # 特定于 Unit 类型的部分
  - \[Service]、\[Socket]、\[Timer]、\[Mount]、\[Path] 等等

除了 Unit 和 Install 以外的其余每个 **Sections(部分)** 都有其特定的 **Directives(指令)**，同时，也有一些通用的 Directives(指令) 可以用在多个 Sections(部分) 中。

# 通用部分的指令

## \[Unit] 部分的指令

https://man7.org/linux/man-pages/man5/systemd.unit.5.html#[UNIT]_SECTIONOPTIONS

unit 本身的说明，以及与其他相依赖的 daemon 的设置，包括在什么服务之前或之后启动等设置

**Description=\<STRING>** # Unit 描述，用 systemctl list-units 和 systemctl status 查看服务时候的描述内容就是这里定义的

**Documentation=\<STRING>** # 提供该 Unit 可以进一步文件查询的地址或者位置

**After(Before)=\<STRING>** # 在哪些之后(前)启动，说明该 Unit 可以在哪些 daemon 启动后(前)才能够启动，非强制性，只是推荐规范

**Requires=\<STRING>** # 需要启动哪些，说明启动该 Unit 前需要启动哪些 Unit，强制性的，如果不启动该项定义的 Unit 则无法启动该 Unit

**Wants=\<STRING>** # 想要启动哪些，与 Requires 相反，说明启动该 Unit 后想要启动哪些 Unit，非强制

**Conflicts=\<STRING>** # 代表该 Unit 与列表中的 daemon 有冲突，如果该设置里的服务启动了，那么这个 Unit 就不能启动

## \[Install] 部分的指令

https://man7.org/linux/man-pages/man5/systemd.unit.5.html#[INSTALL]_SECTION_OPTIONS

Install 部分包含 Unit 的启动信息。通常是配置文件的最后一个区块，用来定义如何启动，以及是否开机启动等等。

**WantedBy=\<STRING>** # 它的值是一个或多个 Target，当前 Unit 激活时（enable）符号链接会放入/etc/systemd/system/目录下面以 Target 名 + .wants 后缀构成的子目录中

**RequiredBy=\<STRING>** # 它的值是一个或多个 Target，当前 Unit 激活时，符号链接会放入/etc/systemd/system 目录下面以 Target 名 + .required 后缀构成的子目录中

**Alias=\<STRING>** # 当前 Unit 可用于启动的别名

**Also=\<STRING>** # 当前 Unit 激活（enable）时，会被同时激活的其他 Unit

# 特殊部分的指令

不同的 Unit 类型就使用对应的部分，在这里面设定 启动程序的命令、环境配置、重启方式 等等。

每个 **特殊 Sections(部分)** 都有其**特定的 Directives(指令)**，同时，也有一些通用的 Directives(指令) 可以用在多个 特殊 Sections(部分) 中。

## 通用指令

> TODO: 在哪里可以找到官方堆这些特殊指令类型的统一描述？~我这里是根据具体的指令找到的执行类型，比如 我想看 KillMode 指令，就找到了 systemd.kill 类指令。

这些通用指令可以用在多个特殊 Sections 中。

### systemd.exec 类指令

详见 [systemd.exec 类指令](/docs/1.操作系统/Systemd/Unit%20File/systemd.exec%20类指令.md)

### systemd.kill 类指令

https://man7.org/linux/man-pages/man5/systemd.kill.5.html

systemd.kill 类的指令用于配置进程停止时，应该使用方式方法。

- 用于 service、socket、swap、mount、scope 部分

**KillMode=\<STRING>** # 这个选项用来处理 Containerd 进程被杀死的方式。默认情况下，systemd 会在进程的 cgroup 中查找并杀死 Containerd 的所有子进程，这肯定不是我们想要的。`KillMode`字段可以设置的值如下。我们需要将 KillMode 的值设置为 `process`，这样可以确保升级或重启 Containerd 时不杀死现有的容器。

- **control-group**（默认值）# 当前控制组里面的所有子进程，都会被杀掉
- **process** # 只杀主进程
- **mixed** # 主进程将收到 SIGTERM 信号，子进程收到 SIGKILL 信号
- **none** # 没有进程会被杀掉，只是执行服务的 stop 命令。

### systemd.resource-control 类指令

https://man7.org/linux/man-pages/man5/systemd.resource-control.5.html

用于对 Unit 启动的进程进行资源限制相关的指令

**Delegate=\<STRING>** # 这个选项允许进程(比如 containerd)以及运行时自己管理自己创建的容器的 `cgroups`。如果不设置这个选项，systemd 就会将进程移到自己的 `cgroups` 中，从而导致该进程无法正确获取容器的资源使用情况。

## \[Mount] 部分的指令

TODO: 待整理

## \[Service] 部分的指令

https://man7.org/linux/man-pages/man5/systemd.service.5.html

详见 [service Unit](/docs/1.操作系统/Systemd/Unit%20File/service%20Unit.md)

## \[timer] 部分指令

https://man7.org/linux/man-pages/man5/systemd.timer.5.html

详见 [timer Unit](/docs/1.操作系统/Systemd/Unit%20File/timer%20Unit.md)
