---
title: Resource Record(资源记录)
---

# 概述

> 参考：

RR 定义的格式：NAME \[TTL] CLASS RR-TYPE VALUE（注意：格式中的域名都要带根域名，即域名最后都要加一个 . ）

1. NAME 和 VALUE # 不同的 RR-TYPE 有不同的格式
2. CLASS：IN
3. TYPE 资源记录类型：A，AAAA，PTR，SOA，NS，CNAME，MX 等：
   1. SRV 格式：域名系统中用于指定服务器提供服务的位置（如主机名和端口）
      1. name：\_服务.\_协议.名称.
      2. value：优先级 权重 端口 主机.
   2. SOA 格式：Start Of Authority：起始授权记录，一个区域解析库有且只能有一个 SOA 记录，而且必须为解析库第一条记录
      1. name：当前区域的名字，例如”baidu.com.“
      2. value (属性)
         1. 当前区域的主 DNS 服务器的 FQDN，也可以使用当前区域的名字
         2. 当前区域管理员的邮箱地址，但是地址中不能用@符号，@符号用.替换
         3. （主从服务协调属性的定义以及否定结果的统一的 TTL）
   3. NS 格式：Name Server：专用于标明当前区域的 DNS 服务器
      1. name：当前区域的名字
      2. value：当前区域的某 DNS 服务器的名字，例如 ns.baidu.com.;(一个区域可以有多个 NS 记录)
   4. MX 格式：Mail eXchanger：邮件交换器
      1. TTL 可以从全局继承
   5. A/AAAA 格式：Address，A 格式用于实现将 FQDN 解析为 IPv4(AAAA 格式用于将 FQDN 解析为 IPv6)
      1. name:
      2. value:
   6. PTR 格式：PoinTeR，用于将 IP 解析为 FQDN
      1. name：IP，特殊格式，反写 IP，比如 1.2.3.4 要写成 4.3.2.1，跟后缀 in-addr.arpa.
      2. value:FQDN
   7. CNAME 格式：Canonical Name，别名记录
      1. name：别名的 FQDN
      2. value：正式名字的 FQDN
4. 注意：
   1. @可用于引用当前区域的名字
   2. 相邻的两个资源记录的 name 相同时，后续的可省略
   3. 同一个名字可以通过多条记录定义多个不同的值，此时 DNS 服务器会轮循响应
   4. 同一个值也有可能有多个不同的定义名字，通过多个不同的名字指向同一个

### EXAMPLE

    baidu.com 86499        IN            SOA    ns.baidu.com.        desistdaydream.qq.com.  （
                                                                     2018072001    #序列号，当序列号变化时，即代表资源有变化，主DNS会主动同步数据给备
                                                                     2H            #刷新时间，2小时
                                                                     10M           #重试时间，10分钟
                                                                     1W            #过期时间，1周
                                                                     1D ）         #否定结果的TTL值，1天
    baidu.com.             IN            NS     ns1.baidu.com
    ns1.baidu.com          IN            A      1.1.1.0
    www.baidu.com          IN            A      1.1.1.1
                           IN            A      1.1.1.2
    baidu.com.             IN            MX 10  mx1.baidu.com
                           IN            MX 20  mx2.baidu.com
    4.3.2.1.in-addr.arpa.  IN            PTR    www.baidu.com
    web.baidu.com.         IN            CNAME  www.baidu.com.

子域授权：每个域的 DNS 服务器，都是通过在其上级 DNS 服务器中的解析库添加该域的 DNS 服务器信息进行授权

EXAMPLE，在根域的 DNS 服务器中，记录了.com.域的资源记录，类似下面的方式，不是绝对的
.com. IN NS ns1.com. 定义.com.域的域名服务器为 ns1.com.
.com. IN NS ns2.com.
ns1.com. IN A 2.2.2.1 定于.com.域中的域名服务器 ns1.com.的 IP 地址为 2.2.2.1
ns2.com. IN A 2.2.2.2

当www.baidu.com.的DNS请求到根的DNS服务器的时候，根的DNS服务器查找自己解析库中.com的域中的DNS服务器资源，然后看到该DNS服务器所对应的IP，然后把该请求转发到.com域中的DNS服务器进行下一步解析，然后.com域的DNS服务器在从解析库中找到baidu的资源再转发到baidu的DNS服务器上(或者直接返回baidu的IP地址)

以上例子是.com 域在.根域中解析库中的资源记录，如果还有 baidu.com 的域名，则该域的资源记录写在.com 域中的解析库，以此类推

反向区域：

区域名称：网络地址反写.in-addr-arpa.

EXAMPLE：172.16.100. 写成 100.16.172.in-addr-arpa.
