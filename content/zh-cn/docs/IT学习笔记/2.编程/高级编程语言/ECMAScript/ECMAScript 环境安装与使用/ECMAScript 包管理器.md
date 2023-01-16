---
title: ECMAScript 包管理器
---

# 概述

> 参考：
> -

当我们使用包管理命令安装各种第三方库(依赖包)及其衍生物通常会保存在两类地方
- **Locally(本地)** # 这是默认的行为，安装的东西放在当前目录的 `./node_modules/` 目录中
  - 当我们想要在代码中使用 require() 或 import 导入模块时，通常安装在本地
- **Globally(全局)** # 使用 `-g` 选项，将安装的东西放在 `${PREFIX}/lib/node_modules/` 目录中；若安装的东西中具有可以在 CLI 执行的工具，则同时会在 `${PREFIX}/bin/` 目录下生成指向原始文件的软链接，`${PREFIX}/bin/` 目录通常都会加入到 `${PATH}` 变量中。
  - 当我们想要在命令行上运行安装的命令行工具，通常安装在全局

随着时代的发展，出现了 pnpm、(期待有更好的)等 工具，可以让我们将各种不同的项目的依赖放在同一个路径下进行统一管理。

# NPM

> 参考：
> - [官网](https://www.npmjs.com/)
> - [官方文档，cli](https://docs.npmjs.com/cli)

**Node.js Package Manager(简称 NPM)** 是 Node.js 自带的包管理工具，通常与 Node.js 一同安装，最初版本于 2010 年 1 月发行。NPM 本质是一个第三方模块，可以在 **${PREFIX}/lib/node_modules/npm/\*** 目录下找到 npm 的所有文件。

> 注意：
> - `${PREFIX}` 指 **Node.js 的安装路径**，Linux 中通常装在 /usr/local/nodejs 目录下，Windows 则随意了~~

## npx

> 参考：
> - <https://www.ruanyifeng.com/blog/2019/02/npx.html>

npx 是 NPM 中自带的工具

通过 `npx serve` 命令(与 `npm exec serve` 命令类似)可以启动一个 HTTP 服务，以访问当前目录下的所有静态资源文件。便于本地开发调试。

## npm 关联文件与配置

> 参考：
> - [官方文档，cli-配置 npm-文件夹](https://docs.npmjs.com/cli/v8/configuring-npm/folders)
> - [官方文档，cli-使用 npm-配置](https://docs.npmjs.com/cli/v8/using-npm/config)

当我们使用 `npm install` 命令安装各种第三方库(依赖包)及其衍生物通常会保存在两类地方

- **Locally(本地)** # 这是默认的行为，安装的东西放在当前目录的 `./node_modules/` 目录中
  - 当我们想要在代码中使用 require() 或 import 导入模块时，通常安装在本地
- **Globally(全局)** # 使用 `-g` 选项，将安装的东西放在 `${PREFIX}/lib/node_modules/` 目录中；若安装的东西中具有可以在 CLI 执行的工具，则同时会在 `${PREFIX}/bin/` 目录下生成指向原始文件的软链接，`${PREFIX}/bin/` 目录通常都会加入到 `${PATH}` 变量中。
  - 当我们想要在命令行上运行安装的命令行工具，通常安装在全局

npm 从 命令行、环境变量、npmrc 文件、某些情况下从 package.json 文件 这些地方获取其配置信息

npm 从以下地方获取其运行时配置

- **命令行标志**
- **环境变量**
- **npmrc 文件** # npm 从以下几个地方依次读取 nmrc 文件
  - **/PATH/TO/NPM/npmrc** # npm 内置的配置文件
  - **${PREFIX}/etc/npmrc** # 全局配置文件，可以通过 `--globalconfig` 命令行选项或 `${NPM_CONFIG_GLOBALCONFIG}` 环境变量改变其值
  - **~/.npmrc** # 用户
  - 配置文件，可以通过 `--userconfig` 命令行选项或 `${NPM_CONFIG_USERCONFIG}` 环境变量改变其值
  - **/PATH/TO/MY/PROJECT/.npmrc** # 每个项目自己的配置

**${PREFIX}/lib/node_modules/npm/\* **# npm 作为一个第三方模块，跟随 Node.js 一起安装，被放在该目录下。
**${PREFIX}/bin\* **# npm 安装的各种依赖包中若包含命令行工具，则会在此目录创建软链接。该目录通常都会加入到 `${PATH}` 变量中。

### 配置文件详解

所有可供配置的信息可从 <https://docs.npmjs.com/cli/v8/using-npm/config#config-settings> 查看

# PNPM

> 参考：
> - [GitHub 项目，pnpm/pnpm](https://github.com/pnpm/pnpm)
> - [官网](https://pnpm.io/)
> - [稀土掘金，pnpm 对比 npm/yarn 好在哪里](https://juejin.cn/post/7047556067877716004)

当使用 npm 或 Yarn 时，如果你有 100 个项目使用了某个依赖（dependency），就会有 100 份该依赖的副本保存在硬盘上。 而在使用 pnpm 时，依赖会被存储在内容可寻址的存储中，所以：

1. 如果你用到了某依赖项的不同版本，只会将不同版本间有差异的文件添加到仓库。 例如，如果某个包有 100 个文件，而它的新版本只改变了其中 1 个文件。那么 pnpm update 时只会向存储中心额外添加 1 个新文件，而不会因为仅仅一个文件的改变复制整新版本包的内容。
2. 所有文件都会存储在硬盘上的某一位置。 当软件包被被安装时，包里的文件会硬链接到这一位置，而不会占用额外的磁盘空间。 这允许你跨项目地共享同一版本的依赖。

因此，您在磁盘上节省了大量空间，这与项目和依赖项的数量成正比，并且安装速度要快得多！

## 安装 pnpm

使用 `corepack enable` 指令启用 pnpm。
设置包的存储路径：

- Windows：`pnpm config set store-dir D:\Projects\.pnpm-store`
- Linux：`pnpm config set store-dir /mnt/d/Projects/.pnpm-store`

配置镜像源 `pnpm config set registry="https://registry.npmmirror.com"`

若 Windows 无法执行 pnpm，报错：`pnpm : 无法加载文件 D:\Tools\nodejs\pnpm.ps1，因为在此系统上禁止运行脚本。有关详细信息，请参阅 https:/go.microsoft.com/fwlink/?LinkID=135170 中的 about_Execution_Policies。`

- 此时需要在 PowerShell 中执行 `Set-ExecutionPolicy -Scope CurrentUser RemoteSigned` 指令。详见[微软官网解释](https://learn.microsoft.com/zh-cn/powershell/module/microsoft.powershell.core/about/about_execution_policies?view=powershell-7.2)。
  - 其中 `-Scope CurrentUser` 是指针对当前用户的操作，若使用管理员运行 VSCode 或 PowerShell，则不用加这个选项。

### 更新

corepack prepare pnpm@7.14.1 --activate

## pnpm 关键文件与配置

**/PATH/TO/.pnpm-store** # 存放各项目依赖的目录

# npm 与 pnpm Syntax(语法)

> 参考：
> -

通常，适用于 npm 的选项，也适用于 pnpm

**npm [OPTIONS] COMMAND [OPTIONS]**

**OPTIONS**

- **-g, --global** # 指定命令作用范围为全局。默认情况下 npm 的所有子命令作用范围是当前目录

## npm config

npm config 用来管理 npm 的配置文件，i.e.npmrc 文件。

### Syntax(语法)
**npm config set <KEY>=<VALUE> [<KEY>=<VALUE> ...]**
**npm config get [<KEY> [<KEY> ...]]**
**npm config delete <KEY> [<KEY> ...]**
**npm config list [--json]**
**npm config edit**

OPTIONS
- **-g, --global** # 对全局配置文件(${PREFIX}/etc/npmrc) 执行操作

### EXAMPLE
- 配置镜像源为淘宝的
  - `npm config set registry="https://registry.npmmirror.com"`

## npm exec
从本地或远程 npm 包运行命令

### Syntax(语法)
**npm exec -- <pkg>[@<version>] [args...]**
**npm exec --package=<pkg>[@<version>] -- <cmd> [args...]**
**npm exec -c '<cmd> [args...]'**
**npm exec --package=foo -c '<cmd> [args...]'**

OPTIONS

-

## npm install
### Syntax(语法)

OPTIONS
- **-D, --save-dev** # 安装的包将会出现在 `devDependencies` 中

## npm init

创建一个 package.json 文件。用来初始化一个项目

### Syntax(语法)
**npm init [--force|-f|--yes|-y|--scope]**
**npm init <@scope> (same as `npx <@scope>/create`)**
**npm init [<@scope>]<name> (same as `npx [<@scope>/]create-<name>`)**

OPTIONS

-

## npm list

列出所有已经安装的包。默认列出当前项目中已安装的包，通常检查如下目录：`node_modules/`

## npm update

更新已安装的包

# Yarn

> 参考：
> - [官网](https://yarnpkg.com/)

管理 Yarn 的首选方式是通过 [Corepack](https://nodejs.org/dist/latest/docs/api/corepack.html)，这是从 16.10 开始随所有 Node.js 版本一起提供的新二进制文件。它充当我们和 Yarn 之间的中介，让我们在多个项目中使用不同的包管理器版本，而无需再签入 Yarn 二进制文件。

Node.js >=16.10
Corepack 默认包含在所有 Node.js 安装中，但目前是可选的。要启用它，请运行以下命令：
`corepack enable`

Node.js <16.10
在 16.10 之前的版本中，Node.js 不包含 Corepack；为了解决这个问题，运行：
`npm i -g corepack`

配置镜像源以加速下载各种依赖包

```bash
yarn config set registry https://registry.npmmirror.com -g
```

配置 $PATH 以便可以直接执行通过 yarn 安装的各种工具

```bash
export PATH=$PATH:~/.config/yarn/global/node_modules/.bin
```

## Yarn 关联文件与配置

**~/.yarnrc** # 配置文件
**~/.config/yarn/*** #

# yarn Syntax(语法)
