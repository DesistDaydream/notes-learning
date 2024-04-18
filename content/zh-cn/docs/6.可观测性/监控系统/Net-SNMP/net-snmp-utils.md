---
title: net-snmp-utils
---

# snmpwalk 与 snmpget

snmpwalk 与 snmpget 使用 snmp 协议的 GETNEXT 请求，向 SNMP 代理发送查询请求，以便获取 SNMP 数据。

- snmpget 获取指定 OID 的数据
- snmpwalk 可以获取大量 OID 的数据

## snmpwalk 配置

**/etc/snmp/snmp.conf** # snmpwalk 运行时配置文件。若不存在则手动创建

**/usr/local/share/snmp/mibs** # MIB 文件的默认路径。这里现阶段包含 66 个 MIB 文件。

`net-snmp-config --default-mibdirs` 命令可以列出工具在使用中会读取 MIB 文件的路径，包括如下几个。

- 注意，CentOS 和 Ubuntu 的路径可能不相同。这种老程序是真滴坑。
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
mibs +HH3C-SERVER-TRAP-MIB
```

## snmpwalk

**snmpwalk \[OPTIONS] AGENT \[OID]**

OPTIONS

- **-l \<noAuthNoPriv | authNoPriv | authPriv>** # 设置安全级别
- **-c \<STRING>** # 团体名
- **-v <1 | 2c | 3>** # snmp 版本
- 认证相关选项
  - **-A PASSWORD** # 指定认证的算法
  - -**a PROTOCOL** # 指定认证的密码
  - -**X PASSWORD** # 指定加密的算法
  - **-x PROTOCOL** # 指定加密的密码

EXAMPLE

- snmpwalk -v 3 -u nm -l authPriv -a SHA -A nm@tjiptv -x AES -X nm@tjiptv 10.10.100.101 #
- snmpwalk -v 2c -c public 192.168.0.2
- walk 第三方 MIB 内容
  - snmpwalk -v 2c -c public 172.19.42.241 HH3C-SERVER-AGENT-MIB:hh3c2014

## snmpget

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

# snmptranslate # 转换 OID 的格式

> 参考：
> 
> - [官方文档](http://net-snmp.sourceforge.net/docs/man/snmptranslate.html)
> - [Manual(手册)](https://man.cx/snmptranslate)

在数字格式和文本格式之间转换 MIB 的 OID

**snmptranslate \[OPTIONS] OID \[OID]...**

**OPTIONS**

**EXAMPLE**

- 根据当前配置，显示所有的 OID 的两种格式
  - snmptranslate -Tz -m all
