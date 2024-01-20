---
title: Control structure
linkTitle: Control structure
date: 2023-11-20T21:08
weight: 20
---

# 概述

> 参考：
>
> -

Python 语言提供了多种条件结构和分支结构用作 [Control structure(控制结构)](/docs/2.编程/计算机科学/Control%20structure.md)

# with

> 参考：
>
> - [官方文档，参考 - 8.复合语句 - 8.5 with 语句](https://docs.python.org/3/reference/compound_stmts.html#the-with-statement)

[`with`](https://docs.python.org/zh-cn/3/reference/compound_stmts.html#with) 语句用于包装带有使用上下文管理器 (参见 [with 语句上下文管理器](https://docs.python.org/zh-cn/3/reference/datamodel.html#context-managers) 一节) 定义的方法的代码块的执行。 这允许对普通的 [`try`](https://docs.python.org/zh-cn/3/reference/compound_stmts.html#try)...[`except`](https://docs.python.org/zh-cn/3/reference/compound_stmts.html#except)...[`finally`](https://docs.python.org/zh-cn/3/reference/compound_stmts.html#finally) 使用模式进行封装以方便地重用。语法如下：

```python
with EXPRESSION as TARGET:
    SUITE
```

EXPRESSION 的返回值赋值给 TARGET 变量，在 SUITE 中可以处理 TARGET。

用[白话解释](https://www.runoob.com/python3/python-with.html)：Python 中的 with 语句用于异常处理，封装了 try…except…finally 编码范式，提高了易用性。**with** 语句使代码更清晰、更具可读性， 它简化了文件流等公共资源的管理。在处理文件时使用 with 关键字是一种很好的做法。

```python
file = open('./test_runoob.txt', 'w')
file.write('hello world !')
file.close()
# 使用 try 的话
file = open('./test_runoob.txt', 'w')
try:
    file.write('hello world')
finally:
    file.close()
```

以上代码我们对可能发生异常的代码处进行 try 捕获，发生异常时执行 except 代码块，finally 代码块是无论什么情况都会执行，所以文件会被关闭，不会因为执行异常而占用资源。

```python
filename = "./testFile.txt"
with open(filename, 'r') as f:
    output = f.read()
```

使用 **with** 关键字系统会自动调用 f.close() 方法， with 的作用等效于 try/finally 语句是一样的。
