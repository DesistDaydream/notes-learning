---
title: 比Linus更厉害的程序员，这家伙一定是穿越回来的！
---

Linus Torvalds 是个非常厉害的程序员，因为他有两个名扬天下的作品：Linux 和 Git。

如果单论技术能力，有一个人，也许比 Linus 更强。

我在看他主页项目列表的时候，感觉头都炸了。

他开发了著名的模拟器 QEMU 和音视频处理库 FFmpeg，仅仅是这两项就超越绝大部分程序员了，他还写过 C 编译器，OpenGL 实现，LTE 软基站，JS 引擎，让 Linux 在浏览器中跑起来，甚至还创造了计算圆周率的世界纪录......

也就是说，这位老兄在操作系统、模拟器、多媒体、计算机图形学、编译器、编程语言、通信、甚至数学等领域跳来跳去，一年开发一个我一辈子都写不出的软件！

他写的程序还总是比别的程序小几个数量级，快几个数量级！

这也太变态了吧？！

不得不承认，这个世界上真的有天才的存在。

他就是法国程序员 Fabrice Bellard。

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/8bd22c62-e692-4080-9eca-a4a4b1a2355d/640)

我们来看看他的神奇之路。

**01**

**压缩软件**

Fabrice Bellard 出生于法国东南部的一个小城市格勒诺布尔，17 岁上高中的时候发现自己的电脑硬盘又小又贵，他就想着如何节省硬盘空间，于是用汇编语言开发了一个压缩程序 LZEXE。

LZEXE 压缩效果极好，他的朋友也 Copy 过来使用，并且放到了 BBS 上，一下子火了，成了 DOS 时代最火的压缩工具之一。

1996 年，24 岁的他写了一个 Java 虚拟机，可以把 Java 代码编译成 C 代码执行。

**02**

**圆周率算法**

1997 年，他对圆周率的计算产生了兴趣，通过改良 Bailey–Borwein–Plouffe 公式，提出了最快速的计算圆周率的算法，极大地降低了时间复杂度。

从此他在数学社区建立了自己的声望，新算法也被命名为 Bellard 公式。

**03**

