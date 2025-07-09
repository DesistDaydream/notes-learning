---
title: Go 网络编程
---

# 概述

> 参考：
> 
> - [Go 标准库 ，net](https://pkg.go.dev/net)
> - [公众号，开发内功修炼-在 golang 中是如何对 epoll 进行封装的？](https://mp.weixin.qq.com/s/hjWhh_zHfxmH1yZFfvu_zA)(关于 go 实现 net 的底层逻辑分析)
> - [Go 标准库，net/url](https://pkg.go.dev/net/url)(URL 解析与转译)
>   - [公众号-马哥 Linux 运维，go 标准库 net/url 学习笔记](https://mp.weixin.qq.com/s/p4F3lv_DBmWEwbj9v8273Q)

在协程没有流行以前，传统的网络编程中，同步阻塞是性能低下的代名词，一次切换就得是 [3 us](https://mp.weixin.qq.com/s/uq5s5vwk5vtPOZ30sfNsOg)  左右的 CPU 开销。各种基于 epoll 的异步非阻塞的模型虽然提高了性能，但是基于回调函数的编程方式却非常不符合人的的直线思维模式。开发出来的代码的也不那么容易被人理解。

Golang 的出现，可以说是将协程编程模式推向了一个高潮。这种新的编程方式既兼顾了同步编程方式的简单易用，也在底层通过协程和 epoll 的配合避免了线程切换的性能高损耗。换句话说就是既简单易用，性能又还不挺错。

## net 包

net 包中包含如下几个包

- **http** # http 包提供 HTTP 客户端和服务端的实现。
- **mail** # Package mail implements parsing of mail messages.
- **netip** # Package netip defines an IP address type that's a small value type.
- **rpc** # Package rpc provides access to the exported methods of an object across a network or other I/O connection.
- **smtp** # Package smtp implements the Simple Mail Transfer Protocol as defined in RFC 5321.
- **textproto** # Package textproto implements generic support for text-based request/response protocols in the style of HTTP, NNTP, and SMTP.
- **url** # 解析 URL 并实现查询转义

这些包基于 net，实现了更加抽象的能力，以便我们可以直接调用。

# net 的使用方式

考虑到不少读者没有使用过 golang，那么开头我先把一个基于官方 net 包的 golang 服务的简单使用代码给大家列出来。为了方便大家理解，我只保留骨干代码。

```go
package main

import (
	"net"
	"log"
)

// 处理连接
func handleConn(conn net.Conn) {
	defer conn.Close()
	// 定义缓冲区
	buf := make([]byte, 1024)

	// 读取客户端数据
	conn.Read(buf[:1024])

	// 将数据写回客户端
	len, err := conn.Write([]byte("hello,i am server"))
}

func main() {
	// 实例化监听器
	listener, err := net.Listen("tcp", ":8080")

	// 监听并接受连接
	for {
		// 等待客户端连接
		conn, err := listener.Accept()

		// 创建goroutine处理客户端连接
		go handleConn(conn)
	}
}
```

在这个示例服务程序中，先是使用 net.Listen 来监听了本地的 9008 这个端口。然后调用 Accept 进行接收连接处理。如果接收到了连接请求，通过 go process   来启动一个协程进行处理。在连接的处理中我展示了读写操作（Read 和 Write）。

整个服务程序看起来，妥妥的就是一个同步模型，包括 Accept、Read 和 Write 都会将当前协程给“阻塞”掉。比如 Read 函数这里，如果服务器调用时客户端数据还没有到达，那么 Read 是不带返回的，会将当前的协程 park 住。直到有了数据 Read 才会返回，处理协程继续执行。

你如果在其它语言，例如 C 和 Java 中写出这样类似的服务器代码，估计会被打死的。因为每一次同步的 Accept、Read、Write 都会导致你当前的线程被阻塞掉，会浪费大量的 CPU 进行线程上下文的切换。

但是在 golang 中这样的代码运行性能却是非常的不错，为啥呢？我们继续看本文接下来的内容。

## Listen 底层过程

在传统的 C、Java 等传统语言中，listen 所做的事情就是直接调用内核的 listen 系统调用。参见[《为什么服务端程序都需要先 listen 一下？》](https://mp.weixin.qq.com/s?__biz=MjM5Njg5NDgwNA==&mid=2247485737&idx=1&sn=baba45ad4fb98afe543bdfb06a5720b8&scene=21#wechat_redirect)。但是如果你也这么同等地理解 golang net 包里的 Listen， 那可就大错特错了。

和其它语言不同，在 golang net 的 listen 中，会完成如下几件事：

- 创建 socket 并设置非阻塞，
- bind 绑定并监听本地的一个端口
- 调用 listen 开始监听
- epoll_create 创建一个 epoll 对象
- epoll_etl 将 listen 的 socket 添加到 epoll 中等待连接到来

一次 Golang 的 Listen 调用，相当于在 C 语言中的 socket、bind、listen、epoll_create、epoll_etl 等多次函数调用的效果。封装度非常的高，更大程度地对程序员屏蔽了底层的实现细节。

> 插一句题外话：现在的各种开发工具的封装程度越来越高，真不知道对码农来说是好事还是坏事。好处是开发效率更高了，缺点是将来的程序员想了解底层也越来越难了，越来越像传统企业里流水线上的工人。

口说无凭，我们挖开 Golang 的内部源码瞅一瞅，这样更真实。

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/xmy6hv/1649232340561-6df2b65d-4643-4a1a-abe1-4cb6203a6fd1.png)

Listen 的入口在 golang 源码的 net/dial.go 文件中，让我们展开来看更细节的逻辑。

### Listen 入口执行流程

源码不用细看，看懂大概流程就可以。

```go
//file:go1.14.4/src/net/dial.go
func Listen(network, address string) (Listener, error) {
 var lc ListenConfig
 return lc.Listen(context.Background(), network, address)
}
```

可见，这个 Listen 只是一个入口。接下来会进入到 ListenConfig 下的 Listen 方法中。在 ListenConfig 的 Listen 中判断这是一个 TCP 类型的话，会进入到 sysListener 下的 listenTCP 方法里（src/net/tcpsock_posix.go）。然后再经过两三次的函数调用跳转，会进入到 net/sock_posix.go 文件下的 socket 函数中。我们直接看它。

```go
//file:go1.14.4/src/net/sock_posix.go
func socket(ctx context.Context, net string, family, ...) (fd *netFD, err error) {
 //创建 socket，见 2.2 小节
 s, err := sysSocket(family, sotype, proto)

 ...

 //TCP 绑定和监听，见 2.3 小节
 //epoll对象的创建以及文件描述符的添加 见 2.4 小节
 if laddr != nil && raddr == nil {
  switch sotype {
  case syscall.SOCK_STREAM, syscall.SOCK_SEQPACKET:
   fd.listenStream(laddr, listenerBacklog(), ctrlFn);
  ......
 }
}
```

接下来我们分别在 2.2 和 2.3 小节来介绍 sysSocket 和 listenStream 这两个函数。

### 创建 socket

sysSocket 这个函数和其它语言中的 socket 函数有很大的不同。在这个一个函数内就完成了三件事，创建 socket、bind 和 listen 监听。我们来看 sysSocket 的具体代码。

```go
//file:net/sys_cloexec.go
func sysSocket(family, sotype, proto int) (int, error) {
 //创建 socket
 s, err := socketFunc(family, sotype, proto)

 //设置为非阻塞模式
 syscall.SetNonblock(s, true)
}
```

在 sysSocket 中，调用的 socketFunc 其实就是 socket 系统调用。见如下代码。

```go
//file:net/hook_unix.go
var (
 // Placeholders for socket system calls.
 socketFunc        func(int, int, int) (int, error)  = syscall.Socket
 connectFunc       func(int, syscall.Sockaddr) error = syscall.Connect
 listenFunc        func(int, int) error              = syscall.Listen
 getsockoptIntFunc func(int, int, int) (int, error)  = syscall.GetsockoptInt
)
```

创建完 socket 之后，再调用 syscall.SetNonblock 将其设置为非阻塞模式。

```go
//file:syscall/exec_unix.go
func SetNonblock(fd int, nonblocking bool) (err error) {
 ...
 if nonblocking {
  flag |= O_NONBLOCK
 }
 fcntl(fd, F_SETFL, flag)
}
```

### 绑定和监听

我们接着再来看 listenStream。这个函数一进来就调用了系统调用 bind 和 listen 来完成了绑定和监听。

```go
//file:net/sock_posix.go
func (fd *netFD) listenStream(laddr sockaddr,...) error
{
 ...

 //等同于 c 语言中的：bind(listenfd, ...)
 syscall.Bind(fd.pfd.Sysfd, lsa);

 //等同于 c 语言中的：listen(listenfd, ...)
 listenFunc(fd.pfd.Sysfd, backlog);

 //这里非常关键：初始化socket与异步IO相关的内容
 if err = fd.init(); err != nil {
  return err
 }
}
```

其中 listenFunc 是一个宏，指向的就是 syscall.Listen 系统调用

```go
//file:go1.14.4/src/net/hook_unix.go
import "syscall"
var (
    // Placeholders for socket system calls.
    socketFunc        func(int, int, int) (int, error)  = syscall.Socket
    connectFunc       func(int, syscall.Sockaddr) error = syscall.Connect
    listenFunc        func(int, int) error              = syscall.Listen
    getsockoptIntFunc func(int, int, int) (int, error)  = syscall.GetsockoptInt
)
```

### epoll 创建和初始化

接下来在 fd.init 这一行，经过多次的函数调用展开以后会执行到 epoll 对象的创建，并还把在 listen 状态的 socket 句柄添加到了 epoll 对象中来管理其网络事件。

我们来看它是如何完成的。

```go
//file:go1.14.4/src/internal/poll/fd_poll_runtime.go
func (pd *pollDesc) init(fd *FD) error {
    serverInit.Do(runtime_pollServerInit)
    ctx, errno := runtime_pollOpen(uintptr(fd.Sysfd))
    ...
    return nil
}
```

serverInit.Do 这个是用来保证参数内的函数只执行一次的。不过多展开介绍。其参数 runtime_pollServerInit 是对 runtime 包的函数 poll_runtime_pollServerInit 的调用，其源码位于 runtime/netpoll.go 下。

```go
//file:runtime/netpoll.go
//go:linkname poll_runtime_pollServerInit internal/poll.runtime_pollServerInit
func poll_runtime_pollServerInit() {
 netpollGenericInit()
}
```

该函数会执行到 netpollGenericInit， epoll 就是在它的内部创建的。

```go
//file:netpoll_epoll.go
func netpollinit() {
 // epoll 对象的创建
 epfd = epollcreate1(_EPOLL_CLOEXEC)
 ...
}
```

再来看 runtime_pollOpen。它的参数就是前面 listen 好了的 socket 的文件描述符。在这个函数里，它将被放到 epoll 对象中。

```go
//file:runtime/netpoll_epoll.go
//go:linkname poll_runtime_pollOpen internal/poll.runtime_pollOpen
func poll_runtime_pollOpen(fd uintptr) (*pollDesc, int) {
 ...
 errno = netpollopen(fd, pd)
 return pd, int(errno)
}

//file:runtime/netpoll_epoll.go
func netpollopen(fd uintptr, pd *pollDesc) int32 {
 var ev epollevent
 ev.events = _EPOLLIN | _EPOLLOUT | _EPOLLRDHUP | _EPOLLET
 *(**pollDesc)(unsafe.Pointer(&ev.data)) = pd

 // listen 状态的 socket 被添加到了 epoll 中。
 return -epollctl(epfd, _EPOLL_CTL_ADD, int32(fd), &ev)
}
```

## Accept 过程

服务端在 Listen 完了之后，就是对 Accept 的调用了。该函数主要做了三件事

- 调用 accept 系统调用接收一个连接
- 如果没有连接到达，把当前协程阻塞掉
- 新连接到来的话，将其添加到 epoll 中管理，然后返回

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/xmy6hv/1649232340549-0d23de92-f60a-4b12-980a-f932d7a4e555.png)

通过 Golang 里的单步调试可以看到它进入到了 TCPListener 下的 Accept 里了。

```go
//file: net/tcpsock.go
func (l *TCPListener) Accept() (Conn, error) {
 c, err := l.accept()
 ...
}
func (ln *TCPListener) accept() (*TCPConn, error) {
 //以 netFD 的形式返回一个新连接
 fd, err := ln.fd.accept()
}
```

我们上面说的三步都是在 netFD 的 accept 函数里处理的。

```go
//file:net/fd_unix.go
func (fd *netFD) accept() (netfd *netFD, err error) {
 //3.1 接收一个连接
 //3.2 如果连接没有到达阻塞当前协程
 d, rsa, errcall, err := fd.pfd.Accept()

 //3.2 将新到的连接也添加到 epoll 中进行管理
 netfd, err = newFD(d, fd.family, fd.sotype, fd.net);
 netfd.init();

 ...
 return netfd, nil
}
```

接下来我们详细看每一步的细节。

### 接收一个连接

经过单步跟踪后发现 Accept 进入到了 FD 对象的 Accept 方法下。在这里将调用操作系统的 accept 系统调用。

```go
//file:internal/poll/fd_unix.go
// Accept wraps the accept network call.
func (fd *FD) Accept() (int, syscall.Sockaddr, string, error) {

 for {
  //调用 accept 系统调用接收一个连接
  s, rsa, errcall, err := accept(fd.Sysfd)

  //接收到了连接就返回它
  if err == nil {
   return s, rsa, "", err
  }

  switch err {
  case syscall.EAGAIN:
   //如果没有获取到，那就把协程给阻塞起来
   if fd.pd.pollable() {
    if err = fd.pd.waitRead(fd.isFile); err == nil {
     continue
    }
   }
  ...
 }
 ...
}
```

其中 accept 方法内部会触发 linux 操作系统的 accept 系统调用，我们就不过度展开了。调用 accept 目的是获取一个来自客户端的连接。如果接收到了，就把他返回回去。

### 阻塞当前协程

我们来说说如果没 accept 调用的时候，客户端的连接请求还一个都没有过来怎么办。

这时候，accept 系统调用会返回 syscall.EAGAIN。Golang 在对这个状态的处理中，会把当前协程给阻塞起来。关键代码在这里

```go
//file: internal/poll/fd_poll_runtime.go
func (pd *pollDesc) waitRead(isFile bool) error {
 return pd.wait('r', isFile)
}
func (pd *pollDesc) wait(mode int, isFile bool) error {
 if pd.runtimeCtx == 0 {
  return errors.New("waiting for unsupported file type")
 }
 res := runtime_pollWait(pd.runtimeCtx, mode)
 return convertErr(res, isFile)
}
```

runtime_pollWait 的源码在 runtime/netpoll.go 下。gopark（协程的阻塞）就是在这里完成的。

```go
//file:runtime/netpoll.go
//go:linkname poll_runtime_pollWait internal/poll.runtime_pollWait
func poll_runtime_pollWait(pd *pollDesc, mode int) int {
    ...
    for !netpollblock(pd, int32(mode), false) {
    }
}

func netpollblock(pd *pollDesc, mode int32, waitio bool) bool {
    ...
    if waitio || netpollcheckerr(pd, mode) == 0 {
        gopark(netpollblockcommit, unsafe.Pointer(gpp), waitReasonIOWait, traceEvGoBlockNet, 5)
    }
}
```

gopark 这个函数就是 golang 内部阻塞协程的入口。

### 将新连接添加到 epoll 中。

我们再来说说假如客户端连接已经到来了的情况。这时 fd.pfd.Accept 会返回新建的连接。然后会将该新连接也一并加入到 epoll 中进行高效的事件管理。

```go
//file:net/fd_unix.go
func (fd *netFD) accept() (netfd *netFD, err error) {
 //3.1 接收一个连接
 //3.2 如果连接没有到达阻塞当前协程
 d, rsa, errcall, err := fd.pfd.Accept()

 //3.2 将新到的连接也添加到 epoll 中进行管理
 netfd, err = newFD(d, fd.family, fd.sotype, fd.net);
 netfd.init();

 ...
 return netfd, nil
}
```

我们来看 netfd.init

```go
//file:internal/poll/fd_poll_runtime.go
func (pd *pollDesc) init(fd *FD) error {
 ...
 ctx, errno := runtime_pollOpen(uintptr(fd.Sysfd))
 ...
}
```

runtime_pollOpen 这个 runtime 函数我们在上面的 2.4 节介绍过了，就是把文件句柄添加到 epoll 对象中。

```go
//file:runtime/netpoll_epoll.go
//go:linkname poll_runtime_pollOpen internal/poll.runtime_pollOpen
func poll_runtime_pollOpen(fd uintptr) (*pollDesc, int) {
 ...
 errno = netpollopen(fd, pd)
 return pd, int(errno)
}

func netpollopen(fd uintptr, pd *pollDesc) int32 {
 var ev epollevent
 ev.events = _EPOLLIN | _EPOLLOUT | _EPOLLRDHUP | _EPOLLET
 *(**pollDesc)(unsafe.Pointer(&ev.data)) = pd

 //新连接的 socket 也被添加到了 epoll 中。
 return -epollctl(epfd, _EPOLL_CTL_ADD, int32(fd), &ev)
}
```

## Read 和 Write 内部过程

当连接接收完成后，剩下的就是在连接上的读写了。

### Read 内部过程

我们先来看 Read 大体过程。
![](https://notes-learning.oss-cn-beijing.aliyuncs.com/xmy6hv/1649232340624-b95d7cbb-46bf-4d19-85ef-456023ca2e73.png)
来看详细的代码。

```go
//file:/Users/zhangyanfei/sdk/go1.14.4/src/net/net.go
func (c *conn) Read(b []byte) (int, error) {
 ...
 n, err := c.fd.Read(b)
}
```

Read 函数会进入到 FD 的 Read 中。在这个函数内部调用 Read 系统调用来读取数据。如果数据还尚未到达则也是把自己阻塞起来。

```go
//file:internal/poll/fd_unix.go
func (fd *FD) Read(p []byte) (int, error) {
 for {
  //调用 Read 系统调用
  n, err := syscall.Read(fd.Sysfd, p)
  if err != nil {
   n = 0

   //将自己添加到 epoll 中等待事件，然后阻塞掉。
   if err == syscall.EAGAIN && fd.pd.pollable() {
    if err = fd.pd.waitRead(fd.isFile); err == nil {
     continue
    }
   }
  ......
 }
}
```

其中 waitRead 是如何将当前协程阻塞掉的，这个和我们前面 3.2 节介绍的是一样的，就不过多展开叙述了。

### Write 内部过程

Write 的大体过程和 Read 是类似的。先是调用 Write 系统调用发送数据，如果内核发送缓存区不足的时候，就把自己先阻塞起来，然后等可写时间发生的时候再继续发送。其源码入口位于 net/net.go。

```go
//file:net/net.go
func (c *conn) Write(b []byte) (int, error) {
 ...
 n, err := c.fd.Write(b)
}
```

```go
//file:internal/poll/fd_unix.go
func (fd *FD) Write(p []byte) (int, error) {
 for {
  n, err := syscall.Write(fd.Sysfd, p[nn:max])
  if err == syscall.EAGAIN && fd.pd.pollable() {
   if err = fd.pd.waitWrite(fd.isFile); err == nil {
    continue
   }
  }
 }
}
```

```go
//file:internal/poll/fd_poll_runtime.go
func (pd *pollDesc) waitWrite(isFile bool) error {
 return pd.wait('w', isFile)
}
```

pd.wait 之后的事情就又和 3.2 节介绍的过程一样了。调用 runtime_pollWait 来讲当前协程阻塞掉。

```go
func (pd *pollDesc) wait(mode int, isFile bool) error {
 ...
 res := runtime_pollWait(pd.runtimeCtx, mode)
}
```

## Golang 唤醒

前面我们讨论的很多步骤里都涉及到协程的阻塞。例如 Accept 时如果新连接还尚未到达。再比如像  Read 数据的时候对方还没有发送，当前协程都不会占着 cpu 不放，而是会阻塞起来。

那么当要等待的事件就绪的时候，被阻塞掉的协程又是如何被重新调度的呢？相信大家一定会好奇这个问题。

Go 语言的运行时会在调度或者系统监控中调用 sysmon，它会调用 netpoll，来不断地调用 epoll_wait 来查看 epoll 对象所管理的文件描述符中哪一个有事件就绪需要被处理了。如果有，就唤醒对应的协程来进行执行。

其实除此之外还有几个地方会唤醒协程，如

- startTheWorldWithSema
- findrunnable   在 schedule 中调用 有 top 和 stop 之分。其中 stop 中会导致阻塞。
- pollWork

不过为了简便起见，我们只选择 sysmon 来作为一个切入口。sysmon 是一个周期性的监控协程，来看源码。

```go
//file:src/runtime/proc.go
func sysmon() {
 ...
 list := netpoll(0)
}
```

它会不断触发对  netpoll 的调用，在 netpoll 会调用 epollwait 看查看是否有网络事件发生。

```go
//file:runtime/netpoll_epoll.go
func netpoll(delay int64) gList {
 ...
retry:
 n := epollwait(epfd, &events[0], int32(len(events)), waitms)
 if n < 0 {
  //没有网络事件
  goto retry
 }

 for i := int32(0); i < n; i++ {

  //查看是读事件还是写事件发生
  var mode int32
  if ev.events&(_EPOLLIN|_EPOLLRDHUP|_EPOLLHUP|_EPOLLERR) != 0 {
   mode += 'r'
  }
  if ev.events&(_EPOLLOUT|_EPOLLHUP|_EPOLLERR) != 0 {
   mode += 'w'
  }

  if mode != 0 {

   pd := *(**pollDesc)(unsafe.Pointer(&ev.data))
   pd.everr = false
   if ev.events == _EPOLLERR {
    pd.everr = true
   }
   netpollready(&toRun, pd, mode)
  }
 }
}
```

在 epoll 返回的时候，ev.data 中是就绪的网络 socket 的文件描述符。根据网络就绪 fd 拿到 pollDesc。在 netpollready 中，将对应的协程推入可运行队列等待调度执行。

```go
//file:runtime/netpoll.go
func netpollready(toRun *gList, pd *pollDesc, mode int32) {
 var rg, wg *g
 if mode == 'r' || mode == 'r'+'w' {
  rg = netpollunblock(pd, 'r', true)
 }
 if mode == 'w' || mode == 'r'+'w' {
  wg = netpollunblock(pd, 'w', true)
 }
 if rg != nil {
  toRun.push(rg)
 }
 if wg != nil {
  toRun.push(wg)
 }
}
```

## 本文总结

同步编码方式的优点是符合人的直线思维。在这种模式下的代码很容易写，写出来也容易理解，但是缺点就是性能奇差。因为会导致频繁的线程上下文切换。

所以现在 epoll 是 Linux 下网络程序工作的最主要的模式。现在各种语言下的流行的网络框架模型都是基于 epoll 来工作的。区别就是各自对 epoll 的使用方式上存在一些差别。主流各种基于 epoll 的异步非阻塞的模型虽然提高了性能，但是基于回调函数的编程方式却非常不符合人的的直线思维模式。开发出来的代码的也不那么容易被人理解。

Golang 开辟了一种新的网络编程模型。这种模型在应用层看来仍然是同步的方式。但是在底层确实通过协程和 epoll 的配合避免了线程切换的性能高损耗，因此并不会阻塞用户线程。代替的是切换开销更小的协程。协程的切换开销大约只有线程切换的三十分之一，参见[《协程究竟比线程牛在什么地方？》](https://mp.weixin.qq.com/s?__biz=MjM5Njg5NDgwNA==&mid=2247483805&idx=1&sn=3e62e6712335ee8520e5d525c078c110&scene=21#wechat_redirect)

我个人一直觉得，Golang 封装的网络编程模型非常之精妙，是世界级的代码。它非常值得你好好学习一下。学完了觉得好的话，转发给你的朋友们一起来了解了解吧！

### 往期相关文章

- [进程/线程切换究竟需要多少开销？](https://mp.weixin.qq.com/s?__biz=MjM5Njg5NDgwNA==&mid=2247483804&idx=1&sn=f2d64fc244d381157bb0c16ff26a33bd&scene=21#wechat_redirect)
- [协程究竟比线程牛在什么地方？](https://mp.weixin.qq.com/s?__biz=MjM5Njg5NDgwNA==&mid=2247483805&idx=1&sn=3e62e6712335ee8520e5d525c078c110&scene=21#wechat_redirect)
- [为什么服务端程序都需要先 listen 一下？](https://mp.weixin.qq.com/s?__biz=MjM5Njg5NDgwNA==&mid=2247485737&idx=1&sn=baba45ad4fb98afe543bdfb06a5720b8&scene=21#wechat_redirect)
- [图解 | 深入理解高性能网络开发路上的绊脚石 - 同步阻塞网络 IO](https://mp.weixin.qq.com/s?__biz=MjM5Njg5NDgwNA==&mid=2247484834&idx=1&sn=b8620f402b68ce878d32df2f2bcd4e2e&scene=21#wechat_redirect)
- [图解 | 深入揭秘 epoll 是如何实现 IO 多路复用的！](https://mp.weixin.qq.com/s?__biz=MjM5Njg5NDgwNA==&mid=2247484905&idx=1&sn=a74ed5d7551c4fb80a8abe057405ea5e&scene=21#wechat_redirect)
- [漫画 | 看进程小 P 讲述它的网络性能故事！](https://mp.weixin.qq.com/s?__biz=MjM5Njg5NDgwNA==&mid=2247485035&idx=1&sn=d126a429f6803f54a053e75723fac288&scene=21#wechat_redirect)

# Web 框架

[Gin](/docs/2.编程/高级编程语言/Go/Go%20网络编程/Gin.md)