---
title: 如何构建以应用为中心的“Kubernetes”?（内含 QA 整理）
---

参考：[原文链接](https://mp.weixin.qq.com/s/ql_AIFc0s5HwZgsML63zQA)
本文整理自 2020 年 7 月 22 日《基于 Kubernetes 与 OAM 构建统一、标准化的应用管理平台》主题线上网络研讨会。

文章共分为上下两篇。上篇文章《[**灵魂拷问，上 Kubernetes 有什么业务价值？**](http://mp.weixin.qq.com/s?__biz=MzUzNzYxNjAzMg==&mid=2247492713&idx=1&sn=63d26542a935a6b3d1cfd7a72f71425b&chksm=fae6efa6cd9166b0c66e73ad47be04d029d40066b7697f2f4c7cd7a53d08ba6e019419166bb8&scene=21#wechat_redirect)》，主要和大家介绍了上 Kubernetes 有什么业务价值，以及什么是 “以应用为中心” 的 Kubernetes。本文为下篇，将跟大家具体分享如何构建 “以应用为中心” 的 Kubernetes。

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/527401b0-9af0-4240-98b5-5f8936c0d0c2/640)

**如何构建 “以应用为中心” 的 Kubernetes？**

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/527401b0-9af0-4240-98b5-5f8936c0d0c2/640)

构建这么一个以用户为中心的 Kubernetes，需要做几个层级的事情。

### **1. 应用层驱动**

首先来看最核心的部分，上图中蓝色部分，也就是 Kubernetes。可以在 Kubernetes 之上定义一组 CRD 和 Controller。可以在 CRD 来做用户这一侧的 API，比如说 pipeline 就是一个 API，应用也是一个 API。像运维侧的扩容策略这些都是可以通过 CRD 的方式安装起来。

### **2. 应用层抽象**

所以我们的需要解决第一个问题是应用抽象。如果在 Kubernetes 去做应用层抽象，就等同于定义 CRD 和 Controller，所以 Controller 可以叫做应用层的抽象。本身可以是社区里的，比如 Tekton，istio 这些，可以作为你的应用驱动层。这是第一个问题，解决的是抽象的问题。不是特别难。

### **3. 插件能力管理**

很多功能不是 K8s 提供的，内置的 Controller 还是有限的，大部分能力来自于社区或者是自己开发的 Controller。这时我的集群里面就会安装好多好多插件。如果要构建以应用为中心的 Kubernetes，那我必须能够管理起来这些能力，否则整个集群就会脱管了。用户想要这么一个能力，我需要告诉他有或者是没有。需要暴露出一个 API 来告诉他，集群是否有他需要的能力。假设需要 istio 的流量切分，需要有个接口告诉用户这个能力存不存在。不能指望用户去 get 一下 crd 合不合适，检查 Controller 是否运行。这不叫以应用为中心的 K8s，这叫裸 K8s。

所以必须有个能力，叫做插件能力管理。如果我装了 Tekton，kEDA，istio 这些组件，我必须将这些组件注册到能力注册中心，让用户能够发现这些能力，查询这些能力。这叫做：插件能力管理。

### **4. 用户体验层**

有了应用层驱动，应用层抽象，插件能力管理，我们才能更好地去考虑，如何给用户暴露一个友好的 API 或者是界面出来。有这么几种方式，比如 CLI 客户端命令行工具，或者是一个 Dashboard，又或者是研发侧的 Docker Compose。或者可以让用户写代码，用 python 或者 go 等实现 DSL，这都是可以的。

用户体验层怎么做，完全取决于用户接受什么样的方式。关键点在于以应用为中心的 Kubernetes，UI 层就可以非常方便的基于应用层抽象去做。比如 CLI 就可以直接创建一个流水线和应用，而不是兜兜转转去创建 Deployment 和 Pod，这两个的衔接方式是完全不一样的。pipeline 只需要生成一下就结束了。然后去把 Pod 和 Deployment 组成一个 Pipeline，那这个工作就非常繁琐了。这是非常重要的一点，当你有了应用层驱动，应用层抽象，插件能力管理，再去构建用户体验层就会非常非常简单。

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/527401b0-9af0-4240-98b5-5f8936c0d0c2/640)

**Open Application Model(OAM)**

