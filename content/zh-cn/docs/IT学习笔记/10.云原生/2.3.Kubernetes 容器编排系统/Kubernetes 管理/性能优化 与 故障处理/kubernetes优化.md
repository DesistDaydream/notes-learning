---
title: kubernetes优化
---

# 概述

> 参考：

增加可以打开的文件数与线程数,防止 pod 无故无法启动

```bash
cat >> /etc/security/limits.conf << EOF
* soft nproc 1000000
* hard nproc 1000000
* soft nofile 1000000
* hard nofile 1000000
EOF
```

## 馆长推荐的优化参数

<https://github.com/moby/moby/issues/31208>&#x20;
\# ipvsadm -l --timout
\# 修复 ipvs 模式下长连接 timeout 问题 小于 900 即可
{% if proxy.mode is defined and proxy.mode == 'ipvs' %}
net.ipv4.tcp_keepalive_time = 600
net.ipv4.tcp_keepalive_intvl = 30
net.ipv4.tcp_keepalive_probes = 10
{% endif %}
===========

net.ipv4.tcp_fin_timeout = 30
net.ipv4.tcp_max_tw_buckets = 5000
net.ipv4.tcp_syncookies = 1
net.ipv4.tcp_max_syn_backlog = 1024
net.ipv4.tcp_synack_retries = 2

net.core.somaxconn = 10000
net.ipv6.conf.all.disable_ipv6 = 1
net.ipv6.conf.default.disable_ipv6 = 1
net.ipv6.conf.lo.disable_ipv6 = 1
net.ipv4.neigh.default.gc_stale_time = 120
net.ipv4.conf.all.rp_filter = 0
net.ipv4.conf.default.rp_filter = 0
net.ipv4.ip_forward = 1
net.bridge.bridge-nf-call-ip6tables = 1
net.bridge.bridge-nf-call-iptables = 1
net.netfilter.nf_conntrack_max = 2310720
fs.inotify.max_user_watches=89100
fs.may_detach_mounts = 1
fs.file-max = 52706963
fs.nr_open = 52706963
net.bridge.bridge-nf-call-arptables = 1

{% if not kubelet.swap %}
vm.swappiness = 0
{% endif %}

vm.overcommit_memory=1
vm.panic_on_oom=0
vm.max_map_count = 262144
