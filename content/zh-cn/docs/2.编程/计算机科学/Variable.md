---
title: Variable
linkTitle: Variable
weight: 4
---

# 概述

> 参考：
>
> - [Wiki，Variable_(computer_science)](https://en.wikipedia.org/wiki/Variable_(computer_science))

在计算机编程中，**Variable(编程)** 是一个抽象的存储位置，与一个相关的符号名称配对，变量中包含一些称为 **Value(值)** 的已知或未知数量的信息。或者可以说，变量是一个有名字的容器，用于存放特定 [Data Type(数据类型)](/docs/2.编程/计算机科学/Data%20type/Data%20type.md) 的数据。

变量是一个可以改变内容的固定规定，比如我定义“这台电脑的名字”叫“XXX”，“这台电脑的名字”就是不变的内容，“XXX”就是可以改变的内容，给不变的内容定义不同的内容

- 比如 X=1，X=2，X=3 等等，X 就是不变的，1，2，3 等等都是可变的，X 就是一个变量，可以被赋予某些内容
- [环境变量](/docs/1.操作系统/Operating%20system/环境变量.md)就是在当前环境下所定义的内容，比如 linux 启动了一个 shell，在 shell 这个环境下，有一些规定被定义了，这些规定就是环境变量；不在这个 shell 下变量就不生效
  - 比如：LANG 是一个语言的规定，你赋予他一个内容，就相当于定义了这个 shell 环境下所显示的语言，比如 LANG=US，LANG=CN 等等。LANG 这叫定义语言，这是不变的，可变的是后面的自定义内容，语言(不变)=英语、汉语、日语(可变)。
