---
title: Callback
linkTitle: Callback
weight: 20
---

# 概述

> 参考：
>
> - [Wiki, Callback_(computer_programming)](https://en.wikipedia.org/wiki/Callback_(computer_programming))
> - [博客园，究竟什么是callback function(回调函数)](https://www.cnblogs.com/ArsenalfanInECNU/p/14650501.html)
> - [公众号-码农的荒岛求生，回调函数 callback 的实现原理是什么？](https://mp.weixin.qq.com/s/zS7URRO5sNzobUNIqSJHIg)
>   - 原文非简化版：《[10 张图让你彻底理解回调函数](http://mp.weixin.qq.com/s?__biz=Mzg4OTYzODM4Mw==&mid=2247485712&idx=1&sn=3d2750dfb693f41b2483b51b60a4f44c&chksm=cfe99590f89e1c860277fe1b22c3731ec4e3b61dbb5cd2a6d9548efbc709104a38d6da812517&scene=21#wechat_redirect)》

**Callback function(回调函数)** 也是函数，只不过函数的参数不是变量，而是另一个函数。这种调用函数的方式有多种好处

- 异步调用。调用 A 时，只要传递的参数中的函数没有阻塞逻辑，那就不用等待 A 函数全部执行完成，即可继续处理后续代码。
- 不同实体调用函数 A 时，想要执行一些不同的特定的代码，不用在函数 A 里加很多 if else 的判断
- 等等

**其实回调函数和普通函数没有本质的区别。**

首先让我们来看看普通的函数调用，假设我们在 A 函数中调用函数 func：

```swift
void A() {
   ...
   func();
   ...
}
```

想一想，你怎么知道可以调用 func 呢？哦，原来 func 是你自己定义的：

```swift
void func() {
  blablabla;
}
```

这很简单吧，现在假设你编写的这段代码无比之牛逼，全世界的程序员都无比疯狂的想引入到自己的项目中，这时你会把 A 函数编写成一个库供全世界的码农使用。

但此时所有人都发现一个问题，**那就是他们都想在 A 函数中的某个特定点上执行一段自己的代码**，作为这个库的创作者你可能会这样实现：

```cs
void A() {
   ...

   if (张三) {
     funcA();
   } else if (李四) {
     funcB();
   }
   ...
}
```

假设全世界有一千万码农，那你是不是要有**一千万个 if else**。。。想想这样的代码就很刺激有没有！

更好的办法是什么呢？**把函数也当做变量**！你可以这样定义 A 函数：

```cs
void A(func f) {
   ...
   f();
   ...
}
```

任何一个码农在调用你的 A 函数时传入一个函数变量，A 函数在合适的时机调用你传入的函数变量，**从而节省了一千万行代码**。

**为了让这个过程更加难懂一些，这个函数变量起了一个无比难懂的称呼：回调函数。**

现在你应该明白了回调函数是怎么一回事了吧，相比回调函数来说我更愿意将其看做**函数变量**。

以上就是回调函数的基本原理，有想看更详细版本的请参考[公众号 - 码农的荒岛求生，10张图让你彻底理解回调函数](https://mp.weixin.qq.com/s/eFYM4uOIF09t8b9tTD523A)

以上仅仅是回调函数的一种用途，回调函数在基于事件驱动编程以及异步编程时更是必备的，关于事件驱动编程你可以参考[公众号 - 码农的荒岛求生，高并发高性能服务器是如何实现的](https://mp.weixin.qq.com/s/Z07Hc9SRfGz6n8XhFHGVyA)，GUI 编程的同学对此肯定很熟悉。
