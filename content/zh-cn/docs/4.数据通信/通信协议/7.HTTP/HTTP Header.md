---
title: HTTP Header
---

# 概述

> 参考：
>
> - [RFC 2616-Message Headers](https://tools.ietf.org/html/rfc2616#section-4.2)
> - [RFC 7231,第五章-请求头字段](https://tools.ietf.org/html/rfc7231#section-5)
> - [RFC 7231,第七章-响应头字段](https://datatracker.ietf.org/doc/html/rfc7231#section-7)
> - [Wiki,List of HTTP header fields](https://en.wikipedia.org/wiki/List_of_HTTP_header_fields)
> - [Mozilla 官方文档](https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers)

一般情况下，在浏览器按 F12 后查看到的首部大部分都是请求和响应首部,这俩首部的信息通常包含了通用首部中的信息
![](https://notes-learning.oss-cn-beijing.aliyuncs.com/rzfx0t/1616161207830-cf569808-255f-4e34-9f02-d7da892c9170.jpeg)
HTTP 请求和响应报文的 Header(头) 基本都是 Key/Value Pair(键/值对) 格式的 **Field(字段)**，每个字段都是以冒号分割的 **键/值对**。此外，除了标准的 Header 字段之外，还可以添加自定义 Header，这就给 HTTP 带来了无限的扩展可能。注意，**Key 不区分大小写**。

自定义 Header 历来以 `X-` 开头，但是该约定在 2012 年 6 月被弃用，因为它在非标准字段成为标准字段时会造成不必要的麻烦，详见 [RFC 6648](https://datatracker.ietf.org/doc/html/rfc6648)。IANA 维护了一个通用的 [**HTTP Header 列表**](https://www.iana.org/assignments/message-headers/message-headers.xhtml)，其中包括 RFC 中定义的标准头以及不在 RFC 中定义的扩展头；并且在同一个页面还有新的提议增加的 HTTP Header 列表。

HTTP (RFC 2616 版本)规定了非常多的 Header 字段，可以实现各种各样的功能，但基本上可以分为以下四类

1. **General Header(通用头)** # 在请求头和响应头里都可以出现；
2. **Request Header(请求头)** # 仅能出现在请求头里，进一步说明请求信息或者额外的附加条件；
3. **Response Header(响应头)** # 仅能出现在响应头里，补充说明响应报文的信息；
4. **Entity Header(实体头)** # 它实际上属于通用字段，但专门描述 body 的额外信息。
5. **Extension Header(扩展头)** # 不在标准规范中，可以通过自定义头实现更多定制化需求的 Header 信息。

**对 HTTP 报文的解析和处理其实本质上就是对头字段的处理**，HTTP 的连接管理，缓存控制，内容协商等都是通过头字段来处理的，**理解了头字段，基本上也就理解了 HTTP**，所以理解头字段非常重要。

## 新版 HTTP Header 规范

注意：上面描述分类是在 RFC 2616 标准中定义的，在新的 RFC 7231 中，已经看不到这种分类了，仅仅简单得分为请求头和响应头

- [Request Header(请求头)](https://tools.ietf.org/html/rfc7231#section-5)
- [Response Header(响应头)](https://tools.ietf.org/html/rfc7231#section-7)

而这两类头字段下，又有各自的子分类

在 2014 年之后的新版规范中，并非所有出现在请求中的 Header 都称为请求头，比如 Content-Length，在 RFC 2616 中称为 Entity Header(实体头)，而在新版规范中，称之为元数据。这也为 HTTP 2.0 的 [Header 压缩](/docs/4.数据通信/通信协议/7.HTTP/HTTP2.md) 打下了基础

# Request Header(请求头)

Request Header(请求头) 主要是在每个 HTTP 的请求中指定。包含要获取的资源或请求某个资源的客户端本身的信息，以便服务端可以根据这些内容，来定制响应。

Request Header 中，将各种 Header 分为多个类别

- Controls # 控制本次 HTTP 请求的行为
- Conditionals # 条件相关 Header
- Content Negotiation # 内容协商相关 Header
- Authentication Credentials # 认证相关的 Header
- Request Context # 请求上下文

## Controls(控制)

Controls 类型的请求头用来指定客户端如何处理本次 HTTP 请求。
**Cache-Control** # 缓存控制
**Expect** # 期待服务器的特定行为
**Host** # 请求资源所在服务器。客户端指定自己想访问的服务器的 `域名` 或者 `IP:PORT`。例如：`Host：www.baidu.com`

- Note：输入的什么网址，请求的就是什么，输入域名就是域名，输入 IP 就是 IP
- Note：当服务器接到这个请求时，如果自身无法处理 ip 或者无法处理域名，则该请求就会丢弃(比如 k8s 的 ingress)。所以在测试的时候一般使用 curl 命令请求 IP 时加上 -H 参数以自己制定 URL 内容即可，否则如果服务器不处理 IP 的话，就会返回 404

**Max-Forwards**# 最大传输逐跳数
**Pragma** #
**Range** # 实体的字节范围请求
**TE**# 传输编码的优先级

## Conditionals(条件)

**If-Match** # 比较实体标记(ETag)
**If-None-Match** # 比较实体标记(与 If-Match 相反)
**If-Modified-Since** # 比较资源的更新时间
**If-Unmodified-Since**# 比较资源的更新时间(与 If-Modified-Since 相反)
**If-Range**# 资源未更新时发送实体 Byte 的范围请求

## Content Negotiation(内容协商)

**Accept** # 用户代理可处理的媒体类型
**Accept-\[ Charset | Encoding | Language]** # 通知 server 自己可接收的媒体类型\[字符集|编码格式|语言]

## Authentication Credentials(认证)

**Authorization** # Web 认证信息。IANA 维护了一个[身份验证方案的列表](https://www.iana.org/assignments/http-authschemes/http-authschemes.xhtml)

- **Basic** # 基本认证。就是用户名和密码。如果用户名和密码为 `admin/admin1234` 的话，该字段应该是这样的：
  - `Authorization: Basic YWRtaW46YWRtaW4xMjM0`
  - 也就是说，用户名和密码是 `admin:admin1234` 这样的 base64 编码后的格式
- **Bearer** # 不记名令牌。
- ......

**Proxy-Authorization** # 代理服务器要求客户端的认证信息

## Request Context

**From** # 用户的电子邮箱地址
**Referer** # 对请求中 URI 的原始获取方
**User-Agent** # HTTP 客户端程序的信息

# Response Header(响应头)

包含有关响应的补充信息，如其位置或服务器本身(名称和版本等)的消息头。

Response Header 中，将各种 Header 分为多个类别：

- Control Data
- Validator Header Fields
- Authentication Challenges
- Response Context

## Control Data(控制数据)

控制服务端如何处理 HTTP 响应
Age # 推算资源创建经过时间
Cache-Control #
Expires #
Data #
Location #
Retry-After #
Vary #
Warning #

## Validator Header Fields(验证器头字段)

**ETag** # 资源的匹配信息
**Last-Modified** # 最后一次修改时间

## Authentication Challenges

**WWW-Authenticate** # 服务器对客户端的认证信息
**Proxy-Authorization** # 代理服务器要求客户端的认证信息

## Response Context

**Accept-Ranges** # 是否接受字节范围请求
**Allow** #
**Server**# HTTP 服务器的安装信息

# **Extension Header(扩展头)**

通常情况下，一个 web 界面有 N 多个资源，比如 index.html 首页是一个资源，首页中有各种 img，css，js 等静态或动态资源，当用户访问一个网站后，除了要请求主页资源，还需要再请求主页上各个图片，板式，功能等资源，每一个资源都有一个单独的报文。也就说有可能每一类资源都可以单独存在在一个相对应的集群服务器群中，比如通过一个图片的 URL，可以直接访问该图片，而不用在网站主页才能看到；而这些资源都需要前端负载均衡器来进行调度用户请求到相应的服务器中去拿去资源，这时候 LB 的压力会非常大，为了解决这种问题，人们就可以根据自己的协商，定义一些标准之外的 Header。

# 老版本的两个 Header 笔记

## General Header(通用头)

同时适用于请求和响应消息，但与最终消息主体中传输的数据无关的消息头。
**Cache-Control** # 控制缓存的行为

- no-cache：不缓存
- max-age：缓存最大保留时间

**Connection** # 逐跳首部、连接的管理。连接状态。HTTP/2 中不在支持该 Header

- keepalive(长连接)：需要实时更新的网络资源，比如网络游戏，具体点就是网页上的游戏，需要实时连接，如果游戏中每一步操作都要进行 TCP 建连，然后各种请求和响应，那么用户体验会非常不好，也会极大得消耗网络带宽，所以需要使用类似 keepalive 这种长连接
- close(短连接)：不需要实时更新的网络资源，当浏览静态网页资源的时候，当把整个网页发送给用户之后即可关闭连接，因为网页资源已经发送到 Client 本地，不需要实时更新了

**Date**# 创建报文的日期时间
**Pragma**# 报文指令
**Trailer**# 报文末端的首部一览
**Transfer-Encoding**# 指定报文主体的传输编码方式
**Upgrade**# 升级为其他协议
**Via**# 报文经过的中间节点。也就是代理服务器的相关信息。
**Warning** # 错误通知

## Entity Header(实体头)

包含有关实体主体的更多信息，比如主体长(Content-Length)度或其 MIME 类型。
Allow：列出对此实体可使用的请求方法
Location：告诉 client 真正的实体位于何处
Content-
缓存相关

1. ETag # 实体扩展标签
2. Expires # 实体过期时间
3. Last-Modified # 最后一次修改时间

**Content-Length**# 实体主体的大小(单位 :字节)
