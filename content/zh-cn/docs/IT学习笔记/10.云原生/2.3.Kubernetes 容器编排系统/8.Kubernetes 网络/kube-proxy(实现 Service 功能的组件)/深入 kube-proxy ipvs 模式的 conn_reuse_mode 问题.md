---
title: 深入 kube-proxy ipvs 模式的 conn_reuse_mode 问题
---

在高并发、短连接的场景下，kube-proxy ipvs 存在 rs 删除失败或是延迟高的问题，社区也有不少 Issue 反馈，比如**kube-proxy ipvs conn_reuse_mode setting causes errors with high load from single client**。文本对这些问题进行了梳理，试图介绍产生这些问题的内部原因。由于能力有限，其中涉及内核部分，只能浅尝辄止。

## 背景

### 端口重用

一切问题来源于端口重用。在 TCP 四次挥手中有个`TIME_WAIT`的状态，作为先发送`FIN`包的一端，在接收到对端发送的`FIN`包后进入`TIME_WAIT`，在经过`2MSL`后才会真正关闭连接。`TIME_WAIT`状态的存在，一来可以避免将之前连接的延迟报文，作为当前连接的报文处理；二是可以处理最后一个 ACK 丢失带来的问题。
![](https://notes-learning.oss-cn-beijing.aliyuncs.com/mo2k6r/1620787619057-62e254bb-bf63-4c3f-a453-14d99d185c20.png)
而在短连接、高并发的场景下，会出现大量的`TIME-WAIT`连接，导致资源无法及时释放。Linux 中内核参数`net.ipv4.tcp_tw_reuse`提供了一种减少`TIME-WAIT`连接的方式，可以将`TIME-WAIT`连接的端口分配给新的 TCP 连接，来复用端口。

    tcp_tw_reuse - BOOLEAN
     Allow to reuse TIME-WAIT sockets for new connections when it is
     safe from protocol viewpoint. Default value is 0.
     It should not be changed without advice/request of technical
     experts.

### ipvs 如何处理端口重用？

ipvs 对端口的复用策略主要由内核参数`net.ipv4.vs.conn_reuse_mode`决定

    conn_reuse_mode - INTEGER
     1 - default
     Controls how ipvs will deal with connections that are detected
     port reuse. It is a bitmap, with the values being:
     0: disable any special handling on port reuse. The new
     connection will be delivered to the same real server that was
     servicing the previous connection. This will effectively
     disable expire_nodest_conn.
     bit 1: enable rescheduling of new connections when it is safe.
     That is, whenever expire_nodest_conn and for TCP sockets, when
     the connection is in TIME_WAIT state (which is only possible if
     you use NAT mode).
     bit 2: it is bit 1 plus, for TCP connections, when connections
     are in FIN_WAIT state, as this is the last state seen by load
     balancer in Direct Routing mode. This bit helps on adding new
     real servers to a very busy cluster.

当`net.ipv4.vs.conn_reuse_mode=0`时，ipvs 不会对新连接进行重新负载，而是复用之前的负载结果，将新连接转发到原来的 rs 上；当`net.ipv4.vs.conn_reuse_mode=1`时，ipvs 则会对新连接进行重新调度。
相关的，还有一个内核参数`net.ipv4.vs.expire_nodest_conn`，用于控制连接的 rs 不可用时的处理。在开启时，如果后端 rs 不可用，会立即结束掉该连接，使客户端重新发起新的连接请求；否则将数据包**silently drop**，也就是 DROP 掉数据包但不结束连接，等待客户端的重试。
另外，关于**destination 不可用**的判断，是在 ipvs 执行删除`vs`（在`__ip_vs_del_service()`中实现）或删除`rs`（在`ip_vs_del_dest()`中实现）时，会调用`__ip_vs_unlink_dest()`方法，将相应的 destination 置为不可用。

    expire_nodest_conn - BOOLEAN
            0 - disabled (default)
            not 0 - enabled
            The default value is 0, the load balancer will silently drop
            packets when its destination server is not available. It may
            be useful, when user-space monitoring program deletes the
            destination server (because of server overload or wrong
            detection) and add back the server later, and the connections
            to the server can continue.
            If this feature is enabled, the load balancer will expire the
            connection immediately when a packet arrives and its
            destination server is not available, then the client program
            will be notified that the connection is closed. This is
            equivalent to the feature some people requires to flush
            connections when its destination is not available.

关于 ipvs 如何处理端口复用的连接，这块主要实现逻辑在`net/netfilter/ipvs/ip_vs_core.c`的`ip_vs_in()`方法中：

    /*
     * Check if the packet belongs to an existing connection entry
     */
    cp = pp->conn_in_get(ipvs, af, skb, &iph);  //找是属于某个已有的connection
    conn_reuse_mode = sysctl_conn_reuse_mode(ipvs);
    //当conn_reuse_mode开启，同时出现端口复用（例如收到TCP的SYN包，并且也属于已有的connection），进行处理
    if (conn_reuse_mode && !iph.fragoffs && is_new_conn(skb, &iph) && cp) {
     bool uses_ct = false, resched = false;
     //如果开启了expire_nodest_conn、目标rs的weight为0
     if (unlikely(sysctl_expire_nodest_conn(ipvs)) && cp->dest &&
         unlikely(!atomic_read(&cp->dest->weight))) {
      resched = true;
      //查询是否用到了conntrack
      uses_ct = ip_vs_conn_uses_conntrack(cp, skb);
     } else if (is_new_conn_expected(cp, conn_reuse_mode)) {
     //连接是expected的情况，比如FTP
      uses_ct = ip_vs_conn_uses_conntrack(cp, skb);
      if (!atomic_read(&cp->n_control)) {
       resched = true;
      } else {
       /* Do not reschedule controlling connection
        * that uses conntrack while it is still
        * referenced by controlled connection(s).
        */
       resched = !uses_ct;
      }
     }
     //如果expire_nodest_conn未开启，并且也非期望连接，实际上直接跳出了
     if (resched) {
      if (!atomic_read(&cp->n_control))
       ip_vs_conn_expire_now(cp);
      __ip_vs_conn_put(cp);
      //当开启了net.ipv4.vs.conntrack，SYN数据包会直接丢弃，等待客户端重新发送SYN
      if (uses_ct)
       return NF_DROP;
      //未开启conntrack时，会进入下面ip_vs_try_to_schedule的流程
      cp = NULL;
     }
    }
    if (unlikely(!cp)) {
     int v;
     if (!ip_vs_try_to_schedule(ipvs, af, skb, pd, &v, &cp, &iph))
      return v;
    }
    IP_VS_DBG_PKT(11, af, pp, skb, iph.off, "Incoming packet");
    /* Check the server status */
    if (cp->dest && !(cp->dest->flags & IP_VS_DEST_F_AVAILABLE)) {
     /* the destination server is not available */
     __u32 flags = cp->flags;
     /* when timer already started, silently drop the packet.*/
     if (timer_pending(&cp->timer))
      __ip_vs_conn_put(cp);
     else
      ip_vs_conn_put(cp);
     if (sysctl_expire_nodest_conn(ipvs) &&
         !(flags & IP_VS_CONN_F_ONE_PACKET)) {
      /* try to expire the connection immediately */
      ip_vs_conn_expire_now(cp);
     }
     return NF_DROP;
    }

### kube-proxy ipvs 模式下的优雅删除

Kubernetes 提供了 Pod 优雅删除机制。当我们决定干掉一个 Pod 时，我们可以通过`PreStop Hook`来做一些服务下线前的处理，同时 Kubernetes 也有个`grace period`，超过这个时间但未完成删除的 Pod 会被强制删除。
而在 Kubernetes 1.13 之前，kube-proxy ipvs 模式并不支持优雅删除，当 Endpoint 被删除时，kube-proxy 会直接移除掉 ipvs 中对应的 rs，这样会导致后续的数据包被丢掉。
在 1.13 版本后，Kubernetes 添加了**IPVS 优雅删除**的逻辑，主要是两点：

- 当 Pod 被删除时，kube-proxy 会先将 rs 的`weight`置为 0，以防止新连接的请求发送到此 rs，由于不再直接删除 rs，旧连接仍能与 rs 正常通信；
- 当 rs 的`ActiveConn`数量为 0（后面版本已改为`ActiveConn+InactiveConn==0`)，即不再有连接转发到此 rs 时，此 rs 才会真正被移除。

