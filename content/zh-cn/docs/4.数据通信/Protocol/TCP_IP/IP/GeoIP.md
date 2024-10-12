---
title: GeoIP
linkTitle: GeoIP
date: 2024-10-08T15:06
weight: 20
---

# 概述

> 参考：
>
> - [Wiki, GeoIP](https://en.wikipedia.org/wiki/GeoIP)
>   - 重定向到 [Internet geolocation](https://en.wikipedia.org/wiki/Internet_geolocation)

在 IT 中，**Internet geolocation(互联网地理定位，简称 GeoIP)** 是能够推断连接到互联网的设备的地理位置的软件。例如，设备的 IP 地址可用于确定国家、城市或邮政编码，从而确定其地理位置。其他方法包括检查 Wi-Fi 热点、

# IP 地址分配机制

> 参考：
>
> - [IANA, 号码资源](https://www.iana.org/numbers)
> - [APNIC, 搜索](https://wq.apnic.net/static/search.html)(通过给定的 IP 地址搜索谁拥有这个 IP)
> - [面包板，你知道中国大陆一共有多少 IPv4 地址吗？](https://www.eet-china.com/mp/a54338.html)
> - [公众号，k8s 中文社区-居然还有 2 亿多 IPv4 地址未分配](https://mp.weixin.qq.com/s/GHYYgZwAuEV4qPCwdI8Bjg)

IPv4 和 IPv6 地址通常以分层方式分配。**ISP(互联网服务提供商)** 为用户分配 IP 地址。ISP 从 **LIR(本地互联网注册机构)** 或 **NIR(国家互联网注册机构)** 或 **RIR(相应的区域互联网注册机构)** 获取 IP 地址分配

![500](https://notes-learning.oss-cn-beijing.aliyuncs.com/ip/ip_rir.png)

| 登记处                                | 覆盖面积                                                    |
| ---------------------------------- | ------------------------------------------------------- |
| [AFRINIC](http://www.afrinic.net/) | Africa Region(非洲地区)                                     |
| [APNIC](http://www.apnic.net/)     | Asia/Pacific Region(亚洲/太平洋地区，亚太地区)                      |
| [ARIN](http://www.arin.net/)       | Canada, USA, and some Caribbean Islands(加拿大、美国、一些加勒比岛屿) |
| [LACNIC](http://www.lacnic.net/)   | Latin America and some Caribbean Islands(拉丁美洲、一些加勒比岛屿)  |
| [RIPE NCC](http://www.ripe.net/)   | Europe, the Middle East, and Central Asia(欧洲、中东、中亚)     |

对 IP 地址的主要作用是根据[全球政策](http://www.icann.org/en/general/global-addressing-policies.html)所述的需求将未分配地址池分配给 RIR，并记录 [IETF](/docs/Standard/Internet/IETF.md) 所做的协议分配。当 RIR 需要在其区域内分配或分配更多 IP 地址时，我们会向 RIR 进行额外分配。我们不会直接向 ISP 或最终用户进行分配，除非在特定情况下，例如分配多播地址或其他协议特定需求。

APNIC 是全球 5 个地区级的 Internet 注册机构（RIR）之一，负责亚太地区的以下事务：

1. 分配 IPv4 和 IPv6 地址空间，AS 号；
2. 为亚太地区维护 Whois 数据库；
3. 反向 DNS 指派；
4. 在全球范围内作为亚太地区的 Internet 社区的代表。

所以，中国大陆境内的地址都会登记在 APNIC 的地址库内。地址库获取方式: http://ftp.apnic.net/apnic/stats/apnic/delegated-apnic-latest

例如在 Linux 系统中，使用 wget 命令可以拉取 delegated-apnic-latest 文件。

```bash
~]# wget http://ftp.apnic.net/apnic/stats/apnic/delegated-apnic-latest
--2024-10-08 15:27:32--  http://ftp.apnic.net/apnic/stats/apnic/delegated-apnic-latest
Resolving ftp.apnic.net (ftp.apnic.net)... 203.119.102.40, 2001:dd8:8:701::40
Connecting to ftp.apnic.net (ftp.apnic.net)|203.119.102.40|:80... connected.
HTTP request sent, awaiting response... 200 OK
Length: 3875233 (3.7M) [text/plain]
Saving to: ‘delegated-apnic-latest’

delegated-apnic-latest             100%[===============================================================>]   3.70M   550KB/s    in 8.5s

2024-10-08 15:27:42 (448 KB/s) - ‘delegated-apnic-latest’ saved [3875233/3875233]

~]# ls -lh
total 3788
-rw-r--r-- 1 root root 3.7M Oct  7 23:17 delegated-apnic-latest
```

文件内容条目参考如下：

```text
apnic|JP|asn|173|1|20020801|allocated
apnic|ID|ipv4|43.240.228.0|1024|20140818|allocated
apnic|HK|ipv6|2001:df5:b800::|48|20140514|assigned
```

条目格式如下：

`注册机构|国家代码|类型|起始位|长度|分配日期|状态`

- **注册机构**：亚太地区一般为 apnic
- **国家代码**：ISO-3166 定义的两位国家或地区代码，如中国为 CN
- **类型**：asn（Autonomous System Number，自治系统编号），也就是 BGP 的 AS 编号；ipv4，IPv4 地址；ipv6，IPv6 地址
- **起始位**：第一个 ASN 编号或 IP 地址
- **长度**：从第一个起始位开始，申请分配多少的编号或地址
- **分配日期**：国家或地区向 APNIC 申请的日期
- **状态**：allocated 和 assigned，都是已分配

所以，需要将 delegated-apnic-latest 文件中所有国家为 CN、且类型为 ipv4 的条目导出，并转换为静态路由格式。

例如使用命令将符合条件的条目导入到 china 文件中。

```bash
~]# cat delegated-apnic-latest | grep CN > chinaCN
~]# cat chinaCN | grep ipv4 > china
~]# tail china
apnic|CN|ipv4|223.223.176.0|2048|20100813|allocated
apnic|CN|ipv4|223.223.184.0|2048|20100813|allocated
apnic|CN|ipv4|223.223.192.0|4096|20100806|allocated
apnic|CN|ipv4|223.240.0.0|524288|20100803|allocated
apnic|CN|ipv4|223.248.0.0|262144|20100713|allocated
apnic|CN|ipv4|223.252.128.0|32768|20110131|allocated
apnic|CN|ipv4|223.254.0.0|65536|20100723|allocated
apnic|CN|ipv4|223.255.0.0|32768|20100810|allocated
apnic|CN|ipv4|223.255.236.0|1024|20110311|allocated
apnic|CN|ipv4|223.255.252.0|512|20110414|allocated
```

可以查看文件行数，代表有多少条明细条目

```bash
~]# wc -l china
8657 china
```

然后根据起始位和长度，转换出静态路由所需的目的地址和掩码即可。在 excel 中通过对长度进行函数运算，可以得到掩码长度，如：`=32-LOG(E2,2)`，代入 2048 的话，可得到掩码长度为 21。操作后得到类似下图的表格：

![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/ip/1646295854710-c698ce59-e5fb-4e59-97cc-77a1a81bfc81.png)

先将表格内容复制到记事本中，再从记事本粘贴到 Word 中，即可得到带有内容字段、tab 制表符和段落标记的内容。如下：

- 1.0.1.0 CN 24 apnic
- 1.0.2.0 CN 23 apnic
- 1.0.8.0 CN 21 apnic

这就简单了，使用 Word 的替换功能，对对应字段进行替换就可以得到形如下文的配置：

- int loop 1
- ip add 1.12.0.1 14
- int loop 2
- ip add 1.24.0.1 13
- int loop 3
- ip add 1.48.0.1 15
- int loop 4
- ip add 1.56.0.1 13
- int loop 5
- ip add 1.68.0.1 14

再把配置分别刷入到 11 台设备当中，配置好 OSPF 和 BGP 就可以了。

![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/ip/1646295854663-30f55fe1-a908-42e8-8b91-e95b33552417.png)

## IP 应用场景

| 标记 | 中文     | 描述                                                    |
| ---- | -------- | ------------------------------------------------------- |
| ANY  | 任播网络 | 属于数据中心的一部分，任播网络；如：8.8.8.8             |
| CDN  | 内容分发 | 属于数据中心的一部分，内容分发网络                      |
| COM  | 商业公司 | 以盈利为目的的公司                                      |
| DNS  | 域名解析 | 用户提供域名解析服务的 IP；如：8.8.8.8，114.114.114.114 |
| EDU  | 教育机构 | 学校/教育机构使用的 IP                                  |
| GTW  | 企业专线 | 固定 IP，中大型公司专线上网的 IP                        |
| GOV  | 政府机构 | 政府单位使用的 IP                                       |
| DYN  | 动态 IP  | 家庭住宅用户使用的 IP                                   |
| IDC  | 数据中心 | 机房/云服务商使用的 IP                                  |
| IXP  | 交换中心 | 网络交换中心使用的 IP                                   |
| MOB  | 移动网络 | 基站出口 IP（2G/3G/4G/5G）                              |
| NET  | 基础设施 | 网络设备骨干路由使用的 IP                               |
| ORG  | 组织机构 | 非营利性组织机构                                        |
| SAT  | 卫星通信 | 通过卫星上网的出口 IP                                   |
| VPN  | 代理网络 | 属于数据中心的一部分，专门做 VPN 业务的                 |

# 最佳实践

> 参考:
>
> - [公众号 -  golang学习记，Go实现全球IP属地，这个强大到离谱的离线库准确率高达99.99%](https://mp.weixin.qq.com/s/_Zjg2HgKlPzZVDwdUj88vQ)

[GitHub 项目，lionsoul2014/ip2region](https://github.com/lionsoul2014/ip2region) # 一个离线IP地址定位库和IP定位数据管理框架，10微秒级别的查询效率，提供了众多主流编程语言的 `xdb` 数据生成和查询客户端实现。

https://github.com/Loyalsoldier/geoip # GeoIP 规则文件加强版，支持自行定制 V2Ray dat 格式文件 geoip.dat、MaxMind mmdb 格式文件、sing-box SRS 格式文件、mihomo MRS 格式文件、Clash ruleset、Surge ruleset 等。
