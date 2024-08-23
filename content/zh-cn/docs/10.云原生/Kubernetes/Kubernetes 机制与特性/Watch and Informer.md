---
title: Watch  and Informer
linkTitle: Watch  and Informer
date: 2024-08-23T09:05
weight: 20
---

# 概述

> 参考：
>
> - [官方文档, 参考 - API 概述 - Kubernetes API 概念, 高效监测变化](https://kubernetes.io/docs/reference/using-api/api-concepts/#efficient-detection-of-changes)
> - [K8s list-watch 机制和 Informer 模块](https://xuliangtang.github.io/posts/k8s-list-watch/#list-watch-%E6%9C%BA%E5%88%B6)

**Reflector**

在 Kubernetes 中，有5个主要的组件，分别是 master 节点上的 kube-api-server、kube-controller-manager 和 kube-scheduler，node 节点上的 kubelet 和kube-proxy 。这其中 kube-apiserver 是对外和对内提供资源的声明式 API 的组件，其它4个组件都需要和它交互。为了保证消息的实时性，有两种方式：

*   客户端组件 (kubelet, scheduler, controller-manager 等) 轮询 apiserver
*   apiserver 通知客户端

为了降低 kube-apiserver 的压力，有一个非常关键的机制就是 list-watch。list-watch 本质上也是 client 端监听 k8s 资源变化并作出相应处理的生产者消费者框架

list-watach 机制需要满足以下需求：

1.  实时性 (即数据变化时，相关组件越快感知越好)
2.  保证消息的顺序性 (即消息要按发生先后顺序送达目的组件。很难想象在Pod创建消息前收到该Pod删除消息时组件应该怎么处理)
3.  保证消息不丢失或者有可靠的重新获取机制 (比如 kubelet 和 kube-apiserver 间网络闪断，需要保证网络恢复后kubelet可以收到网络闪断期间产生的消息)
## list-watch 机制

list-watch 由两部分组成，分别是 list 和 watch。list 非常好理解，就是调用资源的 list API 罗列资源 ，基于 HTTP 短链接实现，watch 则是调用资源的 watch API 监听资源变更事件，基于 HTTP 长链接实现

etcd 存储集群的数据信息，apiserver 作为统一入口，任何对数据的操作都必须经过 apiserver。客户端通过 list-watch 监听 apiserver 中资源的 create, update 和 delete 事件，并针对事件类型调用相应的事件处理函数

## informer 机制

k8s 的 informer 模块封装 list-watch API，用户只需要指定资源，编写事件处理函数，AddFunc, UpdateFunc 和 DeleteFunc 等。如下图所示，informer 首先通过 list API 罗列资源，然后调用 watch API 监听资源的变更事件，并将结果放入到一个 FIFO 队列，队列的另一头有协程从中取出事件，并调用对应的注册函数处理事件。Informer 还维护了一个只读的 Map Store 缓存，主要为了提升查询的效率，降低 apiserver 的负载

![](https://raw.githubusercontent.com/xuliangTang/picbeds/main/picgo/202212181902721.png)

![](https://raw.githubusercontent.com/xuliangTang/picbeds/main/picgo/202212181854816.jpeg)

Reflector 从 API Server 中通过 List&Watch 得到资源的状态变化，把数据塞到 Delta Fifo 队列里（Reflector 相当于生产者），由 Informer 进行消费。更新时在回调里可以获得新值和旧值，旧值从 Indexer（store） 中获取

*   FIFO ：先入先出队列，拥有队列基本方法（ADD，UPDATE，DELETE，LIST，POP，CLOSE 等）
*   Delta ： 存储对象的行为（变化）类型（Added，Updated，Deleted，Sync 等）

如果要对一个资源支持多种监听方式，需要使用到 **SharedInformer**（SharedIndexInformer）

*   支持多个EventHandler . 可以认为是支持多个消费者，多个消费者之间共享 Indexer， Reflector 统一下发数据统一处理
*   内置一个 Indexer（有一个叫做 threadSafeMap 的 struct 来实现 cache/thread\_safe\_store.go）

里面有个属性 sharedProcessor，用于协调和管理若干个处理器对象 processorListener（这是真正干活的对象）

*   run()：阻塞运行
*   pop()：好比不断从队列里取数据，完成对应的回调操作
*   addCh：一个 channel，外部向它插入数据

好比 Reflector 向 deltaFifo 插入数据后，后分发给 ProcessListener，由 ProcessListener 执行具体的回调。Processor 里面可以放若干个 Listener，因此可以使用多个回调

如果要对多个资源支持多种监听方式，需要使用到 **SharedInformerFactory**，里面有个属性 informers 包含多个 SharedIndexInformer 对象

client-go 使用 k8s.io/client-go/tools/cache 包里的 informer 对象进行 list-watch 机制的封装

最粗暴的解释：

1.  初始化时，调 List API 获得全量 list，缓存起来(本地缓存)，这样就不需要每次请求都去请求 ApiServer
2.  调用 Watch API 去 watch 资源，发生变更后会通过一定机制维护缓存
