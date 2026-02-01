---
title: Windows MGMT
linkTitle: Windows MGMT
weight: 1
---

# 概述

> 参考：
>
> - 

休眠

- [这么多年了，为啥Windows笔记本们连个休眠都做不好？【差评君】](https://www.bilibili.com/video/BV1k46GY5Eg1)

# 快捷键

[知乎，Windows 用 Spy++ 查找全局快捷键冲突](https://zhuanlan.zhihu.com/p/704643938)

Spy++ 随 Visual Studio 安装而存在，社区版即可。使用 everything 工具搜 `spyxx`，使用 `spyxx_amd64.exe`

- 关闭所有窗口
- 打开 `监视 — 日志消息`
- 在 “窗口” 标签页，勾选 `其他窗口 — 系统中的所有窗口`
- 在 “消息” 标签页，点击 `全部清除`，勾选 `消息组 - 键盘`（主要是其中的 `WM_HOTKEY` 消息）
- 点击 “确定” 后，创建出来一个新的监视窗口。

将焦点移动到 Spy++ 程序之外，按任意快捷键

---

使用 OpenARK，进入内核模式，列出所有已知快捷键

---

我的电脑 `ctrl + ,` 快捷键无法设置也无法使用，哪怕在安全模式下也不行。任何组合加逗号都行，唯独 ctrl 不行。。。

# 通知

> 参考：
>
> - [官方文档，通知与勿扰](https://support.microsoft.com/en-us/windows/notifications-and-do-not-disturb-in-windows-feeca47f-0baf-5680-16f0-8801db1a8466)

Windows 通知是 Toast 类型

[B 站，输入法利用系统通知弹广告，还禁止关闭？倒查！](https://www.bilibili.com/video/BV1ez6nBmEiy)

注册表项

- **ShowInSettings**(0|1) # 是否在 设置 - 系统 - 通知 中显示。0 不显示，1 显示。
    - 该项无法在官方文档中找到说明

