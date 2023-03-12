---
title: k8s CPU limit和throttling的迷思
---

原文链接：

- <https://nanmu.me/zh-cn/posts/2021/myth-of-k8s-cpu-limit-and-throttle/>
- <https://mp.weixin.qq.com/s/QYJycJCaxB42xdEo3qHHHA>

你应当小心设定 k8s 中负载的 CPU limit，太小的值会给你的程序带来额外的、无意义的延迟，太大的值会带来过大的爆炸半径，削弱集群的整体稳定性。

## request 和 limit

k8s 的一大好处就是资源隔离，通过设定负载的 request 和 limit，我们可以方便地让不同程序共存于合适的节点上。

其中，request 是给调度看的，调度会确保节点上所有负载的 CPU request 合计与内存 request 合计分别都不大于节点本身能够提供的 CPU 和内存，limit 是给节点（kubelet）看的，节点会保证负载在节点上只使用这么多 CPU 和内存。例如，下面配置意味着单个负载会调度到一个剩余 CPU request 大于 0.1 核，剩余 request 内存大于 200MB 的节点，并且负载运行时的 CPU 使用率不能高于 0.4 核（超过将被限流），内存使用不多余 300MB（超过将被 OOM Kill 并重启）。

`resources:   requests:     memory: 200Mi     cpu: "0.1"   limits:     memory: 300Mi     cpu: "0.4"`

## CPU 的利用率

CPU 和内存不一样，它是量子化的，只有“使用中”和“空闲”两个状态。


我和老婆聊了聊 CPU 和内存的不同，她帮我画了一张插图 图/我的妻子

当我们说内存的使用率是 60%时，我们是在说内存有 60%在**空间上**已被使用，还有 40%的空间可以放入负载。但是，当我们说 CPU 的某个核的使用率是 60%时，我们是在说采样时间段内，CPU 的这个核在**时间上**有 60%的时间在忙，40%的时间在睡大觉。

你设定负载的 CPU limit 时，这个时空区别可能会带来一个让你意想不到的效果——过分的降速限流， 节点 CPU 明明不忙，但是节点故意不让你的负载全速使用 CPU，服务延时上升。

## CPU 限流

k8s 使用 CFS（Completely Fair Scheduler，完全公平调度）限制负载的 CPU 使用率，CFS 本身的机制比较复杂\[1]，但是 k8s 的文档中给了一个简明的解释\[2]，要点如下：

- CPU 使用量的计量周期为 100ms；
- CPU limit 决定每计量周期（100ms）内容器可以使用的 CPU 时间的上限；
- 本周期内若容器的 CPU 时间用量达到上限，CPU 限流开始，容器只能在下个周期继续执行；
- 1 CPU = 100ms CPU 时间每计量周期，以此类推，0.2 CPU = 20ms CPU 时间每计量周期，2.5 CPU = 250ms CPU 时间每计量周期；
- 如果程序用了多个核，CPU 时间会累加统计。

举个例子，假设一个 API 服务在响应请求时需要使用 A, B 两个线程（2 个核），分别使用 60ms 和 80ms，其中 B 线程晚触发 20ms，我们看到 API 服务在 100ms 后可给出响应：

没有 CPU 限制的情况，响应时间为 100ms

如果 CPU limit 被设为 1 核，即每 100ms 内最多使用 100ms CPU 时间，API 服务的线程 B 会受到一次限流（灰色部分），服务在 140ms 后响应：


CPU limit = 1，响应时间为 140ms

如果 CPU limit 被设为 0.6 核，即每 100ms 内最多使用 60ms CPU 时间，API 服务的线程 A 会受到一次限流（灰色部分），线程 B 受到两次限流，服务在 220ms 后响应：


CPU limit = 0.6，响应时间为 220ms

注意，**即使此时 CPU 没有其他的工作要做，限流一样会执行**，这是个死板不通融的机制。

这是一个比较夸张的例子，一般的 API 服务是 IO 密集型的，CPU 时间使用量没那么大（你在跑模型推理？当我没说），但还是可以看到，限流会实打实地延伸 API 服务的延时。因此，对于延时敏感的服务，我们都应该尽量避免触发 k8s 的限流机制。

