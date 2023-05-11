---
title: IP
---

# 概述

> 参考：
>
> - [RFC,791](https://datatracker.ietf.org/doc/html/rfc791)(IP 规范)
> - [Wiki,Internet Protocol](https://en.wikipedia.org/wiki/Internet_Protocol)
> - [Wiki,IPv4](https://en.wikipedia.org/wiki/IPv4)
> - [Wiki,Mask(掩码)](<https://en.wikipedia.org/wiki/Mask_(computing)>)
> - [Wiki,Classful Network(分类网络)](https://en.wikipedia.org/wiki/Classful_network)
> - [IANA,IPv4 地址空间分配情况](https://www.iana.org/assignments/ipv4-address-space/ipv4-address-space.xhtml)
>   - [APNIC](https://www.apnic.net/)(管理亚太地区的 IP 地址注册机构)
>     - [APNIC,帮助-FTP 数据库](https://ftp.apnic.net/stats/apnic/)(亚太地区所有分配的 IP 地址信息)
> - [IANA,IPv4 特殊用途地址注册表](https://www.iana.org/assignments/iana-ipv4-special-registry/iana-ipv4-special-registry.xhtml)

**Internet Protocol(互联网协议，简称 IP)** 是[互联网协议套件](https://en.wikipedia.org/wiki/Internet_protocol_suite)(其中包含 TCP/IP)中的主要通信协议，用于跨网络边界中继数据报。它的路由功能可实现互联网络，并实质上建立了 Internet。

> **Internet protocol suite(互联网协议套件)** 是互联网和类似计算机网络中使用的概念模型和通信协议集。由于该套件中的基本协议是 **TCP(传输控制协议)** 和 **IP(互联网协议)**，因此通常被称为 **TCP/IP**。在其开发过程中，其版本被称为国防部（DoD）模型，因为联网方法的开发是由美国国防部通过 DARPA 资助的。它的实现是一个协议栈。

IP 基于数据包的 Header 中的 IP 地址，将数据包从源主机发送到目标主机。基于此目的，IP 还定义了数据包的封装结构、以及一种寻址方法。寻址方法用来使用源和目标的信息标记数据报。

从历史上看，IP 是在 1974 年由 Vint Cerf 和 Bob Kahn 引入的原始 **Transmission Control Program(传输控制程序)** 中的[无连接](https://en.wikipedia.org/wiki/Connectionless_communication)数据报服务。该服务由一项面向连接的服务补充，成为 [**Transmission Control Protocol(传输控制协议，简称 TCP)**](/docs/4.数据通信/通信协议/3_4.TCP_IP/TCP/TCP.md) 的基础。因此 IP 套件通常称为 TCP/IP。IP 的第一个版本是 IPv4，继任者是 IPv6

# IPv4 地址

IPv4 地址最多使用 32 bit 表示，即最多 32 个 1，这 32 bit 以 `点` 分割为 4 组，每组 8 bit，在使用时，使用十进制表示。比如：`192.168.0.1`。
![](https://notes-learning.oss-cn-beijing.aliyuncs.com/nahgxm/1633534934848-ca44c51a-c787-47e7-a82b-589b6f78124b.jpeg)

## IPv4 地址结构

IPv4 地址的这 32 bit 可以分为两部分

- 网络位 # n bit
- 主机位 # 32 - n bit

这两个部分通过 **Subnet Mast(子网掩码)** 来区分，子网掩码由一连串的 1 和 0 组成，遵从以下规则：

- 1 对应网络位
- 0 对应主机位
- 1 和 0 不能交叉出现

将子网掩码和 IP 地址作“与”操作后，IP 地址的主机部分将被丢弃，剩余的是网络地址和子网地址。

例如：一个 IP 地址为 10.2.45.1，子网掩码为 255.255.252.0，“与” 运算得到：10.2.44.0，则网络设备认为该 IP 地址的网络号与子网号为 10.2.44.0，属于 10.2.44.0/22 网络，其中/22 表示子网掩码长度为 22 位，即从前向后连续的 22 个 1。

00001010.00000010.00101101.00000001
与运算
11111111.11111111.11111100.00000000
结果为
00001010.00000010.00101100.00000001 即 10.2.44.0

## IPv4 地址分类

- **单播地址**

| 类 | 开头的 bit | 网络位 bit 数 | 主机位 bit 数 | 子网数量 | 每个子网的地址数 | 总地址数 | 起始地址 | 结束地址 | 默认子网掩码 | [CIDR](https://en.wikipedia.org/wiki/CIDR_notation) |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |
| Class A | 0 | 8 | 24 | 128 (27) | 16,777,216 (224) | 2,147,483,648 (231) | 0.0.0.0 | 127.255.255.255[\[a\]](https://en.wikipedia.org/wiki/Classful_network#cite_note-5) | 255.0.0.0 | /8 |
| Class B | 10 | 16 | 16 | 16,384 (214) | 65,536 (216) | 1,073,741,824 (230) | 128.0.0.0 | 191.255.255.255 | 255.255.0.0 | /16 |
| Class C | 110 | 24 | 8 | 2,097,152 (221) | 256 (28) | 536,870,912 (229) | 192.0.0.0 | 223.255.255.255 | 255.255.255.0 | /24 |

- **组播地址**

  - **D 类 224-239 224.0.0.0 ~ 239.255.255.255**

- **保留地址**

  - **E 类 240 - 254 240.0.0.0 ~ 255.255.255.255**

- **特殊地址**

  - **网络地址** # 网络位不变，主机位全为 0 的 IP 地址代表网络本身
  - **Broadcast Address(广播地址)** # 网络位不变，主机位全为 1 的 IP 地址代表本网络的广播。是专门用于同时向网络中所有[工作站](https://baike.baidu.com/item/%E5%B7%A5%E4%BD%9C%E7%AB%99/217955)进行发送的一个**地址**。在使用[TCP/IP 协议](https://baike.baidu.com/item/TCP%2FIP%20%E5%8D%8F%E8%AE%AE/2116790)的网络中，[主机](https://baike.baidu.com/item/%E4%B8%BB%E6%9C%BA/455151)[标识](https://baike.baidu.com/item/%E6%A0%87%E8%AF%86/6396929)段 host ID 为全 1 的 IP 地址为广播地址，广播的分组传送给 host ID 段所涉及的所有[计算机](https://baike.baidu.com/item/%E8%AE%A1%E7%AE%97%E6%9C%BA/140338)。例如，对于 10.1.1.0 （255.0.0.0 ）网段，其直播[广播](https://baike.baidu.com/item/%E5%B9%BF%E6%92%AD/656406)地址为 10.255.255.255 （255 即为 2 进制的 11111111 ），当发出一个目的地址为 10.255.255.255 的分组（[封包](https://baike.baidu.com/item/%E5%B0%81%E5%8C%85/2017669)）时，它将被分发给该[网段](https://baike.baidu.com/item/%E7%BD%91%E6%AE%B5/11026985)上的所有计算机。
  - **Link Local(链路本地地址)** # 169.254.0.0 ~ 169.254.255.255。用于[链路本地地址](https://en.wikipedia.org/wiki/Link-local_address)[\[9\]](https://en.wikipedia.org/wiki/IPv4#cite_note-rfc3927-9)两台主机之间的单个链路上时，否则指定 IP 地址，如将有通常被从检索到的[DHCP](https://en.wikipedia.org/wiki/DHCP)服务器。

- **Private Network(私人网络地址)**

| 名称 | [CIDR](https://en.wikipedia.org/wiki/Classless_Inter-Domain_Routing) | 地址范围 | 地址数量 | 描述 |
| --- | --- | --- | --- | --- |
| 24-bit block | 10.0.0.0/8 | 10.0.0.0 – 10.255.255.255 | 16777216 | 一个完整的 A 类地址 Single Class A. |
| 20-bit block | 172.16.0.0/12 | 172.16.0.0 – 172.31.255.255 | 1048576 | Contiguous range of 16 Class B blocks. |
| 16-bit block | 192.168.0.0/16 | 192.168.0.0 – 192.168.255.255 | 65536 | Contiguous range of 256 Class C blocks. |

# IPv4 Datagram 结构

IPv4 数据报被封装在链路层的 Frame 中

IPv4 数据报首部共 14 个字段，其中 13 个是必须的，第 14 个是可选的。前 13 个字段长度固定为 20 Bytes，即 160 bit；第 14 个字段长度在 0 ~ 40 Bytes 之间。

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/nahgxm/1633533991076-2a9cb606-463a-4dd4-89c8-c3aae590c113.jpeg)

对照 WireShark 中展示的内容看，排除 `[]` 中的内容，每一行就是首部中的一个字段

![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/nahgxm/1633532862295-9b420f37-7a97-43b9-85c8-1e973ea1aa59.png)

- **Version(版本)** # IP 协议的版本号。IPv4 其版本号为 4，因此在这个字段上的值为“6”。
- **Internet Header Length(首部长度，简称 IHL)** # 由于 Options 字段的长度是可变的。所以 IPv4 的首部长度也是可变的。该字段的值在 5 ~ 15 之间(该字段只有 4 bits，1111 即为 15)
  - 首部长度的计算方式如下：`IHL * 32 bits`。
    - 若 IHL 的值为 5，也就是说 Options 字段为 0，那么 IPv4 首部长度就是 5 \* 32 bits = 160 bits = 20 Bytes
  - 就像上面的 IPv4 的 Datagram 结构图一样，每行都是 32 bit，不算 Options 字段和 Payload，那么刚好是 5 行。
- **Differentiated Services Field** # 差异化的服务字段，基本没啥用。。。。o(╯□╰)o
  - **Differentiated Services Code Point** # 最初定义为 Type Of Service(服务类型，简称 TOS)，
  - **Explicit Congestion Notification** # 该字段定义在 RFC3168 中，
- **Total Length** # 定义了整个 IP 数据报的大小，最小为 20 字节(Payload 字段无内容)，最大为 65535 字节。
- **Identification**# 主要用于唯一标识单个 IP 数据报的片段组。
  - 一些实验工作建议将 ID 字段用于其他目的，例如添加数据包跟踪信息以帮助跟踪具有欺骗源地址的数据报，\[31] 但 RFC 6864 现在禁止任何此类使用。
- **Flags**# 用来控制或识别 IP 分片之后的每个片段，这 3 个 bit 分别表示不同的含义，若字段值为 0 表示未设置，值为 1 表示设置，类似 TCP 首部中 Flags 字段的用法。
  - 第一个 # Reserved，保留字段，必须为 0
  - 第二个 # Don't Fragment(DF)
  - 第三个 # More Fragment(MF)
- **Fragment Offset(分片偏移)** # IP 分片之后的偏移量
- **Time To Live(存活时间，简称 TTL)** # 其实用 Hop Limit 的描述更准确，封包每经过一个路由器，怎会将 TTL 字段的值减 1，减到 0 是，该包将会被丢弃。
- **Protocol**# 封装 IP 数据报的上层协议，比如 6 表示 TCP、1 表示 ICMP
  - 每种协议根据 [RFC 1700](https://datatracker.ietf.org/doc/html/rfc1700) 都分配了一个固定的编号，该 RFC 1700 最终被 RFC 3232 废弃，并将协议编号的维护工作，转到[IANA 的在线数据库](https://www.iana.org/assignments/protocol-numbers/protocol-numbers.xhtml)中
- **Header Checksum** # 当数据包到达路由器时，路由器会计算标头的校验和，并将其与校验和字段进行比较。如果值不匹配，则路由器会丢弃该数据包。
- **Source Address(源地址)** # 发送端 IP 地址。
- **Destination Address(目标地址)** # 接收端 IP 地址。
- **Options(选项)**# 可变长度，0-40 Bytes。

# IPv4 Fragment

IP Fragment(分片) 主要通过首部中的 Identification、Flags、Fragment Offset 这三个字段对每一个分片进行唯一标识

# IP 地址分配机制

> 参考：
>
> - [IANA,号码资源](https://www.iana.org/numbers)
> - [面包板，你知道中国大陆一共有多少 IPv4 地址吗？](https://www.eet-china.com/mp/a54338.html)
> - [公众号，k8s 中文社区-居然还有 2 亿多 IPv4 地址未分配](https://mp.weixin.qq.com/s/GHYYgZwAuEV4qPCwdI8Bjg)
> - [APNIC,搜索](https://wq.apnic.net/static/search.html)(通过给定的 IP 地址搜索谁拥有这个 IP)

IPv4 和 IPv6 地址通常以分层方式分配。**ISP(互联网服务提供商)** 为用户分配 IP 地址。ISP 从 **LIR(本地互联网注册机构)** 或 **NIR(国家互联网注册机构)** 或 **RIR(相应的区域互联网注册机构)** 获取 IP 地址分配
![](https://notes-learning.oss-cn-beijing.aliyuncs.com/nahgxm/1646384624162-21c9bca1-0960-45e4-87bb-3802eca96278.svg)

| 登记处                             | 覆盖面积                                                              |
| ---------------------------------- | --------------------------------------------------------------------- |
| [AFRINIC](http://www.afrinic.net/) | Africa Region(非洲地区)                                               |
| [APNIC](http://www.apnic.net/)     | Asia/Pacific Region(亚洲/太平洋地区，亚太地区)                        |
| [ARIN](http://www.arin.net/)       | Canada, USA, and some Caribbean Islands(加拿大、美国、一些加勒比岛屿) |
| [LACNIC](http://www.lacnic.net/)   | Latin America and some Caribbean Islands(拉丁美洲、一些加勒比岛屿)    |
| [RIPE NCC](http://www.ripe.net/)   | Europe, the Middle East, and Central Asia(欧洲、中东、中亚)           |

[对 IP 地址的主要作用是根据全球政策](http://www.icann.org/en/general/global-addressing-policies.html)所述的需求将未分配地址池分配给 RIR，并记录 [IETF](http://www.ietf.org/) 所做的协议分配。当 RIR 需要在其区域内分配或分配更多 IP 地址时，我们会向 RIR 进行额外分配。我们不会直接向 ISP 或最终用户进行分配，除非在特定情况下，例如分配多播地址或其他协议特定需求。

APNIC 是全球 5 个地区级的 Internet 注册机构（RIR）之一，负责亚太地区的以下事务：

1. 分配 IPv4 和 IPv6 地址空间，AS 号；
2. 为亚太地区维护 Whois 数据库；
3. 反向 DNS 指派；
4. 在全球范围内作为亚太地区的 Internet 社区的代表。

所以，中国大陆境内的地址都会登记在 APNIC 的地址库内。地址库获取方式：<http://ftp.apnic.net/apnic/stats/apnic/delegated-apnic-latest>
例如在 Linux 系统中，使用 wget 命令可以拉取文件。

![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/nahgxm/1646295854669-69e90c43-6765-415a-be08-ba15cdf1f6c3.png)

文件内容条目参考如下：

apnic|JP|asn|173|1|20020801|allocated
apnic|ID|ipv4|43.240.228.0|1024|20140818|allocated
apnic|HK|ipv6|2001:df5:b800::|48|20140514|assigned

条目格式如下：

注册机构|国家代码|类型|起始位|长度|分配日期|状态

- **注册机构**：亚太地区一般为 apnic
- **国家代码**：ISO-3166 定义的两位国家或地区代码，如中国为 CN
- **类型**：asn（Autonomous System Number，自治系统编号），也就是 BGP 的 AS 编号；ipv4，IPv4 地址；ipv6，IPv6 地址
- **起始位**：第一个 ASN 编号或 IP 地址
- **长度**：从第一个起始位开始，申请分配多少的编号或地址
- **分配日期**：国家或地区向 APNIC 申请的日期
- **状态**：allocated 和 assigned，都是已分配

所以，需要将 delegated-apnic-latest 文件中所有国家为 CN、且类型为 ipv4 的条目导出，并转换为静态路由格式。

例如使用命令将符合条件的条目导入到 china 文件中。

![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/nahgxm/1646295854726-90496001-56ba-4bbb-9e7c-48568a601999.png)

可以查看文件行数，代表有多少条明细条目。

![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/nahgxm/1646295854636-89ef24bb-cb5c-4832-a3af-bcc97cffc042.png)

然后根据起始位和长度，转换出静态路由所需的目的地址和掩码即可。在 excel 中通过对长度进行函数运算，可以得到掩码长度，如：=32-LOG(E2,2)，代入 2048 的话，可得到掩码长度为 21。操作后得到类似下图的表格：

![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/nahgxm/1646295854710-c698ce59-e5fb-4e59-97cc-77a1a81bfc81.png)

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

![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/nahgxm/1646295854663-30f55fe1-a908-42e8-8b91-e95b33552417.png)

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
