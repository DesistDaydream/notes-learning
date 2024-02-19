---
title: WSL

---

# 概述

> 参考：
> 
> - [GitHub 项目，microsoft/WSL](https://github.com/microsoft/WSL)
> - [官方文档，windows-wsl](https://docs.microsoft.com/zh-cn/windows/wsl/)


# 安装 WSL

现在默认使用 WSL2，也推荐安装和使用 WSL2。

打开 “启用或关闭Windows功能”，开启 “适用于 Linux 的 Windows 子系统” 和 “虚拟机平台”。

![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/wsl/20230601094318.png)

之后在 Microsoft Store 中使用 `Windows Susystem for Linux` 关键字搜索并安装 WSL 的最新版。

若不开启“虚拟机平台”  或 安装最新版 WSL，在安装后启动时，将可能会出现下图错误

![wsl-error.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/wsl/wsl-error.png)


## 安装 Linux 发行版

在 PowerShell 执行指令

安装 Ubuntu 发行版的 WSL

```
wsl --install -d Ubuntu
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

## Windows 下的关联文件

**%LOCALAPPDATA%/Packages/CanonicalGroupLimited.XXXXX** # 通过应用商店安装后的程序数据保存路径。比如 CanonicalGroupLimited.Ubuntu_79rhkp1fndgsc

- **./LocalState/ext4.vhdx** # WSL 虚拟机文件，WSL 启动的虚拟机后数据都在该文件中，类似于 kvm/qemu 的 .qcow2 文件

**%UserProfile%/.wslconfig** # 用于在作为 WSL2 版本运行的所有已安装 Linux 发行版中全局配置设置。

## LInux 发行版下的关联文件

**/etc/wsl.conf** # 作为 Unix 文件存储，用于为每个发行版配置各自独立的设置。

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

- **--install \[DISTRIBUTION] [OPTIONS]** # 
- **--shutdown** # 立即终止所有正在运行在 wsl 子系统上的 Linux 发行版
- **--status** # 显示 wsl 子系统的状态。
- **--update** # 更新 wsl 子系统程序包。

发行版管理选项

- **-l, --list [OPTIONS]** # 列出发行版，可以根据子参数指定需要列出哪些发行版。
  - **--all** # 列出所有
  - **-o, --online** # 列出所有可以安装的发行版。
- **--set-default-version** # 
- **-s, --set-default DISTRO** # 将指定的发行版设为默认
- **--unregister DISTRO** # 将指定的发行版取消注册。
  - <font color="#ff0000">注意</font>：若从应用商店删除特定发行版后再安装失败的话，需要通过 wsl 命令手动 unregister 一下，即可成功。

