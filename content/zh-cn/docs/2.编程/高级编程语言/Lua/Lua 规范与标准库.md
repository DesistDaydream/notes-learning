---
title: Lua 规范与标准库
weight: 20
linkTitle: Lua 规范与标准库
date: 2023-10-16T09:07
---
# 概述

> 参考：

# Lua 语言关键字

> 参考：

Lua的关键字有下面几个，这些关键字不能作为常量或者变量或其他用户自定义标识符

1. and
2. break
3. do
4. else
5. elseif
6. end
7. false
8. for
9. function
10. if
11. in
12. local
13. nil
14. not
15. or
16. repeat
17. return
18. then
19. true
20. until
21. while
22. goto

# Lua 语言规范

> 参考：

一般来说，以下划线开头连接一串大写字母的名字（比如_VERSION）被保留用于Lua内部全局变量,效果如下：
```lua
~]# lua
Lua 5.1.4  Copyright (C) 1994-2008 Lua.org, PUC-Rio
> print(_VERSION)
Lua 5.1
```

# Lua 标准库

> 参考：
> 
> - [lua 参考手册，6 - 标准库](https://www.lua.org/manual/5.4/manual.html#6)

## 数据类型

- number
- string
- table
- function

## 变量

lua的变量有三种类型：

1. 全局变量。
2. 局部变量。
3. 表中的域。

默认情况下，lua定义的变量为全局变量。全局变量不需要声明，给一个变量赋值后即创建了这个变量，访问一个没有初始化的变量并不会报错，只不过得到的结果是`nil`,效果如下：

```lua
> print(b)
nil
> b=1
> print(b)
1
```

Note：如果想要删除一个变量，只需要将`nil`赋值给变量即可

也就是说，当且当一个变量不等于`nil`时，这个变量即存在

使用`local`关键字即可将变量声明为局部变量，效果如下:`local b = 5 `

## table

table 是 Lua 语言中的一种数据结构，可以通过 table 这种数据结构实现 数组、字典 等等常见的数据类型。同时，Lua 也是通过 table 来解决 模块、包、对象 等面向对象相关的问题。

- Lua 本身没有专门的数组类型，而是使用 table 的数据结构来实现数组的功能
- Lua 的面向对象功能也通过 table 来实现。

table 使用 `{}` 符号声明

```lua
-- 初始化表
mytable = {}
```

表可以用通过 `.` 和 `[]` 这两种符号调用 table 中的元素

```lua
address={} -- empty address
address.Street="Wyman Street"
address.StreetNumber=360
address.AptNumber="2a"
address.City="Watertown"
address.State="Vermont"
address.Country="USA"

print(address.StreetNumber, address["AptNumber"])



-------- Output ------

360     2a
```

## 面向对象

Lua 中最基本的结构是 table，所以需要用 table 这个数据结构来描述对象的属性。

- lua 中的 function 可以用来表示方法。那么LUA中的类可以通过 table + function 模拟出来。
- 至于继承，可以通过 metetable 模拟出来（不推荐用，只模拟最基本的对象大部分实现够用了）。

Lua 中的表不仅在某种意义上是一种对象。像对象一样，表也有状态（成员变量）；也有与对象的值独立的本性，特别是拥有两个不同值的对象（table）代表两个不同的对象；一个对象在不同的时候也可以有不同的值，但他始终是一个对象；与对象类似，表的生命周期与其由什么创建、在哪创建没有关系。表也可以有自己的 Method(方法)（用 C++ 的叫法就叫：成员函数）：

```lua
Account = {balance = 0}  
function Account.withdraw (v)  
    Account.balance = Account.balance - v  
end

-- 然后可以像这样调用 Account 的 withdraw 方法
Account.withdraw(100.00)
```

以下简单的类包含了三个属性： area, length 和 breadth，和一个 printArea 成员函数用于打印计算结果：

> Lua 使用冒号 `:` **定义和调用对象的方法**。使用冒号可以隐式地传递一个名为 self 的参数，该参数指向调用方法的对象自身。这样，在方法内部就可以通过 self 来访问对象的属性和调用其他方法。而在方法调用时，使用冒号可以简化语法，不需要显式地传递 self 参数，Lua会自动将调用者作为 self 参数传入方法。这种语法糖的使用可以使代码更加简洁易读。

```lua
-- 元类
Rectangle = {area = 0, length = 0, breadth = 0}

-- 派生类的方法 new
function Rectangle:new (o,length,breadth)
  o = o or {}
  setmetatable(o, self)
  self.__index = self
  self.length = length or 0
  self.breadth = breadth or 0
  self.area = length*breadth
  return o
end

-- 派生类的方法 printArea
function Rectangle:printArea ()
  print("矩形面积为 ",self.area)
end
```

**创建对象**

创建对象是为类的实例分配内存的过程。每个类都有属于自己的内存并共享公共数据。

```lua
r = Rectangle:new(1,10,20)
```

**访问属性**

使用点号 `.` 来访问类的属性：

```lua
print(r.length)
```

**访问成员函数**

使用冒号 `:` 来访问类的成员函数：

```lua
r:printArea()
```

内存在对象初始化时分配。