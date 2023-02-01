---
title: Jinja
---

# 概述

> 参考：
> - [GitHub 项目](https://github.com/pallets/jinja)
> - [官网](https://jinja.palletsprojects.com/)
> - [Wiki,Jinja](<https://en.wikipedia.org/wiki/Jinja_(template_engine)>)
> - [国人翻译官网](http://docs.jinkan.org/docs/jinja2/)
> - <https://www.junmajinlong.com/ansible/9_power_of_jinja2/>

Jinja 是一个用于 Python 变成语言中的 **Template Engine(模板引擎)**。Jinja 通常被用来作为 Python 的 Web 框架(e.g.Flask、Django)的数据渲染的底层调用。

> Django 其实自带模板引擎(DTL)，只不过由于 Jinja 的流行，通常都让 Django 的模板引擎使用 Jinja2

Jinja 模板引擎允许定制标签、过滤器、测试和全局变量。此外，与 Django 模板引擎不同，Jinja 允许模板设计器调用带有对象参数的函数。Jinja 是 Flask 的默认模板引擎，同时，也被 Ansible、Trac、Salt 使用。

## Jinja 是什么？模板是什么？

何为模板？举个例子就知道了。

假设要发送一个文件给一个或多个目标节点，要发送的文件内容如下：

    hello, __NAME__

其中 `__NAME__` 部分想要根据目标节点的主机名来确定，比如发送给 www 节点时内容应该为 `hello, www`，发送给 wwww 节点时，内容应该为`hello, wwww`。换句话说，`__NAME__` 是一个能够根据不同场景动态生成不同字符串的代码小片段。而根据特殊的代码片段动态生成字符串便是模板要实现的功能。

现在解释模板便容易了：**所谓 Template(模板)，只是文本文件，可以在文本字符串中嵌入一些 Expressions(表达式)，然后使用模板引擎去解析整个模板，将其中嵌入的表达式替换成对应的结果**。其中，**解析并替换模板表达式的过程**，称为**渲染**。从编程语言的角度说，**表达式就是代码中的 function。**

当模板引擎解析表达式时，每个**表达式**都**有返回值**。

为了让模板引擎只替换表达式而不操作普通字符串，所以模板引擎需要能够区分模板表达式和普通字符串，所以模板表达式通常会使用 **Delimiters(分隔符)** 包围起来。例如上面示例中，`__NAME__`使用了前后两个下划线包围，表示这部分是模板表达式，它是需要进行替换的，而”hello” 是普通字符串，模板引擎不会去管它。

Jinja 模板引擎提供了三种** Delimiters(分隔符) **来包围 **模板表达式**：

- `{{ XXX }}` # Jinja 的 **Basic(基本)** 表达式。比如引用变量、使用过滤器
- `{% XXX %}` # Jinja 的 **Statements(语句)** 表达式，比如 定义变量、执行控制结构语句(if 语句、for 循环)、等等
- `{# XXX #}` # Jinja 的注释符号

这些表达式可以用来定义变量、引用变量、定义函数、调用函数、执行控制结构等等。**说白了，模板其实也是一种变成语言，只不过代码是夹杂在常规字符串中的，并且使用特殊的符号，将这些代码包围起来。**

模板更多用在 web 编程中来生成 HTML 页面，但绝不限于 web 编程，它可以用在很多方面，比如 Ansible 就使用 Jinja2 模板引擎来解析 YAML 中的字符串，也用在 template 模块渲染模板文件。每种编程语言都有模板，比如 Python 的模板语言称为 Jinja、而 go 的模板语言就成为 go 模板、等等，通常来说，这些变成语言的模板表达式的语法，都是 `{{ XX }}` 符号。

## 模板文件的扩展名

任何文件都可以作为模板加载，无论扩展名是什么。但是使用 `.jinja` 作为扩展名，可以是某些 IDE 的插件更容易识别模板文件以提供 代码高亮、代码补全 等功能。

另外一个识别模板的好方法，是将它们都放在 `templates` 目录中，而不用管扩展名是什么。这是一个项目最常见的用法。

# Literal(字面量)

Jinja 的 Literal(字面量) 是 最简单、最直接 的表达式形式。但是，这个其实没啥用~~~~毕竟是在文本文件中使用模板表达式，如果 Lieral 都是 字符串、数值、字典、等等 的话，直接在文本中写就好啦~~~~

> 所谓的 Literal(字面量) 从中文角度看，就是所见即所得，比如我输入 "Hello World"，看到的就是这几个字母，这是一个字符串。Literal 更容易理解的词是 Data Type(数据类型)。

所以这里主要是定义一下解析表达式后可以返回的数据类型。Jinja 的基本数据类型有如下几种：

- 字符串 # 双引号或单引号中间的一切都是字符串
  - `"Hello World"`
- 整数和浮点数 # 直接写下数值即可
  - `42` 或者 `42.23`
- 列表
  - `['list','of','objects']`
- 元组
  - `('tuple','of','values')`
- 字典
  - `{'dict': 'of', 'key': 'and', 'value': 'pairs'}`
- 布尔 # 不带引号的 true 与 false
  - `true` 和 `false`

# Variable(变量) 和 作用域

模板中的变量可以通过两种方式获得

- 在模板中使用 `{% set VAR = VALUE %}` 表达式进行定义
- 在编程语言中，由代码在调用模板引擎函数时，在参数中定义。

## 定义变量

```python
{% set VAR = VALUE %}
{% set VAR_1,VAR_2 = VALUE_1,VALUE_2 %}
```

例如：

```python
{% set mystr = "hello world" %}
{% set mylist1,mylist2 = [1,2,3],[4,5,6] %}
```

## 引用变量

Jinja 模板语言中，引用变量是最基本、最简单的一种表达式。可以直接使用 `{{ Variable }}` 引用一个变量。比如：

```python
{% set NAME = "DesistDaydrem" %}
{{ NAME }}
```

Jinja 模板引擎允许使用点 `.` 来访问列表或字典类型的变量，比如 `mylist=["a","b","c"]` 列表，在 Jinja 中既可以使用 `mylist[1]` 来访问第二个元素，也可以使用`mylist.1`来访问它。

在之前的文章中曾解释过这两种访问方式的区别，这里再重复一遍：

- **使用 **`**X.Y**`** 时，先搜索 Ptyhon 对象的属性名或方法名，搜索不到时再搜索 Jinja 变量**
- **使用 **`**X["Y"]**`** 时，先搜索 Jinja 变量，搜索失败时，再搜索 Pythons 对象的属性名或方法名。**

所以，使用 `X.Y` 方式时需要小心一些，使用 `X["Y"]` 更保险。当然，使用哪种方式都无所谓，出错了也知道如何去调整。

## 变量的作用域

如果是在 if、for 等语句块外进行的变量赋值，则可以在 if、for 等语句块内引用。例如：

```python
{% set mylist = [1,2,3] %}
{% for item in mylist %}
  {{item}}
{% endfor %}
```

但是除 if 语句块外，其它类型的语句块都有自己的作用域。比如 for 语句块内部赋值的变量在 for 退出后就消失。

例如：

```python
{% set mylist = [1,2,3] %}
{% set mystr = "hello" %}
{% for item in mylist %}
  {{item}}
  {%set mystr="world"%}
{% endfor %}
{{mystr}}
```

最后一行渲染的结果是`hello`而不是`world`。

# 运算符

## Math(算术) 运算

- `**+**
` # 将两个对象相加。通常对象是数字，但如果两者都是字符串或列表，您可以通过这种方式连接它们。然而，这不是连接字符串的首选方式！对于字符串连接，请查看 `~` 运算符。 `{{ 1 + 1 }}` 表达式的返回值为 2。
  - `+` 操作符也可用于字符串串联、列表相加，例如`"a"+"b"`得到”ab”，`[1,2]+[3,4]`得到`[1,2,3,4]`
- `-`

#

- `*` #
  - `*` 也可用于重复字符串，例如`"-" * 10`得到 10 个连续的短横线
- `/`

#

- `/` 是浮点数除法，例如 3/2 得到 1.5
- `//
` #
  - `//` 是截断式整除法，例如 20/7 得到 2
- `%`

#

<a name="iaI8V"></a>

## Comparisons(比较) 运算

- `>` #
- `<
` #
- `>=`

#

- `<=
` #
- `==
` #
- `!=` #
  需要说明一点：比较操作不仅仅只能比较数值，也能比较其它对象，比如字符串。

例如 `"hey" > "hello"` 返回 True。

<a name="FdfCg"></a>

## Logic(逻辑)

- `not
` #
- `and`
  3
- `or
` #- `(expr)` # <a name="RfC05"></a>

## 其他运算符

- `in` # 成员测试，测试是否在容器内
- `is` # 做 is 测试，参见后文
- `|` # 过滤器，参见后文
- `~` # 字符串串联 <a name="uDAtA"></a>

## 总结与说明

- `in` 运算符可测试多种容器，常见的包括：
  - 列表测试 `3 in [1,2,3]` - 字符串测试 `"h" in "hey"`
  - 字典测试 `"name" in {"name":"j333333unma","age":28}`
  - 上述几种测试结果都是 True
- `is` 运算符可以做很多测试，比如测试是否是数值、是否是字符串、变量是否定义、等等
- `+` 可以做字符串串联，`~` 也可以做字符串串联，例如 `"ab" ~ "cd"` 运算结果为 `"abcd"`
- `not` 运算符和 `is`、`in` 结合时，可以放在两个位置。例如：
  - `not ("h" in "hey")` 和 `"h" not in "hey"` 两者是等价的
  - `not (3 is number())` 和 `3 is not number()` 两者是等价的 <a name="D9uBA"></a>

# Control Structures(控制结构)

> 官方文档：<https://jinja.palletsprojects.com/en/latest/templates/>

- For(循环)
- If(条件判断)
- Macros(宏)- Call(调用)
- Filters(过滤)
- Assignments(赋值)
- Block Assignments(块赋值)
- Extends(继承)
- Blocks(块)
- Inclued(包含)
- Import(导入) <a name="b2e8771d"></a>

## 条件判断

Jinja 中可以使用 if...else... 语句进行条件判断，其语法为：

    {% if CONDITION_1 %}
      string_or_expression1
    {% elif CONDITION_2 %}
      string_or_expression2
    {% elif CONDITION_3 %}
      string_or_expression3
    {% else %}
      string_or_expression4
    {% endif %}

其中 elif 和 else 分支都是可省略的。CONDITION 部分是条件表达式，关于 Jinja 支持的条件表达式，后面会介绍。

例如，模板文件 a.txt.j2 内容如下：

    今天星期几：
    {% if whatday == "0" %}
      星期日
    {% elif whatday == "1" %}
      星期一
    {% elif whatday == "2" %}
      星期二
    {% elif whatday == "3" %}
      星期三
    {% elif whatday == "4" %}
      星期四
    {% elif whatday == "5" %}
      星期五
    {% elif whatday == "6" %}
      星期六
    {% else %}
      错误数值
    {% endif %}

上面判断变量 whatday 的值，然后输出对应的星期几。因为 whatday 变量的值是字符串，所以让它和字符串形式的数值进行等值比较。当然，也可以使用筛选器将字符串转换为数值后进行数值比较：`whatday|int == 0`。

playbook 内容如下：

    ---
    - hosts: localhost
      gather_facts: no
      vars_prompt:
        - name: whatday
          default: 0
          prompt: "星期几(0->星期日,1->星期一...):"
          private: no
      tasks:
        - template:
            src: a.txt.j2
            dest: /tmp/a.txt

如果 if 语句的分支比较简单 (没有 elif 逻辑)，那么可以使用行内 if 表达式。

其语法格式为：

    string_or_expr1 if CONDITION else string_or_expr2

因为行内 if 是表达式而不是语句块，所以不使用 {%%} 符号，而使用 {{}} 。
例如：

```yaml
- debug:
    msg: "{{'周末' if whatday|int > 5 else '工作日'}}"
```

<a name="yz4zS"></a>

### is 运算符

jinja 使用 `is` 关键字，对表达式的渲染结果进行测试，测试结果有两种 true 和 false。常用在 `{% if %}` 表达式中。

比如 `name is defined` 则表示对 name 这个表达式进行测试这个表达式，会根据名为 name 的变量是否被定义，返回 true 或 false。

`is` 运算符可以做很多测试操作，比如测试是否是数值，是否是字符串等等。下表列出了所有 Jinja2 内置的测试函数。

| callable()    | even() | le()       | none()   | string()    |
| ------------- | ------ | ---------- | -------- | ----------- |
| defined()     | ge()   | lower()    | number() | undefined() |
| divisibleby() | gt()   | lt()       | odd()    | upper()     |
| eq()          | in()   | mapping()  | sameas() | escaped()   |
| iterable()    | ne()   | sequence() |          |             |

其中 callable()、escaped() 和 sameas() 在 Ansible 中几乎用不上，所以不解释。

除了 Jinja2 的内置测试函数，Ansible 还有自己扩展的 is 测试函数，在后文我会统一列出来。

这些测试函数有些可带参数、有些不带参数，当不带参数或只带一个参数时，括号可以省略。

下面演示两个测试是否是字符串的示例，让大家知道该如何使用这些测试函数进行测试，之后直接解释各测试函数的意义便可。

示例一：在 when 条件中测试

```yaml
- debug:
    msg: "a string"
  when: name is string
  vars:
    name: junmajinlong
```

示例二：在 jinja 模板的 if 中测试。测试 name 变量，如果不是字符串类型，则为真。

    {% if name is not string %}
    HELLOWORLD
    {% endif %}

下面是各测试函数的作用：

- `defined()`和`undefined()` # 测试变量是否已定义
- `number()`、`string()`、`none()` # 测试是否是一个数值、字符串、None
- `lt()`、`le()`、`gt()`、`ge()`、`eq()`、`ne()` # 分别测试是否小于、小于等于、大于、大于等于、等于、不等于
- `lower()`和`upper()` # 测试字符串是否全小写、全大写
- `even()`和`odd()` # 测试 value 是偶数还是奇数
- `divisibleby(num)` # 测试是否能被 num 整除，例如`18 is divisibleby 3`
- `in(seq)` # 测试是否在 seq 中。例如`3 is in([1,2,3])`、`"h" is in("hey")`
- `mapping()` # 测试是否是一个字典
- `iterable()` # 测试是否可迭代，在 Ansible 中一般就是 list、dict、字符串，当然，在不同场景下可能还会遇到其它可迭代的结构
- `sequence()` # 测试是否是一个序列结构，在 Ansible 中一般是 list、dict、字符串 (注：字符串、dict 在 python 中不是序列，但是在 Jinja2 测试中，属于序列)

纵观上面的内置测试函数，似乎并没有提供直接测试是否是一个列表类型的功能，但在 Ansible 中却会经常需要去判断所定义的变量是否是一个列表。所以这是一个常见的需求，可参考如下间接测试是否是列表的代码片段：

    (VAR is sequence) and (VAR is not string) and (VAR is not mapping)

如果大家以后深入到 Ansible 的开发方面，可以自定义 Ansible 的模块和插件，那么就可以自己写一个更完善的 Filter 插件来测试是否是 list。下面我简单演示下基本步骤，大家能依葫芦画瓢更好，不理解也没任何关系。

首先创建 filter_plugins 目录并在其内创建一个 py 文件，例如 collection.py，内容如下：

```python
def islist(collection):
  '''
  test data type is a list or not
  '''
  return isinstance(collection, list)

class FilterModule(object):
  '''
  custom jinja2 filter for test list type
  '''
  def filters(self):
    return {
      'islist': islist
    }
```

然后在 playbook 中便可使用 islist() 这个筛选器来判断是否是列表类型。例如：

        msg: "a list"
      when: p | islist
      vars:
        p:
          - p1
          - p2

<a name="GDuBh"></a>

## Filters(过滤器)

通常，模板语言都会带有过滤器，JinJa 也不例外，每个过滤器函数都是一个功能，作用就类似于函数，而且它也可以接参数。变量可以通过 **Filters(过滤器)** 修改。Jinja 中有两种使用 Filters 的方式：

- `**|**`** 符号** # 过滤器 与 变量 之间使用 `|` 符号分割，并且可以使用 `()` 符号传递参数。多个过滤器可以链式调用，前一个过滤器的返回值会作为有一个过滤器的输入。
- `**filter**`** 关键字 **# <a name="CLyDw"></a>

### `|` 符号

例如，Jinja 有一个内置 `lower()` 过滤器函数，可以将字符串全部转化成大写字母。

```yaml
- debug:
    msg: "{{'hello world'|upper()}}"
```

如果过滤器函数没有给定参数，则括号可以省略，例如 `"HELLO"|upper`。 <a name="eHdI7"></a>

### `filter` 关键字

我们还可以使用 filter 关键字，在语句表达式中使用过滤器，以对模板中的 一块数据(而不是一行或一个变量) 进行筛选操作，比如：

```python
{% filter upper %}
    这部分文本内容中，小写字母将会变成大写的
    This text becomes uppercase
{% endfilter %}
```

```yaml
{% if result %}
{{result|replace('no', 'yes')}}
{%endif%}
```

<a name="UKdta"></a>

### Jinja 内置过滤器

JinJa 内置了多个过滤器函数，Ansible 自身也扩展了一些方便的筛选器函数，所以数量非常多。如下：

| abs()            | float()       | lower()      | round()      | tojson()    |
| ---------------- | ------------- | ------------ | ------------ | ----------- |
| attr()           | forceescape() | map()        | safe()       | trim()      |
| batch()          | format()      | max()        | select()     | truncate()  |
| capitalize()     | groupby()     | min()        | selectattr() | unique()    |
| center()         | indent()      | pprint()     | slice()      | upper()     |
| default()        | int()         | random()     | sort()       | urlencode() |
| dictsort()       | join()        | reject()     | string()     | urlize()    |
| escape()         | last()        | rejectattr() | striptags()  | wordcount() |
| filesizeformat() | length()      | replace()    | sum()        | wordwrap()  |
| first()          | list()        | reverse()    | title()      | xmlattr()   |

我会将它们中绝大多数的含义列举出来 (剩下一部分是我觉得在 Ansible 中用不上的，比如 escape() 转义为 HTML 安全字符串)，各位没必要全都测试一遍，但是速看一遍并大概了解它们的含义和作用是有必要的。

- `float(default=0.0)` # 将数值形式的字符串转换为浮点数。如果无法转换，则返回默认值 0.0。可使用 default 参数自定义转换失败时的默认值。
  - 例如`"abcd"|float`、`""|float`都转换为 0.0，`""|float('NaN')`返回的是字符串 NaN，表示非数值含义。
- `int(default=0,base=10)` # 将数值形式的字符串直接截断为整数。如果无法转换，则返回默认值 0。可使用 default 参数自定义转换失败时的默认值。
  - 此外，还可以指定进制参数 base，比如 base=2 表示将传递过来的参数当作二进制进行解析，然后转换为 10 进制数值。
  - 例如`'3.55'|int`结果为 3，`'0b100'|int(base=2)`结果为 4。
- `abs()` # 计算绝对值。
  - 注意，只能计算数值，如果传递的是字符串，可使用筛选器 int() 或 float() 先转换成数值。例如`'-3.14'|float|abs`。
- `round(precision=0,method='common')` # 对数值进行四舍五入。第一个参数指定四舍五入的精度，第二个参数指定四舍五入的方式，有三种方式可选择：
  - ceil：只入不舍
  - floor：只舍不入
  - common：小于五的舍，大于等于 5 的入
  - 注意：
    - 只能计算数值，如果传递的是字符串，可使用筛选器 int() 或 float() 先转换成数值
    - 计算的是整数，则返回值是整数，计算的是浮点数，则返回值是浮点数
  - 例如`42.55|round`的结果为 43.0，`45|round`的结果是 45，`42.55|round(1,'floor')`的结果是 42.5。
- `random()` # 返回一个随机整数。竖线左边的值 X 决定了随机数的范围为`[0,X)`。
  - 例如`5|random`生成的随机数可能是 0、1、2、3、4。
- `list()` # 转换为列表。如果要转换的目标是字符串，则返回的列表是字符串中的每个字符。
  - 例如`range(1,4)|list`的结果是`[1,2,3]`，`"hey"|list`的结果是`["h","e","y"]`。
- `string()` # 转换为字符串。
  - 例如`"333aa"`结果为”333aa”。
- `tojson()` # 转换为 json 格式。
- `lower()`、`upper()`、`title()`、`capitalize()` # lower() 将大写字母转换为小写。upper() 将小写字母转换为大写。title() 将每个首字母转为大写。capitalize() 将第一个单词首字母转为大写。
- `min()`、`max()` # 从序列中取最小、最大值。
  - 例如`["a","abddd","cba"]|max`得到 cba。
- `sum(start=0)` # 计算序列中各元素的算术和。可指定 start 参数作为算术和的起点。
  - 例如`[1,2,3]|sum`得到 6，`[1,2,3]|sum(start=3)`得到 9。
- `trim()` # 移除字符串前缀和后缀空白。
  - 例如`"abcd"|trim ~ "DEF"`得到”abcdDEF”。
- `truncate()` # 截断字符串为指定长度。主要用于 web 编程，Ansible 用不到。
- `replace(old,new,count=None)` # 将字符串中的 old 替换成 new，count 参数指定替换多少次，默认替换所有匹配成功的。
  - 例如`"you see see you"|replace("see","look")`得到`you look look you`，而`replace("see","look",1)`则得到`you look see you`。
- `first()`、`last()` # 返回序列中的第一个、最后一个元素。
  - 例如`"hello world" | last`返回字母 d，`[2,3,4]|last`返回数值 4。
- `map(attribute='xxx')` # 如果一个列表中包含了多个 dict，map 可根据指定的属性名 (即 dict 的 key)，从列表中各 dict 内筛选出该属性值部分。
  - 例如，对于如下变量：


    p:
      - name: "junma"
        age: 23
      - name: woniu
        age: 22
        weight: 45
      - name: tuner
        age: 25
        weight: 50

- `p|map(attribute="name")|list`将得到`["junma","woniu","tuner"]`。
- `select()`、`reject()` # 从序列中选中、排除满足条件的项。例如：


    {{ numbers|select("odd") }}         ->选出奇数
    {{ numbers|select("even") }}        ->选出偶数
    {{ numbers|select("lt", 42) }}      ->选出小于42的数
    {{ strings|select("eq", "mystr") }} ->选出"mystr"元素
    {{ numbers|select("divisibleby", 3) }} ->选出被3整除的数

- 其中测试参数可以指定为支持的测试函数，在前文已经介绍过。
- `selectattr()`、`rejectattr()` # 根据对象属性筛选、排除序列中的多个元素。这个有时候很好用。比如：


    p:
      - name: "junma"
        age: 23
      - name: woniu
        age: 22
        weight: 45
      - name: tuner
        age: 25
        weight: 50

筛选所有 age 大于 22 岁的对象：

    p|selectattr('age','gt',22)|list

得到的结果：

    [
      {"age": 23,"name": "junma"},
      {"age": 25,"name": "tuner","weight": 50}
    ]

- `batch(count,fill_withs=None)` # 将序列中每 count 个元素打包成一个列表。最后一个列表可能元素个数不够，默认不填充，如果要填充，则指定`fill_with`参数。

  - 例如`[1,2,3,4,5]|batch(2)|list`得到`[[1,2],[3,4],[5]]`，`[1,2,3,4,5]|batch(2,'x')|list`得到`[[1,2],[3,4],[5,'x']]`。

- `default('default_value',bool=False)`或`d()` # 如果竖线左边的变量未定义，则返回 default() 指定的默认值。默认只对未定义变量其作用，如果想让 default() 也能对布尔类型的数据生效，需将第二个参数设置为 true。

  - `d()`是`default()`的简写方式。
  - 例如`myvar|d('undefined')`在 myvar 不存在时返回 undefined 字符串，`""|d("empty")`中因为是空字符串而不是未定义变量，所以仍然返回空字符串，`""|d("empty",true)`则返回 empty 字符串。

- `unique()` # 对序列中进行去重操作。
  - 例如`[1,2,3,3,1,2]|unique`得到结果`[1,2,3]`。

(22).`join(d="")`
将序列中的元素使用 d 参数指定的符号串联成字符串，默认连接符为空字符串。

例如`[1,2,3]|join("-")`得到`1-2-3`，`[1,2,3]|join`得到 123。

(23).`length()`和`count()`
返回序列中元素的数量或字符串的字符个数。length 和 count 是别名等价的关系。

(24).`wordcount`
计算字符串中的单词个数。

(25).`reverse()`
颠倒序列元素。

例如`"hello"|reverse`得到`olleh`，`[1,2,3]|reverse|list`得到`[3,2,1]`。

(26).`filesizeformat(binary=False)`
将数值大小转换为 kB、MB、GB 单位。默认按照 1000 为单位进行换算，如果指定 binary 参数为 True，则按 1024 进行换算。

(27).`slice(N, fill_with=None)`
将序列均分为 N 个列表，可指定`fill_with`参数来填充均分时不足的列表。

例如：

    [1,2,3,4,5]|slice(3)    |list -> [[1,2],[3,4],[5]]
    [1,2,3,4]  |slice(3,"x")|list -> [[1,2],[3,"x"],[4,"x"]]

(28).`groupby(attribute)`
根据指定的属性对 dict 进行分组。看下面的例子：

```yaml
- debug:
    msg: '{{person|groupby("name")}}'
  vars:
    person:
      - name: "junma"
        age: 23
      - name: "junma"
        age: 33
      - name: woniu
        age: 22
        weight: 45
      - name: tuner
        age: 25
        weight: 50
```

得到的结果：

    [
      ["junma", [{"age": 23,"name": "junma"},{"age": 33,"name": "junma"}]],
      ["tuner", [{"age": 25,"name": "tuner","weight": 50}]],
      ["woniu", [{"age": 22,"name": "woniu","weight": 45}]]
    ]

(29).`indent(width=4)`
对字符串进行缩进格式化。默认缩进宽度为 4。

(30).`format()`
类似于`printf`的方式格式化字符串。对于熟悉 Python 的人来说，使用`%`或`str.format()`格式化字符串要更方便些。

下面是等价的：

    {{ "%s, %s!"|format("hello", "world") }}
    {{ "%s, %s!" % ("hello", "world") }}
    {{ "{}, {}!".format("hello", "world") }}

(31).`center(witdth=80)`
将字符串扩充到指定宽度并放在中间位置，默认 80 个字符。

例如`"hello"|center(10)|replace("","-")`得到`"--hello---"`。

(32).`sort(reverse=False, case_sensitive=False, attribute=None)`
对序列中的元素进行排序，attribute 指定按照什么进行排序。

例如`["acd","a","ca"]|sort`得到`["a","acd","ca"]`。

再例如，对于下面的变量：

    person:
      - name: "junma"
        age: 23
      - name: woniu
        age: 22
        weight: 45
      - name: tuner
        age: 25
        weight: 50

可以使用`person|sort(attribute="age")`对 3 个元素按照 age 进行升序排序。

(33).`dictsort(case_sensitive=False, by='key', reverse=False)`
对字典进行排序。默认按照 key 进行排序，可以指定为`value`按照 value 进行排序。

例如：

    person:
      p2:
        name: "junma"
        age: 23
      p1:
        name: woniu        age: 22
        weight: 45
      p3:
        age: 25
        weight: 50

可以使用`person|dictsort`按照 key(即 p1、p2、p3) 进行排序，结果是先 p1，再 p2，最后 p3。 <a name="f5abdb0e"></a>

## 循环

<a name="4c83c1c0"></a>

### for 迭代列表

for 循环的语法：

```python
{% for i in LIST %}
    string_or_expression
{% endfor %}
```

还支持直接条件判断筛选要参与迭代的元素：

```python
{% for i in LIST if CONDITION %}
    string_or_expression
{% endfor %}
```

此外，Jinja2 的 for 语句还允许使用 else 分支，如果 for 所迭代的列表 LIST 是空列表 (或没有元素可迭代)，则会执行 else 分支。

    {% for i in LIST %}
        string_or_expression
    {% else %}
        string_or_expression
    {% endfor %}

例如，在模板文件 a.txt.j2 中有如下内容：

    {% for file in files %}
    <{{file}}>
    {% else %}
    no file in files
    {% endfor %}

playbook 文件内容如下：

```yaml
---
- hosts: localhost
  gather_facts: no
  tasks:
    - template:
        src: a.txt.j2
        dest: /tmp/a.txt
      vars:
        files:
          - /tmp/a1
          - /tmp/a2
          - /tmp/a3
```

执行 playbook 之后，将生成包含如下内容的 / tmp/a.txt 文件：

    </tmp/a1>
    </tmp/a2>
    </tmp/a3>

如果将 playbook 中的`files`变量设置为空列表，则会执行 else 分支，所以生成的 / tmp/a.txt 的内容为：

    no file in files

如果 files 变量未定义或变量类型不是 list，则默认会报错。针对未定义变量，可采用如下策略提供默认空列表：

    {% for file in (files|default([])) %}
    <{{file}}>
    {% else %}
    no file in files
    {% endfor %}

如果不想迭代文件列表中的`/tmp/a3`，则可以加上条件判断：

    <{{file}}>
    {% else %}
    no file in files
    {% endfor %}

Jinja2 的 for 循环没有提供 break 和 continue 的功能，所以只能通过 {% for...if...%} 来间接实现类似功能。

<a name="adcb82d3"></a>

### for 迭代字典

默认情况下，Jinja2 的 for 语句只能迭代列表。

如果要迭代字典结构，需要先使用字典的`items()`方法进行转换。如果没有学过 python，我下面做个简单解释：

对于下面的字典结构：

    p:
      name: junmajinlong
      age: 18

如果使用`p.items()`，将计算得到如下结果：

    [('name', 'junmajinlong'), ('age', 18)]

然后 for 语句中使用两个迭代变量分别保存各列表元素中的子元素即可。下面设置了两个迭代变量 key 和 value：

    {% for key,value in p.items() %}

那么第一轮迭代时，key 变量保存的是 name 字符串，value 变量保存的是 junmajinlong 字符串，那么第二轮迭代时，key 变量保存的是 age 字符串，value 变量保存的是 18 数值。

如果 for 迭代时不想要 key 或不想要 value，则使用`_`来丢弃对应的值。也可以使用`keys()`方法和`values()`方法分别获取字典的 key 组成的列表、字典的 value 组成的列表。例如：

    {% for key,_ in p.items() %}
    {% for _,values in p.items() %}
    {% for key in p.keys() %}
    {% for value in p.values() %}

将上面的解释整理成下面的示例。playbook 内容如下：

```yaml
- hosts: localhost
  gather_facts: no
  tasks:
    - template:
        src: a.txt.j2
        dest: /tmp/a.txt
      vars:
        p1:
          name: "junmajinlong"
          age: 18
```

模板文件 a.txt.j2 内容如下：

    {% for key,value in p1.items() %}
    {% endfor %}

执行结果：

    key: name, value: junmajinlong
    key: age, value: 18

<a name="1171d6bf"></a>

### for 的特殊控制变量

在 for 循环内部，可以使用一些特殊变量，如下：

| Variable            | Description                                                            |
| ------------------- | ---------------------------------------------------------------------- |
| loop                | 循环本身                                                               |
| loop.index          | 本轮迭代的索引位，即第几轮迭代 (从 1 开始计数)                         |
| loop.index0         | 本轮迭代的索引位，即第几轮迭代 (从 0 开始计数)                         |
| loop.revindex       | 本轮迭代的逆向索引位 (距离最后一个 item 的长度，从 1 开始计数)         |
| loop.revindex0      | 本轮迭代的逆向索引位 (距离最后一个 item 的长度，从 0 开始计数)         |
| loop.first          | 如果本轮迭代是第一轮，则该变量值为 True                                |
| loop.last           | 如果本轮迭代是最后一轮，则该变量值为 True                              |
| loop.length         | 循环要迭代的轮数，即 item 的数量                                       |
| loop.previtem       | 本轮迭代的前一轮的 item 值，如果当前是第一轮，则该变量未定义           |
| loop.nextitem       | 本轮迭代的下一轮的 item 值，如果当前是最后一轮，则该变量未定义         |
| loop.depth          | 在递归循环中，表示递归的深度，从 1 开始计数                            |
| loop.depth0         | 在递归循环中，表示递归的深度，从 0 开始计数                            |
| loop.cycle          | 一个函数，可指定序列作为参数，for 每迭代一次便同步迭代序列中的一个元素 |
| loop.changed(\*val) | 如果本轮迭代时的 val 值和前一轮迭代时的 val 值不同，则返回 True        |

之前曾介绍过，在 Ansible 的循环开启`extended`功能之后也能获取一些特殊变量。不难发现，Ansible 循环开启`extended`后可获取的变量和此处 Jinja2 提供的循环变量大多是类似的。所以这里只介绍之前尚未解释过的几个变量。

首先是`loop.cycle()`，它是一个函数，可以传递一个序列 (比如列表) 作为参数。在 for 循环迭代时，每迭代一个元素的同时，也会从参数指定的序列中迭代一个元素，如果序列元素迭代完了，则从头开始继续迭代。

例如，playbook 内容如下：

    - hosts: localhost
      gather_facts: no
      tasks:
        - template:
            src: a.txt.j2
            dest: /tmp/a.txt
          vars:
            p:
              - aaa
              - bbb
              - ccc

模板文件 a.txt.j2 内容如下：

    {% for i in p %}
    item: {{i}}
    cycle: {{loop.cycle("AAA","BBB")}}
    {% endfor %}

渲染后得到的 / tmp/a.txt 文件内容如下：

    item: aaa
    cycle: AAA
    item: bbb
    cycle: BBB
    item: ccc
    cycle: AAA

然后是`loop.changed(val)`，这也是一个函数。如果相邻的两轮迭代中 (即当前一轮和前一轮)，参数 val 的值没有发生变化，则当前一轮的`loop.changed()`返回 False，否则返回 True。举个例子很容易理解：

playbook 内容如下：

```yaml
- hosts: localhost
  gather_facts: no
  tasks:
    - template:
        src: a.txt.j2
        dest: /tmp/a.txt
      vars:
        persons:
          - name: "junmajinlong"
            age: 18
          - name: "junmajinlong"
            age: 22
          - name: "wugui"
            age: 23
```

模板文件 a.txt.j2 内容如下：

    {% for p in persons %}
    {% if loop.changed(p.name) %}
    index: {{loop.index}}
    {% endif %}    {% endfor %}

渲染后得到的 / tmp/a.txt 结果：

    index: 1
    index: 3

显然，第二轮迭代时的 p.name 和前一轮迭代时的 p.name 值是相同的，所以渲染结果中没有`index: 2`。 <a name="rqns6"></a>

### 如何跨作用域

那如何在 for 循环内做一个自增操作呢？这应该也是非常常见的需求。但只能说 Jinja2 里这不方便，只能退而求其次找其它方式，这里我提供两种：

    {# 使用loop.index，它本身就是自增的 #}
    {% set mylist = [1,2,3] %}
    {% for item in mylist %}
    name{{loop.index}}
    {% endfor %}

    {# 使用Jinja2 2.10版的namespace，它可以让变量跨作用域 #}
    {% set num = namespace(value=3) %}
    {% for item in mylist %}
    name{{num.value}}
    {% set num.value = num.value + 2 %}
    {% endfor %}

使用上面第二种方案时要注意 Jinja2 的版本号，Ansible 所使用的 Jinja2 很可能是低于 2.10 版本的。

<a name="W4CXG"></a>

## Macro(宏)

计算机科学当中，Macro(宏) 表示的是一段指令的简写，它会在特定的时候替换成其所代表的一大段指令。

如果各位之前不曾知道 Macro 的概念，我这里用一个不严谨、不属于同一个范畴但最方便大家理解的示例来解释：Shell 中的命令别名可以看作是 Macro，Shell 会在命令开始执行之前 (即在 Shell 解析命令行的阶段) 先将别名替换成其原本的值。比如将`ls`替换成`rm -rf`。

Jinja2 是支持 Macro 概念的，宏类似于函数，比如可以接参数，具有代码复用的功能。但其本质和函数是不同的，Macro 是替换成原代码片段，函数是直接寻找到地址进行调用。这就不多扯了，好像离题有点远，这可不是编程教程。总的来说，Jinja2 的 Macro 需要像函数一样先定义，在使用的时候直接调用即可，至于底层细节，管它那么多干嘛，又不会涨一毛钱工资。

举一个比较常见的案例，比如某服务的配置文件某指令可以接多个参数值，每个值以空格分隔，每个指令以分号`;`结尾。例如：`log 'host' 'port' 'date';`。如果用模板去动态配置这个指令，可能会使用 for 循环迭代，但要区分空格分隔符和尾部的分号分隔符。于是，编写如下 Macro：

    {% macro delimiter(loop) -%}
    {{ ' ' if not loop.last else ';' }}
    {%- endmacro %}

上面表示定义了一个名为 delimiter 的 Macro，它能接一个表示 for 循环的参数。

上面的 Macro 定义中还使用了 -%} 和 {%- ，这是用于处理空白符号的，稍后会解释它的用法，现在各位只需当这个短横线不存在即可。

定义好这个 Macro 之后，就可以在任意需要的时候去” 调用” 它。例如：

    log {% for item in log_args %}
    '{{item}}'{{delimiter(loop)}}
    {%- endfor %}

    gzip {% for item in gzip_args %}
    '{{item}}'{{delimiter(loop)}}
    {%- endfor %}

提供一个 playbook，内容如下：

```yaml
- hosts: localhost
  gather_facts: no
  tasks:
    - template:
        src: a.txt.j2
        dest: /tmp/a.txt
      vars:
        log_args:
          - host
          - port
          - date
        gzip_args: ["css", "js", "html"]
```

渲染出来的结果如下：

    log 'host' 'port' 'date';
    gzip 'css' 'js' 'html';

Macro 的参数还可以指定默认值，” 调用”Macro 并传递参数时，还可以用 key=value 的方式传递。例如：

    {# 定义Macro时，指定默认值 #}
    {% macro delimiter(loop,sep=" ",deli=";") -%}
    {{ sep if not loop.last else deli }}
    {%- endmacro %}

    {# "调用"Macro时，使用key=value传递参数值 #}
    log {% for item in log_args %}
    '{{item}}'{{delimiter(loop,sep=",")}}
    {%- endfor %}

    gzip {% for item in gzip_args %}
    '{{item}}'{{delimiter(loop,deli="")}}
    {%- endfor %}

渲染得到的结果：

    log 'host','port','date';
    gzip 'css' 'js' 'html'

关于 Macro，还有些内容可以继续深入 (一些变量和 call 调用的方式)，但应该很少很少用到，所以我这就不再展开了，如果大家有意愿，可以去官方手册学习或网上搜索相关资料，有编程基础的人应该很容易理解，没有编程基础的，就别凑这个热闹了。

<a name="spWbH"></a>

## Block(块) 与 Extends(继承)

有些服务程序的配置文件可以使用 include 指令来包含额外的配置文件，这样可以按不同功能来分类管理配置文件中的配置项。在解析配置文件的时候，会将 include 指令所指定的文件内容加载并合并到主配置文件中。

Jinja2 的 block 功能有点类似于 include 指令的功能，block 的用法是这样的：先在一个类似于主配置文件的文件中定义 block，称为 base block 或父 block，然后在其它文件中继承 base block，称为子 block。在模板解析的时候，会将子 block 中的内容填充或覆盖到父 block 中。

例如，在 base.conf.j2 文件中定义如下内容：

    server {
      listen       80;
      server_name  www.abc.com;

    {% block root_page %}
    location / {
      root   /usr/share/nginx/html;
      index  index.html index.htm;
    }
    {% endblock root_page %}

      error_page   500 502 503 504  /50x.html;
    {% block err_50x %}{% endblock err_50x %}
    {% block php_pages %}{% endblock php_pages %}

    }

这其实是一个 Nginx 的虚拟主机配置模板文件。在其中定义了三个 block：

- (1). 名为 root_page 的 block，其内部有内容，这个内容是默认内容
- (2). 名为 err_50x 的 block，没有内容
- (3). 名为 php_pages 的 block，没有内容

如果定义了同名子 block，则会使用子 block 来覆盖父 block，如果没有定义同名子 block，则会采用默认内容。

下面专门用于定义子 block 内容的 child.conf.j2 文件，内容如下：

    {% extends 'base.conf.j2' %}

    {% block err_50x %}
    location = /50x.html {
      root   /usr/share/nginx/html;
    }
    {% endblock err_50x %}

      {% block php_pages %}
      location ~ \.php$ {
        fastcgi_pass   "192.168.200.43:9000";
        fastcgi_index  index.php;
        fastcgi_param  SCRIPT_FILENAME /usr/share/www/php$fastcgi_script_name;
        include        fastcgi_params;
      }
      {% endblock php_pages %}

子 block 文件中第一行需要使用 jinja2 的`extends`标签来指定父 block 文件。这个子 block 文件中，没有定义名为`root_page`的 block，所以会使用父 block 文件中同名 block 的默认内容，`err_50x`和`php_pages`则直接覆盖父 block。

在 template 模块渲染文件时，需要指定子 block 作为其源文件。例如：

```yaml
- hosts: localhost
  gather_facts: no
  tasks:
    - template:
        src: child.conf.j2
        dest: /tmp/nginx.conf
```

渲染得到的结果:

```nginx
server {
  listen       80;
  server_name  www.abc.com;

location / {
  root   /usr/share/nginx/html;
  index  index.html index.htm;
}

  error_page   500 502 503 504  /50x.html;
location = /50x.html {
  root   /usr/share/nginx/html;
}
  location ~ \.php$ {
    fastcgi_pass   "192.168.200.43:9000";
    fastcgi_index  index.php;    fastcgi_param  SCRIPT_FILENAME /usr/share/www/php$fastcgi_script_name;
    include        fastcgi_params;
  }

}
```

jinja2 的 block 是很出色的一个功能，但在 Ansible 中应该不太可能用到 (或机会极少)，所以多的就不介绍了，有兴趣的可自行找资料了解。 <a name="xaJjJ"></a>

# Jinja 的空白处理

通常在模板文件中，会将模板代码片段按照编程语言的代码一样进行换行、缩进，但因为它们是嵌套在普通字符串中的，模板引擎并不知道那是一个普通字符串中的空白还是代码格式规范化的空白，而有时候这会带来问题。

比如，模板文件 a.txt.j2 文件中的内容如下：

    line start
    line left {% if true %}
      <line1>
    {% endif %} line right
    line end

这个模板文件中的代码部分看上去非常规范，有换行有缩进。一般来说，这段模板文件想要渲染得到的文本内容应该是：

    line start
    line left
    <line1>
    line right
    line end

或者是：

    line start
    line left <line1> line right
    line end

但实际渲染得到的结果：

    line start
    line left   <line1>
     line right
    line end

渲染的结果中格式很不规范，主要原因就是 Jinja2 语句块前后以及语句块自身的换行符处理、空白符号处理导致的问题。

Jinja2 提供了两个配置项：`lstrip_blocks`和`trim_blocks`，它们的意义分别是：

- (1).`lstrip_blocks`：设置为 true 时，会将 Jinja2 语句块前面的本行前缀空白符号移除
- (2).`trim_blocks`：设置为 true 时，Jinja2 语句块后的换行符会被移除掉

对于 Ansible 的 template 模块，`lstrip_blocks` 默认设置为 False，`trim_blocks` 默认设置为 true。也就是说，默认情况下，template 模块会将语句块后面的换行符移除掉，但是会保留语句块前的本行前缀空白符号。

例如，对于下面这段模板片段：

    line start
        {% if true %}
      <line1>
    {% endif %}
    line end

`{% if`前的 4 个空格会保留，`true %}`后的换行符会被移除，于是'line1'(注意前面两个空格) 渲染的时候会移到第二行去。再看`{% endif %}`，`{%`前面的空白符号会保留，`%}`后面的换行符会被移除，所以`line end`在渲染时会移动到第三行。第二行和第三行的换行符是由'line1'这行提供的。

所以结果是：

    line start
          <line1>
    line end

一般来说，将`lstrip_blocks`和`trim_blocks`都设置为 true，比较符合大多数情况下的空白处理需求。例如：

```yaml
- template:
    src: a.txt.j2
    dest: /tmp/a.txt
    lstrip_blocks: true
    trim_blocks: true
```

渲染得到的结果：

    line start
      <line1>
    line end

更符合一般需求的模板格式是，**Jinja2 指令部分 (比如 if、endif、for、endfor 等) 不要使用任何缩进格式，非 Jinja2 指令部分按需缩进**。

    line start
    {% if true %}
      <line1>
    {% endif %}
    line end

除了`lstrip_blocks`以及`trim_blocks`可以控制空白外，还可以使用`{%- xxx`可以移除本语句块前的所有空白 (包括换行符)，使用`-%}`可以移除本语句块后的所有空白 (包括换行符)。

注意，`xxx_blocks`这两个配置项和带`-`符号的效果是不同的，总结下：

- (1).lstrip_blocks 只移除语句块前紧连着的且是本行前缀的空白符
- (2).{%- 移除语句块前所有空白
- (3).trip_blocks 只移除语句块后紧跟着的换行符
- (4).-%}\`移除语句块后所有的空白

例如，下面两个模板片段：

    line1
    line2 {%- if true %}
            line3
            line4
          {%- endif %}
    line44
    line5

在两个`xxx_blocks`设置为 true 时，渲染得到的结果是：
line1
line2line3
line4line44
line5

最后想告诉各位，如果渲染后得到的结果是合理的 (比如配置文件语法不报错)，就不要追求精确控制空白符号。比如别为了将多个连续的空格压缩成看上去更显规范的单个空格而想方设法(如果你是强迫症，就要小心咯)。如果你还没遇到过这个问题，那以后也肯定会遇到的，其实只要模板稍微写的复杂一点，就能体会到什么叫做” 众口难调”。

<a name="u40nx"></a>

在前面解释过，当使用`x.y`的方式访问 y 的时候，会先寻找 x 的 y 属性或 y 方法，找不到才开始找 Jinja2 中定义的属性。

所以，在 Ansible 中有些时候是可以直接使用 Python 对象自身方法的，比如字符串对象可以使用`endswith`判断字符串是否以某字符串结尾。

这也为 Ansible 提供了非常有用的功能。但是有些人可能没学过 Python，所以也不知道有哪些方法可用，也不理解有些代码是什么作用。这一点我也没有办法帮助各位，但大家也不用太过在意，几个方法而已，细节罢了。事实上也就字符串对象的方法比较多。

<a name="cc8ed0c9"></a>

### Python 字符串处理

在 Jinja2 中，经常会使用到字符串。如何使用字符串对象的方法？

例如，Python 字符串对象有一个 upper 方法，可以将字符串改变为大写字母，直接使用`"abc".upper()`，注意不要省略小括号，这一点和 Jinja2 和 Shell 函数都是不一样的。

例如：

```yaml
- debug:
    msg: '{{ "abc".upper() }}'
- debug:
    msg: "{{ 'foo bar baz'.upper().split() }}"
```

得到：

    TASK [debug] ****************
    ok: [localhost] => {
        "msg": "ABC"
    }

    TASK [debug] *****************
        "msg": ["FOO", "BAR", "BAZ"]
    }

下面是字符串对象的各种方法，我简单说明了它们的功能，关于它们的用法和示例可参见我的博客：<https://www.junmajinlong.com/python/string_methods/>。

    lower：将字符串转换为小写字母
    upper：将字符串转换为大写字母
    title：将字符串所有单词首字母大写，其它字母小写
    capitalize：将字符串首单词首字母大写，其它字母小写
    swapcase：将字符串中所有字母大小写互换
    isalpha：判断字符串中所有字符是否是字母
    isdecimal：判断字符串中所有字符是否是数字
    isdigit：判断字符串中所有字符是否是数字
    isnumeric：判断字符串中所有字符是否是数字
    isalnum：判断字符串中所有字符是否是字母或数字
    islower：判断字符串中所有字符是否是小写字母
    isupper：判断字符串中所有字符是否是大写字母
    istitle：判断字符串中是否所有单词首字符大写，其它小写
    isspace：判断字符串中所有字符是否是空白字符(空格、制表符、换行符等)
    isprintable：判断字符串中所有字符是否是可打印字符(如制表符、换行符不是可打印字符)
    isidentifier：判断字符串中是否符合标识符定义规则(即只包含字母、数字或下划线，且字母或下划线开头)
    center：在左右两边填充字符串到指定长度，字符串居中
    ljust：在右边填充字符串到指定长度
    rjust：在左边填充字符串到指定长度
    zfill：使用0填充在字符串左边
    count：计算字符串中某字符或某子串出现的次数
    endswith：字符串是否以某字符或某子串结尾
    startswith：字符串是否以某字符或某子串开头
    find：从左向右检查字符串中是否包含某子串，搜索不到时返回-1
    rfind：从右向左检查字符串中是否包含某子串，搜索不到时返回-1
    index：功能类似于find，搜索不到时报错
    rindex：功能类似于rfind，搜索不到时报错
    replace：替换字符串
    expandtabs：将字符串中的制表符\t替换成空格，默认替换为8空格
    split：将字符串分割得到列表
    splitlines：按行分割字符串，每行作为列表的一个元素
    join：将列表各元素使用指定符号连接起来，例如`"-".[1,2,3]`得到`1-2-3`
    strip：移除字符串前缀和后缀指定字符，如果没有指定移除的字符，则默认移除空白
    lstrip：移除字符串指定的前缀字符，如果没有指定移除的字符，则默认移除空白
    rstrip：移除字符串指定的后缀字符，如果没有指定移除的字符，则默认移除空白
    format：格式化字符串，此外还可使用Python的`%`格式化方式，如{{ '%.2f' % 1.2345 }}

<a name="a148a52f"></a>

### list 对象方法

虽然 Python 中 list 对象有很多操作方式，但应用到 Ansible 中，大概也就两个个方法值得了解：

- `count()`：(无参数) 计算列表中元素的个数或给定参数在列表中出现的次数
- `index()`：检索给定参数在列表中的位置，如果元素不存在，则报错

例如：

```yaml
- debug:
    msg: "{{ [1,2,1,1,3,4].count() }}"
- debug:
    msg: "{{ [1,2,1,1,3,4].count(1) }}"
- debug:
    msg: "{{ [1,2,1,1,3,4].index(2) }}"
```
