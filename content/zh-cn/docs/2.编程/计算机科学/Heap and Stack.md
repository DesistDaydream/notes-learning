---
title: "Heap and Stack"
linkTitle: "Heap and Stack"
date: "2023-08-25T10:27"
weight: 20
---

# 概述

> 参考：
>
> - [StackOverflow](https://stackoverflow.com/questions/79923/what-and-where-are-the-stack-and-heap)

**Heap(堆) 与 Stack(栈)**

[汇编语言](/docs/2.编程/低级编程语言/汇编语言.md) 这里有一些 堆栈 的通俗解释，非常好理解

## 内存模型：Heap

寄存器只能存放很少量的数据，大多数时候，CPU 要指挥寄存器，直接跟内存交换数据。所以，除了寄存器，还必须了解内存怎么储存数据。

程序运行的时候，操作系统会给它分配一段内存，用来储存程序和运行产生的数据。这段内存有起始地址和结束地址，比如从`0x1000`到`0x8000`，起始地址是较小的那个地址，结束地址是较大的那个地址。

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/f3d22e4e-0b4d-4787-9050-29da7e23160e/bg2018012208.png)

程序运行过程中，对于动态的内存占用请求（比如新建对象，或者使用`malloc`命令），系统就会从预先分配好的那段内存之中，划出一部分给用户，具体规则是从起始地址开始划分（实际上，起始地址会有一段静态数据，这里忽略）。举例来说，用户要求得到 10 个字节内存，那么从起始地址`0x1000`开始给他分配，一直分配到地址`0x100A`，如果再要求得到 22 个字节，那么就分配到`0x1020`。

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/f3d22e4e-0b4d-4787-9050-29da7e23160e/bg2018012209.png)

这种因为用户主动请求而划分出来的内存区域，叫做 Heap（堆）。它由起始地址开始，从低位（地址）向高位（地址）增长。Heap 的一个重要特点就是不会自动消失，必须手动释放，或者由垃圾回收机制来回收。

## 内存模型：Stack

除了 Heap 以外，其他的内存占用叫做 Stack（栈）。简单说，Stack 是由于函数运行而临时占用的内存区域。

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/f3d22e4e-0b4d-4787-9050-29da7e23160e/bg2018012210.png)

请看下面的例子。

     int main() {
       int a = 2;
       int b = 3;
    }

上面代码中，系统开始执行`main`函数时，会为它在内存里面建立一个帧（frame），所有`main`的内部变量（比如`a`和`b`）都保存在这个帧里面。`main`函数执行结束后，该帧就会被回收，释放所有的内部变量，不再占用空间。

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/f3d22e4e-0b4d-4787-9050-29da7e23160e/bg2018012211.png)

如果函数内部调用了其他函数，会发生什么情况？

     int main() {
       int a = 2;
       int b = 3;
       return add_a_and_b(a, b);
    }

上面代码中，`main`函数内部调用了`add_a_and_b`函数。执行到这一行的时候，系统也会为`add_a_and_b`新建一个帧，用来储存它的内部变量。也就是说，此时同时存在两个帧：`main`和`add_a_and_b`。一般来说，调用栈有多少层，就有多少帧。

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/f3d22e4e-0b4d-4787-9050-29da7e23160e/bg2018012212.png)

等到`add_a_and_b`运行结束，它的帧就会被回收，系统会回到函数`main`刚才中断执行的地方，继续往下执行。通过这种机制，就实现了函数的层层调用，并且每一层都能使用自己的本地变量。

所有的帧都存放在 Stack，由于帧是一层层叠加的，所以 Stack 叫做栈。生成新的帧，叫做 "入栈"，英文是 push；栈的回收叫做 "出栈"，英文是 pop。Stack 的特点就是，最晚入栈的帧最早出栈（因为最内层的函数调用，最先结束运行），这就叫做 "后进先出" 的数据结构。每一次函数执行结束，就自动释放一个帧，所有函数执行结束，整个 Stack 就都释放了。

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/f3d22e4e-0b4d-4787-9050-29da7e23160e/bg2018012213.jpg)

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/f3d22e4e-0b4d-4787-9050-29da7e23160e/bg2018012214.jpg)

Stack 是由内存区域的结束地址开始，从高位（地址）向低位（地址）分配。比如，内存区域的结束地址是`0x8000`，第一帧假定是 16 字节，那么下一次分配的地址就会从`0x7FF0`开始；第二帧假定需要 64 字节，那么地址就会移动到`0x7FB0`。

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/f3d22e4e-0b4d-4787-9050-29da7e23160e/bg2018012215.png)
