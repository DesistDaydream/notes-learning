---
title: Router And Switch
linkTitle: Router And Switch
date: 2024-02-22T13:23
weight: 3
---

# 概述

> 参考：
> 
> - [Wiki，Router_(computing)](https://en.wikipedia.org/wiki/Router_(computing))
> - [Wiki，Network switch](https://en.wikipedia.org/wiki/Network_switch)

**Router(路由器)**

**Switch(交换机)**

# 华为交换机

## 创建 trunk

```bash
sy
[HUAWEI] int eth-trunk 1
[HUAWEI-Eth-Trunk1] mode manual load-balance
[HUAWEI-Eth-Trunk1] port link-type trunk
[HUAWEI-Eth-Trunk1] quit
[HUAWEI] int g1/0/0
[HUAWEI-GigabitEthernet1/0/0] eth-trunk 1
[HUAWEI-GigabitEthernet1/0/0] q
[HUAWEI] int g 1/0/1
[HUAWEI-GigabitEthernet1/0/1] eth-trunk 1
[HUAWEI-GigabitEthernet1/0/1] q
```

配置 SNMP

```bash
~]# telnet 172.19.42.200
Trying 172.19.42.200...
Connected to 172.19.42.200.
Escape character is '^]'.

Warning: Telnet is not a secure protocol, and it is recommended to use Stelnet.

Login authentication


Username:root
Password:
Warning: The initial password poses security risks.
The password needs to be changed. Change now? [Y/N]: n
Info: The max number of VTY users is 10, and the number
      of current VTY users on line is 1.
      The current login time is 2000-06-28 00:17:28+00:00.
# 进入系统视图
<HUAWEI>sys
# 开启 snmp
[HUAWEI]snmp-agent
# 配置 snmp 版本
[HUAWEI]snmp-agent sys-info version v2c
Warning: SNMPv1/SNMPv2c is not secure, and it is recommended to use SNMPv3.
# 关闭团体名密码复杂度检查功能，网上的傻逼资料都不关，不管就没法配置简单的团体名
[HUAWEI]snmp-agent community complexity-check disable
Warning: Does not recommend to disable complexity check. A simple community name may result in security threats.
# 配置 read 模式的团体名
[HUAWEI]snmp-agent community read public
```
