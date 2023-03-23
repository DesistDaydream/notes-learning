---
title: 问题实例：Keepalived 非抢占模式 VIP 不漂移
---

# Keepalived 非抢占模式下 VIP 不漂移问题

Keepalived 主要是通过虚拟路由冗余来实现高可用功能。本文将不对 keepalived 的基本原理进行阐述，可参考文章 Keepalived 详细介绍简介、keepalived vip 漂移基本原理及选举算法。本文记录了在实践过程中使用 keepalived 时，在 weight 值变化的情况下 vip 不漂移的问题及解决方法。

场景

3 个 keepalived 节点, vip 为 172.31.23.6：

- devops1a-zoocassa0 172.31.23.22

- devops1a-zoocassa1 172.31.23.23

预期

1. 两个节点初始都设为 BACKUP，按照优先级（priority）选举 MASTER；

2. 在两个节点上检查 memcached 服务状态，失败则降低优先级；

3. 如果 MASTER(假设为 devops1a-zoocassa0)上检查失败，BACKUP 上检查成功，则优先级高的 BACKUP 节点(假设为 devops1a-zoocassa1)切换为 MASTER 节点；

4. 之前检查失败的 MASTER(devops1a-zoocassa0)上的服务恢复时, 之前的 BACKUP 节点(devops1a-zoocassa1)服务检查也成功，即使 devops1a-zoocassa0 优先级恢复到高于 devops1a-zoocassa1,也不再成为 MASTER(不抢占)。

## 不成功配置范例

主节点 dr-1 配置

    global_defs {
      router_id k8s-master-dr
    }
    vrrp_script check_nginx {
     script "pidof nginx"
     interval 3
     weight -2
     fall 2
     rise 2
    }
    vrrp_instance VI_K8S {
       state BACKUP
       interface eth0
       virtual_router_id 60
       priority 101
       nopreempt
       authentication {
           auth_type PASS
           auth_pass 4be37dc3b4c90194d1600c483e10ad1d
       }
       virtual_ipaddress {
           172.40.0.60
       }
       track_script {
           check_nginx
       }
    }

备节点 dr-2 配置

    global_defs {
      router_id k8s-master-dr
    }
    vrrp_script check_nginx {
     script "pidof nginx"
     interval 3
     weight -2
     fall 2
     rise 2
    }
    vrrp_instance VI_K8S {
       state BACKUP
       interface eth0
       virtual_router_id 60
       priority 100
       nopreempt
       authentication {
           auth_type PASS
           auth_pass 4be37dc3b4c90194d1600c483e10ad1d
       }
       virtual_ipaddress {
           172.40.0.60
       }
       track_script {
           check_nginx
       }
    }

以上述配置文件内容作为 keepalived 配置文件 /etc/keepalived/keepalived.conf，在两个个节点上启动 keepalived：`systemctl restart keepalived`

会发现存在如下问题：

1. 优先级高的 dr-1 可能没有成为 MASTER 节点（多试几次，可能每次选举的 MASTER 节点都不同），不符合预期中的第 1 点；

2. 假设 dr-1 成为了 MASTER 节点，关掉 dr-1 上的 memcached 服务：`systemctl stop keepalived`

此时运行 service keepalived status，发现 dr-1 的 weight 值降低且低于 dr-2 ，但是 dr-2 并没有成为 MASTER 节点，不符合预期中的第 3 点。

1. 将配置文件中的 nginx 去掉以后，可以解决上述问题，符合预期中的第 1，2，3 点，但是当原 MASTER 节点上服务恢复后，原 MASTER 会重新成为 MASTER 角色，这不符合预期中的第 4 点（不抢占）；

## 问题原因：

在网上查阅到的资料中，大都认为按照上述配置后可以完全符合预期中的 4 个点，不会出现 MASTER 节点服务检查失败后 VIP 不漂移的问题。但是实践是检验真理的唯一标准，配置 nopreemt 后，不仅是会让原 MASTER 节点服务恢复后不抢占，而是会完全的不选举新 MASTER(从头到尾永远不切换，除非 BACKUP 认为当前集群中不存在 MASTER, 才会重新选举)，这样便可以解释出现的问题 1 和问题 2 了：

问题 1 的原因在于：

1. 先启动的节点将自己选举为 MASTER, 在收到其他节点的 vrrp 报文后不会按照优先级调整自己的角色；

2. 后启动的节点收到了 MASTER 的 vrrp 报文，发现已经存在 MASTER，由于不抢占，自动进入 BACKUP 状态；

问题 2 的原因在于：

1. 设置了 nopreempt, 永远不发生角色切换；

下面是官方文档中对于 nopreempt 的解释：

    "nopreempt" allows the lower priority machine to maintain the master role,
    even when a higher priority machine comes back online.
    NOTE: For this to work, the initial state of this entry must be BACKUP.

## 解决方案

要想同时满足预期中的效果，其实只要做到两点：

当 MASTER 上的服务检查失败时，触发重新选举；

设置不抢占（已经做到）；

那么如何实现第一点呢？重新选举意味着:

1. BACKUP 成为 MASTER，要求 BACKUP 节点认为当前节点中没有 MASTER 节点；

2. MASTER 成为 BACKUP，要求 MASTER 节点感知到环境中存在别的 MASTER 节点，从而进入 BACKUP 状态；

节点之间通过 VRRP 报文获得相互的优先级及状态信息，因此，可以通过在服务检查失败时，配置防火墙，禁止本机的 VRRP 报文发出即可。这样，BACKUP 节点收不到 MASTER 节点的 VRRP 报文，认为 MASTER 节点不存在，同时 MASTER 节点能收到其他节点的 VRRP 报文，感知到新 MASTER 的产生，从而进入 BACKUP 状态。

配置详见：[keepalived+nginx 配置示例](https://www.yuque.com/go/doc/33183799)

重启 keepalived 服务，测试成功。