如果想构建一个应用为中心的 Kubernetes，有没有一个标准化的、简单的方案呢？

**下面就要为大家介绍：** **Open Application Model(OAM)。**

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/527401b0-9af0-4240-98b5-5f8936c0d0c2/640)

OAM 的本质是帮助你构建一个 “以应用为中心 “的 Kubernetes 标准规范和框架，相比较前面的方案，OAM 专注于做这三个层次。

### **1. 应用组件 Components**

第一个叫做应用层抽象，OAM 对用户暴露出自己定义的应用层抽象，第一个抽象叫做 Components。Components 实际上是帮助我们定义 Deployment、StatefulSet 这样的 Workload 的。暴露给用户，让他去定义这些应用的语义。

### **2. 应用特征 Traits**

第二个叫做应用特征，叫做 Traits。运维侧的概念，比如扩容策略，发布策略，这些策略通过一个叫做 Traits 的 API 暴露给用户。首先 OAM 给你做了一个应用层定义抽象的方式，分别叫做 Components 和 Traits。由于你需要将 Traits 应用特征关联给应用组件 Components，例如 Deployment 需要某种扩容策略或者是发布策略，怎么把他们关联在一起呢？

### **3. 应用配置 Application Configuration**

这个就需要第三种配置叫做 Application Configuration 应用配置。最终这些概念和配置都会变成 CRD，如果你的 K8s 里面安装了 OAM 的 Kubernetes Runtime 组件，那么那就能解析你 CRD 定义的策略和 Workload，最终去交给 K8s 去执行运行起来。就这么一个组件帮助你更好地去定义抽象应用层，提供了几个标准化的方法。

### **4. 能力定义对象 Definitions**

这些抽象和能力交给 K8s 去处理之后，我这些能力需要的 Controller 插件在哪？有没有 Ready？这些版本是不是已经有了，能不能自动去安装。这是第四个能力了：能力定义对象。这是 OAM 提供的最后一个 API，通过这个 API 可以自己去注册 K8s 所有插件，比如 Tekton、KEDA、istio 等。

把它注册为组件的一个能力，或者是某一个特征。比如说 Flager，可以把它注册为金丝雀发布的能力，用户只要发现这个发布策略存在，说明这个集群支持 Flager，那么他就可以去使用。这就是一个以应用为中心的一个玩法。以用户侧为出发点，而不是以集群侧为出发点，用户侧通过一个上层的 api，特征和组件来去了解他的系统，去操作他的系统。以上就是 OAM 提供的策略和方法。

总结下来就是 OAM 可以通过标准化的方式帮助平台构建者或者开发者去定义用户侧，应用侧的抽象。第二点是提供了插件化能力注册于管理机制。并且有了这些抽象和机制之后，可以非常方便的构建可扩展的 UI 层。这就是 OAM 最核心的功能和价值。

### **5. OAM 会怎样给用户提供一个 API 呢？**

### **1）Components**

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/527401b0-9af0-4240-98b5-5f8936c0d0c2/640)

Component 是工作负载的版本化定义，例如上图中创建一个 Component，实际上就是创建一个 Deployment。这样一个 Component 交给 K8s 之后，首先会创建一个 Component 来管理这个 Workload，当你修改 Component 之后就会生成一个对应版本的 deployment。这个 Component 实际上是 Deployment 的一个模板。比如我把 image 的版本修改一下，这个操作就会触发 OAM 插件，生成一个新的版本的 Deployment，这是第一个点。其实就版本化管理机制去管理 Component。

第二点是 Workload 部分完全是自定义的，或者是是可插拔的。

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/527401b0-9af0-4240-98b5-5f8936c0d0c2/640)

今天可以定义为 Deployment，明天可以定义为一个非常简单的版本。也就是说我 Components 的抽象程度完全取决于用户自己决定的。后期也可以改成 Knative Service，甚至改成一个 Open PaaS。所以说在 Components 的 Workload 部分你可以自由的去定义自己的抽象。只要你提前安装了对应 CRD 即可，这是一个非常高级的玩法。

此外在 OAM 中，” 云服务 “也是一种 Workload， 只要你能用 CRD 定义你的云服务，就可以直接在 OAM 中定义为一个应用所依赖的组件。比如上图中的 redis 实际上是阿里云的 Redis 服务，大概是这么一个玩法。

