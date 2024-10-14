---
title: Python 环境安装与使用
weight: 1
---

# 概述

> 参考：
>
> - [官方文档，Python 的安装与使用](https://docs.python.org/3/using/index.html)

# 安装 Python

## Linux

各 Linux 发行版通常都会自带 Python

### 自定义 Python

由于 [Python 模块与包](/docs/2.编程/高级编程语言/Python/Python%20环境安装与使用/Python%20模块与包.md) 的管理非常混乱，我们有没有办法像 Go 一样，依靠一个 GOPATH 即可统一管理呢？可以，当我们**了解了模块搜索路径的底层原理之后**，即可开始着手将 Python 的依赖都移动到指定的目录

对于 root 用户来说，GOPATH 默认在 /root/go，那我们就将 PYTHONHOME 设为 /root/python，开始吧(注意这里要用绝对路径，不要使用 `~`)

```bash
export PYTHON_VERSION="3.10"
mkdir -p /root/python/lib
cp /usr/bin/python${PYTHON_VERSION} /root/python
cp -ax -r /usr/lib/python${PYTHON_VERSION} /root/python/lib
```

准备工作完成了，此时我们只需要修改 `${PYTHONHOME}` 或者将 `/root/python` 加入 `${PATH}` 变量中即可

```bash
export PATH=/root/python:$PATH
~]# python${PYTHON_VERSION}
Python 3.10.6 (main, Nov 14 2022, 16:10:14)
>>> import sys
>>> sys.prefix
'/root/python'
>>> sys.exec_prefix
'/root/python'
```

## Windows

从官网下载 Windows 版的 exe 安装包，勾选 `Add Python ${版本号} to PATH`。安装包中包括 IDLE、pip、Python 文档。

安装完成后会提示关闭 Path 变量值的长度限制

![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/gzv1ih/1659885506889-188054a3-8a67-4039-ab87-16f8ff3a3e38.png)

### 安装完成后自动设置的内容

修改安装路径到 `D:\Tools\Python\Python${版本号}` 目录下

- 比如我安装 Python 3.10 版本，则会安装到 `D:\Tools\Python\Python310` 目录下

同时将会在用户变量中的 Path 环境变量中添加如下值：

- `D:\Tools\Python\Python310\Scripts\`
- `D:\Tools\Python\Python310`

在 `%USERPROFILE%\AppData\Local\Programs\Python\Python310\Scripts\` 目录中将会添加 pip 二进制命令。

### 安装完成后基本目录结构

Python 自动安装 pip 和 setuptools 两个包，这两个包在 site-packages 目录中生成如下目录

- pip 包包含目录
  - pip
  - pip-22.3.1.dist-info
  - 另外还会在 Scripts 目录生成 3 个可执行文件。
- setuptools 包包含目录
  - `_distutils_hack`
  - distutils-precedence.pth
  - pkg_resources
  - setuptools-65.5.0.dist-info
  - setuptools

## Windows 可嵌入的包

嵌入式分发是一个包含最小 Python 环境的 ZIP 文件。它旨在充当另一个应用程序的一部分，而不是由最终用户直接访问。

## 安装多个版本的 Python

# 初始化项目

无

# 编译 Python

> 参考：
>
> - [GitHub 项目，pyinstaller/pyinstaller](https://github.com/pyinstaller/pyinstaller)

Python 官方并未提供可以将代码编译成一个独立的可执行文件的工具。如果想要运行 Python 代码，通常来说都要先在目标机器上安装 Python 环境以及所需的依赖库。

但是，我们可以使用第三方工具，如 PyInstaller、cx_Freeze 等将 Python 代码编译成独立的可执行文件。这些工具会将 Python 解释器和所需的所有库打包到单个文件中，因此在未安装 Python 环境的计算机上也可以运行它。

使用方法：

安装

`pip install pyinstaller`

使用命令行将 Python 文件编译为可执行文件：

`pyinstaller --onefile myscript.py`

这样，您就可以在任何没有 Python 环境的计算机上运行生成的可执行文件了。

# Python 关联文件与配置

**${sys.path}** # 模块的存储目录列表，以及加载模块时读取的目录列表。为 Python 提供寻找模块位置，也就是 python 运行代码时，import 指令导入的包的保存路径。这是多个文件的组合，也就是[搜索模块的路径](/docs/2.编程/高级编程语言/Python/Python%20环境安装与使用/Python%20模块与包.md#Python%20模块管理)

Python 模块与包的关联文件通常都是在编译 Python 解释器时设置的，各个平台的默认值不太一样，这些值都会以变量的形式，存储在 Python 解释器中，当我们使用的时候，直接读取 Python 中的某些变量即可。通常都在 sys、site 模块中变量。

注意：Python 与其他语言不太一样，可以搜索包或者模块的路径非常多。这些模块或者包不一定只存放在一个单一的目录里。

- **${PYTHONPATH}** # 手动设置的目录列表
- **site.getsitepackages()** # 这个函数可以列出所有 site 模块生成的存放模块的目录列表

## pyproject.toml

> 参考：
>
> - [PEP-518](https://peps.python.org/pep-0518/)
> - [Python 包管理指南，编写 pyproject.toml](https://packaging.python.org/en/latest/guides/writing-pyproject-toml/)

根据 [PEP-518](https://peps.python.org/pep-0518/)，Python 项目的配置推荐放到项目根目录 pyproject.toml 文件中，这是一个 [TOML](docs/2.编程/无法分类的语言/TOML.md) 格式的配置文件

该文件通常有 3 个 Table

- **build-system** # **强烈推荐使用**。它允许您声明您使用哪个[构建后端](https://packaging.python.org/en/latest/glossary/#term-Build-Backend)以及构建项目所需的其他依赖项。
- **project** # 大多数构建后端用来指定项目的基本元数据的格式，例如依赖项、您的姓名、etc. 。
- **tool** # 具有特定于工具的子表，例如 `[tool.hatch]`、`[tool.black]`、`[tool.mypy]`。它的内容是由每个工具定义的。请查阅特定工具的文档以了解它可以包含什么。

### project

**name**(STRING)

**dynamic**(\[]STRING)

**dependencies**(\[]STRING)

## requirements.txt 文件（弃用）

> 参考：
>
> - [pip 官方文档，用户指南-Requirements 文件](https://pip.pypa.io/en/latest/user_guide/#requirements-files)
> - [知乎，Python 中的 requirement.txt](https://zhuanlan.zhihu.com/p/69058584)

Python 也需要维护项目相关的依赖包。通常我们会在项目的根目录下放置一个 requirement.txt 文件，用于记录所有依赖包和它的确切版本号。

requirement.txt 的内容长这样：

```python
alembic==1.0.10
appnope==0.1.0
astroid==2.2.5
attrs==19.1.0
backcall==0.1.0
bcrypt==3.1.6
bleach==3.1.0
cffi==1.12.3
Click==7.0
decorator==4.4.0
defusedxml==0.6.0
entrypoints==0.3
...
```

### 如何使用？

那么 requirement.txt 究竟如何使用呢？

当我们拿到一个项目时，首先要在项目运行环境安装 requirement.txt 所包含的依赖：

`pip install -r requirement.txt`

当我们要把环境中的依赖写入 requirement.txt 中时，可以借助 freeze 命令：

`pip freeze > requirements.txt`

### 环境混用怎么办？

在导出依赖到 requirement.txt 文件时会有一种尴尬的情况。你的本地环境不仅包含项目 A 所需要的依赖，也包含着项目 B 所需要的依赖。此时我们要如何做到只把项目 A 的依赖导出呢？

[pipreqs](https://github.com/bndr/pipreqs) 可以通过扫描项目目录，帮助我们仅生成当前项目的依赖清单。

通过以下命令安装：

`pip install pipreqs`

运行：

`pipreqs ./ --encoding utf8`

