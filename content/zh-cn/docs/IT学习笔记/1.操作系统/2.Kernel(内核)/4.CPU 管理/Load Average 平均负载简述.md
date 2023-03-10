---
title: Load Average 平均负载简述
---

# 概述

<https://blog.csdn.net/u011183653/article/details/19489603>

<https://blog.csdn.net/slvher/article/details/9199439>

# Load 与 Load Average

    LoadAverage = calc_load(TASK_RUNNING + TASK_UNINTERRUPTIBLE,n)

Load 是此时此刻 CPU 正在处理的进程数。进程可运行状态时，它处在一个运行队列 run queue 中，与其他可运行进程争夺 CPU 时间。 系统的 load 是指正在运行 running 和准备好运行 runnable 以及 不可中断睡眠 的进程的总数。比如现在系统有 2 个正在运行的进程，3 个可运行进程，那么系统的 load 就是 5

Load Average 为在特定时间间隔内运行队列中(在 CPU 上运行或者等待运行多少进程)的平均进程数。如果一个进程满足以下条件则其就会位于运行队列中：

1. 它没有在等待 I/O 操作的结果

2. 它没有主动进入等待状态(也就是没有调用’wait’)

3. 没有被停止(例如：等待终止)

在 Linux 中，进程分为三种状态，一种是阻塞的进程 blocked process，一种是可运行的进程 runnable process，另外就是正在运行的进程 running process。当进程阻塞时，进程会等待 I/O 设备的数据或者系统调用。

一、查看系统负荷

如果你的电脑很慢，你或许想查看一下，它的工作量是否太大了。

在 Linux 系统中，我们一般使用 uptime 命令查看（w 命令和 top 命令也行）。（另外，它们在苹果公司的 Mac 电脑上也适用。）

你在终端窗口键入 uptime，系统会返回一行信息。

    [root@lichenhao ~]# uptime
     17:00:00 up 2 days,  2:53,  1 user,  load average: 0.09, 0.05, 0.01

这行信息的后半部分，显示"load average"，它的意思是"系统的平均负荷"，里面有三个数字，我们可以从中判断系统负荷是大还是小。

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/pagncp/1616167091810-b44dc550-47cf-44e2-94c4-de77c383db2c.jpeg)

为什么会有三个数字呢？你从手册中查到，它们的意思分别是 1 分钟、5 分钟、15 分钟内系统的平均负荷。

如果你继续看手册，它还会告诉你，当 CPU 完全空闲的时候，平均负荷为 0；当 CPU 工作量饱和的时候，平均负荷为 1。

那么很显然，"load average"的值越低，比如等于 0.2 或 0.3，就说明电脑的工作量越小，系统负荷比较轻。

但是，什么时候能看出系统负荷比较重呢？等于 1 的时候，还是等于 0.5 或等于 1.5 的时候？如果 1 分钟、5 分钟、15 分钟三个值不一样，怎么办？

二、一个类比

判断系统负荷是否过重，必须理解 load average 的真正含义。下面，我根据"Understanding Linux CPU Load"这篇文章，尝试用最通俗的语言，解释这个问题。

首先，假设最简单的情况，你的电脑只有一个 CPU，所有的运算都必须由这个 CPU 来完成。

那么，我们不妨把这个 CPU 想象成一座大桥，桥上只有一根车道，所有车辆都必须从这根车道上通过。（很显然，这座桥只能单向通行。）

系统负荷为 0，意味着大桥上一辆车也没有。

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/pagncp/1616167091810-e462799a-5a71-46cc-8f5a-560ab51b7f6b.jpeg)

系统负荷为 0.5，意味着大桥一半的路段有车。

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/pagncp/1616167091819-cf628c32-b0c1-4c81-a5fa-35c30d663944.jpeg)

系统负荷为 1.0，意味着大桥的所有路段都有车，也就是说大桥已经"满"了。但是必须注意的是，直到此时大桥还是能顺畅通行的。

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/pagncp/1616167091800-0fa2f777-ecc8-43c3-8c10-029b54de5eb4.jpeg)

