---
title: 用户空间和内核空间通讯--netlink
---

**目录**

- [什么是 Netlink](https://e-mailky.github.io/2017-02-14-netlink-user-kernel1#%E4%BB%80%E4%B9%88%E6%98%AFnetlink)
- [Netlink 通信类型](https://e-mailky.github.io/2017-02-14-netlink-user-kernel1#netlink%E9%80%9A%E4%BF%A1%E7%B1%BB%E5%9E%8B)
- [Netlink 的消息格式](https://e-mailky.github.io/2017-02-14-netlink-user-kernel1#netlink%E7%9A%84%E6%B6%88%E6%81%AF%E6%A0%BC%E5%BC%8F)
- [Netlink 的消息头](https://e-mailky.github.io/2017-02-14-netlink-user-kernel1#netlink%E7%9A%84%E6%B6%88%E6%81%AF%E5%A4%B4)
- [Netlink 的消息体](https://e-mailky.github.io/2017-02-14-netlink-user-kernel1#netlink%E7%9A%84%E6%B6%88%E6%81%AF%E4%BD%93)
- [Netlink 提供的错误指示消息](https://e-mailky.github.io/2017-02-14-netlink-user-kernel1#netlink%E6%8F%90%E4%BE%9B%E7%9A%84%E9%94%99%E8%AF%AF%E6%8C%87%E7%A4%BA%E6%B6%88%E6%81%AF)
- [Netlink 编程需要注意的问题](https://e-mailky.github.io/2017-02-14-netlink-user-kernel1#netlink%E7%BC%96%E7%A8%8B%E9%9C%80%E8%A6%81%E6%B3%A8%E6%84%8F%E7%9A%84%E9%97%AE%E9%A2%98)
- [Netlink 的地址结构体](https://e-mailky.github.io/2017-02-14-netlink-user-kernel1#netlink%E7%9A%84%E5%9C%B0%E5%9D%80%E7%BB%93%E6%9E%84%E4%BD%93)
- [实验](https://e-mailky.github.io/2017-02-14-netlink-user-kernel1#%E5%AE%9E%E9%AA%8C)
  - [第一步](https://e-mailky.github.io/2017-02-14-netlink-user-kernel1#%E7%AC%AC%E4%B8%80%E6%AD%A5)
  - [第二步](https://e-mailky.github.io/2017-02-14-netlink-user-kernel1#%E7%AC%AC%E4%BA%8C%E6%AD%A5)
  - [第三步](https://e-mailky.github.io/2017-02-14-netlink-user-kernel1#%E7%AC%AC%E4%B8%89%E6%AD%A5)
- [Netlink 多播](https://e-mailky.github.io/2017-02-14-netlink-user-kernel1#netlink%E5%A4%9A%E6%92%AD)

[请尊重原创版权，转载注明出处。](https://e-mailky.github.io/2017-02-14-netlink-user-kernel1)

> 参考：
> - 原文链接：<https://e-mailky.github.io/2017-02-14-netlink-user-kernel1>

Alan Cox 在内核 1.3 版本的开发阶段最先引入了 Netlink，刚开始时 Netlink 是以 字符驱动接口的方式提供内核与用户空间的双向数据通信；随后，在 2.1 内核开发过程中，Alexey Kuznetsov 将 Netlink 改写成一个更加灵活、且易于扩展的基于消息通信接口，并将其应用到高级路由子系统的基础框架里。 自那时起，Netlink 就成了 Linux 内核子系统和用户态的应用程序通信的主要手段之一。
2001 年，ForCES IETF 委员会正式对 Netlink 进行了标准化的工作。Jamal Hadi Salim 提议将 Netlink 定义成一种用于网络设备的路由引擎组件和其控制管理组件之间通信的协议。不过他的建议 最终没有被采纳，取而代之的是我们今天所看到的格局：Netlink 被设计成一个新的协议域，domain。
Linux 之父托瓦斯曾说过“Linux is evolution, not intelligent design”。 什么意思？就是说，Netlink 也同样遵循了 Linux 的某些设计理念，即没有完整的规范文档，亦没有设计文档。 只有什么？你懂得—“Read the f\*\*king source code”。
当然，本文不是分析 Netlink 在 Linux 上的实现机制，而是就“什么是 Netlink”以及 “如何用好 Netlink”的话题和大家做个分享，只有在遇到问题时才需要去阅读内核源码弄清个所以然。

## 什么是 Netlink

关于 Netlink 的理解，需要把握几个关键点：

1. 面向数据报的无连接消息子系统
2. 基于通用的 BSD Socket 架构而实现

关于第一点使我们很容易联想到 UDP 协议，能想到这一点就非常棒了。按着 UDP 协议 来理解 Netlink 不是不无道理，只要你能触类旁通，做到“活学”，善于总结归纳、联想，最后实现知识迁移这就是 学习的本质。Netlink 可以实现内核->用户以及用户->内核的双向、异步的数据通信，同时它还支持两个用户进程之间、 甚至两个内核子系统之间的数据通信。本文中，对后两者我们不予考虑，焦点集中在如何实现用户<->内核之间的数据通信。
看到第二点脑海中是不是瞬间闪现了下面这张图片呢？ 如果是，则说明你确实有慧根；当然，不是也没关系，慧根可以慢慢长嘛，呵呵。
![](https://notes-learning.oss-cn-beijing.aliyuncs.com/oextqn/1621609972955-cbd73443-4468-439d-87a1-cc83f2daa893.jpeg)
在后面实战 Netlink 套接字编程时我们主要会用到 socket()，bind()，sendmsg() 和 recvmsg()等系统调用，当然还有 socket 提供的轮训(polling)机制。

## Netlink 通信类型

Netlink 支持两种类型的通信方式：单播和多播。
单播：经常用于一个用户进程和一个内核子系统之间 1:1 的数据通信。用户空间发送命令到内核，然后从内核接受命令的返回结果。
多播：经常用于一个内核进程和多个用户进程之间的 1:N 的数据通信。内核作为会话的发起者， 用户空间的应用程序是接收者。为了实现这个功能，内核空间的程序会创建一个多播组， 然后所有用户空间的对该内核进程发送的消息感兴趣的进程都加入到该组即可接收来自内核发送的消息了。如下：
![](https://notes-learning.oss-cn-beijing.aliyuncs.com/oextqn/1621609976015-6753cbdb-f771-4657-a4b6-d7c3735f8d40.jpeg)
其中进程 A 和子系统 1 之间是单播通信，进程 B、C 和子系统 2 是多播通信。 上图还向我们说明了一个信息。从用户空间传递到内核的数据是不需要排队的，即其操作是同步完成； 而从内核空间向用户空间传递数据时需要排队，是异步的。了解了这一点在开发基于 Netlink 的应用模块时 可以使我们少走很多弯路。假如，你向内核发送了一个消息需要获取内核中某些信息，比如路由表，或其他信息， 如果路由表过于庞大，那么内核在通过 Netlink 向你返回数据时，你可以好生琢磨一下如何接收这些数据的问题， 毕竟你已经看到了那个输出队列了，不能视而不见啊。

## Netlink 的消息格式

Netlink 消息由两部分组成：消息头和有效数据载荷， 且整个 Netlink 消息是 4 字节对齐，一般按主机字节序进行传递。消息头为固定的 16 字节，消息体长度可变：
![](https://notes-learning.oss-cn-beijing.aliyuncs.com/oextqn/1621609971315-01096608-1381-4ee9-9d56-d0c80b0ffe78.jpeg)

## Netlink 的消息头

消息头定义在文件里，由结构体 nlmsghdr 表示：
struct nlmsghdr { \_\_u32 nlmsg*len; /* Length of message including header _/ \_\_u16 nlmsg_type; /_ Message content _/ \_\_u16 nlmsg_flags; /_ Additional flags _/ \_\_u32 nlmsg_seq; /_ Sequence number _/ \_\_u32 nlmsg_pid; /_ Sending process PID \_/ };&#x20;
消息头中各成员属性的解释及说明：
nlmsg_len：整个消息的长度，按字节计算。包括了 Netlink 消息头本身。
nlmsg_type：消息的类型，即是数据还是控制消息。目前(内核版本 2.6.21)Netlink 仅支持四种类型的控制消息，如下：
NLMSG_NOOP- 空消息，什么也不做； NLMSG_ERROR- 指明该消息中包含一个错误； NLMSG_DONE- 如果内核通过 Netlink 队列返回了多个消息，那么队列的最后一条消息的类型为 NLMSG_DONE， 其余所有消息的 nlmsg_flags 属性都被设置 NLM_F_MULTI 位有效。 NLMSG_OVERRUN-暂时没用到。&#x20;
nlmsg_flags：附加在消息上的额外说明信息，如上面提到的 NLM_F_MULTI。摘录如下：
&#x20; 标记 | 作用及说明 --------------------|----------- NLM_F_REQUEST | 如果消息中有该标记位，说明这是一个请求消息。所有从用户空间到内核空间的消息都要设置该位，否则内核将向用户返回一个 EINVAL 无效参数的错误 NLM_F_MULTI | 消息从用户->内核是同步的立刻完成，而从内核->用户则需要排队。如果内核之前收到过来自用户的消息中有 NLM_F_DUMP 位为 1 的消息，那么内核就会向用户空间发送一个由多个 Netlink 消息组成的链表。除了最后个消息外，其余每条消息中都设置了该位有效。 NLM_F_ACK | 该消息是内核对来自用户空间的 NLM_F_REQUEST 消息的响应 NLM_F_ECHO | 如果从用户空间发给内核的消息中该标记为 1，则说明用户的应用进程要求内核将用户发给它的每条消息通过单播的形式再发送给用户进程。和我们通常说的“回显”功能类似。 ... | ...&#x20;
大家只要知道 nlmsg_flags 有多种取值就可以，至于每种值的作用和意义， 通过谷歌和源代码一定可以找到答案，这里就不展开了。上一张 2.6.21 内核中所有的取值情况：
![](https://notes-learning.oss-cn-beijing.aliyuncs.com/oextqn/1621609971427-92b211e9-5d53-41f0-9d5d-9c784df8397a.jpeg)
nlmsg_seq：消息序列号。因为 Netlink 是面向数据报的，所以存在丢失数据的风险，但是 Netlink 提供了如何确保消息不丢失的机制， 让程序开发人员根据其实际需求而实现。消息序列号一般和 NLM_F_ACK 类型的消息联合使用， 如果用户的应用程序需要保证其发送的每条消息都成功被内核收到的话，那么它发送消息时 需要用户程序自己设置序号，内核收到该消息后对提取其中的序列号，然后在发送给用户程序 回应消息里设置同样的序列号。有点类似于 TCP 的响应和确认机制。
**注意：当内核主动向用户空间发送广播消息时，消息中的该字段总是为 0。**
nlmsg_pid：当用户空间的进程和内核空间的某个子系统之间通过 Netlink 建立了数据交换的通道后， Netlink 会为每个这样的通道分配一个唯一的数字标识。其主要作用就是将来自用户空间的请求消息 和响应消息进行关联。说得直白一点，假如用户空间存在多个用户进程，内核空间同样存在多个进程， Netlink 必须提供一种机制用于确保每一对“用户-内核”空间通信的进程之间的数据交互不会发生紊乱。
![](https://notes-learning.oss-cn-beijing.aliyuncs.com/oextqn/1621609971166-de3027d7-ff6a-4459-bd99-766a18a23b6d.jpeg)
即，进程 A、B 通过 Netlink 向子系统 1 获取信息时，子系统 1 必须确保回送给进程 A 的 响应数据不会发到进程 B 那里。主要适用于用户空间的进程从内核空间获取数据的场景。通常情况下，用户空间的进程 在向内核发送消息时一般通过系统调用 getpid()将当前进程的进程号赋给该变量，即用户空间的进程希望得到 内核的响应时才会这么做。从内核主动发送到用户空间的消息该字段都被设置为 0。

## Netlink 的消息体

Netlink 的消息体采用 TLV(Type-Length-Value)格式：
![](https://notes-learning.oss-cn-beijing.aliyuncs.com/oextqn/1621609974851-b2f8574a-cbf8-4cc5-a9fe-c21246356b52.jpeg)
Netlink 每个属性都由文件里的 struct nlattr{}来表示：
struct nlattr { \_\_u16 nla_len; \_\_u16 nla_type; };

## Netlink 提供的错误指示消息

当用户空间的应用程序和内核空间的进程之间通过 Netlink 通信时发生了错误， Netlink 必须向用户空间通报这种错误。Netlink 对错误消息进行了单独封装，：
struct nlmsgerr { int error; //标准的错误码，定义在 errno.h 头文件中。可以用 perror()来解释 struct nlmsghdr msg; //指明了哪条消息触发了结构体中 error 这个错误值 };

## Netlink 编程需要注意的问题

基于 Netlink 的用户-内核通信，有两种情况可能会导致丢包：

1. 内存耗尽；
2. 用户空间接收进程的缓冲区溢出。导致缓冲区溢出的主要原因有可能是：用户空间的进程运行太慢；或者接收队列太短。

如果 Netlink 不能将消息正确传递到用户空间的接收进程，那么用户空间的接收进程在 调用 recvmsg()系统调用时就会返回一个内存不足(ENOBUFS)的错误，这一点需要注意。换句话说， 缓冲区溢出的情况是不会发送在从用户->内核的 sendmsg()系统调用里，原因前面我们也说过了，请大家自己思考一下。
当然，如果使用的是阻塞型 socket 通信，也就不存在内存耗尽的隐患了，这又是为什么呢？ 赶紧去谷歌一下，查查什么是阻塞型 socket 吧。学而不思则罔，思而不学则殆嘛。

## Netlink 的地址结构体

在 TCP 博文中我们提到过在 Internet 编程过程中所用到的地址结构体和标准地址结构体， 它们和 Netlink 地址结构体的关系如下：
![](https://notes-learning.oss-cn-beijing.aliyuncs.com/oextqn/1621609973262-e91727eb-0e30-4423-92e3-47ff1b08ded4.jpeg)
struct sockaddr*nl{}的详细定义和描述如下：
struct sockaddr_nl { sa_family_t nl_family; /*该字段总是为 AF*NETLINK */ unsigned short nl*pad; /* 目前未用到，填充为 0*/ \_\_u32 nl_pid; /* process pid _/ \_\_u32 nl_groups; /_ multicast groups mask \_/ };&#x20;
nl_pid：该属性为发送或接收消息的进程 ID，前面我们也说过，Netlink 不仅可以实现用户-内核空间的通信还可 使现实用户空间两个进程之间，或内核空间两个进程之间的通信。该属性为 0 时一般适用于如下两种情况：

1. 我们要发送的目的地是内核，即从用户空间发往内核空间时，我们构造的 Netlink 地址结构体中 nl_pid 通常 情况下都置为 0。这里有一点需要跟大家交代一下，在 Netlink 规范里，PID 全称是 Port-ID(32bits)， 其主要作用是用于唯一的标识一个基于 netlink 的 socket 通道。通常情况下 nl_pid 都设置为当前进程的进程号。 然而，对于一个进程的多个线程同时使用 netlink socket 的情况，nl_pid 的设置一般采用如下这个样子来实现：
   | pthread_self() « 16 | getpid(); |
   | --- | --- |

2. 从内核发出的多播报文到用户空间时，如果用户空间的进程处在该多播组中，那么其地址结构体中 nl_pid 也设置为 0，同时还要结合下面介绍到的另一个属性。

nl_groups：如果用户空间的进程希望加入某个多播组，则必须执行 bind()系统调用。 该字段指明了调用者希望加入的多播组号的掩码(注意不是组号，后面我们会详细讲解这个字段)。 如果该字段为 0 则表示调用者不希望加入任何多播组。对于每个隶属于 Netlink 协议域的协议， 最多可支持 32 个多播组(因为 nl_groups 的长度为 32 比特)，每个多播组用一个比特来表示。

## 实验

在文件里包含了 Netlink 协议簇已经定义好的一些预定义协议：
\#define NETLINK*ROUTE 0 /* Routing/device hook _/ #define NETLINK_UNUSED 1 /_ Unused number _/ #define NETLINK_USERSOCK 2 /_ Reserved for user mode socket protocols _/ #define NETLINK_FIREWALL 3 /_ Firewalling hook _/ #define NETLINK_INET_DIAG 4 /_ INET socket monitoring _/ #define NETLINK_NFLOG 5 /_ netfilter/iptables ULOG _/ #define NETLINK_XFRM 6 /_ ipsec _/ #define NETLINK_SELINUX 7 /_ SELinux event notifications _/ #define NETLINK_ISCSI 8 /_ Open-iSCSI _/ #define NETLINK_AUDIT 9 /_ auditing _/ #define NETLINK_FIB_LOOKUP 10 #define NETLINK_CONNECTOR 11 #define NETLINK_NETFILTER 12 /_ netfilter subsystem _/ #define NETLINK_IP6_FW 13 #define NETLINK_DNRTMSG 14 /_ DECnet routing messages _/ #define NETLINK_KOBJECT_UEVENT 15 /_ Kernel messages to userspace _/ #define NETLINK_GENERIC 16 /_ leave room for NETLINK*DM (DM Events) */ #define NETLINK*SCSITRANSPORT 18 /* SCSI Transports _/ #define NETLINK_ECRYPTFS 19 #define NETLINK_TEST 20 /_ 用户添加的自定义协议 \_/&#x20;
如果我们在 Netlink 协议簇里开发一个新的协议，只要在该文件中定义协议号即可， 例如我们定义一种基于 Netlink 协议簇的、协议号是 20 的自定义协议，如上所示。同时记得，将内核头文件目录中 的 netlink.h 也做对应的修改，在我的系统中它的路径是：/usr/src/linux-2.6.21/include/linux/netlink.h
接下来我们在用户空间以及内核空间模块的开发过程中就可以使用这种协议了，一共分为三个阶段。

### 第一步

我们首先实现的功能是用户->内核的单向数据通信，即用户空间发送一个消息给内核，然后内核将其打印输出， 就这么简单。用户空间的示例代码如下【mynlusr.c】
\#include \<sys/stat.h> #include \<unistd.h> #include \<stdio.h> #include \<stdlib.h> #include \<sys/socket.h> #include \<sys/types.h> #include \<string.h> #include \<asm/types.h> #include \<linux/netlink.h> #include \<linux/socket.h> #define MAX*PAYLOAD 1024 /*消息最大负载为 1024 字节*/ int main(int argc, char* argv\[]) { struct sockaddr*nl dest_addr; struct nlmsghdr *nlh = NULL; struct iovec iov; int sock_fd=-1; struct msghdr msg; if(-1 == (sock_fd=socket(PF_NETLINK, SOCK_RAW,NETLINK_TEST))){ //创建套接字 perror("can't create netlink socket!"); return 1; } memset(\&dest_addr, 0, sizeof(dest_addr)); dest_addr.nl_family = AF_NETLINK; dest_addr.nl_pid = 0; /*我们的消息是发给内核的*/ dest_addr.nl_groups = 0; /*在本示例中不存在使用该值的情况*/ //将套接字和 Netlink 地址结构体进行绑定 if(-1 == bind(sock_fd, (struct sockaddr*)\&dest_addr, sizeof(dest_addr))){ perror("can't bind sockfd with sockaddr_nl!"); return 1; } if(NULL == (nlh=(struct nlmsghdr *)malloc(NLMSG*SPACE(MAX_PAYLOAD)))){ perror("alloc mem failed!"); return 1; } memset(nlh,0,MAX_PAYLOAD); /* 填充 Netlink 消息头部 */ nlh->nlmsg_len = NLMSG_SPACE(MAX_PAYLOAD); nlh->nlmsg_pid = 0; nlh->nlmsg_type = NLMSG_NOOP; //指明我们的 Netlink 是消息负载是一条空消息 nlh->nlmsg_flags = 0; /*设置 Netlink 的消息内容，来自我们命令行输入的第一个参数*/ strcpy(NLMSG_DATA(nlh), argv\[1]); /*这个是模板，暂时不用纠结为什么要这样用。有时间详细讲解 socket 时再说*/ memset(\&iov, 0, sizeof(iov)); iov.iov_base = (void *)nlh; iov.iov*len = nlh->nlmsg_len; memset(\&msg, 0, sizeof(msg)); msg.msg_iov = \&iov; msg.msg_iovlen = 1; sendmsg(sock_fd, \&msg, 0); //通过 Netlink socket 向内核发送消息 /* 关闭 netlink 套接字 */ close(sock_fd); free(nlh); return 0; }&#x20;
上面的代码逻辑已经非常清晰了，都是 socket 编程的 API，唯一不同的是我们这次编程 是针对 Netlink 协议簇的。这里我们提前引入了 BSD 层的消息结构体 struct msghdr{}，定义在文件里， 以及其数据块 struct iovec{}定义在头文件里。这里就不展开了，大家先记住这个用法就行。以后有时间再深入 到 socket 的骨子里去转悠一番。
另外，需要格外注意的就是 Netlink 的地址结构体和其消息头结构中 pid 字段为 0 的情况，很容易让人产生混淆，再总结一下：
nl_pid | 0 ————————-|—– netlink 地址结构体.nl_pid | 1、内核发出的多播报文 2、消息的接收方是内核，即从用户空间发往内核的消息 netlink 消息头体 nlmsg_pid | 来自内核主动发出的消息
这个例子仅是从用户空间到内核空间的单向数据通信，所以 Netlink 地址结构体中 我们设置了 dest_addr.nl_pid = 0，说明我们的报文的目的地是内核空间；在填充 Netlink 消息头部时， 我们做了 nlh->nlmsg_pid = 0 这样的设置。
需要注意几个宏的使用：
NLMSG_SPACE(MAX_PAYLOAD): 该宏用于返回不小于 MAX_PAYLOAD 且 4 字节对齐的最小长度值，一般用于向内存系统申请空间是指定所申请的内存字节数
NLMSG_LENGTH(len): 所不同的是，前者所申请的空间里不包含 Netlink 消息头部所占的字节数，后者是消息负载和消息头加起来的总长度
NLMSG_DATA(nlh): 该宏用于返回 Netlink 消息中数据部分的首地址，在写入和读取消息数据部分时会用到它。
它们之间的关系如下：
![](https://notes-learning.oss-cn-beijing.aliyuncs.com/oextqn/1621609977751-61452e76-5d04-4003-8e2e-80e3e173fd59.jpeg)
内核空间的示例代码如下【mynlkern.c】
\#include \<linux/kernel.h> #include \<linux/module.h> #include \<linux/skbuff.h> #include \<linux/init.h> #include \<linux/ip.h> #include \<linux/types.h> #include \<linux/sched.h> #include \<net/sock.h> #include \<linux/netlink.h> MODULE_LICENSE("GPL"); MODULE_AUTHOR("Koorey King"); struct sock *nl*sk = NULL; static void nl_data_ready (struct sock *sk, int len) { struct sk*buff \_skb; struct nlmsghdr \_nlh = NULL; while((skb = skb_dequeue(\&sk->sk_receive_queue)) != NULL) { nlh = (struct nlmsghdr *)skb->data; printk("%s: received netlink message payload: %s \n", **FUNCTION**, (char*)NLMSG_DATA(nlh)); kfree_skb(skb); } printk("recvied finished!\n"); } static int \_\_init myinit_module() { printk("my netlink in\n"); nl_sk = netlink_kernel_create(NETLINK_TEST,0,nl_data_ready,THIS_MODULE); return 0; } static void \_\_exit mycleanup_module() { printk("my netlink out!\n"); sock_release(nl_sk->sk_socket); } module_init(myinit_module); module_exit(mycleanup_module);&#x20;
在内核模块的初始化函数里我们用
nl_sk = netlink_kernel_create(NETLINK_TEST,0,nl_data_ready,THIS_MODULE);&#x20;
创建了一个内核态的 socket，第一个参数我们扩展的协议号；第二个参数为多播组号，目前我们用不上，将其置为 0； 第三个参数是个回调函数，即当内核的 Netlink socket 套接字收到数据时的处理函数；第四个参数就不多说了。
在回调函数 nl_data_ready()中，我们不断的从 socket 的接收队列去取数据， 一旦拿到数据就将其打印输出。在协议栈的 INET 层，用于存储数据的是大名鼎鼎的 sk_buff 结构，所以我们通过 nlh = (struct nlmsghdr *)skb->data；可以拿到 netlink 的消息体，然后通过 NLMSG_DATA(nlh)定位到 netlink 的消息负载。
将上述代码编译后测试结果如下：
![](https://notes-learning.oss-cn-beijing.aliyuncs.com/oextqn/1621609975598-daa41911-732a-416a-993e-1b07150c3775.jpeg)

### 第二步

我们将上面的代码稍加改造就可以实现用户<->内核的双向数据通信。
首先是改造用户空间的代码：
\#include \<sys/stat.h> #include \<unistd.h> #include \<stdio.h> #include \<stdlib.h> #include \<sys/socket.h> #include \<sys/types.h> #include \<string.h> #include \<asm/types.h> #include \<linux/netlink.h> #include \<linux/socket.h> #define MAX*PAYLOAD 1024 /*消息最大负载为 1024 字节*/ int main(int argc, char* argv\[]) { struct sockaddr*nl dest_addr; struct nlmsghdr *nlh = NULL; struct iovec iov; int sock*fd=-1; struct msghdr msg; if(-1 == (sock_fd=socket(PF_NETLINK, SOCK_RAW,NETLINK_TEST))){ perror("can't create netlink socket!"); return 1; } memset(\&dest_addr, 0, sizeof(dest_addr)); dest_addr.nl_family = AF_NETLINK; dest_addr.nl_pid = 0; /*我们的消息是发给内核的*/ dest_addr.nl_groups = 0; /*在本示例中不存在使用该值的情况*/ if(-1 == bind(sock_fd, (struct sockaddr*)\&dest_addr, sizeof(dest_addr))){ perror("can't bind sockfd with sockaddr_nl!"); return 1; } if(NULL == (nlh=(struct nlmsghdr *)malloc(NLMSG_SPACE(MAX_PAYLOAD)))){ perror("alloc mem failed!"); return 1; } memset(nlh,0,MAX_PAYLOAD); /* 填充 Netlink 消息头部 */ nlh->nlmsg_len = NLMSG_SPACE(MAX_PAYLOAD); `nlh->nlmsg_pid = getpid();//我们希望得到内核回应，所以得告诉内核我们ID号` nlh->nlmsg_type = NLMSG_NOOP; //指明我们的 Netlink 是消息负载是一条空消息 nlh->nlmsg_flags = 0; /*设置 Netlink 的消息内容，来自我们命令行输入的第一个参数*/ strcpy(NLMSG_DATA(nlh), argv\[1]); /*这个是模板，暂时不用纠结为什么要这样用。*/ memset(\&iov, 0, sizeof(iov)); iov.iov*base = (void *)nlh; iov.iov*len = nlh->nlmsg_len; memset(\&msg, 0, sizeof(msg)); msg.msg_iov = \&iov; msg.msg_iovlen = 1; sendmsg(sock_fd, \&msg, 0); //通过 Netlink socket 向内核发送消息 //接收内核消息的消息 printf("waiting message from kernel!\n"); memset((char*)NLMSG*DATA(nlh),0,1024); recvmsg(sock_fd,\&msg,0); printf("Got response: %s\n",NLMSG_DATA(nlh)); /* 关闭 netlink 套接字 */ close(sock_fd); free(nlh); return 0; }&#x20;
内核空间的修改如下：
\#include \<linux/kernel.h> #include \<linux/module.h> #include \<linux/skbuff.h> #include \<linux/init.h> #include \<linux/ip.h> #include \<linux/types.h> #include \<linux/sched.h> #include \<net/sock.h> #include \<net/netlink.h> /*该文头文件里包含了 linux/netlink.h，因为我们要用到 net/netlink.h 中的某些 API 函数，nlmsg*put()*/ MODULE*LICENSE("GPL"); MODULE_AUTHOR("Koorey King"); struct sock *nl_sk = NULL; //向用户空间发送消息的接口 void sendnlmsg(char *message,int dstPID) { struct sk_buff *skb; struct nlmsghdr *nlh; int len = NLMSG_SPACE(MAX_MSGSIZE); int slen = 0; if(!message || !nl_sk){ return; } // 为新的 sk_buffer 申请空间 skb = alloc_skb(len, GFP_KERNEL); if(!skb){ printk(KERN_ERR "my_net_link: alloc_skb Error./n"); return; } slen = strlen(message)+1; //用 nlmsg_put()来设置 netlink 消息头部 nlh = nlmsg_put(skb, 0, 0, 0, MAX_MSGSIZE, 0); // 设置 Netlink 的控制块 NETLINK_CB(skb).pid = 0; // 消息发送者的 id 标识，如果是内核发的则置 0 NETLINK_CB(skb).dst_group = 0; //如果目的组为内核或某一进程，该字段也置 0 message\[slen] = '\0'; memcpy(NLMSG_DATA(nlh), message, slen+1); //通过 netlink_unicast()将消息发送用户空间由 dstPID 所指定了进程号的进程 netlink_unicast(nl_sk,skb,dstPID,0); printk("send OK!\n"); return; } static void nl_data_ready (struct sock *sk, int len) { struct sk_buff \_skb; struct nlmsghdr \_nlh = NULL; while((skb = skb_dequeue(\&sk->sk_receive_queue)) != NULL) { nlh = (struct nlmsghdr *)skb->data; printk("%s: received netlink message payload: %s \n", **FUNCTION**, (char\*)NLMSG_DATA(nlh)); sendnlmsg("I see you",nlh->nlmsg_pid); //发送者的进程 ID 我们已经将其存储在了 netlink 消息头部里的 nlmsg_pid 字段里，所以这里可以拿来用。 kfree_skb(skb); } printk("recvied finished!\n"); } static int \_\_init myinit_module() { printk("my netlink in\n"); nl_sk = netlink_kernel_create(NETLINK_TEST,0,nl_data_ready,THIS_MODULE); return 0; } static void \_\_exit mycleanup_module() { printk("my netlink out!\n"); sock_release(nl_sk->sk_socket); } module_init(myinit_module); module_exit(mycleanup_module);&#x20;
重新编译后，测试结果如下：
![](https://notes-learning.oss-cn-beijing.aliyuncs.com/oextqn/1621609976403-7d22f39a-4e7c-4c7c-aec3-376db6379b04.jpeg)

### 第三步

前面我们提到过，如果用户进程希望加入某个多播组时才需要调用 bind()函数。 前面的示例中我们没有这个需求，可还是调了 bind()，心头有些不爽。在前几篇博文里有关于 socket 编程时 几个常见 API 的详细解释和说明，不明白的童鞋可以回头去复习一下。
因为 Netlink 是面向无连接的数据报的套接字，所以我们还可以用 sendto()和 recvfrom() 来实现数据的收发，这次我们不再调用 bind()。将 Stage 2 的例子稍加改造一下，用户空间的修改如下：
\#include \<sys/stat.h> #include \<unistd.h> #include \<stdio.h> #include \<stdlib.h> #include \<sys/socket.h> #include \<sys/types.h> #include \<string.h> #include \<asm/types.h> #include \<linux/netlink.h> #include \<linux/socket.h> #define MAX*PAYLOAD 1024 /*消息最大负载为 1024 字节*/ int main(int argc, char* argv\[]) { struct sockaddr*nl dest_addr; struct nlmsghdr *nlh = NULL; //struct iovec iov; int sock*fd=-1; //struct msghdr msg; if(-1 == (sock_fd=socket(PF_NETLINK, SOCK_RAW,NETLINK_TEST))){ perror("can't create netlink socket!"); return 1; } memset(\&dest_addr, 0, sizeof(dest_addr)); dest_addr.nl_family = AF_NETLINK; dest_addr.nl_pid = 0; /*我们的消息是发给内核的*/ dest_addr.nl_groups = 0; /*在本示例中不存在使用该值的情况*/ /*不再调用 bind()函数了 if(-1 == bind(sock*fd, (struct sockaddr*)\&dest*addr, sizeof(dest_addr))){ perror("can't bind sockfd with sockaddr_nl!"); return 1; }*/ if(NULL == (nlh=(struct nlmsghdr *)malloc(NLMSG*SPACE(MAX_PAYLOAD)))){ perror("alloc mem failed!"); return 1; } memset(nlh,0,MAX_PAYLOAD); /* 填充 Netlink 消息头部 */ nlh->nlmsg_len = NLMSG_SPACE(MAX_PAYLOAD); nlh->nlmsg_pid = getpid();//我们希望得到内核回应，所以得告诉内核我们 ID 号 nlh->nlmsg_type = NLMSG_NOOP; //指明我们的 Netlink 是消息负载是一条空消息 nlh->nlmsg_flags = 0; /*设置 Netlink 的消息内容，来自我们命令行输入的第一个参数*/ strcpy(NLMSG_DATA(nlh), argv\[1]); /*这个模板就用不上了。*/ /\_ memset(\&iov, 0, sizeof(iov)); iov.iov*base = (void *)nlh; iov.iov*len = nlh->nlmsg_len; memset(\&msg, 0, sizeof(msg)); msg.msg_iov = \&iov; msg.msg_iovlen = 1; */ //sendmsg(sock*fd, \&msg, 0); //不再用这种方式发消息到内核 sendto(sock_fd,nlh,NLMSG_LENGTH(MAX_PAYLOAD),0,(struct sockaddr*)(\&dest*addr),sizeof(dest_addr)); //接收内核消息的消息 printf("waiting message from kernel!\n"); //memset((char*)NLMSG_DATA(nlh),0,1024); memset(nlh,0,MAX_PAYLOAD); //清空整个 Netlink 消息头包括消息头和负载 //recvmsg(sock_fd,\&msg,0); recvfrom(sock_fd,nlh,NLMSG_LENGTH(MAX_PAYLOAD),0,(struct sockaddr*)(\&dest_addr),NULL); printf("Got response: %s\n",NLMSG_DATA(nlh)); /* 关闭 netlink 套接字 \_/ close(sock_fd); free(nlh); return 0; }&#x20;
内核空间的代码完全不用修改，我们仍然用 netlink_unicast()从内核空间发送消息到用户空间。
重新编译后，测试结果如下：
![](https://notes-learning.oss-cn-beijing.aliyuncs.com/oextqn/1621609976660-e321adca-1224-41cf-8fda-74048e4f6523.jpeg)
和 Stage 2 中代码运行效果完全一样。也就是说，在开发 Netlink 程序过程中， 如果没牵扯到多播机制，那么用户空间的 socket 代码其实是不用执行 bind()系统调用的， 但此时就需要用 sendto()和 recvfrom()完成数据的发送和接收的任务；如果执行了 bind()系统调用， 当然也可以继续用 sendto()和 recvfrom()，但给它们传递的参数就有所区别。这时候一般使用 sendmsg()和 recvmsg() 来完成数据的发送和接收。大家根据自己的实际情况灵活选择。

## Netlink 多播

在上面我们所遇到的情况都是用户空间作为消息进程的发起者， Netlink 还支持内核作为消息的发送方的情况。这一般用于内核主动向用户空间报告一些内核状态， 例如我们在用户空间看到的 USB 的热插拔事件的通告就是这样的应用。
先说一下我们的目标，内核线程每个一秒钟往一个多播组里发送一条消息， 然后用户空间所以加入了该组的进程都会收到这样的消息，并将消息内容打印出来。
Netlink 地址结构体中的 nl*groups 是 32 位，也就是说每种 Netlink 协议最多支持 32 个多播组。 如何理解这里所说的每种 Netlink 协议？在里预定义的如下协议都是 Netlink 协议簇的具体协议， 还有我们添加的 NETLINK_TEST 也是一种 Netlink 协议。
\#define NETLINK_ROUTE 0 /* Routing/device hook _/ #define NETLINK_UNUSED 1 /_ Unused number _/ #define NETLINK_USERSOCK 2 /_ Reserved for user mode socket protocols _/ #define NETLINK_FIREWALL 3 /_ Firewalling hook _/ #define NETLINK_INET_DIAG 4 /_ INET socket monitoring _/ #define NETLINK_NFLOG 5 /_ netfilter/iptables ULOG _/ #define NETLINK_XFRM 6 /_ ipsec _/ #define NETLINK_SELINUX 7 /_ SELinux event notifications _/ #define NETLINK_ISCSI 8 /_ Open-iSCSI _/ #define NETLINK_AUDIT 9 /_ auditing _/ #define NETLINK_FIB_LOOKUP 10 #define NETLINK_CONNECTOR 11 #define NETLINK_NETFILTER 12 /_ netfilter subsystem _/ #define NETLINK_IP6_FW 13 #define NETLINK_DNRTMSG 14 /_ DECnet routing messages _/ #define NETLINK_KOBJECT_UEVENT 15 /_ Kernel messages to userspace _/ #define NETLINK_GENERIC 16 /_ leave room for NETLINK*DM (DM Events) */ #define NETLINK*SCSITRANSPORT 18 /* SCSI Transports _/ #define NETLINK_ECRYPTFS 19 #define NETLINK_TEST 20 /_ 用户添加的自定义协议 _/&#x20;
在我们自己添加的 NETLINK_TEST 协议里，同样地，最多允许我们设置 32 个多播组， 每个多播组用 1 个比特表示，所以不同的多播组不可能出现重复。你可以根据自己的实际需求，决定哪个多播组是用来做什么的。 用户空间的进程如果对某个多播组感兴趣，那么它就加入到该组中，当内核空间的进程往该组发送多播消息时， 所有已经加入到该多播组的用户进程都会收到该消息。
再回到我们 Netlink 地址结构体里的 nl_groups 成员，它是多播组的地址掩码， 注意是掩码不是多播组的组号。如何根据多播组号取得多播组号的掩码呢？在 af_netlink.c 中有个函数：
static u32 netlink_group_mask(u32 group) { return group ? 1 << (group - 1) : 0; }&#x20;
也就是说，在用户空间的代码里，如果我们要加入到多播组 1，需要设置 nl_groups 设置为 1； 多播组 2 的掩码为 2；多播组 3 的掩码为 4，依次类推。为 0 表示我们不希望加入任何多播组。理解这一点很重要。 所以我们可以在用户空间也定义一个类似于 netlink_group_mask()的功能函数，完成从多播组号到多播组掩码的转换。 最终用户空间的代码如下：
\#include \<sys/stat.h> #include \<unistd.h> #include \<stdio.h> #include \<stdlib.h> #include \<sys/socket.h> #include \<sys/types.h> #include \<string.h> #include \<asm/types.h> #include \<linux/netlink.h> #include \<linux/socket.h> #include \<errno.h> #define MAX_PAYLOAD 1024 // Netlink 消息的最大载荷的长度 unsigned int netlink_group_mask(unsigned int group) { return group ? 1 << (group - 1) : 0; } int main(int argc, char* argv\[]) { struct sockaddr_nl src_addr; struct nlmsghdr *nlh = NULL; struct iovec iov; struct msghdr msg; int sock_fd, retval; // 创建 Socket sock_fd = socket(PF_NETLINK, SOCK_RAW, NETLINK_TEST); if(sock_fd == -1){ printf("error getting socket: %s", strerror(errno)); return -1; } memset(\&src_addr, 0, sizeof(src_addr)); src_addr.nl_family = PF_NETLINK; src_addr.nl_pid = 0; // 表示我们要从内核接收多播消息。注意：该字段为 0 有双重意义，另一个意义是表示我们发送的数据的目的地址是内核。 src_addr.nl_groups = netlink_group_mask(atoi(argv\[1])); // 多播组的掩码，组号来自我们执行程序时输入的第一个参数 // 因为我们要加入到一个多播组，所以必须调用 bind()。 retval = bind(sock_fd, (struct sockaddr*)\&src_addr, sizeof(src_addr)); if(retval < 0){ printf("bind failed: %s", strerror(errno)); close(sock_fd); return -1; } // 为接收 Netlink 消息申请存储空间 nlh = (struct nlmsghdr *)malloc(NLMSG_SPACE(MAX_PAYLOAD)); if(!nlh){ printf("malloc nlmsghdr error!\n"); close(sock_fd); return -1; } memset(nlh, 0, NLMSG_SPACE(MAX_PAYLOAD)); iov.iov_base = (void *)nlh; iov.iov_len = NLMSG_SPACE(MAX_PAYLOAD); memset(\&msg, 0, sizeof(msg)); msg.msg_iov = \&iov; msg.msg_iovlen = 1; // 从内核接收消息 printf("waitinf for...\n"); recvmsg(sock_fd, \&msg, 0); printf("Received message: %s \n", NLMSG_DATA(nlh)); close(sock_fd); return 0; }&#x20;
可以看到，用户空间的程序基本没什么变化，唯一需要格外注意的就是 Netlink 地址结构体 中的 nl_groups 的设置。由于对它的解释很少，加之没有有效的文档，所以我也是一边看源码，一边在网上搜集资料。 有分析不当之处，还请大家帮我指出。
内核空间我们添加了内核线程和内核线程同步方法 completion 的使用。内核空间修改后的代码如下：
\#include \<linux/kernel.h> #include \<linux/module.h> #include \<linux/skbuff.h> #include \<linux/init.h> #include \<linux/ip.h> #include \<linux/types.h> #include \<linux/sched.h> #include \<net/sock.h> #include \<net/netlink.h> MODULE_LICENSE("GPL"); MODULE_AUTHOR("Koorey King"); struct sock *nl_sk = NULL; static struct task_struct *mythread = NULL; //内核线程对象 //向用户空间发送消息的接口 void sendnlmsg(char *message/_,int dstPID*/) { struct sk_buff *skb; struct nlmsghdr _nlh; int len = NLMSG_SPACE(MAX_MSGSIZE); int slen = 0; if(!message || !nl_sk){ return; } // 为新的 sk_buffer 申请空间 skb = alloc_skb(len, GFP_KERNEL); if(!skb){ printk(KERN_ERR "my_net_link: alloc_skb Error./n"); return; } slen = strlen(message)+1; //用 nlmsg_put()来设置 netlink 消息头部 nlh = nlmsg_put(skb, 0, 0, 0, MAX_MSGSIZE, 0); // 设置 Netlink 的控制块里的相关信息 NETLINK_CB(skb).pid = 0; // 消息发送者的 id 标识，如果是内核发的则置 0 NETLINK_CB(skb).dst_group = 5; //多播组号为 5，但置成 0 好像也可以。 message\[slen] = '\0'; memcpy(NLMSG_DATA(nlh), message, slen+1); //通过 netlink_unicast()将消息发送用户空间由 dstPID 所指定了进程号的进程 //netlink_unicast(nl_sk,skb,dstPID,0); netlink_broadcast(nl_sk, skb, 0,5, GFP_KERNEL); //发送多播消息到多播组 5，这里我故意没有用 1 之类的“常见”值，目的就是为了证明我们上面提到的多播组号和多播组号掩码之间的对应关系 printk("send OK!\n"); return; } //每隔 1 秒钟发送一条“I am from kernel!”消息，共发 10 个报文 static int sending_thread(void \_data) { int i = 10; struct completion cmpl; while(i--){ init_completion(\&cmpl); wait_for_completion_timeout(\&cmpl, 1 _ HZ); sendnlmsg("I am from kernel!"); } printk("sending thread exited!"); return 0; } static int \_\_init myinit*module() { printk("my netlink in\n"); nl_sk = netlink_kernel_create(NETLINK_TEST,0,NULL,THIS_MODULE); if(!nl_sk){ printk(KERN_ERR "my_net_link: create netlink socket error.\n"); return 1; } printk("my netlink: create netlink socket ok.\n"); mythread = kthread_run(sending_thread,NULL,"thread_sender"); return 0; } static void \_\_exit mycleanup_module() { if(nl_sk != NULL){ sock_release(nl_sk->sk_socket); } printk("my netlink out!\n"); } module_init(myinit_module); module_exit(mycleanup_module);&#x20;
关于内核中 netlink_kernel_create(int unit, unsigned int groups,…)函数里的 第二个参数指的是我们内核进程最多能处理的多播组的个数，如果该值小于 32，则默认按 32 处理， 所以在调用 netlink_kernel_create()函数时可以不用纠结第二个参数，一般将其置为 0 就可以了。
在 skbuff{}结构体中，有个成员叫做”控制块”，源码对它的解释如下：
struct sk_buff { /* These two members must be first. _/ struct sk_buff *next; struct sk_buff *prev; … … /\* _ This is the control buffer. It is free to use for every _ layer. Please put your private variables there. If you _ want to keep them across layers you have to do a skb*clone() * first. This is owned by whoever has the skb queued ATM. _/ char cb\[48]; … … }&#x20;
当内核态的 Netlink 发送数据到用户空间时一般需要填充 skbuff 的控制块，填充的方式是通过强制类型转换， 将其转换成 struct netlink_skb_parms{}之后进行填充赋值的：
struct netlink_skb_parms { struct ucred creds; /_ Skb credentials _/ \_\_u32 pid; \_\_u32 dst_group; kernel_cap_t eff_cap; \_\_u32 loginuid; /_ Login (audit) uid _/ \_\_u32 sid; /_ SELinux security id \*/ };&#x20;
填充时的模板代码如下：
NETLINK_CB(skb).pid=xx; NETLINK_CB(skb).dst_group=xx;&#x20;
这里要注意的是在 Netlink 协议簇里提到的 skbuff 的 cb 控制块里保存的是属于 Netlink 的私有信息。 怎么讲，就是 Netlink 会用该控制块里的信息来完成它所提供的一些功能，只是完成 Netlink 功能所必需的一些私有数据。 打个比方，以开车为例，开车的时候我们要做的就是打火、控制方向盘、适当地控制油门和刹车，车就开动了， 这就是汽车提供给我们的“功能”。汽车的发动机，轮胎，传动轴，以及所用到的螺丝螺栓等都属于它的“私有”数据 cb。 汽车要运行起来这些东西是不可或缺的，但它们之间的协作和交互对用户来说又是透明的。就好比我们 Netlink 的 私有控制结构 struct netlink_skb_parms{}一样。
目前我们的例子中，将 NETLINK_CB(skb).dst_group 设置为相应的多播组号和 0 效果都是一样，用户空间都可以收到该多播消息， 原因还不是很清楚，还请 Netlink 的大虾们帮我点拨点拨。
编译后重新运行，最后的测试结果如下：
![](https://notes-learning.oss-cn-beijing.aliyuncs.com/oextqn/1621609979369-87b7a52d-2caa-4bd9-be17-1069e4aee0ad.jpeg)
注意，这里一定要先执行 insmod 加载内核模块，然后再运行用户空间的程序。如果没有加载 mynlkern.ko 而直接执行./test 5 在 bind()系统调用时会报如下的错误：
bind failed: No such file or directory
因为网上有写文章在讲老版本 Netlink 的多播时用法时先执行了用户空间的程序， 然后才加载内核模块，现在(2.6.21)已经行不通了，这一点请大家注意。
小结：通过这三篇博文我们对 Netlink 有了初步的认识，并且也可以开发基于 Netlink 的基本应用程序。 但这只是冰山一角，要想写出高质量、高效率的软件模块还有些差距，特别是对 Netlink 本质的理解还需要提高一个层次， 当然这其中牵扯到内核编程的很多基本功，如临界资源的互斥、线程安全性保护、用 Netlink 传递大数据时的处理等等 都是开发人员需要考虑的问题。
