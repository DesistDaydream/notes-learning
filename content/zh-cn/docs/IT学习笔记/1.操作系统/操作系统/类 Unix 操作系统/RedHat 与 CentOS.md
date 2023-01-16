---
title: RedHat 与 CentOS
---

# 概述

> 参考：
> - [RedHat 官方文档](https://access.redhat.com/products/red-hat-enterprise-linux/#knowledge)(在这里点击 Product Documentation)
> - [RedHat7 生产环境文档](https://access.redhat.com/documentation/en-us/red_hat_enterprise_linux/7)
> - [RedHat8 生产环境文档](https://access.redhat.com/documentation/en-us/red_hat_enterprise_linux/8)
> - [CentOS 官方，法律](https://www.centos.org/legal/)

[CentOS: 永远有多远就离它多远](https://mp.weixin.qq.com/s/heX7Qtc7Fizx43EgGkIiMQ)
[CentOS7 好日子到头了，如何优雅的抛弃 CentOS7？](https://mp.weixin.qq.com/s/DUUYW_OBV_wUu1wZaP6gAg)
[CentOS 8 退役倒计时，开发者们又吵起来了](https://mp.weixin.qq.com/s/FMvNx-kzz7DZZqGGpxjbuw)

CentOS 居然还用 python2
装 Python3 很费劲
装 python-libvirt 很费劲

![41212703dee962f84b5c4a49a80707d.jpg](https://notes-learning.oss-cn-beijing.aliyuncs.com/wlni0v/1654840849145-f536a3be-b969-40b8-813d-959985b4f429.jpeg)

# 安装 CentOS

> 参考：
> - [GitHub 项目，rhinstaller/anaconda](https://github.com/rhinstaller/anaconda)(RedHat 相关发行版的操作系统安装器)

RedHat 系列的 .iso 安装文件中包含了 Anaconda 安装器。

## 注意事项

/bin,/dev,/sbin,/etc,/lib,/root, /mnt, lost+found, /proc 这些目录不能创建单独的分区并挂载，只能创建一个 / 以包含这些目录

- <https://unix.stackexchange.com/questions/121318/this-mount-point-is-invalid-the-root-directory-must-be-on-file-system>
- 代码：[https://github.com/rhinstaller/anaconda/blob/rhel6-branch/storage/**init**.py#L1084](https://github.com/rhinstaller/anaconda/blob/rhel6-branch/storage/__init__.py#L1084)
- 高于 6 版本的分之代码将这个行为封装了
  - <https://github.com/rhinstaller/anaconda/blob/rhel-9/pyanaconda/modules/storage/checker/utils.py#L31>

# 关联文件

**/etc/sysconfig/\*** # Red Hat Linux 发行版的各种系统配置文件

# CentOS 法律

![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/wlni0v/1647171407465-5e7ad5f3-913d-4c93-a893-e3334b09bfbb.png)

# Centos 变为上游项目

Centos 的创始人新建了一个 [**Rocky 项目**](https://github.com/rocky-linux/rocky)，旨在作为 RedHat 下游 Linux 发行版

Frequently Asked Questions

> **Q:** What do you mean, "CentOS has shifted direction?"

The CentOS project recently announced a shift in strategy for CentOS. Whereas previously CentOS existed as a _downstream_ build of its upstream vendor (it receives patches and updates after the upstream vendor does), it will be shifting to an _upstream_ build (testing patches and updates _before_ inclusion in the upstream vendor).

Additionally, support for CentOS Linux 8 has been cut short, from May 31, 2029 to December 31, 2021.

> **Q:** So where does Rocky Linux come in?

Rocky Linux aims to function as a downstream build as CentOS had done previously, building releases after they have been added to the upstream vendor, not before.

> **Q:** When will it be released?

There is not currently an ETA for release.

> **Q:** What is the vision for Rocky Linux?

A **solid**, **stable**, and **transparent** alternative for production environments, developed _by_ the community _for_ the community.

> **Q:** Who drives Rocky Linux?

We all do, Rocky Linux is a community-driven project and always will be. Rocky Linux will not be sold or driven by corporate interest.

> **Q:** How can I get involved?

Please view the contributing section below.

# Centos Stream 问题汇总

CentOS Stream 使用了别人已经用了很久的 system-resolved.service 服务，但是从 centos8 升级到 centos stream 后，服务有了，但是却并没有自动启动该服务，也就导致了没有 /run/systemd/resolve/resolv.conf 文件，很多程序在发现 system-resolved.service 后，会去读取这个文件，比如 kubelet 程序。
