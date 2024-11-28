---
title: http 模块指令
weight: 4
---

# 概述

> 参考：
>
> - [org 官方文档，http core 模块](http://nginx.org/en/docs/http/ngx_http_core_module.html)
> - [官方文档，管理指南-负载均衡-HTTP 负载均衡](https://docs.nginx.com/nginx/admin-guide/load-balancer/http-load-balancer/)

http 模块及其子模块通过 `http {}` 配置环境中的指令控制行为

`http{}` 配置环境下的每个 `server{}` 指令块控制访问特定域名或者 ip 地址上对资源请求的处理。`server{}` 指令块中的一个或多个 `location{}` 指令块定义了根据 URL 来处理流量的规则

- 比如用户访问 map.baidu.com 和 baike.baidu.com。看上去是访问了两台服务器，但是实际上，这是经过作为代理设备的 ngxin 来进行选择后的虚拟服务器。一般情况下，baike.baidu.com 与 map.baidu.com 这俩域名所解析出来的 ip 应该是同一个公网 ip(比如 123.123.123.123)(baidu 有钱用很多公网 IP 除外)。所以可以想到，用户在浏览器输入任何一个域名，访问请求都会来到 123.123.123.123，然后根据请求报文中的 Request-URL 字段中的域名与 server_name 进行配对，用户输入的 URL 中域名与哪个 server_name 相同，则该请求就会通过这个 server 来进行处理，然后根据该 server 中 location 的关键字来决定把改请求转发给哪里。

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

### server

https://nginx.org/en/docs/http/ngx_http_core_module.html#server

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

### listen

https://nginx.org/en/docs/http/ngx_http_core_module.html#listen

- 语法：`listen TARGET[PARAMETER]`
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

### server_name

https://nginx.org/en/docs/http/ngx_http_core_module.html#server_name

- 语法：`server_name NAME`
- 默认值：`server_name "";`

入口名称(也就是 Virtual Server 的名字)。用来匹配一个请求 Header 中的 Host 字段。

STRING 可以是完整（精确）名称，通配符或正则表达式。通配符是一个字符串，在字符串的开始，结尾或同时包括和都包括星号（\*）。星号匹配任何字符序列。 NGINX Plus 对正则表达式使用 Perl 语法;在其前面加上波浪号（〜）。

server_name 指令是用来匹配用户在浏览器浏览网站时，输入的 域名 或者 IP:PORT 的。比如用户访问 `www.baidu.com`，那么 `server_name` 就可以设置为 `*.baidu.com`

如果有多个域名，则可以写多个 server_name 指令，也就是说所有来自这些域名的流量都会被统一处理。

## 流量处理指令

流量处理的指令通常都在 `location URI {}` 块指令内。

### location

http://nginx.org/en/docs/http/ngx_http_core_module.html#location

- 语法：`location [=|~|~\*|^~] URI {}`
- 作用范围：server{}, location{}

location 指令让 Nginx 在 [find_config 阶段](/docs/Web/Nginx/Nginx%20处理%20HTTP%20请求的流程.md#find_config%20阶段) 根据用户请求的 URI 进行匹配，匹配到时，此请求将被响应的 `location{}` 块指令中的指令所处理。对于用户请求的匹配优先级：

- `location =`：精准匹配
- `location ^~`：带参前缀匹配
- `location ~`：正则匹配（区分大小写）
- `location ~*`：正则匹配（不区分大小写）
- `location /XXX`：普通前缀匹配，优先级低于带参数前缀匹配
- `location /`：任何没有匹配成功的，都会匹配这里处理

当接收到一个请求时，nginx会按照配置文件中server块的顺序逐个匹配 location 块：

- nginx 会将请求的 URI **根据优先级**与每个 location 块中的匹配规则进行比较，找到第一个匹配的 location 块。
- 如果找到匹配的 location 块，nginx 会按照该 location 块中的配置指令来处理请求，例如代理请求、返回静态文件等。
- 如果没有找到匹配的location块，则会使用默认的配置（即 `location /`）若没定义默认的配置，则返回404错误。

<font color="#ff0000">注意</font>：若 location 块中存在 ngx_http_rewrite_moudle 模块中的 [break](#break) 指令（不管指令单独使用，还是作为 rewrite 的标志使用），那么当遇到 break 后，则不会再查找任何 location，而是继续执行当前 location 中的其他指令，并像客户端返回请求结果。具体事例详见 [七层代理配置中的 rewrite 与 break 配置](/docs/Web/Nginx/Nginx%20配置详解/最佳实践/七层代理配置.md#rewrite%20与%20break)

下面详解的各种 简单指令 或 块指令 一般情况，都将会定义在 `location URL {}` 块指令中

```nginx
location / { # 用户请求 / 目录下的文件的时候如何处理
    limit_except GET POST HEAD{
        deny all;
    }
    if http_version == 1.0 then
        return ngx.exec("/hcs_proxy_10", args)
    else
        return ngx.exec("/hcs_proxy_11", args)
    end
}
location ~ \.php$ { #用户请求的是.php文件的时候是如何处理的
    fastcgi_pass 127.0.0.1:9000;
    .......
}
```

### proxy_pass

http://nginx.org/en/docs/http/ngx_http_proxy_module.html#proxy_pass

> 参考：
>
> - [公众号-YP小站，详解Nginx proxy_pass 使用](https://mp.weixin.qq.com/s/D3dri6v0Tk45TOWsDb0HqQ)

- 语法：`proxy_pass URL;`
- 默认值：`无`
- 作用范围：location{}

该指令是 [ngx_http_proxy_module 模块](#ngx_http_proxy_module%20模块指令)的核心指令，也是 http{}、stream{} 指令块中用来执行流量处理的指令。代替用户把对 location 定义的请求下的 URL 交给指定的 [后端服务器](/docs/Web/Nginx/Nginx%20配置详解/Nginx%20配置详解.md#后端服务器) 来处理请求。该指令属于 ngx_http_proxy_modeule 模块

URL 指定后端服务器的语法有如下几种(下面的 Protocol 通常都是 http 或 https)

- **Protocol://ServerGroup/URI;** # 将流量代理到 [一组服务器](/docs/Web/Nginx/Nginx%20配置详解/多用途模块的指令/upstream%20模块指令.md) 上。每个流量都会根据负载均衡的规则交给 upstream{} 指令块中定义的服务器。
- **Protocol://IP:PORT/URI;** # 将流量代理到指定的服务器上。当只有一台可用的后端服务器时可以使用这种方式，这样就不用再使用 upstream 指令块定义了
- **Protocol:unix:/PATH/TO/FILE:/URI;** # 将流量代理到本地的 UNIX-domain Socket 上。socket 的路径需要使用 `:` 包裹起来。

proxy_pass 指令指定的 URL 有多种，不同场景有不同的工作方式。

- URL 中 **有/无 PATH 部分**。
  - Notes: 根据 [URL 与 URI](/docs/4.数据通信/Protocol/HTTP/URL%20与%20URI.md) 中 URL 的规范可知，就算只有一个 `/` 也算有 PATH。
- 还有一些**特殊场景**

#### 工作方式

现假设客户端请求 URL 为：`https://172.16.1.1/hello/world.html`

**场景一、无 PATH 部分**

```nginx
location /hello/ {
    proxy_pass http://127.0.0.1;
}
```

- 代理到 http://127.0.0.1/hello/world.html

**场景二、有 PATH 部分**

```nginx
location /hello/ {
    proxy_pass http://127.0.0.1/;
}
```

- 代理到 http://127.0.0.1/world.html

```nginx
location /hello/ {
    proxy_pass http://127.0.0.1/test/;
}
```

- 代理到 URL：http://127.0.0.1/test/world.html

```nginx
location /hello/ {
    proxy_pass http://127.0.0.1/test;
}
```

- 代理到 URL：http://127.0.0.1/testworld.html

**场景三、特殊情况，与 rewrite 指令一起使用**

在某些情况下，无法确定请求 URI 中要替换的部分，比如使用 rewrite 指令通过正则表达式匹配路径时，不应该在 proxy_pass 的 URL 中使用 PATH

```nginx
location /name/ {
    rewrite    /name/([^/]+) /users?name=$1 break;
    proxy_pass http://127.0.0.1;
}
```

在这种情况下，指令中指定的 PATH 将被忽略，完整更改的请求 PATH 将传递到服务器。

**场景四、特殊情况，proxy_pass 的 URL 中有变量**

```nginx
location /name/ {
    proxy_pass http://127.0.0.1$request_uri;
}
```

在这种情况下，如果指令中指定了 PATH，它将按原样传递到服务器，替换原始请求 PATH。

#### 注意

- 在[其他指令](#其他指令)中，以 `proxy_` 开头的简单指令，都可以作为 `proxy_pass` 指令的扩充，以定义更加丰富多样的流量处理功能。
- WebSocket 代理需要特殊配置。详见[官方文档-websocket](https://nginx.org/en/docs/http/websocket.html)

### 特殊的流量处理

**fastcgi_pass 127.0.0.1:9000;** # 反向代理重定向该请求到 127.0.0.1:9000 上,Nginx 本身不支持 PHP 等语言，但是它可以通过 FastCGI 来将请求扔给某些语言或框架处理（例如 PHP、Python、Perl)
代理 fastcgi 协议的指令，当协议不是 http 协议的时候，则要使用 fastcgi 模块，注意与 proxy_pass 的区别。由于 nginx 本身并不支持动态内容的 php 等文件，需要由专门的服务器来提供，nginx 收到 .php 等请求的时候，则将该请求通过 fastcgi 协议，转发给后端能处理动态内容的服务器。比如可以在 location 中设定 .php 访问的条件，然后 {} 内中写明 fastcgi 所定义的服务器

**fastcgi_index index.php;** #

**fastcgi_param SCRIPT_FILENAME /scripts$fastcgi_script_name;** #

**fastcgi_cache_path path PATH ARGS...;** #

# 其他指令

这些指令一般都直接定义在顶层的 `http{}` 配置环境中，与 `server{}`、`upstream{}` 等块指令平级。还有一些指令是可以适用于多个指令块中的。定义在顶层的 `http{}` 配置环境中时，效果将会应用在每个 Virtual Server 中，同时也可以定义在单独的 指令块让，让指令作用于局部。

指令使用时的不成文规范：

- 通常来说，凡是作用范围包含 `location {}` 块指令的 简单指令，都直接定义在 `location{}` 块指令中。

## ngx_http_core_module 模块指令

> 代码：<https://github.com/nginx/nginx/blob/master/src/http/ngx_http_core_module.c>

### alias

https://nginx.org/en/docs/http/ngx_http_core_module.html#alias

- 语法：`alias PATH`

用于 loation 上下文，定义 location 指令定义的路径的别名，注意与 root 指令的区别

### client_body_in_file_only

https://nginx.org/en/docs/http/ngx_http_core_module.html#client_body_in_file_only

- 语法：`client_body_in_file_only on|clean|off`
- 默认值：`client_body_in_file_only off;`
- 作用范围：http{}、server{}、location{}

确定 Nginx 是否应该将整个客户端请求正文保存到文件中。可以在调试期间或使用 `$request_body_file` 变量或模块 ngx_http_perl_module 的$ r-> request_body_file 方法时使用此指令。

设置为 on 时，请求处理后不会删除临时文件

clean 值将导致请求处理后留下的临时文件被删除。

### client_header_timeout

https://nginx.org/en/docs/http/ngx_http_core_module.html#client_header_timeout

读取 http 请求报文首部的超时时长

- 语法：`client_header_timeout NUM`

### error_page CODE ... URI

https://nginx.org/en/docs/http/ngx_http_core_module.html#error_page

根据 http 响应状态码来指名特用的错误页面

### ignore_invalid_headers

https://nginx.org/en/docs/http/ngx_http_core_module.html#ignore_invalid_headers

是否忽略无效的请求头。

- 语法：`ignore_invalid_headers on|off`
- 默认值：`ignore_invalid_headers on;`
- 作用范围：http{}、server{}

这里指的无效的请求头，主要是针对请求头的 key 来说，有效的请求头的 key 只能是由 英文字母、数字、连字符、下划线 这其中的 1 个或多个，而下划线的有效性，由 underscores_in_headers 指令控制。

### keepalive_disable

https://nginx.org/en/docs/http/ngx_http_core_module.html#keepalive_disable

为指定类型的 User Agent(说白了就是浏览器) 禁用长连接

- 语法：`keepalive_disable msie6|safari|none`

### keepalive_requests

https://nginx.org/en/docs/http/ngx_http_core_module.html#keepalive_requests

在一个长连接上所能够允许的最大资源数

- 语法：`keepalive_requests NUMBER`
- 默认值：`keepalive_requests 1000;`
- 作用范围：http{}、server{}、location{}

### keepalive_timeout

https://nginx.org/en/docs/http/ngx_http_core_module.html#keepalive_timeout

设定长连接的超时时长为默认 75 秒

- 语法：`keepalive_timeout DURATION`
- 默认值：`keepalive_timeout 75s;`
- 作用范围：http{}、server{}、location{}

### root

https://nginx.org/en/docs/http/ngx_http_core_module.html#root

指明请求的 URL 所对应的资源所在文件系统上的起始路径。

- 语法：`root PATH`
- 作用范围：http{}、server{}、location{}

把 root 配置指令写到 `location / {} 指令块` 中，即表明当用户请求的是 / 下的资源时候，去 root 定义的本地的那个路径去找对应的资源。

### sendfile

https://nginx.org/en/docs/http/ngx_http_core_module.html#sendfile

开启或关闭 sendfile() 功能，即[DMA 与 零拷贝](/docs/1.操作系统/Kernel/Filesystem/DMA%20与%20零拷贝.md)功能。

- 语法：`sendfile on|off`
- 默认值：`sendfile off;`
- 作用范围：http{}、server{}、location{}

在此配置中，使用 SF_NODISKIO 标志调用 sendfile()，这将导致它不会在磁盘 I / O 上阻塞，而是报告该数据不在内存中。然后，nginx 通过读取一个字节来启动异步数据加载。第一次读取时，FreeBSD 内核将文件的前 128K 字节加载到内存中，尽管接下来的读取只会加载 16K 块中的数据。可以使用 read_ahead 指令更改此设置。

### server_names_hash_bucket_size

http://nginx.org/en/docs/http/ngx_http_core_module.html#server_names_hash_bucket_size

设置 server_name 指定设定的服务器名称哈希表的桶容量。默认值取决于处理缓存线的大小。

- 语法：`server_names_hash_bucket_size SIZE`
- 默认值：`server_namers_hash_bucket_size 32|64|128;`
- 作用范围：http{}

### tcp_nodelay

http://nginx.org/en/docs/http/ngx_http_core_module.html#tcp_nodelay

是否开启长连接使用 tcp_nodelay 选项

- 语法：`tcp_nodelay on|off`

### underscores_in_headers

http://nginx.org/en/docs/http/ngx_http_core_module.html#underscores_in_headers

是否允许请求头中的 key 带有下划线。

- 语法：`underscores_in_headers on|off`
- 默认值：`underscores_in_headers off;`
- 作用范围：http{}、server{}

默认不允许，所有请求头中带有下划线的请求虽然可以被正常代理，但是其中带有下划线的请求头无法被传递到后端服务器。该指令受 ignore_invalid_headers(忽略无效请求头) 指令约束。若关闭 ignore_invalid_headers 指令，则 underscores_in_headers 指令不管如何配置都没有用。

## ngx_http_log_module 模块指令

> 代码：<https://github.com/nginx/nginx/blob/master/src/http/modules/ngx_http_log_module.c>

更多日志格式设置方法，见 [log 相关模块](/docs/Web/Nginx/Nginx%20配置详解/多用途模块的指令/log%20相关模块.md)
### access_log

http://nginx.org/en/docs/http/ngx_http_log_module.html#access_log

设置 access 日志的写入路径。

- 语法：`access_log PATH FORMAT [PARAMETER]`
- 默认值：`access_log logs/access.log combined;`
- 作用范围：http{}、server{}、location{}

FORMAT 是 `log_format` 指令定义的日志格式名称，若不指定则默认是名为 combined 的日志格式

### log_format

http://nginx.org/en/docs/http/ngx_http_log_module.html#log_format

设定 Nginx 的日志格式。

- 语法：`log_format NAME [escape=default|json|none] STRING`
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

## ngx_http_proxy_module 模块指令

> 参考：
>
> - [org 官方文档，http-ngx_http_proxy_module](https://nginx.org/en/docs/http/ngx_http_proxy_module.html)
> - [GitHub 代码：nginx/nginx/src/http/modules/ngx_http_proxy_module.c](https://github.com/nginx/nginx/blob/master/src/http/modules/ngx_http_proxy_module.c)

### proxy_cache_path

http://nginx.org/en/docs/http/ngx_http_proxy_module.html#proxy_cache_path

设定代理服务缓存路径和其它参数

- 语法：`proxy_cache_path PATH ARGS`
- 作用范围：http{}

### proxy_http_version

http://nginx.org/en/docs/http/ngx_http_proxy_module.html#proxy_http_version

设置用于代理的 HTTP 协议版本。

- 语法：`proxy_http_version VERSION`
- 默认值：`proxy_http_version 1.0;`
- 作用范围：http{}、server{}、location{}

> 建议将 1.1 版与 Keepalive 连接和 NTLM 身份验证配合使用。

### proxy_intercept_errors

http://nginx.org/en/docs/http/ngx_http_proxy_module.html#proxy_intercept_errors

确定是否应将代码大于或等于 300 的代理响应传递给客户端，还是应拦截并重定向到 nginx，以便使用 error_page 指令进行处理

- 语法：`proxy_intercept_errors on|off`
- 作用范围：http{}、server{}、location{}

### proxy_redirect

http://nginx.org/en/docs/http/ngx_http_proxy_module.html#proxy_redirect

修改被代理服务器的响应头中 Location 和 Refresh 字段的值。

- 语法：`proxy_redirect REDIRECT REPLACEMENT`
- 默认值：`proxy_redirect default;`
- 作用范围：http{}、server{}、location{}

假如一个被代理的服务器响应头为 `Location: http://localhost:8000/two/some/uri/`。那么如果配置了如下指令：`proxy_redirect http://localhost:8000/two/ http://frontend/one/;` 之后。Nginx 响应给客户端的头变成了 `Location: http://frontend/one/some/uri/`

**EXAMPLE**

- `proxy_redirect http:// https://;`
  - 所有 3XX 跳转的 http 的请求都会被转为 https

### proxy_set_header

http://nginx.org/en/docs/http/ngx_http_proxy_module.html#proxy_set_header

用来重定义发往后端服务器的请求 Header 内容。这是**常用指令**

- 语法：`proxy_set_header FIELD VALUE`
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

### proxy_ssl_certificate

http://nginx.org/en/docs/http/ngx_http_proxy_module.html#proxy_ssl_certificate

指定 PEM 格式的证书文件，Ngxin 作为客户端向被代理的 HTTPS 服务器发起请求时，用来进行身份验证

- 语法：`proxy_ssl_certificate FILE`
- 作用范围：http{}、server{}、location{}

### proxy_ssl_certificate_key

http://nginx.org/en/docs/http/ngx_http_proxy_module.html#proxy_ssl_certificate_key

指定 PEM 格式的密钥，Ngxin 作为客户端向被代理的 HTTPS 服务器发起请求时，用来进行身份验证

- 语法：`proxy_ssl_certificate_key FILE`
- 作用范围：http{}、server{}、location{}

### proxy_ssl_trusted_certificate

http://nginx.org/en/docs/http/ngx_http_proxy_module.html#proxy_ssl_trusted_certificate

指定想要信任的 CA 证书文件，Ngxin 作为客户端向被代理的 HTTPS 服务器发起请求时，用来进行身份验证

- 语法：`proxy_ssl_trusted_certificate FILE`
- 作用范围：http{}、server{}、location{}

### 代理超时相关指令

#### proxy_connect_timeout

http://nginx.org/en/docs/http/ngx_http_proxy_module.html#proxy_connect_timeout

与被代理服务器建立连接的超时时间。

- 语法：`proxy_connect_timeout DURATION`
- 默认值：`proxy_connect_timeout 60s;`
- 作用范围：http{}、server{}、location{}

注意：这个超时时间通常不应该超过 75 秒

#### proxy_read_timeout

http://nginx.org/en/docs/http/ngx_http_proxy_module.html#proxy_read_timeout

从被代理服务器读取响应的超时时间

- 语法：`proxy_read_timeout DURATION`
- 默认值：`proxy_read_timeout 60s;`
- 作用范围：http{}、server{}、location{}

该超时时间仅在两个连续的**读取**操作时间，而不是用于整个响应的传输。如果被代理服务器在这段时间内**未传输**任何内容，则连接将关闭。

所谓的两个连续读取操作，就是发送请求后，尝试读取响应的操作，其实就是读取 socket 中的数据。所以才被称为 等待被代理服务器响应的超时时间。

当一个请求从 Client 发送到 Nginx 后，Nginx 再转发给被代理服务器，如果被代理服务器的响应时间超过了 proxy_read_timeout，则 Nginx 将会返回给 Client 一个 **504 状态码**。

#### proxy_send_timeout

http://nginx.org/en/docs/http/ngx_http_proxy_module.html#proxy_send_timeout

将请求发送到被代理服务器的超时时间。

- 语法：`proxy_send_timeout DURATION`
- 默认值：`proxy_send_timeout 60s;`
- 作用范围：http{}、server{}、location{}

该超时时间仅在两个连续的**写入**操作时间，而不是用于整个响应的传输。如果被代理服务器在这段时间内**未收到**任何内容，则连接将关闭

## ngx_http_rewrite_moudle 模块指令

> 代码：<https://github.com/nginx/nginx/blob/master/src/http/modules/ngx_http_rewrite_module.c>

ngx_http_rewrite_moudle 模块可以根据 [PCRE](/docs/8.通用技术/Regular%20Expression(正则表达式).md) 进行匹配和替换，实现 URL 重写、重定向等功能。还可以根据条件语句选择要执行的代理行为。

### break

https://nginx.org/en/docs/http/ngx_http_rewrite_module.html#break

- 语法：`break`
- 作用范围：server{}、location{}、if{}

使用了 break 指令后，将会停止处理 ngx_http_rewrite_module 指令的当前集合。常作为 rewrite 的选项使用，当一个 location 块中有多个 rewrite 规则时，可以在匹配到某个 rewrite 规则后，防止继续执行后续的 rewrite 规则，比如

```nginx
server {
    listen 80;
    server_name example.com;

    location / {
        rewrite ^/old-path/(.*)$ /new-path/$1 break;
        rewrite ^/(.*)$ /legacy-path/$1;
    }
}
```

在上述示例中，如果请求 URL 匹配到 `/old-path/` 开头的路径，Nginx 会将其重写为 `/new-path/` 开头的路径，并使用 `break` 选项停止后续 `rewrite` 规则的处理。这样，请求将不会再匹配第二个 `rewrite` 规则，即不会再重写为 `/legacy-path/` 开头的路径。

需要注意的是，`break` 选项只停止当前 `rewrite` 规则的处理，并不会终止整个请求处理过程。请求仍然会继续执行后续的 Nginx 配置指令，如 `location` 块中的其他指令。

### if

http://nginx.org/en/docs/http/ngx_http_rewrite_module.html#if

用于 server 和 location 上下文中，类似于 if..else..这种编程语言

- 语法：`if (Condition) {...}`
- 作用范围：server{}、location{}

Condition 是具体的匹配条件

```nginx
if ($remote_addr ~ "^(12.34|56.78)" && $http_user_agent ~* "spider") {
  return 403;
}
```

### return

http://nginx.org/en/docs/http/ngx_http_rewrite_module.html#return

停止处理，并讲指定的状态码返回给客户端。常与 listen 指令的 default_server 参数一起使用，并指定状态码非 200，当客户端访问的域名不存在时，通过默认的 Virtual Server 处理，返回非 200 的状态码。

- 语法：`return CODE [TEXT | URL]`
- 作用范围：server{}、location{}、if{}

### rewrite

http://nginx.org/en/docs/http/ngx_http_rewrite_module.html#rewrite

- 语法：`rewrite RegEx Replacement [FLAG]`
- 作用范围：server{}、location{}、if{}

rewrite 指令用于在 Nginx 的 [rewrite 阶段](/docs/Web/Nginx/Nginx%20处理%20HTTP%20请求的流程.md#rewrite%20阶段) 重写(i.e. 更改)客户端发起 HTTP 请求的 URL。根据 RegEx 的匹配规则，凡是匹配到的 URL 都改为 Replacement。**注意：rewrite 指令的执行优先级要高于 proxy_pass**。

假如指令为: `rewrite ^/prom/(.*)$ /$1`，若请求是 http://localhost/prom/abc ，则 Nginx 会将该请求改为 http://localhost/abc 后进行后续处理。若不进行 3XX 重定向。那么客户端的 URL 不变，仅仅是 Nginx 内部执行逻辑时使用的 URL 改变了。

> 在访问一个网页时，除非收到 3XX 重定向的响应，否则浏览器地址栏中的地址是不会改变的。比如 Nginx 中的 rewrite 功能，如果不使用 **redirect** 或 **permanent** 标志，那么所有的 URL 改变都是针对 Nginx 内部来说的。

**Flag**

- **last** # 此 rewrite 规则重写完成后，将会使用新的 URL 继续匹配 location。
- **break** # 与 [break](#break) 指令效果一样。此 rewrite 规则重写完成后，不会匹配新的 location，继续执行后面的 break 所在 location 中的后续指令。
- **redirect** # 以 302 响应码，返回新的 URL，即在 web 界面地址栏上显示的 URL 也变了，注意跟前面两个 Flag 的区别
- **permanent** # 以 301 响应码，返回新的 URL
- tips: 在一般情况下，如果希望重写 URL 后继续根据新的 URL 匹配 location 块，可以使用 `last` 标志；如果希望终止重写过程，直接使用重写后的 URL 处理请求，可以使用 `break` 标志。

[rewrite 与 proxy_pass](https://blog.csdn.net/weixin_41565755/article/details/120679379)

rewrite 和 proxy_pass 都可以重写整个url，区别是：（1）rewrite 重写整个 url 后，重定向的请求由浏览器发送，不常用，一般适用于访问公网其他服务器，如用于解决跨域问题；proxy_pass 重写整个 url 后，由代理服务器发起重定向请求，浏览器是无感知的，以便于访问内网和隐藏调用链；（2）rewrite 常用于重写 path，此时使用 break 和 last 也可以隐藏重定向的调用链，使用 redirect 和 permanent 则会暴露调用链;

注意：只要不是 redirect 和 permanent 这两种重定向，那么客户端在网页上看到的 URL 是不变的，比如有一个 rewrite 规则为:

```nginx
  location /prom {
      rewrite ^/prom/(.*)$ /$1 break;
      proxy_pass http://192.168.254.253:9090;
  }
```

Chrome 请求 http://192.168.254.1/prom/graph ，那么 nginx 会在 rewrite 后，向 http://192.168.254.253:9090/graph 发起请求，但是 Chrome 地址栏行看到的依然是 http://192.168.254.1/prom/graph

## 其他模块指令

### add_header

http://nginx.org/en/docs/http/ngx_http_headers_module.html#add_header

重定义发往 client 的响应首部报文

- 语法：add_header NAME VALUE [always]
- 作用范围：http{}、server{}、location{}

### index

http://nginx.org/en/docs/http/ngx_http_index_module.html#index

设定默认主页面

- 语法：index FILE

### stub_status

http://nginx.org/en/docs/http/ngx_http_stub_status_module.html#stub_status

开启或关闭监控模块，仅能用于 location 上下文

- 语法：stub_status on|off

# 配置示例

```nginx
http {
  server {
    location / { # 所有访问 / 的请求都执行下面的行为
      proxy_pass http://192.168.254.254:10000; # 代理到 http://192.168.254.254:10000
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
