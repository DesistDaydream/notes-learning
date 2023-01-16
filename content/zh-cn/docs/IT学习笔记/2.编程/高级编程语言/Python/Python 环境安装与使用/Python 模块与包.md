---
title: Python 模块与包
---

# 概述

> 参考：
> - [官方文档，教程-6.模块](https://docs.python.org/3/tutorial/modules.html)
> - [廖雪峰 Python 教程，模块](https://www.liaoxuefeng.com/wiki/1016959663602400/1017454145014176)

在计算机程序的开发过程中，随着程序代码越写越多，在一个文件里代码就会越来越长，越来越不容易维护。

为了编写可维护的代码，我们把很多函数分组，分别放到不同的文件里，这样，每个文件包含的代码就相对较少，很多编程语言都采用这种组织代码的方式。在 Python 中，**一个 **`**.py**`** 文件**就**称之为**一个 **Module(模块)**。

使用模块有什么好处？

最大的好处是大大提高了代码的可维护性。其次，编写代码不必从零开始。当一个模块编写完毕，就可以被其他地方引用。我们在编写程序的时候，也经常引用其他模块，包括 Python 内置的模块和来自第三方的模块。

使用模块还可以避免函数名和变量名冲突。相同名字的函数和变量完全可以分别存在不同的模块中，因此，我们自己在编写模块时，不必考虑名字会与其他模块冲突。但是也要注意，尽量不要与内置函数名字冲突。[这里](http://docs.python.org/3/library/functions.html)可以查看 Python 的所有内置函数。

你也许还想到，如果不同的人编写的模块名相同怎么办？为了避免模块名冲突，Python 又引入了按目录来组织模块的方法，称为 **Package(包)**。

举个例子，一个`abc.py`的文件就是一个名字叫`abc`的模块，一个`xyz.py`的文件就是一个名字叫`xyz`的模块。

现在，假设我们的`abc`和`xyz`这两个模块名字与其他模块冲突了，于是我们可以通过包来组织模块，避免冲突。方法是选择一个顶层包名，比如`mycompany`，按照如下目录存放：

    mycompany
    ├─ __init__.py
    ├─ abc.py
    └─ xyz.py

引入了包以后，只要顶层的包名不与别人冲突，那所有模块都不会与别人冲突。现在，`abc.py`模块的名字就变成了`mycompany.abc`，类似的，`xyz.py`的模块名变成了`mycompany.xyz`。

请注意，每一个包目录下面都会有一个`__init__.py`的文件，这个文件是必须存在的，否则，Python 就把这个目录当成普通目录，而不是一个包。`__init__.py`可以是空文件，也可以有 Python 代码，因为`__init__.py`本身就是一个模块，而它的模块名就是`mycompany`。

类似的，可以有多级目录，组成多级层次的包结构。比如如下的目录结构：

```bash
mycompany
 ├─ web
 │  ├─ __init__.py
 │  ├─ utils.py
 │  └─ www.py
 ├─ __init__.py
 ├─ abc.py
 └─ utils.py
```

文件`www.py`的模块名就是`mycompany.web.www`，两个文件`utils.py`的模块名分别是`mycompany.utils`和`mycompany.web.utils`。

自己创建模块时要注意命名，不能和 Python 自带的模块名称冲突。例如，系统自带了 sys 模块，自己的模块就不可命名为 sys.py，否则将无法导入系统自带的 sys 模块。

`mycompany.web`也是一个模块，该模块对应的 .py 文件是 **init**.py。

## Package(包)

包是一种通过使用“带点的模块名称”来构建 Python 模块命名空间的方式。例如，模块名称 AB 指定了一个名为 A 的包中的一个名为 B 的子模块。就像使用模块使不同模块的作者不必担心彼此的全局变量名称一样，使用点分模块名称可以节省作者多模块包，如 NumPy 或 Pillow，不必担心彼此的模块名称。

假设您要设计一组模块（“包”）来统一处理声音文件和声音数据。有许多不同的声音文件格式（通常通过它们的扩展名识别，例如：.wav、.aiff、.au），因此您可能需要创建和维护不断增长的模块集合，以便在各种文件格式之间进行转换。您可能还想对声音数据执行许多不同的操作（例如混合、添加回声、应用均衡器功能、创建人工立体声效果），因此您将编写一个永无止境的模块流来执行这些操作。这是您的包的可能结构（以分层文件系统表示）：

导入包时，Python 会在 sys.path 上的目录中搜索包子目录。

需要 `__init__.py` 文件才能使 Python 将包含该文件的目录视为包。这可以防止具有通用名称（例如字符串）的目录无意中隐藏了稍后出现在模块搜索路径上的有效模块。在最简单的情况下，`__init__.py` 可以只是一个空文件，但它也可以执行包的初始化代码或设置 `__all__` 变量，稍后将介绍。

注意：

- import 指令导入模块时是否成功，取决于执行 Python 代码时的所在位置。反正设计的挺恶心的。。。o(╯□╰)o
- 说白了：Package 就是 Module，只不过是对 Module 进行了分类~~~导入模块时，模块名称中有 `.` 就是包了~~其实还是模块。也可以这么理解：一个目录就是一个包，一个文件就是一个模块。
- 若没有 `__init__.py` 文件，我们使用 IDE 追踪代码时，如果想要从 import 指令中的导入的包进行追踪，会提示找不到该文件，无法追踪到。

## 总结

模块是一组 Python 代码的集合，可以使用其他模块，也可以被其他模块使用。

创建自己的模块时，要注意：

- 模块名要遵循 Python 变量命名规范，不要使用中文、特殊字符；
- 模块名不要和系统模块名冲突，最好先查看系统是否已存在该模块，检查方法是在 Python 交互环境执行`import abc`，若成功则说明系统存在此模块。

# Python 模块的使用

Python 本身就内置了很多非常有用的模块，只要安装完毕，这些模块就可以立刻使用。

我们以内建的`sys`模块为例，编写一个`hello`的模块：

```python
#!/usr/bin/env python3
# -*- coding: utf-8 -*-

' a test module '

__author__ = 'Michael Liao'

import sys

def test():
    args = sys.argv
    if len(args)==1:
        print('Hello, world!')
    elif len(args)==2:
        print('Hello, %s!' % args[1])
    else:
        print('Too many arguments!')

if __name__=='__main__':
    test()
```

第 1 行和第 2 行是标准注释，第 1 行注释可以让这个`hello.py`文件直接在 Unix/Linux/Mac 上运行，第 2 行注释表示.py 文件本身使用标准 UTF-8 编码；

第 4 行是一个字符串，表示模块的文档注释，任何模块代码的第一个字符串都被视为模块的文档注释；

第 6 行使用`__author__`变量把作者写进去，这样当你公开源代码后别人就可以瞻仰你的大名；

以上就是 Python 模块的标准文件模板，当然也可以全部删掉不写，但是，按标准办事肯定没错。

后面开始就是真正的代码部分。

你可能注意到了，使用`sys`模块的第一步，就是导入该模块：

```python
import sys
```

导入`sys`模块后，我们就有了变量`sys`指向该模块，利用`sys`这个变量，就可以访问`sys`模块的所有功能。

`sys`模块有一个`argv`变量，用 list 存储了命令行的所有参数。`argv`至少有一个元素，因为第一个参数永远是该.py 文件的名称，例如：

运行`python3 hello.py`获得的`sys.argv`就是`['hello.py']`；

运行`python3 hello.py Michael`获得的`sys.argv`就是`['hello.py', 'Michael']`。

最后，注意到这两行代码：

```python
if __name__=='__main__':
    test()
```

当我们在命令行直接使用 `python3 hello.py` 运行模块文件时，Python 解释器把 `__name__` 这个特殊变量的值设为 `__main__`，而如果在其他地方导入`hello`模块时 `__name__` 变量的值则是本模块的名称(i.e.hello)，因此，这种 `if` 测试可以让一个模块通过命令行运行时执行一些额外的代码，最常见的就是运行单元测试。

我们可以用命令行运行`hello.py`看看效果：

```bash
$ python3 hello.py
Hello, world!
$ python3 hello.py Michael
Hello, Michael!
```

如果启动 Python 交互环境，再导入`hello`模块：

```python
$ python3
Python 3.4.3 (v3.4.3:9b73f1c3e601, Feb 23 2015, 02:52:03)
[GCC 4.2.1 (Apple Inc. build 5666) (dot 3)] on darwin
Type "help", "copyright", "credits" or "license" for more information.
>>> import hello
>>>
```

导入时，没有打印`Hello, word!`，因为没有执行`test()`函数。

调用`hello.test()`时，才能打印出`Hello, word!`：

```python
>>> hello.test()
Hello, world!
```

## Compiled(已编译) 的 Python 文件

为了加快加载模块的速度，Python 将每个模块的编译版本缓存在名为 module.version.pyc 的 **pycache** 目录下，其中版本编码了编译文件的格式；它通常包含 Python 版本号。例如，在 CPython 3.3 版中，spam.py 的编译版本将被缓存为 **pycache**/spam.cpython-33.pyc。这种命名约定允许来自不同版本和不同 Python 版本的编译模块共存。

Python 会根据编译后的版本检查源代码的修改日期，以查看它是否已过时并需要重新编译。这是一个完全自动的过程。此外，编译后的模块与平台无关，因此可以在不同架构的系统之间共享同一个库。

Python 在两种情况下不检查缓存。首先，它总是重新编译并且不存储直接从命令行加载的模块的结果。其次，如果没有源模块，它不会检查缓存。要支持非源（仅编译）分发，编译模块必须在源目录中，并且不能有源模块。

## 作用域

在一个模块中，我们可能会定义很多函数和变量，但有的函数和变量我们希望给别人使用，有的函数和变量我们希望仅仅在模块内部使用。在 Python 中，是通过`_`前缀来实现的。

正常的函数和变量名是公开的（public），可以被直接引用，比如：`abc`，`x123`，`PI`等；

类似`__xxx__`这样的变量是特殊变量，可以被直接引用，但是有特殊用途，比如上面的`__author__`，`__name__`就是特殊变量，`hello`模块定义的文档注释也可以用特殊变量`__doc__`访问，我们自己的变量一般不要用这种变量名；

类似`_xxx`和`__xxx`这样的函数或变量就是非公开的（private），不应该被直接引用，比如`_abc`，`__abc`等；

之所以我们说，private 函数和变量“不应该”被直接引用，而不是“不能”被直接引用，是因为 Python 并没有一种方法可以完全限制访问 private 函数或变量，但是，从编程习惯上不应该引用 private 函数或变量。

private 函数或变量不应该被别人引用，那它们有什么用呢？请看例子：

```python
def _private_1(name):
    return 'Hello, %s' % name

def _private_2(name):
    return 'Hi, %s' % name

def greeting(name):
    if len(name) > 3:
        return _private_1(name)
    else:
        return _private_2(name)
```

我们在模块里公开`greeting()`函数，而把内部逻辑用 private 函数隐藏起来了，这样，调用`greeting()`函数不用关心内部的 private 函数细节，这也是一种非常有用的代码封装和抽象的方法，即：

外部不需要引用的函数全部定义成 private，只有外部需要引用的函数才定义为 public。

# Python 模块管理

> 参考：
> - <https://frostming.com/2019/03-13/where-do-your-packages-go/>
> - [官方文档，Python 教程-6.模块-模块搜索路径](https://docs.python.org/3/tutorial/modules.html#the-module-search-path)
> - [官方文档，Python 的安装和使用-命令行工具和环境](https://docs.python.org/3/using/cmdline.html)
>   - [PYTHONPATH](https://docs.python.org/3/using/cmdline.html#envvar-PYTHONPATH)
>   - [PYTHONHOME](https://docs.python.org/3/using/cmdline.html#envvar-PYTHONHOME)

我们通过 **Python 模块的搜索路径**来管理 Python 模块，或者称为管理 Python 包。Python 模块(Python 包) 的路径通常保存在 Python 的 `sys.path` 变量中(就好像 Go 包的路径保存在 GOPATH 变量中似的)。

当导入一个名为 spam 的模块时，Python 解释器首先搜索具有该名称的“内置模块”(内置模块可以用过 `sys.builtin_module_names` 获取)，若没找到，则会在 sys 包中的 path 变量(i.e.`sys.path`)中给出的目录列表中搜索名为 `spam.py` 的文件。

Python 程序启动后，通过如下几个步骤生成 `sys.path` 变量的值：

- **添加执行 Python 代码文件所在的绝对路径**
- **添加 ${PYTHONPATH} 环境变量指定的路径**
  - 这是一个目录列表，类似于类 Unix 中的 `$PATH` 变量，可以通过 `os.path` 获取其值。
  - 可以使用 `os.path.append()` 为 `$PYTHONPATH` 变量添加新的目录条目以便导入想要的模块。也可以直接设置 Linux 系统中的 `$PYTHONPATH` 变量。当项目大，需要对文件进行分类时，非常有用。
- **添加保存 Python 标准库和第三方库的路径**
  - 这是根据 Python 安装路径决定的默认值。这些是 sys.path 中的绝大部分目录。
  - 按照惯例，通常包含：
    - 与平台相关的 Python 基本库保存路径。
    - 由[站点模块](https://docs.python.org/3/library/site.html#module-site)处理的 site-packages 目录。使用 pip 安装的第三方库通常来说会存在 site-packages 目录中。

上述三个路径在 Python 启动时被初始化。我们可以通过 Python 中的 `sys.path` 数组变量查看这些路径。

```bash
$ cat test/test.py
import sys
print(sys.path)

$ python3 test/test.py
['/mnt/d/Projects/DesistDaydream/python-learning', '/usr/lib/python38.zip', '/usr/lib/python3.8', '/usr/lib/python3.8/lib-dynload', '/home/lichenhao/.local/lib/python3.8/site-packages', '/usr/local/lib/python3.8/dist-packages', '/usr/lib/python3/dist-packages']
```

```python
>>> import sys
>>> sys.path
['', '/usr/lib/python38.zip', '/usr/lib/python3.8', '/usr/lib/python3.8/lib-dynload', '/home/lichenhao/.local/lib/python3.8/site-packages', '/usr/local/lib/python3.8/dist-packages', '/usr/lib/python3/dist-packages']
```

模块的路径通常由 Python 中的 `${sys.prefix}` 变量的值作为前缀，并在 Python 启动后推导出完整路径：

> ${} 中的 X 和 Y 分别是 Python 的大版本号和小版本号，比如 Python3.10、Python310

- **${PWD}** # 当前工作目录
- **${sys.prefix}/lib/**
  - **./python${XY}.zip **#
  - **./python${X.Y}/** # \*\* \*\*标准库路径
  - **./python${X.Y}/lib-dynload/** #
  - **./python${X.Y}/sist-packages/** # 第三方库路径。该目录在 Ubuntu 系统中名称为 dist-packages
- **${site.USER_SITE}** # 启动用户 site 后，保存第三方库的路径。

注意：从这里可以看到，不同 Python 版本的三方库路径不同，如果把 Python 从 3.8 升级到 3.9 那么之前装的三方库都没法用了。当然你可以整个文件夹都拷贝过去，大部分情况不会出问题。

## sys.path 列表生成逻辑

> 参考：
> - [官方文档，Python 标准库-导入模块-sys.path 模块搜索路径初始化](https://docs.python.org/3/library/sys_path_init.html)
> - [官方文档，Python 标准库-Python 运行时服务-site—特定于 site 的配置](https://docs.python.org/3/library/site.html)

在上面的示例中，`sys.path` 的第一个元素（`/mnt/e/Projects/DesistDaydream/python-learning/test`）是 test.py 文件所在路径，即执行的 Python 代码文件所在路径，如果不是执行的 Python 代码文件，则第一个元素为空。每当执行一个 Python 文件时，就相当于默认执行了 `sys.path.append("文件所在绝对路径")` 代码。

`sys.path` 变量中除了第一个元素之外，其他的元素是**通过 Python 解释器(i.e.python 可执行文件)推导出来**的。启动 Python 交互环境或者用解释器运行脚本时，将会为如下几个变量生成值

- `sys.prefix` # Python 标准模块(标准库)目录前缀。可以用过 `${PYTHONHOME}` 变量覆盖初始值
- `sys.exec_prefix` # Python 扩展模块(第三方库)目录前缀。可以用过 `${PYTHONHOME}` 变量覆盖初始值
- `sys.executable` # Python 解释器的路径。

### 推导 sys.prefix 变量

假如 Python 解释器的路径是 `/usr/bin/python3`，生成的值如下：

> 这里用的是 Linux 上的默认 Python，`${sys.prefix}` 就是 `/usr`。

```python
>>> import sys
>>> sys.prefix
'/usr'
>>> sys.exec_prefix
'/usr'
>>> sys.executable
'/usr/bin/python3'
```

Windows 的话则直接实在 ${sys.prefix}/ 目录下查找其他路径

```python
>>> import sys
>>> sys.prefix
'C:\\Users\\DesistDaydream\\AppData\\Local\\Programs\\Python\\Python310'
>>> sys.executable
'C:\\Users\\DesistDaydream\\AppData\\Local\\Programs\\Python\\Python310\\python.exe'
```

推导出 `sys.prefix` 之后，**通常会在 **`**${sys.prefix}/lib/**`** 目录下查找其他的路径。**Windows 则是直接在 `${sys.prefix}/` 目录下查找其他路径

> 有的发行版则通常会在 ${sys.prefix}/lib64/ 目录下查找其他路径，比如 CentOS。

### 生成基本路径

首先，查找 `python${XY}.zip` 文件，通常为默认为 `${sys.prefix}/lib/python${XY}.zip`。注意：即使该文件不存在，通常也会添加默认值。该文件的大小版本号之间没有点。

然后，添加 `{sys.prefix}/lib/python${X.Y}`路径和 `${sys.prefix}/lib/python${X.Y}/lib-dynload` 路径。此时如果通过 `python -S` 选项关闭 site 模块，则会看到如下路径：

```python
python3 -S
Python 3.8.10 (default, Jun 22 2022, 20:18:18)
[GCC 9.4.0] on linux
>>> import sys
>>> sys.path
['', '/usr/lib/python38.zip', '/usr/lib/python3.8', '/usr/lib/python3.8/lib-dynload']
```

### 生成 site 模块路径（生成的路径乱七八糟，千奇百怪）

site 模块将会剩下的模块搜索路径添加到 sys.path 变量中。site 模块添加的路径也是最混乱的部分。不同的发行版，生成的路径也千奇百怪。

Python 默认的基本路径是 `${sys.prefix}/lib/python${X.Y}/site-packages/` 目录。

- 如果是 Ubuntu，则会将 site-packages 改为 dist-package，搞不懂为啥。。。o(╯□╰)o
- 对于 CentOS，则除了 lib，还会有一个 lib64，也是搞不懂为啥。。。o(╯□╰)o
- 真 TM 乱

除了基础路径，site 模块还会尝试导入 usercustomize 模块，以添加与用户相关的模块搜索路径。如果 site 模块中的 [ENABLE_USER_SITE](https://docs.python.org/3/library/site.html#site.ENABLE_USER_SITE) 变量为真，且 USER_SITE 定义的文件存在，则会将 USER_SITE 添加到 sys.path 中。

对于 usercustomize 模块，`sys.prefix` 不再使用，取而代之的是 `site.USER_BASE`，`site.USER_BASE` 的值通常为 `~/.local/`，生成的 `site.USER_SITE` 的值通常是 `site.USER_BASE/lib/python${X.Y}/site-packages`

```python
>>> site.USER_SITE
'/root/.local/lib/python3.6/site-packages'
>>> site.USER_BASE
'/root/.local'
```

下面的示例是 CentOS7 中的 Python

```bash
~]# python3 -m site
    '/root',
    '/usr/lib64/python36.zip',
    '/usr/lib64/python3.6',
    '/usr/lib64/python3.6/lib-dynload',
    '/usr/lib64/python3.6/site-packages',
    '/usr/lib/python3.6/site-packages',
]
USER_BASE: '/root/.local' (doesn't exist)
USER_SITE: '/root/.local/lib/python3.6/site-packages' (doesn't exist)
ENABLE_USER_SITE: True
~]# mkdir -p /root/.local/lib/python3.6/site-packages
~]# python3 -m site
sys.path = [
    '/root',
    '/usr/lib64/python36.zip',
    '/usr/lib64/python3.6',
    '/usr/lib64/python3.6/lib-dynload',
    '/root/.local/lib/python3.6/site-packages',
    '/usr/lib64/python3.6/site-packages',
    '/usr/lib/python3.6/site-packages',
]
USER_BASE: '/root/.local' (exists)
USER_SITE: '/root/.local/lib/python3.6/site-packages' (exists)
ENABLE_USER_SITE: True
```

我们可以使用 -s 选项禁用

```bash
~]# python3 -s -m site
sys.path = [
    '/root',
    '/usr/lib64/python36.zip',
    '/usr/lib64/python3.6',
    '/usr/lib64/python3.6/lib-dynload',
    '/usr/lib64/python3.6/site-packages',
    '/usr/lib/python3.6/site-packages',
]
USER_BASE: '/root/.local' (exists)
USER_SITE: '/root/.local/lib/python3.6/site-packages' (exists)
ENABLE_USER_SITE: False
```

### 总结

最后生成的 `sys.path` 具有类似如下的值：

```python
>>> import sys
>>> sys.path
['', '/usr/lib/python38.zip', '/usr/lib/python3.8', '/usr/lib/python3.8/lib-dynload', '/home/lichenhao/.local/lib/python3.8/site-packages', '/usr/local/lib/python3.8/dist-packages', '/usr/lib/python3/dist-packages']
```

到这里可以发现，**关于包路径搜索最重要的就是这个 **`**${sys.prefix}**`**路径前缀**，而这个值又是从使用的 Python 解释器路径推导出来的。所以要找到包的路径，只需要知道解释器的路径就可以了，如果遇到改变包的路径，只需要通过正确的 PATH 设置，指定你想要的 Python 解释器即可。

若 `sys.path` 中的所有目录都无法找到想要导入的模块，将会出现如下报错：

```bash
ModuleNotFoundError: No module named 'XXXXX'
```

## 安装 Python 模块

> 参考：
> - [官方文档，安装 Python 模块](https://docs.python.org/3.10/installing/index.html)
> - <https://frostming.com/2019/03-13/where-do-your-packages-go/>

管理 Python 的模块和包所在路径非常乱，不知道是何原因。

[pip](#fgjXo) 是首选的安装程序。从 Python 3.4 开始，它默认包含在 Python 二进制安装程序中。就算你是用 pipenv，poetry，底层依然是 pip，一律适用。

运行 pip 有两种方式：

- pip ...
- python -m pip ...

第一种方式和第二种方式大同小异，区别是第一种方式使用的 Python 解释器是写在 pip 文件的 shebang 里的，一般情况下，如果你的 pip 路径是 $path\_prefix/bin/pip，那么 Python 路径对应的就是 $path\_prefix/bin/python。如果你用的是 Unix 系统则 cat $(which pip) 第一行就包含了 Python 解释器的路径。第二种方式则显式地指定了 Python 的位置。这条规则，对于所有 Python 的可执行程序都是适用的。流程如下图所示。
![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/loffuc/1669286382022-472bf4de-24cf-4652-bc94-3d52d01f7df1.png)
那么，不加任何自定义配置时，使用 pip 安装包就会自动安装到 $path\_prefix/lib/pythonX.Y/site-packages 下（$path_prefix 是从上一段里得到的），可执行程序安装到 $path_prefix/bin 下，如果需要在命令行直接使用 my_cmd 运行，记得加到 PATH。

## requirements.txt 文件

> 参考：
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
`pip freeze >requirements.txt`

### 环境混用怎么办？

在导出依赖到 requirement.txt 文件时会有一种尴尬的情况。
你的本地环境不仅包含项目 A 所需要的依赖，也包含着项目 B 所需要的依赖。此时我们要如何做到只把项目 A 的依赖导出呢？
[pipreqs](https://github.com/bndr/pipreqs) 可以通过扫描项目目录，帮助我们仅生成当前项目的依赖清单。
通过以下命令安装：
`pip install pipreqs`
运行：
`pipreqs ./ --encoding utf8`

# pip

> 参考：
> - [GitHub 项目，pypa/pip](https://github.com/pypa/pip)
> - [官网](https://pip.pypa.io/)

pip 是 Python 的包管理程序。可以使用它来安装来自 Python 包索引和其他索引的包。从 Python 3.4 开始，它默认包含在 Python 二进制安装程序中。

## 关联文件与配置

**~/.pip/pip.conf **# 配置文件

pip 安装的模块我们可以从如下目录中找到，该目录下的目录名或文件名通常来说即是包名

- Windows
  - **%USERPROFILE%\AppData\Local\Programs\Python\Python${版本号}\Lib\site-packages\***
- Linux
  - root 用户：**/usr/local/lib/python${VERSION}/dist-packages/\***
  - 普通 用户：**~/.local/lib/python${PythonVersion}/site-packages/\***

有些包会生成一些可以执行程序，这些二进制文件则默认保存在如下目录

- Windows
  - **%USERPROFILE%\AppData\Local\Programs\Python\Python310\Scripts\***
- Linux
  - root 用户：**/usr/local/bin/\***
  - 普通 用户：**~/.local/bin/\***

## Syntax(语法)

### 应用示例

更新 pip：`pip install --upgrade pip`

对于 Python 开发用户来讲，PIP 安装软件包是家常便饭。但国外的源下载速度实在太慢，浪费时间。而且经常出现下载后安装出错问题。所以把 PIP 安装源替换成国内镜像，可以大幅提升下载速度，还可以提高安装成功率。

国内源：
新版 ubuntu 要求使用 https 源，要注意。

- 清华：<https://pypi.tuna.tsinghua.edu.cn/simple>
- 阿里云：<http://mirrors.aliyun.com/pypi/simple/>
- 中国科技大学 <https://pypi.mirrors.ustc.edu.cn/simple/>
- 华中理工大学：<http://pypi.hustunique.com/>
- 山东理工大学：<http://pypi.sdutlinux.org/>
- 豆瓣：<http://pypi.douban.com/simple/>

临时使用：
可以在使用 pip 的时候加参数 `-i https://pypi.tuna.tsinghua.edu.cn/simple`

例如：`pip install -i https://pypi.tuna.tsinghua.edu.cn/simple pyspider`，这样就会从清华这边的镜像去安装 pyspider 库。

永久修改，一劳永逸：
Linux 下，修改 ~/.pip/pip.conf (没有就创建一个文件夹及文件。文件夹要加“.”，表示是隐藏文件夹)，内容如下：

```properties
[global]
index-url = https://pypi.tuna.tsinghua.edu.cn/simple
[install]
trusted-host=mirrors.aliyun.com
```

windows 下，直接在 user 目录中创建一个 pip 目录，如：C:\Users\xx\pip，新建文件 pip.ini。内容同上。
