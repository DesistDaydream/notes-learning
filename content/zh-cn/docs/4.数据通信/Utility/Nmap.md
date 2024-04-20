---
title: Nmap
linkTitle: Nmap
date: 2024-04-18T09:03
weight: 20
tags:
  - Network_analyzer
---

# 概述

> 参考：
>
> - [官网](https://nmap.org/)
> - [Manual(手册)，NMAP(1)](https://nmap.org/book/man.html)

**Network Mapper(网络映射器，简称 Nmap)** 是一个用于网络探索和安全审计的开源工具。旨在快速扫描大型网络。

Nmap 项目除了有自己的 nmap 工具外，还有很多实用程序，比如 [Netcat](/docs/4.数据通信/Utility/Netcat.md)、etc.

# Syntax(语法)

**nmap \[Scan Type...] \[OPTIONS] {TARGET}**

- **Scan Type(扫描类型)** #
- **TARGET** # 扫描目标

直接使用 `nmap IP` 即可开始一个简单的扫描任务

## OPTIONS

> 参考：
>
> - https://nmap.org/book/man-briefoptions.html


### 规避防火墙/IDS 与 伪装

> FIREWALL/IDS EVASION AND SPOOFING

TODO

```
-f; --mtu <val>: fragment packets (optionally w/given MTU)
-D <decoy1,decoy2[,ME],...>: Cloak a scan with decoys
-S <IP_Address>: Spoof source address
-e <iface>: Use specified interface
-g/--source-port <portnum>: Use given port number
--proxies <url1,[url2],...>: Relay connections through HTTP/SOCKS4 proxies
--data <hex string>: Append a custom payload to sent packets
--data-string <string>: Append a custom ASCII string to sent packets
--data-length <num>: Append random data to sent packets
--ip-options <options>: Send packets with specified ip options
--ttl <val>: Set IP time-to-live field
--spoof-mac <mac address/prefix/vendor name>: Spoof your MAC address
--badsum: Send packets with a bogus TCP/UDP/SCTP checksum
```

# EXAMPLE
