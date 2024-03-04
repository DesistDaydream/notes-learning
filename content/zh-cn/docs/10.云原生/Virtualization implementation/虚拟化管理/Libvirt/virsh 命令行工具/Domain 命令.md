---
title: Domain 命令
---

# 概述

> 参考：
>
> - https://github.com/libvirt/libvirt/blob/master/docs/manpages/virsh.rst#domain-commands

可以为虚拟机重命名、查看虚拟机信息，状态、等等

## 简单的子命令

**autostart** # 指定 Domain 是否在开机后自启动，可以使用 --disable 选项关闭 Domain 的开机自启功能。

# console - 连接到 VM 的终端

console 用于把虚拟机屏幕上的信息投射到宿主机上，可以直接在宿主机的终端上操作虚拟机。
注意：如果无法通过 console 连接到 VM，则需要在 VM 上启动 serial-getty@.service 服务并在开启服务的时候指定一个终端。e.g.**systemctl enable serial-getty@ttyS0.service --now**

# cpu-stats - 显示 Domain 的 CPU 统计信息

默认显示所有 CPU 的统计信息和总数。仅使用 --total 获取总统计信息，仅使用 start 获取从 --start 开始的 CPU 的 per-cpu 统计信息，仅使用 --count CPU 的统计信息。

# create - 从一个 XML 文件里创建一个 domain

通过 XML 直接启动一台 VM，VM 关闭后，virsh list 列表中该 VM 会消失。

# define - 从一个 XML 文件里定义一个 domain(仅定义不启动)

通过 XML 文件创建一台 VM。

## Syntax(语法)

**define \<file> \[--validate]**

EXAMPLE

- virsh define cirros.xml # 通过 cirros.xml 来定义一个 VM 的属性，如果 VM 不存在则创建

# desc - 显示或者设置一个 domain 的 description 或 title

**desc \<domain> \[--live] \[--config] \[--current] \[--title] \[--edit] \[\[--new-desc] \<STRING>]...**

Note：

- 当不指定--title 时，desc 命令默认修改或查看 domain 的 description
- description 和 title 有三种状态，live、config、current，当修改或者显示的时候，默认为 live 状态
- config 的状态指的是，修改或查看 domain 的 xml 文件。

OPTIONS

- [--domain] \<string> domain name, id or uuid
- --live # 指定当前操作为，运行时状态
- --config # 指定当前操作为，持久配置状态
- --current # 指定当前操作为，当前状态
- --title \[STRING] # 修改或显示 title。指定 STRING 时则会将 domain 的 title 修改为 STRING，不指定则显示 domain 的 title
- --edit # 打开一个编辑器来修改 description 或 title
- [--new-desc] \<STRING> message

EXAMPLE

- virsh desc lchTest # 显示 lchTest 这台虚拟机的当前状态的描述信息
- virsh desc lchTest 10.10.100.200 --config # 为 lchTest 这台 VM 设定描述信息为 10.10.100.200，并将信息写入到 xml 文件中
- virsh desc lchTest --title 10.10.100.200 --config # 指定 lchTest 这台虚拟机的 title 为 10.10.100.200，并将该信息写入到 xml 配置中

# destroy - 摧毁一个 domain，类似于直接拔掉电源

# domblklist - 列出 domain 的所有 blocks(块设备)

列出 Domain 的所有 blocks(块设备)。i.e.指定 domain 所使用的磁盘文件

## Syntax(语法)

**domblklist \<DOMAIN> [OPTIONS]**
OPTIONS

- --inactive #
- --details # 列出的信息还包括类型和设备

# domifaddr - 从正在运行的 domain 中获取网络接口的 IP 地址

**domifaddr \<domain> \[--interface \<string>] \[--full] \[--source \<string>]**

该信息包括：Name MAC address Protocol Address

# domiflist - 列出 domain 所有的虚拟接口

**domiflist \<domain> \[--inactive]**

该信息包括：Interface Type Source Model MAC

EXAMPLE

- virsh domiflist testvm # 列出名为 testvm 这台虚拟机的虚拟接口信息

# domifstat - 获取指定的 domain 的网络接口状态信息

domifstat \<domain> \<interface>

EXAMPLE

# dominfo - 返回指定 domain 的基本信息，包括该 domain 的 name、uuid、mem、cpu 等

