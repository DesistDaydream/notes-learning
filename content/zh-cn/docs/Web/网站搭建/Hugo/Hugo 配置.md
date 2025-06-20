---
title: "Hugo 配置"
linkTitle: "Hugo 配置"
weight: 20
---

# 概述
>
> 参考：
>
> - [官方文档，入门 - 配置](https://gohugo.io/getting-started/configuration/)

Hugo 支持 [TOML](/docs/2.编程/无法分类的语言/TOML.md)、[YAML](/docs/2.编程/无法分类的语言/YAML.md)、[JSON](/docs/2.编程/无法分类的语言/JSON.md) 格式的配置文件，默认配置文件名为 **hugo.SUFFIX**。所有的配置指令都可以写在 hugo.SUFFIX 文件中。

我们可以拆分配置文件，并将其保存在 `${ConfigDir}` 目录下(默认为站点根目录下的 `config/` 文件夹)。 ConfigDir 下的每个文件名代表配置中的根指令，比如：

hugo.toml 文件中有如下指令

```toml
[Params]
foo = 'bar'
```

那么拆分到 ConfigDir 目录时，则使用 `params.toml` 文件，内容为：

```toml
foo = 'bar'
```

除了 Hugo 本身会用到配置文件，有些主题也会使用，并具有各自可以识别的配置指令。比如 Docsy 主题。

在[官方文档，配置-所有配置设置](https://gohugo.io/getting-started/configuration/#all-configuration-settings)我们可以找到 Hugo 的所有配置指令

## config 目录结构

```bash
├── config
│   ├── _default
│   │   ├── hugo.toml
│   │   ├── languages.toml
│   │   ├── menus.en.toml
│   │   ├── menus.zh.toml
│   │   └── params.toml
│   ├── production
│   │   ├── hugo.toml
│   │   └── params.toml
│   └── staging
│       ├── hugo.toml
│       └── params.toml
```

`_default/` 目录是站点的默认配置，可以直接保存单个 config 文件。

production 与 staging 文件夹用来区分运行时配置，比如使用 `hugo --environment staging` 命令时，Hugo 将会使用 `config/_default/` 和 `config/staging` 这两个下的所有文件，将所有文件合并后生成一个单独的临时 config 文件，作为站点运行时的配置文件。

这种目录结构，可以帮助我们区分不同环境(比如开发环境、生产环境等)下运行网站所需要使用的配置。

> Hugo 有个默认值，执行 `hugo server` 命令时，是在本地运行网站，所以默认为开发环境，使用 `config/development/` 目录；而执行 `hugo` 命令时，是要构建静态文件，所以默认为生产环境，使用 `config/production` 目录。

# 基础

**baseURL**(STRING) # 我们发布的网站的绝对 URL（协议，主机，路径和斜杠），比如 https://www.example.org/docs/

**title**(STRING) #

**contentDir = "content/zh-cn"** # **必须的**。Hugo 从该配置指定的目录中读取内容文件。可以为各个语言单独配置。`默认值：content`

**defaultContentLanguage = "zh-cn"**

**defaultContentLanguageInSubdir = false**

**enableGitInfo** # `默认值：false`

- 注意：[issue 3071](https://github.com/gohugoio/hugo/issues/3071)，如果文件名、目录名是中文的话，将会无法获取到 Git 信息。也就无法让文档下面显示最后编辑时间。需要修改 git 的配置 core.quotePath 的值为 false。

**markup**([markup](#markup%20部分))

## URL 生成逻辑控制

**disablePathToLower**(BOOLEAN) # 是否关闭 URL 转换为小写字母的逻辑

**uglyURLs**(BOOLEAN) # URL 路径中是否要带 .html 后缀

# markup 部分

markup 部分的配置用于处理 Markdown 和其他 Markup(标记) 相关配置。

**\[goldmark]**

Goldmark 部分用于配置适用于 Go 的 Markdown 解析库，Hugo 从 0.60 开始使用。它速度快，符合 CommonMark 标准，而且非常灵活。

**\[parser.attribute]**

**block**(BOOLEAN) # 是否为 block 启用 [Markdown 属性](https://gohugo.io/content-management/markdown-attributes/)

**\[renderer]**

**unsafe**(BOOLEAN) # 是否让 Goldmark 渲染器将在 Markdown 中渲染原始 HTML。开启后是不安全的。`默认值: false`。

**\[tableOfContents]**

tableOfContents 部分配置目录相关指令。这些设定只适用于 Goldmark 渲染器。（这个目录是指文章的 *大纲*，并不是文件系统中的目录）

**startLevel**(INT) # 目录中显示的标题级别。从该指令的级别开始显示

**endLevel**(INT) # 目录中显示的标题级别。到该指令的级别结束显示

**ordered**(BOOLEAN) # 是否生成有序列表而不是无序列表。

**\[highlight]**

highlight 部分配置高亮部分（通常都是代码块）

**style**(STRING) # 高亮部分的样式

# module 部分

> 参考：
>
> - [官方文档，模块-配置](https://gohugo.io/hugo-modules/configuration/)

module 部分的配置用于处理 Hugo 模块的运行逻辑。

**proxy = \<STRING>** # 定义用于下载模块的代理服务器。与 go 模块的 proxy 原理一样。`默认值：direct`

**\[\[imports]]**

**path**(STRING) # Hugo 指定 Hugo 要使用的的 [Go Module](/docs/2.编程/高级编程语言/Go/Go%20环境安装与使用/Go%20Module.md)。值是标准的 Go 模块路径，可以是网络上的，也可以是本地的。

# sitemap 部分

> 参考：
>
> - [官方文档，模板 - 站点地图模板](https://gohugo.io/templates/sitemap-template/)

**filename = \<STRING>** # 生成的 sitemap 文件的名称。`默认值：sitemap.xml`

# params 部分

这部分的配置通常用来配置各种 Hugo 主题。这里面的指令通常会被主题读取并生成样式。