下面这张图是我工作中一个 API 服务在 pod 级别的 CPU 使用率和 CPU 限流比率（CPU Throttling），我们看到，CPU 限流的情况在一天内的大部分时候都存在，限流比例在 10%上下浮动，这意味着服务的工作没能全速完成，在速度上打了 9 折。值得一提，这时 pod 所在节点仍然有富余的 CPU 资源，节点的整体 CPU 使用率没有超过 50%.


一个实际的降速限流的例子，服务的处理速度被 kubelet 降低了 10%

你可能注意到，监控图表里的 CPU 使用率看上去没有达到 CPU limit（橙色横线），这是由于 CPU 使用率的统计周期（1min）太长造成的信号混叠（Aliasing）\[3]，如果它的统计统计周期和 CFS 的一样（100ms），我们就能看到高过 CPU limit 的尖刺了。（这不是 bug，这是 feature）

不过，内核版本低于 4.18 的 Linux 还真有个 bug 会造成不必要的 CPU 限流\[4]。┑(￣ Д ￣)┍

## 避免 CPU 限流

有的开发者倾向于完全弃用 CPU limit\[5]，裸奔直接跑，特别是内核版本不够有 bug 的时候\[6]。

我认为这么做还是太过放飞自我了，如果程序里有耗尽 CPU 的 bug（例如死循环，我不幸地遇到过），整个节点及其负载都会陷入不可用的状态，爆炸半径太大，特别是在大号的节点上（16 核及以上）。

我有两个建议：

1. 监控一段时间应用的 CPU 利用率，基于利用率设定一个合适的 CPU limit（例如，日常利用率的 95 分位 \* 10），同时该 limit 不要占到节点 CPU 核数的太大比例（例如 2/3），这样可以达到性能和安全的一个平衡。
2. 使用 automaxprocs\[7]一类的工具让程序适配 CFS 调度环境，各个语言应该都有类似的库或者执行参数，根据 CFS 的特点调整后，程序更不容易遇到 CPU 限流\[8]。

## 结语

上面说到的信号混叠（采样频率不足）和 Linux 内核 bug 让我困扰了一年多，现在想想，主要还是望文生义惹的祸，文档还是应该好好读，基础概念还是要搞清，遂记此文章于错而知新\[9]。

题外话，性能和资源利用率有时是相互矛盾的。对于延时不敏感的程序，CPU 限流率控制在 10%以内应该都是比较健康可接受的，量体裁衣，在线离线负载混合部署，可以提升硬件的资源利用率。有消息说腾讯云研发投产了基于服务优先级的抢占式调度\[10]，这是一条更难但更有效的路，希望有朝一日在上游能看到他们的相关贡献。

### 引用链接

\[1]

CFS 本身的机制比较复杂: [_https://en.wikipedia.org/wiki/Completely_Fair_Scheduler_](https://en.wikipedia.org/wiki/Completely_Fair_Scheduler)

\[2]

简明的解释: [_https://kubernetes.io/docs/concepts/configuration/manage-resources-containers/#how-pods-with-resource-limits-are-run_](https://kubernetes.io/docs/concepts/configuration/manage-resources-containers/#how-pods-with-resource-limits-are-run)

\[3]

信号混叠（Aliasing）: [_https://en.wikipedia.org/wiki/Aliasing_](https://en.wikipedia.org/wiki/Aliasing)

\[4]

内核版本低于 4.18 的 Linux 还真有个 bug 会造成不必要的 CPU 限流: [_https://github.com/kubernetes/kubernetes/issues/67577#issuecomment-466609030_](https://github.com/kubernetes/kubernetes/issues/67577#issuecomment-466609030)

\[5]

完全弃用 CPU limit: [_https://amixr.io/blog/what-wed-do-to-save-from-the-well-known-k8s-incident/_](https://amixr.io/blog/what-wed-do-to-save-from-the-well-known-k8s-incident/)

\[6]

内核版本不够有 bug 的时候: [_https://medium.com/omio-engineering/cpu-limits-and-aggressive-throttling-in-kubernetes-c5b20bd8a718_](https://medium.com/omio-engineering/cpu-limits-and-aggressive-throttling-in-kubernetes-c5b20bd8a718)
