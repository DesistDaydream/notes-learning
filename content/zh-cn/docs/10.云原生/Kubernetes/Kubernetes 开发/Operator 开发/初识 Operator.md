---
title: 初识 Operator
---

Kubebuilder 代码示例详见 [GitHub 上我的 kubernetesAPI 仓库](https://github.com/DesistDaydream/kubernetesAPI/tree/master/operator)

本文我们将首先了解到 Operator 是什么，之后逐步了解到 Operator 的生态建设，Operator 的关键组件及其基本的工作原理。

## 介绍

基于 Kubernetes 平台，我们可以轻松的搭建一些简单的无状态应用，比如对于一些常见的 web apps 或是移动端后台程序，开发者甚至不用十分了解 Kubernetes 就可以利用 Deployment，Service 这些基本单元模型构建出自己的应用拓扑并暴露相应的服务。由于无状态应用的特性支持其在任意时刻进行部署、迁移、升级等操作，Kubernetes 现有的 ReplicaSets、Deployment、Services 等资源对象已经足够支撑起无状态应用对于自动扩缩容、实例间负载均衡等基本需求。

在管理简单的有状态应用时，我们可以利用社区原生的 StatefulSet 和 PV 模型来构建基础的应用拓扑，帮助实现相应的持久化存储，按顺序部署、顺序扩容、顺序滚动更新等特性。

而随着 Kubernetes 的蓬勃发展，在数据分析，机器学习等领域相继出现了一些场景更为复杂的分布式应用系统，也给社区和相关应用的开发运维人员提出了新的挑战：

- 不同场景下的分布式系统中通常维护了一套自身的模型定义规范，如何在 Kubernetes 平台中表达或兼容出应用原先的模型定义？

- 当应用系统发生扩缩容或升级时，如何保证当前已有实例服务的可用性？如何保证它们之间的可连通性？

- 如何去重新配置或定义复杂的分布式应用？是否需要大量的专业模板定义和复杂的命令操作？是否可以向无状态应用那样用一条 kubectl 命令就完成应用的更新？

- 如何备份和管理系统状态和应用数据？如何协调系统集群各成员间在不同生命周期的应用状态？

而所有的这些正是 Operator 希望解决的问题，本文我们将首先了解到 Operator 是什么，之后逐步了解到 Operator 的生态建设，Operator 的关键组件及其基本的工作原理，下面让我们来一探究竟吧。

## 初识 Operator

首先让我们一起来看下什么是 Operator 以及它的诞生和发展历程。

### 1. 什么是 Operator

CoreOS 在 2016 年底提出了 Operator 的概念，当时的一段官方定义如下：

“An Operator represents human operational knowledge in software, to reliably manage an application.”

对于普通的应用开发者或是大多数的应用 SRE 人员，在他们的日常开发运维工作中，都需要基于自身的应用背景和领域知识构建出相应的自动化任务满足业务应用的管理、监控、运维等需求。在这个过程中，Kubernetes 自身的基础模型元素已经无法支撑不同业务领域下复杂的自动化场景。

与此同时，在云原生的大背景下，生态系统是衡量一个平台成功与否的重要标准，而广大的应用开发者作为 Kubernetes 的最直接用户和服务推广者，他们的业务需求更是 Kubernetes 的生命线。于是，谷歌率先提出了 `Third Party Resource` 的概念，允许开发者根据业务需求以插件化形式扩展出相应的 K8s API 对象模型，同时提出了自定义 controller 的概念用于编写面向领域知识的业务控制逻辑，基于 Third Party Resource，Kubernetes 社区在 1.7 版本中提出了`custom resources and controllers` 的概念，这正是 Operator 的核心概念。

基于 custom resources 和相应的自定义资源控制器，我们可以自定义扩展 Kubernetes 原生的模型元素，这样的自定义模型可以如同原生模型一样被 Kubernetes API 管理，支持 kubectl 命令行；同时 Operator 开发者可以像使用原生 API 进行应用管理一样，通过声明式的方式定义一组业务应用的期望终态，并且根据业务应用的自身特点进行相应控制器逻辑编写，以此完成对应用运行时刻生命周期的管理并持续维护与期望终态的一致性。这样的设计范式使得应用部署者只需要专注于配置自身应用的期望运行状态，而无需再投入大量的精力在手工部署或是业务在运行时刻的繁琐运维操作中。

简单来看，**Operator 定义了一组在 Kubernetes 集群中打包和部署复杂业务应用的方法**，它可以方便地在不同集群中部署并在不同的客户间传播共享；同时 Operator 还提供了一套应用在运行时刻的监控管理方法，**应用领域专家通过将业务关联的运维逻辑编写融入到 operator 自身控制器中**，而一个运行中的 Operator 就像一个 7\*24 不间断工作的优秀运维团队，它可以时刻监控应用自身状态和该应用在 Kubernetes 集群中的关注事件，并在毫秒级别基于期望终态做出对监听事件的处理，比如对应用的自动化容灾响应或是滚动升级等高级运维操作。

进一步讲，Operator 的设计和实现并不是千篇一律的，开发者可以根据自身业务需求，不断演进应用的自定义模型，同时面向具体的自动化场景在控制器中扩展相应的业务逻辑。很多 Operator 的出现都是起源于一些相对简单的部署和配置需求，并在后续演进中不断完善补充对复杂运维需求的自动化处理。

### 2. Operator 的发展

时至今日，Kubernetes 已经确立了自己在云原生领域平台层开源软件中的绝对地位，我们可以说 Kubernetes 就是当今容器编排的事实标准；而在 Kubernetes 项目强大的影响力下，越来越多的企业级分布式应用选择拥抱云原生并开始了自己的容器化道路，而 Operator 的出现无疑极大的加速了这些传统的复杂分布式应用的上云过程。无论在生态还是生产领域，Operator 都是容器应用部署上云过程中广受欢迎的实现规范，本小节就让我们来一起回顾下 Operator 的诞生和发展历史。

2014 到 2015 年，Docker 无疑是容器领域的绝对霸主，容器技术自身敏捷、弹性和可移植性等优势使其迅速成为了当时炙手可热的焦点。在这个过程中，虽然市场上涌现了大量应用镜像和技术分享，我们却很难在企业生产级别的分布式系统中寻找到容器应用的成功案例。容器技术的本质是提供了主机虚拟层之上的隔离，这样的隔离虽然带来了敏捷和弹性的优势，但同时也给容器和外部世界的交互带来了多一层的障碍；尤其是面向复杂分布式系统中，在处理自身以及不同容器间状态的依赖和维护问题上，往往需要大量的额外工作和依赖组件。这也成为了容器技术在云原生应用生产化道路上的一个瓶颈。

与此同时，谷歌于 2014 年基于其内部的分布式底层框架 Borg 推出了 Kubernetes 并完成了第一次代码提交。

2015 年，Kubernetes v1.0 版本正式发布，同时云原生计算基金会 Cloud Native Computing Foundation，简称 CNCF）正式成立，基于云原生这个大背景，CNCF 致力于维护和集成优秀开源技术以支撑编排容器化微服务架构应用。

2016 年是 Kubernetes 进入主干道，开始蓬勃发展的一年。这一年的社区，开发者们从最初的种种疑虑转为对 Kubernetes 的大力追捧，无论从 commit 数量到个人贡献者数量都有了显著增长；同时越来越多的企业选择 Kubernetes 作为生产系统容器集群的编排引擎，而以 Kubernetes 为核心构建企业内部的容器生态已经开始逐渐成为云原生大背景下业界的共识。也正是在这一年，CoreOS 正式推出了 Operator，旨在通过扩展 Kubernetes 原生 API 的方式为 Kubernetes 应用提供创建、配置以及运行时刻生命周期管理能力，与此同时用户可以利用 Operator 方便的对应用模型进行更新、备份、扩缩容及监控等多种复杂运维操作。

在 Kubernetes 实现容器编排的核心思想中，会使用控制器（Controller）模式对 etcd 里的 API 模型对象变化保持不断的监听（Watch），并在控制器中对指定事件进行响应处理，针对不同的 API 模型可以在对应的控制器中添加相应的业务逻辑，通过这种方式完成应用编排中各阶段的事件处理。而 Operator 正是基于控制器模式，允许应用开发者通过扩展 Kubernetes API 对象的方式，将复杂的分布式应用集群抽象为一个自定义的 API 对象，通过对自定义 API 模型的请求可以实现基本的运维操作，而在 Controller 中开发者可以专注实现应用在运行时刻管理中遇到的相关复杂逻辑。

在当时，率先提出这种扩展原生 API 对象进行应用集群定义框架的并不是 CoreOS，而是当时还在谷歌的 Kubernetes 创始人 Brendan Burns；正是 Brendan 早在 1.0 版本发布前就意识到了 Kubernetes API 可扩展性对 Kubernetes 生态系统及其平台自身的重要性，并构建了相应的 API 扩展框架，谷歌将其命名为 `Third Party Resource`，简称“TPR”。CoreOS 是最早的一批基于 Kubernetes 平台提供企业级容器服务解决方案的厂商之一，他们很敏锐地捕捉到了 TPR 和控制器模式对企业级应用开发者的重要价值；并很快基于 TPR 实现了历史上第一个 Operator：`etcd-operator`。它可以让用户通过短短的几条命令就快速的部署一个 etcd 集群，并且基于 kubectl 命令行一个普通的开发者就可以实现 etcd 集群滚动更新、灾备、备份恢复等复杂的运维操作，极大的降低了 etcd 集群的使用门槛，在很短的时间就成为当时 K8s 社区关注的焦点项目。

与此同时，Operator 以其插件化、自由化的模式特性，迅速吸引了大批的应用开发者，一时间很多市场上主流的分布式应用均出现了对应的 Operator 开源项目；而很多云厂商也迅速跟进，纷纷提出基于 Operator 进行应用上云的解决方案。Operator 在 Kubernetes 应用开发者中的热度大有星火燎原之势。

虽然 Operator 的出现受到了大量应用开发者的热捧，但是它的发展之路并不是一帆风顺的。对于谷歌团队而言，Controller 和控制器模式一直以来是作为其 API 体系内部实现的核心，从未暴露给终端应用开发者，Kubernetes 社区关注的焦点也更多的是集中在 PaaS 平台层面的核心能力。而 Operator 的出现打破了社区传统意义上的格局，对于谷歌团队而言，Controller 作为 Kubernetes 原生 API 的核心机制，应该交由系统内部的 Controller Manager 组件进行管理，并且遵从统一的设计开发模式，而不是像 Operator 那样交由应用开发者自由地进行 Controller 代码的编写。

另外 Operator 作为 Kubernetes 生态系统中与终端用户建立连接的桥梁，作为 Kubernetes 项目的设计和捐赠者，谷歌当然也不希望错失其中的主导权。同时 Brendan Burns 突然宣布加盟微软的消息，也进一步加剧了谷歌团队与 Operator 项目之间的矛盾。

于是，2017 年开始谷歌和 RedHat 开始在社区推广 Aggregated apiserver，应用开发者需要按照标准的社区规范编写一个自定义的 apiserver，同时定义自身应用的 API 模型；通过原生 apiserver 的配置修改，扩展 apiserver 会随着原生组件一同部署，并且限制自定义 API 在系统管理组件下进行统一管理。之后，谷歌和 RedHat 开始在社区大力推广使用聚合层扩展 Kubernetes API，同时建议废弃 TPR 相关功能。

然而，巨大的压力并没有让 Operator 昙花一现，就此消失。相反，社区大量的 Operator 开发和使用者仍旧拥护着 Operator 清晰自由的设计理念，继续维护演进着自己的应用项目；同时很多云服务提供商也并没有放弃 Operator，Operator 简洁的部署方式和易复制，自由开放的代码实现方式使其维护住了大量忠实粉丝。在用户的选择面前，强如谷歌，红帽这样的巨头也不得不做出退让。最终，TPR 并没有被彻底废弃，而是由 `Custom Resource Definition`（简称 CRD）这个如今已经广为人知的资源模型范式代替。

CoreOS 官方博客也第一时间发出了回应文章指导用户尽快从 TPR 迁移到 CRD：<https://coreos.com/blog/custom-resource-kubernetes-v17> 。

2018 年初，RedHat 完成了对 CoreOS 的收购，并在几个月后发布了 Operator Framework，通过提供 SDK 等管理工具的方式进一步降低了应用开发与 Kubernetes 底层 API 知识体系之间的依赖。至此，Operator 进一步巩固了其在 Kubernetes 应用开发领域的重要地位。

### 3. Operator 的社区与生态

Operator 开放式的设计模式使开发者可以根据自身业务自由的定义服务模型和相应的控制逻辑，可以说一经推出就在社区引起了巨大的反响。

一时间，基于不同种类的业务应用涌现了一大批优秀的开源 Operator 项目，我们可以找到其中很多的典型案例，例如对于运维要求较高的数据库集群，我们可以找到像 etcd、Mysql、PostgreSQL、Redis、Cassandra 等很多主流数据库应用对应的 Operator 项目，这些 Operator 的推出有效的简化了数据库应用在 Kubernetes 集群上的部署和运维工作；在监控方向，CoreOS 开发的 prometheus-operator 早日成为社区里的明星项目，Jaeger、FluentD、Grafana 等主流监控应用也或由官方或由开发者迅速推出相应的 Operator 并持续演进；在安全领域，Aqua、Twistlock、Sisdig 等各大容器安全厂商也不甘落后，通过 Operator 的形式简化了相对门槛较高的容器安全应用配置，另外社区中像 cert-manager、vault-operator 这些热门项目也在很多生产环境上得到了广泛应用。

可以说 **Operator 在很短的时间就成为了分布式应用在 Kubernetes 集群中部署的事实标准**，同时 Operator 应用如此广泛的覆盖面也使它超过了分布式应用这个原始的范畴，成为了整个 Kubernetes 云原生应用下一个重要存在。

随着 Operator 的持续发展，已有的社区共享模式已经渐渐不能满足广大开发者和 K8s 集群管理员的需求，如何快速寻找到业务需要的可用 Operator？如何给生态中大量的 Operator 定义一个统一的质量标准？这些都成为了刚刚完成收购的 RedHat 大佬们眼中亟需解决的问题。

于是我们看到 RedHat 在年初联合 AWS、谷歌、微软等大厂推出了 OperatorHub.io，希望其作为 Kubernetes 社区的延伸，向广大 operator 用户提供一个集中式的公共仓库，用户可以在仓库网站上轻松的搜索到自己业务应用对应的 Operator 并在向导页的指导下完成实例安装；同时，开发者还可以基于 Operator Framework 开发自己的 Operator 并上传分享至仓库中。

下图为一个 Operator 项目从开发到开源到被使用的全生命周期流程：

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/wmo1ab/1616115017879-73ed4f47-23c8-409a-9deb-8f2a2475e386.png)

