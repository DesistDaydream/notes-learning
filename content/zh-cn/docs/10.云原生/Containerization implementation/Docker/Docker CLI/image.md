---
title: image
linkTitle: image
weight: 20
---

# 概述

> 参考：

image 命令用于管理容器镜像。

# build - 使用 Dockerfile 构建一个镜像

# history - Show the history of an image

# import - Import the contents from a tarball to create a filesystem image

# inspect - Display detailed information on one or more images

# load - Load an image from a tar archive or STDIN

# ls - 列出所有镜像

# prune - 移除未使用的镜像

## Syntax(语法)

**docker image prune \[OPTIONS]**

OPTIONS

- **-a, --all** # Remove all unused images, not just dangling ones
- **-f, --force** # Do not prompt for confirmation

## EXAMPLE

docker image prune -a # 清理所有没有使用的镜像

# pull - Pull an image or a repository from a registry

# push - Push an image or a repository to a registry

# rm - Remove one or more images

# save - Save one or more images to a tar archive (streamed to STDOUT by default)

# tag - Create a tag TARGET_IMAGE that refers to SOURCE_IMAGE
