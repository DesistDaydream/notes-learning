---
title: "Proxy Client"
linkTitle: "Proxy Client"
created: "2026-04-11T12:20"
weight: 100
---
> [!question]
> 很多协议，比如 [Tunneling Protocol](/docs/4.数据通信/Protocol/Tunneling%20Protocol/Tunneling%20Protocol.md)(e.g. Shadowsocks, VMess, Trojan, etc.), etc. 可能只有自己的 Server 与 Client 程序，如果我们想要在一台设备上，通过多种协议协议链接多个目标，当某个目标不可用时，可以自动切换，怎么办呢？

# 概述

> 参考：
>
> -

代理客户端就是为了解决开头提到的问题。这些客户端都是

[Clash](/docs/Web/Proxy/Clash.md) # 支持各种混淆协议 VMess, VLESS, Shadowsocks, Trojan, Snell, TUIC, Hysteria 的客户端

[VMess](/docs/4.数据通信/Protocol/Tunneling%20Protocol/VMess.md) 协议关联的客户端，i.e. V2Ray 客户端

- [Qv2ray](https://github.com/Qv2ray/Qv2ray)
    - 跨平台 V2Ray 客户端，支持 Linux、Windows、macOS，可通过插件系统支持 SSR / Trojan / Trojan-Go / NaiveProxy 等协议
- [SagerNet](https://github.com/SagerNet/SagerNet)
    - 已归档
- SagerNet 是一个基于 V2Ray 的 Android 通用代理应用。
- [V2rayN](https://github.com/2dust/v2rayN)
- V2RayN 是一个基于 V2Ray 内核的 Windows 客户端。
- [v2rayA](https://github.com/v2rayA/v2rayA)
    - 基于 web GUI 的跨平台 V2Ray 客户端，在 Linux 上支持全局透明代理，其他平台上支持系统代理。

Shadowrocket # 俗称 ”小火箭“。一开始只支持 Shadowsocks 协议，后来支持的协议逐渐变多

Sing-box # Go 语言写的

- https://github.com/sagernet/sing-box
- https://sing-box.sagernet.org/

https://github.com/drunkdream/turbo-tunnel

- https://github.com/turbo-tunnel/docs

https://github.com/turbo-tunnel/telnet-go # 这是一个用 go 实现的 `telnet` 程序，你可以把它当作一个普通的 telnet 客户端来用（访问中文 telnet 服务端可能会有乱码）。当然，它的真正用途并不在此，而是用于当 SSH 服务端不支持端口转发时建立一个 TCP 隧道。实现原理是通过将 socket 双向通信转换为对 `stdin` 和 `stdout` 的读写，而 `stderr` 则用于日志或错误信息的输出

