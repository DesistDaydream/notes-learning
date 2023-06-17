---
title: WSL

---

# 概述

> 参考：
> 
> - [官方文档，windows-wsl](https://docs.microsoft.com/zh-cn/windows/wsl/)


# 安装 WSL

现在默认使用 WSL2，也推荐安装和使用 WSL2。

打开 “启用或关闭Windows功能”，开启 “适用于 Linux 的 Windows 子系统” 和 “虚拟机平台”。

![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/wsl/20230601094318.png)

若不开启“虚拟机平台”，在安装后启动时，将可能会出现下图错误

![wsl-error.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/wsl/wsl-error.png)

## 安装 Linux 发行版

在 PowerShell 执行指令

安装 Ubuntu 发行版的 WSL

```
wsl --install -d Ubunt
```

## 常见问题

若安装后 linux 无法启动，报错：`WslRegisterDistribution failed with error: 0x800701bc`

（可选）设置 wsl 版本

```shell
wsl --set-default-version 2
```

忘记密码时，可以在 PowerShell 中使用 wsl 命令直接以 root 用户登录 wsl

```shell
wsl.exe --user root
```

更多常见问题，可以参考[官方文档，疑难解答](https://learn.microsoft.com/zh-cn/windows/wsl/troubleshooting)

# 访问 WSL 文件系统

在 Windows 资源管理器中，访问 `\\wsl$` 即可访问 WSL 文件系统

![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/tqwpkc/1654930585949-71f955ca-97c4-45d8-be77-a637670803eb.png)


# WSL 关联文件与配置

Windows 下的关联文件

- `%UserProfile%/.wslconfig` # - 用于在作为 WSL 2 版本运行的所有已安装 Linux 发行版中全局配置设置。

LInux 发行版下的关联文件

- **/etc/wsl.conf** # 

# wsl 命令行工具

> 参考：
> 
> - [官方文档-WSL，概述-基本 WSL 命令](https://learn.microsoft.com/zh-cn/windows/wsl/basic-commands)

## Syntax(语法)

**wsl [OPTOINS]**

**OPTONS**

- **-d, --distribution STRING** # 指定要运行的发行版
- **-u, --user STRING** # 使用指定的用户运行发行版。可以直接以 root 用户运行。常用来在忘记密码时候使用。

WSL 子系统管理选项

- **--status** # 显示适用于 Linux 的 Windows 子系统的状态。
- **--update** # 更新适用于 Linux 的 Windows 子系统程序包。

发行版管理选项

- **-l, --list [OPTIONS]** # 列出发行版，可以根据子参数指定需要列出哪些发行版。
  - **--all** # 列出所有
  - **-o, --online** # 列出所有可以安装的发行版。
- **--set-default-version** # 
- **-s, --set-default DISTRO** # 将指定的发行版设为默认

# WSL 配置网络

## 配置桥接网络和静态 IP

> 参考：
> - [博客园，WSL2使用桥接网络，并指定IP](https://www.cnblogs.com/lic0914/p/17003251.html)
>   - 该文章参考的原文: https://github.com/luxzg/WSL2-fixes/blob/master/networkingMode%3Dbridged.md

**TODO: 这种做法会导致 WSL 无法和 Windows 主机通信，原理暂时未知。**

使用 Hyper-V 创建一个桥接类型的网络，假如命名为 **br0**

在 `%UserProfile%/.wslconfig` 文件中添加如下内容：

```ini
[wsl2]
networkingMode=bridged
vmSwitch=br0
dhcp=false
ipv6=true
```

然后再启动 WSL，即可使用桥接模式的网络。

此时，还需要为 WSL 启用 Systemd，否则系统中没有可以配置 IP 的程序

```
sudo tee /etc/wsl.conf <<EOF
[boot]
systemd=true
EOF
```

然后需要通过 systemd-network 服务配置 IP，那么还需要进行一些配置：

```bash
sudo tee /lib/systemd/network/wsl_external.network <<EOF
[Match]
Name=eth0
[Network]
Description=bridge
DHCP=false
Address=192.168.254.254/24 # 自行修改
Gateway=192.168.254.1 # 自行修改
EOF

sudo systemctl enable systemd-networkd --now
```

然后重启 wsl 即可。

