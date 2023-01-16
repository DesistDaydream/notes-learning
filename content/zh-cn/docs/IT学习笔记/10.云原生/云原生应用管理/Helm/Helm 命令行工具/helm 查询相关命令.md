---
title: helm 查询相关命令
---

# helm get # 获取指定 release 的信息。

该命令由多个子命令组成，这些子命令可用于获取有关 release 的扩展信息，包括：

1. 用于生成 release 的值
2. 生成 release 的 manifest 文件
3. 生成 release 的 chart 的注释信息
4. 与 release 相关的 hooks。

**helm get \[COMMAND]**

该命令与 kubectl get XXX XXX -o yaml 效果类似，可用的 COMMAND 为二级标题

## all # download all information for a named release

## hooks # 获取指定 release 的所有 hooks

**helm get hooks RELEASE_NAME \[FLAGS]**

## manifest # 获取指定 release 的 manifest 文件

**helm get manifest RELEASE_NAME \[FLAGS]**

FLAGS

1. \--revision INT # 指定一个 release 的版本

EXAMPLE

1. helm get manifest myapp # 获取 myapp 这个 release 的所有 manifest 文件。

## notes # download the notes for a named release

## values # download the values file for a named release
