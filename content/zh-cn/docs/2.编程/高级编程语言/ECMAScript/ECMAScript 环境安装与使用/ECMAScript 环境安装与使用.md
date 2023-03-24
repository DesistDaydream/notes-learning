---
title: ECMAScript 环境安装与使用
weight: 1
---

# 概述

> 参考：

有两种运行时环境可以运行 ECMAScript 代码(Javascript、Typescript)

- **Browser**# 浏览器就是 ECMAScript 的运行时环境。
- **Node.js** # 在服务器或 PC 上安装 Node.js 环境以运行 ECMAScript 代码
- **Deno** # [GtiHub 项目，denoland/deno](https://github.com/denoland/deno)。据说要替代 Node.js，很尴尬的是。。。早期 18 年的 issue 中被各种国人无意义灌水
- **Bun** # [GtiHub 项目，oven-sh/bun](https://github.com/oven-sh/bun)

但是这两者可以提供的 API 能力各不相同

- Browser 可以提供的 dockument、window 和其他关于 DOM 或其他 Web 平台 API 的对象。
- Node.js 则可以提供很多操作系统相关的 API，比如文件的读写、进程管理、网络通信等等。

Node.js 与 Browser 都是基于 Chrome V8 引擎的 ECMAScript 运行时环境

你也许会想，直接在我的硬盘上创建好 HTML 和 JavaScript 文件，然后用浏览器打开，不就可以看到效果了吗？

这种方式运行部分 JavaScript 代码没有问题，但由于浏览器的安全限制，以 file:// 开头的地址无法执行如联网等 JavaScript 代码，最终，你还是需要架设一个 Web 服务器，然后以 http:// 开头的地址来正常执行所有 JavaScript 代码。

所以，想要运行 JavaScript 编写的静态文件，通常都需要一个可以提供 HTTP 服务的程序，以便响应客户端的请求返回这些静态文件。通常在生产环境中，我们都会将静态资源文件放到 Nginx 的目录中，由 Nginx 为客户端提供 HTTP 服务。

而本地测试时，我们通过 npm 与 yarn 工具即可安装一个简易的 HTTP 服务，这个 HTTP 服务都是由 JS 代码写的，必须要保证本地有 Node.js 环境，即可启动一个 HTTP 服务

```bash
# Npm 安装 http-server
~]# npm install http-server

# Yarn 安装 http-server
~]# yarn add http-server
```

