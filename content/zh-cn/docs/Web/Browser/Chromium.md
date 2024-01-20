---
title: Chromium
linkTitle: Chromium
date: 2024-01-06T14:49
weight: 2
---

# 概述

> 参考：
>
> - [源码](https://chromium.googlesource.com/chromium/src/)
>   - [GitHub 项目，chromium/chromium](https://github.com/chromium/chromium)
> - [Chromium 官网](https://www.chromium.org/Home/)
> - [Chromium 文档源码](https://chromium.googlesource.com/chromium/src/+/refs/heads/main/docs)
>   - [文档](https://chromium.googlesource.com/chromium/src/+/HEAD/docs/README.md)

https://github.com/chromium/permission.site 用于测试 Web API 和浏览器权限交互的站点。比如通过浏览器调用位置信息、蓝牙、等等。

# 用户数据

https://chromium.googlesource.com/chromium/src/+/HEAD/docs/user_data_dir.md

Chrome 运行产生的用户数据包含 Profile 数据、运行时状态数据。Chrome 可以支持多人同时使用

- Profile 数据则是特定于某个具体用户的数据，包括 历史记录、书签、Cookie、扩展程序、等等。
- 每个 Profile 数据所在位置都是用户数据目录的一个子目录。

保存用户数据的目录根据不同环境，有不同的默认值，一般来说，取决于如下几点：

- 操作系统。Linux、Windows、Macos 等等
- 基于 Chromium 的各种品牌。Chrome、Edge、等等
- Release 版本。比如 stable、beta、dev、canary、等等。

不同系统下的默认路径详见官网，这里就不写了，在 [Chrome](/docs/Web/Browser/Chrome.md) 里有概述

可以通过 --user-data-dir 命令行标志改变用户数据目录的位置，通过 --profile-directory 命令行标志改变启动 Chrome 时要使用具体哪个用户运行。
