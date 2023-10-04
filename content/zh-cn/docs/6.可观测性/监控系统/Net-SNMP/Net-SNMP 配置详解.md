---
title: Net-SNMP 配置详解
---

# snmpd.conf 文件

# snmp.conf 文件

**mibdirs +/usr/share/snmp/mibs/h3c** # 添加 MIB 目录，以便 snmp 工具可以从这些目录中读取 MIB 文件。

**mibs +HH3C-OID-MIB** # 添加 snmp 工具可以 walk 等操作可以控制的 OID。

**mibs +HH3C-SERVER-AGENT-MIB**

**mibs +HH3C-SERVER-TRAP-MIB**

