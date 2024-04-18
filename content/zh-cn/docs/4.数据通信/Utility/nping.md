---
title: nping
linkTitle: nping
date: 2024-04-18T11:21
weight: 20
tags:
  - Network_analyzer
---

# 概述

> 参考：
>
> - [Nmap，nping](https://nmap.org/nping/)

由 [Nmap](docs/4.数据通信/Utility/Nmap.md) 项目组开发的类似 [hping](docs/4.数据通信/Utility/hping.md) 的工具

每次的发包都称为一次 Probe(探测)

# Syntax(语法)

> 参考：
> 
> - [Manual(手册)，nping](https://nmap.org/book/nping-man.html)

`nping [OPTIONS] {<targets>}`

## OPTIONS

**--tcp** # 使用 TCP 探针模式

### TCP 探针模式选项

**-g, --source-port PORT_NUMBER** # 设置源端口

**--flags FLAG_LIST** # 设置 TCP 的 Flag, 多个 Flag 以逗号分隔


### IPv4 选项

**-S, --source-ip** # 设置源 IP 地址

## Timing 和 Performance

**--delay TIME** # 前后两次发送探测行为的间隔. i.e. 每隔 TIME 探测一次

### 其他

**-c, --count N** # 在发送 N 次后停止程序. 若 N 设为 0 则永不停止.

# EXAMPLE

nping --tcp -p 10080 --flags SYN --source-ip 192.168.180.131 --delay 1000ms --count 0 172.24.194.179