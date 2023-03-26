---
title: ECMAScript 包管理器
weight: 3
---

# 概述

> 参考：
>
> - [官方文档，cli-配置 npm-文件夹](https://docs.npmjs.com/cli/v8/configuring-npm/folders)

ECMAScript 的模块与包相关概念与工具是相当混乱的，最早期是由 Node.js 安装时自带的 NPM 工具集进行管理，并且 NPM 工具集中的工具除了模块与包的管理，还可以提供运行时环境等功能。

在后期，出现了 yarn、pnpm 等新型的管理工具，可以通过 Node.js 自带的 `corepack enable` 命令启用这些新的包管理工具。

ECMAScript 的包管理器安装各种包、依赖时，早期都是分开的，可以安装在项目当前目录，或某一个统一目录。但是如果项目想要正常启动，一般都只能安装在项目的当前目录（历史原因已不可考，反正很恶心）。

后来出现的 pnpm 工具，可以让我们将各种不同的项目的依赖放在同一个路径下进行统一管理。

## 模块与包的存储路径

npm 工具会通过 $PREFIX 与 node_modules 组合来决定其所管理的各种依赖包应该保存在什么位置。

其他工具也基本都符合这两点最基本的定义。

### PREFIX 配置

npm 有一个自带的配置 PREFIX，PREFIX 用来定位目录前缀，以决定将文件放在文件系统的何处。可以通过 `npm config get prefix` 命令查看 PREFIX 的值。

**PREFIX 通常默认为 Node.js 的安装路径**

- Linux 中，我个人通常装在 `/usr/local/nodejs/` 目录下。
- Windows 由于某些原因，使用 msi 安装包安装的 Node.js 会将该 PREFIX 改为 `%APPDATA%/npm/`，而不是安装目录。
  - 可以从 nodejs 安装路径下的 node_modules/npm/npmrc 文件中看到有这么一条配置：`prefix=${APPDATA}\npm`
  - 但是我们可以使用 zip 包，手动安装 Node.js，详情见：[ECMAScript 环境安装与使用](/docs/2.编程/高级编程语言/ECMAScript/ECMAScript%20环境安装与使用/ECMAScript%20环境安装与使用.md#Windows)

### node_modules 目录

当我们使用包管理命令安装各种第三方库(依赖包)及其衍生物通常会保存在名为 `node_modules/` 目录下，通常会有两个地方有 node_modules 目录：

- **Locally(本地)** # 这是默认的行为，安装的东西放在当前目录的 `./node_modules/` 目录中
    - 当我们想要在代码中使用 require() 或 import 导入模块时，通常安装在本地
- **Globally(全局)** # 使用 `-g` 选项，将安装的东西放在 `${PREFIX}/lib/node_modules/` 目录中
    - 若安装的东西中具有可以在 CLI 执行的工具，则同时会在 `${PREFIX}/bin/` 目录下生成指向原始文件的软链接，`${PREFIX}/bin/` 目录通常都会加入到 `${PATH}` 变量中。
    - 当我们安装的包可以在命令行执行时，通常安装在全局

> 注意：Windows 的全局 node_modules 目录与 Linux 不太一样，全局路径是 ${PREFIX}/node_modeuls。也就是说生成的链接文件就在 ${PREFIX} 下。

# NPM

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
  - **${PREFIX}/etc/npmrc** # 全局配置文件，可以通过 `--globalconfig` 命令行选项或 `${NPM_CONFIG_GLOBALCONFIG}` 环境变量改变其值
  - **~/.npmrc** # 用户配置文件，可以通过 `--userconfig` 命令行选项或 `${NPM_CONFIG_USERCONFIG}` 环境变量改变其值
  - **/PATH/TO/MY/PROJECT/.npmrc** # 每个项目自己的配置

**${PREFIX}/bin/*** # npm 安装的各种依赖包中若包含命令行工具，则会在此目录创建软链接。该目录通常都会加入到 `${PATH}` 变量中。

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

- **项目中的 node_models/ 应该使用与 store-dir 目录在同一个分区中**

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

## 更新 pnpm

```
corepack prepare pnpm@7.14.1 --activate
```

## pnpm 关联文件与配置

**/PATH/TO/.pnpm-store/** # 存放各项目依赖的目录。默认为根目录下的 .pnpm-store/ 目录。可以通过 `pnpm config -g set store-dir <DIR>` 修改

# npm 与 pnpm 语法

> 参考：
>
> - [官方文档，cli](https://docs.npmjs.com/cli)

通常，适用于 npm 的选项，也适用于 pnpm

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

# Yarn

> 参考：
>
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

# 实用工具
>
> 参考：
>
> - https://www.ruanyifeng.com/blog/2019/02/npx.html

npx 是 NPM 中自带的工具

通过 `npx serve` 命令(与 `npm exec serve` 命令类似)可以启动一个 HTTP 服务，以访问当前目录下的所有静态资源文件。便于本地开发调试。
