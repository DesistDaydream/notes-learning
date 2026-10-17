---
title: Windowing system
created: 2026-10-17T14:30
weight: 2
---

# 概述

> 参考：
>
> - [Wiki, Windowing system](https://en.wikipedia.org/wiki/Windowing_system)

**Windowing system(窗口系统)** 是一个软件套件，可以让应用程序以 GUI 的形式展示。

窗口系统的主要组件称为 **Display server(显示服务器)**，任何通过 Display server 以 GUI 形式运行的应用程序，都是该 Display server 的 **Client(客户端)**

Display server 与 Client 之间通过 [API](/docs/2.编程/API/API.md) 或 [Communication protocol](/docs/4.数据通信/Protocol/Communication%20protocol.md) 进行通信（这种协议通常称为 [Display server protocol](#display%20server%20protocol)）。Display server 是 Client 与用户之间的中介。

> [!Tip] 虽然叫 Display server，但是通常还包含输入设备（键盘、鼠标、etc.）的管理

Display server 从内核接收所有输入：这些输入由内核从所连接的各种输入设备（如键盘、指点设备或触摸屏）采集而来；随后 Display server 将输入转发给正确的 Client。

Display server 还负责把 Client 的输出送到计算机显示器。

> 声音输出通常不由显示服务器管理，但音量一般通过 GUI 小程序来调节；而决定哪些应用程序位于最上层（即窗口叠放次序）的，正是显示服务器。

Windowing system 让计算机用户能够同时使用多个程序。每个程序都在自己的窗口中呈现其 GUI，这个窗口一般是屏幕上的一个矩形区域。

## 反直觉的 C/S 架构

> [!Attention] 在 Display server 环境中的 C/S 架构与通常意义上的 C/S 架构不太一样，有点反着来的感觉。
>
> - Display server 部署在个人电脑上
> - Client 部署在服务器上
>
> 当然了，没有服务器的话，Client 也可以装在 PC 上嘛，就像把个人电脑当服务器似的
>
> 人们通常认为 ”Server“ 应该是大型远程计算机。而 Display server 却是一个小型本地系统，其大部分客户端反倒运行在更大的中央服务器上。这是因为 Display server 为这些客户端提供显示和输入设备的**服务**，或者说提供共享资源。

这种错位感觉的来源主要因为 *我以人的视角作为参考系*，但是对 C/S 架构来说，**共享资源** 所在位置才是参考系：

- 在一般 C/S 架构中，共享资源是 HTML 之类的静态文件、程序计算服务、远程文件、etc.
- 在 Windowing system 中，共享资源是 显示硬件（i.e. Display这里的显示硬件包括：键鼠和显示器）

而通常情况下，Display server 一般都是与 Client 一起装在个人 PC 上的，我们日常很难感知到这种 C/S 架构。像 浏览器、Office、Obsidian、etc. 都可以理解为一种 Client。<font color="#ff0000">当我们把这些图形应用程序放到远端，并通过网络协议连接，只需要保证个人 PC 有实现了 Display server 的程序，就可以直接以 GUI 形式打开这些程序</font>。此时远端的应用程序与本端的图形服务合在一起，就形成了我们在个人 PC 上安装并打开一个应用程序的效果。

# Display server protocol

X11

[Wayland](/docs/11.多媒体/图形处理/Linux%20图形处理/freedesktop/Wayland/Wayland.md)

# 各类系统常见的 Display server 及连接方式

TODO: 来自 DSH + DS API，待逐条校验

| 系统              | Display server                         | Client 如何连它                                                                                                                   | 跨机方式                                                               |
| --------------- | -------------------------------------- | ----------------------------------------------------------------------------------------------------------------------------- | ------------------------------------------------------------------ |
| Linux + X11     | Xorg（X server）                         | X11 协议；Unix socket `/tmp/.X11-unix/X0`（现代发行版默认用 `-nolisten tcp` 关闭 TCP）；`DISPLAY=:0`；`~/.Xauthority` 里的 MIT-MAGIC-COOKIE-1 授权 | 原生支持：`ssh -X`（转发 X11 通道，DISPLAY 变 `localhost:10.0`）、XDMCP（udp 177） |
| Linux + Wayland | Wayland compositor（Mutter、KWin、sway 等） | Wayland 协议；Unix socket `$XDG_RUNTIME_DIR/wayland-0`；X11 应用由 Xwayland 代管                                                       | 协议本身不跨网络；要跨机需先用 waypipe 做中继，或改用 RDP/VNC 把整个桌面传过去                   |
| Windows         | `win32k.sys`（内核态）+ `dwm.exe`（DWM，合成器）  | GDI / DirectX 系统调用，没有协议，也没有 socket                                                                                            | 应用层没有通道；只能使用整个桌面的 RDP；单应用需 RAIL（RemoteApp，且需要 RDS 或第三方程序作为宿主）      |
| macOS           | WindowServer（即 Quartz Compositor）      | CoreGraphics / Metal，经 Mach IPC（mach port）                                                                                    | 没有单应用通道；整个桌面用 Screen Sharing（ARD/VNC）                              |

**只有 X11 把 “Display server” 做成了可跨网络的独立进程**，因此只有它天然存在 “X11 转发” 这种能力；其它的 Display server 要么是内核组件，要么是只在本机通过 IPC 通信的进程，远程只能整套桌面（RDP、VNC 或 ARD），或者依靠 waypipe 这类中继。
