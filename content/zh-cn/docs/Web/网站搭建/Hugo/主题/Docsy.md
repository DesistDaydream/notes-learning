---
title: "Docsy"
linkTitle: "Docsy"
weight: 20
---

# 概述

> 参考：
> 
> - [GitHub 项目，google/docsy](https://github.com/google/docsy)
> - [官网](https://www.docsy.dev/)

Kubernetes 的官网就是 Docsy 主题。

注意：Docsy 必须使用 扩展版 hugo，即 hugo_extended。

# 预览和部署 Docsy 主题网站

> 参考：
> 
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
cat >> hugo.toml <<EOL
[module]
proxy = "direct"
[[module.imports]]
path = "github.com/google/docsy"
EOL
```

## 预览

```bash
hugo server
```

## 部署

https://lucumt.info/post/hugo/using-github-action-to-auto-build-deploy/

https://tomial.github.io/posts/hugo%E4%BD%BF%E7%94%A8github-action%E8%87%AA%E5%8A%A8%E9%83%A8%E7%BD%B2%E5%8D%9A%E5%AE%A2%E5%88%B0github-pages/

使用 Hugo 搭建 GitHub Pages https://zz2summer.github.io/github-pages-hugo-%E6%90%AD%E5%BB%BA%E4%B8%AA%E4%BA%BA%E5%8D%9A%E5%AE%A2/#%E4%B8%83%E6%97%A5%E5%B8%B8%E6%93%8D%E4%BD%9C

## 常见问题

[构建站点时出错：# POSTCSS: failed to transform "scss/main.css"](https://github.com/google/docsy/issues/235)

# Docsy 配置与关联文件

Docsy 也会使用 Hugo 的 CONFIG 文件来配置站点。参考[官网，内容和定制](https://www.docsy.dev/docs/adding-content/)章节来修改 CONFIG 文件，以改变主题样式。

Docsy 的配置主要在 Hugo CONFIG 配置文件中的 `[params]` 部分

## 多语言支持
 
https://www.docsy.dev/docs/language/#content-and-configuration

Docsy 主题中的页面有些地方是没有翻译的，可以在 Hugo 项目根目录创建 `/i18n/` 目录，并按照语言目录名称对应配置文件（比如 `zh-cn/` 目录下的文章就会找 `zh-cn.toml` 文件），即可为这些地方显示对应语言的内容。

> 配置文件内容可以参考 Kubernetes 的 https://github.com/kubernetes/website/blob/main/data/i18n/zh-cn/zh-cn.toml