系统负荷为 1.7，意味着车辆太多了，大桥已经被占满了（100%），后面等着上桥的车辆为桥面车辆的 70%。以此类推，系统负荷 2.0，意味着等待上桥的车辆与桥面的车辆一样多；系统负荷 3.0，意味着等待上桥的车辆是桥面车辆的 2 倍。总之，当系统负荷大于 1，后面的车辆就必须等待了；系统负荷越大，过桥就必须等得越久。

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/pagncp/1616167091822-b290f2e0-2edb-4775-a7a1-fecd0d32bb3c.jpeg)

CPU 的系统负荷，基本上等同于上面的类比。大桥的通行能力，就是 CPU 的最大工作量；桥梁上的车辆，就是一个个等待 CPU 处理的进程（process）。

如果 CPU 每分钟最多处理 100 个进程，那么系统负荷 0.2，意味着 CPU 在这 1 分钟里只处理 20 个进程；系统负荷 1.0，意味着 CPU 在这 1 分钟里正好处理 100 个进程；系统负荷 1.7，意味着除了 CPU 正在处理的 100 个进程以外，还有 70 个进程正排队等着 CPU 处理。

Note：如果 CPU 一分钟的时间一直处理 1 个进程，那么 1 分钟的负载也是 1。

为了电脑顺畅运行，系统负荷最好不要超过 1.0，这样就没有进程需要等待了，所有进程都能第一时间得到处理。很显然，1.0 是一个关键值，超过这个值，系统就不在最佳状态了，你要动手干预了。

三、系统负荷的经验法则

1.0 是系统负荷的理想值吗？

不一定，系统管理员往往会留一点余地，当这个值达到 0.7，就应当引起注意了。经验法则是这样的：

当系统负荷持续大于 0.7，你必须开始调查了，问题出在哪里，防止情况恶化。

当系统负荷持续大于 1.0，你必须动手寻找解决办法，把这个值降下来。

当系统负荷达到 5.0，就表明你的系统有很严重的问题，长时间没有响应，或者接近死机了。你不应该让系统达到这个值。

四、多处理器

上面，我们假设你的电脑只有 1 个 CPU。如果你的电脑装了 2 个 CPU，会发生什么情况呢？

2 个 CPU，意味着电脑的处理能力翻了一倍，能够同时处理的进程数量也翻了一倍。

还是用大桥来类比，两个 CPU 就意味着大桥有两根车道了，通车能力翻倍了。

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/pagncp/1616167091818-7b6aaceb-456c-4625-a2b4-cba8501b8858.jpeg)

所以，2 个 CPU 表明系统负荷可以达到 2.0，此时每个 CPU 都达到 100%的工作量。推广开来，n 个 CPU 的电脑，可接受的系统负荷最大为 n.0。

五、多核处理器

芯片厂商往往在一个 CPU 内部，包含多个 CPU 核心，这被称为多核 CPU。

在系统负荷方面，多核 CPU 与多 CPU 效果类似，所以考虑系统负荷的时候，必须考虑这台电脑有几个 CPU、每个 CPU 有几个核心。然后，把系统负荷除以总的核心数，只要每个核心的负荷不超过 1.0，就表明电脑正常运行。

怎么知道电脑有多少个 CPU 核心呢？

"cat /proc/cpuinfo"命令，可以查看 CPU 信息。"grep -c 'model name' /proc/cpuinfo"命令，直接返回 CPU 的总核心数。

六、最佳观察时长

最后一个问题，"load average"一共返回三个平均值----1 分钟系统负荷、5 分钟系统负荷，15 分钟系统负荷，----应该参考哪个值？

如果只有 1 分钟的系统负荷大于 1.0，其他两个时间段都小于 1.0，这表明只是暂时现象，问题不大。

如果 15 分钟内，平均系统负荷大于 1.0（调整 CPU 核心数之后），表明问题持续存在，不是暂时现象。所以，你应该主要观察"15 分钟系统负荷"，将它作为电脑正常运行的指标。

\==========================================

\[参考文献]

1. Understanding Linux CPU Load

2. Wikipedia - Load (computing)

（完）

# 简单说下 CPU 负载和 CPU 利用率的区别

