---
title: Opensatck 介绍
weight: 1
---

# 概述

> 参考
>
> - [GitHub 项目，openstack/openstack](https://github.com/openstack/openstack)
> - [官网](https://www.openstack.org/)
> - 原文：[OpenStack 关键技术系列: 最全 OpenStack 知识科普](https://cloud.tencent.com/developer/article/1395617)

Ubuntu 成为 OpenStack 部署排名第一的操作系统

https://cn.ubuntu.com/blog/ubuntu-becomes-number-one-os-for-openstack-deployment-post-1

# OpenStack 关键技术系列: 最全 OpenStack 知识科普

最近几年，OpenStack 这个词大家早都熟的不能再熟，越来越多人开始关注。

对于大部分人来说，这还是一个很陌生的词，不知道它到底是什么，从哪里来，有什么用，和自己的工作有什么关系。


有人可能知道，它和现在非常火的云计算有很大的关系。伴随它一起出现的，还有很多新词，例如NFV、Nova、Neutron、Horizon等，更加让人云里雾里。


为了消除大家的疑惑，今天我们就来一个“大揭秘”——通过这篇通俗易懂的科普文，帮助大家轻松入门OpenStack。

OpenStack 的起源

这玩意到底是从哪冒出来的？

我们先来说说OpenStack的起源吧。

2002年，美国著名的电商公司亚马逊（Amazon）干了一件“不务正业”的事。他们向客户推出了一项全新的业务——包括存储空间、计算能力等资源服务的Web Service。这就是大名鼎鼎的 AWS（Amazon Web Service）。

说白了，这个Web Service服务，就是为大家提供“远程电脑”。你可以远程控制它，有硬盘，有CPU，有内存啥的。你在上面配置你的各种服务，然后给你的用户使用，例如网站、FTP等。这个就是云计算的一种早期形式。


后来，到了2006年，亚马逊又推出了弹性计算云（Elastic Compute Cloud），也称 EC2 。EC2配置界面更简单，使用起来更方便，关键一点，它开始有了“弹性”！

什么是“弹性”？别急哈，等会我们再解释。

同样是 2006 年，8 月 9 日，Google 首席执行官埃里克·施密特在搜索引擎大会上首次提出“云计算”（Cloud Computing）的概念。从此，云计算进入了高速发展阶段。

云计算

到了 2010 年，当时有一家名叫 Rackspace 的公司，他们一直在做和亚马逊一样的云主机和云储存服务，但是始终都干不过亚马逊，排名第二。他们一气之下，干脆就把它们的云储存服务给开源了。

啥叫开源（Open Source）？开源就是开放源代码，把程序的代码公开了，给所有人免费查看和使用。和他们一起开放源代码的，还有一个家伙，就是——NASA。

好吧，又是一个“不务正业”的家伙。

NASA 之前在云计算方面投入了大量的资金，但是后来发现这玩意好像是个无底洞，太烧钱了。而且，他们也似乎意识到这不是他们该干的事。所以，NASA 和 Rackspace 一起，选择开放源代码。

其实还有一个原因：以前 NASA 是使用 Eucalyptus 云计算管理平台，不过这个平台分成两个版本，一个开源的版本，一个收费的版本。这就导致 NASA 很不爽，向 Eucalyptus 贡献代码，结果 Eucalyptus 认为这个代码和收费版本冲突，不接受。NASA 给气得不行，所以选择了将代码开源。

Rackspace 和 NASA 并不是简单地代码一丢完事，而是联手共同成立了一个开源项目。这个项目，就是 OpenStack。

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/orsf4g/1616123375003-fa1189ba-3087-41b2-b6b9-a86a312ce567.jpeg)

OpenStack 的版本

开源后的 OpenStack，到底经历了什么？

开源项目的玩法，和企业内部研发是完全不一样的。开源项目中，地球上所有人都可以为这个项目贡献自己的力量，也可以使用这个项目的开发成果。也就是说，“人人为我，我为人人”。

## 开源(Open Source)

但是，为了保证项目能规范、有序地推进下去，还是需要有人“牵头”和“打杂”的。OpenStack 作为一个开源项目，它是由开源社区来负责推进和维护的。这个社区也并不是一盘散沙，它有自己的组织形态。

首先，有一个 OpenStack 基金会，下面设立了董事会、技术委员会、用户委员会。基金会享有话语权，进行目标和发展的引导。

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/orsf4g/1616123375003-25081c7e-1339-496f-8c52-5c0dcc647744.jpeg)

基金会成员有三种形式。首先是独立个体，也就是以个人名义为 OpenStack 做出贡献。

其次是铂金会员。主要由对 OpenStack 作出重要承诺的公司组成，他们提供资金与资源。目前，OpenStack 基金会主要有 7 家铂金会员。

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/orsf4g/1616123374992-d2cd2b2d-554f-4427-a861-f6992500f61c.jpeg)

