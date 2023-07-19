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

Hugo 创建站点时主要依赖两部分，**Content(内容)** 与 **Layout(布局)**

- **Content(内容)** 表示数据。存在 `content/` 目录下。
  - 该目录下的每个文件都会抽象为一个 **Page(页面)**。其实我们在浏览到的页面就是 Hugo 中的 Page 的概念，而 `content/` 目录就是存放这些 Page 的地方。内容的管理详见[内容管理](docs/Web/网站搭建/Hugo/内容管理/内容管理.md)章节
- **Layout(布局)** 表示页面。存在 `layouts/` 目录下。
  - 该目录下的每个文件都会抽象为一个 **Template(模板)**

通过多种渠道获取到数据(i.e. Content)后，需要在页面(i.e. Layout)中填充数据，这就是模板渲染的过程，渲染完成后，可供浏览的页面称之为 **View(视图)**。

Hugo 的这种渲染行为与 Go 的模板渲染机制一致，并提供了更丰富的功能。

# Hugo 的基本使用

> 参考：
>
> - [官方文档，入门-快速开始](https://gohugo.io/getting-started/quick-start/)

这里的示例并没有安全按照官方文档走，而是在我学习之后改编的，官方文档的示例其实会让新手对于渲染逻辑和顺序产生迷惑。

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
~]# hugo server
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

我们可以通过浏览器，访问默认的 1313 端口浏览我们的网站，但是此时我们只能看到一个 Hugo 默认的 `Page Not Found`，因为我们还没有为网站添加任何内容。

Hugo 默认将 `content/` 目录的文件作为内容数据渲染与布局结构一起渲染成页面，所以需要在 layouts/ 目录和 content/ 目录下创建文件。

> 若只在 content/ 目录下创建文件，后台将会有警告提示：`WARN 2023/05/25 14:45:44 found no layout file for "HTML" for kind "page": You should create a template file which matches Hugo Layouts Lookup Rules for this combination.`。并且页面也没有任何东西，这是因为没有布局将这些内容渲染成页面展示出来。

首先在 `layouts/` 目录下新建 `index.html` 文件，并写入如下内容：

```html
Hello Hugo!

{{ .Content }}
```

此时刷新页面，将会看到如下页面

![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/hugo/20230529112224.png)

其中 `{{ .Content }}` 这部分并没有被渲染内容是因为在 content/ 目录还是空的。

我们可以使用命令 `hugo new _index.md` 创建内容文件，也可以手动在 `content/` 目录下创建 `_index.md` 文件。推荐自己手动创建，

