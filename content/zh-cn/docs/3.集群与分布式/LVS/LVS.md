---
title: LVS
linkTitle: LVS
weight: 1
---

# 概述

> 参考：
>
> - [Wiki, LVS](https://en.wikipedia.org/wiki/Linux_Virtual_Server)
> - [Wiki, IPVS](https://en.wikipedia.org/wiki/IP_Virtual_Server)
> - [官网](http://www.linuxvirtualserver.org/)
> - [官网,中文](http://www.linuxvirtualserver.org/zh/)
> - [官方文档，HOWTO](http://www.austintek.com/LVS/LVS-HOWTO/HOWTO/)

**Linux Virtual Server(Linux 虚拟服务器，简称 LVC)** 是一个可以实现虚拟的服务器集群功能的项目，用于实现负载均衡的软件技术。一般情况下，LVS 代之一组服务器，对于外部客户端来说，这似乎是一台服务器，所以，也称为 **。**

目前，LVS 项目已经被集成到 Linux 内核中，并通过 [IPVS](/docs/3.集群与分布式/LVS/IPVS.md)模块实现。LVS 具有良好的可靠性、可扩展性和可操作性，加上其实现最优的集群服务性能所需的低廉成本， LVS 的负载均衡功能经常被用于高性能、高可用的服务器群集中。

LVS 项目在 1998 年 5 月由[章文嵩](https://baike.baidu.com/item/%E7%AB%A0%E6%96%87%E5%B5%A9/6689425)博士成立，是中国国内最早出现的自由软件项目之一。在 linux2.2 内核时，IPVS 就已经以内核补丁的形式出现。从 2.4 版本以后 IPVS 已经成为 Linux 内核官方标准内核的一部分

## 名词解释

调度器的称呼：scheduler，director，dispatcher，balancer

- **Director(指挥器)** # 运行 IPVS 的节点。

  - **IPVS(IP 虚拟服务)** # 实现调度功能的程序。是一个 Linux 内核模块。实际上，IPVS 就是一个 **Schedulers(调度器)**。
    - **Forwarding Method(转发方法)** # Forwarding Method 用来确定 Director 如何将数据包从客户端转发到 Real Servers。如果把 Director 比做路由器，其转发数据包的规则与普通路由器有所不同。
      - Forwarding Method 其实就是指 LVS 的工作模式，当前有 LVS-NAT、LVS-DR、LVS-TUN 这几种。
  - **ipvsadm** # 为 IPVS 程序配置调度规则的用户端应用程序。

- **Real Server(真实服务器，简称 RS)** # 处理来自客户端请求的节点
- **Linux Virtual Server(简称 LVS)** # Director 与 Real Server 共同组成 LVS 集群。这些机器一起构成虚拟服务器，对于客户端来说，它表现为一台机器。
- **Client IP** # CIP,客户端 IP，用户发送请求报文的 IP

- **Director IP** # DIP,调度器 IP

- **Virual IP** # VIP,虚拟 IP，用于提供提供虚拟服务的 IP,该 IP 存在于 Director 和 RS 上

  - 为什么叫虚拟的 IP，因为这个 IP 可以代表 Director，也可以代表很多 RS，把 Director 和 RS 的很多 IP 合成 一个 IP，就称为虚拟的 IP。
  - 为什么需要虚拟 IP 呢，这就涉及到为什么要有 LVS 了，VIP 就是集群服务的一种体现，1.Cluster 集群，LB 负载均衡，HA 高可用.note 在这篇文章中第一段就是说明了集群的作用，为了让用户不用直接找 RS，而把所有的设备当做一个整体，用户看到的只有一个 IP，而不是那么多 RS 的 IP。

- **Real Server IP** # RIP,调度 IP，真实服务器 IP

## LVS Architecture(架构)

典型的 LVS 集群架构如图 所示。在 LVS 负载均衡集群架构中，尽管整个集群内部有多个物理节点在处理用户发出的请求，但是在用户看来，所有的内部应用都是透明的，用户只是在使用一个虚拟服务器提供的高性能服务，这也是 Linux 虚拟服务器项目，即 LVS 项目的主要名称来源，如下是对 LVS 集群架构中各个层次的功能描述。
![](https://notes-learning.oss-cn-beijing.aliyuncs.com/zzd89g/1622466145377-4664be5a-f6bd-4537-94eb-72054d13d096.png)

在基于 LVS 项目架构的服务器集群系统中，通常包含三个功能层次：

- **Load Balancer(负载均衡)** # 是整个集群系统的前端机器，在一组服务器之间均衡来自客户端的请求，让客户端认为所有服务都来自同一个 IP。
  - Director(调度器) 就是在该层工作
  - 负载均衡层位于整个集群系统的最前端，由一台或者多台 Director 组成， IPVS 模块就安装在 Director Server 的系统上，而 Director Server 的主要功能类似路由器，其包含了完成 LVS 负载转发功能所设定的路由表， Director 利用这些路由表信息把用户的请求分发到 Sever Cluster 层的物理服务器(Real Server) 上。此外，为了监测各个 Real Server 服务器的健康状况，在 Director Server 上还要安装监控模块 Ldirectord，当监控到某个 Real Server 不可用时，该服务器会被从 LVS 路由表中剔除，恢复时又会重新加入。
- **Server Cluster(服务器集群)** # 这是一组运行实际网络服务的服务器，如 Web、邮件、FTP、DNS 和媒体服务。
  - Real Server(真实服务器) 就是在该层工作
  - 服务器阵列或服务器池由一组实际运行应用服务的物理机器组成，Real Server 可以是 Web 服务器、Mail 服务器、FTP 服务器、DNS 服务器以及视频服务器中的一个或者多个的组合。每个 Real Server 之间通过高速的 LAN 或分布在各地的 WAN 相连接。在实际应用中，为了减少资源浪费， Director Server 也可以同时兼任 Real Server 的角色，即在 Real Server 同时部署 IPVS 模块。
- **Shared Storage(共享存储)** # 为服务器提供共享的存储空间，便于提供相同的服务。
  - 共享存储可以是数据库系统、网络文件系统或分布式文件系统。服务器节点需要动态更新的数据应该存储在基于数据的系统中，当服务器节点在数据库系统中并行读写数据时，数据库系统可以保证并发数据访问的一致性。静态数据通常保存在 NFS、CIFS 等网络文件系统中，以便所有服务器节点共享数据。但是，单个网络文件系统的可扩展性是有限的，例如单个 NFS/CIFS 只能支持 4 到 8 个服务器的数据访问。对于大型集群系统，分布式/集群文件系统可以用于共享存储，例如 GPFS，Coda 和 GFS，然后共享存储也可以根据系统需求进行扩展。

通常情况下，一个 LVS 集群由两类节点组成：

- **Director(指挥器)** # 前端接收客户端请求的节点，并将请求转发给后端 Real Server。Director 通过 IPVS 与 ipvsadm 来实现。
- **Real Server(真实服务器)**# 处理客户端请求的节点。

这些服务器一起组成了一个虚拟服务器，对于访问他们的客户端来说，它表现为一台机器。

LVS 的工作模式

LVS 的 IP 负载均衡技术是通过 IPVS 模块来实现的， IPVS 是 LVS 集群系统的核心软件，其主要安装在集群的 Director Server 上，并在 Director Server 上虚拟出一个 IP 地址，用户对服务的访问只能通过该虚拟 IP 地址实现。这个虚拟 IP 通常称为 LVS 的 VIP(Virtual IP)，用户的访问请求首先经过 VIP 到达 Director，然后由 Director 从 Real Server 列表中按照一定的负载均衡算法选取一个服务节点响应用户的请求。在这个过程中，当用户的请求到达 Director Server 后， Director Server 如何将请求转发到提供服务的 Real Server 节点，而 Real Server 节点又如何将数据返回给用户， 这是 IPVS 实现负载均衡的核心技术。

IPVS 实现数据路由转发的机制如下几种：

1. **NAT** # 支持端口映射但是 DIP 与 RIP 必须要在同一网段
2. **DR** # 不支持端口影响且调度与 RS 必须在同一网络
3. **TUN** # 各 RS 可以放在不同的地域且都在公网上被人直接访问
4. **FullNAT** # 可以在内部构建复杂网络，比如不同 RS 可以跨机房跨网络，而且可以隐藏 RS 不被公网直接访问)

NAT(Network Address Translation)

即通过网络地址转换的虚拟服务器技术。在这种负载转发方案中，当用户的请求到达调度器时，调度器自动将请求报文的目标 IP 地址（ VIP ）替换成 LVS 选中的后端 Real Server 地址，同时报文的目标端口也替换为选中的 Real Server 对应端口， 最后将报文请求发送给选中的 Real Server 进行处理。当 Real Server 处理完请求并将结果数据返回给用户时，需要再次经过负载调度器，此时调度器进行相反的地址替换操作，即将报文的源地址和源端口改成 VIP 地址和相应端口，然后把数据发送给用户，完成整个负载调度过程。可以看出，在这种方式下，用户请求和响应报文都必须经过 Director Server 进行地址转换，请求时进行目的地址转换（ Destination Network Address Translation, DNAT ），响应时进行源地址转换（ Source Network Address Translation, SNAT ）。在这种情况下，如果用户请求越来越多，调度器的处理能力就会成为集群服务快速响应的瓶颈。

LVS-NAT(Network Address Translation)实测可调度 10 台以内的 RS

多目标 IP 的 DNAT,通过将请求报文中的目标地址和目标端口改为某挑出的 RS 的 RIP 和 PORT 实现转发

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/zzd89g/1616132542821-5af6ee6a-a6d2-4e55-808d-cc00112912f7.jpeg)
lvs-nat 的特性

1. RS 应该使用私有地址

2. RS 的网关必须指向 DIP

3. RIP 和 DIP 必须在同一网段内

4. 请求和响应的报文都得经过 Director，在高负载场景中，Director 很可能成为性能瓶颈(因为既要处理请求报文也要处理响应服方的转发,请求报文一般很小,但响应报文一般都比较大)

5. 支持端口映射,即可修改请求报文的目标端口.

6. Director 必须是 Linux 系统，RS 可以是任意支持集群服务的操作系统.

lvs-nat 修改请求报文的目的 IP

1. 注意：该类型中 DIP 与 RIP 必须在同一网段且 RS 的网关为 DIP2，所以 Director 需要有两块网卡，DIP1 与 CIP 想通，DIP2 与 RS 想通。所有请求都经过调度器，包括请求报文和响应报文,调度器压力很大
2. 数据包到达 Director 时，做 dnat(将 VIP 改为 RIP)，然后发送给 RS。
3. RS 处理完数据包返回响应给 Director，源 IP 是 RIP，目标 IP 是 CIP
4. 这时候 Director 收到响应包后，做 snat(将源 IP 改为 VIP)

DR(Direct Routing)

即直接路由技术实现的虚拟服务器。这种技术在调度连接和管理上与 VSNAT 和 VSTUN 技术是一样的，不过它的报文转发方式与前两种均不同， VSDR 通过改写请求报文的 MAC 地址，将请求直接发送到选中的 Real Server ，而 Real Server 则将响应直接返回给客户端。因此，这种技术不仅避免了 VSNAT 中的 IP 地址转换，同时也避免了 VS TUN 中的 IP 隧道开销，所以 VSDR 是三种负载调度机制中性能最高的实现方案。但是，在这种方案下， Director Server 与 Real Sever 必须在同一物理网段上存在互联。

LVS/DR(Direct Routing 直接路由) 实测可以调度 7、8 十台 RS

通过为请求报文重新封装一个 MAC 首部进行转发,源 MAC 是 DIP 所在的接口的 MAC,目标 MAC 是某挑选出的 RS 的 RIP 所在接口的 MAC 地址;源 IP/PORT,以及目标 IP/PORT 均保持不变,请求报文经过 Dirctor 但响应报文不再经过 Dirctor
![](https://notes-learning.oss-cn-beijing.aliyuncs.com/zzd89g/1616132542832-728aa1cf-07b1-47b3-8cdf-d8918c6789e6.jpeg)

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/zzd89g/1616132542795-f5e7f61e-e8c9-4b6e-bcd5-6c038fb30f7f.jpeg)

### DR 类型工作流程

如上图所示：当客户端请求 VIP 时，通过互联网到达前端路由 Route1，再通过交换机到达 Dirctor 上；而 Dirctor 在收到请求报文后，通过定义在 ipvs 规则中的各 rip 去获得各 RS 的 MAC 地址，并在此报文外再封装一个 MAC 地址,源 MAC 为 Dirctor 的 DIP 端口的 MAC 而目标 MAC 改为其中被调度算法选中一个 RS 的 MAC，但该报文的目标 ip(VIP)不变，最后通过 DIP 接口发送给 RS；为了 RS 能接收 Dirctor 发来的报文，需要在各 RS 上也配置 VIP，但 RS 上的 VIP 是需要隔离前端 arp 广播的，所以需要将各 RS 上的 VIP 隐藏（RS 上的 VIP 通常配置到 lo 网卡接口的别名上，并配合修改 Linux 内核参数来实现隔离 arp 广播）；而 RS 封装响应报文时，源 IP 为 VIP，目标 ip 为 CIP，并通过 RIP 的网络接口直接向外发送响应，不再经过 Dirctor。

需要注意的是：因为 Route1 的 A 点的 IP 和 Dirctor 的 VIP 在同一网段内，VIP 通常是公网 IP；而 DIP、RIP 通常是私有 IP，且这两个 IP 通常也应在同一物理网络内；假设 RIP 与 Route1 的 A 接口(同 Director 的 VIP DIP)在同一网段，则这时可将 RS 的网关指向 Route1，否则，Route2 只能其它路由器(如 Route2)接口访问互联网，且 Route2 的 C 点的 IP 需要与 RIP 在同一网段内，此时 RIP 响应的报文就通过 Route2 发送。

- 1.RIP 与 DIP 在同一 IP 网络，RS 可以使用私有地址，也可以使用公网地址，此时可以直接通过互联网连入 RS，以实现配置、监控等
- 2.RS 的网关一定不能指向 DIP
- 3.RS 跟 Director 要在同一物理网络内（不能有路由器分隔,因为要将报文封装 MAC 首部进行报文转发）
- 4.请求报文必须经过 Director，但响应报文不能经过 Director 而是由 RS 直接发往 Client 以释放 Directory 的压力。
- 5.不支持端口映射(因为响应报文不经过 Director)
- 6.RS 可以使用大多数的操作系统
- 7.Director 的 VIP 对外可见，RS 的 VIP 对外不可见
- 8.RS 跟 Director 都得配置使用 VIP
- 9.确保前端路由器将目标 IP 为 VIP 的请求报文发往 Director(上文的设置)

lvs-dr(direct routing) # 操纵新的 MAC 地址，直接路由，默认的 LVS 类型，通过请求报文的目标 MAC 地址进行转发，即需要 ARP 的 IP 与 MAC 映射表才能转发，由于调度器是基于二层 MAC 来调度的，所以调度器与 RS 必须要在同一个 VLAN 中
![](https://notes-learning.oss-cn-beijing.aliyuncs.com/zzd89g/1616132542782-786ee136-62a7-4fb6-8faf-2a43c660d620.jpeg)

1. 如图所示，请求报文直接到调度器，然后调度器选择一台 RS，让这台 RS 来响应该请求，RS 与用户直接交互，不再经过调度器，只有请求报文经过调度器，响应报文是不经过调度器的(RS 的网关不能指定到 DIP)所以用户访问的业务 IP 都是该业务的调度器的浮动 IP，通过调度器来给用户选择一台提供服务的主机，这样调度器没有压力。
2. RS 构建响应报文响应用户请求的时候，需要使用 VIP 来进行响应，因为用户请求的是 VIP，RS 只能用 VIP 来进行响应。每台 RS 都有一个 VIP，但是用户请求的 VIP 必须要到调度器上，那么这时候在 VIP 上就需要绑定 mac 地址了，以 mac 地址来区分调度器与 RS。
3. 所以 Director 在调度主机的时候，会把选择的 RS 的物理网卡的 MAC 地址加进请求报文中作为目的 mac 地址，然后转发给 RS。（用户依然会通过 VIP 来访问，但是数据到交换机的时候，是根据 MAC 地址来转发该数据到 RS，这样就实现了 RS 与用户的直接交互；所以当用户断开连接再次请求后，由于请求报文目的 MAC 地址没了，所以交换机会根据本身的 arp 表把，找到 mac 地址，这时候依然是 Director，则数据包到 Director 的时候需要重新配分新的 RS）
4. 由于 Linux 从哪网卡收的就要从哪个网卡发，为了解决响应报文中的源地址不能是 RS 的 IP 的问题，那么每台 RS 的 VIP 则不配置在物理网卡上，而是配置在这台机器 lo(loopback)接口上，给 lo 起一个别名用作 VIP，具体原因如下
   1. 由于路由器一般是动态学习 ARP 包的（一般动态配置 DHCP 的话），当内网的机器要发送一个到外部的 ip 包，那么它就会请求 路由器的 Mac 地址，发送一个 arp 请求，这个 arp 请求里面包括了自己的 ip 地址和 Mac 地址，而 linux 默认是使用 ip 的源 ip 地址作为 arp 里面 的源 ip 地址，而不是使用发送设备上面的 ，这样在 lvs 这样的架构下，所有 RS 发送包都是同一个 VIP 地址，那么 arp 请求就会包括 VIP 地址和设备 Mac，而路由器收到这个 arp 请求就会更新自己的 arp 缓存，这样就会造成 ip 欺骗了，VIP 被抢夺，这样调度器的 VIP 就被被 RS 抢走，这样就会出现问题，下一个数据包就无法正确发送给调度器了。所以需要给每台 RS 都要配置 arp 的通告以及响应规则以实现该功能
   2. **arp_ignore** # arp 忽略，响应 arp 请求时的动作，由该项参数决定
      1. 参数 0：默认参数，只要收到 arp 请求，无论是哪个地址，都做出响应
      2. 参数 1：推荐设置，只响应目的 IP 地址为接收网卡上的本地地址的 arp 请求。(如果外面有人询问 VIP，由于 VIP 在 lo 上，不在接收网卡上，所以不会响应询问 VIP 的 arp 通告。这样防止 VIP 的 arp 被抢夺)。
      3. 参数 2：只响应目标 IP 地址是来访网络接口本地地址的 ARP 查询请求,且来访 IP 必须在该网络接口的子网段内
      4. 参数 3：不回应该网络界面的 arp 请求，而只对设置的唯一和连接地址做出回应
      5. 参数 4-7：保留未使用
      6. 参数 8：不回应所有（本地地址）的 arp 查询
   3. **arp_announce**# arp 宣告，通告 arp 给别人时的动作；以及是否接收 arp 通告，并记录；由该项参数决定
      1. 参数 0：默认参数，把本机所有网卡上的所有地址通告给网络中（不管任何情况，使用发送或者转发的数据包的源 IP 作为发送 ARP 请求的 Sender IP）(与参数 1 的区别：不管目的 IP 与本地接口的 IP 在不在同一个网段，都是用发送源 IP 作为 Sender IP）
      2. 参数 1：尽量避免从本网络的外部接口向非本网络中的网络,通告非本网络中的接口的地址（只有当数据包的目的 IP 与本地某个接口的网段相同时，才使用发送或者转发的数据包的源 IP 作为发送 ARP 请求的 Sender IP，不属于则按参数 2 处理）
         1. 本网络的意思就是：比如 192.168.0.0/24 是一个网络，192.168.1.0/24 是另一个网络，0.0 网络中的地址尽量不通告给 1.0 网络中的地址，但是当需要发送数据的时候，还是需要进行通告
      3. 参数 2：推荐设置，在发送 arp 宣告的时候不使用数据包的源 IP，使用能与目标主机通信的最佳地址来作为发送 ARP 的 Sender IP，优先选择对外接口的主 IP；（loopback 不是对外接口）(e.g.在 RS 给 client 发送响应数据包的时候，默认情况下，会先给发送 arp 通告，询问网关在哪。由于数据包的源 IP 是 VIP，MAC 是发送数据包的物理网卡的 MAC；目的 IP 是 client 的 ip，目的 mac 未知，所以 arp 通告的源 IP 也是 VIP，那么这时候，交换机就会更新 VIP 与 MAC 的对应关系，此时产生问题，因为 VIP 应该与 director 的 MAC 绑定才对，但是现在收到的这个 arp 通告说是 VIP 应该与 RS 绑定，这明显是不应该发生的。所以在发送数据包之前的 arp 通告，不能使用 VIP，而是使用本机的物理网卡来进行 arp 通告。不过这个从 RS 发出的数据包的封包其实源 IP 还是 VIP、源 MAC 是物理网卡的 MAC，这样在 client 收到 RS 的响应包之后与 RS 交互发送数据包，目的 IP 则是 VIP，目的 MAC 则是 RS 物理网卡的 MAC，当交换机收到 client 发的包时，解开封包看到目的 mac 地址是 RS 的，则直接就把数据包交给对应的网口了，至于 IP 则是在三层路由的时候才用的，当交换机已经收到这个包时，就会把 IP 拆开直接使用 MAC 来传输数据包。而当 RS 与 client 断开连接后，client 再次主动发的数据包到交换机时，目的 MAC 是未知的，交换机就会把数据包交给 director 来进行处理，因为交换机的 arp 表里已经把 VIP 与 director 的 mac 绑定了)。arp 原理详见 ARP.note
         1. 三者 ARP 的通告规则区别：参数 0 是不管什么时候把所有 IP 都通告，参数 1 是不同网段需要通信的时候才通告，有死亡时间，过一段时间，该 ARP 表自动消失，参数 2 是使用最优的 IP 进行 ARP 通告，不是对外接口(比如 loopback)的永不通告
5. 所以，为了满足 dr 类型的需要，arp_ignore 设置为 1（RS 响应 arp 通告的时候 VIP 不在接收 arp 请求这个接口上就不会响应），arp_announce 需要设置为 2(RS 在发送 arp 通告的时候不使用 VIP 作为源 IP)，该配置为内核参数配置，在/proc/sys/net/ipv4/conf/all 目录和/proc/sys/net/ipv4/conf/lo 目录下的两个文件进行配置
6. 再次注意：进行 RS 配置的时候，需要先修改 arp 的配置，再配置 lo 的 VIP。否则如果直接配置 VIP，则会使用 lo 的 VIP 来响应询问 VIP 在哪的 arp 通告，这时候 VIP 与 MAC 的对应关系就会一直变化，这样立刻就会发生 arp 抢夺

## TUN(IP Tunneling)

即 IP 隧道技术实现的虚拟服务器。VS TUN 与 VSNAT 技术的报文转发方法不同，在 VS TUN 方式中，调度器采用 IP 隧道技术将用户请求转发到某个选中的 Real Server 上，而这个 Real Server 将直接响应用户的请求，不再经过前端调度器。此外， IP TUN 技术对 RealServer 的地域位置没有要求，其既可以与 Director Server 位于同一个网段，也可位于独立网络中。因此，在 VS TUN 方式中，调度器将只处理用户的报文请求，而无需进行转发， 故集群系统的响应速率相对而言得到极大提高。

LVS/TUN
模型：在原请求 IP 报文之外新加一个 IP 首部(这个新添加的 IP 首部其源 IP 是 DIP,目标 IP 是 RIP),将报文发往挑选出的目标 RS.
![](https://notes-learning.oss-cn-beijing.aliyuncs.com/zzd89g/1616132542822-64d36467-fa78-48c8-af3c-ea93ccfeea8b.jpeg)
TUN
TUN 类型工作流程(主要是为了容灾,因为 Director 与各 RS 是在不同网段中,所以可以存在于不同的物理空间)

LVS TUN 类型特性

1.RIP，DIP，VIP 都得是公网地址

2.RS 的网关不会指向也不可能指向 DIP

3.请求报文经过 Directory，但响应报文一定不经过 Director

4.不支持端口映射

5.RS 的 OS 必须得支持隧道功能 ??

2.4 LVS/FULLNAT (LVS 默认不支持此类型)

通过同时修改请求报文的源 IP 地址和目标 IP 地址进行转发

报文件从客户端到调度器时的源 目标 IP : CIP-->VIP

报文从 Director 到 RS 时的源 目标 IP 分别是: DIP-->RIP

特点

1. VIP 是公网地址,RIP 和 DIP 是私网地址,且通常不在同一 IP 网络,因此,RIP 的网关一般不会指向 DIP
2. RS 收到的请求报文源地址是 DIP,因此只需响应给 DIP, 但 Dirctor 还要将其发往 Client
3. 请求和响应报文都经由 Dirctor
4. 支持端口映射.

lvs-tun(ip tunneling) # 在原请求 IP 报文之外新加一个 IP 首部，IP 隧道技术

1. 不修改请求报文的 IP 首部，通过在原有的 IP 首部之外，再封装一个 IP 首部(比如为了运送一袋米，我扛着米运相当于 IP 首部，我骑着驴送，驴相当于新封装的 IP 首部)
2. 调度器收到请求报文时，再封装一层 IP 首部，把源 IP 至目标 IP 中的的 CIP 至 VIP 放在 DIP 至 RIP 的报文里面，DIP 至 RIP 相当于路由，所以可以不在同一网段,并且 RS 必须支持隧道技术，在解封装的时候，必须明白为什么在拆开 IP 首部之后还有一层 IP 首部
3. 不支持端口映射，且 RS 的网关不能指向 DIP。

lvs-fullnat # 修改请求报文的源和目标 IP

1. 调度器同时修改请求报文的目标地址和源地址进行转发。把源 IP 和目标 IP 从 CIP—VIP 改成 DIP—RIP
2. VIP 是公网地址，RIP 和 DIP 是私网地址，二者无须在同一网络中
3. RS 接收到的请求报文的源地址为 DIP，因为要响应给 DIP
4. 调度器一样要承担很大压力

## Note

1. 为什么在 Director 上除了 DIP 本身还需要一个单独的 VIP，而不可以把 DIP 当做 VIP 来用
   1. 如果 DIP 与 VIP 一样，那么在 Director 发送 arp 广播的时候，RS 在收到 arp 广播后，回应的报文会回给自己(因为 RS 设备上 lo 网卡上的 VIP 就是 DIP)，这样调度器上的 arp 表里就无法获得后端 RS 的 MAC 地址。如果多了一个 VIP，那么 RS 在回应 Director 的 arp 广播时，就不会回应到自己身上，因为 RS 的 lo 网卡上的 IP 为单独的 VIP，而不同于 DIP。
   2. 结论：在调度器上，除了本身的 DIP 以外，必须要一个 VIP
2. http 本身是 stateless 无状态的，无法追踪目标来源
   1. session 保持机制：会话保持机制，保证 http 协议可以在用户终端连接再次连接后还能存有之前的操作记录
      1. session 绑定：将来自于同一个 client 的请求始终绑定在一个 RS 上，不会被调度到别的 RS 上
      2. session 集群:
      3. session 服务器

# Scheduling(调度) 方法

LVS 的调度方法分为两类(静态算法、动态算法)，共 10 种

在转发方式选定的情况下，采用哪种调度算法将决定整个负载均衡的性能表现。不同的算法适用于不同的生产环境，有时可能需要针对特殊需求自行设计调度算法。

1. 静态方法：仅根据算法本身进行调度（注重起点公平）
   1. RR(Round Robin) # 轮询，论调，轮流调度，第一个请求给 RS1，第二个请求给 RS2，第 n 个请求给 RSn，第 n+1 个请求给 RS1。。。。。
   2. WRR(Weighted RR) # 加权(Weight)轮询，能者多劳，给 RS1 一个请求，就给 RS2 几倍的请求
   3. SH(Source hash) # 源地址哈希，实现 session 保持的机制,来自同一个 IP 的请求将始终调度到同一个 RS
   4. DH(Destination Hash) # 目标地址哈希，只要请求的是同一个资源，则将请求调度到同一个 RS,比如 CDN 中所有用户都请求一个资源被调度到一台 RS 上
2. 动态方法：根据算法以及各 RS 的当前负载状态进行调度 Overhead
   1. LC(Least Connection) # 最小连接数，新来的请求报文调度给连接数最小的 RS
      1. Overhead=Active\*256+Inactive
   2. WLC(Weighted LC) # 加权(Weight)最小连接数 默认的调度器类型
      1. Overhead=(Active\*256+Inactive)/weight
   3. SED(Shortest Expection Delay) # 最短期望延迟
      1. Overhead=(Active+1)\*256/weight
   4. NQ(Never Queue) # SED 算法的改进
   5. LBLC(Locality-Based LC) # 基于本地的最小连接数，动态的 DH 算法，正向代理情形下的 cache server 调度
   6. LBLCR(Locality-Based LC with Replication) # 带复制功能的 LBLC，相当于几台 RS(HCS)中的资源可以互相共享

# LDirectorD 技术介绍，以及产生的原因

IPVS 有一个缺陷，无法检查后端 Real Server 的健康状态，就是使用 HA 给 LVS 中的 Director 实现了高可用，也不一定能保证后端的各 RS 可以正常响应用户的请求，当其中一台 RS 不能使用时，访问 vip，还会去调度 down 掉的这台 RS，并返回一个错误的页面。这种情况是不合理的，所以我们需要一个 LVS 的健康检查机制，以便当 RS 无法响应时，可以及时通知给 Director，让其不再把请求调度给这台坏掉的 RS 上。为了实现这个功能，就用到了 heartbea 中的 ldirectord，ldirectord 以守护进程运行在后台，提供生成规则以及 Health check 健康检查
![](https://notes-learning.oss-cn-beijing.aliyuncs.com/zzd89g/1616132542796-003ca683-5cd4-4f31-85ec-6b92628f2017.jpeg)
该程序依赖于自己的配置文件生成 ipvs 规则，因此，定义集群服务、添加 RS、调度方法等都在配置文件中指定，而无须手动用 ipvsadm 命令更改

## 配置

/etc/ha.d/ldirectord.cf

## ldirectord.cf 文件说明

**Global Directives** # 全局指令，对所有的 Virtual Services 都有效

- checktimeout=3 # 检查的超时时间，当对 RS 的健康检查时间超过 3 秒的时候的则认为该 RS 不可用
- checkinterval=1 # 检查时间间隔，即每 1 秒都对后端 RS 进行一次健康检查
- \#fallback=127.0.0.1:80 #
- autoreload=yes # 该配置文件是否自动装载
- \#logfile="/var/log/ldirectord.log" # 指明日志文件的 PATH
- \#logfile="local0" #
- \#emailalert="admin@x.y.z" # 警告信息发送的邮箱地址
- \#emailalertfreq=3600 # 每隔多久发送一次警告信息到邮箱
- \#emailalertstatus=all # 通知的 email 信息是全部
- quiescent=no # 静默工作模式

**Sample for an XXXXX** # 对于多种虚拟服务的配置样例，直接修改这一部分内容，可以实现健康检查的基本模式，其中前三行为必须要定义的 LVS 的定义以及调度规则，剩下的所有行定义的都是为 ldirectord 对后端 RS 的健康检查方式，当这些健康检查方式失败的时候，则说明该 RS 不可用

- virtual=IP:PORT # 定义 VIP 的地址和端口
- real=IP\[\[->IP]:\[PORT]] TYPE # 定义 RS 的 IP 地址和 LVS 类型，类型名介绍详见 LB 的 Packet-Forwarding-Method(LVS Type)内容，其中->IP 可以实现从哪个 IP 至哪个 IP 的地址段的定义
  - gate # TYPE 为 DR 类型
- fallback=IP:PORT TYPE # 定义当 RS 全部失效时，使用的 server 的地址，端口，LVS 类型。
- scheduler=SCHEDULER # 定义 LB 集群中的调度规则，规则类型详见 LB 中的 Director 调度方法
- service=TYPE # 定义健康检查的应用层 Protocol，注意：只有当 checktype 指定为 negotiate 的时候，该定义才有意义
  - TYPE 类型包括：ftp|http|stmp|mysql 等
- protocol=tcp # 定义健康检查的传输层 Protocol
- checktype=negotiate # 定义健康检查的方法
  - connect # 传输层检查，向对方端口尝试发送连接请求
  - negotiate # 应用层检查协商方法
  - ping # 网络层检查，ICMP 协议
- checkport=80 # 定义健康检查的端口号
- request="index.html" # 定义健康检查请求目标 server 的哪个页面
- receive="Test Page" # 定义健康检查中 request 中所定义的页面请求后回复的内容包含什么信息
- virtualhost=www.x.y.z # 定义健康检查虚拟主机的主机名
