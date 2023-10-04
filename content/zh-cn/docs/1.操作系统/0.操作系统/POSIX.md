---
title: POSIX
---

# 概述

> 参考：
>
> - [Wiki，POSIX](https://en.wikipedia.org/wiki/POSIX)
> - [腾讯云，什么是 POSIX](https://cloud.tencent.com/developer/ask/26856)

**Portable Operating System Interface(便携式操作系统接口，POSIX)** 是 [IEEE](docs/x_标准化/IT/IEEE.md) 的计算机协会指定的一系列标准，用于维护操作系统之间的兼容性。 POSIX 定义了应用程序编程接口 (API)，以及命令行 shell 和实用程序接口，以实现与 类 Unix 操作系统 和 其他操作系统的软件兼容性。

一般情况下，应用程序通过应用编程接口(API)而不是直接通过系统调用来编程。这点很重要，因为应用程序使用的这种编程接口实际上并不需要和内核 提供的系统调用对应。一个 API 定义了一组应用程序使用的编程接口。它们可以实现成一个系统调用，也可以通过调用多个系统调用来实现，而完全不使用任何系 统调用也不存在问题。实际上，API 可以在各种不同的操作系统上实现，给应用程序提供完全相同的接口，而它们本身在这些系统上的实现却可能迥异。

在 Unix 世界中，最流行的应用编程接口是基于 POSIX 标准的。从纯技术的角度看，POSIX 是由 IEEE 的一组标准组成，其目标是提供一套大体上基于 Unix 的可移植操作系统标准。Linux 是与 POSIX 兼容的。

POSIX 是说明 API 和系统调用之间关系的一个极好例子。在大多数 Unix 系统上，根据 POSIX 而定义的 API 函数和系统调用之间有着直接关 系。实际上，POSIX 标准就是仿照早期 Unix 系统的界面建立的。另一方面，许多操作系统，像 Windows NT，尽管和 Unix 没有什么关系，也提供了与 POSIX 兼容的库。

Linux 的系统调用像大多数 Unix 系统一样，作为 C 库的一部分提供如图 5-1 所示。如图 5-1 所示 C 库实现了 Unix 系统的主要 API，包括标 准 C 库函数和系统调用。所有的 C 程序都可以使用 C 库，而由于 C 语言本身的特点，其他语言也可以很方便地把它们封装起来使用。此外，C 库提供了 POSIX 的绝大部分 API。

从程序员的角度看，系统调用无关紧要；他们只需要跟 API 打交道就可以了。相反，内核只跟系统调用打交道；库函数及应用程序是怎么使用系统调用不是内核所关心的。

简单总结：

完成同一功能，不同内核提供的系统调用（也就是一个函数）是不同的，例如创建进程，linux 下是 fork 函数，windows 下是 creatprocess 函数。好，我现在在 linux 下写一个程序，用到 fork 函数，那么这个程序该怎么往 windows 上移植？我需要把源代码里的 fork 通通改成 creatprocess，然后重新编译...

POSIX 标准的出现就是为了解决这个问题。linux 和 windows 都要实现基本的 POSIX 标准，Linux 把 fork 函数封装成 posix_fork（随便说的），windows 把 creatprocess 函数也封装成 posix_fork，都声明在 unistd.h 里。这样，程序员编写普通应用时候，只用包含 unistd.h，调用 posix_fork 函数，程序就在源代码级别可移植了