Operator 开源生命周期流程图

主要流程包括：

- 开发者首先使用 Operator SDK 创建一个 Operator 项目；
- 利用 SDK 我们可以生成 Operator 对应的脚手架代码，然后扩展相应业务模型和 API，最后实现业务逻辑完成一个 Operator 的代码编写；
- 参考社区测试指南进行业务逻辑的本地测试以及打包和发布格式的本地校验；
- 在完成测试后可以根据规定格式向社区提交 PR，会有专人进行 review；
- 待社区审核通过完成 merge 后，终端用户就可以在 OperatorHub.io 页面上找到业务对应的 Operator；
- 用户可以在 OperatorHub.io 上找到业务 Operator 对应的说明文档和安装指南，通过简单的命令行操作即可在目标集群上完成 Operator 实例的安装；
- Operator 实例会根据配置创建所需的业务应用，OLM 和 Operator Metering 等组件可以帮助用户完成业务应用对应的运维和监控采集等管理操作。

## 总结

本文主要介绍了 Operator 的基本概念，让您了解 Operator 的应用场景和发展历程。Operator 已经成为 Kubernetes 生态的一个重要设计模式，Kubernetes 从 PaaS 层面提供整套集群、应用编排的框架，而用户通过 Operator 的方式扩展自己的应用，并实现与 Kubernetes 的融合。文章同时也介绍了使用 Operator 的生命周期流程，您可以结合自己的业务场景实现自己的 Operator 组件。
