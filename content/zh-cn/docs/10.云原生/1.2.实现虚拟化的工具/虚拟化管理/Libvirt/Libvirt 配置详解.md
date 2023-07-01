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

# 最佳实践

## 开启 libvirtd 的 TCP 监听

> 参考：
> 
> - [StackOverflow，could-not-add-the-parameter-listen-to-open-tcp-socket](https://stackoverflow.com/questions/65663825/could-not-add-the-parameter-listen-to-open-tcp-socket)
> - [libvirtd 官方手册](https://libvirt.org/manpages/libvirtd.html#system-socket-activation)

`systemctl stop libvirtd.service`

在 `/etc/libvirt/libvirtd.conf` 文件中添加 `auth_tcp="none"`

`systemctl enable libvirtd-tcp.socket --now`

`systemctl start libvirtd.service`