### **2）Trait 和 Application Configuration**

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/527401b0-9af0-4240-98b5-5f8936c0d0c2/640)

首先 Trait 是声明式运维能力的描述，其实就是 Kubernetes API 对象。任何管理和运维 Workload 的组件和能力，都可以以这种 CER 的方式定义为一个 Trait。所以像 HPA，API gateway，istio 里面的 Virtual Services 都是 Trait。

Application Configuration 就像是一个信封，将 Traits 绑定给 Component，这个是显式绑定的。OAM 里面不建议去使用 Label 这样的松耦合的方式去关联你的工作负载。建议通过这种结构化的方式，通过 CRD 去显式的绑定你的特征和工作负载。这样的好处是我的绑定关系是可管理的。可以通过 kubectl get 看到这个绑定关系。作为管理员或者用户，就非常容易的看到某一个组件绑定的所有运维能力有哪些，这是可以直接展示出来的，如果通过 label 是很难做到的。同时 Label 本身有个问题是，本身不是版本化的，不是结构体，很难去升级，很难去扩展。通过这么结构化定义，后面的升级扩展将会变得非常简单。

在一个用户配置里面，可以关联多个 Components。它认为一个应用运行所需要的所有组件和所依赖的运维能力，都应该定义为一个文件叫做 ApplicationConfiguration。所以在任何环境，只要拥有这个文件，提交之后，这个应用就会生效了。OAM 是希望能够提供一个自包含的应用声明方式。

### **3）Definition Object**

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/527401b0-9af0-4240-98b5-5f8936c0d0c2/640)

除此之外，还提到了对应管理员提供了 Definition Object，这是用来注册和发现插件化能力的 API 对象。

比如我想讲 Knative Service 定义为平台支持的一种工作负载，如上图只需要简单的写一个文件即可。其中在 definitionRef 中引用 service.serving.knative.dev 即可。这样的好处就是可以直接用 kubectl get Workload 查看 Knative Service 的 Workload。所以这是一个用来注册和发现插件化能力的机制，使得用户非常简单的看到系统中当前有没有一个工作负载叫做 Knative Service。而不是让用户去看 CRD，看插件是否安装，看 Controller 是否 running，这是非常麻烦的一件事情。所以必须有这么一个插件注册和发现机制。

这一部分还有其他额外的能力，可以注册 Trait，并且允许注册的 Trait-A 和 Trait-B 是冲突的。这个信息也能带进去，这样部署的时候检查到 A 和 B 是冲突的，会产生报错信息。否则部署下去结果什么都不知道，两个能力是冲突的，赶紧删了回滚重新创建。OAM 在注册的时候就会暴露出来运维能力的冲突，这也是靠 Definition 去做的。

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/527401b0-9af0-4240-98b5-5f8936c0d0c2/640)

除此之外，OAM 的 model 这层其他的一些附加能力，能够让你定义更为复杂的应用。

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/527401b0-9af0-4240-98b5-5f8936c0d0c2/640)

**总结**

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/527401b0-9af0-4240-98b5-5f8936c0d0c2/640)

前面我们提到很多企业等等都在基于 Kubernetes 去构建一个上层应用管理平台。Kubernetes 实际上是面向平台开发者，而不是面向研发和应用运维的这么一个项目。它天生就是这么设计的，所以说需要基于 Kubernetes 去构建应用管理平台。去更好的服务研发和运维，这也是一个非常自然的选择。不是说必须使用 K8s 去服务你的用户。如果你的用户都是 K8s 专家，这是没问题的。如果不是的话，你去做这样一个应用平台是非常自然的事情。

但是我们不想在 K8s 之前架一个像 Cloud Foundry 传统的 PaaS。因为它会把 K8s 的能力完全遮住。它有自己的一套 API，自己的理念，自己的模型，自己的使用方式。跟 Kubernetes 都是不太一样的，很难把 Kubernetes 的能力给暴露出去。这是经典 PaaS 的一个用法，但是我们不想要这么一个理念。我们的目标是既能给用户提供一个使用体验，同时又能把 Kubernetes 的能力全部发挥出来。并且使用体验跟 Kubernetes 是完全一致的。OAM 本质上要做的是面向开发和运维的，或者说是面向以应用为中心的 Kubernetes。

