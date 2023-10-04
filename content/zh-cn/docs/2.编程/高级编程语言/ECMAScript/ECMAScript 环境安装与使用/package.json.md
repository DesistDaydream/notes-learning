---
title: package.json
weight: 4
---

## 概述

> 参考：
> 
> - [dev 官网，学习-package.json 指南](https://nodejs.dev/learn/the-package-json-guide)
> - [pnpm 官方文档，配置-package.json](https://pnpm.io/package_json)

package.json 文件是项目的清单。 它包含包的所有元数据，包括依赖项、标题、作者等等。例如，它是用于工具的配置中心。 它也是 npm 和 yarn 等包管理工具管理依赖的地方。

想要运行带有 ES6 语法规则的代码（比如导入包是使用的 import 关键字），需要添加 `"type": "module"` 配置。
