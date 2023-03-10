---
title: ECMAScript
---

# 概述

> 参考：
> - [Wiki,ECMAScript](https://en.wikipedia.org/wiki/ECMAScript)
> - [JavaScript 官网](https://www.javascript.com/)
> - [TypeScript 官网](https://www.typescriptlang.org/)

ECMAScript 是一种编程语言的**标准**，起源于 JavaScripts。

1996 年 8 月，微软模仿 JavaScript 开发了一种相近的语言，取名为 JScript（JavaScript 是 Netscape 的注册商标，微软不能用），首先内置于 IE 3.0。Netscape 公司面临丧失浏览器脚本语言的主导权的局面。

1996 年 11 月，Netscape 公司决定将 JavaScript 提交给国际标准化组织 ECMA（European Computer Manufacturers Association），希望 JavaScript 能够成为国际标准，以此抵抗微软。ECMA 的 39 号技术委员会（Technical Committee 39）负责制定和审核这个标准，成员由业内的大公司派出的工程师组成，目前共 25 个人。该委员会定期开会，所有的邮件讨论和会议记录，都是公开的。

1997 年 7 月，ECMA 组织发布 262 号标准文件（ECMA-262）的第一版，规定了浏览器脚本语言的标准，并将这种语言称为 ECMAScript。这个版本就是 ECMAScript 1.0 版。之所以不叫 JavaScript，一方面是由于商标的关系，Java 是 Sun 公司的商标，根据一份授权协议，只有 Netscape 公司可以合法地使用 JavaScript 这个名字，且 JavaScript 已经被 Netscape 公司注册为商标，另一方面也是想体现这门语言的制定者是 ECMA，不是 Netscape，这样有利于保证这门语言的开放性和中立性。因此，ECMAScript 和 JavaScript 的关系是，前者是后者的规范，后者是前者的一种实现。在日常场合，这两个词是可以互换的。

ECMAScript 只用来标准化 JavaScript 这种语言的基本语法结构，与部署环境相关的标准都由其他标准规定，比如 DOM 的标准就是由 W3C 组织（World Wide Web Consortium）制定的。

ECMA-262 标准后来也被另一个国际标准化组织 ISO（International Organization for Standardization）批准，标准号是 ISO-16262。

## ES6 标准

ECMAScript 6.0 是 ECMA 的最新标准，于 2015 年 6 月发布，官方称为 ES2015 标准(ES6 的叫法更民间)。

2011 年，ECMAScript 5.1 版发布后，就开始制定 6.0 版了。因此，ES6 这个词的原意，就是指 JavaScript 语言的下一个版本。

但是，因为这个版本引入的语法功能太多，而且制定过程当中，还有很多组织和个人不断提交新功能。事情很快就变得清楚了，不可能在一个版本里面包括所有将要引入的功能。常规的做法是先发布 6.0 版，过一段时间再发 6.1 版，然后是 6.2 版、6.3 版等等。

但是，标准的制定者不想这样做。他们想让标准的升级成为常规流程：任何人在任何时候，都可以向标准委员会提交新语法的提案，然后标准委员会每个月开一次会，评估这些提案是否可以接受，需要哪些改进。如果经过多次会议以后，一个提案足够成熟了，就可以正式进入标准了。这就是说，标准的版本升级成为了一个不断滚动的流程，每个月都会有变动。

标准委员会最终决定，标准在每年的 6 月份正式发布一次，作为当年的正式版本。接下来的时间，就在这个版本的基础上做改动，直到下一年的 6 月份，草案就自然变成了新一年的版本。这样一来，就不需要以前的版本号了，只要用年份标记就可以了。

ES6 的第一个版本，就这样在 2015 年 6 月发布了，正式名称就是《ECMAScript 2015 标准》（简称 ES2015）。2016 年 6 月，小幅修订的《ECMAScript 2016 标准》（简称 ES2016）如期发布，这个版本可以看作是 ES6.1 版，因为两者的差异非常小（只新增了数组实例的 includes 方法和指数运算符），基本上是同一个标准。根据计划，2017 年 6 月发布 ES2017 标准。

因此，ES6 既是一个历史名词，也是一个泛指，含义是 5.1 版以后的 JavaScript 的下一代标准，涵盖了 ES2015、ES2016、ES2017 等等，而 ES2015 则是正式名称，特指该年发布的正式版本的语言标准。本书中提到 ES6 的地方，一般是指 ES2015 标准，但有时也是泛指“下一代 JavaScript 语言”。

## ECMAScript 的三大核心组成

**ECMAScript**

- JS 的书写语法和规则

**Browser Ojbect Model(浏览器对象模型，简称 BOM)**

- JS 控制浏览器的属性和方法。
  - 比如浏览器右侧的滚动条，可以通过 JS 代码来控制。比如某些网页有个叫回到顶部的按钮，按一下，就等于是 JS 操作滚动条移动到最上面了。
  - 比如很多手机，我们从屏幕最左侧往右滑，一般返回上一页。这是因为这个滑动行为被 JS 代码捕获后，操作浏览器点击了一下后腿按钮。

**Document Object Model(文档对象模型，简称 DOM)**

- JS 控制文档流的属性和方法。
  - 比如很多网页最上面中间都有一个图片，按一下图片左右两边的箭头，就会换到另一个图片
  - 也就是说，JS 控制什么时候，让页面元素发生一些变化

# 学习资料

[MDN 官方文档，Web 开发技术](https://developer.mozilla.org/en-US/docs/Web)((通常指的是网站首页的 References 标签中的文档))

- [JavaScript](https://developer.mozilla.org/en-US/docs/Web/JavaScript)
- [JavaScript-参考-词汇文法-关键字](https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Lexical_grammar#keywords)
- [JavaScript-参考-语句和声明](https://developer.mozilla.org/zh-CN/docs/Web/JavaScript/Reference/Statements)(这里就是 JS 的关键字的用法)

[廖雪峰，JavaScript](https://www.liaoxuefeng.com/wiki/1022910821149312)
[网道，JavaScript](https://wangdoc.com/javascript/index.html)

- [网道，JavaScript-JavaScript 的基本语法](https://wangdoc.com/javascript/basic/grammar.html)

电子书

- [GitHub 项目，javascript/zh.javascript.info](https://github.com/javascript-tutorial/zh.javascript.info/tree/master)(现代 JavaScript 教程)

# Hello World

```html
<!DOCTYPE html>
<html>
  <head>
    <meta charset="UTF-8" />
    <title>Hello World</title>
  </head>
  <body>
    <!-- 有多种方式可以在书写 JS 代码 -->
    <!-- 行内式， JS 代码写在标签上-->
    <!-- a 标签，书写在 href 属性上 -->
    <a href="javascript: alert('Hello World，行内式，a 标签');">点我</a>
    <!-- 非 a 标签，书写在行为属性上 -->
    <div onclick="alert('Hello World，行内式，非 a 标签')">点我</div>

    <!-- 内嵌式，JS 代码写在 script 标签中 -->
    <script>
      // 在前端页面上显示的内容
      alert("Hello World 内嵌式")
      // 在后端控制台显示的内容
      console.log("Hello World backend")
    </script>

    <!-- 外链试，JS 代码写在单独的 .js 文件中，并通过 script 标签的 src 属性引入 .js 文件 -->
    <script src="./hello_world.js"></script>

    <!-- 总结：
            行内式 # 强烈不推荐。不利于代码维护，也会导致 HTML 文件过于臃肿。
            内嵌式 # 一般测试或者学习时，使用这种方式，不用建立很多 .js 文件。
            外链式 # 强烈推荐
         -->
  </body>
</html>
```

# JavaScript 语言关键字

break
case
catch
class
const
continue
debugger
default
delete
do
else
enum
export
extends
false
finally
for
function
if
implements
import
in
instanceof
interface
let
new
null
package
private
protected
public
return
super
switch
static
this
throw
try
true
typeof
var
void
while
with
yield

# JavaScript 基本语法规范

## 语句

JavaScript 程序的执行单位为行（line），也就是一行一行地执行。一般情况下，每一行就是一个语句。

语句（statement）是为了完成某种任务而进行的操作，比如下面就是一行赋值语句。

    var a = 1 + 3;

这条语句先用`var`命令，声明了变量`a`，然后将`1 + 3`的运算结果赋值给变量`a`。

`1 + 3`叫做表达式（expression），指一个为了得到返回值的计算式。语句和表达式的区别在于，前者主要为了进行某种操作，一般情况下不需要返回值；后者则是为了得到返回值，一定会返回一个值。凡是 JavaScript 语言中预期为值的地方，都可以使用表达式。比如，赋值语句的等号右边，预期是一个值，因此可以放置各种表达式。

语句以分号结尾，一个分号就表示一个语句结束。多个语句可以写在一行内。

    var a = 1 + 3 ; var b = 'abc';

分号前面可以没有任何内容，JavaScript 引擎将其视为空语句。

    ;;;

上面的代码就表示 3 个空语句。

表达式不需要分号结尾。一旦在表达式后面添加分号，则 JavaScript 引擎就将表达式视为语句，这样会产生一些没有任何意义的语句。

    1 + 3;
    'abc';

上面两行语句只是单纯地产生一个值，并没有任何实际的意义。

## 变量

### 概念

变量是对“值”的具名引用。变量就是为“值”起名，然后引用这个名字，就等同于引用这个值。变量的名字就是变量名。

    var a = 1;

上面的代码先声明变量`a`，然后在变量`a`与数值 1 之间建立引用关系，称为将数值 1“赋值”给变量`a`。以后，引用变量名`a`就会得到数值 1。最前面的`var`，是变量声明命令。它表示通知解释引擎，要创建一个变量`a`。

注意，JavaScript 的变量名区分大小写，`A`和`a`是两个不同的变量。

变量的声明和赋值，是分开的两个步骤，上面的代码将它们合在了一起，实际的步骤是下面这样。

    var a;
    a = 1;

如果只是声明变量而没有赋值，则该变量的值是`undefined`。`undefined`是一个特殊的值，表示“无定义”。

    var a;
    a

如果变量赋值的时候，忘了写`var`命令，这条语句也是有效的。

    var a = 1;

    a = 1;

但是，不写`var`的做法，不利于表达意图，而且容易不知不觉地创建全局变量，所以建议总是使用`var`命令声明变量。

如果一个变量没有声明就直接使用，JavaScript 会报错，告诉你变量未定义。

    x

上面代码直接使用变量`x`，系统就报错，告诉你变量`x`没有声明。

可以在同一条`var`命令中声明多个变量。

    var a, b;

JavaScript 是一种动态类型语言，也就是说，变量的类型没有限制，变量可以随时更改类型。

    var a = 1;
    a = 'hello';

上面代码中，变量`a`起先被赋值为一个数值，后来又被重新赋值为一个字符串。第二次赋值的时候，因为变量`a`已经存在，所以不需要使用`var`命令。

如果使用`var`重新声明一个已经存在的变量，是无效的。

    var x = 1;
    var x;
    x

上面代码中，变量`x`声明了两次，第二次声明是无效的。

但是，如果第二次声明的时候还进行了赋值，则会覆盖掉前面的值。

    var x = 1;
    var x = 2;

    var x = 1;
    var x;
    x = 2;

### Variables Hoisting(变量提升)

JavaScript 引擎的工作方式是，先解析代码，获取所有被声明的变量，然后再一行一行地运行。这造成的结果，就是所有的变量的声明语句，都会被提升到代码的头部，这就叫做变量提升（hoisting）。

    console.log(a);
    var a = 1;

上面代码首先使用`console.log`方法，在控制台（console）显示变量`a`的值。这时变量`a`还没有声明和赋值，所以这是一种错误的做法，但是实际上不会报错。因为存在变量提升，真正运行的是下面的代码。

    var a;
    console.log(a);
    a = 1;

最后的结果是显示`undefined`，表示变量`a`已声明，但还未赋值。

## 标识符

标识符（identifier）指的是用来识别各种值的合法名称。最常见的标识符就是变量名，以及后面要提到的函数名。JavaScript 语言的标识符对大小写敏感，所以`a`和`A`是两个不同的标识符。

标识符有一套命名规则，不符合规则的就是非法标识符。JavaScript 引擎遇到非法标识符，就会报错。

简单说，标识符命名规则如下。

- 第一个字符，可以是任意 Unicode 字母（包括英文字母和其他语言的字母），以及美元符号（`$`）和下划线（`_`）。
- 第二个字符及后面的字符，除了 Unicode 字母、美元符号和下划线，还可以用数字`0-9`。

下面这些都是合法的标识符。

    arg0
    _tmp
    $elem
    π

下面这些则是不合法的标识符。

    1a
    23
    ***
    a+b
    -d

中文是合法的标识符，可以用作变量名。

    var 临时变量 = 1;

> JavaScript 有一些保留字，不能用作标识符：arguments、break、case、catch、class、const、continue、debugger、default、delete、do、else、enum、eval、export、extends、false、finally、for、function、if、implements、import、in、instanceof、interface、let、new、null、package、private、protected、public、return、static、super、switch、this、throw、true、try、typeof、var、void、while、with、yield。

## 注释

源码中被 JavaScript 引擎忽略的部分就叫做注释，它的作用是对代码进行解释。JavaScript 提供两种注释的写法：一种是单行注释，用`//`起头；另一种是多行注释，放在`/*`和`*/`之间。

此外，由于历史上 JavaScript 可以兼容 HTML 代码的注释，所以`<!--`和`-->`也被视为合法的单行注释。

    x = 1; <!-- x = 2;
    --> x = 3;

上面代码中，只有`x = 1`会执行，其他的部分都被注释掉了。

需要注意的是，`-->`只有在行首，才会被当成单行注释，否则会当作正常的运算。

    function countdown(n) {
      while (n --> 0) console.log(n);
    }
    countdown(3)

上面代码中，`n --> 0`实际上会当作`n-- > 0`，因此输出 2、1、0。

## 区块

JavaScript 使用大括号，将多个相关的语句组合在一起，称为“区块”（block）。

对于`var`命令来说，JavaScript 的区块不构成单独的作用域（scope）。

    {
      var a = 1;
    }

    a

上面代码在区块内部，使用`var`命令声明并赋值了变量`a`，然后在区块外部，变量`a`依然有效，区块对于`var`命令不构成单独的作用域，与不使用区块的情况没有任何区别。在 JavaScript 语言中，单独使用区块并不常见，区块往往用来构成其他更复杂的语法结构，比如`for`、`if`、`while`、`function`等。

## 条件语句

JavaScript 提供`if`结构和`switch`结构，完成条件判断，即只有满足预设的条件，才会执行相应的语句。

### if 结构

`if`结构先判断一个表达式的布尔值，然后根据布尔值的真伪，执行不同的语句。所谓布尔值，指的是 JavaScript 的两个特殊值，`true`表示“真”，`false`表示“伪”。

    if (布尔值)
      语句;

    if (布尔值) 语句;

上面是`if`结构的基本形式。需要注意的是，“布尔值”往往由一个条件表达式产生的，必须放在圆括号中，表示对表达式求值。如果表达式的求值结果为`true`，就执行紧跟在后面的语句；如果结果为`false`，则跳过紧跟在后面的语句。

    if (m === 3)
      m = m + 1;

上面代码表示，只有在`m`等于 3 时，才会将其值加上 1。

这种写法要求条件表达式后面只能有一个语句。如果想执行多个语句，必须在`if`的条件判断之后，加上大括号，表示代码块（多个语句合并成一个语句）。

    if (m === 3) {
      m += 1;
    }

建议总是在`if`语句中使用大括号，因为这样方便插入语句。

注意，`if`后面的表达式之中，不要混淆赋值表达式（`=`）、严格相等运算符（`===`）和相等运算符（`==`）。尤其是赋值表达式不具有比较作用。

    var x = 1;
    var y = 2;
    if (x = y) {
      console.log(x);
    }

上面代码的原意是，当`x`等于`y`的时候，才执行相关语句。但是，不小心将严格相等运算符写成赋值表达式，结果变成了将`y`赋值给变量`x`，再判断变量`x`的值（等于 2）的布尔值（结果为`true`）。

这种错误可以正常生成一个布尔值，因而不会报错。为了避免这种情况，有些开发者习惯将常量写在运算符的左边，这样的话，一旦不小心将相等运算符写成赋值运算符，就会报错，因为常量不能被赋值。

    if (x = 2) {
    if (2 = x) {

至于为什么优先采用“严格相等运算符”（`===`），而不是“相等运算符”（`==`），请参考《运算符》章节。

### if...else 结构

`if`代码块后面，还可以跟一个`else`代码块，表示不满足条件时，所要执行的代码。

    if (m === 3) {

    } else {

    }

上面代码判断变量`m`是否等于 3，如果等于就执行`if`代码块，否则执行`else`代码块。

对同一个变量进行多次判断时，多个`if...else`语句可以连写在一起。

    if (m === 0) {

    } else if (m === 1) {

    } else if (m === 2) {

    } else {

    }

`else`代码块总是与离自己最近的那个`if`语句配对。

    var m = 1;
    var n = 2;

    if (m !== 1)
    if (n === 2) console.log('hello');
    else console.log('world');

上面代码不会有任何输出，`else`代码块不会得到执行，因为它跟着的是最近的那个`if`语句，相当于下面这样。

    if (m !== 1) {
      if (n === 2) {
        console.log('hello');
      } else {
        console.log('world');
      }
    }

如果想让`else`代码块跟随最上面的那个`if`语句，就要改变大括号的位置。

    if (m !== 1) {
      if (n === 2) {
        console.log('hello');
      }
    } else {
      console.log('world');
    }

### switch 结构

多个`if...else`连在一起使用的时候，可以转为使用更方便的`switch`结构。

    switch (fruit) {
      case "banana":

        break;
      case "apple":

        break;
      default:

    }

上面代码根据变量`fruit`的值，选择执行相应的`case`。如果所有`case`都不符合，则执行最后的`default`部分。需要注意的是，每个`case`代码块内部的`break`语句不能少，否则会接下去执行下一个`case`代码块，而不是跳出`switch`结构。

    var x = 1;

    switch (x) {
      case 1:
        console.log('x 等于1');
      case 2:
        console.log('x 等于2');
      default:
        console.log('x 等于其他值');
    }

上面代码中，`case`代码块之中没有`break`语句，导致不会跳出`switch`结构，而会一直执行下去。正确的写法是像下面这样。

    switch (x) {
      case 1:
        console.log('x 等于1');
        break;
      case 2:
        console.log('x 等于2');
        break;
      default:
        console.log('x 等于其他值');
    }

`switch`语句部分和`case`语句部分，都可以使用表达式。

    switch (1 + 3) {
      case 2 + 2:
        f();
        break;
      default:
        neverHappens();
    }

上面代码的`default`部分，是永远不会执行到的。

需要注意的是，`switch`语句后面的表达式，与`case`语句后面的表示式比较运行结果时，采用的是严格相等运算符（`===`），而不是相等运算符（`==`），这意味着比较时不会发生类型转换。

    var x = 1;

    switch (x) {
      case true:
        console.log('x 发生类型转换');
        break;
      default:
        console.log('x 没有发生类型转换');
    }

上面代码中，由于变量`x`没有发生类型转换，所以不会执行`case true`的情况。这表明，`switch`语句内部采用的是“严格相等运算符”，详细解释请参考《运算符》一节。

### 三元运算符 ?:

JavaScript 还有一个三元运算符（即该运算符需要三个运算子）`?:`，也可以用于逻辑判断。

    (条件) ? 表达式1 : 表达式2

上面代码中，如果“条件”为`true`，则返回“表达式 1”的值，否则返回“表达式 2”的值。

    var even = (n % 2 === 0) ? true : false;

上面代码中，如果`n`可以被 2 整除，则`even`等于`true`，否则等于`false`。它等同于下面的形式。

    var even;
    if (n % 2 === 0) {
      even = true;
    } else {
      even = false;
    }

这个三元运算符可以被视为`if...else...`的简写形式，因此可以用于多种场合。

    var myVar;
    console.log(
      myVar ?
      'myVar has a value' :
      'myVar does not have a value'
    )

上面代码利用三元运算符，输出相应的提示。

    var msg = '数字' + n + '是' + (n % 2 === 0 ? '偶数' : '奇数');

上面代码利用三元运算符，在字符串之中插入不同的值。

## 循环语句

循环语句用于重复执行某个操作，它有多种形式。

### while 循环

`While`语句包括一个循环条件和一段代码块，只要条件为真，就不断循环执行代码块。

    while (条件)
      语句;

    while (条件) 语句;

`while`语句的循环条件是一个表达式，必须放在圆括号中。代码块部分，如果只有一条语句，可以省略大括号，否则就必须加上大括号。

    while (条件) {
      语句;
    }

下面是`while`语句的一个例子。

    var i = 0;

    while (i < 100) {
      console.log('i 当前为：' + i);
      i = i + 1;
    }

上面的代码将循环 100 次，直到`i`等于 100 为止。

下面的例子是一个无限循环，因为循环条件总是为真。

    while (true) {
      console.log('Hello, world');
    }

### for 循环

`for`语句是循环命令的另一种形式，可以指定循环的起点、终点和终止条件。它的格式如下。

    for (初始化表达式; 条件; 递增表达式)
      语句

    for (初始化表达式; 条件; 递增表达式) {
      语句
    }

`for`语句后面的括号里面，有三个表达式。

- 初始化表达式（initialize）：确定循环变量的初始值，只在循环开始时执行一次。
- 条件表达式（test）：每轮循环开始时，都要执行这个条件表达式，只有值为真，才继续进行循环。
- 递增表达式（increment）：每轮循环的最后一个操作，通常用来递增循环变量。

下面是一个例子。

    var x = 3;
    for (var i = 0; i < x; i++) {
      console.log(i);
    }

上面代码中，初始化表达式是`var i = 0`，即初始化一个变量`i`；测试表达式是`i < x`，即只要`i`小于`x`，就会执行循环；递增表达式是`i++`，即每次循环结束后，`i`增大 1。

所有`for`循环，都可以改写成`while`循环。上面的例子改为`while`循环，代码如下。

    var x = 3;
    var i = 0;

    while (i < x) {
      console.log(i);
      i++;
    }

`for`语句的三个部分（initialize、test、increment），可以省略任何一个，也可以全部省略。

    for ( ; ; ){
      console.log('Hello World');
    }

上面代码省略了`for`语句表达式的三个部分，结果就导致了一个无限循环。

### do...while 循环

`do...while`循环与`while`循环类似，唯一的区别就是先运行一次循环体，然后判断循环条件。

    do
      语句
    while (条件);

    do {
      语句
    } while (条件);

不管条件是否为真，`do...while`循环至少运行一次，这是这种结构最大的特点。另外，`while`语句后面的分号注意不要省略。

下面是一个例子。

    var x = 3;
    var i = 0;

    do {
      console.log(i);
      i++;
    } while(i < x);

### break 语句和 continue 语句

`break`语句和`continue`语句都具有跳转作用，可以让代码不按既有的顺序执行。

`break`语句用于跳出代码块或循环。

    var i = 0;

    while(i < 100) {
      console.log('i 当前为：' + i);
      i++;
      if (i === 10) break;
    }

上面代码只会执行 10 次循环，一旦`i`等于 10，就会跳出循环。

`for`循环也可以使用`break`语句跳出循环。

    for (var i = 0; i < 5; i++) {
      console.log(i);
      if (i === 3)
        break;
    }

上面代码执行到`i`等于 3，就会跳出循环。

`continue`语句用于立即终止本轮循环，返回循环结构的头部，开始下一轮循环。

    var i = 0;

    while (i < 100){
      i++;
      if (i % 2 === 0) continue;
      console.log('i 当前为：' + i);
    }

上面代码只有在`i`为奇数时，才会输出`i`的值。如果`i`为偶数，则直接进入下一轮循环。

如果存在多重循环，不带参数的`break`语句和`continue`语句都只针对最内层循环。

### 标签（label）

JavaScript 语言允许，语句的前面有标签（label），相当于定位符，用于跳转到程序的任意位置，标签的格式如下。

    label:
      语句

标签可以是任意的标识符，但不能是保留字，语句部分可以是任意语句。

标签通常与`break`语句和`continue`语句配合使用，跳出特定的循环。

    top:
      for (var i = 0; i < 3; i++){
        for (var j = 0; j < 3; j++){
          if (i === 1 && j === 1) break top;
          console.log('i=' + i + ', j=' + j);
        }
      }

上面代码为一个双重循环区块，`break`命令后面加上了`top`标签（注意，`top`不用加引号），满足条件时，直接跳出双层循环。如果`break`语句后面不使用标签，则只能跳出内层循环，进入下一次的外层循环。

标签也可以用于跳出代码块。

    foo: {
      console.log(1);
      break foo;
      console.log('本行不会输出');
    }
    console.log(2);

上面代码执行到`break foo`，就会跳出区块。

`continue`语句也可以与标签配合使用。

    top:
      for (var i = 0; i < 3; i++){
        for (var j = 0; j < 3; j++){
          if (i === 1 && j === 1) continue top;
          console.log('i=' + i + ', j=' + j);
        }
      }

上面代码中，`continue`命令后面有一个标签名，满足条件时，会跳过当前循环，直接进入下一轮外层循环。如果`continue`语句后面不使用标签，则只能进入下一轮的内层循环。

## 参考链接

- Axel Rauschmayer, [A quick overview of JavaScript](http://www.2ality.com/2011/10/javascript-overview.html)
