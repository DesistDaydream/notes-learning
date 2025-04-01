---
title: "灵魂拷问，上 Kubernetes 有什么业务价值？"
linkTitle: "灵魂拷问，上 Kubernetes 有什么业务价值？"
weight: 20
---

[原文链接](https://mp.weixin.qq.com/s/a3NE5fSpZIM9qlOofGTMWQ)

本文整理自 2020 年 7 月 22 日《基于 Kubernetes 与 OAM 构建统一、标准化的应用管理平台》主题线上网络研讨会。文章共分为上下两篇，本文为上篇，主要和大家介绍上Kubernetes有什么业务价值，以及什么是“以应用为中心”的 Kubernetes。下篇将跟大家具体分享如何构建“以应用为中心”的 Kubernetes。

****关注公众号，回复****“0722”****即可下载 PPT****

非常感谢大家来到 CNCF 的直播，我是张磊，阿里云的高级技术专家，Kubernetes 项目资深维护者。同时也是 CNCF 应用交付领域 co-chair。我今天给大家带来的分享主题是《基于 Kubernetes 与 OAM 构建统一、标准化的应用管理平台》。在封面上有个钉钉群组二维码。大家可以通过这个二维码进入线上交流群。

![](https://mmbiz.qpic.cn/mmbiz_png/yvBJb5Iiafvk7Pu7LACjibshibapmweer8QooolHAJELazlEgJaCtxTJSXlz0ze2ryrlyNs08awKug6GMpxVDqYDg/640?wx_fmt=png)

**上 Kubernetes 有什么业务价值？**

今天要演讲的主题是跟应用管理或者说是云原生应用交付是相关的。首先我们想要先回答这么一个问题：为什么我们要基于 Kubernetes 去构建一个应用管理平台？

![](https://mmbiz.qpic.cn/mmbiz_png/yvBJb5Iiafvk7Pu7LACjibshibapmweer8QjhcMMlwCRicqK1hhjiaYqVdUD9gSEhsLSDZvdlich8ZFgOWpzic3fomQsg/640?wx_fmt=png)

上图是一个本质的问题，我们在落地 K8s 经常遇到的一个问题。尤其是我们的业务方会问到这么一个问题，我们上 Kubernetes 有什么业务价值？这时候作为我们 K8s 工程师往往是很难回答的。原因在哪里呢？实际上这跟 K8s 的定位是相关的。K8s 这个项目呢，如果去做一个分析的话，我们会发现 K8s 不是一个 PaaS 或者应用管理的平台。实际上它是一个标准化的能力接入层。什么是能力接入层呢？大家可以看一下下图。

![](https://mmbiz.qpic.cn/mmbiz_png/yvBJb5Iiafvk7Pu7LACjibshibapmweer8Qibw7FRA1G9hpUenLEOUI0PKbdfr1OReJhtM36pjzQAswxYvM9ib5PdsA/640?wx_fmt=png)

实际上通过 Kubernetes 对用户暴露出来的是一组声明式 API，这些声明式 API 无论是 Pod 还是 Service 都是对底层基础设施的一个抽象。比如 Pod 是对一组容器的抽象，而 Deployment 是对一组 pod 的抽象。而 Service 作为 Pod 的访问入口，实际上是对集群基础设施：网络、网关、iptables 的一个抽象。Node 是对宿主机的抽象。Kubernetes 还提供了我们叫做 CRD（也就是 Custom Resource）的自定义对象。让你自己能够自定义底层基础设施的一个抽象。

而这些抽象本身或者是 API 本身，是通过另外一个模式叫做控制器(Controller)去实现的。通过控制器去驱动我们的底层基础设施向我的抽象逼近，或者是满足我抽象定义的一个终态。

所以本质来讲，Kubernetes 他的专注点是“如何标准化的接入来自于底层，无论是容器、虚机、负载均衡各种各样的一个能力，然后通过声明式 API 的方式去暴露给用户”。这就意味着 Kubernetes 实际用户不是业务研发，也不是业务运维。那是谁呢？是我们的平台开发者。希望平台开发者能够基于 Kubernetes 再去做上层的框架或者是平台。那就导致了今天我们的业务研发和业务运维对 Kubernetes 直接暴露出来的这一层抽象，感觉并不是很友好。

这里的关键点在于，Kubernetes 对这些基础设施的抽象，跟业务研发和业务运维看待系统的角度是完全不同的。这个抽象程度跟业务研发和业务运维希望的抽象程度也是不一样的。语义完全对不上，使用习惯也是有很大的鸿沟。所以说为了解决这样一个问题，都在思考一些解决方法。怎么能让我 Kubernetes 提供的基础设施的抽象能够满足我业务研发和业务运维的一个诉求呢？怎么能让 Kubernetes 能够成为业务研发和业务运维喜欢的一个平台呢？

### **方法一：把所有人都变成 Kubernetes 专家**

![](https://mmbiz.qpic.cn/mmbiz_png/yvBJb5Iiafvk7Pu7LACjibshibapmweer8QicqQTtcXytBOyTdQQe7zLsWCVoicvvVMlZYVzm6Pz1XQ0mc9LD0siavcA/640?wx_fmt=png)

假如我们所有人都是 Kubernetes 专家，那当然会喜欢 Kubernetes 对我提供的服务，这里给他发个 Kubernetes 的 PhD 博士。这里我强烈推荐阿里云和 CNCF 主办的云原生技术公开课。大家试试学完这门课程后，能不能变成 Kubernetes 专家。

这个方法门槛比较高，因为每个人对于这个系统本身感兴趣程度不太一样，学习能力也不太一样。

### **方法二：构建一个面向用户的应用管理平台**

业界常见的方法，大家会基于 Kubernetes 构建一个面向用户的应用管理平台，或者说是一个 PaaS，有人直接做成一个 Serverless。

![](https://mmbiz.qpic.cn/mmbiz_png/yvBJb5Iiafvk7Pu7LACjibshibapmweer8QYKexjNJtDyz9V5AMQ1xAh470gONwXHlgQ60WRyt6N5TIWPibNXWdRuQ/640?wx_fmt=png)

那这个具体是怎么做呢？还是在 Kubernetes 之上，会搭建一个东西叫做上层应用管理平台，这个上层应用平台对业务研发和业务运维暴露出来一个上层的 API。比如说业务研发这一侧，他不太会暴露 Pod，Deployment 这样的抽象。只会暴露出来 CI/CD 流水线。或者说一个应用，WordPress，一个外部网站，暴露出这样一个上层的概念，这是第一个部分。

第二部分，它也会给业务运维暴露出一组运维的 API。比如说：水平扩容，发布策略，分批策略，访问控制，流量配置。这样的话有一个好处，业务研发和业务运维面对的 API 不是 Kubernetes 底层的 API，不是 Node，不是 Service，不是 Deployment，不是我们的 CRD。是这样一组经过抽象经过封装后的 API。这样的业务研发和业务运维用起来会跟他所期望的 Ops 流水线，它所熟悉的使用体检有个天然的结合点。

所以说只有这么做了之后，我们才能够跟我们的业务老大说，Kubernetes 的业务价值来了。实际上业务价值不是在 Kubernetes 这一层，而是在 Kubernetes 往上的这一层--"**你的解决方案**"。所以说这样的一个系统构建出来之后呢，实际上是对 Kubernetes 又做了一层封装。变成了很多公司都有的，比如说 Kubernetes 应用平台。这是一个非常常见的做法。相比于我们让研发运维变成 Kubernetes 专家来说会更加实际一点。

但是我们在阿里也好，在很多社区的实际场景也好，它往往会伴随着这么一个问题。这个问题是：今天 Kubernetes 的生态是非常非常繁荣的，下图是我在 CNCF 截的图，好几百个项目，几千个可以让我们 Kubernetes 即插即用的能力。比如 istio，KEDA，Promethues 等等都是 Kubernetes 的插件。正是基于这么一个扩展性非常高的声明式 API 体系才会有了这么繁荣的 Kubernetes 生态。所以可以认为 Kubernetes 能力是无限的，非常强大。

![](https://mmbiz.qpic.cn/mmbiz_png/yvBJb5Iiafvk7Pu7LACjibshibapmweer8QLHL0ibjiauseNYYibibyBwA6lazv8vlz5k9X1Fx9MWCG8FBz3aj2NqpsXQ/640?wx_fmt=png)

可是这么一个无限能力，如果对接到一个非常传统的，非常经典的一个应用管理平台。比如说我们的 PaaS 上，如 Cloud Foundry。立刻就会发现一个问题，PaaS 虽然对用户提供的是很友好的 API，但是这个 API 本身是有限的，是难以扩展的。比如说 Cloud Foundry 要给用户使用，就有 Buildpack 这么一个概念，而不是 Kubernetes 所有的能力都能给用户去使用。其实几乎所有的 PaaS 都会存在这么一个问题。它往上暴露的是一个用户的API，是不可扩展的，是个有限集。

下面一个非常庞大繁荣的 Kubernetes 生态，没办法直接给用户暴露出去。可能每使用一个插件就要重新迭代开发你的 PaaS，重新交付你的 PaaS。这个是很难接受的。

### **传统 PaaS 的“能力困境”**

这问题是一个普遍存在的问题，我们叫做传统 PaaS 的“能力困境”。

![](https://mmbiz.qpic.cn/mmbiz_png/yvBJb5Iiafvk7Pu7LACjibshibapmweer8QoDERgH8CalNEQRVsZ3UdnxkcLXerLqIUdB52bzKpxSVzPEXvaxDFaQ/640?wx_fmt=png)

本质上来说这个困境是什么意思呢？K8s 生态繁荣多样的应用基础设施能力，与业务开发人员日益增长的应用管理诉求，中间存在一个传统的 PaaS，他就会变成一个瓶颈。K8s 无限的能力无法让你的研发与运维立刻用到。所以传统 PaaS 就会成为一个显而易见的瓶颈。

这样给我带来一个思考：我们能不能抛弃传统 PaaS 的一个做法，基于 K8s 打造高可扩展的应用管理平台。我们想办法能把 K8s 能力无缝的透给用户，同时又能提供传统 PaaS 比较友好的面向研发运维的使用体验呢？

其实可以从另外一个角度思考这个问题：如何基于 K8s 打造高可扩展的应用管理平台，实际上等同于 如何打造一个“以应用为中心的”的 Kubernetes。或者说能不能基于 Kubernetes 去封装下，让它能够像 PaaS 一样，去面向我的实际用户去使用呢？这个就是我们要聊的关键点。

![](https://mmbiz.qpic.cn/mmbiz_gif/US10Gcd0tQGY9ddd5GpbmVRuaRfuaESAUBGE7uHX5G0nxxLSub2QTKZdu538V7GaHXS5jsTCebYCUibaHsjg0ow/640?wx_fmt=gif)

**什么是“以应用为中心”的 Kubernetes**

### **特征一：通过原生的声明式 API 和插件体系，暴露面向最终用户的上层语义和抽象**

![](https://mmbiz.qpic.cn/mmbiz_png/yvBJb5Iiafvk7Pu7LACjibshibapmweer8QhBSxNlRvpGWfJDkyx8ftxpO0rswAky1rVmVWho2Ey6RWqE7Dia6LBIw/640?wx_fmt=png)

我们不是说要在 Kubernetes 上盖一个 PaaS，或者说是盖一个大帽子，不干这件事情。因为 K8s 本身可以扩展，可以写一组 CRD，把我们要的 API 给装上去。比如 CI/CD 流水线，就可以像 Tektong 系统直接使用 pipeline。应用也可以通过某些项目直接暴露出来。运维这一侧的发布扩容等，都可以通过安装一个 Operator 去解决问题。当然也需要一些技术将这些运维策略绑定到应用或者流水线中。

这就是我们第一个点，以应用为中心的 K8s 首先是暴露给用户的语义和 API，而不是非常底层的，比如 Service、Node 或者是 Ingress。可能用户都不知道什么意思，也不知道怎么写的。

### **特征二：上层语义和抽象可插拔，可扩展，没有抽象程度锁定和任何能力限制**

![](https://mmbiz.qpic.cn/mmbiz_png/yvBJb5Iiafvk7Pu7LACjibshibapmweer8Q0Z4N3HbqpJnstTLAEicg9IR7QVTMWGzN6toXCticVeGHtCxLXUI2kBpQ/640?wx_fmt=png)

第二个点很重要，上层语义和抽象必须是可插拔的，必须是可扩展的，是无缝兼容利用 K8s 的可扩展能力的。并且也不应该有对抽象程度的锁定。

举个例子：比如一个应用本身既可以是 Deployment，这是一个比较低程度的抽象。也可以是 Knative Service，这是一个相对来说高程度的抽象，相对于 deployment 来说比较简单，只有一个 PodTemplate。甚至可以更简单，可以是一个 Service，或者是个 Function。这个时候抽象程度就很高。如果基于 K8s 做一个以应用为中心的框架的话，它应该是能够暴露工作负载的多种抽象程度的。而不是说单独去使用 Knative，只能暴露出 Knative Service。假如我想使用 Knative 部署一个 Statefulset，这当然是不可以的。抽象程度是完全不一致的。所以我希望这个以应用为中心的 K8s 是没有抽象程度的锁定的。

同时也不应该有能力的限制，什么叫没有能力的限制呢？比如从运维侧举个例子，运维侧有很多很多扩容策略、发布策略等等。如果我想新加一个策略能力，它应该是非常简单的，就像在 K8s 安装一个 Operator 一样非常简单，能 helm insatll 就能搞定，答案是必须的。假如需要添加一个水平扩容，直接 helm install vpa 就能解决。通过这种方式才能做一个以应用为中心的 Kubernetes。

可以看到它跟我们的传统 PaaS 还是有很大区别的，它的可扩展能力非常非常强。它本质上就是一个 K8s，但是它跟专有的 Service，Knative，OpenFaaS 也不一样。它不会把抽象程度锁定到某一种 Workload 上，你的 Workload 是可以随意去定义。运维侧的能力也可以随意可插拔的去定义。这才是我们叫做一个以应用为中心的 Kubernetes。那么这么一个 Kubernetes 怎么做呢？

后续我们将会在下篇文章中详细为大家解读如何构建“以应用为中心”的 Kubernetes？以及构建这么一个以用户为中心的 Kubernetes，需要做几个层级的事情。
