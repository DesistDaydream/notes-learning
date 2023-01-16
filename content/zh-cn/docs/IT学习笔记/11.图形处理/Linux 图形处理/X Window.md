---
title: X Window
---

# 概述

> 参考：
> - [Wiki,X Window 系统协议和架构](https://en.wikipedia.org/wiki/X_Window_System_protocols_and_architecture)
> - [Wiki,X Window 管理器](https://en.wikipedia.org/wiki/X_window_manager)
> - [鸟哥的 Linux 私房菜,第二十三章、X Window](http://linux.vbird.org/linux_basic/0590xwindow.php)

Unix Like 操作系统不是只能进行服务器的架设而已，在美编、排版、制图、多媒体应用上也是有其需要的。 这些需求都需要用到图形介面 (Graphical User Interface, GUI) 的操作的， 所以后来才有所谓的 X Window System 这玩意儿。那么为啥图形窗口介面要称为 X 呢？因为就英文字母来看 X 是在 W(indow) 后面，因此，人们就戏称这一版的窗口介面为 X (有下一版的新窗口之意)！

事实上， X Window System 是个非常大的架构，他还用到网络功能呢！也就是说，其实 X 窗口系统是能够跨网络与跨操作系统平台的！

## X Window 的发展简史

X Window 系统最早是由 MIT (Massachusetts Institute of Technology, 麻省理工学院) 在 1984 年发展出来的， 当初 X 就是在 Unix 的 System V 这个操作系统版本上面开发出来的。在开发 X 时，开发者就希望这个窗口介面不要与硬件有强烈的相关性，这是因为如果与硬件的相关性高，那就等於是一个操作系统了， 如此一来的应用性会比较局限。因此 X 在当初就是以应用程序的概念来开发的，而非以操作系统来开发。

由於这个 X 希望能够透过网络进行图形介面的存取，因此发展出许多的 X 通讯协议，这些网络架构非常的有趣， 所以吸引了很多厂商加入研发，因此 X 的功能一直持续在加强！一直到 1987 年更改 X 版本到 X11 ，这一版 X 取得了明显的进步， 后来的窗口介面改良都是架构於此一版本，因此后来 X 窗口也被称为 X11 。这个版本持续在进步当中，到了 1994 年发布了新版的 X11R6 ，后来的架构都是沿用此一释出版本，所以后来的版本定义就变成了类似 1995 年的 X11R6.3 之类的样式。 (注 1)

1992 年 XFree86 (http://www.xfree86.org/) 计画顺利展开， 该计画持续在维护 X11R6 的功能性，包括对新硬件的支持以及更多新增的功能等等。当初定名为 XFree86 其实是根据『 X + Free software + x86 硬件 』而来的呢。早期 Linux 所使用的 X Window 的主要核心都是由 XFree86 这个计画所提供的，因此，我们常常将 X 系统与 XFree86 挂上等号的说。

不过由於一些授权的问题导致 XFree86 无法继续提供类似 GPL 的自由软件，后来 Xorg 基金会就接手 X11R6 的维护！ Xorg (http://www.x.org/) 利用当初 MIT 发布的类似自由软件的授权， 将 X11R6 拿来进行维护，并且在 2004 年发布了 X11R6.8 版本，更在 2005 年后发表了 X11R7.x 版。 现在我们 CentOS 5.x 使用的 X 就是 Xorg 提供的 X11R7 喔！ 而这个 X11R6/X11R7 的版本是自由软件，因此很多组织都利用这个架构去设计他们的图形介面喔！包括 Mac OS X v10.3 也曾利用过这个架构来设计他们的窗口呢！我们的 CentOS 也是利用 Xorg 提供的 X11 啦！

从上面的说明，我们可以知道的是：

- 在 Unix Like 上面的图形使用者界面 (GUI) 被称为 X 或 X11；
- X11 是一个『软件』而不是一个操作系统；
- X11 是利用网络架构来进行图形介面的运行与绘制；
- 较著名的 X 版本为 X11R6 这一版，目前大部分的 X 都是这一版演化出来的 (包括 X11R7)；
- 现在大部分的 distribution 使用的 X 都是由 Xorg 基金会所提供的 X11 软件；
- X11 使用的是 MIT 授权，为类似 GPL 的自由软件授权方式。

# X 原理

X Window system 是个利用网络架构的图形使用者接口软件，基本上是分成 X Server 与 X Client 两个元件而已喔！

1. X Server 管理硬件
2. X Client 则是应用程序。

在运行上，X Client 应用程序会将所想要呈现的画面告知 X Server ，最终由 X server 来将结果通过他所管理的硬件绘制出来！整体的架构我们大约可以使用如下的图示来作个介绍：\[2]
![](https://notes-learning.oss-cn-beijing.aliyuncs.com/qpbz9c/1616164721534-04137e6f-04ce-41e4-8b0a-2a4765c43e41.png)
上面的图示非常有趣喔！我们在用户端想要取得来自服务器的图形数据时，我们用户端使用的当然是用户端的硬件设备啊， 所以，X Server 的重点就是在管理用户端的硬件，包括接受键盘/鼠标等设备的输入信息， 并且将图形绘制到屏幕上（请注意上图的所有元件之间的箭头指示）。但是到底要绘制个啥东西呢？绘图总是需要一些数据才能绘制吧？此时 X Client(就是 X 应用程序) 就很重要啦！他主要提供的就是告知 X Server 要绘制啥东西。那照这样的想法来思考，我们是想要取得远端服务器的绘图数据来我们的计算机上面显示嘛！所以啰，远端服务器提供的是 X client 软件啊！

注意：在 X Window 环境中的 C/S 架构与通常意义上的 C/S 架构不太一样，有点反着来的感觉。

## X Server # 硬件管理、屏幕绘制与提供字体功能

既然 X Window System 是要显示图形接口，因此理所当然的需要一个元件来管理我主机上面的所有硬件设备才行！这个任务就是 X Server 所负责的。而我们在 X 发展简史当中提到的 XFree86 计划及 Xorg 基金会，主要提供的就是这个 X Server 啦！那么 X Server 管理的设备主要有哪些呢？其实与输入/输出有关喔！包括键盘、鼠标、手写板、显示器（monitor） 、屏幕分辨率与色彩深度、显卡（包含驱动程序） 与显示的字体等等，都是 X Server 管理的。

咦！显卡、屏幕以及键盘鼠标的设置，不是在开机的时候 Linux 系统以 systemd 的相关设置处理好了吗？为何 X Server 还要重新设置啊？这是因为 X Window 在 Linux 里面仅能算是“一套很棒的软件”， 所以 X Window 有自己的配置文件，你必须要针对他的配置文件设置妥当才行。也就是说， Linux 的设置与 X Server 的设置不一定要相同的！因此，你在 CentOS 7 的 multi-user.target 想要玩图形接口时，就得要载入 X Window 需要的驱动程序才行～总之， X Server 的主要功能就是在管理“主机”上面的显示硬件与驱动程序。

既然 X Window System 是以通过网络取得图形接口的一个架构，那么用户端是如何取得服务器端提供的图形画面呢？由于服务器与用户端的硬件不可能完全相同，因此我们用户端当然不可能使用到服务器端的硬件显示功能！举例来说，你的用户端计算机并没有 3D 影像加速功能，那么你的画面可能呈现出服务器端提供的 3D 加速吗？当然不可能吧！所以啰 X Server 的目的在管理用户端的硬件设备！也就是说：“每部用户端主机都需要安装 X Server，而服务器端则是提供 X Client 软件， 以提供用户端绘图所需要的数据数据”。

X Server/X Client 的互动并非仅有 client --> server，两者其实有互动的！从上图 23.1.1 我们也可以发现， X Server 还有一个重要的工作，那就是将来自输入设备（如键盘、鼠标等） 的动作告知 X Client， 你晓得， X Server 既然是管理这些周边硬件，所以，周边硬件的动作当然是由 X Server 来管理的， 但是 X Server 本身并不知道周边设备这些动作会造成什么显示上的效果， 因此 X Server 会将周边设备的这些动作行为告知 X Client ，让 X Client 去伤脑筋。

## X Client # 负责 X Server 要求的“事件”之处理

前面提到的 X Server 主要是管理显示接口与在屏幕上绘图，同时将输入设备的行为告知 X Client， 此时 X Client 就会依据这个输入设备的行为来开始处理，最后 X Client 会得到“ 嗯！这个输入设备的行为会产生某个图示”，然后将这个图示的显示数据回传给 X Server ， X server 再根据 X Client 传来的绘图数据将他描图在自己的屏幕上，来得到显示的结果。

也就是说， X Client 最重要的工作就是处理来自 X Server 的动作，将该动作处理成为绘图数据， 再将这些绘图数据传回给 X Server 啰！由于 X Client 的目的在产生绘图的数据，因此我们也称呼 X Client 为 X Application(X 应用程序)。而且，每个 X Client 并不知道其他 X Client 的存在， 意思是说，如果有两个以上的 X client 同时存在时，两者并不知道对方到底传了什么数据给 X Server ， 因此 X Client 的绘图常常会互相重叠而产生困扰喔！

举个例子来说，当我们在 X Window 的画面中，将鼠标向右移动，那他是怎么告知 X Server 与 X Client 的呢？首先， X server 会侦测到鼠标的移动，但是他不知道应该怎么绘图啊！此时，他将鼠标的这个动作告知 X Client， X Client 就会去运算，结果得到，嘿嘿！其实要将鼠标指标向右移动几个像素，然后将这个结果告知 X server ， 接下来，您就会看到 X Server 将鼠标指标向右移动啰～

这样做有什么好处啊？最大的好处是， X Client 不需要知道 X Server 的硬件配备与操作系统！因为 X Client 单纯就是在处理绘图的数据而已，本身是不绘图的。所以，在用户端的 X Server 用的是什么硬件？用的是哪套操作系统？服务器端的 X Client 根本不需要知道～相当的先进与优秀～对吧！^\_^ 整个运行流程可以参考下图：用户端用的是什么操作系统在 Linux 主机端是不在乎的！
![](https://notes-learning.oss-cn-beijing.aliyuncs.com/qpbz9c/1622302517121-61886ebd-b390-479f-a848-0558c2909ded.gif)
图 23.1.2、X Server 用户端的操作系统与 X client 的沟通示意

## X Window Manager # 特殊的 X Client ，负责管理所有的 X client 软件

刚刚前面提到，X Client 的主要工作是将来自 X Server 的数据处理成为绘图数据，再回传给 X Server 而已， 所以 X Client 本身是不知道他在 X Server 当中的位置、大小以及其他相关信息的。这也是上面我们谈到的， X Client 彼此不知道对方在屏幕的哪个位置啊！为了克服这个问题，因此就有 Window Manager(窗口管理器) 的产生了。窗口管理器也是 X client ，只是他主要在负责全部 X client 的控管，还包括提供某些特殊的功能，例如：

- 提供许多的控制元素，包括工作列、背景桌面的设置等等；
- 管理虚拟桌面（virtual desktop）；
- 提供窗口控制参数，这包括窗口的大小、窗口的重叠显示、窗口的移动、窗口的最小化等等。

当 X Window Manager 开始运作时，[X-Server](https://zh.wikipedia.org/wiki/X_Window%E7%B3%BB%E7%B5%B1%E7%9A%84%E5%8D%94%E8%AD%B0%E5%92%8C%E6%9E%B6%E6%A7%8B) 和 X-Client 之间的交互，会重定向 X Window Manager。每当要显示一个新视窗时，这个请求便会被重定向到 X Window Manager，它会决定视窗的初始位置。此外，大部分较新的 X Window Manager 会改变视窗的亲属关系，通常会在视窗顶部加上标题栏，并在视窗周围加上装饰性的框架。这两个部分皆由视窗管理器来控制，而不是其它程序。因此，当用户点击或拖曳那些组件时，X Window Manager 会进行适当的动作（如移动或改变视窗的大小）。

X Window Manager 也负责处理[图标](https://zh.wikipedia.org/wiki/%E5%9C%96%E7%A4%BA)，图标并不存在于 [X Window 核心协议](https://zh.wikipedia.org/wiki/X_Window%E6%A0%B8%E5%BF%83%E5%8D%94%E8%AD%B0)的层次中。当用户将视窗最小化时，X Window Manager 会取消视窗的映射（使其不可见），并完成适当的动作，将视窗改显示成图标。某些 X Window Manager 并不支持图标功能。

X Window Manager 主要的目标，就如同其名，是用来管理视窗的。许多 X Window Manager 提供附加的功能，如处理鼠标在[根视窗](https://zh.wikipedia.org/w/index.php?title=%E6%A0%B9%E8%A6%96%E7%AA%97&action=edit&redlink=1)上的点击，呈现出窗格以及其它的可视化组件，处理按键（例如 Alt-F4 可关闭视窗），判定哪一个应用程序在引导时运行等等。

我们常常听到的 KDE, GNOME, XFCE 还有阳春到爆的 twm 等等，都是一些窗口管理员的专案计划啦！这些专案计划中，每种窗口管理员所用以开发的显示发动机都不太相同，所著重的方向也不一样， 因此我们才会说，在 Linux 下面，每套 Window Manager 都是独特存在的，不是换了桌面与显示效果而已， 而是连显示的发动机都不会一样喔！下面是这些常见的 Window Manager 全名与链接：

- **GNU Network Object Model Environment(简称 GNOME)**：<http://www.gnome.org/>
- **K Desktop Enviroment(简称 KDE)**：<http://kde.org/>
- **Tab Window Manager(简称 TWM)**：<http://xwinman.org/vtwm.php>
- **XForms Common Environment(简称 XFCE)**：<http://www.xfce.org/>

由于 Linux 越来越朝向 Desktop 桌面电脑使用方向走，因此窗口管理员的角色会越来越重要！目前我们 CentOS 默认提供的有 GNOME 与 KDE ，这两个窗口管理员上面还有提供非常多的 X client 软件， 包括办公室生产力软件（Open Office） 以及常用的网络功能（firefox 浏览器、 Thunderbird 收发信件软件） 等。现在使用者想要接触 Linux 其实真的越来越简单了，如果不要架设服务器，那么 Linux 桌面的使用与 Windows 系统可以说是一模一样的！不需要学习也能够入门哩！^\_^

那么你知道 X Server / X client / window manager 的关系了吗？我们举 CentOS 默认的 GNOME 为例好了， 由于我们要在本机端启动 X Window system ，因此，在我们的 CentOS 主机上面必须要有 Xorg 的 X server 核心， 这样才能够提供屏幕的绘制啊～然后为了让窗口管理更方便，于是就加装了 GNOME 这个计划的 window manager ， 然后为了让自己的使用更方便，于是就在 GNOME 上面加上更多的窗口应用软件，包括输入法等等的， 最后就建构出我们的 X Window System 啰～ ^\_^！所以你也会知道，X server/X client/Window Manager 是同时存在于我们一部 Linux 主机上头的啦！

### Display Manager # 提供登陆需求

谈完了上述的数据后，我们得要了解一下，那么我如何取得 X Window 的控制？在本机的命令行下面你可以输入 startx 来启动 X 系统，此时由于你已经登陆系统了，因此不需要重新登陆即可取得 X 环境。但如果是 graphical.target 的环境呢？你会发现在 tty1 或其他 tty 的地方有个可以让你使用图形接口登陆（输入帐号密码） 的咚咚，那个是啥？是 X Server/X client 还是什么的？其实那是个 Display Manager 啦！这个 display manager 最大的任务就是提供登陆的环境， 并且载入使用者选择的 Window Manager 与语系等数据喔！

几乎所有的大型窗口管理员专案计划都会提供 display manager 的，在 CentOS 上面我们主要利用的是 GNOME 的 GNOME Display Manager(gdm) 这支程序来提供 tty1 的图形接口登陆喔！至于登陆后取得的窗口管理员，则可以在 gdm 上面进行选择的！我们在第四章介绍的登陆环境，那个环境其实就是 gdm 提供的啦！再回去参考看看图示吧！^\_^！所以说，并非 gdm 只能提供 GNOME 的登陆而已喔！

- GDM3 # GNOME 管理器
- KDM # KDE 管理器
- LightDM

## XDMCP # X 显示管理控制协议

当 X server, X client 都在同一部主机上面的时候，你可以很轻松的启动一个完整的 X Window System。但是如果你想要透过这个机制在网络上面启动 X 呢？此时你得先在客户端启动一个 X server 将图形接口绘图所需要的硬件装置配置好， 并且启动一个 X server 常见的接收埠口(通常是 port 6000)，然后再由服务器端的 X client 取得绘图数据，再将数据绘制成图啰。透过这个机制，你可以在任何一部启动 X server 登入服务器喔！而且不管你的操作系统是啥呢！意义就像下图， 如此一来，你就可以取得服务器所提供的图形接口环境啦！
![](https://notes-learning.oss-cn-beijing.aliyuncs.com/qpbz9c/1622307045187-e725761e-968e-455a-bdf4-797f37215985.png)
但是如果你是使用最笨的方法在客户端自己启动 X server ，然后在告诉服务器将 X client 程序一个一个的加载回来，那就太累人了吧！我们之前上面不是提到过可以用 display manager 来管理使用者的登入与启动 X 吗？那服务器能不能提供一个类似的服务，那我们直接透过服务器的 display manager 就能够提供我们登入的认证与加载自己选择的 window manager 的话，这样就太棒了！能够达到吗？当然可以啊！那就是透过 Xdmcp (X display manager control protocol) ([注 3](http://cn.linux.vbird.org/linux_server/0310telnetssh_3.php#ps3) )啦！

Xdmcp 启动后会在服务器的 udp 177 开始监听，然后当客户端的 X server 联机到服务器的 port 177 之后， 我们的 Xdmcp 就会在客户端的 X server 放上用户输入账密的图形接口程序啰！那你就能透过这个 Xdmcp 去加载服务器所提供的类似 Window Manager 的相关 X client 啰！那你就能够取得图形接口的远程联机服务器哩！赞吧！

那么什么时候会出现多使用者连入服务器取得 X 的情况呢？以鸟哥的例子来说，鸟哥实验室有一组 Linux 在进行数值模拟，他输出的结果是 NetCDF 档案，我们必须使用 PAVE 这一套软件去处理这些数据。但是我们有两三个人同时都会使用到那个功能，偏偏 Linux 主机是放在机架柜里面的，要我们挤在那个小小的空间前面『站着』操作计算机，可真是讨人厌啊～这个时候，我们就会架设图形接口的远程登录服务器，让我们可以『多人同时以图形接口登入 Linux 主机』来操作我们自己的程序！很棒，不是吗！

## X Server、X Client、X Window Manager 之间的关系

X Server、X Client、X Window Manager 之间的关系很像操作系统与应用程序之间的关系。很早之前，一台电脑只能运行一个程序，就好比只能运行 一个 X Client。所以，为了同时运行多个程序，就产生了操作系统。而为了同时打开多个图形界面，也就产生了 X Window Manager。

# X 的安装

## x Server 安装

yum install xorg-x11-xauth xorg-x11-server-utils

有了服务端之后，一些简单的 X 程序就可以通过 X client 打开了。

Node：xorg-x11-xinit 包含 xorg-x11-xauth xorg-x11-server-utils 这两个包，有时候安装 xorg-x11-xinit 即可

X server 安装完成后，可以安装 xorg-x11-apps 包，这里面有一些 x 程序，可以测试系统的图形接口。

Note: 在 centos8 下安装 apps 包，需要先开启 PowerTools 库(yum config-manager --set-enabled PowerTools)

- /usr/bin/luit
- /usr/bin/oclock
- /usr/bin/x11perf
- /usr/bin/x11perfcomp
- /usr/bin/xbiff
- /usr/bin/xclipboard
- /usr/bin/xclock
- /usr/bin/xconsole
- /usr/bin/xcursorgen
- /usr/bin/xcutsel
- /usr/bin/xdpr
- /usr/bin/xeyes
- /usr/bin/xfd
- /usr/bin/xfontsel
- /usr/bin/xload
- /usr/bin/xlogo
- /usr/bin/xmag
- /usr/bin/xmessage
- /usr/bin/xpr
- /usr/bin/xvidtune
- /usr/bin/xwd
- /usr/bin/xwud

## X Client 安装

1. xshell 工具在官方提供了一个 Xmanager 的软件，该软件可以作为 X Client 所用。

##
