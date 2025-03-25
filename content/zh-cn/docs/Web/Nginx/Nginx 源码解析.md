---
title: Nginx 源码解析
linkTitle: Nginx 源码解析
weight: 20
---

# 概述

> 参考：
>
> - [GitHub 项目，nginx/nginx](https://github.com/nginx/nginx)

Nginx 的架构设计是高度模块化的，从 Nginx 的源码目录与 Nginx 模块化及其功能的划分是紧密结合的。

```bash
# tree -d
.
├── auto
├── conf
├── contrib
├── docs
├── misc
└── src
    ├── core
    ├── event
    │   └── modules
    ├── http
    │   ├── modules
    │   │   └── perl
    │   └── v2
    ├── mail
    ├── misc
    ├── os
    │   ├── unix
    │   └── win32
    └── stream
```

每个模块必须要具备 `ngx_module_t` 这个数据结构，比如 src/core/nginx.c 文件中，在 ngx_module_t 数据结构中，定义了 ngx_core_module 模块

```c
ngx_module_t  ngx_core_module = { // 定义模块名称为 ngx_core_module
    NGX_MODULE_V1,
    &ngx_core_module_ctx,                  /* module context */
    // ngx_core_commands 是该模块的指令，这里一般会定义一个数组，数组中每个元素就是一个指令
    ngx_core_commands,                     /* module directives */
    // NGX_CORE_MODULE 是该模块的类型
    NGX_CORE_MODULE,                       /* module type */
    NULL,                                  /* init master */
    NULL,                                  /* init module */
    NULL,                                  /* init process */
    NULL,                                  /* init thread */
    NULL,                                  /* exit thread */
    NULL,                                  /* exit process */
    NULL,                                  /* exit master */
    NGX_MODULE_V1_PADDING
};
```

再看 src/http/ngx_http.c 文件

```c
ngx_module_t  ngx_http_module = {
    NGX_MODULE_V1,
    &ngx_http_module_ctx,                  /* module context */
    ngx_http_commands,                     /* module directives */
    // ngx_http_module 模块属于 NGX_CORE_MODULE 类型
    NGX_CORE_MODULE,                       /* module type */
    NULL,                                  /* init master */
    NULL,                                  /* init module */
    NULL,                                  /* init process */
    NULL,                                  /* init thread */
    NULL,                                  /* exit thread */
    NULL,                                  /* exit process */
    NULL,                                  /* exit master */
    NGX_MODULE_V1_PADDING
};
```

http 模块属于 NGX_CORE_MODULE 类型如果看 src/stream/ngx_stream.c、src/event/ngx_event.src、src/mail/ngx_mail_module.c 这三个文件，里面都有类似的定义

最后看 src/http/modules/ 目录，这里就是都有属于 NGX_HTTP_MODULE 模块的子模块。定义方式类似。

# 总结

所以，从代码结构可以看出来，Nginx 模块定义了如下几大类型，并且每一类模块包含不同的模块

- **NGX_CORE_MODULE**
- **NGX_EVENT_MODULE**
- **NGX_HTTP_MODULE**
- **NGX_MAIL_MODULE**
- **NGX_STREAM_MODULE**

与此同时，每种类型的模块下，包含很多子模块，模块类型是抽象的概念，主要是为了帮助模块进行分类。而模块，则是真正处理任务的：

- **NGX_CORE_MODULE**
  - ngx_core_module # core 模块，这个比较特殊，core 模块类型下，其实就是包含了一个 core 模块。
  - ngx_event_module # event 模块
  - ngx_http_module # http 模块
  - ngx_mail_module # mail 模块
  - ngx_stream_module # stream 模块
- **NGX_EVENT_MODULE**
  - ngx_event_core_module # event core 模块
  - ...... 太多了就不一一列举了，详见官网
- **NGX_HTTP_MODULE**
  - ngx_http_core_module # http core 模块
  - ...... 太多了就不一一列举了，详见官网
- **NGX_MAIL_MODULE**
  - ngx_mail_module # mail core 模块
  - ...... 太多了就不一一列举了，详见官网
- **NGX_STREAM_MODULE**
  - ngx_stream_module # stream core 模块
  - ...... 太多了就不一一列举了，详见官网

从上面的分类也可以看出来，每个类型的模块都有一个被称为 core 的模块来实现这个类型的功能。这其实就可以组成一个树状的结构：

![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/eygybr/1619798135171-e81a3491-2566-4f50-ab0f-4b0868a5670e.png)

这个图不太好，有机会自己画一个。

画成图就是这个样子的：

![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/eygybr/1619768579247-4047676f-2ece-4b9c-bed2-7d53d7f32ae4.png)

就像 [官方文档的模块参考](http://nginx.org/en/docs/) 页面一样，我们可以看到下面这样的命名规则 `ngx_XXX_YYY_module`。其中 XXX 就是顶级模块名，YYY 就是属于 XXX 模块的次级模块名。效果如下：

```bash
# 包含 core 和 event 模块
Core functionality
# http 模块
ngx_http_core_module
ngx_http_access_module
ngx_http_addition_module
......
# mail 模块
ngx_mail_core_module
ngx_mail_auth_http_module
......
# stream 模块
ngx_stream_core_module
ngx_stream_access_module
ngx_stream_geo_module
......
```

Nginx 通过这种架构，将各个模块组织成一条链，当有请求到达时，请求依次经过这条链上的部分或者全部模块，进行处理。每个模块实现特定的功能。

而这些模块的运行行为，就是通过配置文件中的 [**Directive(指令)**](/docs/Web/Nginx/Nginx%20配置详解/Nginx%20配置详解.md) 实现的。一般每个顶级模块都有一个同名的 顶级指令，比如 http 模块对应 http{} 顶级指令。

Nginx 模块的架构，也使得 Nginx 实现了灵活得可扩展性。

## 指令

在模块代码中，指令的定义在其所属模块代码内的 `static ngx_command_t` 类型的数组中。

比如  ngx_stream_log_module 模块，可以在 [nginx/src/stream/ngx_stream_log_modul.c](https://github.com/nginx/nginx/blob/stable-1.26/src/stream/ngx_stream_log_module.c#L145) 处看到定义了一个 名为 ngx_stream_log_commands 的数组，名称前半部分是模块名称，后面跟一个 `_commands`。数组内部可以看到指令的名称。

```c
static ngx_command_t  ngx_stream_log_commands[] = {
    { ngx_string("log_format"),
      ......},
    { ngx_string("access_log"),
      ......},
    ......
}
```
