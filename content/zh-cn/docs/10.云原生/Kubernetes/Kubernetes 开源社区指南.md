---
title: Kubernetes 开源社区指南
weight: 100
---

# 概述

> 参考：

[如何玩转 Kubernetes 开源社区？这篇文章一定要看！](https://mp.weixin.qq.com/s/aZGBkBpFZOEyoa1xj-16kQ)

近日，「DaoCloud 道客」成功进入 Kubernetes 开源榜单累计贡献度全球前十，亚洲前三。基于在 Kuberntes 开源社区的长期深耕细作，「DaoCloud 道客」积累了一些心得，特写此文章，旨在帮助对开源贡献感兴趣的同学快速⼊⻔，并为之后的进阶之路提供⼀些参考和指导意义。

这⼀章节，你将了解整个 Kubernetes 社区是如何治理的：

## 1.1. 分布式协作

与公司内部集中式的项⽬开发模式不同，⼏乎所有的开源社区都是⼀个分布式、松散的组织，为此  ，Kubernetes 建⽴了⼀套完备的社区治理制度。协作上，社区⼤多数的讨论和交流主要围绕 issue 和 PR 展开。由于 Kubernetes ⽣态⼗分繁荣，因此所有对 Kubernetes 的修改都⼗分谨慎，每个提交的 PR 都需要通过两个以上成员的 Review 以及经过⼏千个单元测试、集成测试、端到端测试以及扩展性测试，所有这些举措共同保证了项⽬的稳定。

## 1.2. Committees

委员会由多人组成，主要负责制定组织的行为规范和章程，处理一些敏感的话题。常见的委员会包括行为准则委员会，安全委员会，指导委员会。

## 1.3. SIG

SIG 的全称是 Special Interest Group，即特别兴趣⼩组，它们是 Kubernetes 社区中关注特定模块的永久组织，Kubernetes 作为⼀个拥有⼏⼗万⾏源代码的项⽬，单⼀的⼩组是⽆法了解其实现的全貌的。Kubernetes ⽬前包含 20 多个 SIG，它们分别负责了 Kubernetes 项⽬中的不同模块，这是我们参与 Kubernetes 社区时关注最多的⼩组。作为刚刚参与社区的开发者，可以选择从某个 SIG 入手，逐步了解社区的⼯作流程。

## 1.4. KEP

KEP 的全称是 Kubernetes Enhancement Proposal，因为 Kubernetes ⽬前已经是⽐较成熟的项⽬了，所有的变更都会影响下游的使⽤者，因此，对于功能和  API 的修改都需要先在 kubernetes/enhancements 仓库对应 SIG 的⽬录下提交提案才能实施，所有的提案都必须经过讨论、通过社区 SIG Leader 的批准。

## 1.5. Working Group

这是由社区贡献者⾃由组织的兴趣⼩组，对现阶段的⼀些⽅案和社区未来发展⽅向进⾏讨论，并且会周期性的举⾏会议。会议⼤家都可以参加，⼤多是在国内的午夜时分。以 scheduling 为例，你可以查看文档 Kubernetes Scheduling Interest Group 了解例次会议纪要。会议使⽤ Zoom 进⾏录制并且会上传到 Youtube, 过程中会有主持⼈主持，如果你是新⼈，可能需要你进行自我介绍。

🔗<https://docs.google.com/document/d/13mwye7nvrmV11q9_Eg77z-1w3X7Q1GTbslpml4J7F3A/edit%23heading%25253Dh.ukbaidczvy3r>

## 1.6. MemberShip

| 角色       | 职责                                         | 要求                                             |
| ---------- | -------------------------------------------- | ------------------------------------------------ |
| Member     | 社区积极贡献者                               | 对社区有多次贡献并得到两名 reviewer 的赞同       |
| Reviewer   | 对其他成员贡献的代码积极的 review            | 在某个子项目中长期 review 和贡献代码             |
| Approver   | 对提交的代码进行最后的把关，有合并代码的权限 | 属于某一个子项目经验丰富的 reviewer 和代码贡献者 |
| Maintainer | 制定项目的优先级并引领项目发展方向           | 证明自己在这个项目中有很强的责任感和技术能力     |

每种⻆⾊承担不同的职责，同时也拥有不同的权限。⻆⾊晋升主要参考你对社区的贡献，具体内容可参考 KubernetesMemberShip。

🔗<https://github.com/kubernetes/community/blob/master/community-membership.md>

## 1.7. Issue 分类

🔗 文章链接：<https://hackmd.io/O_gw_sXGRLC_F0cNr3Ev1Q>

## 1.8. 其他关注项

更多详情可参见官⽅⽂档，⽂档详细描述了该如何提交 PR，以及应该遵循什么样的原则，并给到了⼀些最佳实践。

🔗 官方文档：<https://www.kubernetes.dev/docs/guide/contributing/>

2.1.  申请  CLA

当你提交 PR 时，Kubernetes 代码仓库 CI 流程会检查是否有 CLA 证书，如何申请证书可以参考官⽅⽂档。

🔗 官方文档：<https://github.com/kubernetes/community/blob/master/CLA.md>

2.2.  搜索 first-good-issue 「你可以选择你感兴趣的或你所熟悉的 SIG」

first-good-issue 是 Kubernetes 社区为培养新参与社区贡献的开发⼈员⽽准备的 issue，⽐较容易上⼿。

以 sig/scheduling 为例，在 Issues  中输⼊：

```swift
is:issue is:open label:sig/scheduling label:"good first issue" no:assignee
```

该 Filters 表示筛选出没有被关闭的，属于 sig/scheduling，没有 assign 给别⼈的 good first issue。

如果没有相关的 good-first-issue，你也可以选择 kind/documentation 或者 kind/cleanup 类型 issue。

🔗Command-Hlep：<https://prow.k8s.io/command-help?repo=kubernetes%25252Fkubernetes>

| /retitle               | 重命名标题                           |
| ---------------------- | ------------------------------------ |
| /close                 | 关闭 issue                           |
| /assign                | 将 issue assign 给⾃⼰               |
| /sig scheduling        | 添加标签 sig/scheduling              |
| /remove-sig scheduling | 去掉标签                             |
| /help                  | 表示需要帮助，会打上标签 help wanted |
| /good-first-issue      | 添加标签 good first issue            |
| /retest                | 重新测试出错的测试⽤例               |
| /ok-to-test            | 准备好开始测试                       |

a.  Fork 代码仓库

将 kubernetes/Kubernetes fork 到⾃⼰的 GitHub 账号名下。

b. Clone ⾃⼰的代码仓库

```nginx
git clone git@github.com:<your github id>/kubernetes.git
```

c.   追踪源代码仓库代码变动

- 添加 upstream：

```cs
git remote add upstream https://github.com/kubernetes/kubernetes.git
```

- git remote -v 检查是否添加成功，成功则显示：

```properties
origin git@github.com:<your github id>/kubernetes.git
(fetch)
origin git@github.com:<your github id>/kubernetes.git
(push)
upstream
https://github.com/kubernetes/kubernetes.git (fetch)
upstream
https://github.com/kubernetes/kubernetes.git (push)
```

- 同步  upstream kubernetes 最新代码

```properties
git checkout master
git pull upstream master git push
```

d.   切分支，编码

```xml
git checkout -b <branch name>
```

e.  Commit，并提交 PR

- 命令行：

```nginx
git commit -s -m '<change me>'
```

- 注意：
- commit push  前先执⾏一些检查，如  make update 等
- 如果本次修改还没有完成，可以使⽤  Github Draft   模式，并添加  \[WIP]  在标题中
- Commit  信息过多，且没有特别⼤的价值，建议合成⼀条  commit  信息

```nginx
git rebase -i HEAD~2
```

f.   提交  PR

在  GitHub  ⻚⾯按照模版提交  PR

a. PR 提交后需要执⾏ Kubernetes CI 流程，此时需要 Kubernetes Member 输入  /ok- to-test 命令，然后会⾃动执⾏ CI，包括验证和各种测试。可以 @ 社区成员帮助打标签。

b. ⼀旦测试失败，修复后可以执⾏  /retest 重新执⾏失败的测试，此时，你已经可以⾃⼰操作。

a. 每次提交需要有 2 个 Reviewer 进⾏ Code Review， 如果通过，他们会打上 /lgtm 标签，表示 looks good to me, 代码审核完成。

b.   另外需要⼀个 Approver 打上 /approve 标签，表示代码可以合⼊主⼲分⽀，GitHub 机器⼈会⾃动执⾏ merge 操作。

c.  PR 跟进没有想像中那么快，有时候 1-2 周也正常。

d.   恭喜，你完成了第⼀个 PR 的提交。

```cs
Sponsored by 2 reviewers and multiple contributions to the project
```

```css
PR
Issue
Kep
Comment
```

a.   多参与社区的讨论，表达⾃⼰的观点

b.   多参与  issue 的解答，帮助提问者解决问题，社区的意义也在于此

c.   可以在看源代码的时候多留意⼀些语法、命名和重复代码问题，做⼀些重构相关的⼯作

d.   从测试⼊⼿也是⼀个好办法，如补全测试，或者修复测试

e.   参与⼀些代码的 review，可以学到不少知识

f. 最有价值的肯定是 feature 的实现，可以提交 kep

参考  Makefile  ⽂件

🔗<https://github.com/kubernetes/kubernetes/blob/master/build/root/Makefile>

- 单元测试（单方法）

```go
go test -v --timeout 30s k8s.io/kubectl/pkg/cmd/get -run
^TestGetSortedObjects$
```

- 集成测试（单⽅法）

```bash
make test-integration WHAT=./vendor/k8s.io/kubectl/pkg/cmd/get GOFLAGS="-v" KUBE_TEST_ARGS="-run ^TestRuntimeSortLess$"
```

- E2E 测试

可以使⽤ GitHub 集成的 E2E 测试：

```bash
/test pull-kubernetes-node-kubelet-serial
```

Kubernetes 和 Linux 一样, 早已成为 IT 技术的重要事实标准，而开源 Kubernetes 是整个行业的  “上游”，灌溉了数亿的互联网和企业应用。「DaoCloud 道客」融合自身在各行各业的实战经验，持续贡献 Kubernetes 开源项目，致力于让以 Kubernetes 为代表的云原生技术更平稳、高效地落地到产品和生产实践中。此外，「DaoCloud 道客」全面布局开源生态，是最早一批加入 CNCF 基金会的云原生企业，拥有云原生基金会成员，Linux 基金会成员，Linux 基金会培训合作伙伴，Kubernetes 培训合作伙伴，Kubernetes 兼容性认证，以及 Kubernetes 认证服务提供商等资质，坚持构建并维护云原生开源生态圈。在开源贡献这条道路上，「DaoCloud 道客」会一直走下去，也愿意成为开源社区的守护者、暸望塔，并始终坚信开源的力量、原生的力量会与这个时代产生共鸣，迸发出属于自己的光彩。

作者简介

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/bb6e356a-0809-470a-8674-9af62b340848/640)

