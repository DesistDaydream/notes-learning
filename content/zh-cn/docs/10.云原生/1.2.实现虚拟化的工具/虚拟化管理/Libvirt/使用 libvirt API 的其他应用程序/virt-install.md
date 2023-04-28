---
title: "virt-install"
linkTitle: "virt-install"
weight: 20
---

# 概述

> 参考：
>
> - [GitHub 项目，virt-manager/virt-manager 中的 Manual(手册)](https://github.com/virt-manager/virt-manager/blob/main/man/virt-install.rst)
> - [Manual(手册),virt-install(1)](https://man.cx/virt-install)（另一个网站的 Manual）

virt-install 是一个命令行工具，用于使用 Libvirt 管理程序管理库创建新的 KVM、Xen 或 Linux 容器。请参阅本文档末尾的[示例部分](https://github.com/virt-manager/virt-manager/blob/main/man/virt-install.rst#examples)以快速入门。

virt-manager 在图形化界面创建的虚拟机本质上就是调用的 virt-install 命令在系统中执行的。virt-manager 创建的虚拟机生成的 xml 文件可以推导出 virt-install 创建同样虚拟机所需要使用到的参数。

virt-install 命令中很多参数都可以在 [XML 文件](/docs/10.云原生/1.2.实现虚拟化的工具/虚拟化管理/Libvirt/XML%20文件详解/XML%20文件详解.md)中找到对应的配置。比如 `--memory` 的全部可配置参数可以在 https://libvirt.org/formatdomain.html#memory-allocation 这里找到。

# Syntax(语法)

**virt-install --name NAME --memory MB STORAGE INSTALL \[OPTIONS]**

许多参数都有子选项，要查看与该参数相关联的子选项的完整列表使用例子中类似的命令，例如：virt-install --disk=?

> 随着版本的更新，很多选项都会被更好的选项而替代，那些被弃用的选项可能不会在笔记中出现，具体详见官方文档。

有几个参数是在使用 libvirt 工具安装虚拟机时必须指定的：

- --name is required
- --memory amount in MiB is required
- --disk storage must be specified (override with --disk none)
- 安装方式
  - --location URL, --cdrom CD/ISO, --pxe, --import, --boot hd|cdrom|...

> 注意：
> 在创建虚拟机时，我们一般都会指定一下虚拟机的系统类型，以便优化 virtio 等性能相关功能。所有 virt-manager 支持的虚拟机列表可以通过 `virt-install --osinfo list` 命令列出。

## GENERAL OPTIONS(通用选项)

- **-n, --name STRING** # 指定 VM 名称，必须要全局唯一
- **--memory** # 设定 VM 的内存大小，单位为 MiB
- -**-vcpus=VCPU\[,maxvcpus=MAX]\[,sockets=NUM]\[,cores=NUM]\[,threads=NUM]** # 设定 VCPU 的个数，最大数，CPU 插槽数，内核数，线程数
- -**-cpu=CPU** # 设定 CPU 的型号及特性，如 coreduo 等，可以使用 kvm -cpu ? 来获取支持的 CPU 模式

## INSTALLATION OPTIONS(安装选项)

- **-c, --cdrom=STRING** # 设定从光盘介质安装
- **-l, --location OPTIONS** # 指定本地安装介质
- **--import** # 使用已经存在的磁盘镜像构建 VM。比如通过一个已经正常运行 VM 的文件

## Guest OS OPTIONS(虚拟机操作系统选项)

- **--os-variant, --osinfo OSNAME** # 指定要虚拟机的操作系统的信息。常用来优化 virtio 等性能相关功能。OSNAME 可用的值可以用过 `virt-install --osinfo list` 命令列出。
  - 注意，--osinfo 是新版本的名称

## STORAGE OPTIONS(存储选项)

### --disk

-**-disk /Some/Storage/Path\[,OPT1=VAL1]\[,OPT2=VAL2]\[,.....]**

设置 VM 存储介质，比如最常见的就是 /var/lib/libvirt/images/XXX.qcow2 这种，virt-install 会自动生成目标文件所在路径，并记录成存储池。

除了 VM 本身 OS 的存储介质，还可以设置诸如 cdrom 之类的存储介质。

**SUB_OPTIONS：**

- **device=STRING** # 指定设备类型，如 cdrom，disk。`默认值: disk`
- **bus=STRING** # 指定磁盘总线类型，如 ide、scsi、usb、virtio、xen
- **size=** # 指定新磁盘映像的大小，单位为 GB
- **cache=** # 指定缓存类型，如 none、writethrouth、writeback
- **format=** # 指定磁盘映像格式，如 qcow2、raw、vmdk 等

## NETWORKING OPTIONS(网络选项)

设置 VM 要使用的网络以便连接到宿主机上。说白了，就是告诉 qemu-system 要模拟什么样的网卡，以及要连接到宿主机的哪个网络设备上。可以在 [Domain](docs/10.云原生/1.2.实现虚拟化的工具/虚拟化管理/Libvirt/XML%20文件详解/Domain.md#Network%20interfaces) 中找到更多关于网络的 XML 配置。

如果省略 --network，则将在来 VM 中创建单个NIC(网卡)。如果主机中有一个连接了物理接口的桥接设备，则该设备将用于连接。否则，将使用称为 default 的虚拟网络。可以多次指定此选项以设置多个NIC。

### --network

网络选项分两部分，这是网络类型及其名称，以及为这个类型的网络设置运行时行为

**-w,--network TYPE,OPT1=VAL1,OTP2=VAL2,......**

**TYPE**

- **bridge=BridgeName** # 指定网络连接类型为 bridge 桥接模式，并选择用哪个桥
- **network=STRING** # 连接到 `virsh net-list` 命令列出的虚拟网络中。选项的值是虚拟网络的名称

**SUB_OPTIONS：**

- **model=STRING** # 指定 GuestOS 中的设备型号，如 e1000、virtio、rt18193 等
- **mac=STRING** # 指定 mac 地址，`默认值：随机`

## GRAPHICS OPTIONS(图形选项)

这个选项并不是为 VM 设置任何与显示有关的虚拟硬件，而是指**我们如何访问 VM 的图形界面**。可以在[Domain](docs/10.云原生/1.2.实现虚拟化的工具/虚拟化管理/Libvirt/XML%20文件详解/Domain.md#Graphical%20framebuffers) 中找到更多关于连接 VM 图形界面的 XML 配置。

> 注意，如果想要使用图形界面安装系统，则必须要配置图形选项，否则无法连接到为 VM 虚拟显卡。

如果未指定图形相关选项，则 virt-install 将会在 `${DISPLAY}` 变量被设置时选择合适的图形与 VM 的虚拟显卡对接。否则，默认为 `--graphics none`。

### --graphics

**--graphics TYPE,OPT1=ARG1,OPT2=ARG2,...** #

**TYPE：**

- **vnc** # 在来宾中设置虚拟控制台并将其导出为主机中的 VNC 服务器。除非同时提供端口参数，否则 VNC 服务器将运行在 5900 或以上的第一个空闲端口号上。分配的实际 VNC 显示可以使用 vncdisplay 命令到 virsh 获得（或者可以使用 virt-viewer(1) 来处理这个细节以供使用）。
- **spice** # 使用 Spice 协议导出来宾的控制台。 Spice 允许高级功能，如音频和 USB 设备流，以及改进的图形性能。
- **none** # 不会为来宾分配图形控制台。来宾可能需要在来宾的第一个串行端口上配置文本控制台（这可以通过 --extra-args 选项完成）。命令“virsh console NAME”可用于连接串行设备。

**SUB_OPTIONS：**

- **listen=STRING** # 指定 vnc 监听的地址(默认值通常为 127.0.0.1。i.e.仅限本地主机使用)，如果配置 0.0.0.0，则可以被非宿主机的设备通过宿主机的 IP 与 PORT 来进行 vnc 访问
- **port=NUM** # 指定访问该 VM 的 vnc 所使用的端口

### 其他选项

**--autoconsole STRING** # 在使用 virt-install 创建虚拟机时，将要默认启用的交互式控制台。可用的值有 graphical、text、none。

- 这个选项不是必须的，默认行为是自适应的，取决于 VM 的配置方式。

**--noautoconsole** # 与 `--autoconsole none` 一样。

推荐使用 --noautoconsole，这样执行 virt-install 命令创建虚拟机时不会自动打开 virt-viewer，x11 转发还是比较卡的。推荐使用 VNC 连接端口以访问虚拟机的图像界面。

## VIRTUALIZATION OPTIONS(虚拟化选项)

## DEVICE OPTIONS(设备选项)

指定文本控制台、声音设备、串行接口、并行接口、显示接口等

- **--serial TYPE,OPT1=VAL1,OPT2=VAL2,...** # 指定一个串行设备附加到 VM，TYPE 包括 pty(伪终端)等
- -**-console=** # 指定启动的控制台

## MISCELLANEOUS OPTIONS(其他选项)

# 应用示例

使用本地文件当做磁盘镜像，VM 名字为 test，1024M 大小的内存，1 个 CPU，指定操作系统版本为 centos7，指定所使用的存储文件为 /var/lib/images/test.qcow2(该文件会自动创建)，指定网络桥接到 br0 上、模式为 virtio，指定图形模式为 vnc、把 vnc 暴露到宿主机上、监听端口为 5910

```bash
virt-install --import --name test \
--memory 2048 --vcpus 2 \
--os-variant centos7.0 \
--disk /var/lib/libvirt/images/test.qcow2,size=20 \
--network bridge=br0,model=virtio \
--graphics vnc,listen=0.0.0.0,port=5910
```

## 使用 cdrom 安装系统

```bash
virt-install --name centos7 \
--memory 4096 --vcpus 2 \
--os-variant centos7.0 \
--disk /var/lib/libvirt/images/test/centos7.qcow2,size=100,bus=virtio \
--network bridge=br0,model=virtio \
--graphics vnc,listen=0.0.0.0,port=5911 \
--noautoconsole \
--cdrom /root/iso/CentOS-7-x86_64-DVD-2009.iso
```

创建完成后，可以使用 virt-viewer 访问虚拟机，也可以使用 VNC 连接到 5911 以访问虚拟机，然后开始安装系统。

> 这里使用了 --noautoconsole，所以不会自动打开 virt-viewer，x11 转发还是比较卡的，推荐使用 VNC 连接端口以访问虚拟机的图像界面

这是最基本的创建方式，virt-install 会自动创建很多默认的虚拟设备以满足所需。我们只需要指定网络、连接显示的方式、系统版本、cpu、内存即可

## ChatGPT 通过 xml 文件推倒出来的 virt-install 参数

```yaml
virt-install \
--name=centos7-2009 \
--memory=67108864 \
--vcpus=2 \
--cpu host-passthrough \
--os-type=linux \
--os-variant=centos7 \
--boot menu=off \
--disk path=/var/lib/libvirt/images/centos7-2009.qcow2,device=disk,bus=virtio \
--network bridge=virbr0,model=virtio \
--graphics none \
--console pty,target_type=serial
```
