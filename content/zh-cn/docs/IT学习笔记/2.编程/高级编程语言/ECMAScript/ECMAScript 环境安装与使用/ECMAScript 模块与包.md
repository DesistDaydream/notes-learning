---
title: ECMAScript 模块与包
weight: 2
---

# 概述

> 参考：
> - [MDN-参考，JavaScript-JavaScript 指南-JavaScript 模块](https://developer.mozilla.org/en-US/docs/Web/JavaScript/Guide/Modules)
> - [网道，ES6 教程-Module 的语法](https://wangdoc.com/es6/module.html)
> - <https://www.zhangxinxu.com/wordpress/2018/08/browser-native-es6-export-import-module/>

历史上，JavaScript 一直没有 Module(模块) 体系，无法将一个大程序拆分成互相依赖的小文件，再用简单的方法拼装起来。其他语言都有这项功能，比如 Ruby 的 require、Python 的 import，甚至就连 CSS 都有@import，但是 JavaScript 任何这方面的支持都没有，这对开发大型的、复杂的项目形成了巨大障碍。

在 ES6 之前，社区制定了一些模块加载方案，最主要的有 2009 年 1 月发起的 CommonJS 和 AMD 两种，前者用于服务器，后者用于浏览器。

> 2013 年 5 月，npm 的作者宣布 Node.js 已经废弃 CommonJS，详见 [GitHub issue-5132，nodejs/node-v0.x-archive](https://github.com/nodejs/node-v0.x-archive/issues/5132#issuecomment-15432598) > [Wiki,Asynchronous_module_definition](https://en.wikipedia.org/wiki/Asynchronous_module_definition)(异步模块定义，简称 AMD)

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

# 模块的加载方式

> 参考：
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
- 编译代码，js 代码被编译后，导入的模块的一组文件将会被打包、压缩为一个文件；并且 from 后面的模块名将被修改为模块文件的路径，即可在浏览器中运行

这几种解决方式通常都是通过编译工具实现的，比如 Webpack、Vite 等工具。
