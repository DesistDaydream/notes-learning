---
title: "通过 veth 设备查找对应的虚拟机"
linkTitle: "通过 veth 设备查找对应的虚拟机"
weight: 20
---

# 概述

通过 xml 文件查看虚拟机的 mac 地址

```bash
~]# grep "mac addr" lchTest.xml
<mac address='52:54:00:6a:86:89'/>
```

筛选网络设备 mac

```bash
~]# ip a | grep "86:89" -B 1
127: vnet1: <BROADCAST,MULTICAST,UP,LOWER_UP> mtu 1500 qdisc pfifo_fast master br1 state UNKNOWN group default qlen 1000
    link/ether fe:54:00:6a:86:89 brd ff:ff:ff:ff:ff:ff
```

由此可见，vnet 设备的 mac 与虚拟机的 mac 在后面几位是永远保持一致的，所以可以通过 vnet 设备的 mac ，从所有虚拟机中的 xml 进行筛选就行可以找到 vnet 设备对应的虚拟机了。

## 应用实例

比如我现在能看到 4 个 vnet 设备。

```bash
~]# ip a
.....
62: vnet3: <BROADCAST,MULTICAST,UP,LOWER_UP> mtu 1500 qdisc pfifo_fast master br1 state UNKNOWN group default qlen 1000
    link/ether fe:54:00:5c:11:85 brd ff:ff:ff:ff:ff:ff
    inet6 fe80::fc54:ff:fe5c:1185/64 scope link
       valid_lft forever preferred_lft forever
79: vnet0: <BROADCAST,MULTICAST,UP,LOWER_UP> mtu 1500 qdisc mq master br1 state UNKNOWN group default qlen 1000
    link/ether fe:54:00:68:20:51 brd ff:ff:ff:ff:ff:ff
    inet6 fe80::fc54:ff:fe68:2051/64 scope link
       valid_lft forever preferred_lft forever
127: vnet1: <BROADCAST,MULTICAST,UP,LOWER_UP> mtu 1500 qdisc pfifo_fast master br1 state UNKNOWN group default qlen 1000
    link/ether fe:54:00:6a:86:89 brd ff:ff:ff:ff:ff:ff
    inet6 fe80::fc54:ff:fe6a:8689/64 scope link
       valid_lft forever preferred_lft forever
139: vnet2: <BROADCAST,MULTICAST,UP,LOWER_UP> mtu 1500 qdisc pfifo_fast master br1 state UNKNOWN group default qlen 1000
    link/ether fe:54:00:3a:95:ef brd ff:ff:ff:ff:ff:ff
    inet6 fe80::fc54:ff:fe3a:95ef/64 scope link
       valid_lft forever preferred_lft forever
```

我想看看 vnet3 是关联到哪个虚拟机上了，就可以进行如下操作：(首先看到 vnet3 的 mac 为 fe:54:00:5c:11:85)

```bash
~]# grep -r "5c:11:85" /etc/libvirt/qemu/*
/etc/libvirt/qemu/cobbler.test.tjiptv.net.xml:      <mac address='52:54:00:5c:11:85'/>
```

这时候可以看到，是 cobbler.test.tjiptv.net.xml 这台虚拟机在使用 vnet3 。
