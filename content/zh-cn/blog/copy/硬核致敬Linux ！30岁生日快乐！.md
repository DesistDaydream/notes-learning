---
title: "硬核致敬Linux ！30岁生日快乐！"
linkTitle: "硬核致敬Linux ！30岁生日快乐！"
date: 2021-08-26
---

原文链接：<https://mp.weixin.qq.com/s/cE4x63tYxoqrDinifeWqeg>

![](https://mmbiz.qpic.cn/mmbiz_png/Pn4Sm0RsAuhH4SOtTAkhF5RQnT4PAWdG2NT9Smu9eqEV5PAwKq3PbC6iagpqsfWz47RFLIZibibcIDAn3IyVS0ahw/640?wx_fmt=png)

1991年8月25日，21岁的Linus Torvalds（以下简称Linus）做了一个免费的操作系统“Linux”，并在这一天向外界公布这个由“业余爱好”主导的个人项目；如今，全球超级计算机500强和超过70%的智能手机都在运行Linux，因此，8月25日也被许多Linux的爱好者视为Linux真正的诞生日期。
# 你好
![](https://mmbiz.qpic.cn/mmbiz_png/cYSwmJQric6nhH4RQfgaJfjrfmLsALibgH5rcBjAWV0lF8QBtlXSJgRrJqBP90P2rfTd8WpVRAtyzqxhbXd6QnNg/640?wx_fmt=png)

30 年前，Linus Torvalds 首次发布 Linux 内核时还是赫尔辛基大学的一名 21 岁学生。他的声明是这样开始的，“我正在做一个（免费的）操作系统（只是一个爱好，不会很大和专业......）”。三十年后，排名前 500 的超级计算机都在运行 Linux，所有智能手机的 70% 以上都是如此。Linux 显然既庞大又专业。

三十年来，Linus Torvalds 领导了 Linux 内核开发，激励了无数其他开发人员和开源项目。2005 年，Linus 还创建了 Git来帮助管理内核开发过程，此后它成为最受欢迎的版本控制系统，受到无数开源和专有项目的信赖。

### Linux历史

**OS史前历史**

![](https://mmbiz.qpic.cn/mmbiz_png/cYSwmJQric6nhH4RQfgaJfjrfmLsALibgHE70XibBtS3DT8Nf3r5k48PGFo8ON6CPEsuyBOxIia8eLIQOuuz6JV1aA/640?wx_fmt=png)

**Linux的历史**  

![](https://mmbiz.qpic.cn/mmbiz_png/cYSwmJQric6nhH4RQfgaJfjrfmLsALibgHRicpx9aScqDmIdIib1M2UdibnVGVHJoTR5j94qiaCosHsT4G1XlPL1vYzA/640?wx_fmt=png)

### Linux系统

**Linux系统软件架构**

![](https://mmbiz.qpic.cn/mmbiz_jpg/cYSwmJQric6nhH4RQfgaJfjrfmLsALibgH1QmL07tiarw5K00x2LuwJEaRFCR4eev2O1DI1uEd6rZPOKnNZHzv5Eg/640?wx_fmt=jpeg)

Linux系统由硬件、kernel、系统调用、shell、c库、应用程序组成，架构层次分明，Linux内的各种层功能独立，程序在用户空间和内核空间之间的分离，能支持更多应用。  

| 用户模态 | **用户应用** | 例如：Bash，LibreOffice，GIMP，Blender，0 A.D.，Mozilla Firefox等 |
| 低层系统构件 | **系统守护进程**：  
systemd，runit，logind，networkd，PulseAudio等 | **窗口系统**：  
X11，Wayland，SurfaceFlinger(Android) | **其他库**：  
GTK+, Qt, EFL, SDL, SFML, FLTK, GNUstep等 | **图形**：  
Mesa，AMD Catalyst等 |
| **C标准库** | open()，exec()，sbrk()，socket()，fopen()，calloc()，... (直到2000个子例程)  
glibc目标为POSIX/SUS兼容，musl和uClibc目标为嵌入式系统，bionic为Android而写等 |
| 内核模态 | **Linux内核** | stat, splice, dup, read, open, ioctl, write, mmap, close, exit等（大约380个系统调用）  
Linux内核系统调用接口（SCI，目标为POSIX/SUS兼容） |
| 进程调度子系统 | IPC子系统 | 内存管理子系统 | 虚拟文件子系统 | 网络子系统 |
| 其他构件：ALSA，DRI，evdev，LVM，device mapper，Linux Network Scheduler，Netfilter  
Linux安全模块：SELinux，TOMOYO，AppArmor, Smack |
| 硬件（CPU，内存，数据存储设备等。） |

**Linux内核代码架构**  

**![](https://mmbiz.qpic.cn/mmbiz_jpg/cYSwmJQric6nhH4RQfgaJfjrfmLsALibgHSXWCabp9jDHopYuOGHSSq3HgXcQKLkWedlzYiaNOBStEEod9YkB8JAw/640?wx_fmt=jpeg)**

Linux代码包含核心几个子系统，比如**内存子系统**，**I/O子系统**，**CPU(调度）子系统**，**设备驱动子系统**，**网络子系统**，**虚拟文件子系统**等。这里简单介绍一些比较重要的子系统。

### 调度子系统

**进程调度**是Linux内核中最重要的子系统，它主要提供对CPU的访问控制。因为在计算机中，CPU资源是有限的，而众多的应用程序都要使用CPU资源，所以需要“进程调度子系统”对CPU进行调度管理。

**进程调度子系统**包括4个子模块（见下图），它们的功能如下：

![](https://mmbiz.qpic.cn/mmbiz_jpg/cYSwmJQric6nhH4RQfgaJfjrfmLsALibgHzhTgCrceib1hkg2RZuwSgf86iaN2JfyHMptZFZGdw0fhQW1hNLgMBlHA/640?wx_fmt=jpeg)

1.  Scheduling Policy，实现进程调度的策略，它决定哪个（或哪几个）进程将拥有CPU。
    
2.  Architecture-specific Schedulers，体系结构相关的部分，用于将对不同CPU的控制，抽象为统一的接口。这些控制主要在suspend和resume进程时使用，牵涉到CPU的寄存器访问、汇编指令操作等。
    
3.  Architecture-independent Scheduler，体系结构无关的部分。它会和“Scheduling Policy模块”沟通，决定接下来要执行哪个进程，然后通过“Architecture-specific Schedulers模块”resume指定的进程。
    
4.  System Call Interface，系统调用接口。进程调度子系统通过系统调用接口，将需要提供给用户空间的接口开放出去，同时屏蔽掉不需要用户空间程序关心的细节。
    

### 内存子系统

**内存管理**同样是Linux内核中最重要的子系统，它主要提供对内存资源的访问控制。Linux系统会在硬件物理内存和进程所使用的内存（称作虚拟内存）之间建立一种映射关系，这种映射是以进程为单位，因而不同的进程可以使用相同的虚拟内存，而这些相同的虚拟内存，可以映射到不同的物理内存上。

**内存管理子系统**包括3个子模块（见下图），它们的功能如下：

![](https://mmbiz.qpic.cn/mmbiz_jpg/cYSwmJQric6nhH4RQfgaJfjrfmLsALibgHNY1OZLaKPcddjClNYLNRmEPuUt7fQf7iafZ7eJtrP46fWqws0wD5yww/640?wx_fmt=jpeg)

1.  Architecture Specific Managers，体系结构相关部分。提供用于访问硬件Memory的虚拟接口。
    
2.  Architecture Independent Manager，体系结构无关部分。提供所有的内存管理机制，包括：以进程为单位的memory mapping；虚拟内存的Swapping。
    
3.  System Call Interface，系统调用接口。通过该接口，向用户空间程序应用程序提供内存的分配、释放，文件的map等功能。
    

### 虚拟文件子系统（Virtual Filesystem, VFS）

传统意义上的文件系统，是一种存储和组织计算机数据的方法。它用易懂、人性化的方法（文件和目录结构），抽象计算机磁盘、硬盘等设备上冰冷的数据块，从而使对它们的查找和访问变得容易。因而文件系统的实质，就是“存储和组织数据的方法”，文件系统的表现形式，就是“从某个设备中读取数据和向某个设备写入数据”。

随着计算机技术的进步，存储和组织数据的方法也是在不断进步的，从而导致有多种类型的文件系统，例如FAT、FAT32、NTFS、EXT2、EXT3等等。而为了兼容，操作系统或者内核，要以相同的表现形式，同时支持多种类型的文件系统，这就延伸出了**虚拟文件系统（VFS）**的概念。VFS的功能就是管理各种各样的文件系统，屏蔽它们的差异，以统一的方式，为用户程序提供访问文件的接口。

我们可以从磁盘、硬盘、NAND Flash等设备中读取或写入数据，因而最初的文件系统都是构建在这些设备之上的。这个概念也可以推广到其它的硬件设备，例如内存、显示器（LCD）、键盘、串口等等。我们对硬件设备的访问控制，也可以归纳为读取或者写入数据，因而可以用统一的文件操作接口访问。Linux内核就是这样做的，除了传统的磁盘文件系统之外，它还抽象出了设备文件系统、内存文件系统等等。这些逻辑，都是由VFS子系统实现。

**VFS子系统**包括6个子模块（见下图），它们的功能如下：

![](https://mmbiz.qpic.cn/mmbiz_jpg/cYSwmJQric6nhH4RQfgaJfjrfmLsALibgHxGOz2c8yZ6yxrrHeibxu3ac8nibP9KGUkVkrXPmgIH9iasFgLjaofdWqA/640?wx_fmt=jpeg)

1.  Device Drivers，设备驱动，用于控制所有的外部设备及控制器。由于存在大量不能相互兼容的硬件设备（特别是嵌入式产品），所以也有非常多的设备驱动。因此，Linux内核中将近一半的Source Code都是设备驱动，大多数的Linux底层工程师（特别是国内的企业）都是在编写或者维护设备驱动，而无暇估计其它内容（它们恰恰是Linux内核的精髓所在）。
    
2.  Device Independent Interface， 该模块定义了描述硬件设备的统一方式（统一设备模型），所有的设备驱动都遵守这个定义，可以降低开发的难度。同时可以用一致的形势向上提供接口。
    
3.  Logical Systems，每一种文件系统，都会对应一个Logical System（逻辑文件系统），它会实现具体的文件系统逻辑。
    
4.  System Independent Interface，该模块负责以统一的接口（快设备和字符设备）表示硬件设备和逻辑文件系统，这样上层软件就不再关心具体的硬件形态了。
    
5.  System Call Interface，系统调用接口，向用户空间提供访问文件系统和硬件设备的统一的接口。
    

### 网络子系统（Net）

**网络子系统**在Linux内核中主要负责管理各种网络设备，并实现各种网络协议栈，最终实现通过网络连接其它系统的功能。在Linux内核中，网络子系统几乎是自成体系，它包括5个子模块（见下图），它们的功能如下：

![](https://mmbiz.qpic.cn/mmbiz_jpg/cYSwmJQric6nhH4RQfgaJfjrfmLsALibgHy9IqgNnoDPLwVHDKqVBzD6iaJtzNk7wm3h9aSn1Nf6xOsNM870ScwEA/640?wx_fmt=jpeg)

1.  Network Device Drivers，网络设备的驱动，和VFS子系统中的设备驱动是一样的。
    
2.  Device Independent Interface，和VFS子系统中的是一样的。
    
3.  Network Protocols，实现各种网络传输协议，例如IP, TCP, UDP等等。
    
4.  Protocol Independent Interface，屏蔽不同的硬件设备和网络协议，以相同的格式提供接口（socket)。
    
5.  System Call interface，系统调用接口，向用户空间提供访问网络设备的统一的接口。
    

Linux内核版本时间线：

![](https://mmbiz.qpic.cn/mmbiz_png/cYSwmJQric6nhH4RQfgaJfjrfmLsALibgHRqUJPPCSOutrOnet4ICAkEbUHf9LLgGX4rorHj7nXDOcyTvLEJ29Bg/640?wx_fmt=png)

**Linux内核支持各种硬件架构**

**Linux内核**最成功的地方之一就是支持各种硬件架构，为软件提供了公共的平台：

![](https://mmbiz.qpic.cn/mmbiz_png/cYSwmJQric6nhH4RQfgaJfjrfmLsALibgHlp75WAFV1ibBJ0IAN2BVxzKia8ET9rFicziaQxtA0sEl0icYnHula6HzsIQ/640?wx_fmt=png)

基于Linux的系统是一个模块化的类Unix操作系统。**Linux操作系统**的大部分设计思想来源于20世纪70年代到80年代的Unix操作系统所建立的基本设计思想。Linux系统使用宏内核，由Linux内核负责处理进程控制、网络，以及外围设备和文件系统的访问。在系统运行的时候，设备驱动程序要么与内核直接整合，要么以加载模块形式添加。

**Linux具有设备独立性**，它内核具有高度适应能力，从而给系统提供了更高级的功能。GNU用户界面组件是大多数Linux操作系统的重要组成部分，提供常用的C函数库，Shell，还有许多常见的Unix实用工具，可以完成许多基本的操作系统任务。大多数Linux系统使用的图形用户界面建立在X窗口系统之上，由X窗口(XWindow)系统通过软件工具及架构协议来建立操作系统所用的图形用户界面.

**基于Linux内核各种衍生OS系统**  

![](https://mmbiz.qpic.cn/mmbiz_png/cYSwmJQric6nhH4RQfgaJfjrfmLsALibgHiaA5ia2r6z3iarNHHZN7gHOISWONZ0vfBibAz5wm0tJfxaBao2KOF4fM8Q/640?wx_fmt=png)

各种发行版本

![](https://mmbiz.qpic.cn/mmbiz_png/cYSwmJQric6nhH4RQfgaJfjrfmLsALibgHJQp1UUv39Vx04CB4W3DSdh2twrDS7kZV01ygjFXbyeUCfg92kJpiaqA/640?wx_fmt=png)

当前比较流行发行版是：**Debian**、**Ubuntu**、**Fedora**、**CentOS**、**Arch Linux**和**openSUSE**等，每个发行版都有自己优势地方，都有一批忠实用户。  

**基于Linux内核著名OS**  

**Android**

![](https://mmbiz.qpic.cn/mmbiz_png/cYSwmJQric6nhH4RQfgaJfjrfmLsALibgHwsGP7L9tNBESaK6GlTHaZWjXcial9ia9bWDQoXqPDUSh4rxfYtCOG1XA/640?wx_fmt=png)

**Android**（读音：英：\['ændrɔɪd\]，美：\[ˈænˌdrɔɪd\]），中文用户多以非官方名称“安卓”称之，是一个基于Linux内核与其他开源软件的开放源代码的移动操作系统，Android的内核是根据Linux内核的长期支持的分支，具有典型的Linux调度和功能。截至2018年，Android的目标是Linux内核的4.4、4.9或是4.14版本。

**ChromeOS**

![](https://mmbiz.qpic.cn/mmbiz_png/cYSwmJQric6nhH4RQfgaJfjrfmLsALibgHGtbgyZUUeFSicfcZCeRpTGicmJI5vzDibiaCibBfoA06demaPcn1iaVficviag/640?wx_fmt=png)

**Chrome OS** 是由Google设计基于Linux内核的操作系统，并使用Google Chrome浏览器作为其主要用户界面。因此，Chrome OS主要支持Web应用程序\[6\]，2016年起开始陆续兼容Android应用程序（可通过Google Play商店下载）和Linux应用程序。  

**鸿蒙OS**  

![](https://mmbiz.qpic.cn/mmbiz_png/cYSwmJQric6nhH4RQfgaJfjrfmLsALibgHdmDMzZHEZR94PaMFm9CYibMfvicxDs2ULz9jQJVZKUQNOtAUduCqRJUA/640?wx_fmt=png)

**鸿蒙**（**HarmonyOS**，开发代号Ark\[1\]）是华为自2012年开发的一款可兼容Android应用程序的跨平台分布式操作系统\[2\]。系统性能包括利用“分布式”技术将各款设备融合成一个“超级终端”，便于操作和共享各设备资源。\[3\]\[4\]\[5\]系统架构支持多内核，包括Linux内核、LiteOS和鸿蒙微内核，可按各种智能设备选择所需内核，例如在低功耗设备上使用LiteOS内核。\[6\]\[7\]2019年8月华为发布首款搭载鸿蒙操作系统的产品“荣耀智能屏”，之后于2021年6月发布搭载鸿蒙操作系统的智能手机、平板电脑和智能手表。  

Linux 内核是最大且变动最快的开源项目之一，它由大约 53,600 个文件和近 2,000 万行代码组成。在全世界范围内超过 15,600 位程序员为它贡献代码，Linux 内核项目的维护者使用了如下的协作模型。

![](https://mmbiz.qpic.cn/mmbiz_jpg/cYSwmJQric6nhH4RQfgaJfjrfmLsALibgHHIWib5fYsnBg0ziccyM2IEj3VRhGUQlkzPSyDpQ5cicLo148yoJ8z33LQ/640?wx_fmt=jpeg)

如果你有**深入linux内核的激情**和**极客精神**，可以为Linux项目贡献源码，具体如何提交第一个补丁，可以详细阅读下面文章，这里篇幅有限不展开：

https://opensource.com/article/18/8/first-linux-kernel-patch

Linux 开源代码仓库：  

https://github.com/torvalds/linux

提交给kernel的补丁，刚开始可能不需要高深的技术，比如这个补丁，可以 是简单的对于已有内容的格式或拼写错误的修正，比如这个来自4岁小朋友的补丁：

![](https://mmbiz.qpic.cn/mmbiz_png/cYSwmJQric6nhH4RQfgaJfjrfmLsALibgHv1787MroX0hXfDbFQCZhqCW68Otw5fpggEib5QicCewZGj0ZRWDVBduQ/640?wx_fmt=png)

![](https://mmbiz.qpic.cn/mmbiz_png/cYSwmJQric6nhH4RQfgaJfjrfmLsALibgHhCiab6dlIM2ibMY4KPSrup1iaHNp6ciayzJ0oBRcvyIibRLltBDSw1zyxMA/640?wx_fmt=png)

**Linux内核学习资源**

**源码：** 

https://elixir.bootlin.com/linux/latest/source  

在线交叉索引看源码，包括Linux几乎所有版本；

https://github.com/torvalds/linux

内核github仓库，可以下载本地，编译，修改和开发。  

**网站**

**http://www.kernel.org**

可以通过这个网站上下载内核的源代码和补丁、跟踪内核bug等

**http://lwn.net**

Linux 内核最新消息，提供给了定期的与内核开发相关的报道

**https://www.wiki.kernel.org/**

各种子模块wiki列表

**http://www.linuxdoc.org**

Linux Documentation Project(Linux文档项目)，拥有大量称为“HowTo”  
的文档，其中一些是技术性的，并涉及到一些内核相关的主题。

**http://www.kerneltravel.net/**

国内Linux内核之旅开源社区

**http://www.linux-mm.org**  
该页面面向Linux内存管理开发，其中包含大量有用的信息，并且还包含大量与内核相关的Web站点链接。

**http://www.wowotech.net**

博客专注分享linux内核知识（偏嵌入式方向）, 很多文章都非常精华和透彻，值得内核学习者学习；

**https://blog.csdn.net/gatieme**

操作系统优质博客，可以学习linux 调度相关内核知识；

**https://blog.csdn.net/dog250**

dog250的文章都比较深刻，属于Linux内核进阶，可能不太适合入门，建议入门后，再看这里文章，会让你醍醐灌顶。

**https://www.kernel.org/doc**

内核文档  

**书籍**

《深入理解Linux内核》

《深入Linux内核架构》

《Linux内核设计与实现》

《Linux内核源代码情景分析》

《深入理解LINUX网络内幕》

《深入理解Linux虚拟内存管理》

《Linux设备驱动程序》

### Git分布式版本控制系统

![](https://mmbiz.qpic.cn/mmbiz_png/cYSwmJQric6nhH4RQfgaJfjrfmLsALibgHplRZhR4pd07qoWDXz6icgRfqLyyC9SbZnxx2PIVCOyYHuWQjYDC3IDw/640?wx_fmt=png)

2005 年，Linus还创建了 **Git**，这是非常流行的分布式源代码控制系统。迅速将 Linux 内核源代码树从专有 Bitkeeper 迁移到新创建的开源 Git。

![](https://mmbiz.qpic.cn/mmbiz_png/cYSwmJQric6nhH4RQfgaJfjrfmLsALibgHK6ewlMwmzxJxkuJ2Wb4dUxdUmOuUcVybz2NpicpKYHpbW3Kf8Rg9iabA/640?wx_fmt=png)

**git 架构**

**Git** 是出于需要而创建的，不是因为发现源代码控制很有趣，而是因为其他多数源代码控制系统不好用，不能满足当时开发需求，并且 git 在 Linux 开发模型中确实运行得相当好，BitKeeper变得站不住脚。  

完美适应现代开源软件的开发模式，分布式版本管理：

![](https://mmbiz.qpic.cn/mmbiz_png/cYSwmJQric6nhH4RQfgaJfjrfmLsALibgHq1UTS3ZCgoIbcKBtU67SJcsLUS7osnicsG9LzkJM75hkjT1O9wVat9g/640?wx_fmt=png)

### Linux内核名人堂

让我们膜拜一下对Linux内核做出核心贡献的大神们：  

![](https://mmbiz.qpic.cn/mmbiz_jpg/cYSwmJQric6nhH4RQfgaJfjrfmLsALibgHtEVHKmWufP3VMY68yCQnh3fgdg40AFwfHhLmz2PaqLqXZJHCKQSn9Q/640?wx_fmt=jpeg)

**林纳斯·班奈狄克·托瓦兹**（1969年12月28日－），生于芬兰赫尔辛基市，拥有美国国籍，Linux内核的最早作者，随后发起了这个开源项目，担任Linux内核的首要架构师与项目协调者，是当今世界最著名的电脑程序员、黑客之一。他还发起了开源项目Git，并为主要的开发者。

![](https://mmbiz.qpic.cn/mmbiz_jpg/cYSwmJQric6nhH4RQfgaJfjrfmLsALibgHxticTE1khTNhUIXzRscgET2fazkOk0ISAzIrORD6X51wTUPSlfzCUxg/640?wx_fmt=jpeg)

**大卫·史提芬·米勒**（英语：David Stephen Miller，1974年11月26日－），网络昵称为 DaveM，生于美国新泽西州新布朗斯维克，著名程式员与骇客，负责Linux核心网络功能以及SPARC平台的实作。他也参与其他开源软件的开发，是GCC督导委员会的成员之一。根据2013年8月的统计，米勒是Linux核心源代码第二大的贡献者，自2005年开始，已经提交过4989个patch。

![](https://mmbiz.qpic.cn/mmbiz_jpg/cYSwmJQric6nhH4RQfgaJfjrfmLsALibgHC4qCA4jxF7rIrqwJpzFhzESIQx7ibZ5dIXkicz0XsBcwAibveORZ4NA5g/640?wx_fmt=jpeg)

**葛雷格·克罗-哈曼**（英语：Greg Kroah-Hartman，姓名缩写为GKH）， Linux核心开发者，目前为 Linux 核心中稳定分支（`-stable`）的维护者\[2\]，他也是staging 子系统\[2\]、USB\[2\]driver core、debugfs、kref、kobject、sysfs kernel 子系统\[2\]、 TTY layer \[2\]、linux-hotplug、Userspace I/O（与 Hans J. Koch 共同维护）等专案的维护者\[2\]，也创立了udev专案。除此之外，他亦协助维护Gentoo Linux中上述程式及 kernel 的套件。

![](https://mmbiz.qpic.cn/mmbiz_jpg/cYSwmJQric6nhH4RQfgaJfjrfmLsALibgHiaR4nLoIu5MJ2cAFdjykndwROoicmZXR7rpqqXyReD4RCiaPwOrsc2urQ/640?wx_fmt=jpeg)

**艾伦·考克斯**（英语：Alan Cox，1968年7月22日－）是一名英国程序员，生于英格兰索利赫尔。他自1991年开始投入Linux内核的开发工作，在开发者社群中有很高的地位，是Linux开发工作中的关键人物之一。他负责维护Linux内核 2.2版这个分支，在2.4版中也提供许多程式码，拥有自己的分支版本。他住在威尔斯斯旺西，他的妻子于2015年逝世\[1\]\[2\]\[3\]。2020年他再婚\[4\]\[5\]。他于1991年在斯旺西大学获得计算机科学理学学士学位，2005年在那里获得工商管理硕士学位\[6\]。

![](https://mmbiz.qpic.cn/mmbiz_jpg/cYSwmJQric6nhH4RQfgaJfjrfmLsALibgHuOQOEoMQhVfJibrcPjBxIy3Y61qGos8lFbz70cSf8Tycwy8OQsX5Nhg/640?wx_fmt=jpeg)

**英格·蒙内**（匈牙利语：Ingo Molnár），匈牙利软件程序员与骇客，在linux内核上有许多贡献，也拥有自己的linux分支版本。对于操作系统的安全性与效能提升方面，他的声名卓著，在linux内核中，他于Linux-2.6.0版加入O(1)排程器，在 Linux-2.6.23版中加入**完全公平调度器CFS**（Completely Fair Scheduler）。

![](https://mmbiz.qpic.cn/mmbiz_png/cYSwmJQric6nhH4RQfgaJfjrfmLsALibgHBWhcnGvYiaDtbmsqnbjoRFnyXSTs5ibMAZ7Tm0vJuU573vKyrpzvvIjw/640?wx_fmt=png)

**米格尔·德伊卡萨**（西班牙语：Miguel de Icaza ，1972年11月23日－），生于墨西哥市，著名墨西哥籍自由软件开发者，为GNOME项目与Mono项目的发起人。但后来\[何时？\]退出了GNOME项目。  

![](https://mmbiz.qpic.cn/mmbiz_jpg/cYSwmJQric6nhH4RQfgaJfjrfmLsALibgHdkfW1X7mmFsq1ShM2iaHl6XJD5RPK1xaOJ60hAxiczdEDzGWrw2HorhA/640?wx_fmt=jpeg)

**罗伯特·马修·拉姆**（英语：Robert Matthew Love，1981年9月25日－），生于美国佛罗里达州，为著名自由软件程式开发者、作家，现职为google软件工程师。现居于波士顿。他是linux核心的主要开发者之一，主要负责程式排程、先占式核心、虚拟内存子系统、核心事件层。他也加入了GNOME计划。目前他在google，主要负责Android系统的开发。

![](https://mmbiz.qpic.cn/mmbiz_jpg/cYSwmJQric6nhH4RQfgaJfjrfmLsALibgHn9h6QWB8CfkNDIeoLLNgW4p3VkHveT61gVHiccEtjfkF708FfXUW3Aw/640?wx_fmt=jpeg)

**安德鲁·基斯·保罗·莫顿**（英语：Andrew Keith Paul Morton，1959年－），生于英国英格兰，澳洲软件工程师与著名骇客。他是Linux核心开发社群的领导者之一，现为ext3的共同维护者，负责区块装置的日志层（Journaling layer for block devices，JBD）。他也是mm tree的负责人。

![](https://mmbiz.qpic.cn/mmbiz_jpg/cYSwmJQric6nhH4RQfgaJfjrfmLsALibgH2O1Ze7oATYXWrxfiaHkZqB5ggRna9RRft2huTliaQWyNrV061Z0q3icKQ/640?wx_fmt=jpeg)

**埃里克·斯蒂芬·雷蒙**（英语：Eric Steven Raymond，1957年12月4日－），是一名程序员，《大教堂与市集》的作者、《新黑客词典》（"Jargon File"）的维护人、著名黑客。作为《新黑客词典》的主要编撰人以及维护者，雷蒙很早就被认为是黑客文化的历史学家以及人类学家。但是在1997年以后，雷蒙被广泛公认为是开放源代码运动的主要领导者之一，并且是最为大众所知道（并最具争议性）的黑客。

![](https://mmbiz.qpic.cn/mmbiz_jpg/cYSwmJQric6nhH4RQfgaJfjrfmLsALibgHutx5g8HmjnzMFaeXZ7kc8TOttGwG1GUg5FswuJtelvTo0GWbZTOjPQ/640?wx_fmt=jpeg)

**西奥多·曹**（英语：Theodore Y. Ts'o，1968年1月23日－），小名**泰德·曹**（Ted Tso），汉名**曹子德**\[1\]，生于美国加利福尼亚州帕罗奥图，著名的自由软件工程师，专长于文件系统设计。他是Linux内核在北美最早的开发者，负责ext2、ext3与ext4文件系统的开发与维护工作。他也是e2fsprogs的开发者。为自由标准组织的创始者之一，也曾担任Linux基金会首席技术官。

由于互联网发达，当前不管是从个人爱好，还是工作原因，对内核贡献的国人越来越多：

![](https://mmbiz.qpic.cn/mmbiz_png/cYSwmJQric6nhH4RQfgaJfjrfmLsALibgHQ2BEn5hagVPEelr1qE2EHlGWclItckm0saPvuB5zACbZWzIB6kbNeQ/640?wx_fmt=png)

http://www.remword.com/kps\_result/all\_whole\_line\_country.html

### 最后

30年的时间，Linux从一个个人玩具变成现在庞然大物，估值超过100亿美元，Linux还带来一股开源潮流，让开源软件百花齐放，对计算机发展和开源文化起到极大促进作用。  

![](https://mmbiz.qpic.cn/mmbiz_jpg/cYSwmJQric6nhH4RQfgaJfjrfmLsALibgHibFAnWeJmNESo5lWBCcEYdI2MpVkiabQ6n4B3FRpv1FwfDpYaTDWXG4w/640?wx_fmt=jpeg)

Linux 庞大的生态与发展过程，Linus伟大而富有创造力并不足以在一篇文章中尽述。

匆匆30 年，Linux 已经不仅仅是改变了世界，而且已经成为了这个世界不可或缺的一部分感谢 Linus Torvalds，感谢为之致力的一切贡献者！

![](https://mmbiz.qpic.cn/mmbiz_jpg/cYSwmJQric6nhH4RQfgaJfjrfmLsALibgHxfSUWjnl9ia6I5GPetib9tqehO96tNWdaPEQzicHIjk2QQ1eNq5WyQTDw/640?wx_fmt=jpeg)

最后，为了致敬Linux，希望大家三连支持，让更多人可以看到！

### 参考和扩展

http://www.atguigu.com/jsfx/5694.html  

https://opensource.com/article/16/12/yearbook-9-lessons-25-years-linux-kernel-development

https://www.reddit.com/r/linux/comments/2pqqla/kernel\_commit\_4\_year\_old\_girl\_fixes\_formatting\_to/utm\_source=amp&utm\_medium=&utm\_content=post\_title

http://oss.org.cn/ossdocs/linux/kernel/a1/index.html

http://www.wowotech.net/linux_kenrel/11.html

https://www.wikiwand.com/zh/Linux

https://zh.wikipedia.org/wiki/Category:Linux%E6%A0%B8%E5%BF%83%E9%A7%AD%E5%AE%A2

http://www.chromium.org/chromium-os/chromiumos-design-docs/software-architecture

\- END -
