---
title: Lua
linkTitle: Lua
date: 2023-10-02T09:06
weight: 20
---

# 概述

> 参考：
> 
> - [GitHub 项目，lua/lua](https://github.com/lua/lua)
> - [官网](https://www.lua.org/about.html)
> - [Wiki，Lua](https://en.wikipedia.org/wiki/Lua_(programming_language))
> - [知乎，Lua 是怎样一门语言？](https://www.zhihu.com/question/19841006)

Lua 是一种强大、高效、轻量级、可嵌入的脚本语言。它支持过程式编程、面向对象编程、函数式编程、数据驱动编程和数据描述。

> tips: “Lua”（发音为 LOO-ah）在葡萄牙语中的意思是“月亮”。因此，它既不是首字母缩略词也不是缩写词，而是一个名词。更具体地说，“Lua”是一个名字，是地球月球的名字，也是语言的名字。与大多数名称一样，它应该以小写字母开头，即“Lua”。请不要将其写为“LUA”，这样既难看又令人困惑，因为这样它就变成了一个缩写词，对不同的人有不同的含义。所以，请把“Lua”写对！

## 学习资料

[菜鸟教程，Lua](https://www.runoob.com/lua/lua-tutorial.html)

安装完 [rjpcomputing/luaforwindows](https://github.com/rjpcomputing/luaforwindows) 编译器后，可以在 `安装目录/examples/quickluatour.lua` 文件中看到非常全面的 Lua 使用示例，直接运行就可体验基本的 Lua 语法。

# Hello World

```lua
print("Hello World!")
```

Lua可以在交互模式下输入代码直接查看效果，也可以将lua代码写入以.lua结尾的文件，然后使用`lua FILE.lua`命令查看代码效果，效果如下：

交互式：

```lua
~]# lua
Lua 5.1.4  Copyright (C) 1994-2008 Lua.org, PUC-Rio
> print("hello world!")
hello world!
```

脚本式：

```lua
~]# cat hello_world.lua
print("hello world!")
~]# lua hello_world.lua
hello world!
```

# Lua 环境安装与使用

**LuaJIT**(Just-In-Time) 是 Lua 编程语言的跟踪即时编译器。

window 下你可以使用一个叫 "SciTE" 的 IDE环 境来执行 lua 程序，下载地址为：

- Github 下载地址：[https://github.com/rjpcomputing/luaforwindows/releases](https://github.com/rjpcomputing/luaforwindows/releases)
