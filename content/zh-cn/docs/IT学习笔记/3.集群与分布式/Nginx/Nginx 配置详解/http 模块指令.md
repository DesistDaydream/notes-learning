---
title: http 模块指令
---

# 概述

> 参考：
> - [org 官方文档,http core 模块](http://nginx.org/en/docs/http/ngx_http_core_module.html)
> - [官方文档,管理指南-负载均衡-HTTP 负载均衡](https://docs.nginx.com/nginx/admin-guide/load-balancer/http-load-balancer/)

http 模块及其子模块通过 `http {}` 配置环境中的指令控制行为

`http{}` 配置环境下的每个 `server{}` 指令块控制访问特定域名或者 ip 地址上对资源请求的处理。`**server{}**`** 指令块中的一个或多个 **`**location{}**`** 指令块定义了根据 URL 来处理流量的规则**

1. 比如用户访问 map.baidu.com 和 baike.baidu.com。看上去是访问了两台服务器，但是实际上，这是经过作为代理设备的 ngxin 来进行选择后的虚拟服务器。一般情况下，baike.baidu.com 与 map.baidu.com 这俩域名所解析出来的 ip 应该是同一个公网 ip(比如 123.123.123.123)(baidu 有钱用很多公网 IP 除外)。所以可以想到，用户在浏览器输入任何一个域名，访问请求都会来到 123.123.123.123，然后根据请求报文中的 Request-URL 字段中的域名与 server_name 进行配对，用户输入的 URL 中域名与哪个 server_name 相同，则该请求就会通过这个 server 来进行处理，然后根据该 server 中 location 的关键字来决定把改请求转发给哪里。

对于 `http{}` 配置环境来说，server{}、server_name、location{}、proxy_pass 是实现 7 层代理的关键指令。server_name 指定接受流量的域名，`location{}` 匹配路径，然后通过 proxy_pass 将流量代理到指定的后端。

## 简单的 http{} 配置环境示例

```nginx
http {
    access_log /dev/stdout main;

    upstream backend { # 后端配置
        server backend1.example.com;
        server backend2.example.com;
        server 192.0.0.1 backup;
    }

    server {
        server_name localhost; # 流量入口
        location / { # 流量处理
            proxy_pass http://backend;
        }
    }

    include /etc/nginx/conf.d/*.conf
}
```

# Virtual Server 基本配置

## 流量入口指令

流量入口的指令通常都定义在 `**server{} 块指令**` 中。

### [server {}](https://nginx.org/en/docs/http/ngx_http_core_module.html#server)

- 作用范围：http{}

server{} 指令块用来定义 Virtual Server

```nginx
server {
    listen [::]:80;
    server_name  "baike.baidu.com";
    location / {
        proxy_pass http://192.168.0.100:8080
    }
}
server {
    .......
}
```

下面详解的各种 简单指令 或 块指令 一般情况，都将会定义在 `server{}` 块指令中

### [listen TARGET\[PARAMETER\];](https://nginx.org/en/docs/http/ngx_http_core_module.html#listen)

- 默认值：`listen *：80 | *：8000;`

指定 Virtual Server 监听的端口，也可加上 IP:PORT

- **TARGET** # 每个 Virtual Server 的 listen 指令都会让 Nginx 监听一个 TARGET。TARGET 可以有多种格式：
  - ADDRESS:\[PORT] # 监听在指定的 IP 和端口上，ADDRESS 可以使用通配符。
  - PORT # 省略地址，即监听在所有 IP 的指定端口上。
  - UNIX:PATH # 监听在以 unix: 为前缀的 UNIX 套接字上。
- **PARAMETER** # 可以为指定的监听配置参数，多个参数以空格分割：
  - **default_server** # 将该 Virtual Server 设为默认。若客户端的请求没有匹配到任何 Virtual Server，则该请求由默认 Virtual Server 处理。
    - 注意：若没有任何 Virtual Server 配置了 listen 指令的 defautl_server 参数，那么当匹配不到 Virtual Server 时，默认的 Virtual Server 就是
  - **ssl** # 启动 SSL 模块，让此监听上接受的所有连接都应在 SSL 模式下工作
  - ......

### [server_name STRING;](https://nginx.org/en/docs/http/ngx_http_core_module.html#server_name)

- 默认值：`server_name "";`

入口名称(也就是 Virtual Server 的名字)。用来匹配一个请求 Header 中的 Host 字段。

STRING 可以是完整（精确）名称，通配符或正则表达式。通配符是一个字符串，在字符串的开始，结尾或同时包括和都包括星号（\*）。星号匹配任何字符序列。 NGINX Plus 对正则表达式使用 Perl 语法;在其前面加上波浪号（〜）。

server_name 指令是用来匹配用户在浏览器浏览网站时，输入的 域名 或者 IP:PORT 的。比如用户访问 www.baidu.com。server\_name 就可以设置为 \*.baidu.com

如果有多个域名，则可以写多个 server_name 指令，也就是说所有来自这些域名的流量都会被统一处理。

## 流量处理指令

流量处理的指令通常都在 `**location URI {}**`\*\* \*\*块指令内。

### [location \[=|~|~\*|^~\] URI {}](http://nginx.org/en/docs/http/ngx_http_core_module.html#location)

根据用户请求的 URI 进行匹配，匹配到时，此请求将被响应的 `location{}` 块指令中的指令所处理。对于用户请求的匹配优先级：

- **=** #
- **^~** #
- **~** # 正则匹配。区分大小写的匹配
- **~\*** # 正则匹配。不区分大小写
- **无符号** # 精确匹配,区分大小写,不区分大小写

> 下面详解的各种 简单指令 或 块指令 一般情况，都将会定义在 `location URL {}` 块指令中

```nginx
location / {	# 用户请求 / 目录下的文件的时候如何处理
    limit_except GET POST HEAD{
        deny all;
    }
    if http_version == 1.0 then
        return ngx.exec("/hcs_proxy_10", args)
    else
        return ngx.exec("/hcs_proxy_11", args)
    end
}
location ~ \.php$ {	#用户请求的是.php文件的时候是如何处理的
    fastcgi_pass 127.0.0.1:9000;
    .......
}
```

### [proxy_pass URL;](https://nginx.org/en/docs/http/ngx_http_proxy_module.html#proxy_pass)

代替用户把对 location 定义的请求下的 URL 交给指定的 UPSTREAM 来处理请求。该指令属于 ngx_http_proxy_modeule 模块
URL 有多种表示方法(下面的 Protocol 通常都是 http 或 https)

- **Protocol://ServerGroup/URI;** # 将流量代理到[一组服务器](https://www.yuque.com/go/doc/34075747)上。每个流量都会根据负载均衡的规则交给 upstream{} 指令块中定义的服务器。
  - 新版中，也可以省略 Protocol://，直接使用 ServerGroup 的名称即可。
- **Protocol://IP:PORT/URI;** # 将流量代理到指定的服务器上。当只有一台可用的后端服务器时可以使用这种方式，这样就不用再使用 upstream 指令块定义了
- **Protocol:unix:/PATH/TO/FILE:/URI;** # 将流量代理到本地的 UNIX-domain Socket 上。socket 的路径需要使用 `:` 包裹起来。

在[其他指令](#1l7qd)中，以 `proxy_` 开头的简单指令，都可以作为 `proxy_pass` 指令的扩充，以定义更加丰富多样的流量处理功能。

注意：

- WebSocket 代理需要特殊配置。详见[官方文档-websocket](https://nginx.org/en/docs/http/websocket.html)

### 特殊的流量处理

**fastcgi_pass 127.0.0.1:9000;** # 反向代理重定向该请求到 127.0.0.1:9000 上,Nginx 本身不支持 PHP 等语言，但是它可以通过 FastCGI 来将请求扔给某些语言或框架处理（例如 PHP、Python、Perl)
代理 fastcgi 协议的指令，当协议不是 http 协议的时候，则要使用 fastcgi 模块，注意与 proxy_pass 的区别。由于 nginx 本身并不支持动态内容的 php 等文件，需要由专门的服务器来提供，nginx 收到 .php 等请求的时候，则将该请求通过 fastcgi 协议，转发给后端能处理动态内容的服务器。比如可以在 location 中设定 .php 访问的条件，然后 {} 内中写明 fastcgi 所定义的服务器

**fastcgi_index index.php;** #

**fastcgi_param SCRIPT_FILENAME /scripts$fastcgi_script_name;** #

**fastcgi_cache_path path PATH ARGS...; **#

# 其他指令

这些指令一般都直接定义在顶层的 `http{}` 配置环境中，与 `server{}`、`upstream{}` 等块指令平级。还有一些指令是可以适用于多个指令块中的。定义在顶层的 `http{}` 配置环境中时，效果将会应用在每个 Virtual Server 中，同时也可以定义在单独的 指令块让，让指令作用于局部。

指令使用时的不成文规范：

- 通常来说，凡是作用范围包含 `location {}` 块指令的 简单指令，都直接定义在 `location{}` 块指令中。

## ngx_http_core_module 模块指令

> 代码：<https://github.com/nginx/nginx/blob/master/src/http/ngx_http_core_module.c>

### [alias PATH;](https://nginx.org/en/docs/http/ngx_http_core_module.html#alias)

用于 loation 上下文，定义 location 指令定义的路径的别名，注意与 root 指令的区别

### [client_body_in_file_only on | clean | off;](https://nginx.org/en/docs/http/ngx_http_core_module.html#client_body_in_file_only)

- 默认值：`client_body_in_file_only off;`
- 作用范围：http{}、server{}、location{}

确定 Nginx 是否应该将整个客户端请求正文保存到文件中。可以在调试期间或使用 `$request_body_file` 变量或模块 ngx_http_perl_module 的$ r-> request_body_file 方法时使用此指令。
设置为 on 时，请求处理后不会删除临时文件
clean 值将导致请求处理后留下的临时文件被删除。

### [client_header_timeout NUM;](https://nginx.org/en/docs/http/ngx_http_core_module.html#client_header_timeout)

读取 http 请求报文首部的超时时长

### [error_page CODE ... URI;](https://nginx.org/en/docs/http/ngx_http_core_module.html#error_page)

根据 http 响应状态码来指名特用的错误页面

### [ignore_invalid_headers on|off;](https://nginx.org/en/docs/http/ngx_http_core_module.html#ignore_invalid_headers)

是否忽略无效的请求头。

- 默认值：`ignore_invalid_headers on;`
- 作用范围：http{}、server{}

这里指的无效的请求头，主要是针对请求头的 key 来说，有效的请求头的 key 只能是由 英文字母、数字、连字符、下划线 这其中的 1 个或多个，而下划线的有效性，由 underscores_in_headers 指令控制。

### [keepalive_disable msie6|safari|none;](https://nginx.org/en/docs/http/ngx_http_core_module.html#keepalive_disable)

为指定类型的 User Agent(说白了就是浏览器) 禁用长连接

### [keepalive_requests NUMBER;](https://nginx.org/en/docs/http/ngx_http_core_module.html#keepalive_requests)

在一个长连接上所能够允许的最大资源数

- 默认值：`keepalive_requests 1000;`
- 作用范围：http{}、server{}、location{}

### [keepalive_timeout DURATION;](https://nginx.org/en/docs/http/ngx_http_core_module.html#keepalive_timeout)

设定长连接的超时时长为默认 75 秒

- 默认值：keepalive_timeout 75s;
- 作用范围：http{}、server{}、location{}

### [root PATH;](https://nginx.org/en/docs/http/ngx_http_core_module.html#root)

指明请求的 URL 所对应的资源所在文件系统上的起始路径。

- 作用范围：http{}、server{}、location{}

把 root 配置指令写到 `location / {} 指令块` 中，即表明当用户请求的是 / 下的资源时候，去 root 定义的本地的那个路径去找对应的资源。

### [sendfile on|off;](https://nginx.org/en/docs/http/ngx_http_core_module.html#sendfile)

开启或关闭 sendfile() 功能，即 [零拷贝](✏IT 学习笔记/📄1.操作系统/2.Kernel(内核)/6.File%20System%20 管理/10.1.零拷贝.md System 管理/10.1.零拷贝.md) 功能。

- 默认值：`sendfile off;`
- 作用范围：http{}、server{}、location{}

在此配置中，使用 SF_NODISKIO 标志调用 sendfile()，这将导致它不会在磁盘 I / O 上阻塞，而是报告该数据不在内存中。然后，nginx 通过读取一个字节来启动异步数据加载。第一次读取时，FreeBSD 内核将文件的前 128K 字节加载到内存中，尽管接下来的读取只会加载 16K 块中的数据。可以使用 read_ahead 指令更改此设置。

### [server_names_hash_bucket_size SIZE;](http://nginx.org/en/docs/http/ngx_http_core_module.html#server_names_hash_bucket_size)

设置 server_name 指定设定的服务器名称哈希表的桶容量。默认值取决于处理缓存线的大小。

- 默认值：`server_namers_hash_bucket_size 32|64|128;`
- 作用范围：http{}

### [tcp_nodelay on|off;](http://nginx.org/en/docs/http/ngx_http_core_module.html#tcp_nodelay)

是否开启长连接使用 tcp_nodelay 选项

### [underscores_in_headers on|off;](http://nginx.org/en/docs/http/ngx_http_core_module.html#underscores_in_headers)

是否允许请求头中的 key 带有下划线。

- 默认值：`underscores_in_headers off;`
- 作用范围：http{}、server{}

默认不允许，所有请求头中带有下划线的请求虽然可以被正常代理，但是其中带有下划线的请求头无法被传递到后端服务器。该指令受 ignore_invalid_headers(忽略无效请求头) 指令约束。若关闭 ignore_invalid_headers 指令，则 underscores_in_headers 指令不管如何配置都没有用。

## ngx_http_log_module 模块指令

> 代码：<https://github.com/nginx/nginx/blob/master/src/http/modules/ngx_http_log_module.c>

### [access_log PATH FORMAT \[PARAMETER\];](http://nginx.org/en/docs/http/ngx_http_log_module.html#access_log)

设置 access 日志的写入路径。

- 默认值：`access_log logs/access.log combined;`
- 作用范围：http{}、server{}、location{}

FORMAT 是 `log_format` 指令定义的日志格式名称，若不指定则默认是名为 combined 的日志格式

### [log_format NAME \[escape=default|json|none\] STRING ...;](http://nginx.org/en/docs/http/ngx_http_log_module.html#log_format)

设定 Nginx 的日志格式。

- 默认值：`log_format combined "...";`
- 作用范围：http{}

定义一个日志格式并将该格式命名为 NAME，格式名称可以在 access_log 等指令中直接引用。
STRING 就是具体的日志格式，其中可以引用一些自带的变量，类似于编程语言中的 printf 关键字效果。具体可用变量详见官方指令详解。

combined 是 Nginx 默认的日志格式名称，格式如下：

```nginx
log_format combined '$remote_addr - $remote_user [$time_local] '
                    '"$request" $status $body_bytes_sent '
                    '"$http_referer" "$http_user_agent"';
```

更多日志格式设置方法，见 [log_format 指令详解](https://www.yuque.com/go/doc/33182060)。

## ngx_http_proxy_module 模块指令

> 参考：
> - [org 官方文档，http-ngx_http_proxy_module](https://nginx.org/en/docs/http/ngx_http_proxy_module.html)
> - [GitHub 代码：nginx/nginx/src/http/modules/ngx_http_proxy_module.c](https://github.com/nginx/nginx/blob/master/src/http/modules/ngx_http_proxy_module.c)

### [proxy_pass URL;](http://nginx.org/en/docs/http/ngx_http_proxy_module.html#proxy_pass)

- 默认值：`无`
- 作用范围：location{}

该指令是 ngx_http_proxy_module 模块的核心指令，也是 http{}、stream{} 指令块中用来执行流量处理的指令。

> 参考：<https://mp.weixin.qq.com/s/D3dri6v0Tk45TOWsDb0HqQ>

根据 URL 最后有没有 `/` 会分为多种情况（现假设客户端请求 URL 为：`https://172.16.1.1/hello/world.html`）：

- 有 `/`

```nginx
location /hello/ {
    proxy_pass http://127.0.0.1/;
}
```

- 代理到 http://127.0.0.1/world.html
- 无 `/`

```nginx
location /hello/ {
    proxy_pass http://127.0.0.1;
}
```

- 代理到 URL：http://127.0.0.1/hello/world.html
- 有其他路由，且有 `/`

```nginx
location /hello/ {
    proxy_pass http://127.0.0.1/test/;
}
```

- 代理到 URL：http://127.0.0.1/test/world.html
- 有其他路由，且无 `/`

```nginx
location /hello/ {
    proxy_pass http://127.0.0.1/test;
}
```

- 代理到 URL：http://127.0.0.1/testworld.html

### [proxy_cache_path PATH ARGS...;](http://nginx.org/en/docs/http/ngx_http_proxy_module.html#proxy_cache_path)

设定代理服务缓存路径和其它参数

- 作用范围：http{}

### [proxy_http_version VERSION;](http://nginx.org/en/docs/http/ngx_http_proxy_module.html#proxy_http_version)

设置用于代理的 HTTP 协议版本。

- 默认值：`proxy_http_version 1.0;`
- 作用范围：http{}、server{}、location{}

> 建议将 1.1 版与 Keepalive 连接和 NTLM 身份验证配合使用。

### [proxy_intercept_errors on|off;](http://nginx.org/en/docs/http/ngx_http_proxy_module.html#proxy_intercept_errors)

确定是否应将代码大于或等于 300 的代理响应传递给客户端，还是应拦截并重定向到 nginx，以便使用 error_page 指令进行处理

- 作用范围：http{}、server{}、location{}

### [proxy_redirect REDIRECT REPLACEMENT;](http://nginx.org/en/docs/http/ngx_http_proxy_module.html#proxy_redirect)

修改被代理服务器的响应头中 Location 和 Refresh 字段的值。

- 默认值：`proxy_redirect default;`
- 作用范围：http{}、server{}、location{}

假如一个被代理的服务器响应头为 `Location: http://localhost:8000/two/some/uri/`。那么如果配置了如下指令：`proxy_redirect http://localhost:8000/two/ http://frontend/one/;` 之后。Nginx 响应给客户端的头变成了 `Location: http://frontend/one/some/uri/`

**EXAMPLE**

- `proxy_redirect http:// https://;`
  - 所有 3XX 跳转的 http 的请求都会被转为 https

### [proxy_set_header FIELD VALUE;](http://nginx.org/en/docs/http/ngx_http_proxy_module.html#proxy_set_header)

用来重定义发往后端服务器的请求 Header 内容。**常用指令**

- 默认值：
  - `proxy_set_header Host $proxy_host;`
  - `proxy_set_header Connection close;`
- 作用范围：http{}、server{}、location{}

**FIELD(字段)** # 指定要重新定义的请求 Header 的字段
**VALUE(值)** # Header 字段的值。可以是包含文本、变量（nginx 的内置变量）或者它们的组合。

- 注意：
  - 在 nginx 的配置文件中，如果当前模块中没有 proxy_set_header 的设置，则会从上级别继承配置。继承顺序为：http, server, location。
  - 由于 UPSTREAM 服务器收到的请求报文所含 IP 为代理服务器的 IP，那么就需要在代理服务器上配置该项，把用户 IP 暴露给 UPSTREAM 服务器
  - 该指令最常用在 `location{}` 块指令中，以便为每个路径的 HTTP 请求，都设置各自的 请求头。

### [proxy_ssl_certificate file;](http://nginx.org/en/docs/http/ngx_http_proxy_module.html#proxy_ssl_certificate)

指定 PEM 格式的证书文件，Ngxin 作为客户端向被代理的 HTTPS 服务器发起请求时，用来进行身份验证

- 作用范围：http{}、server{}、location{}

### [proxy_ssl_certificate_key file;](http://nginx.org/en/docs/http/ngx_http_proxy_module.html#proxy_ssl_certificate_key)

指定 PEM 格式的密钥，Ngxin 作为客户端向被代理的 HTTPS 服务器发起请求时，用来进行身份验证

- 作用范围：http{}、server{}、location{}

### [proxy_ssl_trusted_certificate FILE;](http://nginx.org/en/docs/http/ngx_http_proxy_module.html#proxy_ssl_trusted_certificate)

指定想要信任的 CA 证书文件，Ngxin 作为客户端向被代理的 HTTPS 服务器发起请求时，用来进行身份验证

- 作用范围：http{}、server{}、location{}

### 代理超时相关指令

#### [proxy_connect_timeout DURATION;](http://nginx.org/en/docs/http/ngx_http_proxy_module.html#proxy_connect_timeout)

与被代理服务器建立连接的超时时间。

- 默认值：`proxy_connect_timeout 60s;`
- 作用范围：http{}、server{}、location{}

注意：这个超时时间通常不应该超过 75 秒

#### [proxy_read_timeout DURATION;](http://nginx.org/en/docs/http/ngx_http_proxy_module.html#proxy_read_timeout)

从被代理服务器读取响应的超时时间

- 默认值：`proxy_read_timeout 60s;`
- 作用范围：http{}、server{}、location{}

该超时时间仅在两个连续的**读取**操作时间，而不是用于整个响应的传输。如果被代理服务器在这段时间内**未传输**任何内容，则连接将关闭。

所谓的两个连续读取操作，就是发送请求后，尝试读取响应的操作，其实就是读取 socket 中的数据。所以才被称为 等待被代理服务器响应的超时时间。

当一个请求从 Client 发送到 Nginx 后，Nginx 再转发给被代理服务器，如果被代理服务器的响应时间超过了 proxy_read_timeout，则 Nginx 将会返回给 Client 一个 **504 状态码**。

#### [proxy_send_timeout DURATION;](http://nginx.org/en/docs/http/ngx_http_proxy_module.html#proxy_send_timeout)

将请求发送到被代理服务器的超时时间。

- 默认值：`proxy_send_timeout 60s;`
- 作用范围：http{}、server{}、location{}

该超时时间仅在两个连续的**写入**操作时间，而不是用于整个响应的传输。如果被代理服务器在这段时间内**未收到**任何内容，则连接将关闭

## ngx_http_rewrite_moudle 模块指令

> 代码：<https://github.com/nginx/nginx/blob/master/src/http/modules/ngx_http_rewrite_module.c>

### [if (Condition) {...}](http://nginx.org/en/docs/http/ngx_http_rewrite_module.html#if)

用于 server 和 location 上下文中，类似于 if..else..这种编程语言

- 作用范围：server{}、location{}

Condition 是具体的匹配条件

```nginx
if ($remote_addr ~ "^(12.34|56.78)" && $http_user_agent ~* "spider") {
  return 403;
}
```

### [return CODE \[ TEXT | URL \];](http://nginx.org/en/docs/http/ngx_http_rewrite_module.html#return)

停止处理，并讲指定的状态码返回给客户端。常与 listen 指令的 default_server 参数一起使用，并指定状态码非 200，当客户端访问的域名不存在时，通过默认的 Virtual Server 处理，返回非 200 的状态码。

- 作用范围：server{}、location{}、if{}

### [rewrite RegEx Replacement \[FLAG\];](http://nginx.org/en/docs/http/ngx_http_rewrite_module.html#rewrite)

URL 重写，把 RegEx 匹配到的资源重定向到 Replacement 定义的位置

1. Flag 的用法：
   1. last，此 rewrite 规则重写完成后，不再被后面的其他 rewrite 规则进行处理，由 User Agent 重新对重写后 URL 发起新请求
   2. break，一旦此 rewrite 规则重写完成后，由 User Agent 重新对重写后的 URL 发起新请求，该新请求不再进行 rewrite 检查
   3. redirect，以 302 响应码，返回新的 URL，即在 web 界面地址栏上显示的 URL 也变了，注意跟前面两个 Flag 的区别
   4. permanent，以 301 响应码，返回新的 URL
2. EXAMPLE
   1. rewrite ^/images/(.\*.jpg)$ /imgs/$1 break; #把请求到 images 目录下的所有资源重定向到 imgs 目录下

## 其他模块指令

### [add_header NAME VALUE \[always\];](http://nginx.org/en/docs/http/ngx_http_headers_module.html#add_header)

重定义发往 client 的响应首部报文

- 作用范围：http{}、server{}、location{}

### [index FILE;](http://nginx.org/en/docs/http/ngx_http_index_module.html#index)

设定默认主页面

### [stub_status on|off](http://nginx.org/en/docs/http/ngx_http_stub_status_module.html#stub_status)

开启或关闭监控模块，仅能用于 location 上下文

# 配置示例

```nginx
http {
  server {
    location / {
      proxy_pass http://wss_svr.desistdaydream.ltd; # 转发
      proxy_http_version 1.1; # 代理所用的 http 版本设为 1.1
      proxy_set_header Host $host;
      proxy_set_header X-Real_IP $remote_addr;
      proxy_set_header X-Forwarded-For $remote_addr:$remote_port;
      proxy_set_header Upgrade $http_upgrade; # set_header表示将http协议头升级为websocket协议
      proxy_set_header Connection upgrade; # set_header表示将http协议头升级为websocket协议
    }
  }
}

```
