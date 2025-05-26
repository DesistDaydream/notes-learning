---
title: firewalld 命令行工具
---

# firewall-cmd

> 参考：
>
> - [Manual, firewall-cmd](https://firewalld.org/documentation/man-pages/firewall-cmd)

所有命令都是对当前默认 ZONE(通过 firewall-cmd --get-default-zone 命令获得当前默认 zone)进行操作，如果想要对指定 ZONE 操作，需要使用 --zone=XXX

## OPTIONS

### 状态选项

- **--reload** # 重新加载防火墙规则并保留连接状态信息。注意：reload 会删除所有 runtime 模式的配置并应用 permanent 模式的配置。但是已经建立的连接不受影响(e.g.已经在对本机长 ping 的设备不会断开连接)
- **--complete-reload** # 重新加载防火墙规则并丢弃连接状态信息。注意：与 reload 不同，已经建立的连接会被丢弃(e.g.已经在对本机长 ping 的设备会断开连接)

### Log Denied Options

### Permanent Options

- **--permanent** # 开启永久模式，在该模式的配置都会永久保留

### Zone Options

### 查询

- **--get-default-zone** # 打印出当前默认的 ZONE
- **--list-all** # 列出所有已添加或已启用的内容
- **--list-service**s # 列出一个 ZONE 中已经添加了的 service
- **--list-interfaces** # 列出一个 ZONE 中已经绑定的网络设备
- **--list-rich-rules** # 列出一个 ZONE 中已经添加的丰富语言规则

TODO

- --add-source=<SOURCE> #  绑定 SOURCE 到一个 ZONE。SOURCE 可以使 MASK、MAC、ipset

## EXAMPLE

- firewall-cmd --zone=drop --change-interface=eth1 # 将经由 eth1 网卡的所有流量放在 drop 区域中进行处理
- firewall-cmd --zone=public --add-service=cockpit --permanent # 在永久模式下，允许 cockpit 服务的流量通过 public 这个 ZONE 内的网络设备
- 增
  - firewall-cmd --add-source=10.10.10.0/24 # 绑定 10.10.10.0/24 网段的 IP 到默认 ZONE 中
  - firewall-cmd --add-port=9100/tcp --permanent # 永久添加 9100 端口到默认 zone 中。
- 查
  - firewall-cmd --list-services # 列出当前默认 ZONE 中已经添加了的 service
- rich-rule(丰富语言规则)i.e.高级模式，可以配置更详细的规则
  - firewall-cmd --add-rich-rule='rule family=ipv4 source address="10.10.10.10" port port=1234 protocol=tcp accept'
    - -A IN_public_allow -s 10.10.10.10/32 -p tcp -m tcp --dport 1234 -m conntrack --ctstate NEW -j ACCEPT
