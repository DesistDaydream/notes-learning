---
title: TSConfig
weight: 20
---

# 概述

> 参考：
> 
> - [官方文档，项目配置-tsconfig.json 是什么](https://www.typescriptlang.org/docs/handbook/tsconfig-json.html)
> - [官方文档，TSConfig 参考](https://www.typescriptlang.org/tsconfig)

TSConfig 通常是名为 tsconfig.json 或 jsconfig.json 的文件，当目录中出现了 tsconfig.json 文件，则说明该目录是 TypeScript 项目的根目录。tsconfig.json 文件指定了编译项目所需的根目录下的文件以及编译选项。

## 简单示例

可以运行 ES6 语法（导入包时用的 import 关键字）逻辑的 TS 代码的配置

```json
{
  "compilerOptions": {
    // "target": "es2016",
    "module": "ES6",
    "esModuleInterop": true,
  }
}
```

注意，若环境中有 package.json 文件，需要搭配该文件中的 `"type": "module"` 配置，才可以正常使用 ES6 语法。

# compilerOptions

## baseUrl

## paths

配置路径别名。

若使用 Vite 打包代码，则需要在 vite.config.ts 文件中也同步配置 `resolve.alias`：

```typescript
export default defineConfig({
  resolve: {
    alias: {
      // 让我们在导入时使用可以使用 @ 符号作为 src 目录的别名，而不是相对路径，比如：
      // import App from '@/App.vue'
      // 而不是
      // import App from '../../App.vue'
      "@": fileURLToPath(new URL("./src", import.meta.url)),
    },
  },
})
```

# 最佳实践

## Vue 环境配置示例

```json
{
  "extends": "@vue/tsconfig/tsconfig.web.json",
  "include": [
    "env.d.ts",
    // 两个 * 表示任意层数的所有目录
    "src/**/*",
    "src/**/*.vue",
    // 导出所有接口目录到全局
    "src/api/**/*"
  ],
  "compilerOptions": {
    "baseUrl": ".",
    // 配置路径别名。这样在使用 import 导入时，我们可以通过 @/xxx 来代替 /src/xxx
    "paths": {
      "@/*": ["src/*"]
    }
  },
  "references": [
    {
      "path": "./tsconfig.config.json"
    }
  ]
}
```

**include** # 通过在 include 中指定文件或目录，我们可以将声明的 interface 等默认导出并在其他地方使用，而不用在每个需要使用的地方使用 import 关键字显式导入了。

