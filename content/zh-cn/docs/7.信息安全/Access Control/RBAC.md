---
title: RBAC
linkTitle: RBAC
weight: 20
---

# 概述

> 参考：
>
> - [Wiki, RBAC](https://en.wikipedia.org/wiki/Role-based_access_control)

基于角色的权限访问控制（Role-Based Access Control）作为传统访问控制（自主访问，强制访问）的有前景的代替受到广泛的关注。在 RBAC 中，权限与角色相关联，用户通过成为适当角色的成员而得到这些角色的权限。这就极大地简化了权限的管理。在一个组织中，角色是为了完成各种工作而创造，用户则依据它的责任和资格来被指派相应的角色，用户可以很容易地从一个角色被指派到另一个角色。角色可依新的需求和系统的合并而赋予新的权限，而权限也可根据需要而从某角色中回收。角色与角色的关系可以建立起来以囊括更广泛的客观情况。

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/wn3hwi/1616125478657-b931db83-6f72-44d0-9371-0e19ae04ee25.jpeg)

RBAC 支持三个著名的安全原则：最小权限原则，责任分离原则和数据抽象原则。

1. 最小权限原则之所以被 RBAC 所支持，是因为 RBAC 可以将其角色配置成其完成任务所需要的最小的权限集。
2. 责任分离原则可以通过调用相互独立互斥的角色来共同完成敏感的任务而体现，比如要求一个计帐员和财务管理员共参与同一过帐。
3. 数据抽象可以通过权限的抽象来体现，如财务操作用借款、存款等抽象权限，而不用操作系统提供的典型的读、写、执行权限。然而这些原则必须通过 RBAC 各部件的详细配置才能得以体现。

RBAC 有许多部件(BUCU)，这使得 RBAC 的管理多面化。尤其是，我们要分割这些问题来讨论：用户与角色的指派；角色与权限的指派；为定义角色的继承进行的角色与角色的指派。这些活动都要求把用户和权限联系起来。然而在很多情况下它们最好由不同的管理员或管理角色来做。对角色指派权限是典型的应用管理者的职责。银行应用中，把借款、存款操作权限指派给出纳角色，把批准贷款操作权限指派给经理角色。而将具体人员指派给相应的出纳角色和管理者角色是人事管理的范畴。角色与角色的指派包含用户与角色的指派、角色与权限的指派的一些特点。更一般来说，角色与角色的关系体现了更广泛的策略。

## 基本概念

RBAC 认为权限授权实际上是 Who、What、How 的问题。在 RBAC 模型中，who、what、how 构成了访问权限三元组,也就是“Who 对 What(Which)进行 How 的操作”。

Who：权限的拥用者或主体（如 Principal、User、Group、Role、Actor 等等）。

What：权限针对的对象或资源（Resource、Class）。

How：具体的权限（Privilege,正向授权与负向授权）。

Operator：操作。表明对 What 的 How 操作。也就是 Privilege+Resource

Role：角色，一定数量的权限的集合。权限分配的单位与载体,目的是隔离 User 与 Privilege 的逻辑关系.

Group：用户组，权限分配的单位与载体。权限不考虑分配给特定的用户而给组。组可以包括组(以实现权限的继承)，也可以包含用户，组内用户继承组的权限。User 与 Group 是多对多的关系。Group 可以层次化，以满足不同层级权限控制的要求。

RBAC 的关注点在于 Role 和 User, Permission（允许/权限）的关系。称为 User assignment(UA)和 Permission assignment(PA).关系的左右两边都是 Many-to-Many 关系。就是 user 可以有多个 role，role 可以包括多个 user。

凡是用过 RDBMS 都知道，n:m 的关系需要一个中间表来保存两个表的关系。这 UA 和 PA 就相当于中间表。事实上，整个 RBAC 都是基于关系模型。

Session 在 RBAC 中是比较隐晦的一个元素。标准上说：每个 Session 是一个映射，一个用户到多个 role 的映射。当一个用户激活他所有角色的一个子集的时候，建立一个 session。每个 Session 和单个的 user 关联，并且每个 User 可以关联到一或多个 Session.

在 RBAC 系统中，User 实际上是在扮演角色(Role)，可以用 Actor 来取代 User，这个想法来自于 Business Modeling With UML 一书 Actor-Role 模式。考虑到多人可以有相同权限，RBAC 引入了 Group 的概念。Group 同样也看作是 Actor。而 User 的概念就具象到一个人。

这里的 Group 和 GBAC（Group-Based Access Control）中的 Group（组）不同。GBAC 多用于操作系统中。其中的 Group 直接和权限相关联，实际上 RBAC 也借鉴了一些 GBAC 的概念。

Group 和 User 都和组织机构有关，但不是组织机构。二者在概念上是不同的。组织机构是物理存在的公司结构的抽象模型，包括部门，人，职位等等，而权限模型是对抽象概念描述。组织结构一般用 Martin fowler 的 Party 或责任模式来建模。

Party 模式中的 Person 和 User 的关系，是每个 Person 可以对应到一个 User，但可能不是所有的 User 都有对应的 Person。Party 中的部门 Department 或组织 Organization，都可以对应到 Group。反之 Group 未必对应一个实际的机构。例如，可以有副经理这个 Group，这是多人有相同职责。

引入 Group 这个概念，除了用来解决多人相同角色问题外，还用以解决组织机构的另一种授权问题：例如，A 部门的新闻我希望所有的 A 部门的人都能看。有了这样一个 A 部门对应的 Group，就可直接授权给这个 Group
