---
title: Microsoft Management Console
linkTitle: Microsoft Management Console
weight: 1
---

# 概述

> 参考：
>
> - [Wiki, Microsoft Management Console](https://en.wikipedia.org/wiki/Microsoft_Management_Console)

**Microsoft Management Console(微软管理控制台，简称 MMC)** 是 Microsoft Windows 的一个组件，它为系统管理员和高级用户提供了一个用于配置和监控系统的界面。它于 1998 年首次与 Windows NT 4.0 的 Option Pack 一起推出，后来与 Windows 2000 及其后续版本预捆绑在一起。

# MSC

https://jingyan.baidu.com/article/7e440953dcc56e6ec1e2ef17.html

**MSC(Microsoft Snap-In Control)** 是 MMC(Microsoft Management Console) 用来添加/删除的嵌入式管理单元文件。通常通过 MMC 来管理，可点击“文件”菜单中的“添加/删除管理单元”操作来管理当前系统中已经安装的 MSC 文件。可以点击开始/运行，然后输入下列文件名就可以打开相应的控制窗口。

这些文件通常都以 `.msc` 为后缀。

除第三个文件外，其他均在 `C:\WINDOWS\system32` 文件夹下

- certmgr.msc
  - 作用：系统认证证书编辑。
- ciadv.msc
  - 作用：索引服务，链接文件*:\System Volume Information
- comexp.msc
  - 所在文件夹：C:\WINDOWS\system32\Com
  - 作用：组件服务，可以打开本地服务。
- compmgmt.msc
  - 作用：本地计算机硬件和服务管理，功能很强大。
- devmgmt.msc
  - 作用：设备管理器
- dfrg.msc
  - 作用：磁盘碎片整理程序
- diskmgmt.msc
  - 作用：磁盘管理器，可以修改盘符，格式化和分区等。
- eventvwr.msc
  - 作用：事件查看器
- fsmgmt.msc
  - 作用：共享文件夹管理
- gpedit.msc
  - 作用：组策略管理器，功能强大。TODO: 家庭版没有咋办？
- lusrmgr.msc
  - 作用：本地用户和组管理器
- ntmsmgr.msc
  - 作用：可移动存储管理器
- ntmsoprq.msc
  - 作用：可移动存储管理员操作请求
- perfmon.msc
  - 作用：性能察看管理器
- rsop.msc
  - 作用：各种策略的结果集
- secpol.msc
  - 作用：本地安全策略设置
- services.msc
  - 作用：各项本地服务管理器
- wmimgmt.msc
  - 作用：Windows管理体系结构（WMI）

# gpedit.msc

**Group policy edit(本地组策略编辑器，简称 gpedit)** 仅有专业版可用。

可以用来禁用小组件的重复安装
