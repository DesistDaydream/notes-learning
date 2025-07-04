---
title: WebSocket
linkTitle: WebSocket
weight: 1
---

# 概述

> 参考：
>
> - [RFC 6455, The WebSocket Protocol](https://datatracker.ietf.org/doc/html/rfc6455)
> - [Wiki, WebSocket](https://en.wikipedia.org/wiki/WebSocket)
> - [公众号-小林 coding，有了 HTTP 协议，为什么还要有 websocket 协议？](https://mp.weixin.qq.com/s/TtRKkVxS6H-miQ8luQgY1A)

WebSocket 是一种计算机通信协议，通过单个 TCP 连接提供**全双工**通信通道。

WebSocket 与 [HTTP](/docs/4.数据通信/Protocol/HTTP/HTTP.md) 不同。这两种协议都位于 [OSI 模型](/docs/4.数据通信/数据通信/OSI%20模型.md)的第 7 层，并依赖于第 4 层的 TCP。尽管它们不同，但 RFC 6455 指出 WebSocket “旨在通过 HTTP 端口 443 和 80 工作，并支持 HTTP 代理和中介” ，从而使其与 HTTP 兼容。为了实现兼容性，WebSocket [握手](https://en.wikipedia.org/wiki/Handshaking)使用 [HTTP Upgrade 头](https://en.wikipedia.org/wiki/HTTP/1.1_Upgrade_header)从 HTTP 协议更改为 WebSocket 协议。

WebSocket 协议支持 [Web 浏览器](https://en.wikipedia.org/wiki/Web_browser)（或其他客户端应用程序）和 [Web 服务器](https://en.wikipedia.org/wiki/Web_server)之间的交互，其开销比半双工替代方案（例如 HTTP[轮询）低](<https://en.wikipedia.org/wiki/Polling_(computer_science)>)，从而促进从服务器到服务器的实时数据传输。这是通过为服务器提供一种标准化的方式来向客户端发送内容而无需客户端首先请求，并允许消息在保持连接打开的同时来回传递而实现的。通过这种方式，可以在客户端和服务器之间进行双向正在进行的对话。通信通常通过 TCP[端口](<https://en.wikipedia.org/wiki/Port_(computer_networking)>)完成数字 443（或在不安全连接的情况下为 80），这对于使用[防火墙](<https://en.wikipedia.org/wiki/Firewall_(computing)>)阻止非网络 Internet 连接的环境有益。类似的浏览器-服务器双向双向通信已经使用[Comet](<https://en.wikipedia.org/wiki/Comet_(programming)>)或[Adobe Flash Player](https://en.wikipedia.org/wiki/Adobe_Flash_Player)等临时技术以非标准化方式实现。[\[2\]](https://en.wikipedia.org/wiki/WebSocket#cite_note-2)

大多数浏览器都支持该协议，包括[Google Chrome](https://en.wikipedia.org/wiki/Google_Chrome)、[Firefox](https://en.wikipedia.org/wiki/Firefox)、[Microsoft Edge](https://en.wikipedia.org/wiki/Microsoft_Edge)、[Internet Explorer](https://en.wikipedia.org/wiki/Internet_Explorer)、[Safari](<https://en.wikipedia.org/wiki/Safari_(web_browser)>)和[Opera](https://en.wikipedia.org/wiki/Opera_web_browser)。[\[3\]](https://en.wikipedia.org/wiki/WebSocket#cite_note-3)
与 HTTP 不同，WebSocket 提供全双工通信。[\[4\]](https://en.wikipedia.org/wiki/WebSocket#cite_note-4)[\[5\]](https://en.wikipedia.org/wiki/WebSocket#cite_note-quantum-5) 此外，WebSocket 支持基于 TCP 的消息流。TCP 单独处理字节流，而没有消息的固有概念。在 WebSocket 之前，使用[Comet](<https://en.wikipedia.org/wiki/Comet_(programming)>)通道可以实现端口 80 全双工通信；然而，Comet 的实现并不简单，并且由于 TCP 握手和 HTTP 标头开销，对于小消息来说效率低下。WebSocket 协议旨在在不损害 Web 安全假设的情况下解决这些问题。

WebSocket 协议规范将 ws(WebSocket) 和 wss(WebSocket Secure) 定义为两种新的[统一资源标识符](https://en.wikipedia.org/wiki/Uniform_resource_identifier)(URI) 方案[\[6\]](https://en.wikipedia.org/wiki/WebSocket#cite_note-6)，分别用于未加密和加密连接。除了方案名称和[片段](https://en.wikipedia.org/wiki/Fragment_identifier)（即#不支持），其余的 URI 组件被定义为使用[URI 通用语法](https://en.wikipedia.org/wiki/Path_segment)。[\[7\]](https://en.wikipedia.org/wiki/WebSocket#cite_note-7)
使用浏览器开发人员工具，开发人员可以检查 WebSocket 握手以及 WebSocket 帧。[\[8\]](https://en.wikipedia.org/wiki/WebSocket#cite_note-8)

# 其他文章

> <https://mp.weixin.qq.com/s/38hgQ2zDBfl09tSruiwtDA>

随着科技发展，人们需求越来越多，生活的方方面面都离不开一些实时信息。比如：疫情期间在家协同办公、疫情监控目标人的实时运动轨迹、社交中的实时消息、多玩家互动游戏、每秒瞬息万变的股市基金报价、体育实况播放、音视频聊天、视频会议、在线教育等等，都可以借用 WebSocket TCP 链接可以让数据飞起来。下面就聊一下 WebSocket 协议。

WebSocket 是 HTML5 开始提供的一种浏览器与服务器间进行全双工通讯的网络技术,一种基于 TCP 连接上进行全双工通信的协议，相对于 HTTP 这种非持久的协议来说，WebSocket 是一个持久化网络通信的协议。依靠这种技术可以实现客户端和服务器端的长连接，双向实时通信。

它不仅可以实现客户端请求服务器，同时可以允许服务端主动向客户端推送数据。是真正的双向平等对话，属于服务器推送技术的一种。在 WebSocket API 中，客户端和服务器只需要完成一次握手，两者之间就直接可以创建持久性的连接，并进行双向数据传输。

**「其他特点包括：」**

> - 建立在 TCP 协议之上，服务器端的实现比较容易。
> - 与 HTTP 协议有着良好的兼容性。默认端口也是 80 和 443，并且握手阶段采用 HTTP 协议，因此握手时不容易屏蔽，能通过各种 HTTP 代理服务器。
> - 数据格式比较轻量，性能开销小，通信高效。
> - 可以发送文本，也可以发送二进制数据。
> - 没有同源限制，客户端可以与任意服务器通信。
> - 协议标识符是 ws（如果加密，则为 wss），服务器网址就是 URL。

协议标识符是 ws（如果加密，则为 wss），服务器网址就是 URL

`ws://xxx.ayunw.cn:80/some/path wss://xxx.ayunw.cn:443/some/path`

另外客户端不只是浏览器，只要实现了 ws 或者 wss 协议的客户端 socket 都可以和服务器进行通信。

## 先说一下为什么需要 WebSocket 协议?

在 Web 应用架构中，连接由 HTTP/1.0 和 HTTP/1.1 处理。HTTP 是客户端/服务器模式中 请求一响应 所用的协议，在这种模式中，客户端(一般是浏览器)向服务器提交 HTTP 请求，服务器响应请求的资源(例如 HTML 页面)。

HTTP 是无状态的，也就是说，它将每个请求当成唯一和独立的。无状态协议具有一些优势，例如，服务器不需要保存有关会话的信息，从而不需要存储数据。但是，这也意味着在每次 HTTP 请求和响应中都会发送关于请求的冗余信息，比如使用 Cookie 进行用户状态的验证。

随着客户端和服务器之间交互的增加，HTTP 协议在客户端和服务器之间通信所需要的信息量快速增加。

从根本上讲，HTTP 还是 半双工 的协议，也就是说，在同一时刻信息的流向只能单向的：客户端向服务器发送请求(单向)，然后服务器响应请求(单向)。半双工方式的通信效率是非常低的。

同时 HTTP 协议有一个缺陷：通信只能由客户端发起。

这种单向请求的特点，注定了如果服务器有状态变化，是无法主动通知客户端的。

为了能够及时的获取服务器的变化，我们尝试过各种各样的方式：

> - 轮询(polling)：每隔一段时间，就发出一个请求，了解服务器有没有新的信息。不精准，有延时，大量无效数据交换。
> - 长轮询( long polling)：客户端向服务器请求信息，并在设定的时间段内保持连接。直到服务器有新消息响应，或者连接超时，这种技术常常称作“挂起 GET”或“搁置 POST”。占用服务器资源，相对轮询并没有优势，没有标准化。
> - 流化技术：在流化技术中，客户端发送一个请求，服务器发送并维护一个持续更新和保持打开(可以是无限或者规定的时间段)的开放响应。每当服务器有需要交付给客户端的信息时，它就更新响应。服务器从不发出完成 HTTP 响应。代理和防火墙可能缓存响应，导致信息交付的延迟增加。

上述方法提供了近乎实时的通信，但是它们也涉及 HTTP 请求和响应首标，包含了许多附加和不必要的首标数据与延迟。此外，在每一种情况下，客户端都必须等待请求返回，才能发出后续的请求，而这显著地增加了延退。同时也极大地增加了服务器的压力。

## 什么是 websocket 协议?

Websocket 其实是一个新协议，借用了 HTTP 的协议来完成一部分握手，只是为了兼容现有浏览器的握手规范而已。Websocket 是一种自然的全双工、双向、单套接字连接，解决了 HTTP 协议中不适合于实时通信的问题。

**「一个典型的 Websocket 握手如下：」**

```bash
GET /chat HTTP/1.1
Host: server.example.com
Upgrade: websocket
Connection: Upgrade
Sec-WebSocket-Key: x3JJHMbDL1EzLkh9GBhXDw==
Sec-WebSocket-Protocol: chat, superchat
Sec-WebSocket-Version: 13
Origin: http://example.com
```

其中 Websocket 的核心如下，它告诉 Apache、Nginx 等服务器：注意，我发起的是 Websocket 协议，快点帮我找到对应的助理处理而不是那个老土的 HTTP。

```bash
Upgrade: websocket
Connection: Upgrade
```

**「服务器返回如下：」**

```bash
HTTP/1.1 101 Switching Protocols
Upgrade: websocket
Connection: Upgrade
Sec-WebSocket-Accept: HSmrc0sMlYUkAGmm5OPpG2HaGWk=
Sec-WebSocket-Protocol: chat
```

至此，HTTP 已经完成它所有工作了，接下来就是完全按照 Websocket 协议进行了。

# TCP Segment 结构

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/xgaggc/1669021238963-8fd674c5-92c0-408a-9c39-a1df68ed7e24.png)
