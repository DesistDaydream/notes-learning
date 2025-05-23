---
title: Process
linkTitle: Process
weight: 1
---

# 概述

> 参考：
>
> - [Wiki, Process_(computing)](https://en.wikipedia.org/wiki/Process_(computing))
> - [原文连接](https://blog.csdn.net/ljianhui/article/details/46718835)，本文为 IBM RedBook 的 [Linux Performanceand Tuning Guidelines](http://users.polytech.unice.fr/~bilavarn/fichier/elec5_linux/linux_perf_and_tuning_IBM.pdf) 的 1.1 节的翻译
> - [阿里技术，CPU 飙高，系统性能问题如何排查？](https://mp.weixin.qq.com/s/fzLcAkYwKhj-9hgoVkTzaw)

进程管理是操作系统的最重要的功能之一。有效率的进程管理能保证一个程序平稳而高效地运行。它包括进程调度、中断处理、信号、进程优先级、上下文切换、进程状态、进度内存等。

**Process(进程)** 实际是运行在 CPU 中的一个 **[Program](/docs/2.编程/Program.md)(程序) 的实体**。在 Linux 系统中，能够同时运行多个进程。

Program(程序) 和 Process(进程) 的区别是什么呢?

- 在很久很久以前，计算机刚出现的时候，是没有操作系统的，那时候一台机器只是运行一个程序，计算后得出数据，后来人们为了同时运行多个程序从而研究出了操作系统，在操作系统之上可以运行多个程序
- 进程是程序的一个具体实现。类似于按照食谱，真正去做菜的过程。同一个程序可以执行多次，每次都可以在内存中开辟独立的空间来装载，从而产生多个进程。不同的进程还可以拥有各自独立的 IO 接口。

<font color="#ff0000">举例说明</font>：

比如:

```bash
root         839       1  0 Mar07 ?        Ssl   28:50 /usr/bin/dockerd -H fd:// --containerd=/run/containerd/containerd.sock
```

这就是一个 **Processs(进程)**，包括其 ID、启动时间、等等信息的集合体。进程的唯一标识符就是 ID，而启动该进程的程序是 dockerd

- 至于 `/usr/bin/dockerd -H fd:// --containerd=/run/containerd/containerd.sock` 这一串则是启动进程的 **Command(命令)**
- 其中命令中的 dockerd 就是启动该进程的 **Program(程序)**，`/usr/bin/` 是程序所在路径，后面的 `-H fd:// --containerd=/run/containerd/containerd.sock` 是程序的参数。
- 这一整串字符串所组成的命令，就是启动进程的必备条件，操作系统当收到命令后，会被进程分配 ID，并记录下各种状态信息。

Linux 通过在短的时间间隔内轮流运行这些进程而实现“多任务”。这一短的时间间隔称为“时间片”，让进程轮流运行的方法称为“进程调度” ，完成调度的程序称为调度程序。

进程调度控制进程对 CPU 的访问。当需要选择下一个进程运行时，由调度程序选择最值得运行的进程。可运行进程实际上是仅等待 CPU 资源的进程，如果某个进程在等待其它资源，则该进程是不可运行进程。Linux 使用了比较简单的基于优先级的进程调度算法选择新的进程。

通过多任务机制，每个进程可认为只有自己独占计算机，从而简化程序的编写。每个进程有自己单独的地址空间，并且只能由这一进程访问，这样，操作系统避免了进程之间的互相干扰以及“坏”程序对系统可能造成的危害。 为了完成某特定任务，有时需要综合两个程序的功能，例如一个程序输出文本，而另一个程序对文本进行排序。为此，操作系统还提供进程间的通讯机制来帮助完成这样的任务。Linux 中常见的进程间通讯机制有信号、管道、共享内存、信号量和套接字等。

内核通过 SCI 提供了一个 API 来创建一个新进程(fork、exec 或 Portable Operating System Interface \[POSⅨ] 函数)、停止进程(kill、exit)、并在它们之间进行通信和同步(signal 或者 POSⅨ 机制)。

计算机实际上可以做的事情实质上非常简单，比如计算两个数的和，再比如在内存中寻找到某个地址等等。这些最基础的计算机动作被称为指令(instruction)。所谓的程序(program)，就是这样一系列指令的所构成的集合。通过程序，我们可以让计算机完成复杂的操作。程序大多数时候被存储为可执行的文件。这样一个可执行文件就像是一个菜谱，计算机可以按照菜谱作出可口的饭菜。

操作系统的一个重要功能就是为进程提供方便，比如说为进程分配内存空间，管理进程的相关信息等等，就好像是为我们准备好了一个精美的厨房。

## 进程的生命周期

每一个进程都有其生命周期，例如创建、运行、终止和消除。这些阶段会在系统启动和运行中重复无数次。因此，进程的生命周期对于其性能的分析是非常重要的。下图展示了经典的进程生命周期。

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/ld23ik/1616167507353-2f676d82-88da-483c-a939-399f284d6425.jpeg)

不会关闭的常驻进程可以称为 **Daemon Process(守护进程，简称 Daemon)**

> 一般 daemon 的名称都会在进程名后加一个字母 d 作为 daemon 的 process，比如 vsftp 的 daemon 就是 vsftpd。

当一个进程创建一个新的进程，创建进程(父进程)的进程调用 一个 fork() 系统调用。当 fork() 系统调用被调用，它得到该新创建进程（子进程）的进程描述并调用一个新的进程 id。它复制该值到父进程进程描述到子进程中。此时整个的父进程的地址空间是没有被复制的；父子进程共享相同的地址空间。

exec() 系统调用复制新的程序到子进程的地址空间。因为父子进程共享地址空间，写入一个新的程序的数据会引起一个分页错误。在这种情况下，内存会分配新的物理内存页给子进程。

这个推迟的操作叫作写时复制。子进程通常运行他们自己的程序而不是与父进程运行相同的程序。这个操作避免了不必要的开销，因为复制整个地址空间是一个非常缓慢和效率低下的操作，它需要使用大量的处理器时间和资源。

当程序已经执行完成，子进程通过调用 exit()系统调用终止。exit()系统调用释放进程大部分的数据并通过发送一个信号通知其父进程。此时，子进程是一个被叫作僵尸进程的进程（参阅 page 7 的“Zombie processes”）。

子进程不会被完全移除直到其父进程知道其子进程的调用 wait()系统调用而终止。当父进程被通知子进程终止，它移除子进程的所有数据结构并释放它的进程描述。

## 父进程与子进程

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/ld23ik/1616167507409-d531245f-abbe-4a2a-b575-d2ae72c6949f.jpeg)

- 用颜色的线标示的两列，左侧的为进程号(PID)右侧的为父进程号(PPID)
- 子进程与父进程的环境变量相同
- 老进程成为新进程的父进程(parent process)，而相应的，新进程就是老的进程的子进程(child process)。一个进程除了有一个 PID 之外，还会有一个 PPID(parent PID)来存储的父进程 PID。如果我们循着 PPID 不断向上追溯的话，总会发现其源头是 init 进程。所以说，所有的进程也构成一个以 init 为根的树状结构。
- 如上图所示，我们查询当前 shell 下的进程：
  - 我们可以看到，第二个进程 ps 是第一个进程 bash 的子进程。
  - 还可以用 `pstree` 命令来显示整个进程树。
- fork() 通常作为一个函数被调用。这个函数会有两次返回，将子进程的 PID 返回给父进程，0 返回给子进程。实际上，子进程总可以查询自己的 PPID 来知道自己的父进程是谁，这样，一对父进程和子进程就可以随时查询对方。
- 通常在调用 fork 函数之后，程序会设计一个 if 选择结构。当 PID 等于 0 时，说明该进程为子进程，那么让它执行某些指令,比如说使用 exec 库函数(library function)读取另一个程序文件，并在当前的进程空间执行 (这实际上是我们使用 fork 的一大目的: 为某一程序创建进程)；而当 PID 为一个正整数时，说明为父进程，则执行另外一些指令。由此，就可以在子进程建立之后，让它执行与父进程不同的功能。

### 子进程的 termination(终结)

当子进程终结时，它会通知父进程，并清空自己所占据的内存，并在内核里留下自己的退出信息(exit code，如果顺利运行，为 0；如果有错误或异常状况，为>0 的整数)。在这个信息里，会解释该进程为什么退出。父进程在得知子进程终结时，有责任对该子进程使用 wait 系统调用。这个 wait 函数能从内核中取出子进程的退出信息，并清空该信息在内核中所占据的空间。但是，如果父进程早于子进程终结，子进程就会成为一个孤儿(orphand)进程。孤儿进程会被过继给 init 进程，init 进程也就成了该进程的父进程。init 进程负责该子进程终结时调用 wait 函数。

当然，一个糟糕的程序也完全可能造成子进程的退出信息滞留在内核中的状况（父进程不对子进程调用 wait 函数），这样的情况下，子进程成为僵尸(zombie)进程。当大量僵尸进程积累时，内存空间会被挤占。

## Thread(线程)

一个线程是一个单独的进程生成的一个执行单元。它与其他的线程并行地运行在同一个进程中。各个线程可以共享进程的资源，例如内存、地址空间、打开的文件等等。它们能访问相同的程序数据集。线程也被叫作轻量级的进程（Light Weight Process，LWP）。因为它们共享资源，所以每个线程不应该在同一时间改变它们共享的资源。互斥的实现、锁、序列化等是用户程序的责任。

从性能的角度来说，创建线程的开销比创建进程少，因数创建一个线程时不需要复制资源。另一方面，进程和线程拥在调度算法上有相似的特性。**内核以相似的方式处理它们**。

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/ld23ik/1616167507380-b6ae3b1e-b47c-454c-b3c7-9942dde4f480.jpeg)

