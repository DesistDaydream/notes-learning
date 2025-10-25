---
title: Socket
linkTitle: Socket
weight: 20
---

# 概述

> 参考：
>
> - [Manual, socket(2)](https://man7.org/linux/man-pages/man2/socket.2.html)
> - [Manual, socket(7)](https://man7.org/linux/man-pages/man7/socket.7.html)
> - [GitHub 项目，torvalds/linux - include/linux/socket.h](https://github.com/torvalds/linux/blob/v6.14/include/linux/socket.h#L189)
> - [Wiki 搜索, Socket](https://en.wikipedia.org/wiki/Socket)

**Socket(套接字) 是两个实体间进行数据通信的 Endpoint(端点)**。是计算机领域中数据通信的一种约定，或者说是一种方法。通过 Socket 这种方法，计算机内的进程可以互相交互数据，不同计算机之间也可以互相交互数据。

Socket(套接字) 原意是 “插座”，所以 Socket 就像插座的作用一样，只要把插头插上，就能让设备获得电力。同理，只要两个程序通过 Socket 互相套接，也就是说两个程序都插在同一个 Socket 上，那么这两个程序就能交互数据。

在软件上，Socket 负责套接计算机中的数据(可以想象成套接管，套接管即为套管，是用来把两个管连接起来的东西，套接字就是把计算机中的字(即最小数据)连接起来，且只把头部连接起来，套管也是，只把两根很长的管的头端套起来接上)

1. 在系统层面，Socket 可以连接系统中的两个进程，进程与进程本身是互相独立的，如果需要传递消息，那么就需要两个进程各自打开一个接口(API)，Socket 把两个进程的 API 套住使之连接起来，即可实现进程间的通信。该 Socket 是抽象的，虚拟的，只是通过编程函数来实现进程的 API 功能，如果进程没有 API，那么就无法通过 Socket 与其余进程通信。
2. 当然，一个进程也可以监听一个名为 .scok 的文件，这个文件就像 API 一样，其他程序想与该进程交互，只要指定该 .sock 文件，然后对这个 .sock 文件进行读写即可。
3. 在网络层面，Socket 负责把不在同一主机上的进程(比如主机 A 的进程 C 和主机 B 的进程 D)连接起来，而两个不同主机上的进程如何被套接起来呢，套接至少需要提供一个头端来让套接管(字)包裹住才行。这时候(协议，IP，端口,例如：ftp://192.168.0.1:22)共同组成了网络上的进程标示，该进程逻辑上的头端即为紫色部分的端口号，不同主机的两个进程可以通过套接字把端口号套起来连接，来使两个网络上不同主机的进程进行通信，该同能同样是在程序编程的时候用函数写好的，程序启动为进程的时候，则该接口会被拿出来监听。

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/socket/1619421243110-2db70bc6-f358-459c-b9a9-e199658b151a.png)

想要在 Linux 中使用 Socket，可以直接利用 socket 系统调用 `int socket(int domain, int type, int protocol);` 创建一个 Socket。

> 参考系统调用的 Manual，可以看到 Linux 将 Socket 进行了两种分类，一个 domain(family)，一个 type，protocol 通常忽略。

## Socket type 与 Socket Family

**Socket Family(族)**/**Socket Domain(域)** 用于定义 Socket 在哪里通信，是 本地、网络、etc. 。可以在 [address_families(7)](https://man7.org/linux/man-pages/man7/address_families.7.html) 找到所有定义了的 Families

- **[Unix Domain Socket](#Unix%20Domain%20Socket)** # 用于同一台设备的不同进程间互相通信
  - AF_UNIX
- **[Network Socket](#Network%20Socket)** # 用于进程在网络间互相通信
  - AF_INET,  AF_INET6
- **AF_NETLINK** # 内核空间与用户空间通信的 Socket。比如 [Netlink](/docs/1.操作系统/Kernel/Network/Linux%20网络栈管理/Netlink/Netlink.md)
- etc.

**Socket type** 用于约定约定一些 **communication semantics(通信语义)**：

> [!Tip] Socket type 用于定义 Socket 通信时，数据的传输机制。 从内核 2.6.27 版本开始，Socket type 还可以修改 Socket 的行为

> [!Note]  Socket family 可以支持一个或多个 Socket type。哪怕是 UDS 也可以利用本地 .sock 文件进行类似 TCP/UDP 的通信。

- **SOCK_STREAM** # 面向流的传输（e.g. TCP）。提供基于连接的、可靠的双向有序字节流传输，可能支持带外数据传输机制。
- **SOCK_DGRAM** # 面向消息的传输（e.g. UDP）。支持数据报（固定最大长度的无连接、不可靠消息）。
- **SOCK_SEQPACKET** # 面向记录的传输（e.g. SCTP）。为固定最大长度的数据报提供基于连接的、可靠的双向有序传输路径；每次输入系统调用需读取完整数据包。
- etc.
- 修改 Socket 行为的类型
  - **SOCK_NONBLOCK**
  - etc.

## Linux 的 socket 接口运行逻辑

https://man7.org/linux/man-pages/man7/socket.7.html

[socket(2)](https://man7.org/linux/man-pages/man2/socket.2.html) 创建一个套接字， [connect(2)](https://man7.org/linux/man-pages/man2/connect.2.html) 将套接字连接到 远程套接字地址， [bind(2)](https://man7.org/linux/man-pages/man2/bind.2.html) 函数将套接字绑定到 本地套接字地址， [listen(2)](https://man7.org/linux/man-pages/man2/listen.2.html) 告诉套接字有新的 连接将被接受，并且 [accept(2)](https://man7.org/linux/man-pages/man2/accept.2.html) 用于获取新的 套接字与新的传入连接。[socketpair(2)](https://man7.org/linux/man-pages/man2/socketpair.2.html) 返回两个 连接的匿名套接字（仅为少数本地实现 类似 **AF_UNIX** 的系列）

[send(2)](https://man7.org/linux/man-pages/man2/send.2.html) 、 [sendto(2)](https://man7.org/linux/man-pages/man2/sendto.2.html) 和 [sendmsg(2)](https://man7.org/linux/man-pages/man2/sendmsg.2.html) 通过套接字发送数据，并且 [recv(2)](https://man7.org/linux/man-pages/man2/recv.2.html) 、 [recvfrom(2)](https://man7.org/linux/man-pages/man2/recvfrom.2.html) 、 [recvmsg(2)](https://man7.org/linux/man-pages/man2/recvmsg.2.html) 从套接字接收数据。 [poll(2)](https://man7.org/linux/man-pages/man2/poll.2.html) 和 [select(2)](https://man7.org/linux/man-pages/man2/select.2.html) 等待数据到达或准备就绪 发送数据。此外，标准 I/O 操作，例如 [write(2)](https://man7.org/linux/man-pages/man2/write.2.html) 、 [writev(2)](https://man7.org/linux/man-pages/man2/writev.2.html) 、 [sendfile(2)](https://man7.org/linux/man-pages/man2/sendfile.2.html) 、 [read(2](https://man7.org/linux/man-pages/man2/read.2.html) 和 [readv(2)](https://man7.org/linux/man-pages/man2/readv.2.html) 可以是 用于读取和写入数据。

[getsockname(2)](https://man7.org/linux/man-pages/man2/getsockname.2.html) 返回本地套接字地址， [getpeername(2)](https://man7.org/linux/man-pages/man2/getpeername.2.html) 返回远程套接字地址。getsockopt [(2)](https://man7.org/linux/man-pages/man2/getsockopt.2.html) 和 [setsockopt(2)](https://man7.org/linux/man-pages/man2/setsockopt.2.html) 用于设置或获取套接字层或协议 选项。ioctl [(2)](https://man7.org/linux/man-pages/man2/ioctl.2.html) 可用于设置或读取一些其他选项。

[close(2)](https://man7.org/linux/man-pages/man2/close.2.html) 用于关闭套接字。shutdown [(2)](https://man7.org/linux/man-pages/man2/shutdown.2.html) 关闭套接字的部分 全双工套接字连接。

## Berkeley Sockets API

**Berkeley Sockets** 是 Network Socket 和 Unix Domain Sockets 的 [API](/docs/2.编程/API/API.md)），用于进程间通信（IPC）。通常将其实现为可链接模块的库。它起源于 1983 年发布的 4.2BSD Unix 操作系统。

套接字是网络通信路径的本地终结点的抽象表示（句柄）。Berkeley 套接字 API 将其表示为 Unix 哲学中的文件描述符（文件句柄），该描述符为输入和输出到数据流提供通用接口。

伯克利套接字几乎没有任何改动，从 _事实上的_ 标准演变为 POSIX 规范的组件。术语 **POSIX 套接字**在本质上是 _Berkeley 套接字的_ 同义词，但是它们也称为 BSD 套接字，这是对 Berkeley Software Distribution 中的第一个实现的认可。

# Socket Families

## Unix Domain Socket

详见 [Unix Domain Socket](/docs/1.操作系统/Kernel/Process/Inter%20Process%20Communication/Socket/Unix%20Domain%20Socket.md)

## Network Socket

详见 [Network Socket](/docs/4.数据通信/数据通信/Network%20Socket.md Socket.md)

# Unix Domain Socket 与 Network Socket 对比

原文：[公众号，这种本机网络 IO 方法，性能可以翻倍！](https://mp.weixin.qq.com/s/fHzKYlW0WMhP2jxh2H_59A)

原创 张彦飞allen _2021年12月21日 09:10_

大家好，我是飞哥！

很多读者在看完 [《127.0.0.1 之本机网络通信过程知多少?》](https://mp.weixin.qq.com/s?__biz=MjM5Njg5NDgwNA==&mid=2247485270&idx=1&sn=503534e9f0560bfcfbd4539e028e0d57&scene=21#wechat_redirect) 这一篇后，让我讲讲 Unix Domain Socket。好了，今天就安排！

在本机网络 IO 中，我们讲到过基于普通 socket 的本机网络通信过程中，其实在内核工作流上并没有节约太多的开销。该走的系统调用、协议栈、邻居系统、设备驱动（虽然说对于本机网络 loopback 设备来说只是一个软件虚拟的东东）全都走了一遍。其工作过程如下图

![图片](https://mmbiz.qpic.cn/mmbiz_png/BBjAFF4hcwqr5dxehUyQGUfk69ibqFibVMYBNJiclZk13tB111axoKtso0QfeeZMKPrLohFlPoNBAGTzhviay28ibicA/640?wx_fmt=png&tp=webp&wxfrom=5&wx_lazy=1&wx_co=1)

那么我们今天来看另外一种本机网络 IO 通信方式 -- Unix Domain Socket。看看这种方式在性能开销上和基于 127.0.0.1 的本机网络 IO 有没有啥差异呢。

本文中，我们将分析 Unix Domain Socket 的内部工作原理。你将理解为什么这种方式的性能比 127.0.0.1 要好很多。最后我们还给出了实际的性能测试对比数据。

相信你已经迫不及待了，别着急，让我们一一展开细说！

## 一、使用方法

Unix Domain Socket（后面统一简称 UDS） 使用起来和传统的 socket 非常的相似。区别点主要有两个地方需要关注。

第一，在创建 socket 的时候，普通的 socket 第一个参数 family 为 AF\_INET， 而 UDS 指定为 AF\_UNIX 即可。

第二，Server 的标识不再是 ip 和 端口，而是一个路径，例如 /dev/shm/fpm-cgi.sock。

其实在平时我们使用 UDS 并不一定需要去写一段代码，很多应用程序都支持在本机网络 IO 的时候配置。例如在 Nginx 中，如果要访问的本机 fastcgi 服务是以 UDS 方式提供服务的话，只需要在配置文件中配置这么一行就搞定了。

```
fastcgi_pass unix:/dev/shm/fpm-cgi.sock;
```

如果 对于一个 UDS 的 server 来说，它的代码示例大概结构如下，大家简单了解一下。只是个示例不一定可运行。

```c
int main()
{
 // 创建 unix domain socket
 int fd = socket(AF_UNIX, SOCK_STREAM, 0);

 // 绑定监听
 char *socket_path = "./server.sock";
 strcpy(serun.sun_path, socket_path);
 bind(fd, serun, ...);
 listen(fd, 128);

 while(1){
  //接收新连接
  conn = accept(fd, ...);

  //收发数据
  read(conn, ...);
  write(conn, ...);
 }
}
```

基于 UDS 的 client 也是和普通 socket 使用方式差不太多，创建一个 socket，然后 connect 即可。

```c
int main(){
 sock = socket(AF_UNIX, SOCK_STREAM, 0);
 connect(sockfd, ...)
}
```

## 二、连接过程

总的来说，基于 UDS 的连接过程比 inet 的 socket 连接过程要简单多了。客户端先创建一个自己用的 socket，然后调用 connect 来和服务器建立连接。

在 connect 的时候，会申请一个新 socket 给 server 端将来使用，和自己的 socket 建立好连接关系以后，就放到服务器正在监听的 socket 的接收队列中。这个时候，服务器端通过 accept 就能获取到和客户端配好对的新 socket 了。

总的 UDS 的连接建立流程如下图。

![图片](https://mmbiz.qpic.cn/mmbiz_png/BBjAFF4hcwqr5dxehUyQGUfk69ibqFibVMNB1r0ZJF5FRQNwMx6SJnUuN1NmM29TVDWlRibkCXuEIaYRGMkeAvibcA/640?wx_fmt=png&tp=webp&wxfrom=5&wx_lazy=1&wx_co=1)

内核源码中最重要的逻辑在 connect 函数中，我们来简单展开看一下。unix 协议族中定义了这类 socket 的所有方法，它位于 net/unix/af\_unix.c 中。

```c
//file: net/unix/af_unix.c
static const struct proto_ops unix_stream_ops = {
 .family = PF_UNIX,
 .owner = THIS_MODULE,
 .bind =  unix_bind,
 .connect = unix_stream_connect,
 .socketpair = unix_socketpair,
 .listen = unix_listen,
 ...
};
```

我们找到 connect 函数的具体实现，unix\_stream\_connect。

```c
//file: net/unix/af_unix.c
static int unix_stream_connect(struct socket *sock, struct sockaddr *uaddr,
          int addr_len, int flags)
{
 struct sockaddr_un *sunaddr = (struct sockaddr_un *)uaddr;

 ...

 // 1. 为服务器侧申请一个新的 socket 对象
 newsk = unix_create1(sock_net(sk), NULL);

 // 2. 申请一个 skb，并关联上 newsk
 skb = sock_wmalloc(newsk, 1, 0, GFP_KERNEL);
 ...

 // 3. 建立两个 sock 对象之间的连接
 unix_peer(newsk) = sk;
 newsk->sk_state  = TCP_ESTABLISHED;
 newsk->sk_type  = sk->sk_type;
 ...
 sk->sk_state = TCP_ESTABLISHED;
 unix_peer(sk) = newsk;

 // 4. 把连接中的一头（新 socket）放到服务器接收队列中
 __skb_queue_tail(&other->sk_receive_queue, skb);
}
```

主要的连接操作都是在这个函数中完成的。和我们平常所见的 TCP 连接建立过程，这个连接过程简直是太简单了。没有三次握手，也没有全连接队列、半连接队列，更没有啥超时重传。

直接就是将两个 socket 结构体中的指针互相指向对方就行了。就是 unix\_peer(newsk) = sk 和 unix\_peer(sk) \= newsk 这两句。

```c
//file: net/unix/af_unix.c
#define unix_peer(sk) (unix_sk(sk)->peer)
```

当关联关系建立好之后，通过 \_\_skb\_queue\_tail 将 skb 放到服务器的接收队列中。注意这里的 skb 里保存着新 socket 的指针，因为服务进程通过 accept 取出这个 skb 的时候，就能获取到和客户进程中 socket 建立好连接关系的另一个 socket。

怎么样，UDS 的连接建立过程是不是很简单！？

## 三、发送过程

看完了连接建立过程，我们再来看看基于 UDS 的数据的收发。这个收发过程一样也是非常的简单。发送方是直接将数据写到接收方的接收队列里的。

![](https://mmbiz.qpic.cn/mmbiz_png/BBjAFF4hcwqr5dxehUyQGUfk69ibqFibVMbUcIFoJibcDKRgIMeZqWiaV9gIjBYZicJPhBjib2W3ibmqrhzlHmHT3jmyw/640?wx_fmt=png&tp=webp&wxfrom=5&wx_lazy=1&wx_co=1)

我们从 send 函数来看起。send 系统调用的源码位于文件 net/socket.c 中。在这个系统调用里，内部其实真正使用的是 sendto 系统调用。它只干了两件简单的事情，

第一是在内核中把真正的 socket 找出来，在这个对象里记录着各种协议栈的函数地址。第二是构造一个 struct msghdr 对象，把用户传入的数据，比如 buffer地址、数据长度啥的，统统都装进去. 剩下的事情就交给下一层，协议栈里的函数 inet\_sendmsg 了，其中 inet\_sendmsg 函数的地址是通过 socket 内核对象里的 ops 成员找到的。大致流程如图。

![](https://mmbiz.qpic.cn/mmbiz_png/BBjAFF4hcwqr5dxehUyQGUfk69ibqFibVMghJ9jR02cFjY2FoJFpsN8yZsYibSPSHpS2hIOxXYExNkVfTkRLY8zSg/640?wx_fmt=png&tp=webp&wxfrom=5&wx_lazy=1&wx_co=1)

在进入到协议栈 inet\_sendmsg 以后，内核接着会找到 socket 上的具体协议发送函数。对于 Unix Domain Socket 来说，那就是 unix\_stream\_sendmsg。我们来看一下这个函数

```c
//file:
static int unix_stream_sendmsg(struct kiocb *kiocb, struct socket *sock,
          struct msghdr *msg, size_t len)
{
 // 1.申请一块缓存区
 skb = sock_alloc_send_skb(sk, size, msg->msg_flags&MSG_DONTWAIT,
      &err);

 // 2.拷贝用户数据到内核缓存区
 err = memcpy_fromiovec(skb_put(skb, size), msg->msg_iov, size);

 // 3. 查找socket peer
 struct sock *other = NULL;
 other = unix_peer(sk);

 // 4.直接把 skb放到对端的接收队列中
 skb_queue_tail(&other->sk_receive_queue, skb);

 // 5.发送完毕回调
 other->sk_data_ready(other, size);
}
```

和复杂的 TCP 发送接收过程相比，这里的发送逻辑简单简单到令人发指。申请一块内存（skb），把数据拷贝进去。根据 socket 对象找到另一端， **直接把 skb 给放到对端的接收队列里了**

接收函数主题是 unix\_stream\_recvmsg，这个函数中只需要访问它自己的接收队列就行了，源码就不展示了。所以在本机网络 IO 场景里，基于 Unix Domain Socket 的服务性能上肯定要好一些的。

## 四、性能对比

为了验证 Unix Domain Socket 到底比基于 127.0.0.1 的性能好多少，我做了一个性能测试。在网络性能对比测试，最重要的两个指标是延迟和吞吐。我从 Github 上找了个好用的测试源码：https://github.com/rigtorp/ipc-bench。我的测试环境是一台 4 核 CPU，8G 内存的 KVM 虚机。

在延迟指标上，对比结果如下图。

![](https://mmbiz.qpic.cn/mmbiz_png/BBjAFF4hcwqr5dxehUyQGUfk69ibqFibVM1X83FjHgcKPXoicpOStLm3vprwoH3tuhia91l2oibfmAOUhpCFjUaHlCA/640?wx_fmt=png&tp=webp&wxfrom=5&wx_lazy=1&wx_co=1)

可见在小包（100 字节）的情况下，UDS 方法的“网络” IO 平均延迟只有 2707 纳秒，而基于 TCP（访问 127.0.0.1）的方式下延迟高达 5690 纳秒。耗时整整是前者的两倍。

在包体达到 100 KB 以后，UDS 方法延迟 24 微秒左右（1 微秒等于 1000 纳秒），TCP 是 32 微秒，仍然高一截。这里低于 2 倍的关系了，是因为当包足够大的时候，网络协议栈上的开销就显得没那么明显了。

再来看看吞吐效果对比。

![](https://mmbiz.qpic.cn/mmbiz_png/BBjAFF4hcwqr5dxehUyQGUfk69ibqFibVMtaHJ8KD0jPibiaqiaYXSuWn4NNNhgEIwU9xsib8XxtMBibzHwn4z5eRHTNw/640?wx_fmt=png&tp=webp&wxfrom=5&wx_lazy=1&wx_co=1)

在小包的情况下，带宽指标可以达到 854 M，而基于 TCP 的 IO 方式下只有 386 M。数据就解读到这里。

## 五、总结

本文分析了基于 Unix Domain Socket 的连接创建、以及数据收发过程。其中数据收发的工作过程如下图。

![](https://mmbiz.qpic.cn/mmbiz_png/BBjAFF4hcwqr5dxehUyQGUfk69ibqFibVMbUcIFoJibcDKRgIMeZqWiaV9gIjBYZicJPhBjib2W3ibmqrhzlHmHT3jmyw/640?wx_fmt=png&tp=webp&wxfrom=5&wx_lazy=1&wx_co=1)

相对比本机网络 IO 通信过程上，它的工作过程要清爽许多。其中 127.0.0.1 工作过程如下图。

![](https://mmbiz.qpic.cn/mmbiz_png/BBjAFF4hcwqr5dxehUyQGUfk69ibqFibVMYBNJiclZk13tB111axoKtso0QfeeZMKPrLohFlPoNBAGTzhviay28ibicA/640?wx_fmt=png&tp=webp&wxfrom=5&wx_lazy=1&wx_co=1)

我们也对比了 UDP 和 TCP 两种方式下的延迟和性能指标。在包体不大于 1KB 的时候，UDS 的性能大约是 TCP 的两倍多。 **所以，在本机网络 IO 的场景下，如果对性能敏感，飞哥建议你使用 Unix Domain Socket。**
