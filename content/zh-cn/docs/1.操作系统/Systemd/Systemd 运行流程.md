---
title: Systemd 运行流程
---

# 概述

> 参考：

# Systemd 运行

这里以 CentOS 7 版本为例

## 确认系统运行级别

systemd 执行的第一个目标是 /etc/systemd/system/default.target，是一个软链接，该文件决定了老版本称为 “运行级别” 的一种行为

```bash
~]# ll /etc/systemd/system/default.target
lrwxrwxrwx. 1 root root 37 Oct 10  2020 /etc/systemd/system/default.target -> /lib/systemd/system/multi-user.target
```

如果想要更改系统启动级别，可以使用 systemctl set-defult XXXXX 命令来修改默认启动级别

## 启动 multi-user.taget

multi-user.target 文件内容如下：

```bash
~]# cat /usr/lib/systemd/system/multi-user.target
[Unit]
Description=Multi-User System
Documentation=man:systemd.special(7)
Requires=basic.target
Conflicts=rescue.service rescue.target
After=basic.target rescue.service rescue.target
AllowIsolate=yes
```

从 multi-user.target 中获取到下一步需要启动的服务。

- 根据 Requires 指令，需要先启动 basic.target 关联的服务
- 启动 /usr/lib/systemd/system/multi-user.target.wants/ 和 /etc/systemd/system/multi-user.target.wants/ 目录中的服务

```bash
~]# ls  /usr/lib/systemd/system/multi-user.target.wants/
dbus.service  getty.target  plymouth-quit.service  plymouth-quit-wait.service  systemd-ask-password-wall.path  systemd-logind.service  systemd-update-utmp-runlevel.service  systemd-user-sessions.service
~]# ls /etc/systemd/system/multi-user.target.wants/
auditd.service  chronyd.service  crond.service  firewalld.service  irqbalance.service  kdump.service  NetworkManager.service  remote-fs.target  rhel-configure.service  rsyslog.service  sshd.service  sysstat.service  tuned.service
```

## 启动 basic.target

basic.target 文件内容如下：

```bash
~]# cat /usr/lib/systemd/system/basic.target
[Unit]
Description=Basic System
Documentation=man:systemd.special(7)

Requires=sysinit.target
After=sysinit.target
Wants=sockets.target timers.target paths.target slices.target
After=sockets.target paths.target slices.target
```

从 basic.target 中获取下一步需要启动的服务。

- 根据 Requires 指令，需要先启动 sysinit.target 关联的服务。
- 根据 Wants 和 After 指令，启动 sockets.target、timers.target、paths.target、slices.target 下关联的服务
- 启动 /etc/systemd/system/basic.target.wants/ 和 /usr/lib/systemd/system/basic.target.wants/ 目录中的服务。

```bash
~]# ls /etc/systemd/system/basic.target.wants/
microcode.service  rhel-dmesg.service
~]# ls /usr/lib/systemd/system/basic.target.wants/
selinux-policy-migrate-local-changes@targeted.service
```

## 启动 sysinit.target

sysinit.target 文件内容如下：

```bash
~]# cat /usr/lib/systemd/system/sysinit.target
[Unit]
Description=System Initialization
Documentation=man:systemd.special(7)
Conflicts=emergency.service emergency.target
Wants=local-fs.target swap.target
After=local-fs.target swap.target emergency.service emergency.target
```

此时没有 Requires 指令了，也就是从该 Unit 开始，同时启动上述所有已经关联到的服务。sysinit.target 会启动重要的系统服务例如系统挂载，内存交换空间和设备，内核补充选项等等。

sysinit.target 将会启动如下服务

- loacl-fs.target 关联服务
- swap.target 关联服务
- emergency.service 服务 和 emergency.target 关联服务
- /etc/systemd/system/sysinit.target.wants/ 与 /usr/lib/systemd/system/sysinit.target.wants/ 目录中的服务

```bash
~]# ls /etc/systemd/system/sysinit.target.wants/
lvm2-lvmetad.socket  lvm2-lvmpolld.socket  lvm2-monitor.service  rhel-autorelabel-mark.service  rhel-autorelabel.service  rhel-domainname.service  rhel-import-state.service  rhel-loadmodules.service
~]# ls /usr/lib/systemd/system/sysinit.target.wants/
cryptsetup.target          plymouth-read-write.service        sys-kernel-config.mount            systemd-firstboot.service               systemd-journal-flush.service      systemd-sysctl.service              systemd-udev-trigger.service
dev-hugepages.mount        plymouth-start.service             sys-kernel-debug.mount             systemd-hwdb-update.service             systemd-machine-id-commit.service  systemd-tmpfiles-setup-dev.service  systemd-update-done.service
dev-mqueue.mount           proc-sys-fs-binfmt_misc.automount  systemd-ask-password-console.path  systemd-journal-catalog-update.service  systemd-modules-load.service       systemd-tmpfiles-setup.service      systemd-update-utmp.service
kmod-static-nodes.service  sys-fs-fuse-connections.mount      systemd-binfmt.service             systemd-journald.service                systemd-random-seed.service        systemd-udevd.service               systemd-vconsole-setup.service
```

