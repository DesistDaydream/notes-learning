---
title: Vue
---

# 概述

> 参考：
> - [GitHub 组织，vuejs](https://github.com/vuejs/)
> - [官网](https://vuejs.org/)
> - [官网-中文](https://staging-cn.vuejs.org/)
> - [Vue 互动教程](https://cn.vuejs.org/tutorial/)
> - [Wiki,Vue.js](https://en.wikipedia.org/wiki/Vue.js)

Vue 是一套用于构建用户界面的渐进式 ECMAScript 框架。Vue3 于 2020 年 9 月发布，已全面采用 TypeScript 编写；在 2022 年 2 月份成为默认版本

## 组件化

Vue 是“组件化”模式，一个页面的各个部分，可以拆分成一个一个的组件：

![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/cye267/1651220067462-1822075c-2b85-4cf4-abd8-eebfa658e531.png)

**Single-File Component(单文件组件，简称 SFC)**。顾名思义，Vue 的单文件组件会将一个组件的逻辑 (JavaScript)，模板 (HTML) 和样式 (CSS) 封装在同一个文件里。

**同时，多个组件可以自由组合拼接，形成一个完整的页面。**

单文件组件是 Vue 的标志性功能。如果你的用例需要进行构建，我们推荐用它来编写 Vue 组件。你可以在后续相关章节里了解更多关于[单文件组件的用法及用途](https://cn.vuejs.org/guide/scaling-up/sfc.html)。但你暂时只需要知道 Vue 会帮忙处理所有这些构建工具的配置就好。

这些组件通常被组织在 **XXX.vue** 文件中，通常保存在项目根目录的 `components/` 目录下。

组件化开发是一个树状结构，从一个“根组件”开始：

```bash
App (root component)
├─ TodoList
│  └─ TodoItem
│     ├─ TodoDeleteButton
│     └─ TodoEditButton
└─ TodoFooter
   ├─ TodoClearButton
   └─ TodoStatistics
```

就像下面这样：

![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/cye267/1666837624381-ce56eb28-5092-4e8a-a1c8-de01ed1e1f7f.png)

## 声明式

![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/cye267/1651220139947-ac46307c-fa52-4370-b7ed-16e2ad92629a.png)
![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/cye267/1666836325564-89bbb0d6-be23-45ea-aa56-bc7183a2eb46.png)
有了 VUE 就不用手动操作 DOM 了？

## API 风格

**选项式 API** 与 **组合式 API**。推荐组合式。

> 参考：
> - <https://www.bilibili.com/video/BV1mK411f7Kt/?p=51>

选项式：

```typescript
export default {
  data() {
    return {
      count: 1,
    }
  },

  // `mounted` 是生命周期钩子，之后我们会讲到
  mounted() {
    // `this` 指向当前组件实例
    console.log(this.count) // => 1

    // 数据属性也可以被更改
    this.count = 2
  },
}
```

组合式：

```typescript
import { reactive } from "vue"

export default {
  // `setup` 是一个专门用于组合式 API 的特殊钩子函数
  setup() {
    const state = reactive({ count: 0 })

    // 暴露 state 到模板
    return {
      state,
    }
  },
}
```

选项式 API 将需要处理的数据放在 `data()` 中 ，关于处理数据的逻辑写在 `methods:`、`computed:`、`watch:` 等等地方，如果数据很多，那么处理数据的逻辑在编辑器中将会非常跳跃，就像下图左侧一样，同样颜色的逻辑，不够集中，那么将会形成非常乱的代码结构。

![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/cye267/1667871595181-70833fc5-41d3-48f1-954f-16c912da8749.png)
所以说，组合式，主要组合的是对于数据处理的逻辑，将处理同一个数据的逻辑组合在一起，以便编写出更易读的代码。

# Vue 指令

指令是带有 `v-` 前缀的特殊 attribute。Vue 提供了许多[内置指令](https://cn.vuejs.org/api/built-in-directives.html)。

# Vue 规范

## 项目结构

使用 Vite 创建的 Vue 项目

```bash
.
├── README.md
├── index.html
├── node_modules
├── package.json
├── pnpm-lock.yaml
├── public
│   └── favicon.ico
├── src
│   ├── App.vue
│   ├── assets/
│   ├── components/
│   └── main.js
└── vite.config.js
```

- public/ 目录存放公共资源
- 所有代码都在 src/ 目录下
  - index.html 指向 main.js，main.js 中创建应用的根组件
  - 根组件的代码在 App.vue 文件中
  - 所有根组件下的组件的代码都放在 components/ 目录下
  - assets 存放静态资源，图片、css 样式 等等
- **.eslintrc.cjs** # ESLint 程序配置
- **.prettierc.json** # Prettier 插件的配置
- **env.d.ts** #
- **index.html** # 程序入口
- **package.json** # 包管理器配置文件，比如 npm、pnpm 等
- **vite.config.ts** # Vite 程序给项目打包时使用的配置
- **tsconfig.json** #
- **tsconfig.config.json** #

# 学习资料

[Swiperjs-Vue](https://swiperjs.com/vue)

## 待学习

[B 站，Vue3.2 + Vite + Element-Plus 实现最基础的 CRUD](https://www.bilibili.com/video/BV1yV4y177jC)

[B 站，Vue3 项目实战、Vue3+Element-plus 项目实战系列课程（数据管理平台）](https://www.bilibili.com/video/BV1sP4y127Re)

[B 站，Vue3.2 后台管理系统](https://www.bilibili.com/video/BV1pq4y1c7oy)

[B 站，一天之内快速搭建 vue 后台管理系统-代码写到起飞,接单接到手软](https://www.bilibili.com/video/BV1md4y1C7wS)

## 可以学习的项目

<https://github.com/HalseySpicy/Geeker-Admin> # Geeker Admin，基于 Vue3.2、TypeScript、Vite2、Pinia、Element-Plus 开源的一套后台管理框架。

![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/cye267/1667712440759-38351016-d3de-4246-96ce-f139fb24099c.png)

<https://github.com/flipped-aurora/gin-vue-admin> # 基于 vite+vue3+gin 搭建的开发基础平台（支持 TS,JS 混用），集成 jwt 鉴权，权限管理，动态路由，显隐可控组件，分页封装，多点登录拦截，资源权限，上传下载，代码生成器，表单生成器等开发必备功能。

![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/cye267/1667712396903-5c478114-4b5d-42c8-9caf-e301ac58f2fc.png)

<https://github.com/go-admin-team/go-admin> # 基于 Gin + Vue + Element UI 的前后端分离权限管理系统脚手架（包含了：多租户的支持，基础用户管理功能，jwt 鉴权，代码生成器，RBAC 资源控制，表单构建，定时任务等）3 分钟构建自己的中后台项目；

![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/cye267/1667712379792-273c7485-92f4-46ed-9a98-65e745b1c8df.png)

