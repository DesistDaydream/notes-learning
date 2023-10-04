---
title: Signal(信号)
---

# 概述

> 参考：
>
> - [Wiki，Signal](https://en.wikipedia.org/wiki/Signal)

**Signal(信号)** 是  [Inter Process Communication(进程间通信)](/docs/1.操作系统/2.Kernel/7.Process/Inter%20Process%20Communication(进程间通信)/Inter%20Process%20Communication(进程间通信).md) 的一种受限形式。信号是发送到进程或同一进程内的特定线程的异步通知，目的是将发生的事件通知给它。发送信号后，操作系统会中断目标进程的正常执行流程以传递信号。在任何非原子指令中，执行都可以中断。如果该进程先前已注册了**信号处理程序**，则将执行该例程。否则，将执行默认信号处理程序。

信号类似于中断，区别在于中断由处理器介导并由内核处理，而信号由内核介导（可能通过系统调用）并由进程处理。内核可能会将中断作为信号传递给引起中断的进程（典型示例为 SIGSEGV，SIGBUS，SIGILL 和 SIGFPE）。

信号类型

Linux 系统共定义了 64 种信号，分为两大类：可靠信号与不可靠信号，前 32 种信号为不可靠信号，后 32 种为可靠信号。

1.1 概念

- 不可靠信号： 也称为非实时信号，不支持排队，信号可能会丢失, 比如发送多次相同的信号, 进程只能收到一次. 信号值取值区间为 1~31；

- 可靠信号： 也称为实时信号，支持排队, 信号不会丢失, 发多少次, 就可以收到多少次. 信号值取值区间为 32~64

  1.2 信号表

在终端，可通过 kill -l 查看所有的 signal 信号

> 使用时，这些信号开头的 3 个大写字符(SIG)可以省略

    [root@master-1 libexec]# kill -l
     1) SIGHUP  2) SIGINT  3) SIGQUIT  4) SIGILL  5) SIGTRAP
     6) SIGABRT  7) SIGBUS  8) SIGFPE  9) SIGKILL 10) SIGUSR1
    11) SIGSEGV 12) SIGUSR2 13) SIGPIPE 14) SIGALRM 15) SIGTERM
    16) SIGSTKFLT 17) SIGCHLD 18) SIGCONT 19) SIGSTOP 20) SIGTSTP
    21) SIGTTIN 22) SIGTTOU 23) SIGURG 24) SIGXCPU 25) SIGXFSZ
    26) SIGVTALRM 27) SIGPROF 28) SIGWINCH 29) SIGIO 30) SIGPWR
    31) SIGSYS 34) SIGRTMIN 35) SIGRTMIN+1 36) SIGRTMIN+2 37) SIGRTMIN+3
    38) SIGRTMIN+4 39) SIGRTMIN+5 40) SIGRTMIN+6 41) SIGRTMIN+7 42) SIGRTMIN+8
    43) SIGRTMIN+9 44) SIGRTMIN+10 45) SIGRTMIN+11 46) SIGRTMIN+12 47) SIGRTMIN+13
    48) SIGRTMIN+14 49) SIGRTMIN+15 50) SIGRTMAX-14 51) SIGRTMAX-13 52) SIGRTMAX-12
    53) SIGRTMAX-11 54) SIGRTMAX-10 55) SIGRTMAX-9 56) SIGRTMAX-8 57) SIGRTMAX-7
    58) SIGRTMAX-6 59) SIGRTMAX-5 60) SIGRTMAX-4 61) SIGRTMAX-3 62) SIGRTMAX-2
    63) SIGRTMAX-1 64) SIGRTMAX

