---
title: 硬核致敬Linux ！30岁生日快乐！
---

[硬核致敬 Linux ！30 岁生日快乐！](https://mp.weixin.qq.com/s/cE4x63tYxoqrDinifeWqeg)

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/1ef62d82-d29c-457a-b63f-3d4dfa09d652/640)

1991 年 8 月 25 日，21 岁的 Linus Torvalds（以下简称 Linus）做了一个免费的操作系统 “Linux”，并在这一天向外界公布这个由“业余爱好” 主导的个人项目；如今，全球超级计算机 500 强和超过 70% 的智能手机都在运行 Linux，因此，8 月 25 日也被许多 Linux 的爱好者视为 Linux 真正的诞生日期。

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/1ef62d82-d29c-457a-b63f-3d4dfa09d652/640)

30 年前，Linus Torvalds 首次发布 Linux 内核时还是赫尔辛基大学的一名 21 岁学生。他的声明是这样开始的，“我正在做一个（免费的）操作系统（只是一个爱好，不会很大和专业......）”。三十年后，排名前 500 的超级计算机都在运行 Linux，所有智能手机的 70% 以上都是如此。Linux 显然既庞大又专业。

三十年来，Linus Torvalds 领导了 Linux 内核开发，激励了无数其他开发人员和开源项目。2005 年，Linus 还创建了 Git 来帮助管理内核开发过程，此后它成为最受欢迎的版本控制系统，受到无数开源和专有项目的信赖。

### Linux 历史

**OS 史前历史**

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/1ef62d82-d29c-457a-b63f-3d4dfa09d652/640)

**Linux 的历史**

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/1ef62d82-d29c-457a-b63f-3d4dfa09d652/640)

### Linux 系统

**Linux 系统软件架构**

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/1ef62d82-d29c-457a-b63f-3d4dfa09d652/640)

Linux 系统由硬件、kernel、系统调用、shell、c 库、应用程序组成，架构层次分明，Linux 内的各种层功能独立，程序在用户空间和内核空间之间的分离，能支持更多应用。

| 用户模态                                                                               | **用户应用**   | 例如：Bash，LibreOffice，GIMP，Blender，0 A.D.，Mozilla Firefox 等                       |            |                |                |            |
| -------------------------------------------------------------------------------------- | -------------- | ---------------------------------------------------------------------------------------- | ---------- | -------------- | -------------- | ---------- |
|                                                                                        | 低层系统构件   | **系统守护进程**：                                                                       |
| systemd，runit，logind，networkd，PulseAudio 等                                        |                | **窗口系统**：                                                                           |
| X11，Wayland，SurfaceFlinger(Android)                                                  | **其他库**：   |
| GTK+, Qt, EFL, SDL, SFML, FLTK, GNUstep 等                                             | **图形**：     |
| Mesa，AMD Catalyst 等                                                                  |
|                                                                                        | **C 标准库**   | open()，exec()，sbrk()，socket()，fopen()，calloc()，... (直到 2000 个子例程)            |
| glibc 目标为 POSIX/SUS 兼容，musl 和 uClibc 目标为嵌入式系统，bionic 为 Android 而写等 |                |                                                                                          |            |                |
| 内核模态                                                                               | **Linux 内核** | stat, splice, dup, read, open, ioctl, write, mmap, close, exit 等（大约 380 个系统调用） |
| Linux 内核系统调用接口（SCI，目标为 POSIX/SUS 兼容）                                   |                |                                                                                          |            |                |
|                                                                                        |                | 进程调度子系统                                                                           | IPC 子系统 | 内存管理子系统 | 虚拟文件子系统 | 网络子系统 |
|                                                                                        |                | 其他构件：ALSA，DRI，evdev，LVM，device mapper，Linux Network Scheduler，Netfilter       |
| Linux 安全模块：SELinux，TOMOYO，AppArmor, Smack                                       |                |                                                                                          |            |                |
| 硬件（CPU，内存，数据存储设备等。）                                                    |                |                                                                                          |            |                |                |            |

