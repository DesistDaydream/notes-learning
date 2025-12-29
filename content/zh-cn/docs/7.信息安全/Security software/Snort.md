---
title: Snort
linkTitle: Snort
weight: 20
---

# 概述

> 参考：
>
> - [GitHub 项目，snort3/snort3](https://github.com/snort3/snort3)
> - [官网](https://www.snort.org/)
> - [Wiki, Snort_(software)](https://en.wikipedia.org/wiki/Snort_(software))

Snort 是世界上最重要的开源 Intrusion Prevention System(入侵防御系统，IPS)。 Snort IPS 使用一系列规则来帮助定义恶意网络活动，并使用这些规则来查找与其匹配的数据包并为用户生成警报。

Snort 也可以内联部署来阻止这些数据包。 Snort 有三个主要用途：作为数据包嗅探器（如 tcpdump）、作为数据包记录器 — 这对于网络流量调试很有用，或者可以用作成熟的网络入侵防御系统。 Snort 可以下载并配置用于个人和商业用途。

# Snort 规则

> 参考：
>
> - [官方文档，规则](https://docs.snort.org/rules/)

Snort 规则主要由两部分组成

- **Rule header** # 定义了流量的基础规则，协议、源/目 的 IP 和 PORT，最基本就是这 5-tuple。
- **Rule body** # 类似于 7 层策略。定义了与指定规则关联的数据包的内容应该如何匹配。用 `()` 包裹起来。
  - Rule body 中包括很多可用的 OPTIONS，详见 https://docs.snort.org/rules/options/ ，比如 msg、flow、etc. 都是 OPTIONS。

以下是具有 Rule header 和 Rule body 定义的完整形式的 Snort 3 规则的示例：

```text
alert tcp $EXTERNAL_NET 80 -> $HOME_NET any
(
    msg:"Attack attempt!";
    flow:to_client,established;
    file_data;
    content:"1337 hackz 1337",fast_pattern,nocase;
    service:http;
    sid:1;
)
```

|  关键字  |   协议    |      源 IP      |      源 PORT      | 关键字 |            目的 IP             |            目的 PORT             |            规则主体            |
| :---: | :-----: | :------------: | :--------------: | :-: | :--------------------------: | :----------------------------: | :------------------------: |
| alert | tcp/udp | 可以用 any 或具体 IP | 可以用 any 或具体 PORT | ->  | 可以用 any 或具体 IP。需要使用 `[]` 括起来 | 可以用 any 或具体 PORT。需要添加 `[]` 括起来 | 除了 源/目 之外的匹配规则。使用 `()` 括起来 |

**匹配规则** # 除了 源/目 IP 或 PORT 之外，还可以匹配数据包的 URL、数据内容 等进行过滤。只有完全匹配到的才会记录成安全事件。

## Rule header

https://docs.snort.org/rules/headers/

规则标题由五个主要部分组成

- **Actions(动作)** # 当规则触发时应该执行的动作
- **Protocol(协议)** # 针对哪些协议评估规则
- **IP address(IP地址)** # 针对哪些网络评估规则
- **Port(端口)** # 针对哪些端口评估规则
- **Direction operator(操作方向)** # 针对流量的哪种传输方向评估规则。
    - e.g. `->`, `<-`, etc. 之类的关键字

### IP Address

https://docs.snort.org/rules/headers/ips

基本示例

```snort
# look for traffic sent from the 192.168.1.0/24 subnet to the
# 192.168.5.0/24 subnet
alert tcp 192.168.1.0/24 any -> 192.168.5.0/24 any (
```

```snort
# look for traffic sent from addresses included in the
# defined $EXTERNAL_NET variable to addresses included in the defined
# $HOME_NET variable
alert tcp $EXTERNAL_NET any -> $HOME_NET 80 (
```

```snort
# look for traffic sent from any source network to the IP address, 192.168.1.3
alert tcp any any -> 192.168.1.3 445 (
```

```snort
alert tcp !192.168.1.0/24 any -> 192.168.1.0/24 23 (
```

```snort
alert tcp ![192.168.1.0/24,10.1.1.0/24] any -> [192.168.1.0/24,10.1.1.0/24] 80 (
```

## Rule body

详见 [Snort Rule body](/docs/7.信息安全/Security%20software/Snort%20Rule%20body.md)

# Snort 规则最佳实践

`alert udp 159.138.48.8 any -> [111.30.108.165] [53] (msg:\"ip test\"; pcre:\"/./\"; sid:99999999;)`
