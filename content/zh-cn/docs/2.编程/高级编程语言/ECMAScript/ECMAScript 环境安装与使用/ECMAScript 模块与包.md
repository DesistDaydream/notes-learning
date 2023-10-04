---
title: ECMAScript 模块与包
weight: 2
---

# 概述

> 参考：
>
> - [MDN-参考，JavaScript-JavaScript 指南-JavaScript 模块](https://developer.mozilla.org/en-US/docs/Web/JavaScript/Guide/Modules)
> - [网道，ES6 教程-Module 的语法](https://wangdoc.com/es6/module.html)
> - <https://www.zhangxinxu.com/wordpress/2018/08/browser-native-es6-export-import-module/>

历史上，JavaScript 一直没有 **Module(模块)** 体系，无法将一个大程序拆分成互相依赖的小文件，再用简单的方法拼装起来。其他语言都有这项功能，比如 Ruby 的 require、Python 的 import，甚至就连 CSS 都有 @import，但是 JavaScript 任何这方面的支持都没有，这对开发大型的、复杂的项目形成了巨大障碍。

在 ES6 之前，社区制定了一些模块加载方案，最主要的有 2009 年 1 月发起的 CommonJS 和 AMD 两种，前者用于服务器，后者用于浏览器。

> 2013 年 5 月，npm 的作者宣布 Node.js 已经废弃 CommonJS，详见 [GitHub issue-5132，nodejs/node-v0.x-archive](https://github.com/nodejs/node-v0.x-archive/issues/5132#issuecomment-15432598) > [Wiki，Asynchronous_module_definition](https://en.wikipedia.org/wiki/Asynchronous_module_definition)(异步模块定义，简称 AMD)

比如，CommonJS 模块就是对象，输入时必须查找对象属性。

```javascript
// CommonJS 标准
let { stat, exists, readfile } = require("fs")

// 等同于 js 代码
let _fs = require("fs")
let stat = _fs.stat
let exists = _fs.exists
let readfile = _fs.readfile
```

上面代码的实质是整体加载 fs 模块（即加载 fs 的所有方法），生成一个对象（\_fs），然后再从这个对象上面读取 3 个方法。这种加载称为“运行时加载”，因为只有运行时才能得到这个对象，导致完全没办法在编译时做“静态优化”。

## ES6 Module

> 参考：
>
> - <https://beginor.github.io/2021/08/16/using-es-modules-in-borwser-with-importmaps.html>

**ES6 Module(ES6 模块，简称 ESM)**，ES6 在语言标准的层面上，实现了模块功能，而且实现得相当简单，完全可以取代 CommonJS 和 AMD 规范，成为浏览器和服务器通用的模块解决方案。这种模块功能与 ES6 一起发布于 2015 年

ES6 模块的设计思想是尽量的静态化，使得编译时就能确定模块的依赖关系，以及输入和输出的变量。CommonJS 和 AMD 模块，都只能在运行时确定这些东西。

通常来说，**一个模块指的一组文件的合集**，只不过在通过编译工具编译后，将合并成一个文件。

ES6 模块不是对象，而是通过 export 命令显式指定输出的代码，再通过 import 命令输入。

```javascript
// ES6 模块
import { stat, exists, readFile } from "fs"
```

上面代码的实质是从 fs 模块加载 3 个方法，其他方法不加载。这种加载称为“编译时加载”或者静态加载，即 ES6 可以在编译时就完成模块加载，效率要比 CommonJS 模块的加载方式高。当然，这也导致了没法引用 ES6 模块本身，因为它不是对象。

> 只支持相对路径或者绝对路径下的 ES 模块 (./, ../, /, http://, https://) ， 同时也受服务器跨域请求策略、 HTTPS 策略的约束。

由于 ES6 模块是编译时加载，使得静态分析成为可能。有了它，就能进一步拓宽 JavaScript 的语法，比如引入宏（macro）和类型检验（type system）这些只能靠静态分析实现的功能。

除了静态加载带来的各种好处，ES6 模块还有以下好处。

- 不再需要 UMD 模块格式了，将来服务器和浏览器都会支持 ES6 模块格式。目前，通过各种工具库，其实已经做到了这一点。
- 将来浏览器的新 API 就能用模块格式提供，不再必须做成全局变量或者 navigator 对象的属性。
- 不再需要对象作为命名空间（比如 Math 对象），未来这些功能可以通过模块提供。

## import

https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Statements/import

## Import maps

<https://beginor.github.io/2021/08/16/using-es-modules-in-borwser-with-importmaps.html>

```html
<script type="importmap">
  {
    "imports": {
      "vue": "https://unpkg.com/vue@3/dist/vue.esm-browser.js"
    }
  }
</script>

<script type="module">
  import { createApp } from "vue"
</script>
```

# node_modules

[稀土掘金，node_modules 困境](https://juejin.cn/post/6914508615969669127)

- [Deno不只是个Javascript运行时](https://cloud.tencent.com/developer/article/2212864) 提到了 node_modules 困境

有了 Deno 后，应该就不用再用恶心人 node_modules 了。

# 模块的加载方式

> 参考：
>
> - <https://stackoverflow.com/questions/47403478/es6-modules-in-local-files-the-server-responded-with-a-non-javascript-mime-typ>
> - <https://localcoder.org/es6-modules-in-local-files-the-server-responded-with-a-non-javascript-mime-typ>
> - <https://axellarsson.com/blog/expected-javascript-module-script-server-response-mimetype-text-html/>

在 Node.js 环境和 Browser 环境中加载 ESM 的方式不太一样

- **Node.js** # 可以使用模块名称。从根目录下的 node_modules/ 中查找模块
- **Browser** # 不可以使用模块名称。必须通过编译工具将模块编译成单一文件，并修改 import 指向单一文件，以便可以发起请求获取这个静态资源

## 浏览器中使用 ESM 的常见问题

使用 `import * as Vue from 'vue'` 将会产生如下报错：

`Failed to resolve module specifier "vue". Relative references must start with either "/", "./", or "../".`

![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/mm0ymr/1651724399014-d2052b6f-cd7c-4ec0-b6fc-b748bd5a11ed.png)
接着修改为 `import * as Vue from '../node_modules/vue'` 将会产生如下报错：
`Failed to load module script: The server responded with a non-JavaScript MIME type of "text/html". Strict MIME type checking is enforced for module scripts per HTML spec.`

> ![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/mm0ymr/1651724430503-b62b86bd-4cc7-48b8-ac73-69fa62564ed5.png "firefox")

![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/mm0ymr/1651724407422-526db251-775f-40d5-a25e-402791aa38cc.png "chrome")

问题原因：

这个情况的原因是浏览器在处理 import 逻辑时导致的。浏览器在发现 import 语句时，将会请求 from 后面的静态文件，当 from 指定的是模块名称而不是模块文件的路径时时，浏览器无法发起请求，因为浏览器不知道如何获取到模块文件。

解决方式：

- 使用配置文件告诉 JavaScript 引擎如何从名为 XX 的模块中获取模块文件
- 打包代码，js 代码被打包后，导入的模块的一组文件将会被制作为一个或多个文件；并且 from 后面的模块名将被修改为模块文件的路径，即可在浏览器中运行

这几种解决方式通常都是通过编译工具实现的，比如 Webpack、Vite 等工具。

# ECMAScript 模块的使用与管理

> 参考：
>
> - [官方文档，cli-配置 npm-文件夹](https://docs.npmjs.com/cli/v8/configuring-npm/folders)

ECMAScript 的模块与包相关概念与工具是相当混乱的，最早期是由 Node.js 安装时自带的 NPM 工具集进行管理，并且 NPM 工具集中的工具除了模块与包的管理，还可以提供运行时环境等功能。

在后期，出现了 yarn、pnpm 等新型的管理工具，可以通过 Node.js 自带的 `corepack enable` 命令启用这些新的包管理工具。

ECMAScript 的包管理器安装各种包、依赖时，早期都是分开的，可以安装在项目当前目录，或某一个统一目录。但是如果项目想要正常启动，一般都只能安装在项目的当前目录（历史原因已不可考，反正很恶心）。

后来出现的 pnpm 工具，可以让我们将各种不同的项目的依赖放在同一个路径下进行统一管理。

## 安装 TypeScript 第三方模块

如果使用 `pnpm install crypto-js` 这种命令安装的模块是 JS 代码，有些第三方库并没有提供 TypeScript 类型声明文件，这会导致 TypeScript 在编译时无法识别该库的类型信息。可以通过安装类型声明文件解决该问题。类型声明文件通常以 `.d.ts` 结尾，存放在 `@types` 组织中。

所以，为了解决这个问题，我们安装第三方库是，一般库名前都要加一个 `@types`，以标识该库要从 @types 中拉取，比如上面的 crypto-js 库，应该使用如下命令拉取 TS 版本的

```
pnpm install @types/crypto-js
```

Notes: 有的时候这种 TS 依赖库还依赖原始的 JS 库，也要同时安装 JS 库。这个 crypto-js 就是这个情况，要想使用 ts-node-esm 正常运行代码， package.json 至少需要如下内容：

```json
{
 "type": "module",
 "dependencies": {
  "@types/crypto-js": "^4.2.1",
  "crypto-js": "^4.2.0"
 }
}
```

## 模块与包的存储路径

npm 工具会通过 ${PREFIX} 与 node_modules/ 组合来决定其所管理的各种依赖包应该保存在什么位置。

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

> 注意：Windows 的全局 `node_modules/` 目录与 Linux 不太一样，全局路径是 `${PREFIX}/node_modeuls/`。也就是说生成的链接文件就在 `${PREFIX}` 下。

