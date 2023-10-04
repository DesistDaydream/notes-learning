---
title: QEMU Guest Agent
---

# 概述

> 参考：
>
> - [官方文档，系统模拟管理和互操作性-QEMU 客户机代理](https://www.qemu.org/docs/master/interop/qemu-ga.html)
> - [官方文档，系统模拟管理和互操作性-QEMU 客户机代理协议参考](https://www.qemu.org/docs/master/interop/qemu-ga-ref.html)(i.e.QGA 的 QMP API 参考文档)
> - <https://wiki.qemu.org/Features/GuestAgent>
> - <https://www.toutiao.com/i6646012291059810823/>
> - <https://www.shuzhiduo.com/A/QV5ZgK76dy/>

**QEMU Guest Agent(QEMU 虚拟机代理，简称 QGA)** 是一个类似于 VMware Tools 的工具，用来辅助 Hypervisor 实现对 VM 的管理。

QEMU Guest Agent 旨在通过标准的 **QEMU Monitor Protocol(QEMU 监控协议，简称 QMP)**命令，实现 VM 与 宿主机 之间数据交互的功能。(比如可以在不登陆 VM 的情况下，让 VM 执行某些命令或者直接获取 VM 的信息)

## QEMU Guest Agent 架构

QGA 功能的实现与 虚拟化 I/O 的实现，是相同的原理。KVM/QEMU 会在 VM 中模拟一个 I/O 设备，并通过 ID 关联到宿主机的某个文件或设备上，这样就可以实现宿主机与虚拟机之间的交互。其实说白了，这个年代基本都是半虚拟化设备的实现方式，通过两部分来实现完整的功能：

1. Host Device # 宿主机设备
   1. **socket(套接字)**是宿主机中实现 QGA 的设备**。**也可以是其他未来待发明的东西。
2. Guest Driver # 虚拟机驱动
   1. **virtio-serial(半虚拟化的串口设备)**是 VM 中实现 QGA 的设备**。**也可以是 isa-serial 等模拟设备。

除了基本的半虚拟化设备，VM 中还需要一个程序来处理宿主机发来的 QMP 命令：

- 一个名为 **qemu-ga** 的二进制文件。

默认情况下，qemu-ga 会监听 VM 中的 virtio-serial(默认为 /dev/virtio-ports/org.qemu.guest_agent.0) 串口设备。这样一来，所有从宿主机向 socket 发送的命令，都会传递到 virtio-serial 中，进而被监听该设备的 qemu-ga 接收，并处理该命令。

> /dev/virtio-ports/org.qemu.guest_agent.0 实际上是 /dev/vport1p1 的符号链接，这类设备名字都是 vport1pX 这种格式，X 从 0 开始

为了可以让 VM 实时处理 virtio-serial 中的数据，所以，qemu-ga 以守护程序的方式运行在宿主机上，这个守护程序默认为 qemu-guest-agent.service。

**这样一来，qemu-ga、qemu-guest-agent.service、virtio-serial、socket 这四个东西，就组成了一个完整的 QGA 技术栈。**
![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/mxgyxv/1624240513497-c71382f6-ebc0-4486-a2de-68de15875a07.png)

### **virtio-serial**和 **socket 之间的数据通信路线**

假如我现在启动如下虚拟机，在 宿主机上创建了一个 socket(/tmp/qga.sock)，并为 VM 模拟了一个 virtio-serial

```bash
~]# qemu-kvm -m 4096 -smp 2 -name test \
-drive file=/var/lib/libvirt/images/test-2.bj-net.qcow2,format=qcow2,if=virtio \
-netdev tap,id=n1 \
-device virtio-net-pci,netdev=n1 \
-vnc :3 \
-chardev socket,path=/tmp/qga.sock,server,nowait,id=qga0 \
-device virtio-serial \
-device virtserialport,chardev=qga0,name=org.qemu.guest_agent.0
```

查看 VM 的进程 和 socket

```bash
~]# fuser qga.sock
/tmp/qga.sock:       267751
~]# ps -ef f | grep 267751
root     267766      2  0 01:45 ?        S      0:00  \_ [kvm-pit/267751]
root     267751 238001 99 01:45 pts/3    Sl+    0:36  |       \_ qemu-kvm -m 4096 -smp 2 -name test -drive file=/var/lib/libvirt/images/test-2.bj-net.qcow2,format=qcow2,if=virtio -netdev tap,id=n1 -device virtio-net-pci,netdev=n1 -vnc :3 -chardev socket,path=/tmp/qga.sock,server,nowait,id=qga0 -device virtio-serial -device virtserialport,chardev=qga0,name=org.qemu.guest_agent.0
```

可以看到，qga.sock 被 267751 进程使用着，而这个进程就是我们启动的一台虚拟机。当连接该 socket 后，读写的数据都会经过 267751 进程，并根据其中的 chardev 与 virtserialport 的关系，将数据送到 qemu 模拟的 virtio-serial 设备上，进而被 VM 内的 qemu-ga 接受并处理。

## **QEMU Monitor Protocol**

**QEMU Monitor Protocol(QEMU 监控协议，简称 QMP)**
可用的 QMP 命令详见：[QMP 命令参考](/docs/10.云原生/1.2.实现虚拟化的工具/KVM_QEMU/QEMU%20Guest%20Agent/QMP%20命令参考.md)

# QEMU Guest Agent 部署

有多种方式可以部署 QEMU Guest Agent

## KVM/QEMU 创建 QGA

首先启动一个 VM

```bash
~]# qemu-kvm -m 4096 -smp 2 -name test \
-drive file=/var/lib/libvirt/images/test-2.bj-net.qcow2,format=qcow2,if=virtio \
-netdev tap,id=n1 \
-device virtio-net-pci,netdev=n1 \
-vnc :3 \
-chardev socket,path=/tmp/qga.sock,server,nowait,id=qga0 \
-device virtio-serial \
-device virtserialport,chardev=qga0,name=org.qemu.guest_agent.0
```

在 VM 中安装 QGA，并启动(直接 yum 即可，**一般通过 libvirt 启动的虚拟机，都默认自带 qemu-guest-aent**)。一般默认配置即可，如果 qemu-ga 未监听默认设备，修改配置文件，并重启服务即可。

```bash
yum install qemu-guest-agent -y
systemctl start qemu-guest-agent.service && systemctl enable qemu-guest-agent.service
```

在宿主机上然后使用 socat 连接 /tmp/qga.sock 即可开始使用 QGA。

```bash
~]# socat - unix:/tmp/qga.sock
{"execute":"guest-get-host-name"}
# 这是发送给 QGA 的数据
{"return": {"host-name": "centos8-2004"}} # 这是 QGA 返回的数据
```

## libvirt 创建 QGA

libvrit 提供了专门的 DomainQemuAgentCommand API（对应 virsh qemu-agent-command 命令）来和 QGA 通讯，另外还有有些 libvirt 内置 api 也可以支持 QGA，例如 virsh 的 reboot、shutdown 等命令。

**通过 libvirt 启动 KVM/QEMU 的虚拟机不需要做任何配置，默认就会自动创建一个 channel**，VM 的 xml 中 channel 配置段如下：

```html
<channel type="unix">
  <target type="virtio" name="org.qemu.guest_agent.0" />
  <address type="virtio-serial" controller="0" bus="0" port="1" />
</channel>
```

默认情况，虚拟机启动后会在 /var/lib/libvirt/qemu/channel/target/DOMAIN/ 目录下生成一个名为 org.qemu.guest_agent.0 的 socket 文件，如果在 virt-manager 中查看该 channel 设备，也可以看到 Source path
![](https://notes-learning.oss-cn-beijing.aliyuncs.com/mxgyxv/1616123963420-5de205c9-b9f1-4910-a984-12f711d617b1.png)
注意：该 socket 被 libvirtd 始终连接着，无法通过 socat 等工具再次连接使用，通过 fuser 命令可以看到占用该 socket 的进程(255192 就是 qemu-kvm 运行虚拟机的进程)：

```bash
~]# fuser org.qemu.guest_agent.0
/var/lib/libvirt/qemu/channel/target/domain-54-lichenhao.bj-net/org.qemu.guest_agent.0: 255192
```

而在 VM 内部，默认情况下，会在 /dev 目录下生成也会自动生成串口设备，并自动生成 qemu-ga 文件，且启动 qemu-guest-agent.service 服务

```bash
~]# ll /dev/virtio-ports/
total 0
lrwxrwxrwx 1 root root 11 Nov 21 00:41 com.redhat.spice.0 -> ../vport1p2
lrwxrwxrwx 1 root root 11 Nov 21 00:41 org.qemu.guest_agent.0 -> ../vport1p1
lrwxrwxrwx 1 root root 11 Nov 21 00:41 org.qemu.guest_agent.1 -> ../vport1p3
[root@lichenhao ~]# systemctl status qemu-guest-agent.service
● qemu-guest-agent.service - QEMU Guest Agent
Loaded: loaded (/usr/lib/systemd/system/qemu-guest-agent.service; disabled; vendor preset: enabled)
Active: active (running) since Fri 2020-11-20 23:40:33 CST; 44min ago
Main PID: 826 (qemu-ga)
Tasks: 1 (limit: 23968)
Memory: 2.7M
CGroup: /system.slice/qemu-guest-agent.service
└─826 /usr/bin/qemu-ga --method=virtio-serial --path=/dev/virtio-ports/org.qemu.guest_agent.0 --blacklist= -F/etc/qemu-ga/fsfreeze-hook
Nov 20 23:40:33 lichenhao.bj-net systemd[1]: Started QEMU Guest Agent.
```

此时，宿主机的 socket 与 VM 中的串口设备(/dev/virtio-ports/org.qemu.guest_agent.0) 之间建立了一条 channel

然后宿主机通过 libvirt 的 API(即 virsh qemu-agent-command 命令)，即可向 VM 中发送指令，效果如下：

```bash
~]# virsh qemu-agent-command lichenhao.bj-net --pretty '{"execute":"guest-get-osinfo"}'
{
  "return": {
    "name": "CentOS Linux",
    "kernel-release": "4.18.0-193.28.1.el8_2.x86_64",
    "version": "8 (Core)",
    "pretty-name": "CentOS Linux 8 (Core)",
    "version-id": "8",
    "kernel-version": "#1 SMP Thu Oct 22 00:20:22 UTC 2020",
    "machine": "x86_64",
    "id": "centos"
  }
}
```

### 通过 socat 等工具连接 socket

由于通过 libvirt 创建的虚拟机的这个 channel 的 socket 一直被 libvirt 占用，所以无法使用别的方式连接。这时候我们可以自己再创建一个 channel

那么在宿主机上的 libvirt 将不会建立与 socket 建立连接。

宿主机上的 libvirt 的 xml：

```html
<channel type="unix">
  <target type="virtio" name="org.qemu.guest_agent.1" />
  <address type="virtio-serial" controller="0" bus="0" port="3" />
</channel>
```

> 我们也可以在其中加入 <source mode='bind' path='/tmp/org.qemu.guest_agent.1'/> 这一行，来手动指定 socket 文件的绝对路径。否则 socket 默认在 /var/lib/libvirt/qemu/channel/target/DOMAIN/ 目录下

此时 VM 内部的 qemu-guest-agent 进程还是在连接 org.qemu.guest_agent.0 设备，为了使用 .1 ，我们需要将 /usr/lib/systemd/system/qemu-guest-agent.service 文件中的所有 .0 改为 .1，然后重启服务。

```bash
~]# sed -i 's/qemu.guest_agent.0/qemu.guest_agent.1/g' /usr/lib/systemd/system/qemu-guest-agent.service
~]# systemctl daemon-reload
~]# systemctl restart qemu-guest-agent.service
```

最后，我们就可以在宿主机上使用 socat 去连接 socket 文件了：

```bash
~]# socat unix:/var/lib/libvirt/qemu/channel/target/domain-55-lichenhao.bj-net/org.qemu.guest_agent.1 readline
{"execute":"guest-get-osinfo"}
{"return": {"name": "CentOS Linux", "kernel-release": "4.18.0-193.28.1.el8_2.x86_64", "version": "8 (Core)", "pretty-name": "CentOS Linux 8 (Core)", "version-id": "8", "kernel-version": "#1 SMP Thu Oct 22 00:20:22 UTC 2020", "machine": "x86_64", "id": "centos"}}
```

# QEMU Guest Agent 关联文件

**/etc/sysconfig/qemu-ga** # qemu-ga 的配置文件

```bash
# rpc 黑名单列表。这里用来定义 qemu-ga 不处理来自宿主机的哪些 QMP 命令。
BLACKLIST_RPC=guest-file-open,guest-file-close,guest-file-read,guest-file-write,guest-file-seek,guest-file-flush,guest-exec,guest-exec-status

# Fsfreeze hook script specification.
#
# FSFREEZE_HOOK_PATHNAME=/dev/null           : disables the feature.
#
# FSFREEZE_HOOK_PATHNAME=/path/to/executable : enables the feature with the
# specified binary or shell script.
#
# FSFREEZE_HOOK_PATHNAME=                    : enables the feature with the
# default value (invoke "qemu-ga --help" to interrogate).
FSFREEZE_HOOK_PATHNAME=/etc/qemu-ga/fsfreeze-hook
```

**/usr/lib/systemd/system/qemu-guest-agent.service** # qemu-ga 守护程序文件

```bash
~]# systemctl cat qemu-guest-agent.service
# /usr/lib/systemd/system/qemu-guest-agent.service
[Unit]
Description=QEMU Guest Agent
BindsTo=dev-virtio\x2dports-org.qemu.guest_agent.0.device
After=dev-virtio\x2dports-org.qemu.guest_agent.0.device
IgnoreOnIsolate=True

[Service]
UMask=0077
EnvironmentFile=/etc/sysconfig/qemu-ga
ExecStart=/usr/bin/qemu-ga \
  --method=virtio-serial \
  --path=/dev/virtio-ports/org.qemu.guest_agent.0 \
  --blacklist=${BLACKLIST_RPC} \
  -F${FSFREEZE_HOOK_PATHNAME}
StandardError=syslog
Restart=always
RestartSec=0

[Install]
WantedBy=dev-virtio\x2dports-org.qemu.guest_agent.0.device
```

# ovirt-guest-agent

ovirt-guest-agent 是和 qemu-guest-agent 并列的一个概念。在使用 oVirt 作为虚拟化管理时，虚拟机内部安装下面三个工具，和 ovirt 配合能够提高虚拟机的用户体验和性能。

- oVirt Guest Agent：原理与 qemu-guest-agent 类似，但是提供的功能有所区别。
- Spice Agent：提高 spice 连接虚拟机的用户体验。
- VirtIO Drivers：包含一些驱程序，VirtIO Serial、VirtIO SCS、VirtIO Network、Memory Ballooning

qemu：
<https://wiki.qemu.org/Features/GuestAgent>
