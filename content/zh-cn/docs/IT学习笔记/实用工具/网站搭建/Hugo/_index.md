# 概述

> 参考：
> - [GitHub 项目，gohugoio/hugo](https://github.com/gohugoio/hugo)

使用 Hugo 搭建 GitHub Pages
- https://zz2summer.github.io/github-pages-hugo-%E6%90%AD%E5%BB%BA%E4%B8%AA%E4%BA%BA%E5%8D%9A%E5%AE%A2/#%E4%B8%83%E6%97%A5%E5%B8%B8%E6%93%8D%E4%BD%9C
- 

喜欢的主题：
- https://themes.gohugo.io/themes/hugo-theme-techdoc/

# 安装 Hugo


# 关联文件与配置
**${TMP}/hugo_cache/*** # 运行时的缓存。包括模块等


# Docsy 主题

## 环境准备
```bash
npm install -D autoprefixer
npm install -D postcss-cli
npm install -D postcss
```

## 生成空白站点
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

## 本地运行
```bash
hugo server
```
## 构建静态文件


## 常见问题
[构建站点时出错：# POSTCSS: failed to transform "scss/main.css"](https://github.com/google/docsy/issues/235)
