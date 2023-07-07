---
title: cobbler 命令说明
---

## distro

- cobbler distro report --name=CentOS7-1810 # 查看安装镜像文件信息

## import

cobbler import [OPTIONS]

OPTIONS

- --path # 指定制作部署系统时的镜像所用到的光盘镜像的路径
- --name # 为安装源定义一个名字，指定部署系统所用的镜像名
- --arch # 指定安装源是 32 位、64 位、ia64, 目前支持的选项有: x86│x86_64│ia64

EXAMPLE

- cobbler import --path=/mnt/ --name=CentOS7-1810 --arch=x86_64

## profile

- cobbler profile report --name=CentOS7-x86_64 # 查看指定的 profile 设置

## system

- cobbler system list # 列出

## ksvalidator FILE

用于测试指定 FILE 的 kickstart 语法是否正确
