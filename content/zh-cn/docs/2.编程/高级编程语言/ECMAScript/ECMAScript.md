---
title: ECMAScript
---

# 概述

> 参考：
> - [Wiki，ECMAScript](https://en.wikipedia.org/wiki/ECMAScript)
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

# JavaScript 范儿
