---
title: Generic
linkTitle: Generic
date: 2023-11-28T21:26
weight: 20
---

# 概述

> 参考：
>
> - [Wiki，Generic proramming](https://en.wikipedia.org/wiki/Generic_programming)(泛型编程)

**Generic(泛型)** 编程是一种计算机编程的风格，在这种风格中，程序中的数据类型是不固定的（i.e. 广泛的类型），若想让这种不固定的类型固定下来，需要为其提供 **类型参数**  以便固定成具体的类型。这种编程风格由 1973 年的 ML 编程语言首创，允许编写通用函数或类型，这些函数或类型在使用时仅在操作的类型集上有所不同，从而减少了重复代码。

Generic 也可以看做是以一种特殊的数据类型，就像 string、int 等等数据类型一样，除了具备 [Data type](docs/2.编程/计算机科学/Data%20Type%20AND%20Literal/Data%20type%20AND%20Literal.md) 的能力，还多了变化的能力。

假如现在有如下函数：

```text
func g<T any>(i T){
  print(typeof(i))
}
```

其中函数多了 `<T any>` 这部分，这就表明，这个函数中使用了泛型，这个广泛的类型标识符为 T（就像字符串的标识符为 string 之类的一样），可以让自己变成 any(任意) 的类型。

> 有的语言，比如 Go 语言，不是使用 `[ ]` 符号而不是 `< >` 来标识泛型语法。

此时我们调用 g 函数

```text
g<int>(1)
g<string>("1")
```

`g<int>(1)` 输出的结果为 int，`g<string>("string")` 输出的结果为 string。此时，调用是使用的 `< >` 可以理解为约束，就是让 g 函数中的 T 变为指定的类型。

比如使用 `g<int>(1)` 调用 g 函数时，这个函数相当于变成了 `func g(i int){ print(typeof(i)) }`。
