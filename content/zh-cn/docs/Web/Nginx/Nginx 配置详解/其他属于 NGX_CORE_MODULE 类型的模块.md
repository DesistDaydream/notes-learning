---
title: 其他属于 NGX_CORE_MODULE 类型的模块
weight: 6
---

# NGX_CONF_MODULE

## inclue

- 语法：**include /PATH/FILE**

http://nginx.org/en/docs/ngx_core_module.html#include

在该配置中包含一个网站-可用的配置文件，即把定义的文件内容引入到这里，(也可以写入多个 include，引入多个配置文件以便管理，包括但不限于 server 配置，还可以是 nginx 的配置文件，mail 的配置文件等)

- 作用范围：可以作用在任意 Contexts 中

该指令可以写在任意 块指令 中，只要被包含的文件格式，符合当前 块指令 应该包含的语法即可。
