---
title: object
linkTitle: object
date: 2023-12-07T11:51
weight: 20
---

# 概述

> 参考：
> 
> - [MDN，使用对象](https://developer.mozilla.org/zh-CN/docs/Web/JavaScript/Guide/Working_with_objects)

**object(对象)** 是 Javascript 语言的核心概念。所有的数据类型都可以称之为 object。

JS 中的 object 可以简单理解为面向对象编程语言中的“对象”，只不过并不用 class 这种关键字进行声明（不过，在 ES6 标注后，可以使用 class 关键字声明一个 object）

JS 的 object 也是一系列 **Property(属性)** 的集合。属性包含一个名和一个值，若属性的值是一个函数，则该属性也称为 **Method(方法)**

```js
var myObject = {
  // 其他属性...
  propertyOne: "Hello",

  // 创建 myObject 对象的方法，方法名为 methodOne
  // 这个其实就类似于一个名为 methodOne 的函数，就像 `function methodOne(t)`
  methodOne: function (a) {
    // 实际的方法体代码
    console.log("调用了 myObject 对象的方法，参数为:" + a);
  },

  methodTwo(a) {
    console.log("调用了 myObject 对象的方法，参数为:", a)
  }

  // 其他方法或属性...
};
```

> 注意：Javascript 还有一个 Object 类型的的数据也可以称为 object。。。挺绕的。。。0.0。我们通常使用 Object 的 O 字母大小写来区分~~~~
>
> - Object 是一种狭义的数据类型，与 字典、映射 等同义
> - object 是
>   - 包含数据和用于处理数据的指令的数据结构。
>   - 一种合成类型，一种在 JS 中最复杂的数据类型。
>   - 一种将任意数据类型构建为 object 类型的方法
>   - 也是一个逻辑上的对象，通过特定方法实例化的类型都可以称为对象。

```javascript
var objectType = {}
console.log(objectType)
console.log(typeof objectType)
var arrayType = []
console.log(arrayType)
console.log(typeof arrayType)
```

上面代码的输出结果如果从浏览器看的话，效果如下图：

![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/th2szn/1641957269181-cffe1052-765a-403f-b3af-7ff732eceb78.png)

可以看到，object 是一种类型，而 Array 和 Object 类型都属于 object 类型的原型类；也就是 Array 类型的 object、Object 类型的 object。

# Constructor(构造函数)

面向对象编程的第一步，就是要生成对象。对象是单个实物的抽象。通常需要一个模板，表示某一类实物的共同特征，然后对象根据这个模板生成。

典型的面向对象编程语言（比如 C++ 和 Java），都有“类”（class）这个概念。所谓“类”就是对象的模板，对象就是“类”的实例。但是，JavaScript 语言的对象体系，不是基于“类”的，而是基于 **Constructor(构造函数)** 和 **Prototype(原型类)**。

JavaScript 语言使用 **Constructor(构造函数)** 作为**对象的模板**。所谓“构造函数”，就是用来描述对象的基本结构。通过一个构造函数，可以生成多个实例对象，这些实例对象都有相同的结构。而 **Prototype(原型类)** 则是这个对象中的一个属性，用来标明该对象实例原始的类型，以便可以调用这个类型对象下的方法。

构造函数的特点有两个：

- 函数体内部使用了 `this` 关键字，代表了所要生成的对象实例。
- 生成对象的时候，必须使用 `new` 关键字。
  - 想要生成什么类型的对象，就使用对应的 Constructor，通常来说，Constructor 的名称与类型名称相同，只不过首字母大写。比如：
    - 我要构造一个 String 类型的 object，则使用 `String()` 函数。

构造函数也是一个普通的函数，只不过具有某些特有的特征和用法：

```javascript
function Person(name, age) {
  // 若没有 this 关键字，则这俩不会变为 Person 的属性，仅仅只是一个赋值操作
  this.name = name
  this.age = age
}

Person.prototype.showInfo = function () {
  return this.name + " is " + this.age + " years old."
}

var p1 = new Person("张三", 18)
```

- `Person()` 是构造函数，第一个字母通常都是大写的，且内部使用 `this` 关键字
- `name` 是构造函数的一个属性。

然后，我们可以通过 `new` 关键字生成 `Person()` 的一个实例。

如果用面向对象的概念来类比的话：

- Person 是一个对象
- name 和 age 是对象中的一个属性
- p1 是对象在实例化后的实体

![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/th2szn/1650787512457-aee108cf-2ec9-488a-82ac-363293db2764.png)

![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/th2szn/1650789606739-0e7d0eb8-f5ef-4151-ba29-d6402090becf.png)

## ES6 语法

上面这种写法跟传统的面向对象语言（比如 C++ 和 Java）差异很大，很容易让新学习这门语言的程序员感到困惑。

ES6 提供了更接近传统语言的写法，引入了 Class（类）这个概念，作为对象的模板。通过 class 关键字，可以定义类。

基本上，ES6 的 class 可以看作只是一个语法糖，它的绝大部分功能，ES5 都可以做到，新的 class 写法只是让对象原型的写法更加清晰、更像面向对象编程的语法而已。上面的代码用 ES6 的 class 改写，就是下面这样：

```javascript
class Person {
  constructor(name, age) {
    this.name = name
    this.age = age
  }
  showInfo() {
    return this.name + " is " + this.age + " years old."
  }
}
var p1 = new Person("张三", 18)
```

## Prototype(原型)

每一个构造函数都会自带一个 prototype 属性。为了解决实例化时，对象上的方法被重复创建占用过多内存空间的问题。所以，想要定义对象上的方法，就是使用 `OjbectName.pototype.MethodName` 语法。

## 内置构造函数示例

用一个最简单的声明字符串变量为例：

```javascript
// 基本字符串
var stringType = "Hello_World"
// 字符串对象
var stringObjType = new String("Hello_World")

console.log("基本字符串的类型:", typeof stringType)
console.log("字符串对象的类型:", typeof stringObjType)

console.log(stringType)
console.log(stringObjType)
```

这里的 `String()` 是一个构造函数

执行结果如下：

![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/th2szn/1640327798825-cba5f4ac-6cb9-4ed2-9c17-62a10d4fd2b2.png)

使用 `new` 构造函数后，一个普通的字符串成为了 object 类型，并且 **Prototype(原型类)** 是 String。

- 字符串字面量 (通过单引号或双引号定义) 是**基本字符串**。
- 通过 `new` 构造出来的是**字符串对象**。

注意：当基本字符串需要调用一个字符串对象才有的方法或者查询值的时候(基本字符串是没有这些方法的)，**JavaScript 会自动将基本字符串转化为字符串对象**并且调用相应的方法或者执行查询。

# 内置对象的方法

Javascript 中，提供了很多方法可以对数据直接进行操作(比如 类型转换、数组排序、遍历 等等)。在 [MDN 官方文档，Javascript 标准内置对象](https://developer.mozilla.org/zh-CN/docs/Web/JavaScript/Reference/Global_Objects) 中可以查看所有 Javascript 中内置的操作，这里是以对象类型进行分类，每个类型的对象都有很多可以操作其自身的方法。

> [这里](/docs/2.编程/高级编程语言/ECMAScript/JavaScript%20 标准库/各种类型的%20object(对象)%20 的常见方法.md 标准库/各种类型的 object(对象) 的常见方法.md)也列出了一些日常使用率非常高的对象方法

这种行为本质就是调用对象上的方法，与其他语言是类似的效果，我们使用 `new` 关键字与构造函数创建出 A 类型的实例化对象(其实就是一个变量)，这个变量是 A 数据类型，然后就可以直接调用 A 数据类型下所有可用的方法。

比如：

```javascript
const arrayObject = new Array(9, 6, 3, 1, 2, 4, 5, 7, 8, 0)

arrayObject.sort()

console.log(arrayObject)
```

我实例化了一个 Array 类型的对象：arrayObject，此时可以直接调用该类型下的方法 `sort()` 对数据进行排序

# 原型链

![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/th2szn/1650791286870-fafece55-15db-47be-ae76-3d5f80ea30de.png)

原型链是用 **proto** 串联起来的对象链状结构，该结构用来在访问对象的成员的时候，提供访问路径。

对方访问机制：

- 首先在自己身上查找，如果有直接使用
- 如果没有，自动去 **proto** 上查找
- 如果还没有，就再去 **proto** 上查找
- 直到 Ojbect.prototype 都没有，那么返回 undefiend

# 总结

其实，Javascript 的 object 可以简单理解为面向对象编程中的对象，通过构建各种类型的对象，以便使用对象上的方法。