| 取值 | 名称      | 解释                             | 默认动作                                                                                 |
| ---- | --------- | -------------------------------- | ---------------------------------------------------------------------------------------- |
| 1    | SIGHUP    | 挂起                             |                                                                                          |
| 2    | SIGINT    | 中断                             |                                                                                          |
| 3    | SIGQUIT   | 退出                             |                                                                                          |
| 4    | SIGILL    | 非法指令                         |                                                                                          |
| 5    | SIGTRAP   | 断点或陷阱指令                   |                                                                                          |
| 6    | SIGABRT   | abort 发出的信号                 |                                                                                          |
| 7    | SIGBUS    | 非法内存访问                     |                                                                                          |
| 8    | SIGFPE    | 浮点异常                         |                                                                                          |
| 9    | SIGKILL   | kill 信号                        | 不能被忽略、处理和阻塞                                                                   |
| 10   | SIGUSR1   | 用户信号 1                       | 程序自定义的信号，常用这种信号来处理日志或加载配置文件。比如 docker 用这种信号来生成日志 |
| 11   | SIGSEGV   | 无效内存访问                     |                                                                                          |
| 12   | SIGUSR2   | 用户信号 2                       |                                                                                          |
| 13   | SIGPIPE   | 管道破损，没有读端的管道写数据   |                                                                                          |
| 14   | SIGALRM   | alarm 发出的信号                 |                                                                                          |
| 15   | SIGTERM   | 终止信号                         |                                                                                          |
| 16   | SIGSTKFLT | 栈溢出                           |                                                                                          |
| 17   | SIGCHLD   | 子进程退出                       | 默认忽略                                                                                 |
| 18   | SIGCONT   | 进程继续                         |                                                                                          |
| 19   | SIGSTOP   | 进程停止                         | 不能被忽略、处理和阻塞                                                                   |
| 20   | SIGTSTP   | 进程停止                         |                                                                                          |
| 21   | SIGTTIN   | 进程停止，后台进程从终端读数据时 |                                                                                          |
| 22   | SIGTTOU   | 进程停止，后台进程想终端写数据时 |                                                                                          |
| 23   | SIGURG    | I/O 有紧急数据到达当前进程       | 默认忽略                                                                                 |
| 24   | SIGXCPU   | 进程的 CPU 时间片到期            |                                                                                          |
| 25   | SIGXFSZ   | 文件大小的超出上限               |                                                                                          |
| 26   | SIGVTALRM | 虚拟时钟超时                     |                                                                                          |
| 27   | SIGPROF   | profile 时钟超时                 |                                                                                          |
| 28   | SIGWINCH  | 窗口大小改变                     | 默认忽略                                                                                 |
| 29   | SIGIO     | I/O 相关                         |                                                                                          |
| 30   | SIGPWR    | 关机                             | 默认忽略                                                                                 |
| 31   | SIGSYS    | 系统调用异常                     |                                                                                          |

对于 signal 信号，绝大部分的默认处理都是终止进程或停止进程，或 dump 内核映像转储。 上述的 31 的信号为非实时信号，其他的信号 32-64 都是实时信号。

## 信号产生

信号来源分为硬件类和软件类：

2.1 硬件方式

- 用户输入：比如在终端上按下组合键 ctrl+C，产生 SIGINT 信号；

- 硬件异常：CPU 检测到内存非法访问等异常，通知内核生成相应信号，并发送给发生事件的进程；

  2.2 软件方式

通过系统调用，发送 signal 信号：kill()，raise()，sigqueue()，alarm()，setitimer()，abort()

