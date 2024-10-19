---
title: WebDAV
---

# 概述

> 参考：
>
> - [Wiki, WebDAV](https://en.wikipedia.org/wiki/WebDAV)
> - [RFC 4918， HTTP Extensions for Web Distributed Authoring and Versioning (WebDAV) ](https://www.rfc-editor.org/rfc/rfc4918.html)
> - <https://www.zhihu.com/question/30719209>

**Web Distributed Authoring and Versioning(Web 分布式创作和版本控制，简称 WebDAV)** 是 HTTP 的一组扩展，它允许[用户代理](https://en.wikipedia.org/wiki/User_agent)通过提供[并发控制](https://en.wikipedia.org/wiki/Concurrency_control)和[命名空间操作的设施，](https://en.wikipedia.org/wiki/Namespace)直接在 [HTTP Web 服务器](https://en.wikipedia.org/wiki/Web_server) 中协作创作内容，从而允许 Web 被视为一种 可写的、协作的 媒体，而不仅仅是一种只读媒体。WebDAV 定义在 [RFC 4918](https://datatracker.ietf.org/doc/html/rfc4918) 中

当我们使用符合 WebDAV 标准的程序部署了服务端之后，通过客户端，就可以使用 HTTP 协议访问服务端

## 应用示例

通过 WebDAV，可以将互联网上的网盘提供商，将自身的网盘，挂载到操作系统上，作为一个盘符

HTTP 协议定义了几种请求: GET, POST,PUT 等用来下载文件上传数据。WebDAV 在标准的 HTTP 协议上扩展了特有的请求方式: PROPFIND, MOVE, COPY 等。 然后用这些请求，操作 web 服务器上的磁盘(像不像一个网盘！！！)

**注意: 在 Nginx 等代理后面的 WebDAV 无法执行那些扩展的请求方式，比如 MOVE 等，实际情况是重命名时将会报错 `Dir.Rename error: DirMove MOVE call failed: Bad Gateway: 502 Bad Gateway`，可能是因为 Nginx 不支持，具体应该如何配置解决这个问题的方法还没找到。**
