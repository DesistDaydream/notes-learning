---
title: virt-manager
---

# 概述

> 参考
>
> - [GitHub 项目，virt-manager/virt-manager](https://github.com/virt-manager/virt-manager)
> - [官网](https://virt-manager.org/)

virt-manager 是一个图形化的应用程序，通过 libvirt 管理虚拟机。

virt-manager 提供了多个配套的工具

- virt-manager #
- virt-viewer # 是一个轻量级的 UI 界面，用于与虚拟客户操作系统的图形显示进行交互。它可以显示 VNC 或 SPICE，并使用 libvirt 查找图形连接详细信息。
- virt-install # 是一个命令行工具，它提供了一种将操作系统配置到虚拟机中的简单方法。
- virt-clone # 是一个用于克隆现有非活动客户的命令行工具。它复制磁盘映像，并使用指向复制磁盘的新名称、UUID 和 MAC 地址定义配置。
- virt-xml # 是一个命令行工具，用于使用 virt-install 的命令行选项轻松编辑 libvirt 域 XML。
- virt-bootstrap # 是一个命令行工具，提供了一种简单的方法来为基于 libvirt 的容器设置根文件系统。

virt-clone、virt-xml、virt-install 属于安装虚拟机的工具，通常都在 virtinst 包中
virt-manager、virt-viewer 属于图形化管理虚拟机的工具，通常都在 virt-manager 包中。

# virt-manager

## 使用 virt-manager 管理多台虚拟机

在一台机器上的 virt-manager 可以通过 add connection 管理其它宿主机上的虚拟机，但是前提是建立 ssh 的密钥认证，因为在 virt-manager 在通过 ssh 连接的时候，需要使用窗口模式输入密码，而一般情况下 ssh 是默认不装该组件的。如果不想添加密钥认证，那么安装 `ssh-askpass-gnome` 组件即可。

输入 virt-manager 打开管理界面。选择 File—Add Connecttion.. 勾选 Connect to remote host

依次填入文本框中内容如下：

- Hypervisor: QEMU/KVM
- Method:SSH
- Username:root
- Hostname:192.168.0.123（需要被操作的服务器地址）

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/abyoqm/1616123543706-5c47d241-1780-40d5-b04e-1cfd4b802d6c.png)
然后点击 connect 连接即可，此时就会出现另一台服务器上的虚拟机供操作。

# virt-viewer

可以显示虚拟机的图形控制台

# virt-install

> 参考：
> - [Manual(手册),virt-install(1)](https://man.cx/virt-install)

为各种虚拟化解决方案(KVM,XEN,Linux Container 等)创建 VM 并完成 GuestOS 安装

## Syntax(语法)

**virt-install --name NAME --memory MB STORAGE INSTALL \[OPTIONS]**

许多参数都有子选项，要查看与该参数相关联的子选项的完整列表使用例子中类似的命令，例如：virt-install --disk=?

有几个参数是在使用 libvirt 工具安装虚拟机时必须指定的：

- \--name is required
- \--memory amount in MiB is required
- \--disk storage must be specified (override with --disk none)

安装方式(--location URL, --cdrom CD/ISO, --pxe, --import, --boot hd|cdrom|...)

### OPTIONS

**GENERAL OPTIONS(通用选项)**

- **-n NAME|--name NAME** # 指定 VM 名称，必须要全局唯一
- **-r MEMORY|--ram MEMORY** # 设定 VM 的内存大小，单位为 MiB
- -**-vcpus=VCPU\[,maxvcpus=MAX]\[,sockets=NUM]\[,cores=NUM]\[,threads=NUM]** # 设定 VCPU 的个数，最大数，CPU 插槽数，内核数，线程数
- -**-cpu=CPU** # 设定 CPU 的型号及特性，如 coreduo 等，可以使用 kvm -cpu ? 来获取支持的 CPU 模式

**INSTALLATION OPTIONS(安装选项)**

- **-c, --cdrom=PATH **# 设定从光盘介质安装
- **-l LOCATION|--location OPTIONS** # 指定本地安装介质
- **--import** # 使用已经存在的磁盘镜像构建 VM。比如通过一个已经正常运行 VM 的文件

**Guest OS OPTIONS(虚拟机操作系统选项)**

- **--os-variant, --osinfo OSNAME** # 指定要虚拟机的操作系统的信息。常用来优化 virtio 等性能相关功能。OSNAME 可用的值有 win10、fedora32、ubuntu 等等，所有可用的值在哪查还不知道
  - 注意，--osinfo 是新版本的名称

**STORAGE OPTIONS(存储选项)**

- -**-disk /Some/Storage/Path\[,OPT1=VAL1]\[,OPT2=VAL2]\[,.....] **# 设定用作 GuestOS 存储的介质,后面是使用的路径\[选项=值(OPTIONs=VALUEs)]，可以在路径中使用 img 映像直接启动，以下是常用的 OPTIONS
  - device= #指定设备类型，如 cdrom，disk，默认为 disk
  - bus= #指定磁盘总线类型，如 ide、scsi、usb、virtio、xen
  - size= #指定新磁盘映像的大小，单位为 GB
  - cache= #指定缓存类型，如 none、writethrouth、writeback
  - format= #指定磁盘映像格式，如 qcow2、raw、vmdk 等

**NETWORKING OPTIONS(网络选项)**

- -**w NETWORK|--network 类型=名称,OPT1=VAL1,OTP2=VAL2,......** #有 4 中网络类型可供选择
  - bridge=BridgeName # 指定网络连接类型为 bridge 桥接模式，并选择用哪个桥
  - network=NMAE #
    - model= #指定 GuestOS 中的设备型号，如 e1000、virtio、rt18193 等
    - mac= #指定 mac 地址，默认使用随机 mac 地址

**GRAPHICS OPTIONS(图形选项)**

- **--graphics TYPE,OPT1=ARG1,OPT2=ARG2,...** # 指定 display VNC 的类型，可用的 TYPE 有 vnc、spice、none。其中 OPT=ARG 还可以为 graphics 指定更多的选项和值。下面是可用的 OPT：
  - listen=IP # 指定 vnc 监听的地址(默认值通常为 127.0.0.1。i.e.仅限本地主机使用)，如果配置 0.0.0.0，则可以被非宿主机的设备通过宿主机的 IP 与 PORT 来进行 vnc 访问
  - port=NUM # 指定访问该 VM 的 vnc 所使用的端口

**VIRTUALIZATION OPTIONS(虚拟化选项)**

**DEVICE OPTIONS(设备选项)**
指定文本控制台、声音设备、串行接口、并行接口、显示接口等

- **--serial TYPE,OPT1=VAL1,OPT2=VAL2,...** # 指定一个串行设备附加到 VM，TYPE 包括 pty(伪终端)等
- -**-console=** # 指定启动的控制台

**MISCELLANEOUS OPTIONS(其他选项)**

## EXAMPLE

使用本地文件当做磁盘镜像，VM 名字为 test，1024M 大小的内存，1 个 CPU，指定操作系统版本为 centos7，指定所使用的存储文件为/var/lib/images/test.qcow2，指定网络桥接到 br0 上、模式为 virtio，指定图形模式为 vnc、把 vnc 暴露到宿主机上、监听端口为 5910

```bash
virt-install --import --name test \
--memory 2048 --vcpus 2 \
--os-variant centos7.0 \
--disk /var/lib/libvirt/images/test.qcow2,size=20 \
--network bridge=br0,model=virtio \
--graphics vnc,listen=0.0.0.0,port=5910
```

### 使用 cdrom 安装系统

```bash
virt-install --cdrom /root/iso/CentOS-7-x86_64-DVD-1908.iso --name centos7 \
--memory 4096 --vcpus 4 \
--os-variant centos7.0 \
--disk /var/lib/libvirt/images/centos7.qcow2,size=100,bus=virtio \
--network bridge=br0,model=virtio \
--graphics vnc,listen=0.0.0.0,port=5911
```

创建完成后，可以使用 virt-viewer 访问虚拟机，也可以使用 VNC 连接到 5911 以访问虚拟机，然后开始安装系统

# virt-clone

> 参考：
> - [Manual(手册),virt-clone(1)](https://man.cx/virt-clone)

## Syntax(语法)

**virt-clone \[OPTION]...**