所以，一个进程创建的线程，也是可以运行在多个 CPU 上的。

![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/ld23ik/1616645843002-c07df4a7-3d7a-4969-8203-4bc20169721a.png)

在现在的 Linux 实现中，线程支持 UNIX 的可移植操作系统接口（POSIX）标准库。在 Linux 操作系统中有几种可用的线程实现。以下是广泛使用的线程库：

Linux Threads 自从 Linux 内核 2.0 起就已经被作为默认的线程实现。Linux Threads 的一些实现并不符合 POSIX 标准。Native POSIX Thread Library（NPTL）正在取代 Linux Threads。Linux Threads 在将来的 Linux 企业发行版中将不被支持。

Native POSIX Thread Libary（NPTL）

NPTL 最初是由红帽公司开发的。NPTL 与 POSIX 更加兼容。通过 Linux 内核 2.6 的高级特性，例如，新的 clone()系统调用、信号处理的实现等等，它具有比 LinuxThreads 更高的性能和伸缩性。

NPTL 与 LinuxThreads 有一些不兼容。一个依赖于 LinuxThreads 的应用可能不能在 NPTL 实现中工作。

Next Generation POSIX Thread（NGPT）

NGPT 是一个 IBM 开发的 POSIX 线程库。现在处于维护阶段并且在未来也没有开发计划。

