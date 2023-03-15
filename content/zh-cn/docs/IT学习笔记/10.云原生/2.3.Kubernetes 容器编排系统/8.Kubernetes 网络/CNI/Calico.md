---
title: Calico
---

# Calico 基本概念

# ![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/bgb09a/1616118512894-684263a6-9d81-4a8a-8be3-088fd6aedef1.png)

基于以 BGP 协议构建网络，主要由三个部分组成

第一部分：Calico 的 CNI 插件。这是 Calico 与 Kubernetes 对接的部分

第二部分：Felix，一个 DaemonSet。负责在 Host 上插入路由规则(即：写入 Linux 内核的 FIB(转发信息库 Forwarding information base)，以及维护 Calico 所需的网络设备等工作

第三部分：BIRD，BGP Client。专门负责在集群内分发路由规则信息

Calico 利用 Linux 内核原生的路由和 iptables 防火墙功能。进出各个容器，虚拟机和主机的所有流量都会在路由到目标之前遍历这些内核规则。

1. calicoctl：允许您从简单的命令行界面实现高级策略和网络。
2. orchestrator plugins：提供与各种流行协调器的紧密集成和同步。
3. key/value store：保存 Calico 的策略和网络配置状态。比如 etcd
4. calico/node：在每个主机上运行，从 key/value store 中读取相关的策略和网络配置信息，并在 Linux 内核中实现它。
5. Dikastes/Envoy：可选的 Kubernetes sidecar，通过相互 TLS 身份验证保护工作负载到工作负载的通信，并实施应用层策略。

## Calico BGP 工作原理

实际上，Calico 项目提供的网络解决方案，与 Flannel 的 host-gw 模式，几乎是完全一样的。也就是说，Calico 也会在每台宿主机上，添加一个格式如下所示的路由规则：

    < 目的容器 IP 地址段 > via < 网关的 IP 地址 > dev eth0

其中，网关的 IP 地址，正是目的容器所在宿主机的 IP 地址。

而正如前所述，这个三层网络方案得以正常工作的核心，是为每个容器的 IP 地址，找到它所对应的、“下一跳”的网关。不过，不同于 Flannel 通过 Etcd 和宿主机上的 flanneld 来维护路由信息的做法，Calico 项目使用了一个“重型武器”来自动地在整个集群中分发路由信息。这个“重型武器”，就是 BGP。详见：BGP 协议

Calico 项目的架构由三个部分组成：

1. Calico 的 CNI 插件。这是 Calico 与 Kubernetes 对接的部分。
2. Felix。它是一个 DaemonSet，负责在宿主机上插入路由规则（即：写入 Linux 内核的 FIB 转发信息库），以及维护 Calico 所需的网络设备等工作。
3. BIRD。它就是 BGP 的客户端，专门负责在集群里分发路由规则信息。

除了对路由信息的维护方式之外，Calico 项目与 Flannel 的 host-gw 模式的另一个不同之处，就是它不会在宿主机上创建任何网桥设备。这时候，Calico 的工作方式，可以用一幅示意图来描述，如下所示（在接下来的讲述中，我会统一用“BGP 示意图”来指代它）
![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/bgb09a/1617285084820-0455ddc1-9808-4eb4-955f-2c37a3697372.png)

其中的绿色实线标出的路径，就是一个 IP 包从 Node 1 上的 Container 1，到达 Node 2 上的 Container 4 的完整路径。可以看到，Calico 的 CNI 插件会为每个容器设置一个 Veth Pair 设备，然后把其中的一端放置

在宿主机上（它的名字以 cali 前缀开头）。此外，由于 Calico 没有使用 CNI 的网桥模式，Calico 的 CNI 插件还需要在宿主机上为每个容器的 Veth Pair 设备配置一条路由规则，用于接收传入的 IP 包。比如，宿主机 Node 2 上的 Container 4 对应的路由规则，如下所示

    # 发往 10.233.2.3 的 IP 包，应该进入 cali5863f3 设备
    10.233.2.3 dev cali5863f3 scope link

Note：基于上述原因，Calico 项目在宿主机上设置的路由规则，肯定要比 Flannel 项目多得多(因为没有单独的网桥来连接这些 veth 设备，也就没法通过网桥来自动将数据包交给对应的容器，只能为每一个 veth 单独配置路由信息)。不过，Flannel host-gw 模式使用 CNI 网桥的主要原因，其实是为了跟 VXLAN 模式保持一致。否则的话，Flannel 就需要维护两套 CNI 插件了。

有了这样的 Veth Pair 设备之后，容器发出的 IP 包就会经过 Veth Pair 设备出现在宿主机上。

然后，宿主机网络栈就会根据路由规则的下一跳 IP 地址，把它们转发给正确的网关。接下来的流程就跟 Flannel host-gw 模式完全一致了。

其中，这里最核心的“下一跳”路由规则，就是由 Calico 的 Felix 进程负责维护的。这些路由规则信息，则是通过 BGP Client 也就是 BIRD 组件，使用 BGP 协议传输而来的。

而这些通过 BGP 协议传输的消息，你可以简单地理解为如下格式：

    [BGP 消息]
    我是宿主机 192.168.1.2
    10.233.2.0/24 网段的容器都在我这里
    这些容器的下一跳地址是我

不难发现，Calico 项目实际上将集群里的所有节点，都当作是边界路由器来处理，它们一起组成了一个全连通的网络，互相之间通过 BGP 协议交换路由规则。这些节点，我们称为 BGP Peer。

### Route Reflector 的出现

需要注意的是，Calico 维护的网络在默认配置下，是一个被称为“Node-to-Node Mesh”的模式。这时候，每台宿主机上的 BGP Client 都需要跟其他所有节点的 BGP Client 进行通信以便交换路由信息。但是，随着节点数量 N 的增加，这些连接的数量就会以 N² 的规模快速增长，从而给集群本身的网络带来巨大的压力。

所以，Node-to-Node Mesh 模式一般推荐用在少于 100 个节点的集群里。而在更大规模的集群中，你需要用到的是一个叫作 Route Reflector 的模式。

在这种模式下，Calico 会指定一个或者几个专门的节点，来负责跟所有节点建立 BGP 连接从而学习到全局的路由规则。而其他节点，只需要跟这几个专门的节点交换路由信息，就可以获得整个集群的路由规则信息了。

这些专门的节点，就是所谓的 Route Reflector 节点，它们实际上扮演了“中间代理”的角色，从而把 BGP 连接的规模控制在 N 的数量级上。

### k8s 节点跨网段的解决方案

此外，我在前面提到过，Flannel host-gw 模式最主要的限制，就是要求集群宿主机之间是二层连通的。而这个限制对于 Calico 来说，也同样存在。

举个例子，假如我们有两台处于不同子网的宿主机 Node 1 和 Node 2，对应的 IP 地址分别是 192.168.1.2 和 192.168.2.2。需要注意的是，这两台机器通过路由器实现了三层转发，所以这两个 IP 地址之间是可以相互通信的。而我们现在的需求，还是 Container 1 要访问 Container 4。按照我们前面的讲述，Calico 会尝试在 Node 1 上添加如下所示的一条路由规则：

    10.233.2.0/16 via 192.168.2.2 eth0

但是，这时候问题就来了。

上面这条规则里的下一跳地址是 192.168.2.2，可是它对应的 Node 2 跟 Node 1 却根本不在一个子网里，没办法通过二层网络把 IP 包发送到下一跳地址。

在这种情况下，你就需要为 Calico 打开 IPIP 模式。我把这个模式下容器通信的原理，总结成了一副示意图，如下所示（接下来我会称之为：IPIP 示意图）：
![image.jpeg](https://notes-learning.oss-cn-beijing.aliyuncs.com/bgb09a/1617285101992-492a004e-7a36-4cd2-8d66-a2ec3772827b.jpeg)

在 Calico 的 IPIP 模式下，Felix 进程在 Node 1 上添加的路由规则，会稍微不同，如下所示：

    10.233.2.0/24 via 192.168.2.2 tunl0

可以看到，尽管这条规则的下一跳地址仍然是 Node 2 的 IP 地址，但这一次，要负责将 IP 包发出去的设备，变成了 tunl0。注意，是 T-U-N-L-0，而不是 Flannel UDP 模式使用的 T-UN-0（tun0），这两种设备的功能是完全不一样的。

Calico 使用的这个 tunl0 设备，是一个 IP 隧道（IP tunnel）设备。在上面的例子中，

1. IP 包进入 IP 隧道设备之后，就会被 Linux 内核的 IPIP 驱动接管。
2. IPIP 驱动会将这个 IP 包直接封装在一个宿主机网络的 IP 包中，类似 flannel 中 vxlan 的封装过程。其中，经过封装后的新的 IP 包的目的地址正是原 IP 包的下一跳地址，即 Node 2 的 IP 地址：192.168.2.2。而原 IP 包本身，则会被直接封装成新 IP 包的 Payload。
3. 这样，原先从容器到 Node 2 的 IP 包，就被伪装成了一个从 Node 1 到 Node 2 的 IP 包。
4. 由于宿主机之间已经使用路由器配置了三层转发，也就是设置了宿主机之间的“下一跳”。所以这个 IP 包在离开 Node 1 之后，就可以经过路由器，最终“跳”到 Node 2 上。
5. 这时，Node 2 的网络内核栈会使用 IPIP 驱动进行解包，从而拿到原始的 IP 包。然后，原始 IP 包就会经过路由规则和 Veth Pair 设备到达目的容器内部。

### 如何让外部真正的数通设备(路由器、三层交换机)也加入到集群的 BGP 中

以上，就是 Calico 项目主要的工作原理了。

不难看到，当 Calico 使用 IPIP 模式的时候，集群的网络性能会因为额外的封包和解包工作而下降。在实际测试中，Calico IPIP 模式与 Flannel VXLAN 模式的性能大致相当。所以，在实际使用时，如非硬性需求，我建议你将所有宿主机节点放在一个子网里，避免使用 IPIP。

不过，通过上面对 Calico 工作原理的讲述，你应该能发现这样一个事实：如果 Calico 项目能够让宿主机之间的路由设备（也就是网关），也通过 BGP 协议“学习”到 Calico 网络里的路由规则，那么从容器发出的 IP 包，不就可以通过这些设备路由到目的宿主机了么？

比如，只要在上面“IPIP 示意图”中的 Node 1 上，添加如下所示的一条路由规则：

    10.233.2.0/24 via 192.168.1.1 eth0

然后，在 Router 1 上（192.168.1.1），添加如下所示的一条路由规则：

    10.233.2.0/24 via 192.168.2.1 eth0

那么 Container 1 发出的 IP 包，就可以通过两次“下一跳”，到达 Router 2（192.168.2.1）了。以此类推，我们可以继续在 Router 2 上添加“下一条”路由，最终把 IP 包转发到 Node 2

上。

遗憾的是，上述流程虽然简单明了，但是在 Kubernetes 被广泛使用的公有云场景里，却完全不可行。

这里的原因在于：公有云环境下，宿主机之间的网关，肯定不会允许用户进行干预和设置。

不过，在私有部署的环境下，宿主机属于不同子网（VLAN）反而是更加常见的部署状态。这时候，想办法将宿主机网关也加入到 BGP Mesh 里从而避免使用 IPIP，就成了一个非常迫切的需求。

而在 Calico 项目中，它已经为你提供了两种将宿主机网关设置成 BGP Peer 的解决方案。

1. 第一种方案，就是所有宿主机都跟宿主机网关建立 BGP Peer 关系。
   1. 这种方案下，Node 1 和 Node 2 就需要主动跟宿主机网关 Router 1 和 Router 2 建立 BGP 连接。从而将类似于 10.233.2.0/24 这样的路由信息同步到网关上去。需要注意的是，这种方式下，Calico 要求宿主机网关必须支持一种叫作 Dynamic Neighbors 的 BGP 配置方式。这是因为，在常规的路由器 BGP 配置里，运维人员必须明确给出所有 BGP Peer 的 IP 地址。考虑到 Kubernetes 集群可能会有成百上千个宿主机，而且还会动态地添加和删除节点，这时候再手动管理路由器的 BGP 配置就非常麻烦了。而 Dynamic Neighbors 则允许你给路由器配置一个网段，然后路由器就会自动跟该网段里的主机建立起 BGP Peer 关系。
   2. 不过，相比之下，我更愿意推荐第二种方案。
2. 第二种方案，是使用一个或多个独立组件负责搜集整个集群里的所有路由信息，然后通过 BGP 协议同步给网关。
   1. 而我们前面提到，在大规模集群中，Calico 本身就推荐使用 Route Reflector 节点的方式进行组网。所以，这里负责跟宿主机网关进行沟通的独立组件，直接由 Route Reflector 兼任即可。
3. 更重要的是，这种情况下网关的 BGP Peer 个数是有限并且固定的。所以我们就可以直接把这些独立组件配置成路由器的 BGP Peer，而无需 Dynamic Neighbors 的支持。
4. 当然，这些独立组件的工作原理也很简单：它们只需要 WATCH Etcd 里的宿主机和对应网段的变化信息，然后把这些信息通过 BGP 协议分发给网关即可。

# Calico 配置

# Calicoctl 命令行工具是使用说明

calicoctl 命令行工具用于管理 Calico 网络和安全策略，查看和管理后端配置以及管理 Calico 节点实例

安装方式：

1. kubectl apply -f \ <https://docs.projectcalico.org/v3.4/getting-started/kubernetes/installation/hosted/calicoctl.yaml>
2. alias calicoctl="kubectl exec -i -n kube-system calicoctl /calicoctl -- "
3. 然后可以直接使用别名来使用该命令，calico 对应的 etcd 后端 ip 以及 etcd 证书已经在 calicoctl 的 pod 中定义完成

配置文件路径：/etc/calico/calicoctl.cfg #该路径一般是对于使用二进制方式把 calicoctl 命令文件放在 linux 的$PATH 中来使用

calicoctl \[OPTIONS] <COMMAND> \[\<ARGS>...]

OPTIONS

1. -h --help # Show this screen.
2. -l --log-level=\<level> # Set the log level (one of panic, fatal, error,warn, info, debug) \[default: panic]

COMMAND

create Create a resource by filename or stdin.

calicoctl create --filename=\<FILENAME> \[--skip-exists] \[--config=\<CONFIG>] \[--namespace=\<NS>]

replace Replace a resource by filename or stdin.

apply Apply a resource by filename or stdin. This creates a resource if it does not exist, and replaces a resource if it does exists.

delete Delete a resource identified by file, stdin or resource type and name.

get Get a resource identified by file, stdin or resource type and name.

1. EXAMPLE

convert Convert config files between different API versions.

ipam IP address management.

node Calico node management.

version Display the version of calicoctl.
