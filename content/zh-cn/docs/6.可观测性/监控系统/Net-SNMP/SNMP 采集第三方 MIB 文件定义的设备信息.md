---
title: SNMP 采集第三方 MIB 文件定义的设备信息
---

<https://lvii.github.io/server/2013-07-05-net-snmp-get-info-from-third-custumed-mib-file/>

厂里机房搞来 PDU 机柜电源插座，这个 bootbox 可以采集一些机柜温度，以及各个插座的电流功率啥的。

要用 SNMP 协议来获取这些信息，因为是第三方设备，要根据厂商的 MIB 文件中定义的 OID 来解析相关字段。

下面开始介绍怎么折腾：

- 把 MIB 放到系统 `mibdirs` 路径
- 让 `net-snmp` 工具自动加载 MIB 文件
- 用 `snmptranslate` 解析 MIB 文件中的 OID
- 用 `snmpwalk` 或 `snmpget` 采集对应 OID 字段信息

准备一下 snmp 相关的命令行工具，要在 CentOS 装上下面几个软件包：

|       rpm        |                       command                        |
|:----------------:|:----------------------------------------------------:|
|    `net-snmp`    |                     `snmpd` 服务                     |
| `net-snmp-libs`  |                系统常用的 MIBS 文件库                |
| `net-snmp-utils` | 提供 `snmpget / snmpwalk / snmptranslate` 等采集工具 |
| `net-snmp-devel` |          提供 `net-snmp-config` 命令行工具           |

## 修改 MIB 文件路径 [Permalink](https://lvii.github.io/server/2013-07-05-net-snmp-get-info-from-third-custumed-mib-file/#%E4%BF%AE%E6%94%B9-mib-%E6%96%87%E4%BB%B6%E8%B7%AF%E5%BE%84)

