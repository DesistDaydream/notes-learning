---
title: Goroutine AND Channel
linkTitle: Goroutine AND Channel
weight: 20
---

# 概述

> 参考:
>
> - [公众号,马哥 Linux 运维-golang channel 使用总结](https://mp.weixin.qq.com/s/67XCr7nc_q3gHu6MMO606A)
> - [公众号，田飞雨-Golang GPM 模型剖析](https://mp.weixin.qq.com/s/CJ9renN4ho7sasGhYSWBqw)

# Go 语言的并发

**不要通过共享来通信，而要通过通信来共享。**

通过 `Goroutines(协程)` 与 `Channels(通道)` 实现 `并发编程`

**并发**与**并行**的区别

- **Concurrency(并发)** # 一个处理器或者内核上，一个并发程序可以使用多个线程来**交替运行**。反义词为顺序
  - 实际例子：你有一张嘴，电话来了，你停了下来接了电话，吃一口说一句话，直到说完话并且吃完饭，这是并发
- **Parallelism(并行)** # 多个处理器或者多核上，一个程序在某个时间、在多个处理器上**同时运行**。反义词为串行
  - 实际例子：你有两张嘴，电话来了，你一边打电话一边吃饭，直到说完话并且吃完饭，这是并行

并行是并发的真子集。并发不全是并行，但并行一定并发。(单核并发不并行，多核并行也属于并发)，除非该程序无法使用多线程执行任务。

并行不一定会加快运行速度，因为并行运行的组件之间可能需要相互通信。比如运行在两个 CPU 上的两个组件之间需要互相通信。并发系统上，这种通信开销很小。但在多核的并行系统上，组件间的通信开销就很高了。所以，并行不一定会加快运行速度！

一个程序是运行在机器上的一个进程，进程是一个运行在自己内存空间里的独立执行体。一个进程由一个或多个操作系统线程组成，这些线程其实是共享同一个内存地址空间在一起工作的执行体。几乎所有`正式`的程序都是多线程的，以便让用户或计算机不必等待，或能够同时服务多个请求(e.g.Web 服务器)，或增加性能和吞吐量。

不要使用全局变量或共享内存，他们会给代码在并发运算的时候带来危险。

# Goroutines(协程)

处理应用程序**并发**功能的就是 `Goroutines(协程)`

Go 协程是与其他函数或方法一起并发运行的函数或方法。

调用函数或者方法时，在前面加上关键字 go，可以让一个新的 Go 协程并发地运行。i.e.有关键字`go`的函数或方法，即算协程，可以并发运行。 main()函数算主协程，可以没有`go`关键字

一个基本的协程代码示例：

```go
package main

import (
 "fmt"
 "sync"
 "time"
)

var wg sync.WaitGroup

func hello(m string) {
 if m == "waitgroup" {
        // 让 WaitGroup 计数器 -1
  defer wg.Done()
 }

 i := 0
 for {
  fmt.Printf("Hello world goroutine,%v\n", i)
  i += 1
  if i == 3 {
   return
  }
 }
}

func timeSleepGotourine() {
 // 调用了 go hello() 之后，程序控制没有等待 hello 协程结束，
 // 立即返回到了代码下一行，打印 main function。
 // 接着由于没有其他可执行的代码，Go 主协程终止，
 // 于是 hello 协程就没有机会运行了。
 go hello("timesleep")

 // 调用了 time 包里的函数 Sleep，该函数会休眠执行它的 Go 协程。
 // 在这里，我们使 Go 主协程休眠了 1 秒。因此在主协程终止之前，
 // 调用 go hello() 就有足够的时间来执行了。
 // 该程序首先打印 Hello world goroutine，等待 1 秒钟之后，接着打印 main function。
 time.Sleep(1 * time.Second)
 // 也可以使用下面的方式让程序在手动按回车才结束。
 // fmt.Scanln()
 fmt.Println("main function")
}

func waitGroupGotoutine() {
    // 为 WaitGroup 计数器 +1
 wg.Add(1)

 go hello("waitgroup")

 // 等待 WaitGroup 计数器归零，归零后，wg.Wait() 释放，并继续执行后面的代码、
 wg.Wait()

 fmt.Println("main function")
}

func main() {
 // 通过睡眠，让 main() 等待协程完成
 timeSleepGotourine()
 // 通过 WaitGroup，让 main() 等待协程完成
 waitGroupGotoutine()
}
```

# Channels(通道)

`Channels(通道)`，可以想象成 Goroutines 之间通信的管道。如同管道中的水会从一端流到另一端，通过使用通道，数据也可以从一端发送，在另一端接收。

如果两个协程需要通信，则必须要给它们同一个通道作为参数才可以。

基本声明格式：`var ChanID chan Type`。声明一个通道，指定这个通道里可以传输的类型是什么

初始化格式：`ChanID = make(chan Type)`。因为 channel 是引用类型，所以可以使用 make()函数来给它分配内存 channel 的操作符：`<-`。这个操作符表示数据按照箭头的方向流动。下面有几个例子：

- `ch <- int1`。流向通道(发送)。表示发送变量 int1 的数据到通道 ch 中。i.e.int1 变量中的数据会发送给通道 ch。
- `int2 = <- ch`。从通道流出(接收)。表示变量 int2 从通道 ch 中接收数据。i.e.ch 中的数据会发送给 int2
- `<- ch`。单独调用通道的值，当前值会被丢弃。

一个基本的通道代码示例：

```go
package main

import (
 "fmt"
)

func correct() {
 // 使用make函数创建一个新的通道。通道类型就是他们需要传递值的类型
 messages := make(chan string)

 // 使用`ChannelID <-`语法发送(send)一个新的值到通道中
 go func() { messages <- "ping" }()

 // 使用`<- ChannelID`语法从通道中接收(receives)一个值
 msg := <-messages
 fmt.Println(msg)
}

func error() {
 messages := make(chan string)

 // 如果不使用协程，则运行程序时，会报“死锁”的错,错误信息如下：
 // fatal error: all goroutines are asleep - deadlock!
 // 因为，代码是一行一行执行的，如果一个没有缓存的通道在接收数据之后，需要同步把数据发送给接收者。
 // 可是当前行的代码还没执行完，怎么能执行后面的呢，没有后面的代码，也就没有接收者，所以这就是错误产生的原因。
 // 当通道使用协程的方式运行时，就算当前时刻没有接收者，这个通过协程运行起来的通道，一会自动阻塞，并等待接收者。否则不通过协程启动通道，那么就跟普通代码一样。
 func() { messages <- "ping" }()

 // 可以使用如下方式直接输出通道内的数据，相当于fmt.Println就是通道的接收者
 fmt.Println(-messages)
}

func buffer() {
 messages := make(chan string, 1)

 // 带缓冲的通道可以不使用协程。
 messages <- "ping"

 fmt.Println(-messages)
}

func main() {
 fmt.Println("1.通道正确的示例")
 // correct()
 fmt.Println("2.通道会导致死锁的示例")
 // error()
 fmt.Println("3.通道缓冲")
 buffer()

}

// 默认发送和接收操作是阻塞的，直到发送方和接收方都准备完毕。
// 这个特性允许我们，不使用任何其它的同步操作，来在程序结尾等待消息"ping"
```

## 通道阻塞

默认情况下，数据通信是同步且无缓冲的，i.e.一边发就需要一边接收，在有接收者接收数据之前，发送不会结束，也不会继续发送新数据，这种情况称为**阻塞**。必须要有一个接收者准备好接收通道的数据，然后发送者可以直接把数据发送给接收者。或者可以使用带缓冲的通道，这样，在给通道发送数据时，可以把数据先存储在缓冲区，而不用直接让接收者接收。

对于通道阻塞有非常严格的规定：

- 对于同一个通道，发送操作（协程或者函数中的），在接收者准备好之前是阻塞的：如果 ch 中的数据无人接收，就无法再给通道传入其他数据：新的输入无法在通道非空的情况下传入。所以发送操作会等待 ch 再次变为可用状态：就是通道值被接收时（可以传入变量）。
- 对于同一个通道，接收操作是阻塞的（协程或函数中的），直到发送者可用：如果通道中没有数据，接收者就阻塞了。

为通道提供数据的也叫**生产者**，从通道中拿去数据的也叫**消费者**

**注意:**

- channel.go 例子中，有通道的错误使用方法，请注意！

## 带缓冲的通道

格式：`ChanID = make(chan Type, CapValue)`。
CapValue 为容量值。如果容量大于 0，通道就是异步的了：缓冲满载（发送）或变空（接收）之前通信不会阻塞，元素会按照发送的顺序被接收。如果容量是 0 或者未设置，通信仅在收发双方准备好的情况下才可以成功。

## 单向通道

- `var VarID chan<- int`。通道仅能发送数据。i.e.关键字：`chan<-`
- `var VarID <-chan int`。通道仅能接收数据。i.e.关键字：`<-chan`

只接收的通道 `<-chan T` 无法关闭，因为关闭通道是发送者用来表示不再给通道发送值了，所以对只接收通道是没有意义的。\

# 通道与协程的配合

展示了一个通道连通了两个协程，使得数据从一个协程进入通道，再从通道另一端出去到另一个协程的过程

```go
package main

import (
 "fmt"
 "sync"
 // "time"
)

func sendData(ch chan string, done chan bool, wg sync.WaitGroup) {
 // 该函数把几个字符串发送给通道ch，使得这些字符串会存在一个管道中，等待输出
 ch <- "Tianjin"
 ch <- "Beijing"
 ch <- "China"

 // 当完成数据发送时，发送通知
 done <- true

 // 让 WaitGroup 计数器 -1
 wg.Done()
}

func recvData(ch chan string, done chan bool, wg sync.WaitGroup) {
I:
 for {
  select {
  // 把通道中的字符串都赋值给变量input
  case input := <-ch:
   fmt.Println(input)

  // 当通道通知完成数据传输时，跳出循环
  case <-done:
   break I
  }
 }

 // 让 WaitGroup 计数器 -1
 wg.Done()
}

func main() {
 // 声明 WaitGroup 计数器
 // WaitGroup 等待一组 goroutine 完成。主 goroutine 调用 Add 来设置要等待的 goroutine 的数量。
 // 然后每个 goroutine 运行并在完成时调用 Done。 同时，Wait 可用于阻塞，直到所有 goroutine 完成。
 var wg sync.WaitGroup
 // 为 WaitGroup 计数器 +2，当计数器减为 0 时，wg.Wait() 将会释放，并继续执行后面的代码。
 wg.Add(2)
 // ch通道用来传递数据；done通道用来传递任务是否完成的消息。
 ch := make(chan string)
 done := make(chan bool)

 // 把通道ch作为参数传递给两个协程函数
 // 可以理解为把通道的两端分别连接到两个协程函数上
 go sendData(ch, done, wg)
 go recvData(ch, done, wg)

 // 通过 wg.Wait() 让 main() 等待协程完成。
 // 如果不让 main() 等待，则无任何输出，或者报错 panic: send on closed channel。因为协程是并发运行，不用等代码运行完成就会执行后续代码。
 // 如果后续代码执行完了，协程中的代码还没执行完成，就会没有任何输出。
 // 如果后续代码中包含了关闭通道的操作，那么程序将会 panic，并报错 send on closed channel
 wg.Wait()
}
```

比如聊天室功能中俩人聊天，两个 goroutine 就是两人的聊天窗口，channel 就是连接两人的管道，直到一个人向 channel 发送消息前，channel 都是阻塞的，当发送消息后，消息就会从 channel 的另一头流出，让另一个人收到。