殷纳

DaoCloud  云原生工程师

Kubernetes Member

专注 Kubernetes 及多集群管理开源工作

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/bb6e356a-0809-470a-8674-9af62b340848/640)

**DaoCloud 公司简介**

上海道客网络科技有限公司，成立于 2014 年底，公司拥有自主知识产权的核心技术，以云计算、人工智能为底座构建数字化操作系统为实体经济赋能，推动传统企业完成数字化转型。成立迄今，公司已在金融科技、先进制造、智能汽车、零售网点、城市大脑等多个领域深耕，标杆客户包括交通银行、浦发银行、上汽集团、东风汽车、海尔集团、金拱门（麦当劳）等。被誉为科技领域准独角兽企业。公司在北京、武汉、深圳、成都设立多家分公司及合资公司，总员工人数超过 300 人，是上海市高新技术企业、上海市 “专精特新” 企业和上海市 “科技小巨人” 企业，并入选了科创板培育企业名单。

网址：www.daocloud.io

邮件：info@daocloud.io

电话：400 002 6898

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/bb6e356a-0809-470a-8674-9af62b340848/640)

文章转载自道客船长。[点击这里阅读原文了解更多](https://mp.weixin.qq.com/s?__biz=MzA5NTUxNzE4MQ==&mid=2659272563&idx=1&sn=9cbdc17729dc631d490ab33897012c73&scene=21#wechat_redirect)。

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/bb6e356a-0809-470a-8674-9af62b340848/640)

