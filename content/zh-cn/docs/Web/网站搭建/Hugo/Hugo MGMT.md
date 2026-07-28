---
title: "Hugo MGMT"
created: "2026-07-28T16:08"
weight: 100
---

# 概述

> 参考：
>
> -

# hugo server 热更新在 WSL 中

执行完 hugo server 进行调试时，若在 WSL 中，且路径是 Windows 中的路径，那么 WSL2 访问 `/mnt/d/` 走的是 9p/drvfs 协议，Linux 的 inotify 文件监听机制在这种跨文件系统的场景下经常收不到事件，尤其是用 Windows 侧的编辑器（如 Windows 版 VS Code、Notepad 等）保存文件时。

若想在 WSL 中监听变化，可以添加轮询参数 `hugo server --poll 700ms` 来启动 hugo 内部 WEB 服务以调试。

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
