---
title: Docker CLI Plugin
linkTitle: Docker CLI Plugin
date: 2024-01-27T19:45
weight: 20
---

# 概述

> 参考：
> 
> - https://github.com/docker/cli/issues/1534


# 关联文件与配置


- Unix-like OSes :
  - `$HOME/.docker/cli-plugins`
  - `/usr/local/lib/docker/cli-plugins` & `/usr/local/libexec/docker/cli-plugins`
  - `/usr/lib/docker/cli-plugins` & `/usr/libexec/docker/cli-plugins`
- On Windows:
  - `%USERPROFILE%\.docker\cli-plugins`
  - `C:\ProgramData\Docker\cli-plugins`