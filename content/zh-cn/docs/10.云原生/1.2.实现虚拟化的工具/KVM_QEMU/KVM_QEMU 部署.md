---
title: KVM/QEMU 部署
weight: 3
---

# 概述

> 参考：
> - 官方文档，安装  TODO: 官方文档里没有教安装 qemu-system 的地方呀~o(╯□╰)o


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

## qemu-system 与 CPU 架构的说明

### ARM

QEMU 可以模拟 32 位和 64 位 Arm CPU。使用 qemu-system-aarch64 可执行文件模拟 64 位 Arm 机器。您可以使用 qemu-system-arm 或 qemu-system-aarch64 来模拟 32 位 Arm 机器：通常，适用于 qemu-system-arm 的命令行在与 qemu-system-aarch64 一起使用时表现相同。