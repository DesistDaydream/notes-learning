---
title: cobbler命令说明
---

## distro

1. cobbler distro report --name=CentOS7-1810 #查看安装镜像文件信息

<br />
## import
<br />cobbler import [OPTIONS]

OPTIONS

1. \--path #指定制作部署系统时的镜像所用到的光盘镜像的路径

2. \--name #为安装源定义一个名字，指定部署系统所用的镜像名

3. \--arch #指定安装源是 32 位、64 位、ia64, 目前支持的选项有: x86│x86_64│ia64

EXAMPLE

1. cobbler import --path=/mnt/ --name=CentOS7-1810 --arch=x86_64 #

<br />
## profile

1. cobbler profile report --name=CentOS7-x86_64 #查看指定的 profile 设置

<br />
## system

1. cobbler system list #列出

<br />
## ksvalidator FILE
<br />用于测试指定FILE的kickstart语法是否正确<br />
