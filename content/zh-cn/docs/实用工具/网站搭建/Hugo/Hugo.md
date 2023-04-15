---
title: "Hugo"
linkTitle: "Hugo"
weight: 1
---

# 概述

> 参考：
> 
> - [GitHub 项目，gohugoio/hugo](https://github.com/gohugoio/hugo)
>   - [GitHub 项目，coderzh/gohugo.org](https://github.com/coderzh/gohugo.org)（一个从19年停更的 Hugo 中文网）
> - [官网](https://gohugo.io/)
> - [Wiki,Hugo(软件)](https://en.wikipedia.org/wiki/Hugo_(software))

Hugo 是用 Go 语言编写的静态站点生成器。Steve Francia 最初于 2013 年将 Hugo 创建为开源项目。

# Hugo 的基本使用

> 参考：
> 
> - [自定义hugo主题--从内容开始](https://hugo.aiaide.com/post/%E8%87%AA%E5%AE%9A%E4%B9%89hugo%E4%B8%BB%E9%A2%98-%E4%BB%8E%E5%86%85%E5%AE%B9%E9%A1%B5%E5%BC%80%E5%A7%8B/)

`hugo new site hello_world` 命令将会创建一个包含以下元素的目录结构，这些目录的作用可以在[下文](#目录结构)找到：

```
hello_world/
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

进入到这个目录之后，执行 `hugo server` 将会启动一个 HTTP 服务端

```bash
~]# hugo server --buildDrafts
Start building sites … 
hugo v0.109.0-47b12b83e636224e5e601813ff3e6790c191e371+extended windows/amd64 BuildDate=2022-12-23T10:38:11Z VendorInfo=gohugoio
WARN 2023/01/22 16:15:45 found no layout file for "HTML" for kind "home": You should create a template file which matches Hugo Layouts Lookup Rules for this combination.
WARN 2023/01/22 16:15:45 found no layout file for "HTML" for kind "taxonomy": You should create a template file which matches Hugo Layouts Lookup Rules for this combination.
WARN 2023/01/22 16:15:45 found no layout file for "HTML" for kind "taxonomy": You should create a template file which matches Hugo Layouts Lookup Rules for this combination.

                   | EN
-------------------+-----
  Pages            |  3
  Paginator pages  |  0
  Non-page files   |  0
  Static files     |  0
  Processed images |  0
  Aliases          |  0
  Sitemaps         |  1
  Cleaned          |  0

Built in 44 ms
Watching for changes in D:\Projects\DesistDaydream\hugo-learning\hello_world\{archetypes,assets,content,data,layouts,static}
Watching for config changes in D:\Projects\DesistDaydream\hugo-learning\hello_world\config.toml
Environment: "development"
Serving pages from memory
Running in Fast Render Mode. For full rebuilds on change: hugo server --disableFastRender
Web Server is available at http://localhost:1313/ (bind address 127.0.0.1)
Press Ctrl+C to stop
```

我们可以通过浏览器，访问默认的 1313 端口浏览我们的网站，但是此时我们只能看到一个 Hugo 默认的 `Page Not Found`，因为我们还没有为网站设置、添加任何内容。

Hugo 从 `content/` 目录中渲染内容到页面，我们使用 `hugo new posts/my-first-post.md` 命令将会创建 `content/posts/my-first-post.md` 文件，我们可以自行在该文件中添加 markdown 格式的内容。

但是我们依然无法看到任何东西，因为 Hugo 提供了非常大的自由度，并不会限制 HTML 的样式，所以我们需要先自己创建一个 HTML 页面(就像写前端一样)。

在 layouts/ 目录下新建 \_default 目录，并创建一个名为 single.html 文件，写下如下内容：

```html
<!DOCTYPE html>
<html lang="zh">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <meta http-equiv="X-UA-Compatible" content="ie=edge">
    <title>{{.Title}}</title>
</head>
<body>
    <div id="post" class="post">
        <article>
            <header>
                <h1 class="post-title">{{ .Title }}</h1>
            </header>
            {{.Content}}
        </article>
    </div>
</body>
</html>
```

此时我们直接访问 `http://localhost:1313/posts/my-first-post/` 即可看到我们刚才添加的 markdown 的内容。只不过没有任何样式，光秃秃的~

![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/hugo/20230122175219.png)


## 使用主题

我们自己编写 HTML 是非常复杂的，咱是要内容管理。。又不是写前端页面\~\~~o(╯□╰)o

Hugo 贴心得提供了主题功能，可以让我们专注于内容的产出，在 Hugo 官方的主题页面中，我们可以挑选我们喜欢的主题并放在 themes/ 目录下，以便使用时供 Hugo 加载

> 除了将主题放在 themes/ 目录下，还可以使用 Hugo 模块功能，将主题当做 Go 模块一样的东西，统一管理。这样在我们创建多个 Hugo 站点并使用同一个主题时，不用重复下载了。

我们使用官方示例中的基本主题：

```bash
git submodule add https://github.com/theNewDynamic/gohugo-theme-ananke themes/ananke
echo "theme = 'ananke'" >> config.toml
```

主题将会被下载到 thems/ 目录中，并且我们在 config.toml 文件中指定要使用的主题名称。

此时再打开 1313 端口，我们就可以看到我们的站点了，第一篇文章以标题和概要的形式被展现在首页中。

![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/hugo/20230122164828.png)

### 通过 Hugo 模块使用主题

想要使用 Hugo 模块，我们需要 Go 环境

执行如下命令

```bash
hugo mod init github.com/DesistDaydream/hugo-learning/hello_world
hugo mod get github.com/theNewDynamic/gohugo-theme-ananke
```

此时主题将会被下载到 `%TMP%/hugo_cache/modules/filecache/modules/pkg/mod/github.com` 目录下，然后我们可以删掉项目目录中 themes/ 目录下的主题文件了~o(∩_∩)o

修改 config.toml 文件
```toml
theme = ["github.com/theNewDynamic/gohugo-theme-ananke"]
```


## 发布我们的网站

在基本示例中，我们只是在本地调试，如果想要将网站运行在服务器上，那么肯定需要像前端代码一样，将这些文件打包才可以。

Hugo 打包非常简单，执行 `hugo` 命令即可在 public/ 目录中生成我们网站的静态页面，将这个目录下的所有文件，统统放到 Nginx 中响应页面的目录，就可以访问我们自己的网站了~

### 最佳实践

很多时候，我们通过工作流（GitHub Action 等）将 `public/` 目录下的文件转存到新项目中，并不会将原始内容与打包好的前端代码放在一起。

## 目录结构

> 参考：
> - [官方文档，入门-目录结构](https://gohugo.io/getting-started/directory-structure/)

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

以 `.html` 文件的形式存储模板，这些文件指定如何将您的内容视图呈现到静态网站中。模板包括 [list pages](https://gohugo.io/templates/list/), your [homepage](https://gohugo.io/templates/homepage/), [taxonomy templates](https://gohugo.io/templates/taxonomy-templates/), [partials](https://gohugo.io/templates/partials/), [single page templates](https://gohugo.io/templates/single-page-templates/),等

如果我们不使用主题，则 Hugo 会从 `layouts` 目录中读取前端代码并渲染页面。

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

Hugo 模块是一个类似 Go 模块一样的存在。模块可以是我们的主项目或其他较小的模块，提供 Hugo 的 7 种组件类型中的一种或多种：

- **static**
- **content**
- **layouts**
- **data**
- **assets**
- **i18n**
- **archetypes**

在 config.toml 文件中的 module 字段添加配置，即可为站点设置引用的模块，我们可以将主题当做一个模块。

**注意：Hugo 模块与 Go 模块一样，也需要一个代理服务器，我们只需要在 module 部分配置 proxy 指令，值与 go proxy 一样即可**

# 安装 Hugo

安装 `hogo` 命令行工具，即可开始使用 Hugo。推荐下载扩展版 hugo，即名字带有 extended 的文件。

从 [release](https://github.com/gohugoio/hugo/releases) 页面下载带 **extended** 后缀的文件。

# 关联文件与配置

**hugo.toml | hugo.yaml | hugo.json** # 站点的配置文件，通常在站点的根目录。在 0.110.0 版本之前，默认的文件名是 config.toml 之类的。

**${Site_Root_dir}/config** # 可以将站点根目录下的 config.toml | config.yaml | config.json 拆分后保存到该目录。

Hugo 运行时所需的缓存目录。包括需要使用的模块等：
- Windows:
	- **%TMP%/hugo_cache/**
- Linux:
	- **${TMP}/hugo_cache/**

# Hugo 与 Obsidian


## URL 与 markdown 链接问题

> 参考：
> - https://cloud.tencent.com/developer/article/1688894

Obsidian 内部链接是这种格式 `[B cd](/A/b/B%20cd.md)`

Hugo 生成的内容资源的 URL 是 https://demo.org/a/b/b-cd

此时，如果我们从页面点击 B cd，将会跳转到 https://demo.org/A/b/B-cd 页面，此时将会看到 404。。。。

解决方式：

在 hugo.config 中添加 `disablePathToLower = true` 配置，以关闭转换为小写的功能。

在 layouts/404.html 中添加如下脚本：

```js
<script>
  var currenturl = location.href.replace(/%20/g, "-").replace(".md", "");
  if (currenturl != location.href) {
    location.href = currenturl;
  }
</script>
```

此时跳转到 404 时，将会去掉 `.md` 后缀，以及将所有的 `%20` 替换成 `-`

