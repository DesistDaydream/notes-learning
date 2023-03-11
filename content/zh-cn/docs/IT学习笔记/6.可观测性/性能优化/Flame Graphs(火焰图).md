---
title: Flame Graphs(火焰图)
---

# 概述

> 参考：
> - [GitHub 项目，brendangregg/FlameGraph](https://github.com/brendangregg/FlameGraph)
> - [官方文档](http://www.brendangregg.com/flamegraphs.html)
> - [论文](https://queue.acm.org/detail.cfm?id=2927301)

<https://www.ruanyifeng.com/blog/2017/09/flame-graph.html>
<https://zhuanlan.zhihu.com/p/73385693>

可以生成火焰图的工具：
- [perf 性能分析工具](/docs/IT学习笔记/1.操作系统/X.Linux%20管理/Linux%20系统管理工具/perf%20性能分析工具/perf%20性能分析工具.md)

# 前言

在没有读《性能之巅》这本书之前，就听说过火焰图。当时学习同事用 go 写的项目代码，发现里边有个文件夹叫火焰图，觉得名字很好玩，就百度了一下，惊叹还有这种操作。不过当时并没有听过 Brendan Gregg 的大名，因为懒也没有深入学习。这次找到了 Brendan Gregg 的 blog，也了解了一点动态追踪技术的知识，决心要好好学习一下。
于是就找到了一切开始的地方： Brendan Gregg 写的论文《[The Flame Graph](https://link.zhihu.com/?target=https%3A//queue.acm.org/detail.cfm%3Fid%3D2927301) 》
作为一个英语菜鸡，从来都没有读过英文论文。正好借这次机会尝试一下，看能不能点个新的技能点。结果尝试才发现，真的好难～～刚开始，读一小段就开始犯困。于是坚持每天强迫自己从头开始重读一遍。花了差不多一周时间，总算能集中注意力的读完。
然后我就想，老是吐槽各种汉化的国外优秀的技术书籍生涩难懂，何不亲自试一试呢？于是就有了今天的这篇学习笔记。

# 火焰图

## 让软件执行情况可视化，是性能分析、调试的利器

Brendan Gregg, Netflix
日常工作中，我们需要理解软件对系统资源的使用情况。比如对于 cpu，我们想知道当前软件究竟使用了多少 cpu？软件更新以后又变化了多少？剖析器(profilers)可以用来分析这样的问题，帮助软件开发者优化代码，指导软件使用者调优运行环境。但是 profile 通常都很长，太长的输出分析和理解起来都很不方便。火焰图作为一种新的 profile 可视化方式，可以让我们更直观，更方便的理解、分析问题。
在像“Netflix 云微服务架构”这种软件升级迭代迅速的环境中，快速理解 profiles 尤为重要。同时，对 profile 的快速的理解也有助于我们更好的研究其他人编写的软件。
火焰图可以用多种 profilers（包括资源和系统事件）的输出生成，本文以 cpu 为例，介绍了火焰图的用法以及其可以解决的各种实际问题。

### profile 的理解

profile 有 剖面、剖面图 的含义，对于医学角度来说，如果不解剖看剖面图，也就无法看到一个生物内部的运行情况。同理，在性能分析领域，想要理解一个程序，也需要解剖它，看看它的剖面图。所以，**profile 就可以理解为一个应用程序的 剖面图**。只有看到剖面图，才能深入程序内部一探究竟~~~

## CPU Profiling

CPU 分析的一种常用技术是，使用像 Linux perf_events 和 DTrace 之类系统追踪工具的对 stack traces 进行采样。stack trace 显示了代码的调用关系，比如下面的 stack trace ,每个方法作为一行并按照父子关系从下到上排序。

    SpinPause
    StealTask::do_it
    GCTaskThread::run
    java_start
    start_thread

综合考虑采样成本、输出大小、应用环境，CPU profile 有可能是这样收集：对所有的 cpu，对 stack traces 以每秒 99 次的速度，连续采样 30 秒（使用 99 而不是 100，是为了防止采样周期与某些系统周期事件重合，影响采样结果）。对于一个 16 核的机器来说，输出结果可能会有 47520 个堆栈采样。可能会输出成千上万行文本。（ps:原文是 not 100, to avoid lock-step sampling 理解不了，所以按照书中的描述写的）
有些 profile 可以压缩输出，比如 DTrace，可以把相同的 stack traces 汇总到一起，只输出次数。这个优化还是蛮有用的，比如长时间的 for 循环，或者系统 idle 状态的 stack traces，都会被简化成一个。
Linux perf_events 还可以进一步压缩输出，通过合并相同的 substack，使用树形结构汇总输出。对于树的每个分枝，都可以统计 count 或百分比。如图一所示：
![](https://notes-learning.oss-cn-beijing.aliyuncs.com/hxnw7s/1619535684402-a78901d6-3770-4ff9-90d9-3089acd9510c.jpeg)
实际上，perf_events 和 DTrace 的输出，在很多情况下，足够分析问题使用了。但是也有时候，面对超长的输出，就像面对一堵写满字的高墙，分析其中某个堆栈就好比盲人摸象、管中窥豹。

### The Problem

为了分析“the Joyent public cloud”的性能问题，我们发明了火焰图。问题简单描述就是:某台服务器上部署了一个 mysql 服务，该服务的 cpu 使用率比预期的情况高了 40%。
我们使用 DTrace，以 997 Hz 的频率连续采样 stack traces 60 秒，尽管 DTrace 对输出进行了压缩，输出还是有 591622 行，包括 27053 个 stack traces，图二展示了输出结果，屏幕最下方显示的是调用最频繁的方法，说不定是问题的关键。
![](https://notes-learning.oss-cn-beijing.aliyuncs.com/hxnw7s/1619535684447-cc99cf85-c6e5-415f-995c-56ed19002211.jpeg)
最频繁的方法是`calc_sum_of_all_status()`,这个方法在执行 mysql 命令`show status`时被调用。也许有个客户端在疯狂执行这个命令做监控？
为了证明这个结论，用该方法采样的次数 5530，除以总的采样次数 348427。算出来这个方法只占用了 1.6%,远远不到 40%。看来得继续分析 profile。
如果继续按照调用频度，一个一个分析 stack traces，完全是一项体力劳动。看下图三就知道这是一项多么庞大的工作量。缩放以后，整个 DTrace 输出就像一个毫无特征的灰色图片。
![](https://notes-learning.oss-cn-beijing.aliyuncs.com/hxnw7s/1619535684409-cdd3aff2-75f7-47b6-9ed4-08b6df33d4ed.jpeg)
这简直是不可能完成的任务！那么，有没有更好的方法呢？
为了充分利用 stack traces 层次的特性，我发明了一种可视化原型，如图四所示，展示了跟图三相同的信息。图片之所以选择暖色调，是因为这种原型解释了为什么 cpu 很“hot”，也正是因为暖色调和火焰一样的形状。这种原型被命名为“火焰图”。（可交互的 svg 格式的图 4 可以在[http://queue.acm.org/downloads/2016/Gregg4.svg](https://link.zhihu.com/?target=http%3A//queue.acm.org/downloads/2016/Gregg4.svg) 体验）
![](https://notes-learning.oss-cn-beijing.aliyuncs.com/hxnw7s/1619535684411-236cd7fb-6362-4a8a-a78f-9a0d7256f2f0.jpeg)
使用火焰图可以很快的找到 profile 的主体部分，图片显示之前找到的 MySQL status 命令，只占 profile 的 3.28%，真正的大头是含有`join`的 mysql 语句。顺着这个线索，我们找到了根本问题，解决以后，cpu 使用率下降了 40%。

### Flame Graphs Explained

火焰图用相邻的 diagram 代表一堆 stack traces 。diagram 的形状像是一个倒着的冰锥 ![](https://notes-learning.oss-cn-beijing.aliyuncs.com/hxnw7s/1619535684456-a7bde9ca-6bbf-459f-8fa6-6ccf8bf48044.svg) 。这个冰锥一样的图形通常用来描述 cpu profile。
火焰图有以下特征：
• 一列 box 代表一个 stack trace，每个 box 代表一个 stack frame。
• y 轴显示了栈的深度，按照调用关系从下到上排列。最顶上的 box 代表采样时，on-CPU 的方法。每个 box 下边是该方法的调用者，以此类推。
• x 轴代表了整个采样的范围，注意 x 轴并不代表时间，所有 box 按照方法名称的字母顺序排列，这样的好处是，相同名称的 box，可以合并为一个。
• box 的宽度代表了该方法在采样中出现的频率。该频率与宽度成比例。
• 如果 box 很长，会显示完整的方法名称，如果很短，只会显示省略号或者 nothing。
• box 的颜色是随机选择的暖色调，这样有助于肉眼区分细长的塔状 boxes。当然也有其他配色方案，后面再说。
• profile 有可能是单线程、多线程、多应用甚至是多 host 的，如果需要，可以分解成子图。
• 还有很多其他的采样方式，box 的宽度除了频率以外，还可以表示多种其他的含义。比如 off-cpu 火焰图中，x 轴的宽度代表方法 block 的时间。
使用火焰图，整个 profile 一目了然，可以方便的定位到感兴趣的位置。火焰图成了软件执行情况的导航图。
除了这种可交互的展示方式，火焰图也可以方便的保存为静态图片的格式，方便打印出来。

### Interactivity

火焰图可以支持交互功能，可以显示更多细节、改进导航和执行计算。
原生的火焰图使用嵌入式 javascript 生成一张可交互的 svg 图片，提供了三种交互特性：鼠标 hover 显示详情、点击缩放和搜索。
**hover 显示详情**
当鼠标 hove 到 box 上，tooltip 内和图片左下方会显示方法的 full name，该方法的采样数量，以及百分比。格式如下：

    Function: mysqld'JOIN::exec (272,959 samples, 78.34 percent).

hover 这项特性有助于用户查看很窄的 box，显示百分比能够帮助用户量化代码路径的资源使用率，指导用户找到代码中急需优化的部分。
**点击缩放**
当用户点击一个 box 时，火焰图水平缩放，以显示局部的细节信息。该 box 下方的父 box 颜色变淡，表示只有部分被展示。点击 reset zoom 可以回到全局视图。
**搜索**
点击 search 或者按 ctrl+f 来使用搜索功能。搜索功能支持正则表达式，所有命中的 box 会被高亮并被显示为紫色。同时，图片右下角会显示命中方法的总百分比。如图五所示。（可交互的图五可以在这里体验： [http://queue.acm.org/downloads/2016/Gregg5.svg](https://link.zhihu.com/?target=http%3A//queue.acm.org/downloads/2016/Gregg5.svg).）
![](https://notes-learning.oss-cn-beijing.aliyuncs.com/hxnw7s/1619535684408-74db187a-b473-40ed-9232-ca926f24b102.jpeg)
搜索不仅可以方便定位方法，还可以高亮逻辑上相关的一组方法。比如输入"^ext4\_"显示所有跟 linux ext4 文件系统相关的方法。
有时候，多个代码路径都以相同的热点函数（比如自旋锁）结束，如果这些方法分布在图片 20 多个位置上，汇总他们的百分比就很麻烦。搜索功能可以解决这个问题。

### Instructions

火焰图有很多实现，原生的火焰图使用 Perl 编写，并且开放源码。包括采样在内，生成一张火焰图总共分三步：

1. 用户采样 stack traces (比如使用 Linux perf_events、DTrace 或者 Xperf)。
2. 将采样输出压缩为指定格式。我们已经编写了很多 perl 脚本处理各种 profiler 的输出。在项目中以"stackcollapse"前缀命名。
3. 使用 flamegraph.pl 生成火焰图. 该脚本使用 javascript 解析前边步骤的输出生成最终输出。

**指定压缩格式**指的是：把一个 stack traces 展示为一行，栈帧之间使用分号间隔，并在末尾的空格之后显示采样数量。应用名称，进程 id 之类的信息，可以用“.”分隔以后， 补充在 stack traces 之前。增加这些信息以后，生成的火焰图中，会按照这些前缀进行分组。
比如说，假如 profile 包含下面三个 stack traces:

    func_c  func_b  func_a  start_thread
    func_d  func_a  start_thread
    func_d  func_a  start_thread

压缩为指定格式以后，是这个样子的:

    start_thread;func_a;func_b;func_c 1
    start_thread;func_a;func_d 2

如果把应用名称（比如：java）也加在里面，则是：

    java;start_thread;func_a;func_b;func_c 1
    java;start_thread;func_a;func_d 2

设计这种中间格式的好处是，如果出现了新的 profiler，只需要编写转换器就可以使用火焰图。目前已经有 DTrace、perf_events、FreeBSD pmcstat, Xperf, SystemTap, Xcode Instruments, Intel VTune, Lightweight Java Profiler, Java jstack, and gdb 这么多可用的转换器。
flamegraph.pl 支持一些用户选项，比如说更改火焰图的 title。
下面是从采样（使用 perf ）到生成图片的一个具体的例子:

    # git clone https://github.com/brendangregg/FlameGraph
    # cd FlameGraph
    # perf record -F 99 -a -g -- sleep 60
    # perf script | ./stackcollapse-perf.pl | ./flamegraph.pl > out.svg

因为中间格式每个记录一行，在生成火焰图之前可以使用 grep/sed/awk 进行修改。使用其他 profiler 的教程参见[官方文档](https://link.zhihu.com/?target=https%3A//github.com/brendangregg/FlameGraph)。

### Flame Graph Interpretation

如何分析生成的火焰图:
• 火焰图顶部显示了采样过程中 on CPU 的方法。对 CPU profiles，这些方法直接占用 cpu 资源。对于其他的 profile，这些方法导致了相关的内核事件。
• 在火焰图顶部寻找“高原”状的方法，位于顶部的某个很宽的方法，表示其在采样中大量出现。对于 CPU profiles，这意味着这个方法经常在 CPU 上运行。
• 自顶向下看显示了调用关系，上边的方法被其下方的方法调用，以此类推。快速的从上往下浏览可以理解某个方法为什么被调用。
• 自底向上看显示了代码逻辑，提供了程序的全局视图。底部的方法会调用其顶部的多个方法，以此类推。自底向上看可以看到代码的分支形成的多个小型的“塔尖”。
• box 的宽度可以用来比较，更宽的 box 意味着在采样结果中更多的比例。
• 对于 cpu profiles 来说，如果 a 方法比 b 方法宽，有可能是因为 a 方法本身执行需要使用比 b 方法更多的 cpu。也有可能是 a 方法被调用的次数比 b 方法更频繁。采样的最终结果并不能体现一个方法被调用多少次，所以这两种情况都有可能。
• 如果一个方法顶部出现了两个“大塔尖”，导致火焰图中出现一个“大分叉”，这样的方法很值得研究。两个“塔尖”可能是被调用的两个子方法，也可能是条件语句的两个不同分支。

### Interpretation Example

为了更直观的让大家理解火焰图的含义，下面以图六作为例子。这是使用 cpu profile 生成的一张火焰图。
![](https://notes-learning.oss-cn-beijing.aliyuncs.com/hxnw7s/1619535684422-6c4e4abe-767b-4f3c-b725-b0380038f19a.jpg)
图片顶部显示说明`g()`使用 cpu 最频繁；虽然`d()`很宽，但是该方法直接使用 cpu 最少。`b()`和`c()`并不直接使用 cpu 资源，而他们的子方法使用。
`g()`底部的方法显示了调用关系：`g() `被` f()`调用,`f()`被 `d()`调用,以此类推。
对比`b()` 和` h()`的宽度可以发现，`b()`对 cpu 的使用率是`h()`的 4 倍。真正在 cpu 上执行的是他们的子方法。
该图的主分支是`a()` 调用了 `b()` 和 `h()`,原因有可能是`a()`中存在条件分支（比如一个 if 语句，如果为 true 执行`b()`,反之执行`h()`）,也有可能`a()`分成了两个步骤`b() `和` h()`。

### Other Code-Path Visualizations

正如图一所示，**Linux perf_events** 使用树形结构展示 cpu 使用率，这是另一种层级可视化方式。与火焰图相比，该方法并不能提供直观的全局视图。
**KCacheGrind** 使用有向无环图实现可视化，使用宽度自适应的 box 表示方法，使用箭头表示调用关系，box 和箭头上都标注了百分比。与 Linux perf_events 一样，图片缩小以后，也很难提供全局信息。
![](https://notes-learning.oss-cn-beijing.aliyuncs.com/hxnw7s/1619535684440-0ff268f5-3021-473d-aeaf-6fa1d4b5c051.jpeg)
**sunburst**布局跟火焰图的冰锥布局很像。不同的是 sunburst 使用了极坐标。sunburst 生成的图形很漂亮，但是却并不利于理解。与角度大小相比，人们更容易区分宽度。所以在直观性上火焰图更胜一筹 ![](https://notes-learning.oss-cn-beijing.aliyuncs.com/hxnw7s/1619535684460-0fa7cc81-0b41-49df-9df3-c52c3dc2c7a9.svg) 。
**Flame charts **灵感来源于火焰图，是跟火焰图相似的可视化方式。区别是，x 轴并不是按照字母表的顺序排列，而是按照时间百分比排序。这样做的优势是：很容易发现时间相关的问题。但同时，缺点也很明显，这种排序减少了方法的合并，当分析多个线程时，劣势更加明显。Flame charts 跟火焰图一起使用应该会更有用。

### Challenges

火焰图面临的挑战，更多是来自于 profilers，而不是火焰图本身。profilers 会面对两类典型的问题：
• **Stack traces 不完整**。
有些 profiler 只提供固定深度（比如 10）的采样，这种不完整的输出不利于分析，如果增大深度，这些 profiler 会直接失败。更糟糕的问题是有些编译器会使用“**重用帧指针寄存器**（frame pointer register）”这样的编译优化，破坏了标准 Stack traces 采样流程。解决方式是关闭这种编译器优化（比如 gcc 使用参数 `-fno-omit-frame-pointer`）或者使用另一种采样技术。
• **方法名称丢失**。
有些 profilers，堆栈信息是完整的，但是方法名称却丢失了，显示为十六进制地址。使用了**JIT** (just-in-time) 技术编译的代码经常有这个问题。因为 JIT 并不会为 profiler 创建符号表。对于不同的 profiler 和 runtime，这个问题有不同的解决方式，比如 Linux perf_events 允许应程序提供一个额外的符号表文件，用于产生采样结果。
我们在 Netflix 的工作过程中，曾经尝试为 Java 生成火焰图，结果两个问题都遇到了 ![](https://notes-learning.oss-cn-beijing.aliyuncs.com/hxnw7s/1619535684464-4d4fe494-900b-49a2-8435-567c5510f0f0.svg) 。第一个问题通过一个新的 jvm 参数：`—XX:+PreserveFramePointer`解决了。该参数禁用了编译器优化。第二个问题由 perf-map-agent 解决，这个 Java agent 可以为 profiler 生成符号表。
火焰图面临的另一个挑战是生成 SVG 文件的大小。一个超大的 profile 有可能会有成千上万的 code paths，最终生成的 svg 可能有几十 mb,浏览器加载要花费比较长的时间。解决方式是忽略掉在途中细到肉眼难以观察的方法，忽略这些方法不会影响全局视图，同时能缩小输出。

### Other Color Schemes

除了使用随机的暖色调外，火焰图支持其他配色方案，比如使用颜色区分代码或者数据维度。
Perl 版本的火焰图支持很多配色选项，其中一个选项是**java**。该选项通过不同颜色区分模块，规则如下：绿色代表 Java，黄色为 C++，红色用于所有其他用户代码，橙色用于内核模块。见图七（可交互的 svg 格式的图 7 可以在[http://queue.acm.org/downloads/2016/Gregg7.svg](https://link.zhihu.com/?target=http%3A//queue.acm.org/downloads/2016/Gregg7.svg) 体验）。
![](https://notes-learning.oss-cn-beijing.aliyuncs.com/hxnw7s/1619535684510-3f8c57a8-1f07-4361-936b-1539013e1dfb.jpeg)
另一个可用的选项是**hash**，该选项根据函数名的 hash 选择颜色，这个选项在比较多张火焰图的时候非常有用，因为在不同的图片上，相同的方法会用相同的颜色表示。
颜色选项对于“差分火焰图”也很重要，“差分火焰图”将在下一节中介绍。

### Differential Flame Graphs

差分火焰图显示了两个 profile 的区别。比如现在有两个 profile A 和 B，Perl 版本的火焰图支持这样的操作：以 A 为基准，使用 B 与 A 的差值生成一张火焰图。在差分火焰图上，红色代表差值为正数，蓝色代表差值为负数。这种差分图的问题是，如果 A 中的某个方法，在 B 中完全没有被调用，差分图就会把这种方法丢弃。丢失的数据会误导用户。
一种改进的方法是`flamegraphdiff`，使用三张图解决这个问题。第一张是 A,第二张是 B，第三张是前边提到的差分图。当鼠标 hover 到任意一个方法时，三张图上该方法都会高亮显示。同时 flamegraphdiff 也支持红/蓝的配色方案说明百分比的增减。

### Other Targets

前边提到，火焰图可以可视化多种 profiler 的输出。profiler 有以下几种：CPU PMC (performance monitoring counter) 溢出事件, 静态追踪（static tracing） 事件, 动态追踪（dynamic tracing） 事件。下面是一些其他 profiler 的例子.
**Stall Cycles**
tall-cycle 火焰图显示被处理器或硬件资源（通常是内存 I/O）block 的代码路径。stack trace 使用 PMC profiler, 比如 Linux perf_events 采集，分析这样的火焰图，开发人员可以使用其他策略优化代码，优化的目的是减少内存 I/O，而不是减少指令。
**CPI**
CPI(cycles per instruction)指令每周期数, 或者其倒数 IPC 可以描述 cpu 的使用情况。cpi 火焰图用宽度表示采样 cpu 周期数，同时用颜色区分每个函数的 cpi：红色表示高 cpi，蓝色表示低 cpi。cpi 火焰图需要两个 profile：CPU 采样 profile 和 指令数量 profile，使用差分火焰图技术生成。
**Memory**
火焰图可以通过可视化许多不同的内存事件来揭示内存增长。
通过跟踪 malloc() 方法可以生成 malloc 火焰图，用于可视化申请内存的代码路径。这个方案可能很难应用与实践，因为 malloc 函数调用很频繁，使得在某些场景中跟踪它们的成本很高。
通过跟踪 brk()和 mmap()方法可以展示导致虚拟内存中的扩展的代码路径。当然如果异步的申请内存就另当别论了。这些方法调用频率很低，很适合追踪。
跟踪内存缺页异常可以展示导致物理内存中的扩展的代码路径。导致内存缺页的代码路径通常会频繁的申请内存。内存缺页异也是低频事件。
**I/O**
与 io 相关的问题，比如文件系统，存储设备和网络，通常可以使用 system tracers 方便的追踪。使用这类 profiles 生成的火焰图显示了使用 I/O 的代码路径。
在实践中，io 火焰图可以定位导致 io 的原因。比如一个磁盘 io 可能由以下事件引起：应用程序发起系统调用，文件系统预读，异步的脏数据 flush，或者内核异步的磁盘清理。通过火焰图上的代码路径可以区分以上导致磁盘 io 的事件类型。
**Off-CPU**
当然，也有许多问题使用上边提到的火焰图是看不见的。分析这些问题需要了解线程被阻塞（而不是 on cpu）的时间。线程阻塞的原因有很多，比如等待 I/O、锁、计时器、打开 CPU 以及等待分页或交换。追踪线程被重新调度时的 stack traces 可以区分这些原因，线程 block 的时间长度可以通过跟踪线程从离开 CPU 到返回 CPU 的时间来测量。系统 profilers 通常使用内核中的静态跟踪点来跟踪这些事件。
通过上边 profile 可以生成 Off-CPU 火焰图用来分析这类问题 ![](https://notes-learning.oss-cn-beijing.aliyuncs.com/hxnw7s/1619535684469-02f3f063-04fb-4a64-8b07-c845e66b25f8.svg) 。
**Wakeups**
在 Off-CPU 火焰图的应用中发现这样一个问题：当线程阻塞的原因是条件变量时，火焰图很难解答“为什么条件变量被其他线程持有这么长时间”这样的问题。
通过跟踪线程唤醒事件，可以生成 Wakeups 火焰图。该图显示了线程阻塞的原因。Wakeups 火焰图可以与 Off-CPU 火焰图一起研究，以获得关于阻塞线程的更多信息。
**Chain Graphs**
持有条件变量的线程可能已在另一个条件变量（由另一个线程持有）上被阻塞。实际上，一个线程可能被第二个、第三个甚至第四个线程阻塞。对于这种复杂场景，只分析一个 wakeup 火焰图可能还不够。
chain 火焰图是分析这种复杂场景的一种尝试，chain 火焰图从 off-CPU 火焰图为基础，将 Wakeups 火焰图放到每个 stack traces 顶部。通过自顶向下的分析可以理解阻塞线程的整个条件变量链路。宽度对应线程 block 的时间。
chain 火焰图可以通过组合 Off-CPU 火焰图 U 和 Wakeups 火焰图来实现。这需要大量的采样，目前来看，在实际应用中不太现实。

### Future Work

与火焰图相关的许多工作，都涉及到不同的 profiler 与不同的 runtimes（例如，对于 NoDE.JS、Ruby、Perl、Lua、Erlang、Python、Java、Gangangand，以及 dTrof、PrimeEvices、PMCSTAT、Xperf、仪器）。等等）。将来可能会增加更多种类。
另一个正在开发的差分火焰图，称为白/黑差分，使用前面描述的单火焰图方案加上右侧的一个额外区域，用来显示丢失的代码路径。差分火焰图（任何类型）在未来也会得到更多的应用；在 Netflix，我们正在努力让微服务每晚生成这些图：用来帮助进行性能问题分析。
其他几个火焰图实现正在开发中，探索不同的特性。比如：bla bla bla...bla bla bla...

### Conclusion

火焰图是堆栈跟踪的可视化的高效工具，支持 CPU 以及许多其他 profile。它创建了软件执行情况的可视地图，并允许用户导航到感兴趣的区域。与其他可视化技术相比，火焰图更直观的传递信息，在处理超大 profile 是优势明显。火焰图作为理解 profile 的基本工具，已经成功分析解决了无数的性能问题。

### Acknowledgments

bla bla bla...bla bla bla...

### References(只列举了感兴趣的)

...bla bla bla...bla bla bla...
6\. [Gregg, B., Spier, M. 2015. Java in flames. The Netflix Tech Blog;](https://link.zhihu.com/?target=http%3A//techblog.netflix.com/2015/07/java-in-flames.html)
7\. [Heer, J., Bostock, M., Ogievetsky, V. 2010. A tour through the visualization zoo. acmqueue 8(5);](https://link.zhihu.com/?target=http%3A//queue.acm.org/detail.cfm%3Fid%3D1805128)
...bla bla bla...bla bla bla... 10.[Odds, G. 2013. The science behind data visualisation. Creative Bloq;](https://link.zhihu.com/?target=http%3A//www.creativebloq.com/design/science-behind-data-visualisation-8135496)
...bla bla bla...bla bla bla... 15.[Zhang, Y. 2013. Introduction to off-CPU time flame graphs;](https://link.zhihu.com/?target=http%3A//agentzh.org/misc/slides/off-cpu-flame-graphs.pdf)

## 结论：

现在终于明白为什么很多汉化的技术书籍，读起来比较拗口了。翻译这么短的一篇介绍性文章，都感觉心力憔悴。。。
从最开始的小心翼翼，逐句翻译生怕表达错了含义，结果觉得翻译结果**啰哩啰嗦**。到后来试图用自己的理解描述，又害怕因为理解错误，词不达意。到最后逐渐佛系，降低心里期望。现在终于明白了译者的辛苦。
在这里正好提一下这个 github 地址：[TranslateProject](https://link.zhihu.com/?target=https%3A//github.com/LCTT/TranslateProject)，之前网上查到的好多帖子都是 LCTT 汉化的，一直都没有注意，这次查资料才发现。后续我也试一试按照他们的规范把这篇文章提个 push 申请啥的。
读完这篇论文，才知道网上很多对火焰图的理解都比较片面。特别是后边 Other Targets 部分，很多内容很有意思，网上查找“火焰图”关键词搜到的文章都没有提及。通过这次经历，也终于理解了之前一位前辈说**要多看 paper**，不要老是在网上看几篇帖子就以为自己理解了。果然还是一手的信息靠谱。
整片文章我最感兴趣的，就是**Challenge**里边提到的，作者在给 java 生成火焰图的时候，遇到的两个问题。当时读的时候，一直不理解帧指针寄存器（frame pointer register）是个啥东西。还好这段有 References [java-in-flames](https://link.zhihu.com/?target=http%3A//techblog.netflix.com/2015/07/java-in-flames.html.)。根据指导我详细学习了一下如何使用本文提到的技术，生成 java 的火焰图。后续会整理一个“《性能之巅》学习笔记之火焰图 其之二”介绍一下。总之也是挺有趣的。_当然，bcc 是什么鬼？我不知道啊，听都没听过\[#doge#]～～_
之前一直很奇怪，为什么论文必须要有 References，现在有点理解了。
这个技能点到底有没有点上，其实心里还是没什么底，希望能够坚持下去，虽然阅读英文的资料耗费心力，但是收获真的很大。就像长跑者思维：因为今天下雨了，所以才要去跑步。因为很困难，所以干就完了！
