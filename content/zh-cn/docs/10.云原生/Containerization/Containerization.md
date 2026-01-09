---
title: Containerization
linkTitle: Containerization
weight: 1
---

# 概述

> 参考：
>
> - [Wiki, Containerization](https://en.wikipedia.org/wiki/Containerization_(computing))
> - [Wiki, OS-level virtualization](https://en.wikipedia.org/wiki/OS-level_virtualization)
> - [MoeLove，一篇搞懂容器技术的基石： cgroup](https://moelove.info/2021/11/17/%E4%B8%80%E7%AF%87%E6%90%9E%E6%87%82%E5%AE%B9%E5%99%A8%E6%8A%80%E6%9C%AF%E7%9A%84%E5%9F%BA%E7%9F%B3-cgroup/)

**Container(容器)** 是一种基础工具；泛指任何可以用于容纳其它物品的工具，可以部分或完全封闭，被用于容纳、储存、运输物品。物体可以被放置在容器中，而容器则可以保护内容物。人类使用容器的历史至少有十万年，甚至可能有数百万年的历史。

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/containerization/containerization_history.png)

自 1979 年，Unix 版本 7 引用 Chroot Jail 以及 [Chroot](/docs/1.操作系统/Kernel/Process/Chroot.md) 系统调用开始，直到 2013 年开源出的 Docker，2014 年开源出来的 Kubernetes，直到现在的云原生生态的火热。容器技术已经逐步成为主流的基础技术之一。

## 一、什么是容器

IT 里的容器技术是英文单词 Linux Container 的直译。Container 这个单词有集装箱、容器的含义（主要偏集装箱意思）。不过，在中文环境下，咱们要交流要传授，如果翻译成“集装箱技术” 就有点拗口，所以结合中国人的吐字习惯和文化背景，更喜欢用容器这个词。不过，如果要形象的理解 Linux Container 技术的话，还是得念成集装箱会比较好。我们知道，海边码头里的集装箱是运载货物用的，它是一种按规格标准化的钢制箱子。集装箱的特色，在于其格式划一，并可以层层重叠，所以可以大量放置在特别设计的远洋轮船中（早期航运是没有集装箱概念的，那时候货物杂乱无章的放，很影响出货和运输效率）。有了集装箱，那么这就更加快捷方便的为生产商提供廉价的运输服务。

因此，IT 世界里借鉴了这一理念。早期，大家都认为硬件抽象层基于 hypervisor 的虚拟化方式可以最大程度上提供虚拟化管理的灵活性。各种不同操作系统的虚拟机都能通过 hypervisor（KVM、XEN 等）来衍生、运行、销毁。然而，随着时间推移，用户发现 hypervisor 这种方式麻烦越来越多。为什么？因为对于 hypervisor 环境来说，每个虚拟机都需要运行一个完整的操作系统以及其中安装好的大量应用程序。但实际生产开发环境里，我们更关注的是自己部署的应用程序，如果每次部署发布我都得搞一个完整操作系统和附带的依赖环境，那么这让任务和性能变得很重和很低下。

基于上述情况，人们就在想，有没有其他什么方式能让人更加的关注应用程序本身，底层多余的操作系统和环境我可以共享和复用？换句话来说，那就是我部署一个服务运行好后，我再想移植到另外一个地方，我可以不用再安装一套操作系统和依赖环境。这就像集装箱运载一样，我把货物一辆兰博基尼跑车（好比开发好的应用 APP），打包放到一容器集装箱里，它通过货轮可以轻而易举的从上海码头（CentOS7.2 环境）运送到纽约码头（Ubuntu14.04 环境）。而且运输期间，我的兰博基尼（APP）没有受到任何的损坏（文件没有丢失），在另外一个码头卸货后，依然可以完美风骚的赛跑（启动正常）。

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/containerization/containerization_imagine.png)

## 二、容器技术的实现方式，lxc、runc、kata 等

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/containerization/1616122918517-21a8c653-b45e-44a9-b54e-86ca53db6fd7.png)