使用 LD_ASSUME_KERNEL 环境变量，你可以选择在应用中使用哪一个线程库。

## Linux 内核代码中的 Process

在 Linux 中，**Process(进程) 属于** **Task(任务)** 的一种类型，都被 task_struct 结构管理，该结构同时被叫作进程描述。一个进程描述包含一个运行进程所有的必要信息，例如进程标识、进程属性和构建进程的资源。如果你了解该进程构造，你就能理解对于进程的运行和性能来说，什么是重要的。

v5.14 代码：[include/linux/sched.h](https://github.com/torvalds/linux/blob/v5.14/include/linux/sched.h#L661)

```c
struct task_struct {
#ifdef CONFIG_THREAD_INFO_IN_TASK
 /*
  * For reasons of header soup (see current_thread_info()), this
  * must be the first element of task_struct.
  */
 struct thread_info  thread_info;
#endif
    ......
 // 进程状态
    unsigned int   __state;
    // 进程唯一标识符
 pid_t    pid;
 pid_t    tgid;
 // 进程名称，上限 16 字符
    char    comm[TASK_COMM_LEN];
    // 打开的文件
 struct files_struct  *files;
 ......
}
```

下图展示了进程结构相关的进程信息概述。

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/ld23ik/1616167507336-aaeec645-b9df-41c3-99ab-6bf39aed4f42.jpeg)

