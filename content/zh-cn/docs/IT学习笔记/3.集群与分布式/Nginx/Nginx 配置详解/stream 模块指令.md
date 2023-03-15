---
title: stream 模块指令
weight: 5
---

# 概述

> 参考：
> - [org 官方文档,stream core 模块](http://nginx.org/en/docs/stream/ngx_stream_core_module.html)
> - [官方文档,管理指南-负载均衡-TCP 与 UDP 负载均衡](https://docs.nginx.com/nginx/admin-guide/load-balancer/tcp-udp-load-balancer/)

stream 模块及其子模块通过 `stream {}` 配置环境中的指令控制行为

## 简单的 stream{} 配置环境示例

```nginx
stream {
    upstream stream_backend {
        least_conn;
        server backend1.example.com:12345 weight=5;
        server backend2.example.com:12345 max_fails=2 fail_timeout=30s;
        server backend3.example.com:12345 max_conns=3;
    }

    upstream dns_servers {
        least_conn;
        server 192.168.136.130:53;
        server 192.168.136.131:53;
        server 192.168.136.132:53;
    }

    server {
        listen        12345;
        proxy_pass    stream_backend;
        proxy_timeout 3s;
        proxy_connect_timeout 1s;
    }

    server {
        listen     53 udp;
        proxy_pass dns_servers;
    }

    server {
        listen     12346;
        proxy_pass backend4.example.com:12346;
    }
}
```

# Virtual Server 基本配置

## 流量入口指令

注意：流量入口的指令通常都定义在 `server{} 块指令` 中。

### [server {...}](http://nginx.org/en/docs/http/ngx_http_core_module.html#server)

定义 Virtual Server

```nginx
stream {
  server {
    listen 12345;
    # ...
  }

  server {
    listen 53 udp;
    # ...
  }
  # ...
}
```

### [listen TARGET \[PARAMETER\];](http://nginx.org/en/docs/stream/ngx_stream_core_module.html#listen)

指定 Virtual Server 监听的端口，也可加上 IP:PORT。每个 Virtual Server 的 listen 指令都会让 Nginx 监听一个 TARGET

- **TARGET** 有多种格式
  - **ADDRESS:\[PORT]** # 监听在指定的 IP 和端口上，ADDRESS 可以使用通配符。
  - **PORT** # 省略地址，即监听在所有 IP 的指定端口上。
  - **UNIX** # 监听在以 unix: 为前缀的 UNIX 套接字上。
- **PARAMETER** 可以指定是 tcp 还是 udp 等等额外的信息。

## 流量处理指令

`stream{}` 配置环境中的流量处理指令直接配置在 `server{} 指令块`中即可。与 `http{}` 配置环境不太一样，并没有 `location{}` 块指令。

### [proxy_pass ADDRESS;](http://nginx.org/en/docs/stream/ngx_stream_proxy_module.html#proxy_pass)

将流量入口进来的流量代理到指定的 ADDRESS 上。该指令是 ngx_stream_proxy_module 模块的核心指令

- 默认值：`无`
- 作用范围：server{}

ADDRESS 有多种表示方法

- **ServerGroup** # 将流量代理到[一组服务器](https://www.yuque.com/go/doc/34075747)上。每个流量都会根据负载均衡的规则交给 upstream{} 指令块中定义的服务器。
  - 新版中，也可以省略 Protocol://，直接使用 ServerGroup 的名称即可。
- **IP:PORT** # ，当只有一台可用的后端服务器时可以使用这种方式，这样就不用再使用 upstream 指令块定义了
- **unix:/PATH/TO/FILE;** # 将流量代理到本地的 UNIX-domain Socket 上

在[其他指令](/docs/IT学习笔记/3.集群与分布式/Nginx/Nginx%20 配置详解/stream%20 模块指令.md 模块指令.md)中，以 `proxy_` 开头的简单指令，都可以作为 `proxy_pass` 指令的扩充，以定义更加丰富多样的流量处理功能。

# 其他指令

这些指令一般都直接定义在顶层的 `stream{}` 配置环境中，与 `server{}`、`upstream{}` 等块指令平级。还有一些指令是可以适用于多个指令块中的。定义在顶层的 `stream{}` 配置环境中时，效果将会应用在每个 Virtual Server 中，同时也可以定义在单独的指令块中，让指令作用于局部。

## ngx_stream_proxy_module 模块指令

> 代码：<https://github.com/nginx/nginx/blob/master/src/stream/ngx_stream_proxy_module.c>

### 超时相关指令

#### [proxy_connect_timeout DURATION;](http://nginx.org/en/docs/stream/ngx_stream_proxy_module.html#proxy_connect_timeout)

与被代理服务器建立连接的超时时间

- 默认值：`proxy_connect_timeout 60s;`
- 作用范围：stream{}、server{}

# 配置示例

```nginx
stream {
    upstream konggang_dashboard {
        server 10.10.16.19:32000;
        server 10.10.16.20:32000;
    }

    server {
        listen 32000;
        proxy_pass konggang_dashboard;
        proxy_connect_timeout 10s;
    }
}
```