最后是金牌会员。同样由公司组成，他们赞助的资金与资源比铂金会员稍微少一些。目前，OpenStack基金会拥有 21 位金牌会员。

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/orsf4g/1616123375063-adb07ccb-a240-40bb-9de3-4ffcb67c5d10.jpeg)

从2010年项目诞生之日起，OpenStack开源社区每年都会开两次设计峰会（Design Summit），发布两个正式版本。迄今为止，一共已经出了17个版本。

## OpenStack 设计峰会

这里我要开启“吐槽”模式了。开源社区这帮搞技术的宅男腐女，不管年龄大小，内心仍然是一群孩子。他们平时在公司上班比较“木鸡”，在社区这种自由环境里是一个比一个“皮”。

从哪可以看出来？就在“取名”上——他们竟然给每个版本都单独取了一个名字(而非商业软件一样按数字序号命名)。每个版本的名字如下:

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/orsf4g/1616123375008-6bc9c88c-ad64-447b-a338-578ddf116de1.jpeg)

不知道大家看出来没有，这些名字都是有“玄机”的！首先，版本号的第一个字母，从 A 开始，然后 B、C、D… 其次，每个名字都是从当次设计峰会所在城市中选一个地名，作为该版本的名字。

例如，第一个版本 Austin，就是根据 Rackspace 公司所在地（也是第一次峰会所在地）——美国德克萨斯州的首府“奥斯丁”确定的。还有第 9 个版本，当时峰会是在香港举办的，用的“Ice House Street(雪厂街)”这个名字。

这么做的直接后果就是，记忆和分辨起来真的很困难，容易看晕。

## OpenStack 的架构

它由哪些部分组成？是如何进行工作的？

接下来，我们看看 OpenStack 的架构。前面说了，OpenStack 从一开始，就是为了云计算服务的。简单来说，它就是一个操作系统，一套软件，一套 IaaS 软件。

什么是 IaaS？Infrastructure as a Service，基础设施即服务(了解更多，看这里:“云计算”)。

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/orsf4g/1616123375043-5e534c13-ef20-44ab-8faf-4c7fc3c42ee0.jpeg)

云计算的三种服务模式：IaaS、PaaS、SaaS

管理“基础设施资源”，便于用户调用和使用，是 OpenStack 的首要任务。基础设施资源，主要包括三个方面：计算、存储、网络。说通俗点，就是 CPU，硬盘，网卡。

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/orsf4g/1616123375054-bc9b5930-658c-4a2e-b064-5921ea5ec792.jpeg)

OpenStack 对资源进行管理，并且以服务的形式提供给上层应用或者用户去使用。

例如前面我们所说的“弹性”。正是因为资源能够被灵活调用，所以用户使用资源时，这个云平台可以根据用户的需要，动态增加和删减资源，不用中断用户的使用，更无需全新申请。这就是“弹性”。

那么，它到底是如何实现的呢？答案是——通过它的众多组件。

前方高能预警……

学习 OpenStack，最痛苦的事情，莫过于看它的架构。不信？好，扔个图给你看:

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/orsf4g/1616123375069-000ceee3-f4f8-4383-b38e-eccb7807d276.jpeg)

OpenStack 系统架构逻辑关系图

吓尿了吧。这还不算是最复杂的，再扔一个给你。

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/orsf4g/1616123375077-819b8d49-5aa3-467f-b8f1-9c0d217fee49.jpeg)

好了好了，不扔了，人都跑光了。OpenStack 拥有众多的组件，通过组件之间协同进行工作，所以看上去架构非常复杂。我还是用一个简单的图吧，看得更明白些，如下:

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/orsf4g/1616123375065-6de47154-6b9a-433f-8f0a-54891864e52b.jpeg)

这个图里面的彩色方块，就是 OpenStack 最核心的组件。说到这些组件的名字，我实在忍不住又要吐槽这帮程序猿了，简直就是“取名狂魔”！他们不仅给每个项目版本单独取名字，连 openstack 内部的组件也难逃他们的魔爪。

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/orsf4g/1616123375093-69331370-989f-4eb8-aa0d-81a432cf3ae3.jpeg)

## OpenStack 关键组件及作用

这些组件里，我挑几个再介绍一下(看不懂也没关系，可以跳过):

**Nova**

Nova 是整个 Openstack 里面最核心的组件。当初 Rackspace 和 NASA 贡献代码时，NASA 贡献的那部分就是 Nova 最早的代码（Rackspace 贡献的代码是 Swift）。OpenStack 云实例生命期所需的各种动作都将由 Nova 进行处理和支撑，它负责管理整个云的计算资源、网络、授权及测度。

**Keystone**

Keystone 为所有的 OpenStack 组件提供认证和访问策略服务，主要对(但不限于)Swift、Glance、Nova 等进行认证与授权。

**Horizon**

Horizon 是一个用以管理、控制 OpenStack 服务的 Web 控制面板。用户可以通过这个界面对 OpenStack 状态进行查看和管理。

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/orsf4g/1616123375102-1e3c1e11-e5bf-4416-ab03-f1e0cb9d7f76.jpeg)

