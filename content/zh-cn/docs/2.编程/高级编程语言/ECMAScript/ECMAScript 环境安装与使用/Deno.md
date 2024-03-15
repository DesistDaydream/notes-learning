---
title: Deno
linkTitle: Deno
date: 2024-03-15T13:06
weight: 20
---

# 概述

> 参考：
>
> - [GitHub 项目，denoland/deno](https://github.com/denoland/deno)
> - [GitHub 项目，denolib/awesome-deno](https://github.com/denolib/awesome-deno)
> - [阮一峰，Deno 运行时入门教程：Node.js 的替代品](https://www.ruanyifeng.com/blog/2020/01/deno-intro.html)
> - https://docs.deno.com/runtime/manual/references/vscode_deno/ VSCode 中使用 Deno 的方式

Deno 也是一个服务器运行时，但是支持多种语言，可以直接运行 JavaScript、TypeScript 和 WebAssembly 程序。

它内置了 V8 引擎，用来解释 JavaScript。同时，也内置了 tsc 引擎，解释 TypeScript。它使用 Rust 语言开发，由于 Rust 原生支持 WebAssembly，所以它也能直接运行 WebAssembly。它的异步操作不使用 libuv 这个库，而是使用 Rust 语言的 [Tokio](https://github.com/tokio-rs/tokio) 库，来实现事件循环（event loop）。

> 内置了 V8 引擎意味着可以随时随地运行那些依赖浏览器环境的库，再也不用在 js 代码中补浏览器环境了！内置了 tsc 引擎意味着再也不用面临 Node.js 那些让人头痛的无法运行 .ts 的问题。

> Deno 运行代码时，并不依赖 tsconfig.json 文件或 package.json 文件。不用面临 CommonJS 和 ES6 语法的冲突配置问题。

Deno 甚至可以通过命令行工具的 `deno compile` 命令将程序编译成可执行的二进制文件（比如 windows 下就是 .exe 文件），然后直接执行！！就像 Go 的二进制文件一样。如果是编译了一个使用浏览器操作的文件，那么就会在 Shell 中显示像浏览器似的控制，浏览器弹窗就是 Shell 中的等待输入，诸如此类的效果。

# 安装 Deno

从 [Release](https://github.com/denoland/deno/releases) 处下载 Deno 二进制文件后即可直接使用。非常简单，[deno-x86_64-pc-windows-msvc.zip](https://github.com/denoland/deno/releases/download/v1.39.4/deno-x86_64-pc-windows-msvc.zip) 就是一个非常简单直接的二进制文件，放到 $PATH 中即可直接使用。

# Deno 关联文件与配置

**deno.json** # 适用于每个项目的配置文件

- https://docs.deno.com/runtime/manual/getting_started/configuration_file 这是关于 deno.json 文件的官方文档
- https://github.com/denoland/fresh/blob/main/deno.json 典型使用示例。其中有 import map 的用法。

**${DENO_DIR}** # Deno 的缓存目录。类似于 GOPATH。可以通过 `deno info` 命令查看当前具体目录。`Windows 默认值: %LocalAppData%/deno/`

- **./deps/** # Deno 直接管理的远程模块缓存数据保存路径。通常根据协议分为 https 和 http 等目录。有点类似 go 中的 `${GOPATH}/pkg/` 目录
  - **./https/\${DOMAIN}/** # 依赖包保存在其所在域名的同名目录下。有点类似 go 中的 `${GOPATH}/pkg/mod/${DOMAIN}` 目录
- **./npm/** # Deno 适配 package.json 文件后，npm 模块缓存数据保存路径。TODO: 待确认？
- **./gen/** # Emitted 模块缓存数据保存路径。TODO: 这是啥？是构建或运行时生成的数据吗？
- **./registries/** # Language server registries cache 这是啥？
- **./location_data/** # Origin storage 这是啥？

# 模块与包

https://deno.land/x 是 Deno 第三方模块的托管仓库，就像 https://pkg.go.dev/ 。Deno 可以从网络上的任何位置导入模块，例如 GitHub，个人网络服务器或 ESM.SH，SKYPACK，JSPM.IO 或 JSDELIVR 等 CDN。为了使消费第三方模块更容易，Deno 提供了一些内置的工具，例如 Deno Info 和 Deno Doc。

https://deno.land/std 是 Deno 的标准库

Deno 有一个类似 go 一样的命令行工具 deno，内置了开发者需要的各种功能，不再需要外部工具。打包、格式清理、测试、安装、文档生成、linting、脚本编译成可执行文件等，都有专门命令。

直接运行 deno 命令，也可直接进入 [REPL](/docs/2.编程/Programming%20environment/REPL.md)。

## 兼容 Node 和 npm 模块

https://docs.deno.com/runtime/manual/node 这里可以找到 Deno 兼容 Node 和 npm 模块的方式。当前有三种方式

- 在 import 关键字语法中，使用 npm 和 node 修饰符，比如 `import CryptoJS from "npm:crypto-js"`
  - 从 1.28 版本开始，Deno 可以直接使用 Nodejs 语法的 import 以导入想要使用的包（e.g. `import CryptoJS from "npm:crypto-js"`）。
- 读取 [package.json 文件](https://docs.deno.com/runtime/manual/node/package_json)，import 关键字中的语法不用变化。但是要让 package.json 文件中包含依赖包信息
- 通过 [CDN](https://docs.deno.com/runtime/manual/node/cdns) 使用 npm 包。

通过上面前两种方法使用 Deno 运行代码时，会自动创建 node_moduels/ 目录，并将 npm 依赖保存其中

TODO: Deno 从哪里下载 npm 模块？从 https://registry.npmjs.org/？ 比如 crypto-js，从 https://registry.npmjs.org/crypto-js 下载？

# 其他

> 参考：
>
> - [官方文档，运行时-手册-CLI 命令参考](https://docs.deno.com/runtime/manual/tools/)
## compile

https://docs.deno.com/runtime/manual/tools/compiler

把脚本编译成独立的可执行文件。

## info

https://docs.deno.com/runtime/manual/tools/dependency_inspector

可用于显示有关缓存位置的信息，若指定了具体模块则可检查 ES 模块及其所有依赖项。

可以显示下面这些缓存信息

```text
DENO_DIR location: C:\Users\DesistDaydream\AppData\Local\deno
Remote modules cache: C:\Users\DesistDaydream\AppData\Local\deno\deps
npm modules cache: C:\Users\DesistDaydream\AppData\Local\deno\npm
Emitted modules cache: C:\Users\DesistDaydream\AppData\Local\deno\gen
Language server registries cache: C:\Users\DesistDaydream\AppData\Local\deno\registries
Origin storage: C:\Users\DesistDaydream\AppData\Local\deno\location_data
```
