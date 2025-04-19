---
title: RPC
linkTitle: RPC
weight: 1
---

# 概述

> 参考：
>
> - [Wiki, RPC](https://en.wikipedia.org/wiki/Remote_procedure_call)
> - [Wiki, gRPC](https://en.wikipedia.org/wiki/GRPC)
> - [gRPC 官网](https://grpc.io/)

在[分布式计算](https://en.wikipedia.org/wiki/Distributed_computing)中，**Remote Procedure Call(远程过程调用，简称 RPC)** 是计算机程序使 [Subroutine](/docs/2.编程/计算机科学/Function/Function.md) 在不同的地址空间（通常在共享网络上的另一台计算机上）执行时，被编码为 **Local Procedure Call(本地过程调用)**，而无需程序员为远程交互显式编写细节。也就是说，程序员可以为程序编写相同的代码，而不用关心自己编写的程序将会被本地调用还是远程调用。

其实 LPC 和 RPC 并不是对比的最佳选择，两者都 IPC 的一种方式，也就是说都是两个进程间通讯的一种方式，可能来说，LPC 与 RPC 最大的区别在于是否基于 TCP/IP 来让两个进程进行通信。而如果从网络间两个进程通信的角度看，RPC 又可以与 HTTP 进行对比。

从某种角度来说， HTTP 其实就是一种 RPC

- HTTP 发起请求的 URL 就是 RPC 发起请求的函数名
- 请求体就是函数的参数
- 响应体就是函数的函数中的处理逻辑或返回值

只不过 HTTP 是一个协议(也可以说是一种交互标准)，而 RPC 是一种方式、方法，可以使用 HTTP 来进行 RPC 通信，也可以使用其他协议进行 RPC 通信。如果使用 HTTP 标准进行 RPC 通信，那 RPC 的 C/S 之间就是通过文本格式进行交互；但是 RPC 通信最常使用的是 Protobuf 数据格式进行通信。

> 这里说的使用“HTTP 进行 RPC 通信”指的是使用 xml、json 等格式的数据进行 RPC 通信。而在很多 RPC 框架中，RPC 之间交互的信息与 HTTP 之间交互的信息，是可以互通的！~

**RPC 最常见的场景就是“微服务”**，将一个大而全的产品拆分成多个服务，如果通过 HTTP 调用，那么调用函数时就需要转换为调用 URL，对于关联性非常难强的多个服务来说，这种交互是灾难性的，如果网络上的多个服务之间，可以直接通过函数调用，那么代码写起来，也是非常简洁的。

通常来说，如果想要调用第三方平台提供的接口，使用 HTTP，而一个产品中关联性非常强，甚至可以合并成一个服务的多个服务之间的接口调用，就要使用 RPC 了，公司内服务之间的 RPC 调用，可以通过定制化的协议来使得通信更高效。