http-server 来源于 [GitHub 项目，http-party/http-server](https://github.com/http-party/http-server)

通过 `npm exec serve` 命令将会安装另一个名为 serve 的库以运行 HTTP 服务，默认在 3000 端口上启动 HTTP 服务，通过浏览器访问 localhost:3000 即可获取到所有自己编写的静态文件，便于让开发人员在本地调试。

## (重点)运行具有 Module(模块) 功能的静态资源

从《ECMAScript 模块》章节可以看到，当我们使用 `import name from './one.js'` 导入的模块是一个文件时，那么这个文件将会被响应给浏览器，如果使用 `import name from 'one'` 导入的模块是一组文件时，此时浏览器无法直接识别，将会产生报错。因为浏览器想要执行这一组文件需要发起很多次的请求将所有文件都加载到本地，这其中的路由路径将是不可控的。

所以，此时我们则需要想办法将这一组文件变为一个文件响应给浏览器以便加载代码。这个转换的操作，我们可以使用打包工具(i.e.Webpack、Vite 等等) 将源代码**打包编译**成新的静态文件即可。

后面的逻辑，与基本运行 ECMAScript 代码的行为就是一样的了。

## 使用 Vite 运行 ECMAScript 代码

npm、yarn 的打包后运行代码的速度非常缓慢，才是推荐使用 Vite 工具启动 HTTP 服务并运行 JS/TS 代码，详见：[《Vite》](/docs/IT学习笔记/2.编程/高级编程语言/ECMAScript/ECMAScript%20 工具/Vite.md 工具/Vite.md) 章节

## 运行 TypeSript

Node.js 和 浏览器都无法直接运行 TypeScript 代码，这是因为 TS 的代码需要先转换为 JS 代码才可以运行。此时就需要一种工具，先转换再运行，或者直接转换运行一体。

- **tsc** # 将 TS 代码转换为 JS 代码。`npm install -g typescript`
- **ts-node** # 可以直接转换并运行 TS 代码，`npm install -g ts-node` 安装即可

# Node.js

> 参考：
>
> - [org 官网](https://nodejs.org/en/)
> - [dev 官网](https://nodejs.dev/)
> - [dev 官网中文翻译](http://nodejs.cn/)
> - [Wiki,Node.js](https://en.wikipedia.org/wiki/Node.js)

Node.js 是基于 Chrome V8 引擎的 ECMAScript 运行时环境，由 RyanDahl 于 2009 年 5 月 27 日发布。转年(i.e.2010 年 1 月)，为 Node.js 环境引入了一个名为 npm 的包管理器。包管理器使程序员更容易发布和共享 Node.js 的源代码，旨在简化包的安装、更新、卸载。

Browser 和 Node.js 都是 ECMAScript 的运行时环境，但是这两者可以提供的 API 能力各不相同

- Browser 可以提供的 dockument、window 和其他关于 DOM 或其他 Web 平台 API 的对象。
- Node.js 则可以提供很多操作系统相关的 API，比如文件的读写、进程管理、网络通信等等。

通过 Node.js，可以让我们使用一种语言编写前端与后端。我们甚至可以通过 npm 与 yarn 安装第三方库后，使用 Node.js 在本地监听端口并响应给客户端静态资源文件。

## 安装 Node.js

### Linux

从[官网](https://nodejs.org/zh-cn/download/)下载 Linux 版的 `.tar.xg` 包，并解压

```bash
export NodejsVersion="v18.15.0"
wget https://nodejs.org/dist/${NodejsVersion}/node-${NodejsVersion}-linux-x64.tar.xz
tar -xvf node-${NodejsVersion}-linux-x64.tar.xz -C /usr/local/

mv /usr/local/node-${NodejsVersion}-linux-x64 /usr/local/nodejs
```

配置环境变量

```bash
sudo tee /etc/profile.d/nodejs.sh > /dev/null <<-"EOF"
export PATH=$PATH:/usr/local/nodejs/bin
EOF
```

### Windows

警告！！！由于 msi 安装包会修改 %PREFIX% 为 `%APPDATA%\npm` ，并将该目录到 $PATH。我个人推荐下载 zip，并自己解压到想要的位置后，手动配置环境变量。

```powershell
$NodejsVersion = "18.14.1"
$NodejsUrl = "https://nodejs.org/dist/v$NodejsVersion/node-v$NodejsVersion-win-x64.zip"
$TempZipFile = "D:\tmp\nodejs.zip"
$ExtractPath = "D:\Tools"

# Download the zip file to a temporary location
Invoke-WebRequest -Uri $NodejsUrl -OutFile $TempZipFile

# Extract the contents of the zip file to the installation directory and rename the top-level directory to "nodejs"
Expand-Archive -Path $TempZipFile -DestinationPath $ExtractPath
Rename-Item -Path "$ExtractPath\node-v$NodejsVersion-win-x64" -NewName "nodejs"
```

将 nodejs/ 目录添加到用户的 PATH 环境变量中

```powershell
$path = [Environment]::GetEnvironmentVariable("Path", "User")
$newPath = "D:\Tools\nodejs"
[Environment]::SetEnvironmentVariable("Path", "$path;$newPath", "User")
```

### 目录结构

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

## NVM

> 参考：
>
> - [GitHub 项目，nvm-sh/nvm](https://github.com/nvm-sh/nvm)

**Node Version Manager(Node.js 版本管理器，简称 NVM)**

## Node.js 关联文件与配置

详见：[npm 关键文件与配置](/docs/2.编程/高级编程语言/ECMAScript/ECMAScript%20环境安装与使用/ECMAScript%20模块与包的使用与管理.md#npm%20关联文件与配置)

# 初始化项目

npm 等包管理工具下载完框架后，会自动生成项目目录，这些目录是已经初始化完成的项目，包含了很多必要的文件，比如 package.json 等。

随着学习深入，会逐步总结一个前端项目如果不使用框架从头构建的话会需要什么。

- **public/** 目录存放公共资源
- **src/** # 所有代码都在 src/ 目录下
- **.eslintrc.cjs** # ESLint 程序配置
- **.prettierc.json** # Prettier 插件的配置
- **env.d.ts** #
- **index.html** # 程序入口
- [**package.json**](/docs/IT学习笔记/2.编程/高级编程语言/ECMAScript/ECMAScript%20 环境安装与使用/package.json.md 环境安装与使用/package.json.md) # 包管理器配置文件，比如 npm、pnpm 等
- **vite.config.ts** # Vite 程序给项目打包时使用的配置
- **tsconfig.json** #
- **tsconfig.config.json** # [TSConfig](/docs/IT学习笔记/2.编程/高级编程语言/ECMAScript/ECMAScript%20 环境安装与使用/TSConfig.md 环境安装与使用/TSConfig.md) 文件

## JavaScript 项目初始化

无

## TypeScript 项目初始化

使用 `npm install -g typescript` 安装 tsc 命令。

使用 `tsc init` 命令将会生成 [TSConfig](/docs/IT学习笔记/2.编程/高级编程语言/ECMAScript/ECMAScript%20 环境安装与使用/TSConfig.md 环境安装与使用/TSConfig.md) 文件。

# 编译与打包

> 参考：
>
> - [GitHub 项目，webpack/webpack](https://github.com/webpack/webpack)
> - [GitHub 项目，rollup/rollup](https://github.com/rollup/rollup)

大型项目通常都要打包，打包工具有很多：

- Webpack
- Rollup
- Vite
- ......等等

# 常见问题

## Node.js 无法使用 ES6 语法问题

无法使用 import 关键字导入模块。参考 <https://nodejs.org/docs/latest-v16.x/api/esm.html#enabling>，在 package.json 文件中设置 `"type": "module"` 或者使用 `--input-type=module` 命令行参数以告诉 Node.js 使用 ECMAScript 模块加载器。默认情况下，Node.js 使用 CommonJS 模块加载器。
