---
title: WSL 配置详解
linkTitle: WSL 配置详解
date: 2024-01-13T17:48
weight: 102
---

# 概述

> 参考：
>
> - [官方文档，WSL 配置](https://learn.microsoft.com/en-us/windows/wsl/wsl-config)

wsl.conf 和 .wslconfig 是 [INI](/docs/2.编程/无法分类的语言/INI.md) 格式的配置文件，两者互相配合以定义 WSL 虚拟机的运行方式

- wsl.conf 是在每个 WSL 中的 Linux 发行版内部的配置。通常是指 **本地配置**（tips: 本地配置就是指在 WSL 系统中，也就是 Linux 系统中的配置）
- .wslconfig 是在 Windows 中为所有 WSL发行版配置。通常是指 **全局配置**

# wsl.conf

## boot

**systemd**(BOOLEAN) # 是否启用 systemd

# .wslconfig

.wslconfig 包含两个部分: [wsl2](#wsl2) 和 [experimental](#experimental)

## \[wsl2]

https://learn.microsoft.com/en-us/windows/wsl/wsl-config#main-wsl-settings

**networkingMode**(STRING) # 如果该值是 mirrored，则这将打开镜像网络模式。`默认值: NAT`

- mirrored 会让虚拟机镜像本地网络。WSL2 和 Windows 主机的网络互通而且 IP 地址相同了，还支持 IPv6 了，并且从外部（比如局域网）可以同时访问 WSL2 和 Windows 的网络。这波升级彻底带回以前 WSL1 那时候的无缝网络体验了，并且 Windows 防火墙也能过滤 WSL 里的包了，再也不需要什么桥接网卡、端口转发之类的操作了。

**dnsTunneling**(BOOLEAN) # 更改 DNS 请求从 WSL 代理到 Windows 的方式。`默认值: false`

**firewall**(BOOLEAN) # 将此设置为 true 允许 Windows 防火墙规则以及特定于 Hyper-V 流量的规则来过滤 WSL 网络流量。`默认值: true`

**autoProxy**(BOOLEAN) # 强制 WSL 使用 Windows 的 HTTP 代理信息。`默认值: false`

## \[experimental]

https://learn.microsoft.com/en-us/windows/wsl/wsl-config#experimental-settings

**autoMemoryReclaim**(STRING) # `默认值: disable`

# 最佳实

## WSL 配置网络

### 配置桥接网络和静态 IP

> 参考：
>
> - [博客园，WSL2使用桥接网络，并指定IP](https://www.cnblogs.com/lic0914/p/17003251.html)
>   - 该文章参考的原文: https://github.com/luxzg/WSL2-fixes/blob/master/networkingMode%3Dbridged.md
>   - 上面这些做法可以弃用了，使用 https://github.com/microsoft/WSL/issues/10753#issuecomment-1814839310

常见问题

- https://github.com/microsoft/WSL/issues/10753#issuecomment-1814839310
- https://zhuanlan.zhihu.com/p/657110386

在 WSL2 的 [Release 2.0.0](https://github.com/microsoft/WSL/releases/tag/2.0.0) 版本更新日志中提到了网络模式，可以镜像主机，这样就不用任何配置即可使用主机网络、主机代理、等

```ini
[experimental]
autoMemoryReclaim=gradual | dropcache | disabled # 好像推荐用 gradual
networkingMode=mirrored
dnsTunneling=true
firewall=true
autoProxy=true
```

再新一点的版本可以把多个配置移到 wsl2 部分

```ini
[wsl2]
networkingMode=mirrored
dnsTunneling=true
firewall=false
autoProxy=true
[experimental]
autoMemoryReclaim=gradual
```