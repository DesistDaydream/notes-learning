---
title: Git 命令行工具
linkTitle: Git 命令行工具
date: 2024-02-13T08:48
weight: 20
---

# 概述

> 参考：
>
> -

git 工具通过多个子命令来使用

# clone - 将一个存储库克隆到一个新的目录

OPTIONS

- **--branch,-b** # 指定名为 NAME 的分支

EXAMPLE

- git clone -b v1.0 XXXX # 克隆 v1.0 分支的代码

commit Record changes to the repository

# config - 配置 git

获取和设置存储库或全局选项

## Syntax(语法)

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

## EXAMPLE

diff Show changes between commits, commit and working tree, etc

fetch Download objects and refs from another repository

grep Print lines matching a pattern

init - 创建一个空的 Git 存储库或重新初始化现有的存储库

# log - 展示所有 commit 的记录。默认展示当前分支

git log \[] \[] \[\[--] ...]

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

# 最佳实践

git 放弃本地修改，强制拉取更新

- git fetch --all # 指令是下载远程仓库最新内容，不做合并
- git reset --hard origin/master # 指令把 HEAD 指向 master 最新版本
- git pull # 可以省略

本地所有修改，没有提交的文件，都返回到原来的状态

- git checkout . #

提交修改并上传代码

- git add . #
- git commit -a -m 'XXXX 修改' #
- git push #

git 回滚到之前某一 commit

- git log # 查看所有 commit 记录，记录其中的 commit 号，比如 commit 号为：d07466766d46710e54a627f913eea5661382331a
- git reset --hard d07466766d46710e54a627f913eea5661382331a # 恢复到这次 commit 的状态

## 修改 git commit 信息

对自己的提交进行修改

- git commit --amend

将修改强制提交，覆盖原先的提交内容

- git push -f

## 撤销已 push 的 commit

查看所有 commit 记录

```bash
~]# git log
commit 3c15aad50ed125938bbedbedffed05a7b9d600da (HEAD -> main, origin/main)
Author: DesistDaydream <373406000@qq.com>
Date:   Tue May 17 14:52:27 2022 +0800

    Update e37-exporter-workflows.yml

commit 8251ddbbbffe59e5ddb25bb70874425012aa035e
Author: DesistDaydream <373406000@qq.com>
Date:   Mon May 16 17:15:07 2022 +0800
```

回退到指定的版本

```bash
~]# git reset --hard 8251ddbbbffe59e5ddb25bb70874425012aa035e
HEAD is now at 8251ddb Update e37-exporter-workflows.yml
```

强制推送，覆盖远端的版本信息

```bash
~]# export BRANCH="main"
~]# git push origin ${BRANCH} --force
Total 0 (delta 0), reused 0 (delta 0)
To https://github.com/DesistDaydream/e37-exporter.git
 + 3c15aad...8251ddb main -> main (forced update)
```

# 命令详解

**提交相关**

前面我们提到过，想要对代码进行提交必须得先加入到暂存区，Git 中是通过命令 add 实现。

添加某个文件到暂存区：

git add  文件路径

添加所有文件到暂存区：

git add .

同时 Git 也提供了撤销工作区和暂存区命令。

撤销工作区改动：

git checkout --  文件名

清空暂存区：

git reset HEAD  文件名

提交：

将改动文件加入到暂存区后就可以进行提交了，提交后会生成一个新的提交节点，具体命令如下：

git commit -m "该节点的描述信息"

**分支相关**

创建分支：

创建一个分支后该分支会与 HEAD 指向同一节点，说通俗点就是 HEAD 指向哪创建的新分支就指向哪，命令如下：

git branch  分支名

切换分支：

当切换分支后，默认情况下 HEAD 会指向当前分支，即 HEAD 间接指向当前分支指向的节点。

git checkout  分支名

同时也可以创建一个分支后立即切换，命令如下：

git checkout -b  分支名

删除分支：