# dommemstat - 获取指定 domain 的内存状态信息

https://libvirt.org/manpages/virsh.html#dommemstat

dommemstat 会获取正在运行的 domain 的内存统计信息。可以获取到如下信息：

- swap_in - The amount of data read from swap space (in KiB)
- swap_out - The amount of memory written out to swap space (in KiB)
- major_fault - The number of page faults where disk IO was required
- minor_fault - The number of other page faults
- unused - The amount of memory left unused by the system (in KiB)
- available - The amount of usable memory as seen by the domain (in KiB)
- actual # 当前 Current balloon value (in KiB)
- rss - Resident Set Size of the running domain's process (in KiB)
- usable - The amount of memory which can be reclaimed by balloon without causing host swapping (in KiB)
- last-update - Timestamp of the last update of statistics (in seconds)
- disk_caches - The amount of memory that can be reclaimed without additional I/O, typically disk caches (in KiB)
- hugetlb_pgalloc - The number of successful huge page allocations initiated from within the domain
- hugetlb_pgfail - The number of failed huge page allocations initiated from within the domain

上述这些字段的可用性取决于管理程序。输出中缺少不支持的字段。如果与更新版本的 libvirtd 通信，可能会出现其他字段。

## Syntax(语法)

**dommemstat domain \[--period seconds] \[\[--config] \[--live] | \[--current]]**

# domstate

返回有关 Domain 的状态。 --reason 选项告诉 virsh 还需要打印状态的原因。

# domif-setlink - 设定 domain 网卡的状态

可以控制 VM 网卡的开关，关闭 VM 的 link，则开机不会自动启动网卡

EXAMPLE

- virsh domif-setlink lichenhao--interface 52:54:00:6a:86:89 --state down # 关闭 lichenhao 这台虚拟机的指定网卡

# domrename - 重命名一个 Domain

重命名一个**未激活状态**的 Domain。`virsh domrename <DOMAIN> <NewName>`

# domxml-to-native - 根据 domain 的 XML 描述文件，转换成 qemu-kvm 创建虚拟机的命令

**domxml-to-native \<FORMAT> \[OPTIONS]**

FORMAT

- qemu-argv # QEMU/KVM 类型虚拟化，必须使用此格式

OPTIONS

- --domain \<STRING> domain name, id or uuid
- --xml \<STRING> xml data file to export from

EXAMPLE

- **virsh domxml-to-native qemu-argv --domain lichenhao.bj-net**# 根据 lichenhao.bj-net 虚拟机的 xml 文件，生成 qemu-kvm 命令行
- lichenhao.bj-net.xml 该文件会转换成如下内容

