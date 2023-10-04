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

# pip 安装包逻辑

pip 默认从 [PyPI](https://pypi.org/) 搜索包。

先下载到 /tmp/pip-unpack-随机数/包名-XXX.whl，然后默认情况下，将这些文件安装到 site-packages 目录。

## 有关 --taget 选项的说明

当我们使用 `--target TARGET_DIR` 选项指定包的安装路径时，则会将随包带的可执行文件安装到 `TARGET_DIR/bin/` 目录下，若两个具有可执行文件的包被同事安装到同一个 TARGET_DIR 中，则后安装的包的二进制文件将不会成功，并且有警告信息：

```
WARNING: Target directory /root/pythonpath/bin already exists. Specify --upgrade to force replacement.
```

假如现在想要安装 black 和 pipreqs 两个包，先安装 black，那么 pipreqs 可执行文件将不会安装成功，除非使用 --upgrade 选项，此时，bin/ 目录下的 black 可执行文件将被删除，并替换为 pipreqs 可执行文件。

综上所述：--target 选项不适合指定为 PYTHONPATH，而是为每个包指定一个独立的目录，并且每个目录下的 bin 目录都要添加到 $PATH 才可以，这是一个很鸡肋的选项。

详见 [pip issue 8063](https://github.com/pypa/pip/issues/8063)

# 安装 PIP

通常来说，安装 Python 时，会自动安装 PIP

> 一般都会被安装到 stie-packages/ 目录中

pip 包通常包含两个目录

- pip/
- pip-${VERSION}.dist-info/

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

**Commands**

- **install** # 安装包
- **download** # 下载包
- **uninstall** # 卸载包
- **freeze** # 以 requirements 格式输出已安装的软件包
- inspect                     Inspect the python environment.
- **list** # 列出已安装的包
- **show** # 显示有关已安装软件包的信息。
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

**OPTIONS**

- **-v, --verbose** # 在执行命令时显示更多信息，可多次使用，最多 3 次，每多一个 v，则显示的信息就更多一些。

## pip config

管理本地或全局配置

pip config list # 列出

pip config edit --editor code # 使用 vscode 打开 pip 配置文件

## pip download

**OPTIONS**

- **-d, --dest \<DIR>** # 将 Python 包下载到 DIR 目录中。
- **-r, --requirement \<FILE>** # 从指定的 requirement 文件中下载 Python 包。

## pip freeeze

以 requirements 格式输出安装的软件包列表。软件包以不区分大小写的排序方式列出。

`pip freeze > requirements.txt`

## pip install

https://pip.pypa.io/en/stable/cli/pip_install/

**OPTIONS**

- **-i, --index-url \<URL>** # Python 包索引的 URL。`默认值：https://pypi.org/simple`
    - 可以通过 -i 配置 pip 安装包时使用国内的源，避免国从国外下载速度太慢
- **-r, --requirement \<FILE>** # 安装指定 requirement 文件中的 Python 包。
- **-t, --target \<DIR>** # 将 Python 包安装到 DIR 目录中。可以添加 --upgrade 选项将现有包替换为 DIR 目录中的新版本。
- **-U, --upgrade** # 将指定的所有 Python 包升级到最新的可用版本。
- **--user** # 将包安装到 user 的 site-packages 目录下。

### EXAMPLE

安装 2.6.1.3 版本的 paddleocr

- `pip install paddleocr==2.6.1.3`

更新 paddleocr 包到最新版本，显示详细信息。

- `pip install paddleocr --upgrade --verbose`

# 最佳实践

安装、下载包时，指定包源

`pip install -i https://pypi.tuna.tsinghua.edu.cn/simple pyspider`，这样就会从清华这边的包镜像安装 pyspider 库。

## 离线安装 Python 包

下载 Python 包文件

`pip download black`

安装 Python 包文件

`pip install *.whl`

## 配置镜像源加速

对于 Python 开发用户来讲，PIP 安装软件包是家常便饭。但国外的源下载速度实在太慢，浪费时间。而且经常出现下载后安装出错问题。所以把 PIP 安装源替换成国内镜像，可以大幅提升下载速度，还可以提高安装成功率。

国内源：

- 阿里云：https://mirrors.aliyun.com/pypi/simple/
- 清华：<https://pypi.tuna.tsinghua.edu.cn/simple>
- 中国科技大学 <https://pypi.mirrors.ustc.edu.cn/simple/>
- 华中理工大学：<http://pypi.hustunique.com/>
- 山东理工大学：<http://pypi.sdutlinux.org/>
- 豆瓣：<http://pypi.douban.com/simple/>

Linux 下，修改 ~/.pip/pip.conf (没有就创建一个文件夹及文件。文件夹要加“.”，表示是隐藏文件夹)，内容如下：

```bash
mkdir -p ~/.pip
tee ~/.pip/pip.conf > /dev/null <<EOF
[global]
index-url = https://mirrors.aliyun.com/pypi/simple/
trusted-host = mirrors.aliyun.com
EOF
```

windows 下，直接在 user 目录中创建一个 pip 目录，如：C:/Users/xx/pip，新建文件 pip.ini。内容同上。

```powershell
New-Item -ItemType File $env:APPDATA\pip\pip.ini -Force
Add-Content $env:APPDATA\pip\pip.ini "[global]"
Add-Content $env:APPDATA\pip\pip.ini "index-url = https://mirrors.aliyun.com/pypi/simple/"
```

## 配置默认安装路径

TODO: --target 没效果，修改 user 的 site-packages 没效果。。。。o(╯□╰)o

先用 [Python 虚拟环境](/docs/2.编程/高级编程语言/Python/Python%20环境安装与使用/Python%20模块与包.md#Python%20虚拟环境)吧

## 显示一个包的可安装版本

pip 没有这种命令，但是可以通过 `pip install PACKAGE== --no-deps --no-cache-dir` 实现类似的效果，因为版本号为空，则会报错，在报错时，会显示所有可用的版本号。

比如执行 `pip install numpy== --no-deps --no-cache-dir`

报错：

```
ERROR: Could not find a version that satisfies the requirement numpy== (from versions: 1.3.0, 1.4.1, 1.5.0, 1.5.1, 1.6.0, 1.6.1, 1.6.2, 1.7.0, 1.7.1, 1.7.2, 1.8.0, 1.8.1, 1.8.2, 1.9.0, 1.9.1, 1.9.2, 1.9.3, 1.10.0.post2, 1.10.1, 1.10.2, 1.10.4, 1.11.0, 1.11.1, 1.11.2, 1.11.3, 1.12.0, 1.12.1, 1.13.0, 1.13.1, 1.13.3, 1.14.0, 1.14.1, 1.14.2, 1.14.3, 1.14.4, 1.14.5, 1.14.6, 1.15.0, 1.15.1, 1.15.2, 1.15.3, 1.15.4, 1.16.0, 1.16.1, 1.16.2, 1.16.3, 1.16.4, 1.16.5, 1.16.6, 1.17.0, 1.17.1, 1.17.2, 1.17.3, 1.17.4, 1.17.5, 1.18.0, 1.18.1, 1.18.2, 1.18.3, 1.18.4, 1.18.5, 1.19.0, 1.19.1, 1.19.2, 1.19.3, 1.19.4, 1.19.5, 1.20.0, 1.20.1, 1.20.2, 1.20.3, 1.21.0, 1.21.1, 1.21.2, 1.21.3, 1.21.4, 1.21.5, 1.21.6, 1.22.0, 1.22.1, 1.22.2, 1.22.3, 1.22.4, 1.23.0rc1, 1.23.0rc2, 1.23.0rc3, 1.23.0, 1.23.1, 1.23.2, 1.23.3, 1.23.4, 1.23.5, 1.24.0rc1, 1.24.0rc2, 1.24.0, 1.24.1, 1.24.2, 1.24.3, 1.24.4, 1.25.0rc1, 1.25.0, 1.25.1, 1.25.2, 1.26.0b1, 1.26.0rc1, 1.26.0, 1.26.1)
ERROR: No matching distribution found for numpy==
```