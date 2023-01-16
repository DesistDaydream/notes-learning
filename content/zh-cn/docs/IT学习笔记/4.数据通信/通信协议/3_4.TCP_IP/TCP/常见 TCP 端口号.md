---
title: 常见 TCP 端口号
---

# Socket 网络端口含义

> 参考：
> - [Wiki,TCP 与 UDP 端口号列表](https://en.wikipedia.org/wiki/List_of_TCP_and_UDP_port_numbers)
> - [IANA,分配-服务名和端口号注册列表](https://www.iana.org/assignments/service-names-port-numbers/service-names-port-numbers.xhtml)

分类情况：

- WellKnownPorts(公认端口)
  - 从 0 到 1023，它们紧密绑定（binding）于一些服务。通常这些端口的通讯明确表明了某种服务的协议。这些端口由 IANA 分配管理，IANA(The Internet Assigned Numbers Authority，互联网数字分配机构)是负责协调一些使 Internet 正常运作的机构
- RegisteredPorts(注册端口)
  - 从 1024 到 49151。是公司和其他用户向互联网名称与数字地址分配机构（ICANN）登记的端口号，利用因特网的传输控制协议（TCP）和用户数据报协议（UDP）进行通信的应用软件需要使用这些端口。在大多数情况下，这些应用软件和普通程序一样可以被非特权用户打开。
- Dynamicand/orPrivatePorts(动态和/或私有端口)
  - 从 49152 到 65535。这类端口号仅在客户进程运行时才动态选择，因此又叫做短暂端口号。被保留给客户端进程选择暂时使用的。也可以理解为，客户端启动的时候操作系统随机分配一个端口用来和服务器通信，客户端进程关闭下次打开时，又重新分配一个新的端口。

# TCP/UDP 端口列表

计算机之间依照互联网[传输层](https://zh.wikipedia.org/wiki/%E4%BC%A0%E8%BE%93%E5%B1%82)[TCP/IP 协议](https://zh.wikipedia.org/wiki/TCP/IP%E5%8D%8F%E8%AE%AE)的协议通信，不同的协议都对应不同的[端口](https://zh.wikipedia.org/wiki/%E9%80%9A%E8%A8%8A%E5%9F%A0)。并且，利用数据报文的[UDP](https://zh.wikipedia.org/wiki/%E7%94%A8%E6%88%B7%E6%95%B0%E6%8D%AE%E6%8A%A5%E5%8D%8F%E8%AE%AE)也不一定和[TCP](https://zh.wikipedia.org/wiki/%E4%BC%A0%E8%BE%93%E6%8E%A7%E5%88%B6%E5%8D%8F%E8%AE%AE)采用相同的端口号码。以下为两种通信协议的端口列表链接

## 0 到 1023 号端口

以下列表仅列出常用端口，详细的列表请参阅[IANA](https://zh.wikipedia.org/wiki/IANA)网站。[\[1\]](https://zh.wikipedia.org/wiki/TCP/UDP%E7%AB%AF%E5%8F%A3%E5%88%97%E8%A1%A8#cite_note-IANA-1)

| 端口                                   | 描述                                                                                                 | 状态   |
| -------------------------------------- | ---------------------------------------------------------------------------------------------------- | ------ |
| 0/TCP,UDP                              | 保留端口；不使用（若发送过程不准备接受回复消息，则可以作为源端口）                                   | 官方   |
| 1/TCP,UDP                              | [TCPMUX](https://zh.wikipedia.org/w/index.php?title=TCPMUX&action=edit&redlink=1)                    |
| （传输控制协议端口服务多路开关选择器） | 官方                                                                                                 |
| 5/TCP,UDP                              | [RJE](https://zh.wikipedia.org/w/index.php?title=Remote_Job_Entry&action=edit&redlink=1)             |
| （远程作业登录）                       | 官方                                                                                                 |
| 7/TCP,UDP                              | [Echo](<https://zh.wikipedia.org/wiki/Echo_(%E5%91%BD%E4%BB%A4)>)                                    |
| （回显）协议                           | 官方                                                                                                 |
| 9/UDP                                  | [DISCARD](https://zh.wikipedia.org/w/index.php?title=DISCARD&action=edit&redlink=1)                  |
| （丢弃）协议                           | 官方                                                                                                 |
| 9/TCP,UDP                              | [网络唤醒](https://zh.wikipedia.org/wiki/%E7%BD%91%E7%BB%9C%E5%94%A4%E9%86%92)                       | 非官方 |
| 11/TCP,UDP                             | [SYSTAT](https://zh.wikipedia.org/w/index.php?title=SYSTAT&action=edit&redlink=1)                    |
| 协议                                   | 官方                                                                                                 |
| 13/TCP,UDP                             | [DAYTIME 协议](https://zh.wikipedia.org/wiki/DAYTIME%E5%8D%8F%E8%AE%AE)                              | 官方   |
| 15/TCP,UDP                             | [NETSTAT](https://zh.wikipedia.org/w/index.php?title=NETSTAT&action=edit&redlink=1)                  |
| 协议                                   | 官方                                                                                                 |
| 17/TCP,UDP                             | [QOTD](https://zh.wikipedia.org/w/index.php?title=QOTD&action=edit&redlink=1)                        |
| （Quote of the Day，每日引用）协议     | 官方                                                                                                 |
| 18/TCP,UDP                             | 消息发送协议                                                                                         | 官方   |
| 19/TCP,UDP                             | [CHARGEN](https://zh.wikipedia.org/w/index.php?title=CHARGEN&action=edit&redlink=1)                  |
| （字符发生器）协议                     | 官方                                                                                                 |
| 20/TCP,UDP                             | [文件传输协议](https://zh.wikipedia.org/wiki/%E6%96%87%E4%BB%B6%E4%BC%A0%E8%BE%93%E5%8D%8F%E8%AE%AE) |
| &#x20;\- 默认数据端口                  | 官方                                                                                                 |
| 21/TCP,UDP                             | [文件传输协议](https://zh.wikipedia.org/wiki/%E6%96%87%E4%BB%B6%E4%BC%A0%E8%BE%93%E5%8D%8F%E8%AE%AE) |
| &#x20;\- 控制端口                      | 官方                                                                                                 |
| 22/TCP,UDP                             | [SSH](https://zh.wikipedia.org/wiki/Secure_Shell)                                                    |

（Secure Shell） - 安全远程登录协议，用于安全文件传输（[SCP](https://zh.wikipedia.org/wiki/%E5%AE%89%E5%85%A8%E5%A4%8D%E5%88%B6)
、[SFTP](https://zh.wikipedia.org/wiki/SSH%E6%96%87%E4%BB%B6%E4%BC%A0%E8%BE%93%E5%8D%8F%E8%AE%AE)
）及[端口转发](https://zh.wikipedia.org/wiki/%E7%AB%AF%E5%8F%A3%E8%BD%AC%E5%8F%91) | 官方 |
| 23/TCP,UDP | [Telnet](https://zh.wikipedia.org/wiki/Telnet)
终端仿真协议 - 未加密文本通信 | 官方 |
| 25/TCP,UDP | [SMTP（简单邮件传输协议）](https://zh.wikipedia.org/wiki/%E7%AE%80%E5%8D%95%E9%82%AE%E4%BB%B6%E4%BC%A0%E8%BE%93%E5%8D%8F%E8%AE%AE)
&#x20;\- 用于传递[电子邮件](https://zh.wikipedia.org/wiki/%E7%94%B5%E5%AD%90%E9%82%AE%E4%BB%B6) | 官方 |
| 26/TCP,UDP | [RSFTP](https://zh.wikipedia.org/w/index.php?title=RSFTP&action=edit&redlink=1)
&#x20;\- 一个简单的类似 FTP 的协议 | 非官方 |
| 35/TCP,UDP | 任意的私有打印机服务端口 | 非官方 |
| 37/TCP,UDP | [TIME 时间协议](https://zh.wikipedia.org/wiki/%E6%99%82%E9%96%93%E5%8D%94%E8%AD%B0) | 官方 |
| 39/TCP,UDP | [Resource Location Protocol（资源定位协议）](https://zh.wikipedia.org/w/index.php?title=%E8%B5%84%E6%BA%90%E5%AE%9A%E4%BD%8D%E5%8D%8F%E8%AE%AE&action=edit&redlink=1) | 官方 |
| 41/TCP,UDP | 图形 | 官方 |
| 42/TCP,UDP | [Host Name Server](https://zh.wikipedia.org/wiki/ARPA%E4%B8%BB%E6%9C%BA%E5%90%8D%E6%9C%8D%E5%8A%A1%E5%99%A8%E5%8D%8F%E8%AE%AE) | 官方 |
| 42/TCP,UDP | [WINS](https://zh.wikipedia.org/wiki/WINS)
（WINS 主机名服务） | 非官方 |
| 43/TCP | [WHOIS](https://zh.wikipedia.org/wiki/WHOIS)
协议 | 官方 |
| 49/TCP,UDP | [TACACS](https://zh.wikipedia.org/wiki/TACACS)
登录主机协议 | 官方 |
| 53/TCP,UDP | [DNS](https://zh.wikipedia.org/wiki/%E5%9F%9F%E5%90%8D%E6%9C%8D%E5%8A%A1%E5%99%A8)
（域名服务系统） | 官方 |
| 56/TCP,UDP | 远程访问协议 | 官方 |
| 57/TCP | MTP，邮件传输协议 | 官方 |
| 67/UDP | [BOOTP](https://zh.wikipedia.org/wiki/BOOTP)
（BootStrap 协议）服务；同时用于[动态主机设置协议](https://zh.wikipedia.org/wiki/%E5%8A%A8%E6%80%81%E4%B8%BB%E6%9C%BA%E8%AE%BE%E7%BD%AE%E5%8D%8F%E8%AE%AE) | 官方 |
| 68/UDP | [BOOTP](https://zh.wikipedia.org/wiki/BOOTP)
客户端；同时用于[动态主机设定协议](https://zh.wikipedia.org/w/index.php?title=%E5%8A%A8%E6%80%81%E4%B8%BB%E6%9C%BA%E8%AE%BE%E5%AE%9A%E5%8D%8F%E8%AE%AE&action=edit&redlink=1) | 官方 |
| 69/UDP | [小型文件传输协议](https://zh.wikipedia.org/wiki/%E5%B0%8F%E5%9E%8B%E6%96%87%E4%BB%B6%E4%BC%A0%E8%BE%93%E5%8D%8F%E8%AE%AE)
（小型文件传输协议） | 官方 |
| 70/TCP | [Gopher](https://zh.wikipedia.org/wiki/Gopher) | 官方 |
| 79/TCP | [手指](https://zh.wikipedia.org/wiki/%E6%89%8B%E6%8C%87)
协议 | 官方 |
| 80/TCP,UDP | [超文本传输协议](https://zh.wikipedia.org/wiki/%E8%B6%85%E6%96%87%E6%9C%AC%E4%BC%A0%E8%BE%93%E5%8D%8F%E8%AE%AE)
（超文本传输协议）或[快速 UDP 网络连接](https://zh.wikipedia.org/wiki/%E5%BF%AB%E9%80%9FUDP%E7%BD%91%E7%BB%9C%E8%BF%9E%E6%8E%A5)
\- 用于传输网页 | 官方 |
| 81/TCP | [XB Browser](https://zh.wikipedia.org/wiki/XB_Browser)
&#x20;\- [Tor](https://zh.wikipedia.org/wiki/Tor) | 非官方 |
| 82/UDP | [XB Browser](https://zh.wikipedia.org/wiki/XB_Browser)
&#x20;\- 控制端口 | 非官方 |
| 88/TCP | [Kerberos](https://zh.wikipedia.org/wiki/Kerberos)
&#x20;\- 认证代理 | 官方 |
| 101/TCP | 主机名 | 官方 |
| 102/TCP | ISO-TSAP 协议 | 官方 |
| 107/TCP | 远程[Telnet](https://zh.wikipedia.org/wiki/Telnet)
协议 | 官方 |
| 109/TCP | [邮局协议](https://zh.wikipedia.org/wiki/%E9%83%B5%E5%B1%80%E5%8D%94%E5%AE%9A)
（POP），第 2 版 | 官方 |
| 110/TCP | [邮局协议](https://zh.wikipedia.org/wiki/%E9%83%B5%E5%B1%80%E5%8D%94%E5%AE%9A)
，第 3 版 - 用于接收[电子邮件](https://zh.wikipedia.org/wiki/%E7%94%B5%E5%AD%90%E9%82%AE%E4%BB%B6) | 官方 |
| 111/TCP,UDP | Sun 远程过程调用协议 | 官方 |
| 113/TCP | [Ident](https://zh.wikipedia.org/w/index.php?title=Ident_protocol&action=edit&redlink=1)
&#x20;\- 旧的服务器身份识别系统，仍然被[IRC](https://zh.wikipedia.org/wiki/IRC)
服务器用来认证它的用户 | 官方 |
| 115/TCP | [简单文件传输协议](https://zh.wikipedia.org/wiki/%E7%AE%80%E5%8D%95%E6%96%87%E4%BB%B6%E4%BC%A0%E8%BE%93%E5%8D%8F%E8%AE%AE) | 官方 |
| 117/TCP | [UNIX 间复制协议](https://zh.wikipedia.org/wiki/UUCP)
（**U**nix to **U**nix **C**opy **P**rotocol，UUCP）的路径确定服务 | 官方 |
| 118/TCP,UDP | [SQL](https://zh.wikipedia.org/wiki/SQL)
服务 | 官方 |
| 119/TCP | [网络新闻传输协议](https://zh.wikipedia.org/wiki/%E7%BD%91%E7%BB%9C%E6%96%B0%E9%97%BB%E4%BC%A0%E8%BE%93%E5%8D%8F%E8%AE%AE)
&#x20;\- 用来收取新闻组的消息 | 官方 |
| 123/UDP | [NTP](https://zh.wikipedia.org/wiki/%E7%B6%B2%E7%B5%A1%E6%99%82%E9%96%93%E5%8D%94%E8%AD%B0)
（Network Time Protocol） - 用于时间同步 | 官方 |
| 135/TCP,UDP | [分布式运算环境](https://zh.wikipedia.org/wiki/%E5%88%86%E6%95%A3%E5%BC%8F%E9%81%8B%E7%AE%97%E7%92%B0%E5%A2%83)
（Distributed Computing Environment，DCE）终端解决方案及定位服务 | 官方 |
| 135/TCP,UDP | [微软](https://zh.wikipedia.org/wiki/%E5%BE%AE%E8%BD%AF)
终端映射器（End Point Mapper，EPMAP） | 官方 |
| 137/TCP,UDP | [NetBIOS](https://zh.wikipedia.org/wiki/NetBIOS)
&#x20;NetBIOS 名称服务 | 官方 |
| 138/TCP,UDP | [NetBIOS](https://zh.wikipedia.org/wiki/NetBIOS)
&#x20;NetBIOS 数据报文服务 | 官方 |
| 139/TCP,UDP | [NetBIOS](https://zh.wikipedia.org/wiki/NetBIOS)
&#x20;NetBIOS 会话服务 | 官方 |
| 143/TCP,UDP | [因特网信息访问协议](https://zh.wikipedia.org/wiki/%E5%9B%A0%E7%89%B9%E7%BD%91%E4%BF%A1%E6%81%AF%E8%AE%BF%E9%97%AE%E5%8D%8F%E8%AE%AE)
（IMAP） - 用于检索[电子邮件](https://zh.wikipedia.org/wiki/%E7%94%B5%E5%AD%90%E9%82%AE%E4%BB%B6) | 官方 |
| 152/TCP,UDP | BFTP, 后台文件传输程序 | 官方 |
| 153/TCP,UDP | 简单网关监控协议（[Simple Gateway Monitoring Protocol](https://zh.wikipedia.org/w/index.php?title=Simple_Gateway_Monitoring_Protocol&action=edit&redlink=1)
，SGMP） | 官方 |
| 156/TCP,UDP | SQL 服务 | 官方 |
| 158/TCP,UDP | DMSP, 分布式邮件服务协议 | 非官方 |
| 161/TCP,UDP | [简单网络管理协议](https://zh.wikipedia.org/wiki/%E7%AE%80%E5%8D%95%E7%BD%91%E7%BB%9C%E7%AE%A1%E7%90%86%E5%8D%8F%E8%AE%AE)
（SNMP) | 官方 |
| 162/TCP,UDP | SNMP 协议的 TRAP 操作 | 官方 |
| 170/TCP | 打印服务 | 官方 |
| 179/TCP | [边界网关协议](https://zh.wikipedia.org/wiki/%E8%BE%B9%E7%95%8C%E7%BD%91%E5%85%B3%E5%8D%8F%E8%AE%AE)
&#x20;(BGP) | 官方 |
| 194/TCP | [IRC](https://zh.wikipedia.org/wiki/IRC)
（互联网中继聊天） | 官方 |
| 201/TCP,UDP | AppleTalk 路由维护 | 官方 |
| 209/TCP,UDP | Quick Mail 传输协议 | 官方 |
| 213/TCP,UDP | [互联网分组交换协议](https://zh.wikipedia.org/wiki/%E4%BA%92%E8%81%94%E7%BD%91%E5%88%86%E7%BB%84%E4%BA%A4%E6%8D%A2%E5%8D%8F%E8%AE%AE)
（IPX） | 官方 |
| 218/TCP,UDP | MPP，信息发布协议 | 官方 |
| 220/TCP,UDP | [因特网信息访问协议](https://zh.wikipedia.org/wiki/%E5%9B%A0%E7%89%B9%E7%BD%91%E4%BF%A1%E6%81%AF%E8%AE%BF%E9%97%AE%E5%8D%8F%E8%AE%AE)
（IMAP），第 3 版 | 官方 |
| 259/TCP,UDP | ESRO, Efficient Short Remote Operations | 官方 |
| 264/TCP,UDP | [BGMP](https://zh.wikipedia.org/w/index.php?title=Border_Gateway_Multicast_Protocol&action=edit&redlink=1)
，边界网关多播协议 | 官方 |
| 308/TCP | Novastor 在线备份 | 官方 |
| 311/TCP | Apple 服务器管理员工具、工作组管理 | 官方 |
| 318/TCP,UDP | TSP，[时间戳协议](https://zh.wikipedia.org/w/index.php?title=%E6%97%B6%E9%97%B4%E6%88%B3%E5%8D%8F%E8%AE%AE&action=edit&redlink=1) | 官方 |
| 323/TCP,UDP | IMMP, Internet 消息映射协议 | 官方 |
| 383/TCP,UDP | HP OpenView HTTPs 代理程序 | 官方 |
| 366/TCP,UDP | ODMR，按需邮件传递 | 官方 |
| 369/TCP,UDP | Rpc2 端口映射 | 官方 |
| 371/TCP,UDP | ClearCase 负载平衡 | 官方 |
| 384/TCP,UDP | 一个远程网络服务器系统 | 官方 |
| 387/TCP,UDP | AURP，AppleTalk 升级用路由协议 | 官方 |
| 389/TCP,UDP | 轻型目录访问协议 LDAP | 官方 |
| 401/TCP,UDP | [不间断电源](https://zh.wikipedia.org/wiki/%E4%B8%8D%E9%97%B4%E6%96%AD%E7%94%B5%E6%BA%90)
（UPS） | 官方 |
| 411/TCP | [Direct Connect](https://zh.wikipedia.org/w/index.php?title=Direct_Connect_network&action=edit&redlink=1)
&#x20;Hub 端口 | 非官方 |
| 412/TCP | [Direct Connect](https://zh.wikipedia.org/w/index.php?title=Direct_Connect_network&action=edit&redlink=1)
&#x20;客户端—客户端 端口 | 非官方 |
| 427/TCP,UDP | [服务定位协议](https://zh.wikipedia.org/wiki/%E6%9C%8D%E5%8B%99%E5%AE%9A%E4%BD%8D%E5%8D%94%E5%AE%9A)
（SLP） | 官方 |
| 443/TCP | [超文本传输安全协议](https://zh.wikipedia.org/wiki/%E8%B6%85%E6%96%87%E6%9C%AC%E4%BC%A0%E8%BE%93%E5%AE%89%E5%85%A8%E5%8D%8F%E8%AE%AE)
或[QUIC](https://zh.wikipedia.org/wiki/QUIC) | 官方 |
| 444/TCP,UDP | [SNPP](https://zh.wikipedia.org/w/index.php?title=Simple_Network_Paging_Protocol&action=edit&redlink=1)
，简单网络分页协议 | 官方 |
| 445/TCP | Microsoft-DS ([Active Directory](https://zh.wikipedia.org/wiki/Active_Directory)
、Windows 共享、[震荡波蠕虫](https://zh.wikipedia.org/wiki/%E9%9C%87%E7%9B%AA%E6%B3%A2%E8%A0%95%E8%9F%B2)
、Agobot、Zobotworm) | 官方 |
| 445/UDP | Microsoft-DS [服务器消息块](https://zh.wikipedia.org/wiki/%E4%BC%BA%E6%9C%8D%E5%99%A8%E8%A8%8A%E6%81%AF%E5%8D%80%E5%A1%8A)
（SMB）文件共享 | 官方 |
| 464/TCP,UDP | Kerberos 更改/设定密码 | 官方 |
| 465/TCP | Cisco 专用协议 | 官方 |
| 465/TCP | [传输层安全性协议](https://zh.wikipedia.org/wiki/%E5%82%B3%E8%BC%B8%E5%B1%A4%E5%AE%89%E5%85%A8%E6%80%A7%E5%8D%94%E5%AE%9A)
加密的[简单邮件传输协议](https://zh.wikipedia.org/wiki/%E7%AE%80%E5%8D%95%E9%82%AE%E4%BB%B6%E4%BC%A0%E8%BE%93%E5%8D%8F%E8%AE%AE) | 官方 |
| 475/TCP | tcpnethaspsrv（[Hasp](https://zh.wikipedia.org/w/index.php?title=Hasp&action=edit&redlink=1)
&#x20;服务，TCP/IP 版本） | 官方 |
| 497/TCP | [dantz](https://zh.wikipedia.org/w/index.php?title=Retrospect&action=edit&redlink=1)
&#x20;备份服务 | 官方 |
| 500/TCP,UDP | [网络安全关系与密钥管理协议](https://zh.wikipedia.org/wiki/%E7%B6%B2%E8%B7%AF%E5%AE%89%E5%85%A8%E9%97%9C%E8%81%AF%E8%88%87%E9%87%91%E9%91%B0%E7%AE%A1%E7%90%86%E5%8D%94%E5%AE%9A)
，IKE-Internet 密钥交换 | 官方 |
| 502/TCP,UDP | [Modbus](https://zh.wikipedia.org/wiki/Modbus)
&#x20;协议 | 官方 |
| 512/TCP | exec，远程进程执行 | 官方 |
| 512/UDP | comsat 和 [biff](<https://zh.wikipedia.org/w/index.php?title=Biff_(%E9%9B%BB%E8%85%A6%E9%81%8B%E7%AE%97)&action=edit&redlink=1>)
&#x20;命令：用于电子邮件系统 | 官方 |
| 513/TCP | login，登录命令 | 官方 |
| 513/UDP | Who 命令，查看当前登录计算机的用户 | 官方 |
| 514/TCP | [远程外壳](https://zh.wikipedia.org/wiki/%E8%BF%9C%E7%A8%8B%E5%A4%96%E5%A3%B3)
&#x20;protocol - 用于在远程计算机上执行非交互式命令，并查看返回结果 | 官方 |
| 514/UDP | [Syslog](https://zh.wikipedia.org/wiki/Syslog)
&#x20;协议 - 用于系统登录 | 官方 |
| 515/TCP | [Line Printer Daemon protocol](https://zh.wikipedia.org/w/index.php?title=Line_Printer_Daemon_protocol&action=edit&redlink=1)
&#x20;\- 用于 LPD 打印机服务器 | 官方 |
| 517/UDP | Talk | 官方 |
| 518/UDP | NTalk | 官方 |
| 520/TCP | efs | 官方 |
| 520/UDP | Routing - [路由信息协议](https://zh.wikipedia.org/wiki/%E8%B7%AF%E7%94%B1%E4%BF%A1%E6%81%AF%E5%8D%8F%E8%AE%AE) | 官方 |
| 513/UDP | 路由器 | 官方 |
| 524/TCP,UDP | [NetWare 核心协议](https://zh.wikipedia.org/wiki/NetWare%E6%A0%B8%E5%BF%83%E5%8D%94%E5%AE%9A)
（NetWare 核心协议）用于一系列功能，例如访问 NetWare 主服务器资源、同步时间等 | 官方 |
| 525/UDP | Timed，时间服务 | 官方 |
| 530/TCP,UDP | [远程过程调用](https://zh.wikipedia.org/wiki/%E8%BF%9C%E7%A8%8B%E8%BF%87%E7%A8%8B%E8%B0%83%E7%94%A8) | 官方 |
| 531/TCP,UDP | AOL 即时通信软件, IRC | 非官方 |
| 532/TCP | netnews | 官方 |
| 533/UDP | netwall，用于紧急广播 | 官方 |
| 540/TCP | [UUCP](https://zh.wikipedia.org/wiki/UUCP)
（Unix-to-Unix 复制协议） | 官方 |
| 542/TCP,UDP | [商业](https://zh.wikipedia.org/wiki/%E5%95%86%E4%B8%9A)
（Commerce Applications） | 官方 |
| 543/TCP | klogin，Kerberos 登陆 | 官方 |
| 544/TCP | kshell，Kerberos 远程 shell | 官方 |
| 546/TCP,UDP | DHCPv6 客户端 | 官方 |
| 547/TCP,UDP | DHCPv6 服务器 | 官方 |
| 548/TCP | 通过传输控制协议（TCP）的 Appletalk 文件编制协议（AFP([苹果归档协议](https://zh.wikipedia.org/wiki/%E8%8B%B9%E6%9E%9C%E5%BD%92%E6%A1%A3%E5%8D%8F%E8%AE%AE)
)) | 官方 |
| 550/UDP | new-rwho，new-who | 官方 |
| 554/TCP,UDP | [即时流协议](https://zh.wikipedia.org/wiki/%E5%8D%B3%E6%99%82%E4%B8%B2%E6%B5%81%E5%8D%94%E5%AE%9A) | 官方 |
| 556/TCP | Brunhoff 的远程文件系统（RFS） | 官方 |
| 560/UDP | rmonitor, Remote Monitor | 官方 |
| 561/UDP | monitor | 官方 |
| 563/TCP,UDP | 通过[TLS](https://zh.wikipedia.org/wiki/TLS)
的[网络新闻传输协议](https://zh.wikipedia.org/wiki/%E7%BD%91%E7%BB%9C%E6%96%B0%E9%97%BB%E4%BC%A0%E8%BE%93%E5%8D%8F%E8%AE%AE)
（NNTPS） | 官方 |
| 587/TCP | 邮件消息提交（[简单邮件传输协议](https://zh.wikipedia.org/wiki/%E7%AE%80%E5%8D%95%E9%82%AE%E4%BB%B6%E4%BC%A0%E8%BE%93%E5%8D%8F%E8%AE%AE)
，[RFC 2476](https://tools.ietf.org/html/rfc2476)
） | 官方 |
| 591/TCP | [FileMaker](https://zh.wikipedia.org/wiki/FileMaker)
&#x20;6.0（及之后版本）网络共享（HTTP 的替代，见 80 端口） | 官方 |
| 593/TCP,UDP | HTTP RPC Ep Map（RPC over HTTP, often used by [Distributed COM](https://zh.wikipedia.org/wiki/Distributed_COM)
&#x20;services and [Microsoft Exchange Server](https://zh.wikipedia.org/wiki/Microsoft_Exchange_Server)
） | 官方 |
| 604/TCP | TUNNEL | 官方 |
| 631/TCP,UDP | [互联网打印协议](https://zh.wikipedia.org/wiki/%E4%BA%92%E8%81%AF%E7%B6%B2%E5%88%97%E5%8D%B0%E5%8D%94%E5%AE%9A) | |
| 636/TCP,UDP | [LDAP](https://zh.wikipedia.org/wiki/%E8%BD%BB%E5%9E%8B%E7%9B%AE%E5%BD%95%E8%AE%BF%E9%97%AE%E5%8D%8F%E8%AE%AE)
&#x20;over [TLS](https://zh.wikipedia.org/wiki/%E4%BC%A0%E8%BE%93%E5%B1%82%E5%AE%89%E5%85%A8%E5%8D%8F%E8%AE%AE)
（加密传输，也被称为 LDAPS） | 官方 |
| 639/TCP,UDP | MSDP，[组播源发现协议](https://zh.wikipedia.org/wiki/%E7%BB%84%E6%92%AD%E6%BA%90%E5%8F%91%E7%8E%B0%E5%8D%8F%E8%AE%AE) | 官方 |
| 646/TCP,UDP | LDP，标签分发协议 | 官方 |
| 647/TCP | DHCP 故障转移协议 | 官方 |
| 648/TCP | RRP（Registry Registrar Protocol），注册表注册协议 | 官方 |
| 652/TCP | DTCP（Dynamic Tunnel Configuration Protocol），[动态主机设置协议](https://zh.wikipedia.org/wiki/%E5%8A%A8%E6%80%81%E4%B8%BB%E6%9C%BA%E8%AE%BE%E7%BD%AE%E5%8D%8F%E8%AE%AE) | 官方 |
| 654/UDP | AODV（Ad hoc On-Demand Distance Vector），[无线自组网按需平面距离向量路由协议](https://zh.wikipedia.org/wiki/%E6%97%A0%E7%BA%BF%E8%87%AA%E7%BB%84%E7%BD%91%E6%8C%89%E9%9C%80%E5%B9%B3%E9%9D%A2%E8%B7%9D%E7%A6%BB%E7%9F%A2%E9%87%8F%E8%B7%AF%E7%94%B1%E5%8D%8F%E8%AE%AE) | 官方 |
| 665/TCP | sun-dr, Remote Dynamic Reconfiguration | 非官方 |
| 666/UDP | [毁灭战士](https://zh.wikipedia.org/wiki/%E6%AF%81%E7%81%AD%E6%88%98%E5%A3%AB)
，电脑平台上的一系列[第一人称射击游戏](https://zh.wikipedia.org/wiki/%E7%AC%AC%E4%B8%80%E4%BA%BA%E7%A7%B0%E5%B0%84%E5%87%BB%E6%B8%B8%E6%88%8F)
。 | 官方 |
| 674/TCP | ACAP（Application Configuration Access Protocol），应用配置访问协议 | 官方 |
| 691/TCP | MS Exchange Routing | 官方 |
| 692/TCP | Hyperwave-ISP | |
| 694/UDP | 用于带有高可用性的聚类的心跳服务 | 非官方 |
| 695/TCP | IEEE-MMS-SSL | |
| 698/UDP | OLSR（Optimized Link State Routing），优化链路状态路由协议 | |
| 699/TCP | Access Network | |
| 700/TCP | EPP, 可扩展供应协议 | |
| 701/TCP | LMP,链路管理协议 | |
| 702/TCP | IRIS over BEEP | |
| 706/TCP | SILC，Secure Internet Live Conferencing | |
| 711/TCP | TDP，标签分发协议 | |
| 712/TCP | TBRPF，Topology Broadcast based on Reverse-Path Forwarding | |
| 720/TCP | SMQP，Simple Message Queue Protocol | |
| 749/TCP,UDP | kerberos-adm，Kerberos administration | |
| 750/UDP | Kerberos version IV | |
| 782/TCP | [Conserver](https://zh.wikipedia.org/w/index.php?title=Conserver&action=edit&redlink=1)
&#x20;serial-console management server | |
| 829/TCP | [证书管理协议](https://zh.wikipedia.org/wiki/%E8%AF%81%E4%B9%A6%E7%AE%A1%E7%90%86%E5%8D%8F%E8%AE%AE)
（CMP）[\[2\]](https://zh.wikipedia.org/wiki/TCP/UDP%E7%AB%AF%E5%8F%A3%E5%88%97%E8%A1%A8#cite_note-rfc4210-2) | |
| 860/TCP | [ISCSI](https://zh.wikipedia.org/wiki/ISCSI)
，Internet 小型计算机系统接口 | |
| 873/TCP | [Rsync](https://zh.wikipedia.org/wiki/Rsync)
，文件同步协议 | 官方 |
| 901/TCP | [Samba](https://zh.wikipedia.org/wiki/Samba)
&#x20;网络管理工具（SWAT） | 非官方 |
| 902 | [VMware](https://zh.wikipedia.org/wiki/VMware)
服务器控制台[\[3\]](https://zh.wikipedia.org/wiki/TCP/UDP%E7%AB%AF%E5%8F%A3%E5%88%97%E8%A1%A8#cite_note-3) | 非官方 |
| 904 | [VMware](https://zh.wikipedia.org/wiki/VMware)
服务器替代（如果 902 端口已被占用） | 非官方 |
| 911/TCP | [Network Console on Acid](https://zh.wikipedia.org/w/index.php?title=Network_Console_on_Acid&action=edit&redlink=1)
（NCA） - local [tty](https://zh.wikipedia.org/wiki/Linux)
&#x20;redirection over [OpenSSH](https://zh.wikipedia.org/wiki/OpenSSH) | |
| 981/TCP | [Check Point](https://zh.wikipedia.org/wiki/Check_Point)
&#x20;Remote HTTPS management for firewall devices running embedded [Checkpoint Firewall-1](https://zh.wikipedia.org/w/index.php?title=Checkpoint_Firewall-1&action=edit&redlink=1)
&#x20;software | 非官方 |
| 989/TCP,UDP | FTP Protocol (data) over TLS/SSL | 官方 |
| 990/TCP,UDP | FTP Protocol (control) over TLS/SSL | 官方 |
| 991/TCP,UDP | NAS (Netnews Admin System) | |
| 992/TCP,UDP | 基于 TLS/SSL 的 Telnet 协议 | 官方 |
| 993/TCP | 基于 [传输层安全性协议](https://zh.wikipedia.org/wiki/%E5%82%B3%E8%BC%B8%E5%B1%A4%E5%AE%89%E5%85%A8%E6%80%A7%E5%8D%94%E5%AE%9A)
的[因特网信息访问协议](https://zh.wikipedia.org/wiki/%E5%9B%A0%E7%89%B9%E7%BD%91%E4%BF%A1%E6%81%AF%E8%AE%BF%E9%97%AE%E5%8D%8F%E8%AE%AE)
&#x20;(加密传输) | 官方 |
| 995/TCP | 基于 [传输层安全性协议](https://zh.wikipedia.org/wiki/%E5%82%B3%E8%BC%B8%E5%B1%A4%E5%AE%89%E5%85%A8%E6%80%A7%E5%8D%94%E5%AE%9A)
的[邮局协议](https://zh.wikipedia.org/wiki/%E9%83%B5%E5%B1%80%E5%8D%94%E5%AE%9A)
&#x20;(加密传输) | 官方 |

## 1025 到 49151 号端口

以下列表仅列出常用端口，详细的列表请参阅网站。

| 端口                                                                                                                                 | 描述                                                                                                                                                      | 状态   |
| ------------------------------------------------------------------------------------------------------------------------------------ | --------------------------------------------------------------------------------------------------------------------------------------------------------- | ------ |
| 1025/tcp                                                                                                                             | NFS-or-IIS                                                                                                                                                | 非官方 |
| 1026/tcp                                                                                                                             | 通常用于微软[Distributed COM](https://zh.wikipedia.org/wiki/Distributed_COM)                                                                              |
| 服务器                                                                                                                               | 非官方                                                                                                                                                    |
| 1029/tcp                                                                                                                             | 通常用于微软[Distributed COM](https://zh.wikipedia.org/wiki/Distributed_COM)                                                                              |
| 服务器                                                                                                                               | 非官方                                                                                                                                                    |
| 1058/tcp                                                                                                                             | nim [IBM AIX](https://zh.wikipedia.org/wiki/IBM_AIX)                                                                                                      | 官方   |
| 1059/tcp                                                                                                                             | nimreg                                                                                                                                                    | 官方   |
| 1080/tcp                                                                                                                             | [SOCKS](https://zh.wikipedia.org/wiki/SOCKS)                                                                                                              |
| 代理                                                                                                                                 | 官方                                                                                                                                                      |
| 1099/tcp,udp                                                                                                                         | [Java 远程方法调用](https://zh.wikipedia.org/wiki/Java%E8%BF%9C%E7%A8%8B%E6%96%B9%E6%B3%95%E8%B0%83%E7%94%A8)                                             |
| &#x20;Registry                                                                                                                       | 官方                                                                                                                                                      |
| 1109/tcp                                                                                                                             | Kerberos POP                                                                                                                                              |        |
| 1140/tcp                                                                                                                             | AutoNOC                                                                                                                                                   | 官方   |
| 1167/udp                                                                                                                             | phone, conference calling                                                                                                                                 |        |
| 1176/tcp                                                                                                                             | Perceptive Automation Indigo home control server                                                                                                          | 官方   |
| 1182/tcp,udp                                                                                                                         | [AcceleNet](https://zh.wikipedia.org/w/index.php?title=AcceleNet&action=edit&redlink=1)                                                                   | 官方   |
| 1194/udp                                                                                                                             | [OpenVPN](https://zh.wikipedia.org/wiki/OpenVPN)                                                                                                          | 官方   |
| 1198/tcp,udp                                                                                                                         | The [cajo project](https://zh.wikipedia.org/w/index.php?title=Cajo_project&action=edit&redlink=1)                                                         |
| &#x20;Free dynamic transparent distributed computing in Java                                                                         | 官方                                                                                                                                                      |
| 1200/udp                                                                                                                             | [Steam](https://zh.wikipedia.org/wiki/Steam)                                                                                                              | 官方   |
| 1214/tcp                                                                                                                             | [Kazaa](https://zh.wikipedia.org/w/index.php?title=Kazaa&action=edit&redlink=1)                                                                           | 官方   |
| 1223/tcp,udp                                                                                                                         | TGP: TrulyGlobal Protocol                                                                                                                                 | 官方   |
| 1241/tcp,udp                                                                                                                         | [Nessus](https://zh.wikipedia.org/wiki/Nessus)                                                                                                            |
| &#x20;Security Scanner                                                                                                               | 官方                                                                                                                                                      |
| 1248/tcp                                                                                                                             | NSClient/NSClient++/NC_Net (Nagios)                                                                                                                       | 非官方 |
| 1270/tcp,udp                                                                                                                         | Microsoft Operations Manager 2005 agent (MOM 2005)                                                                                                        | 官方   |
| 1311/tcp                                                                                                                             | Dell Open Manage Https Port                                                                                                                               | 非官方 |
| 1313/tcp                                                                                                                             | Xbiim (Canvii server) Port                                                                                                                                | 非官方 |
| 1337/tcp                                                                                                                             | [WASTE](https://zh.wikipedia.org/w/index.php?title=WASTE&action=edit&redlink=1)                                                                           |
| &#x20;Encrypted File Sharing Program                                                                                                 | 非官方                                                                                                                                                    |
| 1352/tcp                                                                                                                             | IBM [IBM Lotus Notes](https://zh.wikipedia.org/wiki/IBM_Lotus_Notes)                                                                                      |
| /Domino RPC                                                                                                                          | 官方                                                                                                                                                      |
| 1387/tcp,udp                                                                                                                         | Computer Aided Design Software Inc LM (cadsi-lm )                                                                                                         | 官方   |
| 1414/tcp                                                                                                                             | IBM [MQSeries](https://zh.wikipedia.org/w/index.php?title=MQSeries&action=edit&redlink=1)                                                                 | 官方   |
| 1431/tcp                                                                                                                             | RGTP                                                                                                                                                      | 官方   |
| 1433/tcp,udp                                                                                                                         | [Microsoft SQL](https://zh.wikipedia.org/wiki/Microsoft_SQL_Server)                                                                                       |
| &#x20;数据库系统                                                                                                                     | 官方                                                                                                                                                      |
| 1434/tcp,udp                                                                                                                         | Microsoft SQL 活动监视器                                                                                                                                  | 官方   |
| 1494/tcp                                                                                                                             | [思杰系统](https://zh.wikipedia.org/wiki/%E6%80%9D%E6%9D%B0%E7%B3%BB%E7%BB%9F)                                                                            |
| &#x20;ICA Client                                                                                                                     | 官方                                                                                                                                                      |
| 1512/tcp,udp                                                                                                                         | [WINS](https://zh.wikipedia.org/wiki/WINS)                                                                                                                |        |
| 1521/tcp                                                                                                                             | [nCube](https://zh.wikipedia.org/w/index.php?title=NCube&action=edit&redlink=1)                                                                           |
| &#x20;License Manager                                                                                                                | 官方                                                                                                                                                      |
| 1521/tcp                                                                                                                             | [Oracle 数据库](https://zh.wikipedia.org/wiki/Oracle%E6%95%B0%E6%8D%AE%E5%BA%93)                                                                          |
| &#x20;default listener, in future releases official port 2483                                                                        | 非官方                                                                                                                                                    |
| 1524/tcp,udp                                                                                                                         | ingreslock, ingres                                                                                                                                        | 官方   |
| 1526/tcp                                                                                                                             | [Oracle 数据库](https://zh.wikipedia.org/wiki/Oracle%E6%95%B0%E6%8D%AE%E5%BA%93)                                                                          |
| &#x20;common alternative for listener                                                                                                | 非官方                                                                                                                                                    |
| 1533/tcp                                                                                                                             | IBM [Lotus Sametime](https://zh.wikipedia.org/wiki/Lotus_Sametime)                                                                                        |
| &#x20;IM - Virtual Places Chat                                                                                                       | 官方                                                                                                                                                      |
| 1547/tcp,udp                                                                                                                         | [Laplink](https://zh.wikipedia.org/w/index.php?title=Laplink&action=edit&redlink=1)                                                                       | 官方   |
| 1550                                                                                                                                 | [Gadu-Gadu](https://zh.wikipedia.org/w/index.php?title=Gadu-Gadu&action=edit&redlink=1)                                                                   |
| &#x20;(Direct Client-to-Client)                                                                                                      | 非官方                                                                                                                                                    |
| 1581/udp                                                                                                                             | [MIL STD 2045-47001 VMF](https://zh.wikipedia.org/w/index.php?title=Combat-net_radio&action=edit&redlink=1)                                               | 官方   |
| 1589/udp                                                                                                                             | Cisco [VQP](https://zh.wikipedia.org/w/index.php?title=VQP&action=edit&redlink=1)                                                                         |
| &#x20;(VLAN Query Protocol) / [VMPS](https://zh.wikipedia.org/w/index.php?title=VLAN_Management_Policy_Server&action=edit&redlink=1) | 非官方                                                                                                                                                    |
| 1627                                                                                                                                 | iSketch                                                                                                                                                   | 非官方 |
| 1677/tcp                                                                                                                             | Novell GroupWise clients in client/server access mode                                                                                                     |        |
| 1701/udp                                                                                                                             | [第二层隧道协议](https://zh.wikipedia.org/wiki/%E7%AC%AC%E4%BA%8C%E5%B1%82%E9%9A%A7%E9%81%93%E5%8D%8F%E8%AE%AE)                                           |
| , Layer 2 Tunnelling protocol                                                                                                        |                                                                                                                                                           |
| 1716/tcp                                                                                                                             | [美国陆军系列](https://zh.wikipedia.org/wiki/%E7%BE%8E%E5%9C%8B%E9%99%B8%E8%BB%8D%E7%B3%BB%E5%88%97)                                                      |
| &#x20;MMORPG Default Game Port                                                                                                       | 官方                                                                                                                                                      |
| 1723/tcp,udp                                                                                                                         | Microsoft [点对点隧道协议](https://zh.wikipedia.org/wiki/%E9%BB%9E%E5%B0%8D%E9%BB%9E%E9%9A%A7%E9%81%93%E5%8D%94%E8%AD%B0)                                 |
| &#x20;VPN                                                                                                                            | 官方                                                                                                                                                      |
| 1725/udp                                                                                                                             | Valve Steam Client                                                                                                                                        | 非官方 |
| 1755/tcp,udp                                                                                                                         | [MMS (协议)](<https://zh.wikipedia.org/wiki/MMS_(%E5%8D%8F%E8%AE%AE)>)                                                                                    |
| &#x20;(MMS, ms-streaming)                                                                                                            | 官方                                                                                                                                                      |
| 1761/tcp,udp                                                                                                                         | cft-0                                                                                                                                                     | 官方   |
| 1761/tcp                                                                                                                             | [Novell](https://zh.wikipedia.org/wiki/Novell)                                                                                                            |
| &#x20;Zenworks Remote Control utility                                                                                                | 非官方                                                                                                                                                    |
| 1762-1768/tcp,udp                                                                                                                    | cft-1 to cft-7                                                                                                                                            | 官方   |
| 1812/udp                                                                                                                             | radius, [远端用户拨入验证服务](https://zh.wikipedia.org/wiki/%E8%BF%9C%E7%AB%AF%E7%94%A8%E6%88%B7%E6%8B%A8%E5%85%A5%E9%AA%8C%E8%AF%81%E6%9C%8D%E5%8A%A1)  |
| &#x20;authentication protocol                                                                                                        |                                                                                                                                                           |
| 1813/udp                                                                                                                             | radacct, [远端用户拨入验证服务](https://zh.wikipedia.org/wiki/%E8%BF%9C%E7%AB%AF%E7%94%A8%E6%88%B7%E6%8B%A8%E5%85%A5%E9%AA%8C%E8%AF%81%E6%9C%8D%E5%8A%A1) |
| &#x20;accounting protocol                                                                                                            |                                                                                                                                                           |
| 1863/tcp                                                                                                                             | [Windows Live Messenger](https://zh.wikipedia.org/wiki/Windows_Live_Messenger)                                                                            |
| , MSN                                                                                                                                | 官方                                                                                                                                                      |
| 1900/udp                                                                                                                             | Microsoft [简单服务发现协议](https://zh.wikipedia.org/wiki/%E7%AE%80%E5%8D%95%E6%9C%8D%E5%8A%A1%E5%8F%91%E7%8E%B0%E5%8D%8F%E8%AE%AE)                      |

&#x20;Enables discovery of [UPnP](https://zh.wikipedia.org/wiki/UPnP)
&#x20;devices | 官方 |
| 1935/tcp | [实时消息协议](https://zh.wikipedia.org/wiki/%E5%AE%9E%E6%97%B6%E6%B6%88%E6%81%AF%E5%8D%8F%E8%AE%AE)
&#x20;(RTMP) raw protocol | 官方 |
| 1970/tcp,udp | [Danware Data](https://zh.wikipedia.org/w/index.php?title=Danware_Data_A/S&action=edit&redlink=1)
&#x20;NetOp Remote Control | 官方 |
| 1971/tcp,udp | [Danware Data](https://zh.wikipedia.org/w/index.php?title=Danware_Data_A/S&action=edit&redlink=1)
&#x20;NetOp School | 官方 |
| 1972/tcp,udp | [InterSystems Caché](https://zh.wikipedia.org/w/index.php?title=InterSystems_Cach%C3%A9&action=edit&redlink=1) | 官方 |
| 1975-77/udp | Cisco [TCO](https://zh.wikipedia.org/w/index.php?title=Total_cost_of_ownership&action=edit&redlink=1)
&#x20;([Documentation](http://www.cisco.com/en/US/netsol/networking_solutions_networking_basic09186a00800a3524.html)
) | 官方 |
| 1984/tcp | Big Brother - network monitoring tool | 官方 |
| 1985/udp | [热备份路由器协议](https://zh.wikipedia.org/wiki/%E7%86%B1%E5%82%99%E4%BB%BD%E8%B7%AF%E7%94%B1%E5%99%A8%E5%8D%94%E5%AE%9A) | 官方 |
| 1994/TCP | STUN-SDLC protocol for tunneling | |
| 1998/tcp | Cisco X.25 service (XOT) | |
| 2000/tcp,udp | [Cisco SCCP (Skinny)](https://zh.wikipedia.org/w/index.php?title=Skinny_Client_Control_Protocol&action=edit&redlink=1) | 官方 |
| 2002/tcp | Cisco Secure Access Control Server (ACS) for Windows | 非官方 |
| 2030 | [甲骨文公司](https://zh.wikipedia.org/wiki/%E7%94%B2%E9%AA%A8%E6%96%87%E5%85%AC%E5%8F%B8)
&#x20;Services for [Microsoft Transaction Server](https://zh.wikipedia.org/w/index.php?title=Microsoft_Transaction_Server&action=edit&redlink=1) | 非官方 |
| 2031/tcp,udp | [mobrien-chat](https://zh.wikipedia.org/w/index.php?title=Mobrien-chat&action=edit&redlink=1)
&#x20;\- Mike O'Brien <mike@mobrien.com> November 2004 | 官方 |
| 2049/udp | nfs, [NFS](https://zh.wikipedia.org/wiki/NFS)
&#x20;Server | 官方 |
| 2049/udp | shilp | 官方 |
| 2053/tcp | knetd, [Kerberos](https://zh.wikipedia.org/wiki/Kerberos)
&#x20;de-multiplexor | |
| 2056/udp | [文明 IV](https://zh.wikipedia.org/wiki/%E6%96%87%E6%98%8EIV)
&#x20;multiplayer | 非官方 |
| 2073/tcp,udp | DataReel Database | 官方 |
| 2074/tcp,udp | Vertel VMF SA (i.e. App.. SpeakFreely) | 官方 |
| 2082/tcp | Infowave Mobility Server | 官方 |
| 2082/tcp | [CPanel](https://zh.wikipedia.org/wiki/CPanel)
, default port | 非官方 |
| 2083/tcp | Secure Radius Service (radsec) | 官方 |
| 2083/tcp | [CPanel](https://zh.wikipedia.org/wiki/CPanel)
&#x20;default SSL port | 非官方 |
| 2086/tcp | [GNUnet](https://zh.wikipedia.org/w/index.php?title=GNUnet&action=edit&redlink=1) | 官方 |
| 2086/tcp | [CPanel](https://zh.wikipedia.org/wiki/CPanel)
&#x20;default port | 非官方 |
| 2087/tcp | [CPanel](https://zh.wikipedia.org/wiki/CPanel)
&#x20;default SSL port | 非官方 |
| 2095/tcp | [CPanel](https://zh.wikipedia.org/wiki/CPanel)
&#x20;default webmail port | 非官方 |
| 2096/tcp | [CPanel](https://zh.wikipedia.org/wiki/CPanel)
&#x20;default SSL webmail port | 非官方 |
| 2161/tcp | [问号](https://zh.wikipedia.org/wiki/%E9%97%AE%E5%8F%B7)
-APC Agent | 官方 |
| 2181/tcp,udp | [EForward](https://zh.wikipedia.org/w/index.php?title=EForward&action=edit&redlink=1)
-document transport system | 官方 |
| 2200/tcp | Tuxanci game server | 非官方 |
| 2219/tcp,udp | NetIQ NCAP Protocol | 官方 |
| 2220/tcp,udp | NetIQ End2End | 官方 |
| 2222/tcp | [DirectAdmin](https://zh.wikipedia.org/w/index.php?title=DirectAdmin&action=edit&redlink=1)
's default port | 非官方 |
| 2222/udp | Microsoft Office OS X antipiracy network monitor [\[1\]](https://web.archive.org/web/20080225053131/http://www.ciac.org/ciac/techbull/CIACTech02-003.shtml) | 非官方 |
| 2301/tcp | HP System Management Redirect to port 2381 | 非官方 |
| 2302/udp | [武装突袭](https://zh.wikipedia.org/wiki/%E6%AD%A6%E8%A3%85%E7%AA%81%E8%A2%AD)
&#x20;multiplayer (default for game) | 非官方 |
| 2302/udp | [最后一战：战斗进化](https://zh.wikipedia.org/wiki/%E6%9C%80%E5%BE%8C%E4%B8%80%E6%88%B0%EF%BC%9A%E6%88%B0%E9%AC%A5%E9%80%B2%E5%8C%96)
&#x20;multiplayer | 非官方 |
| 2303/udp | [武装突袭](https://zh.wikipedia.org/wiki/%E6%AD%A6%E8%A3%85%E7%AA%81%E8%A2%AD)
&#x20;multiplayer (default for server reporting) (游戏内定端口 +1) | 非官方 |
| 2305/udp | [武装突袭](https://zh.wikipedia.org/wiki/%E6%AD%A6%E8%A3%85%E7%AA%81%E8%A2%AD)
&#x20;multiplayer (default for VoN) (游戏内定端口 +3) | 非官方 |
| 2369/tcp | Default port for [BMC 软件公司](https://zh.wikipedia.org/wiki/BMC%E8%BB%9F%E4%BB%B6%E5%85%AC%E5%8F%B8)
&#x20;CONTROL-M/Server - Configuration Agent port number - though often changed during installation | 非官方 |
| 2370/tcp | Default port for [BMC 软件公司](https://zh.wikipedia.org/wiki/BMC%E8%BB%9F%E4%BB%B6%E5%85%AC%E5%8F%B8)
&#x20;CONTROL-M/Server - Port utilized to allow the CONTROL-M/Enterprise Manager to connect to the CONTROL-M/Server - though often changed during installation | 非官方 |
| 2381/tcp | HP Insight Manager default port for webserver | 非官方 |
| 2404/tcp | IEC 60870-5-104 | 官方 |
| 2427/udp | Cisco [MGCP](https://zh.wikipedia.org/w/index.php?title=MGCP&action=edit&redlink=1) | 官方 |
| 2447/tcp,udp | ovwdb - [OpenView](https://zh.wikipedia.org/w/index.php?title=OpenView&action=edit&redlink=1)
&#x20;[Network Node Manager](https://zh.wikipedia.org/w/index.php?title=Network_Node_Manager&action=edit&redlink=1)
&#x20;(NNM) daemon | 官方 |
| 2483/tcp,udp | [Oracle 数据库](https://zh.wikipedia.org/wiki/Oracle%E6%95%B0%E6%8D%AE%E5%BA%93)
&#x20;listening port for unsecure client connections to the listener, replaces port 1521 | 官方 |
| 2484/tcp,udp | [Oracle 数据库](https://zh.wikipedia.org/wiki/Oracle%E6%95%B0%E6%8D%AE%E5%BA%93)
&#x20;listening port for SSL client connections to the listener | 官方 |
| 2546/tcp,udp | Vytal Vault - Data Protection Services | 非官方 |
| 2593/tcp,udp | RunUO - [网络创世纪](https://zh.wikipedia.org/wiki/%E7%BD%91%E7%BB%9C%E5%88%9B%E4%B8%96%E7%BA%AA)
&#x20;server | 非官方 |
| 2598/tcp | new ICA - when Session Reliability is enabled, TCP port 2598 replaces port 1494 | 非官方 |
| 2612/tcp,udp | QPasa from MQSoftware | 官方 |
| 2710/tcp | XBT Bittorrent Tracker | 非官方 |
| 2710/udp | XBT Bittorrent Tracker experimental UDP tracker extension | 非官方 |
| 2710/tcp | Knuddels.de | 非官方 |
| 2735/tcp,udp | NetIQ Monitor Console | 官方 |
| 2809/tcp | corbaloc:iiop URL, per the [CORBA](https://zh.wikipedia.org/wiki/CORBA)
&#x20;3.0.3 specification.
Also used by IBM [IBM WebSphere Application Server](https://zh.wikipedia.org/wiki/IBM_WebSphere_Application_Server) Node Agent | 官方 |
| 2809/udp | corbaloc:iiop URL, per the [CORBA](https://zh.wikipedia.org/wiki/CORBA)
&#x20;3.0.3 specification. | 官方 |
| 2944/udp | [Megaco](https://zh.wikipedia.org/wiki/Megaco)
&#x20;Text H.248 | 非官方 |
| 2945/udp | [Megaco](https://zh.wikipedia.org/wiki/Megaco)
&#x20;Binary (ASN.1) H.248 | 非官方 |
| 2948/tcp,udp | [无线应用协议](https://zh.wikipedia.org/wiki/%E6%97%A0%E7%BA%BF%E5%BA%94%E7%94%A8%E5%8D%8F%E8%AE%AE)
-push [彩信](https://zh.wikipedia.org/wiki/%E5%A4%9A%E5%AA%92%E9%AB%94%E7%9F%AD%E8%A8%8A)
&#x20;(MMS) | 官方 |
| 2949/tcp,udp | [无线应用协议](https://zh.wikipedia.org/wiki/%E6%97%A0%E7%BA%BF%E5%BA%94%E7%94%A8%E5%8D%8F%E8%AE%AE)
-pushsecure [彩信](https://zh.wikipedia.org/wiki/%E5%A4%9A%E5%AA%92%E9%AB%94%E7%9F%AD%E8%A8%8A)
&#x20;(MMS) | 官方 |
| 2967/tcp | Symantec AntiVirus Corporate Edition | 非官方 |
| 3000/tcp | Miralix License server | 非官方 |
| 3000/udp | [Distributed Interactive Simulation](https://zh.wikipedia.org/w/index.php?title=Distributed_Interactive_Simulation&action=edit&redlink=1)
&#x20;(DIS), modifiable default port | 非官方 |
| 3000/tcp | Puma Web Server | 非官方 |
| 3001/tcp | Miralix Phone Monitor | 非官方 |
| 3002/tcp | Miralix CSTA | 非官方 |
| 3003/tcp | Miralix GreenBox API | 非官方 |
| 3004/tcp | Miralix InfoLink | 非官方 |
| 3006/tcp | Miralix SMS Client Connector | 非官方 |
| 3007/tcp | Miralix OM Server | 非官方 |
| 3025/tcp | netpd.org | 非官方 |
| 3050/tcp,udp | gds*db (Interbase/Firebird) | 官方 |
| 3074/tcp,udp | [Xbox Live](https://zh.wikipedia.org/wiki/Xbox_Live) | 官方 |
| 3128/tcp | [超文本传输协议](https://zh.wikipedia.org/wiki/%E8%B6%85%E6%96%87%E6%9C%AC%E4%BC%A0%E8%BE%93%E5%8D%8F%E8%AE%AE)
&#x20;used by [Web 缓存](https://zh.wikipedia.org/wiki/Web%E7%BC%93%E5%AD%98)
s and the default port for the [Squid (软件)](https://zh.wikipedia.org/wiki/Squid*(%E8%BD%AF%E4%BB%B6)) | 官方 |
| 3260/tcp | [ISCSI](https://zh.wikipedia.org/wiki/ISCSI)
&#x20;target | 官方 |
| 3268/tcp | msft-gc, Microsoft Global Catalog ([轻型目录访问协议](https://zh.wikipedia.org/wiki/%E8%BD%BB%E5%9E%8B%E7%9B%AE%E5%BD%95%E8%AE%BF%E9%97%AE%E5%8D%8F%E8%AE%AE)
&#x20;service which contains data from [Active Directory](https://zh.wikipedia.org/wiki/Active_Directory)
&#x20;forests) | 官方 |
| 3269/tcp | msft-gc-ssl, Microsoft Global Catalog over SSL (similar to port 3268, [轻型目录访问协议](https://zh.wikipedia.org/wiki/%E8%BD%BB%E5%9E%8B%E7%9B%AE%E5%BD%95%E8%AE%BF%E9%97%AE%E5%8D%8F%E8%AE%AE)
&#x20;over [传输层安全性协议](https://zh.wikipedia.org/wiki/%E5%82%B3%E8%BC%B8%E5%B1%A4%E5%AE%89%E5%85%A8%E6%80%A7%E5%8D%94%E5%AE%9A)
&#x20;version) | 官方 |
| 3300/tcp | [TripleA](<https://zh.wikipedia.org/w/index.php?title=TripleA_(computer_game)&action=edit&redlink=1>)
&#x20;game server | 非官方 |
| 3305/tcp,udp | [ODETTE-FTP](https://zh.wikipedia.org/w/index.php?title=ODETTE-FTP&action=edit&redlink=1) | 官方 |
| 3306/tcp,udp | [MySQL](https://zh.wikipedia.org/wiki/MySQL)
数据库系统 | 官方 |
| 3333/tcp | Network Caller ID server | 非官方 |
| 3386/tcp,udp | [GTP'](https://zh.wikipedia.org/wiki/GTP%27)
&#x20;[3GPP](https://zh.wikipedia.org/wiki/3GPP)
&#x20;[GSM](https://zh.wikipedia.org/wiki/GSM)
/[通用移动通讯系统](https://zh.wikipedia.org/wiki/%E9%80%9A%E7%94%A8%E7%A7%BB%E5%8A%A8%E9%80%9A%E8%AE%AF%E7%B3%BB%E7%BB%9F)
&#x20;[CDR](https://zh.wikipedia.org/w/index.php?title=Call_detail_record&action=edit&redlink=1)
&#x20;logging protocol | 官方 |
| 3389/tcp | [远程桌面协议](https://zh.wikipedia.org/wiki/%E9%81%A0%E7%AB%AF%E6%A1%8C%E9%9D%A2%E5%8D%94%E5%AE%9A)
（RDP） | 官方 |
| 3396/tcp | [Novell](https://zh.wikipedia.org/wiki/Novell)
&#x20;NDPS Printer Agent | 官方 |
| 3689/tcp | [DAAP](https://zh.wikipedia.org/w/index.php?title=Digital_Audio_Access_Protocol&action=edit&redlink=1)
&#x20;Digital Audio Access Protocol used by [苹果公司](https://zh.wikipedia.org/wiki/%E8%98%8B%E6%9E%9C%E5%85%AC%E5%8F%B8)
&#x20;[ITunes](https://zh.wikipedia.org/wiki/ITunes) | 官方 |
| 3690/tcp | [Subversion](https://zh.wikipedia.org/wiki/Subversion)
&#x20;version control system | 官方 |
| 3702/tcp,udp | [Web Services Dynamic Discovery](https://zh.wikipedia.org/w/index.php?title=Web_Services_Dynamic_Discovery&action=edit&redlink=1)
&#x20;(WS-Discovery), used by various components of [Windows Vista](https://zh.wikipedia.org/wiki/Windows_Vista) | 官方 |
| 3724/tcp | [魔兽世界](https://zh.wikipedia.org/wiki/%E9%AD%94%E5%85%BD%E4%B8%96%E7%95%8C)
&#x20;Online gaming MMORPG | 官方 |
| 3784/tcp,udp | [Ventrilo](https://zh.wikipedia.org/w/index.php?title=Ventrilo&action=edit&redlink=1)
&#x20;VoIP program used by [Ventrilo](https://zh.wikipedia.org/w/index.php?title=Ventrilo&action=edit&redlink=1) | 官方 |
| 3785/udp | [Ventrilo](https://zh.wikipedia.org/w/index.php?title=Ventrilo&action=edit&redlink=1)
&#x20;VoIP program used by [Ventrilo](https://zh.wikipedia.org/w/index.php?title=Ventrilo&action=edit&redlink=1) | 官方 |
| 3868 tcp,udp | [Diameter base protocol](https://zh.wikipedia.org/w/index.php?title=Diameter_base_protocol&action=edit&redlink=1) | 官方 |
| 3872/tcp | Oracle Management Remote Agent | 非官方 |
| 3899/tcp | [Remote Administrator](https://zh.wikipedia.org/w/index.php?title=Remote_Administrator&action=edit&redlink=1) | 非官方 |
| 3900/tcp | [Unidata UDT OS](https://zh.wikipedia.org/w/index.php?title=Unidata_UDT_OS&action=edit&redlink=1)
&#x20;udt*os | 官方 |
| 3945/tcp | Emcads server service port, a Giritech product used by [G/On](https://zh.wikipedia.org/w/index.php?title=G/On&action=edit&redlink=1) | 官方 |
| 4000/tcp | [暗黑破坏神 II](https://zh.wikipedia.org/wiki/%E6%9A%97%E9%BB%91%E7%A0%B4%E5%A3%9E%E7%A5%9EII)
&#x20;game
[NoMachine Network Server (nxd)](https://zh.wikipedia.org/w/index.php?title=NX_technology&action=edit&redlink=1) | 非官方 |
| 4007/tcp | PrintBuzzer printer monitoring socket server | 非官方 |
| 4089/tcp,udp | OpenCORE Remote Control Service | 官方 |
| 4093/tcp,udp | PxPlus Client server interface [ProvideX](https://zh.wikipedia.org/w/index.php?title=ProvideX&action=edit&redlink=1) | 官方 |
| 4096/udp | Bridge-Relay Element [ASCOM](https://zh.wikipedia.org/w/index.php?title=ASCOM&action=edit&redlink=1) | 官方 |
| 4100 | WatchGuard Authentication Applet - default port | 非官方 |
| 4111/tcp,udp | [Xgrid](https://zh.wikipedia.org/w/index.php?title=Xgrid&action=edit&redlink=1) | 官方 |
| 4111/tcp | [SharePoint](https://zh.wikipedia.org/wiki/SharePoint)
&#x20;\- 默认管理端口 | 非官方 |
| 4226/tcp,udp | [Aleph One (computer game)](https://zh.wikipedia.org/w/index.php?title=Aleph_One*(computer*game)&action=edit&redlink=1) | 非官方 |
| 4224/tcp | [思科系统](https://zh.wikipedia.org/wiki/%E6%80%9D%E7%A7%91%E7%B3%BB%E7%BB%9F)
&#x20;CDP Cisco discovery Protocol | 非官方 |
| 4569/udp | [Inter-Asterisk eXchange](https://zh.wikipedia.org/w/index.php?title=Inter-Asterisk_eXchange&action=edit&redlink=1) | 非官方 |
| 4662/tcp | OrbitNet Message Service | 官方 |
| 4662/tcp | 通常用于[EMule](https://zh.wikipedia.org/wiki/EMule) | 非官方 |
| 4664/tcp | [Google 桌面搜索](https://zh.wikipedia.org/wiki/Google%E6%A1%8C%E9%9D%A2) | 非官方 |
| 4672/udp | [EMule](https://zh.wikipedia.org/wiki/EMule)
&#x20;\- 常用端口 | 非官方 |
| 4894/tcp | [LysKOM](https://zh.wikipedia.org/w/index.php?title=LysKOM&action=edit&redlink=1)
&#x20;Protocol A | 官方 |
| 4899/tcp | [Radmin](https://zh.wikipedia.org/w/index.php?title=Radmin&action=edit&redlink=1)
&#x20;远程控制工具 | 官方 |
| 5000/tcp | commplex-main | 官方 |
| 5000/tcp | [UPnP](https://zh.wikipedia.org/wiki/UPnP)
&#x20;\- Windows network device interoperability | 非官方 |
| 5000/tcp,udp | [VTun](https://zh.wikipedia.org/w/index.php?title=VTun&action=edit&redlink=1)
&#x20;\- [虚拟专用网](https://zh.wikipedia.org/wiki/%E8%99%9B%E6%93%AC%E7%A7%81%E4%BA%BA%E7%B6%B2%E8%B7%AF)
&#x20;软件 | 非官方 |
| 5001/tcp,udp | Iperf (Tool for measuring TCP and UDP bandwidth performance) | 非官方 |
| 5001/tcp | [Slingbox](https://zh.wikipedia.org/w/index.php?title=Slingbox&action=edit&redlink=1)
及 Slingplayer | 非官方 |
| 5003/tcp | [FileMaker](https://zh.wikipedia.org/wiki/FileMaker)
&#x20;Filemaker Pro | 官方 |
| 5004/udp | [实时传输协议](https://zh.wikipedia.org/wiki/%E5%AE%9E%E6%97%B6%E4%BC%A0%E8%BE%93%E5%8D%8F%E8%AE%AE)
实时传输协议 | 官方 |
| 5005/udp | [实时传输协议](https://zh.wikipedia.org/wiki/%E5%AE%9E%E6%97%B6%E4%BC%A0%E8%BE%93%E5%8D%8F%E8%AE%AE)
实时传输协议 | 官方 |
| 5031/tcp,udp | AVM CAPI-over-TCP ([综合业务数字网](https://zh.wikipedia.org/wiki/%E7%BB%BC%E5%90%88%E4%B8%9A%E5%8A%A1%E6%95%B0%E5%AD%97%E7%BD%91)
&#x20;over [以太网](https://zh.wikipedia.org/wiki/%E4%BB%A5%E5%A4%AA%E7%BD%91)
&#x20;tunneling) | 非官方 |
| 5050/tcp | [Yahoo! Messenger](https://zh.wikipedia.org/wiki/Yahoo!_Messenger) | 官方 |
| 5051/tcp | [ita-agent](https://zh.wikipedia.org/w/index.php?title=Ita-agent&action=edit&redlink=1)
&#x20;Symantec Intruder Alert | 官方 |
| 5060/tcp,udp | [会话发起协议](https://zh.wikipedia.org/wiki/%E4%BC%9A%E8%AF%9D%E5%8F%91%E8%B5%B7%E5%8D%8F%E8%AE%AE)
&#x20;(SIP) | 官方 |
| 5061/tcp | [会话发起协议](https://zh.wikipedia.org/wiki/%E4%BC%9A%E8%AF%9D%E5%8F%91%E8%B5%B7%E5%8D%8F%E8%AE%AE)
&#x20;(SIP) over [传输层安全性协议](https://zh.wikipedia.org/wiki/%E5%82%B3%E8%BC%B8%E5%B1%A4%E5%AE%89%E5%85%A8%E6%80%A7%E5%8D%94%E5%AE%9A)
&#x20;(TLS) | 官方 |
| 5093/udp | [SPSS License Administrator](https://zh.wikipedia.org/w/index.php?title=SPSS_License_Administrator&action=edit&redlink=1)
&#x20;(SPSS) | 官方 |
| 5104/tcp | [IBM NetCOOL / IMPACT HTTP Service](https://zh.wikipedia.org/w/index.php?title=IBM_NetCOOL*/_IMPACT_HTTP_Service&action=edit&redlink=1) | 非官方 |
| 5106/tcp | [A-Talk](https://zh.wikipedia.org/w/index.php?title=Atalk&action=edit&redlink=1)
&#x20;Common connection | 非官方 |
| 5107/tcp | [A-Talk](https://zh.wikipedia.org/w/index.php?title=Atalk&action=edit&redlink=1)
&#x20;远程服务器连接 | 非官方 |
| 5110/tcp | ProRat Server | 非官方 |
| 5121/tcp | [无冬之夜](https://zh.wikipedia.org/wiki/%E6%97%A0%E5%86%AC%E4%B9%8B%E5%A4%9C) | 官方 |
| 5176/tcp | ConsoleWorks default UI interface | 非官方 |
| 5190/tcp | [ICQ](https://zh.wikipedia.org/wiki/ICQ)
&#x20;and [AIM (应用程序)](https://zh.wikipedia.org/wiki/AIM_(%E6%87%89%E7%94%A8%E7%A8%8B%E5%BC%8F)) | 官方 |
| 5222/tcp | [XMPP/Jabber](https://zh.wikipedia.org/wiki/Jabber)
&#x20;\- client connection | 官方 |
| 5223/tcp | [XMPP/Jabber](https://zh.wikipedia.org/wiki/Jabber)
&#x20;\- default port for SSL Client Connection | 非官方 |
| 5269/tcp | [XMPP/Jabber](https://zh.wikipedia.org/wiki/Jabber)
&#x20;\- server connection | 官方 |
| 5351/tcp,udp | [NAT 端口映射协议](https://zh.wikipedia.org/wiki/NAT%E7%AB%AF%E5%8F%A3%E6%98%A0%E5%B0%84%E5%8D%8F%E8%AE%AE)
，允许客户端在[网络地址转换](https://zh.wikipedia.org/wiki/%E7%BD%91%E7%BB%9C%E5%9C%B0%E5%9D%80%E8%BD%AC%E6%8D%A2)
网关上配置传入映射 | 官方 |
| 5353/udp | [mDNS](https://zh.wikipedia.org/w/index.php?title=MDNS&action=edit&redlink=1)
&#x20;\- 多播 DNS | |
| 5402/tcp,udp | StarBurst AutoCast MFTP | 官方 |
| 5405/tcp,udp | [NetSupport](https://zh.wikipedia.org/w/index.php?title=NetSupport&action=edit&redlink=1) | 官方 |
| 5421/tcp,udp | [Net Support 2](https://zh.wikipedia.org/w/index.php?title=NetSupport&action=edit&redlink=1) | 官方 |
| 5432/tcp | [PostgreSQL](https://zh.wikipedia.org/wiki/PostgreSQL)
数据库管理系统 | 官方 |
| 5445/udp | [思科系统](https://zh.wikipedia.org/wiki/%E6%80%9D%E7%A7%91%E7%B3%BB%E7%BB%9F)
&#x20;Vidéo VT Advantage | 非官方 |
| 5495/tcp | [Applix](https://zh.wikipedia.org/w/index.php?title=Applix&action=edit&redlink=1)
&#x20;TM1 Admin server | 非官方 |
| 5498/tcp | [Hotline](https://zh.wikipedia.org/w/index.php?title=Hotline_Communications&action=edit&redlink=1)
&#x20;tracker server connection | 非官方 |
| 5499/udp | [Hotline](https://zh.wikipedia.org/w/index.php?title=Hotline_Communications&action=edit&redlink=1)
&#x20;tracker server discovery | 非官方 |
| 5500/tcp | [VNC](https://zh.wikipedia.org/wiki/VNC)
&#x20;remote desktop protocol - for incoming listening viewer, [Hotline](https://zh.wikipedia.org/w/index.php?title=Hotline_Communications&action=edit&redlink=1)
&#x20;control connection | 非官方 |
| 5501/tcp | [Hotline](https://zh.wikipedia.org/w/index.php?title=Hotline_Communications&action=edit&redlink=1)
&#x20;file transfer connection | 非官方 |
| 5517/tcp | [Setiqueue](https://zh.wikipedia.org/w/index.php?title=Setiqueue&action=edit&redlink=1)
&#x20;Proxy server client for [SETI@home](https://zh.wikipedia.org/wiki/SETI@home)
&#x20;project | 非官方 |
| 5555/tcp | [Freeciv](https://zh.wikipedia.org/wiki/Freeciv)
&#x20;multiplay port for versions up to 2.0, [惠普](https://zh.wikipedia.org/wiki/%E6%83%A0%E6%99%AE)
&#x20;Data Protector, [会话通告协议](https://zh.wikipedia.org/wiki/%E6%9C%83%E8%A9%B1%E9%80%9A%E5%91%8A%E5%8D%94%E8%AD%B0) | 非官方 |
| 5556/tcp | [Freeciv](https://zh.wikipedia.org/wiki/Freeciv)
&#x20;multiplay port | 官方 |
| 5631/tcp | [赛门铁克](https://zh.wikipedia.org/wiki/%E8%B5%9B%E9%97%A8%E9%93%81%E5%85%8B)
&#x20;pcAnywhere | 官方 |
| 5632/udp | [赛门铁克](https://zh.wikipedia.org/wiki/%E8%B5%9B%E9%97%A8%E9%93%81%E5%85%8B)
&#x20;pcAnywhere | 官方 |
| 5666/tcp | NRPE (Nagios) | 非官方 |
| 5667/tcp | NSCA (Nagios) | 非官方 |
| 5800/tcp | [VNC](https://zh.wikipedia.org/wiki/VNC)
&#x20;remote desktop protocol - for use over [超文本传输协议](https://zh.wikipedia.org/wiki/%E8%B6%85%E6%96%87%E6%9C%AC%E4%BC%A0%E8%BE%93%E5%8D%8F%E8%AE%AE) | 非官方 |
| 5814/tcp,udp | [惠普](https://zh.wikipedia.org/wiki/%E6%83%A0%E6%99%AE)
&#x20;Support Automation (HP OpenView Self-Healing Services) | 官方 |
| 5900/tcp | [VNC](https://zh.wikipedia.org/wiki/VNC)
&#x20;remote desktop protocol (used by [ARD](https://zh.wikipedia.org/w/index.php?title=Apple_Remote_Desktop&action=edit&redlink=1)
) | 官方 |
| 6000/tcp | [X 窗口系统](https://zh.wikipedia.org/wiki/X_Window%E7%B3%BB%E7%B5%B1)
&#x20;\- used between an X client and server over the network | 官方 |
| 6001/udp | [X 窗口系统](https://zh.wikipedia.org/wiki/X_Window%E7%B3%BB%E7%B5%B1)
&#x20;\- used between an X client and server over the network | 官方 |
| 6005/tcp | Default port for [BMC 软件公司](https://zh.wikipedia.org/wiki/BMC%E8%BB%9F%E4%BB%B6%E5%85%AC%E5%8F%B8)
&#x20;CONTROL-M/Server - Socket Port number used for communication between CONTROL-M processes - though often changed during installation | 非官方 |
| 6050/tcp | Brightstor Arcserve Backup | 非官方 |
| 6051/tcp | Brightstor Arcserve Backup | 非官方 |
| 6100/tcp | Vizrt System | 非官方 |
| 6110/tcp,udp | softcm [HP SoftBench](https://zh.wikipedia.org/w/index.php?title=HP_SoftBench&action=edit&redlink=1)
&#x20;CM | 官方 |
| 6111/tcp,udp | spc [HP SoftBench](https://zh.wikipedia.org/w/index.php?title=HP_SoftBench&action=edit&redlink=1)
&#x20;Sub-Process Control | 官方 |
| 6112/tcp | dtspcd - a network daemon that accepts requests from clients to execute commands and launch applications remotely | 官方 |
| 6112/tcp | [暴雪娱乐](https://zh.wikipedia.org/wiki/%E6%9A%B4%E9%9B%AA%E5%A8%9B%E6%A8%82)
's [暴雪战网](https://zh.wikipedia.org/wiki/%E6%9A%B4%E9%9B%AA%E6%88%98%E7%BD%91)
&#x20;gaming service, [ArenaNet](https://zh.wikipedia.org/wiki/ArenaNet)
&#x20;gaming service | 官方 |
| 6129/tcp | [Dameware Remote Control](https://zh.wikipedia.org/w/index.php?title=Dameware&action=edit&redlink=1) | 非官方 |
| 6257/udp | [WinMX](https://zh.wikipedia.org/wiki/WinMX)
&#x20;（参见 6699 端口） | 非官方 |
| 6346/tcp,udp | [gnutella-svc](https://zh.wikipedia.org/wiki/Gnutella)
&#x20;([FrostWire](https://zh.wikipedia.org/wiki/FrostWire)
, [LimeWire](https://zh.wikipedia.org/wiki/LimeWire)
, [Bearshare](https://zh.wikipedia.org/w/index.php?title=Bearshare&action=edit&redlink=1)
, etc.) | 官方 |
| 6347/tcp,udp | gnutella-rtr | 官方 |
| 6379/tcp | [Redis](https://zh.wikipedia.org/wiki/Redis)
&#x20;\- Redis | 非官方 |
| 6444/tcp,udp | [Oracle Grid Engine](https://zh.wikipedia.org/wiki/Oracle_Grid_Engine)
&#x20;\- Qmaster Service | 官方 |
| 6445/tcp,udp | [Oracle Grid Engine](https://zh.wikipedia.org/wiki/Oracle_Grid_Engine)
&#x20;\- Execution Service | 官方 |
| 6502/tcp,udp | Danware Data NetOp Remote Control | 非官方 |
| 6522/tcp | [Gobby](https://zh.wikipedia.org/w/index.php?title=Gobby&action=edit&redlink=1)
&#x20;(and other libobby-based software) | 非官方 |
| 6543/udp | [Jetnet](https://zh.wikipedia.org/w/index.php?title=Jetnet&action=edit&redlink=1)
&#x20;\- default port that the [Paradigm Research & Development](https://zh.wikipedia.org/w/index.php?title=Paradigm_Research_%26_Development&action=edit&redlink=1)
&#x20;[Jetnet](https://zh.wikipedia.org/w/index.php?title=Jetnet&action=edit&redlink=1)
&#x20;protocol communicates on | 非官方 |
| 6566/tcp | [SANE](https://zh.wikipedia.org/w/index.php?title=Scanner_Access_Now_Easy&action=edit&redlink=1)
&#x20;(Scanner Access Now Easy) - SANE network scanner daemon | 非官方 |
| 6600/tcp | [Music Playing Daemon (MPD)](https://zh.wikipedia.org/w/index.php?title=Music_Player_Daemon&action=edit&redlink=1) | 非官方 |
| 6619/tcp,udp | ODETTE-FTP over TLS/SSL | 官方 |
| 6665-6669/tcp | [IRC](https://zh.wikipedia.org/wiki/IRC) | 官方 |
| 6679/tcp | [IRC](https://zh.wikipedia.org/wiki/IRC)
&#x20;SSL （安全互联网中继聊天） - 通常使用的端口 | 非官方 |
| 6697/tcp | [IRC](https://zh.wikipedia.org/wiki/IRC)
&#x20;SSL （安全互联网中继聊天） - 通常使用的端口 | 非官方 |
| 6699/tcp | [WinMX](https://zh.wikipedia.org/wiki/WinMX)
&#x20;（参见 6257 端口） | 非官方 |
| 6881-6999/tcp,udp | [BitTorrent](<https://zh.wikipedia.org/wiki/BitTorrent_(%E5%8D%8F%E8%AE%AE)>)
&#x20;通常使用的端口 | 非官方 |
| 6891-6900/tcp,udp | [Windows Live Messenger](https://zh.wikipedia.org/wiki/Windows_Live_Messenger)
&#x20;（文件传输） | 官方 |
| 6901/tcp,udp | [Windows Live Messenger](https://zh.wikipedia.org/wiki/Windows_Live_Messenger)
&#x20;（语音） | 官方 |
| 6969/tcp | [acmsoda](https://zh.wikipedia.org/w/index.php?title=Acmsoda&action=edit&redlink=1) | 官方 |
| 6969/tcp | [BitTorrent](<https://zh.wikipedia.org/wiki/BitTorrent_(%E5%8D%8F%E8%AE%AE)>)
&#x20;tracker port | 非官方 |
| 7000/tcp | Default port for [Azureus](https://zh.wikipedia.org/wiki/Azureus)
's built in [超文本传输安全协议](https://zh.wikipedia.org/wiki/%E8%B6%85%E6%96%87%E6%9C%AC%E4%BC%A0%E8%BE%93%E5%AE%89%E5%85%A8%E5%8D%8F%E8%AE%AE)
&#x20;[BitTorrent tracker](https://zh.wikipedia.org/wiki/BitTorrent_tracker) | 非官方 |
| 7001/tcp | Default port for [BEA](https://zh.wikipedia.org/wiki/BEA_Systems)
&#x20;[WebLogic Server](https://zh.wikipedia.org/wiki/WebLogic)
's [超文本传输协议](https://zh.wikipedia.org/wiki/%E8%B6%85%E6%96%87%E6%9C%AC%E4%BC%A0%E8%BE%93%E5%8D%8F%E8%AE%AE)
&#x20;server - though often changed during installation | 非官方 |
| 7002/tcp | Default port for [BEA](https://zh.wikipedia.org/wiki/BEA_Systems)
&#x20;[WebLogic Server](https://zh.wikipedia.org/wiki/WebLogic)
's [超文本传输安全协议](https://zh.wikipedia.org/wiki/%E8%B6%85%E6%96%87%E6%9C%AC%E4%BC%A0%E8%BE%93%E5%AE%89%E5%85%A8%E5%8D%8F%E8%AE%AE)
&#x20;server - though often changed during installation | 非官方 |
| 7005/tcp,udp | Default port for [BMC 软件公司](https://zh.wikipedia.org/wiki/BMC%E8%BB%9F%E4%BB%B6%E5%85%AC%E5%8F%B8)
&#x20;CONTROL-M/Server and CONTROL-M/Agent's - Agent to Server port though often changed during installation | 非官方 |
| 7006/tcp,udp | Default port for [BMC 软件公司](https://zh.wikipedia.org/wiki/BMC%E8%BB%9F%E4%BB%B6%E5%85%AC%E5%8F%B8)
&#x20;CONTROL-M/Server and CONTROL-M/Agent's - Server to Agent port though often changed during installation | 非官方 |
| 7010/tcp | Default port for Cisco AON AMC (AON Management Console) [\[4\]](https://zh.wikipedia.org/wiki/TCP/UDP%E7%AB%AF%E5%8F%A3%E5%88%97%E8%A1%A8#cite_note-4) | 非官方 |
| 7025/tcp | Zimbra - lmtp \[mailbox] - local mail delivery | 非官方 |
| 7047/tcp | Zimbra - conversion server | 非官方 |
| 7171/tcp | [Tibia](<https://zh.wikipedia.org/w/index.php?title=Tibia_(computer_game)&action=edit&redlink=1>) | |
| 7306/tcp | Zimbra - mysql \[mailbox] | 非官方 |
| 7307/tcp | Zimbra - mysql \[logger] - logger | 非官方 |
| 7312/udp | [Sibelius](https://zh.wikipedia.org/w/index.php?title=Sibelius_notation_program&action=edit&redlink=1)
&#x20;License Server port | 非官方 |
| 7670/tcp | [BrettSpielWelt](https://zh.wikipedia.org/wiki/BrettSpielWelt)
&#x20;BSW Boardgame Portal | 非官方 |
| 7680/tcp | 适用于[Windows 10](https://zh.wikipedia.org/wiki/Windows_10)
更新的[传递优化](https://zh.wikipedia.org/w/index.php?title=%E5%82%B3%E9%81%9E%E6%9C%80%E4%BD%B3%E5%8C%96&action=edit&redlink=1) | 官方 |
| 7777/tcp | Default port used by Windows backdoor program tini.exe | 非官方 |
| 8000/tcp | [iRDMI](https://zh.wikipedia.org/w/index.php?title=IRDMI&action=edit&redlink=1)
&#x20;\- often mistakenly used instead of port 8080 (The Internet Assigned Numbers Authority (iana.org) officially lists this port for iRDMI protocol) | 官方 |
| 8000/tcp | Common port used for internet radio streams such as those using [SHOUTcast](https://zh.wikipedia.org/wiki/SHOUTcast) | 非官方 |
| 8002/tcp | Cisco Systems Unified Call Manager Intercluster Port | |
| 8008/tcp | [超文本传输协议](https://zh.wikipedia.org/wiki/%E8%B6%85%E6%96%87%E6%9C%AC%E4%BC%A0%E8%BE%93%E5%8D%8F%E8%AE%AE)
&#x20;替代端口 | 官方 |
| 8008/tcp | [IBM HTTP Server](https://zh.wikipedia.org/w/index.php?title=IBM_HTTP_Server&action=edit&redlink=1)
&#x20;默认管理端口 | 非官方 |
| 8009/tcp | [阿帕契族](https://zh.wikipedia.org/wiki/%E9%98%BF%E5%B8%95%E5%A5%91%E6%97%8F)
&#x20;[JServ](https://zh.wikipedia.org/w/index.php?title=JServ&action=edit&redlink=1)
&#x20;协议 v13 (ajp13) 例如：Apache mod*jk [Tomcat](https://zh.wikipedia.org/wiki/Tomcat)
会使用。 | 非官方 |
| 8010/tcp | [XMPP/Jabber](https://zh.wikipedia.org/wiki/Jabber)
&#x20;文件传输 | 非官方 |
| 8074/tcp | [Gadu-Gadu](https://zh.wikipedia.org/w/index.php?title=Gadu-Gadu&action=edit&redlink=1) | 非官方 |
| 8080/tcp | [超文本传输协议](https://zh.wikipedia.org/wiki/%E8%B6%85%E6%96%87%E6%9C%AC%E4%BC%A0%E8%BE%93%E5%8D%8F%E8%AE%AE)
&#x20;替代端口 （http_alt） - commonly used for [代理服务器](https://zh.wikipedia.org/wiki/%E4%BB%A3%E7%90%86%E6%9C%8D%E5%8A%A1%E5%99%A8)
&#x20;and [caching](https://zh.wikipedia.org/w/index.php?title=Caching&action=edit&redlink=1)
&#x20;server, or for running a web server as a non-[Root](https://zh.wikipedia.org/wiki/Root)
&#x20;user | 官方 |
| 8080/tcp | [Apache Tomcat](https://zh.wikipedia.org/wiki/Apache_Tomcat) | 非官方 |
| 8086/tcp | [HELM](https://zh.wikipedia.org/w/index.php?title=HELM&action=edit&redlink=1)
&#x20;Web Host Automation Windows Control Panel | 非官方 |
| 8086/tcp | [Kaspersky](https://zh.wikipedia.org/wiki/Kaspersky)
&#x20;AV Control Center TCP Port | 非官方 |
| 8087/tcp | [Hosting Accelerator](https://zh.wikipedia.org/w/index.php?title=Hosting_Accelerator&action=edit&redlink=1)
&#x20;Control Panel | 非官方 |
| 8087/udp | [Kaspersky](https://zh.wikipedia.org/wiki/Kaspersky)
&#x20;AV Control Center UDP Port | 非官方 |
| 8087/tcp | [英迈](https://zh.wikipedia.org/wiki/%E8%8B%B1%E9%82%81)
&#x20;控制面板 | 非官方 |
| 8090/tcp | Another [超文本传输协议](https://zh.wikipedia.org/wiki/%E8%B6%85%E6%96%87%E6%9C%AC%E4%BC%A0%E8%BE%93%E5%8D%8F%E8%AE%AE)
&#x20;Alternate (http_alt_alt) - used as an alternative to port 8080 | 非官方 |
| 8118/tcp | [Privoxy](https://zh.wikipedia.org/wiki/Privoxy)
&#x20;web proxy - advertisements-filtering web proxy | 官方 |
| 8123/tcp | Dynmap[\[5\]](https://zh.wikipedia.org/wiki/TCP/UDP%E7%AB%AF%E5%8F%A3%E5%88%97%E8%A1%A8#cite_note-5)
默认网页端口号(Minecraft 在线地图) | 非官方 |
| 8200/tcp | [GoToMyPC](https://zh.wikipedia.org/w/index.php?title=GoToMyPC&action=edit&redlink=1) | 非官方 |
| 8220/tcp | [Bloomberg](https://zh.wikipedia.org/w/index.php?title=Bloomberg_Terminal&action=edit&redlink=1) | 非官方 |
| 8222 | [VMware](https://zh.wikipedia.org/wiki/VMware)
服务器管理用户界面(不安全网络界面)[\[6\]](https://zh.wikipedia.org/wiki/TCP/UDP%E7%AB%AF%E5%8F%A3%E5%88%97%E8%A1%A8#cite_note-6)
。参见 8333 端口 | 非官方 |
| 8291/tcp | Winbox - Default port on a MikroTik RouterOS for a Windows application used to administer MikroTik RouterOS | 非官方 |
| 8294/tcp | [Bloomberg](https://zh.wikipedia.org/w/index.php?title=Bloomberg_Terminal&action=edit&redlink=1) | 非官方 |
| 8333 | [VMware](https://zh.wikipedia.org/wiki/VMware)
&#x20;服务器管理用户界面（安全网络界面）[\[7\]](https://zh.wikipedia.org/wiki/TCP/UDP%E7%AB%AF%E5%8F%A3%E5%88%97%E8%A1%A8#cite_note-7)
。参见 8222 端口 | 非官方 |
| 8400 | [Commvault](https://zh.wikipedia.org/w/index.php?title=Commvault&action=edit&redlink=1)
&#x20;Unified Data Management | 官方 |
| 8443/tcp | [英迈](https://zh.wikipedia.org/wiki/%E8%8B%B1%E9%82%81)
&#x20;Control Panel | 非官方 |
| 8500/tcp | [Adobe ColdFusion](https://zh.wikipedia.org/wiki/Adobe_ColdFusion)
&#x20;Macromedia/Adobe ColdFusion default Webserver port | 非官方 |
| 8501/udp | [毁灭公爵 3D](https://zh.wikipedia.org/wiki/%E6%AF%81%E7%81%AD%E5%85%AC%E7%88%B53D)
&#x20;\- Default Online Port | 官方 |
| 8767/udp | [TeamSpeak](https://zh.wikipedia.org/wiki/TeamSpeak)
&#x20;\- Default UDP Port | 非官方 |
| 8880 | [IBM WebSphere Application Server](https://zh.wikipedia.org/wiki/IBM_WebSphere_Application_Server)
&#x20;[简单对象访问协议](https://zh.wikipedia.org/wiki/%E7%AE%80%E5%8D%95%E5%AF%B9%E8%B1%A1%E8%AE%BF%E9%97%AE%E5%8D%8F%E8%AE%AE)
&#x20;Connector port | |
| 8881/tcp | [Atlasz Informatics Research Ltd](https://zh.wikipedia.org/w/index.php?title=Atlasz_Informatics_Research_Ltd&action=edit&redlink=1)
&#x20;Secure Application Server | 非官方 |
| 8882/tcp | [Atlasz Informatics Research Ltd](https://zh.wikipedia.org/w/index.php?title=Atlasz_Informatics_Research_Ltd&action=edit&redlink=1)
&#x20;Secure Application Server | 非官方 |
| 8888/tcp,udp | [NewsEDGE](https://zh.wikipedia.org/w/index.php?title=NewsEDGE&action=edit&redlink=1)
&#x20;server | 官方 |
| 8888/tcp | [Sun Answerbook](https://zh.wikipedia.org/w/index.php?title=Sun_Answerbook&action=edit&redlink=1)
&#x20;[网页服务器](https://zh.wikipedia.org/wiki/%E7%B6%B2%E9%A0%81%E4%BC%BA%E6%9C%8D%E5%99%A8)
&#x20;server (deprecated by [docs.sun.com](https://web.archive.org/web/20080222182837/http://docs.sun.com/app/docs)
) | 非官方 |
| 8888/tcp | [GNUmp3d](https://zh.wikipedia.org/w/index.php?title=GNUmp3d&action=edit&redlink=1)
&#x20;HTTP music streaming and web interface port | 非官方 |
| 8888/tcp | [英雄大作战](https://zh.wikipedia.org/wiki/%E8%8B%B1%E9%9B%84%E5%A4%A7%E4%BD%9C%E6%88%B0)
&#x20;Network Game Server | 非官方 |
| 9000/tcp | Buffalo LinkSystem web access | 非官方 |
| 9000/tcp | [DBGp](https://zh.wikipedia.org/w/index.php?title=DBGp&action=edit&redlink=1) | 非官方 |
| 9000/udp | [UDPCast](https://zh.wikipedia.org/w/index.php?title=UDPCast&action=edit&redlink=1) | 非官方 |
| 9000 | [PHP](https://zh.wikipedia.org/wiki/PHP)
&#x20;\- php-fpm | 非官方 |
| 9001 | cisco-xremote router configuration | 非官方 |
| 9001 | [Tor](https://zh.wikipedia.org/wiki/Tor)
&#x20;network default port | 非官方 |
| 9001/tcp | [DBGp](https://zh.wikipedia.org/w/index.php?title=DBGp&action=edit&redlink=1)
&#x20;Proxy | 非官方 |
| 9002 | Default [ElasticSearch](https://zh.wikipedia.org/w/index.php?title=ElasticSearch&action=edit&redlink=1)
&#x20;port | |
| 9009/tcp,udp | [Pichat Server](https://zh.wikipedia.org/w/index.php?title=Pichat&action=edit&redlink=1)
&#x20;\- Peer to peer chat software | 官方 |
| 9043/tcp | [IBM WebSphere Application Server](https://zh.wikipedia.org/wiki/IBM_WebSphere_Application_Server)
&#x20;Administration Console secure port | |
| 9060/tcp | [IBM WebSphere Application Server](https://zh.wikipedia.org/wiki/IBM_WebSphere_Application_Server)
&#x20;Administration Console | |
| 9100/tcp | [Jetdirect](https://zh.wikipedia.org/w/index.php?title=Jetdirect&action=edit&redlink=1)
&#x20;HP Print Services | 官方 |
| 9110/udp | [SSMP](https://zh.wikipedia.org/w/index.php?title=SSMP&action=edit&redlink=1)
&#x20;Message protocol | 非官方 |
| 9101 | [Bacula](https://zh.wikipedia.org/w/index.php?title=Bacula&action=edit&redlink=1)
&#x20;Director | 官方 |
| 9102 | [Bacula](https://zh.wikipedia.org/w/index.php?title=Bacula&action=edit&redlink=1)
&#x20;File Daemon | 官方 |
| 9103 | [Bacula](https://zh.wikipedia.org/w/index.php?title=Bacula&action=edit&redlink=1)
&#x20;Storage Daemon | 官方 |
| 9119/TCP,UDP | [Mxit](https://zh.wikipedia.org/wiki/Mxit)
&#x20;Instant Messenger | 官方 |
| 9535/tcp | man, Remote Man Server | |
| 9535 | mngsuite - Management Suite Remote Control | 官方 |
| 9800/tcp,udp | [基于 Web 的分布式编写和版本控制](https://zh.wikipedia.org/wiki/%E5%9F%BA%E4%BA%8EWeb%E7%9A%84%E5%88%86%E5%B8%83%E5%BC%8F%E7%BC%96%E5%86%99%E5%92%8C%E7%89%88%E6%9C%AC%E6%8E%A7%E5%88%B6)
&#x20;Source Port | 官方 |
| 9800 | [WebCT](https://zh.wikipedia.org/wiki/WebCT)
&#x20;e-learning portal | 非官方 |
| 9999 | [Hydranode](https://zh.wikipedia.org/w/index.php?title=Hydranode&action=edit&redlink=1)
&#x20;\- edonkey2000 telnet control port | 非官方 |
| 9999 | Urchin Web Analytics | 非官方 |
| 10000 | [Webmin](https://zh.wikipedia.org/w/index.php?title=Webmin&action=edit&redlink=1)
&#x20;\- web based Linux admin tool | 非官方 |
| 10000 | [BackupExec](https://zh.wikipedia.org/w/index.php?title=BackupExec&action=edit&redlink=1) | 非官方 |
| 10008 | Octopus Multiplexer - CROMP protocol primary port, hoople.org | 官方 |
| 10017 | AIX,NeXT, HPUX - rexd daemon control port | 非官方 |
| 10024/tcp | Zimbra - smtp \[mta] - to amavis from postfix | 非官方 |
| 10025/tcp | Ximbra - smtp \[mta] - back to postfix from amavis | 非官方 |
| 10050/tcp | Zabbix-Agent | |
| 10051/tcp | Zabbix-Server | |
| 10113/tcp,udp | NetIQ Endpoint | 官方 |
| 10114/tcp,udp | NetIQ Qcheck | 官方 |
| 10115/tcp,udp | NetIQ Endpoint | 官方 |
| 10116/tcp,udp | NetIQ VoIP Assessor | 官方 |
| 10480 | SWAT 4 Dedicated Server | 非官方 |
| 11211 | [Memcached](https://zh.wikipedia.org/wiki/Memcached) | 非官方 |
| 11235 | Savage:Battle for Newerth Server Hosting | 非官方 |
| 11294 | Blood Quest Online Server | 非官方 |
| 11371 | [PGP](https://zh.wikipedia.org/wiki/PGP)
&#x20;HTTP Keyserver | 官方 |
| 11576 | [IPStor](https://zh.wikipedia.org/w/index.php?title=IPStor&action=edit&redlink=1)
&#x20;Server management communication | 非官方 |
| 12035/udp | 《[第二人生](https://zh.wikipedia.org/wiki/%E7%AC%AC%E4%BA%8C%E4%BA%BA%E7%94%9F*(%E4%BA%92%E8%81%AF%E7%B6%B2))
》, used for server UDP in-bound[\[8\]](https://zh.wikipedia.org/wiki/TCP/UDP%E7%AB%AF%E5%8F%A3%E5%88%97%E8%A1%A8#cite_note-http://wiki.secondlife.com/wiki/Authentication_Flow@step_4-8) | 非官方 |
| 12345 | [NetBus](https://zh.wikipedia.org/wiki/NetBus)
&#x20;\- remote administration tool (often [特洛伊木马 (电脑)](<https://zh.wikipedia.org/wiki/%E7%89%B9%E6%B4%9B%E4%BC%8A%E6%9C%A8%E9%A9%AC_(%E7%94%B5%E8%84%91)>)
). Also used by [NetBuster](https://zh.wikipedia.org/w/index.php?title=NetBuster&action=edit&redlink=1)
. Little Fighter 2 (TCP). | 非官方 |
| 12975/tcp | LogMeIn [Hamachi](https://zh.wikipedia.org/w/index.php?title=Hamachi&action=edit&redlink=1)
&#x20;(VPN tunnel software;also port 32976) | |
| 13000-13050/udp | 《[第二人生](<https://zh.wikipedia.org/wiki/%E7%AC%AC%E4%BA%8C%E4%BA%BA%E7%94%9F_(%E4%BA%92%E8%81%AF%E7%B6%B2)>)
》, used for server UDP in-bound[\[9\]](https://zh.wikipedia.org/wiki/TCP/UDP%E7%AB%AF%E5%8F%A3%E5%88%97%E8%A1%A8#cite_note-9) | 非官方 |
| 13720/tcp | [赛门铁克](https://zh.wikipedia.org/wiki/%E8%B5%9B%E9%97%A8%E9%93%81%E5%85%8B)
&#x20;[NetBackup](https://zh.wikipedia.org/w/index.php?title=NetBackup&action=edit&redlink=1)
&#x20;\- bprd (formerly [VERITAS](https://zh.wikipedia.org/w/index.php?title=Veritas_Software&action=edit&redlink=1)
) | |
| 13721/tcp | [赛门铁克](https://zh.wikipedia.org/wiki/%E8%B5%9B%E9%97%A8%E9%93%81%E5%85%8B)
&#x20;[NetBackup](https://zh.wikipedia.org/w/index.php?title=NetBackup&action=edit&redlink=1)
&#x20;\- bpdbm (formerly [VERITAS](https://zh.wikipedia.org/w/index.php?title=Veritas_Software&action=edit&redlink=1)
) | |
| 13724/tcp | [赛门铁克](https://zh.wikipedia.org/wiki/%E8%B5%9B%E9%97%A8%E9%93%81%E5%85%8B)
&#x20;Network Utility - vnet (formerly [VERITAS](https://zh.wikipedia.org/w/index.php?title=Veritas_Software&action=edit&redlink=1)
) | |
| 13782/tcp | [赛门铁克](https://zh.wikipedia.org/wiki/%E8%B5%9B%E9%97%A8%E9%93%81%E5%85%8B)
&#x20;[NetBackup](https://zh.wikipedia.org/w/index.php?title=NetBackup&action=edit&redlink=1)
&#x20;\- bpcd (formerly [VERITAS](https://zh.wikipedia.org/w/index.php?title=Veritas_Software&action=edit&redlink=1)
) | |
| 13783/tcp | [赛门铁克](https://zh.wikipedia.org/wiki/%E8%B5%9B%E9%97%A8%E9%93%81%E5%85%8B)
&#x20;VOPIED protocol (formerly [VERITAS](https://zh.wikipedia.org/w/index.php?title=Veritas_Software&action=edit&redlink=1)
) | |
| 14567/udp | [战地风云 1942](https://zh.wikipedia.org/wiki/%E6%88%B0%E5%9C%B0%E9%A2%A8%E9%9B%B21942)
&#x20;and mods | 非官方 |
| 15000/tcp | [Bounce (网络)](<https://zh.wikipedia.org/wiki/Bounce_(%E7%BD%91%E7%BB%9C)>) | 非官方 |
| 15000/tcp | [韦诺之战](https://zh.wikipedia.org/wiki/%E9%9F%A6%E8%AF%BA%E4%B9%8B%E6%88%98) | |
| 15567/udp | [战地风云：越南](https://zh.wikipedia.org/wiki/%E6%88%B0%E5%9C%B0%E9%A2%A8%E9%9B%B2%EF%BC%9A%E8%B6%8A%E5%8D%97)
&#x20;and mods | 非官方 |
| 15345/udp | [XPilot](https://zh.wikipedia.org/w/index.php?title=XPilot&action=edit&redlink=1) | 官方 |
| 16000/tcp | [Bounce (网络)](<https://zh.wikipedia.org/wiki/Bounce_(%E7%BD%91%E7%BB%9C)>) | 非官方 |
| 16080/tcp | [MacOS Server](https://zh.wikipedia.org/wiki/MacOS_Server)
&#x20;performance cache for [超文本传输协议](https://zh.wikipedia.org/wiki/%E8%B6%85%E6%96%87%E6%9C%AC%E4%BC%A0%E8%BE%93%E5%8D%8F%E8%AE%AE)
[\[10\]](https://zh.wikipedia.org/wiki/TCP/UDP%E7%AB%AF%E5%8F%A3%E5%88%97%E8%A1%A8#cite_note-10) | 非官方 |
| 16384/udp | Iron Mountain Digital - online backup | 非官方 |
| 16567/udp | [战地 2](https://zh.wikipedia.org/wiki/%E6%88%98%E5%9C%B02)
&#x20;and mods | 非官方 |
| 17788/udp | [PPS 网络电视](https://zh.wikipedia.org/wiki/PPS%E7%B6%B2%E8%B7%AF%E9%9B%BB%E8%A6%96) | 非官方 |
| 19132/udp | [Minecraft 基岩版](https://zh.wikipedia.org/wiki/%E6%88%91%E7%9A%84%E4%B8%96%E7%95%8C_%E5%9F%BA%E5%B2%A9%E7%89%88)
默认服务器端口号 | 非官方 |
| 19226/tcp | [熊猫 (消歧义)](<https://zh.wikipedia.org/wiki/%E7%86%8A%E7%8C%AB_(%E6%B6%88%E6%AD%A7%E4%B9%89)>)
&#x20;AdminSecure Communication Agent | 非官方 |
| 19638/tcp | Ensim Control Panel | 非官方 |
| 19813/tcp | 4D database Client Server Communication | 非官方 |
| 20000 | [Usermin](https://zh.wikipedia.org/w/index.php?title=Usermin&action=edit&redlink=1)
&#x20;\- 基于网络的用户工具 | 官方 |
| 20720/tcp | [Symantec i3](https://zh.wikipedia.org/w/index.php?title=Symantec_i3&action=edit&redlink=1)
&#x20;Web GUI server | 非官方 |
| 22347/tcp,udp | WibuKey - default port for WibuKey Network Server of WIBU-SYSTEMS AG | 官方 |
| 22350/tcp,udp | CodeMeter - default port for CodeMeter Server of WIBU-SYSTEMS AG | 官方 |
| 24554/tcp,udp | [binkp](https://zh.wikipedia.org/w/index.php?title=Binkp&action=edit&redlink=1)
&#x20;\- [FidoNet](https://zh.wikipedia.org/wiki/FidoNet)
&#x20;mail transfers over [TCP/IP 协议族](https://zh.wikipedia.org/wiki/TCP/IP%E5%8D%8F%E8%AE%AE%E6%97%8F) | 官方 |
| 24800 | [Synergy](https://zh.wikipedia.org/wiki/Synergy)
：keyboard/mouse sharing software | 非官方 |
| 24842 | [StepMania：Online](https://zh.wikipedia.org/wiki/StepMania)
：《[劲爆热舞](https://zh.wikipedia.org/wiki/%E5%8B%81%E7%88%86%E7%86%B1%E8%88%9E)
》模拟器 | 非官方 |
| 25565/tcp | [Minecraft](https://zh.wikipedia.org/wiki/%E6%88%91%E7%9A%84%E4%B8%96%E7%95%8C)
默认服务器端口号 | 非官方 |
| 25999/tcp | [Xfire](https://zh.wikipedia.org/w/index.php?title=Xfire&action=edit&redlink=1) | 非官方 |
| 26000/tcp,udp | [Id Software](https://zh.wikipedia.org/wiki/Id_Software)
's 《[Quake](https://zh.wikipedia.org/wiki/Quake)
》 server, | 官方 |
| 26000/tcp | [CCP Games](https://zh.wikipedia.org/wiki/CCP_Games)
's [星战前夜](https://zh.wikipedia.org/wiki/%E6%98%9F%E6%88%98%E5%89%8D%E5%A4%9C)
&#x20;Online gaming MMORPG, | 非官方 |
| 27000/udp | (through 27006) [Id Software](https://zh.wikipedia.org/wiki/Id_Software)
's 《[雷神世界](https://zh.wikipedia.org/wiki/%E9%9B%B7%E7%A5%9E%E4%B8%96%E7%95%8C)
》 master server | 非官方 |
| 27010/udp | [Half-Life](https://zh.wikipedia.org/wiki/Half-Life)
及其修改版，如《[反恐精英系列](https://zh.wikipedia.org/wiki/%E5%8F%8D%E6%81%90%E7%B2%BE%E8%8B%B1%E7%B3%BB%E5%88%97)
》 | 非官方 |
| 27015/udp | [Half-Life](https://zh.wikipedia.org/wiki/Half-Life)
及其修改版，如《[反恐精英系列](https://zh.wikipedia.org/wiki/%E5%8F%8D%E6%81%90%E7%B2%BE%E8%8B%B1%E7%B3%BB%E5%88%97)
》 | 非官方 |
| 27017/tcp | MongoDB 数据库 | 非官方 |
| 27374 | [Sub7](https://zh.wikipedia.org/w/index.php?title=Sub7&action=edit&redlink=1)
's default port. Most [脚本小子](https://zh.wikipedia.org/wiki/%E8%84%9A%E6%9C%AC%E5%B0%8F%E5%AD%90)
s do not change the default port. | 非官方 |
| 27500/udp | (through 27900) [Id Software](https://zh.wikipedia.org/wiki/Id_Software)
's 《[雷神世界](https://zh.wikipedia.org/wiki/%E9%9B%B7%E7%A5%9E%E4%B8%96%E7%95%8C)
》 | 非官方 |
| 27888/udp | [Kaillera](https://zh.wikipedia.org/w/index.php?title=Kaillera&action=edit&redlink=1)
&#x20;server | 非官方 |
| 27900 | (through 27901) [任天堂](https://zh.wikipedia.org/wiki/%E4%BB%BB%E5%A4%A9%E5%A0%82)
&#x20;[任天堂 Wi-Fi 连接](https://zh.wikipedia.org/wiki/%E4%BB%BB%E5%A4%A9%E5%A0%82Wi-Fi%E8%BF%9E%E6%8E%A5) | 非官方 |
| 27901/udp | (through 27910) [Id Software](https://zh.wikipedia.org/wiki/Id_Software)
's 《[雷神之锤 II](https://zh.wikipedia.org/wiki/%E9%9B%B7%E7%A5%9E%E4%B9%8B%E9%94%A4II)
》 master server | 非官方 |
| 27960/udp | (through 27969) [动视](https://zh.wikipedia.org/wiki/%E5%8A%A8%E8%A7%86)
's 《[Enemy Territory](https://zh.wikipedia.org/w/index.php?title=Enemy_Territory&action=edit&redlink=1)
》 and [Id Software](https://zh.wikipedia.org/wiki/Id_Software)
's 《[雷神之锤 III 竞技场](https://zh.wikipedia.org/wiki/%E9%9B%B7%E7%A5%9E%E4%B9%8B%E9%94%A4III%E7%AB%9E%E6%8A%80%E5%9C%BA)
》 and 《Quake III》 and some ioquake3 derived games | 非官方 |
| 28910 | [任天堂](https://zh.wikipedia.org/wiki/%E4%BB%BB%E5%A4%A9%E5%A0%82)
&#x20;[任天堂 Wi-Fi 连接](https://zh.wikipedia.org/wiki/%E4%BB%BB%E5%A4%A9%E5%A0%82Wi-Fi%E8%BF%9E%E6%8E%A5) | 非官方 |
| 28960 | [决胜时刻 2](https://zh.wikipedia.org/wiki/%E6%B1%BA%E5%8B%9D%E6%99%82%E5%88%BB2)
&#x20;Common Call of Duty 2 port - (PC Version) | 非官方 |
| 29900 | (through 29901) [任天堂](https://zh.wikipedia.org/wiki/%E4%BB%BB%E5%A4%A9%E5%A0%82)
&#x20;[任天堂 Wi-Fi 连接](https://zh.wikipedia.org/wiki/%E4%BB%BB%E5%A4%A9%E5%A0%82Wi-Fi%E8%BF%9E%E6%8E%A5) | 非官方 |
| 29920 | [任天堂](https://zh.wikipedia.org/wiki/%E4%BB%BB%E5%A4%A9%E5%A0%82)
&#x20;[任天堂 Wi-Fi 连接](https://zh.wikipedia.org/wiki/%E4%BB%BB%E5%A4%A9%E5%A0%82Wi-Fi%E8%BF%9E%E6%8E%A5) | 非官方 |
| 30000 | [Pokemon Netbattle](https://zh.wikipedia.org/w/index.php?title=Pokemon_Netbattle&action=edit&redlink=1) | 非官方 |
| 30564/tcp | [Multiplicity](<https://zh.wikipedia.org/w/index.php?title=Multiplicity_(software)&action=edit&redlink=1>)
：keyboard/mouse/clipboard sharing software | 非官方 |
| 31337/tcp | [Back Orifice](https://zh.wikipedia.org/w/index.php?title=Back_Orifice&action=edit&redlink=1)
&#x20;\- remote administration tool（often [特洛伊木马 (电脑)](<https://zh.wikipedia.org/wiki/%E7%89%B9%E6%B4%9B%E4%BC%8A%E6%9C%A8%E9%A9%AC_(%E7%94%B5%E8%84%91)>)
） | 非官方 |
| 31337/tcp | xc0r3 - xc0r3 security antivir port | 非官方 |
| 31415 | [ThoughtSignal](https://zh.wikipedia.org/w/index.php?title=ThoughtSignal&action=edit&redlink=1)
&#x20;\- Server Communication Service（often [Informational](<https://zh.wikipedia.org/w/index.php?title=Information_(computing)&action=edit&redlink=1>)
） | 非官方 |
| 31456-31458/tcp | [TetriNET](https://zh.wikipedia.org/w/index.php?title=TetriNET&action=edit&redlink=1)
&#x20;ports (in order: IRC, game, and spectating) | 非官方 |
| 32245/tcp | [MMTSG-mutualed](https://zh.wikipedia.org/w/index.php?title=True_Mutualization_Service&action=edit&redlink=1)
&#x20;over [MMT](https://zh.wikipedia.org/w/index.php?title=True_Mutualization_Service&action=edit&redlink=1)
&#x20;(encrypted transmission) | 非官方 |
| 33434 | [Traceroute](https://zh.wikipedia.org/wiki/Traceroute) | 官方 |
| 37777/tcp | [Digital Video Recorder hardware](https://zh.wikipedia.org/w/index.php?title=Digital_Video_Recorder_hardware&action=edit&redlink=1) | 非官方 |
| 36963 | [Counter Strike 2D](https://zh.wikipedia.org/w/index.php?title=Counter_Strike_2D&action=edit&redlink=1)
&#x20;multiplayer port (2D clone of popular CounterStrike computer game) | 非官方 |
| 40000 | SafetyNET p | 官方 |
| 43594-43595/tcp | [RuneScape](https://zh.wikipedia.org/wiki/RuneScape) | 非官方 |
| 47808 | [BACnet](https://zh.wikipedia.org/wiki/BACnet)
&#x20;Building Automation and Control Networks | 官方 |

## 49152 到 65535 号端口

参见：[临时端口](https://zh.wikipedia.org/wiki/%E4%B8%B4%E6%97%B6%E7%AB%AF%E5%8F%A3)
根据定义，该段端口属于“动态端口”范围，没有端口可以被正式地注册占用。 [\[11\]](https://zh.wikipedia.org/wiki/TCP/UDP%E7%AB%AF%E5%8F%A3%E5%88%97%E8%A1%A8#cite_note-11)
