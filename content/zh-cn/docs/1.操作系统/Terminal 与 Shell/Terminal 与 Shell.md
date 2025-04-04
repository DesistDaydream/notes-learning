---
title: Terminal 与 Shell
linkTitle: Terminal 与 Shell
weight: 1
tags:
  - CLI
  - GUI
---

# 概述

> 参考：
>
> - [Wiki, Shell](<https://en.wikipedia.org/wiki/Shell_(computing)>)
> - [Manual(手册)，bash](https://www.man7.org/linux/man-pages/man1/bash.1.html)
> - <https://blog.csdn.net/kangkanglou/article/details/82698177>
> - <http://feihu.me/blog/2014/env-problem-when-ssh-executing-command-on-remote/>
> - <https://www.jianshu.com/p/0c7ea235b473>
> - [公众号，阿里技术-一种命令行解析的新思路](https://mp.weixin.qq.com/s/RxpcqBGhUT-5z4N5kRXvBg)

**Shell(也称为壳层)** 是一种计算机程序，向人类用户或其他程序公开操作系统的服务。通常，操作系统的 Shell 程序会根据计算机的角色和特定操作，分为两类：

- **[command-line interface](https://en.wikipedia.org/wiki/Command-line_interface)(命令行界面，简称 CLI)**
- **[graphical user interface](https://en.wikipedia.org/wiki/Graphical_user_interface)(图形用户界面，简称 GUI)**

Shell 在计算机科学中指“为用户提供用户界面”的软件，通常指的是命令行界面的解析器。一般来说，这个词是指操作系统中提供访问内核所提供之服务的程序。Shell 也用于泛指所有为用户提供操作界面的程序，也就是程序和用户交互的接口。因此与之相对的是 Kernel(内核)，内核不提供和用户的交互功能。

用白话讲：人类操作计算机的地方就是 Shell ，可以是一个输入命令的地方(CLI)，也可以是一个用户用鼠标点点点的图形界面(GUI)。Shell 也是一类程序的统称，实际上，当输入完用户名和密码登录系统的时候，就是 Linux 系统后台自动启动了一个名叫 Bash 的 Shell 程序，来让用户输入指令对计算机进行操作

所以，一个 Shell 也会有一个进程号，在该 Shell 执行的程序的父进程号都是该 Shell 的进程号

如下所示，[登录系统](/docs/1.操作系统/登录%20Linux%20与%20访问控制/登录%20Linux%20与%20访问控制.md)时，会启动一个进程标识当前登录的用户，并启动一个子进程，该子进程就是 Bash 这个 Shell，并且会为该 Shell 分配一个[终端](#Terminal(终端))来与用户进行交互(这里的终端名是 tty1)

```bash
root      1067     1  0 11:20 ?        Ss     0:00 login -- root
root      9622  1067  2 13:19 tty1     Ss+    0:00  \_ -bash
```

所有命令都是在这个 shell 下运行的，如下所使，在 bash 下执行了一个 cat 命令

```bash
root      1067     1  0 11:20 ?        Ss     0:00 login -- root
root      9622  1067  0 13:19 tty1     Ss     0:00  \_ -bash
root     11198  9622  0 13:22 tty1     S+     0:00      \_ cat
```

Note：

- 有一点需要明确，系统下的任何程序运行都需要一个用户，哪怕在刚装完系统第一次启动，所有进程也是基于 root 用户来运行的，所以脱离用户讨论 shell 是不对的。
- **由于 Linux 常用的 Shell 为 bash，下面主要描述的都是关于 bash shell 的配置**

# Terminal(终端)

> 参考：
>
> - [Wiki, Computer Terminal](https://en.wikipedia.org/wiki/Computer_terminal)
> - [Wiki, TTY](https://en.wikipedia.org/wiki/Teleprinter)
> - [Wiki, TTY(Unix)](https://en.wikipedia.org/wiki/Tty_(unix))
> - [Wiki, Psedoterminal](https://en.wikipedia.org/wiki/Pseudoterminal)
> - [Wiki, Terminal emulator](https://en.wikipedia.org/wiki/Terminal_emulator)
> - [TTY 代表什么](https://askubuntu.com/questions/481906/what-does-tty-stand-for)
> - [Manual(手册)，pty(7)](https://man7.org/linux/man-pages/man7/pty.7.html)
> - [Manual(手册)，pts(4)](https://man7.org/linux/man-pages/man4/pts.4.html)

**Computer Terminal(计算机终端)** 是一种 Electronic(电子) 或 Electromechanical(机电) 硬件设备，可用于将数据输入计算机或计算机系统，以及从计算机或计算机系统中转录数据。**Teletype(电传打字机，简称 TTY)** 是早期硬件拷贝终端的一个例子，并且比计算机屏幕的使用早了几十年。早期的终端是廉价设备，与用于输入的打孔卡或纸带相比速度非常慢，但随着技术的改进和视频显示器的引用，终端将这些旧的交互形式推向整个行业，分时系统的发展，弥补了低效的用户打字能力，能够支持统一机器上的多个用户，每个用户都在自己的终端上操作。所以，现代我们将 **TTY 表示为终端**，是那种最基础的终端。

除了传统的硬件终端以外，我们还可以通过计算机程序模拟出硬件终端，这种功能称为 **Terminal Emulator(终端模拟器)**，而很多时候也称为 **Psedoterminal/Pseudotty(伪终端，简称 PTY)**。PTY 是一对提供双向通信通道的虚拟字符设备，通道的一端称为 master(**简称 PTMX**)，另一端称为 slave(**简称 PTS**)。Linux 通过 devpts 文件提供了对 PTS 功能的完整支持。

终端也分为多种类型，有多种程序可以为用户分配一个指定类型的终端

- TTY
- Psedoterminal(伪终端)
- CLI
- TUI

名词有很多，但是至今位置没有一个明确的标准定义，大家都是拿起来就用~

# ANSI 转义码

> 参考：
>
> - [Wiki, ANSI_escape_code](https://en.wikipedia.org/wiki/ANSI_escape_code)
> - [Manual(手册)，console_codes(4)](https://man7.org/linux/man-pages/man4/console_codes.4.html)
> - [<计算机知识>：ANSI转义序列以及输出颜色字符详解](https://www.cnblogs.com/xiaoqiangink/p/12718524.html)
> - https://gist.github.com/fnky/458719343aabd01cfb17a3a4f7296797
> - https://blog.csdn.net/ScilogyHunter/article/details/106874395
> - https://learnku.com/articles/26231
> - [51CTO，使用ANSI改变终端输出样式](https://blog.51cto.com/u_15069477/3439423)

**ANSI escape code(ANSI 转义码)** 也称为 **ANSI escape sequences(ANSI 转义序列)**，是一种[带内信令](https://en.wikipedia.org/wiki/In-band_signaling)标准，用于控制终端上的光标位置、颜色、字体样式、等等终端相关功能。可以算是一种特殊的[控制字符](/docs/8.通用技术/编码与解码/字符的编码与解码/控制字符.md)（也称为控制码）。

ANSI 序列是在二十世纪七十年代引入的标准，用以取代特定终端供应商的序列，并在二十世纪八十年代早期开始在计算机设备市场上广泛使用。与早期缺少光标移动功能的系统相比，新生的电子公告板系统（BBS）使用 ANSI 序列改进其显示。正是因为这个原因，ANSI 序列变成了所有制造商共同采用的标准。

转义码由 **Escape(ESC) 符号**后跟**普通字符**组成。终端在收到一个 ESC 时，就会把其后面的几个字符当作主机发送的命令来对待，并对该字符序列作出诠释。在识别出有效的转义序列结束后，终端执行主机的控制命令。随后所接收到的字符将仍然会显示在屏幕上（除非它们也是控制字符或者转义字符序列）。

转义码具有不同的长度，所有序列都以 [ASCII](/docs/8.通用技术/编码与解码/字符的编码与解码/ASCII%20表.md) 的 **ESC 符号**开头，第二个字节则是 0x40–0x5F（i.e. `@A–Z[\]^_`）范围内的字符。

- ESC 字符可以通过如下几种方式表示出来
  - 使用键盘按键 `Ctrl + [`
  - 十六进制 `0x1B`
  - 十进制 `27`
  - 八进制 `\033`
  - Unicode `\u001b`

可以参考 Man 手册中，有完整的 ANSI 转义码列表

![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/ANSI/202308231738266.png)

## CSI

**Control Sequence Introducer(控制序列引入器，简称 CSI)** 以 `ESC[` 开头，大致可分为四类：

- 光标移动指令
- 清屏指令
- 字符渲染指令，通常指 [SGR](#sgr)
- 终端控制指令

CSI 序列 由 `ESC [` 以及若干个 `参数字节` 和 若干个`中间字节` 以及 一个 `最终字节` 组成.

| 组成部分 | 字符范围  | ASCII                 |
| -------- | --------- | --------------------- |
| 参数字节 | 0x30–0x3F | `0–9:;<=>?`             |
| 中间字节 | 0x20–0x2F | `空格、!"#$%&'()*+,-./` |
| 最终字节 | 0x40–0x7E | `@A–Z[]^_`a–z{}~`       |

一些 CSI 序列（不完整）

![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/ANSI/202308231648309.png)

WIi 中代码的表述模糊，这里用一个例子描述一下。CSI 就是 `ESC[`，比如 EL 功能，我们如果使用 Go 语言的话，可以这么写：

```go
fmt.Printf(" \033[2J")
```

其中 `\033[2J` 就是 `CSI n J`。虽然表中有空格，但是实际使用时并不需要添加空格，表中的空格只是让人们看着更舒服。

- CSI 就是 `ESC[`，i.e. `\033[`
  - 注意，这里的 ESC 使用 8 进制表示，我们可以可以使用 16 进制表示，比如这样：`fmt.Printf(" %c[2J",0x1B)` 或者 `fmt.Printf(" \x1b[2J")`
- n 就是 2
- J 就是 J

这个代码可以让终端实现清屏的效果，就像按 Ctrl + L 或者使用 clear 命令。

### SGR

**Select Graphic Rendition(选择图形呈现，简称 SGR)** 用于字符渲染指令，格式为 `ESC[Nm`，在 CSI 的开头面跟数字，并以 m 结尾。n 的取值范围是 0-107

![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/ANSI/202308231700074.png)

**背景色：**

40:黑 41:深红 42:绿 43:黄色 44:蓝色 45:紫色 46:深绿 47:白色

**字体颜色：**

30:黑 31:红 32:绿 33:黄 34:蓝色 35:紫色 36:深绿 37:白色

linux 终端下输出带颜色的文字只需在文字前面添加如下格式

```none
echo -e "\033[字背景颜色;文字颜色m字符串\033[0m"`
```

其中 `\033` 是ESC的八进制，`\033[`即告诉终端后面是设置颜色的参数，显示方式，前景色，背景色均是数字

例如:

`echo -e "\033[41;36m something here \033[0m"`

其中 41 代表背景色，36 代表字的颜色

可以将所有控制参数都用上，也可以只使用前景色或背景色.

但有一点要注意，如果输出带颜色的字符后并没有恢复终端默认设置，后续的命令输出仍旧会采用之前的颜色，如果是在脚本中设置了颜色而未恢复，则整个脚本的输出都会采用之前的颜色，因此如果不希望影响后面文字的输出，最好是在输出带颜色的文字之后恢复终端默认设置，如下

如果想设置文字颜色:

```shell
echo -e "\033[30m 黑色字 \033[0m"
echo -e "\033[31m 红色字 \033[0m"
echo -e "\033[32m 绿色字 \033[0m"
echo -e "\033[33m 黄色字 \033[0m"
echo -e "\033[34m 蓝色字 \033[0m"
echo -e "\033[35m 紫色字 \033[0m"
echo -e "\033[36m 天蓝字 \033[0m"
echo -e "\033[37m 白色字 \033[0m"
```

如果是简单设置背景颜色:

```shell
echo -e "\033[40;37m 黑底白字 \033[0m"
echo -e "\033[41;37m 红底白字 \033[0m"
echo -e "\033[42;37m 绿底白字 \033[0m"
echo -e "\033[43;37m 黄底白字 \033[0m"
echo -e "\033[44;37m 蓝底白字 \033[0m"
echo -e "\033[45;37m 紫底白字 \033[0m"
echo -e "\033[46;37m 天蓝底白字 \033[0m"
echo -e "\033[47;30;1m 白底黑字加粗 \033[0m"
```

# CLI 的 Args、Flag、Options

关于 CLI 的术语有很多，比如 Argument(参数)、Flag(标志)、Option(选项) 等

本质上，命令及其参数只是一个字符串而已，字符串的含义是由 Shell 来解释的，对于 Shell 来说，命令和参数、参数和参数之间是由空白符分割的。除此之外，什么父命令、子命令、本地参数、单横线、双横线 还是其他字符开头都没关系，就仅仅是字符串而已，这些字符串由 Shell 传递给将要执行的程序中。

- **Argument(参数)** # 就是对命令后面一串空白符分割的字符串的称呼
- **Flag(标志)** # 这种类型的参数可以将某个值跟代码中的某个变量关联起来。
- **Option(选项)** # Flag 赋予了我们通过 CLI 直接给代码中某个变量赋值的能力。那么如果我没有给这个变量赋值呢，程序还能运行下去么？如果不能运行，则这个 Flag 就是必选的，否则就是可选的。那么这些 Flag 或者 Argument 从这种角度将可以称为 Option。也就是可选的 Flag；或者称为可选的 Argument。
