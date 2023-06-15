---
title: Connection Tracking(连接跟踪)机制
weight: 2
---

# 概述

> 参考：
> 
> - [Netfilter 官方文档，连接跟踪工具用户手册](https://conntrack-tools.netfilter.org/manual.html)
> - [云计算基层技术-netfilter 框架研究](https://opengers.github.io/openstack/openstack-base-netfilter-framework-overview/)
> - [arthurchiao.art 的文章](http://arthurchiao.art/index.html)
>   - [连接跟踪（conntrack）：原理、应用及 Linux 内核实现](http://arthurchiao.art/blog/conntrack-design-and-implementation-zh/)

**Connection Tracking(连接跟踪系统，简称 ConnTrack、CT)**，用于跟踪并且记录连接状态。Linux 为每一个经过网络堆栈的数据包，生成一个 **ConnTrack Entry(连接跟踪条目，简称 Entry)**，并把该条目记录在一个 **ConnnTrack Table(连接跟踪表)** 中，条目中主要是包含该连接的协议、源 IP 和 PORT、目标 IP 和 PORT、协议号、数据包的大小等等等信息。此后，在处理数据包时读取该文件，在文件中所有属于此连接的数据包都被唯一地分配给这个连接，并标识连接的状态。该文件中的每一个条目都有一个持续时间，当持续时间结束后，该连接会被自动清除，再有相同的连接进来的时候，则按照新连接来处理。Netfilter 中定义了如下几个连接状态以便对具有这些状态的连接进行处理：

可跟踪的连接状态有以下几个

1. NEW：新发出的请求。在连接跟踪文件中(nf_conntrack)不存在此连接。
2. ESTABLISHED：已建立的。NEW 状态之后，在 nf_conntrack 文件中为其建立的条目失效之前所进行的通信的状态
3. RELATED：有关联的。某个已经建立的连接所建立的新连接；e.g.FTP 的数据传输连接就是控制连接所 RELATED 出来的连接。–icmp-type 8(ping 请求)就是–icmp-type 0(ping 应答) 所 RELATED 出来的。
4. INVALIED：无法识别的连接。
5. UNTRACKED：不跟踪的链接状态，仅在使用 raw 表的时候该状态才有用，即 raw 不进行链接跟踪的时候，则连接跟踪表中没有记录的数据包就是此状态
6. 其他：
   1. NEW 与 ESTABLISHED 的定义：只要第一次请求就算 NEW(e.g.本机往外第一次发送，外部第一次发往本机的请求)，哪怕对第一个 NEW 请求再回应的都算 ESTABLISHED。注意在 INPUT 和 OUTPUT 链上定义 NEW 的情况，INPUT 是对外部访问本机来说第一次是 NEW；OUTPUT 是对本机访问外部来说第一次是 NEW。

注意：ConnTrack 中所定义的状态与 TCP 等协议所定义的状态不一样，这里面定义的状态只是为了可以通过一种新的方式来处理每一个数据包，并进行过滤，这是 Netfilter 中所定义的状态

ConnTrack 功能依靠 **nf_conntrack** 模块来实现的，当启用 iptables 功能时(比如 firewalld)会自动加载该模块

连接跟踪是防火墙模块的状态检测的基础，同时也是地址转换中实现 SNAT 和 DNAT 的基础，如果在 nat 表上没有连接跟踪，那么则没法进行 nat 转换(比如通过 raw 表来关闭连接跟踪)。

# ConnTrack Table(连接跟踪表)

ConnTrack 将连接跟踪表保存于系统的内存当中，可以通过 **cat /proc/net/nf_conntrack** 或 **conntrack -L** 命令查看到当前已跟踪的所有 **ConnTrack Entry(连接跟踪条目)**。不同的协议，条目的内容也不太一样，下面是一个 tcp 协议的条目内容：

- `ipv4     2 tcp      6 299 ESTABLISHED src=192.168.2.40 dst=172.38.40.250 sport=61758 dport=22 src=172.38.40.250 dst=192.168.2.40 sport=22 dport=61758 [ASSURED] mark=0 zone=0 use=2`

nf_conntrack 文件中，每个条目占用单独一行。条目中包含了数据包的原始方向信息(蓝色部分)，和期望的响应包信息(红色部分)，这样内核能够在后续到来的数据包中识别出属于此连接的双向数据包，并更新此连接的状态。

在内核中，**ConnTrackTable(连接跟踪表)** 实际上是一个 **hash table(哈希表)**。收到一个数据包，通过如下步骤判断该数据包是否署一个已有连接(即定位连接跟踪条目)：

- 第一步：内核提取数据包信息(源 IP、目的 IP、port，协议号)进行 hash 计算得到一个 hash 值，在哈希表中以此 hash 值做索引，索引结果为数据包所属的 **Bucket(储存区)**。这一步 hash 计算时间固定并且很短。
  - 一个 **Bucket(储存区)** 里包含一个 **linked list(已链接的列表，简称链表)**，即已经追踪的条目的列表。也就是说，每个 Bucket 里可以存放多个 ConnTrack Entry。所谓 Bucket 的大小，就是指一个 Bucket 中可以存放多少个 ConnTrack Entry。
- 第二步：遍历第一步获取的 Bucket 中的所有条目，查找是否有匹配的条目。这一步是比较耗时的操作，所以说 Bucket 越大，遍历时间越长

## Bucket(储存区)

在 Connection Tracking 系统中的 hash table 中，有若干个 **Bucket(储存区)**，Bucket 的个数通过两个内核参数计算而来

- net.netfilter.nf_conntrack_buckets # 一个表最大的 Bucket 数量。默认通过内存计算得来，内存越高，Bucket 越多。也可以通过设置模块参数指定具体的数值
  - 无法通过 sysctl 修改 nf_conntrack_buckets 的值，该值只能通过加载 nf_conntrack 模块时的参数来决定。使用 `echo "options nf_conntrack hashsize=16384" > /etc/modprobe.d/nf_conntrack.conf` 命令即可设置该内核参数为 16384
- net.netfilter.nf_conntrack_max # 一个表最大的 Entry 数量。默认为 nf_conntrack_buckets 值的 4 倍。也就是说，**Bucket 的大小默认为 4**，即系统默认每个 Bucket 中包含 4 个 ConnTrack Entry。
  - 当不使用系统默认的 nf_conntrack_buckets 值时，则 nf_conntrack_max 的值为 nf_conntrack_buckets 的 8 倍

如果把一个 Bucket 的大小称为 `BucketSize` 的话，那么`BucketSize = nf_conntrack_max / nf_conntrack_buckets`(这意思就是说 `储存区的大小=条目总数 / 储存区的总数`，所以储存区大小就是指能装多少条目)

## ConnTrack Entry(连接跟踪条目)

conntrack 从经过它的数据包中提取详细的，唯一的信息，因此能保持对每一个连接的跟踪。关于 conntrack 如何确定一个连接，对于 tcp/udp，连接由他们的源目地址，源目端口唯一确定。对于 icmp，由 type，code 和 id 字段确定。

```
ipv4     2 tcp      6 33 SYN_SENT src=172.16.200.119 dst=172.16.202.12 sport=54786 dport=10051 [UNREPLIED] src=172.16.202.12 dst=172.16.200.119 sport=10051 dport=54786 mark=0 zone=0 use=2
```

如上是一条 conntrack 条目，它代表当前已跟踪到的某个连接，conntrack 维护的所有信息都包含在这个条目中，通过它就可以知道某个连接处于什么状态

- **ipv4** # 此连接使用 ipv4 协议，是一条 tcp 连接(tcp 的协议类型代码是 6)
- **33** # 这条 conntrack 条目在当前时间点的生存时间(每个 conntrack 条目都会有生存时间，从设置值开始倒计时，倒计时完后此条目将被清除)，可以使用`sysctl -a |grep conntrack | grep timeout`查看不同协议不同状态下生存时间设置值，当然这些设置值都可以调整，注意若后续有收到属于此连接的数据包，则此生存时间将被重置(重新从设置值开始倒计时)，并且状态改变，生存时间设置值也会响应改为新状态的值
- **SYN_SENT** # 到此刻为止 conntrack 跟踪到的这个连接的状态(内核角度)，`SYN_SENT`表示这个连接只在一个方向发送了一初始 TCP SYN 包，还未看到响应的 SYN+ACK 包(只有 tcp 才会有这个字段)。
- **src=172.16.200.119 dst=172.16.202.12 sport=54786 dport=10051** # 从数据包中提取的此连接的源目地址、源目端口，是 conntrack 首次看到此数据包时候的信息。
- **\[UNREPLIED]** # 说明此刻为止这个连接还没有收到任何响应，当一个连接已收到响应时，\[UNREPLIED]标志就会被移除
- **src=172.16.202.12 dst=172.16.200.119 sport=10051 dport=54786** # 地址和端口和前面是相反的，这部分不是数据包中带有的信息，是 conntrack 填充的信息，代表 conntrack 希望收到的响应包信息。意思是若后续 conntrack 跟踪到某个数据包信息与此部分匹配，则此数据包就是此连接的响应数据包。注意这部分确定了 conntrack 如何判断响应包(tcp/udp)，icmp 是依据另外几个字段

上面是 tcp 连接的条目，而 udp 和 icmp 没有连接建立和关闭过程，因此条目字段会有所不同，后面 iptables 状态匹配部分我们会看到处于各个状态的 conntrack 条目

注意：conntrack 机制并不能够修改或过滤数据包，它只是跟踪网络连接并维护连接跟踪表，以提供给 iptables 做状态匹配使用，也就是说，如果你 iptables 中用不到状态匹配，那就没必要启用 conntrack

## 总结

所以，一个 ConnTrack Table 就类似于下面的表：

| **Hash Table** | Bucket 1 | Bucket 2 | Bucket 3 | .......  | Bucket N |
| -------------- | -------- | -------- | -------- | -------- | -------- |
| Entry 1        | 条目内容 | 条目内容 | 条目内容 | 条目内容 | 条目内容 |
| Entry 2        | 条目内容 | 条目内容 | 条目内容 | 条目内容 | 条目内容 |
| ........       | 条目内容 | 条目内容 | 条目内容 | 条目内容 | 条目内容 |
| Entry N        | 条目内容 | 条目内容 | 条目内容 | 条目内容 | 条目内容 |

# 计算连接跟踪表所占内存

`total * mem_used(单位为 Bytes) = nf_conntrack_max * sizeof(struct ip*conntrack) + nf_conntrack_buckets * sizeof(struct list_head)`

1. sizeof(struct ip_conntrack) 连接跟踪对象大小，默认 376
2. sizeof(struct list_head) 链表项大小，默认为 16

上述两个值可以通过如下 python 代码计算出来

```bash
import ctypes

#不同系统可能此库名不一样，需要修改
LIBNETFILTER_CONNTRACK = 'libnetfilter_conntrack.so.3.6.0'

nfct = ctypes.CDLL(LIBNETFILTER_CONNTRACK)
print 'sizeof(struct nf_conntrack):', nfct.nfct_maxsize()
print 'sizeof(struct list_head):', ctypes.sizeof(ctypes.c_void_p) * 2
```

假如，我系统信息如下：

```bash
~]# cat /proc/sys/net/netfilter/nf_conntrack_max
524288
~]# cat /proc/sys/net/netfilter/nf_conntrack_buckets
131072
```

那么，此系统下，连接跟踪表所占内存即为：

```
(524288 * 376 + 131072 * 16) / 1024 / 1024 = 190MiB
```

# ConnTrack 关联文件与配置

> 参考：
> 
> - [内核官方文档，网络-nf_conntrack-sysctl](https://www.kernel.org/doc/Documentation/networking/nf_conntrack-sysctl.txt)

连接跟踪系统的配置大部分都可以通过修改内核参数来进行，还有一部分需要通过指定 模块的参数 来配置。

- **/proc/net/nf_conntrack** # 连接跟踪表，该文件用于记录每一个连接跟踪条目
  - 注意：Ubuntu 中没有该文件，可以通过 `conntrack -L` 命令获取连接跟踪条目。据说该文件已 deprecated(弃用)，但是未找到官方说明
  - <https://forum.ubuntu.com.cn/viewtopic.php?t=480072>
  - <https://askubuntu.com/questions/266991/in-ubuntu-12-10-how-to-enable-proc-net-ip-conntrack>
  - <https://patchwork.ozlabs.org/project/ubuntu-kernel/patch/1341986947-28300-3-git-send-email-bryan.wu@canonical.com/>
  - <https://github.com/kubernetes/kubernetes/pull/69589/files#r418929810>
- **/proc/sys/net/nf_conntrack_max** # 等于 /proc/sys/net/netfilter/nf_conntrack_max 的值。修改这俩参数任意一个值，都会互相同步。
- **/proc/sys/net/netfilter/** # 网络栈的运行时属性所在的目录
  - **./nf_conntrack_count** # 当前连接跟踪数。
  - **./nf_conntrack_max** # 连接跟踪表的大小，即一个表中有可以存放多少个条目。默认值为 nf_conntrack_buckets \*4 。等于 /proc/sys/net/nf_conntrack_max 的值。
  - **./nf_conntrack_buckets** # hash 表的大小，即一个 hash 表中有多少个 Buckets。
  - **./nf_conntrack_tcp_timeout_time_wait** # timewait 状态的条目超时时间。 默认 120 秒

# 应用实例

第一次 TCP 或 UDP 或 ICMP 等协议请求建立连接后，有一个持续时间，在持续时间内，这个连接信息会保存在连接跟踪表(记录在 nf_conntrack 文件中)中，当同一个 IP 再次请求的时候，这个请求的数据包则不会被当成 NEW 状态的数据包来处理(具体的状态有几种详见下文)，这个概念可以用在这么一个真实环境当中。

# conntrack 命令行工具

> 参考：

# 常见问题

## nf_conntrack: table full, dropping packets

参考: https://mp.weixin.qq.com/s/N7jfQCfR-1V5ppw7tLcB6w

显然，调大最大值的限制就可以了。不过更大的限制意味着可以承接更多连接，意味着要耗费更多资源，这点要注意。

查看当前有多少活跃连接：

```
cat /proc/sys/net/netfilter/nf_conntrack_count
```

如果这个值跟上面介绍的 nf_conntrack_max 已经很接近了，就说明快满了，需要调大 nf_conntrack_max。可以使用下面的命令临时调大：

```
echo 524288 > /proc/sys/net/netfilter/nf_conntrack_max
```

如果不想每次重启都要重新设置，可以修改 /etc/sysctl.conf，加入下面的配置：

```
net.netfilter.nf_conntrack_max = 524288
```

为了缓解大量连接的问题，您可能还需要考虑减少服务器等待连接关闭/超时的时间。在 /etc/sysctl.conf 中加入下面的配置：

```
net.netfilter.nf_conntrack_tcp_timeout_close_wait = 60  
net.netfilter.nf_conntrack_tcp_timeout_fin_wait = 60  
net.netfilter.nf_conntrack_tcp_timeout_time_wait = 60
```
