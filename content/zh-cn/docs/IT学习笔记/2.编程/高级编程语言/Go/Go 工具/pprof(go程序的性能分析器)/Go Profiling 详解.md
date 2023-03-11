---
title: Go Profiling 详解
---

原文：<https://mp.weixin.qq.com/s/SveQPLr7abKXccLpYKkNKA>

## 介绍

### Read This

这是一本实用指南，针对的是对使用 profiling、跟踪和其他可观察性技术改进程序感兴趣的忙碌的 gopher。如果你对 Go 的内部实现不熟悉，建议你先阅读整个介绍。之后，你可以自由地跳到你感兴趣的任何章节。

### Go 心智模型

虽然不了解 Go 语言的底层原理，但是你也能熟练地编写 Go 代码。但是当涉及到性能和调试时，你会从内部实现中受益匪浅。我们接下来会概述 Go 基本原理的模型。这个模型应该足以让你避免最常见的错误，但是所有的模型都是错误的\[2]，所以我们鼓励你寻找更深入的材料来解决未来的难题。

Go 的主要工作是对硬件资源进行复用和抽象，类似于操作系统。主要使用两个抽象：

1. **Goroutine Scheduler(goroutine 调度器):** 管理你的代码如何在 CPU 上执行。
2. **Garbage Collector(垃圾回收):** 提供虚拟内存，在需要时自动释放。

#### Goroutine 调度器

让我们先用下面的例子来谈谈调度器 :

`func main() {     res, err := http.Get("https://example.org/")     if err != nil {         panic(err)     }     fmt.Printf("%d\n", res.StatusCode) }`

这里我们有一个 goroutine，我们称之为 `G1`，它运行 `main` 函数。下图显示了这个 goroutine 如何在单个 CPU 上执行的简化的时间线。最初 `G1` 在 CPU 上运行以准备 http 请求。然后 CPU 变得空闲，因为 goroutine 需要网络等待。最后它再次被调度到 CPU 上，打印出状态代码。

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/4aaac374-ac48-48d4-8190-b524f880bf2f/640)

从调度器的角度来看，上述程序的执行情况如下所示。一开始，`G1` 在 `CPU 1` 上 `Executing`。然后，在 `Waiting` 网络的过程中，goroutine 被从 CPU 上取出。一旦调度器注意到网络已经返回（使用非阻塞 I/O，类似于 Node.js），它就把 goroutine 标记为 `Runnable`。一旦有 CPU 核心可用，goroutine 就会再次开始 `Executing`。在我们的例子中，所有的 cpu 核都是可用的，所以 `G1` 可以立即回到其中一个 CPU 上 `Executing`  `fmt.Printf()` 函数，而无需在 `Runnable` 状态下花费任何时间。

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/4aaac374-ac48-48d4-8190-b524f880bf2f/640)

大多数时候，Go 程序都在运行多个 goroutine，所以你会有一些 goroutine 在一些 CPU 核心上 `Executing`，大量 goroutine 由于各种原因 `Waiting`，理想情况下没有 goroutine `Runnable`，除非你的程序显示非常高的 CPU 负载。下面就是一个例子。

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/4aaac374-ac48-48d4-8190-b524f880bf2f/640)

当然，上面的模型掩盖了许多细节。Go 调度器是在操作系统管理的线程之上运行的，甚至 CPUs 本身也能够实现超线程，这可以看作是一种调度形式。所以如果你感兴趣的话，可以通过 Ardan labs 的 Scheduling in Go\[3] 或其他资料继续深入这个主题。

但是，上面的模型应该足以理解本指南的其余部分。特别是，各种 Go profilers 所测量的时间基本上是 goroutines 在 `Executing` 和 `Waiting` 状态中所花费的时间，如下图所示。

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/4aaac374-ac48-48d4-8190-b524f880bf2f/640)

#### 垃圾收集器

