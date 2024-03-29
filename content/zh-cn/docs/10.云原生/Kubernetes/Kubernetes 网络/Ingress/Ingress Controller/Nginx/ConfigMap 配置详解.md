---
title: ConfigMap 配置详解
---

# 概述

> 参考：
> 
> - [官方文档，用户指南-ConfigMap](https://kubernetes.github.io/ingress-nginx/user-guide/nginx-configuration/configmap/)
> - [GitHub 代码中的可用的配置，及其默认值](https://github.com/kubernetes/ingress-nginx/blob/master/internal/ingress/controller/config/config.go)

可以通过 ConfigMap 资源来控制 Nginx Ingress Controller 的运行时行为。Nginx Ingress Controller 将会读取指定 ConfigMap 对象中的 `.data` 字段下的内容，并解析其中的内容，转换为传统 Nginx 的配置。

`.data` 字段下的内容由无数的 **Key/Value Pairs(键/值对)** 组成。绝大部分 **Key** 都会对应一个 Nginx 的 [**Directives(指令)**](/docs/Web/Nginx/Nginx%20配置详解/Nginx%20配置详解.md#Directives(指令))。Key 的 Value 就是指令的参数。假如现在有如下 ConfigMap 配置：

```yaml
data:
  map-hash-bucket-size: "128"
  ssl-protocols: SSLv2
```

这就会生成如下 Ngxin 的配置

```nginx
http {
  ······
  map_hash_bucket_size 128;
  ssl_protocols SSLv2;
  ······
}
```

# 可用的 Key 详解

下面每个 Key 的详解中，若没写对应指令，则表示这个 Key 没有对应的老式 Nginx 指令。

**enable-undersores-in-headers(BOOLEAN)** # 是否接收 key 中带有下划线的请求头。

- 默认值：`"true"`
- 对应指令：[underscores_in_headers](http://nginx.org/en/docs/http/ngx_http_core_module.html#underscores_in_headers)

**log-format-escape-json(BOOL)** # 是否为 log_format 指令开启 escape(转义) 参数

- 默认值：`"false"`
- 对应指令：[log_format 指令](http://nginx.org/en/docs/http/ngx_http_log_module.html#log_format)中的 escape 参数

**log-format-upstream(STRING)** # 设定 Nginx 的日志格式

- 默认值：`'$remote_addr - $remote_user [$time_local] "$request" $status $body_bytes_sent "$http_referer" "$http_user_agent" $request_length $request_time [$proxy_upstream_name] [$proxy_alternative_upstream_name] $upstream_addr $upstream_response_length $upstream_response_time $upstream_status $req_id'`
- 对应指令：[log_format](http://nginx.org/en/docs/http/ngx_http_log_module.html#log_format)

**use-geoip2(BOOL)** # 是否启用 geoip2 模块。

- 默认值：`false`
- 对应指令：无

该配置需要与 `--maxmind-license-key` 命令好标志配合使用。这是由于 MaxMind 已经于 [2019 年 12 月对数据库进行了大改](https://blog.maxmind.com/2019/12/18/significant-changes-to-accessing-and-using-geolite2-databases/)，需要一个 License 才可以访问数据库。所以，`--maxmind-license-key` 标志就是用来指定 License Key 的，可以创建完 MaxMind 账户后，在[此页面](https://www.maxmind.com/en/accounts/545756/license-key)创建一个 License Key。

启用 geoip2 模块后，会自动添加相关 geoip2 指令到 http{} 配置环境，详见 [nginx.tmpl 模板文件](https://github.com/kubernetes/ingress-nginx/blob/main/rootfs/etc/nginx/template/nginx.tmpl) 中的 geoip2 相关指令。

**use-forwarded-headers(BOOL)** # 是否使用 `X-Forwarded-*` 请求头

- 默认值：`false`

注意：

- 当 Nginx Ingress Controller 处于其他 7 层代理 或 负载均衡器 后面时，应为 `true`。
- 当 Nginx Ingress Controller 直接暴露在互联网上是，应为

## SSL

**ssl-redirect**(BOOLEAN) # 当具有 TLS 证书时，是否通过 301 让请求强制跳转到 HTTPS

- `默认值: true`

**force-ssl-redirect**(BOOLEAN) # 当具有 TLS 证书时，是否通过 308 让请求强制跳转到 HTTPS

- `默认值: false`