---
title: Puppeteer
linkTitle: Puppeteer
weight: 20
---

# 概述

> 参考：
>
> - [GitHub 项目，puppeteer/puppeteer](https://github.com/puppeteer/puppeteer)
> - [官网](https://pptr.dev/)
>   - [中文文档？](https://puppeteer.bootcss.com/)
> - [10分钟快速上手爬虫之Puppeteer](https://www.bilibili.com/video/BV13Z4y137Kt)

Puppeteer 是一个 Node.js 库，它提供了一个高级 API 来通过 DevTools 协议控制 Chrome/Chromium。 Puppeteer 默认以无头模式运行，但可以配置为在完整（“有头”）Chrome/Chromium 中运行。

> 注意：Puppeteer 并不是一个 WebDriver 的实现。

# Puppetter 安装

pnpm install puppetter

安装 Puppetter 时，会在 `${HOME}/.cache/puppeteer/` 目录下安装 Chrome 程序。

注意：若是 pnpm 安装，则需要删除原始的文件，只取消项目下的链接后执行 `pnpm install puppetter` 并不会在没有 Chrome 程序时执行下载逻辑。

# 用法

## 元素

https://www.cnblogs.com/totoro-cat/p/11310832.html

定位元素返回 ElementHandle 实例，然后使用 ElementHandle 下的方法处理元素
