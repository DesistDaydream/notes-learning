---
title: GitHub
linkTitle: GitHub
weight: 1
---

# 概述

> 参考：
>
> - [官方文档](https://docs.github.com/)
> - [官方文档，中文](https://docs.github.com/cn)

在代码仓库中，点击 `.` 即可进入 Web 版的 VS Code，在线编辑当前仓库的代码。

# GitHub Desktop

https://desktop.github.com/

专注于重要的事情，而不是与 Git 对抗。无论您是 Git 新手还是经验丰富的用户，GitHub Desktop 都能简化您的开发工作流程。

# GitHub 加速

若没有代理的场景下想要访问 GitHub，可以通过 GitHub 镜像站点，常见可用的：

https://ghproxy.link/

- https://ghfast.top/


这种加速站的使用方式通常都一样，只需要在地址前加上加速站点的地址即可，e.g. 要克隆 ggml-org/llama.cpp 项目，使用 https://ghfast.top/ 加速的场景下，执行如下命令：

```bash
git clone https://ghfast.top/https://github.com/ggml-org/llama.cpp.git
```