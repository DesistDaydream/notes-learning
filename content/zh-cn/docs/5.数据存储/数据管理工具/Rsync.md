---
title: "Rsync"
linkTitle: "Rsync"
weight: 20
---

# 概述

> 参考：
> 
> - [GitHub 项目，WayneD/rsync](https://github.com/WayneD/rsync)
> - [官网](https://rsync.samba.org/)
> - [rsync+inotify 数据实时同步介绍与 K8s 实战应用](https://mp.weixin.qq.com/s/VxnDEQ8e3yQOLJi0JwtyjA)

rsync（remote sync） 远程同步，rsync 是 linux 系统下的数据镜像备份工具。使用快速增量备份工具 Remote Sync 可以远程同步，支持本地复制，或者与其他 SSH、rsync 主机同步。已支持跨平台，可以在 Windows 与 Linux 间进行数据同步。rsync 监听端口：873，rsync 运行模式：C/S。

## 路径结尾带不带 / 的区别

> [!Tip]
> 这里说的通常指 **源路径** 加不加 `/`；目标路径加不加不会有任何影响，因为肯定是同步到指定的目录下面

路径结尾 **加** `/` 类似于 `/*`，将 **源目录下的所有文件** 同步到目标目录中；路径结尾 **不加** `/`，将 **源目录** 统统不到目标目录下。

---

示例

假设 `/src/foo/` 下有 `a.txt`、`b.txt`：

| 命令                              | 结果                                |
| ------------------------------- | --------------------------------- |
| `rsync -av /src/foo /dst/`      | `/dst/foo/a.txt`、`/dst/foo/b.txt` |
| `rsync -av /src/foo/ /dst/`     | `/dst/a.txt`、`/dst/b.txt`         |
| `rsync -av /src/foo/ /dst/foo/` | `/dst/foo/a.txt`、`/dst/foo/b.txt` |

# Syntax



# 最佳实践

指定 SSH 的端口和用户。将 observability 目录及其文件同步到 192.168.254.253 的 /home/DesistDaydream/ 目录下。使用 SSH 的 1070 端口和 command 用户

```bash
rsync -a -e "ssh -p 1070" observability command@192.168.254.253:/home/DesistDaydream/
```

## 基本同步

```bash
rsync -av /home/ /mnt/tmp_home/
```

**参数说明：**

- `-a` — 归档模式，保留权限、时间戳、符号链接、属主等
- `-v` — 显示详细输出

**常用可选参数：**

- `--delete` — 删除目标目录中源目录没有的文件（严格镜像）
- `--progress` — 显示每个文件的传输进度
- `-n` / `--dry-run` — 模拟运行，不实际操作，用于预览

> [!Attention] `/home/`（结尾带斜杠）表示同步目录内的内容到 `/mnt/tmp_home/`，而非在目标下创建 `home` 子目录。
