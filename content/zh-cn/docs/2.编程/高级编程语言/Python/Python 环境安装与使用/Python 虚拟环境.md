---
title: "Python 虚拟环境"
linkTitle: "Python 虚拟环境"
date: "2023-07-01T16:03"
weight: 20
---

# 概述

> 参考：
>
> - [官方文档，Python 教程 - 虚拟环境和包](https://docs.python.org/3/tutorial/venv.html)

Python 没有 go.mod 与 go.sum 这种文件来管理第三方依赖模块的版本。假如现在只有一个 3.10 版本的 Python，那么所有项目的依赖模块都会被安装到 site-packages 目录中，且 site-packages 目录中保存只能保存唯一一个版本的依赖。若项目 A 需要  模块 C 的 1.0 版本，项目 B 需要 模块 C 的 2.0 的版本，这时候就会产生冲突，若同时运行这两个项目，将有其中一个无法正常运行。因为只要安装 C 1.0，当前环境的 C 2.0 就会被覆盖。

为了解决上述**第三方模块的多版本管理问题**，Python 想了一个类似 JS 的 node_modules 方案。。称为 [virtual environment(虚拟环境，简称 venv)](https://docs.python.org/zh-cn/3/glossary.html#term-virtual-environment)。就相当于为每个项目建立一个独立的 Python 环境。。。。

> 但是。。。另一个可怕的问题就是。。。。如果多个项目依赖相同版本的模块。。那么。。就要安装很多遍。。。。。唉。。。。。。
> 
> TODO: 如何解决这个问题？

想要使用 Python 虚拟环境，需要先安装一个名为 [venv](https://docs.python.org/3/library/venv.html#module-venv) 的内置模块（这个模块是在 CPyhon 的源码中的，在 sys 库或用pip命令里看不到），但是有的发行版并没有随着 Python 一起安装，比如 Ubuntu，需要手动安装 `apt install python3.10-venv`

## 创建虚拟环境

假设我们现在有一个项目，放在单独的目录中，project-venv-demo，想要让这个项目有独立的依赖环境，那就执行如下命令即可

```bash
# 创建一个虚拟环境目录
~]# python3 -m venv /root/tmp/project-venv-demo
~]# ls /root/tmp
project-venv-demo
# 激活虚拟环境
# source /root/tmp/project-venv-demo/bin/activate
(project-venv-demo) ~]#
```

此时的 Shell 中的提示符前面出现了 `(project-venv-demo)`，这就说明当前已在 Python 的虚拟环境那种了。此时的虚拟环境中将是已安装的特定 Python 版本

```bash
(project-venv-demo) ~]# python -m site
sys.path = [
    '/root/tmp',
    '/usr/lib/python310.zip',
    '/usr/lib/python3.10',
    '/usr/lib/python3.10/lib-dynload',
    '/root/tmp/project-venv-demo/lib/python3.10/site-packages',
]
USER_BASE: '/root/.local' (doesn't exist)
USER_SITE: '/root/.local/lib/python3.10/site-packages' (doesn't exist)
ENABLE_USER_SITE: False
```

在任意位置执行 `deactivate` 命令将会退出当前虚拟环境。

## 虚拟环境关联文件与配置



**${VENV}** # 虚拟化境创建时指定的目录。如果是 VSCode 创建的一般都是当前目录下的 `.venv/` 目录

- **./pyvenv.cfg** #
- **./Lib/site-packages/** # 第三方库保存路径
- **./Lib/Scripts/** # 部分库所需要的可执行文件的保存路径