Linux Container(LXC) 容器技术的诞生（2008 年）就解决了 IT 世界里“集装箱运输”的问题。Linux Container（简称 LXC）它是一种 内核轻量级的操作系统层 虚拟化技术，也称为容器的运行时(runtime 运行环境)。Linux Container 主要由 Namespace 和 Cgroup 两大机制来保证实现。那么 Namespace 和 Cgroup 是什么呢？刚才我们上面提到了集装箱，集装箱的作用当然是可以对货物进行打包隔离了，不让 A 公司的货跟 B 公司的货混在一起，不然卸货就分不清楚了。那么 Namespace 也是一样的作用，做隔离。光有隔离还没用，我们还需要对货物进行资源的管理。同样的，航运码头也有这样的管理机制：货物用什么样规格大小的集装箱，货物用多少个集装箱，货物哪些优先运走，遇到极端天气怎么暂停运输服务怎么改航道等等... 通用的，与此对应的 Cgroup 就负责资源管理控制作用，比如进程组使用 CPU/MEM 的限制，进程组的优先级控制，进程组的挂起和恢复等等。

经过多年的发展，陆续推出了 runc、kata 等容器底层技术

runc 是 lxc 的替代品，官方说明：<https://www.docker.com/blog/runc/>

kata 是自带内核的虚拟机型的容器 runtime，官方网址：<https://katacontainers.io/>

## 三、容器技术的特点

容器的特点其实我们拿跟它跟硬件抽象层虚拟化 hypervisor 技术对比就清楚了，我们之前也提到过，传统的虚拟化（虚拟机）技术，创建环境和部署应用都很麻烦，而且应用的移植性也很繁琐，比如你要把 vmware 里的虚拟机迁移到 KVM 里就很繁琐（需要做镜像格式的转换）。那么有了容器技术就简单了，总结下容器技术主要有三个特点：

- 极其轻量：只打包了必要的 Bin/Lib；
- 秒级部署：根据镜像的不同，容器的部署大概在毫秒与秒之间（比虚拟机强很多）；
- 易于移植：一次构建，随处部署；
- 弹性伸缩：Kubernetes、Swam、Mesos 这类开源、方便、好使的容器管理平台有着非常强大的弹性管理能力。

## 四、容器的标准化 Open Container Initiative(OCI)

当前，docker 几乎是容器的代名词，很多人以为 docker 就是容器。其实，这是错误的认识(docker 只是可以实现容器的引擎, docker 调用 containerd，containerd 再调用 runc 来启动一个容器)。除了 docker 还有 podman 等等。所以，容器世界里并不是只有 docker 一家。既然不是一家就很容易出现分歧。任何技术出现都需要一个标准来规范它，不然各搞各的很容易导致技术实现的碎片化，出现大量的冲突和冗余。因此，在 2015 年，由 Google，Docker、CoreOS、IBM、微软、红帽等厂商联合发起的 [OCI(Open Container Initiative)](https://www.opencontainers.org/) 项目成立了，并于 2016 年 4 月推出了第一个开放容器标准。标准主要包括 runtime(运行时)标准 和 image(镜像)标准。标准的推出，有助于替成长中市场带来稳定性，让企业能放心采用容器技术，用户在打包、部署应用程序后，可以自由选择不同的容器 Runtime；同时，镜像打包、建立、认证、部署、命名也都能按照统一的规范来做。

两种标准主要包含以下内容：

- 容器运行时标准 （runtime spec）
  - creating：使用 create 命令创建容器，这个过程称为创建中 b). created：容器创建出来，但是还没有运行，表示镜像和配置没有错误，容器能够运行在当前平台 c).
  - running：容器的运行状态，里面的进程处于 up 状态，正在执行用户设定的任务 d)
  - stopped：容器运行完成，或者运行出错，或者 stop 命令之后，容器处于暂停状态。这个状态，容器还有很多信息保存在平台中，并没有完全被删除
  - ....等等
- 容器镜像标准（image spec）
  - 文件系统：以 layer 保存的文件系统，每个 layer 保存了和上层之间变化的部分，layer 应该保存哪些文件，怎么表示增加、修改和删除的文件等;
  - config 文件：保存了文件系统的层级信息（每个层级的 hash 值，以及历史信息），以及容器运行时需要的一些信息（比如环境变量、工作目录、命令参数、mount 列表），指定了镜像在某个特定平台和系统的配置。比较接近我们使用 docker inspect
  - ....等等

