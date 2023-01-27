---
title: WSL
---

# 概述

> 参考：
> - [官方文档,windows-wsl](https://docs.microsoft.com/zh-cn/windows/wsl/)


# 安装 WSL

现在默认使用 WSL2，也推荐安装和使用 WSL2。

打开 “启用或关闭Windows功能”，开启 “适用于 Linux 的 Windows 子系统”

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/0-picgo/20230126182502.png)

## 安装 Linux 发行版

在 PowerShell 执行指令

安装

```
wsl --install
```

## 常见问题
若安装后 linux 无法启动，报错：WslRegisterDistribution failed with error: 0x800701bc，可以参考[官方文档，疑难解答](https://learn.microsoft.com/zh-cn/windows/wsl/troubleshooting)

（可选）设置 wsl 版本

```shell
wsl --set-default-version 2
```

忘记密码时，可以在 PowerShell 中使用 wsl 命令直接以 root 用户登录 wsl

```shell
wsl.exe --user root
```

# 为 WSL2 设置代理

```bash
#!/bin/bash
#
export hostip=$(cat /etc/resolv.conf |grep -oP '(?<=nameserver\ ).*')
export http_proxy="http://${hostip}:7890"
export https_proxy="http://${hostip}:7890"
export all_proxy="sock5://${hostip}:7890"
export ALL_PROXY="sock5://${hostip}:7890"
```

# 访问 WSL 文件系统

在 Windows 资源管理器中，访问 `\\wsl$` 即可访问 WSL 文件系统

![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/tqwpkc/1654930585949-71f955ca-97c4-45d8-be77-a637670803eb.png)


# WSL 关联文件与配置


