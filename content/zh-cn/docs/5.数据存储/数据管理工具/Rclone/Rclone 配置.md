---
title: "Rclone 配置"
linkTitle: "Rclone 配置"
weight: 20
---

# 概述

> 参考：
> 
> - [官方文档，用法](https://rclone.org/docs/)

Rclone 可以通过如下几种方式配置：

- 命令行选项
- 环境变量
- rclone.conf 配置文件（只对 Remote 又有效，无法通过配置文件配置 Rclone 的运行时行为）

rclone 命令的每个选项都可以通过环境变量设置其默认值。选项与环境变量的对应关系规则如下：

- 去掉长选项名称开头的 `--`
- 选项名称中的 `-` 改为 `_`
- 字母改为大写
- 前面加上 `RCLONE_` 前缀

比如：

- --stats 对应 RCLONE_STATS
- --drive-use-trash 对应 RCLONE_DRIVE_USE_TRASH

> 其中 --verbose 的对应关系比较特殊，-v 对应 RCLONE_VERBOSE=1；-vv 对应 RCLONE_VERBOSE=2 

Rclone 的配置有点混乱，不过大体分为两部分，Backend 和 Non Backend，即**与后端相关**和**与后端无关**两种配置。

- **与后端相关** # 由于各种存储系统具有相当复杂的身份验证，因此 Remote 的信息都以配置文件的方式保存到本地 [rclone.conf](#rclone.conf) 文件中
  - TODO: 这些标志唯一的用处是不是只在 `rclone config` 命令中有用？就是管理 Remote 信息？
  - 与后端相关的配置也可以通过命令行标志进行配置，这些命令行标志也是用来生成 rclone.conf 配置文件的，那为何要多此一举呢，直接编辑文件它不香吗？~
- **与后端无关** # 通常是指 Rclone 自身的运行时方式。
  - 这部分配置无法通过 rclone.conf 文件设置。

## 使用方式

想要使用 Rclone，我们需要先创建配置文件，并配置想要管理的 Remote 信息(就是使用 `rclone config` 命令或手动配置 rclone.conf 文件)，在这个文件中配置每个 Remote 的 Backend 信息。有了 Remote 之后，Rclone 就有了控制目标，可以对目标中的数据进行控制了。

## rclone.conf

**rclone.conf** 文件中保存了各种 Remote 信息。可以直接编辑该文件，也可以使用 `rclone config` 命令在命令行的交互模式下进行配置。

rclone 运行后，根据如下顺序从上到下依次查找 rclone.conf 文件

- **rclone 二进制可执行文件所在目录**中的 rclone.conf 文件
- **$APPDATA/rclone/rclone.conf**（该位置仅在 Windows 系统中有效）
- **$XDG_CONFIG_HOME/rclone/rclone.conf**
- **~/.config/rclone/rclone.conf** 
- **~/.rclone.conf**

如果上述位置都没有找到 rclone.conf 文件，则会在以下位置创建一个新的配置文件：

- Windows 上
  - 在 **$APPDATA/rclone/rclone.conf**
- 类 Unix 上
  - 如果定义了 `$XDG_CONFIG_HOME`，则在 **$XDG_CONFIG_HOME/rclone/rclone.conf**
  - 如果未定义 `$XDG_CONFIG_HOME`，则在 **~/.config/rclone/rclone.conf**


# 与后端无关的配置

> 参考：
> 
> - [官方文档，命令行标志，与后端无关的标志](https://rclone.org/flags/#non-backend-flags)

这部分配置无法通过 rclone.conf 文件配置。只能使用环境变量和命令行选项。除了在这里记录的通用配置意外，各个 Rclone 功能(比如 mount 等命令)还有自身的配置用以配置运行时行为。

**--cache-dir** # 用于缓存的目录。`默认值：~/.cache/rclone`

# 与后端有关的配置

> 参考：
> 
> - [官方文档，命令行标志，与后端有关的标志](https://rclone.org/flags/#backend-flags)

这部分配置通常都是指 rclone.conf 文件的配置，我们也可以使用 `rclone config` 命令的选项对 rclone.conf 文件进行配置。

**这部分配置本质就是配置 Remote**，为每个 Remote 配置 Backend，可以用的配置内容非常多，并且每种类型的 Backend 通常都有其独立的配置内容。

## WebDav

https://rclone.org/webdav