```bash
LC_ALL=C PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin
QEMU_AUDIO_DRV=none
/usr/libexec/qemu-kvm -name lichenhao.bj-net \
-machine pc-i440fx-rhel7.0.0,accel=kvm,usb=off,dump-guest-core=off \
-cpu Skylake-Server-IBRS,-ds,-acpi,+ss,-ht,-tm,-pbe,-dtes64,-monitor,-ds_cpl,-vmx,-smx,-est,-tm2,-xtpr,-pdcm,-dca,-osxsave,-tsc_adjust,+clflushopt,-intel-pt,+pku,-ospke,+avx512vnni,+md-clear,+stibp,+ssbd,+hypervisor,-arat \
-m 4096 \
-realtime mlock=off
\
-smp 2,sockets=2,cores=1,threads=1 \
-uuid 51b47472-5564-424e-90b5-19e9eecfd671 \
-no-user-config -nodefaults \

-chardev socket,id=charmonitor,path=/var/lib/libvirt/qemu/domain--1-lichenhao.bj-net/monitor.sock,server,nowait \
-mon chardev=charmonitor,id=monitor,mode=control \

-rtc base=utc,driftfix=slew \
-global kvm-pit.lost_tick_policy=delay \
-no-hpet -no-shutdown \
-global PIIX4_PM.disable_s3=1 \
-global PIIX4_PM.disable_s4=1 \
-boot strict=on \
-device ich9-usb-ehci1,id=usb,bus=pci.0,addr=0x4.0x7 \
-device ich9-usb-uhci1,masterbus=usb.0,firstport=0,bus=pci.0,multifunction=on,addr=0x4 \
-device ich9-usb-uhci2,masterbus=usb.0,firstport=2,bus=pci.0,addr=0x4.0x1 \
-device ich9-usb-uhci3,masterbus=usb.0,firstport=4,bus=pci.0,addr=0x4.0x2 \
-device virtio-serial-pci,id=virtio-serial0,bus=pci.0,addr=0x5 \

-drive file=/var/lib/libvirt/images/lichenhao.bj-net.qcow2,format=qcow2,if=none,id=drive-virtio-disk0 \
-device virtio-blk-pci,scsi=off,bus=pci.0,addr=0x6,drive=drive-virtio-disk0,id=virtio-disk0,bootindex=1 \

-drive if=none,id=drive-ide0-0-0,readonly=on \
-device ide-cd,bus=ide.0,unit=0,drive=drive-ide0-0-0,id=ide0-0-0 \

-netdev tap,fd=28,id=hostnet0 \
-device virtio-net-pci,netdev=hostnet0,id=net0,mac=52:54:00:6d:fa:f0,bus=pci.0,addr=0x3 \

-chardev pty,id=charserial0 \
-device isa-serial,chardev=charserial0,id=serial0 \

-chardev socket,id=charchannel0,path=/var/lib/libvirt/qemu/channel/target/domain--1-lichenhao.bj-net/org.qemu.guest_agent.0,server,nowait \
-device virtserialport,bus=virtio-serial0.0,nr=1,chardev=charchannel0,id=channel0,name=org.qemu.guest_agent.0 \

-chardev spicevmc,id=charchannel1,name=vdagent \
-device virtserialport,bus=virtio-serial0.0,nr=2,chardev=charchannel1,id=channel1,name=com.redhat.spice.0 \

-device usb-tablet,id=input0,bus=usb.0,port=1 \
-vnc 127.0.0.1:0 \
-vga qxl \
-global qxl-vga.ram_size=67108864 \
-global qxl-vga.vram_size=67108864 \
-global qxl-vga.vgamem_mb=16 \
-global qxl-vga.max_outputs=1 \

-chardev spicevmc,id=charredir0,name=usbredir \
-device usb-redir,chardev=charredir0,id=redir0,bus=usb.0,port=2 \

-chardev spicevmc,id=charredir1,name=usbredir \
-device usb-redir,chardev=charredir1,id=redir1,bus=usb.0,port=3 \

-device virtio-balloon-pci,id=balloon0,bus=pci.0,addr=0x7 \

-object rng-random,id=objrng0,filename=/dev/urandom \
-device virtio-rng-pci,rng=objrng0,id=rng0,bus=pci.0,addr=0x8 \

-msg timestamp=on
```

# dumpxml - 显示 domain 的 XML 格式的信息

EXAMLPE

- virsh dumpxml lchTest # 显示 lchTest 这个虚拟机的 xml 信息

# edit - 编辑一个 domain 的 XML 配置

# qemu-agent-command - 向 domain 中执行 QEMU Guest Agent 命令

可用的 QGA 命令详见：[QMP 命令参考](/docs/10.云原生/Virtualization%20implementation/KVM_QEMU/QEMU%20Guest%20Agent/QMP%20命令参考.md)

## Syntax(语法)

**qemu-agent-command DOMAIN \[OPTIONS] CMD**

命令示例详见：应用示例

# reboot - 重新启动一个 domainre

# shutdown - 优雅得关闭 domain

# start - 启动一个 domain

# undefine - 取消定义一个 domain。i.e.删除一台虚拟机

该命令会删除 /etc/libvirt/qemu/ 目录下描述该 domain 的 xml 文件

## Syntax(语法)

**undefine \<DOMAIN> \[OPTIONS]**

OPTIONS

- --managed-save remove domain managed state file
- --storage \<STRING> remove associated storage volumes (comma separated list of targets or source paths) (see domblklist)
- **--remove-all-storage** # 移除所有与该 domain 关联的存储卷。(谨慎使用)
- **--delete-snapshots** # 删除与卷关联的快照，需要--remove-all-storage（必须由存储驱动程序支持）
- --wipe-storage wipe data on the removed volumes
- **--snapshots-metadata** # 删除 domain 所有快照元数据（如果不活动）
- --nvram remove nvram file, if inactive
- --keep-nvram keep nvram file, if inactive

# vncdisplay - 输出 Domain 的 VNC 显示的 IP 和端口
