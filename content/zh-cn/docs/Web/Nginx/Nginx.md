---
title: Nginx
linkTitle: Nginx
weight: 1
date: 2023-11-02T10:39
---

# 概述

> 参考：
>
> - [GitHub 项目，nginx/nginx](https://github.com/nginx/nginx)
>   - 原始代码：<https://hg.nginx.org/nginx/>
> - [org 官方网站](http://nginx.org/)
> - [官方网站](https://www.nginx.com/)
> - [官方网站,动态模块列表](https://www.nginx.com/products/nginx/modules/)

Nginx 称为 Engine X，可以做为 [Web](/docs/Web/Web.md) 服务器、代理服务器、缓存服务器、负载均衡器 等来使用。

传统上基于进程或线程模型架构的 web 服务通过每进程或每线程处理并发连接请求，这势必会在网络和 I/O 操作时产生阻塞，其另一个必然结果则是对内存或 CPU 的利用率低下。生成一个新的进程/线程需要事先备好其运行时环境，这包括为其分配堆内存和栈内存，以及为其创建新的执行上下文等。这些操作都需要占用 CPU，而且过多的进程/线程还会带来线程抖动或频繁的上下文切换，系统性能也会由此进一步下降。

在设计的最初阶段，nginx 的主要着眼点就是其高性能以及对物理计算资源的高密度利用，因此其采用了不同的架构模型。受启发于多种操作系统设计中基于“事件”的高级处理机制，nginx 采用了模块化、事件驱动、异步、单线程及非阻塞的架构，并大量采用了多路复用及事件通知机制。在 nginx 中，连接请求由为数不多的几个仅包含一个线程的进程 worker 以高效的回环(run-loop)机制进行处理，而每个 worker 可以并行处理数千个的并发连接及请求。

Nginx 会按需同时运行多个进程：一个主进程(master)和几个工作进程(worker)，配置了缓存时还会有缓存加载器进程(cache loader)和缓存管理器进程(cache manager)等。所有进程均是仅含有一个线程，并主要通过“共享内存”的机制实现进程间通信。主进程以 root 用户身份运行，而 worker、cache loader 和 cache manager 均应以非特权用户身份运行。

Nginx 特性：

1. 模块化设计，较好的扩展性，所有配置均有指定的模块进行处理。
2. 高可靠 master --> worker，主控进程不接收和响应用户请求，主控进程负责解析配置文件并生成多个工作进程，工作进程来响应用户请求
   1. 主控进程读取并验证配置，创建或绑定套接字，启动及终止和维护 worker 进程的个数，无须重启进程让新配置的配置文件进行加载，以及完成平滑版本升级等等
   2. 工作进程，负责缓存加载的(反向代理时候用)，负责响应用户请求，cache manager 缓存管理
3. 低内存消耗，10000 个 keep-alive 模式下的 connection，仅需 2.5MB 内存
4. 支持热部署，不停机而更新配置文件，日志文件滚动，升级程序版本
5. 支持事件驱动、AIO、mmap

基本功能：

- 静态资源的 web 服务器，能缓存打开的文件描述符
- http、SMTP、pop3 协议的反向代理服务器
- 缓存加速、负载均衡
- 支持 FastCGI(fpm，LNMP)，uWSGI(Python)等
- 模块化(非 DSO 机制)、过滤器 zip、SSI 及图像的大小调整
- 支持 SSL(https)

扩展功能

- 基于名称和 IP 的虚拟主机
- 支持 keepalive
- 支持平滑升级
- 定制访问日志
- 支持 url 重写
- 支持路径别名
- 支持基于 IP 及用户的访问控制
- 支持速率限制，支持并发数限制

## Nginx 架构

详见 [Nginx 源码解析](/docs/Web/Nginx/Nginx%20源码解析.md)

## 总结

看似很复杂，其实总结起来 Nginx 主要就是两个功能，这两个功能也是配置文件中的主要内容，各种指令都离不开这两方面。

1. 定义接收流量的人口(port 或者 域名等)
2. 定义处理流量的规则(转发或者丢弃等)

所以，**Nginx 的本质，就是流量处理**

# Nginx 部署

## docker 运行 Nginx

(可选)需要在宿主机的 /etc/nginx 下准备 nginx 的配置文件。这些基本配置文件可以先启动一个 nginx 容器，将容器内的配置文件全部拷贝到宿主机上即可。

```bash
docker run -d --name nginx --network host --restart=always \
  -v /etc/nginx:/etc/nginx:ro \
  nginx:stable-alpine
```

也可以使用自己的配置来运行 nginx

```bash
mkdir -p /opt/nginx/conf/stream.d
# 生成配置文件
cat > /opt/nginx/conf/nginx.conf <<EOF
user  nginx;
worker_processes  auto;
error_log  /var/log/nginx/error.log warn;
pid        /var/run/nginx.pid;
events {
    worker_connections  1024;
}
stream {
        include /etc/nginx/stream.d/*.conf;
}
http {
}
EOF
# 运行
docker run -d --name nginx --network host --restart=always \
  -v /opt/nginx/conf/stream.d:/etc/nginx/stream.d:ro \
  -v /opt/nginx/conf/nginx.conf:/etc/nginx/nginx.conf:ro \
  nginx:stable-alpine
```

挂载一些静态文件到容器中以响应客户端的 web 请求

```bash
export StaticFiles="/root/projects/DesistDaydream/javascript-learning/9_browser"

docker run -d --name nginx --network host --restart=always \
  -v ${StaticFiles}:/usr/share/nginx/html:ro \
  nginx:stable-alpine
```

# Nginx 关联文件

**/etc/nginx/** # nginx 运行所需配置所在目录

- **./nginx.conf** # nginx 主程序运行所读取的默认配置文件。

配置文件官方介绍：<https://docs.nginx.com/nginx/admin-guide/basic-functionality/managing-configuration-files/>

修改完配置后，可以使用 nginx -s reload 命令使之生效

下面是 nginx 默认的基本配置示例，可以实现一个简单的 web 服务。

```nginx
user nginx; # 指定以nginx用户来运行nginx进程
worker_processes auto;
error_log /var/log/nginx/error.log;
pid /run/nginx.pid;

include /usr/share/nginx/modules/*.conf; #包含/usr/share/nginx/modules/目录下以.conf为结尾的所有文件，加载其中的配置

events {
    worker_connections 1024;
}

http {
    log_format  main  '$remote_addr - $remote_user [$time_local] "$request" '
                      '$status $body_bytes_sent "$http_referer" '
                      '"$http_user_agent" "$http_x_forwarded_for"';

    access_log  /var/log/nginx/access.log  main;

    sendfile            on;
    tcp_nopush          on;
    tcp_nodelay         on;
    keepalive_timeout   65;
    types_hash_max_size 2048;

    include             /etc/nginx/mime.types;
    default_type        application/octet-stream;

    include /etc/nginx/conf.d/*.conf;

    server {
        listen       80 default_server;
        listen       [::]:80 default_server;
        server_name  _;
        root         /usr/share/nginx/html; # 指定nginx的工作的/目录。i.e.location中/目录的起始位置

        include /etc/nginx/default.d/*.conf;

        location / {
        }

        error_page 404 /404.html;
            location = /40x.html {
        }

        error_page 500 502 503 504 /50x.html;
            location = /50x.html {
        }
    }
}
```

# 命令行工具

## 应用示例

让nginx在前台运行，常用于container中

- `nginx -g 'daemon off;'`

# 分类

 #代理 #网络 #Web