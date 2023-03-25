---
title: Vue 环境安装与使用
---

# Vue 命令行工具

> 参考：
> - [官网](https://cli.vuejs.org/)

安装

```bash
npm install -g @vue/cli
# OR
yarn global add @vue/cli
```

配置环境变量

# 创建 Vue 项目

> 参考：
> 
> - [官方文档，快速上手-创建一个 Vue 应用](https://cn.vuejs.org/guide/quick-start.html#creating-a-vue-application)
> - [GitHub 项目，vitejs/awesome-vite](https://github.com/vitejs/awesome-vite)（一些使用 vite 创建的应用模板，可以直接拿来用）
>   - [Vue Naive](https://github.com/zclzone/vue-naive-admin) - 管理模板，基于 Vue 3 + Pinia + Naive UI。

可以通过两种方式创建一个 Vue3 的项目

- **vue 工具** # `vue create NAME` 基于 webpack。不推荐使用。
- **vite 工具** # 通过 [vuejs/create-vue](https://github.com/vuejs/create-vue) 项目，基于 vite 创建 Vue 项目。强烈推荐

## 基于 Vite 创建 Vue 项目

```bash
npm init vue@latest
```

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
