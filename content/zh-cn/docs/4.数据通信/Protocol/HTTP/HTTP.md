---
title: HTTP
linkTitle: HTTP
date: 2023-11-21T15:43
weight: 1
---

# 概述

> 参考：
>
> - [RFC 2616](https://tools.ietf.org/html/rfc2616)
> - [Mozilla 官方 HTTP 开发文档](https://developer.mozilla.org/en-US/docs/Web/HTTP)
> - [公众号-小林 coding，硬核！30 张图解 HTTP 常见的面试题](https://mp.weixin.qq.com/s/bUy220-ect00N4gnO0697A)
> - [公众号-码海，51 张图助你彻底掌握 HTTP 协议](https://mp.weixin.qq.com/s/WQpxfwLArltKEjEAdOO2Pw)

**HyperText Transfer Protocol(超文本传输协议，简称 HTTP)**。是基于 TCP 的用于分布式、协作式、超媒体的信息系统的应用层协议。HTTP 是 [World Wide Web(万维网,简称 WWW.就是我们俗称的 Web)](https://en.wikipedia.org/wiki/World_Wide_Web) 的数据通信基础。

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/rlkp9t/1616161240441-f2958719-b738-4698-9fca-64d90f3471ba.png)

## HTTP 标准的演化

> 参考：
>
> - [InfoQ 中的消息](https://www.infoq.cn/article/2014/06/http-11-updated)

在 2014 年之前，HTTP/1.1 版本的标准为 [RFC 2616](https://tools.ietf.org/html/rfc2616)，但由于[某些原因](https://tools.ietf.org/html/rfc7230#appendix-A.2)，为了让标准更规范，HTTP/1.1 被拆分成了 6 个部分：

- [RFC7230 - HTTP/1.1](https://tools.ietf.org/html/rfc7230): Message Syntax and Routing(消息语法和路由)。这里包含 低级的消息解析 和 链接管理。
- [RFC7231 - HTTP/1.1](https://tools.ietf.org/html/rfc7231): Semantics and Content(语意和内容)。这里面包含了 Methods、Status Codes、Headers
- RFC7232 - HTTP/1.1: Conditional Requests - e.g., If-Modified-Since
- RFC7233 - HTTP/1.1: Range Requests - getting partial content
- RFC7234 - HTTP/1.1: Caching - browser and intermediary caches
- RFC7235 - HTTP/1.1: Authentication - a framework for HTTP authentication

## HTTP 三个部分

### 1. 「协议」

在生活中，我们也能随处可见「协议」，例如：

- 刚毕业时会签一个「三方协议」；
- 找房子时会签一个「租房协议」；

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/rlkp9t/1616161240448-b1263f75-a700-4431-9d6b-a99b36a58214.png)

三方协议和租房协议

生活中的协议，本质上与计算机中的协议是相同的，协议的特点:

- 「**协**」字，代表的意思是必须有**两个以上的参与者**。例如三方协议里的参与者有三个：你、公司、学校三个；租房协议里的参与者有两个：你和房东。
- 「**议**」字，代表的意思是对参与者的一种**行为约定和规范**。例如三方协议里规定试用期期限、毁约金等；租房协议里规定租期期限、每月租金金额、违约如何处理等。

针对 HTTP **协议**，我们可以这么理解。

HTTP 是一个用在计算机世界里的**协议**。它使用计算机能够理解的语言确立了一种计算机之间交流通信的规范（**两个以上的参与者**），以及相关的各种控制和错误处理方式（**行为约定和规范**）。

### 2. 「传输」

所谓的「传输」，很好理解，就是把一堆东西从 A 点搬到 B 点，或者从 B 点 搬到 A 点。

别轻视了这个简单的动作，它至少包含两项重要的信息。

HTTP 协议是一个**双向协议**。

我们在上网冲浪时，浏览器是请求方 A ，百度网站就是应答方 B。双方约定用 HTTP 协议来通信，于是浏览器把请求数据发送给网站，网站再把一些数据返回给浏览器，最后由浏览器渲染在屏幕，就可以看到图片、视频了。
![](https://notes-learning.oss-cn-beijing.aliyuncs.com/rlkp9t/1616161240434-514ef4d5-6830-46f1-9674-7c6564798dfd.png)

Request(请求) - (Response)应答

数据虽然是在 A 和 B 之间传输，但允许中间有**中转或接力**。

就好像第一排的同学想传递纸条给最后一排的同学，那么传递的过程中就需要经过好多个同学（中间人），这样的传输方式就从「A < --- > B」，变成了「A <-> N <-> M <-> B」。

而在 HTTP 里，需要中间人遵从 HTTP 协议，只要不打扰基本的数据传输，就可以添加任意额外的东西。

针对**传输**，我们可以进一步理解了 HTTP。

HTTP 是一个在计算机世界里专门用来在**两点之间传输数据**的约定和规范。

### 3. 「超文本」

HTTP 传输的内容是「超文本」。

我们先来理解「文本」，在互联网早期的时候只是简单的字符文字，但现在「文本」的涵义已经可以扩展为图片、视频、压缩包等，在 HTTP 眼里这些都算作「文本」。

再来理解「超文本」，它就是**超越了普通文本的文本**，它是文字、图片、视频等的混合体，最关键有超链接，能从一个超文本跳转到另外一个超文本。

HTML 格式的文件就是最常见的超文本了，它本身只是纯文字文件，但内部用很多标签定义了图片、视频等的链接，再经过浏览器的解释，呈现给我们的就是一个文字、有画面的网页了。

## 总结

OK，经过了对 HTTP 里这三个名词的详细解释，就可以给出比「超文本传输协议」这七个字更准确更有技术含量的答案：

**HTTP 是一个在计算机世界里专门在「两点」之间「传输」文字、图片、音频、视频等「超文本」数据的「约定和规范」。**

> 这里的两点可以是服务器到本地电脑，本地电脑到服务器、服务器到服务器、电脑到电脑，等等。

# HTTP 请求过程概述

HTTP 事务由一次 **Client 的 Request(请求)** 和 **Server 的 Response(响应)** 组成。

HTTP 是一个客户端—服务器协议：请求由一个实体，即 **User Agent(用户代理)**，或是一个可以代表它的代理方（proxy）发出。大多数情况下，这个**用户代理都是一个网页浏览器**，不过它也可能是任何东西，比如一个爬取网页来充实、维护搜索引擎索引的机器爬虫（其实就是代码写的具有发起 HTTP 请求的程序，毕竟浏览器也是代码写的）。

> 之所以用 [用户代理](docs/Web/Web.md#Glossary(术语)) 这个词，早期主要是用来描述实体的，发起请求的行为应该是用户操作的，而代替用户发起请求的，就是 用户的代理 了嘛。。。早期一般没有爬虫之类的东西，都是用户使用浏览器这个用户代理发起请求的。

每个请求都会被发送到一个服务端，它会处理这个请求并提供一个称作 _响应_ 的回复。在客户端与服务端之间，还有许许多多的被称为[代理](https://developer.mozilla.org/zh-CN/docs/Glossary/Proxy_server)的实体，履行不同的作用，例如充当网关或[缓存](https://developer.mozilla.org/zh-CN/docs/Glossary/Cache)。

## HTTP 的无状态

HTTP 协议是 Stateless(无状态)。(因为连接一次后就断开了，不会持久化存储任何数据)

比如一个用户(Client 客户端)向服务器发起了一个请求，请求一个页面，在该页面输入完用户名和密码后进行登录后，如果刷新页面，那么就需要重新输入用户名和密码，因为 client 向 server 只请求了一个页面，请求完成后，连接就断开了，后续的请求是新的，没法再用以前的信息。这时候为了解决该问题，引用了 Cookie 和 Session 保持 的概念。

相关技术

- Cookie：类似于 Token，相当于一个令牌，当访问一个 web server 的时候，server 发给 client 一个 Cookie，让 client 保存在本地，再次访问的时候，即可通过该 Cookie 识别身份
- Session(会话)保持：关联至 Cookie：当你在动态网页上访问了一些信息，比如购物车，在购物车添加一件物品，即通过 Session 功能来保存该信息，以便下次登录还能使用。否则下次登录购物车中的东西就没了

## 流程简述

- 建立或处理连接，接收请求或拒绝请求
  - 建立 TCP 连接，WEB 浏览器向 Web 服务器发送请求
  - web 浏览器发送请求头部信息
    - 建立连接后，客户机发送一个请求给服务器，请求方式的格式为：URL、协议版本号、后面是 MIME 信息包括请求修饰符、客户机信息和可能得内容
- Server 接收请求，并应答
  - WEB 服务器接收到请求后，给予相应的响应信息，其格式为一个状态行，包括信息的协议版本号、一个成功或错误的代码，后边是 MIME 信息包括服务器信息、实体信息和可能得内容
- Server 处理请求
  - Web 服务器发送应答头信息
  - Web 服务器向浏览器发送数据
- Client 访问资源
  - WEB 服务器关闭 TCP 连接
- 构建响应报文
- 发送响应报文
- 记录日志

# HTTP 报文格式

## Request 与 Response 报文

### Request 请求报文

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/rlkp9t/1616161240468-d9f13310-3f67-43b8-b5b7-b48dde683170.png)

- **Method** # 请求方法，用于表明 Client 希望 Server 对 Resource 执行的动作。常用：GET、POST、DELETE
- **URL** # HTTP 请求的 URL。
  - **Params** # URL 参数。就是 URL 中的 Query 部分
- **Version** # 发送给 Server 的请求的 HTTP 协议版本。
- **Headers** # 请求头。
- **Body** # 请求体。

> 这里面有一个要注意的地方，就是 Params 与 Headers，**Params 是 URL 的一部分**，但是 Headers 不是。虽然两者的作用类似，都是用来定义这个请求中应该发送给对方的一些基本信息、认证信息 等等。但是在一个 HTTP 的请求中，两者所处的位置是不一样，用于不同场景。

**Authorization** # 认证信息。这是一个比较特殊的东西，可以存在于 URL 的 Params 中、Headers 中、Body 中。请求报文的各个部分，都可以填写认证信息。

- 当 Server 需要一个认证信息时，就需要在 HTTP 请求中加入认证相关的信息。

#### EXAMPLE

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/rlkp9t/1616161240442-63d3fc7f-80c2-43e4-bec5-50061f2e02f2.png)

### Response 响应报文

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/rlkp9t/1616161240462-a29c9d65-119a-4b70-993f-bd1a4cfbbd7e.png)

- **Version** # 响应给 Client 的 HTTP 版本。
- **Status** # HTTP 响应状态。用来标记请求过程中发生的情况，由 server 告诉 client。响应状态由两部分组成
- **StatusCode** # 状态码。统一为 3 位的数字。
  - 各个状态码的含义，见 [HTTP Status Codes](/docs/4.数据通信/Protocol/HTTP/HTTP%20Status%20Codes.md)
- **ReasonPhrase** # 原因短语。用来表示产生该状态的原因的简要说明
- **Headers** # 响应头。
- **Body** # 响应体。实体部分，请求时附加的数据或响应时附加的数据

#### EXAMPLE

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/rlkp9t/1616161240440-85fac56c-d3ed-44dc-94f3-20d60017e622.png)

## HTTP Header

> 参考：
>
> - [RFC 2616，14 Header Field Definitions](https://datatracker.ietf.org/doc/html/rfc2616#section-14)

HTTP 请求和响应报文的 Header(头) 基本都是 Key/Value(键值) 格式，Key 与 Value 以冒号分隔，此外，除了标准的头部字段之外，还可以添加自定义头，这就给 HTTP 带来了无限的扩展可能。注意，Value 不区分大小写。

HTTP 协议规定了非常多的 Header 字段，可以实现各种各样的功能，但基本上可以分为以下几类

- **General Header(通用头)** # 在请求头和响应头里都可以出现；
- **Request Header(请求头)** # 仅能出现在请求头里，进一步说明请求信息或者额外的附加条件；
- **Response Header(响应头)** # 仅能出现在响应头里，补充说明响应报文的信息；
- **Entity Header(实体头)** # 它实际上属于通用字段，但专门描述 body 的额外信息。
- **Extension Header(扩展头)** # 不在标准规范中，可以通过自定义头实现更多定制化需求的 Header 信息。

**对 HTTP 报文的解析和处理其实本质上就是对头字段的处理**，HTTP 的连接管理，缓存控制，内容协商等都是通过头字段来处理的，**理解了头字段，基本上也就理解了 HTTP**，所以理解头字段非常重要。

详见：[HTTP Header](/docs/4.数据通信/Protocol/HTTP/HTTP%20Header.md)

> 注意：这种分类是在 RFC 2616 标准中定义的，在新的 RFC 7231 中，已经看不到这种分类了

## HTTP 的请求 Method(方法)

> 参考：
>
> - [RFC 2616，5.1.1 Method](https://datatracker.ietf.org/doc/html/rfc2616#section-5.1.1)
> - [Mozilla 官方 HTTP 开发文档，HTTP 请求方法](https://developer.mozilla.org/en-US/docs/Web/HTTP/Methods)

根据 HTTP 标准，HTTP 请求可以使用多种请求方法。 HTTP 的 1.0 版本中只有三种请求方法： GET, POST 和 HEAD 方法。到了 1.1 版本时，新增加了五种请求方法：OPTIONS, PUT, DELETE, TRACE 和 CONNECT 方法。

- **GET** # 从服务器获取了资源
  - 请求指定的页面信息，并返回实体主体。
  - GET 请求请提交的数据放置在 HTTP 请求协议头中，GET 方法通过 URL 请求来传递用户的输入，GET 方式的提交你需要用 Request.QueryString 来取得变量的值。
  - GET 方法提交数据，可能会带来安全性的问题，数据被浏览器缓存。
  - GET 请求有长度限制。
- **HEAD** # 只从 server 获取文档的响应首部（报文中的 Headers）
  - 类似于 get 请求，只不过返回的响应中没有具体的内容，用于获取报头。
- **POST** # 向 server 发送要处理的数据
  - 向指定资源提交数据进行处理请求（例如提交表单或者上传文件）。
  - POST 请求可能会导致新的资源的建立和/或已有资源的修改。
  - POST 方式提交时，你必须通过 Request.Form 来访问提交的内容
- **PUT** # 将请求的主体存储在 server 上
  - 从客户端向服务器传送的数据取代指定的文档的内容。
- **DELETE** # 请求删除 server 上通过 URL 指定的文档，DELETE 请求一般返回 3 种码
  - 200（OK）——删除成功，同时返回已经删除的资源。
  - 202（Accepted）——删除请求已经接受，但没有被立即执行（资源也许已经被转移到了待删除区域）。
  - 204（No Content）——删除请求已经被执行，但是没有返回资源（也许是请求删除不存在的资源造成的）。
- **OPTIONS** # 请求服务器返回对指定资源支持使用的请求方法
  - 允许客户端查看服务器的性能。
- **TRACE** # 追踪请求到达 server 中间经过的 server agent
  - 回显服务器收到的请求，主要用于测试或诊断。

### GET 与 POST

`Get` 方法的含义是请求**从服务器获取资源**，这个资源可以是静态的文本、页面、图片视频等。

比如，你打开我的文章，浏览器就会发送 GET 请求给服务器，服务器就会返回文章的所有文字及资源。

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/rlkp9t/1649668233003-6c01ab91-90f8-4ea8-8095-5448302146e8.jpeg)

而`POST` 方法则是相反操作，它向 `URI` 指定的资源提交数据，数据就放在报文的 body 里。

比如，你在我文章底部，敲入了留言后点击「提交」（**暗示你们留言**），浏览器就会执行一次 POST 请求，把你的留言文字放进了报文 body 里，然后拼接好 POST 请求头，通过 TCP 协议发送给服务器。
![](https://notes-learning.oss-cn-beijing.aliyuncs.com/rlkp9t/1649668232954-1a090cec-1a9c-4c11-a370-d49b4156cac3.jpeg)
GET 和 POST 方法都是安全和幂等的吗？

先说明下安全和幂等的概念：

- 在 HTTP 协议里，所谓的「安全」是指请求方法不会「破坏」服务器上的资源。
- 所谓的「幂等」，意思是多次执行相同的操作，结果都是「相同」的。

那么很明显 **GET 方法就是安全且幂等的**，因为它是「只读」操作，无论操作多少次，服务器上的数据都是安全的，且每次的结果都是相同的。

**POST** 因为是「新增或提交数据」的操作，会修改服务器上的资源，所以是**不安全**的，且多次提交数据就会创建多个资源，所以**不是幂等**的。

# HTTP 相关的协议或规范

HTTP 无法单独存在，要想让它生效，必须依赖其他的协议或者规范

## URL

详见 [URL 与 URI](/docs/4.数据通信/Protocol/HTTP/URL%20与%20URI.md)

## TCP/IP

详见：[TCP_IP](/docs/4.数据通信/Protocol/TCP_IP/TCP_IP.md)

## DNS

详见：[DNS](/docs/4.数据通信/DNS/DNS.md)
