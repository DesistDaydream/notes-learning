---
title: TODO
linkTitle: TODO
weight: 20
---

# 概述

> 参考：
>
> -

https://github.com/drunkdream/turbo-tunnel

- https://github.com/turbo-tunnel/docs

https://github.com/turbo-tunnel/telnet-go # 这是一个用go实现的`telnet`程序，你可以把它当作一个普通的 telnet 客户端来用（访问中文 telnet 服务端可能会有乱码）。当然，它的真正用途并不在此，而是用于当 SSH 服务端不支持端口转发时建立一个 TCP 隧道。实现原理是通过将 socket 双向通信转换为对`stdin`和`stdout`的读写，而`stderr`则用于日志或错误信息的输出。