## 五、容器的主要应用场景

容器技术的诞生其实主要解决了 PAAS 的层的技术实现。像 OpenStack、Cloudstack 这样的技术是解决 IAAS 层的问题。IAAS 层和 PAAS 层大家估计也听得很多了，关于他们的区别和特性我这里不在描述。那么容器技术主要应用在哪些场景呢？目前主流的有以下几种：

1. 容器化传统应用 容器不仅能提高现有应用的安全性和可移植性，还能节约成本。
    - 每个企业的环境中都有一套较旧的应用来服务于客户或自动执行业务流程。即使是大规模的单体应用，通过容器隔离的增强安全性、以及可移植性特点，也能从 容器 中获益，从而降低成本。一旦容器化之后，这些应用可以扩展额外的服务或者转变到微服务架构之上。
2. 持续集成和持续部署 (CI/CD) 通过 Docker 加速应用管道自动化和应用部署，交付速度提高至少 13 倍。
    - 现代化开发流程快速、持续且具备自动执行能力，最终目标是开发出更加可靠的软件。通过持续集成 (CI) 和持续部署 (CD)，每次开发人员签入代码并顺利测试之后，IT 团队都能够集成新代码。作为开发运维方法的基础，CI/CD 创造了一种实时反馈回路机制，持续地传输小型迭代更改，从而加速更改，提高质量。CI 环境通常是完全自动化的，通过 git 推送命令触发测试，测试成功时自动构建新镜像，然后推送到 Docker 镜像库。通过后续的自动化和脚本，可以将新镜像的容器部署到预演环境，从而进行进一步测试。
3. 微服务 加速应用架构现代化进程。
    - 应用架构正在从采用瀑布模型开发法的单体代码库转变为独立开发和部署的松耦合服务。成千上万个这样的服务相互连接就形成了应用。Docker 允许开发人员选择最适合于每种服务的工具或技术栈，隔离服务以消除任何潜在的冲突，从而避免“地狱式的矩阵依赖”。这些容器可以独立于应用的其他服务组件，轻松地共享、部署、更新和瞬间扩展。Docker 的端到端安全功能让团队能够构建和运行最低权限的微服务模型，服务所需的资源（其他应用、涉密信息、计算资源等）会适时被创建并被访问。
4. IT 基础设施优化 充分利用基础设施，节省资金。
    - Docker 和容器有助于优化 IT 基础设施的利用率和成本。优化不仅仅是指削减成本，还能确保在适当的时间有效地使用适当的资源。容器是一种轻量级的打包和隔离应用工作负载的方法，所以 Docker 允许在同一物理或虚拟服务器上毫不冲突地运行多项工作负载。企业可以整合数据中心，将并购而来的 IT 资源进行整合，从而获得向云端的可迁移性，同时减少操作系统和服务器的维护工作。

# Container 的基本核心概念

## Image(镜像)

镜像就是一个只读的模板。

例如：一个镜像可以包含一个完整的 CentOS 操作系统环境，里面仅安装了 Apache 或用户需要的其他应用程序。

镜像可以用来创建 Container。

### Reference(引用)

> 参考：
>
> - <https://docs.docker.com/engine/reference/commandline/images/>

在互联网上，我们通过 **Reference(引用)** 表示唯一一个 Image，就像 URL 之于 HTTP 的 Resource 一样，**Reference 就是 Image 的 URL**。

#### Syntax(语法)

**Scheme://Registry/\[Namespace/]Repository:{Tag|Digest}**

- **Scheme://**# 访问 Registry 时所使用的协议，比如 HTTP、HTTPS
- **Registry(注册中心)** # 提供 Image 管理服务的提供商，通常是一个域名
  - 现阶段常见的 Registry 有：
    - docker.io
    - k8s.gcr.io
    - quay.io
    - ghcr.io
    - ...... 等等
- **Namespace(名称空间)** # 在一个 Registry 中可能会有多个同名的 Repository，所以需要通过 Namespace 将这些 Repository 隔开。
  - docker.io 将用户注册的账户名称作为 Namespace，若 Namespace 被省略，则 Image 就是这个 Registry 官方的。
