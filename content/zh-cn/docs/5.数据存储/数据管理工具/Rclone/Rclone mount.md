---
title: "Rclone mount"
linkTitle: "Rclone mount"
weight: 20
---

# 概述

> 参考：
>
> - [官方文档，命令-rclone mount](https://rclone.org/commands/rclone_mount/)

Rclone 的 mount 功能使用 FUSE 将 Remote 作为文件系统挂载到操作系统中。

## 使用 VFS

https://rclone.org/commands/rclone_mount/#vfs-virtual-file-system

rclone mount 命令使用文件系统的 VFS，这将 rclone 使用 Remote 调整为看起来更像磁盘文件系统的东西。

### VFS 缓存机制

rclone mount 可以使用两种 VFS 的缓存机制

- 目录缓存
- 文件缓存

**文件缓存**

为 --vfs-cache-mode 选项指定了除 off 以外的值，即表示开启文件缓存。开启文件缓存后，默认是 write-back 策略，对文件的所有<font color="#ff0000">写入操作都是针对本地的缓存文件</font>，只有<font color="#ff0000">关闭缓存文件后等待 5s</font>(该值可由 --vfs-write-back 选项指定)后才会将缓存文件回写到 Remote 以更新原文件。如果 rclone 因为意外情况退出导致文件未上传，那么必须要使用相同选项启动 rclone，才会在开始运行时回写这些文件。

> 用人话说：开启文件缓存后，修改的文件是都是本地缓存文件，需要等一会，才会缓存文件同步给源文件。

rclone 提供了 4 中缓存模式，模式越高，兼容性越好，但磁盘占用越多。

- **off** # 关闭文件缓存。最低的模式，兼容性最差。本地不做任何缓存，所有文件直接从云端获取并写入，网速特别好时（复制粘贴大文件时建议至少100M管以上速度）使用。
- **minimal** # 最小模式。与 off 类似，一般只有新创建的文件才会缓存。
- **writes** # 写模式。打开后有修改的文件都会缓存。如果回写失败，将以指数增加的间隔重试，最多 1 分钟。
- **full** # 完全模式，兼容性最好。缓存全部文件。
  - 在此模式下，缓存中的文件将是 **sparse files(稀疏文件)**，并且 rclone 将跟踪已下载的文件的哪些位。

notes: 对于在 Windows 和 Linux 中的 rclone mount，上述缓存模式的说明不太一样，比如关闭缓存的情况下，linux 用 vim 是无法编辑文件的，但是 windows 中用编辑器是可以编辑的，并且修改内容瞬间会同步到源文件。

# Syntax(语法)

**rclone mount REMOTE:PATH /PATH/TO/MountPoint \[FLAGS]**

## FLAGS

- **--network-mode** # (仅限于 Windows 系统)挂载为远程网络驱动器，而不是磁盘驱动器。
  - 网络驱动器可以在网络邻居中显示，方便访问和共享。
  - 网络驱动器可以避免一些 Windows 对本地驱动器的限制，如文件名长度和字符集。

VFS 性能相关标志

- **--no-modtime** # 不要读/写修改时间（可以加快速度）。
  - 特别是 S3 和 Swift 从 --no-modtime 标志中受益匪浅（或使用 --use-server-modtime 以获得稍微不同的效果），因为每次读取修改时间都需要一个事务。

VFS 文件缓存相关标志

- **--dir-cache-time duration** # Time to cache directory entries for。`默认值: 5m0s`
- **--poll-interval duration** # Time to wait between polling for changes. Must be smaller than dir-cache-time. Only on supported remotes. Set to 0 to disable。`默认值: 1m0s`

VFS 文件缓存相关标志

- **--cache-dir STRING** # 指定用于保存缓存文件的目录。`
  - Linux `默认值: ~/.cache/rclone/`
  - Windows `默认值: %LOCALAPPDATA%/rclone/`
- **--vfs-cache-mode STRING** # 缓存模式。`默认值: off`
  - 可用的值有: off | minimal | writes | full
- **--vfs-cache-max-age DURATION** # 缓存中的对象保存的最大时间，超时的将被删除。`默认值: 1h`
- **--vfs-cache-max-size SizeSuffix** # 缓存占用的最大空间。`默认值: off`
  - 请注意，缓存可能由于两个原因而超出此大小。首先，因为它仅在每个 --vfs-cache-poll-interval 期间进行检查。其次，因为打开的文件无法从缓存中逐出，当超过 --vfs-cache-max-size 时，rclone 将尝试首先从缓存中逐出访问次数最少的文件，将从最长时间未被访问的文件开始。这种缓存刷新策略非常有效，并且更相关的文件可能会保留在缓存中。
- **--vfs-cache-poll-interval duration** # 轮询缓存以查找陈旧对象的时间间隔。`默认值: 1m0s` (default 1m0s)
- **--vfs-write-back DURATION** #  关闭文件后，将文件回写到 Remote 的等待时间。`默认值: 5s`

# 最佳实践

挂载 alist webdav 到本地磁盘

- `rclone mount alist-net:/ Z: --cache-dir D:\appdata\rclone-cache --vfs-cache-mode writes --vfs-cache-max-age 30d --vfs-cache-max-size 100G --vfs-cache-poll-interval 10m --no-modtime --header Referer: -v`
  - 让**文件夹视图**变为**小图标**可以避免打开文件夹时缓存全部文件，尤其是对于存图片的文件夹来说，可以极大得提高打开速度
  - 看视频的话，应该不用把视频缓存下来~ writes 模式就够了
