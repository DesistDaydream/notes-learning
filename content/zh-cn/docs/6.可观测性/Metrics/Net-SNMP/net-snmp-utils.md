---
title: net-snmp-utils
linkTitle: net-snmp-utils
weight: 20
---

# 概述

> 参考：
>
> - [官方文档，手册](https://www.net-snmp.org/docs/man/)

# snmpwalk 与 snmpget

snmpwalk 与 snmpget 使用 snmp 协议的 GETNEXT 请求，向 SNMP 代理发送查询请求，以便获取 SNMP 数据。

- snmpget 获取指定 OID 的数据
- snmpwalk 可以获取大量 OID 的数据

## 关联文件与配置

**/etc/snmp/snmp.conf** # snmpwalk 与 snmpget 运行时配置文件。若不存在则手动创建

**/usr/local/share/snmp/mibs** # MIB 文件的默认路径。这里现阶段包含 66 个 MIB 文件。

`net-snmp-config --default-mibdirs` 命令可以列出工具在使用中会读取 MIB 文件的默认路径，包括如下几个：

- 注意，CentOS 和 Ubuntu 的路径可能不相同。PS: 这种老程序是真滴坑
- **$HOME/.snmp/mibs**
- **/usr/share/snmp/mibs**
- **/usr/share/snmp/mibs/iana**
- **/usr/share/snmp/mibs/ietf**

### 添加自定义 MIB

在 /etc/snmp/snmp.conf 文件中添加如下内容

```bash
mibdirs +/root/.snmp/mibs/h3c
```

在 /root/.snmp/mibs/h3c 目录下添加 MIB 文件后，snmpwalk 就可以获取到第三方 MIB 信息。

若是无法找到 MIB，则可能是版本过来，还需要在文件中添加如下内容，以手动指定要添加的 OID

```text
mibs +HH3C-OID-MIB
mibs +HH3C-SERVER-AGENT-MIB
```

## Syntax(语法)

> snmpget 语法与 snmpwalk 语法基本一致，只不过行为和结果有细微区别

**snmpwalk \[OPTIONS] AGENT \[OID]**

OPTIONS

- **-l \<noAuthNoPriv | authNoPriv | authPriv>** # 设置安全级别
- **-c \<STRING>** # 团体名
- **-v <1 | 2c | 3>** # snmp 版本
- 认证相关选项
  - **-u USERNAME** # 用户名
  - -**a PROTOCOL** # 指定认证的算法
  - **-A PASSWORD** # 指定认证的密码
  - **-x PROTOCOL** # 指定加密的算法
  - -**X PASSWORD** # 指定加密的密码

EXAMPLE

snmpwalk -v 2c -c public 192.168.0.2 # 最简单直接的 walk 方式

使用 V3 版本认证方式获取 SNMP 数据

```bash
# 华为服务器使用 snmpv3 访问。认证密码和加密密码都是 impi 用户的登录密码
snmpwalk -v3 -u root -l authPriv -a SHA -A Huawei12#$ -x AES -X Huawei12#$ 192.168.1.82
```

根据导入的私有第三方 MIB，获取 SNMP 数据（需要在 /etc/snmp/snmp.conf 中添加需加载的 MIB 文件配置；或使用 -m 选项指定 MIB 文件）

```bash
snmpwalk -v 2c -c public 172.19.42.241 HH3C-SERVER-AGENT-MIB:hh3c2014

snmpwalk -v 2c -c public 192.168.1.91 INSPUR-MIB
```

# net-snmp-create-v3-user

**net-snmp-create-v3-user \[-ro] \[-A authpass] \[-a MD5|SHA] \[-X privpass] \[-x DES|AES] \[username]**

OPTIONS

- **-A PASSWORD** # 指定认证的算法
- **-a PROTOCOL** # 指定认证的密码
- **-X PASSWORD** # 指定加密的算法
- **-x PROTOCOL** # 指定加密的密码
- **-ro** # 创建的用户有只读权限

EXAMPLE

- 创建一个 snmp 的 v3 用户，只读模式，认证算法为 SHA，认证密码是 nm@tjiptv，加密算法是 AES，加密密码是 nm@tjiptv，用户名是 nm
  - net-snmp-create-v3-user -ro -a SHA -A nm@tjiptv -x AES -X nm@tjiptv nm

# snmptranslate - 转换 OID 的格式

> 参考：
>
> - [官方手册，snmptranslate](http://net-snmp.sourceforge.net/docs/man/snmptranslate.html)
> - [官方手册，snmpcmd](https://net-snmp.sourceforge.io/docs/man/snmpcmd.html)
> - [Manual(手册)](https://man.cx/snmptranslate)

在数值和文本格式之间转换 MIB OID 的名称

**snmptranslate \[OPTIONS] OID \[OID]...**

**OPTIONS**

- **-T**(OPTIONS) # 控制输出的结果。有多个子选项可用
    - d # 输出详细信息
- **-O**(OPTIONS) # 控制 OID 的显示方式。有多个子选项可用
    - a # 以文本格式输出（ASCII）
    - n # 以数值格式输出
- **-m**(\[]STRING) # 指定额外的 MIB 文件
- **-M**(\[]STRING) # 指定额外的读取 MIB 文件的目录

**EXAMPLE**

- 根据当前配置，显示所有的 OID 的两种格式
  - snmptranslate -Tz

显示指定 OID 的数字及其详情

```bash
~]# snmptranslate -Td -On IF-MIB::ifInOctets
.1.3.6.1.2.1.2.2.1.10
ifInOctets OBJECT-TYPE
  -- FROM       IF-MIB
  SYNTAX        Counter32
  MAX-ACCESS    read-only
  STATUS        current
  DESCRIPTION   "The total number of octets received on the interface,
            including framing characters.

            Discontinuities in the value of this counter can occur at
            re-initialization of the management system, and at other
            times as indicated by the value of
            ifCounterDiscontinuityTime."
::= { iso(1) org(3) dod(6) internet(1) mgmt(2) mib-2(1) interfaces(2) ifTable(2) ifEntry(1) 10 }
```

# snmpgetnext

https://www.net-snmp.org/docs/man/snmpgetnext.html

使用 SNMP GETNEXT 请求与网络实体通信，获取当前 OID 的下一个 OID。