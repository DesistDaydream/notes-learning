---
title: Vite
---

# 概述

> 参考：
> 
> - [GitHub 组织，vite](https://github.com/vitejs)
> - [GitHub 项目，vitejs/vite](https://github.com/vitejs/vite)

Vue 作者主导的打包与编译工具，启动项目开发模式速度非常快、下载依赖非常快、编译非常快。

在 [GitHub 项目，vitejs/awesom-vite](https://github.com/vitejs/awesome-vite) 这里有很多使用 Vite 打包的很棒的项目示例

# 安装 Vite

# Vite 关联文件与配置

**vite.config.js** # 每个项目的根目录通常都会有一个 vite.config.js 文件，以定义打包项目代码时的行为

### Vite 配置详解

> 参考：
> 
> - [官方文档，配置](https://cn.vitejs.dev/config/)

# 常见问题

## 解决 Vite 打包项目代码后，使用的是绝对路径

在 vite.config.js 文件中设置 base

## 解决 vite 在 WSL 环境下热更新失效问题

使用 vite 的好处很多，最明显的就是热更新很快。但是在 wsl 环境的时候，由于[WSL2 的限制](https://github.com/microsoft/WSL/issues/4739)，vite 默认配置是无法监控 windows 文件系统中文件的变化的。这就导致了 vite 的热更新失效。
解决热更新失效的方法主要有两种： 1.文件存储到 WSL 系统环境中 2.配置 vite.config.js 的 [server.watch](https://cn.vitejs.dev/config/#server-watch)

```javascript
export default defineConfig({
  server: {
    watch: {
      { usePolling: true }
    }
  }
})
```

# vite Syntax(语法)

vite 将会启动一个开发服务器，默认响应当前目录的 index.html 文件。

COMMAND:

- build \[root] # 构建项目(生产可用)
- optimize \[root] # pre-bundle dependencies
- preview \[root] # locally preview production build

## build

## preview

在本地预览已经通过 `vite build` 构建后的项目。注意：执行该命令前，需要先通过 `vite build` 打包项目，以便生成 `dist/` 目录
启动监听当前目录下的 `dist/` 目录。
