---
title: Nginx 处理 HTTP 请求的流程
linkTitle: Nginx 处理 HTTP 请求的流程
weight: 20
---

# 概述

> 参考：
>
> - 原文：[Nginx 处理 HTTP 请求的 11 个阶段](https://iziyang.github.io/2020/04/12/5-nginx/)
>   - [博客园对该文章的排版更好](https://www.cnblogs.com/iziyang/p/12933565.html)

前面给大家讲了 [Nginx 是如何处理 HTTP请求头部的](https://iziyang.github.io/2020/04/08/4-nginx/)，接下来就到了真正处理 HTTP 请求的阶段了。先看下面这张图，这张图是 Nginx 处理 HTTP 请求的示意图，虽然简单，但是却很好的说明了整个过程。

1. Read Request Headers：解析请求头。
2. Identify Configuration Block：识别由哪一个 location 进行处理，匹配 URL。
3. Apply Rate Limits：判断是否限速。例如可能这个请求并发的连接数太多超过了限制，或者 QPS 太高。
4. Perform Authentication：连接控制，验证请求。例如可能根据 Referrer 头部做一些防盗链的设置，或者验证用户的权限。
5. Generate Content：生成返回给用户的响应。为了生成这个响应，做反向代理的时候可能会和上游服务（Upstream Services）进行通信，然后这个过程中还可能会有些子请求或者重定向，那么还会走一下这个过程（Internal redirects and subrequests）。
6. Response Filters：过滤返回给用户的响应。比如压缩响应，或者对图片进行处理。
7. Log：记录日志。

![](https://img2020.cnblogs.com/other/1833725/202005/1833725-20200521210215201-1729774620.png)

以上这七个步骤从整体上介绍了一下处理流程，下面还会再说一下实际的处理过程。

# Nginx 处理 HTTP 请求的 11 个阶段

下面介绍一下详细的 11 个阶段，每个阶段都可能对应着一个甚至多个 HTTP 模块，通过这样一个模块对比，我们也能够很好的理解这些模块具体是怎么样发挥作用的。

可以在 Nginx 源码 `src/http/ngx_http_core_module.h` 找到这 11 个阶段的定义

```c
typedef enum {
    NGX_HTTP_POST_READ_PHASE = 0,

    NGX_HTTP_SERVER_REWRITE_PHASE,

    NGX_HTTP_FIND_CONFIG_PHASE,
    NGX_HTTP_REWRITE_PHASE,
    NGX_HTTP_POST_REWRITE_PHASE,

    NGX_HTTP_PREACCESS_PHASE,

    NGX_HTTP_ACCESS_PHASE,
    NGX_HTTP_POST_ACCESS_PHASE,

    NGX_HTTP_PRECONTENT_PHASE,

    NGX_HTTP_CONTENT_PHASE,

    NGX_HTTP_LOG_PHASE
} ngx_http_phases;
```

![500](https://img2020.cnblogs.com/other/1833725/202005/1833725-20200521210216056-1680220103.png)

1. POST_READ：在 read 完请求的头部之后，在没有对头部做任何处理之前，想要获取到一些原始的值，就应该在这个阶段进行处理。这里面会涉及到一个 realip 模块。
2. SERVER_REWRITE：和下面的 REWRITE 阶段一样，都只有一个模块叫 rewrite 模块，一般没有第三方模块会处理这个阶段。
3. FIND_CONFIG：做 location 的匹配，暂时没有模块会用到。
4. REWRITE：对 URL 做一些处理。
5. POST_WRITE：处于 REWRITE 之后，也是暂时没有模块会在这个阶段出现。

接下来是确认用户访问权限的三个模块：

6. PREACCESS：是在 ACCESS 之前要做一些工作，例如并发连接和 QPS 需要进行限制，涉及到两个模块：limt_conn 和 limit_req
7. ACCESS：核心要解决的是用户能不能访问的问题，例如 auth_basic 是用户名和密码，access 是用户访问 IP，auth_request 根据第三方服务返回是否可以去访问。
8. POST_ACCESS：是在 ACCESS 之后会做一些事情，同样暂时没有模块会用到。

最后的三个阶段处理响应和日志：

9. PRECONTENT：在处理 CONTENT 之前会做一些事情，例如会把子请求发送给第三方的服务去处理，try_files 模块也是在这个阶段中。
10. CONTENT：这个阶段涉及到的模块就非常多了，例如 index, autoindex, concat 等都是在这个阶段生效的。
11. LOG：记录日志 access_log 模块。

以上的这些阶段都是严格按照顺序进行处理的，当然，每个阶段中各个 HTTP 模块的处理顺序也很重要，如果某个模块不把请求向下传递，后面的模块是接收不到请求的。而且每个阶段中的模块也不一定所有都要执行一遍，下面就接着讲一下各个阶段模块之间的请求顺序。

# 11 个阶段的顺序处理

如下图所示，每一个模块处理之间是有序的，那么这个顺序怎么才能得到呢？其实非常简单，在源码 ngx_module.c 中，有一个数组 `ngx_module_name`，其中包含了在编译 Nginx 的时候的 with 指令所包含的所有模块，它们之间的顺序非常关键，在数组中顺序是相反的。

```c
char *ngx_module_names[] = {
    … …
    "ngx_http_static_module",
    "ngx_http_autoindex_module",
    "ngx_http_index_module",
    "ngx_http_random_index_module",
    "ngx_http_mirror_module",
    "ngx_http_try_files_module",
    "ngx_http_auth_request_module",
    "ngx_http_auth_basic_module",
    "ngx_http_access_module",
    "ngx_http_limit_conn_module",
    "ngx_http_limit_req_module",
    "ngx_http_realip_module",
    "ngx_http_referer_module",
    "ngx_http_rewrite_module",
    "ngx_http_concat_module",
    … …
}
```

灰色部分的模块是 Nginx 的框架部分去执行处理的，第三方模块没有机会在这里得到处理。

在依次向下执行的过程中，也可能不按照这样的顺序。例如，在 access 阶段中，有一个指令叫 satisfy，它可以指示当有一个满足的时候就直接跳到下一个阶段进行处理，例如当 access 满足了，就直接跳到 try_files 模块进行处理，而不会再执行 auth_basic、auth_request 模块。

在 content 阶段中，当 index 模块执行了，就不会再执行 auto_index 模块，而是直接跳到 log 模块。

整个 11 个阶段所涉及到的模块和先后顺序如下图所示：

![](https://img2020.cnblogs.com/other/1833725/202005/1833725-20200521210217666-1083846687.jpg)

下面开始详细讲解一下各个阶段。先来看下第一个阶段 postread 阶段，顾名思义，postread 阶段是在正式处理请求之前起作用的。

# postread 阶段

postread 阶段，是 11 个阶段的第 1 个阶段，这个阶段刚刚获取到了请求的头部，还没有进行任何处理，我们可以拿到一些原始的信息。例如，拿到用户的真实 IP 地址

## 问题：如何拿到用户的真实 IP 地址？

![](https://img2020.cnblogs.com/other/1833725/202005/1833725-20200521210218190-1346456179.png)

我们知道，TCP 连接是由一个四元组构成的，在四元组中，包含了源 IP 地址。而在真实的互联网中，存在非常多的正向代理和反向代理。例如最终的用户有自己的内网 IP 地址，运营商会分配一个公网 IP，然后访问某个网站的时候，这个网站可能使用了 CDN 加速一些静态文件或图片，如果 CDN 没有命中，那么就会回源，回源的时候可能还要经过一个反向代理，例如阿里云的 SLB，然后才会到达 Nginx。

我们要拿到的地址应该是运营商给用户分配的公网 IP 地址 115.204.33.1，对这个 IP 来进行并发连接的控制或者限速，而 Nginx 拿到的却是 2.2.2.2，那么怎么才能拿到真实的用户 IP 呢？

HTTP 协议中，有两个头部可以用来获取用户 IP：

- X-Forwardex-For 是用来传递 IP 的，这个头部会把经过的节点 IP 都记录下来
- X-Real-IP：可以记录用户真实的 IP 地址，只能有一个

## 拿到真实用户 IP 后如何使用？

针对这个问题，Nginx 是基于变量来使用。

例如 binary_remote_addr、remote_addr 这样的变量，其值就是真实的 IP，这样做连接限制也就是 limit_conn 模块才有意义，这也说明了，limit_conn 模块只能在 preaccess 阶段，而不能在 postread 阶段生效。

## realip 模块

- 默认不会编译进 Nginx

  - 需要通过 `--with-http_realip_module` 启用功能
- 变量：如果还想要使用原来的 TCP 连接中的地址和端口，需要通过这两个变量保存

  - `realip_remote_addr`
  - `realip_remote_port`
- 功能

  - 修改客户端地址
- 指令

  - `set_real_ip_from`

        指定可信的地址，只有从该地址建立的连接，获取的 realip 才是可信的

  - `real_ip_header`

        指定从哪个头部取真实的 IP 地址，默认从 `X-Real-IP` 中取，如果设置从 `X-Forwarded-For` 中取，会先从最后一个 IP 开始取

  - `real_ip_recursive`

        环回地址，默认关闭，打开的时候，如果 `X-Forwarded-For` 最后一个地址与客户端地址相同，会过滤掉该地址

```vbnet
Syntax: set_real_ip_from address | CIDR | unix:;
Default: —
Context: http, server, location

Syntax: real_ip_header field | X-Real-IP | X-Forwarded-For | proxy_protocol;
Default: real_ip_header X-Real-IP;
Context: http, server, location

Syntax: real_ip_recursive on | off;
Default: real_ip_recursive off;
Context: http, server, location
```

## 实战

上面关于 `real_ip_recursive` 指令可能不太容易理解，我们来实战练习一下，先来看 `real_ip_recursive` 默认关闭的情况：

- 重新编译一个带有 realip 模块的 nginx

> 关于如何编译 Nginx，详见：[https://iziyang.github.io/2020/03/10/1-nginx/](https://iziyang.github.io/2020/03/10/1-nginx/)

```shell
# 下载 nginx 源码，在源码目录下执行
./configure --prefix=自己指定的目录 --with-http_realip_module
make
make install
```

- 然后去上一步中自己指定的 Nginx 安装目录

```nginx
#屏蔽默认的 nginx.conf 文件的 server 块内容，并添加一行
include /Users/mtdp/myproject/nginx/test_nginx/conf/example/*.conf;
```

```nginx
# 在 example 目录下建立 realip.conf，set_real_ip_from 可以设置为自己的本机 IP
server {
    listen 80;
    server_name ziyang.realip.com;
    error_log /Users/mtdp/myproject/nginx/nginx/logs/myerror.log debug;
    set_real_ip_from 192.168.0.108;
    #real_ip_header X-Real-IP;
    real_ip_recursive off;
    # real_ip_recursive on;
    real_ip_header X-Forwarded-For;

    location / {
        return 200 "Client real ip: $remote_addr\n";
    }
}
```

在上面的配置文件中，我设置了可信代理地址为本机地址，`real_ip_recursive` 为默认的 off，`real_ip_header` 设为从 `X-Forwarded-For` 中取。

- 重载配置文件

```shell
./sbin/nginx -s reload
```

- 测试响应结果

```shell
➜  test_nginx curl -H 'X-Forwarded-For: 1.1.1.1,192.168.0.108' ziyang.realip.com
Client real ip: 192.168.0.108
```

然后再来测试 `real_ip_recursive` 打开的情况：

- 配置文件中打开 `real_ip_recursive`

```nginx
server {
    listen 80;
    server_name ziyang.realip.com;
    error_log /Users/mtdp/myproject/nginx/nginx/logs/myerror.log debug;
    set_real_ip_from 192.168.0.108;
    #real_ip_header X-Real-IP;
    #real_ip_recursive off;
    real_ip_recursive on;
    real_ip_header X-Forwarded-For;

    location / {
        return 200 "Client real ip: $remote_addr\n";
    }
}
```

- 测试响应结果

```shell
➜  test_nginx curl -H 'X-Forwarded-For: 1.1.1.1,2.2.2.2,192.168.0.108' ziyang.realip.com
Client real ip: 2.2.2.2
```

所以这里面也可看出来，如果使用 `X-Forwarded-For` 获取 realip 的话，需要打开 `real_ip_recursive`，并且，realip 依赖于 `set_real_ip_from` 设置的可信地址。

那么有人可能就会问了，那直接用 `X-Real-IP` 来选取真实的 IP 地址不就好了。这是可以的，但是 `X-Real-IP` 是 Nginx 独有的，不是 RFC 规范，如果客户端与服务器之间还有其他非 Nginx 软件实现的代理，就会造成取不到 `X-Real-IP` 头部，所以这个要根据实际情况来定。

# rewrite 阶段

下面来看一下 rewrite 模块。

首先 rewrite 阶段分为两个，一个是 server_rewrite 阶段，一个是 rewrite，这两个阶段都涉及到一个 rewrite 模块，而在 rewrite 模块中，有一个 return 指令，遇到该指令就不会再向下执行，直接返回响应。

![](https://img2020.cnblogs.com/other/1833725/202005/1833725-20200521210218403-2081889023.png)

![](https://img2020.cnblogs.com/other/1833725/202005/1833725-20200521210218605-1924541596.png)

## return 指令

return 指令的语法如下：

- 返回状态码，后面跟上 body
- 返回状态码，后面跟上 URL
- 直接返回 URL

```vbnet
Syntax: return code [text];
        return code URL;
        return URL;
Default: —
Context: server, location, if
```

返回状态码包括以下几种：

- Nginx 自定义
  - 444：立刻关闭连接，用户收不到响应
- HTTP 1.0 标准
  - 301：永久重定向
  - 302：临时重定向，禁止被缓存
- HTTP 1.1 标准
  - 303：临时重定向，允许改变方法，禁止被缓存
  - 307：临时重定向，不允许改变方法，禁止被缓存
  - 308：永久重定向，不允许改变方法

### return 指令与 error_page

`error_page` 的作用大家肯定经常见到。当访问一个网站出现 404 的时候，一般不会直接出现一个 404 NOT FOUND，而是会有一个比较友好的页面，这就是 `error_page` 的功能。

```lua
Syntax: error_page code ... [=[response]] uri;
Default: —
Context: http, server, location, if in location
```

我们来看几个例子：

```nginx
1. error_page 404 /404.html;
2. error_page 500 502 503 504 /50x.html;
3. error_page 404 =200 /empty.gif;
4. error_page 404 = /404.php;
5. location / {
       error_page 404 = @fallback;
   }
   location @fallback {
       proxy_pass http://backend;
   }
6. error_page 403 http://example.com/forbidden.html;
7. error_page 404 =301 http://example.com/notfound.html;
```

那么现在就会有两个问题，大家看下下面这个配置文件：

```nginx
server {
    server_name ziyang.return.com;
    listen 80;
    root html/;
    error_page 404 /403.html;
    #return 405;
    location / {
        #return 404 "find nothing!";
    }
}
```

1. 当 server 下包含 error_page 且 location 下有 return 指令的时候，会执行哪一个呢？
2. return 指令同时出现在 server 块下和同时出现在 location 块下，它们有合并关系吗？

这两个问题我们通过实战验证一下。

### 实战

- 将上面的配置添加到配置文件 return.conf
- 在本机的 hosts 文件中绑定 ziyang.return.com 为本地的 IP 地址
- 访问一个不存在的页面

```shell
➜  test_nginx curl  ziyang.return.com/text
<html>
<head><title>403 Forbidden</title></head>
<body>
<center><h1>403 Forbidden</h1></center>
<hr><center>nginx/1.17.8</center>
</body>
</html>
```

这个时候可以看到，是 `error_page` 生效了，返回的响应是 403。

那么假如打开了 `location` 下 `return` 指令的注释呢？

- 打开 `return` 指令注释，reload 配置文件
- 重新访问页面

```vbnet
➜  test_nginx curl  ziyang.return.com/text
find nothing!%
```

这时候，`return` 指令得到了执行。也就是第一个问题，当 `server` 下包含 `error_page` 且 `location` 下有 `return` 指令的时候，会执行 `return` 指令。

下面再看一下 `server` 下的 `return` 指令和 `location` 下的 `return` 指令会执行哪一个。

- 打开 `server` 下 `return` 指令的注释，reload 配置文件
- 重新访问页面

```shell
➜  test_nginx curl  ziyang.return.com/text
<html>
<head><title>405 Not Allowed</title></head>
<body>
<center><h1>405 Not Allowed</h1></center>
<hr><center>nginx/1.17.8</center>
</body>
</html>
```

针对上面两个问题也就有了答案：

1. 当 `server` 下包含 `error_page` 且 `location` 下有 `return` 指令的时候，会执行哪一个呢？

    会执行 `location` 下的 `return` 指令。

2. `return` 指令同时出现在 `server` 块下和同时出现在 `location` 块下，它们有合并关系吗？

    没有合并关系，先遇到哪个 `return` 指令就先执行哪一个。

## rewrite 指令

`rewrite` 指令用于修改用户传入 Nginx 的 URL。来看下 `rewrite` 的指令规则：

```makefile
Syntax: rewrite regex replacement [flag];
Default: —
Context: server, location, if
```

它的功能主要有下面几点：

- 将 `regex` 指定的 URL 替换成 `replacement` 这个新的 URL
  - 可以使用正则表达式及变量提取
- 当 `replacement` 以 http:// 或者 https:// 或者 $schema 开头，则直接返回 302 重定向
- 替换后的 URL 根据 flag 指定的方式进行处理
  - last：用 `replacement` 这个 URL 进行新的 location 匹配
  - break：break 指令停止当前脚本指令的执行，等价于独立的 break 指令
  - redirect：返回 302 重定向
  - permanent：返回 301 重定向

### 指令示例

现在我们有这样的一个目录结构：

```bash
html/first/
└── 1.txt
html/second/
└── 2.txt
html/third/
└── 3.txt
```

配置文件如下所示：

```nginx
server {
    listen 80;
 server_name rewrite.ziyang.com;
 rewrite_log on;
 error_log logs/rewrite_error.log notice;

 root html/;
 location /first {
     rewrite /first(.*) /second$1 last;
     return 200 'first!\n';
    }

 location /second {
     rewrite /second(.*) /third$1;
     return 200 'second!\n';
    }

 location /third {
     return 200 'third!\n';
    }
    location /redirect1 {
     rewrite /redirect1(.*) $1 permanent;
    }

 location /redirect2 {
     rewrite /redirect2(.*) $1 redirect;
    }

 location /redirect3 {
        rewrite /redirect3(.*) http://rewrite.ziyang.com$1;
    }

 location /redirect4 {
        rewrite /redirect4(.*) http://rewrite.ziyang.com$1 permanent;
    }
}
```

那么我们的问题是：

1. return 指令 与 rewrite 指令的顺序关系？
2. 访问 /first/3.txt，/second/3.txt，/third/3.txt 分别返回的是什么？
3. 如果不携带 flag 会怎么样？

带着这三个问题，我们来实际演示一下。

### 实战

**准备工作**

- 将上面的配置添加到配置文件 rewrite.conf
- 在本机的 hosts 文件中绑定 rewrite.ziyang.com 为 127.0.0.1

**last flag**

首先访问 rewrite.ziyang.com/first/3.txt，结果如下：

```shell
➜  ~ curl rewrite.ziyang.com/first/3.txt
second!
```

为什么结果是 second! 呢？应该是 third! 呀，可能有人会有这样的疑问。实际的匹配步骤如下：

- curl rewrite.ziyang.com/first/3.txt
- 由于 `rewrite /first(.*) /second$1 last;` 这条指令的存在，last 表示使用新的 URL 进行 location 匹配，因此接下来会去匹配 second/3.txt
- 匹配到 /second 块之后，会依次执行指令，最后返回 200
- 注意，location 块中虽然也改写了 URL，但是并不会去继续匹配，因为后面没有指定 flag。

**break flag**

下面将 `rewrite /second(.*) /third$1;` 这条指令加上 break flag，`rewrite /second(.*) /third$1 break;`

继续访问 rewrite.ziyang.com/first/3.txt，结果如下：

```curl
➜  ~ curl rewrite.ziyang.com/first/3.txt
test3%
```

这时候返回的是 3.txt 文件的内容 test3。实际的匹配步骤如下：

- curl rewrite.ziyang.com/first/3.txt
- 由于 `rewrite /first(.*) /second$1 last;` 这条指令的存在，last 表示使用新的 URL 进行 location 匹配，因此接下来会去匹配 second/3.txt
- 匹配到 /second 块之后，由于 break flag 的存在，会继续匹配 rewrite 过后的 URL
- 匹配 /third location

因此，这个过程实际请求的 URL 是 rewrite.ziyang.com/third/3.txt，这样自然结果就是 test3 了。你还可以试试访问 rewrite.ziyang.com/third/2.txt 看看会返回什么。

**redirect 和 permanent flag**

配置文件中还有 4 个 location，你可以分别试着访问一下，结果是这样的：

- redirect1：返回 301
- redirect2：返回 302
- redirect3：返回 302
- redirect4：返回 301

## **rewrite** 行为记录日志

主要是一个指令 `rewrite_log`：

```nginx
Syntax: rewrite_log on | off;
Default: rewrite_log off;
Context: http, server, location, if
```

这个指令打开之后，会把 rewrite 的日志写入 logs/rewrite_error.log 日志文件中，这是请求 /first/3.txt 的日志记录：

```basic
2020/05/06 06:24:05 [notice] 86959#0: *25 "/first(.*)" matches "/first/3.txt", client: 127.0.0.1, server: rewrite.ziyang.com, request: "GET /first/3.txt HTTP/1.1", host: "rewrite.ziyang.com"
2020/05/06 06:24:05 [notice] 86959#0: *25 rewritten data: "/second/3.txt", args: "", client: 127.0.0.1, server: rewrite.ziyang.com, request: "GET /first/3.txt HTTP/1.1", host: "rewrite.ziyang.com"
2020/05/06 06:24:05 [notice] 86959#0: *25 "/second(.*)" matches "/second/3.txt", client: 127.0.0.1, server: rewrite.ziyang.com, request: "GET /first/3.txt HTTP/1.1", host: "rewrite.ziyang.com"
2020/05/06 06:24:05 [notice] 86959#0: *25 rewritten data: "/third/3.txt", args: "", client: 127.0.0.1, server: rewrite.ziyang.com, request: "GET /first/3.txt HTTP/1.1", host: "rewrite.ziyang.com"
```

## if 指令

if 指令也是在 rewrite 阶段生效的，它的语法如下所示：

```nginx
Syntax: if (condition) { ... }
Default: —
Context: server, location
```

它的规则是：

- 条件 condition 为真，则执行大括号内的指令；同时还遵循值指令的继承规则（详见我之前的文章 [Nginx 的配置指令](https://iziyang.github.io/2020/04/06/3-nginx/)）

那么 if 指令的条件表达式包含哪些内容呢？它的规则如下：

1. 检查变量为空或者值是否为 0
2. 将变量与字符串做匹配，使用 = 或 !=
3. 将变量与正则表达式做匹配
    - 大小写敏感，~ 或者 !~
    - 大小写不敏感，~*或者 !~*
4. 检查文件是否存在，使用 -f 或者 !-f
5. 检查目录是否存在，使用 -d 或者 !-d
6. 检查文件、目录、软链接是否存在，使用 -e 或者 !-e
7. 检查是否为可执行文件，使用 -x 或者 !-x

下面是一些例子：

```nginx
if ($http_user_agent ~ MSIE) { # 与变量 http_user_agent 匹配
    rewrite ^(.*)$ /msie/$1 break;
}
if ($http_cookie ~* "id=([^;]+)(?:;|$)") { # 与变量 http_cookie 匹配
    set $id $1;
}
if ($request_method = POST) { # 与变量 request_method 匹配，获取请求方法
    return 405;
}
if ($slow) { # slow 变量在 map 模块中自定义，也可以进行匹配
    limit_rate 10k;
}
if ($invalid_referer) {
    return 403;
}
```

# find_config 阶段

当经过 rewrite 模块，匹配到 URL 之后，就会进入 find_config 阶段，开始寻找 URL 对应的 [location](/docs/Web/Nginx/Nginx%20配置详解/http%20模块指令.md#location) 配置。

## location 指令

### 指令语法

还是老规矩，咱们先来看一下 location 指令的语法：

```nginx
Syntax: location [ = | ~ | ~* | ^~ ] uri { ... }
        location @name { ... }
Default: —
Context: server, location

Syntax: merge_slashes on | off;
Default: merge_slashes on;
Context: http, server
```

这里面有一个 `merge_slashes` 指令，这个指令的作用是，加入 URL 中有两个重复的 /，那么会合并为一个，这个指令默认是打开的，只有当对 URL 进行 base64 之类的编码时才需要关闭。

### 匹配规则

location 的匹配规则是仅匹配 URI，忽略参数，有下面三种大的情况：

- 前缀字符串
  - 常规匹配
  - =：精确匹配
  - ^~：匹配上后则不再进行正则表达式匹配
- 正则表达式
  - ~：大小写敏感的正则匹配
  - ~*：大小写不敏感
- 用户内部跳转的命名 location
  - @

对于这些规则刚看上去肯定是很懵的，完全不知道在说什么，下面来实战看几个例子。

### 实战

先看一下 Nginx 的配置文件：

```nginx
server {
    listen 80;
 server_name location.ziyang.com;
 error_log  logs/error.log  debug;
    #root html/;
 default_type text/plain;
 merge_slashes off;

 location ~ /Test1/$ {
     return 200 'first regular expressions match!\n';
    }
 location ~* /Test1/(\w+)$ {
     return 200 'longest regular expressions match!\n';
    }
 location ^~ /Test1/ {
     return 200 'stop regular expressions match!\n';
    }
    location /Test1/Test2 {
        return 200 'longest prefix string match!\n';
    }
    location /Test1 {
        return 200 'prefix string match!\n';
    }
 location = /Test1 {
     return 200 'exact match!\n';
    }
}
```

问题就来了，访问下面几个 URL 会分别返回什么内容呢？

```javascript
/Test1
/Test1/
/Test1/Test2
/Test1/Test2/
/test1/Test2
```

例如访问 /Test1 时，会有几个部分都匹配上：

1. 常规前缀匹配：location /Test1
2. 精确匹配：location = /Test1

访问 /Test1/ 时，也会有几个部分匹配上：

1. location ~ /Test1/$
2. location ^~ /Test1/

那么究竟会匹配哪一个呢？Nginx 其实是遵循一套规则的，如下图所示：

![](https://img2020.cnblogs.com/other/1833725/202005/1833725-20200521210219311-1832008537.jpg)

全部的前缀字符串是放置在一棵二叉树中的，Nginx 会分为两部分进行匹配：

1. 先遍历所有的前缀字符串，选取最长的一个前缀字符串，如果这个字符串是 = 的精确匹配或 ^~ 的前缀匹配，会直接使用
2. 如果第一步中没有匹配上 = 或 ^~，那么会先记住最长匹配的前缀字符串 location
3. 按照 nginx.conf 文件中的配置依次匹配正则表达式
4. 如果所有的正则表达式都没有匹配上，那么会使用最长匹配的前缀字符串

下面看下实际的响应是怎么样的：

```lua
➜  test_nginx curl location.ziyang.com/Test1
exact match!
➜  test_nginx curl location.ziyang.com/Test1/
stop regular expressions match!
➜  test_nginx curl location.ziyang.com/Test1/Test2
longest regular expressions match!
➜  test_nginx curl location.ziyang.com/Test1/Test2/
longest prefix string match!
➜  test_nginx curl location.ziyang.com/Test1/Test3
stop regular expressions match!
```

- /Test1 匹配 location = /Test1
- /Test1/ 匹配 location ^~ /Test1/
- /Test1/Test2 匹配 location ~* /Test1/(\w+)$
- /Test1/Test2/ 匹配 location /Test1/Test2
- /Test1/Test3 匹配 location ^~ /Test1/

这里面重点解释一下 /Test1/Test3 的匹配过程：

1. 遍历所有可以匹配上的前缀字符串，总共有两个
    - ^~ /Test1/
    - /Test1
2. 选取最长的前缀字符串 /Test1/，由于前面有 ^~ 禁止正则表达式匹配，因此直接使用 location ^~ /Test1/ 的规则
3. 返回 `stop regular expressions match!`

# preaccess 阶段

下面就来到了 preaccess 阶段。我们经常会遇到一个问题，就是如何限制每个客户端的并发连接数？如何限制访问频率？这些就是在 preaccess 阶段处理完成的，顾名思义，preaccess 就是在连接之前。先来看下 limit_conn 模块。

![](https://img2020.cnblogs.com/other/1833725/202005/1833725-20200521210219647-417258696.png)

## limit_conn 模块

这里面涉及到的模块是 `ngx_http_limit_conn_module`，它的基本特性如下：

- 生效阶段：`NGX_HTTP_PREACCESS_PHASE` 阶段
- 模块：`http_limit_conn_module`
- 默认编译进 Nginx，通过 `--without-http_limit_conn_module` 禁用
- 生效范围
  - 全部 worker 进程（基于共享内存）
  - 进入 preaccess 阶段前不生效
  - 限制的有效性取决于 key 的设计：依赖 postread 阶段的 realip 模块取到真实 IP

这里面有一点需要注意，就是 limit_conn key 的设计，所谓的 key 指的就是对哪个变量进行限制，通常我们取的都是用户的真实 IP。

说完了 limit_conn 的模块，再来说一下指令语法。

### 指令语法

- 定义共享内存（包括大小），以及 key 关键字

```nginx
Syntax: limit_conn_zone key zone=name:size;
Default: —
Context: http
```

- 限制并发连接数

```makefile
Syntax: limit_conn zone number;
Default: —
Context: http, server, location
```

- 限制发生时的日志级别

```nginx
Syntax: limit_conn_log_level info | notice | warn | error;
Default: limit_conn_log_level error;
Context: http, server, location
```

- 限制发生时向客户端返回的错误码

```nginx
Syntax: limit_conn_status code;
Default: limit_conn_status 503;
Context: http, server, location
```

### 实战

下面又到了实战的环节了，通过一个实际的例子来看一下以上的几个指令是怎么起作用的。

老规矩，先上配置文件：

```nginx
limit_conn_zone $binary_remote_addr zone=addr:10m;
#limit_req_zone $binary_remote_addr zone=one:10m rate=2r/m;

server {
    listen 80;
 server_name limit.ziyang.com;
 root html/;
 error_log logs/myerror.log info;
 location /{
     limit_conn_status 500;
     limit_conn_log_level  warn;
     limit_rate 50;
     limit_conn addr 1;
        #limit_req zone=one burst=3 nodelay;
        #limit_req zone=one;
    }
}
```

- 在本地的 hosts 文件中添加 limit.ziyang.com 为本机 IP

在这个配置文件中，做了两条限制，一个是 `limit_rate` 限制为 50 个字节，并发连接数 `limit_conn` 限制为 1。

```shell
➜  test_nginx curl limit.ziyang.com
```

这时候访问 limit.ziyang.com 这个站点，会发现速度非常慢，因为每秒钟只有 50 个字节。

如果再同时访问这个站点的话，则会返回 500。

我在另一个终端里面同时访问：

```shell
➜  ~ curl limit.ziyang.com
<html>
<head><title>500 Internal Server Error</title></head>
<body>
<center><h1>500 Internal Server Error</h1></center>
<hr><center>nginx/1.17.8</center>
</body>
</html>
```

可以看到，Nginx 直接返回了 500。

## limit_req 模块

在本节开头我们就提出了两个问题：

- 如何限制每个客户端的并发连接数？

- 如何限制访问频率？

第一个问题限制并发连接数的问题已经解决了，下面来看第二个问题。

这里面生效的模块是 `ngx_http_limit_req_module`，它的基本特性如下：

- 生效阶段：`NGX_HTTP_PREACCESS_PHASE` 阶段
- 模块：`http_limit_req_module`
- 默认编译进 Nginx，通过 `--without-http_limit_req_module` 禁用
- 生效算法：leaky bucket 算法
- 生效范围
  - 全部 worker 进程（基于共享内存）
  - 进入 preaccess 阶段前不生效

### leaky bucket 算法

leaky bucket 叫漏桶算法，其他用来限制请求速率的还有令牌环算法等，这里面不展开讲。

漏桶算法的原理是，先定义一个桶的大小，所有进入桶内的请求都会以恒定的速率被处理，如果请求太多超出了桶的容量，那么就会立刻返回错误。用一张图解释一下。

![](https://img2020.cnblogs.com/other/1833725/202005/1833725-20200521210219943-2053028943.png)

这张图里面，水龙头在不停地滴水，就像用户发来的请求，所有的水滴都会以恒定的速率流出去，也就是被处理。漏桶算法对于突发流量有很好的限制作用，会将所有的请求平滑的处理掉。

### 指令语法

- 定义共享内存（包括大小），以及 key 关键字和限制速率

```vbnet
Syntax: limit_req_zone key zone=name:size rate=rate ;
Default: —
Context: http
```

> rate 单位为 r/s 或者 r/m（每分钟或者每秒处理多少个请求）

- 限制并发连接数

```makefile
Syntax: limit_req zone=name [burst=number] [nodelay];
Default: —
Context: http, server, location
```

> - burst 默认为 0
> - nodelay，如果设置了这个参数，那么对于漏桶中的请求也会立刻返回错误

- 限制发生时的日志级别

```nginx
Syntax: limit_req_log_level info | notice | warn | error;
Default: limit_req_log_level error;
Context: http, server, location
```

- 限制发生时向客户端返回的错误码

```nginx
Syntax: limit_req_status code;
Default: limit_req_status 503;
Context: http, server, location
```

### 实战

在实际验证之前呢，需要注意两个问题：

- limit_req 与 limit_conn 配置同时生效时，哪个优先级高？
- nodelay 添加与否，有什么不同？

添加配置文件，这个配置文件与上一节的配置文件其实是相同的只不过需要注释一下：

```nginx
limit_conn_zone $binary_remote_addr zone=addr:10m;
limit_req_zone $binary_remote_addr zone=one:10m rate=2r/m;

server {
    listen 80;
 server_name limit.ziyang.com;
 root html/;
 error_log logs/myerror.log info;

 location /{
     limit_conn_status 500;
     limit_conn_log_level  warn;
        #limit_rate 50;
        #limit_conn addr 1;
        #limit_req zone=one burst=3 nodelay;
     limit_req zone=one;
    }
}
```

结论：在 `limit_req zone=one` 指令下，超出每分钟处理的请求数后就会立刻返回 503。

```shell
➜  test_nginx curl limit.ziyang.com
<html>
<head><title>503 Service Temporarily Unavailable</title></head>
<body>
<center><h1>503 Service Temporarily Unavailable</h1></center>
<hr><center>nginx/1.17.8</center>
</body>
</html>
```

改变一下注释的指令：

```nginx
limit_req zone=one burst=3;
#limit_req zone=one;
```

在没有添加 burst 参数时，会立刻返回错误，而加上之后，不会返回错误，而是等待请求限制解除，直到可以处理请求时再返回。

再来看一下 nodelay 参数：

```nginx
limit_req zone=one burst=3 nodelay;
```

添加了 nodelay 之后，请求在没有达到 burst 限制之前都可以立刻被处理并返回，超出了 burst 限制之后，才会返回 503。

现在可以回答一下刚开始提出的两个问题：

- limit_req 与 limit_conn 配置同时生效时，哪个优先级高？
  - limit_req 在 limit_conn 处理之前，因此是 limit_req 会生效
- nodelay 添加与否，有什么不同？
  - 不添加 nodelay，请求会等待，直到能够处理请求；添加 nodelay，在不超出 burst 的限制的情况下会立刻处理并返回，超出限制则会返回 503。

# access 阶段

经过 preaccess 阶段对用户的限流之后，就到了 access 阶段。

## access 模块

这里面涉及到的模块是 `ngx_http_access_module`，它的基本特性如下：

- 生效阶段：`NGX_HTTP_ACCESS_PHASE` 阶段
- 模块：`http_access_module`
- 默认编译进 Nginx，通过 `--without-http_access_module` 禁用
- 生效范围
  - 进入 access 阶段前不生效

### 指令语法

```nginx
Syntax: allow address | CIDR | unix: | all;
Default: —
Context: http, server, location, limit_except

Syntax: deny address | CIDR | unix: | all;
Default: —
Context: http, server, location, limit_except
```

access 模块提供了两条指令 `allow` 和 `deny`，来看几个例子：

```nginx
location / {
    deny 192.168.1.1;
    allow 192.168.1.0/24;
    allow 10.1.1.0/16;
    allow 2001:0db8::/32;
    deny all;
}
```

对于用户访问来说，这些指令是顺序执行的，当满足了一条之后，就不会再向下执行。这个模块比较简单，我们这里不做实战演练了。

## auth_basic 模块

auth_basic 模块是用作用户认证的，当开启了这个模块之后，我们通过浏览器访问网站时，就会返回一个 401 Unauthorized，当然这个 401 用户不会看见，浏览器会弹出一个对话框要求输入用户名和密码。这个模块使用的是 RFC2617 中的定义。

### 指令语法

- 基于 HTTP Basic Authutication 协议进行用户密码的认证
- 默认编译进 Nginx
  - --without-http_auth_basic_module
  - disable ngx_http_auth_basic_module

```nginx
Syntax: auth_basic string | off;
Default: auth_basic off;
Context: http, server, location, limit_except

Syntax: auth_basic_user_file file;
Default: —
Context: http, server, location, limit_except
```

这里面我们会用到一个工具叫 htpasswd，这个工具可以用来生成密码文件，而 `auth_basic_user_file` 就依赖这个密码文件。

> htpasswd 依赖安装包 httpd-tools

生成密码的命令为：

```shell
htpasswd –c file –b user pass
```

生成的密码文件的格式为：

```makefile
# comment
name1:password1
name2:password2:comment
name3:password3
```

### 实战

- 在 example 目录下生成密码文件 auth.pass

```python
htpasswd -bc auth.pass ziyang 123456
```

- 添加配置文件

```nginx
server {
 server_name access.ziyang.com;
    listen 80;
 error_log  logs/error.log  debug;
 default_type text/plain;
 location /auth_basic {
     satisfy any;
     auth_basic "test auth_basic";
     auth_basic_user_file example/auth.pass;
     deny all;
    }
}
```

- 重载 Nginx 配置文件
- 在 /etc/hosts 文件中添加 access.ziyang.com

这时候访问 access.ziyang.com 就会弹出对话框，提示输入密码：

![](https://img2020.cnblogs.com/other/1833725/202005/1833725-20200521210220134-783796093.jpg)

## auth_request 模块

- 功能：向上游的服务转发请求，若上游服务返回的响应码是 2xx，则继续执行，若上游服务返回的响应码是 2xx，则继续执行，若上游服务返回的是 401 或者 403，则将响应返回给客户端
- 原理：收到请求后，生成子请求，通过反向代理技术把请求传递给上游服务
- 默认未编译进 Nginx，需要通过 --with-http_auth_request_module 编译进去

### 指令语法

```nginx
Syntax: auth_request uri | off;
Default: auth_request off;
Context: http, server, location

Syntax: auth_request_set $variable value;
Default: —
Context: http, server, location
```

### 实战

- 在上一个配置文件中添加以下内容

```nginx
server {
 server_name access.ziyang.com;
    listen 80;
 error_log  logs/error.log  debug;
    #root html/;
 default_type text/plain;
 location /auth_basic {
     satisfy any;
     auth_basic "test auth_basic";
     auth_basic_user_file example/auth.pass;
     deny all;
    }
 location / {
     auth_request /test_auth;
    }
 location = /test_auth {
     proxy_pass http://127.0.0.1:8090/auth_upstream;
     proxy_pass_request_body off;
     proxy_set_header Content-Length "";
     proxy_set_header X-Original-URI $request_uri;
    }
}
```

- 这个配置文件中，/ 路径下会将请求转发到另外一个服务中去，可以用 nginx 再搭建一个服务
- 如果这个服务返回 2xx，那么鉴权成功，如果返回 401 或 403 则鉴权失败

## 限制所有 access 阶段模块的 satisfy 指令

### 指令语法

```nginx
Syntax: satisfy all | any;
Default: satisfy all;
Context: http, server, location
```

`satisfy` 指令有两个值一个是 all，一个是 any，这个模块对 acces 阶段的三个模块都生效：

- access 模块
- auth_basic 模块
- auth_request 模块
- 其他模块

如果 `satisfy` 指令的值是 all 的话，就表示必须所有 access 阶段的模块都要执行，都通过了才会放行；值是 any 的话，表示有任意一个模块得到执行即可。

下面有几个问题可以加深一下理解：

1. 如果有 return 指令，access 阶段会生效吗？

    return 指令属于 rewrite 阶段，在 access 阶段之前，因此不会生效。

2. 多个 access 模块的顺序有影响吗？

    ```undefined
    ngx_http_auth_request_module,
    ngx_http_auth_basic_module,
    ngx_http_access_module,
    ```

    有影响

3. 输对密码，下面可以访问到文件吗？

    ```nginx
    location /{
     satisfy any;
     auth_basic "test auth_basic";
     auth_basic_user_file examples/auth.pass;
     deny all;
    }
    ```

    可以访问到，因为 `satisfy` 的值是 any，因此只要有模块满足，即可放行。

4. 如果把 deny all 提到 auth_basic 之前呢？

    依然可以，因为各个模块执行顺序和指令的顺序无关。

5. 如果改为 allow all，有机会输入密码吗？

    没有机会，因为 allow all 是 access 模块，先于 auth_basic 模块执行。

# precontent 阶段

讲到了这里，我们再来回顾一下 Nginx 处理 HTTP 请求的 11 个阶段：

![](https://img2020.cnblogs.com/other/1833725/202005/1833725-20200521210216056-1680220103.png)

现在我们已经来到了 precontent 阶段，这个阶段只有 try_files 这一个指令。

## try_files 模块

### 指令语法

```nginx
Syntax: try_files file ... uri;
        try_files file ... =code;
Default: —
Context: server, location
```

- 模块：`ngx_http_try_files_module` 模块
- 依次试图访问多个 URL 对应的文件（由 root 或者 alias 指令指定），当文件存在时，直接返回文件内容，如果所有文件都不存在，则按照最后一个 URL 结果或者 code 返回

### 实战

下面我们实际看一个例子：

```nginx
server {
 server_name tryfiles.ziyang.com;
 listen 80;
 error_log  logs/myerror.log  info;
 root html/;
 default_type text/plain;
 location /first {
     try_files /system/maintenance.html
            $uri $uri/index.html $uri.html
            @lasturl;
    }
 location @lasturl {
     return 200 'lasturl!\n';
    }
 location /second {
     try_files $uri $uri/index.html $uri.html =404;
    }
}
```

结果如下：

- 访问 /first 实际上到了 lasturl，然后返回 200
- 访问 /second 则返回了 404

这两个结果都与配置文件是一致的。

```shell
➜  test_nginx curl tryfiles.ziyang.com/second
<html>
<head><title>404 Not Found</title></head>
<body>
<center><h1>404 Not Found</h1></center>
<hr><center>nginx/1.17.8</center>
</body>
</html>
➜  test_nginx curl tryfiles.ziyang.com/first
lasturl!
```

## mirror 模块

mirror 模块可以实时拷贝流量，这对于需要同时访问多个环境的请求是非常有用的。

### 指令语法

- 模块：`ngx_http_mirror_module` 模块，默认编译进 Nginx
  - 通过 --without-http_mirror_module 移除模块
- 功能：处理请求时，生成子请求访问其他服务，对子请求的返回值不做处理

```nginx
Syntax: mirror uri | off;
Default: mirror off;
Context: http, server, location

Syntax: mirror_request_body on | off;
Default: mirror_request_body on;
Context: http, server, location
```

### 实战

- 配置文件如下所示，需要再开启另外一个 Nginx 来接收请求

```nginx
server {
    server_name mirror.ziyang.com;
    listen 8001;
    error_log logs/error_log debug;
    location / {
        mirror /mirror;
        mirror_request_body off;
    }
    location = /mirror {
        internal;
        proxy_pass http://127.0.0.1:10020$request_uri;
        proxy_pass_request_body off;
        proxy_set_header Content-Length "";
        proxy_set_header X-Original-URI $request_uri;
    }
}
```

- 在 access.log 文件中可以看到有请求记录日志

# content 阶段

下面开始就到了 content 阶段，先来看 content 阶段的 static 模块，虽然这是位于 content 阶段的最后一个处理模块，但是这里先来介绍它。

![](https://img2020.cnblogs.com/other/1833725/202005/1833725-20200521210220966-1780856001.png)

## static 模块

### root 和 alias 指令

先来一下 root 和 alias 这两个指令，这两个指令都是用来映射文件路径的。

```nginx
Syntax: alias path;
Default: —
Context: location
```

```nginx
Syntax: root path;
Default: root html;
Context: http, server, location, if in location
```

- 功能：将 URL 映射为文件路径，以返回静态文件内容
- 差别：root 会将完整 URL 映射进文件路径中，alias 只会将 location 后的 URL 映射到文件路径

### 实战

下面来看一个问题：

现在有一个文件路径：

```shell
html/first/
└── 1.txt
```

配置文件如下所示：

```nginx
server {
 server_name static.ziyang.com;
    listen 80;
 error_log  logs/myerror.log  info;
 location /root {
     root html;
    }
 location /alias {
        alias html;
    }
 location ~ /root/(\w+\.txt) {
     root html/first/$1;
    }
 location ~ /alias/(\w+\.txt) {
     alias html/first/$1;
    }
 location  /RealPath/ {
     alias html/realpath/;
        return 200 '$request_filename:$document_root:$realpath_root\n';
    }
}
```

那么访问以下 URL 会得到什么响应呢？

```bash
/root
/alias
/root/1.txt
/alias/1.txt
```

```shell
➜  test_nginx curl static.ziyang.com/alias/1.txt
test1%
➜  test_nginx curl static.ziyang.com/alias/
<!DOCTYPE html>
<html>
<head>
<title>Welcome to nginx!</title>
...
➜  test_nginx curl static.ziyang.com/root/
<html>
<head><title>404 Not Found</title></head>
<body>
<center><h1>404 Not Found</h1></center>
<hr><center>nginx/1.17.8</center>
</body>
</html>
➜  test_nginx curl static.ziyang.com/root/1.txt
<html>
<head><title>404 Not Found</title></head>
<body>
<center><h1>404 Not Found</h1></center>
<hr><center>nginx/1.17.8</center>
</body>
</html>
```

访问这四个路径分别得到的结果是：

- /root：404
- /alias：200
- /root/1.txt：404
- /alias/1.txt：200

这是为什么呢？是因为，root 在映射 URL 时，会把 location 中的路径也加进去，也就是：

- `static.ziyang.com/root/` 实际访问的是 `html/root/`
- `static.ziyang.com/root/1.txt` 实际是 `html/first/1.txt/root/1.txt`
- `static.ziyang.com/alias/` 实际上是正确访问到了 `html` 文件夹，由于后面有 `/` 的存在，因此实际访问的是 `html/index.html`
- `static.ziyang.com/alias/1.txt` 实际访问的是 `html/first/1.txt`，文件存在

### 三个相关变量

还是上面的配置文件：

```nginx
location  /RealPath/ {
 alias html/realpath/;
    return 200 '$request_filename:$document_root:$realpath_root\n';
}
```

这里有一个问题，在访问 `/RealPath/1.txt` 时，这三个变量的值各为多少？

为了解答这个问题，我们先来解释三个变量：

- request_filename：待访问文件的完整路径
- document_root：由 URI 和 root/alias 指令生成的文件夹路径（可能包含软链接的路径）
- realpath_root：将 document_root 中的软链接替换成真实路径

为了验证这三个变量，在 html 目录下建立一个软链接指向 first 文件夹：

```shell
ln -s first realpath
```

```ruby
➜  html curl static.ziyang.com/realpath/1.txt
/Users/mtdp/myproject/nginx/test_nginx/html/realpath/1.txt:/Users/mtdp/myproject/nginx/test_nginx/html/realpath/:/Users/mtdp/myproject/nginx/test_nginx/html/first
```

可以看出来，三个路径分别是：

- /Users/mtdp/myproject/nginx/test_nginx/html/realpath/1.txt
- /Users/mtdp/myproject/nginx/test_nginx/html/realpath/
- /Users/mtdp/myproject/nginx/test_nginx/html/first

还有其他的一些配置指令，例如：

**静态文件返回时的 Content-Type**

```nginx
Syntax: types { ... }
Default: types { text/html html; image/gif gif; image/jpeg jpg; }
Context: http, server, location

Syntax: default_type mime-type;
Default: default_type text/plain;
Context: http, server, location

Syntax: types_hash_bucket_size size;
Default: types_hash_bucket_size 64;
Context: http, server, location

Syntax: types_hash_max_size size;
Default: types_hash_max_size 1024;
Context: http, server, location
```

**未找到文件时的错误日志**

```nginx
Syntax: log_not_found on | off;
Default: log_not_found on;
Context: http, server, location
```

在生产环境中，经常可能会有找不到文件的情况，错误日志中就会打印出来：

```vhdl
[error] 10156#0: *10723 open() "/html/first/2.txt/root/2.txt" failed (2: No such file or directory)
```

如果不想记录日志，可以关掉。

### 重定向跳转的域名

现在有另外一个问题，当我们访问目录时最后没有带 `/`，static 模块会返回 301 重定向，那么这个规则是怎么定义的呢，看下面三个指令：

```nginx
# 该指令决定重定向时的域名，可以决定返回哪个域名
Syntax: server_name_in_redirect on | off;
Default: server_name_in_redirect off;
Context: http, server, location
# 该指令决定重定向时的端口
Syntax: port_in_redirect on | off;
Default: port_in_redirect on;
Context: http, server, location
# 该指令决定是否填域名，默认是打开的，也就是返回绝对路径
Syntax: absolute_redirect on | off;
Default: absolute_redirect on;
Context: http, server, location
```

这三个指令的实际用法来实战演示一下，先来看配置文件：

```nginx
server {
 server_name return.ziyang.com dir.ziyang.com;
 server_name_in_redirect on;
 listen 8088;
 port_in_redirect on;
 absolute_redirect off;

 root html/;
}
```

`absolute_redirect` 默认是打开的，我们把它关闭了，看下是怎么返回的：

```yaml
➜  test_nginx curl localhost:8088/first -I
HTTP/1.1 301 Moved Permanently
Server: nginx/1.17.8
Date: Tue, 12 May 2020 00:31:36 GMT
Content-Type: text/html
Content-Length: 169
Connection: keep-alive
Location: /first/
```

这个时候看到返回的头部 `Location` 中没有加上域名。

下面再把 `absolute_redirect` 打开（默认是打开的，因此注释掉就行了），看下返回什么：

- `absolute_redirect on`
- `server_name_in_redirect on`
- `port_in_redirect on`

```shell
➜  test_nginx curl localhost:8088/first -I
HTTP/1.1 301 Moved Permanently
Server: nginx/1.17.8
Date: Tue, 12 May 2020 00:35:49 GMT
Content-Type: text/html
Content-Length: 169
Location: http://return.ziyang.com:8088/first/
Connection: keep-alive
```

可以看到，这时候就返回了域名，而且返回的是我们配置的主域名加端口号，这是因为，`server_name_in_redirect` 和 `port_in_redirect` 这两个指令打开了，如果关闭掉这两个指令，看下返回什么：

- `absolute_redirect on`

- `server_name_in_redirect off`

- `port_in_redirect off`

```shell
➜  test_nginx curl localhost:8088/first -I
HTTP/1.1 301 Moved Permanently
Server: nginx/1.17.8
Date: Tue, 12 May 2020 00:39:31 GMT
Content-Type: text/html
Content-Length: 169
Location: http://localhost/first/
Connection: keep-alive
```

这两个指令都设置为 `off` 之后，会发现返回的不再是主域名加端口号，而是我们请求的域名和端口号，如果在请求头中加上 `Host`，那么就会用 `Host` 请求头中的域名。

## index 模块

- 模块：`ngx_http_index_module`

- 功能：指定 `/` 结尾的目录访问时，返回 index 文件内容

- 语法：

    ```nginx
    Syntax: index file ...;
    Default: index index.html;
    Context: http, server, location
    ```

- 先于 autoindex 模块执行

这个模块，当我们访问以 `/` 结尾的目录时，会去找 root 或 alias 指令的文件夹下的 index.html，如果有这个文件，就会把文件内容返回，也可以指定其他文件。

## autoindex 模块

- 模块：`ngx_http_autoindex_module`，默认编译进 Nginx，使用 `--without-http_autoindex_module` 取消

- 功能：当 URL 以 `/` 结尾时，尝试以 html/xml/json/jsonp 等格式返回 root/alias 中指向目录的目录结构

- 语法：

    ```nginx
    # 开启或关闭
    Syntax: autoindex on | off;
    Default: autoindex off;
    Context: http, server, location
    # 当以 HTML 格式输出时，控制是否转换为 KB/MB/GB
    Syntax: autoindex_exact_size on | off;
    Default: autoindex_exact_size on;
    Context: http, server, location
    # 控制以哪种格式输出
    Syntax: autoindex_format html | xml | json | jsonp;
    Default: autoindex_format html;
    Context: http, server, location
    # 控制是否以本地时间格式显示还是 UTC 格式
    Syntax: autoindex_localtime on | off;
    Default: autoindex_localtime off;
    Context: http, server, location
    ```

### 实战

- 配置文件如下：

```nginx
server {
    server_name autoindex.ziyang.com;
    listen 8080;
    location / {
        alias html/;
        autoindex on;
        #index b.html;
        autoindex_exact_size on;
        autoindex_format html;
        autoindex_localtime on;
    }
}
```

这里我把 `index b.html` 这条指令给注释掉了，而 index 模块是默认编译进 Nginx 的，且默认指令是 `index index.html`，因此，会去找是否有 index.html 这个文件。

- 打开浏览器，访问 autoindex.ziyang.com:8080，html 目录下默认是有 index.html 文件的，因此显示结果为：

![](https://img2020.cnblogs.com/other/1833725/202005/1833725-20200521210221493-1958960782.jpg)

- 打开 `index b.html` 指令注释。由于 html 文件夹下并不存在 b.html 这个文件，所以请求会走到 autoindex 模块，显示目录：

![](https://img2020.cnblogs.com/other/1833725/202005/1833725-20200521210221834-631689202.jpg)

后面的文件大小显示格式就是由 `autoindex_exact_size on;` 这条指令决定的。

## concat模块

下面介绍一个可以提升小文件性能的模块，这个模块是由阿里巴巴开发的，在淘宝网中有广泛应用。

- 模块：ngx_http_concat_module

- 模块开发者：Tengine([https://github.com/alibaba/nginx-http-concat](https://github.com/alibaba/nginx-http-concat)) --add-module=../nginx-http-concat/

- 功能：合并多个小文件请求，可以明显提升 HTTP 请求的性能

- 指令：

    ```nginx
    #在 URI 后面加上 ??，通过 ”,“ 分割文件，如果还有参数，则在最后通过 ? 添加参数
    concat on | off
    default concat off
    Context http, server, location

    concat_types MIME types
    Default concat_types: text/css application/x-javascript
    Context http, server, location

    concat_unique on | off
    Default concat_unique on
    Context http, server, location

    concat_max_files numberp
    Default concat_max_files 10
    Context http, server, location

    concat_delimiter string
    Default NONE
    Context http, server, locatione
    concat_ignore_file_error on | off
    Default off
    Context http, server, location
    ```

打开淘宝主页，会发现小文件都是通过这个模块来提高性能的：

![](https://img2020.cnblogs.com/other/1833725/202005/1833725-20200521210222499-928396668.jpg)

这里就不做实战了，感兴趣的同学可以自己去编译一下这个模块，做一下实验，我把配置文件放在这里：

```nginx
server {
    server_name concat.ziyang.com;
    error_log logs/myerror.log debug;
    concat on;
    root html;
    location /concat {
        concat_max_files 20;
        concat_types text/plain;
        concat_unique on;
        concat_delimiter ':::';
        concat_ignore_file_error on;
    }
}
```

# log 阶段

下面终于来到了 11 个阶段的最后一个阶段，记录请求访问日志的 log 模块。

- 功能：将 HTTP 请求相关信息记录到日志
- 模块：`ngx_http_log_module`，无法禁用

## access 日志格式

```nginx
Syntax: log_format name [escape=default|json|none] string ...;
Default: log_format combined "...";
Context: http
```

默认的 combined 日志格式：

```nginx
log_format combined '$remote_addr - $remote_user [$time_local] '
'"$request" $status $body_bytes_sent ' '"$http_referer"
"$http_user_agent"';
```

## 配置日志文件路径

```nginx
Syntax: access_log path [format [buffer=size] [gzip[=level]] [flush=time] [if=condition]];
        access_log off;
Default: access_log logs/access.log combined;
Context: http, server, location, if in location, limit_except
```

- path 路径可以包含变量：不打开 cache 时每记录一条日志都需要打开、关闭日志文件

- if 通过变量值控制请求日志是否记录

- 日志缓存

  - 功能：批量将内存中的日志写入磁盘

  - 写入磁盘的条件：

        所有待写入磁盘的日志大小超出缓存大小；

        达到 flush 指定的过期时间；

        worker 进程执行 reopen 命令，或者正在关闭。

- 日志压缩

  - 功能：批量压缩内存中的日志，再写入磁盘
  - buffer 大小默认为 64KB
  - 压缩级别默认为 1（1最快压缩率最低，9最慢压缩率最高）
  - 打开日志压缩时，默认打开日志缓存功能

## 对日志文件名包含变量时的优化

```nginx
Syntax: open_log_file_cache max=N [inactive=time] [min_uses=N] [valid=time];
        open_log_file_cache off;
Default: open_log_file_cache off;
Context: http, server, location
```

- max：缓存内的最大文件句柄数，超出后用 LRU 算法淘汰
- inactive：文件访问完后在这段时间内不会被关闭。默认 10 秒
- min_uses：在 inactive 时间内使用次数超过 min_uses 才会继续存在内存中。默认 1
- valid：超出 valid 时间后，将对缓存的日志文件检查是否存在。默认 60 秒
- off：关闭缓存功能

日志模块没有实战。

---

到了这里，我们已经将 Nginx 处理 HTTP 请求的 11 个阶段全部梳理了一遍，每个阶段基本都有对应的模块。相信对于这样一个全流程的解析，大家都能够看懂 Nginx 的配置了，在此之上，还能够按照需求灵活配置出自己想要的配置，这样就真正的掌握了 11 个阶段。

最后，欢迎大家关注我的个人博客：[iziyang.github.io](https://iziyang.github.io/)