所以今天所介绍的 OAM 是一个统一、标准、高可扩展的应用管理平台，能够以应用为中心的全新的 Kubernetes，这是今天讨论的一个重点。OAM 这个项目就是支撑这种理念的核心依赖和机制。简单地来说 OAM 能够让你以统一的，标准化的方式去做这件事情。比如标准化定义应用层抽象，标准化编写底层应用驱动，标准化管理 K8s 插件能力。

对于平台工程师来说，日常的工作能不能以一个标准化的框架或者依赖让平台工程师更简单更快的做这件事情。这就是 OAM 给平台工程师带来的价值。当然它也有些额外的好处，基于 OAM 暴露出来的新的 API 之后，你上层的 UI 构建起来会非常简单。

你的 OAM 天然分为两类，一类叫做工作负载，一类叫做运维特征。所以你的 UI 这层可以直接去对接了，会减少很多前端的工作。如果基于 CI/CD 做 GitOps / 持续集成发现也会变得非常简单。因为它把一个应用通过自包含的方式给定义出来了，而不是说写很多个 yaml 文件。并且这个文件不仅自包含了工作负载，也包括了运维特征。所以创建好了这个文件往 Kubernetes 中提交，这个应用要做金丝雀发布或者是蓝绿发布，流量控制，全部是清清楚楚的定义在这个应用配置文件里面的。因为 GitOps 也好，持续集成也好，是不想管你的 pod 或者是 Deployment 怎么生成的，这个应用怎么运维，怎么 run 起来，还是要靠 Kubernetes 插件或者内置能力去做的。这些能力都被定义到一个自包含的文件，适用于所有集群。所以这就会使得你的 GitOps 和持续集成变得简单。

以上就是 OAM 给平台工程师带来的一些特有的价值。简单来说是统一、标准的 API，区分研发和运维策略，让你的 UI 和 GitOps 特别容易去构建。另一点是向下提供了高可扩展的管理 K8s 插件能力。这样的系统真正做到了标准，自运维，一个以应用为中心和用户为中心的 Kubernetes 平台。

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/527401b0-9af0-4240-98b5-5f8936c0d0c2/640)

**OAM 社区**

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/527401b0-9af0-4240-98b5-5f8936c0d0c2/640)

上面最后希望大家踊跃加入 OAM 社区，参与讨论。上图中有钉钉群二维码，目前人数有几千人，讨论非常激烈，我们会在里面讨论 GitOps，CI/CD，构建 OAM 平台等等。OAM 也有亚太地区的周会，大家可以去参加。上面的链接是开源项目地址，将这个安装到 Kubernetes 中就可以使用上面我们说的这些能力了。

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/527401b0-9af0-4240-98b5-5f8936c0d0c2/640)

**QA 环节**

**Q1：** 例子中提问到了 Function 的例子，是否可以理解为 Serverless 或者是 PaaS？

**A1\*\***：\*\* 这样理解是没错的，可以理解为阿里云的一个 Function，或者是 Knative Service。

**Q2：** 有没有可以让我自由定义出相应的规则这种规范？

**A2：** 有的，在 OAM 里面有个规范，叫做 spec。spec 里面有提交容器化的规范。后面会增加更多抽象的规范。当然也分类别，有一些是非常标准化的，需要严格遵守。有一些是比较松的，可以不用严格遵守。

**Q3：** docker-compose 的例子可否再谈谈？

**A3：** 本次 ppt 中没有 docker-compose 的例子，但是这个其实很容易去理解，因为 OAM 将 Kubernetes API 分为两类，一个叫做 Components，一个叫 T raits。有这么一个 Componets 文件，就可以直接映射 OAM 的概念，docker-compose 中有个概念叫做 Service，其实就是对应了 OAM 中的 Component。这完全是一对一对应关系。Service 下面有个 Deployment，有个部署策略，其实对应的就是 OAM 的 Trait。

**Q4：** 定义阿里云的 redis 是否已经实现了？

**A4：** 已经实现了，但是功能有限。内部已经实现了一个更强大的功能，通过 OAM 将阿里云的所有资源给创建起来。目前这个是在 Crossplane 去做的。但是内部更完整的实现还没有完全的放出去。我们还在规划中，希望通过一个叫做 Alibaba Opreator 的方式暴露出去。

