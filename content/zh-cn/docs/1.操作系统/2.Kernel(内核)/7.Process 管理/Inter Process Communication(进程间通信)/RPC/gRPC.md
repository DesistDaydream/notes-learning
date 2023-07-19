---
title: gRPC
---

# 概述

> 参考：
> 
> - [GitHub 组织，grpc](https://github.com/grpc)
> - [官网](https://grpc.io/)

**Google Remote Procedure Calls(谷歌远程过程调用，简称 gRPC)** 是一个开源的 RPC 系统，最初于 2015 年在 Google 开发，作为下一代 RPC 基础设施 Stubby。它使用 HTTP/2 进行传输，Protocol Buffers 作为接口描述语言，并提供身份验证、双向流和流量控制、阻塞或非阻塞绑定以及取消和超时等功能。它为多种语言生成跨平台的客户端和服务器绑定。最常见的使用场景包括在微服务风格架构中连接服务，或将移动设备客户端连接到后端服务。

gRPC 对 HTTP/2 的复杂使用使得无法在浏览器中实现 gRPC 客户端，而是需要代理。

# 其他文章

[gRPC 长连接在微服务业务系统中的实践](https://mp.weixin.qq.com/s/DNHGBCZDdRjBXX0IaIZhwQ)

[公众号-Apifox，找不到好用的 gRPC 调试工具？Apifox 表示我可以！](https://mp.weixin.qq.com/s/Kt69BhGFaqYo466R0P7HZQ)