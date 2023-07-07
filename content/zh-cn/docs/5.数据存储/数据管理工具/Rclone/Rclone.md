---
title: "Rclone"
linkTitle: "Rclone"
weight: 1
---

# 概述

> 参考：
>
> - [GitHub 项目，rclone/rclone](https://github.com/rclone/rclone)
> - [官网](https://rclone.org/)

**Rclone** 是一个命令行工具，用来管理云存储上的文件。Rclone 也可以看作 **rsync for cloud storage(用于云存储的 rsync)**。Rclone 支持各种存储类型，包括 商业文件存储服务、标准传输协议(比如 WebDAV、S3 等)、等等。从[这里](https://rclone.org/#providers)我们可以查看到所有受支持的存储提供者

Rclone 将存储提供者抽象为 **Remote**，在我们配置 Rclone 时，经常会看到 Remote 这个词，创建、删除 Remote 这种行为，就是在 [INI](/docs/2.编程/无法分类的语言/INI.md) 格式的配置文件中配置 Remote。这些 Remote 由指定类型的 **Backend** 提供支持。

> 比如，我们可以这样描述，我创建了一个名为 alist 的 Remote，使用的是 WebDav 类型的 Backend。

Rclone 还可以将这些 Remote 作为磁盘挂载在 Windows、macOS、Linux 上，并通过 SFTP、HTTP、WebDAV、FTP、DLNA 对外提供存储能力。

# Rclone 关联文件与配置

**rclone.conf** # 各种 Remotes 信息。

- 如果在某些已定义的位置都没有找到 rclone.conf 文件，则会在以下位置创建一个新的配置文件：
  - Windows 上
    - 在 **$APPDATA/rclone/rclone.conf**
  - 类 Unix 上
    - 如果定义了 `$XDG_CONFIG_HOME`，则在 **$XDG_CONFIG_HOME/rclone/rclone.conf**
    - 如果未定义 `$XDG_CONFIG_HOME`，则在 **~/.config/rclone/rclone.conf**

# Syntax(语法)

> 参考：
>
> - [官方文档，命令](https://rclone.org/commands/)

## 全局标志

> 参考：
>
> - [官方文档，全局标志](https://rclone.org/flags/)

详见 [Rclone 配置](docs/5.数据存储/数据管理工具/Rclone/Rclone%20配置.md)，命令行标志通常也是 Rclone 的配置，分为两种，与后端无关的和与后端相关的。绝大部分情况，在使用命令行时，我们一般使用与后端无关的的标志。在这个文章中，主要看 [Backend 无关配置详解](docs/5.数据存储/数据管理工具/Rclone/Rclone%20配置.md#Backend%20无关配置详解) 部分即可

# rclone config

进入交互式会话，用以修改配置文件(默认为 `~/.config/rclone/rclone.conf`)。进入交互式配置会话中，我们可以设置新的 Remotes 并管理现有 Remotes。还可以设置或删除密码以保护我们的配置。

除了基础的交互式，我们还可以使用各种子命令来直接修改配置文件

## Syntax(语法)

**rclone config \[FLAGS] \[COMMAND]**

**COMMAND**

- **file** # 显示正在使用的配置文件的j路径
- **show** # 打印 (解密) 配置文件，或单个 Remote 的配置。

# rclone copy

将源 Remote 的文件复制到目标 Remote 中，跳过相同的文件。不会删除目标中的比源中多的文件。

## Syntax(语法)

**rclone copy SOURCE:SourcePath DEST:DestPath**

# rclone copyto

copyto 可以在上传单个文件到目标目录下时，改变文件的原名。其他情况与 copy 的功能相同。

## Syntax(语法)

# rclone mount

https://rclone.org/commands/rclone_mount/

将 Remote 作为文件系统挂载到操作系统中

## Syntax(语法)

**rclone mount REMOTE:PATH /PATH/TO/MountPoint \[FLAGS]**

**FLAGS**

VFS 文件缓存相关标志

- **--cache-dir STRING** # 指定用于保存缓存文件的目录。`默认值: %LOCALAPPDATA%\rclone\`
  - `Linux 默认值: ~/.cache/rclone`
- **--vfs-cache-mode STRING** # 缓存模式，可用的值有: off|minimal|writes|full。`默认值: off`
- **--vfs-cache-max-age DURATION** # 缓存中的对象保存的最大时间，超时的将被删除。`默认值: 1h`
- --vfs-cache-max-size SizeSuffix      Max total size of objects in the cache (default off)
- --vfs-cache-poll-interval duration   Interval to poll the cache for stale objects (default 1m0s)
- --vfs-write-back duration            Time to writeback files after last use when using cache (default 5s)


## EXAMPLE

挂载 alist webdav 到本地磁盘

- `rclone mount alist:/ Z: --cache-dir D:\appdata\rclone-cache --vfs-cache-mode full --header Referer:`

# rclone sync

让目标 Remote 与源 Remote 保持相同，仅修改目标 Remote 中的数据。

https://rclone.org/commands/rclone_sync/

注意：由于 sync 命令会导致目标数据丢失，最好使用 --dry-run 或 -i, --interactive 标志进行测试

## Syntax(语法)

**rclone sync SOURCE:PATH DEST:PATH \[FLAGS]**

# 列出 Remote 中的数据相关命令

ls

lsl

lsd

lsf

lsjson

# 应用示例

## 两个对象存储同步数据

**rclone sync** 命令会在源和目标之间同步文件。 它会删除目标目录中源目录没有的文件，并且会更新目标目录中的文件。 **rclone copy** 命令只会在源和目标之间复制文件。 它不会删除目标目录中的文件，也不会更新文件。

## webdav 挂载为电脑本地硬盘(非网络硬盘)

> 原文链接：[B 站-捕梦小达人](https://www.bilibili.com/read/cv13661426)

注意：需要安装 winfsp

使用 Alist 的 阿里云网盘时，注意添加 `--header`，参考 [alist discussions 630](https://github.com/alist-org/alist/discussions/630)

```
rclone mount --config rclone.conf alist:/ Z: --cache-dir D:\appdata\rclone --vfs-cache-mode full --header "Referer:"
```

可以参考 PowerShell 的 [Management](/docs/1.操作系统/Y.Windows%20管理/Windows%20管理工具/PowerShell%20内置管理工具/Management.md) 模块下的 Start-Process 命令的，以便在后台运行，效果如下：

```powershell
Start-Process "alist.exe" -ArgumentList "server --data D:\appdata\alist" -WindowStyle Hidden -RedirectStandardOutput "D:\Tools\Scripts\log\alist.log" -RedirectStandardError "D:\Tools\Scripts\log\alist-err.log"

Start-Process "rclone.exe" `
-ArgumentList "mount alist-net:/ Z: --cache-dir D:\appdata\rclone-cache --vfs-cache-mode full --vfs-cache-max-age 24h --header Referer:" `
-WindowStyle Hidden `
-RedirectStandardOutput "D:\Tools\Scripts\log\rclone.log" -RedirectStandardError "D:\Tools\Scripts\log\rclone-err.log"
```

然后可以参考 Windows 管理中的 [设置开机自启动](/docs/1.操作系统/Y.Windows%20管理/设置开机自启动.md) 以便开机时自动挂载。