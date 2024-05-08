---
title: PKM
linkTitle: PKM
date: 2024-02-28T23:27
weight: 1
---

# 概述

> 参考：
>
> - [Wiki，Personal knowledge management](https://en.wikipedia.org/wiki/Personal_knowledge_management)

**Personal knowledge management(个人知识管理，简称 PKM)** 是一个收集信息的过程

# 个人知识记录规律

## 配置文件与命令行参数的记录方式

与 [YAML](/docs/2.编程/无法分类的语言/YAML.md) 文章中的“各种产品官方文档中对 YAML 格式的配置文件的描述”段落类似，各种类型的配置文件（包括 INI、JSON、等等）、命令行参数是类似 Key/Value 的结构。

在我得笔记中，记录格式一般都是这样的：

`键的名称(值的类型)`

`--命令行参数名称(参数值类型)`

对于配置文件来说，有的值的类型比较复杂（比如是一个 object 类型），可以再创建一个自定义的名称以在单独的章节下记录。

# 知识管理工具

[Notion](https://www.notion.so/)

[AppFlowy](https://github.com/AppFlowy-IO/AppFlowy) # 开源版 Notion

- 现阶段(0.0.4)只是一个本地应用程序，无法通过浏览器使用

[语雀](https://www.yuque.com/)

- 如何看待语雀付费策略？<https://www.zhihu.com/question/562238887>
- 文档导出：<https://github.com/yuque/yuque-exporter>
  - 先执行 crawl 生成想要下载的文档源数据
  - 执行 build 根据已存在的源数据生成 markdown 文件
    - 源码执行：
      - `pnpm i`
      - `pnpm build`
      - `export YUQUE_TOKEN="XXXX"`
      - `node ./dist/bin/cli.js crawl desistdaydream/ycpve3`
      - `node ./dist/bin/cli.js build`
- 语雀文档导出：<https://github.com/yuque/yuque-exporter>

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