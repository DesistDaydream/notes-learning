---
title: Cilium
---

# 概述

> 参考：
> 
> - [GitHub 项目，cilium/cilium](https://github.com/cilium/cilium)
> - [官网](https://cilium.io/)
> - [官方文档](https://docs.cilium.io/en/latest/)

<https://docs.google.com/presentation/d/1cZJ-pcwB9WG88wzhDm2jxQY4Sh8adYg0-N3qWQ8593I/edit#slide=id.g7608b8c2de_0_0>
<https://www.youtube.com/watch?v=bIRwSIwNHC0>
<http://arthurchiao.art/blog/ebpf-and-k8s-zh/>

# 常见问题

如果在设备 A 添加了到 PodIP 段的静态路由，从集群外部直接访问 pod ip 是不通的。。。。。icmp 行。。。其他不行。。好像在 datapath 处理的时候，给略过了。。。。
