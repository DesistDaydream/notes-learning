---
title: Filters Plugins(过滤插件)
---

# 概述

> 参考：
> - [官方文档，用户指南-目录-使用过滤器操作数据](https://docs.ansible.com/ansible/latest/user_guide/playbooks_filters.html)
> - <https://www.zsythink.net/archives/2862>

在本博客中，ansible 是一个系列文章，我们会尽量以通俗易懂的方式总结 ansible 的相关知识点。

ansible 系列博文直达链接：ansible 轻松入门系列

现在我有一个需求，我想要将获取到的变量的值中的所有字母都变成大写，如果想要在 playbook 中实现这个需求，我该怎么办呢？我可以借助一个叫做"过滤器"的东西，帮助我完成刚才的需求，"过滤器（filters）"可以帮助我们对数据进行处理，这样解释可能不够直观，不如这样，我们先来看一个过滤器的小例子，然后结合示例解释过滤器是个什么东西，示例如下：

    - hosts: test70
      remote_user: root
      gather_facts: no
      vars:
        testvar: 1a2b3c
      tasks:
      - debug:
          msg: "{{ testvar | upper }}"

如上例所示，testvar 变量的值中包含三个小写字母，在使用 debug 模块输出这个变量的值时，我们使用了一个管道符，将 testvar 变量传递给了一个名为"upper"的东西，"upper"就是一个"过滤器"，执行上例 playbook 后你会发现，testvar 中的所有小写字母都被变成了大写。

通过上述示例，你一定已经明白了，过滤器是一种能够帮助我们处理数据的工具，其实，ansible 中的过滤器功能来自于 jinja2 模板引擎，我们可以借助 jinja2 的过滤器功能在 ansible 中对数据进行各种处理，而上例中的 upper 就是一种过滤器，这个过滤器的作用就是将小写字母变成大写，你一定已经发现了，当我们想要通过过滤器处理数据时，只需要将数据通过管道符传递给对应的过滤器即可，当然，过滤器不只有 upper，还有很多其他的过滤器，这些过滤器有些是 jinja2 内置的，有些是 ansible 特有的，如果这些过滤器都不能满足你的需求，jinja2 也支持自定义过滤器。

这篇文章我们就来总结一些常用的过滤器的用法，在总结时，不会区分它是 jinja2 内置的过滤器，还是 ansible 所独有的，我们总结的目的是在 ansible 中使用这些过滤器，如果你想要了解 jinja2 中有哪些内置过滤器，可以参考 jinja2 的官网链接，如下

<http://jinja.pocoo.org/docs/2.10/templates/#builtin-filters>

# 字符串操作有关的过滤器

    - hosts: test70
      remote_user: root
      vars:
        testvar: "abc123ABC 666"
        testvar1: "  abc  "
        testvar2: '123456789'
        testvar3: "1a2b,@#$%^&"
      tasks:
      - debug:
          #将字符串转换成纯大写
          msg: "{{ testvar | upper }}"
      - debug:
          #将字符串转换成纯小写
          msg: "{{ testvar | lower }}"
      - debug:
          #将字符串变成首字母大写,之后所有字母纯小写
          msg: "{{ testvar | capitalize }}"
      - debug:
          #将字符串反转
          msg: "{{ testvar | reverse }}"
      - debug:
          #返回字符串的第一个字符
          msg: "{{ testvar | first }}"
      - debug:
          #返回字符串的最后一个字符
          msg: "{{ testvar | last }}"
      - debug:
          #将字符串开头和结尾的空格去除
          msg: "{{ testvar1 | trim }}"
      - debug:
          #将字符串放在中间，并且设置字符串的长度为30，字符串两边用空格补齐30位长
          msg: "{{ testvar1 | center(width=30) }}"
      - debug:
          #返回字符串长度,length与count等效,可以写为count
          msg: "{{ testvar2 | length }}"
      - debug:
          #将字符串转换成列表，每个字符作为一个元素
          msg: "{{ testvar3 | list }}"
      - debug:
          #将字符串转换成列表，每个字符作为一个元素，并且随机打乱顺序
          #shuffle的字面意思为洗牌
          msg: "{{ testvar3 | shuffle }}"
      - debug:
          #将字符串转换成列表，每个字符作为一个元素，并且随机打乱顺序
          #在随机打乱顺序时，将ansible_date_time.epoch的值设置为随机种子
          #也可以使用其他值作为随机种子，ansible_date_time.epoch是facts信息
          msg: "{{ testvar3 | shuffle(seed=(ansible_date_time.epoch)) }}"

跟数字操作有关的过滤器，示例如下

| 1
2
3
4
5
6
7
8
9
10
11
12
13
14
15
16
17
18
19
20
21
22
23
24
25
26
27
28
29
30
31
32
33
34
35
36
37
38
39
40
41
42
43
44
45
46
47
48 | ---
\- hosts: test70
&#x20; remote_user: root
&#x20; vars:
&#x20; testvar4: -1
&#x20; tasks:
&#x20; \- debug:
&#x20; \#将对应的值转换成 int 类型
&#x20; \#ansible 中，字符串和整形不能直接计算，比如{{ 8+'8' }}会报错
&#x20; \#所以，我们可以把一个值为数字的字符串转换成整形后再做计算
&#x20; msg: "{{ 8+('8' | int) }}"
&#x20; \- debug:
&#x20; \#将对应的值转换成 int 类型,如果无法转换,默认返回 0
&#x20; \#使用 int(default=6)或者 int(6)时，如果无法转换则返回指定值 6
&#x20; msg: "{{ 'a' | int(default=6) }}"
&#x20; \- debug:
&#x20; \#将对应的值转换成浮点型，如果无法转换，默认返回'0.0'
&#x20; msg: "{{ '8' | float }}"
&#x20; \- debug:
&#x20; \#当对应的值无法被转换成浮点型时，则返回指定值’8.8‘
&#x20; msg: "{{ 'a' | float(8.88) }}"
&#x20; \- debug:
&#x20; \#获取对应数值的绝对值
&#x20; msg: "{{ testvar4 | abs }}"
&#x20; \- debug:
&#x20; \#四舍五入
&#x20; msg: "{{ 12.5 | round }}"
&#x20; \- debug:
&#x20; \#取小数点后五位
&#x20; msg: "{{ 3.1415926 | round(5) }}"
&#x20; \- debug:
&#x20; \#从 0 到 100 中随机返回一个随机数
&#x20; msg: "{{ 100 | random }}"
&#x20; \- debug:
&#x20; \#从 5 到 10 中随机返回一个随机数
&#x20; msg: "{{ 10 | random(start=5) }}"
&#x20; \- debug:
&#x20; \#从 5 到 15 中随机返回一个随机数,步长为 3
&#x20; \#步长为 3 的意思是返回的随机数只有可能是 5、8、11、14 中的一个
&#x20; msg: "{{ 15 | random(start=5,step=3) }}"
&#x20; \- debug:
&#x20; \#从 0 到 15 中随机返回一个随机数,这个随机数是 5 的倍数
&#x20; msg: "{{ 15 | random(step=5) }}"
&#x20; \- debug:
&#x20; \#从 0 到 15 中随机返回一个随机数，并将 ansible_date_time.epoch 的值设置为随机种子
&#x20; \#也可以使用其他值作为随机种子，ansible_date_time.epoch 是 facts 信息
&#x20; \#seed 参数从 ansible2.3 版本开始可用
&#x20; msg: "{{ 15 | random(seed=(ansible\_date\_time.epoch)) }}" |
| --- | --- |

列表操作相关的过滤器，示例如下

Shell

|     |     |
| --- | --- |

| 1
2
3
4
5
6
7
8
9
10
11
12
13
14
15
16
17
18
19
20
21
22
23
24
25
26
27
28
29
30
31
32
33
34
35
36
37
38
39
40
41
42
43
44
45
46
47
48
49
50
51
52
53
54
55
56
57
58
59
60
61
62
63
64
65
66
67
68
69
70
71
72
73
74
75
76
77
78
79
80
81
82
83
84
85
86
87
88
89
90
91
92
93
94
95 | ---
\- hosts: test70
&#x20; remote_user: root
&#x20; vars:
&#x20; testvar7: \[22,18,5,33,27,30]
&#x20; testvar8: \[1,\[7,2,\[15,9]],3,5]
&#x20; testvar9: \[1,'b',5]
&#x20; testvar10: \[1,'A','b',\['QQ','wechat'],'CdEf']
&#x20; testvar11: \['abc',1,3,'a',3,'1','abc']
&#x20; testvar12: \['abc',2,'a','b','a']
&#x20; tasks:
&#x20; \- debug:
&#x20; \#返回列表长度,length 与 count 等效,可以写为 count
&#x20; msg: "{{ testvar7 | length }}"
&#x20; \- debug:
&#x20; \#返回列表中的第一个值
&#x20; msg: "{{ testvar7 | first }}"
&#x20; \- debug:
&#x20; \#返回列表中的最后一个值
&#x20; msg: "{{ testvar7 | last }}"
&#x20; \- debug:
&#x20; \#返回列表中最小的值
&#x20; msg: "{{ testvar7 | min }}"
&#x20; \- debug:
&#x20; \#返回列表中最大的值
&#x20; msg: "{{ testvar7 | max }}"
&#x20; \- debug:
&#x20; \#将列表升序排序输出
&#x20; msg: "{{ testvar7 | sort }}"
&#x20; \- debug:
&#x20; \#将列表降序排序输出
&#x20; msg: "{{ testvar7 | sort(reverse=true) }}"
&#x20; \- debug:
&#x20; \#返回纯数字非嵌套列表中所有数字的和
&#x20; msg: "{{ testvar7 | sum }}"
&#x20; \- debug:
&#x20; \#如果列表中包含列表，那么使用 flatten 可以'拉平'嵌套的列表
&#x20; \#2.5 版本中可用,执行如下示例后查看效果
&#x20; msg: "{{ testvar8 | flatten }}"
&#x20; \- debug:
&#x20; \#如果列表中嵌套了列表，那么将第 1 层的嵌套列表‘拉平’
&#x20; \#2.5 版本中可用,执行如下示例后查看效果
&#x20; msg: "{{ testvar8 | flatten(levels=1) }}"
&#x20; \- debug:
&#x20; \#过滤器都是可以自由结合使用的，就好像 linux 命令中的管道符一样
&#x20; \#如下，取出嵌套列表中的最大值
&#x20; msg: "{{ testvar8 | flatten | max }}"
&#x20; \- debug:
&#x20; \#将列表中的元素合并成一个字符串
&#x20; msg: "{{ testvar9 | join }}"
&#x20; \- debug:
&#x20; \#将列表中的元素合并成一个字符串,每个元素之间用指定的字符隔开
&#x20; msg: "{{ testvar9 | join(' , ') }}"
&#x20; \- debug:
&#x20; \#从列表中随机返回一个元素
&#x20; \#对列表使用 random 过滤器时，不能使用 start 和 step 参数
&#x20; msg: "{{ testvar9 | random }}"
&#x20; \- debug:
&#x20; \#从列表中随机返回一个元素,并将 ansible_date_time.epoch 的值设置为随机种子
&#x20; \#seed 参数从 ansible2.3 版本开始可用
&#x20; msg: "{{ testvar9 | random(seed=(ansible\_date\_time.epoch)) }}"
&#x20; \- debug:
&#x20; \#随机打乱顺序列表中元素的顺序
&#x20; \#shuffle 的字面意思为洗牌
&#x20; msg: "{{ testvar9 | shuffle }}"
&#x20; \- debug:
&#x20; \#随机打乱顺序列表中元素的顺序
&#x20; \#在随机打乱顺序时，将 ansible_date_time.epoch 的值设置为随机种子
&#x20; \#seed 参数从 ansible2.3 版本开始可用
&#x20; msg: "{{ testvar9 | shuffle(seed=(ansible\_date\_time.epoch)) }}"
&#x20; \- debug:
&#x20; \#将列表中的每个元素变成纯大写
&#x20; msg: "{{ testvar10 | upper }}"
&#x20; \- debug:
&#x20; \#将列表中的每个元素变成纯小写
&#x20; msg: "{{ testvar10 | lower }}"
&#x20; \- debug:
&#x20; \#去掉列表中重复的元素，重复的元素只留下一个
&#x20; msg: "{{ testvar11 | unique }}"
&#x20; \- debug:
&#x20; \#将两个列表合并，重复的元素只留下一个
&#x20; \#也就是求两个列表的并集
&#x20; msg: "{{ testvar11 | union(testvar12) }}"
&#x20; \- debug:
&#x20; \#取出两个列表的交集，重复的元素只留下一个
&#x20; msg: "{{ testvar11 | intersect(testvar12) }}"
&#x20; \- debug:
&#x20; \#取出存在于 testvar11 列表中,但是不存在于 testvar12 列表中的元素
&#x20; \#去重后重复的元素只留下一个
&#x20; \#换句话说就是:两个列表的交集在列表 1 中的补集
&#x20; msg: "{{ testvar11 | difference(testvar12) }}"
&#x20; \- debug:
&#x20; \#取出两个列表中各自独有的元素,重复的元素只留下一个
&#x20; \#即去除两个列表的交集，剩余的元素
&#x20; msg: "{{ testvar11 | symmetric\_difference(testvar12) }}" |

变量未定义时相关操作的过滤器，示例如下

Shell

|     |     |
| --- | --- |

| 1
2
3
4
5
6
7
8
9
10
11
12
13
14
15
16
17
18 | ---
\- hosts: test70
&#x20; remote_user: root
&#x20; gather_facts: no
&#x20; vars:
&#x20; testvar6: ''
&#x20; tasks:
&#x20; \- debug:
&#x20; \#如果变量没有定义，则临时返回一个指定的默认值
&#x20; \#注：如果定义了变量，变量值为空字符串，则会输出空字符
&#x20; \#default 过滤器的别名是 d
&#x20; msg: "{{ testvar5 | default('zsythink') }}"
&#x20; \- debug:
&#x20; \#如果变量的值是一个空字符串或者变量没有定义，则临时返回一个指定的默认值
&#x20; msg: "{{ testvar6 | default('zsythink',boolean=true) }}"
&#x20; \- debug:
&#x20; \#如果对应的变量未定义,则报出“Mandatory variable not defined.”错误，而不是报出默认错误
&#x20; msg: "{{ testvar5 | mandatory }}" |

其实，说到上例中的 default 过滤器，还有一个很方便的用法，default 过滤器不仅能在变量未定义时返回指定的值，还能够让模块的参数变得"可有可无"。

这样说不太容易理解，不如我们先来看一个工作场景，然后根据这个工作场景来描述所谓的"可有可无"，就容易理解多了，场景如下：

假设，我现在需要在目标主机上创建几个文件，这些文件大多数都不需要指定特定的权限，只有个别文件需要指定特定的权限，所以，在定义这些文件时，我将变量定义为了如下样子

Shell

|     |     |
| --- | --- |

| 1
2
3
4
5
6 | vars:
&#x20; paths:
&#x20; \- path: /tmp/testfile
&#x20; mode: '0444'
&#x20; \- path: /tmp/foo
&#x20; \- path: /tmp/bar |

如上所示，我一共定义了 3 个文件，只有第一个文件指定了权限，第二个文件和第三个文件没有指定任何权限，这样定义目的是，当这三个文件在目标主机中创建时，只有第一个文件按照指定的权限被创建，之后的两个文件都按照操作系统的默认权限进行创建，为了方便示例，我只定义了 3 个文件作为示例，但是在实际工作中，你获得列表中可能有几十个这样的文件需要被创建，这些文件中，有些文件需要特定的权限，有些不需要，所以，我们可能需要使用循环来处理这个问题，但是在使用循环时，我们会遇到另一个问题，问题就是，有的文件有 mode 属性，有的文件没有 mode 属性，那么，我们就需要对文件是否有 mode 属性进行判断，所以，你可能会编写一个类似如下结构的 playbook

Shell

|     |     |
| --- | --- |

| 1
2
3
4
5
6
7
8
9
10
11
12
13
14
15
16 | - hosts: test70
&#x20; remote_user: root
&#x20; gather_facts: no
&#x20; vars:
&#x20; paths:
&#x20; \- path: /tmp/test
&#x20; mode: '0444'
&#x20; \- path: /tmp/foo
&#x20; \- path: /tmp/bar
&#x20; tasks:
&#x20; \- file: dest={{item.path}} state=touch mode={{item.mode}}
&#x20; with_items: "{{ paths }}"
&#x20; when: item.mode is defined
&#x20; \- file: dest={{item.path}} state=touch
&#x20; with_items: "{{ paths }}"
&#x20; when: item.mode is undefined |

上例中，使用 file 模块在目标主机中创建文件，很好的解决我们的问题，但是上例中，我们一共循环了两遍，因为我们需要对文件是否有 mode 属性进行判断，然后根据判断结果调整 file 模块的参数设定，那么有没有更好的办法呢？当然有，这个办法就是我们刚才所说的"可有可无"，我们可以将上例 playbook 简化成如下模样：

Shell

|     |     |
| --- | --- |

| 1
2
3
4
5
6
7
8
9
10
11
12 | - hosts: test70
&#x20; remote_user: root
&#x20; gather_facts: no
&#x20; vars:
&#x20; paths:
&#x20; \- path: /tmp/test
&#x20; mode: '0444'
&#x20; \- path: /tmp/foo
&#x20; \- path: /tmp/bar
&#x20; tasks:
&#x20; \- file: dest={{item.path}} state=touch mode={{item.mode | default(omit)}}
&#x20; with_items: "{{ paths }}" |

上例中，我们并没有对文件是否有 mode 属性进行判断，而是直接调用了 file 模块的 mode 参数，将 mode 参数的值设定为了"{{item.mode | default(omit)}}"，这是什么意思呢？它的意思是，如果 item 有 mode 属性，就把 file 模块的 mode 参数的值设置为 item 的 mode 属性的值，如果 item 没有 mode 属性，file 模块就直接省略 mode 参数，'omit'的字面意思就是"省略"，换成大白话说就是：\[有就用，没有就不用，可以有，也可以没有]，所谓的"可有可无"就是这个意思，是不是很方便？我觉得聪明如你一定看懂了，快动手试试吧~

施主~~加油吧~~~这篇文章就总结到这里，希望能够对你有所帮助~掰掰~
