---
title: Etcd 中存储的数据之探究
---

# 概述

原文: https://jingwei.link/2018/11/25/kubernetes-etcd-data-save-specification.html#default%E5%91%BD%E5%90%8D%E7%A9%BA%E9%97%B4%E4%B8%ADendpoint%E5%AE%9E%E4%BE%8Bkubernetes%E7%9A%84%E5%80%BC

K8s 的架构复杂，涉及到概念非常多，其基础组件包含 ETCD、kube-apiserver、kube-controller-manager、kube-scheduler、kubelet、kube-proxy 等，其运行时环境为 docker 或 Rkt，当然还包含很多插件。在我看来，k8s 是 DevOps 的未来，因此不禁想写一些它的故事。

ETCD 在 k8s 技术栈的地位，就仿佛数据库（Mysql、Postgresql 或 oracle 等）在 Web 应用中的地位，它存储了 k8s 集群中所有的元数据（以 key-value 的方式）。那么很现实的一个问题是：这些元数据是如何组织并保存的呢？本文就对此问题探究一番。

相关环境

- 两台 2 个核 4G 内存（2C4G）的虚拟机，ip 分别为 192.168.205.137 和 192.168.205.139
- k8s 相关控件-1.8.6
- etcd-3.3.10
- docker-18.06.1-ce

# k8s 中 ETCD 数据的增删改查

首先应该明确，K8s 中所有元数据的增删改查都是由 kube-apiserver 来执行的，那么这些数据在 ETCD 中必然有一套存储规范，这样才能保证在集群中部署成千上万的应用时不会出差错。在此基础上可以认为，只要掌握了 k8s 在 ETCD 中存储数据的规范，便可以像 k8s 一样手动来操作 ETCD 了（虽然不建议这么做）。不过更大的价值是能对 k8s 的理解更进一步，便于以后 debug 或者二次开发 k8s 的某些功能时更有底气。

# 初探 ETCD 中的数据

本文对 ETCD 的操作主要使用了其官方客户端工具 etcdctl，这里不对 etcdctl 进行详解了（需要用一整篇博客来介绍它才行），只就用到的一些命令进行阐释。

## 获取 ETCD 中的所有的 key 值

下面的命令可以获取 ETCD 中的所有的 key 值：

```bash
# 获取ETCD中的所有数据
# --prefix 表示获取所有key值头为某个字符串的数据， 由于传入的是""，所以会匹配所有的值
# --keys-only 表示只返回key而不返回value
# 对输出的结果使用grep过滤掉空行
/$ ETCDCTL_API=3 etcdctl get "" --prefix --keys-only |grep -Ev "^$"
# 输出结果如下所示，实际数据会非常整齐
/registry/apiextensions.k8s.io/customresourcedefinitions/globalbgpconfigs.crd.projectcalico.org
/registry/apiextensions.k8s.io/customresourcedefinitions/globalfelixconfigs.crd.projectcalico.org
/registry/apiextensions.k8s.io/customresourcedefinitions/globalnetworkpolicies.crd.projectcalico.org
# ... 略过很多条目
/registry/namespaces/default/registry/namespaces/kube-public
/registry/namespaces/kube-system/registry/pods/kube-system/canal-mljsv
/registry/pods/kube-system/canal-qlvh6
# ... 略过很多条目
/registry/services/endpoints/kube-system/kube-scheduler
/registry/services/specs/default/kubernetes
/registry/services/specs/kube-system/kube-dnscompact_rev_key
# 总共有240条记录
/$ etcdctl get "" --prefix --keys-only |grep -Ev "^$"|wc -l240
```

## ETCD 中 key 值的规律

通过观察可以简单得出下面几个规律：

- k8s 主要把自己的数据注册在/registry/前缀下面（在 ETCD-v3 版本后没有了目录的概念，只能一切皆前缀了）。
- 通过观察 k8s 中 deployment、namespace、pod 等在 ETCD 中的表示，可以知道这部分资源的 key 的格式为/registry/#{k8s 对象}/#{命名空间}/#{具体实例名}。
- 存在一个与众不同的 key 值 compact_rev_key，搜索可以知道这是 apiserver/compact.go 中用来记录无效数据版本使用的；运行 etcdctl get compact_rev_key 可以发现，输出的是一个整形数值。
- 在查看 ETCD 时，k8s 中除了必要的网络插件 canal，未部署其他的应用，此时 ETCD 中只有 240 条数据。

有了上面的规律，可以初步得出一个结论：在研究 k8s 时重点关注/registry/前缀的 key 及其 value 即可。

## ETCD 中的 value 值

通过上面的内容知道，k8s 在 ETCD 中保存数据时 key 的取值非常讲究，规律非常容易概括出来。那么这些 key 值所对应的值是什么样子呢？我试着输出了/registry/ranges/serviceips 和/registry/services/endpoints/default/kubernetes 这两个 key 所对应的值，效果见下面编码展示区。

```bash
# 获取"/registry/ranges/serviceips"所对应的值
# 发现这里有很多奇怪的字符=。=
# 可以大体推断出来，集群所有service的ip范围为10.96.0.0/12， 与api-server的yaml文件中配置的一致
/$ etcdctl get /registry/ranges/serviceips
/registry/ranges/serviceips
k8s
v1RangeAllocation&
"*28Bz
10.96.0.0/12"
# 获取"/registry/services/endpoints/default/kubernetes"所对应的值
# 发现这里有很多奇怪的字符=。=
# 在default命名空间的kubernetes这个service所对应的endpoint有两个ip
# 分别为192.168.205.137和192.168.205.139
/$ etcdctl get /registry/services/endpoints/default/kubernetes
/registry/services/endpoints/default/kubernetes
k8s
v1	Endpoints�
O
kubernetesdefault"*$0b6bb724-f066-11e8-be14-000c29d2cb3a2ں��z;
192.168.205.137
192.168.205.139
https�2TCP"
```