其实从这里能看出来，从某种角度来看，**对于内核来说并没有线程这个概念。Linux 把所有的线程都当做进程来实现，内核也没有特别的调度算法来处理线程。**线程仅仅被视为一个与其他进程共享某些资源的进程，和进程一样，每个线程也都是有自己的 `task_struct`，所以在内核中，线程看起来就是一个普通的进程。线程也被称作轻量级进程，一个进程可以有多个线程，线程拥有自己独立的栈，切换也由操作系统调度。在 Linux 上可以通过 `pthread_create()` 方法或者 `clone()` 系统调用创建；

# 进程优先级和 nice 值

进程优先级是一个数值，它通过动态的优先级和静态的优先级来决定进程被 CPU 处理的顺序。一个拥有更高进程优先级的进程拥有更大的机率得到处理器的处理。

内核根据进程的行为和特性使用试探算法，动态地调整调高或调低动态优先级。一个用户进程可以通过使用进程的 nice 值间接改变静态优先级。一个拥有更高静态优先级的进程将会拥有更长的时间片（进程能在处理上运行多长时间）。

Linux 支持从 19（最低优先级）到-20（最高优先级）的 nice 值。默认值为 0。把程序的 nice 值修改为负数（使进程的优先级更高），需要以 root 身份登陆或使用 su 命令以 root 身份执行。

# 上下文切换

在进程运行过程中，进程的运行信息被保存于处理器的寄存器和它的缓存中。正在执行的进程加载到寄存器中的数据集被称为上下文。为了切换进程，运行中进程的上下文将会被保存，接下来的运行进程的上下文将被被恢复到寄存器中。进程描述和内核模式堆栈的区域将会用来保存上下文。这个切换被称为上下文切换。过多的上下文切换是不受欢迎的，因为处理器每次都必须清空刷新寄存器和缓存，为新的进程制造空间。它可能会引起性能问题。

下图说明了上下文切换如何工作。

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/ld23ik/1616167507475-6f5a9385-f033-4c00-8344-2953197b973c.jpeg)

# 中断处理

中断处理是优先级最高的任务之一。中断通常由 I/O 设备产生，例如网络接口卡、键盘、磁盘控制器、串行适配器等等。中断处理器通过一个事件通知内核（例如，键盘输入、以太网帧到达等等）。它让内核中断进程的执行，并尽可能快地执行中断处理，因为一些设备需要快速的响应。它是系统稳定的关键。当一个中断信号到达内核，内核必须切换当前的进程到一个新的中断处理进程。这意味着中断引起了上下文切换，因此大量的中断将会引起性能的下降。

在 Linux 的实现中，有两种类型的中断。硬中断是由请求响应的设备发出的（磁盘 I/O 中断、网络适配器中断、键盘中断、鼠标中断）。软中断被用于处理可以延迟的任务（TCP/IP 操作，SCSI 协议操作等等）。你可以在 `/proc/interrupts` 文件中查看硬中断的相关信息。

在多处理器的环境中，中断被每一个处理器处理。绑定中断到单个的物理处理中能提高系统的性能。更多的细节，请参阅 4.4.2，“CPU 的中断处理亲和力”。

# 进程的状态

每一个进程拥有自己的状态，状态表示了进程当前在发生什么。LINUX 2.6 以后的内核中，在进程的执行期间进程的状态会发生改变，进程一般存在 7 种基础状态：D-不可中断睡眠、R-可执行、S-可中断睡眠、T-暂停态、t-跟踪态、X-死亡态、Z-僵尸态，这几种状态在 ps 命令的 man 手册中有对应解释。

