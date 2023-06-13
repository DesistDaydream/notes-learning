---
title: upstream 模块指令
---

# 概述

> 参考：
> - [http 模块下的 upstream 模块](http://nginx.org/en/docs/http/ngx_http_upstream_module.html)
> - [stream 模块下的 upstream 模块](http://nginx.org/en/docs/stream/ngx_stream_upstream_module.html)

http 和 stream 模块下的 upstream 模块都是用来定义一组可以处理流量的后端服务器，称为 **Server Group(服务器组)**。

Nginx 中有一个 **Server Group(服务器组)** 的概念，Server Group 表示一组可以处理流量的后端服务器，`proxy_pass` 等指令可以直接引用 服务器组的名称，以便将流量转发到这一组后端服务器上去。并且，可以根据规则，来决定给组内每个服务器分配多少流量，还可以根据规则如何判断 服务器是否掉线，掉线如何处理 等等。

Server Group 功能通过 upstream 模块实现。而 `upstream NAME {}` 指令是一个非常通用的指令，可以作用在顶级的 `http{}` 和 `stream{}` 配置环境中。因为 ServerGroup 就是定义一组服务器以便被 Nginx 的流量处理相关指令引用，故而属于通用配置，只不过定义在 `http{}` 配置环境中的 ServerGroup 不能被 `stream{}` 配置环境中的指令引用，反之亦然。

> 注意：虽然 upstream 比较通用，但也是逻辑意义上的通用，对于 `http{}` 与 `stream{}` 来说，它们都有各自的 `upstream{}` 模块

`upstream NAME {}` 定义的名为 NAME 的一组服务器，可以被 `proxy_pass `、`fastcgi_pass`、`uwsgi_pass`、`scgi_pass`、`memcached_pass`、`grpc_pass` 这些指令直接引用，以便将流量直接代理到这组服务器上，并且可以根据一定的算法，轮流调度。

## 配置示例

```nginx
upstream backend {
    server backend1.example.com       weight=5;
    server backend2.example.com:8080;
    server unix:/tmp/backend3;

    server backup1.example.com:8080   backup;
    server backup2.example.com:8080   backup;
}

server {
    location / {
        proxy_pass http://backend;
    }
}
```

# 服务器组指令

服务器组相关的指令，一般是定义在 `upstream{} 指令块` 中。
[**upstream NAME {}**](http://nginx.org/en/docs/http/ngx_http_upstream_module.html#upstream) # 定义一组服务器。当请求代理到该 upstream 时，可以负载均衡到定义的各个 server 上。

- 作用范围：http{}、stram{}

```nginx
upstream NAME {
  ip_hash;
  server BACKEND1.EXAMPLE.COM;
  server BACKEND2.EXAMPLE.COM;
}
```

[**server ADDRESS \[PARAMETERS\];**](http://nginx.org/en/docs/http/ngx_http_upstream_module.html#server) # 定义接受流量的服务器的地址和其他参数。

- 作用范围：upstream NAME {}

**ADDRESS** # 可以指定为具有可选端口的域名或 IP 地址，或者使用 `unix:` 前缀指定 UNIX 域套接字路径。如果未指定端口，则默认 80。解析为多个 IP 地址的域名一次定义多个服务器。
**PARAMETERS** # 可以指定一个 server 的调度参数、健康检查时间间隔 等等。

- **backup** # 将 server 标记为备份服务器。
- **fail_timeout=\<TIME>** # 与 server 连接的超时时间。`默认值：10s`。
  - fail_timeout 与 max_fails 两个参数配合，就是指，当一个服务器持续 fail_timeout 时间不可用，并尝试了 max_fails 次之后依然不可用，则将该服务器从本组中提出，不再将请求调度到这个服务器上。
- **max_fails=\<INT>** # 与 server 连接超时的次数。`默认值：1`。
- **max_conns=\<INT>** # 限制到 server 的最大同时活动连接数，也就是限制 server 的并发数。`默认值：0`。默认 0 即表示不限制。
- **weight=\<INT>** # server 的权重。`默认值：1`。权重越高，就会有更多的请求调度到这个 server 上。

## EXAMPLE

```nginx
http {
  server {
      ......
  }
  upstream backend {
    server backend1.example.com weight=5;
    server backend2.example.com;
    server 192.0.0.1 backup;
  }
}
```
