---
title: "QEMU 设备模拟逻辑"
linkTitle: "QEMU 设备模拟逻辑"
weight: 2
---

# 概述

> 参考：
> 
> - [官方文档，系统模拟-设备模拟](https://www.qemu.org/docs/master/system/device-emulation.html)

QEMU 模拟设备主要是半虚拟化设备，从[这里](/docs/IT学习笔记/10.云原生/1.2.实现虚拟化的工具/KVM_QEMU/KVM_QEMU.md#Paravirtualized%20Devices(半虚拟化设备))可以看到简单的介绍。


# QEMU Storage Emulation(QEMU 存储模拟)

与 [网络模拟](https://www.yuque.com/desistdaydream/learning/tkr8dt#03psa) 类似，QEMU 想要让虚拟机获得一块硬盘，也需要由两部分组成一个完整的存储功能。

1. **front-end(前端)** # VM 中的 块设备
2. **back-end(后端)** # 宿主机中的与 VM 中模拟出来的块设备进行交互的设备。

# QEMU Network Emulation(QEMU 网络模拟)

> 参考：
> 
> - <https://wiki.qemu.org/Documentation/Networking>
> - https://www.qemu.org/docs/master/system/net.html
> - <https://www.qemu.org/2018/05/31/nic-parameter/>，老版原理，将弃用

QEMU 想要让虚拟机与外界互通，需要由两部分组成一个完整的网络功能：

- **front-end(前端)** # VM 中的 NIC(Network Interface Controller，即人们常说的`网卡`)。
  - VM 中的 NIC 是由 QEMU 模拟出来的，在支持 PCI 卡的系统上，通常可以是 e1000 网卡、rtl8139 网卡、virtio-net 设备。
- **back-end(后端)** # 宿主机中的与 VM 中模拟出来的 NIC 进行交互的设备。
  - back-end 有多种类型可以使用，这些后端可以用于将 VM 连接到真实网络，或连接到另一个 VM
    - [TAP ](https://www.qemu.org/docs/master/system/net.html#using-tap-network-interfaces)# 将 VM 连接到真实网络的标准方法
    - [User mode network stack](https://www.qemu.org/docs/master/system/net.html#using-the-user-mode-network-stack)

效果如图所示：
![](https://notes-learning.oss-cn-beijing.aliyuncs.com/zuowkm/1616124035097-0a64383e-f37f-4cc3-bdc2-3c7502189b7d.png)

## 基本应用示例

在使用 qemu-kvm 命令创建虚拟机时，通过一组两个选项来为虚拟机创建一个网络设备。比如：

- `-netdev tap,id=n1` # 在宿主机创建一个`后端设备`，这是一个 tap 类型的网络设备(tap 类型的设备路径为 /dev/net/tun)
- `-device virtio-net-pci,netdev=n1` # 在 VM 中模拟一个`前端设备`，这是一个 virtio-net-pci 类型的网卡

完整的命令如下：

```bash
qemu-kvm -m 4096 -smp 2 -name test \
-drive file=/var/lib/libvirt/images/test-1.bj-net.qcow2,format=qcow2,if=virtio \
-netdev tap,id=n1 \
-device virtio-net-pci,netdev=n1 \
-vnc :3
```

此时 qemu-kvm 发现 `后端设备的 id` 与 `前端设备的属性(netdev)的值` 一致，那么 qemu-kvm 就会将 两端设备关联。因此，在 VM 启动时，其打开了设备文件 /dev/net/tun 并获得了读写该文件的文件描述符 (FD)XX，同时向内核注册了一个 tap 类型虚拟网卡 tapX，tapX 与 FD XX 关联，虚拟机关闭时 tapX 设备会被内核释放。此虚拟网卡 tapX 一端连接用户空间程序 qemu-kvm，另一端连接主机链路层。

```bash
## 通过进程，获取该进程所使用的网络设备
# 当先宿主机上有3个虚拟机，分别对应 82649、82747、144776 这三个进程
# 82649 与 82747 使用网卡多队列功能，启用了8个队列，144776 仅有一个网卡队列
# 所以，82649 与 82747 打开了8个 /dev/net/tun 设备，而 144776 打开了1个 /dev/net/tun 设备
[root@host-3 fdinfo]# lsof /dev/net/tun
COMMAND     PID USER   FD   TYPE DEVICE SIZE/OFF   NODE NAME
qemu-kvm  82649 qemu   27u   CHR 10,200      0t0 102414 /dev/net/tun
qemu-kvm  82649 qemu   29u   CHR 10,200      0t0 102414 /dev/net/tun
qemu-kvm  82649 qemu   31u   CHR 10,200      0t0 102414 /dev/net/tun
qemu-kvm  82649 qemu   32u   CHR 10,200      0t0 102414 /dev/net/tun
qemu-kvm  82649 qemu   33u   CHR 10,200      0t0 102414 /dev/net/tun
qemu-kvm  82649 qemu   34u   CHR 10,200      0t0 102414 /dev/net/tun
qemu-kvm  82649 qemu   35u   CHR 10,200      0t0 102414 /dev/net/tun
qemu-kvm  82649 qemu   36u   CHR 10,200      0t0 102414 /dev/net/tun
qemu-kvm  82747 qemu   28u   CHR 10,200      0t0 102414 /dev/net/tun
qemu-kvm  82747 qemu   31u   CHR 10,200      0t0 102414 /dev/net/tun
qemu-kvm  82747 qemu   32u   CHR 10,200      0t0 102414 /dev/net/tun
qemu-kvm  82747 qemu   33u   CHR 10,200      0t0 102414 /dev/net/tun
qemu-kvm  82747 qemu   34u   CHR 10,200      0t0 102414 /dev/net/tun
qemu-kvm  82747 qemu   35u   CHR 10,200      0t0 102414 /dev/net/tun
qemu-kvm  82747 qemu   36u   CHR 10,200      0t0 102414 /dev/net/tun
qemu-kvm  82747 qemu   37u   CHR 10,200      0t0 102414 /dev/net/tun
qemu-kvm 144776 qemu   31u   CHR 10,200      0t0 102414 /dev/net/tun
# 查看 144776 进程的 fdinfo 中的 31 号描述符文件，可以看到该进程关联的网络设备是 vnet2
[root@host-3 fdinfo]# cat /proc/144776/fdinfo/31
pos:	0
flags:	0104002
mnt_id:	20
iff:	vnet2

## 通过网络设备，获取使用该设备的进程
# 可以通过一条命令来直接获取使用指定 pid 设备的进程
# 下面的命令可以获取使用 vnet2 这个 tap 类型的网络设备的进程。
[root@host-3 fdinfo]# egrep -l iff:.*vnet2 /proc/*/fdinfo/* 2> /dev/null | cut -d/ -f3
144776
```

所以 144776 这个进程下的虚拟机经过其内部网卡发送的数据包，都会发送到 /dev/net/tun 设备上，然后根据其文件描述符，找到对应的 tap 设备，将数据包转发过去。这样，虚拟机的数据就通过网络，发送到物理机上了。

> 获取 tap 设备 与 VM 关联性的方法参考：<https://unix.stackexchange.com/questions/462171/how-to-find-the-connection-between-tap-interface-and-its-file-descriptor>

## virbr0 说明

virbr0 是 KVM 默认创建的一个 Bridge，其作用是为连接其上的虚机网卡提供 NAT 访问外网的功能。virbr0 默认分配了一个 IP 192.168.122.1，并为连接其上的其他虚拟网卡提供 DHCP 服务。

- 需要说明的是，使用 NAT 的虚机 VM1 可以访问外网，但外网无法直接访问 VM1。 因为 VM1 发出的网络包源地址并不是 192.168.122.6，而是被 NAT 替换为宿主机的 IP 了。
- 这个与使用 br0 不一样，在 br0 的情况下，VM1 通过自己的 IP 直接与外网通信，不会经过 NAT 地址转换。
