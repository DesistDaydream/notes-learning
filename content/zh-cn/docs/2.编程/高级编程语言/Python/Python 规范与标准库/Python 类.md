---
title: Python 类
---

# 概述

> 参考：
> - [官方文档，教程-9.类](https://docs.python.org/3/tutorial/classes.html)
> - [廖雪峰 Python 教程，面向对象编程](https://www.liaoxuefeng.com/wiki/1016959663602400/1017495723838528)

Classes(类) 提供了一种将数据与功能捆绑到一起的手段。创建一个新的 **class(类)** 就意味着创造了一个新的 **object 的类型**，进而可以使用这个新的类型创建多个 **instances(实例)**。每个类实例都可以添加 **attributes(属性)** 以维护其自身的状态，同时还可以有 **methods(方法)** 用于修改其状态(方法在类中定义)。

与其他编程语言相比，Python 的类机制增加了包含最少新语法和语义的类。它是 C ++和 Modula-3 中发现的类机制的混合物。 Python 类提供面向对象编程的所有标准功能：类继承机制允许多个基类，派生类可以覆盖其基类或类的任何方法，方法可以调用具有相同名称的基类的方法。对象可以包含任意数量和类型的数据。正如模块所面临的那样，类 Python 的动态性质的课程：它们是在运行时创建的，并且可以在创建后进一步修改。

# struct 格式的 class

从 3.7 版本开始，可以使用 dataclass 装饰器让 class 声明中不再写 `__init__` 方法，就像这样

```python
from dataclasses import dataclass

@dataclass
class Employee:
    name: str
    dept: str
    salary: int

john = Employee('john', 'computer lab', 1000)
print(john.dept)
print(john.salary)
```

运行结果：

```shell
'computer lab'
1000
```
