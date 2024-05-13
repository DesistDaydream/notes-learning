---
title: "部署 KVM/QEMU 虚拟化环境"
linkTitle: "部署 KVM/QEMU 虚拟化环境"
weight: 20
---

# 概述

> 参考：
> 
> - [RedHat 官方文档，9-配置和管理虚拟化-开启虚拟化](https://access.redhat.com/documentation/en-us/red_hat_enterprise_linux/9/html/configuring_and_managing_virtualization/assembly_enabling-virtualization-in-rhel-9_configuring-and-managing-virtualization)
> - [Ubuntu 官方文档，虚拟化介绍](https://ubuntu.com/server/docs/virtualization-introduction)

# 前期准备

查看 CPU 是否支持 KVM，筛选出来相关信息才可以正常使用 KVM

- egrep "(svm|vmx)" /proc/cpuinfo

# 安装虚拟化组件

## CentOS

- yum group install -y 'Virtualization Host' # 安装虚拟化组
  - 若安装完成后，模块未装载，则手动装载 KVM 模块
    - modprobe kvm
    - modprobe kvm-intel
  - 验证系统是否已经准备好成为虚拟化主机
    - virt-host-validate
- 启动 libvirt 服务
  - systemctl enable libvirtd --now
- 创建连接使用命令
  - ln -sv /usr/libexec/qemu-kvm /usr/bin/
- 安装 X 服务端程序
  - yum install -y xorg-x11-xauth xorg-x11-server-utils
- 安装图形管理工具
  - yum install virt-manager -y
- 安装 qemu 以模拟 I/O 设备
  - yum install qemu-system-x86 qemu-img -y
- 安装 virt 安装命令
  - yum install virt-install -y
- 安装虚拟机文件系统的管理工具
  - yum install libguestfs-tools -y

## Ubuntu

检查环境

- sudo apt update
- sudo apt install -y cpu-checker
- kvm-ok

安装虚拟化环境

- sudo apt install -y qemu-kvm libvirt-daemon-system libvirt-daemon libvirt-clients

安装虚拟机文件系统的管理工具

- apt install libguestfs-tools -y
- apt install virt-manager -y

# 安装 VNC 与桌面(可选)

## CentOS

- yum install -y tigervnc-server # 安装 vnc 服务端
- yum groups install -y 'GNOME Desktop' #

# 其他

yum -y install mesa-libGLES-devel.x86_64 mesa-dri-drivers

若不安装这两个包，当使用 virt-mangaer 工具是，可能会出现如下报错：

```bash
libGL error: unable to load driver: swrast_dri.so
libGL error: failed to load driver: swrast
```