**Linux 内核代码架构**

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/1ef62d82-d29c-457a-b63f-3d4dfa09d652/640)

Linux 代码包含核心几个子系统，比如**内存子系统**，**I/O 子系统**，**CPU(调度）子系统**，**设备驱动子系统**，**网络子系统**，**虚拟文件子系统**等。这里简单介绍一些比较重要的子系统。

### 调度子系统

**进程调度**是 Linux 内核中最重要的子系统，它主要提供对 CPU 的访问控制。因为在计算机中，CPU 资源是有限的，而众多的应用程序都要使用 CPU 资源，所以需要 “进程调度子系统” 对 CPU 进行调度管理。

**进程调度子系统**包括 4 个子模块（见下图），它们的功能如下：

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/1ef62d82-d29c-457a-b63f-3d4dfa09d652/640)

1. Scheduling Policy，实现进程调度的策略，它决定哪个（或哪几个）进程将拥有 CPU。
2. Architecture-specific Schedulers，体系结构相关的部分，用于将对不同 CPU 的控制，抽象为统一的接口。这些控制主要在 suspend 和 resume 进程时使用，牵涉到 CPU 的寄存器访问、汇编指令操作等。
3. Architecture-independent Scheduler，体系结构无关的部分。它会和 “Scheduling Policy 模块” 沟通，决定接下来要执行哪个进程，然后通过“Architecture-specific Schedulers 模块”resume 指定的进程。
4. System Call Interface，系统调用接口。进程调度子系统通过系统调用接口，将需要提供给用户空间的接口开放出去，同时屏蔽掉不需要用户空间程序关心的细节。

### 内存子系统

**内存管理**同样是 Linux 内核中最重要的子系统，它主要提供对内存资源的访问控制。Linux 系统会在硬件物理内存和进程所使用的内存（称作虚拟内存）之间建立一种映射关系，这种映射是以进程为单位，因而不同的进程可以使用相同的虚拟内存，而这些相同的虚拟内存，可以映射到不同的物理内存上。

**内存管理子系统**包括 3 个子模块（见下图），它们的功能如下：

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/1ef62d82-d29c-457a-b63f-3d4dfa09d652/640)

1. Architecture Specific Managers，体系结构相关部分。提供用于访问硬件 Memory 的虚拟接口。
2. Architecture Independent Manager，体系结构无关部分。提供所有的内存管理机制，包括：以进程为单位的 memory mapping；虚拟内存的 Swapping。
3. System Call Interface，系统调用接口。通过该接口，向用户空间程序应用程序提供内存的分配、释放，文件的 map 等功能。

### 虚拟文件子系统（Virtual Filesystem, VFS）

传统意义上的文件系统，是一种存储和组织计算机数据的方法。它用易懂、人性化的方法（文件和目录结构），抽象计算机磁盘、硬盘等设备上冰冷的数据块，从而使对它们的查找和访问变得容易。因而文件系统的实质，就是 “存储和组织数据的方法”，文件系统的表现形式，就是 “从某个设备中读取数据和向某个设备写入数据”。

随着计算机技术的进步，存储和组织数据的方法也是在不断进步的，从而导致有多种类型的文件系统，例如 FAT、FAT32、NTFS、EXT2、EXT3 等等。而为了兼容，操作系统或者内核，要以相同的表现形式，同时支持多种类型的文件系统，这就延伸出了**虚拟文件系统（VFS）**的概念。VFS 的功能就是管理各种各样的文件系统，屏蔽它们的差异，以统一的方式，为用户程序提供访问文件的接口。

我们可以从磁盘、硬盘、NAND Flash 等设备中读取或写入数据，因而最初的文件系统都是构建在这些设备之上的。这个概念也可以推广到其它的硬件设备，例如内存、显示器（LCD）、键盘、串口等等。我们对硬件设备的访问控制，也可以归纳为读取或者写入数据，因而可以用统一的文件操作接口访问。Linux 内核就是这样做的，除了传统的磁盘文件系统之外，它还抽象出了设备文件系统、内存文件系统等等。这些逻辑，都是由 VFS 子系统实现。

