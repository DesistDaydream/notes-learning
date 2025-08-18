---
title: Swagger
linkTitle: Swagger
weight: 20
---

# 概述

> 参考：
> 
> - [官网](https://swagger.io/)
> - [Swagger 相关工具](https://swagger.io/tools/open-source/open-source-integrations/)

Swagger 有两种含义

- **Swagger Specification(规范)** # 用于描述现代 API 的行业标准，已被广泛采用。
  - 规范就是用 JSON 或 YAML 来描述的一组 API。
- **Swagger ToolsSet(工具集)** # 实现 Swagger 规范 的一系列工具的集合。
  - Swagger 衍生出来的一系列工具，可以做到生成各种格式的接口文档，生成多种语言的客户端和服务端的代码，以及在线接口调试页面等等

注意：Swagger 规范已于 2015 年捐赠给 Linux 基金会后改名为 OpenAPI，并定义最新的规范为 OpenAPI 3.0。

> 所以，3.0 之前的规范，准确来说叫 Swagger 1.0、 Swagger 2.0，而 Swagger 3.0 就应该称为 OpenAPI 3.0 了
> 所以现在，Swagger 这个词语，更多的是用来描述一系列工具的合集。借助 Swagger 开源和专业工具集，为用户，团队和企业简化 API 开发。

Swagger 背景

随着互联网技术的发展，现在的网站架构基本都是由原来的后端渲染，变成了前端渲染、前后端分离的形式，而且前端技术和后端技术在各自的道路上越走越远。前后端的唯一联系，变成了 [API](docs/2.编程/API/API.md) 接口；API 文档变成了前后端开发人员联系的纽带。最早的时候，大家都是手写 API 文档的，在什么地方写的都有，有的在 confluence 上、有的直接在项目的 README.md 文件中写。并且，每个项目的写法、格式都不一样。

此时，就出现问题了，大家都需要写 API、用 API，那么何不出台一个规范，来规范所有项目的 API 文档的写法规范。**这就是 Swagger 的由来**。

Swagger 诞生于 2010 年，起初只是一个用于设计 RESTful 风格 API 的开源规范。为了更好得实现和可视化规范中定义得 API，还开发了诸如 Swagger UI、Swagger Editor、Swagger Codegen 等工具。由 Specification(规范) 和 Tools(工具) 组成的 Swagger 项目非常流行，并由此创造了一个庞大的社区驱动工具的生态系统。

2015 年，Swagger 项目被 SmartBear 软件公司收购。其中 **Swagger Specification(Swagger 规范)** 被捐赠给 Linux 基金会，并更名为 **OpenAPI**。OpenAPI 继续致力于，正式规范 REST API 的描述方式。OpenAPI Initiative(倡议) 的创建则是为了以开放和透明的方式指导 **OpenAPI Specification(简称 OAS)** 的开发。

自此，Swagger 成为最受欢迎的工具套件，可以在整个 API 生命周期中充分利用 OAS 的力量。

# Swagger 与 OpenAPI 的区别

> 参考：
> 
> - [SmartBear 博客](https://smartbear.com/blog/what-is-the-difference-between-swagger-and-openapi/)

2017 年 OpenAPI 3.0 正式发布，这是 [OpenAPI 规范 的第一个版本](https://www.openapis.org/blog/2017/07/26/the-oai-announces-the-openapi-specification-3-0-0)。

对于那些参与 API 开发的人来说，OAS 3.0 的发布是……很大的事情。

为什么？该版本如此重要的最显着原因之一是，OpenAPI 3.0 是该规范的第一个正式版本，因为它是 [由 SmartBear Software 捐赠给 OpenAPI Initiative](https://smartbear.com/news/news-releases/smartbear-launches-open-api-initiative-with-key-in/)， 并在 2015 年从 Swagger 规范重命名为 OpenAPI 规范。

在我们探讨 OpenAPI 3.0 对 API 空间如此重要的一些原因之前，首先必须弄清有关 OpenAPI 及其对 Swagger 的含义的一些问题，这一点很重要。

在过去的两年中，关于从 Swagger 转换为 OpenAPI 的问题很多。而且，关于 OpenAPI 和 Swagger 之间的区别，何时使用一个名称代替另一个名称以及 OpenAPI 和 Swagger 之间的关系还有很多困惑。

## 让我们从澄清 Swagger 与 OpenAPI 开始

理解差异的最简单方法是：

- OpenAPI = 规范

- Swagger = 实现规范的工具

OpenAPI 是规范的正式名称。该规范的开发是由 OpenAPI Initiative 推动的，该倡议涉及更多来自技术领域不同领域的 30 个组织-包括 Microsoft，Google，IBM 和 CapitalOne。领导 Swagger 工具开发的公司 Smartbear Software 也是 OpenAPI Initiative 的成员，帮助领导了规范的发展。

Swagger 是与用于实现 OpenAPI 规范的一些最著名，使用最广泛的工具相关联的名称。Swagger 工具集包括开源工具，免费工具和商业工具的组合，可在 API 生命周期的不同阶段使用。

这些工具包括：

- **Swagger 编辑器**：使用 Swagger 编辑器，您可以在浏览器内的 YAML 中编辑 OpenAPI 规范，并实时预览文档。
- **Swagger UI：** Swagger UI 是 HTML，Javascript 和 CSS 资产的集合，这些资产通过符合 OAS 的 API 动态生成精美的文档。
- **Swagger Codegen：** 允许在给定 OpenAPI 规范的情况下自动生成 API 客户端库（SDK 生成），服务器存根和文档。
- **Swagger Parser：** 用于从 Java 解析 OpenAPI 定义的独立库
- **Swagger Core：** 与 Java 相关的库，用于创建，使用和使用 OpenAPI 定义
- **Swagger Inspector（免费）** ：API 检查工具，可让您从现有 API 生成 OpenAPI 定义
- **SwaggerHub（免费和商业）**：API 设计和文档，为使用 OpenAPI 的团队而构建。

由于 Swagger 工具是由参与创建原始 Swagger 规范的团队开发的，因此通常仍将这些工具视为该规范的代名词。但是 Swagger 工具并不是唯一可用于实现 OpenAPI 规范的工具。有各种各样的 API 设计，文档，测试，管理和监视解决方案都支持该规范的 2.0 版，并且正在积极致力于增加 3.0 支持。

您可以在[GitHub](https://github.com/OAI/OpenAPI-Specification/blob/master/IMPLEMENTATIONS.md)上找到支持最新版本的 OpenAPI 规范的工具的完整列表 。

## Swagger 工具为什么没有将其名称更改为 OpenAPI？

Swagger 生态系统一直由规范和围绕它的核心开源工具组成，其中最著名的是 Swagger UI，Swagger 编辑器和 Swagger Codegen。规范之所以被广泛采用的一个重要原因是其附带的工具。

SmartBear 捐赠了该规范，但由于开发人员，技术作家，测试人员和设计师与该工具有着密切的联系，因此流行的开源 Swagger 工具仍保留了原始品牌。该规范不是，也从未完全与 Swagger 工具相关联。实际上，决定捐赠该规范并组成 OpenAPI Initiative 是为了确保 OpenAPI 完全与供应商无关。这就是为什么我们很高兴看到 API 领域有这么多，包括也支持其他定义格式（例如 API 蓝图和 RAML）的公司加入了该计划的原因。

Swagger 团队始终致力于使用 OpenAPI 规范构建功能最强大，易于使用的工具，以设计，记录，开发和测试 API，并将继续发展和发展我们的工具集以支持 OpenAPI。这些工具将继续保持 Swagger 名称。 Swagger.io 是 Swagger 工具和开源 Swagger 项目的在线主页，还将继续成为学习 Swagger 工具的理想之地，并且我们还将继续为 OpenAPI 规范的知识做出贡献，通过有关使用 OpenAPI 的培训，教程，网络研讨会和文档。

## 了解 OpenAPI 和 Swagger 社区

尽管为 OpenAPI 做出贡献的人们和为 Swagger 工具做出贡献的人们之间总是存在重叠，但是这两个社区是彼此独立的。

如本文所述，OpenAPI Initiative 是一个开放的，与供应商无关的组织，欢迎任何想要帮助发展或利用其 API 开发中的规范的人参与。邀请组织加入不断增加的对规范做出贡献的成员列表，并且欢迎个人参加，方法是在 GitHub 上分享想法和反馈，或者参加每月在世界各地举行的许多 OAS 聚会之一。 在此处了解有关如何做出贡献的更多信息。

Swagger 工具拥有自己的社区，致力于帮助改进某些现有 Swagger 项目，并引入新的想法和功能要求。Swagger 社区是由 SmartBear Software 的团队培育的，该团队投资于开源 Swagger 工具的开发，但也受到世界各地成千上万 Swagger 用户的贡献的推动。如果您想加入 Swagger 社区，我们邀请您 在 GitHub 上找到我们 或加入 Swagger API Meetup 组。 您还可以在 Swagger 博客或 Twitter 上的@SwaggerAPI 上找到最新新闻和更新 。

## 期待 OpenAPI 的美好未来

我们期待看到 OpenAPI 成为 API 领域中每个人都认可的名称，我们很高兴成为不断壮大的 OpenAPI Initiative 成员社区的一部分。

希望本文有助于阐明有关 OpenAPI 及其与 Swagger 的关系的一些问题。

回顾一下：

- 该规范在 2015 年重命名为 OpenAPI 规范。OpenAPI3.0 是该规范的最新版本。

- SmartBear Software 支持的 Swagger 工具是实现 OpenAPI 规范的最受欢迎的工具之一，并且将继续维护 Swagger 名称（Swagger 编辑器，Swagger UI，SwaggerHub 等）。

- 还有数百种与 Swagger 不相关的其他开源和专业工具都支持 OpenAPI 2.0 规范，并且支持 3.0 的工具列表也在不断增长。

- OpenAPI 和 Swagger 都有开源社区，欢迎所有贡献者加入以分享他们的想法并参与其中。

# Swagger 2.0 与 OpenAPI 格式对比

> 参考：
> 
> - [简书大佬](https://www.jianshu.com/p/879baf1cff07)

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/xxryy5/1616163614015-5221ca71-70a5-495b-a656-b5bc9a0122ec.png)

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/xxryy5/1616163614002-f8ab633b-9c2b-42c9-acf2-f9cdd34bf2b1.png)