- kernel,使用 kill_proc_info(）等

- native,使用 kill() 或者 raise()等

- java,使用 Procees.sendSignal()等

## 信号注册和注销

3.1 注册

在进程 task_struct 结构体中有一个未决信号的成员变量 struct sigpending pending。每个信号在进程中注册都会把信号值加入到进程的未决信号集。

- 非实时信号发送给进程时，如果该信息已经在进程中注册过，不会再次注册，故信号会丢失；

- 实时信号发送给进程时，不管该信号是否在进程中注册过，都会再次注册。故信号不会丢失；

  3.2 注销

- 非实时信号：不可重复注册，最多只有一个 sigqueue 结构；当该结构被释放后，把该信号从进程未决信号集中删除，则信号注销完毕；

- 实时信号：可重复注册，可能存在多个 sigqueue 结构；当该信号的所有 sigqueue 处理完毕后，把该信号从进程未决信号集中删除，则信号注销完毕；

## 信号处理

内核处理进程收到的 signal 是在当前进程的上下文，故进程必须是 Running 状态。当进程唤醒或者调度后获取 CPU，则会从内核态转到用户态时检测是否有 signal 等待处理，处理完，进程会把相应的未决信号从链表中去掉。

4.1 处理时机

signal 信号处理时机： 内核态 -> signal 信号处理 -> 用户态：

- 在内核态，signal 信号不起作用；

- 在用户态，signal 所有未被屏蔽的信号都处理完毕；

- 当屏蔽信号，取消屏蔽时，会在下一次内核转用户态的过程中执行；

  4.2 处理方式

进程对信号的处理方式： 有 3 种

- 默认 接收到信号后按默认的行为处理该信号。 这是多数应用采取的处理方式。

- 自定义 用自定义的信号处理函数来执行特定的动作

- 忽略 接收到信号后不做任何反应。

  4.3 信号安装

进程处理某个信号前，需要先在进程中安装此信号。安装过程主要是建立信号值和进程对相应信息值的动作。

信号安装函数

- signal()：不支持信号传递信息，主要用于非实时信号安装；

- sigaction():支持信号传递信息，可用于所有信号安装；

其中 sigaction 结构体

- sa_handler:信号处理函数

- sa_mask：指定信号处理程序执行过程中需要阻塞的信号；

- sa_flags：标示位

- SA_RESTART：使被信号打断的 syscall 重新发起。

- SA_NOCLDSTOP：使父进程在它的子进程暂停或继续运行时不会收到 SIGCHLD 信号。

- SA_NOCLDWAIT：使父进程在它的子进程退出时不会收到 SIGCHLD 信号，这时子进程如果退出也不会成为僵 尸进程。

- SA_NODEFER：使对信号的屏蔽无效，即在信号处理函数执行期间仍能发出这个信号。

- SA_RESETHAND：信号处理之后重新设置为默认的处理方式。

- SA_SIGINFO：使用 sa_sigaction 成员而不是 sa_handler 作为信号处理函数。

函数原型：

int sigaction(int signum, const struct sigaction *act, struct sigaction*oldact);

- signum：要操作的 signal 信号。

- act：设置对 signal 信号的新处理方式。

- oldact：原来对信号的处理方式。

- 返回值：0 表示成功，-1 表示有错误发生。

  4.4 信号发送

- kill()：用于向进程或进程组发送信号；

- sigqueue()：只能向一个进程发送信号，不能像进程组发送信号；主要针对实时信号提出，与 sigaction()组合使用，当然也支持非实时信号的发送；

- alarm()：用于调用进程指定时间后发出 SIGALARM 信号；

- setitimer()：设置定时器，计时达到后给进程发送 SIGALRM 信号，功能比 alarm 更强大；

- abort()：向进程发送 SIGABORT 信号，默认进程会异常退出。

- raise()：用于向进程自身发送信号；

  4.5 信号相关函数

信号集操作函数

- sigemptyset(sigset_t \*set)：信号集全部清 0；

- sigfillset(sigset_t \*set)： 信号集全部置 1，则信号集包含 linux 支持的 64 种信号；

- sigaddset(sigset_t \*set, int signum)：向信号集中加入 signum 信号；

- sigdelset(sigset_t \*set, int signum)：向信号集中删除 signum 信号；

- sigismember(const sigset_t \*set, int signum)：判定信号 signum 是否存在信号集中。

信号阻塞函数

- sigprocmask(int how, const sigset_t *set, sigset_t*oldset))； 不同 how 参数，实现不同功能

- SIG_BLOCK：将 set 指向信号集中的信号，添加到进程阻塞信号集；

- SIG_UNBLOCK：将 set 指向信号集中的信号，从进程阻塞信号集删除；

- SIG_SETMASK：将 set 指向信号集中的信号，设置成进程阻塞信号集；

- sigpending(sigset_t \*set))：获取已发送到进程，却被阻塞的所有信号；

- sigsuspend(const sigset_t \*mask))：用 mask 代替进程的原有掩码，并暂停进程执行，直到收到信号再恢复原有掩码并继续执行进程。

# 信号发送工具

我们可以使用 [kill](/docs/1.操作系统/X.Linux%20管理/Linux%20系统管理工具/procps%20工具集.md#kill-向指定PID的进程发送信号) 命令行工具将指定的信号发送到指定的进程或进程组
