---
title: "virt-install"
linkTitle: "virt-install"
weight: 20
---

# 概述

> 参考：
> - [GitHub 项目，virt-manager/virt-manager 中的 Manual(手册)](https://github.com/virt-manager/virt-manager/blob/main/man/virt-install.rst)
> - [Manual(手册),virt-install(1)](https://man.cx/virt-install)（另一个网站的 Manual）

virt-install 是一个命令行工具，用于使用 Libvirt 管理程序管理库创建新的 KVM、Xen 或 Linux 容器。请参阅本文档末尾的[示例部分](https://github.com/virt-manager/virt-manager/blob/main/man/virt-install.rst#examples)以快速入门。

virt-manager 在图形化界面创建的虚拟机本质上就是调用的 virt-install 命令在系统中执行的。virt-manager 创建的虚拟机生成的 xml 文件可以推导出 virt-install 创建同样虚拟机所需要使用到的参数。

# Syntax(语法)

**virt-install --name NAME --memory MB STORAGE INSTALL \[OPTIONS]**

许多参数都有子选项，要查看与该参数相关联的子选项的完整列表使用例子中类似的命令，例如：virt-install --disk=?

有几个参数是在使用 libvirt 工具安装虚拟机时必须指定的：

- --name is required
- --memory amount in MiB is required
- --disk storage must be specified (override with --disk none)

安装方式(--location URL, --cdrom CD/ISO, --pxe, --import, --boot hd|cdrom|...)

> 注意：
> 在创建虚拟机时，我们一般都会指定一下虚拟机的系统类型，以便优化 virtio 等性能相关功能。所有 virt-manager 支持的虚拟机列表可以通过 `virt-install --osinfo list` 命令列出。

## GENERAL OPTIONS(通用选项)

- **-n NAME|--name NAME** # 指定 VM 名称，必须要全局唯一
- **-r MEMORY|--ram MEMORY** # 设定 VM 的内存大小，单位为 MiB
- -**-vcpus=VCPU\[,maxvcpus=MAX]\[,sockets=NUM]\[,cores=NUM]\[,threads=NUM]** # 设定 VCPU 的个数，最大数，CPU 插槽数，内核数，线程数
- -**-cpu=CPU** # 设定 CPU 的型号及特性，如 coreduo 等，可以使用 kvm -cpu ? 来获取支持的 CPU 模式

## INSTALLATION OPTIONS(安装选项)

- **-c, --cdrom=PATH **# 设定从光盘介质安装
- **-l LOCATION|--location OPTIONS** # 指定本地安装介质
- **--import** # 使用已经存在的磁盘镜像构建 VM。比如通过一个已经正常运行 VM 的文件

## Guest OS OPTIONS(虚拟机操作系统选项)

- **--os-variant, --osinfo OSNAME** # 指定要虚拟机的操作系统的信息。常用来优化 virtio 等性能相关功能。OSNAME 可用的值可以用过 `virt-install --osinfo list` 命令列出。
  - 注意，--osinfo 是新版本的名称

## STORAGE OPTIONS(存储选项)

- -**-disk /Some/Storage/Path\[,OPT1=VAL1]\[,OPT2=VAL2]\[,.....] **# 设定用作 GuestOS 存储的介质,后面是使用的路径\[选项=值(OPTIONs=VALUEs)]，可以在路径中使用 img 映像直接启动，以下是常用的 OPTIONS
  - device= # 指定设备类型，如 cdrom，disk，默认为 disk
  - bus= # 指定磁盘总线类型，如 ide、scsi、usb、virtio、xen
  - size= # 指定新磁盘映像的大小，单位为 GB
  - cache= # 指定缓存类型，如 none、writethrouth、writeback
  - format= # 指定磁盘映像格式，如 qcow2、raw、vmdk 等

## NETWORKING OPTIONS(网络选项)

- -**w NETWORK|--network 类型=名称,OPT1=VAL1,OTP2=VAL2,......** # 有 4 中网络类型可供选择
  - bridge=BridgeName # 指定网络连接类型为 bridge 桥接模式，并选择用哪个桥
  - network=NMAE #
    - model= # 指定 GuestOS 中的设备型号，如 e1000、virtio、rt18193 等
    - mac= # 指定 mac 地址，默认使用随机 mac 地址

## GRAPHICS OPTIONS(图形选项)

- **--graphics TYPE,OPT1=ARG1,OPT2=ARG2,...** # 指定 display VNC 的类型，可用的 TYPE 有 vnc、spice、none。其中 OPT=ARG 还可以为 graphics 指定更多的选项和值。下面是可用的 OPT：
  - listen=IP # 指定 vnc 监听的地址(默认值通常为 127.0.0.1。i.e.仅限本地主机使用)，如果配置 0.0.0.0，则可以被非宿主机的设备通过宿主机的 IP 与 PORT 来进行 vnc 访问
  - port=NUM # 指定访问该 VM 的 vnc 所使用的端口

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
virt-install --cdrom /root/iso/CentOS-7-x86_64-DVD-1908.iso --name centos7 \
--memory 4096 --vcpus 4 \
--os-variant centos7.0 \
--disk /var/lib/libvirt/images/centos7.qcow2,size=100,bus=virtio \
--network bridge=br0,model=virtio \
--graphics vnc,listen=0.0.0.0,port=5911
```

创建完成后，可以使用 virt-viewer 访问虚拟机，也可以使用 VNC 连接到 5911 以访问虚拟机，然后开始安装系统

## ChatGPT 通过 xml 文件推倒出来的 virt-install 参数

```
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