- **D**＃不间断的睡眠（通常是 IO）
- **R** ＃正在运行或可运行（在运行队列上）
- **S** ＃可中断的睡眠（等待事件完成）
- **T** ＃被作业控制信号停止
- **t**＃在跟踪过程中被调试器停止
- **X** ＃已死（永远都不会出现）
- **Z** ＃已终止运行（“僵尸”）的进程，已终止但未由其父进程获得

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/ld23ik/1616167507456-ca89ed8d-d8a1-4cd6-96ab-c78372840f4a.jpeg)

## D (TASK_UNINTERRUPTIBLE)，不可中断睡眠态

顾名思义，位于这种状态的进程处于睡眠中，并且不允许被其他进程或中断(异步信号)打断。因此这种状态的进程，是无法使用 kill -9 杀死的(kill 也是一种信号)，除非重启系统(没错，就是这么头硬)。不过这种状态一般由 I/O 等待(比如磁盘 I/O、网络 I/O、外设 I/O 等)引起，出现时间非常短暂，大多很难被 PS 或者 TOP 命令捕获(除非 I/O HANG 死)。SLEEP 态进程不会占用任何 CPU 资源。

## R (TASK_RUNNING)，可执行态

这种状态的进程都位于 CPU 的可执行队列中，正在运行或者正在等待运行，即不是在上班就是在上班的路上。

在此状态下，表示进程正在 CPU 中运行或在队列中等待运行（运行队列）。

## S (TASK_INTERRUPTIBLE)，可中断睡眠态

不同于 D，这种状态的进程虽然也处于睡眠中，但是是允许被中断的。这种进程一般在等待某事件的发生（比如 socket 连接、信号量等），而被挂起。一旦这些时间完成，进程将被唤醒转为 R 态。如果不在高负载时期，系统中大部分进程都处于 S 态。SLEEP 态进程不会占用任何 CPU 资源。

在此状态下，进程被暂停并等待一个某些条件状态的到达。如果一个进程处于 TASK_INTERRUPTIBLE 状态并接收到一个停止的信号，进程的状态将会被改变并中断操作。一个典型的 TASK_INTERRUPTIBLE 状态的进程的例子是一个进程等待键盘中断。

## T & t (TASK_STOPPED & TASK_TRACED)，暂停 or 跟踪态

这种两种状态的进程都处于运行停止的状态。不同之处是暂停态一般由于收到 SIGSTOP、SIGTSTP、SIGTTIN、SIGTTOUT 四种信号被停止，而跟踪态是由于进程被另一个进程跟踪引起(比如 gdb 断点）。暂停态进程会释放所有占用资源。

TASK_STOPPED 在此状态下的进程被某些信号（如 SIGINT，SIGSTOP）暂停。进程正在等待通过一个信号恢复运行，例如 SIGCONT。

## Z (EXIT_ZOMBIE/TASK_ZOMBIE), 僵尸态

这种状态的进程实际上已经结束了，但是父进程还没有回收它的资源（比如进程的描述符、PID 等）。僵尸态进程会释放除进程入口之外的所有资源。

当一个进程调用 exit()系统调用退出后，它的父进程应该知道该进程的终止。处于 TASK_ZOMBIE 状态的进程会等待其父进程通知其释放所有的数据结构。

当一个进程接收到一个信号而终止，它在结束自己之前，通常需要一些时间来结束所有的任务（例如关闭打开的文件）。在这个通常非常短暂的时间内，该进程就是一个僵尸进程。

进程已经完成所有的关闭任务后，它会向父进程报告其即将终止。有些时候，一个僵尸进程不能把自己终止，这将会引导它的状态显示为 z（zombie）。

使用 kill 命令来关闭这样的一个进程是不可能的，因为该进程已经被认为已经死掉了。如果你不能清除僵尸进程，你可以结束其父进程，然后僵尸进程也随之消失。但是，如果父进程为 init 进程，你不能结束它。init 进程是一个非常重要的进程，因此可能需要重启系统来清除僵尸进程。

## X (EXIT_DEAD), 死亡态

进程的真正结束态，这种状态一般在正常系统中捕获不到。

# 进程内存段

进程使用其自身的内存区域来执行工作。工作的变化根据情况和进程的使用而决定。进程可以拥有不同的工作量特性和不同的数据大小需求。进程必须处理各种数据大小。为了满足需求，Linux 内核为每个进程使用动态申请内存的机制。进程内存分配的数据结构如图 1-7 所示。

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/ld23ik/1616167507458-2bbc9553-910c-4d66-9ad1-8f45893277da.jpeg)

