---
title: PKM
linkTitle: PKM
date: 2023-02-28T20:27:00
weight: 1
---

# 概述

> 参考：
>
> - [Wiki, Personal knowledge management](https://en.wikipedia.org/wiki/Personal_knowledge_management)

**Personal knowledge management(个人知识管理，简称 PKM)** 是一个收集信息的过程

# 个人知识记录规范

## 配置文件与命令行参数

各种类型的配置文件（包括 INI、JSON、等等）、命令行参数是类似 Key/Value 的结构。在我的笔记中，记录格式一般都是这样的：

`键的名称(值的类型)`

`--命令行参数名称(参数值类型)`

对于配置文件来说，有的值的类型比较复杂（比如是一个 OBJECT 类型），可以再创建一个自定义的名称以在单独的章节下记录。

- **node1**(STRING)
- **node2**([sub_node1](#sub_node1))
- **node3**(\[][configs](#configs))
- **node4**(OBJECT)
- **node5**(map\[STRING]STRING)
- **node6**(map\[STRING][sub_node2](#sub_node2))

笔记中的配置内容通常符合如下我自己定义的规范：

- 加粗的是 Key
- 括号中是 Value 的类型，Value 一般是非 Scalar 类型的节点。
  - 若 Value 的类型是 Object，那么一般类型名称是自定义的。
    - 由于 Object 类型的节点中，Value 也可以是一个节点，那么 **Value 就有可能是由一个或多个内容**组成，为了可以方便得复用这些内容，所以给**它们起了一个名字**。这就**好像编程中的使用函数**一样。
    - 若 OBJECT 类型的字段下的字段非常多，我会在单独的标题中记录，[Pod Manifest](/docs/10.云原生/Kubernetes/API%20Resource%20与%20Object/API%20参考/工作负载资源/Pod%20Manifest.md) 是典型的例子。不但在单独的标题记录，而且还为这些字段进行了分组。在我们理解时，只有带有 `(XXX)` 这种写法的，才是 YAML 中真正的字段，而标题，通常不作为真正的字段，只是作为该字段的一个指示物，用以记录该字段下还有哪些字段。
    - 若 Object 类型的字段比较简单，没有复杂的子字段，那么笔记中就直接用 `OBJCET` 这几个字符表示。
  - 若 Value 的类型是 STRING、INT、etc. 简单类型，但是其含义很复杂，也会将该字段值的类型写作连接，在独立章节记录。

这种规范为了文档的整洁性，让相同层级的字段在一起，可以一眼看到同级内容，让 Value 与 Key 分开，将 Value 所包含的具体内容放在单独链接（i.e. 单独章节）中。

不管是 老式的 [INI](/docs/2.编程/无法分类的语言/INI.md)、还是新一些的 [JSON](/docs/2.编程/无法分类的语言/JSON.md)、[YAML](/docs/2.编程/无法分类的语言/YAML.md)、[TOML](/docs/2.编程/无法分类的语言/TOML.md)、etc. 都可以使用这套理论来记录

## 命名与命名使用的符号

### `-` 与 `_`

对于文件名的命名来说

| 符号  | 用途                 | 中文  | 英文                   |
| --- | ------------------ | --- | -------------------- |
| `-` | 将连接两端的单词当作**两个单词** | 短横线 | hyphen               |
| `_` | 将连接两端的单词当作**一个单词** | 下划线 | underscore/underline |

e.g. 在 [这篇文章](https://adoyle.me/Today-I-Learned/others/file-naming-with-underscores-and-dashes.html) 里有提到说，Google 搜索引擎会将 `_` 连接的单词作为一个单词。比如搜索 `web_site` 实际上只会找关键词 `website`。只有 `web-site` 会分为 `web` 和 `site` 来查找。

- This_is_a_single_word
- This-is-a-sentence-with-multiple-words

> Tips: 很多时候，在我们进行编辑时（不管是 ide 写代码还是编辑文件名），利用 “按住 ctrl + 左右方向键” 功能快速跳过单词时，除了会直接跳到下一个空格外，还会直接跳到 `-`，而 `_` 则被当作一整个单词跳过（不过跳到 `_`）。

# 知识管理工具

[Notion](https://www.notion.so/)

[AppFlowy](https://github.com/AppFlowy-IO/AppFlowy) # 开源版 Notion

- 现阶段(0.0.4)只是一个本地应用程序，无法通过浏览器使用

[语雀](https://www.yuque.com/)

- 如何看待语雀付费策略？ https://www.zhihu.com/question/562238887
- 文档导出: https://github.com/yuque/yuque-exporter
  - 先执行 crawl 生成想要下载的文档源数据
  - 执行 build 根据已存在的源数据生成 markdown 文件
    - 源码执行：
      - `pnpm i`
      - `pnpm build`
      - `export YUQUE_TOKEN="XXXX"`
      - `node ./dist/bin/cli.js crawl desistdaydream/ycpve3`
      - `node ./dist/bin/cli.js build`
- 语雀文档导出: https://github.com/yuque/yuque-exporter

[飞书](https://www.feishu.cn/product/docs)

- 飞书转 Markdown: https://github.com/Wsine/feishu2md

[Dendron](https://github.com/dendronhq/dendron) # 开源的、本地优先的、基于 MarkDown 的笔记工具

- [公众号-Github 爱好者，专为开发人员构建的个人知识管理工具-Dendron](https://mp.weixin.qq.com/s/HbM93O49aOgW6w_ZX9lzlA)

[Obsidian](/docs/学习/PKM/Obsidian.md) # Markdown 渲染程序

# 社区

PKMer

- 是一个【知识管理】爱好组织，我们热衷于知识管理，喜欢讨论提升效率软件，以及那些让你觉得欣喜的技术。
- https://pkmer.cn/
- 包含 Markdown、Obsidian、Excaildraw、Zotero、TiddyWiki、etc. 相关专题