**VFS 子系统**包括 6 个子模块（见下图），它们的功能如下：

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/1ef62d82-d29c-457a-b63f-3d4dfa09d652/640)

1. Device Drivers，设备驱动，用于控制所有的外部设备及控制器。由于存在大量不能相互兼容的硬件设备（特别是嵌入式产品），所以也有非常多的设备驱动。因此，Linux 内核中将近一半的 Source Code 都是设备驱动，大多数的 Linux 底层工程师（特别是国内的企业）都是在编写或者维护设备驱动，而无暇估计其它内容（它们恰恰是 Linux 内核的精髓所在）。
2. Device Independent Interface， 该模块定义了描述硬件设备的统一方式（统一设备模型），所有的设备驱动都遵守这个定义，可以降低开发的难度。同时可以用一致的形势向上提供接口。
3. Logical Systems，每一种文件系统，都会对应一个 Logical System（逻辑文件系统），它会实现具体的文件系统逻辑。
4. System Independent Interface，该模块负责以统一的接口（快设备和字符设备）表示硬件设备和逻辑文件系统，这样上层软件就不再关心具体的硬件形态了。
5. System Call Interface，系统调用接口，向用户空间提供访问文件系统和硬件设备的统一的接口。

### 网络子系统（Net）

**网络子系统**在 Linux 内核中主要负责管理各种网络设备，并实现各种网络协议栈，最终实现通过网络连接其它系统的功能。在 Linux 内核中，网络子系统几乎是自成体系，它包括 5 个子模块（见下图），它们的功能如下：

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/1ef62d82-d29c-457a-b63f-3d4dfa09d652/640)

1. Network Device Drivers，网络设备的驱动，和 VFS 子系统中的设备驱动是一样的。
2. Device Independent Interface，和 VFS 子系统中的是一样的。
3. Network Protocols，实现各种网络传输协议，例如 IP, TCP, UDP 等等。
4. Protocol Independent Interface，屏蔽不同的硬件设备和网络协议，以相同的格式提供接口（socket)。
5. System Call interface，系统调用接口，向用户空间提供访问网络设备的统一的接口。

Linux 内核版本时间线：

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/1ef62d82-d29c-457a-b63f-3d4dfa09d652/640)

**Linux 内核支持各种硬件架构**

**Linux 内核**最成功的地方之一就是支持各种硬件架构，为软件提供了公共的平台：

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/1ef62d82-d29c-457a-b63f-3d4dfa09d652/640)

基于 Linux 的系统是一个模块化的类 Unix 操作系统。**Linux 操作系统**的大部分设计思想来源于 20 世纪 70 年代到 80 年代的 Unix 操作系统所建立的基本设计思想。Linux 系统使用宏内核，由 Linux 内核负责处理进程控制、网络，以及外围设备和文件系统的访问。在系统运行的时候，设备驱动程序要么与内核直接整合，要么以加载模块形式添加。

**Linux 具有设备独立性**，它内核具有高度适应能力，从而给系统提供了更高级的功能。GNU 用户界面组件是大多数 Linux 操作系统的重要组成部分，提供常用的 C 函数库，Shell，还有许多常见的 Unix 实用工具，可以完成许多基本的操作系统任务。大多数 Linux 系统使用的图形用户界面建立在 X 窗口系统之上，由 X 窗口 (XWindow) 系统通过软件工具及架构协议来建立操作系统所用的图形用户界面.

**基于 Linux 内核各种衍生 OS 系统**

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/1ef62d82-d29c-457a-b63f-3d4dfa09d652/640)

各种发行版本

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/1ef62d82-d29c-457a-b63f-3d4dfa09d652/640)

当前比较流行发行版是：**Debian**、**Ubuntu**、**Fedora**、**CentOS**、**Arch Linux**和**openSUSE**等，每个发行版都有自己优势地方，都有一批忠实用户。

