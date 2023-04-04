---
title: "PIP"
linkTitle: "PIP"
weight: 20
---

# 概述

> 参考：
> 
> - [GitHub 项目，pypa/pip](https://github.com/pypa/pip)
> - [官网](https://pip.pypa.io/)

**Package Installer for Python(Python 的包安装器，简称 PIP)** 是 Python 的包管理程序。可以使用它来安装来自 Python 包索引和其他索引的包。从 Python 3.4 开始，它默认包含在 Python 二进制安装程序中。

# 安装 PIP

通常来说，安装 Python 时，会自动安装 PIP

# 关联文件与配置

## Windows

全局配置

-   **C:/ProgramData/pip/pip.ini**

用户配置

- **%APPDATA%/pip/pip.ini**
- **%USERPROFILE%/pip/pip.ini** # 传统的用户配置如果存在的话也会被加载

Site 配置

- **%VIRTUAL_ENV%\pip.ini**

## Unix

全局配置

- In a “pip” subdirectory of any of the paths set in the environment variable `XDG_CONFIG_DIRS` (if it exists), for example `/etc/xdg/pip/pip.conf`.
- This will be followed by loading `/etc/pip.conf`.

用户配置

- **~/.config/pip/pip.conf** # which respects the `XDG_CONFIG_HOME` environment variable.
- **~/.pip/pip.conf** # 传统的用户配置如果存在的话也会被加载

Site 配置

- **$VIRTUAL_ENV/pip.conf**

## 其他

pip 安装的模块我们可以从如下目录中找到，该目录下的目录名或文件名通常来说即是包名

- Windows
    - **%USERPROFILE%/AppData/Local/Programs/Python/Python${版本号}/Lib/site-packages/***
- Linux
    - root 用户：**/usr/local/lib/python${VERSION}/dist-packages/\***
    - 普通 用户：**~/.local/lib/python${PythonVersion}/site-packages/\***

有些包会生成一些可以执行程序，这些二进制文件则默认保存在如下目录

- Windows
    - **%USERPROFILE%/AppData/Local/Programs/Python/Python310/Scripts/**
- Linux
    - root 用户：**/usr/local/bin/**
    - 普通 用户：**~/.local/bin/**

# Syntax(语法)

> 参考：
> 
> - [官方文档，cli-pip](https://pip.pypa.io/en/stable/cli/pip/)

**pip COMMAND \[OPTIONS] COMMAND**

Commands:
- **install** # 安装包
- **download** # 下载包
- uninstall                   Uninstall packages.
- freeze                      Output installed packages in requirements format.
- inspect                     Inspect the python environment.
- list                        List installed packages.
- show                        Show information about installed packages.
- check                       Verify installed packages have compatible dependencies.
- **config** # 管理本地或全局配置
- search                      Search PyPI for packages.
- cache                       Inspect and manage pip's wheel cache.
- index                       Inspect information available from package indexes.
- wheel                       Build wheels from your requirements.
- hash                        Compute hashes of package archives.
- completion                  A helper command used for command completion.
- debug                       Show information useful for debugging.
- help                        Show help for commands.

## pip config

管理本地或全局配置

pip config list # 列出

pip config edit --editor code # 使用 vscode 打开 pip 配置文件

## pip download

**OPTIONS**

- **-d, --dest \<DIR>** # 将 Python 包下载到 DIR 目录中。
- **-r, --requirement \<FILE>** # 从指定的 requirement 文件中下载 Python 包。

## pip install

**OPTIONS**

- **-t, --target \<DIR>** # 将 Python 包安装到 DIR 目录中。可以添加 --upgrade 选项将现有包替换为 DIR 目录中的新版本。

# 最佳实践

更新 pip：`pip install --upgrade pip`

## 配置镜像源加速

对于 Python 开发用户来讲，PIP 安装软件包是家常便饭。但国外的源下载速度实在太慢，浪费时间。而且经常出现下载后安装出错问题。所以把 PIP 安装源替换成国内镜像，可以大幅提升下载速度，还可以提高安装成功率。

国内源：
新版 ubuntu 要求使用 https 源，要注意。

- 清华：<https://pypi.tuna.tsinghua.edu.cn/simple>
- 阿里云：https://mirrors.aliyun.com/pypi/simple/
- 中国科技大学 <https://pypi.mirrors.ustc.edu.cn/simple/>
- 华中理工大学：<http://pypi.hustunique.com/>
- 山东理工大学：<http://pypi.sdutlinux.org/>
- 豆瓣：<http://pypi.douban.com/simple/>

临时使用：

可以在使用 pip 的时候加参数 `-i https://pypi.tuna.tsinghua.edu.cn/simple`

例如：`pip install -i https://pypi.tuna.tsinghua.edu.cn/simple pyspider`，这样就会从清华这边的镜像去安装 pyspider 库。

永久修改，一劳永逸：

Linux 下，修改 ~/.pip/pip.conf (没有就创建一个文件夹及文件。文件夹要加“.”，表示是隐藏文件夹)，内容如下：

```bash
mkdir -p ~/.pip
tee ~/.pip/pip.conf <<EOF
[global]
index-url = https://mirrors.aliyun.com/pypi/simple/
[install]
trusted-host=mirrors.aliyun.com
EOF
```

windows 下，直接在 user 目录中创建一个 pip 目录，如：C:/Users/xx/pip，新建文件 pip.ini。内容同上。

```powershell
New-Item -ItemType File $env:APPDATA\pip\pip.ini -Force
Add-Content $env:APPDATA\pip\pip.ini "[global]"
Add-Content $env:APPDATA\pip\pip.ini "index-url = https://mirrors.aliyun.com/pypi/simple/"
```