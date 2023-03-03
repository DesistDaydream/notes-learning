---
title: Flannel
---

# 概述

> 参考：
> - [GitHub 项目，flannel-io/flannel](https://github.com/flannel-io/flannel)

Flannel 是一种专为 Kubernetes 设计的，简单、易于配置的 3 层网络结构，并且为 Kubernetes 提供了 CNI 插件。

支持多种后端：即使用什么方法进行进行数据的接收与发送

- vxlan
- host-gw：host gateway
- UDP

Flannel 在每台主机上运行一个名为 flanneld 的小型二进制程序作为代理，负责从更大的预配置地址空间中为每个主机分配 **subnet lease(子网租期)**。Flannel 直接使用 Kubernetes API 或 etcd 来存储网络配置、已分配的子网、以及任何辅助数据(比如主机的 IP)

## 子网获取逻辑

代码：`./main.go —— WriteSubnetFile()`
Flannel 启动时，在 `./main.go` 中调用 [WriteSubnetFile()](https://github.com/flannel-io/flannel/blob/v0.15.1/main.go#L746) 函数，用来生成 subnet 配置文件(默认在 /run/flannel/subnet.env)。

```go
func main() {
    ......
	if err := WriteSubnetFile(opts.subnetFile, config, opts.ipMasq, bn); err != nil {
		// Continue, even though it failed.
		log.Warningf("Failed to write subnet file: %s", err)
	} else {
		log.Infof("Wrote subnet file to %s", opts.subnetFile)
	}
    ......
}

func WriteSubnetFile(path string, config *subnet.Config, ipMasq bool, bn backend.Network) error {
	......
    nw := config.Network
    sn := bn.Lease().Subnet
    // Write out the first usable IP by incrementing sn.IP by one
    sn.IncrementIP()
    fmt.Fprintf(f, "FLANNEL_NETWORK=%s\n", nw)
    fmt.Fprintf(f, "FLANNEL_SUBNET=%s\n", sn)

	fmt.Fprintf(f, "FLANNEL_MTU=%d\n", bn.MTU())
	_, err = fmt.Fprintf(f, "FLANNEL_IPMASQ=%v\n", ipMasq)
    ......
}
```

其中 `bn.Lease().Subnet` 就是当前节点被分配的子网，该子网通过 `./subnet/subnet.go` 中的 `Manage()` 接口中的 `AcquireLease()` 方法获取的。后续 Flannel 会根据该文件的内容，为节点生成路由规则。

具体 Lease 的获取方法，根据配置存储方式来决定，现阶段主要使用的是直接通过 Kubernetes API 的方式，那么就主要看 `./subnet/kube/kube.go` 中的 `ksm.AcquireLease()` 方法。其中最主要的就是这一部分：

```go
type kubeSubnetManager struct {
	nodeName                  string
	nodeStore                 listers.NodeLister
    ......
}

func (ksm *kubeSubnetManager) AcquireLease(ctx context.Context, attrs *subnet.LeaseAttrs) (*subnet.Lease, error) {
    cachedNode, err := ksm.nodeStore.Get(ksm.nodeName)
	n := cachedNode.DeepCopy()

	var cidr, ipv6Cidr *net.IPNet
	_, cidr, err = net.ParseCIDR(n.Spec.PodCIDR)
	if err != nil {
		return nil, err
	}

	if cidr != nil {
		lease.Subnet = ip.FromIPNet(cidr)
	}

	return lease, nil
}
```

其中 `cidr` 变量就是 Flannel 为当前节点分配的子网，而 cidr 的值来自于 `n.Spec.PodCIDR`，这个是直接通过 kubernetes 的 client-go 调用的 Kubernetes 集群内部给 Node 分配的 CIDR。

> 这个 CIDR 通过 kube-controller-manager 的 `--cluster-cidr` 标志配置
> 并且可以通过 `kubectl get node XXX -ojsonpath='{.spec.podCIDR}'` 命令从 Kubernetes 集群中获取到

这也就是为什么 Flannel 的 CIDR 配置与 kube-controller-manager 的 CIDR 配置要保证一致的原因。

此时，Flannel 就知道如何在当前阶段创建所需的路由了。

> Flannel 中获取 Node 信息的方式与我们平时使用的方式不太一样，并不是 `clientset.CoreV1().Nodes().List(context.TODO(), v1.ListOptions{})` 这种通用方式，具体原因未知。Flannel 中通过 `listers.NodeLister()` 接口的 `Get()` 方法获取信息。而 `listers.NodeLister()` 接口的实例化，则是在 `newKubeSubnetManager()` 函数中进行的。

# 三种模型

## UDP 模型

![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/lmmp31/1616118603039-02bb30c5-6784-4446-98ac-a3fe86c48c28.png)

### 数据流的走向

- 用户态的容器进程发出的 IP 包经过 docker0 网桥进入内核态；
- IP 包根据路由表进入 TUN（flannel0）设备，从而回到用户态的 flanneld 进程；
- flanneld 进行 UDP 封包之后重新进入内核态，将 UDP 包通过宿主机的 eth0 发出去。

Flannel 进行 UDP 封装（Encapsulation）和解封装（Decapsulation）的过程，也都是在用户态完成的。在 Linux 操作系统中，上述这些上下文切换和用户态操作的代价其实是比较高的，这也正是造成 Flannel UDP 模式性能不好的主要原因。

所以说，我们在进行系统级编程的时候，有一个非常重要的优化原则，就是要减少用户态到内核态的切换次数，并且把核心的处理逻辑都放在内核态进行。这也是为什么，Flannel 后来支持的 VXLAN 模式，逐渐成为了主流的容器网络方案的原因

Pod 间通信的情况：

- pod1 与 pod2 不在同一台主机
  - pod1(10.0.14.15)向 pod2(10.0.5.150)发送 ping，查找 pod1 路由表，把数据包发送到 cni0(10.0.14.1)
  - cni0 查找 host1 路由，把数据包转发到 flannel.1
  - flannel.1 虚拟网卡再把数据包转发到它的驱动程序 flannel
  - flannel 程序使用 VXLAN 协议封装这个数据包，向 api-server 查询目的 IP 所在的主机 IP,称为 host2(不清楚什么时候查询)
  - flannel 向查找到的 host2 IP 的 UDP 端口 8472 传输数据包
  - host2 的 flannel 收到数据包后，解包，然后转发给 flannel.1 虚拟网卡
  - flannel.1 虚拟网卡查找 host2 路由表，把数据包转发给 cni0 网桥，cni0 网桥再把数据包转发给 pod2
  - pod2 响应给 pod1 的数据包与 1-7 步类似
- pod1 与 pod2 在同一台主机
  - pod1 和 pod2 在同一台主机的话，由 cni0 网桥直接转发请求到 pod2，不需要经过 flannel。

## VxLan 型后端

VXLAN，即 Virtual Extensible LAN(虚拟可扩展局域网)，是 Linux 内核本身就支持的一种网络虚拟化技术。VXLAN 可以在内核中实现上面 UDP 模型中 flanneld 进程的封装和解封装工作，从而通过与前面相似的隧道机制，构建出叠加网络(Overlay NetworkTunnel 隧道技术与 overlay 叠加网络.note)。Overlay 技术实际上是一种隧道封装技术，一个数据包(内部)封装在另一个数据包内(也就是本身的 IP 外面又套了一个 IP);被封装的包转发到隧道端点后再被拆装。原来的包就发送到了目的地。叠加网络就是使用这种所谓“包内之包”的技术安全地将一个网络隐藏在另一个网络中，然后将网络区段进行迁移。

VXLAN 的设计思想是：在现有的三层网络之上，覆盖一层虚拟的，由内核 VXLAN 模块负责维护的二层网络，使得连接在这个二层网络上的主机(虚拟机或者容器都可以)之间，可以像在同一个局域网里那样通信。虽然这些主机可能分布在不同的宿主机上，甚至是分布在不同的物理机房里。这里面的二层是逻辑上的而是，是指的在

为了能在二层网络上打通隧道，VXLAN 会在 Host 上设置一个特殊的网络设备作为隧道的两端，这个设备叫做 VTEP(VXLAN Tunnel End Point 虚拟隧道端点)。VTEP 设备的作用，与 UDP 类型的 flanneld 进程非常相似，只不过 VTEP 进行封装和解封装的对象是二层数据帧，而且这个过程是在内核里完成(因为 VXLAN 本身就是 Linux 内核中的一个模块)。而 VTEP 设备之间的交互，就是二层交互(不同网段可以理解为不同 VLAN 的交互)，想让他们互通，则是 VXLAN 模块来实现的。

下图是基于 VTEP 设备进行隧道通信的流程，每台宿主机上的 flannel.1 的设备就是 VTEP 设备，既有 IP，也有 MAC
![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/lmmp31/1616118613885-1f230a8b-ecaa-4ffb-a094-d331200fa1e3.png)
flannel 会维护这么几个数据：

- 宿主机到对端 VTEP 设备的路由信息。记录要访问的其他宿主机上的容器的网段，下一跳是对端宿主机的 VTEP 设备的 IP，要通过 flannel.1 设备
- 宿主机 flannel.1 设备学习到的 arp 记录。记录其他宿主机的 VTEP 设备的 IP 与 MAC 对应关系
- 桥设备的 FDB 数据。记录要发送到其他宿主机的 VTEP 设备需要经过哪个宿主机的 IP(通过 VTEP 设备的 MAC 地址来确定)

该模型与 GRE 协议类似，10.168.0.2 就是 GRE 的外网网卡，10.244.0.0 就是 tun 设备的 IP，只不过 GRE 是人为规定好了对端的 IP，而该文章将说的是由 flannel 来维护整个叠加网络的信息，可以不仅仅局限于两台设备。

### 数据流的走向

pod1 与 pod2 不在同一台主机

- Container-1 发出请求后，目的地址是 10.244.1.3，经过 cni0，然后被路由到 flannel.1(VTEP) 设备进行处理。flannel 会根据所规定的子网，自动生成路由，让所有符合其子网的目的地址，都会经过 flannel.1 设备。可以把 container-1 发出的数据包称为“原始 IP 包”。这个“原始 IP 包”到达 flannel.1 设备，也就是来到了隧道的入口。此时开始了 VXLAN 的封装工作
  - 添加 Inner Ethernet Header。为了能够将“原始 IP 包”封装并发送到正确的宿主机上，VXLAN 就需要找到这条隧道的出口(i.e.目的宿主机的 VTEP 设备)，这个设备的信息，就是由每台宿主机的 flanneld 进程维护的。当 node2 启动并加入 flannel 网络后，node1 上会添加一条路由：10.144.1.0/24 via 10.244.1.0 dev flannel.1。这个 10.244.1.0 就是 node2 上的 VTEP 设备的 IP 地址。（可以把 node1 的 VTEP 设备成为“源 VTEP 设备”，node2 的 VTEP 设备成为“目的 VTEP 设备”）。这些 VTEP 设备之间，就需要想办法组成一个虚拟的二层网络。flanneld 进程在 node2 节点启动时，还会在 node1 上添加 arp 记录，记录“目的 VTEP 设备”的 IP 与 MAC(假设为 5e:f8:4f:00:e3:37，可以通过 ip neigh show dev flannel.1 命令查看)。有了这个“目的 VTEP 设备”的 MAC 地址，VXLAN 就可以在内核开始二层封包工作了，VXLAN 模块会在“原始 IP 包”外添上一个“目的 VTEP 设备”的 MAC 地址(Inner Ethernet Header)。但是只有一个 MAC 地址，对于宿主机网络来说没有实际意义，并不能在宿主机的网络里传输，所以需要进一步封装，让其成为宿主机网络里一个普通的数据包，以便通过 eth0 网卡。
  - 添加 VXLAN Header。为了让数据包可以变成宿主机网络里的普通数据包，VXLAN 模块会再给数据包加上一个特殊的 VXLAN 头(VXLAN Header)，用来表示这个数据包实际上是一个 VXLAN 要使用的包。而这个 VXLAN 头里有一个重要的标志，叫做 VNI，它是 VTEP 设备识别某个数据帧是不是应该归自己处理的标志。而在 flannel 中，VNI 的默认值为 1，所以宿主机上的 VTEP 设备都叫做 flannel.1 的原因，这里面的 1，就是 VNI 的值。(其实，添加这个 VXLAN 头，就是为了让宿主机在看到这个数据包是由 VXLAN 程序来发出的，而不是由 container 发出的，因为宿主机无法才开 VXLAN 的头部信息，所以也就读不了 VXLAN 下面的真实目的 IP 和 MAC)(说白了，可以把实现 VXLAN 功能的 flannel 当做一个运行在 linux 上的程序，数据包是由这个 flannel 发出来的。其余 VXLAN 的机制也是同理)
  - 添加 Outer UDP Header。然后 Linux 内核会把这个数据帧封装进一个 UDP 包，跟 UDP 模型一样，在宿主机看来，会认为自己的 flannel.1 设备只是在向外另外一台宿主机的 flannel.1 设备，发起了一次普通的 UDP 链接，并不会知道这个 UDP 包里，还有一个完成的二层数据帧（从宿主机看，就是 vxlan 这个模块或者说 flannel.1 设备，发送了一份数据，数据内容是什么，Linux 内核不关心）。不过，flannel.1 设备知道另一端设备的 MAC 地址，但是却不知道对应的宿主机地址是什么，那么这个 UDP 包应该发给哪台宿主机呢？
  - 添加 Outer IP Header。在这种情况下，flannel.1 实际上扮演了一个网桥的角色，网桥设备进行转发的依据，来自于一个 FDB(Forwarding Database)的转发数据库，这个 FDB 的信息也是由 flanneld 进程维护的，当 node2 加入 flannel 网络后，会在 node1 的 FDB 记录对端 VTEP 的信息 5e:f8:4f:00:e3:37 dev flannel.1 dst 10.168.0.3 self permanent(可以通过 bridge fdb show flannel.1 | grep 5e:f8:4f:00:e3:37 命令查到,意思是：MAC 地址为“目的 VTEP”设备的数据包，会经过 flannel.1 设备，发送到目的地是 10.168.0.3 的主机上)。所以接下来的流程就是一个正产的宿主机网络上的封包工作，flannel.1 设备会把 FDB 的信息告诉 Linux 内核要发送个谁，Linux 内核的网络栈就会进行后续封装，把对端 VTEP 设备的 MAC 地址所在的 IP 封装到数据包的头部。
  - 添加 Outer Ethernet Header。Linux 内核在这个数据包前面加上 Node2 的 MAC 地址，这个 MAC 是本身设备网络栈 ARP 表要学习到的，无需 flannel 维护。
  - 这时候，封包工作完成了。实际上就是 flannel.1 设备发送了一个数据给宿主机，至于数据中的内容，则是宿主机不关心的。当对端宿主机把最外层的封装解开后，发现 VXLAN 的标记，自然会交由本机可以处理 VXLAN 的网络设备来进行处理。
- Node1 上的 flannel.1 设备把封装好后的数据帧从 node1 的 eth0 网卡发出去
- node2 收到数据帧后，拆开发现 VXLAN 头，根据 VNI 值交给本地的 flannel.1 设备，flannel.1 设备进一步拆包获取“原始 IP 包”，并把该包送入对应的 Container-2 中。
- Container-2 的响应，与前面的描述一样，只不过是从 node2 开始封装，到 node1 后解封装

pod1 与 pod2 在同一台主机

- pod1 和 pod2 在同一台主机的话，由 cni0 网桥直接转发请求到 pod2，不需要经过 flannel。

pod 到 service 的网络

- 创建一个 service 时，相应会创建一个指向这个 service 的域名，域名规则为{服务名}.{namespace}.svc.{集群名称}。之前 service ip 的转发由 iptables 和 kube-proxy 负责，目前基于性能考虑，全部为 iptables 维护和转发。iptables 则由 kubelet 维护。
  - pod1 向 service ip 10.16.0.10:53 发送 udp 请求，查找路由表，把数据包转发给网桥 cni0(10.0.14.1)
  - 在数据包进入 cnio 网桥时，数据包经过 PREROUTING 链，然后跳至 KUBE-SERVICES 链
  - KUBE-SERVICES 链中一条匹配此数据包的规则，跳至 KUBE-SVC-TCOU7JCQXEZGVUNU 链
  - KUBE-SVC-TCOU7JCQXEZGVUNU 不做任何操作，跳至 KUBE-SEP-L5MHPWJPDKD7XIFG 链
  - KUBE-SEP-L5MHPWJPDKD7XIFG 里对此数据包作了 DNAT 到 10.0.0.46:53，其中 10.0.0.46 即为 kube-dns 的 pod ip
  - 查找与 10.0.0.46 匹配的路由，转发数据包到 flannel.1
  - 之后的数据包流向就与上面的 pod1 到 pod2 的网络一样了

pod 到外网

- pod 向 qq.com 发送请求
- 查找路由表,转发数据包到宿主的网卡
- 宿主网卡完成 qq.com 路由选择后，iptables 执行 MASQUERADE，把源 IP 更改为宿主网卡的 IP
- 向 qq.com 服务器发送请求

## Host-GW 型后端

![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/lmmp31/1616118626150-5a2dc6c9-df90-4487-b957-d4b0bfb60f50.png)

假设现在，Node 1 上的 Infra-container-1，要访问 Node 2 上的 Infra-container-2。当你设置 Flannel 使用 host-gw 模式之后，flanneld 会在宿主机上创建这样一条规则，以

Node 1 为例：

    $ ip route
    ...
    10.244.1.0/24 via 10.168.0.3 dev eth0

这条路由规则的含义是：目的 IP 地址属于 10.244.1.0/24 网段的 IP 包，应该经过本机的 eth0 设备发出去（即：dev eth0）；并且，它下一跳地址（next-hop）是 10.168.0.3（即：via 10.168.0.3）。

所谓下一跳地址就是：如果 IP 包从主机 A 发到主机 B，需要经过路由设备 X 的中转。那么 X 的 IP 地址就应该配置为主机 A 的下一跳地址。

而从 host-gw 示意图中我们可以看到，这个下一跳地址对应的，正是我们的目的宿主机 Node2。一旦配置了下一跳地址，那么接下来，当 IP 包从网络层进入链路层封装成帧的时候，eth0 设备就会使用下一跳地址对应的 MAC 地址，作为该数据帧的目的 MAC 地址。显然，这个 MAC 地址，正是 Node 2 的 MAC 地址。

这样，这个数据帧就会从 Node 1 通过宿主机的二层网络顺利到达 Node 2 上。

而 Node 2 的内核网络栈从二层数据帧里拿到 IP 包后，会“看到”这个 IP 包的目的 IP 地址是 10.244.1.3，即 Infra-container-2 的 IP 地址。这时候，根据 Node 2 上的路由表，该目的地址会匹配到第二条路由规则（也就是 10.244.1.0 对应的路由规则），从而进入 cni0 网桥，进而进入到 Infra-container-2 当中。

### host-gw 模式的工作原理

其实就是将每个 Flannel 子网（Flannel Subnet，比如：10.244.1.0/24）的“下一跳”，设置成了该子网对应的宿主机的 IP 地址。也就是说，这台“主机”（Host）会充当这条容器通信路径里的“网关”（Gateway）。这也正是“host-gw”的含义。

Flannel 子网和主机的信息，都是保存在 etcd 当中的，flanneld 只需要监控这些数据的变化，然后实时更新路由表即可

host-gw 型与 vxlan 型的区别：

- host-gw 模型没有 flannel.1 设备来对数据包进行封装，直接添加的路由，所以性能更好(因为少了封装解封装的步骤)。
- Note：但是这也导致 Host-GW 型后端有个前提，Host 必须在二层互通，否则数据包经过 Host 的 eth 网卡后，无法路由。

Note:在了解过 Calico 的工作方式之后，其实会有这么一个疑问，calico 的 bgp 与 flannel 的 host-gw 原理其实一样，但是为什么 host-gw 要多了一个 cni 网桥呢? veth 设备出来的数据包明明可以不再经过 cni 往前直接进入宿主机网络栈发出。原因是：Flannel host-gw 模式使用 CNI 网桥的主要原因，其实是为了跟 VXLAN 模式保持一致。否则的话，Flannel 就需要维护两套 CNI 插件了。

# Flannel 关联文件与配置

**/etc/kube-flannel/net-conf.json** # 在 flannel 的容器内进入该目录，可以看到 configmap 中的信息

- **kubectl get configmaps -n kube-system kube-flannel-cfg -o yaml** # 通过容器部署的 flannel 的配置文件,其中包括 cni-conf.json 和 net-conf.json 两个配置文件。可以看到 flannel 的基本信息，包括后端类型

**/run/flannel/subnet.env** # 读取 net-conf.json 配置文件，加载 /run/flannel/subnet.env 环境变量信息。基于加载的信息，生成适用于其 delegate 的 ipam 和 CNI bridge 的 netconf 文件；其中指定 ipam 使用 host-local，CNI bridge type 为 bridge。调用 deletgate（CNI bridge type）对应的二进制文件来挂载容器到网桥上。

- Note:该文件会由将 Flannel 分配的子网信息都记录下来，并交给每个节点的 cni0 或者 flannel0 使用，如果想要修改 flannel 配置，则需要删除每个节点上的这个文件，该文件内容如下所示

```
[root@master-1 CNI]# cat /run/flannel/subnet.env
FLANNEL_NETWORK=10.252.0.0/16
FLANNEL_SUBNET=10.252.0.1/24
FLANNEL_MTU=1500
FLANNEL_IPMASQ=true
```

# 配置详解

> 参考：
> https://github.com/flannel-io/flannel/blob/master/Documentation/configuration.md

net-conf.json 配置文件：

- **Network** # flannel 使用的 CIDR(无类域间路由无类域间路由)格式的网络地址，用于为 Pod 配置网络功能，比如 flannel 使用默认的 10.244.0.0/16 网段，然后需要给每个节点分配一个网段，即使用了 SubnetLen 的配置。
- **SubnetLen** # 把 Network 切分子网供各个节点使用时，使用多长的掩码进行切分。`默认值：24`。比如给 node1 分配 10.244.1.0/24,给 Node2 分配 10.244.2.0/24 等
- **SubnetMin** # 指定切分子网时候的起始 IP
- **SubnetMax** # 指定切分子网时候的结束 IP
- **Backend** # vxlan，host-gw，udp(Backend 下面还有一个字段 Type，在该字段填写那 3 种类型中的一种的名称)

配置示例：

```json
{
  "Network": "10.0.0.0/8",
  "SubnetLen": 20,
  "SubnetMin": "10.10.0.0",
  "SubnetMax": "10.99.0.0",
  "Backend": {
    "Type": "udp",
    "Port": 7890
  }
}
```

## Flannel 命令行参数

- -etcd-cafile string # SSL Certificate Authority file used to secure etcd communication
- -etcd-certfile string # SSL certification file used to secure etcd communication
- -etcd-endpoints string # a comma-delimited list of etcd endpoints (default "http://127.0.0.1:4001,http://127.0.0.1:2379")
- -etcd-keyfile string # SSL key file used to secure etcd communication
- -etcd-password string # password for BasicAuth to etcd
- -etcd-prefix string # etcd prefix (default "/coreos.com/network")
- -etcd-username string # username for BasicAuth to etcd
- -healthz-ip string # the IP address for healthz server to listen (default "0.0.0.0")
- -healthz-port int # the port for healthz server to listen(0 to disable)
- **-iface \<STRING>** # 用于主机间通信的网络设备名称或者 IP。可以指定多个网络设备，Flannel 会按顺序检查，并使用找到的第一个网络设备
  - 注意：这个参数指定的网络设备，就是 Flannel 建立静态路由条目时所使用的网络设备。
- **-iface-regex \<EXP>** # 用于主机间通信的网络设备的正则表达式
  - 可以多次指定以按顺序检查每个正则表达式。返回找到的第一个匹配项。在检查 iface 选项指定的特定接口后，将检查正则表达式。
  - 比如 `^(eth0|bond1)$` 这种格式，可以让具有不通网络设备名称的设备被统一
- -ip-masq # setup IP masquerade rule for traffic destined outside of overlay network
- -iptables-forward-rules # add default accept rules to FORWARD chain in iptables (default true)
- -iptables-resync int # resync period for iptables rules, in seconds (default 5)
- -kube-annotation-prefix string # Kubernetes annotation prefix. Can contain single slash "/", otherwise it will be appended at the end. (default "flannel.alpha.coreos.com")
- -kube-api-url string # Kubernetes API server URL. Does not need to be specified if flannel is running in a pod.
- -kube-subnet-mgr # contact the Kubernetes API for subnet assignment instead of etcd.
- -kubeconfig-file string # kubeconfig file location. Does not need to be specified if flannel is running in a pod.
- -log_backtrace_at value # when logging hits line file:N, emit a stack trace
- -net-config-path string # path to the network configuration file (default "/etc/kube-flannel/net-conf.json")
- -public-ip string # IP accessible by other nodes for inter-host communication
- -subnet-file string # filename where env variables (subnet, MTU, ... ) will be written to (default "/run/flannel/subnet.env")
- -subnet-lease-renew-margin int # subnet lease renewal margin, in minutes, ranging from 1 to 1439 (default 60)
- -v value # log level for V logs
- -vmodule value # comma-separated list of pattern=N settings for file-filtered logging

# 清理 flannel

```bash
ifconfig cni0 down
ip link delete cni0
ifconfig flannel.1 down
ip link delete flannel.1
rm -rf /var/lib/cni/
rm -f /etc/cni/net.d/*
```

# Flannel 问题总结

## 不支持 IPv6

<https://github.com/flannel-io/flannel/issues/248>

## 误删 cni0 网络设备后恢复

> 参考：
> - 原文：[公众号-k8s 中文社区，一起误删 cni0 虚拟网卡引发的 k8s 事故](https://mp.weixin.qq.com/s/TDdatl6Mzfc_4VdTSXDv4A)

误操作的命令：`ip link del cni0`

由于 flannel 使用的是 vxlan 模式，所以创建 cni0 网桥的时候需要注意 mtu 值的设置。如下，创建 cni0 网桥：

```bash
# 创建cni0设备，指定类型为网桥
# ip link add cni0 type bridge
# ip link set dev cni0 up
// 为cni0设置ip地址，这个地址是pod的网关地址，需要和flannel.1对应网段
# ifconfig cni0 172.28.0.1/25
// 为cni0设置mtu为1450
# ifconfig cni0 mtu 1450 up

// 查看创建情况
# ifconfig cni0
cni0: flags=4163<UP,BROADCAST,RUNNING,MULTICAST>  mtu 1450
        inet 172.28.0.1  netmask 255.255.255.128  broadcast 172.28.0.127
        ether 0e:5e:b9:62:0d:60  txqueuelen 1000  (Ethernet)
        RX packets 487334  bytes 149990594 (149.9 MB)
        RX errors 0  dropped 0  overruns 0  frame 0
        TX packets 629306  bytes 925100055 (925.1 MB)
        TX errors 0  dropped 0 overruns 0  carrier 0  collisions 0
// 此时查看路由表，也已经有了去往本机pod网段的cni0信息
# route -n | grep cni0
172.28.0.0      0.0.0.0         255.255.255.128 U     0      0        0 cni0
```

挂载网络设备

```bash
for veth in $(ip addr | grep veth | grep -v master | awk -F'[@|:]' '{print $2}' | sed 's/ //g')
do
 ip link set dev $veth master cni0
done
```
