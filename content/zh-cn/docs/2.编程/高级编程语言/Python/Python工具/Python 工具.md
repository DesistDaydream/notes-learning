---
title: Python 工具
weight: 1
---

# 概述

> 参考：

# python

> 参考：
> 
> - [官方文档，Python 安装与使用-命令行和环境](https://docs.python.org/3/using/cmdline.html)

python 是一个工具，用来管理 Python 语言编写的代码。

## Syntax(语法)

**python \[-bBdEhiIOqsSuvVWx?] \[-c command | -m module-name | script | - ] \[args]**
OPTIONS

- **-S** # Python 启动初始化时，不要导入 site 包
- **-m \<ModuleName>** # 在 `sys.path` 中搜索指定模块，并默认执行模块中 `__name__` 为 `__main__` 的代码
  - `python3 -m site` 等效于 `python3 /usr/lib/python3.8/site.py`

## EXAMPLE

## 启用一个简易的 HTTP 服务器

```
# 使用该命令可以在当前目录搭建一个简易的http服务器，当client访问的时候，就可以直接看到该目录下的内容，还可以下载该目录下的内容
python -m SimpleHTTPServer NUM
```

若报错则使用如下命令：

```bash
python3 -m http.server NUM
```