## kube-proxy ipvs 模式下的问题

看上去 kube-proxy ipvs 的删除是优雅了，但当优雅删除正巧碰到端口重用，那问题就来了。
首先，kube-proxy 希望通过设置`weight`为 0，来避免新连接转发到此 rs。但当`net.ipv4.vs.conn_reuse_mode=0`时，对于端口复用的连接，ipvs 不会主动进行新的调度（调用`ip_vs_try_to_schedule`方法）；同时，只是将`weight`置为 0，也并不会触发由`expire_nodest_conn`控制的结束连接或 DROP 操作，就这样，新连接的数据包当做什么都没发生一样，发送给了正在删除的 Pod。这样一来，只要不断的有端口复用的连接请求发来，rs 就不会被 kube-proxy 删除，上面提到的优雅删除的两点均无法实现。
而当`net.ipv4.vs.conn_reuse_mode=1`时，根据`ip_vs_in()`的处理逻辑，当开启了`net.ipv4.vs.conntrack`时，会 DROP 掉第一个 SYN 包，导致 SYN 的重传，有 1S 延迟。而 Kube-proxy 在 IPVS 模式下，使用了 iptables 进行`MASQUERADE`，也正好开启了`net.ipv4.vs.conntrack`。

    conntrack - BOOLEAN
     0 - disabled (default)
     not 0 - enabled
     If set, maintain connection tracking entries for
     connections handled by IPVS.
     This should be enabled if connections handled by IPVS are to be
     also handled by stateful firewall rules. That is, iptables rules
     that make use of connection tracking.  It is a performance
     optimisation to disable this setting otherwise.
     Connections handled by the IPVS FTP application module
     will have connection tracking entries regardless of this setting.
     Only available when IPVS is compiled with CONFIG_IP_VS_NFCT enabled.

