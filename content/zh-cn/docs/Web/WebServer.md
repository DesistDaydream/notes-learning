---
title: WebServer
linkTitle: WebServer
weight: 20
---

# 概述

> 参考：
>
> - [MDN，Web server](https://developer.mozilla.org/en-US/docs/Glossary/Web_server)
> - [Wiki, Web server](https://en.wikipedia.org/wiki/Web_server)

**Web server(Web 服务器)** 是用以响应静态资源的程序，可以提供 Web 服务。

在最基本的层面上，每当浏览器需要托管在 Web server 上的文件时，浏览器都会通过 HTTP 请求该文件。当请求到达正确的 Web  server 时，HTTP 服务器接受该请求，找到所请求的文档，并将其发送回浏览器（同样通过 HTTP）。 （如果服务器找不到所请求的文档，则会返回 404 响应。）

![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/web/webserver_1.png)

实现 Web server 的软/硬件

- Tomcat
- [Nginx](/docs/Web/Nginx/Nginx.md)

# Tomcat

> 参考：
>
> - [官网](https://tomcat.apache.org/)
> - [Wiki, Apache Tomcat](https://en.wikipedia.org/wiki/Apache_Tomcat)

Apache Tomcat（简称“Tomcat”）是 Jakarta Servlet、Jakarta Expression Language 和 WebSocket 技术的免费开源实现。使用 Java 开发的 HTTP Web server 环境，Java 代码也可以在其中运行。因此，它是一个 Java Web 应用程序服务器，尽管不是完整的 JEE 应用程序服务器。