**基于 Linux 内核著名 OS**

**Android**

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/1ef62d82-d29c-457a-b63f-3d4dfa09d652/640)

**Android**（读音：英：\['ændrɔɪd]，美：\[ˈænˌdrɔɪd]），中文用户多以非官方名称 “安卓” 称之，是一个基于 Linux 内核与其他开源软件的开放源代码的移动操作系统，Android 的内核是根据 Linux 内核的长期支持的分支，具有典型的 Linux 调度和功能。截至 2018 年，Android 的目标是 Linux 内核的 4.4、4.9 或是 4.14 版本。

**ChromeOS**

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/1ef62d82-d29c-457a-b63f-3d4dfa09d652/640)

**Chrome OS**  是由 Google 设计基于 Linux 内核的操作系统，并使用 Google Chrome 浏览器作为其主要用户界面。因此，Chrome OS 主要支持 Web 应用程序\[6]，2016 年起开始陆续兼容 Android 应用程序（可通过 Google Play 商店下载）和 Linux 应用程序。

**鸿蒙 OS**

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/1ef62d82-d29c-457a-b63f-3d4dfa09d652/640)

**鸿蒙**（**HarmonyOS**，开发代号 Ark\[1]）是华为自 2012 年开发的一款可兼容 Android 应用程序的跨平台分布式操作系统\[2]。系统性能包括利用 “分布式” 技术将各款设备融合成一个“超级终端”，便于操作和共享各设备资源。\[3]\[4]\[5]系统架构支持多内核，包括 Linux 内核、LiteOS 和鸿蒙微内核，可按各种智能设备选择所需内核，例如在低功耗设备上使用 LiteOS 内核。\[6]\[7]2019 年 8 月华为发布首款搭载鸿蒙操作系统的产品 “荣耀智能屏”，之后于 2021 年 6 月发布搭载鸿蒙操作系统的智能手机、平板电脑和智能手表。

Linux 内核是最大且变动最快的开源项目之一，它由大约 53,600 个文件和近 2,000 万行代码组成。在全世界范围内超过 15,600 位程序员为它贡献代码，Linux 内核项目的维护者使用了如下的协作模型。

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/1ef62d82-d29c-457a-b63f-3d4dfa09d652/640)

如果你有**深入 linux 内核的激情**和**极客精神**，可以为 Linux 项目贡献源码，具体如何提交第一个补丁，可以详细阅读下面文章，这里篇幅有限不展开：

<https://opensource.com/article/18/8/first-linux-kernel-patch>

Linux 开源代码仓库：

<https://github.com/torvalds/linux>

提交给 kernel 的补丁，刚开始可能不需要高深的技术，比如这个补丁，可以   是简单的对于已有内容的格式或拼写错误的修正，比如这个来自 4 岁小朋友的补丁：

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/1ef62d82-d29c-457a-b63f-3d4dfa09d652/640)

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/1ef62d82-d29c-457a-b63f-3d4dfa09d652/640)

**Linux 内核学习资源**

**源码：**

<https://elixir.bootlin.com/linux/latest/source>

在线交叉索引看源码，包括 Linux 几乎所有版本；

<https://github.com/torvalds/linux>

内核 github 仓库，可以下载本地，编译，修改和开发。

**网站**