参考 net-snmp 文档: [「Where should I put my MIB files?」](http://www.net-snmp.org/FAQ.html#Where_should_I_put_my_MIB_files_)
将 MIB 文件放到系统路径下，用 `net-snmp-config` 命令查询系统默认 `mibdirs` :

```bash
$ net-snmp-config --default-mibdirs
/home/user/.snmp/mibs:/usr/share/snmp/mibs
```
还可以用 `snmpget` 查找路径：

```bash
$ snmpget  -Dparse-mibs  2>&1 | grep directory
parse-mibs: Scanning directory /home/user/.snmp/mibs
parse-mibs: cannot open MIB directory /home/user/.snmp/mibs
parse-mibs: Scanning directory /usr/share/snmp/mibs
```

把第三方 MIB 库文件复制到 `/usr/share/snmp/mibs` 即可

## 自动加载 MIB 模块 [Permalink](https://lvii.github.io/server/2013-07-05-net-snmp-get-info-from-third-custumed-mib-file/#%E8%87%AA%E5%8A%A8%E5%8A%A0%E8%BD%BD-mib-%E6%A8%A1%E5%9D%97)

但默认 net-snmp 工具不会自动加载自定义的 MIB 文件，为了让 `snmpwalk / snmpget ...` 工具能够 **自动加载**
第三方 MIB 文件，需要修改 `snmp.conf` 配置：

```bash
$ cat /etc/snmp/snmp.conf
mibs +NPM3G-MIB
```

**MIB 模块名称** : `NPM3G-MIB` 是 MIB 文件中 `DEFINITIONS` 字段前面的名字：

```bash
$ grep DEFINITIONS ~/NPM3G.mib
NPM3G-MIB DEFINITIONS ::= BEGIN
```

net-snmp 官网文档还提供 **环境变量** 和 **命令行** 方法加载 MIB 文件：

<http://www.net-snmp.org/FAQ.html#How_do_I_add_a_MIB_to_the_tools_>

```bash
Secondly, tell the tools to load this MIB:
snmpwalk -m +MY-MIB .....
(load it for this command only)
export MIBS=+MY-MIB
(load it for this session only)
echo "mibs +MY-MIB" >> $HOME/.snmp/snmp.conf
(load it every time)
```

## 解析 MIB 库中的 OID [Permalink](https://lvii.github.io/server/2013-07-05-net-snmp-get-info-from-third-custumed-mib-file/#%E8%A7%A3%E6%9E%90-mib-%E5%BA%93%E4%B8%AD%E7%9A%84-oid)

使用 `snmptranslate` 解析刚添加的 MIB 文件中定义的 OID 及其对应的文字描述：

```bash
$ snmptranslate -Tz -m NPM3G-MIB|column -t|head
"org"                       "1.3"
"dod"                       "1.3.6"
"internet"                  "1.3.6.1"
"directory"                 "1.3.6.1.1"
"mgmt"                      "1.3.6.1.2"
"mib-2"                     "1.3.6.1.2.1"
"transmission"              "1.3.6.1.2.1.10"
"experimental"              "1.3.6.1.3"
"private"                   "1.3.6.1.4"
"enterprises"               "1.3.6.1.4.1"
```

`snmptranslate` OID 文字描述和 OID 数字格式之间的相互转换：

```bash
$ snmptranslate 1.3.6.1.2.1.1.3
SNMPv2-MIB::sysUpTime
$ snmptranslate -On SNMPv2-MIB::sysUpTime
.1.3.6.1.2.1.1.3
```

`snmptranslate` 还可以查看 MIB 树的节点信息，比如 **数据类型** 、 **单位** …：

```bash
$ snmptranslate -Tp -OS -m SNMPv2-MIB|less -i
+--iso(1)
   |
   +--org(3)
      |
      +--dod(6)
         |
         +--internet(1)
            |
            +--directory(1)
            |
            +--mgmt(2)
            |  |
            |  +--mib-2(1)
            |     |
            |     +--system(1)
            |     |  |
            |     |  +-- -R-- String    sysDescr(1)
            |     |  |        Textual Convention: DisplayString
            |     |  |        Size: 0..255
            |     |  +-- -R-- ObjID     sysObjectID(2)
            |     |  +-- -R-- TimeTicks sysUpTime(3)
            |     |  +-- -RW- String    sysContact(4)
            |     |  |        Textual Convention: DisplayString
            |     |  |        Size: 0..255
```

## 采集信息 [Permalink](https://lvii.github.io/server/2013-07-05-net-snmp-get-info-from-third-custumed-mib-file/#%E9%87%87%E9%9B%86%E4%BF%A1%E6%81%AF)

上面只是在采集服务器上面，添加了第三方的 MIB 文件，使得 net-snmp 可以识别第三方厂商定义的 OID 监控项。 但是还需要开启要监控的设备的 SNMP 服务，让监控服务器和硬件设备进行通讯，采集信息。 在设备上面配置采集服务器的 IP 和 SNMP 协议版本和认证方式，重启设备。 在采集服务器上面查看 SNMP 端口 OK 之后，开始用 net-snmp 工具采集信息：

```bash
$ sudo nmap -sU 10.10.20.30
Starting Nmap 5.51 ( http://nmap.org ) at 2013-07-05 22:08 CST
Nmap scan report for 10.10.20.30
Host is up (0.022s latency).
Not shown: 998 closed ports
PORT    STATE         SERVICE
161/udp open|filtered snmp
```

验证 SNMP 通讯是否正常，正常的话，下面的命令会采集到设备名称、时间、IP 地址等通用信息

```bash
$ snmpwalk -v2c -c public 10.10.20.30
```

然后采集传感器的温湿度信息：

```bash
$ snmpwalk -v2c -c public 10.10.20.30 .1.3.6.1.4.1.30966.4.2.1.22
NPM3G-MIB::npm1Temperature1.0 = STRING: 18
$ snmpwalk -v2c -c public 10.10.20.30 NPM3G-MIB::npm1Humidity1
NPM3G-MIB::npm1Humidity1.0 = STRING: 76
```

后面就可以根据需要采集想要的信息了

## 参考 [Permalink](https://lvii.github.io/server/2013-07-05-net-snmp-get-info-from-third-custumed-mib-file/#%E5%8F%82%E8%80%83)

[net-snmp 中载入第三方 mib 库](http://fs20041242.iteye.com/blog/889041)

<http://www.net-snmp.org/FAQ.html>

[Using and loading MIBS](http://www.net-snmp.org/wiki/index.php/TUT:Using_and_loading_MIBS)