为了保证仓库分支的简洁，当某个分支完成了它的使命后应该被删除。比如前面所说的单独开一个分支完成某个功能，当这个功能被合并到主分支后应该将这个分支及时删除。

删除命令如下：

git branch -d  分支名

**合并相关**

关于合并的命令是最难掌握同时也是最重要的。我们常用的合并命令大概有三个 merge、rebase、cherry-pick。

merge：

merge 是最常用的合并命令，它可以将某个分支或者某个节点的代码合并至当前分支。具体命令如下：

git merge  分支名/节点哈希值

如果需要合并的分支完全领先于当前分支，如图 3 所示：

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/yc1he6/1622380437243-d5cc7b99-b228-4972-a253-f048e1378e68.webp)

由于分支 ft-1 完全领先分支 ft-2 即 ft-1 完全包含 ft-2，所以 ft-2 执行了“git merge ft-1”后会触发 fast forward(快速合并)，此时两个分支指向同一节点，这是最理想的状态。但是实际开发中我们往往碰到的是下面这种情况：如图 4（左）。

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/yc1he6/1622380437303-d391e621-2e78-48d4-8500-77e2a09964b3.webp)
图 4

这种情况就不能直接合了，当 ft-2 执行了“git merge ft-1”后 Git 会将节点 C3、C4 合并随后生成一个新节点 C5，最后将 ft-2 指向 C5 如图 4（右）。

注意点：如果 C3、C4 同时修改了同一个文件中的同一句代码，这个时候合并会出错，因为 Git 不知道该以哪个节点为标准，所以这个时候需要我们自己手动合并代码。

rebase：

rebase 也是一种合并指令，命令行如下：

git rebase  分支名/节点哈希值

与 merge 不同的是 rebase 合并看起来不会产生新的节点（实际上是会产生的，只是做了一次复制），而是将需要合并的节点直接累加，如图 5。

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/yc1he6/1622380437264-b4959853-9005-408b-b49b-10745896aca3.webp)
图 5

当左边示意图的 ft-1.0 执行了 git rebase master 后会将 C4 节点复制一份到 C3 后面，也就是 C4'，C4 与 C4'相对应，但是哈希值却不一样。

rebase 相比于 merge 提交历史更加线性、干净，使并行的开发流程看起来像串行，更符合我们的直觉。既然 rebase 这么好用是不是可以抛弃 merge 了？其实也不是了，下面我罗列一些 merge 和 rebase 的优缺点：

merge 优缺点：

- 优点：每个节点都是严格按照时间排列。当合并发生冲突时，只需要解决两个分支所指向的节点的冲突即可
- 缺点：合并两个分支时大概率会生成新的节点并分叉，久而久之提交历史会变成一团乱麻

rebase 优缺点：

- 优点：会使提交历史看起来更加线性、干净
- 缺点：虽然提交看起来像是线性的，但并不是真正的按时间排序，比如图 3-3 中，不管 C4 早于或者晚于 C3 提交它最终都会放在 C3 后面。并且当合并发生冲突时，理论上来讲有几个节点 rebase 到目标分支就可能处理几次冲突

对于网络上一些只用 rebase 的观点，作者表示不太认同，如果不同分支的合并使用 rebase 可能需要重复解决冲突，这样就得不偿失了。但如果是本地推到远程并对应的是同一条分支可以优先考虑 rebase。所以我的观点是 根据不同场景合理搭配使用 merge 和 rebase，如果觉得都行那优先使用 rebase。

cherry-pick：

cherry-pick 的合并不同于 merge 和 rebase，它可以选择某几个节点进行合并，如图 6。

命令行：

git cherry-pick  节点哈希值

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/yc1he6/1622380437213-1e90e77b-8877-4449-a6b7-1dfe7d760ddb.webp)
图 6

假设当前分支是 master，执行了 git cherry-pick C3(哈希值)，C4(哈希值)命令后会直接将 C3、C4 节点抓过来放在后面，对应 C3'和 C4'。

