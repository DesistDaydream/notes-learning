---
title: "Hugo"
linkTitle: "Hugo"
weight: 20
---

# 概述

> 参考：
> - [GitHub 项目，gohugoio/hugo](https://github.com/gohugoio/hugo)
> - [官网](https://gohugo.io/)


## 目录结构
> 参考：
> - [官方文档，入门-目录结构](https://gohugo.io/getting-started/directory-structure/)

`hugo new site example` 命令将会创建一个包含以下元素的目录结构：
```
example/
├── archetypes/
│   └── default.md
├── assets/
├── content/
├── data/
├── layouts/
├── public/
├── static/
├── themes/
└── config.toml
```
archetypes/
- default.md

assets/ # 存放所有需要由 Hugo Pipes 处理的文件。

content/

data/

layouts/

public/

static/

themes/

config.toml


# Hugo Modules(模块)
Hugo 模块是一个类似 Go 模块一样的存在。模块可以我们的主项目或较小的模块，提供 Hugo 的 7 种组件类型中的一种或多种：
- 

在 config.toml 文件中的 module 字段添加配置，即可为站点设置引用的模块，我们可以将主题当做一个模块，



# 安装 Hugo


# 关联文件与配置
config.toml | config.yaml | config.json # 站点的配置文件

Hugo 运行时所需的缓存目录。包括需要使用的模块等：
- Windows:
	- **%TMP%/hugo_cache/\***
- Linux:
	- **${TMP}/hugo_cache/\***


