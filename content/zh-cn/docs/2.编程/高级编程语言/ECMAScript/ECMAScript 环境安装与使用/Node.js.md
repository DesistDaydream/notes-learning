---
title: "Node.js"
created: "2026-09-26T11:49"
weight: 100
---

# 概述

> 参考：
>
> - [org 官网](https://nodejs.org/en/)
> - [dev 官网](https://nodejs.dev/)
> - [dev 官网中文翻译](http://nodejs.cn/)
> - [Wiki, Node.js](https://en.wikipedia.org/wiki/Node.js)

Node.js 是基于 Chrome V8 引擎的 ECMAScript 运行时环境，由 RyanDahl 于 2009 年 5 月 27 日发布。转年(i.e.2010 年 1 月)，为 Node.js 环境引入了一个名为 npm 的包管理器。包管理器使程序员更容易发布和共享 Node.js 的源代码，旨在简化包的安装、更新、卸载。

Browser 和 Node.js 都是 ECMAScript 的运行时环境，但是这两者可以提供的 API 能力各不相同

- Browser 可以提供的 dockument、window 和其他关于 DOM 或其他 Web 平台 API 的对象。
- Node.js 则可以提供很多操作系统相关的 API，比如文件的读写、进程管理、网络通信等等。

通过 Node.js，可以让我们使用一种语言编写前端与后端。我们甚至可以通过 npm 与 yarn 安装第三方库后，使用 Node.js 在本地监听端口并响应给客户端静态资源文件。

# 安装 Node.js

### Linux

从[官网](https://nodejs.org/zh-cn/download/)下载 Linux 版的 `.tar.xg` 包，并解压

```bash
export NodejsVersion="v24.21.0"
wget https://nodejs.org/dist/${NodejsVersion}/node-${NodejsVersion}-linux-x64.tar.xz
sudo tar -xvf node-${NodejsVersion}-linux-x64.tar.xz -C /usr/local/

sudo mv /usr/local/node-${NodejsVersion}-linux-x64 /usr/local/nodejs
```

配置环境变量

```bash
sudo tee /etc/profile.d/nodejs.sh > /dev/null <<-EOF
export PATH=${HOME}/.local/share/pnpm/bin:/usr/local/nodejs/bin:\${PATH}
export COREPACK_NPM_REGISTRY="https://registry.npmmirror.com"
EOF
source /etc/profile.d/nodejs.sh
```

### Windows

> [!Warning] 由于 msi 安装包会修改 `$PREFIX` 为 `%APPDATA%\npm` ，并将该目录加入到 $PATH。个人推荐下载 zip，并自己解压到想要的位置后，手动配置环境变量。

```powershell
$NodejsVersion = "24.21.0"
$NodejsUrl = "https://nodejs.org/dist/v$NodejsVersion/node-v$NodejsVersion-win-x64.zip"
$TempZipFile = "D:\tmp\nodejs.zip"
$ExtractPath = "D:\Tools"

# 将 zip 文件下载到临时位置
Invoke-WebRequest -Uri $NodejsUrl -OutFile $TempZipFile

# 将 zip 文件内容解压到安装目录并将顶级目录重命名为 "nodejs"
Expand-Archive -Path $TempZipFile -DestinationPath $ExtractPath
Rename-Item -Path "$ExtractPath\node-v$NodejsVersion-win-x64" -NewName "nodejs"
```

将 nodejs/ 目录添加到用户的 PATH 环境变量中

```powershell
$path = [Environment]::GetEnvironmentVariable("Path", "User")
$newPath = "D:\Tools\nodejs"
[Environment]::SetEnvironmentVariable("Path", "$path;$newPath", "User")
```

### 初次安装后的目录结构

Linux 目录结构，node_modules/ 目录在 lib/ 目录下，这点与 Windows 不同。

```bash
]# tree -L 2 -F
.
├── bin/
│   ├── corepack -> ../lib/node_modules/corepack/dist/corepack.js*
│   ├── node*
│   ├── npm -> ../lib/node_modules/npm/bin/npm-cli.js*
│   ├── npx -> ../lib/node_modules/npm/bin/npx-cli.js*
│   ├── pnpm -> ../lib/node_modules/corepack/dist/pnpm.js*
│   ├── pnpx -> ../lib/node_modules/corepack/dist/pnpx.js*
│   ├── yarn -> ../lib/node_modules/corepack/dist/yarn.js*
│   └── yarnpkg -> ../lib/node_modules/corepack/dist/yarnpkg.js*
├── CHANGELOG.md
├── include/
│   └── node/
├── lib/
│   └── node_modules/
├── LICENSE
├── README.md
└── share/
    ├── doc/
    ├── man/
    └── systemtap/
```

Windows 目录结构

```bash
$ tree -L 2 -F
./
├── CHANGELOG.md*
├── LICENSE*
├── README.md*
├── corepack*
├── corepack.cmd*
├── install_tools.bat*
├── node.exe*
├── node_etw_provider.man*
├── node_modules/
│   ├── corepack/
│   └── npm/
├── nodevars.bat*
├── npm*
├── npm.cmd*
├── npx*
└── npx.cmd*
```

### 初次安装后的配置

配置 [NPM](/docs/2.编程/高级编程语言/ECMAScript/ECMAScript%20工具/NPM.md)。通常需要 配置 npm 镜像源、添加 pnpm、修改 pnpm 镜像源、修改 pnpm 存储路径



### NVM

> 参考：
>
> - [GitHub 项目，nvm-sh/nvm](https://github.com/nvm-sh/nvm)

**Node Version Manager(Node.js 版本管理器，简称 NVM)**

# 关联文件与配置

[NPM](/docs/2.编程/高级编程语言/ECMAScript/ECMAScript%20工具/NPM.md) 相关的关联文件与配置

# Corepack

> 参考：
>
> - [GitHub 项目，nodejs/corepack](https://github.com/nodejs/corepack)

Corepack 是一个零运行时依赖的 Node.js 脚本，充当 Node.js 项目与包管理器（e.g. [NPM](docs/2.编程/高级编程语言/ECMAScript/ECMAScript%20工具/NPM.md)、Yarn、etc.）之间的桥梁。

## Corepack 关联文件与配置

**${COREPACK_HOME}** # Corepack 的工作目录

> [!Note] 不同系统的默认路径
>
> Windows: `%LOCALAPPDATA%/node/corepack/`
>
> Unix-like: `${HOME}/.cache/node/corepack/`

- **./lastKnownGood.json** # Known Good Release(已知的良好版本，简称 KGR) 的清单。通过 Corepack 使用的 npm、pnpm、yarn 默认版本记录在该文件中。KGR 中包管理器的版本会随着 Corepack 的更新而更新。
- **./v1/** # 目录组织结构的版本
    - **./corepack-${数字}/** # 暂存目录，通常都是空的。不知道为什么不放在 /tmp/ 下。
    - **./pnpm/${PNPM_VERSION}/** # Corepack 管理的 pnpm 工具保存位置
    - **./yarn/${YARN_VERSION}/** # Corepack 管理的 yarn 工具保存位置

# 常见问题

记录于 2024.1.15: 下面这些问题，最好都不要去解决了，使用 Deno 去吧！！！

## Node.js 无法使用 ES6 语法问题

在使用 import 语法导入包的代码中，使用 node 命令运行，报错: `SyntaxError: Cannot use import statement outside a module`。

本质上上 Node.js 默认无法使用 import 关键字导入模块。参考 <https://nodejs.org/docs/latest-v16.x/api/esm.html#enabling>，在 package.json 文件中设置 `"type": "module"` 或者使用 `--input-type=module` 命令行参数以告诉 Node.js 使用 ECMAScript 模块加载器。默认情况下，Node.js 使用 CommonJS 模块加载器。

## ts-node 无法执行 .ts 脚本

报错: `TypeError [ERR_UNKNOWN_FILE_EXTENSION]: Unknown file extension ".ts"`

https://stackoverflow.com/questions/62096269/cant-run-my-node-js-typescript-project-typeerror-err-unknown-file-extension

在这里有讨论

- https://github.com/TypeStrong/ts-node/issues/935
- https://github.com/TypeStrong/ts-node/issues/1007#issuecomment-1163471306

这是因为我们在 package.json 中使用了 `"type": "module"` 配置，所以需要删除该配置。若是不想删除该配置，则可以在 `tsconfig.json` 文件中添加如下配置

```
  "compilerOptions": {
    "esModuleInterop": true,
  }
```

然后使用 `ts-node-esm` 命令而不是 ts-node 命令执行 .ts 脚本。

## Node.js 运行 ES6 语法的 TS 代码

综合上面两个问题，保证 package.json 和 tsconfig.json 的最低配置。同时使用 `ts-node-esm` 命令运行 .ts 文件。

package.json

```json
{
 "type": "module",
}
```

tsconfig.json

```json
{
  "compilerOptions": {
    // "target": "es2016",
    "module": "ES6",
    "esModuleInterop": true,
  }
}
```

Notes: 有的时候 TS 依赖库还依赖原始的 JS 库，所以也要同时安装 JS 库。crypto-js 就是这个情况，要想使用 `ts-node-esm` 正常运行代码， package.json 至少需要如下内容：

```json
{
 "type": "module",
 "dependencies": {
  "@types/crypto-js": "^4.2.1",
  "crypto-js": "^4.2.0"
 }
}
```
