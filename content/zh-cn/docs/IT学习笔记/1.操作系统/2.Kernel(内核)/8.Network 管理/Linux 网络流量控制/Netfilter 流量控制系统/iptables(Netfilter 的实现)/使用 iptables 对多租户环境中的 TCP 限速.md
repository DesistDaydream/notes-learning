---
title: 使用 iptables 对多租户环境中的 TCP 限速
---

[使用 iptables 对多租户环境中的 TCP 限速](https://mp.weixin.qq.com/s/n7bRJb-u5bzIj4TMb8JE-A)

我们有个服务以类似 SideCar 的方式和应用一起运行，SideCar 和应用通过 Unix Domain Socket 进行通讯。为了方便用户，在开发的时候不必在自己的开发环境中跑一个 SideCar，我用 socat 在一台开发环境的机器上 map UDS 到一个端口。这样用户在开发的时候就可以直接通过这个 TCP 端口测试服务，而不用自己开一个 SideCar 使用 UDS 了。

因为所有人都要用这一个地址做开发，所以就有互相影响的问题。虽然性能还可以，几十万 QPS 不成问题，但是总有憨憨拿来搞压测，把资源跑满，影响别人。我在使用说明文档里用红色大字写了这是开发测试用的，不能压测，还是有一些视力不好的同事会强行压测。隔三差五我就得去解释一番，礼貌地请同事不要再这样做了。

最近实在累了。研究了一下直接给这个端口加上 per IP 的 rate limit，效果还不错。方法是在 Per-IP rate limiting with iptables\[1] 学习到的，这个公司是提供一个多租户的 SaaS 服务，也有类似的问题：有一些非正常用户 abuse 他们的服务，由于 abuse 发生在连接建立阶段，还没有进入到业务代码，所以无法从应用的层面进行限速，解决发现就是通过 iptables 实现的。详细的实现方法可以参考这篇文章。

iptables 本身是无状态的，每一个进入的 packet 都单独判断规则。rate limit 显然是一个有状态的规则，所以要用到 module: `hashlimit`。（原文中还用到了 `conntrack`，他是想只针对新建连接做限制，已经建立的连接不限制速度了。因为这个应用内部就可以控制了，但是我这里是想对所有的 packet 进行限速，所以就不需要用到这个 module）

完整的命令如下：

```bash
$ iptables --new-chain SOCAT-RATE-LIMIT
$ iptables --append SOCAT-RATE-LIMIT \
    --match hashlimit \
    --hashlimit-mode srcip \
    --hashlimit-upto 50/sec \
    --hashlimit-burst 100 \
    --hashlimit-name conn_rate_limit \
    --jump ACCEPT
$ iptables --append SOCAT-RATE-LIMIT --jump DROP
$ iptables -I INPUT -p tcp --dport 1234 --jump SOCAT-RATE-LIMIT
```

第一行是新建一个 iptables Chain，做 rate limit；

第二行处理如果在 rate limit 限额内，就接受包；否则跳到第三行，直接将包 DROP；

最后将新的 Chain 加入到 INPUT 中，对此端口的流量进行限制。

有关 rate limit 的算法，主要是两个参数：

1. `--hashlimit-upto` 其实本质上是 1s 内可以进入多少 packet，`50/sec` 就是 `20ms` 一个 packet；
2. 那如何在 `10ms` 发来 10 个 packet，后面一直没发送，怎么办？这个在测试情景下也比较常见，不能要求用户一直匀速地发送。所以就要用到 `--hashlimit-burst`。字面意思是瞬间可以发送多少 packet，但实际上，可以理解这个参数就是可用的 credit。

两个指标配合起来理解，就是每个 ip 刚开始都会有 `burst` 个 credit，每个 ip 发送来的 packet 都会占用 `burst` 里面的 credit，用完了之后再发来的包就会被直接 DROP。这个 credit 会以 `upto` 的速度一直增加，但是最多增加到 `burst`（初始值），之后就 _use it or lost it_.

举个例子，假如 `--hashlimit-upto 50/sec --hashlimit-burst 20` 的话，某个 IP 以匀速每 ms 一个 packet 的速度发送，最终会有多少 packets 被接受？答案是 70. 最初的 20ms，所有的 packet 都会被接受，因为 `--hashlimit-burst` 是 20，所以最初的 credit 是 20. 这个用完之后就要依赖 `--hashlimit--upto 50/sec` 来每 20ms 获得一个 packet credit 了。所以每 20ms 可以接受一个。

这是限速之后的效果，非常明显：

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/3af7c350-867e-429b-be82-d0a4816715c0/640)

> 原文链接：<https://www.kawabangga.com/posts/4594>

### 参考资料

\[1]

Per-IP rate limiting with iptables: [_https://making.pusher.com/per-ip-rate-limiting-with-iptables/index.html_](https://making.pusher.com/per-ip-rate-limiting-with-iptables/index.html)
