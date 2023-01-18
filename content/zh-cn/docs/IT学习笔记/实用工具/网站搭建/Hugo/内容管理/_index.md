---
title: "内容管理"
linkTitle: "内容管理"
weight: 20
---

# 概述
> 参考：
> - [官方文档，内容管理](https://gohugo.io/content-management/)



# FrontMatter(前页) 配置
在每篇文档的开头，我们可以使用 TOML、YAML、JSON、ORG 这四种格式配置文档的元数据，通过这些元数据，我们可以让 Hugo 或 Hugo 主题提供更强大的功能。

各种格式使用不同的符号来识别：

```
+++
TOML 格式
+++


---
YAML 格式
---

{JSON 格式}
```

## 前页变量
在 FrontMatter 配置的元数据可以生成 FrontMatter Variables(前页变量)，Hugo 具有如下预定义的前页变量

**title** # 文档的标题

**linktitle** # 用于创建内容链接；如果设置，Hugo 将会使用 linkTitle 的值作为菜单栏、导航栏中文章的链接名称。 Hugo 还可以通过链接标题对内容列表进行排序。

**data** # 分配给此页面的日期时间。这通常是从 front matter 的日期字段中获取的，但这种行为是可配置的。

**lastmod** # 上次修改内容的日期


