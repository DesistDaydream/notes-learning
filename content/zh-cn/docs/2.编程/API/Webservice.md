---
title: "Webservice"
linkTitle: "Webservice"
weight: 20
---

# 概述

> 参考：
>
> - [Wiki, Webservice](https://en.wikipedia.org/wiki/Web_service)
> - https://hygraph.com/blog/web-service-vs-api

Webservice 描述了一种在 Internet 协议主干上使用 XML、SOAP、WSDL 和 UDDI 开放标准集成基于 Web 的应用程序的标准化方法。XML 是用于包含数据并提供围绕数据的元数据的数据格式，SOAP 用于传输数据，WSDL 用于描述可用的服务，UDDI 列出可用的服务。

什么是 Web service 的答案在很大程度上取决于正在进行的对话的上下文。在一种情况下，我们可以将 Web 服务视为支付生态系统、存储服务、电子邮件服务、云功能、文本到语音转换器等。这些服务本身就是单独的系统，可以为您节省大量“从头开始构建”的时间。当您编写应用程序代码时。任何想要使用它们的人都可以通过即用即付模式进行注册。

回到 2000 年代，在另一个上下文中，我们可以说 Web 服务是一组促进不同软件应用程序之间数据交换的协议和标准。

以下是构建 Web 服务所需的组件。

- SOAP - Simple Object Access Protocol - 简单对象访问协议
- RDF - Resource Description Framework - 资源描述框架
- UDDI - Universal Description, Discovery and IntegrationI - 通用描述、发现和集成
- WSDL - Web Services Description Language - Web 服务描述语言

构建 Web 服务需要遵循严格的规则，并且它们往往更加流程密集和代码繁重。 Web 服务与 SOA（面向服务的架构）大致相同，并且主要基于 XML-RPC 和 SOAP 等标准。通常使用 SOAP，数据交换通过 HTTP 协议以 XML 形式进行。

API 一词指的是应用程序编程接口，类似于 Web 服务，当两个软件组件需要相互通信时，它们可以借助 API 来实现。 API 利用技术并遵循协议来促进通信。有不同的 API，例如 GraphQL、REST、WebSocket API、RPC API 等。所有 API 都有自己的一套协议和定义，因此每种 API 类型都会有不同的运行机制。

最终，API 可以被认为包括但不限于 Web service，可以说 Web service 是 API 的子集，是一种特殊情况下的 API。

随着时代的发展，Web service 过于笨重的逻辑非常影响程序交互。能用 [REST](/docs/2.编程/API/REST.md)Ful 风格的 API 就不要再用 Web service 了。
