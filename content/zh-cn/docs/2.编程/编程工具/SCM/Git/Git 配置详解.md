---
title: Git 配置详解
---

# 概述

> 参考：
> - [官方 Book，自定义 Git-Git 配置](https://git-scm.com/book/en/v2/Customizing-Git-Git-Configuration)
> - [官方文档，git-config](https://git-scm.com/docs/git-config)

Git 使用一系列

# ~/.gitconfig 文件配置详解

该配置文件可以手动修改，也可以通过 git config --global XXX 命令修改。

```ini
[filter "lfs"]
required = true
clean = git-lfs clean -- %f
smudge = git-lfs smudge -- %f
process = git-lfs filter-process
[user]
name = DesistDaydream
email = XXXXXXXX@qq.com
[core]
autocrlf = input
[credential]
helper = store
```

**git config --global credential.helper store**

# core 部分

**quotePath**(BOOLEAN) # 决定了 git 在控制台中如何显示非 ASCII 的文件名。`默认值：true`

- 如果 core.quotePath 为 true，那么 git 会用八进制的引号表示法来显示文件名，例如 "\344\270\255\345\233\276.txt"。如果 core.quotePath 为 false，那么 git 会直接显示文件名的 UTF-8 编码，例如 "中文.txt"。
- core.quotePath 的设置会影响一些 git 命令的输出，例如 git status, git diff, git log 等。