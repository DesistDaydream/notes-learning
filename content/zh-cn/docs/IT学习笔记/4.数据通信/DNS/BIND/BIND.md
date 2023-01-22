---
title: BIND
---

# 概述

> 参考：
> - [ISC-BIND9 官方网站](https://www.isc.org/bind/)
> - [Wiki,BIND](https://en.wikipedia.org/wiki/BIND)

**Bekerley Internat Name Domain(伯克利互联网名字域，简称 BIND)** 是实现 DNS 服务的应用程序。该程序最著名的组件被称为 **named**，主要用来实现两个最主要的 DNS 功能：**NameServer(名称服务器)** 与 **Resolver(解析器)**。

该软件最初是在 1980 年代初在加州大学伯克利分校 (UCB) 设计的。该名称起源于 Berkeley Internet Name Domain 的首字母缩写词，反映了该应用程序在 UCB 中的使用。最新版本是 BIND 9，于 2000 年首次发布，仍然由 Internet Systems Consortium (ISC) 积极维护，每年发布数次新版本。

BIND9 已进化为一个非常灵活，全功能的 DNS 系统。无论您的应用程序是什么，绑定 9 可能具有所需的功能。作为第一个，最旧，最常见的解决方案，还有更多网络工程师已经熟悉绑定 9，而不是与任何其他系统。

# BIND 部署

dns 服务，包名 bind，程序名 named

基础程序包：bind 提供服务，bing-libs 提供库文件，bind-utils 提供测试程序

rndc:remote name domain controller,默认与 bind 安装在同一主机，且只能通过 127.0.1 来连接 named 进程，提供辅助性的管理功能，使用 tcp 的 953 端口

配置文件：

- /etc/named.conf
- /etc/named.rfc1912.zomes #该文件的引用，定义在 named.conf 的最后几行
- /etc/rndc.key

解析库文件：/var/named/ZONE_NAME.ZONE,有以下注意事项

- 一台物理服务器可以同时为多个区域提供解析
- 必须要有根区域文件 named.ca
- 应该有两个实现 localhost 和本地回环地址的解析库

/etc/named.conf # 配置文件 keywords 说明

1. options{ #用于全局 BIND 配置，BIND 的工作目录在 /var/named

- listen-on port NUM { IP1; IP2; }; #设置 DNS 服务监听的端口号 NUM 和监听该端口的 IP 地址
- allow-query { any; }; #设置任何人都可以来这台服务器解析
- forward { first|only }； #转发服务器配置，加了此项则定义先进行域名解析请求转发,转发不了再去迭代查询,可用 first 或者 only 模式
- forwarders { IP； } #转发服务器配置，转发的 IP 地址
- }；

1. loggin{} #配置哪些需要记录，哪些需要忽略
2. zone "ZONE_NAME" IN { #定义 DNS 区域。本机能够为哪些 zone(区域)进行解析，就要定义哪些 zone，比如域名 google.com，它包含子域名 mail.google.com 和 analytics.google.com 等。这几个域名都有一个由 zone 语句定义的区域，该定义可以直接写在.zone 的解析库文件里

- type ; #定义该服务器是什么职责，包括主，辅助，根，转发这四种
- file ”ZONE_NAME.zone“； #区域解析库文件，默认在该文件在/var/named 目录下，所以该位置直接写文件名即可
- };

1. include #在 named.conf 中包含另一个文件。比如 named.rfc1912.zomes 该文件包含其余定义的 zone 信息

/var/named/ZONE_NAME.zone #解析库文件说明

EXAMPLE

- $TTL 86400 #该条目告诉 BIND 每个单独记录的 TTL 值（time to live，生存时间值）。它是以秒为单位的数值，比如 14,400 秒（4 个小时），因此 DNS 服务器最多缓存你的域文件 4 个小时，之后就会向你的 DNS 服务器重新查询。
- $ORIGIN baidu.com. #定义该项后，资源格式里可以省略后面的根域名，顶级域名，所有写的域名自带该变量定义的域名
- @ IN SOA ns1.baidu.com. admin.baidu.com (
- 2015042201
- 1H
- 5M
- 7D
- 1D )
-       	IN		NS		ns1.baidu.com
-       	IN		NS		ns2.baidu.com
- ns1 IN A 1.1.1.1
- ns2 IN A 1.1.1.2
- www IN A 1.1.1.3
-       IN		A		1.1.1.4
-       IN		MX	10	mx1.baidu.com
-       IN		MX	20	mx2.baidu.com
- 4.3.2.1in-addr.arpa. IN PTR www.baidu.com
- web IN CNAME www.baidu.com

主从复制机制：

1. 应该为一台独立的名称服务器
2. 主服务器的区域解析库文件中必须有一条 NS 记录是指向从服务器
3. 从服务器只需定义区域，而无需提供解析库文件，解析库问文件应该放置于/var/named/slaves/目录中
4. 主服务器得允许从服务器作区域传送
5. 主从服务器时间应该同步，可通过 NTP 进行
6. bind 程序的版本应该保持一致，如果无法一致，至少应该从高主低

日志文件说明

当您写入域文件时，也许您忘记了一个句号或空格或其他任意错误。

你可以从日志诊断 Linux DNS 服务器错误。BIND 服务通过/var/log/messages 上的错误，可以使用 tail 命令来查看实时错误日志，须使用-f 选项：$ tail -f /var /log/messages。

因此，当你编写域文件或修改/etc/named.config 并重新启动服务时，显示错误之后，你可以从日志中轻松识别错误类型。

定义一个主域服务器

我们知道 DNS 服务器类型有主域名服务器、辅助域名服务器和缓存域名服务器。不同于缓存域名服务器，主域名服务器和辅助域名服务器在应答过程中是处于同等地位的。

在 /etc/named.conf 的配置文件中，你可以使用如下语法定义一个主域服务器：

包含主要区域信息的文件存放在 /var/named 目录下，从 options 可知，这是一个工作目录。

注意：软件服务器或者托管面板会根据你的域名自动为你创建主域服务器信息的文件名，因此如果你的域名是 example.org，那么你主域服务器信息的文件就为 /var/named/example.org.db。

类型为 master，也就是说这是一个主域服务器。

定义一个辅助域服务器

同定义一个主域服务器一样，辅助域服务器的定义稍微有些变化

zone "ZONE_NAME" IN {

type slave;

masters I{ MASTER-IP; } ;

file “slaves/ZONE_NAME.zone”;

};

对于辅助域服务器来说，它的域名和主域服务器是一样的。上述语法里的的 slave 类型表示这是一个辅助域服务器，“masters IP Address list”表示辅助域服务器中区域文件内的信息都是通过主域服务器中区域文件内的信息复制过来的。

定义一个缓存服务器

即使你已经配置了主域或者辅助域服务器，你仍有必要（不是必须）定义一个缓存服务器，因为这样你可以减少 DNS 服务器的查询次数。

在定义缓存服务器之前，你需要先定义三个区域选择器，第一个：

zone "." IN {

type hint;

file "root.hint";

};

zone "." IN {

type hint;

file "root.hint";

};

zone "." IN {

type hint;

file "root.hint";

};

zone "localhost" IN {

type master;

file "localhost.db";

};

定义第三个区域是为了反向查找到本地主机。这种反向查找是把本地的 IP 地址执向本地主机。

zone "0.0.127.in-addr.arpa" IN {

type master;

file "127.0.0.rev";

};

把这三个区域信息放到/etc/named.conf 文件里，你的系统就可以以缓存服务器来工作了。

TXT 记录

您可以将任何信息存储到 TXT 记录中，例如你的联系方式或者你希望人们在查询 DNS 服务器时可获得的任意其他信息。

你可以这样保存 TXT 记录：example.com. IN TXT ” YOUR INFO GOES HERE”.

此外，RP 记录被创建为对 host 联系信息的显式容器：example.com. IN RP mail.example.com. example.com。

Linux DNS 解析器

我们已经知道 Linux DNS 服务器的工作原理以及如何配置它。另一部分当然是与 DNS 服务器交互的（正在与 DNS 服务器通信以将主机名解析为 IP 地址的）客户端。

在 Linux 上，解析器位于 DNS 的客户端。要配置解析器，可以检查/etc/resolv.conf 这个配置文件。

在基于 Debian 的发行版上，可以查看/etc/resolvconf/resolv.conf.d/目录。

/etc/resolv.conf 文件中包含客户端用于获取其本地 DNS 服务器地址所需的信息。

第一个表示默认搜索域，第二个表示主机名称服务器(nameserver)的 IP 地址。

名称服务器行告诉解析器哪个名称服务器可使用。 只要你的 BIND 服务正在运行，你就可以使用自己的 DNS 服务器。

## BIND 的基础安全相关配置

acl：把一个或多个地址归位一个集合，并通过一个统一的名册很难过调用，这些配置放在 named.conf 配置文件中

- acl ACL-NAME {
- IP；
- IP;
- net/perlen;
- };

bind 有 4 个内置的 acl

1. none:没有一个主机
2. any：任意主机
3. local：本机
4. localnet：本机的 IP 同掩码运算后得到的网络地址

注意：只能先定义后使用，因此，一般定义在配置文件中 options 的前面

访问控制的指令

allow-query {} 允许查询的主机，白名单

allow-transfer {} 允许区域传送的主机，白名单

allow-recursion {} 允许递归的主机，白名单

allow-update {} 允许更新区域数据库中的内容

BIND VIEW：用于同一个域名解析成多个不同区域的 IP
