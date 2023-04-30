---
title: Rclone
---

# 概述

> 参考：
> 
> - [GitHub 项目，rclone/rclone](https://github.com/rclone/rclone)
> - [官网](https://rclone.org/)

Rclone 是一个命令行工具，用来管理云存储上的文件。Rclone 也可以看作 rsync for cloud storage(用于云存储的 rsync)。Rclone 支持各种存储类型，包括 商业文件存储服务、标准传输协议(比如 WebDAV、S3 等)、等等。从[这里](https://rclone.org/#providers)我们可以查看到所有受支持的存储提供者
Rclone 将存储提供者抽象为 Remote，在我们配置时，经常会看到 remote 这个词，创建、删除、等行为一个 remote，就是在配置文件中配置 remote~~
Rclone 还可以将这些 remote 作为磁盘挂载在 Windows、macOS、Linux 上，并通过 SFTP、HTTP、WebDAV、FTP、DLNA 对外提供存储能力。

## Rclone 关联文件与配置

**~/.config/rclone/rclone.conf** # 保存各种 Remotes 信息的配置文件

## Syntax(语法)

> 参考：
> 
> - [官方文档，命令](https://rclone.org/commands/)

### 全局标志

> 参考：
> 
> - [官方文档，全局标志](https://rclone.org/flags/)

**-n, --dry-run** # 试运行，不会真的执行
**-i, --interactive** # 开启交互模式
**-p, --progress** # 显示传输进度、传输速度

# rclone config

进入交互式会话，用以修改配置文件(默认为 ~/.config/rclone/rclone.conf)。进入交互式配置会话中，我们可以设置新的 Remotes 并管理现有 Remotes。还可以设置或删除密码以保护我们的配置。
除了基础的交互式，我们还可以使用各种子命令来直接修改配置文件

## Syntax(语法)

**rclone config \[FLAGS] \[COMMAND]**

# rclone copy

## Syntax(语法)

**rclone copy SOURCE:SourcePath DEST:DestPath**

# rclone copyto

copyto 可以在上传单个文件到目标目录下时，改变文件的原名。其他情况与 copy 的功能相同。

## Syntax(语法)

# rclone mount

将 Remote 作为文件系统挂载到操作系统中

## Syntax(语法)

**rclone mount REMOTE:PATH /PATH/TO/MountPoint \[FLAGS]**

# rclone sync

注意：由于 sync 命令会导致目标数据丢失，最好使用 --dry-run 或 -i, --interactive 标志进行测试：`rclone sync -i SOURCE remote:DESTINATION`

## Syntax(语法)

**rclone sync SOURCE:PATH DEST:PATH \[FLAGS]**

# 应用示例

## webdav 挂载为电脑本地硬盘(非网络硬盘)

> 原文链接：[B 站-捕梦小达人](https://www.bilibili.com/read/cv13661426)

注意：需要安装 winfsp

使用 Alist 的 阿里云网盘时，注意添加 `--header`，参考 [alist discussions 630](https://github.com/alist-org/alist/discussions/630)

```
rclone mount --config rclone.conf alist:/ Z: --cache-dir D:\app_data\rclone --vfs-cache-mode full --header "Referer:"
```