1）CPU 利用率：显示的是程序在运行期间实时占用的 CPU 百分比

2）CPU 负载：显示的是一段时间内正在使用和等待使用 CPU 的平均任务数。

CPU 利用率高，并不意味着负载就一定大。

举例来说：

如果有一个程序它需要一直使用 CPU 的运算功能，那么此时 CPU 的使用率可能达到 100%，但是 CPU 的工作负载则是趋近于"1"，因为 CPU 仅负责一个工作！

如果同时执行这样的程序两个呢？CPU 的使用率还是 100%，但是工作负载则变成 2 了。所以也就是说，当 CPU 的工作负载越大，代表 CPU 必须要在不同的工作之间

进行频繁的工作切换。

\------------------------下面通过一个电话亭打电话的比喻来说明这两者之间的区别------------------------

某公用电话亭，有一个人在打电话，四个人在等待，每人限定使用电话一分钟，若有人一分钟之内没有打完电话，只能挂掉电话去排队，等待下一轮。

电话在这里就相当于 CPU，而正在或等待打电话的人就相当于任务数。

在电话亭使用过程中，肯定会有人打完电话走掉，有人没有打完电话而选择重新排队，更会有新增的人在这儿排队，这个人数的变化就相当于任务数的增减。

为了统计平均负载情况，我们 5 分钟统计一次人数，并在第 1、5、15 分钟的时候对统计情况取平均值，从而形成第 1、5、15 分钟的平均负载。

有的人拿起电话就打，一直打完 1 分钟，而有的人可能前三十秒在找电话号码，或者在犹豫要不要打，后三十秒才真正在打电话。如果把电话看作 CPU，人数看

作任务，我们就说前一个人（任务）的 CPU 利用率高，后一个人（任务）的 CPU 利用率低。当然， CPU 并不会在前三十秒工作，后三十秒歇着，只是说，有的程

序涉及到大量的计算，所以 CPU 利用率就高，而有的程序牵涉到计算的部分很少，CPU 利用率自然就低。但无论 CPU 的利用率是高是低，跟后面有多少任务在排队

没有必然关系。

# 理解 LINUX LOAD AVERAGE 的误区

Load average 的概念源自 UNIX 系统，虽然各家的公式不尽相同，但都是用于衡量正在使用 CPU 的进程数量和正在等待 CPU 的进程数量，一句话就是 runnable processes 的数量。所以 load average 可以作为 CPU 瓶颈的参考指标，如果大于 CPU 的数量，说明 CPU 可能不够用了。

但是，Linux 上不是这样的！

Linux 上的 load average 除了包括正在使用 CPU 的进程数量和正在等待 CPU 的进程数量之外，还包括 uninterruptible sleep 的进程数量。通常等待 IO 设备、等待网络的时候，进程会处于 uninterruptible sleep 状态。Linux 设计者的逻辑是，uninterruptible sleep 应该都是很短暂的，很快就会恢复运行，所以被等同于 runnable。然而 uninterruptible sleep 即使再短暂也是 sleep，何况现实世界中 uninterruptible sleep 未必很短暂，大量的、或长时间的 uninterruptible sleep 通常意味着 IO 设备遇到了瓶颈。众所周知，sleep 状态的进程是不需要 CPU 的，即使所有的 CPU 都空闲，正在 sleep 的进程也是运行不了的，所以 sleep 进程的数量绝对不适合用作衡量 CPU 负载的指标，Linux 把 uninterruptible sleep 进程算进 load average 的做法直接颠覆了 load average 的本来意义。所以在 Linux 系统上，load average 这个指标基本失去了作用，因为你不知道它代表什么意思，当看到 load average 很高的时候，你不知道是 runnable 进程太多还是 uninterruptible sleep 进程太多，也就无法判断是 CPU 不够用还是 IO 设备有瓶颈。

参考资料：<https://en.wikipedia.org/wiki/Load_(computing)>

“Most UNIX systems count only processes in the running (on CPU) or runnable (waiting for CPU) states. However, Linux also includes processes in uninterruptible sleep states (usually waiting for disk activity), which can lead to markedly different results if many processes remain blocked in I/O due to a busy or stalled I/O system.“
