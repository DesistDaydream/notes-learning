---
title: DevOps
linkTitle: DevOps
date: 2024-05-11T13:38
weight: 20
---

# 概述

> 参考：
>
> -

**Developer(开发者)**

**Operator(操作者(运维))**

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/mcelwu/1616077542789-b79d4008-5e43-4380-a6a5-ab0c2b95cdd2.jpeg)

敏捷开发（Agile Development）、持续集成（Continuous Integration）、持续交付（Continuous Delivery）、开发运维一体化（Development Operations），所覆盖的软件生命周期的阶段不同。

DevOps 是一种文化，是一组过程、方法与系统的统称，用于促进开发（应用程序/软件工程）、技术运营和质量保障（QA）部门之间的沟通、协作与整合。

它是一种重视“软件开发人员（Dev）”和“IT 运维技术人员（Ops）”之间沟通合作的文化、运动或惯例。透过自动化“软件交付”和“架构变更”的流程，来使得构建、测试、发布软件能够更加地快捷、频繁和可靠。

它的出现是由于软件行业日益清晰地认识到：为了按时交付软件产品和服务，开发和运营工作必须紧密合作。

计划-开发-构建-测试-修复-测试-交付运维-部署

- CI：持续集成
  - 当程序员开发完提交代码到类似 github 的地方上之后，有一款工具可以自动拉去代码进行构建，当出现问题时，会报告给程序员，重复执行“构建-测试”的过程
- CD：持续交付 Delivery
  - 当有一款工具自动把测试好的东西(比如一个容器用的镜像)交付给运维，自动执行“构建-交付”的过程
- CD：持续部署 Deployment
  - 当有一款工具，运维都不用了，可以直接部署的时候

# GitOps

> 参考：
>
> - [公众号-云原声实验室，大妈都能看懂的 GitOps 入门指南](https://mp.weixin.qq.com/s/JkZP9X2g9TOj6QkrbCDRQQ)

GitOps 这个概念最早是由 Kubernetes 管理公司 Weaveworks 公司在 2017 年提出的，如今已经过去了 5 个年头，想必大家对这个概念早有耳闻，但你可能并不知道它到底是什么，它和 DevOps 到底是啥关系，本文就来帮大家一一解惑。

GitOps = IaC + Git + CI/CD，即基于 IaC 的版本化 CI/CD。它的核心是使用 Git 仓库来管理基础设施和应用的配置，并且**以 Git 仓库作为基础设施和应用的单一事实来源**，你从其他地方修改配置（比如手动改线上配置）一概不予通过。

Git 仓库中的声明式配置描述了目标环境当前所需基础设施的期望状态，借助于 GitOps，如果集群的实际状态与 Git 仓库中定义的期望状态不匹配，Kubernetes reconcilers 会根据期望状态来调整当前的状态，最终使实际状态符合期望状态。
另一方面，现代应用的开发更多关注的是迭代速度和规模，拥有成熟 DevOps 文化的组织每天可以将代码部署到生成环境中数百次，DevOps 团队可以通过版本控制、代码审查以及自动测试和部署的 CI/CD 流水线等最佳实践来实现这一目标，这就是 GitOps 干的事情。

# CI/CD 工具

- [Jenkins](/docs/10.云原生/DevOps/Jenkins/Jenkins.md)
- Tekton # 云原生推荐
- Argocd # 云原生推荐
- [Drone](/docs/10.云原生/DevOps/Drone/Drone.md)
- [Skaffold](/docs/10.云原生/DevOps/Skaffold/Skaffold.md)(k8s 专属)
- [GitLab](/docs/2.编程/Programming%20tools/SCM/GitLab/GitLab.md)-ci
- Dagger # Docker 创始人 Solomon Hykes 推出的产品
    - [公众号，k8s 技术圈-Docker 创始人的新产品 Dagger 好用吗？](https://mp.weixin.qq.com/s/4hwtgV6WJ-60FL1lGHoAQw)
    - [公众号，k8s 技术圈-数据约束语言 CUE 是何方神圣？](https://mp.weixin.qq.com/s/J2Hid1dO8ebkWL5UrVBeyA)
- Zadig # [GitHub 项目，koderover/zadig](https://github.com/koderover/zadig)
    - 一款面向开发者设计的云原生持续交付(Continuous Delivery)产品，具备高可用 CI/CD 能力，提供云原生运行环境，支持开发者本地联调、微服务并行构建和部署、集成测试等。Zadig 不改变现有流程，无缝集成 Github/Gitlab、Jenkins、多家云厂商等，运维成本极低。我们的目标是通过云原生技术的运用和工程产品赋能，打造极致、高效、愉悦的开发者工作体验，让工程师成为企业创新的核心引擎。
