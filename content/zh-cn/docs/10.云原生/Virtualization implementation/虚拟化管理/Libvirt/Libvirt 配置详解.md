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

**log_level**(INT) # 程序运行日志的输出级别。`默认值: 2`。1 debug, 2 information, 3 warnings, 4 errors




