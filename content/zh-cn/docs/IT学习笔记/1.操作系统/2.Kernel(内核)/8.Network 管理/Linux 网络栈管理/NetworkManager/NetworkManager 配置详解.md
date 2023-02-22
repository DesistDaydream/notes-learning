---
title: NetworkManager 配置详解
---

# 概述

> 参考：
> - [Manual(手册),NetworkManager.conf(5)](https://networkmanager.dev/docs/api/latest/NetworkManager.conf.html)
> - 在 [GNOME 开发者中心官网](https://developer-old.gnome.org/NetworkManager/)中，也可以查到 Manual
> - <https://wiki.gnome.org/Projects/NetworkManager/DNS>
> - <https://cloud.tencent.com/developer/article/1710514>

NetworkManager 的配置文件是 INI 格式的，由 **Sections(部分)** 和 **Key/Value Pairs(键/值对)** 组成。

可用的 Sections 有如下几种：

- main #
- keyfile # 用于配置 keyfile 插件。通常只在不使用任何特定 Linux 发行版的插件时才进行配置。
- ifupdown #
- logging # 控制 NetworkManager 的日志记录。此处的任何设置都被 --log-level 和 --log-domains 命令行选项覆盖。
- connection #
- device #
- connectivity #
- global-dns #
- global-dns-domain #
- .config #

# main 部分

**dns=\<MODE>** # 设置 DNS 处理模式。
可用的模式有如下几种：

- default #
- dnsmasq #
- systemd-resolved #
- unbound #
- **none** # NetworkManager 程序不会修改 resovl.conf 文件。

**plugins=\<STRING>** # 设置控制 Connections 配置文件的插件，多个插件以 **`,`** 分隔。

注意，NetworkManager 原生的 keyfile 插件始终附加到此列表的末尾。也就意味着，NetworkManager 始终都会加载 keyfile 插件。

# keyfile 部分

**path=\<STRING>** # 读取和存储 Connection 配置文件的目录`默认值：/etc/NetworkManager/system-connections`

**unmanaged-devices=\<STRING>** # 指定 keyfile 不管理的网络设备

# logging 部分
