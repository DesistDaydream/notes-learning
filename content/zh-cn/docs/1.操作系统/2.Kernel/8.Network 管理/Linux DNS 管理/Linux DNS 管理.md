---
title: Linux DNS 管理
---

# 概述

> - [公众号,重新夺回 /etc/resolv.conf 的控制权](https://mp.weixin.qq.com/s/L9TpAFqT-5V7ppEGdT0cnw)

在 Linux 中，进行域名解析工作的是 reslover(解析器)。

**reslover(解析器)** 是 C 语言库中用于提供 DNS 接口的程序集，当某个进程调用这些程序时将同时读入 reslover 的配置文件，这个文件具有可读性并且包含大量可用的解析参数

# Linux DNS 配置

**/etc/resolv.conf** # reslover 配置文件
**/etc/hosts **# 更改本地主机名和 IP 的对应关系，用于解析指定域名
例：当 ping TEST-1 时，则 ping 192.168.2.3

    127.0.0.1   localhost localhost.localdomain localhost4 localhost4.localdomain4
    ::1         localhost localhost.localdomain localhost6 localhost6.localdomain6
    192.168.2.3 TEST-1

**\*/nsswitch.conf **# 名称服务切换配置。GUN C 库(glibc) 和 某些其他应用程序使用该配置文件来确定从哪些地方获取解析信息。比如是否要读取 /etc/hosts 文件

- 该文件属于 glibc 包中的一部分。但是由于 CentOS 与 Ubuntu 中 glibc 的巨大差异，该文件所在路径也不同：
  - CentOS 在 **/etc/nsswitch.conf**
  - Ubuntu 在 **/usr/share/libc-bin/nsswitch.conf**
