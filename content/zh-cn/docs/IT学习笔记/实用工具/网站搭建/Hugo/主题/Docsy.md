---
title: "Docsy"
linkTitle: "Docsy"
weight: 20
---

# 概述

> 参考：
> - [GitHub 项目，google/docsy](https://github.com/google/docsy)
> - [官网](https://www.docsy.dev/)

Kubernetes 的官网就是 Docsy 主题。

注意：Docsy 必须使用 扩展版 hugo，即 hugo_extended。

# 预览和部署 Docsy 主题网站

> 参考：
> - [官方文档，预览和部署](https://www.docsy.dev/docs/deployment/)

## 准备环境

```bash
npm install -D autoprefixer
npm install -D postcss-cli
npm install -D postcss
```

> 若不在本地安装，则使用 `hugo` 命令构建静态文件时将会报错

## 生成站点文件

### 生成模板站点

```bash
export MY_SITE_DIR="docsy"
git clone https://github.com/google/docsy-example.git ${MY_SITE_DIR}
cd  ${MY_SITE_DIR}
hugo server
```

### 生成空白站点

```Bash
hugo new site .
hugo mod init github.com/me/my-new-site
hugo mod get github.com/google/docsy@v0.6.0
cat >> config.toml <<EOL
[module]
proxy = "direct"
[[module.imports]]
path = "github.com/google/docsy"
[[module.imports]]
path = "github.com/google/docsy/dependencies"
EOL
```

## 预览

```bash
hugo server
```

## 部署

https://lucumt.info/post/hugo/using-github-action-to-auto-build-deploy/

https://tomial.github.io/posts/hugo%E4%BD%BF%E7%94%A8github-action%E8%87%AA%E5%8A%A8%E9%83%A8%E7%BD%B2%E5%8D%9A%E5%AE%A2%E5%88%B0github-pages/

https://www.bloghome.com.cn/post/git-zi-mo-kuai-yi-ge-cang-ku-bao-han-ling-yi-ge-cang-ku.html

使用 Hugo 搭建 GitHub Pages https://zz2summer.github.io/github-pages-hugo-%E6%90%AD%E5%BB%BA%E4%B8%AA%E4%BA%BA%E5%8D%9A%E5%AE%A2/#%E4%B8%83%E6%97%A5%E5%B8%B8%E6%93%8D%E4%BD%9C

## 常见问题

[构建站点时出错：# POSTCSS: failed to transform "scss/main.css"](https://github.com/google/docsy/issues/235)

# Docsy 配置与关联文件

Docsy 也会使用 Hugo 的 CONFIG 文件来配置站点。参考[官网，内容和定制](https://www.docsy.dev/docs/adding-content/)章节来修改 CONFIG 文件，以改变主题样式。


