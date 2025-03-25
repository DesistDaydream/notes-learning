---
title: Obsidian
linkTitle: Obsidian
weight: 20
---

# 概述

> 参考：
>
> - [官网](https://obsidian.md/)
> - [官方文档](https://help.obsidian.md/)
> - [开发者文档](https://docs.obsidian.md/)

Obsidian 也是基于 Chromium 的，使用 `Ctrl + Shift + i` 快捷键可以打开 [DevTools](/docs/Web/Browser/DevTools.md)

[英文论坛](https://forum.obsidian.md/)

[中文论坛](https://forum-zh.obsidian.md/)

中文论坛与英文论坛的账户不共享

# Obsidian 关联文件与配置

> 参考：
>
> - [官方文档，Obsidian 如何存储数据](https://publish.obsidian.md/help-zh/%E9%AB%98%E7%BA%A7%E7%94%A8%E6%B3%95/Obsidian+%E5%A6%82%E4%BD%95%E5%AD%98%E5%82%A8%E6%95%B0%E6%8D%AE)

Obsidian 本身的运行时数据保存路径（我们假定设为 `${ObsidianData}`）[^官方文档]

- **~/.config/Obsidian/** # Linux 系统
- **%APPDATA%/Obsidian/** # Windows 系统

**${REPO}/.obsidian/** # 特定于每个仓的配置的保存路径。在每个仓库的根目录下的 .obsidian/ 目录中。

- **workspaces.json** # 工作区布局的配置文件。通常在每个仓库各自 .obsidian/ 目录下。

[https://www.bilibili.com/video/BV1Dy4y1375P](https://www.bilibili.com/video/BV1Dy4y1375P)

# Vault

Obsidian 将本地仓库称为 **vault**，一个 vault 就是文件系统上的一个文件夹。这个 vault 中保存了所有记录的 文档、附件、插件、etc. 。

# 编辑与格式

https://help.obsidian.md/Editing+and+formatting/Basic+formatting+syntax

Obsidian 识别 [Markdown](/docs/2.编程/标记语言/Markdown.md) 语法并渲染成文章。

## Callouts

https://help.obsidian.md/Editing+and+formatting/Callouts

Obsidian 扩展了 Markdown 的 **Callouts(标注)** 效果。

> [!note]
> Lorem ipsum dolor sit amet

> [!tip]
> Lorem ipsum dolor sit amet

> [!success]
> Lorem ipsum dolor sit amet

> [!warning]
> Lorem ipsum dolor sit amet. 还可以用 `caution`, `attention` 这俩关键字

> [!bug]
> Lorem ipsum dolor sit amet

> [!example]
> Lorem ipsum dolor sit amet

还有很多样式可以参考官方文档

## Properties vs Tags

https://help.obsidian.md/Editing+and+formatting/Tags

https://help.obsidian.md/Editing+and+formatting/Properties

https://forum.obsidian.md/t/the-remaining-advantages-of-tags-over-properties-in-obsidian/69436?page=2

先有的 Tag 后有的 Properties

老式 Tag 需要在文章内部使用 `#STRING` 来标识

Property 可以在开头的 `---` 包裹的文章元数据中使用 tags 关键字添加 Tag。

个人感觉只有在一个知识点被用在多个大类的情况下，才需要添加标签，否则该知识点通过目录即可进行分类，比如一个程序又属于网络工具，又属于安全工具，那么若是放在安全目录下则可以添加网络标签，放在网络目录下则可以添加安全标签。

或者某个知识点的文章，不太好分类，甚至具有嵌套类型等，都可以通过标签解决。

## Footnote(脚注)

**第一种方式，声明与使用分开**

在这里声明一个脚注，声明的脚注必须独立占用一行

[^1]: 参考文献1

在这里使用脚注 [^1]。

**第二种方式，声明与使用合并**

这是一个脚注 ^[参考文献1]。

# 第三方插件

清理没有引用的图片

-

## 自定义排序目录

众人的需求: https://forum.obsidian.md/t/file-explorer-custom-sort/1602

解决方案：

https://forum.obsidian.md/t/file-explorer-custom-sort/1602/212

- 主要是想用 Bartender 插件解决，而安装 Bartender 插件需要使用 Brat 插件（也可以直接下载文件放到 plugin/ 目录中）
- Install BRAT from community plugin panel and then install Bartender beta plugin using BRAT
  - [https://github.com/nothingislost/obsidian-bartender](https://github.com/nothingislost/obsidian-bartender)
  - [https://github.com/TfTHacker/obsidian42-brat](https://github.com/TfTHacker/obsidian42-brat)
- then, step 1-2-3
  - ![image](https://forum.obsidian.md/uploads/default/original/3X/9/1/9150dde8b90e4a93b6edc58cd4cc51c9f4f61abb.png)
- do not forget step 3
- ![image](https://forum.obsidian.md/uploads/default/original/3X/2/d/2d251736195adb913c336f7d309be7ab7c4f25ef.png)

then you can drag as you like

nothingislost/obsidian-bartender 不维护之后其他开发者的 fork

- https://github.com/Mara-Li/obsidian-bartender

## Templater

- https://github.com/SilentVoid13/Templater

比官方自带的模板功能更强大

## Excalidraw

让 Obsidian 可以渲染 Excalidraw 图，并且可以在 Obsidian 中编辑 Excalidraw 图。

## Image auto upload Plugin

图片自动上传到图床，依赖系统中启动的 PicGo 或 PicList

## Image Toolkit

可以在 Obsidian 中打开图片

# 最佳实践

[用 Git 在 Android 和 Windows 间同步 Obsidian 数据库](https://sspai.com/post/68989)

- Android 上使用了 MGit 工具

# 常见问题

## Win11 打开多个仓库的任务栏图标无法合并

https://forum.obsidian.md/t/opening-multiple-vaults-creates-multiple-taskbar-icons-is-this-intended-windows-11/55346/3

- 取消所有任务栏固定
- 打开一个仓库，假如为 A
- 再打开另一个仓库，假如为 B
- 把 B 固定到任务栏
- 关闭所有 obsidian 窗口，再从任务栏打开时，就会发现所有 Obsidian 的仓库都合并到一起了。
