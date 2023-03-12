---
title: KVM/QEMU 部署
weight: 3
---

# 概述

> 参考：
> - <https://access.redhat.com/documentation/en-us/red_hat_enterprise_linux/8/html-single/configuring_and_managing_virtualization/index#enabling-virtualization-in-rhel8_virt-getting-started>
> - [Ubuntu 官方文档，虚拟化介绍](https://ubuntu.com/server/docs/virtualization-introduction)

# 前期准备

查看 CPU 是否支持 KVM，筛选出来相关信息才可以正常使用 KVM

- egrep "(svm|vmx)" /proc/cpuinfo

# 安装虚拟化组件

## CentOS

yum install qemu-kvm

## Ubuntu

检查环境

- sudo apt update
- sudo apt install -y cpu-checker
- kvm-ok

sudo apt install qemu-system
