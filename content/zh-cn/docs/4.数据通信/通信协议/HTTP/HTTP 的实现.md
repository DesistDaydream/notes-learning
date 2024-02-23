---
title: HTTP 的实现
---

HTTP 服务器的程序(想提供 web 服务必须要安装一下程序中的一个)

- httpd(apache)
- nginx
- lighttpd
- 应用程序服务器：上面的程序如果不附加插件则只支持静态的网页，装上下面的程序还能解析 PHP 等动态界面
  - IIS
  - tomcat，jetty，就 boss，resin
  - webshpere,weblogic,oc4j

# httpd

apache(a patchy server)的特性

- 高度模块化：core+modules
- DSO：Dynamic Shared Object
- MPM：Multipath Processing Modules 多路处理模块，不同的工作方式，可以切换，使用不同模块可以满足不同需求
  - prefork：多进程模型，每个进程响应一个请求
    - 一个主进程：负责生成 n 个子近侧很难过，子进程也成为工作进程，每个子进程处理一个用户请求，即便没有用户请求，也会预先生成多个空闲进程，随时等待请求到达，最大不超过 1024 个
  - worker：多线程模型
    - 一个主进程：负责生成子进程；负责创建套接字；负责接收请求，并将其派发给某子进程进行处理；
    - 多个子进程：每个子进程负责生成多个线程；
    - 每个线程：负责响应用户请求；
    - 并发响应数量：m\*n
    - m：子进程数量
    - n：每个子进程所能创建的最大线程数量；
  - event：事件驱动模型，多进程模型，每个进程响应多个请求（老版本系统不支持，systemd 系统支持）
    - 一个主进程 ：负责生成子进程；负责创建套接字；负责接收请求，并将其派发给某子进程进行处理；
    - 子进程：基于事件驱动机制直接响应多个请求；

## Httpd 配置

程序环境(.init 系统下)

- 配置文件
  - /etc/httpd/conf/httpd.conf
  - /etc/httpd/conf.d/\*.conf
- 服务脚本
  - /etc/rc.d/init.d/httpd
  - 配置文件/etc/sysconfig/httpd
- 主程序文件
  - /usr/sbin/httpd
  - /usr/sbin/httpd.event
  - /usr/sbin/httpd.worker
- 日志文件目录
  - /var/log/httpd
    - access_log:访问日志
    - error_log:错误日志
  - 站点文档目录
    - /var/www/html
