---
title: DNSmasq
---

# 概述

> 参考：
> - [Manual(手册)](https://thekelleys.org.uk/dnsmasq/docs/dnsmasq-man.html)

DNSmasq 是一个轻量的 DHCP 和 DNS 缓存 服务。

# DNS 泛解析实例

最近遇到一个问题，需要在服务器上对域名进行泛解析，比如访问百度的域名统统解析到 6.6.6.6，然而发现 hosts 文件根本就不支持类似 \*.baidu.com 的这种写法。

于是乎就在网上找了下资料，发现可以通过 Dnsmasq 来解决这个问题，原理其实就是本机的 DNS 指向 Dnsmasq 服务器，然后 Dnsmasq 通过类似通配符 (_) 的方式进行匹配，凡是匹配到 _.baidu.com 的都解析到 6.6.6.6。 **利用 Dnsmasq 实现 hosts 泛解析**

**环境介绍**

    $ uname \-a
    Linux ansheng 3.10.0\-957.1.3.el7.x86\_64 #1 SMP Thu Nov 29 14:49:43 UTC 2018 x86\_64 x86\_64 x86\_64 GNU/Linux
    $ whoami
    root
    $ cat /etc/redhat\-release
    CentOS Linux release 7.6.1810 (Core)

## **安装 Dnsmasq**

安装非常简单，通过 yum 即可。

    $ yum install dnsmasq \-y

**配置 Dnsmasq**
先把配置文件备份一份

    $ cp /etc/dnsmasq.conf /etc/dnsmasq.conf\_bak

Dnsmasq 的配置在配置文件中都有详细的说明，你可以通过阅读配置文件的注释更改自己想要的配置，我只是想做泛解析，所以我的配置如下：

    $ vim /etc/dnsmasq.conf
    # 严格按照 resolv\-file 文件中的顺序从上到下进行 DNS 解析, 直到第一个成功解析成功为止
    strict\-order

    # 监听的 IP 地址
    listen\-address\=127.0.0.1

    # 设置缓存大小
    cache\-size\=10240

    # 泛域名解析，访问任何 baidu.com 域名都会被解析到 6.6.6.6
    address\=/baidu.com/6.6.6.6

[域名解析](https://cloud.tencent.com/product/cns?from=10680)默认读取 /etc/hosts 文件到本地域名配置文件（不支持泛域名）。

DNS 配置默认读取 /etc/resolv.conf 上游 DNS 配置文件，如果读取不到 /etc/hosts 的地址解析，就会转发给 resolv.conf 进行解析地址。

- DNS 配置文件
  $ vim /etc/resolv.conf 这些都是常用的 DNS，可以配置很多

nameserver 127.0.0.1  # 一定要放在第一个
nameserver 8.8.8.8
nameserver 8.8.4.4
nameserver 1.1.1.1

- 启动服务
  $ systemctl enable --now dnsmasq
  Created symlink from /etc/systemd/system/multi-user.target.wants/dnsmasq.service to /usr/lib/systemd/system/dnsmasq.service.
  查看运行状态
  $ systemctl status dnsmasq
  ● dnsmasq.service - DNS caching server.
  Loaded: loaded (/usr/lib/systemd/system/dnsmasq.service; enabled; vendor preset: disabled)
  Active: active (running) since 日 2018-12-23 09:00:12 UTC; 3s ago
  Main PID: 3844 (dnsmasq)
  CGroup: /system.slice/dnsmasq.service
  └─3844 /usr/sbin/dnsmasq -k
  12 月 23 09:00:12 ansheng systemd\[1]: Started DNS caching server..
  12 月 23 09:00:12 ansheng dnsmasq\[3844]: started, version 2.76 cachesize 10000
  12 月 23 09:00:12 ansheng dnsmasq\[3844]: compile time options: IPv6 GNU-getopt DBus no-i18n IDN DHCP DHCPv6 no-Lua TFTP no-conntrack ipset auth no-DNSSEC loop-detect inotify
  12 月 23 09:00:12 ansheng dnsmasq\[3844]: reading /etc/resolv.conf
  12 月 23 09:00:12 ansheng dnsmasq\[3844]: ignoring nameserver 127.0.0.1 - local interface
  12 月 23 09:00:12 ansheng dnsmasq\[3844]: using nameserver 8.8.8.8#53
  12 月 23 09:00:12 ansheng dnsmasq\[3844]: using nameserver 8.8.4.4#53
  12 月 23 09:00:12 ansheng dnsmasq\[3844]: using nameserver 1.1.1.1#53
  12 月 23 09:00:12 ansheng dnsmasq\[3844]: read /etc/hosts - 6 addresses

## **测试**

    $ ping baidu.com
    PING baidu.com (6.6.6.6) 56(84) bytes of data.
    ^C
    \--\- baidu.com ping statistics \--\-
    2 packets transmitted, 0 received, 100% packet loss, time 1000ms

    $ ping www.baidu.com
    PING www.baidu.com (6.6.6.6) 56(84) bytes of data.
    ^C
    \--\- www.baidu.com ping statistics \--\-
    2 packets transmitted, 0 received, 100% packet loss, time 999ms

    $ ping pan.baidu.com
    PING pan.baidu.com (6.6.6.6) 56(84) bytes of data.
    ^C
    \--\- pan.baidu.com ping statistics \--\-
    2 packets transmitted, 0 received, 100% packet loss, time 999ms

由上可以看到，几乎访问任何 baidu.com 的域名都会被解析到 6.6.6.6，基本上就达到了我们最初的目的。

## **利用 Dnsmasq 缓存特性实现 DNS 加速**

Dnsmasq 还有一项非常有用的功能就是可以对已经解析过的域名进行缓存，下次在访问这个域名的时候就可以直接返回 IP 地址，而不再需要经过 DNS 查询，这对于扶墙的来说，其实也算是一点优化，默认已经配置好了，我们只需要来演示下缓存的效果。

- 安装 dig 工具

```bash
$ yum install bind-utils -y
```

演示效果

```bash
$ dig www.centos.com | grep "Query time"
;; Query time: 88 msec
$ dig www.centos.com | grep "Query time"
;; Query time: 0 msec
$ dig www.centos.com | grep "Query time"
;; Query time: 0 msec
$ dig www.centos.com | grep "Query time"
;; Query time: 0 msec
$ dig www.youtube.com | grep "Query time"
;; Query time: 28 msec
$ dig www.youtube.com | grep "Query time"
;; Query time: 0 msec
$ dig www.qq.com | grep "Query time"
;; Query time: 71 msec
$ dig www.qq.com | grep "Query time"
;; Query time: 0 msec
```

看看上面的对比，查询时间缩小了不少倍，可见缓存已经产生作用。

> 来源：_安生博客_ 原文：[_http://t.cn/AiCohacf_](http://t.cn/AiCohacf) 题图：_来自谷歌图片搜索_ 版权：_本文版权归原作者所有_ 投稿：\_欢迎投稿，投稿邮箱: \_[_editor@hi-linux.com_](mailto:editor@hi-linux.com)

本文分享自微信公众号 - 运维之美（Hi-Linux），作者：ansheng

原文出处及转载信息见文内详细说明，如有侵权，请联系 <yunjia_community@tencent.com> 删除。

原始发表时间：2019-11-02

本文参与[腾讯云自媒体分享计划](https://cloud.tencent.com/developer/support-plan)，欢迎正在阅读的你也加入，一起分享。
<https://cloud.tencent.com/developer/article/1534150>