Go 的另一个主要抽象是垃圾收集器。在 C 语言中，程序员需要使用 `malloc()` 和 `free()` 手动分配和释放内存。这提供了很好的控制，但在实践中很容易出错。垃圾收集器可以减少这种负担，但内存的自动管理很容易成为性能瓶颈。本指南将为 Go 的 GC 提供一个简单模型，这个模型对于识别和优化内存管理相关问题非常有用。

##### The Stack

让我们从最基础开始。Go 可以将内存分配到堆栈或者堆的任意一个位置。每个 goroutine 都有自己的堆栈，栈是一个连续的内存区域。此外，goroutine 之间还有一大块内存共享区域，这就是堆。如下所示。

当一个函数调用另一个函数时，它会在堆栈中获得自己的部分，称为 stack frame(堆栈帧) ，在这里它可以创建局部变量等。堆栈指针用于标识帧中的下一个空闲位置。当一个函数返回时，通过简单地将堆栈指针移回到前一帧的末尾，最后一帧中的数据就会被丢弃。帧的数据本身可以留在堆栈中，并被下一个函数调用覆盖。这是非常简单和有效的，因为 Go 不需要跟踪每个变量。

为了让这个更直观一点，让我们看看下面的例子 :

\`func main() {
 sum := 0
 sum = add(23, 42)
 fmt.Println(sum)
}

func add(a, b int) int {
 return a + b
}

\`

我们有一个 `main()` 函数，它一开始就在堆栈中为变量 `sum` 创建了空间。当调用 `add()` 函数时，它获得自己的帧来保存本地的 `a` 和 `b` 参数。一旦 `add()` 返回，通过将堆栈指针移回 `main()` 函数帧末尾，它的数据就会被丢弃，`sum` 变量就会得到更新。同时，`add()` 的旧值逗留在堆栈指针之外，以便被下一个函数调用覆盖。下面是这个过程的可视化图 :

上面的例子是高度简化了，省略了许多关于返回值、帧指针、返回地址和函数内联的细节。实际上，在 Go 1.17 中，上面的程序甚至可能不需要堆栈上的任何空间，因为少量的数据可以由编译器使用 CPU 寄存器来管理。但这也没关系。这个模型应该还是能让你对简单的 Go 程序在堆栈上分配和丢弃局部变量的方式有一个直观认识。

此时你可能会想，如果堆栈上的空间用完了会怎么样。在像 C 这样的语言中，这会导致堆栈溢出错误。而 Go 会自动处理这个问题，扩容成两倍堆栈。所以一般 goroutine 开始都是很小的，通常是 2 KiB，这也是 goroutine 比操作系统线程更具可伸缩性的关键因素\[4]之一。

##### The Heap

堆栈分配很 nice，但是在很多情况下 Go 不能全部都使用它们。最常见比如函数返回指针。把上面的 `add()` 函数示例修改一下，如下 :

\`func main() {
 fmt.Println(\*add(23, 42))
}

func add(a, b int) \*int {
 sum := a + b
 return \&sum
}

\`

通常 Go 可以把 `add()` 函数内部的 `sum` 变量分配到堆栈中。我们已经知道当 `add()` 函数返回时，这些数据将被丢弃。因此，为了安全返回 `&sum` 指针，Go 必须从堆栈外为其分配内存。这就是堆的作用。

堆用于存储可能比创建它的函数存活时间更长的数据，以及任何使用指针在 goroutine 之间共享数据。然而，这就提出了一个问题 : 这些内存是如何被释放的。因为与堆栈分配不同，当创建堆分配的函数返回时，堆分配是不能被丢弃。

Go 使用其内置的垃圾收集器解决了这个问题。其实现的细节非常复杂，但是从宏观角度来看，它可以跟踪你的内存，如下图所示。在这里你可以看到三个 goroutines，它们具有指向堆上绿色分配的指针。其中一些分配还有指向其他分配的指针，用绿色显示。此外，灰色分配可能指向绿色分配或相互指向，但它们本身并不被绿色分配引用。这些分配曾经是可到达的，但现在被认为是垃圾。如果在堆栈上分配它们的指针的函数返回，或者它们的值被覆盖，就会发生这种情况。GC 负责自动识别和释放这些分配。

执行 GC 涉及大量开销很大的图遍历和缓冲区的处理。它甚至需要定期 stop-the-world 阶段来停止整个程序执行。Go 现在的版本已经把 GC 的过程打到很恐怖的速度了 (毫秒级)，剩余的大部分开销都是任何 GC 算法都跑不掉的的。事实上，Go 程序执行中 20-30% 的时间都开销在内存管理上并不罕见。

一般来说，GC 的成本与程序执行的堆分配量成正比。因此，在优化程序的内存开销时，需要注意的是 :

- **Reduce**: 尝试将堆分配转换为堆栈分配，或者干脆避免它们。把堆上的指针数量打下来也会有帮助。
- **Reuse:** 复用堆分配，而不是用新的堆分配来替换它们。
- **Recycle:** 有些堆分配是无法避免的。让 GC 自动回收它们，并关注其他问题。

与本指南中之前的心智模型一样，上面的流程都是对实际的情况做了简化。但是希望它能够很好地帮你理解本指南的其余部分，并激励你阅读更多关于这个主题的文章。推荐你必读的一篇文章《Getting to Go: The Journey of Go's Garbage Collector\[5]》 ，它让你很好地了解 Go 的 GC 多年来是如何进步的，以及它的改进速度。

## Go Profilers

以下是 Go runtime 内置 profilers 的概述。有关更多详细信息，请点击链接。

|

| CPU               | Memory  | Block   | Mutex   | Goroutine | ThreadCreate |
| ----------------- | ------- | ------- | ------- | --------- | ------------ |
| Production Safety | ✅      | ✅      | ✅      | ✅        | ⚠️ (1.)      |
| Safe Rate         | default | default | `10000` | `10`      | `1000`       |
|goroutines  |
| Accuracy          | ⭐️⭐   | ⭐⭐⭐  | ⭐⭐⭐  | ⭐⭐⭐    | ⭐⭐⭐       |
| Max Stack Depth   | `64`    | `32`    | `32`    | `32`      | `32`         |

\- `100`
(3.) |

1. 一个 O(N) stop-the-world 的暂停，其中 N 是 goroutine 的数量。预计每个 goroutine 会有 ~1-10 微妙的暂停。
2. 完全坏了，不要尝试使用它。
3. 取决于 API 的情况。

### CPU Profiler

Go 的 CPU profiler 可以帮助你确定代码的哪些部分占用了大量的 CPU 时间。

⚠️ 请注意，CPU 时间跟我们体验的实际时间是不同的。例如，一个典型的 http 请求可能需要 `100ms` 才能完成，但是在数据库上等待 `95ms` 时只消耗 `5ms` 的 CPU 时间。一个请求也有可能 `100ms`，但是如果两个 goroutine 并行地执行 CPU 密集型工作，则需要花费 `200ms` 的 CPU。如果你对此感到困惑，请参阅 Goroutine Scheduler\[6] 部分。

你可以通过各种 APIs 来控制 CPU profiler:

- `go test -cpuprofile cpu.pprof` 将运行你的测试并将 CPU profile 写入名为 `cpu.pprof` 的文件。
- `pprof.StartCPUProfile(w)` 将 CPU profile 抓取到 `w`，涵盖的时间跨度直到 `pprof.StopCPUProfile()` 被调用。
- `import _ "net/http/pprof"` 允许你通过点击默认 http 服务器的 `GET /debug/pprof/profile?seconds=30` 来请求 30s CPU profile，你可以通过 `http.ListenAndServe("localhost:6060", nil)` 来启动。
- `runtime.SetCPUProfileRate()` 可以让你控制 CPU profile 的采样率。
- `runtime.SetCgoTraceback()` 可以用来获取 cgo 代码中的堆栈痕迹。benesch/cgosymbolizer\[7] 有一个针对 Linux 和 macOS 的实现。

如果你需要一个可以立马看到效果的代码贴到 `main()` 函数里，你可以使用下面的代码 :

`file, _ := os.Create("./cpu.pprof") pprof.StartCPUProfile(file) defer pprof.StopCPUProfile()`

无论如何激活 CPU profiler，最终的 profile 文件本质上都是一个以二进制 pprof\[8] 格式的堆栈跟踪表。这种表的简化版本如下 :

| stack trace  | samples/count | cpu/nanoseconds |
| ------------ | ------------- | --------------- |
| main;foo     | 5             | 50000000        |
| main;foo;bar | 3             | 30000000        |
| main;foobar  | 4             | 40000000        |

CPU profiler 通过请求操作系统监视应用程序的 CPU 使用情况来捕获这些数据，并为每占用 `10ms` 的 CPU 时间发送一个 `SIGPROF` 信号。操作系统还将内核代表应用程序所消耗的时间包括在这个监测中。由于信号传输速率取决于 CPU 的消耗，因此它是动态的，可以达到 `N * 100Hz`，其中 `N` 是系统上逻辑 CPU 核心的数量。当 `SIGPROF` 信号到达时，Go 的信号处理程序捕获当前活动 goroutine 的堆栈跟踪，并在 profile 中增加相应的值。`cpu/nanoseconds` 的值目前是直接从样本计数中推导出来的，因此它是冗余的，但很方便。

#### CPU Profiler 标签

Go 的 CPU profiler 的一个很吊的特性是你可以将任意键值对附加到 goroutine。这些标签将由从这个 goroutine 产生的任何 goroutine 继承，并显示在产生的 profile 文件中。

让我们考虑下面的示例\[9]，它代表 `user` 执行一些 CPU `work()`。通过使用 pprof.Labels()\[10]和 pprof.Labels()\[11] API，我们可以将 `user` 与执行 `work()` 函数的 goroutine 关联起来。此外，同一块代码中产生任何 goroutine 都会自动继承这些标签，例如 `backgroundWork()` goroutine。

`func work(ctx context.Context, user string) {  labels := pprof.Labels("user", user)  pprof.Do(ctx, labels, func(_ context.Context) {   go backgroundWork()   directWork()  }) }`

得到的 profile 将包括一个新的标签列，可能看起来像这样 :

| stack trace               | label      | samples/count | cpu/nanoseconds |
| ------------------------- | ---------- | ------------- | --------------- |
| main.childWork            | user:bob   | 5             | 50000000        |
| main.childWork            | user:alice | 2             | 20000000        |
| main.work;main.directWork | user:bob   | 4             | 40000000        |
| main.work;main.directWork | user:alice | 3             | 30000000        |

使用 pprof’s Graph 视图查看相同的档案也会包括以下标签 :

如何使用这些标签取决于你。你可以包含一些东西比如 `user ids`、`request ids`、 `http endpoints`, `subscription plan` 或其他数据，这些东西可以让你更好地理解哪些类型的请求导致了高 CPU 利用率，即使它们是由相同的代码路径处理的。也就是说，使用标签会增加 pprof 文件的大小。因此，你可能应该先从低 cardinality 标签（比如 endpoints）开始，一旦你确信它们不会影响你的应用程序的性能，就应该转向高 cardinality 标签。

#### CPU 利用率

由于 CPU profiler 的采样速率适应程序消耗的 CPU 数量，因此可以从 CPU profile 中获得 CPU 利用率。事实上，pprof 会自动为你做这件事。例如，下面的 profile 取自一个平均 CPU 利用率为 `147.77%` 的程序 :

`$ go tool pprof guide/cpu-utilization.pprof Type: cpu Time: Sep 9, 2021 at 11:34pm (CEST) Duration: 1.12s, Total samples = 1.65s (147.77%) Entering interactive mode (type "help" for commands, "o" for options) (pprof)`

另一种流行的表示 CPU 利用率的方法是 CPU 核心。在上面的例子中，程序在 profiling 期间平均使用了 `1.47` 个 CPU 核心。

⚠️ 如果这个数值接近或高于 `250%` ，那么你不应该过于信任这个数值。但是，如果你看到的数字非常低，比如 `10%` ，这通常表明 CPU 消耗对你的应用程序来说是小 case。一个常见的错误是忽略这个数字，并开始担心某个特定的函数占用了相对于 profile 文件其余部分的很长时间。当总体 CPU 利用率较低时，这通常是浪费时间，因为通过优化这个函数不会获得太多好处。

#### CPU Profiles 的系统调用

如果你看到诸如 `syscall.Read()` 或 `syscall.Write()` 这样的系统调用在你的 CPU profiles 中使用了大量的时间，请注意这只是内核中这些函数中占用的 CPU 时间。I/O 时间本身没有被跟踪。在系统调用上花费大量时间通常表明调用过多，因此增加缓冲区大小可能会有所帮助。对于这种更复杂的情况，你应该考虑使用 Linux perf，因为它还可以向你显示内核堆栈跟踪，从而可能为你提供更多的线索。

#### CPU Profiler 局限

有一些已知的问题和 CPU profiler 的局限性，需要注意的是 :

- 🐞 在 Linux 上的一个已知问题是 CPU profile 难以实现超过 `250Hz` 的采样率。这通常不是问题，但如果你的 CPU 利用率非常高，就会导致偏差。有关这方面的更多信息，可以看看 GitHub issue\[12]。同时你可以使用支持更高采样频率的 Linux perf。
- ⚠️️ 你可以在调用 runtime.SetCPUProfileRate()\[13] 之前调用 `runtime.StartCPUProfile()` 来调整 CPU profile 的速率。这将打印一个警告：`runtime: cannot set cpu profile rate until previous profile has finished`。然而，它仍然在上述 bug 的限制下工作。这个问题最初是在 这里\[14] 提出的，并且有一个 被接受的改进 API 的建议\[15]。
- ⚠️ 目前，CPU profile 可以在堆栈跟踪中捕获的嵌套函数调用的最大数量是 64\[16]。如果你的程序使用了大量的递归或其他导致调用函数堆栈的方法，你的 CPU profile 将包括堆栈跟踪被截断。这意味着你将错过调用链中导致采样时处于活动状态的函数的部分。

### Memory Profiler

Go memory(内存) profiler 可以帮助你识别代码中哪些部分执行了大量堆分配，以及在上一次垃圾收集期间有多少分配是仍可访问的。因此，memory profiler 生成的 profile 通常也称为 heap(堆) profile。

堆内存管理相关的活动通常占用 Go 进程消耗的 CPU 时间的 20-30% 。此外，由于减少了垃圾收集器扫描堆时发生的缓存抖动，干掉堆分配会产生二阶效应，从而加快代码的其他部分。这意味着优化内存分配通常比优化程序中与 cpu 绑定的代码路径更好使。

⚠️memory profiler 不显示堆栈分配，因为这些分配通常比堆分开销小得多。有关详细信息，请参阅本指南的 GC 章节。

你可以通过各种 api 来控制 memory profiler:

- `go test -memprofile mem.pprof` 将运行你的测试并将 memory profile 写进 `mem.pprof`。
- `pprof.Lookup("allocs").WriteTo(w, 0)` 将一个涵盖进程开始以来的时间的 memory profile 写入到 `w`。
- `import _ "net/http/pprof"` 允许你通过点击默认的 http 服务器 `GET /debug/pprof/allocs?seconds=30` 来请求 30 秒的 memory profile，你可以通过 `http.ListenAndServe("localhost:6060", nil)` 启动。这在内部也被称为 delta profile 。
- `runtime.MemProfileRate` 允许你控制 memory profilee 的采样率。

如果你需要一个可以立马看到效果的代码贴到 `main()` 函数里，你可以使用下面的代码 :

`file, _ := os.Create("./mem.pprof") defer pprof.Lookup("allocs").WriteTo(file, 0) defer runtime.GC()`

无论如何激活 memory profiler，得到的 profile 文件本质上都是一个二进制 pprof 格式的格式化的堆栈跟踪表。这种表的简化版本如下 :

| stack trace  | alloc_objects/count | alloc_space/bytes | inuse_objects/count | inuse_space/bytes |
| ------------ | ------------------- | ----------------- | ------------------- | ----------------- |
| main;foo     | 5                   | 120               | 2                   | 48                |
| main;foo;bar | 3                   | 768               | 0                   | 0                 |
| main;foobar  | 4                   | 512               | 1                   | 128               |

memory profile 包含两条主要信息 :

- `alloc_*`: 程序从进程最开始以来所有的分配。
- `insue_*`: 程序从上次 GC 完到现在分配。

你可以将此信息用于各种用途。例如，你可以使用 `alloc_*` 数据来确定哪些代码路径可能会产生大量 GC 处理，并且随着时间的推移查看 `inuse_*` 数据可以帮助你调查程序的内存泄漏或内存使用率过高。

#### Allocs vs Heap Profile

pprof.Lookup()\[17] 函数以及 net/http/pprof\[18] 包使用两个名称 (`allocs` 和 `heap`) 公开 memory profile。两个 profile 都包含相同的数据，唯一的区别是 `allocs` profile 将 `alloc_space/bytes` 设置为默认的样本类型，而 `heap` profile 默认设置为 `inuse_space/bytes`。pprof 工具使用它来决定默认情况下显示哪个示例类型。

#### Memory Profiler 采样

为了保持比较低的开销，内存 profiler 使用泊松采样，因此平均每 `512KiB` 只有一次分配触发堆栈跟踪并将其添加到 profile 中。然而，在将 profile 写入最终的 pprof 文件之前，runtime 将收集到的样本值除以抽样概率，从而对其进行扩展。这意味着报告的分配数量应该是对实际分配数量的很好的估算。而不管你使用的是 runtime.MemProfileRate\[19]。

对于生产中的分析，通常不必修改取样速率。这样做的唯一原因是，如果你担心在很少进行分配的情况下没有收集到足够的样本。

#### Memory Inuse vs RSS

一个常见的混淆是查看 `inuse_space/bytes` 样本类型的内存总量，并注意到它与操作系统报告的 RSS 内存使用量不匹配。有多种原因可能导致 :

- 根据定义，RSS 包括了很多不仅仅是 Go 堆内存的使用，例如 goroutine 堆栈、程序可执行文件、共享库以及 C 函数分配的内存所使用的内存。
- GC 可能不会立把空闲内存返回给操作系统，但在 Go 1.16 的 runtime\[20] 改回来了，这个就算小问题了。
- Go 使用 non-moving GC，因此在某些情况下，空闲堆内存可能会碎片化，从而阻止 Go 将其释放到操作系统。

#### Memory Profiler Implementation

下面的伪代码应该捕获 memory profiler 实现的基本东西，让你有个大概的印象。如下所示，Go runtime 中的 `malloc()` 函数使用 `poisson_sample(size)` 来确定是否应该对分配进行取样。如果是，则以堆栈跟踪 `s` 作为 `mem_profile` 哈希表中的 key，用来增加 `allocs` 和 `alloc_bytes` 计数器。此外，`track_profiled(object, s)` 调用将 `object` 标记为堆上的采样分配，并将堆栈跟踪 `s` 与其关联。

\`func malloc(size):
  object = ... // allocation magic

if poisson_sample(size):
    s = stacktrace()
    mem_profile\[s].allocs++
    mem_profile\[s].alloc_bytes += size
    track_profiled(object, s)

return object

\`

当 GC 确定是时候释放分配的对象时，它调用 `sweep()` ，使用 `is_profiled(object)` 来检查 `object` 是否被标记为采样对象。如果是，它将检索导致分配的堆栈跟踪 `s`，并在 `mem_profile` 内为它增加 `frees` 和 `free_bytes` 计数器。

\`func sweep(object):
  if is_profiled(object)
    s = alloc_stacktrace(object)
    mem_profile\[s].frees++
    mem_profile\[s].free_bytes += sizeof(object)

// deallocation magic

\`

`free_*` 计数器本身不包含在最终 memory profile 中。相反，它们通过简单的 `allocs - frees` 减法计算 profile 中的 `insue_*` 计数器。另外，最终的输出值是通过取样概率除以它们的比例而得到的。

#### Memory Profiler 的局限

memory profiler 有一些已知的问题还有局限性，需要注意的是 :

- ⚠️ `runtime.MemProfileRate` 必须尽可能早地在程序启动时修改一次，例如在 `main()` 的开头。写入这个值本质上是一个产生的数据竞争很小，在程序执行期间多次更改它会产生不正确的配置文件。
- ⚠ 在调试潜在的内存泄漏时，memory profiler 可以显示这些分配的是哪里创建的，但它无法告诉你哪些指针还在保持引用。这么些年，一直有人想解决掉这个问题，但没有一个适用于最新版本的 Go。如果你知道有什么好使的工具，请告诉 我\[21]。
- ⚠ CPU Profiler Labels\[22] 或者其他类似的东西不受 memory profiler 支持。在目前的实现中很难添加这个功能，因为它可能会在保存 memory profiler 数据的内部哈希映射中造成内存泄漏。
- ⚠ cgo C 代码所做的分配不会显示在 memory profile 文件中。
- ⚠ Memory profile 可能是两个垃圾收集周期前的数据。如果你想要一个更一致的时间点快照，可以考虑在请求内存配置文件之前调用 `runtime.GC()` 。net/http/pprof\[23] 接受 `?gc=1` 的参数。更多信息请参阅 runtime.MemProfile()\[24] 文档， 以及 `mprof.go` 中关于 `memRecord` 的注释。
- ⚠️ memory profiler 可以在堆栈跟踪中捕获的嵌套函数调用的最大数量目前是 `32`, 有关超过此限制时会发生什么情况的更多信息，请参阅 CPU Profiler Limitations\[25]。
- ⚠️ 保存 memory profile 文件的内部哈希表没有大小限制。这意味着它的大小会不断增长，直到它涵盖您的代码库中的所有分配代码路径。这在实践中不是问题，但如果您正在观察进程的内存使用情况，它可能看起来像一个比较小的内存泄漏。

### ThreadCreate Profiler

🐞 Threadcreate profile 旨在显示导致创建新 OS 线程的堆栈跟踪。然而，它从 2013 年\[26]就已经不好使了，所以大家记得不要用。

## 免责声明

原作者是 felixge\[27]，在 Datadog\[28] 做 Go 的 Continuous Profiling\[29]。同时公司也在招人\[30] : ).

欢迎对此指南\[31]进行反馈！

### 引用链接

\[1]

twitter: [_https://twitter.com/felixge/status/1435537024388304900_](https://twitter.com/felixge/status/1435537024388304900)

\[2]

所有的模型都是错误的: [_https://en.wikipedia.org/wiki/All_models_are_wrong_](https://en.wikipedia.org/wiki/All_models_are_wrong)

\[3]

Scheduling in Go: [_https://www.ardanlabs.com/blog/2018/08/scheduling-in-go-part1.html_](https://www.ardanlabs.com/blog/2018/08/scheduling-in-go-part1.html)

\[4]

伸缩性的关键因素: [_https://golang.org/doc/faq#goroutines_](https://golang.org/doc/faq#goroutines)

\[5]

Getting to Go: The Journey of Go's Garbage Collector: [_https://go.dev/blog/ismmkeynote_](https://go.dev/blog/ismmkeynote)

\[6]

Goroutine Scheduler: [_https://github.com/DataDog/go-profiler-notes/blob/main/guide/README.md#goroutine-scheduler_](https://github.com/DataDog/go-profiler-notes/blob/main/guide/README.md#goroutine-scheduler)

\[7]

benesch/cgosymbolizer: [_https://github.com/benesch/cgosymbolizer_](https://github.com/benesch/cgosymbolizer)

\[8]

pprof: [_https://github.com/DataDog/go-profiler-notes/blob/main/pprof.md_](https://github.com/DataDog/go-profiler-notes/blob/main/pprof.md)

\[9]

示例: [_https://github.com/DataDog/go-profiler-notes/blob/main/guide/cpu-profiler-labels.go_](https://github.com/DataDog/go-profiler-notes/blob/main/guide/cpu-profiler-labels.go)

\[10]

pprof.Labels(): [_https://pkg.go.dev/runtime/pprof#Labels_](https://pkg.go.dev/runtime/pprof#Labels)

\[11]

pprof.Labels(): [_https://pkg.go.dev/runtime/pprof#Labels_](https://pkg.go.dev/runtime/pprof#Labels)

\[12]

GitHub issue: [_https://github.com/golang/go/issues/35057_](https://github.com/golang/go/issues/35057)

\[13]

runtime.SetCPUProfileRate(): [_https://pkg.go.dev/runtime#SetCPUProfileRate_](https://pkg.go.dev/runtime#SetCPUProfileRate)

\[14]

这里: [_https://github.com/golang/go/issues/40094_](https://github.com/golang/go/issues/40094)

\[15]

被接受的改进 API 的建议: [_https://github.com/golang/go/issues/42502_](https://github.com/golang/go/issues/42502)

\[16]

64: [_https://sourcegraph.com/search?q=context:global+repo:github.com/golang/go+file:src/_](https://sourcegraph.com/search?q=context:global+repo:github.com/golang/go+file:src/)+maxCPUProfStack+%3D\&patternType=literal\*

\[17]

pprof.Lookup(): [_https://pkg.go.dev/runtime/pprof#Lookup_](https://pkg.go.dev/runtime/pprof#Lookup)

\[18]

net/http/pprof: [_https://pkg.go.dev/net/http/pprof_](https://pkg.go.dev/net/http/pprof)

\[19]

runtime.MemProfileRate: [_https://pkg.go.dev/runtime#MemProfileRate_](https://pkg.go.dev/runtime#MemProfileRate)

\[20]

Go 1.16 的 runtime: [_https://golang.org/doc/go1.16#runtime_](https://golang.org/doc/go1.16#runtime)

\[21]

我: [_https://github.com/DataDog/go-profiler-notes/issues_](https://github.com/DataDog/go-profiler-notes/issues)

\[22]

CPU Profiler Labels: _#cpu-profiler-labels_

\[23]

net/http/pprof: [_https://pkg.go.dev/net/http/pprof_](https://pkg.go.dev/net/http/pprof)

\[24]

runtime.MemProfile(): [_https://pkg.go.dev/runtime#MemProfile_](https://pkg.go.dev/runtime#MemProfile)

\[25]

CPU Profiler Limitations: _#cpu-profiler-limitations_

\[26]

2013 年: [_https://github.com/golang/go/issues/6104_](https://github.com/golang/go/issues/6104)

\[27]

felixge: [_https://github.com/felixge_](https://github.com/felixge)

\[28]

Datadog: [_https://www.datadoghq.com/_](https://www.datadoghq.com/)

\[29]

Continuous Profiling: [_https://www.datadoghq.com/product/code-profiling/_](https://www.datadoghq.com/product/code-profiling/)

\[30]

招人: [_https://www.datadoghq.com/jobs-engineering/#all\&all_locations_](https://www.datadoghq.com/jobs-engineering/#all&all_locations)

\[31]

此指南: [_https://github.com/DataDog/go-profiler-notes/blob/main/guide/README.md_](https://github.com/DataDog/go-profiler-notes/blob/main/guide/README.md)
