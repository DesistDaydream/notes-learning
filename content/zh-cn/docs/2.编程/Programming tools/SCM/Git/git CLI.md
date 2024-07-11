---
title: git CLI
linkTitle: git CLI
date: 2024-02-13T08:48
weight: 20
---

# 概述

> 参考：
>
> - https://git-scm.com/docs/git

git 工具通过多个子命令来使用，可以按照功能对这些命令进行划分

- [设置与配置 Git](#设置与配置)
- [获取和创建项目](#获取和创建项目)
- [分支与合并](#分支与合并)
- [分享和更新项目](#分享和更新项目)

# 设置与配置

## config - 配置 git

获取和设置存储库或全局选项

### Syntax(语法)

```bash
git config [<file-option>] [--type=<type>] [--fixed-value] [--show-origin] [--show-scope] [-z|--null] <name> [<value> [<value-pattern>]]
git config [<file-option>] [--type=<type>] --add <name> <value>
git config [<file-option>] [--type=<type>] [--fixed-value] --replace-all <name> <value> [<value-pattern>]
git config [<file-option>] [--type=<type>] [--show-origin] [--show-scope] [-z|--null] [--fixed-value] --get <name> [<value-pattern>]
git config [<file-option>] [--type=<type>] [--show-origin] [--show-scope] [-z|--null] [--fixed-value] --get-all <name> [<value-pattern>]
git config [<file-option>] [--type=<type>] [--show-origin] [--show-scope] [-z|--null] [--fixed-value] [--name-only] --get-regexp <name-regex> [<value-pattern>]
git config [<file-option>] [--type=<type>] [-z|--null] --get-urlmatch <name> <URL>
git config [<file-option>] [--fixed-value] --unset <name> [<value-pattern>]
git config [<file-option>] [--fixed-value] --unset-all <name> [<value-pattern>]
git config [<file-option>] --rename-section <old-name> <new-name>
git config [<file-option>] --remove-section <name>
git config [<file-option>] [--show-origin] [--show-scope] [-z|--null] [--name-only] -l | --list
git config [<file-option>] --get-color <name> [<default>]
git config [<file-option>] --get-colorbool <name> [<stdout-is-tty>]
git config [<file-option>] -e | --edit
```

**OPTIONS**

- **-l, --list** # 列出配置文件中设置的所有变量及其值。

### EXAMPLE

diff Show changes between commits, commit and working tree, etc

fetch Download objects and refs from another repository

grep Print lines matching a pattern

init - 创建一个空的 Git 存储库或重新初始化现有的存储库

# 获取和创建项目
## clone - 将一个存储库克隆到一个新的目录

OPTIONS

- **--branch,-b** # 指定名为 NAME 的分支

EXAMPLE

- git clone -b v1.0 XXXX # 克隆 v1.0 分支的代码

commit Record changes to the repository

# 分支与合并

## log - 展示所有 commit 的记录。默认展示当前分支



EXAMPLE

- git log -p -2 prometheus-rules.yaml # 查看 prometheus-rules.yaml 文件最近两次的修改记录

merge Join two or more development histories together

mv Move or rename a file, a directory, or a symlink

pull Fetch from and merge with another repository or a local branch

push Update remote refs along with associated objects

rebase Forward-port local commits to the updated upstream head

reset Reset current HEAD to the specified state

rm Remove files from the working tree and from the index

show Show various types of objects

status Show the working tree status

tag Create, list, delete or verify a tag object signed with GPG


## tag - 管理仓库的 Tag 信息

git tag -d v0.7.0 删除 v0.7.0 这个 Tag

# 分享和更新项目

## fetch

https://git-scm.com/docs/git-fetch

从仓库中下载对象和引用

## pull

https://git-scm.com/docs/git-pull

## push

https://git-scm.com/docs/git-push

### Syntax(语法)

**git push REPOSITORY**

- REPOSITORY # 

**OPTIONS**

- **-f, --force** # 强制覆盖远程仓库的代码，即使本地代码与仓库代码冲突。
  - Notes: 有的远程仓库具有保护分支功能，会阻止接收通过 --force 传过来的强制覆盖请求，比如 [GitLab 的受保护的分支](https://docs.gitlab.com/ee/user/project/protected_branches.html)、etc.

## remote

https://git-scm.com/docs/git-remote/en

可以在本地项目中添加多个远程仓库（称为 Remote），将代码推送到多个远程仓库中，或者从远程仓库拉取代码。

COMMAND

- **add** # 为本地存储库添加一个远程存储库
- **set-url** # 更改 Remote 的 URL

`git remote add [-t <branch>] [-m <master>] [-f] [--[no-]tags] [--mirror=(fetch|push)] <name> <URL>`

`git remote set-head <name> (-a | --auto | -d | --delete | <branch>)`