可以很明显看出来，ETCD 中保存的并不是输出友好的数据（比如 json 或 xml 等就是输出友好型数据）。当然，如果进一步研究可以知道，ETCD 保存的是 Protocol Buffers 序列化后的值。如果大家对 Protobuf 有研究，可以知道这个协议也是个 key-value 的协议，只不过会把其 key-value 值按照特定的算法进行压缩，不过并没有压缩的很厉害，显式输出这些值多少也能获取到一些信息；比如 /registry/services/endpoints/default/kubernetes 对应的 192.168.205.137、192.168.205.139 等值。

# ETCD 中其他的 key 及其 value

通过上面的探索，对 ETCD 中存储的数据有了大体的了解，接下来就可以开始更加刺激的冒险了。

## 据说 flannel 需要使用 ETCD 保存网络元数据

那么，flannel 在 ETCD 中保存的数据是什么，保存在哪个 key 中了呢？下面把所有网关相关的几个关键词 canal|flannel|calico 输出可以知道，里面只有一个可能包含 flannel 所需数据的 key，即/registry/configmaps/kube-system/canal-config，输出内容后对比关于 flannel 的 etcd 配置这篇文章，很大程度可以认为就是它了（需要进一步去 canal 项目的源码中去确认）。

```bash
/$ etcdctl get "" --prefix --keys-only |grep -Ev "^$"|grep"canal|flannel|calico"
# ... 忽略很多条/registry/configmaps/kube-system/canal-config
# ... 忽略很多条
# 可以看到里面有一个配置项 net-conf， 对比flannel的配置，可以认为这个地方很可能就是canal项目中flannel在etcd中需要的值。这里设置了网段为"10.244.0.0/16"
/$ etcdctl get /registry/configmaps/kube-system/canal-config
# ... 省略很多 net-conf.jsonI{"Network":"10.244.0.0/16", "Backend":{"Type":"vxlan"}}
```

## default命名空间中endpoint实例kubernetes的值

```
 /$ kubectl describe svc kubernetes
Name:              kubernetes
Namespace:         default
Labels:            component=apiserver
                   provider=kubernetes
Annotations:       <none>
Selector:          <none>
Type:              ClusterIP
IP:                10.96.0.1
Port:              https  443/TCP
TargetPort:        6443/TCP
Endpoints:         192.168.205.137:6443,192.168.205.139:6443
Session Affinity:  ClientIP
Events:            <none>

/$ etcdctl get "" --prefix --keys-only |grep -Ev "^$" |grep "default" |grep "kubernetes"
/registry/services/endpoints/default/kubernetes
/registry/services/specs/default/kubernetes 
```

我把`/registry/services/endpoints/default/kubernetes`输入到搜索引擎搜索了一下，发现[有人](https://github.com/kubernetes/kubernetes/issues/19989)在github上抛出类似的问题，从其`The three apiservers (Ip Adresses .31,.32,.33) are constantly overwriting the etcd-key /registry/services/endpoints/default/kubernetes`可以推测出来，这个值是api-server这个控件主动去写入的。

这能得出一个结论：在部署高可用集群时，如果想把多个api-server注册到集群，那么所有的api-server的服务都将会出现在default命名空间的kubernetes这个endpoints上；这也就意味着难以把集群中的一个api-server单独隔离出来而不让它对外提供服务（我当前想debug的一个问题需要这么操作，得出这个结论表示很无奈）。

# 小结

文本对 k8s 数据仓库 ETCD 进行了探究，总结了 ETCD 保存 k8s 数据时 key 值的规范，并尝试查看了 value 值的内容。最后对几个感兴趣的 key 值及其 value 值进行了探索。通过探究可以知道，k8s 把集群的信息非常有条理地保存在 ETCD 中，key 值定义有章可循，比较方便 debug；同时，虽然 ETCD 中的 value 值是 protobuf 序列化后的数据，不适合展示，不过输出到文本后依然有一定的参考价值。

# 参考

*   [Production-Grade Container Orchestration - Kubernetes](https://kubernetes.io/) Kubernetes官网
*   [Kubernetes是什么 \_ Kubernetes(K8S)中文文档\_Kubernetes中文社区](http://docs.kubernetes.org.cn/227.html) k8s中文文档
*   [GitHub - cookeem/kubeadm-ha: Kubernetes high availiability deploy based on kubeadm (English/中文 for v1.11.x/v1.9.x/v1.7.x/v1.6.x)](https://github.com/cookeem/kubeadm-ha/) k8s高可用部署方案
*   [Installing Calico for policy and flannel for networking](https://docs.projectcalico.org/v3.3/getting-started/kubernetes/installation/flannel) 网络插件的安装
*   [flannel Container Networking | Configuring flannel Networking](https://coreos.com/flannel/docs/latest/flannel-config.html) 描述了flannel配置etcd作为datastore的做法，可以推敲出etcd中保存的值可能的样子
*   [flannel/configuration.md at master · coreos/flannel · GitHub](https://github.com/coreos/flannel/blob/master/Documentation/configuration.md) flannel配置etcd作为datastore的文档
*   [Kubernetes service multiple apiserver endpoints · Issue #19989 · kubernetes/kubernetes · GitHub](https://github.com/kubernetes/kubernetes/issues/19989) 从这里的描述可以看出api-server本身主动向ETCD写数据