---
title: NAT
---

# 概述

> 参考：
> - [Wiki,Network Address Translation](https://en.wikipedia.org/wiki/Network_address_translation)
> - [RFC 4787](https://www.rfc-editor.org/rfc/rfc4787.html)
> - [RFC 5382,TCP 的 NAT 行为要求](https://www.rfc-editor.org/rfc/rfc5382.html)
> - [RFC 5508,ICMP 的 NAT 行为要求](https://www.rfc-editor.org/rfc/rfc5508.html)
> - [公众号,云原生实验室-NAT 穿透是如何工作的：技术原理及企业级实践](https://mp.weixin.qq.com/s/IzdUBvnRze4GGC6yCqmJwA)

**Network address translation(网络地址转换，简称 NAT)** 是一种通过修改数据包的 IP 报头中的网络地址信息，将一个 IP 地址空间重新映射到另一个 IP 地址空间的方法，同时它们在流量路由设备中传输。\[1] 该技术最初用作快捷方式，以避免在移动网络时重新寻址每个主机。面对 IPv4 地址耗尽，它已成为保护全球地址空间的流行和必不可少的工具。NAT 网关的一个可互联网路由的 IP 地址可用于整个专用网络。

IP 伪装是一种隐藏整个 IP 地址空间的技术，通常由私有 IP 地址组成，位于另一个 IP 地址的后面，通常是公共地址空间。必须隐藏的地址被更改为单个（公共）IP 地址作为传出 IP 数据包的“新”源地址，因此它看起来不是来自隐藏主机而是来自路由设备本身。由于这种技术的普及，以节省 IPv4 地址空间，术语 NAT 实际上已成为 IP 伪装的同义词。

由于网络地址转换修改了数据包中的 IP 地址信息，因此对 Internet 连接的质量产生严重影响，需要特别注意其实现的细节。NAT 实现在各种寻址情况下的特定行为及其对网络流量的影响方面差异很大。包含 NAT 实现的设备供应商通常不记录 NAT 行为的细节。

# NAT 实现分类

NAT 按照 **NAT 映射行为 **和 **有状态防火墙行为** 可以分为多种类型

- Full-cone
- Retricted-cone
- Port-restricted cone
- Symmetric

但对于 NAT 穿透来说根本不需要关心这么多类型，只需要看 **NAT 或者有状态防火墙是否会严格检查目标 Endpoint**，根据这个因素，可以将 NAT 分为 **Easy NAT** 和 **Hard NAT**。

- **Easy NAT** 及其变种称为 “Endpoint-Independent Mapping” (**EIM，终点无关的映射**) 这里的 Endpoint 指的是目标 Endpoint，也就是说，有状态防火墙只要看到有客户端自己发起的出向包，就会允许相应的入向包进入，**不管这个入向包是谁发进来的都可以**。
- **hard NAT** 以及变种称为 “Endpoint-Dependent Mapping”（**EDM，终点相关的映射**） 这种 NAT 会针对每个目标 Endpoint 来生成一条相应的映射关系。在这样的设备上，如果客户端向某个目标 Endpoint 发起了出向包，假设客户端的公网 IP 是 2.2.2.2，那么有状态防火墙就会打开一个端口，假设是 4242。那么只有来自该目标 Endpoint 的入向包才允许通过 2.2.2.2:4242，其他客户端一律不允许。这种 NAT 更加严格，所以叫 Hard NAT。
