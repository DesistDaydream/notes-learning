---
title: "Hugo"
linkTitle: "Hugo"
weight: 20
---

# 概述

> 参考：
> - [GitHub 项目，gohugoio/hugo](https://github.com/gohugoio/hugo)
> - [官网](https://gohugo.io/)
> - [Wiki,Hugo(软件)](https://en.wikipedia.org/wiki/Hugo_(software))

Hugo 是用 Go 语言编写的静态站点生成器。Steve Francia 最初于 2013 年将 Hugo 创建为开源项目。

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
### archetypes

### assets
存放所有需要由 Hugo Pipes 处理的文件。

### config
> 注意：config 目录并不会随着 `hugo new site example` 命令自动创建出来。而是在我们需要将单一配置文件拆分时，手动创建并使用的。

Hugo 附带了大量的配置指令。config 目录是将这些指令存储为 JSON、YAML 或 TOML 文件的地方。每个根设置对象都可以作为自己的文件并按环境构建。设置最少且不需要环境的项目可以在其根目录下使用单个 config.toml 文件。

许多站点可能几乎不需要配置，但 Hugo 附带了大量 [configuration directives(配置指令)](https://gohugo.io/getting-started/configuration/#all-configuration-settings)，用于更详细地指导我们希望 Hugo 如何构建网站。注意：默认情况下不创建 config 目录。

### content/
我们使用 Hugo 创建的网站的所有内容通常都要放在 content 目录中。content 目录下的每个顶级文件夹称为 [content section(内容部分)](https://gohugo.io/content-management/sections/)。

比如，如果我的网站有三大块，分别是 blog、articles、tuorials，那么我们需要创建如下的目录结构
```
example/
├── content/
│   ├── blog/
│   ├── articles/
│   └── tuorials/
```
Hugo 使用 section 的名称作为默认的 [content types(内容类型)](https://gohugo.io/content-management/types/)。假如有这么一个文件 content/blog/my-first-event.md，则这篇文章的内容类型就是 blog 类型。

通过这种对网站内容的分类方式，更利于搜索、整理等。

### data/
该目录用于存放 Hugo 在生成我的网站时可以使用的配置文件。可以用 YAML、JSON 或 TOML 格式编写这些文件。除了添加到此文件夹的文件外，还可以创建从动态内容中提取的数据模板。

### layouts/
以 .html 文件的形式存储模板，这些文件指定如何将您的内容视图呈现到静态网站中。模板包括 [list pages](https://gohugo.io/templates/list/), your [homepage](https://gohugo.io/templates/homepage/), [taxonomy templates](https://gohugo.io/templates/taxonomy-templates/), [partials](https://gohugo.io/templates/partials/), [single page templates](https://gohugo.io/templates/single-page-templates/),等

### public/
使用 `hugo` 命令生成网站的静态文件后，将会保存到 public 目录。public 目录生成的静态文件，可以直接通过 web 服务访问到。

### static/
存储所有静态内容：图像、CSS、JavaScript 等。当 Hugo 构建您的站点时，静态目录中的所有资产都会按原样复制。使用静态文件夹的一个很好的例子是在 Google Search Console 上验证网站所有权，您希望 Hugo 在不修改其内容的情况下复制完整的 HTML 文件。

### themes/
> 更推荐的是使用 Hugo 模块使用主题，该目录不推荐使用。

Hugo 主题可以安装到该目录，使用 `hugo server --themes` 指定使用的主题时，将会从该目录出寻找。


### config.toml
Hugo 运行站点时所使用的配置文件。

推荐使用 config/ 目录，以便拆分 config.toml 文件。可以将 config.tom 文件移动到在 config/\_default/config.toml 处作为默认配置。


# Hugo Modules(模块)
Hugo 模块是一个类似 Go 模块一样的存在。模块可以我们的主项目或较小的模块，提供 Hugo 的 7 种组件类型中的一种或多种：
- 

在 config.toml 文件中的 module 字段添加配置，即可为站点设置引用的模块，我们可以将主题当做一个模块，



# 安装 Hugo


# 关联文件与配置
**config.toml | config.yaml | config.json** # 站点的配置文件，通常在站点的根目录
**${Site_Root_dir}/config** # 可以将站点根目录下的 config.toml | config.yaml | config.json 拆分后保存到该目录。

Hugo 运行时所需的缓存目录。包括需要使用的模块等：
- Windows:
	- **%TMP%/hugo_cache/\***
- Linux:
	- **${TMP}/hugo_cache/\***


