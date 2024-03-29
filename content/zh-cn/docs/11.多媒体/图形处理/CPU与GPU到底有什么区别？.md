---
title: CPU与GPU到底有什么区别？
---

原文链接：[公众号-码农的荒岛求生，CPU 与 GPU 到底有什么区别？](https://mp.weixin.qq.com/s/Vokg9qdQaWt3pPtWeHPkew)

大家好，我是小风哥，今天简单聊聊 CPU 与 GPU。

CPU 的故事我们聊得比较多了，之前也发布过很多关于 CPU 的文章，因此这里重点聊聊 GPU。

**教授 vs 小学生**

你可以简单的将 CPU 理解为学识渊博的教授，什么都精通，而 GPU 则是一堆小学生，只会简单的算数运算，可即使教授再神通广大，也不能一秒钟内计算出 500 次加减法，**因此对简单重复的计算来说单单一个教授敌不过数量众多的小学生**，在进行简单的算数运算这件事上，500 个小学生(并发)可以轻而易举打败教授。

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/ada20417-f607-4763-9bff-f503a64ca4e5/640)

因此我们可以看到，CPU 和 GPU 的最大不同在于架构，CPU 适用于广泛的应用场景(学识渊博)，可以执行任意程序，而 GPU 则专为多任务而生，并发能力强，具体来讲就是多核，一般的 CPU 有 2 核、4 核、8 核等，而 GPU 则可能会有成百上千核：

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/ada20417-f607-4763-9bff-f503a64ca4e5/640)

可以看到，CPU 内部 cache 以及控制部分占据了很大一部分片上面积，因此计算单元占比很少，再来看看 GPU，GPU 只有很简单的控制单元，剩下的大部分都被计算单元占据，因此 CPU 的核数有限，而 GPU 则轻松堆出上千核：

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/ada20417-f607-4763-9bff-f503a64ca4e5/640)

只不过 CPU 中每个核的能力好比教授，而 GPU 的每个核的能力好比一个小学生。

你可能会想，为什么 GPU 需要这么奇怪的架构呢？

#####

**为什么 GPU 需要这么多核心？**

想一想计算机上的一张图是怎么表示的？无非就是屏幕上的一个个像素：

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/ada20417-f607-4763-9bff-f503a64ca4e5/640)

我们需要为每个像素进行计算，**而且是相同的运算**，就好比刚才例子中的小学生计算计加法一样，注意，对于屏幕来说一般会有上百万个像素，如果我们要串行的为每一个像素进行运算效率就太低了，因此我们可以让 GPU 中的每一个核心去为相应的像素进行计算，由于 GPU 中有很多核心，因此并行计算可以大幅提高速度。

现在你应该明白为什么 GPU 要这样工作了吧。

除了 GPU 的核心数比较多之外，GPU 的工作方式也比较奇怪。

#####

**奇怪的工作方式**

对 CPU 来说，不同的核心可以执行不同的机器指令，coreA 在运行 word 线程的同时 coreB 上可以运行浏览器线程，这就是所谓的多指令多数据，MIMD，(Multiple Instruction, Multiple Data)。

而 GPU 则不同，**GPU 上的这些核心必须整齐划一的运行相同的机器指令**，只是可以操作不同的数据，这就好比这些小学生在某个时刻必须都进行加法计算，不同的地方在于有的小学生可能需要计算 1+1，有的要计算 2+6 等等，变化的地方仅在于操作数，这就是所谓的单指令多数据，SIMD，(Single Instruction, Multiple Data)。

因此我们可以看到 GPU 的工作方式和 CPU 是截然不同的。

除了这种工作方式之外，GPU 的指令集还非常简单，不像 CPU 这种复杂的处理器，如果你去看 CPU 的编程手册就会发现，CPU 负责的事情非常多：中断处理、内存管理、IO 等等，这些对于 GPU 来说都是不存在的，可以看到 GPU 的定位非常简单，就是纯计算，GPU 绝不是用来取代 CPU 的，CPU 只是把一些 GPU 非常擅长的事情交给它，GPU 仅仅是用来分担 CPU 工作的配角。

CPU 和 GPU 是这样配合工作的：

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/ada20417-f607-4763-9bff-f503a64ca4e5/640)

#####

**GPU 擅长什么**

比较适合 GPU 的计算场景是这样的：1)计算简单；2）重复计算，因此如果你的计算场景和这里的图像渲染相似那么使用 GPU 就很合理了。

因此对于图形图像计算、天气预报以及神经网络等都适用于 GPU，哦对了，GPU 还适合用来挖矿。