用 Horizon 管理 OpenStack

也就是说，OpenStack 的组件都有自己的功能定位。其实，每个组件都可以算是独立的一个程序(Software)。

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/orsf4g/1616123375100-d253c5a7-51fc-467d-9240-b6626776e530.jpeg)

Open 为开放之意，Stack 则是堆砌

也就是许多 Open 的 Softwares 进行集合和堆砌。

关于技术细节，就先说这么多吧，再说下去估计人都跑光啦。

## OpenStack 的发展

现在的它，是一个什么规模和状态？

经过八年的努力，如今的 OpenStack 已经今非昔比。很多企业和个人纷纷加入 Openstack 开源社区，使之成为了目前仅次于 LINUX 的全球第二大开源社区。

按官网最新数据，现在有 180 多个国家，677 家企业，87426 名社区会员通过各种方式支撑着这个项目。项目的代码也已经超过了 2000 万行。

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/orsf4g/1616123375099-586ae5b0-692c-4044-9c7f-04ce0b3e9456.jpeg)

全球一半以上的 500 强企业，都采用了 OpenStack 技术。而且，根据调查，有 75%以上的企业打算今后使用这项技术。

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/orsf4g/1616123375150-1db2a541-efd2-4249-9401-321411d8edb2.jpeg)

OpenStack 在各行业的应用情况占比（2017 年）

小枣君作为一枚通信汪，这里要特别强调一下，虽然 OpenStack 是云计算技术，主要是 IT 的概念，但对于通信行业来说极为重要。

通信网络中的核心网，已经全面开始了向虚拟化和云计算的演进。小枣君之前就介绍过，现在通信行业里火热的 NFV 技术，就是基于虚拟化的，采用了 IT 里面的很多理念和设计。而核心网的 IT 化，将是整个通信系统 IT 化的第一步。

华为的 FusionSphere 平台和中兴的 TECS 平台，都是基于 OpenStack 进行二次开发的商业系统。这些平台都已经被自家的核心网和云计算产品采用，目前处于替代传统平台的阶段。

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/orsf4g/1616123375097-c7b08789-dcb4-4103-a48d-7392e1f7e9ae.jpeg)

OpenStack 之所以这么受欢迎，主要原因有三个方面:

- 首先是快速。OpenStack 安装部署所需要的时间很少，而时间就是价值。
- 其次是灵活。OpenStack 获得了各大领导厂商的广泛支持，兼容性和适用性极强，使用起来非常方便可靠。
- 最后是便宜。作为开源项目，OpenStack 的使用成本相对低廉，还能获得源源不断的更新，因为开源社区在为项目贡献活力。

总而言之，Openstack 拥有非常大的发展潜力，目前处于高速发展的上升期。在未来很长一段时间内，这种趋势都不会改变。

## OpenStack 的学习

到底该如何对它进行学习呢？

经过上面的介绍，想必大家热血沸腾，跃跃欲试了吧？OpenStack 这么牛掰，到底该如何学习它呢？它看上去那么复杂，会不会很难学会呢？

其实，虽然前面看到的架构很复杂，但是真心要学习 OpenStack 的话，并没有想象得那么困难。

因为 OpenStack 是开源的项目，所以互联网上相关的学习资料非常丰富。无论是官方文档，还是非官方资料，都数不胜数。所以，问题不在于资料缺乏，而是资料太多你看不完…

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/orsf4g/1616123375164-e57c02e6-7476-44d5-96d9-5b22e66d8f8f.jpeg)

官方网站强大的资料库和技术支撑

<https://www.openstack.org/>

网上也有很多手把手进行教学的文档和视频，可以方便新人学习时进行参考。推荐几个大咖，大家可以百度找他们的博客来看: 陈沙克、何明桂、孔令贤，Cloudman。有了官方资料，加上大咖的博客，你只需要一台电脑，你就可以开始 OpenStack 的学习——直接下载，直接安装，直接配置，直接使用，没有任何门槛要求。如果遇到问题，先别急着找人问，先自己尝试找资料解决，一定会学得嗖嗖快。

不过，OpenStack 入门虽然很容易，但是精通就很难了。需要长时间不断地钻研和积累，还需要进行大量的实践部署，才有可能成为专家。

到底哪些人需要学习 OpenStack 呢？有三种人最应该立刻开始对它的学习。

1、IT 行业从业者

这就不用多说了，未来网络就是云计算，大数据的天下，只要是从事 IT 方面的工作，肯定会和云打交道，OpenStack 作为云计算技术的代表，是一个合适的切入点。

2、通信、电子、计算机专业的大学生

云计算技术在目前大部分高校都没有合适的教学规划，所以，在校大学生应该注意提前进行此类趋势技术的学习，既有利于就业，又能紧跟时代节奏，选择将来进修的合适方向。

3、通信行业从业者

啥都别说了，通信人赶紧滚去学习吧。好啦，关于 OpenStack 的介绍，就到这里，谢谢大家的观看！
