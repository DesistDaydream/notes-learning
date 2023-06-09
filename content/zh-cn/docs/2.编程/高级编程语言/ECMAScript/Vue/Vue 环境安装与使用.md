---
title: Vue 环境安装与使用
---

# 概述

> 参考：
> 
> - [官方文档，应用规模化-工具链](https://cn.vuejs.org/guide/scaling-up/tooling.html)

Vue3 早期使用 [Vue CLI](https://cli.vuejs.org/)，创建 Vue 项目，后来 Vue 作者尤雨溪开了一个新的工具 Vite，Vite 通过 [vuejs/create-vue](https://github.com/vuejs/create-vue) 项目，基于 Vite 创建 Vue 项目。

# Vite

> 参考：
> 
> - [GitHub 项目，vitejs/vite](https://github.com/vitejs/vite)
> - [官网](https://vitejs.dev/)

Vite 是一种新型前端构建工具，可显着改善前端开发体验。它由两个主要部分组成：

-   一个开发服务器，通过[原生 ES 模块](https://developer.mozilla.org/en-US/docs/Web/JavaScript/Guide/Modules)为您的源文件提供服务，具有[丰富的内置功能](https://vitejs.dev/guide/features.html)和惊人的快速[热模块替换 (HMR)](https://vitejs.dev/guide/features.html#hot-module-replacement)。
-   一个[构建命令](https://vitejs.dev/guide/build.html)，将你的代码与[Rollup](https://rollupjs.org/)捆绑在一起，预先配置为输出高度优化的生产静态资产。

此外，Vite 通过其[插件 API](https://vitejs.dev/guide/api-plugin.html)和[JavaScript API](https://vitejs.dev/guide/api-javascript.html)具有高度的可扩展性，具有完整的类型支持。

## 基于 Vite 创建 Vue 项目

> 参考：
> 
> - [官方文档，快速上手-创建一个 Vue 应用](https://cn.vuejs.org/guide/quick-start.html#creating-a-vue-application)
> - [GitHub 项目，vitejs/awesome-vite](https://github.com/vitejs/awesome-vite)（一些使用 vite 创建的应用模板，可以直接拿来用）
>     - [Vue Naive](https://github.com/zclzone/vue-naive-admin) - 管理模板，基于 Vue 3 + Pinia + Naive UI。

```bash
npm init vue@latest
```

> 注意：通过 npm 将 vite 作为模块安装到 node_modules/ 目录下，然后执行 npm run dev、npm build 等命令时，可以直接调用。但是想在 CLI 直接调用 vite 命令是需要通过 `npm install -g vite` 单独安装 vite 命令行工具的
> 注意：如果不在全局安装 vite，那也可以直接使用 `node ./node_moduels/vite/bin/vite.js` 运行在本项目依赖中安装的 vite。

该指令会安装并执行 [create-vue](https://github.com/vuejs/create-vue)(这是 Vue 官方的项目脚手架工具)。我们会看到如下的可选功能：

```bash
✔ Project name: … <your-project-name>
✔ Add TypeScript? … No / Yes
✔ Add JSX Support? … No / Yes
✔ Add Vue Router for Single Page Application development? … No / Yes
✔ Add Pinia for state management? … No / Yes
✔ Add Vitest for Unit testing? … No / Yes
✔ Add Cypress for both Unit and End-to-End testing? … No / Yes
✔ Add ESLint for code quality? … No / Yes
✔ Add Prettier for code formatting? … No / Yes

Scaffolding project in ./<your-project-name>...
Done.
```

如果不确定是否要开启某个功能，你可以直接按下回车键选择 No。在项目被创建后，通过以下步骤安装依赖并启动开发服务器：

```bash
> cd <your-project-name>
> npm install
> npm run dev
```

`npm install` 用以安装本项目的依赖，`npm run dev` 则是执行 `vite` 命令，此时 Vite 将会启动一个监听程序用以响应 html 文件。

# IDE 插件

## Volar

![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/bvqpem/1667541398977-a8c81df3-0834-4a86-9179-c50e3c9c9c20.png)