这样看来，目前的情况似乎是，如果你需要实现优雅删除中的“保持旧连接不变，调度新连接”能力，那就要付出 1s 的延迟代价；如果你要好的性能，那么就不能重新调度。

## 如何解决

从 Kubernetes 角度来说，Kube-proxy 需要在保证性能的前提下，找到一种能让新连接重新调度的方式。但目前从内核代码中可以看到，需要将参数设置如下

    net.ipv4.vs.conntrack=0
    net.ipv4.vs.conn_reuse_mode=1
    net.ipv4.vs.expire_nodest_conn=1

但 Kube-proxy ipvs 模式目前无法摆脱 iptables 来完成 k8s service 的转发。此外，Kube-proxy 只有在`ActiveConn+InactiveConn==0`时才会删除 rs，除此之外，在新的 Endpoint 和`GracefulTerminationList`（保存了`weight`为 0，但暂未删除的 rs）中的 rs 冲突时，才会立即删除 rs。这种逻辑似乎并不合理。目前 Pod 已有优雅删除的逻辑，而 kube-proxy 应基于 Pod 的优雅删除，在网络层面做好 rs 的优雅删除，因此在 kubelet 完全删除 Pod 后，Kube-proxy 是否也应该考虑同时删除相应的 rs？
另外，从内核角度来说，ipvs 需要提供一种方式，能在端口复用、同时使用 conntrack 的场景下，可以对新连接直接重新调度。

## 即将到来

这个问题在社区讨论一段时间后，目前出现的几个相关的解决如下：

### 内核两个 Patch

- **ipvs: allow connection reuse for unconfirmed conntrack**修改了`ip_vs_conn_uses_conntrack()`方法的逻辑，当使用`unconfirmed conntrack`时，返回 false，这种修改针对了 TIME_WAIT 的 conntrack。
- **ipvs: queue delayed work to expire no destination connections if expire_nodest_conn=1**提前了`expire connection`的操作，在 destination 被删除后，便开始将`expire connection`操作入队列。而不是等到数据包真正发过来时，才做`expire connection`，以此来减少数据包的丢失。

### Kubernetes

- **Graceful Termination for External Traffic Policy Local**
- **Add Terminating Condition to EndpointSlice**

正如前面所说的，Kube-proxy 需要能够感知到 Pod 的优雅删除过程，来同步进行 rs 的删除。目前，已有一个相应的 KEP 在进行中，通过在`Endpoint.EndpointConditions`中添加`terminating`字段，来为 kube-proxy 提供感知方式。

### 脚注

\[1]
kube-proxy ipvs conn*reuse_mode setting causes errors with high load from single client: *<https://github.com/kubernetes/kubernetes/issues/81775>_
\[2]
IPVS 优雅删除: _<https://github.com/kubernetes/kubernetes/pull/66012>_
\[3]
ipvs: allow connection reuse for unconfirmed conntrack: _<http://patchwork.ozlabs.org/project/netfilter-devel/patch/20200701151719.4751-1-ja@ssi.bg/>_
\[4]
ipvs: queue delayed work to expire no destination connections if expire_nodest_conn=1: _<http://patchwork.ozlabs.org/project/netfilter-devel/patch/20200708161638.13584-1-kim.andrewsy@gmail.com/>_
\[5]
Graceful Termination for External Traffic Policy Local: _<https://github.com/kubernetes/enhancements/pull/1607>_
\[6]
Add Terminating Condition to EndpointSlice: _<https://github.com/kubernetes/kubernetes/pull/92968>\_

原文链接：[https://maao.cloud/2021/01/15/%E6%B7%B1%E5%85%A5kube-proxy%20ipvs%E6%A8%A1%E5%BC%8F%E7%9A%84conn_reuse_mode%E9%97%AE%E9%A2%98/](https://maao.cloud/2021/01/15/%25E6%25B7%25B1%25E5%2585%25A5kube-proxy%2520ipvs%25E6%25A8%25A1%25E5%25BC%258F%25E7%259A%2584conn_reuse_mode%25E9%2597%25AE%25E9%25A2%2598/)
