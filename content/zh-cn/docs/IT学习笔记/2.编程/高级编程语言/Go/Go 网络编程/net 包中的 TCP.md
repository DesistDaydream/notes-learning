---
title: net 包中的 TCP
---

# 概述

> 参考：
> - [知乎,TCP 漫谈之 keepalive 和 time_wait](https://zhuanlan.zhihu.com/p/126688315)

TCP 是一个有状态通讯协议，所谓的有状态是指通信过程中通信的双方各自维护连接的状态。
**一、TCP keepalive**

先简单回顾一下 TCP 连接建立和断开的整个过程。（这里主要考虑主流程，关于丢包、拥塞、窗口、失败重试等情况后面详细讨论。）
首先是客户端发送 syn（Synchronize Sequence Numbers：同步序列编号）包给服务端，告诉服务端我要连接你，syn 包里面主要携带了客户端的 seq 序列号；服务端回发一个 syn+ack，其中 syn 包和客户端原理类似，只不过携带的是服务端的 seq 序列号，ack 包则是确认客户端允许连接；最后客户端再次发送一个 ack 确认接收到服务端的 syn 包。这样客户端和服务端就可以建立连接了。整个流程称为“三次握手”。
![](https://notes-learning.oss-cn-beijing.aliyuncs.com/hmo2md/1626270211217-6cd6f47b-f74f-4f7a-841c-f95aaec1bd50.jpeg)
建立连接后，客户端或者服务端便可以通过已建立的 socket 连接发送数据，对端接收数据后，便可以通过 ack 确认已经收到数据。数据交换完毕后，通常是客户端便可以发送 FIN 包，告诉另一端我要断开了；另一端先通过 ack 确认收到 FIN 包，然后发送 FIN 包告诉客户端我也关闭了；最后客户端回应 ack 确认连接终止。整个流程成为“四次挥手”。TCP 的性能经常为大家所诟病，除了 TCP+IP 额外的 header 以外，它建立连接需要三次握手，关闭连接需要四次挥手。如果只是发送很少的数据，那么传输的有效数据是非常少的。是不是建立一次连接后续可以继续复用呢？的确可以这样做，但这又带来另一个问题，如果连接一直不释放，端口被占满了咋办。为此引入了今天讨论的第一个话题 TCP keepalive。所谓的 TCP keepalive 是指 TCP 连接建立后会通过 keepalive 的方式一直保持，不会在数据传输完成后立刻中断，而是通过 keepalive 机制检测连接状态。Linux 控制 keepalive 有三个参数：保活时间 net.ipv4.tcp_keepalive_time、保活时间间隔 net.ipv4.tcp_keepalive_intvl、保活探测次数 net.ipv4.tcp_keepalive_probes，默认值分别是 7200 秒（2 小时）、75 秒和 9 次探测。如果使用 TCP 自身的 keepalive 机制，在 Linux 系统中，最少需要经过 2 小时 + 9\*75 秒后断开。譬如我们 SSH 登录一台服务器后可以看到这个 TCP 的 keepalive 时间是 2 个小时，并且会在 2 个小时后发送探测包，确认对端是否处于连接状态。
![](https://notes-learning.oss-cn-beijing.aliyuncs.com/hmo2md/1626270211249-79d9a329-86e0-4158-afb3-ea2a75f868e3.png)
之所以会讨论 TCP 的 keepalive，是因为发现服器上有泄露的 TCP 连接：

```go
# ll /proc/11516/fd/10
lrwx------ 1 root root 64 Jan  3 19:04 /proc/11516/fd/10 -> socket:[1241854730]
# date
Sun Jan  5 17:39:51 CST 2020
```

已经建立连接两天，但是对方已经断开了（非正常断开）。由于使用了比较老的 go（1.9 之前版本有问题）导致连接没有释放。解决这类问题，可以借助 TCP 的 keepalive 机制。新版 go 语言支持在建立连接的时候设置 keepalive 时间。

以 go 1.16 版本为例。首先查看 net 包中建立 TCP 连接的 [DialContext()](https://github.com/golang/go/blob/release-branch.go1.16/src/net/dial.go#L369) 方法，其中 [defaultTCPKeepAlive](https://github.com/golang/go/blob/release-branch.go1.16/src/net/dial.go#L17) 是 15s：

```go
if tc, ok := c.(*TCPConn); ok && d.KeepAlive >= 0 {
   setKeepAlive(tc.fd, true)
   ka := d.KeepAlive
   if d.KeepAlive == 0 {
      ka = defaultTCPKeepAlive
   }
   setKeepAlivePeriod(tc.fd, ka)
   testHookSetKeepAlive(ka)
}
```

如果是 HTTP 连接，使用默认的 `http.Client`，那么它会将 keepalive 时间设置成 30s，

> 代码：
>
> - 默认 http.Client 中的 Transport 的默认值：<https://github.com/golang/go/blob/release-branch.go1.16/src/net/http/transport.go#L42>
> - 在 do() 方法中会调用 send() 函数，send() 函数中调用 transport() 方法来返回 Transport 的值。
>   - do() 方法：<https://github.com/golang/go/blob/release-branch.go1.16/src/net/http/client.go#L590>
>     - 调用 send() 函数：<https://github.com/golang/go/blob/release-branch.go1.16/src/net/http/client.go#L717>
>   - send() 函数：<https://github.com/golang/go/blob/release-branch.go1.16/src/net/http/client.go#L169>
>     - 调用 transport() 方法：<https://github.com/golang/go/blob/release-branch.go1.16/src/net/http/client.go#L175>
>   - transport() 方法：<https://github.com/golang/go/blob/release-branch.go1.16/src/net/http/client.go#L194>

```go
var DefaultTransport RoundTripper = &Transport{
	Proxy: ProxyFromEnvironment,
	DialContext: (&net.Dialer{
		Timeout:   30 * time.Second,
		KeepAlive: 30 * time.Second,
	}).DialContext,
	ForceAttemptHTTP2:     true,
	MaxIdleConns:          100,
	IdleConnTimeout:       90 * time.Second,
	TLSHandshakeTimeout:   10 * time.Second,
	ExpectContinueTimeout: 1 * time.Second,
}
```

下面通过一个简单的 demo 测试一下，代码如下：

```go
package main

import (
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"sync"
	"time"
)

func main() {
	wg := &sync.WaitGroup{}

	dialer := &net.Dialer{
		Timeout:   30 * time.Second,
        // TCP 连接进入 keepalive 状态前的等待时间
		KeepAlive: 30 * time.Second,
		DualStack: true,
	}

	transport := &http.Transport{
		DialContext:           dialer.DialContext,
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}

	client := &http.Client{
		Transport: transport,
	}

	for i := 0; i < 2; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				r, err := client.Get("http://**.**.**.**:****")
				if err != nil {
					fmt.Println(err)
					return
				}
				_, err = ioutil.ReadAll(r.Body)
				r.Body.Close()
				if err != nil {
					fmt.Println(err)
					return
				}
				time.Sleep(30 * time.Millisecond)
			}
		}()
	}
	wg.Wait()
}
```

执行程序后，可以查看连接。初始设置 keepalive 为 30s。
![](https://notes-learning.oss-cn-beijing.aliyuncs.com/hmo2md/1626270211294-fb52b3fb-7b5b-416c-91b8-6b5b2869d3f7.png)
然后不断递减，至 0 后，又会重新获取 30s。
![](https://notes-learning.oss-cn-beijing.aliyuncs.com/hmo2md/1626270211298-e367c78a-af97-4fb6-92a0-42f7417c59e6.png)
整个过程可以通过 tcpdump 抓包获取。

    # tcpdump -i bond0 port 35832 -nvv -A

其实很多应用并非是通过 TCP 的 keepalive 机制探活的，因为默认的两个多小时检查时间对于很多实时系统是完全没法满足的，通常的做法是通过应用层的定时监测，如 PING-PONG 机制（就像打乒乓球，一来一回），应用层每隔一段时间发送心跳包，如 websocket 的 ping-pong。
**二、TCP time_wait**

第二个希望和大家分享的话题是 TCP 的 Time_wait 状态。
![](https://notes-learning.oss-cn-beijing.aliyuncs.com/hmo2md/1626270211222-e48fb524-94c3-465f-9117-36709a4e9cc3.png)
为啥需要 time_wait 状态呢？为啥不直接进入 closed 状态呢？直接进入 closed 状态能更快地释放资源给新的连接使用了，而不是还需要等待 2MSL（Linux 默认）时间。有两个原因：一是为了防止“迷路的数据包”。如下图所示，如果在第一个连接里第三个数据包由于底层网络故障延迟送达。等待新的连接建立后，这个迟到的数据包才到达，那么将会导致接收数据紊乱。
![](https://notes-learning.oss-cn-beijing.aliyuncs.com/hmo2md/1626270211341-c2c59e69-f95c-451d-8116-e01caed05c40.jpg)
第二个原因则更加简单，如果因为最后一个 ack 丢失，那么对方将一直处于 last ack 状态，如果此时重新发起新的连接，对方将返回 RST 包拒绝请求，将会导致无法建立新连接。

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/hmo2md/1626270211264-46a07efe-aa21-4c7b-a743-12c3eb9d563a.jpg)

为此设计了 time_wait 状态。在高并发情况下，如果能将 time_wait 的 TCP 复用，time_wait 复用是指可以将处于 time_wait 状态的连接重复利用起来，从 time_wait 转化为 established，继续复用。Linux 内核通过 net.ipv4.tcp_tw_reuse 参数控制是否开启 time_wait 状态复用。读者可能很好奇，之前不是说 time_wait 设计之初是为了解决上面两个问题的吗？如果直接复用不是反而会导致上面两个问题出现吗？这里先介绍 Linux 默认开启的一个 TCP 时间戳策略 net.ipv4.tcp_timestamps = 1。

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/hmo2md/1626270211233-8114fc51-79a6-4719-a8fd-19fe2b64a7f4.png)

时间戳开启后，针对第一个迷路数据包的问题，由于晚到数据包的时间戳过早会被直接丢弃，不会导致新连接数据包紊乱；针对第二个问题，开启 reuse 后，当对方处于 last-ack 状态时，发送 syn 包会返回 FIN,ACK 包，然后客户端发送 RST 让服务端关闭请求，从而客户端可以再次发送 syn 建立新的连接。最后还需要提醒读者的是，Linux 4.1 内核版本之前除了 tcp_tw_reuse 以外，还有一个参数 tcp_tw_recycle，这个参数就是强制回收 time_wait 状态的连接，它会导致 NAT 环境丢包，所以不建议开启。
作者：陈晓宇陈晓宇著作《云计算那些事儿：从 IaaS 到 PaaS 进阶》
