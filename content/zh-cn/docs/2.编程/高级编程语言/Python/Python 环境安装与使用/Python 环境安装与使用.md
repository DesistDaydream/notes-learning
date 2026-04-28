---
title: Python 环境安装与使用
linkTitle: Python 环境安装与使用
weight: 1
---

# 概述

> 参考：
>
> - [官方文档，Python 的安装与使用](https://docs.python.org/3/using/index.html)

# 安装 Python

## Linux

各 Linux 发行版通常都会自带 Python。要想在 Linux 中安装指定版本的 Python，需要使用源码编译安装，因为 Python 依赖很多 C 库，而很多 Linux 发行版中的程序又依赖自带的 Python。

所以我们通常需要把自己想要安装的 Python 装在其他目录，而不是更新发行版自带的 Python。

### 编译 Python 环境

> [!Question]
> 由于 [Python 模块与包](/docs/2.编程/高级编程语言/Python/Python%20环境安装与使用/Python%20模块与包.md) 的管理非常混乱，我们有没有办法像 Go 一样，依靠一个 GOPATH 即可统一管理呢？可以，当我们**了解了模块搜索路径的底层原理之后**，即可开始着手将 Python 的依赖都移动到指定的目录

首先，下载 Python 源码包。解压后进入源码目录

生成 Makefile 文件

```bash
./configure --enable-optimizations --prefix=/usr/local/python3.12
```

构建 Python 环境，但不安装到系统中

```bash
make
```

确认构建无误后（参考下文中的 [编译注意事项](#编译注意事项)），将 Python 安装到系统环境中（i.e. ./configure 的 --prefix 参数执行的目录）

```bash
make install
```

之后会在 `/usr/local/python3.12/` 目录中生成可用的 Python 环境。

准备工作完成了，此时我们只需要将 `which python3` 的软链接指到 `/usr/local/python3.10/bin/python3` 即可

```bash
export PATH=/root/python:$PATH
~]# python
Python 3.10.6 (main, Nov 14 2022, 16:10:14)
>>> import sys
>>> sys.prefix
'/usr/local/python3.10'
>>> sys.exec_prefix
'/usr/local/python3.10'
```

> [!Attention] yum, dnf 的问题
> 更换 Python 后，可能会由于依赖不足无法使用。而 `dnf` 模块本身是用 C 扩展写的，它的 `.so` 文件编译时绑定了特定的 Python ABI。系统里的 dnf 是针对 python3.9 编译安装的，直接用 3.10 加载会报 ABI 不兼容或者直接找不到模块。
> 可以使用如下方式将 yum 和 dnf 调用的 Python 改回 3.0
> ```bash
> sed -i '1s|.*|#!/usr/bin/python3.9|' /usr/bin/yum
> sed -i '1s|.*|#!/usr/bin/python3.9|' /usr/bin/dnf
> ```

### 编译注意事项

> [!Attention] 编译出现问题时
> 要清理环境（执行 `make clean` 命令）后，重新执行 `./configure XXX && make` 命令。注意：不要执行 make install，避免直接生成最终 Python 环境。 

**编译 Python 时所需的依赖**

在使用 `make` 命令编译的时候可能会出现这种提示：

```  bash
The necessary bits to build these optional modules were not found:  
_curses               _curses_panel         _dbm                 
_gdbm                 _tkinter              nis                  
readline                                                         
To find the necessary bits, look in configure.ac and config.log.  
```

这是提示编译 Python 环境的时候缺少一些模块，虽然这些模块的确实不会导致编译失败。但是在后续使用的过程中，如果代码依赖了这些模块，那么将无法使用，必须在安装完这些模块后<font color="#ff0000">**重新编译**</font>才可以解决。

> [!Tip] Python 环境所需的这些 necessary 的依赖模块需要使用 [Package 管理](/docs/1.操作系统/Package%20管理/Package%20管理.md) 安装，这些模块通常都是与系统自身 C 环境相关的，无法通过外部环境引入或包含在源码中。

> [!Attention] 各种 Linux 发行版的环境缺失的模块可能并不一样，Python 版本之间所需要的必要模块可能也不一样。这里只能尽量总结一些。

```bash
dnf install gcc make openssl-devel libffi-devel zlib-devel \
  bzip2-devel readline-devel sqlite-devel ncurses-devel \
  xz-devel tk-devel
```

## Windows

从官网下载 Windows 版的 exe 安装包，勾选 `Add Python ${版本号} to PATH`。安装包中包括 IDLE、pip、Python 文档。

安装完成后会提示关闭 Path 变量值的长度限制

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/gzv1ih/1659885506889-188054a3-8a67-4039-ab87-16f8ff3a3e38.png)

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

# Python 编译

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
> - [Python 包用户指南，编写 pyproject.toml](https://packaging.python.org/en/latest/guides/writing-pyproject-toml/)
> - [Python 包用户指南，规范 - pyproject.tom 规范](https://packaging.python.org/en/latest/specifications/pyproject-toml/)

根据 [PEP-518](https://peps.python.org/pep-0518/)，Python 项目的配置推荐放到项目根目录 pyproject.toml 文件中，这是一个 [TOML](/docs/2.编程/无法分类的语言/TOML.md) 格式的配置文件

该文件通常有 3 个 Table

- **build-system** # **强烈推荐使用**。它允许您声明您使用哪个[构建后端](https://packaging.python.org/en/latest/glossary/#term-Build-Backend)以及构建项目所需的其他依赖项。
- **project** # 大多数构建后端用来指定项目的基本元数据的格式，例如依赖项、您的姓名、etc. 。
- **tool** # 具有特定于工具的子表，例如 `[tool.hatch]`、`[tool.black]`、`[tool.mypy]`。它的内容是由每个工具定义的。请查阅特定工具的文档以了解它可以包含什么。

最简单的 pyproject.tom 文件应该至少包含如下内容（不能只写 name 不写 version）：

```toml
[project]
name = "python-learning"
version = "0.1.0"
```

### build-system

### project

**name**(STRING)

**version**(STRING)

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
