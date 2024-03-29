---
title: LVS 配置示例
linkTitle: LVS 配置示例
date: 2024-03-29T17:42
weight: 2
---

# 概述

> 参考：
> 
> -

# DR 模型配置样例

Director 配置

```bash
ipvsadm -A -t 10.10.100.107:80 -s rr
ipvsadm -a -t 10.10.100.107:80 -r 10.10.100.111 -g
ipvsadm -a -t 10.10.100.107:80 -r 10.10.100.112 -g
```

RS 的配置

设置 arp 参数

```bash
cat > /etc/sysctl.d/lvs-sysctl.conf << EOF
net.ipv4.conf.all.arp_ignore = 1
net.ipv4.conf.lo.arp_ignore = 1
net.ipv4.conf.all.arp_announce = 2
net.ipv4.conf.lo.arp_announce = 2
EOF
sysctl -p /etc/sysctl.d/*
```

配置 lo 网卡

    cat > /etc/sysconfig/network-scripts/ifcfg-lo:0 << EOF
    DEVICE=lo:0
    IPADDR=10.10.100.107
    NETMASK=255.255.255.255
    ONBOOT=yes
    NAME=loopback
    EOF
    ifup ifcfg-lo\:0

# 其他配置样例

## 配置一个 NAT 类型（2 台 RS）的集群

LVS(Thinkpad): if1=172.16.53.128/24(VIP) if2=10.1.1.77/16(DIP)

RS1(ibm1): 10.1.1.78/16(RIP1)

RS2(ibm2): 10.1.1.79/16(RIP2)

RS1 与 RS2 主机使用同一个 vmnet 网卡

如下图：
![](https://notes-learning.oss-cn-beijing.aliyuncs.com/rlona8/1616132564982-5b0ee7e2-e822-4976-8cf0-b10f9108d986.jpeg)

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/rlona8/1616132564990-81c808d9-696f-4cda-921e-1f3532fc2a6f.jpeg)

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/rlona8/1616132565016-ee9b26be-4fb0-4104-b1f8-73ccd019bc04.jpeg)

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/rlona8/1616132565003-365215ad-80a1-43c2-b7d0-2167171bbd4f.jpeg)

分别在 ibm1、ibm2 上装上几个服务

同步 ibm1 与 ibm2 的时间!这个在集群服务中很重要！！下面两步即同

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/rlona8/1616132564982-ee55e0f5-4497-47b9-bb50-7c46bf1028ab.jpeg)

2.重启 chronyd 服务

打开 ibm1 与 ibm2 上的配置 httpd 服务如下图效果

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/rlona8/1616132564958-21036f72-e72e-4b9e-a7e6-c529ced28984.jpeg)

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/rlona8/1616132564984-b8b52892-fdaa-4068-9dee-a1eeb12c4694.jpeg)

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/rlona8/1616132564979-267630af-71bf-4a2c-9e9e-e068c6f1947e.jpeg)

在 Thinkpad 主机上安装 ipvsadm

将 RS1 RS2 的默认网关指向 Director 的 DIP

route add default gw 10.1.1.77/16

当主机及拓扑结构及 ip 与各种软件都设置好的时候，LVS 的设置是非常简单与快速的一件事

ipvsadm -A -t 10.1.1.77:80 -s rr # 添加 LVS 集群与调度算法 ipvsadm -a -t 10.1.1.77:80 -r 10.1.1.78 -m -w 1 # 添加 LVS 集群主机与 LVS 调度模式及 RS 权重 ipvsadm -a -t 10.1.1.77:80 -r 10.1.1.79 -m -w 1 # 同上\[root@Thinkpad ~]# echo 1 > /proc/sys/net/ipv4/ip_forward # 打开 Director 的转发功能

至此,一个简单的根据 LVS rr 算法调度的负载均衡集群完成了.

## 配置一个 DR 类型（2 台 RS）的集群,三台虚拟机都桥接到物理机上.

Directory: DIP: 192.168.31.101/24 调度主机只需要一个网卡接口 VIP 用 DIP 的别名生成

RS1: 192.168.31.194/24

RS2: 192.168.31.220/24

编写以下脚本文件并分别在各 RS 上执行.(setparam.sh)

```bash
#!/bin/bash
#
vip=192.168.10.21
mask='255.255.255.255'
case $1 in
start)
    echo 1 > /proc/sys/net/ipv4/conf/all/arp_ignore
    echo 1 > /proc/sys/net/ipv4/conf/lo/arp_ignore
     echo 2 > /proc/sys/net/ipv4/conf/all/arp_announce
    echo 2 > /proc/sys/net/ipv4/conf/lo/arp_announce
    /sbin/ifconfig lo:0 $vip netmask $mask  broadcast $vip up
    route add -host $vip dev lo:0
    ;;
stop)
    /sbin/ifconfig lo:0  down
    echo 0 > /proc/sys/net/ipv4/conf/all/arp_ignore
    echo 0 > /proc/sys/net/ipv4/conf/lo/arp_ignore
    echo 0 > /proc/sys/net/ipv4/conf/all/arp_announce
    echo 0 > /proc/sys/net/ipv4/conf/lo/arp_announce
    route del -host $vip dev lo:0
    ;;
status)
    # Status of LVS-DR real server.
    islothere=`/sbin/ifconfig lo:0 | grep $vip`
    isrothere=`netstat -rn | grep "lo:0" | grep $vip`
    if [ ! "$islothere" -o ! "isrothere" ]; then
        # Either the route or the lo:0 device
        # not found.
        echo "LVS-DR real server Stopped."
    else
        echo "LVS-DR real server Running."
    fi
;;
*)
    echo "Usage $(basename $0) start|stop"
    exit 1
    ;;
esac
```

