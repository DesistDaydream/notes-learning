---
title: Data type
linkTitle: Data type
weight: 20
date: 2023-11-20T21:32
---

# 概述
>
> 参考：
>
> - [MDN 官方文档，JavaScript-JavaScript 数据类型和数据结构](https://developer.mozilla.org/en-US/docs/Web/JavaScript/Data_structures)
> - [MDN 官方文档，Javascript 标准内置对象](https://developer.mozilla.org/zh-CN/docs/Web/JavaScript/Reference/Global_Objects)(所有类型的对象的列表)
> - [网道，JavaScript 教程-面向对象编程-实例对象与 new 命令](https://wangdoc.com/javascript/oop/new.html)
> - <https://www.bilibili.com/video/BV1W54y1J7Ed?p=50>

**Primitive Type(原始类型，有的地方也称为基本数据类型)**

- **Number(数值)** # 十进制数字、科学计数法、其他进程表示方式的数字
- **String(字符串)** # 单引号或双引号内的一切内容
- **Boolean(布尔)** # ture 和 false
- **Null(空)** #
  - Undefined
  - Null # Null 类型是 Object，这是由于历史原因造成的。1995 年的 JavaScript 语言第一版，只设计了五种数据类型（对象、整数、浮点数、字符串和布尔值），没考虑 null，只把它当作 object 的一种特殊值。后来 null 独立出来，作为一种单独的数据类型，为了兼容以前的代码，typeof null 返回 object 就没法改变了。

**Complex Type(合成类型，有的地方也称为引用数据类型)**

- **object(对象)** # 各种值组成的集合，也就是下文提到的 [标准内置对象](#p3TIB)。在很多场景下，第一个 O 是小写的。
  - object 又划分为很多子类型：
    - **Ojbect(对象)** # 与 字典、映射 等同义的那个对象。
    - **Array(数组)** #
    - **Functiom(函数)** # JavaScript 中将 Function 当做一种类型来处理
    - **其他** #
  - 通常我们这么描述： Object 类型的 object、Array 类型的 object、String 类型的 object、等等。简化一点就是 Object 对象、Array 对象、String 对象、等等。

## 数据类型检测

