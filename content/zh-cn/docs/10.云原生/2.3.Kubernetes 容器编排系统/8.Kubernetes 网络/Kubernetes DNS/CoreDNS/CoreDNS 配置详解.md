---
title: CoreDNS 配置详解
---

# 概述

> 参考：
> - [官方配置文档](https://coredns.io/manual/toc/#configuration)

CoreDNS 的配置文件为 Corefile。
Corefile 源于 Caddy 框架的的配置文件 Caddyfile。Corefile 将会定义如下 CoreDNS 的行为：

- CoreDNS 的运行逻辑很像 Nginx，会抽象出 `server` 的概念并运行。可以同时定义多个 Server 以实现不同功能，每个 Server 主要定义下面几种行为：
  - Server 以什么协议监听在哪个端口
  - Server 负责哪个 zone 的权威 DNS 解析
  - Server 将会加载哪些插件

一个典型的最基础的 Corefile 格式如下所示：

```bash
ZONE[:PORT] {
    [PLUGIN]...
}
```

- **ZONE**# 定义 Server 的 zone。`默认值： .`
- **PORT**# 定义 Server 监听的端口。`默认值：53`。即 -dns.port 标志的值。
- **PLUGIN**# 定义 Server 要加载的[插件](https://coredns.io/plugins/)。这是可选的，但是如果不加载任何插件，那么 coredns 将为所有查询返回 SERVFAIL 。
  - 并且，不同的 Plugins 还可以定义不同的参数以改变其运行行为。

比如：

```go
. {}
```

上述配置文件表达的是：server 负责根域 `.` 的解析，监听在 53 端口，并且不使用任何插件。

# 配置示例

## 定义 server

一个最简单的配置文件可以为：

    .{}

即 server 监听 53 端口并不使用插件。
如果此时在定义其他 server，要保证监听端口不冲突；如果是在原来 server 增加 zone，则要保证 zone 之间不冲突，如：

    .    {}
    .:54 {}

另一个 server 运行于 54 端口并负责根域 `.` 的解析。

又如：

    example.org {
        whoami
    }
    org {
        whoami
    }

同一个 server 但是负责不同 zone 的解析，有不同插件链。

### 定义 Reverse Zone

跟其他 DNS 服务器类似，Corefile 也可以定义 `Reverse Zone`（反向解析 IP 地址对应的域名）：

    0.0.10.in-addr.arpa {
        whoami
    }

或者简化版本：

    10.0.0.0/24 {
        whoami
    }

可以通过 `dig` 进行反向查询：

    $ dig -x 10.0.0.1

### 使用不同的通信协议

CoreDNS 除了支持 DNS 协议，也支持 `TLS` 和 `gRPC`，即 DNS-over-TLS\[3] 和 DNS-over-gRPC 模式：

    tls://example.org:1443 {
    	#...
    }

### 添加其他的域名解析

    .:53 {
            hosts {
               123.138.66.130 xsky.xa.ehualu.it
               123.138.66.130 prometheus.xa.ehualu.it
               116.182.4.38 xsky.heb.ehualu.it
               116.182.4.38 prometheus.heb.ehualu.it
               fallthrough
            }
    }

## Forwarding(转发)

> 参考：
> - <https://coredns.io/manual/toc/#forwarding>
