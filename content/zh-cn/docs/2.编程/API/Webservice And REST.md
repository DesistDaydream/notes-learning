---
title: Webservice And REST
linkTitle: Webservice And REST
date: 2024-04-15T19:32
weight: 20
---

# 概述

> 参考：
>
> - [阮一峰，RESTful API 设计指南](http://www.ruanyifeng.com/blog/2014/05/restful_api.html)
> - [思否，REST架构的思考](https://segmentfault.com/a/1190000004311893)

## Server Side(服务端)

**Server Side(服务端) 的 WebAPI** 是由一个或多个公开暴露的 **Endpoints(端点)** 组成的编程接口

### Endpoints(端点)

**Endpoints(端点，简称 ep)** 是与服务端 WebAPI 交互的重要方面，因为 Endpoints 指定了客户端可以访问的资源位置。通常，是通过 URI 进行访问，HTTP 请求发送到这个 URI 上，从而期望从这个 URI 上得到响应。

Web services expose one or more endpoints to which messages can be sent. A web service endpoint is an entity, processor, or resource that can be referenced and to which web services messages can be addressed. Endpoint references convey the information needed to address a web service endpoint. Clients need to know this information before they can access a service.

Endpoint 这个词以前经常被用来描述进程间通信。例如

- 在客户端与服务器之间的通讯，客户端是一个 Endpoint 和服务器是另外一个 Endpoint。
- 根据不同的情况下，一个 Endpoint 可能指的地址，如一个 TCP 通信的（主机：端口）对，也可能是指与这个地址相对应的一个软件实体。例如，如果大家使用“ www.example.com:80 ”来描述一个 Endpoint。这些 Endpoint 可能是指实际的端口上的主机名称（即地址）
- 也可能是指与地址相关的的网页服务器（即在这个地址之上运行的软件地址） 。
- 邮寄信件时，将信件投递到邮箱，那么邮箱就是一个 ep
- 寄送快递时，快递员上门取件，快递员就是一个 ep
- webservice 服务，一个服务地址：<http://www.url.com/service1>是一个 ep

## Client Side(客户端)

**Client Side(客户端) 的 WebAPI** 也是一个编程接口，用于扩展 Web 浏览器或其他 HTTP 客户端内的功能。

# Webservice

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

随着时代的发展，Web service 过于笨重的逻辑非常影响程序交互。能用 RestFul 风格的 API 就不要再用 Web service 了。

# REST

> 参考：
>
> - [Wiki, REST](https://en.wikipedia.org/wiki/Representational_state_transfer)

**Representational State Transfer(表述性状态传递，简称 REST)** 是交互式应用程序（通常使用多个 Web 服务实现）的软件架构的事实标准。说白了就是两个应用程序的交互标准。规定了两个程序互相传递数据的格式(JSON、XML 等)、内容、方法(增删改查)等等。这种格式，称为 REST 风格的接口

REST 是由 Roy Thomas Fielding 博士在他的论文 《Architectural Styles and the Design of Network-based Software Architectures》中提出的一个术语。REST 本身只是为分布式超媒体系统设计的一种架构风格，而不是标准。

## REST 规范

原有的 B/S 架构有两种规范：

**1. 无状态性**

无状态性是在客户－服务器约束的基础上添加的又一层规范。他要求通信必须在本质上是无状态的，即从客户到服务器的每个 request 都必须包含理解该 request 所必须的所有信息。这个规范改善了系统的可见性（无状态性使得客户端和服务器端不必保存对方的详细信息，服务器只需要处理当前 request，而不必了解所有的 request 历史），可靠性（无状态性减少了服务器从局部错误中恢复的任务量），可伸缩性（无状态性使得服务器端可以 很容易的释放资源，因为服务器端不必在多个 request 中保存状态）。同时，这种规范的缺点也是显而易见得，由于不能将状态数据保存在服务器上的共享上 下文中，因此增加了在一系列 request 中发送重复数据的开销,严重的降低了效率。

**2. 缓存**

为了改善无状态性带来的网络的低效性，我们填加了缓存约束。缓存约束允许隐式或显式地标记一个 response 中的数据，这样就赋予了客户端缓存 response 数据的功能，这样就可以为以后的 request 共用缓存的数据，部分或全部的消除一部分交互，增加了网络的效率。但是用于客户端缓存了信 息，也就同时增加了客户端与服务器数据不一致的可能，从而降低了可靠性。

B/S 架构的优点是其部署非常方便，但在用户体验方面却不是很理想。为了改善这种情况，我们引入了 REST。 

REST 在原有的架构上增加了三个新规范：统一接口，分层系统和按需代码。

**3. 统一接口**

REST 架构风格的核心特征就是强调组件之间有一个统一的接口，这表现在 REST 世界里，网络上所有的事物都被抽象为资源，而 REST 就是通过通用的链接器接口对 资源进行操作。这样设计的好处是保证系统提供的服务都是解耦的，极大的简化了系统，从而改善了系统的交互性和可重用性。并且 REST 针对 Web 的常见情况 做了优化，使得 REST 接口被设计为可以高效的转移大粒度的超媒体数据，这也就导致了 REST 接口对其它的架构并不是最优的。

**4. 分层系统**

分层系统规则的加入提高了各种层次之间的独立性，为整个系统的复杂性设置了边界，通过封装遗留的服务，使新的服务器免受遗留客户端的影响，这也就提高了系统的可伸缩性。

**5. 按需代码**

REST 允许对客户端功能进行扩展。比如，通过下载并执行 applet 或脚本形式的代码，来扩展客户端功能。但这在改善系统可扩展性的同时，也降低了可见性。所以它只是 REST 的一个可选的约束。

## REST 的设计准则  

REST 架构是针对 Web 应用而设计的，其目的是为了降低开发的复杂性，提高系统的可伸缩性。REST 提出了如下设计准则：

- 网络上的所有事物都被抽象为资源（resource）；
- 每个资源对应一个唯一的资源标识符（resource identifier）；
- 通过通用的连接器接口（generic connector interface）对资源进行操作；
- 对资源的各种操作不会改变资源标识符；
- 所有的操作都是无状态的（stateless）。

#### 注意：

1. REST 中的资源所指的不是数据，而是数据和表现形式的组合，比如“最新访问的 10 位会员”和“最活跃的 10 为会员”在数据上可能有重叠或者完全相同，而 由于他们的表现形式不同，所以被归为不同的资源，这也就是为什么 REST 的全名是 Representational State Transfer 的原因。资源标识符就是 URI(Uniform Resource Identifier)，不管是图片，Word 还是视频文件，甚至只是一种虚拟的服务，也不管你是 xml 格式,txt 文件格式还是其它文件格式，全部通过 URI 对资源进行唯一的标识。
2. REST 是基于 Http 协议的，任何对资源的操作行为都是通过 Http 协议来实现。以往的 Web 开发大多数用的都是 Http 协议中的 GET 和 POST 方 法，对其他方法很少使用，这实际上是因为对 Http 协议认识片面的理解造成的。Http 不仅仅是一个简单的运载数据的协议，而是一个具有丰富内涵的网络软 件的协议。他不仅仅能对互联网资源进行唯一定位，而且还能告诉我们如何对该资源进行操作。Http 把对一个资源的操作限制在 4 个方法以内：GET, POST,PUT 和 DELETE，这正是对资源 CRUD 操作的实现。由于资源和 URI 是一一对应的，执行这些操作的时候 URI 是没有变化的，这和以往的 Web 开发有很大的区别。正由于这一点，极大的简化了 Web 开发，也使得 URI 可以被设计成更为直观的反映资源的结构，这种 URI 的设计被称作 RESTful 的 URI。这位开发人员引入了一种新的思维方式：通过 URL 来设计系统结构。当然了，这种设计方式对一些特定情况也是不适用的，也就是说不 是所有的 URI 都可以 RESTful 的。

## 思考

1. 以上的信息参考自<http://www.cnblogs.com/EasyLive2006/archive/2009/11/03/1595152.html>并整理；
2. 对于网站开发，我们常用的操作就是 GET，和 POST 方式，比如获取数据采用 GET 方式，提交数据采用 POST 方式，但不管是哪种方式，提交数据还是获取数据，后端都要对参数进行处理并对这些操作进行相应。而 REST 的架构把 PUT 和 DELETE 两种数据提交方式用上了，整个操作就会更加的清晰明了，非常的有逻辑性。
3. REST 的 HTTP 协议操作与数据库的 CURD 操作对比:
   HTTP 请求数据库请求 GETSELECTPOSTINSERTPUTUPDATEDELETEDELETE
4. 从以上对比表来看，前端与后端的数据交互更加像是另外一层的数据库交互操作，然而，知道这个没什么卵用，这个是基于 REST 思想的一种 HTTP 协议规范，但如果后端不按照其对应的数据操作来进行处理，没有什么卵用；
   举个例子：比如采用 DELETE 方式发送请求，但后端却处理成增加数据或者修改数据，POST 数据后端处理成删除；但是，我们目前是没有采用这种思想来进行操作的，不管是修改，删除，增加，查询，采用的都是 POST 或者 GET 方式。
5. 数据交互形式：后端返回 json/xml 或者其他前端所期望的数据，但不管什么数据，需要有一个明确的规范和完善的资源管理机制
6. 如果运用了 REST 这种设计思想，我们可以干什么呢？

> 1\). 我们的前端服务器完全可以和数据服务器（REST 服务器）分离，前端服务器处理服务器请求，加载前端骨架（这里不叫框架，叫骨架我觉得更加贴切，REST 服务器上的就是肉），然后前端再更具不同的服务需求，像 REST 请求数据或者更具不同的操作，像 REST 服务器提交增删改的请求等等
> 2\). 不过我觉得把 REST 叫做面向 Api(接口)服务设计来说应该也是很贴切的，REST 服务就是接口。
> 3\). 有了 REST 服务器，不管是电脑端还是手机端，或者是 APP，按照 REST 的接口来进行数据交互，完全不用关心后端实现，也就是说，前端和后端真正的实现了完全的分离设计。

## 后记

关于之后的信息，应该要整理一个简版的前后端代码模型来理解，这个我晚上抽点时间来理一个模型，还有几个方面需要整理：

1. 代码模型实现
2. 架构的不合理性
3. 不适合 REST 架构的场景
4. 其他
