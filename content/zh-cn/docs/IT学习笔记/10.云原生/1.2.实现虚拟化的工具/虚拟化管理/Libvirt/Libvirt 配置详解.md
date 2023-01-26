---
title: Libvirt 配置详解
---

# 概述

> 参考：

/etc/libvirt/libvirt.conf

```bash
# 设置别名
uri_aliases = [
"vs-1=qemu+ssh://10.10.100.201/system",
]

# 可以对 10.10.100.201 使用 virsh 命令
virsh -c vs-1 list
```
