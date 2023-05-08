---
title: Netplan 配置详解
---

# 概述

> 参考：
>
> - [官网，参考](https://netplan.io/reference)

Netplan 的配置文件使用 YAML 格式。`/{lib,etc,run}/netplan/*.yaml` 都是 Netplan 程序读取配置文件的路径。

# 配置文件详解

该 YAML 文件只有一个顶级节点：`network: <Object>`，其中包括 version、设备类型(例如 ethernets、modems、wifis、birdge 等)、renderer。

## version: \<INT>

## renderere: \<STRING>

## ethernetes: \<OBJECT>

以太网设备的专用属性

## bridge: \<OBJECT>

桥设备的专用属性

## 所有设备的通用属性

**addresses: <\[]OBJECT>** #
**dtcp4: \<BOOL>** # 为 IPv4 启用 DHCP。`默认值：false`
**dhcp6: \<BOOL>** # 为 IPv6 启用 DHCP。`默认值：false`
**gateway4 | gateway6: \<STRING>** # **已弃用**。使用 `routes` 字段。
**nameservers: \<OBJECT>** # 设置 DNS 服务器和搜索域，用于手动地址配置
**routes: <\[]OBJECT>** # 为设备配置静态路由；请参阅下面的路由部分。

# 配置示例

> 参考：
>
> - [官网，示例](https://netplan.io/examples)

```yaml
# This is the network config written by 'subiquity'
network:
  ethernets:
    ens3:
      addresses:
      - 172.19.42.248/24
      dhcp4: no
      dhcp6: no
      optional: true
     routes:
        - to: default
          via: 172.19.42.1
      nameservers:
        addresses:
        - 8.8.8.8
  version: 2
  renderer: networkd
```

生成如下配置

```bash
~]# cat /run/systemd/network/10-netplan-ens3.network
[Match]
Name=ens3

[Link]
RequiredForOnline=no

[Network]
LinkLocalAddressing=ipv6
Address=172.19.42.248/24
Gateway=172.19.42.1
DNS=8.8.8.8
```

## Bridge 配置示例

/etc/netplan/br0.yaml

```yaml
network:
  version: 2
  ethernets:
    eno1:
      dhcp4: false
      dhcp6: false
  bridges:
    br0:
      interfaces:
        - eno1
      dhcp4: false
      addresses:
        - 172.38.180.100/24
      routes:
        - to: default
          via: 172.38.180.254
      nameservers:
        addresses:
          - 8.8.8.8
      parameters:
        stp: false
      dhcp6: false
```

应用配置

**netplan apply**