**回退相关**

分离 HEAD：

在默认情况下 HEAD 是指向分支的，但也可以将 HEAD 从分支上取下来直接指向某个节点，此过程就是分离 HEAD，具体命令如下：

git checkout  节点哈希值

//也可以直接脱离分支指向当前节点

git checkout --detach

由于哈希值是一串很长很长的乱码，在实际操作中使用哈希值分离 HEAD 很麻烦，所以 Git 也提供了 HEAD 基于某一特殊位置（分支/HEAD）直接指向前一个或前 N 个节点的命令，也即相对引用，如下：

//HEAD 分离并指向前一个节点

git checkout  分支名/HEAD^

//HEAD 分离并指向前 N 个节点

git checkout  分支名～ N

将 HEAD 分离出来指向节点有什么用呢？举个例子：如果开发过程发现之前的提交有问题，此时可以将 HEAD 指向对应的节点，修改完毕后再提交，此时你肯定不希望再生成一个新的节点，而你只需在提交时加上--amend 即可，具体命令如下：

git commit --amend

回退：

回退场景在平时开发中还是比较常见的，比如你巴拉巴拉写了一大堆代码然后提交，后面发现写的有问题，于是你想将代码回到前一个提交，这种场景可以通过 reset 解决，具体命令如下：

//回退 N 个提交

git reset HEAD~N

reset 和相对引用很像，区别是 reset 会使分支和 HEAD 一并回退。

**远程相关**

当我们接触一个新项目时，第一件事情肯定是要把它的代码拿下来，在 Git 中可以通过 clone 从远程仓库复制一份代码到本地，具体命令如下：

git clone 仓库地址

前面的章节我也有提到过，clone 不仅仅是复制代码，它还会把远程仓库的引用（分支/HEAD）一并取下保存在本地，如图 7 所示：

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/yc1he6/1622380437328-fac36536-b8d8-4ee2-b0cc-9e52d5828b2f.webp)
图 7

其中 origin/master 和 origin/ft-1 为远程仓库的分支，而远程的这些引用状态是不会实时更新到本地的，比如远程仓库 origin/master 分支增加了一次提交，此时本地是感知不到的，所以本地的 origin/master 分支依旧指向 C4 节点。我们可以通过 fetch 命令来手动更新远程仓库状态。

小提示：并不是存在服务器上的才能称作是远程仓库，你也可以 clone 本地仓库作为远程，当然实际开发中我们不可能把本地仓库当作公有仓库，说这个只是单纯的帮助你更清晰的理解分布式。

fetch：

说的通俗一点，fetch 命令就是一次下载操作，它会将远程新增加的节点以及引用(分支/HEAD)的状态下载到本地，具体命令如下：

git fetch  远程仓库地址/分支名

pull：

pull 命令可以从远程仓库的某个引用拉取代码，具体命令如下：

git pull  远程分支名

其实 pull 的本质就是 fetch+merge，首先更新远程仓库所有状态到本地，随后再进行合并。合并完成后本地分支会指向最新节点。

另外 pull 命令也可以通过 rebase 进行合并，具体命令如下：

git pull --rebase  远程分支名

push：

push 命令可以将本地提交推送至远程，具体命令如下：

git push  远程分支名

如果直接 push 可能会失败，因为可能存在冲突，所以在 push 之前往往会先 pull 一下，如果存在冲突本地解决。push 成功后本地的远程分支引用会更新，与本地分支指向同一节点。

综上所述

- 不管是 HEAD 还是分支，它们都只是引用而已，引用+节点是 Git 构成分布式的关键
- merge 相比于 rebase 有更明确的时间历史，而 rebase 会使提交更加线性应当优先使用
- 通过移动 HEAD 可以查看每个提交对应的代码
- clone 或 fetch 都会将远程仓库的所有提交、引用保存在本地一份
- pull 的本质其实就是 fetch+merge，也可以加入--rebase 通过 rebase 方式合并
