---
title: Browser
linkTitle: Browser
date: 2023-11-02T01:23
weight: 1
---

# 概述

> 参考：
>
> - [Wiki, Web_browser](https://en.wikipedia.org/wiki/Web_browser)

在访问一个网页时，除非收到 3XX 重定向的响应，否则浏览器地址栏中的地址是不会改变的。比如 Nginx 中的 rewrite 功能，如果不使用 **redirect** 或 **permanent** 标志，那么所有的 URL 改变都是针对 Nginx 内部来说的。

# 解决网页播放【鼠标移开屏幕或点击屏外视频暂停播放】

原文： https://www.jianshu.com/p/945851ea95da

从网页的 F12 中，元素-事件监听器 中

- 将【blur】所有内容【remove】掉
- 单击【mouseout】左边的倒三角，将出现的子元素全部remove掉，
- 将【mouseup】也用同样的操作移除掉子元素，现在就可以成功切换页面而不受限制啦!

（点击Remove要精准，remove会把blur清除，不会进入其他设置）

注：可同时点开多个网页播放器并行播放不暂停，提高效率