- **Repository(仓库)** # 顾名思义，存放镜像的仓库
- **Tag(标签)** #
- **Digest(摘要)** # Image 内容的 sha256 计算结果。通常是互联网唯一的

假如我在 docker.io 注册了一个账号 lchdzh 用来存放容器镜像，有一个 k8s-debug 的镜像，版本号是 v1，我想把镜像放在 dd_k8s 仓库中。

- 那么正常的 Image Reference 是：`docker.io/lchdzh/dd_k8s:k8s-debug-v1`

但是，后来人们一般情况 Repository 都存放同一个软件的 Image，把 Tag 仅仅当做了镜像的版本

- 那么上面例子的 Image Reference 就变成了：`docker.io/lchdzh/k8s-debug:v1`

### Registry(注册中心)

Registry 可以理解为一个网站，通过 https 协议与 docker daemon 交互;也可以自己搭建私有单位 registry，提供多个功能

- 用于存储 image 的 Repository 功能，一个 Registry 上有多个 Repository
- 用户来获取 image 时的认证功能
- 当前 registry 所有 image 的索引 功能

Registry 上有多个 Repository，每个 Repository 中又包含了多个 TAG(标签)。一个 registry 中分两种：顶层仓库与用户仓库，顶层仓库里的 Repository 是这个 Registry 官方所创建的，用户仓库里的 Repository 是在该 Registry 创建的用户所创建的。image 名字中有 namespace 的就是用户仓库，没有就是顶层仓库

### Repository(仓库)

想要定位一个 Registry 下的一个 Repository，至少需要两部分

- Namespace(名称空间) # 有的也称为 ProjectID。
  - Docker 将用户注册的账户名称作为 Namespace，若 Namespace 被省略，则就是这个 Registry 官方的。所以也可以这么理解。
- Repository(仓库) # 仓库名称

很多时候都将 Namespace 和 Repository 合起来，统一称为 Repository

仓库分为公开仓库(Public)和私有仓库(Private)两种形式。当用户创建了自己的镜像之后就可以使用 push 命令将它上传到公有或者私有仓库，这样下载在另外一台机器上使用这个镜像时候，只需需要从仓库上 pull 下来就可以了。

注意：Docker 仓库的概念跟 Git 类似，Registry 可以理解为 GitHub 这样的托管服务。

### Tag(标签)

Repository 可以存放不同的 Image(比如 nginx,redis,centos 等)，通过 Tag 来区分这些 Image。说白了，Tag 就是 Image 的名称。

## Container(容器)

容器是从镜像创建的运行实例。它可以被启动、开始、停止、删除。每个容器都是相互隔离的，保证安全的平台。可以把容器看做是一个简易版的 Linux 环境（包括 root 用户权限、进程空间、用户空间和网络空间等）和运行在其中的应用程序。注意：镜像是只读的，容器在启动的时候创建一层可写层作为最上层。

Image 与 Container 的关系，就好比是程序与进程之间的关系。Image 类似程序是静态的。Container 类似进程是动态的，是有生命周期的。

## 联合文件系统

- 当我们在下载镜像的时候，会发现每一层都有一个 id，这是 **Layer(层)** 的概念，是 **UnionFS(联合文件系统)** 中的重要概念
- 联合文件系统（UnionFS）是一种分层、轻量级并且高性能的文件系统，它支持对文件系统的修改作为一次提交来一层层的叠加，同时可以将不同目录挂载到同一个虚拟文件系统下
- 联合文件系统是 Docker 镜像的基础。镜像可以通过分层来进行继承，基于基础镜像（没有父镜像），可以制作各种具体的应用镜像。
- 不同容器就可以共享一些基础的文件系统层，同时再加上自己独有的改动层，大大提高了存储的效率。

# Rootless Containers

> 参考:
>
> - [GitHub 项目，rootless-containers](https://github.com/rootless-containers)

**Rootless Containers(无根容器)** 是指非特权用户能够创建、运行和以各种方式管理容器。这个术语还包括围绕容器的各种工具，这些工具也可以作为非特权用户运行。

运行 Rootless Containers 通常需要弃用 CGroupV2 来限制 CPU、内存、I/O、PID 这些资源的消耗。