中国 KubeCon + CloudNativeCon + Open Source Summit 虚拟大会

12 月 9 日至 10 日

<https://www.lfasiallc.com/kubecon-cloudnativecon-open-source-summit-china/>

**标准票和免费的 “主题演讲 + 解决方案展示仅用票”**

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/bb6e356a-0809-470a-8674-9af62b340848/640)

选定合适门票，马上扫码注册！

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/bb6e356a-0809-470a-8674-9af62b340848/640)

（[https://www.bagevent.com/event/7680821）](https://www.bagevent.com/event/7680821%EF%BC%89)

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/bb6e356a-0809-470a-8674-9af62b340848/640)

CNCF 概况（幻灯片）

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/bb6e356a-0809-470a-8674-9af62b340848/640)

扫描二维码联系我们！

---

**\_CNCF (Cloud Native Computing Foundation) 成立于 2015 年 12 月，隶属于 Linux Foundation，是非营利性组织。 \_**

**\_\***CNCF\***\*\*\*\***（\***\*\*\*\***云原生计算基金会**\*）致力于培育和维护一个厂商中立的开源生态系统，来推广云原生技术。我们通过将最前沿的模式民主化，让这些创新为大众所用。请长按以下二维码进行关注。\_**

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/bb6e356a-0809-470a-8674-9af62b340848/640)
