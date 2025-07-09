---
title: Python 工具
weight: 1
---

# 概述

> 参考：

# 安装 Python 包/模块

> 参考：
>
> - [官方文档，安装 Python 模块](https://docs.python.org/3.10/installing/index.html)
> - <https://frostming.com/2019/03-13/where-do-your-packages-go/>

管理 Python 的模块和包所在路径非常乱，不知道是何原因。

[PIP](/docs/2.编程/高级编程语言/Python/Python工具/PIP.md) 是首选的安装程序。从 Python 3.4 开始，它默认包含在 Python 二进制安装程序中。就算你是用 pipenv，poetry，底层依然是 pip，一律适用。

运行 pip 有两种方式：

- pip ...
- python -m pip ...

第一种方式和第二种方式大同小异，区别是第一种方式使用的 Python 解释器是写在 pip 文件的 shebang 里的，一般情况下，如果你的 pip 路径是 $path\_prefix/bin/pip，那么 Python 路径对应的就是 $path\_prefix/bin/python。如果你用的是 Unix 系统则 cat $(which pip) 第一行就包含了 Python 解释器的路径。第二种方式则显式地指定了 Python 的位置。这条规则，对于所有 Python 的可执行程序都是适用的。流程如下图所示。

![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/python/1669286382022-472bf4de-24cf-4652-bc94-3d52d01f7df1.png)

那么，不加任何自定义配置时，使用 pip 安装包就会自动安装到 `$path_prefix/lib/pythonX.Y/site-packages` 下（$path_prefix 是从上一段里得到的），可执行程序安装到 $path_prefix/bin 下，如果需要在命令行直接使用 my_cmd 运行，记得加到 PATH。

刚刚安装完的 Python 一般只有 pip 和 setuptools 模块，site-packages 目录下内容如下：

```
$ tree -L 1
.
├── README.txt
├── _distutils_hack
├── distutils-precedence.pth
├── pip
├── pip-23.0.1.dist-info
├── pkg_resources
├── setuptools
└── setuptools-65.5.0.dist-info
```

## 包管理工具

[PIP](/docs/2.编程/高级编程语言/Python/Python工具/PIP.md)

Rye

- https://github.com/astral-sh/rye

poetry

- https://github.com/python-poetry/poetry

pipenv

pyenv

pdm

hatch

uv # 依赖库有全局缓存

- https://github.com/astral-sh/uv

# Wheel 包

Wheel 是一种类似压缩包的 Python 用于分发包的文件，有点类似 .rpm、.deb。

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

```bash
# 使用该命令可以在当前目录搭建一个简易的http服务器，当client访问的时候，就可以直接看到该目录下的内容，还可以下载该目录下的内容
python -m SimpleHTTPServer NUM
```

若报错则使用如下命令：

```bash
python3 -m http.server NUM
```