[**http://www.kernel.org**](http://www.kernel.org)

可以通过这个网站上下载内核的源代码和补丁、跟踪内核 bug 等

[**http://lwn.net**](http://lwn.net)

Linux 内核最新消息，提供给了定期的与内核开发相关的报道

[**https://www.wiki.kernel.org/**](https://www.wiki.kernel.org/)

各种子模块 wiki 列表

[**http://www.linuxdoc.org**](http://www.linuxdoc.org)

Linux Documentation Project(Linux 文档项目)，拥有大量称为 “HowTo”
的文档，其中一些是技术性的，并涉及到一些内核相关的主题。

[**http://www.kerneltravel.net/**](http://www.kerneltravel.net/)

国内 Linux 内核之旅开源社区

[**http://www.linux-mm.org**](http://www.linux-mm.org)
该页面面向 Linux 内存管理开发，其中包含大量有用的信息，并且还包含大量与内核相关的 Web 站点链接。

[**http://www.wowotech.net**](http://www.wowotech.net)

博客专注分享 linux 内核知识（偏嵌入式方向）, 很多文章都非常精华和透彻，值得内核学习者学习；

[**https://blog.csdn.net/gatieme**](https://blog.csdn.net/gatieme)

操作系统优质博客，可以学习 linux 调度相关内核知识；

[**https://blog.csdn.net/dog250**](https://blog.csdn.net/dog250)

dog250 的文章都比较深刻，属于 Linux 内核进阶，可能不太适合入门，建议入门后，再看这里文章，会让你醍醐灌顶。

[**https://www.kernel.org/doc**](https://www.kernel.org/doc)

内核文档

**书籍**

《深入理解 Linux 内核》

《深入 Linux 内核架构》

《Linux 内核设计与实现》

《Linux 内核源代码情景分析》

《深入理解 LINUX 网络内幕》

《深入理解 Linux 虚拟内存管理》

《Linux 设备驱动程序》

### Git 分布式版本控制系统

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/1ef62d82-d29c-457a-b63f-3d4dfa09d652/640)

2005 年，Linus 还创建了 **Git**，这是非常流行的分布式源代码控制系统。迅速将 Linux 内核源代码树从专有 Bitkeeper 迁移到新创建的开源 Git。

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/1ef62d82-d29c-457a-b63f-3d4dfa09d652/640)

**git 架构**

**Git** 是出于需要而创建的，不是因为发现源代码控制很有趣，而是因为其他多数源代码控制系统不好用，不能满足当时开发需求，并且 git 在 Linux 开发模型中确实运行得相当好，BitKeeper 变得站不住脚。

完美适应现代开源软件的开发模式，分布式版本管理：

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/1ef62d82-d29c-457a-b63f-3d4dfa09d652/640)

### Linux 内核名人堂

让我们膜拜一下对 Linux 内核做出核心贡献的大神们：

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/1ef62d82-d29c-457a-b63f-3d4dfa09d652/640)

**林纳斯 · 班奈狄克 · 托瓦兹**（1969 年 12 月 28 日－），生于芬兰赫尔辛基市，拥有美国国籍，Linux 内核的最早作者，随后发起了这个开源项目，担任 Linux 内核的首要架构师与项目协调者，是当今世界最著名的电脑程序员、黑客之一。他还发起了开源项目 Git，并为主要的开发者。

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/1ef62d82-d29c-457a-b63f-3d4dfa09d652/640)

**大卫 · 史提芬 · 米勒**（英语：David Stephen Miller，1974 年 11 月 26 日－），网络昵称为 DaveM，生于美国新泽西州新布朗斯维克，著名程式员与骇客，负责 Linux 核心网络功能以及 SPARC 平台的实作。他也参与其他开源软件的开发，是 GCC 督导委员会的成员之一。根据 2013 年 8 月的统计，米勒是 Linux 核心源代码第二大的贡献者，自 2005 年开始，已经提交过 4989 个 patch。

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/1ef62d82-d29c-457a-b63f-3d4dfa09d652/640)

**葛雷格 · 克罗 - 哈曼**（英语：Greg Kroah-Hartman，姓名缩写为 GKH）， Linux 核心开发者，目前为 Linux 核心中稳定分支（`-stable`）的维护者\[2]，他也是 staging 子系统\[2]、USB\[2]driver core、debugfs、kref、kobject、sysfs kernel 子系统\[2]、 TTY layer \[2]、linux-hotplug、Userspace I/O（与 Hans J. Koch 共同维护）等专案的维护者\[2]，也创立了 udev 专案。除此之外，他亦协助维护 Gentoo Linux 中上述程式及 kernel 的套件。

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/1ef62d82-d29c-457a-b63f-3d4dfa09d652/640)

**艾伦 · 考克斯**（英语：Alan Cox，1968 年 7 月 22 日－）是一名英国程序员，生于英格兰索利赫尔。他自 1991 年开始投入 Linux 内核的开发工作，在开发者社群中有很高的地位，是 Linux 开发工作中的关键人物之一。他负责维护 Linux 内核 2.2 版这个分支，在 2.4 版中也提供许多程式码，拥有自己的分支版本。他住在威尔斯斯旺西，他的妻子于 2015 年逝世\[1]\[2]\[3]。2020 年他再婚\[4]\[5]。他于 1991 年在斯旺西大学获得计算机科学理学学士学位，2005 年在那里获得工商管理硕士学位\[6]。

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/1ef62d82-d29c-457a-b63f-3d4dfa09d652/640)

**英格 · 蒙内**（匈牙利语：Ingo Molnár），匈牙利软件程序员与骇客，在 linux 内核上有许多贡献，也拥有自己的 linux 分支版本。对于操作系统的安全性与效能提升方面，他的声名卓著，在 linux 内核中，他于 Linux-2.6.0 版加入 O(1) 排程器，在 Linux-2.6.23 版中加入**完全公平调度器 CFS**（Completely Fair Scheduler）。

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/1ef62d82-d29c-457a-b63f-3d4dfa09d652/640)

**米格尔 · 德伊卡萨**（西班牙语：Miguel de Icaza ，1972 年 11 月 23 日－），生于墨西哥市，著名墨西哥籍自由软件开发者，为 GNOME 项目与 Mono 项目的发起人。但后来\[何时？]退出了 GNOME 项目。

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/1ef62d82-d29c-457a-b63f-3d4dfa09d652/640)