**Q5：** 是否可以理解 OAM 通过管理元数据通过编写 CRD 来打包 Components 和 Traits。

**A5：** 可以说是对的。你把自己的 CRD 也好，社区里面的 CRD 也好，稍微做个分类或者封装，暴露给用户。所以对于用户来说只要理解两个概念——Components 和 Traits。Components 里面的内容是靠你的 CRD 来决定的，所以说这是一个比较轻量级的抽象。

**Q6：** 假设 Components 有四个，Traits 有五个，是否可以理解为可封装能力有 20 项。

**A6：** 这个不是这么算的，不管有多少 Components 和 Trait，最终有几个能力取决于你注册的实际 CRD。Components 和 Traits 与背后的能力是解耦开的。

**Q7：** OAM 能使用 Kustomize 生成么？

**A7：** 当然可以了，Kustomize 使一个 yaml 文件操作工具。你可以用这个工具生成任何你想要的 yaml 文件，你也可以用其他的，比如 google 的另一个项目叫 kpt，比如你用 DSL，json。所有可以操作 yaml 文件的工具都可以操作 OAM 文件，OAM 的 yaml 文件跟正常的 K8s 中的 yaml 没有任何区别。在 K8s 看来 OAM 无非就是一个 CRD。

**Q8：** OAM 是否可以生产可用？

**A8：** 这里面分几个点，OAM 本身分两个部分。第一部分是规范，是处于 alpha 版本，计划在 2020 年内发布 beta 版本。beta 就是一个稳定版本，这是一个比较明确的计划。现在的 spec 是有可能会变的，但是有另外一个版本叫做 oam-kubernetes-runtime 插件，这是作为独立项目去运营的，计划在 Q3 发布稳定版本。即使我的 spec 发生的改变，但是插件会做向下兼容，保证 spec 变化不会影响你的系统，我们的 runtime 会提前发布稳定版本，应该是比较快的。如果构建平台化建议优先使用 runtime。

**Q9：** OAM 有没有稳定性考虑？比如说高可用。

**A9：** 这个是有的，目前 runtime 这个项目就在做很多稳定性的东西，这是阿里内部和微软内部的一个诉求。这块都是在做，肯定是有这方面考虑的，包括边界条件的一个覆盖。

**Q10：** 可不可介绍下双十一的状态下，有多少个 Pod 在支持？

**A10：** 这个数量会比较大，大概在十几万这样一个规模，应用容器数也是很多的。这个对大家的参考价值不是很大，因为阿里的架构和应用跟大多数同学看到的是不太一样的，大多数是个单元化的框架，每个应用拆分的微服务非常非常细。pod 数和容器数都是非常多的。

**Q11：** 目前 OAM 只有阿里和微软，以后像 google 这些大厂会加入么？

**A11：** 一定会的，接下来的计划会引入新的合作方。目前 google 和 aws 都对 OAM 有一些社区的支持。本身作为云原生的一个规范，也是有一些想法的。在初期的时候，大厂加入的速度会比较慢，更希望的是用户使用起来。大厂并不一定是 OAM 的主要用户，他们更多的是商业考虑。

**Q12：** OAM 是否会关联 Mesh？

**A12：** 一定会的，但是并不是说直接 Mesh 一个核心能力，更多的说作为 OAM trait 使用, 比如描述一个流量的拓扑关系。

**Q13：** OAM 的高可用方案？

**A13：** OAM 本身就是个无状态服务，本身的高可用方案不是很复杂。

**Q14：** OAM 考虑是单集群还是多集群？

**A14：** 目前是单集群，但是我们马上也会发布多集群的模型，在阿里内部已经是多集群模型。简单来说多集群是两层模型。多集群的概念是定义在 Scope 里面的，通过 Scope 来决定 Workload 或者是 Components 放到哪个集群里面。我们会在社区尽快放出来。

如果有其他问题，建议大家加入我们的钉钉群进行讨论。

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/527401b0-9af0-4240-98b5-5f8936c0d0c2/640)
<https://mp.weixin.qq.com/s/ql_AIFc0s5HwZgsML63zQA>
