---
title: TCPDump
weight: 1
---

# 概述

> 参考：
>
> - [官网](https://www.tcpdump.org/)
> - [Manual(手册)，tcpdump(1)](https://www.tcpdump.org/manpages/tcpdump.1.html)
> - [Manual(手册)，PCAP-FILTER](https://www.tcpdump.org/manpages/pcap-filter.7.html)，主要描述过滤表达式的语法
> - [Wiki，tcpdump](https://en.wikipedia.org/wiki/Tcpdump)
> - <https://www.middlewareinventory.com/blog/tcpdump-capture-http-get-post-requests-apache-weblogic-websphere/>

TCPDump 是一个在命令行界面下的 [Packet analyzer](/docs/7.信息安全/Packet%20analyzer/Packet%20analyzer.md)(数据包分析器)。tcpdump 适用于大多数类 Unix 操作系统，在这些系统中，tcpdump 使用 [libpcap](https://en.wikipedia.org/wiki/Libpcap) 库来捕获数据包。对于 Windows 操作系统来说，tcpdump 使用的 pcap API 是 [WinPcap](https://en.wikipedia.org/wiki/WinPcap)(即 libpcap 的 Windows 版本)。

tcpdump 最初由[Van Jacobson](https://en.wikipedia.org/wiki/Van_Jacobson)、[Sally Floyd](https://en.wikipedia.org/wiki/Sally_Floyd)、[Vern Paxson](https://en.wikipedia.org/wiki/Vern_Paxson)和[Steven McCanne](https://en.wikipedia.org/w/index.php?title=Steven_McCanne&action=edit&redlink=1)于 1988 年编写，他们当时在[劳伦斯伯克利实验室](https://en.wikipedia.org/wiki/Lawrence_Berkeley_Laboratory)网络研究小组工作。到 1990 年代后期，有许多版本的 tcpdump 作为各种操作系统的一部分分发，以及许多没有很好协调的补丁。 [Michael Richardson (mcr)](<https://en.wikipedia.org/w/index.php?title=Michael_Richardson_(mcr)&action=edit&redlink=1>)和[Bill Fenner](https://en.wikipedia.org/w/index.php?title=Bill_Fenner&action=edit&redlink=1)于 1999 年创建 www.tcpdump.org

说明：Dump 有 **转出，倾卸；转储；内容全部打印** 的含义，在官方文档中，通过 TCPDump 程序输出的数据包，通常称为 dump line，转储的行。说白了就是程序抓到的包，每个包都是一行~~~~~

# Syntax(语法)

**tcpdump \[OPTIONS] \[Filter-Expression]**

## OPTIONS

- **-A** # 以 ASCII 打印每个数据包 (减去其链接级别标头)。方便捕捉网页。
- **-c \<INT>** # 指定程序将会捕获的数据包数量。
- **-D, --list-interfaces** # 列出可用于抓包的接口。将会列出接口的数值编号和接口名，它们都可以用于"-i"后
- **-e** # 在每条 dump 出来的行上显示二层头信息。这个选项可以输出 以太网 和 IEEE802.11 等协议的 MAC 层信息。
  - 通常用来抓取 VLAN 的 Tag。
- **-F \<FEIL>** # 指定一个包含 Filter-Expression 语法的文件。程序将会使用该文件的内容作为过滤表达式，并忽略命令行给出的过滤表达式。
- **-i, --interface \<DEV>** # 抓取指定网卡 DEV 的包，`默认值：any`，即抓取所有设备
  - 注意：在有 Bond 的服务器上，不要抓所有设备的包，否则使用 Wireshark 读取抓包文件时，会显示出很多乱序和重传，这是因为 Bond 设备和 Bond Salve 设备的包是相同的，但是咱都抓了。相同的包，时间不同，Wireshark 就识别成乱序了。
- -l # 使用标准输出列的缓冲区；
- **-n** # 不把主机的网络地址转换成 IP。可以多次指定，-nn 表示不转换 IP 地址和端口号的名称
- -O # 不将数据包编码最佳化
- -p # 不让网络界面进入混杂模式
- -q # 快速输出，仅列出少数的传输协议信息
- **-r \<FILE>** # 从 FILE 读取数据包。FILE 是通过 -w 选项保存的文件，或者任何使用 pcap API 的应用程序生成的文件。
- **-s <数据包大小>** # 设置每个数据包的大小；
- **-S, --absolute-tcp-sequence-numbers** # 输出 TCP sequence 号的绝对值，而不是相对值
- **-tt** # 显示每个抓到的包的时间戳(自 1970 年 1 月 1 日 00:00:00 以来)
- **-tttt** # 显示每个抓到的包的绝对时间
- **-ttttt** # 显示每个抓到的包的相对时间，单位微妙，第一个包是 `00:00:00.000000`，
- **-T <数据包类型>** # 强制将表达方式所指定的数据包转译成设置的数据包类型；
- **-v\[vv]** # 从一个 v 开始，每多一个 v 则抓出的包的信息则做出一部分最多 3 个 v，包信息最多
- **-w \</Path/TO/FILE>** # 把数据包数据写入指定的文件

## Filter-Expression(过滤表达式)

**Filter Expression(过滤表达式)** 用于过滤抓取到的数据包。可以让 tcpdump 程序只输出通过表达式匹配到的数据包。

过滤表达式由一个或多个 **Primitives(原语)** 组成。原语通常是一个具有一个或多个 **Qualifiers(限定词)** 的 ID。ID 就是 端口号、主机名、IP 地址 等等

### Qualifiers(限定词)

一共有三种类型的限定词：

- **TYPE** # 类型限定词。指定 ID 的类型。默认值：`host`
  - **host** # 匹配主机名称或 IP 地址
  - **net** # 匹配网段，CIDR 模式。比如：1.2.3.0/24
  - **port** # 匹配端口号
  - **portrange** # 匹配端口号的范围。两个端口号中间以 `-` 连接
- **DIR** # direction(方向) 限定词。指定 ID 的传输方向。默认值：`src or dst`
  - **src**
  - **dst**
  - **...... 等等**
- **PROTO** # 协议限定词。指定要匹配的协议
  - **icmp**
  - **ip**
  - **ip6**
  - **arp**
  - **tcp**
  - **udp**
  - **......等等**

**不同类型的限定符之间可以互相组合**，以实现更复杂的描述，比如：

- 类型匹配和方向匹配组合
  - src host 1.1.1.1 # 匹配源主机为 1.1.1.1 的数据包
- 类型匹配、方向匹配、协议匹配全部组合在一起
  - udp dst portrange 100-400 # 匹配 目的端口范围在 100-400 之间的 UDP 协议的数据包

**相同类型的限定符之间也可以互相组合**，但是需要使用逻辑表达式 **and、or、not**(或者用符号表示 `&&`、`||`、`!`) 来连接两个限定符，比如：

- icmp or tcp # 匹配 icmp 或 tcp 的数据包
- ip host 1.1.1.1 and ! 2.2.2.2 # 匹配 1.1.1.1 主机，但不包括 2.2.2.2 主机之间的通信
- src and dst port 1111 # 匹配 源端口 和 目的端口 都是 1111 的数据包
- tcp src or dst portrange 1111-2222 # 匹配 源 或者 目的 任意一个方向的端口范围在 1111-2222 之间的 TCP 协议的数据包
- 注意：为了节省输入，可以**省略相同的限定符**，下面两个表达式的作用完全相同：
  - tcp dst port 21 or 22 or 23
  - tcp dst port 21 or tcp dst port 22 or tcp dst port 33
- icmp or tcp dst port 23 and host 1.1.1.1 # 匹配 icmp 或 tcp，目的端口为 23 且主机为 1.1.1.1 的数据包

对于比较复杂的过滤器表达式，为了逻辑的清晰，可以使用括号。不过默认情况下，tcpdump 把 () 当做特殊的字符，所以必须使用单引号 ' 来消除歧义：

- tcpdump -nvv -c 20 'src 10.0.2.4 and (dat port 3389 or 22)'

#### TYPE 限定词

#### DIR 限定词

#### PROTO 限定词

各种协议的数据包，都有其各自的包头，tcpdump 可以根据包头中的字段的值进行过滤，比如只抓 TCP 标志为 SYN 的包，只抓 HTTP 的包，等等。

协议限定词扩展语法如下：

- **PROTO\[EXPR:SIZE]**
  - PROTO # 协议名
  - EXPR # 表达式
  - SIZE #

简单示例：

- 匹配 TCP 头中标志为 tcp-syn 的包
  - **'tcp\[tcpflags] & tcp-syn != 0'**
- 匹配 TCP 头中标志为 SYN 或 FIN 的包。也就是说一个 TCP 会话的开始和结束的包。
  - **tcp\[tcpflags] & (tcp-syn|tcp-fin) != 0**

# 应用示例

从 bond0 网卡抓取 源地址为 10.10.10.10 且 端口号为 18999 的 UDP 的包，每个包大小为 500

- tcpdump -i bond0 udp port 18999 and src host 10.10.10.10 -s 500 0nvvvSe -T snmp

抓取 eth9 网卡，到 111.30.199.159 的包，记录到/tmp/srcIP 中

- tcpdump -i eth9 host 111.30.199.159 -Avns 0 -w /tmp/srcIP
- tcpdump -i eth0 -tnn dst port 80 -c 1000 | awk -F”.” '{print $1"."$2"."$3"."$4"."}' | sort |uniq -c | sort -nr | head-5

获取 vnet1 网卡的 icmp 的包，i.e.抓 ping 命令的包

- tcpdump icmp -i vnet1

在 eth0 上抓取目的端口是 53 的 udp 包

- tcpdump -i eth0 udp dst port 53

抓取多个 host 的数据包

- tcpdump -i any host 172.19.42.202 or 172.19.42.201 -nn

从 192.168 网段到 10 或者 172.16 网段的数据报

- tcpdump -nvX src net 192.168.0.0/16 and dat net 10.0.0.0/8 or 172.16.0.0/16

抓取 LACP 包，是否有用待验证？

- `tcpdump -i bond1 ether host 01:80:c2:00:00:02`
- 注意：ether 只能在某些设备上运行，需要指定具体网络设备，不能用 any

使用 tcpdump 截取数据报文的时候，默认会打印到屏幕的默认输出，你会看到按照顺序和格式，很多的数据一行行快速闪过，根本来不及看清楚所有的内容。不过，tcpdump 提供了把截取的数据保存到文件的功能，以便后面使用其他图形工具（比如 wireshark，Snort）来分析。

-r 可以读取文件里的数据报文，显示到屏幕上。

tcpdump -nXr capture_file.pcap host web30

NOTE：保存到文件的数据不是屏幕上看到的文件信息，而是包含了额外信息的固定格式 pcap，需要特殊的软件来查看，使用 vim 或者 cat 命令会出现乱码。

从 Mars 或者 Pluto 发出的数据报，并且目的端口不是 22

tcpdump -vv src mars or pluto and not dat port 22

## 高级过滤

匹配 TCP 标志为 RST 的所有包。

- 意思就是数据包的包头中，tcpflags 字段的 tcp-rst 的值为 1
- `tcpdump 'tcp[tcpflags] & (tcp-rst) != 0' -nn`

## 抓 HTTP 的包

**抓取 HTTP 的包**

- `tcpdump -AnnSs 70 'tcp[20:4]=0x48545450'`
  - `tcp[20:4]` # 由 [TCP Segment 结构](/docs/4.数据通信/通信协议/TCP_IP/TCP/TCP.md#TCP%20Segment%20结构) 可知，TCP 的首部长度为 20 字节，那么 TCP 首部后面的 4 字节通常用来标识
- 或者
- `tcpdump -AnnSs 70 'tcp[((tcp[12:1] & 0xf0) >> 2):4] = 0x48545450'`
  - `tcp[((tcp[12:1] & 0xf0) >> 2):4]` # 首先确定我们感兴趣的字节的位置 (在 TCP 头之后)，然后选择我们希望匹配的 4 个字节。TODO: 还没搞明白这个描述
  - `0x47455420` # 0x 表示 16 进制，47455420 是 ASCII 中以 16 进制表示的 `GET` 这四个字符(G、E、T、空格)。

> Notes: 其中 -s 选项的值设为 70 是为了避免输出过多数据信息干扰。70 长度一般足够看到请求的 PATH，如果 PATH 太长，可以适当提高 -s 的值。

**抓取 HTTP 的 GET 请求的包**

- `tcpdump -AnnSs 0 tcp[20:2]=0x4745 or tcp[20:2]=0x4854`
  - 这里是用了 TCP 首部后两位进行匹配，即 GE 和 HT
- 或者
- `tcpdump -AnnSs 0 'tcp[((tcp[12:1] & 0xf0) >> 2):4] = 0x47455420'`

**抓取 HTTP 的 GET 或 POST 请求的包**

- `tcpdump -nnAs 0 'tcp[((tcp[12:1] & 0xf0) >> 2):4] = 0x47455420 or tcp[((tcp[12:1] & 0xf0) >> 2):4] = 0x504F5354'`

为什么要用 -Xe 选项？

- `tcpdump -XvvennSs 0 tcp[20:2]=0x4745 or tcp[20:2]=0x4854`

# tcpdump 输出内容详解

截取数据只是第一步，第二步就是理解这些数据，下面就解释一下 tcpdump 命令输出各部分的意义。

```bash
21:27:06.995846 IP (tos 0x0, ttl 64, id 45646, offset 0, flags [DF], proto TCP (6), length 64)
    192.168.1.106.56166 > 124.192.132.54.80: Flags [S], cksum 0xa730 (correct), seq 992042666, win 65535, options [mss 1460,nop,wscale 4,nop,nop,TS val 663433143 ecr 0,sackOK,eol], length 0
21:27:07.030487 IP (tos 0x0, ttl 51, id 0, offset 0, flags [DF], proto TCP (6), length 44)
    124.192.132.54.80 > 192.168.1.106.56166: Flags [S.], cksum 0xedc0 (correct), seq 2147006684, ack 992042667, win 14600, options [mss 1440], length 0
21:27:07.030527 IP (tos 0x0, ttl 64, id 59119, offset 0, flags [DF], proto TCP (6), length 40)
    192.168.1.106.56166 > 124.192.132.54.80: Flags [.], cksum 0x3e72 (correct), ack 2147006685, win 65535, length 0
```

最基本也是最重要的信息就是数据报的源地址/端口和目的地址/端口，上面的例子第一条数据报中，源地址 ip 是 192.168.1.106，源端口是 56166，目的地址是 124.192.132.54，目的端口是 80。 > 符号代表数据的方向。

此外，上面的三条数据还是 tcp 协议的三次握手过程，第一条就是 SYN 报文，这个可以通过 Flags \[S] 看出。下面是常见的 TCP 报文的 Flags:

- **\[S]** # SYN（开始连接）
- **\[.]** # 没有 Flag
- **\[P]** # PSH（推送数据）
- **\[F]** # FIN （结束连接）
- **\[R]** # RST（重置连接）

而第二条数据的 \[S.] 表示 SYN-ACK，就是 SYN 报文的应答报文。

tcpdump 很详细的 <http://blog.chinaunix.net/uid-11242066-id-4084382.html>

<http://www.cnblogs.com/ggjucheng/archive/2012/01/14/2322659.html> Linux tcpdump 命令详解

Tcpdump usage examples（推荐） <http://www.rationallyparanoid.com/articles/tcpdump.html>

使用 TCPDUMP 抓取 HTTP 状态头信息 <http://blog.sina.com.cn/s/blog_7475811f0101f6j5.html>
