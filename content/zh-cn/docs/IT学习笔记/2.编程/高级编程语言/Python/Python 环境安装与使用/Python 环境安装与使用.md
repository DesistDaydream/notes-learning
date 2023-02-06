---
title: Python 环境安装与使用
---

# 概述

> 参考：
> - [官方文档，Python 的安装与使用](https://docs.python.org/3/using/index.html)

# 安装 Python

## Linux

各 Linux 发行版通常都会自带 Python

## Windows

从官网下载 Windows 版的 exe 安装包，勾选 `Add Python ${版本号} to PATH`。安装包中包括 IDLE、pip、Python 文档。

安装完成后会提示关闭 Path 变量值的长度限制
![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/gzv1ih/1659885506889-188054a3-8a67-4039-ab87-16f8ff3a3e38.png)

### 安装完成后自动设置的内容

默认安装到 `%USERPROFILE%\AppData\Local\Programs\Python\Python${版本号}` 目录下

- 比如我安装 Python 3.10 版本，则会安装到 `%USERPROFILE%\AppData\Local\Programs\Python\Python310` 目录下

同时将会在用户变量中的 Path 环境变量中添加如下值：

- %USERPROFILE%/AppData/Local/Programs/Python/Python310/Scripts/
- %USERPROFILE%/AppData/Local/Programs/Python/Python310/

在 `%USERPROFILE%\AppData\Local\Programs\Python\Python310\Scripts\` 目录中将会添加 pip 二进制命令。

## Windows 可嵌入的包

嵌入式分发是一个包含最小 Python 环境的 ZIP 文件。它旨在充当另一个应用程序的一部分，而不是由最终用户直接访问。

## 安装多个版本的 Python

# 编译 Python
> 参考：
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

**XXXX** # 寻找 Python 代码的位置，也就是 python 编译或运行代码时，import 指令导入的包的保存路径。这是多个文件的组合，也就是[搜索模块的路径](/docs/IT学习笔记/2.编程/高级编程语言/Python/Python%20环境安装与使用/Python%20模块与包.md#Python%20模块管理)

- 注意：Python 与其他语言不太一样，可以搜索包或者模块的路径非常多。这些模块或者包不一定只存放在一个单一的目录里。
