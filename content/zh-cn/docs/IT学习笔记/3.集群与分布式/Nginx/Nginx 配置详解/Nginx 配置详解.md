---
title: Nginx 配置详解
---

# 概述

> 参考：
> - [org 官方文档,初学者指南-配置文件结构](http://nginx.org/en/docs/beginners_guide.html#conf_structure)
> - [org 官方文档,全部指令列表](http://nginx.org/en/docs/dirindex.html)
> - [org 官方文档,全部变量列表](http://nginx.org/en/docs/varindex.html)
> - [官方文档,管理指南-基础功能-创建 NGINX 配置文件](https://docs.nginx.com/nginx/admin-guide/basic-functionality/managing-configuration-files/#)

Nginx 由 **Modules(模块)** 组成， Modules 由配置文件中的 **Directives(指令) **控制其运行行为。有的 Directives 可以控制多个模块，只不过在控制不同模块时，产生的效果也许会不尽相同。

## Directives(指令)

Directives(指令) 分为如下几种：

- Simple Directives(简单指令)
- Block Directives(块指令)
- Conexts(配置环境 | 上下文)

### Simple Directives(简单指令)

由空格分割的 **Name(指令名称)** 和 **Parameters(指令参数)**，以 `;` 符号结尾。

- 如果从 Nginx 的代码角度看，指令就相当于结构体中的属性，参数就是该属性的值。

### Block Directives(块指令)

将多个相关的简单指令组合在一起的容器，并将它们用 `{}` 符号包围起来。

### Top Level Directives(顶级指令) # 也称为 Contexts(配置环境 | 上下文)。

将多个相关的 块指令 和 简单指令 组合在一起的指令。一共分为 4 类 Contexts：

- [**events {}**](/docs/IT学习笔记/3.集群与分布式/Nginx/Nginx%20 配置详解/events%20 模块指令.md 配置详解/events 模块指令.md)\*\* \*\*# 用于配置如何处理常规连接。
- [**http {}**](/docs/IT学习笔记/3.集群与分布式/Nginx/Nginx%20 配置详解/http%20 模块指令.md 配置详解/http 模块指令.md)\*\* \*\*# http 流量处理配置，通常用来配置 7 层代理。由 ngx_http_core_module 模块处理其中配置
- **mail {} **# mail 流量处理配置。由 ngx_mail_core_module 模块处理其中配置
- [**stream {}**](/docs/IT学习笔记/3.集群与分布式/Nginx/Nginx%20 配置详解/core%20 模块指令.md 配置详解/core 模块指令.md)\*\* \*\*# TCP 和 UDP 流量处理配置，通常用来配置 4 层代理。由 ngx_stream_core_module 模块处理其中配置

[**main**](/docs/IT学习笔记/3.集群与分布式/Nginx/Nginx%20 配置详解/core%20 模块指令.md 配置详解/core 模块指令.md) # 如果某些指令在上述 4 类 Contexts 之外，则称之为 main Context。可以说，events{}、http{}、mail{}、stream{} 四个 Contexts，都属于 main 上下文中的指令。说白了，main 上下文就是 Nginx 的配置文件~~~其实，main 就是指最顶层的 core 模块指令

每一个 Context 类型的指令都对应控制一个 NGX_CORE_MODULE 类型的模块

- main 指令 —> 控制 core 模块
- events {} 指令 —> 控制 events 模块
- http {} 指令 —> 控制 http 模块
- mail {} 指令 —> 控制 mail 模块
- stream {} 指令 —> 控制 stream 模块

所以，配置文件的格式实际上也是一个树状结构：
![树形结构.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/tt8mpd/1619841196176-fa5e09e5-44b9-49e2-bdc3-bccc385d0218.png)
最顶层的 main 指令，包含 简单指令 和 4 个 Contexts，每个 Context 又包含 简单指令/块指令。

## Inheritance(继承)

通常，子块指令 将会继承其 父块指令 中的指令设置。某些 简单指令 可以出现在多个 块指令 中，这种情况下，可以通过在 子块指令 中设置该 简单指令，以便覆盖其从 父块指令 中继承过来的设置。

## Variables(变量)

在 Nginx 的配置文件中，还可以设置并引用变量，通过 `set` 指令，可以定义一个变量，并在其他指令中使用 `$变量名` 引用变量。同样，也有环境变量的概念，Nginx 的很多模块在加载并时候后，会产生环境变量，也可以直接引用

变量最常用的地方就是通过 `log_format` 指令定义日志内容~~

## 基本配置示例

```nginx
# 在 4 种配置环境之外的指令属于 main 配置环境，一般用于配置 nginx 运行的基础信息
user nobody; # 这是一个 main context 中的指令

# 事件配置环境。用于配置连接参数等信息
events {
    # 关于连接处理的配置
}

# http 流量处理配置环境
http {
    # 影响所有 virtual servers 的 http 流量的配置
    server {
        server_name localhost;
        # 配置处理 http 流量的 Virtual Server 1
        location /one {
            # 配置用于处理以'/one'开头的URI的流量
        }
        location /two {
            # 配置用于处理以'/two'开头的URI的流量
        }
    }

    server {
        # 配置处理 http 流量的 Virtual Server 2
    }
}

# mail 流量处理环境
mail {
    ....
}

# TCP 和 UDP 流量处理配置环境
stream {
    # 影响所有 virtual servers 的  TCP/UDP 流量的配置
    server {
        #配置处理 TCP 流量的 Virtual Server 1
    }
}
```

Note：

1. 通常 `{}` 中的指令只对大括号内部内容生效；不在 `{}` 中且在文件开头的，则对全局生效；配置指令要以分号结尾
2. 配置技巧：为了使配置更易于维护，还可以将大段的配置拆分为一组一组存储在 /etc/nginx/conf.d/ 目录下的文件，并在 nginx.conf 这个主配置文件中使用 `include` 指令来引用这些文件。

## 总结

其实，Nginx 的配置本质就是对模块的配置，这些指令就是模块可以接收的形参；而指令的值就是代码中，应该给模块传递的实参。

而 Nginx 配置文件的组织结构更像 INI 格式的配置文件，Nginx 的 顶级指令就是 Selections(部分) 的概念；简单指令就是 Key/Value Pairs(键/值对) 的概念；块指令其实就是一种嵌套形式的简单指令。

# 配置文件详解

本质上，Nginx 的配置文件，就是由 **Directives(指令)** 与 **Contexts(配置环境 | 上下文)** 组成，而这些指令和配置环境，又是围绕 Virtual Servers 运转。

## Virtual Servers(虚拟主机) 抽象概念

**Virtual Servers(虚拟主机) **是 Nginx 的抽象概念，Virtual Servers 用来定义 **流量入口 **和 **流量处理**。而这两块功能也是 Nginx 正常工作的最基本功能。

### 流量入口

流量入口包括监听的地址、域名等。对于发送到 Nginx 的流量，根据流量中的内容，分配到不同的入口，进行后续流量处理。

### 流量处理

流量处理包括将流量代理到何处、是否丢弃、连接超时时长等。

#### 后端服务器

由于流量处理需要将流量代理到指定的服务器，所以还需要配置后端服务器来接收流量。

## 总结

Virtual Server 的概念，通过 **`server{}` **指令来实现。在每个流量处理的配置环境中，都应该包含一个或多个 server{} 指令。server{} 指令是 nginx 正常运行的基础配置。虚拟主机，顾名思义，对于用户来说，访问的就是一台一台服务器，但是对于 nginx 来说，是虚拟出来的。

1. **对于 http 流量(http 配置环境)** # 每个 `server{} 指令块` 控制访问特定域名或者 ip 地址上对资源请求的处理。server 指令块中的一个或多个 location 指令块定义了根据 URI 来处理流量的规则
   1. 比如用户访问 map.baidu.com 和 baike.baidu.com。看上去是访问了两台服务器，但是实际上，这是经过作为代理设备的 ngxin 来进行选择后的虚拟服务器。一般情况下，baike.baidu.com 与 map.baidu.com 这俩域名所解析出来的 ip 应该是同一个公网 ip(比如 123.123.123.123)(baidu 有钱用很多公网 IP 除外)。所以可以想到，用户在浏览器输入任何一个域名，访问请求都会来到 123.123.123.123，然后根据请求报文中的 Request-URL 字段中的域名与 server_name 进行配对，用户输入的 URL 中域名与哪个 server_name 相同，则该请求就会通过这个 server 来进行处理，然后根据该 server 中 location 的关键字来决定把改请求转发给哪里。
2. **mail 和 TCP/UDP 流量(mail 和 stream 配置环境)** # 每个 `server{} 指令块` 控制处理到达指定 TCP port 或 UNIX socket 的流量。
   1. 比如用户访问 30000 端口，则可以根据其中的规则，将 对 30000 端口发起的请求，代理到其他设备的某些端口上。

其实说白了，每个 Virtual Server 都相当于一个独立运行的服务，用来处理客户端的请求。具体如何处理，则在每个 server{} 指令块中定义。可以这么说，Nginx 中所有指令，其实都是为 Virtual Servers 服务的。

# 指令详解

## [main 模块指令](/docs/IT学习笔记/3.集群与分布式/Nginx/Nginx%20 配置详解/core%20 模块指令.md 配置详解/core 模块指令.md)\*\*

## [events 模块指令](/docs/IT学习笔记/3.集群与分布式/Nginx/Nginx%20 配置详解/events%20 模块指令.md 配置详解/events 模块指令.md)

[http 模块指令](/docs/IT学习笔记/3.集群与分布式/Nginx/Nginx%20 配置详解/http%20 模块指令.md 配置详解/http 模块指令.md)

常用来配置七层代理、web 应用

## mail 模块指令

[stream 模块指令](/docs/IT学习笔记/3.集群与分布式/Nginx/Nginx%20 配置详解/stream%20 模块指令.md 配置详解/stream 模块指令.md)

常用来配置四层代理