**罗伯特 · 马修 · 拉姆**（英语：Robert Matthew Love，1981 年 9 月 25 日－），生于美国佛罗里达州，为著名自由软件程式开发者、作家，现职为 google 软件工程师。现居于波士顿。他是 linux 核心的主要开发者之一，主要负责程式排程、先占式核心、虚拟内存子系统、核心事件层。他也加入了 GNOME 计划。目前他在 google，主要负责 Android 系统的开发。

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/1ef62d82-d29c-457a-b63f-3d4dfa09d652/640)

**安德鲁 · 基斯 · 保罗 · 莫顿**（英语：Andrew Keith Paul Morton，1959 年－），生于英国英格兰，澳洲软件工程师与著名骇客。他是 Linux 核心开发社群的领导者之一，现为 ext3 的共同维护者，负责区块装置的日志层（Journaling layer for block devices，JBD）。他也是 mm tree 的负责人。

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/1ef62d82-d29c-457a-b63f-3d4dfa09d652/640)

**埃里克 · 斯蒂芬 · 雷蒙**（英语：Eric Steven Raymond，1957 年 12 月 4 日－），是一名程序员，《大教堂与市集》的作者、《新黑客词典》（"Jargon File"）的维护人、著名黑客。作为《新黑客词典》的主要编撰人以及维护者，雷蒙很早就被认为是黑客文化的历史学家以及人类学家。但是在 1997 年以后，雷蒙被广泛公认为是开放源代码运动的主要领导者之一，并且是最为大众所知道（并最具争议性）的黑客。

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/1ef62d82-d29c-457a-b63f-3d4dfa09d652/640)

**西奥多 · 曹**（英语：Theodore Y. Ts'o，1968 年 1 月 23 日－），小名**泰德 · 曹**（Ted Tso），汉名**曹子德**\[1]，生于美国加利福尼亚州帕罗奥图，著名的自由软件工程师，专长于文件系统设计。他是 Linux 内核在北美最早的开发者，负责 ext2、ext3 与 ext4 文件系统的开发与维护工作。他也是 e2fsprogs 的开发者。为自由标准组织的创始者之一，也曾担任 Linux 基金会首席技术官。

