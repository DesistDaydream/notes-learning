---
title: 3.System Call
---

# 概述

> 参考：
> - [Manual(手册)，syscalls(2)](https://man7.org/linux/man-pages/man2/syscalls.2.html)
> - [Wiki,System_call](https://en.wikipedia.org/wiki/System_call)
> - <http://arthurchiao.art/blog/system-call-definitive-guide-zh/>

**System Call(系统调用，简称 syscall)** 是 Application(应用程序) 和 Linux Kernel(内核) 之间的基本接口。是操作内核的唯一入口。其实，所谓 syscall 就是各种编程语言中的 **Function(函数)** 概念。一个 syscall 也有名称、参数、返回值。syscall 即可以是名词，用来描述一个具体的 syscall；也可以是动词，用来表示某物调用了某个 syscall。当用户进程需要发生系统调用时，CPU 通过软中断切换到内核态开始执行内核系统调用函数。

> syscall 还有另一种意思，是一种编程方式，比如我们常说的 API，就是 syscall 的一种实现。

在 [syscalls(2) 手册中的 System call list 章节](https://man7.org/linux/man-pages/man2/syscalls.2.html#DESCRIPTION)可以看到 Linux 可用的完整的 syscall 列表。也就是说所有 Kernel 暴露出来的可供用户调用的 Function。

## 用户程序、内核和 CPU 特权级别

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/bbar3l/1616168230254-e3c38b73-8092-41bd-a17d-d3c4768de743.jpeg)

用户程序（例如编辑器、终端、ssh daemon 等）需要和 Linux 内核交互，内核代替它们完 成一些它们自身无法完成的操作。

例如，如果用户程序需要做 IO 操作（open、read、write 等），或者需要修改它的 内存地址（mmpa、sbrk 等），那它必须触发内核替它完成。

为什么禁止用户程序做这些操作呢？

因为 x86-64 CPU 有一个特权级别 （privilege levels）的概念。这个概念很复杂，完全可以单独写一篇博客。 出于本文讨论目的，我们将其（大大地）简化为如下：

1. 特权级别是权限控制的一种方式。当前的特权级别决定了允许执行哪些 CPU 指令和操作
2. 内核运行在最高级别，称为 “Ring 0”；用户程序运行在稍低的一个级别，一般称作 “Ring 3”

- 内核空间（Ring 0）具有最高权限，可以直接访问所有资源；
- 用户空间（Ring 3）只能访问受限资源，不能直接访问内存等硬件设备，必须通过系统调用陷入到内核中，才能访问这些特权资源。

用户程序要进行特权操作必须触发一次特权级别切换（从 “Ring 3” 到 “Ring 0”）， 由内核（替它）执行。触发特权级别切换有多种方式，我们先从最常见的方式开始：中断。

# Interrupts(中断)

详见：[Interrupts(中断) 概念详解](/docs/IT学习笔记/1.操作系统/2.Kernel(内核)/4.CPU%20 管理/Interrupts(中断).md 管理/Interrupts(中断).md)

# syscall 的方式

通过 glibc 提供的库函数

glibc 是 Linux 下使用的开源的标准 C 库，它是 GNU 发布的 libc 库，即运行时库。glibc 为程序员提供丰富的 API，除了例如字符串处理、数学运算等用户态服务之外，最重要的是封装了操作系统提供的系统服务，即系统调用的封装。那么 glibc 提供的系统调用 API 与内核特定的系统调用之间的关系是什么呢？

- 通常情况，每个特定的系统调用对应了至少一个 glibc 封装的库函数，如系统提供的打开文件系统调用 sys_open 对应的是 glibc 中的 open 函数；
- 其次，glibc 一个单独的 API 可能调用多个系统调用，如 glibc 提供的 printf 函数就会调用如 sys_open、sys_mmap、sys_write、sys_close 等等系统调用；
- 另外，多个 API 也可能只对应同一个系统调用，如 glibc 下实现的 malloc、calloc、free 等函数用来分配和释放内存，都利用了内核的 sys_brk 的系统调用。

举例来说，我们通过 glibc 提供的 chmod 函数来改变文件 etc/passwd 的属性为 444

```c
#include <sys/types.h>
#include <sys/stat.h>
#include <errno.h>
#include <stdio.h>
int main()
{
    int rc;
    	rc = chmod("/etc/passwd", 0444);
    if (rc == -1)
    	fprintf(stderr, "chmod failed, errno = %d\n", errno);
    else
    	printf("chmod success!\n");
    return 0;
}
```

在普通用户下编译运用，输出结果为：

    chmod failed, errno = 1

上面系统调用返回的值为-1，说明系统调用失败，错误码为 1，在 /usr/include/asm-generic/errno-base.h 文件中有如下错误代码说明：

    #define EPERM       1                /* Operation not permitted */

即无权限进行该操作，我们以普通用户权限是无法修改 /etc/passwd 文件的属性的，结果正确。

## 使用指定的 SyscallName 直接调用

使用上面的方法有很多好处，首先你无须知道更多的细节，如 chmod 系统调用号，你只需了解 glibc 提供的 API 的原型；其次，该方法具有更好的移植性，你可以很轻松将该程序移植到其他平台，或者将 glibc 库换成其它库，程序只需做少量改动。

但有点不足是，如果 glibc 没有封装某个内核提供的系统调用时，我就没办法通过上面的方法来调用该系统调用。如我自己通过编译内核增加了一个系统调用，这时 glibc 不可能有你新增系统调用的封装 API，此时我们可以利用 glibc 提供的 syscall 函数直接调用。该函数定义在 unistd.h 头文件中，函数原型如下：

long int syscall (long int sysno, ...)

- sysno 是系统调用号，每个系统调用都有唯一的系统调用号来标识。在 sys/syscall.h 中有所有可能的系统调用号的宏定义。
- ... 为剩余可变长的参数，为系统调用所带的参数，根据系统调用的不同，可带 0~5 个不等的参数，如果超过特定系统调用能带的参数，多余的参数被忽略。
- 返回值 该函数返回值为特定系统调用的返回值，在系统调用成功之后你可以将该返回值转化为特定的类型，如果系统调用失败则返回 -1，错误代码存放在 errno 中。

还以上面修改 /etc/passwd 文件的属性为例，这次使用 syscall 直接调用：

    #include <stdio.h>
    #include <unistd.h>
    #include <sys/syscall.h>
    #include <errno.h>
    int main()
    {
    int rc;
            rc = syscall(SYS_chmod, "/etc/passwd", 0444);
    if (rc == -1)
    fprintf(stderr, "chmod failed, errno = %d\n", errno);
    else
    printf("chmod succeess!\n");
    return 0;
    }

在普通用户下编译执行，输出的结果与上例相同。

## 通过 syscall() 间接调用

## 通过 int 指令陷入

如果我们知道系统调用的整个过程的话，应该就能知道用户态程序通过软中断指令 int 0x80 来陷入内核态（在 Intel Pentium II 又引入了 sysenter 指令），参数的传递是通过寄存器，eax 传递的是系统调用号，ebx、ecx、edx、esi 和 edi 来依次传递最多五个参数，当系统调用返回时，返回值存放在 eax 中。

仍然以上面的修改文件属性为例，将调用系统调用那段写成内联汇编代码：

```c
#include <stdio.h>
#include <sys/types.h>
#include <sys/syscall.h>
#include <errno.h>
int main()
{
long rc;
char *file_name = "/etc/passwd";
unsigned short mode = 0444;
asm(
"int $0x80"
: "=a" (rc)
: "0" (SYS_chmod), "b" ((long)file_name), "c" ((long)mode)
);
if ((unsigned long)rc >= (unsigned long)-132) {
                errno = -rc;
                rc = -1;
}
if (rc == -1)
fprintf(stderr, "chmode failed, errno = %d\n", errno);
else
printf("success!\n");
return 0;
}
```

如果 eax 寄存器存放的返回值（存放在变量 rc 中）在 -1~-132 之间，就必须要解释为出错码（在/usr/include/asm-generic/errno.h 文件中定义的最大出错码为 132），这时，将错误码写入 errno 中，置系统调用返回值为 -1；否则返回的是 eax 中的值。

上面程序在 32 位 Linux 下以普通用户权限编译运行结果与前面两个相同！
