---
title: EBNF
---

# 概述

> 参考：
> - [Wiki,EBNF](https://en.wikipedia.org/wiki/Extended_Backus%E2%80%93Naur_form)
> - [Wiki,Metasyntax](https://en.wikipedia.org/wiki/Metasyntax)

**Extended Backus-Naur Form(扩展的 Backus-Naur 格式，简称 EBNF)** 是一组 [Metasyntax(元语法)](https://en.wikipedia.org/wiki/Metasyntax) 表示法。EBNF 用于对计算机[编程语言](https://en.wikipedia.org/wiki/Programming_language)等[形式语言](https://en.wikipedia.org/wiki/Formal_language)进行形式化描述。EBNF 是基于 BNF 的扩展。

EBNF 是一种表达形式语言语法的代码。EBNF 由两部分组成

- Terminal Symbols(终结符号)
- non-terminal production rules(非终结表达式规则) # 其实就相当于一个表达式

这两部分组合起来，其实就是一句话，最后跟一个句号~~~一行内容就是一个 EBNF 表示法，比如：

    digit excluding zero = "1" | "2" | "3" | "4" | "5" | "6" | "7" | "8" | "9" ;
    digit                = "0" | digit excluding zero ;

# Symbols(符号)

下面定义的符号意义中，`...` 仅仅用来表示符号中可以是任意内容，不属于被定义的符号的一部分。
`=` # Definition(定义)
`,` # Concatenation(串接)
`;` # Termination(终止)
`|` # Alternation(交替)，就是“或者”的意思。
`[...]` # Optional(可选)
`{}` # Repetition(重复)
`(...)` # Grouping(分组)
`'...'` # Terminal String(终端字符串)
`"..."` # Terminal String(终端字符串)
`(*...*)` # Comment(注释)
`?...?` # Special Sequence(特殊序列)
