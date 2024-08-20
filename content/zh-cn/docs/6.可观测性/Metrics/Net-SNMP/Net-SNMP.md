---
title: Net-SNMP
linkTitle: Net-SNMP
date: 2024-04-19T09:41
weight: 1
---

# 概述

> 参考：
>
> - [官网](http://net-snmp.sourceforge.net/)

Net-SNMP 是实现 [SNMP(传统监控标准)](/docs/6.可观测性/Metrics/监控系统概述/SNMP(传统监控标准).md) 的工具和库的集合。包含如下内容：

- **net-snmp** # SNMP 代理，用于采集设备的 SNMP 信息。包含两个守护程序。
- **net-snmp-utils** # 是一组工具的集合，包括下面这些命令行工具：
  - snmpwalk # 获取 SNMP 信息，可以根据 OID 获取指定 OID 的 SNMP 信息
  - snmptranslate # 转换 OID 的两种格式
  - encode_keychange、snmpbulkget、snmpbulkwalk、snmpdelta、snmpdf、snmpget、snmpgetnext、snmpinform、snmpnetstat、snmpset、snmpstatus、snmptable、snmptest、snmptls、snmptrap、snmpusm、snmpvacm

## net-snmp

是一种可以通过 snmp 协议来实现基础监控功能的守护程序。包含两个守护程序以及几个命令行工具

- **snmpd** # 用于响应请求的 SNMP 代理。说白了提供 SNMP 数据的程序。监听一个端口(默认监听 `161/udp`)，当别人向 161 发送 SNMP 请求时，snmpd 会将采集到的数据发送给对方。
- **snmptrapd** # 用于接收 SNMP 通知的通知接收器
- **agentxtrap** #
- **net-snmp-create-v3-user** # 创建 v3 用户
- **snmpconf** # 用于生成配置文件

# Net-SNMP 安装

Net-SNMP 提供了各种安装文件，对于 Linux 系统来说，直接使用 yum、apt 等工具安装即可

CentOS

```bash
yum install -y net-snmp*
```

Ubuntu

```bash
apt-get install snmp snmp-mibs-downloader
```

# net-snmp 关联文件

**/etc/snmp/snmpd.conf** # snmpd 根据该文件定义运行时行为

**/etc/snmp/snmp.conf** # [net-snmp-utils](/docs/6.可观测性/Metrics/Net-SNMP/net-snmp-utils.md) 包中的相关工具所用的配置文件。snmpd 也会从这个文件中指定的 MIB 路径中加载 MIB 信息

**/etc/snmp/snmptrapd.conf** # snmptrapd 运行时配置文件
