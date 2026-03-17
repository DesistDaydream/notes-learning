---
title: UV
linkTitle: UV
created: 2026-03-17T08:21
weight: 100
---

# 概述

> 参考：
>
> - [GitHub 项目，astral-sh/uv](https://github.com/astral-sh/uv)

只需要一个二进制文件即可使用！具有缓存管理依赖库的功能，完美解决依赖多版本、虚拟环境重复安装依赖占用空间 的问题。

uv init 初始化项目，生成 pyproject.toml 文件

uv add XXX 添加依赖库，同时更新 pyproject.toml 中的依赖部分

uv sync 根据 pyproject.toml 同步依赖

uv tool install XXX 以工具的形式安装 XXX。不作为依赖库

使用 --index 指定镜像源

uv sync --index "https://mirrors.aliyun.com/pypi/simple/"

# uv 关联文件与配置

缓存的储存目录

- Windows: `%LOCALAPPDATA%\uv\cache\`
- Unix: `$XDG_CACHE_HOME/uv/` 或 `$HOME/.cache/uv/`

缓存的储存目录可以通过如下几种方式修改：

- `UV_CACHE_DIR` 环境变量
- pyproject.toml 文件中的 `tool.uv.cache-dir` 键
- uv CLI 的 `--cache-dir` 标志

# 缓存

https://docs.astral.sh/uv/reference/settings/#link-mode

UV 提供四种缓存模式

- clone # 将 Packet 从缓存克隆（i.e. 写时复制）到目标
- copy # 将 Packet 从缓存复制到目标
- hadrlink # 将 Packet 从缓存硬连接到目标
- symlink # 用 [Symbolic link](/docs/1.操作系统/Kernel/Filesystem/文件管理/Symbolic%20link.md) 将 Packet 从缓存软链接到目标

[Unix-like OS](/docs/1.操作系统/Operating%20system/Unix-like%20OS/Unix-like%20OS.md) 上默认使用 clone 模式，[Microsoft OS](/docs/1.操作系统/Operating%20system/Microsoft%20OS/Microsoft%20OS.md) 上默认使用 hardlink 模式。官方不建议使用 symlink 模式

这四种模式，只有 copy 是会消耗双倍磁盘空间的，其他三种方式，基本上只有缓存中的文件才会占用磁盘空间。

> [!Note]
>
> 从 [文件管理](/docs/1.操作系统/Kernel/Filesystem/文件管理/文件管理.md) 可知，hardlink 要求两个文件必须在同一个文件系统下。也就是说 Windows 下，**“Python 项目” 与 “UV 缓存” 必须在同一个分区中**。否则会出现提示: `warning: Failed to hardlink files; falling back to full copy.`

缓存模式配置方式

- 环境变量 `UV_LINK_MODE`
- pyproject.toml 文件中的 `tool.uv.link-mode` 键
- CLI 使用 `--link-mode` 参数

其他

uv cache clean 清除缓存

`uv cache prune` 会删除所有未使用的缓存条目


