---
title: Git 配置详解
---

# 概述

# ~./gitconfig 文件配置详解

该配置文件可以手动修改，也可以通过 git config --global XXX 命令修改。

```bash
[filter "lfs"]
	required = true
	clean = git-lfs clean -- %f
	smudge = git-lfs smudge -- %f
	process = git-lfs filter-process
[user]
	name = DesistDaydream
	email = 373406000@qq.com
[core]
	autocrlf = input
[credential]
	helper = store
```

**git config --global credential.helper store**
