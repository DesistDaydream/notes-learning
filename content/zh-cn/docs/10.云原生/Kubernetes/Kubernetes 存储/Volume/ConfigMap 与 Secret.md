---
title: ConfigMap 与 Secret
---

# 概述

> 参考：
>
> - [官方文档，概念-存储-卷-configMap](https://kubernetes.io/zh/docs/concepts/storage/volumes/#configmap)
> - [官方文档，任务-配置 Pods 和容器-使用 ConfigMap 配置 Pod](https://kubernetes.io/zh/docs/tasks/configure-pod-container/configure-pod-configmap/)

ConfigMap 与 Secret 这两种资源是 Kubernetes 的配置管理中心，是一种特殊类型的 Volume。用来提供给从 k8s 集群外部到 pod 内部的应用，注入各种信息(配置文件、变量等)的功能。

这种类型的 Volume 不是为了存放 Container 中的数据，也不是用来进行 Container 与 Host 之间的数据交换。而是用来为 Container 提供预先定义好的数据。从 Container 的角度看，这些 Volume 里的信息就仿佛是被 Kubernetes "投射"进容器当中一样。

## 为什么需要 ConfigMap？

1. 几乎所有的应用开发中，都会涉及到配置文件的变更，比如说在 web 的程序中，需要连接数据库，缓存甚至是队列等等。而我们的一个应用程序从写第一行代码开始，要经历开发环境、测试环境、预发布环境直到最终的线上环境。而每一个环境都要定义其独立的各种配置。如果我们不能很好的管理这些配置文件，运维工作将顿时变的无比的繁琐。为此业内的一些大公司专门开发了自己的一套配置管理中心，如 360 的 Qcon，百度的 disconf 等。很多应用程序的配置需要通过配置文件，命令行参数和环境变量的组合配置来完成（“十二要素应用”等均要求去配置）。这些配置应该从 image 内容中解耦，以此来保持容器化应用程序的可移植性。

2. Kubernetes 也提供了自己的一套方案，即 ConfigMap。kubernetes 通过 ConfigMap 来实现对容器中应用的配置管理，configmap 是 kubernetes 中的一个资源，可以通过 yaml 来进行配置。每个运行 Pod 的环境都可以有自己的一套 configmap，只需要当 Pod 运行在此环境的时候，自动加载对应的 configmap 即可实现 Pod 中 container 的配置变更。secret 是加密的 configmap。

# ConfigMap

ConfigMap 的特性与用法：(是 kubernetes 集群中的一个资源，作为 kubernetes 的配置中心),使用 kubectl get configmaps -A 命令查看 configmap 的资源列表

特性：

1. ConfigMap 资源提供了将配置数据注入容器的方式，同时保证该机制对容器来说是透明的。ConfigMap 可以被用来保存单个属性，也可以用来保存整个配置文件或者 JSON 二进制大对象。

2. ConfigMap 资源存储 **Key/Value Pairs(键/值对)** 配置数据，这些数据可以在 Pod 里使用。KEY/VAL PAIR 的 KEY 由自己定义，类似于变量的变量名；KEY/VAL PAIR 的 VAL 类似于变量的值。一个 KEY/VAL PAIR 就是 configmap 中 data 字段中的一个数据，当需要使用 configmap 中的数据时。只需要引用其中一个 KEY，就等于把 KEY 所对应的 VAL 交给引用方

用法：

1. ConfigMap 作为 container 内变量使用。KEY 是 Pod 中定义 env 字段中 key 的值，VAL 是 Pod 定义后 container 显示出的 Pod 中的 env 定义的变量的 VAL。即引用时通过 KEY 名引用，然后显示 VAL。

2. configmap 作为命令行的参数使用。注意：需要先把数据保存在变量中，再引用变量作为命令行参数

3. configmap 作为 volume 挂载使用。KEY 是 container 中的文件名，VAL 是该文件的内容。Pod 中的 volumes 的类型设定为 configMap,选择引用的 configMap 的名称。当 configmap 的 data 中有多个 KEY/VAL PAIR 时，每个 KEY 都是一个文件名。

   1. pod 使用 ConfigMap，通常用于：设置环境变量的值、设置命令行参数、创建配置文件。
      1. 当需要修改 container 中的配置文件时，只需要修改 container 引用的 configmap 中的 KEY 的 VAL，过一会就会在 container 中生效

Configmap 的理解概述：可以这么说，ConfigMap 也是存储数据的一种东西，只不过既不是外部的硬件(比如磁盘阵列、存储设备)，也不是外部软件(比如分布式文件系统)，而是 kubernetes 自带的一种键值型存储工具，所以 volume 同样可以把 configmap 当做卷来使用并挂载进 pod 中，还可以让 pod 把 configmap 当做外部数据，引入后作为变量使用。

# secret # 与 ConfigMap 类似，不过是把数据进行 base64 编码之后保存
