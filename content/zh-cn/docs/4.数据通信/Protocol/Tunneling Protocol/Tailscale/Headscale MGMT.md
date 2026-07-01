---
title: Headscale MGMT
created: 2026-07-01T14:34
weight: 101
---

# 概述

> 参考：
>
> - 

# 最佳实践

## tailscale 默认网段与阿里云 DNS 网段冲突

https://nyan.im/p/troubleshoot-tailscale

本质：tailscale 使用 100.64.0.0/10 网段，阿里云 DNS 100.100.2.138。阿里云的 DNS 地址刚好被规则丢了。

解决方式：

- 治标
    - 删除 ts-input 链中 `-A ts-input -s 100.64.0.0/10 ! -i tailscale0 -j DROP` 这条规则
- 治本
    - 暂无

# 修改节点 IP

> [!Tip]
>
> https://github.com/juanfont/headscale/issues/3133 未合并。现阶段只能通过修改 Sqlite 实现。

查看各节点 IP

```bash
sqlite3 -header -column /var/lib/headscale/db.sqlite "SELECT id, hostname, ipv4, ipv6 FROM nodes;"
```

```bash
id  hostname          ipv4        ipv6             
--  ----------------  ----------  -----------------
1   LAPTOP-5113MLDP   100.64.0.1  fd7a:115c:a1e0::1
2   aliyun-ubuntu-01  100.64.0.2  fd7a:115c:a1e0::2
```

更新 1 和 2 节点的 ipv4

```bash
sqlite3 -header -column /var/lib/headscale/db.sqlite "UPDATE nodes SET ipv4='100.64.0.2' WHERE id=1;"
sqlite3 -header -column /var/lib/headscale/db.sqlite "UPDATE nodes SET ipv4='100.64.0.1' WHERE id=2;"
```