## 启动 local-fs.target 与 swap.target 关联服务

local-fs.target 文件内容如下：

```bash
~]# cat /usr/lib/systemd/system/local-fs.target
[Unit]
Description=Local File Systems
Documentation=man:systemd.special(7)
DefaultDependencies=no
Conflicts=shutdown.target
After=local-fs-pre.target
OnFailure=emergency.target
OnFailureJobMode=replace-irreversibly
```

local-fs.target 不会启动用户相关服务，它只处理底层核心服务，这个 target 会根据 /etc/fstab 来执行相关磁盘挂载操作。它通过如下一个目录决定哪些 Unit 会被启动。

- /usr/lib/systemd/system/local-fs.target.wants/

swap.target 文件内容如下：

```bash
~]# cat /usr/lib/systemd/system/swap.target
[Unit]
Description=Swap
Documentation=man:systemd.special(7)
```

# 总结

虽然 systemd 的引用 target 的顺序如上，但是真正的启动顺序为从下到上，所以可以通过设置 default.target 文件来确定开机后默认的登录级别。其中管理单元可以并行启动，从而使效率大大提高。同时 Systemd 是向下兼容 System V 的。

具体顺序应该如下：

- **local-s.target + swap.target **# 这两个 target 主要在挂载本机 /etc/fstab 里面所规范的文件系统与相关的内存交换空间。
- **sysinit.target** # 这个 target 主要在侦测硬件,载入所需要的核心模块等动作。核心所侦测到的各硬件设备，会被记录在 /proc/ 与 /sys/ 目录中，内核参数的修改详见 sysctl 命令。该 target 包括但不限于以下 Unit，详见/usr/lib/systemd/system/sysinit.target.wants/目录
    - 特殊文件系统装置的挂载：包括 dev-hugepages.mount、dev-mqueue.mount 等，主要在挂载跟巨量内存分页使用与消息队列的功能。成功后，会在/dev/目录下简历/dev/hugepages/、/dev/mqueue/等目录
    - 特殊文件系统的启动：包括磁盘阵列、网络驱动器(iscsi)、LVM 文件系统、文件系统对照服务(multipath)等等
    - 开机过程的讯息传递与动画执行：使用 plymouthd 服务搭配 plymouth 指令来传递动画与讯息
    - 日志式登录文件的使用：systemd-journald
    - 加载额外的内核模块：通过 /etc/modules-load.d/\*.conf 文件的设定，让内核额外加载管理员所需要的内核模块
    - 加载额外的内核参数设定：包括 /etc/sysctl.conf 以及 /etc/sysctl.d/\*.conf 内的设定
    - 启动系统的随机数生成器：随机数生成器可以帮助系统进行一些密码加密演算的功能
    - 设定终端(console)字形
    - 启动动态设备管理器：udevd。用来动态对应实际设备存取与设备文件名对应的一个服务
- **basic.target** # 载入主要的周边硬件驱动程序与防火墙相关任务。该 target 包括但不限于以下 Unit，详见/usr/lib/systemd/system/basic.target.wants/目录
    - 加载 alsa 音效驱动程序：这个 alsa 是个音效相关的驱动程序，会系统产生音效
    - 载入 firewalld 防火墙
    - 加载 CPU 微指令功能
    - 启动与设定 SELinux 的安全文本
    - 将目前的开机过程所产生的开机信息写入到/var/log/dmesg 当中
    - 由/etc/sysconfig/module/\*.module 以及/etc/rcmodules 载入管理员指定的模块
    - 加载 systemd 支持的 timer 功能
- **multi-user.target** # 下面的其它一般系统或网络服务的载入。在加载核心驱动硬件后,经过 sysinit.target 的初始化流程让系统可以存取之后,加上 basic.target 让系统成为操作系统的基础, 之后就是服务器要顺利运作时,需要的各种主机服务以及提供服务器功能的网络服务的启动了。这些服务的启动则大多是附挂在 multi-user.target 这个操作环境底下, 可以到 /etc/systemd/system/multi-user.target.wants/ 里头去瞧瞧预设要被启动的服务。针对主机的本地服务与网络服务的各项 Unit 若要 enable 的话，就是将该 Unit 放到这个目录下做个软链接。该 target 包括但不限于以下 Unit，详见 /usr/lib/systemd/system/multi-user.target.wants/ 目录
    - 相容 systemV 的 rc-loacl.service，开机自动执行的命令
    - 提供 tty(终端)界面与登录的服务

## 查看启动顺序的时间

- 要查看具体的启动顺序可以通过如下命令输入到文件，然后通过浏览器打开查看。
    - systemd-analyze plot > boot.html
- 列出所有正在运行的单元，按从初始化开始到启动所花的时间排序。
    - systemd-analyze blame

也就是说，如果想让一个 systemd 的系统正常运行，则通过 default.target 来一步一步决定运行那些 Unit，最后从决定的末尾开始，一步一步启动各个 Unit，如图所示：

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/az9p3g/1616634160644-69ffd65b-b9c4-490b-aad8-77e8e218bb02.jpeg)
