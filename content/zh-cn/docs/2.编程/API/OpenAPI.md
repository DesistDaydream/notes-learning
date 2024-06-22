---
title: OpenAPI
linkTitle: OpenAPI
date: 2024-06-21T20:23
weight: 20
---

# 概述

> 参考：
>
> - [官网](https://www.openapis.org/)
> - [GitHub 社区，OAI](https://github.com/OAI)
> - [Swagger](https://swagger.io/)

**OpenAPI Initiative(开放应用程序接口倡议，简称 OAI)**。是由具有前瞻性的行业专家联合创建的，他们认识到标准化 API 描述方式的巨大价值。作为 Linux 基金会下的开放治理结构，OAI 致力于创建、发展 和 推广供所有人可用的中立的描述格式。

OAI 现阶段包含一个规范

- **OpenAPI Specification(简称 OAS)** # 最初基于 SmartBear Software 捐赠的 Swagger 规范。
  - OAS 可以描述为一个文件，可以描述为一个规范的内容。人们常常在不同场景下，不加区分得统一用 OAS 来描述。比如 Swagger Codegen 项目中的描述中。
    - This is the Swagger Codegen project, which allows generation of API client libraries (SDK generation), server stubs and documentation automatically given an OpenAPI Spec.(这是 Swagger Codegen 项目，该项目允许在给定 OpenAPI 规范的情况下自动生成 API 客户端库（生成 SDK），服务器存根和文档。)
    - 这里描述的`指定的的 OpenAPI Spec` 就是指 OAS，也就是说，而已根据 OAS 文件来生成代码。

# OpenAPI Specification

> 参考:
>
> - [官方文档 v3.0.3](http://spec.openapis.org/oas/v3.0.3)
> - [Swagger 官网上的 OAS 文档](https://swagger.io/specification/)
> - [GitHub OAS 项目](https://github.com/OAI/OpenAPI-Specification)

**The OpenAPI Specification(开放应用程序接口规范，简称 OAS), 以前称为 Swagger 规范**，是定义 **RESTful 风格接口**的**世界标准**。OAS 使开发人员能够设计与技术无关的 API 接口，从而构成其 API 开发和使用的基础。

> 由于 OAS 是接着 Swagger 的，所以，OAS 最早的版本就是 3.0，所谓的 OAS 2.0 实际上是 Swagger 2.0。因为那时候 Swagger 还没有变为 OpenAPI

OAS 是 Linux 基金会合作项目 OpenAPI Initiative 中的一个社区驱动的开放规范。

OpenAPI 规范（OAS）为 HTTP API 定义了标准的，与编程语言无关的接口描述，使人和计算机都可以发现和理解服务的功能，而无需访问源代码，其他文档或检查网络流量。通过 OpenAPI 正确定义后，使用者可以使用最少的实现逻辑来理解远程服务并与之交互。与接口描述针对低级编程所做的类似，OpenAPI 规范消除了调用服务时的猜测。

机器可读的 API 定义文档的用例包括但不限于：交互式文档；文档，客户端和服务器的代码生成；和测试用例的自动化。OpenAPI 文档描述了 API 服务，并以 YAML 或 JSON 格式表示。这些文档可以静态生成和提供，也可以从应用程序动态生成。

OpenAPI 规范不需要重写现有的 API。它不需要将任何软件绑定到服务-所描述的服务甚至可能不是其描述的创建者所拥有。但是，它确实需要在 OpenAPI 规范的结构中描述服务的功能。并非所有服务都可以用 OpenAPI 来描述-该规范并非旨在涵盖 HTTP API 的所有可能样式，而是包括对 REST API 的 支持。OpenAPI 规范没有规定特定的开发过程，例如设计优先或代码优先。通过与 HTTP API 建立清晰的交互，它确实有助于这两种技术。

这个 GitHub 项目是 OpenAPI 的起点。在这里，您将找到有关 OpenAPI 规范所需的信息，其外观的简单示例以及有关该项目的一些常规信息。

## 最佳实践

通常，设计 API 规范有两个方向，Design-First（设计优先） 或 Code-First（编码优先）

### Design-First(设计优先)

即优先设计 API 规范，设计完成后再着手进行代码开发工作。采用 Design-First 就意味着，将设计 API 路由、参数等工作提前，后续整个软件开发的流程都需要围绕着 API 规范为核心，当然这需要有一定的设计经验的开发人员才能胜任。Design-First 有很多好处：更多可参考：<https://swagger.io/blog/api-design/design-first-or-code-first-api-development>

1. 提高开发效率。开发团队将根据 API 规范进行并行开发和对接工作，而无需等待接口逻辑开发完毕。
2. 降低接口开发的成本，无需修改代码逻辑即可轻松地修改 API 规范，因为 API 描述语言（如：OpenAPI）与编码语言无关
3. 开发人员更加专注于设计 API 规范，对比 Code-First 可以描写更多 API 的细节，如：校验规则、范例数据等，同时开发人员对 API 的全局一致性、扩展性、规范性也有更好的把控。
4. 在联调开发的过程中可以提前发现和解决问题，避免问题在开发完毕后修改的成本过高。
5. 由于 API 描述更加标准化，可以方便做更多的 API 生态延伸，比如基于 API 规范生成 Mock API Server，开发 API 自动化测试代码，接入 API 网关等。

### Code-First(编码优先)

即通过代码中关于 API 描述特性、注解或注释自动生成 API 描述文件的设计方式，如：JAVA 生态的 SpringFox。适合倾向于在代码上编写 API 规范，通过自动化设施自动生成文档的团队。Code-First 的优点：虽然 Code-First 省去了开发者设计 API 描述文件的阶段，提高了 API 开发者的效率，但是从整个团队的角度来看，效率并不一定提升了，反而有可能降低了效率。不同 API 开发者的经验和习惯的不同，极有可能在编码的过程中对 API 的限制条件考虑不全，又或者框架生成 API 文档的程序完善度不够，种种因素导致最终生成的 API 的描述无法达到理想标准。而很多 API 开发者习惯开发完成后才推送代码，并生成 API 文档，也导致了团队的进程阻塞于此，拖后了整个团队的开发进程。另一方面，API 在开发完成如果没有测试，很有可能导致 API 对接者在对接的过程中遇到重重阻碍。CODING 也希望各位开发者重视 API 设计的重要性，如果您喜欢 Code-First 设计方向，我们的建议是：

1. 节省时间。对于 API 开发者，编码的同时可以获得一份满足基本要求的 API 文档。
2. 方便接入自动化 CI/CD 流程中。
3. 选用完善程度比较高的生成组件
4. 对 API 的描述尽可能的细致、完整
5. 优先设计路由、请求、响应等规则，详细的逻辑代码在 API 设计完成后再着手开发。

Design-First 和 Code-First 针对不同的场景有着各自的优势，不同团队对两者考虑的方向也不同，但是对 API 描述的完善不管是哪个方向都是最重要的。
