---
title: Linux 图形处理
linkTitle: Linux 图形处理
weight: 1
---

# 概述

> 参考：
>
> - 

# Xorg, X11, Wayland? Linux Display Servers And Protocols Explained

原文链接：<https://linuxiac.com/xorg-x11-wayland-linux-display-servers-and-protocols-explained/>

**Have you ever wondered what exactly X server, Xorg, X11, Wayland and stuff like that does? Wayland vs Xorg, what is better? This guide is for you.**
You always stumble upon those terms, and know they have something to do regarding the graphics, but you’d like to know more.

### What is display server in Linux?

A display server  is a program whose primary task is to coordinate the input and output of its clients to and from the rest of the operating system, the hardware, and each other. The display server communicates with its clients over the display server protocol.
The display server is a key component in any graphical user interface, specifically the windowing system. It is the basic component of Graphical User Interface (GUI) which sits between the graphical interface and the kernel. So, thanks to a display server, you can use your computer with GUI. Without it, you would only be restricted to a command line interface.
It is very important to not confuse display server with desktop environment. The desktop environments ([Gnome](https://linuxiac.com/gnome-3-38-is-here-with-new-app-grid-and-better-performence/), KDE, Xfce, MATE, etc.) uses display server underneath it.
The display server communicates with its clients over the display server protocol. There are three display server protocols available in Linux. X11 and Wayland are two of them. The third, [Mir](https://mir-server.io/), is beyond the scope of this tutorial.

### X Window System, Xorg, X11, explained

X Window System, often referred to merely as X, is really old. First originating in 1984, it ended up being the default windowing system for most UNIX-like operating systems, including Linux.
**X.Org server** is the free and open-source implementation of the X Window System display server stewarded by the [X.Org Foundation](https://www.x.org/). It is an application that interacts with client applications via the X11 protocol to draw things on a display and to send input events like mouse movements, clicks, and keystrokes. Typically, one would start an X server which will wait for clients applications to connect to it. Xorg is based on a client/server model and thus allows clients to run either locally or remotely on a different machine.
If it’s not obvious, it’s implicit in the design of X11 that the application and the display don’t have to be on the same computer. At the time X was developed, it was very common that the X server would run on a workstation and the users would run applications on a remote computer with more processing power.
**X11** is a network protocol. It describes how messages are exchanged between a client (application) and the display (server). These messages typically carry primitive drawing commands like “draw a box”, “write these character at this position”, “the left mouse button has been clicked”, etc.
But X11 is old, and it was still a pile of hacks sitting on top of a protocol not truly overhauled for over 30 years. Most of the features that the X Server protocol provided were not used anymore. Pretty much all of the work that X11 did was redelegated to the individual applications and the window manager. And yet all of those old features are still there, weighing down on all of these applications, hurting performance and security.

### Wayland, the next-generation display server

**Wayland** was begun by Kristian Hogsberg, an X.Org developer, as a personal project in 2008. [It is a communication protocol](https://wayland.freedesktop.org/) that specifies the communication between a display server and its clients. Wayland is developed as a free and open-source community-driven project with the aim of replacing the X Window System (also known as X11, or Xorg ) with a modern, secure, and simpler windowing system.
In Wayland, the compositor is the display server. **Compositor**, is a window manager that provides applications with an off-screen buffer for each window. The window manager composites the window buffers into an image representing the screen and writes the result into the display memory.
The Wayland protocol lets the compositor send the input events directly to the clients and lets the client send the damage event directly to the compositor.
As in the X case, when the client receives the event it updates the user interface (UI) in response. But, in the Wayland rendering happens in the client, and the client just sends a request to the compositor to indicate the region that was updated.
Wayland’s main advantage over X is that it is starting from scratch. One of the main reasons for X’s complexity is that, over the years, its role has changed. As a result, today, X11 acts largely as “a really terrible” communications protocol between the client and the window manager.
Wayland is also superior when it comes to security. With X11, it’s possible to do something known as “keylogging” by allowing any program to exist in the background and read what’s happening with other windows open in the X11 area. With Wayland this simply won’t happen, as each program works independently.

### Conclusion

不过，X Window 系统相对于 Wayland 仍然有很多优点。尽管 Wayland 消除了 Xorg 的大部分设计缺陷，但它也有自己的问题。尽管 Wayland 项目已经运行了十多年，但并不是 100% 稳定。截至 2020 年，大多数 Linux 视频游戏和图形密集型应用程序仍然是为 X11 编写的。此外，许多闭源图形驱动程序（例如 NVIDIA GPU 的驱动程序）尚未提供对 Wayland 的完整支持。

X 无法持续，而 Wayland 在很多方面是一个改进。但目前绝大多数现有的原生应用程序都是为 Xorg 编写的。在这些应用程序全部移植之前，Xorg 需要继续维护。与 Xorg 相比，Wayland 目前还不太稳定。