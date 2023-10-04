---
title: NPM
linkTitle: NPM
date: 2024-01-15T20:51
weight: 20
---

# 概述

> 参考：
>
> - [官网](https://www.npmjs.com/)
> - [官方文档，cli-使用 npm-配置](https://docs.npmjs.com/cli/v8/using-npm/config)

**Node.js Package Manager(简称 NPM)** 是 Node.js 自带的包管理工具，通常与 Node.js 一同安装，最初版本于 2010 年 1 月发行。NPM 本质是一个第三方模块，可以在 NodeJS 安装目录下的 **lib/node_modules/npm/*** 目录下找到 npm 的所有文件。

配置镜像源为淘宝的： `npm config -g set registry="https://registry.npmmirror.com"`

## npm 关联文件与配置

npm 从 命令行、环境变量、npmrc 文件 这些地方获取其配置信息：

- **命令行标志**
- **环境变量**
- **npmrc 文件** # npm 从以下几个地方依次读取 npmrc 文件
  - **/PATH/TO/NPM/npmrc** # npm 内置的配置文件。这内置的文件是不是不可见的？o(╯□╰)o
  - **`${PREFIX}/etc/npmrc`** # 全局配置文件，可以通过 `--globalconfig` 命令行选项或 `${NPM_CONFIG_GLOBALCONFIG}` 环境变量改变其值
  - **~/.npmrc** # 用户配置文件，可以通过 `--userconfig` 命令行选项或 `${NPM_CONFIG_USERCONFIG}` 环境变量改变其值
  - **/PATH/TO/MY/PROJECT/.npmrc** # 每个项目自己的配置

**`${PREFIX}/bin/`** # npm 安装的各种依赖包中若包含命令行工具，则会在此目录创建软链接。该目录通常都会加入到 `${PATH}` 变量中。

## 配置文件详解

所有可供配置的信息可从 https://docs.npmjs.com/cli/v8/using-npm/config#config-settings 查看

# PNPM

> 参考：
>
> - [GitHub 项目，pnpm/pnpm](https://github.com/pnpm/pnpm)
> - [官网](https://pnpm.io/)
> - [稀土掘金，pnpm 对比 npm/yarn 好在哪里](https://juejin.cn/post/7047556067877716004)

当使用 npm 或 Yarn 时，如果你有 100 个项目使用了某个依赖（dependency），就会有 100 份该依赖的副本保存在硬盘上。 而在使用 pnpm 时，依赖会被存储在内容可寻址的存储中，所以：

- 如果你用到了某依赖项的不同版本，只会将不同版本间有差异的文件添加到仓库。 例如，如果某个包有 100 个文件，而它的新版本只改变了其中 1 个文件。那么 pnpm update 时只会向存储中心额外添加 1 个新文件，而不会因为仅仅一个文件的改变复制整新版本包的内容。
- 所有文件都会存储在硬盘上的某一位置。 当软件包被被安装时，包里的文件会硬链接到这一位置，而不会占用额外的磁盘空间。 这允许你跨项目地共享同一版本的依赖。

因此，您在磁盘上节省了大量空间，这与项目和依赖项的数量成正比，并且安装速度要快得多！

store-dir 说明：

- **项目中的 `node_models/` 目录应该使用与 store-dir 目录在同一个分区中**

## 安装 pnpm

使用 `corepack enable` 指令启用 pnpm

设置包的存储路径：

- Windows：`pnpm config -g set store-dir D:\Projects\.pnpm-store`
- Linux：`pnpm config -g set store-dir /mnt/d/Projects/.pnpm-store`

配置镜像源：

- `pnpm config -g set registry="https://registry.npmmirror.com"`

若 Windows 无法执行 pnpm，报错：`pnpm : 无法加载文件 D:\Tools\nodejs\pnpm.ps1，因为在此系统上禁止运行脚本。有关详细信息，请参阅 https:/go.microsoft.com/fwlink/?LinkID=135170 中的 about_Execution_Policies。`

- 此时需要在 PowerShell 中执行 `Set-ExecutionPolicy -Scope CurrentUser RemoteSigned` 指令。详见[微软官网解释](https://learn.microsoft.com/zh-cn/powershell/module/microsoft.powershell.core/about/about_execution_policies?view=powershell-7.2)。
  - 其中 `-Scope CurrentUser` 是指针对当前用户的操作，若使用管理员运行 VSCode 或 PowerShell，则不用加这个选项。

**更新 pnpm**

```bash
export PNPM_VERSION="XXXX"
corepack prepare pnpm@${PNPM_VERSION} --activate
```

## pnpm 关联文件与配置

> Notes: pnpm 的官方文档中依然使用 .npmrc 这个单词作为自己得配置文件，但是实际上，配置文件的名称只是 rc。

pnpm 使用 [INI](/docs/2.编程/无法分类的语言/INI.md) 格式的配置文件。若想使用独立于本地项目的配置文件，可以建立一个名为 `.npmrc` 的文件在项目的根目录。

**/PATH/TO/rc** # pnpm 配置文件。`pnpm config` 命令控制的配置即是控制该文件。在[官方文档，CLI 命令-其他-pnpm config](https://pnpm.io/next/cli/config)中可以看到 pnpm 配置文件的保存路径：

- 如果设置了 `$XDG_CONFIG_HOME` 环境变量，则为 `$XDG_CONFIG_HOME/pnpm/rc`
- 在 Windows 上：**%LOCALAPPDATA%/pnpm/config/rc**
- 在 macOS 上：**~/Library/Preferences/pnpm/rc**
- 在 Linux 上：**~/.config/pnpm/rc**

**/PATH/TO/.pnpm-store/** # pnpm 管理的依赖保存路径。可以通过 `pnpm config -g set store-dir <DIR>` 修改。`pnpm store path` 命令会显示当前 .pnpm-store 的完整路径。在 [官方文档，配置-.npmrc](https://pnpm.io/next/npmrc#store-dir) 中可以看到 pnpm 存储目录的默认值位置：

- 如果存在 `$PNPM_HOME` 环境变量，则为 `$PNPM_HOME/pnpm/store`
- 如果设置了 `$XDG_DATA_HOME` 环境变量，则为 `$XDG_DATA_HOME/pnpm/store`
- 在 Windows 上: **~/AppData/Local/pnpm/store**
- 在 macOS 上: **~/Library/pnpm/store**
- 在 Linux 上: **~/.local/share/pnpm/store**
- 特殊默认值规则：pnpm 的管理的包的存储位置应始终与项目目录保持在同一磁盘上，因此每个磁盘将有一个存储。 如果在使用磁盘中具有主目录，存储目录就会创建在这里。 如果磁盘上没有主目录，那么将在文件系统的根目录中创建该存储。 例如，如果安装发生在挂载在 `/mnt` 的文件系统上，那么存储将在 `/mnt/.pnpm-store/` 目录中创建。 Windows 系统上也是如此。
  - 用人话说是：假如我的项目在 `D:/Projects/DesistDaydream/javascript-learning/`，那模块将会默认下载到 `D:/.pnpm-store/` 目录中。
- 注意：可以从不同的磁盘设置同一个存储，但在这种情况下，pnpm 将复制包而**不是硬链接**它们，因为**硬链接只能发生在同一文件系统上**。

# npm 与 pnpm 语法

> 参考：
>
> - [官方文档，cli](https://docs.npmjs.com/cli)

通常，适用于 npm 的选项，也适用于 pnpm。不过也有部分命令是独属于 npm 或 pnpm 的

**npm [OPTIONS] COMMAND [OPTIONS]**

**OPTIONS**

- **-g, --global** # 指定命令作用范围为全局。默认情况下 npm 的所有子命令作用范围是当前目录

## npm config

npm config 用来管理 npm 的配置文件，i.e.npmrc 文件。

### Syntax(语法)

**npm config COMMAND [KEY=VALUE]**

**COMMAND**

- set
- get
- delete
- list
- edit

**OPTIONS**

- **-g, --global** # 对全局配置文件(${PREFIX}/etc/npmrc) 执行操作

### EXAMPLE

列出所有已知配置

- npm config ls -l

获取 prefix 配置的值

- npm config get prefix

## npm exec

从本地或远程 npm 包运行命令

### Syntax(语法)

**npm exec -- \<pkg>[@\<version>] [args...]**
**npm exec --package=\<pkg>[@\<version>] -- \<cmd> [args...]**
**npm exec -c '\<cmd> [args...]'**
**npm exec --package=foo -c '\<cmd> [args...]'**

**OPTIONS**

## npm install

安装项目的所有依赖

install 可以简写为 i。

### Syntax(语法)

**OPTIONS**

- **-D, --save-dev** # 安装的包将会出现在 `devDependencies` 中

## npm init

> 别名：create、innit

创建一个 package.json 文件。用来初始化一个项目。

可以指定一个 **Initializer(初始化程序)** 并执行其 bin/ 下的 js/ts 文件以运行其他与初始化相关的操作。

### Syntax(语法)

npm init [INITIALIZER]

**OPTIONS**

## npm list

列出所有已经安装的包。默认列出当前项目中已安装的包，通常检查如下目录：`node_modules/`

## npm outdated

此命令将检查注册表，查看是否有任何（或指定的）已安装包目前已过时。

默认情况下，只显示根项目的直接依赖项和已配置的工作区的直接依赖项。使用--all查找所有过时的元依赖项。

在输出中:

- wanted # 是满足package.json中指定的semver范围的软件包的最大版本。如果没有可用的semver范围（即您在运行npm outdated --global或包未包含在package.json中），则wanted显示当前安装的版本。
- latest # 是注册表中标记为当前版本的软件包的版本。在没有特殊配置的情况下运行npm publish将使用最新的dist-tag发布该软件包。这可能是软件包的最大版本，也可能是软件包的最近发布版本，这取决于包开发人员如何管理最新npm help dist-tag。
- location 是软件包在物理树中的位置。
- depended by # 显示依赖于显示依赖项的软件包
- package type # （使用--long / -l时）告诉您此软件包是依赖项还是开发/同行/可选依赖项。始终标记未包含在package.json中的软件包的依赖项。
- homepage # （使用--long / -l时）是包的packument中包含的主页值
- 红色表示存在符合您的 semver 要求的更新版本，因此您应立即更新。
- 黄色表示有一个新的版本超出了您的 semver 要求（通常是新主要版本或新的0.x次要版本），因此请小心进行。

## npm update

更新已安装的包

## pnpm stroe

https://pnpm.io/cli/store

管理包的存储

## 最佳实践

一个全局模块的声明周期

```bash
# 安装
npm install -g vite
# 列出
npm list -g
# 更新
npm update -g vite
# 卸载
npm uninstall -g vite
```

# 实用工具

> 参考：
>
> - https://www.ruanyifeng.com/blog/2019/02/npx.html

npx 是 NPM 中自带的工具

通过 `npx serve` 命令(与 `npm exec serve` 命令类似)可以启动一个 HTTP 服务，以访问当前目录下的所有静态资源文件。便于本地开发调试。