> 使用命令创建的优势有一个前提，是前提是在 `archetypes/` 目录下放一些原型，这并不属于基本使用的内容，以后再介绍；并且自动创建的文件带有 `draft: true` 指令。Hugo 默认不会构建被标记为 [draft(草稿)、future(未来)、expired(过期)](https://gohugo.io/getting-started/usage/#draft-future-and-expired-content) 的内容，我们还要手动删掉，比较麻烦。

我们创建完文件后，写入如下内容：

```md
---
title: "主页"
date: 2023-05-29T11:14:20+08:00
---

这些字符是在 content/_index.md 文件中的。

## 概述

这是 **bold(粗体)** 文本text, 这是 *斜体* 文本.

访问 [Hugo](https://gohugo.io) 网站!
```

刷新页面后，将会看到如下内容：

![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/hugo/20230529113139.png)

可以看到，在布局文件中的 `{{ .Content }}` 部分被替换成内容了，这就是 Hugo 的模板渲染能力。

## 使用主题

我们自己编写 HTML 是非常复杂的，咱是要内容管理。。又不是写前端页面 o(╯□╰)o

Hugo 贴心得提供了主题功能，可以让我们专注于内容的产出，说白了，就是不用再关心 `layouts/` 目录，而是只在 `content/` 目录下创建我们的内容就可以了。

注意：由于前文我们创建了 `layouts/index.html` 文件，这回覆盖主题的注意布局，所以需要删掉该文件，再进行后面的测试。

在 Hugo 官方的主题页面中，我们可以挑选我们喜欢的主题并放在 `themes/` 目录下，以便使用时供 Hugo 加载

> 除了将主题放在 `themes/` 目录下，还可以使用 Hugo 模块功能，将主题当做 Go 模块一样的东西，统一管理。这样在我们创建多个 Hugo 站点并使用同一个主题时，不用重复下载了。

我们使用官方示例中的基本主题。

```bash
git submodule add https://github.com/theNewDynamic/gohugo-theme-ananke themes/ananke
echo "theme = 'ananke'" >> hugo.toml
```

主题将会被下载到 `thems/` 目录中，并且我们在 hugo.toml 文件中指定要使用的主题名称。

此时再打开 1313 端口，我们就可以看到我们的站点带上了样式

![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/hugo/20230529113945.png)

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

使用 Hugo 的模块时，也要像 Go 一样配置 Proxy，在 hugo.toml 文件中添加如下内容：：

```toml
[module]
proxy = 'https://goproxy.cn,https://goproxy.io,direct'
```

## 发布我们的网站

在基本示例中，我们只是在本地调试，如果想要将网站运行在服务器上，那么肯定需要像前端代码一样，将这些文件打包才可以。

Hugo 打包非常简单，执行 `hugo` 命令即可在 `public/` 目录中生成我们网站的静态页面，将这个目录下的所有文件，统统放到 Nginx 中响应页面的目录，就可以访问我们自己的网站了~

### 最佳实践

很多时候，我们通过工作流（GitHub Action 等）将 `public/` 目录下的文件转存到新项目中，并不会将原始内容与打包好的前端代码放在一起。

## 目录结构

> 参考：
>
> - [官方文档，入门-目录结构](https://gohugo.io/getting-started/directory-structure/)

### archetypes

**archetypes** 译为**原型**，是创建新 Content 时使用的模板。我们在使用 `hugo new` 命令创建新的 Content 时，会使用该目录下的 default.md 文件作为原型创建新的文件。

比如快速开始中，我们创建了一个名为 my-first-post.md 的文件，其内容为：

```md
---
title: "My First Post"
date: 2023-05-25T01:40:23+08:00
draft: true
---
```

这是因为使用了 archetypes/default.md 文件作为原型：

```md
---
title: "{{ replace .Name "-" " " | title }}"
date: {{ .Date }}
draft: true
---
```

[官方文档，内容管理-原型](https://gohugo.io/content-management/archetypes/)

### assets

存放所有需要由 Hugo Pipes 处理的文件。

### config

> 注意：config 目录并不会随着 `hugo new site example` 命令自动创建出来。而是在我们需要将单一配置文件拆分时，手动创建并使用的。

Hugo 附带了大量的配置指令。config 目录是将这些指令存储为 JSON、YAML 或 TOML 文件的地方。每个根设置对象都可以作为自己的文件并按环境构建。设置最少且不需要环境的项目可以在其根目录下使用单个 config.toml 文件。

许多站点可能几乎不需要配置，但 Hugo 附带了大量 [configuration directives(配置指令)](https://gohugo.io/getting-started/configuration/#all-configuration-settings)，用于更详细地指导我们希望 Hugo 如何构建网站。注意：默认情况下不创建 config 目录。

### content/

我们使用 Hugo 创建的网站的所有内容通常都要放在 content/ 目录中。content/ 目录下的每个顶级文件夹称为 [content section(内容部分)](https://gohugo.io/content-management/sections/)。

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

以 `.html` 文件的形式存储模板，这些文件指定如何将您的内容视图呈现到静态网站中。模板包括 [list pages](https://gohugo.io/templates/list/), your [homepage](https://gohugo.io/templates/homepage/), [taxonomy templates](https://gohugo.io/templates/taxonomy-templates/), [partials](https://gohugo.io/templates/partials/), [single page templates](https://gohugo.io/templates/single-page-templates/) 等等，不同名称的模板对应渲染不同的页面。可以在[这里](https://gohugo.io/templates/lookup-order/#hugo-layouts-lookup-rules-with-theme)找到 Hugo 在渲染不同页面时需要使用 layouts/ 目录下的哪些文件。

如果我们不使用主题，我们则需要在 `layouts/` 目录中自己编写 HTML 文件以供 Hugo 渲染前端页面。

### public/

使用 `hugo` 命令生成网站的静态文件后，将会保存到 public 目录。public 目录生成的静态文件，可以直接通过 web 服务访问到。

### static/

存储所有静态内容：图像、CSS、JavaScript 等。当 Hugo 构建您的站点时，静态目录中的所有资产都会按原样复制。使用静态文件夹的一个很好的例子是在 Google Search Console 上验证网站所有权，您希望 Hugo 在不修改其内容的情况下复制完整的 HTML 文件。

### themes/

> 更推荐的是使用 Hugo 模块使用主题，该目录不推荐使用。

Hugo 主题可以安装到该目录，使用 `hugo server --themes` 指定使用的主题时，将会从该目录出寻找。

### hugo.toml

Hugo 运行站点时所使用的配置文件。

推荐使用 `config/` 目录，以便拆分 hugo.toml 文件。可以将 hugo.tom 文件移动到在 `config/_default/hugo.toml` 处作为默认配置。

# Hugo Modules(模块)

Hugo 模块是一个类似 Go 模块一样的存在。模块可以是我们的主项目或其他较小的模块，提供 Hugo 的 7 种组件类型中的一种或多种：

- **static**
- **content**
- **layouts**
- **data**
- **assets**
- **i18n**
- **archetypes**

在 hugo.toml 文件中的 module 字段添加配置，即可为站点设置引用的模块，我们可以将主题当做一个模块。

**注意：Hugo 模块与 Go 模块一样，也需要一个代理服务器，我们只需要在 module 部分配置 proxy 指令，值与 go proxy 一样即可**

详见 [模块](/docs/Web/网站搭建/Hugo/模块/模块.md)

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
>
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

