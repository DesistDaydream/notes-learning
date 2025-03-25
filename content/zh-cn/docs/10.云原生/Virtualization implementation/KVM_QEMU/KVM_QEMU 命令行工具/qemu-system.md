---
title: qemu-system
linkTitle: qemu-system
weight: 20
---

# 概述

> 参考：
>
> - [官方 Manual(手册)，qemu-system-x86_64](https://www.qemu.org/docs/master/system/qemu-manpage.html)
> - [官方文档，系统模拟-调用](https://www.qemu.org/docs/master/system/invocation.html)（这个其实也是 man 手册）

qemu-system 的名称在不同的 CPU 架构上有不同的名称：

- amd64 架构
  - qemu-system-x64_64
- arm64 架构
  - qemu-system-aarch64
  - qemu-system-arm

注意：在 CentOS 系统中，该二进制文件的名字是 qemu-kvm，是一个在 /usr/local/bin/qemu-kvm 这个位置并指向 /usr/libexec/qemu-kvm 的软链接

# Syntax(语法)

**qemu-system-x86_64 \[OPTIONS] \[DISK_IMAGE]**

DISK_IMAGE 是 IDE 硬盘 0 的原始硬盘映像。有些目标不需要磁盘映像。

## Standard OPTIONS(标准选项)

https://www.qemu.org/docs/master/system/invocation.html#hxtool-0

**-name** #

**-m** # 指定内存大小，单位是 M

**-cpu MDOEL** # 设定要模拟的 CPU 类型

**-smp \[,cores=核心数]\[,threads=线程数]\[,sockets=有几个 CPU 插槽]\[,maxcpus=指定正在使用的 CPU 数]** # 设定 vCPU 数量

**-M** # 指定要模拟的主机类型,可以使用 `qemu-system-x86_64 -M ?` 命令查看所支持的所有可模拟的主机类型

**-boot \[order=DRIVES]\[,once=DRIVES]\[,menu=on|off]** # 定义启动设备的引导次序

**-device DRIVER\[,PROPERTY\[=VALUE]\[,....]]** # 为虚拟机添加指定设备的驱动程序，`PROPERTY=VALUE` 用于设置驱动程序的属性。

> 使用 `qemu-system-x86_64 -device help` 获取 qemu-kvm 可以模拟的所有设备列表
> 使用 `qemu-system-x86_64 -device DRIVER,help` 获取指定 DRIVER(设备) 的驱动程序信息

- Storage Devices(存储设备相关驱动)
- Network Devices(网络设备相关驱动)
  - e1000 # bus PCI, desc "Intel Gigabit Ethernet"
  - rtl8139 # bus PCI
  - virtio-net-device # bus virtio-bus
  - virtio-net-pci # bus PCI, 别名：virtio-net
- Input Devices(输入设备相关驱动)
  - virtio-serial-pci # bus PCI。别名：virtio-serial
  - virtserialport # bus virtio-serial-bus

## Block Device OPTIONS(块设备选项)

https://www.qemu.org/docs/master/system/qemu-manpage.html#hxtool-1

**-cdrom FILE** # 将指定文件作为 VM 的 CD-ROM。可以使用 /dev/cdrom 作为 FILE 让 VM 使用宿主机的 CD-ROM。

**-drive OPTIONS\[,OPTION\[,.......]]** # 定义一个硬盘设备。创建一个 块设备节点 和 VM 设备

> 该选项是新版推荐的，用来代替 -blockdev 和 -device 两个选项，以使操作更简便。
> 并且 -drive 接受 -blockdev 选项的所有选项。

- file=/PATH/FileName # 硬盘映像文件路径
- if=INTERFACE # 指定硬盘设备所连接的接口类型，即控制器类型，如 ide、scsi、sd、mtd、floppy、pflash、virtio、none 等
- media=disk|cdrom # 定义介质类型为硬盘(disk)还是光盘(cdrom)
- snapshot=on|off # 定义是否支持快照功能，on 开启，off 关闭
- cache=CACHE # 定义如何使用物理机缓存来访问块数据，如 none、writeback、writethrough、unsafe
- format=FORMAT # 指定映像文件的格式，具体格式参考 qemu-img

## Display OPTIONS(显示选项)

https://www.qemu.org/docs/master/system/qemu-manpage.html#hxtool-3

**-vnc DISPLAY\[,option\[,option\[,...]]]**

该选项可以让 QEMU 监听一个端口，并通过 VNC 会话重定向 VGA 显示。

> 在使用此选项时启用 usb tablet 设备是非常有用的。(使用 -device usb-tablet)。这样可以让鼠标移动更迅速，否则会出现两个鼠标的情况，就是虚拟机外面一个，虚拟机内部一个，外部移动到哪，内部移动到哪，而且非常慢。

在使用 VNC 显示时，如果不使用 en-us，则必须使用-k 参数设置键盘布局。显示的有效语法为

## Network OPTIONS(网络选项)

https://www.qemu.org/docs/master/system/qemu-manpage.html#hxtool-5

> 注意：
> -net 选项不再推荐使用，详见：<https://www.qemu.org/2018/05/31/nic-parameter/>

**-netdev tap,id=ID\[,fd=H]\[,ifname=NAME]\[,script=FILE]\[,downscript=DFILE]\[,helper=HELPER]**

在宿主机上自动创建一个 tap 类型的网络设备，并使用 ID 作为该 netdev 的标识符，用于与 -device 进行关联

在 VM 启动时使用 script=FILE 指定的脚本(默认为 /etc/qemu-ifup)来配置当前网络设备，且在虚拟机停止时使用 downscript=DFILE 指定的脚本(默认为/etc/qemu-ifdown)来撤销网络设备配置。

> 可以使用 script=no 和 downscript=no 来禁止执行脚本，以便由自己手动配置网络设备

- id=ID # tap 设备关联的 VLAN 号，`默认值：0`。
- ifname=NAME # tap 设备的名称。`默认值：tapN`。N 从 0 开始，第一个启动的 VM 会创建 tap0 第二个创建 tap2，以此类推)
- 注意，ifname 的值会作为 script 和 downscript 这俩脚本中的第一个参数，可以在脚本中使用 $1 引用该 ifname 选项指定的值。

**-nic \[tap|bridge|user|l2tpv3|vde|netmap|vhost-user|socket]\[,...]\[,mac=macaddr]\[,model=mn]**

> 该选项是新版推荐的，用来代替 -netdev 和 -device 两个选项，以使操作更简便。
> 并且 -nic 接受 -netdev 选项的所有选项。

**-vnc \[IP]:PORT** # 指定 VNC 暴露的端口，0 为 5900，1 为 5901，依此类推。

**-nographic** # 让虚拟机在前台运行，虚拟机的输出信息会在 宿主机直接显示，类似于 virsh console 命令

## Character Device OPTIONS(字符设备选项)

https://www.qemu.org/docs/master/system/qemu-manpage.html#hxtool-6

**-chardev BACKEND,id=ID\[,mux=on|off]\[,OPTIONS]**
BACKEND 可以是以下设备之一：`null`, `socket`, `udp`, `msmouse`, `vc`, `ringbuf`, `file`, `pipe`, `console`, `serial`, `pty`, `stdio`, `braille`, `tty`, `parallel`, `parport`, `spicevmc`, `spiceport`。特定的后端将确定适用的选项。

**-chardev socket,id=ID\[,TCP options or unix options]\[,server]\[,nowait]\[,telnet]\[,websocket]\[,reconnect=seconds]\[,tls-creds=id]\[,tls-authz=id]**
创建一个双向流套接字，可以是 TCP 或 UNIX 套接字。如果指定了路径，则会创建一个 Unix 套接字。如果为 unix 套接字指定了 TCP 选项，则行为是不确定的。

- server # 指定 socket 为监听套接字
- nowait # 指定 QEMU 不应该在等待客户端连接到监听 socket 时阻塞。
- unix options: path=path\[,abstract=on|off]\[,tight=on|off] # 使用 unix socket 时的特定选项
- path # 指定 unix socket 的路径，该选项时必须的。

## TPM Device OPTIONS

## Boot Image or Kernel OPTIONS

## Debug/Expert OPTIONS

## Generic Object OPTIONS(通用对象选项)

-object

# Deivce URL Syntax
