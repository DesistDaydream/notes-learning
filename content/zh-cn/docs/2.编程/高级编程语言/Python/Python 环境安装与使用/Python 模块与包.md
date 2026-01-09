---
title: Python 模块与包
weight: 2
---

# 概述

> 参考：
>
> - [官方文档，教程 - 6.模块](https://docs.python.org/3/tutorial/modules.html)
> - [GitHub 项目，pypa/packaging.python.org](https://github.com/pypa/packaging.python.org)
>   - [Python 包管理指南](https://packaging.python.org) “Python Packaging User Guide”(PyPUG) 旨在成为有关如何使用当前工具在 Python 中打包和安装发行版的权威资源。
>   - [Python Packaging Authority](https://www.pypa.io/en/latest/)
> - [廖雪峰 Python 教程，模块](https://www.liaoxuefeng.com/wiki/1016959663602400/1017454145014176)

Python 的 [Modular programming](/docs/2.编程/Programming%20technology/Modular%20programming.md) 设计

在 Python 中，**一个 `.py` 文件**就称之为一个 **Module(模块)**。

假如现在有 a.py 和 b.py，b.py 中有定义了名为 bFun 的函数，现在想要在 a.py 中使用 bFun 函数，则只需要在开始使用 `from b import bFun` 即可。

**使用模块有什么好处？**

最大的好处是大大提高了代码的可维护性。其次，编写代码不必从零开始。当一个模块编写完毕，就可以被其他地方引用。我们在编写程序的时候，也经常引用其他模块，包括 Python 内置的模块和来自第三方的模块。

使用模块还可以避免函数名和变量名冲突。相同名字的函数和变量完全可以分别存在不同的模块中，因此，我们自己在编写模块时，不必考虑名字会与其他模块冲突。但是也要注意，尽量不要与内置函数名字冲突。[这里](http://docs.python.org/3/library/functions.html)可以查看 Python 的所有内置函数。

**如果不同的人编写的模块名相同怎么办？**

为了避免模块名冲突，Python 又引入了按目录来组织模块的方法，称为 **Package(包)**。

举个例子，一个 `abc.py` 的文件就是一个名字叫 `abc` 的模块，一个 `xyz.py` 的文件就是一个名字叫`xyz`的模块。

现在，假设我们的`abc`和`xyz`这两个模块名字与其他模块冲突了，于是我们可以通过包来组织模块，避免冲突。方法是选择一个顶层包名，比如`mycompany`，按照如下目录存放：

```bash
mycompany
├─ __init__.py
├─ abc.py
└─ xyz.py
```

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

文件 `www.py` 的模块名就是 `mycompany.web.www`，两个文件 `utils.py` 的模块名分别是 `mycompany.utils` 和 `mycompany.web.utils`。

自己创建模块时要注意命名，不能和 Python 自带的模块名称冲突。例如，系统自带了 sys 模块，自己的模块就不可命名为 sys.py，否则将无法导入系统自带的 sys 模块。

`mycompany.web`也是一个模块，该模块对应的 .py 文件是 `__init__.py`。

## Package(包)

包是一种通过使用“带点的模块名称”来构建 Python 模块命名空间的方式。例如，模块名称 `A.B` 指定了一个名为 `A` 的包中的一个名为 `B` 的子模块。就像使用模块使不同模块的作者不必担心彼此的全局变量名称一样，使用点分模块名称可以节省作者多模块包，如 NumPy 或 Pillow 包，不必担心彼此的模块名称相同会影响到对方。

假设您要设计一组模块集（“包”）来统一处理声音文件和声音数据。有许多不同的声音文件格式（通常通过它们的扩展名识别，例如：.wav、.aiff、.au），因此您可能需要创建和维护不断增长的模块集合，以便在各种文件格式之间进行转换。您可能还想对声音数据执行许多不同的操作（例如混合、添加回声、应用均衡器功能、创建人工立体声效果），因此您将编写一个永无止境的模块流来执行这些操作。这是您的包的可能结构（以分层文件系统表示）：

```bash
sound/                          顶级包
      __init__.py               初始化名为 sound 的包，有这个文件的存在，Python 才可以将 sound/ 目录识别为包
      formats/                  用于格式转换的 Subpackage(子包，子包也可以称为 Module(模块))
              __init__.py
              wavread.py
              wavwrite.py
              aiffread.py
              aiffwrite.py
              auread.py
              auwrite.py
              ...
      effects/                  用于声音效果的 Subpackage(子包)
              __init__.py
              echo.py
              surround.py
              reverse.py
              ...
      filters/                  用于过滤的 Subpackage
              __init__.py
              equalizer.py
              vocoder.py
              karaoke.py
              ...
```

导入包时，Python 会在 `${sys.path}` 上的目录中搜索包子目录。

需要 `__init__.py` 文件才能使 Python 将包含该文件的目录视为包。这可以防止具有通用名称（例如字符串）的目录无意中隐藏了稍后出现在模块搜索路径上的有效模块。在最简单的情况下，`__init__.py` 可以只是一个空文件，但它也可以执行包的初始化代码或设置 `__all__` 变量，稍后将介绍。

注意：

- import 指令导入模块时是否成功，取决于执行 Python 代码时的所在位置。反正设计的挺恶心的。。。o(╯□╰)o
- 说白了：Package 就是 Module，只不过是对 Module 进行了分类。。。导入模块时，模块名称中有 `.` 就是包了。。其实还是模块。也可以这么理解：一个目录就是一个包，一个文件就是一个模块。
- 若没有 `__init__.py` 文件，我们使用 IDE 追踪代码时，如果想要从 import 指令中的导入的包进行追踪，会提示找不到该文件，无法追踪到。

## 总结

模块是一组 Python 代码的集合，可以使用其他模块，也可以被其他模块使用。

创建自己的模块时，要注意：

- 模块名要遵循 Python 变量命名规范，不要使用中文、特殊字符；
- 模块名不要和系统模块名冲突，最好先查看系统是否已存在该模块，检查方法是在 Python 交互环境执行`import abc`，若成功则说明系统存在此模块。

不管是 Package 还是 Module，都可以统一称为我们常说的 **Library(库)**，毕竟在下面的模块管理章节，也能看到 Python 保存**包和模块的目录通常都是 `lib/` 目录**。并且，Python 官方也有一个专门的页面列出了所有的[标准库](https://docs.python.org/3/library/index.html)

**Python Package Index(简称 PyPI)** 是 Python 编程语言官方的的软件存储库。<https://pypi.org/>

# Python 模块使用

Python 本身就内置了很多非常有用的模块，只要安装完毕，这些模块就可以立刻使用。

我们以内建的 `sys` 模块为例，编写一个 `hello` 的模块：

```python
#!/usr/bin/env python3
# -*- coding: utf-8 -*-

' a test module '

__author__ = 'DesistDaydream'

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

第 1 行和第 2 行是标准注释，第 1 行注释可以让这个`hello.py`文件直接在 Unix/Linux/Mac 上运行，第 2 行注释表示 .py 文件本身使用标准 UTF-8 编码；

第 4 行是一个字符串，表示模块的文档注释，任何模块代码的第一个字符串都被视为模块的文档注释；

第 6 行使用`__author__`变量把作者写进去，这样当你公开源代码后别人就可以瞻仰你的大名；

以上就是 Python 模块的标准文件模板，当然也可以全部删掉不写，但是，按标准办事肯定没错。

后面开始就是真正的代码部分。

你可能注意到了，使用`sys`模块的第一步，就是导入该模块：

```python
import sys
```

导入 `sys` 模块后，我们就有了变量 `sys` 指向该模块，利用 `sys` 这个变量，就可以访问 `sys` 模块的所有功能。

`sys` 模块有一个 `argv` 变量，用 list 存储了命令行的所有参数。`argv` 至少有一个元素，因为第一个参数永远是该 .py 文件的名称，例如：

运行`python3 hello.py`获得的`sys.argv`就是`['hello.py']`；

运行`python3 hello.py Michael`获得的`sys.argv`就是`['hello.py', 'Michael']`。

最后，注意到这两行代码：

```python
if __name__ == '__main__':
    test()
```

当我们在命令行直接使用 `python3 hello.py` 运行模块文件时，Python 解释器把 `__name__` 这个特殊变量的值设为 `__main__`，而如果在其他地方导入`hello`模块时 `__name__` 变量的值则是本模块的名称(i.e. hello)，因此，这种 `if` 测试可以让一个模块通过命令行运行时执行一些额外的代码，最常见的就是运行单元测试。

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
Python 3.10.6 (main, Nov 14 2022, 16:10:14)
>>> import hello
>>>
```

导入时，没有打印`Hello, word!`，因为没有执行`test()`函数。

调用`hello.test()`时，才能打印出`Hello, word!`：

```python
>>> hello.test()
Hello, world!
```

## \_\_name__ 变量说明

一个 .py 文件可以作为一个单独的脚本运行，也可以作为模块被其他代码引用，这个 `__name__` 变量，就是用来判断当前是哪一种情况的。

常用来在多文件的项目中，测试独立文件中的代码。

## Compiled(已编译) 的 Python 文件

为了加快加载模块的速度，Python 将每个模块的编译版本缓存在名为 module.version.pyc 的 **pycache** 目录下，其中版本编码了编译文件的格式；它通常包含 Python 版本号。例如，在 CPython 3.3 版中，spam.py 的编译版本将被缓存为 **pycache**/spam.cpython-33.pyc。这种命名约定允许来自不同版本和不同 Python 版本的编译模块共存。

Python 会根据编译后的版本检查源代码的修改日期，以查看它是否已过时并需要重新编译。这是一个完全自动的过程。此外，编译后的模块与平台无关，因此可以在不同架构的系统之间共享同一个库。

Python 在两种情况下不检查缓存。首先，它总是重新编译并且不存储直接从命令行加载的模块的结果。其次，如果没有源模块，它不会检查缓存。要支持非源（仅编译）分发，编译模块必须在源目录中，并且不能有源模块。

## 作用域

在一个模块中，我们可能会定义很多函数和变量，但有的函数和变量我们希望给别人使用，有的函数和变量我们希望仅仅在模块内部使用。在 Python 中，是通过`_`前缀来实现的。

正常的函数和变量名是公开的（public），可以被直接引用，比如：`abc`，`x123`，`PI`等；

类似`__xxx__`这样的变量是特殊变量，可以被直接引用，但是有特殊用途，比如上面的 `__author__`，`__name__` 就是特殊变量，`hello`模块定义的文档注释也可以用特殊变量`__doc__`访问，我们自己的变量一般不要用这种变量名；

类似 `_xxx` 和 `__xxx` 这样的函数或变量就是非公开的（private），不应该被直接引用，比如 `_abc`，`__abc` 等；

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
>
> - [官方文档，Python 教程 - 6.模块 - 模块搜索路径](https://docs.python.org/3/tutorial/modules.html#the-module-search-path)
> - [官方文档，Python 的安装和使用 - 命令行工具和环境](https://docs.python.org/3/using/cmdline.html)
>   - [PYTHONPATH](https://docs.python.org/3/using/cmdline.html#envvar-PYTHONPATH)
>   - [PYTHONHOME](https://docs.python.org/3/using/cmdline.html#envvar-PYTHONHOME)
> - https://stackoverflow.com/questions/59104100/what-is-the-idea-behind-the-installation-dependent-default-directory-layout
> - https://stackoverflow.com/questions/897792/where-is-pythons-sys-path-initialized-from

我们通过 **Python 模块的搜索路径**来管理 Python 模块，或者称为管理 Python 包。Python 模块的搜索路径在 **Python 解释器(i.e.python 可执行文件)启动时初始化**，并将路径字符串保存在 **`${sys.path}`** 这个数组类型的变量中。

Python 模块通常分两大类

- **内置模块** # 可以当做 [Python 标准库](/docs/2.编程/高级编程语言/Python/Python%20规范与标准库/Python%20规范与标准库.md#Python%20标准库)。
- **第三方模块** # 一般由 [PIP](/docs/2.编程/高级编程语言/Python/Python工具/PIP.md) 程序管理

> tips: 一般情况，模块默认都保存在 lib/ 目录下，这目录包含了了标准库的模块和第三方模块。而第三方模块是保存在 lib/site-pacages/ 目录下。至于目录的前缀，在不同环境一般各不相同。这个说法仅供参考并没有官方文档说明，仅是个人总结。

当导入一个名为 spam 的模块时，Python 解释器首先搜索具有该名称的 **Built-in module(内置模块)**(内置模块可以用过 `sys.builtin_module_names` 获取)，若没找到，则会在 **sys 内置模块中的 `${path}` 数组变量**下的目录列表中搜索名为 `spam.py` 的文件（从第一个元素开始逐一搜索，找到后就不再找了）。

> 这里面说的内置模块，属于[Python 规范与标准库](/docs/2.编程/高级编程语言/Python/Python%20规范与标准库/Python%20规范与标准库.md) 的一部分。这部分内置模块内嵌到解释器里面（也就是说无法在文件系统中找到与模块名相同的同名文件），它们给一些虽并非语言核心但却内嵌的操作提供接口，要么是为了效率，要么是给操作系统基础操作例如系统调入提供接口。 这些模块集是一个配置选项， 并且还依赖于底层的操作系统。 例如，[`winreg`](https://docs.python.org/zh-cn/3/library/winreg.html#module-winreg "winreg: Routines and objects for manipulating the Windows registry. (Windows)") 模块只在 Windows 系统上提供。一个特别值得注意的模块 [`sys`](https://docs.python.org/zh-cn/3/library/sys.html#module-sys "sys: Access system-specific parameters and functions.")，它被内嵌到每一个 Python 编译器中，**sys 模块是 CPython 非常重要的内置模块，也是很多功能的基础模块**。

**`${sys.path}` 变量是我们使用 Python 模块的最重要一环**。通常来说，该变量的值来源于以下位置：

- **运行 Python 代码文件所在的绝对路径**
- **${PYTHONPATH} 环境变量指定的路径**
  - 这是一个手动指定的目录列表，类似于类 Unix 中的 `$PATH` 变量，可以通过 `os.path` 获取其值。
  - 可以使用 `os.path.append()` 为 `$PYTHONPATH` 变量添加新的目录条目以便导入想要的模块。也可以直接设置 Linux 系统中的 `$PYTHONPATH` 变量。当项目大，需要对文件进行分类时，非常有用。
- **编译、安装 Python 时设置的默认路径**
  - 官方文档的这个说法挺模糊的，详见下面 sys.path 列表生成逻辑中的详解。这些路径主要取决于 [prefix](#1.确认并生成%20prefix) 的设置。
  - 按照惯例，通常包含：
    - Python 标准库保存路径。
    - [site 模块](https://docs.python.org/3/library/site.html#module-site)处理的第三方库保存路径。使用各种方式安装(比如 pip)的第三方库通常来说会存在 site-packages 目录中。

上述三种路径在 Python 启动时被初始化。我们可以通过 Python 中的 `${sys.path}` 数组变量查看这些路径。

**注意：各个平台生成的 sys.path 的列表并不统一**，具体原因详见 [模块管理混乱说明](#模块管理混乱说明)

## sys.path 列表生成逻辑

> 参考：
>
> - [官方文档，Python 标准库 - 导入模块 - sys.path 模块搜索路径初始化](https://docs.python.org/3/library/sys_path_init.html)
> - [官方文档，Python 标准库 - Python 运行时服务 - site—特定于 site 的配置](https://docs.python.org/3/library/site.html)
> - [源码，Moduels/getpath.py](https://github.com/python/cpython/blob/3.11/Modules/getpath.py)

如果说 `${sys.path}` 是我们使用 Python 模块最重要的东西，那 **prefix(路径的前缀)** 就是 `${sys.path}` 这个变量最重要的东西。prefix 通常表示 `${sys.prefix}` 与 `${sys.exec_prefix}` 这两个变量。**prefix** 的值来源于 **Python 解释器自身**。

我们可以在 CPyhon 源码 [Modules/getpath.py](https://github.com/python/cpython/blob/3.11/Modules/getpath.py) 查看在 Python 解释器启后生成 sys.path 的整体逻辑。这段代码会利用 [configure](https://github.com/python/cpython/blob/3.11/configure#L571) 与 [Include/cpython/initconfig.h](https://github.com/python/cpython/blob/3.11/Include/cpython/initconfig.h#L188)。经过如下几个步骤：

- PLATFORM CONSTANTS
- HELPER FUNCTIONS (note that we prefer C functions for performance)
- READ VARIABLES FROM config
- CALCULATE program_name
- **CALCULATE executable** # 计算 executable 变量的值
- CALCULATE (default) home
- **READ pyvenv.cfg** # 读取 pyvenv.cfg 文件，若读取到则说明这是一个 Python 虚拟环境
- CALCULATE base_executable, real_executable AND executable_dir
- DETECT _pth FILE
- CHECK FOR BUILD DIRECTORY
- **CALCULATE prefix AND exec_prefix** # 计算 prefix 和 exec_prefix 变量的值。
- **UPDATE pythonpath (sys.path)** #
- **POSIX prefix/exec_prefix QUIRKS** #
- **SET pythonpath FROM _PTH FILE** #
- **UPDATE config FROM CALCULATED VALUES** # 将计算后的值赋值给相关变量

最终，Python 解释器启动后，将为如下几个变量赋值

- **`sys.executable`** # Python 解释器的路径。
- **`sys.prefix`** # Python 标准模块(标准库)目录前缀。默认通过运行的 python 解释器生成出来。可以用过 `${PYTHONHOME}` 变量覆盖初始值
- **`sys.exec_prefix`** # Python 扩展模块目录前缀。可以用过 `${PYTHONHOME}` 变量覆盖初始值
- 等等......
- **`sys.path`**

下面我笔记中关于 sys.path 的生成逻辑中，只做简单描述，详细生成逻辑就像[这里](https://stackoverflow.com/questions/897792/where-is-pythons-sys-path-initialized-from)说的，跟论文一样，而且每个平台编译出来的 Python 解释器的设置也不一样，所以最后路径也不一样，由于 Python 复杂的自推导逻辑导致也不好总结普适性，所以，可以直接参考下文[sys.path 总结](#sys.path%20总结)中的各平台示例

### 1.确认并生成 prefix

> 这个行为准确描述应该是确认 prefix 的值之后，生成 `${sys.prefix}` 和 `${sys.exec_prefix}` 变量的值，这两个变量，通常用在标准库和第三方库的保存路径中，作为路径的前缀。

在我们运行 Python 解释器时，prefix 的值可以通过如下两种方式被设置：

- **一、若 `${PYTHONHOME}` 环境变量不为空，则使用该变量的值**
- **二、若 `${PYTHONHOME}` 环境变量为空，则使用在构建 Python 解释器时设定的默认值。**

**！！！重点来了！！！** 设置好 prefix 之后，需要通过 **Landmark(地标)** 文件确认 prefix 的值是否可用。Python 在 prefix 下找到某些 landmark，才说明 prefix 是可用的。因为这些 Landmark 所在的目录中包含了正常运行 Python 所需的条件（比如标准库等）。如果没找到这些，那么 Python 不应该正常启动，否则就算运行起来也会因为缺少各种依赖库导致无法使用。

下面将会对 prefix 的确认进行简单的分布概述。参考 Python 解释器运行后的路径生成[源码](https://github.com/python/cpython/blob/3.11/Modules/getpath.py#L176)

第一步，声明下列几个 landmark 的：（下面  `{VERSION_MAJOR}`  和 `{VERSION_MINOR}` 分别是 Python 的大版本号和小版本号，比如 Python3.10、Python310）

```python
if os_name == 'posix' or os_name == 'darwin':
    BUILD_LANDMARK = 'Modules/Setup.local'
    STDLIB_LANDMARKS = [f'{STDLIB_SUBDIR}/os.py', f'{STDLIB_SUBDIR}/os.pyc']
    PLATSTDLIB_LANDMARK = f'{platlibdir}/python{VERSION_MAJOR}.{VERSION_MINOR}/lib-dynload'
    BUILDSTDLIB_LANDMARKS = ['Lib/os.py']
    VENV_LANDMARK = 'pyvenv.cfg'
    ZIP_LANDMARK = f'{platlibdir}/python{VERSION_MAJOR}{VERSION_MINOR}.zip'

elif os_name == 'nt':
    BUILD_LANDMARK = f'{VPATH}\\Modules\\Setup.local'
    STDLIB_LANDMARKS = [f'{STDLIB_SUBDIR}\\os.py', f'{STDLIB_SUBDIR}\\os.pyc']
    PLATSTDLIB_LANDMARK = f'{platlibdir}'
    BUILDSTDLIB_LANDMARKS = ['Lib\\os.py']
    VENV_LANDMARK = 'pyvenv.cfg'
    ZIP_LANDMARK = f'python{VERSION_MAJOR}{VERSION_MINOR}{PYDEBUGEXT or ""}.zip'
```

第二步，将 prefix 与声明的 landmark 的路径拼接得出 landmark 的绝对路径。

- **os.py 文件** # 标准库 landmark，os.py 文件所在目录就是标准库目录。
  - Windows 中，默认在 `${prefix}/Lib/os.py`
  - Linux 中，默认在 `${prefix}/lib/python${VERSION_MAJOR}.${VERSION_MINOR}/os.py`
- **lib-dynload/ 目录** # 未知
  - Windows 中，没找到该目录
  - Ubuntu 中，默认在 `${prefix}/lib/python${VERSION_MAJOR}.${VERSION_MINOR}/lib-dynload/`
- **python${XY}.zip 文件** # ZIP landmark，zip 文件所在目录就是 prefix 目录
  - Windows 中，默认在`${prefix}/python${VERSION_MAJOR}${VERSION_MINOR}.zip`
  - Ubunut 中，默认在 `${prefix}/python${VERSION_MAJOR}${VERSION_MINOR}.zip`
- 等等......

第三步，Python 解释器将会逐一检查这些 landmark 的绝对路径。

- 这里有点要注意：若设置了 PYTHONHOME 的话，则先寻找由 PYTHONHOME 拼接的 landmark 绝对路径，若找不到 landmark 文件的话，还会去安装 Python 时设定的 prefix 默认值中再寻找一遍 landmark 文件。
- 若某些必须的 landmark 未找到的话，Python 解释器将会启动失败并报错。
- 假如我们将 os.py 文件移动走，或者设置一个没有 os.py 存在的 `${PYTHONHOME}`，那么 Python 解释器都是启动不起来的，假如现在设置 `export PYTHONHOME="/error_python_home"`，Python 将会报错：

```bash
~]# export PYTHONHOME="/error_python_home"
~]# python3
Python path configuration:
    ......
Fatal Python error: init_fs_encoding: failed to get the Python codec of the filesystem encoding
Python runtime state: core initialized
ModuleNotFoundError: No module named 'encodings'

Current thread 0x00007fe5a7de31c0 (most recent call first):
  <no Python frame>
```

第四步，这些 landmark 定位后，prefix 即表示可用，此时为 sys.prefix 和 sys.exec_prefix 两个变量赋值。

> 到这里就有很大的疑惑，为什么 Ubuntu、CentOS 的 Python 的可执行文件，即.Python 解释器在 /usr/bin/ 目录下，但是生成的 prefix 却是 /usr 呢？
>
> 在使用源码构建 Python 解释器之前，会执行 `./configure --prefix=/usr 命令，其中 --prefix 就是设置 prefix 的默认值。其实 CPython 官方在 configure 文件中为 prefix 设置了 /usr/local 这个默认值。。。但是。。。
>
> 各种 Linux 发行版的 Python 都在生成自己的路径，所以 Python 的 prefix 混乱根源也在这。这时 Python 的官方文档和各个发行版实际就对不上，给初学者造成很大困扰。。o(╯□╰)o。。
>
> 所以，就算 Python 解释器在 /usr/bin 目录下，只要 prefix 为 /usr，那么根据 prefix 拼接出的各种 sys.path 路径，也会在 /usr/lib/PythonX.Y 目录下，这从[Python 环境安装与使用](/docs/2.编程/高级编程语言/Python/Python%20环境安装与使用/Python%20环境安装与使用.md) 中的自定义 python 部分可以看到实际改变情况。
>
> 也许是出于复用 /usr/lib 目录的目的把，毕竟 python 属于系统级别依赖的程序，并且 Python 的模块大多是 ${prefxi}/lib/ 目录的，把 /usr 当做 prefix 的话，可以统一管理所有 lib。
>
> 而 Windows 中，Python 并不是系统自带的，可以手动安装，并且安装时可以指定安装目录，这时，这个安装目录就是 prefix。

#### prefix 结果示例

假如 Python 解释器的路径

- 在 Ubuntu 中是 `/usr/bin/python3`
- 在 Windows 中是 `D:\Tools\Python\Python310\python.exe`

Ubuntu 生成的值为：

```python
>>> sys.executable
'/usr/bin/python3'
>>> sys.prefix
'/usr'
>>> sys.exec_prefix
'/usr'
>>> sys.platlibdir
'lib'
```

Windows 生成的值为：

```python
>>> sys.executable
'D:\\Tools\\Python\\Python310\\python.exe'
>>> sys.prefix
'D:\\Tools\\Python\\Python310'
>>> sys.exec_prefix
'D:\\Tools\\Python\\Python310'
>>> sys.platlibdir
'DLLs'
```

### 2.添加要运行的 Python 文件所在的路径

添加运行 Python 代码文件所在的绝对路径，若直接运行的 Python 解释器，则 `$sys.path` 的第一个元素为空

- 在下面的 [Ubuntu 示例](#Ubuntu%20示例)中，`sys.path` 的第一个元素（`/root/scripts`）是 module-path-demo.py  文件所在路径，即执行的 Python 代码文件所在路径，如果不是运行的 Python 代码文件，则第一个元素为空。每当运行一个 Python 文件时，就相当于默认执行了 `sys.path.append("文件所在绝对路径")` 代码。

### 3.添加 PYTHONPATH 环境变量设置的路径

添加 `${PYTHONPATH}` 变量中的值

### 4.添加标准库的存放路径

添加包含 Python 标准模块以及这些模块所依赖的任何扩展模块的文件和目录，**这些路径是很重要的，包含了 Python 解释器启动成功所依赖的模块**。通常包含如下文件和目录：

- `${sys.prefix}/lib/python${VERSION_MAJOR}${VERSION_MINOR}.zip` # Python 库文件的归档文件，其中包含了许多 Python 标准库和已安装的第三方库的模块。注意：即使该文件不存在，通常也会添加默认值。该文件的大小版本号之间没有点。
- `${sys.prefix}/lib/python${VERSION_MAJOR}.${VERSION_MINOR}/` # 标准库保存路径
- `${sys.prefix}/lib/python${VERSION_MAJOR}.${VERSION_MINOR}/lib-dynload/` # 使用 C 语言编写的模块的存放路径。

这里的扩展模块是指用 C 或 C++ 编写的模块，使用 Python 的 C API 与核心和用户代码交互。并不是指 Python 的第三方模块

- Windows 上的扩展模块是后缀名为 `.pyd` 的文件
- Linux 上的扩展模块是后缀名为 `.so` 的文件

此时通过 `python -S` 命令在运行解释器时不自动加载 site 模块，则会看到如下路径：

```python
export PYTHONPATH="/pythonpath-demo"
~]# python3 -S
Python 3.10.6 (main, Nov 14 2022, 16:10:14)
>>> import sys
>>> sys.path
['', '/pythonpath-demo', '/usr/lib/python310.zip', '/usr/lib/python3.10', '/usr/lib/python3.10/lib-dynload']
```

其中前两个元素是基本路径，后面三个元素是 prefix 相关的路径。

### 5.添加第三方库的存放路径

Python 的第三方模块通常保存在由 site 模块生成的目录中，该目录名称通常为 site-packages。但是有的系统，比如 Ubuntu，会将 site-packages 目录的名称改为 dist-packages。site 模块在 Python 解释器启动时自动调用，并将生成目录添加到 sys.path 列表中。

**第一、调用 site 模块的 `main()` 函数将 `${sys.prefix}/lib/site-package/` 目录添加到 `sys.path` 变量中**。`site.main()` 函数从 Python3.3 版本开始被自动调用，除非运行 Python 解释器时添加 -S 标志。

Ubuntu 效果如下：

```python
~]# python3 -S
Python 3.10.6 (main, Nov 14 2022, 16:10:14)
>>> import sys,site
>>> sys.path
['', '/usr/lib/python310.zip', '/usr/lib/python3.10', '/usr/lib/python3.10/lib-dynload']
>>> site.main()
>>> sys.path
['/root', '/usr/lib/python310.zip', '/usr/lib/python3.10', '/usr/lib/python3.10/lib-dynload', '/usr/local/lib/python3.10/dist-packages', '/usr/lib/python3/dist-packages']
```

在这个示例中，我们可以看到 site-packages 目录跟官方文档的说明并不一样是吧？在 [模块管理混乱说明](#模块管理混乱说明) 中详细说明

**第二、site 模块在将全局的 site 目录添加前，会先添加用户 site 相关的模块搜索路径**。如果 site 模块中的 [ENABLE_USER_SITE](https://docs.python.org/3/library/site.html#site.ENABLE_USER_SITE) 变量为真，且 USER_SITE 定义的文件存在，则会将 USER_SITE 添加到 sys.path 中，用户 site 不再依赖 prefix，取而代之的是 `site.USER_BASE`，`site.USER_BASE` 的值通常为 `~/.local/`，生成的 `site.USER_SITE` 的值通常是 `site.USER_BASE/lib/python${X.Y}/site-packages`

```python
>>> site.USER_SITE
'/root/.local/lib/python3.6/site-packages'
>>> site.USER_BASE
'/root/.local'
```

```bash
~]# python3 -m site
sys.path = [
    '/root',
    '/usr/lib/python310.zip',
    '/usr/lib/python3.10',
    '/usr/lib/python3.10/lib-dynload',
    '/usr/local/lib/python3.10/dist-packages',
    '/usr/lib/python3/dist-packages',
]
USER_BASE: '/root/.local' (doesn't exist)
USER_SITE: '/root/.local/lib/python3.10/site-packages' (doesn't exist)
ENABLE_USER_SITE: True
~]# mkdir -p /root/.local/lib/python3.10/site-packages
~]# python3 -m site
sys.path = [
    '/root',
    '/usr/lib/python310.zip',
    '/usr/lib/python3.10',
    '/usr/lib/python3.10/lib-dynload',
    '/root/.local/lib/python3.10/site-packages',
    '/usr/local/lib/python3.10/dist-packages',
    '/usr/lib/python3/dist-packages',
]
USER_BASE: '/root/.local' (exists)
USER_SITE: '/root/.local/lib/python3.10/site-packages' (exists)
ENABLE_USER_SITE: True
```

**最后，尝试导入名为 usercustomize 与 sitecustomize 模块**。该模块由用户自行编写代码实现，以添加自己的路径。

## sys.path 总结

总的来说，与 prefix 相关的路径还是有一定规律的，通常都是在 `${prefix}/lib/` 目录下，官方文档通常都会省略 prefix，直接用 lib/、Lib/ 等作为开头来描述特定文件的位置。

> 注意：下面  `{VERSION_MAJOR}`  和 `{VERSION_MINOR}` 分别是 Python 的大版本号和小版本号，比如 Python3.10、Python310

- **`${PWD}`** # 当前工作目录
- ${PYTHONPATH} # 手动添加的目录。
- **`${sys.prefix}/lib/python${VERSION_MAJOR}${VERSION_MINOR}.zip`** # Python 库文件的归档文件，其中包含了许多 Python 标准库和已安装的第三方库的模块。
- **`${sys.prefix}/lib/python${VERSION_MAJOR}.${VERSION_MINOR}/`** # 标准库保存路径。没有这个的话 Python 解释器无法正常运行
- **`${sys.prefix}/lib/python${VERSION_MAJOR}.${VERSION_MINOR}/lib-dynload/`** # 使用 C 语言编写的模块的存放路径。
- **`${sys.prefix}/lib/python${VERSION_MAJOR}.${VERSION_MINOR}/sist-packages/`** # 第三方库保存路径。该目录在 Ubuntu 系统中名称为 dist-packages
- **`${site.USER_SITE}`** # 启动用户 site 后，保存第三方库的路径。

> 注意：从这里可以看到，不同 Python 版本的三方库路径不同，如果把 Python 从 3.8 升级到 3.9，那么之前装的三方库都没法用了。当然可以整个文件夹都拷贝过去，或者添加统一的 PYTHONPATH 路径大部分情况不会出问题。

最后生成的 `sys.path` 具有类似如下的值：

```python
>>> import sys
>>> sys.path
['', '/usr/lib/python310.zip', '/usr/lib/python3.10', '/usr/lib/python3.10/lib-dynload', '/usr/local/lib/python3.10/dist-packages', '/usr/lib/python3/dist-packages']
```

到这里可以发现，**关于包路径搜索最重要的就是这个 `${sys.prefix}` 路径前缀**，而这个值是由 Python 解释器生成出来的。

若 `sys.path` 中的所有目录都无法找到想要导入的模块，将会出现如下报错：

```bash
ModuleNotFoundError: No module named 'XXXXX'
```

### Ubuntu 示例

`${sys.prefix}/lib/python${X.Y}/site-packages/`

```bash
~]# cat module-path-demo.py
import sys
print(sys.path)
print(sys.prefix)
print(sys.exec_prefix)
print(sys.executable)
~]# python3 module-path-demo.py
['/root/scripts', '/usr/lib/python310.zip', '/usr/lib/python3.10', '/usr/lib/python3.10/lib-dynload', '/usr/local/lib/python3.10/dist-packages', '/usr/lib/python3/dist-packages']
/usr
/usr
/usr/bin/python3
```

### CentOS 示例

TODO

## 模块关联的可执行文件

部分模块具有可执行文件，可以当做命令使用，在不同平台上，会在不同位置生成可执行文件

Windows

- **${sys.prefix}/Scripts/**

Ubuntu

- **/usr/local/bin/**

# 模块管理混乱说明

在 [公众号-OSC开源社区，Flask之父凭一己之力击败各种GPT，称Python包管理比LLM更火热](https://mp.weixin.qq.com/s/i5fWKWs-D9ZphDqhtYoj9Q) 这篇文件中描述了 Flask 框架作者 Armin 画了一张图来描述他对 Python 包管理现状的感受，意思就是由于缺乏统一的标准，因此诞生了满足不同需求和场景的许多不同工具——不过每个抱着“统一”初心的标准最后都是适得其反。

Python 的模块管理非常混乱和复杂，不像 Go 只需要指定 GOPATH 变量之后，所有安装的依赖库都会存放到 GOPATH 目录下。上述三种路径如果说最像 GOPATH 的，那应该是安装 Python 是默认值中的 PYTHONPATH 或者 prefix 变量了。

并且，**各个发行版会修改 Python 的代码，这就导致编译后的很多默认值并不统一**，下面我举几个例子来说明

## 内置模块的差异

site 模块向 sys.path 添加的路径是非常混乱。不同的发行版，生成的路径也千奇百怪。而且不一定只生成一个 site-packages 目录

主要差异可以在 site.py 文件中关于 `def getsitepackages(prefixes=None)` 函数看到

3.10 版本 CPython 的原始 [site.py](https://github.com/python/cpython/blob/3.10/Lib/site.py)

```python
        if os.sep == '/':
            for libdir in libdirs:
                path = os.path.join(prefix, libdir,
                                    "python%d.%d" % sys.version_info[:2],
                                    "site-packages")
                sitepackages.append(path)
        else:
            sitepackages.append(prefix)

            for libdir in libdirs:
                path = os.path.join(prefix, libdir, "site-packages")
                sitepackages.append(path)
```

Ubuntu 中的 site.py

```python
        if os.sep == '/':
            if is_virtual_environment:
                sitepackages.append(os.path.join(prefix, "lib",
                                                 "python%d.%d" % sys.version_info[:2],
                                                 "site-packages"))
            sitepackages.append(os.path.join(prefix, "local/lib",
                                             "python%d.%d" % sys.version_info[:2],
                                             "dist-packages"))
            sitepackages.append(os.path.join(prefix, "lib",
                                             "python3",
                                             "dist-packages"))
            # this one is deprecated for Debian
            for libdir in libdirs:
                path = os.path.join(prefix, libdir,
                                    "python%d.%d" % sys.version_info[:2],
                                    "dist-packages")
                sitepackages.append(path)
        else:
            sitepackages.append(prefix)

            for libdir in libdirs:
                path = os.path.join(prefix, libdir, "site-packages")
                sitepackages.append(path)

```

Rocky 的 site.py

```python
        if os.sep == '/':
            sitepackages.append(os.path.join(prefix, "lib64",
                                        "python" + sys.version[:3],
                                        "site-packages"))
            sitepackages.append(os.path.join(prefix, "lib",
                                        "python%d.%d" % sys.version_info[:2],
                                        "site-packages"))
        else:
            sitepackages.append(prefix)
            sitepackages.append(os.path.join(prefix, "lib64", "site-packages"))
            sitepackages.append(os.path.join(prefix, "lib", "site-packages"))
```

Windows 的 site.py

```python
# 与 CPython 的原始 [site.py](https://github.com/python/cpython/blob/3.10/Lib/site.py) 代码保持一致
```

上面的例子中，Ubuntu 的 dist-packages 应该是 site-packages 才对，为啥这么改搞不懂为啥。。。o(╯□╰)o。。。对于 CentOS，则除了 lib，还会有一个 lib64，也是搞不懂为啥。。。o(╯□╰)o。。。

## 关于 Python 包的在线浏览网站

Python 好像没有像 [Go Package](https://pkg.go.dev/) 类似的网站

**Python Package Index(PyPI)** 是 Python 编程语言的软件存储库。但是内容非常少，只是提供了包的简介、版本、网站连接等信息。

问了问 NewBing 之后， [Read the Docs](https://readthedocs.org/search/?q=python+package) 这个网站可能还有点好内容，下面是完整的回答：

Python 有很多不同的网站可以查看各种包的信息，比如：

- [PyPI](https://pypi.org/)：Python 包索引，可以搜索、安装和发布 Python 软件包。
- [Python Package Index](https://packaging.python.org/en/latest/tutorials/packaging-projects/)：Python 打包用户指南，可以教您如何打包和分发 Python 项目。
- [Stack Overflow](https://stackoverflow.com/questions/tagged/python)：一个编程问答网站，可以找到关于 Python 包的问题和答案。
- [GitHub](https://github.com/search?q=python+package)：一个代码托管平台，可以浏览和下载 Python 包的源代码。

如果您想要像 go package 网站那样，可以直接在网页上看到 Python 包的类型、函数等文档，您可以试试以下几个网站：

- [Read the Docs](https://readthedocs.org/search/?q=python+package)：一个文档托管平台，可以查看很多 Python 包的在线文档。
- [Python Module of the Week](https://pymotw.com/3/)：一个网站，可以学习 Python 标准库中的模块的用法和示例。
- [Awesome Python](https://awesome-python.com/)：一个网站，可以发现一些优秀的 Python 包和资源。

\[1]: [How do I find the location of my Python site-packages directory?](https://stackoverflow.com/questions/122327/how-do-i-find-the-location-of-my-python-site-packages-directory/)
\[2]: [What is python's site-packages directory? - Stack Overflow](https://stackoverflow.com/questions/31384639/what-is-pythons-site-packages-directory)
\[3]: [go-python/gopy - GitHub](https://github.com/go-python/gopy)
\[4]: [PyPI · The Python Package Index](https://pypi.org/)
\[5]: [Packaging Python Projects — Python Packaging User Guide](https://packaging.python.org/en/latest/tutorials/packaging-projects/)
\[6]: [How To Package And Distribute Python Applications](https://www.digitalocean.com/community/tutorials/how-to-package-and-distribute-python-applications)

# 安装 Python 模块/包

详见 [Python 工具](/docs/2.编程/高级编程语言/Python/Python工具/Python%20工具.md#安装%20Python%20包/模块)
