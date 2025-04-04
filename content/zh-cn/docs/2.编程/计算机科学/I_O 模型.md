---
title: I/O 模型
linkTitle: I/O 模型
weight: 20
---

# 概述

> 参考：
>
> - [Wiki, Asynchronous_I/O](https://en.wikipedia.org/wiki/Asynchronous_I/O)

编程中的 [I_O](/docs/0.计算机/I_O.md)

同步/异步：关注的是消息通信机制，被调用者在收到调用请求后，是否立即返回，还是得到最终结果后才返回,立即返回为异步，等待结果再返回为同步，异步不会影响调用者处理后续

同步和异步通常用来形容一次方法调用。

同步方法调用一旦开始，调用者必须等到方法调用返回后，才能继续后续的行为。

异步方法调用更像一个消息传递，一旦开始，方法调用就会立即返回，调用者就可以继续后续的操作。而，异步方法通常会在另外一个线程中，“真实”地执行着。整个过程，不会阻碍调用者的工作。举个例子

- 你打电话问书店老板有没有《分布式系统》这本书，如果是同步通信机制，书店老板会说，你稍等，”我查一下"，然后开始查啊查，等查好了（可能是5秒，也可能是一天）告诉你结果（返回结果）。
- 而异步通信机制，书店老板直接告诉你我查一下啊，查好了打电话给你，然后直接挂电话了（不返回结果）。然后查好了，他会主动打电话给你。在这里老板通过“回电”这种方式来回调。

阻塞/非阻塞：关注的是程序在等待调用结果（消息，返回值）时的状态.，调用者发起调用请求后，在收到响应结果之前是否会被挂起，被挂起为阻塞，不被挂起为非阻塞。举个例子

- 同步阻塞：
  - 老张在厨房用普通水壶烧水，一直在厨房等着（阻塞），盯到水烧开（同步）；
- 异步阻塞：
  - 老张在厨房用响水壶烧水，一直在厨房中等着（阻塞），直到水壶发出响声（异步），老张知道水烧开了；
- 同步非阻塞：
  - 老张在厨房用普通水壶烧水，在烧水过程中，就到客厅去看电视（非阻塞），然后时不时去厨房看看水烧开了没（轮询检查同步结果）；
- 异步非阻塞：
  - 老张在厨房用响水壶烧水，在烧水过程中，就到客厅去看电视（非阻塞），当水壶发出响声（异步），老张就知道水烧开了。
- 所谓同步异步，只是对于水壶而言。
  - 普通水壶，同步；响水壶，异步。
  - 虽然都能干活，但响水壶可以在自己完工之后，提示老张水开了。这是普通水壶所不能及的。
  - 同步只能让调用者去轮询自己（情况 2 中），造成老张效率的低下。
- 所谓阻塞非阻塞，仅仅对于老张而言。
  - 立等的老张，阻塞；看电视的老张，非阻塞。
  - 情况 1 和情况 3 中老张就是阻塞的，媳妇喊他都不知道。虽然 3 中响水壶是异步的，可对于立等的老张没有太大的意义。所以一般异步是配合非阻塞使用的，这样才能发挥异步的效用。