图 1-7 进程地址空间

进程内存区由以下几部分组成：

Text 段

该区域用于存储运行代码。

Data 段

数据段包括三个区域。

– Data：该区域存储已被初始化的数据，如静态变量。

– BSS：该区域存储初始化为 0 的数据。数据被初始化为 0。

– Heap：该区域用于根据需求使用 malloc()动态申请的内存。堆向高地址方向增长。

Stack 段

该区域用于存储局部变量、函数参数和返回函数的地址。栈向低地址方向增长。

用户进程的地址空间内存分布可以使用 pmap 命令来查看。你可以使用 ps 命令来查看内存段的大小。可以参阅 2.3.10 的“pmap”，“ps 和 pstree”。

# 进程的 exit code(退出码)

在 Linux 系统中，程序可以在执行终止后传递值给其父进程，这个值被称为 **exit code(退出码)** 或 **exit status(退出状态)**或 **reture status(返回码)**。在 POSIX 系统中，惯例做法是当程序成功执行时 **exit code 为 0**，当程序执行失败时 **exit code 非 0**。

传递状态码为何重要？如果你在命令行脚本上下文中查看状态码，答案显而易见。任何有用的脚本，它将不可避免地要么被其他脚本所使用，要么被 bash 单行脚本包裹所使用。特别是脚本被用来与自动化工具 SaltStack 或者监测工具 Nagios 配合使用。这些工具会执行脚本并检查它的状态，来确定脚本是否执行成功。

其中最重要的原因是，即使你不定义状态码，它仍然存在于你的脚本中。如果你不定义恰当的退出码，执行失败的脚本可能会返回成功的状态，这样会导致问题，问题大小取决于你的脚本做了什么。

Linux 提供了一个专门的变量$?来保存上个已执行命令的退出状态码。

对于需要进行检查的命令，必须在其运行完毕后立刻查看或使用$?变量，它的值会变成由 shell 所执行的最后一条命令的退出状态码。

一个成功结束的命令的退出状态码是 0，如果一个命令结束时有错误，退出状态码就是一个正数值（1-255）。

Linux 上执行 exit 可使 shell 以指定的状态值退出。若不设置状态值参数，则 shell 以预设值退出。状态值 0 代表执行成功，其他值代表执行失败。exit 也可用在 script，离开正在执行的 script，回到 shell。

Linux 错误退出状态码没有什么标准可循，但有一些可用的参考。

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/ld23ik/1616167507500-9f1aab01-171b-4ece-a6fa-9f576852a403.webp)

关于具体的服务，相应的退出码，由开发者代码决定。

**Linux 进程退出码**

<https://jin-yang.github.io/post/linux-process-exit-code-introduce.html>

**Linux 退出状态码及 exit 命令**

<https://www.cnblogs.com/01-single/p/7206664.html>

**理解 Exit Code 并学会如何在 Bash 脚本中使用**

<http://blog.jayxhj.com/2016/02/understanding-exit-codes-and-how-to-use-them-in-bash-scripts>

**Appendix E. Exit Codes With Special Meanings**

<http://www.tldp.org/LDP/abs/html/exitcodes.html>

**What is the authoritative list of Docker Run exit codes?**

<https://stackoverflow.com/questions/31297616/what-is-the-authoritative-list-of-docker-run-exit-codes>

**Identifying Exit Codes and their meanings**

<https://support.circleci.com/hc/en-us/articles/360002341673-Identifying-Exit-Codes-and-their-meanings>

**OpenShift Exit Status Codes**

<https://access.redhat.com/documentation/en-US/OpenShift_Online/2.0/html/Cartridge_Specification_Guide/Exit_Status_Codes.html>
