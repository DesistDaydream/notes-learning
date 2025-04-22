---
title: Proxy
linkTitle: Proxy
weight: 20
---

# 概述

> 参考：
>
> - [Wiki, Proxy_server](https://en.wikipedia.org/wiki/Proxy_server)

在计算机网络中，**Proxy server(代理服务器)** 是一种服务器应用程序，充当资源请求的客户端和提供该资源的服务器之间的中介。

**Proxy(代理)** 有很多种理解，还可以表示一种服务、一个概念。

Proxy 服务在整个 IT 圈子中非常常见，隧道、VPN 等等都可以看做是代理的一种。

## Forward/Reverse proxy

**Forward proxy(正向代理)** 与 **Reverse proxy(反向代理)**。

|      | 正向代理                      | 反向代理                       |
| ---- | ------------------------- | -------------------------- |
| 代理对象 | 客户端                       | 服务端                        |
| 隐藏对象 | 客户端 IP                    | 服务端 IP                     |
| 主要用途 | 突破限制、匿名访问                 | 负载均衡、安全防护、加速               |
| 典型工具 | Clash, Squid, Shadowsocks | Nginx, HAProxy, Cloudflare |

> [!Note]
> 在 [Web](/docs/Web/Web.md) 中还有一个 User-Agent 的概念，Agent 可以看作是一种代理，只不过代理形式与 Proxy 有点不太一样，Agent 更强调作为用户的代理人执行操作。虽然都是代替真实人类发起网络请求，Agent 更靠近人类。
>
> 比如我可以这么描述：DesistDaydream 通过 Chrome 这个 **Agent**，利用 Clash 这个 **Forward proxy** 访问 Google 网站，Google 网站使用 Nginx 这个 **Reverse proxy** 返回其站点的资源给我的 Agent 后，由 Agent 展现给我。

# Squid

> 参考：
>
> - [GitHub 项目，squid-cache/squid](https://github.com/squid-cache/squid)
> - [Wiki, Squid](https://en.wikipedia.org/wiki/Squid_(software))

Squid 是一款老牌的可以提供代理服务的程序。Squid 版本 1.0.0 于 1996 年 7 月发布。

在服务端安装完成后，将 /etc/squid/squid.conf 文件中的 `http_access deny all` 修改为 `http_access allow all`；之后在客户端通过 [Linux 代理配置](/docs/1.操作系统/Linux%20管理/Linux%20管理案例/Linux%20代理配置.md) 指定服务端的 3128 端口即可。

## 关联文件与配置

**/etc/squid/**

- **./squid.conf** # 主要配置文件

# 其他

[GitHub 项目，ginuerzh/gost](https://github.com/ginuerzh/gost)

- Golang 语言编写，简单隧道
- `gost -L http://:8080 -L socks5://:1080` 使用命令直接启动一个简单的代理。
  - 然后在 Shell 中配置代理即可

```bash
export http_proxy="http://${hostip}:8080"
export https_proxy="http://${hostip}:8080"
export all_proxy="sock5://${hostip}:1080"
```

[GitHub 项目，vacuityv/vacproxy](https://github.com/vacuityv/vacproxy)

- Go 语言编写，简单的 http 代理

## proxychains

项目地址: https://github.com/haad/proxychains

凡是通过 proxychains 程序运行的程序都会通过 proxychains 配置文件中设置的代理配置来发送数据包。

apt install proxychains 即可

修改配置文件

sock5 127.0.0.1 10808 # 指定本地代理服务所监听的地址

proxychains /opt/google/chrome/chrome # 即可通过代理打开 chrome 浏览器

`proxychains curl -I https://www.google.com` 会成功

# Reverse Proxy(反向代理)

> 参考：
>
> - [Wiki, Reverse proxy](https://en.wikipedia.org/wiki/Reverse_proxy)

[Nginx](/docs/Web/Nginx/Nginx.md)

[FRP](/docs/4.数据通信/Utility/FRP.md)

https://github.com/yosebyte/nodepass