由于互联网发达，当前不管是从个人爱好，还是工作原因，对内核贡献的国人越来越多：

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/1ef62d82-d29c-457a-b63f-3d4dfa09d652/640)

<http://www.remword.com/kps_result/all_whole_line_country.html>

### 最后

30 年的时间，Linux 从一个个人玩具变成现在庞然大物，估值超过 100 亿美元，Linux 还带来一股开源潮流，让开源软件百花齐放，对计算机发展和开源文化起到极大促进作用。

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/1ef62d82-d29c-457a-b63f-3d4dfa09d652/640)

Linux 庞大的生态与发展过程，Linus 伟大而富有创造力并不足以在一篇文章中尽述。

匆匆 30 年，Linux 已经不仅仅是改变了世界，而且已经成为了这个世界不可或缺的一部分感谢 Linus Torvalds，感谢为之致力的一切贡献者！

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/1ef62d82-d29c-457a-b63f-3d4dfa09d652/640)

最后，为了致敬 Linux，希望大家三连支持，让更多人可以看到！

### 参考和扩展

<http://www.atguigu.com/jsfx/5694.html>

<https://opensource.com/article/16/12/yearbook-9-lessons-25-years-linux-kernel-development>

<https://www.reddit.com/r/linux/comments/2pqqla/kernel_commit_4_year_old_girl_fixes_formatting_to/utm_source=amp&utm_medium=&utm_content=post_title>

<http://oss.org.cn/ossdocs/linux/kernel/a1/index.html>

<http://www.wowotech.net/linux_kenrel/11.html>

<https://www.wikiwand.com/zh/Linux>

<https://zh.wikipedia.org/wiki/Category:Linux%E6%A0%B8%E5%BF%83%E9%A7%AD%E5%AE%A2>

<http://www.chromium.org/chromium-os/chromiumos-design-docs/software-architecture>

- END -

---

**看完一键三连 \*\***在看 \***\*，**转发 \***\*，点赞 \*\***

**是对文章最大的赞赏，极客重生感谢你**![](https://notes-learning.oss-cn-beijing.aliyuncs.com/1ef62d82-d29c-457a-b63f-3d4dfa09d652/640)

推荐阅读

[图解 Linux 内核 TCP/IP 协议栈实现 | Linux 网络硬核系列](http://mp.weixin.qq.com/s?__biz=MzkyMTIzMTkzNA==&mid=2247510568&idx=1&sn=79f335aaab5c0a36c0a66c5bfb1619ae&chksm=c1845d79f6f3d46f81b6fd24335eb8994c9daf21b6846d80af2cad73d9f638c5dda48b02892c&scene=21#wechat_redirect)

[Linux 网络子系统](http://mp.weixin.qq.com/s?__biz=MzkyMTIzMTkzNA==&mid=2247532046&idx=2&sn=04ffe282ce1278297d124f0c382ba665&chksm=c184895ff6f300497eb2bcc63d352b6d6b374606399cb7dd5b5bb59a773e674a368f9f4c9169&scene=21#wechat_redirect)

[开源, yyds!](http://mp.weixin.qq.com/s?__biz=MzkyMTIzMTkzNA==&mid=2247530537&idx=1&sn=11ed00203af160568e114e093dc706b4&chksm=c1848f78f6f3066ee21c11b603a683d28cee63924ab8397b5702d10227cdb69a9aab66e1bb0f&scene=21#wechat_redirect)

[操作系统的起源 | 开源运动的兴起](http://mp.weixin.qq.com/s?__biz=MzkyMTIzMTkzNA==&mid=2247508731&idx=1&sn=b562efbcdf5183ea4db6c9bced62894c&chksm=c18455aaf6f3dcbcbd225fcc4a8176ffa4401db3d7c9a0da80f0d4a302ce7bee4307ee4e38ee&scene=21#wechat_redirect)

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/1ef62d82-d29c-457a-b63f-3d4dfa09d652/640)
