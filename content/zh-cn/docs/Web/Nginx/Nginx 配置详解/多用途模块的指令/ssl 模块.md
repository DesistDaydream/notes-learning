---
title: ssl 模块
---

# 概述

> 参考：
> 
> - [http 模块下的 ssl 模块](http://nginx.org/en/docs/http/ngx_http_ssl_module.html)
> - mail 模块下的 ssl 模块
> - [stream 模块下的 ssl 模块](http://nginx.org/en/docs/stream/ngx_stream_ssl_module.html)

ssl 模块可以作用在 http、mail、stream 模块下，称为 **ngx_http_ssl_module**、**ngx_mail_ssl_module**、**ngx_stream_ssl_module**。

ssl 模块提供了 SSL/TLS 的必要支持，可以通过 Nginx 来为后端提供 HTTP 的服务配置 SSL、也可以为普通的 4 层服务配置 SSL。

# ngx_http_ssl_module

该模块需要 [OpenSSL](/docs/7.信息安全/Crypto%20mgmt/OpenSSL/OpenSSL.md) 库的支持才可以正常使用。

ssl 模块启用后，Nginx 将可以处理 TLS/SSL 请求。当客户端发起 TLS/SSL 请求时，Nginx 中启用了 ssl 模块的 Virtual Server 作为服务端将会用配置好的证书与客户端进行认证。然后 Nginx 再作为客户端，向被代理的后端 Server 发起 TLS/SSL 请求。

注意：

- 如果一个 Virtual Server 代理的后端服务器是 HTTPS 的，那么就必须为 Virtual Server 启用 ssl 模块。因为客户端发起的请求，总是会被该 Virtual Server 重定向到 443 端口。

## SSL 指令

**ssl on | off;** # 启用 ssl 策略。

http://nginx.org/en/docs/http/ngx_http_ssl_module.html#ssl

- 作用范围：http{}、server{}

注意：

- 该指令在 1.15.0 版已过时。应该使用 `listen` 指令中的 ssl 参数。
- 如果一个端口监听了多个 server，只要有任意一个 server 启用了 ssl 策略，则其他都默认启用。

**ssl_certificate FILE;** # 为 Virtual Server 指定 PEM 格式的证书文件

http://nginx.org/en/docs/http/ngx_http_ssl_module.html#ssl_certificate

- 作用范围：http{}、server{}

**ssl_certificate_key FILE;** # 为 Virtual Server 指定 PEM 格式的密钥文件

http://nginx.org/en/docs/http/ngx_http_ssl_module.html#ssl_certificate_key

- 作用范围：http{}、server{}

**ssl_client_certificate FILE;** # 指定一个受信任的 PEM 格式 CA 证书的文件，如果启用 ssl_stapling，该文件用于验证客户端证书和 OCSP 响应。通常用于双向认证。

http://nginx.org/en/docs/http/ngx_http_ssl_module.html#ssl_client_certificate

- 作用范围：http{}、server{}