编写以下脚本在 Director(调度器)上添加 VIP 添加 LVS 规则等操作 (ipvs.sh)

```bash
#!/bin/bash
#
vip=192.168.10.21
iface='eth0:0'
mask='255.255.255.255'
port='80'
rs1='192.168.10.22'
rs2='192.168.10.23'
scheduler='rr'
case $1 in
start)
    ifconfig $iface $vip netmask $mask broadcast $vip up
    iptables -F
    ipvsadm -A -t ${vip}:${port} -s $scheduler
    ipvsadm -a -t ${vip}:${port} -r $rs1 -g
    ipvsadm -a -t ${vip}:${port} -r $rs2 -g
    ;;
stop)
    ipvsadm -C
    ifconfig $iface down
    ;;
*)
    echo ''Usage: $(basename $0) {start|stop|status}
    exit 1
    ;;
esac
```

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/rlona8/1616132565011-7f0f7b82-c341-4628-ab9b-a729530ca1d5.jpeg)

## FWM(FireWall Mark)借助于防火墙标记来分类报文，而后基于标记定义集群服务，可将多个不同的应用使用同一个集群服务进行调度。

示例：

    iptables -t mangle -A PREROUTING -d 10.1.0.5 -p tcp -m multiport --dports 80,443 -j MARK --set-mark 11
    ipvsadm -A -f 11 -s rr
    ipvsadm -a -t 192.168.31.100 -r 192.168.31.194 -g -w 1
    ipvsadm -a -t 192.168.31.100 -r 192.168.31.220 -g -w 1

实现目标 2：firewallmark（fwm）

fwm 是为指定的某一群体打标签，这一群体可以是某几个 ip 地址或者是几个端口，实现的作用是，可以将需要访问的几个目标地址或者端口定向到同一个标签，比如给 80,8080,443 端口打上同一个标签，那么就可以让用户访问这三个端口中的任意一个端口定向到同一个主机。

在上面的基础上开始实现 fwm：

1.host1 上面清理掉之前的规则：

ipvsadm -C

2.使用 mangle 在 PREROUTING 上面打上标签

iptables -A PREROUTNG -t mangle -d 172.16.52.57 -p tcp --dport 80 -j MARK --set-mark 9

ipvsadm -A -f 3 -s rr

ipvsadm -a -f 3 -r 172.16.52.60 -g

ipvsadm -a -f 3 -r 172.16.52.61 -g

## LVS 的持久连接机制 persiistence 这个功能与 FWM 结合可实现端口姻亲关系

LVS 基于其的持久连接模板，实现无论使用任何调度算法，在一段时间内，都能实现将来自同一个地址的请求始终发往同一个 RS。不同于 sh 算法在于其是可自定义的且没有超时时长。

    ipvsadm -A -t 10.1.0.5:80 -s rr -p 60     # 持久连接定义成60秒，在这个时间之后 rr算法才会生效
    ipvsadm -a -t 10.1.0.5:80 -r 10.1.0.7 -g -w 1
    ipvsadm -a -t 10.1.0.5:80 -r 10.1.0.8 -g -w 1

port Affinity：端口姻亲关系

每端口持久：每集群服务单独定义，并定义其持久性

每防火墙标记持久：基于防火墙标记定义持久的集群服务，可实现将多个端口上的应用统一调度，即所谓的 port Affinity

每客户端持久

实现目标 3：持久连接

eg1：

\[root@localhost scripts]# ipvsadm -A -t 172.16.52.57:80 -s rr -p

\[root@localhost scripts]# ipvsadm -a -t 172.16.52.57:80 -r 172.16.52.60 -g

\[root@localhost scripts]# ipvsadm -a -t 172.16.52.57:80 -r 172.16.52.61 -g

说明：基于端口的持久连接，-p 是指定持久时长，默认是 360s，但是 man 手册里面说的 300s，常用方式是： -p timeout

eg2：

\[root@localhost scripts]# ipvsadm -A -t 172.16.52.57:0 -s rr -p

\[root@localhost scripts]# ipvsadm -a -t 172.16.52.57:0 -r 172.16.52.60 -g

\[root@localhost scripts]# ipvsadm -a -t 172.16.52.57:0 -r 172.16.52.61 -g

说明：基于客户端的持久连接，使用端口号 0，后面必须要跟-p ，不加可以试试看，会报错的。。。

eg3：

\[root@localhost scripts]# iptables -A PREROUTING -t mangle -d 172.16.52.57 -p tcp --dport 80 -j MARK --set-mark 20

\[root@localhost scripts]# iptables -A PREROUTING -t mangle -d 172.16.52.57 -p tcp --dport 443 -j MARK --set-mark 20

\[root@localhost scripts]# ipvsadm -A -f 20 -s rr -p

\[root@localhost scripts]# ipvsadm -a -f 20 -r 172.16.52.60 -g -p

\[root@localhost scripts]# ipvsadm -a -f 20 -r 172.16.52.61 -g -p

说明：通过防火墙标记的方式来进行的持久连接,可以 cat /proc/net/ip_vs_conn 查看连接状态

总结：持久连接分为上面三种方式，但跟调度算法无关，调度算法无法保持持久连接

#
