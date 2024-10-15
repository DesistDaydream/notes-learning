---
title: Data type
linkTitle: Data type
weight: 20
---

# 概述

> 参考：
>
> - [官方文档，参考 - 数据模型](https://docs.python.org/3/reference/datamodel.html)

**object(对象)** 是 Python 中对数据的抽象。Python 程序中的所有数据都是由对象或对象间关系来表示的。 （从某种意义上说，按照冯·诺依曼的“存储程序计算机”模型，代码本身也是由对象来表示的。）

> Python 中的 object 概念类似于 JS 中的 [object](/docs/2.编程/高级编程语言/ECMAScript/JavaScript%20规范与标准库/object) 概念。但是又不完全一样。
>
> Python 变量的本质是对对象的引用

每个 object 都有一个 **Identity(标识符)**、**Type(类型)**、**Value(值)**。一个对象被创建后，它的 Identity 就绝不会改变；你可以将其理解为该对象在内存中的地址。 '[`is`](https://docs.python.org/zh-cn/3/reference/expressions.html#is)' 运算符可以比较两个对象的标识号是否相同；[`id()`](https://docs.python.org/zh-cn/3/library/functions.html#id "id") 函数能返回一个代表其标识号的整数。

**CPython 实现细节：** 在 CPython 中，`id(x)` 就是存放 `x` 的内存的地址。从某个角度来看，获取变量的值，就是获取变量所引用的对象的值。

> 对于 Python 和 JS 中的 object 来说，这个 object 就像全能的超人一样。。。。o(╯□╰)o。。。而 Go 语言中的全能超人则是 struct

对象的 Type 决定该对象所支持的操作 (例如 "对象是否有长度属性？" 比如数组类型的长度) 并且定义了该类型的对象可能的取值。[`type()`](https://docs.python.org/zh-cn/3/library/functions.html#type "type") 函数能返回一个对象的类型 (类型本身也是对象)。与 Identity 一样，一个对象的 Type 也是不可改变的。

下面的代码可以让我们对 Python 对象有更形象的感受：

```python
s = "Hello, World!"
# 变量的本质是对象的引用
# 在 Python 中有一个与 ESMAScript 中类似的 object(对象) 概念。
# 在 Python 中，所有的数据都是对象，每一个对象都有唯一的标识符、类型和值。与 JavaScript 不同的是，在 Python 中，变量本身并不拥有内存空间，它只是指向一个对象的引用。因此，我们在 Python 中声明变量时，并不需要显式地指定它的类型。
print("对象标识符: ", id(s))
print("对象的类型: ", type(s))
print("对象的值: ", s)
# 由于变量就是对对象的引用，那么就可以调用这个对象的属性和方法。例如：
print(s.upper())  # 输出 "HELLO, WORLD!"
print(s.lower())  # 输出 "hello, world!"
print(s.capitalize())  # 输出 "Hello, world!"
```

# 特殊方法名

> 参考：
>
> - [官方文档，参考 - 数据模型 - 特殊方法名](https://docs.python.org/3.12/reference/datamodel.html#special-method-names)

`__init__`、等

`__call__` # 当对象作为一个函数被调用时，执行 `__call__` 方法中的逻辑。

https://docs.python.org/3.12/reference/datamodel.html#emulating-callable-objects

比如 YOLO 的 `ultralytics.engine.model.Model` 类，若调用本身 Model()，则会触发其 `__call__` 方法，该方法直接返回 `self.predict()`。这里其实相当于让 Model 自身作为函数时默认调用 predict() 方法。


# 类型提示

痛点：Python 是动态类型语言，可以在运行时修改变量的类型，若不为函数的参数、变量指定类型，阅读代码会造成障碍，IDE 也无法给出正确的提示。

所以，从 Python 3.5 版本开始，Python 添加了 [typing](https://docs.python.org/3/library/typing.html) 库以支持 **Type hints(类型提示)**。

```python
# greeting 函数，参数 name 的类型应是 str，返回类型是 str。子类型也可以作为参数。
def greeting(name: str) -> str:
    return 'Hello ' + name

# 变量 stringType 的类型是 str
stringType: str = "Hello World!"

greeting(stringType)
```

<font color="#ff0000">注意：虽然我们可以在为变量、函数添加类型，但 Python 依然是动态语言，类型提示的功能就如其名，仅仅作为提示，哪怕函数参数的形参和实参类型不一样，也只是 IDE 会有错误提示，但是程序还是可以正常运行的。比如：</font>

```python
stringType: int = "Hello World!"
print("字符串类型: ", type(stringType))
```

输出结果为：`字符串类型:  <class 'str'>`
