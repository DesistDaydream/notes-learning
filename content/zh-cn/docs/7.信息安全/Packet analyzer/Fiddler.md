---
title: Fiddler
linkTitle: Fiddler
weight: 20
---

# 概述

> 参考：
>
> - [官网](https://www.telerik.com/fiddler)
> - [Wiki, Telerik](https://en.wikipedia.org/wiki/Telerik)

Fiddler 在 2012 年被 Telerik 收购

# 安装

Fiddler Classic 版本可免费使用

# HTTPS 抓包

打开 Tools - Options

在 HTTPS 标签中，勾选 `Capture HTTPS CONNECTs` 和 `Decrypt HTTPS traffic`，Windows 会自动安装 Root 证书。

> 在 Actions 中可以执行 重置证书、下载证书、等等 操作

![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/fiddler/capture_https_1.png)

在 Connections 标签中，勾选 `Allow remote computers to connect` 以便通过 PC 的 Fiddler 抓取移动设备的包。在这里还可以配置 Fiddler 的监听端口

> 取消勾选 `Act as system proxy on startup` 以便 Fiddler 启动时不要配置系统代理

![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/fiddler/202401191101395.png)

## IOS

### 安装证书

为无线连接配置手动代理，Fiddler 默认监听在 :8888

IOS 访问 `IP:8888` ，点击 FiddlerRoot certificate 下载证书并安装

设置 —— 通用 —— 关于本机 —— 证书信任设置，开启信任证书

# Fiddler Add-ons

https://www.telerik.com/fiddler/add-ons
