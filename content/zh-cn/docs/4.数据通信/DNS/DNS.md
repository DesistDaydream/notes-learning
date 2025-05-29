---
title: DNS
---

# 概述

> 参考：
>
> - [Wiki, DNS](https://en.wikipedia.org/wiki/Domain_Name_System)
> - [Wiki, Name Server](https://en.wikipedia.org/wiki/Name_server)
> - 《DNS 与 BIND》(第 5 版)

**Domain Name System(域名系统，简称 DNS)** 是一个分层的和去中心化的命名系统，以便人们可以更方便得访问互联网。DNS 主要用来将更容易让人类记忆的 `域名` 与 `IP地址` 互相映射，以便可以通过域名定位和访问互联网上的服务。自 1985 年以来，域名系统通过提供全球性、分布式的域名服务，已成为 Internet 功能的重要组成部分。

从另一个方面说，DNS(域名系统) 其实是一个**分布式的数据库**。这种结构允许对整体数据库的各个部分进行本地控制，并且在各个部分中的数据通过 C/S 模式变得对整个网络都可用。通过复制和缓存等机制，DNS 将会拥有健壮性和充足的性能。

> 注：这段描述不好理解，需要看完后面才能体会。所谓的分布式，其实就是指 DNS 的模式，每个 Name Server 都可以是分布式数据库的一个节点。

当我们访问一个网站时，在浏览器上输入 `https://www.baidu.com/` 时，`www.baidu.com` 就是域名。而管理 域名与 IP 对应关系的系统，就是域名系统。

DNS 其实是一个规范、概念，具体想要让 DNS 在世界上应用起来，则至少要保证两个方面

- 其一是服务端，通过 NameServer 为大家提供解析服务、存储域名与 IP 的对应关系，
- 其二是客户端，客户端上的应用程序将会调用符合 DNS 标准的库以便向 NameServer 发起域名查询请求，程序收到解析后的 IP 后将会发起请求。

## 背景

网络诞生之初并没有 DNS，那时候访问对方只需要 IP 地址就可以了，但是后来接入互联网的主机太多了，IP 没法记，所以研究了 DNS。

**Internet Assigned Numbers Authority(互联网数字分配机构，简称 IANA)** 是负责协调一些使 Internet 正常运作的机构。同时，由于 Internet 已经成为一个全球范围的不受集权控制的全球网络，为了使网络在全球范围内协调，存在对互联网一些关键的部分达成技术共识的需要，而这就是 IANA 的任务

IANA 的所有任务可以大致分为三个类型：
一、域名。IANA 管理 DNS 域名根和.int，.arpa 域名以及 IDN（国际化域名）资源。
二、数字资源。IANA 协调全球 IP 和 AS（自治系统）号并将它们提供给各区域 Internet 注册机构。
三、协议分配。IANA 与各标准化组织一同管理协议编号系统。

## NameServer

**NameServer(名称服务器)** 是 DNS 中最重要的概念，有的时候也称为 **DNS Server**。Name Server 是 DNS 的一个组件。该组件的最重要功能就是将人类记忆的域名解析为 IP 地址。

> 简单得说，想要让全世界都接入 DNS(域名系统)，就需要有一个服务，这个服务就是一个翻译器。比如我要访问谷歌，就会问 NameServer：谷歌在哪。此时这个翻译器(NameServer)就会告诉我谷歌的 IP 地址，这时我再去访问这个 IP 地址即可。

可以实现上述的程序有很多，凡是可以实现 DNS 功能的程序，我们一般成为 **NameServer(名称服务器)：**

- **Bind** # 最常用 NameServer，Bind 的具体介绍详见 《[BIND](/docs/4.数据通信/DNS/BIND/BIND.md)》章节。
- **DNSmasq** # 轻量的 NameServer，用于提供 DNS 缓存功能。
- ...... 等等

一个 NameServer 通常包含了整个 DNS 数据库中的某些部分的信息，并让这些信息可以被客户端所用(客户端通常称为 **Resolver(解析器)**)。解析器通常只是一组 Library，这些库产生查询请求，并将请求通过网络发送给名称服务器。

- 比如 Linux 中 curl、ping 等等命令就会调用 Resolver，毕竟这些命令如何请求一个域名，就需要要知道对应的 IP 才可以。
- 解析器一般是被内嵌在系统中，当在系统中运行任何需要域名解析的程序时，都会调用这个解析器。

**用白话说，Name Server 通常是一个默认监听在 53/UDP 与 53/TCP 上的服务，用来处理客户端发来的域名解析请求(根据自身数据库中 [RR](#Resource%20Record(资源记录)) 解析域名)，或者向上级 Name Server 发起域名查询请求。比如 8.8.8.8、114.114.114.114 等等，都属于 NameServer**

**NameServer 是实现 DNS 的具体实现**

注意：

- 上述描述的都是自建 Name Server。如果我们在一个域名注册商出购买域名后，都需要指定用来处理自己购买域名的 Name Server，只不过通常情况，域名注册商本身自己就有 Name Server，这样也方便进行身份验证(即域名属于自己)。
- 但是也有例外，比如免费的域名注册商 eu.org 就没有自己的 Name Server，在 eu.org 出买的域名必须要指定其他的 Name Server。

### NameServer 类型

> 这部分内容需要理解域名结构后，再回来看。否则无法理解为什么 NameServer 会有这么多分类

1. 主 DNS Server：维护其所负责解析的域内解析库服务器，解析库由管理员维护
2. 辅助 DNS Server：从主 DNS Server 或其他辅助 DNS Server 那里复制(区域传送)一份解析库
   1. 序列号：解析库的版本号，主服务器的解析库内容发生变化，其序列号发生变化，当序列号发生变化的时候，辅助 DNS Server 则去复制解析库
   2. 区域传送：
      1. 全量传送：传送整个解析库
      2. 增量传送：传送解析库变化的那部分内容
   3. 刷新时间间隔：辅助 DNS Server 从主 DNS Server 请求同步解析库的时间间隔
   4. 重试时间间隔：辅助 DNS Server 从主 DNS Server 请求同步解析库失败时，再次尝试的时间间隔
   5. 过期时长：辅助 DNS Server 始终联系不到主 DNS Server 时，多久之后放弃辅助 DNS Server 职责，变成主服务器
3. 缓存 DNS Server：这些服务器上不存放特定域名的配置文件。当客户端请求缓存服务器来解析域名时，该服务器将首先检查其本地缓存。如果找不到匹配项便会询问主服务器。接着这条响应将被缓存起来。您也可以轻松地将自己的系统用作缓存服务器，如果该域内不存在主 DNS 服务器，那么则直接去找根域的 DNS 服务器进行域名解析。
4. 转发器 Server：该服务器会转发收到的域名解析请求到别的服务器，这杯被转发到的服务器需要能够为请求做递归处理，否则，转发请求不予进行
   1. 全部转发：凡是对非本机所有负责解析的区域的请求，统统转发给指定的服务器
   2. 区域转发：仅转发对特定的区域的请求至某服务器

# DNS 架构与概念

![1024px-Domain_name_space.svg.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/msw4yv/1617092212523-dcc0c65c-8807-4964-8766-8613f94c8cd0.png)
DNS(域名系统) 是一个树型架构，这颗大树被称为 **Domain Namespace(域名称空间)**，树中的每个节点都具有 **Label(标签)** 和 0 个或多个 **ResourceRecords(资源记录，简称 RR)**，其中保存与域名相关联的信息。域名本身由 Label(标签) 组成，多个标签通过 `.` 符号串联而成，`.` 符号右侧的标签就是其左侧标签的父节点名称。

> 这里所描述的树形结构中的节点，实际上就是一个个的 NameSever

## Domain(域) 的概念

根负责管理他下面的一些域，这些域再负责其所在域的管理，依次类推，域就是该 NameServer 可管辖的范围。而被管的域则称为 **SubDomain(子域)**。在树型结构中，也可以说低级节点就是高级节点的 SubDomain(子域)。

Domain 还可以用现实中的地理来类比，比如我前文所说的我的名字”四.李.海淀区.北京市.中国.“，根域里包括多个国家域，每个国家域又包括很多城市域，每个城市域包括很多区域，每个区域里有很多姓氏域，姓氏域中，名为李姓氏域中，叫四的就是我(某台设备的名字就是这个区域的一台)

我们光有 Domain 的概念还不行，在真正使用上，还需要给每个域起一个名称，就叫 **DomainName(域名)**。所以，DNS 这整套系统，就是围绕着 Domain 来进行的。而**每个 Domain 的名称，就是其 Label 的名称。**

## Zone(区) 的概念

DNS 这个树型结构还使用了一个 **Zone(区)** 的概念来进行分区管理，一个 Zone 可以包含一个域，也可以包含很多域及其子域。Zone 可以包含的内容具体取决于 **Delegation(授权)**。
详见：[Zone 与 Domain](/docs/4.数据通信/DNS/Zone%20与%20Domain.md)

## DomainName(域名) 的结构

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/msw4yv/1616161023756-88e23654-7efc-4c25-826f-632042129088.jpeg)
域名由一个或多个部分组成，每个部分由 `.` 符号连接。组成域名的每个部分称为 **Label(标签)**。这个 Label 就是 NameServer 的 Label。
结构如下：`Label-N...Label-3.Label-2.Label-1.`
最右侧的空标签是为 root 保留的，长度为 0。其余每个 Label 最多可以写 63 个文本字符。

1. **RootDomain(根域名)**# 空标签，管理根以下的所有顶级域名，全世界一共有 13 组根域名 Server，不包括中国
   1. 实际上，当你访问任何网站时，浏览器会默认在域的末尾添加一个不可见的点，因此域名像 `www.baidu.com.` 一样。 最后边的点被称为根域。
   2. 注意，任何 NameServer 在本地数据库没有查到请求的域名的时候，必然会去根 NameServer 请求，不会去顶级或者二级直接发请求，这是规定
3. **Top-Level Domain(顶级域名)** # Label-1，管理该域名下的所有二级域名。比如 `com.`、`net.` 等都是顶级域名
4. **二级域名** # Label-2，常用来标识该网站名称，比如 baidu.com、google.com
5. **三级域名** # Label-3，常用来标识域名作用的，比如提供网页服务的用 www.baidu.com、提供邮箱服务的用 mail.baidu.com
6. **其他级域名** # 理论上来说，一共可以划分 127 个等级。

比如 `www.baidu.com` 实际上就是 `www.baidu.com.`。baidu.com. 就是二级域名，www.baidu.com. 就是三级域名。

> 这就好比 Linux 的目录树结构一样，只不过 Linux 中，域 称为 目录，子域 称为 子目录。最顶层的称为 根目录，而且分隔符是 `/`。只不过 Linux 的结构是从左至右的。
> 域名就是区域的名字，就好像本人的名字李四一样，我是四.李.南开区.天津市.中国.，域名就是这样的结构。

比如我现在想知道 `www.baidu.com` 在哪：

1. 该域名查询的时候，先去问根.，根告诉我这个归 .com 管
2. 然后去问 .com，.com 告诉我这个域名归 baidu 管
3. 然后去问 .baidu，.baidu 告诉我，这是他们的 web 服务器，IP 地址是这个，然后我就可以跟 www.baidu.com 建立链接了

**Fully Qualified Domain Name(完全限定域名，简称 FQDN)** 。类似于 `www.baidu.com.` 就是完全限定域名，即该域名直接代表一台提供业务服务的设备。这就好比 Linux 中的绝对路径概念一样。其他的都是相对域名

## Resource Record(资源记录)

与域名相关的数据都被保存在每个 NameServer 的 **ResourceRecord(资源记录，简称 RR)** 中。域名与 IP 之间的对应关系，称为 **Record(记录)**。

RR 用于标示 NameServer 解析库中每一个解析条目的具体类型，比如本域的域名与 IP 的对应关系，下级域的域名中的资源等等，标注出这些资源，以便用这些资源解析用户的 DNS 请求

RR 具有标准格式，遵守 RFC1035，每一条记录都可以包含如下几个字段：

- **NAME** # 本条记录中节点的完全限定域名
- **CLASS** # 本条记录的类。对于互联网来说，通常为 IN
- **TYPE**# 本条记录的类型
- **TTL**# 本条记录的保存时长
- **RDLENGTH**# RDATA 字段的长度
- **RDATA**# 本条记录的额外内容
- **VALUE** # 本条记录的具体内容。不同的记录类型，该字段的内容也不同。

RR 定义的格式：NAME \[TTL] CLASS RR-TYPE VALUE（注意：格式中的域名都要带根域名，即域名最后都要加一个 . ）

- NAME 和 VALUE # 不同的 RR-TYPE 有不同的格式
- CLASS：IN
- TYPE 资源记录类型：A，AAAA，PTR，SOA，NS，CNAME，MX 等：
  - SRV：域名系统中用于指定服务器提供服务的位置（如主机名和端口）
    - name # \_服务.\_协议.名称.
    - value # 优先级 权重 端口 主机.
  - SOA：Start Of Authority：起始授权记录，一个区域解析库有且只能有一个 SOA 记录，而且必须为解析库第一条记录
    - name # 域名，例如”baidu.com.“
    - value # (属性)
      - 当前区域的主 DNS 服务器的 FQDN，也可以使用当前区域的名字
      - 当前区域管理员的邮箱地址，但是地址中不能用@符号，@符号用.替换
      - （主从服务协调属性的定义以及否定结果的统一的 TTL）
  - NS：Name Server：专用于标明当前区域的 DNS 服务器
    - name # 域名
    - value # 当前区域的某 DNS 服务器的名字，例如 ns.baidu.com.;(一个区域可以有多个 NS 记录)
  - MX：Mail eXchanger：邮件交换器
    - TTL 可以从全局继承
  - A/AAAA：Address，A 格式用于实现将 FQDN 解析为 IPv4(AAAA 格式用于将 FQDN 解析为 IPv6)
    - name # 域名
    - value # 域名对应的 IP 地址
  - PTR：PoinTeR，用于将 IP 解析为 FQDN
    - name # IP，特殊格式，反写 IP，比如 1.2.3.4 要写成 4.3.2.1，跟后缀 in-addr.arpa.
    - value # FQDN
  - CNAME：Canonical Name，别名记录
    - name # 别名的 FQDN
    - value # 正式名字的 FQDN
- 注意：
  - @可用于引用当前区域的名字
  - 相邻的两个资源记录的 name 相同时，后续的可省略
  - 同一个名字可以通过多条记录定义多个不同的值，此时 DNS 服务器会轮循响应
  - 同一个值也有可能有多个不同的定义名字，通过多个不同的名字指向同一个

### EXAMPLE

```bash
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
```

子域授权：每个域的 DNS 服务器，都是通过在其上级 DNS 服务器中的解析库添加该域的 DNS 服务器信息进行授权

EXAMPLE，在根域的 DNS 服务器中，记录了.com.域的资源记录，类似下面的方式，不是绝对的

- .com. IN NS ns1.com. 定义.com.域的域名服务器为 ns1.com.
- .com. IN NS ns2.com.
- ns1.com. IN A 2.2.2.1 定于.com.域中的域名服务器 ns1.com.的 IP 地址为 2.2.2.1
- ns2.com. IN A 2.2.2.2

当www.baidu.com.的DNS请求到根的DNS服务器的时候，根的DNS服务器查找自己解析库中.com的域中的DNS服务器资源，然后看到该DNS服务器所对应的IP，然后把该请求转发到.com域中的DNS服务器进行下一步解析，然后.com域的DNS服务器在从解析库中找到baidu的资源再转发到baidu的DNS服务器上(或者直接返回baidu的IP地址)

以上例子是.com 域在.根域中解析库中的资源记录，如果还有 baidu.com 的域名，则该域的资源记录写在.com 域中的解析库，以此类推

反向区域：

区域名称：网络地址反写.in-addr-arpa.

EXAMPLE：172.16.100. 写成 100.16.172.in-addr-arpa.

# DNS 解析过程

DNS 查询过程：主机发送请求到根域名解析服务器，然后重定向到二级域名解析服务器，再重定向到三级域名解析服务器，以此类推

- 在本机上查询 DNS 的配置文件(比如/etc/hosts)，有没有 IP 地址与 Domain Name 的对应关系
  - EXAMPLE：如果把本机 IP 的对应域名改 baidu.com.，那么在访问百度的时候，就只会访问本机了
- 如果在本机无法查询到 Domain Name 与 IP 的对应关系，那么需要通过 DNS 代理来进行查询，总共分为两种查询类型
  - 递归查询：主机只发送一次 DNS 解析请求，就获得最后的结果。
    - 在本机配置一个运行 DNS 服务的 Server 的 IP 地址，把请求直接发送给该 server，
    - 由该 server 去找.根域名服务器进行查询，然后.根域名服务器再根据该请求中的顶级域名把该请求重定向到顶级域名服务器上
    - 如果该请求还有二级域名，那么顶级域名 server 会 再把该请求重定向二级域名 server 上
    - 直到查询到最终结果后，把该结果返回给 DNS Server，然后 DNS Server 把结果直接告诉发送请求的主机
  - 迭代查询：主机发送一次 DNS 解析请求后，被重定向到另一台 DNS 服务器继续发送请求，直到获得最后结果。
    - 该查询主机直接发送请求到.根域名 server，然后进行递归查询中的 2,3,4 步骤

一次完整的查询请求经过的流程

- 首先，客户端查询本地的 DNS 服务器。默认情况下，即/etc/resolv.conf 文件中所列的第一个名称服务器。
- 接着本地名称服务器会查询本地解析库和缓存，如果本地数据库中有该资源记录（该名称服务器对该域进行权威解析），则返回查询结果。
- 如果没有，则查询缓存，看看缓存中是否有以前对该资源记录的查询结果，如有，则返回查询结果；如果仍然没有，则会向其他 DNS 服务器进行递归解析。
- 进入递归解析，本地 DNS 服务器向根域 DNS 服务器（Root Nameserver）提出查询请求，根域 DNS 服务器会返回顶级域（TLD）DNS 服务器地址（例如.com 的 DNS 服务器地址）；本地 DNS 服务器再次向 TLD Nameserver 发出查询请求，TLD Nameserver 会返回下一级域的 DNS 服务器地址；依此类推，直到查询到权威的名称服务器（Authoritative Nameserver）;

解析类型：

1. Name --> IP 正向解析
2. IP --> Name 反向解析

## 域名解析结果

DNS-Rcode 作为 DNS 应答报文中有效的字段，主要用来说明 DNS 应答状态，这可是小编排查域名解析失败的重要指标。通常常见的 Rcode 值如下：

- Rcode 值为 0，对应的 DNS 应答状态为 NOERROR，意思是成功的响应，即这个域名解析是成功
- Rcode 值为 2，对应的 DNS 应答状态为 SERVFAIL，意思是服务器失败，也就是这个域名的权威服务器拒绝响应或者响应 REFUSE，递归服务器返回 Rcode 值为 2 给 CLIENT
- Rcode 值为 3，对应的 DNS 应答状态为 NXDOMAIN，意思是不存在的记录，也就是这个具体的域名在权威服务器中并不存在
- Rcode 值为 5，对应的 DNS 应答状态为 REFUSE，意思是拒绝，也就是这个请求源 IP 不在服务的范围内

### DNS 请求失败的具体分析

常见的请求失败包括：
1、域名记录不存在，即 Rcode 值为 3（NXDOMAIN）的情况，这种情况下域名权威服务器及托管的主域名均正常，但是权威并不存在这条具体的域名记录，于是权威返回了 NXDOMAIN，值的注意的是这个 NXDOMAIN 的报文中会包含一个 AUTHORITY SECTION，内容为改主域名的 SOA 记录，这个应答结果会在递归服务器中被缓存，缓存时间周期为域名的 SOA 记录的 TTL：
![](https://notes-learning.oss-cn-beijing.aliyuncs.com/msw4yv/1616161023846-76ea5416-093f-4be5-93ca-8f2ec94d745e.jpeg)
2、权威解析失败，即 Rcode 值为 2（SERVFAIL）的情况，递归服务器会给请求源这个结果的原因是向权威解释请求异常，包括且不限于权威不响应/或者权威返回 refuse/或者权威返回 servfail，这个 SERVFAIL 的应答结果当然是一个空结果，不过 BIND 会强制给这个结果增加一个 1S 的 TTL，所以 SERVFAIL 的应答会在递归服务器中被缓存，缓存时间周期为 1S

- 2.1）权威不响应。包括递归服务器至权威服务器中间的网络异常在内，递归服务器在发出递归请求并完成重试超时后，给请求源一个 SERVFAIL 的应答，并缓存 1S ：

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/msw4yv/1616161023768-36d66d1b-0e60-4295-9fd7-899f67b2fe9b.jpeg)

- 2.2）权威向递归服务器应答 REFUSE。当权威服务器不存在主域名及对应的 SOA 记录时，权威会向递归服务器返回 REFUSE，即不在我服务的范围内拒绝，递归服务器在收到这个 REFUSE 应答后，给请求源一个 SERVFAIL 的应答，并缓存 1S：

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/msw4yv/1616161023786-286f827b-fc20-42ba-8997-e83267a91a55.jpeg)

- 2.3）权威向递归服务器应答 SERVFAIL。当权威服务器存在主域名但是由于 zonefile 被破坏导致权威服务器上域名的 NS 记录异常时，权威会向递归服务器返回 SERVFAIL，即解析失败，递归服务器在收到这个 SERVFAIL 应答后，给请求源一个 SERVFAIL 的应答，并缓存 1S：

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/msw4yv/1616161023782-ace75691-f594-4871-a17a-ee667e422bde.jpeg)

- 2.4）权威向递归服务器应答其他的错误 Rcode。由于不常见本文就不展开了，递归服务器在收到其他错误应答后，给请求源一个 SERVFAIL 的应答，并缓存 1S：

3、拒绝服务，即 Rcode 值为 5（REFUSE）的情况。除了记录不存在（NXDOMAIN）和解析失败（SERVFAIL）以外，如果请求源不在递归服务器的服务范围内，这种情况下递归服务器会直接给请求源一个 REFUSE 的应答，本地直接应答无缓存：
![](https://notes-learning.oss-cn-beijing.aliyuncs.com/msw4yv/1616161023827-0e3a338b-f567-46de-96fc-bf08dbd86720.jpeg)
4、响应成功，但是没有解析结果，这是一种比较特殊的情况，这种情况是 Rcode 值为 0（NOERROR）的情况。这种情况下域名权威服务器及托管的主域名均正常，权威本身也存在这条具体的域名记录，但是没有对应的记录类型（不包含 CNAME，CNAME 是特殊情况，可以响应任意类型的请求），这是权威返回了 NOERROR，值的注意的是这个 NOERROR 的报文中没有 ANSWER SECTION。但是会包含一个 AUTHORITY SECTION，内容为改主域名的 SOA 记录，这个应答结果会在递归服务器中被缓存，缓存时间周期为域名的 SOA 记录的 TTL：
![](https://notes-learning.oss-cn-beijing.aliyuncs.com/msw4yv/1616161023799-a2ec8dd5-576d-42db-8c4f-224070eb0e21.jpeg)
5、还有最后的一种情况，就是递归服务器本身不响应了，这个比较容易理解，如果递归服务器不响应，那么请求段收不到任何应答，这个时候请求端终端如果有超时机制则会跑出一个 dns 请求 timeout 的结果：
![](https://notes-learning.oss-cn-beijing.aliyuncs.com/msw4yv/1616161023809-badb7042-8c8f-44c0-888f-802672eddd6b.jpeg)
结论：
![](https://notes-learning.oss-cn-beijing.aliyuncs.com/msw4yv/1616161023767-6dcbac05-682c-4b55-8225-1e50694e0f68.jpeg)

# Domain name registrar(域名注册商)

> 参考：
>
> - [Wiki, Domain name registrar](https://en.wikipedia.org/wiki/Domain_name_registrar)
> - [Name.com](https://www.name.com/)
> - [eu.org](https://nic.eu.org/arf/en/)
> - <https://iweec.com/144.html>

**Domain name registrar(域名注册商**) 是一个商业实体或组织，它们由[互联网名称与数字地址分配机构](https://zh.wikipedia.org/wiki/%E4%BA%92%E8%81%94%E7%BD%91%E5%90%8D%E7%A7%B0%E4%B8%8E%E6%95%B0%E5%AD%97%E5%9C%B0%E5%9D%80%E5%88%86%E9%85%8D%E6%9C%BA%E6%9E%84)（ICANN）或者一个国家性的 [Country code top-level domain](https://en.wikipedia.org/wiki/Country_code_top-level_domain) (国家代码顶级域名，简称 ccTLD) [域名注册局](https://zh.wikipedia.org/wiki/%E5%9F%9F%E5%90%8D%E6%B3%A8%E5%86%8C%E5%B1%80)委派，以在指定的[域名注册数据库](https://zh.wikipedia.org/wiki/%E5%9F%9F%E5%90%8D%E6%B3%A8%E5%86%8C%E6%95%B0%E6%8D%AE%E5%BA%93)中管理[互联网](https://zh.wikipedia.org/wiki/%E4%BA%92%E8%81%94%E7%BD%91)[域名](https://zh.wikipedia.org/wiki/%E5%9F%9F%E5%90%8D)，向公众提供此类服务。

国内域名注册商：

- 万网
- 新网
- DNSPod

国外域名注册商：

- eu.org # 没有自己的 Name Server
- godaddy
- Name.com
- freenom

## eu.org 注册方式

> 参考：
>
> - <https://iweec.com/363.html>
> - <https://www.bilibili.com/video/BV1JB4y1m7e9>

eu.org 免费域名从 1996 年就有了，由此可见是非常非常早，计划是专门给无力承担费用的一些组织使用的，现在我们来申请一个。

注册地址：<https://nic.eu.org/arf/en/login/?next=/arf/en/>

点击：Register

填表挺简单的，可以参考我的视频教程。

B 站：<https://www.bilibili.com/video/BV1JB4y1m7e9/>

Youtube：<https://www.youtube.com/watch?v=xWgeCUpM81I>

然后成功后到邮箱（垃圾箱）找到邮件、激活，然后登录。

点击登录，然后点击：New domain

填写理想的完整域名例如：abcde.eu.org 同意协议；

域名服务器建议填写下面两个 dnspod，否则无法转到 cloudflare。

edmund.dnspod.net

dempsey.dnspod.net

这里先只选择 server names，然后 Submit，注意看检查页面，如果出现 No error,Done.说明成功了~

如果有错误，请返回修改！

域名审核 1 天——30 天都有可能，所以慢慢等吧，经过我的测试，一个账号内最多可以申请 4 个免费域名。

域名投资参考：<https://iweec.com/144.html>

eu.org 域名通过后是有邮件通知的，但是都在垃圾箱。若是想转到 cloudflare 出现问题，请参考视频：

Bilibili：<https://www.bilibili.com/video/BV1ST4y1z7Ra/>

Youtube：<https://www.youtube.com/watch?v=EOsBJxtiOho>

注意：

- eu.org 本身没有自己的 Name Server，在这里免费获取到域名后，需要指定其他的 Name Server(比如 DNSPod 中的 NameServer)

## 可以为其他域名提供 NameServer 的域名注册商

DNSPod https://www.dnspod.cn/

# 国内域名备案说明

在国内的域名注册机构是必须要备案的，备案时需要关联自己的服务器或者通过公司进行备案。如果不备案，那么域名只有解析 IP 的功能，无法在公网被访问，效果如下：

![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/msw4yv/1654614073447-ad968e6f-5cda-4fa1-a31a-84c148d8e6a8.png)

若访问的是 https，则会提示

![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/msw4yv/1656167438257-cf4f1287-4ce6-4ca1-b7b7-3814937ebe74.png)

这个是由于在国内的服务器，都会收到服务器所在 IDC 的限制，这些 IDC 会在最外层部署一套检测服务，用以检查每个标准端口(80 和 443)的请求域名是否已经备案，若没有备案，该请求 IDC 则不会放行到服务器上。

DNS 污染、GFW、SNI 阻断，这是三座大山
