---
title: nping
linkTitle: nping
date: 2024-04-18T11:21
weight: 20
---

# 概述

> 参考：
>
> - [Nmap，nping](https://nmap.org/nping/)

由 [Nmap](docs/4.数据通信/Utility/Nmap.md) 项目组开发的类似 [hping](docs/4.数据通信/Utility/hping.md) 的工具

# EXAMPLE

nping --tcp -p 10080 --flags SYN --source-ip 192.168.180.131 --delay 1000ms --count 0 172.24.194.179