**TinyGL(计算机图形学）**

1998 年，Bellard 在开发一个虚拟现实引擎项目的时候，需要用到 OpenGL，这是一个用于渲染 2D、3D 矢量图形的跨语言、跨平台的 API，OpenGL 的开源实现是 Mesa，Bellard 觉得 Mesa 太慢了，于是自己写了一个，这就是 TinyGL。

TinyGL 实现了 OpenGL 的子集，比 Mesa 或者其他商用实现（如 Solaris OpenWin OpenGL）快得多，占用的资源少得多，并且比任何一个都小几个数量级，Bellard 再次展示了他在编写高效 C 代码方面的超高技能。

**04**

**FFmpeg(音频视频多媒体)**

2000 年，他化名 Gérard Lantau，开始了他最重要和最受广泛认可的项目之⼀：FFmpeg。

FFmpeg 是名副其实的数字视频和音频的“瑞士军刀”，在视频软件和商业网站中无处不在：VLC，YouTube， iTunes ，它的强大之处不用我介绍了吧。

**05**

**C 语言混淆竞赛**

在创建 FFmpeg 不久，2000 年和 2001 年，Bellard 向国际 C 语言代码混淆竞赛 (IOCCC) 提交了两个参赛作品，并且连续两次获奖。

要知道，IOCCC 是最具创造性的 C 语言混淆竞赛，是程序员充分展示自己智力的最佳舞台，能赢一次就非常厉害了，而 Bellard 轻轻松松就搞定了两次。

下面是 Bellard 在 2000 年的获奖作品，使用快速傅里叶变换在较短时间内计算出已知的最大素数(2^6972593-1)

\`int m = 754974721, N, t\[1 << 22], a, \*p, i, e = 1 << 22, j, s, b, c, U;

f (d)

{

for (s = 1 << 23; s; s /= 2, d = d _ 1LL _ d % m)

if (s < N)

for (p = t; p < t + N; p += s)

for (i = s, c = 1; i; i--)

b = _p + p\[s], p\[s] = (m + _p - p\[s]) \*

1LL _ c % m, _p++ = b % m, c = c _ 1LL _ d % m;

for (j = 0; i < N - 1;)

{

for (s = N / 2; !((j ^= s) & s); s /= 2);

if (++i < j)

a = t\[i], t\[i] = t\[j], t\[j] = a;

}

}

main ()

{

\*t = 2;

U = N = 1;

while (e /= 2)

{

N \*= 2;

U = U _ 1LL _ (m + 1) / 2 % m;

f (362);

for (p = t; p < t + N;)

\_p++ = (\_p \* 1LL ** p % m) ** U % m;

f (415027540);

for (a = 0, p = t; p < t + N;)

a += (6972593 & e ? 2 : 1) \** p, *p++ = a % 10, a /= 10;

}

while (!\*--p);

t\[0]--;

while (p >= t)

printf ("%d", \*p--);

}

\`

**06**

**TinyCC(最快的编译器）**

从 2001 年的比赛中还产生了一个副产品：TinyCC，这是世界上最快、最小的 C 语言编译器，比其他大多数 C 编译器都要小几个数量级。

为了证明 TinyCC 的威力，Bellard 基于 TinyCC 开发了一个只有 138K 的 TCCBoot，可以在 15 秒以内编译完 Linux 内核并且启动，实在太吓人了。

**07**

**QEMU（模拟器）**

2005 年，Bellard 又发布了一个爆炸性项目 QEMU，这是一个开源的模拟器，可以用软件的方式来模拟 CPU，内存，I/O 设备，给操作系统营造一个运行在硬件中的假象。

可以想想，开发这样的软件不但需要对操作系统极其了解，还得掌握极其广泛的硬件知识，细节非常多，其难度甚至比操作系统都高。

**08**

**创造 PI 的世界纪录**

2009 年，Bellard 又去玩数学了，他宣布把圆周率小数点后 2.7 万亿位以后，仅仅使用了一台普通的 PC，创造了世界纪录。

此前的世界纪录是由排名世界第 47 位的 T2K Open 超级计算机创造的，而 Bellard 这台桌面电脑不到 2000 欧元，配置仅为：2.93GHz Core i7 CPU，6GB 内存，7.5TB 硬盘。

一个人加一台电脑，竟然击败了超级计算机。

**09**

**JSLinux（模拟器）**

2011 年，他的兴趣又转到了 JavaScript 身上，居然用 JavaScript 写了一个 PC 模拟器，让 Linux，Windows 可以在浏览器中运行起来。

这个模拟器仿真了一个 32 位的 x86 兼容处理器，一个 8259 可编程中断控制器，一个 8254 可编程中断计时器，和一个 16450 UART。

不仅支持命令行，还支持图形界面，看到 Windows 2000 竟然在浏览器中跑了起来，那种震撼的感觉，只能用卧槽来形容了！

**10**

**LTE 软基站（通信）**

2012 年，Bellard 的兴趣再次转移，一个人花了 10 个月时间，在一台 PC 上居然实现了一个运行效率极高 LTE 软基站。支持 LTE TDD/FDD，NB-IoT、eMTC，最大可支持 5 载波 2x2 MIMO 或 3 载波 4x4 MIMO。

这一切，只需要一个拥有 i7 4 核 CPU 的 PC 就够了。

**11**

**QuickJS （JavaScript 解释器）**

2019 年，Bellard 发布了一个嵌入式的 JavaScript 执行引擎 QuickJS。

QuickJS 支持 ES2020，小巧并且易于嵌入，只有几个 C 文件，没有任何其他外部依赖。

它运行速度很快，在一个单核 CPU 上可以在 95 秒内完成 69000 个 ECMAScript 测试。

我只是挑了 Bellard 开发的部分软件，在他的网站 bellard.org 还有很多，最让人震撼的是，这些软件覆盖了计算机科学的各个领域，千差万别。

Bellard 给人的印象是，他可以轻易进入一个他觉得有趣的领域，成为这个领域的专家，留下一个让其他人愿意花费数年时间维护的软件，自己则轻飘飘地离开，进入下一个领域。

有人问他为什么要研究这么多不同的东西时，他说：我讨厌一直做同样的事情，所以一定要切换不同的项目来玩......

Just for fun，这是 Linus 的口号，看来各个大神的追求都是一样的啊。

Bellard 对金钱或者名声不感兴趣（用化名做开源项目就是证明），他极少接受媒体的采访，互联网上他的资料非常少，远不如 Linus 那样声名远扬。

但是如果你如果你看过他那简陋的个人主页，bellard.org，看看那些展示了惊人的深度和广度的项目，绝对会被震撼。

Bellard 一定是穿越回来的，或者一定是在上帝模式下编程。

[程序员，你得选准跑路的时间！](http://mp.weixin.qq.com/s?__biz=MzAxOTc0NzExNg==&mid=2665518834&idx=1&sn=5dcedf3fc7e318f50cd32fffe38c429d&chksm=80d66eb1b7a1e7a72165940271001720e7853da6d324e6d9a0055b4a8c636d4bdbb7c1f1bfb9&scene=21#wechat_redirect)

[两年，我学会了所有的编程语言！](http://mp.weixin.qq.com/s?__biz=MzAxOTc0NzExNg==&mid=2665517522&idx=1&sn=758759e90d395311bfd33a68e2eeb116&chksm=80d66991b7a1e0874b452a9d3627fe995231eb5d80d88b5234b6b54e631b5c04827b16800e8d&scene=21#wechat_redirect)

[你们这些偷代码的程序员！](http://mp.weixin.qq.com/s?__biz=MzAxOTc0NzExNg==&mid=2665520720&idx=1&sn=404be15f144ffc7d0dccae6b5bd1a46f&chksm=80d66613b7a1ef05df940c46fe3e0611915ce037535c8e16a3b6233d894208a38be21e9376ce&scene=21#wechat_redirect)

[程序员的宿命](http://mp.weixin.qq.com/s?__biz=MzAxOTc0NzExNg==&mid=2665520490&idx=1&sn=6dd71186fd2d648c29d69f541a8c7baf&chksm=80d66529b7a1ec3f9e266072ff2d4839df4082fa702d84be76e51d15326b34b7034278ce7ee4&scene=21#wechat_redirect)

[芯片战争 70 年，真正的王者即将现身](http://mp.weixin.qq.com/s?__biz=MzAxOTc0NzExNg==&mid=2665520306&idx=1&sn=42b0fefa3be505d8a1926214fbf58e3e&chksm=80d664f1b7a1ede73cfe7488991cf78ead735357be8ed2528ab7d275d47712cc73754f802be6&scene=21#wechat_redirect)！

[宇宙第一 IDE 到底是谁？](http://mp.weixin.qq.com/s?__biz=MzAxOTc0NzExNg==&mid=2665519602&idx=1&sn=f85e22d32d4919116afbb2808f1a151a&chksm=80d661b1b7a1e8a70568d243a20c210eaf474650efc61e809dcf8e5d534be1d7f15ce1145295&scene=21#wechat_redirect)

[HTTP Server ：一个差生的逆袭](http://mp.weixin.qq.com/s?__biz=MzAxOTc0NzExNg==&mid=2665519532&idx=1&sn=d17efeb01b6114fe20ea60d08c7836b2&chksm=80d661efb7a1e8f9981903e2d7f0754449769973c80e5c7539930e89d7a3ffe0bdd74aabd8d5&scene=21#wechat_redirect)

[Javascript: 一个屌丝的逆袭](http://mp.weixin.qq.com/s?__biz=MzAxOTc0NzExNg==&mid=2665513059&idx=1&sn=a2eaf97d9e3000d15a33681d1b720463&chksm=80d67820b7a1f136d73b874e5784f06ef34eb38d7750712a5b9c7532ad44f18a3d43d2edbfbd&scene=21#wechat_redirect)

[我是一个线程](http://mp.weixin.qq.com/s?__biz=MzAxOTc0NzExNg==&mid=416915373&idx=1&sn=f80a13b099237534a3ef777d511d831a&chksm=06ef67ee3198eef85a668f7197709b54bcd9508087ace8f6eaa11d34eb4e5392c5190a77d530&scene=21#wechat_redirect)

[TCP/IP 之大明邮差](http://mp.weixin.qq.com/s?__biz=MzAxOTc0NzExNg==&mid=2665513094&idx=1&sn=a2accfc41107ac08d74ec3317995955e&chksm=80d678c5b7a1f1d3e122ae99f368671acb7f6baecebb6488c8ca63caa7f1ef28fea1a4dc1ef0&scene=21#wechat_redirect)

[一个故事讲完 Https](http://mp.weixin.qq.com/s?__biz=MzAxOTc0NzExNg==&mid=2665513779&idx=1&sn=a1de58690ad4f95111e013254a026ca2&chksm=80d67b70b7a1f26697fa1626b3e9830dbdf4857d7a9528d22662f2e43af149265c4fd1b60024&scene=21#wechat_redirect)

[CPU 阿甘](http://mp.weixin.qq.com/s?__biz=MzAxOTc0NzExNg==&mid=2665513017&idx=1&sn=5550ee714abd36d0b580713f673e670b&scene=21#wechat_redirect)
